---
idea: tui-agent-controls
author: codex
created: 2026-05-14
participants: [codex, claude, gemini, hermes]
status: final
---

## Problem / idea

The current `parley tui` dashboard is mostly read-only. It shows ideas and agents, but the user cannot meaningfully navigate or operate agent runtime settings from inside the TUI. The user wants CLI functionality exposed in the TUI, especially the ability to move between configured models/agents and deal with agents that are not running in headless mode.

Add a focused first slice of interactive TUI controls so the dashboard can select agents, show the selected agent's resolved runtime details, and switch the selected agent's launch mode between `headless`, `interactive`, and `manual` for the current TUI session.

## Constraints

- Keep the slice small and implementable in one PR.
- Preserve existing `parley agents list`, `parley run`, `parley resume`, and live-run TUI behavior.
- Do not start paid or hosted agent processes from the static dashboard.
- Do not persist changes to `agents.local.toml` in this first slice unless consensus explicitly decides the added write risk is justified now.
- Show enough runtime detail for a user to understand what will happen outside headless mode: launch mode, command shape, prompt mode, invoke strategy, model/profile, sandbox, approval, timeout, backend, and notes.
- Make keyboard behavior predictable: visible selection marker, `j/k` and arrow keys to move, `tab`/`shift+tab` to switch dashboard focus, and mode toggle keys.
- Keep initialized-workspace and newly initialized TUI behavior working.
- Add focused tests for navigation, focus switching, mode override behavior, and rendered details.

## Non-goals

- Do not build a full workflow creation wizard.
- Do not implement terminal multiplexing or embed external CLIs inside the Bubble Tea dashboard.
- Do not secretly automate interactive sessions through a PTY.
- Do not change agent runtime schema in this idea.
- Do not change billing semantics or make claims about provider accounting.
