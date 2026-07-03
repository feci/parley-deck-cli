---
agent: claude-1
idea: protocol-restructure-appendices
round: 1
date: 2026-07-03
---

## Summary

Do the least-risky thing that delivers the physical progressive-disclosure layout:
**relocate the reference sections to the back, keep their numbers, change no rule text.** A
deterministic reorder script (split on `## ` headers, reassemble, apply identically to both
copies) makes this a provably content-preserving move — verified by the drift guard plus a
sorted-line diff. Renumbering and compression are explicitly out of scope (they alter refs/text).

## 1. Numbering — keep-numbers-relocate (do NOT renumber)

Keep every section's number and physically move it. `## 11. Transport mechanics` stays titled
"11" but lives at the back, so all `§11.B`-style refs resolve unchanged. Renumbering would force
rewriting every cross-reference (there are many) — pure risk for cosmetic sequential order. The
non-sequential reading order is mitigated by an `## Appendices` banner and the already-shipped
Quickstart "core vs reference" map. **Decision: keep numbers.**

## 2. Target order

Core (reading path), unchanged relative order:
Quickstart → §0 → §1 → §2 → §3 → §4 (+§4.0/§4.0.1) → §5 → §6 → §7 → §8 → §10 (TL;DR).
Then a new `## Appendices` divider (a short banner: "reference material — read when a task needs
it; section numbers are unchanged so all §N cross-references resolve"), then:
Appendix A (already an appendix) → §9 → §11 → §12 → §13 → §14.

Rationale: §10 TL;DR stays in the core (it is orientation, not reference). §9 (session checklist)
moves to the back per the devx-speed reading-guide, but is the first appendix after Appendix A so
it stays easy to find. §4 stays in core (it is the workflow; splitting it is a separate idea).

## 3. Mechanical method + verification

A deterministic Python reorder script (run once per copy, or once and header-adjusted):
1. Read the file; split into blocks at top-level `^## ` headers (keep the pre-§0 header + Quickstart as block 0).
2. Index blocks by a stable key (the header line).
3. Emit blocks in the target order (§3 above), inserting the `## Appendices` banner before the first moved block.
4. Write back. Apply to BOTH `internal/protocol/defaults/COOPERATION.md` and
   `parley-deck/COOPERATION.md` with the SAME reordering (the two differ only in the allowlisted
   header/roster zones, which are inside their own blocks and move intact).

**Verification (all must pass):**
- `go test ./internal/protocol/...` (drift guard `TestEmbeddedDefaultMatchesLiveDeck`) green — proves both copies byte-identical outside the allowlist.
- Content-preservation: `diff <(sort pre.md) <(sort post.md)` is EMPTY for each copy (only the added `## Appendices` banner lines differ — account for them explicitly).
- The drift-guard's asserted anchors still appear exactly once each: `## 2. Active agents (roster)`, the two roster table header rows, the Workspace/Created prefixes, one `**Protocol synced:**` (live) / zero (embedded).
- Skill fallback re-synced (`diff <(tail -n+7 embedded) <(tail -n+7 skill)` empty).

## 4. Cross-reference audit

`grep -oE '§[0-9]+(\.[0-9A-Z]+)*|Appendix [A-Z]' COOPERATION.md | sort -u`, then confirm each
referenced section/appendix header still exists in the reordered file. Since we KEEP numbers,
every `§N` target survives by construction; the audit is a belt-and-suspenders check. Also grep
for "above"/"below"/"earlier"/"later" positional language that the move could falsify (e.g. "the
table above") and fix any that now point the wrong way — these are the only content-ish edits
allowed, and each must be logged.

## 5. Risks

- **Drift-guard anchors** — the guard asserts the §2 roster section line + both table header rows
  appear exactly once, and the header prefixes. The move keeps §2 in core intact, so anchors are
  preserved; still, run the guard as the gate. (R: mitigated.)
- **`protocolSha256`** — reordering changes the file hash; refresh `meta/version.json`
  `protocolSha256` (source deck; advisory per §9.0) at release, as devx-speed did.
- **Positional prose** ("see the table above/below") could dangle after a move — the §4 → §14
  distance changes. Audit + fix positional language (logged as the only text touch).
- **Consumers** — a consumer deck syncs via §9.0 freshness; this is an additive/compatible
  reorg (no rule change), so a minor `deckVersion` bump + normal sync covers it.
- **Split-brain edit** — never hand-edit one copy; the script writes both from the same reorder.

## Position

Keep-numbers relocate via a deterministic script; core = Quickstart..§10, appendices =
A,§9,§11,§12,§13,§14 under a banner; no rule text changed; verify by drift guard + sorted-line
diff + ref audit. Compression and §4 phase-split are explicitly separate future ideas.
