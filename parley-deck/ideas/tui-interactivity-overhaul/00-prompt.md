---
idea: tui-interactivity-overhaul
author: user
created: 2026-06-04
participants: [claude, codex, hermes]
roles:
  claude: UX state-model and event-projection lens (facilitator)
  codex: TUI architecture + continuity with prior finalized TUI ideas
  hermes: adversarial interaction-design lens (Codex/Claude-Code parity)
status: final
---

## Problem / idea

The Parley Deck live-run TUI (`internal/tui`: `live.go`, `app.go`, `layout.go`)
feels passive and weakly interactive. The user wants it reworked to feel like a
modern agent TUI (Codex CLI / Claude Code): a real input composer, scrollback,
per-agent focus/detail view, and follow mode.

Three concrete defects motivate the rethink (all verified in code):

1. **No real "what is the agent doing right now" view.** The only per-agent
   detail is a fixed 6-line / 4096-byte tail of `stdout.log`
   (`live.go:576-594`, `tailLogFile`). There is no scrollable viewport, no full
   transcript, no follow/tail toggle, and no per-agent focus screen. The event
   stream (`Latest events`) and the agent's last `LatestEvent` (64-char) are the
   only other signals.

2. **`[FINISHED]` badge is sticky after continue/resume.** Agent state is
   projected from `events.jsonl` by `ProjectEvents` (`internal/runstate`),
   keyed per agent. When the user continues/resumes a run, prior-round terminal
   events (`agent.finished`) keep the badge green/`[FINISHED]` even though the
   agent should be pending/running again for the new segment. State is not
   scoped to the current round/segment, so stale terminal state "stays lit."

3. **No way to inject a follow-up prompt mid-run.** Once a task starts there is
   no composer to give an agent (or the deck) its next instruction. The only
   text input is the HITL answer mode (`a`/`?`) bound to an *open question*
   (`live.go:148-154`, `472-495`). There is no general "send the agent its next
   prompt / steer it" affordance after a task has begun.

Secondary facts (context, not all in scope):
- Refresh is a 250 ms poll of `events.jsonl` (`live.go:618-627`), re-projected
  by `ProjectEvents`. No fsnotify/push.
- Agent output is captured to `parley-deck/runs/<run-id>/agents/<id>/stdout.log`
  + `stderr.log` (`runner.go:217-218`); artifacts land in
  `ideas/<slug>/round-NN/<id>.md` (`runner.go:220`).
- ACP "thinking"/thought chunks ARE received but **discarded** —
  `acpRunnerHandler.thoughtBuf` is accumulated and never flushed to disk or
  events (`acp.go:186-211`). Only visible `agent.acp.message_chunk` is kept.
- Headless CLI agents are **one-shot** subprocesses (Generic CLI Invocation
  Contract). Some CLIs expose resume/continue (e.g. `claude` config lists
  `--resume`/`--continue`/`--session-id`/`--fork-session`); others do not.
- Prior FINAL designs this must build on (not reinvent): `live-run-tui`,
  `tui-agent-controls`, `interactive-agent-mode`, `continuous-run-tui`,
  `tui-layout-refresh`, `hitl-tui-questions`, `tui-action-execution`,
  `tui-workspace-sessions`.

## Goal

A consensus design (→ `FINAL.md`) for an interactive, Codex/Claude-Code-style
live TUI, concretely covering: (a) a rich, scrollable, follow-capable per-agent
view that shows what each agent is doing; (b) a correct per-segment agent state
model that fixes the sticky `[FINISHED]` badge; (c) an input composer to send
follow-up/steering prompts to an agent or the deck mid-run, reconciled with the
one-shot vs resume-capable invocation contract; (d) the overall interaction
model — view states, keymap, scrollback, follow mode, help overlay — and how it
maps onto the existing `events.jsonl` projection + runner without breaking
`--no-tui` headless behavior.

## Constraints

- Bubble Tea (charmbracelet) TUI; Go 1.26 CLI.
- Keep `--no-tui` headless path and the canonical `events.jsonl` projection
  model intact; enrich, do not replace, the runner contract.
- Respect the Generic CLI Invocation Contract: one-shot by default; only use
  hidden-session resume where the agent CLI supports it AND the user opts in.
- Side-effects / mutations stay driver-owned and gated; a steering composer must
  not bypass HITL risk gates.
- Prefer building on the finalized TUI ideas above; if any must change, say which
  and whether it needs a protocol/amendment note.
- Keep the change implementable in incremental, reviewable slices.

## Non-goals

- No web/remote UI; terminal only.
- No change to the Phase 0–8 cooperation protocol semantics themselves.
- Not redesigning action-block execution or the provider/effects ledger.
- No mandatory dependency on a specific agent supporting session resume.
