---
idea: repo-map-mvp
cycle: 2
drafted-by: codex
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/22
reviewed-commit: 07b26ad
status: ready
---

## Agreed fixes

No remaining agreed fixes. Review round 02 verified fix-up cycle 1.

## Deferred follow-ups

- Optional deeper non-regular-file walker coverage with a FIFO/socket-style entry when the platform and test environment make that practical. The current symlink plus directory coverage satisfies this MVP's `where practical` scope.
- `--max-files 0` behavior remains deferred from cycle 1.
- JSON optional-field `omitempty` behavior remains accepted and deferred from cycle 1.

## Dismissed findings

- Claude round-02 NIT on non-regular-file coverage is accepted as non-blocking and deferred. No participant requested a fix for this slice.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-17
Status: ✅ ACCEPT
Notes: Review round 02 has no remaining agreed fixes. The only new item is a deferred NIT on deeper non-regular-file coverage, so the implementation is ready to mark complete after participant signoffs.

### Signoff: claude — 2026-05-17
Status: ✅ ACCEPT
Notes: My round-02 verdict was ACCEPT and the consensus draft matches it. All four consensus fixes (parse_error sanitization, symlink/non-regular walker coverage, bare/unknown `parley context` usage tests, removal of the `--format md` alias) verified against the `07b26ad` inline diff plus facilitator transcript. My round-02 NIT on FIFO/socket coverage is correctly recorded as dismissed/deferred under FINAL's "where practical" scope. The two carry-forward deferrals (`--max-files 0`, JSON `omitempty`) remain appropriately out of this slice.

### Signoff: gemini — 2026-05-17
Status: ✅ ACCEPT
Notes: The fix-up cycle 1 successfully addresses all findings from the review rounds. As the context-efficiency and output-shape lens, I confirm that the deterministic output (both Markdown and JSON) and the data model for symbols strictly follow the design in `FINAL.md`. The implementation is ready for completion.

### Signoff: hermes — 2026-05-17
Status: ✅ ACCEPT
Notes: round-02 verified fix-up cycle 1; no remaining agreed fixes; implementation ready to complete (deferred NIT noted)
