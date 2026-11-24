package backend

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var embeddedBackend []byte

func getOrExtractBackend() (string, error) {
	if len(embeddedBackend) == 0 {
		return "", fmt.Errorf("embedded backend binary is missing: run 'make backend' before building")
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

func StartEmbedded(ctx context.Context) (*Backend, error) {
	path, err := getOrExtractBackend()
	if err != nil {
		return nil, err
	}

	return Start(ctx, path)
}

func SetupEmbedded(ctx context.Context) error {
	path, err := getOrExtractBackend()
	if err != nil {
		return err
	}

	fmt.Println("==> Setting up yt-tracks backend runtime...")
	fmt.Println("    Downloading CPython and dependencies via uvbox (if not cached)...")

	backend, err := StartWithOutput(ctx, path, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}
	defer backend.Stop()

	fmt.Println("\n✓ Setup completed successfully! You can now run 'yt-tracks'.")
	return nil
}
