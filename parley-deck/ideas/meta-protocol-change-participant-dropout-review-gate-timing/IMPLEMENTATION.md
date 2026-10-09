---
idea: meta-protocol-change-participant-dropout-review-gate-timing
status: fix-up-cycle-1
implementer: codex-1
started: 2026-10-09
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing#review-gate-timing
head-commit: b89e2abaa056813a4b238c2fa4278195b0f7096b
design-pr: https://github.com/feci/parley-deck-cli/pull/77
implementation-pr: https://github.com/feci/parley-deck-cli/pull/78
skill-pr: https://github.com/feci/parley-deck-skill/pull/10
---

## Summary of work

Plan recorded before product edits. Implement frozen FINAL D1–D5 in the owner-allocated CLI and skill worktrees. Both design ACCEPTs exist and design PR77 merged at ac0d374.

## Implementation plan / checklist

- [x] Shared membership cause predicate: validated history/latest exact automatic transition or settled prospective decision, typed evidence allowlist, protected roles/current identities and distinct known snapshot models. Integrate precommit and driver review/goal/close gates without changing track policy.
- [x] Scoped300s participant-step stall default and supervised headless signoff command with real terminal watchdog classification and process cleanup. Correct buffered Zcode/default Claude text declarations; preserve explicit config, existing bounds and other launch paths.
- [x] Adversarial membership/driver/app/runner/adapter tests: positive and negative cause/model cases; prospective commit and kickoff; revisions/pending history; count-only auto exception; strict/reserved/goal/dissent/floor; batch rebind; real watchdog, buffered success/hard timeout, two-attempt/replay behavior.
- [x] Apply the14 frozen protocol substitutions identically to the three copies. Compact skill instructions while retaining all duties/required headings; update guidance/notices/version/changelog/compatibility metadata for CLI1.53.0 and skill2.17.0. Preserve caps/map. Record exact compaction/evidence.
- [x] Run focused tests, full Go suite with45m package timeout, build/vet, packet/drift guards and npm test. Use task-local native test temp with shared-volume cache; retain all failed and passing evidence.
- [ ] Publish both implementation PRs; separate Zcode full-scope review and own review consensus/signoff. Apply agreed fixes with max5 cycles. Never author a self-review. Fresh independent goal process, both ACCEPTs and current-tree evidence precede the authorized attended close.
- [ ] After close: clean-source builds and channel delivery; all19 runtime targets/114 SKILL.md hashes; exact staged-core reconstruction; separate Zcode channel verification; owner-only unpublished commands in2.15→2.16→2.17 order; released note and verified shared memory.

## Deviations from FINAL.md

None. Adjacent test/helper filenames may follow the smallest implementation seam, as FINAL permits. No policy schema/track/accounting/credentials/worktree-allocation changes.

## Notes for reviewers

Review every FINAL criterion with refutation attempts; no finding is suppressed by a prior disposition. Main risks are false causal credit, precommit/close divergence, model-identity ambiguity, a false permanent buffered dropout, signoff cleanup/classification and compaction rule loss. Product base85c5a8f; design-only merge ac0d374. Skill base dbdb919. Private commands/results live under .parley-runtime/review-gate-timing; independent artifacts remain participant-owned.

The controlling brief permits THIS run's attended close only after no open CRITICAL/MAJOR, both current review-consensus ACCEPTs, a fresh Zcode goal PASS and current-tree AC evidence. It has no automatic history/snapshot and does not exercise the new exception itself. Retry policy remains the brief's15min wait and at most8 relaunches on own/Zcode429/503/timeout; timeout relaunch2400s then3600s. AC14 is post-close delivery.

## Initial implementation — ready for independent review

The product follows D1–D5. `internal/membership/review_gate.go` derives causal qualification from inspected history or the settled prospective decision; gates and driver consumers share it. The runner reuses its supervisor for eligible signoffs, including terminal classification after cleanup, and scopes the new stall default to participant steps. Zcode and default Claude text are buffered; Codex/Kimi remain streaming. No Claude task was launched.

One adjacent seam was necessary: standalone signoff runs in `internal/app/quota_signoff.go` inherit the existing history run's roster snapshot/revision, and signoff launches apply that snapshot. A new orchestration run must not discard the configured model basis needed at settlement. Missing snapshots remain missing. This adds no policy, schema or model observation claim.

Skill implementation is commit 74cc831 in the sibling worktree, based on dbdb919. CLI product base is 85c5a8f. Exact protocol and sequential skill transformations are reproducible from source-context/protocol-hunks.json, skill-compaction.json and implementation-compaction-proof.json. All 14 substitutions match once in each protocol and staged-core preview. Phase-1 guard is 68,846/70,000 B; skill is 18,990/20,000 B. Applicability map unchanged. The 2.17.0 core preview is not yet staged/published; AC14 follows close.

Validation evidence is retained privately under `.parley-runtime/review-gate-timing/`. Focused new tests passed in membership, driver, app, runner and agents; build/vet and actual packet/drift checks passed. Initial full-suite run v1 exposed two regression integration details: the goal execution fixture needed an idea prompt now that cause validation inspects membership, and the notice must retain its LE-7/LE-11 names. Both were corrected and the exact two tests pass (validation-fix-v1-result.json, exit 0). Initial skill suite failed against stale payload hashes; the manifest was regenerated and the complete suite rerun. Failed results remain preserved as failed results, never counted as passes. Final committed-source full validation is pending and must pass before close; its results will be appended here.

No independent code review or completion verdict is claimed by codex-1. Zcode's separate full-scope review will assess both product diffs and every AC, with unrestricted finding scope. Fix-up cycles used: 0 of 5.

The complete skill v2 run passed 385 Node tests but failed loading the missing commonmark dev dependency; this is not a passing suite. npm ci and the complete v3 suite are running. The final current-source Go suite is full-host-tests-v2 at a26f588, with build/vet afterward. Initial Go v1 remains failed history.

## Committed-source validation completed

CLI product a26f5881835560dfabfe483a6d537aff39c45c35, skill 74cc831b18ce33e48e705fc8992c89be85cadf6f. Later CLI changes are idea evidence only. source-context/implementation-validation.json records commands, counts, exits and log identities.

- `go test ./... -json -count=1 -timeout 45m`: actual exit0 in847.270s; 34 passing packages, 3 with no tests, 3311 passing test/subtest events, 4 built-in skips, zero fail events. Log full-host-tests-v2.jsonl SHA256 04e3fa050564d222f09fc7d0b08bc665a0c26e6e180c28dec1cd28e262c28d23. This is the qualifying full run; failed v1 stays failure history.
- Receipt storage anomaly disclosed: the process emitted its exact JSON receipt to the tool, but the shared volume lists full-host-tests-v2-result.json while open/stat reports ENOENT. The identical emitted receipt is preserved as full-host-tests-v2-recovered-result.json, with provenance in full-host-tests-v2-receipt-recovery.md. The full3.4MB JSONL is readable and its SHA256 verified. No result was inferred or new execution claimed. Independent review must assess this evidence, not treat it as an exemption.
- `go build ./...` and `go vet ./...`: exit0 on the same product; build-v2/vet-v2-result.json retain the actual receipt.
- npm ci then complete npm test: exit0 at skill74cc831; 399 Node tests, 54 Python tests and six manifests. Log skill-tests-v3.log SHA256 3509b5ce4c7db449fb6f72ab9055215ce6ad7d0c29bbe367ff83a0aa0cc35a4b. Failed v1/v2 are retained.
- Packaging preflight: npm pack dry-run, native portable build/version and all-target install dry-run pass. All210 packed files match skill source bytes. Private release preparation: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-review-gate-timing; no publication/runtime installation yet.

These are implementer-run checks. Independent review, both ACCEPT blocks and the separate fresh goal process remain outstanding; no self-review or self-certified exemption.

## Fix-up cycle1 plan

The revised review/consensus.md has both ACCEPTs. Preserve the original-manifest integrity gate confirmed against baseline by Zcode; distinguish it from absent snapshots, warn when a standalone signoff's valid manifest lacks a snapshot, and cover absent/missing/corrupt/foreign variants. Add legacy short/contradictory reset, adapter and provenance negatives at the shared predicate. Insert the Markdown blank line. Re-run covering checks and full current-source validation, then independent full-scope review. No policy/historical recovery relaxation.

## Fix-up cycle 1
status: complete
completed: 2026-10-09

### Fixes applied

The signed revised cycle-1 plan is implemented. The snapshot-inheritance comment now explicitly distinguishes an absent snapshot in a valid manifest from a missing/corrupt/foreign manifest. Standalone signoff emits a warning when the inherited snapshot is empty; it creates no model identity and cannot qualify the exception. End-to-end tests cover absent snapshot through successful signoffs and empty inheritance, plus missing/corrupt/foreign manifests failing for their expected integrity reasons before a child starts. A membership test removes the snapshot after a valid causal dropout and confirms refusal. The existing same-tick rebind test passes.

The legacy evidence seam now tests all three allowed rules at their valid boundary plus wrong adapter, wrong provenance, short reset or contradictory no-reset evidence. The documentation heading has its missing blank line.

### Deviations from agreed fixes

None from the revised, doubly signed plan. The initial proposed relaxation was withdrawn by its author before Codex signed, with the baseline proof and prior draft preserved. No policy/history recovery behavior changed.

### Validation and next gate

Focused app/membership TestReviewGate run v3 passes (16.733s, log SHA2568548172a7dc135e57fef6063131f49535fbc33c5b1cbb62ac0671c71cb048c41). Failed v1/v2 are retained: the new fixture initially overwrote immutable prompt metadata and then lacked attributable consensus drafter metadata; both setup defects were corrected. Full final-source Go/build/vet and a separate full-scope Zcode review follow this commit. Skill product74cc831 and its full399 Node/54 Python/six-manifest pass are unchanged. Cycles used:1 of5.

Cycle1 product commit: b89e2abaa056813a4b238c2fa4278195b0f7096b. Full-source validation log: full-host-tests-cycle1.jsonl with checked-durable receipts beside the log and in the native task temp validation-receipts directory.
