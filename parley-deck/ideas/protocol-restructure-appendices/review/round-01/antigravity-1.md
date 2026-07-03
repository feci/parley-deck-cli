---
agent: antigravity-1
idea: protocol-restructure-appendices
review-round: 1
date: 2026-07-03
reviewed-commit: bfaea9b
---

## Summary

Successfully verified the relocation of §9 (Session-start checklist) to follow §10 (TL;DR) in `COOPERATION.md`. The move is confirmed to be a pure content-preserving relocation with zero line changes. All drift guards, cross-references, and the full test suite remain green.

## Refutation attempts

- **PURE MOVE**: Verified via sorted diff comparison for both `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`. Both diffs were empty, confirming no lines were added, removed, or altered (only reordered).
- **NEW ORDER**: Verified header sequence via `grep`. The order is correctly Quickstart → §0–§8 → §10 → §9 → §11 → Appendix A → §12–§14.
- **BYTE IDENTITY**: Verified via `go test ./internal/protocol/...`. The embedded defaults match the live deck protocol.
- **SKILL FALLBACK**: Verified via body-diff between the embedded protocol and the skill reference file; they are identical.
- **CROSS-REFS**: Audited all `§N` and `Appendix` references. Every reference still resolves to a valid header; no dangling pointers were created.
- **POSITIONAL PROSE**: Audited keywords ("above", "below", "table", etc.) in and around the relocated sections. No relative references were invalidated by the reordering.
- **FULL SUITE**: `go test ./...` passed without failures.

## Findings

### [NIT] Core range description in Quickstart is stale
In [COOPERATION.md:L37](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md#L37), the text still defines the "core" range as **§0–§8**. Since **§10** has been moved before **§9** and is logically part of the core material, this summary description is technically incomplete. However, updating it would have broken the "pure move" (empty sorted diff) guarantee.
**Fix**: None required for this PR (priority: purity); should be updated in a subsequent protocol-cleanup idea.

## Signoff

Status: ✅ ACCEPT
