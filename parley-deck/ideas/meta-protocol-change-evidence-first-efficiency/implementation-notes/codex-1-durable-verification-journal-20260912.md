---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-12
status: partial
scope: durable-captured-verification-journal
---

# Durable captured verification journal

## Implemented behavior

The charge-bound captured-source execution API now has shared persistent
verification storage. `PrepareCapturedVerification` reserves exactly one
`VerificationTicket` per original fixup charge and freezes its complete original
request, selected verifier, driver run and canonical origin root before launch.
An existing reservation is never reused, even after an interrupted request write.
`ReserveCapturedVerificationLaunch` reserves one distinct invocation before
spawn. `ExecuteCapturedVerification` takes an exclusive durable helper claim
before preparing sources or executing commands. Required file and directory
persistence barriers must succeed before proceeding.

The helper restores the original captured sources and records their private
roots in `prepared.json` before checks. Each actual AB/BA execution is published
as a bounded, exclusive, ordinal and hash-linked step. A receipt binds the exact
ticket, launch, helper PID/claim, source preparation and complete retained step
chain. Failed preparation, incomplete execution, authority drift and failed
observation writes remain failures. Receipt publication errors cannot become a
successful return. Failed records, partial files and interrupted claims remain
present; there is no automatic replacement, retry or release of the claim.

`ReadCapturedVerification` checks original state/ledger/archives and exact
prepared authority, then reconstructs the ordered observation from retained
steps. Missing, ambiguous, noncanonical, symlinked, extra, reordered or oversized
artifacts refuse. A missing/failed terminal returns an error and any already-read
partial observation. Complete receipts are assessed using the existing captured
criterion assessment, preserving clean/inconclusive/regression distinctions.
Reading or executing a journal never resolves an attempt or changes a charge.

## Failure and trust limits

This is an internal API. The instrumented runner and independent model/helper
CLI still need to be wired to the reserved invocation, inherited markers and
observed terminal. Caller-supplied invocation labels are not authentication.
A matching library receipt is not proof that an independent model called the
helper. Synthetic executable helpers in tests do not satisfy real model launch,
independent participant review/signature or live experiment requirements.
Same-UID callers can fabricate or delete private artifacts; this mechanism is
cooperative process lineage and crash-safe refusal, not a cryptographic identity
or hostile-user sandbox. Metadata remains private and is not exported telemetry.

A killed helper may leave private prepared directories and only the steps whose
writes completed before termination. `prepared.json` retains the directory
locations for explicit recovery. The unobserved in-flight outcome remains unknown;
no reconstruction from a textual PASS, new invocation or later source is allowed.
The raw retained journal survives original-authority failure, but the checked
reader refuses that authority; explicit recovery inspection remains necessary.
No generic operator recovery/retry/continuation or migrated state is introduced.

The journal requires a POSIX execution host. Windows compilation is checked
separately and does not establish Windows runtime support. Original source/index
reproduction remains exactly as the preceding captured-source checkpoint: staged
status mismatches or missing necessary history explicitly refuse. No arbitrary
historical Git/index/object backup is invented.

## Remaining work

Connect the prepared request and exclusive launch reservation to the actual
instrumented non-implementer invocation before spawn. Require the model to invoke
the matching helper, validate its inherited attribution and terminal, and retain
failed launches without replay. Derive the full expected patch sequence from
charge history, incorporate checked receipts and preserve the first two-regression
review trigger. Implement explicit durable review/recovery/continuation without
resetting charges or dropping quorum. AC-B2 and all six delivery areas remain
incomplete. Current-source independent acceptance, exact packet experiment,
full-six comparative pilot, report refresh and actual-delivery follow-ups remain.

Existing Claude provider-limit, historical Hermes/current Zcode, full-six roster
amendment and funding decisions remain unanswered. No retry, participant
substitution, real model invocation, actual operator activation, final merge,
release, deployment, global install or immutable-core publication occurred here.

The installer/runtime skill is 2.11.0. Session status/dry-run retains the known
source-role metadata drift and expected source-versus-packaged protocol difference.
The live source remains authoritative; no metadata sync or package substitution
was applied. OpenViking tools were unavailable; no shared-memory write is claimed.

## Verification

The initial focused run passed before two additional failure-boundary tests were
added. The complete final-source checks are recorded below.

Validation covers exactly 370 Go/module files, manifest SHA256
7158926537d9076a55286b7297c6875f17f1cc34df478217b54d4a70ad92f501.
Full JSON Go suite PASS (182.834s; all 32 package terminal events, including the
no-test cmd/parley package), six-package race PASS (195.817s), scoped vet PASS
(1.420s). Windows amd64 trajectory/app cross-builds passed and their PE amd64
headers and binary hashes were checked; Windows runtime remains unverified.
Compiled shared-volume trajectory/runner/driver/app fixtures passed in
222.010/9.169/1.525/3.229s. All ten new top-level test entries (nine behavioral
tests and their synthetic process-helper entry) passed in full, race and shared
runs. Complete native/shared logs match and contain no NUL bytes or failed-test
events. `final-verification.json`, the exact validation and verification scripts,
toolchain identity and 38 checksummed evidence files are retained under
`.parley-runtime/durable-verification-validation-20260912/`.

The original focused run is retained separately within that evidence directory.
No test failure was hidden or rerun to obtain a passing checkpoint. The preceding
sparse-log anomaly remains unclassified at its original source checkpoint; these
passes do not retroactively explain it. The older HTML report has the same
629641 bytes and SHA256 06ed736c5ef7379adf1ce3c3d1adc51f5afbb34ef6869b70798f8ad28b87a0be.
