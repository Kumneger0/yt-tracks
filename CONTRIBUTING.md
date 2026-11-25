# Contributing to yt-tracks

Thank you for your interest in contributing to **yt-tracks**! This document provides guidelines and setup instructions to help you get started.

---

## Architecture Overview

`yt-tracks` consists of two main components:
- **Go TUI Frontend**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), Lipgloss, and Cobra.
- **Python Backend Service**: A gRPC / Connect RPC server managing YouTube Music API interactions, authentication, and metadata.

---

## Prerequisites

Before starting, ensure you have the following installed:

- **Go**: `1.26` or higher
- **Python**: `3.13` or higher
- **uv**: Fast Python package manager ([install guide](https://docs.astral.sh/uv/getting-started/installation/))
- **[buf](https://buf.build/)**: Protobuf code generation (required if modifying `.proto` files or building the backend)
- **System Dependencies**:
  - **[yt-dlp](https://github.com/yt-dlp/yt-dlp)**: Required for extracting and fetching YouTube Music audio stream URLs
  - **[ffmpeg](https://ffmpeg.org/)**: Required for audio stream decoding and playback

---

## Development Setup

1. **Fork & Clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/yt-tracks.git
   cd yt-tracks
   ```

2. **Install Python Environment**:
   ```bash
   uv sync
   ```

3. **Install Git Hooks** *(Enforces commit message formatting)*:
   ```bash
   make hooks
   ```

---

## Common Makefile Commands

| Command | Description |
| --- | --- |
| `make backend` | Build the Python backend wheel and fetch `uv` binaries for embedding |
| `make run` | Build and run the TUI locally |
| `make proto` | Generate Protobuf & Connect RPC code for both Go and Python |
| `make hooks` | Set up local git commit hooks |

---

## Development Workflow

### Modifying Protobuf Schemas
If you update `proto/music.proto`:
```bash
make proto
```

### Running Tests
Run the Go unit test suite:
```bash
go test -v ./...
```

---

## Pull Request Guidelines

1. **Conventional Commits**: Commit messages must follow Conventional Commit rules (e.g., `feat: ...`, `fix: ...`, `docs: ...`, `test: ...`).
2. **Code Formatting**: Ensure Go code is formatted (`go fmt ./...`) and tests pass (`go test ./...`).
3. **Keep PRs Focused**: Limit PRs to a single feature or bugfix for fast review.

---

## Reporting Issues

Use [GitHub Issues](https://github.com/kumneger0/yt-tracks/issues) for bug reports and feature requests. Include your OS, Go version, and log output (`yt-tracks log` or `~/.config/yt-tracks/yt-tracks.log`) when reporting bugs.
