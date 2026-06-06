---
agent: hermes
idea: tui-command-picker
phase: review
round: 2
date: 2026-06-06
---

## Summary
All round-01 items addressed; implementation now matches FINAL and introduces no regressions.

## Verification
- **FINAL §8 async refresh (codex MAJOR)**: FIXED — refreshPickerItems() rebuilds Items from openItems/answerItems, preserves Filter/Index/Offset, calls reclamp; wired into questionsMsg (339) and refreshHomeRuns (373). Test TestPickerAnswerRefreshesOnBackgroundUpdate asserts filter preservation, rebuild, and Enter selection on new list.
- **Answer lifecycle (agy #1)**: FIXED — submitInput no longer pre-clears composing/answerQID; answerQuestion clears only on success (1134-1135). Test TestPickerAnswerFailureKeepsCompose confirms failed write keeps compose state.
- **launchIdea on failure (agy #2)**: FIXED — keeps composing + inputText on err (1080-1082).
- **renderInputRow hint (agy #4)**: FIXED — picker.Active case (681) switches to picker hint.
- **renderPicker local reclamp (codex MINOR, agy #5)**: FIXED — p := m.picker; p.reclamp(visible) (1396-1397) is the single source; WindowSizeMsg also reclamps (299).
- **Label order (codex NIT)**: FIXED — openItems now emits "run   <run-id>  <idea>  [status]" (1344).
- **Mixed receivers (agy #6)**: DISMISSED — value receivers on Update-path methods match the rest of the file and Bubble Tea idiom; no divergence risk.

## New findings
None.

## Verdict
ACCEPT — all claims verified, tests green, no new bugs or stranded state paths.