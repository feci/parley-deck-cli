---
agent: codex
idea: tui-agent-controls
review-round: 1
date: 2026-05-14
reviewed-commit: f4315d4
---

## Summary

Implementer self-review added because the current `parley consensus draft --review` command requires a review artifact for every participant before drafting review consensus. Non-implementer reviews from Claude, Gemini, and Hermes are the primary defect-finding inputs for this cycle.

## Findings

### [MINOR] Apply non-implementer rendering and test coverage fixes before merge

Claude and Gemini identified small but valid gaps: the stale `initModel.View` post-init branch bypasses `newModel`, headless command rendering omits the base command when no explicit headless fields exist, backend fallback can display `cli-default`, and tests do not cover arrow keys, `h`/`m`, or focus-gated mode changes. These do not invalidate the feature, but addressing them will make the slice more robust before merge.

Suggested fix: include those items in review consensus and apply a single fix-up cycle.

## Open questions

None.
