---
agent: agy
idea: tui-live-steering
phase: review
round: 1
date: 2026-06-06
---

## Summary
The implementation is highly correct and concurrency-safe, but has a major UX gap where the `KILLED` badge is never displayed (rendering as `FAILED` instead) and lacks unit tests for the queued steering path.

## Findings

### MAJOR: Missing KILLED Badge in State Projection
* **Location**: [internal/runstate/runstate.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runstate/runstate.go#L340-L357) & [internal/tui/live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1871)
* **Problem**: The spec (§B) and `IMPLEMENTATION.md` state that a killed agent should display a `KILLED` badge. However, `runstate.ProjectEvents` does not process the `agent.killed` event type. Furthermore, when the runner reaps a killed agent, it appends a standard `agent.failed` event (because it exited with a "killed by user" error), which overwrites the state to `"failed"`. Consequently, the TUI displays a red `FAILED` badge instead of a yellow/orange `KILLED` badge.
* **Fix**: Update `applyAgentEvent` in `internal/runstate/runstate.go` to inspect the `"error"` field of the `agent.failed` event (or add a `"killed"` boolean field to the event data). If it matches `"killed by user"`, map the agent's projected state to `"killed"` instead of `"failed"`.

---

### MINOR: Untested Queued Wait Loop in Steering
* **Location**: [internal/runner/steer_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/steer_test.go)
* **Problem**: The waiter loop in `runSteerAgent` that implements the depth-1 queue (polling `h.agentRunning(req.AgentID)` until the active round attempt finishes) is completely untested.
* **Fix**: Add a unit test `TestRunSteerAttemptQueued` in `steer_test.go` that:
  1. Registers a fake round attempt in `h.active`.
  2. Spawns a steer attempt via `RunSteerAttempt` and asserts it is accepted with status `"queued"`.
  3. Simulates the round completion by calling `h.finish`.
  4. Asserts that the steer attempt transitions to running and completes.

---

### MINOR: Untested Steer Reply Event Projection in TUI
* **Location**: [internal/tui/live_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live_test.go)
* **Problem**: `applySteerReplyEvents` in `live.go` (which transitions the inline steer reply buffer to a terminal state on `steer.replied` / `steer.reply_failed` events) is not exercised by any automated test in `live_test.go`.
* **Fix**: Add a test `TestSteerReplyEventTransitionsState` in `live_test.go` that registers an active steer reply on a model, sends an `eventsMsg` containing a `steer.replied` or `steer.reply_failed` event, and asserts that the `steerReplies` state updates `done=true` and `failed` appropriately.

---

### NIT: Missing Validation for steerBusy Clearing on Steer Failures
* **Location**: [internal/runner/steer_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/steer_test.go)
* **Problem**: Although `runSteerAgent` uses `defer h.clearSteerBusy` to clear the busy flag on exit, there is no unit test confirming that `steerBusy` is correctly cleared when a steer attempt fails (e.g., agent command fails).
* **Fix**: Add a test in `steer_test.go` verifying that `steerBusy` is released when `execAgentProcess` returns an error.

---

### VERIFICATION (Brief Check-ins)
* **Concurrency & Safety**: Per-agent kill is race-safe under `h.mu` (first writer wins vs completion). Steer segment-id allocation is correctly serialized under `h.segmentMu`. Per-steer directories correctly isolate `stdout.log` and `reply.md` without clobbering round assets. No goroutine leaks exist in the wait loop or steer execution path.
* **Decoupling**: Excellent. `internal/tui` does not import `runner` or `app` packages. All communication happens through plain function seams on `LiveOptions` and `LaunchResult`, which are correctly copied in `activateRun`.
* **Keymap & UX**: Conditional-Tab operates correctly (switches tabs only when input is not slash-prefixed). The modal confirm-kill blocks other keys and intercepts commands correctly. Suggestion mode only claims designated navigation/execution keys, leaving printable inputs unaffected.

## Verdict
REQUEST-CHANGES
