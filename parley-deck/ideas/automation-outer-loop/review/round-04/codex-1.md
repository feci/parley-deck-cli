---
agent: codex-1
idea: automation-outer-loop
review-round: 4
date: 2026-06-24
---

## Summary

REQUEST-CHANGES. Round-03 did close the exact slug symlink escape, the prompt-file symlink case, the non-directory slug case, the old frontmatter injection, and the 32-bit digest collision. I found one remaining AF10 boundary break: a symlink at `parley-deck/ideas` is still followed and causes `00-prompt.md` to be written outside the deck. I also found a smaller AF11 contract gap for C0 line separators in `Detail`.

The section 14 action brake still holds in the code path I inspected: `runLoopTick` only reads config/signals and calls `loop.Tick`; `loop.Tick` only validates/drafts candidates and has no run, push, merge, finalize, implementation, roster, or quorum path. The remaining MAJOR is about where the candidate file can be written, not about automatic promotion or execution.

## Refutation attempts

- Read `review/consensus.md` round-03, `IMPLEMENTATION.md` fix-up cycle 3, `git show d32b082`, current `internal/loop/loop.go`, and `internal/app/loop_cmd.go`.
- AF10 slug-dir symlink: planted `parley-deck/ideas/<slug> -> /tmp/target`; actual CLI exited 1 with the AF10 refusal and did not write the target.
- AF10 prompt-file symlink: planted real slug dir plus `00-prompt.md -> /tmp/outside-prompt.md`; `O_CREATE|O_EXCL` refused by treating it as an existing candidate, and the target was not written.
- AF10 non-directory slug: planted a regular file at `ideas/<slug>`; the path was rejected and the file was not modified.
- AF10 parent symlink: planted `parley-deck/ideas -> /tmp/outside-ideas`; actual CLI reported one candidate drafted and wrote `/tmp/outside-ideas/<slug>/00-prompt.md`. This is Finding 1.
- AF10 TOCTOU: I did not win a separate live race between `Mkdir`/`Lstat` and `OpenFile`, but the final prompt create is still pathname-based after validation. A same-path swap remains in the same class unless the fix anchors the final create to opened, no-follow directory handles.
- AF11/AF12: CR, CRLF, U+2028, U+2029, and U+0085 are normalized and the produced physical LF lines are indented. `TrimSpace` removes leading/trailing blank detail lines as claimed. However `Detail` containing vertical tab, form feed, and U+001C/U+001D/U+001E leaves heading/fence/frontmatter-looking text unindented under broad line splitting. This is Finding 2.
- AF12: no separate TrimSpace issue found.
- AF9/AF2: the previous 32-bit collision pair `probe-55599` and `probe-100565` now produced distinct 128-bit slug suffixes; `a/b` versus `a:b` also produced distinct slugs; the invalid `ci:` source was rejected while valid `ci` plus `:build` was drafted.
- Round-01 CRITICAL: end-to-end Source/ID/Title injection stayed closed. A hostile ID containing `status: round-01`, `participants: [evil]`, and `checks:` produced one frontmatter `status: candidate` key, no `participants:` key, and no `checks:` key.
- Tests: `go test -count=1 ./internal/loop ./internal/app` passed.

## Findings

### [MAJOR] AF10 still follows a symlinked `ideas/` parent

`internal/loop/loop.go:211` uses `os.MkdirAll(ideasDir, 0o755)` before the slug-level `os.Mkdir`/`Lstat` guard. `MkdirAll` accepts an existing `ideas` symlink to a directory. After that, `dir := filepath.Join(ideasDir, slug)` and `os.Mkdir(dir, ...)` operate through the symlinked parent, so the slug guard never sees a symlink at `ideas/<slug>`.

I reproduced this with the real CLI: with `parley-deck/ideas -> /tmp/outside-ideas`, `parley loop tick --enable` returned success and wrote `/tmp/outside-ideas/<slug>/00-prompt.md`. That violates the section 14 write boundary that loop-drafted candidates stay under `parley-deck/ideas/<slug>/`. It is the same security class as AF10, just one path component higher.

Concrete fix: validate/create `ideasDir` with the same no-follow rule as the slug dir: create it with `os.Mkdir`, and on `ErrExist` use `os.Lstat` to reject symlinks and non-directories before any slug work. Prefer a stronger helper that opens `ideasDir` and the slug directory with no-follow directory handles and creates `00-prompt.md` relative to that handle (`openat`/`O_NOFOLLOW` style), because that also closes the remaining path-swap window before `OpenFile`. Add a regression test for `parley-deck/ideas` as a symlink and assert the outside target remains empty.

### [MINOR] AF11 misses C0 line separators in `Detail`

`internal/loop/loop.go:312` normalizes CR, U+2028, U+2029, and U+0085, then splits only on `\n`. It does not normalize vertical tab (`\x0b`), form feed (`\x0c`), or the information separators U+001C/U+001D/U+001E. These are not live frontmatter breaks for the repo's current LF scanners, but they are treated as line boundaries by common broad splitters such as Python `splitlines()`.

Probe detail:
`ok\n## lf heading\r## cr heading\u2028## ls heading\u2029--- ps fence\u0085status: nel\x0b## vt heading\x0c--- ff fence\x1cstatus: fs\x1dparticipants: gs\x1echecks: rs`

Physical LF splitting had no unindented lines, but broad splitting produced unindented `## vt heading`, `--- ff fence`, `status: fs`, `participants: gs`, and `checks: rs`. That leaves AF11's "no column-0 heading/fence/key under a renderer-like line split" contract incomplete.

Concrete fix: before trimming/splitting, normalize every line-break-like control separator to `\n`: at least `\v`, `\f`, U+001C, U+001D, and U+001E in addition to the existing CR/NEL/LS/PS set. Consider mapping other non-tab C0 controls in `Detail` to spaces so free-form logs remain readable without carrying terminal/control semantics. Add a regression that checks the hostile `Detail` with a broad line-boundary set, not only `strings.Split(..., "\n")`.

## Open questions

None blocking beyond the two findings above.
