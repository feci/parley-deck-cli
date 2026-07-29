---
agent: claude-1
idea: parley-design-skills
round: 1
date: 2026-07-28
---

## Summary

Ship **`PDP/1` — the Parley Design Protocol** as the substance of `parley-design`, and a
thin, tiered checker as `parley-design-check`. The two skills share exactly one object: a
**numbered rule registry that is pure data**. The doctrine skill owns the registry and can
be fully useful with nothing installed; the check skill implements detectors keyed to the
same rule ids and declares, per rule, which engine tier it can actually reach.

My central claim: **the only defensible reason for multi-agent design is diversity
generation, and everything else in the protocol exists to stop that diversity from
collapsing before it is used.** The published evidence in the kickoff is decisive — debate
does not improve taste, deliberation homogenises stance across rounds, and human agreement
on aesthetics is ~38 %. So the protocol should spend its complexity budget on *isolation,
a distinctness gate, and a clean selection rule*, and spend almost nothing on discussion.
A design protocol whose centre of gravity is "the agents talk until they agree" is
building the exact failure the skill exists to prevent.

## Proposed approach

### 1. The protocol shape (D4)

`parley-design/reference/PROTOCOL.md` is written as a specification, not advice, borrowing
the AG-UI structure that the owner pointed at:

- **Versioned identifier `PDP/1`** in the title, and a stated compatibility policy.
- **RFC 2119 keywords in uppercase** (`MUST`, `MUST NOT`, `SHOULD`, `MAY`) with a
  normative-language section. AG-UI deliberately uses lowercase prose instead; the research
  digest calls that its weakness and I agree — with four heterogeneous model families
  reading the same text, ambiguity is not a style choice, it is a defect.
- **Every phase gets an identically-shaped entry**, never varying: `### P<N> · NAME`, then
  one sentence of what it is, then one paragraph of *why it exists and what the next phase
  does with its output*, then a table of the artifact's required fields. AG-UI applies this
  four-part shape uniformly across 33 events and it is the single most copyable thing in
  that spec.
- **A `Contract rules` section**: numbered rules, each with a **bolded short name** then
  one sentence (`1. **Isolation.** A Proposer MUST NOT read another Proposer's direction
  before the exhibit opens.`). Short names are what cross-review findings cite.
- **An extension policy stated in three sentences**, copied in shape from AG-UI's `reason`
  taxonomy: a small set of spec-defined values; any other string is a valid extension;
  extensions SHOULD be namespaced `<project>:<name>`; the `pdp:` prefix is reserved; and
  **an unknown value MUST NOT cause an error** — it is rendered and passed through.
- **A conformance section** that defines what it means to conform: which artifacts must
  exist, which fields are required, which gates must have run. Conformance is a property of
  a *run*, not of a tool, so any roster can claim it.

### 2. Phases and artifacts

`PDP/1` has six phases. I deliberately cut the research strawman's eight down, because
every phase is a round-trip across four agents and the evidence says extra rounds cost more
than they return.

| Phase | Artifact | Owner | Purpose |
|---|---|---|---|
| **P0 · BRIEF** | `design/BRIEF.md` | facilitator, human-ratified | problem, audience, goals findings are judged against, hard constraints, and the **`divergence_axes`** — the named axes on which directions MUST differ. No axes, no P1. |
| **P1 · DIVERGE** | `design/directions/<agent-id>.md` | each participant, in isolation | one committed direction per agent, fixed schema, hard length cap. |
| **P2 · EXHIBIT** | `design/EXHIBIT.md` | facilitator, deterministic | anonymise to stable slugs, randomise order per reader, and run the **distinctness gate**. |
| **P3 · CRITIQUE** | `design/critique/<agent-id>.md` | each participant, on others only | one round, assigned lens, typed entries, author silent on its own. |
| **P4 · DECIDE** | `design/DECISION.md` | the Decider | one winner whole, dissent recorded, losers archived as `maybe-later`, ≤3 named grafts. |
| **P5 · RATIFY** | `DESIGN-CONTRACT.md` (+ `design.tokens.json`) | winner's author | the pre-build binding contract; rides into `FINAL.md`. |

Then, after Phase 5 of the *Parley* protocol has produced code:

| **P6 · DOCUMENT** | `DESIGN-SYSTEM.md` | the design reviewer | written **from the shipped code**, per D7. |

**The direction artifact schema** (P1) is the load-bearing decision, because five prose
essays cannot be compared but five instances of one schema can be diffed field by field.
Required fields, each capped:

```
handle:            one word, not the agent's name  (anonymity survives the exhibit)
thesis:            ≤2 sentences — the idea it owns AND the category default it refuses
axes:              one declared position per divergence_axis in BRIEF.md
world:             palette strategy · type voice · spatial posture · motion posture
first-impression:  what the user sees first and what it asks them to do
tokens:            a real token table — scale, roles, radii, durations
refusals:          ≥2 named things this direction will not do
```

Two rules from the prior art earn their place here verbatim in spirit: *"if a block reads
like a mood, the direction is not decided yet"*, and the direction must be
**self-explanatory without its author present**. Both are reviewable by any agent without
taste judgment, which is what makes them good rules.

### 3. The distinctness gate is the heart of it

**G1 — DISTINCTNESS.** After P1, before any critique: if any two directions hold the same
position on **every** declared axis, or share a slop signature (same font strategy + same
colour strategy + same structural posture), **P1 has FAILED**. The protocol MUST NOT
proceed to critique. It re-runs P1 with **assigned** distinct positions on the primary
axis — the facilitator hands agent *i* the *i*-th position, derived deterministically from
the idea slug and round number so the assignment is reproducible and auditable, with no
RNG service anywhere (F4).

This is my answer to F2 as well: **convergence is an alarm, not a pass.** Four models
independently reaching for the same look is the training distribution talking. The gate is
the mechanism; the escape hatch is narrow and explicit — a converged set may proceed only
if it survives the ban list *and* the "could someone guess this aesthetic from the category
alone, or from category-plus-avoidance?" test, *and* a human ratifies the convergence on
the record.

This also happens to be the one thing in this whole design that no single-agent tool can
do, which is the honest justification for the add-on existing at all.

### 4. The rule registry — the shared object between the two skills

`parley-design/reference/rules.md`, one row per rule, **data only, no detection logic**:

```
id          pd-041                      stable forever; append-only; ids are never reused or re-meant
class       defect | tell | system      determines authority (see F1)
surface     core | web                  D5: core rules hold anywhere; web rules live in the annex
tier        stated | source | rendered | measured   the cheapest evidence that can decide it
name        Focus ring is animated
why         one sentence: why this reads as unfinished or machine-made
fix         one sentence: the corrective move
```

Three classes, three burdens of proof (F1):

- **`defect`** — objectively wrong (contrast below threshold, missing interaction state,
  text occluded, animated focus ring). A single participant MAY BLOCK on a `defect`.
- **`tell`** — taste with a strong prior (a banned font, a purple-to-cyan gradient, the
  icon-tile feature card). A `tell` is **advisory**; it becomes binding only by quorum.
- **`system`** — conformance to *this project's* ratified contract (off-scale spacing, a
  colour outside the ramp, a font outside the allowlist). Binding **once a contract has
  been ratified**, meaningless before. At least one `system` rule MUST be deliberately
  system-blind — a legibility floor that cannot be legalised by widening the ramp — because
  an implementer will otherwise widen the system to legalise its own output (F6).

This split is what stops the protocol ratcheting toward the safe category standard. If
every rule could block, four reviewers each holding a veto converge on the option nobody
objects to, which is precisely the slop we are trying to defeat.

### 5. `parley-design-check`

A small Node CLI, following the `parley-tracker` precedent (`bin/validate.js` +
`bin/claim.js`), not an impeccable clone:

- `bin/check.js` — runs detectors, emits JSON findings `{rule, class, tier, file, line, evidence}`.
- `bin/rules.js` — loads the registry from `parley-design`'s `rules.md` **if installed**,
  else from its own vendored copy, and reports which source it used.
- `conformance.json` — the machine-readable contract: for each rule id, the highest tier
  this version can actually reach, or `manual`. This is how the zero-dependency doctrine
  skill names a check it cannot itself run: it cites `pd-041` and the checker declares
  whether `pd-041` is machine-verifiable here and now.

**Tier discipline is the honesty mechanism (F3).** `stated` needs nothing. `source` needs
a file read. `rendered` needs a browser. `measured` needs pixels. A participant that cannot
reach a rule's tier writes `unjudgeable: rendered` and is **compliant, not silent**, and any
artifact produced below full capability MUST lead with a degradation banner. v1 ships
`stated` and `source` only; `rendered`/`measured` are declared-and-unimplemented rather
than pretended.

### 6. Answers to the remaining forks

- **F5 — size.** SKILL.md ≤ 12 KB as a dispatcher. Exactly four references: `PROTOCOL.md`,
  `rules.md`, `FLOOR.md` (≤50 lines, the only always-loaded file), `WEB-ANNEX.md`. Hard
  ceiling ~60 KB total. hallmark's ~400 KB economy depends on one agent lazily loading one
  file; under four agents that inverts into 4× cost plus divergent reads.
- **F7 — fast path.** The full ritual runs **only** when the work creates a new visual
  system or changes a ratified rule. Everything else — a component, a page inside an
  existing system, a copy fix — runs `FLOOR.md` + the checker, single agent, no
  deliberation. Deciding whether a hero is centred does not deserve four models and a
  consensus round.
- **F8 — the Decider.** The human is the Decider by default; agent scores and votes are
  **advisory and must be labelled as such**. In an unattended run the facilitator agent may
  decide, but MUST record `decider: agent (unattended)` and the decision is **provisional**
  until a human ratifies. Critics MUST have no decision authority, and no agent may score
  its own direction — its self-scores are *discarded*, not down-weighted.
- **Scoring.** Absolute against an anchored rubric, never pairwise, order randomised per
  reader, artifacts truncated to the cap before scoring. `ABSTAIN` is a legitimate verdict
  that escalates to the Decider rather than being coerced into a vote.
- **Grafts.** ≤3, each a discrete nameable detail, **never a system layer** — no colour
  ramp, no type scale, no grid. Each graft must state which winner token it is re-expressed
  in; a graft that cannot be re-expressed in the winner's tokens is rejected. This is the
  bound that keeps "graft" from quietly becoming "average".

### 7. Tokens

`design.tokens.json` in **W3C DTCG** shape (`$type` / `$value` / `$description`, groups,
aliasing) because it is the one interchange format with real tooling behind it, and because
it makes the token table mechanically checkable. Deviations, stated explicitly rather than
silently: we do not adopt the component-token sub-schema, and anything the format cannot
hold (motion posture, focus treatment, density) lives in `DESIGN-CONTRACT.md` prose as
**Named Rules** — `**The [Name] Rule.** <one forceful sentence>` — because a named rule can
be cited, contested and violated *by name*, and a bullet list cannot.

## Concerns / open questions

1. **Two artifacts (D7) is right but under-specified.** `DESIGN-CONTRACT.md` binds before
   the build; `DESIGN-SYSTEM.md` describes after it. The unanswered question is what happens
   when they disagree: my position is **the build wins and the divergence is recorded as a
   finding**, because a contract defended against reality is exactly the failure impeccable
   documented — but that means a ratified artifact can be overruled by an implementer, which
   is a real hole in the Parley authority model and needs an explicit rule.
2. **The core/annex split (D5) may not survive contact.** I am not yet convinced a
   genuinely surface-neutral core has enough hard content to change any output. If the core
   collapses into platitudes, the honest move is to admit the annex is the product and the
   core is a preface.
3. **Anonymisation is weaker than it looks.** Four model families have recognisable prose
   styles; stripping the filename will not fool a reader that has seen the other three. I
   would rather state this limit than claim blind review we do not have.
4. **Who writes `DESIGN-SYSTEM.md`?** The winning author knows the intent but is the worst
   party to document what was actually built. I lean toward a non-author, but that agent may
   lack the build context.
5. **The registry will be wrong at first.** Every threshold imported from prior art was
   tuned on someone else's corpus. We need a stated path for contesting a rule on the record
   — the §13 retro machinery is the obvious home, advisory-only, quorum over self-preference.

## Risks

- **The protocol becomes ceremony.** Six phases × four agents is a lot of machinery for a
  landing page. If the fast path is not genuinely the default, nobody will use this twice.
- **A new house style.** Rules that merely invert today's defaults produce a new,
  equally recognisable look. The registry must ship *invariants* and let the deck ratify the
  taste (D6), and the "guess it from category-plus-avoidance" test has to be a real gate,
  not a slogan.
- **Four blind reviewers, five times as confident.** Nothing in the roster renders anything.
  Agreement among four agents that never saw the page is the most dangerous artifact this
  add-on could produce; the tier/`unjudgeable`/banner machinery is the only thing standing
  between us and a signed lie.
- **Registry drift.** Two copies of the rule list (doctrine + vendored fallback) is exactly
  the drift surface this repo has been bitten by before. It needs a generated copy and a
  drift test, on the `TestEmbeddedDefaultMatchesLiveDeck` precedent, or it will rot.
- **Scope.** If I had to cut, I would cut in this order: the web annex depth, then
  `parley-design-check`'s `source` tier, then P3 critique entirely — keeping diverge,
  distinctness gate, and decide, which is where the actual value is.
