---
idea: protocol-mutation-diversity
author: user
created: 2026-08-31
participants: [claude-1, codex-1, hermes-1, kimi-1]
roles:
  claude-1: protocol fit — where would this live, and what does it cost the protocol
  codex-1: skeptic — does the GA analogy survive n=4 and 2 generations?
  hermes-1: mechanism design — the concrete operator, if one is warranted
  kimi-1: measurement — how would we know diversity actually increased
status: final
track: deliberation
---

## Problem / idea

The owner's proposal, in his words (translated from Slovak):

> "I've got another idea — to bring a bit of **randomness** into the ideation process so that some
> participant can come up with a **completely different angle**, and maybe a unique solution. What
> about integrating properties of **genetic algorithms** into the protocol, specifically
> **mutation** — that's what injects randomness in GAs. Study genetic algorithms and mutation and
> try to work out whether it could bring the protocol something positive."

The question is genuinely open: **could a mutation-like operator improve Parley Deck, and if so,
what exactly is the operator?** "No, and here is the measurement that says so" is a fully
acceptable answer. So is "yes, but not as randomness."

External material is copied into `reference/ga-notes.md` with provenance tags. **Read it first.**
Several participants cannot browse. Note especially §5 of that file: the literature says that for
LLMs, mutation is better understood as a *learned proposal operator with semantic priors* than as
blind random perturbation, and that random mutation is **inefficient** in large spaces. A proposal
that reduces to "raise the temperature" must argue against that finding, not ignore it.

### The problem is real and measured (PRIMARY, this deck, 2026-08-31)

Premature convergence is not hypothetical here:

- Of **80** ideas that opened `round-01`, **29 never opened a cross-review round at all**.
- **44** more stopped at exactly two rounds.
- Only **7 of 80** ever reached a third round; **2** reached a fourth.

So ~91% of all deliberation in this deck is finished within two rounds.

Two live examples from the last 48 hours, both closed:

- `protocol-generation-bias` — the idea that produced §15.6 — was itself about correlated
  agreement, hidden profile and design fixation.
- `openviking-context-structure` closed today. All four participants **independently** rejected the
  same things in round 1 and converged on the same mechanism. The consensus had to record its own
  unanimity as a *shared prior, not evidence* (§15.6b). Diversity did not come from the protocol
  there; it came from one participant reading the runner source instead of the protocol.

Also relevant: **effort and speed are near-uniform across the roster** (`parley roster show` —
every active agent runs `deep` speed, and effort is `max` for all but hermes-1 at `high`). There is
no configured variance of any kind today.

### What to actually decide

For any mechanism you propose:

- **The operator.** What exactly is perturbed — the prompt, the role/lens, the participant's
  assigned position, the model config, the reading order? Name it. "Add randomness" is not an
  operator.
- **When it fires.** Every round? Only on detected convergence? Only above a track threshold?
- **Who pays.** An extra participant round costs minutes and money. State the cost.
- **The measurement.** How would we know diversity rose *and quality did not fall*? Diversity is
  trivially maximised by producing nonsense.
- **The failure mode.** Specifically: a mutated participant argues a position it does not hold,
  the round spends itself refuting an artificial idea, and consensus is *worse* than without it.

Consider, and say explicitly which you are adopting or rejecting, the QD alternatives in the notes
— they are what the field moved to *instead of* raw mutation:

- **novelty search / local competition** — reward being different, compete only locally;
- **MAP-Elites archive** — keep the best proposal *per behavioural niche* rather than one winner
  (this would change what `consensus.md` and `FINAL.md` even are);
- **FunSearch island model** — behavioural hashes to collapse trivial variants, and periodic island
  reset on stagnation.

## Constraints

- **A protocol change is a new core version**, published attended by the owner. Do not propose
  editing `COOPERATION.md` directly. A CLI-only change that quietly alters deliberation semantics
  is worse than a protocol change, not better — see `openviking-context-structure` A6.
- **Zero new runtime dependencies.** Single Go binary + Node installer.
- **Do not raise the floor cost of a `fast` two-round idea.** Any mechanism must be track-aware or
  trigger-gated.
- **Determinism where it matters.** Artifacts must stay reproducible and auditable; a mechanism
  whose behaviour cannot be replayed from the record is not acceptable. If you propose an RNG, say
  where the seed is recorded.
- **§15 applies:** tag claims PRIMARY / SECONDARY / RECALL; untagged is RECALL. Do not state a
  result about a paper the notes file does not support.
- **§15.6(a) applies to your round-01 file** — see the required shape. Note in particular the
  removed steelman clause (notes §7): if you are re-proposing it, say so.
- English only for everything under `parley-deck/`.

## Non-goals

- Building a real genetic algorithm, a population manager, or a fitness function.
- Crossover between participants' artifacts. There is no defined splice of half a design onto
  another and no reason to believe one exists.
- Changing the roster, adding participants, or per-agent model/effort defaults as the *primary*
  mechanism (that is a config change and is already available today — but it may be cited as an
  existing alternative).
- Automating the owner out of any ratification step.
