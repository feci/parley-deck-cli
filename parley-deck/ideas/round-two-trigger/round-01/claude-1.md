---
agent: claude-1
idea: round-two-trigger
round: 1
date: 2026-09-01
---

## Summary

**The premise of my own kickoff is wrong, and the real finding is worse.**

I framed this as "the facilitator decides, and the facilitator is conflicted". Reading the driver
says otherwise: **nobody decides. A counter does.**

`advanceRound` (`internal/driver/driver.go:289-310`, PRIMARY) opens the next round on exactly one
condition — `c.CurrentRound >= 1+d.cfg.CrossReviewRounds`. That budget comes from
`cross_review_rounds` in `00-prompt.md`, **default 1** (`internal/driver/transport.go:34`), clamped
by track policy (`internal/track/track.go`): `fast` forces **0**, `standard` caps at 2,
`deliberation` caps at 3.

**No content is read at any point.** `grep -rni "substantive" --include='*.go' internal cmd`
returns **nothing** (PRIMARY). §15.6(b)'s existing language — *"round 1 closes with no substantive
disagreement"* — has **zero** code acting on it.

So the decision to stop deliberating is made **before the deliberation begins**, as a number in
frontmatter, and the driver counts down to it. Whether the participants actually agreed is never
consulted by anything.

**This reframes the measurement that motivated the idea.** 44 of 80 ideas stopped at *exactly* two
rounds. I read that as convergence. The mechanism above says it is more likely **the default budget
expressing itself**: default `cross_review_rounds: 1` → one independent round plus one cross-review
round → stop. That is a counter reaching zero, not four agents agreeing.

## Proposed approach

**Make the stop condition observable before making it smarter.** Three steps, cheapest first, and
step 1 may end the idea.

1. **Separate "budget exhausted" from "converged" in the record.** Today both produce an idea that
   simply has no `round-03`. The driver already emits a round digest with a `nextAction` string
   (`driver.go:300-307`) — it knows which case it is at the moment it decides. Record that decision:
   which condition fired, what the budget was, and whether it was the default or explicit.
   This is a **CLI/record change, not protocol text**, and it is reversible.
2. **Then measure what the corpus actually contains.** With step 1, "how many ideas closed because
   the counter ran out, versus because participants agreed?" becomes answerable. Today it is not,
   and every number in my kickoff — including the ones I froze — cannot distinguish the two.
3. **Only then decide whether a trigger is warranted**, and what it reads.

**If a trigger is eventually warranted**, my current preference is the cheapest possible form: not
an extra round, but a **single assigned artifact** when the budget would end round 1 on
`deliberation`. But I am explicitly not proposing that yet — step 2 has to run first, and it may
show the budget default is simply mis-set, in which case the fix is a default, not a mechanism.

**What I will not propose**: protocol text in this round. §15.6's preamble already states that only
clause (a) is machine-validated, and the carrier thesis says an unvalidatable prose duty decays to
single-digit compliance. A rule that says "consider opening another round" is exactly that rule.

## Existing alternatives

Mechanisms I would build by hand, and what already ships (all locators verified today, PRIMARY):

**1. A convergence detector.**
- Already ships: **nothing.** §15.6(b) contains the *language* ("round 1 closes with no substantive
  disagreement") but `grep -rni "substantive" --include='*.go' internal cmd` returns no non-test
  hit. The condition is written in the protocol and unimplemented. Constraint-forced: this is the
  actual gap.

**2. A round budget.**
- Already ships and is the **current answer**: `driver.ReadCrossReviewRounds`
  (`internal/driver/transport.go:34`), default 1, per-idea via `cross_review_rounds` frontmatter,
  clamped per track (`internal/track/track.go:182,197,217`). Anyone proposing a trigger must first
  argue why tuning this number is insufficient — it is one frontmatter field and needs no new
  machinery. Inherited, not constraint-forced.

**3. Re-opening a closed deliberation.**
- Already ships: **`parley consensus reopen --reason TEXT`** (`parley --help`, PRIMARY). This is a
  real, existing escape hatch: if consensus closed too early it can be reopened with a recorded
  reason. It is manual and after-the-fact, but it means "closed too early" is **already
  recoverable** and the cost of a wrong stop is lower than my kickoff implied.

**4. Detecting a stalled or degenerate deliberation after the fact.**
- Already ships: **`parley retro scan`** — classifies closed ideas by failure type
  (`blocked-or-abandoned`, `escalation`, `fix-up-heavy`, `review-churn`). Verified running. It does
  **not** have a "converged too fast" class, which is precisely the class this idea is about.
  Extending its taxonomy may be cheaper than a live trigger. Inherited.

**5. Escalating a judgment to the owner.**
- Already ships: the inbox path, `inbox/<from>-to-user_<slug>_<topic>.md` with `blocking: yes`, and
  `ActionEscalated` in the driver (`driver.go:134`, used when consensus triage cannot auto-advance).
  So an escalation channel exists and the driver already uses it for a different condition.

**6. An outer supervisory loop.**
- Already ships: `parley loop tick` and `internal/driver/loop.go` with budgets
  (`loopCostUSD`, `emitLoopBudget`). Note its known limit, carried from the previous idea: runners
  emit no `agent.usage`, so cost enforcement is 0 in practice.

**Sources consulted:** `internal/driver/driver.go`, `internal/driver/transport.go`,
`internal/driver/consensus.go`, `internal/driver/cursor.go`, `internal/track/track.go`,
`internal/runplan/runplan.go`, `parley --help`, `parley retro scan`, `COOPERATION.md` §15.6-15.7.

## Concerns / open questions

1. **My kickoff's framing was wrong and the corpus numbers inherit that.** I described a conflicted
   human judgment; the mechanism is a counter. Every number I froze at `2d17478` is still correct
   as a count, but its *interpretation* — "premature convergence" — is not established. It may be
   "the default budget is 1". Someone should challenge whether this idea has a problem at all.
2. **How many of the 28 are `fast`-policy zero-budget ideas?** `fast` forces
   `CrossReviewRounds: 0`, so a `fast` idea closing after one round is the system working exactly
   as designed. Only 1 of the 28 carries `track: fast` — but 19 carry **no track at all** and
   predate the field, so their effective policy is the `-1`/default path. That needs checking
   before anyone calls those 19 premature.
3. **Is `consensus reopen` the answer already?** If closing early is cheaply reversible, the value
   of a pre-emptive trigger drops a lot. I do not know how often it has been used — the record
   should be checked rather than assumed.
4. **Whose cost?** My kickoff asserted the facilitator bears the cost and is therefore conflicted.
   Under auto-drive the facilitator does not decide at all. The conflict I described applies only
   to *manual* facilitation — which, per the previous idea, is almost all real runs (only 3
   `_index.md` files exist, and those are written only on the `parley run` path).

## Risks

- **Building a detector for a problem that is a mis-set default.** The cheapest outcome of this
  idea is "raise `cross_review_rounds` for `deliberation` and stop". If step 2 shows that, a
  mechanism would be pure cost.
- **A rule nothing can check.** Any prose duty of the form "consider whether to continue" lands in
  §15.6's unvalidated half and decays. If a trigger ships, it must be a gate with a call site, not
  a sentence — the defect class this deck has hit repeatedly.
- **Ceremony at the round boundary.** A detector that fires often and changes nothing taxes every
  deliberation. It must be measured on whether firing changed a `FINAL.md`.
- **Over-deliberation is also a failure.** 52 of 80 ideas did open a second round and cross-review
  moves positions in most artifacts. Forcing more rounds on top of a system that is mostly working
  risks paying for rounds that ratify what round 2 already settled.
- **My own conflict, stated plainly:** I proposed this idea, wrote its kickoff, froze its corpus and
  am now its first participant. The framing error in concern 1 is mine, and I found it only by
  reading code I should have read before writing the kickoff. Weight my round-01 accordingly.
