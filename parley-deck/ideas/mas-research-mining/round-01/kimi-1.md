---
agent: kimi-1
idea: mas-research-mining
round: 1
date: 2026-08-15
---

## Summary

Three proposals, ranked, all aimed at the measured term (review/fix-up rounds 1.6→5.1, max 24 —
operator figure, SECONDARY; the distribution's shape is PRIMARY: I counted review round directories
per idea on disk — `integrate-parley-bidding-addon` 24, `skills-cli-install-path` 21,
`parley-design-skills` 11, tail falling to 1–2). Nothing here touches design rounds, which are flat.

1. **P1 — Mine the review loop's own artifacts before changing anything.** The decomposition the
   brief's §9 says has never been done (acted-on fraction, same-claim re-open rate,
   harness-vs-agent-forced rounds) is computable from files the protocol already mandates. Zero
   protocol bytes, nothing in the normative path. This is the only defensible *first* move (T1),
   and it is the trigger condition for P3.
2. **P2 — Replace the refuted convergence heuristic in §6.5, don't add to it.** Lines 656–659
   define converging as "total findings dropping sharply each pass" — precisely the acceptance-count
   signal VRR-Stop proves can rise while validity collapses (SECONDARY). Swap that illustrative
   trigger for two label-free signals computed from dispositions already on disk. Near-zero net
   bytes ×3 because it replaces text rather than adding (T4).
3. **P3 — Conditional: a witness requirement for cycle-opening force, not for reporting.** LE-1's
   "attempt a failing case **or** run the relevant check" (line 527, PRIMARY) is a diligence rule
   for the review as a whole; a single CRITICAL asserted without any executed witness can still
   open a full fix-up cycle today. Gate the *consequence*, never the report — and only open this
   §7 successor if P1 shows witnessless findings underperform (T3).

Explicitly **not** proposed: asymmetric reviewer context (T2 — dies by the standard that killed the
context selector, and the four-line convergence is weaker than the prompt frames it), protocol
compression or a precedence layer (T5 — its own evidence refutes the instinctive fix on strong
models; the ratified-but-unbuilt packet and subtractive-maintenance ideas already cover the
legitimate residue). Roughly a dozen brief items are "confirmation, already have it" — listed at
the end. That is a finding, not a failure of ambition.

## Proposed approach

### P1 (rank 1) — Review-loop telemetry from artifacts already on disk

**Mechanism.** A read-only analysis pass (script or disciplined manual pass — an implementation
detail, not a protocol question) over `ideas/*/review/round-*/<agent>.md`,
`ideas/*/review/consensus.md`, `ideas/*/IMPLEMENTATION.md` fix-up sections, and
`runs/*/events.jsonl`, producing per idea: findings per review round by severity; **acted-on
fraction** (findings landing in `## Agreed fixes` vs `## Deferred follow-ups` vs
`## Dismissed findings`); **same-claim re-open count** (a finding in round N+1 that cites or
textually re-litigates a finding "fixed" in round N — the observable proxy for the brief's
detect-miscorrect rate, SECONDARY: P(DM|D=1) 53–94% in an answer-extractive regime whose authors
explicitly flag non-extractive tasks as an open question); and **harness-vs-agent decomposition**
(rounds forced by validation failures / driver retries vs rounds carrying new substantive
findings — the residue the dropped 639k-telemetry finding left behind, SECONDARY).

**Measured Parley problem it touches.** The 1.6→5.1 review-round series itself, which today is a
single undifferentiated number. Every historical claim about whether more rounds helped is
compute-confounded (brief §4.5, SECONDARY) and — my addition — *cause-confounded*: the series has
never been split into "the loop found real defects" vs "the loop churned" vs "the harness forced
re-runs" (brief §9). P1 is what converts P2/P3 and several ratified-but-unbuilt ideas from
plausible-sounding structure into falsifiable ones.

**What Parley already has.** The entire input, machine-shaped: review artifacts are driver-validated
to contain a `## Findings` heading and a non-empty `## Refutation attempts` section
(`internal/runner/phase58.go:435-441`, PRIMARY); `review/consensus.md` mandates the
Agreed/Deferred/Dismissed sections with each agreed fix citing its originating finding
(`COOPERATION.md:570-577`, PRIMARY — I read them); dispositions have a recorded concurrence path
(`:545-556`, PRIMARY); run event logs exist (`parley-deck/runs/*/events.jsonl`, PRIMARY — I read
one). What does **not** exist and P1 cannot conjure retroactively: token/cost normalization —
`internal/driver/loop.go:174-175` states the runners "do not yet emit agent.usage, so this is 0 in
practice" (PRIMARY); `agent.acp.usage` (`internal/runner/acp.go:389`, PRIMARY) carries context
used/size, not cost. A compute-matched single-agent control is a *prospective experiment*, not
telemetry; I recommend deferring it explicitly rather than pretending P1 delivers it.

**Successor-idea shape.** New idea `review-loop-telemetry`, track fast, no `COOPERATION.md` delta.
Deliverable: a report with per-claim locators (every number auditable to a file:line, matching this
deck's provenance culture), plus an optional small Go addition emitting `agent.usage` — that is
observability tooling, not normative-path, so it fits the constraint; but it is a *separable*
deliverable and should not hold the report hostage.

**Cost of being wrong.** Wasted analysis effort, and misclassification risk on ambiguous artifacts
("same claim" matching is judgment, not string equality). Mitigation: publish the classification
with locators so any participant can re-read any cell — the same standard the brief's verifiers
held the researchers to.

**What evidence would show it did not work.** (a) An independent re-classification of a sample of
rounds disagrees materially with the instrument's labels — then the artifacts don't support the
decomposition and every conclusion drawn from it is void. (b) No subsequent idea, retro, or §13.3
acceptance gate cites its numbers — then it was trivia. (c) The sharp one: if acted-on fraction
≈ 1, same-claim re-opens ≈ 0, and harness-forced rounds ≈ 0, then 5.1 rounds is the loop *earning*
its cost and the entire T3/T4 concern should be closed as a non-problem. I state (c) as a live
possibility, not a hedge: the AIDev result (60.2% of agent-only PRs in the 0–30% acted-on range,
SECONDARY, correlational, and structurally unlike Parley's mandatory-disposition population) is the
only population-scale measurement we have of anything adjacent, and it does not pre-judge which
side Parley is on.

### P2 (rank 2) — Swap the §6.5 convergence trigger for computed signals

**Mechanism.** `COOPERATION.md:656-659` currently gives agents an illustrative trigger:
*converging* = "total findings dropping sharply each pass, new ones low-severity and confined to
fresh fix code". VRR-Stop (SECONDARY, all figures verified by the brief's verifier; reasoning/code/
tool-use benchmarks, not design deliberation — extrapolation) proves acceptance can rise
monotonically while true validity falls (0.87 round 2 → 0.12 round 6; 55% of instances had a
correct plan repaired into an incorrect one). The finding-count trajectory **is** that acceptance
signal. Replace the illustrative trigger with: before judging convergence, compute per pass
(i) acted-on fraction of the previous pass's findings (from the mandated consensus.md sections) and
(ii) same-claim re-open count. A pass whose findings were mostly dismissed/deferred is reviewers
stopping, not quality rising; a fix re-BLOCKed on the same claim is an observed miscorrection and
weighs *against* convergence, not for it.

**Measured Parley problem it touches.** The exploding term, directly: fix-up rounds. It sharpens
the exact judgment point (`:646-666`, PRIMARY) that decides whether a 6th, 12th, or 24th cycle
happens. `MaxFixupCycles` defaults are 1/2/5 by track (`internal/track/track.go:149-172`, PRIMARY)
yet the observed max is 24 — so escalation judgment, not the budget, is where rounds are actually
decided, which is where P2 lands.

**What Parley already has.** Trajectory-based stopping, not fixed-K (`:646-666`); the anti-signal
half-present as prose — "the same ground is re-litigated despite open rebuttals" already appears at
`:651-652` (PRIMARY) but as an uncomputed vibe; a fail-closed deterministic finding-scan on the
closing round (`:640-644`, PRIMARY), which stays untouched; grep confirms zero occurrences of
"acted on", "miscorrection", or "re-block" in `COOPERATION.md` (PRIMARY, `grep -i`, rc=1) — the
computed inputs are genuinely absent, not renamed.

**Successor-idea shape.** §7 meta-protocol-change, e.g. `meta-protocol-change-convergence-signals`.
The draft *replaces* the two illustrative-trigger sentences at `:656-659` with the computed-signal
rule. Net bytes ≈ 0 ×3 (I verified two in-repo copies, `parley-deck/COOPERATION.md` and
`internal/protocol/defaults/COOPERATION.md`, 105,382 and 104,805 bytes, differing only in
workspace/roster template zones — `diff`, 19 diff lines, PRIMARY; the third copy lives outside this
repo per project memory, SECONDARY). Pays for itself under the byte constraint by aiming at the
5.1 term: one avoided fix-up cycle per idea dwarfs any plausible byte delta.

**Cost of being wrong.** The dropping-count heuristic may carry a real signal the new inputs miss
(e.g., a genuinely clean pass after a genuine fix). Mitigation: the new rule *adds evidence to the
same judgment*, the fail-closed scan and human-escalation paths are unchanged, and the change is a
two-sentence revert if it misbehaves. Deeper risk: acted-on fraction can be gamed by reviewers
filing trivial findings — but dispositions are already adversarially concurred (`:552-556`), so
gaming must survive the same adversary.

**What evidence would show it did not work.** Stratify by track: post-adoption mean review rounds
per idea does not fall against the pre-adoption baseline (which P1 supplies), or same-claim
re-opens — now countable — do not fall; or a compliance read of closing trajectories shows agents
still judging on raw count, i.e., the text changed and the behavior didn't (the instruction-
following literature's silent-omission mode, SECONDARY, says exactly this failure is the default —
so measure it rather than assume it).

### P3 (rank 3, conditional) — Witness requirement for cycle-opening force

**Mechanism.** Today a reviewer must *as a review* show refutation attempts (LE-1, `:527`,
machine-enforced at `phase58.go:437`, both PRIMARY), but any individual finding — including a
CRITICAL asserted without an executed check, failing case, or located source — can open a full
fix-up cycle, because §15.4 deliberately declines to gate what a reviewer may report
(`:1321-1322`, PRIMARY; the no-suppression stance is repeated in §15's preamble, `:1233-1236`,
PRIMARY). The proposal keeps reporting fully ungated and gates only the **consequence**: a finding
with no executed witness is recorded, dispositioned, and concurred exactly as today, but cannot
*by itself* force a fix-up cycle; if maintained against a rebuttal it goes straight to the
DISPUTED / operator-ruling path (`:1297-1304`, PRIMARY) instead of through a repair round. Rationale
from the corpus: Refute-or-Promote's surviving residual is that the demonstrated filter is
execution, not framing ("One test killed what 80+ agents' reasoning could not", SECONDARY, no
ablation, causality disclaimed by the paper itself); VRR-Stop adds that repair under a verifier
with no discrimination destroys correct work (β 0.615–0.938, SECONDARY). A fix for an
unexecutable claim is a repair whose verifier is noise by construction.

**Measured Parley problem it touches.** Same term: fix-up round count. The brief notes this is the
only mechanism in the entire survey aimed at round count rather than quality (SECONDARY).

**What Parley already has.** The rebuttal → concurrence → DISPUTED → operator path (so the gate
adds no new machinery, only a shorter route to it); the LE-4 asymmetry ("can only fail a close
claim, never auto-pass one", `:592`, PRIMARY) as precedent that consequence-gating by evidence type
is an accepted shape in this protocol.

**Successor-idea shape.** §7 meta-protocol-change `meta-protocol-change-finding-witness-gate`,
~2 sentences in the Phase 6 / LE-1 area ×3. **Trigger condition: do not open it until P1 reports
acted-on fraction split by witness presence.** If witnessless findings are acted on at the same
rate as witnessed ones, the gate cuts real value and must not be opened — that is P1 earning its
keep as a kill switch.

**Cost of being wrong.** A true defect asserted from unverifiable-but-real experience (production
lore, a security smell that only reproduces under conditions the reviewer can't stage) gets
deferred instead of fixed, and escapes. This is the deliberate-dissent-protection concern that
motivated §15.4's non-gating, and it is why reporting stays ungated and the operator path stays
open. The cost is bounded by the follow-up-idea mechanism for deferred findings (`:573-574`).

**What evidence would show it did not work.** Pre-opening: P1's by-witness acted-on rates
contradict the premise (abort). Post-adoption: escaped-defect incidence — defects discovered after
close that a witness-gated finding had named — rises above the pre-adoption baseline (again P1's
baseline); or review rounds per idea fail to fall, meaning the gate was not the cost driver.

### The five tensions, taken as questions

- **T1 — measurement before mechanism?** Yes, with a correction to the framing: it is not
  "instrumentation vs building", because the first instrument is a *read* of artifacts the protocol
  already mandates, and its output is the trigger that decides whether the only mechanism proposals
  on the table (P2/P3) live or die. What T1 gets right: a compute-matched single-agent control and
  token normalization are not recoverable retroactively (`loop.go:174-175`, PRIMARY) — anyone
  claiming P1-style telemetry settles the compute-confound is overclaiming; that part needs a
  prospective experiment I am explicitly not proposing this round.
- **T2 — asymmetric context.** Dies by the standard, and the convergence is weaker than framed.
  Of the four lines: Anthropic's verification-minimal-context is an unmeasured vendor claim
  (weakest tier in the corpus, SECONDARY); Refute-or-Promote's cold-start reviewers come with no
  ablation and author-disclaimed causality (SECONDARY); the ICML teams-hold-experts-back result
  measures deliberation diluting a group's best member — mapping it onto *reviewer* context is the
  framing's stretch, not the paper's, and the item was dropped ALREADY-HAVE; and Parley's own
  1.6-vs-5.1 split compares two phases that differ in kind (consensus-on-prose vs
  defect-finding-in-code, different artifacts, different gates), so attributing the split to
  independence enforcement is one hypothesis among several — I verify the independence rule exists
  at design round 1 (`:324`, PRIMARY) and the phases' artifact shapes differ (PRIMARY, I read both),
  and note the confound is exactly what P1's decomposition addresses. Whatever survives of the
  intuition is already being acted on by the ratified-but-unbuilt phase-packet idea, which cuts
  *protocol* bytes from context without selecting artifacts — the only version that dodges the
  deleted selector's proof burden, reaffirmed by §7's conformal-filtering and MetaGPT drops.
  Proposal: none.
- **T3 — admission bar.** Real gap, correctly framed; my P3. But the tension's either/or
  (feature vs cost driver) is a false dichotomy — the protocol already separates *reporting* from
  *consequence*; only the latter is ungated, and only the latter is the cost driver. Sequence
  behind P1.
- **T4 — stopping on a signal that can lie.** The framing survives its own fact-check: I read
  `:656-659` and it is the acceptance-count heuristic, verbatim. My P2. One honest limit:
  VRR-Stop's regime is noisy-verifier repair on verifiable tasks; Parley reviewers are noisy
  verifiers but the workload is non-extractive, so the transfer is plausible, not proven
  (SECONDARY, scope flag from the source authors themselves).
- **T5 — rule count and collisions.** The tension contains its own refutation and says so: the
  compression recovery on the strongest model is **−1.2pp** with ρ = −0.85 (SECONDARY) — the
  instinctive fix does not work on the models this deck seats. And both results are extrapolations
  beyond their measured range (depth 20 / 500 keyword instructions vs 105,382 bytes of conditional
  prose; no source measures compliance vs protocol length — brief §9). The byte side is ratified-
  but-unbuilt (`meta-protocol-change-subtractive-maintenance`, agreed 2026-08-14, not opened —
  confirmed absent from `parley-deck/ideas/`, PRIMARY `ls`); the rules-in-context side is ratified-
  but-unbuilt (the phase-packet idea). The one genuinely uncovered residue is WIRE's static
  collision audit (stages 1–3 read only the document and decide nothing about what any agent sees,
  so no proof burden attaches): a **one-time manual collision read of `COOPERATION.md`** would fit
  naturally inside the subtractive-maintenance idea when it opens, as a report feeding it. I
  deliberately do not propose it as a separate idea — that would split authority with ratified
  work, which the prompt correctly names as worse than proposing nothing. A document-wide
  precedence order (the constitution finding) rests on first-party rationale with zero empirical
  support (SECONDARY) — not worth bytes on that tier.

### Confirmations, not opportunities (sampled)

Established by the brief and spot-verified by me where load-bearing: trajectory-based stopping
already beats fixed-K (`:646-666`, PRIMARY — two findings assumed otherwise and were wrong);
provenance caps the verdict and never resolves by count (`:1294-1295`, PRIMARY); no self-verdicts
(`:1259`, PRIMARY); recorded concession and dispute-close paths (`:552-556`, PRIMARY); cold-start
independence with anchoring named (`:324`, PRIMARY); fail-closed machine gates on close and on
review-artifact shape (`:640-644`, `phase58.go:435-441`, PRIMARY); roster authority as a generated
non-authoritative view (brief PRIMARY at `:1094`, not re-verified by me, SECONDARY-for-me).
Confidence-weighted aggregation, judge/selector layers, and MoA-style synthesis are all correctly
absent — each collides with `:1259` or `:1294`, and the brief's own sub-corpus shows the strongest
of them buying +2.4pp over a single model call on 7–8B models (SECONDARY). Declined.

## Concerns / open questions

- **"Same claim" matching is the soft joint in P1 and P2.** Agreed fixes cite originating findings
  (`:570-571`), which helps, but re-litigation under new wording needs judgment. If two independent
  classifiers can't agree, both proposals lose their signal and should be narrowed to what is
  mechanical (acted-on fraction only).
- **Harness-vs-agent decomposition may be underdetermined by events.jsonl.** I verified the logs
  exist and carry structured events (PRIMARY, one file read), not that every re-run cause is
  distinguishable. Unknown cells must be reported as unknown, not interpolated.
- **P2's compliance risk is the default, not the edge case.** The instruction literature's dominant
  failure is silent omission (SECONDARY); a two-sentence heuristic swap may change nothing agents
  do. P2 must ship with its own compliance read or it is unfalsifiable text churn — the thing this
  deck has been damaged by before.
- **The third protocol copy is outside this repo and outside the in-repo guard** (brief PRIMARY;
  my diff covered only the two in-repo copies). Any §7 successor carrying byte accounting should
  state how the third copy is counted; I could not verify it from here.
- **Review findings are not systematically provenance-tagged.** P3's "witnessless" class therefore
  needs an operational definition (no executed check / failing case / located source attached to
  *that finding*), which is a protocol-text subtlety the successor idea must get right — sloppy
  wording there reopens the suppression question §15.4 deliberately closed.

## Risks

- **Measurement becomes the project.** T1's "instrument first" is also the classic excuse to never
  decide; P1 is bounded (read-only, one idea, one report) and its kill-switch outputs are defined
  in advance so it cannot silently become permanent infrastructure.
- **Survivorship in the round-count baseline.** P1 compares against ideas that closed; abandoned
  ideas (max-24 territory) may carry different mechanics. The report must include abandoned and
  escalated ideas, not only clean closes.
- **P3's escape hatch could swallow the gate.** If every witnessless maintained finding routes to
  the operator, the operator becomes the fix-up cycle. The successor draft needs a bound (e.g.,
  operator ruling once per finding, then deferred-follow-up), or the cost just relocates.
- **Optimizing the measured term distorts it.** Once acted-on fraction and re-open count are
  convergence inputs, reviewers face pressure on those numbers (Goodhart). Mitigation is the
  existing adversarial concurrence path, but the first retro after P2 ships should look for
  disposition-inflation specifically.
- **External evidence is all SECONDARY to me** — I read none of the ~30 papers, only the brief and
  its verifier notes. Where a proposal rests on a single extrapolated source (VRR-Stop for P2,
  Refute-or-Promote for P3), I have said so inline; if the deck wants any of these load-bearing at
  ratification time, someone must fetch the source and re-verify before the §7 successor closes,
  per §15.2's own standard.
