---
agent: claude
idea: tui-live-steering
round: 2
date: 2026-06-06
---

## Cross-review

Round-01 converged hard on the architecture: steer = a fresh **single-agent attempt**
whose stdout is the reply shown inline in the agent's tab; per-agent **kill** via a
runner registry of per-attempt cancel funcs (never the run-wide cancel); autocomplete
gated on `inputText` starting with `/`; everything reaches `internal/tui` through
injected `LiveOptions` seams. codex's runner design is the backbone; agy nailed the UX;
hermes nailed the keymap + concurrency. Below I resolve the open forks and lock the
contract. codex's `Handle`-based design is adopted as the spine.

### @codex
Adopting your runner design wholesale — it's the authoritative one. Specifically:
- **`Handle` with per-attempt registry** (`active map[string]AgentAttempt`, `steerBusy`,
  mutex), `KillAgent(agentID) KillResult`, `RunSteerAttempt(ctx, req) (SteerAttemptResult,
  error)`. Kill cancels ONLY that attempt's child ctx; if normal completion already
  unregistered, `KillAgent` returns `Killed:false` and synthesizes no event. Never touches
  `cancelRun`. **Consensus adopts this.**
- **Per-steer dirs are mandatory** —
  `runs/<runID>/agents/<id>/steers/<steerID>/{stdout.log,stderr.log,reply.md}`. Never reuse
  `agents/<id>/stdout.log` (it's `os.Create`d → would truncate the live transcript).
- **Steer attempt mode bypasses artifact validation** and the protocol-artifact fallback,
  and must NOT skip-if-`round-01/<agent>.md`-exists. A steer writes only `reply.md` +
  stdout; it never writes/validates a round artifact.
- **Segment-ID race is real and in scope**: `nextSegmentID` counts events unlocked. Since
  steer attempts can fire while a round runs, guard `appendSegmentStarted` for steers with
  a handle-level segment mutex (and keep the round path sequential as today). FINAL
  mandates this. `run.segment_started` for a steer carries `reason:"steer"` + `steer_id`.
- **Events**: keep `steer.requested`/`steer.delivered`; add `steer.reply_started`,
  `steer.replied` (reply_path + stdout + duration), `steer.reply_failed`. Projection learns
  the steer stdout path from these (the TUI does not guess paths).
- **Seams + the activateRun catch**: `LiveOptions.SubmitSteer`/`KillAgent` AND
  `LaunchResult.SubmitSteer`/`KillAgent`; `activateRun` MUST copy them (else Home-launched
  runs silently lose steering/kill). Great catch — it becomes a required test.
- **Scope calls I'm taking from your concerns**: (a) CLI `parley steer` stays
  record-only (no in-process handle). (b) Observational `/open` runs have no handle →
  record-only with a clear message. (c) `RunFixup` duplicates exec and doesn't use
  `runAgent`; kill/steer target the **async round `Handle` path** that `parley tui` uses —
  fix-up/implementation phases (driver-driven, headless) are explicitly OUT of scope for
  v1, noted in FINAL as a follow-up. (d) Queue depth-1; a second steer while one is
  pending/running is **rejected** ("<agent> is already replying"), not replaced.

### @agy
Adopting your UX wholesale, with two refinements:
- **ADOPT** the slim floating suggestion menu (NOT the modal picker), inline filtering as
  you type, ctrl+k + `kill <agent>? (y/N)` warn-flip confirm, the `KILLED` badge, the
  steer-input banner + "<agent> is replying…" spinner, and **reply appended inline in the
  same agent tab** (reject a separate pane/tab — you, hermes and I all agree).
- **ADOPT** your "accidental steer" guard: the `steer <id> › ` row gets a distinct colour
  (cyan) so it's never confused with plain input; the hint says Enter sends to the agent.
- **REFINE — Tab semantics**: you proposed Tab cycles through matches after completing the
  common prefix. I'm going with: Tab completes the **longest common prefix**; if that
  leaves exactly one match, complete the full command (+ trailing space if it takes args);
  if several remain, ↑/↓ + Enter pick (or keep typing) — **no Tab-cycling**. Cycling
  fights the "Tab again hands off to the picker" gesture and is less predictable. Your
  trailing-space → picker hand-off for `/open `/`/answer ` is adopted.
- **REFINE — force-steer**: your "ctrl+k to force a queued steer" is nice but couples two
  features; v1 keeps it simple — a steer on a busy agent queues (depth-1); the user can
  separately ctrl+k the round attempt if they want it to run sooner. Documented in the hint.

### @hermes
Adopting your keymap table and concurrency rules, with one override:
- **ADOPT** the conditional-Tab rule: Tab/shift+tab switch tabs ONLY when `inputText` does
  not start with `/`; when it does, Tab drives autocomplete. ←/→ remain primary tab-nav
  always (so muscle memory survives). Adopt the `confirmKillAgentID` modal sub-mode
  (blocks all other keys; y/enter confirm, n/esc cancel) — strictly modal so a stray
  Enter can't send a steer mid-confirm.
- **ADOPT** your concurrency rules: fresh segment per steer; per-agent lock; depth-1 queue;
  kill cancels only that agent's ctx; late post-run steer allowed (TUI-launched run);
  every attempt writes its own per-steer log + a `steer.replied` event.
- **OVERRIDE — reuse `pickerState` for suggest**: you proposed `pickerSuggest` as a picker
  kind. codex and agy are right that the modal picker CLEARS and OWNS `inputText` (its
  `openPicker` wipes input and treats typing as a filter), but autocomplete must keep the
  command text live and editable in the input row. So autocomplete is a **dedicated,
  non-modal `suggest` sub-mode** (its own small state: `suggest bool` + a filtered command
  list + `suggestIndex`), distinct from `picker`. It only claims Tab + (when visible) ↑/↓
  + Enter; printable keys still edit `inputText`.

## Resolved decisions (for FINAL)
1. **Autocomplete**: dedicated non-modal `suggest` sub-mode driven by a single
   `commandSpecs` table (also feeds `runCommand` + `/help` + hints). Shows while
   `inputText` starts with `/` and has no space. Tab = complete longest-common-prefix
   (single match → full cmd + space if it takes args); ↑/↓+Enter pick; trailing space on
   `/open`/`/answer` hands off to the existing picker. Tab/shift+tab switch tabs only when
   input isn't slash-prefixed; ←/→ always switch tabs.
2. **Kill**: `ctrl+k` on an agent tab whose state is running/steering → modal
   `confirmKillAgentID` ("kill <agent>? (y/N)", warn flip) → `opts.KillAgent(id)`. Runner
   `Handle` registry of per-attempt cancel funcs; `agent.killed` event; terminal
   killed/failed badge; run continues; race-safe; never run-wide cancel.
3. **Steer**: `opts.SubmitSteer` → app records via `steer.Submit` (audit trail kept) then
   `Handle.RunSteerAttempt`: fresh segment (`reason:"steer"`, handle-locked segment id),
   per-steer dir, `BuildSteerPrompt` (steer text + idea `00-prompt.md` + the agent's own
   latest artifact + ~50-line/4KB stdout tail), no artifact validation/fallback. Emits
   `steer.reply_started`/`replied`/`reply_failed`. Returns `{ID,Status,SegmentID,
   StdoutPath}`; the TUI tails `StdoutPath` into the agent's buffer behind a
   "── steer <id>: <text> ──" divider, with a replying spinner. Depth-1 queue; second →
   reject; observational/CLI → record-only.
4. **Seams**: `SubmitSteer SteerFunc` + `KillAgent KillAgentFunc` on BOTH `LiveOptions` and
   `LaunchResult`; `activateRun` copies them (tested). `internal/tui` imports neither
   runner nor app.
5. **Tests**: fake agent command (sleeper → kill; printer → reply); runner unit tests for
   `RunSteerAttempt` + `KillAgent` (+ `-race`); TUI model tests for Tab-completion, suggest
   mode, confirm-kill dispatch, `submitSteer` via a fake seam, and the activateRun copy.

## Risks (carried into FINAL)
Segment-ID race (handle-locked for steers); per-steer dirs (no log clobber); steer mode
must not write/validate round artifacts; kill must never bleed into run-wide cancel
(test: a second sleeping agent survives a kill); orphan grandchildren on kill — v1 uses
ctx cancellation (`exec.CommandContext`), with process-group signaling noted as a
follow-up; conditional-Tab muscle-memory — mitigated by the hint line.
