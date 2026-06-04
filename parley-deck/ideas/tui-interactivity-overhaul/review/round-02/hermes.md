---
agent: hermes
idea: tui-interactivity-overhaul
review-round: 2
date: 2026-06-04
responding-to: [claude/IMPLEMENTATION fix-up-cycle-1, codex/review-round-01]
---

## Position changes since prior review round
All CRITICAL/MAJOR/MINOR findings from round-01 resolved by the fix-up (AF1-AF9). No remaining agreed fixes. Composer wording now honest; viewport bounded and non-fragmenting; steer queuing useful as durable intent capture (honestly labeled "not auto-run yet").

## Responses to other reviewers
### @codex
Agree on AF2/AF4 resolution: byte cap + completeLinesFrom eliminates jumpy/oversized/partial-line risks. No regressions from the fix-up itself.

## Updated findings
### [CRITICAL] Composer over-promise
RESOLVED (AF1). Wording in live.go:824+ and app.go:1111+ now explicitly states "recorded ... (queued; auto-exec not wired yet)" and "enter record". CLI/TUI both honest; user cannot expect dead-letter execution.

### [MAJOR] Segment reset verification
RESOLVED (AF6+AF7). ACP tagging added; LoadRunAt integration test confirms badge unsticks on real events.jsonl path.

### [MAJOR] Viewport rapid-output/rotation
RESOLVED (AF2+AF4+AF9). readAppendedLines caps to maxFocusBytes, capFocusLines enforces dual budget, partial lines merged, preview height-derived. No unbounded growth or fragmentation.

### [MINOR] parley steer arg handling
RESOLVED (AF8). Empty-text error + non-participant warning present.

### [NIT] Mode enum comment
RESOLVED (dismissed per consensus).

## Open questions
None. Zero agreed fixes remain; Phase 7 can finalize. Queuing remains net-useful (captures intent today) because it is now honestly communicated as deferred. No new findings.