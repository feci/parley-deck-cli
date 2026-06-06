//go:build windows

package acp

import (
	"os/exec"
	"strconv"
)

// setSysProcAttr is a no-op on Windows; the /dev/tty issue does not apply.
func setSysProcAttr(_ *exec.Cmd) {}

// killProcessGroup best-effort kills the process tree via taskkill on Windows.
func killProcessGroup(pid int) error {
	if pid <= 0 {
		return nil
	}
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
}
