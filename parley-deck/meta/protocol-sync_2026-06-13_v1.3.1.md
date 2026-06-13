# Protocol sync — parley-deck-skill 1.3.1

**Date:** 2026-06-13
**By:** claude
**Scope:** project deck (`parley-deck/COOPERATION.md`) + project metadata only.

## What changed

The shared skill bumped to **1.3.1**, whose only substantive protocol change is
the new additive section **§12 "Pipeline blocks & action stages"** (opt-in;
nothing about existing ideas/reviews/phases changes).

- **§12 was already present and byte-identical to packaged 1.3.1.** We had
  ratified §12 locally earlier via the `pipeline-protocol-change` idea, and that
  text matches upstream 1.3.1 exactly (`diff` of the two §12 sections is empty,
  45 lines each). So no append was needed — lockstep held.
- **Header:** added a `**Protocol synced:** 2026-06-13 — parley-deck-skill 1.3.1
  (claude)` line to record provenance (neither our header nor the packaged
  reference previously carried one).
- **Metadata:** `parley-deck-skill sync-project --project . --yes` refreshed
  `parley-deck/meta/version.json` only (deckVersion 1.2.0 → 1.3.1; protocol /
  skill / packaged shas re-stamped). It does not touch `COOPERATION.md`.

## Verification

- `diff` of our §12 vs packaged 1.3.1 §12 → empty (byte-identical).
- Full `diff parley-deck/COOPERATION.md <packaged>` hunks are **only** our
  intentional customizations: the §2 roster, the Phase 0 `strict_gate` line,
  Phase 6 "Review briefs and dispositions", Phase 8 strict-gate + "Stopping
  judgment", and §8 "Consults" (all from the 1.24.0 kindly-adoption work). No
  unexpected drift; nothing after §8 differs.
- `parley-deck-skill status` after sync: `project-metadata-stale` and
  version-drift reasons are gone; `metadataStatus: valid`,
  `metadataMatchesProtocol: true`, `deckVersion: 1.3.1`. The sole remaining
  compatibility reason is `project-protocol-differs-from-packaged-reference`,
  which is expected for any customized deck.

## Not in scope (flagged follow-up)

The **embedded default** protocol `internal/protocol/defaults/COOPERATION.md`
still lacks §12 entirely (and also predates the Phase 6/8/§8 amendments). This
is the pre-existing drift already recorded in the kindly inbox note
`claude-to-all_review-gate-honesty_external-skill-snapshot-sync.md` and is left
for its own dedicated protocol-change idea — it is a broader resync than this
spec-only §12 carry, and §12 is SPEC-only (no local-dir driver pipeline code
ships in 1.3.1), so day-to-day facilitation is unchanged.
