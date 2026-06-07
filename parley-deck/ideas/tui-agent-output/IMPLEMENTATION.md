---
idea: tui-agent-output
phase: implementation
status: complete
implementer: claude
date: 2026-06-07
---

# IMPLEMENTATION — live, scrollable, Codex-CLI-style agent transcript

Implements FINAL.md. Scope is `internal/tui/live.go` (+ `live_test.go`). No agent-flag
changes, no new deps; `internal/tui` stays decoupled. The runner/round/steer paths and
`--no-tui` are untouched.

## What was built
- **Transcript model.** `agentBuffer` rebuilt: `stdout`/`stderr`/`steer tailCursor`s;
  committed `lines []transcriptLine{Text, Stream}`; per-stream byte-capped live `partial`;
  `scroll`/`follow`/`trunc` + `hideStderr`/`showArtifact`. New `transcriptStream` enum.
- **CR-aware ingestion (headline).** `readAppendedChunk(*tailCursor)` advances the offset to
  EOF (bounded to `maxFocusBytes`; drops the leading partial on a jump), reports
  rotation/jump. `ingestTranscriptBytes` scans a chunk: `\r\n` and `\n` **commit** the
  partial (dropping empty/ANSI-only lines); a lone `\r` **rewrites** the live partial in
  place ("potom sa to prepíše"); committed lines are immutable; the partial is byte-capped.
  `cleanLogText` strips ANSI after CR handling. `capTranscriptLines` keeps the bound.
- **Two-stream tail + live partial.** `advanceBuffer` ingests stdout, stderr AND the
  in-flight steer cursor each tick; `visibleLines()` = committed (stderr filtered when
  hidden) + the live partials at the bottom. `bufferBottom` counts the visible lines minus
  the status-header row; follow pins to the live partial; scroll-up never yanks.
- **Status header (never blank).** `renderAgentStatusHeader` from `AgentState`:
  `● <id> working… M:SS` (+ `STALE`) / `✓ <id> finished … · wrote <rel-artifact>` /
  `✗ <id> failed …` / `◌ killed` / skipped / pending — replaces the old "no output yet".
- **stderr visible + `/stderr`.** stderr merged, rendered dimmed + `[err]`; `/stderr`
  toggles `hideStderr`.
- **Woven steer conversation.** The 1.18 replace-panel is GONE. A steer appends a cyan
  `❯ you: <text>` line into the agent's transcript and tails the reply stdout as a
  `transcriptSteer` stream (CR/partial-handled, streams in place); `appendSteerEvents`
  weaves a `[reply complete]` / `[reply failed]` marker on the terminal event. It scrolls
  with the rest.
- **Artifact view + `/artifact`.** Toggles the body to a bounded tail of
  `AgentState.ArtifactPath` with a `[Viewing Artifact: <rel>]` header; missing/empty →
  `[Artifact not yet written]`. `/stderr` + `/artifact` joined `commandSpecs` (autocomplete
  + `/help`); no new single-letter keys.
- **Agent flags UNCHANGED** (protects the stdout artifact-capture fallback). codex exec
  streams natively → it now shows live; one-shot agents become non-blank via status +
  stderr + artifact. claude `--output-format stream-json` is a deferred follow-up.

## Tests
- CR scanner cases (`a\rb`, `a\rb\n`, `a\r\nb\n`, multi-CR, `a\nb\rc`, split-across-ticks,
  ANSI-stripped, progress→final) — `TestIngestTranscriptCRCases`.
- Live partial surfaced + rewrites in place; stderr merged + `[err]` + `/stderr` hides;
  status header never blank (no "no output yet"); `/stderr`+`/artifact` toggle; steer woven
  (`❯ you:` line + steer cursor, NOT a replacing panel) + reply markers; file-replace
  pickup. Existing transcript scroll/follow + kill/steer/picker/suggest tests updated and green.

## Verification
`go build ./... && go vet ./... && go test ./...` green.

## Manual-smoke note
Not run in a live terminal (the model-driven tests exercise the ingester + render headlessly).
To smoke: `parley run`; while an agent works, its tab shows the status header + live stderr;
a streaming line (codex) rewrites in place; type a message → `❯ you:` + the reply stream into
the same scroll; `/stderr` hides stderr; `/artifact` shows the produced file.

## Fix-up cycle 1 (Phase 6 review → addressing codex + agy)
- **MAJOR/CRITICAL — split `\r\n` across ticks** (codex+agy): the ingester now carries a
  per-stream `crPending` (a `\r` at a chunk boundary is deferred and resolved against the
  next chunk's first byte) — `"alpha\r"`+`"\nb\n"` commits `alpha` then `b` instead of
  losing `alpha`. Test `TestIngestSplitCRLFAcrossTicks`.
- **MAJOR — `readAppendedChunk` offset duplication** (agy): the offset now advances by the
  bytes ACTUALLY read (`start + len(data)`), not the pre-read stat size, so bytes appended
  between Stat and ReadAll aren't re-read. Test `TestReadAppendedChunkNoDuplication`.
- **MAJOR — steer reply lost on terminal event** (codex): the eventsMsg handler now drains
  buffers BEFORE `appendSteerEvents`, which also drains + commits the trailing reply line
  before clearing the steer cursor — both the reply text and the `[reply complete]` marker
  are kept. Test `TestSteerReplyTextAndMarkerBothKept`.
- **MAJOR — artifact view showed the head, not the tail** (codex+agy): `renderArtifactView`
  now slices `lines[len-rows:]`.
- **MINOR — UTF-8 cap split** (codex): the partial byte-cap advances to a `utf8.RuneStart`
  boundary (no mojibake). Test `TestPartialCapIsRuneSafe`.
- **MINOR — `relArtifact` relative `ideas/`** (agy): handles a leading `ideas/` (no slash).
- **NIT — `/artifact` outside an agent tab** (codex): the error is no longer cleared.
- **Tests added** (codex+agy): split-CRLF, offset-no-dup, UTF-8 cap, steer text+marker,
  scroll-up-disables-follow-no-yank, /artifact missing-file graceful.

`go build/vet/test ./...` green. Ready for re-review.

## Deviations / scope (per FINAL)
- v1 strips ALL ANSI from log data (colour-preservation needs ANSI-aware width → deferred);
  the TUI's own lipgloss styling (status colours, `[err]`, `❯`) is intact.
- Per-tick stdout-then-stderr merge (no exact cross-stream interleave — separate files; a
  runner `combined.log` is a deferred option). On stdout file rotation the committed
  scrollback is kept (continuous transcript) and new content appends.
- claude `--output-format stream-json` (true live streaming for claude) = deferred follow-up
  behind a safe parser. Agent `HeadlessArgs` unchanged.
