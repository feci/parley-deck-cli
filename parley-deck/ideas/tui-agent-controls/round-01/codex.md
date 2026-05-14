---
agent: codex
idea: tui-agent-controls
round: 1
date: 2026-05-14
---

## Summary

The next TUI slice should make the static dashboard operational without turning it into a full workflow wizard. I propose adding selectable agent rows, a selected-agent details panel, and session-local launch-mode switching so the user can inspect and reason about headless versus interactive/manual agents directly in the TUI.

## Proposed approach

Extend `internal/tui/app.go`'s dashboard model with focus and selection state: `focus = ideas|agents`, selected idea index, selected agent index, and a map of per-agent session launch-mode overrides. Use `j/k`, up/down, and tab/shift-tab for navigation. The agent list should mark the selected agent and show its effective launch mode.

Replace the current agent summary panel with a two-part agent view: a compact selectable list and a details section for the selected agent. Details should include installed state, version/probe error, effective launch mode, configured headless command, configured interactive command, prompt mode, invoke strategy, model/profile/reasoning, sandbox, approval, timeout, home isolation, backend, notes, and config source highlights where available.

Add keys on the agent focus only: `h` sets the selected agent to `headless`, `i` sets it to `interactive`, `m` sets it to `manual`, and `x` clears the session override back to configured/default mode. These mode switches are intentionally session-local in this slice. Persisting to `agents.local.toml` should be a later, explicit config-editing feature because it touches user-owned local configuration and needs careful TOML preservation.

For newly initialized workspaces, `initModel` should transition to the real dashboard model after success rather than only rendering the dashboard view from `initModel`. That avoids the known limitation where post-init dashboard keybindings would not work after the new controls are added.

## Concerns / open questions

Session-local launch-mode switching is useful for understanding and previewing runtime behavior, but it will not affect `parley run` after the TUI exits. The UI copy must be clear enough to avoid implying persistence. A later slice should add explicit persistence, likely with a separate `parley agents set-mode` command and a TUI `s` key that writes `agents.local.toml`.

## Risks

The main risk is overbuilding the dashboard into a full run launcher. Keep this slice to navigation, inspection, and session-local mode overrides. Another risk is that text may overflow in the agent detail panel; tests should assert key strings and the view should wrap or truncate long command/details lines.
