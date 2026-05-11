---
idea: hitl-tui-questions
author: user
created: 2026-05-11
participants: [codex, gemini, hermes]
status: final
---

## Problem / idea

Implement the next Parley Deck CLI slice: human-in-the-loop questions and answers in the live TUI. When an agent or the runner needs user input, the user should be able to see the question while a run is active and answer it from the TUI. There must also be a simple CLI fallback so answers can be recorded without an interactive terminal.

## Constraints

- Transport is `github-pr`; canonical artifacts are files under `parley-deck/ideas/hitl-tui-questions/`.
- The implementation must build on the existing Go + Bubble Tea live TUI.
- Keep the slice small and verifiable.
- Preserve fully automatic mode: `--auto` should continue through low-risk questions using defaults or documented assumptions where possible.
- Preserve HITL default behavior: without `--auto`, questions should remain visible and answerable by the human.
- Persist questions and answers durably under the run directory so they survive TUI detach and can be inspected later.
- Keep `parley run --no-tui` usable; it should not depend on terminal UI input.
- Include focused tests for question/answer persistence and TUI state transitions.
- Claude CLI authentication failed with HTTP 401 during round-01 launch on 2026-05-11, so quorum for this idea is `codex`, `gemini`, and `hermes`.

## Non-goals

- No token accounting.
- No resume or re-attach implementation beyond making Q&A data durable.
- No GitHub/GitLab API automation.
- No general chat UI.
- No complex multi-turn agent protocol beyond the first durable Q&A mechanism.
