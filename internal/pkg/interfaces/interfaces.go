package interfaces

import (
	"context"

	"github.com/deneb-cygnus-dev/paper-analyzer/internal/pkg/entities"
)

// MetadataFetcher is the interface for fetching paper metadata
type MetadataFetcher interface {
	// Fetch fetches the metadata of the paper by the given configuration
	Fetch(ctx context.Context, config entities.FetchConfig) ([]entities.Paper, error)
}
