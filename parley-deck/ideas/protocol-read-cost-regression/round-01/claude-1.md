---
agent: claude-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

**I withdraw the working diagnosis I put in `00-prompt.md` as the primary cause.** I attacked it as
that document asked others to, and it does not survive: protocol read cost grew 2.1×, but review
volume grew **7.2×**, and design rounds did not grow at all. The protocol's growth is real and is a
multiplier, but the driver is that Phase 6–8 now runs five times as many rounds as it did in May,
with an unbounded cap on the `deliberation` track. I was measuring the wrong phase.

## Proposed approach

### The measurement that changed my mind (PRIMARY)

Across all 76 ideas in `parley-deck/ideas`, split at 2026-07-01:

```
                 ideas   mean design rounds   mean review rounds   max   mean review bytes
May–Jun (older)    55            1.4                  1.6            5         20,237
Jul–Aug (recent)   21            1.6                  5.1           24        146,290
```

Design rounds are **flat**: 1.4 → 1.6. Review rounds are **3.2×**, and the artifact volume those
rounds produce and re-read is **7.2×**. The worst cases are not outliers so much as the new shape:

```
2026-07-29  24 review rounds  699,565 bytes  integrate-parley-bidding-addon
2026-07-29  21               386,665         skills-cli-install-path
2026-07-28  11               467,514         parley-design-skills
2026-08-07   9               301,506         meta-protocol-change-global-core-protocol
```

Set that against the protocol's own growth over the same period — 720 → 1,359 lines, ~12,300 →
~26,100 tokens, 2.1×. Both are real. They are not the same size.

### Why my original framing misled me

I measured one deliberation's **design** side — three rounds — and generalized. That idea happened
to have *zero* review rounds because it stopped at FINAL. So I sampled the half of the lifecycle
that did not grow, found a real 2.24× compounding inside it, and presented it as the answer.

The compounding I found is genuine and it applies with more force to review, where round counts are
several times higher. But "the protocol is 71% of the round-1 read" is a fact about round 1
specifically, and round 1 is the cheapest round there is.

### What I think is actually happening

Three candidate mechanisms, in the order I would test them:

1. **The `deliberation` track has no fix-up cap.** §4.0's table caps fix-up at 1 cycle for `fast`
   and 2 for `standard`, and leaves `deliberation` **unbounded**. Combined with "repeat review until
   there are zero agreed fixes", any reviewer who keeps finding something keeps the loop open.
   A 24-round review is not an accident under that rule; it is the rule working as written.
2. **Refutation-default review rewards finding something.** The discipline is correct and I do not
   want it removed — it is what caught three of my own errors this month. But it has no notion of
   *diminishing severity*: a MINOR found in round 19 reopens the cycle exactly as a CRITICAL in
   round 2 does.
3. **Read cost multiplies whatever the round count is.** This is my original diagnosis, demoted to
   its proper place: it is the coefficient, not the exponent.

If (1) is the driver, the fix is bounded and cheap. If (2) is, the fix is a severity-aware
termination rule, which is a protocol change and needs care.

### What is load-bearing per phase

Round 1 needs: §1 scope, §4 phases (the round-1 shape), §6 conflict-avoidance, §15 provenance. That
is a small fraction of the document. It does **not** need §11 transport mechanics, §12 pipelines,
§13 retrospective optimization, Appendix A, or most of §4.0.1's LE rules.

Review rounds need: §4 Phase 6–8, the severity vocabulary, §15, §6. Again a fraction.

My recommendation for mechanism is therefore **phase-scoped views generated from the core**, not a
hand-maintained "short version" — a second hand-maintained document is the drift problem this deck
has already been bitten by twice. The overlay work just built the machinery for exactly this shape:
a deck file that is a generated view of a core, with a report of anything the generation drops.

### On the compounding re-read

Round N re-reading every prior round is what makes round 3 cost 2.24× round 1. I do not think it
can be replaced with a digest for **design** rounds: cross-review requires addressing what an agent
actually wrote, and a digest is where a position gets flattened into agreement it never gave. That
matters — §15.6 exists because correlated agreement is hard to detect even with the full text.

For **review** rounds I think it can, and that is where the volume is. A reviewer in round 12 does
not need rounds 1–11 in full; it needs the open findings, their severities, and what changed since.
That is a mechanical digest, not a judgement call.

## Concerns / open questions

- I have measured *artifact counts and bytes*, not wall-clock. The owner reports feeling slower, and
  I have not established that felt slowness tracks bytes rather than, say, agent-side latency
  changes. Somebody should check whether the agent CLIs themselves got slower over this period —
  that is a competing explanation I cannot rule out from this repository.
- "Review rounds" here counts directories under `review/`. An idea that reused a directory would
  under-count. I did not verify the counting against the round files themselves.
- The 24-round case may be pathological rather than representative, and the mean is only 21 ideas.
  I would not want a protocol change justified on one outlier.

## Risks

- **Cutting read cost by cutting rules.** §15, §7, §6 rule 3 and §14 were each bought with a real
  failure. A "shorter protocol" that drops them is not an optimization, it is a regression with a
  performance justification.
- **A hand-maintained short version drifts.** Two separate incidents in this deck came from a second
  copy of protocol text going stale. Any digest must be *generated*, with a report of what it omits.
- **My own framing.** I wrote the prompt, chose the measurements, and stated the diagnosis. I have
  now corrected it with a measurement I could equally have taken before opening the idea, and did
  not. Treat the new diagnosis with the same suspicion as the old one; it is still mine, and it is
  one afternoon old.
