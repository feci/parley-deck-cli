---
agent: codex-1
idea: protocol-restructure-appendices
review-round: 2
date: 2026-07-03
reviewed-commit: c6587da
---

## Summary

The reorder itself remains a pure content-preserving move. I re-ran the sorted-line diff checks against `main` for both in-repo protocol copies, and both were empty. The section order still matches FINAL.md: core sections first, then `## 10. TL;DR`, `## 9. Session-start checklist`, `## 11. Transport mechanics`, `## Appendix A`, and `## 12`-`## 14`.

The drift guard is green with a fresh `go test -count=1 ./internal/protocol/...` run. The sibling skill fallback is now directly verifiable from `../parley-deck-skill/references/COOPERATION.md`; its body comparison against the embedded default copy is empty.

## Verification of round-01 findings

Round-01 MAJOR, full suite not green: accurately recorded as an accepted environment exception. I re-ran `go test ./...` in this codex sandbox. The only failure is still `internal/runner/TestDurableKillEndToEndRealProcess`, with `process verification failed (no recorded boot id); not killed`. The implementation does not touch runner code, and `IMPLEMENTATION.md` now records this as the known codex-sandbox limitation while noting green runs in the implementer, hermes-1, and antigravity-1 environments. That resolves my claim-accuracy block.

Round-01 MINOR, `§6.6` header audit: accurately recorded. `IMPLEMENTATION.md` now states that `§6.6` is a pre-existing sub-item convention for item 6 under `## 6. Conflict-avoidance mechanics`, not a markdown header and not a move-created dangling reference. The pure sorted diff confirms this text was preserved byte-for-byte by the reorder.

## New findings (if any)

None. I found no new defect in the reorder.

## Signoff

Status: ✅ ACCEPT
