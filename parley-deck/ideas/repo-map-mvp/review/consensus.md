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
