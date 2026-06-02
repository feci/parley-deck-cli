---
agent: hermes
idea: 2026-06-02T21-14-49-finish-the-6-rem
review-round: 1
date: 2026-06-02
---

## Summary

Phase 6 implementation reviewed for safety invariants. Core guarantees (production gate non-bypassable, auto never advances past action blocks, breach dedupe persistence, DAG validation, linear back-compat) hold. Minor surface issues noted; no CRITICAL or MAJOR blockers.

## Findings

### MINOR Action block stop message slightly inconsistent with implementation intent

The printed message in runPipelineAuto for action blocks claims "plan finalized" but the preceding code still calls autoDriveDeliberationBlock only when !done. If already done, the message is emitted without re-driving. This is cosmetic but could mislead operators reading logs.

Fix: Align the message text to "action block reached; awaiting human gate" regardless of done state.

### NIT Breach fingerprint persistence uses map keyed by fp without explicit collision handling comment

In runPipelineWatch, current[fp] = true and priorOpen iteration assumes fingerprints are unique. While Breach.Fingerprint() is deterministic, a comment noting the dedupe contract (same signal+threshold+class+observed yields identical fp) would prevent future drift.

Fix: Add one-line doc comment above the current map.

### NIT DAG validateDAG not exercised by any auto path in this change

validateDAG and ReadyBlocks exist and are correct (Kahn), but the shipped auto still uses linear Driver.Advance. Backward-compat is preserved, yet the new DAG engine is untested in the auto flow.

Fix: No action required for this round; note for future DAG-auto integration test.

## Open questions

None. All safety properties verified via diff + IMPLEMENTATION.md.