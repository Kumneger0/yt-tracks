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

const healthCheckTimeout = 2 * time.Second

func Start(ctx context.Context, path string) (*Backend, error) {
	svcBackend := &Backend{}
	return startWithOutput(ctx, path, svcBackend, io.Discard, &svcBackend.stderrBuf)
}

func StartWithOutput(ctx context.Context, path string, stdout, stderr io.Writer) (*Backend, error) {
	svcBackend := &Backend{}
	return startWithOutput(ctx, path, svcBackend, stdout, stderr)
}

func startWithOutput(ctx context.Context, path string, svcBackend *Backend, stdout, stderr io.Writer) (*Backend, error) {
	port, err := pickFreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to pick a free port: %w", err)
	}

	cmd := exec.Command(path)
	cmd.Env = append(os.Environ(),
		"HOST=127.0.0.1",
		fmt.Sprintf("PORT=%d", port),
	)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start backend %q: %w\nstderr: %s", path, err, svcBackend.Stderr())
	}

	svcBackend.cmd = cmd
	svcBackend.done = make(chan struct{})
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	exitChan := make(chan error, 1)
	go func() {
		exitChan <- cmd.Wait()
		close(svcBackend.done)
	}()

	if err := waitForReady(ctx, url, exitChan); err != nil {
		svcBackend.stopNow()
		return nil, fmt.Errorf("backend at %s did not become ready: %w\nstderr: %s", url, err, svcBackend.Stderr())
	}

	svcBackend.url = url
	return svcBackend, nil
}

func (b *Backend) stopNow() {
	if b.cmd == nil || b.cmd.Process == nil {
		return
	}

	_ = terminateProcess(b.cmd)

	if b.done == nil {
		_ = b.cmd.Process.Kill()
		return
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-b.done:
	case <-timer.C:
		_ = b.cmd.Process.Kill()
		<-b.done
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
		deadline = time.Now().Add(10 * time.Minute)
	}

	endpoint := url + "/yt_tracks.music.v1.MusicService/HealthCheck"
	client := &http.Client{Timeout: healthCheckTimeout}

	for time.Now().Before(deadline) {
		select {
		case err := <-exitChan:
			return fmt.Errorf("backend process exited prematurely: %v", err)
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if checkBackendHealth(ctx, client, endpoint) {
			return nil
		}

		timer := time.NewTimer(250 * time.Millisecond)
		select {
		case err := <-exitChan:
			timer.Stop()
			return fmt.Errorf("backend process exited prematurely: %v", err)
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return fmt.Errorf("timeout waiting for backend at %s", endpoint)
}

func checkBackendHealth(ctx context.Context, client *http.Client, endpoint string) bool {
	reqCtx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	// Re-create reader on every iteration so request body is never drained
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, strings.NewReader("{}"))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}
