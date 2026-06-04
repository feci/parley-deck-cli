---
agent: claude
idea: tui-interactivity-overhaul
round: 1
date: 2026-06-04
---

## Summary

The current TUI is a passive dashboard: a poll-projected table + a 6-line tail.
To feel like Codex/Claude-Code it needs three things the dashboard model can't
express: (1) a **focusable, scrollable, follow-capable transcript per agent**;
(2) a **segment-scoped state model** so badges reflect the *current* run segment,
not the last terminal event ever seen; (3) an **input composer** that turns user
text into a typed intent (answer / steer / new-round) reconciled with the
one-shot invocation contract. I propose a small **view-state machine**
(Overview → AgentFocus → Compose) layered over the existing `events.jsonl`
projection, plus three additive event types, and a per-agent log viewport backed
by the full `stdout.log` rather than a tail. No change to `--no-tui` or the
runner's artifact contract.

## Proposed approach

### A. Rich per-agent view (watch what it's doing)
- Introduce a **view-state enum** in the live model: `ViewOverview` (today's
  two-pane), `ViewAgentFocus` (one agent, large), optionally `ViewCompose`.
- `enter` on a selected agent → `ViewAgentFocus`: a `bubbles/viewport` filling
  the body, fed by the **entire** `stdout.log` (+ interleaved `stderr`), not the
  4 KB/6-line tail. Keep ANSI strip. Header shows id, state, elapsed, current
  round, artifact path.
- **Follow mode** (`f`, default ON in focus): auto-scroll to bottom as the log
  grows, like `less +F`; any manual scroll (`pgup`,`k`) drops follow until `f`
  or `G`. `g`/`G` jump top/bottom.
- Back to overview with `esc`. Overview's right pane keeps a *preview* but the
  preview is now "last N lines that fit" rather than a hard 6.
- Read strategy: track a per-agent file offset + ring buffer; on each tick read
  appended bytes only (cheap), so focus view scales to large logs.

### B. Fix sticky `[FINISHED]` — segment-scoped state
- Root cause: `ProjectEvents` folds all of `events.jsonl` into a single
  per-agent `State`, so the newest terminal event for an agent ID wins forever.
  On continue/resume the new segment's `agent.started` may not arrive (skipped /
  artifact exists) or arrives but the projection still shows the prior
  `finished`.
- Fix: **scope agent state to a segment key** (round label or a monotonic
  `segment`/`run-cursor`). Emit a **`round.started` / `segment.started`** event
  that resets the participating agents to `pending` for that segment. The TUI
  renders the *current* segment's state, not the global last state. Historical
  states remain in the projection for the focus/transcript view.
- Define precisely: badge = state of (agent, currentSegment). If the agent has
  no event in the current segment yet → `pending`/`queued`, never inherited
  `[FINISHED]`.

### C. Input composer (steer / next prompt) — reconciled with one-shot
- Generalize the HITL answer box into a **Compose bar** (`i` to open) that
  produces a typed **intent**, routed by the driver:
  1. **Answer** — if a HITL question is open for the focused agent, same as today.
  2. **New round / re-run** — enqueue the typed instruction as the seed for the
     agent's *next* invocation. Because agents are one-shot, "steering a running
     agent" = (a) if the CLI is **resume-capable** (per `headless-agents` caps,
     e.g. claude `--resume`) and the user opted in → resume that session with the
     new turn; (b) otherwise → a fresh one-shot invocation seeded with the
     prior artifact + the steer text (a cross-review-style round).
  3. **Deck-level** — a note that becomes the next round's prompt for all agents.
- The composer never executes side-effects directly; if an intent maps to a
  gated action it goes through the existing HITL/gate path. Emit a
  **`user.steer`** event so the steer is auditable in `events.jsonl` and mirrored
  into the canonical round/inbox file (protocol already requires phase-affecting
  decisions to be mirrored).
- Make the one-shot-vs-resume choice explicit and visible ("queued next turn"
  vs "resuming session") so the user understands latency.

### D. Interaction model / keymap (Codex/Claude-Code parity)
- Keep `q` detach, `ctrl+c` cancel. Add: `enter` focus, `esc` back, `f` follow,
  `g`/`G` top/bottom, `pgup`/`pgdn`/`k`/`j` scroll-or-select (context-sensitive),
  `i` compose, `tab` cycle agents, `?` help overlay. A persistent **footer
  hint + `?` overlay** documents the contextual keymap (the discoverability gap).
- Lower latency: prefer **fsnotify** on the run dir with the 250 ms tick as a
  fallback heartbeat; keeps CPU low and feels live.
- Persist ACP **thoughts**: flush `thoughtBuf` to a `thoughts.log` (or
  `agent.thought` events) so the focus view can show a "thinking" stream — this
  is the closest thing to "what it's reasoning about" and is currently thrown
  away.

### Sequencing (incremental slices)
1. View-state machine + AgentFocus viewport over full `stdout.log` + follow mode.
2. Segment-scoped state + `round.started` reset (fixes `[FINISHED]`).
3. `?` help overlay + keymap cleanup + fsnotify.
4. Compose bar → answer/new-round intents + `user.steer` event + mirror.
5. ACP thought persistence + thinking stream in focus view.

## Concerns / open questions

- **Steering semantics is the crux.** Does "give it the next prompt" mean resume
  the same session (continuity, but violates one-shot default and only some CLIs
  support it) or seed a fresh round? I lean: default = fresh seeded round;
  resume only when CLI-capable AND user opts in per invocation.
- **Segment key definition** — is it the `round-NN` label, a run cursor, or a
  new `segment` counter for continue/resume within a round? Need codex's read of
  how continue/resume currently emits events.
- **Does this amend `interactive-agent-mode` / `continuous-run-tui`?** Likely
  extends them; need to confirm we're not contradicting a finalized decision.
- Scrollback memory bounds for very long logs (ring buffer / cap?).
- Should focus view interleave stdout+stderr+thoughts or tab between them?

## Risks

- **Scope creep**: this is 5 slices; risk of a mega-PR. Mitigate by shipping
  slices 1–2 first (the user's loudest pain: visibility + sticky badge).
- **One-shot contract violation**: a "resume" composer could silently create
  hidden-session continuity the protocol forbids by default — must be opt-in and
  audited.
- **State-model regression**: re-keying agent state by segment could break
  existing resume/skip projection and prior TUI ideas' assumptions; needs tests
  locking current behavior first.
- **Latency/CPU**: reading full logs each tick is wasteful; offset-incremental
  reads + fsnotify required, not optional, for large runs.
- **Gate bypass**: composer-driven intents must not become a path around
  production HITL gates.
