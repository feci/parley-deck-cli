---
idea: named-roster-presets
review-cycle: 1
drafted-by: claude-1
date: 2026-07-04
reviewed-commit: 8ab0e8c
---

## Agreed fixes

All applied in fix-up cycle 1, confirmed resolved in review round-02 (codex-1 + hermes-1: "zero remaining"):
- [CRITICAL, hermes-1] provenance comment moved outside the frontmatter fence.
- [MAJOR, codex-1] fail-closed when the §2 roster is unparseable.
- [MAJOR, hermes-1] `parley preset list` warns when §2 validation is skipped.

## Deferred follow-ups

- [MINOR, codex-1] full `parley run --preset` integration test (components are unit-covered).

## Dismissed findings

None.

## Signoffs

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Implementer. All agreed fixes applied; build/vet/test green.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: zero remaining.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: zero remaining on both in-scope items.
