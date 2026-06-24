---
agent: hermes-1
idea: loop-budgets
review-round: 2
date: 2026-06-24
---

## Summary

Re-review (round-02, refutation) of codex-1's three fix-up commits
(cfcbe56..HEAD) for the loop-budgets idea. I reported zero findings in round-01;
this round focuses exclusively on whether the fixes introduced a fail-open path
in budget enforcement — specifically the presence-aware config merge (F-T2-1),
the flag.Visit run-flag override (F-T2-1), and emitLoopBudget's unconditional
loopCostUSD load (F-T2-2).

Verdict: the fixes hold. I could not construct any path where a budget the user
wanted is accidentally disabled, nor any path where the unconditional cost load
weakens enforcement. Zero findings. `go test -count=1 ./internal/config/...
./internal/driver/...` and `go vet ./internal/config/... ./internal/driver/...
./internal/app/...` are green.

## Refutation attempts

A. Presence-aware merge can disable a budget the user wanted (F-T2-1, config
   layer).
   `loopBlock` fields are now `*int`/`*float64` (runtime.go:50-54). With
   go-toml/v2, a present TOML key (even `= 0`) unmarshals to a non-nil pointer;
   an absent key leaves the pointer nil. `mergeDefaults` (runtime.go:205-217)
   copies `*gd.Loop.<Field>` only when the pointer is non-nil.
   - Central seed `max_driver_steps = 200` → out.MaxDriverSteps = 200.
   - Deck `max_driver_steps = 0` → pointer non-nil → out.MaxDriverSteps = 0
     (unlimited). This is the user's explicit override, not an accident.
   - Deck omits `max_driver_steps` → pointer nil → falls through, central 200
     retained. The user's seeded budget is preserved.
   The only way a budget becomes 0 (unlimited) at this layer is the user
   explicitly writing `= 0`, which is the documented escape hatch. Omitting the
   key preserves the lower layer. Could not break.
   Test `TestLoopDefaultsDeckZeroOverridesCentralSeed` asserts exactly:
   central 200 + deck 0 → 0; central 7200000 + deck-absent wall → 7200000.

B. flag.Visit can accidentally disable a budget (F-T2-1, CLI layer).
   `fs.Visit` (app.go:1738-1745) only visits flags the user explicitly set on
   the command line. Omitting `--max-driver-steps` means the override block
   never executes; `lbSteps` retains `loopBudget(root)` (the merged config).
   Only an explicit `--max-driver-steps=0` sets `lbSteps = 0`, which is the
   user's deliberate unlimited choice. No accidental disabling path.
   Verified all three `driver.New` construction sites receive the
   post-override values:
   - run --no-tui path (app.go:1823): downstream of the flag.Visit block
     (1738). lbSteps/lbWall/lbCost reflect the override.
   - run TUI path (app.go:1877): same variables captured in the closure,
     downstream of 1738. Correct.
   - continueAuto (app.go:1150): `loopBudget(root)` with no per-invocation
     flags (documented deviation). Not a fail-open — ~/.parley defaults still
     apply.
   Could not break.

C. emitLoopBudget's unconditional loopCostUSD load weakens enforcement (F-T2-2).
   `emitLoopBudget` (loop.go:157-158) now calls `d.loopCostUSD()` always.
   Enforcement lives in a separate function, `loopBudgetBreach` (loop.go:137-152),
   which still gates cost enforcement on `MaxCostUSD > 0` (line 146). The
   unconditional load in `emitLoopBudget` only affects what is *reported* in the
   `loop.budget` observability event, not what is *enforced*. The two code paths
   are disjoint. No enforcement weakening.
   Safety of the load itself: `loopCostUSD` (loop.go:176-191) is read-only —
   Load error → return 0; `agent.usage` events' `cost_usd` is summed via a
   comma-ok type assertion (`f, ok := e.Data["cost_usd"].(float64)`), so a
   non-float value is skipped, not a panic. No mutation, no side effects.
   Test `TestEmitLoopBudgetReportsCostWhenUnlimited` confirms cost is reported
   even when MaxCostUSD == 0. Could not break.

D. Stale pre-override values reach a driver.New site.
   The flag.Visit block (1738) runs before both run-path driver.New sites
   (1823, 1877). continueAuto (1150) has no flags. No stale values possible.
   Could not break.

E. Pointer change in loopBlock breaks existing TOML parsing or introduces a
   zero-value ambiguity in CentralDefaults.
   CentralDefaults fields remain plain int/float64 (runtime.go:69-71); the
   pointers live only in loopBlock, which is the merge-time decision layer.
   After merge, a flat 0 means "unlimited" whether it came from "no layer set
   it" or "a layer explicitly set 0" — both are correct per the documented
   `0 = unlimited` semantics. A positive value is the only path to enforcement.
   No ambiguity that leads to fail-open. The config + driver test suites pass.
   Could not break.

## Findings

None. Zero CRITICAL/MAJOR/MINOR/NIT.

The three fixes are correctly implemented and do not introduce a fail-open path:
- F-T2-1: presence-aware merge (`*int`/`*float64` + nil-check) and flag.Visit
  distinguish "explicit 0 = unlimited" from "absent = use lower layer" at both
  the config and CLI layers. A budget the user wanted is preserved unless the
  user explicitly opts out with `= 0`.
- F-T2-2: the unconditional loopCostUSD load is observability-only and
  panic-safe; enforcement remains gated by `MaxCostUSD > 0` in the separate
  `loopBudgetBreach` function.
- F-T2-3: the three new tests cover the Run-level breach path, the
  unlimited-cost telemetry, and the deck-zero-override config path.

## Open questions

None. The fixes resolve all three round-01 agreed fixes. The deferred
follow-ups (agent.usage emission, loopCostUSD O(steps × events) re-reads) remain
out of scope for this safety-floor tier, as recorded in consensus.md.
