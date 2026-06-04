---
idea: tui-claude-cli-layout
author: user
created: 2026-06-04
participants: [claude, codex, hermes]
roles:
  claude: layout/state-machine + reuse-of-1.12.0 plumbing (facilitator)
  codex: TUI architecture + Bubble Tea input/focus model
  hermes: adversarial UX — real Claude-CLI/Codex-CLI parity
status: final
---

## Problem / idea

The live-run TUI (shipped in 1.12.0 via idea `tui-interactivity-overhaul`) still
*lands* on a two-pane dashboard (agent table + events + questions + a small log
preview). The owner finds it unclear and wants a **fundamentally different,
Claude-CLI-style layout** where the **live generated text is the main screen**,
not a dashboard. Verbatim intent:

- The MAIN screen shows **what the currently-selected agent is generating** — its
  streaming output/transcript — front and centre.
- Agents are **tabs across the top**; switch between them with **arrow keys**;
  the active tab's transcript fills the main area.
- The BOTTOM has a **prompt input line + a status line**, like Claude CLI —
  an always-present compose row and a status bar.
- Overall feel like **Claude Code / Codex CLI**: persistent tabs, persistent
  input, the generated text always visible, ability to flip to other tabs
  (status, etc.).

This is a **default-layout/UX overhaul of `internal/tui/live.go`**, NOT new
engine work. It must **reuse** the plumbing already shipped in 1.12.0:
- `runstate.ProjectEvents` + segment-scoped state (the sticky-`[FINISHED]` fix);
- the offset-incremental focus reads `loadFocusTail` / `readAppendedLines` /
  `capFocusLines` (bounded 20k lines / 4 MiB);
- the steering composer + `internal/steer` (`steer.Submit`, queued `new_attempt`);
- the `liveMode` state machine + `?` help overlay;
- HITL questions (`hitl` pkg) + answer flow.

## Goal

A consensus design (→ FINAL.md) then implementation for a Claude-CLI-style
tabbed, transcript-centric default layout. Decide concretely:

1. **Tab bar** (top): one tab per agent (id + state colour); plus special tabs
   ("Status" = the old dashboard, optionally "Events"). Rendering when many
   agents / narrow terminals. Navigation keys.
2. **Main transcript pane**: the active agent's full stdout via the existing
   offset-incremental reads + bounded scrollback; follow mode (auto-tail);
   scrollback keys. Per-agent buffer/scroll state vs reload-on-switch.
3. **Bottom prompt row + status line**: an always-present input (Claude-CLI
   style) wired to the steer composer; Enter answers an open HITL question for
   the active agent if one exists, else queues a steer. Status line: run/idea,
   round/segment, active agent state, follow on/off, queued-steers, errors.
4. **The central UX conflict**: an always-typeable input vs. tab/scroll
   navigation keys. Define exactly which keys switch tabs, scroll the
   transcript, and type into the input, so it feels like Claude CLI (not vim).
5. Where the old dashboard, the events stream, and HITL questions live in the
   new model (a "Status" tab and/or banners). Nothing from 1.12.0 regresses.

## Constraints

- Bubble Tea / lipgloss; Go 1.26. Terminal only.
- Keep `--no-tui` and the runner/`events.jsonl`/artifact contracts intact.
- REUSE 1.12.0 plumbing; do not re-fix the segment/badge or re-write the steer
  package. This is presentation + input routing.
- Preserve every 1.12.0 capability (segment fix, steer queue + `parley steer`,
  help) — the dashboard may move to a tab but must remain reachable.
- Keep it implementable in reviewable slices; bounded memory for N per-agent
  transcript buffers.
- Steering still records-only (no execution) per the prior consensus; honest
  wording stays.

## Non-goals

- No web/remote UI.
- No new engine/event types unless strictly required by the layout (prefer none).
- Not implementing slice-5 execution of queued steers, live ACP delivery, or the
  "thinking" stream (still deferred).
- No protocol (Phase 0–8) changes.
