---
idea: tui-claude-cli-layout
drafted-by: claude
date: 2026-06-04
participants: [claude, codex, hermes]
status: accepted
---

# Consensus — Claude-CLI-style tabbed transcript layout

Round-01 (independent) + round-02 (cross-review) converged with no blockers. This
is the basis for FINAL.md. It is a default-layout/UX refactor of
`internal/tui/live.go`; the engine, events, `internal/steer`, `runstate` segment
plumbing, `hitl`, resume, and `--no-tui` are reused UNCHANGED.

## Agreed decisions

### D1 — Three-region layout, transcript-first
- **Top:** a one-line tab strip. **Main:** the active tab's content — an agent's
  live stdout transcript, or the `Status` dashboard. **Bottom:** a status line
  plus an always-present prompt input row (Claude-CLI style).

### D2 — Tabs
- One tab per agent + a `Status` tab. **Agents render first, `Status` last.**
- **Stable IDs** (`agent:<id>`, `status`); never couple behavior/tests/commands
  to the rendered numeric index.
- **Default active tab** = first running agent, else first known agent, else
  `Status` (the owner lands on a live transcript, not the dashboard).
- Tab order follows `runstate.ProjectEvents` agent order. Overflow on narrow
  terminals: keep the active tab + neighbors, show `… +N` on the clipped side;
  never wrap the strip (wrapping steals transcript rows).

### D3 — Navigation keys (honoring the owner)
- **`↑`/`↓` = primary previous/next tab switch** (owner asked for the up arrow to
  switch agents). **`←`/`→` and `tab`/`shift+tab` = aliases.**
- **`PgUp`/`PgDn`/`Home`/`End` (+ `ctrl+u`/`ctrl+d`) = scroll** the active
  transcript/Status pane; scrolling up drops follow, `End`/bottom re-enables it.
- **`esc`** clears a non-empty input, else detaches/closes the TUI.
- **`ctrl+c`** cancels the run (live) / closes resume — unchanged.
- **All printable characters always append to the input**, including the legacy
  `q`/`?`/`a`/`f`/`i`/`j`/`k`. No single-letter global hotkeys remain.
- v1 input is **append-only** (backspace + `ctrl+u`-clear); no in-line cursor
  editing. Documented trade-off: arrow line-scroll is intentionally dropped
  because follow/auto-tail dominates usage; fine scroll is `PgUp`/`PgDn`.

### D4 — Persistent input + Enter routing
- One always-present input row. **Enter**: if the active agent has an open HITL
  question AND the input is non-empty → **answer** the oldest open question for
  that agent (`hitl.Answer`); else → **steer** the active agent (`steer.Submit`);
  on the `Status` tab → **deck** steer. Empty Enter never submits an empty answer
  (show a status hint, stay in answer mode).
- **Prompt labels** make the route visible BEFORE Enter: `steer <agent> ›`,
  `answer <agent>/<qid> ›`, `deck ›`. In answer mode the **whole prompt row is
  colour-flipped** (warn/red) + a one-line banner, so an answer is never
  accidental (hermes hardening).

### D5 — Slash commands
- Core set: `/help`, `/status`, `/follow`, `/deck <text>`, `/answer <qid> <text>`,
  `/quit`. Typing a leading `/` pops a one-line command-hint list. (`/answer` is
  the explicit escape hatch for the answer-vs-steer ambiguity / multiple
  questions.) Defer `/agent <id>` and `/tab <id>` unless implementation needs
  them.

### D6 — HITL questions
- Banner when the active agent has an open question + `q:N` in the status line +
  the colour-flipped answer prompt (D4). The full question list stays on the
  `Status` tab. `/answer <qid> <text>` is the unambiguous escape hatch.

### D7 — Per-agent transcript buffers
- Per agent: a bounded transcript buffer (reuse `loadFocusTail` /
  `readAppendedLines` / `capFocusLines`, 20k lines / 4 MiB) + scroll + follow
  state, so switching tabs is instant and retains position. Load **active +
  visited** agents only; refresh those each 250 ms tick (don't preload every
  agent). Bounded by the existing cap × visited.

### D8 — Reuse unchanged
- `internal/steer` (queued `new_attempt`, honest "recorded/queued" wording),
  `runstate` projection + segment/`[FINISHED]` fix, `hitl`, resume behavior,
  `--no-tui`. **No new event types. No `bubbles`/`bubbles/textinput` dependency
  in v1.** Keep `modeHelp` as an overlay (opened by `/help`); retire
  `modeOverview`/`modeAgentDetail`/`modeCompose`/`modeAnswerQuestion` as
  top-level states (folded into the tabbed view + persistent input).

### D9 — Tests-first, lock the contract
- BEFORE/ALONGSIDE the rewrite, lock `steer.Submit` + `hitl.Answer` + segment
  projection with tests so the input-routing change can't silently regress them.
- Add a **routing-table test**: printable → input (incl. `q`/`?`/`a`/`f`); `↑`/`↓`
  → tab switch; `PgUp`/`PgDn` → scroll; Enter → answer-before-steer (and deck on
  Status); `/`-commands route separately; `esc` clears-then-detaches; `ctrl+c`
  preserves cancel.
- Tab-strip overflow / narrow-width (80-col) tests; transcript follow tests;
  `Status` dashboard regression test.

## Trade-offs

- **Arrow line-scroll dropped** (↑/↓ now switch tabs). Accepted because follow is
  the default and dominant mode; `PgUp`/`PgDn` covers manual review. Documented
  in `/help`.
- **Answer-vs-steer in one box** — mitigated by the colour-flipped prompt row +
  banner + `q:N` + `/answer` escape, not the label alone.
- **Per-agent buffers** multiply the 4 MiB cap by visited agents — bounded and
  acceptable for 2–4 rosters; load active+visited only.
- **Default-surface change** breaks the existing dashboard/footer tests — these
  are intentional rewrites, not string edits.

## Deferred follow-ups

- `bubbles/textinput` cursor/paste editing (revisit arrow ownership then).
- `/agent <id>` and `/tab <id>` commands; a dedicated `Events` tab; mouse;
  configurable keymap; per-agent buffer LRU eviction for very large rosters.
- (Carried from tui-interactivity-overhaul) executing queued steers, live ACP
  delivery, opt-in "thinking" stream.

## Slices (implementation)
1. **Tab shell + frame** — tab strip (agents first, Status last), `activeTab`
   with stable IDs, default = first running agent, bottom status/input rows;
   old dashboard still renderable as `Status`.
2. **Transcript buffers** — active transcript in the agent-tab main area;
   per-agent visited buffers + follow/scroll + scroll keys.
3. **Input routing** — append-only input; Enter answer-before-steer (deck on
   Status); `esc`/`ctrl+c`; remove printable legacy hotkeys.
4. **Slash commands + help** — `/help`,`/status`,`/follow`,`/deck`,`/answer`,
   `/quit` + `/` hint list; `modeHelp` overlay behind `/help`.
5. **Tests + polish** — routing-table, overflow/narrow, follow, Status
   regression; status-line copy (`steer.List` counts), answer-mode colour.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept as drafted. Captures the round-02 convergence: ↑/↓ tab switch (owner),
agents-first/Status-last with stable IDs, append-only persistent input with
answer-before-steer + colour-flip hardening, slash commands incl. `/answer`,
per-agent buffers, all 1.12.0 plumbing reused unchanged, tests-first. No blockers.

### codex — ACCEPT (2026-06-04)
Accept as drafted. Confirms D1-D9: ↑/↓ tabs, agents-first/Status-last with stable IDs, append-only input with answer-before-steer, slash set including `/answer`, per-agent buffers, unchanged reuse of existing plumbing, and tests-first coverage.

### hermes — ACCEPT (2026-06-04)
From adversarial UX lens: answer-mode colour-flip + /answer escape, /-hint discoverability, documented arrow-scroll trade-off, and tests-first locking are all explicitly captured with no residual ambiguity.
