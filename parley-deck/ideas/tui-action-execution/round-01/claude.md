---
agent: claude
idea: tui-action-execution
round: 1
date: 2026-05-25
---

## Summary

The TUI shows next-actions from each run but they are inert: no selection cursor, no execution path, no feedback. The goal is to make the action list navigable and to let the operator trigger safe, non-interactive actions directly from the TUI, while clearly blocking actions that require a shell prompt or carry mutation risk the operator hasn't confirmed.

## Proposed approach

### 1. Selection state and focus mode

Add `focusActions` alongside the existing `focusIdeas`/`focusAgents`. Track `selectedAction int` on the model. When the action panel is focused, j/k move the selection cursor; the `>` marker in `renderQuestions` follows `selectedAction` instead of being hardcoded to index 0. Tab cycles through focus zones in the order: ideas → agents → actions (skip actions when the current run has none).

### 2. Injected action executor

Define a callback type matching the existing pattern:

```go
type ExecuteActionFunc func(runID string, action runaction.NextAction) tea.Cmd
```

Pass it into `WorkspaceOptions` from `runTUIViewWithDiscovery`, where `app.actionCommand` is already available. This keeps the `tui` package free of `app` imports. The callback returns a `tea.Cmd` that performs the work asynchronously and emits a result message (`actionResultMsg`) carrying success/error.

### 3. Action classification at trigger time

When the operator presses Enter on a selected action:

| Kind | Behaviour |
|---|---|
| `inspect` | Execute immediately (read-only, no risk). |
| `draft-consensus`, `request-signoffs`, `finalize` | Execute immediately; these are CLI-driven mutations with no interactive prompts. |
| `answer-question` | Block with inline message: "This action requires interactive input — run it from the CLI." Show the exact command so the operator can copy it. |
| `retry-agent` | Block with inline message: "No automation command available for retry-agent yet." |

Risk metadata from `runaction` (if present) can gate `finalize` behind a one-key confirmation (y/n) rendered inline, but the first slice should not introduce a full modal dialog.

### 4. Feedback cycle

- On Enter: replace the action line with a spinner or `[running...]` marker. Disable re-trigger of the same action.
- On `actionResultMsg` success: fire `refreshRunsCmd` to reload state. Show a transient status line ("action completed") that clears on next render cycle.
- On `actionResultMsg` error: show the error inline in the action area, styled distinctly (e.g. red). The action stays selected so the operator can retry or read the message.

### 5. Footer update

Add `enter` to the footer key hints when `focusActions` is active. Remove it (or grey it) when focus is elsewhere so the key map stays contextual.

### 6. Testing

- `app_test.go` (tui): test that j/k move `selectedAction`, that Enter on a safe action calls the injected executor, and that Enter on `answer-question`/`retry-agent` produces the blocking message without calling the executor.
- `app_test.go` (app layer): test that the callback wired through `runTUIViewWithDiscovery` correctly delegates to `actionCommand` and returns appropriate success/error.

## Concerns / open questions

1. **`finalize` confirmation UX.** A single y/n prompt inline is minimal but could be accidentally confirmed. Should the first slice skip `finalize` entirely and require CLI execution, or is inline confirmation acceptable?
2. **Action staleness.** If another agent advances the run between refreshes, the displayed actions may be stale. Executing a stale action will fail at the protocol layer, but the error message should be clear ("run state changed, refreshing...") rather than a raw CLI error.
3. **Multiple actions per run.** `NextActions` can return several actions. The current `renderQuestions` flattens them into a list. If two actions conflict (e.g. `finalize` and `retry-agent` on the same run), selection order matters — should the TUI warn, or trust the protocol to reject invalid sequences?
4. **Output capture.** `inspect` and other commands produce stdout. The TUI currently has no output pane. Options: (a) ignore output and tell the operator to check the CLI, (b) append output to a scrollable log region, (c) dump to a temp file and show the path. Option (a) is simplest for the first slice; (b) is better UX but adds layout complexity.

## Risks

- **Import cycle.** The injected-callback approach avoids it, but if someone later tries to call `app` functions directly from `tui`, the build breaks. The boundary should be documented with a comment on `WorkspaceOptions`.
- **Blocking the event loop.** If `ExecuteActionFunc` shells out synchronously instead of returning a `tea.Cmd`, the TUI freezes. The implementation must wrap the subprocess call in a `tea.Cmd` goroutine, same as `StartRunFunc` already does.
- **Partial execution.** If the operator quits the TUI while an action command is running, the subprocess may orphan. The first slice can accept this (CLI commands are short-lived), but a follow-up should wire context cancellation through the Bubble Tea shutdown path.
- **UX confusion from blocked actions.** Showing actions the operator cannot trigger may feel broken. The inline message must be specific enough ("run from CLI: `parley action answer-question --run X`") that the operator has a clear next step.
