---
idea: meta-protocol-change-evidence-first-efficiency
status: in-progress
implementer: codex-1
started: 2026-09-05
branch: parley-deck-cli#integration/meta-protocol-change-evidence-first-efficiency
head-commit: 47661e5f7d72d81ef864baf9825fc223e91e7dc0
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

### 2026-09-10 — recovered interrupted Codex session

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
attempts remain untouched. Both resumed participant handoffs are pending.

No raw credentials, prompts or private session output are copied into these
canonical records. The original session's workspace spend-cap failures are
historical; no current failure is labelled a spend-cap failure without evidence.

Use this manifest and actual git worktree/branch state. Never delete or reset an
unexpected worktree. Keep failed agent attempts and authored artifacts. The design
FINAL is immutable; implementation details and genuine deviations belong here.
No package, global skill/core installation or production deployment is requested.

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
