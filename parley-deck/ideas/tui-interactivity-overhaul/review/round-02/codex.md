---
agent: codex
idea: tui-interactivity-overhaul
review-round: 2
date: 2026-06-04
responding-to: [claude/IMPLEMENTATION fix-up-cycle-1, hermes/review-round-01]
---

## Position changes since prior review round

The fix-up resolves most of my round-01 findings. AF3, AF4, AF5, and AF9 are addressed in code, and Hermes' AF1, AF6, AF7, and AF8 are also addressed.

I cannot say zero agreed fixes remain because AF2 is only partially resolved: the append/tail read window is now bounded, but `capFocusLines` still allows a retained single line larger than `maxFocusBytes`, so the documented "both bytes and lines" helper invariant is not actually held in isolation.

Verification attempted:
- Inspected `git diff 09b28e6..HEAD` and re-read the changed implementation/tests in `internal/tui/live.go`, `internal/steer/steer.go`, `internal/runstate/runstate.go`, `internal/runner/acp.go`, and `internal/app/app.go`.
- Ran `GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go test ./...`. It still cannot complete offline because uncached modules require `proxy.golang.org`; packages not requiring those downloads passed, including `internal/runstate`, `internal/runner`, and `internal/steer`.

## Responses to other reviewers

### @hermes

Your AF1 execution-overpromise finding is resolved: TUI, help, and CLI copy now say steers are recorded/queued and not auto-run yet.

Your AF6 ACP tagging finding is resolved for the agent lifecycle events that drive the sticky badge: `agent.started` and terminal `agent.finished`/`agent.failed` now carry `segment_id`. The auxiliary ACP detail events remain untagged, but they are not projected as agent state, so I do not consider that a remaining blocker for this cycle.

Your AF7 request for an integration-style sticky-badge test is resolved by `TestLoadRunSegmentUnsticksFinishedBadge`.

Your AF8 `parley steer` validation finding is resolved for empty text and non-participant warning behavior. The warning remains non-fatal, which matches the consensus text.

I agree the carried deferrals are acceptable: F4 cross-process append ordering/locking and slice 5 execution/delivery are explicitly outside this fix-up cycle, as long as the UI keeps the current honest wording.

## Updated findings

### [MAJOR] Focus byte cap still fails for a retained single oversized line

Status from my prior MAJOR AF2 finding: NOT-RESOLVED.

The fix-up bounds the read window in `readAppendedLines` and uses `completeLinesFrom` to drop a leading partial when jumping into the last 4 MiB. That covers the intended >4 MiB burst path and many >4 MiB single-line file-read cases. However, `capFocusLines` stops evicting when one line remains:

- `internal/tui/live.go:1158-1167` computes retained bytes, then loops only while `total > maxFocusBytes && len(lines) > 1`.
- Therefore `capFocusLines([]string{strings.Repeat("x", maxFocusBytes+1)})` returns a buffer over the 4 MiB cap.

So the 4 MiB invariant is not actually held by the buffer-capping helper for a >4 MiB single retained line, even though a >4 MiB multi-line burst is now capped by eviction. The current tests miss this exact case: `TestFocusTailDropsOversizedLine` covers an oversized leading line plus a short tail after a tail-window jump, and `TestCapFocusLinesByteBudget` covers many modest lines, but not one oversized retained line.

Concrete fix: make `capFocusLines` drop or trim the final oversized line when `total > maxFocusBytes`, and add a test that passes exactly one line larger than `maxFocusBytes`. If preserving the newest oversized line is desired, truncate it to the newest `maxFocusBytes-1` bytes and mark `truncated=true`; otherwise return an empty capped buffer.

### [MAJOR] Steering IDs can collide across the new CLI/TUI entry points

Status from my prior MAJOR AF3 finding: RESOLVED for this cycle.

IDs are now `steer-NNNN-<rand>` and the projection groups by the full ID (`internal/steer/steer.go:126-198`). Concurrent submitters may still share the monotonic prefix because F4 locking is deferred, but distinct requests should not merge unless the random suffix collides. That is acceptable within the consensus scope.

### [MINOR] Untagged events do not generally inherit the latest segment

Status from my prior MINOR AF5 finding: RESOLVED for participant/known-agent projection.

`ProjectEvents` tracks `currentSegment` and `applyAgentEvent` uses it when an agent event lacks `segment_id` (`internal/runstate/runstate.go:328-357`, `internal/runstate/runstate.go:432-440`). The previous unknown-agent first-event skip remains, but the agreed AF5 test and FINAL compatibility path for normal participants are covered.

### [MINOR] Initial partial stdout lines are split incorrectly after completion

Status from my prior MINOR AF4 finding: RESOLVED.

`loadFocusTail` and `readAppendedLines` now share `completeLinesFrom`; offsets advance only through complete newline-terminated regions. The new partial-line test covers `partial` followed by ` line\n`, so the fragmentation case is gone.

### [NIT] Overview preview still uses the old hard six-line tail

Status from my prior NIT AF9 finding: RESOLVED.

`refreshLogPreview` now derives a line budget from `previewLineBudget` and passes it to stdout/stderr tails instead of the hard-coded six-line preview.

## Open questions

- Does the team want `capFocusLines` to drop an oversized newest line entirely, or keep a truncated tail of that line? Either is compatible with the 4 MiB invariant, but it should be explicit.
- No objection from me to leaving F4 cross-process locking and slice 5 execution open for a later cycle; the remaining AF2 issue is the only agreed-fix blocker I see.
