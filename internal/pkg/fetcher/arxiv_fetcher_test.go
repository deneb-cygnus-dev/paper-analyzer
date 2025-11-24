package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestArxivFetcher_Fetch(t *testing.T) {
	mockResponse := `
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry>
    <id>http://arxiv.org/abs/2511.17464v1</id>
    <title>A Patient-Centric Blockchain Framework</title>
    <summary>We present a patient-centric architecture...</summary>
    <published>2025-11-21T18:09:25Z</published>
    <updated>2025-11-21T18:09:25Z</updated>
    <author>
      <name>Tanzim Hossain Romel</name>
    </author>
    <link href="https://arxiv.org/abs/2511.17464v1" rel="alternate" type="text/html"/>
    <category term="cs.CR"/>
  </entry>
</feed>
`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		q := r.URL.Query()
		assert.Contains(t, q.Get("search_query"), "cat:cs.CR")
		assert.Equal(t, "submittedDate", q.Get("sortBy"))
		assert.Equal(t, "descending", q.Get("sortOrder"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())
	fetcher.baseURL = server.URL + "?" // Override base URL for testing

	config := entities.FetchConfig{
		Category:   "cs.CR",
		MaxResults: 10,
	}

	papers, err := fetcher.Fetch(context.Background(), config)
	assert.NoError(t, err)
	assert.Len(t, papers, 1)

	paper := papers[0]
	assert.Equal(t, "http://arxiv.org/abs/2511.17464v1", paper.ID)
	assert.Equal(t, "A Patient-Centric Blockchain Framework", paper.Title)
	assert.Equal(t, "We present a patient-centric architecture...", paper.Summary)
	assert.Equal(t, "Tanzim Hossain Romel", paper.Authors[0].Name)
	assert.Equal(t, "https://arxiv.org/abs/2511.17464v1", paper.Links[0].Href)
	assert.Equal(t, "cs.CR", paper.Categories[0])

	expectedTime, _ := time.Parse(time.RFC3339, "2025-11-21T18:09:25Z")
	assert.Equal(t, expectedTime, paper.PublishDate)
}

func TestArxivFetcher_Fetch_WithKeywords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Contains(t, q.Get("search_query"), "cat:cs.SE")
		assert.Contains(t, q.Get("search_query"), "AND all:fuzzing")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<feed></feed>`))
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())
	fetcher.baseURL = server.URL + "?"

	config := entities.FetchConfig{
		Category:   "cs.SE",
		MaxResults: 5,
		Keywords:   []string{"fuzzing"},
	}

	_, err := fetcher.Fetch(context.Background(), config)
	assert.NoError(t, err)
}

func TestArxivFetcher_Fetch_WithTimeSpan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Contains(t, q.Get("search_query"), "cat:cs.AI")
		assert.Contains(t, q.Get("search_query"), "AND submittedDate:[")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<feed></feed>`))
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())
	fetcher.baseURL = server.URL + "?"

	config := entities.FetchConfig{
		Category: "cs.AI",
		TimeSpan: "last_7_days",
	}

	_, err := fetcher.Fetch(context.Background(), config)
	assert.NoError(t, err)
}

func TestArxivFetcher_Fetch_ValidationErrors(t *testing.T) {
	fetcher := NewArxivFetcher(nil)

	// Missing Category
	_, err := fetcher.Fetch(context.Background(), entities.FetchConfig{
		MaxResults: 10,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category is required")

	// Missing Limit (TimeSpan or MaxResults)
	_, err = fetcher.Fetch(context.Background(), entities.FetchConfig{
		Category: "cs.LG",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either TimeSpan or MaxResults must be specified")
}

func TestArxivFetcher_Fetch_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())
	fetcher.baseURL = server.URL + "?"

	config := entities.FetchConfig{
		Category:   "cs.CR",
		MaxResults: 1,
	}

	_, err := fetcher.Fetch(context.Background(), config)
	assert.Error(t, err)
}

func TestArxivFetcher_Fetch_EndToEnd(t *testing.T) {
	fetcher := NewArxivFetcher(nil)

	config := entities.FetchConfig{
		Category:   "cs.SE",
		MaxResults: 5,
	}

	papers, err := fetcher.Fetch(context.Background(), config)
	assert.NoError(t, err)
	assert.Len(t, papers, 5)
}
