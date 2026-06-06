// Package procctl is parley's cross-platform process-control layer: it spawns
// agents in their own process group, captures a durable process identity, and
// kills the whole process tree — including across a parley restart, gated by a
// strict attribution check so a reused PID is never killed by mistake.
//
// OS specifics live in build-tagged files: procctl_unix.go (group spawn + kill),
// procctl_linux.go and procctl_darwin.go (the identity probe), procctl_windows.go.
package procctl

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Spawned is the durable identity of one agent process. It is persisted in the
// agent.started event so a restarted parley can re-attribute and kill it.
type Spawned struct {
	PID       int       `json:"pid"`
	PGID      int       `json:"pgid,omitempty"`
	BootID    string    `json:"boot_id,omitempty"`
	ProcStart string    `json:"proc_start,omitempty"` // captured via the same probe used to verify
	Command   string    `json:"command,omitempty"`
	Marker    string    `json:"marker,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
}

// probe is the OS-specific identity surface (swappable in tests). All methods
// return portable types so this file stays build-tag-free.
type probe interface {
	bootID() string
	alive(pid int) bool
	procStart(pid int) (string, bool)
	command(pid int) (string, bool)
	pgid(pid int) (int, bool)
	supportsDurableKill() bool
}

// active is set by the per-OS file's init(); tests may swap it.
var active probe

// MarkerEnv returns the identifying env vars set on every spawned agent (read
// back on Linux via /proc/<pid>/environ as defense-in-depth; not required on
// darwin where start-time + pgid + command are the decisive guards).
func MarkerEnv(runID, agentID, marker string) []string {
	return []string{
		"PARLEY_RUN_ID=" + runID,
		"PARLEY_AGENT_ID=" + agentID,
		"PARLEY_PROC_MARKER=" + marker,
	}
}

// CurrentBootID identifies the current OS boot; "" when unsupported.
func CurrentBootID() string {
	if active == nil {
		return ""
	}
	return active.bootID()
}

// Alive reports whether the recorded PID currently exists (identity unchecked).
func Alive(s Spawned) bool {
	return active != nil && s.PID > 0 && active.alive(s.PID)
}

// SupportsDurableKill reports whether cross-restart kill is available on this OS.
func SupportsDurableKill() bool { return active != nil && active.supportsDurableKill() }

// Capture records the durable identity of a just-started command. Call it
// immediately after cmd.Start(); ProcStart/Command come from the same probe that
// Attributed uses later, so capture and verify compare identically.
func Capture(cmd *exec.Cmd, marker string) Spawned {
	if cmd == nil || cmd.Process == nil {
		return Spawned{Marker: marker, StartedAt: time.Now().UTC()}
	}
	return CaptureByPID(cmd.Process.Pid, marker)
}

// CaptureByPID records the durable identity of an already-running pid (used by the
// ACP path, which spawns via its own helper).
func CaptureByPID(pid int, marker string) Spawned {
	sp := Spawned{PID: pid, Marker: marker, StartedAt: time.Now().UTC()}
	if active == nil || pid <= 0 {
		return sp
	}
	sp.BootID = active.bootID()
	if pg, ok := active.pgid(pid); ok {
		sp.PGID = pg
	}
	if ps, ok := active.procStart(pid); ok {
		sp.ProcStart = ps
	}
	if c, ok := active.command(pid); ok {
		sp.Command = c
	}
	return sp
}

// Attributed is the cross-restart safety gate: it returns (true,"") ONLY when the
// live PID is provably the same process we recorded. It FAILS CLOSED — every
// identity facet must be recorded AND readable from the live process and match:
// same boot, alive, exact (microsecond) start time, same process group, the
// process is its own session leader (our Setsid signature), and matching command.
// Any miss returns (false, reason) and the caller MUST NOT signal — this is what
// prevents killing a reused/unrelated PID.
func Attributed(s Spawned) (bool, string) {
	if active == nil || !active.supportsDurableKill() {
		return false, "durable kill is unsupported on this platform"
	}
	// Fail closed: an incomplete recorded identity can never be attributed.
	if s.PID <= 0 {
		return false, "no pid recorded"
	}
	if s.BootID == "" {
		return false, "no recorded boot id"
	}
	if s.ProcStart == "" {
		return false, "no recorded process start time"
	}
	if s.PGID <= 0 {
		return false, "no recorded process group"
	}
	if s.Command == "" {
		return false, "no recorded command"
	}
	if s.BootID != active.bootID() {
		return false, "boot id differs (process is from a prior boot)"
	}
	if !active.alive(s.PID) {
		return false, "process is not running"
	}
	ps, ok := active.procStart(s.PID)
	if !ok {
		return false, "cannot read live start time"
	}
	if ps != s.ProcStart {
		return false, "process start time mismatch (pid was reused)"
	}
	pg, ok := active.pgid(s.PID)
	if !ok {
		return false, "cannot read live process group"
	}
	if pg != s.PGID {
		return false, "process group mismatch (pid was reused)"
	}
	if pg != s.PID {
		return false, "process is not its own session leader (not one of ours)"
	}
	c, ok := active.command(s.PID)
	if !ok {
		return false, "cannot read live command"
	}
	if !commandMatches(c, s.Command) {
		return false, "command mismatch (pid was reused)"
	}
	return true, ""
}

// commandMatches accepts an exact match or the live command extending the
// recorded one (tolerating ps/cmdline truncation of the live read). It does NOT
// accept a shorter recorded command matching a longer live one (that would let a
// truncated record match an unrelated process).
func commandMatches(live, recorded string) bool {
	live = strings.TrimSpace(live)
	recorded = strings.TrimSpace(recorded)
	if live == "" || recorded == "" {
		return false
	}
	return live == recorded || strings.HasPrefix(live, recorded)
}

// KillTreeAttributed gates on Attributed, then group-kills the tree. Use this on
// the reattach path (a restarted/observational parley). Returns (killed, reason).
func KillTreeAttributed(s Spawned) (bool, string, error) {
	if ok, reason := Attributed(s); !ok {
		return false, reason, nil
	}
	if err := KillGroup(s); err != nil {
		return false, "kill failed: " + err.Error(), fmt.Errorf("kill group: %w", err)
	}
	return true, "", nil
}
