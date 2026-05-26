---
reviewer: claude
idea: tui-layout-refresh
round: 1
date: 2026-05-26
verdict: conditional-approve
---

# Review: tui-layout-refresh

Implementation matches the consensus direction: height-aware compact rendering, two-column dashboard, semantic color badges, and focused compact tests. The core rendering logic is sound. Findings below are all addressable without rearchitecting.

## Findings

### 1. Dead code: three orphaned dashboard render methods (medium)

`renderEvents()` (app.go:782), `renderRunDetails()` (app.go:799), and `renderQuestions()` (app.go:844) are no longer called from anywhere in the package. The old `View()` called them; the new `renderDashboard()` and `renderCompactDashboard()` use different methods. Per the project's surgical-changes guideline, code made unused by this change should be removed.

**Action:** delete the three methods, confirm `go build ./internal/tui` still compiles.

### 2. Removed footer safety text without consensus basis (low-medium)

The informational line `"Quit closes the TUI and cancels TUI-started runs; it does not detach child processes."` was removed from `renderFooter()`. The consensus specifies "a single compact footer line" but does not call for removing this user-facing process-lifecycle warning. The horizontal `|`-joined footer already has room for it as a second segment, or it could be kept in normal mode only and dropped in compact mode.

**Action:** restore the text (at least in normal mode) or confirm with the group that its removal is intentional.

### 3. Inconsistent compact layout=compact label styling (low)

Dashboard compact embeds `layout=compact` inside `headerStyle.Render()` via `dashboardHeader("compact")`. Live compact appends `mutedStyle.Render("layout=compact")` outside the header render call. The result is a visual mismatch: same label, different colors/styling between the two views.

**Action:** unify the approach. Embedding inside the header string (dashboard pattern) is cleaner.

### 4. Missing threshold-boundary test (low)

Compact tests use heights well below the 30-line threshold (20 for dashboard, 18 for live). No test verifies behavior at the exact boundary (`height == 30` should be normal mode, `height == 29` should be compact). An off-by-one in the `<` comparison would only surface at the boundary.

**Action:** add one test case at `height = compactDashboardHeight` confirming normal mode is used, and one at `height = compactDashboardHeight - 1` confirming compact mode.

### 5. Footer horizontal join can exceed terminal width (low)

`renderFooter()` changed from `"\n"` to `"  |  "` as the separator. When error text is present, the single resulting line can easily exceed terminal width without any wrapping or truncation.

**Action:** consider truncating the joined footer to `width` or keeping errors on a separate line.

### 6. "Latest events" heading renders only one event in compact live (cosmetic)

`renderCompactLatestEvent()` renders exactly one event. The section title says "Latest events" (plural). Either rename to "Latest event" or show 2-3 events.

## No blocking findings

All findings are low-to-medium severity. The core layout refresh works correctly: height propagation, compact thresholds, clipping, badge rendering, and test coverage of the happy path are all solid.

Conditional approve: merge after addressing finding 1 (dead code cleanup) and deciding on finding 2 (footer text removal).
