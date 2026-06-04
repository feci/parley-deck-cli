---
agent: codex
idea: tui-interactivity-overhaul
review-round: 1
date: 2026-06-04
---

## Summary

The implementation covers the intended slices and the main segment reset path appears directionally correct: targeted agents are reset on `run.segment_started`, standard runner events are tagged, `?` is split from HITL answer mode, and steering is queued-only rather than executed. I found no CRITICAL issue, but there are MAJOR durability/correctness gaps in the focus viewport bounds and steering ID persistence that should be fixed before accepting the implementation.

Verification attempted:
- `git diff main...HEAD --stat` was inspected.
- `go build ./...`, `go vet ./...`, and `go test ./...` first failed because the default Go cache is outside the sandbox.
- Rerunning with `GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache` failed because dependencies were not cached locally and network access to `proxy.golang.org` is unavailable.

## Findings

### [MAJOR] Focus append path is not bounded by the required 4 MiB cap

`internal/tui/live.go:939-946` only caps appended focus output by `maxFocusLines`, while `readAppendedLines` uses `io.ReadAll` from the saved offset at `internal/tui/live.go:1098`. A long line, a large burst between ticks, or a process that writes megabytes without newlines can exceed the mandatory 4 MiB scrollback cap and can repeatedly allocate the full unread suffix every 250 ms. This violates FINAL.md's "20,000 lines OR 4 MiB per stream" invariant and can still make the TUI grow or stall on noisy agent output.

Concrete fix: keep the focus buffer byte-budgeted as well as line-budgeted. Read at most a bounded window from the current offset, preserve a separate partial-line buffer, and after appending complete lines evict oldest content until both `len(lines) <= maxFocusLines` and retained bytes <= `maxFocusBytes`. Add tests for a >4 MiB single line and a >4 MiB appended burst.

### [MAJOR] Steering IDs can collide across the new CLI/TUI entry points

`internal/steer/steer.go:93-118` computes `steer-NNNN` by loading and counting existing `steer.requested` events, then appends `steer.requested` and `steer.delivered` as two separate writes. Because `store.Append` only has a process-local mutex (`internal/store/events.go:24-57`), two concurrent `parley steer` processes, or CLI plus live TUI, can both assign the same ID and interleave the requested/delivered pair. The branch now exposes exactly those multiple entry points, so this is no longer only a future live-delivery problem; duplicate IDs corrupt `steer.List` projection and make `continue --json` ambiguous.

Concrete fix: either use collision-resistant IDs that do not require read-count-write state, or add a file lock / atomic append batch covering "load latest id + append both events" across processes. If monotonic numeric IDs are required, put the lock around the full critical section and add a concurrency test with two independent stores/process-like submitters.

### [MINOR] Untagged events do not generally inherit the latest segment

FINAL.md requires compatibility where "events without explicit ids inherit the latest known segment." The projection currently sets `Segment` from `run.segment_started` only for listed targets (`internal/runstate/runstate.go:357-378`) and from explicit event data in `applyAgentEvent` (`internal/runstate/runstate.go:430-434`). An untagged `agent.*` event for a non-targeted or late materialized agent after a segment boundary does not inherit the current segment; unknown agents are also created as `StateUnknown` and the first event is skipped at `internal/runstate/runstate.go:344-350`. This leaves edge cases mis-scoped and falls short of the stated backward-compat rule.

Concrete fix: track `currentSegment` while iterating events and pass it into `applyAgentEvent`; when `segment_id` is absent, set the agent segment to `currentSegment` before applying the event. Add tests for an untagged `agent.started` after `run.segment_started`, including a late/unknown agent if those are intentionally projected.

### [MINOR] Initial partial stdout lines are split incorrectly after completion

`loadFocusTail` returns `offset=size` even when the last loaded line does not end with `\n` (`internal/tui/live.go:1078-1083`). If the agent later completes that line, `readAppendedLines` starts at the old EOF and displays only the appended suffix as a new line. Logs often stream progress text without newlines, so the focus view can show duplicated or fragmented output.

Concrete fix: treat the final unterminated line as a partial tail: either do not advance the offset past it, or store it separately and merge it with the next complete append. Add a test where the initial file is `partial` and the next append is ` line\n`.

### [NIT] Overview preview still uses the old hard six-line tail

FINAL.md asks the overview right-pane preview to show "as many lines as fit" rather than the old hard six-line tail, but `refreshLogPreview` still calls `tailLogFile(..., 4096, 6)` for stdout and stderr at `internal/tui/live.go:648-665`. This is less severe because the new focus view handles full logs, but it is a visible spec miss in the overview.

Concrete fix: compute the available preview rows from the current layout height and pass that to `tailLogFile`, splitting stdout/stderr within the pane budget.

## Open questions

- Is cross-process steering intentionally acceptable for Phase 6 despite FINAL.md's F4 warning? If yes, the implementation should explicitly downgrade the promise of monotonic steering IDs or document that simultaneous CLI/TUI steering is unsupported until the file-lock follow-up lands.
- Should the focus view be stdout-only for this cycle? IMPLEMENTATION.md documents stderr streams as deferred, but FINAL.md described stdout/stderr per-stream scrollback.
