package backend

import (
	"os/exec"
	"sync"
)

const maxDiagnosticBytes = 64 * 1024

type tailBuffer struct {
	mutex sync.Mutex
	limit int
	data  []byte
}

func (buf *tailBuffer) Write(chunk []byte) (int, error) {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	maxBytes := buf.limit
	if maxBytes <= 0 {
		maxBytes = maxDiagnosticBytes
	}
	if len(chunk) >= maxBytes {
		buf.data = append(buf.data[:0], chunk[len(chunk)-maxBytes:]...)
		return len(chunk), nil
	}
	if overflow := len(buf.data) + len(chunk) - maxBytes; overflow > 0 {
		copy(buf.data, buf.data[overflow:])
		buf.data = buf.data[:len(buf.data)-overflow]
	}
	buf.data = append(buf.data, chunk...)
	return len(chunk), nil
}

func (buf *tailBuffer) String() string {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()
	return string(buf.data)
}

type Backend struct {
	cmd       *exec.Cmd
	url       string
	stderrBuf tailBuffer
	done      chan struct{}
	stopOnce  sync.Once
}

func (b *Backend) URL() string { return b.url }

func (b *Backend) Stop() {
	b.stopOnce.Do(func() {
		b.stopNow()
	})
}

func (b *Backend) IsRunning() bool {
	if b.cmd == nil || b.cmd.Process == nil {
		return false
	}
	if b.done != nil {
		select {
		case <-b.done:
			return false
		default:
		}
	}
	return isProcessAlive(b.cmd.Process)
}

func (b *Backend) Stderr() string {
	return b.stderrBuf.String()
}
