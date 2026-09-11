---
idea: meta-protocol-change-evidence-first-efficiency
status: in-progress
implementer: codex-1
started: 2026-09-05
branch: parley-deck-cli#integration/meta-protocol-change-evidence-first-efficiency
head-commit: 78a7b13
design-pr: https://github.com/feci/parley-deck-cli/pull/72
implementation-pr: https://github.com/feci/parley-deck-cli/pull/73
---

# Evidence-First Delivery

## Summary of work

The signed design is merged at 3a09a0cf2ef938e2456b50ec6eb34a6b5f37038c.
Implementation is partial on draft PR #73: measured protocol-task launch
boundaries, live/full context attestation, cursor persistence, fail-closed goal
checking, packet publication and source-instruction updates are integrated and
locally tested. Claude's partial independent source reviews are preserved with
their dispositions. Kimi's completion-transition source is integrated and wired
through the actual independent helper. The post-completion evidence gate and
original-contract deletion/restart counterexamples now pass after corrections.
Readiness schema/capture corrections, shared launch/step/cycle policies, finite
operator cycle extensions and monetary-default enforcement are integrated.
Source 78a7b13 refuses a configured dollar ceiling before an unreserved first
process launch and requires a matching persistent operator policy. The full Go
suite, scoped race, vet and shared-volume fixtures pass. Legacy migration,
launch/step extensions, safe guard recovery, durable semantic action replay,
canonical refusal publication/recovery, regression trajectory, actual live
experiments and final independent acceptance remain open. The frozen historical
evaluation is unchanged.

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

1. Complete Kimi's mixed package/build-failure correction after the 58-event
   independent pass. The preceding three issues are fixed in uncommitted owner
   source; the new masked-exit package-failure probe still requires correction.
2. Wire the tested shared reservation store into per-launch/per-action budgets,
   persistent loop limits and manual/driver/resume/BLOCK paths. Complete readiness
   validation and independent re-review of the corrected terminal process path.
3. Reconcile at least 20 actual attempts across launch surfaces. The recovery
   inventory currently has fourteen terminal attempts, five with unknown cost.
4. Obtain the pending historical Hermes-to-Zcode decision and pilot amendment
   direction. Freeze exact tasks/resource ceilings before either experiment;
   preserve the signed packet trial and full-six design until lawfully amended.
5. Finish independent per-slice review and the full Phase-6/7/8 review/signoff
   cycle. Keep PR #73 and the source-skill changes in draft; no final merge,
   global installation, immutable-core publication or release is authorized.

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
