---
idea: live-run-tui
drafted-by: codex
date: 2026-05-10
---

## Agreed decisions

- Default `parley run` should launch a live Bubble Tea TUI while the runner executes in the background. `--no-tui` must preserve the existing synchronous runner path and text summary behavior.
- The TUI must derive runtime state from the active run directory, with `events.jsonl` as the source of truth. It should poll the event file with a retained byte offset and handle a missing or not-yet-created event file as "no events yet".
- The runner/TUI bridge should be small: either an app-layer goroutine wrapper or a narrow runner handle that exposes the run directory and a completion result through `Done()`/`Wait()`. The implementation should avoid a second in-memory event model for this slice.
- Agent display states are `pending`, `running`, `finished`, `failed`, `skipped`, and `unknown`. Round status is derived only from `round.completed` or `round.incomplete`, never inferred from all agents being terminal.
- The log preview should read bounded tails from stdout/stderr logs for the selected agent only, cap reads to a small byte window such as 4 KiB, strip ANSI/control noise where practical, and show whole lines.
- The TUI layout should include a header with idea, round, run ID, and round status; an agent table with state and elapsed time; a summarized latest-events pane; a selected-agent log preview; and a footer explaining keys.
- Keyboard controls for this slice are `j`/`k` and `tab` to change selected agent, `q`/`esc` to detach from the TUI while the runner continues, and `ctrl+c` to cancel the run context and exit.
- Testing should cover event-to-state reduction, including `pending` and `unknown`; log-tail behavior; stable rendering at a few terminal sizes; and regression coverage that `--no-tui` remains synchronous.

## Agreed trade-offs

- Disk polling is intentionally chosen over fsnotify or Go channels because it keeps the durable event log honest and makes future re-attach/resume simpler.
- The implementation should not add token accounting, HITL answer input, resume, GitHub/GitLab automation, or packaging work in this slice.
- The first version should render summarized event rows, not raw JSON, to keep the TUI readable. Raw event data remains available in `events.jsonl`.
- The first version should not add an artifact status column unless the existing event schema already exposes reliable data for it. Avoid inventing states the runner cannot prove.

## Open items deferred to implementation

- If the run directory or event file is deleted while the TUI is open, render a clear error instead of crashing.
- If terminal width is narrow or color support is weak, preserve readability with text markers such as a `>` selected-agent prefix.
- A future slice may add re-attach/resume, raw JSON event expansion, token accounting, and HITL response entry.

## Signoffs

### Signoff: codex - 2026-05-10
Status: ACCEPT
Notes: This consensus matches the shared round-02 direction and keeps the implementation slice narrow enough to verify.

### Signoff: gemini - 2026-05-10
Status: ACCEPT
Notes: The plan effectively balances real-time visibility with implementation simplicity by using the event log as the source of truth.

### Signoff: hermes - 2026-05-10
Status: ACCEPT
Notes: The narrow TUI slice using event polling aligns with safe incremental delivery.

### Signoff: claude - 2026-05-10
Status: ACCEPT
Notes: The events.jsonl-as-source-of-truth approach with polling and a narrow runner handle keeps this slice testable and avoids speculative scope.
