---
agent: hermes
idea: tui-command-picker
round: 1
date: 2026-06-06
---
## Summary
A modal picker sub-state (single reusable `picker` struct on liveModel) is the correct model. It mirrors the existing `composing` bool + inputText pattern and must be mutually exclusive with it. While active it fully owns ↑/↓/Enter/esc/printable/backspace; all other keys (tab nav, N, scrolls, ctrl+c) keep their global meaning or are explicitly ignored to avoid ambiguity.

## Proposed approach
Introduce `picker active bool, title string, items []pickerItem, index int, filter string, onSelect func...` exactly as sketched. Activate only on bare arg-taking commands (/open, /answer). Key table while picker open:
- ↑/↓: move selection (not tabs)
- Enter: confirm (call onSelect, deactivate)
- esc: cancel (deactivate, restore prior tab/run focus)
- printable (except N handled specially): append to filter (N must type, never trigger new-idea)
- backspace: edit filter
- / : type into filter (never dispatch command)
- PgUp/PgDn/ctrl+u/d/Home/End/shift+↑↓: ignored (or optionally page the filtered list; safer to ignore first)
- ctrl+c: global cancel/quit (unchanged)
- all other keys: no-op

Re-entrancy: esc always restores the exact prior state (activeTab, attached run, inputText cleared). Opening from Home vs mid-run is identical because picker is orthogonal to tabs. Render as overlay (reuse help style) above input row; empty list shows "no matches".

## Concerns / open questions
- Two-step /answer: picker then pre-fill `/answer <qid> ` into inputText feels cleaner than nested state; confirm with UX lens.
- Filter: substring is sufficient; fuzzy adds complexity without library.
- Discoverability: update /help text only; no other changes.

## Risks
- N collision during filter: handled by special-case check before append.
- esc re-entrancy leak: must unit-test "open picker from attached run → esc → still attached".
- Scroll keys during picker: ignoring them is safest to keep model unambiguous; page-list later if needed.