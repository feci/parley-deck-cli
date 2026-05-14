---
agent: hermes
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
reviewed-commit: 20a568c
implementation-pr: https://github.com/feci/parley-deck-cli/pull/16
---

## Findings

No blocking findings.

## Verification

- Confirmed the PR routes missing workspace errors in `internal/app/app.go` to the TUI init flow instead of exiting with the old `parley init` message.
- Confirmed `internal/tui/app.go` adds an explicit setup state with `i`/`enter` initialization, `q`/`esc`/`ctrl+c` quit handling, duplicate-init guarding, inline failure display, and success transition to the existing dashboard view.
- Confirmed initialization reuses `protocol.InitWorkspace` and reloads via `protocol.ReadWorkspaceStatus`.
- Confirmed `internal/tui/app_test.go` covers prompt rendering, successful initialization and dashboard transition, and failure staying on the setup screen.

## Residual risk

- Agent discovery still happens before the setup screen, so a slow version probe can delay first render in an uninitialized directory.
- The active model remains `initModel` after success and only renders the dashboard view. This is acceptable for the current dashboard keys, but future dashboard-specific keybindings should be delegated or transitioned to the normal model.
- Agents are not re-discovered after initialization; this is acceptable while discovery is independent of newly created workspace files.