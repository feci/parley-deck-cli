---
agent: claude
idea: tui-workspace-sessions
review-round: 1
date: 2026-05-17
reviewed-commit: 8a9fe81
---

## Findings

### 1. MAJOR — TUI `startRun` discards its `CancelFunc`; quit cannot cancel TUI-started runs, footer claim is false, and child processes may orphan

`internal/app/app.go:1395`:

```
runCtx, _ := context.WithCancel(startCtx)
```

The `CancelFunc` is assigned to `_`. Nothing — neither the TUI nor the runner — ever calls it, so:

- `internal/tui/app.go:549` renders `"Quit closes the TUI and cancels TUI-started runs; it does not detach child processes."`. This is not what the code does. Pressing `q`/`esc` simply invokes `tea.Quit` (`internal/tui/app.go:166-167`); no cancellation runs against `runCtx`.
- The agent subprocesses launched by `runner.RunRoundOneAsync(runCtx, ...)` are not signaled when the TUI exits. When the parent process exits, those children re-parent to PID 1 and continue (they are not in the parent's TTY foreground group and `exec.Cmd` does not set up process-group kill here), which can leave them writing artifacts after the user believes the run is over.
- FINAL.md §5 ("Define quit and resume semantics clearly") requires either a prompt-before-cancel/close path or an accurate footer rendering of the actual behavior. Neither holds: the footer asserts a behavior the implementation does not perform, and there is no confirmation prompt.
- This is also a `go vet` "lostcancel" violation.

Suggested fix: keep the cancel handles on the model keyed by `RunID` (return them out of `StartRunFunc`), call them from the `q`/`esc`/`ctrl+c` branch before returning `tea.Quit`, or — at minimum — rewrite the footer to describe what the code actually does and document the orphan risk.

### 2. MAJOR — Refresh tick can silently shift the user's selected run because selection is tracked by index, not `RunID`

`internal/tui/app.go:204-212`:

```
case refreshRunsMsg:
    if msg.err != nil {
        m.errText = msg.err.Error()
    } else {
        m.errText = ""
        m.runs = msg.runs
        m.clampSelections()
    }
```

`runstate.ListRuns` sorts descending by `RunID` (`internal/runstate/runstate.go:213-215`), and `RunID` is timestamp-derived, so any newly-created run — whether started from this TUI, another terminal, or another `parley run` invocation — is prepended on the next 1 s refresh tick. `clampSelections` only clamps the integer index; it does not look up the previously selected `RunID`. The result: a user sitting on `selectedIdea = 2` watching one run can find their selection silently moved to a different run a second later, with `selectedAgent` then clamped against an unrelated run's agent set. The same hazard applies to `selectedAgent` inside a single run if a new agent first appears.

The `startRunMsg` handler already gets this right (`indexRun(m.runs, msg.run.RunID)`, `internal/tui/app.go:228`). Apply the same pattern in the refresh path: remember the currently-selected `RunID` before assigning `m.runs = msg.runs`, then re-resolve the index after assignment (and similarly for `selectedAgent` by agent ID).

### 3. MINOR — `registerWorkspaceSessions` performs an O(N) read-modify-rewrite of `sessions.json` per refresh tick

`internal/app/app.go:1418-1444` iterates all runs and calls `sessionStore.Upsert(...)` once per run. Each `Upsert` (`internal/sessionstore/sessionstore.go:85-128`) reads `sessions.json`, mutates the in-memory document, then atomically rewrites the entire file via temp-file + `Rename`. With the 1 Hz refresh tick (`internal/tui/app.go:808-812`) and a workspace of, say, 50 runs, this is 50 reads + 50 full rewrites per second of `~/.parley-deck/sessions.json`. It is functionally correct but worth batching: take the doc once, apply all run upserts in memory, write once.

## Test gaps

- No test asserts that quitting the TUI (or pressing `ctrl+c`) cancels `runCtx` for TUI-started runs, or that subprocesses are terminated. The current discard makes such a test impossible to write; adding the test would catch the regression in Finding 1.
- No test asserts that `selectedIdea`/`selectedAgent` follow the previously-selected `RunID`/agent ID across a `refreshRunsMsg` that reorders or prepends entries (Finding 2). The unit tests in `internal/tui` should include a fixture where a new run appears at index 0 between ticks.
- No test exercises `sessionstore.Upsert` under repeated calls within a single refresh cycle or asserts that a batched-update path keeps `LastEventAt`/`Terminal` correct (Finding 3 follow-up).
- `sessionstore.Upsert` has no zero-value guard on `LastEventAt` analogous to its guards on `Task`/`Participants`/`IdeaSlug`/`CreatedAt`; a regression that passes a zero `LastEventAt` would silently clobber a populated value. Worth a table-test case even though no current caller triggers it.

## Verdict

minor-fixes-needed — Findings 1 and 2 should be addressed before merge: Finding 1 because the footer makes a promise the code does not keep and subprocesses can outlive the TUI, and Finding 2 because the selection-by-index regression is a routine UX hazard on a 1 Hz refresh. Finding 3 is a polish item.
