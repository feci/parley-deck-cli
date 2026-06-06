---
agent: agy
idea: tui-command-picker
phase: review
round: 2
date: 2026-06-06
---

## Summary
ACCEPT: All round-01 findings are fully resolved, with robust tests added, and the implementation strictly adheres to the FINAL.md spec.

## Verification

### Finding 1: Unconditional Clearing of Answer Composition State on Error
* **Status**: FIXED
* **Detail**: `submitInput()` no longer pre-clears `composing` and `answerQID`. Instead, `answerQuestion()` only clears them on a successful write. If HITL writing fails, the user remains in composing mode with their text intact. The new test `TestPickerAnswerFailureKeepsCompose` verifies this behavior.

### Finding 2: Failed `launchIdea` Clears Input Text
* **Status**: FIXED
* **Detail**: `launchIdea()` now keeps the user in composition mode and preserves `m.inputText` on a launch failure, preventing data loss.

### Finding 3: Background Updates Do Not Refresh Active Picker Items
* **Status**: FIXED
* **Detail**: Background updates for `questionsMsg` (during `pickerAnswer`) and `refreshHomeRuns()` (during `pickerOpen`) now invoke `refreshPickerItems()`, which rebuilds the picker candidate list from the cached state and calls `reclamp()` without resetting the search `Filter` or selection `Index`. This is covered by the new state-machine test `TestPickerAnswerRefreshesOnBackgroundUpdate`.

### Finding 4: Missing Input-Row Hint when Picker is Active
* **Status**: FIXED
* **Detail**: `renderInputRow()` now checks `m.picker.Active` first to switch the displayed hint to `"↑/↓ select · type filter · Enter choose · esc cancel"`, resolving the issue with stale hints.

### Finding 5: Duplicated and Inconsistent Scroll Window Math
* **Status**: FIXED
* **Detail**: Duplicate scroll calculations in `renderPicker()` have been removed. It now operates on a local copy `p` that is defensively reclamped, ensuring `reclamp()` acts as the single source of truth for window math. Furthermore, `WindowSizeMsg` now triggers `reclamp()` when the picker is active to handle terminal resizes correctly.

### Finding 6: Mixed Value and Pointer Receivers on `liveModel`
* **Status**: DISMISSED (AGREE)
* **Detail**: As discussed in the round-01 follow-up, the mix of value-receiver `Update`/view helpers and pointer-receiver mutating methods is consistent with the established idiom of this file and is safe since value copies returned from `Update` undergo pointer-resolution automatically for mutating helper calls.

### Finding 7: Run/Idea Label Order
* **Status**: FIXED
* **Detail**: The picker items for runs and ideas now correctly format their label orders (e.g., `run  <run-id>  <idea>  [status]`) to align with the spec layout.

## New findings
None.

## Verdict
ACCEPT
