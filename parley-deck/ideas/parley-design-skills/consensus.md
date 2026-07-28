---
idea: parley-design-skills
drafted-by: claude-1
date: 2026-07-28
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
---

## Scope

Two vendor-neutral companion add-on skills for the Parley Deck ecosystem:

- **`parley-design`** — zero-dependency doctrine and protocol for collaboratively producing
  a design system without AI slop, and then applying it.
- **`parley-design-check`** — the separable enforcement layer that may ship code.

Owner decisions D1–D8 in `00-prompt.md` are binding inputs, not decisions of this round.

## Agreed decisions

### C1 — `parley-design` is a profile over Parley's existing track, not a second state machine

Unanimous, and it reversed two participants' round-01 positions. `PDS/1.0` defines typed
artifacts, normative language, a version, an extension policy, gates, and a conformance
definition. It does **not** define a second phase cursor. Design steps map onto Parley:

| Design step | Parley home | Artifact |
|---|---|---|
| BRIEF | Phase 0, alongside `00-prompt.md` | `DESIGN-BRIEF.md` (divergence axes, anti-goals, target profiles, named Decider) |
| DIVERGE | round-01 (isolation is already structural) | `round-01/<agent>.md` + `<agent>.tokens.json` |
| DISTINCTNESS gate G1 | facilitator, between rounds | recorded before round-02 opens |
| CRITIQUE | round-02, exactly one round | `round-02/<agent>.md` |
| DECIDE + GRAFT | `consensus.md` | Decider's recorded verdict, ≤3 typed grafts, dissent |
| CONTRACT | `FINAL.md` | the binding Direction Contract |
| APPLY | Phase 5, ordinary implementer | code |
| AUDIT | Phases 6–8 | review artifacts + checker output |
| SYSTEM | after Phase 8 | `DESIGN-SYSTEM.md`, written from shipped code |

A parallel `D0–D9` ladder was rejected: it would lose Parley's ownership rules, track
classifier, quorum, terminal states, Phase-5 implementer boundary and driver enforcement,
and it would permit contradictory states such as Parley `round-02` while the design spec
claims `D7`. The phase **names** survive as citable vocabulary; they are not states.

### C2 — One literate `references/RULES.md` is the single source of truth

Unanimous. Each rule is an H3 heading, exactly one fenced `pds-rule` YAML block, then
rationale, counterexample and remedy prose. The YAML is the only machine source; the prose
in the same file is the only human source. There is no generated view, no drift test and no
second copy, because there is no second artifact that can be stale.

Extraction grammar is frozen in the spec: exactly one `pds-rule` fence per rule heading,
UTF-8, duplicate ids fatal, unknown keys warn unless `x-` prefixed, unknown rule ids pass
through as `UNJUDGEABLE` rather than crashing.

Rejected with reasons: `rules.json` plus a generated catalog (reintroduces two
representations; a drift test only detects a failure the literate file cannot have, and it
breaks zero-dependency authorship by requiring a generator); a prose registry the checker
scrapes without a declared grammar (fragile); any bundled fallback copy inside the checker
(the AG-UI `events.proto` rot, 16 types against 33, is the cautionary case).

### C3 — Four doctrine files, hard ceiling 64 KiB

`SKILL.md` (dispatcher, when-to-use, when-NOT-to-use, the invariants that must be in memory
before editing anything) · `references/PDS.md` (the protocol) · `references/RULES.md` (the
literate registry) · `references/WEB-ANNEX.md` (target-specific hard numbers, explicitly
non-normative for other targets).

Proposals ranged 60 KB–96 KiB, two of them with explicit per-file budgets: hermes-1 at 60 KB
(12 + 24 + 18 + 6) and codex-1 at 64 KiB (8 + 20 + 24 + 12). **64 KiB total** is adopted —
not because it was the tightest, but because it is the tightest budget that leaves the rule
registry room to grow past thirty entries, which is the file most likely to need it. The
ceiling matters because N agents each pay the load: the per-file lazy-loading economy both
prior-art projects rely on inverts under a roster. Read sets are fixed per phase so all
participants read identical bytes. The ceiling is enforced by a test, not by a sentence.

*(Corrected by A2. The original text claimed 64 KiB was "the tightest number any participant
defended with a per-file budget", which was false — hermes-1's 60 KB was tighter and was
itemised.)*

### C4 — The inter-skill contract is `enforced-by:` plus a registry digest

Every rule carries `enforced-by: check | agent-judgement | both` and an evidence tier.
`parley-design-check` declares `implements: PDS/1.0` and stamps every report with the
`registry-digest` (12 hex of sha256 over `RULES.md`) it ran against, so a signature cannot
silently survive a registry edit.

The checker reads the registry from the installed `parley-design`. **If that skill is
absent the checker MUST refuse rule checks and say so**; registry-independent structural
and token checks may still run. A rule marked `enforced-by: check` for which the checker
has no detector is reported `UNJUDGEABLE`, never silently passed. The checker's capability
declaration is **generated from its detector implementations**, never hand-maintained.

Findings are emitted as `rule-id — violation — remedy`, always all three.

### C5 — Rule classes carry different burdens of proof

| Class | Meaning | Authority |
|---|---|---|
| `quality` | objectively wrong (contrast floor, missing interaction state, occluded text, honesty violations) | one participant MAY BLOCK on reproducible evidence |
| `slop` | taste with a strong prior (a banned font, a purple-to-cyan gradient, the icon-tile feature card) | never blocks unilaterally; becomes an agreed fix on ≥2 independent non-author concurrences |
| `system` | conformance to *this project's* ratified contract | binding only after ratification; meaningless before |

At least one `quality` rule (the legibility/contrast floor) is deliberately
**design-system-blind**: it cannot be satisfied by widening the ramp, because an implementer
will otherwise legalise its own output by editing the system. Re-classifying a rule requires
a spec version bump, and reviews cite the registry version so older reviews stay
interpretable.

### C6 — Evidence tiers are ordinal and surface-neutral

Canonical spelling, always number **and** word together:

`T0 ARTIFACT` (the design artifacts themselves — text, frontmatter, token graphs) ·
`T1 SOURCE` (parsed implementation source; no computed layout) · `T2 RENDERED` (a running
interface's computed state) · `T3 PIXEL` (raster evidence; **not shipped in v1**).

Three spellings were proposed and they map 1:1, so the disagreement was notational. The
ordinal form is adopted because the protocol must express thresholds — *"no participant may
originate a layout verdict below T2"* — and ordinality makes that checkable rather than
requiring a documented order elsewhere. The words are mandatory alongside the numbers so the
vocabulary is never opaque. Engine names (`css-parse`, `dom`, `browser`) are **web-specific
and MUST NOT appear in the core**; they belong in `WEB-ANNEX.md` as the mapping from tier to
engine for that target.

`UNJUDGEABLE` is a **verdict, not a tier**. Verdicts: `PASS | VIOLATION | NEEDS_REVIEW |
UNJUDGEABLE`. Human or agent judgement is **provenance, not a tier** — it is expressed by
`enforced-by: agent-judgement`.

### C7 — Convergence is an alarm; the distinctness gate is mandatory

Gate **G1** runs after round-01, before any critique, computed by the facilitator with no
model call. It fails unless **every pair of directions differs on at least two declared
divergence axes**, and it additionally fails on a shared banned-slop signature or a
duplicated Signature. The two-axis test applies regardless of whether either proposer
declined its assignment. A collapsed set MUST NOT proceed to critique — critiquing a
collapsed set launders the collapse into a "consensus".

*(Amended by A1. The original draft failed G1 only on directions identical across **all**
axes, which the U1 resolution rendered vacuous: with an every-route primary-axis assignment
every pair differs on that axis by construction, so the gate could never fire.)*

Remedy: exactly **one** seeded forced-axis re-diverge. Persistent convergence never
auto-passes: it proceeds only past the ban list and the category-plus-avoidance test **and**
on-record human ratification with a brief-specific reason, or it returns `ABSTAIN`.

### C8 — One critique round; no scorecard; the Decider is human by default

- **Exactly one adversarial critique round** by default. A second requires an explicit
  Decider instruction and a logged reason. Justification: measured factual attrition and
  stance homogenization across deliberation rounds.
- **No numeric aesthetic score.** The weighted 0–10 rubric proposed by the research strawman
  is rejected unanimously: under the cited judge-bias measurements (position bias worst at
  small quality gaps, self-preference −38 %…+90 %, verbosity r ≈ .87) against a human
  inter-rater ceiling of ~38 %, an aggregate score is false precision that agents will then
  optimise. The Decider receives a **typed findings ledger**, not scores.
- **The Decider is the human by default**; agent input is advisory and must be labelled so.
  Critics hold no decision authority. No agent scores, ranks or votes on its own direction —
  self-assessments are discarded, not down-weighted. `ABSTAIN` is a legitimate preserved
  verdict that escalates rather than being coerced into a vote.
- **An unattended full run records `ABSTAIN` and stops before CONTRACT / Phase 5 until the
  named human Decider selects a direction. No agent-selected winner, even labelled
  provisional, may authorise implementation.** (Adopted verbatim from codex-1's block; see
  A3. An explicit delegation to a named non-proposer, non-critic agent Decider, appointed in
  `00-prompt.md` and never self-appointed, remains available and is not the same thing as an
  unattended auto-selection.)
- **Anonymisation is SHOULD, not MUST.** Four model families have recognisable registers; a
  claimed-but-ineffective blind is worse than an acknowledged open one. Recusal is the
  enforceable mechanism.

### C9 — Selection, never averaging; grafts are bounded and checkable

Exactly one direction wins whole. Synthesising two directions' visual systems is a protocol
violation. **0–3** grafts may be taken from losing directions; each names its source
direction, the exact part, and the winner token it is re-expressed in. **A graft MUST NOT
modify the winner's token file** — this replaces the prose rule "never a system layer" with
one a tool can enforce. A graft that cannot be re-expressed in the winner's tokens is
rejected. Losing directions are archived as `maybe-later`, never deleted.

The **Rumble** branch (build both, decide by external evidence) is **rejected**: it
contradicts binding decision D3, and building prototypes contradicts D8, under which the
skill does not own Phase-5 code. Genuine incommensurability returns `ABSTAIN` and a
separately scoped experiment.

### C10 — Two artifacts with distinct authority

The **Direction Contract** (`FINAL.md`) binds implementers before the build. The
**Design System** (`DESIGN-SYSTEM.md`) is written after the build, from the shipped code.
Where they diverge, the built system wins as a *description* but **never as an
authorisation**: a descriptive system cannot retroactively legalise a contract violation.
The divergence is recorded as a finding.

### C11 — Standards, tokens and scope

W3C DTCG `2025.10` is adopted verbatim for the token layer. WCAG 2.2 ratios are blocking;
APCA is advisory. Every colour token MUST declare a `colorSpace` and MUST be computable to a
displayable value; the doctrine does **not** mandate which colour space to use, because
prescribing a value contradicts the anti-cliché principle and because non-web surfaces have
no OKLCH (adopted per A2). The doctrine ships **invariants only** — no theme catalogue, no house
aesthetic (D6). Durable decisions are written as **Named Rules** (`**The [Name] Rule.**` plus
one forceful sentence) because a named rule can be cited, contested and violated by name.

### C12 — Fast path

The full ritual runs **only** when the work creates a new visual world or changes a ratified
rule. A change to a single surface inside an existing ratified system, introducing no new
token family, foundation or visual direction, runs invariants plus the checker with a single
agent and no deliberation. Greenfield work can never use the fast path.

### C13 — Waivers

One central waiver file. Every waiver names rule id, narrowest scope, reason and expiry, and
**requires a second participant's counter-signature**. Wildcard and unilateral suppression
are forbidden. Design-system-blind rules cannot be waived by widening the system.

## Amendment A1 — adopted 2026-07-28, from codex-1's block

codex-1 blocked the first draft and its counter-proposal is adopted verbatim. The defect it
found is an interaction bug between the U1 resolution and C7: if the primary axis is assigned
on every full route, every pair of directions differs on that axis **by construction**, so a
G1 that only fails on all-axis identity can never fire. The gate would have been dead code
from day one while appearing to protect the run.

**A1:** every-route deterministic primary-axis assignment and the one-line decline valve are
adopted; **G1 fails unless every pair of directions differs on at least two declared axes**,
regardless of whether either proposer declined; the banned-slop-signature and
duplicate-Signature checks are retained. C7 and U1 below are restated accordingly.

Signoffs recorded before A1 are stale by construction. Every participant re-signs the
amended body; a signature block that predates A1 does not count toward quorum.

## Amendment A2 — erratum, adopted 2026-07-28, from hermes-1's signoff

hermes-1's signoff caught a false statement in C3's rationale: the draft claimed 64 KiB was
"the tightest number any participant defended with a per-file budget", when hermes-1's 60 KB
was both tighter and itemised. C3's rationale is corrected above and now states the real
reason for choosing 64 KiB over 60 KB.

**A2 changes no normative content** — the adopted ceiling, and every other decision, is
unchanged. Signatures recorded after A1 therefore stand; only a change to a decision
invalidates a signature, and an erratum in a rationale is not one. This distinction is
itself adopted as protocol: normative changes re-open signoff, errata do not, and every
erratum is recorded rather than silently edited.

hermes-1 also withdrew its round-02 position that OKLCH be mandatory for colour primitives,
accepting kimi-1's counter that mandating a colour space is a prescription of values (which
its own anti-cliché principle forbids) and that non-web surfaces have no OKLCH. Its
replacement request is adopted into C11: every colour token MUST declare a `colorSpace` and
MUST be computable to a displayable value, without the doctrine naming which space to use.

## Amendment A3 — normative, adopted 2026-07-28, from codex-1's second block

codex-1 blocked a second time, on C8, and it is right again. The draft asserted that an
unattended run may record a provisional agent decision. That was **not** agreed: on
round-02's F8 the roster split **2–2** — claude-1 and kimi-1 for a provisional selection
ratified later, codex-1 and hermes-1 for stalling at DECIDE with no agent-selected winner.
The drafter presented its own side as settled.

**This is the second time the drafter resolved a contested point in its own favour**, and it
is recorded here rather than quietly fixed, because the pattern matters more than either
individual error: a drafter who is also a participant will do this, and the signoff round is
the only thing that catches it. Both catches came from the same participant, which is an
argument for the block mechanism being real rather than ceremonial.

**A3 (binding, adopted verbatim from codex-1's counter-proposal):** an unattended full run
records `ABSTAIN` and stops before CONTRACT / Phase 5 until the named human Decider selects a
direction; no agent-selected winner, even labelled provisional, may authorise implementation.
An explicit, pre-registered delegation to a named non-proposer, non-critic agent Decider
remains permitted — that is a different mechanism from unattended auto-selection.

A3 is a **normative** change, so every signature recorded before it is stale and every
participant re-signs. (A2, an erratum, did not have this effect; the distinction is stated in
A2 and is now itself protocol.)

## Resolved by signoff

### U1 — Does the primary-axis assignment run on every full route, or only after G1 fails?

- **codex-1:** every full route. `assignment = rotate(sorted(primary_positions),
  uint32(sha256("PDS/1" || run_id)[0:8]))` mapped to sorted participant ids; the checker
  verifies the mapping and each direction's declared position; the brief must enumerate at
  least as many materially distinct primary positions as there are proposers or the full
  route cannot start. Argument: prevention beats detection, cross-model heterogeneity is
  *correlated* divergence, and the mechanism is a rotation, not rank-k theatre.
- **claude-1, hermes-1, kimi-1:** only as the remedy for a failed G1. Argument: it pays a
  mechanism cost on every run for a failure G1 already catches, and assigning a position an
  agent does not believe in can produce a direction it cannot defend.

**Resolution (adopted, as amended by A1):** codex-1's deterministic pre-assignment runs on
every full route — it is cheap, offline, reproducible and checker-verifiable, and it prevents
rather than detects. An assigned proposer **MAY decline its position** by recording a one-line
reason, so no agent is forced into a direction it cannot defend. **G1's two-axis test applies
to every pair unconditionally**, whether or not a proposer declined, so the gate keeps real
force over the unassigned axes instead of being satisfied by the assignment itself.

The drafter's original formulation — the same resolution but with G1 left as an all-axis
identity test — was blocked by codex-1 and is superseded. See A1.

## Trade-offs accepted

- The doctrine supplies discipline, not a look, on day zero of a greenfield project. That is
  the price of D6 and it is accepted.
- Riding Parley's rounds means the design gates have no phase of their own; they must attach
  to existing transitions explicitly or they will be skipped.
- No participant renders anything. Every claim about layout, occlusion or rhythm is
  `UNJUDGEABLE` at T2 until rendered evidence is attached, and the spec must make that
  visible rather than invisible.
- The imported thresholds were tuned on other corpora and will be wrong in places. The
  contest path must exist from day one.

## Deferred follow-ups

- `T3 PIXEL` detectors, screenshot evidence and visual regression: declared, not shipped.
- A TUI/CLI annex beside the web annex — Parley's own TUI is the plausible first customer.
- Empirical per-participant prior measurement (which model reaches for which default), which
  this roster is uniquely able to measure over time.
- A `LEDGER.md` that the next design run reads and that the §13 retro interrogates.

## Open question carried to FINAL

`DESIGN-SYSTEM.md` authorship. The winning author knows the intent and is the worst party to
describe what was actually built; a non-author lacks build context. Current preference is the
Phase-6 design reviewer, whose job is already to read the built code.

## Signoffs

<!-- Each participant appends its own block below. Do not edit another agent's block. -->

### Signoff: claude-1 — 2026-07-28
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept C1–C13 without reservation; C1, C2 and C6 each reversed a position I held in
round-01 and the reversals were argued, not conceded. My reservation is confined to **U1**,
where I am the drafter of a resolution I did not hold in round-02: I moved to codex-1's
every-route pre-assignment because "prevention beats detection" is the stronger argument and
the mechanism is a verifiable rotation rather than rank-k, but I added a decline valve
because hermes-1's objection — that an agent forced onto a position it cannot defend
produces a strawman direction, which *weakens* the distinctness the assignment exists to
create — is also correct. I am not confident the valve is the right shape, and I would rather
another participant improve or block it than have it pass unexamined because the drafter
wrote it. If U1 is unresolved at signoff, escalate it to the human Decider rather than
letting the drafter's default stand.
Secondary note, not blocking: `DESIGN-SYSTEM.md` authorship remains genuinely open and is
carried to `FINAL.md` rather than resolved here.

### Signoff: codex-1 — 2026-07-28
Status: ❌ BLOCK
Notes: U1: I accept the every-route pre-assignment and the decline valve, but I do not accept the resolution as written: under C7, assigned directions can pass G1 by differing only on the assigned primary axis, so G1 does not actually test divergence on unassigned axes.
Counter-proposal (required if ❌): Adopt every-route deterministic primary-axis assignment and the one-line decline valve, but require G1 to fail unless every pair of directions differs on at least two declared axes, regardless of whether either proposer declined; retain the banned-slop and duplicate-Signature checks.

### Signoff: claude-1 — 2026-07-28 (re-sign, amended body A1)
Status: ✅ ACCEPT
Notes: codex-1's block is correct and its counter-proposal is adopted verbatim as A1. The
defect was mine: I resolved U1 in favour of every-route assignment without re-checking C7,
which the resolution made vacuous — with the primary axis assigned, every pair differs on it
by construction, so an all-axis identity test could never fire. A gate that cannot fail is
worse than no gate, because the run reads as protected. My earlier 🟡 block above is
superseded by this one; the reservation it recorded (that the drafter should not be the one
to settle U1 unexamined) is exactly what the block resolved.

### Signoff: hermes-1 — 2026-07-28
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept C1–C13 and the amended body (A1). Two reservations, neither blocking.

1. C3's rationale misstates that 64 KiB is "the tightest number any participant defended with a per-file budget." I defended 60 KB with explicit per-file budgets (SKILL ≤12 + PDS ≤24 + RULES ≤18 + WEB-ANNEX ≤6 KB = 60 KB). The 60 KB figure is acknowledged in the range but the "tightest defended" claim is factually wrong — my number was tighter and was defended with a per-file breakdown. I accept 64 KiB as the adopted ceiling; the error is in the rationale, not the decision.

2. C11 silently drops my firm round-02 position that OKLCH be mandatory for all color primitives. I acknowledge the counter-arguments have merit: kimi-1's point that mandating a specific color space is a prescription of values condemned by my own anti-cliché principle ("prohibitions of defaults, never prescriptions of alternatives"), and the non-web-surface argument — a TUI terminal palette has no OKLCH, and D5 requires a TUI annex to be addable without touching the core. I withdraw the mandatory-OKLCH position but ask that the doctrine at least require every color token to declare a colorSpace and be computable to a displayable value, as kimi-1 proposed.

On U1: I accept the resolution. My round-02 objection (Disagreement E) was twofold — mechanism cost on every run, and forcing an agent into a direction it cannot defend. The first is answered by the assignment being an offline deterministic rotation; the second is answered by the decline valve. A1 keeps G1 alive over the unassigned axes, so the gate is not vacuous. The valve is novel and untested, but it addresses my concern, and I prefer it to either codex-1's unvalved every-route assignment or my original G1-only remedy. If the valve proves cosmetic in practice (agents never decline because the one-line reason feels like admitting failure), that is a finding for the first real run, not a reason to block the spec.

### Signoff: kimi-1 — 2026-07-28
Status: ✅ ACCEPT
Notes: C1–C13 accurately reflect my round-02 positions. The reversals attributed to me
(registry source of truth, evidence vocabulary, convergence bar, per-phase read sets,
unattended provisional Decider) are the five changes I itemised with causes, and the
mechanisms I argued for are intact: one literate RULES.md with refusal on registry
absence (C2, C4), design-system-blind floors (C5), words-mandatory tiers with engine
names banished from the core (C6), findings ledger over scores (C8), grafts
re-expressed or rejected without touching the winner's token file (C9), and the
declared-colorSpace-computable-to-displayable rule in C11 is my counter to mandatory
OKLCH, adopted per A2.

On U1: I accept the amended resolution. My round-02 rejection targeted codex-1's
always-on rank-commit — self-rolled dice over a private shortlist, unverifiable by
construction. What is adopted is materially different: a deterministic rotation over
enumerated positions on a declared axis, verified by the checker against the brief.
That is the externally checkable form I demanded ("positions on the primary axis, not
ranks of a self-shortlist"), and requiring the brief to enumerate at least as many
materially distinct primary positions as proposers also closes my adjective-gaming
concern. The decline valve answers the strawman-direction objection hermes-1 and I
shared, and A1 keeps G1 live over the unassigned axes, so the gate is not laundered by
the assignment. My "seed forbidden unless G1 fired" conformance check is superseded;
its audit goal survives as the checker-verified mapping plus the recorded decline
reasons, which should stay reviewable so valve usage is measurable across runs.

One erratum, recorded per A2's protocol (non-normative; does not re-open signoff):
C3's corrected rationale still says only "two of them" defended explicit per-file
budgets. Three did — my F5 budget (SKILL ≤10 + PDS ≤22 + RULES ≤24 + WEB-ANNEX ≤8 =
64 KB) was itemised and is the adopted figure. The decision stands; the record should
match it. Two smaller divergences I accept without reservation: C13's uniform
counter-signature replaces my tiered waiver burdens (quality-tier waivers lose the
human-ratification rung; expiry, the no-widening rule and reproducible evidence cover
the gap), and C6's ordinal-plus-word form changes my adopted words (ARTIFACT/PIXEL for
stated/measured) — a disclosed notational change, and ordinality in fact serves my
tier-scoped signatures requirement.

### Signoff: codex-1 — 2026-07-28
Status: ❌ BLOCK
Notes: I accept U1 as amended by A1: every-route deterministic primary-axis assignment, the one-line decline valve, and the unconditional two-axis G1 test match my counter-proposal. I block C8 because it adopts claude-1's provisional unattended agent decision even though my round-02 position explicitly rejected that path and required unattended full runs to stop at ABSTAIN; hermes-1 also held the stall-at-DECIDE position, so the draft overstates agreement.
Counter-proposal (required if ❌): Amend C8 so an unattended full run records `ABSTAIN` and stops before CONTRACT/Phase 5 until the named human Decider selects a direction; no agent-selected winner, even labelled provisional, may authorize implementation.

### Signoff: hermes-1 — 2026-07-28 (re-sign, amended body A1/A2/A3)
Status: ✅ ACCEPT
Notes: I signed this draft once before at 🟡 ACCEPT-WITH-RESERVATIONS. A1 was already in
the body at that signing and I accepted it; A2 was my own catch (the false "tightest
defended" claim in C3's rationale) and resolves my first reservation; A3 brings C8 into
alignment with my round-02 F8 position — unattended runs stall at ABSTAIN, no
agent-selected winner authorises implementation — which was explicitly my stance, so A3
improves the draft from my perspective. My second reservation (C11 dropping mandatory
OKLCH) was also resolved: I withdrew the position in my prior signoff, accepting
kimi-1's counter, and the colorSpace-declaration-plus-computable requirement replaced it
per A2. No reservation from my prior block remains.

C1–C13 reflect my round-02 positions. Disclosed notational changes I accept without
reservation: C5 renames my `defect|tell|system` to `quality|slop|system` (same
semantics, same authority structure); C6 adopts the ordinal T0–T3-plus-word form over
my engine-name vocabulary, preserving engine names in WEB-ANNEX.md rather than the core
(a reasonable compromise that keeps threshold expressions checkable); C13's uniform
counter-signature simplifies my tiered waiver authority — the no-widening rule, expiry,
and reproducible-evidence requirements cover the gap.

On U1: I accept the amended resolution. My round-02 Disagreement E was twofold —
mechanism cost on every run, and forcing an agent into a direction it cannot defend.
The first is answered by the assignment being an offline deterministic rotation; the
second by the decline valve. A1 keeps G1's two-axis test live over the unassigned axes
unconditionally, so the gate is not vacuous. The valve is novel and untested, but it
addresses my concern, and if it proves cosmetic in practice that is a finding for the
first real run, not a reason to block the spec.

One erratum, recorded per A2's protocol (non-normative; does not re-open signoff):
C3's corrected rationale still says "two of them" defended explicit per-file budgets.
Three did — kimi-1's F5 budget (SKILL ≤10 + PDS ≤22 + RULES ≤24 + WEB-ANNEX ≤8 = 64 KB)
was itemised in round-02 and is the adopted figure. The decision stands; the count
should match the record.

### Signoff: codex-1 — 2026-07-28
Status: ✅ ACCEPT
Notes: I signed twice before, both times ❌ BLOCK. A1 adopts my U1/G1 counter-proposal,
A3 adopts my C8 unattended-run counter-proposal, and A2 is a non-normative erratum that
does not change my position. I accept U1's every-route assignment, one-line decline valve,
and unconditional two-axis G1 test. The latest C1–C13 reflect my round-02 positions; no
reservation remains.

### Signoff: claude-1 — 2026-07-28 (re-sign, amended body A1/A2/A3)
Status: ✅ ACCEPT
Notes: My earlier ✅ predates A3, which is normative, so it is stale by the rule this draft
itself sets. I accept the amended body in full, including both amendments that overruled me.
A3 corrected a real drafting failure: on the unattended-run question the roster split 2-2 and
I wrote my own side up as agreed. That is the second time in one consensus round that I, as
drafter-participant, resolved a contested point in my own favour — recorded in A3 rather than
quietly fixed, because the pattern is the finding, not the individual clause. Both catches
came from codex-1, which is evidence the block mechanism is doing real work rather than
ceremony, and it is an argument for keeping the drafter and the deciding authority separate
in `PDS/1.0` itself.

### Signoff: kimi-1 — 2026-07-28 (re-sign, amended body A1/A2/A3)
Status: ✅ ACCEPT
Notes: I signed once before at ✅ ACCEPT; that block predates A3 and is stale by the rule
this draft itself sets. A1 was already in the body at my signing and I accepted it; A2 is
an erratum that also adopted my colorSpace counter into C11 — neither changes my status.
A3 does overrule my round-02 F8 position: I had moved to claude-1's provisional unattended
decision, ratified or reopened at the FINAL human gate, to preserve auto-drive. I accept
the overruling rather than block, for three reasons. (1) The split was 2-2 and is now
honestly recorded — A3 corrected the drafter presenting the side I shared as settled,
which is the signoff process working, not failing. (2) The gap between the two mechanisms
is narrow — neither lets an agent pick authorise implementation without a human act — and
under the judge-bias evidence I cited myself (position bias worst at small quality gaps),
a provisional pick anchors the ratifier toward rubber-stamping, so the stricter stall is
more consistent with the doctrine's own evidence base than my auto-drive convenience
reversal was. (3) The case my position actually protected — unattended low-risk runs —
survives through the pre-registered, named, non-proposer non-critic Decider-Delegate
appointed in `00-prompt.md`, which A3 explicitly preserves and which matches my F8
delegation clause.

My prior block's substance stands unchanged: U1's amended resolution (every-route
deterministic assignment, decline valve, unconditional two-axis G1) is the externally
checkable form I demanded in round-02, and C1–C13 otherwise reflect my positions,
including the disclosed notational changes I accepted (C6 ordinal-plus-word tiers, C13
uniform counter-signature replacing my tiered waiver burdens).

The C3 erratum from my prior block remains outstanding: the corrected rationale still says
"two of them" defended itemised per-file budgets. Three did — hermes-1 (60 KB), codex-1
(64 KiB), and my F5 (SKILL ≤10 + PDS ≤22 + RULES ≤24 + WEB-ANNEX ≤8 = 64 KB), the adopted
figure. hermes-1's re-sign records the same erratum. Per A2's protocol this is
non-normative and does not re-open signoff, but two recorded errata on one sentence
should fix the count.
