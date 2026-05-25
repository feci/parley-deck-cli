---
agent: claude
idea: tui-action-execution
review-round: 1
date: 2026-05-25
reviewed-implementation: IMPLEMENTATION.md
---

## Findings

### 1. MAJOR — `request-signoffs` bypasses `RequiresYes` confirmation in the TUI

**File:** `internal/app/app.go:1795`

`consensusActionArgs` unconditionally injects `--yes` for `KindRequestSignoffs`. The FINAL.md explicitly defers a confirmation modal for `RequiresYes` actions, but the current implementation silently auto-confirms without any user-facing warning. If a user accidentally presses `enter` on a `request-signoffs` action that targets hosted backends, agents will be launched and billed with no chance to cancel.

The planner already sets `RequiresYes: true` on these actions (`internal/runplan/runplan.go:149`). The TUI should either refuse to execute actions with `RequiresYes` until the modal exists, or at minimum display a prominent warning in `actionText` before dispatching.

**Suggested mitigation:** Return an advisory result with the manual command when `action.RequiresYes` is true, matching the pattern used for `answer-question` and `retry-agent`.

---

### 2. MAJOR — No test for `actionRunning` busy-state rejection

**File:** `internal/tui/app.go:978-979`

When `actionRunning` is true, `triggerSelectedAction` sets `actionText` and returns without a command. However, the Bubble Tea `Update` method returns the model as a value type (`tea.Model`), and `triggerSelectedAction` is called on a value receiver. The mutation to `actionText` on line 979 happens on a local copy of `m` which is then returned — this works correctly. But there is no test asserting that pressing `enter` twice rapidly (before the first action completes) is blocked. The `TestDashboardEnterRunsSelectedAction` test sets `actionRunning` implicitly via the first `enter`, but then immediately processes the result message in the same flow, so the busy guard is never exercised.

---

### 3. MINOR — `actionText` is never cleared between successful actions

**File:** `internal/tui/app.go:289-305`

After a successful `actionRunMsg`, `actionText` is set to the result message. If the user subsequently triggers a different action, `actionText` is overwritten to "Running action: ..." on line 996, which is fine. But if the user navigates away from the actions panel and back, stale feedback text from a prior action remains visible. There is no mechanism to clear `actionText` when focus changes or when a new run is selected.

This is a minor UX annoyance rather than a correctness bug. Clearing `actionText` in `moveSelection` when `focusIdeas` changes (line 708-709) would prevent confusion.

---

### 4. MINOR — `runaction.Command` returns empty string for `KindRetryAgent` but no test covers the TUI's advisory path end-to-end

**File:** `internal/runaction/action.go:71`, `internal/app/app.go:1748-1750`

`TestCommandLeavesRetryUnmapped` (action_test.go:30) confirms the command builder returns `""` for retry. `TestRunTUIActionReturnsAdvisoryForRetry` (app_test.go:1207) confirms the app runner returns an advisory message. But there is no TUI-level test confirming that pressing `enter` on a retry action in the TUI shows the advisory message in the rendered view.

---

### 5. MINOR — `RoundNumber` silently returns `"0"` for `"round-00"`, which may not be a valid round

**File:** `internal/runaction/action.go:74-84`

`RoundNumber("round-00")` strips all leading zeros and returns `"0"`. Passing `--round 0` to `parley consensus draft` may be unexpected since rounds are 1-indexed in the protocol. There is no validation that the round is positive, and no test for the `"round-00"` edge case.

---

### 6. NIT — Duplicate helper functions across packages

**Files:** `internal/tui/app.go:1075` and `internal/app/app.go:1927`

Both `tui` and `app` packages define their own `valueOr(value, fallback string) string` with identical implementations. Similarly, `agentProgress`, `agentDuration`, and `displayRunState` are duplicated between the two packages. This is acknowledged as pre-existing and not introduced by this slice, but worth noting for future deduplication.

---

### 7. NIT — `commandOutput` truncation at 240 characters may clip consensus status

**File:** `internal/app/app.go:1810-1813`

`commandOutput` hard-truncates at 240 characters. Consensus draft/signoff output from `printConsensusSummary` can exceed this if multiple participants are listed. The truncation is applied before the result is shown in the TUI's `actionText` area. This is cosmetic but could hide whether a consensus operation fully succeeded.

---

## Summary

The implementation is well-structured and matches the accepted FINAL.md closely. The `runaction` package provides clean command construction, the TUI integration follows established Bubble Tea patterns correctly, and the app-backed runner dispatches consensus operations through the existing CLI layer rather than duplicating logic.

The most significant concern is the `--yes` bypass for `request-signoffs` (finding 1), which circumvents the safety gate designed to prevent accidental hosted-backend launches. The busy-state guard (finding 2) works correctly in code but lacks test coverage, making it fragile against future refactors.

Test coverage is otherwise solid: command mapping, focus cycling, action selection, runner invocation, refresh-on-success, and advisory feedback paths are all exercised.
