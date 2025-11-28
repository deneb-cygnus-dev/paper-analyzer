package interfaces

import (
	"context"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
)

// MetadataFetcher is the interface for fetching paper metadata
type MetadataFetcher interface {
	// Fetch fetches the metadata of the paper by the given configuration
	// Parameters:
	//   - ctx: the context
	//   - config: the configuration for fetching paper metadata
	// Returns:
	//   - papers: the metadata of the paper
	//   - error: the error if any
	Fetch(ctx context.Context, config entities.FetchConfig) ([]entities.Paper, error)
}

// PDFDownloader is the interface for downloading PDF files
type PDFDownloader interface {
	// Download downloads the PDF file of the paper by the given configuration
	// Parameters:
	//   - ctx: the context
	//   - papers: the papers to download PDF files
	// Returns:
	//   - paths: the paths of the downloaded PDF files
	//   - error: the error if any
	Download(ctx context.Context, papers []entities.Paper) ([]string, error)
}
