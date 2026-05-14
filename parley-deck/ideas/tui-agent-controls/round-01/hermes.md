---
agent: hermes
idea: tui-agent-controls
round: 1
date: 2026-05-14
---

## Summary
The proposal adds a minimal interactive layer to the existing TUI dashboard. Users can navigate the list of configured agents/models, view resolved runtime configuration for the selection, and toggle the launch mode (headless/interactive/manual) for the current session only. No paid or hosted agents are started from the dashboard, and no changes are written to agents.local.toml.

## Proposed approach
- Extend the current agents pane with a selectable list using j/k and arrow keys.
- On selection change, resolve and render the runtime details pane showing: launch mode, command shape, prompt mode, invoke strategy, model/profile, sandbox, approval, timeout, backend, and any notes.
- Add a mode-toggle key (e.g. m) that cycles the override for the selected agent and immediately updates the rendered details.
- Keep all existing CLI commands and live-run behavior untouched; the TUI only supplies an in-session override flag when the user later invokes `parley run` or `parley resume` from the same session.
- Navigation and focus follow the documented tab/shift+tab rules; a visible selection cursor is shown at all times.

## Concerns / open questions
- How should the resolved details pane handle agents that have no local definition yet (purely remote references)?
- Should the mode override be visible in the status bar or only inside the details pane?
- Is there a safe default when an agent has never been run before?

## Risks
- Even session-only overrides could confuse users if they forget they changed the mode before launching.
- Rendering too much detail at once may make the pane overflow on narrow terminals; we must keep the slice strictly to the listed fields.