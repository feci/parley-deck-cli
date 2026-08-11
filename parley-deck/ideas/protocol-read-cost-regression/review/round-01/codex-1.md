---
agent: codex-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

BLOCK. The derived marker ledger is not acceptable for v1. It cannot establish that it retained
every live objection, and the implementation does not implement most of G3, knowingly fails G6,
and has not completed G7. Those are release gates, not deferred enhancements.

The narrow boundary from my signoff is otherwise respected in one important sense: no artifact
validator or close predicate reads this ledger, so ledger extraction does not determine whether an
artifact is valid or consensus has closed. The design-consensus drafter also still reads every round
artifact directly. However, `review-consensus` is routed through the compacting review-context path;
that is a consensus-drafting regression even though it is not a new consensus rule.

## Refutation attempts

- **G1 / design round 2:** pass. `frontierContext` calls the original full renderer for
  `round <= 2`, so design round 2 remains full-history.
- **G2:** the packet has the advertised mechanical shape for marker-friendly fixtures: the previous
  round is full and older bodies are absent. It does not prove safe compaction because the ledger is
  incomplete and line-based.
- **G3:** fail. Only two coarse absence cases fall back. The required missing/invalid/ambiguous/
  challenged/hash/locator/conflict states are not validated.
- **G4 / design consensus:** pass by direct call-path inspection.
  `buildConsensusDraftPrompt` still tells the drafter to read every `round-*` artifact and does not
  call `frontierContext`. The no-op test in `frontier_test.go` does not prove this.
- **G5:** fail outside the hand-picked `CRITICAL:` fixture. A real objection in this idea,
  `round-01/kimi-1.md:128` — “I reject two of the listed options as primary levers” — matches no
  marker and silently leaves round-3 context if it is not repeated in round 2.
- **G6:** fail, as `IMPLEMENTATION.md` acknowledges. There is no claim identity, verdict join, or
  conflict detector.
- **G7:** fail, as `IMPLEMENTATION.md` acknowledges. The fallback-banner reversion check was not
  completed.
- **`parley run` paths:** design round 2 and review round 1 retain the old full-context path; the
  existing targeted runner tests pass. The design-consensus drafter is unchanged. Review consensus
  after review round 2 is not unchanged: `round+1 == 3` activates compaction of review round 1.
- `go test ./internal/runner -run 'Frontier|Orphaned|Fallback|Ledger' -count=1` passed.
  The focused round/review tests also passed. `go test ./... -count=1` failed in the unrelated
  `TestDurableKillEndToEndRealProcess` because the fixture had no recorded boot ID; all other
  reported packages passed.

## Findings

### [CRITICAL] Marker extraction silently drops live dissent and cannot implement owner disposal

`ledgerMarkers` is a fixed, case-sensitive substring list (`frontier.go:54-74`), and
`extractItems` keeps only individual matching lines (`frontier.go:120-140`). That is not a
conservative extractor: it has no way to prove that an unmarked line is non-material. It misses
ordinary objection language such as “I reject…”, “I object…”, “I oppose…”, “I cannot accept…”,
“this is unsafe”, “this concern remains”, lower-case severities, and headings followed by the real
proposition in later paragraphs. The current corpus supplies a concrete counterexample at
`round-01/kimi-1.md:128`.

Even a matched finding is truncated to one line. A heading such as `### [MAJOR] Incomplete
validation` is carried while its following what/why/fix paragraphs disappear. The filename-based
`Owner` also cannot distinguish an author's own position from a quotation of somebody else's
objection. Any of those paths lets an objection leave the context without its owner disposing of
it, which converts omission into agreement under Phase 2 rule 1.

The lack of lifecycle is not an acceptable v1 simplification either. `ledgerItem` has no immutable
ID, state, or transition history (`frontier.go:43-49`), so the owner cannot record `RESOLVED` or
`SUPERSEDED`; recognized objections persist forever while unrecognized ones vanish immediately.
That produces both false consensus and false non-convergence.

**Concrete fix:** replace marker extraction with the participant-authored structured ledger from
the signed contract: stable owner-namespaced IDs, exact propositions, explicit owner identities,
lifecycle and append-only transitions, locators/hashes, claims/verdicts/evidence, and tombstones.
Treat that ledger as optional optimization input, not as an artifact-validity rule: if a complete
validated ledger is absent, use visibly announced full history. Add G5 cases using the real
“I reject…” wording, a multi-line review finding, a quoted peer objection, and an attempted
non-owner disposition.

### [CRITICAL] G3 fallback does not cover the mandated uncertainty states

The implementation comment says every uncertainty selects full history, but the code only falls
back when the entire previous round renders empty (`frontier.go:183-185`) or when all older files
together yield zero marker matches (`frontier.go:112-115`). The required states behave as follows:

| Required uncertainty | Actual behavior |
| --- | --- |
| missing | A wholly empty previous round or wholly empty extracted ledger falls back. A missing older round is silently skipped; a missing participant among otherwise present files is not detected. `frontierContext` receives no expected-participant set. |
| invalid | No ledger schema or transition validation exists. Filesystem read errors return an error instead of selecting and announcing full history. |
| ambiguous | Not detected. Substring matches are accepted as complete state; unmatched or multi-line meaning is silently discarded. |
| challenged | Not detected, and no challenged locator is expanded verbatim. |
| unresolved hash/locator | No SHA-256 is stored or checked. Generated line locators are never resolved or verified. |
| verdict conflict not marked `DISPUTED` | Not detected; no fallback fires. |

The single G3 test covers only “zero markers in all older text”, so its name materially overstates
coverage.

**Concrete fix:** make `frontierContext` accept expected participants and consume a parsed ledger
with a structured validation result. Validate completeness, IDs, ownership, lifecycle transitions,
hashes, locators, supersession/SECONDARY graphs, challenges, and verdict conflicts before compacting.
Any validation uncertainty must call the original full renderer and include the precise reason in
the prompt. Add one table-driven test per state above, including a partial-missing-participant case.

### [CRITICAL] G6 is knowingly unmet and permits false convergence over conflicting PRIMARY verdicts

The extractor merely carries lines containing `PRIMARY`, `WRONG`, `CONFIRMED`, or `DISPUTED`.
It has no claim record, stable claim identity, equivalence/join step, or relationship between a
verdict and the proposition it verifies. Rewording the same material claim under a new ID therefore
lets opposing `PRIMARY` verdicts survive as unrelated strings, with neither `DISPUTED` nor fallback.
That violates an explicit FINAL gate and can hide a §15.3 conflict from later participants.

Shipping without G6 is CRITICAL, not a v1 limitation: the acceptance criterion was created for the
exact false-consensus failure this optimization can introduce.

**Concrete fix:** implement authored claim and verdict records with immutable IDs, old/new claim
links, verifier/provenance/evidence fields, and deterministic conflict validation. If changed wording
cannot be safely joined to an earlier material claim, classify identity as ambiguous and fall back
to full history. Add the exact cross-round, reworded-claim/opposing-PRIMARY fixture required by G6
and prove both the `DISPUTED` and ambiguity-fallback branches.

### [CRITICAL] Review-consensus drafting loses undisposed finding bodies after review round 2

`buildPromptForRound` sends `review-consensus` through `gatherReviewContext` with
`roundNumber(opts)+1` (`runner.go:919-924`). For consensus after review round 2, that becomes round
3, so review round 1 is compacted. `BuildReviewConsensusPrompt` then tells the drafter to synthesize
the findings “below”, but older findings have been reduced to marker-bearing lines. Their evidence,
rationale, concrete fix, rebuttals, and disposition text can disappear even though the reviewer
never disposed of them.

The design-consensus path is correctly unchanged; the review-consensus path is not. This also
violates the FINAL packet rule “Consensus drafter — unchanged, full history” for Phase 7 and meets
the explicit CRITICAL condition for an objection leaving context without owner disposal.

**Concrete fix:** have the `review-consensus` case call `gatherReviewContextFull` directly for every
cycle. Restrict frontier selection to participant review rounds. Add an integration test with two
review rounds where a round-1 finding has a multi-line body and assert that the complete body reaches
the Phase 7 prompt.

### [CRITICAL] Content sniffing can remove FINAL.md and IMPLEMENTATION.md from later review prompts

`gatherReviewContext` decides whether fallback occurred with
`strings.Contains(rounds, "carry-forward fallback")` (`phase58.go:502-505`). `rounds` contains
participant-authored text. Any prior review that discusses the banner using that exact phrase makes
a normal compacted packet look like a fallback packet, causing the function to return without
prepending `FINAL.md` and `IMPLEMENTATION.md`.

This phrase is especially likely here because reviewers are required to assess fallback behavior.
At review round 3, or review consensus after round 2, the current implementation and its acceptance
criteria can therefore disappear from the prompt.

**Concrete fix:** return structured metadata from `frontierContext`, for example
`{Text, UsedFallback, Reason}`, and branch on `UsedFallback`; never infer control flow from
participant content. Add a test whose previous review contains the literal phrase while asserting
that FINAL and IMPLEMENTATION are still present exactly once.

### [CRITICAL] G7 has not been completed

`IMPLEMENTATION.md` explicitly records that only two of three required reversion checks completed;
the fallback-announcement reversion was left unverified after the harness was killed. G7 says every
fix must be reverted, the revert must apply and compile, and the relevant test must go red. The
current evidence does not satisfy that gate.

**Concrete fix:** rerun the missing mutation in an isolated temporary worktree or copied checkout
with kill-safe cleanup, record that the mutation applied and compiled, show the targeted test failing
for the intended assertion, restore the tree, and rerun the green targeted suite. Do not rely on a
language-level `finally` as the only restoration mechanism.

### [MAJOR] The tests do not exercise the changed call sites and one claimed G4 test is inert

`TestFrontierDoesNotTouchTheConsensusDrafterPath` contains no assertion; assigning two string
literals to `_` cannot fail if frontier compaction is later wired into consensus. The remaining new
tests call `frontierContext` directly with design-round fixtures. None exercises
`gatherReviewContext`, `review-consensus`, partial owner loss, the content-sniffing branch, or the
actual `buildPromptForRound` phase dispatch. This is why the two review-path CRITICAL defects remain
green.

**Concrete fix:** replace the inert test with call-path tests. Assert design round 2 and review round
1 use full context, design consensus retains its full-history instruction, review consensus uses
full review history, and review round 3 includes FINAL/IMPLEMENTATION plus safe frontier context.
Include all G3/G5/G6 adversarial fixtures above and rerun the full suite after resolving or
quarantining the unrelated durable-kill environment failure.

## Open questions

None. The safe v1 is clear: either implement the authored, validated ledger contract and all gates,
or keep full history. A best-effort substring extractor cannot safely sit between those choices.
