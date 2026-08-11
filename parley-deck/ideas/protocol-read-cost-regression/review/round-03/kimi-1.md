---
agent: kimi-1
idea: protocol-read-cost-regression
round: 3
date: 2026-08-11
---

verdict: CLEAN — no CRITICAL, no MAJOR, 3 NIT. My round-2 MAJOR is resolved by a stronger latch
than the one I demanded, and nothing found in refutation would make me refuse a release today.

## Summary

The cycle-2 response does what it claims. `compactionEnabled` is a `const false`
(`frontier.go:80`), referenced nowhere else in the tree; no flag, env var, config key, or file
reaches it. I verified the shipped prompt paths against git HEAD byte-for-byte by reading the
diff: the review context assembly is identical to the pre-idea walker, the design path is the old
`gatherPriorRounds` call behind a pass-through, and the review-consensus drafter gets the
unabridged walk. Two literal exceptions to "byte-identical" exist and both check out as safe (Q1).
The dormant machinery is tripwired so that flipping the constant turns two tests red — enablement
cannot happen silently. Build, vet, the full suite, and every frontier test are green on my own
run. The three NITs: a vacuous banner clause in the one changed prompt sentence, an overstated
"exercised" claim for the dormant path, and my round-2 design-dispatch MINOR restated at NIT now
that the call site is a constant pass-through.

## Refutation attempts

- Re-ran `go build ./...` — green. `go vet ./internal/runner/` — green. `go test ./... -count=1` —
  exit 0. `go test ./internal/runner/ -count=1` — ok (7.9s). All eleven frontier-related tests
  verbosely — all PASS, including both structural guards.
- Grepped the whole repo for `compactionEnabled` / `ledgerFileName` / `_ledger.md` in Go: the
  constant exists only in `frontier.go` plus the tripwire assertion in `frontier_test.go`. No env
  var reads anywhere near it (repo-wide `os.Getenv` audit: SHELL, HOME, GEMINI/GOOGLE keys,
  PARLEY_HOME, VISUAL/EDITOR — all unrelated). No config key. Nothing outside tests writes
  `_ledger.md`; `find` over `parley-deck/` and `tmp-test-plugin/` shows none exists in any deck.
- Diffed `runner.go` and `phase58.go` against git HEAD (the idea's work is uncommitted, so HEAD is
  the pre-idea state) and compared the prompt assembly line by line (Q1).
- Tried to construct a compaction-triggering input (Q1): none exists. Tried to construct an input
  that changes what an agent receives: exactly one found — a hand-planted `_ledger.md`, which is
  now hidden rather than rendered. Safe direction; analysed below.
- Traced revert-redness for every guard statically (reviewers may not touch source): constant flip
  → 2 reds; dispatch revert → both structural assertions red; removing any ledger skip → red.
  Computed the structural test's 700-byte window over `runner.go` — it covers exactly the
  `review-consensus` case (lines 919–931) and does not bleed into the default case.
- Checked the two `_index.md` walkers the structural test does NOT scan (`internal/driver/impl.go:364`,
  `internal/app/pipeline_cmd.go:374`): neither feeds an agent prompt (Q1).
- Confirmed `internal/app/driver_consensus.go` has zero frontier/ledger references — the design
  drafter boundary holds.

## Q1 — Is the feature genuinely inert as shipped?

**Compaction: yes, unreachable by construction.** `frontierContext` returns `full()` on the first
line; `authoredLedger`, `fallbackTo`, `renderRound`, and both banners are runtime-unreachable.
There is no input — file, env var, config, or phase — that reaches them. I tried.

**Byte-identical: two literal exceptions, both verified safe.** Against git HEAD, what an agent
receives differs from pre-idea in exactly two ways:

1. **A file named `_ledger.md` in any round dir is now excluded** from all four prompt walkers
   (`gatherPriorRounds`, `reviewRoundsOnly`, `gatherReviewContextFull`, dormant `renderRound`).
   Pre-idea it would have been rendered as if a participant artifact. This IS a file-triggered
   behavioural delta — the only one I could construct. It does not affect any real deck: nothing
   creates the file, the name was minted by this idea, and none exists on disk. The direction is
   the safe one (machinery is never presented as participant-authored). The claim "byte-identical"
   holds for every input that exists.
2. **The round ≥ 2 design instruction text changed** (`runner.go:1002`) — my own round-2 NIT-1
   fix. The context bytes below the instruction are identical to pre-idea.

The review path is byte-identical apart from those: `reviewHead` + `reviewRoundsOnly` reproduces
the old `gatherReviewContext` exactly (I compared the loops line by line), and the drafter's
`gatherReviewContextFull` is the old walker verbatim plus the ledger skip.

The two non-prompt walkers that don't skip the ledger — the strict-gate findings veto
(`driver/impl.go:352`) and the pipeline round counter (`pipeline_cmd.go:365`) — never reach an
agent prompt; a planted ledger there fails closed (veto) or is contrived (a round dir containing
only a ledger). Not findings.

## Q2 — Dead code that claims to be a guard: acceptable here

The dormant branch is not rot hiding behind a constant; it is a **latched** deferral:

- Flipping the constant turns `TestCompactionIsOffEvenWithAnAuthoredLedgerPresent` red on an
  explicit `if compactionEnabled` assertion AND `TestAuthoredLedgerPathIsRetainedButDormant` red
  behaviourally — enablement forces a deliberate test edit through review. That is stronger than
  the config-flag opt-in my round-2 MAJOR asked for (fix (a)); a file or deck cannot speak to it
  at all.
- The enablement preconditions (validator including G3/G5/G6) are recorded at the constant's
  comment and in IMPLEMENTATION.md.
- Deleting the machinery instead would discard the reviewed fallback semantics and invite a
  fresh, unreviewed re-implementation later. Shipping it compiled, unreachable, and tripwired is
  the right call.

One honesty caveat, NIT 2 below: IMPLEMENTATION.md says the machinery and tests "remain compiled
and exercised" — no test executes `authoredLedger`, `fallbackTo`, or `renderRound`; behavioural
rot inside the dormant path stays green until enablement. Acceptable because enablement itself
requires the validator plus review, but "compiled and tripwired, not exercised" is the true
statement. Consequently `TestPartialOrEmptyLedgerFallsBack` is now inert as a guard of the
emptiness/BOM check it originally tested — harmless, dormant-consistent, but it guards nothing
today.

## Q3 — The two structural tests are sound as drift guards

- `TestReviewConsensusDispatchUsesTheFullWalker`: I computed the window — `src[i:i+700]` covers
  exactly the `review-consensus` case (through `runner.go:931`), contains the real
  `gatherReviewContextFull(` call, and the negative pattern `gatherReviewContext(opts` cannot
  false-positive on `gatherReviewContextFull(opts` (the `(opts` anchor). Reverting the case trips
  BOTH assertions. It is backed by the behavioural dispatch test with planted ledgers, so the
  guard does not rest on source-grep alone. Sound.
- `TestBothRoundWalkersExcludeTheLedgerFile`: all four prompt walkers use the anchored idiom
  within the 120-byte window (verified in the grep output), and the guard already earned its keep
  — it found the `gatherReviewContextFull` leak before shipping. Its limits are inherent to
  source-grep guards: a fixed three-file list, one exact formatting (`e.Name() == "_index.md"`),
  so a differently-shaped or differently-located future walker evades it. That guards drift, not
  malice — acceptable, since the enablement review is the real gate and the two non-prompt walkers
  it doesn't scan hand nothing to an agent.

These tests do not "pass while the behaviour is wrong" in any way I could construct: each has a
concrete, traced reversion that turns it red, and the tripwire pair covers the constant itself.

## Q4 — Does the prompt text an agent receives say something true?

Almost entirely. The only new agent-visible text is the `runner.go:1002` sentence. "READ
everything below" — true; everything is below. "Older rounds appear either in full or as a
carry-forward ledger" — satisfied (always in full). "An objection is live until ITS OWN OWNER
withdraws it. Open any full artifact on disk if you need it" — true and protocol-consistent. The
one false clause: **"a banner above says which" — no banner is ever emitted** while the constant
is off (`fallbackTo` is unreachable), so 100% of shipped round ≥ 2 prompts reference a banner that
does not exist. Harmless direction — no content claim is wrong and nothing is missing — but it is
a false statement in a shipped prompt. NIT 1. The review prompt's context line ("FINAL.md,
IMPLEMENTATION.md, prior review rounds") is accurate.

## Q5 — Would I refuse a release today?

No. Build, vet, full suite green on my own run; the only reachable behavioural delta is
safe-direction and synthetic-only; G6's deferral remains correctly recorded and is genuinely
non-load-bearing while the boundary never operates; the latch on enablement is stronger than what
review demanded. FINAL.md's rank-2 objective (cutting the quadratic read) remains unmet — the idea
should still close as "shipped inert, enablement gated," not as delivered — but that is a
bookkeeping fact, not a release defect.

## Round-2 finding disposition

| Round-2 finding | Status |
| --- | --- |
| MAJOR — gate fail-closed but not latched | **Resolved**, stronger than requested: a constant with a two-test tripwire, not a config flag. |
| MINOR — review fallback drops an existing ledger under a full-history banner | **Resolved.** Exclusion is now uniform across all four walkers, and the contradicting banner is unreachable; the contradiction my finding named cannot occur. |
| MINOR — no dispatch-level test for the design path | Not added. Re-assessed at NIT: with the constant off the call site is a pure pass-through and a mis-wire would break every round ≥ 2 prompt loudly. Restated as NIT 3. |
| NIT — instruction claims a ledger on the fallback path | Fixed in substance; replacement clause references a never-present banner (NIT 1). |
| NIT — `TestPartialOrEmptyLedgerFallsBack` tests only emptiness | Still one case; now folded into the dormancy observation (NIT 2) since the guarded code is unreachable. |
| NIT — BOM-only ledger defeats emptiness check | **Fixed** (`frontier.go:99` trims a leading BOM); correct by inspection; lives in the dormant path. |

## Findings

### [NIT] The round ≥ 2 instruction references a banner that never exists as shipped

`runner.go:1002`: "a banner above says which" is false in 100% of current prompts — with
compaction hard-off, no banner is ever emitted. An agent looking for it finds nothing, but all
content is present, so the blast radius is confusion at worst. **Fix:** while the constant is off,
drop the clause or reword to "Older rounds appear in full below"; restore the banner sentence at
enablement.

### [NIT] "Compiled and exercised" overstates the dormant path

No test executes `authoredLedger`, `fallbackTo`, or `renderRound`; behavioural rot there stays
green until enablement, and `TestPartialOrEmptyLedgerFallsBack` now guards nothing. **Fix:** one
word in IMPLEMENTATION.md ("compiled and tripwired, not exercised"); optionally a direct unit test
of `authoredLedger` so the BOM/emptiness fixes are pinned before they become load-bearing.

### [NIT] Design-path dispatch test still absent (restated from round 2, severity lowered)

`buildPromptForRound`'s default case (`runner.go:938`) has no dispatch-level test. With the
constant off this guards a pass-through, so the gap is small. **Fix:** when enablement lands,
mirror the consensus dispatch test with `Phase: ""`, `Round: 3`, ledgers present/absent.

## Open questions

- None new. The enablement mechanics question from round 2 (one shared `_ledger.md` vs per-agent)
  remains the right first question for the follow-up idea that builds the validator.
