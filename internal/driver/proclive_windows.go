//go:build windows

package driver

import (
	"errors"

	"golang.org/x/sys/windows"
)

// processAlive is the truthful P-A probe (FINAL §D.2) — golang.org/x/sys is a
// direct dependency (go.mod), so the old "no portable probe without x/sys"
// placeholder is retired. OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION) +
// GetExitCodeProcess: a PID that cannot be opened with ERROR_INVALID_PARAMETER
// does not exist (dead); any other observation failure — ERROR_ACCESS_DENIED
// first among them — is unverifiable, and unverifiable is ALIVE for refusal
// purposes (fail closed: a stale-looking lock must never be treated as
// releasable while a live owner might hold it). STILL_ACTIVE (259) is the
// documented sentinel caveat: a genuine exit code 259 is indistinguishable
// and reported alive.
const stillActive = 259 // STILL_ACTIVE sentinel (§D.2 documented caveat)

func processAlive(pid int) bool {
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
		return true
	}
	return code == stillActive
}
