---
agent: antigravity-1
idea: loop-budgets
review-round: 1
date: 2026-06-24
---

## Summary

This review analyzes the `loop-budgets` design and implementation (Phase 6, refutation mode) through a strict risk, safety, and guardrails lens. The focus is to identify edge cases, logic flaws, and gaps where the auto-drive loop could escape enforcement, spin unbounded, or fail unexpectedly.

The implementation introduces vital guardrails to prevent infinite deliberation or execution loops. However, the review identified one **MAJOR** defaults precedence gap that prevents users from disabling safety budgets once a global default is seeded, alongside a few minor efficiency and robustness findings.

## Refutation attempts

I systematically attempted to break the budget controls and verify the correctness of the invariants:

1. **Escaping Budget Ceilings or Spinning Unbounded:**
   - I checked if any progress action could run without incrementing the step budget. All progress-generating `Action` types return `true` in `isProgressAction` and increment the step count.
   - I evaluated whether the driver could spin unbounded during await phases (e.g. `ActionAwait`). The loop is safely bounded by the 30-minute `roundDeadline`, which halts the driver loop and escalates to a user inbox note if reached.
   - I checked if clock shifts or temporal changes could bypass `MaxWallClock`. Go's `time.Since` utilizes a monotonic clock, ensuring immunity to system time changes.

2. **Verifying LE-10 (Candidate Status and Quorum Safeguard):**
   - I verified whether watcher-opened ideas with `status: candidate` could be active or auto-driven.
   - Because the watcher omits the `participants` list, `ReadWorkspaceStatus` parses `participants` as `nil`.
   - The CLI `ResolveRun` logic will reject any attempt to resume or continue the idea slug because no run directory exists under `runs/`, returning an error indicating the idea has no runs yet.
   - In a hypothetical case where a run directory is manually forced with empty participants, the driver's `roundComplete` check immediately detects that the participant list is empty and returns `(false, nil)`, placing the driver in an `ActionAwait` state rather than driving or completing it.

3. **Analyzing Budget Defaults Precedence:**
   - I traced the precedence chain: built-in (`0`) $\rightarrow$ `~/.parley` seed (`200` steps, `2h` wall-clock) $\rightarrow$ deck config $\rightarrow$ run flags.
   - A major merging gap was identified where setting a limit to `0` (unlimited) at a higher precedence layer is ignored because the parser cannot distinguish "not set" from "set to 0".

---

## Findings

### MAJOR: Defaults Precedence Gap — Impossible to Disable Budgets (Set to Unlimited) from Higher Precedence Layers
* **What is wrong:**
  Once a user initializes parley and seeds a global default in `~/.parley/agents.toml` (which sets `max_driver_steps = 200` and `max_wall_clock_ms = 7200000`), it is impossible for a project deck (`parley-deck/agents.toml`) or a CLI invocation (using `--max-driver-steps 0`) to disable loop budgets (set them to `0` / unlimited).
* **Why:**
  In `internal/config/runtime.go`, `mergeDefaults` only overrides values if the higher precedence layer's value is strictly greater than zero:
  ```go
  if gd.Loop.MaxDriverSteps > 0 {
      out.MaxDriverSteps = gd.Loop.MaxDriverSteps
  }
  ```
  Similarly, `internal/app/app.go:loopBudget` checks:
  ```go
  if stepsOverride > 0 {
      steps = stepsOverride
  }
  ```
  Because Go unmarshals missing TOML keys as `0` by default on integer/float fields, `0` is used to mean "not specified." Consequently, an explicit setting of `0` (unlimited) is treated as "not specified" and falls back to the seeded global default.
* **Concrete Fix:**
  Define the fields in `loopBlock` (within `internal/config/runtime.go`) as pointer types:
  ```go
  type loopBlock struct {
      MaxDriverSteps *int     `toml:"max_driver_steps"`
      MaxWallClockMS *int     `toml:"max_wall_clock_ms"`
      MaxCostUSD     *float64 `toml:"max_cost_usd"`
  }
  ```
  Then, check for non-nil pointers in `mergeDefaults`:
  ```go
  if gd.Loop != nil {
      if gd.Loop.MaxDriverSteps != nil {
          out.MaxDriverSteps = *gd.Loop.MaxDriverSteps
      }
      if gd.Loop.MaxWallClockMS != nil {
          out.MaxWallClockMS = *gd.Loop.MaxWallClockMS
      }
      if gd.Loop.MaxCostUSD != nil {
          out.MaxCostUSD = *gd.Loop.MaxCostUSD
      }
  }
  ```
  For the CLI flags, allow `-1` to explicitly signify "unlimited" or check if the flag was explicitly set by the user using `flag.Visit`.

### MINOR: Resource Exhaustion / Sleeping in Await Loop when Roster is Empty
* **What is wrong:**
  If a run is initialized or continued on an idea that has no participants (e.g. `len(d.cfg.Participants) == 0`), the driver's round verification logic in `internal/driver/driver.go:roundComplete` returns `(false, nil)`. This causes the driver to yield `ActionAwait`, leading it to spin in a `time.Sleep(2 * time.Second)` loop for 30 minutes until the `roundDeadline` trips.
* **Why:**
  `roundComplete` assumes that an empty participant list is simply a pending/incomplete state rather than a protocol validation failure.
* **Concrete Fix:**
  In `internal/driver/driver.go:roundComplete`, validate the participant list size and return an error immediately if it is empty:
  ```go
  if len(d.cfg.Participants) == 0 {
      return false, fmt.Errorf("cannot drive round: no participants configured for idea %s", d.cfg.IdeaSlug)
  }
  ```

### MINOR: Inefficient Event Log Reloading on Await Ticks
* **What is wrong:**
  If a cost budget is configured (`MaxCostUSD > 0`), the driver checks the budget before every single advance tick in the `Run` loop. This calls `loopCostUSD()` which reads and parses the entire `events.jsonl` file from disk on every single loop iteration. During `ActionAwait` periods, the loop sleeps for 2 seconds, wakes up, and immediately re-reads and re-parses the file, even if no progress has occurred. On long runs with large event logs, this causes significant, unnecessary disk I/O and CPU overhead.
* **Why:**
  The driver does not cache the cost value or events log, nor does it check if any new progress steps or agent runs occurred before reloading.
* **Concrete Fix:**
  Cache the calculated cost, or only reload `events.jsonl` when `Advance` indicates that a progress action has completed or when a new agent has run, rather than loading the file on every tick (including pure sleep ticks).

### NIT: Best-effort Cost Telemetry Parsing Robustness
* **What is wrong:**
  In `loopCostUSD()`, the agent usage cost is read by checking `e.Data["cost_usd"].(float64)`. If a future runner writes `cost_usd` as a string (e.g. `"1.50"` or `"0.05"`), or if an integer value is unmarshaled as something other than `float64` (e.g. if the JSON parser changes), the cost will be silently ignored (added as `0`) rather than parsed.
* **Why:**
  The driver relies on a strict type assertion `.(float64)` without attempting to parse strings or check other numeric types.
* **Concrete Fix:**
  Implement a robust numeric parser helper that checks for both string formats and float formats (e.g., using `strconv.ParseFloat` or a type switch).

---

## Open questions

* **Candidate Promotion UX:** While `status: candidate` successfully isolates watcher-created remediation ideas from being auto-driven, there is currently no CLI subcommand to initialize a run for an existing idea folder (e.g., `parley run-idea <slug>`). Since `parley run` always creates a new timestamped idea directory from a new task string, how are operators expected to promote and execute these candidate ideas in-place once a quorum is staffed?
