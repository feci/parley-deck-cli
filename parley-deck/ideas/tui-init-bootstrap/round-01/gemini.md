---
agent: gemini
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
---

## Summary

Allow `parley tui` to start in an uninitialized directory by deferring the workspace check. If no workspace is found, the TUI will display a setup screen that allows the user to initialize the current directory using existing protocol logic. This provides a smoother onboarding experience for new users who want to explore the project through the TUI immediately.

## Proposed approach

- **App Entry Point:** Modify `runTUIViewWithDiscovery` in `internal/app/app.go` to handle `os.ErrNotExist` from `protocol.ReadWorkspaceStatus`. Instead of exiting with an error, it should proceed to `tui.Run` with an uninitialized state. We must ensure `WorkspaceStatus.Root` is correctly set even when the workspace doesn't exist yet.
- **TUI Model (internal/tui/app.go):**
    - Add an `Uninitialized bool` field to the `model` struct.
    - Set `Uninitialized` to `true` if `protocol.ReadWorkspaceStatus` failed with `os.ErrNotExist`.
- **TUI View:**
    - If `m.Uninitialized` is true, render a special "Setup" view instead of the normal dashboard.
    - This view will display the target workspace path and a clear prompt: "Parley Deck is not initialized in this directory. Press 'i' to initialize now."
    - Use `lipgloss` to style this prompt to be visually distinct and welcoming.
- **TUI Update:**
    - Handle the 'i' key in `Update`.
    - When 'i' is pressed, call `protocol.InitWorkspace(m.status.Root)`.
    - If initialization succeeds:
        - Call `protocol.ReadWorkspaceStatus(m.status.Root)` to refresh the status.
        - Set `m.Uninitialized = false`.
        - Re-run agent discovery if necessary (or reuse existing results).
    - If initialization fails (e.g., permission denied), display an error message in the TUI instead of quitting.
- **Protocol:** Update `protocol.ReadWorkspaceStatus` in `internal/protocol/workspace.go` to return a `WorkspaceStatus` containing at least the `Root` path even when `COOPERATION.md` is missing, or change the caller to handle the missing status gracefully.
- **Preserve Command Behavior:** Other commands like `run` or `status` should continue to require manual initialization as per the non-goals.

## Concerns / open questions

- **Agent Discovery:** Should we run agent discovery before initialization? Currently, `runTUIView` runs discovery first. Discovery might depend on project-specific configuration that isn't present yet. We should ensure discovery handles an uninitialized directory gracefully (which it seems to do via `config.LoadAgentSpecs`).
- **Confirmation Flow:** Should we use a more formal confirmation (e.g., 'y/n' prompt) or just the 'i' shortcut? Given the constraints, a clear "Press 'i' to initialize" seems sufficient and consistent with TUI patterns.

## Risks

- **Accidental Initialization:** A user might start `parley tui` in the wrong directory and accidentally initialize it. Mitigated by showing the full path in the setup view and requiring an explicit 'i' keypress.
- **Partial Initialization:** If `protocol.InitWorkspace` fails halfway through, the workspace might be in an inconsistent state. Reusing the existing robust `InitWorkspace` logic minimizes this risk.
- **UI Consistency:** Ensure the transition from the setup view to the dashboard is seamless and doesn't flicker or lose state.
- **Testing Complexity:** Testing the interactive initialization flow in a headless TUI environment requires careful setup of mock directories and event simulation.
