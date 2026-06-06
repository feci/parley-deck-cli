---
idea: tui-live-steering
author: user
created: 2026-06-06
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + Bubble Tea TUI wiring
  codex: runner/app seams — per-agent process control + steer-execution path, Go correctness, tests
  agy: UX correctness + consistency (autocomplete, kill affordance, where the reply shows)
  hermes: keymap/interaction-model fidelity + concurrency/safety (no key collisions, no race when steering a running agent)
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem / idea (owner's words, three features)

In the unified `parley tui` (`internal/tui/live.go`, shipped 1.17.0) the owner wants:

1. **Slash-command autocomplete** — like the Claude Code CLI: typing the start of a
   command (e.g. `/o`) and pressing **Tab** completes it, and as you type `/…` a menu of
   matching commands is offered to pick from (arrow-select + confirm). No memorizing the
   exact command names.
2. **Kill a long-running agent** — from the TUI, the owner can terminate ONE agent that
   has been running too long, while the rest of the run continues.
3. **Steer round-trip that actually works** — when the owner types a message to an agent
   in the TUI (e.g. to `agy` or `hermes`), they expect Parley to **send it to that agent
   and show the agent's reply**. Today nothing visible happens and the owner cannot tell
   where the text went.

## Current state (VERIFIED against the code — design against these facts)

- **Steers are write-only.** `steer.Submit` (`internal/steer/steer.go:87-124`) appends
  `steer.requested` + `steer.delivered` events to `events.jsonl` and returns
  `status:"queued"`. The TUI `submitSteer` (`internal/tui/live.go`) records the steer and
  shows "queued; auto-exec not wired yet" — **no consumer re-invokes an agent or injects
  the text anywhere.** This is literally true repo-wide. THIS is why the owner sees no reply.
- **Agents are one-shot processes.** `runner.RunRoundOne` spawns participants concurrently
  in goroutines; each `runAgent` (`runner.go:275-432`) builds the command via
  `CommandFor` (`runner.go:659-685`), runs it to completion with a per-agent **timeout**
  context that is a child of the run-wide context, and exits. Per-agent stdout/stderr go
  to `<agentDir>/stdout.log` / `stderr.log` (the `StdoutPath` the TUI tails). There is NO
  long-lived interactive agent session.
- **Segments exist for re-invocation bookkeeping.** `appendSegmentStarted` / `nextSegmentID`
  (`runner.go:203-237`) stamp a `run.segment_started` (segment-NNNN, round label, target
  agents) before each round-run; `RunImplementation` / `RunReviewRound` / `RunFixup`
  (`internal/runner/phase58.go`) each open a new segment. This is the natural seam for a
  steer-driven single-agent re-invocation.
- **Cancellation is run-wide only.** `app.go` makes `runCtx, cancelRun :=
  context.WithCancel(ctx)`; the TUI's `opts.Cancel` calls `cancelRun()` (ctrl+c). There is
  **no per-agent cancel handle** — the per-agent timeout ctx inside `runAgent` is local and
  not tracked, so one agent cannot be killed while others continue.
- **HITL answers are durable but not read back** (`hitl.Answer`). `answerQuestion` in the
  TUI writes the answer + a `hitl.answered` event; agents don't poll answers mid-run.
- The TUI already has the **picker** sub-mode (1.17.0): `pickerState` on `liveModel`,
  intercepted first in `updateMain`, ↑/↓/Enter/esc/filter. Reusable for feature 1.

## Proposed direction (a STARTING proposal — challenge it in round-01)

- **Feature 1 — autocomplete.** A lightweight command-suggestion affordance in `live.go`:
  while `inputText` starts with `/`, show the matching commands (prefix match over the
  known command set: `/help /status /follow /deck /answer /open /home /quit` + any new
  ones) as a small menu above the input row. **Tab** completes the longest common prefix
  (or, if one match, the whole command); ↑/↓ + Enter pick a command. Decide: reuse
  `pickerState` (a `pickerCommand` kind) or a dedicated `suggest` sub-mode? Keep it from
  colliding with the picker and with normal typing.
- **Feature 2 — kill an agent.** Give the runner per-agent process control: hold a
  registry of per-agent `context.CancelFunc` (or `*exec.Cmd`) keyed by agent id, created
  in `runAgent`, removed on exit. Expose `KillAgent(agentID)` on the run handle, wired to
  the TUI (a key on the focused agent tab, e.g. `ctrl+k` or `K`-when-on-an-agent-tab, with
  a confirm). A killed agent is recorded (`agent.killed` event → state failed/killed) and
  the round continues with the others. Must not cancel the whole run.
- **Feature 3 — steer round-trip.** Make a steer to an agent EXECUTE, not just record.
  Because agents are one-shot, a steer spawns a fresh single-agent invocation (a new
  segment, "steer attempt") of that agent's CLI with a prompt built from the steer text +
  the agent's context (its latest artifact / recent transcript), reusing
  `CommandFor`/`runAgent`. The attempt's stdout is captured (per-attempt log) and surfaced
  in that agent's transcript tab so the **reply appears there**; keep the existing
  `steer.requested`/`steer.delivered` events and add a `steer.replied` (or
  `agent.reply`) event. Wire `submitSteer` to trigger this via a new app/runner seam
  (injected like the driver's `ImplOps`), so `internal/tui` does not import the runner.

## Round-01 focus questions (answer independently)

1. **Steer execution model (the crux).** When the owner steers a running agent: do we
   (a) spawn an immediate fresh single-agent invocation in parallel, (b) wait for the
   current attempt to finish then run the steer, or (c) queue and apply at the next round?
   What prompt/context does the steer attempt get (just the text? + last artifact? +
   transcript tail?)? How is the reply represented and surfaced? Keep it durable + simple.
2. **Per-agent kill mechanism.** Registry of cancel funcs vs `*exec.Cmd` + signal? What
   state does a killed agent land in, what event is emitted, and how does the round/driver
   treat a killed agent (skip? failed? allow a later steer/re-run)? Confirm it never kills
   the whole run and is race-safe with normal completion.
3. **Autocomplete UX + mechanism.** Reuse `pickerState` (new kind) or a dedicated suggest
   row? Exact Tab semantics (complete-common-prefix vs cycle); do we also suggest/complete
   arguments (e.g. after `/open ` fall back to the existing run/idea picker)? No collision
   with the picker, `N`, or normal typing.
4. **Reply surfacing.** Where does the agent's steer reply render — appended into the
   agent's existing transcript buffer (same `stdout.log` tail), a new per-attempt log
   shown in the same tab, or a dedicated "conversation" pane? What does the owner see
   while the steer attempt is running (a spinner / "agy is replying…")?
5. **Concurrency & safety.** Steering an agent that is CURRENTLY running its round —
   block, queue, or run a second process? Locking/segment rules so two attempts don't
   clobber the same agent dir/artifact. What happens to a steer if the run already ended?
6. **Seams & testability.** Exact new seam(s) so `internal/tui` stays decoupled from the
   runner (mirror the driver's injected `ConsensusOps`/`ImplOps`). What can be tested
   headlessly (model-driven TUI tests + runner unit tests with a fake agent command)?

## Constraints (non-negotiable)
- Reuse existing seams: `CommandFor`, `runAgent`, the segment machinery, per-agent dirs,
  the `steer` events, the new picker. Do NOT rewrite the normal round path.
- `internal/tui` must not import `internal/runner`/`internal/app`; new capabilities reach
  the TUI through injected function seams on `LiveOptions` (like `Start`/`Cancel` today).
- Keep steers durable (the `steer.requested`/`steer.delivered` audit trail stays) AND make
  them execute. Per-agent kill must never tear down the whole run.
- `--no-tui` and non-interactive paths untouched. No new third-party deps.
- Testable without a live terminal/real agents: model-driven TUI tests + runner tests with
  a fake/stub agent command; build/vet/test green.

## Non-goals
- No long-lived interactive ACP chat sessions; no multi-turn agent memory beyond the
  context we pass into a steer attempt. No protocol change. No mouse. No fuzzy-finder lib.
- Not changing what a normal round/implementation does; steering is an additive side path.

## Deliverables (in order)
1. FINAL.md: the design for all three features — the steer-execution model + new
   events/seam, the per-agent kill mechanism, the autocomplete UX, reply surfacing,
   concurrency rules, and the exact seams/signatures.
2. Implement in `internal/{tui,runner,app,steer}` (+ tests), build/vet/test green; a
   manual-smoke note. Then Phases 6-8 review to zero agreed fixes.
