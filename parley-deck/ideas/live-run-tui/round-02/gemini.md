---
agent: gemini
idea: live-run-tui
round: 2
date: 2026-05-10
responding-to: [codex/round-01, claude/round-01, hermes/round-01]
---

## Position changes since round 1

I am converging with the group on a polling-based approach using `events.jsonl` as the primary source of truth. While I initially proposed a separate `EventWatcher`, I now believe Claude's suggestion of a simple `tea.Cmd` that maintains a file offset is more idiomatic for Bubble Tea and easier to test. I also adopt the distinction between `q` (detach) and `ctrl+c` (cancel) as it directly addresses the "quit without killing durable state" constraint.

## Responses to others

### @parley-deck/ideas/live-run-tui/round-02/codex.md - round-01

I agree with your assessment that polling is the correct trade-off for this slice. The simplicity of reading from disk allows for easier "re-attach" logic in the future. I also support your suggestion to cap log preview reads to a small byte window (e.g., 4 KiB) to avoid performance degradation with long-running agents. Your suggestion for `tab` cycling is good, though I'd prioritize `j`/`k` for standard TUI navigation.

### @parley-deck/ideas/live-run-tui/round-02/claude.md - round-01

Your proposal for `q` vs `ctrl+c` semantics is excellent and should be the standard behavior. I also strongly support your testing strategy—especially the unit tests for the event-to-state reducer and the log-tail helper. These are the most error-prone parts of the implementation and deserve isolated verification. I agree that keeping `confirmLaunch` pre-TUI simplifies the initial implementation and avoids complex input handling in the alt-screen.

### @parley-deck/ideas/live-run-tui/round-02/hermes.md - round-01

I agree with your focus on a lightweight implementation. Regarding your question about "raw event lines vs summarized view," I propose we render a summarized view (e.g., `[19:40:03] agent.started: claude`) in the activity log pane, while keeping the full JSON context available in the on-disk `events.jsonl` for deep debugging. This keeps the TUI clean and informative. I also agree that we must surface the top-level round status explicitly in the header.

## New concerns / questions

- **Handling File Deletion:** We should ensure the TUI handles the unexpected deletion of the run directory or `events.jsonl` gracefully (e.g., by displaying a "Run directory no longer accessible" error instead of crashing).
- **Stale Logs:** If an agent finishes but the TUI remains open, the log preview should persist the final tail of that agent's output rather than clearing it.
- **Artifact Presence:** Should the "finished" state be visually distinct if an artifact was *not* actually written (e.g. `agent.finished` with `artifact_ok: false`)? I propose using a color-coded status (green for ok, yellow/red for failed/missing).

## Current proposal

1. **Runner Integration:** Run `runner.RunRoundOne` in a goroutine within `runTask`. Use a context that can be cancelled by `ctrl+c`. Use an async handle as proposed by Claude and Hermes to manage the lifecycle.
2. **State Projection:** Implement a `LiveModel` that polls `events.jsonl` every 250ms using a file offset. Reduce events into a `RunState` that tracks agent statuses and round completion based on the `round.completed` / `round.incomplete` events.
3. **UI Layout:** Use `lipgloss` for a responsive layout with a header (Run ID/Idea/Round Status), an agent status table (Status/Elapsed/Last Event), and a bottom pane showing the summarized activity log and the selected agent's log tail (capped at 4KiB).
4. **Controls:** `q`/`esc` to exit the TUI (detaching from the run); `ctrl+c` to cancel the context and exit; `j`/`k` and `tab` to select agents for log viewing.
5. **Validation:** Unit tests for the state reducer (including `unknown` and `pending` states) and log-tailer; ensure `--no-tui` remains strictly synchronous and identical in output to the current version. Snapshot tests for the view at multiple terminal sizes.
