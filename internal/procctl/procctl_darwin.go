//go:build darwin

package procctl

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func init() { active = darwinProbe{} }

// darwinProbe reads process identity via `ps` and boot identity via `sysctl`
// (no /proc on macOS, no third-party deps).
type darwinProbe struct{}

func (darwinProbe) supportsDurableKill() bool { return true }
func (darwinProbe) alive(pid int) bool        { return unixAlive(pid) }
func (darwinProbe) pgid(pid int) (int, bool)  { return unixPgid(pid) }

// bootID is the kernel boot time, stable across the boot.
func (darwinProbe) bootID() string {
	out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// psField runs `ps -p <pid> -o <fmt>=` and returns the trimmed single line.
func psField(pid int, format string) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", format+"=").Output()
	if err != nil {
		return "", false
	}
	v := strings.TrimSpace(string(out))
	if v == "" {
		return "", false
	}
	return v, true
}

// procStart is the process's start time at MICROSECOND resolution via the
// kern.proc.pid kinfo_proc (p_starttime), so a reused PID cannot share the
// recorded start identity. Capture and verify call this identically.
func (darwinProbe) procStart(pid int) (string, bool) {
	kp, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil || kp == nil {
		return "", false
	}
	st := kp.Proc.P_starttime
	return fmt.Sprintf("%d.%06d", int64(st.Sec), int64(st.Usec)), true
}

func (darwinProbe) command(pid int) (string, bool) { return psField(pid, "command") }
