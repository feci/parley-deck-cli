---
idea: 2026-06-02T21-54-52-continue-the-4-r
review-cycle: 1
drafted-by: claude
date: 2026-06-02
reviewed-commit: e299d58
outstanding_agreed_fixes: 0
blocked: false
---

## Agreed fixes

Both codex MAJOR findings accepted and applied in fix-up cycle 1; hermes had only NITs (no change). `outstanding_agreed_fixes: 0` → review complete.

```yaml
agreed_fixes:
  - id: blockcomplete-honors-blocked
    severity: MAJOR
    from: codex
    finding: blockCompleteFunc treated an implementation block with outstanding_agreed_fixes:0 AND blocked:true as complete, bypassing a reviewer BLOCK.
    fix: completion predicate now requires `!blocked` (count==0 && !blocked && found). Regression test TestBlockCompleteRespectsBlockedReviewConsensus.
    status: applied
  - id: stdout-fallback-validate-before-persist
    severity: MAJOR
    from: codex
    finding: the stdout fallback wrote captured stdout to the artifact path on a mere "---" prefix, then validated afterward, leaving a malformed file that poisons retries.
    fix: write a .stdout-candidate temp, require the FIRST LINE to be exactly "---" (firstLineIsFence, no leading narration), validate the candidate via validateArtifactForPhase, rename only on success, else remove. Regression test TestStdoutFallbackRejectsInvalidFrontmatter (no artifact + no leftover temp).
    status: applied
```

## Deferred follow-ups

- hermes NITs: confirmations only, no action.
- Parallel multi-active DAG launches deliberation/watcher waves concurrently; action/implementation ready blocks stop the DAG drive with guidance (future: integrate impl fix-up loop into DAG waves).

## Dismissed findings

None.

## Signoffs

### Signoff: codex — 2026-06-02
Status: ✅ ACCEPT
Notes: Both MAJOR findings fixed with regression tests; fail-closed behavior on blocked consensus and invalid stdout candidates is now correct.

### Signoff: hermes — 2026-06-02
Status: ✅ ACCEPT
Notes: No CRITICAL/MAJOR; production safety + back-compat preserved; fixes verified.

### Signoff: claude — 2026-06-02
Status: ✅ ACCEPT
Notes: Implementer. Fix-up cycle 1 applied both MAJOR fixes; outstanding_agreed_fixes: 0; go build/vet + full suite green.
