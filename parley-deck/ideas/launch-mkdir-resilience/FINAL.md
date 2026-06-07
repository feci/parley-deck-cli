---
idea: launch-mkdir-resilience
phase: final
date: 2026-06-07
drafter: claude
implementer: claude
status: final
participants: [claude, codex, agy, hermes]
---

## Problem

`parley tui` → **N** (new idea) intermittently shows `launch failed: mkdir
…/runs/<runID>: <errno>` even though the run directory is created and the run actually
runs in the background (orphaned). Root cause: the deck lives on an `AppleVirtIOFS`
(virtio-fs) shared mount with weak cross-process cache coherence; `os.MkdirAll` in
`Store.Append` returns a transient `ENOENT`/`ENOTDIR` from a stale dentry/attribute cache
while the directory is (or becomes) present on the host. `runcontrol.Create` propagates the
spurious error and the TUI aborts the launch.

## Solution

A small, stdlib-only resilient mkdir helper whose **healthy path is exactly one
`os.MkdirAll`** (zero overhead), and which on error treats "the directory exists" as
success and retries a small bounded number of times before surfacing the real error. Swap
it in on the new-idea launch path and the live run path.

### Helper — `internal/fsutil/fsutil.go`

```go
package fsutil

import (
    "errors"
    "io/fs"
    "os"
    "time"
)

// Seams (overridable in tests only).
var (
    mkdirAll = os.MkdirAll
    stat     = os.Stat
    sleep    = time.Sleep
)

// retryDelays are the sleeps before each retry after the initial attempt. The first retry
// is immediate (the retry itself often forces a virtio-fs cache revalidation); later
// retries give a weakly-coherent attribute cache a brief window to settle. 5 total
// MkdirAll attempts, worst-case added latency 75ms — only ever on the error path.
var retryDelays = []time.Duration{0, 5 * time.Millisecond, 20 * time.Millisecond, 50 * time.Millisecond}

// MkdirAllResilient behaves like os.MkdirAll but tolerates spurious/transient failures
// from weakly-coherent filesystems (virtio-fs, NFS, SMB) whose dentry/attribute cache can
// be momentarily stale across processes. The healthy path is a single os.MkdirAll with no
// Stat and no sleep. On error it returns nil whenever a fresh Stat shows the path is a
// directory, fails fast on permission errors, and otherwise retries before returning the
// last os.MkdirAll error.
func MkdirAllResilient(path string, perm os.FileMode) error {
    err := mkdirAll(path, perm)
    if err == nil {
        return nil
    }
    for _, d := range retryDelays {
        if isDir(path) {
            return nil
        }
        if errors.Is(err, fs.ErrPermission) {
            return err
        }
        if d > 0 {
            sleep(d)
        }
        if err = mkdirAll(path, perm); err == nil {
            return nil
        }
    }
    if isDir(path) {
        return nil
    }
    return err
}

func isDir(path string) bool {
    info, e := stat(path)
    return e == nil && info.IsDir()
}
```

**Semantics (locked):**
- First `mkdirAll` success → return immediately. No `Stat`, no `sleep` (the invariant all
  four insisted on; hermes's core concern fully honored).
- Success oracle is ALWAYS a fresh `Stat`+`IsDir` — `fs.ErrExist` is never trusted blindly
  (a regular file colliding at the path must still surface an error).
- Fail-fast on `errors.Is(err, fs.ErrPermission)` (after the `isDir` check), so a genuine
  `EACCES`/`EPERM` does not burn the retry window.
- After the bounded retries, return the **last `mkdirAll` error** (not a stat error).
- Worst case 5 `mkdirAll` attempts / 75 ms sleep, only on the error path.

### Call sites (swap `os.MkdirAll` → `fsutil.MkdirAllResilient`)

New-idea launch path (inside `runcontrol.Create`):
- `internal/protocol/workspace.go:67` — `CreateIdea` idea/round-01 dir (runs BEFORE
  `Append`). **Only this one in workspace.go — NOT `InitWorkspace`/line 44** (`parley init`
  bootstrap, not the per-launch hot path).
- `internal/store/events.go:41` — `Store.Append` events dir (proven failure point).
- `internal/runmanifest/manifest.go:142` — manifest dir (runs right after `Append`).

Live run path:
- `internal/runner/runner.go:340` (agentDir) and `:343` (`filepath.Dir(outputPath)`).
- `internal/runner/round_index.go:85` (roundDir).
- `internal/runner/steer.go:152` (steerDir).
- `internal/runner/phase58.go:58` (fix-up agentDir).
- `internal/runner/handoff.go:36` (handoff agentDir).
- `internal/hitl/hitl.go:183` (questions dir).

**Out of scope** (unchanged): `pipeline/*`, `consensus/*`, `sessionstore`, `driver/*`, the
`app.go:1855` probe, `workspace.go:44` (`InitWorkspace`), and the `os.MkdirTemp`
isolated-home dirs (`runner.go:861,901`). They may adopt the helper later if a failure
proves they need it.

### Unchanged
- The trailing `os.OpenFile(events.jsonl, O_CREATE|O_APPEND)` in `Store.Append` (a
  Stat-confirmed dir makes the open succeed; retrying risks duplicate events).
- TUI / launch error handling (no auto-attach, no message change in v1).

## Tests — `internal/fsutil/fsutil_test.go`

Seam-driven, `sleep` stubbed to a counter/no-op so the suite is instant:
1. **common path:** first `mkdirAll` returns nil → `nil`, **`stat` and `sleep` never
   called** (assert via counters).
2. **transient-then-success:** `mkdirAll` errors once then succeeds → `nil`; assert attempt
   count.
3. **host-succeeded-guest-lied:** `mkdirAll` ALWAYS errors but `stat` reports an existing
   dir → `nil` via `isDir` **without exhausting retries / without sleeping**.
4. **already-exists:** real `os.MkdirAll` on a pre-created temp dir → `nil`, zero sleeps.
5. **genuine failure:** `mkdirAll` always returns a non-permission error and `stat` says
   not-exist → returns the **last mkdir error** after the bounded retries (no infinite
   loop); assert the recorded sleeps are `5ms,20ms,50ms`.
6. **fail-fast permission:** `mkdirAll` returns `fs.ErrPermission`, `stat` not-exist →
   returns the error immediately, **zero sleeps**, a single `mkdirAll` attempt.
7. **non-directory collision:** `mkdirAll` errors and `stat` shows a regular file → returns
   an error (not masked as success).

## Verification
`go build ./...`, `go vet ./...`, `go test ./...` all green. Manual smoke: a normal `parley`
round still creates `runs/<id>/events.jsonl` (healthy path unaffected).

## Ship
- Bump `internal/app/version.go` → `1.21.0`.
- IMPLEMENTATION.md, then Phases 6–8 review (codex/agy/hermes review `internal/fsutil` +
  the call-site swaps; fix-up until zero agreed findings).
- GitHub PR (transport `github-pr`), tag `v1.21.0`, GitHub release + Windows x64/arm64
  assets, Homebrew formula bump (sha256 verified).

## Deferred follow-ups (NOT v1)
Auto-attach to an orphaned run on launch error; untruncated launch-error text + a `/open`
hint; adopting the helper in non-launch mkdir sites if they prove flaky.
