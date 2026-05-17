---
idea: tui-workspace-sessions
author: user
created: 2026-05-17
participants: [codex, claude, gemini, hermes]
roles:
  codex: implementation and codebase-fit lens
  claude: interaction model and correctness lens
  gemini: state model and concurrency lens
  hermes: operations and recovery lens
status: final
---

## Problem / idea

Upgrade the Parley Deck CLI TUI from a mostly static workspace/run view into a workspace session console.

The user wants to start and supervise Parley Deck ideas directly from the TUI:

- Start new ideas/runs from the TUI.
- See status for what is currently running.
- Start multiple ideas in parallel.
- Clearly see which idea/run needs user action.
- See what agents are doing while work is active.
- Show a time-series/event stream as events happen.
- Provide a per-idea view where each agent has its own view/tab.
- Support multiple running ideas, where each idea can contain per-agent views.
- Persist enough session information under a user-local directory such as `~/.parley-deck` so sessions can be found again and continued later.
- Allow continuing/resuming previously started sessions.

## Current context

- Existing `internal/tui/app.go` renders a workspace dashboard with idea and agent panels.
- Existing `internal/tui/live.go` renders a live view for a single run, including agent table, recent events, HITL questions, answer mode, and selected-agent log preview.
- Existing `internal/runstate` loads run summaries from `parley-deck/runs/<run-id>/events.jsonl` and derives agent state, recent events, open HITL questions, and terminal/liveness state.
- Existing `parley run` can create one run and optionally open a live TUI for that run.
- Existing `parley resume` can resolve a run ID or idea slug and re-open a live run view.
- Current repository transport is `github-pr`.
- The implementation should fit the existing Go + Bubble Tea/Bubbles/Lip Gloss stack.

## Constraints

- Keep this as a pragmatic MVP that can be shipped and tested locally.
- Preserve existing command behavior for `parley run`, `parley resume`, and `parley tui`.
- Use existing run artifacts (`parley-deck/runs/<run-id>/events.jsonl`, logs, questions) as the source of truth where possible.
- Add user-local session metadata under `~/.parley-deck` only for cross-workspace discovery/resume convenience; do not make it the canonical protocol artifact.
- Do not store secrets or large logs in the user-local session metadata.
- Multiple concurrently running ideas should be visible in one TUI.
- The UI must make user-action-required states obvious.
- The UI must expose per-agent activity and logs for the selected idea/run.
- Tests should cover the state model and key TUI rendering/interaction behavior without requiring real agent CLIs.

## Non-goals

- No hosted service or background daemon in this slice unless all participants agree it is unavoidable.
- No terminal automation of interactive third-party CLIs.
- No replacement of `parley-deck/` canonical protocol files.
- No full project management system beyond Parley Deck run/session supervision.
- No complex persistence database; a small JSON file under `~/.parley-deck` is preferred if sufficient.
