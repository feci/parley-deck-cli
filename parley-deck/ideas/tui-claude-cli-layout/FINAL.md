---
idea: tui-claude-cli-layout
status: final
author: claude
consensus-date: 2026-06-04
participants: [claude, codex, hermes]
---

## Final plan / specification

Rework the `parley` live-run TUI default surface into a **Claude-CLI-style
tabbed, transcript-centric layout**. Unanimous consensus (claude/codex/hermes;
see consensus.md). This is a default-layout/UX refactor of `internal/tui/live.go`
only — the engine, `events.jsonl`, `internal/steer`, `runstate` segment plumbing,
`hitl`, resume, and `--no-tui` are reused UNCHANGED. No new event types, no new
dependency (`bubbles`/`bubbles/textinput`) in v1.

### Layout (three regions)
```
[codex ●RUN] [claude ○FIN] [hermes ◍PEN]            Status      ← top tab strip
─────────────────────────────────────────────────────────────
<active agent's live stdout transcript — follows by default>   ← main
...
─────────────────────────────────────────────────────────────
run … idea … segment … codex:running follow:on q:1 queued:2    ← status line
steer codex › <input>                                          ← input row
```
- **Top:** one-line tab strip — agent tabs first (id + state-coloured marker),
  `Status` tab last. Active tab highlighted. Overflow: keep active + neighbors,
  `… +N` on the clipped side, never wrap.
- **Main:** active tab content — an agent's stdout transcript (via the existing
  `loadFocusTail`/`readAppendedLines`/`capFocusLines`, bounded 20k lines/4 MiB,
  follow default-on), or the `Status` dashboard (old agent table + events +
  questions + queued steers).
- **Bottom:** a status line + an always-present prompt input row.

### Tabs & state
- Stable IDs `agent:<id>` and `status`; behavior/tests never keyed to rendered
  index. Tab order = `runstate.ProjectEvents` agent order, then `Status`.
- Default active tab = first running agent, else first known agent, else
  `Status`.
- Per-agent bounded transcript buffer + scroll + follow state; load active +
  visited agents only; refresh those each 250 ms tick.

### Key routing (append-only input; no vim modes)
- Printable chars (incl. legacy `q`/`?`/`a`/`f`/`i`/`j`/`k`) → append to input.
- `↑`/`↓` = primary tab switch (owner request); `←`/`→`, `tab`/`shift+tab` aliases.
- `PgUp`/`PgDn`/`Home`/`End` (+ `ctrl+u`/`ctrl+d`) = scroll active pane;
  scroll-up drops follow, bottom re-enables it.
- `enter` = submit: if active agent has an open HITL question AND input non-empty
  → answer the oldest open question (`hitl.Answer`); else steer active agent
  (`steer.Submit`); on `Status` → deck steer. Empty Enter never submits an empty
  answer.
- `esc` = clear non-empty input, else detach/close. `ctrl+c` = cancel run
  (unchanged). `backspace` = delete one rune.
- Slash commands (leading `/`): `/help`, `/status`, `/follow`, `/deck <text>`,
  `/answer <qid> <text>`, `/quit`; a `/` pops a one-line command-hint list.
  Defer `/agent`, `/tab`.

### Input/HITL hardening
- Prompt labels show the route before Enter: `steer <agent> ›`,
  `answer <agent>/<qid> ›`, `deck ›`. In answer mode the **whole prompt row is
  colour-flipped** (warn) + a one-line banner; `q:N` in the status line; `/answer`
  is the explicit escape hatch. Answers must never be accidental.

### Reuse / retire
- Unchanged: `internal/steer`, `runstate` (segment/`[FINISHED]` fix), `hitl`,
  resume, `--no-tui`. Keep `modeHelp` overlay (now opened by `/help`); retire
  `modeOverview`/`modeAgentDetail`/`modeCompose`/`modeAnswerQuestion` as
  top-level states (folded into the tabbed view + persistent input). The old
  two-pane dashboard becomes the `Status` tab content.

### Slices
1. Tab shell + frame (strip, `activeTab` stable IDs, default = first running
   agent, bottom rows; old dashboard renderable as `Status`).
2. Transcript buffers (active transcript in the main area; per-agent visited
   buffers + follow/scroll + scroll keys).
3. Input routing (append-only input; Enter answer-before-steer / deck on Status;
   `esc`/`ctrl+c`; remove printable legacy hotkeys).
4. Slash commands + help (`/help`,`/status`,`/follow`,`/deck`,`/answer`,`/quit`
   + `/` hint list; `modeHelp` behind `/help`).
5. Tests + polish (routing-table, tab overflow/narrow-width, follow, Status
   regression; status-line copy via `steer.List`; answer-mode colour).

### Tests-first
Lock `steer.Submit` + `hitl.Answer` + segment projection with tests before/while
rewriting input routing; add the routing-table test (D9). Existing
dashboard/footer tests are intentionally rewritten (default surface changed).

### Invariants
- `--no-tui` + `events.jsonl`/artifact + runner contracts unchanged.
- No new event types; no new dependency in v1.
- No regression of 1.12.0 features (segment fix, steer queue + `parley steer`,
  help).

### Deferred
`bubbles/textinput` editing; `/agent`,`/tab`; `Events` tab; mouse; configurable
keymap; per-agent LRU; (carried) executing queued steers, live ACP delivery,
opt-in thoughts.

### References
- Consensus: ./consensus.md (claude/codex/hermes ACCEPT, 2026-06-04)
- Rounds: ./round-01/, ./round-02/
- Builds on shipped idea `tui-interactivity-overhaul` (1.12.0).
