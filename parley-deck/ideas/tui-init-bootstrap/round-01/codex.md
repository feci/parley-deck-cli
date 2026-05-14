---
agent: codex
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
---

## Summary

The TUI should become a valid first entry point for a new Parley Deck workspace. The current failure path in `runTUIViewWithDiscovery` should be replaced, only for `parley tui`, with a Bubble Tea setup screen that asks the user to initialize and then reloads the normal workspace status.

## Proposed approach

Add a small setup model in `internal/tui` that is shown when `protocol.ReadWorkspaceStatus(root)` returns `os.ErrNotExist` from `runTUIViewWithDiscovery`. The setup model should render the target path, a concise "initialize workspace" action, and the quit key. Pressing `i` or `enter` should execute a command that calls `protocol.InitWorkspace(root)` and then `protocol.ReadWorkspaceStatus(root)`.

After successful initialization, the model can render the existing dashboard by reusing the current `model` view with the newly loaded `WorkspaceStatus` and already discovered agents. If initialization fails, stay on the setup screen and display the error. If the workspace is already initialized, the existing path should still call `tui.Run(status, results)` unchanged.

Tests should cover the uninitialized model's view, the confirm key path creating a real workspace in a temp directory, and the initialized view transition. An app-level test can assert that the old stderr failure is no longer produced only if it can be done without trying to run an interactive Bubble Tea program.

## Concerns / open questions

Running a full TUI program in `app` tests may be brittle because Bubble Tea expects terminal semantics. Unit tests on the TUI model are likely enough for this narrow change.

## Risks

The main risk is accidentally changing behavior for `parley run` or `parley resume`, which should still require an initialized workspace. Keep the bootstrap path local to `parley tui`.
