---
idea: hitl-tui-questions
status: implemented
implementer: codex
started: 2026-05-11
completed: 2026-05-11
branch: parley-deck-cli#feature/hitl-tui-questions
head-commit: 4d65b7228a5c0e4d75b29e245aa9d5e2b889e657
design-pr: https://github.com/feci/parley-deck-cli/pull/3
implementation-pr: https://github.com/feci/parley-deck-cli/pull/4
---

## Summary of work

Implemented the HITL TUI Q&A slice from `FINAL.md`:

- Added `internal/hitl` for canonical question JSON records under `parley-deck/runs/<run-id>/questions/`.
- Added stable list/create/answer/auto-answer behavior with local atomic temp-file writes.
- Mirrored HITL changes into the run event stream with `hitl.question` and `hitl.answered`.
- Added the CLI fallback `parley answer [--dir DIR] RUN_ID QUESTION_ID ANSWER`.
- Extended the live TUI with a Questions panel, question selection, and answer-entry mode.
- Added `--auto` low-risk default auto-answer support during active runs.
- Updated runner prompts so agents know where to write question files and how to poll for answers.
- Added unit coverage for HITL storage, CLI answering, prompt content, and TUI answer mode.

## Deviations from FINAL.md

None known.

## Notes for reviewers

- Agent-authored question files are canonical and are displayed by direct polling from the TUI. `hitl.question` events are emitted when questions are created through the internal HITL package; manually authored files still appear in the Questions panel.
- This slice intentionally does not inject answers into arbitrary running CLIs. Agents can poll the question JSON file when they choose to wait for human input.
- Verification run: `GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
