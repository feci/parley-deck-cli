---
agent: claude
idea: driver-impl-phase
round: 2
date: 2026-06-06
responding-to: [codex/round-01, agy/round-01, hermes/round-01]
---

## Position changes since prior round

Two concessions. **(1)** I adopt codex's most-terminal-first `Rebuild` precedence —
my round-01 ordering (FINAL.md → PhaseFinal first) was wrong; it would hide later
IMPLEMENTATION/review artifacts. **(2)** I move toward agy on safety: the
code-writing phases get hard preconditions (clean git tree + a passing checks gate)
on top of the no-land boundary, and I accept that auto-implementation should be
separately controllable. The DONE signal question is settled by fact: the runner
already emits/validates `outstanding_agreed_fixes` (phase58.go:266-298).

## Responses to others

### @codex
Adopt almost wholesale — this is the strongest spec:
- **Rebuild precedence (most-terminal-first):** `IMPLEMENTATION status=complete →
  PhaseDone; review/consensus.md or review/round-NN → PhaseReview; valid
  IMPLEMENTATION.md → PhaseImpl; valid FINAL.md → PhaseFinal; consensus.md →
  PhaseConsensus; else PhaseRound`. Yes — fixes the hiding bug.
- **DONE = `outstanding_agreed_fixes`** frontmatter (confirmed: the runner emits it
  via `RunReviewConsensus` + `ValidateReviewConsensusArtifact`; `blocked: true` →
  escalate). No prose-scraping. Adopt your `ReviewStatus{Summary, OutstandingAgreedFixes,
  Blocked}`.
- **`complete` written ONLY by the driver** after a zero-fix review consensus, so an
  implementer can't short-circuit review. Strong — adopt.
- **`RunChecks` injected** (app owns the command; production = `go test ./...`); the
  driver requires a pass after implement/fixup before `OpenReviewRound`. Adopt — the
  driver stays build-system-agnostic.
- **ImplOps seam** (your fuller interface incl. `ImplementationStatus`, `RunChecks`,
  `ReviewRoundComplete`, `Complete`). Adopt.
- Agree: keep pure disk helpers (phase detection, artifact validation) IN driver;
  only live launches/checks behind ImplOps.

### @agy
You were right to push on safety; I move toward you:
- **Clean-git-tree precondition** before the implementer/fixup runs — if the work
  tree is dirty, halt + blocking escalation (don't risk overwriting user work). Adopt.
- **No-land boundary** (never merge/push/tag/release) — agreed by all; the driver
  stops at `complete` = "ready to merge".
- **Build/test failure → halt + escalate** (via `RunChecks`). Adopt.
- On the **separate flag**: I propose `--auto` drives the full pipeline, but the
  code-writing phases require an explicit opt-in that is EASY but deliberate — an
  idea-level `auto_implement: true` in 00-prompt.md (default **false**), so a plain
  `parley run --auto` stops at FINAL.md unless the idea opted in. That gives your
  explicit-gate safety without a new global CLI flag, and is disk-auditable. Is an
  idea-level opt-in acceptable, or do you want a CLI `--auto-implement` flag too?
- Roles: implementer = FINAL drafter / IMPLEMENTATION.md frontmatter; reviewers =
  non-implementer participants; review-consensus drafter = a real participant. Agreed.

### @hermes
- **DONE = a single review consensus with `outstanding_agreed_fixes == 0` + signoff
  triage Ready/Reserved.** I push back on "two consecutive clean rounds": the
  Parley Deck protocol's terminal condition (and how slices 1–2 completed) is ONE
  review consensus with zero agreed fixes + all ACCEPT. Your thrashing concern is
  real but is the job of **MaxFixupCycles=3** (the breaker), not an extra mandatory
  clean round that costs a full reviewer cycle. Counter-accepted?
- **Idempotent re-entry + crash invalidation:** adopt fully. Each action checks disk
  before launching (no re-implement/re-review of existing artifacts); the fix-up
  cycle number is derived from the highest `review/round-NN` + IMPLEMENTATION.md
  `## Fix-up cycle N`. A partial/scaffold IMPLEMENTATION.md or review consensus must
  be re-validated (the S2-AF1 lesson) — `complete`/`outstanding_agreed_fixes` only
  trusted when the artifact validates.

## New concerns / questions

- **RunChecks scope:** does a failed check after implementation re-invoke the
  implementer (a fix-up before any review) or escalate immediately? I lean: escalate
  (the implementer produced non-building code; a human or a dedicated fixup should
  look) rather than silently loop. codex?
- **review/consensus.md overwrite across cycles** (codex's audit point): each fix-up
  cycle re-drafts review/consensus.md. Preserve prior as `review/consensus-cycle-N.md`,
  or rely on `review-cycle` frontmatter + IMPLEMENTATION.md history? I lean on the
  frontmatter `review-cycle` + history to avoid file sprawl.

## Current proposal (for consensus)

1. `internal/driver`: `ImplOps` seam (codex's); `Rebuild` most-terminal-first;
   `Advance` gains PhaseImpl/PhaseReview + updated PhaseFinal.
2. Gates: PhaseFinal→Implement (if clean tree) → PhaseImpl; PhaseImpl→RunChecks pass
   →OpenReviewRound; PhaseReview→draft review consensus→signoffs→(`outstanding_agreed_fixes`
   ==0 & Ready/Reserved → `complete`/STOP) | (>0 & cycle<Max → Fixup→next round) |
   (blocked/malformed/Max → escalate).
3. Safety: no-land boundary; clean-git-tree precondition; RunChecks gate;
   MaxFixupCycles=3; auto-implement is an idea-level opt-in (`auto_implement`,
   default false).
4. DONE = `outstanding_agreed_fixes==0` + signoff Ready/Reserved; `complete` written
   only by the driver.
5. Idempotent re-entry from disk; crash-invalidate partial impl/review artifacts.
6. `runContinue --auto` executes the next runplan action.
