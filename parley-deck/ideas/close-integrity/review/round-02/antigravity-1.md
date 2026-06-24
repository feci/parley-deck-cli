---
agent: antigravity-1
idea: close-integrity
review-round: 2
date: 2026-06-24
---

## Summary

This is the round-02 refutation-mode review of the `close-integrity` idea (Phase 6, Round 2) by agent **antigravity-1**. All round-01 agreed fixes (CF1-CF6) implemented in commit `9267034` have been verified. Refutation testing failed to break the updated implementation, and no new regressions or bugs were discovered.

## Refutation attempts

I evaluated the fixes in `internal/app/driver_impl.go` and `internal/driver/impl.go` by assuming they were incorrect or incomplete, and attempting to break them:

1. **CF1 (Reviewer Deduplication):** 
   - *Attempt:* Test if duplicate/reordered/missing reviewer lists can bypass the `ReviewerCount < 2` guard under auto-drive, or cause concurrent write races in review rounds.
   - *Result:* Fails to break. The deduplication loop in `newDriverImplOps` iterates over the input slice, excludes the resolved implementer, and collapses all duplicates into a unique ordered subset using a `seen` map. This guarantees that `ReviewerCount` represents the count of distinct independent reviewers, and `OpenReviewRound` fans out exactly one goroutine per unique reviewer.

2. **CF2 (Wrapped Verdict Parsing):**
   - *Attempt:* Supply various wrapped strings to `parseGoalVerdict` (e.g. bolded prefix `**GOAL-CHECK:** FAIL`, backticked `` `GOAL-CHECK: FAIL` ``, quoted `"GOAL-CHECK: PASS"`, bulleted `- GOAL-CHECK: PASS`, and other weird wrappers) to see if they are missed or trigger a fail-open.
   - *Result:* Fails to break. The line `TrimLeft` in `parseGoalVerdict` strips `#*->_ \t`\"'` before looking for `GOAL-CHECK:`. Then, the `rest` `TrimLeft` strips `*`\"'_ ` (the closing markdown/quote characters that now lead the rest of the string). As a result:
     - `` `GOAL-CHECK: FAIL` ``: `t` trims to `GOAL-CHECK: FAIL``, prefix matches, `rest` trims left to `FAIL``, `strings.HasPrefix` matches `"FAIL"`. Correct.
     - `**GOAL-CHECK:** FAIL`: `t` trims to `GOAL-CHECK:** FAIL`, prefix matches, `rest` starts as `"** FAIL"`, trims left to `"FAIL"`, `strings.HasPrefix` matches `"FAIL"`. Correct.
     - `"GOAL-CHECK: PASS"`: `t` trims to `GOAL-CHECK: PASS"`, prefix matches, `rest` trims to `"PASS"`, `strings.HasPrefix` matches `"PASS"`. Correct.
     - A line containing only `GOAL-CHECK: ***` is correctly parsed as empty (`""` / ambiguous) and does not incorrectly match `PASS` or `FAIL`.

3. **CF3 (Goal-Check Timeout):**
   - *Attempt:* Check if `GoalCheck` can still hang and block the driver's execution loop under slow network/API conditions.
   - *Result:* Fails to break. `runner.RunConsult` is now invoked with an explicit `Timeout: 2 * time.Minute`. If the agent or network hangs, the context is cancelled after 2 minutes, and the driver immediately fails open (advisory) instead of blocking for the default 15–30 minutes.

4. **CF4 (Last Verdict Wins & Reset):**
   - *Attempt:* Supply inputs containing multiple matched verdict lines, such as `PASS` then `FAIL`, or `FAIL` then an ambiguous statement like `GOAL-CHECK: RE-EVALUATING`.
   - *Result:* Fails to break. Because `verdict = ""` is executed immediately when any line matches the `GOAL-CHECK:` prefix, the verdict is reset for every matched line. A trailing ambiguous line (e.g., `GOAL-CHECK: RE-EVALUATING`) correctly resets the verdict to `""` rather than leaving the prior `FAIL` or `PASS` intact. Last-wins (e.g. `PASS` then `FAIL` -> `FAIL`) continues to work correctly.

5. **CF5 (Strict Gate Design-Only Path):**
   - *Attempt:* Trace if a strict design-only idea (StrictGate=true, AutoImplement=false) correctly executes the close-path tests.
   - *Result:* Fails to break. The new tests `TestStrictDesignOnlyCompletesAndRunsGoalCheck` and `TestStrictDesignOnlyGoalCheckFailEscalates` successfully verify that under StrictGate, the LE-7 goal-check runs regardless of `AutoImplement`, but the LE-11 reviewer-count checks do not apply (which is the expected conditional-rigor boundary).

6. **CF6 (Goal-Check Implementer Fallback):**
   - *Attempt:* Trigger the fallback where the implementer checks their own code by setting `drafter == implementer`.
   - *Result:* Fails to break. The top of `GoalCheck` checks `if checker == o.implementer` and immediately logs a warning and fails open (`true, "advisory: goal-check has no independent checker"`) instead of executing the implementer agent.

## Findings

No findings of severity **CRITICAL**, **MAJOR**, **MINOR**, or **NIT** were identified in this round. The fix-up completely resolves all agreed points from round-01, is robust against refutation, and introduces no regressions or new bugs.

## Open questions

1. **Timeout Fail-Open Behavior:** A 2-minute timeout is appropriate to prevent a driver tick hang. However, in heavily throttled or high-latency environments, it might flip a slow-but-valid `PASS` or `FAIL` into an advisory fail-open (PASS). This is accepted per the advisory fail-open design, but should be noted as a known trade-off.
