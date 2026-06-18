---
idea: meta-protocol-change-fusion-execplans
status: final
author: claude
consensus-date: 2026-06-18
participants: [claude, codex, agy, hermes]
quorum-note: "claude ✅, codex ✅, hermes ✅; agy live-signoff waived per operator ruling (headless hang) — agy round-01 endorsement on record."
---

## Final plan / specification

A **meta-protocol-change** (§7) ratifying inspiration from OpenRouter Fusion +
OpenAI ExecPlans. All changes are **additive protocol-text guidance**; no Go logic
is required beyond keeping the embedded default in lockstep (drift guard). Apply each
edit **byte-identically** to both protocol copies in their shared (non-allowlisted)
zones:
- `parley-deck/COOPERATION.md` (live deck)
- `internal/protocol/defaults/COOPERATION.md` (embedded `parley init` default; stays
  genericized — only the allowlisted workspace/roster/transport zones differ, which
  these edits do not touch)

Guiding principle (ratified): **conditional rigor.** The new sections are REQUIRED
only for complex / `auto_implement` / driver-managed / pipeline ideas; trivial or
design-only ideas keep today's lean templates ("N/A" remains valid). This preserves
the "don't add overhead where it doesn't pay" rule.

### Edit 1 — Phase 4 (Finalization): self-contained, still-static `FINAL.md`
**Where:** §4 → Phase 4, the `FINAL.md` template + surrounding prose.
**What:** Extend the template with design-time sections (written before consensus
close; `FINAL.md` remains immutable after Phase 4 — the existing "open a v2" rule is
unchanged and reaffirmed). Add to the template:
- `## Purpose / user-visible outcome` — what is true after implementation.
- `## Context & orientation` — relevant paths, modules, terms, constraints found in
  deliberation.
- `## Observable acceptance criteria` — concrete commands / UI / API / artifact checks
  a reviewer or the driver can verify (Edit 3).
- `## Idempotence & recovery` — what state matters, what is safe to rerun, what needs
  a human gate (Edit 4).
- `## Known risks / de-risking` — spikes / high-risk dependencies to validate first.

Add one sentence: *"For complex, `auto_implement`, or driver-managed ideas, `FINAL.md`
plus `IMPLEMENTATION.md` MUST be self-contained enough that a fresh agent or the driver
can implement/resume from them alone, without session transcripts. For trivial or
design-only ideas these sections may be `N/A`."*

### Edit 2 — Phase 5 (Implementation): living `IMPLEMENTATION.md` sections
**Where:** §4 → Phase 5, the `IMPLEMENTATION.md` template.
**What:** Add living sections (kept current at every stopping point), conditional per
the rigor rule:
- `## Progress` — timestamped checklist, ISO `(YYYY-MM-DD HH:MMZ)`, partial-step
  notation `(completed: X; remaining: Y)`.
- `## Decision Log` — decisions made *after* `FINAL.md`: `Decision / Rationale /
  Date·Author` (deviations still also recorded under the existing "Deviations" head).
- `## Surprises & Discoveries` — unexpected findings with evidence, esp. when they
  change implementation choices.
- `## Validation evidence` — which acceptance criteria (Edit 1/3) were met, with the
  commands run and what they proved.
- `## Outcomes & Retrospective` — at completion: achievements, gaps, lessons; framed
  to feed §13 `parley retro` (Edit 6).

Note in prose: these make `IMPLEMENTATION.md` the **living** companion to the static
`FINAL.md`, giving the auto-drive driver task-level resume context (§12 already
supplies the low-level effects ledger + idempotency keys).

### Edit 3 — Phase 4/6 bridge: behavior-focused acceptance criteria
**Where:** Phase 4 template (the `## Observable acceptance criteria` head from Edit 1)
plus one sentence in Phase 6.
**What:** Acceptance criteria are stated as observable behavior (e.g. "after X, Y is
true"), not just intent. Phase 6 reviewers MAY classify findings against them; this
**does not change** the severities `CRITICAL/MAJOR/MINOR/NIT` — it makes severity
assignment less subjective. Phase 6 sentence: *"Where `FINAL.md` states observable
acceptance criteria, reviewers should check the implementation against them and may
cite a criterion in a finding."*

### Edit 4 — Idempotence & recovery (covered by Edit 1's `## Idempotence & recovery`)
No separate section; it is the Phase-4 head added in Edit 1. Prose ties it to the
driver: *"For `auto_implement` / action / pipeline ideas this section is required and
the driver treats it as the recovery contract."*

### Edit 5 — Phase 3 + Phase 7: advisory "compare-don't-merge" + blind-spots lens
**Where:** §4 → Phase 3 `consensus.md` template, and Phase 7 `review/consensus.md`
template.
**What (advisory drafting discipline, NOT a gate or a judge authority):**
- Phase 3 `consensus.md`: add `## Comparison & blind spots` with prompts —
  *contradictions* (don't smooth real disagreement into a vague trade-off),
  *partial coverage* (what only one participant covered), *unique insights* (minority
  contributions worth keeping), and **blind spots** (*"what did no participant
  address?"* — a forward-looking check before finalize).
- Phase 7 `review/consensus.md`: add `## Coverage & blind spots` — findings everyone
  saw vs. only one reviewer saw, and areas no reviewer inspected deeply.
- Explicit clause: this is advisory; append-only signoffs remain the only gate; any
  participant may block if the comparison is inaccurate; raw round files are never
  hidden behind a summary.

### Edit 6 — §13: "confident-error" retro signal
**Where:** §13 (Retrospective optimization), evidence corpus / guardrails.
**What:** Add that retro evidence SHOULD surface **confident-error** signals — a
dismissed `CRITICAL`/`MAJOR` finding, an unsupported assumption that shaped `FINAL.md`,
or a missed risk that caused fix-up churn — sourced from the new `IMPLEMENTATION.md`
Outcomes & blind-spots fields. Explicit clause: this is **diagnostic evidence only**,
**never** a new review severity, a blame label, or a merge gate. (The `parley retro`
tooling that mines it is a deferred follow-up, governed by §13 but specified
separately — out of scope for this idea.)

### Edit 7 — §3 directory layout annotations
**Where:** §3 layout comments.
**What:** Update the inline comments to reflect FINAL.md = "static, self-contained
authoritative artifact" and IMPLEMENTATION.md = "living execution document
(Progress/Decision Log/Surprises/Outcomes)". Cosmetic; keeps the map honest.

### Explicitly NOT in scope (ratified rejects)
Confidence-by-breadth / majority gates; a dedicated judge role with authority; hiding
raw rounds behind a judge summary; Fusion panel / recursion-depth / cost / default
web-search machinery; collapsing the deck into one ExecPlan file;
"proceed-without-prompting" autonomy across gates; the anti-list/anti-table prose
maximalism; and any new `parley retro` mining code (deferred follow-up).

### Acceptance criteria for THIS change (observable)
1. Both COOPERATION.md copies contain Edits 1–7, byte-identical in shared zones.
2. `go test ./internal/protocol/ -run TestEmbeddedDefaultMatchesLiveDeck` is green
   (drift guard holds).
3. `go test ./...` is green; `go build ./...` succeeds.
4. The embedded default remains genericized (no parley-deck roster/workspace leaked
   into `internal/protocol/defaults/COOPERATION.md`).
5. §7 footers added to any changed numbered section per convention are accurate; the
   change is reviewed (Phase 6) by ≥1 non-implementer and signed off (Phase 7).

### Idempotence & recovery (for this change)
Pure text edits under git; rerunnable safely. Recovery = `git checkout` the two .md
files. No external side effects. No `auto_implement`.

## References
- Consensus: ./consensus.md  (signoffs: claude ✅, codex ✅, hermes ✅; agy waived)
- Round-01: ./round-01/{claude,codex,agy,hermes}.md
- Research: ./reference/research.md
