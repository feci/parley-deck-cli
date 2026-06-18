---
idea: meta-protocol-change-fusion-execplans
status: complete
implementer: claude
started: 2026-06-18
completed: 2026-06-19
branch: parley-deck-cli#meta/fusion-execplans
head-commit: 8477ea6 (impl) → fix-up cycle 1 (see below)
design-pr: n/a (single-repo; design + impl in one PR)
implementation-pr: <pending>
---

## Summary of work
Ratified the Fusion + ExecPlans inspiration as additive protocol-text guidance.
Applied Edits 1–7 from `FINAL.md` byte-identically to both protocol copies
(`parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`).
No Go logic changed; the embedded default stays genericized. This `IMPLEMENTATION.md`
intentionally **dogfoods** the new living-section format it introduces.

## Implementation plan / checklist
- [x] Files or areas to change: both COOPERATION.md copies (Phase 3/4/5/6/7 + §13 + §3 layout).
- [x] Checks to run: drift guard, `go test ./...`, `go build ./...`.
- [x] Review or risk notes: pure additive text; risk is wording/lockstep, caught by the drift guard.

## Deviations from FINAL.md
None. All seven edits applied as specified; Edit 4 (idempotence) folded into Edit 1's
`## Idempotence & recovery` head as planned.

## Notes for reviewers
- The change is **additive guidance** + a `conditional rigor` rule; no existing rule
  was removed or weakened (FINAL.md immutability reaffirmed, severities unchanged,
  signoff gate unchanged).
- Please confirm: (a) both copies are byte-identical in shared zones (drift guard is
  the machine check); (b) no invariant from the REJECT list was accidentally imported;
  (c) wording is clear and the conditional-rigor escape hatch ("N/A" for trivial
  ideas) is unambiguous.

## Progress
- [x] (2026-06-18 20:10Z) Research captured → `reference/research.md`.
- [x] (2026-06-18 20:25Z) Round-01 ×4 (claude/codex/agy/hermes) collected.
- [x] (2026-06-18 20:40Z) `consensus.md` drafted.
- [x] (2026-06-18 22:30Z) Phase-3 signoffs: claude ✅, codex ✅, hermes ✅; agy waived.
- [x] (2026-06-19 00:05Z) `FINAL.md` (edit spec) written.
- [x] (2026-06-19 00:30Z) Edits 1–7 applied to BOTH COOPERATION.md copies.
- [x] (2026-06-19 00:35Z) Drift guard + `go test ./...` + `go build ./...` green.
- [ ] Phase 6 review (codex, hermes) → Phase 7 consensus → merge → release.

## Decision Log
- Decision: Keep `FINAL.md` static; put living sections in `IMPLEMENTATION.md` only.
  Rationale: consensus snapshot = audit trail (agy/codex; hermes endorsed the flip).
  Date/Author: 2026-06-18 / claude (from consensus item B).
- Decision: F (confidently-wrong) ships as §13 evidence guidance only; no `parley
  retro` mining code. Rationale: consensus scoped tooling as a deferred follow-up.
  Date/Author: 2026-06-18 / claude.
- Decision: agy live-signoff waived. Rationale: repeated headless hang; operator
  authorized skip; agy round-01 endorsement on record. Date/Author: 2026-06-18 / claude.

## Surprises & Discoveries
- The two protocol copies differed in **only** the header + the two §2 roster tables
  (confirmed by `diff`); every Phase/§13/§3 section was already byte-identical, so the
  same edit strings applied cleanly to both. Evidence: `diff` showed changes only at
  header lines and roster rows; drift guard green after identical edits.

## Validation evidence
- Acceptance #1 (both copies carry Edits 1–7, byte-identical shared zones): drift
  guard `TestEmbeddedDefaultMatchesLiveDeck` → ok.
- Acceptance #2/#3 (`go test ./...`, `go build ./...`): all packages ok; build OK.
- Acceptance #4 (embedded stays genericized): no roster/workspace leaked; drift guard
  enforces the allowlisted zones.

## Outcomes & Retrospective
Shipped all 7 edits as additive, conditional-rigor protocol guidance across both
COOPERATION.md copies; drift guard + `go test ./...` + build green; Phase 6 review
(codex clean, hermes 1 MINOR + 2 NITs) resolved in fix-up cycle 1. Gaps: none
outstanding; the `parley retro` confident-error *mining* tooling stays a deferred
follow-up. Lessons: (1) agy headless `--print` is unreliable for append tasks —
prefer codex/hermes for canonical-artifact writes, waive agy with an operator-recorded
note when it hangs; (2) keeping the two protocol copies' shared zones byte-identical
makes a fix trivially replayable (same edit strings to both); (3) this idea dogfooded
its own new IMPLEMENTATION.md living-section format end to end.

## Fix-up cycle 1
status: complete
completed: 2026-06-19
head-commit: <pending fix-up commit>

### Fixes applied
- [MINOR] Decision Log now carries the deviations cross-reference ("Deviations still
  go under `## Deviations from FINAL.md` above.") — both copies.
- [NIT] General rigor trigger harmonized to include "pipeline" in Phase 4 + both
  Phase 5 spots (matches the ratified preamble's four triggers) — both copies.
- [NIT] Phase 4 paragraph split: `Idempotence & recovery` moved to its own paragraph
  for scannability — both copies.

### Deviations from agreed fixes
None. Drift guard + `go build ./...` re-verified green; diff stayed symmetric (13/13).
