---
agent: hermes
idea: 2026-06-02T21-54-52-continue-the-4-r
review-round: 1
date: 2026-06-02

## Summary
Implementation correctly addresses all focus areas with additive changes, strict fail-closed behaviors, and passing tests. No blocking issues found. Production safety and back-compat preserved.

## Findings
### NIT stdout fallback strictness
What: `runAgent` rejects non-`---` stdout; only persists on exact prefix match.
Why: Prevents narration injection as specified.
Fix: Already strict; tests cover reject case. No change needed.

### NIT Phase8Decision/ReviewAgreedFixes fail-closed
What: `ReviewAgreedFixes` returns found=false on absent field or non-frontmatter; Phase8NoData case stops.
Why: Matches "fail-closed if contract field absent".
Fix: Correct; no action.

### NIT cursor back-compat
What: Ready/active fields are additive to existing CompletedBlocks etc.
Why: Old runs load without error.
Fix: Already additive; status output conditional.

### NIT ComputeDAGStep topological+gate correctness
What: Uses ReadyBlocks + gate load/auto-approve policy; waves respect deps and supervised gates.
Why: Topo + non-bypassable as required.
Fix: Tests (diamond, supervised) pass.

## Open questions
None.