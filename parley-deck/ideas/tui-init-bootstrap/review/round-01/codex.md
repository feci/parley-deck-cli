---
agent: codex
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
reviewed-commit: 20a568c
implementation-pr: https://github.com/feci/parley-deck-cli/pull/16
---

## Findings

No blocking findings.

## Notes

- The uninitialized `parley tui` path now routes to `tui.RunInit` instead of returning the old "run `parley init` first" error.
- The setup model requires an explicit `i` or `enter` keypress, calls `protocol.InitWorkspace`, reloads status, and renders the existing dashboard view on success.
- The failure path remains on the setup screen, records the error, and allows retry.
- Existing initialized-workspace behavior still calls `tui.Run(status, results)`.

## Verification

- `go test ./internal/tui`
- `go test ./...`
- PTY smoke with a temp directory: setup view appeared, pressing `i` transitioned to the dashboard, pressing `q` exited, and `parley-deck/COOPERATION.md` existed afterward.

## Residual risk

Agent discovery still runs before the setup screen because the existing TUI entry point discovers agents before reading workspace status. This preserves the current initialized path and keeps the patch small, but it means the setup screen can wait on agent version probes.
