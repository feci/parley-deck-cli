---
agent: claude
idea: repo-map-mvp
review-round: 2
date: 2026-05-17
reviewed-commit: 07b26ad
implementation-pr: https://github.com/feci/parley-deck-cli/pull/22
responding-to: [claude/review/round-01, gemini/review/round-01, hermes/review/round-01, review/consensus]
---

## Summary

Round-02 verdict: **✅ ACCEPT**. The previous process blocker (Bash EPERM in
this session, working tree on `main` where the repomap files do not exist) is
**resolved by inline patch review** — the facilitator supplied the full diff
for `07b26ad` plus an independently-run test transcript, so I can now read the
four fixes against the consensus list without shell access.

The patch is small (4 files, +67/−2 lines) and every line maps directly to one
of the four agreed consensus fixes. I verified each fix against the inline
diff below. The facilitator's local run confirms `gofmt`, `git diff --check`,
`go test ./internal/repomap ./internal/app`, and `go test ./...` all pass, and
that `go run ./cmd/parley context repo-map --dir . --format md` now exits
non-zero with `invalid format "md"; expected markdown or json` — matching
consensus fix #4 exactly. No new code-level findings.

## Fix Verification

Status legend: ✅ verified, ⚠️ self-reported only, ❌ regressed.

1. **Sanitize `parse_error` so the absolute developer-machine root cannot leak
   via `os.ReadFile`** (claude/round-01 MINOR #1, consensus fix #1).
   - Status: ✅ verified (via inline diff).
   - `internal/repomap/repomap.go` now routes the read failure through a new
     `readErrorMessage(path, err)` helper that uses `errors.As` to peel off
     `*os.PathError` and formats `"read %s: %v"` against the relative
     `file.Path` and the underlying `pathErr.Err` — so the wrapped absolute
     path inside `*os.PathError.Path` never reaches `file.ParseError`. The
     non-`*os.PathError` branch also uses the relative `path`, so a wrapped
     error from a different read codepath cannot bypass the sanitisation.
   - `internal/repomap/repomap_test.go` adds `TestReadErrorMessageUsesRelativePath`,
     which constructs a `*os.PathError{Path: "/tmp/private/repo/broken.go"}`
     and asserts the formatted message neither contains `/tmp/private/repo`
     nor drops `broken.go`. That's a tight contract test for the leak I flagged
     in round-01.
   - Nit (non-blocking): the test is a direct unit test on `readErrorMessage`
     rather than an end-to-end repro that forces `enrichGoFile` to fail mid-walk,
     but the helper is the only write site for `file.ParseError`, so the
     contract is adequate.

2. **Add walker coverage for symlinks and practical non-regular paths**
   (claude/round-01 MINOR #2, consensus fix #2).
   - Status: ✅ verified (via inline diff), with one NIT below.
   - `TestBuildSkipsSymlinksAndDirectories` creates a `target.txt` regular
     file, a `plain-dir` directory, and (on non-Windows) a symlink
     `linked.txt -> target.txt`, then runs `Build` and asserts that neither
     `linked.txt` nor `plain-dir` shows up as an emitted `File`. The symlink
     branch is gated by `runtime.GOOS != "windows"` and falls back to
     `t.Skipf` on symlink-creation errors, which is the right portability
     posture.
   - This covers the symlink branch I asked for in round-01 and exercises a
     non-file entry (directory) as the "practical non-regular" case. It does
     not exercise a true non-regular `os.ModeType` bit such as a FIFO/socket,
     but FINAL says "where practical" and a directory IS non-regular under
     `Mode().IsRegular()`, so this is a defensible reading of the spec — same
     judgment the implementer telegraphed in IMPLEMENTATION.md cycle 1.

3. **Add CLI usage coverage for bare `parley context` and an unknown
   `parley context <bogus>` subcommand** (claude/round-01 MINOR #3,
   consensus fix #3).
   - Status: ✅ verified (via inline diff).
   - `internal/app/app_test.go` adds `TestContextUsage` covering both
     branches: `Run([]string{"context"}, ...)` and
     `Run([]string{"context", "bogus"}, ...)`. Both assert exit code 2 and
     `usage: parley context repo-map` on stderr. That matches the round-01
     MINOR #3 ask exactly (both the no-subcommand and the unknown-subcommand
     paths print usage to stderr and exit non-zero).

4. **Drop the undocumented `--format md` alias** (claude/round-01 NIT,
   consensus fix #4).
   - Status: ✅ verified (via inline diff + smoke check).
   - `internal/app/context.go` collapses `case "markdown", "md":` to
     `case "markdown":`, which routes `--format md` to the `default` arm. The
     facilitator's `go run ./cmd/parley context repo-map --dir . --format md`
     transcript confirms the resulting message is
     `invalid format "md"; expected markdown or json` and exit code is
     non-zero — exactly the round-01 NIT formulation.
   - No dedicated regression test was added for the rejection itself, but the
     existing markdown-and-validation test plus the live smoke check are
     enough at this slice; adding a `--format md` assertion is trivial and can
     ride along with future CLI-validation tests if desired.

The two consensus deferrals (`--max-files 0` semantics and JSON `omitempty`
policy) are appropriately untouched in this commit. No regressions, no
unrelated drift, and the test suite green per the facilitator's transcript.

## Findings

No new CRITICAL, MAJOR, or MINOR findings.

### [NIT] Non-regular-file coverage in `TestBuildSkipsSymlinksAndDirectories` is symlink + directory only

- The test name and body cover the symlink branch (the main one I called out
  in round-01) plus a directory as the "non-regular" stand-in. It does not
  exercise a true `os.ModeType` non-regular file such as a FIFO or socket,
  so the `!info.Mode().IsRegular()` walker branch is only proven for the
  directory case in tests.
- FINAL's "where practical" qualifier permits this scoping, and the symlink
  case is the one with realistic developer-machine exposure (cloned repos
  with checked-in symlinks). I'm flagging this as a NIT, not as a fix
  request — if a future change tightens or generalises the walker's mode
  filter, a FIFO/socket test (likely gated by `runtime.GOOS` and
  `syscall.Mkfifo` availability) would be the natural follow-up.

## Open questions

- Resolved: the round-01-style per-fix verification standard is now applied
  to `07b26ad` via inline-patch review, matching the standard the other
  participants applied to `b11bc19`. No further verification is needed for
  this slice.
- Carry-forward (not a round-02 action): the `parse_error` / `package` /
  `imports` / `symbols` `omitempty` contract and the `--max-files 0` UX
  remain consensus-deferred. I agree they should stay out of this slice.
