---
agent: claude
idea: tui-layout-refresh
round: 1
date: 2026-05-26
---

## Summary

The TUI stacks five full-width regions vertically (header, 3-col panel row, actions/questions, logs, footer), consuming far more rows than a typical 24–40 line terminal can display. Each panel renders multi-line detail inline, so scanning status requires scrolling. The live view repeats the pattern: agent table + events side-by-side, then full-width questions, full-width logs, footer — four vertical sections minimum.

The core information-architecture problem is that every piece of data is shown at the same visual weight at the same time, regardless of how often it changes or how urgent it is.

## Proposed approach

### 1. Two-tier layout: status strip + detail pane

Replace the 3-column top row with a single **status strip** (1–2 lines per run) that uses inline badges for attention, agent progress, and open-question count. Below it, a single **detail pane** shows the content of whichever entity is selected (run details, agent detail, event stream, or question). This converts the layout from "everything visible, nothing readable" to "summary always visible, detail on demand."

Vertical cost: header (1) + status strip (N runs × 2 lines, capped at ~8) + detail pane (remaining height) + footer (1–2). Fits in 24 rows with room for a 10-line detail area.

### 2. Compact mode triggered by terminal height

When `m.height < 30` (or a similar threshold), collapse the status strip to one line per run (badge + slug only, no run-ID or timestamp) and hide the footer help text (show on `?` keypress). This avoids an unusable layout on short terminals without a separate code path — just conditional truncation in the existing `View()`.

### 3. Status strip design (dashboard view)

Each run line:
```
[WAIT]  session-resume-cache  ██░░ 2/4  !2   3m ago
```
- `[WAIT]` — colored attention badge (green/yellow/red via Lip Gloss)
- `██░░` — agent progress as a micro-bar (done/total, 4 chars wide)
- `!2` — open question count, highlighted if > 0
- `3m ago` — last-event age, right-aligned

This compresses the current 4-line-per-run `renderIdeas()` into one dense scannable line.

### 4. Status strip design (live view)

Agent table compresses to one line per agent with inline state badge:
```
[RUN] codex   12m  artifact-written
[WAIT] gemini  8m   waiting-for-answer
```
Move the event stream into the detail pane (selected agent's events only), eliminating the side-by-side split that wastes width when the event pane is sparse.

### 5. Detail pane is context-sensitive

- **Run selected:** show run metadata + agent sub-list (current `renderRunDetails` content, reformatted)
- **Agent selected:** show agent detail + log tail (merge `renderAgentDetails` and `renderLogPane`)
- **Question selected:** show question detail + answer input

This avoids rendering three panels when only one is being read.

### 6. Footer as a single adaptive line

Current footer has two long lines. Compress to one line of the most relevant keys for the current focus, e.g.:
- Focus on runs: `j/k sel  enter detail  N new  r refresh  q quit`
- Focus on agents: `j/k sel  h/i/m mode  x clear  tab back  q quit`

Show full key reference only on `?`.

## Concerns / open questions

1. **Discoverability vs. density.** Hiding panels behind selection means new users won't see all data at once. Mitigate with clear selection indicators and the `?` help overlay, but this is a real trade-off.

2. **Three-column removal.** The current `JoinHorizontal` layout is the expected "dashboard" look. Removing it is a UX regression for wide terminals (>160 cols) where all three columns fit comfortably. Consider keeping the 3-column layout as a `wide mode` when `width > 140 && height > 40`, falling back to the strip+detail layout otherwise.

3. **Micro-bar for agent progress.** The `██░░` notation is compact but may be unfamiliar. Alternative: `2/4` plain text with color on the denominator. Needs input from the color/state hierarchy lens (Gemini).

4. **Height detection.** Bubble Tea provides `WindowSizeMsg` for width but the current code only tracks `m.width`. We need to also track `m.height` and thread it through the view logic. Minor code change but worth flagging.

5. **Event stream depth.** If events move into the detail pane, how many lines should be shown? Suggest: fill remaining height after metadata lines, minimum 5 lines. This requires computing available height dynamically rather than using fixed counts.

## Risks

- **Scope creep.** Adaptive layout (compact mode + wide mode + default) could balloon into a broad rewrite. Recommend shipping strip+detail as the single default first, then adding wide-mode in a follow-up if users on large terminals complain.
- **Regression in HITL flow.** The questions/answer panel is the most operationally critical part. If it moves into the detail pane, the user must navigate to it explicitly. Mitigate by auto-switching detail pane to the question view when a new open question arrives (attention-driven focus).
- **Testing.** Layout logic in `View()` is currently untested. Refactoring it adds risk. Suggest adding a few golden-file snapshot tests (`View()` output at 80x24, 120x40, 200x50) before changing anything, so regressions are caught.
