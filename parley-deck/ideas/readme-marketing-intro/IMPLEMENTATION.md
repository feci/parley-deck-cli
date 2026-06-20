---
idea: readme-marketing-intro
status: implemented
implementer: claude-1
started: 2026-06-20
completed: 2026-06-20
branch: parley-deck-cli#readme-marketing-intro
head-commit: pending
design-pr: pending
implementation-pr: pending
---

## Summary of work

Added the marketing-style intro from `FINAL.md` to both READMEs:
- `parley-deck-cli/README.md` — replaced the dry one-line opening paragraph with
  Deliverable A (hook → "What you get" → "Inspired by — adopted & adapted").
- `parley-deck-skill/README.md` — replaced the existing intro paragraphs with
  Deliverable B (installer-framed variant + lineage block).

## Implementation plan / checklist

- [x] Files to change: `parley-deck-cli/README.md`, `parley-deck-skill/README.md`
- [x] Insert intro before existing `## Install` sections (additive, no usage removed)
- [ ] Checks: render-sanity (headings/quote/code fences balanced)
- [ ] Push: cli via branch+PR (github-pr transport); skill README to its repo

## Deviations from FINAL.md

None planned; text is used verbatim from FINAL.md deliverables A and B.

## Notes for reviewers

Doc-only change. Verify no over-claim slipped in and that every bullet still maps
to a real section/command per codex-1's feature map.

## Signoff status at implementation time

claude-1 ✅, codex-1 ✅, hermes-1 ✅ (consensus.md). antigravity-1 signoff was
still appending when implementation began; its attribution lens is already fully
reflected in FINAL.md (the "Inspired by — adopted & adapted" lineage block is its
round-01 proposal). Final push gated on confirming its signoff or recording a
per-idea waiver.
