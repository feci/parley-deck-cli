---
idea: tui-command-picker
phase: consensus
drafter: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
---

## Consensus

After round-01 (four independent positions) and round-02 (cross-review), all four
participants agree on the design below. No open disagreements remain. agy accepted both
facilitator overrules (single-esc-cancels; reuse `composing` for `/answer` step 2);
hermes's full key-collision sweep found no ambiguity; codex confirmed the state shape
and added cleanup-helper discipline.

### Goal
In the unified `parley tui` (`internal/tui/live.go`), a bare arg-taking slash command
(`/open`, `/answer`) opens an arrow-key **picker** instead of erroring, so the user
selects with ↑/↓ + Enter and never has to type a slug/run-id/qid. Explicit-arg forms
keep working directly. Pure TUI change; no new deps; runner/driver/agents untouched.

### 1. State shape (one reusable picker sub-mode)
On `liveModel`:
```go
type pickerKind string
const (
    pickerOpen   pickerKind = "open"
    pickerAnswer pickerKind = "answer"
)
type pickerItem struct{ Label, Value string }          // no per-item Kind, no callback
type pickerState struct {
    Active bool
    Kind   pickerKind
    Title  string
    Items  []pickerItem
    Index  int      // selection into the FILTERED view
    Filter string   // case-insensitive substring
    Offset int      // top of the visible scroll window
}
// liveModel gains:  picker pickerState   and   answerQID string
```
`composing` and `picker.Active` are mutually exclusive input sub-modes. No callback is
stored on the picker; command effects live in a `selectPickerItem(item)` method that
switches on `picker.Kind`. `pickerItem.Kind` is intentionally omitted (a visible type
like `idea`/`run` lives in `Label`).

### 2. Key routing (the highest-risk detail)
In `updateMain`, the picker branch is intercepted FIRST — immediately after the
`ctrl+c` case and BEFORE the existing `esc`/`N`/arrow/`enter`/printable handling:
```go
case "ctrl+c": ...          // global, unchanged
}
if m.picker.Active { return m.updatePicker(msg) }
switch msg.String() { /* existing esc / N / arrows / enter / printable */ }
```
`updatePicker` owns ONLY these keys; everything else is a no-op while active:

| Key | Picker active | Inactive (today) |
|---|---|---|
| ↑ / ↓ | move selection in filtered list | switch tab |
| Enter | confirm highlighted → dispatch by Kind | submit / steer / answer |
| esc | cancel the whole picker, restore prior state | clear input / detach |
| printable rune (incl. `N`, `/`) | append to Filter, reset Index→0 | append to inputText / dispatch |
| backspace / ctrl+h | pop last Filter rune | pop inputText |
| PgUp/PgDn/ctrl+u/ctrl+d/Home/End/shift+arrows | no-op (ignored) | scroll transcript |
| ctrl+c | global quit | global quit |

The `N` guard keeps a redundant `&& !m.picker.Active` as documentation of the invariant
"`N` opens the new-idea composer only when `!composing && inputText=="" && !picker.Active`".

### 3. Eligible commands + triggers
- Bare `/open` → picker over **ideas (from Status) + recent runs (`homeRuns`)**, ideas
  first then runs by recency, deduped by `Value`. Explicit `/open <slug|run-id>` stays
  direct. Empty candidates → set `inputErr` ("nothing to open yet"), do NOT open the picker.
- Bare `/answer` → picker over **open questions** (`questions` with status open).
  Explicit `/answer <qid> <text>` stays direct. Empty candidates → `inputErr`
  ("no open questions"), do NOT open the picker.
- No-target commands (`/home /status /follow /help /quit`) and free-text `/deck`
  are unchanged.

### 4. Two-step `/answer`
Selecting a question does NOT submit. It transitions to answer composition by reusing
the existing `composing` machinery with `answerQID` set:
`selectPickerItem` sets `answerQID=<qid>`, `composing=true`, `inputText=""`, clears the
picker. `renderInputRow` shows `answer <qid> › ` (warning style) when
`composing && answerQID != ""`. `submitInput` runs `answerQuestion(qid, text)` when
`answerQID != ""`, else `launchIdea(text)`. Not pre-filling `inputText` with
`/answer <qid> ` (avoids leaking command syntax / accidental qid edits).

### 5. State-reset invariant (correctness gate)
Two helpers enforce clean transitions:
- `clearPicker()` zeroes `picker.Active/Kind/Title/Items/Index/Filter/Offset`.
- `clearComposition()` sets `composing=false`, `answerQID=""`, `inputText=""`, and
  clears `inputErr` where appropriate.
Every path that exits or supersedes compose/picker — `esc` from composing, successful
`launchIdea`, successful `answerQuestion`, `activateRun`/open, returning Home — must
leave `composing=false, answerQID="", picker.Active=false, inputText=""`. This prevents
a later normal composition from submitting to a stale qid.

### 6. Filter
Case-insensitive substring over `Label`+`Value`. Printable appends and resets
`Index`→0; backspace pops; single `esc` cancels the entire picker (no two-step
filter-clear). A single filtered match still requires Enter (no auto-select). `Index`
is clamped to the filtered length on every mutation and again before render; the visible
`Offset` window is computed from the clamped index.

### 7. Render
A list above the input row (reuse the existing overlay/box style), inserted before the
status line. Visible rows capped at ~8 AND dynamically clamped to the available height
(`min(8, transcriptHeight())`, floor ~3) so tiny terminals don't overflow; the window
scrolls to keep the selection in view, with `↑`/`↓` boundary indicators when the list
overflows. Empty filtered results show a muted message
(`(no recent runs or ideas to open)` / `(no open questions to answer)`). While the
picker is active, the tab strip is dimmed (`mutedStyle`) to signal focus is on the
picker. The input-row hint switches to `↑/↓ select · type filter · Enter choose · esc
cancel`; `/help` and the slash-command hint advertise that bare `/open` / `/answer`
open a picker.

### 8. Data freshness
Use cached `homeRuns` / `questions` (already refreshed by ticks/messages) — no new
fetch on open. If a background `eventsMsg`/`questionsMsg` arrives while the picker is
open, rebuild `picker.Items` reactively WITHOUT resetting `picker.Filter` or destroying
the cursor (re-clamp Index after rebuild). Async refresh-on-open is a deferred follow-up.

### 9. Tests (state-machine first, no terminal)
Drive `Update(tea.KeyMsg{...})` and assert model fields (mirror `live_test.go`):
1. Bare `/open` with candidates activates `pickerOpen`; explicit `/open <target>` does NOT.
2. Bare `/answer` with open questions activates `pickerAnswer`; selecting sets
   `composing=true`, `answerQID=<qid>`, `inputText=""`.
3. Picker ↑/↓ changes `picker.Index` and does NOT change the active tab.
4. Picker printable (incl. `N` and `/`) changes `picker.Filter`, not `inputText`, and
   starts no command.
5. Filter narrowing clamps Index when it exceeds the filtered length.
6. Picker Enter for `/open` uses the same open path as explicit `/open <value>`.
7. Picker `esc` cancels only the picker and preserves attached run / tab (re-entrancy).
8. Empty filtered results show the empty state; Enter on empty selects nothing and
   mutates no run/answer state.
9. Answer submit calls `answerQuestion`, then clears `answerQID`+`composing`+`inputText`
   (and the answer persists via `hitl`).
10. Answer cancel clears `answerQID` so a later `N` composition can't answer the old qid.
11. Empty candidates set `inputErr` and do NOT open the picker.

### 10. Risks (mitigated)
- Picker branch placed too late → arrows still switch tabs / `N` starts an idea / Enter
  submits wrong. Mitigated by branch-first placement + active-tab-invariance test.
- Stale `answerQID` answering the wrong question. Mitigated by `clearComposition()` on
  every cancel/run-change path + tests 9-10.
- Scroll-window off-by-one when the filtered list shrinks. Mitigated by clamping Index
  before selection and before render, and computing Offset from the clamped index.

### Non-goals
No mouse. No fuzzy-finder library. No new commands. No change to what a command does
once a target is chosen. `--no-tui` / non-interactive paths untouched.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-06. As drafter I accept this consensus. It reflects the round-02 resolution:
one reusable picker sub-mode, picker-branch-first routing, `/open`+`/answer` only,
reuse-`composing`+`answerQID` for the two-step answer, single-esc-cancels, substring
filter with Enter-to-confirm, ~8-row dynamically-clamped scroll window, and the
`clearPicker()`/`clearComposition()` state-reset invariant with state-machine tests. No
blockers. Ready for FINAL + implementation.

### codex — ✅ ACCEPT
I ACCEPT this consensus and have no blocking objections. Please preserve the picker-branch-first routing and the `clearPicker()`/`clearComposition()` cleanup discipline during Phase 5, because those are the parts most likely to prevent stale answer state or accidental tab/input handling. Also make the render path clamp `Index` before computing the visible window so filtered-list shrinkage cannot produce an invalid selection.

### agy — ✅ ACCEPT
I ACCEPT this consensus design. I am particularly pleased to see the adoption of the dynamic ~8-row viewport clamp, which prevents UI overflow on tiny terminals. Preserving the filter and cursor during background tick updates, dimming the tab strip, and updating the input-row hint and `/help` ensure a robust, discoverable, and polished user experience.

### hermes — ✅ ACCEPT
I accept the consensus as drafted. The picker-branch-first routing immediately after ctrl+c, combined with single-esc-cancels and the explicit no-op list for all scroll/paging keys while the picker is active, guarantees zero key collisions. N-during-filter typing text and the two required re-entrancy and filter tests complete the invariant set.
