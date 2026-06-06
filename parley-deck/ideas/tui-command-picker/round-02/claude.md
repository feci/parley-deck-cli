---
agent: claude
idea: tui-command-picker
round: 2
date: 2026-06-06
---

## Cross-review

Round-01 converged hard: all four of us want ONE reusable picker sub-mode on
`liveModel`, picker keys intercepted FIRST in `updateMain`, only `/open` + `/answer`
eligible, case-insensitive substring filter, "still press Enter on a single match",
and an empty-state message instead of a broken list. The remaining work is resolving
three concrete disagreements and locking the exact state shape. I'm changing my
round-01 position on two points below.

### @codex
Agreed and adopting your two strongest refinements:
1. **No callback in the picker struct.** My round-01 (and the 00-prompt) put
   `onSelect func(pickerItem)` on the struct. You're right — a func field is awkward
   for tests, model copies, and debugging. Adopt your `Kind` discriminator +
   `selectPickerItem(item)` method that switches on `Kind`. This keeps command effects
   in ordinary model methods and makes the state a plain value. **Consensus adopts this.**
2. **Key-routing order — the `if m.picker.Active { return m.updatePicker(msg) }` branch
   sits immediately after the `ctrl+c` case and before the `esc`/`N`/arrow/`enter`
   switch.** This is the single highest-risk detail; your placement is exactly right.
   Keep the redundant `&& !m.picker.Active` in the `N` guard as documentation of the
   invariant.
3. Adopt `answerQID string` on `liveModel`, reuse `composing` for `/answer` step 2,
   and your `submitInput`/`renderInputRow` branches verbatim in spirit. Your point
   that all cancel paths (`esc` from composing, `activateRun`, successful
   `launchIdea`/`answerQuestion`) must clear `composing=false, answerQID="",
   picker.Active=false, inputText=""` is a real correctness gate — I'm lifting it into
   FINAL as an explicit invariant with a test.

One counter on naming: you used both `pickerKind` (typed const) and a `Kind string` on
`pickerItem`. Let's keep `pickerState.Kind pickerKind` (typed) but drop the unused
`pickerItem.Kind string` field — the item only needs `Label` + `Value`; the dispatch
discriminator lives on the picker, not per-item. Less surface, nothing reads item.Kind.

### @agy
Three of your points are going into FINAL; one I'm overruling:
- **ADOPT — height cap + scroll.** Your "cap at 8 visible rows, scroll the viewport
  with ↑/↓ indicators" is the right answer to long lists and small terminals. FINAL
  will specify a visible window (cap ~8 rows) that scrolls to keep the selection in
  view, and a guard so the picker never exceeds available height on tiny terminals.
- **ADOPT — empty-state wording** (`(no recent runs or ideas to open)` /
  `(no open questions to answer)`) and the discoverability/help updates.
- **ADOPT — `/answer` clarity**: the answer-compose row must be unmistakable. But I'm
  taking codex/hermes's *mechanism* (reuse `composing` + `answerQID`, label
  `answer <qid> › `) rather than your auto-tab-switch + agent-tab focus. Auto-switching
  the active tab on selection is more moving parts (which tab? what if the question
  isn't agent-scoped?) and couples the picker to tab state. Reusing the existing
  compose row is simpler, testable without a terminal, and already where the user's
  eyes are. Your fallback ("pre-fill `/answer <qid> `") I'm also declining for codex's
  reason: it leaks command syntax and lets the user accidentally edit the qid.
- **OVERRULE — two-step `esc`.** You want esc to clear a non-empty filter first, then
  exit on a second press. codex, hermes and I want a single esc to cancel the picker
  outright. I'm going with single-esc-cancels: it matches this TUI's existing esc
  semantics everywhere else (esc = back out one level: clear input / detach), the lists
  are short so "lost my filter" is cheap to retype, and a mode where esc sometimes
  closes and sometimes doesn't is exactly the ambiguity hermes warns about. Backspace
  already edits the filter; esc always exits. Noted as a deliberate, recorded decision.

### @hermes
Agreed across the board — your keymap table matches codex's and mine. Two of your risk
items become required tests in FINAL:
- "open picker from an attached run → esc → still attached" (esc must restore exact
  prior state, not detach). This is the re-entrancy test.
- "`N` during filter types an N, never starts a new idea." Covered by the early picker
  branch; we test it explicitly.
On your open question (scroll keys during picker): agreed — **ignore** PgUp/PgDn/
ctrl+u/d/Home/End while the picker is active for v1 (no-op), rather than paging. Keeps
the model unambiguous; paging can come later. The only navigation keys the picker
consumes are ↑/↓ (and Enter/esc/printable/backspace).

## Resolved decisions (going into FINAL)

1. State: `picker pickerState{Active bool; Kind pickerKind; Title string;
   Items []pickerItem; Index int; Filter string}` + `answerQID string` on `liveModel`.
   `pickerItem{Label, Value string}`. No callbacks. No per-item Kind.
2. Routing: picker branch first (after ctrl+c). Picker consumes ↑/↓/Enter/esc/
   printable/backspace only; all other keys no-op while active.
3. Eligible: bare `/open` → ideas + recent runs; bare `/answer` → open questions.
   Explicit-arg forms unchanged. Empty candidates → `inputErr`, picker does not open.
4. `/answer` step 2: reuse `composing` + `answerQID`; label `answer <qid> › `;
   Enter runs `answerQuestion(qid, text)`. All cancel paths clear composing+answerQID+
   picker+inputText.
5. Filter: case-insensitive substring over `Label`+`Value`; printable appends, reset
   index to 0; backspace pops; single esc cancels the whole picker; single match still
   needs Enter; clamp index to filtered length on every mutation.
6. Render: list above the input row, ~8-row scrolling window, selection kept in view,
   empty-state message; hint switches to `↑/↓ select · type filter · Enter choose ·
   esc cancel`; `/help` + slash hints advertise bare `/open`/`/answer`.
7. Tests: codex's table-driven state-machine set + hermes's two re-entrancy tests,
   driving `Update(tea.KeyMsg{...})` with no terminal, asserting model fields.

## Risks (unchanged from consensus)
Key-routing order is the one that turns "looks done" into "still switches tabs" —
covered by an active-tab-invariance test under picker arrows. Stale `answerQID` answering
the wrong question — covered by clearing it on every cancel/run-change path.
