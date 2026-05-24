---
agent: codex
idea: continuous-run-tui
round: 1
date: 2026-05-23
---

## Summary

The next slice should make Parley continuable at the workflow-step level, not by pretending it can always reattach to dead OS processes. The TUI should become a workbench over many runs: select a run, see its next safe action, answer questions, retry only missing/failed participant work, request signoffs, or start another idea while existing work keeps progressing inside the current TUI process. Durable recovery after closing the TUI should be implemented as reconciliation plus step-level retry/resume, with native agent session continuation treated as best-effort.

## Proposed approach

### 1. Add a run action planner

Create a small package, likely `internal/runplan`, that turns durable state into explicit next actions.

Inputs:

- `parley-deck/runs/<run-id>/events.jsonl`
- `parley-deck/runs/<run-id>/run.json`
- `parley-deck/ideas/<slug>/00-prompt.md`
- round/review files under `parley-deck/ideas/<slug>/`
- consensus and review consensus files
- HITL questions under the run directory
- configured agent runtime capabilities

Output:

```go
type NextAction struct {
    ID          string
    Kind        string // answer-question, retry-agent, resume-agent, draft-consensus, request-signoffs, finalize, open-next-round, inspect
    RunID       string
    IdeaSlug    string
    Phase       string
    Round       string
    AgentID     string
    Artifact    string
    Risk        string // low, normal, high
    RequiresYes bool
    Summary     string
}
```

The key rule: the planner reconciles files first. If an expected artifact exists and validates, it is done even if the old run event says the agent failed. If an artifact is missing or malformed, only the owning agent can retry or repair it.

### 2. Extend the run manifest without replacing events

Keep `events.jsonl` as the append-only audit stream and `run.json` as a current snapshot.

Extend `run.json` with conservative fields:

- `status`: `running | waiting | action_required | incomplete | failed | completed | cancelled | stale`
- `phase`: `round | consensus | finalization | implementation | review | fixup`
- `idea_status`
- `current_round`
- `active_steps`
- `last_action_at`
- `next_actions`

Do not make `run.json` the only truth. It is a cache/snapshot that can be regenerated from files and events.

### 3. Introduce step attempts

Add a per-run attempt ledger:

```text
parley-deck/runs/<run-id>/
  attempts/
    <step-id>.<attempt-n>.json
  input-packs/
    <step-id>.prompt.txt
```

Each attempt records:

- step ID such as `round-01.claude` or `review-round-02.gemini`
- command path and configured launch mode
- prompt hash and input-pack path
- expected artifact path
- start/end timestamps
- stdout/stderr paths
- exit state and validation result
- optional native resume handle if the adapter can capture one

This lets the CLI continue a run by retrying only `round-01.gemini`, not by creating a new idea or rerunning all participants.

### 4. Generalize the runner from "round one" to "step execution"

The current `runner.RunRoundOne` is hard-coded around round-01 prompt generation and selected participants. Keep it working, but build a second layer:

```go
type Step struct {
    ID           string
    Kind         string // round, review, signoff
    AgentID      string
    Prompt       string
    OutputPath   string
    Validate     func(path string) error
}

func RunStepsAsync(ctx context.Context, opts StepRunOptions) *Handle
```

Then round-01 becomes one producer of steps. Later the same executor can drive:

- retry missing round files;
- review round files;
- consensus signoff requests;
- mechanical handoff validation.

### 5. TUI workbench model

Upgrade `parley tui` from mostly observe/start to a run workbench.

Layout:

- Left pane: all current-workspace runs and ideas, sorted by attention then recency.
- Main pane: selected run timeline and latest events.
- Right pane: selected run details, selected agent/step, logs, artifact path.
- Bottom pane: action queue for selected run plus global pending HITL questions.

Key bindings:

- `N`: start new idea/run.
- `enter`: open action palette for selected run.
- `c`: continue the top recommended action for selected run.
- `r`: retry selected failed/missing step.
- `a`: answer selected HITL question.
- `s`: switch run/idea focus.
- `p`: pause/cancel a live TUI-owned attempt after confirmation.
- `R`: refresh/reconcile state.

The selected run remains stable by `RunID`, not by list index. Refreshes must not move the user's focus if attention sorting changes.

### 6. Parallel runs inside one TUI process

Add a TUI-local supervisor:

```go
type Supervisor struct {
    handles map[string]*runner.Handle
    cancels map[string]context.CancelFunc
    limit int
}
```

It should:

- allow multiple TUI-started runs at once;
- cap concurrent agent processes, defaulting to a conservative value such as 2 runs or a configurable agent-process limit;
- reject starting the same run twice;
- expose per-run live/cancel state to the TUI;
- update `run.json` and session index when terminal events arrive.

If the TUI exits, v1 may cancel owned live attempts, but the TUI must clearly say so. Continuation after reopening is then step-level recovery, not process reattachment. A later daemon/supervisor can change this promise.

### 7. CLI commands

Keep `parley resume` as "reopen or inspect durable state"; add an explicit mutating command:

```text
parley continue [--dir DIR] [--yes] [--action ACTION_ID] RUN_OR_IDEA
```

Behavior:

- no `--action`: print planner result and the recommended next action;
- with `--yes` and a safe default action: execute it;
- with `--action`: execute the selected action;
- never silently edits another agent's canonical artifact.

The TUI should call the same app-layer action executor, not duplicate continuation logic.

### 8. Implementation slices

Slice 1: planner and snapshots

- Add `internal/runplan`.
- Extend `runmanifest.Manifest` conservatively.
- Add tests for next-action derivation: open question, missing round file, malformed artifact, all round files ready for consensus, missing signoffs, blocked consensus.

Slice 2: retry one missing participant step

- Add input-pack persistence and attempt metadata for round-01.
- Add `parley continue --action retry-agent`.
- Reuse existing runner validation.

Slice 3: TUI actions

- Add action pane and action palette.
- Wire `continue`, `retry`, and `answer`.
- Preserve selection by run ID/action ID across refresh.

Slice 4: bounded parallel supervisor

- Move TUI-started run handles into a `Supervisor`.
- Add concurrency limits and duplicate-run protection.
- Add cancellation confirmation.

Slice 5: native agent resume handles

- Extend agent discovery with `resume_capability`.
- Capture handles when reliable.
- Prefer native resume, fall back to fresh retry with cached prompt.

Slice 6: optional detach daemon

- Only after the previous slices are reliable, design a real background supervisor if the user wants runs to survive TUI exit.

## Concerns / open questions

- The word "resume" is overloaded. I recommend using `resume` for read-only state restoration and `continue` for mutating progression.
- Native agent resume should not be MVP-critical. The useful guarantee is "do not lose Parley workflow progress"; exact model-session continuity is best-effort.
- Parallelism needs a cap. Starting four ideas with four participants each can easily spawn too many hosted CLIs.
- GitHub PR transport adds a later question: should `continue` fetch remote PR/MR state before deciding phase advancement? I would defer that to a `--refresh-transport` flag.

## Risks

- If the TUI starts mutating phases too aggressively, it can bypass Parley consensus rules. Keep every action explicit and render the artifact it will create.
- If run state is split between `run.json`, events, and artifacts without a clear precedence rule, recovery will become inconsistent. Precedence should be: canonical artifacts, then events, then snapshots.
- If prompts are regenerated on retry without recording prompt hashes, retries may silently differ from the original task. Persist input packs from the first attempt.
- If the first version promises detached background execution, users will lose work. Do not promise this until a real daemon/supervisor exists.
