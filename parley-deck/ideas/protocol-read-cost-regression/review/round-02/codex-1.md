verdict: CRITICAL — hard-disable compaction until the signed ledger validator, including G3/G5/G6, exists; CRITICAL — replace the remaining inert tests and complete guard-level G7 reversions

agent: codex-1
idea: protocol-read-cost-regression
review-round: 2
date: 2026-08-10
reviewed-state: d4256a2 plus the current uncommitted implementation

## Summary

BLOCK. The approach change removes the marker extractor, and with no `_ledger.md` (or an empty one)
the code really does fall back to full history. That makes the old marker-specific failures
unreachable in the normal runner today, because no production path authors `_ledger.md`.

It does not make the compaction feature fail-closed. Any non-empty bytes at that pathname are accepted
as a participant-authored ledger without parsing, provenance, expected-participant coverage, ownership,
lifecycle, locator/hash, challenge, or verdict-conflict validation. That file immediately enables
compaction. The safe dormant version must be mechanically off even when such a file exists; the real
optimization may be enabled only after the signed ledger contract and G6 are implemented.

The implementation boundary remains correct: the ledger is used only to assemble prompt context.
No artifact validator or consensus/close predicate reads it. This is still an implementation-scoped
context optimization, not an artifact-validity or consensus rule.

## Round-1 CRITICAL disposition

| Round-1 CRITICAL | Round-2 judgment |
| --- | --- |
| Marker extraction drops dissent / cannot enforce owner disposal | **Marker-specific defect resolved by deletion. Underlying safety condition not resolved.** It is unreachable only while `_ledger.md` is absent; a non-empty unvalidated file re-enables the same orphaned-dissent outcome. |
| G3 uncertainty fallback incomplete | **Not resolved.** Missing and empty files fall back. Invalid, ambiguous, challenged, unresolved-hash/locator, unauthorized transition, incomplete-participant, and unmarked verdict-conflict states are never recognized. |
| G6 claim-ID forking | **Not resolved and not safely deferred in this code.** Deferral is correct only while compaction is mechanically impossible. Here any non-empty file enables it. |
| Review-consensus drafter compacted | **Resolved.** The real `buildPromptForRound` dispatch for `Phase: "review-consensus"` calls `gatherReviewContextFull`; the dispatch-level test is sensitive because ledgers are present and it asserts both old bodies survive. |
| Content sniffing can remove FINAL/IMPLEMENTATION | **Code defect resolved.** The head is emitted once and control flow no longer sniffs participant text. The intended banner-quote regression test is still inert, described below. |
| G7 incomplete | **Not resolved.** Additional guards can be reverted while the claimed tests remain green. |

## Refutation attempts

- **No ledger:** `TestNoAuthoredLedgerMeansNothingIsEverCompacted` passed and retained the unmarked
  round-1 objection with the visible fallback reason. This confirms the default behavior only for
  absence.
- **Arbitrary non-empty file:** `frontier.go:70-90` accepts any non-whitespace file and returns it as
  ledger state. `TestAuthoredLedgerCompactsOlderRoundsOnly` itself writes the unvalidated line
  `- a-1 OPEN: the lock is unchecked` and passes only when the round-1 body is removed. There is no
  production authoring or validation path elsewhere in the repository. This is a concrete input that
  activates compaction without establishing a participant-authored ledger.
- **Orphaned dissent:** put a live objection in round 1, put unrelated non-empty text in
  `round-01/_ledger.md`, and omit the objection from round 2. `authoredLedger` accepts the file,
  `frontierContext` emits the unrelated text plus round 2, and the round-1 objection leaves the
  prompt. `renderRound` also accepts a previous round with only one participant artifact; it has no
  expected-participant set.
- **Real review-consensus dispatch:** executed
  `go test ./internal/runner -run 'Test(ReviewConsensusDrafterGetsFullHistoryThroughDispatch|NoAuthoredLedgerMeansNothingIsEverCompacted|AuthoredLedgerCompactsOlderRoundsOnly|PartialOrEmptyLedgerFallsBack|ReviewHeadIsNotDoubledOnFallback|QuotingTheBannerCannotStripTheHead)$' -count=1 -v`.
  All six tests passed. In particular, the review-consensus test went through `buildPromptForRound`
  and retained both `REVIEW1 finding body` and `REVIEW2 finding body` with ledgers present.
- **Boundary:** repository-wide references to `_ledger.md`/`ledgerFileName` are confined to
  `frontier.go`, review-context rendering, tests, and `IMPLEMENTATION.md`. The phase validators at
  `phase58.go:312-322` do not consume ledger state.

## Findings

### [CRITICAL] A non-empty file, not a validated participant-authored ledger, enables compaction

`authoredLedger` (`frontier.go:70-90`) equates successful `os.ReadFile` plus non-whitespace content
with a safe ledger. It deliberately makes no judgment about content. `frontierContext`
(`frontier.go:43-59`) then drops every older body. There is no schema, author/provenance check,
expected-participant coverage, immutable ID or owner check, lifecycle validation, hash/locator
resolution, challenge expansion, or G6 conflict join. `renderRound` likewise checks only that the
combined previous-round rendering is non-empty.

Therefore the claimed fail-closed state is conditional, not enforced. An invalid or incomplete
`_ledger.md` makes G3, G5, and G6 reachable immediately. In particular, a ledger that omits an
objection, closes it as a non-owner, or forks an opposing PRIMARY verdict is accepted and the source
body disappears. Calling the pathname “authored” does not establish authorship or validity.

**Concrete fix:** for this release, make `frontierContext` use full history unconditionally (or put
the compaction branch behind an internal gate that cannot be enabled by the presence of a round file).
Add a test proving that even a non-empty, plausible-looking `_ledger.md` cannot compact while the
validator is unavailable. Before enabling the optimization for real, implement the participant
authoring/ownership mechanism and the complete signed validator: expected participants, stable IDs,
exact propositions, lifecycle and authorized transitions, tombstones, hashes/locators, challenge
expansion, claim equivalence, and G6 `DISPUTED` handling. Keep all of that strictly in prompt assembly.

### [CRITICAL] G7 still has guards that can be reverted without a red test

Three gaps remain:

1. `TestQuotingTheBannerCannotStripTheHead` creates no ledger. Round 3 therefore takes a genuine
   fallback, so it never constructs the normal **compacted** packet in which a participant's quoted
   banner fooled the old content sniff. Reintroducing that sniff still leaves FINAL in this fixture;
   the test passes for the wrong reason.
2. The ordinary design and review tests call `frontierContext` or `gatherReviewContext` directly.
   Reverting the real design dispatch at `runner.go:938-942`, or the review frontier wiring at
   `phase58.go:494-498`, to the old full walkers leaves the frontier tests green. The review-consensus
   guard is now tested through dispatch; the two optimization dispatches are not.
3. `TestPartialOrEmptyLedgerFallsBack` has only one table row, `empty ledger`. It contains no partial
   case such as a round-1 ledger present and the round-2 ledger missing. A regression that accepts the
   prefix and compacts at round 4 can pass this test.

These are not optional coverage improvements: G7 requires each fix to be reverted, compile, apply,
and make the relevant test red. The current suite does not establish that gate.

**Concrete fix:** add a compacting banner-quote fixture with all required ledger files present; add
real `buildPromptForRound` dispatch tests for design round 3 and review round 3; and add a multi-round
partial-ledger case. For each one, revert the exact production guard/wiring (not a helper), show the
targeted test fail for the intended assertion, restore it, and rerun the green suite.

## G6 deferral

Deferring G6 is correct only for a mechanically disabled optimization. It would be correct after the
current release is changed to ignore `_ledger.md` and always use full history. It is not correct while
arbitrary non-empty ledger content activates compaction. G6 must be implemented and mutation-tested
before that activation path ships for real.

## Validation evidence

- Targeted frontier and dispatch tests: PASS.
- `go build ./...`: PASS.
- `go test ./internal/runner -count=1`: FAIL only at the pre-existing
  `TestDurableKillEndToEndRealProcess` environment case (`no recorded boot id`); the same unrelated
  failure was recorded in round 1.

## Open questions

None.
