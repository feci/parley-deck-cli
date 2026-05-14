---
agent: hermes
idea: tui-agent-controls
review-round: 1
date: 2026-05-14
reviewed-commit: f4315d4
---

## Summary
Implementation matches FINAL.md scope with session-local overrides and detail rendering. No launch paths touched. Tests and init transition verified.

## Findings
### [MINOR] Agent pane markers lack explicit focus indicator
Markers use only text but selected row in agents pane does not visually distinguish focus vs idea focus; tab switch is silent.
Fix: prepend focused pane marker like ">" to the selected row in the active list only.

### [NIT] Detail view omits explicit "session-only" label on every override line
Override is noted once at top of details but individual effective mode lines do not repeat "(session)".
Fix: append the qualifier to the effective launch-mode line when an override exists.

## Open questions
None.