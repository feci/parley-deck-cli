---
idea: unified-tui-home
status: complete
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
- [x] **Slice 5 — delete the old workspace model + tests** (done in fix-up
      cycle 1, AF5): `internal/tui/app.go` reduced to the init wizard + shared
      styles/helpers; the workspace `model`, `RunWorkspace`/`Run`/
      `WorkspaceOptions`, `Start*`/`Action*` launch types, workspace msg types,
      and the app-side `runTUIAction`/`consensusActionArgs`/`commandOutput`/
      `applySessionLaunchOverrides` all deleted, with their tests. `tui/app.go`
      1443 → 150 lines.
- [x] Checks: `go build ./...`, `go vet ./...`, `go test ./...` — all green.

## Fix-up cycle 1 (Phase 8)

Applied the Phase 7 review-consensus agreed fixes (see `review/consensus.md`):

- **AF1** — TUI detach no longer cancels `N`-launched runs: `newLaunchFunc` reaps
  each handle in a background goroutine (`Wait` → `registerWorkspaceSessions`);
  the unconditional cancel-all-on-exit defer was removed. Only the attached run's
  `Cancel` (ctrl+c) cancels.
- **AF2** — `parley run` no longer owns secondary launches: `Start` is not passed
  into `parley run`'s live TUI; `Root` is kept so Home still lists ideas/runs and
  `N` reports that new ideas start from `parley tui`.
- **AF3** — added `TestTranscriptPopulatesFromOnDiskRun`: writes a run dir with
  `events.jsonl` (`agent.started` → stdout path) + `stdout.log`, drives an
  `eventsMsg` read, and asserts the agent buffer is non-empty and the view shows
  the output (closes the owner's #3 transcript gate end-to-end).
- **AF4** — done-state exit hint: status line shows `[done]` and the input row
  shows the `/quit or esc to exit` hint so a finished run never feels stuck.
- **AF5** — retired the workspace model (see Slice 5 above and Deviations).
- **AF6** — help wording is now `ctrl+c cancel the attached run, else quit`.

Checks after fix-up: `go build ./...`, `go vet ./...`, `go test ./...` — green.

## Fix-up cycle 2 (Phase 8)

codex's round-02 re-review accepted AF2–AF6 but raised a valid MAJOR on AF1:
removing the explicit defer-cancel-all was not sufficient because `N`-launched
runs derive from the top-level signal context, whose `defer cancel()` fires when
`parley tui` returns — so a normal `/quit` still canceled in-flight launched runs
via context propagation.

- **AF1 (deepened)** — added `launchReaper` in `internal/app/app.go`: each
  launched run is tracked (`reaper.track`) and, after `RunLive` returns,
  `runTUIViewWithDiscovery` calls `reaper.waitForActive(stdout)` to wait for any
  still-running launched runs to finish (and record their sessions) **before** the
  command returns and the parent cancel fires. The parent context stays live
  during the wait, so a real ctrl+c (SIGINT, once the TUI releases the terminal)
  still aborts; the attached run's in-TUI `Cancel` (ctrl+c) is unchanged. New
  test `TestLaunchReaperWaitsForInFlightRuns` asserts the wait blocks until the
  in-flight run finishes (detach waits, never abandons).

Checks after cycle 2: `go build ./...`, `go vet ./...`, `go test ./...` — green.

## Tests

`internal/tui/live_test.go`: `TestHomeDefaultWhenNoRunAndTabOrder`,
`TestActiveRunTabOrderHomeFirst`, `TestDoneDoesNotQuit`,
`TestNLaunchesViaStartFunc`, `TestRunTokenDropsStaleEvents`. Existing tabbed-layout
tests still pass (Home prepended to tab order without breaking defaults). Existing
`TestTabbedDefaultShowsTranscriptAndInput` already proves a transcript renders.

## Deviations from FINAL.md

- **The old workspace model is fully retired and deleted (fix-up cycle 1, AF5).**
  During slice 4 the workspace model was retired *functionally* (`parley tui`
  launches `RunLive(Home)`) but left in place as dead-but-compiling code so the
  reviewers could pinpoint exactly what to delete. In fix-up cycle 1 it was
  removed: `internal/tui/app.go` now holds only the init wizard and the shared
  styles/helpers `live.go` needs (`valueOr`, `truncateText`, the styles), and the
  app-side `runTUIAction`/`consensusActionArgs`/`commandOutput`/
  `applySessionLaunchOverrides` are gone. The `Launch*` launch types (renamed
  during slice 2 to avoid colliding with the old `StartRequest`/`StartRunFunc`)
  are now the only launch types and keep their names.
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

## Phase 8 complete (2026-06-04)

Fix-up cycle 2 re-review: codex and hermes both signed ACCEPT (see
`review/consensus.md` "Fix-up cycle 2" + `review/round-03/`). Zero remaining
agreed fixes. Shipped as parley-deck-cli **1.14.0**.
