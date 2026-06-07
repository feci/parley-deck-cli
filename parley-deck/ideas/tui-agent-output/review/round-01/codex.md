---
agent: codex
idea: tui-agent-output
phase: review
round: 1
date: 2026-06-07
---

## Summary
REQUEST-CHANGES: the implementation is close, but terminal steer events can drop the whole reply stream, and the CR ingester is not correct when a CRLF pair is split across tail ticks.

## Findings

MAJOR, internal/tui/live.go:405 and internal/tui/live.go:1533, steer reply stdout is cleared before it is drained. `Update(eventsMsg)` calls `appendSteerEvents` before `refreshBuffers`; `appendSteerEvents` then sets `b.steer = tailCursor{}` and clears `partial[transcriptSteer]`. If the `steer.replied` event arrives in the same event read as the reply file becoming available, the TUI appends `[reply complete]` but never tails the reply stdout, so the woven conversation loses the agent's answer. Fix by advancing/draining the steer cursor before clearing it, or by refreshing buffers before appending terminal markers; add a test that submits a steer with a populated reply log, delivers `steer.replied`, and asserts the reply text plus marker are both in the transcript.

MAJOR, internal/tui/live.go:2253, CRLF handling is not split-safe across ticks. The scanner treats any `\r` not followed by `\n` in the same chunk as a lone rewrite, so chunks `[]string{"a\r", "\nb\n"}` drop `a` instead of committing it as the spec's `\r\n` newline case. Tail reads can split anywhere, including between `\r` and `\n`. Fix by carrying pending-CR state per stream, or by encoding that state in the partial/cursor model so a CR at chunk end is resolved against the next byte before deciding whether it is rewrite or newline; add split-CRLF tests.

MINOR, internal/tui/live.go:2281, partial byte capping can split UTF-8. `seg = seg[len(seg)-partialMaxBytes:]` may start in the middle of a multibyte rune, and the next `string(seg)` can render replacement characters in the live line or committed line. The spec explicitly called out UTF-8 safety. Fix by moving the retained start forward to a `utf8.RuneStart` boundary, or by storing/limiting runes rather than raw byte slices; add a test with a multibyte character crossing the cap boundary.

MINOR, internal/tui/live.go:788, `/artifact` does not show the bottom of the bounded artifact view when the artifact tail has more lines than the terminal body. `loadFocusTail` already returns the bounded tail, but `renderArtifactView` then keeps `lines[:rows]`, which shows the oldest rows in that bounded tail instead of the newest rows. Fix by slicing `lines[len(lines)-rows:]` and add a small render test with `rows` smaller than the artifact tail.

MINOR, internal/tui/live_test.go:1147, the tests cover the advertised simple CR examples and basic stderr/status/steer toggles, but they miss the highest-risk contracts: split `\r\n`, terminal steer event draining, per-stream rotation reset, follow-no-yank while a partial rewrites, `/artifact` missing-file and bounded-tail behavior, and UTF-8 cap safety. Add focused model tests for those cases before accepting this as the contract.

NIT, internal/tui/live.go:1787, `/artifact` outside an agent tab sets `m.inputErr` and immediately clears it via `m.inputText, m.inputErr = "", ""`, so the intended error is never shown. Keep the error when there is no active agent tab.

Correct points: `internal/tui` remains decoupled from `internal/runner` and `internal/app`; stdout, stderr, and steer cursors are independent; `readAppendedChunk` advances offsets to EOF and drops a jumped leading fragment; committed lines are not mutated by CR rewrites; `cleanLogText` strips ANSI after CR processing; stderr is merged and dimmed with `[err]`; `/stderr` and `/artifact` are wired into command specs and command dispatch; the old replace-panel steer reply path is gone; short terminal handling and the always-on status header are directionally correct.

## Verdict
REQUEST-CHANGES with blocking items: fix the steer terminal-event drain ordering and split-across-ticks CRLF ingestion before acceptance.
