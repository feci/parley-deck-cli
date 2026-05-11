---
idea: hitl-tui-questions
drafted-by: codex
date: 2026-05-11
---

## Agreed decisions

- Use one JSON file per question as the canonical source of truth under `parley-deck/runs/<run-id>/questions/`.
- Mirror Q&A changes into the existing event stream with `hitl.question` and `hitl.answered` events so the live TUI timeline can react without making `events.jsonl` the only storage.
- Add a small internal HITL package for question ID generation, atomic persistence, listing, and answering.
- Add `parley answer RUN_ID QUESTION_ID ANSWER` as the non-TUI fallback.
- Extend the live TUI with a Questions panel and a clear answer-entry mode. Keep the existing run monitoring controls intact.
- Use conservative auto mode: only questions with `risk: low` and a non-empty `default_answer` can be auto-answered, and those answers must be marked `auto_answered`.
- Do not promise live answer injection into arbitrary already-running external CLIs. This slice provides the durable file contract; agent prompts may instruct agents to create question files and poll for answers.

## Agreed trade-offs

- Per-question JSON files are more direct for current question status, CLI fallback, and external-agent polling than storing only in `events.jsonl`.
- Event mirroring preserves the useful timeline and TUI reactivity from Gemini's round-01 proposal.
- Runner-wide blocking, stdin fallback inside `--no-tui`, answer validation rules, and timeout enforcement are deferred.
- Question IDs should include time plus a short random suffix to avoid rapid-collision risk.

## Open items deferred to implementation

- A future slice can make agent-side waiting more robust and provider-specific.
- A future slice can add validation metadata, enforced timeouts, and richer priority display.
- A future slice can add resume/re-attach UI for existing runs.

## Signoffs

### Signoff: codex - 2026-05-11
Status: ACCEPT
Notes: The consensus keeps HITL durable and usable without overpromising process-level answer injection.

### Signoff: gemini - 2026-05-11
Status: ACCEPT
Notes: The proposed hybrid approach of per-question JSON files and event mirroring ensures both durability and real-time UI reactivity.

### Signoff: hermes - 2026-05-11
Status: ACCEPT
Notes: Hybrid JSON+events approach balances durability with TUI reactivity while staying conservative on scope.
