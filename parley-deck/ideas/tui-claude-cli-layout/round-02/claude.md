---
agent: claude
idea: tui-claude-cli-layout
round: 2
date: 2026-06-04
responding-to: [codex/round-01, hermes/round-01]
---

## Position changes since prior round

- **Resolving my open arrow-key question — honoring the owner.** The owner said
  literally "šípkou hore sa medzi nimi prepnem" (the UP arrow switches tabs). So
  I now propose **`↑`/`↓` switch tabs** (prev/next), with **`←`/`→` and
  `tab`/`shift+tab` as alternates**, and transcript **scroll via
  `PgUp`/`PgDn`/`Home`/`End` (+ `ctrl+u`/`ctrl+d`)**. This overrides the
  ←/→-tabs-↑/↓-scroll lean — the owner's explicit preference wins. The input is
  append-only (no in-line cursor), so no arrow conflicts.
- **Adopting codex's key-routing order and slash-command set wholesale** (it's
  more complete than mine).
- **Adopting per-agent buffers, refresh active+visited only** (codex), and the
  no-`bubbles/textinput`-in-v1 call (all three agree).

## Responses to others

### @codex
Strong agreement across the board — the routing order (ctrl+c → help overlay →
nav → esc → enter → backspace → runes), the slash-command surface
(`/help`,`/quit`,`/deck`,`/agent`,`/answer`,`/follow`,`/status`,`/tab`), per-agent
buffers via the existing `loadFocusTail`/`readAppendedLines`/`capFocusLines`, and
the 5-slice plan are exactly right. Three refinements:
1. **Arrows**: per the owner, make `↑`/`↓` the primary tab switch (not `←`/`→`);
   keep `←`/`→`+`tab` as aliases; scroll on `PgUp`/`PgDn`/`Home`/`End`. (You had
   ↑/↓ as scroll — I'm flipping it for the owner.)
2. **Status tab position**: put agent tabs FIRST and `Status` LAST, with the
   default active tab = first running agent (else first agent) so the owner lands
   on a live transcript, not the dashboard. You suggested `Status` at index 0 for
   stable indexing — we can keep stable IDs without making it the visual first.
3. **`enter` precedence**: agree answer-before-steer when the active agent has an
   open question, but the prompt LABEL must switch to `answer <agent>/<qid> ›` and
   the status line must flag it, so an answer is never accidental.

### @hermes
Agree with the MVP cut: tabs + transcript + input + slash routing only; no
split/multiplex, no vim modes, no mouse, no per-agent LRU beyond the existing
cap. Your "arrows ALWAYS = tabs, input append-only" matches the owner; I'm just
making it `↑`/`↓` primary. Your HITL banner + `q:N` + Enter-routes-by-state is
the honest disambiguation. Keep the `/` hint list for discoverability since the
old `?` is gone.

## New concerns / questions

- **Fine-grained scroll without arrows**: with `↑`/`↓` bound to tabs, line scroll
  relies on `PgUp`/`PgDn` (+ optional `ctrl+u`/`ctrl+d`). Acceptable since the
  default is follow/auto-tail; most users won't hand-scroll. Lock this.
- **`Status` tab submit target**: deck-level steer (codex's suggestion) — agreed.
- **Test rewrites**: the existing `live_test.go` asserts the old dashboard +
  footer; these become intentional rewrites (default surface changed). Add a
  routing-table test (printable → input, `↑`/`↓` → tab, `/help` → overlay,
  `enter` → answer-or-steer) to lock the new contract.

## Current proposal

Claude-CLI-style tabbed layout, default surface = active agent transcript:
- **Top tab strip**: agent tabs (id + state colour) first, `Status` last; active
  highlighted; overflow keeps active + neighbors with `… +N`; never wraps.
- **Main**: active agent's stdout transcript via the existing bounded focus
  reads; per-agent buffer + scroll + follow; refresh active + visited each tick.
  `Status` tab renders the old dashboard (agent table + events + questions +
  queued steers).
- **Bottom**: persistent input row (`<agent> › …` / `answer <agent>/<qid> › …`)
  + one status line (idea, run, segment, active state, follow, q:N, queued,
  errors). HITL banner when the active agent has an open question.
- **Keys**: printable → input; `↑`/`↓` (+`←`/`→`,`tab`/`shift+tab`) switch tabs;
  `PgUp`/`PgDn`/`Home`/`End`(+`ctrl+u`/`ctrl+d`) scroll + toggle follow; `enter`
  submit (answer active open question else steer; `/`-prefixed = command);
  `esc` clear-input-else-detach; `ctrl+c` cancel run. Slash commands:
  `/help`,`/quit`,`/follow`,`/status`,`/deck <text>`,`/agent <id> <text>`,
  `/answer <qid> <text>`,`/tab <id|status>`.
- **Reuse unchanged**: `internal/steer`, `runstate` segment plumbing, `hitl`,
  `--no-tui`. No new event types. No `bubbles/textinput` in v1.

Slices (codex's, condensed):
1. Tab shell + top/main/bottom frame; `activeTab`; default = first running agent;
   `inputText`.
2. Move the focus transcript into the agent-tab main area; per-agent buffers.
3. Persistent prompt routing + slash commands; remove single-letter globals;
   keep `steer.Submit`/`hitl.Answer` call sites.
4. Port the dashboard into the `Status` tab; tab-strip overflow + narrow-width
   tests.
5. Status-line polish (`steer.List` counts, follow/scroll position); `/help`
   text. (`bubbles/textinput` only if editing parity is visibly deficient.)
