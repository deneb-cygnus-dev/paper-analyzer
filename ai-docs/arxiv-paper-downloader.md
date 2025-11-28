# Arxiv Paper Downloader

This document describes the implementation and usage of the Arxiv Paper Downloader for the paper analyzer.

## Overview

The Arxiv Paper Downloader is responsible for downloading PDF files of papers from arXiv.org. It implements the `PDFDownloader` interface and handles the retrieval and storage of paper PDFs.

## Architecture

The implementation is part of the `downloader` package and interacts with the `entities` package.

### Package Structure

```text
internal/pkg/
├── downloader/
│   ├── arxiv_downloader.go       # ArxivDownloader implementation
│   └── arxiv_downloader_test.go  # Tests
└── entities/
    └── entities.go               # Paper entity definition
```

## API Reference

### Interface (`internal/pkg/downloader`)

The `PDFDownloader` interface defines the contract for downloading papers:

```go
// PDFDownloader is an interface for downloading PDF files.
type PDFDownloader interface {
 // Download downloads the PDF files from the given papers.
 // Returns:
 //   - paths: the paths of the downloaded PDF files, a map from paper ID to path
 //   - errors: a map from paper ID to error
 Download(ctx context.Context, papers []entities.Paper) (map[string]string, map[string]error)
}
```

### Implementation (`internal/pkg/downloader`)

The `ArxivDownloader` struct implements the `PDFDownloader` interface.

```go
// ArxivDownloader implements the PDFDownloader interface
type ArxivDownloader struct {
 downloadDir string
}

// NewArxivDownloader creates a new ArxivDownloader
func NewArxivDownloader(downloadDir string) *ArxivDownloader {
 return &ArxivDownloader{
  downloadDir: downloadDir,
 }
}
```

#### Error Handling

The downloader uses the internal error handling system (`internal/pkg/errors`). Common errors include:

- `ErrPaperDownload` (600001): Returned when a paper fails to download or has no PDF link.

### Usage Example

```go
package main

import (
 "context"
 "fmt"
 "log"
 "os"

 "github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/downloader"
 "github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
)

func main() {
 // Create download directory
 downloadDir := "./downloads"
 if err := os.MkdirAll(downloadDir, 0755); err != nil {
  log.Fatal(err)
 }

 // Initialize downloader
 d := downloader.NewArxivDownloader(downloadDir)

 // Prepare papers (usually fetched via ArxivFetcher)
 papers := []entities.Paper{
  {
   ID: "http://arxiv.org/abs/2110.06449",
   Links: []entities.Link{
    {
     Href: "http://arxiv.org/pdf/2110.06449",
     Type: "application/pdf",
    },
   },
  },
 }

 // Download papers
 paths, errs := d.Download(context.Background(), papers)
 
 // Handle errors
 if len(errs) > 0 {
  for id, err := range errs {
   log.Printf("Failed to download paper %s: %v", id, err)
  }
 }

 // Handle successful downloads
 for id, path := range paths {
  fmt.Printf("Paper %s downloaded to: %s\n", id, path)
 }
}
```

## Testing

Unit and integration tests are located in `internal/pkg/downloader/arxiv_downloader_test.go`.

The tests cover:

1. **Successful Download**: Verifies that a real paper can be downloaded from arXiv.
2. **No PDF Link**: Verifies that an error is returned if the paper has no PDF link.
3. **Compare Test**: Downloads a specific paper and compares it byte-by-byte with a pre-downloaded artifact to ensure integrity.
4. **Partial Failure**: Verifies that the downloader continues to download other papers even if one fails, and correctly reports both successes and errors.

Run tests with:

```bash
go test -v ./internal/pkg/downloader/...
```
