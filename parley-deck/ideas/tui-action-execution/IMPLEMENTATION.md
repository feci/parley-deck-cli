---
idea: tui-action-execution
implemented-by: codex
date: 2026-05-25
status: complete
---

## Summary

Implemented the accepted first slice:

- shared planner action command construction in `internal/runaction`;
- dashboard TUI action focus, selection, asynchronous runner invocation, busy
  state, action feedback, and refresh-on-success;
- app-backed execution for consensus draft, request-signoffs, and finalize;
- advisory TUI feedback for inspect, answer-question, retry-agent, and unknown
  unsupported actions;
- focused unit coverage for command mapping, TUI interaction, and app runner
  dispatch helpers.
- review fix-up for `RequiresYes` action confirmation, busy-state coverage, and
  invalid zero-round command construction.

## Files changed

- `internal/runaction/action.go`
- `internal/runaction/action_test.go`
- `internal/tui/app.go`
- `internal/tui/app_test.go`
- `internal/app/app.go`
- `internal/app/app_test.go`

## Verification

- `go test ./internal/runaction ./internal/tui ./internal/app`
- `go test -v ./internal/app`
- `go test ./...`
- `git diff --check` on the implementation files

## Review fix-up

- Added a TUI-level two-step confirmation for actions marked `RequiresYes`.
  The first `enter` now shows a confirmation message and command; the second
  `enter` on the same selected action dispatches the runner.
- Added TUI test coverage for `actionRunning` duplicate-enter rejection.
- Tightened `runaction.RoundNumber("round-00")` to return no round instead of
  `0`, with regression coverage.

## Notes

The working tree also contains unrelated pre-existing changes outside this
slice. They are intentionally excluded from this implementation summary and
should not be reviewed as part of `tui-action-execution`.
