---
agent: codex
idea: launch-mkdir-resilience
round: 1
phase: review
date: 2026-06-07
---

## Summary

The helper implementation in `internal/fsutil/fsutil.go` matches FINAL.md: healthy success is one `mkdirAll` call with no stat/sleep, the retry loop is bounded at 5 total mkdir attempts with delays `[0,5ms,20ms,50ms]`, `isDir` is checked before permission fail-fast, and exhaustion returns the last mkdir error. The scoped call-site swaps are present and the explicitly out-of-scope sites remain unchanged. I do not see a source-level correctness regression in the mkdir implementation.

## Findings

### CRITICAL

No findings at CRITICAL.

### MAJOR

#### `go test ./...` is not green in this workspace

`internal/runner/durablekill_test.go:116`

The implementer claims `go test ./...` is green, but in this `darwin/arm64` workspace it fails consistently:

```text
--- FAIL: TestDurableKillEndToEndRealProcess
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
```

I reproduced it with `go test ./internal/runner -run TestDurableKillEndToEndRealProcess -count=1 -v`. This failure is outside the mkdir-resilience files and does not look caused by this change, but the release verification claim is false in the current workspace. Before finalizing, either make the full suite pass here or explicitly document why this environment-specific existing failure is accepted and provide a green supported-environment/CI run.

### MINOR

#### Some fsutil tests do not fully lock the stated semantics

`internal/fsutil/fsutil_test.go:97`, `internal/fsutil/fsutil_test.go:123`, `internal/fsutil/fsutil_test.go:140`

The seven tests cover the main behavior and the immediate-retry path, but a few locked details are weaker than their claims:

- `Test_GenuineFailure` returns the same `errBoom` from every mkdir attempt, so it proves exhaustion and sleep counts, but not that the returned error is specifically the last mkdir error.
- `Test_FailFastPermission` uses `stat = notExist`, so it does not prove the required ordering that `isDir` wins before permission fail-fast. A `mkdirAll -> fs.ErrPermission` plus `stat -> existing dir` case would lock that.
- `Test_NonDirCollision` returns `errBoom`, not `fs.ErrExist`, so it does not specifically catch an implementation that blindly trusts `fs.ErrExist`.

Suggested fix: add narrow seam-driven cases with per-attempt distinct errors, permission-plus-existing-dir, and `fs.ErrExist` plus regular-file stat. This is test hardening only; the current helper code implements these semantics correctly.

### NIT

No findings at NIT.

## Scope And Call Sites

All FINAL.md call sites are swapped: `internal/protocol/workspace.go:69`, `internal/store/events.go:43`, `internal/runmanifest/manifest.go:143`, `internal/runner/runner.go:341` and `internal/runner/runner.go:344`, `internal/runner/round_index.go:86`, `internal/runner/steer.go:153`, `internal/runner/phase58.go:59`, `internal/runner/handoff.go:37`, and `internal/hitl/hitl.go:184`.

The intended exclusions are still stdlib mkdirs: `internal/protocol/workspace.go:46` (`InitWorkspace`), `internal/runner/runner.go:862`, and `internal/runner/runner.go:902`; `os.MkdirTemp` also remains for isolated homes. Remaining `os.MkdirAll` sites are in tests or FINAL.md out-of-scope areas such as pipeline, consensus, sessionstore, driver, and the app probe.

`internal/fsutil` is stdlib-only, and `internal/protocol` importing it does not create an import cycle. The unchanged `os.OpenFile` in `Store.Append` is acceptable per FINAL.md: retrying the append itself would risk duplicate events, and the resilient mkdir now stat-confirms the directory before the open.

## Verification

- `go build ./...`: PASS
- `go vet ./...`: PASS
- `go test ./internal/fsutil -count=1 -v`: PASS
- `go test ./...`: FAIL, reproducible at `internal/runner/TestDurableKillEndToEndRealProcess`

## Overall Verdict

REQUEST-CHANGES
