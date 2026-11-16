#!/bin/bash
# Build the yt-tracks-backend Python launcher with uvbox and generate Go
# //go:embed source files so the backend is compiled into the Go binary.
#
# Outputs:
#   internal/backend/binaries/<goos>_<goarch>/yt-tracks-backend(.exe)  raw launcher binaries
#   internal/backend/embedded_<goos>_<goarch>.go                      //go:embed source (build-tagged)
#
# Run this before `go build` so the embedded backend bytes are available.
# In CI this runs as a goreleaser before.hook.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "--> Syncing Python dependencies"
uv sync --frozen

echo "--> Generating Python protobuf / Connect modules"
buf generate --template buf.gen.python.yaml

echo "--> Building Python wheel"
WHEEL_DIR="$(mktemp -d)"
uv build --wheel --out-dir "$WHEEL_DIR"

WHEEL=$(ls "$WHEEL_DIR"/yt_tracks-*.whl 2>/dev/null | head -1)
if [ -z "$WHEEL" ]; then
    echo "ERROR: no wheel found in $WHEEL_DIR" >&2
    rm -rf "$WHEEL_DIR"
    exit 1
fi

echo "--> Building uvbox launchers"
rm -rf uvbox-tmp
uvbox wheel --config uvbox.toml --output uvbox-tmp --linux --amd "$WHEEL"
uvbox wheel --config uvbox.toml --output uvbox-tmp --darwin --arm "$WHEEL"
uvbox wheel --config uvbox.toml --output uvbox-tmp --windows --amd "$WHEEL"
rm -rf "$WHEEL_DIR"

BIN_DIR="internal/backend/binaries"
EMBED_DIR="internal/backend"
rm -rf "$BIN_DIR"
rm -f "$EMBED_DIR"/embedded_*.go
mkdir -p "$BIN_DIR"

# Format: <goos> <goarch> <uvbox_arch> <exe_name> <archive_ext>
TARGETS=(
    "linux amd64 x86_64-unknown-linux-gnu yt-tracks-backend tar.gz"
    "darwin arm64 aarch64-apple-darwin yt-tracks-backend tar.gz"
    "windows amd64 x86_64-pc-windows-msvc yt-tracks-backend.exe zip"
)

for target in "${TARGETS[@]}"; do
    set -- $target
    goos="$1"
    goarch="$2"
    uvbox_arch="$3"
    exe="$4"
    format="$5"

    archive="uvbox-tmp/yt-tracks-backend-${uvbox_arch}.${format}"
    target_dir="$BIN_DIR/${goos}_${goarch}"
    mkdir -p "$target_dir"

    tmpdir="$(mktemp -d)"
    if [ "$format" = "zip" ]; then
        unzip -q -o "$archive" -d "$tmpdir"
    else
        tar xzf "$archive" -C "$tmpdir"
    fi
    cp "$tmpdir/$exe" "$target_dir/$exe"
    chmod +x "$target_dir/$exe"
    rm -rf "$tmpdir"

    # Generate //go:embed source file (build-tagged per platform)
    embed_path="binaries/${goos}_${goarch}/${exe}"
    var_name="${goos}${goarch}Backend"
    embed_file="$EMBED_DIR/embedded_${goos}_${goarch}.go"
    cat > "$embed_file" <<EOF
//go:build $goos && $goarch

package backend

import _ "embed"

//go:embed $embed_path
var $var_name []byte

func init() {
	embeddedBackend = $var_name
}
EOF

    echo "   $goos/$goarch -> $target_dir/$exe + $embed_file"
done

rm -rf uvbox-tmp
echo "--> Embedded backend sources ready"
