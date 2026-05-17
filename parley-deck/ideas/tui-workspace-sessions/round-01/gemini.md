---
agent: gemini
idea: tui-workspace-sessions
round: 1
date: 2026-05-17
---

## Summary

The current Parley Deck CLI provides a static workspace dashboard and a live view for a single run. To support "Workspace Sessions," we must transition to a state model that manages multiple concurrent runs, tracks session metadata across process restarts, and provides a race-free mechanism for aggregating events and HITL requirements from diverse run directories. My focus is on the state model, persistence strategy, and concurrency patterns required to make this robust and responsive.

## Proposed approach

### 1. User-Local Session Registry
We should introduce a `sessions.json` file in `~/.parley-deck/` to track runs across different workspace roots. This allows the TUI to "remember" active ideas even if they were started in a different terminal or if the TUI was restarted.

**State Model for a Session:**
- `RunID` (string): Primary key, matches the `events.jsonl` directory.
- `WorkspaceRoot` (string): Absolute path to the repository.
- `IdeaSlug` (string): Short name for the idea.
- `Task` (string): The original task description.
- `StartTime` (time.Time): When the run was created.
- `Status` (string): Last known terminal state or "Active".
- `PID` (int, optional): The process ID of the runner if it's an external process.

### 2. Multi-Run State Projection
The TUI needs a "Session Manager" that periodically refreshes the status of all registered active runs.
- **Incremental Loading**: Reuse the existing `readEventsFromOffset` pattern to ensure we only read new events from each `events.jsonl`.
- **Background Refresh**: Each active run should have its own polling cycle (via `tea.Tick`) to fetch new events and update its local projection.
- **Unified Status Derivation**: Enhance `runstate.ProjectEvents` to explicitly identify "User Action Required" states (e.g., `OpenQuestions > 0` or a new `run.blocked` event).

### 3. Concurrency and Liveness
To avoid "Zombie" runs (runs that appear active but the process has died), we need a liveness check.
- **Heartbeat Events**: The runner should optionally emit periodic `run.heartbeat` events.
- **Process Check**: If a `PID` is recorded, the TUI can verify if the process still exists.
- **Timeout Fallback**: If a run is in `Running` state but has seen no events for a threshold (e.g., 5 minutes), it should be flagged as `Stale` or `Unverified`.

### 4. Event Stream Aggregation
A global "Event Stream" view in the TUI will require a merged view of events from all active runs.
- **Time-Ordered Buffer**: Maintain a shared, sorted buffer of recent `EventSummary` objects across all sessions.
- **Efficient Updates**: Only re-sort the buffer when new events arrive from any run.

## Concerns / open questions

- **Persistence Performance**: If many runs are active, polling dozens of `events.jsonl` files every 250ms might cause high I/O. We should probably increase the polling interval for background (non-selected) runs.
- **Multi-Process Safety**: If the user runs `parley answer` in another terminal, the TUI should reflect the change immediately. The current file-based HITL status handles this well, but we must ensure we don't hold file locks.
- **Session Cleanup**: When does a session leave `sessions.json`? We should keep a history of "Recent" finished runs but eventually prune old ones to keep the file small.

## Risks

- **Race Conditions in State Updates**: If the TUI's state projection logic is not strictly sequential within the Bubble Tea loop, we might see inconsistent UI states (e.g., an agent appearing finished in one panel but running in another).
- **Stale Metadata**: If a run directory is manually deleted, `sessions.json` will point to a non-existent path. The state model must gracefully handle missing files.
- **TUI Responsiveness**: Large log previews or many parallel runs could block the main UI thread. Log tailing should be strictly limited to a few KB and executed asynchronously.
