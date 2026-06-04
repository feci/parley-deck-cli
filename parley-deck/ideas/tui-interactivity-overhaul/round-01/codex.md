---
agent: codex
idea: tui-interactivity-overhaul
round: 1
date: 2026-06-04
---

## Summary

This redesign should EXTEND the finalized TUI decisions, not contradict them. The durable source of truth remains `events.jsonl` plus run-local files; `--no-tui` remains a headless/status path; the TUI remains an observer/controller facade over app/runner primitives, not a second protocol engine. No protocol/amendment note is needed unless we decide that mid-run steering changes Phase 0-8 semantics. My recommendation is to keep steering as run-control input/events and follow-on subprocess scheduling, so no amendment is needed.

The four requested fixes are related but should be split:

1. Replace the current selected-agent `stdout.log` / `stderr.log` preview (`internal/tui/live.go:576-594`) with a real per-agent detail view backed by scrollable viewports and bounded log readers.
2. Fix stale `[FINISHED]` by projecting agent state for the active run segment/attempt, not by letting the last terminal event for an agent win forever (`internal/runstate/runstate.go:320-424`).
3. Add a composer that records steering intent immediately, then routes it through the runner/app only when the target agent supports safe delivery. One-shot agents get queued follow-up attempts, not stdin injection.
4. Introduce explicit live-view modes and help/key state so agent selection, log scrolling, question answering, action execution, and steering do not collide.

## Proposed approach

Use a segmented run projection. Add an event concept such as `run.segment_started` with `segment_id`, `reason` (`initial`, `continue`, `retry`, `steer`), `round`, and optional `agent`; include `segment_id` or `attempt_id` on `agent.started`, `agent.finished`, `agent.failed`, and `agent.skipped`. Then teach `ProjectEvents` to derive the current segment before applying terminal events. Existing runs without segment events should map all events to an implicit `segment-0001`, preserving compatibility. This directly fixes the sticky badge: when a continue/retry/steer segment starts, participants targeted by that segment return to `pending` or `running` until new terminal events land, instead of preserving an earlier `agent.finished`.

For a smaller first fix, `ProjectEvents` can also reset an agent on any newer `agent.started` even after `finished`; `applyAgentEvent` already does that (`internal/runstate/runstate.go:407-411`). The observed sticky behavior likely comes from continue/resume flows that do not append a fresh start event for the resumed work. Segment events are the cleaner long-term answer because they also let the TUI show history per attempt.

Promote log display from preview string to per-agent buffers:

- Add a `logViewport` per selected stream using Bubble Tea's `viewport` component, with `follow` default true.
- Read `stdout.log` and `stderr.log` incrementally by byte offset, similar to `readEventsFromOffset` (`internal/tui/live.go:640-678`), not by rereading a 4 KiB tail every tick.
- Keep bounded in-memory scrollback, e.g. 5,000-20,000 lines or a byte cap per agent; show truncation markers when old lines are dropped.
- Provide tabs or a small selector for `Transcript`, `Events`, `stdout`, `stderr`, and `Artifact`. `Transcript` should merge summarized events, ACP message chunks, ACP tool/plan events, and log lines by arrival time. ACP thoughts should either be flushed as `agent.acp.thought_chunk` or deliberately omitted with a visible "thoughts unavailable" design; currently they are buffered and discarded (`internal/runner/acp.go:207-211`).
- Keep polling at 250 ms initially (`internal/tui/live.go:618-627`). fsnotify can remain deferred; offset-based reads are enough for this slice.

Add live modes to the model rather than overloading booleans:

- `overview`: current two-column table/events/questions/log summary.
- `agentDetail`: full selected-agent view with scrollback.
- `compose`: bottom composer overlay for steering/follow-up.
- `answerQuestion`: existing HITL answer flow, renamed from `answerMode`.
- `help`: modal keymap overlay.

Suggested keys:

- `j/k` or arrows: move selection in overview; scroll in detail when viewport focused.
- `enter` or `o`: open selected agent detail.
- `esc`: close detail/help/compose, or detach live view from overview.
- `f`: toggle follow mode for the active viewport.
- `g` / `G`: top / bottom.
- `/`: search/filter within current viewport, optional later.
- `i`: open composer for selected agent; `I`: composer targeting deck/run.
- `a`: answer selected open HITL question, preserving the prior decision in `hitl-tui-questions`.
- `?`: help overlay, not answer mode; this intentionally changes the old `a/?` answer binding from `internal/tui/live.go:148-154` because `?` is the conventional help affordance.

The composer should be a run-control action, not direct process control. Add a durable steering record under the run directory and mirror events:

- `steer.requested`: `{id, target:"agent|deck", agent, text, created_by:"tui", segment_id, risk}`.
- `steer.accepted` / `steer.rejected`: app/runner decision with reason.
- `steer.delivered`: delivery mode (`acp_prompt`, `native_resume`, `new_attempt`, `queued_only`, `manual_handoff`).

Delivery rules:

- ACP sessions can receive another `Prompt` while the session is alive if the runner owns the session handle. This requires the TUI/app to call a runner supervisor API, not write to files and hope the subprocess notices.
- Headless one-shot subprocesses cannot receive mid-run stdin after `cmd.Run()` owns fixed stdin/stdout (`internal/runner/runner.go:263-304`). For these, the composer queues a follow-up attempt after the current attempt reaches terminal state, or creates a HITL question/task requiring user-approved continuation.
- Resume-capable CLIs may be relaunched with configured resume flags only when runtime config declares that capability and the user opts in. This extends, rather than contradicts, `session-resume-cache-plan`: workflow recovery stays primary; native continuation is opportunistic.
- Interactive/manual launch modes remain user-driven handoffs per `interactive-agent-mode`; the TUI may generate the handoff packet or command preview but must not automate an interactive TTY.

Keep mutation ownership in `internal/app` / runner. `internal/tui` should expose callbacks like `SubmitSteering(ctx, SteeringRequest) (SteeringResult, error)` similar to existing `ActionRunner` (`internal/tui/app.go:39-50`). This preserves `tui-action-execution`'s driver-owned execution model and avoids making Bubble Tea code responsible for provider/session semantics.

Do not break `--no-tui`: status/resume/continue should still load `RunSummary` via `runstate.ResolveRun` and print next actions (`internal/app/app.go:965-1045`). New steering should have a CLI equivalent, e.g. `parley steer RUN --agent codex TEXT` or `parley continue --steer`, but it should be additive. Headless `continue --json` can expose queued steering and segment state without needing a TUI.

Incremental slices:

1. Projection: add segment/attempt events, compatibility reducer, tests for stale terminal events, and render segment labels.
2. Viewport: replace `logPreview` with offset-based selected-agent stdout/stderr readers and follow/scroll state.
3. Agent detail/help: add live modes, keymap overlay, and focused event/log transcript view.
4. Composer persistence: record `steer.requested` and render queued/accepted/rejected states; no delivery yet.
5. Delivery: ACP live prompt and resume-capable relaunch where config declares support; one-shot queued follow-up otherwise.

## Concerns / open questions

The largest design question is whether a steering prompt is meant to alter the current canonical artifact or create a new attempt/segment. I think it must create a new segment/attempt unless delivered to a live ACP session before artifact finalization; otherwise the audit trail becomes ambiguous.

We need a runtime capability schema for native resume before implementing provider-specific delivery. Existing config shows interactive `--resume` examples, but the runner does not currently model "safe hidden-session resume" as a headless capability.

ACP thought chunks are useful for "what is the agent doing", but they may expose sensitive chain-of-thought-like content depending on provider semantics. We should only surface provider-approved thought summaries or events, and treat raw thought text as opt-in.

The current event store uses a single append mutex only within one process (`internal/store/events.go:24-57`). If future steering commands are sent from another terminal while the TUI is open, local append behavior may need cross-process locking or atomic append discipline review.

## Risks

Segment projection can regress old run status if compatibility is under-tested. Add table tests for old unsegmented runs, continued runs, retry after failure, skipped agents, and round completion.

A composer that looks like chat can overpromise. The UI must label delivery state clearly: "queued for next attempt", "delivered to live ACP session", "native resume required", or "unsupported for this agent".

Large logs can make the TUI sluggish if we store/render too much. Bound scrollback, incremental reads, and viewport rendering are required; do not render every line in a giant string.

Keymap growth can make the live view harder to use. The help overlay should be implemented early, and `?` should move to help even though that changes the current answer shortcut.

Driver-owned gating must remain intact. Steering must not bypass HITL questions, action risk gates, or manual signoff rules established by `hitl-tui-questions` and `tui-action-execution`.
