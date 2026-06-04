---
agent: claude
idea: unified-tui-home
round: 1
date: 2026-06-04
---

## Summary

Make `internal/tui/live.go`'s `liveModel` the ONE model. Add a **Home tab** (first
tab) for the no-run state — lists ideas/recent runs, `N` starts an idea — and let
the model carry an **optional active run** so it transitions Home → agent tabs
when a run launches and back (Home stays a tab). Stop quitting on `done`. Repoint
`parley tui` to this model (Home mode) via a `StartRun` callback; retire
`tui.RunWorkspace` + the old `app.go` workspace model. `parley run "task"` opens
straight into agent tabs; `parley tui` opens to Home; both are the SAME model.

## Proposed approach

### One model, two states (no-run / active-run)
- `liveModel` gains: `home homeState` (ideas + recent runs list + selection),
  `start StartRunFunc` callback, and the run fields become **per-active-run** (a
  small `run` struct: RunID, RunDir, participants, Done, offset, events, state,
  buffers, questions). Empty run = no-run state.
- **Tab order:** `Home` (id `home`, always first) → agent tabs (only the active
  run's participants) → `Status` (when a run is active). No active run ⇒ just
  `Home` (+ nothing else). Default active tab = `Home` on `parley tui`; = first
  agent on `parley run`.
- `LiveOptions` gains: `Home bool` (open to Home), `Status protocol.WorkspaceStatus`
  (already present) for the ideas list, and `StartRun StartRunFunc` where
  `StartRunFunc(task string) (RunHandle, error)` launches **all available agents**
  through the existing runner path (the same one `runTask` uses) and returns the
  RunID/RunDir/Done/participants. The TUI never runs a parallel engine.

### Home tab (owner #1, #4, #5)
- Renders: open ideas (from `Status`/`protocol.ReadWorkspaceStatus`) + recent runs
  (`runstate.ListRuns`), selectable with `↑/↓`? — no, `↑/↓` are tab switches;
  use the input + a selection cursor via `j/k`? those type. Use `Home`-local
  selection on `tab`-within? Simplest: Home is mostly informational + the input
  row is the launcher. **`N`** (uppercase, not used elsewhere) opens new-idea
  mode: the input row becomes `new idea › <task>`; Enter → `m.start(task)` →
  switches to the new run's agent tabs.
- Selecting a recent run on Home (optional) could open it; v1 can defer that —
  `N`-launch is the must-have.

### Don't quit on done (owner #2/#3 persistence)
- On `doneMsg`: set `m.run.done = true` (status line shows "done"), **do not
  `tea.Quit`**. Keep ticking so the user can read transcripts, steer, or press
  `N` to start another idea (which replaces the active run). Only `ctrl+c`
  (cancel) and `/quit`/`esc`-on-empty-at-Home detach.

### Available-agents only (owner #1/#2)
- The default participant set (when none passed) = **all discovered `Found`
  agents** (`agents.Discover` → `.Found==true`), computed in `app.go`'s StartRun
  and runTask. The tab list is the run's participants, so unavailable agents
  never appear. (Change runTask's bounded 2–4 default to "all available".)

### Retire the old workspace TUI (owner #5)
- `parley tui` (runTUI) builds `LiveOptions{Home:true, Status, StartRun}` and
  calls `tui.RunLive`. Delete `tui.RunWorkspace` + the `app.go` model/panes;
  fold its ideas/runs listing into Home and its launch into `StartRun`. Keep
  `ActionRunner`? The old "next actions" execution is separate — defer/keep as a
  Home affordance only if cheap; otherwise drop for v1 (owner said discard the
  rest).

### Slices
1. **Run struct + no-run state**: extract the active-run fields into `m.run`;
   `RunDir==""` ⇒ no-run; tab order Home/agents/Status; default active tab.
2. **Stay open on done**: stop `tea.Quit` on `doneMsg`; status shows done; ticks
   continue.
3. **Home tab + `N` launch**: Home renders ideas/runs; `N` → new-idea input →
   `StartRun` → switch to agent tabs. Wire `StartRun` in `app.go` to launch all
   available agents via the runner path.
4. **Repoint + retire**: `parley tui` → `RunLive(Home)`; delete `RunWorkspace` +
   workspace model; default participants = all available in runTask/StartRun.
5. **Tests + transcript verify + polish**: state-transition tests, available-only
   tab test, a real-run transcript check, Home render test.

## Concerns / open questions

- **`StartRun` signature/threading**: launching agents returns a `Done` channel +
  RunID/RunDir; the model must swap to the new run and re-init the event tick for
  it. Need a clean `RunHandle` and a `tea.Cmd` that begins reading the new run's
  events.
- **Per-agent transcript not visible (#3)**: likely the buffer loads fine; verify
  the real run writes `stdout` paths into `agent.started` and that
  `ensureActiveBuffer`/`refreshBuffers` run for the active run. Possible bug: if
  tabs were built before any `agent.started`, the buffer path is empty until the
  event arrives — refreshBuffers must pick up the path lazily (it does, via
  `agentByID`).
- **Selection on Home without arrows** (arrows are tabs): use `N` to launch; run
  selection can be a later slice or via the input (`/open <run>`).
- **ActionRunner / next-actions**: keep or drop? Owner said discard the rest; I
  lean drop for v1, keep only ideas/runs listing + `N`.

## Risks

- **Big merge / regression**: deleting the workspace model + repointing `parley
  tui` risks breaking `parley tui` users; lock behavior with tests; keep
  `--no-tui` + `parley status`/`resume --no-tui` paths intact.
- **Run-swap state bugs**: switching the active run (N during an active run) must
  reset offset/buffers/questions cleanly; test it.
- **"All available agents" could be many**: tabs may overflow — the Slice-from
  tui-claude-cli-layout windowing handles it, but verify with 4–5 agents.
- **Launch failures inside the TUI**: a failed `StartRun` must surface in the
  status line, not crash; the model stays at Home.
