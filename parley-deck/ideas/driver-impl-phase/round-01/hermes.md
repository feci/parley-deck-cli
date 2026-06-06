---
agent: hermes
idea: driver-impl-phase
round: 1
date: 2026-06-06
---
## Summary
The fix-up loop (review → fixup → re-review) is the second back-edge after the consensus BLOCK→reopen. Bounding by MaxFixupCycles=3 + escalate is insufficient for faithful convergence; it needs a "consecutive-clean" / no-new-fixes termination (loop-until-dry) to guarantee termination on stable zero-fixes state rather than arbitrary cycle count.

## Proposed approach
- DONE detection: "zero agreed fixes" in review/consensus.md is the terminal signal (robust recompute_ready), but requires explicit "consecutive clean rounds" counter (e.g., two successive zero-fix reviews) to distinguish genuine convergence from thrashing reviewer cycles that re-produce the same fixes.
- dispatch_once for implement/review/fixup: each action must be idempotent on re-entry by checking disk artifacts (IMPLEMENTATION.md presence/status, review/round-NN/<agent>.md existence, consensus.md triage) before invoking RunImplementation/RunReviewRound/RunFixup — faithful to driver pattern in Advance/Rebuild.
- Crash stranding: impl/review back-edge can leave partial IMPLEMENTATION.md or review/consensus.md (analogous to stale consensus.md/FINAL.md); needs explicit invalidation on re-entry (e.g., status=fix-up-cycle-N guard + scaffold detection) before advancing.

## Concerns / open questions
- How does Rebuild derive PhaseImpl/PhaseReview purely from disk (status frontmatter + artifact presence) without stranding on partial writes?
- Is implementer always the FINAL drafter, or configurable?

## Risks
- Non-convergent thrashing if reviewer produces cyclic fixes without consecutive-clean check.
- Stale state on crash mid-fixup mirroring the BLOCK→reopen stale-file issue.