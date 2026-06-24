---
idea: close-integrity
status: complete
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
branch: parley-deck-cli#loop-engineering-impl
head-commit: (this commit)
review-round: 2
fixup-cycle: 1
---

## Summary of work

Tier 3 (close-decision integrity): LE-11 HITL-fatigue guardrails + LE-7 goal-done gate,
per `FINAL.md`. `go build`, `go vet`, `go test -count=1 ./...`, and the drift guard are
green.

## Implementation plan / checklist

- [x] **LE-11 guardrails** — `driver.ReviewStatus` gains `ReviewerCount int` (set by
  `driverImplOps.ReviewStatus` = `len(o.reviewers)`). In `advanceReview`'s
  `outstanding_agreed_fixes == 0` close branch, **under `AutoImplement`**: a
  `TriageReserved` triage escalates (reservations need a human), and `ReviewerCount < 2`
  escalates (a single checker is too weak for unattended completion).
  (`internal/driver/impl.go`, `internal/app/driver_impl.go`.)
- [x] **LE-7 goal-done gate** — new `ImplOps.GoalCheck(ctx) (bool, string)`. The app
  adapter runs the **review drafter** (a non-implementer) via the existing
  `runner.RunConsult` path with a new `BuildGoalCheckPrompt` (verdict line
  `GOAL-CHECK: PASS|FAIL`). `parseGoalVerdict` reads the last verdict (case-insensitive,
  marker-tolerant). **Fail-open:** a checker error / ambiguous verdict returns
  `(true, advisory)` and warns — only a confident FAIL returns `(false, …)`. In
  `advanceReview`, under `AutoImplement || StrictGate`, after the LE-11 guards and before
  `Complete()`, a `!ok` escalates. (`internal/driver/{impl}.go`,
  `internal/app/driver_impl.go`, `internal/runner/consult.go`.)
- [x] **Protocol** — §4 (after the loop-budget invariant) gains the close-decision
  integrity paragraph (LE-7/LE-11). Both `COOPERATION.md` copies (drift green).
- [x] **Tests** — `internal/driver/close_integrity_test.go` (Reserved→escalate,
  `<2`→escalate, goal-check fail→escalate, gates-pass→complete, non-auto single-reviewer
  completes without goal-check); `internal/app/goal_check_test.go` (`parseGoalVerdict`).
  Updated the existing complete-path fakes with `ReviewerCount: 2`.

## Deviations from FINAL.md

- **LE-11a (batch/rate-limit driver-opened HITL questions) — deferred** (documented in
  FINAL non-goals): per-idea fatigue is already bounded (the driver halts on escalation,
  it does not flood); cross-idea question batching belongs to the Tier-4 outer loop.
- **Goal-check is gated on `AutoImplement || StrictGate`** and fail-open on its own error:
  it is a defense-in-depth signal on top of an already-passed review consensus, not the
  sole gate — a broken/timed-out/ambiguous checker warns and proceeds rather than blocking
  a review-clean idea. (hermes's "Adapt, not full Adopt".)

## Notes for reviewers

- **Behavior change:** an `auto_implement` idea with a single reviewer (e.g. a 2-agent
  1-implementer-1-reviewer deck) no longer auto-completes — it escalates for human signoff
  or a second reviewer. This is intended (LE-11). Verify a design-only idea is unaffected.
- **Fresh checker ≠ implementer:** GoalCheck uses `o.drafter` (a non-implementer). Confirm
  it never uses the implementer, and that the prompt is read-only.
- Try to break: a goal-check answer with no verdict line (must fail-open, not block); a
  Reserved triage under `auto_implement` (must escalate); exactly 2 reviewers + Ready +
  PASS (must complete).

## Fix-up cycle 1 (round-01 review consensus)

All six agreed fixes from `review/consensus.md` applied. `gofmt`, `go build ./...`,
`go vet`, `go test -count=1 ./...` (incl. drift guard) green.

- **CF1 (CRITICAL/MAJOR)** — `newDriverImplOps` now dedupes to distinct non-implementer
  IDs via a `seen` map. Closes both the `ReviewerCount` count-bypass of the LE-11 `< 2`
  guard (codex) and the duplicate-goroutine concurrent-write race on the same log/artifact
  files (antigravity). (`internal/app/driver_impl.go`.)
- **CF2 (MAJOR)** — `parseGoalVerdict` strips leading markdown/quote wrappers
  (`` ` "' _ * # > ``) from the line AND a leading wrapper run from `rest`, so
  `` `GOAL-CHECK: FAIL` ``, `"GOAL-CHECK: PASS"`, and `**GOAL-CHECK:** FAIL` all parse
  instead of fail-opening. (`internal/app/driver_impl.go`.)
- **CF3 (MAJOR)** — the goal-check `RunConsult` now passes `Timeout: 2 * time.Minute` so
  the advisory gate fails open fast rather than blocking the driver for the agent's full
  15–30 min timeout on a hang. (`internal/app/driver_impl.go`.)
- **CF4 (MINOR)** — `parseGoalVerdict` resets `verdict` on each matched verdict line so a
  trailing ambiguous line correctly returns "" (true last-verdict-wins).
- **CF5 (MINOR)** — added strict-design-only (`StrictGate: true, AutoImplement: false`)
  close-path tests + `newStrictDesignDriver` helper. (`internal/driver/close_integrity_test.go`,
  `internal/driver/strict_gate_test.go`.)
- **CF6 (NIT)** — `GoalCheck` fails open without running an agent when `o.drafter ==
  o.implementer` (never run the implementer as its own checker). (`internal/app/driver_impl.go`.)

New tests: `goal_check_test.go` (wrapper + reset cases), `driver_impl_le_test.go`
(`TestNewDriverImplOpsDedupesReviewers`, `TestGoalCheckNoIndependentChecker`),
`close_integrity_test.go` (`TestStrictDesignOnly*`).

**Deferred (own follow-up ideas, documented in `review/consensus.md`):** DF1 reject
duplicate participant IDs at the load boundary; DF2 extend LE-7/LE-11 to the `pipeline
auto` block-completion path (a separate §12 subsystem); DF3 exact-token verdict match in
`parseGoalVerdict` (pre-existing `PASSED`/`FAILURE` prefix coarseness — a future parser
hardening idea). **Resolved:** the codex durable-kill test failure is a codex-sandbox
limitation (passes in the implementer's env).

## Round-02 re-review — complete

All three reviewers re-reviewed the fix-up in refutation mode and returned clean:
codex-1 "No findings", antigravity-1 "No findings / no regressions", hermes-1 "No
CRITICAL/MAJOR/MINOR" (one pre-existing NIT → DF3, dismissed for this idea). Zero new
agreed fixes → close-integrity marked complete. `gofmt`, `go build ./...`, `go vet`,
`go test -count=1 ./...` (incl. drift guard) green.