---
idea: tui-agent-output
phase: consensus
drafter: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
---

## Consensus

After round-01 (4 positions) and round-02 (cross-review) all four participants agree.
codex's tail-pipeline model is the backbone; agy owns the conversation/status UX; hermes
the `\r` correctness + follow/scroll safety. codex's round-02 counter (lock the exact
partial-line cursor accounting) is adopted; agy's 3 minor UX refinements are folded in.

### Goal
Make the `parley tui` agent tab a **live, scrollable, Codex-CLI-style transcript**: a
streaming line that **rewrites itself in place** (honor `\r`), stderr merged-and-dimmed so
silent agents are visible, steers woven into the same scroll, and an **always-on status
header** so the tab is never blank. Today the TUI tails only stdout, strips `\r`, never
shows stderr — so the tab is blank during a run.

### A. Transcript model (rebuild `agentBuffer`)
```go
type transcriptStream int
const ( transcriptStdout transcriptStream = iota; transcriptStderr; transcriptSteer; transcriptEvent )
type transcriptLine struct { Text string; Stream transcriptStream; Live bool }
type tailCursor struct { Path string; Offset int64; Info os.FileInfo; Loaded bool }
// agentBuffer: stdout/stderr tailCursors; committed lines []transcriptLine (bounded);
// partial map[transcriptStream]string (byte-capped live line per stream); scroll/follow/trunc.
```

### B. CR-aware ingestion (the headline — new helper, do NOT touch shared `splitLogLines`)
- `readAppendedChunk(cursor) (chunk []byte, rotated bool, err error)` advances
  `cursor.Offset` to **EOF** for bytes actually read (bounded to `maxFocusBytes`; on a
  huge burst, jump to the last `maxFocusBytes`). No re-reading the same bytes.
- `ingestTranscriptBytes(lines, partial, stream, chunk) (lines, partial)` scans the chunk:
  - `\r\n` and `\n` → **commit** the current partial as a `transcriptLine{Stream}` and
    clear it; lone `\r` → **replace** the live partial with the text after it (rewrite in
    place); plain bytes → append to the partial.
  - Committed lines are **immutable** — a later `\r` only rewrites the live partial.
  - Strip ANSI **after** CR tokenization (or with a stripper that preserves `\r`/`\n`);
    drop a line that is empty after CR-collapse (no blank clutter). The byte-capped partial
    keeps the tail of an oversized newline-less line.
- **Exact scanner cases (must test):** `a\rb`→partial `b`; `a\rb\n`→commit `b`;
  `a\r\nb\n`→commit `a` then `b`; `a\rb\rc`→partial `c`; `a\nb\rc`→commit `a`, partial `c`;
  a partial split across ticks renders the same as if it arrived in one tick; a
  rotation/truncation reset clears ONLY that stream's partial + offset, not the other's.
- **ANSI scope:** strip all ANSI from the raw **log data** (colour + cursor) for v1
  (colour-preservation needs ANSI-aware width → deferred). The TUI's OWN lipgloss styling
  (status colours, `[err]` dim, `❯` prefix) stays rich — stripping applies only to parsed
  log bytes.

### C. Two-stream merge (stderr visible)
`ensureBuffer` sets `stdout.Path = a.StdoutPath`, `stderr.Path = a.StderrPath`.
`refreshBuffers` advances BOTH cursors each tick (reuse the rotation/truncation check per
cursor) and ingests stdout then stderr (deterministic per-tick order; separate files
cannot prove true interleave — no false precision; a runner `combined.log` is a deferred
option). stderr lines render **dimmed + `[err]`-tagged**. `/stderr` toggles stderr
visibility (filters `Stream==stderr`).

### D. Conversation weave (steers in the scroll, not a panel)
The 1.18 `steerReplies` REPLACE-panel is removed. A steer appends into the SAME scrollback:
a `❯ you: <question>` line (`transcriptSteer`), the live streamed reply (via the same
CR/partial ingester over the steer attempt's stdout), and a final `[reply complete — 12s]`
/ `[reply failed]` marker. It scrolls with the rest; on completion it stays as committed
lines. (Reuse the steer stdout path already returned by `SubmitSteer`.)

### E. Always-on status header (never blank)
A styled line below the tab strip, OUTSIDE the scrollback, from `AgentState`:
- `● <agent> working… 1:12` (warn) [+ `STALE` when the liveness seam says so];
- `✓ <agent> wrote round-01/<agent>.md (2m30s)` (ok) — relative artifact path, readable;
- `✗ <agent> failed: <error> (45s)` (danger); `◌ <agent> killed` (muted).
Replaces the bare "no output yet from <agent>".

### F. Artifact view (`/artifact` toggle)
Toggles the viewport between live logs and the produced artifact (bounded tail of
`AgentState.ArtifactPath`), with a `[Viewing Artifact: <path>]` header; a missing/empty
file shows `[Artifact not yet written]` (graceful, no crash/fd leak). `/stderr` +
`/artifact` join `commandSpecs` (autocomplete + `/help`); no new single-letter keys.

### G. Follow / scroll safety
`bufferBottom` and the rendered length count committed lines + visible live partials +
woven steer lines (status header is outside scrollback). `follow` pins the viewport to the
live partial; scrolling up disables follow and a rewriting partial NEVER yanks the
viewport.

### H. Agent flags — UNCHANGED in v1
No `HeadlessArgs` edits (protects the stdout artifact-capture fallback + validation). codex
exec streams text natively → it shows live once the tailer lands; one-shot agents
(agy/hermes/claude-text) become non-blank via status + stderr + artifact. claude
`--output-format stream-json` is a **deferred follow-up** behind a safe per-agent
text-delta parser that cannot break stdout capture.

### I. Tests (headless, model-driven)
The exact CR scanner cases (B); no prefix duplication + bounded partial across ticks; two
cursors advance independently; per-stream rotation/truncation; stderr `[err]` tag +
`/stderr` filter; follow-pins-partial + scroll-up-no-yank; steer woven (NOT replacing) with
a live reply; status header per state; `/artifact` toggle + missing-file graceful; runner
stdout-fallback + validation regression intact.

### Non-goals
No ANSI colour passthrough (v1), no exact cross-stream interleave, no agent-flag changes,
no live interactive chat beyond steering, no mouse. claude stream-json = follow-up.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-07. As drafter I accept. Locks the round-02 resolution: transcriptLine scrollback
+ two tailCursors + byte-capped per-stream live partial; the CR-aware ingester
(`readAppendedChunk` to EOF + `ingestTranscriptBytes`, exact scanner cases, never mutate
committed lines, strip-ANSI-after-CR); stderr merged+dimmed with `/stderr`; woven steers
(no replace-panel); always-on status header; `/artifact` toggle; agent flags unchanged.
Ready for FINAL + implementation.

### codex — ✅ ACCEPT
I accept the drafted consensus as written. It captures the v1 implementation boundary I asked for: exact cursor accounting, isolated CR-aware ingestion, immutable committed lines, byte-capped live partials, stderr filtering, woven steers, follow-safe scrolling, `/artifact`, and unchanged `HeadlessArgs`. Ready for FINAL + implementation.

### agy — ✅ ACCEPT
I accept the drafted consensus. It fully folds in my round-02 UX refinements—specifically the artifact-mode header, relative paths, and dropping empty-after-CR lines—along with codex's cursor accounting and hermes's follow/scroll protections. The design is clean, well-bounded, and ready for FINAL and implementation.

### hermes — ✅ ACCEPT
I accept the drafted consensus. All my round-02 points (CR-aware ingester isolation, exact scanner cases, follow/scroll safety, immutable committed lines, byte-capped partials, stderr merge+dim, woven steers, artifact toggle, status header, unchanged HeadlessArgs) are folded in. The spec is precise, bounded, and ready for FINAL + implementation.
