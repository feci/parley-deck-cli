---
agent: hermes
idea: consensus-request-signoffs
review-round: 2
date: 2026-05-13
reviewed-commit: 7cb518507fc8f8389eee0d014162334e5ed20e71
---

## Findings

No findings.

Fix-up cycle 1 fully resolved the consensus-agreed items from round 1:
- Exactly-one-new-signoff enforcement and foreign-signoff rejection now active via `validateAppendOnlyContent` + `newSignoffAgents`.
- Existing-content edit detection added (prefix + suffix header validation).
- Partial-progress summary printed on mid-loop failure.
- BLOCK counter-proposal line present in prompt template.
- Runner preflight errors consistently wrapped as `errRequestUsage` (exit 2).
- Agent IDs rendered unquoted in prompt prose.
- Regression tests added for forged extra signoff and non-append edits.

All prior tests continue to pass; no regressions introduced.

## Open questions

## Summary

Implementation after fix-up cycle 1 is complete and correct. All review consensus fixes landed cleanly. Residual risks remain only in explicitly deferred areas: durable child-process logs, cross-process locking, and Git auto-commits. Ready for merge.