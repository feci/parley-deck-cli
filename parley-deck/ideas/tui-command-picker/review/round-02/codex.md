---
agent: codex
idea: tui-command-picker
phase: review
round: 2
date: 2026-06-06
---

## Summary
ACCEPT — the round-01 blocking findings are addressed and I found no new regression.

## Verification
- codex MAJOR async refresh: FIXED. `questionsMsg` now calls `refreshPickerItems()` for an active `/answer` picker, and `refreshHomeRuns()` does the same for an active `/open` picker. The helper rebuilds from cached sources, preserves `Filter`, `Index`, and `Offset`, then re-clamps against the current filtered list. `TestPickerAnswerRefreshesOnBackgroundUpdate` covers preserved filter plus selecting from the rebuilt list.
- codex MINOR render reclamp: FIXED. `renderPicker` reclamps a local `pickerState` copy and renders from that copy, while model-mutating paths (`refreshPickerItems`, `updatePicker`, and `WindowSizeMsg`) reclamp the stored picker. I do not see a divergence bug on the next keypress: Enter re-clamps before indexing, and movement/filter mutations re-clamp before returning.
- codex MINOR missing async test coverage: FIXED for the blocking async answer path. `/open` refresh is also wired through `refreshHomeRuns`; the successful open reset path still relies on `activateRun`, which clears picker, compose state, `answerQID`, and input.
- codex NIT label order: FIXED. Run labels now put the run id before the idea slug, and idea labels put the slug before status.
- agy MAJOR answer failure lifecycle: FIXED. `submitInput` no longer clears `composing` or `answerQID` before `answerQuestion`; `answerQuestion` clears them only after `hitl.Answer` succeeds. `TestPickerAnswerFailureKeepsCompose` verifies failed writes keep compose state and the qid.
- agy MAJOR launch failure loses typed task: FIXED. `launchIdea` preserves `composing` and `inputText` on `Start` errors, allowing retry. The `Start == nil` branch still exits compose, but `N` cannot normally enter compose when launching is unavailable, so I do not treat that as a regression.
- agy MAJOR picker hint: FIXED. `renderInputRow` switches to the picker hint while `picker.Active`, before compose and slash-command hints.
- agy MINOR duplicated scroll math: FIXED. Window math is centralized in `pickerState.reclamp`; `renderPicker` uses a reclamped local copy, and `WindowSizeMsg` reclamps the model when a picker is open.
- agy NIT mixed receivers: dismissed. I agree with the dismissal: value-receiver Bubble Tea update methods mutate and return a model copy, while pointer helpers are used on that copy for localized state updates. This matches the file's existing idiom and is not a defect here.
- hermes round-01 ACCEPT: still valid. The fix-up changes do not weaken the key-routing, dedup, empty-state, or reset invariants hermes called out.

## New findings
None.

## Verdict
ACCEPT — no blocking items remain.
