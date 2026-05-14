---
agent: gemini
idea: tui-agent-controls
round: 1
date: 2026-05-14
---

# Summary

The proposed "Agent Inspector" slice for the Parley Deck TUI introduces interactive agent selection and runtime configuration visibility. By allowing users to navigate the agent list and inspect resolved settings, we bridge the gap between static configuration files and live execution behavior. This slice focuses on keyboard-driven navigation and session-local launch mode overrides, enabling rapid experimentation with `headless`, `interactive`, and `manual` modes without permanent side effects to the workspace configuration.

# Proposed approach

1.  **Dual-Pane Navigation and Focus**:
    *   Implement focus management between the **Protocol** list and the **Agents** list using `Tab` / `Shift+Tab`.
    *   In the **Agents** list, use `j`/`k` or arrow keys to move a selection cursor (`>`).
    *   Highlight the focused pane with a distinct border style (e.g., changing from a muted color to a bright color) to provide clear visual feedback.

2.  **Agent Inspector Pane**:
    *   Add an "Inspector" panel that appears when an agent is selected. For 80+ column terminals, this can be a third vertical column; for narrower screens, it can be a bottom-docked panel above the footer.
    *   For the selected agent, display a structured view of its **Spec** (already embedded in the `Discovery` struct):
        *   **Launch Details**: Current launch mode (`headless`, `interactive`, `manual`), resolved command line with placeholders (`{root}`, `{tempdir}`) expanded, and timeout.
        *   **Capabilities**: Prompt mode (stdin/arg), Sandbox mode, and Approval policy.
        *   **Identity**: Model name, profile, and backend (local/hosted/unknown).
        *   **Notes**: Display both `InteractiveNotes` and general `Notes` for operational context.

3.  **In-Session Launch Mode Toggles**:
    *   Add a key binding (`m`) to cycle the `LaunchMode` for the selected agent for the current TUI session.
    *   Store these overrides in a `map[string]string` in the TUI model.
    *   Provide a visual indicator (e.g., a "MODIFIED" tag or color change) in the agent list when the mode differs from the underlying `Spec`.
    *   Ensure these overrides are passed into the `runner` when an idea is executed from the TUI.

4.  **UI Feedback and Legends**:
    *   Add a dynamic footer legend: `tab: switch focus`, `j/k: select agent`, `m: cycle launch mode`, `q: quit`.
    *   Use `lipgloss` styles to color-code launch modes: `headless` (green), `interactive` (yellow), `manual` (blue).

# Concerns / open questions

- **Layout Density**: A three-pane layout might be too dense. Should the Inspector be a toggleable overlay (`i`) or always visible? Given the goal of "showing resolved details", always visible (with responsive stacking) is preferred.
- **Model Navigation**: If "navigating models" implies changing the `Model` field, we should decide if we want to support cycling through a set of predefined models for the first slice, or if just viewing the configured model is sufficient.
- **Manual Mode Clarification**: Does `manual` mode require any special handling in the TUI other than displaying a "Manual Action Required" message when a run reaches that agent?

# Risks

- **Stale State**: External changes to `agents.local.toml` won't be reflected without a TUI restart or a manual refresh key (`r`).
- **Input Conflict**: As we add more shortcut keys, we must ensure they don't conflict with core navigation or future features.
- **Accidental Overrides**: Users might accidentally toggle a mode and forget they've done so. Clear visual differentiation between "config-defined" and "session-overridden" modes is critical.
