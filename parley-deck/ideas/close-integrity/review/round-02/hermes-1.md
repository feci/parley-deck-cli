---
agent: hermes-1
idea: close-integrity
review-round: 2
date: 2026-06-24
---

## Summary

Verification lens: REFUTATION — assume the round-01 fixes (CF1–CF6) are WRONG or
incomplete until I fail to break them. I read `review/consensus.md`, the "Fix-up
cycle 1" section of `IMPLEMENTATION.md`, `git show 9267034`, and the CURRENT
`internal/app/driver_impl.go` (newDriverImplOps, parseGoalVerdict, GoalCheck). I
also traced `RunConsult`/`timeoutForAgent` in `internal/runner/consult.go` to
confirm CF3 actually takes effect, and ran the driver/app test suites
(`ok internal/app`, `ok internal/driver`).

Verdict: all six agreed fixes hold under refutation. CF1 collapses duplicate
non-implementer IDs to a distinct set preserving first-seen order; CF2 parses
all three named wrapper variants; CF3's 2m timeout is honored (not silently
overridden by the agent default); CF4 resets per matched line; CF6 fails open
when drafter==implementer. I found no CRITICAL, MAJOR, or MINOR regressions
introduced by the fixes. I report one NIT (pre-existing, surfaced sharper by
CF2's rest-trim) and document the acceptable fail-open behavior of the 2m
timeout. I explicitly state below what I tried and why it failed to break each
fix, plus the new-bug hunts that came up clean.

## Refutation attempts

### CF1 — duplicate non-implementer IDs collapse to a distinct set

I tried to break this three ways against the current `newDriverImplOps`
(driver_impl.go:38-64):

1. `[impl, rev, rev]` inflates ReviewerCount past the LE-11 `< 2` guard.
   FAILS to break. The `seen` map (line 46-49) skips the second `rev`, so
   `o.reviewers == [rev]`, `ReviewerCount == 1` (line 325 `len(o.reviewers)`),
   and the `< 2` guard at impl.go:240 escalates. The count-bypass is closed.

2. `[impl, rev, rev]` launches two goroutines writing the same
   `agents/rev/stdout.log` and `review/round-NN/rev.md` concurrently.
   FAILS to break. `o.reviewers` holds one `rev`, so `OpenReviewRound`
   (driver_impl.go:242) passes a single-element list to `RunReviewRound`; one
   goroutine, one writer. The concurrent-write race is closed. I confirmed
   `RunReviewRound` fans out per list element (the consensus cites runner.go:204)
   so deduping the input list is the correct fix locus.

3. CF1 drops a legitimately distinct reviewer or reorders them wrongly.
   FAILS to break. The loop appends in participant order, gated only by
   `p != implementer && !seen[p]`; `seen[p] = true` is set on first append. So
   `[claude, codex, agy, codex]` → `[codex, agy]` (codex first-seen at index 1,
   agy at index 2, second codex skipped). First-seen order is preserved; no
   distinct reviewer is dropped. The regression test
   `TestNewDriverImplOpsDedupesReviewers` (driver_impl_le_test.go) asserts both
   the `[rev]` collapse and the `[codex, agy]` order; it passed in my run.

   One subtlety I checked: `drafter` is set to `reviewers[0]` (line 58). With
   `[impl, rev, rev]` the drafter is `rev` (correct, a non-implementer). With
   `[impl, revA, revB]` the drafter is `revA`. Dedup does not change which agent
   becomes drafter for the realistic non-duplicate case. No reordering hazard.

### CF2 — wrapped verdicts parse instead of fail-opening

I simulated `parseGoalVerdict` (driver_impl.go:391-418) in Python against the
three named variants plus ~25 adversarial inputs. Results for the named cases:

- `` `GOAL-CHECK: FAIL` `` → "FAIL"  (backtick-wrapped, the line-level TrimLeft
  eats the leading `` ` ``, prefix matches, rest = "FAIL`", rest TrimLeft eats
  nothing (`` ` `` is in the cut set but `FAIL`… starts with F), HasPrefix FAIL
  → FAIL). Correct.
- `"GOAL-CHECK: PASS"` → "PASS"  (line-level TrimLeft eats `"`, prefix matches,
  rest = `PASS"`, HasPrefix PASS → PASS). Correct.
- `**GOAL-CHECK:** FAIL` → "FAIL"  (line-level TrimLeft eats leading `**`, t =
  `GOAL-CHECK:** FAIL`, prefix matches, rest after TrimPrefix = `** FAIL`, rest
  TrimLeft eats `** ` → `FAIL`, HasPrefix FAIL → FAIL). Correct — this is the
  antigravity case and it now parses.

So CF2 holds for all three named wrappers. The fail-open swallow path for a
confident wrapped FAIL is closed.

### CF3 — the goal-check RunConsult passes Timeout: 2*time.Minute

I traced the timeout plumbing, not just the call site. `GoalCheck`
(driver_impl.go:361-373) sets `Timeout: 2 * time.Minute` in `ConsultOptions`.
`RunConsult` (consult.go:101) calls `timeoutForAgent(opts.Timeout, opts.Agent)`.
`timeoutForAgent` (runner.go:1105-1110) is: `if override > 0 { return override }`
— so a positive 2m override is returned directly and the agent's 15–30m default
(TimeoutMS / DefaultTimeoutMS) is NOT consulted. `context.WithTimeout(ctx,
hardTimeout)` (consult.go:103) then enforces 2m. CF3 genuinely bounds the
consult; it is not a no-op. A hung checker now fails open at 2m instead of
blocking the driver tick for the full agent timeout.

### CF4 — verdict resets per matched line so a trailing ambiguous line returns ambiguous

I tried to break this two ways:

1. A trailing ambiguous line leaves the prior verdict stuck.
   FAILS to break. `verdict = ""` runs on every matched `GOAL-CHECK:` line
   (driver_impl.go:405) BEFORE the switch. So `GOAL-CHECK: FAIL` then
   `GOAL-CHECK: RE-EVALUATING`: line 2 matches the prefix, resets to "", the
   switch finds no PASS/FAIL prefix on "RE-EVALUATING", verdict stays "". Returns
   "". The test case `GOAL-CHECK: FAIL\nGOAL-CHECK: RE-EVALUATING` → "" is in
   goal_check_test.go and passed.

2. CF4's reset breaks the existing PASS-then-FAIL last-wins.
   FAILS to break. The reset to "" is immediately followed by the switch, which
   reassigns verdict for the current line. `GOAL-CHECK: PASS\nGOAL-CHECK: FAIL`:
   line 1 resets to "" then sets "PASS"; line 2 resets to "" then sets "FAIL".
   Returns "FAIL". Last-wins is preserved — the reset only matters when the
   current line is itself ambiguous, which is exactly the intended behavior. The
   existing test case `GOAL-CHECK: PASS\nGOAL-CHECK: FAIL` → "FAIL" still passes.

### CF6 — GoalCheck fails open when drafter==implementer

I tried to reach GoalCheck with `o.drafter == o.implementer` and have it run the
implementer as its own checker. FAILS to break. The guard at driver_impl.go:349
returns `(true, "advisory: goal-check has no independent checker")` before any
`discoveryFor`/`RunConsult` call. No agent runs. The test
`TestGoalCheckNoIndependentChecker` (driver_impl_le_test.go) sets `drafter =
"claude"` == implementer and asserts `ok == true`; it passed. The contract is
now enforced locally, not just by upstream guards.

### New-bug hunts (fixes introducing fresh defects)

I hunted specifically for the four classes the round-02 brief named.

1. Does CF2's rest TrimLeft over-trim and misread a legit verdict
   (e.g. a line that is only "GOAL-CHECK: ***")?
   No misread found. `GOAL-CHECK: ***` → rest after TrimPrefix = "***", rest
   TrimLeft eats all `*` → "", no PASS/FAIL prefix → verdict "". That is the
   correct ambiguous/fail-open outcome for a line with no actual verdict word —
   it is not a legit verdict being misread. `GOAL-CHECK: ****` → "" likewise.
   The dangerous case would be a TrimLeft that eats part of a real verdict: I
   checked `GOAL-CHECK: *PASS` → "PASS" (eats `*`, "PASS" prefix matches),
   `GOAL-CHECK: *FAIL*` → "FAIL" (eats leading `*`, "FAIL" matches; trailing `*`
   is irrelevant to HasPrefix), `GOAL-CHECK: *** FAIL` → "FAIL" (eats `*** `,
   "FAIL" matches). The cut set `*`"`"'_ ` contains no letter, so it can never
   consume the P of PASS or the F of FAIL. The rest-trim is safe.

2. Does CF4's reset break the existing PASS-then-FAIL last-wins?
   No — covered under CF4 attempt 2 above. The reset is immediately overwritten
   by the switch on the same line; it only changes the outcome when the current
   line is ambiguous, which is the intended "true last-verdict-wins" semantics.

3. Does CF1 ever drop a legitimately distinct reviewer or reorder them wrongly?
   No — covered under CF1 attempt 3. First-seen order preserved; `seen` only
   blocks exact duplicates; distinct IDs are never dropped.

4. Is the 2m timeout too short and would it flip a slow-but-valid PASS into a
   fail-open?
   This is the one behavior worth noting explicitly. A goal-check consult that
   takes > 2m (a large repo, a slow model, a cold-start) will hit
   `context.WithTimeout` at 2m → `res.ExitError != ""` → GoalCheck returns
   `(true, "advisory: goal-check checker error")` (driver_impl.go:374-376) → the
   driver proceeds to Complete. This is a fail-OPEN, not a fail-closed: a
   slow-but-valid PASS becomes an advisory pass, and the idea auto-completes
   anyway. Per the fail-open design (LE-7 is defense-in-depth on top of an
   already-passed review consensus), this is acceptable — the gate's purpose is
   to catch a confident FAIL, not to enforce a deadline on a PASS. I note it as
   a design observation, not a defect: a 2m budget is tight for a genuinely
   large workspace but the failure mode is the safe one (proceed, not block).
   If the team later wants goal-check to be a hard gate rather than advisory,
   the 2m ceiling would need revisiting; today it is fine.

## Findings

### [NIT] parseGoalVerdict matches PASSED/FAILURE as PASS/FAIL (pre-existing, sharper now)

What is wrong: the PASS/FAIL check is `strings.HasPrefix(rest, "PASS")` /
`HasPrefix(rest, "FAIL")` (driver_impl.go:411-414). So `GOAL-CHECK: PASSED` →
"PASS" and `GOAL-CHECK: FAILURE` → "FAIL". This is pre-existing behavior (the
round-01 parser did the same prefix check); CF2's rest TrimLeft does not change
it. I flag it only because CF2 widened the input shapes that reach the switch,
so a wrapped `GOAL-CHECK: **FAILURE**` now also resolves to "FAIL" where before
it would have fail-opened — marginally more exposure to the prefix-only match.

Why it is only a NIT: an LLM emitting "PASSED" or "FAILURE" instead of "PASS"/
"FAIL" is off-spec (the prompt asks for EXACTLY `GOAL-CHECK: PASS` or
`GOAL-CHECK: FAIL — <reason>`), and the direction of error is asymmetric and
mostly safe: "FAILURE" → FAIL escalates (conservative — surfaces a human rather
than auto-completing), and "PASSED" → PASS is the only direction that could let
a borderline-true verdict through, which requires the checker to have written a
non-spec word that nonetheless prefix-matches the pass token. The review
consensus already passed, so the blast radius is bounded by LE-7's defense-in-
depth role. Not worth a code change for this idea; noting for the record.

Concrete fix (if desired in a future hardening idea): match on the exact token
boundary — `rest == "PASS" || strings.HasPrefix(rest, "PASS ") || ...` — or
require the verdict word to be followed by whitespace, punctuation, or end-of-
string. Out of scope for close-integrity's fix-up cycle.

No CRITICAL, MAJOR, or MINOR findings. The six agreed fixes are correctly
implemented and survive refutation.

## Open questions

1. The 2m goal-check timeout (CF3) is an advisory-gate budget, and its timeout
   fails open (proceed). If a future idea promotes LE-7 from advisory to a hard
   close gate, this ceiling becomes the effective deadline for a valid PASS and
   may need to scale with workspace size. Worth recording as a follow-up note on
   the goal-check design, not a change to close-integrity.

2. CF1 dedupes `o.reviewers` at the driver-adapter boundary, which closes the
   close-integrity-relevant race. The deferred DF1 (reject duplicate participant
   IDs at the `parseList`/`workspace.go` load boundary) remains the right place
   for global defense-in-depth across all phases. I concur with the consensus
   deferral; no action here.

3. The prefix-only PASS/FAIL match (NIT above) is a pre-existing parser
   coarseness that CF2 did not introduce. Should a future hardening idea tighten
   the verdict token to an exact-word match? Low priority given the fail-open
   design and review-consensus backstop, but it is the natural next parser
   refinement.
