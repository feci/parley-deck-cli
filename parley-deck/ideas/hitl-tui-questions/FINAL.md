---
idea: hitl-tui-questions
status: final
author: codex
consensus-date: 2026-05-11
participants: [codex, gemini, hermes]
---

## Final plan / specification

Implement the first human-in-the-loop question/answer slice for Parley Deck CLI.

The canonical Q&A state lives under each run directory:

```text
parley-deck/runs/<run-id>/questions/<question-id>.json
```

Each question JSON record should include:

- `id`
- `agent`
- `prompt`
- `details`
- `default_answer`
- `risk` (`low`, `normal`, `high`)
- `status` (`open`, `answered`, `auto_answered`)
- `answer`
- `created_at`
- `answered_at`

Question IDs should include the UTC timestamp, the source agent, and a short random suffix. Writes must be atomic enough for local file use: write a temporary file in the same directory and rename it into place.

Add an internal HITL package responsible for:

- generating question IDs;
- creating question files;
- listing run questions in stable order;
- answering a question;
- auto-answering a low-risk question only when `default_answer` is non-empty.

Questions and answers should also be mirrored into the run event stream:

- `hitl.question`
- `hitl.answered`

The JSON files remain canonical; events are a timeline/projection aid for the live TUI.

Add a CLI fallback:

```text
parley answer RUN_ID QUESTION_ID ANSWER
```

It finds the run directory, updates the matching question JSON file, and appends `hitl.answered` to `events.jsonl`. This replaces the previous placeholder implementation for `parley answer`.

Extend the live TUI:

- Poll questions in the active run directory alongside `events.jsonl`.
- Show a compact Questions panel with open and answered questions.
- Let the user select questions without breaking existing agent/log navigation.
- Add a modal answer-entry mode for the selected open question.
- On submit, persist the answer through the same HITL package and append the answer event.
- Keep existing live run controls: `j`/`k`/`tab`, `q`/`esc`, `ctrl+c`.

Auto mode:

- `--auto` may auto-answer only questions with `risk: low` and a non-empty `default_answer`.
- Auto answers must be persisted as `status: auto_answered`.
- High/normal risk questions remain open.

Agent contract:

- Agent prompts should mention the run questions directory and the JSON schema.
- This slice does not guarantee live answer injection into every running external CLI.
- Agents that want to block for user input can create a question file and poll it for an answered status.

## Verification plan

- Unit tests for HITL question create/list/answer/auto-answer persistence.
- Unit tests for atomic write/read behavior where practical.
- TUI tests for question projection and answer-entry state transitions.
- CLI test or app-level test for `parley answer RUN_ID QUESTION_ID ANSWER`.
- Existing tests must continue to pass with `go test ./...`.

## Out of scope

- Token accounting.
- Resume/re-attach UI.
- Provider-specific process injection.
- Timeout enforcement.
- Full answer validation rules.
- General chat UI.

## References

- Consensus: ./consensus.md
- Round 1: ./round-01/
- Round 2: ./round-02/
