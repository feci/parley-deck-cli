---
from: user
to: all
idea: parley-deck-cli-plan
phase: post-final-user-direction
blocking: no
date: 2026-05-10
---

## User direction

The user answered the final planning questions with these priorities:

> 1 self contained would be great, 2 use Go + Bubble Tea/Bubbles/Lip Gloss, 3 polished TUI, 4 yes continue in auto mode through low risk questions, 5 yes, 6 codex, claude, gemini, hermes agent, 7 all three including windows, 8 local dir and github/gitlab

## Interpreted implementation decisions

- Prioritize a self-contained executable.
- Use Go for the implementation.
- Use Bubble Tea, Bubbles, and Lip Gloss for the TUI.
- Build a polished TUI, not only the smallest reliable dashboard.
- In automatic mode, continue through low-risk model questions instead of pausing every time.
- Best-effort token tracking is acceptable; explicit `unknown` values are acceptable when a CLI does not expose usage.
- First-class agent adapters: Codex, Claude, Gemini, and Hermes Agent.
- Release gate includes macOS, Linux, and Windows.
- Product transport support should include local directory plus GitHub PR and GitLab MR workflows.

## Notes

The current coordination transport for this Parley Deck workspace remains `local-dir` per `parley-deck/COOPERATION.md`. The user's answer about GitHub/GitLab is interpreted as a product requirement for `parley-deck-cli`, not as a request to switch this current protocol transport.

Hermes Agent should be included as a first-class adapter target, but no `hermes`, `hermes-agent`, or `hermesagent` executable was found on PATH during this session. Implementation should make the adapter configurable by command path until the exact CLI invocation is known.
