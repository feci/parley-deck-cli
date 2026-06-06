//go:build linux

package procctl

import (
	"os"
	"strconv"
	"strings"
)

func init() { active = linuxProbe{} }

// linuxProbe reads process identity from /proc.
type linuxProbe struct{}

func (linuxProbe) supportsDurableKill() bool { return true }
func (linuxProbe) alive(pid int) bool        { return unixAlive(pid) }
func (linuxProbe) pgid(pid int) (int, bool)  { return unixPgid(pid) }

func (linuxProbe) bootID() string {
	data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// procStart is field 22 (starttime, in clock ticks since boot) of /proc/<pid>/stat.
// The comm field (2) may contain spaces/parens, so parse after the final ')'.
func (linuxProbe) procStart(pid int) (string, bool) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return "", false
	}
	s := string(data)
	close := strings.LastIndexByte(s, ')')
	if close < 0 {
		return "", false
	}
	fields := strings.Fields(s[close+1:])
	// After ')': index 0 = state (field 3); starttime (field 22) is index 19.
	if len(fields) < 20 {
		return "", false
	}
	return fields[19], true
}

func (linuxProbe) command(pid int) (string, bool) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/cmdline")
	if err != nil || len(data) == 0 {
		return "", false
	}
	// cmdline is NUL-separated argv.
	return strings.TrimSpace(strings.ReplaceAll(string(data), "\x00", " ")), true
}
