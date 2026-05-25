---
agent: codex
idea: tui-action-execution
round: 1
date: 2026-05-25
---

## Summary

The first operational TUI action slice should add action focus/selection and an
injected asynchronous action runner, while keeping all CLI orchestration in
`internal/app`. This avoids an import cycle, keeps `internal/tui` testable, and
lets the app layer decide which planner actions are executable versus advisory.

## Proposed approach

Add an `ActionRunner` option to `tui.WorkspaceOptions`:

- `ActionRequest` carries the selected `runstate.RunSummary` and
  `runaction.NextAction`.
- `ActionResult` carries a short message, an optional command string, and a
  `Refresh` flag.
- Bubble Tea starts the runner via a `tea.Cmd`, sets an `actionRunning` flag, and
  handles an `actionRunMsg` when the command completes.

Extend the dashboard model with `focusActions` and `selectedAction`. `tab` should
cycle through sessions, actions, and agents when those panes have useful content.
`j/k` should move the action selection while action focus is active. The action
pane should continue showing the top planner actions even when no runner is
configured, but pressing `enter` should return a clear message instead of doing
nothing.

Move the command-string construction out of `internal/app.actionCommand` into a
small shared helper, preferably in `internal/runaction`, so both CLI printing and
the app-provided TUI runner use the same mapping. Keep this helper limited to
planner action metadata; do not make `internal/tui` import `internal/app`.

In `internal/app.runTUIViewWithDiscovery`, pass an action runner that dispatches
supported actions through existing app command handlers or equivalent internal
functions. For this first slice:

- `inspect` is low risk and can return the existing status command as guidance;
  direct navigation is not necessary yet.
- `draft-consensus`, `request-signoffs`, and `finalize` can execute via the app
  layer with `--dir <root>` and should request a refresh on success.
- `answer-question` should remain advisory because it needs user text input.
- `retry-agent` should remain unsupported until the runtime has a real retry
  primitive.

Add focused tests in `internal/tui/app_test.go` for action focus cycling,
selection, enter-triggered runner invocation, busy-state blocking, result
messaging, and refresh-on-success. Add app/helper tests that verify command
mapping still avoids hardcoded agents and preserves round-aware consensus draft
commands.

## Concerns / open questions

The main UX question is whether pressing `enter` should execute normal/high-risk
actions immediately. My recommendation for this slice is conservative but still
useful: if the planner marks an action as `RequiresYes`, the app runner may only
execute it when the underlying CLI command already has an explicit noninteractive
confirmation flag, such as `consensus request-signoffs --yes`. Otherwise the TUI
should show the command and explain that direct execution is not implemented.

The second concern is output capture. Calling the top-level `Run` function with
buffers is acceptable for the first slice and prevents duplicated app logic, but
the result message must stay short and must not dump full agent logs into the
dashboard footer.

## Risks

- Running long signoff requests from the TUI can take time; it must be fully
  asynchronous and visibly busy.
- Reusing the top-level CLI dispatcher from inside the app layer can accidentally
  recurse into the TUI if the command map is too broad. Only consensus/status
  subcommands should be allowed.
- If command construction remains duplicated between CLI output and TUI runner,
  planner actions will drift again.
