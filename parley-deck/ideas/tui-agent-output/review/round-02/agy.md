---
agent: agy
idea: tui-agent-output
phase: review
round: 2
date: 2026-06-07
---

## Summary
The fixes applied in Fix-up Cycle 1 by the implementer (`claude`) successfully resolve all critical and major issues identified during Round 1 of the review. The TUI agent output is now split-safe for CRLF sequences, offset tracking avoids log chunk duplication, the steer events drain ordering ensures replies are not lost, and the artifact view displays the correct tail. The codebase remains decoupled, and the extensive test additions verify the robustness of these fixes.

## Verification

### 1. Ingest split `\r\n` line loss in `ingestTranscriptBytes`
* **Status**: FIXED
* **Detail**: A per-stream `crPending` state is now tracked in `agentBuffer` and threaded through `ingestTranscriptBytes`. Trailing carriage returns at chunk boundaries are deferred. If the subsequent chunk starts with `\n`, it is correctly committed as a newline, otherwise the carriage return is treated as a lone rewrite and the line is cleared. This is verified by `TestIngestSplitCRLFAcrossTicks` in `internal/tui/live_test.go`.

### 2. `readAppendedChunk` offset tracking duplication
* **Status**: FIXED
* **Detail**: `c.offset` is now advanced by the bytes actually read (`start + int64(len(data))`) instead of the pre-read stat size. This prevents duplicate reading of logs when new bytes are written after the stat check but before the read completes. This is verified by `TestReadAppendedChunkNoDuplication`.

### 3. `renderArtifactView` showing head of the tail, not the tail
* **Status**: FIXED
* **Detail**: Slicing logic was updated to slice `lines[len(lines)-rows:]`, ensuring that the newest lines (tail) of the bounded artifact view are rendered. The missing-file edge case is gracefully handled as verified by `TestArtifactViewMissingFile`.

### 4. Steer reply lost when `steer.replied` arrives same-tick
* **Status**: FIXED
* **Detail**: The `eventsMsg` handler in `Update` now executes `refreshBuffers()` (the drain) before calling `appendSteerEvents()`. Moreover, `appendSteerEvents()` calls `m.advanceBuffer(b)` and commits any trailing steer partial line before clearing the cursor. The test `TestSteerReplyTextAndMarkerBothKept` confirms both reply text and completion markers are kept.

### 5. UTF-8 cap split
* **Status**: FIXED
* **Detail**: Inside `ingestTranscriptBytes`, `capSeg()` has been modified to advance the start of the retained slice to a `utf8.RuneStart` boundary, eliminating any risk of mojibake. This is verified by `TestPartialCapIsRuneSafe`.

### 6. `relArtifact` relative `ideas/`
* **Status**: FIXED
* **Detail**: `relArtifact()` now correctly checks for the prefix `ideas/` (without a leading slash) and returns the correct context path `path[len("ideas/"):]`.

### 7. `/artifact` outside an agent tab error cleared
* **Status**: FIXED
* **Detail**: Command dispatching logic for `/artifact` handles the error state properly. If not inside an agent tab, it sets `m.inputErr = "open an agent tab to view its artifact"` without immediately clearing it, so it renders correctly.

### 8. Testing Gaps
* **Status**: FIXED
* **Detail**: Extensive model-driven test coverage was added in `internal/tui/live_test.go` covering all critical edge cases (including scroll follow-no-yank on partial updates via `TestScrollUpDisablesFollowNoYank`).

## New findings
None.

## Verdict
ACCEPT
