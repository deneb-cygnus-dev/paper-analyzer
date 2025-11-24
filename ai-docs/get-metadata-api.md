# Get Metadata API

This document describes the implementation and usage of the Get Metadata API for the paper analyzer.

## Overview

The Get Metadata API abstracts the process of fetching paper metadata from external sources (currently limited to arXiv.org). It provides a unified interface and data structures for handling paper information.

## Architecture

The implementation follows a modular design with clear separation of concerns:

- **Entities**: Define the data structures for papers and related information.
- **Interfaces**: Define the contract for fetching metadata.
- **Fetcher**: Provides the concrete implementation for specific sources (arXiv).

### Package Structure

```text
internal/pkg/
├── entities/
│   └── entities.go       # Data structures (Paper, Author, Link)
├── interfaces/
│   └── interfaces.go     # MetadataFetcher interface
└── fetcher/
    ├── arxiv_fetcher.go  # ArxivFetcher implementation
    └── arxiv_fetcher_test.go
```

## API Reference

### Entities (`internal/pkg/entities`)

The core data structure is the `Paper` struct:

```go
type Paper struct {
    ID          string    `json:"id"`           // e.g., http://arxiv.org/abs/2511.17464v1
    Title       string    `json:"title"`
    Summary     string    `json:"summary"`      // Abstract
    Authors     []Author  `json:"author"`
    PublishDate time.Time `json:"publish_date"`
    UpdatedDate time.Time `json:"updated_date"`
    Links       []Link    `json:"links"`
    Categories  []string  `json:"categories"`
}
```

### Interface (`internal/pkg/interfaces`)

The `MetadataFetcher` interface defines the contract:

```go
type MetadataFetcher interface {
    // Fetch fetches the metadata of the paper by the given query
    Fetch(ctx context.Context, query string) ([]entities.Paper, error)
}
```

### Implementation (`internal/pkg/fetcher`)

The `ArxivFetcher` implements `MetadataFetcher` for arXiv.org.

#### Initialization

```go
// Create a new ArxivFetcher with a custom HTTP client (optional)
client := &http.Client{Timeout: 10 * time.Second}
fetcher := fetcher.NewArxivFetcher(client)

// Or use default client
fetcher := fetcher.NewArxivFetcher(nil)
```

#### Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/fetcher"
)

func main() {
    f := fetcher.NewArxivFetcher(nil)
    
    // Query arXiv API
    // See https://arxiv.org/help/api/user-manual#query_details for query format
    query := "http://export.arxiv.org/api/query?search_query=cat:cs.SE&max_results=1"
    
    papers, err := f.Fetch(context.Background(), query)
    if err != nil {
        log.Fatalf("Failed to fetch papers: %v", err)
    }

    for _, p := range papers {
        fmt.Printf("Title: %s\n", p.Title)
        fmt.Printf("ID: %s\n", p.ID)
        fmt.Printf("Abstract: %s\n", p.Summary)
    }
}
```

## Testing

Unit tests are located in `internal/pkg/fetcher/arxiv_fetcher_test.go`.
Run tests with:

```bash
go test ./internal/pkg/fetcher/...
```
