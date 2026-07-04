---
idea: parley-learn-playbooks
review-cycle: 1
drafted-by: claude-1
date: 2026-07-04
reviewed-commit: 8825dd3
---

## Agreed fixes

All applied in fix-up cycle 1, confirmed resolved in review round-02 (codex-1 + hermes-1: "zero remaining"):
- [MAJOR, codex-1] hardened playbook write boundary (parent-symlink refusal + O_EXCL atomic create).
- [MINOR, hermes-1] softened the §13.5 v1-skeleton wording.
- [NIT, hermes-1] clarified the skill-fallback-in-sibling-repo wording.

## Deferred follow-ups

- parley learn --refresh; cross-idea playbook synthesis; Phase-0 auto-suggestion.

## Dismissed findings

- [MINOR, hermes-1] playbooks/ absent from the §3 directory tree — dismissed as consistent with consults/ (also advisory, also not in the tree).

## Signoffs

### Signoff: claude-1 — 2026-07-04
Status: ACCEPT (✅)
Implementer. All agreed fixes applied; build/vet/gofmt/test + drift guard green.

### Signoff: codex-1 — 2026-07-04
Status: ACCEPT (✅)
Review round-02: zero remaining; the symlink write-boundary hole is closed.

### Signoff: hermes-1 — 2026-07-04
Status: ACCEPT (✅)
Review round-02: zero remaining; all three round-01 items resolved.
