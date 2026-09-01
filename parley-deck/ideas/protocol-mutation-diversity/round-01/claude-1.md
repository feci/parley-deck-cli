---
agent: claude-1
idea: protocol-mutation-diversity
round: 1
date: 2026-08-31
---

## Summary

**Yes, something should be added. No, it should not be randomness.** And the most uncomfortable
finding first: **we deleted this mechanism two days ago.**

The superseded §15.6 (core 2.10.0, replaced on 2026-08-29 by `protocol-generation-bias`) contained,
verbatim, the mechanism the owner is now reaching for — including its trigger and its
deliberation-track form:

> "if round 1 closes with **no substantive disagreement** ... consensus MUST NOT close until (a)
> the strongest rejected or unconsidered alternative is **steelmanned** ... On `deliberation`, one
> participant is **assigned** and files it as a canonical round artifact."

So the protocol already had a premature-convergence detector (*no substantive disagreement in round
1*) wired to a divergence duty, and we removed it. Any proposal in this idea must start there, not
from population genetics.

What GA theory adds on top is one real correction to the owner's framing, and it comes from the
literature rather than from me: for LLMs, mutation is **a learned proposal operator with semantic
priors**, not blind perturbation, and "random mutations [are] inefficient in large solution spaces"
(notes §5). Raising temperature or shuffling prompts is the version of this idea that does not
work. A *named, drawn, structured* divergence is the version that does.

## Proposed approach

**A divergence assignment with a recorded draw.** One artifact, `deliberation` track only.

1. **Trigger — reuse the one we already ratified.** Round 1 closes with no substantive
   disagreement. This is the premature-convergence detector, it is already written, and it was
   already agreed once.
2. **The draw is where randomness lives, and it is replayable.** One participant is selected, and
   one **divergence axis** is drawn, by a deterministic PRNG seeded from `idea-slug + round`. The
   seed is written into the artifact frontmatter. This satisfies the determinism constraint: the
   run is reproducible from the record, and nobody chooses who gets the awkward job.
3. **The operator is a fixed, enumerated axis list — this is the semantic prior.** The assignee
   must produce a position that differs along the *drawn* axis, e.g.:
   - invert the load-bearing assumption every participant shared;
   - treat as variable the constraint everyone treated as fixed;
   - solve it at 100× the assumed scale;
   - solve it with zero new code;
   - argue the case of the stakeholder nobody was speaking for.
   A drawn axis forces a *specific kind* of difference. "Be creative" does not, and is the form the
   literature says fails.
4. **Local competition, not global.** The divergent artifact is **not** judged against the
   consensus on correctness. It is judged on one question: *did it name something the consensus had
   not considered?* This is novelty-search-with-local-competition (notes §2) expressed as a
   protocol rule, and it is what stops the assignee being punished for losing an argument it was
   ordered to have.
5. **A null result is a finding, not a failure** — the removed clause already said this, and it is
   the release valve that keeps the mechanism honest: "this axis yields nothing here, and here is
   why" is compliant.
6. **Cost: one artifact section, no extra round.** Not per-round, not on `fast`, not on `standard`.

I am **not** proposing a population, a fitness function, crossover, or a generational loop. Those
require the things notes §6 says we do not have.

## Existing alternatives

The mechanisms I would build by hand, and what already ships:

**1. A divergence duty triggered by round-1 unanimity.**
- Already shipped, then **removed**: superseded §15.6 in `~/.parley/protocol/core/2.10.0/`
  (verified verbatim today, PRIMARY). It had the trigger, the deliberation-track assignment, the
  `## Adversarial alternative` section form for `standard`, and the null-result release valve.
  **Constraint-forced: no.** This is the single strongest existing alternative and my proposal is
  substantially a restoration of it with a drawn axis added. Anyone who wants to reject my proposal
  should argue that the 2026-08-29 removal was right.

**2. An archive of distinct approaches rather than one winner (MAP-Elites-shaped).**
- Already ships: **§15.6(c)**, `## Alternatives disposition` in `consensus.md` — "an `ALT-` id and
  an adopt or reject with the decisive reason for each alternative" (`COOPERATION.md:1359`,
  PRIMARY). This is a proto-archive: it keeps rejected alternatives addressable by id.
  What it does not do: it records dispositions of alternatives that *were raised*; it creates no
  pressure to raise a distinct one. Inherited, not constraint-forced.

**3. Population diversity via participant selection.**
- Already ships: the roster spans five model families and five companies (`parley roster show`,
  PRIMARY). This is the deck's real diversity operator and it is doing most of the work today.
- Already ships: `parley preset list` — **but it returns "No roster presets defined"** (PRIMARY,
  run today). The mechanism for swapping in a deliberately different cohort exists and has zero
  configured instances.

**4. Per-individual variance via config.**
- Already ships: `~/.parley/agents.toml` per-agent `model` / `reasoning` / `speed`. Today every
  active agent is `deep` speed and `max` effort except hermes-1 at `high` — **no variance is
  configured at all**. Explicitly named a non-goal as a *primary* mechanism in `00-prompt.md`, and
  I agree it is not one, but it is the cheapest existing knob and must be on the record.

**5. Isolation between candidates.**
- Already ships: round-01 independence is an island model in everything but name — participants
  must not read each other. FunSearch's addition over this is *periodic island reset on stagnation*
  (notes §3), which we have no analogue for. Constraint-forced: yes, this is protocol.

**6. Disclosure of correlated agreement.**
- Already ships: **§15.6(b)** — unanimity among related models is recorded as a shared prior.
  It **describes** the problem and takes no action on it. That gap is precisely this idea.

**7. Asking a specific agent a targeted question.**
- Already ships: `parley consult [--timeout D] AGENT [QUESTION]` (`parley --help`, PRIMARY).
  Advisory, does not produce a canonical artifact, and is facilitator-initiated — so it cannot
  serve as an automatic response to detected convergence.

**Sources consulted:** `~/.parley/protocol/core/2.10.0/COOPERATION.md` (superseded §15.6),
`parley-deck/COOPERATION.md` §15.6, `parley roster show`, `parley preset list`, `parley --help`,
`~/.parley/agents.toml`, and `reference/ga-notes.md`.

## Concerns / open questions

1. **Is the 2026-08-29 removal being quietly reversed?** I want this argued explicitly. That idea
   replaced a *judgment* duty (steelman) with a *checkable* duty (enumerate what already ships),
   and the checkable one is machine-validated while the steelman never was. Restoring an
   unvalidatable duty may re-create the exact decay §15.6's preamble warns about.
2. **Does the trigger even fire?** "No substantive disagreement in round 1" was never
   machine-detected — it is a judgment call by the facilitator, who is also a participant and has
   an interest in not opening extra work. Measured: 29 of 80 ideas never opened round 2 at all. I
   cannot tell from the record how many of those *should* have fired this trigger.
3. **n=4 and 2 generations.** Notes §6 is not a footnote. Every QD result assumes populations and
   generation counts we will never have. What survives the analogy is the *framing* (reward
   difference, keep an archive, detect stagnation), not any algorithm.
4. **Who draws?** If the facilitator runs the PRNG, the facilitator can re-roll. The seed must be
   derived from fixed inputs, not chosen — and the derivation must be in the record.

## Risks

- **Manufactured dissent.** The assignee argues a position it does not hold, the round spends
  itself refuting an artificial idea, and the FINAL is worse. This is the central risk. Mitigations:
  local competition (item 4), the null-result valve (item 5), and the axis list containing only
  genuine reframings rather than contrarian poses.
- **Ceremony.** The duty is satisfied with a paragraph nobody reads and nothing changes. The
  measurement must therefore be *did the divergent artifact change `FINAL.md`* — and if it never
  does across N ideas, the rule should be deleted, not defended.
- **Re-creating a rule nothing can check.** §15.6's own preamble says only (a) is machine-validated
  and says so explicitly rather than implying a gate that does not exist. A divergence duty is
  judgment-shaped; it will sit in the unvalidated half. That is an argument for keeping it small.
- **Cost on the wrong track.** If this ever fires on `fast` or `standard`, it taxes the 91% of
  ideas that finish in two rounds to benefit the 9% that do not.
- **The owner asked for randomness and I am proposing structure.** If the drawn axis list is the
  wrong shape, this degenerates into a checklist and loses exactly the "completely different angle"
  the request was about. Worth a counter-proposal from someone who thinks genuine stochasticity
  beats a curated list.
