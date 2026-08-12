# Round 4 — two residuals, then consensus

Round 3 closed **D1 (cap 5, unanimous)**, **D5 (cross-review cap 3)**, **D2**, **D3**, and the
correctness veto. @hermes-1 adopted 5 and re-ran the distribution itself across all 69 ideas
(n=69, every value above 5 in {9, 14, 15, 25}); @codex-1's scoped deliberation-only count agrees.
**Do not reopen any of it.**

Two things remain. Both are narrow. Answer only these; a long file is a wrong file.

## R1 — the pre-registered threshold numbers

You converged on the unit (median wall-clock ratio, packet arm over full arm), the phases (1 and 6),
n=6 paired runs per phase, the canary as a standalone veto, live hash-bound source, and the
implementer running it before release. **Three numbers still differ:**

| | Ship if | Refute if | Middle band |
| --- | --- | --- | --- |
| @codex-1 | ratio ≤ **0.60** both phases | otherwise | none — binary |
| @kimi-1 | ratio ≤ **0.50** both phases | ratio > **0.80** either phase | returns to the user with the measured number; never rounded up to the estimate |
| @hermes-1 | ratio ≤ **0.50** both phases | ratio > **0.67** either phase | replan the packet's section set and re-run |

Note for @codex-1: you adopted "@kimi-1's `0.60` cut" in round 3, and @kimi-1 tightened to 0.50 with
a three-way gate in the same round — the same parallel-write that produced the 5-versus-6 detour.
Say whether 0.60 survives that, or adopt 0.50.

Settle, in one short answer each:

1. **Ship threshold** — 0.50 or 0.60. One number.
2. **Refute threshold** — 0.67, 0.80, or "no middle band". One number.
3. **If there is a middle band, what happens in it** — @kimi-1's return-to-user-with-the-number, or
   @hermes-1's replan-and-re-run. These are materially different: one hands the decision to the
   owner, the other spends another experiment cycle first.

§15.3 applies: two participants holding 0.50 does not settle it. Give the reason the number is
right, or adopt the other and say what changed your mind.

## R2 — scope of the §4.0 audit

- @codex-1: named follow-up `meta-protocol-change-track-gate-enforcement-audit`; this idea keeps the
  packet plus the two cap cells; the follow-up blocks any claim that the whole table is enforced.
- @kimi-1: named follow-up `meta-protocol-change-track-gate-audit`; same split; adds that no further
  §4.0 cell edits land until the audit runs, and that this idea's FINAL records the found
  divergences as the follow-up's seed inventory. Its reason: `Timeout per agent`,
  `Reviewers (Phase 6)` and `Review consensus (Phase 7)` have **no known-correct per-track values**,
  so choosing them is the anchoring work D1 just cost two rounds, times N cells.
- @hermes-1: in scope, sequenced — packet first commit, audit second commit on the same branch;
  survey every cell, publish divergences, fix only the two agreed cells, each other divergence gets
  a FINAL entry with a recommendation.

@hermes-1's version and the follow-up version differ less than they look: both enumerate every cell,
both publish the divergence list, both fix only the two agreed cells. The real question is whether
the **structural test** and the enumeration ship in this idea's diff or the next one.

Settle: **one idea or two.** If two, agree the follow-up's exact slug. If one, say what happens to
the pre-registered experiment if the audit's diff grows — the experiment measures the packet, and a
larger diff around it is a confound.

## C4 restated — last call

Anything that would stop you signing a consensus, that you have not already listed. Round 5 is
signoff. An objection raised there that could have been raised here costs everyone a cycle.

## Constraints

- Read all three `round-03/*.md` files before answering.
- Read-only outside your own `round-04/<agent-id>.md`. No git write commands.
- English only. Redact obvious secrets.
- §15.2: tag provenance. Untagged reads as `RECALL`.
- Quote literally or say you condensed.
