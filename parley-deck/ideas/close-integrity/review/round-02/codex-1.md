---
agent: codex-1
idea: close-integrity
review-round: 2
date: 2026-06-24
---

## Summary

Round-02 refutation covered the agreed CF1-CF6 fixes from `review/consensus.md`,
the "Fix-up cycle 1" implementation note, `git show 9267034`, and the current
`internal/app/driver_impl.go` implementations of `newDriverImplOps`,
`parseGoalVerdict`, and `GoalCheck`.

The targeted app/driver regression tests passed. The full `go test -count=1 ./...`
run still fails in `internal/runner/TestDurableKillEndToEndRealProcess` with
`process verification failed (no recorded boot id); not killed`, which matches the
previously documented sandbox limitation and is not in the fix-up surface.

## Refutation attempts

- CF1 reviewer dedupe: inspected `newDriverImplOps` and ran the targeted test that
  feeds `[claude, codex, codex]` and `[claude, codex, agy, codex]`. Exact duplicate
  non-implementer IDs collapse before `ReviewerCount` and before
  `withParticipants(o.reviewers...)` reaches `RunReviewRound`. First occurrence order
  is preserved (`codex, agy`), and no legitimately distinct reviewer ID is dropped.
- CF2 wrapped verdict parsing: verified existing tests for `` `GOAL-CHECK: FAIL` ``,
  `"GOAL-CHECK: PASS"`, and `**GOAL-CHECK:** FAIL`. I also traced the over-trim case
  `GOAL-CHECK: ***`: after `TrimLeft(rest, "*\`\"'_ ")`, `rest` is empty, so it stays
  ambiguous rather than becoming PASS or FAIL.
- CF3 goal-check timeout: inspected `GoalCheck` and `runner.RunConsult`. The call now
  passes `Timeout: 2 * time.Minute`, and `RunConsult` gives a positive override
  precedence over the agent default timeout.
- CF4 last matched verdict wins: ran the parser regression covering
  `GOAL-CHECK: PASS` followed by `GOAL-CHECK: FAIL`, which still returns FAIL. The new
  reset only clears on another matched `GOAL-CHECK:` line, so a trailing ambiguous
  `GOAL-CHECK: RE-EVALUATING` returns ambiguous as intended.
- CF6 no implementer-as-checker: inspected the top of `GoalCheck` and ran
  `TestGoalCheckNoIndependentChecker`. When `drafter == implementer`, it returns
  `(true, "advisory: goal-check has no independent checker")` before discovery or
  `RunConsult`, so it fails open without launching the implementer.
- New-bug hunt around CF1/CF2/CF4: checked that CF1 does not reorder distinct reviewers,
  CF2's wrapper trim does not turn a wrapper-only line into a verdict, and CF4 does not
  regress the existing PASS-then-FAIL last-wins behavior.
- Commands run: `go test ./internal/app ./internal/driver -run 'Test(ParseGoalVerdict|NewDriverImplOpsDedupesReviewers|GoalCheckNoIndependentChecker|StrictDesignOnly|CompletesWhenGatesPassUnderAuto|GoalCheckFailEscalatesUnderAuto|SingleReviewerEscalatesUnderAuto|ReservedEscalatesUnderAuto)' -count=1`
  passed; `go test -count=1 ./...` failed only in the known runner durable-kill sandbox
  case noted above.

## Findings

No findings.

## Open questions

- The 2m goal-check timeout can turn a slow but valid checker PASS into an advisory
  fail-open. I do not consider that a defect in CF3 because the agreed design explicitly
  favors fail-open liveness for this defense-in-depth gate, but it does make goal-check
  coverage best-effort for complex acceptance criteria.
