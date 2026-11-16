package backend

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var embeddedBackend []byte

func getCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = filepath.Join(os.TempDir(), fmt.Sprintf("yt-tracks-%d", os.Getuid()))
	}
	binDir := filepath.Join(cacheDir, "yt-tracks", "bin")
	if err := os.MkdirAll(binDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create private cache directory: %w", err)
	}
	return binDir, nil
}

func getOrExtractBackend() (string, error) {
	if len(embeddedBackend) == 0 {
		return "", fmt.Errorf("embedded backend binary is missing: run 'make backend' before building")
	}

	binDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(embeddedBackend)
	exeName := fmt.Sprintf("yt-tracks-backend-%x", hash[:8])
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	targetPath := filepath.Join(binDir, exeName)

	if info, err := os.Stat(targetPath); err == nil && info.Size() == int64(len(embeddedBackend)) {
		if file, err := os.Open(targetPath); err == nil {
			h := sha256.New()
			if _, err := io.Copy(h, file); err == nil && subtle.ConstantTimeCompare(h.Sum(nil), hash[:]) == 1 {
				_ = file.Close()
				return targetPath, nil
			}
			_ = file.Close()
		}
	}

	tmpFile, err := os.CreateTemp(binDir, exeName+"-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary extraction file: %w", err)
	}
	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(embeddedBackend); err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("failed to write embedded backend: %w", err)
	}
	if err := tmpFile.Chmod(0755); err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("failed to set backend permissions: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		return "", fmt.Errorf("failed to finalize cached backend: %w", err)
	}
	cleanup = false

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

	backend, err := StartWithOutput(ctx, path, os.Stdout, os.Stderr, true)
	if err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}
	defer backend.Stop()

	fmt.Println("\n✓ Setup completed successfully! You can now run 'yt-tracks'.")
	return nil
}
