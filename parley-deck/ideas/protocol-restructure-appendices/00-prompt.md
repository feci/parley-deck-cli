---
idea: protocol-restructure-appendices
author: claude-1
created: 2026-07-03
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
status: proposed
---

## Problem / idea

Follow-up ratified by `meta-protocol-change-devx-speed` (2026-07-03). That idea delivered the
**functional** progressive-disclosure outcome (Quickstart + "core vs reference" reading-guide
map) but deliberately deferred the **physical** reorganization to keep the diff atomic and the
drift guard tractable. This idea performs the physical move:

- Physically relocate §9 (session checklist), §11 (transport mechanics), §12 (pipelines),
  §13 (retro), §14 (outer loop) to clearly-marked **appendices** at the end, so the core is
  ≤ ~200 lines before the first appendix (FINAL acceptance criterion 4, layout half).
- Compress §1 / §5 / §6 / §7 as planned; keep every cross-reference valid via an audited
  renumber (or "Appendix X (§N)" labels that preserve existing `§11.B`-style references).

## Constraints
- **Atomic edit + cross-reference audit**; both COOPERATION.md copies stay byte-identical
  (drift guard `TestEmbeddedDefaultMatchesLiveDeck` green) and the skill fallback re-syncs.
- No rule text changed, added, or dropped — this is reorganization only (a §7 restructure,
  not new semantics).

## Non-goals
- Any new track/classifier semantics (delivered by `meta-protocol-change-devx-speed`).
- CLI enforcement (that is `track-aware-driver`).
