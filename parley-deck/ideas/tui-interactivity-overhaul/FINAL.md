---
idea: tui-interactivity-overhaul
status: final
author: claude
consensus-date: 2026-06-04
participants: [claude, codex, hermes]
---

## Final plan / specification

Rework the `parley` live-run TUI from a passive dashboard into an interactive,
Codex/Claude-Code-style controller. Unanimous consensus (claude, codex, hermes).
The TUI stays an **observer/controller facade** over `internal/app`/runner
primitives; the canonical `events.jsonl` + artifact contracts and the `--no-tui`
path are preserved and only **extended** (additive event types/fields). No
Phase 0–8 cooperation-protocol amendment.

### Problem recap (verified in code)
1. Only a 6-line / 4 KiB `stdout.log` tail per agent (`live.go:576-594`); no
   scrollable transcript, focus view, or follow mode.
2. Sticky `[FINISHED]` after continue/resume. Real cause: `applyAgentEvent`
   already resets `finished→running` on a newer `agent.started`
   (`runstate.go:407-411`), but continue/resume/skip paths begin new work
   **without emitting a fresh segment boundary/start event**, so global per-agent
   state shows a stale terminal badge.
3. No composer to inject a follow-up/steering prompt mid-run; the only text input
   is HITL answer mode bound to an open question (`live.go:148-154`).
4. Weak, undiscoverable interactivity overall.

### Design

**1. Segmented run-state projection (fixes sticky `[FINISHED]`).**
- New event `run.segment_started` with fields: `segment_id` (monotonic
  `segment-NNNN` per run), `reason` (`initial|continue|retry|steer`), `round`,
  and `targets` (target-agent list).
- Tag `agent.started|agent.finished|agent.failed|agent.skipped` with `segment_id`
  (and optional `attempt_id`).
- `ProjectEvents` computes the **current segment** before applying terminal
  events; a badge reflects the state of `(agent, current segment)`.
- Reset rules: a new segment resets ONLY its `targets` to `pending`/`running`;
  non-targeted agents keep their state; targeted-but-skipped agents render
  `skipped` for that segment (never inherited `finished`).
- Compatibility: old unsegmented runs map to implicit `segment-0001`; events
  without explicit ids inherit the latest known segment.
- Emit the segment-start on continue/resume/retry in the SAME slice so the
  user-visible bug is fully fixed.
- Tests (table): old unsegmented run, continue-after-finished, retry-after-failed,
  skip-in-current-segment, unknown/late agent, round completion.
- Use a new monotonic `segment_id`; do NOT overload `round-NN`.

**2. Per-agent focus view (replaces the tail).**
- New `agentDetail` view: a `bubbles/viewport` over the **full**
  `stdout.log`/`stderr.log`, read **offset-incrementally** (mirror
  `readEventsFromOffset`, `live.go:640-678`), not a 4 KiB re-read each tick.
- Follow mode default-on (`f`); auto-scroll on new bytes; any manual scroll
  disables follow until `f`/`G`. `g`/`G` jump top/bottom.
- **Bounded scrollback (mandatory):** cap 20,000 lines OR 4 MiB per stream, with
  a visible truncation marker.
- Overview right-pane preview becomes "as many lines as fit" (not a hard 6).
- fsnotify deferred; 250 ms poll + offset reads suffice for v1.

**3. View-state machine + keymap + help.**
- Modes: `overview | agentDetail | compose | answerQuestion | help` (replace
  overloaded booleans; rename `answerMode → answerQuestion`).
- Keymap: `enter`/`o` open detail; `esc` back/detach; `f` follow; `g`/`G`
  top/bottom; `j`/`k`/arrows select (overview) or scroll (detail); `i` compose
  for selected agent; `I` compose for deck/run; `a` answer open HITL question
  (unchanged); `?` help overlay; `q` detach, `ctrl+c` cancel (unchanged).
- `?` moves from answer to help; `a` stays the sole HITL answer key (preserves
  `hitl-tui-questions`). Record the `a/?` change in the changelog. Persistent
  footer hint for contextual keys.

**4. Driver-owned steering composer (the "next prompt" affordance).**
- `internal/tui` exposes `SubmitSteering(ctx, SteeringRequest) (SteeringResult,
  error)` mirroring the `ActionRunner` boundary (`app.go:39-50`); app/runner owns
  all mutation.
- Durable events: `steer.requested {id,target:"agent|deck",agent,text,created_by,
  segment_id,risk}` → `steer.accepted|steer.rejected` (reason) → `steer.delivered`
  (mode).
- Delivery modes: `new_attempt` (universal default), `acp_prompt`, `native_resume`,
  `queued_only`, `manual_handoff`.
- One-shot truth: after `cmd.Run()` owns fixed stdin (`runner.go:263-304`) there
  is NO mid-run stdin injection. One-shot agents get a **queued `new_attempt`**
  after the current attempt reaches terminal state — never fake live steering.
- UI labels delivery state explicitly ("queued new attempt", "delivered to live
  ACP session", "native resume pending approval", "unsupported for this agent").
- Steering NEVER bypasses HITL questions or action/risk gates.
- CLI parity: add `parley steer RUN --agent AGENT -- TEXT` (additive); surface
  queued steering + segment state in `continue --json`.

**5. ACP "thinking" (opt-in).**
- Thought chunks are currently buffered then discarded (`acp.go:207-211`).
  Persist only as opt-in `agent.acp.thought_chunk` (or provider-approved
  summaries), gated by agent capability/config. Default disabled → transcript
  shows message chunks, tool calls, plan events, and "thoughts unavailable".

### Implementation slices (ordered)
1. Segment projection + sticky-badge fix (events + segment-aware `ProjectEvents` +
   emit-on-continue + compatibility + table tests).
2. Per-agent focus viewport (offset reads, follow, bounded scrollback, `g`/`G`/`f`).
3. Modes + keymap + help (`?` overlay, footer hints, `a` preserved).
4. Composer persistence + queued delivery (`SubmitSteering`, `steer.*` events,
   default queued `new_attempt`, `parley steer` CLI; no live delivery yet).
5. Live delivery + opt-in thoughts (`acp_prompt`/`native_resume` capability-gated +
   per-delivery user opt-in; opt-in `agent.acp.thought_chunk`).

**Committed scope: slices 1–4** (slice 4 satisfies the user's explicit "where do
I give the next prompt" requirement via queued `new_attempt`). **Slice 5 deferred.**

### Invariants
- No Phase 0–8 protocol amendment.
- `--no-tui` + canonical `events.jsonl`/artifact contracts preserved, only
  extended additively.
- Mutation ownership stays in `internal/app`/runner; `internal/tui` is a facade.
- No HITL/action/risk-gate bypass via steering.

### Deferred follow-ups
- F1 live `acp_prompt`; F2 `native_resume`; F3 opt-in ACP thoughts (all slice 5).
- F4 cross-process append safety: event-store append mutex is process-local
  (`events.go:24-57`); `parley steer` beside a live TUI needs atomic-append /
  file-lock review before multi-entry-point steering is durable.
- F5 fsnotify refresh; F6 viewport search/filter (`/`).

### References
- Consensus: ./consensus.md (claude/codex/hermes all ACCEPT, 2026-06-04)
- Rounds: ./round-01/, ./round-02/
- Extends finalized: live-run-tui, tui-agent-controls, interactive-agent-mode,
  continuous-run-tui, tui-layout-refresh, hitl-tui-questions, tui-action-execution,
  session-resume-cache-plan. Contradicts none.
