---
agent: claude
idea: runtime-status-resume
round: 1
date: 2026-05-12
---

## Summary

`parley status` already lists workspace ideas, and `parley run` already writes
durable `events.jsonl`, agent logs, and HITL `questions/*.json` under
`parley-deck/runs/<run-id>/`. The smallest useful next step is: (1) extend
`parley status` to project those durable files into a per-run summary that
works in CI and plain terminals; (2) implement `parley resume` as a
**read-only re-attachment to the durable state**, not to live OS processes.
Reuse the existing `tui.ProjectEvents` reducer by lifting it into a
non-TUI package so both `status` (plain/JSON) and `resume` (TUI or
`--no-tui`) share one tested projection.

## Proposed approach

### 1. Carve out a shared projection package

Move `ProjectEvents`, `RunState`, `AgentState`, `EventSummary` from
`internal/tui/live.go` into a new `internal/runstate/` package. The TUI keeps
its existing render code and depends on `runstate` for the reducer. No
behavior change for `tui` — just an import shuffle. Existing tests in
`internal/tui/live_test.go` move with the reducer.

Add to `internal/runstate/`:

- `LoadRun(runDir string) (RunSummary, error)` — wraps `store.Load`,
  reads `run.created` to recover idea/task/participants/mode, calls
  `ProjectEvents`, also lists `questions/` via `hitl.New(runDir).List()`.
- `ListRuns(root string) ([]RunSummary, error)` — enumerates
  `parley-deck/runs/*`, run-ids are RFC-3339-ish sortable strings so the
  lexicographic order matches chronological order; latest first.
- `RunSummary` fields: `RunID`, `IdeaSlug`, `Mode`, `Participants`,
  `State` (`runstate.RunState`), `OpenQuestions int`, `LastEventAt time.Time`,
  `LastEventAge time.Duration`, `Liveness` (see below).

### 2. `Liveness` — conservative state model

Three values, computed from events only:

- `complete` — terminal `round.completed`, `round.incomplete`, or `run.failed`
  event present, and every started agent has a `finished|failed|skipped`.
- `idle` — no `round.*` event yet; every started agent has a terminal event;
  there is at least one agent still in `pending`. Means: round was launched,
  partial completion, no runner is provably progressing it.
- `unverified` — at least one agent has `agent.started` with no terminal
  event. We do **not** claim it is alive. The renderer prints
  `running (unverified — last event <age> ago)`. A `--stale-after` flag
  (default 10m) only changes how the label is colored/worded; it does not
  let us assert a process is dead.

This is the whole story of "is it still running?" in this slice. We avoid
PID files, lockfiles, and supervisor daemons. They belong in a follow-up
idea once the surface area is justified by user need.

### 3. `parley status` surface

Replace the body of `runStatus` in `internal/app/app.go`:

- `parley status` (no args) — current ideas table + a runs table. Latest
  run per idea: run-id, mode, participants, liveness, open-questions,
  last-event age. Plain text, deterministic order.
- `parley status --run <run-id>` — full detail: agent table (state /
  elapsed / last event / log paths), 8 most recent events (already
  produced by `ProjectEvents`), open HITL questions with IDs, suggested
  next command (one of: `parley answer …`, `parley resume <run-id>`,
  `parley run --participants <missing-ids> …`, or "nothing to do").
- `parley status --idea <slug>` — same detail view for the latest run of
  that idea.
- `parley status --json` — same data as a single JSON object so CI can
  parse it. Keep the field schema small (RunSummary + agent list +
  question list); document it in a code comment, not a new doc file.

`--dir`, `--no-tui` semantics already in place stay unchanged.

### 4. `parley resume <idea-or-run>`

Resolution: if the argument matches an existing `runs/<id>/` directory,
use it; else treat it as an idea slug and pick the lexicographically
greatest run-id whose `run.created.data.idea` matches. If neither
resolves, return a clear error listing candidate run IDs.

Two modes:

- **Default (TUI)** — call `tui.RunLive` with a `LiveOptions` synthesized
  from `runstate.LoadRun`. `Done` is `nil` (the existing `waitDoneCmd`
  already handles nil by returning nil), `Cancel` is a no-op. The TUI is
  then a read-only tail of `events.jsonl`, log files, and `questions/`,
  with the existing `a` keybinding still working to answer open HITL
  questions because `hitl.Answer` writes through the durable JSON files
  regardless of which process is "running" the agent.
- **`--no-tui`** — print the same plain-text body as
  `parley status --run <id>` and exit. Useful in CI / for piping.

What resume is **not**:

- It does not adopt or signal subprocesses started by an earlier `parley
  run`. On Unix, a child belongs to its parent process group; we cannot
  reparent it. If `parley run` was killed mid-round, its agent
  subprocesses are either already gone or orphaned and will exit on
  their own; either way `parley resume` only observes their durable
  output and accepts HITL answers.
- It does not retry failed or skipped agents. Retry is "run with
  `--participants <id>` against the same idea", which is a separate
  follow-up.

The CLI usage line in `printUsage` is updated to reflect this:
`parley resume RUN_OR_IDEA [--no-tui]`.

### 5. Test plan

Driven by deterministic `events.jsonl` fixtures:

1. `runstate_test.go` — given a fixture events file + a `run.created`
   line, `LoadRun` returns expected `RunSummary` with correct liveness for
   each of: round-completed, partial-with-running-agent, partial-with-no-
   running-agent, run-failed.
2. `app_test.go` — `parley status` against a tempdir workspace containing
   one finished run + one `unverified` run produces deterministic text.
3. `app_test.go` — `parley status --json` parses back to the same struct.
4. `app_test.go` — `parley resume <run-id> --no-tui` matches
   `parley status --run <run-id>`.
5. `app_test.go` — `parley resume nonexistent` exits non-zero with a
   message that lists available run IDs (or says "no runs yet").
6. TUI integration: move existing `live_test.go` projection tests to the
   new package; keep `live_test.go` covering only rendering.

### Verification

- `go test ./...` green.
- Manual: `parley run --no-tui …` to produce a run, then
  `parley status`, `parley status --run <id>`, `parley status --json`,
  `parley resume <id> --no-tui`, `parley resume <slug>`. The TUI resume
  case is verified by opening it against a completed run and confirming
  the agent table / events / questions panes match the live-run
  rendering.

## Concerns / open questions

- **Where does run→idea binding live?** Today the only durable link is
  `run.created.data.idea`. I am relying on that. If a run was created by
  an older binary that did not record the idea slug, `runstate` should
  fall back to `unknown` rather than fail. Worth confirming we never
  shipped a binary without that field; if we did, we may need to scan
  `00-prompt.md`/`round-01/<agent>.md` headers as a backup.
- **Suggested next command logic.** "Suggest a next step" can creep into
  policy. Keep it to four cases listed above; anything richer is a
  follow-up.
- **JSON schema stability.** Once we publish `--json`, downstream
  callers will depend on it. Mark it `unstable` in a code comment for
  this slice; revisit when there is a real consumer.
- **Multiple-run-per-idea UX.** `parley status` lists latest only. Should
  we add `--all-runs`? My instinct is no for this slice — defer until a
  user asks.
- **Run-id prefix matching for resume.** Tempting but adds ambiguity
  rules. Skip in v1; require exact run-id or idea slug.
- **HITL answer rendering for already-answered questions.** Resume's
  question pane defaults to "open only", matching live-run behavior.
- **Concurrent producer.** If another `parley run` is still appending to
  the same `events.jsonl`, the offset-based tail in `tui.RunLive` keeps
  working — same code path as today. No new locking needed.

## Risks

- **Implying liveness we cannot prove.** If the renderer says "running"
  for an `agent.started` with no terminal event, users will assume a
  process is alive. The `unverified — last event <age> ago` label
  mitigates this but does not eliminate it. We must not call the state
  `running` unqualified.
- **Refactor regression.** Lifting `ProjectEvents` out of
  `internal/tui/live.go` touches the most-tested file in the repo. Keep
  the function signature identical, move tests with it, and run the
  whole suite before opening the PR.
- **Scope creep into supervision.** Reviewers may push for PID files or
  a daemon "while we are in here". Resist: that is a separate idea with
  its own design surface (file-lock semantics, crash-recovery, signal
  forwarding). This slice deliberately stops at durable-state recovery.
- **`events.jsonl` growth.** Each `status` call reads the whole file. At
  current event volumes this is fine; flag as a known future
  optimization (offset cache or partial tail). Do not optimize now.
- **TUI read-only resume can still emit HITL answers.** That is correct
  behavior — `hitl.Answer` writes durable files and the in-flight agent
  (if any) polls them — but worth documenting so reviewers do not
  mistake it for "writing while in read-only mode".
