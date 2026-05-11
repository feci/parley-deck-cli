---
agent: hermes
idea: hitl-tui-questions
round: 1
date: 2026-05-11
---

## Summary
The idea proposes adding human-in-the-loop (HITL) question/answer support directly into the existing Go + Bubble Tea TUI for Parley Deck CLI runs. Questions triggered by agents or the runner must surface visibly during active execution, allow answers from the TUI, and include a non-interactive CLI fallback. Answers and questions must persist under the run directory for durability across TUI detachments. The slice must remain small, preserve `--auto` automatic mode with safe defaults, keep `--no-tui` functional, and include targeted tests.

## Proposed approach
- Extend the current TUI model to include a dedicated question pane or modal that activates on incoming questions while preserving live log streaming.
- Store each question/answer pair as a small JSON artifact inside the run directory (e.g., `qa/NN.json`) using atomic writes for durability.
- Implement a simple stdin/stdout fallback path for `--no-tui` and non-interactive terminals that records the same JSON artifacts.
- Add a `--auto` path that skips or auto-answers low-risk questions using documented defaults, falling back to human input only for flagged items.
- Expose a minimal TUI keybinding or input field for answers that triggers state transitions and persistence.
- Write focused unit tests covering persistence round-trips and TUI state machine transitions when questions arrive or are answered.

## Concerns / open questions
- How should the TUI visually distinguish high-priority vs low-risk questions without cluttering the live view?
- What exact schema and directory layout for persisted Q&A artifacts will best support later inspection and potential resume features (even if resume itself is out of scope)?
- Should the auto-default logic live inside the runner or be delegated to individual agents?
- Is there a clean way to surface questions without blocking the primary log stream in Bubble Tea?

## Risks
- Adding interactive input handling to a live TUI could introduce responsiveness issues or deadlock risks if not carefully isolated from the event loop.
- Overly broad persistence might accidentally capture sensitive data; clear scoping of what is stored is needed.
- Maintaining both TUI and CLI fallback paths risks divergent behavior if the fallback is not kept minimal and tested equally.
- Scope creep toward general chat or multi-turn protocols could occur if the initial durable Q&A mechanism is not strictly bounded.