---
agent: claude
idea: runtime-status-resume
review-round: 2
date: 2026-05-12
reviewed-commit: bd8cd27fc934467093c1b2da83dd4c10f09c979e
responding-to: [review/round-01/claude.md, review/consensus.md]
---

## Summary
Fix-up cycle 1 cleanly addresses every consensus item. The MAJOR finding from round 1 (resume TUI rendering `status=running`) is fixed by threading `m.opts.Resume` into `displayRoundStatus` and is locked in by a regression test that asserts both the absence of `status=running` and the presence of `status=unverified`. The minor coverage gaps now have deterministic tests (incomplete/failed/idle, `status --idea`, workspace `status --json`, nonexistent resume target, started-only elapsed duration). The CLI elapsed-duration fix, the `idea "X" has no runs yet` error, the resume footer `ctrl+c` advertisement, and the `ideaForRun` empty-path fallback all land as agreed. I did not find any regressions, but I noted one small semantic gap that is adjacent to (not introduced by) the fix-up.

## Findings

### [CRITICAL]
None.

### [MAJOR]
None.

### [MINOR]
- Resume TUI header collapses `idle` and `failed` into `unverified` — internal/tui/live.go:644-655, internal/runstate/runstate.go:219-242
  `displayRoundStatus` now returns `runstate.LivenessUnverified` whenever `Resume==true` and `state.RoundStatus` is empty or `pending`. Two adjacent cases come out wrong:
  1. A run with only `run.created` (no `agent.started`) is `liveness=idle` on the CLI side (deriveLiveness in runstate.go:327-334) but `status=unverified` in the resume TUI.
  2. `ProjectEvents` does not set `RoundStatus` on `run.failed`, so a failed run resumed in the TUI also renders as `status=unverified` rather than `failed` (the CLI correctly says `outcome=failed`).
  Neither violates the explicit FINAL.md rule ("Never print unqualified 'running' after a restart"), and both are improvements over the previous `running` output, but they leave the resume header inconsistent with `parley status --run` for the same run. A small follow-up that either (a) plumbs the run-level outcome/liveness into `LiveOptions` for the header, or (b) handles `run.failed` in `ProjectEvents`, would close this. Non-blocking for this slice.

### [NIT]
- `ResolveRun` says `idea "X" has no runs yet` even if the only runs for that idea are errored — internal/runstate/runstate.go:188-195
  The loop skips runs with `run.Error != ""`, then `ideaExists` returns true, so a slug whose only run is corrupted reports "no runs yet" instead of pointing at the errored run. Not introduced by the fix-up and not in consensus, just worth recording for a future pass.

## Prior findings verification
- [MAJOR] Resume TUI header `status=running` — **fixed.** `displayRoundStatus` now takes a `resume` argument (live.go:644-655) and returns `runstate.LivenessUnverified` for the pending/empty case. `TestResumeViewHasExplicitExitPath` (live_test.go:180-207) seeds an `agent.started` event, asserts the rendered view does **not** contain `status=running`, asserts it **does** contain `status=unverified`, and also asserts the updated `q/esc/ctrl+c close resume view` footer. Non-resume code paths are unchanged, so no regression for live runs.
- [MINOR] Outcome/liveness coverage gaps — **fixed.** `TestLoadRunDerivesIncompleteOutcome`, `TestLoadRunDerivesFailedOutcome`, and `TestLoadRunDerivesIdleLiveness` (runstate_test.go:69-128) add the three previously missing branches.
- [MINOR] CLI coverage gaps — **fixed.** `TestStatusAndResumeUseRunState` now also exercises `status --idea` and workspace-level `status --json` (app_test.go:243-283), and `TestResumeReportsKnownIdeaWithNoRuns` (app_test.go:286-307) covers the nonexistent-runs-for-idea path end to end.
- [MINOR] `parley status --run` `duration=-` for a started-only agent — **fixed.** `agentDuration` now returns `time.Since(agent.StartedAt)` when the agent is `StateRunning` and `StartedAt` is set (app.go:438-448), guarded by `TestAgentDurationUsesElapsedForRunningSnapshot` (app_test.go:309-317).
- [MINOR] `ResolveRun` cannot distinguish "idea exists, no runs" from "idea unknown" — **fixed.** `ResolveRun` now calls `ideaExists` before falling back to generic error messages (runstate.go:188-200, 352-364) and `TestResolveRunReportsKnownIdeaWithNoRuns` (runstate_test.go:149-160) plus the app-level test above lock in the new error string.
- [NIT] Resume footer omits `ctrl+c` — **fixed.** Footer now reads `q/esc/ctrl+c close resume view` (live.go:225-228) and is asserted in `TestResumeViewHasExplicitExitPath`.
- [NIT] `ideaForRun` fallback path under `parley-deck/ideas/unknown` — **fixed.** `ideaForRun` now only fills `Path` when the slug is known (app.go:478-487).

## Open questions
- Was the choice to keep the resume TUI header at a constant `unverified` (rather than threading the run-level liveness/outcome through `LiveOptions`) deliberate for this slice? Either interpretation is defensible — the strong FINAL.md rule is satisfied — but the small inconsistency above is the only place the fix-up surface diverges from the rest of the projection.

## Verdict
🟢 ACCEPT — the MAJOR resume-header issue is fixed with a tight regression test, all agreed minor/nit consensus items are applied, and no regressions were introduced. The one remaining MINOR inconsistency (idle/failed both render as `unverified` in the resume TUI header) is non-blocking and a reasonable follow-up.
