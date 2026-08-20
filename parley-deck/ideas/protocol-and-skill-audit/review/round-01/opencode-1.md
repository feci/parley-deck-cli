---
agent: opencode-1
idea: protocol-and-skill-audit
review-round: 1
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
date: 2026-08-21
---

## Summary

Slice question: did @claude-1 weaken any existing test so a fix would pass?

Reviewed-commit `0bb99031` only edits `IMPLEMENTATION.md`. Scope is existing (not newly added) `*_test.go` files @claude-1 modified in the last two days. Checked exactly the seven named files. **No existing test was weakened.** Every edit either tightens the contract, updates a fixture to a new legitimate contract, or only changes a comment.

## Refutation attempts

Tried to read each diff as a pass-making loosen: exact-match → contains, removed assertions, scaffold fixtures padded so a gate would stop firing, one-step close replaced by a weaker status check. None held.

- `TestFinalizeCreatesFinalAndUpdatesStatus` no longer asserts `status=final` after the first `Finalize`. That is the F5 contract change (scaffold must not close), not a dropped check: the same test then asserts refusal, writes a full FINAL, and still requires `status=final` on the second call.
- Scaffold fixtures in `internal/driver/consensus_test.go` gained N/A required sections. Placeholder / `<todo>` / `<fill in…>` text was kept, and those tests still require escalation/rejection. Padding does not hide the original failure mode.
- `reviewed-commit` added to previously-valid review fixtures is a new required field, plus new cases that fail on missing/empty commit.

## Findings

No CRITICAL / MAJOR / MINOR / NIT on this slice.

Per-file verdicts (existing tests only):

### internal/consensus/consensus_test.go — not weakened

- `TestFinalizeCreatesFinalAndUpdatesStatus` (7b9e6c2, F5): old contract was “one Finalize writes FINAL and sets status final”. New contract is two-step (scaffold stays `consensus`, unwritten scaffold refused, written FINAL then closes). Assertions added, not removed. Section list moved from a short hardcoded set to `protocol.RequiredFinalSections` (stricter).
- `TestReservedFinalizeSucceedsWithOpenItems` (7b9e6c2 + 7f164f9, F5/F9): same two-step wrap; open-item bullet changed from `- Carry this reservation…` to `- codex: carry this reservation…` to match the new identity-prefixed reservation contract. Tighter fixture, same pass/fail polarity.

### internal/driver/consensus_test.go — not weakened

- `validFinal`, `TestConsensusReadyButFinalScaffoldEscalates`, `TestConsensusReadyRevalidatesExistingScaffoldFinal`, `TestFinalScaffoldReason` (31f0083, F22): fixtures updated to the new required-section + matching-slug contract. Rejection cases still reject (short, non-final status, `<...>`). Acceptance case still accepts legitimate `<option>`/`<path>`. Directory for `TestFinalScaffoldReason` now uses slug `demo` so the new slug check does not false-fail. Old contract encoded “one padded heading is enough”; that was the bug F22 fixed.

### internal/runner/phase58_test.go — not weakened

- `TestValidateImplementationAndReviewArtifacts`, `TestFakeReviewHelper` (8bf5501, F18): valid review bodies now include `reviewed-commit`. New cases require missing/empty `reviewed-commit` to fail. Fake reviewer now emits the field the prompt asks for.

### internal/runner/phase58_le_test.go — not weakened

- `TestValidateReviewArtifactRequiresRefutation` (8bf5501, F18): the already-valid “heading with content” fixture gained `reviewed-commit: cafe123` so it still tests refutation, not the new empty-commit failure. Assertion unchanged in polarity.

### internal/track/track_test.go — not weakened

- `TestPolicyForAbsentIsLegacy` (3bcbf7b): comment-only (`F14 deferred`). Assertions unchanged.
- `TestInvalidTrackIsAnErrorNotASilentDefault` (3bcbf7b, F15) and `TestExplicitStandardIsRefusedUnderDeliberationTriggers` (b58e921, F13): new tests. Existing tests not loosened.

### internal/app/app_test.go — not weakened

- `TestConsensusCLIWorkflowAndIdeaStatus` (7b9e6c2, F5): first `consensus finalize` must print `NOT closed` and must not set `status=final`; after a written FINAL, second call still requires `status=final`. Stronger than the old “file exists + status final after one call”.

### internal/protocol/drift_test.go — not weakened

- `TestDefaultCooperationForInit` (b7f25e2, zcode-1/F6): `mustContain` still exact-substring. Expected slug changed from `meta-protocol-change-end-to-end-pipeline` to `2026-06-02T12-07-14-meta-protocol-ch` because the protocol text was corrected (old slug named no idea that existed). Equal strictness, new legitimate string.

## Open questions

None on this slice.
