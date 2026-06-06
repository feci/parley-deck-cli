---
idea: tui-command-picker
phase: implementation
status: complete
implementer: claude
date: 2026-06-06
---

# IMPLEMENTATION — TUI command picker

Implements FINAL.md. Pure TUI change in `internal/tui/live.go` (+ tests in
`live_test.go`). No new dependencies; runner/driver/agents/`--no-tui` untouched.

## What was built

**State** (`live.go`):
- `pickerKind` (`pickerOpen`/`pickerAnswer`), `pickerItem{Label, Value}`,
  `pickerState{Active, Kind, Title, Items, Index, Filter, Offset}` with
  `filtered()` (case-insensitive substring over Label+Value) and `reclamp(visible)`
  (keeps Index in range + a scroll window containing it).
- `liveModel` gains `picker pickerState` and `answerQID string`. No callbacks;
  dispatch is by `picker.Kind` in `selectPickerItem`.

**Routing** (`updateMain`): `ctrl+c` handled first (global), then
`if m.picker.Active { return m.updatePicker(msg) }` BEFORE the existing
`esc`/`N`/arrow/`enter`/printable switch. `updatePicker` owns only
↑/↓/Enter/esc/printable/backspace; all other keys (PgUp/PgDn/ctrl+u/d/Home/End/
shift+arrows/tab/left/right) are no-ops while active. The `N` guard keeps the
documentation term `&& !m.picker.Active`.

**Triggers** (`runCommand`): bare `/open` → `openPicker(pickerOpen, …, openItems())`;
bare `/answer` → `openPicker(pickerAnswer, …, answerItems())`. Empty candidates set
`inputErr` ("nothing to open yet" / "no open questions") and do NOT open the picker.
Explicit `/open <slug|run>` and `/answer <qid> <text>` unchanged.

**Two-step `/answer`**: `selectPickerItem` for `pickerAnswer` sets `answerQID`,
`composing=true`, clears input + picker. `renderInputRow` shows `answer <qid> › `
(warning style) when `composing && answerQID != ""`. `submitInput` routes to
`answerQuestion(qid, text)` when `answerQID != ""`, else `launchIdea`.

**State-reset invariant**: `clearPicker()` and `clearComposition()` helpers;
`activateRun` and the `esc`-from-composing path both clear
`composing`+`answerQID`+`picker`+`inputText`.

**Render** (`renderPicker`): a bordered list above the input row, inserted in
`renderTabbed` before the status line; the transcript area is shrunk by the picker's
height so tiny terminals don't overflow. Visible rows = `min(8, transcriptHeight())`
with `↑ more`/`↓ more` boundary markers, defensively clamped each render. Empty
filtered list → muted message. The tab strip is dimmed (`styleTab` checks
`picker.Active`) while picking. Input hint switches to
`↑/↓ select · type filter · Enter choose · esc cancel`; `/help` and the slash hint
advertise the bare-command pickers.

**Data freshness**: uses cached `homeRuns` / `questions`; no fetch on open.

## Tests (`live_test.go`, all green, no terminal)
Covers FINAL §9 items 1-11: bare-vs-explicit `/open`; empty candidates don't open;
arrows move selection not tab; printable (incl. `N`,`/`) filters not inputText;
`reclamp` clamps Index; picker Enter dispatches through `openRun` (error-path
equivalence with explicit `/open`); two-step `/answer` select→compose→submit (answer
persisted via `hitl`, compose state cleared); esc preserves the attached run/tab;
empty filter shows the empty state and Enter no-ops; answer cancel clears `answerQID`.

## Verification
```
go build ./...   # ok
go vet ./...     # ok
go test ./...    # ok (all packages)
```
(with `GOCACHE`/`GOMODCACHE` set to the repo-local caches.)

## Manual-smoke note
Not run in a live terminal as part of this slice (the model-driving tests exercise
the same `Update` paths headlessly, mirroring `live_test.go`). To smoke manually:
`parley tui` → type `/open` + Enter → arrow-select an idea/run + Enter; type `/answer`
+ Enter on a run with an open question → pick the question → type the answer + Enter.

## Deviations from FINAL
None of substance. Minor: the empty-state message distinguishes a filter that matched
nothing (`(no matches)`) from the never-reached no-candidate case (since a picker only
ever opens when candidates exist), which is clearer than always showing the
kind-specific text.

## Fix-up cycle 1 (Phase 8 — addressing review/round-01)
Applied the agreed fixes from codex + agy (hermes accepted as-is):
- **MAJOR (FINAL §8 async refresh)** — added `refreshPickerItems()` (rebuilds an open
  picker's `Items` from cached data, preserving `Filter`/cursor, re-clamping); wired
  into the `questionsMsg` handler (pickerAnswer) and `refreshHomeRuns` (pickerOpen).
  New test `TestPickerAnswerRefreshesOnBackgroundUpdate`.
- **MAJOR (answer lifecycle)** — `submitInput` no longer clears `composing`/`answerQID`
  before the write; `answerQuestion` clears them only on success, so a failed answer
  keeps the user in compose. New test `TestPickerAnswerFailureKeepsCompose`.
- **MAJOR (lose-on-error)** — `launchIdea` keeps `composing` + the typed task on failure.
- **MAJOR (hint)** — `renderInputRow` switches the hint to
  `↑/↓ select · type filter · Enter choose · esc cancel` while the picker is active.
- **MINOR (window math)** — `renderPicker` now reclamps a local copy of the picker and
  renders from it (single source of truth); `WindowSizeMsg` re-clamps the open picker.
- **NIT (label order)** — run rows now read `run  <run-id>  <idea>  [status]`; idea rows
  `idea  <slug>  [status]` (matches FINAL §3).
- **NIT (mixed receivers)** — dismissed: the value-receiver Update-path methods +
  pointer-receiver mutating helpers match the existing Bubble Tea idiom in this file
  (hermes concurred the receivers are appropriate).

`go build/vet/test ./...` green. Ready for Phase 8 re-review.
