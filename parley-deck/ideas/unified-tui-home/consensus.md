---
idea: unified-tui-home
drafted-by: claude
date: 2026-06-04
participants: [claude, codex, hermes]
status: accepted
---

# Consensus — unified `parley tui` (Home tab + live agent tabs)

Round-01 (independent) + round-02 (cross-review) converged with no blockers. This
is the basis for FINAL.md. It unifies the two TUIs into one model in
`internal/tui/live.go`; the runner/events/`steer`/`runstate`/`hitl` contracts and
`--no-tui` are reused unchanged.

## Agreed decisions

### D1 — One model, nullable run axis
Unify around `liveModel`. It has a no-run state (Home only) and an active-run
state (Home + agent tabs + Status). `parley tui` opens it in Home mode;
`parley run "task"` and `parley resume` open it with an active run. Same model.

### D2 — Tabs derived from RUN participants
Tab order: `Home` (id `home`, always first) → agent tabs → `Status`. No-run ⇒
`[Home]` only. Tabs come from the **active run's recorded participants**
(`run.created`/manifest), NOT the global installed roster — so unavailable
agents never render (owner #2). `activeTabResolved` defaults to `Home` on
`parley tui`, first running agent on `parley run`.

### D3 — Default participants unchanged
`installedAgentIDs` already = "all discovered `Found` agents except legacy
gemini". KEEP it; do not churn `TestDefaultParticipantSelectionSkipsLegacyGemini`.
gemini stays excluded as **roster-legacy policy** (not a hidden capability
filter); reactivation is a separate roster idea.

### D4 — Launch contract (no parallel engine)
`StartRunFunc(ctx, StartRequest) (StartedRun, error)` where `StartedRun` carries
the real `runner.Handle` + run metadata (RunID, RunDir, Participants, Idea,
Cancel). Implemented by a shared `startRoundOne` extracted from `runTask` (the
SAME path: `runcontrol.Create` + optional auto-answerer + `runner.RunRoundOneAsync`).
`activateRun(started)` resets the model in place (opts/offset/events/questions/
buffers/done/activeTab) and re-batches the ticks (`readEventsCmd`/`readQuestionsCmd`/
`eventTickCmd`/`elapsedTickCmd`/`waitDoneCmd(handle.Done())`).

### D5 — runToken (stale-tick suppression)
Tag event/question/done/tick messages with a per-run identity (`runToken`).
After `activateRun` swaps runs, stale ticks/messages from the previous run are
ignored. Essential — otherwise the old `eventTickCmd` loop keeps firing.

### D6 — Stay open on done (owner #2/#3)
`doneMsg` sets `m.run.done` + one final refresh, and does **NOT** `tea.Quit`.
Render a per-participant `DONE` badge + `[done]` status; keep ticking. Reap the
runner result asynchronously so failures surface (status/Home) before the user
detaches. `N` starts another idea (replaces the active run with a full reset).
`/quit`, or `esc` on an empty input at Home, detaches; `ctrl+c` cancels (see D7).

### D7 — Attached vs observational runs
A Home/`N`-launched (or `parley run`) run is **attached**: it has a `Handle` +
`Cancel`; `ctrl+c` cancels ONLY it; done-wait via `handle.Done()`. A run opened
via `/open` (or `parley resume`'s read view) is **observational**: no cancel, no
done-wait, status must not imply a live process. This distinction drives
cancellation, done-handling, and status text.

### D8 — Home tab
Lists open ideas (`protocol.ReadWorkspaceStatus`) + recent runs
(`runstate.ListRuns`, newest first: slug + status DONE/active) with hint text
`N=new · /open <slug>`. **`N`** (uppercase) → the input row shows `new idea ›`;
Enter → `startRun` → `activateRun` switches to agent tabs. **`/open <slug|run-id>`**
opens an existing run (`runstate.ResolveRun`, observational). Refresh the recent-
runs list on `doneMsg` / successful `N` / `/open` (explicit refresh on
transitions; no file watcher in MVP). No in-Home selection cursor in v1 (arrows
are tab switches). Slash/launch errors stay in-model (status line; active run
untouched).

### D9 — Retire the old workspace TUI
Delete `RunWorkspace`/`WorkspaceOptions`/the old `internal/tui/app.go` model/
`ActionRunner`/panes. Repoint `runTUIViewWithDiscovery` → `RunLive(Home …)`.
KEEP `RunInit` (used when the workspace is missing). Keep `--no-tui`,
`parley status`, `parley resume --no-tui` paths intact (shared-helper extraction
must preserve their behavior).

### D10 — Per-agent transcript (owner #3) — verify, don't guess
`refreshBuffers` already lazily resolves a buffer's stdout path from `agentByID`
each tick, so transcripts SHOULD populate once `agent.started` (carrying
`stdout`) lands. A missing path is a transient "no output yet", never a permanent
blank. **Phase-5 gate:** reproduce with a real multi-agent run and confirm
transcripts populate; fix the projection only if `stdout` is genuinely missing.
Refresh buffers on agent-state change; don't assume paths exist at initial tab
construction.

## Trade-offs

- Stronger `tui → runner` dependency (intentional — one engine, not two).
- `RunLive` blocks until the user exits, so `parley run` can't print
  `handle.Wait()` results synchronously; reaped async + shown in-TUI; the CLI
  summary prints on detach.
- `/open` slash instead of an in-Home cursor for MVP (arrows are tabs).

## Deferred follow-ups

- In-Home selection cursor via an explicit "browse mode" toggle.
- hermes's timed "previous run" snapshot footer when `N` replaces a finished run.
- Re-adding consensus-action execution (old `ActionRunner`) as a Home workflow;
  session browser; richer Home.
- (Carried) executing queued steers, live ACP delivery, opt-in thoughts.

## Dismissed findings

- Including gemini by default (codex round-01): withdrawn — gemini is roster-
  legacy; the run-participant-derived tabs give available-agent behavior without
  changing roster policy.

## Slices (implementation)
1. Nullable active-run state in `live.go`: Home tab first; agent tabs from active
   run participants; Status when useful; no empty-path event reads; `doneMsg`
   no longer quits. Tests.
2. `StartRunFunc` (real `runner.Handle` + metadata), `activateRun`, `runToken`
   stale-message filtering, async result reap (RunLive attached until exit).
3. Home MVP: ideas + recent runs display, `N` launch (clear `new idea ›` prefix),
   `/open <slug|run-id>` via `runstate.ResolveRun`, errors in status line,
   `DONE` badges.
4. Repoint `parley tui` → `RunLive(Home)`; share the round-one launch path with
   `parley run`; preserve `RunInit`; keep `--no-tui` + `installedAgentIDs`
   unchanged.
5. Delete the old workspace model + `ActionRunner` after coverage proves: tab
   derivation from run participants, stale-tick suppression across `activateRun`,
   no-quit-on-done, `/open`, transcript refresh after `agent.started`, no
   non-TUI regression. Real-run transcript verification.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept as drafted. Captures the round-02 convergence: one nullable-run model with
a Home tab, run-participant-derived tabs (default unchanged, gemini legacy),
StartRunFunc+activateRun+runToken via the shared launch path, stay-open-on-done
with DONE badges + async reap, attached-vs-observational runs, `/open` MVP, retire
the workspace model keeping RunInit, and the Phase-5 transcript verification gate.
No blockers.

### codex — ACCEPT (2026-06-04)
Accept. D1-D10 correctly lock the design to one nullable-run `liveModel` with a Home tab; tabs derived from run participants while keeping default participant selection unchanged and gemini legacy; shared `startRoundOne` via `StartRunFunc`/`activateRun` plus `runToken`; stay-open-on-done with DONE badges and async reap; attached vs observational runs; `/open` MVP; retiring the workspace model while keeping `RunInit`; and the Phase-5 real-run transcript gate.

<!-- hermes appends its signoff below -->

### hermes — ACCEPT (2026-06-04)
From the adversarial UX lens: Home+N+/open MVP with `new idea ›` prefix clarity is solid for discoverability; stay-open-on-done + DONE badges keeps context without abrupt quits; Phase-5 transcript verification gate prevents projection guesswork. The timed "previous run" snapshot is accepted as a deferred follow-up.
