---
idea: tui-live-steering
phase: implementation
status: complete
implementer: claude
date: 2026-06-06
---

# IMPLEMENTATION — TUI live steering, kill, autocomplete

Implements FINAL.md across `internal/{runner,app,tui}` (+ tests). No new deps; the
normal round path and `--no-tui` are intact. `internal/tui` imports neither runner nor
app — kill/steer reach it only through injected `LiveOptions` func seams.

## Runner (`internal/runner/runner.go`, new `steer.go`)
- `Handle` extended to own attempt state: stores `opts`, `rootCtx`, `active
  map[string]*attempt`, `steerBusy map[string]bool`, `segmentMu`. `RunRoundOneAsync`
  sets `opts.tracker = handle`.
- `attempt{agentID,segmentID,kind,steerID,cancel,killed}` + `Result.Killed bool`.
- Extracted `execAgentProcess` (CommandFor → stdio → run); `runAgent` now appends
  `agent.started`, registers its cancel via the tracker, runs through
  `execAgentProcess`, and `finish()`es (reporting whether it was killed). Shared by the
  steer path so kill works for both.
- `KillAgent(agentID) KillResult`: cancels ONLY that attempt's child ctx, emits
  `agent.killed` (only when it wins the race vs normal completion), never run-wide cancel.
- `RunSteerAttempt(ctx, SteerAttemptRequest) (SteerAttemptResult, error)`: depth-1 queue
  per agent (second → rejected "already replying"); returns deterministic per-steer paths
  synchronously, runs async. `runSteerAgent` waits out any in-flight round attempt,
  allocates a fresh segment under `segmentMu` (`reason:"steer"`, `steer_id`), emits
  `steer.reply_started`, runs the agent in its OWN dir
  `runs/<id>/agents/<a>/steers/<sid>/{stdout,stderr,reply.md}` (no round-artifact
  skip/fallback/validation), then emits exactly one `steer.replied`/`steer.reply_failed`.
- `BuildSteerPrompt`: steer text + idea `00-prompt.md` + the agent's stdout tail (~4KB),
  with strict "no protocol writes / answer only this" instructions.

## App (`internal/app/app.go`)
- `liveSteerKillSeams(runCtx, handle)` builds `tui.SteerFunc` + `tui.KillAgentFunc`:
  SubmitSteer records via `steer.Submit` first (audit trail) then, for agent targets,
  calls `handle.RunSteerAttempt`; KillAgent calls `handle.KillAgent`. Wired at BOTH TUI
  launch sites (`parley run` ≈1751 and `newLaunchFunc` ≈2039).

## TUI (`internal/tui/live.go`, styles in `app.go`)
- Seams: `SteerRequest`/`SteerResult`/`SteerFunc`/`KillAgentFunc`; `LiveOptions` and
  `LaunchResult` gain `SubmitSteer`+`KillAgent`; `activateRun` copies them (regression
  test) and resets the new state.
- **Autocomplete**: `commandSpecs` table; non-modal `suggest` state recomputed on every
  input edit (visible while input is a bare `/cmd`, ≥1 match, terminal tall enough). Tab =
  longest-common-prefix (single match → full name + space if it takes args); ↑/↓+Enter
  pick; conditional Tab (tab/shift+tab switch tabs only when input isn't slash-prefixed;
  ←/→ always switch). `renderSuggest` slim menu above the input row.
- **Kill**: `ctrl+k` on a running agent tab → modal `confirmKillAgent` (highest-priority
  interceptor, blocks all keys) rendering `kill <agent>? (y/N)`; y/enter → `opts.KillAgent`.
- **Steer**: `submitSteer` calls `opts.SubmitSteer`; on an agent target with a returned
  `StdoutPath`, registers a `steerReply` for that agent. `refreshSteerBuffers` tails the
  attempt stdout; `applySteerReplyEvents` flips done/failed from `steer.replied`/
  `reply_failed`. `renderTranscript` shows the reply inline (divider `── steer <id>:
  "<query>" ──`, faint reply text, "<agent> is replying…" / complete / failed status);
  esc dismisses it back to the transcript. The steer input row prefix is cyan
  (`steerStyle`) so it's never confused with plain input. Observational `/open` runs fall
  back to record-only with a clear message.
- Hints + `/help` updated for steer/kill/autocomplete.

## Tests
- runner (`steer_test.go`, fake agent via `os.Args[0]` re-exec): steer reply captured
  (per-steer dir + `steer.reply_started`/`replied` events), second steer rejected,
  `KillAgent` kills one agent while the other's context survives, kill-vs-completion race
  emits no false event. Passes under `-race`.
- tui (`live_test.go`): Tab single-match + common-prefix completion, suggest clears on
  space, conditional-Tab switches tabs, confirm-kill modal (y/n via the seam),
  `submitInput` steer via a fake seam registers the reply, `activateRun` copies both seams.

## Verification
`go build ./... && go vet ./... && go test ./...` green; `go test -race ./internal/runner`
green for the kill/steer concurrency tests.

## Manual-smoke note
Not run in a live terminal here (model + runner tests exercise the same code headlessly).
To smoke manually: `parley run` (or `parley tui` → N), then on an agent tab type a message
+ Enter → the agent's reply streams into that tab; type `/op`+Tab → completes `/open `;
`ctrl+k` on a running agent → `kill <agent>? (y/N)` → y kills just that agent and the run
continues.

## Fix-up cycle 1 (Phase 8 — addressing review/round-01)
codex + agy REQUEST-CHANGES (hermes ACCEPT). All agreed items fixed:
- **MAJOR — KILLED badge (codex+agy)**: added `runstate.StateKilled` + a sticky
  `AgentState.Killed`; `ProjectEvents`/`applyAgentEvent` now project `agent.killed` as
  killed and keep it across the trailing (canceled) `agent.failed`; segment reset clears
  it. TUI `shortState`→`KILL` and `stateBadge`→warn. Tests: `TestProjectEventsKilledIsSticky`,
  `TestKilledAgentShortState`.
- **MAJOR — ACP agents not killable (codex)**: `runACPAgent` now registers its cancel in
  the attempt registry (and deregisters), so `KillAgent` targets ACP round agents too.
- **MAJOR — steer participant validation (codex)**: `RunSteerAttempt` rejects an agent not
  in `Idea.Participants` (`TestRunSteerAttemptRejectsNonParticipant`).
- **MAJOR — duplicate-kill idempotency (codex)**: a second `KillAgent` while the first is
  in flight is a no-op (no second `agent.killed`) — `TestKillAgentIdempotent`.
- **MINOR — ctx honored (codex)**: `RunSteerAttempt`'s ctx threads into the queue-wait
  loop; a canceled caller ctx (or ended run) emits `steer.reply_failed` and clears busy
  instead of allocating a segment.
- **MINOR — unique steer id (codex)**: an empty `SteerID` now gets a unique
  `steer-auto-NNNN` (no per-steer dir clobber).
- **MINOR — SegmentID contract (codex)**: the steer segment is allocated synchronously, so
  `SteerAttemptResult.SegmentID` is populated on return.
- **MINOR — short-terminal suppress (codex)**: `recomputeSuggest` checks the raw
  (pre-clamp) transcript rows, so the suggest menu is actually suppressed on tiny terminals.
- **Tests added (codex+agy+hermes)**: queued-then-runs, steerBusy-cleared-on-failure,
  non-participant reject, dup-kill, killed projection, steer-reply event→done flip.
- hermes's confirm-vs-suggest ordering NIT was already satisfied (confirm-kill is checked
  before the suggest block); no change.

`go build/vet/test ./...` green; `go test -race ./internal/runner` green. Ready for re-review.

## Fix-up cycle 2 (Phase 8 round-02 → addressing codex MAJOR)
codex round-02 found that the synchronous segment allocation (cycle-1 fix #7) emitted
`run.segment_started` at accept time, which for a QUEUED steer reset/reordered the
still-running round agent's projected state. Fix: a steer no longer emits a
`run.segment_started` boundary at all — it is a side conversation, not a round re-run.
The steer events carry a steer-scoped `segment_id = "steer/<steerID>"` (known
synchronously, no event-count race) purely for correlation, and `SegmentID` is still
populated on the result. Removed the now-unused `Handle.segmentMu`. Regression test
`TestQueuedSteerEmitsNoSegmentBoundary`. hermes ACCEPTed round-02; build/vet/test +
`-race` green.

## Deviations / scope (per FINAL)
- Kill/steer target the async round `Handle` path used by `parley tui`. Fix-up /
  implementation phases (driver-driven, headless) and deck-level steer fan-out are out of
  v1; CLI `parley steer` stays record-only. The steer-reply panel replaces the round
  transcript in that tab while active (esc returns) — chosen over a split pane for clarity
  on small terminals. Orphan grandchildren on kill use ctx cancellation
  (`exec.CommandContext`); process-group signaling is a noted follow-up.
