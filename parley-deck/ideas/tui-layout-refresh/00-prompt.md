---
idea: tui-layout-refresh
author: codex
created: 2026-05-26
participants: [codex, claude, gemini, hermes]
roles:
  codex: Go TUI implementation and terminal constraints
  claude: UX information architecture and low-height readability
  gemini: color/state hierarchy and responsive layout
  hermes: operational dashboard and long-running session visibility
status: final
---

## Problem / idea

The user reports that the current Parley Deck TUI is cramped, visually unclear,
too vertical, and does not fit in a short terminal. The dashboard stacks too
many panels, uses long raw text, and makes status scanning difficult.

Design and implement a more readable terminal layout for both the dashboard and
live run views. The UI should be split into clearer regions, use color and
status badges more intentionally, reduce vertical waste, and degrade gracefully
on narrow or low-height terminals.

## Constraints

- Keep the implementation in the current Bubble Tea/Lip Gloss stack.
- Avoid a broad rewrite; make a focused layout refresh that can be tested.
- Text must not overflow or force unusable horizontal/vertical scrolling.
- Support short terminals by showing a compact mode with the most important
  state first.
- Preserve existing keyboard behavior unless a layout change requires clearer
  footer text.
- Protocol files remain English.

## Non-goals

- Do not introduce a web UI.
- Do not add new workflow semantics.
- Do not redesign agent runtime configuration.
