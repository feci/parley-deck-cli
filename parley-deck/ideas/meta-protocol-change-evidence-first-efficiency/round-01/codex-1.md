---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
round: 1
date: 2026-09-05
---

# Independent Proposal

## Summary

Implement one auditable improvement program with independently testable slices,
then run the preregistered packet trial and 12-task pilot. Do not turn telemetry
coverage, document completeness, majority agreement, or benchmark score into a
claim of production correctness. Keep the frozen September assessment intact.
This proposal was written before reading peers' round-01 files.

## Existing Alternatives

- PRIMARY: inspected `internal/runner/runner.go`, especially `runExecAttempt`,
  `finalizeExecResult`, and `failEarly`. Exec attempts emit process identity and
  terminal outcomes, but attempt ordinals repeat across segments and neither
  requested/resolved model identity nor usage is recorded there. Instrument the
  shared launch boundary rather than another report-only wrapper.
- PRIMARY: inspected `internal/runner/supervision.go`. Its 120-second first-event
  default measures output visibility, not actual process health. Known buffered
  adapters need an explicit observation contract; disabling all deadlines would
  hide hangs and is not acceptable.
- PRIMARY: inspected `internal/app/driver_checks.go`. `runChecksContract` returns
  true after successful commands even when `writeValidationEvidence` failed.
  Evidence contains names, exit codes, durations and tails but no source-tree
  identity or independent verifier attribution. Fix this in the existing path.
- PRIMARY: inspected `internal/driver/loop.go` and budget test locations. The
  consumer already sums `agent.usage`; a comment explicitly says producers do
  not yet emit it. Missing cost cannot enforce a monetary ceiling reliably.
- PRIMARY: read the previous packet/budget `FINAL.md`. Its exact six AB/BA pairs
  per phase, phases 1 and 6, canaries, and dual R <= 0.50 gate are binding. A
  simpler demonstration or synthetic timing must not replace this experiment.
- SECONDARY: frozen assessment `REPORT.md` and `CASES.md` identify independence,
  false liveness failures, and uncertain historical model identity. Those
  historical proxies do not establish a causal protocol effect.

No new test suite or mutation probe was executed for this independent proposal.

## Proposed Approach

### 1. Launch Evidence

Use a small shared `internal/telemetry` package, with typed records and a
concurrency-safe append-only per-run sink. Reserve a globally unique invocation
ID before attempting spawn. Emit requested/start/terminal records even for
build errors and failed starts; retain existing event names for compatibility.
Do not use PID, agent ID, round number, or ordinal as the sole identity. Retries
get separate IDs and a parent/previous ID.

Fields: schema version, invocation ID, run/segment/idea/phase, adapter and launch
mode, exact requested model/effort/speed, reported resolved identity plus source,
UTC timestamps and monotonic duration, exit/failure class, artifact hash,
first-visible-output latency, token counters, and cost with basis and coverage.
Unknown is null, never inferred from the request. Costs reported by subscription
CLIs are list-price estimates, not invoices. Cache counters and model subcalls
need an explicit aggregation policy; avoid double counting terminal summaries.

Read only documented structured result/usage envelopes or ACP usage updates.
Parse output incrementally and boundedly; never serialize raw prompts, argv,
environment, provider config, stdout, stderr, or secret-bearing URLs in the
telemetry sink. Add real failure-path fixtures and secret canaries. Provide
`parley agent exec` (or the nearest existing launch command) for manual skill
facilitation through the same instrumentation, not a second launch engine.
All raw direct invocations must be marked uninstrumented, not silently counted.

### 2. Completion Evidence

Keep existing review and completion contracts; add a typed evidence record to
the same checks path. Bind it to a code-tree fingerprint excluding the evidence
file itself, with Git commit and dirty-state annotation. Commit-only binding
creates an endless stale-evidence cycle when the evidence is committed.

Each criterion has an ID, explicit status (pass/fail/skipped/not-run), command
digest, independent executor identity, verifier provenance, timing, bounded
scrubbed diagnostic, and relevant output hashes. Exit zero alone cannot prove
that a test actually ran. Supported structured test reports must show executed
cases; unknown command types need explicit independently attested evidence,
not an unreliable universal text regex. Missing evidence, changed code, reused
self-report, no executed cases, and a failed evidence write all fail closed.

Separate lifecycle dimensions: design accepted, implementation present, scope
coverage, independently verified, deployed. An explicitly scoped partial result
does not satisfy the original whole contract. Deployment is optional evidence,
never implied by merged code. Register 14/30-day observations without marking
them observed before their dates.

### 3. Packets and Budgets

Create one renderer adjacent to `protocolcore`, resolving the live authoritative
protocol via existing APIs. Represent top-level sections and normative child
blocks with stable locators and a complete disposition index. Classify every
block; unknown content is included and invalidates optimized mode. Uncertain
hashes, parsing, phase, track, dependency or authority fall back visibly to the
entire live protocol. No stale bundled substitution for an initialized deck.

Return `{sourceSHA256, mode, reasons, included, omitted, text}`. First deployment
is shadow/audit. Wire the CLI command, every prompt builder, skill standing
instructions and session-start instructions together. Keep full Section 15
where the ratified FINAL requires it. Measure render + launch + completion, not
precomputed packet read time. Failed speed gates leave audit mode in place.

Audit budget enforcement across manual commands, initial schedules, consensus
BLOCK backedges, resume, and failed fixups. Reuse trusted driver charges rather
than parsing self-reported iteration counts. The proposed two-consecutive-
material-regression trigger needs confirmed review evidence tying a new defect
to a patch; mere repeated criticism is not a regression. It escalates, never
completes. Keep this trajectory policy opt-in until the new pilot validates it.

### 4. Liveness

Separate installed, runnable, authentication/provider outcome, protocol reply,
and observation timeout. A timed-out ping is inconclusive and cannot remove an
agent from quorum. Preserve exact PONG semantics after parsing only recognized
vendor envelopes; do not accept an echoed instruction containing PONG. Add
quiet/buffered success, late first output, true hang, malformed envelope,
auth/rate-limit/provider failure, and nonzero-exit tests. ACP semantic activity
must not include our own heartbeat. Keep hard ceilings and process-tree cleanup.

### 5. Comparative Pilot

Freeze 12 tasks across bug fixing, requirement interpretation, design tradeoffs,
and review/refutation. Freeze input hashes, hidden deterministic acceptance
tests where meaningful, scoring rubric, adapter/model versions and order seed
before the first treatment call. Run solo, duo, and full active roster with the
same per-task total output-token and wall-time ceilings. Counterbalance solo
identity, duo pairs, drafter, critic, and treatment order. Total budget is an
independent stop rule; never quietly lower one arm's budget after seeing output.

Use real independent proposals, scoped cross-review, and an owned final output;
declare this a bounded pilot protocol variant and measure conformance rather
than pretending full unrestricted eight-phase runs fit a tiny budget. The
full-roster arm retains all six active participants. Failures and timeouts stay
in intention-to-treat denominators. Blind final output IDs and strip roster/
treatment metadata before two independent graders. Use deterministic tests as
primary evidence for code tasks, subjective grades as secondary; adjudicate
disagreements with a third view or an explicit unresolved result.

Report paired per-task deltas, uncertainty, quality per actual elapsed/cost,
valid findings, escaped regressions, rejected false findings, and missing usage.
Small n and model/task confounding prohibit broad vendor superiority claims.
Do not manufacture model identities or usage when tools do not report them.

## Work Boundaries

- Codex: telemetry, launch integration, budget reconciliation, pilot/report.
- Claude: packet renderer/tests, CLI packet helper, proposed skill wiring.
- Hermes: liveness helper/tests and supervision defaults/observation policy.
- Kimi: evidence schema/validator/tests and completion integration helper.
- Exact disjoint file claims go in IMPLEMENTATION.md after consensus. Shared
  runner/app entrypoint changes are serialized by Codex. Another participant
  reviews every implementer's change; nobody approves their own slice alone.

## Concerns and Risks

The scope is material and must not be reduced to a dashboard of planned work.
However, an external spending stop or immutable-core human-ratification gate
cannot be bypassed to claim completion. Preserve runnable code and explicit
pending evidence at every stopping point. Subscription pricing, hidden test
contamination, provider nondeterminism, grader correlated errors, and active
workspace mutations require explicit provenance and limitations. Concurrent
worktrees must not share test-mutated state or global configuration. No release
or global installation is implied by this implementation request.
