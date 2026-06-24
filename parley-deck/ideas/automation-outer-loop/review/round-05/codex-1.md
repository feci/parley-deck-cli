---
agent: codex-1
idea: automation-outer-loop
review-round: 5
date: 2026-06-24
---

## Summary

I re-reviewed fix-up cycle 4 in refutation mode against `git show 580600a`, the current
`internal/loop/loop.go` (`safeMkdir`, `assertInsideDeck`, `writeCandidate`, `indentDetail`,
`cleanField`, `dedupeDigest`), and `internal/app/loop_cmd.go`.

I found no remaining issues: 0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT. AF14 and AF15 held under
the attacks I tried. The round-01 CRITICAL frontmatter/quorum injection remains closed, the
128-bit digest is collision-resistant for this use, and COOPERATION.md §14 still holds: the
loop path drafts candidate prompts only and has no run, push, merge, finalize, or quorum path.

Validation: `go test -count=1 ./internal/loop ./internal/app` passed. `go test -count=1 ./...`
failed only in the pre-existing environment-sensitive `internal/runner`
`TestDurableKillEndToEndRealProcess` case (`no recorded boot id`), not in the loop/app path.

## Refutation attempts

### AF14 symlink and path containment

I tried to make `parley loop tick` write `00-prompt.md` outside the intended candidate tree
using temporary decks and a scratch-built CLI binary.

- `parley-deck/ideas/` as a symlink to an outside target: rejected by `safeMkdir`; the target
  stayed empty.
- `parley-deck/ideas/<slug>/` as a symlink to an outside target: rejected by `safeMkdir`; the
  target stayed empty.
- Dangling symlink at `parley-deck/ideas/`: rejected before any target was created.
- Relative/absolute-looking signal data (`../../..` in `id`, absolute-looking fingerprint):
  produced a digest-only slug with no slash and no path traversal.
- `parley-deck/` itself as a symlink, and a symlink above the workspace root: the prompt was
  written inside the resolved deck root. I treat this as intended base-path resolution rather
  than an escape; `assertInsideDeck` uses the resolved deck as the trust anchor.
- `EvalSymlinks` edge cases: `writeCandidate` creates `deck`, `ideas/`, and `ideas/<slug>/`
  before `assertInsideDeck`, so the checked slug path exists. Dangling `ideas`/slug symlinks
  are rejected before `EvalSymlinks`. A `Rel == "."` result is not reachable through the
  sanitized slug plus real-directory creation path; `Rel == ".."` and `../...` are rejected.
- TOCTOU swap: I did not reproduce a post-`assertInsideDeck` swap into an outside symlink.
  Such a race would require a concurrent writer replacing checked directories during the tiny
  assert-to-open window. I am not raising it without a concrete repro under this threat model.

### AF15 Detail indentation

I generated a `Detail` payload with hostile `##`, `---`, `status:`, `participants:`, and
`checks:` tokens after LF, lone CR, vertical tab, form feed, U+001C, U+001D, U+001E, U+0085,
U+2028, and U+2029. Reading the produced prompt with Python-style broad `splitlines()` found
only the two expected frontmatter fences and no unindented hostile tokens. The current
`indentDetail` normalizes every tested line-break-like separator to `\n`, maps other non-tab
C0 controls to spaces, trims boundary blank lines, and prefixes every resulting line.

### Original CRITICAL, digest, and §14 brake

The round-01 CRITICAL remains closed: `Source`, `ID`, and `Title` are flattened before
frontmatter/body scalar use; `Source` must be in the closed set; `Detail` is below the closed
frontmatter and indented as literal body content. I found no way for signal content to add
frontmatter `participants:`, `status: round-01`, or `checks:`.

`dedupeDigest` now uses a 32-hex-character (128-bit) SHA-256 prefix over `strconv.Quote`d
fields, with explicit fingerprint separated from source/id identity. I did not attempt an
infeasible 128-bit collision; the prior 32-bit birthday-collision class is closed.

For §14, `internal/loop` depends only on stdlib plus `internal/protocol` and
`internal/fsutil`; it does not import or call driver, runner, process, git, push, merge,
finalize, or run code. The app dispatch is `loop` -> `runLoopTick` -> `loop.Tick`, and the
only artifact written is a `status: candidate` prompt with no frontmatter `participants:`.

## Findings

No findings.

## Open questions

None.
