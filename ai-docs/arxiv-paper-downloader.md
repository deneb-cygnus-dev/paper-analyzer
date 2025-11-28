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
 // Download downloads the PDF files from the given papers and returns the paths of the downloaded PDF files.
 // If failed, returns an error.
 Download(ctx context.Context, papers []entities.Paper, downloadDirPath string) ([]string, error)
}
```

### Implementation (`internal/pkg/downloader`)

The `ArxivDownloader` struct implements the `PDFDownloader` interface.

```go
// ArxivDownloader implements the PDFDownloader interface
type ArxivDownloader struct {
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
 // Initialize downloader
 d := &downloader.ArxivDownloader{}

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

 // Create download directory
 downloadDir := "./downloads"
 if err := os.MkdirAll(downloadDir, 0755); err != nil {
  log.Fatal(err)
 }

 // Download papers
 paths, err := d.Download(context.Background(), papers, downloadDir)
 if err != nil {
  log.Fatalf("Failed to download papers: %v", err)
 }

 for _, path := range paths {
  fmt.Printf("Downloaded to: %s\n", path)
 }
}
```

## Testing

Unit and integration tests are located in `internal/pkg/downloader/arxiv_downloader_test.go`.

The tests cover:

1. **Successful Download**: Verifies that a real paper can be downloaded from arXiv.
2. **No PDF Link**: Verifies that an error is returned if the paper has no PDF link.
3. **Comparison Test**: Downloads a specific paper and compares it byte-by-byte with a pre-downloaded artifact to ensure integrity.

Run tests with:

```bash
go test -v ./internal/pkg/downloader/...
```
