package fetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/errors"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/interfaces"
)

// ArxivFetcher implements MetadataFetcher for arXiv.org
type ArxivFetcher struct {
	client  *http.Client
	baseURL string
}

// Ensure ArxivFetcher implements MetadataFetcher
var _ interfaces.MetadataFetcher = (*ArxivFetcher)(nil)

// NewArxivFetcher creates a new ArxivFetcher
func NewArxivFetcher(client *http.Client) *ArxivFetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &ArxivFetcher{
		client:  client,
		baseURL: "http://export.arxiv.org/api/query?",
	}
}

// Fetch fetches the metadata of the paper by the given configuration
func (f *ArxivFetcher) Fetch(ctx context.Context, config entities.FetchConfig) ([]entities.Paper, error) {
	queryURL, err := f.buildQueryURL(config)
	if err != nil {
		return nil, err
	}

	req, err := f.buildRequest(ctx, queryURL)
	if err != nil {
		return nil, err
	}

	body, err := f.doRequest(req)
	if err != nil {
		return nil, err
	}

	return f.parseResponse(body)
}

func (f *ArxivFetcher) buildQueryURL(config entities.FetchConfig) (string, error) {
	if config.Category == "" {
		return "", errors.ErrMissingRequiredField
	}
	if config.TimeSpan == "" && config.MaxResults == 0 {
		return "", errors.ErrInvalidInput
	}

	// Build search query
	searchQuery := fmt.Sprintf("cat:%s", config.Category)
	if len(config.Keywords) > 0 {
		for _, kw := range config.Keywords {
			searchQuery += fmt.Sprintf(" AND all:%s", kw)
		}
	}

	// Handle TimeSpan if specified (e.g., "last_5_days")
	if config.TimeSpan != "" {
		var days int
		if _, err := fmt.Sscanf(config.TimeSpan, "last_%d_days", &days); err == nil {
			// Calculate start date
			startDate := time.Now().AddDate(0, 0, -days)
			// Format: YYYYMMDDHHMM
			startStr := startDate.Format("200601021504")
			// Append to search query: submittedDate:[START TO *]
			// Note: arXiv API uses "submittedDate" for submission time
			searchQuery += fmt.Sprintf(" AND submittedDate:[%s0000 TO *]", startStr)
		}
	}

	// Use url.Values to encode parameters
	v := url.Values{}
	v.Set("search_query", searchQuery)
	v.Set("sortBy", "submittedDate")
	v.Set("sortOrder", "descending")

	if config.MaxResults > 0 {
		v.Set("max_results", fmt.Sprintf("%d", config.MaxResults))
	}

	// Let's parse the baseURL
	u, err := url.Parse(f.baseURL)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternalServer)
	}

	// Add new params to existing query params (if any)
	q := u.Query()
	for key, values := range v {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (f *ArxivFetcher) buildRequest(ctx context.Context, queryURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer)
	}
	return req, nil
}

func (f *ArxivFetcher) doRequest(req *http.Request) ([]byte, error) {
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNetwork)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(errors.ErrExternalAPI.Code, fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNetwork)
	}
	return body, nil
}

func (f *ArxivFetcher) parseResponse(body []byte) ([]entities.Paper, error) {
	var feed atomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, errors.Wrap(err, errors.ErrExternalAPIParsing)
	}

	var papers []entities.Paper
	for _, entry := range feed.Entry {
		paper := entities.Paper{
			ID:          entry.ID,
			Title:       entry.Title,
			Summary:     entry.Summary,
			PublishDate: entry.Published,
			UpdatedDate: entry.Updated,
		}

		for _, author := range entry.Author {
			paper.Authors = append(paper.Authors, entities.Author{
				Name: author.Name,
			})
		}

		for _, link := range entry.Link {
			paper.Links = append(paper.Links, entities.Link{
				Href: link.Href,
				Rel:  link.Rel,
				Type: link.Type,
			})
		}

		for _, cat := range entry.Category {
			paper.Categories = append(paper.Categories, cat.Term)
		}

		papers = append(papers, paper)
	}

	return papers, nil
}

// Internal structures for XML parsing

type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entry   []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID        string         `xml:"id"`
	Title     string         `xml:"title"`
	Summary   string         `xml:"summary"`
	Published time.Time      `xml:"published"`
	Updated   time.Time      `xml:"updated"`
	Author    []atomAuthor   `xml:"author"`
	Link      []atomLink     `xml:"link"`
	Category  []atomCategory `xml:"category"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type atomCategory struct {
	Term string `xml:"term,attr"`
}
