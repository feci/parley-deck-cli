---
idea: live-run-tui
status: fix-up-cycle-1
implementer: codex
started: 2026-05-10
completed: 2026-05-10
branch: parley-deck-cli#feature/live-run-tui
head-commit: 1ec8284d66ab331b5a27a4326d48162a93cd0eef
design-pr: https://github.com/feci/parley-deck-cli/pull/1
implementation-pr: https://github.com/feci/parley-deck-cli/pull/2
---

## Summary of work

Implemented the `live-run-tui` FINAL.md slice:

- Added a small async runner handle with `RunID`, `RunDir`, `Done()`, `Wait()`, and result storage.
- Kept `parley run --no-tui` on the existing synchronous runner path.
- Changed default `parley run` to start the runner asynchronously and open a live Bubble Tea view before the round finishes.
- Added a live TUI model that polls `events.jsonl`, derives per-agent state and round status, shows recent event summaries, and previews bounded stdout/stderr log tails for the selected agent.
- Added TUI controls for `j`/`k`/`tab` selection, `q`/`esc` TUI detach, and `ctrl+c` run cancellation.
- Added focused tests for async runner completion, event projection, partial event reads, log tailing, and view content.

## Deviations from FINAL.md

- None known.

## Notes for reviewers

- The live TUI intentionally reads from disk instead of using an in-memory event channel.
- If the runner completes before the final `round.completed` or `round.incomplete` event is read, the TUI performs one final event read before exiting.
- Verification run: `GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`

## Fix-up cycle 1

Review round 1 found blocking issues in skipped-agent duration, sticky `unknown` state, and table alignment. The fix-up:

- Keeps skipped duration at zero unless the event supplies `duration_ms`.
- Keeps non-participant agents in `unknown` state across all later events.
- Removes ANSI styling from padded state table cells and aligns headers/rows with fixed-width plain columns.
- Derives the round label from `IdeaStatus.Status`.
- Shows terminal zero-duration states as `0s`.
- Extends reducer coverage for skipped-without-start and unknown-after-multiple-events.

Verification: `GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
