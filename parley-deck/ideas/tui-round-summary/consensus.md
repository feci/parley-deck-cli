---
idea: tui-round-summary
drafted-by: claude-1
date: 2026-07-04
track: standard
participants: [claude-1, codex-1, hermes-1]
---

## Agreed decisions

Strong convergence. Two round-01 refinements are adopted as binding: (a) build the
digest in the driver/event layer, never in the Home renderer (codex — single source of
truth); (b) render stance as keyword *flags/hints*, not ACCEPT/counter/block *verdicts*
(hermes — free text can't be classified reliably without an LLM, which is a v1 non-goal).

1. **Source of truth = the driver's completion gate.** The digest is built only after
   the driver's `roundComplete(round)` predicate passes (which already validates one
   artifact per participant + round-2+ cross-review evidence + terminal-event
   reconciliation). Home NEVER decides completion from file presence — that would
   diverge from the advancement gate and can mislead on late writes.

2. **Emit a presentation event** `round.digest` keyed by `(idea, round)`, appended to
   the run event stream after the round-complete action, so it is durable (survives
   reopen), replayable, and consumed like existing `round.completed` / `run.phase`
   events. Emission is **idempotent** — a stable per-(idea,round) key; re-running the
   driver never appends a duplicate.

3. **Digest content** (position map, not argument summary):
   - Header `Round NN — complete (X/Y)` — completeness is the highest-value signal.
   - One line per agent: first sentence/paragraph under `## Summary`, cap ~120 chars at
     a sentence boundary. No `## Summary` ⇒ fall back to first prose paragraph AND tag
     it `[no Summary — fell back]` so degraded extraction is visible, not authoritative.
   - `flags:` line = raw keyword counts over a small curated set
     (`block`/`blocker`, `counter-proposal`, `accept`/`agree`, `escalat`), presented as
     **hints/mentions, never verdicts**. UI wording must make the hint nature unmissable.
   - `next:` line straight from driver phase state (opening round N+1 / consensus
     drafting / re-round) — the most robust signal, touches no prose.
   - Round 2+: derive engagement gaps from the structured `responding-to:` frontmatter
     ("@codex-1 did not address @hermes-1") — deterministic convergence signal.

4. **Extraction is per-agent-degrading and never blocks advancement.** An extraction
   failure for one agent yields a tagged fallback line; it must NOT be treated as a
   driver failure (a display feature cannot block protocol progress).

5. **Rendering** — extend `ProtocolSnapshot` with recent round digests; `renderHome`
   stays a pure renderer over model/snapshot state (no disk scan, no markdown parse per
   render). History bounded by a **viewport budget** (scrollable sub-region), NOT a
   fixed count, so chips/roster/runs list are never pushed off-screen (regression guard).

6. **No LLM, no protocol change, read-only over canonical files.**

## Deferred follow-ups

- A structured `stance:` frontmatter field agents set (accept/counter/block/question)
  would let the digest show real verdicts — a protocol change, deferred.
- Navigable one-liners (enter → jump to that agent's tab) so a truncated position is
  recoverable in one keystroke (hermes) — nice-to-have, deferred unless cheap.

## Dismissed / non-goals

- No stance verdict classification in v1 (flags/hints only).
- No cross-round analytics, no LLM executive summaries, no per-agent timestamps.
- No change to per-agent tabs or transcript view.

## Signoffs

<!-- each participant appends its own block -->

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Drafter. Adopts driver-layer digest + round.digest event + flags-not-verdicts +
viewport-bounded history. Read-only, no protocol change.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
This matches my round-01 position: the digest belongs in the driver/event-state layer, keyed to the canonical completion gate, with Home rendering cached snapshot state only. I also accept the idempotent `round.digest` event keyed by `(idea, round)` so reruns do not duplicate history.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
Matches my round-01: stance as keyword flags/hints (not verdicts) avoids unreliable free-text classification, and viewport-bounded history (scrollable sub-region, not a fixed count) keeps chips/roster/runs list visible. Read-only, no protocol change.
