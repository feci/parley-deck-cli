---
idea: tui-action-execution
drafted-by: codex
date: 2026-05-25
---

## Agreed decisions

- Add operational action focus to the dashboard TUI:
  - introduce `focusActions` and `selectedAction`;
  - render the selected planner action instead of always marking index 0;
  - let `tab` include the action area when the selected run has actions;
  - let `j/k` move action selection while the action area is focused;
  - let `enter` trigger the selected action while action focus is active.
- Keep `internal/tui` reusable and free of app imports:
  - add an injected asynchronous action runner to `WorkspaceOptions`;
  - define request/result types in `internal/tui` using `runstate.RunSummary`
    and `runaction.NextAction`;
  - execute actions through a Bubble Tea command and handle a result message,
    so rendering does not block during long-running work.
- Keep CLI orchestration in `internal/app`:
  - `runTUIViewWithDiscovery` wires the app-backed action runner into the TUI;
  - the runner may call existing CLI/app command handlers with captured output;
  - the TUI receives only concise success/error/manual-command feedback.
- Move planner action command construction out of `internal/app.actionCommand`
  into a shared helper under `internal/runaction`, so CLI status output and TUI
  execution use one mapping.
- First-slice executable policy:
  - `draft-consensus`, `request-signoffs`, and `finalize` are supported direct
    app-layer actions when a command can be constructed from the selected
    `NextAction`;
  - `request-signoffs` must keep using the existing explicit `--yes` command
    path;
  - `inspect` remains advisory and shows the status command rather than adding a
    new output pane;
  - `answer-question` remains advisory because it needs user-provided text;
  - `retry-agent` remains unsupported because the runtime has no concrete retry
    primitive yet.
- After a successful mutating action, request a run refresh and clamp action
  selection against the refreshed action list.
- Add focused tests for action focus, selection, runner invocation, unsupported
  actions, busy-state behavior, result messages, refresh-on-success, and command
  mapping.

## Agreed trade-offs

- Pressing `enter` on an action is the first-slice confirmation. A richer modal
  confirmation flow can be added later if normal-risk direct actions prove too
  easy to trigger accidentally.
- The TUI will not add a live stdout/log pane in this slice. Runner output is
  captured and reduced to a short status message to keep the dashboard layout
  stable.
- The runner will not implement pre-execution re-planning in this slice. The
  underlying CLI command remains the source of truth and must return a clear
  error if the action is stale.
- The action runner may reuse the top-level app command dispatcher for narrow
  consensus subcommands, but it must not dispatch recursively into `tui`.

## Open items deferred to implementation

- Decide the exact helper API shape in `internal/runaction` while implementing;
  it should be small and avoid importing higher-level app packages.
- Decide the exact footer wording and status line placement while preserving the
  existing dense dashboard layout.
- Consider follow-up slices for retry-agent runtime support, explicit
  confirmation modal, live action logs, and cancel/recovery behavior for
  long-running actions.

## Signoffs

### Signoff: codex — 2026-05-25
Status: ACCEPT

This consensus matches my proposed implementation slice and incorporates the
other participants' safety concerns by keeping retry-agent, answer-question,
live logs, and confirmation modals out of the first implementation.

### Signoff: claude — 2026-05-25
Status: ACCEPT

Consensus preserves the injected-callback architecture and async execution model from my proposal while sensibly narrowing the first slice — making inspect advisory and extracting command construction into a shared runaction helper are good refinements.

### Signoff: gemini — 2026-05-25
Status: ACCEPT

This consensus document aligns with my understanding and proposed solutions, providing a clear path forward for implementation.

### Signoff: hermes — 2026-05-25
Status: ACCEPT

Consensus safely scopes first slice to direct actions while preserving decoupling and async execution.