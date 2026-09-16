# kimi-1 — N1 APP fixture correction + independent runner-skip evaluation (2026-09-16)

Allocation: `internal/app/driver_precheck_test.go`, this note. No other file touched. Prior note
`kimi-1-n1-app-review-tests-20260916.md` preserved untouched; this note supersedes only its
description slips, not its assertions.

## Cause of the FAIL (`TestDriverPrecheckImplementationRefusalAttributionAndRestore`, line 101)

PRIMARY (file:line locators consulted directly):

- The test built its ops via `precheckImplOpsFixture` → `trajectoryHelperFixture`
  (`trajectory_verify_test.go:44`) → `trajectoryHelperFixtureScope` (`:52`) → `gateScratchRepo`
  (`driver_evidence_test.go:24`), which unconditionally writes
  `parley-deck/ideas/idea-x/IMPLEMENTATION.md` with `status: implemented`
  (`driver_evidence_test.go:40`). My test never removed it.
- Production truth: `RunImplementation` (`phase58.go:24-40`) launches `selected[0]` via
  `runAgent`, which SKIPS the launch when the artifact exists and `Overwrite=false`
  (`runner.go:406-419`) — a pure no-op that renders nothing.
- The coordinator's corrected `precheckSelectedLaunch` (`launch_precheck.go:54-74`) now mirrors
  that skip: existing `IMPLEMENTATION.md` + `!Overwrite` → `continue` → nil. So
  `PrecheckImplementation` on my fixture correctly returned nil, and the test's `err == nil`
  branch fired at line 101.
- Conclusion: the behavioral assertion was never wrong; the fixture PRECONDITION was. The test
  only passed under the old unconditional-render precheck — the exact false-refusal case my
  prior note flagged for the runner owner. The fixture was reused for its discovered
  agents/git/protocol activation; the pre-existing `IMPLEMENTATION.md` came along unnoticed.
  This is the same class of fixture-reuse trap I recorded in `kimi-1-recovery-fix-20260910.md`
  (inherited state silently contradicting the scenario under test).

## Exact scope of the correction (test preconditions only)

1. `TestDriverPrecheckImplementationRefusalAttributionAndRestore`: remove the fixture's
   `IMPLEMENTATION.md` before tampering, establishing a genuinely PENDING implementer launch —
   the only state in which production `Implement` would launch and render. Every behavioral
   assertion is byte-identical (exactly one unstarted refusal attributed
   builder/idea-x/implementation/run-id with `protocol_context_refused`; restore → clean pass,
   zero new records). No assertion weakened, no production file touched.
2. NEW `TestDriverPrecheckImplementationSkipsExistingArtifact`: the fixture's `IMPLEMENTATION.md`
   left in place (precondition asserted via `os.Stat`, so fixture drift fails loudly instead of
   silently vacating the case); tampered phase-5 cache; asserts nil error AND zero retained
   records for the run. This pins the app wiring to the runner skip contract and guards the
   false-refusal regression at the app level — meaningful because the runner-side skip test has
   no real-run cross-check for the implementation entry point (see residual observations below).

No edit to `driver_precheck.go`: the wiring remains faithful per the prior note's review.

## Correction of my prior note (description slip)

Prior note §"Tests added" item (4) says the consensus-drafter refusal is attributed
"`one-shot`/reviewer". The actual assertion (`driver_precheck_test.go:201`) expects agent
**builder** — `firstHeadlessAgent` over participants `[builder, reviewer]` selects builder, the
first participant. "reviewer" was a description slip in the note; the test assertion itself was
and is correct. Corrected record: drafter refusal attributed `one-shot`/**builder** via
`firstHeadlessAgent`.

## Independent evaluation — coordinator runner skip correction (verdict: faithful; no defect found)

Checked each precheck entry point against its production twin (PRIMARY, file:line):

- **Real options/paths.**
  - `PrecheckRound` (`launch_precheck.go:11-19`) == `RunRound` (`runner.go:925-933`): same
    round<2→2 clamp and RoundLabel default; `RunRoundOne` launches EVERY selected agent through
    `runAgent`'s per-agent skip; the precheck iterates the same list with the same predicate
    (`err==nil && !Overwrite` → skip; stat errors proceed to the render, `launch_precheck.go:64-68`
    mirroring `runner.go:406` exactly), rendering for the first pending agent only — documented
    sufficient because the renderer depends on phase/idea/root (`launch_precheck.go:50-53`).
  - `PrecheckImplementation` (`:22-29`) == `RunImplementation` (`phase58.go:24-40`): identical
    Phase/ArtifactName/RoundLabel defaults.
  - `PrecheckReviewRound` (`:32-39`) == `RunReviewRound` (`phase58.go:44-51`): same round<1→1
    clamp, Phase="review", `review/round-NN` label. `runAgent`'s review-snapshot isolation
    (`runner.go:433`) runs AFTER the skip stat (`:406`), so the canonical-path predicate is the
    operative one in both.
  - Phase attribution via `protocolLaunchPhase` (`protocol_context.go:161`): rounds →
    `roundLabel(round)`; named phases pass through — identical to the phase the real launches
    carry into `LaunchInfo`.
- **First-only implementation/drafter.** The REAL `RunImplementation` and `RunReviewConsensus`
  both hard-select `selected[0]` with NO fallback to `selected[1]` (`phase58.go:39`, `:369`).
  The precheck truncates `selected[:1]` BEFORE the skip loop (`launch_precheck.go:56-58`), so an
  existing artifact on the first implementer yields nil — it does NOT fall through to a second
  participant. Exact mirror; correct.
- **Review-consensus forced overwrite.** The real entry forces `opts.Overwrite = true`
  (`phase58.go:363`) so an existing `review/consensus.md` never skips; the precheck forces the
  same (`launch_precheck.go:46`) — pinned by the skip test's
  `review-consensus-always-overwrites` case (existing file + caller `Overwrite=false` still
  yields a first-attributed refusal).
- **Correctly untouched:** `PrecheckFixup` (`phase58.go:55-64`) — `RunFixup` (`:71+`) never
  stats the artifact; it always launches `selected[0]` through its own hardened path. No skip
  exists to mirror.
- **Skip test quality** (`launch_precheck_skip_test.go`): real fixtures; `round-all-skipped`
  additionally runs the REAL `RunRound` and asserts both results `Skipped` plus zero terminal
  records — a true production cross-check, not a precheck-side-only assertion; existing
  artifacts verified byte-untouched in every case; refusal attribution (agent/phase, unstarted)
  checked per case.

Residual observations (reported, not defects; no suppression):

- The real-run cross-check exists only for rounds; `implementation-first-only-skipped` asserts
  the precheck side only (a real `RunImplementation` would pull in `budget.GroupStepSession`
  machinery first). The implementation mirror therefore rests on the shared `runAgent`
  predicate plus my new app-level skip test — adequate coverage, recorded so it is not mistaken
  for a direct real-run proof.
- Both precheck and real entry points swallow `selectedAgents`' unresolved list identically
  (`selected, _ :=`); with zero resolved agents the precheck returns nil while the real run
  errors — safe direction (the precheck never authorizes; the real launch re-renders at
  `beginProtocolLaunch`, `protocol_context.go:116-129`).
- Multi-launch rounds: the precheck renders only for the FIRST pending agent while the real
  round launches all pending agents — documented design; the retained refusal names an agent
  that would genuinely launch.

## Precise remaining gaps

- Driver-level step-wrapper proof (precheck fires before `ChargeStep`, `budget.go:97-153`)
  still not added — needs a configured step-binding fixture (MaxDriverSteps); unchanged from
  the prior note, the coordinator owns that harness correction.
- Untested app branches unchanged: non-review `PrecheckConsensusSignoffs`,
  `PrecheckReviewLaunch`, `PrecheckReviewConsensus`, `PrecheckConsensusFinal` (same shape as
  tested siblings; recorded, not inflated).
- Preserved untouched: real-launch rechecks at `beginProtocolLaunch`, spent late refusals,
  quorum, historical unknowns, source-only status. No amendment round, migration, cap change,
  or acceptance/signoff.

Source-only correction; the coordinator runs the combined tests after terminal.
