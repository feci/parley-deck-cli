---
idea: driver-impl-phase
status: final
drafted-by: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
implementer: claude
consensus: parley-deck/ideas/driver-impl-phase/consensus.md
---

## Final plan / specification

Extend `internal/driver` (shipped 1.15.0) to auto-drive Parley Deck Phases 5–8
(implementation → review → fix-up) after FINAL.md, behind `parley run --auto` /
`parley continue --auto` on a local-dir idea that has opted in. Orchestration over
existing runner building blocks (`RunImplementation`/`RunReviewRound`/`RunFixup`/
`RunReviewConsensus`) + the built-in `outstanding_agreed_fixes` DONE signal.
Implements consensus D1–D10.

### ImplOps seam (`internal/driver/impl.go`) — D1
```go
type ImplOps interface {
    Implement(ctx context.Context) error            // RunImplementation (one implementer)
    ImplementationStatus() (string, error)          // IMPLEMENTATION.md status frontmatter
    RunChecks(ctx context.Context) (bool, string)   // build/test gate (app: go test ./...)
    OpenReviewRound(ctx context.Context, round int) error // RunReviewRound (non-implementers)
    ReviewRoundComplete(round int) (bool, error)    // all reviewer artifacts present+valid
    DraftReviewConsensus(ctx context.Context, round int) error // RunReviewConsensus
    ReviewStatus() (ReviewStatus, error)            // {Summary, OutstandingAgreedFixes, Blocked}
    RequestReviewSignoffs(ctx context.Context, missing []string) error
    Fixup(ctx context.Context, cycle int) error     // RunFixup (re-invoke implementer)
    Complete(ctx context.Context) error             // driver writes IMPLEMENTATION status=complete
}
type ReviewStatus struct { Summary consensus.Summary; OutstandingAgreedFixes int; Blocked bool }
```
Driver owns pure disk helpers (phase detection, artifact validation via
`runner.Validate*`); only live launches/checks are behind ImplOps. The driver never
imports `internal/app`; the app injects the adapter.

### Rebuild precedence (most-terminal-first) — D2
In `cursor.go`, ABOVE the existing FINAL/consensus/round cases:
```
valid IMPLEMENTATION.md status=complete            → PhaseDone
review/consensus.md OR highest review/round-NN     → PhaseReview
valid IMPLEMENTATION.md present                     → PhaseImpl
valid FINAL.md (existing)                           → PhaseFinal
consensus.md (existing)                             → PhaseConsensus
else                                                → PhaseRound
```
New phases: `PhaseImpl`, `PhaseReview`, `PhaseDone`. Round/cycle numbers derived
from disk.

### Advance gates (`impl.go`: advanceImpl / advanceReview; driver.go routes) — D6
- **PhaseFinal**: if no IMPLEMENTATION.md → require `auto_implement` (idea opt-in) +
  `!--no-implement` + clean git tree → `Implement`; else surface-only (opt-in absent)
  or escalate (dirty tree). If IMPLEMENTATION.md exists → fall through (Rebuild moves
  to PhaseImpl next tick).
- **PhaseImpl**: validate IMPLEMENTATION.md status is review-ready
  (`implemented`/`ready-for-review`); malformed/empty → escalate; not-ready → await.
  Run `RunChecks`; fail → escalate (no implicit fix-up). Pass + `review/round-NN`
  absent → `OpenReviewRound(nextReviewRound)`.
- **PhaseReview**: latest review round incomplete → await. Round complete + no
  review/consensus.md → `DraftReviewConsensus`. Consensus signoffs missing →
  `RequestReviewSignoffs` then re-check; still missing → escalate. `ReviewStatus`:
  `blocked`/malformed → escalate; Ready/Reserved & `outstanding_agreed_fixes==0` →
  `Complete` (driver writes status=complete) → PhaseDone; `>0` & cycle<MaxFixupCycles
  → `Fixup(cycle)` → (next tick) OpenReviewRound(round+1); cycle≥MaxFixupCycles(3) →
  escalate.
- **PhaseDone** → surface-only ("ready to merge"; human lands it).

### Safety model — D3/D4
- Code-writing phases require `--auto` AND idea `auto_implement: true` (default
  false), checked before BOTH Implement and Fixup. `--no-implement` CLI flag forces
  stop-at-FINAL.
- Clean git working tree precondition before Implement/Fixup (driver helper
  `gitTreeClean(root)`; dirty → escalate). Fresh non-git workspaces: treat "not a git
  repo" as clean (nothing to clobber).
- No-land boundary: never merge/push/tag/release. Stop at driver-written
  `status: complete`.
- `RunChecks` injected; runs after implement + each fixup before review; fail →
  immediate blocking escalation with the result text.

### DONE — D5
`outstanding_agreed_fixes` + `blocked` read from `review/consensus.md` frontmatter
(`protocol.ReadFrontmatter`; runner emits + `ValidateReviewConsensusArtifact`
requires it). `complete` written only by the driver.

### Config additions
`Config` gains `Impl ImplOps`, `AutoImplement bool` (from 00-prompt `auto_implement`
+ `!--no-implement`), `MaxFixupCycles int` (default 3). `Rebuild` reads
`auto_implement` for nothing (it's a gate, not a phase); the driver reads it from
00-prompt via a helper `ReadAutoImplement(ideaDir)`.

## File-by-file plan
1. `internal/driver/cursor.go` — Rebuild precedence + PhaseImpl/PhaseReview/PhaseDone;
   helpers: `implementationStatus`, `highestReviewRound`, `ReadAutoImplement`,
   `validImplementation`.
2. `internal/driver/impl.go` — `ImplOps`, `ReviewStatus`, `advanceImpl`,
   `advanceReview`, `gitTreeClean`, review-completeness + consensus-frontmatter
   helpers, new Actions (ActionImplemented, ActionReviewOpened, ActionReviewDrafted,
   ActionFixup, ActionComplete).
3. `internal/driver/driver.go` — Advance routes PhaseFinal(impl)/PhaseImpl/PhaseReview;
   Config fields; New defaults (MaxFixupCycles=3).
4. `internal/driver/loop.go` — log + deadline-reset for the new progress Actions;
   PhaseDone → return.
5. `internal/app/driver_impl.go` — production `ImplOps` adapter (RunImplementation/
   RunReviewRound/RunFixup/RunReviewConsensus + request-review-signoffs + RunChecks=
   `go test ./...` + Complete writes status). Implementer = FINAL drafter/first
   participant; reviewers = non-implementers.
6. `internal/app/app.go` — `runTask` passes `Impl` + `AutoImplement` (read
   `auto_implement` & `--no-implement`); `runContinue --auto` executes next action (D9).
7. Tests: `internal/driver/impl_test.go` — fake ImplOps, all gate branches (codex's
   list); Rebuild precedence tests. Live acceptance: a tiny `auto_implement: true`
   idea driven FINAL → implementation → review → complete.

## Acceptance
A live `parley run --auto --no-tui` on a fresh local-dir idea with
`auto_implement: true` drives a tiny task to `IMPLEMENTATION.md status: complete`
with a clean review consensus (`outstanding_agreed_fixes: 0`), running real
implement + review agents, never merging/pushing, stopping at "ready to merge".
Without `auto_implement` (or with `--no-implement`) it stops at FINAL.md.

## Non-goals
No merge/push/release. No github-pr/gitlab-mr auto-drive. No new storage. discover.go
untouched. Not unifying with the §12 pipeline.

## Decisions reference
Full rationale + dismissed alternatives in `consensus.md` (D1–D10). round-01/round-02
artifacts under `round-01/`, `round-02/`.
