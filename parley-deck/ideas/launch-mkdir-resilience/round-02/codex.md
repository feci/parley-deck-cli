---
agent: codex
idea: launch-mkdir-resilience
round: 2
date: 2026-06-07
---

## Decision

I agree with the converged shape: add `internal/fsutil.MkdirAllResilient(path, perm)` and swap it into the launch/run directory creation sites. The helper must keep the healthy path as exactly one successful `os.MkdirAll` call: no `Stat`, no sleep, no allocation-heavy abstraction.

My final positions on the open forks:

1. **BACKOFF:** use the compromise: first retry/revalidation is immediate, later retries use `5ms`, `20ms`, `50ms`. That is four total `MkdirAll` attempts: initial attempt, immediate retry after failed `Stat`, then three slept retries. After each failed `MkdirAll`, `Stat+IsDir` is checked before deciding whether to continue. This preserves zero cost on first success, gives virtio-fs a short cache-settle window, and caps worst-case sleep at 75ms.
2. **SCOPE:** include the full `runcontrol.Create` launch path: `protocol.CreateIdea` idea-dir mkdir (`internal/protocol/workspace.go:67`), `store.Append` (`internal/store/events.go:41`), and `runmanifest.Write` (`internal/runmanifest/manifest.go:142`), plus the agreed runner/run sites (`runner.go:340,343`, `round_index.go:85`, `steer.go:152`, `phase58.go:58`, `handoff.go:36`) and `hitl.go:183`. Keep unrelated repo-wide mkdirs out.
3. **EEXIST:** do not treat `fs.ErrExist`/`os.IsExist` as success blindly. It is success only if `os.Stat(path)` succeeds and `info.IsDir()` is true. A colliding regular file must still fail.
4. **FAIL-FAST:** agree with fail-fast for `fs.ErrPermission`, but only after the immediate `Stat+IsDir` check. If the target is confirmed as a directory, success wins; otherwise a permission error should return immediately and not burn the retry window.
5. **OPEN FILE:** keep the trailing `os.OpenFile` in `Store.Append` unchanged. Retrying open/write in append code risks duplicate events or muddier semantics. A confirmed directory is enough for this fix; any later proven open-specific cache bug should get a separate, narrower design.
6. **TEST PLAN:** confirm the listed fsutil seam tests: transient-then-success; dir-exists-via-`Stat` without exhausting retries; already-exists; genuine-failure-returns-last-mkdir-error; fail-fast-permission; common-path-no-stat-no-sleep. Add non-directory collision as the concrete `EEXIST` regression case if it is not covered by "already-exists".

## Responses

**claude:** Agree on `internal/fsutil`, bounded retry, no TUI attach path, and run-path-wide use. I disagree with `20/40/80ms`: it is more latency than needed for the first mitigation pass. Counter-proposal: immediate retry, then `5/20/50ms`. I also disagree with keeping `protocol.CreateIdea` out if that was implied by "run path"; it executes before `Append` during launch and can produce the same spurious launch failure.

**codex:** I keep my round-01 positions on `EEXIST` requiring `Stat+IsDir`, not touching `OpenFile`, preserving the last mkdir error, and avoiding a repo-wide sweep. I change scope from excluding protocol idea creation to including `protocol.CreateIdea`, because the user explicitly identified it as inside `runcontrol.Create` before `Append`; that makes it launch-path, not general bootstrap.

**agy:** Agree on `internal/fsutil`, small backoff only on the error path, full run-path coverage, and fail-fast permission handling. I disagree with treating `fs.ErrExist` as immediate success. Counter-proposal: `ErrExist` is only a reason to `Stat`; success requires `IsDir`. I also disagree with adding cleanup or auto-attach UX in this change; prevent the false launch error first, and defer orphan recovery unless failures remain.

**hermes:** Agree on a narrow helper, no generic file-operation retry, and no TUI recovery. I disagree on zero sleep and immediate `ErrExist` success. Counter-proposal: one immediate retry covers the cheap revalidation argument, then `5/20/50ms` gives weak caches time to settle without making genuine errors hang. `ErrExist` still needs `Stat+IsDir` to avoid masking file collisions.

## Current Proposal

`MkdirAllResilient` should:

- call `os.MkdirAll` once and return nil immediately on success;
- after any error, call `os.Stat(path)` and return nil only for an existing directory;
- if not a directory and the mkdir error is permission, return that error immediately;
- retry with no sleep before the first retry, then `5ms`, `20ms`, `50ms`;
- after the final failure, do one last `Stat+IsDir` check and otherwise return the last `MkdirAll` error.

This is surgical, launch-path complete, and testable without a flaky mount.
