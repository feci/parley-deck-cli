---
agent: hermes
idea: launch-mkdir-resilience
round: 1
date: 2026-06-07
---

## Summary
Independent analysis under filesystem/edge-case lens. The root cause is virtio-fs stale dentry cache producing transient ENOENT/ENOTDIR on the first os.MkdirAll in Store.Append (events.go:41) even when the directory exists on the host. The same cache can affect the subsequent O_CREATE|O_APPEND open and other run-path MkdirAll sites. A minimal, bounded retry helper that only re-Stats on error (no sleep on first success) is sound and sufficient.

## Proposed approach
Introduce `internal/fsutil.MkdirAll(path string, perm os.FileMode) error` that:
- Calls os.MkdirAll once.
- On any error, does os.Stat + IsDir; if true, return nil (success).
- Otherwise retry the full MkdirAll at most 3 times total (initial + 2 retries) with no sleep.
- Treat fs.ErrExist and os.IsExist as immediate success.
- Return the last observed error only after retries exhausted; never wrap or mask.

Place the helper in a new internal/fsutil package (single file, stdlib only). Adopt it at exactly these call sites: Store.Append (events.go:41), runner.go:340+343, round_index.go:85, steer.go:152, phase58.go:58, handoff.go:36. Do not touch non-run MkdirAll sites (consensus, sessionstore, pipeline_cmd, tests). This is one surgical chokepoint change honoring CLAUDE.md.

The trailing OpenFile in events.go:44 does not need identical treatment; its ENOENT would be a genuine persistent failure after the resilient MkdirAll succeeded.

Re-Stat-after-error is sound on virtio-fs: the cache invalidation window is sub-second under cross-process write; a second MkdirAll attempt forces the host to observe the directory. Same-cache Stat lie is mitigated by the bounded retries and the fact that a real persistent ENOTDIR (file-in-path) or EACCES will repeat identically. No risk of partial/dup directories.

## Focus question answers
1. Helper: `fsutil.MkdirAll(path string, perm os.FileMode) error` lives in internal/fsutil. Retry count=3 total attempts, zero backoff/sleep (hot-path invariant). Success cases: first MkdirAll nil, or post-error os.Stat succeeds + IsDir, or fs.ErrExist. Genuine persistent failures (EACCES, real ENOTDIR on a file, permission denial) surface after retries. No sleep ever on success path.

2. Scope: exactly the six run-creation sites listed above. Trade-off: single helper + six call-site swaps is still minimal and surgical vs. blanket replacement. Non-run MkdirAll (consensus writes, sessionstore, tests, pipeline) remain untouched to avoid scope creep.

3. Correctness: re-Stat is sound because the failing MkdirAll itself is the coherence trigger; subsequent calls see fresh state. Short backoff unnecessary and forbidden on hot path. O_CREATE open inherits safety once its parent dir is resiliently ensured. No duplicate-run or partial-dir races introduced.

4. Orphan UX: prevent at mkdir only. No attach/recovery path. The mkdir fix eliminates the observed failure mode; adding recovery would violate "minimum change" and "Simplicity First".

5. Headless test plan: add fsutil/mkdir_test.go with a seam (var failNextMkdir atomic.Bool). Test cases: (a) transient failure on first attempt then success asserts nil error; (b) directory already exists asserts nil; (c) persistent EACCES after 3 attempts returns the real error. Use os.MkdirTemp + manual removal between attempts; no real virtio-fs required.

6. Safety: zero cost on success (single MkdirAll + immediate return). Genuine errors unchanged after bounded attempts. No new deps. tui/live.go untouched. All run-path sites now consistent.

## Concerns / risks
None under the assigned lens; the design is the narrowest fix that closes the exact failure observed on AppleVirtIOFS. Counter-proposal to starting sketch: drop any sleep/backoff language entirely (already ruled out by constraints) and make the helper package-level rather than in store to keep surgical scope.