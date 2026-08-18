---
agent: codex-1
idea: mas-research-mining
round: 1
date: 2026-08-15
---

## Summary

I recommend **two successor ideas, only one of them immediate**:

1. Extend the existing read-only retrospective scanner into a review-loop measurement and frozen-replay harness. Measure the mechanism before changing the normative process.
2. Only if that baseline finds a material same-claim miscorrection rate, replace the current finding-count convergence example with a scoped repair-damage tripwire. This should be a small, byte-neutral-or-negative §7 change, not VRR-Stop's Bayesian controller.

I do **not** recommend a cold-context reviewer policy, a document-wide precedence order, or a finding-admission gate in this round. The first would decide what a reviewer sees without a no-drop result; the second has no measured Parley collision corpus and overlaps ratified subtractive work; the third risks suppressing genuine dissent before we know whether unsupported findings actually cause the 1.6→5.1 review expansion. The literature's strongest negative lesson is to add less process, not to transplant every mechanism it describes (SECONDARY: `reference/research-brief.md:263-280`).

The target is the measured review/fix-up term: review rounds rose from 1.6 to 5.1, with a maximum of 24, while design rounds remained nearly flat (PRIMARY repository record: `parley-deck/ideas/protocol-read-cost-regression/FINAL.md:18-27`; the current task repeats the constraint at `00-prompt.md:142-147`). No proposal below adds design-round ceremony.

All claims about external research below are **SECONDARY**, dependent on the adversarially verified research brief; I did not read the cited papers. Repository claims carry PRIMARY file locators or the exact check run.

## Proposed approach

### Rank 1 — review-loop observability plus frozen replay

**Mechanism.** Open an ordinary tooling successor, tentatively `review-loop-observability-and-replay`, which extends the existing optional, read-only `parley retro scan` path. It should produce one row per review finding and transition, with raw numerators and denominators:

- finding source, severity, round, disposition, and whether it became an agreed fix;
- whether the resulting fix was later accepted, marked `NOT-FIXED`, re-blocked on the same scoped claim, superseded by a different claim, or followed by a regression;
- whether each extra cycle came from an agent finding/fix, a strict-clean pass, an agent/watchdog failure, an artifact-validation retry, or another harness-forced transition;
- elapsed time and usage when actually emitted, never fabricated; and
- per-idea totals for acted-on fraction, fix-up cycles, same-claim re-blocks, harness-forced cycles, and unique evidence-bearing findings per unit of time/usage.

The classifier must not pretend open-ended code review has an oracle. `NOT-FIXED` is a candidate signal, not automatically “DM”; a small, predeclared hard-case sample should be independently labelled from the frozen pre-fix and post-fix commits, the scoped finding, and available checks. Report inter-rater disagreement rather than resolving it by vote.

The same successor should run a **non-canonical frozen replay** on a small set of completed ideas with executable acceptance checks. Hold the implementation commit and reviewer model fixed, predeclare at least three repetitions per arm, and compare:

1. today's review prompt: `FINAL.md` + `IMPLEMENTATION.md` + all prior review rounds;
2. a cold-start review prompt: `FINAL.md` + `IMPLEMENTATION.md` + frozen diff/check evidence, with no prior review prose or dispositions; and
3. one reviewer given the same aggregate time/token budget as the multi-review arm.

Replay outputs must never become signoffs, dispositions, or consent. This is how T2 can be tested without reviving the deleted selector in the normative path. The runner currently embeds all prior review-round artifacts after `FINAL.md` and `IMPLEMENTATION.md` (PRIMARY: `internal/runner/phase58.go:276-306`); design cross-review similarly concatenates every earlier artifact (PRIMARY: `internal/runner/runner.go:936-965`).

**Measured Parley problem touched.** The direct target is review/fix-up expansion, not the flat design series. Secondary external evidence supports the *measurement discipline*, not a winner: compute-normalized multi-agent gains often disappear, but the tested workloads are not Parley's (`reference/research-brief.md:208-213`); real-PR review work supplies an acted-on-fraction instrument but not a correctness label (`reference/research-brief.md:134-138`); the four-way context asymmetry is explicitly weak in every individual line (`reference/research-brief.md:193-198`). The experiment therefore measures Parley rather than assuming transfer.

**What Parley already has.** This is an extension, not a new telemetry stack:

- canonical review findings, dispositions, agreed fixes, and per-cycle audit artifacts already exist (PRIMARY: `parley-deck/COOPERATION.md:535-607`);
- `internal/retro` already counts design rounds, review rounds, fix-up cycles, `NOT-FIXED`, dismissed findings, escalations, blocks, and run failures (PRIMARY: `internal/retro/retro.go:19-35,83-143`);
- the driver emits elapsed loop-budget data (PRIMARY: `internal/driver/loop.go:154-171`), while its own comment states that runners do not yet emit `agent.usage`, so cost is zero in practice (PRIMARY: `internal/driver/loop.go:174-190`); and
- the driver preserves per-cycle review consensus and runs the post-fix checks before opening the next review round (PRIMARY: `internal/driver/impl.go:324-355`).

I ran this repository check:

```text
rg -n -i "acted[-_ ]on|signal ratio|detect[-_ ]?correct|detect[-_ ]?miscorrect|miscorrect|pass.?2|run[-_ ]to[-_ ]run reproduc|single[-_ ]agent control" internal cmd parley-deck/COOPERATION.md
```

Its only hit was an unrelated `// Pass 2` comment in `internal/retro/retro.go:190`; the requested measures are not named in the scanned code/protocol. The positive existing fields above are why this should extend `internal/retro`, not create a parallel authority.

**Concrete successor-idea shape.** Standard track, tooling only, no protocol edit, no dependency, no service, and no tool in the normative path. Scope it to `internal/retro`, existing structured run events, fixtures, and a replay command that writes outside canonical idea artifacts. Acceptance criteria should require:

- raw counts plus denominators, with every derived label traceable to canonical paths/commits;
- an explicit `unknown` bucket for unclassifiable DC/DM and finding-to-fix links;
- separate agent-error, harness-retry, strict-clean, and actual fix-up counts;
- replay arms frozen before execution and compute matched using observed, not guessed, usage;
- repeated-run variance reported alongside arm means;
- no mutation of closed ideas and no replay output consumed by consensus; and
- a short result stating which, if any, successor mechanism is now justified.

**Cost of being wrong.** The scanner can create false precision from prose artifacts, the manual audit can be expensive, and acted-on fraction can reward bad fixes or punish correct rejected findings. A replay can overfit a small coreset and still lack a true code-quality oracle. Missing `agent.usage` means token normalization may initially be impossible; the tool must say so rather than substitute wall clock silently.

**Evidence that would show it did not work.** Treat the successor as failed if any of the following occurs:

- finding→disposition→fix linkage has too much `unknown` data or independent auditors cannot reproduce DC/DM labels;
- harness-forced versus agent-caused cycles cannot be reconstructed from canonical artifacts/events;
- repeated-run variance is as large as, or larger than, the cold/full-context or single/multi-agent differences;
- cold-start reviewers lose later-validated unique findings or create more retracted/dismissed findings at matched spend;
- unsupported or non-executed findings do not predict extra fix-up cycles, undermining T3's proposed causal story; or
- the added measurement cost exceeds the review time it can plausibly save.

### Rank 2 — a scoped repair-damage tripwire, contingent on Rank 1

**Mechanism.** Open `meta-protocol-change-review-damage-tripwire` only if Rank 1 finds repeated same-claim re-blocks at a material rate. Make two narrow changes to the existing stopping paragraph:

1. Finding count and severity trajectory are **workload/churn signals only**, not evidence that artifact validity is rising.
2. When a fix is re-blocked as `NOT-FIXED` on the same scoped claim, pause that claim's repair thread, preserve the pre-fix commit as the incumbent comparison, and require an evidence-bearing pre/post comparison or operator ruling before another repair. Unrelated findings may continue.

This is deliberately smaller than VRR-Stop. It adds no posterior, learned noise parameters, fixed-K stop, majority rule, automatic revert, new service, or new dependency. A review consensus can use its existing `blocked` state to force escalation; the driver already escalates when review consensus is blocked (PRIMARY: `internal/driver/impl.go:190-196`). Before that branch, the driver has no semantic same-claim input: with outstanding agreed fixes it checks the cap, clean tree, auto-implementation flag, applies the fix, runs existing checks, and opens another round (PRIMARY characterization of the decision inputs: `internal/driver/impl.go:290-350`).

**Measured Parley problem touched.** The tripwire attacks repeated repair, the term that reaches 5.1 average and 24 maximum rounds, and it corrects a currently named heuristic: the protocol calls a sharply falling finding count “converging” (PRIMARY: `parley-deck/COOPERATION.md:646-659`). Secondary evidence shows why that signal can be anti-correlated with validity in noisy verify-repair tasks, but also explicitly warns that the source did not study open-ended design/code synthesis (`reference/research-brief.md:103-107,122-132`). That is why the trigger is contingent on Parley's own baseline rather than adopted directly.

**What Parley already has.** Parley is already ahead of the source's fixed-K baseline: budgets escalate rather than certify completion, deterministic checks can veto a close, and minority disagreement can remain live (PRIMARY: `parley-deck/COOPERATION.md:648-669,1270-1309`). `internal/retro` already recognizes `NOT-FIXED` as a re-review signal (PRIMARY: `internal/retro/retro.go:38-52,99-114`). The ratified objection-ledger contract already specifies owner-namespaced IDs, exact propositions, lifecycle, hashes, and forced `DISPUTED` (`parley-deck/ideas/protocol-read-cost-regression/FINAL.md:60-83`); a tripwire must reuse that authority if it later becomes machine-enforced, not create a second ID system. Git commit `41e6cd6` deleted the dormant frontier implementation after its authored-ledger validation was still fail-open; this proposal must not turn that deleted machinery back on.

**Concrete successor-idea shape.** Deliberation track because it changes `COOPERATION.md`. Replace the illustrative finding-count sentence rather than append a new section. Target **net ≤ 0 protocol bytes per copy**, hence net ≤ 0 across the three lockstep copies required by this idea (`00-prompt.md:142-144`). Use existing review finding locators and `blocked`/operator escalation manually; defer any automation until the ratified ledger has its separate failure analysis and validator. Acceptance criteria:

- aggregate finding decline is never described as validity evidence;
- same-claim `NOT-FIXED` pauses only that claim thread;
- a pre/post comparison records which check/evidence discriminates the incumbent and candidate, or records that none exists and escalates;
- no automatic rollback and no suppression of other findings;
- no new artifact or tracking vocabulary duplicating the objection ledger; and
- the prior wording removed offsets all new normative bytes.

**Cost of being wrong.** “Same claim” can be subjective before the ledger exists. A tripwire may escalate correct iterative refinement, increase operator burden, or preserve a faulty incumbent. Checks can also be incomplete, so “pre-fix passes” is not proof that the incumbent is better.

**Evidence that would show it did not work.** Reject or revert the change if, over a preregistered cohort:

- same-claim third-or-later repair attempts and total fix-up rounds do not fall;
- most tripwires are later adjudicated as different claims rather than miscorrections;
- operator escalations rise without a reduction in accepted bad fixes or later regressions;
- deterministic checks and later review consistently favour the repaired candidate despite the tripwire; or
- genuine CRITICAL/MAJOR fixes are delayed enough that time-to-correctness worsens at equal review quality.

## Concerns / open questions

### T1 — measurement before mechanism

Yes for the first move, but not as permanent deferral. The existing scanner already has most coarse inputs, so a bounded extension is cheaper and more falsifiable than changing reviewer context or admission rules. The result should contain a go/no-go decision for Rank 2 and T3, not merely a dashboard. The external literature does not measure Parley's workload and explicitly leaves all six critical Parley measurements open (SECONDARY: `reference/research-brief.md:341-347`).

### T2 — asymmetric context

The convergence is a hypothesis worth testing, not a normative result. The vendor evidence is the weakest surviving tier, the artifact-reference primitive is unablated, and every pointer still creates a decision about what the reviewer sees (SECONDARY: `reference/research-brief.md:193-198`). The repository's deleted frontier work demonstrates the local failure mode: the first derived ledger was fail-open, then authored-ledger presence was accepted without validation, and the feature was hard-disabled before deletion (PRIMARY: `git show 41e6cd6 -- internal/runner/frontier.go`; see also the ratified boundary and fallback at `protocol-read-cost-regression/FINAL.md:77-100`). A frozen replay does not create consent and therefore survives the no-drop standard; a production selector does not yet.

### T3 — finding admission

The framing is a false binary. Keep the reporting channel ungated: Phase 6 explicitly forbids suppressing findings and gives the owner-controlled close paths (PRIMARY: `parley-deck/COOPERATION.md:535-556`). The possible future lever is **promotion from reported finding to agreed code-changing fix**, not permission to speak. Do not open that protocol idea unless Rank 1 shows that RECALL-only or non-executed findings materially predict extra cycles. Refute-or-Promote's executed-test residual is suggestive, but its paper disclaims causal isolation and its flagship unanimity failed until a test was run (SECONDARY: `reference/research-brief.md:109-113`).

If the baseline does implicate unsupported promotion, the later successor should replace the existing `Agreed fixes` instruction with a byte-neutral requirement for a finding-linked executed failure, stable source witness, or explicit operator risk ruling. It must not alter the no-suppression paragraph. That is a conditional direction, not a third recommendation today.

### T4 — stopping on a signal that can lie

The fixed-K premise was false, but the corrected tension is real. “Fewer findings” is useful for workload forecasting; it is not evidence that the latest artifact is better. Rank 2 removes only that inference and adds an observable local damage signal. It does not claim that the secondary paper's estimated probabilities transfer to Parley.

### T5 — rule count and collisions

No new successor is justified. Instruction-density results are extrapolated far outside their measured range, and compilation was neutral or worse on the strongest tested model (SECONDARY: `reference/research-brief.md:83-87`). WIRE establishes that silent within-policy collision is possible, but not that a specific Parley collision caused review expansion; its pipeline has a 0.55% candidate yield and no replication (SECONDARY: `reference/research-brief.md:80-81`). The document-wide precedence example is first-party rationale with zero empirical support (`reference/research-brief.md:95-99`).

The correct action is to complete already-ratified work rather than split authority: the phase-scoped packet and subtractive-maintenance/rule-inventory ideas are already named as open (`00-prompt.md:116-129`). I would not add a global priority order until that inventory produces concrete witnesses showing which rules co-govern one decision and conflict. The repository check `wc -c -l parley-deck/COOPERATION.md; rg -c '^## ' parley-deck/COOPERATION.md` returned `105382` bytes, `1372` lines, and `18` top-level sections; size is measured, semantic collision is not.

### Confirmations, not opportunities

Do not reverse or re-propose the mechanisms for which Parley is already ahead: provenance that caps verdicts, no self-verdicts, never resolving by participant count, `DISPUTED`, round-1 independence, owner-controlled concession, minority escalation, append-only artifacts, acknowledgment receipts, and deterministic close vetoes (PRIMARY examples: `parley-deck/COOPERATION.md:324,552-556,592,1259,1270-1309,1356-1360`; SECONDARY comparison summary: `reference/research-brief.md:284-302`). Section 5's negative evidence strengthens the burden of proof for adding ceremony; it does not establish that these protections should be deleted because none of the sources measures multi-turn repository review (`reference/research-brief.md:263-280`).

## Risks

- **Evidence transfer.** All source-side results are SECONDARY here, and almost every benchmark has a correctness label or workload unlike Parley's. Neither proposal treats a paper's headline number as a Parley effect size.
- **Metric gaming.** Acted-on fraction, low finding count, short duration, and low token use can all improve while correctness worsens. Rank 1 must publish them separately and never collapse them into one optimization score.
- **Observer and replay effects.** Frozen replay is not a live idea: reviewers know the artifact is fixed, and selected completed ideas may overrepresent difficult cases. Results must be labelled accordingly.
- **False tripwires.** Without validated claim identities, “same claim” is a reviewable judgment. Rank 2 therefore stays manual, local to one finding thread, reversible, and contingent.
- **Protocol bloat.** Any normative sentence is multiplied across three copies. Rank 2 pays for itself by replacing and shortening the existing illustrative convergence text; Rank 1 changes no protocol bytes.
- **Building the wrong layer.** The Go runner does not read `COOPERATION.md` to construct a phase-scoped protocol packet (PRIMARY task record: `00-prompt.md:121-123`). Reviewer-prompt experiments belong in runner/tooling; protocol-size reduction belongs in the already-open instruction-layer idea.
