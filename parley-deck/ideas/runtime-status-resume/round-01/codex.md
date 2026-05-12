---
agent: codex
idea: runtime-status-resume
round: 1
date: 2026-05-12
---

## Summary

This slice should turn the existing event stream into a reusable runtime snapshot and then expose that snapshot through `parley status` and `parley resume`. The current code has the right raw materials: protocol frontmatter, `runs/<run-id>/events.jsonl`, HITL question files, agent logs, and `tui.ProjectEvents`. The important product boundary is to make resume truthful: it can restore the user's view and next recoverable action after a restart, but it should not claim to reattach to subprocesses that the current process no longer supervises.

## Proposed approach

Add a small runtime projection layer, likely under `internal/runtime` or `internal/store`, that scans `parley-deck/runs/*/events.jsonl` and derives immutable snapshots:

- run ID and run directory;
- task, mode, idea slug, participant list, and runtime matrix summary from `run.created`;
- latest event time and terminality from `round.completed`, `round.incomplete`, `run.failed`, and per-agent terminal events;
- per-agent state, artifact path, stdout/stderr paths, duration, and error/reason;
- open and answered HITL question counts by reading `runs/<run-id>/questions/*.md`;
- protocol artifact completion by reading the idea directory: `00-prompt.md`, `round-NN/<agent>.md`, `consensus.md`, `FINAL.md`, `IMPLEMENTATION.md`, and review artifacts.

Use this projection in both commands:

1. `parley status [--dir DIR] [--run RUN_ID] [--json]`
   - Keep the current workspace overview, but add a "Runs" section sorted newest first.
   - For each run, show `run`, `idea`, `phase/round`, `participants done/total`, `state`, pending questions, and next action.
   - Preserve plain terminal output as the default; add JSON only if it is cheap and directly useful for scripts/tests.

2. `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA`
   - Resolve `RUN_OR_IDEA` as an exact run ID first, then as an idea slug with the newest matching run.
   - If `--no-tui` is set, print the same snapshot plus next commands.
   - If TUI is requested and the run directory exists, reuse the live TUI in a read-only resume mode by feeding it the derived `LiveOptions`, a closed `Done` channel, and no cancel function.
   - For a terminal run, show the completed/incomplete result and the next protocol action.
   - For a non-terminal but stale run, show "not supervised by this process" and the recoverable commands rather than pretending execution can continue.

Do not add `state.json` as required source of truth in this slice. It is safe to design a later cache, but the first implementation should derive from events and protocol files every time so interrupted writes cannot corrupt resume behavior. If performance becomes an issue, `state.json` can be an explicitly disposable cache.

The current `tui.ProjectEvents` should either move into a shared package or become a thin wrapper over the new projection. Avoid duplicating event-state rules in `status`, `resume`, and TUI. The display layer can stay in `internal/app`; the state rules should be testable without Bubble Tea.

Suggested outputs:

```text
Transport: github-pr
Ideas:
  runtime-status-resume  status=round-01  participants=codex,claude,gemini,hermes

Runs:
  20260512T101500Z  idea=runtime-status-resume  state=incomplete  round=round-01  agents=3/4  questions=1 open
    next: parley answer 20260512T101500Z <question-id> <answer>
```

For `resume`:

```text
Run: 20260512T101500Z
Idea: runtime-status-resume
State: incomplete, not supervised by this process
Artifacts: round-01/codex.md ok, round-01/claude.md ok, round-01/gemini.md missing, round-01/hermes.md ok
Next: resolve missing participant output or rerun a new Parley Deck round manually.
```

Testing should use temporary workspaces with fixture event streams and question files:

- no runs;
- one completed run;
- one incomplete run with missing artifact;
- run ID resolution;
- idea slug resolution to newest run;
- malformed event line handling;
- open HITL question display;
- `resume --no-tui` output.

## Concerns / open questions

- The current live TUI receives an active `Done` channel from `runner.RunRoundOneAsync`. After process restart, there is no supervisor handle. Resume TUI should therefore be read-only unless a future supervisor/daemon writes durable process metadata.
- Existing older run events may not include full participant/runtime data. The projection should degrade gracefully: infer participants from the idea frontmatter when missing and show `unknown` rather than failing status entirely.
- `parley status` can become noisy if it prints all historical runs. A default of the newest 5 or 10 runs, with `--all` or `--run`, would keep it usable.
- If `resume` resolves an idea slug with multiple runs, it should choose the newest by default and print the selected run ID so the user can be explicit next time.

## Risks

- Duplicating event projection logic between TUI and CLI output will drift quickly. Shared projection is the main engineering safeguard.
- Over-promising "reattach" will create a false sense of recovery. The wording should be precise: resume state and UI now, execution reattachment only after a future supervised-run mechanism.
- Adding JSON output before the snapshot shape stabilizes could create premature API commitments. If included, label it as best-effort developer output or keep it minimal.
- Too much protocol validation in this slice could expand into full phase automation. Artifact completion should be read-only reporting, not automatic consensus or transport progression.
