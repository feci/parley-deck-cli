---
agent: codex
idea: driver-impl-phase
round: 1
date: 2026-06-06
---
## Summary

The existing driver shape is already right for this slice: `Advance` is a small ordered switch over a disk-derived cursor, and the production runner helpers for Phase 5-8 already exist. I would extend that shape with an injected `ImplOps` seam that mirrors `ConsensusOps`, keeping `internal/driver` testable and free of agent-launch or app-layer imports.

The main design issue is not launching agents; it is defining crisp disk gates for ready-for-review, review consensus, zero-fix completion, and re-entry after partial cycles. I think review-mode `consensus.Status(..., review=true)` is necessary but not sufficient because it only tells us signoff triage; the unattended fix-up loop also needs a machine-readable agreed-fix count or equivalent outcome.

## Proposed approach

Add an `ImplOps` interface in `internal/driver`, not in runner/app:

```go
type ImplOps interface {
	Implement(ctx context.Context) error
	ImplementationStatus() (ImplementationStatus, error)
	RunChecks(ctx context.Context) (CheckResult, error)
	OpenReviewRound(ctx context.Context, round int) error
	ReviewRoundComplete(round int) (bool, error)
	DraftReviewConsensus(ctx context.Context, round int) error
	ReviewStatus() (ReviewStatus, error)
	RequestReviewSignoffs(ctx context.Context, missing []string) error
	Fixup(ctx context.Context, cycle int) error
	Complete(ctx context.Context) error
}
```

The exact structs can stay small. `ImplementationStatus` should include implementer, status, and maybe commit/branch metadata parsed from `IMPLEMENTATION.md`. `ReviewStatus` should wrap `consensus.Summary` plus `OutstandingAgreedFixes int`, `Blocked bool`, and parse errors from `review/consensus.md` frontmatter.

`Rebuild` should classify later phases in priority order from most-terminal / most-specific to earlier artifacts:

```text
valid IMPLEMENTATION.md status=complete                 => PhaseDone
review/consensus.md exists                              => PhaseReview
latest review/round-NN exists                           => PhaseReview
valid IMPLEMENTATION.md exists                          => PhaseImpl
valid FINAL.md exists                                   => PhaseFinal
consensus.md exists                                     => PhaseConsensus
otherwise                                               => PhaseRound
```

This ordering prevents a valid `FINAL.md` from hiding later implementation artifacts. It also keeps cursor state rebuildable after a crash because the latest review round number and fix-up cycle are derived from directories/frontmatter rather than persisted cursor fields.

`Advance` should then gain only two new switch cases plus an updated `PhaseFinal` branch:

```text
PhaseFinal:
  if no valid IMPLEMENTATION.md, call ImplOps.Implement
  otherwise move to PhaseImpl

PhaseImpl:
  validate IMPLEMENTATION.md status is review-ready, not merely present
  run a local build/test gate if configured
  open the next review round only if its artifacts are not already present

PhaseReview:
  if latest review round incomplete, await
  if no review/consensus.md for current cycle, draft it
  if review signoffs missing, request them once and re-check
  if BLOCK/malformed, escalate
  if outstanding_agreed_fixes == 0 and all signoffs accept/reserve, mark complete
  if outstanding_agreed_fixes > 0 and cycle < MaxFixupCycles, run Fixup
  if cycle exhausted, escalate
```

For DONE detection, I would not scrape the `## Agreed fixes` body. The runner's `BuildReviewConsensusPrompt` and `ValidateReviewConsensusArtifact` already point toward a better contract: `outstanding_agreed_fixes: <int>` in `review/consensus.md` frontmatter. I would keep that field, validate that it parses as a non-negative integer, and treat `blocked: true` as escalation regardless of signoff triage. Review-mode `consensus.Status` still handles participant signoff completeness and ACCEPT/BLOCK semantics, but it cannot alone distinguish "all reviewers signed off on three agreed fixes" from "all reviewers signed off and no fixes remain."

For implementation readiness, use `IMPLEMENTATION.md` frontmatter rather than prose. The current prompt writes `status: implemented`; I would make the driver accept a narrow ready set such as `implemented` or `ready-for-review`, reject empty/malformed values, and reserve `fix-up-cycle-N` or `fixup-ready` for post-fixup work that is again ready for review. `complete` should be written only by the driver after review consensus reaches zero agreed fixes, so implementers cannot accidentally short-circuit review.

For idempotent fix-up re-entry, derive the next review round from disk:

```text
latestReviewRound = highest review/round-NN directory
currentCycle = latestReviewRound
if review/round-NN exists but is incomplete: await
if review/round-NN complete and review/consensus.md missing: draft consensus
if consensus has fixes and fixup for cycle N already recorded in IMPLEMENTATION.md:
  open review/round-(N+1) if absent
else:
  run Fixup(N)
```

This avoids reopening or rerunning an existing review round. If `RunFixup` succeeds but cursor save fails, the next tick sees the implementation update and opens the next review round. If review artifacts already exist, it does not dispatch reviewers again.

I prefer a build/test gate before spending reviewer agents. It should be an injected operation or command policy owned by the app layer, not hard-coded `go test ./...` inside the driver. The driver should require a passing check after implementation and after each fix-up before `OpenReviewRound`; if no check command is configured or discoverable, it may warn and proceed only if the implementation artifact explicitly records checks run. For this Go module, the production default can reasonably run `go test ./...` before review, but the driver seam should expose only `RunChecks`.

Tests should use a fake `ImplOps`, following the existing `fakeConsensus` pattern:

- `Rebuild` classifies valid `FINAL.md`, valid `IMPLEMENTATION.md`, review rounds, review consensus, and complete status in the right order.
- `PhaseFinal` calls `Implement` once, and does not call it when `IMPLEMENTATION.md` already exists.
- `PhaseImpl` awaits non-ready or malformed implementation status.
- `PhaseImpl` runs checks before `OpenReviewRound`; failed checks escalate and do not spend reviewer agents.
- `PhaseImpl` does not reopen review round 1 when `review/round-01` artifacts already exist.
- `PhaseReview` awaits incomplete review artifacts.
- `PhaseReview` drafts review consensus once, then requests missing signoffs and re-checks.
- `PhaseReview` marks done when signoff triage is ready/reserved and `outstanding_agreed_fixes == 0`.
- `PhaseReview` runs `Fixup` when agreed fixes are positive, then opens the next review round on re-entry.
- `PhaseReview` escalates on BLOCK, malformed consensus, malformed fix count, or `MaxFixupCycles`.

## Concerns / open questions

The proposed `ImplOps` should not become a grab bag of protocol parsing that belongs in `internal/driver`. I would keep pure disk helpers in driver for phase detection and artifact validation, and keep live agent launch/check execution behind `ImplOps`.

There is a naming mismatch to resolve: runner currently emits `outstanding_agreed_fixes`, while the prompt asks whether an `agreed-fixes:` count is cleaner. I prefer keeping `outstanding_agreed_fixes` because it encodes the decision the loop needs and avoids ambiguity between total findings ever agreed and fixes still pending.

The review consensus file is currently a single `review/consensus.md` overwritten by `RunReviewConsensus`. That is acceptable for the machine gate, but it weakens auditability across fix-up cycles unless prior drafts are preserved or the file records `review-cycle`. If the loop overwrites, tests should assert that cycle metadata changes and signoffs are recollected for the current review round.

Implementer identity should be disk-derived from `FINAL.md` or consensus metadata when available, otherwise use the protocol default FINAL drafter. Reviewers are all participants except the implementer. Review consensus drafter can be a non-implementer or configured drafter, but it must still be a real participant action in production.

## Risks

The largest risk is false DONE: accepting all review signoffs without also checking an explicit zero-fix signal. That would skip required fix-up work.

The second risk is duplicate spending of agent runs after a crash. The driver must always check for existing implementation, review round artifacts, and review consensus before invoking agents.

Another risk is making the driver run project-specific shell commands directly. Keeping the check gate behind `ImplOps` preserves unit-testability and avoids turning the driver into a build-system detector.

Finally, safety expectations are higher once agents can write code. `--auto` plus local-dir is a reasonable authoring opt-in only if the driver never merges, pushes, or releases, and if every transition remains disk-auditable with a hard `MaxFixupCycles` breaker.
