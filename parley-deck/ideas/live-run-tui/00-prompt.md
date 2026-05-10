---
idea: live-run-tui
author: user
created: 2026-05-10
participants: [codex, claude, gemini, hermes]
status: final
---

## Problem / idea

Implement the next `parley-deck-cli` slice: a live TUI view for active agent runs.

The user must be able to start a task and see, while agents are still running:

- the active idea, phase, and round;
- per-agent state such as pending, running, finished, failed, skipped, or unknown;
- per-agent elapsed time;
- the latest protocol/runtime events from `events.jsonl`;
- stdout/stderr log previews for active or failed agents;
- whether the current round completed or is incomplete.

This slice should make the existing runner observable. It should not attempt full token accounting yet.

## Constraints

- Keep implementation small and verifiable.
- Files under `parley-deck/` remain canonical.
- Runtime events come from `parley-deck/runs/<run-id>/events.jsonl`.
- Existing `parley run --no-tui` behavior must remain intact.
- The default `parley run` should open a live TUI while the runner goroutine is active.
- The TUI must remain usable in a normal terminal and should quit without killing durable state.
- Token accounting is explicitly deferred to a later slice.

## Non-goals

- Do not implement `parley answer` or HITL question capture in this slice.
- Do not implement `parley resume` in this slice.
- Do not implement GitHub/GitLab API automation in this slice.
- Do not add npm/Homebrew/release automation in this slice.
