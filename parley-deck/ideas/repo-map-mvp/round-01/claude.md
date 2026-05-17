---
agent: claude
idea: repo-map-mvp
round: 1
date: 2026-05-17
---

## Summary

Ship `parley context map` as a thin command over a new `internal/repomap` package that has three reviewable layers: walker, language extractor (Go-only for this slice), and renderer (Markdown + JSON). Keep the data model schema-versioned and explicit about what is intentionally absent so later slices (`context-pack-wiring`) can wire it without renaming fields. Optimize for a small, boring diff that a reviewer can hold in their head.

## Proposed approach

### 1. Command shape and outputs

- Command group `context`, single subcommand `map`. Matches the existing switch-based dispatch in `internal/app/app.go:41`.
  - `parley context map [--root PATH] [--format md|json] [--out PATH] [--max-files N] [--include-tests] [--ignore-from PATH]`
  - Defaults: `--root .`, `--format md`, stdout, `--max-files 2000`, exclude `*_test.go` by default off (i.e. tests included; trade-off below), no extra ignore file.
- Exit codes: `0` success, `2` usage error, `3` file-limit exceeded with a clear stderr message naming the cap (no truncation surprise).
- Stdout-only writes for piping; `--out` writes atomically (write-to-temp + rename) so partial runs never leave a half-rendered file. This matches reviewer expectations from existing writers.
- No global flags added; everything lives inside `runContextMap` in `internal/app/context_map.go` to match the file-per-command pattern already used (`consensus_request_signoffs.go`, `version_status.go`).

### 2. Minimal data model

One file, `internal/repomap/model.go`, with a schema version baked in from day one so consumers can branch on it:

```go
type Map struct {
    SchemaVersion int     `json:"schema_version"` // start at 1
    Root          string  `json:"root"`           // repo-relative, always "."
    GeneratedAt   string  `json:"generated_at,omitempty"` // optional, omit by default for determinism
    Files         []File  `json:"files"`
    Truncated     bool    `json:"truncated"`      // true if file limit hit
    Counts        Counts  `json:"counts"`
}

type Counts struct {
    Files   int `json:"files"`
    Symbols int `json:"symbols"`
}

type File struct {
    Path     string   `json:"path"`               // forward-slash, repo-relative
    Lang     string   `json:"lang"`               // "go" | "" (unknown)
    SizeB    int64    `json:"size_b"`
    SHA1     string   `json:"sha1,omitempty"`     // optional, off by default
    Symbols  []Symbol `json:"symbols,omitempty"`  // omitted for non-Go in this slice
}

type Symbol struct {
    Kind     string `json:"kind"`     // "func" | "method" | "type" | "const" | "var"
    Name     string `json:"name"`
    Recv     string `json:"recv,omitempty"` // method receiver type, no pointer star
    Exported bool   `json:"exported"`
    Line     int    `json:"line"`
    Doc      string `json:"doc,omitempty"` // first line of leading doc comment, trimmed
}

```

Trade-offs:
- Including `SchemaVersion` and `Truncated` from v1 costs almost nothing and avoids a breaking change when `context-pack-wiring` lands.
- Omitting `GeneratedAt` from default output preserves byte-for-byte determinism, which makes golden-file tests trivial. Add it later behind a flag if humans ask.
- `SHA1` left as optional/off keeps the slice small; later slices that want cache keys can flip it on.

### 3. Default ignore rules and file-limit behavior

- Built-in ignore set (literal directory names anywhere in the path): `.git`, `node_modules`, `vendor`, `dist`, `build`, `.next`, `.venv`, `__pycache__`, `.idea`, `.vscode`, `target`, `parley-deck/ideas/*/round-*` excluded only when not the root being mapped. Trade-off: hard-coding is uglier than a config file but is reviewable in one screen; we can graduate to a config later.
- Skip files where size > 1 MiB (configurable via `--max-file-bytes`, default 1 MiB).
- `--ignore-from` reads a `.gitignore`-style file but only honors plain prefix/suffix and exact-name patterns (no negation, no `**` outside leading position). State this explicitly in `--help`; a partial gitignore parser is the largest realistic source of "why was this excluded?" reviewer questions.
- File limit: hard cap at `--max-files` (default 2000). On exceed: stop walking, set `Truncated: true`, render with a clear footer/banner, and exit `3`. Deterministic order means truncation is reproducible (walk in `filepath.WalkDir` order then sort before rendering, so the truncation set is the first N by sorted path — not by walk order).

### 4. Tests for confidence

- `internal/repomap/walker_test.go`: fixtures under `testdata/walker/` covering ignore rules, symlink loops (skip with no follow), size cap, and `--max-files` truncation determinism.
- `internal/repomap/gosyms_test.go`: a small Go file with exported/unexported funcs, a method, a type, a const block, a var, and one with a doc comment. Asserts exact `Symbol` slice (line numbers included).
- `internal/repomap/render_test.go`: golden files for both Markdown and JSON. JSON test compares bytes; Markdown test compares bytes too — determinism is the headline feature.
- `internal/app/app_test.go` (extend): one end-to-end test that runs `parley context map --root testdata/repo --format json` on a tiny embedded repo and diffs against a golden.
- Negative tests: malformed Go file does not abort the run; it produces a `File` with empty `Symbols` and the run continues. Reviewer sanity check.
- Cross-platform path test: assert all emitted paths use `/` even on Windows-style separators (use `filepath.ToSlash` consistently).

### 5. Architecture and reviewability notes

- Package layout (each file < ~200 LOC, suggested):
  - `internal/repomap/model.go` — types only.
  - `internal/repomap/walk.go` — filesystem walk + ignore + limits, returns `[]File` with `Lang` set, no symbols.
  - `internal/repomap/gosyms.go` — `func ExtractGo(src []byte) ([]Symbol, error)` using `go/parser` + `go/ast`.
  - `internal/repomap/render_md.go` and `render_json.go` — pure functions `Render(Map, io.Writer) error`.
  - `internal/repomap/build.go` — top-level `Build(opts Options) (Map, error)` that composes the above.
- The CLI layer (`internal/app/context_map.go`) is a thin flag parser plus one call to `repomap.Build` and one call to a renderer. Reviewer can read the CLI in one page.
- Concurrency: do not parallelize Go parsing in this slice. Single-goroutine keeps determinism trivial and the diff small. Parallelism is a measured follow-up if benchmarks justify it.
- Logging: nothing on stdout besides the rendered map; warnings (e.g. unparseable Go file) go to stderr, prefixed `parley context map:`, suppressed by default and shown with `--verbose`. Avoids polluting pipes.

## Concerns / open questions

- Should `parley-deck/ideas/**/round-*/*.md` be ignored by default? They are noisy and self-referential. Lean yes, but flag this for codex/hermes — operational impact.
- Tests-included default: including `*_test.go` symbols makes the map noisier but more honest for reviewers reading "what's in this package." Acceptable for MVP; revisit when a real consumer (context-pack-wiring) has a budget signal.
- `Doc` field truncation policy: first line only, trimmed, max 200 chars. Cheap, deterministic, good enough for human scan; revisit when an LLM consumer arrives.
- Path canonicalization for the `Root`: always `"."` in output, even if user passes an absolute path. Avoids leaking developer-machine paths into committed golden files.

## Risks

- **Scope creep into a tree-sitter clone.** Mitigation: lock `gosyms.go` to `go/parser` only this slice; explicit non-goal in the PR description.
- **Determinism regressions** from map iteration, locale-sensitive sort, or `filepath.Walk` ordering. Mitigation: sort all slices by `path` (and by `(line, name)` within a file) immediately before render; assert via golden tests on CI Linux and macOS.
- **Ignore-rule confusion** ("why did my file get skipped?"). Mitigation: `--verbose` lists every excluded path with the reason. Cheap to add, expensive to debug without.
- **Hidden coupling to later slices.** `context-pack-wiring` will want stable field names and a way to filter by package. Mitigation: schema-version the JSON; do not add filtering in this slice but document the intended extension point in the package doc comment.
- **Performance cliff on monorepos.** Mitigation: the file-limit-plus-exit-3 contract makes the failure mode loud and reproducible rather than a silent slow build; tune defaults after the first real user.
- **Symlink cycles / non-regular files** causing infinite walks. Mitigation: `filepath.WalkDir` with no symlink-follow, skip non-regular files; covered by walker_test.
