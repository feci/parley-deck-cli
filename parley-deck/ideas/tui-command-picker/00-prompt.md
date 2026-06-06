---
idea: tui-command-picker
author: user
created: 2026-06-06
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + Bubble Tea model/picker state
  codex: Go idioms, the picker state machine + key routing + tests
  agy: UX correctness + consistency across commands + edge cases
  hermes: keymap/interaction-model fidelity (no collisions with existing keys)
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem / idea

In the unified `parley tui` (internal/tui/live.go, shipped 1.16.0), arg-taking slash
commands force the user to TYPE a slug/id, which is annoying: e.g. `/open <slug|run>`
makes you remember and type a long timestamped slug. The owner wants: typing a command
that needs a target (most importantly `/open`) and confirming should pop an **arrow-key
selectable picker** (a list of the candidates), and the user just ↑/↓-selects and
Enter-confirms. The same pattern should apply to the other arg-taking commands.

## Current state (verified)

`liveModel` (live.go:102) already has a single-purpose input mode: `composing bool`
(the new-idea input opened by `N`). It has the data the pickers need: `homeRuns
[]runstate.RunSummary` (recent runs), `questions []hitl.Question` (open HITL
questions), and the ideas list (Status). Key routing is in `updateMain` (live.go:618);
slash commands are dispatched in `runCommand` (`/help /status /follow /deck /answer
/open /home /quit`); the bottom prompt renders in `renderInputRow` (live.go:538).
Today ↑/↓ switch tabs and Enter answers/steers; typing `/open foo` works but `/open`
alone errors with "usage: /open <slug|run-id>".

## Proposed direction (a STARTING proposal — challenge it in round-01)

Add a generic **picker mode** to `liveModel`, mirroring `composing`:
- A `picker` struct: `{active bool, title string, items []pickerItem, index int,
  filter string, onSelect func(pickerItem) (...)}` where `pickerItem` is
  `{label, value, kind}`.
- When an arg-taking command is run with NO argument (e.g. `/open`), populate the
  picker with the candidates and activate it instead of erroring.
- Key routing while the picker is active (in `updateMain`): **↑/↓ move the selection**
  (not the tabs), **Enter** confirms the highlighted item (runs the command with that
  value), **esc** cancels, **printable chars** type into a filter that narrows the
  list, **backspace** edits the filter.
- Render the picker as an overlay/list above the input row (reuse the help-overlay
  rendering style), showing the highlighted row.
- Commands covered: `/open` → ideas + recent runs; `/answer` → open questions
  (pick the question, then type the answer — two-step). Commands with no target
  (`/home /status /follow /help /quit`) are unchanged. `/deck`/steer take free text,
  not a pick.

## Round-01 focus questions (answer independently)

1. **The picker state shape** — a single reusable `picker` mode vs per-command modes?
   How does it compose with the existing `composing` mode and `inputText`?
2. **Key routing & collisions** — while the picker is active, ↑/↓ must move the
   selection, not switch tabs; Enter selects, esc cancels. Any collision with N,
   PgUp/PgDn, slash, or the steer/answer Enter? How to keep it unambiguous?
3. **Which commands get a picker, and the two-step `/answer`** — `/answer` needs a
   question id AND answer text; pick the qid then type the text. Is the two-step
   flow right, or pre-fill `inputText` with `/answer <qid> ` and let the user type?
4. **Type-to-filter** — should the picker filter as you type (fuzzy/substring)? It
   must not collide with the cancel/confirm/navigation keys.
5. **Data + empty states** — `/open` lists ideas + runs (dedup? ordering? labels with
   status?); empty list → message, not a broken picker. Re-fetch on open or use the
   cached `homeRuns`/`questions`?
6. **Discoverability** — does the input-row hint / `/help` need updating so users know
   a bare `/open` opens a picker?

## Constraints (non-negotiable)
- Pure TUI change in `internal/tui` (Bubble Tea / lipgloss). No new deps. No change to
  the runner/driver/agents. Keep `--no-tui` and the non-interactive paths untouched.
- Backward compatible: `/open <slug>` with an explicit arg still works directly (no
  picker); the picker is the no-arg affordance.
- The picker must not break existing keys (N, tab nav, scroll, steer/answer, slash).
- Testable with the Bubble Tea model (drive key msgs; assert selection/among items)
  WITHOUT a live terminal — mirror the existing live_test.go patterns.

## Non-goals
- No mouse support. No new commands beyond pickers for existing ones. No fuzzy-finder
  library. Not changing what the commands DO once a target is chosen.

## Deliverables (in order)
1. FINAL.md: the picker state machine (mode shape, key routing table, which commands,
   the two-step `/answer`, filter behavior, render).
2. Implement it in `internal/tui/live.go` (+ tests in `live_test.go`), build/vet/test
   green; a manual-smoke note. Then Phases 6-8 review.
