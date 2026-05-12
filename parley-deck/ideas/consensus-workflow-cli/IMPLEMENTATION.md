---
idea: consensus-workflow-cli
implemented-by: codex
status: ready-for-review
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
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

## Ready for review

Review should focus on:

- signoff parser/validator strictness;
- CLI command shape and flag behavior;
- `parley status --idea` behavior when an idea has no runs;
- review consensus assumptions;
- finalization/reopen file mutations.
