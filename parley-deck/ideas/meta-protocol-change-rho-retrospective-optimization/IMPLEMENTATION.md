---
idea: meta-protocol-change-rho-retrospective-optimization
status: implemented
implementer: claude
started: 2026-06-16
completed: 2026-06-16
branch: parley-deck-cli#feature/rho-protocol-proposal
head-commit: e1c7582
design-pr: https://github.com/feci/parley-deck-cli/pull/52
implementation-pr: https://github.com/feci/parley-deck-cli/pull/52
---

## Summary of work

Implemented the protocol half of the approved RHO proposal (FINAL.md): added a new
top-level **§13 "Retrospective optimization"** section to BOTH protocol copies,
byte-identically, so the drift guard stays green.

- `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`
  both gain §13 with subsections 13.1 (what a retro pass is), 13.2 (harness
  layers — protocol / runtime-shared "Repository Instruction Files" / local
  "Agent Local Memory" / evidence corpus), 13.3 (acceptance gate = consensus +
  all-participant signoff + human approval + no-regression; self-preference is a
  diagnostic note only), 13.4 (guardrails: audit, adversarial-trajectory hygiene,
  reversibility, multi-agent diagnosis), and the standard "Changing this section
  follows §7 … ratified by idea …" provenance line.
- §13 is normative but does NOT contain tool schemas (those live with the
  separate `parley retro` tooling idea), per consensus D1/D5.

## Implementation plan / checklist

- [x] Author §13 from consensus D1–D7 (advisory-only; never auto-applies).
- [x] Append §13 identically to both COOPERATION.md copies.
- [x] `go test ./internal/protocol/` green (drift guard confirms lockstep).
- [ ] Reviewers' round-01 (codex/agy/hermes), fix-up, complete.

## Deviations from FINAL.md

- None. The tooling half is a separate idea (`rho-retro-tooling`), per the design's
  implementation split.

## Notes for reviewers

- Confirm §13 is byte-identical in both copies (the drift guard test
  `TestEmbeddedDefaultMatchesLiveDeck` enforces it; it is green).
- Confirm §13 is advisory-only and routes all acceptance through the existing
  gate, and that it adds no tool schemas (kept in the tooling idea).
- Confirm placement (new top-level section at end, after §12) and the §7
  provenance line, matching §12's convention.
