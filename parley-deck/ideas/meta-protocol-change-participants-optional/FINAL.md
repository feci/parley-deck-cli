---
idea: meta-protocol-change-participants-optional
status: final
author: claude-1
implementer: claude-1
consensus-date: 2026-10-31
participants: [claude-1]
---

## Final plan / specification

Solo idea under the recorded owner exception. FINAL freezes consensus decisions D1–D6.

### Protocol (all three COOPERATION.md copies and staged core 2.21.0)

- §1 non-solo: a §9.0 v3 automatic exclusion is a standing owner exception.
- §4 Phase 6: §9.0 v3 is the only exception to "no reviewer, no completion".
- §9 quorum: the last-non-facilitator rule and the one-participant rule name the v3 exception.
- §9.0 Policy: new ideas freeze `participant-failure-v3`; v3 keeps every v2 rule; saved v1/v2 stay.
- §9.0 Batch and floor: v3 needs one usable organizer or implementer holder; participants are optional.
- §9.0 Review gates: the v3 zero-reviewer proof, the organizer close, `independent review: none`, never a pass.
- Reflow seven hard-wrapped §9.0 bullets without word changes (packet guardrail).

### CLI 1.57.0

- `internal/quota`: `ParticipantFailureV3`, `NewOptionalParticipantPolicy`, `OptionalParticipants`,
  `MinUsableSurvivors`, `MinMembers`; policy-aware floors in Evaluate, kickoff, batch and role validation.
- `internal/membership`: `NoReviewerAfterDropout` (shared proof with `remaining = 0`); the gate accepts a proved
  zero-reviewer loss; v3 floor text in blocked notices.
- `internal/driver`: the track non-solo check accepts a proved v3 sole survivor.
- `internal/app`: new ideas freeze v3; review dispatch with zero proved reviewers stops with an organizer-close
  message.

### Skill 2.21.0

SKILL.md v3 pointers; `references/EXHAUSTION_AND_HANDOVER.md` "Optional participants (v3)";
`references/ROSTER_AND_PROTOCOL.md` floor text; regenerated manifests.

## Acceptance criteria

- AC1. A v3 batch that leaves only the usable implementer applies; the same batch under v2 blocks.
- AC2. A v3 batch that leaves only a usable facilitator applies and gives it implementation; with no usable
  holder it blocks.
- AC3. The zero-reviewer proof holds only for v3 with an automatic transition; without one, the gate vetoes.
- AC4. The driver keeps a proved v3 sole survivor alive; an unproved solo member escalates.
- AC5. `go test ./... -count=1 -timeout 45m` passes; the skill `npm test` passes; SKILL.md ≤ 20,000 B; the
  phase-1 meta facilitator packet ≤ 70,000 B.
