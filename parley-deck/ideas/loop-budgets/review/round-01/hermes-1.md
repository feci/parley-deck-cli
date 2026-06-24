---
agent: hermes-1
idea: loop-budgets
review-round: 1
date: 2026-06-24
---

## Summary

Refutation review (Phase 6) of loop-budgets (LE-5/6/10) at commit range
c109ffe..HEAD. I assumed the implementation was WRONG and tried to break the
escalate-not-complete invariant, the breach-detection logic, and the LE-6 cost
honesty claims. I traced the full Run() control flow, the config precedence
chain, the store typing, and the pipeline watcher fix; I ran the test suite and
vet/build.

What I tried and could NOT break:
- No fail-open path to Complete() after a budget breach. The breach check
  (loop.go:46) sits BEFORE Advance() at the top of every iteration, and on a
  breach it returns escalateLoopBudget immediately — the ActionComplete branch
  (loop.go:99) is below and only reachable when Advance returns Complete on a
  tick that passed the breach check. Because Complete is not a progress action,
  it never increments `steps`, so a clean completion on an under-budget idea is
  not falsely blocked. Escalate-not-complete holds.
- loopBudgetBreach returns the correct result for all three ceilings. The
  `steps >= MaxDriverSteps` bound (not `>`) means a ceiling of N allows exactly N
  progress steps before escalation on step N+1's tick — correct. 0 ceilings are
  unlimited (the `> 0` guards short-circuit). Wall-clock uses a single `start`
  set once at Run entry (loop.go:34), never reset — matching FINAL.md's "total
  run budget, distinct from per-tick roundDeadline".
- LE-6 cost honesty: loopCostUSD (loop.go:177) uses a type-assertion
  `e.Data["cost_usd"].(float64)` with the comma-ok form, so a non-float cost_usd
  (string, int, nil, missing key) is skipped, not a panic. The "telemetry-gated
  / inert until agent.usage emitted" limit is accurately documented in FINAL.md
  §LE-6, IMPLEMENTATION.md, the Config field comment, the seed template comment,
  and COOPERATION.md §4. No runner emits agent.usage today (grep confirms), so
  the ceiling is genuinely 0-bound in practice.
- The MaxDriverSteps=1 acceptance edge: a SINGLE-progress-step idea completes
  (iter1: breach check steps=0 passes, Advance→Complete→return nil). This is
  correct — an idea that finishes within budget should not escalate. A
  MULTI-step idea (≥2 progress actions) escalates on the second tick (steps=1,
  1>=1). The FINAL.md wording "a multi-step idea escalates" is precise and
  holds.

Build/vet/tests green for the touched packages (driver, config, app).

## Refutation attempts

1. Fail-open after breach: I looked for any branch between the breach check
   (loop.go:46-48) and the ActionComplete return (loop.go:99-101) that could
   bypass escalation. There is none — a breach returns directly. The only way to
   reach Complete is an Advance that returns ActionComplete on a tick that
   already passed the breach gate. Could not break.
2. Breach check ordering vs Complete: considered whether Complete increments
   steps and thus could trip the check on the SAME tick it completes. It does
   not — isProgressAction(ActionComplete)==false (loop.go:128), and the breach
   check is at loop top, before Advance, so the completing tick's check ran with
   the pre-Complete step count. Could not break.
3. Wall-clock reset bug: hunted for any reassignment of `start` or a
   per-progress reset. `start` (loop.go:34) is assigned once and only read
   thereafter; the per-tick `deadline` (loop.go:33) is the only thing that
   resets on progress. Distinct, correct.
4. Cost panic on malformed event: confirmed the comma-ok type assertion
   (loop.go:187) means a non-float cost_usd is ignored, total unchanged. No
   panic. Could not break. (JSON unmarshal into map[string]any yields float64
   for JSON numbers by default, so the happy path also works.)
5. Zero-value Cursor on first-tick breach: if a wall-clock ceiling is already
   exceeded before the first Advance (start in the past is impossible since
   start=time.Now(), but a MaxWallClock of ~0ns could trip on a slow first
   iteration), escalateLoopBudget receives the zero Cursor `last` (loop.go:36).
   escalate() writes phase="" into the inbox note — cosmetic, not a correctness
   bug; the note is still durable and blocking. Not a finding.
6. Config precedence: mergeDefaults runs per layer (central → deck), so a deck
   [defaults.loop] overrides the central seed, and `loopBudget` only applies CLI
   overrides when >0. The FINAL.md precedence chain (built-in → ~/.parley → deck
   → run flags) is accurately implemented. Could not break.
7. continueAuto path: continues inherit ~/.parley defaults via loopBudget(root,
   0,0) but have no per-invocation flag (documented deviation). Not a
   fail-open — defaults still apply. Could not break.

## Findings

No CRITICAL or MAJOR findings. The core safety invariant (escalate-not-complete)
and the breach logic are sound. Below are MINOR/NIT observations.

MINOR — No end-to-end Run() breach test exercises the real Advance loop.
The new tests (loop_budget_test.go) unit-test loopBudgetBreach, isProgressAction,
emitLoopBudget, loopCostUSD, and escalateLoopBudget in isolation, but none drive
d.Run(ctx) with a real multi-step fakeRunner to the point of a mid-loop breach
and assert (a) the inbox note is written AND (b) Run returns nil without
reaching ActionComplete. The FINAL.md acceptance ("with MaxDriverSteps: 1, a
multi-step idea escalates with a loop-budget inbox note instead of completing")
is the headline guarantee and is only verified compositionally, not
end-to-end. The unit pieces are correct, so this is a coverage gap, not a
defect — but a reviewer refuting escalate-not-complete would want the
integration assertion. Concrete fix: add a test that constructs a driver with
MaxDriverSteps=1 over a fakeRunner/ImplOps that emits ≥2 progress actions, calls
d.Run(ctx), and asserts the loop-budget inbox note exists and the idea status is
not "done"/"implemented".

MINOR — loopCostUSD re-reads the entire event log on every breach check and
every emitLoopBudget. loopBudgetBreach is called once per tick and calls
loopCostUSD (Load + full scan) when MaxCostUSD>0; emitLoopBudget also calls it
per progress step. For a long run (hundreds of steps, thousands of events) this
is O(steps × events) re-reads. Today it is inert (no agent.usage events), so it
is not a live problem, but the moment agent-usage-telemetry lands this becomes a
quadratic hot path inside the auto-drive loop. Concrete fix: when LE-6
telemetry is wired, either maintain a running cost total (incremented on each
agent.usage append) or cache the last-seen event count. Not blocking now; flag
for the agent-usage-telemetry follow-up.

NIT — escalateLoopBudget ignores the error from escalate(). escalate() (the
shared path) currently swallows all its own errors (os.MkdirAll / WriteFile
return values are discarded with `_ =`), so escalateLoopBudget returning nil is
consistent with the existing escalateDeadline/driver-error callers. No behavior
change, but the "durable inbox note" guarantee is best-effort-by-omission across
all escalation paths, not just this one. Not a regression introduced here; noted
for completeness.

## Open questions

1. Should `continue --auto` gain `--max-driver-steps` / `--max-wall-clock`
   flags for parity with `run`? The IMPLEMENTATION.md deviation notes this is
   scoped down intentionally. If a user relies on `continue` as the primary
   re-drive entry, the absence of per-invocation overrides may surprise — but
   ~/.parley defaults still apply, so it is not unsafe. Confirm the scoping
   decision is acceptable for the merge.

2. The seed template sets max_cost_usd = 0 (unlimited). Given the ceiling is
   inert today, a user who later sets a non-zero max_cost_usd in ~/.parley
   before agent-usage telemetry lands will get a ceiling that never trips (cost
   stays 0). Is that the intended user experience, or should the field be
   documented as "no-op until agent.usage is emitted" directly in the seed
   template comment (it currently says "telemetry-gated (LE-6)" which is
   accurate but terse)?
