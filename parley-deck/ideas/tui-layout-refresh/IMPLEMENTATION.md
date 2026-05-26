---
idea: tui-layout-refresh
date: 2026-05-26
implementer: codex
status: implemented
---

# Implementation

Implemented the consensus plan as a focused TUI rendering refresh.

## Changed

- Added shared layout helpers in `internal/tui/layout.go`:
  - minimum/default terminal dimensions;
  - height-aware clipping;
  - compact line helpers;
  - semantic attention/state/risk badges.
- Updated dashboard rendering in `internal/tui/app.go`:
  - tracks terminal height from `tea.WindowSizeMsg`;
  - normal mode renders a header, a sessions context column, a combined workspace column, and one compact footer line;
  - compact mode activates below 30 rows and renders sessions, run summary, agents, actions, questions, and footer without nested boxes.
- Updated live run rendering in `internal/tui/live.go`:
  - normal mode renders agents/questions on the left and latest events/log preview on the right;
  - compact mode activates below 30 rows and renders agents, latest event, selected question, log summary, and footer as clipped lines;
  - selected HITL question detail now includes the full prompt on its own line for readability.
- Added focused tests for compact dashboard and compact live output fitting short terminal heights.
- After review, removed obsolete dashboard render methods, restored the normal-mode process lifecycle footer warning, unified compact header labeling, and added threshold-boundary tests.

## Verification

- `go test ./internal/tui` passed.
- `go test ./...` passed.
