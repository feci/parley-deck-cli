---
agent: hermes
idea: unified-tui-home
review-cycle: 2
date: 2026-06-04
---

## Summary
Re-review of fix-up commit 91e0e64. Verified AF1–AF6 via diffs, source, `TestTranscriptPopulatesFromOnDiskRun`, and green `go build && go test ./...`.

## Findings per Agreed Fix

### AF1 — TUI detach must not cancel N-launched runs
VERIFIED. `newLaunchFunc` no longer takes `track`; launched handles are reaped in a background goroutine (`handle.Wait` + `registerWorkspaceSessions`). Only the attached run's `Cancel` (ctrl+c) cancels. Defer-cancel-all removed from `runTUIViewWithDiscovery`.

### AF2 — `parley run` must not own secondary N-launched runs
VERIFIED. `Start` removed from `parley run` LiveOptions (comment notes "start new ideas from `parley tui`"); `Root` retained for Home.

### AF3 — Real on-disk transcript test (`TestTranscriptPopulatesFromOnDiskRun`)
VERIFIED. Writes real `events.jsonl` (`agent.started` with `stdout` path) + `stdout.log`, drives `eventsMsg`, asserts non-empty `buffers["codex"]` AND rendered View contains output. Genuine projection → buffer → render path.

### AF4 — Done-state exit hint
VERIFIED. `renderStatusLine` adds `[done]` tag; `renderInputRow` shows `/quit or esc to exit` hint when `m.done`.

### AF5 — Dead workspace model deleted
VERIFIED. `internal/tui/app.go` reduced to init wizard + shared styles/helpers (1443→150 lines). All workspace symbols, `runTUIAction`/`consensusActionArgs`/`applySessionLaunchOverrides`, and dead tests removed. Init wizard now `tea.Quit`s; caller re-reads status and opens unified Home.

### AF6 — Help wording
VERIFIED. `ctrl+c cancel the attached run, else quit`.

All agreed fixes correctly and completely applied. Build and tests green. No remaining gaps.