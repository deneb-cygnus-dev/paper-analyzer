# AGENTS.md for paper-analyzer

## Context

The "paper-analyzer" project is a polyglot application (Go and Python) designed to process and analyze research papers, specifically targeting metadata and content from arXiv.org.

## Goals

- Provide a robust API for fetching paper metadata.
- Analyze paper content using Python-based tools.
- Maintain a clean, modular, and testable codebase.

## Principles

- **Modularity**: Separation of concerns is paramount. Use interfaces to define contracts and separate implementation details.
- **Testability**: All implementations must be covered by unit tests.
- **Simplicity**: Keep the design simple and understandable.

## Conventions

### Coding Style

#### Go Style

- Follow standard Go conventions (Effective Go).
- **File Naming**: Use `snake_case` for file names (e.g., `arxiv_fetcher.go`).
- **Test File Naming**: Test files must match the source file name with `_test` suffix (e.g., `arxiv_fetcher_test.go` for `arxiv_fetcher.go`).
- **Interfaces**: Interfaces should be defined in a dedicated `interfaces` package (e.g., `internal/pkg/interfaces`) to avoid circular dependencies and promote decoupling.
- **Implementations**: Concrete implementations should be in their own packages (e.g., `internal/pkg/fetcher`).

### Python Style

- Follow PEP 8 guidelines.
- Use virtual environments for dependency management.

### Project Structure

- `cmd/`: Main Go application entry points.
- `internal/`: Private Go packages.
  - `pkg/entities`: Data structures (e.g., `Paper`, `Author`).
  - `pkg/interfaces`: Interface definitions (e.g., `MetadataFetcher`).
  - `pkg/fetcher`: Concrete implementations (e.g., `ArxivFetcher`).
- `python/`: Python source code.
- `ai-docs/`: Documentation for AI agents.
- `mise.local.toml`: Tool version management.

### Documentation Style

- **File Naming**: Use `kebab-case` for documentation files (e.g., `get-metadata-api.md`).
- **Linting**: Markdown files should strictly follow the convention of `markdownlint`.
- **Structure**:
  - **Title**: H1 header.
  - **Overview**: Brief description of the document's purpose.
  - **Architecture**: (Optional) Package structure and relationships.
  - **API Reference**: Details about interfaces, structs, and functions.
  - **Usage**: Code examples showing how to use the component.
  - **Testing**: Information about tests and how to run them.
  - **Error Handling**: Details about specific errors and codes.
- **Code Blocks**: Always specify the language (e.g., `go`, `bash`, `text`).
- **Formatting**:
  - **Trailing Spaces**: Remove trailing spaces from all lines.
  - **Ordered Lists**: Start all ordered list elements with "1." (not "1. 2. 3. ...").

### Documentation Updates

- **Post-Implementation Update**: Whenever a branch is completely implemented, the documentation in `ai-docs/` must be scanned and updated to reflect the newly implemented features.

## Tech Stack

- **Go**: Version 1.25.1
- **Python**: Version 3.14
- **Version Management**: `mise`

## Building and Running

### Go Building, Testing, and Running

```bash
# Build the Go application
go build ./...

# Run tests
go test ./...
```

### Python Building, Testing, and Running

```bash
# Setup virtual environment
python -m venv .venv
source .venv/bin/activate

# Install dependencies
pip install -r python/requirements.txt

# Run main script
python python/main.py
```
