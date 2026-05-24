---
agent: gemini
idea: continuous-run-tui
review-round: 1
date: 2026-05-24
reviewed-commit: 1b08f410c7b9cf046b487c0f4d976eebe0793724
---

## Summary

The implementation successfully adds the core `runplan` logic and integrates it into `runstate`, the CLI `continue` command, and the TUI. The durable planning surface is a significant step toward managing multi-agent workflows continuously. However, the current implementation is heavily focused on `round-01` and lacks the robust tracking and planning required for multi-round workflows. Additionally, the `run.json` manifest extensions are not yet fully utilized as a source of truth during state reconstruction.

## Findings

### [CRITICAL]
None.

### [MAJOR]
1. **Hardcoded Round-01 Logic in Planner**: `internal/runplan/runplan.go` (L77) hardcodes `round := "round-01"` when checking for missing artifacts. This prevents the planner from correctly identifying missing or failed artifacts in `round-02` and subsequent rounds, which is a core requirement for Phase 2+ cross-review workflows.
2. **Brittle CLI Command Derivation**: `internal/app/app.go`'s `actionCommand` (L1439-1451) hardcodes `--round 1` and `--by codex` for `KindDraftConsensus` and `KindFinalize`. This ignores the actual round state and the protocol's drafter rules, which specify that the first submitter or initiator should draft. Hardcoding `codex` will cause failures in environments where `codex` is not the preferred or available drafter.

### [MINOR]
1. **Redundant `NextAction` Structs**: The `NextAction` struct is defined identically in both `internal/runplan` and `internal/runmanifest`. This redundancy increases the risk of drift as the planner evolves. These should be unified in a shared package.
2. **Underutilized `run.json` Manifest**: While `internal/runmanifest` was extended with continuation snapshot fields (Phase, IdeaStatus, CurrentRound, etc.), `internal/runstate.LoadRun` (L86) continues to rely solely on `events.jsonl`. The "source-of-truth" order in `FINAL.md` lists `run.json` as a fallback, but the current implementation doesn't yet attempt to read it to populate missing state or cached actions.
3. **Incomplete Participant Inference for Commands**: `actionCommand` does not use the `run.Participants` list to suggest a valid participant for drafting or finalization, relying instead on a hardcoded string.

### [NIT]
1. **Constant Use in `runplan.go`**: `riskFromQuestion` uses `hitl.RiskLow` and `hitl.RiskHigh` (L176) instead of the local `Risk*` constants, which is slightly inconsistent with the rest of the file.
2. **Event Summary Formatting**: `internal/runstate/runstate.go` uses hardcoded strings like `"idea=%s"` in `SummarizeEvent`. These would be better as constants if they are intended to be part of a stable developer surface.
3. **TUI Action Visibility**: The TUI limits visible actions to 4. While the "more action(s)" indicator is present, this may be too restrictive for complex runs where multiple high-risk actions might be pending.
