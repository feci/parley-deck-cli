---
agent: hermes
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
---

## Summary
The goal is to allow `parley tui` (including with `--dir`) to launch in an uninitialized directory. The TUI detects the missing `parley-deck/COOPERATION.md`, presents a minimal confirmation prompt showing the target path, performs initialization using the existing `protocol.InitWorkspace` call, and transitions directly into the normal dashboard view. `parley init` behavior and all other commands remain unchanged. Scope is limited to the TUI bootstrap path plus targeted tests.

## Proposed approach
- In the TUI entry point, before rendering the main dashboard, invoke `protocol.ReadWorkspaceStatus` (or equivalent lightweight check) on the resolved target directory.
- If the workspace is uninitialized, render a focused prompt screen that displays the absolute target path and offers a single "Initialize workspace" action (plus quit).
- On confirmation, call the same `protocol.InitWorkspace` logic currently used by `parley init`, handling any errors by showing them inline before returning to the prompt.
- After successful initialization, immediately switch to the normal TUI dashboard for that workspace.
- Add two narrow tests: one verifying the prompt appears for an uninitialized directory, and one confirming that confirming the prompt produces the expected files and enters the dashboard state.
- No changes to command routing, protocol definitions, or non-TUI paths.

## Concerns / open questions
- How should transient errors during `InitWorkspace` (e.g., permission issues) be presented to avoid trapping the user in the prompt?
- Should the prompt support a "Choose different directory" option, or keep it strictly to yes/no on the provided `--dir`?

## Risks
- Low risk: the change touches only the TUI bootstrap sequence and reuses existing initialization semantics, so no impact on `parley init`, other commands, or the Parley Deck artifact protocol. Tests will be added to guard the new flow.