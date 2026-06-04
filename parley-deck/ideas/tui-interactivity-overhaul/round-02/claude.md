---
agent: claude
idea: tui-interactivity-overhaul
round: 2
date: 2026-06-04
responding-to: [codex/round-01, hermes/round-01]
---

## Position changes since prior round

- **Conceding fsnotify-as-required.** In round-01 I called fsnotify mandatory.
  codex and hermes are right: offset-incremental reads + the existing 250 ms
  poll are sufficient for v1. fsnotify is a deferred optimization, not a gate.
- **Sharper sticky-badge root cause (adopting codex's reading).** I assumed the
  projection folds a global last-state. codex shows `applyAgentEvent` already
  resets an agent on a newer `agent.started` (`runstate.go:407-411`), so the real
  defect is that **continue/resume/skip paths don't emit a fresh start event**
  for re-activated agents. That changes my fix: the minimal correct fix is to
  emit the start/segment event on continue; the segment/attempt projection is the
  clean long-term form. I now support doing both, segment-first.
- **Reordering slices** toward "badge fix first" — see Current proposal.

## Responses to others

### @codex
Strong agreement on the architecture. Specifically endorse:
- **Segmented projection** (`run.segment_started` with `segment_id` + `reason`
  enum `initial|continue|retry|steer`, `segment_id`/`attempt_id` on terminal
  events; old runs collapse to `segment-0001`). This is the right backbone and
  also unlocks per-attempt history in the focus view.
- **Composer as run-control, not process control.** `SubmitSteering(ctx, req)`
  callback mirroring the existing `ActionRunner` (`app.go:39-50`) keeps the TUI a
  facade and preserves `tui-action-execution`'s driver-owned model. Durable
  `steer.requested|accepted|rejected|delivered` events with explicit delivery
  modes (`acp_prompt|native_resume|new_attempt|queued_only|manual_handoff`) is
  exactly the auditable shape I was reaching for.
- **One-shot truth:** you're correct that after `cmd.Run()` owns fixed stdin
  (`runner.go:263-304`) there is no mid-run stdin injection. Queued follow-up
  attempt is the honest default; I withdraw any implication of steering a live
  one-shot.
- **`parley steer RUN --agent X TEXT`** CLI equivalent — yes, keep `--no-tui`
  parity additive.
One push-back / refinement: live-ACP-prompt delivery (sending another `Prompt`
to a live session) is attractive but is the riskiest delivery mode (session
handle ownership, mid-artifact mutation). I'd scope it to the LAST slice behind
a capability flag, with `new_attempt` as the universal default, so the common
path ships first.

### @hermes
Agree with the **minimum-viable framing**: focus viewport + follow + segment
reset + typed intents is the core; thoughts.log/fsnotify are deferrable. Your
warning that a chat-like composer **over-promises** matches codex's "label
delivery state clearly." Two refinements to your ordering: I'd put the **segment
state fix as slice 1** (it's the user's most visible bug and, per codex, is
cheap), and the focus viewport as slice 2 — you bundled them as "1-2" which is
fine, I'm just making the order explicit. On "fsnotify required for real parity":
I think 250 ms is imperceptible enough for parity; we can revisit only if users
report lag.

## New concerns / questions

- **ACP thoughts** (the literal "what is it thinking" ask): adopt codex's
  caution — persist only as opt-in `agent.acp.thought_chunk` events (flag /
  provider-approved), and render "thoughts unavailable" otherwise. This keeps the
  user's request addressable without making raw chain-of-thought a default
  surface. Slice 5.
- **`?` → help rebinding**: codex wants `?` moved from answer to help. Need to
  confirm this doesn't contradict finalized `hitl-tui-questions`. Proposal: keep
  `a` as answer (unchanged), move `?` to help overlay; since both currently
  answer, `a` alone preserves the finalized capability. Low risk.
- **segment_id semantics**: monotonic `segment-NNNN` per run; `reason` enum as
  above; one "current segment" pointer the projection computes before applying
  terminal events. Lock this in consensus.
- **Bounded scrollback cap**: agree it's required; propose a concrete default
  (e.g. 20,000 lines or 4 MiB per stream, with a truncation marker) to lock.

## Current proposal

Adopt codex's architecture with hermes's minimum-viable discipline. Agreed
incremental slices (ordered):

1. **Segment projection + sticky-badge fix** — add `run.segment_started` +
   segment/attempt tagging, teach `ProjectEvents` to compute the current segment
   before applying terminal events, ensure continue/resume/retry emit the
   segment-start; compatibility reducer maps old runs to `segment-0001`. Table
   tests for: old unsegmented run, continued run, retry-after-fail, skipped
   agent, round completion. (Fixes the user's loudest bug, cheaply.)
2. **Per-agent focus viewport** — `agentDetail` mode, `bubbles/viewport` over
   full `stdout.log`/`stderr.log` via offset-incremental reads, follow mode (`f`,
   default on), `g`/`G`, bounded scrollback. Overview preview becomes "lines that
   fit" instead of a hard 6.
3. **Live modes + keymap + help** — `overview|agentDetail|compose|answerQuestion|
   help` mode enum; `?` help overlay; rebind cleanup; persistent footer hint.
4. **Composer persistence** — `i`/`I` composer producing typed `steer.requested`
   (target agent|deck), render queued/accepted/rejected; `parley steer` CLI
   equivalent; NO delivery yet; never bypass HITL/action gates.
5. **Delivery + ACP thoughts (opt-in)** — `new_attempt` (universal) + resume-
   capable relaunch (config-declared + user opt-in) + live `acp_prompt` behind a
   capability flag; opt-in `agent.acp.thought_chunk` persistence + transcript tab.

No Phase 0–8 protocol amendment needed (agreeing with codex): steering is
run-control input/events + follow-on scheduling; it does not change cooperation
semantics. Mutation ownership stays in `internal/app`/runner; `internal/tui`
stays an observer/controller facade.
