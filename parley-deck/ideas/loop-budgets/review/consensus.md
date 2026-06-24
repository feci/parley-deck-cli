---
idea: loop-budgets
review-cycle: 2
drafted-by: claude-1
date: 2026-06-24
outstanding_agreed_fixes: 0
blocked: false
---

## Cycle 2 resolution

All 3 cycle-1 agreed fixes (F-T2-1/2/3) applied; round-02 re-review was **3/3 zero
findings** (codex-1, hermes-1, antigravity-1) — codex's MAJOR `0`-override and both MINORs
confirmed resolved, no fail-open or safety regression introduced. Outstanding agreed
fixes: 0.

## Agreed fixes

From review/round-01 (refutation): antigravity-1 and hermes-1 reported **zero findings**;
codex-1 raised 1 MAJOR + 2 MINOR, all accepted.

- **F-T2-1 [MAJOR] (codex-1)** — `0` cannot override a seeded loop budget. `mergeDefaults`
  only copied loop values `> 0` and `loopBudget` only applied run overrides `> 0`, so once
  `parley init` seeds `max_driver_steps=200`, neither a deck `= 0` nor `--max-driver-steps 0`
  could opt back into unlimited — contradicting the documented `0 = unlimited` escape hatch.
  Fix: `loopBlock` fields are now `*int`/`*float64` (presence-aware merge — explicit 0
  overrides); run flags use `flag.Visit` so an explicit `--max-driver-steps=0` /
  `--max-wall-clock=0` wins while absence still means "use defaults"; flag help corrected.
- **F-T2-2 [MINOR] (codex-1)** — `loop.budget` reported `cost_usd: 0` whenever cost was
  unlimited, hiding observed spend once `agent.usage` lands. Fix: `emitLoopBudget` now
  computes `loopCostUSD()` unconditionally for the event; only enforcement stays gated by
  `MaxCostUSD > 0`.
- **F-T2-3 [MINOR] (codex-1)** — the breach path was helper-tested but not `Run`-tested.
  Fix: added `TestRunEscalatesOnLoopBudget` (a 1ns wall-clock budget trips the first
  pre-Advance check → inbox note, never Complete) + `TestEmitLoopBudgetReportsCostWhenUnlimited`
  + `TestLoopDefaultsDeckZeroOverridesCentralSeed`.

## Deferred follow-ups
- **`agent.usage` emission from the runner** (LE-6 full activation). The cost ceiling is
  plumbed but inert until a runner emits `agent.usage{cost_usd}`. Tracked as
  `agent-usage-telemetry`; out of scope for the safety-floor tier.

## Dismissed findings
- **codex-1 open question** (separate `--no-loop-budget` flag vs explicit `=0`). Resolved
  per FINAL.md: explicit `--max-...=0` IS the unlimited override (now implemented via
  flag.Visit); no separate flag is added — fewer knobs, matches the documented semantics.

## Coverage & blind spots
codex-1 (CLI/engine lens) carried the review; antigravity-1 (risk) and hermes-1
(verification) independently confirmed the escalate-not-complete invariant holds and found
no fail-open path. The `0`-override gap was a single-reviewer catch but a real
acceptance-criterion miss — fixed with presence-aware config + flag.Visit.

## Signoffs

**Outcome: 3/3 ACCEPT (unanimous, no BLOCK) → Phase 8 complete.**

### Signoff: codex-1 — 2026-06-24
Status: ACCEPT
The round-01 MAJOR (`0`-override) is fixed and my round-02 reported zero findings: deck
config can now override the central seed to unlimited and `parley run --max-driver-steps=0`
/ `--max-wall-clock=0` override inherited defaults; the two MINORs (cost observability,
Run-level breach test) also hold.
_(Facilitator note: codex-1 authored this signoff but its filesystem sandbox blocked the
write — the target path `../signoff-codex-1.md` resolved outside its writable root.
Transcribed verbatim from its run output per the protocol's tooling-failure exception;
source at `tasks/b7r6ehq3q.output`.)_

### Signoff: hermes-1 — 2026-06-24
Status: ACCEPT
I reported zero findings in both round-01 and round-02; the three agreed fixes (F-T2-1
presence-aware `0`-override, F-T2-2 unconditional cost observability, F-T2-3 Run-level
breach test) are correctly implemented with no fail-open path, and all three reviewers
independently converged on zero findings in round-02 with touched-package tests green. The
deferred follow-ups (agent.usage emission, loopCostUSD O(steps × events) re-reads) remain
explicitly out of scope and tracked.

### Signoff: antigravity-1 — 2026-06-24
Status: ACCEPT
The fix-up implementation successfully addresses the unlimited budget escape hatch. The
presence-aware merges and CLI override logic have been thoroughly verified and introduce no
safety regressions.
