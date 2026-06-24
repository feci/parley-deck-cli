---
idea: loop-budgets
drafter: claude-1
implementer: claude-1
date: 2026-06-24
status: final
spawned_from: loop-engineering-research
---

## Purpose

Implementation spec for Tier 2 (the safety floor). Adds explicit, escalate-on-breach
loop ceilings to the driver and fixes the non-solo violation in the pipeline watcher.

## LE-5 — Unified loop-budget contract

**Files:** `internal/driver/driver.go` (Config + New defaults), `internal/driver/loop.go`
(enforcement + event), `internal/config/runtime.go` (`[defaults.loop]` + seed template),
`internal/app/app.go` (apply defaults at the 3 driver-construction sites + `run` flags);
`COOPERATION.md` §4/§9.0/§12.

- `driver.Config` gains:
  - `MaxDriverSteps int` — total `Advance` iterations that produce progress (non-await);
    `0` = unlimited.
  - `MaxWallClock time.Duration` — total run wall-clock budget across all ticks (distinct
    from the existing per-tick `roundDeadline`, which resets on progress); `0` = unlimited.
  - `MaxCostUSD float64` — total external-backend cost budget; `0` = unlimited.
- `loop.go Run` enforces them: a step counter (incremented on each progress action), a
  `start := time.Now()` wall-clock origin, and a running cost total. **Before each
  `Advance`**, if any non-zero ceiling is exceeded → `escalate(c, "loop-budget", …)`
  (durable inbox note, same path as the deadline escalation) and return — **never
  Complete**. A `loop.budget` event (steps, elapsed_ms, cost_usd, the ceilings) is
  appended on each progress step so the TUI/state can show budget burn.
- **Semantics:** `0` everywhere = unlimited (backward-compatible: existing decks/tests are
  unchanged). Enforcement turns on for users whose `~/.parley [defaults.loop]` seeds
  values — `parley init` seeds generous safety-net defaults (steps 200, wall-clock 2h,
  cost 0/unlimited). `run --max-driver-steps N` / `--max-wall-clock D` override per run.
- Defaults precedence: built-in (unlimited) → `~/.parley [defaults.loop]` (seeded) →
  deck `agents.toml` → `run` CLI flags.

**Acceptance:** with `MaxDriverSteps: 1`, a multi-step idea escalates with a loop-budget
inbox note instead of completing; with `MaxWallClock` in the past, the loop escalates on
the next tick; `0` ceilings preserve today's behavior exactly (existing driver tests green).

## LE-6 — Best-effort cost telemetry

**Files:** `internal/driver/loop.go` (sum cost from `agent.usage` events), docs.

- The Run loop sums any `agent.usage` events' `cost_usd` from the store into the running
  total used by the `MaxCostUSD` ceiling and the `loop.budget` event.
- **Honest limit:** the CLI runners do not yet emit `agent.usage` (the agents' cost output
  is agent-specific and not parsed), so the summed cost is `0` in practice and the
  `MaxCostUSD` ceiling is effectively inert until a runner emits usage. The contract +
  plumbing land here; **emitting `agent.usage` from the runner is a deferred follow-up**
  (`agent-usage-telemetry`). Documented so the knob is not mistaken for active enforcement.

## LE-10 — Fix `openRemediationIdea` non-solo violation

**Files:** `internal/app/pipeline_cmd.go` (`openRemediationIdea`); `COOPERATION.md` §12.11.

- The watcher currently writes `participants: []` + `status: round-01` — an "active" idea
  with no quorum, violating the non-solo Phase-0 invariant.
- Fix: write `status: candidate` (a non-active status the driver does not pick up) and omit
  the empty-quorum claim; add a "## Promotion" note that a human or the pipeline manifest
  must set `participants:` and flip `status: round-01` before deliberation starts.
- §12.11: watcher-created remediation ideas are **candidates** until they satisfy Phase-0
  quorum.

**Acceptance:** an auto-opened remediation idea has `status: candidate` (not `round-01`)
and no empty `participants: []` round-01 claim; a test asserts the generated prompt.

## Conditional rigor / non-goals
- Budgets are a SAFETY NET, not a tight leash: seeded defaults are generous; `0` = off.
- This idea does NOT add the outer-loop trigger (that is Tier 4 / LE-9) — it only bounds
  the existing auto-drive loop and fixes the existing watcher's quorum bug.

## Tests
- `driver`: `MaxDriverSteps`/`MaxWallClock` breach escalates (fake ImplOps), `0` is
  unlimited, the `loop.budget` event is emitted.
- `config`: `[defaults.loop]` parses + merges into `CentralDefaults`; seed template
  contains the block.
- `app`: `openRemediationIdea` emits `status: candidate` and no `participants: []`.
- Drift guard stays green (edit BOTH COOPERATION.md copies).
