---
agent: claude
idea: consensus-workflow-cli
review-round: 1
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
reviewed-commit: 86f4523
---

## Summary

The slice delivers the deterministic consensus primitives described in `FINAL.md`: the `internal/consensus` package, the five `parley consensus` subcommands, and the `parley status --idea` integration. The code is readable, well-bounded, and the test set covers the main happy and error paths. Two real protocol deviations from `FINAL.md`/`COOPERATION.md` need a fix, plus a small set of robustness/coverage gaps worth tightening before merge.

## Findings

### [MAJOR] `reopen` sets idea status to `discussion`, which is not a documented status value

`internal/consensus/consensus.go:271-275` calls `updateIdeaStatus(idea.Path, "discussion")` after a successful reopen. `COOPERATION.md` §4 Phase 0 enumerates the only valid `00-prompt.md` `status:` values as `round-N | consensus | final | abandoned`. `discussion` is not in that set. `FINAL.md` uses the word "discussion" colloquially ("returns the idea to discussion"), but `COOPERATION.md` is the canonical protocol and downstream tooling (`readIdeas`, `consensusTriageLabel`, future status-driven logic) keys off the documented enum.

Suggested fix: revert the prompt status to the most recent `round-NN/` directory the idea has (or fall back to `round-01` if none). At minimum, pick one of the documented values — `round-NN` is the closest fit for "back to deliberation"; `abandoned` would be wrong because reopen explicitly continues the idea.

### [MAJOR] `FINAL.md` scaffold is missing the `Non-goals` section required by `FINAL.md`

`FINAL.md` line 110 states the scaffold must contain "explicit sections for goal, scope, implementation details, tests, **non-goals**, and verification." `finalTemplate` in `internal/consensus/consensus.go:513-539` only emits `### Goal`, `### Scope`, `### Implementation details`, `### Tests`, `### Verification` — no `### Non-goals`. `TestFinalizeCreatesFinalAndUpdatesStatus` only checks that the file exists, so the gap is not caught.

Suggested fix: add `\n### Non-goals\n` to the template between `### Tests` and `### Verification`, and extend the finalize test to assert the section headers it must contain.

### [MINOR] Counter-proposal parser only accepts the exact prefix `Counter-proposal (required if ❌):`

`signoffHeader`'s value loop at `internal/consensus/consensus.go:342-345` matches the field by the literal prefix `"Counter-proposal (required if ❌):"`. The CLI writer in `signoffBlock` (lines 541-552) emits exactly that string, so CLI-written signoffs round-trip. But the parser silently drops any signoff block that uses a shorter human-friendly form like `Counter-proposal:` — the validator then flags the entry with "counter-proposal is required for ❌ BLOCK" even though a counter-proposal is present. Given the protocol expects humans/agents to also append signoffs by hand on transports B/C, accept the shorter prefix too (`strings.HasPrefix(line, "Counter-proposal")` followed by `strings.SplitN(line, ":", 2)`), or document the strict format requirement.

### [MINOR] No test covers the positive `reserved + open-items` finalize path

`TestReservedFinalizeRequiresOpenItems` (`consensus_test.go:114-128`) asserts the failure case. `FINAL.md` line 84 explicitly requires the positive path too ("it may succeed for `reserved` only when reservations are visibly captured in open items"). Add a sibling test that writes a non-empty `## Open items deferred to implementation` section and verifies `Finalize` succeeds. Without it, a regression in `hasSectionContent` (lines 588-606) — e.g. an off-by-one on the heading match — would be invisible.

### [MINOR] `consensus.Status` is recomputed for every idea in the workspace listing

`runStatus` at `internal/app/app.go:447-450` calls `consensusTriageLabel` per idea, and that function calls `consensus.Status` → `findIdea` → `protocol.ReadWorkspaceStatus`. For N ideas the listing performs O(N²) workspace reads. Functional, but cache the workspace status once and pass it in, or expose a `Status` overload that accepts a preloaded `IdeaStatus`. Same pattern exists in `printConsensusIfPresent`.

### [MINOR] `consensusTriageLabel` silently swallows non-`ErrNotExist` errors

`internal/app/app.go:650-656`: any error from `consensus.Status` (e.g. a read error on a partially-readable `consensus.md`) returns an empty label, so the listing row looks clean. `printConsensusIfPresent` correctly distinguishes `ErrNotExist` from other errors and prints `Consensus: error: ...` — the listing path should do the same, or at least surface `consensus=error`.

### [MINOR] `selectRound` relies on `sort.Strings` over lexicographic round labels

`internal/consensus/consensus.go:412-440` collects `round-NN` directories and uses `sort.Strings`. This is only correct because the CLI uses zero-padded `%02d` names. A hand-created `round-1` would sort after `round-10` (`round-1` < `round-10` lexically; `round-2` > `round-10`), which is a foot-gun for any participant who creates an unpadded round directory. Parse the numeric suffix and sort numerically.

### [NIT] `signoffBlock` writes notes as a single line; multi-line notes silently truncate on re-parse

`signoffBlock` (`consensus.go:541-552`) emits `Notes: <one-line>\n`; the parser's per-line `CutPrefix` matches only that first line. A user later editing the file to add a wrapped Notes line will have the continuation absorbed into the next field or dropped. Not a regression caused by this slice, but worth a doc comment so future readers know the format constraint.

### [NIT] `AppendSignoff`'s `appendMu` is process-local

`internal/consensus/consensus.go:87,158-159`: two parallel `parley consensus signoff ...` invocations against the same file from different processes can both pass the duplicate check and both append. Likely acceptable for the human-driven CLI surface — note it explicitly or use a short-lived `flock` on the consensus file.

### [NIT] `nextAbortedPath` has no termination guard

`internal/consensus/consensus.go:608-615`: the `for i := 1; ; i++` loop has no upper bound. In practice it always terminates at the first absent path, but a defensive cap (e.g. 999) would prevent a runaway loop if `os.Stat` ever returned an unexpected non-nil non-ErrNotExist error path that was wrongly treated as "exists."

## Open questions

- Is the `Reopen` aborted-file naming `round-NN-consensus-aborted.md` intentional even though it lives at the idea root (not under `round-NN/`)? If the file is meant to live alongside other consensus history, `consensus-aborted-NN.md` would be a clearer name; if it's meant to be tied to a specific round, it should sit inside that round's directory.
- The `consensus.*` event emission listed in `FINAL.md` is deferred in `IMPLEMENTATION.md` because "the current CLI command context has no project-level event store outside individual run directories." That rationale is fine for this slice, but consider tracking a follow-up idea so the deferral does not get lost; future automation around request-signoffs (the next slice) will likely want these events.
- `Finalize` does not reject the case where the idea is already `final` (the `00-prompt.md` status check is implicit, only via the FINAL.md existence check). If somebody manually deletes `FINAL.md` after finalize, `Finalize` could re-run and overwrite the status. Intentional, or worth an additional guard?

## Tests reviewed

- `internal/consensus/consensus_test.go`: `TestDraftAndAppendSignoffTriage`, `TestMalformedSignoffs`, `TestFinalizeCreatesFinalAndUpdatesStatus`, `TestReservedFinalizeRequiresOpenItems`, `TestReopenBlockedConsensus`, `TestReviewDraftUsesReviewPath`.
- `internal/app/app_test.go`: new `TestConsensusCLIWorkflowAndIdeaStatus` plus existing tests for status/resume/idea-without-runs (`TestResumeReportsKnownIdeaWithNoRuns`, the updated branch in `TestStatusAndResumeUseRunState`).

Coverage gaps worth filling before approving fix-up complete:

- positive `reserved + open-items` finalize path (see MINOR above);
- a parser test for a signoff block written with the shorter `Counter-proposal:` prefix (currently silently dropped — captures behavior either way);
- a test asserting `FINAL.md` contains every section required by `FINAL.md` line 110 (would have caught the missing `Non-goals`);
- a test asserting `Reopen`'s post-condition `00-prompt.md` `status:` is a value defined in `COOPERATION.md` (would have caught the `discussion` value).

## Verdict

REQUEST_CHANGES
