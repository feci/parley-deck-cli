---
agent: codex
idea: launch-orphan-hardening
round: 1
phase: review
date: 2026-06-07
---

## Findings

- No findings at CRITICAL.
- No findings at MAJOR.
- No findings at MINOR.
- [NIT] `internal/fsutil/fsutil_test.go:95` still describes `Test_GenuineFailure` as sleeping `5ms,20ms,50ms`, but the implemented and asserted schedule is now `15ms,35ms,100ms,250ms,500ms,1000ms` at `internal/fsutil/fsutil_test.go:123`.
- [NIT] `TestCreateBestEffortManifest` proves there is no loadable manifest, but not strictly that `run.json` is absent: `internal/runcontrol/runcontrol_test.go:97` accepts any `runmanifest.Load` error. The seam makes absence the expected outcome here, but an `os.Stat` / `errors.Is(err, os.ErrNotExist)` assertion would match the test name and implementation note more precisely.

## Review Notes

Best-effort manifest correctness: I do not see a path that hard-requires `run.json` after `runcontrol.Create` has successfully appended `run.created`. The new code keeps the initial event append fatal at `internal/runcontrol/runcontrol.go:55`-`67`, then makes only the manifest write best-effort at `internal/runcontrol/runcontrol.go:68`-`92`; the returned `CreatedRun` and `runner.Options` are built from `idea`, `runID`, and `runStore`, not from a loaded manifest, at `internal/runcontrol/runcontrol.go:96`-`110`. `registerSession` also uses those same in-memory values at `internal/runcontrol/runcontrol.go:173`-`187`.

The main projection path degrades correctly. `runstate.LoadRunAt` treats the manifest as optional defaults at `internal/runstate/runstate.go:105`-`120`, then derives idea/task/mode/participants from `run.created` at `internal/runstate/runstate.go:121`-`134`; `loadManifestSnapshot` explicitly returns `(empty, false)` on any load error at `internal/runstate/runstate.go:180`-`185`. `ListRuns` and `ResolveRun` operate through `LoadRunAt` / `LoadRun` at `internal/runstate/runstate.go:265`-`323`.

The CLI paths I checked also tolerate a missing manifest. `status --run` and workspace status resolve/load through `runstate` at `internal/app/app.go:623`-`634` and `internal/app/app.go:670`-`680`. `sessions inspect` loads `runmanifest` only for optional display, records non-ENOENT errors, and still loads the run state at `internal/app/app.go:874`-`892`; the text path prints a missing-manifest message and continues to `printRunDetail` at `internal/app/app.go:924`-`935`. `sessionStatus` falls back to `"active"` / `"terminal"` when manifest load fails at `internal/app/app.go:939`-`947`. `resume` and `continue` both resolve via `runstate.ResolveRun` at `internal/app/app.go:983`-`997` and `internal/app/app.go:1041`-`1064`; `continue --auto` reconstructs the driver from `runstate.RunSummary` and `store.New(run.RunDir)` at `internal/app/app.go:1070`-`1108`.

The normal launch paths do not regress. `parley run --no-tui` starts `runner.RunRoundOne` from `created.RunOptions` at `internal/app/app.go:1706`-`1717`, and the TUI path starts `runner.RunRoundOneAsync` with the same options at `internal/app/app.go:1755`-`1784`. The TUI launch-from-home path similarly calls `runcontrol.Create`, then starts the async runner and returns `LaunchResult` fields from `created` / `handle` at `internal/app/app.go:2113`-`2141`.

The driver has no manifest dependency. Its config is built around `RunDir`, `Root`, and the event store at `internal/driver/driver.go:40`-`48`; `Advance` rebuilds from idea artifacts at `internal/driver/driver.go:100`-`122`, and the cursor rebuild reads idea/final/implementation/review artifacts only at `internal/driver/cursor.go:85`-`121`. Round completion reads and appends through the event store at `internal/driver/driver.go:174`-`229`.

Retry window: the 8-attempt schedule is reasonable for the stated virtio-fs / weak-coherence failure. `retryDelays` is now `0,15,35,100,250,500,1000ms` at `internal/fsutil/fsutil.go:29`-`37`, so the sleeps total about 1.9s after the initial attempt. The healthy path still returns immediately on first `os.MkdirAll` success at `internal/fsutil/fsutil.go:49`-`52`, and permission errors still fail before sleeping at `internal/fsutil/fsutil.go:54`-`60`. The updated failure test asserts 8 attempts and the six positive sleeps at `internal/fsutil/fsutil_test.go:120`-`130`.

Deferred audit append: ignoring the second append error at `internal/runcontrol/runcontrol.go:87`-`91` is acceptable. `store.Append` itself recreates the run directory via `MkdirAllResilient` at `internal/store/events.go:36`-`46`, so this gives the audit event a chance to land without reintroducing the same launch-fatal mkdir problem. If it still fails, the already-written `run.created` event remains the source of truth; treating the audit event as mandatory would recreate the orphaning failure.

Tests: `TestCreateBestEffortManifest` does prove the important behavior: the seam is installed and restored at `internal/runcontrol/runcontrol_test.go:81`-`83`, `Create` succeeds at `internal/runcontrol/runcontrol_test.go:85`-`93`, and both `run.created` and `run.manifest_deferred` are required at `internal/runcontrol/runcontrol_test.go:100`-`115`. The `writeManifest` seam is restored with `t.Cleanup`, and no tests use `t.Parallel`, so I do not see a leak between package tests.

Verification run locally: `gofmt -l` on touched Go files printed nothing, `git diff --check` printed nothing, `GOCACHE=/private/tmp/parley-go-cache go build ./...` passed, `GOCACHE=/private/tmp/parley-go-cache go vet ./...` passed, and `GOCACHE=/private/tmp/parley-go-cache go test ./internal/fsutil ./internal/runcontrol -count=1` passed. Full `go test ./... -count=1` did not pass on this host because `internal/runner` failed `TestDurableKillEndToEndRealProcess` at `internal/runner/durablekill_test.go:116` with "process verification failed (no recorded boot id); not killed"; that package was not touched by this diff, so I am treating it as an existing host-specific verification failure rather than a regression in this change.

Overall verdict: ACCEPT
