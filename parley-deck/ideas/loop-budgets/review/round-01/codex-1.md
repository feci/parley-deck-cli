---
agent: codex-1
idea: loop-budgets
review-round: 1
date: 2026-06-24
---

## Summary

I found one MAJOR semantics bug: `0` is not truly unlimited once seeded defaults exist, because both layered config merge and run flag handling treat zero as "unset/use defaults". The driver breach path itself looks structurally non-completing once `loopBudgetBreach` returns a reason, and the three production driver construction sites all receive the ceilings.

I also found a cost telemetry gap and a test gap around the full `Run` path. Validation run: `go test -count=1 ./internal/driver/... ./internal/app/... ./internal/config/... ./internal/protocol/...` passed. The broader requested command including `./internal/runner/...` failed in `TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`, which appears environment-sensitive and unrelated to this diff.

## Refutation attempts

- LE-5 step breach / never Complete: I traced `Driver.Run` before the `Advance` call. If `loopBudgetBreach` returns a reason, `Run` returns through `escalateLoopBudget` before the switch can reach `ActionComplete`. With `MaxDriverSteps: 1`, the first progress action increments `steps` to 1; the next tick checks `steps >= max` before advancing again, so a multi-step idea halts before the next action can complete it. I did not find a Complete path after an actual breach.
- LE-5 wall-clock vs `roundDeadline`: `start` is created once per `Run` invocation and is not reset when progress actions reset `deadline`. The per-tick `roundDeadline` reset remains separate. A wall budget that is already elapsed should trip before the next `Advance`.
- LE-5 zero / backward compatibility: direct driver config with zero ceilings is unlimited, but the config and CLI layers break this once a seeded nonzero default exists. See MAJOR finding.
- LE-5 construction sites: `rg` found exactly three production `driver.New(driver.Config{...})` sites in `internal/app/app.go`, and all three pass `MaxDriverSteps`, `MaxWallClock`, and `MaxCostUSD`.
- LE-5 defaults / flags precedence: nonzero central, deck, and run values follow the intended layering by inspection, but explicit zero cannot win at higher precedence. See MAJOR finding.
- LE-6 cost: `loopCostUSD` sums `agent.usage` events when they are JSON numeric values. I tried the unlimited-cost case and found the `loop.budget` event would still report `cost_usd: 0` even if usage events exist. See MINOR finding.
- LE-10 watcher remediation: `openRemediationIdea` now writes `status: candidate`, omits `participants: []`, and the new test asserts those properties. I did not find an active empty-quorum remediation path in this diff.

## Findings

### [MAJOR] `0` cannot override seeded loop budgets

`0` is specified as unlimited and the precedence is built-in -> `~/.parley [defaults.loop]` -> deck `agents.toml` -> run flags. The implementation cannot express an explicit zero at either higher-precedence layer: `mergeDefaults` only copies loop values when they are `> 0` (`internal/config/runtime.go:202`), and `loopBudget` only applies run overrides when they are `> 0` (`internal/app/app.go:1125`). The new run flag help also says `0 = use ~/.parley [defaults.loop]` (`internal/app/app.go:1720`), which contradicts FINAL.md and the protocol text that say `0` means unlimited.

Why it matters: after `parley init` seeds `max_driver_steps = 200` and `max_wall_clock_ms = 7200000`, a deck cannot opt back into backward-compatible unlimited behavior with `[defaults.loop] max_driver_steps = 0`, and a user cannot force unlimited for one run with `--max-driver-steps 0` or `--max-wall-clock 0`. The safety ceiling works, but the escape hatch required by the acceptance criteria does not.

Concrete fix: make loop TOML fields presence-aware, for example `*int` / `*float64` in `loopBlock`, and merge when the pointer is non-nil so explicit zero overwrites earlier layers. For CLI flags, track flag presence with `FlagSet.Visit` or an equivalent wrapper and pass presence booleans into `loopBudget`, so explicit `--max-driver-steps=0` and `--max-wall-clock=0` override defaults to unlimited while absence still means "use defaults". Add tests for central nonzero -> deck zero and central nonzero -> run zero.

### [MINOR] `loop.budget` hides cost telemetry when cost is unlimited

`emitLoopBudget` only calls `loopCostUSD` when `MaxCostUSD > 0` (`internal/driver/loop.go:155`). That means once runners start emitting `agent.usage`, any run with unlimited cost (`MaxCostUSD == 0`) will still emit `loop.budget` events with `cost_usd: 0`, even though LE-6 says the event should carry the summed cost telemetry.

Why it matters: unlimited cost should disable enforcement, not hide observed spend from the TUI/state. This would make the telemetry look inert even after runner usage events land.

Concrete fix: compute `cost := d.loopCostUSD()` unconditionally in `emitLoopBudget`; keep only enforcement gated by `MaxCostUSD > 0`. Add a test with an `agent.usage` event and `MaxCostUSD == 0` asserting the `loop.budget` event reports the summed cost.

### [MINOR] The critical breach behavior is only helper-tested, not `Run`-tested

The new tests cover `loopBudgetBreach` directly and call `escalateLoopBudget` directly, but they do not run a multi-step `Driver.Run` scenario with `MaxDriverSteps: 1` and assert that the loop writes the budget inbox note and never reaches `ActionComplete` / `status: complete`.

Why it matters: the highest-risk acceptance criterion is the integration of the pre-`Advance` budget check, action classification, switch handling, and durable escalation. Helper tests would not catch a future regression where a progress action is misclassified or the switch completes before the budget check is applied.

Concrete fix: add a `Run`-level test using fake ops/state for a multi-step idea: set `MaxDriverSteps: 1`, drive until halt, assert the loop-budget inbox note exists, no complete status was written, and no `driver: idea ... complete` path was reached. Add the analogous tiny wall-clock test with an already-expired effective budget.

## Open questions

- Should `parley run --max-driver-steps=0` and `--max-wall-clock=0` be the explicit unlimited override over seeded defaults, as FINAL.md says, or should the CLI expose a separate `--no-loop-budget` / `--unlimited-loop-budget` flag and update the protocol text?
