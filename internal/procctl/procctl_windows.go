//go:build windows

package procctl

import (
	"errors"
	"os/exec"
	"strconv"

	"golang.org/x/sys/windows"
)

func init() { active = windowsProbe{} }

// windowsProbe: durable cross-restart kill is not supported on Windows (P-B
// follow-up idea: Job Objects with JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
// GetProcessTimes attribution, QueryFullProcessImageName); the live
// in-process handle still kills via the process tree below.
type windowsProbe struct{}

func (windowsProbe) supportsDurableKill() bool        { return false }
func (windowsProbe) bootID() string                   { return "" }
func (windowsProbe) procStart(pid int) (string, bool) { return "", false }
func (windowsProbe) command(pid int) (string, bool)   { return "", false }
func (windowsProbe) pgid(pid int) (int, bool)         { return 0, false }

// alive is the truthful P-A probe (FINAL §D.2): OpenProcess with
// PROCESS_QUERY_LIMITED_INFORMATION + GetExitCodeProcess. Fail-closed shape:
// ERROR_ACCESS_DENIED (or any other observation failure) means unverifiable,
// and unverifiable is ALIVE for refusal purposes — a PID we cannot probe is
// never treated as dead. ERROR_INVALID_PARAMETER means the PID does not exist
// (provably dead). Caveat, documented per §D.2: STILL_ACTIVE (259) is a
// sentinel, so a process that genuinely exited with code 259 is
// indistinguishable from a live one and is reported alive (fail closed).
const stillActive = 259 // STILL_ACTIVE sentinel (§D.2 documented caveat)

func (windowsProbe) alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return !errors.Is(err, windows.ERROR_INVALID_PARAMETER)
	}
	defer windows.CloseHandle(h)
	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return true // unobservable: fail closed, treat as alive
	}
	return code == stillActive
}

// SetNewProcessGroup is a no-op placeholder on Windows.
func SetNewProcessGroup(_ *exec.Cmd) {}

// KillGroup best-effort kills the process tree via taskkill (live path only).
func KillGroup(s Spawned) error {
	if s.PID <= 0 {
		return nil
	}
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(s.PID)).Run()
}
