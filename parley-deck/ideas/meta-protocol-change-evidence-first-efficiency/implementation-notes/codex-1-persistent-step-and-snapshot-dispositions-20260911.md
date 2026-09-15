---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: implemented-pending-independent-review
---

# Persistent steps and live-origin review accounting

This continuation addresses the executed Driver.Run resume counterexample at
2194dca, plus further production-boundary failures found while validating the
uncommitted correction. It is a D6 step/time slice, not completion of D6 or of
the six audit recommendations. No participant-owned review or signature is
changed. Claude's failed review remains terminal and the recovery choice is
pending; these changes have not received fresh independent acceptance.

## Behavior and scope

- Runtime step/time ceilings freeze an idea-scoped policy in the common Git
  administration (or local runtime storage for a non-Git deck). Omitted flags,
  a new Run or direct Advance load the saved policy. Changed nonzero limits
  refuse an implicit extension. Unknown/missing policy or ledger refuses reset.
- A private synchronous step session precharges the first mutable operation.
  Nested rounds/checks/signoffs share that charge; a later attempt gets its own
  identity. Await/no-op grouping does not reserve a step. Failed work stays spent.
  Session expiry prevents a retained context from authorizing later work.
- A charged transition may finish its own operations at the inclusive step cap,
  but cannot authorize a later launch after lifetime time expires or required
  accounting state disappears. This does not replace per-process hard timeouts.
- Run reports the persisted count and original activation time, including a
  failed charged attempt, instead of resetting reported totals on resume. Event
  append failures halt the driver. Exhaustion cannot mark an idea complete.
- Legacy cursors have no native idea identity. A validated, consistent event
  history associates them with their run; an unrelated identified cursor no
  longer blocks a new idea. Missing/conflicting identities refuse migration.
  Bounded stable file reads reject aliases, malformed and duplicate/case-aliased
  JSON. Requests catch failed actions that never published a cursor. Initial
  preflight and round-01 precede the driver and are not counted as old driver
  transitions. Absent history is not authenticated proof of never having run.
- Manual grouped runners retain a policy refusal until the common launch
  boundary has recorded a request. Every attempted child then has terminal
  refusal evidence without a spawn. An already-present artifact is still a
  no-op. Unobserved handoffs remain distinct from executed processes.
- A disposable Phase-6 clone keeps its live runtime origin in private context.
  Its process still executes in the snapshot, while policy and invocation
  records belong to the original idea and survive clone cleanup. Metadata is
  not a participant-controlled substitute for that lineage.
- Review publication now stages a unique file, uses the existing SyncFile
  helper, and performs a synced atomic replacement. On Darwin ENOTTY alone
  retries ordinary fsync; other errors remain errors. Failed publication keeps
  the snapshot recovery copy. This is cooperative filesystem behavior, not a
  certification of every remote mount or hostile same-UID tampering.

## Executed counterexamples and corrections

The original overlay against 2194dca observed round calls [2 3] across two Runs
despite MaxDriverSteps=1. The canonical regression now observes one charged
dispatch; same/new-run resume and omitted flags cannot grant the second.
Direct Advance returns the budget refusal; Run retains its clean escalation
return behavior with no further dispatch.

The new actual snapshot fixture failed in two ways before live-origin binding:
the step policy rejected the first valid review because its clone lacked the
origin policy, while the launch policy permitted a second launch. The local
failure log is retained at:

`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-step-snapshot-negative-v5lm4x37/negative.log`

SHA256: `f361a8eacbcb62e9238a83976ba4fb7671976ca53f373aeb8fbb83a9ea31be46`.
The corrected test requires one actual spawn and one recorded refusal under
each policy, successful owned artifact publication, real snapshot creation
without fallback, snapshot deletion, and both terminals retained at the origin.

Shared-volume execution then found a separate publication failure: File.Sync
returned ENOTTY at snapshot move-back. This negative result is preserved at:

`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-step-shared-check-lmbdmafr/runner.log`

SHA256: `b552cb0dc111582e91a37c66c7fb882a916c098f8d892e6b0c780ed63a0f836d`.
The corrected shared-volume runner fixture passes in 4.794s, using a locally
compiled binary and a fresh fixture root under AI_WORKSPACE verified outside
any enclosing Git repository. Its log is:

`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-step-shared-corrected-ve4ukvjm/runner.log`

Shared budget and driver fixtures passed in 2.468s and 0.647s. Each successful
binary emitted exactly `PASS\n`, SHA256
`c26de83abdc9496cd1301470918ec39ecca1cf389ef0ae1c6504da1800d1c431`;
identical hashes reflect identical tiny output, not a distinct evidence identity.
The initial attempt to create a fixture directly at the shared-volume root was
denied by the filesystem and did not execute tests. Failed fixture roots/logs
remain available; the successful script removes only its own disposable root.

## Validation of the final source in this slice

The final full `go test -count=1 ./...` passes (app 97.897s, budget 5.165s,
driver 6.543s, evidence 1.798s, runner 29.673s; wall time 100.504s). Its captured
1402-byte log has SHA256
`e1bf602e6a365fbe61f6b517913f472f020715d0514f51d6d3f3c759c8bd6c27` at:

`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-step-final-full-i16g7ih_/full.log`

Budget/driver race passed at the same respective source (13.154s / 11.095s).
After the later review-publication correction, runner race passed again in
54.408s, SHA256
`9a390ccc5e79304ce56c4ea03e5cd199bcfb025b13a3250f7636629a9c6f7df3`.
Scoped vet and Windows amd64 app cross-build pass; Windows execution was not
performed. The final scoped logs reside at:

`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-step-final-scoped-64i0u9tx/`

Successful earlier whole-suite output is retained separately from the final
rerun. All log capture used process pipes/in-memory output and one flushed write
to local storage. Copies and checksums are retained under the integration
worktree's ignored `.parley-runtime/persistent-step-validation-20260911/`.
These fixtures are deterministic/local-child evidence, not live-model
independent verification or satisfaction of the experiment gates.

## Remaining obligations

Persistent fixup and cross-review precharges, explicit legacy migration and
operator extension controls, monetary default mapping, canonical durable
refusal publication/recovery, and independently confirmed regression trajectory
remain open. Runtime step grouping alone does not establish durable semantic
action replay, exactly-once execution, or all AC-B1/B2 boundaries. This slice
does not certify arbitrary mount aliases or Windows runtime.

The exact phase-packet trial, frozen equal pilot ceilings, twelve-task
solo/duo/full-six treatments, blind nonauthor grading, real-model independent
verification/closure, participant-owned final reviews/signatures and final
offline-report browser QA remain required. The numerical inventory stays at
34 actual terminal model attempts with 17 unknown costs, USD 46.1887585 known
CLI estimates and unknown total. Local process fixtures add no provider calls.
Historical quorum and the full-six FINAL remain unchanged; no final merge,
release, deployment, global installation or immutable-core publication occurs.
