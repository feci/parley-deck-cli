---
idea: protocol-restructure-appendices
drafted-by: claude-1
date: 2026-07-03
---

## Consensus — design round 1 (unanimous)

Four independent round-01 analyses (claude-1, codex-1, hermes-1, antigravity-1) converged.

### D1 — keep-numbers-relocate (unanimous)
Physically relocate reference sections to the back but **keep every section's existing number**
(`## 11. Transport mechanics` stays "11"), so every `§11.B` / `§9.0` / `§12.11` cross-reference
resolves by construction. Renumbering is rejected: it rewrites content (ref strings are content
lines), has a large blast radius, and violates the "no rule text changed" constraint.

### D2 — MINIMAL move: relocate ONLY §9, to sit right after §10 (hermes-1's insight)
The current tail order is §8 → §9 → §10 → §11 → Appendix A → §12 → §13 → §14. The **only**
reference section sitting inside the core reading path is **§9** (session-start checklist, between
§8 and §10); §11, Appendix A, and §12–§14 are **already** after §10. So the entire reorganization
is a **single block move: §9 jumps past §10.** Final order:
Quickstart → §0 → §1 → §2 → §3 → §4 → §5 → §6 → §7 → §8 → §10 → **§9** → §11 → Appendix A → §12 → §13 → §14.
This is the least-churn way to achieve core-first / reference-last (codex proposed the same
functional result; hermes identified it is a one-block move).

### D3 — NO `## Appendices` banner (codex-1 + hermes-1)
A new banner heading would add non-blank content lines, breaking the "sorted-line diff is EMPTY"
guarantee and making this more than a pure move. **Omit it.** The core/reference boundary is
already (a) explained by the shipped Quickstart "core vs reference" map and (b) marked by the
`---` separator trailing §10. This keeps the change a provably-content-preserving reorder.
(claude-1 and antigravity-1 had proposed a banner; both accept dropping it — the purity/safety
win outweighs the cosmetic gain, and the Quickstart already carries the explanation.)

### D4 — Mechanical method + verification (unanimous)
A deterministic reorder script splits on top-level `## ` headers and re-emits blocks in the D2
order; applied identically to both `COOPERATION.md` copies (they differ only in the allowlisted
header/roster zones, which travel inside their blocks) and then re-synced to the skill fallback.
Gates (ALL must pass):
- `go test ./internal/protocol/...` (drift guard) green — both copies byte-identical outside the allowlist.
- **Content-preservation: `diff <(sort pre) <(sort post)` is EMPTY** (a true pure move — zero additions/deletions).
- Cross-reference audit: `grep -oE '§[0-9]+(\.[0-9A-Z]+)*|Appendix [A-Z]'` → every target header still exists.
- Skill fallback body-identical (`diff <(tail -n+7 embedded) <(tail -n+7 skill)` empty).
- Drift-guard asserted anchors (the §2 roster section line + both table header rows + Workspace/Created prefixes) still appear exactly once — preserved because §2 does not move.

### D5 — Scope boundaries (unanimous)
`core ≤200 lines` is NOT achievable here — §4 alone is ~505 lines (it is the core workflow). That
needs a separate §4 phase-split idea. Compression of §1/§5/§6/§7 is also out of scope (it alters
text). `protocolSha256` in `meta/version.json` is refreshed at release (source deck; advisory per §9.0).

## Agreed trade-offs
- Non-sequential reading order (§10 before §9) — accepted; the Quickstart map + `---` boundary
  explain it, and keeping numbers is worth far more than sequential elegance.

## Comparison & blind spots
- Unanimous: keep-numbers-relocate, script+sorted-diff verification, ref audit, scope boundaries.
- hermes-1 uniquely reduced it to a single §9 block move (minimal churn) — adopted.
- codex-1 + hermes-1 argued against the banner (purity) — adopted over the claude-1/antigravity-1 banner.
- Blind spot: positional prose ("see below" etc.) — pre-audited by the facilitator: the only two
  hits are temporal (line ~375) or intra-§11 ("see below" within the transport table), neither
  falsified by moving §9. Re-confirm in review.

## Signoffs
<!-- Each participant appends its own signoff (collected via signoffs/<agent>.md). -->
