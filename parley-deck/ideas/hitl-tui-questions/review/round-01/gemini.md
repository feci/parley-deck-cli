---
agent: gemini
idea: hitl-tui-questions
review-round: 1
date: 2026-05-11
reviewed-commit: c4d6f5fd203450f7d2843a17354c5e9447a4444b
---

## Summary
The implementation strictly follows the `FINAL.md` specification and `IMPLEMENTATION.md` plan. It introduces a robust HITL package (`internal/hitl`) with atomic JSON persistence for question records, event mirroring to the run log, and a functional CLI fallback via `parley answer`. The live TUI is extended with a new Questions panel and a modal answer-entry mode that integrates seamlessly with the existing agent and log navigation. Unit tests cover storage, auto-answering logic, and CLI command behavior.

## Findings
### [MINOR] Missing HITL event summaries in TUI
The `summarizeEvent` function in `internal/tui/live.go` does not explicitly handle `hitl.question` or `hitl.answered` event types. As a result, these events appear in the "Latest events" pane with only the agent ID and the event type, but no descriptive text (like the question prompt or status). While the dedicated Questions pane provides the necessary detail, mirroring some context into the timeline would improve the narrative flow of the run.

### [MINOR] Brittle TUI list error handling
In `internal/hitl/hitl.go`, the `List()` method returns an error if any JSON file in the questions directory is malformed. In the TUI (`internal/tui/live.go`), this causes `m.answerErr` to be set, potentially obscuring the entire questions list if a single file is corrupt or partially written by an agent. A more resilient approach would be to log/skip malformed files during listing.

### [NIT] Unicode-unsafe backspace in TUI
The `updateAnswerMode` in `internal/tui/live.go` uses `m.answerText[:len(m.answerText)-1]` for backspace handling. This is unsafe for multi-byte Unicode characters and should ideally use rune slicing or the `textinput` component from the bubbles library.

### [NIT] Agent slugification in Question IDs
The `slugify` function in `internal/hitl/hitl.go` can return an empty string if an agent ID consists only of symbols (e.g., "@#!"), leading to IDs like `20260511...Z--randomsuffix`. While unlikely with current agent naming conventions, a fallback to "agent" or "unknown" would be safer.

## Open questions
None. The implementation is well-structured and aligns with the project's incremental design philosophy.
