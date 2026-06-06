---
idea: driver-impl-phase
status: implemented
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

## Live acceptance
(appended after the run.)
