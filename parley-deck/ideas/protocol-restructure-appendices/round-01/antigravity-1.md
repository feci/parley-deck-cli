---
agent: antigravity-1
idea: protocol-restructure-appendices
round: 1
date: 2026-07-03
---

## Summary
The proposed physical reorganization of `COOPERATION.md` aims to implement a progressive-disclosure layout without altering rule text. By moving reference-heavy sections (§9, §11, §12, §13, §14) to appendices, we improve the core reading experience for new participants while maintaining the integrity of the protocol rules and cross-references.

## Proposed approach

### 1. Numbering: keep-numbers-relocate vs renumber
**Recommendation:** Keep existing section numbers and physically relocate.
**Why:** 
- **Integrity:** Renumbering is a text change that carries a high risk of breaking the many internal cross-references (e.g., `§11.B`, `§4 Phase 6-8`). 
- **Constraint Compliance:** The idea is explicitly a MOVE, not a rewrite. Relocating §11 to the back while keeping it titled "## 11. Transport mechanics" satisfies the constraint perfectly.
- **Mitigation:** The "non-sequential" nature of the numbering is a minor UX trade-off that is easily mitigated by the `## Appendices` banner and the existing "Core vs reference" map in the Quickstart.

### 2. Exact target order
The document should be split into a **Core** block and an **Appendices** block.

**Target Core Order:**
1.  **Quickstart**
2.  **§0. Choose the transport**
3.  **§1. Scope and purpose**
4.  **§2. Active agents (roster)**
5.  **§3. Directory layout**
6.  **§4. Protocol — phases of an idea**
7.  **§5. Quorum and async participation**
8.  **§6. Conflict-avoidance mechanics**
9.  **§7. Changing this protocol**
10. **§8. Inbox (lightweight channel)**
11. **§10. TL;DR** (Stays at the end of Core as the final summary before reference material).

**Target Appendices Order:**
12. **## Appendices** (New Banner)
13. **§9. Session-start checklist for every agent**
14. **§11. Transport mechanics**
15. **Appendix A — Adopting this protocol in a new project**
16. **§12. Pipeline blocks & action stages**
17. **§13. Retrospective optimization**
18. **§14. Automated outer loop (loop engineering) — the human brake**

### 3. Mechanical method + verification
To ensure zero content changes and byte-identity across copies:
- **Method:** Use a deterministic reorder script (e.g., Python or Zsh) that splits the file into chunks based on `## ` headers, then concatenates them in the new order.
- **Verification A (Content):** `diff <(grep -v "^#" old.md | grep -v "^$" | sort) <(grep -v "^#" new.md | grep -v "^$" | sort)` must be empty. This proves every non-header, non-blank line of rule text is preserved and unchanged.
- **Verification B (Identity):** Calculate `shasum -a 256` for both the `parley-deck/COOPERATION.md` and the skill's reference copy after the move; they must be identical.

### 4. Cross-reference audit method
We must prove no reference dangles:
- **Extraction:** `grep -oE "§[0-9]+(\.[0-9A-Z]+)?|Appendix [A-Z]" COOPERATION.md | sort -u > refs.txt`
- **Validation:** For every entry in `refs.txt`, verify that a corresponding `## N. Name` or `### N.X Name` header exists in the new file.
- **Scripted Check:** A simple script can automate this by checking each extracted reference against the set of headers in the reorganized document.

### 5. Risks
- **Drift-guard anchors:** The current drift guard likely asserts the content of specific sections. If it relies on line numbers, it will break. We must update the drift guard to use pattern-based anchors that are independent of section position.
- **protocolSha256:** The file's overall hash will change. Any consumer (like `parley-deck-skill status`) that validates the protocol version by hash will need to be updated with the new hash in `meta/version.json`.
- **Consumer breakage:** Automated tools that expect §11 or §14 at specific offsets or after specific sections may fail. We should prioritize pattern-matching headers over positional assumptions in all Parley-related tooling.
- **§2 Roster Anchors:** The roster table in §2 is a frequent target for automated writes. Since §2 remains in the Core block, this risk is low, but we must ensure the "Core" block remains stable.

## Concerns / open questions
- **Core Length:** As noted in the prompt, §4 alone is 505 lines. Even with this reorganization, the "Core" reading path will remain > 600 lines. The `core ≤ 200 lines` target is unreachable without splitting §4, which should be a priority for a follow-up idea.

## Risks
- **Reference Rot:** If participants stop reading the appendices because they are "at the back," subtle protocol changes in transport mechanics (§11) might be missed. The Quickstart "Core vs reference" map must be very clear about when to dive into appendices.
