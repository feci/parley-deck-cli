---
agent: hermes
idea: live-run-tui
round: 2
date: 2026-05-10
responding-to: [codex/round-01, claude/round-01, gemini/round-01]
---

## Position changes since round 1

Round-1 position remains mostly unchanged. I now align with the group's file-backed polling approach.

## Responses to others

### @codex - round-01
I agree with the core polling strategy using events.jsonl as source of truth and the minimal state machine (pending/running/finished/failed/skipped). The suggestion to keep log previews bounded and to retain --no-tui behavior unchanged is sound. One concrete counter-proposal: instead of a generic "done channel" plus separate result buffer, expose a single handle type that implements both Wait() and Done() so callers cannot accidentally forget to drain results; this reduces the risk of the TUI exiting before final summary printing.

### @claude - round-01
Claude's detailed reduction table and quit-semantics distinction (q = detach vs ctrl+c = cancel) align closely with the prompt's "quit without killing durable state" requirement. I support the pre-TUI confirmation flow and the decision to defer resume. Minor disagreement on test scope: the proposed "reduce events to state" unit test is valuable, but it should also assert that unknown agents (not in the current round) render as unknown rather than pending; adding that case prevents later drift when more agents join.

### @gemini - round-01
Gemini's high-level EventWatcher and three-pane layout are compatible with the other two proposals. I agree that a circular buffer for recent events is a clean way to keep the activity log bounded. Counter-proposal on log streaming: polling stdout/stderr tails (as Codex and Claude also suggest) is preferable to an in-memory shared buffer for this slice; the latter would couple the runner and TUI more tightly than the prompt's "observable runner" goal allows. Keep the watcher strictly file-based.

## New concerns / questions
- If multiple agents finish within the same 250 ms poll window, the TUI must still produce a stable sorted view (e.g., by agent name) to avoid visual reordering on every tick.
- The prompt lists "whether the current round completed or is incomplete" as required display data; none of the round-1 proposals explicitly surface a top-level round status banner derived from round.completed / round.incomplete events. This should be added to the header.

## Current proposal
Adopt the common polling + events.jsonl reduction approach outlined by Codex and Claude, with the layout and quit semantics from Claude, and the circular event buffer from Gemini. Implement a small RunLive entry point that receives a runner handle exposing RunDir and Done(). The TUI model will maintain per-agent maps populated only from observed events, cap all file reads to 4 KiB, and render a single static screen until q or ctrl+c. --no-tui path remains untouched. This keeps the slice minimal while satisfying every requirement in 00-prompt.md.