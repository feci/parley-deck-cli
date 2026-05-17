---
idea: roadmap-implementation-plan
implemented-by: codex
date: 2026-05-17
status: complete
design-pr: https://github.com/feci/parley-deck-cli/pull/19
design-final: https://github.com/feci/parley-deck-cli/blob/6dd1621decea3ef53e41df0edab2e5c15fe4c704/parley-deck/ideas/roadmap-implementation-plan/FINAL.md
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
implementation-branch: feature/roadmap-implementation-plan
implementation-commit: 8d57b63
fix-up-cycle: 1
completed: 2026-05-17
---

## Summary

Implemented the first roadmap slice from `FINAL.md`: `round-index-artifact-pruning`.

The implementation PR separates two reviewable behaviors:

- Deterministic `round-NN/_index.md` generation after `RunRoundOne`.
- Context-only hidden-reasoning fence sanitization used by the index builder.

## Implementation plan / checklist

- Files or areas changed:
  - `internal/runner/round_index.go`
  - `internal/runner/runner.go`
  - `internal/app/app.go`
  - focused runner tests
- Checks run:
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./internal/runner`
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./...`
- Review or risk notes:
  - `_index.md` is runner-owned and derived, not participant-owned.
  - Sanitization is context-only and does not perform secret redaction.
  - Index write failure is warning-only and does not fail an otherwise successful round.

## Behavior Delivered

- Added `SanitizeForContext` with the closed initial fence set:
  - `<think>...</think>`
  - `<thought>...</thought>`
  - `<thinking>...</thinking>`
- Added deterministic round index rendering with:
  - frontmatter marking it as a derived `round-index` artifact;
  - a participant summary table;
  - per-agent status sections;
  - H2-only section extraction;
  - deterministic approximate token estimate using `(sanitized_bytes + 3) / 4`;
  - skipped and failed participant rows.
- Wired index generation into `RunRoundOne` after participant execution.
- Added `round.index_written` and `round.index_failed` events.
- Added a non-failing `runner.Result.Warning` path and CLI output for warning-only runner conditions.

## Verification

- `go test ./internal/runner`: passed.
- `go test ./...`: passed.

## Fix-up Cycle 1

Review consensus: `review/consensus.md`.

Applied agreed fixes:

- Malformed opening fences are preserved as literal text instead of dropping the artifact tail.
- Fence matching now uses ASCII-only case folding without indexing through Unicode-lowercased copies.
- H2 extraction uses direct `## ` detection without the dead `### ` guard.
- `_index.md` writing now creates the round directory before writing.
- `_index.md` frontmatter includes `token-heuristic: bytes_div_4`.
- The warning-only synthetic result now uses `runner/index` as a runner-owned sentinel.
- `approxTokens` uses `len(value)` for byte length without allocating.
- `SanitizeForContext` no longer trims leading/trailing whitespace.
- Added tests for no-fence identity, case-mixed fences, malformed opening fences, and updated index metadata/sentinel behavior.

Verification after fix-up:

- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./internal/runner`: passed.
- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./...`: passed.

## Deviations From FINAL.md

None. Repo-map, prompt context-pack wiring, secret redaction, hooks, memory, and direct SDK backend work remain deferred to later roadmap slices.

## Completion

Review cycle 2 reached zero agreed fixes in `review/consensus.md`. Implementation PR #20 is ready to merge.
