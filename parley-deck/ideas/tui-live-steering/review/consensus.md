---
idea: tui-live-steering
phase: review-consensus
drafter: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
---

## Review consensus

Phase 6 (round-01) reviews → fix-up cycle 1 → Phase 8 re-review (round-02) → fix-up
cycle 2 → re-review (round-03). The implementation is accepted: **codex and hermes both
✅ ACCEPT at round-03 with zero remaining findings.** agy's round-01 findings are all
FIXED and independently verified; agy's round-02/03 re-review artifacts could not be
produced (recorded headless tooling exception — see
`inbox/claude-to-all_tui-live-steering_agy-rereview-tooling-exception.md`).

### Agreed fixes — fix-up cycle 1 (all applied, verified FIXED in round-02/03)
- **MAJOR — KILLED badge** (codex + agy): sticky `runstate.StateKilled` projected from
  `agent.killed`, preserved across the trailing `agent.failed`, reset on a real new
  segment; TUI `shortState`→`KILL`, `stateBadge`→warn.
- **MAJOR — ACP agents not killable** (codex): `runACPAgent` registers/deregisters its
  cancel in the attempt registry.
- **MAJOR — steer participant validation** (codex): `RunSteerAttempt` rejects
  non-participants.
- **MAJOR — duplicate-kill idempotency** (codex): a second in-flight `KillAgent` is a
  no-op (one `agent.killed`).
- **MINOR** — ctx honored in the queue-wait loop; unique `steer-auto-NNNN` id;
  `SegmentID` populated; short-terminal suggest suppression uses raw pre-clamp rows.
- **Tests** — queued-then-runs, busy-cleared-on-failure, non-participant reject, dup-kill,
  killed projection (sticky + reset), steer-reply event→done flip, killed badge.

### Agreed fix — fix-up cycle 2 (codex round-02 MAJOR, verified FIXED round-03)
- A steer no longer emits `run.segment_started` (it is a side conversation, not a round
  re-run); previously the synchronous boundary reset/reordered the still-running round
  agent. Steer events carry a steer-scoped `segment_id = "steer/<steerID>"` for
  correlation; `SegmentID` still populated. Regression test
  `TestQueuedSteerEmitsNoSegmentBoundary`. Removed unused `Handle.segmentMu`.

### Dismissed / no-change
- hermes round-01 confirm-vs-suggest ordering NIT — already satisfied (confirm-kill is
  the highest-priority interceptor, before the suggest block).

### Deferred follow-ups (non-blocking)
- Fix-up / implementation phases and deck-steer fan-out kill/steer (driver-driven,
  headless) are out of v1. CLI `parley steer` stays record-only. Process-group kill for
  orphan grandchildren (v1 uses ctx cancellation). agy round-03 re-review when its CLI
  recovers (no outstanding fix).

### Verification
`go build ./... && go vet ./... && go test ./...` green; `go test -race ./internal/runner`
green.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-06. As implementer I confirm both fix-up cycles are applied and green
(build/vet/test + -race). All agreed fixes done; no outstanding agreed fix. Proceeding
to release on codex + hermes ACCEPT with agy's findings resolved and its re-review
tooling exception recorded.

### agy — ⚠️ FINDINGS RESOLVED (round-03 artifact not produced — tooling)
agy's round-01 code review (KILLED-badge MAJOR + test gaps) is fully implemented and
verified FIXED by codex and hermes in rounds 02/03. agy's round-02/03 re-review could
not be generated (3 attempts, empty headless output — recorded exception). No agreed
fix is outstanding. This block is the facilitator's record, not a proxy signoff.

### codex — ✅ ACCEPT
The agreed fix-up cycles are in, including my round-01 and round-02 findings, and the recorded verification is green. I accept marking the idea complete and releasing with no outstanding agreed fixes.

### hermes — ✅ ACCEPT
Keymap and concurrency invariants hold per round-03 verification; all findings fixed with green build/vet/test -race. I accept releasing.
