---
agent: antigravity-1
idea: loop-budgets
review-round: 2
date: 2026-06-24
---

## Summary

This review covers Phase 6, Round 2 of the `loop-budgets` idea, focusing on the fix-up implementation (commit `455a947` compared against `cfcbe56`). The primary objective of this review is to verify that the presence-aware pointer merging for layered configuration (`*int`/`*float64`) and `flag.Visit` logic for command-line arguments correctly resolve the `0`-override (unlimited) escape hatch without introducing safety issues (such as silent or unexpected budget disabling).

I have audited the code changes, conducted refutation risk analysis, verified the integration tests, and run the project's test suite (`go test ./...` which passes cleanly). I report **zero findings** for this round.

## Refutation attempts

I systematically attempted to break the safety of the new budget precedence and override mechanics through the following vectors:

1. **Silent / Unexpected Disabling of Budgets by Project Decks:**
   - **Scenario:** Can a project-level `agents.toml` configuration silently zero out safety budgets that a system operator expected to be globally enforced via `~/.parley/agents.toml`?
   - **Verification:** I analyzed the presence-aware TOML parsing structure in [internal/config/runtime.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/config/runtime.go). If a project deck has `max_driver_steps = 0`, it overrides the central default value to `0` (unlimited). While this disables the budget, it is the correct and documented precedence model (project configuration overrides central/user-global settings). If the operator wants to override the project's defaults, they must use command-line flags. 
   - **Absence Check:** If the project deck completely omits the `[defaults.loop]` section or specific keys (e.g. `max_driver_steps`), the pointer remains `nil`. The presence-aware merge logic in [mergeDefaults](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/config/runtime.go#L178-L218) ignores `nil` pointers, correctly allowing the central default values (e.g., `200` steps) to fall through and remain active. This is verified by [TestLoopDefaultsDeckZeroOverridesCentralSeed](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/config/loop_defaults_test.go#L39-L64) in `internal/config/loop_defaults_test.go`.

2. **Exploiting `flag.Visit` CLI Flag Bypasses:**
   - **Scenario:** Can command-line flags interact with configuration files in a way that causes safety budgets to be unexpectedly disabled or misconfigured?
   - **Verification:** I traced the flag resolution inside `runTask` in [internal/app/app.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go). By calling `fs.Visit` to iterate over explicitly passed flags, the application distinguishes between a user omitting a flag (which defaults to the value `0` but is NOT visited, leaving the configuration-layer budget intact) and a user explicitly passing `--max-driver-steps 0` (which is visited, overriding the budget to `0` / unlimited). This behaves exactly as intended, protecting the safety floor while preserving the escape hatch.

3. **Malformed Configuration Files Bypassing Budgets:**
   - **Scenario:** If a configuration file contains a TOML syntax or structural error, does `LoadDefaults` fail-open and silently disable all safety budgets?
   - **Verification:** Although the [loopBudget](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L1121-L1128) function swallows errors returned by `LoadDefaults`, any TOML syntax or structural errors in the config layers are also caught during agent discovery via `LoadAgentSpecs` / [discoverConfigured](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L1757). Because `discoverConfigured` returns the error and halts the execution before any driver ticks are run, there is no risk of running ideas with disabled safety budgets due to malformed configs.

4. **Transient Flag Persistence on `continue --auto`:**
   - **Scenario:** Does running `parley continue --auto` on an idea that previously hit a transient CLI-specified budget run open/unlimited?
   - **Verification:** `continue --auto` does not receive run flags and defaults back to the persistent central/deck configuration. Since these configuration files either have standard safety budgets configured (like the `200` step central seed) or project-specific defaults, it does not default to unlimited unless the config files themselves explicitly opt for unlimited. This is safe and prevents unbounded execution.

5. **Cost Observability when Unlimited:**
   - **Scenario:** Does setting a cost budget to unlimited disable telemetry collection?
   - **Verification:** I confirmed that [emitLoopBudget](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/loop.go#L157-L172) in [internal/driver/loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/loop.go) now computes `loopCostUSD` unconditionally, writing the burn to the event log even when `MaxCostUSD == 0` (enforcement is disabled). This ensures continuous visibility of spend on the TUI, verified by [TestEmitLoopBudgetReportsCostWhenUnlimited](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/loop_budget_test.go#L102-L125).

## Findings

Zero findings. The fixes (F-T2-1, F-T2-2, F-T2-3) are clean, robust, and correctly implement all safety invariants and escape hatches.

## Open questions

- The scoping deviation where `continue --auto` does not support transient command-line overrides for loop budgets remains acceptable and consistent with the consensus design.
