---
idea: launch-mkdir-resilience
phase: review-consensus
date: 2026-06-07
drafter: claude
participants: [claude, codex, agy, hermes]
---

## Review consensus

### Round-01 review verdicts
- **hermes — ACCEPT** (only NITs).
- **agy — ACCEPT** (only a documentation NIT).
- **codex — REQUEST-CHANGES** (1 MAJOR + 1 MINOR; no CRITICAL).

All three confirmed: the helper matches FINAL, every call-site swap is correct and on the
intended launch/run path, the out-of-scope sites (`InitWorkspace`, the `os.MkdirTemp`
isolated homes) are correctly left as `os.MkdirAll`, imports are clean, and
`internal/protocol` importing `internal/fsutil` creates no cycle.

### Agreed fixes (applied in fix-up cycle 1)
1. **codex MINOR — test hardening** (`internal/fsutil/fsutil_test.go`):
   - `Test_GenuineFailure` returns a distinct error on the final attempt and asserts that
     specific last error is returned.
   - New `Test_DirExistsBeatsPermission` locks the ordering: `isDir` success wins before the
     permission fail-fast.
   - `Test_NonDirCollision` uses `fs.ErrExist` + a regular file to prove `fs.ErrExist` is
     never trusted blindly.
   → 8/8 fsutil tests green.
2. **hermes + agy NIT** — inline comment on the leading `0` in `retryDelays`
   (`internal/fsutil/fsutil.go`) marking the immediate first retry.

### Dismissed finding (with evidence)
- **codex MAJOR — `go test ./...` not green (`TestDurableKillEndToEndRealProcess`).**
  **Dismissed: false positive from codex's sandboxed execution, not a defect and unrelated
  to this change.** Evidence: the test passes in a normal dev shell both WITH this change
  (0.354s) and on the ORIGINAL code (`git stash`, 0.422s); a fresh `go test ./... -count=1`
  is fully green (0 FAIL). Root cause of codex's failure: the procctl fail-closed kill gate
  (shipped 1.19.0) needs the darwin boot id via `sysctl kern.boottime`, which the seatbelt
  sandbox of `codex exec` restricts → "no recorded boot id" → the gate correctly refuses to
  kill → a test that expects a kill fails. This is the safety gate working as designed under
  a restricted sandbox, not a regression. The release green claim is verified in the
  supported environment.

### Deferred follow-ups (not v1)
Auto-attach to an orphaned run on launch error; untruncated launch-error text + `/open`
hint; adopting the helper at non-launch mkdir sites if they ever prove flaky.

### Status
Zero agreed fixes remain after cycle 1 (the MINOR + NIT are applied; the MAJOR is a
dismissed sandbox artifact). Re-review confirms below.

## Signoffs

### claude — ACCEPT
Facilitator/implementer. Fix-up cycle 1 applied the agreed MINOR (test hardening) and the
NIT (retryDelays comment); the MAJOR is a documented sandbox artifact (verified green in a
normal shell, with and without the change). Zero agreed fixes remain. Marking complete.

### codex — ACCEPT (review/round-02)
Round-01 MINOR addressed (distinct-last-error, isDir-before-permission, fs.ErrExist not
masked — 8/8 fsutil tests). Accepts the dismissal of the round-01 MAJOR as a sandbox
restriction on `sysctl kern.boottime` (procctl fail-closed), not a code defect. Verdict:
ACCEPT. (Authored in review/round-02/codex.md.)

### hermes — ACCEPT (review/round-01, carried forward)
Round-01 verdict ACCEPT with only NITs; the retryDelays-comment NIT was applied in fix-up.
No code-logic change in the fix-up affects the round-01 review. (Authored in
review/round-01/hermes.md.)

### agy — ACCEPT (review/round-01, carried forward)
Round-01 verdict ACCEPT with only a documentation NIT, applied in fix-up. (Authored in
review/round-01/agy.md.)

## Status: COMPLETE
Zero CRITICAL/MAJOR/MINOR remain; all participants ACCEPT. Ready to ship v1.21.0.
