---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
source-manifest-sha256: 83bba004bdb2392f825f6d5752ec6c5de6b1780c1c247a2a5ab3f27e3e651ebc
---

# Resource guard continuity disposition

This is a Codex-owned implementation prerequisite for safe lock recovery. It is
not an independently accepted migration command, review consensus or completed
full-goal requirement. The original quorum and frozen experiments remain intact.

## Executed counterexamples

An actual child process acquired a synchronization-only resource guard and kept
its descriptor open while the parent removed the origin and local lock pathname.
The uncorrected parent then acquired a new lock, overlapping the still-live
holder. Unlike a charged ledger, a resource-only directory had no durable witness
that distinguished prior use from first use. The first test also cancelled its
child context before registered cleanup, producing a killed-holder diagnostic;
cleanup ordering was corrected separately. The observed second acquisition had
already occurred before that cleanup failure.

Six further executed counterexamples changed or removed the shared origin after
its initial validation, at either the main kernel acquisition or its exclusion
probe. The unchanged local token/inode still allowed permission. All failures are
retained in .parley-runtime/guard-continuity-negative-20260911.log.

## Implemented behavior

AcquireResourceGuard publishes guard-established with exact origin bytes under
the held, verified kernel lock before returning the first permission. A retained
witness with no origin blocks both origin and local-inode recreation. An existing
witness must match the current origin. Established old guards acquire this witness
on first successful use by the new runtime; this is not retroactive protection
for an already erased old guard.

All lock users recheck the exact pinned host/path/token origin after acquisition
and after the exclusion probe. The continuity callback also finishes inside the
held lock and is followed by descriptor, pinned-origin and cancellation checks.
These are read-only checks; an old waiting caller never republishes a missing
origin. Publication failures return no permission and release kernel ownership.
Exact unchanged-origin retry can finish interrupted witness publication. No new
ledger, charge, policy, model invocation or authorization is created by the guard.

Tests cover actual live-process overlap refusal, origin removal/token/host drift,
released kernel ownership, witness publication errors before/after persistence,
exact retry, conflicting witnesses, origin drift during publication, cancellation
and absence of staging/ledger creation. Existing budget and evidence integration
suites exercise the shared production API.

## Limits and remaining work

This continuity witness does not establish quiescence or implement lock migration.
Do not delete it, the origin or an established local inode to recover a lock.
Old binaries and complete erasure/consistent forgery of all records are outside
this cooperative runtime boundary; no malicious same-UID or distributed-writer
protection is claimed. A previously held guard is not revoked by external edits;
safe operator maintenance must still stop its users. Windows runtime is unverified.

Safe origin migration, changed-inactive-import reconciliation, durable semantic
action replay, independently confirmed two-patch regressions, fresh independent
acceptance, real-model launch/concurrency/closure coverage, exact packet trials,
12-task solo/duo/full-six pilot, frozen equal ceilings/role rotation, blind nonauthor
grading, participant-owned final reviews/signatures, final populated report QA and
actual elapsed 14/30-day follow-ups remain binding. Historical quorum/pilot/funding
and Claude recovery choices remain unanswered. The previously checked table report
is explicitly bound to source 60640709; it does not certify this newer guard change.

## Current-source validation

All 345 Go/module files match manifest
83bba004bdb2392f825f6d5752ec6c5de6b1780c1c247a2a5ab3f27e3e651ebc.
Full Go suite PASS (132.785s). Budget/evidence/app/driver/runner race PASS
(152.416s). Scoped vet PASS (1.251s). Windows amd64 app cross-build PASS;
Windows runtime remains unverified. The compiled budget binary passed isolated
shared-volume resource-guard, origin-change, lock and concurrent-reservation
fixtures (2.063s). Shared fixtures ran outside any real repository ancestry.
Focused current guard fixtures passed (0.641s). Logs and checksums remain under
.parley-runtime/guard-continuity-validation-20260911/ and the referenced sibling
negative/focused logs. No new model or real operator accounting action ran.
Fresh independent acceptance remains pending.
