---
idea: runtime-status-resume
drafted-by: codex
date: 2026-05-12
---

## Agreed decisions

- Build this slice around a shared runtime projection package, named `internal/runstate`, rather than duplicating event-state logic in the CLI and TUI.
- Move or wrap the existing `tui.ProjectEvents` reducer into `internal/runstate` so `parley status`, `parley resume`, and the live TUI share one tested state projection.
- Add `LoadRun(runDir)` and `ListRuns(root)` style APIs that derive run summaries from durable files:
  - `runs/<run-id>/events.jsonl`
  - `run.created` event data for run-to-idea binding, mode, task, and participants
  - HITL question state through the existing `hitl` package
  - protocol files and idea frontmatter for idea status and artifact presence
- Represent run state conservatively:
  - terminal runs expose `outcome=completed|incomplete|failed`;
  - non-terminal runs expose liveness such as `unverified` or `idle`;
  - `round.incomplete` is terminal but not successful.
- `parley status` remains plain-terminal first and gains run awareness:
  - default workspace overview with transport, ideas, and recent/latest run summaries;
  - `--run <run-id>` detail view;
  - `--idea <slug>` detail view for the newest run for that idea;
  - `--json` with a small unstable developer schema.
- `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA` resolves exact run ID first, then newest run for an idea slug.
- Default `resume` opens a TUI view over durable events/logs/questions. `--no-tui` prints the same run detail as `status --run`.
- Resume is durable-state recovery, not OS process reattachment. It must not claim to adopt, signal, or resurrect subprocesses from a previous CLI invocation.
- The resumed TUI path must be explicit and tested. Do not rely on accidental `Done=nil` behavior without a regression test or dedicated resume option.
- Use HITL package APIs rather than hardcoding question filename patterns in the design.
- Keep implementation small and covered by deterministic fixture tests, then verify with `go test ./...`.

## Agreed trade-offs

- No PID files, process heartbeats, lockfiles, hidden supervisor daemon, or process reattachment in this slice. Those features require a separate supervised-run design.
- No retry of failed or skipped agents in this slice. Retry mutates event-stream semantics and needs its own recovery policy.
- No `parley status --watch` in this slice. `parley resume` TUI is the live-ish view; `status` remains a snapshot.
- No required `state.json` cache in this slice. Recompute from durable events and protocol files first; a cache can be added later if performance demands it.
- No run-id prefix matching or multi-run picker yet. Exact run ID or idea slug keeps resolution predictable.
- JSON output is useful for tests and CI, but it is not yet a stable public API. Mark it unstable in code/help text.

## Open items deferred to implementation

- Exact field names for `RunSummary` and JSON output, provided the semantics above are preserved.
- Exact default count for runs in the workspace overview, as long as detail flags can inspect a specific run or idea.
- Exact stale age wording/default for non-terminal unverified runs.
- Whether artifact presence is reported in the default overview or only in detail mode.

## Deferred follow-up ideas

- `runtime-retry-failed`: retry failed/skipped agents with explicit event and artifact overwrite semantics.
- `supervised-run-state`: process supervision, PID/heartbeat/lock lifecycle, and cross-process liveness.
- `runtime-watch`: workspace-wide live status/watch view if users need it after resume exists.
- `runstate-cache`: optional disposable `state.json`/summary cache if event-log scans become too slow.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-12
Status: ✅ ACCEPT
Notes: Accept. The consensus keeps the slice narrow, shares projection logic, and avoids unsupported process reattachment claims.
