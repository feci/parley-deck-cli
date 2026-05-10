---
agent: hermes
idea: live-run-tui
round: 1
date: 2026-05-10
---

## Summary
The live-run-tui slice aims to make parley-deck-cli runs observable in real time via a terminal UI. It surfaces the current idea/phase/round, agent states (pending/running/finished/failed/skipped/unknown), elapsed times, recent events from events.jsonl, and log previews. Existing non-TUI runs remain unchanged, and the default `parley run` launches the TUI while the runner continues in the background. Token accounting and HITL features are out of scope.

## Proposed approach
Implement a lightweight TUI using a Go library such as bubbletea for model/view/update loop. On `parley run`, spawn the runner goroutine and immediately attach the TUI; poll events.jsonl and per-agent stdout/stderr tails at short intervals. Maintain separate display panes for overview, per-agent status, and scrolling event log. Graceful exit on Ctrl-C or 'q' must leave the run directory and any durable state intact. Keep the implementation under a few hundred lines with clear separation between runner and TUI concerns.

## Concerns / open questions
- How to handle terminal resizing and wide output without breaking layout?
- What is the minimal polling interval that feels live without excessive I/O?
- Should the TUI show raw event lines or render a summarized view?
- How to surface "round complete" vs "still running" without false positives on slow agents?

## Risks
- TUI bugs could mask runner failures or make it hard to diagnose issues compared to plain --no-tui logs.
- File polling on events.jsonl might miss rapid writes or cause brief display lag under heavy agent output.
- Adding the TUI as default may surprise users who expect immediate CLI completion in scripts or CI.