---
agent: hermes
idea: driver-impl-phase
round: 2
date: 2026-06-06
responding-to: [claude/round-01, codex/round-01, agy/round-01, hermes/round-01]
---

## Position changes since prior round

I concede on the single-review-consensus termination if MaxFixupCycles=3 is treated as a hard breaker; the protocol precedent (slices 1-2) favors one clean consensus. I retain the consecutive-clean preference only as a strong recommendation, not a blocker.

## Responses to others

### @claude
- Adopt your most-terminal-first Rebuild precedence and `outstanding_agreed_fixes` frontmatter DONE signal (runner already emits it).
- On safety opt-in: the idea-level `auto_implement: true` (default false) satisfies my requirement for explicit deliberate choice before code-writing agents run. It is auditable and avoids a new global CLI flag. Agree.
- On RunChecks failure: agree with your lean to escalate immediately rather than auto-fixup; non-compiling output from implementer should surface to human.
- On review/consensus.md: prefer `review-cycle` frontmatter + IMPLEMENTATION.md history over preserving prior drafts as separate files (avoids sprawl).

### @codex
- Strong alignment on ImplOps seam, RunChecks injection, clean `outstanding_agreed_fixes` contract, and `complete` written only by driver.
- On RunChecks failure behavior: codex's path (checks gate before OpenReviewRound; failure escalates) is correct — no silent fix-up loop on build breakage. Adopt.
- On review/consensus.md across cycles: agree with frontmatter + history approach.

### @agy
- Your explicit safety gate push was correct; I endorse the idea-level opt-in (`auto_implement: true`) as the minimal deliberate control that meets the `--allow-code-writes` spirit without new CLI surface.
- Clean-git-tree precondition and build/test escalation on failure both adopted.
- Roles and strict verification logic accepted.

## New concerns / questions

None — all live disagreements resolved in this round.

## Current proposal

Converged proposal (all points accepted or counter-proposed above):
- Safety: idea-level `auto_implement: true` opt-in + clean-git + RunChecks gate + no-land + MaxFixupCycles=3.
- DONE = single review consensus with `outstanding_agreed_fixes==0` + Ready/Reserved (protocol precedent).
- Rebuild most-terminal-first; idempotent disk-driven re-entry; crash invalidation.
- review/consensus.md uses `review-cycle` frontmatter.