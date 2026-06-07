---
agent: agy
idea: tui-agent-output
phase: review
round: 1
date: 2026-06-07
---

## Summary
The Codex-CLI feel is mostly delivered (status header always-on, stderr merged/dimmed, and steer conversation woven), but critical ingestion edge cases (split `\r\n` line loss, cursor offset race duplication, and reversed artifact tail display) require fixes.

## Findings
* **CRITICAL** [internal/tui/live.go:2268](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L2268) `ingestTranscriptBytes`: When a `\r\n` sequence is split across two read chunks (e.g. `\r` is the last byte of chunk 1, and `\n` is the first byte of chunk 2), the `\r` at the end of chunk 1 is incorrectly processed as a lone `\r`. This clears the partial line (`seg = seg[:0]`), meaning the line content before it is permanently discarded. On the next tick, the `\n` commits the empty buffer.
  *Fix*: Defer processing `\r` if it is the last byte of a chunk. Keep it at the end of the `partial` string, and in the next invocation of `ingestTranscriptBytes`, check if the chunk starts with `\n`. If it does, strip the `\r` and commit. If it does not, treat the trailing `\r` as a lone `\r` (clear `seg`).
* **MAJOR** [internal/tui/live.go:2235](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L2235) `readAppendedChunk`: The tail cursor's offset `c.offset` is set to the pre-read file size (`c.offset = size`). If the target process writes new bytes to the log file after `f.Stat()` but before `io.ReadAll(f)` completes, those bytes are read in the current tick, but `c.offset` is set only to the older, smaller `size`. The next tick seeks to `size` and re-reads those bytes, causing transcript duplication.
  *Fix*: Set `c.offset` to the actual post-read position. Record the read length before slicing/editing the data: `c.offset = start + int64(len(data))`, or query the file offset via `f.Seek(0, io.SeekCurrent)` right after reading.
* **MAJOR** [internal/tui/live.go:789](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L789) `renderArtifactView`: When the number of loaded lines from the artifact tail exceeds the available window size (`rows`), the view slices from the beginning (`lines[:rows]`) rather than the end. This displays the oldest lines (head) of the loaded tail rather than the newest lines (tail), showing outdated info.
  *Fix*: Slice from the end of the lines slice:
  ```go
  if len(lines) > rows {
      lines = lines[len(lines)-rows:]
  }
  ```
* **MINOR** [internal/tui/live.go:804](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L804) `relArtifact`: The code checks `strings.LastIndex(path, "/ideas/")`. If the artifact path is relative and starts with `ideas/` (without a leading slash), this check returns `-1` and falls back to the base name (e.g. `artifact.md`) rather than keeping the slug context (e.g. `tui-agent-output/artifact.md`).
  *Fix*: Add a fallback to check for `ideas/` prefix:
  ```go
  if i := strings.LastIndex(path, "/ideas/"); i >= 0 {
      return path[i+len("/ideas/"):]
  } else if strings.HasPrefix(path, "ideas/") {
      return path[len("ideas/"):]
  }
  ```
* **MINOR** [internal/tui/live_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live_test.go): There are testing gaps for crucial specification guarantees:
  * No test verifies that manual scrolling up disables follow-mode and prevents rewriting partial lines from yanking the viewport.
  * No test verifies that file rotation resets the `partial` state *only* for the affected stream.
  * No test verifies that `/artifact` handles a missing file gracefully.
  *Fix*: Add test cases covering these scenarios in `live_test.go`.

## Verdict
REQUEST-CHANGES

**Blocking items for ACCEPT:**
1. Fix the split `\r\n` line loss in `ingestTranscriptBytes`.
2. Fix the file offset tracking in `readAppendedChunk` to prevent duplicate log chunks.
3. Fix the slicing logic in `renderArtifactView` to display the actual tail of the artifact.
