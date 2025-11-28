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
│   └── entities.go       # Data structures (Paper, Author, Link, FetchConfig)
├── interfaces/
│   └── interfaces.go     # MetadataFetcher interface
└── fetcher/
    ├── arxiv_fetcher.go  # ArxivFetcher implementation
    └── arxiv_fetcher_test.go
```

## API Reference

### Entities (`internal/pkg/entities`)

#### `Paper`

The core data structure representing a paper:

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

#### `FetchConfig`

Configuration for fetching papers:

```go
type FetchConfig struct {
    // Category to search for (e.g., "cs.SE")
    // Required.
    Category string

    // TimeSpan to filter papers (e.g., "last_5_days")
    // Mutually inclusive with MaxResults (at least one required).
    TimeSpan string

    // MaxResults to limit the number of papers
    // Mutually inclusive with TimeSpan.
    MaxResults int

    // Keywords to search for
    // Optional.
    Keywords []string
}
```

### Interface (`internal/pkg/interfaces`)

The `MetadataFetcher` interface defines the contract:

```go
type MetadataFetcher interface {
    // Fetch fetches the metadata of the paper by the given configuration
    Fetch(ctx context.Context, config entities.FetchConfig) ([]entities.Paper, error)
}
```

#### Error Handling

The `Fetch` method returns `CustomError` types defined in `internal/pkg/errors`. Common errors include:

- **Validation Errors**:
  - `ErrMissingRequiredField` (400002): Returned when `Category` is missing.
  - `ErrInvalidInput` (400001): Returned when neither `TimeSpan` nor `MaxResults` is specified.
- **Infrastructure Errors**:
  - `ErrNetwork` (500004): Returned when network communication fails.
  - `ErrExternalAPI` (500006): Returned when the arXiv API returns a non-200 status code.
  - `ErrExternalAPIParsing` (500007): Returned when the response XML cannot be parsed.
- **Internal Errors**:
  - `ErrInternalServer` (100001): Returned for URL parsing or request creation failures.

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

    "github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
    "github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/fetcher"
)

func main() {
    f := fetcher.NewArxivFetcher(nil)
    
    // Configure fetch parameters
    config := entities.FetchConfig{
        Category:   "cs.SE",
        MaxResults: 5,
        Keywords:   []string{"fuzzing", "testing"},
    }
    
    papers, err := f.Fetch(context.Background(), config)
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
