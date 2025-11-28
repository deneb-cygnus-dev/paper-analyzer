package downloader

import (
	"context"
	std_errors "errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/errors"
	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/fetcher"
)

func TestArxivDownloader_Download_Success(t *testing.T) {
	// Setup temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "arxiv_download_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloader := NewArxivDownloader(tempDir)

	// Use real fetcher to get a valid paper
	f := fetcher.NewArxivFetcher(http.DefaultClient)
	config := entities.FetchConfig{
		Category:   "cs.SE",
		MaxResults: 1,
	}
	papers, err := f.Fetch(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to fetch papers: %v", err)
	}
	if len(papers) == 0 {
		t.Fatalf("No papers fetched")
	}

	// Perform download
	paths, downloadErrors := downloader.Download(context.Background(), papers)
	if len(downloadErrors) > 0 {
		t.Fatalf("Download failed with errors: %v", downloadErrors)
	}

	// Verify results
	if len(paths) != len(papers) {
		t.Errorf("Expected %d paths, got %d", len(papers), len(paths))
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("File not found at %s: %v", path, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("File at %s is empty", path)
		}
		// Check if file extension is .pdf
		if filepath.Ext(path) != ".pdf" {
			t.Errorf("File at %s does not have .pdf extension", path)
		}
	}
}

func TestArxivDownloader_Download_NoPDFLink(t *testing.T) {
	// Setup temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "arxiv_download_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloader := NewArxivDownloader(tempDir)

	// Create a dummy paper without PDF link
	paper := entities.Paper{
		ID:    "http://arxiv.org/abs/dummy",
		Title: "Dummy Paper",
		Links: []entities.Link{
			{
				Href: "http://arxiv.org/abs/dummy",
				Rel:  "alternate",
				Type: "text/html",
			},
		},
	}

	paths, downloadErrors := downloader.Download(context.Background(), []entities.Paper{paper})
	if len(paths) != 0 {
		t.Errorf("Expected 0 paths, got %d", len(paths))
	}

	if len(downloadErrors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(downloadErrors))
	}

	err = downloadErrors[paper.ID]
	if err == nil {
		t.Errorf("Expected error for paper ID '%s', got nil", paper.ID)
	}

	// Verify error type
	var customErr *errors.CustomError
	if !std_errors.As(err, &customErr) {
		t.Errorf("Expected *errors.CustomError, got %T", err)
	}

	if customErr != nil && customErr.Code != errors.ErrPaperDownload.Code {
		t.Errorf("Expected error code %d, got %d", errors.ErrPaperDownload.Code, customErr.Code)
	}
}

func TestArxivDownloader_Download_Compare(t *testing.T) {
	// Setup temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "arxiv_download_compare_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloader := NewArxivDownloader(tempDir)

	// Target paper
	paperID := "2110.06449"
	paperURL := "http://arxiv.org/abs/" + paperID
	pdfURL := "http://arxiv.org/pdf/" + paperID

	paper := entities.Paper{
		ID:    paperURL,
		Title: "Constrained Detecting Arrays: Mathematical Structures for Fault Identification in Combinatorial Interaction Testing",
		Links: []entities.Link{
			{
				Href: pdfURL,
				Rel:  "related",
				Type: "application/pdf",
			},
		},
	}

	// Perform download
	paths, downloadErrors := downloader.Download(context.Background(), []entities.Paper{paper})
	if len(downloadErrors) > 0 {
		t.Fatalf("Download failed with errors: %v", downloadErrors)
	}

	if len(paths) != 1 {
		t.Fatalf("Expected 1 path, got %d", len(paths))
	}

	downloadedPath := paths[paper.ID]

	// Read downloaded file
	downloadedContent, err := os.ReadFile(downloadedPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	// Read artifact file
	// Note: The path is relative to the project root, but tests run in the package directory.
	// We need to go up 3 levels: internal/pkg/downloader -> internal/pkg -> internal -> root
	artifactPath := "../../../testdata/artifacts/Constrained Detecting Arrays.pdf"
	artifactContent, err := os.ReadFile(artifactPath)
	if err != nil {
		t.Fatalf("Failed to read artifact file: %v", err)
	}

	// Compare contents
	if len(downloadedContent) != len(artifactContent) {
		t.Errorf("Downloaded file size (%d) does not match artifact size (%d)", len(downloadedContent), len(artifactContent))
	} else {
		// Only check bytes if sizes match to avoid massive output on failure
		for i := range downloadedContent {
			if downloadedContent[i] != artifactContent[i] {
				t.Errorf("Downloaded file content mismatch at byte %d", i)
				break // Stop after first mismatch
			}
		}
	}
}

func TestArxivDownloader_Download_PartialFailure(t *testing.T) {
	// Setup temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "arxiv_download_partial_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloader := NewArxivDownloader(tempDir)

	// Create one valid paper and one invalid paper
	validPaper := entities.Paper{
		ID:    "http://arxiv.org/abs/2110.06449",
		Title: "Valid Paper",
		Links: []entities.Link{
			{
				Href: "http://arxiv.org/pdf/2110.06449",
				Rel:  "related",
				Type: "application/pdf",
			},
		},
	}

	invalidPaper := entities.Paper{
		ID:    "http://arxiv.org/abs/dummy",
		Title: "Invalid Paper",
		Links: []entities.Link{
			{
				Href: "http://arxiv.org/abs/dummy",
				Rel:  "alternate",
				Type: "text/html",
			},
		},
	}

	paths, downloadErrors := downloader.Download(context.Background(), []entities.Paper{validPaper, invalidPaper})

	// Check valid paper
	if _, ok := paths[validPaper.ID]; !ok {
		t.Error("Expected valid paper to be downloaded")
	}
	if _, ok := downloadErrors[validPaper.ID]; ok {
		t.Error("Expected no error for valid paper")
	}

	// Check invalid paper
	if _, ok := paths[invalidPaper.ID]; ok {
		t.Error("Expected invalid paper not to be in paths")
	}
	if _, ok := downloadErrors[invalidPaper.ID]; !ok {
		t.Error("Expected error for invalid paper")
	}
}
