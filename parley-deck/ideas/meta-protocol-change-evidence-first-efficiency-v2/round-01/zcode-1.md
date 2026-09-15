---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency-v2
round: 1
date: 2026-09-15
---

## Summary

I accept the narrow amendment: the quorum is codex-1, claude-1, kimi-1, zcode-1
(user waiver recorded in 00-prompt.md); the full-four arm supersedes only the
prospective full-six references in D8/AC-X1 of the signed audit FINAL; Hermes'
historical contributions keep their attribution and I own none of them
retrospectively. Everything else in
`../meta-protocol-change-evidence-first-efficiency/FINAL.md` and the packet
FINAL stays binding: 12 tasks, three arms, equal enforced ceilings, owned
artifacts, blind nonauthor grading, intention-to-treat denominators, paired
analysis, and the exact packet trial (six AB/BA pairs in each of phases 1 and
6, three canaries plus full control, ship gate R ≤ 0.50 in both phases,
correctness veto, disputed (0.67, 0.80] band returned to the user with both
positions). My concrete reservation: USD 15 total most likely cannot run the
full frozen scope; the design must make that outcome honest rather than shrink
it.

## Proposed approach

**Rotation (freeze before treatment).** Identities C, K, Z, X. Full-four: all
four on every task. Solo: task i (1–12) → identity `i mod 4` (three solos
each). Duo: six unordered pairs, task i → pair `i mod 6`, each pair twice,
first-author seat counterbalanced by task parity; every identity appears in six
duo cells. Drafter in multi-participant cells rotates by
`(task, arm) mod n`, critic duty rotates among the remainder; each identity
drafts six of the 24 multi-participant cells. Arm order is a 12×3 Latin square
(four tasks per arm-position each), fixed seed, recorded in the write-once
freeze file with task bytes, rubrics, models/efforts, hashes and ceilings.

**Grader disjointness.** Two grading-only configurations, frozen in the freeze
file, disjoint from the four candidate configurations and from task
authorship; they see anonymized outputs, task bytes, rubric and deterministic
test results only — never arm/model/agent/cost/order metadata; the hash mapping
lives outside candidate workspaces, and hidden tests/keys never enter them.
Style leakage risk is reported as the DRAFT requires. Disagreement follows the
predeclared DRAFT rule (critical-invariant conflict, pass/fail split, or >15
rubric points): one third blind review if contingency remains, otherwise an
explicitly unresolved cell; never averaged into a score.

**Enforceable spend policy (USD 15) and 15-minute clock.** Words/token limits
are guidance, not a cap. Enforcement is a durable cross-experiment ledger with
per-call pre-charge: each launch first reserves a conservative per-config
amount (p95 of measured same-config calls; cold-start reservations set in the
freeze file); terminal `agent.usage` reconciles, and unknown or missing cost —
e.g. any adapter whose CLI reports no cost field — keeps the full reservation
charged, never zero (audit FINAL D2). When spent plus outstanding reservations
reaches USD 15 the boundary fails closed: no further launches. The 15-minute
task-arm limit is one shared monotonic stopwatch per cell, started at first
launch, wall-clock enforced by the runner, never reset per call, orchestration
included; frozen segment budgets (proposed: proposals ≤ 6, cross-review ≤ 4,
final ≤ 2, acceptance ≤ 1, orchestration slack ≤ 2 minutes) make
deadline-no-output a retained outcome, not an extension.

**Feasibility gaps, stated honestly.** USD 15 across 36 pilot cells, 28
packet/canary/control calls, grading and retries implies ≲ USD 0.20 per
chargeable event before any retry — below realistic per-call cost for
opus/gpt/kimi/glm configs carrying full protocol context (our own actual zcode
terminal envelope recorded 4.1M input tokens in one turn). Heterogeneous cost
coverage worsens it: conservative reservations outrun actuals. The 15-minute
full-four cell leaves roughly three minutes per participant phase, so many
cells will end deadline-no-output. Proposal: packet trial first (it is the
ratified ship gate, ~28 calls), then pilot cells in the frozen Latin order,
batched grading, small contingency — allocation numbers belong to the freeze
file. Expected honest outcome: a large not-run fraction. Unrun cells stay
not-run with planned denominators intact; no scope shrink, no rerun of
negatives, no complete claim; preserved cells can proceed only if the user
later raises the ceiling.

## Existing alternatives

- Spend-stop semantics already specified: audit FINAL D2 (visible stop or
  documented conservative reservation, never zero; fail-closed telemetry at
  spend boundaries). Only the cross-run USD ledger is new hand-building.
- Machine hooks exist: §4.0 LE-5 loop budgets (`--max-wall-clock`, steps,
  wall-clock ceilings; cost enforcement telemetry-gated on `agent.usage`).
- Preparation tooling: sibling evaluation repo `delivery/2026-09-05/pilot/`,
  `scripts/pilot_{acceptance,grading,analysis}.py`; DRAFT copied to
  `source-context/pilot-preparation.md` (sha256 6fe1e8d8…) — reuse; it is
  preparation, not preregistration, so a new write-once freeze file is still
  required.
- Grading design: DRAFT §Grading (two fresh grading-only configs, anonymity
  set, disagreement thresholds) — adopt unchanged.

## Concerns / open questions

Who authored the 12 tasks and hidden tests — grader disjointness needs that
list frozen. Cold-start reservation values need a measured or explicitly
conservative basis. Readiness of the three in-flight implementation calls must
come from their actual completion, not silence. Nothing here may re-decide the
packet ship gate or disputed band.

## Risks

Post-truncation small n leaves paired analysis underpowered — report as-is.
Equal elapsed ceilings are not equal compute (DRAFT); identity can leak through
style. Declared interest: I am the incoming member whose adapter's usage parser
is still under correction; until it reports cost, my arm's spend is reserved
conservatively, and that interaction should be recorded, not hidden.

## SELF-CORRECTION (2026-09-15, appended; every byte above is an unmodified prefix)

Three scoped challenges arrived; each item below names the statement it replaces.
All three are weakenings or arithmetic corrections, so per §15.1 they take effect
immediately. The scope parameters are unchanged by this correction: USD 15 total
ceiling, one shared 15-minute stopwatch per task-arm, the original 12 tasks, and
the exact separate packet trial.

### SC-1 — the USD 15 boundary caps reservations, not provider spend

Replaces: "each launch first reserves a conservative per-config amount (p95 of
measured same-config calls …)", and "When spent plus outstanding reservations
reaches USD 15 the boundary fails closed" insofar as it read as an enforced
monetary ceiling on provider billing.

Corrected position: a p95 reservation is a quantile estimate — by construction a
fraction of same-config calls exceed it, and heavy cost tails make exceedance worse
— and while a call is in flight the provider's actual charge is invisible to the
ledger until terminal `agent.usage` reconciles. The ledger bounds reserved plus
settled spend and gates further launches; it does not bound what the provider
actually bills inside that window. My round cites no provider-side enforced USD or
token ceiling on in-flight calls; absent one, actual-spend feasibility at USD 15 is
UNVERIFIED, not enforced. The fail-closed boundary on known exposure stands, and my
"large not-run fraction" expectation is strengthened, not weakened. Dependent
calls stay not-run, consistent with my not-run posture; the packet trial should
record per-config reservation-vs-actual divergence so the freeze file's feasibility
statement rests on measurement. I find no counterexample to the challenge and
concede it.

### SC-2 — the 4.1M-token envelope is context-weight anecdote, not a cost or impossibility measurement

Replaces: "(our own actual zcode terminal envelope recorded 4.1M input tokens in
one turn)" as support for "≲ USD 0.20 per chargeable event … below realistic
per-call cost".

Corrected position: that envelope was largely cached tokens, which price far below
uncached input, and it measured a source-audit task shape, not the pilot's. One
largely-cached envelope establishes neither the frozen pilot's per-cell cost nor
that USD 15 cannot fund the frozen scope. I withdraw the impossibility reading.
What survives is only the direction: the per-event budget is small relative to
plausible uncached full-context calls, so feasibility is unmeasured and likely
tight — a hypothesis the packet trial's terminal `agent.usage` must settle per
config before pilot cells spend. No counterexample: conceded.

### SC-3 — task parity does not counterbalance duo first-author seats; formula corrected

Replaces: "Duo: six unordered pairs, task i → pair `i mod 6`, each pair twice,
first-author seat counterbalanced by task parity".

Counterexample conceded: pair p recurs at tasks i and i+6, and 6 is even, so both
occurrences share parity — a parity-keyed seat assigns the same first author to
both repetitions of every pair, so nothing alternates. Corrected formula: seat by
repetition half — first author = `pair[p][(i-1)/6]` (integer division), i.e. the
pair's first-listed member holds the first-author seat for its task in 1–6 and the
second-listed member for its task in 7–12; equivalently, swap seats on the second
occurrence. The independent balance claim is unaffected and retained: each
identity sits in 3 of the 6 pairs × 2 repetitions = 6 duo cells. The full-four
drafter rotation `(task, arm) mod n` and the "drafts six of the 24
multi-participant cells" claim are not touched by this correction.

### Standing after correction

Unchanged: USD 15 total ceiling; one shared 15-minute stopwatch per task-arm; the
original 12 tasks; the exact separate packet trial as the ratified ship gate; unrun
cells stay not-run with planned denominators intact. Changed: monetary feasibility
is UNVERIFIED absent a provider-side bound; pilot cost is unmeasured pending
packet-trial `agent.usage`; the duo first-author seat formula is as corrected
above.
