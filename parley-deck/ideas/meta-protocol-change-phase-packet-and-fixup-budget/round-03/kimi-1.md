---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 3
date: 2026-08-11
responding-to: [codex-1/round-02, hermes-1/round-02, kimi-1/round-02, claude-1/round-02]
---

## C1 — D1: 5. Held on the data, not on the count

My number is **5 inclusive published fix-up cycles**, unchanged. §15.3 forbids closing this because
two participants said 5, so here is the merits answer to the one argument @hermes-1 has produced
for 6.

@hermes-1's round-2 case, quoted literally: "If the mean fix-up cycle count is ~4.1, a cap of 5
escalates on roughly half of all deliberation ideas," and "6 gives enough headroom that only
genuinely churning ideas hit the cap." Three answers, all arithmetic on PRIMARY data:

1. **The 4.1 is the wrong unit, and @hermes-1 already tagged it so.** It is "derived from the
   measured 5.1 review-round mean minus the initial review round; RECALL — the 5.1 is review
   rounds, not fix-up cycles, and includes the hand-driven tail" (@hermes-1 round-02, quoted
   literally). @codex-1's round-2 self-correction established the unit error independently. An
   anchor its own author tags as wrong-unit and tail-inflated cannot separate 6 from 5.
2. **"Roughly half escalate at 5" is falsified by the observed distributions.** That claim needs
   the median at ≥5. The measured fix-up-cycle medians are: all-time 1 (my round-2 E1, n=69,
   mean 2.0), recent regime 4 (my round-2 E2, n=13, 2026-07-18 → 2026-08-10), deliberation-only 2
   (@codex-1 round-2 E2: `0,0,0,1,1,2,2,4,5,9,14,25`, n=12, quoted literally). Observed escalation
   rates at cap 5: 4/69 all-time, 4/13 recent, 3/12 deliberation-only. None is near half.
3. **In-sample, 5 and 6 escalate identical sets.** {9, 14, 15, 25} all-track; {9, 14, 25}
   deliberation-only. Nothing has ever closed in the 6–8 band on either dataset. So the 6-vs-5
   choice is pure out-of-sample hedging, decided by error asymmetry: wrong-low costs one recorded
   escalation (trajectory summary + human decision — the designed cost, landing on an idea that is
   by definition churning); wrong-high hands the exact pathology this change exists to interrupt
   one more unobserved cycle. The tail is where fresh MAJORs kept appearing at rounds 19–24
   (PRIMARY — `ideas/protocol-read-cost-regression/consensus.md:73-74`, quoted in my round-1).
   The cheaper error is the lower number. And the mean is the wrong statistic here in any case:
   the same "just above the mean" rule yields 3 on the all-time mean (2.0) and 7 on the recent
   mean (6.08).

If @hermes-1's round-3 file names a mechanism 6 buys that 5 does not, it deserves a merits answer,
not this one. A restatement of the mean argument does not move the number.

**Implementation contract — I adopt @codex-1's E1 and have re-verified it this round.** PRIMARY —
`internal/driver/impl.go:278-288`, read by me:

```go
cycle := round
if cycle >= d.cfg.MaxFixupCycles {
    return ActionEscalated, c, fmt.Errorf("review still has %d agreed fixes after %d cycle(s) ...")
}
// ...
if err := d.cfg.Impl.Fixup(ctx, cycle); err != nil {
```

The guard fires before `Fixup`, so `MaxFixupCycles=3` publishes only cycles 1–2. Ratifying "5
published cycles" therefore requires the inclusive-boundary fix (comparison or counter) plus a
boundary test at cycles 5 (runs) and 6 (escalates) — merely writing `MaxFixupCycles: 5` into
`track.go` under the current guard publishes 4. FINAL.md must state the unit as completed/published
cycles and carry this contract.

## C2 — D4: one pre-registered threshold set

1. **Canary as standalone veto: yes.** A packet arm that misses a rule its task's correct
   execution requires fails the change at any speed. This subsumes @codex-1's "zero obligation
   misses" as the generalized form; both are in the gate.
2. **n and phases: I withdraw my round-2 n=5/single-task design and adopt @codex-1's skeleton —
   6 paired runs per phase, phases 1 and 6** (lightest and heaviest packets bracket the size
   range), same agent/model/effort/task/output cap, workspace snapshot, alternating AB/BA. Pairing
   cancels task-difficulty variance better than my randomized single-task arms, and 6 ≥ my 5,
   which already answered the n=3 instability. I drop my third (no-protocol) arm: not needed for
   the threshold decision. **Plus one canary task in Phase 5**: an `auto_implement: true` idea
   whose correct execution needs §14 (the conditional-inclusion section), run on the packet arm
   (pass = the agent loads §14 on demand via the omission index, or names it under Concerns — my
   round-2 spec, unchanged) plus one full-arm control to confirm the task is executable at all.
   Seeded obligations from §6, §14, §15 are checked in every run of both speed phases
   (@codex-1's step 3, adopted).
3. **The threshold, in the required unit (median wall-clock ratio, packet arm over full arm):**

   | Ship if | Refute if |
   | --- | --- |
   | ratio ≤ 0.50 in **both** phases (P1, P6) **and** canary passed **and** zero seeded-obligation misses | ratio > 0.80 in either phase; **any correctness miss (canary or obligation) is a veto at any speed** |

   Middle band (0.50, 0.80]: no ship, no refute — the measured number replaces the 50–70% planning
   estimate verbatim in FINAL.md and the ship decision returns to the user with that number; it is
   never rounded up to the estimate (@codex-1's rule, adopted). This three-way gate is stricter
   than my round-2 binary (≤0.60 ship / otherwise not) and I prefer it: under it a 0.55 result
   goes to the user with the measurement instead of shipping on my interpolation.
4. **Who runs it, against which source:** the Phase 5 implementer of this idea, on the
   implementation branch, before release, reported as validation evidence in IMPLEMENTATION.md.
   The packet is rendered from the **live resolved protocol**, never a snapshot; the generator
   records `sourceSha256` per run and the full arm reads the same live file. Snapshot sourcing is
   disqualified by @claude-1's PRIMARY: "the drift guard covers the deck and the Go-embedded copy
   but **not** the skill's bundled snapshot, and that all seven installed runtime snapshots were a
   full minor version stale" (@claude-1 round-02, quoted literally; locator given there).

This row is pre-registered: it goes into this idea's FINAL.md (Phase 3) before the experiment runs
(Phase 5), and the result is reported against it unchanged.

## C3 — the §4.0 audit: named follow-up, not in scope

**Decision: a named follow-up idea — `meta-protocol-change-track-gate-audit`.** This idea ships in
one piece; the packet work ships first.

- **What stays in this idea (and always was):** the cells this idea changes — deliberation fix-up
  (5) and deliberation cross-review (3) — land text and code in the same patch. That is the
  round-1 lock plus my round-2 E4 lesson (a printed cap of 2 ran 15 cycles with no recorded
  escalation; PRIMARY locator in my round-2). This idea's FINAL also records the divergences found
  during it — the silent `MaxRounds: 4` default (@hermes-1's open item 2), cross-review defaulted
  to 1 vs "unbounded", the app layer passing only `CrossReviewRounds` (@claude-1 round-2, PRIMARY,
  `internal/app/*.go` grep quoted there) — as the follow-up's seed inventory.
- **Why the general audit is not folded in:** @claude-1's proposal has three parts, and two of
  them are not mechanical. Enumerating every cell is mechanical; but `Timeout per agent`,
  `Reviewers (Phase 6)`, `Review consensus (Phase 7)` have no known-correct per-track values —
  choosing them is the same anchoring work D1 just cost two rounds, times N cells, and the
  structural test needs an enforcement inventory that does not exist yet. Folding that in
  mid-flight delays the packet, which is the measured cost lever and is independent of the audit.
- **What blocks on the follow-up:** nothing in this idea's ship. In the other direction: no
  further §4.0 cell edits land until the audit runs, so the divergence list is published once,
  complete, and the test then keeps it empty. @hermes-1's open items 1 and 2 are absorbed: item 1
  is settled below (D5), item 2 is the follow-up's first row.

## D5 residue — settled at 3, no change

@codex-1's "three cross-review rounds after round 1" and my round-2 "escalate past 3" are the same
boundary. @hermes-1 asked for the anchoring the number needed; the data supplies it: observed
cross-review maximum 3 across the deck (my round-2 E-family) and 2 deliberation-only (@codex-1
round-2 E6, `max_cross_review_rounds=2`). Three keeps one round of headroom over the observed
maximum; same escalate-never-close semantics and recorded finite extension as the fix-up cap.
Implementation contract: deliberation needs the clamp wired in the driver — today
`ReadCrossReviewRounds` honors an explicit 00-prompt value with no ceiling and only `standard`
carries the cap path (my round-2 E5 and concern 3, PRIMARY locators there) — in the same patch as
the fix-up cap, with the boundary test at rounds 3/4. Settled, no change.

## D2, D3 — settled, no change

D2: one standalone generator exposed as `parley protocol packet`; prompt builders call the shared
renderer and never read `COOPERATION.md` themselves. @claude-1's snapshot evidence closes
@hermes-1's open item 4 in favor of the round-1 lock: render from the live resolved protocol,
hash-bound, fail open to full. D3: the §15.1–§15.4 + §15.7 verdict kernel is load-bearing in P5/P8
packets, present before any validation, resolution, or completion claim is authored. Settled, no
change on either.

## C4 — what would stop my signature

1. **A FINAL that sets `MaxFixupCycles: 5` without the inclusive-boundary fix and the 5/6 boundary
   test.** Under the verified `>=` guard that delivers 4 published cycles — a different number than
   ratified. Same for the cross-review boundary at 3/4.
2. **Caps landing text-only.** Both cap cells change in `COOPERATION.md` and the driver in one
   patch; E4 is the demonstrated cost of splitting them.
3. **Any softening of the correctness veto in C2.** No ship on speed alone; the pre-registered row
   lands in FINAL.md before the experiment exists.
4. **Any whole-idea savings claim.** The FINAL reports the per-call ratio as per-call; the
   re-read term (`gatherPriorRounds`/`gatherReviewContext`) is untouched by this change. @claude-1
   adds that regression rank 2 "was implemented and then **deleted** in 1.43.1" (quoted literally)
   — I have not independently verified the deletion and cite it as @claude-1's claim; the scope
   guard stands regardless.

Nothing else. The round-1 locks stand as signing conditions already; I do not relist them.
