---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
base-commit: 6962f2af9b6e46548588c5152285e6808f3d6f7f
kind: implementation-candidate (code written; tests NOT executed in this launch)
extends: zcode-1-verifier-refusal-recovery-correction-20260916.md (and candidate before it; both left byte-unchanged)
---

# Verifier parent-lineage recovery candidate (zcode-1)

## Scope and ownership

Owned files written this launch: `internal/trajectory/verification.go`
(authority-reader chronology check + effective-launch resolver),
`internal/trajectory/verification_recovery.go` (one signature extension),
`internal/trajectory/reconcile.go` (derivation wiring),
`internal/trajectory/verification_parent_recovery_test.go` (new),
`internal/trajectory/verification_recovery_test.go` (one stale doc comment
updated on `TestFreshResolutionValidatesRecoveredLineageAndFailsClosed`; its
assertions are unchanged), and this note. The coordinator's mechanical fixes —
the runner-test RunMeasured changed-source prerequisite, the ReadCapturedVerification
arity fix, and gofmt — are preserved untouched, as instructed. No Bash, no
subagents, no model calls, no Git mutations, no fabricated results.

Provenance: PRIMARY (source) — my own reads of `verification.go`,
`verification_recovery.go`, `reconcile.go`, `parent_recovery.go`, `state.go`,
`continuation.go`, `captured_test.go`, `verification_test.go`,
`parent_recovery_test.go`, `snapshot_test.go`, `telemetry/record.go` and the
retained coordinator logs under `.parley-runtime/input/current-checks/` in this
checkout. Coordinator-quoted results below are coordinator execution, not mine.

## What was implemented (the resolution seam)

**The seam.** `deriveParentEvidence` (reconcile.go) previously read
`launch.json` and required THAT invocation's successful observed terminal.
After a verifier-launch recovery, `launch.json` still pins the refused
invocation (retained byte-identical by design), so fresh resolution failed
closed — the open obligation both earlier notes recorded. This candidate wires
that final hop without weakening any check:

- `resolveVerificationEffectiveLaunch` (verification.go, new) resolves which
  invocation's observed lifecycle can complete the verification: the original
  reservation when no recovery exists, or the recovery's bound replacement
  invocation afterwards. Its discovery read of `recovery.json` only SELECTS the
  invocation; every validity decision — the recovery binding against the actual
  `launch.json` bytes, full revalidation of the refused evidence, bound-digest
  equality, and chronology — stays in the corrected common authority reader
  `readVerificationLaunch`, re-run for the selected invocation. No weaker
  parallel parser was introduced. A malformed or contradictory retained
  recovery fails this resolution exactly as it fails every other recovered
  handle.
- `deriveParentEvidence` now follows that resolution: the effective invocation
  must additionally satisfy `runtimeID` (it names an invocations directory),
  `derived.LaunchSHA256` becomes the effective authority SHA (the recovery
  artifact itself when recovered — the same SHA the helper receipt binds), and
  the successful observed terminal/requested/started lifecycle must belong to
  the EFFECTIVE invocation with all pre-existing checks intact (metadata
  binding to run/idea/verifier/phase/headless, process-exited exit 0, PID,
  requested/terminal/started agreement).
- **Recovered chronology, defined explicitly.** The original
  requested/reservation/refused-terminal consistency is unchanged
  (`refusedVerificationEvidence`); NEW: `readVerificationLaunch` now also
  returns-receives the refusal completion time and refuses a retained recovery
  whose `At` predates it ("verification recovery predates the retained refusal
  completion") — recovery cannot predate the original refusal completion. For
  the replacement, the authority-window predicate is uniform: the authority
  reservation timestamp must lie inside `[terminal.RequestedAt,
  started.StartedAt]` — `launch.json`'s own `At` for an ordinary launch, the
  recovery's `At` for a recovered one. A replacement therefore cannot have been
  requested or started before the recovery, and the old reservation's timestamp
  is never applied to the replacement's requested record (it is not in the
  window predicate at all once recovered).
- `refusedVerificationEvidence` (verification_recovery.go) returns the refused
  terminal's `CompletedAt` as a fourth value; both call sites updated
  (admission ignores it; the authority reader enforces it). Admission semantics
  and the recovery artifact schema are unchanged.
- **An ordinary unrecovered journal is byte-compatibly unaffected**: without
  `recovery.json` the resolver returns the original reservation, the window
  predicate uses `launch.At` exactly as before, and every other reader is
  untouched.

**Preview/acceptance binding.** ParentDerivation/preview/acceptance bind the
effective recovery authority and all refused evidence through re-derivation
itself: every `PreviewReconciliation`, `Reconcile` and later state-guarded
operation (`withState` runs `checkResolutions` on every call) re-derives the
complete lineage through the common reader, and the ordinary parent-result
binding requires `parent-result.json` to equal the freshly derived
`ParentResult` — which pins the effective invocation, terminal digest and
receipt digest. I deliberately did NOT overload
`ReconciliationPreview.RecoverySHA256` (that field belongs to the
trajectory-level parent-recovery record, a different mechanism that can
coexist); the verifier-launch recovery binds through `derived.LaunchSHA256`
(the effective authority SHA) and the exact parent-result equality. Mutation
after preview therefore fails the exact recheck rather than being compared
against a stale snapshot.

**Preserved invariants.** Old records/accounting and the original unresolved
state on refusal are untouched (nothing deletes, rewrites or refunds; the
ledger and charge count are asserted byte-identical through recovery, helper
execution and resolution in the new positive test). Same-UID hashes remain
non-authentication (unchanged honesty boundary). A partial helper-only receipt
without the replacement's successful observed requested/started/terminal still
fails resolution at the terminal binding — `TestFreshResolution...` remains a
genuine negative (its stale framing comment was updated; assertions unchanged).

## Tests written (NOT executed here)

New `verification_parent_recovery_test.go`, all through public APIs, reusing
the existing fixtures (`journalFixture`, `dirtyCapturedChild`,
`refusedVerifierLaunchFixture`, `snapshotRead/Write`,
`rewriteRefusedEvidenceCanonical`) and call signatures read from source:

- **Positive** (`TestRecoveredReplacementLifecycleResolvesAndReconciles`): real
  refusal → explicit recovery → fully observed replacement lifecycle (requested/
  started before, successful terminal after — telemetry's exact canonical
  encoding) around a REAL local helper execution under the bound invocation →
  exact parent-result publication → fresh preview validates (fields, parent
  digest, empty trajectory-recovery field) → explicit `Reconcile` accepts,
  replays idempotently, and the retained resolution re-derives on every later
  state guard; original fixup accounting byte-identical, charge count 1,
  refused evidence retained, launch/recovery bytes unchanged, no fabricated
  started lifecycle; the regression outcome legitimately keeps the trajectory
  open (asserted, not hidden).
- **Failure table** (`TestRecoveredParentResolutionRequiresExactLineage`), each
  from a clean fixture whose preview provably resolves first: refused
  requested/terminal deleted, refused evidence canonically rewritten, launch
  deleted, recovery deleted (degrades to the unrecovered reader refusing the
  refusal terminal), recovery rebound to a decoy invocation (plus stale current
  handle refusing), replacement lifecycle moved away, recovery backdated into
  the refusal window (the new chronology check, plus handle refusal),
  replacement lifecycle predating the recovery (backdated timestamps), missing
  replacement terminal, unknown journal record. Every variant also asserts no
  auto-resolution (RequireResolved fails, zero recorded resolutions).
- **Stale preview** (`TestRecoveredParentStalePreviewFailsExactRecheck`): a
  preview accepted on the clean lineage then mutated (refused evidence
  rewritten / recovery deleted / parent-result receipt digest forged) is
  refused by `Reconcile` with the specific error and leaves the trajectory
  unresolved.

Nothing in this launch ran: no `go build`, `go vet`, `gofmt`, or any test. Like
the previous candidates, compile or assertion failures are plausible and must
be fixed before integration; the two prior rounds of coordinator-caught compile
errors (scoped `err`, unused variable, call arity) are exactly why every call
signature above was re-read against source rather than assumed.

## Coordinator evidence since my correction (their execution, quoted; PRIMARY only for my reading of the retained logs)

`.parley-runtime/input/current-checks/zcode-runner-fixture-control/result.json`
records `source_unchanged: true` over the four files, and: the repaired runner
mutation test `TestVerifierRecoveryRefusedEvidenceMutationFailsHandles` PASSED
(3.067s, exit 0), and an exact pre-correction shape-only-reader overlay of the
same test FAILED at the intended stale-stop lineage assertion (3.001s, exit 1)
— demonstrating the digest-equality binding, not just the superseded error,
does the refusing. `zcode-runner-corrected/` retains the earlier trajectory
selection (93 pass events, 110.536s). I claim none of this as my own execution
or as acceptance of this candidate; it predates this note's changes.

## Honest limits and precise follow-ups (NOT done, not designed away)

1. **A real positive runner relaunch riding the recovery is still unbuilt.**
   The runner generates invocation IDs internally (`telemetry.Begin` cannot
   predetermine an ID), so a real `RunMeasured` verifier relaunch today would
   arrive at `ReserveCapturedVerificationLaunch` with a fresh random ID that
   differs from the recovery's bound intended ID and be refused before its
   budget store — by design. Wiring the runner/attended relaunch to launch
   under the bound intended ID (and to produce the replacement lifecycle
   records natively rather than as reconstructed retained bytes, as these tests
   do) is the acknowledged follow-up. The recovery API's caller contract is
   unchanged from admission: it accepts an intended fresh invocation ID chosen
   before that launch exists; this candidate adds no new caller requirement.
2. **Attended app CLI control** (handing the bound invocation to a real
   verifier relaunch) remains separate work, unchanged from the correction
   note.
3. **Same-UID forgery** of a fully self-consistent recovery+evidence+lifecycle
   set remains beyond byte evidence — unchanged honesty boundary.
4. This candidate completes neither the audit slice nor D-series acceptance;
   integration, publishing and signoff are not attempted, and no treatment
   experiment was started.
