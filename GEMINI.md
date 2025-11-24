# Context

The "paper-analyzer" project is a polyglot application (Go and Python) designed to process and analyze research papers, specifically targeting metadata and content from arXiv.org.

# Goals

- Provide a robust API for fetching paper metadata.
- Analyze paper content using Python-based tools.
- Maintain a clean, modular, and testable codebase.

# Principles

- **Modularity**: Separation of concerns is paramount. Use interfaces to define contracts and separate implementation details.
- **Testability**: All implementations must be covered by unit tests.
- **Simplicity**: Keep the design simple and understandable.

# Conventions

## Coding Style

### Go
- Follow standard Go conventions (Effective Go).
- **File Naming**: Use `snake_case` for file names (e.g., `arxiv_fetcher.go`).
- **Test File Naming**: Test files must match the source file name with `_test` suffix (e.g., `arxiv_fetcher_test.go` for `arxiv_fetcher.go`).
- **Interfaces**: Interfaces should be defined in a dedicated `interfaces` package (e.g., `internal/pkg/interfaces`) to avoid circular dependencies and promote decoupling.
- **Implementations**: Concrete implementations should be in their own packages (e.g., `internal/pkg/fetcher`).

### Python
- Follow PEP 8 guidelines.
- Use virtual environments for dependency management.

## Project Structure

- `cmd/`: Main Go application entry points.
- `internal/`: Private Go packages.
    - `pkg/entities`: Data structures (e.g., `Paper`, `Author`).
    - `pkg/interfaces`: Interface definitions (e.g., `MetadataFetcher`).
    - `pkg/fetcher`: Concrete implementations (e.g., `ArxivFetcher`).
- `python/`: Python source code.
- `ai-docs/`: Documentation for AI agents.
- `mise.local.toml`: Tool version management.

# Tech Stack

- **Go**: Version 1.25.1
- **Python**: Version 3.14
- **Version Management**: `mise`

# Building and Running

## Go

```bash
# Build the Go application
go build ./...

# Run tests
go test ./...
```

## Python

```bash
# Setup virtual environment
python -m venv .venv
source .venv/bin/activate

# Install dependencies
pip install -r python/requirements.txt

# Run main script
python python/main.py
```
