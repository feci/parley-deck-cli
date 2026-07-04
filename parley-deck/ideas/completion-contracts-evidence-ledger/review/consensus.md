---
idea: completion-contracts-evidence-ledger
review-cycle: 1
drafted-by: claude-1
date: 2026-07-04
reviewed-commit: 1efc0fc
---

## Agreed fixes

All applied in fix-up cycle 1, confirmed resolved in review round-02 (codex-1 + hermes-1):
- [CRITICAL, codex-1] zero-fixes completion path now gated by the completion contract (list-form checks veto).
- [MAJOR, codex-1] secret scrubber hardened (bearer/sk-/ghp_/AWS/JWT) + real multi-shape test.
- [MAJOR, codex-1] list-form YAML syntax error fails closed (no legacy fallback).
- [MAJOR, hermes-1] driver commits the evidence write → tree stays clean across fix-up cycles.
- [MAJOR-2, hermes-1] vacuous scrub test replaced.
- [NIT, codex-1] gofmt checks_test.go — applied.

## Deferred follow-ups

- `parley check-contract` pre-flight command; richer expectation matchers (stdout/regex/path_exists).

## Dismissed findings

None.

## Signoffs

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Implementer. All agreed fixes applied; build/vet/gofmt/test + drift guard green.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: no remaining CRITICAL or MAJOR; the earlier NIT (gofmt) is now applied.

### Signoff: hermes-1 — 2026-07-04
Status: ✅ ACCEPT
Review round-02: both my MAJOR findings resolved; design stayed minimal.
