---
agent: hermes
idea: tui-layout-refresh
round: 1
date: 2026-05-26
---

## Summary
Operational lens on long-running TUI sessions: focus on persistent dashboard visibility for active agents, live status streams, and session history without disrupting core layout.

## Proposed approach
- Add a collapsible right pane (or bottom dock) for live agent status and recent events.
- Use existing live.go streams; render compact rows with agent ID, model, elapsed time, last action.
- Trade-off: 15-20% width cost for always-on visibility vs. full-width focus mode toggle.

## Concerns / open questions
- How to handle >8 concurrent agents without scrolling or truncation?
- Does the current app.go layout use fixed vs. flex panes?

## Risks
- Overcrowding the TUI on narrow terminals; may require minimum width guard.
- Live updates could increase render churn if not throttled.