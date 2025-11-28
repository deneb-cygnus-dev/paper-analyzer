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

	downloader := &ArxivDownloader{}

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
	paths, err := downloader.Download(context.Background(), papers, tempDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
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

	downloader := &ArxivDownloader{}

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

	_, err = downloader.Download(context.Background(), []entities.Paper{paper}, tempDir)
	if err == nil {
		t.Error("Expected error, got nil")
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

	downloader := &ArxivDownloader{}

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
	paths, err := downloader.Download(context.Background(), []entities.Paper{paper}, tempDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("Expected 1 path, got %d", len(paths))
	}

	downloadedPath := paths[0]

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
	// Note: PDF files might have different metadata (creation date, etc.) if downloaded at different times.
	// However, the user asked to compare them. If they are identical byte-by-byte, great.
	// If not, we might need to be more lenient or check PDF structure, but let's try byte comparison first.
	// Wait, the artifact is likely just a manually downloaded copy.
	// If the server serves the exact same file every time, it should match.
	// But sometimes dynamic watermarks or timestamps are added.
	// Let's assume exact match for now as per request.

	// Actually, let's just check if the size is roughly the same or if it's a valid PDF.
	// But the user said "compare it with already downloaded file".
	// Let's try byte comparison.

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
