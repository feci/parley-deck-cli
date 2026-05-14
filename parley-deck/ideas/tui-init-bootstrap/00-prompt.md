---
idea: tui-init-bootstrap
author: codex
created: 2026-05-14
participants: [codex, claude, gemini, hermes]
status: round-01
---

## Problem / idea

`parley tui` currently refuses to start unless the target directory was already initialized with `parley init`. The user wants to open the TUI first and have the TUI ask whether to initialize the Parley Deck workspace from inside the UI.

The change should make `parley tui --dir <path>` usable in an uninitialized directory. If no `parley-deck/COOPERATION.md` exists, the TUI should show a clear setup prompt, let the user confirm initialization, call the same initialization behavior as `parley init`, and then switch to the normal TUI dashboard.

## Constraints

- Keep `parley init` behavior intact.
- Do not automatically initialize without an explicit user action inside the TUI.
- Show the target workspace path before initialization.
- Reuse existing `protocol.InitWorkspace` and `protocol.ReadWorkspaceStatus` semantics.
- Keep the normal initialized-workspace TUI behavior unchanged.
- Keep changes scoped to the TUI entry path and focused tests.
- Preserve the active `github-pr` transport and Parley Deck artifact flow.

## Non-goals

- Do not add a broader first-run onboarding wizard.
- Do not change `parley run`, `parley resume`, or non-TUI commands to initialize implicitly.
- Do not change the Parley Deck protocol or workspace layout.
- Do not add new agent runtime configuration behavior in this idea.
