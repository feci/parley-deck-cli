---
idea: live-run-tui
review-round: 1
drafted-by: codex
date: 2026-05-10
---

## Review summary

Review round 1 produced two blocking implementation bugs and one blocking UI integrity issue:

- Skipped agents without a preceding `agent.started` event displayed an invalid multi-year elapsed duration.
- Unknown agents could become normal running/finished agents after a second event, contradicting FINAL.md.
- Agent table alignment could break because styled ANSI strings were padded with `fmt`.

Hermes approved without findings. Claude and Gemini requested changes.

## Agreed fixes

- Set skipped-agent duration from `duration_ms` when present, otherwise keep it at zero.
- Keep `unknown` sticky for agents outside the selected participant set.
- Remove pre-padding ANSI styling from the agent state table and align headers/rows with plain fixed-width columns.
- Derive the displayed round label from the idea status instead of hardcoding `round-01`.
- Render terminal zero-duration states as `0s` instead of `-`.
- Extend reducer tests to cover skipped-without-start and unknown-after-multiple-events.

## Deferred follow-ups

- JSON parse resync after a malformed complete event line.
- Dedicated `run.failed` status projection beyond showing the event in the latest-events pane.
- Wider render snapshot coverage across terminal sizes.
- Footer mention for `shift+tab`.
- Optional final-state pause after runner completion.

## Signoffs

### Signoff: codex - 2026-05-10
Status: ACCEPT
Notes: Agreed fixes have been applied in fix-up cycle 1; deferred follow-ups are not blockers for this slice.

### Signoff: gemini - 2026-05-10
Status: ACCEPT
Notes: Implementation correctly addresses blocking issues identified in round 1, including sticky unknown states and terminal zero-duration formatting.

### Signoff: hermes - 2026-05-10
Status: ACCEPT
Notes: Fix-up cycle resolves all prior blockers with correct skipped duration, sticky unknown, and plain column alignment; tests cover the edge cases.

### Signoff: claude - 2026-05-10
Status: ACCEPT
Notes: Verified the six agreed fixes in live.go and live_test.go — skipped duration falls back to 0, unknown stays sticky across later events, table uses plain fixed-width columns, round label is derived from Idea.Status, terminal zero-duration renders as `0s`, and reducer tests cover skipped-without-start and unknown-after-multiple-events.
