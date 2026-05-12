---
agent: claude
idea: consensus-workflow-cli
review-round: 2
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
reviewed-commit: 1337d6d
---

## Summary

Fix-up cycle 1 applies every item in `review/consensus.md` — the CRITICAL triage bypass, the two `MAJOR` protocol deviations I raised, and the agreed `MINOR` robustness items — with matching tests for the load-bearing paths. Deferred follow-ups are still scoped as deferred. I find no new regressions; one observation about the finalized-status guard going beyond the consensus list (a welcome bonus, not a defect).

## Findings

### Agreed-fix verification

1. **CRITICAL — triage bypass via manual aliases** — Fixed. `validateDocument` at `internal/consensus/consensus.go:390-395` now normalizes through `CanonicalStatus` before switching, and `TestManualStatusAliasesAffectTriage` (`consensus_test.go:90-112`) regresses the `Status: block` case to `TriageBlocked` with zero errors.
2. **Parser tolerance**
   - Trailing parenthetical statuses: `CanonicalStatus` at `consensus.go:290-294` strips the `(...)` segment before matching. Exercised indirectly by the regression test (the canonical `Status: block` flows through).
   - Short `Counter-proposal:` form: `parseDocument` at `consensus.go:357-363` now matches any `Counter-proposal*` prefix and splits on the first `:`. `TestManualStatusAliasesAffectTriage` covers this prefix end-to-end (the manual signoff writes `Counter-proposal: …`).
3. **`reopen` protocol state** — Fixed. `Reopen` at `consensus.go:262-286` resolves the latest round via `selectRound` (with a `round-01` fallback) and writes that label into `00-prompt.md`; `nextAbortedPath` at `consensus.go:646-655` produces `<round>-consensus-aborted-<NN>.md`. `TestReopenBlockedConsensus` (`consensus_test.go:223-226`) asserts both the new filename and the `round-01` prompt status.
4. **`### Non-goals` scaffold + headings test** — Fixed. `finalTemplate` at `consensus.go:568` emits the section; `TestFinalizeCreatesFinalAndUpdatesStatus` (`consensus_test.go:137-144`) asserts every required heading.
5. **Positive `reserved + open-items` finalize test** — Added: `TestReservedFinalizeSucceedsWithOpenItems` (`consensus_test.go:165-189`) writes an open-items entry, signs off with `reserve`, and asserts `Finalize` succeeds and updates the prompt status to `final`.
6. **Review consensus frontmatter (numeric cycle + `reviewed-commit`)** — Fixed. `draftTemplate`/review branch at `consensus.go:502-527` emits numeric `cycle:` derived from the round label and a `reviewed-commit:` field; the CLI exposes `--reviewed-commit` (`internal/app/app.go:258,270`). `TestReviewDraftUsesReviewPath` (`consensus_test.go:236-250`) asserts both fields.
7. **`consensus=error` surfaced in workspace listing** — Fixed. `consensusTriageLabel` at `internal/app/app.go:651-664` now distinguishes `os.ErrNotExist` (empty label) from other read errors and from malformed summaries (`len(summary.Errors) > 0`), returning `consensus=error` for both error paths.
8. **Numeric round sort** — Fixed. `selectRound` at `consensus.go:444-465` parses the numeric suffix and sorts by integer. `TestDraftSelectsLatestRoundNumerically` (`consensus_test.go:254-271`) exercises the `round-2` vs `round-10` ordering.

### [NIT] `Finalize` finalized-status guard added beyond the consensus list

`Finalize` at `consensus.go:201-203` now rejects re-running when `idea.Status == "final"`. The consensus explicitly deferred this as a "remaining non-blocking follow-up" in `IMPLEMENTATION.md`. The new guard addresses the open question I raised in round-01 and is strictly an improvement, but it is an undocumented deviation from the agreed scope. Either drop the `## Deviations from FINAL.md` follow-up line listing this as still-deferred, or note in `IMPLEMENTATION.md` that the guard was added during fix-up. Non-blocking.

### [NIT] `consensus=error` path is not directly tested

`consensusTriageLabel` distinguishes three branches (not-exist → empty; read error → `consensus=error`; malformed summary → `consensus=error`). The behavior is correct on inspection, but no test exercises the malformed-summary branch through `runStatus`. A future regression could re-silence the label without breaking any test. Non-blocking.

## Tests reviewed

- `internal/consensus/consensus_test.go`:
  - existing: `TestDraftAndAppendSignoffTriage`, `TestMalformedSignoffs`, `TestFinalizeCreatesFinalAndUpdatesStatus` (extended), `TestReservedFinalizeRequiresOpenItems`, `TestReopenBlockedConsensus` (extended), `TestReviewDraftUsesReviewPath` (extended);
  - new: `TestManualStatusAliasesAffectTriage`, `TestReservedFinalizeSucceedsWithOpenItems`, `TestDraftSelectsLatestRoundNumerically`.
- `internal/app/app_test.go`: unchanged for this cycle; existing `TestConsensusCLIWorkflowAndIdeaStatus` continues to cover the CLI surface.

Coverage I would still like to see in a follow-up (non-blocking):

- a direct test for the workspace listing rendering `consensus=error` when `consensus.md` is malformed or unreadable;
- a dedicated parser-level test for a hand-written signoff using the short `Counter-proposal:` form against a non-`block` status, to lock down the prefix-relaxation behavior independently from the triage regression.

## Verdict

APPROVE
