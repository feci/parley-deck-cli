---
agent: claude-1
idea: meta-protocol-change-participants-optional
round: 1
date: 2026-10-31
---

## Summary

The owner wants an idea to continue while an organizer or an implementer exists. Core 2.20.0 stops it earlier:
the batch floor needs two usable non-facilitators, and the review gates need at least one independent reviewer.
I propose a new policy generation, `participant-failure-v3`, for new ideas only.

## Proposed approach

1. **Policy.** New ideas freeze `participant-failure-v3`. V3 keeps every v2 rule (no-retry exhaustion,
   succession, typed history). Saved v1 and v2 ideas keep their floor, gates, bytes and hashes. This follows the
   existing rule that an upgrade never widens a saved policy.
2. **Floor.** V3 requires one usable organizer or implementer holder after succession. A usable frozen
   facilitator counts when it holds a role. If nobody can hold a role, the batch applies nothing and escalates
   once, as today.
3. **Review.** Zero independent reviewers suffice only when validated immutable history proves that the latest
   automatic transition removed every remaining independent reviewer. I reuse the single-reviewer proof with
   `remaining = 0`; markers, manual exclusions and stale history earn nothing.
4. **Close.** The CLI driver does not auto-close an idea with no reviewer. It stops at review and hands the close
   to the organizer. The organizer closes with green checks and current-tree AC evidence and records
   `independent review: none`. A missing review is never a pass (§15).
5. **Driver.** The non-solo track check must not stop a v3 idea whose sole member is a proved survivor.
6. **Size.** Add about 640 bytes of protocol text. Reflow seven hard-wrapped §9.0 bullets without word changes so
   the phase-1 meta facilitator packet stays under 70,000 bytes. Shorten one SKILL.md paragraph to a pointer to
   references that already hold the same rules.

## Concerns / open questions

- "Participants are optional" can also mean an idea may start with no participant. I exclude that case: a
  deliberate solo start still needs the §1 user exception. Automatic exclusion at kickoff is covered.
- The driver could auto-close with the marker. I prefer an organizer close because no independent process checks
  the work; LE-11 already refuses a quiet single-checker close.

## Risks

- A v3 idea can close with no independent review. The record makes this visible, and the owner accepted the
  trade-off in the rule.
- This idea itself has no independent review (solo exception). A later idea may review it.

## Existing alternatives

- Keep v2 and ask the owner at each loss. The owner rejected this: participants are optional.
- An on-demand appendix for v3. It needs a new packet block and map change; the reflow is smaller.
