---
idea: repo-map-mvp
drafted-by: codex
date: 2026-05-17
---

## Agreed decisions

- Implement a deterministic local repository map as the second roadmap slice.
- Add a new `context` command group with:

```text
parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]
```

- Defaults:
  - `--dir .`
  - `--format markdown`
  - `--max-files 1000`
- Output goes to stdout only in this slice. `--out`, verbose exclusion reports, external ignore files, and strict non-zero limit exits are deferred.
- Keep output byte-deterministic: no timestamps, elapsed time, absolute developer-machine root in JSON/Markdown bodies, or map-iteration ordering.
- Use Go standard-library parsing only: `go/parser`, `go/ast`, and `go/token`.
- Extract Go symbols from both source and test files by default.
- Do not type-check, rank files, compute embeddings, perform LLM summarization, or wire the map into agent prompts in this slice.

## Data model

The JSON output should include:

- `schema_version`: start at `1`.
- `root`: normalized root display value, `.` for the selected root.
- `max_files`: configured max.
- `truncated`: true when the map reached `max_files`.
- `counts`: at least `files` and `symbols`.
- `files`: sorted by slash-normalized relative path.

Each file entry should include:

- `path`
- `kind`: `go`, `markdown`, `json`, `text`, or `other`
- `size_bytes`
- `package` for Go files
- `imports` for Go files
- `symbols` for Go files
- `parse_error` when Go parsing fails but the file is still included

Each symbol entry should include:

- `kind`: `type`, `func`, `method`, `const`, or `var`
- `name`
- `receiver` for methods
- `exported`
- `line`

## Ignore and limit behavior

- Built-in directory ignores:
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
- Skip `parley-deck/runs`.
- Skip non-regular files and do not follow symlinks.
- Apply `max-files` deterministically after collecting and sorting candidate files. When truncated, render a clear `truncated` signal in JSON and Markdown. Keep exit code `0` so the command remains useful in pipes; a strict/failing mode can be added later if needed.

## Markdown output

Markdown should be compact and scannable:

- title/header with schema, root, file count, symbol count, max-files, and truncated status;
- one bullet per file;
- for Go files, include package/import summary and compact symbol bullets.

## Tests expected

- Go symbol extraction covers functions, methods, receiver names, types, consts, vars, exported flags, line numbers, imports, and parse-error degradation.
- Walker tests cover built-in ignores, `parley-deck/runs`, symlink/non-regular skipping where practical, sorted relative slash paths, and deterministic truncation.
- Renderer tests cover deterministic Markdown and JSON bytes.
- CLI tests cover usage, `--format json`, `--format markdown`, invalid format, and `--max-files`.

## Deferred follow-ups

- `--out` atomic writes.
- external ignore file or `.gitignore` compatibility.
- verbose exclusion diagnostics.
- exported-only filters, package filters, depth filters, or token budget filters.
- prompt wiring in `context-pack-wiring`.
- incremental indexing or cache keys.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. This keeps the MVP dependency-free, deterministic, and small enough to implement/review while preserving a stable JSON shape for later context-pack wiring.

### Signoff: claude — 2026-05-17
Status: ✅ ACCEPT
Notes: Schema-versioned, byte-deterministic, Go-stdlib only, with a clean three-layer split (walker / extractor / renderer) — matches my round-01 architecture and reviewability lens. I'd have preferred exit code 3 on truncation, but the `truncated` field plus exit 0 keeps the command pipe-friendly and a strict mode can be added later without breaking the schema. Deferring `--out`, gitignore compat, and verbose exclusion diagnostics is the right call for slice size.

### Signoff: gemini — 2026-05-17
Status: ✅ ACCEPT
Notes: The emphasis on byte-determinism and standard-library parsing aligns perfectly with context efficiency and project stability. Including exported flags and imports in the JSON schema provides the necessary "signal" for future agentic wiring without bloating the MVP.

### Signoff: hermes — 2026-05-17
Status: ✅ ACCEPT
Notes: From operations/resilience lens: byte-deterministic output, built-in ignores, symlink skipping, truncation with exit-0 pipe-friendliness, and Go-stdlib only all reduce runtime surprises and keep the MVP reviewable/operable. Good foundation.
