---
idea: protocol-mutation-diversity
author: claude-1
created: 2026-08-31
purpose: Copied external material so participants without web access work from the same facts.
---

# Genetic algorithms, mutation and diversity — source notes

**Why this file exists.** §11 requires copying external snippets when participants may lack access.
Not every participant can browse. Work from this file; if you verify something independently, say
so and tag it PRIMARY.

Tags: **PRIMARY** = quoted from the paper/source itself. **SECONDARY** = third-party summary.
**RECALL** = unverified, treat as a claim to check.

---

## 1. The classical picture (SECONDARY — textbook consensus, not re-derived here)

In a genetic algorithm, **mutation** applies a small random perturbation to an individual
(bit-flip, Gaussian noise on a real vector) at some rate. Its job is **not** to improve a solution.
Its job is to keep the population from collapsing onto one point — to preserve the *raw material*
that crossover and selection can exploit. Selection is the exploiter; mutation is the explorer.

A common rule of thumb is a per-gene rate of about `1/L` for a genome of length `L`, so roughly one
change per offspring. Rates can be fixed, scheduled, or adapted from observed diversity.

**The failure it exists to prevent is `premature convergence`:** the population loses diversity
early, every individual is a near-copy, and the search stalls on a local optimum it can no longer
escape. This is the term to keep in mind — it is the exact shape of our problem.

## 2. Quality-Diversity: the part that has actually moved on (SECONDARY — search summaries)

Modern work does **not** treat "add more randomness" as the answer. Two ideas dominate:

**Novelty Search with Local Competition** — "explicitly treats exploration and exploitation as
different objectives as a solution to premature convergence." You do not reward being *better*;
you reward being *different*, and compete only locally among similar solutions.

**MAP-Elites** — "discretizes the space of possible behavioral descriptors into a grid (also called
archive), with the goal of filling each cell of the grid with the highest performing individuals."
You keep **the best solution per behavioural niche**, not the single best overall. The archive
"guarantees ... as many behaviour niches as are accessible by the search, with the best solutions
found per niche", and the diverse collection acts as **stepping stones** that let search escape
local optima.

The framing worth stealing: *the output of the process is an archive of diverse good solutions,
not one winner.*

## 3. FunSearch — the island model (SECONDARY, DeepMind's LLM-driven search)

"Unlike traditional genetic algorithms anchored to manually engineered mutation and crossover
operators, FunSearch leverages modern LLMs as universal code generators, orchestrated within an
**island-based evolutionary architecture**." Two mechanisms are directly relevant to us:

- **behavioural hashes to eliminate trivial variants** — two candidates that behave identically are
  collapsed, regardless of surface wording;
- **islands are periodically reset to prevent stagnation** — a sub-population that has converged is
  wiped and re-seeded from elsewhere.

## 4. Promptbreeder (PRIMARY-ish — arXiv:2309.16797, title and mechanism from the paper listing)

"PROMPTBREEDER: SELF-REFERENTIAL SELF-IMPROVEMENT VIA PROMPT EVOLUTION". It evolves *task-prompts*
and *mutation-prompts*, and describes **five classes of mutation operator**, including direct
mutation, estimation-of-distribution mutation, **hyper-mutation** (mutating the mutation prompt
itself) and **Lamarckian mutation**. It is explicitly framed as "a diversity maintaining
evolutionary algorithm ... to [solve] the problem of diminishing returns".

## 5. THE CENTRAL CAVEAT — read this before proposing anything (SECONDARY, survey material)

For LLMs, the literature does **not** support blind randomness:

> "Mutation and recombination are better viewed as **learned proposal operators** than as blind
> random perturbations, distinguishing them from traditional genetic algorithm operators."

and

> "Mutation operators can use **semantic priors** to effectively guide exploration ... addressing
> the **inefficiency of random mutations** in large solution spaces."

A bit-flip is cheap and a wasted LLM round is not. In a space this large, uniform random
perturbation is close to worthless; a *structured* operator that forces a specific, named kind of
difference is what the field actually uses. **Any proposal here that amounts to "raise the
temperature" should expect to be measured against this.**

## 6. What the analogy costs — where the mapping is weak (claude-1's own analysis, SECONDARY)

State these honestly rather than discovering them later:

- **No fitness function.** GAs need a cheap, automatic, numeric fitness. Parley Deck's "fitness" is
  an agent's or the owner's judgment, evaluated at most once per round. Selection pressure without
  measurable fitness is just preference.
- **Population size ~4, generations ~2-3.** GA intuitions assume populations of 10^2-10^6 over
  10^2+ generations. At n=4 and 2 generations, most GA machinery has no room to operate, and
  variance dominates.
- **Cost per individual is enormous.** One "individual" is one agent round: minutes of wall clock
  and a real bill. A GA evaluates millions of cheap individuals.
- **No crossover.** There is no defined way to splice half of codex's design onto half of kimi's
  and get something coherent.
- **Convergence is often correct.** A GA population converging on one answer is a failure mode. Four
  agents agreeing may simply mean the answer is right. Diversity is not free: forcing a participant
  away from the correct answer costs a round and can inject a wrong idea that must then be refuted.

## 7. What Parley Deck already does that IS an evolutionary mechanism (claude-1, PRIMARY-local)

Do not propose these as new. Verify them and say what they do not cover:

- **Round-01 independence is an island model.** Participants must not read each other in round 1.
- **The roster is a diversity operator.** Five model families / five companies (Anthropic, OpenAI,
  Thinking Machines Lab, Moonshot, Zhipu/xAI) — `parley roster show`.
- **`roles:` in `00-prompt.md`** — per-idea advisory lenses.
- **§15.6(b)** — unanimity among related models must be recorded as a shared prior, not evidence.
- **§15.6(a)** — the enumerate-existing-alternatives duty, machine-validated.
- **A steelman duty existed and was REMOVED on 2026-08-29.** The superseded §15.6 required, on
  `standard` and `deliberation`, that "the strongest rejected or unconsidered alternative is
  steelmanned, with its best supporting evidence and an observation that would change the
  recommendation", as an assigned artifact or an `## Adversarial alternative` section. The
  `protocol-generation-bias` idea replaced it. Anyone proposing a diversity mechanism must say
  whether they are re-introducing that clause, and if so why it should return in a different form.

## Sources

- https://arxiv.org/pdf/2309.16797 (Promptbreeder)
- https://arxiv.org/pdf/2007.05352 (Multi-Emitter MAP-Elites)
- https://arxiv.org/html/2505.15741v1 (Evolutionary Computation and LLMs: a survey)
- https://arxiv.org/pdf/2401.10510 (When LLMs Meet Evolutionary Algorithms)
- https://www.emergentmind.com/topics/funsearch-algorithm
- https://www.emergentmind.com/topics/map-elites-algorithm
