---
agent: claude
idea: runtime-status-resume
review-round: 1
date: 2026-05-12
reviewed-commit: 4e367ae
---

## Summary
The shared `internal/runstate` projection cleanly absorbs the reducer that previously lived in the TUI, and both `parley status` (default, `--run`, `--idea`, `--json`) and `parley resume` (default TUI, `--no-tui`) are wired through it. The CLI surface honors the FINAL.md cautious wording (`unverified — last event … ago`) and the JSON surface is explicitly labelled unstable. One real regression slipped through into the resume TUI, and there are several test-coverage gaps relative to the FINAL.md test list.

## Findings

- [MAJOR] Resume TUI header renders `status=running` for non-terminal runs — internal/tui/live.go:212, internal/tui/live.go:644-652, internal/app/app.go:293-300
  `runResume` constructs `tui.LiveOptions{... Resume: true}` with no `Done` channel, so `waitDoneCmd` is never spawned, `m.done` stays `false` for the entire session, and any run that has not emitted `round.completed`/`round.incomplete` keeps `RoundStatus == "pending"`. `displayRoundStatus("pending", false)` then returns the literal string `"running"`, which the header prints verbatim. FINAL.md is explicit: "Never print unqualified 'running' after a restart … status must not imply a subprocess is alive unless a future supervised-run design can prove it." The CLI detail path correctly says `unverified — last event … ago`, but the TUI a user actually `resume`s into still asserts liveness. Suggested fix: thread `m.opts.Resume` (or a derived "post-restart" flag) into `displayRoundStatus` and render the same `unverified|idle` wording (or simply `paused`/`stopped`) instead of `running` whenever Resume is true and no terminal event has been observed. A regression test that asserts the resume header does not contain `running` should accompany the fix.

- [MINOR] FINAL.md test list is only partially covered — internal/runstate/runstate_test.go, internal/app/app_test.go
  FINAL.md §Tests lists thirteen deterministic test cases; the implementation lands three runstate cases (`completed`, `unverified`+questions, `ResolveRun by idea`) and one combined app-level case. Missing: `outcome=incomplete`, `outcome=failed`, `liveness=idle`, missing/partial `run.created` fields, `ListRuns` ordering newest first (only covered indirectly), `parley status --idea` end-to-end, `parley status --json` for the workspace overview (only `--run --json` is asserted), and the nonexistent-resume-target error path. These are cheap to add given the helpers already in `runstate_test.go` and would lock in the very semantics this slice exists to provide. None of them are blockers individually, but together they leave the most important branches in `deriveOutcome`/`deriveLiveness`/`ResolveRun` untested.

- [MINOR] `parley status --run` detail shows `duration=-` for a still-running agent — internal/app/app.go:362-368, internal/app/app.go:438-443
  `agentDuration` only returns `agent.Duration`, which is set exclusively on terminal agent events. For an `agent.started`-only state the CLI prints `state=running duration=- latest=agent.started`, while the TUI computes a live `now - StartedAt`. FINAL.md asks for "elapsed/duration" in the agent table. A conservative `now.Sub(agent.StartedAt)` (with a clear formatting hint, e.g. `elapsed=…`) would be more useful than a bare dash and would not imply process liveness any more than the existing `state=running` label does. If the team prefers to keep CLI strictly snapshot-only, a one-line comment in `agentDuration` documenting that choice would be enough.

- [MINOR] `ResolveRun` cannot distinguish "idea exists, no runs" from "idea unknown" — internal/runstate/runstate.go:175-197, internal/app/app.go:221-236
  FINAL.md: "If the idea exists but no run exists, return a clear non-zero error." Today both cases collapse to either `no runs found for "X"` (when there are zero runs in the workspace) or `no run or idea "X" found; available runs: …`. The exit code is non-zero, so the contract is technically met, but the message does not differentiate the two scenarios. Cheap improvement: cross-check `target` against `protocol.ReadWorkspaceStatus(...).Ideas` and, when the slug is a known idea with no run, return something like `idea "X" has no runs yet`. Useful for users who just ran `parley init` and an idea kickoff but no `parley run` yet.

- [NIT] Resume mode hides `ctrl+c` from the footer even though the key still quits — internal/tui/live.go:225-229, internal/tui/live.go:129-132
  The resume footer reads `Keys: j/k/tab agent  n/p question  a answer  q/esc close resume view` and omits `ctrl+c`. The `ctrl+c` handler is still wired (it just no-ops `Cancel` because it is nil and calls `tea.Quit`), so users who reflexively press `ctrl+c` will still exit cleanly — the only cost is a small surprise. Either advertise `ctrl+c close resume view` alongside `q/esc`, or strip the handler entirely.

- [NIT] `ideaForRun` synthesises a path under `parley-deck/ideas/unknown` when the slug is missing — internal/app/app.go:467-479
  When `run.created` lacks an `idea` field, `RunSummary.IdeaSlug` is `"unknown"` and `ideaForRun` builds a `Path` pointing at a directory that almost certainly does not exist. Nothing in this slice consumes that `Path`, so the only effect today is a slightly misleading struct value; worth either leaving empty or commenting that this is intentionally a best-effort fallback.

## Questions
- Is the resume TUI expected to ever stream new events from an ongoing process (e.g. someone resumes a run while a background `parley run` is still writing), or is it strictly a viewer over the durable files at the moment of open? Either interpretation is consistent with FINAL.md, but the header semantics depend on it: a pure viewer should never say "running" without a terminal event, while a "may attach to ongoing writer" mode could legitimately distinguish "stream is alive" from "stale".
- Was the choice to skip a CLI-side `--idea` end-to-end test deliberate (e.g. covered manually) or just an oversight? The runstate-level `TestResolveRunChoosesNewestRunForIdea` proves the resolver, but it does not exercise the `runStatus` codepath that prints the detail body — which is the user-visible contract from FINAL.md.

## Verdict
🟡 ACCEPT-WITH-RESERVATIONS — the MAJOR finding (resume TUI header printing `running` after a restart) is a direct contradiction of an explicit FINAL.md requirement and should be fixed before merge. The MINOR test gaps would ideally be closed in the same fix-up cycle but are not individually blocking.
