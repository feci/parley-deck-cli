---
idea: meta-protocol-change-evidence-first-efficiency
status: in-progress
implementer: codex-1
started: 2026-09-05
branch: parley-deck-cli#integration/meta-protocol-change-evidence-first-efficiency
head-commit: 8ccd1ea
design-pr: https://github.com/feci/parley-deck-cli/pull/72
implementation-pr: https://github.com/feci/parley-deck-cli/pull/73
---

# Evidence-First Delivery

## Summary of work

Design is merged at 3a09a0cf2ef938e2456b50ec6eb34a6b5f37038c with all four
owned signatures. Production implementation begins here. No production feature,
live pilot, packet measurement, independent final acceptance or deployment is
claimed complete. The frozen historical evaluation is preserved.

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

1. Provision the three claimed CLI worktrees and the sibling skill worktree.
2. Implement disjoint slices and collect owned handoffs with actual test results.
3. Integrate and independently review each slice; implement agreed fixes.
4. Build an uninstalled binary, verify real telemetry, freeze and run experiments.
5. Populate/verify the single HTML, complete independent acceptance and the
   protocol review cycle before merging the implementation PR.

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

Pending production implementation. The separate pilot harness has real passing
tests on Python 3.9 and 3.11 and independent reviewer reruns, but no live treatment
has been executed. Claude's prior quota failure is preserved; his post-reset
signature succeeded on the same requested model without a route change.

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
