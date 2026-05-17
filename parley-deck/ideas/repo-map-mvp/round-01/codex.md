---
agent: codex
idea: repo-map-mvp
round: 1
date: 2026-05-17
---

## Summary

The repo-map MVP should be a deterministic local command and small internal package, not a prompt-injection feature yet. It should provide enough structure to help agents orient themselves: file paths, coarse file type, Go package names, imports, and top-level Go symbols. Keep the implementation dependency-free and make JSON the future tooling contract.

## Proposed approach

Command shape:

```text
parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]
```

Defaults:

- `--dir .`
- `--format markdown`
- `--max-files 400`

Data model:

- `root`: absolute root used for the walk.
- `truncated`: boolean.
- `maxFiles`: configured max.
- `files`: sorted by slash-normalized relative path.
- file entry:
  - `path`
  - `kind`: `go`, `markdown`, `json`, `text`, or `other`
  - `package` for Go files
  - `imports` for Go files
  - `symbols` for Go declarations
- symbol entry:
  - `kind`: `type`, `func`, `method`, `const`, `var`
  - `name`
  - `receiver` for methods

Ignore rules:

- Always skip `.git`, `node_modules`, `vendor`, common build/cache dirs (`dist`, `build`, `target`, `.cache`, `.next`), and local Go caches.
- Skip `parley-deck/runs` because it is transient run telemetry.
- Include source, tests, docs, and committed `parley-deck/ideas` files unless the file limit is reached. The MVP is a repository map, not only a source-code map.

Implementation:

- Add `internal/repomap`.
- Use `filepath.WalkDir`, convert paths with `filepath.ToSlash`, sort before output.
- Parse `.go` files using `go/parser` and `go/ast` with imports only plus declarations. Do not type-check.
- If a Go file has parse errors, still include the file with a parse error note rather than failing the whole map.
- Add CLI handling in `internal/app`: `parley context repo-map ...`.
- Keep output byte-stable: no timestamps, no local durations.

## Concerns / open questions

- The default `max-files` should be conservative enough for this repo but not surprising. `400` is a reasonable first default; callers can raise it.
- JSON output should be stable but can be documented as an early developer surface. We do not need a long-term schema guarantee yet.
- It is tempting to add token estimates now, but that belongs to the later `context-pack-wiring` slice unless reviewers strongly prefer reusing the current byte heuristic.

## Risks

- Scope creep into ranking/import graph/type-checking. Keep this as a structural outline only.
- Over-broad ignore rules can hide important code. Default ignores should target obvious generated or cache directories.
- Markdown output can get large. The file limit and compact symbol rendering are required for usefulness.
