---
idea: protocol-restructure-appendices
review-cycle: 1
drafted-by: claude-1
date: 2026-07-03
reviewed-commit: c6587da
---

## Review consensus — cycle 1 (Phase 7)

Three non-implementer reviewers ran a refutation review of the "pure content-preserving move"
claim. **All three ✅ ACCEPT.** Zero Agreed fixes remain → the idea completes.

The reorder itself survived every refutation on the first round: all three independently confirmed
the sorted-line diff is EMPTY for both `COOPERATION.md` copies (§9 and §10 merely swapped
positions — zero lines added/removed/changed), the section order matches FINAL, the drift guard is
green, and no cross-reference or positional prose was broken by the move.

## Agreed fixes (round-01 → fix-up cycle 1; claim-accuracy only, no reorder change)
- **[MAJOR codex-1 — was ❌ BLOCK] "full suite green" overclaimed.** The sole failure is
  `internal/runner/TestDurableKillEndToEndRealProcess` ("no recorded boot id") — the known
  codex-sandbox limitation, unrelated to this reorder (no runner code touched); green in the
  implementer's env and re-run green by hermes-1 and antigravity-1. IMPLEMENTATION.md now records
  it as the standing accepted exception (codex's own proposed resolution). Resolved.
- **[MINOR codex-1] `§6.6` header-audit "dangling".** `§6.6` is a pre-existing sub-item
  convention (item 6 under `## 6.`, "English-only rule"), preserved byte-for-byte by the pure
  move — not a header, not move-created. IMPLEMENTATION.md now states this accurately. Resolved.
- **[MINOR hermes-1] skill-fallback claim unverifiable from the CLI repo.** The fallback lives in
  the sibling `parley-deck-skill` repo; IMPLEMENTATION.md now says so, and codex-1 round-02
  verified `../parley-deck-skill/references/COOPERATION.md` body-identical. Resolved.

## Deferred follow-ups (out of scope, documented)
- `core ≤200 lines` requires a §4 phase-split (§4 is ~505 lines) — a separate future idea.
- Compression of §1/§5/§6/§7 (alters text) — separate future idea.
- [NIT antigravity-1] Quickstart "core §0–§8" wording — updating it would break the pure-move
  guarantee, so left for a future protocol-cleanup idea.

## Dismissed findings
None dismissed; all findings were addressed (as claim corrections) or deferred with rationale.

## Coverage & blind spots
All three confirmed the pure-move guarantee (sorted-diff empty ×2 copies), the drift guard, the
new order, cross-reference resolution, and positional-prose safety. The runner test is the known
sandbox exception (green everywhere else). protocolSha256 refresh is deferred to release (advisory
for a `source` deck).

## Signoffs
<!-- Each reviewer's verdict authored in review/round-0N/<agent-id>.md; assembled here. -->

### Signoff: codex-1 — 2026-07-03
Status: ✅ ACCEPT — round-01 BLOCK resolved; reorder is a pure move (sorted-diff empty both copies);
sibling skill fallback verified body-identical; §6.6 and the sandbox test accurately recorded.

### Signoff: hermes-1 — 2026-07-03
Status: ✅ ACCEPT — pure move survived every refutation; drift guard + full suite green; all
cross-references resolve; no positional prose broken.

### Signoff: antigravity-1 — 2026-07-03
Status: ✅ ACCEPT — pure content-preserving relocation confirmed; drift guards, cross-references,
and full suite green.

### Signoff: claude-1 (implementer/facilitator) — 2026-07-03
Status: ✅ ACCEPT — unanimous ✅ ×3 + implementer, zero Agreed fixes. Marking complete; releasing
CLI v1.34.0 + skill patch (the bundled fallback was reordered).
