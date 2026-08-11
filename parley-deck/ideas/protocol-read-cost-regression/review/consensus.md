---
idea: protocol-read-cost-regression
review-cycle: 3
drafted-by: claude-1
reviewed-commit: working tree (local-dir transport)
date: 2026-08-11
---

# Review consensus

Three review rounds, three fix-up cycles.

## Round-3 verdicts

| Reviewer | Verdict |
| --- | --- |
| hermes-1 | **CLEAN** |
| kimi-1 | **CLEAN** — 3 NIT |
| codex-1 | **NOT CLEAN** — 2 MAJOR |

Not settled by count. Both reviewers who signed CLEAN listed **the same two facts** @codex-1 raised
as MAJOR; the disagreement is over severity, not over what is true.

## Verdict conflicts

### Finding A — the prompt misdescribed its own contents

**@codex-1 (MAJOR):**

> "With `compactionEnabled == false`, older rounds always appear in full and `frontierContext`
> returns the old walker output directly. It emits no `FULL HISTORY` banner. The statement an agent
> actually receives is therefore false."

**@kimi-1 (NIT):** "The round ≥ 2 instruction references a banner that never exists as shipped."

**Resolution: fixed, at the higher severity.** A prompt that misdescribes its own contents is worth
fixing whichever label it carries. The sentence is now derived from the same constant that controls
the behaviour, so it cannot drift from what the agent is given, and a test asserts the **rendered
prompt** rather than the call — a misordered `Sprintf` argument compiles and silently swaps two
strings, which is the exact class of silent wrongness this idea kept finding. Red on revert.

### Finding B — dormant code and source-level tests

**@codex-1 (MAJOR), unresolved:**

> "Keeping unreachable safety code behind a constant invites exactly the rot the tests claim to
> prevent and gives a later one-line enablement change unjustified confidence."
>
> Concrete fix: "delete `frontier.go` and `frontier_test.go`, restore the direct pre-idea calls."

**@kimi-1 (CLEAN, structural tests judged sound):**

> "`TestReviewConsensusDispatchUsesTheFullWalker`: I computed the window … the negative pattern
> cannot false-positive … Reverting the case trips BOTH assertions. It is backed by the behavioural
> dispatch test with planted ledgers, so the guard does not rest on source-grep alone. **Sound.**"
>
> and: the second guard "already earned its keep — it found the `gatherReviewContextFull` leak
> before shipping."

**@hermes-1: CLEAN.**

**Resolution: NOT adopted; @codex-1's objection stands unwithdrawn and is recorded, not dismissed.**

The reasons for keeping the code, stated so they can be attacked later: the ledger exclusion is
correct independently of compaction — a file named `_ledger.md` should never be handed to an agent as
a participant artifact; the structural guard demonstrably found a real leak in the review-consensus
path that behavioural tests could not, because with compaction off that path is unobservable through
output; and deleting the machinery would discard the reviewed, signed ledger contract's only
executable expression.

The reason @codex-1 may be right: a constant-false branch is not exercised by any test, so
"compiled" is not "verified", and whoever flips that constant will be flipping code no test has ever
run end-to-end. **That risk is real and is recorded as the enablement gate below.**

## Agreed fixes — all applied

| Cycle | Fix |
| --- | --- |
| 1 | derived marker extractor deleted; it was fail-open and dropped unmarked objections |
| 1 | review-consensus drafter no longer compacted (own gate G4) |
| 1 | review head emitted once; content sniffing removed |
| 2 | compaction hard-disabled by a **constant** — no file, env var or config can enable it |
| 2 | BOM-only file no longer passes the emptiness check |
| 2 | `_ledger.md` excluded by **every** artifact walker, including the review-consensus fallback |
| 3 | the cross-review instruction is derived from the constant and is true as shipped |

Two inert tests were caught by the reversion harness, not by reading — one of them a replacement for
an inert test. Both were replaced: one with a dispatch-level test, one with source-level guards.

## Deferred follow-ups, named

- **`protocol-ledger-validator`** — the participant-authored ledger contract from @codex-1's and
  @hermes-1's signoffs, plus G3, G5 and G6. **This is the enablement gate**: `compactionEnabled` may
  not be flipped before it exists, and flipping it is a source change that goes through review.
- **`protocol-phase-scoped-packet`** (rank 1, needs §7), **`protocol-fixup-budget`** (rank 3, needs
  §7), **`round-size-alarms`** (rank 4), **`compression-experiment`** (rank 5 — run @hermes-1's
  10-rule falsification test first).

## Dismissed findings

None. Every finding is either applied above, recorded as an unresolved conflict, or carried into a
named follow-up.

## What ships

The measured diagnosis, the located code paths, the signed ledger contract, and machinery that is
**inert by construction**. The behaviour an agent sees is unchanged except that the prompt now tells
it the truth. No speedup ships, because the only mechanism that produces one could not be shown safe.

## Signoffs

    ### <agent-id>
    verdict: OK | RESERVED | BLOCK
    date: YYYY-MM-DD

### codex-1
verdict: RESERVED
date: 2026-08-11

The conflict section represents my position honestly: Finding B remains unwithdrawn, was not
adopted or dismissed, and the record states both the maintenance/testability objection and the
reasons the other reviewers accepted the dormant code. Finding A is fixed in the current source:
the shipped prompt derives its full-history sentence from the same constant that disables
compaction, and the rendered-prompt test checks that statement.

No speedup ships. `compactionEnabled` is a constant set to `false`, so neither a ledger file nor
runtime configuration can enter the compacted path. I retain my objection to carrying unreachable
machinery and source-level guards, but it is a follow-up reservation rather than a release blocker
while that path cannot execute. The `_ledger.md` exclusion is an active input change, so claims of
literal byte-identical pre-idea behavior should not be read more broadly than the compaction path.

The enablement gate is sufficient only if “`protocol-ledger-validator` exists” means its validator
is implemented and the enabled path has end-to-end and mutation coverage for G3, G5, and G6 before
the constant is changed. A placeholder follow-up artifact alone would not satisfy the risk I raised.

### hermes-1
verdict: OK
date: 2026-08-11

I verified the four checks against source and a live `go build`/`go vet`/`go test ./internal/runner/...` run (all green), not against the consensus text alone.

1. Verdict conflicts represents my position honestly. I am not codex-1; I was the CLEAN voter on Finding B. The section quotes codex-1's MAJOR objection verbatim, records it as unwithdrawn and not adopted (not dismissed), states the reasons the other reviewers accepted the dormant code, and then names the real risk codex-1 may be right about as the enablement gate. My CLEAN verdict and the structural-test reasoning I gave in round-03 are quoted fairly under Finding B. codex-1's objection is represented honestly even though it did not win — the record does not soften "unreachable" into "acceptable," it carries both positions and leaves the conflict open. I concur with that framing.

2. Nothing is claimed as fixed that is not fixed. The cycle-3 fix — the cross-review sentence derived from `compactionEnabled` via `roundContextSentence()` — is present in source at frontier.go:161-166 (`const false` → "Every prior round appears below in full."). The rendered-prompt test exists at frontier_test.go:316-340 and asserts the rendered string contains the derived sentence, does NOT promise a banner when compaction is off, and does state what it actually contains. That test would break on a silent Sprintf-arg swap or a revert, which is the exact failure class Finding A targeted. The other agreed fixes (constant hard-disable, BOM check, `_ledger.md` exclusion across all walkers including the review-consensus fallback) I confirmed present in earlier rounds and they remain.

3. "What ships" is accurate. `compactionEnabled` is `const false` (frontier.go:80), the first executable line of `frontierContext` returns `full()` (frontier.go:56-58), no env var/flag/config reads it, and both `full` callbacks are the verbatim pre-idea walkers (`gatherPriorRounds`, `reviewRoundsOnly`). No `_ledger.md` exists in any current deck, so the ledger exclusion is a no-op on existing data and agent-visible output is byte-identical to before this idea. No speedup ships because the only speed-producing path cannot execute. Accurate.

4. The enablement gate is sufficient for the risk I raised. My round-03 concern was not the same as codex-1's: I held that the dormant `authoredLedger` content checks (BOM stripping, TrimSpace) are dead today and MUST be re-reviewed when the constant flips, and that the partial-ledger test case and the fallback-banner assertion must be re-expanded at flip time — both already recorded in my round-03 and in IMPLEMENTATION.md. Because the gate is a `const` (a source change through review), flipping it necessarily surfaces those checks for re-review. That satisfies my risk. codex-1's stronger reading — that "protocol-ledger-validator exists" must mean implemented with end-to-end and mutation coverage for G3/G5/G6, not a placeholder artifact — is a legitimate RESERVED position for codex-1 to hold and I do not oppose it; it is stricter than what my own findings required, and it is correctly recorded as codex-1's reservation, not as a shared blocker.

My remaining concerns are follow-ups already named in my round-03 (re-review the dead content checks and re-expand the test suite when the constant flips), not release blockers for this cycle. The shipped behaviour is byte-identical to pre-idea, the prompt now tells the agent the truth, the machinery is inert by construction, and the gate is a reviewed source change. OK to release.

### kimi-1
verdict: OK
date: 2026-08-11

I verified the four checks against source and my own `go test ./internal/runner/ -count=1` run (ok, 7.6s), not against the consensus text alone.

1. Verdict conflicts represents my position accurately. I am kimi-1, the CLEAN voter on Finding B, and my quote is faithful: the condensed passage keeps my actual words (window computation, the negative pattern cannot false-positive, reverting trips both assertions, backed by the behavioural dispatch test, "Sound."), and the second-guard line "already earned its keep — it found the `gatherReviewContextFull` leak before shipping" is verbatim from my round-03 Q3. Finding A quotes my NIT-1 title verbatim and correctly attributes the severity disagreement ("agreed on the fact, disagreed only on severity"). Finding B records codex-1's MAJOR in their own words as UNWITHDRAWN and NOT adopted, with both sides' reasons stated and codex-1's strongest point — a constant-false branch is compiled, not verified — carried into the enablement gate rather than argued away. The record does not soften the losing position; it leaves the conflict open. Honest.

2. Nothing is claimed as fixed that is not fixed. The cycle-3 Finding A fix is present in source: `roundContextSentence()` (frontier.go:161-166) derives the sentence from the same `compactionEnabled` constant, so with it false the prompt says "Every prior round appears below in full." — true as shipped, and the false banner clause is gone. The rendered-prompt test (frontier_test.go:316-345) asserts the rendered string, asserts no banner is promised while the constant is off, and checks the participant-list slot, which catches the silent Sprintf-swap class the consensus names. The fixes I verified byte-for-byte in round-03 remain: `const compactionEnabled = false` (frontier.go:80), first-line `return full()` (frontier.go:56-58), BOM trim (frontier.go:99), and the ledger skip in `renderRound` (frontier.go:134). One nit of mine is correctly NOT claimed as fixed: the "compiled and exercised" wording in IMPLEMENTATION.md (my NIT-2). The consensus does not list it as applied, so no overstatement — it stands as an open follow-up.

3. "What ships" is accurate. No speedup ships because the only speed-producing path is unreachable by construction: a constant, referenced nowhere else, no env var, flag, config key, or file reaches it (I audited this repo-wide in round-03). Agent-visible behaviour is unchanged except the two literal deltas I documented: the prompt sentence (now truthful) and the `_ledger.md` exclusion, which no existing deck can trigger and which fails in the safe direction. codex-1's caveat on byte-identity claims covers that nuance and I concur with it.

4. The enablement gate is sufficient for the risk I raised. My round-2 MAJOR asked for a fail-closed latch stronger than a config flag; the constant plus the two-test tripwire is exactly that — flipping it turns the explicit tripwire assertion (frontier_test.go:113-114) red and the dormancy test red behaviourally, so enablement cannot happen silently and must pass through review, where the never-executed dormant path (my NIT-2 caveat: compiled and tripwired, not exercised) is necessarily re-examined. The gate's precondition — a validator including G3, G5 and G6 — is recorded at the constant's comment (frontier.go:76-79), in IMPLEMENTATION.md, and in this consensus. I endorse codex-1's reading as the correct interpretation of that precondition: "validator exists" means implemented with end-to-end and mutation coverage, not a placeholder artifact. That is a clarification of the gate, not a defect in what ships today.

My remaining concerns (NIT-2 wording, NIT-3 design-path dispatch test at enablement) are follow-ups already on record, not release blockers. The machinery is inert by construction, the prompt now tells the agent the truth, and the only conflict on record is preserved honestly rather than resolved by count. OK to release.
