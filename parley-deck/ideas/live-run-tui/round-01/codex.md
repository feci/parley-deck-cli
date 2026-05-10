---
agent: codex
idea: live-run-tui
round: 1
date: 2026-05-10
---

## Summary

The smallest useful slice is to run the existing `runner.RunRoundOne` in a goroutine and let the Bubble Tea program poll the run directory for `events.jsonl` and per-agent log files. This avoids changing the runner contract while giving the user live visibility into agent status, elapsed time, failures, and recent activity.

## Proposed approach

Add a `tui.RunLive(...)` entry point that receives workspace status, discovered agents, the active `IdeaStatus`, the run ID, the run directory, and a done channel. Its model should poll on a short ticker, load `events.jsonl` through the existing store reader, derive per-agent display state, and read bounded tails from each agent's stdout/stderr logs. `parley run` should keep `--no-tui` synchronous, but for default TUI mode it should start the runner in a goroutine, launch `RunLive`, then print final results after the TUI exits or the round completes.

The initial derived states can stay simple:

- no event yet: `pending`;
- `agent.started`: `running`;
- `agent.finished` with `artifact_ok=true`: `finished`;
- `agent.failed`: `failed`;
- `agent.skipped`: `skipped`;
- round event: `round.completed` or `round.incomplete`.

The TUI layout should reuse the current dashboard style but replace the static agent discovery panel with a runtime table: agent, state, elapsed, artifact, and log path. A lower panel should show the latest events and a bounded log preview for the selected/first active agent. Keyboard controls can stay minimal: `q` quits the view, `tab` cycles the selected agent if time permits.

## Concerns / open questions

Polling `events.jsonl` every few hundred milliseconds is simpler than wiring a subscription channel, but repeated file reads can become noisy later. That trade-off is acceptable for the first live UI because event volume is low and the durable event log remains the source of truth.

The current runner only emits coarse events, so "writing-artifact" cannot be known reliably yet. The TUI should not invent that state. It can show `running` until an artifact appears or a terminal event arrives.

## Risks

The main risk is blocking forever if the TUI exits before the runner finishes or if the runner finishes before the TUI is ready. Keep the done channel one-way, buffer the runner results in `runTask`, and make sure final result printing happens after the TUI returns. Another risk is reading very large logs; cap log preview reads to a small byte window.
