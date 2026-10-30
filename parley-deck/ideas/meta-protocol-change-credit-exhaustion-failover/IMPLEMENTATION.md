---
idea: meta-protocol-change-credit-exhaustion-failover
status: complete
implementer: codex-1
started: 2026-10-29
completed: 2026-10-30
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/credit-failover#feature/meta-protocol-change-credit-exhaustion-failover
head-commit: b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666
design-pr: https://github.com/feci/parley-deck-cli/pull/86
implementation-pr: https://github.com/feci/parley-deck-cli/pull/87
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
- [x] Checks: focused tests and fault cases, relevant race tests, full Go suite, skill suite, package checks, portable builds.
- [x] AC-15 pre-release gates: independent review, signed review consensus, agreed fixes, fresh independent goal check, and current-tree AC evidence.
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

CLI implementation commit: `931550cb4c57f041f9a8a1d1337571cb271a873a`. Skill implementation commit: `70d29c1aebd9a632f6f6239123fd2f52e3775a48`.


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
- The reviewed-source full Go suite exited 1. VERSION and unreadable-implementation phase checks failed.
- Kimi and Zcode filed independent round-01 reviews. Review consensus proposes six fixes.
- The capture unsafe-mode fixture also fails under umask 0077. Fix-up and current-source verification remain required.

## Outcomes & Retrospective

Both independent reviewers resolved all findings after two fix-up cycles. All three participants accepted the zero-fix consensus.
Zcode issued `GOAL-CHECK: PASS` in a fresh process at the reviewed product trees.
The intermittent admission race required repeated concurrent tests; one full-suite pass could not dispose of its reproduced counterexample.
The checked publication path now preserves complete records, exclusive refusal, and directory durability.
Release follows this implementation close under FINAL D12. Its evidence belongs to the separate delivery record.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none

## Fix-up cycle 1
status: complete
completed: 2026-10-29
head-commit: 31c039c41a455a03385401085d132234a17603f6

### Fixes applied

- F1: separate protected roles from fallback/default candidates at both protection sites. Only exhaustion-backed losses request succession.
  Selection uses the surviving membership after every settled loss. Immutable history re-derives the same result.
  Readiness tests count children and check five-second spacing. They verify ordinary global-default fallthrough through the actual driver resolver.
  Settlement tests cover mixed losses, both-role loss, an unusable counterpart, explicit designees, pins, and immutable kickoff validation.
- F2: root VERSION now says 1.56.0, as the CLI binary does.
- F3: retain positive stat evidence when ExecutionStarted returns an error. Preserve the integrity error for the caller.
- F4: force the unsafe fixture to mode 0644 and check setup errors. Tests pass under umasks 0077 and 0022.
- F5: wait names the active execution record. The integration test covers absent records and a validated handover before Phase 5.
- F6: remove the empty Unreleased heading. Preserve the release entries.

### Deviations from agreed fixes

F1 retains the eager read-only candidate calculation in app.go and membership/succession.go.
It grants no transient protection, pin, stored role assignment, or succession authority.
An exhaustion-backed transition binds the candidate identity and source in immutable history.
This preserves the existing resolver and avoids a second resolver or callback.
The reviewers must assess this timing detail against F1's instruction to resolve and bind only when succession requires it.

The completed Windows CI log exposed one new capture-test expectation mismatch.
Windows cannot report the required POSIX private modes through this storage check.
The test now expects refusal and verifies that no capture file exists. It does not skip the failure or weaken storage checks.
The CLI guide states this limit. No Windows ACL implementation enters this scope.

### Current verification and remaining gates

Focused checks passed at the fixed source, including the new wait and portable readiness fixtures.
The quota, membership, runner, and telemetry race checks passed after the F1 changes. go vet ./... passed.
The full command is active at this product commit: `go test ./... -count=1 -timeout 45m`.
Its private log is `fixup01-full-go-v2.log`. It must pass before close.
The earlier `fixup01-full-go.log` was stopped after a new test fixture needed correction. It is not acceptance evidence.
The corrected wait fixture now carries the frozen policy projection; its focused test passes.
Skill source remains 70d29c1aebd9a632f6f6239123fd2f52e3775a48. Its prior independent suite and package results still apply.

Windows run 40956924144 completed with failure at the previous product source.
Its runner test timed out at 45 minutes; private-mode, durable directory sync, and other platform tests also failed.
Many affected source files are unchanged by this idea. This record does not claim a clean Windows baseline or Windows support.
The native Windows capture refusal assertion needs hosted execution. CLI Windows remains experimental; CLI WinGet stays held.

Fresh full-scope Kimi and Zcode reviews, review consensus, and an independent goal check remain required.
AC-15 release and installation evidence remains pending. No release or installation has occurred.

### Review-round-2 concurrent admission observation

One Ubuntu job failed at this unchanged product source; the corresponding push job passed.
A 100-run local repetition reproduced `TestStepBudgetConcurrentIndependentGroupsShareCap` failures.
Diagnostic test-only output reports `unreadable invocation request during exhaustion admission`.
The record writer exposes requested.json before its write completes. Common admission can read those incomplete bytes.
The product remains at 31c039c while both independent reviewers work.
A private atomic-publication proposal is under test; no product fix is accepted or applied yet.
The next review consensus must address this result before close, regardless of any one full-suite pass.
Evidence and proposal: `../../inbox/codex-1-to-all_meta-protocol-change-credit-exhaustion-failover_review02-concurrent-admission.md`.

### Review-round-2 current state — 2026-10-29 22:45Z

The required full Go command completed exit 0 at 31c039c, with 34 passing packages.
Receipt: private `fixup01-full-go-v2.receipt.json`. The repeated concurrent-admission failure still blocks acceptance.
Zcode independently reproduced that failure in 7 of 100 runs. Its supplemental review rates it CRITICAL.
Zcode accepted the private writer-side atomic-publication proposal and the capture-status documentation correction.
Kimi reached its 1800-second ceiling without an artifact. Retry 1 may start after 22:55:50 UTC with 2400 seconds.
Its process and test child stopped. Native timeout recovery changes no membership or driver accounting.
Review consensus, the accepted fix-up, fresh review, goal check, and all release work remain open.


## Fix-up cycle 2 plan — 2026-10-29 23:34Z

All three participants accepted F7-F9. Both native signers completed exit 0 without timeout.
Exact-prefix and receipt-hash checks passed. Zcode also placed its source assessment in its own review.
The second fix plan preserves the original protocol hunks and FINAL decisions.

- [x] F7: stage complete private invocation records, then publish exclusively and sync both directory levels.
- [x] Reuse the quota directory-creation helper for the invocation store and its missing ancestors.
- [x] Recheck the store barrier before reservation, including after an earlier incomplete creation.
- [x] Add request/terminal fault tests, collision checks, staging cleanup, and inert crash-residue cases.
- [x] F8: state the capture-status cleanup effects in both references and command help.
- [x] F9: format the driver import and check all changed Go files.
- [x] Run persistence, grouped accounting, repeated concurrency, race, full Go, skill, package, and portable checks.
- [ ] Obtain fresh full-scope reviews, review consensus, and the independent goal check before release.

The store barrier also covers Zcode's source addendum about openInvocationStore.
DurableMkdir exposes the existing quota helper without changing that helper's behavior.
New invocation publication retains every error and refuses every existing destination.
These details implement F7; the reviewers must assess them against D3 and the agreed plan.
The phase-8 driver refused: `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`.
The authorized native fallback applies. It supplies no driver accounting or membership authority.


## Fix-up cycle 2
status: complete
completed: 2026-10-30
head-commit: b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666

### Fixes applied

F7 publishes complete private records through an exclusive hard link. The write, file sync, close, and directory barriers return errors.
The checked directory helper also creates missing store ancestors. An existing store receives a fresh durability barrier before invocation reservation.
No record reader consumes the private staging names. Published malformed records still refuse admission.
Request and terminal tests inject write, sync, close, publication, and directory failures. Store-sync failure prevents reservation.
Collision tests preserve files, symlinks, and directories. Failed terminal persistence remains an error on repeated Finish calls.

F8 explains status cleanup in CLI help, the CLI guide, and the packaged reference.
Skill commit: 406b6976cc75344d76f908297fe46ec2443e62cb. Its generated payload manifest reflects the reference change.
F9 places the protocol import with project imports. Every changed Go file passes gofmt.

### Deviations from agreed fixes

None. Reuse of the quota directory helper covers the store-creation site in Zcode's supplementary source assessment.
No helper behavior changes. No protocol/core hunk or size limit changes.

### Current verification and remaining gates

Persistence tests, relevant race tests, go vet, package dry-run, and the native portable build passed.
The repeated concurrent/grouped check, complete Go command, and skill suite are in progress.
Private evidence uses the fixup02 prefix. No pending check is recorded as passed.
Fresh full-scope reviews must inspect the entire CLI and skill changes, including both fix-up cycles.
The fresh independent goal check and all release/install evidence remain open.


### Fixed-source verification complete — 2026-10-29 23:46Z

The complete command `go test ./... -count=1 -timeout 45m` passed at b7dce40, with 34 packages and exit 0.
The 100-run concurrent-budget and grouped-accounting checks passed. Relevant race checks and go vet passed.
The skill suite passed 402 Node tests and 54 Python tests at 406b697. Package checks and the native portable build passed.
All changed Go files pass gofmt. All four protocol/core hashes retain the signed bytes; SKILL.md remains 19,989 bytes.
Content-free receipts and the full Go package summary are in `source-context/fixup02-validation.json`.
Both fresh independent review processes remain active. These execution results grant no independent review or close verdict.


## Implementation close — 2026-10-30

Status: complete for Phases 5–8. Two fix-up cycles resolved F1–F9 and both audit-record NITs.
All three current participants filed plain ACCEPT in cycle-3 review consensus. No agreed fix remains.
Kimi and Zcode supplied full-scope current-tree reviews at CLI `b7dce40` and skill `406b6976`.
Zcode then issued `GOAL-CHECK: PASS` in a separate native process, with exit 0 and no timeout.
Its full assessment remains in `review/round-03/zcode-1.md`; its receipt and prefix check remain in `source-context/`.
The required full Go suite passed at the reviewed product: 34 packages. The skill suite passed 402 Node and 54 Python tests.
Current validation, limits, and command-bearing receipts remain in `source-context/fixup02-validation.json` and the verifier-owned reviews.

The product source remains CLI `b7dce40a1bf1dd4d3ba32dc792b98891ed4a9666` and skill `406b6976cc75344d76f908297fe46ec2443e62cb`.
Later commits contain audit records only. Both implementation PRs now qualify for merge commits under the signed gates.

### Ordered release handoff

The open release checklist items below the initial plan are deliberate. D12 orders their execution after this close.
This complete status makes no claim that channels or runtime installation already passed.
Continue with clean merged-source builds, both GitHub releases, both Homebrew formulae, skill-only WinGet, and all 19 runtime destinations.
Then run the separate channel verifier and write `../../inbox/codex-1-to-user_meta-protocol-change-credit-exhaustion-failover_released.md`.
That delivery record carries the remaining AC-15 channel evidence, owner commands, deferred scope, and usage.
Preserve this closed execution record and frozen FINAL. Append later channel assessment only to the verifier's own review artifact.
Never execute npm or global-core publication. CLI Windows remains experimental; CLI WinGet stays held.

The phase-8 driver refused: `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`.
The owner's native fallback authority remains in force. It supplies no driver accounting or automatic membership authority.
The host reboot caused no participant exclusion.
