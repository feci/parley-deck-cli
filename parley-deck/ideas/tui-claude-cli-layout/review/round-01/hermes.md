---
agent: hermes
idea: tui-claude-cli-layout
review-round: 1
date: 2026-06-04
---

## Summary
Reviewed FINAL.md + IMPLEMENTATION.md + git diff main...HEAD + internal/tui/live.go + live_test.go. Builds and relevant tests pass. UX gaps exist around Claude-CLI fidelity and key routing coverage; documented deviations are acceptable but some warrant MAJOR.

## Findings
### [MAJOR] Tab strip truncation can hide active tab on narrow terminals
Tab strip overflow logic keeps active + neighbors + …+N but does not guarantee the active tab label remains fully visible when width < ~40 cols; active marker can be clipped.

### [MAJOR] No line-scroll with arrows in transcript view
Claude CLI supports ↑/↓ for line scroll in transcript; current impl only page-up/down + follow toggle. Feels incomplete vs target.

### [MINOR] Answer-vs-steer distinction relies on colour-flip + banner only
Single input box + Enter routes to answer or steer based on mode; colour + "Answer" banner reduces but does not eliminate accidental-submit risk when multiple questions open or rapid tab switch.

### [MINOR] Key-routing test insufficient
TestKeyRoutingPrintableAppendsNotHotkey covers printable but omits Enter answer-before-steer, slash-command routing, and scroll keys in agent buffers.

### [NIT] Per-agent buffer cap comment stale
maxFocusLines comment says "20k" but const is 20000; minor doc drift.

## Open questions
None.