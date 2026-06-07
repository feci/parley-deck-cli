---
idea: tui-agent-output
phase: review-consensus
drafter: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
---

## Review consensus

Phase 6 (round-01) → fix-up cycle 1 → Phase 8 re-review (round-02). **All three reviewers
(codex, agy, hermes) ✅ ACCEPT at round-02 with zero remaining findings.** Ready to mark
complete.

### Agreed fixes (fix-up cycle 1 — all applied, verified FIXED in round-02)
- **MAJOR/CRITICAL — split `\r\n` across ticks** (codex + agy): per-stream `crPending`
  defers a chunk-boundary `\r` and resolves it against the next byte — no lost line.
  `TestIngestSplitCRLFAcrossTicks`.
- **MAJOR — `readAppendedChunk` offset duplication** (agy): offset advances by bytes
  actually read (`start + len(data)`), not the stat size. `TestReadAppendedChunkNoDuplication`.
- **MAJOR — steer reply lost on same-tick terminal event** (codex): drain buffers before
  `appendSteerEvents`, which drains + commits the trailing reply line before clearing the
  cursor — reply text AND marker kept. `TestSteerReplyTextAndMarkerBothKept`.
- **MAJOR — artifact view showed head not tail** (codex + agy): `lines[len-rows:]`.
- **MINOR — UTF-8 cap split** (codex): cap advances to a `utf8.RuneStart`. `TestPartialCapIsRuneSafe`.
- **MINOR — `relArtifact` relative `ideas/`** (agy); **NIT — `/artifact` outside a tab
  cleared its error** (codex).
- **Tests added** (codex + agy): split-CRLF, offset-no-dup, steer text+marker, rune-safe
  cap, scroll-up-no-yank, /artifact missing-file.

### Deferred follow-ups (non-blocking, out of v1)
- ANSI colour passthrough (needs ANSI-aware width); exact cross-stream interleave (a runner
  `combined.log`); claude `--output-format stream-json` (true live streaming for claude)
  behind a safe parser. Agent `HeadlessArgs` unchanged in v1.

### Verification
`go build ./... && go vet ./... && go test ./...` green; the CR scanner + streaming behavior
are covered by model-driven tests (no terminal). `internal/tui` imports neither runner nor
app; the runner round/steer/artifact-capture paths are untouched.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-07. As implementer I confirm fix-up cycle 1 is applied and green (build/vet/test).
All agreed fixes done with tests; no outstanding agreed fix. Ready to mark complete and release.

### codex — ✅ ACCEPT
The agreed fix-up cycle 1 changes are in, verified by round-02 re-review, and my round-01 findings are resolved. I accept completing and releasing this implementation.

### agy — ✅ ACCEPT
I confirm that all agreed fixes for my round-01 findings—including the split-CRLF newline tracking and offset tracking duplication—have been successfully verified in round-02. The codebase shows excellent stability with the added test suite, and I accept completing and releasing this implementation.

### hermes — ✅ ACCEPT
All round-01 findings are fixed and round-02 acceptance is recorded. The agent tab is now a live scrollable Codex-CLI-style transcript with streaming rewrite via \r, merged+dimmed stderr, woven steers, always-on status header, and /stderr + /artifact—all green. I accept completing and releasing.
