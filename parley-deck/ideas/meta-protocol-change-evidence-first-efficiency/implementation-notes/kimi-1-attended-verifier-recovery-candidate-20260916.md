---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: implementation-candidate (code written; tests NOT executed in this launch)
extends: zcode-1-bound-verifier-runner-candidate-20260916.md (SALVAGED slice; all prior notes byte-unchanged)
---

# Attended verifier-launch recovery: app lifecycle + read-only recovery reader (kimi-1)

Status: COMPLETE CANDIDATE — all four tasks implemented; tests written, NONE executed
in this launch (no Bash/toolchain ran; compile or assertion failures remain plausible
and are the coordinator's to surface on the immutable pre-slice copy). This note was
written early and finished by the same agent; nothing here is PASS evidence,
integration, acceptance or signoff.

## Launch attestation and salvage provenance

- Launch: kimi-1, protocol attestation `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`
  (packet_sha256 identical; the shadow audit block is an unapplied diagnostic).
- Method: Read/Grep/Glob/Write/Edit only. No Bash, tests, builds, MCP, web, Git
  mutations, subagents, model calls, installs, or external messages in this launch.
- zcode-1's latest slice is SALVAGED, not a clean terminal success: invocation
  `15b1ca5f-07d1-45bc-9a3e-a3233e2f146b` started; the wrapper terminated with ENOSPC
  while persisting output; the normalized terminal is missing and exit/cost are unknown.
  Its owned code and note exist in this tree and are the base of this slice; I do not
  invent a terminal for it and claim no tests passed for it. Storage was recovered by
  byte-verified relocation of one unrelated synthetic sizing archive; no worktree or
  history was pruned. All prior notes/failures/source corrections are preserved.
- Source remains UNINTEGRATED: this candidate is no acceptance, no signoff, no merge.
  The coordinator tests an immutable pre-slice copy separately.

## What this slice implements (the four authorized tasks)

### Task 1 — genuinely read-only exported recovery reader + runner constructor fix

New `trajectory.ReadCapturedVerificationRecovery(ctx, ticket)` (in
`internal/trajectory/verification_recovery.go`) returns
`CapturedRecoveryIdentity{InvocationID, RefusedInvocationID, RefusedRequestedSHA256,
RefusedTerminalSHA256, PriorLaunchSHA256, RecoverySHA256, At}`. It runs under the
existing guarded complete authority (`withVerification`: charge/state/archive +
quorum + exact retained ticket), re-reads `launch.json`, and resolves through the
EXISTING `resolveVerificationEffectiveLaunch` → `readVerificationLaunch`, so the
refused-byte digests and refusal/recovery chronology are revalidated in full. It
performs no reservation and no publication: no `writeVerificationArtifact`, no
`Mkdir` with create=true, no state write — the callback only reads. Missing launch
reservation or missing recovery is a refusal, never a recreated authority.

`runner.WithCapturedVerificationRecovery` (`internal/runner/trajectory_verification.go`)
drops its unlocked `os.Lstat(launch.json/recovery.json)` probes plus the
`trajectory.ReserveCapturedVerificationLaunch` "probe" — which is a mutator whenever
`recovery.json` is absent (it writes a fresh `launch.json`; if both authority files
disappeared between the unlocked checks and the probe it would recreate a launch
reservation pinning the intended invocation with no recovery and no refused
evidence). The constructor now calls the read-only reader and compares the returned
bound identity: intended == refused → the superseded message; intended != bound →
"differs from its durable reservation"; reader refusal → wrapped refusal. File
presence alone is never authority, and the documented same-UID custody limit is not
reinterpreted as permission to skip missing-evidence refusal.

### Task 2 — attended app CLI preview/apply/relaunch (budget_refused-before-start only)

New `internal/app/trajectory_verifier_recovery.go`, wired in
`internal/app/trajectory.go`:

- `parley trajectory recover-verifier --dir DIR --idea ID --run RUN [--replacement ID]`
  prints an exact `trajectory.CapturedRecoveryPreview` (ticket/run/root/idea,
  original refused invocation + its bound requested/terminal digests + refusal
  completion, and the selected replacement identity — caller-supplied via
  `--replacement` or a fresh `telemetry.NewID`). No process, no write, no budget.
- `... --replacement ID --sha256 SHA --yes` applies: the trajectory layer recomputes
  the preview from live evidence and refuses any drift vs the inspected SHA, then
  performs the single immutable `RecoverCapturedVerificationLaunch` write. Exact
  replay (same replacement + same SHA against the retained recovery) is a no-write
  idempotent return; a conflicting replacement is refused. Recovery grants no
  process, no budget reset/refund, no resolution.
- `parley trajectory relaunch-verifier --dir DIR --idea ID --run RUN` previews the
  relaunch plan (ticket SHA + validated recovery identity + verifier + root/idea/run,
  original failed parent re-validated, live quorum/material scope re-checked via
  `CheckHelperScope`, verifier resolved from configured roster). `--sha256 SHA --yes`
  rechecks the exact plan, then binds `runner.WithCapturedVerificationRecovery` and
  runs the REAL verifier through `runner.RunConsult` exactly like the production
  verify path (same helper command/prompt builders, shared). Budget admission is the
  ordinary immutable operator budget binding (`reserveBudget` untouched — a persisted
  operator binding or the caller's explicit `WithLaunchBudget` policy); recovery is
  not a bypass. Fresh requested/started/terminal come from the real runner under the
  pinned ID. Preview/recovery/replay never launch a model.

### Task 3 — full app lifecycle: immutable recovered parent observation

The original failed `parent-result.json` (written by `verifyTrajectoryWithAgent`
even on a pre-start budget refusal, stage `launch`) is never overwritten, deleted,
or treated as missing, and `compatibleUnpublishedParent` is not relaxed. New
artifact `parent-recovered.json` in the run directory ties BOTH sides:

- the retained original failure: raw-bytes digest + exact decoded `ParentResult`,
  required to be Version 1, same run/request path/request SHA, `InvocationID` ==
  recovery's refused invocation, `TerminalSHA256` == the recovery's bound refused
  terminal digest, empty receipt SHA, nil assessment, `FailureStage == "launch"`,
  `TrajectoryPending`; and
- the fully observed replacement/helper lineage: the validated
  `CapturedRecoveryIdentity` plus the freshly derived `ParentDerivation`
  (replacement terminal/requested/started lifecycle, helper receipt, assessment)
  with `ParentSHA256 = originalParentDigest(derived.Result)`.

`trajectory.PreviewRecoveredParent` / `trajectory.PublishRecoveredParent(ctx, root,
idea, run, expected)` (in `internal/trajectory/parent_recovery.go`, following the
existing parent-recovery stage+rename+sync conventions) re-derive and validate the
complete lineage on every call; a retained record is validated and replayed, never
rewritten. `readParentEvidence` (`internal/trajectory/reconcile.go`) gains a branch:
after the trajectory-level parent-recovery check and before the ordinary
same-bytes parent check, a retained `parent-recovered.json` is validated against
the fresh derivation + the retained original failure + the journal recovery
identity; fresh previews/reconciliation and every later state guard
(`checkResolutions`) re-run that validation, and new resolutions pin the record's
digest in the existing `ReconciliationPreview.RecoverySHA256` provenance slot.
Stale/deleted/contradictory original OR replacement evidence fails closed.
Derivation gains an internal variant returning the effective launch
(`deriveParentEvidenceWithLaunch`); the ordinary unrecovered path is byte-identical.
This is deliberately separate from the missing-parent recovery: a complete failed
parent still makes `PreviewParentRecovery` refuse ("conflicting identity, failure
or assessment"), and unrelated old resolutions are untouched.

### Task 4 — tests (written, NOT executed in this launch; coordinator runs them)

- `internal/trajectory/verification_recovery_read_test.go` (new):
  `TestReadCapturedVerificationRecoveryReturnsValidatedIdentity` (identity equals
  the retained recovery; journal/accounting byte-unchanged; no bound directory);
  `TestReadCapturedVerificationRecoveryRefusesMissingAuthority` (never-reserved,
  never-recovered, and both-authority-files-deleted-after-recovery each refuse with
  NOTHING written — the deterministic form of the old probe's race);
  `TestReadCapturedVerificationRecoveryRefusesDrift` (canonical evidence rewrite,
  refused-terminal deletion, forged prior-launch, backdated recovery At);
  `TestCapturedRecoveryPreviewApplyExactRecheck` (preview binds refused launch +
  replacement; stale SHA and mismatched replacement refuse write-free; exact apply;
  exact write-free replay; conflicting preview/apply refuse);
  `TestCapturedRecoveryPreviewRefusesUnavailableAuthority` (never-reserved,
  ordinary never-refused reservation, invalid replacement identities).
- `internal/trajectory/verification_parent_recovery_test.go` (extended):
  `TestRecoveredParentObservationResolvesAndReconciles` (positive: retained failure
  blocks resolution pre-publication; missing-parent recovery still refuses the
  failed parent; exact preview/publish/replay; fresh reconciliation with provenance
  pinned; explicit resolution + replay; old failure byte-exact);
  `TestRecoveredParentObservationRequiresExactLineage` (8 mutation/deletion
  variants each fail fresh validation and never auto-resolve);
  `TestRecoveredParentResolutionGuardsRevalidate` (post-reconciliation mutations
  fail the later state guards); `TestRecoveredParentPublicationRefusesNonRefusalOriginal`
  (a successful retained parent is never conflated; the ordinary same-bytes path
  is unaffected).
- `internal/runner/verification_bound_recovery_test.go` (extended; all five prior
  tests retained): `TestBoundVerifierRecoveryContextAdmissionIsReadOnly`,
  `TestBoundVerifierRecoveryContextRefusesMissingAuthorityWithoutWrite`,
  `TestBoundVerifierRecoveryContextRefusesEvidenceDriftWithoutWrite`.
- `internal/app/trajectory_verifier_recovery_test.go` (new), all with REAL local
  fake-CLI processes through the app entrypoints and the complete production
  parent path (`verifyTrajectoryWithAgent`; the helper runs the real built CLI's
  `trajectory verify-helper` inside the verifier child; no provider calls):
  `TestAppVerifierRecoveryFullLifecycle` (refused original launch → attended
  exact preview/apply through the `runTrajectory` dispatcher → no process on
  recovery-only/replay → stale-SHA/wrong-identity/conflict refusals → separately
  budget-authorized real relaunch under the pinned ID with AB/BA helper trace →
  immutable recovered observation tying the original failure to the replacement
  lineage → explicit fresh reconciliation + replay with original accounting
  byte-identical and the old failure retained → duplicate relaunch refused
  uncharged with no second process → publish replay exact);
  `TestAppVerifierRelaunchRequiresRecovery`; `TestAppVerifierRecoveryRefusedEvidenceMutationFailsRelaunch`;
  `TestAppVerifierRecoveryFailedReplacementRemainsUnresolved` (budget-refused
  relaunch consumes the pinned invocation with its own failed terminal, no
  started lifecycle, no observation, no resolution, no re-recovery, no charged
  replay); `TestAppVerifierRecoveredLineageMutationFailsReconciliation`
  (original/replacement evidence deletions and canonical mutations fail fresh
  reconciliation and the later state guards after a recorded resolution).
  All prior assertions in the pre-existing app/runner/trajectory tests are
  retained byte-for-byte.

## Invariants held

No budget extension/reset/refund; no model/quorum substitution; no general retry;
no history rewrite; no acceptance/signoff/final audit closure. The refused
invocation's telemetry, the failed parent result, the launch reservation, the
immutable recovery and the cycle ledger stay byte-exact through every step.
Every restart is separately checked (reservation boundary revalidates the whole
recovered lineage) and separately charged (ordinary budget admission/settlement).
Same-UID hashes bind bytes, never writers — the documented honesty boundary is
unchanged and is not used to skip missing-evidence refusal.

## Expected operator CLI flow

1. `parley trajectory verify --dir D --idea I --verifier V --yes` → verifier launch
   budget-refused pre-start; failed `parent-result.json` retained (stage `launch`).
2. `parley trajectory recover-verifier --dir D --idea I --run RUN` → inspect preview
   (+ printed replacement identity).
3. `parley trajectory recover-verifier --dir D --idea I --run RUN --replacement ID
   --sha256 SHA --yes` → one immutable `recovery.json`.
4. `parley trajectory relaunch-verifier --dir D --idea I --run RUN` → inspect plan.
5. `parley trajectory relaunch-verifier --dir D --idea I --run RUN --sha256 SHA
   --yes [--timeout T]` → real budget-authorized relaunch + helper + recovered
   parent observation.
6. `parley trajectory reconcile --dir D --idea I --run RUN [--sha256 SHA --yes]` →
   ordinary fresh reconciliation over the complete validated lineage.

## Files owned/written in this slice

`internal/trajectory/verification_recovery.go`, `.../verification_recovery_read_test.go` (new),
`.../reconcile.go`, `.../parent_recovery.go`, `.../verification_parent_recovery_test.go`,
`internal/runner/trajectory_verification.go`, `.../verification_bound_recovery_test.go`,
`internal/app/trajectory.go`, `.../trajectory_verify.go`,
`.../trajectory_verifier_recovery.go` (new), `.../trajectory_verifier_recovery_test.go` (new),
and this note. All other files byte-unchanged.

## Remaining limitations (explicit, not designed away)

- A failed bound relaunch consumes the pinned invocation (requested + failed
  terminal, no started); the recovery cannot re-recover it — one-shot by design
  (carried from zcode-1's note, unchanged).
- Same-UID forgery of a fully self-consistent recovery+evidence+lifecycle set is
  beyond byte evidence (unchanged honesty boundary).
- The old probe's race (both authority files disappearing between the unlocked
  Lstat and the mutating probe) is fixed by construction (single read-only path);
  the static missing/drift cases are tested deterministically, the race interleave
  itself is not reproducible without instrumentation and is documented, not faked.
- Publication of the recovered observation happens inside the attended relaunch;
  a separate operator preview/apply command pair for the observation alone is a
  possible later ergonomic, not a correctness gap (every read revalidates).
- The relaunch reuses the production `RunConsult` verifier path and the
  `os.Executable()` helper resolution; tests inject the built helper binary via
  `runTrajectoryVerifierRelaunchWith`, the exact same seam convention as
  `verifyTrajectoryWithAgent` — the app entrypoint logic itself is exercised
  unchanged.
- Budget for the relaunch is the ordinary operator binding: in production the
  persisted `parley budget` launch binding governs admission/settlement; the app
  commands add no budget logic of their own (tests inject explicit context
  policies for the denied original and the authorized replacement).
- App negative fixtures each rebuild the real fixture repo and drive real local
  processes; that is deliberate (refutation over reconstructed lifecycles) at
  the cost of suite runtime.
- No test executed in this launch; compile or assertion failures remain plausible
  and are the coordinator's to surface on the immutable pre-slice copy. Nothing
  here is PASS evidence, integration, or D-series acceptance, and no treatment
  experiment was started.
