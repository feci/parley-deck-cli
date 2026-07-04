---
agent: hermes-1
idea: tui-round-summary
round: 1
date: 2026-07-04
---

## Summary

A returning user needs a *position map*, not an argument summary. In 5 seconds they
want to answer four questions: which round just finished, did every expected agent
contribute, where do agents stand relative to each other, and what happens next. The
per-agent tabs already hold the full reasoning; the digest is an index that tells the
user *whether* to drill in, not a substitute for drilling in. The strongest
mechanically-extractable signals are structural (file/section presence, frontmatter
`responding-to`, driver phase state); the weakest is semantic stance classification
(ACCEPT/counter/block), which free-text round files do not support reliably without an
LLM. The v1 digest should lean on structure and present semantic hints as *flags*, not
verdicts.

## Proposed approach

Digest block shape (fixed, scannable, rendered in Home tab when the driver declares a
round complete):

```
Round 03 — complete (4/4)
──────────────────────────────────────────────────────
@claude-1  Accepts v2 framing; proposes splitting extraction in two stages.
@codex-1   Raises blocker: completion detection races with late writes.
@hermes-1  Digest scope = position map, not argument summary.
@agent-4   [no Summary — fell back to first paragraph]

flags: 1 blocker · 2 counter-proposal · 0 accept-mention
next: opening round 04
```

Design rules:

1. One line per agent. Extract the first sentence/paragraph under `## Summary`; cap at
   ~120 chars and truncate at the nearest sentence boundary. If `## Summary` is absent,
   fall back to the first prose paragraph after frontmatter and tag it
   `[no Summary — fell back to first paragraph]` so the user sees the extraction was
   degraded, not authoritative. Honesty about extraction quality is part of information
   design.

2. Completeness header (`Round N — complete (X/Y)`) so a missing agent is visible at a
   glance. This is the single highest-value signal for a returning user.

3. `flags:` line = raw keyword counts over a small curated set
   (`block`/`blocker`, `counter-proposal`, `accept`/`agree`, `escalat`). Presented as
   hints, never as classified stances. The user interprets them; the tool does not claim
   it understood intent. This is the honest v1 substitute for the "ACCEPT/counter/block
   counts" in the prompt.

4. `next:` line comes straight from driver phase state (round complete → opening
   round N+1; final round → consensus drafting; blocker present → re-round). This is the
   most robust signal in the whole digest because it does not touch file prose at all.

5. Engagement structure (round 2+): the `responding-to:` frontmatter field is
   structured YAML. Mechanically derive who-engages-whom and surface gaps ("@codex-1 did
   not address @hermes-1") — a stronger convergence signal than keyword counting and
   fully deterministic.

6. History: keep the last 3 round digests scrollable above the current one. More than 3
   becomes noise for a 5-second catch-up; older rounds live in per-agent tabs.

Explicitly noise for v1 (do not render): full prose paragraphs, verbatim quotes longer
than the one-liner, cross-round analytics, per-agent timestamps, risk-list details,
attempted tone/strength-of-disagreement classification.

## Concerns / open questions

- The prompt asks for "counts of ACCEPT/counter-proposal/blockers where detectable." Free
  text does not make these reliably detectable without an LLM. My counter-proposal: v1
  renders keyword *flags* (hints), not stance *verdicts*. A cleaner long-term fix is a
  structured `stance:` frontmatter field agents set (accept/counter/block/question), but
  that is a protocol change and out of scope for v1. Open: does the group accept flags as
  the v1 substitute, or does someone have a deterministic classifier I'm not seeing?

- Keyword flags false-positive on negation ("this is *not* a blocker"). Mitigation is
  presentation (label them "mentions", not "findings") plus drill-in. Is that acceptable,
  or do we suppress flags inside negated clauses? Suppressing negation mechanically is
  fragile.

- Truncating a nuanced Summary to one line can misrepresent a position. Should the
  one-liner be navigable (enter/jump to that agent's tab) so misrepresentation is
  recoverable in one keystroke? I believe yes.

- "Last 3 digests" is a guess. The right number depends on terminal height and how much
  existing Home-tab content (chips, roster, runs list) the digest competes with. Needs a
  viewport budget, not a fixed count.

- Round-completion detection racing with late writes is codex-1's lens, but it affects
  information design: a digest built on an incomplete file set misleads. The digest must
  render only when the *driver* declares the round complete, never merely because files
  appeared on disk.

## Risks

- Stance flags presented as counts may be read as authoritative verdicts by users who
  skip the "flags/hints" framing, leading to wrong mental models of convergence. UI
  wording must make the hint nature unmissable.

- Summary-section convention drift: if agents write verbose multi-paragraph summaries,
  the one-liner becomes a multi-liner and the digest grows past the 5-second budget. Hard
  cap + truncation handles it, but the convention should be documented.

- History of N digests grows the Home tab vertically; on small terminals it pushes chips,
  roster, and runs list off-screen, regressing existing content (an explicit
  non-goal). Needs a bounded viewport with the digest as a scrollable sub-region, not
  unbounded append.

- Fallback extraction (no `## Summary`) pulling the "first prose paragraph" can grab a
  boilerplate line or a `responding-to`-style preamble instead of substance, producing a
  misleading one-liner with no quality tag beyond the generic fallback notice.

- Engagement-structure signal (`responding-to`) only exists from round 2 onward; round 1
  has no cross-review graph, so the `flags:` line carries the whole semantic load in
  round 1 — its weakest round.
