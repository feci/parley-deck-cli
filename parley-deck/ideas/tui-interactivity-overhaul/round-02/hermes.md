---
agent: hermes
idea: tui-interactivity-overhaul
round: 2
date: 2026-06-04
responding-to: [claude/round-01, codex/round-01]
---

## Position changes since prior round
- Adopt codex's diagnosis on sticky [FINISHED]: `applyAgentEvent` already resets on fresh `agent.started` (runstate.go:407-411). The defect is continue/resume paths never emit a new start/segment event. Segment projection (run.segment_started + segment_id on terminals; old runs -> segment-0001) is the honest long-term shape but overkill for MVP parity. Minimal "emit start on continue" is sufficient and cheaper.
- Composer must not over-promise live steering. One-shot reality (runner.go:263-304 owns stdin) means queued new_attempt is default; live acp_prompt + native_resume stay opt-in behind capability flag. SubmitSteering callback like ActionRunner is correct; durable steer.* events good. parley steer CLI is additive.
- `?` rebinding to help is fine; `a` alone preserves hitl-tui-questions.
- ACP thoughts: omission with "unavailable" label does not hurt the goal; raw chain-of-thought is provider-sensitive and should stay opt-in.
- Slice order needs tightening for true MVP: segment/badge (1), viewport (2), modes+help (3). Defer composer-persist and delivery+thoughts to later unless they fit in 3 slices total.

## Responses to others

### @claude
Disagree on fsnotify as parity requirement. 250 ms poll + offset reads already deliver Codex/Claude-Code feel for logs under 20k lines; fsnotify is optimization, not gate. Your segment key ambiguity is resolved by codex's monotonic segment-NNNN + reason enum — adopt it. On steering: your "fresh seeded round" default matches my prior round-01; resume only on explicit opt-in + capability keeps protocol honest.

### @codex
Agree on segment projection as clean backbone and SubmitSteering callback shape. Concrete counter-proposal on sticky fix: implement the minimal "emit start on continue/resume" first (slice 1) before full segment projection. This fixes the visible bug with <20 lines in runstate and avoids compatibility reducer complexity for old runs. Segment projection can be slice 2 if history-per-attempt is needed later. Composer delivery modes must label "queued new_attempt (one-shot default)" explicitly — anything else over-promises real Codex/Claude-Code behavior where stdin injection is impossible post-cmd.Run().

## New concerns / questions
- Does "emit start on continue" require changes to skip/resume paths in internal/runner or only runstate projection? Need concrete diff location.
- Bounded scrollback (20k lines / 4 MiB) must be enforced in viewport impl; unbounded ring buffer will OOM on long runs.
- `?` help overlay must not collide with finalized answer flow; `a` remains sole answer key.

## Current proposal
Adopt codex architecture + hermes MVP discipline. Tightened slices for parity without scope creep:

1. Segment/badge fix (minimal emit-start + projection) — table tests for continue/resume cases.
2. Per-agent focus viewport (offset reads, follow, bounded scrollback, g/G).
3. Modes + help + keymap (`overview|agentDetail|compose|answerQuestion|help`; `?` overlay; `a` answer preserved).
4+. Composer persistence + delivery + opt-in thoughts deferred.

No protocol amendment. TUI stays facade. --no-tui and one-shot contract untouched. This slice set actually ships something that feels like Codex/Claude-Code instead of promising chat-like steering.