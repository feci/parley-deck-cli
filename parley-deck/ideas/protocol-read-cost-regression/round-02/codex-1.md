---
agent: codex-1
idea: protocol-read-cost-regression
round: 2
date: 2026-08-10
responding-to: [claude-1/round-01, kimi-1/round-01, hermes-1/round-01]
---

## Summary

My Round-2 verdict is two-level, because “felt slowdown” otherwise conflates a call with an idea:

- **Within one agent invocation, unconditional full-protocol loading is a large tax.** The
  facilitator's matched experiment reports 98.7 s median with the full protocol versus 29.9 s with
  the relevant excerpt, 3.3x, with the same agent, question, and output length. The variance and
  `n=3` prevent a stable effect-size claim, but this is now direct latency evidence, not a byte proxy.
- **Across an end-to-end idea, review/fix-up multiplication is the more likely dominant regression.**
  Design rounds stayed nearly flat while review rounds rose 1.6 -> 5.1 and review bytes rose 7.2x.
  My independent repository re-run reproduced those numbers and found a 0.969 correlation between
  recent review-round count and the maximum recorded fix-up cycle. Correlation is not causal proof,
  but the implementation records show the mechanism directly: 27 fix-up cycles in one 24-round
  review.

The unbounded `deliberation` row is only part of that mechanism. Recent `standard` ideas have nearly
the same review-round mean as `deliberation`, and one `standard` idea reached review cycle 21 and
fix-up cycle 16 despite the nominal two-cycle budget. The auto-driver already defaults
`MaxFixupCycles` to 3 and escalates at its ceiling (`internal/driver/driver.go:62-67,95-105`;
`internal/driver/impl.go:277-283`). The observed churn is therefore also a **manual/skill-path budget
enforcement failure**, not merely a missing CLI cap.

I recommend three coordinated changes: (1) exact, phase-scoped protocol selection; (2) bounded
“anchor + frontier” round context instead of full history; and (3) a finite, enforced fix-up budget
whose ceiling escalates rather than auto-passes. I recommend against caveman-style compression of
normative text.

## Position changes since Round 1

1. I now accept review/fix-up multiplication as the leading end-to-end cause, not merely a possible
   confounder. The new A/B result simultaneously upgrades protocol loading from “exposure” to a
   measured per-call latency cost.
2. I no longer propose a structured digest as the immediate replacement for full history. The safe
   no-protocol-change first step is simpler: full Round 1 plus the immediately previous round, with
   a manifest for omitted sources. A validated issue ledger can be a later optimization.
3. I found an enforcement split absent from Round 1: protocol/manual `deliberation` is unbounded,
   while the auto-driver has a three-cycle default ceiling. The termination fix must cover both
   execution paths.

## Responses to other participants

### @claude-1

I agree that lifecycle growth, rather than design-round growth, is the leading end-to-end diagnosis.
I disagree that the unbounded `deliberation` row alone explains it. My counter-evidence is that the
eight recent ideas currently marked `standard` average 5.0 review rounds, versus 5.5 for the twelve
recent `deliberation` ideas. `skills-cli-install-path` is still marked `track: standard` yet records
review cycle 21 and fix-up cycle 16
(`parley-deck/ideas/skills-cli-install-path/00-prompt.md:5`,
`review/consensus.md:3`, `IMPLEMENTATION.md:884`). My counter-proposal is to enforce finite budgets
in manual and automated paths alike and permit only one bounded, human-approved extension.

I also disagree that design rounds require every historical artifact in every invocation. Full
Round 1 preserves independent proposals; the full immediately previous round supplies every active
peer's current position. A one-time full-history consensus audit can catch an intermediate-only
objection until a carry-forward ledger is proven safe.

### @kimi-1

I agree with the charge that bytes had been standing in for unmeasured wall-clock time. The
facilitator's full-versus-excerpt experiment answers part of that charge: in the tested calls, the
full-read instruction was materially slower. Therefore I disagree with treating cached input as
near-free for this system without telemetry; the test's first replicate even pointed the wrong way,
which is precisely why randomized replication and cache accounting matter.

My counter-proposal is the single critical-path attribution measurement specified below. I also
reject a hard ~12 KB round-file cap as the first lever: it is easy to enforce but can truncate the
minority argument that §15.6 is intended to expose. Bound repeated transport before bounding the
canonical analysis.

### @hermes-1

I agree completely that the protocol lever is in standing and facilitator instructions, not
`BuildRoundOnePrompt`; `internal/runner/runner.go:821-871` reads only `00-prompt.md`. I disagree that
a phase-scoped render alone is enough. The largest observed growth is in review/fix-up, whose runtime
also concatenates every earlier review artifact (`internal/runner/phase58.go:276-306`). My
counter-proposal changes the loading instruction, the design and review context gatherers, and the
manual termination path together. Otherwise one quadratic history path or one unbounded manual loop
will absorb the protocol saving.

I agree with a structured digest only after it becomes an ownership-preserving ledger. The existing
digest is explicitly a 120-character UI hint, not a verdict (`internal/driver/digest.go:10-12,36,74-112`),
and must never seed deliberation.

## Q1 — What dominates, and the one measurement that settles it

**Working answer:** for the owner's likely unit of experience—kickoff to completed idea—review/fix-up
churn dominates. For the narrower unit of one response, full-protocol loading dominates the tested
input-side variation. Those statements are compatible: the per-call tax is paid again inside every
extra cycle.

I independently reproduced claude-1's 76-idea split. PRIMARY, quoted read-only command output; dates
came from `00-prompt.md`, round counts from `round-NN/` directories, and bytes from all Markdown under
each `review/` tree:

```text
included=76
older ideas=55 design=1.418 review=1.618 max=5 review_bytes=20236.8
recent ideas=21 design=1.619 review=5.143 max=24 review_bytes=146290.0
```

I then compared review-round count with the maximum numbered `## Fix-up cycle` in each recent
`IMPLEMENTATION.md`. PRIMARY, quoted output:

```text
recent_n=21 pearson_review_rounds_vs_max_fixup=0.969
2026-07-29 review=24 max_fixup=27 track=deliberation integrate-parley-bidding-addon
2026-07-29 review=21 max_fixup=16 track=standard skills-cli-install-path
2026-07-28 review=11 max_fixup=14 track=deliberation parley-design-skills
2026-08-07 review=9 max_fixup=9 track=deliberation meta-protocol-change-global-core-protocol
```

The raw record corroborates the first row: review consensus says cycle 24 and the implementation
reaches fix-up cycle 27 (`integrate-parley-bidding-addon/review/consensus.md:3` and
`IMPLEMENTATION.md:1881`). The facilitator's A/B timing remains **facilitator-reported, not rerun by
me**: full 98.7 s (27.3–105.3) versus excerpt 29.9 s (21.1–39.2), `n=3` each.

The **one settling measurement** should be a paired critical-path trace of one representative,
completed `deliberation` idea. For every real agent/facilitator invocation, record monotonic start,
first-token, and end times; phase/round/fix-up; exact model/effort; uncached and cached input tokens;
output tokens; retries; and queue delay. Replay that same frozen invocation once in randomized order
with only the full protocol replaced by the relevant exact excerpt. Aggregate by the actual round
critical path—maximum participant duration per parallel round, plus sequential drafting/fix-up
spans—not by summing parallel agents. Report exactly two attributable totals:

```text
protocol tax = critical-path(full packets) - critical-path(excerpt packets)
repeat-cycle tax = critical-path time in review/fix-up cycles after the initial review + budgeted fixes
```

Keep human idle time as a third, separately labelled bucket. Whichever of the first two totals is
larger settles the technical dominance question. Token or byte regression alone does not.

## Q2 — Replace full history with an anchor and a frontier

`gatherPriorRounds` currently concatenates rounds `1..N-1` (`internal/runner/runner.go:936-965`) and
the prompt orders every artifact read (`runner.go:968-990`). Phase 2 requires addressing every active
peer, not every historical version (`parley-deck/COOPERATION.md:326-352`). Change the runtime packet
as follows:

- **Round 2:** every participant's Round-1 artifact in full.
- **Round N >= 3:** every participant's Round-1 artifact in full (the independence/correlation
  anchor), plus every participant's Round-(N-1) artifact in full (the current frontier), plus a small
  manifest of every omitted artifact: path, participant, round, byte length, and SHA-256.
- **Triggered expansion:** if the frontier cites, disputes, supersedes, or challenges the provenance
  of an older artifact, include that older artifact verbatim. Missing frontier files or malformed
  `responding-to` data fail back to full history.
- **Consensus:** retain one full-history read for now. It is linear once, not repeated quadratic
  history. This preserves the strongest audit for §15.6 while a carry-forward scheme is unproven.

This is intentionally not a semantic digest. It bounds normal round input to two generations,
preserves the raw independent proposals needed to detect correlated agreement, and leaves every
omitted source content-addressed and readable on demand.

What can break is an objection, qualifier, or contradiction that appeared only in an intermediate
round and was silently dropped from the frontier. The final full-history audit catches that late,
but may have to reopen a round. A later optimization may replace that audit with participant-owned,
machine-validated issue records—stable ID, exact proposition, OPEN/RESOLVED/DEFERRED state, owner,
source path/hash, and resolution locator. An issue may disappear only through an explicit disposition;
invalid carry-forward falls back to full source. A facilitator- or model-written prose digest is not
acceptable because it can manufacture agreement.

Apply the same bounded-history principle to review. `gatherReviewContext` currently embeds `FINAL.md`,
the ever-growing `IMPLEMENTATION.md`, and every earlier review round (`phase58.go:276-306`). Fixing
only design `gatherPriorRounds` would leave the measured hotspot untouched. Never cut `FINAL.md`, the
current implementation diff/check results, the latest full review round, or any open finding.

## Q3 — Fix-up termination without weakening refutation-default review

The unbounded `deliberation` rule is a driver, but not **the** driver. PRIMARY, quoted grouping of the
same 21 recent ideas:

```text
deliberation n=12 mean_review=5.500 max=24 total=66
fast n=1 mean_review=2.000 max=2 total=2
standard n=8 mean_review=5.000 max=21 total=40
```

The manual records exceed even the bounded tracks, while the auto-driver already escalates after its
configured maximum. The repair is one termination state machine shared by both paths:

1. Set an initial escalation budget of 1 fix-up for `fast`, 2 for `standard`, and 3 for
   `deliberation`. A ceiling is never a pass criterion.
2. At the ceiling, halt and produce a trajectory record: each unresolved finding, severity, whether
   it is an acceptance-criterion failure, whether it was introduced by the latest fix, affected code,
   and the attempted root-cause fixes.
3. The operator may choose exactly one of: abandon/defer the implementation; approve a new root-cause
   plan with observable exit checks and **at most two additional cycles**; or split eligible low
   severities into linked follow-up ideas. A second extension is forbidden; unresolved blocking work
   requires a new idea/redesign rather than continuing the patch spiral.
4. CRITICAL/MAJOR findings, acceptance-criterion failures, security/correctness defects, and any
   regression introduced by the latest fix remain blocking regardless of cycle number. A post-budget
   MINOR/NIT is still reported, but does not reopen the current implementation unless it is such a
   regression or criterion failure; it receives an explicit follow-up/disposition instead. Under an
   opted-in `strict_gate`, all findings remain blocking unless the operator uses the existing explicit
   ruling path.

This changes continuation, not inspection. Reviewers still run full-scope refutation attempts, may
report every severity, and may rebut prior dispositions (`COOPERATION.md:527-556`). A MINOR in round
19 therefore remains visible but no longer automatically has the same continuation effect as a
CRITICAL. The stop can never convert unresolved blocking findings into “complete.”

## Q4 — Caveman-style semantic compression

**Recommendation: against deployment for the protocol, and against using it as the primary round-file
optimization.** The claimed 40–58% token saving is from the facilitator brief and is not independently
verified here. Articles and auxiliaries are not filler in a normative document: “must,” “must not,”
“may,” “is not,” scope, exception, and actor phrases encode the rule. A lossy compressed protocol
would be a fourth representation beside the global core, generated deck view, and bundled fallback;
hash equality could prove byte identity but not semantic fidelity.

Peer rounds are analysis rather than law, so the risk is lower but not low. They carry negation,
uncertainty, minority objections, provenance, and counterfactuals; stripping those can create exactly
the correlated agreement §15.6 guards against. If tested at all, compression should be ephemeral and
prompt-only for **older, explicitly resolved explanatory prose**. Never compress Round 1, the latest
round, an open objection/finding, evidence and provenance locators, quoted user direction, code/path/hash,
numbers, severities, conditions, negation, or normative text. The canonical raw file remains available,
and any dispute triggers full-source loading. Given the anchor+frontier design, that safe remainder may
be too small to justify another mechanism.

The deciding test is a blinded replay on real Parley packets, comparing (a) full source, (b) exact
anchor+frontier selection, and (c) anchor+frontier plus caveman compression:

- For normative samples, extract `(actor, modal, action, object, condition, exception)` tuples and
  answer rule-application cases. **One changed tuple or answer rejects protocol compression.**
- For peer samples, require 100% preservation of open objections/findings, severity, negation,
  numerical thresholds, provenance locators, and proposal-family relationships. Measure independent
  refutation yield and final decision, not summary similarity.
- Randomize arm order, hold model/effort/output cap constant, record cached/uncached tokens, first-token
  and completion time, and use enough repetitions that the latency interval no longer spans no effect.
  Deploy for old resolved analysis only if arm (c) materially beats arm (b) in wall-clock while meeting
  every semantic gate. Token reduction without latency and decision preservation is not a win.

## Q5 — Intervention ranking by expected saving / rule-weakening risk

| Rank | Intervention | Expected saving | Risk to real-error detection | Reason |
| --- | --- | --- | --- | --- |
| 1 | Exact phase-scoped protocol packet, with the skill/facilitator instruction changed to point to it | High on every call; the facilitator's small A/B found a 3.3x median latency ratio | Low if selection is extractive, dependency-checked, and fails back to full authority | It removes irrelevant text without rewriting an applicable rule. |
| 2 | Anchor + frontier context in design and review, with full-source triggers and a final full-history audit | Very high in later rounds; removes repeated quadratic history | Medium | Intermediate nuance can disappear, but raw anchors, frontier, manifests, and fallback constrain the loss. |
| 3 | Enforced finite fix-up budget plus severity/scope-aware post-budget disposition | Largest absolute end-to-end saving in churned ideas | Medium-high | Late review has caught real defects; escalation, strict-gate preservation, and no auto-pass are mandatory. |
| 4 | Concision guidance and regression alarms on round-file size | Moderate | Medium | Useful against verbosity, but a hard cap can suppress the only dissenting argument. Alarm first; never truncate. |
| 5 | Caveman compression of old resolved analysis only | Unknown incremental saving after ranks 1–2 | Medium-high | Lossy semantics and another transformation layer; experiment only. Never apply to normative text. |

Never cut or suppress: the authoritative protocol; applicable normative modals, negation, conditions,
and exceptions; Round-1 independent proposals; every active objection and finding; provenance and
verdict-conflict data; the §15.6 adversarial-alternative/family audit; `FINAL.md`, current code diff,
acceptance criteria and check results; full-scope refutation attempts; explicit user rulings; canonical
ownership/no-secrets/human-brake rules; or access to raw historical artifacts.

## Cost of the selective protocol read

I did not read `COOPERATION.md` in full. I opened only §4.0, Phase 2, Phases 6–8 including stopping
judgment, and §15.6 (`parley-deck/COOPERATION.md:200-257,326-356,501-686,1336-1369`), plus the exact
runner/driver paths cited above. That cost me confidence that no unmarked dependency elsewhere changes
an edge case. It did not cost the central answer; in fact the targeted source read exposed the
auto-driver's three-cycle default, which a protocol-only reading would have missed. The existing
Graphify graph confirmed the `PolicyFor` -> driver/review/fix-up relationship, but warned that its
installed skill metadata (0.8.38) trails the package (0.8.41), so I treated source lines—not graph
inference—as authority. This residual uncertainty is another reason all proposed selectors fail back
to full source rather than silently omit it.
