---
idea: verification-honesty
status: implemented
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
