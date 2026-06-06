---
agent: codex
idea: tui-command-picker
phase: review
round: 1
date: 2026-06-06
---

## Summary
needs-fixes — key routing and answerQID reset paths look correct, but the async picker item refresh/reclamp contract from FINAL §8 is not implemented.

## Findings

### MAJOR
Location: `internal/tui/live.go:304-337` (`Update` handling for `eventsMsg` / `questionsMsg`)

Problem: FINAL §8 requires that if an `eventsMsg` or `questionsMsg` arrives while the picker is open, `picker.Items` is rebuilt without resetting `Filter` or the cursor, followed by re-clamping. The implementation updates `m.events`, `m.state`, and `m.questions`, but never updates `m.picker.Items`. In an active `/answer` picker, a background question refresh can leave the picker showing stale open questions; Enter can then attempt to answer a question that is no longer open, or fail to show a newly opened question until the picker is closed and reopened. The same missing hook exists for `/open` picker source changes such as refreshed recent runs.

Fix: After successful `questionsMsg`, when `m.picker.Active && m.picker.Kind == pickerAnswer`, preserve `Filter`, `Index`, and `Offset`, replace `m.picker.Items = m.answerItems()`, then call `m.picker.reclamp(m.pickerRows())`. Add the equivalent rebuild for `pickerOpen` after any update path that changes its sources (`homeRuns` / status-derived candidates), including the relevant event/done refresh paths. Add a state-machine test where an active answer picker receives `questionsMsg` with a changed open-question set and asserts the filter is preserved, the cursor is clamped, and Enter selects from the rebuilt list.

### MINOR
Location: `internal/tui/live.go:1343-1405` (`renderPicker`)

Problem: FINAL §6 says to clamp `Index` again before render and compute `Offset` from the clamped index. `renderPicker` only clamps a local `off` variable; it does not clamp `m.picker.Index`, and because the receiver is by value it cannot repair model state. If the filtered list shrinks due to an async rebuild before the next keypress, the render can show no selected row until a key event runs `reclamp`.

Fix: Ensure the model is reclamped before render after every picker item rebuild, and preferably have `renderPicker` work from a local copy such as `p := m.picker; p.reclamp(m.pickerRows())`, then render from `p.Index` / `p.Offset`. This keeps render defensive without mutating state from a view method.

### MINOR
Location: `internal/tui/live_test.go:646-899` (`tui-command-picker tests`)

Problem: The tests cover most FINAL §9 state-machine cases and correctly stay terminal-free, but they do not assert the FINAL §8 async refresh contract at all. They also only exercise `/open` picker Enter through the failure path (`/open zzz` with no resolvable run), so successful picker-to-`activateRun` reset behavior is covered indirectly by older tests, not by the picker path itself.

Fix: Add tests for active picker rebuild on `questionsMsg` and for a successful `/open` picker selection against a real run summary, asserting `picker.Active=false`, `composing=false`, `answerQID=""`, and `inputText=""` after activation.

### NIT
Location: `internal/tui/live.go:1316-1318` (`openItems`)

Problem: The run label orders fields as idea/status/run-id, while FINAL §3's example puts the run id immediately after the visible `run` type. This is not a behavior bug because `Value` is correct and the user does not type the ID, but it makes the picker less aligned with the ratified display shape.

Fix: Prefer a run label like `run   <run-id>   <idea-slug>   [<status>]`, keeping the existing dedup by `Value`.

## Verdict
REQUEST-CHANGES — blocking item: the MAJOR async picker item rebuild/reclamp gap in `Update` must be fixed before acceptance.
