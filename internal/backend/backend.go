package backend

import (
	"bytes"
	"os/exec"
	"sync"
)

type Backend struct {
	cmd       *exec.Cmd
	url       string
	stdoutBuf bytes.Buffer
	stderrBuf bytes.Buffer
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
	if b.cmd.ProcessState != nil && b.cmd.ProcessState.Exited() {
		return false
	}
	return isProcessAlive(b.cmd.Process)
}

func (b *Backend) Stderr() string {
	return b.stderrBuf.String()
}
