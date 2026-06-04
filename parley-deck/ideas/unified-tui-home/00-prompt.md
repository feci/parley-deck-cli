---
idea: unified-tui-home
author: user
created: 2026-06-04
participants: [claude, codex, hermes]
roles:
  claude: TUI state model + merge of the two TUIs (facilitator)
  codex: app wiring (runTUI/runTask/runResume), agent discovery, run lifecycle
  hermes: adversarial UX — does the unified app actually feel right end-to-end
status: final
---

## Problem / idea

There are currently TWO TUIs and the owner wants ONE:
- the OLD workspace dashboard (`internal/tui/app.go` model, launched by `parley
  tui` via `internal/app` runTUI) — ideas/runs/sessions/actions panes, `N`=new
  idea, `StartRunFunc`, `ActionRunner`;
- the NEW Claude-CLI-style tabbed live view (`internal/tui/live.go` `RunLive`,
  launched by `parley run`/`parley resume`) — agent tabs + live transcript +
  bottom input/status (shipped 1.13.0, idea `tui-claude-cli-layout`).

Merge them into a single unified `parley tui`, retiring the old dashboard.

## Owner decisions (locked)

1. **Entry model = Home tab + agent tabs (one app).** `parley tui` opens to a
   **Home tab** listing ideas / recent runs with **`N` = start a new idea**.
   Starting an idea picks **all available agents**, launches them, and
   auto-switches to their **live transcript tabs**. Home is always a tab you can
   flip back to (`↑/↓` cycles Home ↔ agent tabs ↔ Status). After a run finishes
   the TUI **stays open**; `N` starts another.
2. **Default participants = ALL AVAILABLE (installed/discoverable) agents** when
   the user passes none. The TUI must show **only available agents** as tabs —
   never render unavailable / not-found / inactive agents.
3. **Do not exit the TUI when the run finishes.** Keep it open so the user can
   keep typing, steer, and start new work.
4. **Each agent tab shows its live generated output** (the per-agent transcript).
   Verify it actually populates during a real run and fix if not visible.
5. **`N` starts a new idea from inside the TUI** (the old workspace TUI had this;
   port that flow).
6. **Retire the old `parley tui` dashboard.** The new unified tabbed app IS
   `parley tui`. Fold in only the useful bits of the old `app.go` model
   (start-idea via `N`, ideas/runs listing on Home); discard the rest.

## Goal

A consensus design (→ FINAL.md) then implementation for the unified `parley tui`.
Decide concretely:

- **Home tab**: what it lists (open ideas + recent runs), selection + `N` flow
  (prompt for task → discover all available agents → launch → switch to agent
  tabs). How it reuses the old model's idea/run listing + StartRunFunc.
- **Run lifecycle in one model**: a `RunLive`-style model that can be in a
  no-run state (Home only) AND an active-run state (agent tabs). `parley run
  "task"` launches directly into agent tabs; `parley tui` opens to Home; both are
  the SAME model. Don't quit on done; allow starting a new run in-place.
- **Available-agents filtering**: where `internal/agents` discovery feeds the
  default participant set and the tab list; never show unavailable agents.
- **Per-agent transcript correctness**: confirm/fix that each agent tab shows
  live stdout (the focus-buffer path) during a real multi-agent run.
- **Key/slash model**: extend the existing one — `N` new idea, `↑/↓` tabs,
  Enter steer/answer, `/help` etc. Keep `--no-tui` and the runner/events
  contracts intact.
- **Retiring app.go**: what to delete vs fold into the unified model; keep
  `parley run`/`parley resume`/`parley tui` all wired to the one model.

## Constraints

- Bubble Tea / lipgloss; Go 1.26. Terminal only.
- Reuse UNCHANGED: events.jsonl projection + runstate segment plumbing,
  `internal/steer` (+ `parley steer`), `hitl`, the focus buffers
  (loadFocusTail/readAppendedLines/capFocusLines), the slash/key-routing model.
- Keep `--no-tui` and the runner/events/artifact contracts intact.
- Starting a run from the TUI must go through the existing runner launch path
  (the same one `parley run` uses), not a parallel engine.
- Incremental, reviewable slices; tests-first for the new state transitions.

## Non-goals

- No web/remote UI; no new event types unless strictly required.
- Not implementing queued-steer execution / live ACP delivery / thinking stream
  (still deferred).
- No Phase 0–8 protocol changes.
