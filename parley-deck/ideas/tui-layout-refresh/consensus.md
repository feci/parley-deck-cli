---
idea: tui-layout-refresh
status: consensus
date: 2026-05-26
participants: [codex, claude, gemini, hermes]
---

# Consensus

The current TUI should be made more readable by reducing vertical stacking, adding height-aware compact rendering, and introducing a restrained semantic color hierarchy. This should be a focused refactor of the existing `View()` rendering code, not a new layout engine or a navigation redesign.

## Agreed direction

- Track terminal height in the dashboard model and use it in both dashboard and live views.
- Replace the dashboard's three top panels plus full-width bottom panels with a clearer default:
  - one header line;
  - a left context column for sessions;
  - a right workspace column for the selected run, agents, latest events, actions, and questions;
  - a single compact footer line instead of a boxed footer.
- Add compact rendering for short terminals:
  - dashboard compact mode should show the selected run/session summary, compact agent/action/question rows, and one key line;
  - live compact mode should show compact agent status, the latest event, the selected question or log summary, and one key line.
- Keep important labels visible in normal mode: `Sessions`, `Agents`, `Run details`, `Event stream`, `Actions`, `Questions`, and `Log preview`.
- Use semantic color only where it improves scanning:
  - green for success/found/done;
  - cyan or blue for active/focus/section titles;
  - amber for waiting, action, question, or warning states;
  - red for failure, error, or missing states;
  - grey for muted or inactive text.
- Preserve current keyboard behavior and existing focus state.
- Add focused tests for low-height compact rendering and retain existing substring coverage.

## Non-goals

- Do not introduce tabs or new pane toggles in this slice.
- Do not move HITL answering behind a different workflow.
- Do not add scrolling or viewport widgets unless plain clipped rendering proves insufficient.
- Do not optimize for very wide terminals first; the user complaint is low-height readability.

## Implementation plan

1. Add shared layout/style helpers for default dimensions, clipping, section titles, and semantic badges.
2. Add `height` to the dashboard model and propagate `WindowSizeMsg.Height`.
3. Rewrite dashboard `View()` into normal and compact branches.
4. Rewrite live `View()` into normal and compact branches.
5. Keep old rendering helpers where useful, but add compact variants to avoid dumping full logs or long metadata into short terminals.
6. Add tests for compact dashboard and compact live output, then run the TUI package tests and full Go test suite.

## Signoff

- codex: approved
