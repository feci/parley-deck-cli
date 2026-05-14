---
idea: tui-init-bootstrap
status: complete
implementer: codex
started: 2026-05-14
completed: 2026-05-14
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#impl-tui-init-bootstrap
head-commit: 20a568c
design-pr: https://github.com/feci/parley-deck-cli/pull/15
implementation-pr: https://github.com/feci/parley-deck-cli/pull/16
---

## Implemented changes

- `parley tui` now opens a setup TUI when the target workspace is missing instead of exiting with "run `parley init` first".
- The setup screen shows the target path, supports `i` or `enter` to initialize, and supports `q`, `esc`, or `ctrl+c` to quit.
- Initialization reuses `protocol.InitWorkspace`, reloads with `protocol.ReadWorkspaceStatus`, and transitions to the existing dashboard in the same TUI session.
- Initialization failures stay on the setup screen, show an inline error, and allow retry.
- Already-initialized workspaces continue to use the existing dashboard path.

## Tests and verification

- `go test ./internal/tui`
- `go test ./...`
- PTY smoke: `parley tui --dir "$(mktemp -d)"`, press `i`, then `q`; verified `parley-deck/COOPERATION.md` was created.

## Deviations from FINAL

- None.

## Review result

- Review consensus: ./review/consensus.md
- Required fixes: none.
- Deferred follow-ups are non-blocking and recorded in review consensus.
