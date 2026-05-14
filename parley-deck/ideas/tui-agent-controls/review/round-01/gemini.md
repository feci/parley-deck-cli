---
agent: gemini
idea: tui-agent-controls
review-round: 1
date: 2026-05-14
reviewed-commit: f4315d4
---

# Summary

The implementation of `tui-agent-controls` successfully introduces interactive navigation and detailed agent inspection to the Parley Deck TUI. It fulfills the core requirements of the `FINAL.md` specification, including focus switching, row selection, and session-local launch mode overrides. The transition from initialization to the dashboard is handled correctly, ensuring keybindings work immediately. However, there are some minor inconsistencies in the detail rendering and some redundant logic in the TUI model that should be addressed for better maintainability and UX.

# Findings

### 1. Incomplete headless command rendering [MINOR]
In `renderAgentDetails`, the headless command shape is only rendered if `HeadlessMode` or `HeadlessArgs` is explicitly set. If an agent is configured using only the base `Commands` (common for simple headless agents), the "headless" line is omitted entirely from the detail view. This contrasts with the `interactive` line which uses `agents.InteractiveCommandOrDefault` to ensure a value is always shown.
**Suggested fix:** Use a similar helper or logic for headless to show the base command when mode/args are empty.

### 2. Redundant and inconsistent logic in `initModel.View` [MINOR]
The `initModel.View` method contains a branch `if m.status != nil` that renders a dashboard view by manually constructing a `model` struct. This branch is likely redundant since `Update` already transitions to the real `model` upon success. More importantly, this manual construction bypasses the `newModel` constructor, skipping critical initialization steps like focus setting and selection clamping.
**Suggested fix:** Remove the `if m.status != nil` branch in `initModel.View` or ensure it uses `newModel`.

### 3. Confusing default value for `ExternalBackend` [NIT]
The `valueOrDefault` helper returns `"cli-default"` for empty strings. While appropriate for model names, it results in `backend: cli-default` when the backend is unknown or not specified. Since "cli-default" is not a valid backend type (unlike "hosted" or "local"), this is confusing.
**Suggested fix:** Provide a specific fallback for backend (e.g., `"unknown"`) or use `agents.ExternalUnknown`.

### 4. `shift+tab` does not navigate backwards [NIT]
The `tab` and `shift+tab` keys both call `switchFocus`, which simply toggles between the two panes. While functional for exactly two panes, it violates the conventional expectation that `shift+tab` moves the focus in reverse.
**Suggested fix:** Update `switchFocus` to take a direction or implement a separate `prevFocus` method.

### 5. Idea selection has no associated action [NIT]
Users can switch focus to the Protocol pane and select ideas, but unlike agents, there is no detail panel or footer action for the selected idea. This makes the "selection" state for ideas feel slightly "hollow" or like a dead end.
**Suggested fix:** Add a simple detail view for ideas (e.g., showing the path or full participant list) or a note in the footer.

# Open questions

- Are there plans to persist the launch mode overrides to `agents.local.toml` in a future iteration, or will they remain strictly session-local?
- Should the `Ideas` selection eventually allow triggering a `parley run` for that specific idea?
