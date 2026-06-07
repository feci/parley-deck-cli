---
agent: codex
idea: launch-mkdir-resilience
round: 1
date: 2026-06-07
---

## Summary

The fix should be a small standard-library helper in a shared internal package, adopted on the run path rather than only in `Store.Append`. The common path must remain one plain `os.MkdirAll` call with no sleep; only an error path may pay for re-stat and bounded retry. I would not add a TUI attach/recovery feature in this change; make launch creation robust at the filesystem boundary and keep `/open` as the manual recovery path.

## 1. Helper design and semantics

Put the helper in `internal/fsutil`, not inline in `store`. `Store.Append` is the proven failure point, but `runmanifest.Write`, runner artifact directories, steers, fixups, and handoff packets all create directories under the same `parley-deck/runs/<runID>` tree on the same weakly coherent mount. A tiny package avoids duplicating retry logic and avoids making `store` own runner/runmanifest concerns.

Exact API:

```go
package fsutil

func MkdirAllResilient(path string, perm fs.FileMode) error
```

Implementation policy:

- First call `os.MkdirAll(path, perm)` and return immediately on success. No stat, no sleep.
- On error, call `os.Stat(path)`. If stat succeeds and `info.IsDir()` is true, return nil.
- If stat succeeds and `!info.IsDir()`, return the original mkdir error, not nil.
- Retry only when the path was not confirmed to be a directory. Use three sleeps after failed attempts: `5ms`, `20ms`, `50ms`, for four total mkdir attempts including the first.
- After every failed retry, re-stat before sleeping/continuing. If the directory appears, return nil.
- After the final attempt, do one final stat; if it is a directory, return nil; otherwise return the last `MkdirAll` error.

`fs.ErrExist`/`EEXIST` should not be accepted blindly. It counts as success only through the same `os.Stat(path) && IsDir()` confirmation, because `EEXIST` can also mean a real file collision at the target path. `fs.ErrPermission`, real `ENOTDIR`, and persistent `ENOENT` are not special-cased into success; they only get the bounded retry window and then surface as the last real mkdir error.

A tiny sleep is acceptable on the launch hot path only after an actual mkdir failure. The observed failure is already exceptional and user-visible; adding at most 75ms in that path is a better trade-off than orphaning a run. There must be zero sleep cost when the first `MkdirAll` succeeds.

## 2. Scope of application

Adopt `fsutil.MkdirAllResilient` at run-path directory creation sites:

- `internal/store/events.go`: `Store.Append` for the run events directory. This is mandatory and is the proven failure point.
- `internal/runmanifest/manifest.go`: manifest directory creation. It immediately follows `Store.Append` in `runcontrol.Create` and hits the same `runs/<runID>` directory.
- `internal/runner/runner.go`: `agentDir` and `filepath.Dir(outputPath)`.
- `internal/runner/round_index.go`: `roundDir`.
- `internal/runner/steer.go`: `steerDir`.
- `internal/runner/phase58.go`: fixup `agentDir`.
- `internal/runner/handoff.go`: handoff `agentDir`.
- `internal/hitl/hitl.go`: questions directory, because launch-time auto-answer and interactive questions are part of the run tree.

Do not sweep all repo-wide `os.MkdirAll` calls. Leave workspace bootstrap, protocol idea creation, pipeline, driver, consensus, app probe, tests, and provider home/cache setup alone unless a later failure proves they need the same treatment. This honors "Simplicity First" while avoiding a half-fix where launch succeeds but the same run immediately fails creating its manifest or first agent directory.

If the team wants the most surgical possible patch, `Store.Append` plus `runmanifest.Write` is the minimum viable subset. My recommendation is the run-path list above because the helper is already shared and the additional call-site changes are mechanical, small, and directly tied to the same weakly coherent `runs/` tree.

## 3. Correctness on virtio-fs, NFS, and SMB

Re-stat-after-error is not a perfect coherence proof. On weak caches, `Stat` can also lie. It is still sound as a success condition because a positive `Stat` saying "this path is a directory" is exactly the state `MkdirAll` was trying to establish; the subsequent file operation remains the real end-to-end check. A false negative stat does not mask failure; it just drives the bounded retry.

The retry is useful because these filesystems often have short dentry/attribute cache windows and cross-process visibility lag. A 5/20/50ms backoff gives the cache a chance to settle without turning launch into a long hang. Retries are idempotent for directory creation: duplicate `MkdirAll` attempts either observe the directory or continue failing; they do not create duplicate run IDs or duplicate events.

The trailing `os.OpenFile(events.jsonl, O_CREATE|O_APPEND|O_WRONLY, 0644)` should not get the same generic retry in this change. If resilient mkdir returns success because stat confirms the directory, the open should normally work. If the same cache makes open fail transiently, `Store.Append` will still return an error, but that is a distinct unproven file-open resilience problem with more risk: retrying append can duplicate events if the first write partially succeeded after open. For now, keep append semantics unchanged. If logs later prove open itself can fail before any write, handle that with a narrower "open file after confirmed directory" retry, not a broad write retry.

## 4. Orphaned-run UX

Keep this change to prevention at mkdir. Do not add an automatic attach/recovery path in `launchIdea`.

The TUI already has `/open <slug|run>` for manual recovery, and adding automatic attach would require reliable detection that `runcontrol.Create` actually created and registered a specific run after returning an error. That is broader than the diagnosed bug and risks coupling `internal/tui` to run internals. After the mkdir resilience patch, the original orphan mode should become very unlikely because `Create` should not report failure when the run directory exists.

Minimal UX follow-up, if desired later: make the error text mention `/open` when launch fails. I would defer even that unless users still hit orphan cases after the filesystem fix.

## 5. Testability

Add focused tests for `internal/fsutil.MkdirAllResilient` using package-level seams inside `internal/fsutil`:

```go
var mkdirAll = os.MkdirAll
var stat = os.Stat
var sleep = time.Sleep
```

Tests should reset seams with `t.Cleanup`. Stub `sleep` to record durations and avoid real delays.

Required headless tests:

- Success after one transient mkdir failure: first `mkdirAll` returns a synthetic error after creating the directory or with `stat` reporting an existing directory; helper returns nil, and no extra failure leaks.
- Success when target already exists as a directory and `mkdirAll` returns an `fs.ErrExist`-style error: helper returns nil only because `stat` confirms `IsDir`.
- Genuine failure still errors: `mkdirAll` always returns a permission or ENOTDIR-like error and `stat` returns not-exist or a non-directory; helper retries the bounded number of attempts, sleeps exactly `5ms`, `20ms`, `50ms`, then returns the last mkdir error.
- Common path has no overhead: first `mkdirAll` returns nil; assert `stat` and `sleep` were not called.
- Existing non-directory is not masked: `mkdirAll` returns `fs.ErrExist`, `stat` returns a regular file, helper returns an error.

I would also add one store-level test only if current tests do not already cover `Store.Append`: inject the helper seam is harder if the helper is imported directly, so the main value is in fsutil unit coverage. The call-site changes are simple enough that `go test ./...` plus fsutil tests should be sufficient.

## 6. Safety and counter-proposals

I agree with the starting proposal's shape but would tighten two points:

- Do not treat `fs.ErrExist` as success by itself. Confirm the path is a directory. This preserves correct behavior when a file collides with the intended run directory.
- Do not apply the helper to every `os.MkdirAll` in the codebase. Apply it to the run tree. That is broader than only `Store.Append`, but still surgical and directly linked to launch/run execution.

The helper should return the last mkdir error after exhausting retries, not a stat error, unless implementation finds a clearer way to preserve both. The mkdir error is what callers already expect and what points at the attempted operation. A real permission problem, non-directory parent, or persistent missing parent therefore remains visible.

No third-party dependency, no TUI import changes, no configurable retry policy, and no sleeps on success. The implementation should be small enough to review by inspection: one helper file, one helper test file, and mechanical replacements at the run-path mkdir sites.
