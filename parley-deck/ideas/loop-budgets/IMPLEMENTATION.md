---
idea: loop-budgets
status: complete
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
branch: parley-deck-cli#loop-engineering-impl
head-commit: (this commit)
---

## Summary of work

Tier 2 (the safety floor): explicit, escalate-on-breach loop ceilings on the auto-drive
driver (LE-5/6) + the non-solo fix for the pipeline watcher (LE-10), per `FINAL.md`.
`go build`, `go vet`, `go test -count=1 ./...`, and the drift guard are green.

## Implementation plan / checklist

- [x] **LE-5 loop ceilings** — `driver.Config` gains `MaxDriverSteps int`,
  `MaxWallClock time.Duration`, `MaxCostUSD float64` (`internal/driver/driver.go`; `0` =
  unlimited, backward-compatible). `internal/driver/loop.go` enforces them in `Run`: a
  step counter + wall-clock origin, checked before each `Advance`; a breach calls
  `escalateLoopBudget` (durable inbox note, shared `escalate` path) and halts — never
  `Complete`. `isProgressAction` gates step counting; `emitLoopBudget` records a
  `loop.budget` event per progress step.
- [x] **LE-6 cost (best-effort)** — `loopCostUSD` sums `cost_usd` across `agent.usage`
  events; the `MaxCostUSD` ceiling + `loop.budget` cost field consume it. **Honest limit:**
  the runners do not yet emit `agent.usage`, so cost is `0` in practice and `MaxCostUSD`
  is inert until that telemetry lands (deferred follow-up `agent-usage-telemetry`). The
  contract + plumbing are here.
- [x] **Defaults + flags** — `internal/config/runtime.go`: `[defaults.loop]` block
  (`max_driver_steps`/`max_wall_clock_ms`/`max_cost_usd`) parsed into `CentralDefaults` +
  merged + seeded by `centralDefaultTemplate` (steps 200, wall-clock 2h, cost 0). `app.go`:
  `loopBudget(root, stepsFlag, wallFlag)` resolves central defaults + `run` flag overrides
  and applies them at all 3 driver-construction sites; `run --max-driver-steps` /
  `--max-wall-clock` added.
- [x] **LE-10** — `internal/app/pipeline_cmd.go` `openRemediationIdea` now writes
  `status: candidate` (not `round-01`) with no `participants: []` claim + a `## Promotion`
  note (a human/the manifest staffs the quorum before deliberation).
- [x] **Protocol** — §4 "Stopping judgment" gains the loop-budget invariant; §12.11 gains
  the candidate-remediation rule. Mirrored in both `COOPERATION.md` copies (drift green).
- [x] **Tests** — `internal/driver/loop_budget_test.go` (breach/step/wall-clock/0-unlimited,
  `isProgressAction`, `loop.budget` event, `loopCostUSD`, inbox escalation);
  `internal/config/loop_defaults_test.go` (`[defaults.loop]` parse/merge + seed template);
  `internal/app/pipeline_remediation_test.go` (candidate, no empty quorum).

## Deviations from FINAL.md

- **MaxCostUSD is plumbed but inert** (no runner `agent.usage` emission yet) — see LE-6
  honest limit above. The field/event/summing exist so the ceiling activates the moment a
  runner emits usage; full emission is the deferred `agent-usage-telemetry` follow-up.
- **CLI flags only on `run`** (not `continue`): `continueAuto` inherits the `~/.parley`
  defaults but has no per-invocation flag override. Scoped down — `run` is the primary
  auto-drive entry; adding the same flags to `continue` is a trivial later add if wanted.

## Notes for reviewers

- **Backward-compat:** `0` ceilings = unlimited; existing decks/tests (which construct
  `Config` without the fields) are unchanged — verify the driver suite is green and that
  a default-constructed driver never escalates on budget.
- **Wall-clock vs roundDeadline:** `MaxWallClock` is a TOTAL run budget (does not reset);
  the existing 30-min `roundDeadline` is per-tick (resets on progress). They are distinct.
- **Escalate-not-complete:** confirm a breach returns via `escalateLoopBudget` (inbox note)
  and that no `Complete` path is reachable after a breach.
- Try to break: a step/wall-clock ceiling of 1; a `[defaults.loop]` deck override vs the
  central seed; an `agent.usage` event with a non-float `cost_usd`.

## Fix-up cycle 1
status: complete
completed: 2026-06-24

### Fixes applied
(From review/round-01: antigravity-1 + hermes-1 zero findings; codex-1 1 MAJOR + 2 MINOR.)
- **F-T2-1 [MAJOR]** — `0` now overrides a seeded budget. `loopBlock` fields are
  `*int`/`*float64` (presence-aware `mergeDefaults` — explicit `= 0` overrides a lower
  layer); `run` flags use `flag.Visit` so an explicit `--max-driver-steps=0` /
  `--max-wall-clock=0` means unlimited while omission uses defaults; flag help corrected to
  "explicit 0 = unlimited; omit to use ~/.parley". (`internal/config/runtime.go`,
  `internal/app/app.go`.)
- **F-T2-2 [MINOR]** — `emitLoopBudget` now computes `loopCostUSD()` unconditionally so the
  `loop.budget` event reports observed spend even when cost is unlimited; only enforcement
  stays gated by `MaxCostUSD > 0`. (`internal/driver/loop.go`.)
- **F-T2-3 [MINOR]** — added `TestRunEscalatesOnLoopBudget` (Run-level breach → inbox note,
  never Complete), `TestEmitLoopBudgetReportsCostWhenUnlimited`, and
  `TestLoopDefaultsDeckZeroOverridesCentralSeed`.

### Deviations from agreed fixes
- codex-1's open question (separate `--no-loop-budget` flag) was resolved per FINAL.md:
  explicit `=0` is the unlimited override; no extra flag added.

`go build`, `go vet`, `go test -count=1 ./...`, and the drift guard are green.
