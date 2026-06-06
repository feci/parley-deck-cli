---
idea: tui-command-picker
phase: final
status: final
drafter: claude
implementer: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
supersedes: consensus.md
---

# FINAL — TUI command picker

Ratified by claude, codex, agy, hermes (all ✅ ACCEPT in consensus.md). This is the
build spec for Phase 5. Scope: `internal/tui/live.go` (+ `live_test.go`). No new deps;
runner/driver/agents/`--no-tui` paths untouched.

## Goal
A bare arg-taking slash command (`/open`, `/answer`) in the unified `parley tui` opens
an arrow-key **picker** instead of erroring. The user selects with ↑/↓ + Enter and never
types a slug / run-id / qid. Explicit-arg forms keep working directly (backward compatible).

## 1. State (one reusable picker sub-mode)
Add near `liveModel`:
```go
type pickerKind string
const (
    pickerOpen   pickerKind = "open"
    pickerAnswer pickerKind = "answer"
)
type pickerItem struct{ Label, Value string }   // no per-item Kind, no callback
type pickerState struct {
    Active bool
    Kind   pickerKind
    Title  string
    Items  []pickerItem
    Index  int      // selection into the FILTERED view
    Filter string   // case-insensitive substring
    Offset int      // top of the visible scroll window
}
```
`liveModel` gains `picker pickerState` and `answerQID string`. `composing` and
`picker.Active` are mutually exclusive input sub-modes. Command effects live in
`selectPickerItem(item pickerItem)` which switches on `picker.Kind` (no func field).

## 2. Key routing (highest-risk detail — get the order right)
In `updateMain`, intercept the picker FIRST — right after the `ctrl+c` case, BEFORE the
existing `esc`/`N`/arrow/`enter`/printable switch:
```go
// ... ctrl+c handled and returns above ...
if m.picker.Active {
    return m.updatePicker(msg)
}
// existing switch (esc, N, arrows, enter, printable, slash) unchanged below
```
`updatePicker(msg)` owns ONLY: ↑/↓ (move Index), Enter (confirm filtered[Index] →
`selectPickerItem`), esc (`clearPicker()`), printable runes (append to Filter, reset
Index→0), backspace/ctrl+h (pop Filter). Every other key is a no-op while active
(PgUp/PgDn/ctrl+u/ctrl+d/Home/End/shift+arrows explicitly ignored). Clamp Index after
every mutation.

The `N` guard keeps the redundant `&& !m.picker.Active` term as documentation:
`N` opens the new-idea composer only when `!composing && inputText=="" && !picker.Active`.

## 3. Eligible commands
- Bare `/open` → items = ideas (from `opts.Status`) then recent runs (`homeRuns`),
  deduped by `Value`. Explicit `/open <slug|run-id>` unchanged. **Empty candidates →
  set `inputErr` ("nothing to open yet") and do NOT open the picker.**
- Bare `/answer` → items = open questions (`questions`, open status). Explicit
  `/answer <qid> <text>` unchanged. **Empty candidates → `inputErr` ("no open
  questions") and do NOT open the picker.**
- `/home /status /follow /help /quit` and free-text `/deck` are unchanged.

`pickerItem` labels carry the visible type, e.g. `idea   tui-command-picker   [final]`
and `run    20260606-…   tui-command-picker   [running]`. `Value` = slug (ideas) / run-id
(runs) / qid (questions).

## 4. Two-step `/answer`
`selectPickerItem` for `pickerAnswer` sets `answerQID=<qid>`, `composing=true`,
`inputText=""`, then `clearPicker()`. `renderInputRow` shows `answer <qid> › ` (warning
style) when `composing && answerQID != ""`. `submitInput` runs `answerQuestion(qid,text)`
when `answerQID != ""`, else `launchIdea(text)`. Do NOT pre-fill `inputText` with
`/answer <qid> `.

## 5. State-reset invariant (correctness gate)
- `clearPicker()` zeroes `Active/Kind/Title/Items/Index/Filter/Offset`.
- `clearComposition()` sets `composing=false`, `answerQID=""`, `inputText=""` (and clears
  `inputErr` where appropriate).
Every path that exits/supersedes compose or picker — `esc` from composing, successful
`launchIdea`, successful `answerQuestion`, run activation/open, return Home — leaves
`composing=false, answerQID="", picker.Active=false, inputText=""`. Prevents a later
normal composition from submitting to a stale qid.

## 6. Filter
Case-insensitive substring over `Label`+`Value`. Printable appends + resets Index→0;
backspace pops; single `esc` cancels the whole picker (no two-step). Single filtered
match still requires Enter (no auto-select). Clamp Index to filtered length on every
mutation and again before render; compute `Offset` from the clamped Index.

## 7. Render
List above the input row, inserted before the status line (reuse existing box/overlay
style). Visible rows = `min(8, transcriptHeight())` with a floor (~3) so tiny terminals
don't overflow; window scrolls to keep selection in view; show `↑`/`↓` boundary
indicators when the list overflows. Empty filtered results → muted message
(`(no recent runs or ideas to open)` / `(no open questions to answer)`). While active,
dim the tab strip (`mutedStyle`). Input-row hint switches to
`↑/↓ select · type filter · Enter choose · esc cancel`. `/help` + the slash hint
advertise that bare `/open` / `/answer` open a picker.

## 8. Data freshness
Use cached `homeRuns` / `questions`; no fetch on open. If a background
`eventsMsg`/`questionsMsg` arrives while the picker is open, rebuild `picker.Items`
WITHOUT resetting `Filter` or the cursor (re-clamp Index after rebuild). Async
refresh-on-open is a deferred follow-up.

## 9. Tests (state-machine first, no terminal; mirror `live_test.go`)
1. Bare `/open` (candidates present) activates `pickerOpen`; explicit `/open <target>` does NOT.
2. Bare `/answer` (open questions present) activates `pickerAnswer`; select → `composing=true`, `answerQID=<qid>`, `inputText=""`.
3. Picker ↑/↓ changes `picker.Index`, NOT the active tab.
4. Picker printable (incl. `N`, `/`) changes `picker.Filter`, not `inputText`; starts no command.
5. Filter narrowing clamps Index when it exceeds filtered length.
6. Picker Enter for `/open` uses the same open path as explicit `/open <value>`.
7. Picker `esc` cancels only the picker; attached run / tab preserved (re-entrancy).
8. Empty filtered results show the empty state; Enter selects nothing, mutates no state.
9. Answer submit calls `answerQuestion`, then clears `answerQID`+`composing`+`inputText` (persisted via `hitl`).
10. Answer cancel clears `answerQID` so a later `N` composition can't answer the old qid.
11. Empty candidates set `inputErr` and do NOT open the picker.

## 10. Acceptance
`go build ./... && go vet ./... && go test ./...` green (with
`GOCACHE`/`GOMODCACHE` set to the repo-local caches). Manual-smoke note in
`IMPLEMENTATION.md`. Then Phases 6-8 review by codex/agy/hermes to zero agreed fixes.

## Non-goals
No mouse, no fuzzy-finder library, no new commands, no change to what a command does
once a target is chosen.
