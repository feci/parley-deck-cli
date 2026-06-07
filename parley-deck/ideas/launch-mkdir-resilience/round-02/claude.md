---
agent: claude
idea: launch-mkdir-resilience
round: 2
date: 2026-06-07
---

## Cross-review

All four round-01 analyses converge on the same shape: a stdlib-only
`internal/fsutil.MkdirAllResilient(path, perm)` that runs `os.MkdirAll` once, and on
error re-`Stat`s + retries with the directory's existence as the success oracle — zero
cost on the healthy path. I'm locking the design and resolving the open forks below.

### @codex
Adopting your design as the spine, in full:
- **Helper in `internal/fsutil`** (not inline in `store`) so `store`/`runmanifest`/`runner`
  share one tested implementation. **Adopt.** Name **`MkdirAllResilient`** (you, agy,
  hermes all chose it — clearer than my `MkdirAll` at the call sites). **Adopt.**
- **Do NOT blind-trust `fs.ErrExist`** — confirm success only via `os.Stat`+`IsDir`,
  because `EEXIST` can mean a *file* collides at the target path. **Adopt — overrides
  hermes/agy's "treat ErrExist as immediate success".** (My `isDir` re-check already
  subsumes this; we won't special-case `ErrExist`.)
- **Return the last `os.MkdirAll` error** after exhausting retries (not a stat error) —
  callers expect the mkdir error pointing at the attempted op. **Adopt.**
- **Open file stays unchanged** — a `Stat`-confirmed dir makes the trailing
  `O_CREATE|O_APPEND` open succeed, and retrying risks event duplication. **Adopt.**
- **Scope = the run tree, not every repo mkdir.** **Adopt**, and I extend it by one
  category (see SCOPE below).

### @hermes
- **The "retry itself revalidates the cache"** insight is right and I'm honoring it in the
  backoff compromise (below): the **first retry is immediate (0 ms)**. But I'm **not**
  taking the strict zero-sleep stance for *all* retries — a sub-100 ms settle window on the
  already-failing error path is free insurance against a still-stale attr cache, and three
  of four of us want it. So: immediate first retry, then a tiny backoff. **Partial adopt.**
- **Bounded, idempotent, no partial/dup risk; never wrap/mask the error.** **Adopt.**
- Your minimal site list (Append + 5 runner sites) is the floor; I'm adding `runmanifest`
  and the idea-dir (SCOPE) because they sit on the same launch path.

### @agy
- **Fail-fast on `fs.ErrPermission`** — don't spend the retry/backoff window on a genuine
  `EACCES`. **Adopt** (cheap, correct: permission errors are not cache artifacts).
- **Untruncated error messaging / orphan recovery / auto-attach** — good UX lens, but the
  mkdir fix *removes the cause*, so I'm keeping v1 to prevention and deferring the UX
  additions (auto-attach is the invasive one codex/hermes/I all flagged as out of scope).
  One cheap item I'll keep as a **deferred follow-up, not v1**: have `launchIdea`'s error
  hint at `/open`. **Defer.**
- Your 5/15/45/100 ms schedule is a bit long; I'm taking a shorter one (below).

## Resolved decisions (for FINAL)

1. **Helper:** `internal/fsutil.MkdirAllResilient(path string, perm os.FileMode) error`,
   new package, stdlib only (`os`, `io/fs`, `errors`, `time`). Package-var seams
   `mkdirAll = os.MkdirAll`, `stat = os.Stat`, `sleep = time.Sleep` for tests.
2. **Algorithm:**
   - `os.MkdirAll(path, perm)`; if nil → return nil. **No `Stat`, no sleep on success.**
   - On error: if `errors.Is(err, fs.ErrPermission)` → return err immediately (fail-fast).
   - Else `isDir(path)` (fresh `os.Stat` + `IsDir`) → if true return nil.
   - Else retry: **4 total attempts** (1 initial + 3 retries), retry delays **`[0, 25ms,
     75ms]`** (first retry immediate per hermes; small settle after). Before each sleep and
     after each failed `MkdirAll`, re-check `isDir`; on true → nil.
   - After the last attempt, final `isDir`; if true → nil; else return the **last
     `os.MkdirAll` error**. Worst-case added latency ≈100 ms, only on the error path.
   - Success oracle is ALWAYS "dir exists per a fresh Stat" — never blind `ErrExist`.
3. **Scope = the new-idea launch path + the live run path.** Swap `os.MkdirAll` →
   `fsutil.MkdirAllResilient` at:
   - `internal/store/events.go:41` — `Store.Append` events dir (the proven failure point).
   - `internal/runmanifest/manifest.go:142` — manifest dir (runs in `runcontrol.Create`
     right after `Append`; if it fails the launch still fails).
   - **`internal/protocol/workspace.go`** — the idea-dir mkdir(s) reached by
     `protocol.CreateIdea` (runs in `runcontrol.Create` BEFORE `Append`). **I'm including
     this over codex/hermes's "leave protocol alone"**, because it is literally on the
     launch the user triggered and has identical virtio-fs exposure; excluding it leaves a
     known gap. (Limit to the CreateIdea path, not unrelated `parley init` workspace
     bootstrap, to keep the diff tight — confirm exact lines at implementation.)
   - `internal/runner/runner.go:340,343`, `round_index.go:85`, `steer.go:152`,
     `phase58.go:58`, `handoff.go:36`, `internal/hitl/hitl.go:183` — the live run path.
   - **Leave alone:** `pipeline/*`, `consensus/*`, `sessionstore`, `driver/*`,
     `app.go:1855` probe, and the `os.MkdirTemp` isolated-home dirs (native /tmp). They can
     adopt the helper later if a failure proves they need it (`log` it if so).
4. **Open file unchanged** (no write-retry → no event-dup risk).
5. **Orphan UX:** prevented at mkdir; auto-attach / untruncated-error / "/open" hint are
   deferred follow-ups, not v1.
6. **Tests** (headless, seam-driven, `sleep` stubbed to no-op):
   - common path: first `MkdirAll` nil → returns nil, **`stat` and `sleep` never called**.
   - transient-then-success: fails once then succeeds → nil; assert attempt count.
   - host-succeeded-guest-lied: `mkdirAll` *always* errors but dir pre-exists on disk →
     nil via `isDir` **without exhausting retries**.
   - already-exists: real `os.MkdirAll` on an existing temp dir → nil, zero sleeps.
   - genuine failure: `mkdirAll` always errors, dir never exists → returns the **last
     mkdir error** after exactly 3 retries (no infinite loop).
   - fail-fast permission: `mkdirAll` returns `fs.ErrPermission` → returns immediately,
     **no sleeps, no further attempts**.
   - non-dir collision: `mkdirAll` errors, `stat` shows a regular file → returns error
     (not masked).
7. **Ship:** bump `internal/app/version.go` → `1.21.0`; Phases 6–8 review; GitHub PR
   (deck transport `github-pr`), tag `v1.21.0`, GitHub release + Windows assets, Homebrew
   bump.

## Risks carried to FINAL
The same stale cache could in principle make the *final* `Stat` lie long enough to exhaust
4 attempts (~100 ms). If that ever reproduces we raise the attempt count/backoff — but the
observed virtio-fs attr-cache window is sub-second and the act of `MkdirAll`+`Stat` forces
revalidation, so 4 attempts is ample. Fail-fast on permission must use `errors.Is(err,
fs.ErrPermission)` (covers `EACCES`/`EPERM` wrapped in `*PathError`). The `protocol`
inclusion must touch only the CreateIdea path, not unrelated workspace bootstrap.
