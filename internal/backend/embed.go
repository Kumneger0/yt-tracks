package backend

import (
	"context"
	"crypto/sha256"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gofrs/flock"
)

const (
	pythonVersion         = "3.13"
	defaultStartupTimeout = 10 * time.Minute
)

var errMissingAssets = errors.New("embedded backend assets are missing: run 'make backend' before building")

func readSingleAsset(fsys embed.FS, dir, suffix string) (string, []byte, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return "", nil, errMissingAssets
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || strings.HasPrefix(name, ".") || (suffix != "" && !strings.HasSuffix(name, suffix)) {
			continue
		}
		data, readErr := fs.ReadFile(fsys, path.Join(dir, name))
		if readErr != nil {
			return "", nil, fmt.Errorf("failed to read embedded %s: %w", name, readErr)
		}
		return name, data, nil
	}
	return "", nil, errMissingAssets
}

func cacheRootDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user cache directory: %w", err)
	}
	root := filepath.Join(cacheDir, "yt-tracks")
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}
	return root, nil
}

func writeCachedFile(targetPath string, data []byte, perm os.FileMode) error {
	if info, err := os.Stat(targetPath); err == nil && info.Size() == int64(len(data)) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
	}

	tmpPath := targetPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetPath, err)
	}
	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to finalize %s: %w", targetPath, err)
	}
	return nil
}

func venvPaths(venvDir string) (string, string) {
	if runtime.GOOS == "windows" {
		return filepath.Join(venvDir, "Scripts", "python.exe"),
			filepath.Join(venvDir, "Scripts", "yt-tracks-backend.exe")
	}
	return filepath.Join(venvDir, "bin", "python"),
		filepath.Join(venvDir, "bin", "yt-tracks-backend")
}

func isBackendReady(readyMarker, entrypoint string) bool {
	if _, err := os.Stat(readyMarker); err != nil {
		return false
	}
	_, err := os.Stat(entrypoint)
	return err == nil
}

func runUV(ctx context.Context, uvPath string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, uvPath, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func acquireBackendLock(ctx context.Context, cacheRoot string) (*flock.Flock, error) {
	fileLock := flock.New(filepath.Join(cacheRoot, "backend.lock"))
	locked, err := fileLock.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire backend setup lock: %w", err)
	}
	if !locked {
		return nil, errors.New("failed to acquire backend setup lock")
	}
	return fileLock, nil
}

func ensureBackendReady(ctx context.Context, stdout, stderr io.Writer, notify bool) (string, error) {
	uvAssetName, uvBytes, err := readSingleAsset(embeddedUVFS, fmt.Sprintf("binaries/%s_%s", runtime.GOOS, runtime.GOARCH), "")
	if err != nil {
		return "", err
	}
	wheelName, wheelBytes, err := readSingleAsset(embeddedWheelFS, "binaries/wheel", ".whl")
	if err != nil {
		return "", err
	}

	root, err := cacheRootDir()
	if err != nil {
		return "", err
	}

	wheelHash := sha256.Sum256(wheelBytes)
	venvDir := filepath.Join(root, "venvs", fmt.Sprintf("%x", wheelHash[:8]))
	readyMarker := filepath.Join(venvDir, ".ready")
	venvPython, entrypoint := venvPaths(venvDir)

	if isBackendReady(readyMarker, entrypoint) {
		return entrypoint, nil
	}

	fileLock, err := acquireBackendLock(ctx, root)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fileLock.Unlock()
	}()

	if isBackendReady(readyMarker, entrypoint) {
		return entrypoint, nil
	}

	uvHash := sha256.Sum256(uvBytes)
	uvPath := filepath.Join(root, "bin", fmt.Sprintf("%x-%s", uvHash[:8], uvAssetName))
	if writeErr := writeCachedFile(uvPath, uvBytes, 0755); writeErr != nil {
		return "", writeErr
	}

	tmpDir, err := os.MkdirTemp("", "yt-tracks-wheel-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp wheel directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	wheelPath := filepath.Join(tmpDir, wheelName)
	if writeErr := os.WriteFile(wheelPath, wheelBytes, 0644); writeErr != nil {
		return "", fmt.Errorf("failed to write temp wheel: %w", writeErr)
	}

	if notify {
		fmt.Fprintln(os.Stderr, "==> Initializing backend runtime (first run may take a minute to download CPython and dependencies)...")
		fmt.Fprintln(os.Stderr, "    Tip: Run 'yt-tracks setup' to see real-time download progress.")
	}

	if err := runUV(ctx, uvPath, []string{"venv", "--clear", "--python", pythonVersion, venvDir}, stdout, stderr); err != nil {
		return "", fmt.Errorf("failed to create Python %s venv: %w", pythonVersion, err)
	}
	if err := runUV(ctx, uvPath, []string{"pip", "install", "--python", venvPython, wheelPath}, stdout, stderr); err != nil {
		return "", fmt.Errorf("failed to install backend wheel: %w", err)
	}
	if err := os.WriteFile(readyMarker, []byte("ok\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to write venv ready marker: %w", err)
	}

	return entrypoint, nil
}

func withDefaultStartupTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultStartupTimeout)
}

func StartEmbedded(ctx context.Context) (*Backend, error) {
	ctx, cancel := withDefaultStartupTimeout(ctx)
	defer cancel()

	var errBuf tailBuffer
	entrypoint, err := ensureBackendReady(ctx, io.Discard, &errBuf, true)
	if err != nil {
		if stderrOut := errBuf.String(); stderrOut != "" {
			return nil, fmt.Errorf("%w\nuv output:\n%s", err, stderrOut)
		}
		return nil, err
	}

	return Start(ctx, entrypoint)
}

func SetupEmbedded(ctx context.Context) error {
	fmt.Println("==> Setting up yt-tracks backend runtime...")
	fmt.Println("    Downloading CPython and dependencies via uv (if not cached)...")

	entrypoint, err := ensureBackendReady(ctx, os.Stdout, os.Stderr, false)
	if err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}

	backend, err := StartWithOutput(ctx, entrypoint, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("setup verification failed: %w", err)
	}
	defer backend.Stop()

	fmt.Println("\n✓ Setup completed successfully! You can now run 'yt-tracks'.")
	return nil
}
