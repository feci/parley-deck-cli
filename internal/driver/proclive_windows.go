//go:build windows

package driver

// processAlive on Windows conservatively assumes any existing lock is held by a
// live process. There is no portable signal-0 liveness probe on Windows without
// golang.org/x/sys, and os.FindProcess always succeeds, so the Unix probe would
// always report "dead" and defeat the lock entirely (AF3). Returning true instead
// errs toward refusing to start a second driver (safe) rather than risking two
// concurrent drivers corrupting the workspace. A genuinely stale lock left by a
// crashed driver must be removed manually (the lock-acquire error names the path).
func processAlive(pid int) bool { return true }
