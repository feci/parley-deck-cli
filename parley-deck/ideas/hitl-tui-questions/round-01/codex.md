---
agent: codex
idea: hitl-tui-questions
round: 1
date: 2026-05-11
---

## Summary

The first useful HITL slice should add a durable question/answer queue tied to each run, then surface that queue in the live TUI. I would avoid trying to stream answers into arbitrary already-running agent CLIs in this slice; instead, make the file protocol reliable and teach agent prompts how to create questions and observe answers when they choose to wait.

## Proposed approach

Add an internal `hitl` package that stores question records under `parley-deck/runs/<run-id>/questions/`. Use one JSON file per question so agent CLIs can write or read a single file without coordinating append locks. A question should include at least:

- `id`
- `agent`
- `prompt`
- `details`
- `default_answer`
- `risk`
- `status` (`open`, `answered`, `auto_answered`)
- `answer`
- `created_at`
- `answered_at`

The live TUI should poll questions alongside `events.jsonl`. It can show a compact "Questions" panel below or next to logs, with the selected open question and an answer input mode. Keep controls small: `?` or `a` enters answer mode, `enter` submits, `esc` cancels answer entry, and `tab`/`j`/`k` continue to navigate. On submit, the TUI updates the JSON file and appends a `hitl.answered` event.

Add `parley answer RUN_ID QUESTION_ID ANSWER` as a non-TUI fallback. It should find the run directory, update the question record, and append the same answer event. This makes `--no-tui` workflows viable and keeps future resume simple.

For `--auto`, the runner should be able to auto-answer low-risk questions only when a `default_answer` is present and `risk` is low. The first implementation can expose this as a helper and prompt contract rather than building deep consensus automation. In headless prompts, tell agents: if blocked, create a question file; in auto mode, include enough default/assumption text that the runner can answer low-risk questions or the agent can proceed with its stated assumption.

## Concerns / open questions

The hardest part is not the TUI input; it is making arbitrary external agent CLIs actually pause and consume answers. A file protocol can support that, but agents may choose not to wait. We should not overpromise "live injection" into child processes yet.

Question IDs should be deterministic enough to reference from CLI commands but unique enough to avoid collisions. A timestamp plus sanitized agent ID is enough for this slice.

## Risks

If the question schema is too loose, later resume/answer features will have to migrate data. Keep it explicit now. If the TUI steals too many keys from the live-run view, it can make run monitoring worse; answer mode should be clearly modal and easy to exit.
