---
idea: tui-interactivity-overhaul
drafted-by: claude
date: 2026-06-04
participants: [claude, codex, hermes]
status: accepted
---

# Consensus — Interactive live-run TUI redesign

After round-01 (independent) and round-02 (cross-review), all three participants
converged on a single design. No blockers were raised. This consensus is the
basis for `FINAL.md`.

## Agreed decisions

### D1 — Sticky `[FINISHED]` is a run-state correctness bug; fix with segmented projection
- Refined root cause (codex, confirmed by all): `applyAgentEvent` already resets
  `finished → running` on a newer `agent.started` (`runstate.go:407-411`). The
  real defect is that **continue/resume/skip paths can begin new work without
  emitting a fresh segment boundary** (and sometimes without a fresh
  `agent.started` before the TUI renders), so a global per-agent state shows a
  stale terminal badge. The fix must make new work *explicit*, not rely on a
  lucky start event.
- **Fix:** introduce a `run.segment_started` event carrying `segment_id`,
  `reason` (`initial|continue|retry|steer`), `round`, and a **target-agent list**.
  Tag `agent.started|finished|failed|skipped` with `segment_id` (and optional
  `attempt_id`). `ProjectEvents` computes the *current segment* before applying
  terminal events; the badge reflects state of `(agent, current segment)`.
- **Compatibility:** old, unsegmented runs project into an implicit
  `segment-0001`; events lacking explicit IDs inherit the latest known segment.
- **Skipped/targeting rule (codex):** non-targeted agents are NOT reset by a
  segment that doesn't include them; agents targeted-but-skipped in a segment
  render `skipped` for that segment, never inherited `finished`.
- The minimal "emit start/segment on continue" patch ships **in the same slice**
  as the projection so the user-visible bug is fully (not partially) addressed.

### D2 — Segment key is a new monotonic id, not an overloaded round label
- Use `segment-NNNN` (monotonic per run) + per-agent `attempt_id`. Do **not**
  overload `round-NN`: a round may contain continue/retry/skip/steer attempts
  that each need their own current-state reset without masquerading as review
  rounds.

### D3 — Per-agent focus view replaces the 6-line tail
- New `agentDetail` view backed by `bubbles/viewport` over the **full**
  `stdout.log`/`stderr.log`, read **offset-incrementally** (mirroring
  `readEventsFromOffset`, `live.go:640-678`) — not a 4 KiB/6-line re-read.
- **Follow mode** default-on (`f`), auto-scroll on new bytes; any manual scroll
  disables follow until `f`/`G`. `g`/`G` top/bottom.
- **Bounded scrollback** is mandatory: cap **20,000 lines or 4 MiB per stream**
  with a visible truncation marker (prevents OOM on long runs).
- Overview right-pane preview becomes "as many lines as fit" instead of a hard 6.
- fsnotify stays **deferred**; 250 ms poll + offset reads are sufficient for v1.

### D4 — View-state machine + keymap + help overlay
- Model modes: `overview | agentDetail | compose | answerQuestion | help`
  (replacing overloaded booleans; `answerMode` renamed `answerQuestion`).
- Keymap: `enter`/`o` open agent detail; `esc` back/detach; `f` follow; `g`/`G`
  top/bottom; `j`/`k`/arrows select (overview) or scroll (detail); `i` compose
  for selected agent, `I` compose for deck/run; `a` answer open HITL question
  (unchanged); `?` **help overlay**; `q` detach, `ctrl+c` cancel (unchanged).
- **`?` moves from answer to help.** `a` remains the sole HITL answer key, so the
  finalized `hitl-tui-questions` capability is preserved. Record the `a/?`
  binding change in the feature's design record/changelog. **No Phase 0–8
  protocol amendment** (keymap/UI change, not a protocol decision).
- Persistent footer hint advertises contextual keys (discoverability).

### D5 — Composer is driver-owned run-control, never process control
- `internal/tui` exposes `SubmitSteering(ctx, SteeringRequest) (SteeringResult,
  error)` mirroring the existing `ActionRunner` boundary (`app.go:39-50`); the
  **app/runner** owns all mutation. TUI stays an observer/controller facade.
- Durable, auditable events: `steer.requested {id,target:"agent|deck",agent,text,
  created_by,segment_id,risk}`, then `steer.accepted|steer.rejected` (with
  reason) and `steer.delivered` (with mode).
- **Delivery modes:** `new_attempt` (universal default), `acp_prompt` (live
  runner-owned ACP session), `native_resume` (CLI resume), `queued_only`,
  `manual_handoff`.
- **One-shot truth:** after `cmd.Run()` owns fixed stdin (`runner.go:263-304`)
  there is NO mid-run stdin injection. One-shot agents get a **queued
  `new_attempt`** after the current attempt reaches terminal state — never fake
  "live" steering.
- The UI must **label delivery state explicitly** ("queued new attempt",
  "delivered to live ACP session", "native resume pending approval",
  "unsupported for this agent") — no chat-like over-promising.
- Steering must **never bypass** HITL questions or action/risk gates
  (`hitl-tui-questions`, `tui-action-execution`).
- **CLI parity:** add `parley steer RUN --agent AGENT -- TEXT` (additive), and
  surface queued steering + segment state in `continue --json`, keeping
  `--no-tui` first-class.

### D6 — ACP "thinking" is opt-in, not a default surface
- ACP thought chunks are currently buffered then discarded (`acp.go:207-211`).
  Persist only as **opt-in** `agent.acp.thought_chunk` events (or
  provider-approved thought summaries), gated by agent capability/config.
- Default disabled: the transcript shows ACP message chunks, tool calls, and
  plan events, and renders "thoughts unavailable" rather than pretending.
- Rationale: raw chain-of-thought is provider-sensitive.

### D7 — Agreed slice ordering (incremental, reviewable)
1. **Segment projection + sticky-badge fix** (D1, D2) — events, segment-aware
   `ProjectEvents`, emit-segment-on-continue, old-run compatibility, table tests
   (old unsegmented, continue-after-finished, retry-after-failed, skip-in-segment,
   unknown agent, round completion).
2. **Per-agent focus viewport** (D3) — `agentDetail` mode, offset reads, follow,
   bounded scrollback, `g`/`G`/`f`.
3. **Modes + keymap + help** (D4) — mode enum, `?` help overlay + footer hints,
   `a` answer preserved.
4. **Composer persistence + queued delivery** (D5) — `SubmitSteering`, `steer.*`
   events, default queued `new_attempt`, `parley steer` CLI; no live delivery yet.
5. **Live delivery + opt-in thoughts** (D5 `acp_prompt`/`native_resume`, D6) —
   capability-gated, per-delivery user opt-in.

Slices 1–4 are the **committed scope** (slice 4 satisfies the user's explicit
"where do I give the next prompt" requirement via queued `new_attempt`). Slice 5
is deferrable gold-plating.

### D8 — Invariants
- No Phase 0–8 cooperation-protocol amendment (unanimous).
- `--no-tui` headless path and the canonical `events.jsonl` + artifact contracts
  are preserved and only **extended** (additive event types/fields).
- Mutation ownership stays in `internal/app`/runner; `internal/tui` is a facade.

## Trade-offs

- **Segment-first vs viewport-first ordering:** hermes initially favored shipping
  visibility first; resolved in favor of **badge/projection first** because the
  stale badge is a correctness bug affecting both TUI and headless status, it is
  the smaller, more testable change, and the segment identity is the foundation
  the composer's queued attempts depend on.
- **Composer scope:** hermes flagged risk of an over-promising chat-like
  composer; mitigated by explicit delivery-state labels and by making
  `new_attempt` (queued) the honest default while live `acp_prompt`/resume are
  capability-gated and deferred to slice 5. The composer itself stays in
  committed scope because it is the user's explicit requirement.
- **fsnotify vs poll:** deferred; 250 ms + offset reads judged sufficient for
  parity, revisit only on reported lag.

## Deferred follow-ups

- F1 — Live `acp_prompt` delivery to a runner-owned ACP session (slice 5).
- F2 — `native_resume` relaunch for resume-capable CLIs, capability-declared +
  per-delivery user opt-in (slice 5).
- F3 — Opt-in ACP thought persistence / provider-approved summaries (slice 5).
- F4 — Cross-process append safety: the event-store append mutex is process-local
  (`events.go:24-57`); `parley steer` running beside a live TUI needs atomic
  append / file-locking review before multi-entry-point steering is durable.
- F5 — fsnotify-based refresh as a latency optimization.
- F6 — Viewport search/filter (`/`) within the active stream.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept as drafted. This captures my round-02 position: codex's segmented
projection as the backbone with the emit-on-continue fix in the same slice,
hermes's MVP discipline on slice scope, and the user's composer requirement
preserved as committed slice 4 with queued `new_attempt` as the honest default.
No blockers.

### codex — ACCEPT (2026-06-04)
Accept as drafted. This confirms the segmented projection with emit-on-continue as
the correctness fix, the driver-owned composer with queued `new_attempt` as the
default delivery mode, the `?` to help rebind with `a` preserved for HITL
answers, and slices 1–4 as the committed implementation scope.

### hermes — ACCEPT (2026-06-04)
From adversarial-interaction lens: slices 1-4 committed with slice 5 deferred is the correct minimal viable boundary; the honest queued-new_attempt composer default plus explicit delivery-state labels eliminates over-promise risk; bounded scrollback (20k lines / 4 MiB) is non-negotiable for long-run safety. All invariants preserved.
