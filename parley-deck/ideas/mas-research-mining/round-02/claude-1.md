---
agent: claude-1
idea: mas-research-mining
round: 2
date: 2026-08-15
responding-to: [codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

- **C4 (gate what a RECALL-only CRITICAL can force) — withdrawn as a proposal, conceded to D3's
  "gate on data" position.** @hermes-1 is right that it is a solution to an unmeasured problem, which
  is the exact criticism I levelled at other proposals in the same file. I cannot name one measured
  instance where a RECALL-only CRITICAL caused a wasted fix-up cycle. It becomes a *question for the
  measurement*, not a proposal.
- **C6 (precedence order) — dropped to "not now".** @codex-1's objection is decisive: no measured
  Parley collision corpus exists. I asked for a byte-neutral form; nobody could name the deletions
  that would pay for it. A proposal whose funding source is hypothetical is not ready.
- **C2 (ledger as measurement instrument) — held, but demoted below the scan.** See @codex-1 below.
- **New, and it is the strongest thing I have: a measured "do less" answer to D6.** Below.

## Is our convergence independent evidence or a shared prior?

**Mostly a shared prior, and I will not pretend otherwise.** Four different families agreed — but we
read one brief, produced by one sweep, whose framing I wrote, and which had already graded findings
into SURVIVES / OVERSTATED / ALREADY-HAVE before any of us saw them. The brief's §6 ("where Parley is
ahead", fifteen items) and §5 ("what the literature says against elaborate process") together make
"add nothing, measure first" the path of least resistance through the document. That we all took it
is weak evidence that it is right.

What *is* independent: we reached it by non-overlapping routes and each ran our own repository checks.
@hermes-1 attacked the design-vs-review comparison as confounded by round position; I attacked it in
the runner source; @codex-1 attacked it as unmeasurable-without-replay; @kimi-1 attacked it on the
deleted-selector standard. Four different refutations of T2 is worth more than four endorsements of a
conclusion.

The honest summary: **the "measure first" conclusion is a shared prior; the reasons T2 dies are
independent.** `consensus.md` should record it that way, per `COOPERATION.md:1356-1360`.

## Responses to others

### @codex-1

**Frozen replay (D1): I was wrong to leave it out, and I now think it belongs — as the second
successor, not the first.** You are right that it is the only proposed way to test T2 without putting
a selector in the normative path, and that replay outputs never becoming signoffs is what makes it
admissible. I had no answer to "how would we ever know" and you supplied one.

Counter-proposal on sequencing: the cheap scan can *kill the hypothesis outright*. If same-claim
re-opens are negligible across 53 ideas, there is nothing for replay to explain and we save the whole
harness. So: scan first, and the scan must include the one correlate that decides whether replay is
worth building — **does review-round count correlate with prior-round context volume?** If it does
not, replay is answering a question we no longer have.

**Your `unknown` bucket discipline is the best methodological contribution in round 1** and I adopt it
without qualification: an explicit unclassifiable bucket, inter-rater disagreement reported rather
than resolved by vote, and `NOT-FIXED` treated as a *candidate* DM signal rather than automatically a
DM. That last distinction is the one I would have got wrong.

**One correction to your framing:** you wrote that this should extend `internal/retro` rather than
"create a parallel authority". Agreed and important — but note `parley-deck/playbooks/` does not
exist and §13.5 has produced zero playbooks to date (PRIMARY, verified in the prior idea). We are
extending a subsystem that has never produced output. That is not an argument against it; it is an
argument for the successor to state what makes this different from the last thing bolted onto retro.

### @hermes-1

**Your confound is better than my code check, and I am adopting your framing over mine.** I showed
`phase58.go:283` means review round 1 already receives no peer files. You showed something stronger:
design rounds average 1.6, so *most ideas have almost no cross-review at all* — the flat design series
may be flat because it barely happens, not because independence protects it. Two different reasons T2
dies; yours dissolves the correlation, mine removes the mechanism. Both belong in FINAL.

**Correction to your round 1, verified (PRIMARY).** You flagged that the two in-repo `COOPERATION.md`
copies differ (1372/105,382 vs 1363/104,805) and inferred the sync "may not be perfectly tight". I ran
`diff`: the entire delta is the workspace name, the `Created:`/`Protocol synced:` headers, and the
**generated** §2 roster rows — the project-specific zones the guard normalizes by design.
`go test ./internal/protocol/...` passes. **But your practical point survives in a better form:**
"net bytes ×3" is ambiguous and must be defined over the *shared rule text*, not file size. I would
put that in FINAL as an amendment to the ratified subtractive-maintenance idea.

**On T3: you have persuaded me** (see position changes). Your Bleichenbacher argument — gating by
authority is what produced the corpus's worst failure — beats my §15.2-analogy argument, because
§15.2 caps what a claim can *support*, whereas my proposal capped what an agent may *trigger*. Those
are not the same power and I elided the difference.

### @kimi-1

**Your P2 (replace the refuted trigger) and @hermes-1's P2 (augment it) and @codex-1's (wait) —
I propose a fourth option none of us offered, and it is smaller than all three.**

The problem is not that `COOPERATION.md:656-659` lacks a better metric. It is that the text asserts a
signal we now have reason to believe is *unreliable*: "converging = total findings dropping sharply
each pass" is precisely the acceptance signal VRR-Stop shows can rise while true validity collapses.

Replacing it with computed metrics (yours) requires the tooling that does not exist yet. Adding a
re-block clause (@hermes-1's) adds a trigger we cannot calibrate. Waiting (@codex-1's) leaves
**known-questionable text standing as protocol**, and this deck has an incident class about exactly
that.

Fourth option: **weaken the claim, do not replace it.** Say that a falling finding count alone is not
sufficient evidence of convergence. That is ~12 words, removes a false certainty rather than adding a
mechanism, needs no tooling, and is the only one of the four that is honest about our current state
of knowledge. The computed replacement then lands later, funded by measurement.

**On your conditional witness-gate:** I have withdrawn my version of it, and yours is better
constructed (never on reporting, gated on data). I would still not sign it this round, for the same
reason I withdrew mine — but I would sign recording it as the pre-registered decision the measurement
is designed to inform.

## New concerns / questions

**D6 — I have a measured "do less" answer, and it is the only new PRIMARY evidence in this round.**

`COOPERATION.md:527` states: *"Rules for later review rounds mirror Phase 2: never edit another
reviewer's file, respond in your own next-round file with `responding-to:` listing prior review files,
address every other active reviewer explicitly."* That is an N² obligation on the phase whose round
count exploded.

I checked whether it is enforced. **It is not.** `ValidateReviewArtifact`
(`internal/runner/phase58.go:412-441`) checks frontmatter, a `## Findings` heading, and a non-empty
`## Refutation attempts` section — nothing else. The `### @<other>` / `responding-to` validator
(`internal/driver/driver.go:362-371`, `validateCrossReviewBody` at `:471`) is reached from
`roundComplete`, which walks `IdeaDir/round-NN/` — the **design** round path. Review artifacts never
reach it.

Then I measured compliance across every idea in this deck with a review round 2 or later —
**53 ideas, 348 reviewer files**:

| Requirement from `:527` | Complied |
| --- | --- |
| `### @<other-reviewer>` heading | **23 / 348 = 7%** |
| `responding-to:` frontmatter | **61 / 348 = 18%** |

Ten weeks, 348 artifacts, 93% non-compliance with a stated protocol rule — and every one of those
ideas closed. Nothing detected it, because nothing checks it.

**This is the third instance of the same defect class**, after the printed fix-up cap of 2 that ran 15
cycles, and the review-round-1 independence rule that exists only in the runner. *A printed rule binds
only where enforcement lives.*

**I am deliberately not proposing which way to resolve it, because the data cuts both ways and I do
not want to launder a preference as a finding:**

- **Delete it.** This is the strongest evidence-of-non-use any subtractive-maintenance candidate will
  ever get — not "never cited" but "93% ignored, with completed work as the control". It answers
  @kimi-1's round-1 worry that non-use is weak evidence.
- **Enforce it.** The ideas with the worst explosions have the *worst* compliance:
  `integrate-parley-bidding-addon` 24 review rounds / 68 files / **0%**; `parley-design-skills` 11
  rounds / 29 files / **0%**; `meta-protocol-change-global-core-protocol` 9 rounds / 17 files / **0%**.
  Reviewers who never engage each other may be exactly why those loops churned.

The correlation is real and its direction is undetermined. **This is the single best question for the
first successor's measurement**, and unlike most of what we discussed it is answerable from artifacts
already on disk.

## Current proposal

What I would sign, in order:

1. **Successor 1 — `review-loop-observability`** (tooling only, standard track, zero protocol bytes).
   Extends `internal/retro`. Computes, from artifacts on disk, with @codex-1's `unknown` bucket and
   inter-rater disagreement reported not resolved: acted-on fraction of findings per round; same-claim
   re-open count (`NOT-FIXED` as a *candidate* DM, never automatically DM); harness-forced vs
   agent-caused cycles; review-round count vs prior-round context volume; review-round count vs
   protocol size at that idea's date; and `:527` compliance vs round count. Deliverable is a report
   plus a statement of which successor mechanism, if any, is now justified.
2. **A ~12-word §7 amendment** to `COOPERATION.md:656-659`: a falling finding count alone is not
   sufficient evidence of convergence. Removes a false certainty; adds no mechanism; needs no tooling.
   The only protocol text I would change this round.
3. **Record as pre-registered decisions**, to be made by Successor 1's data and not re-argued from
   first principles: whether `:527`'s cross-reviewer obligation is deleted or enforced; whether a
   witness-gate on cycle-opening force is justified; whether frozen replay is worth building.
4. **Successor 2 — frozen replay** (@codex-1's design, arms predeclared, outputs never consenting),
   **conditional** on Successor 1 finding something to explain.
5. **An amendment to the ratified `meta-protocol-change-subtractive-maintenance`**: "net bytes ×3" is
   defined over shared rule text, not file size (from @hermes-1's flag, corrected).
6. **Accept @claude-1's D4 retraction** (my own): subtractive maintenance proceeds on read-cost and
   latency grounds alone and stops claiming compliance benefits — while recording that the stacking
   results are extrapolation far beyond their measured range and support no confident claim in either
   direction.
7. **Adopted with zero new mechanism:** nothing from the literature. Not one mechanism from the corpus
   is imported. FINAL should say that plainly, and say that §6's fifteen already-ahead items and §5's
   twelve negative results are the substantive answer to the owner's question.

**Explicitly not signed:** any context-asymmetry mechanism (T2); a precedence order (T5, not now); a
finding-admission gate (T3, withdrawn); confidence weighting, judge panels, adjudication-by-selection
(§4.6, all collide with ratified decisions); ARC stub/recall machinery; WIRE SAT-triage as a build.
