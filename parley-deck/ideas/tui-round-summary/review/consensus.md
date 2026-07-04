---
idea: tui-round-summary
review-cycle: 1
drafted-by: claude-1
date: 2026-07-04
reviewed-commit: 2b25f08
---

## Agreed fixes

All applied in fix-up cycle 1, confirmed resolved in review round-02 (codex-1 + hermes-1: "zero remaining"):
- [MAJOR, codex-1] renderHome reserves rows for Recent runs below the digest.
- [MINOR, codex-1] rune-safe position truncation.
- [MINOR-1, hermes-1] degraded-extraction marker as a truncation-proof prefix.
- [MINOR-2, hermes-1] next-action honesty when the consensus gate is unwired.

## Deferred follow-ups

- Last-N scrollable digest history + round-2+ responding-to engagement gaps (additive enrichments).

## Dismissed findings

- [NIT-1, hermes-1] aggregate mentions without per-agent attribution — ratified design.

## Signoffs

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Implementer. All agreed fixes applied; build/vet/test green incl. the Recent-runs regression test.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: zero remaining; four agreed fixes present.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: MINOR-1/MINOR-2 resolved.
