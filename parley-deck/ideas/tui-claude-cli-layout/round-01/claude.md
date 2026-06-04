---
agent: claude
idea: tui-claude-cli-layout
round: 1
date: 2026-06-04
---

## Summary

Flip the TUI inside-out: the **default and primary view is the active agent's
live transcript**, with a **tab strip on top** (one tab per agent + a `Status`
tab that IS the old dashboard) and a **persistent bottom = input row + status
line**, like Claude Code. The hard part is not rendering — it's **key routing**:
with an always-typeable input, single-letter global keys (`j`,`g`,`a`,`?`,`q`)
no longer work. I propose resolving that with **slash-commands in the input**
(`/help`, `/status`, `/follow`, `/deck …`) plus a small set of non-printable
nav keys (arrows, Tab, PgUp/Dn, ctrl+c, esc). All 1.12.0 plumbing is reused;
this is presentation + input routing only.

## Proposed approach

### Layout (three regions)
```
┌ tabs ─────────────────────────────────────────────┐
│ [codex●] [claude○] [hermes◍]   Status   (▸ active) │   top: tab strip
├───────────────────────────────────────────────────┤
│ <active agent's live stdout transcript, follows>   │   main: transcript
│ ...                                                │   (the focus viewport,
│ ...                                                │    now full-screen)
├───────────────────────────────────────────────────┤
│ run=… idea=… segment=… codex:running follow:on q:1 │   status line
│ › type a steer for codex, or /help        [queued] │   input row
└───────────────────────────────────────────────────┘
```
- **Tab strip**: agent tabs (id + state-coloured dot) + a `Status` tab (the
  existing dashboard: agent table, events, questions). Active tab highlighted.
  Overflow on narrow terminals → truncate ids / show `‹ N more ›`.
- **Main**: the active agent's transcript, backed by the existing
  `loadFocusTail`/`readAppendedLines`/`capFocusLines` (bounded 20k/4 MiB),
  follow mode default-on. The `Status` tab renders the old two-pane dashboard.
- **Status line** (above input): `run`, `idea`, `round/segment`, active agent
  state + elapsed, `follow on/off`, open-question count, transient errors /
  "queued steer-NNNN".
- **Input row**: always present; placeholder hints the action. Plain text =
  steer the active agent; `/`-prefixed = command.

### State model (extend, don't fight, `liveMode`)
- Replace "overview is default" with **tabs**: `activeTab` ∈ {agent ids…,
  `status`}. The default landing tab is the first agent (transcript front-and-
  centre); `Status` is one tab over.
- The bottom input is **always live** (no separate compose mode). The prior
  `modeCompose`/`modeAnswerQuestion` collapse into the persistent input;
  `modeHelp` stays as an overlay (opened by `/help`).
- Keep `internal/steer` exactly as-is: Enter on non-empty input →
  `steer.Submit` for the active agent (or `/deck …` for the deck). Honest
  "recorded/queued" wording unchanged.

### Key routing (the crux) — Claude-CLI-like, no vim modes
- **Typing**: every printable char goes into the input.
- **Tabs**: `←`/`→` switch agent/Status tabs (and `tab`/`shift+tab` as
  always-works alternates). Arrows are non-printable so they never conflict with
  typing; I drop in-input left/right cursor (steers are short; `backspace`,
  `ctrl+u` clear suffice). [Open question: keep arrows for tabs always, or only
  when input is empty?]
- **Scroll transcript**: `PgUp`/`PgDn`, `ctrl+u`/`ctrl+d`; auto-follow when at
  bottom, PgUp drops follow; a `/follow` command (or `ctrl+f`) re-pins.
- **Enter**: if a HITL question is open for the active agent → submit the input
  as its answer; else if input non-empty → steer; else no-op.
- **esc**: clears the input if non-empty, else detaches the TUI (run continues).
- **ctrl+c**: cancel the run (unchanged).
- **Commands** (input starting with `/`): `/help`, `/status` (jump to Status
  tab), `/follow`, `/deck <text>` (steer the deck), `/quit`. A `/` opens a tiny
  command hint list. This replaces the old single-letter global keys and is the
  most Claude-CLI-authentic affordance.

### Per-agent transcript buffers
- Keep a buffer + offset + scroll + follow flag **per agent**, all advanced each
  250 ms tick via the existing incremental reads, so **switching tabs is instant
  and retains scroll position**. Memory is bounded by the existing cap (20k/4 MiB)
  × number of agents (2–4 → ≤ ~16 MiB worst case). If we ever support many
  agents, lazy-load + LRU-evict inactive buffers.

### HITL questions
- A one-line **banner** above the status line when a question is open for the
  active agent ("? Which branch? — type the answer below, Enter to submit"), and
  a `q:N` counter in the status line. Enter routes to answer when open. The
  `Status` tab keeps the full questions pane. `←/→` to the agent whose question
  is open.

### Slices
1. **Tabbed shell + transcript-as-default**: tab strip, `activeTab`, main =
   active agent transcript (reuse focus reads), `Status` tab = the old
   dashboard moved verbatim. Arrows/Tab switch tabs; PgUp/Dn scroll; follow.
   (Biggest visible win; no input yet — Enter/keys for nav only.)
2. **Persistent input row + status line**: always-on input wired to
   `steer.Submit` (active agent) and HITL answer; status line; banner.
3. **Slash commands + help**: `/help`,`/status`,`/follow`,`/deck`,`/quit`; `/`
   hint list; fold `modeHelp` overlay behind `/help`.

## Concerns / open questions

- **Arrows for tabs vs input cursor**: I lean "arrows always = tabs, input has no
  left/right cursor" (steers are short) — but codex/hermes may prefer the
  "empty-input heuristic" or `tab`/`shift+tab` only. This is THE decision.
- **Default landing tab**: first agent transcript (owner wants text front-and-
  centre) vs `Status`. I pick first-agent.
- **Detach key**: `esc`-when-empty vs a `/quit` command vs `ctrl+c` twice. Need
  one unambiguous detach that doesn't fight typing.
- Do we keep the compact (short-terminal) layout? The tabbed layout must degrade.
- Scope: this rewrites `View()` + a lot of `Update()`; risk of regressing the
  1.12.0 features. Lock current behavior with tests first.

## Risks

- **Key-conflict UX**: the single biggest risk — get routing wrong and it feels
  broken. Mitigate with slash-commands + a visible footer/hint and tests for the
  routing table.
- **Regression**: the segment fix, steer, help must survive the rewrite. Keep
  `internal/steer`/`runstate` untouched; add tests asserting steer + answer still
  work through the new input.
- **Memory**: N per-agent buffers — bounded, but verify the cap holds with many
  agents; lazy-evict if needed.
- **Terminal width**: tab strip + status line must truncate gracefully; test
  narrow widths.
- **Discoverability**: if actions hide behind `/commands`, users must find them —
  the input placeholder + `/` hint list must advertise them.
