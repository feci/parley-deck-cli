---
agent: codex-1
idea: automation-outer-loop
review-round: 3
date: 2026-06-24
---

## Summary

I re-reviewed `14f8295` and the current `internal/loop/loop.go` / `internal/app/loop_cmd.go` in refutation mode. AF6-AF9 mostly hold: current repo parsers still see hostile candidates as `status: candidate` with no `participants:`, same-run and concurrent dedupe work with the 128-bit slug, and I found no `parley run`, push, merge, finalize, roster, or quorum-staffing path reachable from `parley loop tick`.

I did find two fix-up regressions/gaps: one MAJOR in the AF7 claim path, where a pre-existing symlinked candidate directory is followed and the prompt is written outside the canonical idea directory; and one MINOR in AF6, where Unicode line separators inside `Detail` are not normalized before indentation.

## Refutation attempts

- Read `review/consensus.md`, especially AF6-AF9 plus the dismissed/deferred list, and the `## Fix-up cycle 2` section of `IMPLEMENTATION.md`.
- Inspected `git show 14f8295`, current `cleanField`, `dedupeDigest`, `writeCandidate`, `indentDetail`, and `runLoopTick`.
- Re-confirmed the round-01 CRITICAL source/id/title frontmatter injection is closed: committed tests and CLI probes keep exactly one frontmatter `status: candidate` and no `participants:` / `checks:` keys.
- Tried hostile LF-delimited `Detail` values containing `---`, `status: round-01`, `participants: [evil]`, `## Injected heading`, a fake `## Promotion`, and `status: final`. The rendered lines are 4-space indented; `parley status --idea` still reports `Status: candidate` and empty participants.
- Checked downstream parser behavior: `protocol.ReadFrontmatter` and the driver-style `readFrontmatterField` stop at the first closing frontmatter fence and do not scan body fences or headings.
- Ran 20 compiled `parley loop tick --enable` processes concurrently against the same signal: exit failures 0, prompt files 1, created sum 1, skipped sum 19, stderr 0. I did not find a normal concurrent double-create or clobber.
- Checked same-run dedupe after AF9 with two signals sharing one explicit fingerprint: one prompt file, one created, one skipped. Slug length for the longest valid source was 44 chars (`loop-manual-` + 32 hex), which is acceptable.
- Checked `cleanField` against YAML 1.1 line breaks: LF/CR plus U+0085/U+2028/U+2029 are covered for frontmatter fields. I did not find another YAML 1.1 line break missing from `cleanField`.
- Ran focused tests: `go test -count=1 ./internal/loop ./internal/app ./internal/protocol ./internal/driver ./internal/runplan ./internal/consensus` passed.
- Ran the full suite: `go test -count=1 ./...` failed only in `internal/runner/TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`. I did not treat that as a loop finding, but it means I could not reproduce the implementation note's full-suite green in this environment.

## Findings

### [MAJOR] AF7 follows symlinked candidate directories and writes outside the idea path

`writeCandidate` now does `os.MkdirAll(dir)` and then `os.OpenFile(filepath.Join(dir, "00-prompt.md"), O_CREATE|O_EXCL, ...)` (`internal/loop/loop.go`). If `ideas/<slug>` already exists as a symlink to a directory, `MkdirAll` succeeds and `OpenFile` follows the parent symlink. I reproduced this through the compiled CLI: pre-create `parley-deck/ideas/loop-commit-d77b252477634790043c542314258e19 -> /tmp/target`, run a matching signal, and `loop tick` reports one candidate drafted while `/tmp/target/00-prompt.md` is created.

Why it matters: AF7 moved the durable claim to `00-prompt.md`, but the claim is only safe if the parent slug directory is the real candidate directory. This can make the automated loop write outside `parley-deck/ideas/<slug>/`, violating the "draft candidate prompts only" write boundary and reintroducing a path-level TOCTOU/symlink class that nearby code (`retro propose`) already fails closed on with `Lstat`.

Concrete fix: create/validate the slug directory without following a pre-existing symlink. For example, `MkdirAll` only the `ideas` parent, then `Lstat` `ideas/<slug>`; if absent, `os.Mkdir` that exact slug dir; if present, require a real directory and reject `ModeSymlink` or non-directory entries. Keep the `O_CREATE|O_EXCL` prompt-file claim after that so empty real dirs are healed and concurrent writers still serialize. Add a regression test mirroring `retroPropose`'s symlink test: a symlinked slug must return an error/skip and must not create `00-prompt.md` in the symlink target.

### [MINOR] AF6 does not indent after Unicode line separators inside Detail

`indentDetail` normalizes CR/LF and prepends four spaces to each `\n`-split line, but it does not normalize U+0085, U+2028, or U+2029. A detail like `safe\u2028## Unicode heading\u2029---\u0085status: round-01` is written as one physical Go line with only the initial four spaces. The current repo parsers still report `status: candidate` because they split only on `\n`, so this is not a live frontmatter/quorum bypass.

Why it matters: AF6's stated guarantee is that indentation keeps Detail literal so it cannot inject headings or fences. That guarantee is incomplete for Markdown/rendering/tooling that treats Unicode line separators as line breaks: after the separator, `## Unicode heading`, `---`, or `status:` can be interpreted as starting at column 0 rather than as part of the indented block.

Concrete fix: in `indentDetail`, normalize the same Unicode line separators to `\n` before `strings.Split`, preserving multiline readability while ensuring every logical line receives the four-space prefix. Add a test with U+0085/U+2028/U+2029 in `Detail` asserting the prompt contains normalized, individually indented lines and no raw Unicode separators in the detail block.

## Open questions

- Is the `internal/runner` durable-kill full-suite failure expected on this host? It appears unrelated to AF6-AF9, but it prevented an exact reproduction of the claimed `go test -count=1 ./...` green run.
