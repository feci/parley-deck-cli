---
idea: meta-protocol-change-driver-unstall
review-cycle: 1
outstanding_agreed_fixes: 3
blocked: false
drafted-by: codex-1
date: 2026-10-09
reviewed-commit: c659bc84b44203681cbdc4e5f68d135078253113
---

## Agreed fixes

Proposed for both current participants' own signoffs before implementation:

1. R1-F3 (Zcode NIT): document the exact Git-common durable legacy record path in docs/legacy-history.md.
2. R1-F4 (Zcode NIT): add the malformed-track refusal to the README goal-check summary.
3. CI-TIMEOUT (additional implementer-discovered validation defect): the macOS PR job failed `TestGoalUnstallBufferedHardDeadlineAndCancellation` at the shell counter assertion (`hard timeout attempts: "child\n"`), while the corrected local full suite, macOS push CI and ten focused repetitions passed. A shell may be launched and killed by the 100ms ceiling before completing stdin consumption and its counter append. Count attempts using immutable terminal telemetry with a positive StartedAt/PID, explicit timeout failure, distinct attempt ordinals and the same frozen ceiling. Replaying must preserve the exact two invocation identities and never create a third. Keep the hard-deadline/failed-close checks and cancellation refusal; retain the shell counter only as supporting evidence bounded by the number of actual starts. This proposal is a test correction, not a change in product timeout/retry behavior. The reviewer may challenge the diagnosis or fix.

## Finding dispositions and close conditions

- R1-F1 (Zcode MINOR): the corrected full suite has now exited 0 at reviewed commit c659bc8: source-context/go-full-fixed.log, app686.116s, trajectory685.676s. The earlier failed version-race run is retained as source-context/go-full.log; it is not reported as passing. Skill399 Node/54 Python/six manifests passed. The new macOS CI failure above remains open until fixed and checked. Final current-tree verification and CI results are required after the fix.
- R1-F2 (Zcode MINOR): accepted as a binding sequencing condition. Pre-merge independent goal verification covers AC1–AC8 and AC9's checks/drift portion. AC9 live release/install evidence stays explicitly PENDING in the eventual implementation-complete record. A separate fresh Zcode process must verify delivered channels after source merge/publication/installation and before any released handoff is written. A textual pre-merge goal PASS never claims live AC9 satisfaction. This follows the controlling brief's source-review → merge → release sequence; the design FINAL is unchanged.
- The actual `membership.SingleReviewerAfterDropout` probe returned allowed=false, error empty (source-context/reviewer-gate-result.json). The brief-authorized attended close may be used only after the final independent review has no open CRITICAL/MAJOR, both current participants ACCEPT, a fresh independent goal PASS and current-tree independent AC evidence. No product exception or automatic history is fabricated.

## Deferred follow-ups

- Zcode open question1: owner recovery after both durable attempts are exhausted. TBD follow-up; current contract escalates or needs an owner ruling/new idea, without resetting the attempt budget.
- Zcode open question2: cycle-scoped goal identity for a future explicitly authorized re-check. TBD follow-up; this release preserves the ratified durable logical step and does not authorize new attempts from changed run/config/input. Current-tree independent criterion evidence remains a separate close duty; an old textual PASS alone cannot close changed work.
- The design FINAL's phase-pointer/planner, ready/partial consensus reopen and other listed driver gaps remain deferred, with their recorded hazards.

## Dismissed findings

None. F1/F2 are explicit process conditions, F3/F4 are accepted fixes. CI-TIMEOUT is an additional proposal disclosed after review, not attributed to Zcode.

## Coverage & blind spots

Zcode independently probed AC1–AC7, ran the real121-second witness and compared all protocol/core copies; no CRITICAL/MAJOR. One independent model remains after Kimi's two403 failures. Agreement is not evidence of correctness; current-tree adversarial checks and the later independent channel process supply the evidence. macOS PR CI demonstrated why local green alone is insufficient.

Role concentration: codex-1 organizes, implements and drafts this consensus; Zcode independently reviews. Both own signoffs govern the transition.

## Drafter position changes

The last canonical round, round-02/codex-1.md, proposed D1/D2 with fail-closed accounting and derived bounded goal execution; those product positions are unchanged. New implementation-only position: accept Zcode F3/F4 and correct the CI test's start-count oracle. IMPLEMENTATION.md's prior statement "Full test results remain a merge gate" is retained; now the fixed-tree local full run is green and macOS PR CI provides a further defect to resolve.

## Signoffs

<!-- Each current participant APPENDS their own signoff. -->

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: Accept the three bounded fixes and all explicit close conditions. This is agreement to fix, not implementation completion or live AC9 delivery.

### Signoff: zcode-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: Accept the three agreed fixes and every stated close condition; this is agreement to fix, not completion. Independent verification before signing (PRIMARY): source-context/go-full-fixed.log ends green at reviewed commit c659bc8 (internal/app 686.116s, internal/trajectory 685.676s, zero FAIL lines), closing F1's local gate while the failed version-race run stays honestly on record; source-context/reviewer-gate-result.json records allowed=false with empty error, so the attended close path and its four conditions remain the only route — no mechanical exception is claimed; source-context/ci-macos-pr-failed.log:199-205 reproduces CI-TIMEOUT (TestGoalUnstallBufferedHardDeadlineAndCancellation, `hard timeout attempts: "child\n"`), and I concur with both the diagnosis (shell killed by the 100ms ceiling before consuming stdin and appending its counter) and the bounded fix — telemetry-oracle counting with positive StartedAt/PID and explicit timeout failure, the exact two invocation identities preserved and never a third, the shell counter retained only as supporting evidence bounded by actual starts — a test-oracle correction that strengthens rather than weakens AC6/AC7. Binding before any complete: post-fix current-tree verification and green macOS CI, the separate fresh channel process before any released handoff (F2/AC9), and the attended-close conditions. Deferred follow-ups (durable-budget owner recovery, cycle-scoped goal identity) are correctly recorded as decisions, not accidents.
