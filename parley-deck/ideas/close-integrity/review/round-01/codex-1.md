---
agent: codex-1
idea: close-integrity
review-round: 1
date: 2026-06-24
---

## Summary

I found one MAJOR close-integrity gap. The main driver ordering is otherwise correct: `TriageReserved` and reviewer-count guards run before `GoalCheck`, `GoalCheck` runs before `Complete`, confident FAIL escalates, and checker errors/ambiguous verdicts fail open.

Acceptance-focused tests passed. `go test -count=1 ./...` did not fully pass in this environment because `internal/runner/TestDurableKillEndToEndRealProcess` failed with `process verification failed (no recorded boot id); not killed`; that looks unrelated to this diff but means I did not reproduce the claimed full-suite green state.

## Refutation attempts

- LE-11 `TriageReserved` under `auto_implement`: inspected `internal/driver/impl.go:236-239` and ran `TestReservedEscalatesUnderAuto`; it escalates before `Complete`.
- LE-11 `< 2` reviewers under `auto_implement`: inspected `internal/driver/impl.go:240-242`, `internal/app/driver_impl.go:40-44`, and `internal/app/driver_impl.go:315-318`. A normal `parley run --participants` path de-duplicates explicit IDs, but the guard itself counts slice entries, not distinct reviewer identities; see finding.
- LE-11 `>= 2` reviewers + `Ready`: ran `TestCompletesWhenGatesPassUnderAuto`; it calls `goal-check` then `complete`.
- Non-auto design-only close: ran `TestNonAutoCompletesWithSingleReviewer`; it completes without goal-check, matching conditional rigor.
- LE-7 confident FAIL: ran `TestGoalCheckFailEscalatesUnderAuto`; `GoalCheck(false, detail)` escalates before `Complete`.
- LE-7 PASS/error/ambiguous split: inspected `driverImplOps.GoalCheck`; `RunConsult` errors and missing verdicts return advisory pass, while parsed `FAIL` returns false. `TestParseGoalVerdict` passes for PASS, FAIL, no verdict, marker-prefixed lines, and last-verdict-wins.
- Checker non-implementer: `newDriverImplOps` chooses the first participant not equal to the resolved implementer as `drafter`/checker. The fallback to implementer exists only when there are no reviewers; the auto close path should hit the reviewer-count guard first.
- Parser mis-parse: tried the likely malformed-output cases against the code path: missing verdict fails open, multiple verdicts use the last one, leading markdown markers are tolerated. I did not find a direct acceptance failure there.
- Drift guard shape: both live and default `COOPERATION.md` contain the new close-decision integrity paragraph; their remaining diff is the expected project-specific header/roster data.

## Findings

### [MAJOR] ReviewerCount can count duplicate reviewer IDs as independent reviewers

What is wrong: `newDriverImplOps` appends every participant not equal to the implementer into `o.reviewers` (`internal/app/driver_impl.go:40-44`), and `ReviewStatus` reports `ReviewerCount: len(o.reviewers)` (`internal/app/driver_impl.go:315-318`). That is not necessarily the number of independent reviewers. `parseList` preserves duplicate participant IDs from existing idea frontmatter (`internal/protocol/workspace.go:327-342`), and `runner.selectedAgents` also preserves duplicate participant IDs (`internal/runner/runner.go:327-342`). A malformed or manually authored run with participants like `[impl, rev, rev]` can therefore report `ReviewerCount == 2` while only one agent identity and one canonical `rev.md` review artifact exist.

Why it matters: LE-11 is specifically guarding unattended completion against a single weak checker. Counting duplicate IDs lets an `auto_implement` close satisfy the `>= 2` reviewer guard with only one actual independent reviewer, then proceed to goal-check and `Complete`.

Concrete fix: normalize participant/reviewer IDs to a unique set before building `o.reviewers`, or compute `ReviewerCount` from distinct non-implementer reviewer IDs with validated review artifacts. Also reject duplicate participants when loading/creating run state, and add a regression test with duplicate reviewer IDs asserting the auto close escalates.

## Open questions

- Is `pipeline auto` implementation-block completion in scope for close-integrity? It has a separate completion path in `internal/app/pipeline_cmd.go` that does not appear to use these new guards.
- Is the `internal/runner` durable-kill test failure expected on this machine, or should this branch include an environment-tolerant fix before claiming `go test -count=1 ./...` is green?
