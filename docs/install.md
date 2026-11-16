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

yt-tracks embeds its YouTube Music Connect RPC backend (`yt-tracks-backend`)
inside the Go binary via `//go:embed`, built with
[uvbox](https://github.com/AmadeusITGroup/uvbox). The backend runs on `127.0.0.1`
only and requires no user-installed Python — on first launch it downloads a
CPython runtime and installs the backend dependencies (one-time network
required).

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/kumneger0/yt-tracks.git
cd yt-tracks
```

### 2. Build the Python Backend

Install `uvbox` and build the backend launcher to embed it into the Go binary:

```bash
uv tool install uvbox
make backend
```

This generates `internal/backend/embedded_*.go` sources with `//go:embed`
directives pointing at the platform-specific launcher binaries.

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
  (uvbox downloads CPython on first launch).
- Open an issue on [GitHub Issues](https://github.com/Kumneger0/yt-tracks/issues).
