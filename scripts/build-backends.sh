#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

UV_VERSION="${UV_VERSION:-0.11.6}"
UV_CACHE_DIR="$ROOT/build/uv-cache/$UV_VERSION"
BIN_DIR="internal/backend/binaries"
mkdir -p "$UV_CACHE_DIR" "$BIN_DIR/wheel"
rm -f internal/backend/embedded_*.go

echo "--> Generating Python protobuf / Connect modules"
buf generate --template buf.gen.python.yaml

echo "--> Building Python wheel"
rm -f "$BIN_DIR"/wheel/*.whl "$BIN_DIR/wheel/.gitignore"
tmp_wheel="$(mktemp -d)"
uv build --wheel --out-dir "$tmp_wheel"
cp "$tmp_wheel"/yt_tracks-*.whl "$BIN_DIR/wheel/"
rm -rf "$tmp_wheel"

echo "--> Preparing uv $UV_VERSION binaries"
# Format: <goos>_<goarch> <uv_target> <exe_name> <archive_ext>
TARGETS=(
    "linux_amd64 x86_64-unknown-linux-musl uv tar.gz"
    "darwin_arm64 aarch64-apple-darwin uv tar.gz"
    "windows_amd64 x86_64-pc-windows-msvc uv.exe zip"
)

for target in "${TARGETS[@]}"; do
    read -r platform uv_target exe format <<< "$target"
    archive_name="uv-${uv_target}.${format}"
    archive_path="$UV_CACHE_DIR/$archive_name"
    url="https://github.com/astral-sh/uv/releases/download/${UV_VERSION}/${archive_name}"

    if [ ! -f "$archive_path" ]; then
        echo "   Downloading $url"
        curl -fsSL -o "${archive_path}.tmp" "$url"
        mv "${archive_path}.tmp" "$archive_path"
    fi

    target_dir="$BIN_DIR/$platform"
    mkdir -p "$target_dir"
    if [ "$format" = "zip" ]; then
        unzip -jqo "$archive_path" "$exe" -d "$target_dir"
    else
        tar -xzf "$archive_path" -C "$target_dir" --strip-components=1 "uv-${uv_target}/$exe"
    fi
    chmod +x "$target_dir/$exe"
    echo "   $platform -> $target_dir/$exe"
done
