---
agent: gemini
idea: runtime-status-resume
round: 1
date: 2026-05-12
---

## Summary
Propose an observability-first `parley status` and a recovery-focused `parley resume`. The design focuses on rebuilding state from `events.jsonl` and protocol frontmatter, enabling the TUI to re-hydrate from any previous run. `parley status` will provide a project-wide overview of ideas, active rounds, and pending HITL questions, while `parley resume` will allow continuing or retrying specific runs without losing context.

## Proposed approach

### 1. `parley status` Enhancement
Expand `parley status` to provide a high-signal overview of the workspace.
- **Global Overview**: Show active Transport and a count of pending HITL questions across all ideas.
- **Ideas Section**: For each idea in `ideas/`:
    - **Protocol State**: Current phase and round (from `00-prompt.md`).
    - **Latest Run**: Show the ID of the most recent run, its start time, and its outcome (e.g., `completed`, `failed`, `running/detached`).
    - **Agent Matrix**: A brief summary of agent status for the latest run (e.g., `codex:OK, gemini:FAILED, hermes:WAITING`).
- **Plain Output**: Ensure the output is clean, tabular, and suitable for standard terminals and CI.

### 2. `parley resume <run-id | idea-slug>`
Implement `resume` as the primary recovery mechanism.
- **Target Resolution**:
    - If a `run-id` is provided, use it directly.
    - If an `idea-slug` is provided, resolve it to the most recent run ID associated with that idea (via `run.created` events).
- **State Restoration**:
    - Load and project all events from `runs/<run-id>/events.jsonl`.
    - If the run has open HITL questions, open the TUI (or use `--no-tui` mode) to allow answering them.
- **Recovery Actions**:
    - **Retry Failed Agents**: Provide a way to re-run only the agents that failed or were skipped in the resumed run. This appends new `agent.started` events to the existing `events.jsonl`.
    - **Process Check**: Be conservative about "reattachment". If `parley run` was killed, assume child agents are dead. `resume` should focus on restarting the work that didn't finish.

### 3. Data & Schema Refinement
- **Association**: The `run.created` event already links to an `idea`. We should ensure this link is stable and indexed for fast lookup in `parley status`.
- **Event Enrichment**: 
    - Add `pid` to `agent.started` data. While we won't implement complex supervisor reattachment, knowing the PID allows `parley status` to check if a "detached" agent is still physically running.
- **HITL Discovery**: `parley status` should scan `runs/*/questions/*.json` for `status: "open"` to report blockers globally.

### 4. CLI & TUI Updates
- **`parley status --watch`**: Optionally allow a live-updating status view using Bubble Tea.
- **TUI "History" Mode**: Update `tui.RunLive` to support initializing from an existing `store.Store` without requiring a live `runner.Handle`.

## Concerns / open questions
- **Concurrency**: Should we prevent `parley resume` from starting if another CLI instance is already running for that `run-id`? A simple `.lock` file in the run directory would suffice.
- **Process Orphans**: If `parley run` crashes but an agent keeps running, `parley resume` might try to start a second instance of that agent. How should agents handle this? (Likely out of scope for this slice, but worth noting).

## Risks
- **Event Log Bloat**: Repeatedly retrying agents in the same run-id could lead to a long `events.jsonl`. We should ensure the TUI projection logic handles duplicate/overlapping agent events (last event wins for state).
- **Filesystem Drift**: If a user manually deletes a log file or an artifact, the event log becomes an inaccurate reflection of the current state. The "Retry" logic should check both `events.jsonl` and the presence of artifacts.
