---
agent: claude
idea: tui-workspace-sessions
round: 1
date: 2026-05-17
---

## Summary

Treat the TUI as a small workspace-scoped session console that wraps the
existing per-run live view, not as a new state store. The user's mental model
should be three nested levels — workspace → session → agent — with one rule
held constant at every level: "what currently needs me" is always the most
visible thing on screen. `events.jsonl` per run remains the source of truth;
`~/.parley-deck` is purely a cross-workspace index to help find and resume
those runs.

## Proposed approach

### Screen model and navigation

Three nested Bubble Tea models, with explicit focus zones (mirroring the
pattern already in `internal/tui/app.go:14-29`):

1. **Workspace screen** (replaces the current dashboard for the active
   path). Left rail: list of sessions ordered by "needs user action" first,
   then "running", then "recent finished". Right rail: discovered agents
   panel from today's `app.go`, mostly unchanged. Footer always shows the
   global "needs answer" count across all sessions.
2. **Session screen** (a generalized `live.go`). Header with idea/run/round
   and status. Body split into: agents table, time-series event pane, HITL
   questions pane, log preview of the focused agent. Same keybindings as
   today's `live.go:116-155`, plus `tab`/`shift+tab` to switch zones and
   `enter` on an agent to drill into the per-agent view.
3. **Per-agent screen**. Full-height log view (stdout + stderr) for one
   agent in one session, with scrollback, and a sidebar showing the
   agent's state, artifact path, duration, last event, and any open
   question authored by that agent. `esc` returns to the session.

Global keys (from any screen): `q` close current screen (does not cancel
runs), `?` help, `g` go to global question queue (see below), `s` open
session by run-id, `n` new idea.

This separation keeps the existing "single live view" intuitive for users
already familiar with `parley run` while letting them step up one level to
see and start parallel work, and step down one level to focus on a single
agent without losing context.

### Starting new ideas from the TUI

Bind `n` on the workspace screen to a small modal flow:

1. Free-text task (single line, like the CLI `TASK` arg).
2. Participant multi-select, defaulted to all installed agents from the
   already-loaded `[]agents.Discovery`. Show missing agents greyed out.
3. Mode: HITL (default) or auto.
4. Final confirmation showing the same launch summary that
   `confirmLaunch` prints today in `internal/app/app.go:1479-1493`.

Submit calls the same code path used by `runTask`
(`internal/app/app.go:1140-1208`): `protocol.CreateIdea`, append
`run.created` event, then `runner.RunRoundOneAsync` with the run's
`context.Context` stored on the session. The TUI inserts the new session
at the top of the rail and selects it. Hitting `enter` opens the session
screen for the new run; pressing `esc` keeps it running in the
background.

Important: nothing in the existing `parley run`/`parley resume` flag
shapes changes; the modal is just a new way to invoke the same
operation.

### Multiple parallel sessions

Introduce a small in-memory `sessionRegistry` owned by the workspace
model. Entry per run:

```
type session struct {
    runID       string
    runDir      string
    idea        protocol.IdeaStatus
    participants []string
    handle      *runner.Handle   // nil when the run was discovered, not started here
    cancel      context.CancelFunc
    offset      int64            // events.jsonl tail offset, per session
    state       runstate.RunState
    questions   []hitl.Question
    lastEvent   time.Time
    terminal    bool
    needsUser   bool
    failureText string
}
```

The workspace model holds `sessions []*session`. A single global tea ticker
(250ms, the existing cadence from `live.go:496-499`) walks the registry and
emits per-session `tailMsg{runID, events, questions, offset}` messages.
Update() applies them by calling `runstate.ProjectEvents` per session — no
shared mutable state across sessions, so concurrent runs cannot corrupt
each other's derived state. We never re-read full files; offsets are kept
per session (the same logic already in `live.go:518-556`, just lifted out
of the live model).

For sessions started in this TUI process, we keep `handle.Done()` wired up
to flip `terminal = true` and stop tailing. For sessions that were started
elsewhere (or in a previous TUI run), we tail events.jsonl until we see
`round.completed`, `round.incomplete`, or `run.failed`.

### Surfacing user-action states clearly

This is the lens I care about most. A session is in exactly one **action
state** at a time:

| Action state    | Trigger                                                       |
|-----------------|---------------------------------------------------------------|
| `needs-answer`  | At least one `hitl.Question` with `Status == open`            |
| `needs-signoff` | round terminal **and** consensus triage = `pending-signoffs`  |
| `running`       | any `AgentState` in `running`, no open questions              |
| `blocked`       | `run.failed`, or terminal with `outcome == incomplete/failed` |
| `done`          | terminal and outcome `completed`, no pending consensus        |
| `idle`          | not terminal, no running agent, no event for &gt; idle window |

Rendering rules:

- Workspace rail sorts `needs-*` first, then `running`, then everything
  else; rows render with a single distinct glyph + color per state
  (`!` red for `needs-answer` / `needs-signoff`, `>` yellow for
  `running`, `*` orange for `blocked`, dim for `done` / `idle`). The
  glyph + state token is the first thing on the line.
- The global footer always shows the total open-question count across
  all sessions, so the user notices even on a screen that does not
  belong to the affected session.
- `g` opens a **global questions queue**: one merged list of open
  questions across all sessions, sorted by `created_at`. Selecting one
  jumps to its session and pre-selects the question in the existing
  questions pane (`live.go:288-325`). This is the simplest way to make
  "which idea needs me right now" answerable in one keystroke.

This state derivation is pure: input is `(events, questions,
consensus-summary)`, output is the action state. It's easy to test
exhaustively without launching agents.

### Event time series

Per-session: keep the existing 8-event "Recent" window from
`runstate.ProjectEvents` for compact rendering, but also maintain an
append-only **scrollback** buffer of `EventSummary` items for the
session screen, scrolled by `pgup/pgdn`. New events are highlighted for
~1s.

Workspace-wide: a small bottom strip on the workspace screen showing the
last ~3 events from any session, formatted as `HH:MM:SS  <run-id>
<type>  <text>`. This is the "something happened somewhere" lane.

We deliberately do not introduce a new file format — the scrollback is
built in memory from the per-session offset reads, and lost on TUI exit.
Anything durable already lives in `events.jsonl`.

### Per-agent views

Drilling into an agent (`enter` on an agent row) opens a focused screen
that, for the selected `AgentState`:

- Shows the agent's lifecycle: `started_at`, `duration`, `state`,
  `latest_event`, `artifact`, `error`/`reason`.
- Tails the agent's stdout and stderr using the existing `tailLogFile`
  helper (`live.go:558-598`) but with a larger window (e.g. 64 KiB, 200
  lines) and scrollback.
- Lists the agent's open and recently answered questions so the user can
  answer in place with `a` (same flow as today's `live.go:338-375`).
- `tab` cycles to the next agent in the same session without going back
  up to the session screen.

This gives a real "per-agent tab" feel without introducing actual
Bubbles tab widgets — `j/k` between agents on the session screen and
`enter` to focus is enough.

### Persistence under `~/.parley-deck`

Minimal user-local index, written atomically (tmp + rename), shape:

```
~/.parley-deck/sessions.json
{
  "version": 1,
  "sessions": [
    {
      "workspace": "/abs/path/to/repo",
      "run_id": "20260517T120000.000000000Z",
      "idea_slug": "tui-workspace-sessions",
      "participants": ["claude","gemini","hermes"],
      "started_at": "...",
      "last_seen_at": "...",
      "terminal": false
    }
  ]
}
```

Write rules:

- Append a row on `run.created` (only from TUI-initiated runs in this
  slice; CLI `parley run` can also write it for parity).
- Update `last_seen_at` and `terminal` whenever the TUI observes a
  state transition for that run, throttled to at most one write per
  second per session.
- Never write secrets, logs, prompts, task text, or events here.
- The index is advisory only. Reading it gives `parley tui` a "recent
  sessions across workspaces" list; opening any of them still requires
  the original `events.jsonl` on disk and falls back gracefully if the
  workspace path is missing.

### Resume semantics

From the workspace rail, `r` on a non-running row reopens it via the
same code path as `parley resume`
(`internal/app/app.go:648-696`). It calls `resumePendingConsensusSignoffs`
first, then enters the session screen with `Resume: true` so footer text
swaps to the resume variant exactly as today.

### Test strategy (correctness)

- A new `internal/tui/workspace` package exposes a pure reducer:
  `Update(model, msg) (model, []cmd)`. Tests feed synthetic
  `tailMsg`, `questionsMsg`, `doneMsg`, and key events, and assert on
  derived `actionState` and selected session. No real agent CLIs.
- Golden-text rendering tests at fixed width (e.g. 100 and 140) for
  workspace, session, and per-agent screens, covering each action
  state. The existing pattern of stripping ANSI in `live.go:600-602`
  means we can compare plain text snapshots.
- Action-state derivation gets its own table-driven test covering all
  rows in the table above.
- Persistence: round-trip test for `~/.parley-deck/sessions.json` —
  write, read, dedupe, throttle.

## Concerns / open questions

1. **Ctrl-C semantics across nested screens.** Today `live.go:128-132`
   treats `ctrl+c` as "cancel the run". With many parallel runs in one
   TUI, that key has to mean "close current screen", and cancellation
   must move to a deliberate key (e.g. `ctrl+x` from the session
   screen, with confirmation). This is a small but user-visible
   compatibility change for users coming from today's `parley run`.
2. **Closing the TUI with live children.** When the user quits the TUI
   process while runs are active, do we cancel them or detach? The
   existing single-run behavior is "cancel on ctrl-c". For multi-run we
   should prompt: "N runs still active — detach (keep running) /
   cancel all / back". Detaching keeps `events.jsonl` advancing under
   the existing runner goroutines, which is what users will expect for
   long-running idea rounds.
3. **Global question queue vs. notification model.** Is a single
   queue ergonomic enough, or do we need something more intrusive
   (e.g. a popup) when a new question arrives on a non-focused
   session? My instinct is "no popup, just promote the rail row and
   update the footer counter"; a popup interrupts whatever else the
   user is reading.
4. **Naming.** Internally the prompt uses "session". The codebase
   already calls these "runs". I'd keep "run" as the protocol-level
   noun (because that's what's on disk under `parley-deck/runs/`) and
   call the TUI's per-run shell a "session view". If we conflate them,
   we'll repeatedly have to disambiguate in help text and prompts.
5. **Event source for an in-process run.** A run started from this TUI
   could push events into the model directly without going through the
   file tail; but doing that creates two code paths for "how the model
   learns about an event". I'd prefer to keep one code path — always
   tail `events.jsonl` — even for runs we started ourselves. Simpler
   to reason about and to test.
6. **fsnotify vs. polling.** Per-session polling at 250ms scales fine
   for the expected handful of parallel runs and matches the existing
   live view's tick. fsnotify would be tighter but adds a dependency
   and OS-specific edge cases. MVP says polling; we can revisit if
   users actually run dozens of sessions.

## Risks

- **Correctness of multi-session reducer.** The biggest risk is that
  Update() ends up holding shared state per agent rather than per
  `(runID, agentID)`. The reducer must key everything by run ID first.
  Tests must specifically cover "two sessions, same agent ID, both
  running, one fails" — easy to get wrong, hard to spot by eye.
- **Refresh cost regression.** A naive port of `live.go` to the
  workspace would re-read every events.jsonl on every tick. We must
  keep `offset` per session (as today's live model does for one run)
  or large workspaces will slow the UI noticeably.
- **Stale `~/.parley-deck` entries.** The user-local index can point
  to workspaces that were moved or deleted, or to runs whose events
  file is gone. The UI must show these as "unavailable" rather than
  erroring; a self-pruning rule (drop entries older than N days when
  the index is opened) keeps the file from growing unboundedly.
- **Background goroutine leaks.** Each session owns at least a
  context, possibly a `runner.Handle`, and tail goroutines. Closing
  the session view must not leak any of these, and quitting the TUI
  must reap them deterministically. A single `sessionRegistry.Close()`
  that cancels every context and waits on every `Done()` channel
  before `tea.Quit` returns is the safest pattern.
- **HITL race.** Today the answer write path goes through
  `hitl.Store` (`internal/hitl/hitl.go:112-141`) with atomic
  rename. If two TUIs (one local, one remote ssh into the same
  workspace) try to answer the same question, second write wins.
  Acceptable for MVP, but should be called out in docs so users do
  not expect collaborative answering.
- **Backward compatibility on keybindings.** Adding `n`, `g`, `r`,
  `enter`, `pgup/pgdn` should not collide with anything in today's
  `live.go` `Update`. They do not, but the `n` key today moves to the
  next question (`live.go:142-144`). I propose changing the "new
  idea" key to `N` (capital) to avoid the collision; that's cheap and
  preserves existing muscle memory.
- **Persistence schema lock-in.** Even though `sessions.json` is
  advisory, once we ship it we'll get attached to its shape. Keep a
  top-level `"version"` field from day one and ignore unknown
  versions gracefully.
- **Resume of terminal runs.** Re-opening a finished run as a
  "session" should be visibly distinct from a live one (header
  badge, no tail goroutine spawned). Otherwise users will sit
  waiting for events that will never come — exactly the kind of
  correctness/UX bug the `Resume: true` flag in `live.go` was meant
  to address; we just need to preserve that distinction at the
  workspace level too.
