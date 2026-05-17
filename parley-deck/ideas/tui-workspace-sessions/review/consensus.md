---
idea: tui-workspace-sessions
review-round: 1
drafted-by: codex
date: 2026-05-17
reviewed-commit: 8a9fe81
---

## Agreed findings

1. MAJOR: TUI-started run lifecycle semantics needed to be made true in code.
   - Source: `review/round-01/claude.md`
   - Resolution: fixed.
   - Fix: `runTUIViewWithDiscovery` now stores cancel functions for TUI-started runs and cancels them when the TUI exits. The previous lost-cancel pattern was removed.

2. MAJOR: TUI refresh could shift selection because selected run and agent were tracked only by index.
   - Source: `review/round-01/claude.md`
   - Resolution: fixed.
   - Fix: refresh handling now preserves selected run by `RunID` and selected agent by agent ID before clamping.

3. MAJOR: session registry writes were too aggressive on every TUI refresh.
   - Sources: `review/round-01/claude.md`, `review/round-01/gemini.md`
   - Resolution: fixed.
   - Fix: refresh no longer writes every run back to `~/.parley-deck/sessions.json`; workspace sessions are registered on initial TUI load and run lifecycle updates.

4. MAJOR: session registry updates lacked cross-process coordination.
   - Source: `review/round-01/gemini.md`
   - Resolution: fixed.
   - Fix: `sessionstore.Upsert` now takes a lock file around read-modify-write and preserves existing `LastEventAt` when the update omits it.

5. MINOR: refresh ticks could queue while a prior refresh was still running.
   - Source: `review/round-01/gemini.md`
   - Resolution: fixed.
   - Fix: the next refresh tick is scheduled only after the current refresh result is processed.

## Deferred findings

- Full incremental event-log loading is deferred. The current MVP still uses `runstate.ListRuns`, which rereads current run summaries. This is acceptable for the first slice but should be revisited if real workspaces grow large.
- Batch-upserting all initially discovered workspace runs would be more efficient than repeated `Upsert` calls during startup. The per-tick write problem is fixed; startup batching can wait.
- Benchmarks and large-workspace performance tests are deferred until the TUI has real user data.

## Dismissed findings

- `AttentionAction` priority over failed state is not changed in this fix-up. Open HITL questions are intentionally surfaced first because they are immediately user-actionable; failed runs still render as `FAILED` when there is no open question.

## Verification

- `go test ./...`
- `go run ./cmd/parley status --dir .`
- `go run ./cmd/parley tui --dir .` in a PTY, exited with `q`.

## Signoffs

<!-- Each reviewer APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex - 2026-05-17
Status: ✅ ACCEPT
Notes: Agreed fixes are implemented in the current working tree and tests pass.

### Signoff: hermes - 2026-05-17
Status: ✅ ACCEPT
Notes: All agreed findings addressed adequately.

### Signoff: claude - 2026-05-17
Status: ✅ ACCEPT
Notes: Findings 1-2 (lost cancel, index-tracked selection) and Finding 3 follow-ups (per-tick rewrite, `LastEventAt` clobber, cross-process race) are addressed in `internal/app/app.go`, `internal/tui/app.go`, and `internal/sessionstore/sessionstore.go`. Deferring batch upsert and incremental event-log loading is reasonable now that registry writes no longer occur on every refresh.

### Signoff: gemini - 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The per-refresh session write issue is fixed, sessionstore updates are lock-protected, and deferred performance work is acceptable for this MVP.
