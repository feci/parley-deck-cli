---
agent: claude
idea: continuous-run-tui
review-round: 1
date: 2026-05-24
reviewed-commit: 1b08f410c7b9cf046b487c0f4d976eebe0793724
---

## Summary

This commit implements the first slice of the continuation planner as specified in FINAL.md: a durable `runplan.Plan` function that derives next actions from run state, conservative `run.json` continuation snapshot fields, a read-only `parley continue` CLI command, and planner output surfaced in the TUI actions panel. The implementation is well-structured with good test coverage for each planner action kind. The code aligns closely with the spec's planner model, source-of-truth order (canonical files take precedence), and non-goals (no daemon, no silent advancement). A few issues merit attention.

## Findings

### [MAJOR] Duplicate NextAction struct creates a maintenance and serialization divergence risk

`runplan.NextAction` and `runmanifest.NextAction` are structurally identical structs defined in separate packages. When one is updated, the other must be kept in sync manually. `RunSummary.NextActions` uses `runplan.NextAction`, while `Manifest.NextActions` uses `runmanifest.NextAction` — callers converting between the two must do field-by-field copies with no compile-time safety that fields match.

**Why it matters:** A field added to one but not the other will silently drop data during the `run.json` write or `continue` output path. This contradicts the spec's intent that `run.json` snapshot fields and planner output share the same action model.

**Suggested fix:** Define `NextAction` once (in `runplan` or a shared `runtypes` package) and import it from both `runmanifest` and `runstate`. Alternatively, have `runmanifest` embed or alias `runplan.NextAction`.

### [MAJOR] Hardcoded `--by codex` in generated commands assumes a specific agent identity

`actionCommand` (app.go:1459, 1469) hardcodes `--by codex` for `draft-consensus` and `finalize` commands. If the user or a different agent runs `parley continue`, the suggested command attributes the consensus/finalization to "codex" regardless of who is actually performing the action.

**Why it matters:** This violates FINAL.md's principle that only the owning agent may produce artifacts. The generated command could mislead users into creating artifacts attributed to the wrong agent.

**Suggested fix:** Either omit `--by` from the generated command (let the CLI prompt or infer it), or derive the agent identity from the run's active participant context / the current user session.

### [MINOR] Planner only inspects `round-01` — multi-round continuation is not handled

`runplan.Plan` hardcodes `round := "round-01"` (runplan.go:82). FINAL.md's planner model and consensus handling imply the planner should work across rounds. If a run has progressed to round-02 (e.g., after a blocked consensus triggers another round), the planner will incorrectly check for round-01 artifacts and may suggest unnecessary retries or premature consensus drafting.

**Why it matters:** While FINAL.md scopes the first slice narrowly, the planner's round detection should at minimum read the current round from the manifest or input rather than assuming round-01.

**Suggested fix:** Use `input.RoundStatus` or add a `CurrentRound` field to `runplan.Input` (already available in the manifest as `CurrentRound`) and derive the round directory from it.

### [MINOR] `parley continue` does not consult append-only run events (source-of-truth item 2)

FINAL.md specifies the source-of-truth order as: (1) canonical files, (2) run events, (3) run.json snapshot, (4) session index. The planner receives agent state and questions derived from events indirectly via `RunSummary`, but does not itself read `events.jsonl`. This means certain event-only signals (e.g., a recent retry event that hasn't yet produced an artifact) cannot influence planning.

**Why it matters:** For this first slice (read-only), the impact is low. But as the planner gains write capabilities, missing event context could cause duplicate retries or stale action suggestions.

**Suggested fix:** Document this as a known limitation for slice 1, or pass a filtered event summary into `runplan.Input` for richer planning decisions.

### [NIT] `RunSummary` alignment uses mixed tabs/spaces for struct tag column

The diff shows the `RunSummary` struct fields were reformatted with varying amounts of whitespace before the struct tags (e.g., `Terminal` has many more spaces than `RunID`). While `gofmt` normalizes this, the diff suggests manual editing may have introduced inconsistency that could produce noisy future diffs if the file is reformatted.

## Open questions

1. Is there a plan to unify `runmanifest.NextAction` and `runplan.NextAction` in a follow-up slice, or should this be addressed before merging?
2. Should `actionCommand` support a `--by` flag derived from the current user/agent context, or is the intent that users always edit the suggested command before running it?
3. The planner is invoked on every `LoadRunAt` call. For workspaces with many runs, is there a performance concern with calling `consensus.Status` and `os.Stat` per-participant for every loaded run?
