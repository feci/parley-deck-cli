---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-12
status: automated-validation-passed-independent-review-pending
---

# Lost parent-result recovery

A successful independent verifier and its complete helper journal can survive a
failed publication of the enclosing parent result. The consumed ticket prevents
repeating the helper, while ordinary reconciliation requires the parent file.
`trajectory recover-parent` now derives a separate retained observation from the
original execution. It does not repeat the work or rewrite its failed publication.

## Implementation and scope

- Reconciliation derives from the private request, original charged source
  archives, exact quorum/criterion commands, shared request and launch, matching
  requested/started/successful terminal lifecycle, complete helper receipt, step
  and process records. It recomputes `AssessCaptured` instead of trusting a
  caller-supplied invocation or parent assessment.
- `parent-recovery.json` pins the original parent observation, source state,
  primary derivation and reconstructed parent digest. The parent stays missing,
  an empty failed-publication directory, or its original bounded partial/compatible
  publication-failure bytes. A complete success uses ordinary reconciliation.
- Preview is read-only. Apply takes its exact SHA and publishes under the cycle
  guard through a synced staging file, rename and file/directory sync. Identical
  concurrent applies and exact retries keep one record. A retry after a failed
  persistence barrier repeats the sync without rewriting the record.
- Existing source/charge/archive validation stays unconditional. Targeted recovery
  of a lost, previously reconciled parent reproduces all pinned facts and either
  existing canonical newline encoding. Other resolutions are checked normally.
  New reconciliations additionally pin recovery record bytes. Historical state
  itself is never retroactively rewritten to add a recovery hash.
- Exact receipt replay survives reconciliation and later legitimate work while
  original facts remain valid. First apply refuses state/evidence drift. Changed
  original parent state, lost/mutated accepted recovery or changed primary facts
  refuses. Recovering evidence does not issue reconciliation, new execution,
  budget, continuation acknowledgment, review clearance or final acceptance.
- No attendance gate is added to deterministic recovery. Existing attended
  continuation is unchanged. Content hashes do not authenticate a same-UID actor,
  model identity, human identity or an entire source beyond its captured scope.
  State version 3 has an additive recovery hash; older strict readers refuse it.

## Counterexample and correction

Initial partial-parent parsing decoded the whole JSON into a struct. A syntax
error leaves zero values and can hide already-written contrary attribution or a
non-publication failure. An explicit false trajectory-pending flag was also
indistinguishable from omission. Recovery now uses the standard JSON decoder to
inspect complete fields in order, rejects repeated/unknown fields and explicit
contradictions, and permits only compatible truncated keys/scalar prefixes.
Partial structured assessments and ambiguous encodings refuse for inspection.
Focused tests also exposed incomplete-object and trailing-comma boundary handling;
those failures and corrected output are retained, rather than reported as passes.

## Concurrent initial budget publication correction

The first full suite on the 396-file source failed in the pre-existing
`TestConcurrentReservationsCannotOverspend` (274.123s overall): one reservation
reported a ledger without a lock origin even though a concurrent first writer
had published both. The lock reader had retained an earlier absence observation
and then inspected newly created history. This was not an overspend or a failure
of the new parent-recovery tests, and the failed suite remains a failed result.

After claiming the budget files, the correction refreshes an absent-origin
observation when history has appeared. The newly found origin must still pass
normal location, identity and held-kernel-lock checks. A genuinely missing origin
or identity and a conflicting origin still refuse without recreation. A new
boundary-controlled concurrent test covers all four cases with an actual first
reservation; existing overspend, lost-resource-origin and pinned-migration-waiter
checks pass. An overlay removing the fresh read must reproduce the original
misclassification. Final validation uses the fresh 397-file manifest below.
The failed full run is retained separately under
`.parley-runtime/parent-recovery-final-validation-20260912/`.

## Validation

Final automated validation passed; independent participant acceptance is pending.
Focused PASS 124.831s; full PASS 265.057s (all 32 package terminals);
six-package race PASS 309.809s; vet PASS. Windows trajectory/budget/app/runner/
driver cross-builds PASS; runtime is unverified. Compiled shared-volume
trajectory/app/budget selections PASS 26.816s / 113.069s / 1.010s.
The final verifier accounts for all eleven new top-level test functions (ten
substantive scenarios plus the separate-process harness), six negative controls,
native/shared log equality, stable source, PE amd64 headers and unchanged
historical HTML. Actual production CLI preview/concurrent apply runs inside the
app tests. No earlier-source pass substitutes for this final-source validation.
Source manifest: 397 Go/module files, SHA256
`93dbaa9a63ed1500b72a71ba6eb154e9859c2e7153b4f19508106dd29077f6cd`.
Logs and scripts: `.parley-runtime/parent-recovery-final-validation-v2-20260912/`.
Earlier counterexample/boundary logs:
`.parley-runtime/parent-recovery-validation-20260912/`.

Tests execute real local synthetic helper and criterion processes, including the
compiled production CLI; they are not real-model experiment observations. The
new scenarios cover missing/partial/empty-directory/publication-failure parents,
no repeated execution or accounting, output failure, concurrent CLI apply, primary
evidence drift, failed helpers, original digest preservation, stale preview,
later charged history, changed original scope and lost/mutated accepted provenance.
Package-private publication seams inject failure before publication, at actual
rename, after rename before its final sync, and after successful publication.

## Remaining whole-goal obligations

Incomplete/consumed helper-ticket recovery, orphan reservations, unchanged-source
attempts and other workflow effects remain open. The next mechanism must preserve
original attempts and evidence, retain process reconciliation and charge every
actual retry. Missing terminal/receipt or a dead leader alone is not proof that
all descendants are inactive.

Fresh participant-owned current-source acceptance, complete live launch and
concurrency/closure evidence, exact phase packet experiment, twelve-task
solo/duo/full-six comparison with frozen ceilings/role rotation/blind grading,
final owned signatures, populated offline HTML with ego-browser QA and actual-
delivery-based follow-ups remain open. Historical quorum/default roster and
resource questions retain their existing status. No participant artifacts were
written by proxy, no live model was invoked, and no historical report was changed.
