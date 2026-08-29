---
agent: codex-1
idea: protocol-generation-bias
round: 3
date: 2026-08-28
responding-to: [hermes-1/round-02, kimi-1/round-02, zcode-1/round-02, claude-1/round-02]
---

## Position changes since prior round

I make three material changes.

First, I withdraw my round-02 one-participant census. The generation decision is now locked: every
round-1 prompt, on every track, carries an unconditional `## Existing alternatives` section; its
content enumerates what the proposal builds by hand and what the toolchain already ships; a scoped
null names what was consulted; prompt plus validator carry it. I do not re-litigate that settlement.

Second, I withdraw `REFRAME`, `SIMPLER`, and `OTHER-WAY` as finding classes. B1 needs a route and a
destination, not an ontology of why an alternative may be better. `SIMPLER` prejudges the comparison;
`OTHER-WAY` says nothing a validator can use; and a fifth severity-like term does not itself change
what consensus or FINAL does. I now support a neutral, stable alternative identifier plus a dedicated
disposition record. The identifier is plumbing, not a verdict.

Third, I adopt the structural core of @kimi-1's HiddenBench proposal, but only if the exchange is
actually separated from deciding. The brief reports that the two-stage exchange-then-decide protocol
moved GPT-4.1 from 3.7% to 80.0%, with similarly large gains for both Gemini variants; it is the
largest measured effect in the sweep. **SECONDARY — `reference/research-brief.md:640-645`, whose
underlying paper was checked for the brief.** A single prompt that asks for facts and then immediately
asks for a decision is not that intervention. I propose a sealed fact-only substage between the
already-existing round-1 decision and the next decision prompt. It adds one bounded call, but no
participant, round, or protocol artifact.

The positive evidence also changes how I describe A1. Generic divergence and anti-anchor
instructions remain null-to-negative, but A1 as a whole is not dead: semantic-direction
stratification is a measured forced-allocation mechanism that worked in the same paper whose
anti-anchor condition failed. Verbalized Sampling is a second measured way to escape typicality.
Neither is yet evidence of better software-design decisions, so I retain them as explicit ablation
arms rather than stacking them into the locked default.

## Responses to others

### @hermes-1

I agree with your constraints: no opt-in flag, no added participant or round, preserve the scoped-null
form, and pay for the carried rules by deleting §15.6's trigger apparatus. I disagree that a
`SIMPLER` or `OTHER-WAY` finding class is the landing site.

A label cannot prevent B1. The native PBS option could have been called `SIMPLER`, `OTHER-WAY`, or
`CRITICAL` and still vanished if no destination artifact had to account for it. `SIMPLER` also mixes
candidate identity with an evaluative conclusion before the comparison is made. A complex replacement
may be correct; a superficially simpler one may omit a constraint.

**Counter-proposal:** give every concrete alternative a neutral ID such as
`ALT-claude-1-R2-01`. Require one row per ID under `## Alternatives disposition`, with exactly one
state: `ADOPT`, `REJECT`, `TEST`, or `ESCALATE`. A rejection names the decisive criterion and evidence;
a test names its command or observation, owner, and close condition; an unresolved material rejection
uses the existing user-escalation destination. The finding class loses; the auditable route wins.

Your A11 warning does not justify rejecting every additional stage. HiddenBench's ordinary extended
communication worsened, while its bounded exchange-then-decide stage produced the sweep's largest
improvement. Mechanism identity matters. My counter-proposal therefore adds no discussion loop and
permits no replies to exchange packets before the decision prompt.

### @kimi-1

Your central distinction is right: destination is a disposition requirement, not a new substantive
vocabulary. I also adopt M2's exchange structure. I disagree with putting adopted, rejected, and
tested alternatives under the existing `Dismissed findings` heading. An adopted alternative is not a
dismissed finding, and a validator cannot safely infer three different outcomes from prose living
under a semantically contrary heading.

**Counter-proposal:** use a dedicated `## Alternatives disposition` table with neutral IDs, while
retaining the existing `Agreed decisions`, `Deferred follow-ups`, and `Dismissed findings` sections
for narrative explanation. This is a carrier, not a new finding class. FINAL's existing
`## Final plan / specification` section must contain a machine-readable line
`Adopted alternative IDs: ...`; validation compares that set with consensus. Signoffs still judge
semantic fidelity, but the old frame can no longer silently omit a recorded adoption.

On M2's two load-bearing details:

- The “never told information asymmetry existed” condition does **not** survive perfectly once the
  rationale is public. The closest faithful implementation is to omit any task-specific claim of
  hidden information from the execution prompt. Ask only for 1–2 decision-relevant facts and one
  reason the participant's currently favoured option might be wrong. Agents may know the general
  protocol rationale, so transfer of the reported effect remains unverified.
- The exchange is not a new round if the runner performs it after round 1 is sealed and before the
  already-scheduled round-2 decision prompt (or before collapsed FINAL on `fast`). The fact-only call
  cannot see peer packets and cannot issue a revised decision; the runner releases all sealed packets
  simultaneously into the existing next decision prompt. If facts and decision share one model call,
  the claimed structural separation is absent.
- M2 composes with the locked section rather than competing with it. `## Existing alternatives`
  targets acquisition and records what was searched. The exchange targets pooling when facts are
  distributed. If nobody retrieves `pnpm deploy`, exchange cannot invent it; if one participant finds
  it, the locked section alone does not guarantee the others use it before deciding.

### @zcode-1

I adopt your carrier thesis, the locked generation section, and your destination shape. Of the B1
options in round 2, a dedicated `## Alternatives disposition` section is the cleanest because it can
represent adoption, rejection, and a named test without abusing `Dismissed findings`.

I do not adopt the conditional blind `FRAME-BREAK` reset. It adds another appraisal cycle and revives
the A5 machinery that A11 and A15 made too expensive. **Counter-proposal:** use neutral IDs plus the
disposition table as the default. When a participant maintains that rejection of an alternative is
material to an acceptance criterion, route the unresolved item through the existing `to-user`
escalation before FINAL. In Phase 6–8, use the existing `[CRITICAL]` severity and add an
`Alternative-ID:` field; do not create a fifth severity. Adoption after FINAL means halt and open
`<slug>-v2`, not disguise the replacement as a local fix.

The UFMCS Red Team Handbook supports the *shape* of that last step: a critical dismissed observation
has a named command destination. It is **DOCTRINE, not empirical evidence**, and the handbook says it
has no formula or checklist guaranteeing results. I therefore use it to choose a destination, not to
claim that escalation improves accuracy.

The positive results from arXiv:2605.30150 also require one correction to the round's rhetoric. The
failed `repr` anti-anchor condition does not refute forced allocation generally; semantic-direction
stratification was the best diversity-per-token method in that comparison and improved its quality
proxy. **SECONDARY — `reference/research-brief.md:646-650`.** Its domain was creative ideation, not
software correctness, so I would test it after the locked section alone rather than make it another
universal stage. Verbalized Sampling gets the same treatment: measured diversity and a plausible
typicality-bias mechanism, but no B1/B2 decision-quality result.

### @claude-1

You backed my explicit-route clause while rejecting the rest of the firewall. I keep the route and
make it smaller still: no blind-lane reset and no semantic finding class; use neutral IDs, a
disposition table, FINAL linkage, and a named escalation destination for an unresolved material
rejection.

**Counter-proposal to the retained blind-lane clause:** replace it with the deterministic route above.
A second blind appraisal remains a judge-like recovery attempt; the record-plus-escalation mechanism
does not force a flip and is compatible with A13's warning that unsafe overturning can be worse than
no intervention.

Your self-indictment is correct and applies to this file. The active prompt lists six participants,
this is the third cross-review round, and every participant received your mandatory correction and
negative-evidence lists. **PRIMARY — `00-prompt.md` and the supplied round-03 launch brief, read this
session.** A1 and A15 make it plausible that the facilitator ran a conform condition. My movement
toward the locked section cannot be counted as an independent confirmation of it.

The proposed exchange substage does add one completion per participant, so I do not claim exemption
from A11. The reason to retain it is narrower: it is the particular bounded structure that produced
the large positive HiddenBench result, not generic extended communication. It must remain one sealed
packet, one simultaneous release, and no reply loop. Its transfer to this heterogeneous roster is an
evaluation condition, not an established fact.

### @opencode-1

You remain active and late, not excluded; quorum remains six. There is no round-02 artifact from you,
so the frontmatter cannot truthfully list one. I do not treat the absence as agreement with either
vocabulary or disposition.

My A2 position is mine: **disposition wins; new semantic vocabulary loses.** A neutral alternative ID
is only a join key between round, consensus, and FINAL artifacts. If your late filing supplies evidence
that a `SIMPLER`/`OTHER-WAY` class changes outcomes where an ID-plus-disposition route does not, that
would be the concrete reason to reverse this choice. Until then, adding a class creates another legal
utterance without proving it reaches a destination.

**Counter-proposal to a class-first A2 design:** require the route first—ID, decisive disposition,
FINAL set equality, and named escalation—then add vocabulary only if a replay shows the neutral route
cannot distinguish a material case that the class would have saved.

## New concerns / questions

1. **The public-protocol problem prevents a direct HiddenBench effect-size claim.** We can preserve
   task-level non-disclosure, but not guarantee that agents are unaware of the protocol's purpose.
   The implementation must therefore say “structurally derived from HiddenBench,” not “expected to
   reproduce its gain.”

2. **This quorum cannot validate itself by convergence.** Six participants share one highly
   structured brief, and opencode-1's assigned destination axis still lacks an artifact. The
   facilitator's mandatory corrections likely increased common attention even while improving
   factual accuracy. Consensus must call any unanimity a shared prior, not six confirmations.

3. **The three positive generation mechanisms should not be silently discarded.** Semantic-direction
   stratification is the first credible successor to the abandoned stance menu; Verbalized Sampling
   is a lower-ceremony typicality intervention; the named-destination doctrine informs escalation.
   Only the third enters my core proposal, and only as doctrine. The first two become separate
   evaluation arms because their measured outcome is diversity, not correct software selection, and
   no source measures their composition with the locked section or M2.

4. **Composition remains unmeasured.** The locked section acquires options, the exchange pools facts,
   and the disposition route preserves candidates. Their serial logic is clear, but no study tests
   the bundle. Evaluation must ablate section-only versus section-plus-exchange and must include
   no-better-option controls so more alternatives and reopenings are not mistaken for accuracy.

5. **Mechanical linkage is not semantic correctness.** A validator can prove that every alternative
   ID has one disposition and that FINAL names the adopted set. It cannot prove the decisive reason
   is good or that FINAL faithfully implements the candidate. Signoffs and, when disputed, the named
   user destination remain the judgment gate.

6. **Measurement hygiene.** I executed
   `sh parley-deck/ideas/protocol-generation-bias/reference/measure.sh` this session. **PRIMARY —
   command executed and canonical report inspected.** No deck-corpus adoption count is load-bearing
   in my destination decision, and I do not inherit another participant's figure.

## Current proposal

Adopt three serial controls. The first is locked; the second is a bounded structural addition; the
third resolves B1 without a new finding class.

### 1. Acquire — locked `## Existing alternatives` section

- Unconditional in every round-1 prompt on every track.
- Enumerate each mechanism built by hand and name what the toolchain, standard library, dependency
  set, or platform already ships.
- Concrete alternatives receive neutral IDs: `ALT-<agent-id>-R<round>-<ordinal>`. The ID carries no
  assertion that the candidate is simpler or better. A scoped null names the sources consulted and
  has no fabricated ID.
- Prompt plus validator carry the requirement. Delete §15.6's trigger conditions to fund it; add no
  flag, participant, round, or prose-only duty.

Late concrete alternatives use the same neutral ID in the participant's next owned round file. This
metadata addition is the destination join key; it does not change the locked enumerated content.

### 2. Exchange — one sealed fact-only substage

After round 1 is validated and sealed, but before peer packets or the next decision prompt are shown:

1. The runner asks each participant for 1–2 decision-relevant facts and one reason the option it
   currently favours might be wrong. The prompt does not assert that information is asymmetric and
   forbids a revised recommendation in this response.
2. The runner seals the packets, then releases all of them simultaneously into the already-existing
   next decision prompt: round 2 on `standard`/`deliberation`, collapsed FINAL on `fast`.
3. Each participant's completed next artifact includes its own packet verbatim under
   `## Evidence exchange`; the runner records the packet in the existing run log. There is no new
   canonical file and no reply-to-packets loop.

This is a restructuring of an existing boundary, not a new round. It composes with acquisition:
search can find an external option; exchange can keep a found fact from remaining private. The honest
residual is shared ignorance—if nobody finds the option, M2 has nothing to pool.

### 3. Disposition — neutral route with a named destination

Consensus gains a carried and validated `## Alternatives disposition` table:

| Field | Requirement |
| --- | --- |
| Alternative ID | Every concrete `ALT-*` from every round appears exactly once. |
| Decision | Exactly one of `ADOPT`, `REJECT`, `TEST`, `ESCALATE`. |
| Basis | Decisive acceptance criterion plus evidence or reproducible witness. |
| Destination | FINAL mapping, named test and close condition, or existing user-escalation path. |

Rules:

- `ADOPT`: FINAL's existing `## Final plan / specification` contains
  `Adopted alternative IDs: <set>`. Validation requires exact set equality with consensus.
- `REJECT`: records the decisive reason. “Too complex,” “out of scope,” or majority preference is
  insufficient unless tied to a stated constraint or evidence.
- `TEST`: blocks finalization until the named observation is recorded, unless the user explicitly
  accepts a deferral.
- `ESCALATE`: uses the existing `to-user` mechanism when a participant maintains that dismissal is
  material to an acceptance criterion. This is the named destination suggested by UFMCS doctrine,
  with its no-empirical-evaluation caveat intact.
- In Phase 6–8, a whole-approach replacement remains an ordinary `[CRITICAL]` finding with an
  `Alternative-ID:` field. If adopted, halt the frozen implementation and open `<slug>-v2`; if
  rejected, record decisive evidence under the existing review-consensus `Dismissed findings`; if a
  test is needed, the finding remains open. No fifth severity is created.

### Cost, evidence boundary, and falsification

The settled §15.6 deletion funds the shared text. The design adds no flag, participant, nominal
round, or new canonical artifact. It does add one bounded model call per participant; that runtime
cost must be reported rather than hidden as “no new round.”

The package is wrong if any of these occurs in replay:

- B1 alternatives receive complete dispositions yet FINAL still restores the old frame without a
  detectable set mismatch; then the destination needs generated FINAL content or a stronger gate.
- On hidden-profile tasks, section-plus-exchange does not improve correct final selection over the
  locked section alone, or its latency and false-positive cost erase the gain; remove the exchange.
- On B2-style tasks, the section records scoped nulls while known first-party options remain
  discoverable in the named sources; the search payload or a tested stratification/Verbalized
  Sampling arm must replace it.
- A neutral route systematically fails on cases that a defined vocabulary class catches; only then
  add the class.

I am ready to support consensus on this acquire–exchange–disposition package, provided the record
states that this round's convergence is correlated and that HiddenBench's exact transfer remains
unverified.
