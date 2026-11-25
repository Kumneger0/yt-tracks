package backend

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcquireBackendLockSerializesAndReChecksReady(t *testing.T) {
	cacheRoot := t.TempDir()
	venvDir := filepath.Join(cacheRoot, "venvs", "test")
	if err := os.MkdirAll(filepath.Dir(venvDir), 0755); err != nil {
		t.Fatalf("mkdir venvs: %v", err)
	}
	readyMarker := filepath.Join(venvDir, ".ready")
	_, entrypoint := venvPaths(venvDir)

	firstLock, err := acquireBackendLock(context.Background(), cacheRoot)
	if err != nil {
		t.Fatalf("first acquireBackendLock failed: %v", err)
	}

	acquired := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		secondLock, lockErr := acquireBackendLock(context.Background(), cacheRoot)
		if lockErr != nil {
			errCh <- lockErr
			return
		}
		close(acquired)
		_ = secondLock.Unlock()
	}()

	select {
	case <-acquired:
		t.Fatal("second caller acquired lock while first caller still held it")
	case lockErr := <-errCh:
		t.Fatalf("second acquireBackendLock returned error: %v", lockErr)
	case <-time.After(150 * time.Millisecond):
	}

	if err := os.MkdirAll(filepath.Dir(entrypoint), 0755); err != nil {
		t.Fatalf("mkdir entrypoint dir: %v", err)
	}
	if err := os.WriteFile(entrypoint, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("write entrypoint: %v", err)
	}
	if err := os.WriteFile(readyMarker, []byte("ok\n"), 0644); err != nil {
		t.Fatalf("write ready marker: %v", err)
	}
	if err := firstLock.Unlock(); err != nil {
		t.Fatalf("unlock first lock: %v", err)
	}

	select {
	case <-acquired:
	case lockErr := <-errCh:
		t.Fatalf("second acquireBackendLock failed after unlock: %v", lockErr)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second caller to acquire lock")
	}

	if !isBackendReady(readyMarker, entrypoint) {
		t.Fatal("expected isBackendReady to return true after first caller initialized runtime")
	}
}

func TestWithDefaultStartupTimeout(t *testing.T) {
	baseCtx, baseCancel := withDefaultStartupTimeout(context.Background())
	defer baseCancel()

	deadline, hasDeadline := baseCtx.Deadline()
	if !hasDeadline {
		t.Fatal("expected default startup deadline to be set")
	}
	remaining := time.Until(deadline)
	if remaining <= 9*time.Minute || remaining > defaultStartupTimeout {
		t.Fatalf("unexpected default deadline remaining: %v", remaining)
	}

	earlier := time.Now().Add(2 * time.Second)
	shortCtx, shortCancel := context.WithDeadline(context.Background(), earlier)
	defer shortCancel()

	wrappedCtx, wrappedCancel := withDefaultStartupTimeout(shortCtx)
	defer wrappedCancel()

	gotDeadline, hasWrappedDeadline := wrappedCtx.Deadline()
	if !hasWrappedDeadline {
		t.Fatal("expected deadline on wrapped context")
	}
	if !gotDeadline.Equal(earlier) {
		t.Fatalf("expected earlier deadline %v to be preserved, got %v", earlier, gotDeadline)
	}
}

func TestCacheRootDirReturnsUserCacheDirError(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("LocalAppData", "")

	if _, err := cacheRootDir(); err == nil {
		t.Fatal("expected cacheRootDir() to return an error when user cache directory cannot be resolved")
	}
}
