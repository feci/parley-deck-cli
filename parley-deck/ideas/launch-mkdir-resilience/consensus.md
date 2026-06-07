---
idea: launch-mkdir-resilience
phase: consensus
date: 2026-06-07
drafter: claude
participants: [claude, codex, agy, hermes]
---

## Consensus

After two design rounds the participants converge on a small, stdlib-only resilient mkdir
helper applied across the new-idea launch path and the live run path. Forks resolved by
3/4 majority where not unanimous (hermes's recorded dissents noted).

### Decided

1. **Helper.** New package `internal/fsutil`, one exported function:
   `func MkdirAllResilient(path string, perm os.FileMode) error`. Stdlib only
   (`os`, `io/fs`, `errors`, `time`). Package-var seams `mkdirAll = os.MkdirAll`,
   `stat = os.Stat`, `sleep = time.Sleep` for headless tests.

2. **Algorithm.**
   - `mkdirAll(path, perm)`; on `nil` → return `nil` immediately. **No Stat, no sleep on
     the healthy path** (the universal invariant everyone insisted on).
   - On error, loop the retry schedule. Before each retry: (a) if `isDir(path)` (fresh
     `os.Stat` + `IsDir`) → return `nil`; (b) else if `errors.Is(err, fs.ErrPermission)`
     → return `err` (fail-fast — no point burning the window on a real `EACCES`/`EPERM`);
     (c) else sleep the delay (if > 0) and `mkdirAll` again.
   - **Retry schedule:** `[0, 5ms, 20ms, 50ms]` → **5 total `MkdirAll` attempts** (initial
     + immediate retry + three slept retries), worst-case added latency **75 ms, only on
     the error path**.
   - After the last attempt, a final `isDir` check; else return the **last `MkdirAll`
     error** (not a stat error).
   - **Success oracle is always "a fresh Stat shows a directory"** — `fs.ErrExist` is never
     trusted blindly (a file colliding at the path must still surface an error).

3. **Scope — the new-idea launch path + the live run path** (swap `os.MkdirAll` →
   `fsutil.MkdirAllResilient`):
   - `internal/protocol/workspace.go` — the idea-dir mkdir(s) on the `protocol.CreateIdea`
     path (runs inside `runcontrol.Create` BEFORE `Append`; same launch, same exposure).
   - `internal/store/events.go:41` — `Store.Append` events dir (the proven failure point).
   - `internal/runmanifest/manifest.go:142` — manifest dir (runs right after `Append`).
   - `internal/runner/runner.go:340,343`, `round_index.go:85`, `steer.go:152`,
     `phase58.go:58`, `handoff.go:36` — live run path.
   - `internal/hitl/hitl.go:183` — questions dir.
   - **Out of scope:** `pipeline/*`, `consensus/*`, `sessionstore`, `driver/*`, the
     `app.go:1855` probe, and the `os.MkdirTemp` isolated-home dirs. They may adopt the
     helper later if a failure proves they need it.

4. **Trailing `os.OpenFile` in `Store.Append` is unchanged** (unanimous). A Stat-confirmed
   directory makes the `O_CREATE|O_APPEND` open succeed; retrying the open/write risks
   duplicate events.

5. **Orphan UX = prevention at mkdir only** for v1. Auto-attach recovery, untruncated
   error text, and a `/open` hint in the launch-error message are **deferred follow-ups**.

6. **Tests** (headless, seams, `sleep` stubbed to no-op), in `internal/fsutil`:
   common-path (first success → zero `stat`, zero `sleep`); transient-then-success;
   host-succeeded-guest-lied (mkdir always errors but dir pre-exists → `nil` via `isDir`
   **without exhausting retries**); already-exists; genuine-failure (returns the last
   mkdir error after the bounded retries); fail-fast-permission (returns immediately, zero
   sleeps); non-directory collision (mkdir errors + `stat` shows a regular file → error,
   not masked).

7. **Ship:** bump `internal/app/version.go` → `1.21.0`; Phases 6–8 review; GitHub PR
   (transport `github-pr`), tag `v1.21.0`, GitHub release + Windows assets, Homebrew bump.

### Recorded dissents (overridden by majority, not blockers)
- **hermes** preferred (a) zero backoff / 3 attempts, (b) excluding `protocol.CreateIdea`
  (run-path only), and (c) no fail-fast on permission (exhaust the window for an authentic
  final error). Majority chose the small-backoff compromise, full-launch-path scope, and
  permission fail-fast. hermes's core invariant — **zero cost on the success path** — is
  fully preserved.

## Signoffs

Each participant appends its own ACCEPT/REJECT block below (append-only).

### claude — ACCEPT
Facilitator + diagnosis owner. I drafted this consensus from the two design rounds; it
matches the verified failure (spurious `os.MkdirAll` on `AppleVirtIOFS` stale cache) and
the minimal fix. The healthy path stays a single `os.MkdirAll` with zero overhead. Forks
resolved by 3/4 majority; hermes's dissents recorded. Ready to implement as the FINAL
drafter/implementer.

### codex — ACCEPT
I accept the decided `internal/fsutil.MkdirAllResilient` helper, retry algorithm, launch
and live-run scope, and the focused seam-based test set. The consensus matches my
round-02 position: zero cost on the healthy path, fresh Stat as the success oracle, and
bounded retries only on mkdir errors. My only reservation is that auto-attach/orphan UX
remains deferred, but that is correctly outside the minimal v1 fix.

### hermes — ACCEPT
I accept the consensus: the resilient mkdir helper with the decided retry schedule, full launch+run scope, and permission fail-fast. My round-02 preferences for zero backoff, narrower scope, and no fail-fast are noted but not blockers. The zero-cost success path invariant is preserved exactly as required.

### agy — ACCEPT
I accept the consensus: the fsutil.MkdirAllResilient helper, the retry/backoff algorithm with fail-fast on permission (which I proposed), the full launch+run-path scope, and the seam-based tests. No blockers.
