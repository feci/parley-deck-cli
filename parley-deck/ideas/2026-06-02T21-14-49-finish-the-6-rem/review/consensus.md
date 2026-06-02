---
idea: 2026-06-02T21-14-49-finish-the-6-rem
review-cycle: 1
drafted-by: claude
date: 2026-06-02
reviewed-commit: c872cb7
outstanding_agreed_fixes: 0
---

## Agreed fixes

All raised findings were accepted and applied in fix-up cycle 1 (machine-readable
`agreed_fixes` list below; `outstanding_agreed_fixes: 0` signals review complete).

```yaml
agreed_fixes:
  - id: decider-low-risk-only
    severity: MAJOR
    from: codex
    finding: AutoApproveWithDecider auto-approved normal-risk gates via a decider, broader than the agreed "low-risk only".
    fix: decider path now approves ONLY RiskLow (never ""/normal/high/production); auto-left keeps low/normal. Test updated.
    status: applied
  - id: dag-execution-aware-driver
    severity: MAJOR
    from: codex
    finding: execution=dag validated but Driver.Advance ran linear blocks[]-order, so a dependent listed before its prerequisite could run first.
    fix: added advanceDAG — single-active executor selecting the next block by topological readiness (ReadyBlocks over completed inbound sources), with a per-transition boundary gate. New test TestAdvanceDAGRespectsTopologicalOrder proves dependency order with reverse-listed blocks.
    status: applied
  - id: action-stop-message
    severity: MINOR
    from: hermes
    finding: auto action-block message said "plan finalized" even when already done.
    fix: reworded to "action block ready (plan finalized)".
    status: applied
```

## Deferred follow-ups

- hermes NIT: breach dedupe contract comment — applied (added doc comment), not deferred.
- hermes NIT / codex note: full PARALLEL multi-active DAG auto-drive (auto walking ReadyBlocks for concurrent blocks) — deferred; the shipped DAG executor is single-active but dependency-correct. Tracked for a future idea.

## Dismissed findings

None.

## Signoffs

### Signoff: codex — 2026-06-02
Status: ✅ ACCEPT
Notes: Both MAJOR findings (decider low-risk-only; DAG execution-aware driver) addressed with code + tests; production-execution boundary preserved.

### Signoff: hermes — 2026-06-02
Status: ✅ ACCEPT
Notes: No CRITICAL/MAJOR; MINOR message + NIT comment applied. Safety invariants verified.

### Signoff: claude — 2026-06-02
Status: ✅ ACCEPT
Notes: Implementer. Fix-up cycle 1 applied all 3 agreed fixes; outstanding_agreed_fixes: 0; go build/vet + full suite green. Single-active DAG + production gates intact.
