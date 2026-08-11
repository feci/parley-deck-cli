---
agent: hermes-1
idea: protocol-read-cost-regression
round: 2
date: 2026-08-10
---

verdict: CLEAN

# Review — fix-up cycle 1, round 2

## Summary

The implementer's response to six CRITICALs was an approach change, not a patch: the derived marker extractor is deleted, and compaction now requires a participant-authored `_ledger.md` per round. With no such file, nothing is compacted and every round gets full history. I verified this end-to-end through the real dispatch, through reversion testing, and by adversarial input analysis.

go build ./..., go vet ./..., and the full internal/runner suite are green. I reverted the review-consensus guard to gatherReviewContext (the frontier path), confirmed the dispatch-level test goes red, and restored it — the tree is clean.

I have no CRITICAL or MAJOR findings. The six round-1 CRITICALs are resolved. The fail-closed default is genuine. The boundary holds. G6 deferral is correct. Three NIT-level observations follow.

## 1. Are the round-1 CRITICALs actually resolved?

Each of my round-1 CRITICALs is resolved, not merely made unreachable. I'll address each:

### CRITICAL: marker extraction fail-open — RESOLVED

The extractor is gone (frontier.go has no marker list, no ledgerItem struct, no extraction logic). `authoredLedger` (frontier.go:70-91) reads a file named `_ledger.md` verbatim. It makes no judgement about content (line 68-69). An objection worded without a marker cannot vanish because there are no markers to miss — the ledger is whatever the participant wrote, byte for byte.

This is a true resolution, not an unreachability fix. The fail-open class cannot recur because the mechanism that was fail-open (substring matching against unbounded natural language) no longer exists.

### CRITICAL: review rounds 1-2 double FINAL.md + IMPLEMENTATION.md — RESOLVED

`gatherReviewContext` (phase58.go:492-503) now emits `head = reviewHead(ideaPath)` exactly once, then calls `frontierContext` with `reviewRoundsOnly` as the fallback. `reviewRoundsOnly` (phase58.go:518-541) walks review rounds only — it does NOT include FINAL.md or IMPLEMENTATION.md. So the head appears exactly once whether or not the fallback fires.

I verified this with `TestReviewHeadIsNotDoubledOnFallback` (frontier_test.go:184): it calls `gatherReviewContext(root, 3)` with no ledger (fallback fires) and asserts FINAL.md appears exactly once. I also confirmed the test catches the old doubling: if the fallback were reverted to `gatherReviewContextFull` (which includes FINAL.md), the count would be 2 and the test would fail.

### CRITICAL: fallback did not cover mandated states — RESOLVED (made unreachable, correctly)

My round-1 MAJOR (which I should have called CRITICAL) listed six uncertainty states my signoff required to trigger fallback: missing, invalid, ambiguous, challenged, unresolved hash/locator, verdict conflict. The new code checks only two: missing (`os.ReadFile` error) and empty (`strings.TrimSpace`). The other four states (invalid, ambiguous, challenged, verdict conflict) have no code path.

This is acceptable because the feature is off. With the authored-ledger approach, the code makes NO judgement about ledger content — it only checks existence and non-emptiness. The states "invalid", "ambiguous", "challenged", and "verdict conflict" are content-level judgements that would require the code to parse and evaluate ledger content, which is exactly what the implementer correctly refuses to do (frontier.go:68-69: "It deliberately makes NO judgement about content").

What must be true before this ships for real: when authored ledgers are introduced, the four unhandled states need a resolution path. Either (a) the participant-authored ledger format includes enough structure (IDs, lifecycle states, DISPUTED markers) that the code can validate without interpreting content, or (b) the protocol's §15.6 duty on the consensus drafter is documented as the sole guard against unstructured ledger content. Option (a) is what codex-1's signoff specified; option (b) is a weaker fallback. The implementer should state which path they intend before authored ledgers ship.

### CRITICAL: content sniffing could strip FINAL.md — RESOLVED

The sniffing guard is gone. `gatherReviewContext` always prepends `head`; there is no `strings.Contains` check on the rendered output. `TestQuotingTheBannerCannotStripTheHead` (frontier_test.go:207) verifies this: a reviewer artifact that quotes the fallback banner cannot remove FINAL.md from a later prompt. The test is real — re-introducing the sniffing guard would make it fail.

### CRITICAL: FINAL/IMPLEMENTATION doubled in review rounds 1-2 — RESOLVED

Same fix as the doubling CRITICAL above. For round <= 2, `frontierContext` returns `full()` = `reviewRoundsOnly`, which excludes FINAL.md/IMPLEMENTATION.md. The head is prepended once by `gatherReviewContext`. No doubling.

### MAJOR (round 1): inert G4 test — RESOLVED

`TestReviewConsensusDrafterGetsFullHistoryThroughDispatch` (frontier_test.go:152) goes through `buildPromptForRound` with phase `review-consensus` — the real dispatch. I reverted the guard from `gatherReviewContextFull` to `gatherReviewContext` (the frontier path) and the test went red: "review-consensus drafter lost REVIEW1 finding body — §15.6 binds here". The test is NOT inert.

## 2. Does the code honour the context-optimization / consensus-rule boundary?

YES. The boundary from my signoff — "implementation-scoped context optimization, not an artifact-validity or consensus rule" — is honoured at the code level:

- No code enforces the owner-disposes rule as a close condition. The consensus drafter (both design and review) gets full history via `gatherPriorRounds` / `gatherReviewContextFull`. §15.6 governs close.
- No artifact validity depends on the ledger. `validateArtifactForPhase` (phase58.go:313) validates frontmatter and headings, never ledger content.
- No consensus-close logic consults the ledger. `ValidateReviewConsensusArtifact` (phase58.go:382) checks `outstanding_agreed_fixes`, not ledger state.
- `authoredLedger` makes no content judgement (frontier.go:68-69).

The instruction at runner.go:1002 still says "an objection in it is live until ITS OWN OWNER withdraws it" — the same text from round 1 that I flagged as MINOR for blurring the boundary. It is a reading instruction, not a validation rule, and no code enforces it as a close condition. With the authored-ledger approach, the owner-disposes rule is more appropriate than it was with the derived ledger (the participant IS the author). But the instruction would be clearer if it said "The ledger is a reading aid, not a consensus rule; the consensus drafter's full-history audit governs close." This remains a MINOR wording issue, not a code defect.

## 3. Is the fail-closed default genuinely fail-closed?

YES. I tried to find an input where compaction happens without an authored ledger. I could not.

`frontierContext` (frontier.go:43-60) calls `authoredLedger(dirFor, round-2)` only when round > 2. `authoredLedger` (frontier.go:70-91) returns non-empty ONLY when every round 1..upTo has a `_ledger.md` that exists (`os.ReadFile` succeeds) and is non-empty (`strings.TrimSpace` != ""). Any missing file returns "" with a reason. Any empty file returns "" with a reason. The `if led == ""` guard at line 49 sends every such case to `fallbackTo`.

No code creates `_ledger.md`. No prompt tells agents to create it. The protocol says "create exactly this file and no other protocol artifact." In normal operation, `_ledger.md` does not exist, `authoredLedger` returns "", and every round gets full history.

Edge cases I checked:
- `_ledger.md` exists for round 1 but not round 2, at round 4: `authoredLedger` loops r=1 (succeeds), r=2 (ReadFile fails) -> returns "" -> fallback. Correct, but not tested (see NIT 1).
- `_ledger.md` is a directory: `os.ReadFile` on a directory returns an error -> returns "" -> fallback. Correct.
- `_ledger.md` contains only whitespace: `strings.TrimSpace` check catches it -> returns "" -> fallback. Tested.
- `_ledger.md` contains the fallback banner string: no sniffing exists, so it cannot strip anything. Tested via `TestQuotingTheBannerCannotStripTheHead`.
- Empty head (FINAL.md/IMPLEMENTATION.md absent): `reviewHead` returns "", head + rounds = rounds. No issue.

An objection can leave the context only during compaction, and compaction requires authored ledgers. If a participant authors a ledger that omits an objection, that objection is not in the compacted context. But this is the participant's responsibility — the code requires an authored ledger and makes no content judgement. This is exactly what the signoff specified.

## 4. The review-consensus path — verified through the real dispatch

I verified by reversion, not by reading.

`buildPromptForRound` (runner.go:919-929) with phase `review-consensus` calls `gatherReviewContextFull(opts.Idea.Path, roundNumber(opts)+1)`. This bypasses `frontierContext` entirely. The test `TestReviewConsensusDrafterGetsFullHistoryThroughDispatch` (frontier_test.go:152) creates authored ledgers for rounds 1 and 2, sets opts.Round=2 (so roundNumber+1=3), and asserts both REVIEW1 and REVIEW2 bodies are present.

I reverted the guard to `gatherReviewContext` (the frontier path) and ran the test. It failed: "review-consensus drafter lost REVIEW1 finding body — §15.6 binds here". The reason: with authored ledgers present, `gatherReviewContext` -> `frontierContext(3)` -> `authoredLedger(dirFor, 1)` succeeds -> compacts round 1 to ledger -> "REVIEW1 finding body" (in round-01/a-1.md) is not in the output. The test catches the reversion.

I restored the guard, confirmed the test passes, and verified the tree is clean.

## 5. Are the new tests real?

Every test is real. I checked each for the inertness class that codex-1 found in round 1 (a test that passes with its guard reverted).

| Test | Guard it protects | Reversion scenario | Catches reversion? |
| --- | --- | --- | --- |
| TestNoAuthoredLedgerMeansNothingIsEverCompacted | `if led == "" { fallback }` | Skip fallback when led="" | YES — unmarked round-1 text drops |
| TestAuthoredLedgerCompactsOlderRoundsOnly | compaction path with ledger | Always-fallback or always-compact | YES — assertion 2 or 3 fails |
| TestPartialOrEmptyLedgerFallsBack | empty-ledger check | Remove TrimSpace check | YES — ROUND1 PROSE absent |
| TestReviewConsensusDrafterGetsFullHistoryThroughDispatch | review-consensus uses gatherReviewContextFull | Revert to gatherReviewContext | YES — empirically verified |
| TestReviewHeadIsNotDoubledOnFallback | head emitted once, fallback excludes head | Revert fallback to gatherReviewContextFull | YES — FINAL.md count = 2 |
| TestQuotingTheBannerCannotStripTheHead | no sniffing guard | Re-introduce sniffing | YES — FINAL.md stripped |
| TestFrontierRoundTwoIsUnchangedFullHistory | round <= 2 returns full() | Change round<=2 to round<=1 | WEAK — see NIT 2 |

The first replacement test that was itself inert (called `gatherReviewContextFull` directly) was caught by the reversion check and rewritten to go through `buildPromptForRound`. I verified the rewritten version catches the reversion. No other inert test found.

## 6. Is deferring G6 correct now that nothing is compacted?

YES. G6 requires opposing PRIMARY verdicts to join as DISPUTED or trigger fallback. With no authored ledgers, nothing is compacted, so no verdict conflict can cross a compaction boundary. G6 is not load-bearing today.

The implementer correctly states: "It becomes load-bearing the moment authored ledgers ship, and must be built before then." This is the right framing. G6 deferral is correct for this cycle.

One observation for the record: when authored ledgers do ship, G6 may require the code to make content-level judgements (detecting opposing verdicts), which tensions against the "no content judgement" principle in `authoredLedger`. The resolution will likely be that the ledger format must include structured DISPUTED markers that the code can check syntactically without interpreting content. This is a design problem for the ledger-format proposal, not for this review.

## NITs

### [NIT] renderRound and gatherPriorRounds do not exclude _ledger.md

`renderRound` (frontier.go:113) skips `_index.md` but not `_ledger.md`. `gatherPriorRounds` (runner.go:964) likewise. `reviewRoundsOnly` (phase58.go:530) correctly excludes both.

During compaction, `renderRound` renders the previous round (round N-1) in full. If round N-1 has a `_ledger.md`, it appears in the context as `===== round-NN/_ledger.md =====` alongside the actual artifacts. This is minor redundancy, not a correctness issue — the agent sees the ledger and the full artifacts for the same round. But it is an asymmetry: the review fallback path excludes `_ledger.md`, the compaction and design-fallback paths do not.

Fix: add `e.Name() == ledgerFileName` to the skip conditions in `renderRound` (frontier.go:113) and `gatherPriorRounds` (runner.go:964), matching `reviewRoundsOnly`.

### [NIT] TestPartialOrEmptyLedgerFallsBack does not test the partial case

The test name says "Partial or Empty" but the table has one entry: `{"empty ledger", "   \n"}`. The partial case — ledger exists for round 1 but not round 2, at round 4 (upTo=2) — is not tested. The code handles it correctly (`authoredLedger` returns "" on the first missing file at line 76), but the test doesn't verify it.

Fix: add a table entry that writes a ledger for round 1 only, leaves round 2 without one, calls `frontierContext(4)`, and asserts the fallback fires.

### [NIT] TestFrontierRoundTwoIsUnchangedFullHistory does not check for the fallback banner

The test checks that round-2 bodies are present and "CARRY-FORWARD LEDGER" is absent, but does not check that "carry-forward fallback" is absent. If the `round <= 2` guard were changed to `round <= 1`, round 2 would go through the fallback path (full history with a banner) and the test would still pass — both bodies present, no "CARRY-FORWARD LEDGER". The behavioral difference is only the banner, so this is a NIT, not a correctness gap.

Fix: add `if strings.Contains(got, "carry-forward fallback") { t.Error("round 2 must not fall back") }` to the test.

## What must be true before authored ledgers ship

The implementer states this correctly in IMPLEMENTATION.md:139 — G6 must be built, and the authored-ledger format must be specified. I add three more gates from this review:

1. The four unhandled fallback states (invalid, ambiguous, challenged, verdict conflict) need a resolution path — either structured ledger fields the code can validate syntactically, or an explicit protocol rule that §15.6 is the sole guard.
2. The ledger format must be specified in a protocol proposal (§7 idea), not just in code comments.
3. No prompt should tell agents to create `_ledger.md` until the format and G6 are specified — otherwise an agent might create a malformed ledger that triggers compaction without the safety properties the format provides.
