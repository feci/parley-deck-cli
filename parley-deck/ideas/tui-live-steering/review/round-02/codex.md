---
agent: codex
idea: tui-live-steering
phase: review
round: 2
date: 2026-06-06
---

## Summary
REQUEST-CHANGES. The fix-up cycle addresses the round-01 blocking items: killed projection is sticky and resets on a new segment, ACP attempts are registered for per-agent cancellation, steer participant validation is present, duplicate kill is idempotent, queue-wait cancellation clears busy and emits a failure, and the added tests cover the reported gaps.

However, the synchronous steer segment allocation introduced a new ordering bug for queued steers. `RunSteerAttempt` emits `run.segment_started` at accept time, before the active round attempt has finished. That makes the event log start a newer segment while older-segment `agent.finished`/`agent.failed` events are still expected, which breaks the projection ordering assumptions and can make the TUI lose the fact that the active round agent is still running.

## Verification
FIXED — ACP round agents not killable. `runACPAgent` now creates its timeout context, registers `opts.tracker.register(agent.ID, opts.SegmentID, "round", "", cancel)`, and deregisters with `defer opts.tracker.finish(agent.ID)`. The nil-tracker path is safe because the methods guard `h == nil`, so the sync path is not broken.

FIXED — Missing killed-state projection and badge. `runstate.StateKilled` and `AgentState.Killed` are present; `agent.killed` sets killed state, trailing `agent.failed` preserves it, and `run.segment_started` clears it for targeted agents. `shortState(stateKilled)` returns `KILL`, and `stateBadge("killed")` uses the warning badge. `TestProjectEventsKilledIsSticky` and `TestKilledAgentShortState` cover this.

FIXED — Steer participant validation. `RunSteerAttempt` now rejects discovered-but-non-participant agents with a rejected result before setting `steerBusy` or creating files. Covered by `TestRunSteerAttemptRejectsNonParticipant`.

FIXED — Duplicate kill idempotency. `KillAgent` checks `att.killed` under `h.mu` and returns a no-op without emitting a second `agent.killed`. Covered by `TestKillAgentIdempotent`.

FIXED — Caller context ignored in queue wait. `runSteerAgent` checks both the supplied `ctx` and `rootCtx` while waiting, emits `steer.reply_failed`, and clears `steerBusy` via defer.

FIXED — Empty steer ID clobbers paths. Empty `SteerID` now generates `steer-auto-NNNN` under `h.mu`, avoiding repeated `steer-attempt` paths.

FIXED — Synchronous `SegmentID` result population. `RunSteerAttempt` now returns a populated `SegmentID` after serializing `nextSegmentID` under `segmentMu`, which avoids the earlier empty-result contract.

FIXED — Short-terminal suggest suppression. `recomputeSuggest` checks raw available rows via `tuiHeight(m.height, defaultLiveHeight)-7`, so the pre-clamp small-terminal case is handled.

FIXED — Test gaps from round-01. Added coverage includes queued-then-runs, busy clearing after steer failure, non-participant rejection, duplicate kill, killed projection, killed short state, and steer reply event transitions. I also ran `go test ./internal/runner ./internal/runstate ./internal/tui` and `go test -race ./internal/runner`; both passed.

new concern — queued steer segment ordering. See MAJOR finding below.

## New findings
MAJOR — `internal/runner/steer.go:155`: queued steers emit `run.segment_started` immediately, before the current active round attempt has finished. This is not just an audit-ordering detail: `ProjectEvents` resets targeted agents on `run.segment_started`, but it does not ignore later `agent.finished`/`agent.failed` events from an older segment. A queued steer therefore creates an event order like `agent.started(segment-0001)`, `run.segment_started(segment-0002, reason=steer)`, then later `agent.finished(segment-0001)`. The projection can show the running round agent as pending/finished while it is still active, and `ctrl+k` becomes unavailable because the TUI only opens kill confirmation when projected state is running. This also undermines the segment reset contract by allowing older-segment terminal events to overwrite the newer segment state.

Fix by reserving/returning a unique steer segment ID without emitting `run.segment_started` until the queued steer actually starts, or by making projection ignore stale older-segment agent events after a newer segment boundary. The former is cleaner: keep the returned `SegmentID` contract, but append the segment boundary immediately before `steer.reply_started`, after the queue wait.

## Verdict
REQUEST-CHANGES. Blocking item: fix the queued-steer segment ordering so accepting a queued steer does not reset or overwrite the active round agent state before the steer actually starts. Once that is corrected and covered by a projection/order test, I expect to ACCEPT.
