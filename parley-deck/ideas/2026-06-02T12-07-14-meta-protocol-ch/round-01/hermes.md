---
agent: hermes
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 1
date: 2026-06-02
---

## Summary
Proposed architecture is sound on reuse of existing round engine and additive design. Risks: driver complexity for durable resume/reconcile; side-effect gate semantics may leak into consensus if not strictly separated. Missing: explicit block-to-block seeding contract and effects-ledger idempotency schema.

## Proposed approach
Adopt the spine + gate + effects-ledger model. §12 text: blocks are Phase 0-8 invocations producing typed FINAL.md artifacts (BUSINESS_SPEC.md etc.); action blocks add `execute` sub-phase writing plan + ledger entry only. Driver owns pipeline-run.json cursor, gate.json promotion, and MCP-executed side effects behind non-bypassable prod gates. Auto-advance inside block; pause at every boundary until per-pipeline flag enables left-half auto.

## Concerns / open questions
- How does block N seed block N+1 prompt without facilitator writing another agent's file?
- Effects ledger reconcile needs concrete idempotency key format for deploy/ops actions.
- Monitoring watcher opening new root idea: does it bypass the single sticky transport?

Trade-off: strict "markdown only" for agents keeps ownership clean but pushes all execution surface area to driver (higher harness trust required).

## Risks
1. Restart double-execute despite ledger (mitigate: external idempotency + dry-run mandatory before gate).
2. Gate.json policy drift across blocks (mitigate: central typed policy evaluator).
3. Smallest increment: pipeline.yaml + BUSINESS_SPEC block only; defer action/execute and monitoring loop to v1.1.