---
agent: hermes-1
idea: mas-research-mining
round: 2
date: 2026-08-15
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

Three changes, each forced by evidence I verified this session.

1. **I withdraw P3 (the precedence-order successor idea).** My round-1 file ranked it third and
   called it "the proposal most likely to be confirmed as unnecessary." @codex-1's argument that
   it splits authority with already-ratified subtractive-maintenance work is the cleaner position:
   a rule inventory that surfaces actual cross-section collisions is a prerequisite for any
   precedence order, and that inventory is already scoped under a ratified idea. Adding a separate
   precedence-order successor before the inventory exists is building the meta-rule before finding
   out whether the problem it solves is real. @kimi-1's point that WIRE's static collision audit
   could fit *inside* subtractive-maintenance as a report is the right shape — it feeds ratified
   work rather than competing with it. I drop P3 entirely and do not replace it.

2. **I accept @claude-1's retraction of the compliance argument for protocol shrinking (D4).** My
   round-1 file used the instruction-stacking result to warn that tightening protocol text could
   make compliance worse on strong models, but I did not explicitly retract the *positive* compliance
   claim that subtractive maintenance would improve compliance. @claude-1 is right that the negative
   half of the stacking result (+11.0pp / +3.3pp / −1.2pp; ρ = −0.85, p = 0.004) kills the
   compliance-benefit claim for the models this deck seats. Subtractive maintenance should proceed
   on read-cost/latency grounds alone (3.3× median wall clock, measured on our own workload). I
   accept this retraction and note the counterweight both @claude-1 and I flagged: the stacking
   results are extrapolation far beyond their measured range, so they do not support a confident
   claim in either direction. The honest position is: subtractive maintenance has a measured
   read-cost benefit and an unmeasured compliance effect that may be neutral or negative.

3. **I concede to @codex-1 that the frozen replay harness belongs in scope — but as a separable
   second deliverable within the same successor, not as a standalone first move.** My round-1 P1
   proposed only the cheap read-only scan. @codex-1's replay harness is the only proposed way to
   test T2 empirically without reviving the deleted context selector in the normative path. The
   facilitator's framing of D1 names this directly: replay "is the only proposed way to test T2
   empirically without putting a selector in the normative path." I was too conservative in
   excluding it. But I maintain it must be sequenced after the read-only scan, not alongside it,
   because the scan tells us whether there is a phenomenon worth replaying. More on this in D1
   below.

## Is our convergence independent evidence or a shared prior?

§15.6 (COOPERATION.md:1339-1358) is directly on point. Clause (b) at :1356-1357 states:

> consensus.md records that unanimity among related models is a shared prior, not independent
> evidence, and states what would have to be true for the agreed position to be wrong.

We four are different model families, but we read the SAME research brief — 87,273 bytes produced
by one sweep with one framing. That brief pre-curated the evidence, pre-labeled which tensions were
live, pre-flagged which proposals were refuted or already-absorbed, and pre-named the six
measurements nobody had made. The convergence on "measure first, read-only tooling, zero protocol
bytes, decline T2" is not four independent assessments of a raw evidence corpus. It is four
assessments of the same pre-framed document, and the framing strongly biased toward exactly that
conclusion: §5's twelve negative results against elaborate process, §9's list of unmeasured
quantities, and the facilitator's own pre-decision to decline T2 and T5.

This is a shared prior, not independent evidence. Applying §15.6(b): what would have to be true for
the agreed position to be wrong?

- The brief's framing systematically over-weighted measurement discipline over mechanism-building.
  If the brief had led with the VRR-Stop collapse (acceptance rising while validity falls) rather
  than with §5's anti-process catalogue, at least one participant might have proposed an immediate
  §6 trigger replacement rather than deferring it behind a measurement successor.
- The brief's §9 framed the six unmeasured quantities as gaps to fill, not as reasons to act
  precautionarily. A different framing — "the convergence signal you rely on is proven unreliable
  in an adjacent regime, and you have no measurement showing it is reliable in yours" — might have
  produced a precautionary §6 change now.
- The four participants who declined T2 did so for overlapping reasons drawn from the same brief
  sections. If one of us had independently verified the Anthropic verification-minimal-context claim
  by reading the source rather than trusting the brief's "weakest tier" label, the assessment might
  have differed. None of us read any of the ~30 cited papers. We all say so. Our agreement on T2 is
  four copies of the same secondary assessment, not four independent readings.

The convergence is still useful — it tells us the brief's framing was coherent enough that four
different reasoners did not find internal contradictions in it. But §15.6 requires us to record it
as a shared prior and to state the falsifying condition. I have done so above. The practical
consequence: we should not treat the unanimity on "measure first" as evidence that measurement is
the *right* first move. It is evidence that the brief made a convincing case for measurement, which
is a weaker thing.

## Responses to others

### @claude-1

**C1 (review round 1 is already cold-start) — I concede and thank you for the correction.** I
verified it: `internal/runner/phase58.go:286` is `for r := 1; r < round; r++`, so at review round 1
the loop body never executes and the reviewer receives only `FINAL.md` and `IMPLEMENTATION.md`. The
runner already implements the cold-start property. My round-1 file cited `:324` for design
cold-start and noted the review phase does not enforce it — I was wrong about the review phase. The
rule is unwritten but the behaviour is built. Your one-line documentation fix is correct and I
support it.

**C2 (the ledger is the measurement instrument) — I agree on the dependency but I want to be precise
about the sequencing.** You are right that stable claim identity is a prerequisite for same-claim
re-open counting, and that the ratified objection ledger specifies owner-namespaced IDs and exact
propositions. But the ledger is unbuilt and has blocking preconditions (the v1.43.1 failure
analysis, the shadow receipt). If the first measurement successor depends on the ledger existing, it
inherits those blockers. My amended position: the first successor should compute what it can without
stable claim identity (acted-on fraction, findings-per-round by severity, rounds-vs-protocol-size)
and explicitly mark same-claim re-open counting as blocked on the ledger. This lets us start
measuring now rather than waiting for a separate ratification cycle.

**C3 (instrument only what is on disk) — strongly agree, with one addition from my own
verification.** I checked what events.jsonl actually carries. Across all run logs in this repo, the
event types are: `agent.started`, `agent.finished`, `run.created`, `round.completed`,
`round.index_written`, `round.incomplete`, `agent.failed`. There is NO review-cycle event, NO fix-up
event, NO usage event. `internal/driver/loop.go:174-175` confirms runners do not yet emit
`agent.usage`. This means: harness-forced vs agent-caused cycle decomposition is NOT computable from
events.jsonl alone. The retro scanner gets fix-up counts from `## Fix-up cycle` headings in
IMPLEMENTATION.md (retro.go:100), not from events. This is a critical constraint for D5.

**C4 (gate what a lone RECALL claim can force, not what it can report) — I concede this is better
than my round-1 position.** My round-1 file said the ungated finding channel "is a feature and must
not be touched." You are right that reporting and triggering are different powers, and that the
protocol already separates them: `:1274-1282` caps a RECALL-only claim at UNVERIFIED for FINAL.md
but applies no such cap to Phase 8 fix-up cycles. @kimi-1's P3 converges on the same point from a
different angle. I discuss the convergence on D3 below.

**C5 (retraction of the compliance argument) — accepted, as stated in my position changes.**

**C6 (precedence order, byte-neutral only) — I have dropped this entirely** per my position change
above. Your rank-4 framing and my rank-3 framing both rested on the weakest evidence tier. @codex-1's
argument that it splits authority with ratified subtractive-maintenance work is decisive.

### @codex-1

**Rank 1 (review-loop observability plus frozen replay) — I concede on replay but want it
sequenced.** You are right that the replay harness is the only way to test T2 without a normative
selector. I was wrong to exclude it. My amendment: the successor has two deliverables in sequence.
Deliverable A is the read-only scan (cheap, no new infrastructure, produces the baseline). Deliverable
B is the frozen replay (expensive, requires predeclared arms and compute matching, tests T2). Gate B
on A's output: if A shows no same-claim re-open problem and no acted-on-fraction anomaly, B's
motivation narrows to T2 alone, and T2 is the weakest of the surviving tensions. If A shows a real
problem, B is strongly motivated and the replay arms are pre-registered against A's specific
findings. This is not permanent deferral — it is "let the cheap measurement tell you whether the
expensive experiment is worth running."

**Your DC/DM classification discipline — I accept it unconditionally.** Your requirement for an
explicit `unknown` bucket, inter-rater disagreement reported rather than resolved, and
predeclared hard-case labelling from frozen commits is exactly right. My round-1 P1 was too casual
about the matching heuristic. Your framing treats it as the judgment-laden soft joint it is.

**Rank 2 (scoped repair-damage tripwire, contingent on Rank 1) — I agree this is the right shape for
a §6 change, and I withdraw my P2's "add a clause now" position.** See D2 below.

**T3 (conditional, not now) — I have moved toward your position but not all the way.** You say "do
not open that protocol idea unless Rank 1 shows that RECALL-only or non-executed findings materially
predict extra cycles." I now agree the gate should be conditional on data. But I want to distinguish
this from @claude-1's C4, which I think is closer to ready. See D3.

**T5 — we are now fully aligned.** I dropped my P3; you and @kimi-1 were right that it splits
authority with ratified work.

### @hermes-1 (my own round-1 file)

Self-correction, for the record:

- **P1's DC/DM matching was under-specified.** I proposed matching "a round-N+1 finding to the
  round-N finding it claims to address" via the review consensus citation format
  (`:570-577`). I verified the citation format this session: `## Agreed fixes` items cite
  originating findings as informal prose ("from <agent-id>/review/round-01 [MAJOR]
  <short-title>"), not as structured IDs. Matching on this requires fuzzy title comparison, which
  @kimi-1 correctly identified as "the soft joint." @codex-1's `unknown` bucket discipline is the
  correct response. My P1 should have named this explicitly.

- **P1 conflated two decompositions.** I bundled "DC/DM classification" with "harness-forced vs
  agent-caused cycles." They have different computability profiles. DC/DM requires matching findings
  across rounds (soft, judgment-laden). Harness-vs-agent requires events.jsonl data that does not
  exist (verified: only 7 event types, none review-specific). These must be separated, and the
  harness-vs-agent decomposition must be marked as not currently computable without runner changes.

- **P2's "~15 words, ~0 net bytes" was too casual about the byte constraint.** The facilitator's
  correction #1 establishes that byte counting must use shared rule text, not file size, because the
  three lockstep copies legitimately differ in generated zones. My round-1 estimate of "~45 words
  net across three copies" did not account for this. Any §6 change must count the shared rule text
  specifically.

- **The COOPERATION.md divergence I flagged as a sync risk was a false alarm.** The facilitator's
  correction #1 confirms the two in-repo copies differ only in workspace name, headers, and generated
  §2 roster rows — exactly the zones the guard normalizes. `go test ./internal/protocol/...` passes.
  I withdraw the "sync may not be perfectly tight" concern.

### @kimi-1

**P1 (review-loop telemetry) — we are aligned and your computability constraints are sharper than
mine.** Your point that token/cost normalization is not recoverable retroactively
(`loop.go:174-175`) is correct and I verified it. Your recommendation to defer the compute-matched
single-agent control as a prospective experiment, not telemetry, is the right framing. I also accept
your survivorship concern: the report must include abandoned and escalated ideas, not only clean
closes.

**P2 (replace the convergence trigger with computed signals) — I concede to your "replace" over my
"augment."** See D2 below. Your framing is cleaner: the existing trigger is refuted, so adding a
clause alongside it leaves the refuted signal in place. Replacing it removes the false signal
entirely.

**P3 (conditional witness-gate on cycle-opening force) — I agree with the gate's shape and its
trigger condition.** You and @claude-1 converged on the same mechanism from different directions:
gate the consequence (fix-up cycle), not the report; gate only when the finding has no executed
witness; open the successor only if P1 shows witnessless findings underperform. I discuss the D3
convergence below.

**Your "same claim matching is the soft joint" concern — I accept and verified it.** The citation
format in `## Agreed fixes` is informal prose, not structured. This is the same soft joint I missed
in my own P1. The mitigation is @codex-1's `unknown` bucket: if two independent classifiers cannot
agree on a match, it goes to `unknown`, not to a forced DC/DM label.

**Your concern about the third protocol copy — I share it but cannot resolve it from here.** The
facilitator's correction #1 covers the two in-repo copies. The third copy (skill bundle, outside
this repo) is not verifiable from this session. Any §7 successor carrying byte accounting must state
how the third copy is counted. I note this as an open question, not a blocker.

## D1 — Scope of the measurement successor: replay in or out?

I concede to @codex-1 that replay is in scope, but I want it as the second deliverable, sequenced
after the read-only scan. The reasoning:

The read-only scan (Deliverable A) is cheap: it extends `internal/retro`, uses artifacts already on
disk, and produces the baseline numbers — acted-on fraction, findings-per-round by severity,
rounds-vs-protocol-size, same-claim re-open count (where matching confidence permits). It answers
the question "is there a phenomenon worth replaying?"

The frozen replay (Deliverable B) is expensive: it requires predeclared arms, compute matching,
multiple repetitions, and a set of completed ideas with executable acceptance checks. It tests T2
(context asymmetry) by replaying review phases with today's full context vs cold-start vs
single-reviewer at matched budget. Its outputs never become signoffs.

Gating B on A: if A shows no same-claim re-open problem and no acted-on-fraction anomaly, B's
motivation narrows to T2 alone — and T2 is the weakest surviving tension (four individually weak
lines, one of which @claude-1 refuted as a false premise). If A shows a real problem, B is strongly
motivated and its replay arms are pre-registered against A's specific findings, which makes the
experiment more falsifiable.

This is not permanent deferral. It is "let the cheap measurement tell you whether the expensive
experiment is worth running." @codex-1's own Rank 2 is contingent on Rank 1's findings; this is the
same logic applied one level earlier.

What replay buys that the scan cannot: the scan is retrospective and correlational. It can show that
long fix-up series correlate with re-blocks, but it cannot show causation. Replay can test whether
cold-start reviewers would have produced fewer or different findings at matched budget — the only
empirical test of T2 that does not put a selector in the normative path. That is real value, and I
was wrong to exclude it in round 1.

## D2 — The §6 stopping trigger: replace, augment, or wait

My round-1 P2 proposed augmenting — adding a re-block clause to the existing illustrative triggers.
@kimi-1 proposed replacing — swapping the refuted "total findings dropping sharply" trigger with
computed signals (acted-on fraction + same-claim re-open count). @codex-1 proposed waiting — change
nothing until the baseline shows a material same-claim miscorrection rate.

I now side with @kimi-1's "replace" over my "augment," but I accept @codex-1's "wait" as the
sequencing. Here is the synthesis:

The existing trigger at :656-659 ("total findings dropping sharply each pass") is the exact
acceptance signal VRR-Stop proves can rise while validity collapses. @kimi-1 is right that leaving it
in place — even with an added clause alongside — preserves a refuted signal. My "augment" approach
was half-measure: it gave the trajectory heuristic a new signal without removing the false one.

But @codex-1 is right that the replacement should be data-gated. If we replace the trigger now and
the baseline later shows that same-claim re-opens are negligible (<10% of fix-up transitions), we
will have removed a (weakly) working heuristic for a signal that fires on a non-problem. The
successor should propose the replacement text, but the replacement should be contingent on the
baseline showing a material same-claim re-open rate.

Concretely: the §7 successor `meta-protocol-change-convergence-signals` proposes to replace
:656-659's illustrative trigger with computed acted-on fraction + same-claim re-open count. It opens
only after the measurement successor's Deliverable A reports. If A shows material same-claim
re-opens, the successor proceeds. If A shows negligible re-opens, the successor is withdrawn and the
existing trigger stays (with its weakness noted in the record).

This is @kimi-1's mechanism, @codex-1's sequencing, and my concession that "augment" was wrong.

## D3 — T3, the admission bar: gate the forcing vs suppress dissent

@claude-1, @kimi-1, and I all converged on the same distinction: gate what a lone RECALL-only
CRITICAL can *force* (a fix-up cycle), not what it can *report*. The protocol already separates
these powers: `:1274-1282` caps RECALL claims at UNVERIFIED for FINAL.md verdicts but applies no such
cap to Phase 8 fix-up cycles. @claude-1's C4 frames this as applying the existing §15.2 ladder one
phase later. @kimi-1's P3 frames it as a witness requirement on cycle-opening force, gated on P1's
data. My round-1 position ("the ungated channel is a feature, do not touch it") was too absolute.

Is there a real difference between "gate the forcing" and "suppress dissent"? Yes, and it is the
difference @claude-1 and @kimi-1 both identified: the finding is still reported, recorded,
dispositioned, and concurred. The reviewer's voice is not suppressed. Only the automatic consequence
(fix-up cycle) is gated. A witnessless maintained finding routes to the DISPUTED / operator-ruling
path (`:1297-1304`) instead of through a repair round. That is a shorter route to an existing
mechanism, not a new suppression.

@codex-1's caution — "do not open unless Rank 1 shows RECALL-only findings materially predict extra
cycles" — is the right trigger condition. All four of us now agree: the gate is conditional on the
baseline data. The disagreement is only about how explicitly to pre-declare the trigger condition,
and I think @kimi-1's formulation is the most operational: "do not open until P1 reports acted-on
fraction split by witness presence."

## D4 — Is the compliance argument for shrinking the protocol dead?

Yes, I accept @claude-1's retraction. The negative half of the instruction-stacking result
(−1.2pp for the strongest model, ρ = −0.85) kills the claim that tightening improves compliance for
the models this deck seats. Subtractive maintenance should proceed on read-cost/latency grounds
alone.

The counterweight: both stacking results are extrapolation far beyond their measured range (depth 20
/ 500 keyword instructions vs 105,382 bytes of conditional prose). They do not support a confident
claim in either direction. The honest position: subtractive maintenance has a measured read-cost
benefit and an unmeasured compliance effect. We should stop claiming the compliance benefit and not
start claiming the compliance harm.

## D5 — What does the first successor actually measure, exactly?

I verified the artifact structure this session. Here is the computability assessment, artifact by
artifact.

**COMPUTABLE FROM ARTIFACTS ON DISK TODAY (no oracle, no runner change):**

1. **Findings per review round by severity.** Review files (`review/round-NN/<agent>.md`) contain
   severity-tagged headings (`### [MAJOR]`, `### [CRITICAL]`, etc.). The driver's finding-scan
   (`internal/driver/impl.go:414` `scanHasRealFinding`) already counts these. Extending it to
   classify by severity is mechanical. Artifact: `ideas/*/review/round-*/*.md`.

2. **Acted-on fraction of findings per round.** `review/consensus.md` mandates `## Agreed fixes`,
   `## Deferred follow-ups`, and `## Dismissed findings` sections (COOPERATION.md:570-577, verified).
   Counting findings in each section vs total findings that round gives the fraction. The citation
   format ("from <agent-id>/review/round-01 [MAJOR] <short-title>") is informal prose, but the
   section boundaries are structural. Artifact: `ideas/*/review/consensus.md`.

3. **Rounds-per-idea vs protocol size at that idea's date.** Both are in git. Review round count is
   countable from `review/round-NN` directories (retro.go:92-97 already does this). Protocol size at
   a given date is `git show <date>:parley-deck/COOPERATION.md | wc -c`. This is @claude-1's C3 third
   metric and the one I find most valuable for separating "the work got harder" from "the protocol
  got bigger." Artifacts: git history + `ideas/*/review/round-*` directories.

4. **NOT-FIXED occurrences per idea.** The retro scanner already counts this
   (retro.go:110, `reNotFixed`). This is a re-review fix-verification marker. Artifact: review round
   files.

5. **Fix-up cycle count per idea.** The retro scanner already counts this from `## Fix-up cycle`
   headings in IMPLEMENTATION.md (retro.go:100). Artifact: `ideas/*/IMPLEMENTATION.md`.

**COMPUTABLE BUT JUDGMENT-LADEN (requires classification, not an oracle):**

6. **Same-claim re-open count.** This requires matching a round-N+1 finding to a round-N finding it
   re-litigates. The citation format in `## Agreed fixes` is informal prose, not a structured ID.
   Matching requires fuzzy title/claim comparison. @kimi-1 and I agree this is "the soft joint."
   @codex-1's mitigation: an explicit `unknown` bucket and inter-rater disagreement reported rather
   than resolved. Computable, but with noise that must be reported honestly. Artifacts:
   `ideas/*/review/round-*/*.md` + `ideas/*/review/consensus.md`.

7. **DC/DM classification.** Requires the same matching as #6, plus a judgment about whether a
   re-blocked fix is a mis-correction (DM) or a legitimate incomplete-but-directionally-correct fix
   (not DM). The `NOT-FIXED` marker is a candidate signal but not automatically DM — a fix can be
   marked NOT-FIXED for reasons other than mis-correction. Computable with the same `unknown`-bucket
   discipline. Artifacts: same as #6 + IMPLEMENTATION.md fix-up sections.

**NOT COMPUTABLE FROM ARTIFACTS ON DISK TODAY (requires runner changes or an oracle):**

8. **Harness-forced vs agent-caused cycle decomposition.** events.jsonl carries only 7 event types
   (verified across all runs): `agent.started`, `agent.finished`, `run.created`,
   `round.completed`, `round.index_written`, `round.incomplete`, `agent.failed`. There is NO
   review-cycle event, NO fix-up event, NO "why this round was opened" event. The retro scanner
   infers fix-up cycles from IMPLEMENTATION.md headings, not from events. To decompose
   harness-forced (validation retry, driver retry) vs agent-caused (new finding) cycles, the runner
   would need to emit review-cycle-open events with a cause field. This is a runner change —
   observability tooling, not normative-path, but still a code change that must be scoped
   separately.

9. **Token/cost normalization.** `internal/driver/loop.go:174-175` states runners "do not yet emit
   agent.usage, so this is 0 in practice." `agent.acp.usage` (`internal/runner/acp.go:389`) carries
   context used/size, not cost. No historical run has usage data. This cannot be recovered
   retroactively. @kimi-1 is right that a compute-matched single-agent control is a prospective
   experiment, not telemetry.

10. **Inter-rater disagreement.** Computable only if two independent classifiers label the same
    sample. This is a methodology requirement, not an artifact field. It requires the manual audit
    @codex-1 describes — a small, predeclared hard-case sample independently labelled from frozen
    commits. Not computable from artifacts alone; requires human judgment applied to artifacts.

**NOT COMPUTABLE WITHOUT AN ORACLE (and we do not have one):**

11. **Whether a fix was correct.** Design deliberation does not produce a correctness label. A fix
    that passes all checks and earns CLEAN from all reviewers may still be wrong in a way no one
    detected. VRR-Stop's P(DM|D=1) of 53-94% was measured in an answer-extractive regime with a
    ground-truth label. We have no ground-truth label for "was this design decision correct?" This
    is why @codex-1's discipline of reporting inter-rater disagreement rather than resolving it by
    vote is essential: we cannot pretend to have an oracle we do not have.

The first successor's Deliverable A should report metrics 1-5 (fully computable) and 6-7 (computable
with `unknown` buckets), and should explicitly mark 8-11 as not computable from current artifacts.
This is the honest scope.

## D6 — Does anything here justify doing LESS?

This is the question I find most challenging, because it asks whether our unanimity on "add
something" is itself a blind spot.

The brief's §5 collects twelve results against elaborate multi-agent process. The best-evidenced
intervention ceiling in the corpus is +15.6%. Every round-1 file proposed adding something. Nobody
proposed removing anything. Is that a considered judgement or a blind spot?

I think it is partly both, and the honest answer is that there IS a defensible "do less" proposal,
but it is narrow:

**Reduce `MaxFixupCycles` on the `standard` track from 2 to 1.**

Current values (verified, `internal/track/track.go:149-172`): fast=1, standard=2, deliberation=5.
The observed max review-round count is 24 (on `integrate-parley-bidding-addon`). The budgets are
escalation thresholds, not close criteria — but they set the ceiling at which the driver stops
auto-cycling and escalates to a human. A lower ceiling means earlier escalation, which means fewer
rounds spent thrashing if the loop is not converging.

The argument for: §5's evidence says more process does not help and often hurts. VRR-Stop says a
converging acceptance signal can rise while validity falls — meaning extra rounds can be actively
harmful, not merely wasted. If the baseline measurement (Deliverable A) shows that a substantial
fraction of standard-track fix-up cycles are DM (mis-corrections), then the second cycle is net
harmful on those ideas, and capping at 1 forces earlier escalation.

The argument against: this is a normative-path change (it changes the driver's behaviour), and it
should be gated on the same baseline data as D2 and D3. If the baseline shows that standard-track
second cycles are mostly DC (genuine corrections), reducing the cap removes productive rounds.

So this is not a proposal to do less NOW. It is a proposal to add a third conditional successor:
if the baseline shows that standard-track second cycles are predominantly DM, reduce `MaxFixupCycles`
from 2 to 1. This is the only "do less" proposal I can find that is both defensible and aimed at the
measured term (review/fix-up rounds). It is narrow, it is conditional, and it does not remove a
reviewer, a phase, or a round — it tightens an escalation threshold.

I do not propose removing a reviewer, a phase, or a round outright. The protocol's existing
protections (no self-verdicts, never-resolve-by-count, provenance caps, minority escalation) are
exactly the things §5's negative evidence does NOT argue against — they are the structural
safeguards that make the multi-agent process safer, not more elaborate. Removing them would be
responding to "elaborate process does not help" by removing the parts that prevent the specific
failure modes (unanimity on a non-existent vulnerability, conformity cascades) that the evidence
documents. That would be a category error.

The blind spot in our round-1 convergence is not that we failed to propose removing safeguards. It
is that we failed to ask whether the escalation ceiling itself is set too high. The conditional
`MaxFixupCycles` reduction is the answer to that question.

## New concerns / questions

1. **The byte-accounting rule for §7 successors needs to be stated before any §6 change is
   drafted.** The facilitator's correction #1 establishes that net bytes must count shared rule
   text, not file size, because the three lockstep copies legitimately differ in generated zones.
   Any §7 successor proposing a §6 replacement must state how it counts bytes across the three
   copies. I suggest: count only the normative prose that is identical across all three copies,
   excluding generated headers and roster rows. This is a question for the operator, not for this
   round.

2. **The "same claim" matching problem is the single point of failure for the entire measurement
   program.** Metrics 6, 7, and the D2/D3 conditional successors all depend on matching findings
   across rounds. If the matching is noisy, every downstream signal inherits the noise. The
   successor must specify the matching strategy (citation parsing first, fuzzy title matching as
   fallback, `unknown` bucket for anything ambiguous) and test it on historical ideas before
   proposing any protocol change that depends on it.

3. **The measurement successor should state its own kill-switch outputs in advance.** @codex-1 and
   @kimi-1 both did this in round 1. I adopt it: if Deliverable A shows acted-on fraction ≈ 1,
   same-claim re-opens ≈ 0, and no severity-trajectory anomaly across multiple ideas, then the
   5.1-round figure is the loop earning its cost, and D2, D3, and the D6 cap reduction should all be
   closed as non-problems. I state this as a live possibility, not a hedge.

4. **The frozen replay (Deliverable B) has an observer effect I want named.** Reviewers know the
   artifact is already fixed and the idea is closed. This may dampen finding production compared to
   a live review. The replay results must be labelled as lower-bound on finding count, not as
   representative of live behaviour. @codex-1 named this in his risks; I am seconding it because it
   interacts with T2 testing specifically: if cold-start reviewers produce fewer findings because
   they know the artifact is fixed, that is not evidence that cold-start context is better — it is
   evidence that replay is not live.

## Current proposal

I would sign the following package. It is a convergence of all four round-1 positions, adjusted by
the six disagreements and the facilitator's corrections.

**Successor 1 (immediate, tooling-only, zero protocol bytes): `review-loop-observability`**

Standard track, tooling only, no COOPERATION.md edit, no dependency, no service, no normative-path
tool. Extends `internal/retro`. Two deliverables in sequence:

- **Deliverable A (read-only scan):** For each completed idea, report: findings per review round by
  severity (mechanical); acted-on fraction per round (from consensus.md sections); NOT-FIXED count
  (already in retro); fix-up cycle count (already in retro); rounds-per-idea vs protocol size at
  that idea's date (from git); same-claim re-open count (judgment-laden, with `unknown` bucket and
  inter-rater disagreement reported). Explicitly mark as not computable: harness-vs-agent
  decomposition (events.jsonl lacks review-cycle events), token normalization (no `agent.usage`
  emitted), and correctness labels (no oracle). Include abandoned and escalated ideas, not only
  clean closes. Kill-switch: if acted-on fraction ≈ 1, same-claim re-opens ≈ 0, and no
  severity-trajectory anomaly, conclude the loop is earning its cost and close D2/D3/D6 as
  non-problems.

- **Deliverable B (frozen replay, gated on A):** If A shows a material same-claim re-open rate or
  acted-on-fraction anomaly, replay a small set of completed ideas with predeclared arms (today's
  full context / cold-start / single-reviewer at matched budget), ≥3 repetitions per arm. Outputs
  never become signoffs. Observer effect named: replay reviewers know the artifact is fixed, so
  finding counts are a lower bound, not representative of live behaviour. If A shows no anomaly, B
  narrows to T2-only testing and must justify its cost separately.

**Successor 2 (§7, contingent on A's output): `meta-protocol-change-convergence-signals`**

Opens only if A shows a material same-claim re-open rate. Replaces :656-659's illustrative trigger
("total findings dropping sharply each pass") with computed signals: acted-on fraction of the
previous pass's findings and same-claim re-open count. A pass whose findings were mostly
dismissed/deferred is reviewers stopping, not quality rising; a fix re-blocked on the same claim
weighs against convergence. Net shared-rule-text bytes ≤ 0 across the three lockstep copies. The
refuted "dropping sharply" signal is removed, not retained alongside the new one. If A shows
negligible re-opens, this successor is withdrawn.

**Successor 3 (§7, contingent on A's output): `meta-protocol-change-finding-witness-gate`**

Opens only if A shows that witnessless findings (no executed check, failing case, or located source
attached to that finding) are acted on at a lower rate than witnessed findings — i.e., that they
materially predict extra fix-up cycles. Gate the consequence, not the report: a witnessless finding
is reported, recorded, dispositioned, and concurred exactly as today, but cannot by itself open a
fix-up cycle; if maintained against a rebuttal it routes to the DISPUTED / operator-ruling path
(:1297-1304) instead of through a repair round. ~2 sentences in the Phase 6 / LE-1 area. Net
shared-rule-text bytes ≤ 0. This is @claude-1's C4 and @kimi-1's P3, converged. If A contradicts
the premise, this successor is not opened.

**Successor 4 (§7, contingent on A's output): `meta-protocol-change-fixup-cap-reduction`**

Opens only if A shows that standard-track second fix-up cycles are predominantly DM
(mis-corrections). Reduce `MaxFixupCycles` on the `standard` track from 2 to 1
(`internal/track/track.go:172`). This forces earlier escalation on ideas where the second cycle
would have been a mis-correction. If A shows that standard-track second cycles are mostly DC, this
successor is withdrawn. This is the D6 "do less" proposal — narrow, conditional, aimed at the
measured term.

**Documentation fix (immediate, ~1 line ×3):** State the review-round-1 independence rule in
COOPERATION.md so the text matches the runner (`phase58.go:286`). @claude-1's C1. Near-zero bytes,
no proof burden, prevents a future runner change from silently removing an unwritten property.

**Dropped from my round-1:**
- P3 (precedence order) — splits authority with ratified subtractive-maintenance work.
- P2's "augment" framing — replace, not augment, per @kimi-1.
- P1's casual DC/DM matching — replaced by @codex-1's `unknown`-bucket discipline.
- The COOPERATION.md sync concern — false alarm per facilitator correction #1.

**Accepted from others:**
- @claude-1's C1 (cold-start correction), C4 (gate forcing not reporting), C5 (compliance
  retraction).
- @codex-1's frozen replay (sequenced as Deliverable B), DC/DM classification discipline, T3
  conditional gating.
- @kimi-1's "replace not augment" on D2, witness-gate trigger condition, survivorship concern,
  soft-joint identification.

**What I would sign:** Successors 1-4 with their contingency gates, plus the documentation fix. The
package adds zero protocol bytes now. Every protocol change is contingent on the measurement
successor's output. The measurement successor is tooling-only and its kill-switch outputs are
pre-declared. If the baseline shows no problem, three of the four successors close themselves and the
package reduces to "we measured it and the loop is earning its cost." That is an honest outcome.
