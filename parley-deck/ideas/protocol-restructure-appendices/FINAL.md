---
idea: protocol-restructure-appendices
status: final
author: claude-1
consensus-date: 2026-07-03
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Final plan / specification

Physically reorganize `COOPERATION.md` for progressive disclosure — core-first, reference-last —
as a **pure content-preserving move**: relocate **only §9** (session-start checklist) to sit
after §10 (TL;DR), keeping every section number and every line. No rule text is added, removed,
or altered; no `## Appendices` banner (its lines would break the zero-addition guarantee).

Final top-level order:
Quickstart → §0 → §1 → §2 → §3 → §4 → §5 → §6 → §7 → §8 → §10 → §9 → §11 → Appendix A → §12 → §13 → §14.

## Purpose / user-visible outcome

The document reads core-concepts → TL;DR → reference material, matching the shipped Quickstart
"core vs reference" map. §9/§11/§12/§13/§14 no longer interleave the core reading path (only §9
actually needed to move; §11–§14 + Appendix A were already after §10).

## Context & orientation

Applied by a deterministic reorder script (splits on `## ` headers, re-emits in the target order)
to BOTH `internal/protocol/defaults/COOPERATION.md` and `parley-deck/COOPERATION.md` identically,
then re-synced to `parley-deck-skill/references/COOPERATION.md`. Keep-numbers-relocate preserves
all `§N` cross-references by construction.

## Observable acceptance criteria

1. Section order matches the list above; §9 follows §10; §11–§14 + Appendix A are after §10.
2. **Pure move:** `diff <(sort pre) <(sort post)` is EMPTY for each copy (zero content change).
3. Both copies byte-identical outside the allowlist (`TestEmbeddedDefaultMatchesLiveDeck` green);
   skill fallback body-identical.
4. Every `§N` / `Appendix` cross-reference still resolves (audit → zero dangling).
5. Full suite green; `protocolSha256` refreshed at release.

## Idempotence & recovery

Deterministic reorder; re-running the script on the reordered file is a no-op (order already
matches). No persisted state beyond the doc + version metadata.

## Known risks / de-risking

- Non-sequential numbering (§10 before §9) — accepted; Quickstart map + `---` boundary explain it.
- Positional prose — pre-audited: only temporal / intra-§11 hits, none falsified by moving §9.
- `core ≤200 lines` NOT delivered (out of scope: §4 is ~505 lines; needs a separate phase-split idea).

## References
- Consensus: ./consensus.md
- Rounds: ./round-01/{claude-1,codex-1,hermes-1,antigravity-1}.md
