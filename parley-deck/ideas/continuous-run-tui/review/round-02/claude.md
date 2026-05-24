---
agent: claude
idea: continuous-run-tui
review-round: 2
date: 2026-05-24
reviewed-commit: uncommitted-fixup-cycle-1-after-1b08f410c7b9cf046b487c0f4d976eebe0793724
---

## Summary

Fix-up cycle 1 resolves all four agreed round-01 findings. The duplicate `NextAction` type is unified via a new `internal/runaction` package. The planner now reads `CurrentRound` from input instead of hardcoding `round-01`. Generated continuation commands no longer embed `--by codex` or a hardcoded round number. Each fix has targeted test coverage. No agreed-fix blockers remain.

## Findings

### Finding 1: Duplicate `NextAction` definitions — RESOLVED

`internal/runaction/action.go` now holds the single canonical `NextAction` struct and kind/risk constants. Both `runplan` (line 29: `type NextAction = runaction.NextAction`) and `runmanifest` (line 57: `type NextAction = runaction.NextAction`) re-export it as a type alias. The kind constants are similarly aliased in `runplan` (lines 17-27). This eliminates divergence risk while preserving backward-compatible package-level references.

### Finding 2: Hardcoded `round-01` for planner retry/draft actions — RESOLVED

`runplan.Input` gained a `CurrentRound` field (line 40). `Plan` calls `currentRound(input)` (line 80) which reads the field and falls back to `"round-01"` only when empty (lines 195-200). The round flows from manifest via `applyManifestDefaults` (runstate.go:195) into `summary.CurrentRound`, with a secondary inference from agent artifact paths via `inferCurrentRound` (runstate.go:204-214), and is passed into `runplan.Input.CurrentRound` (runstate.go:169).

### Finding 3: Hardcoded `--by codex` / `--round 1` in generated commands — RESOLVED

`actionCommand` (app.go:1445-1482) no longer emits `--by` for `KindDraftConsensus` or `KindFinalize`. The `--round` flag now uses `roundNumber(action.Round)` (app.go:1459), which extracts the numeric suffix from the action's round field (e.g., `"round-02"` becomes `"2"`). The `roundNumber` helper (app.go:1484-1493) is clean and handles edge cases (empty input, leading zeros).

### Finding 4: Tests for the fixes — RESOLVED

Three targeted tests cover the fixes:

- `TestPlanUsesCurrentRoundForMissingArtifacts` (runplan_test.go:43-59): verifies the planner targets the correct round directory and produces actions with matching `Round` and `ArtifactPath` when `CurrentRound: "round-02"` is set.
- `TestLoadRunUsesManifestCurrentRoundForPlanner` (runstate_test.go:92-131): verifies the full integration path from manifest `CurrentRound` through runstate loading into planner output with correct round-02 artifact paths.
- `TestActionCommandUsesActionRoundAndAvoidsHardcodedAgent` (app_test.go:1103-1124): verifies `actionCommand` generates `parley consensus draft --round 2 sample` from a round-02 action, and asserts the output does not contain `"codex"`.

### [NIT] Re-exported constants in runplan could be removed if no external consumers exist

`runplan` lines 17-27 re-export all six kind constants and three risk constants from `runaction` as package-level `const` aliases. If no code outside the `internal/` tree references `runplan.KindAnswerQuestion` etc., these could be replaced by direct `runaction.Kind*` references at call sites. However, since both `runstate` and `app` already import `runplan`, the aliases avoid churn and maintain a consistent import surface. This is cosmetic, not blocking.
