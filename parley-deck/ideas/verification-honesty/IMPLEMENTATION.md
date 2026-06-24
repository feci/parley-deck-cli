---
idea: verification-honesty
status: fix-up-cycle-1
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
branch: parley-deck-cli#loop-engineering-impl
head-commit: (this commit)
---

## Summary of work

Implemented Tier 1 of the loop-engineering backlog (LE-1/2/3/4) per `FINAL.md`,
closing the three inert-gate "false-green" seams. All four items + protocol edits to
both `COOPERATION.md` copies. `go build ./...`, `go vet`, and `go test -count=1 ./...`
are green; the drift guard `TestEmbeddedDefaultMatchesLiveDeck` passes.

## Implementation plan / checklist

- [x] **LE-1 refutation-default review** — `internal/runner/phase58.go`:
  `BuildReviewPrompt` now opens with an adversarial posture ("assume the implementation
  is WRONG until you fail to break it") and the required shape gains `## Refutation
  attempts`; `ValidateReviewArtifact` requires that section. §Phase 6 normative line added.
- [x] **LE-3 model-diversity guard** — `internal/driver/transport.go`
  `ReadRequireModelDiversity`; `internal/app/driver_impl.go`
  `reviewersShareImplementerModel`/`modelOf` + `OpenReviewRound` warns (default) or, with
  `require_model_diversity: true`, escalates. §Phase 6 line added.
- [x] **LE-2 strict_gate enforcement** — `ReadStrictGate` (transport.go),
  `Config.StrictGate` (driver.go) wired at all 3 driver-construction sites (app.go),
  `ReviewStatus{StrictGateClean, ClosingReviewRound}` (impl.go) parsed in
  `driverImplOps.ReviewStatus`, consensus prompt emits the two fields under strict_gate
  (`BuildReviewConsensusPrompt`, runner.Options.StrictGate), and `advanceReview`'s
  zero-fix branch requires a certified-clean closing round + a deterministic
  `reviewRoundHasFindings` veto, bounded by `MaxFixupCycles`. §Phase 8 driver-enforcement
  note added.
- [x] **LE-4 RunChecks generalize + artifact-wins tie** — `internal/app/driver_impl.go`
  `RunChecks` resolves `checks:` → `go test` → fail-closed (code idea, no checks) →
  pass (design-only). §Phase 4/5 `checks:` field documented in the 00-prompt frontmatter doc.
- [x] Tests: `internal/driver/strict_gate_test.go` (5), `internal/runner/phase58_le_test.go`
  (3), `internal/app/driver_impl_le_test.go` (5); updated 2 existing review fixtures in
  `phase58_test.go` for the new `## Refutation attempts` requirement.
- [x] Protocol edits mirrored in `parley-deck/COOPERATION.md` and
  `internal/protocol/defaults/COOPERATION.md` (drift guard green).

## Deviations from FINAL.md

- **hermes #8 (artifact-wins) handled transitively, not by editing RunFixup.** FINAL.md
  flagged tying the "artifact-wins" fix-up override to a real check. The post-fix-up
  `RunChecks` gate (`advanceReview`) already runs after every fix-up and escalates on
  failure; making `RunChecks` fail-closed for code ideas (LE-4) therefore closes hermes #8
  at the driver level without touching `RunFixup`'s exit-code logic. Fewer moving parts,
  same guarantee — a fix-up that wrote a valid-shaped artifact but cannot be verified now
  escalates instead of passing. (Surgical-changes rationale.)

## Notes for reviewers

- **Refutation-default scope:** universal/cheap (always-on validation). strict_gate and
  RunChecks-fail-closed scale with `strict_gate`/`auto_implement` — conditional rigor was
  the consensus REJECT of uniform rigor; verify trivial/design-only ideas still close
  cheaply.
- **Deterministic scan limits:** `scanHasRealFinding` is a heuristic (severity heading
  with a non-placeholder title). It can only VETO a clean claim (fail closed), never
  auto-pass — verify it cannot be made to pass a dirty round.
- **Strict-close bound:** confirm the strict closing-round loop terminates
  (`round >= MaxFixupCycles` escalates) and that a non-strict idea is byte-for-byte
  unchanged in behavior.
- Try to break: a `checks:` command with shell metacharacters; a same-model 2-agent
  roster; a drafter that lies `strict_gate_clean: true` with a real finding on disk.

## Fix-up cycle 1
status: complete
completed: 2026-06-24

### Fixes applied
(All agreed fixes from `review/consensus.md`, found by the round-01 refutation review.)
- **F1** — `reviewRoundHasFindings` now fails CLOSED: a non-`fs.ErrNotExist` ReadDir
  error or a per-file ReadFile error vetoes (escalates) instead of auto-passing.
  (`internal/driver/impl.go`; +`errors`,`io/fs` imports.)
- **F2** — `scanHasRealFinding` is now case-insensitive (`strings.ToUpper` on the
  severity tag) and whitespace-tolerant (`### ` with any spacing before `[`).
- **F3** — an empty-title `### [SEV]` heading now counts as a finding; only the literal
  `<title>` placeholder is ignored.
- **F4** — `OpenReviewRound` now emits an `agent.model_diversity` event (via the new
  `checkModelDiversity`) in addition to the stdout warning. (`internal/app/driver_impl.go`.)
- **F5** — `ValidateReviewArtifact` now requires a real `## Findings` heading line and a
  non-empty `## Refutation attempts` section (new `hasHeadingLine`/`hasNonEmptySection`
  helpers), not a substring match. (`internal/runner/phase58.go`.)
- **F6** — model comparison uses `strings.EqualFold` (case-insensitive).
- **F7** — `DraftReviewConsensus` escalates immediately when, under `strict_gate`, the
  drafted consensus omits `strict_gate_clean`/`closing_review_round`.
- **F8** — tests added: scan case/empty-title/fail-closed (`strict_gate_test.go`);
  `checkModelDiversity` warn/escalate/silent + event (`driver_impl_le_test.go`);
  substring/empty-section validation (`phase58_le_test.go`).

### Deviations from agreed fixes
- **F6 unknown-model handling:** kept the conservative "don't fire on an unknown
  implementer model" behavior (a warn cannot assert sameness it can't see) rather than
  warning on unresolved models — documented, not a regression.
- The `sh -c` execution concern was **deferred** (not applied) per `review/consensus.md`:
  `checks:` is author-controlled kickoff input today; it becomes untrusted only with
  automated triggers, so the guard belongs to the automation-trust tier (LE-8/9).

`go build ./...`, `go vet`, `go test -count=1 ./...` green; drift guard
`TestEmbeddedDefaultMatchesLiveDeck` green.
