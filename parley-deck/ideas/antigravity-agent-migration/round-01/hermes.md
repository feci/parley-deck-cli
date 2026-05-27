---
agent: hermes
idea: antigravity-agent-migration
round: 1
date: 2026-05-27
---

## Summary
Hermes focuses on operational defaults for headless launches, ensuring reliable Antigravity (`agy`) configuration in CLI workflows, and reviewing release risks for the migration from Gemini.

## Proposed approach
- Set effective launch defaults to model xai/grok-4.3 with cli-default profile, balanced smart-fast speed, and 1800000ms timeout.
- Update parley init and default agent configs to prefer `agy` for new sessions while keeping overrides explicit.
- Ensure headless reliability by validating launch configs in skill installer and runtime targets.
- Review packaging and release metadata for npm/GitHub to surface the new defaults without breaking legacy paths.

## Concerns / open questions
- How will `agy` expose model/thinking controls if not yet in help output?
- Ensuring no disclosure of internal launch flags in user-facing docs.

## Risks
- Migration could impact existing projects using Gemini defaults if not clearly marked legacy.
- Release channels (Homebrew, WinGet) may lag, causing inconsistent user experience post-publish.