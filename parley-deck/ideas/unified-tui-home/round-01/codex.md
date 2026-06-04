---
agent: codex
idea: unified-tui-home
round: 1
date: 2026-06-04
---

## Summary

Unify around `internal/tui/live.go`, not around the old dashboard. The current live model already owns the right low-level mechanics: `LiveOptions` carries status, idea, participants, run id/dir, done channel, cancel, and resume state (`internal/tui/live.go:49`); `Init` starts event, question, elapsed, and done ticks (`internal/tui/live.go:141`); the tabbed renderer uses projected agent state for tabs (`internal/tui/live.go:238`); and per-agent stdout is loaded and tailed through bounded focus buffers (`internal/tui/live.go:581`, `internal/tui/live.go:618`).

The main work is to make that model nullable on the run axis: `Home` exists when there is no active run, and `Home + agent tabs + Status` exists after a run is active. `parley tui` should open that same model in no-run mode. `parley run` and `parley resume` should open it with an active run. `doneMsg` must stop being a quit trigger: today `eventsMsg` exits once `m.done` is true (`internal/tui/live.go:176`), which directly conflicts with the locked requirement to keep the TUI open.

## Proposed approach

Define the app-facing launch contract in `internal/tui/live.go` as a replacement for the old `StartRunFunc` in `internal/tui/app.go:37`:

```go
type StartRequest struct {
    Task string
    Auto bool
}

type StartedRun struct {
    Idea         protocol.IdeaStatus
    Participants []string
    RunID        string
    RunDir       string
    Handle       *runner.Handle
    Cancel       context.CancelFunc
}

type StartRunFunc func(context.Context, StartRequest) (StartedRun, error)
```

This makes the TUI depend on the runner handle shape intentionally. The handle already exposes `RunID`, `RunDir`, `Done()`, `Wait()`, and copied `Results()` (`internal/runner/runner.go:62`, `internal/runner/runner.go:100`, `internal/runner/runner.go:104`). The trade-off is a slightly stronger `tui -> runner` dependency, but it prevents inventing a second async abstraction and keeps the owner decision "no parallel engine" honest.

Move the existing `runTask` launch path into a shared app helper:

```go
func startRoundOne(ctx context.Context, root, task string, discovered []agents.Discovery, auto bool) (tui.StartedRun, error)
```

Inside it, use all found agents as participants, call `runcontrol.Create`, start the auto-answerer when requested, then call `runner.RunRoundOneAsync`. This is exactly the path `runTask` uses today: discovery and participant selection (`internal/app/app.go:1612`, `internal/app/app.go:1619`), `runcontrol.Create` (`internal/app/app.go:1632`), optional `StartAutoAnswerer` (`internal/app/app.go:1667`), and `runner.RunRoundOneAsync` (`internal/app/app.go:1670`). The current TUI start callback already creates through `runcontrol.Create`, but it discards the async handle (`internal/app/app.go:1971`, `internal/app/app.go:1990`); that must change so the model can wait on `handle.Done()` and swap into the new run immediately.

On `N`, the model should enter a compact Home input state, then call `startRunCmd(m.startRun, request)`. The returned `StartedRun` should be applied by a single method:

```go
func (m *liveModel) activateRun(started StartedRun) tea.Cmd
```

`activateRun` resets `opts.Idea`, `opts.Participants`, `opts.RunID`, `opts.RunDir`, `opts.Done`, `opts.Cancel`, `offset`, `events`, `questions`, `done`, `activeTab`, and `buffers`, then recomputes `state = ProjectEvents(started.Participants, nil, m.now)`. It must return a `tea.Batch` with `readEventsCmd`, `readQuestionsCmd`, `eventTickCmd`, `elapsedTickCmd`, and `waitDoneCmd(started.Handle.Done())`, mirroring current `Init` (`internal/tui/live.go:141`). This lets switching runs happen in place without reconstructing the Bubble Tea program.

Model lifecycle should be explicit:

- No-run state: `opts.RunID == ""` and `opts.RunDir == ""`; tabs are `[Home]` only, plus optionally `Status` if we want workspace diagnostics. No event/question ticks should read an empty path.
- Active-run state: tabs are `[Home] + agent tabs + [Status]`. `tabIDs` currently returns only agents plus Status (`internal/tui/live.go:498`); prepend Home and keep `activeTabResolved` defaulting to the first running agent after activation (`internal/tui/live.go:506`).
- Done state: `doneMsg` sets `m.done = true` and performs one final event read (`internal/tui/live.go:199`), but `eventsMsg` must not return `tea.Quit`. It should refresh Home's run list and keep the active run visible.
- Switching runs: selecting a recent run on Home should call a lightweight `activateExistingRun(run runstate.RunSummary)` that reuses the same reset logic with no handle/done channel. That matches `parley resume`, which currently loads a run, reads status, derives the idea, and calls `RunLive` with `Resume: true` (`internal/app/app.go:981`, `internal/app/app.go:998`, `internal/app/app.go:1003`).

Default participants should be every discovered `Found` agent. `agents.Discovery.Found` is set by PATH lookup in discovery (`internal/agents/discover.go:249`, `internal/agents/discover.go:257`). The old workspace model already filters to found agents for display (`internal/tui/app.go:157`). `selectedParticipantIDs` should stop excluding Gemini by default: `installedAgentIDs` currently skips `result.ID == "gemini"` (`internal/app/app.go:2061`), and the test locks that old policy in (`internal/app/app_test.go:73`). The new owner decision supersedes that bounded default. If a deprecated built-in is installed and discoverable, it is available; if we still want to discourage it, show a note from `Discovery.Notes`, not a hidden exclusion.

The tab list must come from run participants, not global discovery. `runcontrol.Create` records participants in the `run.created` event and manifest (`internal/runcontrol/runcontrol.go:52`, `internal/runcontrol/runcontrol.go:65`), and `runstate.LoadRun` reads them back (`internal/runstate/runstate.go:121`). `runstate.ProjectEvents` initializes agent state from the participant list (`internal/runstate/runstate.go:321`). Therefore unavailable agents never appear if participant selection is based only on `Found` discoveries.

Repoint `parley tui` by changing `runTUIViewWithDiscovery` to call the unified `RunLive` with Home options instead of `RunWorkspace`. Keep the status and runs reads (`internal/app/app.go:1938`, `internal/app/app.go:1951`) because Home needs them, keep session registration (`internal/app/app.go:1956`), and pass a `RefreshRuns` callback backed by `runstate.ListRuns` (`internal/runstate/runstate.go:261`). Delete `RunWorkspace`, `WorkspaceOptions`, the old `model`, `ActionRunner`, and the old actions panes from `internal/tui/app.go` after Home has ideas/runs listing and `N` start. Preserve `RunInit`, or move it to a small setup file, because `runTUIViewWithDiscovery` still uses it when the workspace is missing (`internal/app/app.go:1938`, `internal/app/app.go:1941`).

Home should fold in only:

- Ideas from `protocol.ReadWorkspaceStatus`, whose status contains root, transport, and ideas (`internal/protocol/workspace.go:22`, `internal/protocol/workspace.go:104`).
- Recent runs from `runstate.ListRuns`, already sorted newest first (`internal/runstate/runstate.go:261`, `internal/runstate/runstate.go:288`).
- Basic open question and agent progress summaries from `RunSummary`, not the old action runner. `runStatus` already combines workspace status and runs for CLI output without changing runner contracts (`internal/app/app.go:668`).

Keep `--no-tui`, `parley status`, and `parley resume --no-tui` intact by leaving their non-TUI branches untouched. In particular, `runTask --no-tui` should continue to call synchronous `runner.RunRoundOne` and print results (`internal/app/app.go:1647`, `internal/app/app.go:1653`), and `runResume --no-tui` should still print the resolved run detail (`internal/app/app.go:993`).

Incremental slices:

1. Add Home/no-run state and tests for tab lists, no empty event reads, and `doneMsg` not quitting.
2. Replace default participant selection with all `Found` discoveries and update tests.
3. Add `StartRunFunc`/`StartedRun`, wire `N` through shared `startRoundOne`, and activate the returned handle in-place.
4. Repoint `runTUI` to unified `RunLive(Home)` while leaving `RunWorkspace` temporarily unused.
5. Move or keep `RunInit`, then delete old workspace model/action code once tests and a real run prove transcripts populate.

## Concerns / open questions

The current `RunLive` API returns only when the Bubble Tea app exits (`internal/tui/live.go:121`). If `parley run "task"` now stays open after done, `runTask` cannot print `handle.Wait()` results until the user quits. That is acceptable for TUI mode, but the app should still reap results asynchronously on `doneMsg` so failures can be surfaced in Status/Home without waiting for process exit.

The old dashboard allowed session-only launch mode overrides with `h/i/a/m/x` (`internal/tui/app.go:238`). The owner said default all available and retire the dashboard, so I would drop these controls from the first unified slice. If launch overrides matter later, they should return as a small Home agent settings panel, not as a preserved old focus pane.

Home run switching needs a policy for live external processes. A TUI-started active run should keep running when the user switches Home to inspect an older run; `ctrl+c` should cancel only the currently attached TUI-started run if it has a `Cancel`. This is a behavior change from the old footer warning that quitting cancels TUI-started runs (`internal/tui/app.go:889`), but it matches the locked "do not quit on done" direction better.

## Risks

There is a real risk of duplicate ticks if activation starts new `eventTickCmd` loops without invalidating old ones. Use a `runToken` or compare `msg.runID` in event/question/done messages so stale ticks from the previous run are ignored after switching runs.

Per-agent transcripts depend on stdout paths being populated by agent events. `runAgent` writes stdout to `runs/<run>/agents/<agent>/stdout.log` (`internal/runner/runner.go:277`, `internal/runner/runner.go:336`), and the buffer path comes from `AgentState.StdoutPath` (`internal/runstate/runstate.go:49`, `internal/tui/live.go:591`). A real multi-agent run must be part of verification because a missing `stdout_path` in event projection would render only "no output yet" (`internal/tui/live.go:317`).

Deleting `ActionRunner` removes in-TUI consensus action execution. That is deliberate for the unified first pass: slash/input steering and HITL answering already exist in live mode (`internal/tui/live.go:707`), while consensus actions remain available through CLI commands and `parley status`. Re-adding actions should be a later Home workflow if users miss it.
