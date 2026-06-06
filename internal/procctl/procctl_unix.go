//go:build !windows

package procctl

import (
	"errors"
	"os/exec"
	"syscall"
	"time"
)

// SetNewProcessGroup makes the child a session leader (Setsid), so pgid == pid
// and killing -pgid reaps the whole tree (matches the ACP spawn path).
func SetNewProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true
}

// KillGroup terminates the recorded process group: SIGTERM, a short grace
// period, then SIGKILL. Signaling an already-dead group (ESRCH) is a no-op. This
// is ungated — callers that own the live process (the timeout watcher) use it
// directly; the reattach path goes through KillTreeAttributed.
func KillGroup(s Spawned) error {
	target := s.PGID
	if target <= 0 {
		target = s.PID
	}
	if target <= 0 {
		return errors.New("no pid/pgid to kill")
	}
	// Self-protection: never signal parley's own process group (would happen if
	// Setsid was skipped/failed and the agent inherited our group — killing the
	// group would take down parley and its session).
	if self, err := syscall.Getpgid(0); err == nil && target == self {
		return errors.New("refusing to kill parley's own process group")
	}
	// Negative pid signals the whole process group.
	_ = syscall.Kill(-target, syscall.SIGTERM)
	deadline := time.Now().Add(1500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(-target, syscall.Signal(0)); err != nil {
			return nil // group gone
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := syscall.Kill(-target, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	return nil
}

// unixAlive reports whether a pid exists (signal 0 probe).
func unixAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}

// unixPgid returns the process group of pid.
func unixPgid(pid int) (int, bool) {
	pg, err := syscall.Getpgid(pid)
	if err != nil {
		return 0, false
	}
	return pg, true
}
