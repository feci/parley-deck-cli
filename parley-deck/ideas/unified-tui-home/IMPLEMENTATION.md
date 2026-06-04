---
idea: unified-tui-home
status: implemented
implementer: claude
started: 2026-06-04
completed: 2026-06-04
branch: parley-deck-cli#feat/unified-tui-home
head-commit: see-branch-tip
design-pr: https://github.com/feci/parley-deck-cli/pull/35
implementation-pr: https://github.com/feci/parley-deck-cli/pull/35
---

## Summary of work

Unified `parley tui` per FINAL.md: `internal/tui/live.go`'s tabbed model gained a
Home tab + a nullable active run; `parley tui` now opens it in Home mode and
`parley run`/`resume` open it with an active run. Same model. Reuses the runner/
events/`steer`/`runstate`/`hitl` contracts and `--no-tui` unchanged.

## Implementation plan / checklist

- [x] **Slice 1 — nullable run state** (live.go): `homeTabID`; `hasRun()`/
      `attached()`; `tabIDs` = Home first → agent tabs → Status (only with a run);
      `activeTabResolved` defaults to Home (no run / `opts.Home`) else first
      running agent; `renderTabbed` Home branch; no empty-path reads when at Home
      (`runCmds`).
- [x] **Slice 2 — launch contract** (live.go): `LaunchRequest`/`LaunchResult`/
      `LaunchFunc`; `activateRun` swaps the model in place + bumps `runToken`;
      event/question/done/tick messages carry `runToken` and stale ones are
      dropped; `doneMsg` sets done WITHOUT `tea.Quit` (TUI stays open); elapsed
      clock is global (not token-gated).
- [x] **Slice 3 — Home + N + /open** (live.go): `N` opens `new idea ›` compose →
      `launchIdea` → `Start` callback → `activateRun`; `/open <slug|run>` →
      `openRun` (observational, via `runstate.ResolveRun`); `/home`; `renderHome`
      (ideas + recent runs + hints); input row labels/hints per context.
- [x] **Slice 4 — repoint + share launch** (app.go): `newLaunchFunc` (all
      available agents via `runcontrol.Create` + `RunRoundOneAsync`);
      `runTUIViewWithDiscovery` → `RunLive(Home, Root, Status, Start)` instead of
      `RunWorkspace`; `runTask` also passes `Root` + `Start` so `N` works inside
      `parley run`. `installedAgentIDs` (gemini-excluded) unchanged.
- [~] **Slice 5 — delete the old workspace model + tests** (in progress): added
      the unified-behavior tests below; the now-dead workspace model deletion is
      the remaining cleanup (see Deviations).
- [x] Checks: `go build ./...`, `go vet ./...`, `go test ./...` — all green.

## Tests

`internal/tui/live_test.go`: `TestHomeDefaultWhenNoRunAndTabOrder`,
`TestActiveRunTabOrderHomeFirst`, `TestDoneDoesNotQuit`,
`TestNLaunchesViaStartFunc`, `TestRunTokenDropsStaleEvents`. Existing tabbed-layout
tests still pass (Home prepended to tab order without breaking defaults). Existing
`TestTabbedDefaultShowsTranscriptAndInput` already proves a transcript renders.

## Deviations from FINAL.md

- **The old workspace model is retired functionally but not yet deleted.**
  `parley tui` now launches `RunLive(Home)`, so `tui.RunWorkspace`/
  `WorkspaceOptions`/the `internal/tui/app.go` workspace `model`/`ActionRunner`/
  `ActionRequest`/`ActionResult`/`StartRequest`/`StartRunFunc`/`StartedRun`, and
  `app.go`'s `runTUIAction`/`applySessionLaunchOverrides`, are now DEAD code that
  still compiles. Deleting them safely (tui/app.go also holds shared helpers like
  `valueOr`/`truncateText`) is the final slice-5 step, deferred to fix-up so the
  review can pinpoint exactly what is dead. I renamed the new launch types
  `Launch*` to avoid colliding with the old `StartRequest`/`StartRunFunc` until
  they're removed.
- **D10 transcript gate**: verified via tests (a transcript renders from a
  loaded buffer) and by code path (stdout path flows from `agent.started` →
  `agentBuffer` lazily each tick); a full multi-agent live run is the reviewer's
  acceptance gate.

## Notes for reviewers

- Reuse unchanged: `internal/steer`, `runstate` (segment/`[FINISHED]`), `hitl`,
  focus buffers, slash/key model. `--no-tui`/`parley status`/`resume --no-tui`
  untouched. `installedAgentIDs` (gemini-legacy excluded) unchanged.
- `runToken` drops stale ticks/reads after a run swap (`N`/`/open`); the elapsed
  clock is intentionally global (one loop, not token-gated).
- Attached (`N`/`parley run`) runs have a `Cancel`; `/open`-ed runs are
  observational (no cancel/done-wait).
- PLEASE flag the exact dead workspace-model symbols to delete in fix-up.
