---
agent: codex-1
idea: tui-round-summary
review-round: 1
date: 2026-07-04
reviewed-commit: 947effa
---

## Summary

The driver-side placement is mostly correct: `advanceRound` calls `roundComplete` before `emitRoundDigest`, digest append failures are swallowed, and the event blob contains the idea, round, per-agent lines, mention counts, and next action needed for a returning user. The documented deviations from `ProtocolSnapshot` and last-N history are acceptable for v1 because the TUI still renders from the durable run event stream and the latest digest gives the catch-up surface without adding snapshot plumbing.

Refutation checks run: `go build ./...` passed; exact `go vet` failed because the repository root has no Go files; `go vet ./...` passed; `go test ./internal/driver ./internal/tui` passed. I also traced malformed digest handling: invalid JSON is skipped by `latestRoundDigest`, so decode failure does not crash rendering.

## Findings

### [MAJOR] Home can still let the digest push Recent runs off-screen

`renderRoundDigest` is internally capped, but `renderHome` inserts the digest before `Recent runs` using only `rows/3` as the budget (`internal/tui/live.go:1097`-`internal/tui/live.go:1105`), then appends `Recent runs` afterward (`internal/tui/live.go:1108`-`internal/tui/live.go:1123`). The outer renderer clips the whole Home body at `rows` (`internal/tui/live.go:652`), so on a short terminal with enough open ideas, the digest can be the extra block that causes `Recent runs` to be clipped. For example, with a 10-row Home body and four ideas, the pre-digest layout can still show `Recent runs`; adding the 3-line digest consumes the remaining rows before the runs header.

Why it matters: FINAL.md explicitly requires the history/digest budget to avoid pushing chips/roster/runs off-screen. This implementation bounds the digest block, but it does not reserve rows for the content below it.

Suggested fix: compute the Home layout from reserved sections: render Ideas/chips, reserve at least the `Recent runs` header plus one row and the command hint, then allocate the digest from the remaining rows. Add a regression test that renders Home with a digest, several ideas, a small height, and asserts `Recent runs` is still present.

### [MINOR] Long-position truncation can split UTF-8 runes

`firstSentence` caps long text with byte indexing (`internal/driver/digest.go:105`-`internal/driver/digest.go:106`). A summary whose 120-byte boundary lands inside a multi-byte character will produce an invalid UTF-8 string before JSON marshaling replaces it. FINAL.md called out unicode-long-line coverage, but `TestExtractPositionCapsLongLine` only uses ASCII repeated words (`internal/driver/digest_test.go:71`-`internal/driver/digest_test.go:79`).

Why it matters: this is a display feature, so mojibake is not a protocol blocker, but the digest should preserve valid UTF-8 for non-English text, names, and symbols.

Suggested fix: cap by runes, not bytes, and add a test with a long Unicode summary where the boundary crosses multi-byte characters.

## Open questions

- Should `go vet` in project docs and review prompts be standardized as `go vet ./...`? The exact command currently fails at the repo root because there are no root-level Go files, while package vetting succeeds.
- Should last-N digests and round-2 engagement gaps be opened as a follow-up idea now, or remain deferred until users validate the latest-only Home surface?
