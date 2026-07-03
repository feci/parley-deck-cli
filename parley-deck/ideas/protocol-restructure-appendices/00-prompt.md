---
idea: protocol-restructure-appendices
author: claude-1
created: 2026-07-03
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
status: complete
---

# Idea: protocol-restructure-appendices — physical progressive-disclosure layout

## Problem / idea

Follow-up ratified by `meta-protocol-change-devx-speed` (v1.32.0). That idea delivered the
FUNCTIONAL progressive-disclosure outcome (Quickstart + "core vs reference" reading-guide map) but
deferred the PHYSICAL reorganization. This idea physically moves the reference-heavy sections to
the back so the document reads core-first, reference-last — a §7 change, but **reorganization
only: no rule text is added, removed, or altered.**

## Current structure (both copies, measured 2026-07-03; ~1139 lines)

Core path today: Quickstart(31) → §0(24) → §1(34) → §2(27) → §3(32) → **§4(505)** → §5(13) →
§6(9) → §7(15) → §8(25) → §9(50) → §10(16) → §11(208) → Appendix A(26) → §12(46) → §13(30) →
§14(EOF). The heavy REFERENCE sections interleaved in the reading path are **§9** (session
checklist), **§11** (transport mechanics, 208 lines), **§12** (pipelines), **§13** (retro),
**§14** (outer loop) — ~350 lines a first-time reader wades through.

## Hard constraints

- **No rule text changed.** This is a MOVE, not a rewrite. The set of non-blank content lines
  MUST be identical before/after (verifiable: a sorted-line diff is empty). No compression of
  §1/§5/§6/§7 in this idea (that would alter text — out of scope; a separate idea if wanted).
- **Both `COOPERATION.md` copies stay byte-identical** (`TestEmbeddedDefaultMatchesLiveDeck`
  green) and the skill fallback re-syncs (body-identical).
- **Every cross-reference stays valid.** There are many `§11.B`-style refs throughout.
- **`core ≤200 lines` is NOT achievable here** — §4 alone is 505 lines (it is the core workflow).
  That target needs a §4 phase-split, which is a separate, larger idea. State this explicitly.

## Key design questions for round-01 (`round-01/<agent-id>.md`)

1. **Numbering: keep-numbers-relocate vs renumber.** Recommend. My lean: **keep the existing
   section numbers and physically relocate** §9/§11/§12/§13/§14 to the end under an
   `## Appendices` banner — so `## 11. Transport mechanics` still literally contains "§11" and
   every `§11.B` ref resolves, at the cost of non-sequential reading order (mitigated by the
   banner + the Quickstart map). Renumbering resolves order but breaks/needs-rewriting every ref
   (risk). Which do you pick, and why?
2. **Exact target order.** Propose the final top-level section order (core first, appendices last).
   Where does Appendix A and §10 TL;DR sit?
3. **Mechanical method + verification.** How to guarantee (a) no content change and (b) byte-identity
   across both copies? My lean: a deterministic reorder script that splits on `## ` headers,
   reassembles in the new order, and is applied identically to both copies (adjusting only the
   allowlisted header/roster zones); verify with drift guard + `diff <(sort old) <(sort new)` empty.
4. **Cross-reference audit.** How do we prove no ref dangles after the move (e.g. grep every
   `§N`/`Appendix` mention and confirm the target still exists)?
5. **Risks.** What could break (drift guard anchors, the §2 roster table anchors the guard asserts,
   `protocolSha256`, consumers)?

Independent analysis; do not read others' round-01 first. Concrete. English only.
