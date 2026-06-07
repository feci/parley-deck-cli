---
idea: launch-orphan-hardening
phase: implementation
date: 2026-06-07
implementer: claude
status: ready-for-review
---

## Changes
- `internal/fsutil/fsutil.go` — `retryDelays` → `[0,15ms,35ms,100ms,250ms,500ms,1000ms]`
  (8 total attempts, ~1.9s worst case, error-path only; healthy path still one os.MkdirAll
  with no Stat/no sleep). Comment updated.
- `internal/runcontrol/runcontrol.go` — `runmanifest.Write` is now called via a
  `var writeManifest = runmanifest.Write` seam; on failure, Create appends a
  `run.manifest_deferred` event (best-effort, reusing the existing run store) and CONTINUES
  instead of `return CreatedRun{}, err`. The run is no longer orphaned by a manifest hiccup.
- `internal/runcontrol/runcontrol_test.go` — new `TestCreateBestEffortManifest`.
- `internal/fsutil/fsutil_test.go` — `Test_GenuineFailure` updated for the 8-attempt
  schedule (distinct last error on attempt 8; sleeps 15/35/100/250/500/1000ms).
- `internal/app/version.go` + `VERSION` → 1.22.0.

## Notes / open questions for reviewers
- Best-effort `run.manifest_deferred` Append itself does `MkdirAllResilient(runDir)`; if
  that also fails it is ignored (`_ =`) — the run still proceeds (no orphan). Acceptable?
- Is ~1.9s the right window for a virtio-fs ~1s cache timeout, or too long/short?
- Confirm no caller HARD-requires `run.json` (claude verified runstate + app.go:880/940
  degrade gracefully — please double-check resume/continue/driver).

## Verification
gofmt clean; go build/vet ./... OK; go test ./... -count=1 green (incl. new test).
