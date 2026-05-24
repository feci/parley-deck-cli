---
idea: continuous-run-tui
cycle: 2
drafted-by: codex
date: 2026-05-24
reviewed-commit: uncommitted-fixup-cycle-1-after-1b08f410c7b9cf046b487c0f4d976eebe0793724
---

## Agreed fixes

None remaining after fix-up cycle 1. Round 02 reviewers agree that the round-01 agreed findings were resolved:

- `NextAction` has a single canonical definition in `internal/runaction`.
- Planner actions target `CurrentRound` instead of hardcoded `round-01`.
- Generated continuation commands no longer hardcode `--by codex` or `--round 1`.
- Targeted tests cover the fix-up paths.

## Deferred follow-ups

- Direct `events.jsonl` consumption inside `runplan.Plan` remains deferred. Current slice still derives event state through `runstate.RunSummary`, which is acceptable for the read-only planner slice.
- Broader detached/supervisor continuation remains outside this slice, as specified in `FINAL.md`.

## Dismissed findings

- Claude round-02 NIT about replacing `runplan` re-exported constants with direct `runaction` imports is dismissed as cosmetic. The aliases preserve the existing package-level call surface and avoid unrelated churn.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-24
Status: ✅ ACCEPT
Notes: Fix-up cycle 1 resolves the agreed review findings. No agreed fixes remain for this slice.

### Signoff: claude — 2026-05-24
Status: ✅ ACCEPT
Notes: All four round-01 findings (duplicate NextAction, hardcoded agent identity, hardcoded round-01, missing event consumption) are resolved or appropriately deferred. The new `runaction` package cleanly unifies the type, `CurrentRound` flows correctly from manifest through planner, and generated commands no longer embed agent-specific defaults. Test coverage targets each fix path. No blockers remain for this slice.

### Signoff: gemini — 2026-05-24
Status: ✅ ACCEPT
Notes: Agreed fixes from round 01 are resolved. The planning logic is now round-aware, and the `NextAction` model is unified in `internal/runaction`. The deferral of direct `events.jsonl` consumption in the planner is acceptable for this slice.

### Signoff: hermes — 2026-05-24
Status: ✅ ACCEPT
Notes: All round-01 findings resolved in fix-up cycle 1. Long-running recovery and stale-process concerns addressed via round-aware planner and unified NextAction; deferred items are appropriate per scope.
