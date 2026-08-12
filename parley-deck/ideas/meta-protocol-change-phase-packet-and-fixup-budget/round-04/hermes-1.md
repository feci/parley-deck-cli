---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 4
date: 2026-08-11
responding-to: [codex-1/round-03, hermes-1/round-03, kimi-1/round-03]
---

## R1 — the three numbers

**1. Ship: 0.50.**

The 2× speedup is the effect size that motivated this idea — the packet exists to halve the protocol-read cost, which was measured as the top cost lever. @codex-1's overhead argument (round-03, condensed: the real packet has generator/index costs the excerpt arm lacked) explains why the measured ratio may be higher than the excerpt arm's, not why the bar should be lower. The overhead is part of the packet's real cost; the threshold gates the real effect. A 0.60 threshold pre-commits to shipping a 1.67× improvement — real, but not the target that justified the work. If overhead pushes the ratio to 0.55, the middle band handles it honestly. (Note: @codex-1 adopted @kimi-1's 0.60 at round-03 line 49, but @kimi-1 tightened to 0.50 in the same round at line 81. §15.3: the count doesn't settle it; the reason does.)

**2. Refute: 0.67.**

Below 1.5× speedup, the optimization doesn't justify a new generator, failure mode, and packet system. @kimi-1's 0.80 (1.25×) is too loose — a middle band that wide spends the experiment's credibility on a range where the answer is nearly always "ship anyway." @codex-1's no-middle-band binary is cleaner but pairs with 0.60, which I don't hold. 0.67 keeps the middle band narrow: the range where the result is promising but not conclusive.

**3. Middle band: return to user with the measured number.**

I drop my replan-and-re-run and adopt @kimi-1's (round-03 line 85, quoted literally: "the measured number replaces the 50–70% planning estimate verbatim in FINAL.md and the ship decision returns to the user with that number; it is never rounded up to the estimate"). My replan-and-re-run was post-hoc adjustment masquerading as rigor — "change the intervention and try again" is the optimization pre-registration exists to prevent. Return-to-user is a genuine pre-registered rule: the outcome is determined (human gets the number, decides), not open-ended.

## R2 — two ideas

I change my round-3 position (one idea, sequenced). Two things moved me:

1. @kimi-1's anchoring argument (round-03, condensed): `Timeout per agent`, `Reviewers (Phase 6)`, and `Review consensus (Phase 7)` have no known-correct per-track values. The audit's enforcement-path mapping for those cells is the same class of deliberation that cost D1 two rounds. That work should not delay the packet idea's closure.

2. The structural test is audit work, not packet work. My round-3 design put it in a second commit on the same branch, but a test that fails when a §4.0 cell has no enforcing path belongs with the audit. The packet idea's tests are the D1/D5 boundary tests (cycles 5/6, rounds 3/4) — already agreed, shipping with this idea.

This idea ships: the packet (ranks 1–2), D1 (5, inclusive fix, boundary test), D5 (3, `CapCrossReviewRounds` wired for deliberation, boundary test). The experiment runs against the packet commit before ship. FINAL records the known divergences (silent `MaxRounds: 4` default, cross-review defaulted to 1, app layer passing only `CrossReviewRounds`) as the follow-up's seed inventory (@kimi-1, round-03 lines 109–110). No further §4.0 cell edits land until the audit runs (@kimi-1, round-03 line 118).

Follow-up slug: `meta-protocol-change-track-gate-enforcement-audit` (@codex-1, round-03 line 75) — "enforcement" names the audit's purpose: it checks whether the table is code-enforced, not just whether it exists.

## C4 — last call

Nothing new beyond my round-3 list (lines 182–202). I sign a consensus that holds D1 = 5 (inclusive fix, boundary test), D5 = 3 (wired, boundary test), D2/D3 as settled, the R1 thresholds written into FINAL.md before the experiment runs and reported unchanged after, and R2 as two ideas with the conditions above.

One boundary from round-3 I hold: the threshold does not move after the data is seen. If the consensus adopts replan-and-re-run instead of return-to-user for the middle band, I still sign — provided the FINAL records the replan as a new experiment with a new pre-registration, not a continuation of the first.
