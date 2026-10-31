---
idea: meta-protocol-change-participants-optional
status: complete
implementer: claude-1
started: 2026-10-31
completed: 2026-10-31
branch: parley-deck-cli#participants-optional, parley-deck-skill#participants-optional
head-commit: 9db33f2
design-pr: n/a
implementation-pr: https://github.com/feci/parley-deck-cli/pull/88, https://github.com/feci/parley-deck-skill/pull/14
independent-review: none
---

## Summary of work

Implemented FINAL D1–D6 for CLI 1.57.0, skill 2.21.0 and staged core 2.21.0.

## Implementation plan / checklist

- [x] Files or areas to change: `internal/quota` (policy v3, floors), `internal/membership` (zero-reviewer proof,
  gate, notice text), `internal/driver` (track check), `internal/app` (new-idea policy, review dispatch message),
  three COOPERATION.md copies, staged core 2.21.0, skill SKILL.md and two references, version files.
- [x] Checks to run: `go test ./... -count=1 -timeout 45m`, `go vet` on changed packages, `gofmt`, skill `npm test`
  (Node, Python, manifest check), packet guardrail, SKILL.md size.
- [x] Review or risk notes: solo run; no independent review.

## Deviations from FINAL.md

None.

## Notes for reviewers

- The phase-1 meta facilitator packet is 69,856 B (guardrail 70,000 B; before: 69,880 B). The plain facilitator
  packet is 68,424 B. The skill core SKILL.md is 19,983 B (limit 20,000 B).
- Seven §9.0 bullets are reflowed without word changes: Policy, Two-attempt path, Provider exhaustion in v2,
  Legacy quota predicate, Batch and floor, V2 succession, Review gates.
- The staged core `~/.parley/staging/COOPERATION-2.21.0.md` is staged 2.20.0 plus exactly these hunks
  (SHA256 dcc4b14e1334bf2710d21d276c2fcac5a32763b73cc737b27a57d00c16568e5e).

## Fix-up after the internal helper check

The helper found three defects and one text issue. All are fixed with tests:

1. A v3 batch could apply with empty membership when a usable facilitator held no role. `Evaluate` now blocks
   below `MinMembers` (`TestOptionalParticipantFacilitatorWithoutRoleIsNoHolder`).
2. Preflight still hard-stopped below two participants after a valid v3 decision. It now accepts an applied v3
   decision that keeps a holder (`TestOptionalParticipantsPreflightFloor`). A first version also accepted a
   deliberate solo start; `TestRunParticipantsSubsetHardStopsSolo` caught it, and commit 9db33f2 narrows the
   check (FINAL D5).
3. The reviewer-loss proof read the current policy, so a later policy-only revision lost the v3 cause. It now
   uses the deciding policy (`TestOptionalParticipantsProofSurvivesPolicyOnlyRevision`).
4. Owner-facing floor text now names the v3 floor.

## Validation evidence

- `go test ./... -count=1 -timeout 45m` at 9db33f2: all 34 packages ok (log SHA256
  0d2700ecb562321555ce6e26aa8d2be82ec221c21674f4e27e96c5e1c9601829).
- `go vet` on the changed packages and `gofmt` on the changed files: clean.
- Skill `npm test` at a9a0cef: 402 Node tests, the Python tests and the manifest check pass. Skill PR CI passed
  on Python 3.10 and 3.13.

AC1–AC5: met by the tests and measurements above.

## Independent review

`independent review: none`. The owner authorized a solo run because every other agent was out of tokens. An
internal helper (a Claude subagent of claude-1) read the Go diff for defects; it is not a participant and its check
is not an independent review (§1, §15). A later idea may review this change.
