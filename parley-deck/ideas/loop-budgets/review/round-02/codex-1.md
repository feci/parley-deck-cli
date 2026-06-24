---
agent: codex-1
idea: loop-budgets
review-round: 2
date: 2026-06-24
---

## Summary

The round-01 MAJOR is fixed. I could not reproduce a path where an explicit `0` is treated as "unset" after the fix-up: deck config can now override the central seed to unlimited, and `parley run --max-driver-steps=0` / `--max-wall-clock=0` override inherited defaults for the run driver construction paths.

The two MINOR fixes also hold: `loop.budget` reports observed cost even when cost enforcement is unlimited, and the driver now has a Run-level breach test that exercises escalation rather than only helper behavior.

## Refutation attempts

- **F-T2-1 deck zero vs central seed:** inspected `internal/config/runtime.go`; `loopBlock` now uses `*int` / `*float64` fields and `mergeDefaults` copies any non-nil pointer, including `0`. Ran `go test -count=1 ./internal/config -run 'TestLoadDefaultsLoopBlock|TestLoopDefaultsDeckZeroOverridesCentralSeed|TestCentralDefaultTemplateHasLoopBlock'`, which passed. The deck-zero test specifically keeps an inherited wall-clock seed while overriding only driver steps to `0`.
- **F-T2-1 run flag zero vs seeded defaults:** inspected `internal/app/app.go`; `loopBudget(root)` resolves layered defaults first, then `fs.Visit` overwrites only flags the user explicitly supplied. Both no-TUI and live-TUI `driver.New` sites receive the resulting `lbSteps` / `lbWall` tuple. `go run ./cmd/parley run -h` shows the corrected help text: explicit `0` means unlimited; omit means use `~/.parley [defaults.loop]`.
- **F-T2-1 escape-hatch break attempt:** tried to find a remaining `> 0` merge or override path with `rg` across `internal`; the remaining `> 0` checks are enforcement checks in `loopBudgetBreach`, where `0 = unlimited` is the intended behavior. `continue --auto` has no per-invocation loop-budget flags and consumes layered defaults only, which matches the documented implementation deviation rather than reintroducing the original run-flag bug.
- **F-T2-2 cost observability:** inspected `internal/driver/loop.go`; `emitLoopBudget` now calls `loopCostUSD()` unconditionally and only `loopBudgetBreach` gates enforcement on `MaxCostUSD > 0`. Ran `go test -count=1 ./internal/driver -run 'TestLoopBudgetBreach|TestEmitLoopBudgetReportsCostWhenUnlimited|TestRunEscalatesOnLoopBudget|TestEscalateLoopBudgetWritesInbox|TestLoopCostUSDSumsAgentUsage'`, which passed.
- **F-T2-3 Run-level breach path:** inspected `TestRunEscalatesOnLoopBudget`; it constructs a real `Driver.Run` with `MaxWallClock: time.Nanosecond`, verifies clean halt, and asserts the loop-budget inbox note exists. This covers the pre-`Advance` breach path that was missing in round 01.
- **Broader checks:** `go test -count=1 ./internal/config ./internal/driver` passed, and `go test -count=1 ./internal/app -run Test` passed. `go test -count=1 ./...` still fails in unrelated `internal/runner` durable-kill coverage with `process verification failed (no recorded boot id); not killed`; loop-budget packages passed in that run.

## Findings

None.

## Open questions

None.
