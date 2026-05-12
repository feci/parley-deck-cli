---
idea: runtime-status-resume
status: final
author: codex
consensus-date: 2026-05-12
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Implement the next small `parley-deck-cli` slice: runtime status and resume.

The slice turns durable runtime files into a shared state projection and exposes that projection through:

- `parley status`
- `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA`

`resume` means durable-state recovery and view restoration. It does not mean OS process reattachment.

### Shared run-state projection

Add `internal/runstate`.

Move or wrap the current event reducer logic from `tui.ProjectEvents` so CLI status, CLI resume, and TUI rendering use one tested projection. The TUI should import this shared logic instead of owning the only reducer.

Add APIs equivalent to:

- `LoadRun(root, runID string) (RunSummary, error)`
- `ListRuns(root string) ([]RunSummary, error)`
- `ResolveRun(root, target string) (RunSummary, error)` where `target` is exact run ID first, then idea slug selecting the newest matching run.

Derive state from durable files:

- `parley-deck/runs/<run-id>/events.jsonl`
- `run.created` event fields for idea slug, task, mode, participants, and runtime summary when present
- HITL question state via the existing `hitl` package APIs
- protocol files under `parley-deck/ideas/<slug>/`, especially `00-prompt.md` status and artifact presence
- agent stdout/stderr paths recorded in events

Older or malformed runs should degrade gracefully where possible:

- missing `run.created.data.idea` -> `idea=unknown`
- missing participants -> infer from the idea frontmatter if the idea is known, otherwise empty/unknown
- malformed event line -> return a clear error for that run detail; the workspace overview may skip or mark the run errored rather than failing all status output

### Run state semantics

Expose terminal outcome separately from non-terminal liveness.

Terminal runs:

- `outcome=completed` for `round.completed`
- `outcome=incomplete` for `round.incomplete`
- `outcome=failed` for `run.failed`

Non-terminal runs:

- `liveness=unverified` when at least one agent has `agent.started` without a terminal event
- `liveness=idle` when there is no round terminal event and no agent appears to be progressing from the event stream

Never print unqualified "running" after a restart. Use wording such as:

```text
unverified — last event 12m ago
```

The exact stale threshold and wording can be implementation details, but status must not imply a subprocess is alive unless a future supervised-run design can prove it.

### `parley status`

Keep `parley status [--dir DIR]` plain-terminal first, and extend it with runtime information.

Supported flags:

- `--dir DIR`
- `--run RUN_ID`
- `--idea SLUG`
- `--json`

Default output:

- transport
- idea table using `00-prompt.md` `status:`
- recent or latest run summaries with run ID, idea, outcome or liveness, participant progress, open question count, and last event age

Detail output for `--run` or `--idea`:

- run ID, idea slug, task/mode if available
- outcome or liveness
- agent table: state, elapsed/duration, latest event, artifact path, stdout/stderr log paths, error/reason
- open HITL questions with IDs
- recent event summary
- conservative next action:
  - answer open HITL question with `parley answer ...`
  - inspect with `parley resume <run-id>`
  - no recoverable action; inspect artifacts/logs

`--idea SLUG` selects the newest run whose `run.created.data.idea` matches the slug. If the idea exists but no run exists, return a clear non-zero error.

`--json` returns a small developer-oriented schema based on the same `RunSummary` data. Mark it unstable in code/help text for this slice.

### `parley resume`

Replace the current stub with:

```text
parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA
```

Resolution:

1. exact `parley-deck/runs/<RUN_OR_IDEA>/` directory
2. newest run for idea slug `RUN_OR_IDEA`
3. clear error listing available run IDs or stating that no runs exist

Modes:

- default: open a TUI view over the durable run state
- `--no-tui`: print the same detail body as `parley status --run <run-id>`

The resumed TUI must be an explicit code path, not an accidental `Done=nil` dependency. Add either a `Resume` field or equivalent constructor option and cover it with a regression test.

The resumed TUI may still answer HITL questions through the existing answer flow because answers are durable files/events. This is allowed and should be described as durable-state interaction, not process control.

### Explicit non-goals for implementation

Do not implement in this slice:

- retrying failed/skipped agents
- appending retry lifecycles into an existing run's `events.jsonl`
- PID files
- process heartbeats
- lockfiles
- hidden daemon/supervisor
- cross-process reattachment/adoption/signaling
- `parley status --watch`
- `state.json` as required source of truth
- run-id prefix matching
- multi-run picker
- full protocol validation or automatic phase advancement

### Tests

Add deterministic tests for:

- `runstate.LoadRun` outcomes: completed, incomplete, failed
- non-terminal liveness: unverified and idle
- missing or partial `run.created` fields
- open HITL question counts through the `hitl` package
- `ListRuns` ordering newest first
- `ResolveRun` by exact run ID and by idea slug
- `parley status` default output
- `parley status --run`
- `parley status --idea`
- `parley status --json`
- `parley resume <run-id> --no-tui`
- nonexistent resume target error
- resumed TUI option can be constructed and exits cleanly in tests

Move existing reducer tests out of `internal/tui/live_test.go` if the reducer moves to `internal/runstate`; keep TUI tests focused on rendering and interaction.

### Verification

Implementation is complete when:

- `go test ./...` passes
- manual `parley status` in this repo shows ideas and existing runs
- manual `parley status --run <run-id>` works for at least one existing run
- manual `parley resume <run-id> --no-tui` works and matches the run detail view
- manual TUI resume opens against an existing run without claiming live process reattachment

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
