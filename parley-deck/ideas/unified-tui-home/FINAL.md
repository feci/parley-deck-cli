---
idea: unified-tui-home
status: final
author: claude
consensus-date: 2026-06-04
participants: [claude, codex, hermes]
---

## Final plan / specification

Unify the two TUIs into ONE model in `internal/tui/live.go`: the Claude-CLI-style
tabbed live view gains a **Home tab** and an optional active run, becoming the
single `parley tui`. Retire the old `internal/tui/app.go` workspace dashboard.
Unanimous consensus (claude/codex/hermes; see consensus.md). Runner/events/
`steer`/`runstate`/`hitl` contracts and `--no-tui` are reused unchanged.

### Owner requirements (locked in 00-prompt.md)
1. Default participants = all available agents; show only available agents as tabs.
2. Don't exit the TUI when the run finishes — keep it open to keep working.
3. Each agent tab shows its live generated output.
4. Start a new idea from inside the TUI via `N`.
5. Retire the old `parley tui` dashboard; the new tabbed app IS `parley tui`.

### Design

**One model, nullable run axis (D1/D2).** `liveModel` has a no-run state (Home
only) and an active-run state (Home + agent tabs + Status). Tab order:
`Home` (id `home`, first) → agent tabs → `Status`; `[Home]` only when no run.
**Tabs are derived from the active run's recorded participants** (run.created /
manifest), NOT the global roster — so unavailable agents never render (#1, #2).
`activeTabResolved` defaults to `Home` on `parley tui`, first running agent on
`parley run`.

**Default participants unchanged (D3).** `installedAgentIDs` already = all `Found`
agents except legacy gemini; KEEP it (and its test). gemini stays excluded as
roster-legacy policy. Available-agent behavior comes from run-participant tabs.

**Launch contract, no parallel engine (D4).** `StartRunFunc(ctx, StartRequest)
(StartedRun, error)`; `StartedRun` carries the real `runner.Handle` + metadata
(RunID, RunDir, Participants, Idea, Cancel). Implemented by a shared
`startRoundOne` extracted from `runTask` (runcontrol.Create + optional
auto-answerer + `runner.RunRoundOneAsync`). `activateRun(started)` resets the model
in place and re-batches the ticks (readEvents/readQuestions/eventTick/elapsedTick/
waitDone(handle.Done())).

**runToken (D5).** Tag event/question/done/tick messages with a per-run identity;
ignore stale messages after `activateRun` swaps runs.

**Stay open on done (D6, owner #2/#3).** `doneMsg` sets done + one final refresh,
NO `tea.Quit`; per-participant `DONE` badge + `[done]` status; keep ticking; reap
the runner result async (surface failures). `N` starts another idea (replaces the
active run, full reset). `/quit` or `esc` on empty Home input detaches.

**Attached vs observational runs (D7).** A Home/`N`-launched (or `parley run`) run
is attached: `Handle` + `Cancel`; `ctrl+c` cancels ONLY it; done-wait via
`handle.Done()`. An `/open`-ed (or `parley resume`) run is observational: no
cancel, no done-wait, status must not imply a live process.

**Home tab (D8, owner #4).** Lists open ideas (`protocol.ReadWorkspaceStatus`) +
recent runs (`runstate.ListRuns`: slug + DONE/active) with hint `N=new · /open
<slug>`. `N` (uppercase) → input shows `new idea ›`; Enter → `startRun` →
`activateRun` → agent tabs. `/open <slug|run-id>` opens an existing run
(`runstate.ResolveRun`, observational). Refresh recent runs on doneMsg / `N` /
`/open`. No in-Home cursor in v1 (arrows are tabs). Slash/launch errors stay
in-model (status line; active run untouched).

**Retire the old workspace TUI (D9, owner #5).** Delete `RunWorkspace`/
`WorkspaceOptions`/the old `app.go` model/`ActionRunner`/panes. Repoint
`runTUIViewWithDiscovery` → `RunLive(Home …)`. KEEP `RunInit`. Keep `--no-tui`,
`parley status`, `parley resume --no-tui` intact.

**Per-agent transcript (D10, owner #3).** `refreshBuffers` already lazily resolves
the stdout path from `agentByID` each tick, so transcripts populate once
`agent.started` lands; a missing path is a transient "no output yet", never a
permanent blank. **Phase-5 gate:** verify with a real multi-agent run; fix the
projection only if `stdout` is genuinely missing.

### Slices
1. Nullable active-run state in `live.go` (Home tab first; agent tabs from run
   participants; no empty-path reads; `doneMsg` no quit). Tests.
2. `StartRunFunc` (real `runner.Handle` + metadata), `activateRun`, `runToken`
   filtering, async result reap (RunLive attached until exit).
3. Home MVP: ideas + recent runs, `N` launch (`new idea ›`), `/open <slug|run>`
   via `runstate.ResolveRun`, errors in status line, `DONE` badges.
4. Repoint `parley tui` → `RunLive(Home)`; share the round-one launch path with
   `parley run`; preserve `RunInit`; keep `--no-tui` + `installedAgentIDs`.
5. Delete the old workspace model + `ActionRunner` after coverage proves: tabs
   from run participants, stale-tick suppression, no-quit-on-done, `/open`,
   transcript refresh after `agent.started`, no non-TUI regression. Real-run
   transcript verification.

### Invariants
- One engine: TUI-launched runs go through the existing runner path, not a second
  async abstraction.
- `--no-tui`, `parley status`, `parley resume --no-tui` unchanged. No new event
  types. Reuse `steer`/`runstate` segment plumbing/`hitl`/focus buffers/slash+key
  model. No regression of 1.13.0.

### Deferred
In-Home selection cursor (browse mode); the timed "previous run" snapshot footer;
re-adding consensus-action execution / session browser; (carried) executing
queued steers, live ACP delivery, opt-in thoughts.

### References
- Consensus: ./consensus.md (claude/codex/hermes ACCEPT, 2026-06-04)
- Rounds: ./round-01/, ./round-02/
- Builds on `tui-claude-cli-layout` (1.13.0) and `tui-interactivity-overhaul` (1.12.0).
