---
agent: kimi-1
idea: protocol-mutation-diversity
round: 1
date: 2026-08-31
---

## Summary

The problem is real and I re-measured it rather than inheriting it. Running my own count over
`parley-deck/ideas/` (PRIMARY — `for d in */; … ls -d "$d"round-*` over the 80 idea dirs):
**80 ideas opened round-01; 29 never opened a second round; 44 stopped at exactly two; 7 reached
a third.** That matches `00-prompt.md` exactly on those four numbers. On the fifth it does not:
I count **1** idea reaching a fourth round (`meta-protocol-change-phase-packet-and-fixup-budget`),
not 2. Flagged under Concerns as a candidate correction, not a blocker.

My lens is measurement, and the measurement verdict is: **no mutation operator is warranted yet,
because we cannot measure the thing it would be changing.** The deck has no diversity metric, no
convergence detector, no baseline, and — per `openviking-context-structure` finding 6, verified
PRIMARY at `internal/driver/loop.go:174-175` — no provider input-token telemetry on the headless
path, so even the *cost* side of any operator is currently unmeasurable. Proposing an operator
before the instrument exists is exactly the shape `openviking-context-structure` just closed
against: A4 pre-registered the gate (M1/M2/M3) *before* any behaviour change, and A6 ruled that a
passing gate still does not authorise a default flip. That doctrine is one day old and directly
on point.

On the GA analogy itself I defer to codex-1's skeptic lens, but one measurement-side fact is
decisive on its own: with n=4 participants there are exactly **6 pairwise distances** per round.
Per-idea diversity is statistically unreadable; any measurement that means anything must run at
corpus level over the 80 frozen ideas. This kills every design whose success criterion is "this
idea became more diverse" — the signal only exists in aggregate.

On the QD alternatives from the notes, explicitly:

- **Novelty search / local competition — reject.** It rewards *being different*. We have no
  fitness function (notes §6, SECONDARY) and no outcome label per idea; rewarding difference
  without a quality axis is a nonsense generator with extra steps.
- **MAP-Elites archive — reject for this protocol.** Keeping the best proposal per behavioural
  niche changes what `consensus.md` and `FINAL.md` *are*; the deck's ratification model produces
  one ratified artifact with signoffs, not an archive. That is a protocol rewrite, not a
  diversity mechanism, and §15.6(c)'s ALT- disposition already records rejected alternatives with
  decisive reasons — a one-dimensional shadow of the archive that already ships.
- **FunSearch island model — partially adopt, and mostly already shipped.** Round-01 independence
  *is* the island model (see Existing alternatives). The one genuinely new, stealable piece is
  the **behavioural hash**: collapse trivial surface variants before measuring. That piece is a
  *measurement* tool, and it is what I propose building first.

## Proposed approach

**Leg 1 — Build the convergence metric first, no operator.** A deterministic, replayable,
zero-model-token measurement over existing artifacts, reusing the two extractors that already
ship (locators in Existing alternatives):

- Per round, per participant, extract the set of *decision-relevant units*: named mechanisms,
  rejected alternatives (ALT- ids where present), verdicts, locators cited. Reuse
  `BuildRoundIndex` / `BuildRoundDigest` output as the substrate rather than parsing prose again.
- Compute pairwise similarity over those sets (Jaccard or equivalent), with a FunSearch-style
  collapse step: two claims that are the same claim in different words count once. The collapse
  step is the hard part and I am honest about that in Concerns.
- Emit, per idea, one number: **round-01 convergence distance** (how far apart the independent
  round-01 artifacts actually were) and the idea's round count and outcome class.

Deliverable is a report, not a protocol change. Runs offline over the 80-idea corpus; fully
replayable from git history, satisfying the determinism constraint without an RNG or a seed.

**Leg 2 — Answer the only question that justifies an operator.** Correlate round-01 convergence
distance against outcome: did ideas that converged immediately do *worse* — more rounds burned
refuting nothing, more owner overrides, more post-hoc corrections (the C1–C9 pattern in
openviking) — than ideas that started diverse? Two live data points cut both ways and must both
be explained by the metric: `protocol-generation-bias` converged under facilitator pressure and
its own R5 records that the unanimity was partly a shared prior (bad convergence);
`openviking-context-structure` had four participants independently reject the same six things in
round 1 and was *right on every item* (good convergence — §15.6b was invoked precisely because
unanimity felt suspicious and wasn't). If round-01 convergence distance does not separate these
two classes, **no operator can be trigger-gated and the idea closes as "no, and here is the
measurement"** — which `00-prompt.md` explicitly accepts as an outcome.

**Leg 3 — Only if Leg 2 separates: a trigger-gated, structured, non-random operator.** If and
only if premature convergence is measurable and predictive, the operator is: on a measured
round-01 convergence event, on `standard`/`deliberation` only, one participant is assigned a
*named, specific* divergence — e.g. "argue the strongest case for the highest-ranked rejected
alternative, with its best evidence" — never "be creative", never a temperature bump. Notes §5 is
the constraint here: the literature supports *structured proposal operators with semantic
priors* and measures random mutation as inefficient in large spaces (SECONDARY, per the notes);
a raise-the-temperature proposal must argue against that finding and I do not.

**I state plainly, per the prompt's requirement: Leg 3 is a re-proposal of the steelman clause
removed on 2026-08-29, in a different form.** The differences, and why they might justify return:
the removed clause was unconditional (fired every `standard`/`deliberation` round whether or not
convergence was real) and unmeasured; this form is trigger-gated on a measured convergence event
and instrumented against a frozen corpus baseline. Whether that difference is enough is a
question for the mechanism-design lens, not something I assert.

**Track gating and cost.** Legs 1–2 cost zero participant rounds (offline analysis). Leg 3, if it
ever ships, fires only on a detected event on `standard`/`deliberation`; `fast` never pays. That
satisfies the floor-cost constraint by construction.

## Existing alternatives

Everything below was checked today, 2026-08-31, against the live system. Each entry: what my
proposal would build by hand, what already ships, with a locator, and whether the element is
constraint-forced or merely inherited.

1. **Diversity/divergence extraction over round artifacts — ALREADY SHIPS, TWICE, and is
   unconsumed.** `internal/runner/round_index.go`: `BuildRoundIndex` (line 96) returns a
   per-round L1 view (per-agent status, approx-token counts, H2 section lists, first-sentence
   excerpts), deterministic, zero model tokens; `writeRoundIndex` (lines 83–94) writes it at line
   90, sole call site `runner.go:238`. `internal/driver/digest.go:48` `BuildRoundDigest` is a
   second deterministic extractor, sole call site `driver.go:458`. Per openviking-context-structure
   findings 1–5 (PRIMARY, verified there 2026-08-31 and re-verified by me reading the same files
   today): both run only on the `parley run` path, `gatherPriorRounds` (`runner.go:940`) inlines
   full text and **explicitly skips `_index.md` (line 953)** — generated, then discarded; only
   three `_index.md` files exist in the whole deck, all June 2026. My Leg 1 builds on these two
   extractors rather than writing a third. **Constraint-forced**: zero-new-dependency and
   determinism constraints rule out an embedding model or a new extraction service; the existing
   deterministic extractors are the only permitted substrate.

2. **Roster model-family diversity — ALREADY SHIPS and is the strongest existing diversity
   operator.** `parley roster show` (run today): six active agents, six model families from six
   companies — Anthropic (claude-1), OpenAI (codex-1), Thinking Machines Lab (hermes-1), Moonshot
   AI (kimi-1), xAI (opencode-1), Zhipu AI (zcode-1). This is population-level diversity the GA
   analogy would call the primary mechanism, and it already exists. **However**, effort and speed
   are near-uniform: every active agent runs `deep` speed; effort is `max` for all but hermes-1
   (`high`) and opencode-1 (`xhigh`). There is no configured *process* variance even though there
   is model variance. **Merely inherited** — the roster predates this idea and the prompt's
   non-goals exclude roster change as the primary mechanism; cited here as the existing
   alternative any new operator must beat.

3. **Round-01 independence — the island model ALREADY SHIPS.** The rule that participants must
   not read each other in round 1 (this round's own rules restate it) is functionally FunSearch's
   island model: isolated populations, later recombination at cross-review. **Constraint-forced**
   by the protocol; any proposal must build on it, not duplicate it. What it does *not* cover:
   the islands converge at round 2 and nothing ever re-separates them — FunSearch's *island reset
   on stagnation* has no analogue. That reset is the only part of the island model not already
   shipped, and Leg 3 is precisely a reset operator.

4. **Per-idea `roles:` lens map — ALREADY SHIPS and is a structured divergence operator.**
   `00-prompt.md` frontmatter assigns each participant a distinct advisory lens (this idea:
   protocol fit / skeptic / mechanism design / measurement). This is directed, non-random
   divergence — exactly what notes §5 says the field uses instead of blind mutation.
   **Merely inherited** (per-idea, author-written, advisory only, nothing checks that a lens was
   actually applied), but it is the closest existing thing to a working "mutation" and it is
   *structured*, which is the point.

5. **§15.6(b) correlated-agreement disclosure — ALREADY SHIPS, binds by discipline.**
   `parley-deck/COOPERATION.md:1355-1358`: unanimity among related models is recorded as a shared
   prior, not evidence. Applied live yesterday in openviking A8. **Merely inherited** — §15.6's
   own preamble (line 1349-1350) states only (a) is machine-validated; (b) is prose. It detects
   unanimity *after* consensus, not premature convergence *before* it.

6. **`parley preset list` — NULL RESULT.** Ran it today: *"No roster presets defined."* Presets
   are not an existing diversity mechanism in this deck; there is nothing to compose or compare
   against. Sources consulted: `parley preset list` output, `~/.parley/agents.toml` (no
   `[rosters.*]` block).

7. **Per-agent model/effort config in `~/.parley/agents.toml` — EXISTS as a variance knob, unused
   as one.** Read today: `[defaults] speed = "deep"`; per-agent `reasoning` set to `max`
   (claude, codex, kimi), `high` (hermes), `high` (opencode). A per-idea effort or speed
   randomisation would be a config change, no protocol change needed — and the prompt's non-goals
   explicitly exclude per-agent model/effort config as the *primary* mechanism for exactly that
   reason. **Merely inherited.** Cited as the existing alternative, not proposed.

8. **The removed steelman clause — VERIFIED REMOVED, and addressed head-on.** PRIMARY:
   `git grep` over the history of `parley-deck/COOPERATION.md` finds §15.6(a) *"the strongest
   rejected or unconsidered alternative is steelmanned, with its best supporting evidence…"* in
   commits up to and including `5538796` (parley-deck-cli 1.46.0); the current file
   (`COOPERATION.md:1346-1362`) replaces it with (a) enumerate existing alternatives, (b)
   correlated-agreement disclosure, (c) ALT- disposition — ratified by `protocol-generation-bias`
   FINAL.md (2026-08-29), whose Leg 4 is the deletion. **I am partially re-proposing it**: Leg 3
   is that clause made trigger-gated and measured. The reason it should only return in this
   different form is in its own replacement's record — `protocol-generation-bias` R3 warns that
   unconditional scanner-checkable sections ritualise; a clause that fires only on a measured
   convergence event cannot become per-round ritual.

## Concerns / open questions

1. **The behavioural-hash problem is the whole ballgame and I do not have it solved.** Two
   round-01 artifacts can be verbally distant and decision-identical, or verbally near and
   decision-opposed. FunSearch collapses trivial variants by *behaviour*; our analogue —
   "same claims, same rejections, same verdicts" — needs claim extraction that is either
   deterministic-but-shallow (Jaccard over extracted H2s and named mechanisms) or deep-but-costs
   a model call (which breaks the zero-cost and determinism properties of Leg 1). Open question,
   and it gates Leg 2: a metric that cannot tell verbal diversity from decision diversity will
   report false convergence on openviking-style agreement and false divergence on stylistic
   noise.
2. **The outcome label does not exist.** "Quality did not fall" requires a per-idea quality
   ground truth. Candidate proxies in the record — blocks found in review, corrections logged
   (C1–C9), owner overrides, signoff reservations — are all noisy and partly retrocausal (a bad
   idea that converges quietly looks *cleaner* than a good idea that survived a fight). Until
   this is answered, Leg 2's correlation is unfalsifiable.
3. **A convergence trigger would have fired on openviking round 1 and been wrong.** Four
   independent rejections of the same six items, all correct. This is the strongest evidence that
   "convergence at round 01" cannot itself be the trigger — the trigger must be *premature*
   convergence, which is only separable from correct convergence using the outcome label from
   point 2. The two problems are one problem.
4. **Fourth-round count discrepancy (PRIMARY, mine).** The prompt says 2 of 80 ideas reached a
   fourth round; my count over `parley-deck/ideas/*/round-04` finds exactly 1
   (`meta-protocol-change-phase-packet-and-fixup-budget`). Possible explanations: the prompt
   counted a `review/round-04`, or the aborted-consensus artifact in `protocol-generation-bias`.
   Minor, but the deck's own rule is that a locator is a factual claim and must be checked like
   one (openviking C8); recording the delta rather than silently inheriting either number.
5. **Statistical power.** n=4 participants → 6 pairwise distances per round → per-idea diversity
   is noise. Any measurement must be corpus-level (80 ideas), and any operator evaluation must
   run across many ideas before it can claim an effect. This rules out every "try it on one idea
   and see" evaluation plan.
6. **Cost measurability.** openviking finding 6 (PRIMARY, `internal/driver/loop.go:174-175`):
   headless runners emit no provider input-token telemetry. The "who pays" question for any
   operator cannot currently be answered in tokens — only in wall clock and round count. The
   deferred `parley-context-telemetry` follow-up is a prerequisite for honest operator
   accounting, same as it was for M1 there.

## Risks

- **R1 — Measuring the wrong thing confidently.** A shallow similarity metric produces a
  clean-looking convergence number that does not track decision diversity; Leg 2 then either
  justifies an operator on noise or kills a good one. Mitigation: the metric must reproduce the
  known split (protocol-generation-bias vs openviking-context-structure) before it is believed —
  a two-point sanity check, stated in advance.
- **R2 — Goodhart on the trigger.** If participants know a convergence detector exists, round-01
  artifacts acquire cosmetic divergence to dodge it. This is the same ritualisation failure
  `protocol-generation-bias` R3 pre-committed to deleting. Mitigation: deterministic metric,
  corpus-level evaluation, and the same deletion-default commitment if ritualisation is detected.
- **R3 — The operator's named failure mode is real and expensive.** A mutated participant arguing
  a position it does not hold costs a full participant round (minutes plus a real bill, per the
  prompt's framing) and can leave consensus *worse* — the refutation of an artificial idea
  becomes the round's content. At n=4 and ~2 generations (notes §6), one wasted individual is 25%
  of the population for 50% of the idea's lifetime. This asymmetry is why Leg 3 is gated behind
  Leg 2 rather than proposed alongside it.
- **R4 — Convergence is often correct.** Notes §6 and openviking A8 agree: four agents agreeing
  may mean the answer is right. A diversity mechanism that treats agreement as pathology
  optimises against correctness. The metric must carry the §15.6b disclosure logic with it:
  unanimity is a shared prior to *record*, not a defect to *punish*.
- **R5 — Scope creep into protocol change.** Any operator that changes what `consensus.md` /
  `FINAL.md` are (MAP-Elites archive) or what the runner inlines is a protocol change — new core
  version, attended publish (openviking A6). Legs 1–2 as proposed are read-only and outside the
  protocol entirely; that boundary is deliberate and should be defended in cross-review.
