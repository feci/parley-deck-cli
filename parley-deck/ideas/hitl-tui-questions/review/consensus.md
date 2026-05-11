---
idea: hitl-tui-questions
review-cycle: 1
drafted-by: codex
date: 2026-05-11
reviewed-commit: 93741c936f78f36d120ed10cd19f7aac3bc86a87
---

## Agreed fixes

- From gemini/review/round-01 [MINOR] Missing HITL event summaries in TUI: add explicit `hitl.question` and `hitl.answered` summaries for the Latest events pane.
- From gemini/review/round-01 [NIT] Unicode-unsafe backspace in TUI: remove whole runes in answer-entry backspace handling.
- From gemini/review/round-01 [NIT] Agent slugification in Question IDs: fall back to `agent` when slugification produces an empty value.
- From hermes/review/round-01 [NIT] TUI question panel uses stable sort but lacks explicit test for ordering: add focused ordering coverage for HITL question listing.
- From hermes/review/round-01 [MINOR] CLI `parley answer` help text could better document the --dir flag behavior: clarify usage with `ANSWER...` to show multi-word answers are accepted.

## Deferred follow-ups

- Hermes open question about agent polling timeouts is deferred to a future agent-contract or resume/re-attach slice. This implementation keeps the current slice scoped to durable Q&A records and polling support.

## Dismissed findings

- Gemini [MINOR] Brittle TUI list error handling: keep malformed canonical question JSON as a visible error. Temp files are ignored by suffix, and silently skipping corrupt canonical records would hide a state problem the user should see.
- Hermes [MINOR] Question ID generation entropy: no change for this slice. IDs already combine nanosecond UTC timestamp, agent slug, and a random suffix; collision risk is acceptable for local file use.
- Hermes open question about distinct auto-answer display: already covered by the Questions panel status column, which shows `auto_answered`.

## Signoffs

### Signoff: codex — 2026-05-11
Status: ✅ ACCEPT
Notes: Implementer accepts the agreed fixes and will apply them in fix-up cycle 1.

### Signoff: gemini — 2026-05-11
Status: ✅ ACCEPT
Notes: The consensus accurately reflects my review findings; agreed fixes and dismissals are acceptable for this slice.

### Signoff: hermes — 2026-05-11
Status: ✅ ACCEPT
Notes: The consensus captures my review findings; agreed fixes and deferred follow-ups are acceptable for this slice.

## Review cycle 2 final consensus

reviewed-commit: 63b62098bcdcf92a1ce9b1d1b226e09d90d0ad55

### Review summary

Fix-up cycle 1 applied the agreed fixes from review cycle 1. Gemini and Hermes both completed review round 2 with no findings and no open questions.

### Agreed fixes

None for review cycle 2.

### Deferred follow-ups

- Agent polling timeout guidance remains deferred to a future agent-contract or resume/re-attach slice.

### Dismissed findings

- No new findings were dismissed in review cycle 2.

### Signoffs

### Signoff: codex — 2026-05-11
Status: ✅ ACCEPT
Notes: Review round 2 has zero agreed fixes. Implementation can be marked complete and merged.

### Signoff: gemini — 2026-05-11
Status: ✅ ACCEPT
Notes: Review round 2 has no findings; the implementation can be marked complete and merged.
### Signoff: hermes — 2026-05-11
Status: ✅ ACCEPT
Notes: Review round 2 has no findings; the implementation can be marked complete and merged.
