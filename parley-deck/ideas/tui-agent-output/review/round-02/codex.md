---
agent: codex
idea: tui-agent-output
phase: review
round: 2
date: 2026-06-07
---

## Summary
ACCEPT. The fix-up cycle resolves the round-01 blocking issues I raised and the overlapping agy blockers. I did not find a new regression in the CR state machine, file tail offset handling, steer reply drain/marker ordering, artifact rendering, UTF-8 cap, relative artifact path display, or `/artifact` error handling.

## Verification
FIXED — split `\r\n` across ticks losing a line. `ingestTranscriptBytes` now carries `crPending` per stream and defers a trailing `\r` until the next byte. `"\n"` commits the existing segment as CRLF, while any non-newline clears the segment and reprocesses that byte as a lone-CR rewrite. This preserves the accepted scanner cases for same-chunk CRLF, split CRLF, repeated CR rewrites, and multi-chunk CR followed by non-newline. Rotation and jump reset the pending CR for only the affected stream.

FIXED — `readAppendedChunk` offset duplication. The cursor now advances to `start + int64(len(data))`, using the actual bytes read rather than the pre-read stat size. This avoids duplicating bytes appended between `Stat` and `ReadAll`, while still preserving the jump-ahead behavior and dropping only the leading fragment after the offset has been advanced to the actual read end.

FIXED — steer reply lost when `steer.replied` arrives same tick. The `eventsMsg` path refreshes buffers before appending steer terminal events, and `appendSteerEvents` drains the steer cursor again before committing any trailing partial and clearing the cursor. The reply text and `[reply complete]` marker are both retained, and I do not see a double-commit path for newline-terminated replies because the partial is empty after ingestion.

FIXED — artifact view showed the head, not the tail. `renderArtifactView` now slices `lines[len(lines)-rows:]`, so the visible artifact body shows the newest rows from the already bounded tail.

FIXED — UTF-8 cap split. `capSeg` now advances the retained byte slice to a `utf8.RuneStart` boundary after applying the byte cap, keeping the live partial valid UTF-8.

FIXED — `relArtifact` dropped relative `ideas/` context. `relArtifact` now handles paths beginning with `ideas/` as well as paths containing `/ideas/`, so relative artifact paths keep the idea slug context.

FIXED — `/artifact` outside an agent tab cleared its own error. The command now clears `inputText` but leaves `inputErr` set when there is no active agent tab.

FIXED — missing high-risk tests. The added tests cover split CRLF, offset no-duplication, steer reply text plus marker, rune-safe partial cap, scroll-up no-yank, and missing artifact rendering. These directly cover the round-01 blockers and the highest-risk fix-up contracts.

## New findings
None.

## Verdict
ACCEPT.
