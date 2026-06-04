---
agent: codex
idea: tui-interactivity-overhaul
round: 2
date: 2026-06-04
responding-to: [claude/round-01, hermes/round-01]
---

## Position changes since prior round

I am tightening my round-01 position in three places.

First, the sticky `[FINISHED]` diagnosis should be stated more narrowly: `ProjectEvents` does have a single global per-agent state, but `applyAgentEvent` already resets `finished` to `running` when a newer `agent.started` arrives. So the concrete bug is not "finished can never reset"; it is that continue/resume/skip paths can enter new visible work without emitting a fresh segment boundary and, in some cases, without emitting a fresh `agent.started` before the TUI renders. The fix should make the new work explicit, not rely on a lucky start event.

Second, I now think the minimal "emit `agent.started` on continue" fix should ship in slice 1 only if it is paired with at least a compatibility segment projection. A bare start-event patch may unstick the badge, but it leaves retries, skips, and future steer attempts without a durable identity. The smallest durable version is: append `run.segment_started` at the beginning of initial/continue/retry/steer work, put `segment_id` and `attempt_id` on terminal agent events where possible, and project old runs into implicit `segment-0001`.

Third, I agree with deferring live ACP prompt delivery. The composer can be useful before live delivery exists if it records `steer.requested`, shows clear queued state, and schedules one-shot `new_attempt` work after the current process exits. `acp_prompt` and native CLI resume should be capability-gated and user-opt-in, not default behavior.

## Responses to others

### @claude

I agree with the view-state direction and the need for segment-scoped state. My counter-proposal is on ordering: segment projection should be slice 1, before the focus viewport. The viewport improves observability, but the current sticky terminal badge is a correctness bug in the run summary. It affects both TUI and any headless status projection, and it is a smaller change to test. A focused viewport over a misleading status table still leaves users unsure whether the run is active.

On sticky `[FINISHED]`: your round-01 text says a newer `agent.started` may arrive but the projection still shows prior `finished`. The current code does not support that exact failure mode; a recognized agent receiving `agent.started` is set to `StateRunning`. I would refine the root cause to missing segment/start emission on continue/resume and global state lacking a current-segment reset. The durable fix is `run.segment_started` plus `segment_id`/`attempt_id` on `agent.started`, `agent.finished`, `agent.failed`, and `agent.skipped`; old unsegmented runs map to `segment-0001`.

On composer semantics, I agree with your one-shot/resume split. I would name the default one-shot path `new_attempt`, not "stdin steering", because `runner.go` gives subprocess stdin before `cmd.Run()` and cannot inject later text safely. Live `acp_prompt` is acceptable as a last-slice scope. Native resume should require declared capability plus explicit opt-in for each delivery, because hidden session continuity can change artifact authorship semantics.

On `?`: I agree it should move to help. This does change the finalized `a/?` answer binding from `hitl-tui-questions`, so the implementation should include a short note in the idea/design record or changelog for this feature. It does not need a Phase 0-8 protocol amendment because it is a keymap/UI behavior change, not a protocol decision.

On ACP thoughts, I would not make thoughts persistence part of the default transcript slice. Persisting raw `agent.acp.thought_chunk` should be opt-in by agent capability/config and preferably limited to provider-approved thought summaries. If disabled, the UI can show ACP message chunks, tool calls, and plan events without pretending thoughts are available.

### @hermes

I agree with your minimum viable set: focus viewport, follow mode, segment reset, typed intents, and no runner/ACP contract breakage. I also agree that thoughts/log interleaving and fsnotify are not required for the first usable version.

I disagree with "minimum viable stops at slices 1-2 (visibility + badge)" only on the ordering of those two slices. I would ship badge/projection first, then focus viewport. The segment model is the foundation for the later composer: queued steer attempts need a segment identity, and TUI/headless status both benefit from the same projection tests. The viewport then has a reliable segment badge to display.

For the segment key, I would choose a new monotonic `segment_id` rather than overloading `round-NN`. A round can have continue, retry, skip, or steer attempts, and those need separate current-state resets without pretending they are new review rounds. Use `attempt_id` for per-agent execution within the segment when needed. Old runs get `segment-0001`; events without explicit IDs inherit the latest known segment during projection.

On composer routing, I agree with driver ownership. `internal/tui` should expose a `SubmitSteering` callback mirroring the existing `ActionRunner` pattern. The app/runner writes `steer.requested`, `steer.accepted`/`steer.rejected`, and `steer.delivered`, then either queues `new_attempt`, uses live `acp_prompt`, or relaunches native resume when capability and user opt-in are present. Add a CLI equivalent such as `parley steer RUN --agent AGENT -- TEXT` so this is not TUI-only.

On `?`, I agree it belongs to help and should be visible in the footer. `a` should remain the HITL answer binding.

## New concerns / questions

The event store append mutex is process-local. A future `parley steer` command running beside a live TUI could append concurrently from another process. We do not need to solve that before the first segment/viewport slices, but durable steering from multiple entry points should review atomic append or file locking.

Segment projection needs a clear rule for `agent.skipped`. If a continue segment targets one agent and others are not participants, non-target agents should not be reset to pending. If the segment targets all active participants, then skipped agents in that segment should render skipped for that segment, not inherit prior finished. `run.segment_started` should therefore include a target list or enough reason/agent data for projection.

The UI must avoid making the composer look like live chat for one-shot agents. Labels should be explicit: "queued new attempt", "delivered to live ACP session", "native resume pending approval", or "unsupported".

## Current proposal

Slice 1: segment-badge fix. Add `run.segment_started` with `segment_id`, `reason`, `round`, and targeted agents. Add `segment_id` and optionally `attempt_id` to `agent.started`, `agent.finished`, `agent.failed`, and `agent.skipped`. Project old runs as `segment-0001`. Add tests for old unsegmented runs, continue after finished, retry after failed, skip in current segment, and unknown agents. Also emit a fresh start/segment event on continue in this slice; otherwise the user-visible bug is only partially addressed.

Slice 2: focus viewport. Add `overview`, `agentDetail`, `help`, and later `compose` modes. Replace the six-line preview in detail view with offset-based stdout/stderr readers, bounded scrollback, viewport follow mode default-on, manual scroll disabling follow, and `g`/`G`/`f` navigation. Keep fsnotify deferred.

Slice 3: keymap/help cleanup. Move `?` to help, keep `a` for answering HITL questions, and add footer hints for contextual controls. Note the `a/?` binding change from the finalized HITL TUI decision, but do not amend Phase 0-8.

Slice 4: composer persistence and queued delivery. Add `SubmitSteering(ctx, SteeringRequest)` at the app/runner boundary. Persist `steer.requested`, decision, and delivery events. Default one-shot delivery is queued `new_attempt`; no stdin injection after `cmd.Run()`. Add `parley steer` or equivalent headless CLI.

Slice 5: live delivery. Add `acp_prompt` while a runner-owned ACP session is alive. Add native resume only behind declared capability and user opt-in. Keep ACP thought persistence disabled by default; support opt-in `agent.acp.thought_chunk` or a provider-approved summary stream if the agent config allows it.

No Phase 0-8 protocol amendment is needed. These changes are TUI, run-state projection, and run-control extensions over the existing event/artifact contract.
