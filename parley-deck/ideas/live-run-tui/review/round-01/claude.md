---
agent: claude
idea: live-run-tui
review-round: 1
date: 2026-05-10
implementation-pr: https://github.com/feci/parley-deck-cli/pull/2
---

## Verdict
Request changes

## Findings
- MAJOR internal/tui/live.go:402-406 `agent.skipped` projection sets `agent.Duration = now.Sub(agent.StartedAt)`. The runner emits `agent.skipped` *without* a preceding `agent.started` when the artifact already exists (internal/runner/runner.go:188-201, payload contains only `agent`/`reason`/`artifact` — no `duration_ms`). `StartedAt` therefore stays at the zero time and `Duration` becomes ~time-since-year-1, which `formatDuration` will render as a multi-million-hour string in the agent table's ELAPSED column. This regresses the "skipped" UX every time `parley run` is invoked over an idea that already has round-01 artifacts (a common case).
- MAJOR internal/tui/live.go:349-369 Unknown-agent state is not retained across events. On the first event for an unmapped agent the code creates `AgentState{State: stateUnknown}`, sets `LatestEvent`, and `continue`s — but on any *subsequent* event for the same agent it now finds the entry in `agentsByID` and falls through to `applyAgentEvent`, mutating state to `running`/`finished`/`failed`/`skipped`. FINAL.md defines `unknown` as "agent is not part of the current selected round or cannot be mapped from run metadata", which is a permanent classification for the run; the current behavior turns `unknown` into a one-event placeholder and a non-participant emitting two events will silently appear as a normal participant. Either drop the entry into `applyAgentEvent` for unknown agents or guard the call so `unknown` is sticky.
- MINOR internal/tui/live.go:217-246 `renderAgentTable` formats `state` with `%-16s` *after* wrapping it in a Lipgloss style. The injected ANSI escape bytes are counted toward the column width, so any colored state (finished/failed/running) shifts the ELAPSED and LAST EVENT columns left and breaks the table alignment header introduced on line 220. Apply padding before styling, or use Lipgloss's `Width(...)` on the styled value.
- MINOR internal/tui/live.go:562-570 `displayRoundStatus` synthesizes the strings `running` and `unknown` for the round status header before any `round.completed`/`round.incomplete` event. FINAL.md says "Round status must be displayed from `round.completed` or `round.incomplete` events only"; while this isn't an inferred *completion*, surfacing strings that overlap with agent states (`running`/`unknown`) is misleading. Consider a clearly distinct placeholder such as `awaiting` or `—`.
- MINOR internal/tui/live.go:456-494 On a JSON parse error `readEventsFromOffset` returns the *original* offset, so the next 250 ms tick will re-read the same line and surface the same `errText` indefinitely. There is no resync mechanism. Acceptable for the slice but worth either advancing past the broken line or capping retries.
- MINOR internal/runner/runner.go:60-78 The panic recovery emits `run.failed` but the TUI does not handle that event at all (live.go:349-368 only switches on agent and round events), so a crashed runner leaves the round status header at `running` until `Done` flips it to `unknown`. A small case in `ProjectEvents` (or surfacing `run.failed` in the recent-events pane explicitly) would make a crashed background runner observable.
- NIT internal/tui/live.go:32-40 `LiveOptions.Status` (the workspace status) is plumbed through `app.go` but never consumed by the live model. Either use it in the header or drop the field; the implementation file pinkie-promised "no deviations" but ships an unread input.
- NIT internal/tui/live.go:184-189 `round=round-01` is hardcoded in the header. The slice is intentionally round-01-only per FINAL.md, but a magic literal will rot the moment a second round is supported. Consider deriving from `opts.Idea.Status` or threading the round label from the runner.
- NIT internal/app/app.go:226 `workspaceStatus.Ideas = []protocol.IdeaStatus{idea}` mutates the read snapshot only to feed an unused field — see the `LiveOptions.Status` NIT above.
- NIT internal/tui/live.go:160-162 The 250 ms `eventTickCmd` keeps re-arming after `doneMsg`. Quit eventually fires from the next `eventsMsg` branch, but the loop is more wasteful than it needs to be; gating the next tick on `m.done` would tighten the shutdown.

## Test gaps / residual risk
- No test exercises `agent.skipped` through `ProjectEvents`. The MAJOR finding above is uncovered; a unit test that feeds a single `agent.skipped` event (no preceding `agent.started`) and asserts a sane `Duration` (e.g. zero) would have caught it.
- No test asserts the spec's "unknown is sticky" property. `TestProjectEventsDerivesAgentAndRoundState` checks gemini's first event but never feeds gemini a follow-up `agent.finished` to verify the state stays `unknown`.
- FINAL.md asks for "render tests or view snapshots for representative terminal sizes" — `TestLiveViewIncludesRequiredPanels` only checks width=100. There is no narrow-terminal (<80) test to exercise the `bodyWidth < 76` clamp at live.go:179-182, no wide-terminal test, and no test that the ANSI-colored agent table actually aligns.
- No test covers the `doneMsg` path performing the final read-then-quit sequence promised in IMPLEMENTATION.md ("If the runner completes before the final `round.completed` … is read, the TUI performs one final event read before exiting").
- No test covers `tailLogFile` when the path is empty/missing (the common case for skipped or pending agents) — behaviour is correct (returns `""`) but the empty-paths path through `refreshLogPreview` (live.go:312-330) has no coverage.
- No test for `readEventsFromOffset` returning a JSON error and the model surfacing `errText` without losing offset progress already made on prior ticks.

## Summary
The slice mostly tracks FINAL.md: the async runner handle, default-vs-`--no-tui` paths, byte-offset polling, and the four-pane layout are all in place, and the partial-line/log-tail edge cases have focused unit coverage. Two real bugs need fixing before merge: skipped agents whose runner path didn't emit `agent.started` will display a nonsensical multi-year ELAPSED, and the `unknown` state is silently overwritten on the second event for any unmapped agent — directly contradicting the FINAL.md definition. The remaining items are smaller polish and test-gap concerns that block the "no deviations" claim in IMPLEMENTATION.md but not the slice's overall shape.
