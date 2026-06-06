---
agent: claude
idea: tui-command-picker
round: 1
date: 2026-06-06
---

## Summary

This is a contained Bubble Tea state addition: one reusable **picker mode** on
`liveModel`, mirroring the existing `composing` mode, activated when an arg-taking
slash command is run with no argument. The model already holds the data
(`homeRuns`, `questions`, the ideas list), so the work is the mode struct, the key
routing while it's active, the overlay render, and wiring `/open` + `/answer` to open
it. Keep it backward-compatible (explicit `/open <slug>` still works directly).

## Proposed approach

**1. One reusable picker mode (not per-command).**
```go
type pickerItem struct{ label, value, detail string }
type picker struct {
    active bool
    title  string
    kind   string        // "open" | "answer" — what Enter dispatches to
    items  []pickerItem  // full list
    filter string        // type-to-narrow (case-insensitive substring)
    index  int           // selection into the FILTERED view
}
```
`liveModel` gains a `pk picker`. `composing` and `pk.active` are mutually exclusive
input modes; `inputText` is untouched while a picker is open. A `filtered()` helper
returns the visible items; `index` is clamped to it.

**2. Key routing (updateMain), guarded by `m.pk.active` FIRST.**
| Key | Picker active | Else (today) |
|---|---|---|
| ↑ / ↓ | move selection in filtered list | switch tab |
| Enter | confirm highlighted → dispatch by `kind` | answer/steer |
| esc | cancel picker (close, clear filter) | clear input / detach |
| printable | append to `filter`, reset index | type into input |
| backspace | edit `filter` | edit input |
| (tab/PgUp/PgDn/N) | ignored while picking | their normal action |
No collision: `N` only fires today when `!composing && inputText==""`; add
`&& !pk.active`. The picker owns ↑/↓/Enter/esc only while active.

**3. Triggers + the two-step `/answer`.**
- `runCommand("/open", "")` → `openPicker("open", ideasAndRuns())`; `/open <slug>` with
  an arg keeps the direct path.
- `/answer` with no qid → `openPicker("answer", openQuestions())`. Selecting a question
  does NOT finish — it transitions to an **answer-compose** input pre-tagged with the
  chosen qid (reuse the `composing`-style input, label `answer <qid> › `), so the user
  types the answer text + Enter → runs `/answer <qid> <text>`. Two-step because
  `/answer` needs id AND free text. `/answer <qid> <text>` typed in full still works.
- No-target commands (`/home /status /follow /help /quit`) and free-text `/deck` are
  unchanged.

**4. Type-to-filter:** case-insensitive substring over `label`+`detail`; printable
chars append to `filter` and reset `index` to 0; backspace pops. esc cancels the whole
picker (it does not just clear the filter — one consistent exit).

**5. Data, dedup, empty states.** `/open` list = the ideas (from Status) + recent runs
(`homeRuns`), each a `pickerItem` labelled with kind + status, e.g.
`idea  <slug>            [final]` / `run   <run-id>  <idea>  [running]`. Dedup a run
that's just the latest run of a listed idea is optional; ordering = ideas first, then
runs by recency. Use the cached `homeRuns`/`questions` (already refreshed on Home /
each tick) — no new fetch needed. An empty list opens the picker showing
`(nothing to open — press esc)` rather than erroring.

**6. Discoverability.** Update the input-row hint and `/help` so a bare `/open`/
`/answer` advertises the picker (e.g. `/open  pick an idea/run`).

## Concerns / open questions

- **Render placement** — the picker overlay should sit above the input row (like the
  help overlay) and cap height (scroll the window if the list is long). Reuse
  `boxStyle`/the help-overlay renderer; codex/hermes confirm the cleanest reuse.
- **`/answer` selection→compose transition** — does selecting a question switch into
  the existing `composing` machinery (retag it) or a third small mode? I lean: reuse
  `composing` with an `answerQID` field set, so the Enter handler knows to run
  `/answer <qid> <text>` instead of launching an idea.
- **Filter vs immediate single match** — if the filter narrows to one item, still
  require Enter (no auto-select) to stay predictable.
- **Where ideas come from** — confirm the model field holding the ideas list for the
  no-run/Home case vs an active run; the picker must work from Home AND mid-run.

## Risks

- **Key-routing regressions** — the picker guard must be checked before the existing
  ↑/↓/Enter/esc handlers; a missed guard breaks tab nav or steering. Mitigate with
  table-driven model tests (drive key msgs, assert mode + selection) like live_test.go.
- **Scope creep into a fuzzy-finder** — keep it substring + ↑/↓; no library, no mouse.
- **Stale candidates** — using cached lists could show a just-finished run oddly;
  acceptable (esc + reopen), but note it.
- **Two-step `/answer` confusion** — clear labels (`answer <qid> › `) + a hint so the
  user knows step 2 is "type the answer".
