# Project Overview

This project, "paper-analyzer," appears to be a polyglot application utilizing both Go and Python. The directory structure suggests a standard Go layout (`cmd/`, `internal/`) for the Go components and a dedicated `python/` directory for the Python parts.

The specific purpose of the project is not detailed in the `README.md`, but the name "paper-analyzer" suggests it might be a tool for processing and analyzing documents or research papers.

**Technologies:**

*   **Go:** Version 1.25.1
*   **Python:** Version 3.14

**Project Structure:**

*   `cmd/`: Likely contains the main Go application entry points.
*   `internal/`: Likely contains private Go packages for the application.
*   `python/`: Contains the Python source code.
*   `mise.toml`: Defines the project's tool versions (Go and Python).

# Building and Running

**Go:**

```bash
# TODO: Verify the build and run commands for the Go application.
# Build the Go application
go build ./...

# Run the main application (assuming a main package in cmd/paper-analyzer)
go run ./cmd/paper-analyzer
```

**Python:**

```bash
# TODO: Verify the setup and execution for the Python part.
# It's recommended to use a virtual environment.
python -m venv .venv
source .venv/bin/activate

# Install dependencies (assuming a requirements.txt file)
pip install -r python/requirements.txt

# Run the main Python script (assuming a main.py)
python python/main.py
```

# Development Conventions

*   Go code seems to follow the standard project layout (`cmd/`, `internal/`).
*   Dependency versions for Go and Python are managed via `mise.toml`.
