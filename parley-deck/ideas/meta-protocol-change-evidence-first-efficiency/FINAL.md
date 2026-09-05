---
idea: meta-protocol-change-evidence-first-efficiency
status: final
author: codex-1
consensus-date: 2026-09-05
participants: [codex-1, claude-1, hermes-1, kimi-1]
track: deliberation
design-pr: https://github.com/feci/parley-deck-cli/pull/72
---

# Evidence-First Efficiency

## Final plan / specification

All four participants wrote independent round-01, cross-review round-02 and
their own accepting signoff. Claude and Kimi accepted with the recorded
reservations: live experiments and independent verification remain obligations,
not passed gates. The initial Claude signature attempt failed on a session
limit; the post-reset attempt wrote his signature without changing model,
route or quorum. This is design acceptance, not implementation completion.

### D1. Deliver Every Priority, With Independent Gates

The five priorities, the exact ratified packet experiment, the real 12-task
three-arm pilot, and a self-contained offline HTML are all required. A passing
slice cannot hide a missing slice. Design accepted, implemented, independently
verified, deployed and partial are distinct states. Historical files stay frozen.
No package release, global install, immutable-core publication, or production
deployment is authorized or necessary for this implementation.

### D2. One Instrumented Launch Path and Honest Usage

Reserve a unique `invocation_id` before spawn and keep run, segment, idea, phase,
agent, adapter, launch mode, attempt ordinal, and `retry_of`. Requested/start/
terminal records include failed builds and failed starts. Every supported path
is covered: rounds, review, implementation/fixup, steer, consult/goal-check,
signoffs, preflight, interactive handoff, and manual facilitation through a CLI
launch command. A command merely printed for a human is an unobserved handoff,
not a launched or successful process.

Requested model/effort/speed and provider-reported resolved identity are separate.
Record timing, first visible activity, bounded token counters, artifact identity,
failure observations and cost with source/basis/coverage. Unknown values are null.
Never serialize raw prompts, environment, argv or transcripts into telemetry.
Use typed allowlists and secret-shaped metadata rejection/scrubbing. Existing
private process logs are not automatically public report evidence.

Emit one normalized `agent.usage` summary per invocation. Deduplicate by unique
identity, not a repeated integer ordinal; retries remain separate spent attempts.
ACP `used`/`size` is context utilization and stays diagnostic-only, not billed
tokens or cost. A provider cost field is not an invoice unless that provenance
is actually established. With a configured monetary ceiling, unknown spend
requires a visible stop or a documented conservative reservation, never zero.
Telemetry writes that are required for a spend/launch boundary fail closed.

### D3. Completion Requires Current-Tree Independent Evidence

Extend existing checks/review machinery with typed per-criterion records and
scope reconciliation. Bind to a reviewed commit plus tested code-tree digest;
exclude only defined evidence artifacts to avoid evidence-commit invalidation.
Track pass/fail/skipped/not-run, executor/verifier and provenance, command digest,
executed-case counts when the test format supports them, timings and bounded
scrubbed diagnostics. An evidence-write failure is a failure, not a warning plus
PASS. A self verdict, stale code-tree, skipped/no-execution report, missing
criterion, or partial original scope cannot close the whole implementation.

Unknown command output is not semantically certified by a text regex. Material
claims need independent PRIMARY verification, including a concrete serial versus
barrier concurrency fixture. Caller-supplied attribution is not authentication:
state the trust boundary and do not market metadata as proof of a human identity.
No deployment claim without separate evidence. Register 14/30-day follow-ups;
they remain not-yet-due until elapsed and actually observed.

### D4. Single Live-Source Packet Renderer, Default Shadow

Implement `internal/protocolpacket` and `parley protocol packet` / packet check.
Use verbatim blocks, stable heading/block locators, source and packet hashes,
complete disjoint included/omitted index and explicit triggers. Applicability
map changes are protocol changes. Unknown applicability or dependency, malformed
source/map, unexpected phase/track/flag/hash causes visible full-fallback, never
an optimized packet missing unknown obligations. Missing/unprovable authority
blocks rather than substituting a bundled snapshot. Detected secrets refuse
external context emission instead of redacting under a false original hash.

For source-role decks, the live source file is authoritative; expected global
core drift cannot substitute old core text. For consumer decks, resolve and
verify installed core/lock/overlay authority before optimization; unresolved
drift cannot silently pass as a valid optimized packet. Generated bodies stay
in an explicitly ignored ephemeral path. Launch attestation carries
`context_mode`, `source_sha256`, `packet_sha256`, `fallback_reason`.

Change all three instruction sources together: skill standing instructions,
source protocol session-start wording, and runner/handoff prompt wiring. Every
builder uses the shared renderer, not an independent protocol selection logic.
Default remains full context with shadow packet audit until the ratified gate
passes and rollout is authorized. The twelve-task pilot uses full context in
every arm so packet treatment is not a confound.

### D5. Preserve the Ratified Packet Experiment Exactly

The prior `meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md` is binding,
including its never-cut rules, full Section 15 phases, omission audit, and exact
experiment: phases 1 and 6, six matched AB/BA pairs EACH, same model, effort,
task, output ceiling and snapshot, packet generation inside elapsed time, three
packet canaries plus one full control exercising auto_implement and Section 14,
with Section 6/14/15 obligations checked on every run. Freeze inputs before calls.

Ship gate is R = median(packet/full) <= 0.50 in BOTH phases, with zero correctness
misses and passing canaries. Retain all other prior decision bands; the disputed
(0.67,0.80] refutation interval goes to the user with both existing positions.
No post-hoc tuning, synthetic substitute or whole-idea speed claim. A failed gate
is a valid measured result and keeps optimization unshipped; an unrun trial is
unfinished work. Non-implementer recomputation is required.

### D6. Charge Budgets Before Work Across All Entrypoints

Preserve the existing monotonic charged-fixup count and maximum with driver
markers, plus HardCrossReviewCap. Do not derive safety counts from participant
headings or self-authored grants. Add a shared read/reserve boundary for manual,
driver, resume and BLOCK backedges with durable idempotent action identities.
Unknown/corrupt charged state fails closed. A failed charged attempt stays spent.
Signoffs, zero-fix reviews and status queries do not spend a fixup cycle; reaching
the inclusive cap does not prevent verification of its final allowed attempt.

Exhaustion escalates with durable evidence and never closes. Human extensions
must use an actual operator control, not frontmatter impersonation. Audit call
sites and preserve existing track gates; no new Section 4.0 cell edits here.
The opt-in pilot trajectory rule escalates after two consecutive independently
confirmed material patch-induced regressions. Repeated unchanged criticism is
not a new regression. The escalation does not auto-complete or auto-drop quorum.

### D7. Liveness Is Observation, Not Diagnosis From Silence

Distinguish ready, malformed-reply, process-exited-empty, deadline-no-output,
deadline-after-output, classified provider failure and process failure. Only an
exact PONG assistant response extracted from a recognized envelope is ready;
an echoed instruction containing PONG is not. Unknown/empty/parser/deadline
outcomes open a nonzero readiness-resolution gate, never a silent pass or
automatic exclusion. Existing explicit operator exclusion for definite
unavailability can remain distinct. Provider failures stay blocking.

Honor declared `Spec.BuffersStdout` in soft watchdog policy, retaining hard
timeouts and process-group cleanup. ACP agent activity counts; our own heartbeat
does not. Do not search HOME for misplaced artifacts. Diagnostics name the
explicit cwd and expected absolute artifact path. A known-hung fixture verifies
cleanup, not universal diagnosis. Historical events lacking decisive evidence
stay inconclusive; do not retrospectively label them proven false kills.

### D8. Live Comparative Pilot and Presentation

Freeze twelve tasks, three arms (solo, duo, full six active roster IDs), equal
per-task arm resource ceilings, seed/order, roles, versions, rubric and tests
before treatment calls. Counterbalance identities, pairs, drafters and critics.
Use bounded actual independent model proposals, cross-review and owned final
outputs; document the exact experimental protocol variant and conformance.
Retain failures and timeouts in intention-to-treat denominators. No substitution
of mocks or a smaller roster for the full-roster arm.

Use two blind evaluators who did not author the tasks or candidate outputs;
grader-only model configurations may be observers, never extra quorum votes.
Remove treatment/author metadata, preserve answer substance, and adjudicate
disagreement or explicitly report it unresolved. Deterministic acceptance tests
are primary for executable tasks. Report paired deltas, uncertainty, quality,
valid and false findings, regressions, actual elapsed time, coverage and cost
provenance. Small n and model/task confounding limit generalization. No universal
vendor ranking or causal historical trend from this pilot.

The offline HTML embeds baseline, current implementation/experiment results,
denominators, evidence and limitations. It supports keyboard access, filtering,
desktop/mobile/low-height and print; charts are nonblank and legible. Verify with
ego-browser only. A pending result is visually distinct from passed evidence.


### Packet Never-Cut and Decision-Band Invariants

Preserve the authoritative protocol; applicable modals, negations, conditions
and exceptions; Section 4.0 overrides/invariants; current-phase block; non-solo
and files-canonical close guards; Section 6 rule 3, status re-read, English-only
and no-secrets; escalation; Section 14; active transport current-phase mechanics;
applicable close/cap/strict gates; full Section 15 in phases 1,2,3,6,7 and
Sections 15.1-15.4 plus 15.7 in phases 5,8; Section 7 throughout a protocol-change
idea; and on-demand raw historical artifacts. Omission classification cannot
prove its own semantic completeness.

The earlier packet FINAL's disputed speed band remains unresolved. Codex and
Kimi refute at R>0.80 in either phase; Hermes refutes at R>0.67. In (0.67,0.80],
return the measured value and both positions to the user. The implementer cannot
pick one. Above the ship threshold but outside that dispute, preserve the
earlier middle-band reporting rules; do not round the result up to an estimate.
Per-call saving never becomes an unsupported whole-idea saving.

## Purpose / user-visible outcome

Make Parley capable of reporting what model work was actually attempted, what
was independently verified, and what the time/cost/quality trade-off was.
Deliver all priorities and an offline HTML that distinguishes historical
descriptions, tested implementation, measured experiments and future follow-up.
Do not turn more words, more signatures or larger teams into a quality score.

## Context & orientation

The CLI repository is parley-deck-cli, with Go packages under internal/.
The sibling parley-deck-skill source supplies installed skill instructions.
The live source-role protocol is parley-deck/COOPERATION.md. The evaluation
workspace is the sibling parley-deck-evaluation directory; its historical
assessments/2026-09-05 snapshot is frozen, and new work belongs under
delivery/2026-09-05 plus its evaluation scripts/templates.

Code base for design inspection: 257ef8c75ac5e478edf42b59543e89cf94730de6.
Reuse internal/runner, internal/store, internal/driver and the existing
checks/preflight paths. The baseline full Go suite passed, but no proposed
feature is certified by that baseline result.

Codex owns shared runtime/app integration, telemetry, budgets, experiments and
presentation. Claude owns the new packet package/helper/map and source instruction
changes, including skill source in a separate worktree. Hermes owns liveness and
supervision. Kimi owns typed evidence and the checks integration. Exact disjoint
paths, tests, integration order and recovery state belong in IMPLEMENTATION.md
before code edits. Non-owners review every changed slice; no self-issued overall
acceptance is sufficient.

## Observable acceptance criteria

| ID | Required observable result |
| --- | --- |
| AC-T1 | Every supported launch surface yields a unique requested record before spawn and a terminal record after its real outcome; unobserved handoffs are not launches. |
| AC-T2 | At least 20 real new-runtime attempts reconcile without losing failed starts/retries or duplicating usage; requested and reported identity remain distinct. |
| AC-T3 | Structured usage parsers, unknown/null provenance, secret canaries and telemetry-write failure tests pass; monetary boundaries stop or conservatively reserve unknown spend. |
| AC-E1 | Missing/self/stale/skipped/zero-case/partial evidence and evidence-write failure cannot close a new whole implementation; code changes invalidate the tested tree. |
| AC-E2 | An independent verifier executes the concurrency counterexample and corrected case; current-tree criterion evidence covers the delivered scope, with deployment separate. |
| AC-P1 | One live-source renderer supplies complete inclusion/omission records and attestation across CLI, runner/handoff, source protocol and skill instructions; missing authority blocks and unknown applicability falls back visibly. |
| AC-P2 | Execute the exact 12 matched packet pairs and 3 canaries plus full control; a non-implementer recomputes both R values and checks every seeded obligation. Report a failed speed hypothesis honestly. |
| AC-B1 | Driver/manual/resume/BLOCK/failed-attempt boundary fixtures preserve precharged monotonic counts, inclusive 5/6 and 3/4 limits and idempotent actions; exhaustion cannot close or silently reset. |
| AC-B2 | The opt-in two-confirmed-patch-regression trigger escalates only with independent patch-linked evidence; unchanged repeated criticism does not trigger it. |
| AC-L1 | Real fake-child fixtures distinguish quiet late success, hard timeout, partial output, malformed/echo/empty replies, auth/provider/process errors; cleanup and unchanged quorum are verified. |
| AC-X1 | Freeze and run twelve real tasks in solo, duo and full-six arms with equal enforceable ceilings, controlled full context, rotation and owned phase artifacts; retain failures and not-run cells honestly. |
| AC-X2 | Two blind nonauthor grader configurations, deterministic code tests, disagreement handling, paired uncertainty and cost coverage produce auditable results without model-IQ or historical causal claims. |
| AC-H1 | One self-contained offline HTML contains baseline plus actual delivery/results/limits; ego-browser confirms desktop/mobile/low-height/print, nonblank charts, readable controls, filtering and keyboard operation. |
| AC-F1 | Register 14/30-day follow-up from actual delivery, marked not-yet-due until observed; no fake elapsed durability evidence. |

## Idempotence & recovery

Preserve failed invocation records, canonical artifact ownership and completed
attempts. Resume from IMPLEMENTATION.md, current branches and durable run state,
not remembered session narratives. Re-read source hashes before using packets,
and tested-tree hashes before closure. Unique action/invocation IDs prevent
double spending or double counted usage while retries remain distinct attempts.

Freeze experiment inputs and ceilings before treatment; amendments are new,
explicit versions, never silent edits to a completed trial. Keep failed runs in
their assigned cells and do not rerun a negative speed result until it passes.
An external provider or spend stop leaves work pending. Recovery never silently
substitutes a model, drops a selected participant or turns unknown evidence into
success. Do not overwrite the historical assessment.

## Known risks / de-risking

Per-provider billing coverage and resolved-model identity may be unavailable.
Equal elapsed ceilings do not by themselves mean equal compute; report which
resource ceilings are actually enforced. The live pilot is small and curated;
grader/model correlation, task mix and identity leakage constrain inference.
Anonymization preserves substance and records possible leaks. Metadata attribution
is not cryptographic identity, and bounded checks are not universal semantic proof.

Real implementation may reveal interface details but may not silently weaken the
accepted scope or experiment. Keep full/shadow default until the exact ship gate
and any rollout authorization are satisfied. Source-only implementation does not
authorize global core publication or package deployment.

Role concentration: codex-1 is idea author, facilitator, participant and drafter;
this creates no extra signoff weight or dispute-adjudication authority.

## Drafter position changes

None since the drafter's most recent canonical position,
`round-02/codex-1.md`. The earlier title-cased synopsis records the round-01 to
round-02 refinements, not an additional position change after round 2. For
traceability, that file already states: "Accept the single renderer, explicit
applicability map, source-role distinction," and "A manual launch THROUGH the
shared runner is measured telemetry, not honor-system reporting". Both remain
unchanged. Its narrowing of ACP billing is also unchanged: "Those values describe
context utilization, not a monetary total."

## Alternatives disposition

- ALT-01 ADOPT the existing charged cursor plus maximum with driver markers;
  REJECT participant-heading counts and self-authored grants. The charge precedes
  the attempted work and survives a failed publication.
- ALT-02 ADOPT `Spec.BuffersStdout` for soft-watchdog policy; REJECT silence as a
  universal real-hang diagnosis. Hard deadlines and process-group cleanup remain.
- ALT-03 ADOPT exact semantic `isExactPONG` after recognized-envelope parsing;
  REJECT echoed instruction matching and exit-zero unknown readiness.
- ALT-04 ADOPT extensions to existing checks/review machinery; REJECT a parallel
  self-issued pass ledger. Criterion scope and fresh tree identity bind closure.
- ALT-05 ADOPT refusal of detected secret-bearing protocol disclosure; REJECT
  redaction presented under the unchanged original source hash.
- ALT-06 ADOPT full-fallback for unknown applicability/dependencies; REJECT an
  optimized packet whose omission index cannot account for every source block.
- ALT-07 ADOPT the live source file on source-role decks; REJECT stale global or
  bundled substitution. Consumer authority must be verified before optimization.
- ALT-08 ADOPT one normalized `agent.usage` per unique invocation; REJECT dual
  alias summation or billing inference from ACP context-window utilization.

## Correlated agreement and falsifiers

Four-way convergence is a shared prior, not independent outcome evidence. These
nominally independent proposals belong to the same evidence-first, existing-runtime
family; they are not four experimentally compared architectures. The agreed
position is wrong in a scoped way if negative fixtures still permit false closure,
launch paths escape accounting, unknown obligations disappear, declared buffered
success is killed by a soft guard, or equal-ceiling live comparisons do not support
the claimed efficiency benefit. Such observations must be reported rather than
outvoted. An observed speed-gate failure is valid evidence against shipping the
optimization, not a reason to rerun until success.


## References

- Canonical consensus and owned signatures: ./consensus.md
- Independent proposals: ./round-01/
- Cross-review and retractions: ./round-02/
- Binding prior experiment: ../meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md
- Live protocol: ../../COOPERATION.md
- Design PR: https://github.com/feci/parley-deck-cli/pull/72
- Frozen assessment: sibling parley-deck-evaluation/assessments/2026-09-05/REPORT.md
