---
agent: codex
idea: tui-live-steering
phase: review
round: 1
date: 2026-06-06
---

## Summary
REQUEST-CHANGES: the broad TUI seam/autocomplete/steer shape is present, but runner correctness and the killed-state UX still miss important parts of FINAL.md.

## Findings

MAJOR — `internal/runner/runner.go:352` / `internal/runner/acp.go:27`: ACP-launched round agents bypass the new per-attempt registration path entirely. `runAgent` returns into `runACPAgent` before `opts.tracker.register(...)`, so `Handle.KillAgent` cannot kill ACP agents and queued steers do not see ACP agents in `active`. FINAL §B says the round path uses the shared registration helper so round agents are killable; this is not true for `LaunchACP`. Fix by registering the ACP attempt after `agent.started` and before the ACP process/session work, deregistering through the same first-writer `finish` path, and setting `Result.Killed`/`ExitError` when the ACP context was canceled by `KillAgent`.

MAJOR — `internal/runstate/runstate.go:339` and `internal/tui/live.go:615`: `agent.killed` is emitted but never projected, and the later canceled attempt is rendered as ordinary `failed`/`ERR`. FINAL §B requires a distinct killed badge (`KILLED`), and `Result.Killed` currently dies inside the in-memory result without reaching the event model. Fix by carrying killed state into events, for example add `killed:true` to the terminal `agent.failed` event or project `agent.killed` as a distinct terminal/flag state, then render `KILLED` in `shortState`/`stateBadge` and cover it in runstate/TUI tests.

MAJOR — `internal/runner/steer.go:117`: `RunSteerAttempt` validates only that the agent is discovered and found, not that it is an idea participant. FINAL §C explicitly says to validate the agent is a participant. As written, any installed discovered agent can be steered by calling the runner seam directly, creating events and files under the run for an agent outside the idea quorum. Fix by checking `req.AgentID` against `h.opts.Idea.Participants` before accepting the steer.

MAJOR — `internal/runner/steer.go:89`: concurrent duplicate `KillAgent` calls against the same active attempt can both append `agent.killed`, because the first caller marks `att.killed=true` but leaves the attempt in `active` until the process exits, and the second caller does not check that flag before appending. The kill-vs-finish race is mostly first-writer safe, but kill-vs-kill is not idempotent. Fix under `h.mu`: if `att.killed` is already true, return a no-op/already-killing result without emitting another event.

MINOR — `internal/runner/steer.go:113` and `internal/runner/steer.go:161`: the `ctx` parameter to `RunSteerAttempt` is ignored, so a caller-canceled queued request still waits until `rootCtx` is canceled or the active attempt finishes, then may allocate a segment and emit `steer.reply_started`/`reply_failed` for a request the caller abandoned. App wiring passes `runCtx`, so run cancellation eventually stops it, but the public runner contract should honor the supplied context. Fix by threading `ctx` into `runSteerAgent`, selecting on it in the wait loop, and returning/clearing `steerBusy` without starting a segment when it is canceled.

MINOR — `internal/runner/steer.go:121`: an empty `SteerID` is normalized to the constant `steer-attempt`, so repeated direct `RunSteerAttempt` calls with no ID reuse the same `agents/<id>/steers/steer-attempt/` directory and can clobber stdout/stderr/reply files. The app path passes the unique `steer.Submit` ID, but the runner API itself violates FINAL's per-steer isolation expectation. Fix by generating a unique steer ID in the runner when the request does not provide one.

MINOR — `internal/runner/steer.go:34` / `internal/runner/steer.go:143`: `SteerAttemptResult.SegmentID` is always empty on the synchronous return path because the segment is allocated later in the goroutine. FINAL §C includes `SegmentID` in the runner/app/TUI result contract. Either allocate the segment before returning for immediate attempts and report queued segment later by event only, or remove the synchronous `SegmentID` field from the contract; as ratified, the implementation does not satisfy it.

MINOR — `internal/tui/live.go:908` and `internal/tui/live.go:1247`: short-terminal suppression for the suggest menu is ineffective. `recomputeSuggest` checks `m.transcriptHeight() < 3`, but `transcriptHeight` clamps the value to at least 3, so the condition can never be true and suggestions are never suppressed for very short terminals. Fix by checking the raw available height before clamping or adding a dedicated `suggestRowsAvailable` helper.

MINOR — `internal/tui/live_test.go:951`: the tests cover the happy paths but miss several FINAL §E contracts and the bugs above: queued steer actually waiting then starting, `steerBusy` clearing after failed/canceled queued paths, supplied-context cancellation, non-participant rejection, duplicate kill idempotency, killed badge projection, conditional Tab when slash-prefixed, and `steer.replied` event flipping the inline reply from replying to done. Add focused tests for these before accepting.

Correct — `internal/tui/live.go:49` / `internal/tui/live.go:430` and `internal/app/app.go:2018`: the TUI remains decoupled from `internal/runner`/`internal/app` via plain function seams, and `activateRun` copies both `SubmitSteer` and `KillAgent`, including Home-launched runs.

Correct — `internal/runner/steer.go:175`: steer segment ID allocation for concurrent steer attempts is serialized under `segmentMu`, which addresses the `nextSegmentID` collision called out in FINAL.

Correct — `internal/runner/steer.go:138` / `internal/runner/steer.go:196`: steer attempts use per-steer stdout/stderr/reply paths and do not reuse the round artifact skip/fallback/validation path.

Correct — `internal/tui/live.go:1026` / `internal/tui/live.go:1034`: the confirm-kill modal is ordered after only `ctrl+c` and before picker/suggest/input handling, so it blocks other keys as required.

## Verdict
REQUEST-CHANGES: block on the ACP killability gap, missing killed-state projection/badge, participant validation for steer, and duplicate kill idempotency. The remaining MINOR items can be fixed in the same pass or explicitly deferred by review consensus.
