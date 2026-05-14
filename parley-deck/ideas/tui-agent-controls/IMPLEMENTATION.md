---
idea: tui-agent-controls
status: implemented
implementer: codex
started: 2026-05-14
completed: 2026-05-14
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#impl-tui-agent-controls
head-commit: f76d6ab
design-pr: https://github.com/feci/parley-deck-cli/pull/17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/18
---

## Summary of work

- Added focus and selection state to the static dashboard TUI.
- Added `tab` / `shift+tab` focus switching between ideas and agents.
- Added `j`/`down` and `k`/`up` row navigation.
- Added session-local selected-agent launch-mode overrides with `h`, `i`, `m`, and `x`.
- Added selected-agent runtime details covering launch mode, runtime config, command shape, prompt/invoke settings, backend, timeout, notes, and probe errors.
- Kept dashboard mode overrides preview-only; the dashboard does not launch agents or persist config.
- Changed successful first-run initialization to return the real dashboard model so dashboard keybindings work immediately.

## Deviations from FINAL.md

- None.

## Notes for reviewers

- The mode override map intentionally does not feed into `parley run`; this slice is inspection and preview only.
- The detail panel uses already resolved `agents.Discovery` data and does not re-read config.
- Verification run:
  - `go test ./internal/tui`
  - `go test ./...`
  - PTY smoke: opened `parley tui --dir .`, sent dashboard keys, and exited with `q`.
