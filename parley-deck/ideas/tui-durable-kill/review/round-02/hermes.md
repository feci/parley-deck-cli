---
agent: hermes
idea: tui-durable-kill
phase: review
round: 2
date: 2026-06-07
---

## Summary
ACCEPT. Fix-up cycle 1 fully resolves every blocking item from round-01 (codex + agy). Attributed is now fail-closed, ACP is integrated with durable identity and group-kill, darwin start-time is microsecond-resolution, KillGroup self-protects, stale banner exists, commandMatches direction is correct, and refusal semantics are fixed. All tests, builds, cross-compiles, and real macOS process tests pass.

## Verification

**CRITICAL Attributed missing fields (codex+agy)**: FIXED. `Attributed` now refuses unless BootID, ProcStart, PGID, Command are all non-empty in the record AND all live probes succeed and match exactly (plus pgid==pid session-leader check). New tests `TestAttributedFailsClosedOnMissingFields` and `TestDurableKillRefusesUnattributableLiveProcess` confirm a live but unattributable PID is never signaled and produces no event.

**CRITICAL ACP integration (agy)**: FIXED. `runACPAgent` calls `procctl.CaptureByPID` and stores full identity in `agent.started`; `acp.Process.Stop/Kill` now call `killProcessGroup` (sysproc_{unix,windows}.go) to reap the whole tree.

**MAJOR darwin 1s resolution (codex+agy)**: FIXED. `procStart` on darwin now uses `unix.SysctlKinfoProc("kern.proc.pid", pid).Proc.P_starttime` at microsecond resolution (no new dependency).

**MAJOR KillGroup self-protection (agy)**: FIXED. `KillGroup` guards with `syscall.Getpgid(0)` and refuses to signal parley's own process group.

**MAJOR Home stale banner (agy)**: FIXED. `renderHome` now renders `⚠ N stale agent process(es)…` via `staleAgentCount`.

**MINOR commandMatches direction (codex+agy)**: FIXED. Only accepts `live == recorded || HasPrefix(live, recorded)`.

**MINOR refusal auto-clearing (codex+agy)**: FIXED. Attribution refusal returns error (TUI shows red) without writing terminal event; clear-stale only on provably-dead PID. `agent.started` now stamped at actual start via `onStarted`.

**NIT agent.started timestamp (codex+agy)**: FIXED. Timestamp captured inside `onStarted`.

## New findings
None. Live in-memory and watcher paths remain ungated (correct trust boundary); only the reattach durable path uses Attributed. Self-attribution tests pass on this macOS host.

## Verdict
ACCEPT. All round-01 blocking items resolved; no regressions or new concerns.