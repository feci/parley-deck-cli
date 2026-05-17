---
idea: tui-workspace-sessions
status: ready-for-review
implemented-by: codex
branch: feature/tui-workspace-sessions
date: 2026-05-17
---

## Summary

Implemented the first MVP slice of the workspace session console:

- Added `internal/sessionstore` for safe user-local session metadata in `~/.parley-deck/sessions.json`, with `PARLEY_HOME` override for tests.
- Added `internal/runcontrol` so CLI and TUI use a shared run creation path.
- Added shared run attention derivation in `internal/runstate`.
- Updated `parley run` to register session metadata through the shared run creation path.
- Reworked `parley tui` into a workspace run console with session list, event stream, run/agent details, HITL question display, refresh, and `N` start-new-idea mode.
- Kept TUI-started runs in-process and made the footer explicit that quitting does not provide detached execution.

## Deviations from FINAL.md

- The first UI uses a compact selected-agent detail pane rather than a separate full-screen per-agent view.
- The selected-run event stream uses the existing compact recent event projection for now, not a larger scrollback buffer.
- Cross-workspace recent sessions are persisted, but the TUI currently prioritizes current-workspace runs.
- No global question queue was implemented; it remains deferred as planned.

## Verification

- `go test ./...`
- `go run ./cmd/parley status --dir .`
- `go run ./cmd/parley status --dir . --json`
- `go run ./cmd/parley tui --dir .` smoke-tested in a PTY and exited with `q`.

## Ready for review

Reviewers should focus on:

- Whether the session registry can corrupt or mislead users.
- Whether `parley run` behavior stayed compatible.
- Whether TUI-started run lifecycle semantics are clear enough.
- Whether the TUI model has hidden race or stale-state problems.
