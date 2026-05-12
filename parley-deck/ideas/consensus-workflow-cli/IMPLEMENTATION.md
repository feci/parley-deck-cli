---
idea: consensus-workflow-cli
implemented-by: codex
status: ready-for-re-review
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
fix-up-cycle: 1
---

## Summary

Implemented the deterministic consensus/signoff CLI slice from `FINAL.md`.

## Changes

- Added `internal/consensus` for consensus and review-consensus file handling:
  - draft scaffolding;
  - signoff parsing/appending;
  - validation and triage;
  - design finalization scaffold;
  - blocked consensus reopen.
- Added `parley consensus` CLI commands:
  - `status`;
  - `draft`;
  - `signoff`;
  - `finalize`;
  - `reopen`.
- Extended `parley status --idea` and the workspace idea list with consensus triage when consensus files exist.
- Added package and CLI regression tests for the new workflow.

## Deviations from FINAL.md

- `consensus.*` events were not emitted because the current CLI command context has no project-level event store outside individual run directories. The filesystem remains canonical, as specified.
- Review consensus uses the same participant list as the idea prompt in this slice. A narrower review-round participant model can be added if the review workflow starts recording a separate reviewer roster.

## Verification

- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go test ./...`
- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go run ./cmd/parley consensus status --dir . consensus-workflow-cli`
- `GOCACHE=/private/tmp/parley-go-build-cache GOMODCACHE=/private/tmp/parley-go-mod-cache go run ./cmd/parley status --dir . --idea consensus-workflow-cli`

## Fix-up cycle 1

Applied agreed fixes from `review/consensus.md`:

- normalized signoff statuses before triage so manual aliases such as `Status: block` cannot bypass blocked consensus;
- accepted trailing parenthetical status text and the shorter `Counter-proposal:` field;
- changed `reopen` to restore the latest `round-NN` status and preserve aborted consensus files as `<round>-consensus-aborted-<attempt>.md`;
- added `### Non-goals` to generated `FINAL.md` scaffolds;
- added review consensus `cycle` and `reviewed-commit` frontmatter;
- surfaced malformed consensus in the idea list as `consensus=error`;
- sorted round directories numerically;
- added tests for manual blocked aliases, reserved finalization with open items, final scaffold headings, review metadata, and numeric round selection.

Remaining non-blocking follow-ups recorded by review consensus:

- project-level `consensus.*` events;
- cross-process signoff locking;
- multi-line signoff notes;
- native GitHub review submission;
- automated `request-signoffs`;
- workspace status performance for large idea counts;
- explicit finalized-status guard policy beyond the existing `FINAL.md` existence guard.

## Ready for review

Review should focus on:

- signoff parser/validator strictness;
- CLI command shape and flag behavior;
- `parley status --idea` behavior when an idea has no runs;
- review consensus assumptions;
- finalization/reopen file mutations.
