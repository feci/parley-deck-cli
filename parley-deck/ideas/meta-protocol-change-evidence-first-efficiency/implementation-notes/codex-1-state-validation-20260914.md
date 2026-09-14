---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
status: integrated-validated-awaiting-independent-review
baseline-commit: 29361d6fdadd511330c081b9b840a93f250ef52c
---

# N2 response: release the cycle guard during full historical validation

The common full-history validator now snapshots complete structural authority,
releases the cycle guard for source/resolution reads, then reacquires it and
compares store identity, complete policy, complete ledger and complete state.
The original caller runs under the second guard only after equality. Live
Finish and registered-criterion Stop can acquire their control guard while
that historical reader is paused. The already decisive pending-history
BeforeCycle refusal now precedes content reads; positive reservation keeps all
checks. This is an implementer response to Claude's N2, not a withdrawal of the
finding or independent acceptance.

## Replay correction and bounded read retry

The first no-retry prototype passed its focused cases and two ordinary CLI
concurrency tests. A deterministic native-only hook then paused an actual public
ReconcileUnchanged after full validation, allowed a second identical public apply
to finish, and released the first. The first failed before its exact-replay
branch (2.541s), refuting that prototype's compatibility for this interleaving.
Original source/test/failed logs remain preserved.

The integrated response allows exactly one complete read-validation retry on
final authority drift before any caller callback. A second drift refuses.
Disappearing authority, full evidence errors, context errors and caller callback
errors are not retried. Every retry reads fresh authority and all referenced
historical evidence; it never caches permission or repeats a callback, charge,
ticket, launch or external action. The identical instrumented public test passes
(2.499s); the instrumentation is not production code.

## Focused evidence

Six new top-level tests cover actual Finish during a paused reader; policy,
ledger and disappearing authority; full checks and the guarded final callback;
early pending refusal while positive reservation still requires an archive;
bounded drift/fresh evidence/error handling; and actual attributed criterion
process stop during a paused reader. Five cases together passed in 12.133s;
the registered stop case passed in 5.561s.

Twelve predicate controls fail at the intended assertions: original locked
Finish; omitted state/policy/ledger/existence/source/checker predicates; late
pending refusal; missing callback guard; an extra read retry; retried callback
error; and retried checker error. The old locked-source registered-stop case
also fails as intended (6.476s). Two control scaffold compile failures remain
preserved (unused originalLedger 0.206s; unused changed 0.179s), and are excluded
from successful counterexample counts. Only the affected controls were corrected.
All native checks preserve the 424-file published Go/module source manifest.

Evidence: .parley-runtime/state-validation-prototype-20260914/ (original and
revised native pointers, exact overlays, JSONL terminals, predicate results).
The forced public test has the same bytes in the failing and passing runs.

## Limits and remaining work

This removes the full shared validator's historical-content hold, not all work
under the guard. Positive cycle reservation and individual callbacks may still
perform guarded reads or publication. The test's two-second control context is
a deterministic small-fixture witness, not a production deadline guarantee or
real 128-attempt / 256 MiB concurrency measurement. The retained F6 historical
fixture and original timeout results are unchanged. Full source validation is
not a same-UID writer fence and does not prove descendant inactivity.

There is no persisted schema, protocol, quorum, pilot or helper retry change.
N1/N3 remain awaiting independent review. Historical/abnormal terminal recovery,
consumed/incomplete helper authority, custody, full current-scope participant
review and all live-experiment/delivery obligations remain open. PR #73 stays
draft and the whole audit stays in progress.

## Frozen validation

428-file Go/module manifest: 6149bf1db913498f6b4495eaac24853facd5df0e7c1bb91c37f1e66020ff5995.
Full unmodified suite PASS 405.883s (31 package passes, CLI no-test-files skip);
six-package race PASS 504.857s; vet PASS. All six new top-level cases appear in
both runs. Compiled shared trajectory PASS 60.196s and app PASS 14.357s,
including actual registered stop and concurrent production CLI applies.
Windows CLI/trajectory/app cross-builds PASS; PE amd64 checked, runtime unverified.
Exact native/shared source, control and log hashes and protected participant /
historical HTML bytes verified by the executable final verifier.

The initial parallel full suite FAILED 441.648s at the existing
TestAgentsExecRecordsManualLaunch (6.28s, generic invocation failure); every
other test passed. That case passed in the contemporaneous race, a native-only
isolated diagnostic run (6.638s, actual 371ms exit0 with artifact hash), and the
unmodified full rerun without simultaneous race load. The original failure's
cause remains unresolved; it is preserved and is not relabelled or claimed fixed.
No source, test assertion or timeout changed between the complete runs.
The accepted race was not rerun.

Evidence: .parley-runtime/state-validation-final-validation-20260914/.
Final verification reports the initial failure separately. This is an
implemented response awaiting independent review; the full audit is incomplete.
