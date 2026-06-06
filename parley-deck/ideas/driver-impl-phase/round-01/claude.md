---
agent: claude
idea: driver-impl-phase
round: 1
date: 2026-06-06
---

## Summary

The building blocks (`RunImplementation`/`RunReviewRound`/`RunFixup` in phase58.go,
plus `consensus` review-mode) already exist, so S4–S5 is orchestration + gating in
the same shape as slices 1–2: extend `Advance`/`Rebuild` with PhaseImpl/PhaseReview,
fed by an injected `ImplOps` seam. My round-01 takes a firm position on the one
genuinely new risk — **auto-driving an agent that writes code** — and proposes the
disk-derived done/termination signals.

## Proposed approach

**1. Safety boundary (the crux): auto-drive AUTHORING, never landing.**
`--auto` + local-dir is the opt-in (the user explicitly chose auto on a throwaway/
feature workspace). The risk is contained by a hard boundary: the driver
implements on the run's branch, reviews, and fixes up, then **STOPS at "ready to
merge"** (IMPLEMENTATION.md `status: complete`). It NEVER `git merge`/`push`/tags/
releases — a human does that. This mirrors slices 1–2 (the driver writes artifacts
but never merges). I do NOT think a separate `--auto-implement` flag is required,
because the no-land boundary already contains the blast radius; but I'd support an
OPTIONAL idea-level `auto_implement: false` opt-OUT for the cautious. Reviewers:
push hard on whether `--auto` alone is enough license to run a code-writing agent.

**2. Disk-derived phase + done detection.**
`Rebuild` precedence gains, above the consensus branch:
```
valid FINAL.md + IMPLEMENTATION.md present:
   review/round-NN present → PhaseReview     (gate on review completeness + consensus)
   else                    → PhaseImpl        (gate on IMPLEMENTATION.md status)
```
IMPLEMENTATION.md `status:` frontmatter is the implementation signal:
`implemented` → ready for review; `fix-up-cycle-N` → a fix-up is mid-flight;
`complete` → done (PhaseDone, STOP). "Zero agreed fixes" is the review-mode
consensus signal: reuse `consensus.Status(review=true)` — triage `Ready` (all
reviewers ACCEPT) AND an empty `## Agreed fixes` section → complete; a non-empty
agreed-fixes section → fix-up. (I propose a machine-readable convention: the review
consensus carries `agreed-fixes: N` in frontmatter so the driver never has to scrape
prose — codex/hermes, is that cleaner than parsing the section?)

**3. The fix-up loop, bounded and idempotent.**
PhaseReview gate: review round complete (all reviewer artifacts present + valid) →
draft `review/consensus.md` → request review signoffs (real agents) → if agreed
fixes == 0 and all ACCEPT → `complete` → STOP; else `ImplOps.Fixup` (re-invoke
implementer) → open review round N+1 → repeat, bounded by **MaxFixupCycles=3**
(escalate on breach). Idempotent: don't re-open a review round whose artifacts
already exist; the fix-up cycle number is derived from the highest `review/round-NN`
dir + IMPLEMENTATION.md `## Fix-up cycle N` sections.

**4. Roles.** Implementer = the FINAL drafter (protocol default; the driver reuses
the same first-participant selection as the consensus drafter, restricted to idea
participants). Reviewers = the non-implementer participants. Review-consensus
drafter = facilitator agent (same drafter, like the consensus draft).

**5. `ImplOps` seam** (app-side adapter, injected — driver never imports app):
```go
type ImplOps interface {
    Implement(ctx) error                       // RunImplementation (one implementer)
    OpenReviewRound(ctx, round int) error      // RunReviewRound (non-implementers)
    ReviewStatus() (consensus.Summary, error)  // review-mode consensus triage
    DraftReviewConsensus(ctx) error            // draft review/consensus.md (agent)
    RequestReviewSignoffs(ctx, missing) error  // real reviewer signoffs
    Fixup(ctx) error                           // RunFixup (re-invoke implementer)
}
```

## Concerns / open questions

- **Done signal robustness.** Scraping `## Agreed fixes` is fragile; I favor an
  `agreed-fixes:` frontmatter count in `review/consensus.md`. Does that change the
  consensus package contract, or can the driver own that convention?
- **Implementation "ready" detection.** Relying on the implementer agent to set
  IMPLEMENTATION.md `status: implemented` is trusting agent output. Should the
  driver also require a successful `go build`/`go test` gate before opening review?
  (A "the code at least compiles" gate before spending reviewer agents.)
- **Reviewer completeness** mirrors the round D4 gate (all non-implementer artifacts
  present + valid). Reuse the same machinery.
- **Branch reality.** RunImplementation writes "code on a branch" — in a fresh
  `parley run` workspace there may be no git repo. Does auto-implement require a git
  branch, or does it write into the workspace and leave VCS to the human?

## Risks

- **Runaway code-writing agent** is the scariest failure (it edits real files). The
  no-merge boundary + MaxFixupCycles + the build-gate (if adopted) + the existing
  driver.lock contain it, but a misbehaving implementer could still make a mess on
  the branch. The human-merge gate is the backstop.
- **Non-terminating fix-up** if reviewers keep finding fixes — MaxFixupCycles=3 then
  escalate.
- **Scope creep** — keep the ImplOps gate a small ordered switch like the round/
  consensus gates; reuse phase58 + consensus, add no new launch code.
- **Trusting agent-set status** to drive phase transitions — mitigate with the
  build/test gate and disk validation of the produced artifacts.
