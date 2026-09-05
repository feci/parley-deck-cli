---
idea: meta-protocol-change-evidence-first-efficiency
status: in-progress
implementer: codex-1
started: 2026-09-05
branch: parley-deck-cli#integration/meta-protocol-change-evidence-first-efficiency
head-commit: 3a09a0cf2ef938e2456b50ec6eb34a6b5f37038c
design-pr: https://github.com/feci/parley-deck-cli/pull/72
implementation-pr: pending
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
| internal/protocolpacket/**; internal/app/protocol_packet.go; internal/app/protocol_packet_test.go; parley-deck/meta/packet-applicability.yaml; parley-deck/COOPERATION.md outside generated Section 2; parley-deck/meta/protocol-changelog.md; owned claude handoff | claude-1 | feature/meta-protocol-change-evidence-first-efficiency/claude-1 | ../worktrees/evidence-first-claude | claimed |
| internal/app/preflight.go; internal/app/preflight_test.go; internal/app/preflight_hash_test.go; internal/app/preflight_liveness.go; internal/app/preflight_liveness_test.go; internal/runner/supervision.go; internal/runner/buffered_test.go; owned hermes handoff | hermes-1 | feature/meta-protocol-change-evidence-first-efficiency/hermes-1 | ../worktrees/evidence-first-hermes | claimed |
| internal/evidence/**; internal/app/driver_checks.go; internal/app/driver_checks_test.go; internal/app/driver_evidence.go; internal/app/driver_evidence_test.go; owned kimi handoff | kimi-1 | feature/meta-protocol-change-evidence-first-efficiency/kimi-1 | ../worktrees/evidence-first-kimi | claimed |
| Skill repo: skills/parley-deck/SKILL.md; skills/parley-deck/references/COOPERATION.md; test/packet-context.test.js | claude-1 | feature/meta-protocol-change-evidence-first-efficiency/claude-1 | ../worktrees/evidence-first-skill | claimed |

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
