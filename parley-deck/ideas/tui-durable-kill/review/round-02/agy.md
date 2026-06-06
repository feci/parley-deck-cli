---
agent: agy
idea: tui-durable-kill
phase: review
round: 2
date: 2026-06-07
---

## Summary
The implementer (claude) has successfully completed fix-up cycle 1, addressing all critical, major, minor, and nit issues identified in the round-01 reviews. All process group reaping, safety checks (fail-closed attribution and self-kill prevention), and Darwin-specific high-resolution timing enhancements have been verified. The test suite, including real OS integration tests, passes successfully on macOS darwin. I ACCEPT the implementation.

## Verification

### 1. CRITICAL: Attributed Safety Gate (fail-closed)
- **Status**: FIXED
- **Verification**: `procctl.go` [Attributed](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl.go#L107-L158) now explicitly verifies that all identity fields (`PID`, `BootID`, `ProcStart`, `PGID`, `Command`) are non-empty/non-zero before matching. If any field is missing in the persisted record, it fails closed immediately. It also checks that `pg == s.PID`, enforcing that the target process is its own session leader (conforming to the `Setsid` spawn configuration). The test suite includes `TestAttributedFailsClosedOnMissingFields` and `TestDurableKillRefusesUnattributableLiveProcess` to ensure that PID reuse/unattributable states do not result in signals being sent or state being incorrectly cleared.

### 2. CRITICAL: ACP Integration
- **Status**: FIXED
- **Verification**: `runACPAgent` in [acp.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/acp.go#L83-L104) now captures process metadata with `procctl.CaptureByPID` and records it into the `agent.started` event. In [spawn.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/acp/spawn.go#L106-L139), `Stop()` and `Kill()` now invoke `killProcessGroup(p.cmd.Process.Pid)` which routes through `sysproc_unix.go` to group-kill using a negative PID. This ensures that grandchild processes spawned by ACP agents are correctly reaped upon cancellation/timeout.

### 3. MAJOR: Darwin Precise Start Time
- **Status**: FIXED
- **Verification**: macOS timing in [procctl_darwin.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl_darwin.go#L52-L59) has been upgraded from `ps lstart` (1-second resolution) to use `unix.SysctlKinfoProc` querying `kern.proc.pid` at microsecond resolution (`kp.Proc.P_starttime`), closing the process timing reuse gap.

### 4. MAJOR: Home Stale Banner
- **Status**: FIXED
- **Verification**: `renderHome` in [live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L872-L875) checks `m.staleAgentCount()` and renders a warning banner `⚠ N stale agent process(es) — open the agent's tab and ctrl+k to clear` if any running agents in the active run are stale.

### 5. MAJOR: KillGroup Self-Protection
- **Status**: FIXED
- **Verification**: `KillGroup` in [procctl_unix.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl_unix.go#L36-L38) now compares the target PGID against `syscall.Getpgid(0)`. If they match, it immediately exits with an error refusing to signal the parent parley session.

### 6. MINOR: `commandMatches` Match Direction
- **Status**: FIXED
- **Verification**: `commandMatches` in [procctl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl.go#L164-L171) has been updated to check `live == recorded || strings.HasPrefix(live, recorded)`. A shorter recorded command can no longer erroneously match a longer live command.

### 7. MINOR: Refusal Handling and TUI Color
- **Status**: FIXED
- **Verification**: Alive-but-unattributable process verification failures in `KillAgentDurable` return `Failed: true` with a failure message. This is adapted via `killOutcome` in [app.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L2072-L2086) to return a Go error, which is caught by the TUI and rendered as a red warning banner rather than auto-clearing. Only dead processes clear the stale status.

### 8. NIT: agent.started Timestamp
- **Status**: FIXED
- **Verification**: `onStarted` in [runner.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L365-L383) is called in `execAgentProcess` immediately after `cmd.Start()`, aligning the timestamp with actual process startup rather than pre-setup.

### 9. NIT: Steer Agent Process Metadata
- **Status**: FIXED (Acceptable Limitation / Documented)
- **Verification**: Documented in `IMPLEMENTATION.md` that steer attempts do not persist process metadata or append `agent.started` events since they are transient.

## New findings
None.

## Verdict
ACCEPT
