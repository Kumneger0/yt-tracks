//go:build windows

package backend

import (
	"os"
	"os/exec"
	"strconv"
)

func setProcessGroup(cmd *exec.Cmd) {
	// No Setpgid on Windows
}

func terminateProcess(cmd *exec.Cmd) error {
	return killProcessTree(cmd)
}

func killProcess(cmd *exec.Cmd) error {
	return killProcessTree(cmd)
}

func killProcessTree(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	// taskkill /T kills child processes (process tree), /F forces termination
	killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	if err := killCmd.Run(); err == nil {
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
