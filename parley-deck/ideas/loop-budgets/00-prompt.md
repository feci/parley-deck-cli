---
idea: loop-budgets
author: claude-1
created: 2026-06-24
participants: [claude-1, codex-1, hermes-1, antigravity-1]
implementer: claude-1
status: final
spawned_from: loop-engineering-research
checks: go build ./... && go test ./internal/driver/... ./internal/runner/... ./internal/app/... ./internal/config/... ./internal/protocol/...
---

## Problem / idea

Tier 2 of the `loop-engineering-research` backlog: **stopping conditions / the safety
floor**. Loop engineering mandates that an automated loop carry explicit ceilings —
max iterations, max wall-clock, max cost — and that **hitting a ceiling escalates, never
marks complete**. Parley's driver today has scattered, partial caps (`loop.go`'s 30-min
per-tick `roundDeadline`, `MaxRounds`, `MaxFixupCycles`, `pipeline auto`'s `maxCycles`)
but no total driver-step ceiling, no total wall-clock budget, and no cost budget.

Bundled with **LE-10**: `openRemediationIdea` opens watcher ideas with `participants: []`
+ `status: round-01`, violating the non-solo Phase-0 invariant — fix before any trigger
expansion (a prerequisite the audit flagged).

Implemented by claude-1 (Phase 5); reviewed by codex-1, hermes-1, antigravity-1
(Phase 6, refutation). Design = `FINAL.md`.

## Scope (LE-5, LE-6, LE-10)
- **LE-5** Unified loop-budget contract: `MaxDriverSteps`, `MaxWallClock`, `MaxCostUSD`
  on the driver, enforced in the Run loop — breach escalates (inbox note) and halts.
  Numbers seeded from `~/.parley [defaults.loop]`; `run` CLI flags override. A
  `loop.budget` event is emitted each progress step.
- **LE-6** Best-effort cost: the `MaxCostUSD` field + budget event + cost summed from
  `agent.usage` events when present (telemetry-gated — see FINAL for the honest limit).
- **LE-10** `openRemediationIdea` writes a non-active candidate (`status: candidate`, no
  empty participant quorum claim) until a human/quorum promotes it to `round-01`.
