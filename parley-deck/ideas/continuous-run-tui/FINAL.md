---
idea: continuous-run-tui
finalized-by: codex
date: 2026-05-23
status: final
exception: user-authorized-solo-continuation-after-peer-agent-blockers
---

# FINAL: Continuous Run TUI

## Goal

Make `parley-deck-cli` usable continuously across existing runs and ideas. A user should be able to inspect a concrete run, understand the next safe action, switch to another idea while work is pending, intervene when needed, and later continue without creating a duplicate idea or rerunning completed work.

## Scope for the first implementation

The first implementation slice adds the durable planning surface that both CLI and TUI can use:

1. Add a run action planner that derives next actions from canonical files, run events, HITL questions, and run metadata.
2. Extend `run.json` with conservative continuation snapshot fields.
3. Add `parley continue` as the public CLI entrypoint. The initial command may be read-only unless an action is explicitly safe to execute.
4. Surface planner output in workspace/TUI state so the TUI can show what each selected run needs next.

## Planner model

The planner should produce stable action records:

- `answer-question`
- `retry-agent`
- `draft-consensus`
- `request-signoffs`
- `finalize`
- `inspect`

Each action includes an ID, run ID, idea slug, phase, round, optional agent ID, optional artifact path, risk level, and a short summary.

## Source-of-truth order

1. Canonical files under `parley-deck/ideas/<slug>/`.
2. Append-only run events under `parley-deck/runs/<run-id>/events.jsonl`.
3. `run.json` snapshot fields.
4. User-local session index.

If a valid expected artifact exists, it is complete even if an old event says the agent failed. If an artifact is missing or malformed, only the owning agent or an explicit user-approved exception may repair it.

## TUI direction

The workspace TUI should evolve into a workbench:

- left side: runs/ideas sorted by attention and recency;
- main area: selected run timeline and status;
- action area: planner-derived next actions;
- commands for inspect, continue, retry, answer, refresh, start new run, and switch focus.

The selected run must remain stable by run ID across refreshes.

## Parallelism direction

Use a TUI-local supervisor first:

- keep in-memory handles for TUI-started runs;
- reject duplicate starts for the same run;
- cap concurrent agent processes conservatively;
- allow switching focus without blocking on one run;
- clearly state that closing the TUI cancels TUI-owned live attempts until a daemon exists.

## Non-goals for the first implementation

- No detached daemon.
- No silent phase advancement.
- No broad automatic edits to another participant's artifact.
- No promise of native agent conversation continuity.

## Protocol note

This final plan was accepted under a user-authorized exception after configured peer agents were blocked by local auth/model/sandbox failures. The exception is recorded in the idea consensus and inbox.
