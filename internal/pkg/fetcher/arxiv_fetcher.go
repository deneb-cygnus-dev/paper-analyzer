package fetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/interfaces"
)

// ArxivFetcher implements MetadataFetcher for arXiv.org
type ArxivFetcher struct {
	client *http.Client
}

// Ensure ArxivFetcher implements MetadataFetcher
var _ interfaces.MetadataFetcher = (*ArxivFetcher)(nil)

// NewArxivFetcher creates a new ArxivFetcher
func NewArxivFetcher(client *http.Client) *ArxivFetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &ArxivFetcher{
		client: client,
	}
}

// Fetch fetches the metadata of the paper by the given query
func (f *ArxivFetcher) Fetch(ctx context.Context, query string) ([]entities.Paper, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed atomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
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
