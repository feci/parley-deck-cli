---
idea: tui-durable-kill
phase: review-consensus
drafter: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
---

## Review consensus

Phase 6 (round-01) → fix-up cycle 1 → Phase 8 re-review (round-02). **All three reviewers
(codex, agy, hermes) ✅ ACCEPT at round-02 with zero remaining findings.** Ready to mark
complete.

### Agreed fixes (fix-up cycle 1 — all applied, verified FIXED in round-02)
- **CRITICAL — `Attributed` fail-closed** (codex + agy): refuses unless every identity
  facet is recorded AND readable from the live process AND matches (boot-id, microsecond
  start-time, pgid, command) + the process is its own session leader (`pgid==pid`). A live
  unattributable PID is never signaled and writes no event. Tests
  `TestAttributedFailsClosedOnMissingFields`, `TestDurableKillRefusesUnattributableLiveProcess`.
- **CRITICAL — ACP integration** (agy): `runACPAgent` captures the durable identity into
  the ACP `agent.started`; `acp.Process.Stop/Kill` route through `killProcessGroup` (reaps
  the tree).
- **MAJOR — darwin microsecond start-time** (codex + agy): `unix.SysctlKinfoProc`
  `kern.proc.pid` p_starttime (no new dep — `x/sys` was already in the module graph),
  closing the `ps lstart` 1-second reuse window.
- **MAJOR — `KillGroup` self-protection** (agy): refuses to signal parley's own group.
- **MAJOR — Home stale banner** (agy): `staleAgentCount` + the `⚠ N stale…` Home warning.
- **MINOR — `commandMatches`** (codex + agy): only `live==recorded || HasPrefix(live,recorded)`.
- **MINOR — refusal handling** (codex + agy): clear-stale only on a provably-dead pid; an
  alive-but-unattributable agent returns `Failed` → a Go error → red TUI (no event).
- **NIT** (codex + agy): `agent.started` stamped at actual process start.

### Dismissed / documented (agreed)
- **NIT — steer attempts don't persist process identity** (agy): accepted as a documented
  limitation (steers are transient follow-up replies); not durably killable across restart.

### Deferred follow-ups (non-blocking, out of v1)
- windows durable kill refuses (live-handle kill only). Full resumed-run steer EXECUTION
  and the larger steer-conversation-from-events rendering remain follow-ups (the
  durable-kill ask is fully delivered). linux `/proc` start-time parse + the watcher
  group-kill-on-timeout path are covered indirectly; deeper linux integration tests welcome.

### Verification
`go build ./... && go vet ./... && go test ./...` green; `go test -race ./internal/runner
./internal/procctl` green; windows + linux cross-compile; the procctl + durable-kill tests
exercise REAL processes (group kill reaps a grandchild; tampered identity refused; live
end-to-end durable kill via the event log) on this macOS host.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-07. As implementer I confirm fix-up cycle 1 is applied and green
(build/vet/test + -race + cross-compile; real-process tests on macOS). All agreed fixes
done; no outstanding agreed fix. Ready to mark complete and release.

### codex — ✅ ACCEPT
I confirm the agreed fix-up cycle is in: durable kill now fails closed on process attribution, ACP uses durable process-group teardown, Darwin start-time precision is tightened, and the remaining safety/UI fixes are verified green. I accept completing and releasing this work.

### agy — ✅ ACCEPT
I verify that all Phase 6 fixes—including fail-closed process attribution, microsecond-resolution Darwin start times, and self-kill prevention—have been successfully integrated and validated. The macOS integration tests confirm robust process-group cleanup and TUI stale banner behavior. I accept completing and releasing this work.

### hermes — ✅ ACCEPT
All agreed fixes from fix-up cycle 1 are verified in place (fail-closed attribution, ACP integration, microsecond darwin start-time, KillGroup self-protection, Home stale banner) with round-02 ACCEPT recorded. Real-process tests green on macOS. I accept completing and releasing this work.
