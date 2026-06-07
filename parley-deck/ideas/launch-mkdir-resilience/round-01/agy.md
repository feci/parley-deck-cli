---
agent: agy
idea: launch-mkdir-resilience
round: 1
date: 2026-06-07
---

# UX/Failure-Mode Independent Analysis: Launch-Mkdir Resilience

This analysis focuses on the UX impact, orphaned-run recovery, and error messaging surrounding the transient virtio-fs `os.MkdirAll` launch failure.

---

### 1. Helper Design & Semantics

* **Signature & Location:** `func MkdirAllResilient(path string, perm os.FileMode) error` located in a new utility package `internal/fsutil`. This keeps dependencies clean and prevents circular imports across `store`, `runcontrol`, and `runner`.
* **Retry & Backoff:** 
  * Make the first `os.MkdirAll` attempt. If successful, return `nil` immediately (zero overhead/sleep for the happy path).
  * If it fails, perform an immediate `os.Stat(path)`. If the error is nil and `info.IsDir()`, treat it as success.
  * If the directory is still missing, enter a retry loop: up to 4 additional attempts (5 total attempts).
  * Use a short, non-blocking exponential backoff: `5ms`, `15ms`, `45ms`, `100ms`. 
  * Total worst-case latency under transient failure is under 200ms, which is completely unnoticeable to a user launching a run, yet provides ample time for virtio-fs/NFS cache invalidation.
* **Success Conditions:** Success is declared if `os.MkdirAll` succeeds, if `errors.Is(err, fs.ErrExist)` is returned, or if `os.Stat` verifies the path is an existing directory.
* **Fail-Fast for Fatal Errors:** To avoid unnecessary delay, do not retry on fatal errors (e.g., `os.IsPermission(err)` or `errors.Is(err, fs.ErrPermission)`). Return immediately on these.

---

### 2. Scope of Application

* **Recommendation:** Apply `fsutil.MkdirAllResilient` to **all run-path `os.MkdirAll` sites** (specifically `internal/store/events.go:41`, `runner.go`, `round_index.go`, `steer.go`, and manifest directories).
* **Rationale:** Failing at *any* stage of directory setup during run initialization aborts the launch in the same way. Relying on a single chokepoint (`Store.Append`) is insufficient if subsequent steps fail trying to write agent or round outputs. Non-run system setup paths (e.g., config initialization) can continue using vanilla `os.MkdirAll`.

---

### 3. Correctness on virtio-fs/NFS/SMB

* **Stat Cache Synchronization:** On virtualization mounts, calling `os.Stat` on a failed path is highly sound. It forces the virtualization client (e.g., virtio-fs) to revalidate attributes with the host OS, refreshing its local dentry cache.
* **OpenFile Hazards:** A major risk is that the subsequent `os.OpenFile(..., os.O_CREATE|...)` call in `Store.Append` might also fail with `ENOENT` if cache invalidation lags. However, because our helper performs a successful `os.Stat` which forces directory lookup and caches the parent dentry, the subsequent open will succeed. If we want absolute resilience, we must ensure the directory is fully visible before returning.

---

### 4. Orphaned-Run UX

* **The Problem:** Currently, the TUI displays a truncated error message (e.g., `launch failed: mkdir ...: n…`) and aborts, leaving a background run active but invisible.
* **UX Enhancements:**
  1. **Improve Messaging:** Prevent error truncation in the TUI status bar. Display a user-friendly error (e.g., `Launch failed: transient filesystem issue. Check logs.`) and ensure the full, untruncated error stack is written to the parley log.
  2. **Cleanup on Failure:** If launch fails during `runcontrol.Create`, ensure we clean up any orphaned files or background processes spawned during initialization.
  3. **Auto-Attach Recovery:** If a user attempts to launch an idea and a run directory for that timestamp already exists and contains active runs, the TUI should prompt the user: *"An active run was found for this idea. Would you like to attach to it?"* rather than throwing a duplicate path error.

---

### 5. Testability

* **Headless Test Plan:**
  * Introduce a package-level hook in `fsutil`: `var mkdirAll = os.MkdirAll`.
  * Swap `mkdirAll` with mock functions in unit tests:
    * **Test Transient Success:** Return an error on the first call, then succeed on the second. Verify `MkdirAllResilient` returns `nil` and the retry was hit once.
    * **Test Directory Exists (Stat fallback):** Return an error on `os.MkdirAll`, but stub `os.Stat` (or use a real temp dir) to confirm the directory exists. Verify immediate success.
    * **Test Persistent Failure:** Return an error on all calls. Verify the error is eventually returned after the maximum retry count.
    * **Test Fatal Error Fail-Fast:** Return a permission error. Verify it returns immediately without waiting/sleeping.

---

### 6. Safety / No Regressions

* The common path avoids any performance penalty (no sleep or extra system calls on initial success).
* Normal Linux/macOS filesystems operate identically.
* Standard library imports only.

---

### Challenge to Starting Proposal

The starting proposal assumes that `Store.Append` is the only critical point. However:
1. We must wrap *all* run-path directory creations; virtio-fs latency can hit any of them.
2. We must specifically check for permission-denied errors to fail fast. Otherwise, users experience a lagging UI hang on legitimate permission errors.
3. The starting proposal ignores subsequent file operations like `os.OpenFile` failing immediately after directory creation; our helper's `os.Stat` is critical to warming up the dentry cache for those operations.
