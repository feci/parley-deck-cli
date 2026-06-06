---
agent: agy
idea: tui-command-picker
phase: review
round: 1
date: 2026-06-06
---

## Summary
needs-fixes

## Findings

### 1. Unconditional Clearing of Answer Composition State on Error
* **Severity**: MAJOR
* **Location**: [internal/tui/live.go:1036-1041](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1036-L1041) in `submitInput()`
* **Problem**: When a user selects a question from the `/answer` picker and enters answer composition mode (`composing=true`, `answerQID=<qid>`), pressing Enter triggers `submitInput()`. However, `m.composing` and `m.answerQID` are cleared *unconditionally* before calling `m.answerQuestion(qid, text)`. If the actual answer submission fails (e.g., due to write permissions on the HITL state, or validation failure), the model returns the error in `m.inputErr`, but the composition state has already been lost. The user is kicked out of composing mode, and if they correct the text and press Enter again, the input will be routed as a standard steer or command instead of an answer.
* **Fix**: Do not clear `m.composing` and `m.answerQID` in `submitInput()` before calling `m.answerQuestion()`. Instead, pass them to `m.answerQuestion()` (or let `answerQuestion` clear them) and only clear them on a *successful* write, conforming to the spec requirement: *"Every path that exits/supersedes compose or picker — ... successful answerQuestion ... — leaves composing=false, answerQID=""..."*

### 2. Failed `launchIdea` Clears Input Text
* **Severity**: MAJOR
* **Location**: [internal/tui/live.go:1074](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1074) in `launchIdea()`
* **Problem**: If launching an idea fails (e.g., due to a daemon or script error), `launchIdea()` clears `m.inputText = ""` and `m.composing = false` on failure. If the user typed a long, descriptive task, they lose their entire input and must re-type it from scratch.
* **Fix**: Do not clear `m.inputText` and `m.composing` on failure. Keep the user in composing mode with their text intact so they can edit it and retry.

### 3. Background Updates Do Not Refresh Active Picker Items
* **Severity**: MAJOR
* **Location**: [internal/tui/live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go) in `Update()` (`questionsMsg` handling) and `refreshHomeRuns()`
* **Problem**: The spec §8 requires: *"If a background eventsMsg/questionsMsg arrives while the picker is open, rebuild picker.Items WITHOUT resetting Filter or the cursor (re-clamp Index after rebuild)."* The implementation completely omits this. Background updates to questions or runs do not rebuild `m.picker.Items`, leaving the user with a stale list of candidates while the picker is active.
* **Fix**:
  - In `Update()` under the `questionsMsg` case: if `m.picker.Active && m.picker.Kind == pickerAnswer`, rebuild `m.picker.Items = m.answerItems()` and call `m.picker.reclamp(m.pickerRows())`.
  - In `refreshHomeRuns()`: if `m.picker.Active && m.picker.Kind == pickerOpen`, rebuild `m.picker.Items = m.openItems()` and call `m.picker.reclamp(m.pickerRows())`.

### 4. Missing Input-Row Hint when Picker is Active
* **Severity**: MAJOR
* **Location**: [internal/tui/live.go:671-685](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L671-L685) in `renderInputRow()`
* **Problem**: The spec §7 requires: *"Input-row hint switches to ↑/↓ select · type filter · Enter choose · esc cancel."* The implementation does not switch the hint when `m.picker.Active` is true. As a result, the input-row hint displays stale instructions (e.g., claiming that `↑/↓` switches tabs and `N` starts a new idea), which is incorrect and misleading since those keys are captured by the picker.
* **Fix**: Add a case at the top of the switch block in `renderInputRow()` to handle `m.picker.Active`:
  ```go
  case m.picker.Active:
      hint = "↑/↓ select · type filter · Enter choose · esc cancel"
  ```

### 5. Duplicated and Inconsistent Scroll Window Math
* **Severity**: MINOR
* **Location**: [internal/tui/live.go:145-176](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L145-L176) (`reclamp`) vs [internal/tui/live.go:1367-1381](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1367-L1381) (`renderPicker`)
* **Problem**: The logic to calculate the scroll window's top offset is duplicated. Furthermore, they are slightly inconsistent: `reclamp()` clamps `Offset` against `maxOff` at the very end of the calculation, whereas `renderPicker()` clamps it at the *beginning*. If the index checks adjust `off` at the end of `renderPicker()`, it remains unclamped, which can lead to out-of-bound slices if state is out of sync.
* **Fix**: Remove the duplicated scroll math from `renderPicker()`. Let `reclamp()` be the single source of truth. If the terminal height changes on a `WindowSizeMsg`, call `m.picker.reclamp(m.pickerRows())` in the `Update()` handler to keep the offset state correct.

### 6. Mixed Value and Pointer Receivers on `liveModel`
* **Severity**: NIT
* **Location**: [internal/tui/live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go)
* **Problem**: `liveModel` mixes value receivers (e.g., `updatePicker`, `selectPickerItem`, `submitInput`, `answerQuestion`) and pointer receivers (e.g., `clearPicker`, `clearComposition`, `openPicker`). Mixing receiver types on a mutating Bubble Tea model is error-prone because calling pointer methods on value copies will mutate temporary stack copies.
* **Fix**: Standardize on pointer receivers (`*liveModel`) for all methods that update model state or trigger commands.

## Verdict
REQUEST-CHANGES

**Blocking Items:**
- **Finding 1**: Unconditional Clearing of Answer Composition State on Error
- **Finding 2**: Failed `launchIdea` Clears Input Text
- **Finding 3**: Background Updates Do Not Refresh Active Picker Items
- **Finding 4**: Missing Input-Row Hint when Picker is Active
