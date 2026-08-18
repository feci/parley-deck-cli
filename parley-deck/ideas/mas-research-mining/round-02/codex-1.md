---
agent: codex-1
idea: mas-research-mining
round: 2
date: 2026-08-15
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

I change five positions.

1. **Frozen replay moves out of the first successor.** @claude-1, @hermes-1 and @kimi-1 are right
   that the first deliverable should be a bounded read of existing artifacts. Replay is the only
   proposal that can test reduced reviewer context without changing the normative path, but it is a
   prospective experiment with a different cost and validity burden. It should be a separately
   approved second successor, triggered by the baseline's feasibility result.
2. **I no longer support waiting before correcting the §6 stopping language.** The current example
   calls a falling finding count “converging” (PRIMARY: `parley-deck/COOPERATION.md:656-659`). Even
   without deciding whether Parley has a material miscorrection rate, a workload count is not
   evidence of artifact validity. I would replace that example in a small §7 successor, but I would
   not replace it with acted-on fraction, which is also not a correctness measure.
3. **“Gate reporting” and “gate cycle-opening force” are genuinely different.** I overstated the
   suppression objection in round 1. A force gate can preserve the report, disposition, dispute and
   operator paths. It can still suppress a true objection's practical effect by delaying repair, so
   I do not support it before measurement.
4. **DC/DM is not a historical metric we can honestly claim to compute.** A later `NOT-FIXED` or
   same-locator re-block is an observable failed-repair signal, not an oracle-backed label that the
   prior artifact was correct and the repair made it wrong. The first successor should report the
   observable transition and `unknown`, not rename it DM.
5. **There is a concrete “do less” action:** delete the misleading count-as-convergence example
   rather than add a re-block clause or a new computed proxy. Separately, use replay—if justified—to
   test whether `standard` should have fewer reviewers. Do not add T3, T5, or a written cold-start
   rule now.

## Is our convergence independent evidence or a shared prior?

It is mostly a **shared prior**, not four independent confirmations. Model-family diversity does not
make the evidence independent when every participant received the same 87 KB research brief, the
same five-tension framing, the same repository, and constraints that strongly favor reversible,
non-normative work. The brief explicitly foregrounded the six missing measurements and “measurement
before mechanism”; all four answers then selected from that menu. The agreement is useful evidence
that a read-only baseline is a mutually acceptable low-regret decision. It is weak evidence that the
baseline's causal hypotheses are true.

There is some independent residue: participants inspected different repository paths and surfaced
different limitations, and @claude-1 found the decisive runner correction that review round 1 is
already cold-start. But after that correction, the supposed T2 convergence loses the only leg
measured on Parley's workload. The remaining agreement should therefore be recorded under §15.6 as
correlated judgment. It must not be cited as four votes proving that instrumentation—or any T2/T3/T4
mechanism—will reduce review rounds.

For our common recommendation to be wrong, at least one of these would have to hold: existing
artifacts are too ambiguous to yield even auditable descriptive counts; the cost of building and
maintaining the scanner exceeds any decision value; or a cheap deletion/reviewer reduction has such
an obvious favorable risk profile that measuring first is needless delay. The first successor must
be allowed to return that null result and stop.

## Responses to others

### @claude-1

I concede C1. The facilitator's PRIMARY correction at `internal/runner/phase58.go:283` establishes
that review round 1 already receives no earlier review files. T2's Parley-side premise is false, and
the 1.6-versus-5.1 split cannot be attributed to round-1 cold start. I also accept C5's retraction,
with one refinement: the stacking evidence does not establish that tighter text harms Parley's
models either. The compliance benefit should be recorded as **unverified**, not inverted into a
harm claim. The ratified subtractive work should proceed only on the repository-measured read-cost
and latency case unless new Parley evidence appears.

I also concede the ledger point in the narrow form: stable owner-namespaced claim identity is a
prerequisite for a *complete* same-claim rate. I disagree that the baseline should wait for the
ledger. Counter-proposal: emit a high-precision lower bound for exact originating locators, keep all
ambiguous relitigation in `unknown`, and let the result quantify whether the ledger is necessary.

On C4, reporting and forcing are distinct powers. I would preserve the idea as a conditional
successor direction, but not sign it now: a RECALL-only defect may be real, and current artifacts do
not expose a reliable per-finding witness field. On C1's documentation fix, my counter-proposal is
to defer writing cold-start as a normative invariant until replay gives a reason to preserve it;
documenting current runner behavior as policy would itself decide T2 without evidence.

### @codex-1

My round-1 replay proposal bundled a static retrospective, manual correctness annotation, usage
instrumentation, and a prospective three-arm experiment into one successor. That was too broad and
would make a cheap baseline hostage to missing `agent.usage` and an oracle problem. I retract that
scope. I retain replay only as a separate, preregistered experiment whose outputs are non-canonical.

My round-1 “change nothing until a material same-claim miscorrection rate” condition was also too
strict. We do not need a Parley effect estimate to stop calling a falling workload count evidence of
convergence. The replacement must be validity-neutral; acted-on fraction and a guessed DC/DM label
must not become new stopping scores.

### @hermes-1

The facilitator correction disposes of the sync concern. The two in-repo protocol copies differ
only in normalized project-specific zones, and `go test ./internal/protocol/...` passes. I accept
the useful residue: byte accounting must compare shared rule text, not whole-file bytes; “net bytes
×3” is otherwise ambiguous.

I agree with P1's direction but not its claim that DC/DM can be computed by matching prose. A
re-block may mean an incomplete fix, a wrong fix, a changed scope, a newly exposed defect, or a
reviewer mistake. Counter-proposal: call the machine-visible categories `explicit-same-locator
reopen`, `NOT-FIXED`, `new-locator`, and `unknown`; reserve DC/DM for a manually labelled sample with
an external correctness witness, and report annotator disagreement.

I do not support adding P2's re-block clause. The preceding §6 paragraph already treats repeated
ground as churning (PRIMARY: `parley-deck/COOPERATION.md:648-654`). Adding another example preserves
the more fundamental error that proxy trajectories establish validity. My counter-proposal is
subtractive: replace the two illustrative sentences with a shorter statement that finding
trajectories estimate remaining review work but do not establish validity. I agree with your own
risk assessment that P3, the precedence order, is the weakest proposal; I would drop it entirely.

### @kimi-1

I concede the core of P1 and your explicit separation of retrospective telemetry from prospective
compute-matched controls. I also adopt your `unknown` discipline and the requirement that a manual
classification publish locators and disagreement instead of manufacturing consensus.

I disagree with P2's proposed replacement. “Acted on” has at least two meanings in current
artifacts: promoted into `## Agreed fixes`, and recorded as applied in an `IMPLEMENTATION.md` fix-up
section. Neither says the resulting artifact is correct. A high fraction can reward bad repairs; a
low fraction can mean noisy reviews or appropriately resisted findings. Counter-proposal: report
both raw ratios descriptively, never make either a continue/stop input, and replace the current
example only with a warning that count trajectories are workload signals.

I now agree with P3's conceptual distinction between reporting and force, but its trigger is not
fully available from P1. Review findings have no mandatory per-finding provenance or witness field.
“No attached executed witness” can be manually observed; “witnessless” cannot safely be inferred
from absent free-form prose. Counter-proposal: a capped manual audit may code `attached witness:
yes/no/unknown`; open no force-gate successor unless that audit shows a material association with
extra cycles and separately checks for later-valid defects that the gate would have delayed.

That schema claim is PRIMARY: the required review shape contains summary, refutation attempts,
severity headings and open questions at `parley-deck/COOPERATION.md:515-527`, not a per-finding
witness field. I also ran an `rg` check for per-finding witness/provenance/evidence/source fields over
`COOPERATION.md`, `internal/runner`, `internal/driver`, and `internal/retro`; it found no such schema
field (the only code hit was an unrelated `internal/retro/retro.go:247` comment).

## New concerns / questions

### What is computable from disk today

The first successor must distinguish extraction from interpretation. The following is the contract I
would use:

| Candidate | Artifact and field | Honest status today |
| --- | --- | --- |
| Findings per round and severity | `review/round-NN/<agent>.md`; round path plus `[CRITICAL]`, `[MAJOR]`, `[MINOR]`, `[NIT]` headings | Mechanically computable; malformed headings go to `unknown`. |
| Review rounds and fix-up cycles | review directory names; `IMPLEMENTATION.md` `## Fix-up cycle N`; existing `internal/retro` counts | Mechanically computable, with discrepancies reported rather than silently reconciled. |
| Promoted-to-fix fraction | `review/consensus.md` `## Agreed fixes` citations divided by review findings | Computable only for explicitly linked findings and the consensus versions retained at HEAD or in git history. Unlinked items stay unknown. It is not correctness. |
| Recorded-applied fraction | `IMPLEMENTATION.md` `### Fixes applied` linked back through an agreed fix to an originating finding | A high-precision partial metric; missing links are unknown. This is closer to “acted on” than agreement alone, but still not correctness. |
| Exact same-claim reopen lower bound | An agreed fix's required originating locator followed by a later review that cites that same locator | Computable as a lower bound. Reworded or split claims require judgment; a complete rate requires the ratified ledger or manual coding. |
| DC versus DM | No canonical correctness field exists for pre-fix and post-fix artifacts | **Not computable** from design/review prose. `NOT-FIXED` and a same-locator reopen are observations, not DM labels. Executable checks may witness a scoped case but are not a general oracle. |
| Harness-forced versus agent-caused cycles | `runs/*/events.jsonl`, validation/check events, driver failure/retry events, and corresponding review/fix-up artifacts | Only explicitly recorded causes are computable. A schema inventory comes first; every unmatched transition is unknown. |
| Rounds versus protocol size | idea `00-prompt.md` `created:` plus git history for the contemporaneous shared `COOPERATION.md` rule text; round counts above | Descriptively computable. Date-to-commit ambiguity must be exposed, and correlation is confounded by track, task difficulty, model, and runner version. It cannot revive a causal compliance claim. |
| Attached witness presence | Free-form text in each finding; no mandatory per-finding witness/provenance field | Manual `yes/no/unknown` sample only. Absence of text is not proof that no witness existed. |
| Inter-rater disagreement | Not an existing artifact field | An output of a new, capped annotation exercise. Report the disagreement matrix/raw labels; do not resolve it by participant count. |
| Frozen replay or single-agent control | No historical artifact contains the unrun counterfactual; usable token cost is currently absent | Prospective experiment only, not part of the disk scan. |

This eliminates two false promises. First, acted-on fraction is label-free only as a behavior count,
not as a quality measure. Second, a DC/DM classifier without a correctness witness is an oracle in
disguise.

### D2 sequencing

I choose **replace**, but not Kimi's replacement metric and not Hermes's augmentation. A separate
§7 successor should delete the two illustrative count-trigger sentences and substitute, at net
negative shared-rule bytes, approximately:

> Finding trajectories estimate remaining review work; they do not establish artifact validity.
> Continue, pause, or escalate using evidence against the acceptance criteria, unresolved claims,
> and the existing budget/escalation rules.

The exact text belongs to that successor. The already-present “same ground is re-litigated” clause
remains; no new re-block clause is needed. This correction does not wait on the baseline because it
withdraws an inference rather than installing a new mechanism.

### D3 consequence gating

There is a real difference between preserving a finding and preserving its power to cause a repair.
Calling every force gate “suppression” obscures that design space. But preserving text alone is not
enough: if a true CRITICAL is routed into an operator queue that is slower or routinely deferred, its
effect has been suppressed operationally. Therefore the future admission bar, if any, must be framed
as a measured trade: fewer unsupported repair cycles versus later-valid defects delayed or escaped.
Nothing in the current disk schema settles that trade.

### D4 compliance and subtraction

I accept @claude-1's retraction. The +11.0/+3.3/−1.2 percentage-point ladder and its correlation are
SECONDARY here, and both stacking studies are far outside Parley's prose regime. They support
neither a confident benefit nor a confident harm claim for frontier participants. The
`meta-protocol-change-subtractive-maintenance` case should be read-cost/latency and maintainability
only. Even a rounds-versus-bytes correlation from the baseline would be hypothesis-generating, not
causal compliance evidence.

### D6 doing less

The round-1 omission was a blind spot. The corpus's negative results do not identify which Parley
protection to delete, but they raise the burden for every added mechanism. The defensible immediate
subtraction is the invalid count-as-convergence example. The defensible experimental subtraction is
the one-reviewer replay arm, explicitly aimed at whether `standard` can reduce its reviewer count.
Until that result exists, I would remove no dissent, provenance, deterministic veto, signoff, or
minority-escalation protection.

## Adversarial alternative

The strongest alternative is **do no measurement successor at all**: change `standard` from two
reviewers/two fix-up cycles to one reviewer/one fix-up cycle now, enforce escalation strictly, and
spend the saved review budget on deterministic acceptance checks. Its best support is the brief's
twelve negative results about elaborate multi-agent process (SECONDARY; I did not read the cited
sources), the absence of a correctness oracle for our proposed telemetry, and the danger that a
dashboard becomes permanent process without changing decisions. It directly attacks the measured
cost rather than studying it.

I reject immediate adoption because none of those sources measures Parley's repository review, and
the second reviewer may be the only source of a later-validated defect. The observation that would
change my recommendation is a preregistered replay or auditable historical sample showing that the
second reviewer contributes no material later-validated unique findings at matched spend, while one
reviewer does not increase escaped defects or acceptance-check failures. Conversely, if existing
artifacts cannot support the bounded baseline table above with tolerable unknown rates, I would
adopt the alternative's “no measurement infrastructure” half and stop rather than expand the
instrument.

## Current proposal

I would sign exactly the following package:

1. **First successor: `review-loop-baseline`.** `standard` tooling idea; read-only over canonical
   review/consensus/implementation artifacts, run events and git history; no protocol edit,
   dependency, service, normative-path tool, replay, or signoff-producing output. It emits raw counts
   and denominators for mechanically available rows in the table above, uses `unknown` for every
   broken or ambiguous link, and includes a capped double-coded sample for same-claim and attached-
   witness judgments with raw inter-rater disagreement. It must state a go/no-go for a replay and for
   any T3 successor. Failure is an unknown rate or annotation disagreement too high to support the
   intended decision; in that case, stop rather than add fields to the protocol.
2. **Separate subtractive §7 successor: convergence language.** Replace—not augment—the §6
   illustrative count language with the validity-neutral wording above. Shared-rule-text bytes must
   be net negative; project-specific generated/header/roster bytes are excluded from accounting.
   Add no acted-on threshold, DC/DM controller, same-claim automation, or new artifact.
3. **Second experimental successor, conditional: `review-loop-frozen-replay`.** Open only if the
   baseline identifies an auditable cohort and the run can predeclare a real spend-matching method.
   Compare later-round full prior-review context, cold later rounds, and one reviewer at matched
   aggregate spend, with repeated runs and executable acceptance evidence where available. Round 1
   remains cold in every arm because that is already runner behavior. Replay outputs never become
   signoffs, findings, dispositions or consent. The explicit decision target is whether T2 survives
   and whether `standard` can use fewer reviewers; high variance, lost later-validated findings, or
   absent spend telemetry kills the experiment's decision claim.
4. **No T3 or T5 successor now.** Preserve ungated reporting. Reconsider a cycle-opening force gate
   only if the baseline/manual audit shows that no-attached-witness findings materially predict
   extra cycles and quantifies the true-defect delay risk. Drop document-wide precedence absent a
   concrete collision corpus. Do not add a normative cold-start rule before replay.
5. **Retract the compliance benefit for protocol shrinking.** Continue already-ratified
   subtractive maintenance on read-cost/latency/maintainability grounds only. Treat both benefit and
   harm to compliance as unverified until Parley-specific evidence exists.

This package adds zero protocol bytes in the first successor, makes the only proposed protocol
change subtractive, preserves T2 as a non-canonical experiment rather than policy, and leaves a
clear stopping point if the artifacts cannot answer the questions.
