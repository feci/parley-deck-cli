---
agent: codex-1
idea: protocol-read-cost-regression
review-round: 3
date: 2026-08-11
reviewed-state: d4256a2 plus the current uncommitted implementation
responding-to: [codex-1/review/round-02, claude-1/IMPLEMENTATION-fix-up-cycle-2]
---
verdict: NOT CLEAN

## Summary

The round-2 activation finding is resolved: `compactionEnabled` is a package constant set to
`false`, and repository-wide search found no file, environment, configuration, linker, or alternate
call path that can make the compaction branch execute. A non-empty `_ledger.md` no longer enables
compaction.

The stronger claim that the shipped behavior is byte-identical to before this idea is wrong. The
cross-review instruction sent to every design-round agent changed unconditionally, and it now makes
a false banner claim. In addition, a file named `_ledger.md` is filtered from three active full-history
walkers even while compaction is disabled; before this idea it was emitted like any other Markdown
artifact. Those are real inputs that change what an agent receives.

The dormant ledger implementation is not exercised. The constant return makes the rest of
`frontierContext` unreachable, and the only calls to `authoredLedger`, `fallbackTo`, and `renderRound`
are behind that return. The two new source-level tests pattern-match text and can pass while the
behavior they claim to guard is wrong. This scaffolding should be removed, not shipped as a future
safety guard.

## Round-2 disposition

| Round-2 finding | Round-3 judgment |
| --- | --- |
| Any non-empty `_ledger.md` enables unvalidated compaction | **Resolved.** The constant-false gate cannot be changed by runtime input. |
| G7 guards can be reverted while tests remain green | **Not resolved soundly.** The replacement tests are textual heuristics, and the behavioral ledger/fallback tests now return before exercising the guarded implementation. |

## Refutation attempts

- **Activation search — PRIMARY.** `rg` found the only definition at
  `internal/runner/frontier.go:80`: `const compactionEnabled = false`. Its only use is the immediate
  return at lines 56–58. No environment or configuration lookup references it. The targeted frontier
  suite passed, including a non-empty ledger fixture that retained the older round body.
- **Actual design prompt trace — PRIMARY.** `buildPromptForRound` calls `frontierContext`; the
  constant gate returns `gatherPriorRounds` unchanged; `BuildRoundPrompt` then adds the new sentence at
  `runner.go:1002`. No fallback or selection banner is prepended on this path.
- **Pre-idea comparison — PRIMARY.** `git diff` against `HEAD` shows that the old instruction was
  `READ every prior-round artifact below...`. The patch replaces it with ledger/fallback language.
  The same diff shows `|| e.Name() == ledgerFileName` newly added to the design, review, and
  review-consensus full walkers.
- **Structural-test challenge — PRIMARY.** `TestReviewConsensusDispatchUsesTheFullWalker` passes if
  `gatherReviewContextFull(` appears in a comment, dead branch, or overwritten assignment within a
  700-byte slice; it does not establish which context reaches the prompt. It also says nothing about
  whether `gatherReviewContextFull` later compacts. `TestBothRoundWalkersExcludeTheLedgerFile` scans
  only exact textual occurrences of `e.Name() == "_index.md"`, requires no expected occurrence count,
  and merely looks for `ledgerFileName` within the next 120 bytes. Removing both exclusions, using a
  helper, or placing an ineffective ledger expression nearby can leave it green while the ledger is
  emitted.

## Findings

### [MAJOR] The hard-disabled feature still changes the bytes and truthfulness of shipped prompts

`runner.go:1002` now tells every design cross-review agent:

> Older rounds appear either in full or as a carry-forward ledger; a banner above says which.

With `compactionEnabled == false`, older rounds always appear in full and `frontierContext` returns
the old walker output directly. It emits no `FULL HISTORY` banner. The statement an agent actually
receives is therefore false, and the prompt is not byte-identical to the pre-idea prompt.

There is a second concrete input difference: `_ledger.md` was previously included by every full
Markdown walker except the `_index.md` special case. It is now silently excluded by
`runner.go:964`, `phase58.go:299`, and `phase58.go:530`. Thus creating that file changes what an agent
receives compared with the pre-idea implementation even though no compaction occurs.

**Concrete fix:** restore the exact pre-idea cross-review instruction and the pre-idea walker
semantics. Prefer the complete fix in the next finding: remove the dormant frontier feature and
restore the direct full-history call paths until a validated design is ready.

### [MAJOR] The unreachable ledger branch and textual tests are not a sound future guard

The statement that the ledger machinery, fallback, and tests remain “compiled and exercised” is
materially misleading. Go type-checks the code, but the constant return prevents every test through
`frontierContext` from executing ledger selection, validation, fallback banners, or compacted
rendering. Tests named `TestPartialOrEmptyLedgerFallsBack` and
`TestReviewHeadIsNotDoubledOnFallback` take the ordinary full-history path, not `fallbackTo`.

The two structural tests do not repair this: they test nearby source spelling, not dispatch or
artifact behavior, and have the false-positive cases demonstrated above. Keeping unreachable safety
code behind a constant invites exactly the rot the tests claim to prevent and gives a later one-line
enablement change unjustified confidence.

**Concrete fix:** delete `frontier.go` and `frontier_test.go`, restore the direct pre-idea calls to
`gatherPriorRounds` and `gatherReviewContext`, and retain the future ledger contract in the idea
artifacts only. Reintroduce implementation code when the validator exists, at which point semantic
dispatch tests and G3/G5/G6 mutation tests can exercise the real enabled behavior.

## Validation evidence

- Targeted frontier, dispatch, and structural tests: PASS.
- `go build ./...`: PASS.
- `go test ./internal/runner -count=1`: the same environment-only
  `TestDurableKillEndToEndRealProcess` failure as round 2 (`no recorded boot id`); no additional
  package failure.

## Release judgment

Refuse release in the current state. The compaction safety defect is fixed, but the stated inertness
contract is not met and the dormant implementation is guarded only by unsound textual tests.

## Open questions

None.
