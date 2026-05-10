---
agent: gemini
idea: live-run-tui
round: 1
date: 2026-05-10
---

## Summary
The proposed approach focuses on creating a real-time, event-driven TUI using the `bubbletea` framework. It leverages the existing `events.jsonl` as the source of truth, ensuring that the TUI remains a decoupled projection of the underlying runner state.

## Proposed approach
1. **Event Sourcing:** Implement an `EventWatcher` that tail-follows the `events.jsonl` file in the current run directory. This watcher will emit Bubble Tea messages whenever new events are appended. This keeps the TUI logic independent of the runner's execution details.
2. **State Model:** Define a `RunState` model that aggregates these events into the required view data:
    - Agent status map (pending, running, finished, failed, skipped, unknown).
    - Start and end timestamps per agent to calculate elapsed time in real-time.
    - A circular buffer of the latest N events for the activity log.
    - Metadata about the active idea, phase, and round.
3. **UI Layout:** Use `lipgloss` to define a responsive three-pane layout:
    - **Header:** Displays the active idea name, phase, round, and an overall round completion indicator.
    - **Agent Sidebar:** A vertical list of participants showing their status icon, name, and elapsed time.
    - **Main Viewport:**
        - **Activity Log:** A scrolling list of the most recent protocol events.
        - **Log Preview:** A dedicated section showing the tail of stdout/stderr for the currently active or selected agent.
4. **Integration:**
    - Update the `run` command to default to the TUI view.
    - Implement the `--no-tui` flag to bypass the TUI and use the legacy logging output.
    - The runner will execute in a separate goroutine, while the Bubble Tea program manages the terminal state.
5. **Termination:** Ensure that quitting the TUI (e.g., via 'q' or 'Ctrl+C') cleans up the terminal state properly. The TUI should distinguish between "user requested exit" and "run completed."

## Concerns / open questions
- **Log Streaming:** What is the most efficient way to preview agent stdout/stderr? Polling the log files on disk is simple but might lag. An in-memory buffer shared with the runner might be faster but adds coupling.
- **Terminal Resize:** Bubble Tea handles resize events, but we need to ensure the `lipgloss` layouts adapt gracefully without breaking the UI.
- **Runner Decoupling:** Should the TUI strictly read from `events.jsonl`, or should it receive events via a Go channel for lower latency? Reading from disk is more robust for "observable runner" goals.

## Risks
- **Race Conditions:** The `events.jsonl` file might not exist immediately when the TUI starts. The watcher must handle missing files gracefully.
- **I/O Overhead:** Excessive polling of `events.jsonl` and agent log files could impact runner performance on resource-constrained systems.
- **UI Flickering:** If events arrive in very rapid succession, we may need to debounce or throttle UI updates to maintain a smooth experience.
