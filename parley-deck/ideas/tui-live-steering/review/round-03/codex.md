---
agent: codex
idea: tui-live-steering
phase: review
round: 3
date: 2026-06-06
---

## Summary
ACCEPT. Fix-up cycle 2 resolves my round-02 MAJOR: steer attempts no longer emit `run.segment_started`, so accepting a queued steer cannot reset or reorder the still-running round agent's projected state.

The steer path now treats the reply attempt as a side conversation. `RunSteerAttempt` returns a deterministic steer-scoped `SegmentID` (`steer/<steerID>`) synchronously for correlation, and `runSteerAgent` emits only `steer.reply_started`, `steer.replied`, or `steer.reply_failed` with that label.

## Verification
Verified `internal/runner/steer.go`: the prior synchronous segment allocation was removed. `RunSteerAttempt` derives `seg := "steer/" + steerID`, returns it in `SteerAttemptResult.SegmentID`, and starts the goroutine without appending any run segment boundary. The queued path waits while `agentRunning(req.AgentID)` remains true, so the active round attempt keeps its existing segment and projected running state until it naturally finishes or is killed.

Verified regression coverage: `TestQueuedSteerEmitsNoSegmentBoundary` simulates an active round attempt, accepts a queued steer, asserts the returned segment label is `steer/steer-q2`, and fails if any `run.segment_started` event appears. Existing steer tests still cover queued-then-runs, duplicate steer rejection, failure cleanup, participant validation, unique empty steer IDs, and reply capture.

Verified correlation is still intact: `internal/app/app.go` records the durable steer first, passes `rec.ID` into `RunSteerAttempt`, and returns the runner's `SegmentID`/`StdoutPath` to the TUI. `internal/tui/live.go` tracks inline replies by agent plus steer id and flips state from `steer.replied`/`steer.reply_failed`; it does not require `runstate` to interpret steer labels as run segments.

Re-verified prior fixes remain present: `agent.killed` projects to sticky `StateKilled`, trailing `agent.failed` preserves the killed badge, actual `run.segment_started` boundaries clear killed state for real reruns, ACP attempts register with the kill tracker, duplicate kills are idempotent, queue wait observes caller/run cancellation, and short-terminal autocomplete suppression uses the raw pre-clamp height.

I ran `go build ./...`, `go vet ./...`, `go test ./...`, and `go test -race ./internal/runner`; all passed.

## New findings
None.

## Verdict
ACCEPT. The round-02 blocking issue is resolved, and I found no new blocking or non-blocking issues in this re-review.
