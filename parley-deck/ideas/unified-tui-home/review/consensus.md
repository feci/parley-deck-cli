---
idea: unified-tui-home
review-cycle: 1
drafted-by: claude
date: 2026-06-04
reviewed-commit: 16dc034
---

Synthesis of Phase 6 review-round-01 from codex and hermes (review/round-01/).
Agreed fixes applied in fix-up cycle 1 on branch feat/unified-tui-home.

## Agreed fixes

### AF1 — TUI detach must not cancel N-launched runs (codex MAJOR)
`runTUIViewWithDiscovery` tracks every `newLaunchFunc` cancel and unconditionally
cancels them all in a defer after `RunLive` returns, so `/quit`/`esc` cancels an
`N`-launched run even though only `ctrl+c` (the attached run's `Cancel`) should.
Fix: remove the unconditional defer-cancel-all; reap each launched handle in a
goroutine (`handle.Wait()` → `registerWorkspaceSessions`) so its session is
recorded; the active run's `Cancel` (ctrl+c) stays the only cancel path.

**Cycle-2 deepening (codex round-02 MAJOR):** removing the defer was not enough —
launched runs derive from the top-level signal context whose `defer cancel()`
fires when `parley tui` returns, so a normal `/quit` still killed them via context
propagation. Corrected: a `launchReaper` tracks each in-flight launched run; after
`RunLive` returns, `runTUIViewWithDiscovery` calls `reaper.waitForActive(stdout)`
to wait for them to finish (recording sessions) **before** the command returns and
the parent cancel fires. The parent context stays live during the wait, so a real
ctrl+c (SIGINT, after the TUI releases the terminal) still aborts. The attached
run's `Cancel` (ctrl+c inside the TUI) is unchanged. Covered by
`TestLaunchReaperWaitsForInFlightRuns`.

### AF2 — `parley run` should not own secondary N-launched runs (codex MAJOR)
`runTask` passes `Start` into the live TUI, but after `RunLive` exits it waits/
reports only the original handle; a secondary `N`-launched run is never reaped
and the top-level ctx can cancel it without a result. Fix (accepted fallback):
**disable `Start` in `parley run`** — `N`/new-idea is a `parley tui` feature.
Keep `Root` so the Home tab still lists ideas/runs; `N` reports "start new ideas
from `parley tui`".

### AF3 — Integration test that transcripts populate from real on-disk run data (hermes MAJOR / D10)
Add a test that writes a run dir with `events.jsonl` (`agent.started` carrying a
`stdout` path) + a `stdout.log`, builds the model on that RunDir, drives an
`eventsMsg` read, and asserts the active agent's transcript buffer is non-empty
and the rendered view shows the output. Closes the owner's #3 verification
end-to-end (projection → buffer → render), not just a mocked event.

### AF4 — Done-state exit hint (hermes MINOR)
After `done`, the status line / input hint must clearly show the exit options so
the finished run never feels "stuck": show `[done] /quit or esc to exit` (and the
DONE/FIN per-tab badge already present).

### AF5 — Delete the dead workspace model (codex MINOR/NIT + hermes NIT)
Retire the old TUI per FINAL D9 and the owner's "ten stary tui zahod" directive.
Done: `internal/tui/app.go` is reduced to the first-run init wizard plus the
shared styles/helpers `live.go` uses; deleted `Run`, `RunWorkspace`,
`WorkspaceOptions`, the workspace `model` and all its methods/panes,
`focusZone`/`focusIdeas`/`focusActions`/`focusAgents`,
`refreshRunsMsg`/`startRunMsg`/`actionRunMsg`/`refreshTickMsg`, `StartRequest`,
`StartRunFunc`, `ActionRequest`, `ActionResult`, `ActionRunner`, `newModel`,
`startRunCmd`, `actionRunCmd`, `upsertRunSummary`, and the workspace-only string/
index helpers. From `internal/app/app.go`: removed `runTUIAction`,
`consensusActionArgs`, `commandOutput` (only-caller was `runTUIAction`), and
`applySessionLaunchOverrides`. The init wizard no longer hands off to the dead
dashboard — on successful init it quits and `runTUIViewWithDiscovery` re-reads
status and opens the unified Home. Dead tests removed: the workspace-dashboard
tests in `internal/tui/app_test.go` (init-wizard tests kept, one updated for the
quit-handoff) and the TUI-action tests in `internal/app/app_test.go`. KEPT the
shared helpers `live.go` uses: `valueOr`, `truncateText`, and the shared styles
(`headerStyle`/`boxStyle`/`mutedStyle`/`okStyle`/`warnStyle`); `clipLines`,
`tuiWidth`, `tuiHeight`, `clampInt`, `sectionTitle` already live in `layout.go`.
`tui/app.go` went 1443 → 150 lines. `go build/vet/test ./...` all green.

### AF6 — Conditional ctrl+c help wording (codex NIT)
The help overlay always says `ctrl+c cancel the run`, but at Home/no-run and for
`/open` observational runs `ctrl+c` just quits. Fix: word it
`ctrl+c cancel attached run / quit`.

## Deferred follow-ups

- The timed "previous run" snapshot footer when `N` replaces a finished run
  (hermes): deferred — finished runs persist on disk, appear in Home recent-runs,
  and are reopenable via `/open`, so context is not lost.
- A shared app-side launch manager enabling secondary-run ownership inside
  `parley run` (codex): deferred with AF2.
- (Carried) executing queued steers, live ACP delivery, opt-in thoughts;
  in-Home selection cursor.

## Dismissed findings

- `n` (lowercase) alias for new idea (hermes): consensus chose `N` uppercase-only
  so a capital N never collides with steer/answer text; lowercase `n` is input.
- N-replace confirmation (hermes): N replacing a run is the locked behavior;
  mitigated by on-disk persistence + Home recent-runs + `/open`.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept the synthesis. AF1/AF2 fix the run-lifecycle ownership (detach ≠ cancel;
`parley run` doesn't own secondary launches); AF3 closes the #3 transcript gate
with a real on-disk test; AF4 removes the "stuck" feel; AF5 finishes retiring the
workspace model (codex's exact symbol list); AF6 corrects the help wording. The
deferrals (snapshot footer, secondary-run ownership in `parley run`) are the right
scope boundary. No blockers.

### codex — CHANGES REQUESTED (2026-06-04)
AF2-AF6 are correctly applied and the requested build/test command is green. AF1 still has a lifecycle gap: `N`-launched runs inherit the top-level CLI context, so normal `parley tui` detach cancels or abandons them when the command returns.

### hermes — ACCEPT (2026-06-04)
AF1–AF6 all correctly applied (background reap, no Start in `parley run`, real on-disk transcript test, done exit hint, workspace model retired, help wording). Build/tests green.

## Fix-up cycle 2 — re-review (AF1 deepening)

codex's round-02 raised a valid MAJOR: removing the defer-cancel-all left
`N`-launched runs tied to the top-level signal context, so a normal `/quit` still
canceled them when the command returned. Fix-up cycle 2 (commit pending) adds the
`launchReaper` wait-on-exit described in the AF1 cycle-2 deepening above, with
`TestLaunchReaperWaitsForInFlightRuns`. Re-review artifacts in
`review/round-03/`.

### codex — signoff after cycle-2 re-review
### codex — ACCEPT (2026-06-04)
The AF1 wait-on-exit fix resolves the cycle-1 gap: TUI detach now waits for in-flight N-launched runs before the parent cancel fires, while SIGINT and attached-run cancel behavior remain correct. The requested build/test command is green.

### hermes — signoff after cycle-2 re-review
### hermes — ACCEPT (2026-06-04)
AF1 wait-on-exit fix (launchReaper + waitForActive) is sound; in-flight N-launched runs now finish and record before parent cancel, real ctrl+c still works, test passes, build/tests green. AF2–AF6 untouched.
