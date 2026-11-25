package backend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPickFreePort(t *testing.T) {
	port, err := pickFreePort()
	if err != nil {
		t.Fatalf("pickFreePort() failed: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Fatalf("invalid port: %d", port)
	}
}

func TestTailBufferBoundsAndConcurrency(t *testing.T) {
	buf := tailBuffer{limit: 8}
	_, _ = buf.Write([]byte("12345"))
	_, _ = buf.Write([]byte("67890"))
	if got := buf.String(); got != "34567890" {
		t.Fatalf("expected bounded tail %q, got %q", "34567890", got)
	}

	_, _ = buf.Write([]byte("abcdefghijkl"))
	if got := buf.String(); got != "efghijkl" {
		t.Fatalf("expected truncated single write %q, got %q", "efghijkl", got)
	}

	backend := &Backend{stderrBuf: tailBuffer{limit: 64}}
	var waitGroup sync.WaitGroup
	for worker := range 8 {
		waitGroup.Go(func() {
			for iter := range 50 {
				if (worker+iter)%2 == 0 {
					_, _ = backend.stderrBuf.Write([]byte(strings.Repeat("x", 16)))
				} else {
					_ = backend.Stderr()
				}
			}
		})
	}
	waitGroup.Wait()
}

func TestWaitForReadyCancelsStuckHealthCheck(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	var startedOnce sync.Once

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/yt_tracks.music.v1.MusicService/HealthCheck" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		writer.WriteHeader(http.StatusOK)
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
		startedOnce.Do(func() { close(started) })
		select {
		case <-release:
		case <-req.Context().Done():
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-started
		cancel()
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- waitForReady(ctx, server.URL, make(chan error))
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected context cancellation error, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("waitForReady blocked past context cancellation")
	}
}
