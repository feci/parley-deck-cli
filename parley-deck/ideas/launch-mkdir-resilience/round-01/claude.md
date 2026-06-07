---
agent: claude
idea: launch-mkdir-resilience
round: 1
date: 2026-06-07
---

## Summary

I diagnosed the live failure, so I'll state the design decisively and let the others
challenge it. The bug is **not** a logic error in `parley` — it is a spurious
`os.MkdirAll` error returned by the `AppleVirtIOFS` shared mount under a stale
dentry/attribute cache, surfaced as a hard launch failure even though the run directory
is (or becomes) present and the run actually runs. The fix is a **drop-in resilient
mkdir** that re-checks existence after an error and retries a bounded number of times.
This is small, surgical, and has zero cost on the healthy path.

## Q1 — Helper design & semantics

New package `internal/fsutil` with one function (stdlib only):

```go
// MkdirAll behaves like os.MkdirAll but tolerates spurious/transient failures from
// networked or virtualized filesystems (virtio-fs, NFS, SMB) whose dentry/attribute
// cache can be momentarily stale across processes. It returns nil whenever the target
// exists as a directory — even if an intermediate os.MkdirAll call returned an error.
func MkdirAll(path string, perm fs.FileMode) error {
    if err := osMkdirAll(path, perm); err == nil || isDir(path) {
        return nil
    }
    var err error
    for i := 0; i < maxRetries; i++ {          // maxRetries = 3
        sleep(backoff(i))                       // 20ms, 40ms, 80ms — only on the error path
        if err = osMkdirAll(path, perm); err == nil || isDir(path) {
            return nil
        }
    }
    return err
}

func isDir(path string) bool { fi, e := osStat(path); return e == nil && fi.IsDir() }
```

- **Common path cost = zero.** If the first `os.MkdirAll` succeeds we return immediately;
  no `Stat`, no sleep. Backoff only ever runs after a real error.
- **Success conditions:** (a) `os.MkdirAll` returns nil, or (b) a fresh `os.Stat` shows a
  directory. `os.MkdirAll` already maps a pre-existing dir and inner `EEXIST` (via its own
  `Lstat`) to nil; our extra `isDir` re-check covers the case where that internal `Lstat`
  itself was served a stale error. We do **not** special-case `fs.ErrExist` separately —
  `isDir` subsumes it and is simpler.
- **Not masking genuine failures:** bounded to 3 retries / ~140 ms worst case. If the dir
  genuinely cannot be created (bad perms, a *file* sits at a parent component → real
  `ENOTDIR`), every `os.MkdirAll` keeps failing and `isDir` stays false, so we return the
  real, last error. We never loop forever and never invent success.
- A tiny sleep on the launch hot path is acceptable: it happens **only** when the first
  attempt already failed (i.e. we were going to error out anyway); trading ~20–140 ms for
  a successful launch is strictly better than the current hard failure.

## Q2 — Scope of application

Surgical, but consistent. The **proven** failure point is `Store.Append`'s
`os.MkdirAll(s.dir)` (`internal/store/events.go:41`) — and since every run begins with
`Append(run.created)`, fixing that one site fixes the reported launch failure. But the
identical exposure exists at every other **run-path** mkdir reached mid-run; a stale-cache
hiccup there would abort a round just as spuriously. Because the change is a one-token swap
(`os.MkdirAll` → `fsutil.MkdirAll`), I propose adopting it at the run-creation/run-path
sites as a set:

- **Must:** `internal/store/events.go:41` (Append).
- **Should (same run path, trivial swap):** `internal/runner/runner.go:340,343`,
  `internal/runner/round_index.go:85`, `internal/runner/steer.go:152`,
  `internal/runner/phase58.go:58`, `internal/runner/handoff.go:36`, and `runmanifest`'s
  own dir creation if it has one.
- **Leave alone:** non-run mkdirs — `sessionstore`, `pipeline/*`, and the `os.MkdirTemp`
  isolated-home dirs (these are on the local/native fs or are temp dirs; out of scope,
  keep the diff small per CLAUDE.md). I'm open to codex arguing for fewer sites (only
  Append) on strict-surgical grounds.

**The trailing open matters too.** `Append` does `os.OpenFile(events.jsonl,
O_CREATE|O_APPEND)` right after the mkdir. After a `Stat`-confirmed `fsutil.MkdirAll` the
dir is visible to *this* process, so the open should succeed; but to make the *launch* (not
just the mkdir) robust, I lean toward a 1-retry wrapper on that open as well, or making
`Append` retry the mkdir+open pair once on a stale error. I'll defer the exact shape to
round-02 — minimum viable is `fsutil.MkdirAll` for the dir; the open-retry is cheap
insurance worth discussing.

## Q3 — Correctness on virtio-fs/NFS/SMB

Re-`Stat`-after-error is sound precisely because the failure mode is **cache staleness,
not corruption**: the host filesystem is consistent; only the guest's cached view lags.
Three cases:
1. Host created the dir, guest cache lied on `MkdirAll`'s internal stat → our fresh `Stat`
   (or the next one after backoff, once the attr cache TTL elapses) sees it → success.
2. Dir not yet there due to a transient host hiccup → a retry creates it → success.
3. Dir genuinely uncreatable → all retries fail, `isDir` false → real error returned.

The short backoff (20→80 ms) is chosen to exceed a typical virtio-fs attr-cache settle
window without being user-perceptible. Retrying `os.MkdirAll` is idempotent (it no-ops on
an existing dir), so retries can't create partials or duplicates. **hermes** should
pressure-test whether the same stale cache could make even the *final* `Stat` lie long
enough to exhaust 3 retries — if so we bump `maxRetries`/backoff, but I believe 3 is ample
given the TTLs involved.

## Q4 — Orphaned-run UX

Fixing mkdir makes the orphan (run created + running but "launch failed" shown) **not
happen** in the observed case, which is the right primary fix. A secondary
attach-on-error recovery path (TUI offers to open the existing run dir when launch errors
but the run exists) is more invasive — the live handle/seams aren't built on the error
path — and I judge it **out of scope** for v1. Recommend: prevent at mkdir only; note
attach-on-error as a deferred follow-up. (agy owns this lens — if there's a cheap message
improvement, e.g. surfacing the full untruncated error, raise it.)

## Q5 — Testability

Inject a seam in `fsutil`: package vars `osMkdirAll = os.MkdirAll`, `osStat = os.Stat`,
`sleep = time.Sleep` (or a `backoff` override), overridden in tests. Cases:
- **transient-then-success:** stub fails once, second call succeeds → `MkdirAll` returns
  nil and the dir exists; assert exactly the expected number of attempts.
- **host-succeeded-guest-lied:** stub *always* returns an error but the dir is pre-created
  on disk → returns nil via `isDir` short-circuit **without** exhausting retries.
- **already-exists:** real `os.MkdirAll` on an existing temp dir → nil, zero sleeps.
- **genuine failure:** stub always errors and dir never exists → returns the error after
  exactly `maxRetries` retries (assert no infinite loop, real error preserved).
- **fast common path:** first attempt succeeds → no `Stat`, no sleep (assert via counters).
Keep `sleep` a no-op in tests so the suite stays instant.

## Q6 — Safety / no regressions

- Zero added latency when the fs is healthy (the overwhelming common case).
- Behavior on genuine permission/collision errors is unchanged (real error still returned).
- Stdlib only (`os`, `io/fs`, `time`); no new deps; `internal/tui` untouched (the change
  lives in `fsutil`/`store`/`runner`), so the no-runner/app-import rule for tui holds.
- `--no-tui`, the round/steer path, and durable-kill/liveness are unaffected (they just
  get a more reliable mkdir).

## Recommendation

Ship `internal/fsutil.MkdirAll` (3 retries, 20/40/80 ms backoff, `isDir` re-check),
swap it in at `Store.Append` + the run-path mkdir sites, leave non-run mkdirs alone, add
the headless table-test with injected transient failures, bump to 1.21.0. Open questions
for round-02: (a) exact open-retry treatment in `Append`; (b) strict-minimal (Append only)
vs run-path-wide scope; (c) whether `maxRetries`/backoff need tuning for virtio-fs TTLs.
