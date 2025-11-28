package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/errors"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/interfaces"
)

// ArxivDownloader implements the PDFDownloader interface
type ArxivDownloader struct {
	downloadDir string
}

// Ensure ArxivDownloader implements PDFDownloader
var _ interfaces.PDFDownloader = (*ArxivDownloader)(nil)

// NewArxivDownloader creates a new ArxivDownloader
func NewArxivDownloader(downloadDir string) *ArxivDownloader {
	return &ArxivDownloader{
		downloadDir: downloadDir,
	}
}

// Download implements the PDFDownloader interface
func (d *ArxivDownloader) Download(ctx context.Context, papers []entities.Paper) (map[string]string, map[string]error) {
	downloadedPaths := make(map[string]string)
	downloadErrors := make(map[string]error)

	for _, paper := range papers {
		pdfLink := d.findPDFLink(paper)
		if pdfLink == "" {
			downloadErrors[paper.ID] = errors.Wrap(fmt.Errorf("paper %s has no PDF link", paper.ID), errors.ErrPaperDownload)
			continue
		}

		filePath, err := d.downloadPaper(ctx, pdfLink, paper.ID, d.downloadDir)
		if err != nil {
			downloadErrors[paper.ID] = errors.Wrap(err, errors.ErrPaperDownload)
			continue
		}

		downloadedPaths[paper.ID] = filePath
	}

	return downloadedPaths, downloadErrors
}

func (d *ArxivDownloader) findPDFLink(paper entities.Paper) string {
	for _, link := range paper.Links {
		if link.Type == "application/pdf" {
			return link.Href
		}
	}
	return ""
}

func (d *ArxivDownloader) downloadPaper(ctx context.Context, url, paperID, dirPath string) (string, error) {
	// Create the request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Determine filename from paper ID
	// arXiv IDs are like "http://arxiv.org/abs/2101.12345v1" or just "2101.12345v1"
	// We want the last part
	parts := strings.Split(paperID, "/")
	filename := parts[len(parts)-1] + ".pdf"
	filePath := filepath.Join(dirPath, filename)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	return filePath, nil
}
