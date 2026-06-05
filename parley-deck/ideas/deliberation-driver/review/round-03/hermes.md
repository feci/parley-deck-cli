---
agent: hermes
idea: deliberation-driver
review-round: 3
date: 2026-06-05
reviewed-commit: a83efa8
---

## Summary

Slice 2 consensus gate implemented per D6/D7/D9 with live acceptance passing. All build/vet/test green. One deviation from D9 (adapter injection) and two protocol/ordering concerns identified; others are solid or minor.

## Findings

**CRITICAL**  
Content-drafter approach in driverConsensusOps.DraftFinal (and Draft) writes scaffold via consensus.Finalize then overwrites via agent. If drafter agent fails after status=final set, FINAL.md remains scaffold with final status. finalScaffoldReason escalation only triggers on next tick; status left inconsistent until human intervention. Violates D7 non-scaffold guarantee.

Fix: Move status=final write to after successful non-scaffold verification inside the adapter (or make DraftFinal atomic with post-verify status update).

**MAJOR**  
D9 deviation: driver injects ConsensusOps adapter from app layer instead of extracting internal/signoffs seam. Import direction preserved (driver does not import app) but violates "extract internal/signoffs" spec. Acceptable for reuse but not per spec.

Fix: Either update FINAL.md to ratify the adapter pattern or extract signoffs package as originally planned.

**MINOR**  
Reopen path correctly bounds by MaxRounds (next > 1+MaxRounds escalates) and invalidates .bak files before RunRound. No loop/skip possible. Good.

**NIT**  
Partial signoffs: synchronous RequestSignoffs + immediate re-Status check is sound; escalates on persistent Partial. No false-escalate or loop risk. Clear comment explains rationale.

## Open questions

None.