---
agent: gemini
idea: live-run-tui
review-round: 1
date: 2026-05-10
implementation-pr: https://github.com/feci/parley-deck-cli/pull/2
---

## Verdict
Request changes

## Findings
- CRITICAL <internal/tui/live.go:210> Table alignment is broken by ANSI color codes. `fmt.Sprintf("%-16s", state)` calculates padding based on string length including non-printable ANSI sequences (e.g., from `okStyle.Render`), causing columns to shift. Use `lipgloss.Width()` or lipgloss table/placement primitives to maintain alignment.
- CRITICAL <internal/tui/live.go:193> Column headers in `renderAgentTable` are misaligned with data rows. The header starts with 2 spaces but row logic uses `%s ` (marker + space) which results in offset alignment for all subsequent columns.
- MAJOR <internal/tui/live.go:343> Bug in `agent.skipped` duration calculation. If an agent is skipped without ever starting (e.g., artifact already exists), `agent.StartedAt` is zero, causing `now.Sub(agent.StartedAt)` to return a massive duration (time since epoch).
- MAJOR <internal/tui/live.go:167> Hardcoded `round=round-01` in the header. While correct for the current slice, it should use `m.opts.Idea.Status` or be passed via `LiveOptions` to support future rounds as intended by the "phase/round" requirement in FINAL.md.
- MAJOR <internal/tui/live.go:312> `ProjectEvents` does not preserve `stateUnknown` for agents not in the participant list. If an "unknown" agent emits multiple events, subsequent events will call `applyAgentEvent` which may overwrite the unknown state with `running` or `finished`, though the agent remains in the table.
- MINOR <internal/tui/live.go:445> `displayDuration` returns 0 for very fast agents (0ms), which `formatDuration` renders as `-`. This makes finished agents look like they haven't started. It should distinguish between "not started" and "0ms duration" (e.g., by checking state).
- MINOR <internal/tui/live.go:463> `displayRoundStatus` returns "unknown" if the runner finishes (`done=true`) but the `round.completed` event hasn't been processed yet. While technically following the "no inference" rule, it may be confusing to users if the runner finishes successfully but the TUI briefly shows "unknown" before exiting.
- NIT <internal/tui/live.go:187> Footer text mentions `tab` but not `shift+tab` for reverse selection, although the latter is implemented in `Update`.
- NIT <internal/tui/live.go:153> `Update` for `eventsMsg` exits immediately if `m.done` is true. This prevents the user from seeing the final state of the TUI (including the `round.completed` status) if it arrived in the final read.

## Test gaps / residual risk
- No test for `displayRoundStatus` logic.
- `TestProjectEventsDerivesAgentAndRoundState` does not verify that `Error` (for failed) or `Reason` (for skipped) are correctly projected into `AgentState`.
- `TestTailLogFileReturnsBoundedWholeLines` uses a small file (~60 bytes) and does not verify the `Seek` logic for truly large log files (e.g. > 4 KiB) to ensure partial-line stripping from the start of the buffer works as intended.

## Summary
The implementation successfully delivers the core async runner and event-polling TUI architecture. However, the TUI's visual integrity is compromised by ANSI-aware padding issues in the agent table, and there is a significant logic bug in duration calculation for skipped agents. The hardcoding of the round label in the header should also be addressed to align with the design's multi-round intent.
