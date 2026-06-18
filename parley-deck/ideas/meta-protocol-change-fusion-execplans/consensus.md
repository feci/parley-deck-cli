---
idea: meta-protocol-change-fusion-execplans
phase: 3-consensus (facilitator draft — for human review)
drafted-by: claude (facilitator)
date: 2026-06-18
participants: [claude, codex, agy, hermes]
status: consensus-signoff (Phase 3; human approved ratification path)
---

# Consensus draft — Fusion + ExecPlans inspiration for parley-deck

All four participants (claude, codex, agy, hermes) wrote independent round-01
analyses and **converged strongly**. This draft synthesizes them for human review.
**Nothing is being changed in the protocol.** Per §7, any real change needs its own
ratified meta-protocol-change idea + human approval; this idea only answers *"which
concepts could inspire us, and which not."*

## Headline finding

parley-deck already owns the hard part of **Fusion** in a stronger, auditable form
(named participants, independent round-01 files, cross-review, `consensus.md`,
append-only signoffs, human gate). The real prize is **ExecPlans**, adapted *into the
existing artifact chain* — not as a replacement. The unanimous shape: **borrow
artifact-shape discipline (FINAL.md / IMPLEMENTATION.md / consensus.md), reject
mechanism imports (hidden judge, confidence gates, panels, autonomy-across-gates).**

## Agreed — ADOPT / ADAPT (unanimous or near-unanimous)

### A. ExecPlan living sections in `IMPLEMENTATION.md` — *highest leverage*
All four. Make `IMPLEMENTATION.md` a self-contained *living* execution document so a
fresh headless agent **or the auto-drive driver** can resume from the artifact alone:
- **Progress** — timestamped checklist (ISO `(YYYY-MM-DD HH:MMZ)`), partial-step notation.
- **Decision Log** — decisions made *after* `FINAL.md`, with rationale + date/author.
- **Surprises & Discoveries** — with evidence, especially when they change choices.
- **Validation evidence** — which acceptance criteria were satisfied, with what commands.
- **Outcomes & Retrospective** — at completion, framed to feed §13 `parley retro`.

Two-birds: directly fixes cross-invocation **resumption** (the driver's real need —
§12 already gives the driver an effects ledger + idempotency keys; the missing layer
is *task-level* orientation/recovery narrative), **and** hands `parley retro` a far
richer evidence corpus than today's thin density metrics.

### B. `FINAL.md` stays **static**, but gets self-contained design-time sections
**Resolved tension:** hermes initially framed FINAL.md as "living"; agy and codex
(and on reflection claude) hold that **FINAL.md must remain immutable after Phase 4**
— it is the consensus snapshot / audit trail. If invalidated → open a v2 idea, do not
mutate it. The ExecPlan *self-containment* principle still applies, written at design
time (so still static):
- **Purpose / user-visible outcome**, **Context & orientation** (paths, modules, terms).
- **Observable acceptance criteria** (concept C below).
- **Idempotence & Recovery notes** (concept D below).
- **Known risks / de-risking** (spikes, high-risk dependencies to validate first).

### C. Behavior-focused acceptance criteria — bridge Phase 4→5→6→driver
All four. `FINAL.md` should state *how success is observed* (commands, UI checks, API
behavior, artifact checks), not just intent. `IMPLEMENTATION.md` records which were
met + evidence; Phase 6 reviewers classify findings against them; the driver can
verify them. **Does not change** the review severities (CRITICAL/MAJOR/MINOR/NIT) —
it makes severity assignment *less subjective*.

### D. Idempotence & Recovery — harden `auto_implement`
All four. A task-specific recovery section (what state matters, what's safe to rerun,
what needs a human gate) in `FINAL.md`/`IMPLEMENTATION.md`. Complements §12's
low-level effect idempotency with plan-level recovery. Apply mainly to
implementation / action / pipeline / `auto_implement` ideas.

### E. Fusion "compare-don't-merge" lens in `consensus.md` + `review/consensus.md`
All four. Add an **advisory** comparison frame (a drafting discipline, *not* a judge
authority): *agreement / contradictions / partial-coverage / unique-insights /
**blind spots***. The standout is **blind spots** — *"what did no participant address?"*
— the one genuinely new Fusion idea and the single clear coverage gap: our
independent-then-converge flow catches errors in what was written, not shared
omissions. Forward-looking blind-spot check before finalize/complete (vs. §13 which
only finds them after failure). Append-only signoffs remain the real gate.

### F. "Confidently wrong" → §13 retro signal (NOT a 5th severity)
All four. Capture confident-but-wrong process failures (a dismissed CRITICAL/MAJOR
per agent, an unsupported assumption that shaped FINAL.md, a missed risk that caused
fix-up churn) as **retro evidence** for `parley retro`, to improve future harness
proposals — explicitly **not** a blame label, merge blocker, or new review severity.

## Agreed — REJECT (unanimous)

- **Confidence-by-breadth / majority weighting as a gate** — would let a 3-1 majority
  override a correct minority `BLOCK`; erodes all-participant signoff. (Recording
  "all agreed on X" as non-binding context is fine.)
- **A dedicated single-model judge with authority** — the consensus *drafter* already
  synthesizes; a judge role conflicts with equal participant signoff.
- **Hiding raw round files behind a judge summary** — the audit trail is a feature.
- **Fusion panel / recursion-depth / cost / default web-search-on-every-participant** —
  API plumbing; cost, privacy, prompt-injection risk; irrelevant to a file protocol.
- **Collapsing the deck into one ExecPlan file** — the multi-file shape *is* what gives
  independence, append-only signoffs, and collision avoidance.
- **ExecPlan "proceed without prompting" autonomy** — dangerous across parley-deck's
  quorum / human / production / protocol-change gates; the driver advances only inside
  approved boundaries.
- **ExecPlan narrative-prose-over-lists maximalism** — our structured headings/tables
  are a strength for multi-agent + tooling parsing; keep them. Borrow only the
  "each milestone reads goal→work→result→proof" discipline.

## Suggested adoption path (advisory, from codex; endorsed)

1. **Trial** the ExecPlan-style `IMPLEMENTATION.md` discipline on *one* medium/large
   implementation idea (ideally driver/pipeline-related). Test: can a fresh agent
   resume from `FINAL.md` + `IMPLEMENTATION.md` *without* session transcripts?
2. **Trial** the Fusion comparison + blind-spots block in *one* `consensus.md` and one
   `review/consensus.md`. Keep advisory; let signoffs test its accuracy.
3. **Feed** the results into a §13 `parley retro` pass. If retro extracts better
   failure modes from the new structured sections, they earn their keep; if they turn
   to boilerplate, keep them **opt-in for risky/long-running ideas only**.

## Cost realism
Reject Fusion's 4–5× tax. The ExecPlan additions are bounded to two artifacts and a
consensus subsection; keep full rigor **conditional** (complex / `auto_implement` /
driver-managed ideas), lean defaults for trivial ideas — consistent with parley-deck's
existing "don't add multi-agent overhead where it doesn't pay" principle.

## Next step (pending human decision)
This is a brainstorm checkpoint. If you want to proceed toward ratification, the
protocol path is: gather round-02 cross-review + append-only signoffs here, then open
a dedicated **meta-protocol-change** idea that actually edits both COOPERATION.md
copies (drift-guard lockstep) for the agreed items. **No protocol text changes until
you approve.**

---

## Signoffs (Phase 3 — append-only)

The human approved the ratification path. Each participant appends its own block
below after reading all four `round-01/*.md` files and this consensus. Accept with
`✅ ACCEPT` or block with `❌ BLOCKER` + specifics. Scope being ratified: the agreed
ADOPT/ADAPT items (A–F) and the REJECT list, as protocol-text guidance; exact edits
are specified in `FINAL.md`.

### Signoff: claude — 2026-06-18
Status: ✅ ACCEPT
Notes: I drafted this synthesis and stand by it. Two points I want on record for
FINAL.md: (1) the new living/static sections must be **conditional** — full rigor for
complex / `auto_implement` / driver-managed ideas, lean defaults otherwise — so we
don't tax trivial ideas; (2) all edits land byte-identical in both COOPERATION.md
copies' shared zones (drift guard), embedded default stays genericized. F (confidently-
wrong) is ratified as protocol-text guidance only; the `parley retro` mining for it is
a deferred follow-up, not part of this change.

### Signoff: codex — 2026-06-18
Status: ✅ ACCEPT
Notes: I accept the consensus draft as protocol-text guidance for items A-F and the reject list. The key constraint for FINAL.md is preserving static consensus snapshots while letting IMPLEMENTATION.md carry living execution state, recovery notes, and validation evidence.

### Facilitator note: agy live-signoff excluded — 2026-06-18 (claude)
agy (Gemini 3.5 Flash) repeatedly hung in headless `--print` while appending its
signoff (two attempts, no output, killed). The **user explicitly authorized skipping
agy's live signoff** for this idea. This is a documented per-participant tooling
exclusion, **not** a solo run: the quorum proceeds with claude + codex + hermes (three
participants), and agy's own `round-01/agy.md` already independently endorsed the
consensus items (it returned ADAPT for A/B/E, ADOPT for C, ADAPT for D, REJECT for the
confidence/judge mechanisms — i.e. the exact shape ratified here). agy's round-01
remains canonical evidence; only its Phase-3 signoff act is waived by operator ruling.

### Signoff: hermes — 2026-06-18
Status: ✅ ACCEPT
Notes: My round-01 maps cleanly onto items A-F and the reject list; the synthesis is faithful, not lossy. I specifically endorse item B's resolution — keeping FINAL.md immutable after Phase 4 is the stronger invariant than my initial "living FINAL.md" framing, since the consensus snapshot is the audit trail; ExecPlan living-state belongs in IMPLEMENTATION.md alone. The conditional-rigor / lean-defaults principle is the right guard against ceremony tax on trivial ideas.
