---
idea: mas-research-mining
status: final
drafted-by: claude-1
date: 2026-08-15
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
rounds: 2
signoffs: [claude-1 ✅, codex-1 ✅, hermes-1 ✅, kimi-1 🟡]
---

# FINAL — the literature supplied no mechanism; the reading supplied a measurement

## The answer

**Nothing from the surveyed multi-agent-collaboration research is adopted.** Six research lenses,
~48 findings, an independent adversarial verifier per lens, four participants, two rounds. Zero
imports.

Three reasons, in descending order of how much they should change future behaviour:

1. **No source measures our workload.** Zero measure multi-turn software design deliberation. One
   measures real pull requests, and it measures whether a comment was *acted on*, not whether it was
   *correct*. Every number in the corpus is a transfer, and most metrics in it (oracle gap, expected
   correctness, Pass², P(DM|D=1)) require a ground-truth label that a design document does not have.
2. **Where the field and this protocol meet, we are mostly ahead** — fifteen places, several
   established by verifiers *trying and failing* to find a gap. Our provenance grading is verifiable
   where the literature uses self-reported confidence, which one source measures as a useless filter.
   Our refusal to resolve by count, including unanimous count, stands against the corpus's most
   damaging results — all of which are failures of counting. Our completeness standard is **stricter
   than the entire surveyed industry**: LangGraph's default reducer silently discards one of two
   concurrent writes, MetaGPT's publish-subscribe has never been ablated, and the best statistical
   alternative offers only a population-level bound.
3. **The strongest results in the corpus argue for doing less, not more** — and the best-evidenced
   intervention ceiling anywhere in it is +15.6%, with its peer-reviewed authors stating plainly that
   this is not enough.

## What the reading actually produced

Asking "what does the field have that we don't?" sent four agents back through our own record. That
is where the value came from.

### A rule we have printed and not enforced for ten weeks

`COOPERATION.md:531` requires that in later review rounds each reviewer address every other active
reviewer explicitly and carry `responding-to:`. Nothing enforces it: `ValidateReviewArtifact`
(`internal/runner/phase58.go:412-441`) checks frontmatter, `## Findings` and a non-empty
`## Refutation attempts`; `validateCrossReviewBody` (`internal/driver/driver.go:473`) is reached only
from `roundComplete`, which walks the **design** round path.

Measured across every idea in this deck with a review round ≥ 2 — **53 ideas, 348 reviewer files**:

| Requirement from `:531` | Complied |
| --- | --- |
| `### @<other-reviewer>` heading | **23 / 348 = 7%** |
| `responding-to:` frontmatter | **61 / 348 = 18%** (55 non-empty) |

Every one of those ideas closed. **This is the third instance of the class *a printed rule binds only
where enforcement lives*** — after the printed fix-up cap of 2 that ran 15 cycles, and the
review-round-1 independence property that exists only in the runner.

This measurement entered at round 2, too late for cross-review. It was therefore put to the signoff
as a challenge, and **independently re-run by all three other participants with different tooling,
reproducing 53 / 348 / 23 / 61 exactly**; all three independently confirmed the enforcement gap.

**Its direction is undetermined, by unanimous reading**, and it is Successor 1's first question:

- **Delete it** — the strongest evidence-of-non-use a subtractive candidate will ever get: not "never
  cited" but "93% ignored, with completed work as the control".
- **Enforce it** — the ideas with the worst round explosions have the *worst* compliance:
  `integrate-parley-bidding-addon` 24 review rounds / 68 files / **0%**; `parley-design-skills`
  11 / 29 / **0%**; `meta-protocol-change-global-core-protocol` 9 / 17 / **0%**.

### An invalid inference standing as protocol

`COOPERATION.md:656-659` presents "total findings dropping sharply each pass" as what *converging*
looks like. That is precisely the acceptance signal a source shows can rise monotonically while true
validity collapses. A workload count is not evidence of artifact validity.

### We have never measured whether any of this works

Ten weeks of protocol change, and the deck has never measured: acted-on fraction of findings per
review round; same-claim re-opens; harness-forced versus agent-caused cycles; run-to-run
reproducibility; outcome quality normalized by spend; or a compute-matched single-agent control.
**The 1.6 → 5.1 review-round series that motivated the read-cost work and the fix-up budget is itself
compute-confounded** — protocol size grew 112% over the same window.

## Outputs

**S1 — `review-loop-baseline`** (immediate; tooling only; **zero protocol bytes**). Read-only over
canonical artifacts, run events and git history; extends `internal/retro`. Raw counts *and*
denominators, an explicit **`unknown` bucket**, a capped **double-coded sample** with **inter-rater
disagreement reported, not resolved by vote**, and a stated go/no-go for S2 and for any admission-bar
successor. **Pre-agreed failure condition: if the unknown rate or annotator disagreement is too high
to support the decision, stop — do not add protocol fields to make the metric computable.**

**S2 — frozen replay** (conditional). Replay completed review phases with predeclared arms (full
context / cold-start / one reviewer at matched spend); outputs never become signoffs, dispositions or
consent. The only proposed way to test reduced reviewer context without putting a selector in the
normative path. **Blocked on both** S1 finding something to explain **and** a real spend-matching
method existing — runners emit no `agent.usage` (`internal/driver/loop.go:174-175`), so "matched
budget" is currently unimplementable and must not be silently replaced with wall clock.

**S3 — a small subtractive §7 change**: replace the `:656-659` illustrative sentences with
validity-neutral wording at **net-negative shared rule bytes**. No computed metric replaces them.
Withdrawing an invalid inference is not installing a mechanism, so it needs no baseline.

**Zero-byte fix**: record the review-round-1 cold-start property as a **code comment** at
`internal/runner/phase58.go:286`, not as protocol text — writing it normatively would pre-decide the
context-asymmetry question, and a future runner change then shows up as a comment/code mismatch.

**Amendment to the ratified `meta-protocol-change-subtractive-maintenance`**: "net bytes ×3" is
defined over **shared rule text**, not file size. The two in-repo copies legitimately differ (1372 /
105,382 B vs 1363 / 104,805 B) — entirely in the workspace name, the `Created:`/`Protocol synced:`
headers and the **generated** §2 roster rows, which the guard normalizes by design;
`go test ./internal/protocol/...` passes.

**Pre-registered decisions**, to be made by S1's data and not re-argued from first principles:
delete-or-enforce `:531`; whether a force-gate on cycle-opening (never on reporting) is justified;
whether replay is worth building; and whether `MaxFixupCycles` on `standard` should drop from 2 to 1.

## Rejected, with reasons

- **Asymmetric reviewer context** — refuted four independent ways: the runner already gives review
  round 1 no peer files (`phase58.go:286`); round 1 has no prior rounds to send, so context asymmetry
  is trivially impossible there; it is unmeasurable without replay; and it sits exactly where the
  deleted context selector sat.
- **DC/DM as a label** — a later `NOT-FIXED` or same-locator re-block is an *observable failed-repair
  signal*, not an oracle-backed claim that a correct artifact was made incorrect. Report the
  observable and `unknown`.
- **An admission bar on findings** — gating by authority is what produced the corpus's worst failure
  (ten agents unanimously confirming a non-existent vulnerability). The distinction between gating
  *reporting* and gating *cycle-opening force* is real, but no measured instance of the harm exists.
- **A document-wide precedence order** — witness-gated: it opens only if the ratified subtractive
  rule inventory produces a concrete cross-section collision, and then only byte-neutral.
- **Confidence-weighted aggregation, judge panels, adjudication-by-selection** — each asks us to
  reverse a ratified decision (`:1259` no self-verdicts, `:1294` never by count); the martingale
  theorems require agent homogeneity a four-family roster violates.
- **ARC stub/recall machinery, WIRE SAT-triage as a build, temporal search, NodeSets/ACL, memory
  provenance visualisation** — no measured problem, or barred by the completeness standard.

## The compliance claim, retracted symmetrically

The argument that shrinking the protocol improves rule-following is **unverified in both
directions**. The instruction-stacking negative half (+11.0pp / +3.3pp / **−1.2pp** across a
weak→strong ladder; ρ = −0.85, p = 0.004) kills the benefit claim; the same
extrapolation-beyond-measured-range limit kills any harm claim. Subtractive maintenance proceeds on
the repository-measured **read-cost and latency** case alone.

## Method notes worth keeping

**Our convergence was mostly a shared prior.** Four model families agreed, but read one 87 KB brief
with one framing, pre-graded before anyone saw it. What was independent: T2 died four separate ways.
Recorded per `COOPERATION.md:1356-1360`.

**Citation integrity of the sweep: zero fabricated sources** across ~30 papers. But roughly six
numeric or attribution errors across ~48 findings, and **every one favoured the finding it
supported**. Adversarial verification caught them; a single-pass research agent would not have.

**Six locator and attribution corrections were caught at signoff**, all in the drafter's text, all
upheld: `:527`→`:531`; `phase58.go:283`→`:286`; 61 `responding-to:` of which 55 non-empty; a
misparaphrase of @hermes-1's T2 argument; @kimi-1 credited as originator of a form it adopted from
@claude-1; and a @codex-1 formulation merged into a @kimi-1 locator. None changed a decision. The
signoff step is where this deck's attribution errors get caught — asking each signer explicitly to
audit what is attributed to it is what makes it work.
