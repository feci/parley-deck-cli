---
agent: hermes-1
idea: protocol-read-cost-regression
round: 3
date: 2026-08-11
---

verdict: CLEAN

# Review — fix-up cycle 2, round 3

## Summary

The response to round 2 is a constant: `compactionEnabled` is `const false`, so the first
executable line of `frontierContext` returns `full()` before anything else can run. I tried to
find any input, file, env var, or config that causes compaction or changes what an agent
receives compared to before this idea, and I could not. The build, vet, and the full
internal/runner suite are green. My round-2 NITs are addressed or mooted. I have no finding
that blocks release.

## 1. Is the feature genuinely inert as shipped?

YES. `frontierContext` (frontier.go:56-58) checks `if !compactionEnabled { return full() }`
as its first act. `compactionEnabled` is `const false` (frontier.go:80). A Go const is
compile-time; there is no reflection, no `unsafe`, no build-tag conditional that overrides it.
The function returns `full()` before `authoredLedger`, `renderRound`, or `fallbackTo` can
execute on any path.

I traced every caller of `frontierContext`:
- design cross-review: runner.go:938, full = `gatherPriorRounds` — returns every prior round.
- review: phase58.go:494, full = `reviewRoundsOnly` — returns every prior review round.

Both `full` callbacks walk every round and exclude only `_index.md` and `_ledger.md`. Since
no `_ledger.md` exists in any current deck (I searched), the exclusion is a no-op on existing
data and the output is the same bytes the old `gatherPriorRounds` / old `gatherReviewContext`
produced. The review-consensus drafter (runner.go:925) calls `gatherReviewContextFull` directly
and never reaches `frontierContext` at all — unchanged from before this idea.

No env var, flag, or config key reads `compactionEnabled`. No prompt tells an agent to create
`_ledger.md`. No code writes `_ledger.md`. The only way to turn compaction on is to change the
constant in source, which is a review-gated edit. The feature is inert.

## 2. Is anything now dead code that claims to be a guard?

Lines 59-73 of `frontier.go` (`authoredLedger`, `renderRound`, the ledger-formatting block,
`fallbackTo`) are unreachable at runtime today. The tests exercise them by calling
`frontierContext` with `compactionEnabled` effectively false — so even the tests do not
execute the compaction branch; they assert properties of the `full()` return.

This is acceptable and I would not require deletion. The dead code is not a guard that hides
rot — it is the machinery for a feature that is intentionally staged behind a constant. The
comment at frontier.go:76-79 states the enabling condition (a validator covering G3, G5, G6),
and the constant is the gate. Deleting the machinery would mean re-introducing it from scratch
later, which is worse than carrying compiled-but-unreachable code with a clear enablement
contract. Go's compiler and vet are happy with it. The structural tests
(`TestCompactionIsOffEvenWithAnAuthoredLedgerPresent`, `TestAuthoredLedgerPathIsRetainedButDormant`)
explicitly assert that the machinery is retained and dormant, which is the right framing.

One thing that IS worth saying for the record: the `authoredLedger` function's content checks
(BOM stripping, TrimSpace) are dead today and will only become live when the validator exists.
When that happens, those checks must be re-reviewed — they are not a validator, and the
IMPLEMENTATION.md correctly says so. This is not a finding against this cycle.

## 3. Are the two structural source-level tests sound?

YES. They test properties that cannot be observed through output today, and they test the
right properties.

`TestReviewConsensusDispatchUsesTheFullWalker` (frontier_test.go:272) reads `runner.go`,
finds the `case "review-consensus":` block, and asserts that block contains
`gatherReviewContextFull(` and does NOT contain `gatherReviewContext(opts`. This is sound:
the only way to satisfy both is for the dispatch to call the full walker and not the frontier
walker. A pattern that called the frontier walker would fail the second assertion. I verified
the 700-byte window (line 281) is large enough to contain the whole case block (it is ~350
bytes of source). The test would break if someone reformatted the dispatch to call
`gatherReviewContext` — which is exactly what it should catch.

`TestBothRoundWalkersExcludeTheLedgerFile` (frontier_test.go:290) reads `runner.go`,
`frontier.go`, and `phase58.go`, and for every occurrence of `e.Name() == "_index.md"` asserts
that `ledgerFileName` appears within 120 bytes. This is sound for the three current walkers:
each has the ledger check on the same logical line as the `_index.md` check, well within 120
bytes. If someone added a new walker that skips `_index.md` but not the ledger, the test would
catch it. If someone removed the ledger check from an existing walker, the test would catch
it.

These are pattern-matching tests, but they pattern-match the right thing: the source-level
contract that the dispatch and walkers must satisfy. They are the same style as the protocol
drift guard, and they are appropriate when the behavioural output is identical with or without
the guard (which the reversion check proved). A pattern-match test can pass while behaviour is
wrong only if the source matches the pattern but the semantics differ — e.g., the dispatch
calls `gatherReviewContextFull` but that function itself delegates to the frontier. I checked:
`gatherReviewContextFull` (phase58.go:281) is the verbatim pre-idea walker. It does not call
`frontierContext`. So the pattern match and the semantics agree.

## 4. Does the prompt text an agent receives now say something true?

Almost entirely. The cross-review instruction (runner.go:1002) now reads:

  "READ everything below and respond to the other participants by name ... Older rounds appear
  either in full or as a carry-forward ledger; a banner above says which. An objection is live
  until ITS OWN OWNER withdraws it. Open any full artifact on disk if you need it."

With compaction off, every agent receives full history and no banner. The sentence "Older
rounds appear either in full or as a carry-forward ledger; a banner above says which" is
therefore technically true — "either in full" is the case that always holds — but the
"carry-forward ledger" alternative never occurs today, and the "banner above says which"
clause refers to a banner that is never emitted. An agent reading this literally would look
for a banner that does not exist.

This is a NIT, not a blocker. The sentence is not false (it is a disjunction where one
disjunct always holds), and it does not instruct the agent to do anything wrong. It is
forward-looking text for when compaction is enabled, shipped alongside the dormant machinery.
I flagged this wording in round 2 as a MINOR and it persists, but it does not make the prompt
misleading enough to block release of an inert feature. When compaction is enabled, the text
becomes accurate in both disjuncts; until then, it is accurate in one. Acceptable.

The review prompt (phase58.go:236) and the review-consensus prompt (phase58.go:351) are
unchanged from before this idea and say nothing about ledgers or compaction. They are true.

## 5. Anything that would make me refuse release today?

No. The shipped behaviour is byte-identical to before this idea for every deck that lacks
`_ledger.md` (which is all of them). The constant makes compaction impossible without a source
change. The tests compile and pass. The structural tests guard the two protections that are
unobservable through output. The prompt text is slightly forward-looking but not false. I see
no defect that would cause harm if this shipped.

## Round-2 NITs — status

My three round-2 NITs:

1. `renderRound` and `gatherPriorRounds` did not exclude `_ledger.md` — RESOLVED. Both now
   exclude `ledgerFileName` (frontier.go:134, runner.go:964). The structural test
   `TestBothRoundWalkersExcludeTheLedgerFile` guards all three walkers.

2. `TestPartialOrEmptyLedgerFallsBack` did not test the partial case — NOT ADDRESSED, but
   mooted. The test still has only the "empty ledger" entry. With compaction disabled, the
   partial case cannot be exercised through `frontierContext` (the constant returns `full()`
   before `authoredLedger` runs), so the missing test case cannot distinguish correct from
   incorrect behaviour today. It becomes relevant only when the constant flips, at which point
   the test suite must be re-expanded. Not a finding against this cycle.

3. `TestFrontierRoundTwoIsUnchangedFullHistory` did not check for the fallback banner — NOT
   ADDRESSED, but mooted for the same reason: with compaction off, no banner is ever emitted,
   so the absence of the assertion cannot mask a regression.

## What must be true before the constant flips

The IMPLEMENTATION.md states this (fix-up cycle 2, lines 159-162): a validator covering G3,
G5 and G6 must exist. I concur and add nothing new — my round-2 gates (structured ledger
fields, protocol proposal for the format, no prompt telling agents to create `_ledger.md`
until the format is specified) all still apply and are all still unmet, which is correct
because the feature is off.
