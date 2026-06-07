---
agent: agy
idea: launch-mkdir-resilience
round: 1
phase: review
date: 2026-06-07
---

# Review of launch-mkdir-resilience (round-01)

## 1. CORRECTNESS of MkdirAllResilient
No findings at CRITICAL.
No findings at MAJOR.
No findings at MINOR.
- NIT: In [internal/fsutil/fsutil.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/fsutil/fsutil.go#L26), the `retryDelays` slice is defined with `0` as the first delay to trigger an immediate first retry. While correct and aligned with the FINAL spec, adding a brief comment above or inline explaining that the `0` value explicitly implements the immediate retry without sleeping would improve long-term readability.

## 2. TESTS
No findings at CRITICAL.
No findings at MAJOR.
No findings at MINOR.
No findings at NIT.
- Note: The tests in [internal/fsutil/fsutil_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/fsutil/fsutil_test.go) are well-structured. Package-level seams are safely saved and restored using `t.Cleanup` in `saveSeams(t)`. The assertions are strong and precise, specifically validating that the healthy path has zero stats/sleeps, the transient path sleeps zero times for the immediate retry, and a collision with a regular file is not masked.

## 3. CALL SITES
No findings at CRITICAL.
No findings at MAJOR.
No findings at MINOR.
No findings at NIT.
- Note: All 10 swaps to `fsutil.MkdirAllResilient` in the launch/run paths listed in FINAL.md are correctly implemented. `InitWorkspace` in [internal/protocol/workspace.go:46](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/workspace.go#L46) and the isolated temp homes in [internal/runner/runner.go:862](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L862) and [runner.go:902](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L902) are correctly left as `os.MkdirAll`. Imports are clean, alphabetical, and introduce no cycle.

## 4. SCOPE/SAFETY
No findings at CRITICAL.
No findings at MAJOR.
No findings at MINOR.
No findings at NIT.
- Note: The trailing `os.OpenFile` in `Store.Append` in [internal/store/events.go:46](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/store/events.go#L46) remains unchanged, which is the correct design. The preceding `MkdirAllResilient` call acts as a synchronization barrier, so the file open succeeds without risking duplicate appends. The helper depends only on the standard library.

## 5. BUILD
No findings at CRITICAL.
No findings at MAJOR.
No findings at MINOR.
No findings at NIT.
- Note: Verified that `go build ./...`, `go vet ./...`, and `go test ./...` all pass cleanly without errors.

## Overall verdict
ACCEPT
