---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
note: n1-known-exhaustion-fix
context_mode: full
attestation_source_sha256: 4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a
---

# kimi-1 — MINOR-1: known cross-resource exhaustion preflight (fix + tests only)

## Task scope (bounded)
- Non-mutating preflight of known exhausted step/cycle **before EITHER first charge**.
- Three wiring points:
  1. driver `stepRoundRunner` — before `ChargeStep`, must check **cycle**.
  2. driver `reserveFixupCycle` — before `ChargeStep`, must check **cycle**.
  3. runner `reserveBudget` — before `ChargeCycle`, must check **step**.
- Preserve charged-session reuse at cap (a session already charged for this logical
  action stays reusable); late-race failures stay spent.
- If existing session state is private, add narrowly scoped exported preflight helpers
  in `internal/budget/step_session.go` and `internal/budget/cycle_session.go` that
  follow existing Charge authority/reuse semantics without persistence/charging.
- NO budget migration, binding schema, cap change, refunds, or atomicity claims.
  Distinct request receipts are intentional — no dedup.
- Bounds: 15 min / 60 steps. No Bash/tests/Git/subagents/providers. Source candidate
  only; coordinator runs tests after terminal.

## Retained state (do not revert)
- My fresh-scope MAJOR-5 Go edits are retained.
- Coordinator added ONLY `MaxDriverSteps: 5` to the focused fixture, preserving
  assertions; the corrected test is already copied into this editable root.
- Coordinator evidence: compile 1.872s; focused 15 pass / 1 fail (fixture omitted
  finite MaxDriverSteps) → after fixture fix compile 1.901s + 16 focused events
  2.423s passed. Old failed snapshot retained.
- No broad/full/race/vet pass claimed by anyone.

## Plan
1. Read the 6 allocated source/test files.
2. Add read-only exported preflight helpers in budget step/cycle session files
   (reuse-aware: reuse of an already-charged session is allowed at cap; only a
   *new* charge is refused when the counter is at cap).
3. Wire checks at the three call sites (driver ×2, runner ×1).
4. Extend tests: exhausted step cannot spend cycle; exhausted cycle cannot spend
   step; already-charged grouped session at cap reusable; actual late race/refusal
   still spent.

## Status — COMPLETE (source candidate; coordinator runs tests after terminal)

### What changed
- `internal/budget/step_session.go`: new exported `PreflightStepCharge(ctx, *StepBinding)`
  — read-only known-refusal probe. Reuse-first, not naive current-count-only: a live
  already-charged session returns its cached result (nil = reusable at the inclusive
  cap; a cached refusal replays); only a would-be NEW reservation probes the binding's
  existing non-mutating `Check` (the same validation `OpenStepSession` runs). No
  charge, no persistence, no guard, no refund.
- `internal/budget/cycle_session.go`: new exported `PreflightCycleCharge(ctx, *CycleBinding)`
  — same reuse-first shape; a NEW reservation refuses only when
  `Count(state) >= Policy.Maximum` (mirrors `cycleLimits` denial; zero maximum forbids
  the operation). Added the `fmt` import.
- `internal/driver/budget.go` (`stepRoundRunner.RunRound`): new `preflightRunnerCycle`
  gate between `precheckRoundLaunch` and `budget.ChargeStep`.
- `internal/driver/cycle_budget.go`: added `preflightRunnerCycle` (loads the admitted
  CrossReview binding, wraps errors `cross-review accounting/budget` exactly like the
  charge path); `reserveFixupCycle` now calls `budget.PreflightCycleCharge` after
  `OpenCycleSession`, before `ChargeStep`.
- `internal/runner/launch_budget.go` (`reserveBudget`): inside the cycle-kind block,
  loads the step binding and calls `budget.PreflightStepCharge` BEFORE `ChargeCycle`;
  `JoinStepSession`/`ChargeStep` stays the authority for reuse and late races.

### Tests (written, NOT run — coordinator runs them)
- driver `launch_precheck_test.go`:
  - `TestN1ExhaustedCrossReviewCycleRefusesBeforeStep` — exhausted cycle cannot spend
    step (round path via Advance; no dispatch, no step charge, cycle unchanged).
  - `TestN1ReserveFixupCyclePreflightsExhaustionBeforeStep` — same on the fix-up path
    (direct `reserveFixupCycle`; step stays 0).
  - `TestN1ReserveFixupCycleChargedSessionReusableAtCap` — already charged grouped
    session at cap reusable through the full wiring (ordinal 1 twice; 1 step + 1 cycle).
- runner `launch_budget_preflight_test.go` (new):
  - `TestLaunchPreflightExhaustedStepCannotSpendCycle` — full RunMeasured: known
    exhausted step refuses before the cycle charge (cycle 0, step 1, one
    `budget_refused` terminal, never started).
  - `TestLaunchPreflightChargedGroupedSessionReusableAtCap` — both budgets at cap with
    charged grouped sessions: preflights admit reuse; re-charges reuse; counts stay 1/1.
  - `TestLaunchPreflightLateRaceRefusalStaysSpent` — deterministic interleave in
    reserveBudget order: two launches pass preflight at one remaining step; winner
    charges cycle+step; loser charges cycle then refuses at the real step authority;
    the loser's cycle stays spent (cycle 2, step 1).

### Explicit non-claims (per task)
- No budget migration, no binding schema change, no cap change, no refunds, no
  atomicity claims; distinct request receipts remain intentional (no dedup).
- The preflight is a known-state read, NOT a lock: concurrent exhaustion after it
  still refuses at the real charge and stays spent.
- Monetary/launch-kind Reserve ordering unchanged; it stays last in `reserveBudget`.
- Existing MAJOR-5 fresh-scope edits and the coordinator's `MaxDriverSteps: 5`
  fixture are untouched; the old failed snapshot is retained.
- I ran NO compiler or tests (no Bash per launch bounds). Coordinator evidence
  before these edits: compile 1.901s, 16 focused events 2.423s passed. Post-edit
  compile/focused runs are the coordinator's; no broad/full/race/vet pass claimed.

### Residual risk (honest)
- `PreflightStepCharge` replays any cached session error (not only exhaustion) —
  documented; identical to what `ChargeStep` would replay moments later.
- A charged session whose receipt/time revalidation fails at reuse is a late failure
  by design (the other resource may already be spent) — unchanged pre-change semantics.
- `OpenStepSession`'s own `Check` at fresh-open already refused exhausted steps at
  join time; the runner change only moves that knowledge BEFORE the cycle charge.
