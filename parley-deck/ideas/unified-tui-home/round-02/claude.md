---
agent: claude
idea: unified-tui-home
round: 2
date: 2026-06-04
responding-to: [codex/round-01, hermes/round-01]
---

## Position changes since prior round

- **Default-participants: NO code change needed.** I checked `installedAgentIDs`
  (`app.go:2061`): the no-`--participants` default is ALREADY "all Found agents
  except legacy gemini", and the bounded 2–4 is not applied here. So requirement
  #2 is met by the existing default + the new model showing only RUN participants
  as tabs. I withdraw any "broaden the default" change and **keep gemini excluded
  as roster-legacy** (it's not hidden-because-unavailable; it's a retired roster
  member). This also means `TestDefaultParticipantSelectionSkipsLegacyGemini`
  stays — no test churn. (Counter to @codex's "stop excluding gemini": the owner
  wants only-available, and excluding the dead legacy agent serves that.)
- **Adopting @codex's launch contract wholesale** (StartRequest/StartedRun/
  StartRunFunc + `runner.Handle` + `activateRun` + shared `startRoundOne` +
  `runToken`).
- **Adopting @hermes's done-state UX** (per-tab DONE badge + `new idea ›` prefix).

## Responses to others

### @codex
Strong agreement. Adopt: the `StartRunFunc(ctx, StartRequest) (StartedRun, error)`
returning the `runner.Handle`; the shared `startRoundOne` extracted from
`runTask` (runcontrol.Create + optional auto-answerer + RunRoundOneAsync) so the
TUI uses the SAME launch path; `activateRun` resetting opts/offset/events/
questions/buffers + re-batching the ticks; repointing
`runTUIViewWithDiscovery` → `RunLive(Home)` and deleting `RunWorkspace`/
`WorkspaceOptions`/model/`ActionRunner` while KEEPING `RunInit`; the `runToken`
to ignore stale ticks after a run-swap (this is essential — without it the old
run's `eventTickCmd` keeps firing). Two refinements:
1. **Default participants** — see above; do NOT change `installedAgentIDs`/its
   test. Keep gemini excluded.
2. **Run-selection on Home**: I agree with your `activateExistingRun`, but since
   `↑/↓` are tab switches there's no cursor on Home. Resolve: MVP exposes
   open-existing via a **`/open <slug|run-id>`** slash command (reusing
   `runstate.ResolveRun`); an in-Home `↑/↓` selection cursor (valid only in the
   no-run state where there are no agent tabs) is a nice follow-up, not v1.

### @hermes
Agree on the MVP cuts (drop ActionRunner/next-actions, session browser, in-Home
run-cursor for v1) and the done-state: on `doneMsg` set done, render `[done]` +
a per-tab `DONE` badge, keep ticking, allow `N`. On the **#3 transcript bug**:
your timing hypothesis is the likely cause — but note `refreshBuffers` already
resolves a buffer's path lazily from `agentByID` each tick, so once
`agent.started` lands the path fills. The real risk is the OLD path:
`parley run` built `LiveOptions` and opened `RunLive` correctly, so transcripts
SHOULD show — we must reproduce with a real multi-agent run and confirm
`agent.started` carries `stdout`, else fix the projection. This is a Phase-5
verification gate, not a guessed code change. On `n` lowercase: it just appends
to the input (steer/launch text); only uppercase `N` triggers new-idea.

## New concerns / questions

- **`RunLive` now blocks until exit** (@codex): `parley run "task"` stays open
  after done, so `runTask` can't print `handle.Wait()` results synchronously.
  Resolve: reap results async on `doneMsg` and surface failures in the status
  line / Home; the final CLI summary only prints after the user detaches. Lock.
- **ctrl+c semantics**: cancels only the currently-attached TUI-started run (the
  one with a `Cancel`); switching Home to inspect an older run doesn't cancel the
  live one. Lock (matches "don't quit on done").
- **Home in no-run vs active-run**: tabs `[Home]` only when no run; `[Home] +
  agents + Status` when active. `tabIDs` prepends Home; `activeTabResolved`
  defaults to Home on `parley tui`, first agent on `parley run`.

## Current proposal

Unify around `liveModel`; add a Home tab + an optional active run; one model for
`parley tui` (Home) and `parley run`/`resume` (active run). Locked:
- **Home tab** (id `home`, first): lists ideas (`protocol.ReadWorkspaceStatus`) +
  recent runs (`runstate.ListRuns`); `N` → `new idea ›` input → `startRun` →
  `activateRun` switches to agent tabs; `/open <slug|run>` opens an existing run.
- **Launch contract**: `StartRunFunc(ctx, StartRequest)(StartedRun,error)` using
  the shared `startRoundOne` (= the existing `runTask` runner path); `activateRun`
  swaps state in place + re-batches ticks; `runToken` drops stale ticks.
- **Default participants**: unchanged (all Found except legacy gemini); tabs =
  run participants ⇒ unavailable agents never render.
- **Stay open on done**: `doneMsg` sets done (no `tea.Quit`), keeps ticking;
  per-tab `DONE` badge + `[done]` status; reap results async; `N` starts another
  (replaces active run); `ctrl+c` cancels the attached run; `/quit` or `esc` at
  empty-Home detaches.
- **Retire** `RunWorkspace`/`WorkspaceOptions`/old model/`ActionRunner`; repoint
  `parley tui`→`RunLive(Home)`; keep `RunInit`; keep `--no-tui`/`status`/
  `resume --no-tui`.
- **#3**: verify transcripts populate in a real multi-agent run; fix projection
  only if `stdout` is missing.

Slices:
1. Home/no-run state + `run` lifecycle + `doneMsg`-no-quit + tab order + runToken
   + tests (tab lists per state, no empty-path reads, done doesn't quit).
2. `StartRunFunc`/`StartedRun` + `N` launch via shared `startRoundOne` +
   `activateRun` + `DONE` badges + async result reap.
3. Repoint `parley tui`→`RunLive(Home)`; delete `RunWorkspace`/model/
   `ActionRunner`; keep `RunInit`.
4. `/open <slug|run>` + Home ideas/runs listing polish.
5. Tests + real-run transcript verification + polish.
