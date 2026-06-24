---
agent: hermes-1
idea: close-integrity
review-round: 1
date: 2026-06-24
---

## Summary

Verification lens: REFUTATION — assume the close branch is WRONG until I fail to break it.
I hunted for a path where the close branch reaches `Complete()` despite a guard that should
have fired (Reserved triage, <2 reviewers, goal-check FAIL), and for a path where the
goal-check fail-open swallows a confident FAIL.

The gate ordering is correct and the Tier 1 (strict_gate) + Tier 2 (loop-budget) behavior is
preserved: `strict_gate block → LE-11 → LE-7 → Complete` (impl.go:202→236→247→253). The new
close-branch code (lines 232-251) is purely additive — it only inserts escalation returns
before `Complete()` and touches no `MaxFixupCycles`/`MaxDriverSteps`/`MaxWallClock` logic.
Build, vet, and the full driver/app/runner test suites are green.

I found one MAJOR finding: `parseGoalVerdict` does not strip backtick or quote characters
from the verdict line, so a confident `GOAL-CHECK: FAIL` wrapped in backticks or quotes is
parsed as ambiguous → fail-open → `Complete()` proceeds despite the FAIL. This is the exact
"confident FAIL swallowed" path the review targets. One MINOR test-coverage gap and one NIT
round out the report.

## Refutation attempts

1. Reserved triage under AutoImplement reaches Complete().
   FAILS to break. Line 197 allows Ready|Reserved into the close branch; line 237
   (LE-11, gated on `AutoImplement`) catches Reserved and escalates. No bypass.

2. <2 reviewers under AutoImplement reaches Complete().
   FAILS to break. Line 240 (`ReviewerCount < 2`, gated on `AutoImplement`) escalates.
   `ReviewerCount` is `len(o.reviewers)` = participants minus implementer (driver_impl.go:318).
   A 2-participant deck (1 implementer + 1 reviewer) yields ReviewerCount=1 → escalates. Correct.

3. Goal-check confident FAIL reaches Complete() — the fail-open swallow path.
   BREAKS. `parseGoalVerdict` (driver_impl.go:371-388) strips leading `#*-> \t` from each line
   before checking the `GOAL-CHECK:` prefix. Backtick (`` ` ``) and double-quote (`"`) are NOT
   in the trim set. An LLM that emits `` `GOAL-CHECK: FAIL` `` or `"GOAL-CHECK: FAIL"` produces
   a line whose first char is `` ` `` or `"`, so `HasPrefix("GOAL-CHECK:")` is false → verdict
   stays `""` → `switch default` → `return true, "advisory: inconclusive"` → `Complete()`.
   I verified this by simulating the parser against backtick/single-backtick/double-backtick/
   quote/underscore-wrapped inputs; all FAIL verdicts were missed. This is a confident FAIL
   being swallowed by the fail-open path. See Findings #1.

4. Reserved triage under StrictGate-only (AutoImplement=false) reaches Complete().
   Does not break (by design). LE-11 is scoped to `AutoImplement` only (FINAL.md: "under
   AutoImplement"; "A non-auto_implement idea keeps the lighter close"). So a strict-gate
   design-only idea with Reserved + 1 reviewer + goal-check PASS does complete. This is the
   documented conditional-rigor boundary, not a bug.

5. strict_gate block ordering vs LE-11/LE-7.
   FAILS to break. The strict_gate block (lines 202-231) either escalates/opens a fresh round
   or falls through only when `certifiedClean && !findings`. LE-11/LE-7 run strictly after the
   fall-through. Goal-check thus runs only on the certified-clean round ("once before close").
   No re-ordering hazard.

6. Loop-budget (Tier 2) regression.
   FAILS to break. The new code adds no step/wall/cost logic; it only short-circuits to
   escalation. `MaxFixupCycles` bounds (lines 215, 267) and the driver.Run step budget are
   untouched. The strict-close loop still terminates at `MaxFixupCycles`.

7. GoalCheck uses the implementer as checker (drafter fallback).
   Does not break in practice. `newDriverImplOps` falls back to `drafter=implementer` when
   there are 0 reviewers (driver_impl.go:49-52). Under auto, LE-11 (`ReviewerCount < 2`) fires
   before GoalCheck → unreachable. Under strict-only with 0 reviewers, `OpenReviewRound`
   (line 214) returns "no non-implementer reviewers available" earlier → review phase never
   completes → close branch unreachable. Defended upstream, not by GoalCheck itself. See NIT #3.

8. "Last verdict wins" waffle (FAIL then PASS → PASS).
   Not a finding. `parseGoalVerdict` keeps the last verdict (test line 17 confirms PASS-then-FAIL
   → FAIL). A checker that corrects itself to PASS is reasonably treated as PASS. The prompt asks
   for a single FINAL line, so multiple verdicts are already off-spec.

## Findings

### #1 — MAJOR: parseGoalVerdict swallows backtick/quote-wrapped FAIL verdicts (fail-open bypass)

What is wrong: `parseGoalVerdict` (internal/app/driver_impl.go:375) strips only `#*-> \t`
from the start of each line before testing `strings.HasPrefix(t, "GOAL-CHECK:")`. Backtick
(`` ` ``), double-quote (`"`), and single-quote (`'`) are not in that set. An LLM emitting
`` `GOAL-CHECK: FAIL — criteria unmet` `` (a very common way to render a code-like structured
token, even when the prompt says "EXACTLY") produces a line whose first character is a
backtick. `HasPrefix("GOAL-CHECK:")` is false, the line is skipped, `verdict` stays `""`,
and `GoalCheck` returns `(true, "advisory: goal-check inconclusive")` — fail-open. The driver
then calls `Complete()`. A confident FAIL is silently swallowed and the idea auto-completes.

Why it matters: this is the precise refutation target — "the goal-check fail-open accidentally
swallows a confident FAIL." The LE-7 gate is defense-in-depth, not the sole gate, so the blast
radius is bounded (review consensus already passed), but the gate is fully defeated for the most
common LLM formatting variation of a structured verdict line. The prompt's "EXACTLY" instruction
reduces but does not eliminate the risk; backtick-wrapping of machine-readable tokens is
pervasive enough across models that relying on prompt compliance alone is fragile.

Concrete fix: add `` ` ``, `"`, and `'` to the `TrimLeft` cut set, and also strip a trailing
backtick/quote run from `rest` before the PASS/FAIL prefix check (so `` `GOAL-CHECK: FAIL` ``
and `"GOAL-CHECK: FAIL"` both parse). Minimal change at driver_impl.go:375:

```go
t = strings.TrimLeft(t, "#*-> \t`\"'")
```

and after computing `rest` (line 379), trim closing wrappers:

```go
rest = strings.TrimRight(rest, "`\"'")
```

Add test cases to `goal_check_test.go`: `` `GOAL-CHECK: FAIL` `` → "FAIL",
`"GOAL-CHECK: PASS"` → "PASS". This keeps the fail-open semantics intact for genuinely
ambiguous output while ensuring wrapped verdicts are not misclassified as inconclusive.

### #2 — MINOR: no test for the StrictGate-only (AutoImplement=false) close path

What is wrong: `newStrictDriver` (strict_gate_test.go:14-19) sets both `AutoImplement: true`
and `StrictGate: true`. There is no test exercising `StrictGate: true, AutoImplement: false`
— the design-only strict idea. On that path, LE-11 (Reserved / <2 reviewers) is skipped (gated
on `AutoImplement`) while LE-7 goal-check runs (gated on `AutoImplement || StrictGate`). The
behavior is intended per FINAL.md, but it is untested: a regression that accidentally gates
LE-7 on `AutoImplement` only, or that makes LE-11 fire for design-only strict ideas, would not
be caught.

Concrete fix: add a test in `close_integrity_test.go` (or `strict_gate_test.go`) with
`StrictGate: true, AutoImplement: false`, Reserved triage, 1 reviewer, goal-check PASS →
assert `ActionComplete` (documents the conditional-rigor boundary), and a second with
goal-check FAIL → assert `ActionEscalated` (LE-7 still fires without auto).

### #3 — NIT: GoalCheck comment claims "fresh non-implementer" but o.drafter can be the implementer

What is wrong: the `GoalCheck` doc comment (driver_impl.go:332-335) and FINAL.md state the
checker is "a fresh non-implementer agent (the review drafter)." But `newDriverImplOps`
(driver_impl.go:49-52) falls back to `drafter = implementer` when there are no reviewers.
`GoalCheck` itself does not guard against this — it uses `o.drafter` unconditionally.

Why it is only a NIT: the fallback is unreachable at the close branch in practice. Under
`AutoImplement`, LE-11 (`ReviewerCount < 2`) fires before GoalCheck. Under strict-only with
zero reviewers, `OpenReviewRound` fails earlier so the review phase never completes. The
invariant is defended by upstream guards, not by `GoalCheck`. But the code as written does not
enforce its own stated contract locally, which makes it fragile to future reordering of the
close-branch guards.

Concrete fix: either (a) add a defensive check at the top of `GoalCheck` — if `o.drafter ==
o.implementer`, warn and fail-open (advisory) rather than running the implementer as checker;
or (b) amend the doc comment to note the upstream-guard dependency. Option (a) is preferred
for local robustness.

## Open questions

1. Is the prompt's "EXACTLY" instruction considered a strong enough contract that backtick
   wrapping is out of scope, or should the parser be hardening-tolerant regardless? I lean
   strongly toward the latter (Finding #1) given this is a defense-in-depth gate whose whole
   purpose is to catch what the review consensus missed.

2. Should a strict-gate design-only idea (AutoImplement=false, StrictGate=true) with an
   ACCEPT-WITH-RESERVATIONS triage and a single reviewer auto-complete? FINAL.md says yes
   (LE-11 is auto-only), but it is worth confirming this is an intentional acceptance of a
   weaker close for strict design ideas, not an oversight in the gating predicate.

3. The goal-check verdict is an LLM judgment, not a proof. FINAL.md acknowledges this and
   fail-opens on checker failure. Given Finding #1, is the fail-open surface area (any parse
   miss → proceed) too broad, or is the review-consensus precondition a sufficient backstop?
   I assess the backstop is sufficient for CRITICAL prevention but not for the gate's own
   integrity claim.
