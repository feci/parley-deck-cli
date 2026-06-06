---
idea: driver-impl-phase
status: complete
implementer: claude
started: 2026-06-06
completed: 2026-06-06
branch: parley-deck-cli#feature/driver-impl-phase
head-commit: see-branch-tip
design-pr: n/a (local-dir transport)
implementation-pr: n/a (local-dir transport)
---

## Summary of work

Extended `internal/driver` to auto-drive Parley Deck Phases 5–8 (implementation →
review → fix-up) after FINAL.md, behind `parley run --auto` / `parley continue
--auto` on a local-dir idea that opted in. Orchestration over existing runner
building blocks (`RunImplementation`/`RunReviewRound`/`RunFixup`/`RunReviewConsensus`)
plus the built-in `outstanding_agreed_fixes` DONE signal. Implements FINAL D1–D10.

## Implementation plan / checklist

- [x] `internal/driver/cursor.go` — `Rebuild` most-terminal-first precedence
      (IMPLEMENTATION status=complete→PhaseDone; review/consensus.md or review/
      round-NN→PhaseReview; IMPLEMENTATION.md→PhaseImpl; valid FINAL→PhaseFinal; …);
      `implementationStatus`, `highestReviewRound` helpers (D2).
- [x] `internal/driver/impl.go` — `ImplOps` seam + `ReviewStatus`; `advanceFinal`
      (opt-in + clean-tree → Implement), `advanceImpl` (status review-ready →
      RunChecks gate → OpenReviewRound(1)), `advanceReview` (round complete → draft
      review consensus → signoffs → done|fixup|escalate), `gitTreeClean`,
      `archiveReviewConsensus`, new Actions (D1/D3/D4/D5/D6).
- [x] `internal/driver/driver.go` — Advance routes PhaseFinal/Impl/Review; Config
      gains `Impl`/`AutoImplement`/`MaxFixupCycles` (default 3).
- [x] `internal/driver/loop.go` — log/continue past FINAL into impl/review;
      `ActionComplete` → stop ("ready to merge").
- [x] `internal/driver/transport.go` — `ReadAutoImplement` (00-prompt opt-in, D3).
- [x] `internal/app/driver_impl.go` — production `ImplOps` adapter (RunImplementation/
      RunReviewRound/RunFixup/RunReviewConsensus + review-mode consensus + RunChecks
      `go test ./...` gated on go.mod; implementer=first participant, reviewers=rest;
      `Complete` writes IMPLEMENTATION status=complete — D5).
- [x] `internal/app/app.go` — `runTask` passes `Impl`/`AutoImplement` + `--no-implement`;
      `runContinue --auto` → `continueAuto` reconstructs the driver and ticks it (D9).
- [x] `internal/driver/impl_test.go` — fake `ImplOps`, all gate branches.
- [x] `go build/vet/test ./...` green; `GOOS=windows go build ./...` green.
- [ ] Live acceptance (below).

## Safety model (D3)
Code-writing phases require `--auto` AND idea-level `auto_implement: true` (default
false), checked before Implement AND Fixup; `--no-implement` forces stop-at-FINAL.
Clean git working tree precondition (dirty → escalate; non-git workspace = clean).
No-land boundary: never merge/push/tag/release — the driver stops at driver-written
`status: complete`. RunChecks failure escalates (no auto-repair of un-reviewed
failures). MaxFixupCycles=3 breaker.

## Deviations from FINAL.md
- **Consensus archiving adopted (agy's deferred option).** D8 said "overwrite +
  review-cycle frontmatter; no archive files." In practice the draft-if-absent gate
  needs the root `review/consensus.md` cleared to re-draft for the next cycle, and
  the runner's `BuildReviewConsensusPrompt` hardcodes `review-cycle: 1`. So after a
  fix-up the driver archives the just-reviewed consensus to
  `review/round-NN/consensus.md` (agy's middle ground) and opens round N+1; the next
  tick re-drafts a fresh root consensus. This preserves the audit trail AND makes the
  re-draft idempotent. Flag for reviewers.
- `Complete` lives in the `ImplOps` adapter (a deterministic orchestrator file write,
  not an implementer agent) rather than a separate driver helper; satisfies D5
  ("complete written only by the driver, not the implementer").

## Notes for reviewers
- Driver imports only `runner`/`store`/`consensus`/`protocol` — never `internal/app`
  (ImplOps injected, like ConsensusOps).
- Idempotent re-entry: Rebuild's most-terminal-first ordering means re-entry after a
  crash routes to the furthest-along phase; advanceImpl only opens review round 1
  (later rounds come from the fix-up loop); a present review/round-NN routes to
  PhaseReview not PhaseImpl.
- The fix-up cycle number == the current review round; bounded by MaxFixupCycles.

## Live acceptance (PASS — both paths)

`parley run --auto` creates the idea + drives to FINAL; without `auto_implement` it
STOPS at FINAL ("auto-advance not enabled here; idea left at final") — the safety
opt-in confirmed. After adding `auto_implement: true` to 00-prompt.md,
`parley continue --auto <idea>` resumes.

**Happy path (`/tmp/dd-s4b`, clean self-contained task):** the driver drove FINAL →
`implementing via codex` (created TIPS.md + IMPLEMENTATION.md status=implemented) →
`opening review round 1 (reviewers: agy)` → `drafting review consensus via codex` →
`idea … is complete (review consensus clean); ready to merge — the driver does not
merge/push/release`. IMPLEMENTATION.md status=`complete`, the deliverable (TIPS.md)
present, **zero escalations** — full Phase 5-8 auto-drive to complete with real
agents and zero human input.

**Safety path (`/tmp/dd-s4`, task with no real codebase):** same chain through
implement → review → review consensus, but agy's review correctly flagged a CRITICAL
(IMPLEMENTATION.md marked implemented while the workspace has no source to modify);
the drafter set `blocked: true`, and the driver **escalated** (blocking inbox note),
never marking a hollow implementation complete. The blocked-review safety gate works.

`continue --auto` (D9) is the resume mechanism; the no-land boundary held in both
runs (the driver stops at "ready to merge").

## Fix-up cycle 1 (Phase 8)

status: complete

Applied the Phase 7 agreed fixes (review/consensus.md). hermes: no findings;
codex + agy converged on the RunChecks-after-fixup gap + hardening.
- **AF1** (CRITICAL) — `advanceReview` runs `RunChecks` after `Fixup` before opening
  the next review round; failure escalates. Tests:
  `TestPhaseReviewFixupChecksFailEscalates`, fixup test asserts `checks`.
- **AF2** — driver-owned `review/round-NN/.fixup-done` marker; a crash after Fixup
  before the next round re-enters via the marker fast-path (no re-Fixup/re-draft).
  Test: `TestPhaseReviewFixupMarkerSkipsRefixup`.
- **AF3** — review-consensus drafter is now a non-implementer reviewer (was the
  implementer); prevents the implementer filtering reviewer findings.
- **AF4** — `ReviewStatus` strips quotes from `outstanding_agreed_fixes`/`blocked`.
- **AF5** — `OpenReviewRound` (adapter) removes malformed reviewer artifacts before
  re-running so a bad file regenerates instead of spinning to the deadline.
- **AF6** — `resolveImplementer` reads IMPLEMENTATION.md/FINAL.md role metadata
  (validated vs participants), else participants[0]; reviewers = the rest. Test:
  `TestResolveImplementerFromRoleMetadata`.
- **AF7** — known in-progress IMPLEMENTATION statuses await; empty/unknown escalate.
  Test: `TestPhaseImplInProgressAwaits`.
- **AF8** — `ReadAutoImplement`/`ReadCrossReviewRounds` strip quotes.
- **AF9** — `gitTreeClean` probes `--is-inside-work-tree`; a git error INSIDE a repo
  is treated as dirty/unsafe.

Checks: `go build/vet/test ./...` green; `GOOS=windows go build ./...` green.

## Phase 8 complete — 2026-06-06

Fix-up cycle 1 re-review: codex, agy, hermes all signed ✅ ACCEPT (review/consensus.md
+ review/round-02/). Zero remaining agreed fixes. **driver-impl-phase is complete**:
`parley run --auto` / `parley continue --auto` auto-drive a local-dir idea (opted in
via `auto_implement: true`) all the way from a task to reviewed, complete code —
implementation → review → fix-up → status=complete — with the safety boundary
(opt-in + clean tree + checks gate + no-land + MaxFixupCycles) and zero human input.
Shipped as parley-deck-cli 1.16.0.
