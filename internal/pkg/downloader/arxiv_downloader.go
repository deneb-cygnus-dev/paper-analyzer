package downloader

import (
	"context"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
)

// ArxivDownloader implements the PDFDownloader interface
type ArxivDownloader struct {
}

// Download implements the PDFDownloader interface
func (d *ArxivDownloader) Download(ctx context.Context, papers []entities.Paper) ([]string, error) {
	// TODO: check the Links field of the paper and determine
	// if the Type of "application/pdf" exists
	// if not, return error
	// if exists, download the PDF file (go forward)
	return nil, nil
}
