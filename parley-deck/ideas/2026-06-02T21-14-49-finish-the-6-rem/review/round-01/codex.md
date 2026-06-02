---
agent: codex
idea: 2026-06-02T21-14-49-finish-the-6-rem
review-round: 1
date: 2026-06-02
---

## Summary

The implementation preserves the most important production execution boundary: `pipeline auto` stops at action blocks and `pipeline execute` does not perform provider mutation itself. However, two parts of the agreed contract are not actually upheld: decider approval is broader than "low-risk only", and DAG manifests validate but still run through the linear driver.

## Findings

### [MAJOR] Decider auto-approval incorrectly includes normal-risk gates

`AutoApproveWithDecider` treats `RiskNormal` and empty risk as decider-approvable via `nonProdLow := risk == RiskLow || risk == "" || risk == RiskNormal`, and the new test locks this behavior in. This contradicts FINAL.md item 6b and the manifest comment that the decider is "low-risk non-prod only". It matters because a supervised pipeline with a decider can now advance normal-risk boundaries without a human, expanding the safety envelope beyond what was agreed.

Concrete fix: split the policies. Keep `auto-left` compatible with the prior low/normal behavior if that is intentional, but require `hasDecider` to approve only `risk == RiskLow` and never `""`/`RiskNormal`/`RiskHigh`/`RiskProduction`. Update `TestDeciderAutoApprovesLowRiskOnly` to assert normal-risk is rejected for the decider path.

### [MAJOR] `execution: dag` is accepted but the driver ignores DAG dependencies

`Manifest.Validate` accepts `execution: dag` and validates the edge graph, but `Driver.Advance` still computes `idx := blockIndex(CurrentBlock)`, advances to `Blocks[idx+1]`, and creates a gate for only that adjacent pair. The new `ReadyBlocks` helper is not wired into `start`, `continue`, `auto`, or the cursor schema with `ready_blocks`. This means a valid DAG whose `blocks[]` order differs from dependency order can seed and run a dependent block before all inbound prerequisites are complete, and joins/fan-out are serialized or skipped according to list order rather than the graph.

Concrete fix: either reject `execution: dag` until the driver is execution-aware, or implement the agreed DAG cursor behavior: persist `completed_blocks[]` plus `ready_blocks[]`, seed all currently ready roots, advance based on inbound edge completion and gate approval via `ReadyBlocks`, and add tests where block order intentionally differs from topological order to prove dependencies are enforced.

## Open questions

None.
