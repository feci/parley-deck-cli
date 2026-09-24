//go:build linux

package procctl

import (
	"os"
	"strconv"
	"strings"
	"time"
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

// cmdlinePublishBound bounds how long command() waits for /proc/<pid>/cmdline
// to become non-empty on an otherwise readable process. Linux publishes argv
// late in execve: the CLOEXEC exec-status pipe that unblocks cmd.Start() is
// closed in begin_new_exec(), while mm->arg_start/arg_end are set afterwards
// in create_elf_tables() — so a read racing that window returns ZERO bytes for
// a process that provably exists (its stat/pgid read fine microseconds apart)
// and is already executing. Recording that transient emptiness as an empty
// Command later trips Attributed's fail-closed "no recorded command" facet for
// a process we just started (hosted U2: exit -1, empty output, ~1500ms kill
// grace). The bound only tolerates the kernel's publication latency: residual
// exec work after the pipe closes is microseconds of CPU, with preemption
// under load as the variable. Measured, the window is ~15µs at the median but
// tens of milliseconds at the loaded tail (worst independently observed:
// 56.5ms across 10,400 spawns, none past the bound), so 100ms is ~1.8× that
// worst observed sample — not orders of magnitude beyond the tail — and is
// sized to sit under every surrounding budget (KillGroup grace 1500ms,
// WaitDelay 2s/10s): it can only extend paths that already fail, and the
// healthy first-read path pays nothing. A read error (process reaped and
// gone) returns immediately: no amount of retrying makes a reaped pid
// recordable. A killed-but-UNREAPED process is different: /proc/<pid> still
// exists and its cmdline reads zero bytes with no error, so there the poll is
// bounded, not immediate — it spends cmdlinePublishBound and then returns
// ("", false), exactly the pre-poll refusal.
const (
	cmdlinePublishBound = 100 * time.Millisecond
	cmdlinePublishPoll  = time.Millisecond
)

func (linuxProbe) command(pid int) (string, bool) {
	return pollCmdline(func() ([]byte, error) {
		return os.ReadFile("/proc/" + strconv.Itoa(pid) + "/cmdline")
	})
}

// pollCmdline reads NUL-separated argv through read, polling only while a
// live-looking process yields exactly zero bytes (the execve publication
// window above), for at most cmdlinePublishBound. It returns on the FIRST
// non-empty read — it never re-reads after obtaining a value — so it cannot
// turn an untrusted identity into a trusted one: whatever it returns is still
// subject to every unchanged Attributed facet (boot id, alive, exact start
// time, exact process group, session leader, command match). Exhaustion
// returns ("", false): exactly the pre-poll outcome, so a Command that never
// published still records empty and Attributed still refuses it.
//
// Two scope notes reviewers pinned on this shared probe. First, the bound is
// paid by Attributed's LIVE read too: an alive-but-empty-cmdline pid (zombie
// interim, argv-zeroed) used to refuse instantly at "cannot read live
// command"; now that read waits ≤cmdlinePublishBound first and then either
// compares the published argv or refuses at the same facet with the same
// reason — so render-path liveness callers (TUI badges) can spend the bound
// per such pid instead of refusing at once. Second, a re-execing process can
// publish its PRE-exec argv inside the window; the first non-empty value is
// recorded and Attributed then refuses the live mismatch fail-closed — a
// window the pre-poll single read had as well, neither widened nor narrowed
// by polling.
func pollCmdline(read func() ([]byte, error)) (string, bool) {
	deadline := time.Now().Add(cmdlinePublishBound)
	for {
		data, err := read()
		if err != nil {
			return "", false
		}
		if len(data) > 0 {
			// cmdline is NUL-separated argv.
			return strings.TrimSpace(strings.ReplaceAll(string(data), "\x00", " ")), true
		}
		if !time.Now().Before(deadline) {
			return "", false
		}
		time.Sleep(cmdlinePublishPoll)
	}
}
