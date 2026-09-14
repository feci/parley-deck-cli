---
idea: meta-protocol-change-evidence-first-efficiency
status: in-progress
implementer: codex-1
started: 2026-09-05
branch: parley-deck-cli#integration/meta-protocol-change-evidence-first-efficiency
checkpoint-base-commit: 0585deb665375fed4b16d73ff0013027b22f03b7
validation-source-manifest: 35558a3f60bc892c15d68cfb063dec9b94ad54e8cbc75573667fec80a9f7cd7a
design-pr: https://github.com/feci/parley-deck-cli/pull/72
implementation-pr: https://github.com/feci/parley-deck-cli/pull/73
---

# Evidence-First Delivery

## Unchanged-source reconciliation allocation — September 12, 2026

Before edits, codex-1 claims new internal/trajectory/unchanged.go and
unchanged_test.go; serialized trajectory/reconcile.go, continuation.go,
parent_recovery.go and snapshot.go; new app/trajectory_unchanged.go and
trajectory_unchanged_test.go; app/trajectory.go; runtime docs; and owned note
implementation-notes/codex-1-unchanged-source-reconciliation-20260912.md.
No participant artifact, signature, protocol or frozen experiment is changed.

Reconcile an actually terminated charged attempt with identical captured source
bytes as a distinct observation, not an independently verified patch. Retain its
charge, full original before/after identities, invocation lifecycle and original
scope from the retained archive; recheck the current original quorum/criteria.
Do not invent a verifier run or parent result. Leave material criteria unverified,
preserve prior regression counts/triggers without adding a regression, and use
the existing inconclusive continuation acknowledgment. No attendance is added to
deterministic observation publication. Exact preview/apply/replay must preserve
historical decisions and refuse drift, missing evidence and changed scope. Verify
real local process/CLI execution, concurrent apply, later charged history, failure
and output/publication recovery. The original helper-ticket/orphan recovery and
all live/independent/final-delivery obligations remain open.

## Concurrent lock bootstrap correction allocation — September 12, 2026

The parent-recovery full suite exposed a real ordering race in existing budget
bootstrap: a caller could observe an absent origin before another initial caller
published its origin and ledger, then misclassify the new ledger as orphaned.
Before correction, codex-1 claims internal/budget/lock.go and new
internal/budget/lock_bootstrap_test.go, plus this implementation summary and the
owned parent-recovery note. Preserve missing-established-origin refusal and all
pinned-waiter/cutover checks. Add a deterministic concurrent-publication boundary
test and retain the failed full-suite evidence. Revalidate the final source after
the correction; do not reclassify the failing run as a pass.

## Lost parent-result recovery allocation — September 12, 2026

Before edits, codex-1 claims new internal/trajectory/parent_recovery.go and
parent_recovery_test.go; serialized trajectory/reconcile.go, state.go and
continuation.go; new internal/app/trajectory_parent_recovery.go and
trajectory_parent_recovery_test.go; app/trajectory.go; runtime documentation;
and implementation-notes/codex-1-parent-result-recovery-20260912.md.
These files are within the existing Codex trajectory/runtime ownership. No
participant artifact, signature, evidence-package implementation or protocol
may be rewritten by this claim.

Reconstruct a missing/unpublished parent observation only from its original
request, complete helper journal and successful matching process lifecycle.
Retain a separate immutable recovery record and the original absence/failure.
Pin exact inputs for apply/replay, preserve all charges and existing resolutions,
and support recovery of a lost previously reconciled parent only when the
reconstructed original facts match its retained digest. Recovery grants no
execution, budget extension, continuation acknowledgment or final acceptance.
Do not add a human gate for deterministic recovery of already executed facts;
keep the existing attended continuation gate. Test real helper/CLI execution,
conflicting evidence, publication interruption, concurrent replay and later
history revalidation. Other incomplete-ticket/orphan/unchanged-source recovery
and the original live/independent audit requirements remain separate obligations.

## Durable action identity allocation — September 12, 2026

Before edits, codex-1 claims new `internal/budget/action_identity.go` and
`internal/budget/action_identity_test.go`; serialized
`internal/budget/{ledger.go,reservation_receipt.go,reservation_receipt_test.go,step_session.go,cycle_session.go,cycle_binding.go,cycle_observer.go,policy_extension.go}`;
new `internal/app/budget_action.go` and `internal/app/budget_action_test.go`;
`internal/app/budget.go`; runner `runner.go`, `phase58.go`, `cycle_budget.go`,
`launch_budget.go`, `telemetry.go` and new `action_identity.go` / `action_identity_test.go`;
driver `budget.go`, `driver.go` and new `action_identity_test.go`;
runtime documentation and the owned note
`implementation-notes/codex-1-durable-action-identity-20260912.md`.
These claims are within Codex's runtime ownership and do not permit changing
participant-owned reviews, signatures, evidence-package implementation or protocol.

Persist logical operation/input identity separately from each spent attempt;
bind exact read-only replay to original accounting and reject changed payloads.
Integrate actual step/cycle call sites, retaining legacy records honestly.
Reading/replaying a receipt must not grant a duplicate execution, replace a
missing process terminal, reset accounting or claim whole-action completion.
Recovery of effects after lost parent results remains a distinct required gate.

## Guard-origin migration allocation — September 12, 2026

Before edits, codex-1 claims `internal/budget/lock.go`, new
`internal/budget/lock_origin_migration.go`,
`internal/budget/lock_origin_migration_test.go`,
`internal/app/budget.go`, new `internal/app/budget_origin.go`,
`internal/app/budget_origin_test.go`, `docs/agent-runtime-configuration.md`, and
`implementation-notes/codex-1-lock-origin-migration-20260912.md` for same-host,
same-resource cache relocation. These paths are within Codex's serialized
runtime ownership and do not modify another participant's artifacts.

The slice must pin exact origin bytes before waiting, retain old and destination
kernel exclusion through journaled cutover, preserve accounting and witnesses,
refuse incomplete migrated authority, and support exact recovery/replay. Missing
old identities, changed resource roots, and remote-host recovery stay refused.
Only synthetic test resources may be migrated during implementation. No actual
operator migration is authorized or claimed by this allocation.

## Unchanged-source reconciliation implementation — September 12, 2026

The allocated slice adds `trajectory reconcile-unchanged --sequence N` with exact
preview/apply/replay and typed unchanged provenance. It rederives complete original
source/archive/charge, archived original scope and matching normally exited
invocation lifecycle. It grants no execution, budget, continuation or acceptance.
All material criteria remain inconclusive, unchanged attempts preserve regression
streaks without reopening acknowledged review triggers, and original independent
patch verification continues to refuse equal material bytes. Historical parent
recovery and ordinary reads dispatch and revalidate both evidence types.

Final automated validation on 401 Go/module files returned exit zero: focused 155.765s, full 297.826s (32 package terminals), six-package race 352.626s, vet, Windows trajectory/app cross-builds and compiled shared-volume trajectory/app selections (45.861s / 64.312s). Eleven new top-level scenarios are present in focused/full/race and the applicable shared runs. Six removed-protection overlays returned nonzero at their intended assertions. The final verifier checked frozen source/test inventories, native/shared log equality, PE amd64 build headers and the unchanged historical HTML. Windows runtime and independent participant acceptance remain unverified.
The owned note is
implementation-notes/codex-1-unchanged-source-reconciliation-20260912.md.
All whole-goal live experiments, independent acceptance, final report/follow-ups
and other workflow-effect recovery remain open.

## Summary of work

The signed design is merged at 3a09a0cf2ef938e2456b50ec6eb34a6b5f37038c.
Implementation remains partial on draft PR #73. Integrated source includes typed
invocation accounting, actual launch/step/cycle reservations and finite attended
extensions; original-scope/current-tree completion evidence through an independent
helper; live-source full/shadow context packets; precise readiness observations
and buffered-process supervision. Unknown cost, failed attempts, missing authority
and incomplete evidence remain visible and cannot become success.

The current trajectory slice adds durable reconciliation of every original
charged patch, rechecking parent request/result, terminal and complete helper
journal evidence. Separate attended continuation records preserve original dirty
attempts and charges, retain the first two-regression trigger after corrective
work and leave acknowledged inconclusive outcomes visible. Historical exact
continuation replay returns its original preview instead of consulting a later
source. The current snapshot implementation verifies every regular member by
one additional fresh, contained read matching its copied bytes, identity, size
and mode, with stable descriptor metadata and post-read name checks. Full archive
and original-source checks remain required. That snapshot checkpoint passed
its 381-file source manifest checks: focused/full/race, vet, Windows builds
and compiled shared trajectory/runner/driver/evidence/app tests. The four new
reconciliation app tests and eleven snapshot-read tests pass in the applicable
runs; six negative overlays detect removed protections. Earlier failed runs and
counterexamples remain retained. This is automated implementation evidence;
fresh independent current-source acceptance is still pending.

The next implemented slice adds attended same-host guard/lock-origin relocation
with both kernel locks held, exact journaled recovery/replay and preserved ledger
bytes. The final 385-file manifest passed full tests (210.427s), six-package race
(247.623s), vet, Windows budget/app cross-builds and compiled shared-volume
budget/app selections. Eleven substantive new top-level scenarios and their
child-process harness pass. Five removed-protection overlays fail at their
intended assertions. A separately compiled unmodified 498a21c production lock
passed the old-holder/old-waiter/new-writer cutover proof. The actual CLI also
refuses unattended apply without mutation. This is automated implementation
evidence; independent current-source acceptance is pending. Missing old
identities, changed roots, lost journals and stale post-cutover material remain
explicit recovery limits. Independent current-source model acceptance, other
interrupted-attempt recovery and durable semantic action replay remain open.
The current action-identity slice persists declared input and original charge
identity atomically, with read-only exact accounting replay and conflict refusal.
It preserves distinct spent attempts, inherited group input, explicit weaker
policy/manual bases and legacy-unbound records. A reproducer exposed unpinned
charge fields; the correction binds the original epoch, kind, time and reserve.
A real fake-child fixup also exposed a missing returned invocation ID, now wired
through the existing observer. The 392-file source manifest passed full tests
(220.143s), six-package race (261.187s), vet, four Windows package cross-builds,
compiled shared budget/app/runner/driver checks, five intended negative overlays
and actual CLI inspection/replay/refusal checks. Windows runtime and independent
participant acceptance remain unverified. Receipt replay grants no execution
and does not recover a missing parent result or unfinished workflow effect.
The new parent-result recovery slice derives a separate observation from the
original independently executed request/lifecycle/helper records, preserving
missing/failed parent publication and all original charges. Exact replay survives
later legitimate work; new reconciliation pins recovery provenance, while an
older lost parent must reproduce its original digest and every fact. Contrary
partial fields, ambiguous structured prefixes, stale previews and changed scope
or accepted evidence refuse. It grants no execution, budget or continuation.
Full-suite validation also exposed an existing first-reservation origin-read race;
a fresh absence check fixes that misclassification while missing/conflicting
established origins/identities and pinned-waiter checks still refuse. The final
397-file source passed focused/full/race, vet, five Windows cross-builds and
compiled shared trajectory/app/budget tests. Eleven new top-level tests and six
negative overlays are accounted for. Independent participant acceptance is pending.
See `implementation-notes/codex-1-parent-result-recovery-20260912.md`.
The exact packet trial, real twelve-task solo/duo/full-six pilot, final populated
HTML/ego-browser QA and actual-delivery-based follow-ups are also outstanding.
No synthetic fixture supplies a live experiment or a participant-owned signature.
The historical evaluation and old HTML remain frozen. The live invocation inventory
is 34 terminal attempts, 17 unknown costs, USD 46.1887585 known CLI estimates and
an unknown total. Provider recovery, historical quorum/pilot amendment and resource
ceilings still need the previously requested decisions.

## Implementation plan / checklist

- [ ] AC-T1/T2/T3: typed unique invocation telemetry, every launch path, measured
  manual launch command, safe usage parsing, spend reservations and >=20 real calls.
- [ ] AC-E1/E2: independent current-tree evidence, fail-closed checks/goal completion,
  negative self/stale/skip/partial cases and real concurrency counterexample.
- [ ] AC-P1: one live-source packet renderer and CLI plus all three instruction
  sources, hash/omission/dependency guards, default full/shadow, no global publish.
- [ ] AC-B1/B2: precharged/idempotent budget boundary on manual/driver/resume/BLOCK
  paths and optional two-confirmed-patch-regression escalation.
- [ ] AC-L1: readiness observations, buffered soft guards, hard timeout and cleanup.
- [ ] AC-P2: exact ratified packet experiment and non-implementer recomputation.
- [ ] AC-X1/X2: frozen twelve-task real three-arm pilot and two blind nonauthors.
- [ ] AC-H1/F1: populated offline HTML, ego-browser QA and honest follow-up register.
- [ ] Independent per-slice review, agreed fixes, current-tree final verification,
  review consensus with all participant-owned signatures and zero agreed fixes.

Checks: focused Go package tests per slice, full `go test ./...`, race tests where
supported, CLI command fixtures, skill source tests, evaluation unit/containment
tests, real instrumented invocation reconciliation and live experiments. Baseline
Go tests passed on 257ef8c; those results do not certify this implementation.

## Worktree allocation

Base branch: integration/meta-protocol-change-evidence-first-efficiency
Integration worktree: ../worktrees/evidence-first-integration
Assignments below were accepted in consensus. They are exclusive file claims;
the facilitator serializes every shared entrypoint and dependency-file edit.

| Boundary (file set) | Owner | Branch | Worktree | Status |
| --- | --- | --- | --- | --- |
| New internal/telemetry/** and internal/budget/**; runner/app/driver integration paths listed below; .gitignore; owned codex handoff | codex-1 | integration/meta-protocol-change-evidence-first-efficiency | ../worktrees/evidence-first-integration | claimed |
| internal/protocolpacket/**; internal/app/protocol_packet.go; internal/app/protocol_packet_test.go; parley-deck/meta/packet-applicability.yaml; parley-deck/COOPERATION.md outside generated Section 2; parley-deck/meta/protocol-changelog.md; internal/protocol/defaults/COOPERATION.md (§9 item 1 mirror only — the drift guard `internal/protocol/drift_test.go` binds it byte-for-byte to the deck; claimed 2026-09-05 by claude-1, non-overlapping); owned claude handoff | claude-1 | feature/meta-protocol-change-evidence-first-efficiency/claude-1 | ../worktrees/evidence-first-claude | implemented on branch; pending independent review and integration glue (see implementation-notes/claude-1.md) |
| internal/app/preflight.go; internal/app/preflight_test.go; internal/app/preflight_hash_test.go; internal/app/preflight_liveness.go; internal/app/preflight_liveness_test.go; internal/runner/supervision.go; internal/runner/buffered_test.go; owned hermes handoff | hermes-1 | feature/meta-protocol-change-evidence-first-efficiency/hermes-1 | ../worktrees/evidence-first-hermes | claimed |
| internal/evidence/**; internal/app/driver_checks.go; internal/app/driver_checks_test.go; internal/app/driver_evidence.go; internal/app/driver_evidence_test.go; owned kimi handoff | kimi-1 | feature/meta-protocol-change-evidence-first-efficiency/kimi-1 | ../worktrees/evidence-first-kimi | claimed |
| Skill repo: skills/parley-deck/SKILL.md; skills/parley-deck/references/COOPERATION.md; test/packet-context.test.js; skills/parley-deck/parley-addon.json (generated per-file hash manifest, regenerated with `node scripts/build-addon-manifest.js` because the source edits made it stale; claimed 2026-09-05 by claude-1, non-overlapping) | claude-1 | feature/meta-protocol-change-evidence-first-efficiency/claude-1 | ../worktrees/evidence-first-skill | implemented on branch; pending independent review |

Codex's serialized existing paths: internal/runner/{runner.go,consult.go,acp.go,
steer.go,phase58.go,handoff.go,failclass.go} and their tests; new
internal/runner/{launch.go,launch_test.go,telemetry.go,telemetry_test.go,
protocol_context.go,protocol_context_test.go}; internal/app/{app.go,app_test.go,
protocol.go,protocol_test.go,consult.go,consult_test.go,driver_impl.go,
driver_impl_le_test.go,goal_check_test.go,consensus_request_signoffs.go}; new
internal/app/{agents_exec.go,agents_exec_test.go,budget.go,budget_test.go};
internal/driver/{checks.go,checks_test.go,cursor.go,driver.go,impl.go,impl_test.go,
driver_test.go,loop.go,loop_budget_test.go,close_integrity_test.go,strict_gate_test.go};
internal/store/{events.go,events_test.go}. Additional integration-only paths need
a recorded non-overlapping claim before editing, not implicit permission.

Additional Codex allocation (2026-09-05, before edits):
internal/acp/spawn.go and internal/acp/spawn_test.go for content-free stream
observers and observed process exit status used by runner telemetry. These paths
do not overlap any participant's claim; no ACP schema/billing inference is added.

Additional Codex allocation (2026-09-05, before edits):
internal/app/preflight_evidence_test.go for non-owner negative integration cases.
Hermes may refine only preflight_liveness.go and preflight_liveness_test.go in a
follow-up attempt and write implementation-notes/hermes-1-fix-readiness.md.
Codex serializes CommandFor glue in preflight.go after the original handoff.

Every participant writes its own supporting handoff at
implementation-notes/<agent-id>.md. It may change only its own allocation row in
its worktree if useful; it must not overwrite the shared implementation narrative
or another participant's note. Source changes are committed on the named branch;
integration merges are sequential. No shared mutable test database or local env
is needed. All mutation probes use per-worktree temporary directories.

Generated packet bodies and manual runtime records use explicitly ignored
`.parley-runtime/` paths; Codex owns the ignore entry and runtime wiring. A
participant may report an integration dependency in its handoff, never falsely
claim that a missing shared entrypoint is already wired.

## Integration log

- 2026-09-05 codex-1: created integration worktree from the merged design commit;
  claims recorded before production edits. No slice merged yet.
- 2026-09-05 12:05Z codex-1: telemetry storage/parser foundation committed at
  4a9697a. Shared exec wiring is in progress; the existing runner, telemetry and
  store test packages passed locally. This is not independent final acceptance.
- 2026-09-05 12:05Z codex-1: Kimi attempt
  1d04e5ba-c074-4951-8694-543393b34a21 ended on a provider-reported weekly quota
  (403), leaving uncommitted internal/evidence files without an owned handoff.
  These files remain preserved in its worktree; no completion is claimed. The
  user was asked whether to waive Kimi for this idea's remaining quorum. Pending
  that answer, quorum is unchanged and final closure remains gated. Other
  participants may continue their non-blocked implementation work.
- Additional non-overlapping allocation: Claude may update
  internal/protocol/defaults/COOPERATION.md to mirror the source protocol change.
  The worktree already shows this mirror edit; the original allocation omitted
  that required embedded copy. No other owner claims that path. This late claim
  correction is recorded rather than represented as a pre-edit allocation.
- 2026-09-05 codex-1: integrated Claude's renderer at 876aafb and Hermes's
  supervision/readiness at aaeaaaf, then Hermes's readiness follow-up ff3df32
  and owned handoff 9238246. Integration merge HEAD is 47661e5. The four new
  non-owner rejection cases pass with classifier/real-child fixtures (1.489s).
  Full `go test ./...` passes on the integrated worktree including serialized
  CommandFor and packet-command dispatch glue. This is not final independent
  verification or a claim that AC-L1 is satisfied in full.
- 2026-09-05 codex-1: source inspection leaves stricter readiness-schema checks
  open: assistantPayload accepts role-less content fields without a recognized
  message type; assistantRoleOrAbsent treats a present non-string role as absent.
  Capture bounds and diagnostic-tail secrecy also remain open as reported by
  Hermes. These are not waived by a green fixture suite.
- 2026-09-05 codex-1: user requested an English presentation. The evaluation
  report template, status copy, chart labels and document navigation are now
  English. Separate translated documents preserve the frozen assessment; the
  builder checks numeric/link/heading invariants and records both source hashes.
  This is a presentation update, not an expansion of implementation completion.

## Deviations from FINAL.md

None. Source-only skill changes are separate from global installation/release.
The experimental variants and enforceable resource policy are not frozen yet.

## Decisions after FINAL.md

- Per-invocation raw logs for new manual commands must be private runtime files,
  not public telemetry. Metadata allowlists do not authorize exporting raw logs.
- The evaluation harness distinguishes infrastructure failure from candidate
  failure; independent review found and corrected system Python framework launch
  compatibility. Those are preparation results, not real model comparisons.

## Current state & next steps

1. Parent-result recovery, concurrent lock-bootstrap correction and unchanged-source
   observations are published through 2ab9e82 (local and remote, PR #73).
   The unchanged-source observation slice passed final automated validation; see
   `implementation-notes/codex-1-unchanged-source-reconciliation-20260912.md`.
   Keep independent current-source acceptance separate from this code checkpoint.
2. Complete workflow-effect recovery for orphan reservations, failed/consumed
   helper tickets and unavailable source/index/history. Normally exited unchanged
   source now has a separate deterministic reconciliation control; incomplete
   lifecycle/capture cases remain unresolved.
   The implemented parent recovery requires a successful matching invocation
   and complete helper journal. Preserve all charges and refusal evidence;
   every new execution needs a separately spent attempt. The next read-only map
   is `.parley-runtime/helper-ticket-recovery-next-20260912.md`.
3. Obtain fresh participant-owned acceptance of the final code and finish live
   launch-surface coverage plus independent real-model concurrency/closure evidence.
   Claude's latest attempt hit its weekly limit without a review artifact; its
   reported reset is September 14, 06:00 Europe/Berlin. Do not silently replace it.
4. Resolve historical Hermes/Zcode membership, any explicit full-six pilot amendment
   and experiment funding/resource ceilings. Preserve the signed scope until an
   amendment is accepted; current global defaults do not rewrite historical quorum.
5. Run the exact packet trial: six AB/BA pairs each in phases 1 and 6, three packet
   canaries plus one full control, with independent recomputation and original
   ship gates. Run the twelve-task solo/duo/full-six pilot with frozen equal
   enforceable ceilings, role rotation and two blind nonauthor graders.
6. Finish all participant-owned reviews/signatures, populate the final offline HTML,
   verify it through ego-browser and register 14/30-day follow-ups from actual
   delivery. Keep PR #73 and skill-source changes in draft until their gates pass.
   No final merge, global installation, immutable-core publication or release is
   authorized by this implementation task.

## Recovery / resume

### Active driver persistence work

Codex will harden the already-owned `internal/driver/cursor.go` and its driver
tests before adding shared budget entrypoints: make reservation writes durable
with the verified filesystem sync policy, and reject malformed or negative
charged state. Preserve valid historical cursors whose omitted zero field was
written by the prior schema. This is persistence validation, not authentication
against a same-user process capable of replacing the entire state.

### 2026-09-10 — recovered interrupted Codex session

Additional non-overlapping Codex allocation: `internal/fsutil/sync_darwin.go`
and `internal/fsutil/sync_other.go`. Move the tested sync adapter there so the
telemetry and Claude-owned packet publisher share the same filesystem policy.
Claude's packet handoff reproduces the same shared-volume ENOTTY failure and a
fresh-directory mkdir race; those findings require an owner follow-up before
integration. The non-atomic exclusive-create fallback also remains unacceptable
for an immutable publication contract. Actual hardlink creation on the shared
volume succeeded, so unsupported hardlinks can fail closed without a fallback.

Recovered the original local conversation `019e5f54-222b-7802-8342-b98dfb89d306`.
Its last substantive user instruction (2026-09-09 20:54Z) requested all six
audit recommendations through Parley Deck. The session then ended on a workspace
spend-cap error, including two failed continuation attempts. The history still
exists in the May 25 rollout file; this recovery does not rewrite application
history or participant signatures.

The implementation worktrees remain at integration commit `88d650d`; PR #73
is still a draft, with an older remote head `598324e`. Claude has no new resume
handoff; Kimi's uncommitted evidence package remains preserved. Existing roster
and pilot-amendment questions are unresolved and were surfaced again. Unrelated
implementation can continue, but no quorum replacement or full-four pilot is
effective yet.

Next bounded recovery work (Codex-owned telemetry and measured-launch paths):
- [x] Fix a reproduced launch blocker on this shared macOS volume: Go File.Sync
  returns ENOTTY for F_FULLFSYNC while ordinary fsync succeeds. Preserve required
  synchronization and all other write failures; do not silently skip sync.
- [x] Preserve legitimate bracketed model identifiers in requested/reported
  telemetry while retaining credential-shaped metadata rejection.
- [x] Make a pre-record launch failure diagnosable without emitting fake schema-0
  telemetry, raw prompts, or unsafe diagnostic content.
- [ ] Rebuild the uninstalled binary, retry only missing Claude/Kimi actions,
  and record actual results. Continue remaining runtime/budget integration.
- [x] In the already-owned driver loop, reconcile costs by invocation identity;
  unknown, conflicting, malformed or missing accounting must stop a configured
  monetary ceiling and remain null in the total-cost display. This is one part
  of D2/D6; it does not claim per-launch reservations or all budget paths complete.

The shared-volume regression was reproduced with `TMPDIR` on that volume and
`go test ./internal/telemetry -run '^TestLifecycleAndUnknowns$' -count=1`:
`persist requested invocation before launch: sync ...: inappropriate ioctl for device`.
The same filesystem accepted Python `os.fsync` on an actual written file. These
are implementation observations, not independent completion verdicts.

Recovery validation: the failing shared-volume lifecycle test now passes, along
with manual-launch and bracketed-model regressions on that volume. The full Go
suite passes after the telemetry/launcher changes. Telemetry and runner race
tests and their app/runner/telemetry vet checks pass. Driver tests pass after the
additional cost-reconciliation change; duplicate replay, distinct retry,
unknown/null/invalid cost, conflicting identity, missing/corrupt event logs and
inclusive monetary-ceiling cases are covered. No test result here closes an AC
that requires independent verification or live experiments.

Live recovery attempts use the rebuilt, uninstalled binary. Claude invocation
`a23b2caa-9d10-447e-9a98-3101f9883202` started with requested model
`claude-opus-5[1m]`, max effort, native restricted file tools and explicit
`acceptEdits`/no permission prompts; native restricted mode demonstrably confines
file tools and rejects bypassPermissions (per installed CLI help). The Parley
inventory's legacy bypassPermissions declaration reports AUTO=no for this
narrower mode; do not misreport that inventory as verified. Its effective
file-write capability is established by the native mode and actual owned edits.

Kimi invocation `4e56455d-873e-46bf-90e4-f3a1c5058ee8` failed at native storage
setup with `permission denied`, before an owned artifact. Invocation
`f7248aa0-d6cc-4dc4-9d69-7341cce6f61d` then started with `kimi-code/k3`, max from
native config, a dedicated `KIMI_CODE_HOME` under the private local recovery
directory, and the same OS workspace write sandbox. Only that private runtime
directory and its assigned worktree are writable. A concrete probe permits
own-worktree writes and rejects a sibling-worktree write. Historical failed
attempts remain untouched. The Kimi handoff remains pending. Claude has written both the initial resume
handoff and its owner correction; see the validation entry below.

No raw credentials, prompts or private session output are copied into these
canonical records. The original session's workspace spend-cap failures are
historical; no current failure is labelled a spend-cap failure without evidence.

Use this manifest and actual git worktree/branch state. Never delete or reset an
unexpected worktree. Keep failed agent attempts and authored artifacts. The design
FINAL is immutable; implementation details and genuine deviations belong here.
No package, global skill/core installation or production deployment is requested.

### 2026-09-10 — owner packet correction and cursor persistence

Claude's corrective invocation `aa4c4736-9315-4d5e-a834-f5f11081fbe0` exited 0
and wrote its own `implementation-notes/claude-1-recovery-fix-20260910.md`.
Elapsed process time: 433.850s; reported CLI cost estimate: $2.12302 (not an
invoice). Usage reports Opus 5 [1m] and Haiku entries, so no singular resolved
model identity is inferred. The owner explicitly ran no shell tests. Codex then
executed the packet package, 20 repetitions of concurrent-publication tests,
shared-volume publication/flush fixtures, race tests, app packet tests and vet;
all passed. The prior non-atomic direct-write fallback was removed by its owner.
These are scoped executed checks, not a final review verdict or AC-P1 closure:
renderer-to-launch wiring and attestation are still outstanding.

Codex changed the driver cursor reservation writer to unique staging, file sync,
atomic replacement and directory sync. Version 1 writes its charge explicitly;
null, duplicate/aliased fields, negative counters, unknown schemas and incomplete
objects halt instead of recreating a zero budget. Complete historical schema-less
cursors retain the prior omitted-zero convention. Driver tests, race tests and
vet pass. Directory synchronization is required rather than silently ignored;
platforms that reject it fail closed and remain unvalidated (including Windows).
This does not claim same-UID tamper prevention, cross-entrypoint action identities,
persistent loop counters, or completed per-launch monetary reservations.

### Next serialized launch integration

Codex is wiring the shared renderer into the measured manual launch boundary
and handoff path in already-claimed runner files. Full context stays mandatory;
unknown phase/track input is reported by the shadow renderer, never optimized
away. Authority/publication/secret refusal must produce failed attempt evidence
without starting a process. This first integration does not cover the remaining
round/ACP/consult/signoff launch paths; those remain open.

### 2026-09-09 22:45Z — evidence owner timeout and reproduced blockers

Kimi invocation `f7248aa0-d6cc-4dc4-9d69-7341cce6f61d` reached the configured
30-minute hard deadline (1800.112s), exit 143, failure_class=timeout. It produced
partial owned code but no handoff. Reported cost and token usage are unknown;
no quota/spend-cap diagnosis is inferred. The partial work remains preserved.

Codex executed independent probes against the exact evidence source snapshot
SHA256 `a71f397b8419ba715e3ffa1eea99bb1c6fafd0fe2603f206038252fa5f3c84a5`:
- A report with no verifier attestation passed merely by supplying another ID.
- A missing current-tree digest was accepted.
- A malformed final envelope retained an earlier valid pass envelope.
- Retargeting an in-root symlink, and changing a file's executable bit, each
  retained the same tree digest. macOS /var aliasing also caused a false escape.
The probe is private runtime material; no owner code or artifact was rewritten.
Source inspection additionally shows the tested digest is taken after execution
without a before/after stability check, and report synchronization errors are
suppressed. These are owner fix-up requests, not accepted design deviations.

### Next close-gate hardening

Within the already-owned app/driver files, Codex will make missing, self,
failed or ambiguous goal-check execution reject closure and require an exact
PASS token rather than a prefix match. Typed criterion evidence remains a
separate required gate owned by Kimi; a model's textual PASS cannot replace it.

### 2026-09-10 — measured full-context launch and fail-closed goal check

The manual `agents exec` boundary and interactive handoff now obtain full
protocol text from the shared renderer, pass exactly the in-memory attested
bytes, and include the shadow audit in the private prompt. Missing authority,
secret refusal or failed body publication records a failed invocation without
spawning a child or writing a handoff prompt. The canonical source is re-read
for each new attempt; callers cannot overwrite the computed attestation.
Other exec/ACP/consult/signoff paths still require integration, and this does
not claim AC-P1 complete.

A real owner recovery was launched through the new boundary: Kimi invocation
`c16ced7b-14d0-4e82-80d6-0f88d035fb36`, started 2026-09-09 22:47Z, records full
context with source and packet SHA256
`317ac8a2279052c827ae31f173949348989cb83047a65545baa2f79a773367e6`.
Its early owned handoff exists; final outcome is pending. No unknown spend is
represented as zero.

The goal check now rejects absent/self/unavailable checkers, process errors,
nonzero exits even with a printed PASS, and ambiguous/reserved/prefix-only
verdicts. An exact textual PASS is still insufficient to replace typed
criterion evidence; that integration awaits the evidence owner's corrections.
Manual-signoff errors also preserve the underlying launch refusal rather than
replacing it with an empty-signoff validation diagnostic.

Validation after these changes: full `go test ./...` passes; targeted manual/
handoff/goal-check race tests and app/runner/driver vet pass. The initial full
suite exposed two legacy manual-handoff fixtures without source metadata; both
now explicitly declare their test protocol authority and pass. Missing authority
still has its own rejection fixtures. No independent final verdict is claimed.

### Additional non-overlapping Codex allocation — replacement portability

Claim `internal/fsutil/replace_unix.go`, `replace_windows.go`, and
`replace_test.go` for a small platform replacement primitive used by the cursor.
The current POSIX directory-sync sequence cannot be assumed to work through a
Windows read-only directory handle. Use the existing x/sys/windows dependency's
MoveFileEx with REPLACE_EXISTING and WRITE_THROUGH there; keep rename + directory
sync on other platforms. No cross-volume copy fallback. Verify native behavior
and Windows cross-compilation, with Windows runtime behavior explicitly untested.

### 2026-09-10 — replacement portability validation

`ReplaceSyncedFile` now uses rename plus required directory sync on POSIX, and
MoveFileEx(REPLACE_EXISTING | WRITE_THROUGH) on Windows without a copy fallback.
The staged file is synchronized before either replacement. Native fsutil/driver
suites, shared-volume replacement/cursor tests and vet pass; the Windows fsutil
test binary cross-compiles. Windows runtime behavior is not tested here. This
supersedes the earlier deliberate Windows fail-closed limitation of attempting
to flush an ordinary directory handle; it adds no claim of universal power-loss
or network-filesystem durability.

## Outcomes & surprises

Recovery resumed actual implementation and found additional false acceptance
paths through independent checks. Integrated changes have the recorded local
validation below; the Kimi slice is explicitly unintegrated. The separate pilot
harness passes on Python 3.9 and 3.11 with independent reviewer reruns, but no live
treatment has executed. Historical quota failures, timed-out attempts and unknown
usage remain preserved; no missing execution is reclassified as success.

## Validation evidence

Not yet independently verified against the implemented tree.

## Notes for reviewers

Review against every AC in FINAL.md. Do not count a self-authored pass, an
unexecuted task, unknown usage, stale tree or missing scope as success. Reviewers
must exclude their own slice from independent verdicts and verify another
owner's work directly. Correlated agreement is not a correctness result.

### 2026-09-10 — independent runtime review disposition in progress

Claude's authored partial review of 8b11f1d is preserved verbatim in
implementation-notes/claude-1-runtime-review-20260910.md (invocation
48b96bfa-6584-43e8-8c8e-8070eca0bc1e, process exited 0, 438.381 seconds,
CLI-estimated USD 3.3188955; reported model list includes Opus 5 and Haiku).
The reviewer inspected source without executing tests; this is not Phase-6
acceptance. The newer 342fcf0 Windows replacement addresses its directory-handle
portability concern in code; Windows runtime validation remains open. POSIX
directory-flush errors remain fatal because ignoring them would weaken durable
reservation semantics. No universal filesystem durability is claimed.

Codex is correcting the mixed-verdict/trailing-template false PASS, preserving
safe refusal classes, guarding protocol delimiter collisions, labelling shadow
counts as unapplied diagnostics, and retaining the file lifetime during Darwin
fsync fallback. Claude is updating the owned live/embedded LE-7 protocol text
under already-ratified D3; the old fail-open sentence is a temporary unresolved
implementation inconsistency, not an accepted deviation from FINAL. Strict or
auto ideas with unverified goal checks intentionally remain open under D3.

Partial phase mapping, headless-signoff context asymmetry and other launch
surfaces still need integration. Handoff prompt publication remains separately
open. The review's redundant cursor switch is harmless and deferred as a NIT.

Runtime-review corrections above pass app, driver, runner and fsutil package
suites (app 70.646s; runner 18.995s). The regression includes FAIL followed by
a PASS template, ambiguous-before-PASS, and protocol envelope collision refusal.

### Remaining launch-context integration plan

Codex will apply the same renderer at the actual tracked-command, supervised
exec and ACP boundaries, removing duplicate manual preparation. Every boundary
will overwrite caller-supplied context attestation and persist refusals before
returning without process start. App call sites must not replace the prepared
stdin afterward. Test fixtures will explicitly supply synthetic source authority;
missing authority and malicious source remain separate rejection cases. This
also covers round/review/implementation/fixup/steer/consult/headless-signoff
execution through their existing shared launch primitives. No optimized packet
launch or experiment result is introduced. Preflight CommandFor glue remains
Codex's previously recorded serialized allocation; Hermes source/tests stay owned.

Capability/readiness probes are distinguished from protocol tasks: their bounded
PONG/runtime-artifact instruction must work before a deck or installed core
exists. They use an explicit probe-only instrumented boundary, restricted to
preflight/runtime-probe phases, which replaces caller attestation with
context_mode=probe-only and no source/body hashes. This interpretation of D4
keeps every protocol-task builder on the renderer; probe-only is not a full
context claim. Independent review of this distinction remains required.

### 2026-09-10 — protocol-task launch wiring validated locally

The tracked-command, supervised-exec and ACP boundaries now resolve and attest
full live protocol bytes per actual attempt. This covers headless signoff,
round/cross-review, implementation/review/fixup, steer and consult/goal-check
through existing callers. Caller-supplied attestation is overwritten; direct
stdin overwrites were removed at app call sites. ACP's real child fixture checks
that one attested protocol reaches the prompt. Negative command/exec/ACP cases
record refusal with no process start even if the caller supplies a forged hash.

Round numbers now map to Phase 1/2; design/review consensus signoffs and handoffs
carry Phase 3/7 and idea identity. Explicit canonical kickoff tracks populate
shadow metadata; absent/unknown values remain unknown. Consult/steer with no
conceptual phase remains unknown rather than guessing applicability.

Handoffs publish unique per-invocation prompt names through synchronized staging
and replacement; later attempts cannot overwrite an earlier attested prompt.
The stable human instructions use atomic replacement. The regression verifies
both retained attempts. Runner/app suites pass (19.548s / 72.766s). These are
facilitator-run tests, not independent final acceptance.

Claude's protocol-text fix at eba3a0a is integrated; the live/default close rule
now matches D3. Its protocol/packet suites pass under facilitator execution.
The invocation 2ebe2909-6446-4b67-85ed-c319f391d241 exited 0 after 306.630s,
CLI-estimated USD 2.4902795 with Opus 5 and Haiku in the reported model list.
The earlier changelog statement about Codex's stale comment/aggregation refers
to the owner's reviewed worktree; e2f9d19 already corrected those in integration.
Separate skill-reference mirroring and review remain open.

### 2026-09-10 — Kimi second timeout and independent check

Invocation c16ced7b-14d0-4e82-80d6-0f88d035fb36 reached its 30-minute deadline
(1800.905s, exit 143). Usage/cost remain unknown. The early owner note and all
partial source are checkpointed on the Kimi branch, not merged. The new
Codex-owned check in implementation-notes/codex-1-kimi-recovery-check-20260910.md
records four failing fixture groups and two reproduced false-acceptance paths:
duplicate envelope fields, and unrelated/unbound independent rerun attribution.
Earlier missing-verifier, malformed-final-envelope, symlink and mode probes now
reject correctly. No owner-completion or typed-evidence integration is claimed.

At f84560b the full Go suite passes. Focused runner/app race checks, shared-volume
context/handoff/replacement checks and runner/app/driver/fsutil vet also pass.
These results do not cover or certify the unintegrated Kimi checkpoint.

### 2026-09-10 — second Claude review disposition

Claude independently reviewed f84560b by source inspection and preserved its
verdict in implementation-notes/claude-1-context-review-skill-sync-20260910.md.
It confirms prior mixed-verdict, signoff-context, refusal-diagnosis, Darwin file
lifetime and Windows-directory-handle findings are addressed; it agrees that
readiness/runtime probes need a separately labelled no-protocol-task boundary.
This is partial source review with no reviewer-executed tests or final signoff.

Codex accepts the opening-delimiter finding and adds refusal for both markers.
The probe boundary now has negative phase-guard and explicit no-hash telemetry
fixtures. The API remains an internal trusted-caller boundary; a stronger
constructor that cannot accept free-form task text is a remaining refinement.
The spawned interactive TTY path still lacks its own process invocation record,
as separately identified by the reviewer; the unobserved handoff record is not
misrepresented as that process. Handoff retention and cross-round source-hash
reconciliation remain follow-ups, with no full launch-coverage acceptance.

Counter-position to Claude's remaining directory-sync MAJOR: an unsupported
barrier cannot be reported as satisfied for a durable budget reservation.
Native and this shared-volume filesystem pass; other filesystems that reject
directory fsync must stop visibly. Windows has a separate implementation.
Codex claims internal/fsutil/replace_unix_test.go before editing to demonstrate
that EINVAL/EOPNOTSUPP and real I/O errors all propagate after conservative
replacement; no silent unsupported-errno success fallback will be introduced.
This disposition is not a completed review-consensus signature.

Final local edge-case fixes pass runner/fsutil/protocol packages after source
synchronization. Full-suite evidence remains tied to f84560b; subsequent changes
have their relevant package checks, not a fabricated full rerun.

Eight unique real recovery terminal records are reconciled in the sibling
evaluation delivery's runtime-recovery-invocations.json. Five reported CLI
estimates total USD 16.125529; three costs are unknown, so the total is unknown.
This does not satisfy the >=20-attempt gate. No roster change, live-pilot
treatment, packet optimization rollout, release or global installation occurred.

The full separate skill Node suite reached 397 tests: 395 passed and two
packet-context tests failed on obsolete wording assertions. The owner is
correcting those tests and clarifying bundled snapshot use as local orientation
only, not authority for a protocol-task launch. No passing skill-suite claim
is made before that correction is executed and checked.

### 2026-09-10 — recovery checkpoint ready to resume

All recovery participant invocations have terminated. The final Claude owner
follow-up 3e0c01f4-bf68-4422-8df4-ade47c2efe26 completed in 161.221s with
CLI-estimated USD 1.3832965; reported models include Opus 5 and Haiku. The
reconciled recovery inventory now has nine unique terminal attempts, six known
estimates totaling USD 17.5088255 and three unknown costs; the overall total
remains unknown. No >=20-attempt gate is claimed.

The skill source at a5e9fbb is committed and available as draft PR
https://github.com/feci/parley-deck-skill/pull/8, linked to CLI draft PR #73.
Full npm test passes: 397 Node tests, 54 Python tests on Python 3.14, and all
add-on manifest checks. Its six targeted context tests also pass. The earlier
wording failures are corrected by the owner, with the authored handoff preserved.
Bundled snapshots are now explicitly local orientation only, not launch authority.
No globally installed source was changed and neither PR is merged.

The offline report was rebuilt with the final recovery inventory and actual
partial outcomes; its ten source/provenance tests pass. Desktop/mobile browser
checks on the recovery report show working section controls, no horizontal
page overflow, visible unrun pilot status and the frozen assessment hash.
Subsequent report changes are limited to the final owner/test inventory paragraphs.

Resume from the current-state checklist above. Preserve Kimi checkpoint 7f2676f
and its independent rejection note. The pending user proposal remains in
inbox/codex-1-to-user_meta-protocol-change-evidence-first-efficiency_recovery-decisions.md;
there is no approval of a quorum or experimental-arm change by elapsed time.

### 2026-09-11 — continued six-priority implementation

The user renewed the instruction to finish all six priorities and resume existing
work. Integration HEAD a2e7f5d and draft PR #73 were revalidated. The existing
Hermes-to-Zcode / full-four amendment decision remains pending; unrelated work
continues with existing owners. No historical FINAL or signatures change.

Codex next works on the remaining interactive spawn-TTY launch telemetry in
previously claimed runner/launch.go, launch_test.go, handoff.go, handoff_test.go,
and app/consensus_request_signoffs.go and its tests. Preserve actual terminal
file descriptors while recording the separate process invocation; do not count
the unobserved handoff as execution or invent unobserved stream usage. Revalidate
the live rendered protocol and handoff before the process starts. Test actual
child success, failed start, timeout cleanup and required evidence failures.

Kimi's preserved 7f2676f checkpoint now has the current integration merged into
its isolated worktree without conflicts. Retry only its unfinished evidence
corrections through a newly built, uninstalled measured launcher. All ownership
boundaries remain unchanged.

Additional non-overlapping Codex allocation before edits: internal/procctl/terminal_unix.go
and terminal_windows.go for foreground terminal ownership with process-group cleanup.
The terminal launch keeps actual descriptors; it must also restore the prior foreground
group after exit and kill only the owned child group on timeout. A constant shell
no-op used to restore foreground executes no model task and is not agent telemetry.

Additional Codex allocation before editing: docs/agent-runtime-configuration.md
for the spawn-TTY prompt delivery requirement. A process with no configured
file/argument delivery cannot claim receipt of the rendered protocol. Print-only
handoffs remain user-driven; measured protocol-task spawning requires an explicit
matching prompt placeholder.

The terminal process path now creates its own invocation, re-renders current
protocol context, and requires a matching file/argument prompt delivery contract.
It preserves all three real terminal descriptors and explicitly labels stream
observations unavailable. Unix foreground ownership is restored through the
os/exec child setup, avoiding process-wide signal-handler mutation; timeouts
terminate the owned child group. Normal child/failed-start/nonzero/timeout and
required-evidence failure fixtures pass, as do the relevant signoff/handoff
checks. A compiled test binary under Python pty.fork verified actual child TTY
detection, child terminal input and restored-parent terminal input on Darwin.
An initial harness using go test directly failed because the Go test launcher
redirects the test binary's stdin; the corrected harness executes the compiled
binary under the PTY and passes. The failed observation is not discarded.
Focused race checks and vet passed before the final delivery-contract change;
current full-suite and independent owner review are still required. Windows
cross-compilation passed before that final change; Windows runtime is untested.

Full go test ./... passed on terminal-launch commit 122f4d2. Claude's independent
source review is running via measured launcher parley-recovery-4. Kimi's new
owner correction is running via parley-recovery-3 and has created its own early
handoff; neither in-progress action is claimed as completed.

Next Codex budget work (existing internal/budget/** allocation): implement the
shared durable reservation store before wiring callers. It must serialize
concurrent callers with an OS lock, precharge unique actions, retain failed or
interrupted reservations, reject replay as permission to execute again, retain
unknown monetary exposure, reject corrupt persisted state and use the tested
fsutil synchronized replacement boundary. Driver/manual/resume/BLOCK wiring and
operator extensions remain separately required; a passing store alone is not AC-B1.

The mandatory shared-volume reservation test found that flock on this mount
reports success without mutual exclusion: 23 of 24 concurrent calls passed a
five-call monetary cap, and a second process acquired a supposedly held lock.
The package is NOT accepted on that evidence. Move the kernel lock to the
host-local cache, keyed by the canonical ledger directory, and verify competing
handles are excluded before trusting it. Ledger data and synchronization remain
in the worktree. This is a same-host runtime lock, not distributed cross-host
coordination; that limitation must remain explicit in review and documentation.

The reservation-store tests now pass on ordinary local temp storage and the
actual shared worktree volume. They exercise concurrent monetary/call ceilings,
replayed IDs, failed/crashed precharges, kernel-lock release after killing a
separate process, unknown terminal cost, conflicting settlement, durable
replacement failure and malformed/duplicate/aliased state. The local-cache
lock verifies same-host exclusion with competing handles before every use.
The store's data remains on the shared volume. Same-host locking is explicit;
distributed writers are not certified. Driver/manual/launch integration is not
yet present, so this remains a foundation rather than a completed budget gate.

Claude's new owner-authored partial source review at 122f4d2 identified three
MAJOR findings: unobserved byte counters serialized as zero, lack of a real-PTY
timeout/restoration fixture, and unbounded/undocumented argument prompt delivery.
Codex accepts these as correction work, plus the concrete MINOR mode/diagnostic/
selection/timeout/private-prompt issues. No review or signature is edited.
Additional non-overlapping Codex allocation: scripts/test_terminal_launch.py,
a reproducible Unix PTY harness for the actual compiled Go terminal tests.

Terminal review dispositions (Claude reviewed 122f4d2, not the following fixes):
MAJOR-1 corrected with nullable output byte counters; unobserved handoffs and
TTY streams now persist null. MAJOR-2 corrected with the committed PTY harness
covering timeout, failed start and both evidence-write failures: all four restore
the parent's readable terminal and leave no delayed descendant artifact. The
success case also passes /dev/tty at a nonzero parent fd. MAJOR-3 corrected with
a 120 KiB per-argument guard, an explicit file-mode preference and documentation
of argv visibility. This is a conservative bound, not a universal argv capacity
claim; OS total argv/env failures still remain recorded failed starts.

MINOR-1 through MINOR-5: added interactive-mode validation; preserved restoration
error details; shared delivery validation with selection and recorded failed
handoff events; used one session/poll deadline; moved newly rendered handoff
prompts into private ignored invocation directories. Historical prompt artifacts
are preserved. SysProcAttr now preserves unrelated fields, WaitDelay is explicit,
and the argument scan exits on a match. The fixed /bin/sh restoration helper is
retained as an explicit Unix dependency; it runs no task data. Failed exits remain
failures even if an artifact exists, preserving the existing process-integrity gate.
Ctty parent-fd semantics are independently confirmed in local Go source
syscall/exec_libc.go lines 30-34 and by the nonzero-fd PTY fixture. This is Codex's
disposition and test evidence, not Claude's agreement or a final signature.

Validation after the terminal review corrections: full go test ./... PASS;
focused runner/telemetry race checks PASS; vet PASS; Windows runner cross-build
PASS (runtime untested); python3 scripts/test_terminal_launch.py PASS for all
five real-PTY scenarios. The budget-store ordinary/shared-volume/race tests and
Windows cross-build also pass. The English report was rebuilt and all ten report
tests passed; no new browser QA is claimed for this text-only progress update.
Recovery inventory now contains ten terminal attempts (three costs unknown),
with known CLI-estimate subtotal USD 21.381507; overall cost remains unknown.
Claude's review attempt completed at 601.450 seconds and estimated USD 3.8726815;
its reported model list includes Opus 5 and Haiku, so no single-model attribution
is inferred. Kimi's new attempt remains live at this checkpoint.

Live continuation handles at 2026-09-11 08:49Z (revalidate before any retry):
- Kimi invocation e020e34e-b900-44c5-9fad-70d7baac1a2b, agents-exec session
  34323, expected terminal deadline around 08:49:29Z. Its own handoff now reports
  corrected evidence/fixtures, with unsupported-entry coverage skipped on this
  shared filesystem; independent inspection/execution and commit still required.
- Claude follow-up review launched on 9cd1e42 through parley-recovery-5,
  agents-exec session 29245. It reviews budget foundation and terminal fixes;
  expected owned artifact claude-1-budget-tty-followup-20260911.md. Do not restart
  while this actual process is live.
- PR #73 remains draft, pushed through 9cd1e42. No final merge or release.

### 2026-09-11 09:00Z — independent findings and next bounded corrections

Kimi completed invocation e020e34e-b900-44c5-9fad-70d7baac1a2b at
08:49:13.899763Z (1784.296s, exit 0, cost unknown); source checkpoint 89a4305
is preserved locally. Codex independently ran the focused evidence/app suite
on local TMPDIR: 42 pass events, zero failures/skips, including unsupported
entries and the serial/barrier fixture. New overlay probes reproduce three
remaining issues in the Kimi-owned slice; details and exact commands are in
implementation-notes/codex-1-kimi-independent-20260911.md. Integration remains
pending owner correction; earlier handoffs are preserved.

Claude follow-up invocation 6e0126fc-30cb-44a7-8b69-04d6ee4f8e85 completed
at 08:57:05.137977Z (509.599s, exit 0, CLI-estimate USD 2.596291). His own
partial source review concurs with the principal TTY corrections and identifies
budget locking/recovery and residual TTY issues. No shell execution or full
acceptance is attributed to that review. Old sessions 29245 and 34323 are
terminal and must not be polled or restarted.

Codex will correct the existing budget allocation before caller wiring: make
lock identity stable or refuse divergent local-cache origins, canonicalize
absolute paths before resolving aliases, conservatively unify case aliases,
pin the no-op-lock detector with an injected regression, and cover concurrent
process cap enforcement. Bound JSON traversal and diagnose clock/overflow
failures. Add explicit append-only operator cost reconciliation without
erasing spent actions or unknown observed costs. Operator authorization must
come from an actual CLI control, never participant frontmatter. All launch/
action integration and acceptance gates remain open.

Codex next wires the explicit cost-recovery control in the already-claimed
internal/app/{budget.go,budget_test.go,app.go}: read-only inspection and an
attended `budget reconcile` command with an exact ledger/scope/action/decision,
conservative ceiling and reason. Reuse the existing platform terminal probe;
this is an operator control, not cryptographic authentication or a grant from
participant frontmatter. Add the recovery instructions to the already-claimed
docs/agent-runtime-configuration.md. Per-launch enforcement remains separate.

Codex's budget and residual terminal corrections now have executed evidence and
explicit dispositions in implementation-notes/codex-1-budget-review-dispositions-20260911.md.
The independent reviewer has not accepted those corrections yet. The earlier
full-suite pass predates the later operator-control/Inspect additions; those
additions separately pass focused app/budget tests. Kimi is running only the
new probe correction under invocation a4e1d24c-6466-4764-96c9-ae83b3e7ec85
(exec session 76055, started 09:00:25Z, 30-minute deadline). Do not restart it
while live. The terminal inventory now contains twelve real attempts, four
unknown costs, known CLI-estimate subtotal USD 23.977798, total cost unknown.

Claude completed review f47b7f16-bfdd-4ca7-99af-2b0d70d020f8 at 09:28:48Z
(655.093s, exit 0, CLI estimate USD 4.5135205). The own artifact confirms most
corrections and retains three new major findings: cache-lock identity loss,
Inspect origin mutation, and unsupported attended recovery platforms. Codex
will fix the existing owned budget/CLI paths before integration: bind a random
per-lock identity and refuse a missing lock/origin on an existing ledger, make
Inspect an atomic read without writes, and implement a Windows console probe.
Additional non-overlapping allocation before edits: internal/app/budget_attended_windows.go
and budget_attended_other.go for this budget-specific control; global protocol
publication's terminal policy is unchanged.


### Budget identity correction checkpoint — 2026-09-11

Implementation dispositions and executed checks are recorded in
implementation-notes/codex-1-budget-identity-dispositions-20260911.md. Claude's
source review remains unchanged; these changes still require independent review.
The full app suite following the previous operator CLI additions passed (66.253s,
budget-control-app-suite-20260911.log). New focused checks are scoped above.
The recovery inventory now contains 14 unique terminal attempts, five unknown
costs, known CLI-estimate subtotal USD 28.4913185 and unknown total cost.
Kimi is working on the new independently reproduced package/build-failure PASS
under measured run recovery-kimi-package-failure-20260911 (exec session 40766).
No new experiment or quorum amendment is inferred from that correction run.


### Recovery invocation correction — 2026-09-11

The first package-failure follow-up did not start useful work: invocation
fefaab07-cb0a-4f40-9f3f-7f8bdba2b5a1 ended at 09:41:52Z after 0.859s,
exit 1, with "storage write failed: permission denied" and no artifact. It used
the wrong local Kimi home/sandbox pairing. The failed attempt is retained, not
relabelled or removed. Retry f750e9ad-b977-44c4-95f3-cc9ce51a77ee uses the
original recovery KIMI_CODE_HOME and participant-recovery.sb, started 09:56:00Z
(exec 55171, 30-minute ceiling). Claude source review f77101c9-71e8-4252-b7f4-5f661f73b2f4
started 09:53:16Z against 3ea8693 (exec 6497, native restricted tools only).
The report inventory now has 15 terminal attempts, six unknown costs, known
CLI-estimate subtotal USD 28.4913185 and unknown total cost. Report build/tests
pass (10 tests); current refresh is status-only, without new browser QA.
Next actual verifier integration call-site requirements are recorded in
implementation-notes/codex-1-evidence-close-callsite-plan-20260911.md.

Additional source concern for the next budget disposition: Windows currently
locks byte 0 while identity/probe reads include that byte. LockFileEx's mandatory
range exclusion may prevent a second handle from reading identity while held,
including the self-probe; cross-compilation cannot establish runtime correctness.
Evaluate a separate lock range beyond the bounded identity payload and retain
Windows runtime status as untested. The current review source is left stable
at 3ea8693 while that independent reviewer runs; no closure claim is made.


### Existing Hermes owner continuation — 2026-09-11

Read-only discovery confirms the original Hermes CLI remains installed with saved
model fireworks/inkling and high reasoning. Pending replacement is not approval;
the historical idea still has hermes-1. A bounded file-tools-only continuation
asks that existing owner to finish readiness schema/capture/diagnostic fixes in
its original worktree, with no new membership or signatures. Runtime config and
state are isolated; global config/auth files are reused read-only. Only the
owner's already-claimed preflight source/test files and a new own handoff may be
written. Codex will run independent tests and serialize existing preflight glue
when integrating. This does not authorize a pilot amendment or resolve quorum.


### D-series corrections and recovered owner review — 2026-09-11

Claude's original review of 3ea8693 is now published by Claude unchanged, with
exact byte identity independently checked against the original denied Write
payload. Publication retries are separate spent attempts, not new source reviews.
Codex's dispositions are in implementation-notes/codex-1-budget-d-series-dispositions-20260911.md.
They cover Windows lock-byte separation, origin diagnostics, retained unknown-cost
reserves, bounded read-only inspection, honest unsupported recovery, and durable
BLOCK evidence even after failed execution. Focused, race, shared-volume, vet
and Windows cross-build checks pass; Windows runtime remains untested. This
is a partial implementation correction, not final independent acceptance.

Kimi's package-failure follow-up timed out with source retained and no new owner
handoff. Independent local checks on that source pass 53 selected evidence/app
test events and 85 full evidence-package test events, zero skips/failures,
including the unchanged negative overlay probes. Owner handoff recovery is in
progress; the source has not yet been integrated.

Hermes completed a partial owner handoff after an earlier turn-limit exit.
Independent compilation then found a missing strings import in his new tests.
The readiness source still contains a regex duplicate-key scan and incomplete
schema/secret/capture enforcement. A bounded owner correction is running with
the exact compiler failure and source counterexamples; no readiness acceptance
or quorum amendment is inferred.


Additional Codex allocation (2026-09-11, before edits):
internal/app/evidence_verify.go and internal/app/evidence_verify_test.go for the
runtime-only independent-verifier helper and real process closure integration
fixtures. These new paths do not overlap Kimi's evidence package or existing
app helper claims. Codex also uses its existing app.go, driver_impl.go,
internal/driver/impl.go and runner/consult.go allocations for command routing,
closure and launch identity glue. The actual verifier must invoke the helper;
facilitator-side execution with a replacement actor string is not acceptable.


### Independent verifier source checkpoint — 2026-09-11

Kimi's corrected evidence slice and participant-owned handoff are integrated
from 32d00bc via 374a5c5. The full independent evidence-package run has 85 passing
terminal test events with zero skips/failures, including unchanged adversarial
probes. The 53-event selected run overlaps it and is not an additional total.

Codex has implemented a real verifier CLI/helper path through Driver.Advance:
checks, an independently selected agent process, its actual helper execution,
retained criterion reruns, bound receipt validation and a last completion gate.
The full Go suite, scoped vet and Windows cross-build pass. Windows runtime and
a live real-model closure trial remain untested. Details, source trust boundary
and the still-open defects are in
implementation-notes/codex-1-verifier-checkpoint-20260911.md.

A newly executed stronger overlay probe fails after otherwise successful closure:
writing status=complete changes the separately hashed non-evidence document.
The passing suite did not check that final persisted invariant. Do not treat this
checkpoint as AC-E acceptance. Preserve the failing counterexample and require a
correction with an exact, verifier-bound status transition; never broaden the
scope exclusion. Concurrent report replacement also still needs a cooperative
serialization or compare-and-swap boundary, not just read/compare/rename.

The current invocation inventory has 27 unique terminal attempts, 14 unknown
costs, known CLI-estimate subtotal USD 36.6907045 and unknown total cost. The
numerical >=20 threshold is met; actual launch-surface coverage is not. Hermes's
latest owner slice does not compile and is not integrated; its unchanged own
handoffs do not substitute for an executed successful compile.


Additional Codex allocation (2026-09-11, before edits):
internal/runner/launch_budget.go and internal/runner/launch_budget_test.go for
one pre-spawn reservation and cancellation-independent terminal settlement
across runner process implementations. These new paths are non-overlapping.
The policy context is separate from replaceable LaunchInfo metadata. Durable
scope/policy resolution, legacy import and app/operator configuration remain
separate required integration work; an explicit internal context API alone
must not be described as automatic whole-product budget enforcement.


### Serialized readiness correction allocation — 2026-09-11

Hermes's bounded correction invocation de2694cd-9daa-4cb5-8849-c956778674a6
has stopped and its own handoff exists. Current owner source compiles, but
TestClassifyReadiness/bare_content_no_provenance_is_malformed fails: an
unattributed content field still becomes ready. Earlier failed compilation
and the 16-failure compile-only overlay remain in private logs.

Under parley-worktrees Sections 5/6, overlapping edits may be serialized through
the integration-branch owner. Codex records this explicit serialized override
BEFORE editing: internal/app/preflight.go, preflight_liveness.go and
preflight_liveness_test.go will receive the necessary schema/capture fixes in
integration/meta-protocol-change-evidence-first-efficiency at
../worktrees/evidence-first-integration after the owner checkpoint is preserved.
No concurrent Hermes source task remains; the owner worktree is retained.
Original allocation and participant-owned artifacts/signatures are unchanged.
This is source integration/fixup authorized by the existing completion task,
not a quorum replacement or pilot amendment. Codex also updates its already
owned preflight_evidence_test.go for the changed classifier signature and new
adversarial cases. Independent review remains required after these corrections.


Additional serialized completion correction allocation (2026-09-11, before edits):
Kimi's completion-transition handoff and source are preserved at c232e1d and
merged at 195fbed; his native call is terminal. Codex now serializes two
non-concurrent source corrections in internal/evidence/completion.go and adds
internal/evidence/completion_review_test.go: reject a malformed short binding
without panic and require the actual recomputed post-state status to be complete,
even when a modified record's AfterSHA256 matches another status. The pure status
transform must also reject YAML aliases/duplicate spellings of the status key.
This uses the integration-owner serialization rule, preserves Kimi's own note
unchanged and requires fresh independent review of the resulting source.


### Integrated verifier/readiness/process-budget checkpoint — 2026-09-11

See implementation-notes/codex-1-verifier-readiness-budget-dispositions-20260911.md
for all new review dispositions, ownership serialization, executed counterexamples
and corrections. The generated status transition now remains auditable after
completion, and driver cursor v2 prevents deletion/shrinkage of its original named
checks across ticks or a process restart. The current full Go suite, scoped vet
and Windows app cross-build pass; runner race tests pass. All claims remain
partial until independent review accepts the final source, including the open
report concurrency and recovery boundaries.

The reconciled inventory now contains 32 unique terminal attempts, 17 unknown
costs, known CLI-estimate subtotal USD 42.5625535 and unknown total cost. Original
records, failed writes and provider-reported internal model usage are retained.
No participant call remains active at this checkpoint. The exact packet trial,
full-six pilot and pending quorum/pilot amendment decisions are unchanged.


### Serialized report publication allocation — 2026-09-11

Before edits, Codex records a serialized integration override under
parley-worktrees Sections 5/6 for internal/evidence/report.go, new
internal/evidence/report_guard.go and report_guard_test.go, and
internal/app/driver_checks.go and driver_evidence.go. Kimi's canonical own
source/handoffs remain preserved and no owner source task is concurrent.
Additional non-overlapping Codex paths: internal/budget/resource_guard.go,
internal/app/evidence_publication_test.go. Existing evidence_verify.go/tests
and driver_impl.go allocations cover verifier and Complete wiring.

Plan: reuse the tested pinned host-local budget locking primitive for a
synchronization-only guard with no monetary accounting. Keep metadata in git
administration where available (per-report canonical-path identity), and
initialize a stable local runtime guard before checks in non-git fixtures.
All cooperative report writers, evidence-table writes and completion publication
share that guard. Verification publishes with byte-exact compare-and-swap and
keeps the guard through receipt persistence. A persisted driver acceptance
binds completion to the report actually reconciled against that receipt.
Completed reports cannot be overwritten by a delayed writer. No same-UID or
distributed-host authentication is claimed. Require actual process contention,
failure/recovery and post-completion tests, then independent source review.

The same serialized report correction includes internal/evidence/completion.go: return the already validated prior status alongside the existing no-op error, so publication can distinguish completed status from malformed frontmatter without weakening transition validation.


### Report publication checkpoint and next independent review — September 11

See implementation-notes/codex-1-report-publication-20260911.md for the
shared cooperative guard, byte-CAS, exact accepted-report binding and actual
persistence-failure/recovery process tests. Source remains partial. The full
Go command exited zero, but the streamed shared-volume log is damaged; the
in-memory captured rerun is pending. Focused, race, shared-volume, vet and
Windows cross-build checks pass. No Windows runtime or live-model closure claim.

A read-only source-review worktree will be frozen from this checkpoint at
../worktrees/evidence-first-publication-review, branch
review/meta-protocol-change-evidence-first-efficiency/publication-20260911.
Claude may write only its own new implementation-notes/claude-1-publication-review-20260911.md
and its native atomic siblings there. No source claim overlaps: Codex continues
app/budget work on integration, the reviewer inspects this frozen snapshot.
The historical quorum remains unchanged; this is not a full Phase-6 signoff.


### Captured full verification and budget continuity counterexample — September 11

The full captured-output rerun passes at source 0cdc56e: app 101.024s,
budget 3.221s, driver 4.791s, evidence 3.590s, runner 25.249s. The 1403-byte
local and shared logs match SHA256
f9c8e46a9162b6ae019d43e5b3743cd098a264d9b274c9a336d5de15bce10940.
Shared log: .parley-runtime/report-serialization-full-captured-20260911.log.
The damaged earlier log remains retained. Independent Claude review of this
frozen source is running in the dedicated review worktree with a 15-minute
ceiling and native USD 5 cap; its exact atomic-write sandbox canary passed.

Before attaching automatic scope/policy, Codex executed a new budget continuity
probe in the already-owned internal/budget boundary.
Command: go test -count=1 ./internal/budget -run '^TestDeletedLedgerCannotResetCharges$'
Result: FAIL, exit 1: deleting only ledger.json granted another attempt above
the frozen cap. The store retained its lock origin but recreated empty charges.
The correction publishes an immutable bounded continuity witness after durable
ledger persistence and before any grant returns. Missing charged history then
refuses; the witness cannot reconstruct charges or authenticate same-UID edits.
The source is on integration only; the reviewer snapshot stays unchanged.


Additional Codex allocation before policy wiring — September 11:
internal/budget/binding.go and binding_test.go, internal/app/budget_configure.go
and budget_configure_test.go. Existing budget.go, runner/launch_budget.go,
runner/telemetry.go and their tests cover routing and the shared pre-spawn
boundary. All files are disjoint from the frozen Claude review worktree.

Implement an explicit attended configure control that pins a policy per idea
(or auxiliary scope) in shared Git administration, then automatically resolves
that binding at every common launch boundary. A binding cannot be overwritten
by mutable participant files or a different injected context. Configuration
refuses observed legacy activity until its charges can be migrated; absent
policy is not claimed as enforcement. Existing internal explicit policies stay
available for trusted tests/callers, without overriding a persisted binding.
Initial configuration and later operator extensions are distinct: this step
does not silently implement a reset/extension or authorize experiment funding.


### Publication review follow-up allocation — September 11

Claude's own frozen-source review is preserved unchanged at fbb3c31. Its
new triple-deletion path and CRLF/nested-status cases are checked through
actual tests, including Driver.Advance rather than source assertion alone.
Codex serializes the following evidence-boundary follow-up after all owner
calls are terminal: new internal/evidence/contract_pin.go; existing
report_guard.go, completion.go and completion_review_test.go; app driver_checks.go,
evidence_publication_test.go, driver_impl.go and driver/checks.go, driver.go.
The runtime original-contract witness lives with the publication guard in
Git administration and survives removal of cursor, report and current checks.
It is monotonic cooperative state, not protection against deleting all Git
administration or malicious same-UID tampering. Preserve exact CRLF bytes
while accepting one top-level normalized status; nested fields stay ordinary
bound content. No participant canonical review or signature is edited.

Additional serialized correction before editing completion_test.go: Codex
updates Kimi's old CRLF-refusal fixture to the accepted byte-preserving CRLF
behavior, keeping inverse-transition coverage. This source-test overlap is
serialized under the existing integration ownership, not a rewrite of Kimi's
canonical review/handoff. No participant source writer is concurrent.

### Original-contract continuity and configured launch budgets — September 11

See implementation-notes/codex-1-publication-review-dispositions-20260911.md
for every finding in Claude's 0cdc56e review, corrected triple deletion and
CRLF handling, ledger continuity, actual persisted launch policy attachment,
historical-accounting counterexamples, and exact validation logs. The full Go
suite passes (app 92.263s, budget 4.117s, driver 5.357s, evidence 3.038s,
runner 26.011s). Budget/evidence/runner race, scoped vet, Windows cross-build
and isolated shared-volume runtime fixtures pass. Failed earlier runs remain
retained. Windows runtime and final independent acceptance are unverified.

The inventory now reconciles 33 unique terminal attempts, 17 unknown costs
and USD 45.7335065 known CLI estimates. Total cost remains unknown. The added
Claude attempt cf190d54-548d-4840-a5b6-b10f4eb09742 reports Opus 5 plus Haiku
4.5, not an extra participant; its own source review is preserved at fbb3c31.

Next read-only source-review allocation: ../worktrees/evidence-first-binding-review,
branch review/meta-protocol-change-evidence-first-efficiency/binding-20260911,
frozen from this tested checkpoint. Claude may write only its own new
implementation-notes/claude-1-binding-review-20260911.md and its native atomic
siblings there. This is a scoped source review, not final acceptance. Codex
continues independent implementation on integration; no source files overlap
with the reviewer. Historical quorum, full-six and exact packet gates remain.

### Review stopped and next driver-budget counterexample — September 11

The new Claude source-review invocation bdeecc5d-0ce0-456d-a363-af62f6c3742f
actually terminated at the provider weekly limit (6.292s, exit 1, CLI estimate
USD 0.455252); no canonical review file was produced. Provider reset is reported
as September 14 06:00 Europe/Berlin. The inbox contains the concrete pending
recovery choice. No participant process remains active, and no retry or roster
change is inferred. Inventory now has 34 actual terminal attempts, 17 unknown
costs, known CLI-estimate subtotal USD 46.1887585 and unknown total.

Codex also reproduced the still-open persistent driver-step gap using a Go
overlay against production Driver.Run: MaxDriverSteps=1 permits round 2 on the
first Run and round 3 after resume, observed calls [2 3]. See
implementation-notes/codex-1-driver-resume-counterexample-20260911.md for exact
reproduction and required next boundary. This negative probe is separate from
the passing ordinary suite and confirms that D6/AC-B1 is not yet complete.

### Persistent driver-step allocation before edits — September 11

Codex claims new internal/budget/step_binding.go, step_binding_test.go,
step_session.go and step_session_test.go; new internal/driver/budget.go and
budget_test.go; new internal/runner/step_budget.go and step_budget_test.go.
Existing serialized claims cover driver.go, loop.go, impl.go, runner.go and
phase58.go. Add internal/driver/consensus.go for the shared pre-mutation
boundary; it has no concurrent owner writer. All work stays on integration.
The frozen Claude review tree at 2194dca stays unchanged and its invocation
is terminal. No participant's own review/handoff is edited.

Implement frozen lifetime step/time policy with the existing durable store,
one precharged step session per actual driver transition, shared by nested
runner calls. Read-only/await paths do not reserve steps; failed work retains
its charge. Direct Advance, resumed/new runs, BLOCK and supported manual
protocol runners must consult the same idea scope. Existing driver history
cannot silently acquire zero charges. This is the step/time slice; fixup and
cross-review counters, explicit legacy migration/extensions, monetary policy
mapping and regression trajectory remain separately required D6 work.

Continuation allocation: add internal/budget/step_history.go and its tests for
validated run identity and bounded history reads. Codex also extends the existing
internal/runner/launch_budget_test.go process fixture under the same serialized
integration ownership. Refused manual attempts must reach the already-instrumented
launch boundary; no high-level grouping shortcut may erase request/terminal
evidence. These edits do not change the frozen independent review checkpoint.

Snapshot-boundary continuation allocation before edits: internal/runner/telemetry.go
is added to this serialized slice. A Phase-6 disposable shared clone must retain
the live origin for budget scope and invocation evidence; clone cleanup cannot
erase telemetry or create an unbudgeted identity. Reproduce both step-policy
refusal and launch-policy bypass through the actual review snapshot fixture.

Shared-volume follow-up allocation before edits: internal/runner/reviewsnapshot.go
and a new reviewsnapshot_test.go are serialized to Codex. The executed fixture
reached real review publication, where Darwin File.Sync returned ENOTTY. Reuse
the existing strict fsutil.SyncFile fallback and synced atomic replacement;
retain snapshot recovery and propagate every actual synchronization failure.

### Persistent steps and disposable-review origin — September 11

The driver-resume counterexample is corrected by frozen lifetime step/time
policy and durable precharged transition grouping. Run reporting now uses the
saved count and clock; unrelated legacy cursors use validated run identity.
Failed requests lacking a cursor also require migration. Manual refusal still
produces requested/terminal telemetry at the common launch boundary.

The actual review snapshot fixture additionally exposed lost origin accounting:
step policy wrongly denied the first clone review, while launch policy granted
a second call. Private runtime lineage now retains the live policy/evidence
root across snapshot execution and cleanup. Shared-volume execution exposed
ENOTTY during review artifact publication; the existing strict synchronization
helper and synced replacement correct it while preserving failure recovery.
All negative logs and own dispositions are in implementation-notes/
codex-1-persistent-step-and-snapshot-dispositions-20260911.md.

Final full Go suite passes: app 97.897s, budget 5.165s, driver 6.543s, evidence
1.798s, runner 29.673s; captured log SHA256
e1bf602e6a365fbe61f6b517913f472f020715d0514f51d6d3f3c759c8bd6c27.
Budget/driver race pass (13.154s/11.095s); corrected runner race passes in
54.408s. Scoped vet, Windows app cross-build and actual isolated shared-volume
fixtures pass. Windows runtime and fresh independent acceptance are unverified.

This is a tested step/time and snapshot-accounting slice. D6 still requires
persistent fixup/cross-review counters, legacy migration, operator extensions,
monetary default mapping, canonical refusal evidence and regression trajectory.
All exact packet/pilot, real-model verifier and owned-signature gates remain.
No new participant call occurred: inventory remains 34 terminal attempts,
17 unknown costs, USD 46.1887585 known CLI estimates and unknown total.

### Finite operator cycle-extension allocation before edits — September 11

Codex claims internal/budget/cycle_extension.go and cycle_extension_test.go,
internal/app/budget_cycle.go and budget_cycle_test.go, plus the existing
cycle binding/session, driver cycle/budget/consensus/impl and runner cycle files
and their tests under serialized integration ownership. internal/app/budget.go
and docs/agent-runtime-configuration.md join this boundary. No other source
writer is active; independent participant review worktrees remain frozen.

Add an attended explicit operator control with a read-only preview, expected
policy digest, unique decision identity, finite absolute maximum and reason.
Retain the immutable original ceiling and every grant in the atomically saved
policy; keep counts, reservation identities and original clock unchanged.
Exact decision replay is idempotent, conflicts/stale previews refuse, and a
failed publication does not imply permission to retry an action. Driver and
typed manual calls must actually use the recorded grant, preserving track
contradictions, skipped phases, review independence and all closure gates.
Implementing this control grants no extension to the active real idea. Legacy
migration, launch/step extensions and the other D6 obligations remain separate
required continuation work.

### Shared fixup/cross-review allocation before edits — September 11

Codex owns new internal/budget/cycle_binding.go, cycle_session.go,
cycle_history.go and their tests; new internal/runner/cycle_budget.go and tests;
and new internal/driver/cycle_budget.go and tests. Existing serialized ownership
also covers driver.go, budget.go, impl.go, consensus.go, runner.go, phase58.go,
launch_budget.go, and the budget process fixtures. No other source writer or
participant model process is active. Frozen review worktrees stay unchanged.

Use one persistent policy/ledger per idea and charged cycle kind, shared across
manual, driver, resume and BLOCK paths. Preserve current known cursor/marker
counts as an explicit frozen carried count; ambiguous legacy attempts require
migration, not a zero-count policy. A synchronous operation may reuse its own
reservation in nested runner calls, never across a retry/new operation. The
inclusive final fixup remains reviewable; signoffs and zero-fix reviews do not
spend a fixup. Preserve current track caps and all independent closure gates.

Documentation allocation before edits: Codex extends the serialized boundary
to docs/agent-runtime-configuration.md and its own cycle-disposition note below.
No participant review file or signature is changed.

### Persistent protocol cycle checkpoint — September 11

Source a223a0d persists separate fixup and Phase-2 cross-review policies and
charged ledgers in the shared idea scope. Cursor deletion, a new driver run,
ordinary and BLOCK dispatch, and supported typed manual runners cannot reset
failed charges. Nested children of one synchronous operation share one charge;
independent operations spend separate charges. The existing track cells remain
unchanged. The fifth deliberation fixup and third cross-review group are allowed;
the sixth/fourth are refused. Verification after the final allowed fixup remains
possible, while an over-cap state or changed frozen cap cannot close.

Unknown historical phases and ambiguous histories require reconciliation.
Missing track authority cannot activate a new legacy policy. A standalone
unobserved handoff creates evidence without activating a cycle policy, and a
cancelled launch retains its cancellation terminal without spawning or charging.

Full Go tests pass (app 104.956s, budget 11.712s, driver 11.907s,
runner 31.677s); full log SHA256
f5cbd12d7da2b09eeac900fe5cc641f18c6fffdbae3dd012ee20b61010071e7d.
Scoped race passes (budget 15.027s, driver 13.494s, runner 58.611s), as do
vet, Windows app cross-build and isolated shared-volume cycle fixtures.
Negative probes, exact commands, source manifest, hashes and limits are in
implementation-notes/codex-1-persistent-cycle-dispositions-20260911.md.

This is a tested persistent-cycle slice, not all of D6/AC-B1/B2. Explicit legacy
migration, operator extensions, durable semantic action replay, monetary default
mapping, canonical refusal publication/recovery and independently confirmed
patch-regression trajectory remain required. Unknown external/manual work
cannot be inferred from arbitrary prompt text. Fresh independent review,
actual model concurrency/closure, the exact packet trial, full-six comparative
pilot, owned final reviews/signatures and elapsed follow-ups remain pending.
No new model invocation occurred; inventory remains 34 terminal attempts,
17 unknown costs, USD 46.1887585 known CLI estimates and unknown total.


### Finite cycle-extension checkpoint — September 11

Tested source 321aef9 adds read-only `budget cycle inspect` and attended explicit
`budget cycle extend`. A finite absolute grant retains original ceiling, carried
count, all decision history, charged ledger bytes and activation clock. Stale
previews/conflicting replay refuse; exact decision replay does not grant again.
Driver, manual and BLOCK paths actually use the recorded extension. Skipped
phases, track configuration, strict review and independent evidence gates remain.

An existing full-suite diagnostic test caught an omitted refused-attempt ordinal;
it was corrected without weakening the assertion. Final whole Go suite passes
(app 88.994s, budget 12.375s, driver 9.884s, runner 30.044s). Budget/runner race
passed at unchanged source; corrected driver race passes in 12.294s. Vet,
Windows app cross-build, isolated shared-volume fixtures and a real unattended
CLI refusal all pass. Windows runtime and fresh independent acceptance remain
unverified. No real idea received an operator grant and no model call occurred.

Exact validation, failure logs, schema limits and remaining scope are in
implementation-notes/codex-1-cycle-extension-dispositions-20260911.md. Legacy
migration, launch/step extensions, safe recovery, durable semantic replay,
monetary mapping, canonical refusal evidence and regression trajectory remain
required, alongside all outstanding independent/live experiment delivery gates.

### Cycle-extension publication allocation — September 11

Codex owns the serialized documentation/report refresh for tested source 321aef9:
the existing runtime guide and owned disposition, evaluation workspace
delivery/2026-09-05/{PROGRESS.md,priorities.json,report.html}, and the draft PR
description. Check retained validation hashes, rebuild the report and run its
existing tests before publication. Prior browser evidence describes the previous
report hash; a content refresh does not inherit exact-file browser verification.

### Monetary default mapping allocation — September 11

Codex claims the serialized monetary mapping boundary before edits:
internal/budget/monetary_binding{,_test}.go, internal/runner/launch_budget.go,
internal/runner/monetary_defaults_test.go, internal/driver/budget.go and its
owned tests, internal/config/runtime.go and loop_defaults_test.go, and the
runtime guide. The current loop-only MaxCostUSD check is insufficient for a
first manual launch or a new run. Reproduce that with actual local process
fixtures, then require a matching persistent finite launch policy with an
explicit conservative reservation before honoring a nonzero configured dollar
ceiling. Do not invent a per-call price or auto-configure a real grant.

Resolve layered defaults at the common launch boundary using the captured live
origin, reject unreadable/invalid required configuration, and retain terminal
refusal evidence. Driver-supplied monetary limits must use the same requirement.
Explicit zero defaults do not remove a saved policy. Convert configured dollar
ceilings conservatively to integer microdollars; sub-microdollar, nonfinite,
negative and overflowing inputs refuse rather than becoming unlimited.
Validate retries, omitted defaults, shared worktrees, stricter/different frozen
caps, handoff/no-spawn behavior and local child execution. This does not migrate
legacy history, issue an extension, or establish independent acceptance.
Additional test allocation: extend internal/runner/launch_budget_test.go's
existing six-surface process fixture to run with a configured monetary default,
an explicit conservative reservation and no launch-count cap.

### Monetary default enforcement checkpoint — September 11

Source 78a7b13 connects configured max_cost_usd to the shared process launch
boundary. A local-process counterexample at 0a872c5 returned success for an
unreserved manual launch despite a USD 1 default; it now records budget_refused
without starting a child. Positive defaults require a matching persistent finite
policy and explicit conservative reservation. Missing/invalid required config
does not become unlimited. Zero/omitted defaults preserve saved policies.

All six existing process surfaces run the new monetary fixture. New runs retain
known or conservatively bounded unknown exposure, and disposable review clones
resolve defaults at their live origin. Dollar ceilings convert conservatively;
no provider price or real operator grant is inferred. Fresh configuration after
retained refusal history can still require the pending legacy migration.

Full Go suite PASS (app 102.204s, budget 10.856s, driver 11.656s,
runner 32.633s; wall 104.835s). Budget/config/driver/runner race PASS
(14.247s/1.692s/10.728s/60.329s), scoped vet, Windows app cross-build and actual
shared-volume fixtures PASS. Windows runtime and independent acceptance remain
unverified. Exact commands, source manifest, negative probe and log hashes are
in implementation-notes/codex-1-monetary-default-dispositions-20260911.md.
No participant model was called. All remaining D6 and experiment obligations
listed above still apply; this is not all-six-priorities completion.

### Launch/step extension allocation — September 11

Codex claims the serialized D6 extension implementation before edits:
internal/budget/{policy_extension.go,policy_extension_test.go,binding.go,
step_binding.go,step_session.go,monetary_binding.go},
internal/runner/{launch_budget.go,policy_extension_test.go},
internal/driver/{loop.go,budget.go,budget_test.go},
internal/app/{budget.go,budget_policy.go,budget_policy_test.go}, runtime docs and
the owned extension disposition. Existing participant artifacts are unchanged.

Add read-only inspection and attended explicit finite absolute launch/step
extensions. Preserve the original ceilings, reservation amount, charged ledger,
activation time and decision history. At least one finite ceiling must increase;
unlimited axes stay unchanged, no grant lowers another ceiling, and increased
axes must exceed actual spent exposure/elapsed time. Unknown monetary exposure
requires reconciliation before increasing a monetary ceiling. Exact replay is
idempotent; stale/conflicting decisions and malformed history refuse.

Serialize extension and reservation through the existing resource guard and
reload cached bindings. Driver/manual/BLOCK and nested sessions must consume
effective limits without double-charging or relaxing protocol/independent gates.
Original configuration remains a valid authority reference after a recorded grant;
unrecorded intermediate values remain invalid. Test publication recovery, races,
continued execution, nested grants, clock and schema boundaries, actual CLI
attendance refusal and isolated shared-volume behavior. No real policy extension
is authorized by this implementation allocation.

### Evidence cancellation race allocation — September 11

Codex claims serialized integration edits to internal/evidence/execute.go and
the new owned internal/evidence/execute_cancellation_test.go before edits.
The expanded budget/driver/runner/app race run retained in runtime-extension
validation fails: RunCriterion's context watcher reads spawned concurrently
with its post-Start assignment (execute.go:77/90). Preserve that negative run
and Kimi's original source attribution, handoffs and review artifacts.

Use the live command's established process identity for cancellation without
unsynchronized post-start state. Verify cancellation before launch, start failure,
and cancellation of an actual child group, including a TERM-resistant descendant
holding output pipes. Re-run the original app barrier fixture under the race
detector, then the required full and scoped race suites. This correction and
its execution record do not constitute participant-owned independent acceptance.

### Runtime extension checkpoint and report allocation — September 11

Source 9c1990c adds read-only launch/step inspection and finite attended absolute
extensions. Original/effective ceilings, per-launch reservation, charged ledger
bytes, activation time and ordered decisions are retained. Exact replay remains
idempotent after later decisions. Unrecorded settings, stale/conflicting grants,
lost charges, clock regression and malformed history refuse. Extensions serialize
with reservations and settlements; nested children remain one charged step.
The driver reads persistent monetary exposure separately from observed usage.
No real policy was extended and unknown observed costs remain unknown.

Expanded race testing found an existing RunCriterion cancellation race. A new
actual-process fixture reproduced it before the fix, including a TERM-resistant
descendant holding output pipes. Cancellation now uses exec's established process
identity, without the unsynchronized post-Start assignment. Both the new fixture
and the original serial/barrier fixture pass under the race detector. Original
Kimi source/handoffs and independent review artifacts remain preserved.

Full Go suite PASS (app 98.824s, budget 18.914s, driver 15.266s, evidence 6.118s,
runner 35.887s; wall 102.638s). Budget/driver/runner/evidence/app race PASS
(22.138s/17.629s/65.175s/6.517s/122.273s; wall 123.291s). Scoped vet, Windows
app cross-build and actual shared-volume cancellation fixtures PASS. All 323
Go/module files match the tested manifest. Thirty-eight logs/manifests, including
both negative race runs, are retained under the ignored runtime-extension
validation directory. Windows runtime and independent acceptance remain unverified.

Codex claims the owned runtime-extension disposition and the current evaluation
delivery/2026-09-05/{PROGRESS.md,priorities.json,report.html} update before edits,
followed by draft PR #73 description synchronization. Rebuild and run the ten
report checks. No fresh exact-report browser QA is inferred from older evidence.
All remaining D6, real-model, packet, full-six pilot, blind-grading, owned-signature
and actual elapsed follow-up requirements remain binding; no full-goal completion.

### Expanded report layout correction allocation — September 11

Before edits, Codex claims the evaluation template templates/delivery-report.html
and delivery/2026-09-05/browser-runtime-extension-20260911.json plus its owned
browser-qa evidence directory. Actual ego-browser QA of report hash 85e81b4 found
that expanding Technical progress log at width 390 creates page width 541:
unbroken plain-text hashes/URLs overflow paragraphs in #progress.doc. Retain the
negative observation and preceding harness mistakes. Apply the existing document
wrapping behavior to progress prose, rebuild, run the ten report checks, and
verify the changed report's expanded/collapsed sections and interactions at the
four viewport sizes. This is presentation correction, not live experiment proof.

### Verified report checkpoint — September 11

Report 0b9e72687fdd50e19cb74f5938d5d9849d11cd837b39a027b4d092fc73456e43
(609331 bytes; built 2026-09-11T16:00:05.838108+00:00) passes all ten report tests
and 62 retained ego-browser assertions across desktop, low-height and both mobile
widths. Expanded progress prose now fits; current source and the failed race
record are accessible through its summary control. Search/reset, pagination,
keyboard focus, five documents, painted charts and print-media visibility pass.
The owned runtime-extension disposition identifies the exact QA sidecar and its
limits. Earlier harness mistakes and the actual overflow are retained separately.
No final populated experiment result, independent acceptance or physical print
verification is inferred from this presentation checkpoint.

### Legacy launch accounting migration allocation — September 11

Codex claims internal/budget/{launch_migration.go,launch_migration_test.go,
launch_migration_history.go,launch_migration_history_test.go,
binding.go,policy_extension.go}, internal/app/{budget.go,budget_migrate.go,
budget_migrate_test.go}, the owned runner migration fixture and runtime guide
before edits. Add a read-only historical inventory and explicit attended launch
accounting import with a frozen history hash, decision identity/reason, original
accounting epoch, all future policy ceilings and an explicit additional legacy
launch count. Require quiescent writers; never infer user approval from terminal
allocation. No real migration or grant is authorized by this source allocation.

Import unique recognized attempts with original timestamps and observed cost
provenance. Preserve unknown monetary values. Retain proven pre-start refusals
without inventing a model execution or silently dropping their evidence. Unobserved
handoffs remain potentially spent; missing/inconsistent/nonterminal records,
conflicting copies, malformed state and stale inventory refuse. Inventory all
available worktrees and relevant artifacts; the operator must account explicitly
for work preceding invocation telemetry. Do not derive that count from prose.

Publish a hash-bound immutable import record, ledger and continuity witness before
the policy can authorize work. Replay must recover partial publication without
erasing later charges or grants. Existing configured/charged state cannot be
overwritten by migration; lock-origin recovery is a separate obligation. Test
known/unknown costs, refusal-only bootstrap, old failed attempts, shared worktrees,
concurrent activation, stale input, publication recovery and actual resumed
processes. This is the launch-accounting part of the broader D6 migration work;
step/cycle migration and all other remaining requirements stay open.

### Concurrent historical-read refusal allocation — September 11

The first full migration suite failed only because the existing eight-contender
cycle-extension test recognized a stale policy hash but not the bounded reader's
safe rejection of an atomic policy replacement during its initial read. Its
final state assertions still verified one winner, one grant and retained spend.
Codex claims internal/budget/{step_history.go,cycle_extension_test.go} before
edits to expose that existing changed-file outcome as a typed sentinel and let
the contention test recognize exactly that safe refusal. Preserve both bounded
reader checks and the one-winner/count assertions; do not accept arbitrary errors
or retry until a passing result hides the failed run. Retain the original full
log and rerun meaningful focused and whole-suite validation after this correction.

### Launch migration disposition and delivery allocation — September 11

Before edits, Codex claims its launch-migration disposition under
implementation-notes/codex-1-launch-migration-dispositions-20260911.md, the existing
evaluation delivery/2026-09-05/{PROGRESS.md,priorities.json,report.html} content
refresh, and browser-launch-migration-20260911.json with its owned browser-qa
subdirectory. Update the existing draft PR description after tested source is
committed. The report must keep every priority partial/preparation, retain failed
validation history and identify the exact source/report it describes. Browser QA
uses ego-browser only and does not certify an unexecuted experiment or treatment.


### Launch accounting migration checkpoint — September 11

Source fd3e3ce68080a0b7b39d168e153af0280221f94e adds explicit attended historical launch import, including
retained pre-start refusals without invented model charges. Unique failed/started
attempts, original timestamps, unknown costs and CLI-estimate provenance survive.
Exact replay recovers unchanged partial publication and retains subsequent spend,
reconciliation and grants. Missing/corrupt/aliased/nonterminal/stale history,
conflicting event/cursor identities and implicit missing zero fields refuse.
The owned launch-migration disposition records all behavior and limitations.

Full Go suite PASS (93.717s); complete scoped race PASS
(107.389s); focused adversarial/race, scoped vet, Windows app
cross-build, isolated shared-volume publication/worktree/actual-process fixtures
and the real unattended CLI refusal PASS. All 330 Go/module files match the
retained manifest. Initial missing-field/identity failures, the full-suite test's
unrecognized safe read refusal, and the mistakenly nested shared-volume fixture
failure remain retained. No production check was weakened to hide those outcomes.

No real operator migration/grant or participant invocation occurred. Independent
acceptance and Windows runtime remain unverified. Step/cycle migration, guard
recovery, semantic replay, canonical refusal recovery, independent regression
trajectory and the full signed real-model/packet/pilot/grading/signature/elapsed
follow-up requirements remain open. The goal is not complete.


### Current report verification checkpoint — September 11

The 613084-byte report 306c99b7f1cac02452ab724fb999076db00da3f9b82a937a6e40d234f6745b66
passes ten report tests and 55 live ego-browser DOM/layout/canvas/interaction
assertions across desktop, low-height and both mobile widths. The expanded log
includes current source fd3e3ce and retained failures without page overflow.
Its exact sidecar and raw/harness evidence identify the verified file. New
screenshots timed out; physical click placement and fresh visual screenshot
review are not certified. The owned disposition records these limits explicitly.
Report state stays partial/preparation and the all-six goal remains incomplete.

### Historical step/cycle accounting migration allocation — September 11

The preceding goal turn made verified progress (source fd3e3ce, checkpoint
85c8ad9, current draft PR and report evidence). Codex now claims
internal/budget/{protocol_migration.go,protocol_migration_history.go,
protocol_migration_test.go,protocol_migration_history_test.go,step_binding.go,
cycle_binding.go,cycle_extension.go,policy_extension.go},
internal/app/{budget_migrate.go,budget_migrate_protocol.go,
budget_migrate_protocol_test.go}, owned driver/runner protocol-migration fixtures,
and docs/agent-runtime-configuration.md before edits.

Add explicit attended import for driver-step, fixup and cross-review accounting,
retaining original epoch, source hashes and known lower bounds from driver
counters/events and structural cycle markers. Nested participant launches are
not one protocol charge each. Legacy telemetry lacks durable operation grouping;
therefore expose ambiguous invocations and require an explicit operator total
including failed/partially observed attempts. The observed floor is not a claim
of complete lifetime accounting, and a pre-start launch refusal alone cannot
prove no protocol step/cycle was charged before it. Never infer totals from prose.

Use an immutable import record, ledger continuity, a policy hash reference and
final activation marker; preserve exact partial replay and later charges/grants.
Existing policies/ledgers cannot be replaced and real track/runtime gates stay
intact. Import is not semantic action replay or lock-origin recovery. Keep this
implementation serialized under Codex; no real operator grant/migration or
participant invocation is authorized by this allocation. Test refusal-only
history, known/unknown grouping, stale/malformed/conflicting copies, epoch/count
floors, all publication failures, exact replay, and actual runner/driver resume
at inclusive caps. The signed full objective and independent gates stay open.

Codex also claims internal/budget/cycle_session.go before editing its cached
session path: subsequent children must check imported charge continuity, not only
the aggregate current count. This does not add a second charge for a nested child.

Codex claims internal/app/budget_migrate_test.go before updating the old negative
inspect case: `--kind step` is now a supported read-only operation, so the
unknown-kind rejection fixture must use an actually unsupported name. Existing
launch-only apply calls with mixed protocol flags must still refuse.

### Protocol migration disposition and delivery allocation — September 11

Before edits, Codex claims implementation-notes/codex-1-protocol-migration-dispositions-20260911.md,
the evaluation delivery/2026-09-05/{PROGRESS.md,priorities.json,report.html}
refresh, and browser-protocol-migration-20260911.json with its owned browser-qa
subdirectory. Record current-source validation and the retained negative tests,
then synchronize draft PR #73. Previous launch-migration test/report results do
not certify the new step/cycle import. Every unfinished full-scope requirement
and unanswered recovery decision remains open.


### Protocol accounting migration checkpoint — September 11

The step/fixup/cross-review import now preserves explicit lifetime totals and
original epochs through partial publication, exact replay, later charges and
finite grants. Distinct published steps add across runs without charging nested
children separately; fixup recovery events do not invent fresh code-writing
attempts. Disjoint marker paths across worktrees remain counted. The owned
protocol-migration disposition records source-bound behavior, negative tests,
actual-process evidence, track compatibility and remaining recovery limits.

Full Go suite PASS (134.046s); budget/driver/runner/evidence/app race PASS
(147.141s); scoped vet and Windows amd64 app cross-build PASS. Isolated
shared-volume publication/worktree/actual-process fixtures PASS (13.366s).
Actual compiled CLI refuses unattended imports for all three kinds with exit 2
and no created state. All 338 Go/module files match the retained source manifest.
The negative history regressions and initial driver fixture's mistaken refund
expectation are retained; runtime precharge behavior was preserved.

No real policy was imported/extended and no participant/model was invoked.
Independent acceptance, Windows runtime, guard/changed-import recovery, semantic
replay, canonical refusal recovery, independent regression trajectory and every
full-scope model/packet/full-six pilot/grading/signature/elapsed-follow-up
obligation remain pending. This checkpoint does not complete the six-part goal.

Source 1a3d512e8909af3c82e418c04d15aae646675ced is the committed implementation
matching that 338-file manifest. The 617000-byte current English offline report
c2cbdfb81d5c77641bf7294ca72e0aecab7ce10f75c68d5fbc355487c59b7b07 passes
ten report tests and 55 live ego-browser assertions at desktop, low-height and
both mobile widths. The owned disposition and exact-file sidecar record the
current progress, navigation/filter/pagination/keyboard/painted-chart/print-media
checks and their screenshot/physical-print limitations. Model experiments remain
not-run and independent acceptance remains pending.


### Durable verification refusal allocation — September 11

The immediately preceding goal turn only recalled the audit, so it made no
implementation progress. Current clean source and draft PR #73 remain 7b3b084.
Codex now claims internal/evidence/{refusal.go,refusal_test.go},
internal/app/{evidence_refusals.go,evidence_refusals_test.go,evidence_verify.go,
evidence_verify_test.go,evidence_publication_test.go}, the runtime guide and
implementation-notes/codex-1-refusal-dispositions-20260911.md before edits.

Retain bounded typed failure observations in Git administration before guarded
canonical publication. Driver and helper observations share verification identity
when a frozen request exists, without inventing model invocations. Preserve exact
immutable records, missing/failed publication, and explicit hash-bound recovery.
Commit only the new refusal path and verify its committed bytes. Refusals stay
inside the tested tree digest: publication requires fresh checks and never grants
acceptance. Do not serialize raw diagnostics, commands, output or credentials.
Test actual helper failures, parent interruption, publication/commit failures,
concurrent/exact/conflicting recovery, staged-file isolation, ignored-runtime
prerequisites and false-closure rejection. Historical quorum, independent review,
real-model trials and the full six-part objective remain unchanged and open.


### Refusal delivery allocation — September 11

Before edits, Codex claims the evaluation delivery/2026-09-05/{PROGRESS.md,
priorities.json,report.html,report-manifest.json} refresh and browser-refusal-20260911.json with its
owned browser-qa/refusal-20260911 directory. Record current-source validation,
retained initial path-alias and test-instrumentation failures, and explicit
limitations. Update the existing draft PR after committing verified source.
The five presentation priorities (six audit areas when context and budgets are
separated) remain partial/preparation. No model trial, independent acceptance or
closed gate may be inferred from local process fixtures or report rendering.


### Durable refusal source checkpoint — September 11

Bounded driver/helper refusal observations now have guarded canonical publication,
Git byte verification and exact-hash recovery. Pending observations survive a
stopped parent or failed publication; canonical records remain in the tested tree
and require fresh checks. The owned refusal disposition records concurrency,
interrupted staging, conflict and staged-file isolation behavior and limitations.

All 342 Go/module files match the manifest. Full Go suite PASS
(113.824s); evidence/app/driver/runner race PASS (130.440s);
scoped vet, Windows amd64 app cross-build, isolated shared-volume suites and actual
compiled CLI inspect/refuse/recover/replay probes PASS. Initial path-alias and
test-instrumentation failures remain retained. No production lock was weakened.
Independent acceptance, Windows runtime and all remaining full-scope obligations
in the disposition are still pending. No model or real accounting action occurred.


### Refusal report verification checkpoint — September 11

Source 9cef3fc9a90c8b3025796398dd06d8ef4d3b3b3c matches the retained 342-file
source manifest. The current 620834-byte English offline report has SHA256
42f47e73764d063b4e246a971c39f69cbf0dd15597da3b6803d0aa11b5b59513.
Ten report tests and 56 live ego-browser assertions pass at desktop, low-height
and both mobile widths. Its actual embedded-data hash matches the local payload;
the local full HTML hash was checked before and after QA. File-scheme source
fetch is unsupported, so a browser-computed full HTML hash is not claimed.

All sections/documents, expanded current progress, visible partial/not-run gates,
painted charts, filtering/reset, pagination, actual keyboard navigation and
print-media visibility pass. Fresh screenshot review, physical click placement,
physical printing/PDF pagination and future populated experiment QA remain
unverified. The exact-file sidecar and raw observations retain those limits and
the corrected stdout/stderr log capture. Space 29 was closed with done:true.
No source acceptance, model invocation, real accounting action or goal completion
is inferred from this report checkpoint.


### Human evidence projection integrity allocation — September 11

The preceding goal turn made verified committed progress at f976fc3. The current
worktree and draft PR match that checkpoint. Before edits, Codex claims go.mod,
go.sum, internal/app/{evidence_table.go,evidence_table_test.go,driver_checks.go,
driver_checks_test.go,driver_evidence.go,driver_evidence_test.go,driver_impl.go,
evidence_verify.go,evidence_verify_test.go}, the runtime guide and
implementation-notes/codex-1-table-integrity-dispositions-20260911.md.

Render the managed human table deterministically from actual typed records, bind
its exact bytes under a separate ExtraDigests key, and independently reconcile
that projection before helper execution/publication and final closure. Use a
CommonMark parser for actual top-level sections so examples, comments, nested
blocks and alternative heading syntax cannot silently select the wrong span.
Refuse ambiguous managed sections; keep all other implementation bytes bound,
including the authorized exact status transition. Retain failed attempts and
require fresh checks for older reports without a table binding. Test direct
post-recording edits, changed records with recomputed table hashes, missing or
duplicate sections, misleading markup/diagnostics, CRLF and actual helper/close
interleavings. Preserve the original quorum and all remaining full-scope gates.

The same claimed checks path also pins the exact non-evidence remainder before
commands run, compares it after execution, and retains the original binding on
failure. A missing managed section is initialized before that snapshot. This
closes a source-observed gap where commands could alter excluded implementation
scope before its first digest was taken; it does not exclude more content.


### Table integrity delivery allocation — September 11

Before edits, Codex claims the evaluation delivery/2026-09-05/{PROGRESS.md,
priorities.json,report.html,report-manifest.json} refresh and
browser-table-integrity-20260911.json with browser-qa/table-integrity-20260911.
Record exact-source validation and all initial negative cases, then update the
existing draft PR. All five display priorities (six underlying audit areas)
remain partial/preparation until their full independent and live gates pass.


### Human table integrity source checkpoint — September 11

Typed original records now determine the human validation section. Exact section
hashing and deterministic reconstruction run through the independent helper and
final status-write gate. CommonMark boundaries preserve examples and non-evidence
scope; duplicate/missing/rehashed-inconsistent tables refuse. Pre-execution scope
binding, CRLF and JSON-compatible Unicode rendering are covered. The owned table
integrity disposition records behavior, negative cases and limitations.

All 344 Go/module files match the manifest. Full Go suite PASS
(127.709s); app/evidence/driver/runner race PASS (142.938s).
Scoped vet, module verification, Windows amd64 app cross-build and isolated
shared-volume projection/actual-helper fixtures PASS. Original table-forgery and
Unicode failures, diagnostic ordering and Markdown fixture corrections remain
retained. Independent acceptance and all other full-scope obligations remain
pending; no real model or operator accounting operation occurred.


### Human table report verification checkpoint — September 11

Source 60640709ebb8b963311495900a981b763329c314 and its 344-file Go/module manifest
remain unchanged. The current 624432-byte offline report passes ten report tests
and 56 ego-browser assertions across desktop, low-height and two mobile widths.
Current source/projection behavior, retained negative cases, partial/not-run states,
all sections/documents, painted charts, search/reset/pagination, actual keyboard
navigation and print-media visibility passed. The first harness read the wrong
DOM dataset field and retained 21 assertion failures; correcting the harness
produced 56 passes without modifying or rebuilding the report.

Full report SHA256: 7e7c295fa72fe859151a509b4c97998fda3f17e9e2878655734c164aa1f39906.
Browser embedded-data SHA256: 955cbab005c2c3704c1a4fcf816c79e5f2d4f982897ceaa45638f6523c0e8f90.
Local full-file hashes bracket QA. New task space 30 was required because the old
space was absent; its dedicated cleanup returned done:true. Exact logs, retained
first attempt and observations are under delivery/2026-09-05/browser-qa/
table-integrity-20260911, with browser-table-integrity-20260911.json beside them.
No fresh screenshot, physical click placement or print/PDF pagination is certified.
This is current partial-report verification, not independent source acceptance or
future populated experiment verification. All remaining FINAL gates are unchanged.


### Resource guard continuity allocation — September 11

Before edits, Codex claims internal/budget/{lock.go,resource_guard.go,
resource_guard_test.go}, docs/agent-runtime-configuration.md and the owned
implementation-notes/codex-1-guard-continuity-dispositions-20260911.md.

The synchronization-only resource guard has no ledger; losing its origin and
local inode can look like first use even while an old descriptor remains held.
Reproduce that overlap with actual child processes before fixing it. Establish
a durable guard witness before returning permission, refuse missing or changed
origin continuity without recreation, and recheck exact pinned origin after
kernel acquisition/probing so a waiting caller cannot use obsolete origin data.
Retain failures, test write interruption and lock cleanup, and preserve current
budget/evidence publication behavior. This is a required foundation for safe
recovery, not an operator migration command or a claim that lock-origin recovery,
changed-inactive-import recovery or any other full-goal gate is complete.


### Resource guard continuity source checkpoint — September 11

A real second process reproduced split guard ownership after losing origin and
local lock pathname while the first holder remained alive. Six acquisition/probe
origin-mutation fixtures also granted work before correction. Guard witnesses and
read-only pinned-origin rechecks now reject these cases. Witness publication,
retry, conflicting origins, cancellation and released kernel ownership are tested.
The first negative process fixture's cleanup ordering was separately corrected;
the actual overlapping acquisition preceded its cancelled-child diagnostic.

All 345 Go/module files match the source manifest. Full suite PASS (132.785s),
budget/evidence/app/driver/runner race PASS (152.416s), scoped vet PASS (1.251s),
Windows amd64 app cross-build PASS and isolated shared-volume guard/lock/concurrent
reservation fixtures PASS (2.063s). The owned guard disposition records exact
behavior, retained counterexamples and limits. Safe migration/reconciliation,
semantic action replay, patch-regression trajectory and every independent/live
FINAL gate remain open. No model or operator accounting action ran.

### Guard continuity delivery allocation — September 11

Before edits, Codex claims evaluation delivery/2026-09-05/{PROGRESS.md,
priorities.json,report.html,report-manifest.json}, browser-guard-continuity-20260911.json
and browser-qa/guard-continuity-20260911. Preserve the exact previously verified
table report before refreshing the current partial report. Bind report QA to the
new source checkpoint, without implying independent implementation acceptance.


### Resource guard report verification checkpoint — September 11

The current 627683-byte report for source 1b891cf572deb3ee3975b221c051af745a8120ff
passes ten report tests and 56 live ego-browser assertions at 1440×900, 1280×540,
390×844 and 320×720. Current guard behavior and retained negatives, all sections
and documents, partial/not-run gates, painted charts, search/reset/pagination,
actual ArrowRight navigation and print-media visibility passed. The existing
README-documented Python environment built the report without installation.

Report SHA256: 72ff61c3f62df25b94a5304097124621c2944a0bdc058824922c993609e19d71.
Embedded payload SHA256: 5d1f20782bc19ac683312e6435d22c2c227706bcad34749e5aece54e166f757e.
Local full-file hashes bracket browser QA. Dedicated cleanup of task space 31
returned done:true. Exact logs, checksums and limitations are under evaluation
delivery/2026-09-05/browser-qa/guard-continuity-20260911, with
browser-guard-continuity-20260911.json beside them. The prior table report is
archived under browser-qa/table-integrity-20260911/report.html at its original
hash. No fresh screenshot, physical click-placement or print/PDF pagination is
certified. This remains current partial-report QA; independent acceptance and
all live experiments/follow-ups remain open. Source manifest is unchanged.


### Changed inactive import recovery allocation — September 11

Before edits, Codex claims internal/budget/{migration_recovery.go,
migration_recovery_test.go,launch_migration.go,protocol_migration.go},
internal/app/{budget.go,budget_migrate.go,budget_migration_recovery.go,
budget_migration_recovery_test.go}, docs/agent-runtime-configuration.md and
implementation-notes/codex-1-import-recovery-dispositions-20260911.md.

Add read-only preview and attended exact-decision recovery for changed inactive
launch, step, fixup and cross-review imports. Preserve original migration.json,
original policy ceilings and all original charges; append a bounded hash-linked
recovery journal with actual current history and explicit accounting assertions.
Merge newly observed launch attempts without recharging existing identities;
retain unknown costs and refuse conflicting prior observations. Protocol totals
and explicit pre-telemetry launch counts cannot shrink. An accounting epoch may
move earlier to include newly found history, never later. Freeze both current
import-state and current history hashes before publication. Exact replay handles
interrupted persistence; changed history needs a new operator decision. A versioned
activation witness binds the journal and makes old readers refuse recovered state.

Exercise every accounting kind, interrupted write boundaries, changed-again
history, concurrent exact/conflicting controls, active replay with later charges,
no policy extension, unknown costs and malformed/stale/unattended decisions. This
is implementation and fixture verification only; no actual operator accounting
operation or model call is authorized by these test decisions. Lock-origin
migration, semantic action replay, trajectory and all independent/live FINAL
requirements remain intact.


### Changed inactive import source checkpoint — September 11

Read-only preview and attended exact-decision recovery now cover launch, step,
fixup and cross-review imports whose history changed before activation. Original
imports/ceilings and all charges remain; a bounded hash-linked journal plus a new
activation marker binds each reconciled checkpoint. Counts cannot shrink, epochs
cannot move later, conflicting terminal observations refuse, and exact replay
preserves later charges and separate extensions. The owned recovery disposition
records behavior, limits, all publication-failure fixtures and the initial compile
correction. No independent acceptance is claimed.

All 349 Go/module files match the source manifest. Full suite PASS (143.697s),
budget/evidence/app/driver/runner race PASS (159.952s), scoped vet PASS (0.763s),
Windows amd64 app cross-build PASS and isolated shared-volume recovery fixtures
PASS (17.996s). Separate-process exact/conflicting controls and nine actual compiled
CLI fixture probes passed. No actual model or real project accounting action ran.
All other independent/live FINAL gates and original quorum remain unchanged.

### Import recovery delivery allocation — September 11

Before edits, Codex claims evaluation delivery/2026-09-05/{PROGRESS.md,
priorities.json,report.html,report-manifest.json}, browser-import-recovery-20260911.json
and browser-qa/import-recovery-20260911. Preserve the exact prior guard report before
refreshing the partial report. Bind new report verification to this source;
keep independent acceptance and every outstanding experiment/follow-up explicit.


### Changed inactive import report verification checkpoint — September 11

The current 629641-byte report for source bd0018865946764cd65cd5cde7860bb5bfd480c9
passes ten report tests and 56 ego-browser assertions at 1440×900, 1280×540,
390×844 and 320×720. Current recovery behavior, retained negative history and
initial compile correction, all sections/documents, partial/not-run states,
painted charts, search/reset/pagination over 673 ideas, actual ArrowRight
navigation and print-media visibility passed. Normal media and desktop metrics
were restored. Dedicated cleanup of task space 33 returned done:true; unrelated
spaces 1, 2 and 32 were not used or closed.

Report SHA256: 06ed736c5ef7379adf1ce3c3d1adc51f5afbb34ef6869b70798f8ad28b87a0be.
Embedded payload SHA256: 2c01b0d2b78666bf30c6549848922688133e3416c0dac13cb320c4a774e5e685.
Local full-file hashes bracket browser QA, and the exact report was not rebuilt.
All checks, observations, confirmed cleanup and file hashes are retained under
evaluation delivery/2026-09-05/browser-qa/import-recovery-20260911, with
browser-import-recovery-20260911.json beside them. The previous guard report is
archived at its original hash. No fresh screenshot, physical click placement or
physical print/PDF pagination is certified. This verifies the current partial
report, not independent source acceptance or future populated experiment results.
The exact 349-file source set and all validation evidence hashes remain unchanged.
All outstanding FINAL gates, historical quorum and recovery decisions remain open.

### Exact active reservation continuity allocation — September 11

Before edits, Codex claims internal/budget/{reservation_receipt.go,
reservation_receipt_test.go,step_session.go,cycle_session.go,cycle_binding.go}, the owned
implementation-notes/codex-1-reservation-continuity-dispositions-20260911.md and
docs/agent-runtime-configuration.md. Existing session tests show nesting and total
charges, but retained sessions do not bind revalidation to their exact reservation.
Reproduce disappearance/substitution of the session's own entry while retaining
the same aggregate count; fail before further nested execution on any changed
original charge. Retain the original reservation identity, kind, timestamp and
reserved exposure, preserve legitimate later charges/extensions/settlements, and
keep publication errors sticky. Use the existing ledger and locking mechanisms.
This is required active-action continuity, not a claim that durable cross-process
semantic replay, origin migration, trajectory or independent/live FINAL gates
are complete. Record the reproduced negative and current-source validation.

### Exact active reservation continuity source checkpoint — September 11

Fifteen real ledger/session fixtures reproduced further nested admission after
same-count substitution or a changed original timestamp/reserved exposure/epoch.
Step, fixup and cross-review sessions now pin their exact original reservation
from its publication result and recheck it before further nested work. Refusals
remain sticky and never rewrite accounting. Later charges, valid extensions,
settlement/reconciliation and the last allowed live group remain supported.
The owned reservation-continuity disposition records the exact scope and limits.

All 351 Go/module files match manifest
c3e2507ec3c76781e1606e97a466ae34b198d38b68b5fc9f615b1a52436952c8.
Full suite PASS (146.656s), budget/evidence/app/driver/runner race PASS (165.232s),
scoped vet PASS (1.256s), Windows amd64 app cross-build PASS and the compiled
budget fixtures on an isolated shared volume PASS (5.721s). No model or actual
operator accounting operation ran. This closes the identified live-session gap;
durable cross-process semantic replay and all other unfinished FINAL obligations
remain open. The latest browser-verified report remains explicitly bound to the
prior bd00188 source checkpoint; refreshing presentation is still pending.

### Patch-regression trajectory execution allocation — September 11

Before edits, Codex claims the new internal/trajectory/{trajectory.go,
verify.go,trajectory_test.go,verify_test.go} files, the owned
implementation-notes/codex-1-trajectory-execution-dispositions-20260911.md and
docs/agent-runtime-configuration.md. Build the execution and decision core for
FINAL D6/AC-B2 using the existing evidence.RunCriterion machinery, without
changing the ordinary completion attestation's pass-only contract. Freeze exact
before/after commit and tree identities, patch bytes, criterion command hashes,
independent runtime verifier and ordered patch-request ancestry before checks.
Run actual paired baseline/patched checks in AB/BA order and require stable trees
and consistent structured observations; parser/process failures, missing coverage,
ambiguous existing failures and self-attribution are not confirmed regressions.
Derive outcomes from retained executions rather than a supplied verdict. Detect
the first two consecutive confirmed material criterion regressions, preserve that
review requirement, and reject replayed patches or omitted sequence entries.

Use executable Git fixtures for new regression, unchanged prior failure, clean
intervening patch, instability, missing/zero-case/opaque evidence and source drift.
This establishes the paired-execution core only. Persisted opt-in policy and
complete driver/manual/verifier invocation wiring remain required before AC-B2
can be accepted; no prototype API or fixture will be presented as that gate's
completion. No real model invocation is part of this implementation step.

The paired failure classifier also needs typed execution completeness: a failing
test seen in a truncated stream must not establish a confirmed regression.
Before edits, Codex additionally claims internal/evidence/execute.go and the new
internal/evidence/execution_observation_test.go for an additive detailed-execution
API. Preserve existing RunCriterion and persisted report shapes; expose process,
capture and parser completeness as runtime observation rather than parsing the
human diagnostic string. Existing evidence owner artifacts remain untouched.

### Patch-regression execution core checkpoint — September 11

The new trajectory core freezes exact clean before/after Git snapshots and
material command hashes, performs actual paired AB/BA criterion executions, and
derives consecutive regression decisions from the retained observations. Two
confirmed material regressions trigger a retained review requirement; a clean
intervening patch breaks the sequence. Prior failures, unstable/incomplete output,
reduced passing coverage and missing/replayed/changed-scope records cannot supply
false confirmations. The additive detailed criterion API supplies typed capture,
parser and process completeness while preserving ordinary evidence report shapes.

All 356 Go/module files match manifest
6f60552b591b1c23ff448e3c6b94b5672ca5e6164701f32028ba82a643077e84.
Full suite PASS (124.571s), trajectory/evidence/budget/app/driver/runner race PASS
(146.321s), scoped vet PASS (1.237s), Windows amd64 trajectory/app cross-builds
PASS and the compiled trajectory fixtures on an isolated shared volume PASS
(60.063s). Initial SafeLabel pointer-type build failures are retained; all
subsequent runtime validation passed. No real model or operator action ran.

The owned trajectory execution disposition records current behavior, trust
boundaries and concrete remaining integration. This is an execution/decision
core, not a persisted opt-in policy or enforcement across driver/manual paths.
AC-B2 remains incomplete until complete real patch coverage, source snapshots,
independent verifier invocation, durable observations and all enforcement/recovery
paths are wired and independently accepted. Other FINAL gates remain open.
The last browser-verified report remains bound to bd00188, pending refresh.

### Persistent trajectory accounting allocation — September 11

Before edits, Codex claims internal/budget/{cycle_binding.go,cycle_extension.go,
cycle_session.go,cycle_observer.go,cycle_observer_test.go}, new
internal/trajectory/{state.go,state_test.go}, internal/runner/{cycle_budget.go,
telemetry.go,launch_budget.go,trajectory.go,trajectory_test.go},
internal/driver/{cycle_budget.go,impl.go}, internal/app/{app.go,trajectory.go,
trajectory_test.go,driver_impl.go}, docs/agent-runtime-configuration.md and the owned
implementation-notes/codex-1-trajectory-accounting-dispositions-20260911.md.

Persist explicit opt-in scope before the first fixup charge, binding original
material criteria and implementer to the shared cycle authority. Use existing
resource guards and synced publication, and pin required observation in the
cycle policy so missing state or an older writer cannot silently bypass it.
At the actual reservation boundary retain every exact charged attempt, including
failed or interrupted work. Record actual pre/post Git and source observations;
dirty results remain dirty and unresolved, without committing or discarding them.
Bind each launch to its actual invocation and reservation. Refuse new patches and
completion while coverage is missing or a prior attempt awaits independent
verification. Exercise manual/driver accounting, restart, exact ledger identity,
state deletion, publication failure and dirty/nonzero execution fixtures.

The previously validated execution core is committed at 26a6b434dc9433597e17853a86c05d2e3a2caaaa.
This allocation continues AC-B2 integration, not a reduced acceptance target.
Independent model invocation, verified helper receipt publication, restoration
of interrupted/dirty snapshots and the full trajectory decision/explicit recovery
path remain required before declaring AC-B2 or the overall task complete.

Before adding driver integration fixtures, Codex also claims
internal/driver/trajectory_test.go. Verify the actual driver reservation boundary
and fresh-process refusal against its shared state, including a finite budget
extension that must not resolve missing independent trajectory evidence.

### Persistent trajectory accounting source checkpoint — September 11

Explicit opt-in capture now binds the original material criteria and implementer
to the shared fixup policy. Driver/manual precharge records every exact attempted
reservation before execution permission. The instrumented runner binds its real
invocation and retains failed/interrupted terminal and dirty/unavailable source
observations. Missing entries, changed original charges, repeated launches and
pending observations refuse another fixup and both completion paths. A linked
worktree and a separate process retain that refusal after a finite extension.
The owned trajectory-accounting disposition documents the trust boundary and
remaining independent execution/snapshot/disposition integration; AC-B2 is open.

Two live-handle negative cases reproduced incorrect success after the required
accounting directory disappeared before launch/terminal capture. Both now refuse;
an actual child that displaces accounting and exits zero yields trajectory_failure.
All 364 Go/module files match manifest
5bc5352d75ff5f5bc5e8c1da088510199ae97595c7a8ac55039fa4103020e8e2.
Current full JSON suite PASS (152.744s, all 32 package terminal events), scoped
race PASS (160.730s, all six packages), vet PASS (1.100s), Windows amd64 trajectory
and app cross-builds PASS. Compiled shared-volume trajectory/runner/driver/app
fixtures PASS (24.074/10.900/1.852/6.234s). Native/shared JSON logs match exactly;
final-verification.json checks every new top-level test passed, complete package
coverage, unchanged source, and artifact checksums. Windows runtime is unverified.

A preceding broad run exited 1 (156.489s) with an incomplete sparse log containing
NUL bytes; its cause remains unclassified. The preceding race returned 0 but its
log was also sparse. Those attempts remain retained as incomplete evidence;
replacement runs captured complete structured output locally before shared
publication. Later success does not explain away the earlier failure.
No model invocation or actual operator approval/activation ran. Current capture
retains identities/digests, not yet reconstructible archives of dirty contents.
Every captured attempt stays pending until actual independent helper execution,
durable receipt/decision publication and explicit recovery are integrated.
All other FINAL obligations, pending recovery decisions and the older bd00188
HTML report state remain unchanged. No final merge/release/deployment is claimed.

### Reconstructible trajectory source allocation — September 11

Before edits, Codex claims new internal/trajectory/{snapshot.go,snapshot_test.go},
internal/trajectory/{state.go,state_test.go}, internal/runner/trajectory_test.go,
internal/app/trajectory_test.go, docs/agent-runtime-configuration.md and the owned
implementation-notes/codex-1-trajectory-snapshots-dispositions-20260911.md.

Retain bounded private archives of the actual source at activation/precharge and
terminal capture, including dirty additions, deletions, modes and supported
relative symlinks. Use archive/tar, os.Root containment and existing synced
publication. Bind archive bytes to source observations and recheck source drift.
Restore only into a newly allocated empty private directory, validate paths,
member types, counts/sizes, links and complete source digests, and never modify
the live Git index, HEAD, branches or worktree to manufacture a clean snapshot.
Keep archive bodies private and outside public telemetry/protocol artifacts.

Version persisted trajectory state so old digest-only attempts cannot silently
claim reconstructible contents. Archive loss/corruption/publication failure must
leave an explicit unresolved attempt; no invented post-state or acceptance.
Exercise actual process capture, restore-after-live-edit, removed tracked files,
untracked files, executable bits, supported links, malformed/truncated archives,
path traversal, symlink escape, source drift and unchanged user Git state.
This continues AC-B2 toward independent paired verification. Snapshot restoration
alone does not invoke a verifier, derive a final verdict or grant continuation.

Before changing Git observation, Codex additionally claims
internal/trajectory/verify.go. Pass Git's --no-optional-locks for read-only
snapshot commands so status refreshes cannot rewrite the user's index. The
ordinary completion evidence contract remains unchanged; trajectory observation
may explicitly exclude only paths observed absent from the actual working tree,
while bracketing that deletion with source/status checks.

### Reconstructible trajectory snapshot checkpoint — September 11

Opted-in policy/state v2 now binds private before/after source archives to the
original reservation and actual model invocation. Dirty additions, deletions,
modifications, binary files, permission bits and supported relative links remain
restorable after further live edits. Capture/restore preserve the user's Git
index, HEAD, refs and live files. Canonical archive bytes, metadata, inventory,
source digest and restored filesystem are checked; corrupt/missing archives,
traversal, parent conflicts, unsupported links/types, truncation, trailing data,
excessive size and source drift refuse. A failed archive publication retains its
actual terminal/source as archive-unavailable. Old v1 digest-only state is
explicitly refused with recovery guidance and is not rewritten from current files.

The initial archive serializer incorrectly added a regular-file newline; the
actual integration failure is retained and corrected. Private temp-directory
and PAX header construction errors were fixture corrections, also retained.
Nine new archive and three new state tests pass, plus the updated actual runner
manual/grouped nonzero/interrupted restoration and CLI archive-loss tests.
All 366 Go/module files match source manifest
7675f7dc655644e7e8d3f1c1e24da57a5c8f52caf2d943c26c6fad11e7eb3fd2.
Full JSON suite PASS (167.333s, all 32 package terminal events), six-package race
PASS (183.107s), vet PASS (0.865s), Windows amd64 trajectory/app cross-builds PASS,
and compiled shared-volume trajectory/runner/driver/app fixtures PASS
(63.183/13.339/2.332/4.923s). All twelve new and three amended top-level tests
have passed events in both full/race runs. Complete native/shared logs match,
contain no NUL bytes and no failed test events. final-verification.json and forty
checksummed evidence files live in the owned snapshot validation directory.
Windows runtime and current-source independent acceptance remain unverified.

The owned trajectory-snapshots disposition records implementation, trust boundary,
executed failures/corrections and next integration. Restoration creates source
files without asserting Git commit ancestry; the verifier still needs isolated
execution roots explicitly bound to original archived source. Actual independent
runner/helper invocation, durable receipts/observations, expected charged patch
inventory evaluation and review/disposition recovery remain open. Every attempt
stays pending; AC-B2 and the six-area goal are incomplete. No real model call,
operator activation, final merge/release/deployment, global install or core
publication occurred. The older bd00188 HTML, historical quorum/full-six pilot,
Claude rate-limit recovery and unanswered experiment decisions remain unchanged.

### Charge-bound captured-source verification allocation — September 11

Before edits, Codex claims new internal/trajectory/{captured.go,captured_test.go},
internal/trajectory/trajectory.go, docs/agent-runtime-configuration.md and the
owned implementation-notes/codex-1-captured-verification-dispositions-20260911.md.
Connect retained source/charge identity to actual paired criterion execution.
Restore both private source archives, attach independent local copies of the
original Git object history, detach at each actual recorded commit, and require
exact original Source observations before execution. Never manufacture original
clean commits or edit the live worktree/index/refs. If staged state or missing
Git objects cannot be reconstructed exactly, refuse explicitly instead of
substituting a different Git context. Retain original dirty status and commit
identity in the comparison request.

Freeze policy, complete original attempt/charge, invocation, archive references,
material criterion hashes and selected non-implementer verifier. Bind actual
AB/BA executions to these facts and recheck source/authority around execution.
Reuse existing format-aware criterion assessment; ordinary pass-only completion
attestation and clean-Git comparison semantics stay unchanged. Return partial
observations on interruption/drift. Add real archived dirty/deleted/untracked
source execution, charge/criterion/source substitution, missing history, Git-state
preservation and staged-state refusal fixtures.

This continues the required independent runner/helper path. The API itself does
not invoke a model or authenticate a participant. Durable helper requests,
invocation receipts, complete trajectory evaluation and review/recovery remain
required before accepting attempts or declaring AC-B2 complete.

### Charge-bound captured verification checkpoint — September 12

FreezeCaptured now pins the complete shared policy/attempt/charge, original
invocation, material criteria and both source archives. OpenCaptured restores
private before/after source and independently copies original Git objects without
hardlinks, checkout filters or template hooks, detached at the actual commits.
Both roots must reproduce the original Source observation exactly; dirty output
keeps its original HEAD/status, and no synthetic clean commit is introduced.
Staged-index differences, missing original commits, object alternates, excessive
or unsupported storage and nested Git parents refuse. Original source, index,
HEAD and refs are preserved. VerifyCaptured rechecks authority and both sources
around actual AB/BA executions and returns partial observations on drift or
interruption. Assessment is shared with the original clean-commit comparison;
ordinary pass-only completion semantics remain unchanged. A workspace handle is
single-use but does not yet provide cross-process durable helper replay control.

All ten new tests pass, including actual dirty failed child output after live
edits, original committed ancestry, exact execution order, clean/inconclusive/
opaque outcomes, request/authority/scope substitutions, copied object independence,
source mutation, signalled real-process interruption, staging refusal and unchanged
criticism. All 368 Go/module files match manifest
bf97070a3ed66c535b9b49a49e5a79f0a04627ac61274e6a992e862652895970.
Full JSON suite PASS (167.580s, all 32 package terminal events), six-package race
PASS (183.213s), vet PASS (1.493s), Windows amd64 trajectory/app cross-builds PASS
with PE amd64 headers verified. Compiled shared-volume trajectory/runner/driver/
app fixtures PASS (162.394/10.170/1.652/3.290s). Every new test passed in full,
race and the shared-volume run; logs match native copies, contain no NUL bytes
and no failed-test events. final-verification.json and thirty-nine checksummed
evidence files are retained under captured-verification-validation-20260912.

The owned captured-verification disposition records supported source reproduction,
trust/host limits, staged/index/object retention gaps and remaining integration.
This API does not launch a model, publish or authenticate an independent helper
receipt, resolve pending state or authorize continuation. Actual runner/helper
invocation, durable receipts/observations, full charged patch-list evaluation,
retained review trigger and explicit recovery remain required. AC-B2 and the six
areas remain incomplete. No new model call or actual operator activation occurred.
The older report, real experiment requirements and unanswered recovery/roster/
pilot/funding decisions are unchanged. No final merge/release/deployment/global
installation/immutable-core publication is claimed.

### Ownership claim — durable captured verification journal, September 12

Codex-1 owns `internal/trajectory/verification.go`, its new tests, the private
observation hook in `internal/trajectory/captured.go`, runtime documentation and
`implementation-notes/codex-1-durable-verification-journal-20260912.md`.

Persist the exact captured request and selected verifier/run before launch in
shared charge-scoped storage. Reserve one observed invocation and one helper
claim durably, refusing duplicate or interrupted replay without inventing a new
patch charge. Retain bounded ordered execution observations and a terminal
receipt, including preparation failures and partial executions. Re-read exact
request/launch authority around helper execution; accept only a complete matching
receipt under unchanged original authority. Use actual cross-process fixtures
for competing claims, interruption and replay. Record the same-UID trust limit.

This internal journal is a prerequisite for the instrumented model/CLI helper
integration. It does not itself launch/authenticate a model, resolve trajectory
attempts, authorize retries or continuation, change quorum, or satisfy AC-B2.
Existing source archives, pass-only completion and older report remain scoped
as previously recorded. No real model or operator action is part of these tests.

### Durable captured verification journal checkpoint — September 12

The shared journal now reserves one exact captured ticket per original fixup
charge before launch, one distinct invocation before spawn and one exclusive
helper claim before source preparation or criterion execution. Prepared private
source roots are retained before checks. Every actual AB/BA execution is written
as an immutable bounded step with its ordinal and prior digest. A terminal receipt
binds request, launch, claim, preparation and all retained steps. Changed authority,
missing/ambiguous/extra history and lost observation/receipt writes refuse. Killed
helpers leave completed steps and a pending claim; fresh processes cannot silently
retry or replace the attempt. Checked reads return partial observations with an
error for missing/failed terminals and never resolve the underlying charge.

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

Actual synthetic process fixtures cover competing helpers, exactly one execution
sequence, a killed helper after its first retained observation, fresh-process
replay refusal, source/scope/identity/history mutations and required-write failure.
The owned durable-journal note records POSIX, same-UID attribution, interrupted
workspace cleanup and remaining staged/object recovery limits. This is an internal
API; actual instrumented model/helper invocation, inherited-marker/terminal binding,
full charge-derived evaluation and durable review/recovery/continuation remain
required. Every trajectory attempt remains pending. No new model invocation or
actual operator activation occurred. AC-B2, the six areas, independent current-source
acceptance, live experiments and final report/follow-up delivery remain incomplete.
The previously unanswered provider recovery, historical roster/pilot and funding
choices remain unchanged; no merge/release/deployment/global publication is claimed.

### Ownership claim — independent trajectory runner/helper integration, September 12

Codex-1 owns the new `internal/runner/trajectory_verification.go` and its tests,
the shared `internal/runner/telemetry.go` launch boundary, new
`internal/app/trajectory_verify.go` and its tests, the trajectory CLI dispatch,
minimal journal accessors if required, runtime documentation and
`implementation-notes/codex-1-trajectory-helper-integration-20260912.md`.

Bind a prepared captured-verification ticket to the instrumented runner before
spawn, refusing changed origin/run/idea/agent/phase and duplicate launches. Expose
an explicit selected-verifier command and runtime helper. Freeze the parent
request before launch, require inherited process attribution, exact original
criteria and the shared ticket, and invoke the existing durable helper journal.
The parent must compare observed terminal records and complete retained helper
receipts, preserving failed/partial attempts and refusing text-only PASS,
self-verification, replay or changed authority. Tests use actual executable CLI
fixtures, not provider calls or real operator activation.

This integrates the independent execution path; complete charge-derived trajectory
evaluation and explicit review/recovery/continuation remain required. Do not reset
charges or mark AC-B2/overall delivery complete. Historical quorum and unanswered
recovery/pilot/funding decisions remain intact.

### Ownership extension — nested verifier cancellation, September 12

The initial full suite passed, but the exact-source race run reproduced a
TERM-resistant criterion surviving verifier cancellation. The outer runner and
nested criterion used competing termination grace periods; the helper could be
killed before finishing its own child cleanup. Retain this failed evidence.

Codex-1 additionally owns the criterion-start control hook in
`internal/evidence/execute.go`, process-control records and stop handling in the
captured verification journal, and runner cancellation wiring. Register the actual
criterion under the same shared guard used to request stop, prevent starts after
that stop and terminate the owned active criterion before stopping its enclosing
verifier group. Preserve explicit attribution refusal for unknown processes and
record no acceptance on cancellation. Test the reproduced race and normal execution;
do not merely rerun the failing checkpoint or weaken the fixture.


### Process-control fixture ownership and validation boundary — September 12

The nested-cancellation ownership extension also includes
`internal/evidence/execute_control_test.go`,
`internal/trajectory/verification_control_test.go`, the existing journal
history-substitution tests and the existing app cancellation fixture. New
claims/steps/receipts use journal version 2; completed version-1 receipts remain
readable without inventing process-control evidence. The process-control records
contain stable supervisor identity, not raw material command text.

The previous exact-source race and compiled shared-volume app failures are
confirmed independently in their complete retained logs. Focused corrected tests
pass. The first corrected full/race/shared validation is in progress under a new
manifest; it is not final acceptance. Code review also identified the need to
preserve incomplete status for signal-like material exits through the stable
supervisor. Validate that boundary before committing or certifying this change.


### Supervisor signal counterexample — September 12

A separate Go overlay test, without changing the source under the in-flight
full/race/shared run, printed a structured failure and killed its material shell
with SIGKILL. The stable supervisor returned exit 137; the executor incorrectly
marked the observation complete. The counterexample failed as expected and is
retained with its exact overlay/test source under the first corrected validation
checkpoint. Reject signal-like controlled exits conservatively before accepting
this integration. A passing ordinary suite does not settle this new boundary.


### Independent helper integration and controlled cancellation checkpoint — September 12

The selected-verifier CLI now executes through the actual instrumented headless
runner and binds the original frozen ticket to its reserved invocation before
spawn. The helper requires the exact parent request, inherited markers, original
scope/quorum and durable shared reservation. The parent requires matching observed
and retained successful terminal telemetry, complete AB/BA helper observations
and successful publication of the request/terminal/receipt-bound result. A written
PASS, self-verification, changed authority, failed start/exit, lost receipt or
lost parent publication cannot produce an accepted assessment.

The reproduced nested-cancellation race is corrected by guarded stop/start
coordination. A stable waiting supervisor is identified and persisted before
its material command is released under the shared cycle guard. A durable stop
prevents later claims/starts; the parent reaps strictly attributed registered
criterion groups before stopping the enclosing verifier. Required process-write
failure leaves the material command unstarted. The global procctl attribution
checks were not weakened. New journal version-2 steps/receipts bind process
records; version-1 completed journals remain readable without inventing new
control evidence. Stop-marked or mixed/substituted history refuses acceptance.
Signal-like controlled exits >=128 remain incomplete, including an explicit 137,
because the supervisor cannot distinguish those cases from signal termination.

Final validation covers all 376 Go/module files with manifest SHA256
b9cbd3b1372c7fdc7d560b54134915730d2fc28bc017909683cc6ea6026b2404.
The full JSON suite passed (237.313s), with all 32 package
terminal events: 31 tested packages passed and cmd/parley had no tests. The
six-package race run passed (258.405s); scoped vet passed
(1.344s). Compiled shared-volume trajectory/runner/driver/
evidence/app fixtures passed (264.528/
28.928/1.746/
0.595/80.785s).
All thirteen new top-level test entries passed in full, race and their compiled
shared-volume runs, including both previously failing boundaries. Windows amd64
trajectory/app cross-builds passed and their PE headers/hashes were checked;
Windows execution remains unverified. Complete native/shared logs match, contain
no NUL bytes and no failed-test events in the final runs. The verifier checked the
unchanged source manifest, test coverage, package terminals and older HTML hash.
`final-verification.json`, validation/verifier scripts, toolchain identity and
56 checksummed evidence files remain under
`.parley-runtime/trajectory-helper-final-validation-20260912/`. Prior cancellation
failures and the signal counterexample remain distinct retained evidence; none
is relabelled as a pass by this checkpoint.

These tests use synthetic executable participants, actual parent/helper CLIs and
real process groups. No model invocation or actual operator activation occurred;
participant-owned source acceptance remains pending. Same-UID fabrication,
arbitrary daemonized sessions, abrupt outer-runner death, changed cleanup
authority and missing/staged source history remain explicit trust/recovery limits.
All trajectory charges remain pending. Complete charge-derived evaluation,
resolution-time parent-result reconciliation, retained two-regression review and
explicit recovery/continuation remain required; AC-B2 and overall delivery of the
six areas are incomplete. No quorum, historical signature or global roster changed.

The original full-six pilot, exact packet experiment, fresh independent reviews,
final populated report and actual-delivery-based follow-ups remain binding.
Previously unanswered Claude recovery, Hermes/Zcode membership/pilot amendment
and resource-ceiling decisions remain unanswered. The real-model inventory and
older report are unchanged. No merge, release, deployment, global install or
immutable-core publication is claimed.


### Ownership claim — durable trajectory reconciliation and continuation, September 12

Codex-1 owns new trajectory reconciliation/state transition code and tests,
shared parent-request/result types used by the existing helper CLI, the existing
captured-request/journal authority helpers, trajectory state validation and
snapshot checks, CLI dispatch/control and runtime documentation.

Re-read the exact retained parent request/result, observed terminal binding and
complete helper journal at reconciliation and subsequent runtime boundaries.
Derive ordered assessment coverage from every original charged attempt; never
accept a caller-filtered patch list. Preserve attempt hashes and original dirty
after captures. Record separate clean source promotion when the actual committed
source has identical material bytes and valid ancestry. Keep the first confirmed
two-regression trigger durable, and require an attended exact decision to continue
past review or inconclusive evidence without resetting charges or asserting a
clean outcome. Manual/driver/resume reservations must use the same checked state.

Use real synthetic executable participants and fresh-process CLI tests. Preserve
failed/partial reservations and all existing source/index/history refusals. These
changes do not satisfy independent current-source model acceptance, authorize the
live experiments, change quorum, or permit final merge/release/global publication.


### Reconciliation test extension and frozen validation — September 12

The four new executable app tests passed in 39.723s after adding simultaneous
fresh-process reconciliation, exact continuation replay, current criterion/quorum
substitution and retained-state mutation counterexamples. Evidence remains in
`.parley-runtime/trajectory-reconcile-extended-20260912/`. No new model call or
actual operator activation occurred.

Codex-1 additionally owns
`implementation-notes/codex-1-trajectory-reconciliation-20260912.md` for this slice.
Full-suite, six-package race, scoped vet, Windows cross-build and compiled
shared-volume validation are running under the new 380-file source manifest
ff1228819ee604db1f1f9974ad5eec1a50af807ee50905dca66253f19affbbf0.
This is a validation-in-progress record, not final acceptance. The current
Go/module files remain frozen until these results have been inspected.


### Reconciliation validation correction — September 12

The first frozen 380-file run passed the full suite (273.402s), six-package race
(298.428s), Windows cross-builds and shared trajectory/runner/driver/evidence/app
fixtures (256.787/30.621/1.665/0.596/172.201s). Vet failed because the moved
helper-request alias was initialized without field names. Its complete failed
log and the passing behavioral logs remain in the original validation directory.
A separate overlay with named fields passed vet without changing the source
under those running tests.

After all first-checkpoint processes terminated, the named-field correction was
applied. CLI usage now lists reconciliation/history/continuation and optional
acknowledgment flags. Only those two Go files differ from the preceding manifest.
A separate corrected source manifest and validation directory preserve the first
checkpoint as a vet failure rather than relabeling it as a full pass.


### Historical continuation replay correction — September 12

Before starting the corrected full validation, source review found that the CLI
computed the latest continuation preview before applying an older exact decision
replay. This could refuse an otherwise valid replay on a later dirty source, or
return the later preview alongside an earlier successfully replayed decision.
The existing three-patch executable test now exercises an exact first-decision
replay after the third patch and requires the original preview and unchanged
state. Retain its counterexample before correcting the return/control path.
Codex-1 owns that narrow CLI/API correction within the existing claimed files.
The keyed-only corrected validation directory has not executed any test; a new
final manifest will cover this additional behavioral correction.

The historical replay counterexample failed in 22.397s with the expected later-
source refusal. Returning the original preview from guarded apply/replay corrected
that behavior; all four focused app tests passed in 39.472s. Final validation now
uses the new 380-file source manifest 646ce1882aa32959f8c45344e265c4fe55d4e9ced103fef7d7e8271a7c184d10
under `.parley-runtime/trajectory-reconcile-replay-final-validation-20260912/`.
The earlier keyed-only directory was superseded before any validation ran.


### Existing launch-timeout fixture correction — September 12

At manifest 646ce1882aa32959f8c45344e265c4fe55d4e9ced103fef7d7e8271a7c184d10,
the full suite failed TestTrackedCommandForTimeoutKillsChildGroup at command
preparation: its original 150ms context expired before trackedCommandFor returned.
The failure did not observe a started process or demonstrate surviving descendants.
The six-package race run passed, as did the reconciliation app tests in full/race.
Preserve this full-suite failure separately; a passing race run does not erase it.

Codex-1 owns the correction in existing `internal/runner/launch_test.go`. Arm the
test-owned deadline after the descendant confirms readiness, release a survivor
canary only after the timed-out parent has been waited, and keep the assertion
that no descendant writes it. Preserve actual process start and timeout telemetry
checks. No production supervision, timeout or attribution rule changes here.

The corrected startup-gated fixture passed as an overlay in 2.333s. A separate
parent-only-kill control failed as expected in 4.158s with a surviving descendant.
Both complete logs and overlays are retained. Only `internal/runner/launch_test.go`
changed after the preceding checkpoint; production source is unchanged. Final
validation now includes this corrected fixture in its shared runner selection and
explicit test inventory, under manifest a10d0a02e8e68b0498b518778ea29de83b356a9c92080dfbf5678e3f31836504
and `.parley-runtime/trajectory-reconciliation-validation-20260912/`.


### Additional retained validation failures — September 12

The a10d0a02e8e68b0498b518778ea29de83b356a9c92080dfbf5678e3f31836504 full suite
failed TestLaunchBudgetRetainsFailedAndUnreportedSpend: terminal-write observed
no cost, and unknown-cost expired during command preparation. That fixture used
300ms for every scenario, including cases intended to isolate accounting writes
rather than timeout. The corrected child-group timeout fixture passed. The
six-package race run passed (297.308s). Shared trajectory, driver, evidence and
app fixtures passed; shared runner failed while constructing the source fixture
for TestCapturedVerifierLaunchRejectsChangedBinding/root with the existing
"source file changed during capture" refusal. This is not a completed validation.

Codex-1 owns the narrow accounting fixture correction in
`internal/runner/launch_budget_test.go`: allow 30s for non-timeout accounting cases;
use a deliberate 5s timeout against a 30s sleeper only for timeout-known-cost.
Preserve cost-settlement and reservation-exposure assertions. An unchanged-source
overlay passed all five scenarios in 6.352s before applying this fixture change.
No production accounting or deadline policy changes.

A separate diagnostic-only snapshot overlay retained full stat details if the
shared capture refusal recurred. Twenty fixed repetitions of the affected root
fixture passed (41.201s) without reproducing it. These later passes do not explain
the earlier failure. No production snapshot guard was weakened or changed; the
original shared-volume refusal remains unclassified and needs further diagnosis.
All original full/race/shared logs, source manifests and diagnostic overlays remain
retained. Do not call the validation final, commit/push as verified, or close AC-B2.


### Current uncommitted checkpoint — September 12

All test processes for this checkpoint are terminal. The latest targeted runner
race run passed both corrected fixture families (9.571s), with source manifest
794900ef5d33ecf8e4568190a9ae84bbba941f185c8b87b7b1d5d42cfb270698
and evidence under `.parley-runtime/reconciliation-latest-focused-20260912/`.
The preceding full run failed its accounting fixture, and its shared runner run
refused source capture in one fixture. All four new reconciliation tests passed
in that checkpoint's full/race/shared app runs. Complete native/shared logs match
and have no NUL bytes. The disposition and checksums preserve the failures under
`.parley-runtime/trajectory-reconciliation-validation-20260912/`; its validation
is not accepted. The snapshot refusal remains unclassified; no production snapshot
check was changed and the 20 later diagnostic passes are not an explanation.

No current reconciliation change was committed or pushed. PR #73 and the remote
branch remain at 7d54de924910a22aa228919dc42d7ce3b31c2ae2. Final validation, fresh
independent source acceptance and all live/recovery/delivery gates remain open.
No model invocation, actual operator activation, quorum change or deployment
occurred. `.parley-runtime/next-guard-origin-recovery-20260912.md` records a
read-only map for the following migration slice; it is not an implemented or
approved migration. OpenViking remains unavailable; no shared persistence claimed.


### Ownership and evidence — coherent snapshot reads on AppleVirtIOFS, September 12

The mounted shared volume reports AppleVirtIOFS. A controlled standalone probe
found mismatches between successful writes and subsequent rooted reads (764/2000
shared trials; 0/2000 native trials). A second probe compared direct-path and
rooted reads of the same file/inode. An 8000-trial held-root run observed 5957
stale rooted reads and one before/after metadata mismatch satisfying the exact
capture-refusal condition; direct-path reads matched all writes. Fresh roots
opened from the original absolute path had 0 stale reads and 0 capture-guard
refusals in 8000 trials. Reopening "." through the retained root did not solve it
(1215 stale reads and 3 capture-guard refusals). Full sources/logs and native
controls are retained in `.parley-runtime/snapshot-stat-probe-20260912/`.

This reproduces a concrete failure mechanism at the same guard; the original
single failed fixture retained no per-file stat details, so its exact event
is not retrospectively proven to have that cause. The older sparse/NUL anomaly
is separate and remains unclassified.

Codex-1 owns the narrow extraction in `internal/trajectory/snapshot.go`, new
`internal/trajectory/snapshot_read_test.go`, runtime documentation and this owned
disposition. For each material entry, reopen a root from its original path and
require the same pinned directory identity before reading. Keep rooted containment
and file identity checks. Use the opened file descriptor's size/mode/mtime for
the copy and post-read stability comparison; the preceding Lstat selects type
and pins identity, not potentially stale content metadata. Retain descriptor
bounds, required-close errors, complete archive/tree digest recomputation and
the final unchanged Source observation. No retries, source rewrites, budget reset
or relaxed expected source digest are introduced. Verify stale metadata, changed
inode/root, oversize-after-Lstat and actual source mutation boundaries.

### Snapshot read tests and remaining timestamp instability — September 12

The extraction is formatted and compiled. Six new tests cover real same-inode
stale metadata, replaced roots/files, a descriptor exceeding the file bound
after an earlier small stat, size/mtime mutation after descriptor stat, and 128
exact member reads while retaining an older directory handle. Four independent
negative overlays fail at the intended assertions when stale copy metadata is
restored or root/file identity and mtime protections are removed. The production
source is unchanged by these overlays.

Checkpoint `.parley-runtime/snapshot-read-validation-20260912/` freezes 381 Go/module
files at manifest 924031c2c7155b71114b8f15ffe854cbb5956c716fac3b9cd961d8a7c0a9d9ac.
Focused tests passed in 41.113s; shared snapshot tests passed in 8.302s. The full
suite passed in 263.110s (31 test packages, one command package without tests),
six-package race tests passed in 287.028s, vet and Windows trajectory/app builds
passed. Shared runner/driver/evidence/app tests passed, including both corrected
runner fixture families and all four reconciliation app tests.

The broader shared trajectory suite FAILED at
`TestCapturedVerificationRechecksDurableAuthorityAndScope/changed-scope` while
constructing its captured source: `source file changed during capture`.
All six new read tests passed within that same failed suite. Thirty diagnostic
repetitions of the affected fixture passed (87.121s) without a captured stat
failure; they do not explain or erase the earlier event. This checkpoint is not
accepted final validation and is not a reason to commit/push as verified.

A separate controlled descriptor probe now reproduces another sufficient
failure mechanism: write once, open/stat, wait 1.2 seconds without another write,
read/stat the same descriptor. Both fresh-root and absolute-path shared reads
changed mtime in 6/6 trials, with identical bytes, size and inode. Native-root
controls changed mtime in 0/6 trials. Repeating with descriptor fsync before the
first stat gave the same counts. These observations show why fresh roots alone
cannot guarantee stable timestamps on this mount. They do not prove which stats
changed in the original failed test, whose refusal did not retain those details.
The older sparse/NUL anomaly remains separate.

Complete native/shared logs match; the checkpoint disposition retains expected
negative failures and the unexpected shared-suite failure separately. No mtime
check, source digest, identity guard or archive bound was relaxed. No blind retry
or source rewrite was added. All executions are terminal and synthetic; no real
model invocation or actual operator activation occurred. HEAD/remote remain
7d54de924910a22aa228919dc42d7ce3b31c2ae2. Final validation, independent current-source
acceptance, remaining recovery work and all live/delivery gates remain open.

### Ownership and plan — bounded material revalidation, September 12

Codex-1 continues its existing snapshot.go/snapshot_read_test.go allocation.
First retain an executable counterexample through captureSnapshotRegular with
a delayed archive write, on native and shared storage. For timestamp-only drift,
the proposed fix requires one fresh, contained verification read of the bytes
already copied: identical pinned root/file identity, regular type, size, mode,
content hash and stable descriptor metadata throughout that verification read.
Size/mode changes and all I/O errors refuse immediately. A changed second read
also refuses; there is no loop, model retry, source rewrite or extra budget grant.
The whole archived tree must still equal the original expected source, followed
by the final exact source observation. This is material revalidation after an
ambiguous timestamp transition, not an assertion that timestamps never matter.
Add counterexamples for changed bytes, root/file replacement, second-read
mutation, truncation, oversize and cancellation before accepting this boundary.

The first conditional-revalidation checkpoint passed native focused tests but
failed two shared fixtures. Retained diagnostic archive bytes show the first
descriptor still returned original content while the absolute path returned
the changed content; cached mtime did not trigger the conditional second read.
Codex-1 therefore strengthens the owned implementation to require the one fresh
verification read for every regular member. No metadata heuristic decides
whether content verification is needed. The expected-tree fixture will mutate
before opening its member and explicitly prove it archived different bytes,
so it tests the full-tree guard rather than assuming what a stale descriptor read.


### Stable material revalidation checkpoint — September 12

The delayed-sink production-helper test initially used snapshotWrite, whose
additional chmod did not reproduce the timestamp transition. That passing
fixture is retained separately. A plain WriteFile variant then failed on shared
storage and passed natively. Its diagnostic showed the same 39-byte inode, size
and mode with descriptor mtime moving from 882341277 to 882430195 nanoseconds
within the same second. No subsequent source write occurred. Evidence is under
`.parley-runtime/snapshot-timestamp-write-only-20260912/`; the earlier variant
remains under `.parley-runtime/snapshot-timestamp-counterexample-20260912/`.

The first conditional revalidation passed native focused tests but failed shared
fixtures: an open descriptor could retain original bytes and timestamps while
the absolute path returned replacement bytes. Its diagnostic tar contents and
frozen source are retained under
`.parley-runtime/snapshot-revalidation-validation-20260912/`. The source/test
versions reconstructed for that checkpoint match their original manifest hashes.
The expected-tree fixture originally failed to establish that it archived the
changed bytes; it now explicitly checks that precondition before testing the
whole-tree refusal. No failed checkpoint is relabeled pass.

The final implementation performs exactly one fresh verification read for every
regular member. It checks contained root and file identity before reading, exact
copied bytes/size/mode, stable size/mode/mtime during verification, and root/file
name identity afterwards. A changed verification read, I/O error, replacement,
truncation, oversize file or canceled context refuses. Agreement between reads
cannot override the original expected tree or final unchanged source observation.
There is no retry loop, timestamp tolerance, source rewrite or additional grant.
This adds one bounded source read per regular member; it is not a performance
or live-pilot claim.

Accepted automated checkpoint:
`.parley-runtime/snapshot-stable-read-validation-20260912/`, 381 Go/module files,
manifest 4c982f04de789ff07fb8ba1508e19402fbc3a0a3c0bf54b61487c1caf8943576.
Focused tests PASS 42.988s; shared snapshot selection PASS 9.629s; full suite PASS
262.619s; six-package race PASS 289.588s; vet and Windows trajectory/app builds
PASS. Compiled shared trajectory/runner/driver/evidence/app selections all PASS
(254.496s / 36.301s / 1.565s / 0.633s / 160.148s). Windows runtime is unverified.
The final verifier inventories 30 relevant top-level tests, including all four
reconciliation tests, eleven snapshot-read tests and both corrected runner test
families. It checks 32 package terminals (31 test packages plus one command
package without tests), matching native/shared logs, unchanged source hashes,
PE amd64 build outputs, six expected negative-overlay failures and the unchanged
historical HTML. All processes are terminal.

No real model invocation, actual operator activation, final acceptance, live
experiment, quorum change, merge, release, deployment, global installation or
core publication occurred. This progress checkpoint can be committed to the
existing draft PR; the full objective and every remaining gate stay open.
The user was asked asynchronously about the concrete September 10 historical
quorum/pilot proposal and experiment ceilings. No approval is inferred from the
continuation request or elapsed time; the proposal is not a ratified amendment.

## Guard-origin relocation validation — September 12

The allocated slice is implemented and ready for independent review. Current
source manifest `ac1589bbf46569489f94ae29cda45c7f7cff4f3fe3ee04c5e4641d97c7e1e002`
covers 385 Go/module files. Final full/race/vet/platform/shared checks passed;
the final verifier checks all 32 package terminals, all 12 new top-level test
functions (including the process harness), five intended negative failures,
native/shared log equality, unchanged source and historical HTML, actual legacy
compatibility and actual CLI attendance refusal. See the owned note and
`.parley-runtime/lock-origin-migration-final-validation-20260912/`.

The earlier full pass (214.932s) has source_matches=false because the capacity
fix landed during that run; it is retained as superseded, not final evidence.
The accepted final full run is 210.427s on the fixed manifest. No real model call
or actual operator migration occurred. The live inventory remains 34 terminal
attempts, 17 unknown costs and USD 46.1887585 known CLI estimates; total unknown.
No participant artifact/signature or historical evaluation was rewritten.

Next independent work is durable semantic action identity and explicit recovery
of interrupted/lost-result/helper-ticket/unchanged-source/orphan-reservation
states. The read-only map is `.parley-runtime/next-semantic-action-recovery-20260912.md`.
Same-host cache relocation does not solve missing original identities, changed
repository scopes, large/changed snapshot stores, erased journals or Git history.
Existing quorum, pilot and funding decisions remain pending; no amendment or
spending authority is inferred from this checkpoint.


## Durable action accounting validation — September 12

The owned slice is implemented and ready for independent review. Source manifest
`eadb9d9c91d4038e6732c91d9134e2a745c712da03efee12e31ac151e84b3f1e`
covers 392 Go/module files. Accepted evidence is retained under
`.parley-runtime/action-identity-final-validation-v2-20260912/`; native copies are
at `/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-action-final-v2-7gzu_ry7`.

Focused PASS 4.879s; full PASS 220.143s; six-package race PASS 261.187s;
vet PASS; Windows budget/app/runner/driver cross-builds PASS (runtime unverified).
Compiled shared budget/app/runner/driver selections PASS in
2.834s / 0.432s / 0.705s / 0.542s. The final verifier checks all 32 package
terminals, thirteen new top-level tests (twelve substantive and their process
harness), five intended negative-overlay failures, native/shared log hashes,
unchanged source, actual CLI receipt replay/refusals and the unchanged historical
HTML. These are automated implementation checks, not independent participant
acceptance, live model comparisons or proof of exactly-once workflow effects.

The prior targeted fixture failures and five-field replay counterexample remain
under `.parley-runtime/action-identity-validation-20260912/`. The first attempted
final-validation focused run exposed the missing fixup result invocation ID and
is retained in `.parley-runtime/action-identity-final-validation-20260912/` as a
failed checkpoint. The final source includes both corrections. No passing result
from an earlier source is reused to certify this one.

Next required recovery includes missing parent results, consumed/unfinished
helper tickets, orphan charges and unchanged-source observations. The next map
identifies the existing parent/request/journal reconstruction path; it adds no
recovery implementation or permission. New sessions still use distinct spent
attempt IDs. The original six audit areas, exact packet trial, full-six pilot,
independent current-source reviews/signatures, final HTML and delivery-based
follow-ups remain incomplete. No real model invocation or actual operator
migration occurred; historical quorum, pilot and funding decisions remain open.


## Parent-result recovery and concurrent bootstrap validation — September 12

The owned slice is implemented and ready for independent review. Source manifest
`93dbaa9a63ed1500b72a71ba6eb154e9859c2e7153b4f19508106dd29077f6cd`
covers 397 Go/module files. Evidence and the final source/log/test verifier are in
`.parley-runtime/parent-recovery-final-validation-v2-20260912/`; native copies:
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-parent-final-v2-vsch5s7l`.

Focused PASS 124.831s; full PASS 265.057s; six-package race PASS 309.809s;
vet PASS; five Windows package cross-builds PASS (runtime unverified). Compiled
shared trajectory/app/budget selections PASS 26.816s / 113.069s / 1.010s.
The verifier accounts for all 32 package terminals, eleven new top-level tests
(ten substantive scenarios plus one process harness), six intended negative
controls, native/shared log equality, stable source and unchanged historical HTML.
The production CLI actually previews and applies recovery in two concurrent
processes. These synthetic helper/criterion executions are not live-model trials.

Partial-parent contradiction and decoder-boundary failures remain in the initial
validation directory. The first full 396-file source run failed an existing
concurrent-reservation test (274.123s overall) and remains a failed run under
`.parley-runtime/parent-recovery-final-validation-20260912/`. A deterministic
concurrent boundary test reproduces its stale absence observation; the correction
refreshes only that absence and retains original location/identity/kernel checks.
Missing/conflicting origins or identities cannot be recreated. The fresh final
source is validated independently of earlier-source passes.

Recovery preserves the original parent publication state, derives original
independent executions and grants no execution, budget, continuation or acceptance.
It handles lost already-reconciled parents only when every pinned fact/digest
matches. New reconciliations also retain the exact recovery record. Other unfinished
helper-ticket/orphan/unchanged-source and workflow-effect recovery remains open.
The next inspected map covers existing stop/process/charge boundaries without
claiming that a missing terminal or dead leader proves all descendants inactive.

Current-source participant acceptance, full live launch/concurrency/closure
coverage, exact packet experiment, full-six pilot, owned signatures, final HTML
and delivery-based follow-ups remain incomplete. No participant-owned artifact
or signature changed. No real model invocation occurred; live inventory remains
34 terminal attempts, 17 unknown costs, USD 46.1887585 known CLI estimates and an
unknown total. Historical quorum/pilot/funding decisions remain pending.


## September 14 continuation and disjoint ownership

User direction: continue the authorized unfinished audit implementation (translated
from Slovak). Published checkpoint is 2ab9e82f0dd5a885717483a363de9482ae7c938a,
PR #73, with accepted unchanged-source validation retained locally. The reported
Claude reset has elapsed; current provider availability will be measured, not inferred.

| Boundary (file set) | Owner | Branch | Worktree | Status |
| --- | --- | --- | --- | --- |
| Only new implementation-notes/claude-1-reconciliation-review-20260914.md and its native atomic siblings under this idea | claude-1 | review/meta-protocol-change-evidence-first-efficiency/reconciliation-20260914 | ../worktrees/evidence-first-reconciliation-review | review |
| internal/budget/cycle_binding.go, cycle_observer.go; new internal/budget/cycle_intent.go, cycle_intent_test.go; internal/trajectory/state.go; new internal/trajectory/reservation_recovery.go, reservation_recovery_test.go; internal/app/trajectory.go; new internal/app/trajectory_reservation_recovery.go, trajectory_reservation_recovery_test.go; docs/agent-runtime-configuration.md; this IMPLEMENTATION.md; new implementation-notes/codex-1-reservation-recovery-20260914.md | codex-1 | integration/meta-protocol-change-evidence-first-efficiency | ../worktrees/evidence-first-integration | claimed |

The Claude checkout freezes source at 2ab9e82 and preserves the older failed
binding review checkout. This is an independent supporting source review, not
full-scope Phase-6 acceptance, a signature or an experiment arm. Claude uses its
existing roster identity, Opus 5 [1m], max effort, restricted native file tools,
a measured manual launch and a 30-minute process deadline. Codex cannot write
Claude's artifact. Runtime preparation/canaries remain facilitator-owned ignored
files. No source write boundary intersects the reviewer artifact.

Reservation recovery plan: preserve a complete precharge intent under the common
cycle guard before the original Store.Reserve. Pin original entry, epoch, effective
accounting contract, typed action identity and validated precharge state/archive.
Recover only an exact missing charged trajectory row, with deterministic preview,
hash-confirmed apply and exact replay. Preserve intent-without-charge, corrupt or
partial publications and all failed original attempts. Reconstructing accounting
must not infer a process terminal, inactivity, acceptance, a refund or permission
to retry; other helper/execution recovery remains required. Check the original
source archive and every other original charge/resolution. Refuse missing legacy
intent rather than guessing from today's source. Validate the charge/observer
interruption, changed identity/source/scope, replay, publication/output failure
and competing CLI recovery with real process fixtures.

Historical quorum, full-six pilot, packet gates and experiment funding questions
remain unchanged and unresolved. No treatment calls are authorized by a guessed
ceiling. The tool goal still reports an older blocked status; this continuation
is active work and is not a completion claim or a fresh three-turn blocked audit.


### September 14 measured review and reservation implementation checkpoint

Claude invocation dcb0984c-611b-4fe3-8ba5-61182be0642a has published requested
and started metadata in the new frozen review checkout. Source-review result and
artifact are pending; provider availability is not inferred from elapsed reset time.
No old failure/artifact was replaced. The launcher binary was built from 2ab9e82
in native temporary storage; the exact/atomic-file sandbox canary passed.

The new reservation boundary and deterministic missing-row CLI are implemented.
Initial focused budget/trajectory checks passed (29.403s). Wider budget/trajectory/
app integration selection passed (172.754s), including actual concurrent production
CLI applies, lost-output replay, unchanged observations, parent recovery and captured
verification. Subsequent source inspection corrected extended-policy compatibility;
the dedicated original/extended-ceiling test passed (2.526s). Full final validation
and independent acceptance are still required. Development and accepted final
checkpoints must remain distinct.


### Claude supporting source review completed — September 14

Invocation dcb0984c-611b-4fe3-8ba5-61182be0642a exited 0 after 1563.577s and
wrote its own supporting source-only review of 2ab9e82. Artifact SHA256:
cb6143ff207eeddf05d1be49b42a0a035d3e133a3a88de1f0380bc8c837597f0.
Codex imported those exact bytes from the frozen review checkout without editing
Claude's content. Only the assigned artifact changed there; native/shared measured
launch outputs match. This is not review or acceptance of the later reservation
slice, execution verification or a signoff.

Claude reports F1/F2 MAJOR, F3/F4 MINOR, F5/F6 NIT. F1 covers abnormal or
refused-after-charge unchanged attempts with no continuation path; it remains an
explicit owned recovery obligation alongside helper tickets. F2 identifies stale
preview hashes masking lifecycle-predicate omissions in tests; Codex is testing
the proposed counterexample before modifying those tests. F3 challenges activation
quorum binding; F4 concerns non-tracked Git excludes and remains UNVERIFIED in the
review. F5 narrows provider-failure wording; F6 flags unmeasured guarded history
cost. Findings are retained verbatim, not dismissed by the facilitator or treated
as waived by passing tests. Independent review disposition and final scope remain open.

The review reported USD 8.2361355 as a CLI estimate, with Opus 5 [1m] and Haiku 4.5
usage identities reported separately from the requested Opus configuration. Live
inventory now has 35 terminal attempts, 17 unknown costs and USD 54.424894 known
CLI estimates; total unknown. The earlier weekly-limit failure stays counted.
No packet/pilot treatment calls or quorum changes occurred.


### Reservation intent and exact missing-row validation — September 14

Accepted source manifest a5e1f06be6cff9a60322318d2b9dcdb591749ae471ac417c68d46cff97caf1c5
covers 407 Go/module files. Full PASS 337.880s (32 package terminals), six-package
race PASS 387.708s, vet PASS, Windows budget/trajectory/app cross-builds PASS;
Windows runtime unverified. Compiled shared selections PASS 0.556s / 16.633s /
8.412s. Ten new top-level tests and five intended negative failures are inventoried.
Native/shared bytes, unchanged source and historical HTML match. See the owned
reservation note and .parley-runtime/reservation-final-validation-20260914/.

This checkpoint is ready for independent review, not full-scope acceptance.
Claude's F2 counterexample was actually reproduced: two removed-qualification
overlays left the original contrary-lifecycle test passing. Its correction is
next; no prior green run is represented as detecting those missing clauses.
All helper/abnormal-terminal recovery, quorum/exclude-source findings, live
experiments, current-source reviews/signatures and delivery/report gates stay open.
