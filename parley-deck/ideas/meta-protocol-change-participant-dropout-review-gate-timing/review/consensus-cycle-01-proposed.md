---
idea: meta-protocol-change-participant-dropout-review-gate-timing
review-cycle: 1
outstanding_agreed_fixes: 3
blocked: false
drafted-by: zcode-1
date: 2026-10-09
reviewed-commit: a26f588
---

## Agreed fixes

Drafter's note: this cycle-1 consensus codifies the fix plan proposed by codex-1 (the
implementer) in response to `review/round-01/zcode-1.md`. I concur with all three items
below; each matches or extends the suggested fix recorded in the originating finding.
The single independent reviewer this round is also the drafter; my positions here are
unchanged from my round-01 review file (no position changes).

- From `zcode-1/review/round-01` **[MINOR] "Standalone mid-idea signoff now hard-depends
  on the prior run's manifest"** — in `internal/app/quota_signoff.go`, treat only
  `os.IsNotExist` on the prior run's manifest as empty snapshot inheritance when
  creating a standalone mid-idea signoff run, surfaced with a clear diagnostic; keep the
  hard failure for corrupt or foreign-identity manifests. An empty inheritance must never
  invent models and must never qualify the single-reviewer exception (an absent snapshot
  cannot qualify, per §9.0/AC4 — the exception simply fails to derive, which is the
  fail-safe direction). Tests: missing prior manifest (empty inheritance + diagnostic),
  missing snapshot, corrupt/foreign rejection, plus a rerun of the existing same-tick
  rebind test (`TestReviewGateSignoffWatchdogAndSameTickRebind`), which is the covering
  test for this seam (resolves my round-01 Open question 3).
- From `zcode-1/review/round-01` **[NIT] "Missing blank line before new heading in
  docs/agent-cli-mechanics.md"** — insert the Markdown blank line before the
  `## Participant-step timing and buffered output (1.53.0)` heading.
- From `zcode-1/review/round-01` **[NIT] "Legacy-rule negatives not exercised at the
  review-gate seam"** — add adversarial table rows at the `reviewFailureEvidence` seam
  for: a short reset on `quota.named-reset-ge-60m.v1`, a contradictory reset on a
  no-reset rule, a non-zcode adapter, and wrong provenance. Preserve the positive
  boundary cases and all three enumerated legacy rules. (This adds wrong-provenance to
  the three rows my finding suggested — accepted.)

## Deferred follow-ups

- From `zcode-1/review/round-01` **[NIT] "full-host-tests-v2 receipt needed recovery;
  on-disk *-result.json unreadable"** — defer the investigation of the shared-volume
  ENOENT behavior to a follow-up idea (slug `TBD`; open it if the volume repeats this).
  Disposition rationale, which I concur with openly: my own PRIMARY recount re-hashed
  the complete actual exit-0 jsonl (sha256 match, 3,315 distinct test/subtest IDs
  recounted: 3,311 pass + 4 skip, 34/34 packages, 0 fail events), so the pass is
  independently established from the full raw log; the failure mode is disclosed
  evidence-storage trouble on the shared volume, not a product defect in this diff and
  not a test exemption. All original and failed logs, the exact recovered emitted
  receipt, and the hash/recount proof stay retained
  (`.parley-runtime/review-gate-timing/full-host-tests-v2*`, committed
  `source-context/implementation-validation.json`).
- From `zcode-1/review/round-01` **Open question 2** — the deferred timeout-tuning
  follow-up idea should record the fresh goal-check durations this run produces (per
  FINAL.md "Known risks"), alongside the already-disclosed fact that a genuinely hung
  buffered zcode signoff is detected only at its hard ceiling.

## Dismissed findings

None. No finding was withdrawn and none was judged not-an-issue: three findings became
Agreed fixes and one became a Deferred follow-up. No verdict conflicts exist this cycle
(single reviewer; §15.3 requires no `## Verdict conflicts` section).

## Coverage & blind spots

- **Reviewer coverage this round:** one independent reviewer (zcode-1). kimi-1 was
  owner-excluded before Phase 0 (recorded in `00-prompt.md`), so this floor run had
  exactly one non-implementer participant; the run correctly claims no product
  single-reviewer exception for itself (no mid-idea automatic transition exists to
  derive one from — verified in my review, AC13). Cross-reviewer overlap is therefore
  undefined this round; coverage rests on the round-1 refutation sweep across all 14
  ACs plus independent reruns of the focused Go suites, the packet/drift guard tests,
  `go build ./...` and `go vet ./...`.
- **Blind spots disclosed in the review:** no attempt was made to hand-forge a history
  file (tampering is outside this rule's threat model; stops-for-repair is the shipped
  stance); the skill-suite pass at 74cc831 is the dependency-refreshed v3 run after
  `npm ci` (v1/v2 failures retained as history); protocol/skill text equality was
  established by independent byte-exact reconstruction, not by eyeball.
- **Process state after this consensus:** this is a fix cycle, not a close. Still
  required before completion: the fresh re-review of the fix-up diff, a later Phase 7
  consensus with zero Agreed fixes, both current participants' final review-consensus
  ACCEPT blocks, the fresh zcode-1 goal-check process, and the current-tree AC evidence
  (owner-attended close authority per `00-prompt.md` Constraints).

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: zcode-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: ACCEPT approves the cycle-1 fix plan only, not product close. I concur with all three Agreed fixes (they match or extend my own suggested fixes in review/round-01) and with deferring the shared-volume ENOENT receipt investigation: my PRIMARY recount re-hashed the complete exit-0 jsonl, so the headline suite pass is independently established and the failure is disclosed storage trouble, not a test exemption. Awaiting fresh re-review of the fix-up diff, a zero-fix consensus, both final ACCEPT blocks, the fresh goal-check process and current-tree AC evidence.
