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

func Start(ctx context.Context, path string) (*Backend, error) {
	svcBackend := &Backend{}
	return startWithOutput(ctx, path, svcBackend, &svcBackend.stdoutBuf, &svcBackend.stderrBuf, false)
}

func StartWithOutput(ctx context.Context, path string, stdout, stderr io.Writer, isSetup bool) (*Backend, error) {
	svcBackend := &Backend{}
	return startWithOutput(ctx, path, svcBackend, stdout, stderr, isSetup)
}

func startWithOutput(ctx context.Context, path string, svcBackend *Backend, stdout, stderr io.Writer, isSetup bool) (*Backend, error) {
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
		return nil, fmt.Errorf("failed to start backend %q: %w\nstderr: %s", path, err, svcBackend.stderrBuf.String())
	}

	svcBackend.cmd = cmd
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	doneChan := make(chan struct{})
	svcBackend.doneChan = doneChan
	exitErrChan := make(chan error, 1)

	// Sole caller of cmd.Wait
	go func() {
		exitErrChan <- cmd.Wait()
		close(doneChan)
	}()

	if err := waitForReady(ctx, url, exitErrChan, isSetup); err != nil {
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

	select {
	case <-b.doneChan:
	case <-time.After(5 * time.Second):
		_ = killProcess(b.cmd)
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

func waitForReady(ctx context.Context, url string, exitChan <-chan error, isSetup bool) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Minute)
	}

	endpoint := url + "/yt_tracks.music.v1.MusicService/HealthCheck"
	start := time.Now()
	notified := false

	for time.Now().Before(deadline) {
		select {
		case err := <-exitChan:
			return fmt.Errorf("backend process exited prematurely: %v", err)
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !notified && time.Since(start) > 2*time.Second {
			fmt.Fprintln(os.Stderr, "==> Initializing backend runtime (first run may take a minute to download CPython)...")
			if !isSetup {
				fmt.Fprintln(os.Stderr, "    Tip: Run 'yt-tracks setup' to see real-time download progress.")
			}
			notified = true
		}

		// Create request with derived context and bounded per-attempt timeout
		reqCtx, reqCancel := context.WithTimeout(ctx, 2*time.Second)
		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, strings.NewReader("{}"))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					reqCancel()
					return nil
				}
			}
		}
		reqCancel()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exitChan:
			return fmt.Errorf("backend process exited prematurely: %v", err)
		case <-time.After(250 * time.Millisecond):
		}
	}

	return fmt.Errorf("timeout waiting for backend at %s", endpoint)
}
