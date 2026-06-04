---
agent: hermes
idea: tui-interactivity-overhaul
review-round: 1
date: 2026-06-04
---

## Summary
Adversarial review of slices 1-4 vs FINAL.md. Segment reset relies on barrier not per-event tagging (ACP gap); composer queues but never executes (misleading label); viewport follow/rotate edge cases untested in live rapid output; steer CLI arg handling fragile; mode enum migration clean but state invariants thin.

## Findings
### [CRITICAL] Composer over-promises execution that slice 4 never delivers
Label and help say "queued, runs on next attempt" (live.go:1120+, steer.Submit) but queued steers are only persisted as events; no consumer launches new_attempt (slice 5 deferred per IMPLEMENTATION.md). User expects action, gets dead letter. Concrete fix: change all UI strings to "queued (exec deferred to future slice)" or remove compose until slice 5.
### [MAJOR] Segment reset only verified in tests, not real continue/resume paths
runstate.ProjectEvents + runner appendSegmentStarted fix badge in segment_test.go but ACP launch (acp.go) skips tagging; stale [FINISHED] can survive if segment_started races an in-flight agent.finished or non-targeted agent inherits via legacy unsegmented mapping. Concrete fix: add live continue test that emits real events.jsonl then asserts no stale badge on targeted agent.
### [MAJOR] Viewport follow mode and boundary clamp lack rapid-output/rotate coverage
live.go:840-980 implements offset-incremental + follow but no test for: log rotate mid-scroll (new segment), focused agent with zero stdout yet, or >1k lines/sec burst causing jumpy scroll. Can show empty/stale pane. Concrete fix: add live_test.go cases that simulate fs write bursts + rotation and assert follow clamps to last N visible lines.
### [MINOR] `parley steer` arg parsing incomplete for edge cases
steer.go:140+ uses flag parsing but missing handling for: empty TEXT after --, -- placed before RUN, or --agent not in participants list. Produces confusing error or silent no-op. Concrete fix: explicit validation + usage error before Submit.
### [NIT] Mode enum migration comment drift
live.go:42 comment still references old answerMode/focus booleans even after collapse; minor source of future confusion.

## Open questions
- Does any code path actually read the queued steer.* events to trigger new_attempt?
- Are there cross-process append races on events.jsonl when TUI + steer CLI run concurrently?