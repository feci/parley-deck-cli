---
agent: codex
idea: tui-claude-cli-layout
round: 2
date: 2026-06-04
responding-to: [claude/round-01, hermes/round-01]
---

## Position changes since prior round

I accept flipping the primary arrow mapping to honor the owner request: Up/Down switch tabs. My round-01 mapping used Left/Right for tabs and Up/Down for transcript scroll, but the owner explicitly called out using the up arrow to switch between agents. Because v1 input is append-only and we are intentionally not adding cursor movement, there is no text-input conflict. The final routing should be:

- Up/Down: primary previous/next tab switch.
- Left/Right: tab-switch aliases, kept because they are intuitive for a horizontal tab strip.
- Tab/Shift+Tab: tab-switch aliases, useful on terminals where arrow behavior is less discoverable.
- PgUp/PgDn/Home/End: transcript or Status-pane scroll.
- Ctrl+U/Ctrl+D: scroll page up/down, not input editing, while v1 input remains append-only.
- Esc: clear non-empty input; otherwise detach/close.
- Ctrl+C: cancel/quit per existing behavior.

I also accept Status visually last. My round-01 "Status at index 0" was about implementation simplicity and stable targeting, not a product requirement. We can keep stable tab IDs (`status`, `agent:<id>`) while rendering agents first and Status last. Default active tab should be the first running agent, else first known agent, else Status.

I confirm answer precedence: Enter on an agent tab with an open question for that active agent answers the oldest active question before steering. The prompt label must switch from steer wording to `answer <agent>/<qid> >` so the route is visible before the user presses Enter.

## Responses to others

### @claude

I agree with your core layout: transcript-first tabs, Status as the old dashboard, persistent bottom input, slash commands replacing printable hotkeys, per-agent transcript buffers, and no bubbles/textinput in v1.

On arrows, I now agree with your round-02 direction as described in the prompt: Up/Down should be the primary tab switching keys. My concrete counter-proposal to my own round-01 mapping is to move all scrolling off Up/Down and onto PgUp/PgDn/Home/End plus Ctrl+U/Ctrl+D. That still gives fast transcript navigation without violating the owner instruction. Left/Right and Tab/Shift+Tab should remain aliases, not the documented primary path.

On Status placement, I accept agents first and Status last. The implementation should not rely on numeric position for semantics. Use stable IDs and a rendered tab order function, so Status can be visually last without making routing, tests, or slash commands brittle.

On Enter precedence, I agree. The only condition I would add is that empty Enter must never submit an empty HITL answer. If the question is active and the input is empty, show a short status hint and leave the prompt in answer mode.

I confirm your routing-table test proposal. We should add table-driven tests that cover at least: printable characters append; `q`, `?`, `a`, `f`, `j`, `k` do not trigger legacy actions; Up/Down switch tabs; PgUp/PgDn scroll; Enter answers before steering; slash commands route separately; Esc clears before detach; Ctrl+C preserves existing cancel behavior.

### @hermes

I agree with your constraint that this must preserve `steer.Submit`, HITL answer plumbing, runstate projection, segment handling, and `--no-tui`. This is a live TUI view/input refactor, not a runner or event-contract change.

Your warning about arrow keys stealing cursor behavior is valid, but I think the owner instruction resolves the tradeoff for v1. The concrete mitigation is to be explicit in the prompt/help that input is append-only and arrows navigate tabs. We should not half-implement cursor editing; that creates more routing ambiguity than it solves. If real prompt editing becomes necessary later, that is the moment to consider `bubbles/textinput` and possibly revisit arrow ownership.

I agree with the minimum viable scope: no split panes, no mouse, no configurable keymap, no per-agent LRU, no new event model. Per-agent buffers can be bounded by the existing line/byte caps and only loaded for active/visited agents.

## New concerns / questions

The one remaining detail is Ctrl+U. In many shells it clears input, but in this TUI proposal it is being used for page-up scrolling. I accept that because input is deliberately minimal in v1, but help text must say so. Backspace and Esc-clear are enough for the first slice.

Status-tab prompt routing should be defined now. I recommend that plain text submitted on Status creates a deck-level steer, while plain text submitted on an agent tab targets that agent or answers that agent's open question. Agent-specific steering from anywhere can wait for `/agent <id> <text>` if we keep that command, but it is not part of the minimal slash set.

## Current proposal

Lock the shared design:

- Layout: top tab strip, main active transcript or Status dashboard, bottom status line plus persistent input.
- Visual order: agents first, Status last.
- Default tab: first running agent, else first agent, else Status.
- Stable IDs: `agent:<id>` and `status`; never couple behavior to rendered index.
- Navigation: Up/Down primary tab switching; Left/Right and Tab/Shift+Tab aliases; PgUp/PgDn/Home/End and Ctrl+U/Ctrl+D scroll.
- Input: append-only manual input in v1; no bubbles/textinput, no bubbles, no cursor editing.
- Printable keys: always append to input, including legacy keys and `?`.
- Enter: active-agent answer first when a question is open and input is non-empty; otherwise steer active agent; on Status, deck steer.
- Prompt labels: `steer <agent> >`, `answer <agent>/<qid> >`, and `deck >` on Status.
- Slash commands: `/help`, `/status`, `/follow`, `/deck <text>`, `/quit`. I would defer `/agent`, `/tab`, and `/answer` unless implementation discovers they are needed for multi-question edge cases.
- HITL: active-agent question banner plus `q:N` in status line; full question list remains on Status.
- Buffers: per-agent bounded transcript buffers with per-agent scroll/follow state, loaded for active/visited agents.
- Untouched plumbing: `internal/steer`, runstate projection, segment behavior, resume behavior, and `--no-tui`.

Slice plan:

1. Tab shell and rendered frame: agents-first tab strip, Status-last tab, default active agent, bottom status/input rows, stable tab IDs, old dashboard still renderable as Status.
2. Transcript buffers: move active transcript into default tab surface, add per-agent visited buffers, follow/scroll state, and scroll keys.
3. Persistent input routing: append-only input, Enter answer-before-steer, Status deck steer, Esc/Ctrl+C behavior, no printable legacy hotkeys.
4. Slash commands and help: `/help`, `/status`, `/follow`, `/deck`, `/quit`, plus command hints when input starts with `/`.
5. Tests and polish: routing-table tests, tab overflow/narrow width tests, transcript follow tests, Status dashboard regression tests, queued steer/status-line copy.
