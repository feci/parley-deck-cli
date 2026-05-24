---
agent: hermes
idea: continuous-run-tui
review-round: 1
date: 2026-05-24
reviewed-commit: 1b08f410c7b9cf046b487c0f4d976eebe0793724
---

## Summary

The implementation adds the run action planner (`internal/runplan`), conservative continuation snapshot fields in run manifests and state, the `parley continue` read-only CLI command, and TUI action panel integration. It follows the FINAL.md scope for slice 1, the planner model action kinds, and source-of-truth precedence favoring canonical files. Changes are limited to the approved directories and docs. Tests cover the planner cases. No silent advancement or daemon work introduced.

## Findings

### [MAJOR] Duplicate NextAction definitions between runplan and runmanifest

`internal/runplan/runplan.go:28` and `internal/runmanifest/manifest.go:33` each define an identical `NextAction` struct. Runstate copies between them (runstate.go:151). This creates future divergence risk for JSON serialization and field additions.

### [MINOR] Round is hardcoded to round-01 in planner

`internal/runplan/runplan.go:90` sets `round := "round-01"`. While the current slice targets round-01, the planner input already carries `RoundStatus`; using it would make the logic future-proof for later rounds without violating the non-goals.

### [MINOR] Planner does not directly consume events.jsonl

FINAL.md lists run events as source-of-truth #2. Planning logic in runplan.go derives state only via RunSummary/AgentState passed in; no direct `events.jsonl` read occurs inside Plan(). Acceptable for read-only slice but noted for later evolution.

### [NIT] Inconsistent struct tag alignment in runstate.go

The added NextActions field and surrounding RunSummary fields show uneven whitespace before json tags (runstate.go:80-90 range). gofmt would clean this; current state adds minor noise to the diff.

No CRITICAL findings. No MAJOR findings beyond the duplication noted above. No other MINOR or NIT issues found within scope.