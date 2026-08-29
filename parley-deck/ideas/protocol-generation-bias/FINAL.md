---
idea: protocol-generation-bias
drafted-by: claude-1
date: 2026-08-29
status: final
participants: [claude-1, codex-1, hermes-1, kimi-1, zcode-1]
excluded: [opencode-1]
signoffs: 2 ✅ ACCEPT (codex-1, hermes-1), 3 🟡 ACCEPT-WITH-RESERVATIONS (kimi-1, zcode-1, claude-1)
---

## Final plan / specification

Three legs and a deletion. **No new flag, no new agent, no new round, no new file, no new finding
class.**

### Leg 1 — Acquire (targets B2)

An unconditional `## Existing alternatives` section in the round-1 prompt template, **all tracks
including `fast`**:

- Enumerate the mechanisms this proposal builds **by hand** — components **named**, not described.
- For each, name what the toolchain, stdlib, dependencies or platform **already ships**, with a
  locator.
- Mark each load-bearing element **constraint-forced** or **merely inherited**.
- A scoped null result is legal and must name the sources consulted. *"The hand-built route is
  correct"* is a valid outcome.

Carrier: `BuildRoundOnePrompt` (`internal/runner/runner.go:821`), gated by the exported
`HasNonEmptySection` family (`internal/protocol/reviewartifact.go:114`). **Full wording lives in the
template; `COOPERATION.md` carries the compact duty only.**

The enumerated form is load-bearing, not stylistic. Open "consider alternatives" is the
measured-null family; enumeration is the form with positive evidence. See References.

### Leg 2 — Exchange (targets B2's hidden-profile core)

After round 1 is sealed and **before the already-scheduled next decision prompt** (round 2 on
`standard`/`deliberation`; before collapsed FINAL on `fast`), the runner collects one sealed packet
per participant: **1–2 decision-relevant facts, and one reason that participant's currently
favoured option might be wrong.** No positions, no revised recommendation, **no asymmetry
assertion**. All packets released simultaneously; each participant echoes its own under
`## Evidence exchange`.

One bounded model call per participant, reported as runtime cost. **Not a new round** — it
restructures an existing step.

**Protocol text must carry the label verbatim: "structurally derived from HiddenBench; transfer
unverified; instrumented."**

### Leg 3 — Disposition (targets B1)

`consensus.md` carries `## Alternatives disposition`. Every alternative raised gets a stable
identifier, an **adopt / reject with the decisive reason**, and **FINAL may not contradict a
recorded adoption**.

Identifier format: the literal `ALT-`, then the raising agent's roster id, then `-R` and the round
number, then `-` and a per-round index. Example: `ALT-kimi-1-R2-1`. The identifier is the mechanical
handle — a disposition requirement without one is a prose rule, and prose is not a carrier.

**Terminal destination — settled here, per kimi-1's and zcode-1's reservations, not left to
default.** A contradicted adoption **blocks signoff and escalates to the owner**, who decides
dismissal-with-reason or a `-v2` successor idea. **The scanner never auto-halts.** This is the UFMCS
escalation form (a dismissed critical observation requires resolution with the commander) and it is
adopted as **doctrine, not evidence** — the handbook contains no empirical evaluation and says so.

**No new finding-class vocabulary.** `CRITICAL|MAJOR|MINOR|NIT` stands.

### Leg 4 — The deletion that funds it

§15.6's trigger conditions are deleted, including the *"primarily a judgment rather than a
mechanically decidable artifact"* carve-out — the clause that switched the protocol's only
forced-alternative rule **off** for exactly the class of task that started this idea.

**Byte accounting, corrected.** The consensus draft carried **−726 B**, which was wrong: that is
zcode-1's **round-1, section-only** figure inherited into a three-leg package. Both kimi-1 and
zcode-1 flagged it independently; zcode-1 named it *"the 'locks must not carry numbers' defect
repeating inside the consensus text itself"*, and it was the drafter's error.

Re-derived at ratification time and re-checkable in one command:

```
sed -n '1346,1368p' parley-deck/COOPERATION.md | wc -c    # → 1372
```

| | bytes |
|---|---|
| §15.6 deleted (`COOPERATION.md:1346-1368`) | **1,372** |
| replacement text, all three legs | **1,164** |
| **net** | **≈ −237** |

Net-negative holds under every measurement in the record (−726 / −342 / −237). **The ratified
`COOPERATION.md` change MUST embed its own replacement text and publish the actual diff against
1,372 B**, so the claim is verifiable by the command above rather than by trust.

## Purpose / user-visible outcome

A run stops being able to converge on the first plausible frame without anyone having looked at
what already exists.

Concretely, the case that started this: a template deliberation converged on hand-written build
scripts copying dependencies out of `node_modules`, and `pnpm deploy` — first-party, documented, one
command — only appeared after a human said it aloud. Leg 1 turns that from a **knowledge** question
("is there a simpler way?", which the evidence says bounces) into a **lookup** question ("what does
the toolchain itself document?"). Nobody has to *know* `deploy` exists; somebody has to *look*, at a
named place, and record what they found.

And the second case: an agent wrote *"Nobody proposed the option that actually exists"*, withdrew
its own proposal — and `FINAL.md` shipped the original design, never mentioning the alternative at
all. Leg 3 gives that sentence a destination.

## Context & orientation

**The finding of this idea is not the package. It is the carrier thesis.**

Rules carried by Go prompt templates plus a validator reach near-universal compliance. Identical
rules left in `COOPERATION.md` prose run in single digits — same deck, same roster, same corpus, one
variable. kimi-1 refined it decisively: carriage is **necessary and not sufficient**, because the
two halves of one carried cross-review rule diverge 2.5× (`responding-to:` 18.1% vs
`### @<other>` 7.2% of 349 files). **The half a scanner cannot check decays to prose rates even when
the prompt carries it.**

This reframes the kickoff brief's own finding 6. The 0–5% adoption figures are **not** evidence that
this deck picks bad gates; they are evidence that **prose is not a carrier**, and every opt-in flag
is prose. The drafter wrote that finding and the correct reading came from a participant.

**What the evidence killed.** Eight of sixteen negative findings indict a mechanism proposed in
round 1 of this very idea: generic divergence instructions, debiasing training, ACH, premortem,
more agents/rounds, conflictual framing, LLM-judge override, rigid formats. Every one of the five
participants changed position; several abandoned their own assigned axis.

**Late filing by an excluded participant.** `opencode-1` was excluded under §9.0 after ~10 failed
invocations (two models, three prompt strategies, all exiting 0 with no artifact) — then filed a
complete round-01 artifact at 14:05 on 2026-08-29, after the exclusion and during consensus. It
holds axis A2, the axis three participants independently named as the missing piece.

**Its content answers the question this design defers to it**: a `REFRAME` class, a route absorbing
before `FINAL.md` freezes, and a `## Frames considered` destination inside `FINAL.md`. The deferred
vocabulary question therefore has a **live input, not an absent owner**, and implementation must
read `round-01/opencode-1.md` before deciding it. Its summary states the split this design adopts:
*"B1 is an absorption failure. B2 is a generation failure A2 cannot fix alone."*

## Observable acceptance criteria

1. Byte-counting the ratified §15.6 replacement with `sed -n` over its line range piped to
   `wc -c`, plus the published diff, shows a **net reduction against 1,372 B**.
2. A round-1 artifact lacking a non-empty `## Existing alternatives` is **rejected by the validator**,
   not merely warned about — verified by a test that reverts the gate and watches it fail.
3. The exchange adds **zero** rounds, files, agents and artifacts. Verified by diffing the run's
   phase sequence before and after against the same idea.
4. The exchange execution prompt contains **no assertion that information asymmetry exists** —
   grep-checkable against the template.
5. `COOPERATION.md` contains the string *"transfer unverified"* attached to the exchange. The label
   is in the protocol, not only in this idea's artifacts.
6. A `consensus.md` recording an adoption, followed by a `FINAL.md` contradicting it, **blocks
   signoff and escalates to the owner** — and the scanner does **not** auto-halt. Both halves tested.

## Idempotence & recovery

- All four changes are additive-to-template and subtractive-from-protocol; re-running the
  implementation is a no-op on an already-migrated deck.
- **Pre-ratification baseline freeze — owner: the implementer, before the first ratified run.**
  Tag the ratification commit so the pre-adoption corpus is a named "before" state. This is the one
  step that cannot be retrofitted, and without it acceptance criterion 7 below is unfalsifiable.
  Both kimi-1 and zcode-1 raised its absence; it was missing from the consensus draft.
- Rollback is deletion of the template block plus restoration of §15.6 from the tagged commit. No
  data migration, no persisted state.

## Known risks / de-risking

**R1 — The exchange's headline effect may not transfer, and we know it.** HiddenBench moved GPT-4.1
from 3.7% → 80.0% (+76.3pp), but depended on agents **never being told** information asymmetry
existed. codex-1: that condition *"does not survive perfectly once the rationale is public."* The
faithful remainder is to omit the asymmetry claim from the execution prompt. **Ships instrumented
with pre/post recall measurement.**

**R2 — Group size is unread, and it gates the binding.** *Before implementation*, someone must read
the two-stage condition's group size from `arXiv:2505.11556`. Our dose-response figures
(3 agents +34.8% / 4 +25.0% / 7 +0.6%) are `SECONDARY` via the reference brief. **If the effect was
measured at ≤3 agents, the exchange initially binds small decks only.** Raised by kimi-1, endorsed
by zcode-1; **owner: the implementer, as a pre-implementation gate, not a retrospective item.**

**R3 — A scanner-checkable section can still become a ritual.** `## Refutation attempts` is gated
for non-emptiness and never checked for being work. `## Existing alternatives` inherits that risk
exactly. **De-risking, stated in advance so deletion is the default rather than a defeat: at the
first retrospective with enough data, compare the two for ritualisation against the frozen
baseline. If `## Existing alternatives` has become a filled-in ritual, it is deleted.**

**R4 — This deliberation is indicted by its own evidence.** HiddenBench measured +0.6% at 7 agents
against +34.8% at 3; this idea ran six participants across three rounds. The facilitator chose six
on the aesthetics of the topic, unevidenced, and no gate in this protocol checks that choice. **This
is a measurable claim, not a concession: the deck has 88 ideas with recorded participant counts and
outcomes.** Owner required.

**R5 — Correlated agreement.** Two readings are recorded rather than one. Against: the facilitator
circulated a mandatory corrections list and a mandatory sixteen-item negative list — per the
evidence a conform instruction and diversity-suppressing structure — and five files converged. For:
participants moved in **different** directions, several against their own prior work, and two
disagreements survived the pressure (the exchange-carrier split, settled on checkability grounds;
the byte divergence, which caught the drafter). **All three reserving signers independently reached
the same reading: one shared prior plus genuine movement, not five independent confirmations.**

**R6 — Provenance.** Of 24 external citations in round 1, verified by agents that did not write the
citing file, **20 survived and 4 did not** — including two by the drafter. Every failure was caught
by a reader who fetched the source; **none would have been caught by anything currently in the
protocol.**

## References

**Primary artifacts.** `00-prompt.md`; `round-01/` (6 files incl. the late `opencode-1.md`);
`round-02/`, `round-03/` (5 each); `consensus.md`;
`round-03-consensus-aborted-01.md` (the blocked draft, preserved);
`reference/research-brief.md` (918 lines, claims C1–C14 and negative findings A1–A16);
`reference/measure.sh` (canonical deck measurements).

**Evidence the design rests on** — full provenance and verification verdicts in the reference brief:

- Smith, Ward & Schumacher 1993, `10.3758/bf03202751` — instructing subjects to produce ideas *very
  different from the examples* did **not** reduce conformity; the conform instruction **did**
  increase it. The asymmetry rules out manipulation failure. *This is why leg 1 is enumerated.*
- Jansson & Smith 1991, `10.1016/0142-694X(91)90003-F` (**note: the round-1 citation gave a DOI that
  404s**) — fixation persisted onto features the instructions **explicitly forbade** (straws 17% vs
  1%, mouthpieces 39% vs 10%) and replicated in **professional engineers**. Reports means and
  percentages only — no p-values, no CIs, cells of 6–18. Do not call it statistically significant.
- Chrysikou & Weisberg 2005; George & Wiley 2020, `10.3758/s13421-019-01005-4` — **naming** the
  problematic elements diminished fixation where **describing** them did not; a list *plus* an
  avoid-instruction enhanced originality, the same examples without it produced nothing.
- Stasser & Titus 1985, `10.1037/0022-3514.48.6.1467` — discussion *"tended to perpetuate, not to
  correct, members' distorted pictures."* **B2 is a hidden profile, not a creativity failure.**
- HiddenBench, `arXiv:2505.11556` — exchange-then-decide +76.3pp; and 7 agents +0.6% vs 3 agents
  +34.8%. Both the design's foundation and its self-indictment.
- Tversky & Kahneman 1974, `10.1126/science.185.4157.1124` — *"The initial value, or starting point,
  may be suggested by the formulation of the problem."*
- Nemeth 2001, `10.1002/ejsp.58` — the authentic minority beat **all three** devil's-advocate
  variants, including when the advocate's position was known. **Corrected**: the cognitive-bolstering
  result belongs to `10.1111/j.1559-1816.2001.tb02481.x`, not this paper.
- UFMCS Red Team Handbook v5 (2011) — escalation with a named destination, adopted for leg 3.
  **Doctrine, not evidence**: contains no empirical evaluation and says so.
- **Deliberately excluded**: semantic-direction stratification (`arXiv:2605.30150`). Measured to
  work, but its scope is creative-ideation quality proxies rather than task accuracy, and it needs a
  planning call — the "add a stage" shape the evidence warns against. A follow-up idea with its own
  benchmark, not a fourth leg.

**Follow-ups opened by this idea, not resolved in it.**

1. The skill/Go drift: `skills/parley-deck/SKILL.md:768-786` omits `## Refutation attempts`, which
   `internal/protocol/reviewartifact.go:41` rejects a review for lacking. The documented hand-driven
   path emits artifacts the enforced path refuses.
2. Design rounds have **no existence-or-shape validator**, unlike review rounds. `opencode-1` exited
   0 with no artifact about ten times; only a human reading a directory listing stood between a
   phantom participant and the record. `parley consensus draft` did fail closed — three rounds after
   the artifact was due.
3. §15 in its entirety has **zero** Go enforcement (`/usr/bin/grep -rl` over `internal/` + `cmd/`
   for its vocabulary matches no `.go` file).
4. The finding-class vocabulary question, deferred to `opencode-1`, whose late round-01 artifact now
   proposes `REFRAME` + `## Frames considered`.
