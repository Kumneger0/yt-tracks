package backend

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Start(ctx context.Context, path string, args ...string) (*Backend, error) {
	port, err := pickFreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to pick a free port: %w", err)
	}

	svcBackend := &Backend{}
	cmd := exec.Command(path, args...)
	cmd.Env = append(os.Environ(),
		"HOST=127.0.0.1",
		fmt.Sprintf("PORT=%d", port),
	)
	// Capture the child's output instead of attaching to the terminal so the
	// uvbox spinner / uvicorn access logs don't corrupt the Bubbletea TUI.
	cmd.Stdout = &svcBackend.stdoutBuf
	cmd.Stderr = &svcBackend.stderrBuf

	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start backend %q: %w\nstderr: %s", path, err, svcBackend.stderrBuf.String())
	}

	svcBackend.cmd = cmd
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	exitChan := make(chan error, 1)
	go func() {
		exitChan <- cmd.Wait()
	}()

	if err := waitForReady(ctx, url, exitChan); err != nil {
		svcBackend.stopNow()
		return nil, fmt.Errorf("backend at %s did not become ready: %w\nstderr: %s", url, err, svcBackend.stderrBuf.String())
	}

	svcBackend.url = url
	return svcBackend, nil
}

func (b *Backend) stopNow() {
	if b.cmd == nil || b.cmd.Process == nil {
		return
	}

	_ = terminateProcess(b.cmd)

	done := make(chan struct{})
	go func() {
		_ = b.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = b.cmd.Process.Kill()
	}
}

func pickFreePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}

func waitForReady(ctx context.Context, url string, exitChan <-chan error) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(3 * time.Minute)
	}

	endpoint := url + "/yt_tracks.music.v1.MusicService/HealthCheck"

	for time.Now().Before(deadline) {
		select {
		case err := <-exitChan:
			return fmt.Errorf("backend process exited prematurely: %v", err)
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Re-create reader on every iteration so request body is never drained
		resp, err := http.Post(endpoint, "application/json", strings.NewReader("{}"))
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(250 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for backend at %s", endpoint)
}
