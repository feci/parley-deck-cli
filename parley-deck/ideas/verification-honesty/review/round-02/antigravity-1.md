---
agent: antigravity-1
idea: verification-honesty
review-round: 2
date: 2026-06-24
---

## Summary

The implementer has successfully integrated all 8 agreed fixes (`F1` through `F8`) from [consensus.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/verification-honesty/review/consensus.md) into the codebase. Testing verifies that scan robustness, model-diversity events, validation enhancements, strict-gate enforcements, and general test coverage have been correctly applied. However, two new scan-evasion shapes have been identified that still allow real findings to bypass the veto mechanism.

## Refutation attempts

I performed the following verification checks on the updated codebase and diff:
- **Scan Case/Spacing (F2) & Empty Title (F3)**: Validated [scanHasRealFinding](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L338-L366) in [internal/driver/impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go). Verified that lowercase severity tags (e.g. `### [critical]`), extra spacing (e.g. `###   [MAJOR]  `), and empty-title severity headings (e.g. `### [CRITICAL]` followed by finding text on the next line) are successfully recognized.
- **Fail-Closed Read Errors (F1)**: Checked [reviewRoundHasFindings](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L312-L336) in [internal/driver/impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go). Verified that any directory-read error (not matching `fs.ErrNotExist`) and file-read error now cause a fail-closed veto (returns `true` -> escalates) rather than falling through to complete.
- **Model Diversity Events (F4)**: Checked [checkModelDiversity](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go#L127-L157) in [internal/app/driver_impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go). Confirmed that the `agent.model_diversity` event is emitted containing idea, implementer, reviewers, model, and required flag.
- **Validation Headings & Section Checks (F5)**: Checked [ValidateReviewArtifact](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/phase58.go#L412-L441) in [internal/runner/phase58.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/phase58.go). Confirmed that a literal `## Findings` heading line and a non-empty `## Refutation attempts` section are now required (preventing simple prose substring bypasses or blank templates from validating).
- **Strict Gate Field Enforcement (F7)**: Verified that [DraftReviewConsensus](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go#L263-L289) in [internal/app/driver_impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go) correctly raises an error and escalates if `strict_gate_clean` or `closing_review_round` are omitted from the drafted consensus under strict mode.
- **Test Coverage (F8)**: Verified that all newly added tests in [internal/driver/strict_gate_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/strict_gate_test.go), [internal/app/driver_impl_le_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl_le_test.go), and [internal/runner/phase58_le_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/phase58_le_test.go) execute and pass successfully.

## Findings

### [MAJOR] Finding Scan Bypass via Literal Title Placeholder with Subsequent Content

If a reviewer leaves the literal template placeholder `<title>` in a severity heading (e.g. `### [CRITICAL] <title>`) but includes the actual finding description/content on the lines immediately below it, the scanner will bypass it.
- **Why it matters**: [scanHasRealFinding](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L338-L366) evaluates lines independently. It skips the heading line because `rest[closeIdx+1:]` matches `<title>`. It then skips the subsequent description lines because they do not start with `###`. Meanwhile, [ValidateReviewArtifact](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/phase58.go#L412-L441) only checks that `## Findings` exists but does not validate its contents to reject literal placeholders. This allows a real finding to evade the veto during a strict closing round.
- **Suggested Fix**: Update `ValidateReviewArtifact` to reject any review file containing `### [SEV] <title>` where `SEV` is one of `CRITICAL`, `MAJOR`, `MINOR`, or `NIT`.

### [MINOR] Finding Scan Bypass via Non-Standard Heading/Formatting Levels

[scanHasRealFinding](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L338-L366) strictly scans for lines starting with `###` (after trimming).
- **Why it matters**: If a reviewer formats a finding using a different Markdown heading level (e.g. `## [CRITICAL] ...` or `#### [CRITICAL] ...`) or uses bold tags (e.g., `**[CRITICAL]** ...`), the scanner will silently skip it. This allows a finding to bypass the veto while still appearing visually as a finding to a human reader.
- **Suggested Fix**: Update the scanner to tolerate any heading prefix (e.g., `##`, `###`, `####`) when looking for severity brackets.

## Open questions

- **Placeholder validation coverage**: Should we have a more generalized mechanism in `ValidateReviewArtifact` that rejects any remaining template placeholder strings (like `<...>` or `[SEV]`) in both the review round files and the consensus drafts?
