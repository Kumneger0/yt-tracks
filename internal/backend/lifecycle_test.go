package backend

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const fakeBackendScript = `#!/bin/sh
port="${PORT:-8080}"
host="${HOST:-0.0.0.0}"
( python3 - "$port" "$host" <<'PY'
import http.server, json, sys, threading
port = int(sys.argv[1])
host = sys.argv[2]
class H(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == "/yt_tracks.music.v1.MusicService/HealthCheck":
            body = b'{"ok": true}'
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
        else:
            self.send_response(404)
            self.end_headers()
    def log_message(self, fmt, *args): pass
http.server.HTTPServer((host, port), H).serve_forever()
PY
) &
wait
`

func writeFakeBackend(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "fake-backend")
	if err := os.WriteFile(path, []byte(fakeBackendScript), 0755); err != nil {
		t.Fatalf("writing fake backend: %v", err)
	}
	return path
}

func TestStartAndStop(t *testing.T) {
	if os.Getenv("YT_TRACKS_BACKEND_TEST_SKIP") == "1" {
		t.Skip("YT_TRACKS_BACKEND_TEST_SKIP set")
	}

	dir := t.TempDir()
	backendPath := writeFakeBackend(t, dir)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	backend, err := Start(ctx, backendPath)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if backend.URL() == "" {
		t.Fatal("expected non-empty URL")
	}
	t.Logf("backend URL: %s", backend.URL())

	resp, err := http.Post(backend.URL()+"/yt_tracks.music.v1.MusicService/HealthCheck",
		"application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("HealthCheck request failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HealthCheck returned %d: %s", resp.StatusCode, body)
	}
	if !strings.Contains(string(body), `"ok": true`) {
		t.Fatalf("unexpected HealthCheck body: %s", body)
	}

	backend.Stop()
	time.Sleep(200 * time.Millisecond)
	if backend.IsRunning() {
		t.Error("backend process still running after Stop()")
	}
}

func TestPickFreePort(t *testing.T) {
	port, err := pickFreePort()
	if err != nil {
		t.Fatalf("pickFreePort() failed: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Fatalf("invalid port: %d", port)
	}
}
