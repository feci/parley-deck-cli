---
idea: meta-protocol-change-participant-dropout-review-gate-timing
status: ready-for-review
implementer: codex-1
started: 2026-10-09
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing#review-gate-timing
head-commit: ac0d374ded57f37b68fc2c31bd89ea2f552885bb
design-pr: https://github.com/feci/parley-deck-cli/pull/77
implementation-pr: pending
---

## Summary of work

Plan recorded before product edits. Implement frozen FINAL D1–D5 in the owner-allocated CLI and skill worktrees. Both design ACCEPTs exist and design PR77 merged at ac0d374.

## Implementation plan / checklist

- [x] Shared membership cause predicate: validated history/latest exact automatic transition or settled prospective decision, typed evidence allowlist, protected roles/current identities and distinct known snapshot models. Integrate precommit and driver review/goal/close gates without changing track policy.
- [x] Scoped300s participant-step stall default and supervised headless signoff command with real terminal watchdog classification and process cleanup. Correct buffered Zcode/default Claude text declarations; preserve explicit config, existing bounds and other launch paths.
- [x] Adversarial membership/driver/app/runner/adapter tests: positive and negative cause/model cases; prospective commit and kickoff; revisions/pending history; count-only auto exception; strict/reserved/goal/dissent/floor; batch rebind; real watchdog, buffered success/hard timeout, two-attempt/replay behavior.
- [x] Apply the14 frozen protocol substitutions identically to the three copies. Compact skill instructions while retaining all duties/required headings; update guidance/notices/version/changelog/compatibility metadata for CLI1.53.0 and skill2.17.0. Preserve caps/map. Record exact compaction/evidence.
- [ ] Run focused tests, full Go suite with45m package timeout, build/vet, packet/drift guards and npm test. Use task-local native test temp with shared-volume cache; retain all failed and passing evidence.
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
