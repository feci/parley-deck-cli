---
agent: hermes
idea: hitl-tui-questions
review-round: 1
date: 2026-05-11
reviewed-commit: a6a3609aa3aa744589d8197527e740cc1b89f9f5
---

## Summary
Reviewed the HITL TUI questions implementation in parley-deck-cli at the specified commit against FINAL.md. Code changes, tests, TUI behavior, event persistence, and CLI command were inspected locally via source, git diff, and test runs. Implementation is largely faithful to the spec with good adherence to atomic writes, event mirroring, and TUI integration.

## Findings
### [MINOR] Question ID generation includes timestamp and agent but suffix randomness is minimal (uses short hex); consider ensuring higher entropy for collision avoidance in high-frequency scenarios.
### [NIT] TUI question panel uses stable sort but lacks explicit test for ordering under concurrent file creation.
### [MINOR] CLI `parley answer` help text could better document the --dir flag behavior for consistency with other commands.

## Open questions
- How should agents be guided to handle question polling timeouts in practice (out of scope for this slice but worth noting for future agent contracts)?
- Is there a plan to surface auto-answered questions distinctly in the TUI panel for auditability?