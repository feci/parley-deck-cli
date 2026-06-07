---
idea: launch-mkdir-resilience
phase: implementation
date: 2026-06-07
implementer: claude
status: ready-for-review
---

## What was implemented

Per FINAL.md: a stdlib-only resilient mkdir helper that tolerates spurious/transient
failures from the weakly-coherent `AppleVirtIOFS` shared mount, swapped in across the
new-idea launch path and the live run path. The healthy path is unchanged (exactly one
`os.MkdirAll`, no Stat, no sleep).

## Files changed

**New — `internal/fsutil/fsutil.go`**
- `func MkdirAllResilient(path string, perm os.FileMode) error`.
- First `os.MkdirAll`; on success return immediately (no Stat, no sleep).
- On error: loop `retryDelays = [0, 5ms, 20ms, 50ms]` (5 total attempts) — before each
  retry, (a) `isDir(path)` via a fresh `os.Stat`+`IsDir` → success; (b) `errors.Is(err,
  fs.ErrPermission)` → return err (fail-fast); (c) sleep the delay (if > 0) then retry.
- After the loop, a final `isDir`; else return the last `os.MkdirAll` error.
- Success oracle is always a fresh Stat showing a directory — `fs.ErrExist` is never
  trusted blindly. Seams `mkdirAll`/`stat`/`sleep` are package vars for tests.

**New — `internal/fsutil/fsutil_test.go`** (7 tests, all green, instant — `sleep` stubbed):
common-path (1 mkdir, 0 stat, 0 sleep); transient-then-success (2 mkdir, 0 sleep);
host-succeeded-guest-lied (mkdir always errors, dir pre-exists → nil via isDir, no retry/
sleep); already-exists; genuine-failure (returns last mkdir error after retries; sleeps
exactly 5/20/50ms); fail-fast-permission (returns immediately, 1 mkdir, 0 sleep);
non-directory-collision (file at path → error, not masked).

**Call-site swaps** `os.MkdirAll` → `fsutil.MkdirAllResilient` (launch + live run path):
- `internal/protocol/workspace.go:67` — `CreateIdea` idea/round-01 dir.
  **`InitWorkspace` (line 44) intentionally left as `os.MkdirAll`** (`parley init`
  bootstrap, not the per-launch hot path), per FINAL scope.
- `internal/store/events.go:41` — `Store.Append` events dir (the proven failure point).
- `internal/runmanifest/manifest.go:142` — manifest dir.
- `internal/runner/runner.go:340,343` — agentDir + output dir.
- `internal/runner/round_index.go:85` — roundDir.
- `internal/runner/steer.go:152` — steerDir.
- `internal/runner/phase58.go:58` — fix-up agentDir.
- `internal/runner/handoff.go:36` — handoff agentDir.
- `internal/hitl/hitl.go:183` — questions dir.

**Version:** `internal/app/version.go` `1.20.0` → `1.21.0`; root `VERSION` → `1.21.0`
(`TestVersionFileMatchesBinaryVersion`).

## Unchanged (per consensus)
- The trailing `os.OpenFile(events.jsonl, O_CREATE|O_APPEND)` in `Store.Append`.
- TUI / launch error handling (no auto-attach, no message change in v1).
- The `os.MkdirTemp` isolated-home dirs (`runner.go:861,901`) and all non-launch/non-run
  `os.MkdirAll` sites (`pipeline/*`, `consensus/*`, `sessionstore`, `driver/*`, the
  `app.go:1855` probe).

## Deviations from FINAL
None of substance. The retry loop is expressed exactly as the FINAL pseudocode; `os`
remains imported in every edited file (each uses `os` elsewhere), so no import churn beyond
adding `parley-deck-cli/internal/fsutil`.

## Verification
- `gofmt -l` on all edited files: clean.
- `go build ./...`: OK. `go vet ./...`: OK.
- `go test ./...`: all green (7/7 new fsutil tests; full suite incl. app/runner/store/
  runmanifest/hitl/protocol).

Ready for Phase 6 review.

## Fix-up cycle 1 (Phase 8) — responses to review/round-01

- **codex MINOR (test hardening):** applied in `internal/fsutil/fsutil_test.go`:
  - `Test_GenuineFailure` now returns a **distinct** error on the final attempt and asserts
    the returned error is that last one (not an earlier `errBoom`).
  - Added `Test_DirExistsBeatsPermission`: `mkdirAll → fs.ErrPermission` + `stat →` existing
    dir → returns nil, proving the `isDir` check wins **before** permission fail-fast.
  - `Test_NonDirCollision` now returns `fs.ErrExist` (not a generic error) with a regular
    file at the path → still an error, proving `fs.ErrExist` is never trusted blindly.
  Tests are now 8/8, all green.
- **hermes/agy NIT (retryDelays comment):** added an inline `/* immediate first retry, no
  sleep */` on the `0` element in `internal/fsutil/fsutil.go`.
- **codex MAJOR (`go test ./...` not green — `TestDurableKillEndToEndRealProcess`):**
  **dismissed with evidence.** The test passes in a normal dev shell both WITH this change
  (0.354s) and on the ORIGINAL code via `git stash` (0.422s), and a fresh
  `go test ./... -count=1` is fully green (0 FAIL). codex's failure is an artifact of its
  sandboxed `codex exec` execution: the procctl fail-closed kill gate (from 1.19.0) needs
  the darwin boot id via `sysctl kern.boottime`, which the seatbelt sandbox restricts →
  "no recorded boot id" → the gate correctly refuses to kill → the test that expects a kill
  fails. It is pre-existing, environment-specific, and unrelated to mkdir-resilience.

Re-verified after fix-up: `gofmt` clean, `go build/vet ./...` OK, `go test ./... -count=1`
green (0 FAIL). Ready for re-review.
