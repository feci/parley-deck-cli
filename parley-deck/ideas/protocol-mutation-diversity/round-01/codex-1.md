---
agent: codex-1
idea: protocol-mutation-diversity
round: 1
date: 2026-08-31
---

## Summary

[SECONDARY — `reference/ga-notes.md` §§5–6] The genetic-algorithm analogy does not survive the load-bearing details: Parley has about four expensive proposals, one or two discussion rounds rather than generations of offspring, no cheap fitness function, and no crossover whose raw material mutation would preserve. [PRIMARY — repository and commands cited below] Parley already separates round-one proposals, selects different model families, assigns advisory lenses, enumerates shipped alternatives, and discounts correlated unanimity. I therefore recommend **no production mutation operator and no core-version change now**. If the owner wants evidence rather than a permanent null, the least misleading next step is a sealed, benchmark-only test of an endorsable semantic reframe carried by the existing `roles:` field—not random temperature, forced advocacy, or another live round.

## Proposed approach

### Production decision: a scoped null

Do not add mutation semantics to `COOPERATION.md`, the runner, `consensus.md`, or `FINAL.md`. The reported round-count distribution in `00-prompt.md` shows early convergence, but it does not distinguish premature convergence from correct convergence. Without that distinction, diversity is not an objective: a mechanism can improve a novelty count merely by producing an answer the rest of the round must reject.

[INFERENCE from `reference/ga-notes.md` §§1, 5–6 and `COOPERATION.md` Phase 1–3] The mapping fails in five places:

| GA element | Parley Deck analogue | Why it is not load-bearing here |
| --- | --- | --- |
| Population | Four named participants in this idea | Four samples cannot maintain or estimate a useful population distribution. |
| Generation | Round 1 followed by cross-review | Round 2 critiques persistent participants and artifacts; it does not create a fresh population of offspring. |
| Fitness | Judgment, signoff, and eventual implementation evidence | There is no cheap, automatic scalar score with which to select thousands of variants. |
| Mutation plus crossover | A changed prompt or lens, with crossover explicitly out of scope | Classical mutation preserves material for later recombination; Parley has no defined recombination operator. |
| Premature convergence | Early agreement | Agreement is sometimes the correct result, and the protocol has no independent label that separates correctness from collapse. |

[SECONDARY — `reference/ga-notes.md` §5] If any treatment is tested, it should be a learned proposal operator with a semantic prior, not blind perturbation. That means the model searches inside a named transformation and may return a null; it is not ordered to defend a position it does not hold.

### Optional evidence-only pilot

The pilot is not a protocol change and does not authorize later adoption.

1. **Operator — endorsable semantic reframe.** In the treatment arm, add one advisory lens to one participant prompt: identify one load-bearing assumption; apply exactly one named transformation—`boundary`, `mechanism`, `representation`, or `objective`; emit an alternative only if the participant would actually recommend it; otherwise emit a scoped null and the observation that would change that result. The perturbed object is the **role/lens prompt**, not the model, effort, reading order, temperature, or assigned position.
2. **When it fires.** Benchmark only, on `deliberation` prompts reconstructed from their original kickoff material. Each call receives `00-prompt.md` and its kickoff reference bundle, not later rounds or `FINAL.md`. It never fires on `fast`, never opens an extra canonical round, and never changes a live idea. If a later proposal seeks production adoption, the only cost-compatible shape is one preselected round-one participant whose existing call is modified rather than an added call.
3. **Selection and audit.** Use a deterministic selector over `(idea slug, participant id, core version)` to choose the participant and transform. Record that tuple, algorithm version, prompt, model, effort, output, and wall time with the benchmark. The repository already uses Go's `crypto/sha256` in, for example, `internal/loop/loop.go:12,113`; no dependency or hidden RNG is needed. The tuple is the recorded seed. Audit replay consumes the recorded output rather than pretending that a second hosted-model call is bit-reproducible; the assignment, input, and evidence record must be reconstructable exactly.
4. **Who pays.** A first batch of 12 paired prompts costs 24 deep calls—one baseline and one treatment call per prompt—with no live-round latency. A passing first batch must be repeated on 12 new prompts before a core-change idea is opened, for a maximum of 48 deep calls before protocolization. This cost belongs to the optional benchmark, not to every idea.
5. **Diversity measurement.** Blind to arm, reduce each valid proposal to a behavioural tuple: `(problem boundary, changed assumption, mechanism, observable test)`. Collapse surface variants with the same tuple. The primary diversity result is the paired change in unique **valid** tuples, not lexical distance or raw idea count.
6. **Quality measurement.** A blind evaluator scores constraint compliance, factual/provenance discipline, feasibility, decision usefulness, and whether the proposed observable test could falsify the mechanism. Separately count material unsupported claims, violated constraints, and alternatives that exist only because the prompt forced advocacy.
7. **Gate.** The treatment advances only if both independent batches have a median paired gain of at least one valid behavioural tuple, no lower median blind quality score, and no increase in material unsupported claims or constraint violations. Failure or ambiguity preserves the null result. Passing opens a new protocol-change idea; it does not itself change the core.

Load-bearing classifications:

- **Constraint-forced:** no `fast` cost; null-allowed endorsement; deterministic and recorded assignment; behavioural rather than lexical deduplication; blind quality measurement; and a separate ratification idea after a passing replication.
- **Merely inherited:** the `roles:` carrier, the four transformation labels, the 12-prompt batch size, SHA-256 as the convenient deterministic selector, and the existing model roster. These are replaceable implementation choices, not reasons to accept the operator.

### Relation to the removed steelman clause

I am **not re-proposing the removed steelman duty**. [PRIMARY — executed check] `git show 59eb663d0c98cdc6f524498a70ce765d7bde6992 -- parley-deck/COOPERATION.md` shows that the 2026-08-29 implementation deleted the requirement that “the strongest rejected or unconsidered alternative is steelmanned” and replaced it with enumerated existing alternatives plus correlated-agreement and disposition duties. `protocol-generation-bias/FINAL.md:68–77` ratifies that deletion. The experimental reframe is adjacent to the deleted design family, so it must remain labelled as a treatment; its null permission and endorsement requirement are specifically intended to avoid silently restoring forced devil's advocacy.

## Existing alternatives

The production recommendation is a scoped null. These are the components a mutation proposal would otherwise build by hand, the closest mechanisms that already ship, and their disposition.

1. **ALT-01 — Model-family allocator (`merely inherited`).**
   - Hand-built component: a population-diversity scheduler.
   - Already ships: `parley roster show`. [PRIMARY — command executed 2026-08-31] Relevant output:

     ```text
     $ parley roster show
     claude-1   claude    active  yes  claude-opus-5[1m]     Claude Opus  Anthropic              max    deep  yes  inherited-roster
     codex-1    codex     active  yes  gpt-5.6-sol           GPT          OpenAI                 max    deep  yes  inherited-roster
     hermes-1   hermes    active  yes  fireworks/inkling     Inkling      Thinking Machines Lab  high   deep  yes  inherited-roster
     kimi-1     kimi      active  yes  kimi-code/k3          Kimi K       Moonshot AI            max    deep  yes  inherited-roster,effort-from-config
     opencode-1 opencode  active  yes  litellm/xai/grok-4.6  Grok         xAI                    xhigh  deep  yes  inherited-roster,effort-from-config
     zcode-1    zcode     active  yes  zai/glm-5.3           GLM          Zhipu AI               max    deep  yes  inherited-roster,model-from-config,effort-from-config
     ```

     The current idea selects four of these families in `00-prompt.md:4`.
   - **PRIMARY — WRONG for the live roster:** `reference/ga-notes.md` §7 says five families/five companies; the live command shows six active families/companies. This count correction does not prove that their proposals are independent.
   - Disposition: use the existing roster; do not mutate membership or model family as the primary mechanism, which is also excluded by the kickoff non-goals.

2. **ALT-02 — Island isolation (`merely inherited`).**
   - Hand-built component: separated sub-populations and a pre-contact generation.
   - Already ships: `COOPERATION.md:331` requires round 1 to be written without reading other round-one files. The current idea follows that rule.
   - Disposition: retain. It is an island-like isolation boundary, but not a FunSearch island model: there are no evolving island populations, behavioural selection, migration, or periodic reset.

3. **ALT-03 — Semantic lens injection (`merely inherited`).**
   - Hand-built component: prompt mutation and assignment to a participant.
   - Already ships: the per-idea `roles:` map (`COOPERATION.md:288–311`). This kickoff already assigns protocol-fit, sceptic, mechanism-design, and measurement lenses (`00-prompt.md:6–10`).
   - Disposition: this is the correct carrier for a benchmark treatment. Do not add a second mutation field until this carrier is shown inadequate.

4. **ALT-04 — Correlated-agreement detector (`constraint-forced`).**
   - Hand-built component: detection and disclosure of nominally independent convergence.
   - Already ships: §15.6(b) at `COOPERATION.md:1355–1357`. `openviking-context-structure/FINAL.md:110–112` applies it by recording four-model unanimity as a shared prior and naming M2 as the disconfirming test.
   - Disposition: retain. It prevents agreement from masquerading as independent evidence, but it does not generate another proposal.

5. **ALT-05 — Alternative acquisition and archive (`constraint-forced`).**
   - Hand-built components: search for already-shipped mechanisms, stable proposal identifiers, and a terminal disposition.
   - Already ships: §15.6(a) and (c) (`COOPERATION.md:1352–1361`), implemented from `protocol-generation-bias/FINAL.md:16–63`. Round 1 must enumerate shipped alternatives; consensus assigns each an `ALT-` identity and adopts or rejects it.
   - Disposition: retain and measure before adding another generator. The `ALT-` record is an archive of considered proposals, but not a MAP-Elites archive: it has no behavioural cells or best-per-cell fitness.

6. **ALT-06 — Named preset bank (`merely inherited`).**
   - Hand-built component: reusable mutation/treatment profiles.
   - Already ships: the CLI surface exists, but [PRIMARY — command executed] `parley preset list` returned: `No roster presets defined. Add a [rosters.<name>] block to ~/.parley/agents.toml or parley-deck/agents.toml.`
   - Disposition: no semantic-mutation preset exists. Roster presets are the wrong layer for changing an idea's reasoning lens; do not create one before the pilot shows a need.

7. **ALT-07 — Per-agent model/effort variation (`merely inherited`).**
   - Hand-built component: mutation of model, effort, speed, or sampling configuration.
   - Already ships: `~/.parley/agents.toml:15,26–160` configures global `speed = "deep"`, per-agent models, and per-agent reasoning; `parley roster show` reports the effective model/effort/speed. In the live roster every active agent uses `deep`; efforts are mostly `max`, with Hermes `high` and OpenCode `xhigh` in the effective output.
   - Disposition: this variance is configurable today and can be tested without a core change. Reject it as the primary operator: changing capability and sampling at once confounds “different frame” with “different model quality,” and the kickoff explicitly makes it a non-goal.

8. **ALT-08 — Forced steelman/adversarial alternative (`merely inherited historical alternative`).**
   - Hand-built component: assignment of a participant to construct the strongest rejected position.
   - Already shipped and was deliberately removed on 2026-08-29. The deletion is visible in commit `59eb663d0c98cdc6f524498a70ce765d7bde6992`; the ratified rationale is `protocol-generation-bias/FINAL.md:68–77`.
   - Disposition: **reject for production; not re-proposed.** A benchmark treatment must allow a null and must not require a participant to argue a position it does not endorse.

9. **ALT-09 — Novelty search with local competition (`constraint-forced comparison; not adopted`).**
   - Hand-built components: behavioural distance, neighbourhoods, and a local quality score.
   - Closest shipped mechanisms: advisory `roles:` provide named niches, and §15.6(c) preserves alternatives, but there is no local competition or fitness signal. [SECONDARY — `reference/ga-notes.md` §2]
   - Disposition: reject. Without an independently measurable local quality score, “competition” is another agent's preference, not selection.

10. **ALT-10 — MAP-Elites archive (`constraint-forced comparison; not adopted`).**
    - Hand-built components: behavioural descriptors, cells, per-cell fitness, archive selection, and a multi-output terminal artifact.
    - Closest shipped mechanism: `## Alternatives disposition` records multiple options, while `FINAL.md` deliberately remains one authoritative plan. [PRIMARY — `COOPERATION.md:391–414,1359–1361`; SECONDARY for MAP-Elites mechanics — `reference/ga-notes.md` §2]
    - Disposition: reject. It would redefine consensus and finalization, and Parley lacks the fitness required to identify the best proposal per niche.

11. **ALT-11 — FunSearch-style reset and behavioural hash (`behavioural hash: constraint-forced for the pilot`; `reset/reseed: merely inherited comparison, not adopted`).**
    - Hand-built components: stagnation detector, behavioural hash, island reset, and reseeding policy.
    - Closest shipped mechanisms: round-one isolation supplies only the island boundary; `ALT-` identifiers supply identity but not behavioural equivalence. [SECONDARY for FunSearch mechanics — `reference/ga-notes.md` §3]
    - Disposition: reject resets and reseeding because they add calls after convergence without a validated trigger. Adopt only the behavioural-tuple idea as a **constraint-forced benchmark measurement**, not as a protocol archive.

12. **ALT-12 — Blind random or temperature mutation (`merely inherited from the initial proposal`).**
    - Hand-built components: RNG policy, rate, seed record, and replay semantics.
    - Closest shipped mechanism: none is needed; explicit role lenses already provide auditable semantic variation. [SECONDARY — `reference/ga-notes.md` §5 says random mutation is inefficient in large solution spaces and favours learned semantic proposal operators.]
    - Disposition: reject. The pilot uses a recorded deterministic selector; there is no unrecorded randomness.

Scoped-null sources consulted: `00-prompt.md`; `reference/ga-notes.md`; `parley-deck/COOPERATION.md`, especially §15 and §15.6; `protocol-generation-bias/FINAL.md`; `openviking-context-structure/FINAL.md`; live outputs of `parley roster show` and `parley preset list`; `~/.parley/agents.toml`; and the 2026-08-29 protocol diff at commit `59eb663d0c98cdc6f524498a70ce765d7bde6992`. I did not consult any other participant's round-one artifact for this idea.

## Concerns / open questions

- **Correct convergence remains unlabeled.** The 80-idea round-count measurement in the kickoff establishes brevity, not error. What external outcome will label a converged design as premature rather than correct—implementation reversals, later `-v2` ideas, review findings, owner rejection, or something else?
- **Behavioural tuples contain judgment.** Two proposals can share a mechanism but differ in a load-bearing constraint. The tuple schema and blind adjudicator must be fixed before outputs are seen, or the evaluator can manufacture either diversity or convergence.
- **The pilot may not transfer to live cooperation.** Sealed single-call pairs isolate generation, while Parley's value may arise in cross-review and correction. Conversely, testing inside a live idea contaminates the baseline. Both limitations must be reported.
- **Model-family diversity is not evidence independence.** The live roster has six vendors/families, but shared training corpora, instruction conventions, or gateway behaviour may still correlate outputs. §15.6(b) correctly treats unanimity as a shared prior; the repository does not quantify the correlation.
- **A prompt can create performative dissent.** “Emit only what you endorse” is a guard, not a proof that the model has an authentic minority position. The null rate and later withdrawal rate should therefore be reported, not treated as failures to comply.
- **The pilot thresholds are policy choices.** The proposed two batches and gates are deliberately conservative, not literature-derived. The owner should approve the spend and thresholds before calls are made.
- **External support is mostly secondary here.** The GA, QD, and FunSearch claims used above retain the provenance assigned in `reference/ga-notes.md`; I did not independently browse or upgrade them to PRIMARY.

## Risks

- A mutation duty can recreate the exact removed steelman clause under a new name, adding protocol weight without new evidence.
- Artificial dissent can consume the only cross-review round refuting an idea nobody genuinely supports and make consensus worse.
- A behavioural hash can collapse two materially different proposals or reward superficial category changes; both false convergence and false novelty are possible.
- Blind scoring can still be correlated if evaluators are drawn from the same model family or see stylistic cues that reveal the arm.
- Varying model, effort, temperature, and lens together makes the treatment uninterpretable; only the lens should differ in the proposed pair.
- A benchmark can become a permanent ritual even when it is inconclusive. Failure or ambiguity must terminate the line rather than lower the gate.
- The null decision also has a cost: Parley may continue missing valuable frames. That is why the recommendation preserves a bounded, owner-approved falsification path instead of declaring mutation impossible.
