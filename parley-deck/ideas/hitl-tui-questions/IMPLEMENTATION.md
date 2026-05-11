---
idea: hitl-tui-questions
status: complete
implementer: codex
started: 2026-05-11
completed: 2026-05-11
branch: parley-deck-cli#feature/hitl-tui-questions
head-commit: 6ff4a4434c10f815694a86c4a095cef46a9559b2
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

## Fix-up cycle 1

status: complete
completed: 2026-05-11
head-commit: 6ff4a4434c10f815694a86c4a095cef46a9559b2

### Fixes applied

- Added explicit `hitl.question` and `hitl.answered` event summaries in the live TUI.
- Changed answer-entry backspace handling to remove a whole rune.
- Added `agent` fallback when question ID slugification would otherwise be empty.
- Added focused HITL question list ordering coverage.
- Clarified `parley answer` usage as accepting `ANSWER...` words.

### Deviations from agreed fixes

None.

Verification run: `GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`

## Completion

Review round 2 accepted fix-up cycle 1 with no remaining findings. The final review consensus lists zero agreed fixes for review cycle 2, so this implementation is complete and ready to merge.
