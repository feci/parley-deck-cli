---
agent: codex
idea: unified-tui-home
review-round: 1
date: 2026-06-04
---

## Summary

The `live.go` core mostly matches FINAL.md: Home is first, Status only appears with a run, `N` is uppercase-only, `/open` uses `runstate.ResolveRun`, `activateRun` resets events/questions/buffers/offset and bumps `runToken`, and stale event/question/done/event-tick messages are dropped after a swap. The new run gets an immediate first read from `activateRun` -> `runCmds`, and the old event tick loop dies on token mismatch; I did not find a duplicate tick loop in the model itself. The elapsed ticker is intentionally global and not token-gated.

D10 also looks correct by code path: `agent.started` stores `stdout` in `runstate.ProjectEvents`/`applyAgentEvent`, `live.go` projects that into `AgentState.StdoutPath`, then `ensureBuffer`/`refreshBuffers` lazy-load the active or visited tab. A real multi-agent run should show output once each `agent.started` event lands. I did not run Go tests because this review was constrained to create exactly one file, and test caches would create additional workspace artifacts.

## Findings

### [MAJOR] TUI detach cancels `N`-launched runs from `parley tui`

What: `live.go` only calls `Cancel` on `ctrl+c`, while `esc` and `/quit` just return `tea.Quit`, which matches the attached-vs-observational contract. But `runTUIViewWithDiscovery` tracks every `newLaunchFunc` cancel and unconditionally calls all of them in a defer after `RunLive` returns (`internal/app/app.go:1959-1968`, `internal/app/app.go:1977-1982`). That means an attached run started with `N` from Home is canceled even when the user exits via `/quit` or empty-input `esc`, despite FINAL.md saying those paths detach and only `ctrl+c` cancels the attached run.

Why: This turns normal TUI exit into cancellation, so a user cannot safely start a long run from Home and leave the UI open-ended; the process receives cancellation even though they did not press `ctrl+c`.

Concrete fix: replace the unconditional cancel-on-return tracking with explicit ownership semantics. Let `LiveOptions.Cancel` remain the only cancel path for the active attached run, and have the app-side launch manager reap `handle.Wait()` in a goroutine and register the completed run/session. If shutdown cleanup is still needed, separate "cancel because user pressed ctrl+c" from "detach because the Bubble Tea program exited".

### [MAJOR] `parley run` does not own secondary runs launched with `N`

What: `runTask` passes `Start: newLaunchFunc(ctx, *root, discovered, nil)` into the live TUI (`internal/app/app.go:1672-1682`). If the user launches a new run with `N`, `newLaunchFunc` starts `runner.RunRoundOneAsync` and returns its `Done`/`Cancel` (`internal/app/app.go:1992-2018`), so the TUI swaps correctly. After `RunLive` exits, however, `runTask` still waits and reports only the original `handle` (`internal/app/app.go:1689-1694`). The secondary handle is not waited, results are not printed, and its session is not registered from this path.

Why: This is a lifecycle mismatch with FINAL.md's "real runner.Handle" launch contract. The UI can display the new active run, but the command driver still owns only the old run; on process return the top-level context can also cancel the secondary run without a result reap.

Concrete fix: use the same app-side launch manager for `parley run` and `parley tui`: every TUI-launched `LaunchResult` should carry or register its `runner.Handle`, get reaped exactly once, update session state, and be canceled only by the active-run `Cancel` path. If that cannot be done in this fix-up, disable `Start` in `parley run` until the secondary-run ownership is implemented.

### [MINOR] Dead workspace model still needs deletion

What: FINAL.md requires retiring the old workspace TUI. It is functionally bypassed, but the old code and tests still compile. Exact dead symbols to delete in fix-up: `internal/tui/app.go` `Run`, `RunWorkspace`, `WorkspaceOptions`, `model`, `focusZone`, `focusIdeas`, `focusActions`, `focusAgents`, `refreshRunsMsg`, `startRunMsg`, `actionRunMsg`, `refreshTickMsg`, `StartRequest`, `StartRunFunc`, `ActionRequest`, `ActionResult`, `ActionRunner`, `newModel`, `startRunCmd`, `actionRunCmd`, `upsertRunSummary`, and the old model methods/panes that only serve that dashboard. In `internal/app/app.go`, delete `runTUIAction`, `consensusActionArgs`, and `applySessionLaunchOverrides`; also remove their now-dead tests in `internal/tui/app_test.go` and the TUI-action-specific tests in `internal/app/app_test.go`. There is no current `StartedRun` symbol to delete; the live replacement is `LaunchResult`.

Why: Leaving a complete unreachable TUI surface behind makes future changes ambiguous and keeps obsolete `ActionRunner`/launch override semantics alive in tests.

Concrete fix: split `RunInit` plus shared helpers into a retained file, then delete the workspace dashboard code and stale tests. Keep shared helpers used by `live.go`, including `valueOr`, `truncateText`, `clipLines`, `tuiWidth`, `tuiHeight`, `clampInt`, `stripANSI`, `sectionTitle`, and the shared styles (`headerStyle`, `boxStyle`, `mutedStyle`, `okStyle`, `warnStyle`).

### [NIT] Help text overstates cancellation on Home and observational runs

What: The help overlay always says `ctrl+c             cancel the run` (`internal/tui/live.go:1145-1173`). In Home/no-run mode there is no run, and for `/open` observational runs `Cancel` is nil, so `ctrl+c` just quits.

Why: The behavior is correct, but the text contradicts the attached-vs-observational contract and can make users think `/open` can kill a live external process.

Concrete fix: change the wording to `ctrl+c             cancel attached run / quit`, or render a conditional line based on `m.attached()`.

## Open questions

None.
