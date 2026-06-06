//go:build windows

package procctl

import (
	"os/exec"
	"strconv"
)

func init() { active = windowsProbe{} }

// windowsProbe: durable cross-restart kill is not supported (no dependency-free
// way to attribute a PID safely); the live in-process handle still kills via the
// process tree below. Windows is not the owner's platform.
type windowsProbe struct{}

func (windowsProbe) supportsDurableKill() bool      { return false }
func (windowsProbe) bootID() string                 { return "" }
func (windowsProbe) alive(pid int) bool             { return pid > 0 }
func (windowsProbe) procStart(pid int) (string, bool) { return "", false }
func (windowsProbe) command(pid int) (string, bool)   { return "", false }
func (windowsProbe) pgid(pid int) (int, bool)         { return 0, false }

// SetNewProcessGroup is a no-op placeholder on Windows.
func SetNewProcessGroup(_ *exec.Cmd) {}

// KillGroup best-effort kills the process tree via taskkill (live path only).
func KillGroup(s Spawned) error {
	if s.PID <= 0 {
		return nil
	}
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(s.PID)).Run()
}
