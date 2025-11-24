package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())

	papers, err := fetcher.Fetch(context.Background(), server.URL)
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

func TestArxivFetcher_Fetch_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := NewArxivFetcher(server.Client())

	_, err := fetcher.Fetch(context.Background(), server.URL)
	assert.Error(t, err)
}
