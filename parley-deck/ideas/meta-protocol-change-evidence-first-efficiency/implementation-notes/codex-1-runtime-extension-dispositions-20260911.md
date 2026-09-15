---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: 9c1990c8e494048cb4b437389ebdf6f7f31b73e0
status: tested-slice-independent-acceptance-pending
---

# Lifetime launch and step extensions; evidence cancellation correction

## Resulting behavior

An exhausted launch or driver-step policy can receive an explicit finite operator
extension without deleting spent attempts or starting its lifetime clock again.
`parley budget launch inspect|extend` and `parley budget step inspect|extend`
expose the control. Inspection is read-only and creates no scope or lock state.
Extension requires an actual attended terminal, every absolute ceiling for the
selected kind, the inspected canonical policy hash, a unique decision ID, a reason
and `--yes`. Flags and participant frontmatter cannot manufacture attendance.
The runtime guide documents the complete commands and their units.

At least one finite axis must increase. An increased action ceiling must exceed
both the previous maximum and spent count; increased cost must exceed known
exposure; increased duration must exceed elapsed time from original activation.
Unchanged zero retains an originally unlimited axis. A grant cannot lower a
ceiling, turn a finite axis unlimited, or change an unlimited axis into finite.
The per-launch reservation is immutable through this control. Unknown exposure
requires reconciliation before a monetary increase; other extensions retain that
unknown state and cannot make it spendable under a monetary ceiling.

The first extension upgrades the policy from v1 to v2, preserving original limits
and ordered decision records. Each decision carries the prior canonical hash,
absolute limits, actual spent count, nullable exposure, original start, time and
reason. Exact replay is idempotent even after later grants; conflicting reuse or
a stale new decision refuses. The hash freezes policy rather than changing spend.
Resource and ledger locks are held through publication, serializing reservation,
settlement and reconciliation without rewriting charged ledger bytes. Failure
before replacement preserves the old policy; a replacement followed by a reported
failure can retain the grant, recovered through inspection or exact replay.

Bindings reload the current policy. Original or recorded intermediate runtime
settings remain valid references; unrecorded values refuse. The driver refreshes
step/time limits and uses persistent monetary exposure and effective grant,
separately from observed run usage. Telemetry records effective `max_cost_usd`
and exact `max_cost_micros`. Nested synchronous children remain one charged step.
V1 remains readable; older readers refuse v2. Bounds are 128 decisions, a 1 MiB
policy, 128-byte decision IDs and 1024-byte reasons. Invalid fields, lost charges,
clock regression and overflowing increased limits refuse.

A real local-process extension fixture exposed the existing nil/empty-map
comparison mismatch between trusted launch context and persisted policy. The
context copier now preserves the normalized empty-map representation. The
negative probe is retained rather than omitted from the execution history.

## Reproduced cancellation race and correction

The initial full Go suite passed, but the expanded budget/driver/runner/app race
run failed. `TestSerialVsBarrierConcurrencyFixture` exposed a read of `spawned`
inside `RunCriterion`'s `cmd.Cancel` concurrent with its post-Start assignment
at the former execute.go:77/90. The failed log is retained, SHA256
`3b4f2734fe569d53f153b5429868ecc58b4ba6f42944eb54f566feb482856d0d`.

Before changing production code, the new cancellation fixture independently
reproduced the same race with a real TERM-resistant descendant holding output
pipes. Negative log SHA256:
`652bf2ccc65d5085e5e8cf8056ceb441c68cb07b91d8083fcd5d495507448057`.
This is an implementer-authored regression probe, not independent participant
acceptance. Ownership was recorded before edits; Kimi's original work and review
artifacts are preserved.

Cancellation now uses the command Process published by exec before its context
watcher starts. The new session's PID is its process-group ID even if the shell
has exited; no mutable post-Start identity is captured. The fixture proves helper
readiness through a live TCP connection and descendant termination through EOF.
Cancellation before launch and failed Start remain failed evidence. The new tests
and the original serial/barrier fixture pass together under the race detector.

## Executed validation

| Check | Observed result |
| --- | --- |
| `go test -count=1 ./...` after correction | PASS; app 98.824s, budget 18.914s, driver 15.266s, evidence 6.118s, runner 35.887s; wall 102.638s |
| `go test -race -count=1 ./internal/budget ./internal/driver ./internal/runner ./internal/evidence ./internal/app` | PASS; 22.138s / 17.629s / 65.175s / 6.517s / 122.273s; wall 123.291s |
| Cancellation and original app barrier fixture under race | PASS; wall 7.434s |
| Scoped budget/driver/runner/evidence/app vet | PASS |
| Windows amd64 app test cross-build after correction | PASS; Windows runtime untested |
| Shared-volume budget/driver/runner/app extension fixtures | PASS; wall 3.151s / .934s / .673s / .780s |
| Shared-volume cancellation fixture under race | PASS; wall 2.867s |
| Actual unattended launch and step extension CLI calls | Both exit 2 with attended-terminal refusal; no runtime directory created |

Tests cover ledger/start preservation, stale/conflicting/exact replay, publication
recovery, concurrent decisions/reservations, settlement blocked through publication,
worktrees, nested sessions, clock expiry and rollback, schema/reader boundaries,
unknown exposure/reconciliation and actual manual-process continuation. The
earlier extension fixtures still match current source: only execute.go and the
new cancellation test changed after that fixture suite. No real grant was issued.

All 323 Go/module files match the post-correction source manifest, SHA256
`9fcc78dfcf8bfe9ef27093127a51ec6c5c79b10533e208c25b916e2f41cb8acf`.
Full log: 1404 bytes, SHA256
`af0944efeaa92340331d2f5580f4929d180bb593dbd6a6e2ffc3d30f8d612d96`.
Race log: 224 bytes, SHA256
`58648e2bf4cadfa81338883d247fb15e57fc9d46fd9aefeed1b8651cbeaf2de7`.

Thirty-eight logs/manifests/metadata files and verified checksums are retained in
the integration worktree's ignored
`.parley-runtime/runtime-extension-validation-20260911/`.
Original storage:
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-runtime-extension-validation-uisoudex/`.
Shared fixtures used
`/Volumes/My Shared Files/AI_WORKSPACE/parley-runtime-extension-fixtures-p80zb9ml`,
outside an enclosing Git repository. Local test binaries emitted through pipes;
logs were captured then written once. Pre-correction full PASS and race FAIL
remain separate records, not reclassified as passing.

## Remaining obligations

This control is not legacy migration, safe guard/lock-origin recovery, durable
semantic action replay, canonical refusal publication/recovery, or the opt-in
independently confirmed two-patch regression trajectory. Retained refused requests
can still require explicit migration before first configuration; deleting them
is not recovery. Terminal attendance and hash links do not authenticate a human
against a process with the same filesystem access.

Fresh independent source acceptance, human evidence-table integrity, actual
independent real-model concurrency/closure, the exact phase 1/6 packet experiment,
the twelve-task solo/duo/full-six pilot with frozen equal ceilings and rotation,
two blind nonauthor graders, owned reviews/signatures, final populated HTML QA
and real 14/30-day observations remain required. Claude recovery and historical
quorum/pilot/funding decisions are unanswered. No participant model was called
for this slice: inventory remains 34 terminal model attempts, 17 unknown costs,
USD 46.1887585 known CLI estimates and unknown total. This disposition is an
implementation execution record and does not mark the full goal complete.

## Updated report and actual browser correction

The English offline HTML was rebuilt at 2026-09-11T16:00:05.838108+00:00.
It is 609331 bytes, SHA256
`0b9e72687fdd50e19cb74f5938d5d9849d11cd837b39a027b4d092fc73456e43`,
and passes the ten report tests. An actual ego-browser check of the earlier
85e81b4 report found that expanding Technical progress log made a 390px page
541px wide: unbroken plain-text hashes and URLs overflowed document prose.
The existing template now wraps prose/list items with overflow-wrap:anywhere.

Fresh QA of the corrected exact hash records 62 passing assertions: all five
panels at 1440x900, 1280x540, 390x844 and 320x720, expanded log width, painted
charts, five embedded documents, search/reset, pagination, keyboard focus and
print-media visibility. Five screenshots are retained; the deep-log viewport
was captured through ego-browser CDP after the convenience screenshot helper
returned blank output. No Chrome control was used. The final QA record is
the evaluation workspace's delivery/2026-09-05/browser-runtime-extension-20260911.json.
Its raw successful log, screenshots and prior failed-attempt summaries remain
separate; the initial harness used a non-IIFE evaluation, checked collapsed
content and incorrectly captured stdout only. Those mistakes are not product
failures or passing evidence. The real overflow observation is retained distinctly.

Browser inspection certifies this partial report's current presentation only.
Independent source acceptance, future populated experiment results and physical
printing/PDF pagination remain unverified. The HTML is not rebuilt merely to
insert its own QA hash; the sidecar binds the already-tested artifact.
