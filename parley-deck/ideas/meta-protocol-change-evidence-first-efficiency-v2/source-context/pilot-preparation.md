# Comparative Pilot: Preparation, Not Preregistration

Status: draft. No treatment calls have been made. Freeze only after the runtime
binary, test suite, task inputs, rubric, model configurations, allocation and
resource policy are hashed. The final freeze file is write-once; amendments get
new versions and invalidate affected comparisons.

## Target Question

For these twelve bounded engineering tasks, what quality, defect detection,
completion rate, elapsed time and observable cost does solo, a pair, or the full
six-participant roster achieve? This is not a vendor IQ ranking and cannot prove
that historical production quality improved causally since April.

Tasks are curated artificial fixtures inspired by observed engineering failure
classes, not a random sample of production work. Calls are real models through
the instrumented CLI. Synthetic tasks must never be confused with mocked model
responses. The four executable, four review and four design tasks below cover
different output shapes. All candidates see exactly the same task bytes.

## Experimental Units

- Twelve tasks, each assigned all three arms, counterbalanced order.
- Full arm: all six active roster IDs frozen from `parley roster show`.
- Solo identity and duo pair rotate; drafter and critic roles rotate separately.
- Fixed phases: independent proposals, scoped cross-review, final draft, explicit
  acceptance or blocking disposition from all active arm participants. Exact
  bounded variant and deviations are recorded, never called an unrestricted run.
- Full protocol context in every arm; packet generation is shadow-only and kept
  out of participant inputs. Packet A/B is a separate experiment.
- Fresh task workspaces contain only participant-visible inputs and owned output
  paths. Hidden acceptance tests and grading keys are outside those workspaces.
- A failed member is not dropped. Failures remain in intention-to-treat counts.

## Ceilings to Freeze

The available CLIs do not all expose an enforceable total token ceiling. Do not
claim equal compute from equal wall time. Freeze the exact enforceable resource
ceilings and separately record output-word guidance, provider caps, observed
tokens and unknown coverage. If identical token ceilings cannot be enforced,
state that limitation explicitly before measurement rather than after results.
The proposed main ceiling is equal per-task arm elapsed time, with conservative
spend reservations under the user-selected overall estimated-cost stop.

An overall stop may leave some arms not-run. Do not replace missing cells with
mock output or selectively omit expensive failed tasks. Analysis must distinguish
not-run, provider failure, timeout, protocol block and valid final output.

## Grading

Two fresh grading-only configurations must be disjoint from task authors and
candidate authors. They are observers, not extra consensus votes. Graders see
anonymized candidate outputs, task bytes, rubric and deterministic test results,
not arm/model/agent/cost/order metadata. A separate mapping retains original
hashes. Style and protocol wording can still leak identity; report this risk.

Primary for code: hidden executable acceptance tests, failure count and critical
invariant violations. Secondary rubric (0-4 each): correctness 30%, completeness
25%, verification 20%, explicit trade-offs 15%, actionability 10%. Review tasks
also report valid material findings, unsupported findings and missed seeded
defects. No reward for more words, more signoffs or larger teams.

Predeclare disagreement requiring adjudication: any critical-invariant conflict,
pass/fail disagreement, or weighted rubric difference above 15 points. An actual
third blind review or an explicit unresolved cell is required. Do not average a
factual contradiction into a consensus score.

Report paired task-level differences and bootstrap uncertainty with fixed seed,
not unpaired averages pretending 36 independent tasks. Cost is CLI/provider
reported estimate unless invoice evidence exists. Missing cost is missing, not
zero; report coverage alongside any quality-per-dollar calculation.

## Follow-Up

Register 14/30-day observations from actual delivery time. Record regression,
rework and adoption only when observed. This pilot does not fabricate elapsed
durability evidence.
