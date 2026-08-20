---
idea: tui-action-execution
finalized-by: codex
date: 2026-05-25
status: final
---

## Decision

Implement a narrow operational action slice for the dashboard TUI:

- users can focus the action area, select a planner action, and press `enter`;
- `internal/tui` remains app-agnostic through an injected asynchronous runner;
- `internal/app` owns action execution and command dispatch;
- command construction is shared through `internal/runaction`;
- unsupported or interactive actions produce clear in-TUI feedback instead of
  silent no-ops.

## Implementation plan

1. Add shared command construction under `internal/runaction`.
   - Build command strings from `NextAction` plus fallback run/idea IDs.
   - Preserve the current mappings for `answer-question`, `draft-consensus`,
     `request-signoffs`, `finalize`, and `inspect`.
   - Keep `retry-agent` unmapped until a real runtime retry command exists.

2. Extend `internal/tui`.
   - Add `ActionRequest`, `ActionResult`, and `ActionRunner`.
   - Add `ActionRunner` to `WorkspaceOptions`.
   - Add model state for `focusActions`, `selectedAction`, `actionRunning`, and
     short action feedback text.
   - Let `tab` include actions when the selected run has planner actions.
   - Let `j/k` move selected actions while action focus is active.
   - Let `enter` invoke the injected runner asynchronously.
   - Refresh runs after successful mutating actions when requested by the
     runner.

3. Wire app execution in `internal/app`.
   - `runTUIViewWithDiscovery` passes an app-backed runner into the TUI.
   - The runner supports direct app-layer execution for `draft-consensus`,
     `request-signoffs`, and `finalize`.
   - `request-signoffs` must use the existing explicit `--yes` path.
   - `inspect`, `answer-question`, and `retry-agent` return advisory messages
     and the relevant manual command when available.

4. Add tests.
   - Shared command mapping tests cover round-aware draft commands and avoid
     hardcoded agents.
   - TUI tests cover focus cycling, action selection, enter-triggered runner
     invocation, unsupported action feedback, busy-state blocking, and
     refresh-on-success.
   - App tests cover the app-backed runner's dispatch behavior without needing a
     full interactive TUI.

## Deferred follow-ups

- A confirmation modal for `RequiresYes` actions.
- Runtime support for `retry-agent`.
- Native input flow for `answer-question`.
- Live action stdout/log panes and cancellation/recovery controls.
