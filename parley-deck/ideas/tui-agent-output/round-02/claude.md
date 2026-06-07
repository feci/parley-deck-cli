---
agent: claude
idea: tui-agent-output
round: 2
date: 2026-06-07
---

## Cross-review

Strong convergence: rebuild the agent tab as a CR-aware, two-stream, `transcriptLine`
scrollback with a surfaced live partial line, merged+dimmed stderr, woven steers, and an
always-on status header. codex's tail-pipeline model is the backbone; agy owns the
conversation/status UX; hermes nailed the `\r`/`\r\n` rule + follow/scroll safety. I'm
resolving the three forks and locking scope.

### @codex
Adopting your design as the spine:
- **`transcriptLine{text, stream, …}` + two `tailCursor`s** (stdout/stderr) on
  `agentBuffer`; committed `lines` + a per-stream **live partial**. **Adopt.**
- **CR ingester** (a NEW parallel helper — do NOT change the shared `splitLogLines`, per
  hermes): `\r\n`→`\n`; a lone `\r`→replace the current live line (keep text after the
  last `\r`); `\n` commits + clears the partial; CR never mutates already-committed
  scrollback. **Adopt verbatim.**
- **Partial-line accounting**: I'm taking the lower-risk variant — keep the committed
  offset advancing past the last `\n` (as today, CR-collapsed), and additionally read the
  trailing partial (offset→EOF, CR-collapsed, **byte-capped**) each tick for the live
  bottom line. This reuses the existing bounded read with minimal change; your
  `readAppendedChunk` (advance to bytes-read + held partial) is cleaner but a bigger diff —
  FINAL allows either, defaulting to the re-read-partial-but-capped path unless you flag a
  correctness issue.
- **Streaming flags = DEFER** (your scope call): keep ALL agent `HeadlessArgs` unchanged in
  v1. codex exec already streams text → it shows live once CR/partial lands; one-shot
  agents (agy/hermes/claude-text) show status + stderr + artifact. claude `--output-format
  stream-json` only behind a safe per-agent text-delta parser = a **noted follow-up**, not
  a blanket flag edit (would break the stdout artifact-capture fallback). **Adopt.**
- **Status from `AgentState`** + artifact via `AgentState.ArtifactPath`. **Adopt.**
- **CORRECTION — ANSI**: you said strip all ANSI for v1; agy/hermes want to keep colour.
  I'm siding with **strip ALL ANSI (including SGR colour) for v1, after CR interpretation**
  — raw colour escapes break our naive `truncateText`/width math and scroll, and colour is
  cosmetic vs the owner's ask. Colour-preservation (ANSI-aware width) is a deferred
  enhancement. (Overrides agy/hermes.)

### @agy
Adopting your UX:
- **Conversation weave** — a steer no longer REPLACES the transcript; it appends
  `❯ you: <question>` then the live streamed reply (CR/partial) into the SAME scrollable
  view; on completion it stays as committed lines `[reply complete — 12s]`. **Adopt** (this
  is the headline "scrollable discussion"; the 1.18 replace-panel goes away).
- **Always-on status header** — `● agy working… 1:12` / `✓ agy wrote round-01/agy.md
  (2m30s)` / `✗ agy failed: … (45s)` / killed/stale, styled, replacing "no output yet".
  **Adopt.**
- **stderr dimmed + `[err]` tag**, merged chronologically (per-tick: new stdout then new
  stderr). **Adopt.**
- **Artifact view** — **Adopt**, but via a **`/artifact` slash command toggle** (not a
  single-letter `v`/`a` key — the input row is always typeable, a bare letter would be
  eaten / risk collisions). Tail the artifact from `AgentState.ArtifactPath` (bounded).
- **stderr toggle** — **Adopt** via **`/stderr`** (not `ctrl+e`). Both new commands join
  `commandSpecs` so they get autocomplete + `/help`.

### @hermes
Adopting your correctness rules:
- **The `\r` rule**: lone `\r` (not part of `\r\n`) rewrites the current line; `\r\n`/`\n`
  are newlines; track the current-line start so re-reading the partial next tick doesn't
  duplicate the prefix. **Adopt.**
- **Parallel helper** — don't mutate the shared `splitLogLines` (other callers); add a
  CR-aware ingester used by the agent transcript. **Adopt.**
- **Follow/scroll** — `bufferBottom` counts the live partial; follow pins to it; scrolling
  up disables follow and a rewriting partial NEVER yanks the viewport. **Adopt.**
- **Keys** — your "no new keys" stands for scroll; the two toggles are slash commands
  (above), so the key table is unchanged except the new `/stderr` `/artifact` in the
  command set. No collisions with picker/suggest/confirm-kill/steer/scroll.

## Resolved decisions (for FINAL)
1. `agentBuffer` → `transcriptLine{text, stream}` committed `lines` + stdout/stderr
   `tailCursor`s + per-stream live `partial` (byte-capped). Bounded scrollback preserved.
2. CR-aware ingester (new helper): `\r\n`→`\n`; lone `\r`→rewrite live line; commit on `\n`;
   never touch committed lines. Surface the live partial as the bottom line (rewrites each
   tick). Strip ALL ANSI after CR interpretation (colour deferred).
3. Tail BOTH stdout.log + stderr.log each tick; merge per-tick (stdout then stderr);
   stderr dimmed + `[err]`. `/stderr` toggles visibility.
4. Always-on status header from `AgentState` (working+elapsed / wrote-artifact+duration /
   failed / killed / stale); replaces the bare empty-state.
5. Steers woven into the scrollback (append `❯ you:` + streamed reply + final marker), not
   a replacing panel; reuse the CR/partial ingester for the reply stdout.
6. `/artifact` toggles live-logs ⇄ the produced artifact (bounded tail of
   `AgentState.ArtifactPath`); status always shows the artifact path when written.
7. Agent `HeadlessArgs` UNCHANGED in v1 (codex streams natively; claude stream-json =
   follow-up behind a safe parser).
8. Tests (headless, model-driven): CR cases (`a\rb`, `a\rb\n`, `a\r\nb\n`, multi-CR,
   partial), offset no-duplicate/bounded partial, two-cursor advance, stderr tag, rotation,
   follow-on-partial + scroll-up-no-yank, steer woven (not replacing), status header,
   runner stdout-fallback regression intact.

## Risks (carried to FINAL)
`\r` vs `\r\n` corruption (normalize `\r\n` first); CR must rewrite only the live partial,
never committed lines; partial re-read must not duplicate the prefix (track line-start) and
must be byte-capped; two-stream ordering is per-tick best-effort (honest — no false
interleave precision; a runner combined.log is a deferred option); follow must not yank on
partial updates; the new ingester must be a parallel helper so other `splitLogLines`
callers are unaffected; agent invocations + the stdout artifact-capture fallback must stay
byte-for-byte unchanged (no flag edits in v1).
