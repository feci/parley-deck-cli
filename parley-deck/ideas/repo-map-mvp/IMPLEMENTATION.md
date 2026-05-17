---
idea: repo-map-mvp
implemented-by: codex
date: 2026-05-17
status: ready-for-re-review
design-pr: https://github.com/feci/parley-deck-cli/pull/21
design-final: https://github.com/feci/parley-deck-cli/blob/0ec2dc42eee346f1d4ee5f2730ee707360ebfe40/parley-deck/ideas/repo-map-mvp/FINAL.md
implementation-pr: https://github.com/feci/parley-deck-cli/pull/22
implementation-branch: feature/repo-map-mvp
implementation-commit: 07b26ad
fix-up-cycle: 1
---

## Summary

Implemented `repo-map-mvp`: a deterministic local repository map command and package.

## Fix-up Cycle 1

Applied review consensus fixes from `review/consensus.md`:

1. Sanitized Go file read errors so `parse_error` reports the relative repository path and does not leak an absolute developer-machine root.
2. Added walker coverage for skipped symlinks and directories as practical non-regular paths.
3. Added CLI usage coverage for bare `parley context` and an unknown context subcommand.
4. Removed the undocumented `--format md` alias; accepted formats are exactly `markdown` and `json`.

Checks run after fix-up:

- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./internal/repomap ./internal/app`
- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./...`
- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go run ./cmd/parley context repo-map --dir . --format md` exits non-zero with `invalid format "md"; expected markdown or json`.

## Implementation plan / checklist

- Files or areas changed:
  - `internal/repomap/repomap.go`
  - `internal/repomap/repomap_test.go`
  - `internal/app/context.go`
  - `internal/app/app.go`
  - `internal/app/app_test.go`
- Checks run:
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./internal/repomap ./internal/app`
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./...`
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go run ./cmd/parley context repo-map --dir . --format markdown --max-files 20`
  - `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go run ./cmd/parley context repo-map --dir . --format json --max-files 3`
- Review or risk notes:
  - Output is stdout-only in this slice.
  - Built-in ignores and deterministic truncation are implemented.
  - Go parsing continues on parse errors and records `parse_error`.
  - Prompt wiring, gitignore compatibility, `--out`, and verbose diagnostics remain deferred.

## Behavior Delivered

- Added:

```text
parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]
```

- Defaults:
  - `--dir .`
  - `--format markdown`
  - `--max-files 1000`
- JSON schema:
  - `schema_version`
  - `root`
  - `max_files`
  - `truncated`
  - `counts`
  - `files`
- File fields:
  - `path`
  - `kind`
  - `size_bytes`
  - `package`
  - `imports`
  - `symbols`
  - `parse_error`
- Symbol fields:
  - `kind`
  - `name`
  - `receiver`
  - `exported`
  - `line`

## Verification

- `go test ./internal/repomap ./internal/app`: passed.
- `go test ./...`: passed.
- Markdown smoke command: passed.
- JSON smoke command: passed.

## Deviations From FINAL.md

None.
