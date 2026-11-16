package backend

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var embeddedBackend []byte

func getOrExtractBackend() (string, error) {
	if len(embeddedBackend) == 0 {
		return "", fmt.Errorf("embedded backend binary is empty")
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	binDir := filepath.Join(cacheDir, "yt-tracks", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	hash := sha256.Sum256(embeddedBackend)
	exeName := fmt.Sprintf("yt-tracks-backend-%x", hash[:8])
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	targetPath := filepath.Join(binDir, exeName)

	if info, err := os.Stat(targetPath); err == nil && info.Size() == int64(len(embeddedBackend)) {
		return targetPath, nil
	}

	tmpPath := targetPath + ".tmp"
	if err := os.WriteFile(tmpPath, embeddedBackend, 0755); err != nil {
		return "", fmt.Errorf("failed to write cached backend: %w", err)
	}
	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("failed to finalize cached backend: %w", err)
	}

	return targetPath, nil
}

func findLocalPython() string {
	candidates := []string{
		filepath.Join(".venv", "bin", "python"),
		filepath.Join("venv", "bin", "python"),
		filepath.Join(".venv", "Scripts", "python.exe"),
		filepath.Join("venv", "Scripts", "python.exe"),
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return ""
}

func StartEmbedded(ctx context.Context) (*Backend, error) {
	if serverURL := os.Getenv("YT_TRACKS_SERVER_URL"); serverURL != "" {
		slog.Info("using external backend server", "url", serverURL)
		return &Backend{url: serverURL}, nil
	}

	if len(embeddedBackend) == 0 {
		if localPython := findLocalPython(); localPython != "" {
			slog.Info("embedded backend not present; starting local python environment", "path", localPython)
			return Start(ctx, localPython, "-m", "grpc_server.main")
		}

		return nil, fmt.Errorf("embedded backend is missing from this build: run 'make backend' before building, or run 'make server-run' for local development")
	}

	path, err := getOrExtractBackend()
	if err != nil {
		return nil, err
	}

	return Start(ctx, path)
}
