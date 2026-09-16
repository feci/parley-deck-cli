# kimi-1 — N1 APP wiring review + behavior tests (2026-09-16)

Allocation: `internal/app/driver_precheck.go`, `internal/app/driver_precheck_test.go`, this note. No other file touched.

## App wiring review — verdict: faithful, no corrections made

Compared every production precheck against the ACTUAL launch path it mirrors:

- `PrecheckImplementation` == `Implement` (`driver_impl.go:192`): same `withParticipants(o.implementer)`, runner precheck uses `RunImplementation`'s own selection (`selectedAgents`), phase `implementation`, artifact `IMPLEMENTATION.md`.
- `PrecheckReviewLaunch(round)` == `OpenReviewRound` (`driver_impl.go:253`): same deduped/capped `o.reviewers`, same `opts.Round`. The actual model-diversity check and AF5 artifact sweep run before any launch and launch nothing, so the precheck correctly ignores them (nil there is not a false pass — the real gate still fires).
- `PrecheckReviewConsensus(round)` == `DraftReviewConsensus` (`driver_impl.go:304`): same drafter + round. `StrictGate` only changes the prompt, never the rendered phase/idea/track, so omitting it is render-neutral.
- `PrecheckReviewSignoffs`/`PrecheckConsensusSignoffs` == `requestConsensusSignoffs` (`consensus_request_signoffs.go:82`): same `consensus.Status`, same malformed gate, same `requestSignoffTargets(strings.Join(missing,","))`, same `discoverConfigured` + `requestSignoffAgents(rosterMappingFor)`. Differences are all safe-direction: no TriageBlocked check (real call errors pre-launch anyway), no mode overrides/`--yes` (driver passes none/`Yes:true`), non-headless skipped (interactive/manual handoffs render no protocol in driver context: spawn-tty fails `isTerminal` before `RunInteractive`; print-only/manual only write handoffs). Idea/phase match `signoffContext(summary.Path)` because the path is built from the same slug.
- `PrecheckGoalCheck` == `GoalCheck` (`driver_impl.go:379`): identical checker resolution (`ResolveParticipant(drafter, base.Agents, rosterMappingFor)`), identical `LaunchInfo{RunID, Idea, "goal-check", Store}`. Empty/self/unresolvable checker returns nil — matches the real call's no-launch error branches.
- `PrecheckConsensusDraft/Final` == `runDrafter` (`driver_consensus.go:83`): same `firstHeadlessAgent` selection; empty phase/idea + `RunID:"one-shot"` matches `trackedCommandFor` → `beginProtocolLaunch(ctx, root, "one-shot", ...)` (`launch.go:43`) with no `WithLaunchInfo` upstream. Confirmed `Draft` ALWAYS calls `runDrafter` even when `consensus.md` exists (scaffold branch is create-if-absent only), so the unconditional precheck render is correct — no false no-op inferred.

No false-precheck-refusal case found in APP code; no edit made to `driver_precheck.go`.

## Runner concern (out of my edit scope) — REAL

`precheckFirstSelected` (`runner/launch_precheck.go:54`) renders for `selected[0]` unconditionally, but `runAgent` (`runner/runner.go:406`) skips a selected agent whose artifact exists when `Overwrite=false` — launching nothing and rendering nothing. With every selected artifact present and a tampered cache, the precheck refuses while the real run is a pure no-op: a false precheck refusal. Correction for the runner owner: mirror the skip — per entry point, with that entry point's actual `Overwrite` (`RunImplementation`/rounds/`RunReviewRound`: false; `RunReviewConsensus` forces true at `phase58.go:363`), pick the first selected agent whose artifact is absent; if none would launch, return nil.

## Tests added (`driver_precheck_test.go`)

Real fixtures (`trajectoryHelperFixture`, `refuseAppProtocolCache`), no real providers: (1) implementation precheck — tampered phase-5 cache ⇒ exactly one unstarted (`StartedAt==nil`, `protocol_context_refused`) refusal attributed to builder/idea-x/implementation/run-id; restore ⇒ success with zero new records (no launch/charge). (2) goal-check — tampered phase-(-1) cache ⇒ attributed reviewer/goal-check refusal; self-checker and unknown-checker branches stay nil and record nothing; restore ⇒ clean. (3) review-signoff precheck over a real `review/consensus.md` — tampered phase-7 cache ⇒ `review-consensus`/reviewer refusal; restore ⇒ clean. (4) consensus drafter precheck — tampered -1 cache ⇒ refusal attributed `one-shot`/reviewer via `firstHeadlessAgent`.

## Precise remaining gaps

- Driver-level step-wrapper proof (precheck fires before `ChargeStep`, `budget.go:97-153`) not added: needs a configured step-binding fixture (MaxDriverSteps); coordinator is correcting that harness in its own snapshot.
- Untested app branches: non-review `PrecheckConsensusSignoffs`, `PrecheckReviewLaunch`, `PrecheckReviewConsensus`, `PrecheckConsensusFinal` (same shape as tested siblings; recorded, not inflated).
- Preserved untouched: real-launch rechecks at `beginProtocolLaunch`, spent late refusals, quorum, historical unknowns, source-only status. No amendment round, migration, cap change, or acceptance/signoff.
