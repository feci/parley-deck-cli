---
agent: codex
idea: live-run-tui
round: 2
date: 2026-05-10
responding-to: [claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

The round-01 files converged strongly. I now think the design should explicitly define two quit paths: `q`/`esc` detaches from the TUI while the runner continues, and `ctrl+c` cancels the run context and exits. I also agree with limiting log tailing to the selected agent for this slice to keep polling cost bounded and implementation small.

## Responses to others

### @claude - round-01

I agree with using `events.jsonl` plus log tails as the only live projection source, keeping HITL confirmation before the alt-screen starts, and preserving `--no-tui` as the synchronous path. Your `q` detach vs `ctrl+c` cancel proposal should be part of the final plan because it removes ambiguity and keeps durable state intact. I would keep `RunRoundOneAsync` optional: a small app-level goroutine wrapper around the existing runner is enough unless the implementation shows repeated boilerplate.

### @gemini - round-01

I agree with the `EventWatcher`/`RunState` framing and with handling a missing `events.jsonl` file gracefully at TUI startup. I prefer disk polling over an in-memory channel for this slice, even if a channel would be lower latency, because resume/re-attach later becomes a natural extension when the TUI can reconstruct state from disk alone. Debouncing can stay minimal: one event poll around 250 ms and one render tick around 1 s should be sufficient.

### @hermes - round-01

I agree that the implementation should stay lightweight and avoid token accounting or HITL interaction in this slice. For "raw event lines vs summarized view", I propose a summarized event list derived from event type, agent, and message fields rather than printing raw JSONL; raw files remain available on disk for debugging. To avoid false positives, round completion should only come from `round.completed` or `round.incomplete`, not inferred from all selected agents being terminal.

## New concerns / questions

The final plan should name the exact test targets so this does not grow into a broad TUI rewrite. I propose three targets: event-to-state reducer tests, log-tail helper tests, and a render snapshot for a hand-built live model. Integration with the runner can be verified by preserving existing runner tests and adding one small async/default-run path test only if the app layer can do it without brittle terminal coupling.

## Current proposal

Implement a live Bubble Tea TUI that polls the active run directory, reduces `events.jsonl` into per-agent status and round status, shows elapsed time and recent summarized events, and tails only the selected agent's stdout/stderr logs. Default `parley run` should launch the runner in a goroutine and open the TUI; `--no-tui` should keep the existing synchronous output. `q`/`esc` detaches from the view, `ctrl+c` cancels the run context, and final result printing happens after the TUI returns.
