---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: implementation dispositions and executed verification
not-a-signoff: true
---

# Verifier, readiness and process-budget integration checkpoint

This is a partial implementation, not AC-E/AC-B/AC-L acceptance or a Phase-6
signoff. Canonical participant reviews and handoffs remain unchanged. The live
comparative and packet experiments remain unrun and their pending decisions
are not inferred from continued implementation work.

## Independent source review and publication

Claude independently inspected verifier checkpoint 0a2022b in invocation
e04e426f-7bdb-4404-bb4b-39a38277d2aa (806.723s, process exit 0, missing artifact,
CLI-estimated USD 3.9186855). The facilitator's publication config mistakenly
retained the OLD artifact's escaped atomic-temp regex while updating its literal
path. This denied native Write's temp sibling; it was a facilitator configuration
error, not participant consent or successful artifact delivery. The review itself
was retained in the final result. Reported model usage includes Opus 5 and Haiku
4.5; no additional Parley participant or signature is inferred from CLI-internal
model usage.

A corrected exact regex passed a real sandbox canary: its native atomic sibling
could be created, an unrelated sibling could not. Claude resumed its own native
session and published the already-written review, with no source re-review, in
e876f9d6-5df1-4f84-bb21-77ed2b7a35c4 (186.606s, exit 0, CLI-estimated USD 1.9531635).
Codex verified byte equality with the artifact extracted from the original result:
43,184 UTF-8 bytes, SHA256
fb64642a7a60f0d47ee11b9107420dcd9790dd5ade31380410cd37f39fa6b443.
The canonical own review is claude-1-verifier-review-20260911.md. It reviews
0a2022b only, not the corrections described here.

## Verifier dispositions

- CRITICAL, original checks removed between ticks: corrected. Driver cursor v2
  pins the normalized names/commands before agent work, carries the binding
  through reconstruction, and rejects removal/change before more work or closure.
  Existing evidence activates the original-scope gate for legacy runs too. Direct
  Complete checks the actual run store's cursor and prior report. Executed restart
  tests cover deletion, replacement by scalar, renaming, command edit, shrinkage,
  and removal of both list and report; a positive unchanged case still closes.
  This is not same-UID authentication or a claim that all history can be deleted
  safely. Old v1/schema-less cursors remain readable; older binaries reject v2.
- MAJOR, generated completion self-invalidates evidence: corrected using Kimi's
  c232e1d API and handoff, integrated at 195fbed. The actual independent helper
  authorizes the exact status transition after attesting all criteria. Parent
  recomputes the same authorization before accepting it; Complete uses the same
  pure transform and verifies the exact captured before/after bytes before its
  synchronized write. The stronger real-process post-completion gate now passes.
  All other non-evidence scope remains hashed.
- Additional independently executed transition probes exposed a short-digest
  panic, a matched altered AfterSHA256 accepting a non-complete current status,
  and duplicate YAML status spellings/merge keys. The original five failures are
  retained. Serialized fixes validate the digest, actual final status and YAML
  key semantics; all five probes and the owner's transition suite now pass.
- MAJOR, retained before/after tree fields copied from the request: the reviewer
  explicitly confirms current runtime behavior is correct because both real
  bracketing TreeDigest checks compare to that value. No false-accept reproduction
  is established by this finding. Passing those measured values through explicitly
  would improve provenance clarity; this source cleanup remains open for review.
- MAJOR, receipt-write failure allegedly permanently prevents closure: the
  current driver performs fresh RunChecks before each closing Verify call, so a
  fresh driver attempt can replace the prior attested report. Replaying the same
  frozen helper request is intentionally refused. This permanent-failure claim
  remains unverified; an actual persistence-failure/recovery process test is still
  needed. No lost receipt is treated as successful verification.
- MAJOR, concurrent report overwrite: OPEN. Read/compare/Save is not an atomic CAS.
  All cooperating report writers need a shared serialization boundary; completion
  publication also needs that boundary. No concurrency acceptance is claimed.
- MAJOR, Windows verifier command quoting: OPEN. POSIX fixture execution and
  Windows cross-compilation do not establish a Windows runtime. A platform-aware
  invocation contract or explicit supported-platform restriction remains needed.
- MAJOR, refusal record not committed: receipts and invocation records are durable
  local ignored runtime artifacts by the existing telemetry design; the outer
  driver also escalates errors. Whether an additional bounded canonical summary
  is required remains an open disposition. No commit durability is claimed for
  ignored logs, and no cleanup of these worktrees is authorized.
- MAJOR, runtime path not ignored changes the tree: OPEN usability/recovery issue.
  The refusal is safe; validate the ignore prerequisite before creating artifacts
  and give an actionable recovery instruction, without widening exclusions.
- MINOR, positional receipt matching: OPEN cleanup. Original checks are now pinned
  including order, so reordering is diagnosed as a scope change. Name-keyed receipt
  reconciliation would still remove the implicit report-order dependency.
- MINOR, time.Time DeepEqual: current values come from JSON and compare correctly,
  as the reviewer notes; canonical serialization is a future robustness cleanup.
- MINOR, repeated whole-tree hashing: retained to enforce each actual criterion's
  before/after boundary. No measured performance improvement is claimed. Hoist
  only proven invariants; do not remove those observations.
- NIT, stale gate wiring comment and inconsistent path predicates: still open
  documentation/consistency cleanup; neither is silently treated as acceptance.

## Readiness correction and attribution

Hermes's first new attempt 5873e61e-0db0-417b-8014-cdfde97813a4 ended after
104.862s with exit 0, no artifact and four remaining compile errors. A compile-
only overlay then exposed 16 readiness failures: the token parser incorrectly
expected json.Decoder.Token to emit colon/comma tokens. The next owner call
de2694cd-9daa-4cb5-8849-c956778674a6 (99.184s, exit 0) wrote its own note and
corrected compilation, but a real bare-content false READY remained. Both costs
are unknown. The own source/narrative checkpoint is preserved at 0fc1aaa and
merged at fa98c96, without treating it as acceptance.

After those calls stopped, the integration owner serialized the recorded overlap
under parley-worktrees Sections 5/6. Corrections require explicit assistant or
known result provenance, reject semantic aliases/duplicates/null or mistyped
fields and competing outputs, propagate nested provider failures, distinguish
actual whitespace/dropped output on timeout, and preserve bounded capture.
Credential diagnostics decode valid JSON keys before redaction, cover escaped
values and omit malformed JSON text that could expose a truncated credential.
The original own handoff's alleged pre-existing missing rosterEntry symbol has
no actual compiler support and is not adopted as a diagnosis.

## Process-budget boundary

WithLaunchBudget now uses a separate, detached context channel that survives
WithLaunchInfo replacement. Every beginLaunch process path reserves before spawn
with the actual unique invocation ID; denied requests retain failure telemetry.
Printed handoffs explicitly do not charge a process. Terminal settlement runs
with its own bounded context even after timeout/cancellation, failed telemetry
publication or process failure. Unknown cost retains its conservative reserve;
known CLI estimates are rounded upward only for microdollar budget accounting.

Actual process fixtures cover manual execution, the common round process,
consultation, readiness probes, terminal execution and ACP. The first real child
runs, the over-cap second request is recorded without spawning. Further fixtures
cover concurrent roots sharing one ledger, corrupt history, failed start,
known-cost timeout, terminal-write failure, unknown cost and settlement failure.
The full runner race suite passes (41.265s).

LIMITATION: no app/operator caller yet resolves and attaches this policy to all
real product entrypoints. Persistent scope/policy configuration, legacy migration,
driver action/step/time charging and extensions remain required. These tests do
not make the existing manual recovery calls budget-enforced or prove all actual
model launch surfaces have been trialed.

## Executed validation

- Full go test -count=1 ./... PASS: app 84.272s, driver 5.892s, evidence 2.480s,
  runner 25.624s. Log integrated-verifier-readiness-budget-full-20260911.log.
- Scoped app/driver/runner/evidence/store vet PASS.
- Windows app test-binary cross-build PASS; runtime untested.
- Focused original-scope and real verifier transition tests PASS: driver 0.495s,
  app 13.750s. The earlier self-invalidation failure remains retained.
- Readiness adversarial/capture/diagnostic checks PASS (3.503s).
- Corrected completion-transition probes and owner suite PASS.
- Runner full race suite PASS (41.265s). Focused readiness/scope race PASS:
  app 10.305s, driver 1.767s. Actual shared-volume concurrent-root launch-budget
  checks PASS (0.613s).

All logs above live under the ignored .parley-runtime directory unless stated
otherwise. No final independent acceptance, live-model verifier trial, packet
experiment, full-six pilot, release, final merge or global install is claimed.
