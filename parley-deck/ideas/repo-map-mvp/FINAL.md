---
idea: repo-map-mvp
finalized-by: codex
date: 2026-05-17
consensus: consensus.md
status: final
---

## Decision

Implement the second roadmap slice: a deterministic local repository map command and internal package.

Command:

```text
parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]
```

Defaults:

- `--dir .`
- `--format markdown`
- `--max-files 1000`

## Required Behavior

- Output goes to stdout.
- Markdown is the default human-readable format.
- JSON is stable and intended for future tooling.
- Output must be byte-deterministic for fixed inputs:
  - no timestamps;
  - no elapsed time;
  - no absolute developer-machine root in rendered bodies;
  - stable sorting of files, imports, and symbols.
- Use only Go standard-library parsing:
  - `go/parser`
  - `go/ast`
  - `go/token`
- Extract Go symbols from source and test files by default.
- Continue on Go parse errors: include the file and expose `parse_error`.
- Do not type-check.
- Do not add tree-sitter or external parser dependencies.
- Do not wire repo maps into agent prompts in this slice.

## JSON Data Model

Top-level fields:

- `schema_version`: `1`
- `root`: `.`
- `max_files`
- `truncated`
- `counts`
- `files`

Counts:

- `files`
- `symbols`

File fields:

- `path`
- `kind`: `go`, `markdown`, `json`, `text`, or `other`
- `size_bytes`
- `package`
- `imports`
- `symbols`
- `parse_error`

Symbol fields:

- `kind`: `type`, `func`, `method`, `const`, or `var`
- `name`
- `receiver`
- `exported`
- `line`

## Ignore And Limit Rules

Built-in directory ignores:

- `.git`
- `node_modules`
- `vendor`
- `dist`
- `build`
- `target`
- `.cache`
- `.next`
- `.venv`
- `__pycache__`
- `.idea`
- `.vscode`
- `.gocache`
- `.gomodcache`

Also skip:

- `parley-deck/runs`
- non-regular files
- symlink targets

Apply `max-files` after collecting and sorting candidate files. If the candidate list exceeds the limit, include the first `N` files deterministically, set `truncated: true`, and keep exit code `0`.

## Tests Required

- Go symbol extraction:
  - functions;
  - methods and receiver names;
  - types;
  - consts;
  - vars;
  - exported flags;
  - line numbers;
  - imports;
  - parse-error degradation.
- Walker:
  - built-in ignores;
  - `parley-deck/runs`;
  - non-regular and symlink handling where practical;
  - slash-normalized relative paths;
  - deterministic truncation.
- Renderers:
  - deterministic Markdown bytes;
  - deterministic JSON bytes.
- CLI:
  - usage;
  - `--format json`;
  - `--format markdown`;
  - invalid format;
  - `--max-files`.

## Deferred Follow-ups

- `--out` atomic writes.
- external ignore files or `.gitignore` compatibility.
- verbose exclusion diagnostics.
- exported-only filters, package filters, depth filters, or token budget filters.
- prompt wiring in `context-pack-wiring`.
- incremental indexing or cache keys.
