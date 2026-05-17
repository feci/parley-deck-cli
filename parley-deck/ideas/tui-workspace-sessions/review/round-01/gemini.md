---
agent: gemini
idea: tui-workspace-sessions
review-round: 1
date: 2026-05-17
reviewed-commit: 8a9fe81
---

## Findings

### MAJOR: Inefficient session registration during TUI refresh
In `internal/app/app.go`, the `refreshRuns` callback (line 1377) calls `registerWorkspaceSessions` on every 1-second tick. `registerWorkspaceSessions` (line 1418) iterates over **all** runs in the current workspace and calls `sessionStore.Upsert` for each. `sessionStore.Upsert` (in `internal/sessionstore/sessionstore.go`, line 71) performs a full file read, unmarshal, modification, marshal, and atomic write of the global `sessions.json` file for **every single run**.

In a workspace with $N$ runs and a global session file with $M$ sessions, this results in $O(N \times (M+N))$ work every second. This will lead to severe TUI lag, high CPU usage, and excessive disk I/O as the number of runs grows. It also increases the window for race conditions between multiple `parley` processes.

**Recommendation:** Only register sessions when they are created (already handled by `runcontrol.Create`) or when their terminal state changes. Avoid syncing the entire workspace on every TUI tick.

### MAJOR: Lack of concurrency control in `sessionstore.Upsert`
`sessionstore.Store.Upsert` (line 71) follows a "read-modify-write" pattern without any cross-process file locking. While it uses an atomic rename for the final write, two concurrent `parley` processes (e.g., two parallel `parley run` calls or the TUI and a `parley run`) could both read the same file state, make independent additions, and the one that finishes last will overwrite the other's session entry.

**Recommendation:** Use a file-based lock (e.g., `flock` or a `.lock` file) during the `read-modify-write` cycle in `Upsert`.

### MINOR: TUI refresh loop doesn't account for execution time
The TUI `Init` and `Update` (handling `refreshTickMsg`) trigger a new `refreshRunsCmd` every second regardless of how long the previous refresh took. If `ListRuns` (which reads and projects events for every run on disk) becomes slower than 1 second, multiple refresh commands will queue up, causing the TUI to become unresponsive.

**Recommendation:** Consider scheduling the next tick only *after* the previous refresh has completed, or increase the interval.

### MINOR: Full event log re-reading in `LoadRun`
`runstate.LoadRunAt` (line 86) reads and parses the entire `events.jsonl` file every time it is called. Since the TUI calls this for every run in the workspace every second, this will scale poorly with both the number of runs and the length of the event logs.

**Recommendation:** Implement a cache for run summaries or only read new events since the last load if possible.

### NIT: Inconsistent `Attention` prioritization
In `runstate.Attention` (line 148), `AttentionAction` (open questions) is prioritized over `AttentionFailed` (agent failure). While an open question definitely requires action, an agent failure might be equally or more critical depending on the workflow.

## Test gaps
- **Concurrency tests**: No tests verify that multiple processes can safely update the session store simultaneously.
- **Performance benchmarks**: No tests evaluate the performance of `ListRuns` or TUI refreshes with 100+ runs or large event logs.
- **Race conditions**: Tests for `runcontrol.Create` and `StartAutoAnswerer` should be checked for potential races in highly concurrent environments.

## Verdict
MAJOR. The implementation successfully delivers the requested features and architectural improvements (shared `runcontrol`, session registry, rich TUI), but the O(N^2) session synchronization on every TUI tick is a significant performance regression that will break for active users. Fix the sync logic before merging.
