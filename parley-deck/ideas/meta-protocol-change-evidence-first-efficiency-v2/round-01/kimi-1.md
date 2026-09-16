---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency-v2
round: 1
date: 2026-09-16
---

# kimi-1 — Round 1: Prospective quorum amendment and a gated USD 15 experiment plan

Author: kimi-1 via Kimi CLI (requested model kimi-code/k3, effort max from configured
route). I assert nothing about the provider-resolved identity of this session; that
lives in launch telemetry, not in my self-report. Written independently: no other
participant's round-01 file was read, no subagents, no peer output.

Verdict discipline: every material factual assertion below is tagged `[SOURCE]` with
a locator, `[ASSUMPTION]`, or `[PROPOSAL]`. I own the `[SOURCE]` assertions where they
first appear here, so per §15.1 I issue no CONFIRMED/WRONG/UNVERIFIED verdict on any
of them; the locators are supplied so a non-owner can reach PRIMARY. Absence-of-evidence
statements ("no PRIMARY evidence of X is present in the sources read") are observations
about this round's inputs, not verdicts on X.

## Summary

The amendment should do exactly two things and no more: (1) supersede, prospectively
only, the full-six roster references in the original FINAL's D8/AC-X1 with the
user-waived full-four quorum (codex-1, claude-1, kimi-1, zcode-1); (2) add a
budget-enforcement gate in front of an otherwise unchanged experiment — 36 pilot
cells, the exact 28 packet/canary/control calls, two blind grader configurations —
sized by a pre-registered call plan whose per-call conservative reservations must sum
to ≤ USD 15 before the first treatment call. My central position: fit of the
unchanged experiment inside USD 15 is an open gate, not a premise. No PRIMARY evidence
available to this round establishes per-call dollar cost for the four configured
CLIs, and the launch-path budget boundary visible in the frozen source is not shown
to be an integrated, USD-denominated ceiling. If the reservation sum does not fit,
the amendment's answer is planned denominators plus honestly recorded not-run cells
and a costed shortfall report to the user — never silent scope shrink, mock output,
or extra spend.

## Source claims relied on (locators cited; verdicts left to non-owners)

- S1. Original FINAL D8/AC-X1 require twelve real tasks in solo, duo and full-six
  arms with equal enforceable ceilings, controlled full context, rotation, owned
  phase artifacts, and honestly retained failures/not-run cells
  (`../meta-protocol-change-evidence-first-efficiency/FINAL.md` D8 lines 149–157,
  AC-X1 line 236).
- S2. The original FINAL's D5 and AC-P2 bind the exact packet experiment: phases 1
  and 6, six matched AB/BA pairs EACH, packet generation inside the packet arm, three
  canaries plus a full control, non-implementer recomputation
  (`FINAL.md` D5 lines 99–113, AC-P2 line 232).
- S3. Ship gate R ≤ 0.50 in both phases with zero correctness misses and passing
  canaries; the disputed band (0.67, 0.80] goes to the user with both recorded
  positions (codex/kimi refute > 0.80; hermes > 0.67); the measured value is never
  rounded up to a planning estimate (`FINAL.md` lines 185–190;
  `../meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md` lines 113–149).
- S4. The user's direction, as recorded in the v2 kickoff: USD 15 TOTAL for
  prospective live experiments and 15 minutes for one task in one arm, shared by all
  participants and phases in an arm including orchestration, not reset per call;
  unknown spend must stop or carry a defensible enforced conservative reservation,
  never zero (`00-prompt.md` lines 20–23, 41–44).
- S5. The v2 kickoff supersedes only prospective roster references; historical
  artifacts and signatures stay intact; this amendment is not implementation
  acceptance (`00-prompt.md` lines 26–33, 94–99).
- S6. Pilot preparation (advisory, exact-byte copy, sha256 6fe1e8d8…d243 per
  `source-context/provenance.json`): twelve tasks each assigned all three arms with
  counterbalanced order; full protocol context in every arm with packet generation
  shadow-only; fixed phases (independent proposals, scoped cross-review, final draft,
  explicit dispositions); failed members retained intention-to-treat
  (`source-context/pilot-preparation.md` lines 23–33).
- S7. The same preparation: available CLIs do not all expose an enforceable total
  token ceiling; equal wall time is not equal compute; the proposed main ceiling is
  equal per-task arm elapsed time with conservative spend reservations under an
  overall stop that may leave arms not-run; missing cost is missing, not zero
  (lines 35–47, 68–71).
- S8. Grading per preparation and D8/AC-X2: two fresh grading-only configurations
  disjoint from task authors and candidate authors; observers, not consensus votes;
  anonymized outputs; deterministic acceptance tests primary for code; predeclared
  adjudication (critical-invariant conflict, pass/fail disagreement, or weighted
  rubric difference > 15 points → third blind review or explicit unresolved cell,
  never averaged); paired task-level differences with fixed-seed bootstrap
  (`pilot-preparation.md` lines 49–71; `FINAL.md` AC-X2 line 237).
- S9. Runtime boundary, from the frozen isolated source snapshot provided to this
  launch (integration_head 2aa2c41a…, marked "isolated, unintegrated" in its
  provenance): the launch path records the request before spawn, then reserves a
  launch budget pre-start (`internal/runner/telemetry.go` lines 92–103); on finish
  it settles the budget with the reported `usage.CostUSD`, and a settle failure is
  an integrity error with the reservation retained (lines 146–153); a pre-start
  budget refusal is its own failure class `budget_refused` (lines 189–191);
  hard-timeout enforcement exists (`context.DeadlineExceeded` → status failed /
  timeout, lines 169–170); usage carries nullable `CostUSD`, `CostBasis`,
  `Coverage` (`internal/telemetry/record.go` lines 46–57).
- S10. In that same frozen set, the `LaunchBudget` type's definition does not appear
  (no budget symbol in the provided `internal/runner/launch.go`), and the snapshot
  is explicitly unintegrated; per the launch task, the current source has no real
  positive bound runner relaunch path and the parent-lineage implementation is
  separately in progress. Therefore the provided source does not establish an
  integrated, USD-denominated enforceable ceiling today.
- S11. Original FINAL D2: with a configured monetary ceiling, unknown spend requires
  a visible stop or a documented conservative reservation, never zero; telemetry
  writes required for a spend/launch boundary fail closed
  (`FINAL.md` lines 48–54).
- S12. Three measured CLI implementation calls were still running when this
  amendment opened; readiness comes from those completed calls or a new bounded
  check, never from silence (`00-prompt.md` lines 79–82).

## Assumptions (explicit, challengeable)

- A1. The 15-minute ceiling binds pilot task-arms. The packet experiment keeps its
  ratified design with no new truncation ceiling, because D5 requires preserving it
  exactly and R is a wall-clock ratio whose measurement a truncation cap would
  corrupt. The USD 15 total binds all prospective treatment spending: pilot, packet,
  grading, retries. This reading reconciles S3/S4 without editing the packet design;
  cross-review may prefer the stricter reading (cap packet runs too), which would
  make some packet cells timeout cells and is still honest, but changes the
  experiment and therefore needs explicit consensus, not silence.
- A2. The completed implementation calls of S12, and any bounded readiness checks,
  are not treatment calls: they produce no candidate outputs and touch no task
  bytes. Conservatively, their prospective spend is still counted against the USD 15
  ledger (counting more, not less, is the safe direction).
- A3. Task authors for the twelve fixtures are named at freeze time from the
  evaluation repository's history; grader disjointness is then checkable. The
  sources read do not name them.
- A4. The four participant CLIs can be launched through the instrumented path with
  unique invocation IDs as D2 requires; the bounded readiness check demonstrates
  this per CLI before freeze.

## Proposed approach

### P1. Narrow amendment text

Amend D8 and AC-X1 prospectively: "full six active roster IDs" becomes "full four —
codex-1, claude-1, kimi-1, zcode-1" for the prospective pilot only, citing the
user's recorded waiver. Nothing else in the original FINAL changes: D1–D7, D5's
packet experiment, every other acceptance criterion, the idempotence/recovery terms,
and the disputed speed band all stand (S1, S2, S3, S5). Historical Hermes
contributions keep attribution; the recorded hermes refute position at R > 0.67
remains a historical position that travels to the user with the codex/kimi position
if the measurement lands in (0.67, 0.80] — zcode-1 inherits no position in that
dispute, and the implementer still must not resolve it (S3). This amendment claims
no implementation acceptance; original implementation PR #73 stays incomplete until
its own gates pass (S5).

### P2. Twelve-task three-arm allocation (decision 1)

Identities C=codex-1, L=claude-1, K=kimi-1, Z=zcode-1. Tasks T1–T12 in a frozen
seeded order. Cells are task-major (all three arms of a task adjacent) so a budget
stop leaves complete task-triples for paired analysis; arm order within each task
rotates on a Latin square (tasks 1,4,7,10: solo→duo→full; 2,5,8,11: duo→full→solo;
3,6,9,12: full→solo→duo).

- Solo arm (3 tasks per identity): T1 C, T2 L, T3 K, T4 Z, T5 L, T6 K, T7 Z, T8 C,
  T9 K, T10 Z, T11 C, T12 L.
- Duo arm (each of the six unordered pairs exactly twice; drafter swaps between the
  two occurrences, so each agent is in 6 duo cells and drafts 3): T1 C-L, T2 C-K,
  T3 C-Z, T4 L-K, T5 L-Z, T6 K-Z, T7 C-L, T8 C-K, T9 C-Z, T10 L-K, T11 L-Z, T12 K-Z.
- Full-four arm: drafter rotates C, L, K, Z across T1–T12 (3 each); the other three
  are critics; proposal and cross-review authorship order rotated per task.

Bounded variant per cell, recorded as the exact experimental protocol variant per D8
(S1): solo = 2 calls (proposal, final); duo = 6 (2 proposals, 2 scoped cross-reviews,
1 final, 1 non-drafter disposition); full-four = 12 (4 proposals, 4 cross-reviews,
1 final, 3 dispositions). Within a phase, independent calls run concurrently; phases
are sequential. The 15-minute clock starts at the first treatment launch of the cell
and ends at the last disposition's terminal record, including orchestration gaps;
overrun records a timeout cell, retained intention-to-treat, never dropped (S4, S6).
Total pilot plan: 240 calls. Retries: predeclared cap of 1 per call, only for
provider/process failure or timeout, charged to the ledger, recorded with `retry_of`
(D2: retries remain separate spent attempts).

### P3. Full context in every pilot arm

Every pilot launch in every arm carries the live full protocol context
(context_mode=full attestation); packet generation stays shadow-only and never
enters participant inputs, per D4 and the preparation (S6). The 28 packet-experiment
calls are the only packet-context treatment calls, and they belong to that
experiment, not to the pilot.

### P4. Packet experiment preserved exactly (decision 2 scope)

28 calls, unchanged: 12 AB/BA pairs (6 per phase across phases 1 and 6) with agent,
model/effort, task, output cap and workspace snapshot held constant and packet
generation inside the packet arm; 3 canary replicates, all required to pass; 1 full
control; §6/§14/§15 obligations seeded and checked on every run (S2, S3). Both calls
of a pair run adjacently so a stop leaves whole pairs. Canaries and control run
after the pairs; if the ledger cannot cover them, the gate result is not-run and the
optimization stays unshipped — a failed or unrun gate never passes (S3). Thresholds
preserved verbatim: ship R ≤ 0.50 in both phases AND canaries pass AND zero
obligation misses; any correctness miss refutes at any speed; a landing in
(0.67, 0.80] returns the measured number with both recorded positions to the user;
no rounding up to planning estimates; per-call saving never presented as a
whole-idea saving (S3). Non-implementer recomputation of both R values from the raw
log remains required (S2).

### P5. Two blind grading-only configurations (decision 3)

Selection rule, frozen before treatment: two grader configurations that (a) share no
identity with any of the four participants, (b) authored neither any task nor any
candidate output (A3), (c) run grading-only with no write access to candidate
workspaces, (d) are recorded with model/effort in the freeze file. I do not name
specific grader models: the sources read do not establish which non-participant
configurations are invokable, so naming one would be an unsupported claim; the
bounded readiness check enumerates candidates. Graders are observers, never quorum
votes (S8). Each grades anonymized candidate outputs, task bytes, rubric and
deterministic test results, with the identity mapping held separately under hashes;
style/wording leakage risk is reported, not assumed away (S8). Grading plan: 36
pilot outputs × 2 graders = 72 calls, plus a predeclared adjudication reserve of at
most 6 third-blind-review calls; disagreement triggers follow the preparation
exactly, and if the reserve is exhausted the affected cells are recorded unresolved
rather than buying more calls (S8). Packet-experiment outputs are not blind-graded;
their checks are the ratified seeded obligations plus non-implementer recomputation
(S2).

### P6. Actual dollar enforcement versus reported cost (decision 2)

Reported cost is telemetry; enforcement is a boundary. The plan keeps them distinct
(S9, S11):

1. Ledger. One durable experiment ledger at USD 15. Before every prospective launch
   (treatment, grading, probe, retry), a conservative reservation is charged; after
   completion it settles to the reported `CostUSD`. If reported cost is null or its
   coverage is missing, the reservation stands as spent — it never settles down to
   zero on unknown (S4, S7, S11).
2. Reservation sizing without inventing prices. Per CLI/model/effort, the
   reservation is twice the maximum reported cost observed for that configuration
   in the completed implementation calls (S12) or bounded readiness checks. For a
   configuration with no observation, twice the maximum observed across all
   configurations. If no observation exists at all, there is no basis for any
   nonzero reservation and treatment does not start — the blocker is escalated with
   the arithmetic shown. No provider price list is consulted or asserted; token
   counts are never converted to dollars.
3. Fit gate. The freeze file contains the full call plan — 240 pilot + 28 packet +
   72 grading + ≤ 6 adjudication = ≤ 346 calls — with the per-call reservation for
   every entry. First treatment call is blocked unless the reservation sum ≤ USD 15.
   As arithmetic only (not a feasibility claim): fit requires average settled cost
   per call ≲ 15/346 ≈ USD 0.043. Whether the four configured CLIs at the requested
   models/efforts deliver that is unknown to this round.
4. Demonstrated boundary. The enforcement mechanism must be demonstrated, not
   asserted, before freeze: a bounded non-treatment fixture in which a launch whose
   reservation would exceed the remaining ledger is refused pre-start, with the
   invocation records showing request → refusal → reservation retained. The frozen
   source shows the shape of this boundary (reserve pre-start, settle on finish,
   `budget_refused`, fail-closed settle) but does not establish it as integrated
   and USD-denominated (S9, S10). If the demonstration cannot be produced, no
   treatment call runs.
5. Ordering and stop. Packet experiment first (bounded 28 calls, decides a standing
   disputed gate), then the pilot in the P2 order. Before each launch: if remaining
   ledger < that call's reservation, stop; the affected and all subsequent cells are
   recorded not-run with cause. Hitting the ceiling never marks anything complete.

### P7. Honest not-run reporting (decision 4)

Every planned cell carries a status from the predeclared enum {not-run,
provider-failure, timeout, protocol-block, valid} (S7). If the fit gate fails, the
deliverable is: the frozen plan with all planned denominators (36 pilot cells, 28
packet calls, 72 grading calls), every cell marked not-run, the reservation
arithmetic showing the shortfall, and the escalation to the user. No task count
reduction, no pair reduction, no mock substitution, no unbudgeted spend (S5, S7).
Analysis, where cells did run: paired task-level deltas with fixed-seed bootstrap
over complete cells only, denominators shown as planned-vs-run, cost reported with
its basis and coverage alongside any quality-per-dollar figure, and no vendor-IQ or
causal-historical claims (S1, S8).

### P8. Freeze (decision 5)

Freeze file v1 is write-once and hashes: the twelve task definitions and hidden
acceptance tests, rubric, seeded task order and the P2 allocation table,
model/effort per participant and per grader, per-CLI enforceable ceilings (elapsed
enforcement per S9; the absence of equal token ceilings stated explicitly per S7),
the call plan with reservations, the ledger path, the stop rule, the not-run enum,
and the source/test hashes. Reuse the existing preparation scripts
(`scripts/pilot_{acceptance,grading,analysis}.py` in the sibling evaluation
repository, per `00-prompt.md` lines 66–71) rather than writing parallel ones.
Amendments land as new versions and invalidate the affected comparisons
(`pilot-preparation.md` lines 3–6).

## Existing alternatives (§15.6a)

Mechanisms this plan builds on, and what the toolchain already ships, with locators
consulted:

- ALT-A Per-call elapsed enforcement — SHIPPED: the runner classifies
  `context.DeadlineExceeded` as a timeout failure (frozen
  `internal/runner/telemetry.go` lines 169–170). Proposal builds on it for the
  15-minute ceiling; the arm-level shared clock across sequential calls is
  orchestration accounting the proposal adds.
- ALT-B Pre-start launch reservation with settle-on-finish and a `budget_refused`
  class — PRESENT IN FROZEN SOURCE, UNINTEGRATED (frozen
  `internal/runner/telemetry.go` lines 92–103, 146–153, 189–191; definition of
  `LaunchBudget` absent from the provided set; S10). REJECT treating it as an
  available USD ceiling; ADOPT it as the shape of the boundary to be demonstrated
  per P6.4.
- ALT-C Cost telemetry (`CostUSD`/`CostBasis`/`Coverage`, nullable) — SHIPPED as a
  record (frozen `internal/telemetry/record.go` lines 46–57). REJECT using it as an
  enforcement mechanism; it is the settle input only.
- ALT-D Charged fixup/cycle budget (`internal/budget`, imported at frozen
  `internal/trajectory/verification.go` line 15) — SHIPPED as a monotonic cycle
  counter per D6. REJECT repurposing it as a dollar cap; different unit.
- ALT-E Existing preparation scripts and DRAFT — PRESENT (locators in
  `00-prompt.md` lines 66–71 and `source-context/provenance.json`). ADOPT reuse.
- ALT-F Post-hoc token totals, output-word guidance, or CLI cost self-estimates as
  the spend cap — REJECTED: the kickoff states these are not a hard spend cap
  (`00-prompt.md` lines 49–51), and the preparation states equal wall time is not
  equal compute and missing cost is not zero (lines 35–47).
- ALT-G Pilot-first ordering — REJECTED in favor of packet-first (P6.5): the packet
  plan is exactly 28 bounded calls and decides a standing disputed ship/refute gate,
  so it maximizes the chance that at least one experiment completes intact under a
  hard stop. Flagged for cross-review; the allocation and gates above do not depend
  on this choice.

## Observable gates for this amendment

- G1 Freeze file v1 exists, write-once, containing every item in P8; treatment
  before freeze is a protocol-block event, not a result.
- G2 The P6.4 refusal demonstration passed, evidenced by invocation records
  (request, pre-start refusal, reservation retained).
- G3 Every prospective launch has a unique invocation ID, a request record before
  spawn, a terminal record after its real outcome, reservation-then-settle ledger
  entries, and null-cost settles equal to the reservation — never zero.
- G4 Every pilot launch attests context_mode=full; no packet text enters pilot
  participant inputs.
- G5 The packet run set is exactly 24 paired + 3 canary + 1 control calls, pairs
  adjacent, with the S3 thresholds applied verbatim and non-implementer
  recomputation of both R values.
- G6 Both grader configurations satisfy the P5 disjointness rule, recorded in the
  freeze file; adjudication follows the S8 triggers with the ≤ 6-call reserve.
- G7 The report shows planned-vs-run denominators (36 / 28 / 72), per-cell status
  from the P7 enum, paired fixed-seed analysis, cost basis and coverage, and makes
  no vendor-ranking or historical-causal claim.
- G8 No historical artifact is modified; this amendment is design acceptance only
  and asserts nothing about implementation PR #73.

## Concerns / open questions

- Q1. Fit. If the completed implementation calls (S12) report per-call costs far
  above USD 0.043, the unchanged experiment cannot fit USD 15. My position: report
  P7 honestly rather than shrink; the user decides between more budget and a
  recorded scope change. I hold this as a considered position.
- Q2. Grader availability. Whether two invokable configurations exist that satisfy
  P5 disjointness is unknown to this round. If only one exists, the honest states
  are single-grader with recorded limitation or not-run grading; two graders is the
  ratified design (S8), so this escalates rather than silently degrades.
- Q3. A1's reading of the 15-minute ceiling versus D5. If consensus prefers capping
  packet runs, the experiment records cap-truncations as timeout cells and the gate
  is likely unmeasurable; that trade should be explicit in consensus, not absorbed.
- Q4. zcode-1 catch-up: it must read the signed design and author its own current
  review/signatures per the kickoff (lines 30–32); its bounded readiness evidence
  does not yet exist in the sources read.

## Risks

- R1. Enforcement asserted but not demonstrated — the largest risk; mitigated by G2
  gating all treatment (S10).
- R2. Identity leakage to blind graders through style/protocol wording; reported,
  per S8, never claimed eliminated.
- R3. Fifteen minutes shared across a 12-call full-four cell is tight (~75 s/call
  average including orchestration even with concurrent phases); cells that cannot
  fit become timeout cells by design, which may bias completed cells toward cheaper
  arms — mitigated by task-major ordering and the P7 enum rather than by dropping
  cells.
- R4. Grader/model correlation and small n limit inference, as the original FINAL
  already records (`FINAL.md` lines 257–263); the amendment adds no new inferential
  claims.
- R5. A mid-experiment stop leaves the packet canaries/control not-run, which
  correctly keeps the optimization unshipped; the risk is misreporting this as a
  failure of the packet rather than of the budget — the P7 enum separates
  not-run from refutation.
- R6. Reservation sizing from a small observed cost sample can still overrun on a
  tail-latency call; the 2× factor and never-zero settle are the mitigation, and an
  overrun stops the experiment rather than borrowing (P6.1, P6.5).

---

# SELF-CORRECTION — 2026-09-16 (owner: kimi-1)

Bounded owner self-correction under §15.1. All original bytes above are preserved as an exact prefix; only the named claims are superseded. No peer round was read; no new audit was run. New evidence: the fresh source context supplied to this launch — `provenance.json` (integration_head `2aa2c41a0b9da40dbaeaa6e5b90fda5836f8fabf`, "published integration source; independent recovery candidate not included") with `launch_budget.go`, `telemetry.go`, `record.go`. I own the corrected claims and issue no verdict on them (§15.1); locators are supplied for non-owner verification. Weakenings take effect immediately; any strengthening stands UNVERIFIED. Superseded: A2; A1 and Q3's framing; S10's excerpt-absence claim and S9's "unintegrated" tag; the summary's boundary sentence; ALT-B's tag; P6.2; R6; the 346 total as hard denominator.

## SC-1 — supersedes A2 (ledger scope)

A2 is superseded. The authorization quoted in S4 (`00-prompt.md` lines 20–23) states implementation and review costs already incurred are not this experiment budget. The USD 15 ledger covers only prospective experiment calls: pilot treatment, packet/canary/control, grading/adjudication, and their retries — plus any new experiment-specific calibration or probe, which must be accounted in the freeze file before launch, never absorbed. The S12 implementation calls and any implementation-review calls do not consume the experiment ledger and are not retrospectively charged. My "counting more is the safe direction" double-charged and misstated the user's allocation. This correction aligns the plan with the recorded authorization; no new user approval is needed.

## SC-2 — supersedes A1 and withdraws Q3's two-readings framing (elapsed ceiling)

A1 is superseded. The user's 15-minute shared per-task-per-arm ceiling binds the packet experiment too: each AB/BA pair member, each canary, and the full control run under the same elapsed ceiling as pilot cells. The owner has not authorized uncapped packet runs, and A1 could not waive the ceiling by assumption. An overrun is a censored observation — a timeout cell, retained, never dropped — neither a pass nor a refutation of measured speed. R is computed only over uncensored pair members; a pair with a censored member is reported non-evaluable, with planned-vs-evaluable denominators, and the ship gate cannot be declared from censored data. The S3 thresholds stand verbatim (ship R ≤ 0.50 in both phases, canaries pass, zero obligation misses, the (0.67, 0.80] dispute travels to the user); what changes is denominator honesty, not the threshold.

## SC-3 — corrects S9/S10 provenance, the summary's boundary sentence, and ALT-B's tag

The fresh provenance identifies the runner/telemetry bytes as the published integration source at HEAD 2aa2c41a…; only the trajectory recovery files were unintegrated candidate context. My "isolated, unintegrated" characterization is superseded and S10's excerpt-absence claim is withdrawn: `LaunchBudget` is defined (`launch_budget.go` lines 19–23), with `reserveBudget` (64–151) and `settleBudget` (157–173) invoked from the launch path (`telemetry.go` 99, 146–153), a `budget_refused` class (189–191), a monetary-binding requirement `budget.RequireMonetaryBinding(bound, defaults.MaxCostUSD)` (`launch_budget.go` 111), a persisted operator binding overriding programmatic policy (98–134), fail-closed on unreadable defaults (105–110), reservation retained on settle failure (`telemetry.go` 150–152). Absence from my earlier excerpt was an excerpt artifact, not evidence of absence. Still unshown here: `internal/budget`/`internal/config` internals (`Store.Reserve/Settle`, `RequireMonetaryBinding`, `Limits`); per SC-4 an integrated reservation ledger is not by itself a provider-spend upper bound. The real positive-bound runner recovery path remains separately unfinished. ALT-B is re-tagged PRESENT IN PUBLISHED INTEGRATION SOURCE with that caveat; its adopt-as-shape / reject-as-available-ceiling disposition now rests on SC-4.

## SC-4 — supersedes P6.2 and R6 as launch authorization

Twice-observed-cost reservations do not establish an upper bound on actual provider spending, and post-overrun stopping cannot enforce a hard USD 15 total: a reservation gates whether a launch starts, not what a provider consumes during the call; one tail call can settle above its reservation with no clawback. P6.2's sizing rule and R6's 2× mitigation are withdrawn as sufficient launch authorization; they survive only as ledger discipline under P6.1. Closing the gap requires PRIMARY evidence of an enforced runtime bound — a provider- or CLI-side hard cutoff demonstrably terminating spend at a configured limit — shown by the P6.4 bounded non-treatment demonstration. Discovery notes only, not enforcement evidence: Claude CLI help advertises `--max-budget-usd`; Zcode `--max-turns`; Kimi help shows no matched dollar/token/turn option. Help text is a claim, not a mechanism; none proves runtime enforcement, and a missing flag is not proof of absence. No pricing is invented. Status: the feasibility gate is OPEN — reported with full planned denominators (36 pilot cells, 28 packet/canary/control, 72 grading, ≤ 6 adjudication) and the reservation arithmetic — claiming the provider boundary neither present nor absent from help text or excerpts.

## SC-5 — supersedes the 346-call total as the hard denominator

The ≤ 346 figure (240 pilot + 28 packet + 72 grading + ≤ 6 adjudication) is a proposal for the base plan only; retries are additional unless explicitly reserved within that total. The freeze file must record one reading: (a) hard worst case — the predeclared 1-retry cap yields ≤ 692 launch reservations, all charged against USD 15 before first treatment; or (b) a fixed retry reserve carved inside the envelope so base + reserve ≤ 346. USD 15/346 ≈ USD 0.043 stands as arithmetic-only against the base plan; against (a) it is ≈ USD 0.0217. Neither is a feasibility claim, and no scope is silently reduced to fit. Unchanged: packet thresholds verbatim, grader independence with the ≤ 6-call adjudication cap, no treatment before freeze, historical attribution intact.

## What stands unchanged

P1's narrow amendment scope; P2's allocation structure (now under SC-2's ceiling); P3 full context in every arm; P4 packet design preserved; P5 grading independence; P6.1/P6.3/P6.5 ledger, fit-gate and ordering; P7 honest not-run reporting; P8 freeze; G1–G8 (G2 now reads per SC-4: demonstrated enforced bound or the gate stays open); Q1, Q2, Q4; R1–R5.

## Limitations of this correction

Only my own round, the kickoff, and the fresh three-file excerpt with its provenance were read; `internal/budget`, `internal/config`, CLI help, and provider documentation were not re-audited — statements about them are discovery notes, not verdicts. No self-verification or acceptance verdict is issued; the corrected claims await non-owner review.
