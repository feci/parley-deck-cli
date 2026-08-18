---
agent: claude-1
idea: mas-research-mining
round: 1
date: 2026-08-15
---

## Summary

**T2 is built on a false premise, and the premise is mine — I wrote it into `00-prompt.md`.** I
verified it before writing this file and it does not hold. That correction changes what the rest of
the idea should be about, so it leads.

Beyond that: three proposals, one reprioritization, one honest retraction of an argument this deck
has already agreed to act on.

## Proposed approach

### C1 — Correction: review round 1 is ALREADY cold-start. T2's Parley-side leg fails.

The convergence I put in the prompt — four unrelated lines pointing at *give reviewers less context,
not designers* — rests on a Parley-side leg: that design rounds stayed flat at 1.6 **because** round-1
independence is enforced there, while review rounds exploded to 5.1 **because** it is not.

That leg is false. **PRIMARY**, `internal/runner/phase58.go:283`:

```go
for r := 1; r < round; r++ {   // review round 1 → loop body never executes
```

`gatherReviewContext` sends `FINAL.md` and `IMPLEMENTATION.md`, then every *prior* review round. In
review round 1 there are no prior rounds, so **a round-1 reviewer already receives zero peer review
files.** Review round 1 is cold-start by construction of the runner, and has been all along.

What is actually missing is narrower: the **rule** is unwritten. `COOPERATION.md:324` states round-1
independence for design; `:527` extends Phase-2 mirroring only to "later review rounds"; `:895` says
the design rule "is a social one… no enforcement beyond agent discipline". So the protocol text has a
gap the runner does not.

**Consequences, and they matter more than the correction:**

1. **The 1.6-vs-5.1 split cannot be explained by round-1 anchoring.** Both phases start cold. Whatever
   drives 5.1 lives in review rounds 2+ and in the fix-up loop — which is where T3 and T4 point.
2. The brief's verifier flagged that no single line of the four-line convergence is strong and that
   the *agreement* was the signal. One of the four was Parley's own split. **Removing it does not
   weaken the convergence by a quarter; it removes the only leg that was measured on our workload.**
   The remaining three are a vendor claim with no isolated measurement, a paper that explicitly ran no
   ablation, and an ICML result about team size. T2 should be closed, not pursued.
3. What survives is a one-line documentation fix: state the review-round-1 independence rule so the
   text matches the runner. Near-zero bytes, no proof burden, and it prevents a future runner change
   from silently removing a property nobody wrote down.

I am the reason this was in the prompt as a live tension. Round 1 is where that gets caught.

### C2 — Reprioritization: the ratified-but-unbuilt ledger is the MEASUREMENT instrument (rank 1)

The objection ledger (`protocol-read-cost-regression/FINAL.md:53-100`) has been carried in this deck
as a context-transmission fix. The brief shows that is the smaller half of its value.

T4's decomposition needs to identify **the same claim re-raised after a fix**. T1 needs **acted-on
fraction of findings per round**. Both require exactly what the ledger already specifies and nothing
else in this system has: an immutable **owner-namespaced ID** and an **exact scoped proposition,
never a generated paraphrase**, with an append-only transition history.

Without stable claim identity you cannot tell a re-BLOCK on the same claim from a new finding, so
P(DM|D=1) — reported at 53–94% across 14/14 cohorts — is **not merely unmeasured here, it is
unmeasurable**. Same for "was this concern still present at consensus", which the brief notes is
proposed twice by different routes.

This does not change the ledger's design or its two blocking preconditions (the v1.43.1 failure
analysis; the shadow receipt). It changes its **justification and its rank**: it is the prerequisite
for T1 and T4, not a nice-to-have for context size. Everything else I propose below is cheaper, and
also weaker.

### C3 — T1: instrument only what is already on disk (rank 2)

The brief's §9 list is the most uncomfortable thing in it. My position is that "measure first" is
correct **only** for metrics computable from artifacts we already have, with no new protocol text and
no new agent behaviour. Anything requiring re-runs (Pass², a compute-matched single-agent control) is
real cost and should be deferred, stated as deferred, not quietly dropped.

Computable today, tooling-only, zero protocol bytes — extend `parley retro`, which already assembles
per-idea signals (rounds, review rounds, fix-up cycles, dismissed counts, escalations):

- **acted-on fraction of findings per review round** — dispositions are already recorded;
- **findings per round by severity**, to test §6.5's "dropping sharply each pass" against the
  trajectory it actually claims;
- **rounds-per-idea plotted against protocol size at that idea's date** — both are in git.

That third one is the one I care about. It is the only cheap way to separate two hypotheses that the
brief puts in direct competition and that this deck has never distinguished: review rounds tripled
because *the work got harder*, or because **the protocol got 112% bigger and compliance degrades with
instruction count by silent omission**. If the second is even partly true, then subtractive
maintenance is not a read-cost nicety — it is the intervention. If it is false, we stop citing 1.6→5.1
as evidence for anything, because as the brief notes, that series is compute-confounded.

### C4 — T3: do not gate reporting; gate what a lone RECALL claim can *force* (rank 3)

§15.4's carve-out (`:1321`) is deliberate and I would not touch it. The corpus's most damaging results
— unanimity on a non-existent padding oracle, 63.6% correct-to-wrong conformity, consensus inflation
— are all failures of suppressing a minority. An admission bar on *reporting* buys round count by
spending exactly the thing this protocol is best at.

But reporting and *triggering* are different powers, and the protocol already separates them
elsewhere: `:1274-1282` caps a RECALL-only material claim at `UNVERIFIED` for `FINAL.md`. The same cap
is simply not applied to Phase 8. So: a `CRITICAL` supported only by `RECALL` is reported, recorded,
disputed and answered — but does not by itself open a fix-up cycle; it opens a **verification
obligation** on its raiser. That is the existing §15.2 ladder applied one phase later, not a new
gate, and it is the only lever in the survey aimed at round *count* rather than review *quality*.

Cost of being wrong: a real defect known-but-not-yet-provable gets slowed by one step. That is a
genuine cost and I do not want to minimize it — it is why this is rank 3 and not rank 1.

### C5 — Retraction: the compliance argument for shrinking the protocol does not survive

This one cuts against work this deck agreed to on 2026-08-14
(`meta-protocol-change-subtractive-maintenance`, from `cognee-mechanism-mining`).

The instruction-stacking result has a negative half that we would have to ignore to keep our story:
rewriting an instruction set tighter recovers compliance **+11.0pp / +3.3pp / −1.2pp** across a weak→
strong ladder, with ρ(raw follow rate, recovery) = **−0.85, p = 0.004**. The better the model, the
less consolidation buys — and it went *negative* for the strongest.

So: **the read-cost argument for shrinking survives** (3.3× median wall clock is measured on our own
workload). **The compliance argument does not**, and I have used it. Subtractive maintenance should
proceed on the byte/latency case alone and should stop claiming it will make agents follow rules
better. If C3's third metric comes back showing rounds track protocol size, that changes — but that is
a measurement we have not made.

Honesty requires the counterweight too: both instruction-stacking results are extrapolation far beyond
their measured range (depth 20, or 500 string-matched keyword instructions, versus our hundreds of
conditional cross-referencing prose rules). Neither supports a confident claim in *either* direction.

### C6 — T5: a document-wide precedence order, only if it is byte-neutral (rank 4)

Real defect: no precedence order across 18 sections, and §15 has already had to declare pairwise
subordination twice (`:1235`, `:1322`) while §4.0 carries its own override clause (`:233-237`). Those
scattered clauses are the symptom — you write local overrides when you have no global order.

The largest deployed agent constitution resolves conflicts with one document-wide priority order
rather than by deduplication, which is a genuinely different move from anything we have tried. **But
its evidence tier is first-party design rationale with zero empirical support**, and the brief caught
the reasons-over-rules half being cited with an omission that inverts it.

So I would propose this **only** in the form where it pays for itself: add the order, delete the
scattered override clauses it subsumes, and require the net byte delta ×3 to be ≤ 0. If the deletions
do not cover the addition, drop it. Rank 4 because the evidence is weak and the defect is currently
theoretical — I know of no incident caused by cross-section ambiguity, and I did not find one.

### What I would drop

- **T2** (see C1).
- **Everything in brief §4.6** — confidence-weighted aggregation, judge panels, adjudication by
  selection. All four ask us to reverse ratified decisions (`:1259` no self-verdicts; `:1294` never by
  count), all four were graded OVERSTATED, and the martingale theorems require agent homogeneity that
  a four-family roster violates. The one exception worth a sentence in FINAL: **CDP localization**
  does not collide with our standard — it selects what an *adjudicator* reads after the fact, from
  artifacts that stay on disk, so no agent's input is reduced.
- **WIRE's SAT-triage** as a build. The diagnosis (collisions resolve silently, 64.6% non-jointly, no
  signal emitted) is worth carrying into subtractive maintenance as motivation. Building clause
  encoding plus pairwise SAT over prose rules is a research project with a 0.55% yield, and our
  duplicate-rule instances were found by four agents reading the file.
- **ARC's stub/recall machinery.** The genuinely new idea in it is small and I would keep only that:
  an **audit log of what was and was not opened** changes the unverifiable step from "the summarizer
  preserved every objection" into "agent X never opened §ab12c3" — a logged fact. That belongs to the
  shadow receipt in C2's preconditions, not to a new mechanism.

## Concerns / open questions

1. **The brief is 87,273 bytes and I am handing it to four agents in an idea about over-transmission.**
   I did it deliberately — a research brief is the one artifact where completeness beats brevity, and
   it lives on disk to be read rather than pasted into every later round. But if a participant thinks
   this invalidates the round, say so now rather than at consensus.
2. **Almost nothing in the corpus tests our workload.** Zero sources measure multi-turn software
   design deliberation; one measures real PRs and measures whether comments were *acted on*, not
   whether they were *correct*. Every number in the brief is a transfer. I have tried to rank by
   *what is cheap and reversible* rather than by *what has the best paper behind it*, and I think that
   is right, but it is a judgement call worth attacking.
3. **C3 is the load-bearing proposal and it is the least exciting.** If the deliberation converges on
   a mechanism instead, I want the reason stated explicitly, because the brief's §9 says we have been
   changing this protocol for ten weeks without measuring whether any change helped.
4. Does anyone read the four-line convergence in T2 as surviving C1? I do not, but I built the prompt
   and I may be over-correcting against my own error.

## Risks

- **Instrumentation as procrastination.** C3 could become a reason to ship nothing for a month. Guard:
  it is tooling-only, it uses artifacts already on disk, and it should be timeboxed — if the three
  metrics are not computable from existing files, that is itself the finding.
- **C4 slows a true defect.** Named above; it is why it ranks third.
- **The ledger is unbuilt for a reason nobody has established.** C2 raises its priority but does not
  touch its blocking precondition. If the v1.43.1 failure analysis says the design is unimplementable,
  C2 and much of T1/T4 collapse together, and we should know that before building on it.
- **My own correction bias.** C1 is a correction of my own prompt, and there is a temptation to
  over-weight it to demonstrate rigor. The check on me is whether the other three participants
  independently reach C1 from the same source line.
