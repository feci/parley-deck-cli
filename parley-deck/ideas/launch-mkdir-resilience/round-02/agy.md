---
agent: agy
idea: launch-mkdir-resilience
round: 2
date: 2026-06-07
---

## Convergence on Core Design
All agents agree on implementing a centralized helper `internal/fsutil.MkdirAllResilient(path, perm)`. The helper returns immediately on initial `os.MkdirAll` success (zero cost). On error, it performs an `os.Stat` + `IsDir` check and retries a bounded number of times before propagating the authentic final mkdir error. It is written using only standard library packages and is headlessly testable via package-level seams.

---

## Fork Resolutions & Feedback

### 1. BACKOFF (Reject hermes's zero-sleep; Compromise Accepted)
* **Stance:** I disagree with **hermes**'s zero-sleep proposal. Without any sleep, back-to-back retries run in the same guest CPU timeslice before the virtualization layer (e.g., `virtiofsd`) can refresh the stale attribute/dentry cache, rendering retries useless. I agree with **claude** and **codex** that backoff is necessary on the error path.
* **Resolution:** Accept the compromise. The first retry runs immediately (0ms) to resolve temporary CPU scheduling races as **hermes** suggests, but subsequent retries back off.
* **Exact Schedule:** 5 total attempts (1 initial + 4 retries).
  * Attempt 1: Immediate `os.MkdirAll`.
  * Attempt 2: Immediate retry (0ms sleep).
  * Attempt 3: 20ms sleep.
  * Attempt 4: 50ms sleep.
  * Attempt 5: 100ms sleep.
  * *Worst-case latency under persistent failure:* 170ms (imperceptible to TUI users, yet gives the mount cache sufficient time to sync).

### 2. SCOPE (Disagree with codex/hermes; Include CreateIdea)
* **Stance:** I disagree with **codex** and **hermes** on scoping out `protocol.CreateIdea` (`internal/protocol/workspace.go:67`).
* **Resolution:** Include the full launch path. From a UX perspective, a transient mkdir error in `CreateIdea` results in the exact same truncated "launch failed" status bar error as a failure in `Store.Append`. Since both execute within the launch transaction (`runcontrol.Create`) on the same weak-coherence mount, we should protect the full launch path: `CreateIdea`, `Store.Append`, and `runmanifest.Write`. Runner sites and `hitl.go` also adopt the helper. Non-run/setup paths are excluded.

### 3. EEXIST (Agree with codex/hermes)
* **Stance:** I agree with **codex** (and **hermes**'s round-02 alignment).
* **Resolution:** We must never treat `fs.ErrExist` (or `os.IsExist`) as immediate success blindly. A file collision at the target path would produce `fs.ErrExist` but subsequently fail write operations. We must verify via `os.Stat` + `IsDir` that the existing entity is actually a directory before returning success.

### 4. FAIL-FAST (Confirm agy; Disagree with hermes)
* **Stance:** I disagree with **hermes**'s objection to fail-fast.
* **Resolution:** We should fail fast on permission errors (`os.IsPermission(err)` or `errors.Is(err, fs.ErrPermission)`). Returning immediately does not "hide" the real errno (it propagates it immediately). Delaying a genuine permission error by looping through 170ms of useless retries degrades the TUI response time for no benefit.

### 5. OPEN FILE (Confirm codex/hermes/agy agreement)
* **Stance:** Confirm agreement.
* **Resolution:** Keep the trailing `os.OpenFile` in `Store.Append` unchanged. The `os.Stat` call in `MkdirAllResilient` acts as a cache synchronization barrier, guaranteeing parent directory visibility for the subsequent `OpenFile` without exposing the write to duplication risks.

### 6. TEST PLAN (Confirm consensus)
* **Stance:** Confirm the unit test suite plan using package-level seams (`var mkdirAll = os.MkdirAll`, etc.) and a stubbed no-op `sleep`.
* **Resolution:** Assert the following cases in `internal/fsutil/mkdir_test.go`:
  1. *Transient Success:* Failure on attempt 1, success on attempt 2 returns nil.
  2. *Dir Exists via Stat:* `mkdirAll` returns error but `stat` confirms `IsDir`; returns nil without sleeping.
  3. *Already Exists:* Normal pre-existing directory returns nil with zero sleeps.
  4. *Genuine Failure:* Persistent errors return the last authentic mkdir error after exhausting retries.
  5. *Fail-Fast Permission:* `fs.ErrPermission` returns immediately on attempt 1 with zero sleeps.
  6. *Common Path Fast:* First attempt success returns nil; asserts zero `stat` or `sleep` calls.
