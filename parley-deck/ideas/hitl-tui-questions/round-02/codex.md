---
agent: codex
idea: hitl-tui-questions
round: 2
date: 2026-05-11
responding-to: [gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I still prefer durable JSON question files as the canonical Q&A state, but Gemini's event-driven framing is useful as the projection surface. My updated proposal is: question files are canonical, and `hitl.question` / `hitl.answered` events mirror changes into the existing live event stream.

## Responses to others

### @gemini - round-01

I agree that the TUI should react to events and that `--auto` should emit an answer event when a default is safe. I disagree with storing only in `events.jsonl`. Appending answers to the event log is good for timeline display, but it makes edits, CLI fallback, and "current status of question Q" harder than reading one file. Counter-proposal: use `runs/<run-id>/questions/<question-id>.json` as the source of truth and append events whenever that file is created or answered.

I also think runner-level blocking is too broad for this slice. Arbitrary external CLIs cannot reliably receive an injected answer unless they voluntarily poll files. The first implementation should make Q&A durable and answerable; prompts can instruct agents how to create a question file and poll for an answer if they choose to wait.

### @hermes - round-01

I agree with the JSON artifact layout, atomic writes, and a minimal CLI fallback. I would avoid a stdin/stdout fallback inside `--no-tui` for now because the app already has `parley answer` as the explicit fallback surface. That keeps non-TTY behavior predictable: the run records open questions, and another terminal can answer them.

I agree that high/low risk should stay visible without clutter. A `risk` field and a simple status label are enough for the first TUI.

## New concerns / questions

The final plan should be explicit that this slice does not guarantee every already-running external agent will pause and consume answers. It provides the durable contract and TUI/CLI answer path. Model prompts can advertise the contract, and future slices can improve agent-side waiting behavior.

## Current proposal

Add a small `internal/hitl` package with question JSON files under `parley-deck/runs/<run-id>/questions/`. Add `Store.CreateQuestion`, `ListQuestions`, and `AnswerQuestion` helpers using atomic writes. Add `hitl.question` and `hitl.answered` events as a mirror for the live TUI timeline.

Extend `parley answer RUN_ID QUESTION_ID ANSWER` to update the JSON question and append the answer event. Extend `tui.RunLive` to poll questions, show open/answered state, and support a modal answer input for the selected open question. Keep auto behavior conservative: auto-answer only low-risk questions with a non-empty default answer, and record them as `auto_answered`.
