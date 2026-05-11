---
agent: hermes
idea: hitl-tui-questions
round: 2
date: 2026-05-11
responding-to: [codex/round-01, gemini/round-01]
---

## Position changes since round 1

My round-1 position remains mostly aligned. I now favor file-backed queue plus event mirrors.

## Responses to others
### @codex - round-01

I agree with the hitl package and one-JSON-per-question design under runs/<run-id>/questions/. This matches the "small and verifiable" constraint and supports the non-TUI fallback cleanly. The proposed keybindings and modal answer mode are practical. Counter-proposal on one point: make question ID generation use a short UUID suffix rather than only timestamp+agent to reduce collision risk even in rapid succession.

### @gemini - round-01

The event-driven approach with QuestionRequested/AnswerProvided is elegant for future extensibility but adds unnecessary complexity for round 2. Appending everything to events.jsonl risks bloat and harder partial reads by external agents. I prefer codex's explicit per-question files; they can still emit events if desired. On runner blocking: we should not implement full cross-agent blocking yet—keep it advisory via the question status for this slice.

## New concerns / questions

- How do we surface multiple pending questions without overwhelming the TUI? A simple list + detail pane seems sufficient.
- Should the schema include an optional `timeout` field even if not enforced yet?
- For --auto, clarify the exact risk threshold logic before implementation.

## Current proposal

Adopt codex's file-per-question storage and TUI polling model, augmented with a lightweight hitl package. Add `parley answer` CLI. For agents, document a simple "create question JSON then poll status" contract. --auto answers only low-risk questions with defaults. Include tests for persistence and TUI modal transitions. Keep scope strictly to durable Q&A visibility and input.