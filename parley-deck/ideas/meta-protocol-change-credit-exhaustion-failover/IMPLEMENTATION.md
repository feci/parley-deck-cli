---
idea: meta-protocol-change-credit-exhaustion-failover
status: implemented
implementer: codex-1
started: 2026-10-29
completed: 2026-10-29
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/credit-failover#feature/meta-protocol-change-credit-exhaustion-failover
head-commit: 24c6184332ff4fc8084ead7458c1c2b207b949f3
design-pr: https://github.com/feci/parley-deck-cli/pull/86
implementation-pr: n/a
---

## Summary of work

Implementation starts from merged design commit `24c6184332ff4fc8084ead7458c1c2b207b949f3`. No product change precedes this plan.
Codex-1 organizes, drafts, and implements. Kimi-1 and Zcode-1 independently review and goal-check.
The skill worktree remains on its allocated `credit-failover` branch.

## Implementation plan / checklist

- [x] AC-01/02/11: add v2 policy and typed exhaustion proof; preserve legacy bytes, hashes, and behavior.
- [x] AC-03/04/05: bind invocation/validation proof and common admission; preserve transient retry behavior and code-step freshness.
- [x] AC-06/07: freeze roles/order and commit whole-batch succession through the existing immutable history.
- [x] AC-08/09/10: add the checked active execution resolver, successor handover, draft recovery, and causal review proof.
- [x] AC-12: add machine-only private Zcode capture, bounded cleanup, status, and purge.
- [x] AC-13/14: apply the exact signed hunks; update skill guidance, hard size tests, versions, and release notes.
- [ ] Checks: focused tests and fault cases, relevant race tests, full Go suite, skill suite, package checks, portable builds.
- [ ] AC-15: independent review, signed review consensus, agreed fixes, fresh independent goal check, and current-tree AC evidence.
- [ ] Release: CLI 1.56.0 and skill/core 2.20.0; GitHub, both Homebrew formulae, skill-only WinGet, all runtime installs.
- [ ] Separate reviewer process verifies delivery. Report owner-only npm/core commands without execution.

Review risks: suppression must survive restarts and blocked batches. Every current execution consumer must use checked history.
A stale valid artifact must not excuse a failed new code step. Pre-implementation handover must not imply review readiness.
Retain predecessor bytes, all dissent, and all prior findings. A successor never reviews its own code.

## Deviations from FINAL.md

No decision change. Added a shared telemetry admission helper so membership cannot count stale artifacts from suppressed identities.
Added `internal/quotatest/exhaustion.go`, focused capture/succession/handover tests, and `docs/credit-exhaustion.md`.
The skill adds `references/EXHAUSTION_AND_HANDOVER.md`; existing references and templates link to it.
The capture scrubber accepts only proven safe framing. Uncertain input creates no file, as D10 requires.

## Notes for reviewers

Use FINAL AC-01 through AC-15 and decisions D1–D12. The exact 29-hunk SHA256 is
`f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
The historical DISPUTED interpretation supplies no decision or AC premise.
Codex support is source-backed at exact CLI version 0.165.0 with explicit JSON argv.
The ordinary smoke does not prove a native exhaustion capture or binary equivalence.
Zcode keeps its existing grammar and native-positive limitation.

## Progress

- 2026-10-29 20:41Z: design PR 86 merged with all three ACCEPT signatures.
- 2026-10-29 20:41Z: implementation plan filed before product edits.
- Baseline package checks passed: quota, telemetry, membership, runner, and protocol.

- Resume continuation: all design rounds and three ACCEPT signatures remain intact. No participant relaunch was needed.
- Added v2 exhaustion evidence, durable admission, owner-order role history, active execution resolution, draft recovery, and private capture.
- Core integration passed: quota, membership, protocol, telemetry, runner, and driver.
- New tests passed for role selection, execution chains, draft integrity, capture source, scrubbing, permissions, symlinks, retention, and capacity.
- Applied the exact 29 signed hunks to all three protocol copies and staged core 2.20.0.
- Added shared synthetic evidence fixtures under `internal/quotatest/`; they make no provider calls.
- Succession integration, packet limits, package guidance, and native portable build passed. Full Go checks and independent review remain pending.

## Decision Log

The phase-5 driver attempt reports `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`.
The existing owner answer authorizes native fallback. It supplies no driver accounting or automatic membership authority.
Use the same selected identities, models, effort, timeouts, retry limits, and independent roles.

## Surprises & Discoveries

No new product finding at plan publication. The execution consumer inventory includes evidence, status, TUI, and retrospective paths.
The requested all-home runtime scan stalled in unavailable CloudStorage. The completed local inventory excludes that location.
Do not claim all-home completeness or overwrite source/archive copies.

## Validation evidence

Pre-change command:
`go test ./internal/quota ./internal/telemetry ./internal/membership ./internal/runner ./internal/protocol -count=1 -timeout 10m`.
All five packages passed. This establishes a baseline only; it proves no new AC.
Initial package checks passed after exhaustion admission, role history, handover, and capture changes.
The full Go suite, package tests, and independent review remain pending.

- Current focused tests pass for suppression, stale code output, handover dispatch/replay, whole-batch succession, retained findings, and owner widening.
- Targeted race checks pass for quota, membership, runner, and telemetry.
- Packet tests: 68,448 bytes unflagged; 69,880 bytes with `protocol_change`. Both use the test source display fields.
- Skill `npm test` passed: 402 Node tests and 54 Python tests. Payload manifests pass.
- Skill size is 19,989 bytes. Package dry-run and native portable build passed.
- CLI native build reports 1.56.0. Capture status reports disabled by default and no retained files.
- The full Go suite is in progress. It must pass at the reviewed tree before close.

## Outcomes & Retrospective

Pending implementation and independent review.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none
