---
idea: launch-mkdir-resilience
author: user
created: 2026-06-07
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + diagnosis owner; run-creation/store path; final synthesis
  codex: Go correctness — os.MkdirAll semantics, retry/idempotency design, tests, where the helper lives
  agy: failure-mode/UX — what the user sees on a spurious launch failure; orphaned-run recovery; messaging
  hermes: filesystem/edge-case correctness — virtio-fs/NFS/SMB coherence, ENOENT/ENOTDIR/EEXIST handling, races
transport: local-dir
cross_review_rounds: 1
status: kickoff
---

## Problem (owner's words, verified by facilitator)

In `parley tui`, pressing **N** to start a new idea sometimes fails with:

```
new idea › hello world  launch failed: mkdir /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-test/parley-deck/runs/20260607T191828.148968000Z: n…
```

(the syscall error is truncated in the TUI status row at `n…`).

## Verified diagnosis (design against these facts)

The facilitator investigated the live failure on the owner's machine. Findings:

1. **The run did NOT actually fail.** The directory
   `runs/20260607T191828.148968000Z/` exists and its `events.jsonl` is ~38 KB and
   contains `run.created` plus a full round of events. So `parley` **did** create the
   idea, create the run directory, and run the round in the background. The TUI showed
   "launch failed" and refused to open the live view → the run was **orphaned** (running
   and persisted, but never displayed).

2. **The mount is `AppleVirtIOFS`** (macOS Virtualization.framework virtio-fs shared
   folder), confirmed via `mount`. It has **weak cache coherence**: the facilitator
   directly observed the SAME directory return an **empty** listing from `ls` and a
   **full** listing (3 ideas, 3 runs) from `find` seconds apart. virtio-fs serves stale
   dentry/attribute cache across processes.

3. **The failing call is `os.MkdirAll(s.dir, 0o755)` in `internal/store/events.go:41`**
   (`Store.Append`). The launch path is:
   `tui` → `newLaunchFunc` (`internal/app/app.go:2107`) → `runcontrol.Create`
   (`internal/runcontrol/runcontrol.go:34`) → `runStore.Append(run.created)` →
   `os.MkdirAll(runDir)`. On the stale cache, `MkdirAll`'s internal `Stat`→`Mkdir`→`Lstat`
   dance returns a transient `mkdir …: ENOENT/ENOTDIR` error **even though the directory
   is (or becomes) present on the host**. `runcontrol.Create` propagates the error;
   `newLaunchFunc` returns it; the TUI prints `launch failed: <err>`
   (`internal/tui/live.go:1677`) and aborts `activateRun`.

4. **A synchronous stress test passed clean** (200× mkdir/rmdir, 20× rapid `ls`, 100×
   nested `mkdir -p` — all 0 failures). So the fault is **not** reproducible
   synchronously; it only manifests under **cross-process timing** when parley writes
   while the virtio-fs cache is stale. It is real but intermittent.

## Relevant code (already located — verify, don't re-discover)

- `internal/store/events.go:34-58` — `Store.Append`: `os.MkdirAll(s.dir)` then open+append
  `events.jsonl`. **This is the proven failure point and the gateway for every run** (the
  first thing a new run does).
- `internal/runcontrol/runcontrol.go:34-100` — `Create`: `CreateIdea` → `NewRunID` →
  `runStore.Append(run.created)` → `runmanifest.Write` → `registerSession`. Returns error
  on any failure; the error bubbles to the TUI.
- Other `os.MkdirAll` on the run path that share the same exposure:
  `internal/runner/runner.go:340,343` (agentDir, output dir), `round_index.go:85`
  (roundDir), `steer.go:152` (steerDir), `phase58.go:58`, `handoff.go:36`,
  `runmanifest.Write` (writes the manifest file — check its own dir creation).
- `internal/tui/live.go:1671-1700` — `launchIdea`: on launch error sets
  `m.inputErr = "launch failed: " + err`; on success calls `activateRun`.

## Proposed direction (a STARTING proposal — challenge it in round-01)

Make directory creation on the run path **resilient to spurious/stale-cache failures** on
networked/virtualized filesystems, without masking genuine errors:

1. **A small resilient mkdir helper** (e.g. `fsutil.MkdirAllResilient(path, perm)`): call
   `os.MkdirAll`; on error, **re-`Stat` the path** — if it now exists and is a directory,
   treat as success (this alone fixes the observed case, since the dir does exist). If
   still absent, **retry** a small bounded number of times with a short backoff, then
   return the last error. Treat `EEXIST`/`ErrExist` as success.
2. **Apply it on the run-creation hot path** — at minimum `Store.Append`'s dir creation
   (the proven failure point and the gateway every run hits), and consider the other
   run-path `os.MkdirAll` sites for the same wrapper.
3. **Decide the orphan question:** even with (1), should the TUI/launch be more forgiving
   if a run dir already exists / the run is actually running? Is a recovery path worth it,
   or does fixing mkdir-resilience make orphans impossible in practice?

## Round-01 focus questions (answer independently)

1. **Helper design & semantics.** Signature, where it lives (`internal/fsutil`? inline in
   `store`?), retry count/backoff, which errno/`errors.Is` cases count as "success"
   (`fs.ErrExist`), how to re-check (`os.Stat` + `IsDir`). How to NOT mask a genuine
   persistent failure (bad permissions, real ENOTDIR on a file collision) — bounded
   retries + final real error. Is a tiny sleep acceptable on the launch hot path?
2. **Scope of application.** Just `Store.Append`, or all run-path `os.MkdirAll` sites?
   Trade-off: one chokepoint (simplest, surgical — CLAUDE.md favors this) vs blanket
   resilience. Any non-run `os.MkdirAll` that should NOT change.
3. **Correctness on virtio-fs/NFS/SMB.** Is re-Stat-after-error sound given the SAME stale
   cache might also make `Stat` lie? Does a short backoff actually let the cache settle
   (virtio-fs attr timeout)? Any risk a retry creates a partial/dup? `EEXIST` vs the
   `O_CREATE|O_APPEND` open that follows (does the open also need the same treatment?).
4. **Orphaned-run UX.** What should the user see if a launch errors but the run exists and
   is running? Stay with "prevent it at mkdir" only, or add a recovery/attach path? Keep
   it minimal — the headline fix is mkdir-resilience.
5. **Testability.** How to test headlessly without a real flaky mount: inject a mkdir that
   fails once then succeeds (seam/var), or a temp dir + a wrapper that simulates a
   transient error. Assert: success after one transient failure, success when dir already
   exists, genuine-failure still returns an error after bounded retries.
6. **Safety / no regressions.** Must not slow the normal (non-flaky) path meaningfully (no
   sleep when the first `MkdirAll` succeeds). Must not change behavior on a genuine
   permission/collision error. No new third-party deps. Keep `internal/tui` free of
   runner/app imports (this change is in store/runcontrol/runner/fsutil, not tui).

## Constraints (non-negotiable)
- Surgical (CLAUDE.md): the minimum change that makes launch robust on virtio-fs. Prefer
  one well-placed helper over scattering retries.
- No sleep/backoff cost on the common success path (first `MkdirAll` works → return
  immediately).
- Do NOT mask genuine, persistent errors — bounded retries, then surface the real error.
- No new third-party deps. Standard library only (`os`, `errors`, `io/fs`, `time`).
- Headlessly testable (no dependence on an actual flaky mount).

## Non-goals
- Not "fix virtio-fs". Not a general filesystem abstraction layer. Not retrying every
  file op in the codebase — only the run-creation directory path that the launch failure
  proved fragile. Not changing how agents run.

## Deliverables (in order)
1. FINAL.md: the helper (signature, retry policy, success conditions), exactly which call
   sites adopt it, the orphan decision, and the test plan.
2. Implement in `internal/{fsutil|store|runcontrol|runner}` (+ tests); build/vet/test
   green; bump to 1.21.0. Then Phases 6–8 review. Ship via GitHub PR (deck transport),
   tag v1.21.0, GitHub release, Homebrew bump.
