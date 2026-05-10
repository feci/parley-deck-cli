---
idea: live-run-tui
status: final
author: codex
consensus-date: 2026-05-10
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

Implement a live TUI for `parley run` that observes an active runner execution through the run directory. The default `parley run` flow should confirm launch, start the runner in the background, and open a Bubble Tea view. `parley run --no-tui` must keep the existing synchronous behavior and final text summary.

The TUI must use `events.jsonl` as its source of truth. It should poll the file with a retained byte offset, treat a missing file as an empty event stream, and reduce events into a display model. It must not depend on a separate in-memory event bus for this slice.

Agent states are:

- `pending`: selected participant has no start event yet.
- `running`: latest relevant event is `agent.started`.
- `finished`: latest relevant event is `agent.finished`.
- `failed`: latest relevant event is `agent.failed`.
- `skipped`: latest relevant event is `agent.skipped`.
- `unknown`: agent is not part of the current selected round or cannot be mapped from run metadata.

Round status must be displayed from `round.completed` or `round.incomplete` events only. The TUI must not infer round completion from all agents reaching terminal states.

The layout should stay compact and operational:

- Header: idea, phase/round, run ID, and round status.
- Agent table: selected marker, agent ID, state, elapsed time, and latest event summary.
- Event pane: bounded list of summarized recent events.
- Log pane: bounded stdout/stderr tail for the selected agent only.
- Footer: `j/k` and `tab` select agent, `q`/`esc` detach, `ctrl+c` cancel.

Log preview reads should be capped to a small byte window such as 4 KiB per stream, should prefer whole lines, and should strip terminal control noise where practical. Reading only the selected agent's logs keeps polling cost bounded.

Runner integration should be minimal. Use either an app-layer goroutine wrapper or a small runner handle that exposes the run directory and completion through `Done()`/`Wait()`. The important behavior is that final result printing happens after the TUI exits and the terminal has been restored.

Quit semantics are fixed for this slice:

- `q` or `esc`: leave the TUI while the runner continues.
- `ctrl+c`: cancel the run context and exit.

## Verification plan

- Unit test event reduction from fixed event sequences, including `pending`, `unknown`, all terminal agent states, and round status.
- Unit test log-tail behavior for capped reads and partial-line handling.
- Add render tests or view snapshots for representative terminal sizes.
- Preserve existing runner tests and add focused coverage only where the async/default TUI path introduces new lifecycle behavior.
- Run `go test ./...` before opening implementation review.

## Out of scope

- Token accounting.
- HITL question/answer entry inside the TUI.
- Resume or re-attach to an existing run.
- GitHub/GitLab API automation.
- Packaging and installer work.

## References

- Consensus: ./consensus.md
- Round 1: ./round-01/
- Round 2: ./round-02/
