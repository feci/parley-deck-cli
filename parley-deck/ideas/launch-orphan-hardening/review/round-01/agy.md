---
agent: agy
idea: launch-orphan-hardening
round: 1
phase: review
date: 2026-06-07
---

## Summary

The implemented changes successfully address the launch orphaning issue on weakly-coherent mounts (like virtio-fs) under heavy load. The solution correctly approaches the problem from two directions: extending the resilience window of the filesystem creation helpers and making the manifest writing step best-effort (non-fatal) to ensure that the run is never orphaned once the initial directory and events file have been created. The implementation is highly robust, compiles cleanly, and is fully validated by updated and new unit tests.

## Findings

### [CRITICAL]
No findings at CRITICAL.

### [MAJOR]
No findings at MAJOR.

### [MINOR]
No findings at MINOR.

### [NIT]
1. **Error Swallowing in `runStore.Append` for Deferred Manifest**: In [internal/runcontrol/runcontrol.go:87](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runcontrol/runcontrol.go#L87), if the manifest write fails, we try to append a `run.manifest_deferred` event to `events.jsonl` and ignore the error (`_ =`). While ignoring the failure here is the correct strategy (we should not crash or abort the launch if our fallback audit event fails to write), it does mean `MkdirAllResilient` will run again under the hood. Since the `run.created` event has already been successfully appended, the directory is highly likely to exist and the call should return immediately. There are no performance or logical correctness concerns here, but adding a brief code comment explaining the choice to swallow this error helps future maintainers.

## Verification Checklist & Detailed Analysis

### 1. Correctness of Best-Effort Manifest
* **Safety without `run.json`**: Verified. The source of truth for the run is `events.jsonl`.
* **Caller Requirements**:
  * `runstate.loadManifestSnapshot` in [internal/runstate/runstate.go:180](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runstate/runstate.go#L180) returns `(Manifest, false)` if the manifest file is absent. In `LoadRunAt` ([internal/runstate/runstate.go:107](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runstate/runstate.go#L107)), this only skips applying manifest defaults. The fields `IdeaSlug`, `Task`, `Mode`, and `Participants` are subsequently populated by scanning the `run.created` event in `events.jsonl` (lines 121-134), while the other fields (`Outcome`, `Questions`, `NextActions`, etc.) are dynamically projected from the folder and events stream.
  * `inspectSession` in [internal/app/app.go:880](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L880) handles manifest load errors gracefully: `os.ErrNotExist` leaves `payload.Manifest` as `nil` and `payload.ManifestError` empty. `printSessionDetail` ([internal/app/app.go:927](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L927)) then prints a fallback message.
  * `sessionStatus` in [internal/app/app.go:940](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L940) falls back to checking `session.Terminal` or returns `"active"` if the manifest is missing.
  * `continueAuto` and `driver.New` (such as in [internal/app/app.go:1070](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L1070)) rely solely on the reconstructed `RunSummary` properties and do not read or require the manifest file.
* **Error Hiding**: Making this non-fatal does not mask real filesystem failures. Any severe write issues (e.g. read-only filesystem or out of space) will already fail on the preceding `runStore.Append(run.created)` call in [internal/runcontrol/runcontrol.go:55](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runcontrol/runcontrol.go#L55), aborting the launch cleanly.

### 2. Retry Window of `MkdirAllResilient`
* **Virtio-fs Cache Outlasting**: The new `retryDelays` schedule in [internal/fsutil/fsutil.go:29-37](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/fsutil/fsutil.go#L29-37) runs 8 attempts with a worst-case duration of 1.9s (0ms + 15ms + 35ms + 100ms + 250ms + 500ms + 1000ms), which is more than enough to outlast typical virtio-fs cache timeouts (~1s).
* **Healthy Path Latency**: Unaffected. If the initial `os.MkdirAll` succeeds, the function returns immediately without executing `isDir`, checking permissions, or sleeping.
* **Permission Errors**: Still fail fast. `errors.Is(err, fs.ErrPermission)` is checked inside the loop at [internal/fsutil/fsutil.go:58](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/fsutil/fsutil.go#L58) before sleeping, ensuring permission errors exit instantly.

### 3. Deferral Append behavior
* The fallback `_ = runStore.Append(...)` is safe. Since the directory has already been successfully created in the prior append, the `MkdirAllResilient` inside `Append` will return immediately. If it does fail for an unexpected system reason, swallowing the error ensures we avoid failing the launch of a run that is already initialized and registered.

### 4. Tests Verification
* `TestCreateBestEffortManifest` in [internal/runcontrol/runcontrol_test.go:72](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runcontrol/runcontrol_test.go#L72) correctly mocks `writeManifest`, cleans up the seam via `t.Cleanup`, asserts `Create` succeeds with no error, verifies `run.json` is missing, and checks that both `run.created` and `run.manifest_deferred` events are written to the events file.
* `Test_GenuineFailure` in [internal/fsutil/fsutil_test.go:98](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/fsutil/fsutil_test.go#L98) has been correctly updated for 8 attempts. The mock `mkdirAll` receives 8 calls, returns the unique `errLast` on attempt 8, and verifies that the 6 non-zero sleep delays (`[15ms, 35ms, 100ms, 250ms, 500ms, 1000ms]`) match.

### 5. Build/Regression
* Both `go vet ./...` and `go build ./...` compile cleanly without issues.
* All unit tests pass (`go test ./... -count=1` is 100% green).
* No regressions are introduced to normal run flows, TUI interactions, status checks, resume/continue paths, or workspace sessions.

## Verdict

ACCEPT
