---
agent: kimi-1
idea: protocol-read-cost-regression
round: 2
date: 2026-08-10
---

verdict: NOT CLEAN — no CRITICAL, 1 MAJOR, 2 MINOR, 3 NIT

## Summary

The approach change is the right call and it is honestly executed. The fail-open marker extractor is
gone, not patched; compaction is gated on a participant-authored `_ledger.md` in every round 1..N−2;
absence and emptiness both fall back per round, naming the round; the head is emitted exactly once
with no content sniffing; and the review-consensus drafter is structurally routed to
`gatherReviewContextFull`. I re-ran `go build ./...`, `go vet ./internal/runner/`, the full suite
(exit 0), and every frontier test including the dispatch-level consensus test — all green on my own
run. My two round-1 CRITICALs are resolved: one structurally (the drafter), one by deletion (the
extractor). What remains is one MAJOR about the gate being fail-closed but not *latched*, and the
explicit record of what must be true before compaction may ever fire for real — because today the
optimization is inert: no runtime path, prompt, or template authors `_ledger.md`, so every round gets
full history plus a banner, and the quadratic read cost this idea exists to cut is unchanged.

## Refutation attempts

- Re-ran `go build ./...`, `go vet ./internal/runner/` — green. `go test ./... -count=1` — exit 0.
  `go test ./internal/runner/ -count=1` — ok (8.0s). Ran the seven frontier tests verbosely — all PASS.
- Grepped the whole repo for `_ledger` / `ledgerFileName` / `frontierContext` / `authoredLedger`
  (Go and markdown): the symbol exists only in `frontier.go`, `phase58.go` (skip list), and
  `frontier_test.go`. **No prompt, template, docs page, or runner path ever writes `_ledger.md`** —
  so the compaction branch is unreachable in normal operation. Fail-closed confirmed structurally.
- Grepped `internal/app` for `frontier` / `gatherPriorRounds`: only `driver_consensus.go:50,110`
  building the design drafter prompt the old way. Design consensus drafter untouched — boundary holds.
- Traced the review-consensus dispatch: `RunReviewConsensus` (`phase58.go:328-337`) sets
  `Phase: "review-consensus"` → `runAgent` → `buildPromptForRound` case `"review-consensus"`
  (`runner.go:919-929`) → `gatherReviewContextFull(opts.Idea.Path, roundNumber(opts)+1)`. No frontier
  on that path. Verified through the real dispatch by executing
  `TestReviewConsensusDrafterGetsFullHistoryThroughDispatch` — PASS (see Q4).
- Tried to construct compaction-without-ledger inputs (Q3): none exists in code. The closest edges:
  a stray hand-authored `_ledger.md` (Finding 1), an existing ledger dropped by the review *fallback*
  walker (Finding 2), a BOM-only ledger passing the emptiness check (NIT 3).
- Revert-redness analysis per test, without touching source (hard rule): traced each test's
  assertions against the plausible reversion of its guard (Q5). No additional inert tests found.
- Confirmed the stale round-1 instruction text survives at `runner.go:1002` (NIT 1) and that
  `TestPartialOrEmptyLedgerFallsBack` tests only the empty case despite its name (NIT 2).

## Q1 — Are my round-1 findings resolved, or only made unreachable?

| Round-1 finding | Status now |
| --- | --- |
| CRITICAL — review-consensus drafter compacted | **Resolved structurally**, not by unreachability. The `"review-consensus"` case calls `gatherReviewContextFull` unconditionally; even with authored ledgers present in every review round the drafter gets full history (the new test's fixture deliberately plants ledgers to prove exactly this). This fix is live today, independent of the feature being off. |
| CRITICAL — marker extraction fail-open | **Resolved by deletion; the feature it served is unreachable.** No extractor exists, so the extraction-gap class cannot recur in code. Compaction itself can only fire when every round 1..N−2 carries a non-empty `_ledger.md`, and nothing creates those files. Acceptable as "unreachable because off" — with the binding preconditions listed under Q6 before it ships for real. |
| MAJOR — FINAL/IMPLEMENTATION doubled in review rounds 1–2 | **Resolved for real, on the live path.** `gatherReviewContext` emits the head once and the round≤2/full-fallback walkers (`reviewRoundsOnly`) carry no head. This path executes today; not gated on ledgers. |
| MAJOR — head-drop via banner sniffing | **Resolved for real.** No `strings.Contains` decision remains anywhere in the assembly; `TestQuotingTheBannerCannotStripTheHead` guards reintroduction. |
| MAJOR — G3 triggers mostly absent | **Partially resolved.** Missing (per round, named in the reason) and empty now fire. Invalid/garbage, ambiguous, challenged, and verdict-conflict triggers do not exist. Unreachable today; binding precondition for enablement. |
| MAJOR — owner attribution by filename | **Resolved by deletion** (no extractor, nothing to misattribute). The authored-ledger equivalent — who may write items about whom in a single shared per-round file — is an enablement precondition, not a live defect. |
| MAJOR — G6 unmet | Still unimplemented. Deferral now correct — see Q6. |
| MINOR — NIT/RESERVED markers missing; MINOR — TrimSpace breaks verbatim; MINOR — vacuous G4 test | **Moot / resolved.** The renderer and marker list are deleted; the G4 test was replaced with a dispatch-level test (the first replacement was itself inert and was caught by the reversion check — I verified the current one exercises the switch; Q5). |
| MINOR — no `gatherReviewContext` coverage | **Resolved.** Three tests now cover head-exactly-once-on-fallback, compaction, and banner-quoting. |
| NIT — instruction claims a ledger even on fallback | **Not fixed** — restated as NIT 1. |
| NIT — G7 harness unsafe | Accepted as reported: fix-up cycle records 3/3 reverts applied and restored. I could not re-run reversions (review rules forbid touching source); I verified statically that each surviving test's assertions depend on its guard. |

## Q2 — The boundary

Holds. `frontierContext` and everything it touches only assemble prompt strings (`runner.go:938`,
`phase58.go:494`). No validator, close condition, driver gate, or artifact-acceptance rule references
`_ledger.md` or the frontier; `validateArtifactForPhase` is unchanged in shape; the design consensus
drafter in `internal/app/driver_consensus.go` has zero frontier references; the review-consensus
drafter is now explicitly full-history. The owner-disposes sentence in the ledger footer is prompt
text addressed to participants, not a rule the code enforces — context, not consensus machinery. The
boundary comment at `frontier.go:11-14` accurately describes what the code does.

## Q3 — Is the fail-closed default genuinely fail-closed?

**On the letter of the design: yes.** I could not construct an input where compaction happens without
a `_ledger.md` file: `authoredLedger` returns non-empty only when *every* round 1..N−2 yields a
readable, non-whitespace file; any failure returns `""` plus a reason; `found == 0` is unreachable
dead defense since `round ≥ 3 ⇒ upTo ≥ 1`. Emptiness uses `TrimSpace`, which covers NBSP and Unicode
spaces. With no ledger anywhere, every design round ≥ 3 and every review round ≥ 3 gets the old full
walk plus a visible banner naming the missing round — byte-for-byte the pre-idea behavior plus an
announcement.

Three residual inputs, in descending order of concern:

1. **A stray or self-authored `_ledger.md` flips compaction on with none of the deferred
   protections.** The gate infers "participant-authored" from file existence alone. A participant
   that writes `_ledger.md` on its own initiative (agents deviate from "create exactly this file";
   a participant could also be *told* via quoted text that ledgers are good practice) enables
   compaction for all later rounds with no G6 joining, no invalid/ambiguous/challenged triggers, and
   no authoring contract — the round-1 fail-open class re-enters through the front door. Compaction
   is at least visible in-prompt via the `CARRY-FORWARD LEDGER` banners, which is the falsifiability
   my signoff required; but the consent damage happens in the same round it first appears. Finding 1.
2. **The review *fallback* drops an existing ledger.** Input: review round 3; `review/round-02/
   _ledger.md` exists (someone authored one); `review/round-01/_ledger.md` does not → fallback fires
   on round 1 → `reviewRoundsOnly` skips `_ledger.md` (`phase58.go:530`) → content that exists on
   disk is omitted while the banner claims "Every prior-round artifact follows in full." The design
   fallback (`gatherPriorRounds`) and the compaction-path renderer (`renderRound`) both *include*
   it — the review fallback is the odd one out, and the unsafe direction. Finding 2.
3. A BOM-only or zero-width-only ledger passes the emptiness check (Go's `unicode.IsSpace` excludes
   U+FEFF). Contrived — requires authoring a degenerate file — NIT 3.

## Q4 — Review-consensus, verified through the real dispatch

Executed, not just read: `go test -run TestReviewConsensusDrafterGetsFullHistoryThroughDispatch -v`
→ PASS. The test calls `buildPromptForRound` with `Phase: "review-consensus"` — the same switch
`runAgent` reaches from `RunReviewConsensus` — with authored ledgers planted in both review rounds,
and asserts both review-round bodies arrive. Revert-redness by trace: reverting the guard (routing
the case back through `gatherReviewContext`) makes `frontierContext` compact — the planted
round-01 ledger satisfies the gate, "REVIEW1 finding body" never reaches the prompt, and the test
goes red. This is exactly the property the inert first replacement lacked: the fixture forces the
reverted guard to *change observable output*. The one path not covered at dispatch level is the
design side (Finding 3).

## Q5 — Are the new tests real?

I traced each against the plausible reversion of its guard (live reverts are forbidden to reviewers):

- `TestNoAuthoredLedgerMeansNothingIsEverCompacted` — revert to compaction-without-gate drops the
  unmarked round-1 line → red on the first assertion; also asserts the banner and the named reason.
  Real.
- `TestAuthoredLedgerCompactsOlderRoundsOnly` — revert to always-full resends "ROUND1 PROSE" → red.
  Also pins the owner-disposes footer. Real.
- `TestPartialOrEmptyLedgerFallsBack` — reverting the emptiness check ships the whitespace ledger and
  drops "ROUND1 PROSE" → red. Real, but see NIT 2: the name promises a *partial* case (some ledgers
  present, one missing) that no test exercises.
- `TestReviewConsensusDrafterGetsFullHistoryThroughDispatch` — Q4. Real, through the real dispatch.
- `TestReviewHeadIsNotDoubledOnFallback` — reverting to a head-carrying fallback walker double-counts
  "THE FINAL BODY" → red. Real.
- `TestQuotingTheBannerCannotStripTheHead` — reintroducing the sniff skips the head (the fallback
  banner itself contains the phrase) → red. Real. Note its fixture attacks via the fallback path, not
  the original compaction-path arrangement; acceptable since the sniffing is deleted outright.
- `TestFrontierRoundTwoIsUnchangedFullHistory` — pins round 2 to full history; no plausible
  single-line revert keeps it green while breaking G1.

No further inert tests found.

## Q6 — Is deferring G6 correct now?

Yes. With compaction unreachable, no verdict conflict can cross a compaction boundary, because the
boundary itself never operates. G6 becomes load-bearing precisely when the first `_ledger.md` is
honored. Deferral is correct **on the record** that this idea must not be booked as delivering its
goal: runtime behavior today is the pre-idea behavior plus a banner; the quadratic read cost is
uncut. What shipped is safe, inert infrastructure — that is a legitimate outcome of "an optimization
that cannot prove it is safe does nothing," but FINAL.md's rank-2 objective is unmet and the idea's
status should say so.

**Binding preconditions before any `_ledger.md` is honored (enablement gate):**

1. Finding 1 fixed — an explicit, deliberate enablement, not file presence alone.
2. A sanctioned authoring path: the round/review prompt instructs each participant to carry their
   *own* live items forward (owner-namespaced), and nothing else writes `_ledger.md`. One shared
   file per round with unattributed items reopens round-1 Finding 6.
3. G6 implemented with its fixture (same claim reworded under a new ID, opposing PRIMARY verdicts
   across the boundary → joined DISPUTED or fallback).
4. G3's remaining triggers: invalid/garbage, ambiguous, challenged.
5. Finding 2 fixed — "full history" must actually be full on the review fallback.

## Findings

### [MAJOR] The authored-ledger gate is fail-closed but not latched: any non-empty `_ledger.md` enables compaction with every content protection deferred

`authoredLedger` (`frontier.go:70-91`) treats file existence + non-emptiness as proof of
participant authorship. Nothing structural prevents an unsanctioned `_ledger.md` — written by a
deviating participant, a confused one, or one nudged by quoted text — from switching compaction on
for every later round, with no G6 joining, no invalid/ambiguous/challenged fallback, no
owner-namespacing, and no validation. That is the round-1 fail-open class reachable without any code
change; the only barrier is behavioral ("create exactly this file"), and this protocol routinely
assumes agents deviate from instructions everywhere else. Visibility of the `CARRY-FORWARD LEDGER`
banner mitigates detection, not prevention. **Fix (either):** (a) gate compaction on an explicit
opt-in — e.g. a deck-config or idea-frontmatter flag read at the call sites, without which
`frontierContext` always selects full history — so honoring ledgers is a deliberate act that can be
coupled to the Q6 preconditions; or (b) at minimum, record the Q6 preconditions in FINAL.md /
IMPLEMENTATION.md as binding on enablement and add a test asserting that a plausible non-empty
`_ledger.md` alone does not compact while the idea is unenabled. (a) is the real fix; (b) is the
paper trail.

### [MINOR] The review fallback silently drops an existing `_ledger.md` while claiming full history

`reviewRoundsOnly` (`phase58.go:530`) skips `ledgerFileName`; the design fallback
(`gatherPriorRounds`) and the compaction renderer (`renderRound`) include it. When the fallback
fires precisely because *some* round lacks a ledger, any ledger that *does* exist — content a
participant authored — is omitted under a banner that reads "Every prior-round artifact follows in
full." Fail-safe in direction today only because ledgers barely exist; exactly wrong once they do.
**Fix:** stop skipping `_ledger.md` in `reviewRoundsOnly`; add a test: ledger present in round-02,
absent in round-01 → fallback fires and the round-02 ledger content is present.

### [MINOR] No dispatch-level test for the design path

The review-consensus phase got a real dispatch test; the design cross-review path
(`buildPromptForRound` default case, `runner.go:937-946`) has none — `frontierContext` is only
tested direct. A future refactor that mis-wires the `dirFor` or `full` closures there compiles and
stays green. **Fix:** mirror the consensus dispatch test with `Phase: ""`, `Round: 3`, no ledgers →
assert the fallback banner and both prior rounds' bodies.

### [NIT] The round ≥ 2 instruction still claims a ledger on the fallback path — i.e., always, today

`runner.go:1002` tells every cross-review participant "older rounds appear as a verbatim
carry-forward ledger" — false whenever the fallback fired (currently 100% of round ≥ 3 prompts,
whose adjacent banner says the opposite) and vacuous for round 2, which has no older rounds.
Restated from round 1. **Fix:** "older rounds appear either in full or as a carry-forward ledger; a
banner says which."

### [NIT] `TestPartialOrEmptyLedgerFallsBack` tests only emptiness

The name promises the partial case; the table has one entry. **Fix:** add the case — round-01 ledger
present, round-02 absent, round 4 → fallback whose reason names round 2
("no authored _ledger.md for round 2").

### [NIT] BOM-only ledger defeats the emptiness check

`strings.TrimSpace` does not strip U+FEFF, so a degenerate file passes as "authored." **Fix:** trim
a leading BOM before the emptiness check (one line), or reject files with no printable rune.

## Open questions

- Enablement mechanics: when the time comes, is `_ledger.md` one shared file per round (who writes
  it, and how are items owner-namespaced?) or one file per participant (`_ledger-<agent>.md`)?
  The current filename implies the former; the owner-disposes rule implies the latter.
- Does this idea close as "shipped inert, enablement gated" with a follow-up idea for the authoring
  contract + G6 + opt-in, or stay open until compaction actually fires once? My recommendation:
  close with the preconditions recorded, since the safety work — the part review can meaningfully
  gate — is done.
