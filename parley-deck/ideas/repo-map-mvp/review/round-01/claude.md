---
agent: claude
idea: repo-map-mvp
review-round: 1
date: 2026-05-17
reviewed-commit: b11bc19
implementation-pr: https://github.com/feci/parley-deck-cli/pull/22
---

## Summary

Implementation matches the spirit and structure of FINAL.md: the
`parley context repo-map` command lands with documented defaults, a deterministic
walker, Go AST-based symbol extraction, stable Markdown and JSON renderers, and
no external parser dependencies. The JSON schema mirrors FINAL byte-for-byte
(`schema_version`, `root: "."`, `max_files`, `truncated`, `counts`, `files`),
parse errors degrade gracefully, and built-in directory ignores plus
`parley-deck/runs` are honored. Findings below are limited to a small
determinism leak in one error path and modest gaps in test coverage relative to
the FINAL's "Tests Required" list.

## Findings

### [MINOR] Absolute path may leak into `parse_error` via `os.ReadFile`

`internal/repomap/repomap.go` builds the absolute filesystem path with
`filepath.Join(rootAbs, ...)` and, if `os.ReadFile` fails, stores
`err.Error()` directly in `file.ParseError`. The standard library's
`*os.PathError` embeds the absolute path it was called with, so a read failure
(e.g. a file removed or chmod'd between walk and read) renders the developer's
absolute root into the JSON/Markdown body. FINAL explicitly requires:

> no absolute developer-machine root in rendered bodies

Why it matters: the file is emitted to stdout and consumed by downstream
tooling that the FINAL declares "byte-deterministic for fixed inputs"; an
absolute path makes the output host-specific and breaks reproducibility.

Concrete fix: redact the path before storing the error, e.g.
`file.ParseError = fmt.Sprintf("read %s: %v", file.Path, errors.Unwrap(err))`
after stripping `rootAbs`, or simply
`file.ParseError = "read error: " + filepath.Base(... )`-style message that
references `file.Path` (already slash-relative) rather than the wrapped OS
error. The same defensive treatment can be applied uniformly to the
`parser.ParseFile` branch, even though it is currently safe because the
relative `file.Path` is passed as filename.

### [MINOR] Walker tests omit symlink and non-regular handling required by FINAL

FINAL's "Tests Required → Walker" section lists:

> non-regular and symlink handling where practical

`internal/repomap/repomap_test.go` covers built-in ignores, `parley-deck/runs`,
slash-normalized paths (implicitly, via `a.txt`/`b.txt`/`c.txt`), and
deterministic truncation, but never creates a symlink (`os.Symlink`) or a
non-regular file. The walker code paths `entry.Type()&os.ModeSymlink != 0` and
`!info.Mode().IsRegular()` are therefore untested.

Why it matters: the FINAL-mandated guarantees that symlink targets are skipped
and that only regular files are included are silently regression-prone. A
later refactor of the walker (or someone enabling symlink following) would not
be caught by the suite.

Concrete fix: add a test that creates `os.Symlink(target, link)` inside the
root and asserts the link does not appear in `m.Files`, plus (where the
platform allows) a `syscall.Mkfifo` or other non-regular entry guarded by
`runtime.GOOS != "windows"`, with a `t.Skip` fallback if creation is not
permitted (per FINAL's "where practical" qualifier).

### [MINOR] CLI "usage" path is not covered by tests

FINAL's "Tests Required → CLI" lists `usage` first.
`internal/app/app_test.go` exercises `--format json`, `--format markdown`,
invalid format, and `--max-files`, but never calls `Run([]string{"context"})`
or otherwise asserts that the usage path returns exit code 2 and writes the
`usage: parley context repo-map ...` string to stderr.

Why it matters: `printContextUsage` and the bare-`context` branch in
`runContext` are silent dead code from the test suite's perspective; a typo
that breaks the usage message (or its routing) ships green.

Concrete fix: add a small test that runs `Run([]string{"context"}, ...)` and
asserts `code == 2` plus `strings.Contains(stderr, "usage: parley context
repo-map")`. Optionally also assert the same for an unknown subcommand like
`context bogus`.

### [NIT] Undocumented `--format md` alias

`internal/app/context.go` accepts both `markdown` and `md`:

```go
case "markdown", "md":
```

FINAL only documents `markdown|json`, the in-binary usage string only lists
`markdown|json`, and no test exercises `md`. It is a harmless extra, but it
either belongs in the documented surface (FINAL + usage string + test) or
should be dropped to keep the CLI surface exactly aligned with FINAL.

Concrete fix: remove the `"md"` arm to match FINAL, or, if the alias is
intentional, raise it in a follow-up and update the usage string and tests.

## Open questions

- Should `repomap.Build` treat `--max-files 0` as "include nothing" (and
  perhaps `--max-files < 0` as an error), or is the current silent fallback
  to `1000` the desired UX? FINAL does not specify, but the current behavior
  means `--max-files 0` produces a 1000-file output, which could surprise a
  caller probing for an empty-set degenerate case.
- The `parse_error` field uses `omitempty` (as do `package`, `imports`,
  `symbols`). FINAL lists these as part of the file shape without specifying
  whether they must always be emitted. Is the current "omit when empty" the
  intended JSON contract, or should the schema be made strictly fixed-shape
  for downstream consumers?
