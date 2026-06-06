//go:build !windows

package driver

import (
	"errors"
	"os"
	"syscall"
)

// processAlive reports whether the given PID belongs to a live process. On Unix
// it probes with signal 0; EPERM means the process exists but is owned by another
// user (still alive), so it must NOT be treated as dead (AF3).
func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}
