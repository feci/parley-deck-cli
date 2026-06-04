---
idea: tui-claude-cli-layout
status: implemented
implementer: claude
started: 2026-06-04
completed: 2026-06-04
branch: parley-deck-cli#feat/tui-claude-cli-layout
head-commit: see-branch-tip
design-pr: https://github.com/feci/parley-deck-cli/pull/34
implementation-pr: https://github.com/feci/parley-deck-cli/pull/34
---

## Summary of work

Reworked the `parley` live-run TUI default surface (`internal/tui/live.go`) into
the Claude-CLI-style tabbed transcript layout from FINAL.md. The engine,
`internal/steer`, `runstate` segment plumbing, `hitl`, resume, and `--no-tui`
are unchanged.

## Implementation plan / checklist

- [x] **Slice 1 — tab shell + frame**: `agentBuffer` type + `activeTab`/`inputText`/
      `buffers` model fields; `tabIDs`/`activeTabResolved` (agents first, Status
      last; stable `agent:<id>`/`status` ids; default = first running agent);
      `renderTabbed` (tab strip + main + status line + input row); `renderTabStrip`
      + `shortState`; `Status` tab = the old dashboard panes (`renderAgentTable`/
      `renderEventPane`/`renderQuestionsPane`).
- [x] **Slice 2 — transcript buffers**: per-agent `agentBuffer` over the existing
      `loadFocusTail`/`readAppendedLines`/`capFocusLines` (20k/4 MiB); `ensureBuffer`/
      `ensureActiveBuffer`/`refreshBuffers` (active+visited, reload on rotation);
      `renderTranscript`; follow + scroll state per agent.
- [x] **Slice 3 — input routing**: `updateMain` — printable→input; `↑/↓`(+`←/→`,
      `tab`/`shift+tab`) switch tabs; `PgUp/PgDn`/`ctrl+u`/`ctrl+d`/`Home`/`End`
      scroll (+ follow toggle); `Enter` → `submitInput` (answer active agent's open
      question else steer; deck steer on Status); `esc` clears-or-detaches;
      `ctrl+c` cancels; `backspace` deletes a rune. Legacy single-letter hotkeys
      removed.
- [x] **Slice 4 — slash commands + help**: `runCommand` for `/help`,`/status`,
      `/follow`,`/deck`,`/answer`,`/quit`; `/`-prefix shows a command-hint line;
      `modeHelp` overlay rewritten for the new keymap.
- [x] **Slice 5 — tests + polish**: status line (`run/idea/round/state/follow/q:N`),
      answer-mode colour-flipped input row + question banner; routing-table test
      and the rewritten behavior tests below.

- [x] Checks run: `go build ./...`, `go vet ./...`, `go test ./...` — all green.

## Tests

Rewrote `internal/tui/live_test.go` for the new default surface (the old
dashboard/footer/focus/compose/answer-mode tests were intentional rewrites, per
consensus D9): `TestTabbedDefaultShowsTranscriptAndInput`,
`TestTabSwitchWithUpDownArrows`, `TestTranscriptScrollAndFollow`,
`TestInputSteersActiveAgent`, `TestAnswerViaInputRowWhenQuestionOpen`,
`TestKeyRoutingPrintableAppendsNotHotkey` (D9 routing table),
`TestSlashHelpOpensOverlay`, `TestEscClearsInputThenDetaches`,
`TestInputBackspaceRemovesWholeRune`, `TestResumeStatusAndEscDetach`. Kept the
focus-read pipeline + projection + summarize tests unchanged.

## Deviations from FINAL.md

- **Tab-strip overflow** is a simple right-truncate of the joined strip for v1
  (FINAL described keeping active + neighbors with `… +N`). The strip never
  wraps; a far-right active tab on a very narrow terminal could be clipped. Noted
  as a polish follow-up.
- The retained dashboard renderer `renderQuestionsPane` still has an answer-mode
  branch (`mode == modeAnswerQuestion`) that is now never entered from the
  default surface; harmless, slated for cleanup.

## Notes for reviewers

- Removed 35 now-dead functions from the old modal layout (compact renderers,
  focus-mode handlers, composer/answer modal handlers, single-letter selectors).
  A few model fields (`answerText`/`answerErr`/`logPreview`/`focus*`/`compose*`)
  and the unused `liveMode` constants remain referenced only by the retained
  dashboard renderer; flag any you want fully removed.
- `internal/steer`, `runstate` (segment/`[FINISHED]`), `hitl`, resume, and
  `--no-tui` are untouched. No new event types; no new dependency.
- Per-agent buffers are bounded by the existing 20k-line / 4 MiB cap × visited
  agents; only active+visited buffers are refreshed each tick.
