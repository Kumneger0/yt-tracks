# Building from Source

## Prerequisites

Before building `yt-tracks` from source, ensure you have the following installed:

- **[Go](https://go.dev/dl/)** (v1.26 or higher)
- **[uv](https://docs.astral.sh/uv/)** (Python package manager, for building the backend)
- **[buf](https://buf.build/)** (protobuf code generation, for building the backend)

- **System Dependencies**:
  - **[yt-dlp](https://github.com/yt-dlp/yt-dlp)**: Required for extracting and fetching YouTube Music audio stream URLs
  - **[ffmpeg](https://ffmpeg.org/)**: Required for audio stream decoding and playback

## Backend

yt-tracks embeds its YouTube Music Connect RPC backend wheel (`yt-tracks-backend`)
and the official [uv](https://docs.astral.sh/uv/) binary inside the Go binary via
`//go:embed`. The backend runs on `127.0.0.1` only and requires no
user-installed Python — on first launch it downloads a CPython runtime and
installs the backend dependencies into your cache directory (one-time network
required).

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/kumneger0/yt-tracks.git
cd yt-tracks
```

### 2. Build the Python Backend

Build the Python wheel and fetch the `uv` binaries to embed into the Go binary:

```bash
make backend
```

This generates the Python protobuf stubs, builds `yt_tracks-*.whl` into
`internal/backend/binaries/wheel/`, and downloads platform-specific `uv`
binaries into `internal/backend/binaries/<goos>_<goarch>/` for Go's static
`//go:embed` files.

### 3. Build the Go Application

Build the application binary:

```bash
make build
```

> **Note:** Run `make backend` before `make build`. Without the embedded backend,
> there is no Connect RPC server to connect to.

> **Note for Linux Users:**
> Ensure `gcc` is installed on your system to support CGO audio bindings:
> ```bash
> CGO_ENABLED=1 make build
> ```

### 4. Install & Launch

To install `yt-tracks` to `/usr/local/bin`:

```bash
make install
```

Or run the built binary directly:

```bash
./yt-tracks
```

Optionally, pre-download and initialize the local Python backend runtime with live progress before launching the TUI:

```bash
./yt-tracks setup
```

### 5. Authenticate with YouTube Music

To access your playlists, liked songs, and library, extract session cookies from your browser:

```bash
yt-tracks extract-cookies
```

---

## Support & Troubleshooting

If you encounter issues during installation or build:
- Ensure `yt-dlp` and `ffmpeg` are in your system `PATH`.
- Check application logs via `yt-tracks log`.
- For local backend first-run failures, ensure network access is available
  (`uv` downloads CPython and Python dependencies on first launch; run `yt-tracks setup` to see live progress).
- Open an issue on [GitHub Issues](https://github.com/Kumneger0/yt-tracks/issues).
