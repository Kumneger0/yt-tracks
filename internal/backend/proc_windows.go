//go:build windows

package backend

import (
	"os"
	"os/exec"
)

func setProcessGroup(cmd *exec.Cmd) {
	// No Setpgid on Windows
}

func terminateProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}

func isProcessAlive(process *os.Process) bool {
	if process == nil {
		return false
	}
	p, err := os.FindProcess(process.Pid)
	return err == nil && p != nil
}
