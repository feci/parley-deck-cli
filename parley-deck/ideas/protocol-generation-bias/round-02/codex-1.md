---
agent: codex-1
idea: protocol-generation-bias
round: 2
date: 2026-08-28
responding-to: [hermes-1/round-01, kimi-1/round-01, zcode-1/round-01, claude-1/round-01]
---

## Position changes since prior round

I **abandon the round-01 Anchor Firewall as the primary/default solution**. I retain only a narrow
A5 constraint—do not seed round 1 with a facilitator-authored solution unless the user made that
mechanism a constraint—but I withdraw staged frame receipts, per-participant hidden evidence calls,
pairwise order-swapped appraisal, and late-candidate reset as a package. I now carry **A2
reframe-vocabulary explicitly**, because its owner has not filed, and pair it with the smallest
fact-finding form of A4. This is a temporary argument for the missing axis, not a proxy artifact for
@opencode-1 and not an attribution of my proposal to that participant.

The evidence changed the decision:

- My original package tried to control anchoring with more structure and more participant
  completions. **SECONDARY — independently checked by claude-1's
  `reference/research-brief.md:A11,A15`.** HiddenBench found that seven agents improved only 0.6%
  pre-to-post discussion versus 34.8% for three, and longer communication eventually declined;
  Yun et al. found rigid chat templates suppressed diversity even at high temperature, while also
  warning that removing structure costs quality. Those findings do not prove that six heterogeneous
  frontier agents fail, but they directly indict my proposal's unmeasured assumption that two more
  structured completions per participant are harmless.
- An instruction to avoid anchors is not an adequate substitute. **SECONDARY — same checked brief,
  A1–A2.** Generic “be different” instructions were a measured null in Smith, Ward & Schumacher;
  in the LLM-specific Anchorless Diversification comparison, the shared-representative-anchor
  avoidance condition was the worst anchorless method, and adding a divergence instruction hurt it.
  That does not test actual incumbent removal, but it defeats the instructional half of my round-01
  design.
- Two narrow A5 facts survive. **SECONDARY — same checked brief, C4 and C10.** Removing the familiar
  incumbent changed human chess performance from 0–50% to 100% in the affected skill groups, while
  merely reporting a search did not distinguish solvers from non-solvers. Separately, Wang et al.
  measured LLM evaluation-order effects and used swapped-order calibration. The first supports not
  injecting optional solution examples; the second supports order control when comparison order is
  genuinely decision-relevant. Neither supports the full firewall or proves better software-design
  selection.
- My late-candidate reset also leaned too hard on an LLM judge-like second hearing. **SECONDARY —
  same checked brief, A13.** Minority Sentinel's LLM-as-Judge baseline had negative net gain. A
  deterministic requirement to *disposition* a reframe is safer than asking another judge to
  overturn a majority.

The citation audit found **no substantive defect in my five citations**. It confirmed Tversky &
Kahneman, Stasser & Titus, Rowe & Wright, Wang et al., and Qu, Fu & Hu. Bibliographic correction:
the full Stasser & Titus title includes “Biased Information Sampling During Discussion”; my round-01
title omitted that subtitle. My Rowe & Wright text did name all four Delphi features, including
statistical aggregation. **SECONDARY — claude-1's independent citation audit and
`reference/research-brief.md:38-68`.** I make no claim that those human results transfer to this
deck.

I also re-derived the disputed deck counts rather than inheriting another participant's number.
**PRIMARY — executed**
`sh parley-deck/ideas/protocol-generation-bias/reference/measure.sh` on 2026-08-28. Under the
script's frontmatter-key definition, `require_model_diversity` is **1/88 including this idea and
0/87 excluding it**, with one separate prose-only mention. Under the literal level-2-heading
definition, `## Adversarial alternative` occurs in **4 files spanning 3/89 idea directories**.
Those are the only adoption figures I use below.

I am not moving to A2 because three participants named it. That convergence shares a brief and is
not independent evidence. I am moving because B1, as stipulated, contains a candidate that already
existed before the old frame won. Once existence is held constant, a generation-only intervention
has no causal handle on the remaining failure; a route and a destination are logically necessary.
B2 then proves that the route is not sufficient. The benchmark pair requires two serial mechanisms.

## Responses to others

### @hermes-1

I agree with three parts of your null position: do not add another opt-in flag, pay for any surviving
rule by deletion, and do not pretend protocol prose is enforcement. The canonical measurement makes
the opt-in point more precise than round 1: the diversity key is 1/88 including this self-measuring
idea and 0/87 excluding it, not 0/88; the adversarial section is present in 3/89 ideas, not a generic
“about 5%” population.

Two citation corrections must travel with your argument. **SECONDARY — claude-1's independent audit.**
Cognitive bolstering belongs to Nemeth, Connell, Rogers & Brown, JASP,
`10.1111/j.1559-1816.2001.tb02481.x`; the EJSP paper `10.1002/ejsp.58` instead found the authentic
minority superior to all three devil's-advocate variants, including variants where the advocate's
true position was disclosed. Also, arXiv:2605.00914 is by **Bertalanič & Fortuna**, not Chen et al.
Its homogeneous-debate result survives, but your statement that divergence enforcement is another
failed debate round remains your inference, not the paper's finding.

I disagree that broadening and enforcing the existing steelman is enough. It asks for an argument
after an incumbent exists, while B2 may be an information-acquisition failure; and generic
alternative instructions have direct negative/null evidence. Deleting only the mechanically
decidable exclusion would therefore make a late, abstract duty universal without changing what is
searched or what happens to a found replacement.

**Counter-proposal:** replace §15.6 rather than merely deleting its carve-out. Move its one assigned
filer to blind round 1 and turn the task into an enumerated, locator-bearing existing-option census.
Then give any evidence-linked structural replacement a `REFRAME` identifier that consensus must
adopt, reject with decisive evidence, or send to a named test. This remains subtractive in shared
text and supplies the two things A6 concedes it cannot supply: B2 discovery work and B1 destination.

### @kimi-1

Your strongest components survive: independent generation before interaction, an accountable owner,
and a carried mechanical check. **SECONDARY — checked brief C8 and “positive evidence”.** Human
brainstorming meta-analysis strongly favors individuals generating alone, and semantic-direction
stratification is the best diversity-per-token LLM intervention in the cited comparison. Your
distinction between direction diversity and model-roster diversity is also important.

I do not accept the four fixed stances, four-way occupancy floor, semantic distinctness gate, and
dropped-candidates ledger as one default package. They spend the same scarce round-1 attention four
times; a fixed menu can become the anchor; rigid formats have a measured diversity cost; and adding
agents/rounds is specifically harmful in HiddenBench's hidden-profile regime. The evidence also says
over-generation saturates: 4,000 ideas deduplicated to about 200 in one study. None of that proves
your four candidates fail, but it removes the presumption that more occupied slots mean more useful
coverage. DCI is the sharp warning: it produced decision packets reliably, yet the full protocol lost
to a single agent overall and cost about 62 times as many tokens. **SECONDARY — checked brief C9,
A10, A11, A15.**

The authentic-dissent evidence does not directly refute direction assignment, because `REUSE` is a
search direction rather than a performed opposing belief. It does refute treating nominal occupancy
as authentic disagreement. A file can wear `REPLACE` while remaining inside the incumbent frame.

**Counter-proposal:** allocate only one bounded `existing-option census` to one current participant;
leave all other round-1 proposals free-form. The census checks first-party docs/CLI, existing
platform/dependency capability, and delete/no-change, with locators or a scoped null. Do not add a
semantic distinctness gate or a ledger for every candidate. Preserve only evidence-linked structural
reframes, using an ID and the existing consensus/review destination sections. This retains your
accountability and isolation while avoiding fourfold mandatory generation.

### @zcode-1

Your carrier diagnosis survives and becomes a hard constraint on my proposal. **PRIMARY — local code
inspection.** `internal/runner/runner.go:821-871` embeds the whole `00-prompt.md` in
`BuildRoundOnePrompt`; `internal/runner/phase58.go:240-249` carries the review posture and allowed
severities; `internal/protocol/reviewartifact.go:17-54` rejects a review without a non-empty
`## Refutation attempts` and `reviewed-commit`. Conversely,
`internal/app/driver_consensus.go:112-133` currently asks for `Agreed decisions`, `Trade-offs
accepted`, `Deferred follow-ups`, and `Dismissed findings` but carries none of §15.6. A ratified rule
that is absent from prompts and validators is not an operational intervention.

Your evidence corrections are material. **SECONDARY — claude-1's independent audit.** The Jansson &
Smith DOI is `10.1016/0142-694X(91)90003-F`, not `...90011-F`; its raw gaps are powerful, but it
reports no p-values, confidence intervals, or effect sizes and has cells of 6–18, so it must not be
called statistically significant. Smit et al. do say untuned MAD systems do not reliably outperform
other prompting strategies, but the appended “dissent tends to be shallow” was unsupported, and the
omitted qualification says tuned systems can improve and even surpass the non-debate protocols
tested. That paper cannot carry a blanket anti-debate conclusion.

I disagree with a mandatory `## Existing alternatives` section in *every* round-1 file. It is a
generic instruction repeated across all participants, exactly where A1 and A15 warn about null
divergence instructions and rigid templates. The engineering precedent also shows why presence is
not enough: 42% of 657 Kubernetes KEPs had no substantive Alternatives content despite the template.
**SECONDARY — checked brief A1, A15 and engineering precedent.** A non-empty-section validator would
make artifact production reliable, not make the search competent.

**Counter-proposal:** keep your default-on carrier and validator, but require the census in exactly
one assigned round-1 file, with three enumerated search scopes and locators. Let other participants
remain structurally free. Add a deterministic scan for `REFRAME` IDs and their consensus
dispositions. This is more enforceable than prose, less repetitive than a universal heading, and it
addresses your admitted B1 boundary rather than stopping at earlier timing.

The canonical re-measurement also changes your round-1 adoption premise: the literal adversarial
heading spans **3/89 ideas and 4 files**, not 5/89. The carrier conclusion remains; the corrected
number strengthens it.

### @claude-1

Your revision from “appoint a stance” to “appoint a question” is the right part of A4 to keep. The
citation audit strengthens the rejection of a devil's advocate: authentic dissent beat all three
role-play variants even when true position was transparent. It also requires separating that EJSP
result from the JASP cognitive-bolstering result. **SECONDARY — your independent audit.** A
locator-bearing census is fact-finding, not a command to oppose, so it does not claim an exemption
from Nemeth; it changes the task being appointed.

I reject the premortem half. **SECONDARY — checked brief A8.** The famous “30% more reasons” source
could not be opened, and the controlled premortem outcomes measured confidence rather than decision
change. Group premortem was no better than individual premortem (`p=.772`). It does not earn shared
protocol bytes here.

I also withdraw any support my earlier position might have lent to “Anchorless Diversification” as
confirmation of peer-anchored avoidance. **SECONDARY — checked brief A2.** Its representative-anchor
avoidance condition was the worst anchorless method, and the added divergence instruction hurt it.
The positive result was semantic-direction stratification, not “tell agents to avoid shared anchors.”

Your seventh-role problem has a smaller answer: use one of the already-required participants, in
round 1 before peer proposals are visible, and validate the participant's census section. The
appointment may still be shallow; locators make that failure inspectable but do not cure competence.

**Counter-proposal:** keep only the missing-option census; drop the premortem; perform the census in
the existing blind round rather than after proposals; enumerate first-party/toolchain,
already-present platform/dependency, and delete/no-change search scopes; and pair every resulting
structural candidate with the A2 disposition route below. That pair hits B2 and B1 without a new
agent, round, artifact, or flag.

### @opencode-1

There is no round-01 artifact to answer. I do **not** read that absence as agreement, withdrawal, or
evidence against A2. The owner explicitly kept quorum at six and directed the group to wait for the
backend, so `opencode-1` remains late. The `responding-to` frontmatter therefore lists only the four
round-01 artifacts that exist, while this heading preserves the required explicit response.

I am carrying A2 openly because otherwise the round's structurally necessary argument would be
missing. My provisional vocabulary is not a fifth review severity. In design rounds the exact marker
is `### [REFRAME] RF:<agent>:<round>:<ordinal> — <candidate>`; in implementation review it is the
existing blocking severity plus a subtype:
`### [CRITICAL] REFRAME RF:<agent>:<review-round>:<ordinal> — <candidate>`. The body must name a
materially different mechanism, a source or reproducible witness, and the acceptance criterion it
could change.

The route is deterministic. Before design consensus closes, each ID must appear exactly once under
`Agreed decisions` (adopt), `Dismissed findings` (reject with decisive evidence), or `Deferred
follow-ups` (named test and owner). After `FINAL.md` freezes, an accepted reframe is not a “fix”: the
current implementation halts and the destination is `<slug>-v2`; a rejected one goes under
`Dismissed findings` with the decisive evidence. A mechanically decidable claimed advantage cannot
be rejected with rhetoric alone; cite the test result or route it to the named test.

This is my concrete provisional proposal for the absent axis. A late A2 artifact may replace,
improve, or reject it; nothing here claims to be @opencode-1's position.

## New concerns / questions

1. **B2 is not automatically a hidden-profile problem.** If one participant knows a native option
   and the group fails to pool it, HiddenBench applies directly. If nobody retrieves or recognizes
   it, the failure is information acquisition, not distribution. The census changes acquisition;
   the separated fact-before-decision lane changes pooling. Conflating them would make a successful
   exchange protocol look like a search engine.

2. **This idea's six-agent design is conditionally indicted by its own evidence.** HiddenBench does
   not establish that six heterogeneous frontier models are worse than three, but it removes any
   right to count six agents as six independent confirmations. Here all six received the same
   high-structure brief, the supposedly decisive axis has no artifact, and the other five cannot
   repair that by repetition. The present proposal adds no agent and no round; future evaluation
   should compare three-participant and six-participant runs rather than assume scale helps.

3. **The census list may itself anchor.** The three scopes are deliberately about *where to look*,
   not three required solution shapes, and only one participant receives the rigid section. Still,
   the enumerated-blacklist evidence is about avoiding named features, not about finding software
   options. Mapping it to search scopes is an inference. A wildcard/free proposal lane must remain.

4. **Artifact validity is not decision quality.** The driver can prove a section and disposition
   exist; it cannot prove a search was competent or a rejection was correct. Locators and named tests
   make defects reviewable. They do not automate judgment. The DCI result is the warning against
   claiming more.

5. **Do not use position stability as the outcome.** **SECONDARY — checked brief A14.** Thirty-seven
   percent of observations flipped under self-reflection alone, and reducing peer adoption cannot
   distinguish harmful from beneficial influence. Success must be known-option recall and final
   decision quality, not fewer position changes, lower consensus, or more dissent artifacts.

6. **False-positive alternatives need measurement.** Sherbino et al.'s debiasing trial included a
   false-positive arm; an option-forcing protocol should likewise measure unnecessary reopenings and
   inferior native-option proposals, not only recovered misses. **SECONDARY — checked brief A5.**

7. **What would make me conclude this position is wrong:**
   - In a preregistered B2-style replay set where each task has a known first-party/native route, the
     one-person census does not improve pre-discussion native-option recall over ordinary independent
     round 1, or its false-positive/cost burden erases the gain. Remove the census.
   - In B1-style replays where a known-better structural candidate arrives late, the `REFRAME` gate
     does not reduce silent frame restoration or wrong final selection relative to an ungated
     consensus. Remove or redesign the route; a recorded corpse is not success.
   - A carried, universal `## Existing alternatives` section matches or beats the one-person census
     at lower token and false-positive cost. Prefer @zcode-1's simpler carrier.
   - Suppressing facilitator-authored solution examples reduces solution quality without improving
     known-option recall or candidate-family coverage. Abandon the remaining A5 constraint.
   - A single mechanism passes both B1 and B2 replays. Then my two-link claim is over-designed.

## Current proposal

### Decision and axis budget

Adopt a **Search–Route pair**, not an omnibus anti-bias protocol.

| Axis | Function | Budget decision |
| --- | --- | --- |
| A5 | Prevent an optional facilitator solution from becoming the round-1 incumbent | Keep one sentence; abandon staged firewall and late blind reset |
| A4 / A1 | Make one external/native option search happen before peer exposure | One assigned census, not four stances or N mandatory candidates |
| A3 | Carry and validate the duties | Required implementation constraint, not a separate participant ceremony |
| A2 | Give a structural candidate a route and destination | Required second link; cannot be traded away for more generation |
| A6 | Force cost discipline | Delete and replace §15.6; no flag, agent, round, or artifact |

A1, A3, and A4 compete for the same round-1 attention budget when each proposes its own alternative
section, stance menu, or appointee. I choose one assigned census and A3's carrier. A2 is complementary,
not competing: it acts after a candidate exists. The full A5 reset and A2 disposition compete for the
B1 retention budget; the explicit route is smaller and more auditable, so A5 loses there.

### Mechanism 1 — search before proposals (B2)

1. Unless the user mandates a mechanism, the round-1 view of `00-prompt.md` contains outcome,
   constraints, acceptance criteria, and raw evidence—not a facilitator-authored candidate.
2. One existing participant is assigned `## Existing-option census` before seeing peer proposals.
   It checks exactly three scopes: first-party docs/CLI; capabilities already present in the platform
   or dependency tree; and delete/no-change. It files concrete candidates with locators, or a scoped
   null naming what was checked.
3. Everyone else remains free-form. No semantic-distinctness floor, candidate count, added agent, or
   extra round.
4. The runner carries the assignment in that participant's prompt and reuses
   `protocol.HasNonEmptySection`-style validation. Round 2 does not open until the assigned file has a
   non-empty census. Presence is the mechanical floor; cross-review evaluates substance.

This would have put a first-party `pnpm` docs/help lookup into B2 before hand-written scripts became
peer context. It does not guarantee recognition of `deploy`; failure to find it remains an observable
search failure rather than proof no native route exists.

### Mechanism 2 — route every material reframe (B1)

1. A design-round reframe uses `### [REFRAME] RF:<agent>:<round>:<ordinal> — <candidate>` and supplies
   a materially different mechanism, source/reproducible witness, and affected acceptance criterion.
   An objection, slogan, or restyled incumbent does not qualify.
2. The consensus prompt and validator scan all prior round files. Every ID appears exactly once in
   `Agreed decisions` (adopt), `Dismissed findings` (reject with decisive evidence), or `Deferred
   follow-ups` (named test, owner, and close condition). Signoffs remain the human judgment gate.
3. In Phase 6–8, use `### [CRITICAL] REFRAME <id> — <candidate>` so the existing severity scanner
   blocks a false clean close. Because `FINAL.md` is frozen, the only destinations are evidence-backed
   dismissal or a halted current implementation plus `<slug>-v2`; it is never laundered into a local
   fix or an unowned deferral.

In B1, the native S3 candidate would have received an ID. `FINAL.md` could retain `vzdump + rclone`
only after the consensus record named the candidate and supplied decisive evidence or a test. The
mechanism guarantees a hearing and destination, not that agents choose correctly.

### Smallest shared-rule change

Delete §15.6 and replace it with the following exact text; simplify the §15.7 row to `yes` on all
tracks. No new opt-in exists, because `fast` is where a small mechanically verifiable B2-like task is
likely to land.

```markdown
### 15.6 Search and reframe

On every track, a kickoff states outcomes, constraints, and raw evidence; it does not supply a facilitator-authored solution unless the user made that mechanism a constraint. Before peer proposals are visible, one existing participant files `## Existing-option census` in round 1. With locators, it checks first-party docs/CLI, capabilities already in the platform or dependencies, and delete/no-change. A null names the scope searched.

A materially different candidate is marked `REFRAME <id>` and names the candidate, a source or reproducible witness, and the acceptance criterion it changes. Consensus must put every id exactly once under Agreed decisions (adopt), Dismissed findings (reject with decisive evidence), or Deferred follow-ups (named test and owner). In review the heading is `[CRITICAL] REFRAME <id>`; it is not a fix to frozen `FINAL.md`: dismiss it with evidence or halt and route it to `<slug>-v2`. Prompts carry these duties and validators block a missing census or disposition.
```

**PRIMARY — local byte measurements.** The live §15.6 is **1,372 bytes**, re-derived with
`awk 'BEGIN{p=0} /^### 15\\.6/{p=1} /^### 15\\.7/{p=0} p{print}' parley-deck/COOPERATION.md | wc -c`
→ `1372`. The exact ASCII replacement above is **1,030 bytes**, so the shared-rule change is at most
**−342 bytes** before the shorter §15.7 row. Prompt, scanner, validator, and tests are enforcement
code, not a second copy of normative prose; generated protocol views mirror the same rule.

### Benchmark and evaluation contract

| Benchmark | Search link | Route link | Honest residual failure |
| --- | --- | --- | --- |
| B1: candidate exists, old frame wins | Already satisfied; generation is not the problem | ID plus adopt/reject/test disposition blocks silent absorption | Group can still reject the better candidate for a bad stated reason |
| B2: native option never appears | One blind, enumerated, locator-bearing census targets the missing lookup | Activates only if a candidate is found | Shared ignorance or shallow search can still miss `pnpm deploy` |

Before ratification, replay both benchmarks and a preregistered set of known-native-option and
no-better-option tasks. Measure native-option recall before discussion, oracle/human-rated final
decision quality, unnecessary reopens, tokens, and latency. Do not measure dissent count or position
stability. If both links do not earn their separate effects, delete the one that fails rather than
turning it into another permanent ceremony.
