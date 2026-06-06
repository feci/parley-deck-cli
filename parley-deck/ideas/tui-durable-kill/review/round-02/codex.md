---
agent: codex
idea: tui-durable-kill
phase: review
round: 2
date: 2026-06-07
---

## Summary
ACCEPT. I re-reviewed the fix-up cycle against my round-01 findings and agy's blockers. The dangerous durable-kill paths now fail closed on incomplete or unreadable process identity, ACP agents now persist durable identity and kill process groups, and the Darwin start-time reuse window has been removed with microsecond-resolution `kern.proc.pid` data.

## Verification
FIXED — CRITICAL, `internal/procctl/procctl.go` `Attributed` fail-closed behavior: `Attributed` now refuses unless `PID`, `BootID`, `ProcStart`, `PGID`, and `Command` are all recorded, then separately requires current boot match, liveness, live start-time readability and equality, live pgid readability and equality, `pgid == pid`, and command match. This resolves the partial-metadata kill risk. `TestAttributedFailsClosedOnMissingFields` and `TestDurableKillRefusesUnattributableLiveProcess` cover the key safety case: a live but unattributable PID is not signaled and no kill/clear event is written.

FIXED — CRITICAL, ACP integration: `runACPAgent` now captures `procctl.CaptureByPID` into ACP `agent.started` events, including `pid`, `pgid`, `boot_id`, `proc_start`, `proc_marker`, and `command`. `internal/acp/spawn.go` now routes `Stop` timeout and `Kill` through `killProcessGroup`, with unix `Setsid` preserving the same `pgid == pid` shape as the headless runner. This resolves the missing durable metadata and direct-child-only teardown blocker.

FIXED — MAJOR, Darwin start-time precision: `internal/procctl/procctl_darwin.go` now uses `unix.SysctlKinfoProc("kern.proc.pid", pid)` and stores `p_starttime` as seconds plus microseconds. That replaces `ps lstart`'s one-second resolution and materially strengthens PID-reuse attribution on the owner's macOS host.

FIXED — MAJOR, `KillGroup` self-protection: `internal/procctl/procctl_unix.go` refuses to signal if the target group equals parley's own process group. This covers the Setsid-skipped/failed edge.

FIXED — MAJOR, Home stale banner: `internal/tui/live.go` implements `staleAgentCount` and renders the stale-agent Home warning when projected-running agents are stale by the liveness seam.

FIXED — MINOR, `commandMatches`: matching is now only exact match or `strings.HasPrefix(live, recorded)`, so a longer recorded value no longer weakens the live command check.

FIXED — MINOR, refusal handling: `KillAgentDurable` clears only provably dead PIDs. Alive but unattributable PIDs return `Failed` without writing `agent.killed` or synthetic `agent.failed`, and `killOutcome` maps that to a Go error for red TUI display.

FIXED — NIT, `agent.started` timestamp: the headless path stamps `agent.started` inside `onStarted` after `cmd.Start()` and `procctl.Capture`, rather than using the pre-setup result timestamp.

## New findings
None.

## Verdict
ACCEPT.
