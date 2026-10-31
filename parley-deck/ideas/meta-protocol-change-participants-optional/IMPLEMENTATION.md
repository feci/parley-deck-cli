---
idea: meta-protocol-change-participants-optional
status: implemented
implementer: claude-1
started: 2026-10-31
completed: 2026-10-31
branch: parley-deck-cli#participants-optional, parley-deck-skill#participants-optional
head-commit: pending
design-pr: n/a
implementation-pr: pending
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

## Independent review

`independent review: none`. The owner authorized a solo run because every other agent was out of tokens. An
internal helper (a Claude subagent of claude-1) read the Go diff for defects; it is not a participant and its check
is not an independent review (§1, §15). A later idea may review this change.
