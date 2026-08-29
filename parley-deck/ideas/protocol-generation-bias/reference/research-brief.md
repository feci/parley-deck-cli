---
idea: protocol-generation-bias
artifact: reference
author: claude-1 (facilitator, via verification workflow)
date: 2026-08-28
---

# Research brief — external evidence, checked

## What this is

This is a **reference artifact**. It proposes nothing, endorses no design, and recommends no rule.
Where the evidence supports competing designs it says so and presents both.

Every claim below was **fetched and checked by an agent other than the one that first cited it**.
Round-01 citations were re-resolved from DOI/arXiv ID to publisher record, and — wherever the host
permitted — to full text. Sweeps of adjacent literature (human fixation and premature closure;
institutional alternatives-forcing practice; LLM multi-agent diversity) were run independently of
round-01's reading list. Deck-internal numbers under dispute in round 01 were re-derived by a
reproducible script.

This file exists so that round 02 does not build on an unchecked citation, an inverted finding, or a
number whose denominator nobody agreed on. It is a floor to argue from, not a position to argue for.

**Provenance conventions used here.** A claim marked **MEASURED** comes from a study reporting
outcome data. **MODELED** comes from a proof or simulation with no human or LLM subjects.
**DOCTRINE** comes from a manual or standard with no evaluation attached. **UNVERIFIED** means the
source could not be opened and the claim rests on secondary indexing — these are retained rather than
dropped, and are never to be quoted as established.

**One standing caution that applies to the whole file.** The great majority of the strongest,
best-measured results here are about **human** groups. Their transfer to LLM agents is **not
established** by those studies. Where an LLM-specific result exists it is flagged as such. Round 02
should not silently upgrade a human finding into a claim about this deck.

---

## Citation audit

| Citation | Cited by | Verdict | What it actually is |
|---|---|---|---|
| 10.1126/science.185.4157.1124 | codex-1 | **CONFIRMED** | Tversky & Kahneman, "Judgment under Uncertainty: Heuristics and Biases", *Science* 185(4157):1124–1131, 1974. Quote verbatim and complete. |
| 10.1037/0022-3514.48.6.1467 | codex-1 | **CONFIRMED** | Stasser & Titus, "Pooling of Unshared Information in Group Decision Making: **Biased Information Sampling During Discussion**", *JPSP* 48(6):1467–1478, 1985. Subtitle was dropped in the citation. |
| 10.1016/S0169-2070(99)00018-7 | codex-1 | **CONFIRMED** | Rowe & Wright, "The Delphi technique as a forecasting tool: **issues and analysis**", *IJF* 15(4):353–375, 1999. Title truncated; **four** defining features, not three. |
| 10.1002/ejsp.58 | claude-1, hermes-1 | **MISQUOTED** | Nemeth, Brown & Rogers, EJSP 31(6):707–720, 2001. Real paper, **wrong finding attached**. See below. |
| 10.1111/j.1559-1816.2001.tb02481.x | (facilitator) | **CONFIRMED** | Nemeth, Connell, Rogers & Brown, "Improving Decision Making by Means of Dissent", *JASP* 31(1):48–58, 2001. This is the correct home of the cognitive-bolstering finding. |
| 10.1126/science.1091721 | zcode-1 | **CONFIRMED (identity only)** | Johnson & Goldstein, "Do Defaults Save Lives?", *Science* 302(5649):1338–1339, 2003. A choice-architecture paper; **no group, no dissent, no conformity**. |
| 10.5465/255641 | zcode-1 | **CONFIRMED (metadata only)** | Van de Ven & Delbecq, *AMJ* 17(4):605–621, 1974. Full text not obtained; findings from indexed summaries. |
| 10.1016/0142-694X(91)90011-F | zcode-1 | **WRONG-ATTRIBUTION (dead DOI)** | The DOI **does not exist**. Correct DOI for Jansson & Smith, "Design fixation" is **10.1016/0142-694X(91)90003-F**. |
| 10.1073/pnas.0403723101 | kimi-1 | **CONFIRMED** | Hong & Page, "Groups of diverse problem solvers can outperform groups of high-ability problem solvers", *PNAS* 101(46):16385–16389, 2004. **MODELED**, and contested. |
| Madrian & Shea, "The Power of Suggestion" | zcode-1 | **CONFIRMED** | *QJE* 116(4):1149–1187, 2001; DOI 10.1162/003355301753265543. Locator guess was correct. |
| arXiv:2305.14325 | kimi-1 | **CONFIRMED** | Du, Li, Torralba, Tenenbaum, Mordatch, "Improving Factuality and Reasoning in Language Models through Multiagent Debate". Cite as **ICML 2024**, not a bare 2023 preprint. |
| arXiv:2305.19118 | kimi-1 | **CONFIRMED (one use over-extended)** | Liang et al., "Encouraging Divergent Thinking in LLMs through Multi-Agent Debate" (EMNLP 2024). First use accurate; second use inverts the thesis. See below. |
| arXiv:2203.11171 | kimi-1 | **CONFIRMED** | Wang et al., "Self-Consistency Improves Chain of Thought Reasoning in Language Models", ICLR 2023. Tagged RECALL by kimi-1; the tag was honest and can now be upgraded. |
| arXiv:2310.13548 | kimi-1, zcode-1 | **CONFIRMED** | Sharma et al., "Towards Understanding Sycophancy in Language Models" (ICLR 2024). Both agents' quotes verbatim. |
| arXiv:2404.13076 | kimi-1 | **CONFIRMED** | Panickssery, Bowman & Feng, "LLM Evaluators Recognize and Favor Their Own Generations" (NeurIPS 2024). |
| 10.18653/v1/2024.acl-long.511 | codex-1 | **CONFIRMED** | Wang et al., "Large Language Models are not Fair Evaluators", ACL 2024. Cleanly quoted; no defect found. |
| arXiv:2311.17371 | zcode-1 | **MISQUOTED** | Smit et al., "Should we be going MAD?" (ICML 2024). Real quote plus an **invented clause**, and the qualifying sentences were omitted. See below. |
| arXiv:2605.00914 | claude-1 | **CONFIRMED** | Bertalanič & Fortuna, "The Cost of Consensus". claude-1 was *too* conservative — the headline numbers are in the public abstract. |
| arXiv:2605.00914 (2nd citation) | hermes-1 | **WRONG-ATTRIBUTION** | Cited as "Chen et al. (2025/2026)". **There is no author named Chen on this paper.** See below. |
| arXiv:2605.30150 | claude-1 | **CONFIRMED (characterisation disputed)** | Ibrahim, Azad & Baten, "Anchorless Diversification for Parallel LLM Ideation". Its headline finding points **against** the mechanism claude-1 cited it as prior art for. |
| arXiv:2606.19494 | claude-1 | **CONFIRMED (semantic trap)** | Pokharel & Dantu, "Hidden Anchors in Multi-Agent LLM Deliberation". Its "anchor" is an agent's own **internal** belief, and its effect is **positive**. Correctly tagged RECALL by claude-1. |
| arXiv:2606.01637 | codex-1 | **CONFIRMED** | Qu, Fu & Hu, "Easier to Mislead Than to Correct". Quote verbatim word-for-word; the best-executed citation of the round. |

Two citations in the batch were **facilitator-supplied**, not cited by any agent in round 01:
arXiv:2606.00820 (Hao et al., "Not All Flips Are Conformity") and arXiv:2606.29270 (He et al.,
"Minority Sentinel"). Both **CONFIRMED**. They matter because two numbers already loose in the
deliberation belong to them — see the note under "Numbers already circulating without owners".

### Citations that did not survive

**1. 10.1002/ejsp.58 — MISQUOTED. Cited by claude-1 and hermes-1.**

The bibliographic record is exactly right: Nemeth, Brown & Rogers, "Devil's advocate versus authentic
dissent: stimulating quantity and quality", *EJSP* 31(6):707–720, 2001. The **claim pinned to it is
not this paper's finding**. Its abstract attributes cognitive bolstering to "**a prior study**", and
states that the present study *replaced* cognitive activity with quantity and quality of solutions as
the dependent measure.

What this paper actually reports, verbatim: *"Results indicated that the authentic minority was
superior to all three forms of 'devil's advocate,' again underscoring the value and importance of
authenticity and the difficulty in cloning such authenticity by role-playing techniques."* The three
forms varied whether the advocate's true position was known and whether it was consistent with the
assigned role — so the failure is **not** repaired by making the assigned role transparent.

The cognitive-bolstering result belongs to the companion paper, **10.1111/j.1559-1816.2001.tb02481.x**
(Nemeth, Connell, Rogers & Brown, *JASP* 31(1):48–58, 2001), whose abstract states it in its own
voice: *"Devil's advocate was found to stimulate cognitive bolstering of the initial position."*

Correct split for round 02: cite **JASP** for bolstering; cite **EJSP** for "authentic minority beat
all three devil's-advocate forms; role-play cannot clone authenticity." A sentence of the form
"Nemeth et al. (2001, EJSP) found that assigned devil's advocacy produces cognitive bolstering"
attributes the wrong result to the wrong paper.

**2. 10.1016/0142-694X(91)90011-F — WRONG-ATTRIBUTION (dead DOI). Cited by zcode-1.**

The DOI **returns 404 from both doi.org and the Crossref API**, in upper and lower case. The paper is
real and the intended attribution is right; the locator is wrong. Correct DOI:
**10.1016/0142-694X(91)90003-F** — Jansson & Smith, "Design fixation", *Design Studies* 12(1):3–11,
1991. Full text was read.

**3. arXiv:2311.17371 — MISQUOTED. Cited by zcode-1.**

Attribution is correct (Smit, Duckworth, Grinsztajn, Barrett, Pretorius, "Should we be going MAD? A
Look at Multi-Agent Debate Strategies for LLMs", ICML 2024). Two defects.

- The quoted clause *"do not reliably outperform other proposed prompting strategies"* **is verbatim**.
  The appended clause *"and their dissent tends to be shallow"* **appears nowhere in the abstract**
  and was inside the same PRIMARY bracket. It is unsupported by what was read.
- The omitted continuation reverses the weight: *"However, when performing hyperparameter tuning,
  several MAD systems, such as Multi-Persona, perform better. This suggests that MAD protocols might
  not be inherently worse than other approaches, but that they are more sensitive to different
  hyperparameter settings and difficult to optimize."* And: adjusting agent agreement levels *"can
  significantly enhance performance and even surpass all other non-debate protocols we evaluated."*

Quoted in full, the paper says debate is *not reliably better as currently configured, but tunable* —
not *debate cannot work*. Round 02 should either re-quote with the qualifier or stop using it for that
step.

**4. arXiv:2605.00914 cited as "Chen et al. (2025/2026)" — WRONG-ATTRIBUTION. hermes-1.**

**Highest-priority fix in the batch.** The same arXiv ID is cited by two agents with two different
author lists. The paper is by **Blaž Bertalanič and Carolina Fortuna** (Jožef Stefan Institute). There
is no author named Chen. hermes-1 tagged the citation SECONDARY and stated the "abstract and title
verified via web search" — so the invented byline is a **verification-integrity failure**, not a typo:
the provenance tag asserts a check that cannot have happened. The substantive claim hermes-1 draws
(homogeneous debate produces sycophantic conformity at poor cost-to-accuracy ratios; isolated
self-correction outperforms unguided debate) **is** supported by the abstract, so the argument
survives correction of the byline. hermes-1's further inference — that a divergence-enforcement
mechanism "is functionally another debate round — exactly the mechanism the paper says fails" — is the
agent's own reasoning and does not inherit the citation's authority.

### Two citations that are correct but semantically trapped

Neither is a misquote. Both are one careless sentence away from becoming one.

- **arXiv:2606.19494, "Hidden Anchors in Multi-Agent LLM Deliberation".** The most on-topic *title* in
  the batch and the easiest to invert. Its "anchor" is an agent's **own hidden internal belief** that
  pulls its opinion against its neighbours — not a first-mover, ordering, or seed-proposal anchor. And
  its notable consequence is **positive**: verbatim, *"an agent's confidence in the correct answer can
  climb past where any agent started, escaping the space (convex hull) formed by the initial
  beliefs."* Citing it as evidence that anchors are a bias to be designed out **inverts the paper**.
  claude-1 handled it correctly, marking it unread beyond title and authors and tagging findings
  RECALL.
- **10.1126/science.1091721 ("Do Defaults Save Lives?").** Real, correctly attributable, and relevant
  to **defaults and status quo** — the paper's own mechanism list names the status quo explicitly. It
  has **no** relevance to dissent or conformity: no group, no minority influence, no deliberation. It
  arrived batched under dissent-conformity. Any dissent- or conformity-flavoured claim resting on this
  DOI should be treated as unsupported until zcode-1's exact sentence is checked.

### Numbers already circulating without owners

claude-1.md:102–104 flags "29% strict conformity, 57–77% correct-to-wrong" and a "32.3pp oracle gap"
as unattributed figures circulating in search results. **All three numbers are real and correctly
valued, and they come from two different papers.**

- **29% strict conformity** and **57–77% correct-to-wrong** are from **arXiv:2606.00820** (Hao et al.,
  "Not All Flips Are Conformity").
- **32.3pp oracle gap** is from **arXiv:2605.00914** (Bertalanič & Fortuna).

Round 02 must attribute them separately rather than treating them as one undifferentiated pool.

---

## What the evidence actually supports

Organised by claim. Each claim carries its sources, the effect measured, and the limits.

### C1. An anchor can be introduced by the *framing* of a problem, not only by an explicit number

**Source.** Tversky & Kahneman 1974 (10.1126/science.185.4157.1124), "Adjustment and Anchoring",
verbatim: *"In many situations, people make estimates by starting from an initial value that is
adjusted to yield the final answer. The initial value, or starting point, may be suggested by the
formulation of the problem, or it may be the result of a partial computation. In either case,
adjustments are typically insufficient."*

**Limits.** The supporting experiments are **numerical estimation** (the 8! ascending/descending
product task, median 512 vs 2,250 against a correct 40,320; the wheel-of-fortune UN task). Applying
this to open-ended text generation by LLM agents is an **extrapolation beyond the paper's evidence**.
Cite it for the mechanism, never as measured evidence about agents.

### C2. Discussion can entrench a group's initial skew rather than correct it

**Source.** Stasser & Titus 1985 (10.1037/0022-3514.48.6.1467), verbatim from the abstract:
*"discussion tended to perpetuate, not to correct, members' distorted pictures of the candidates."*
**Corroborated and quantified** by Lu, Yuan & McLeod 2012 (10.1177/1088868311417243), a meta-analysis
of 65 studies / 101 effects / **3,189 groups**: groups mentioned **two standard deviations** more
common than unique information, and hidden-profile groups were **eight times** less likely to find the
correct solution than groups given full information.

**A finding that changes what to measure.** Of the two pooling measures that predicted decision
quality — *information coverage* (share of unique information that reached the table) and *discussion
focus* (share of discussion airtime spent on it) — **coverage was the stronger predictor**. Rewarding
an agent for *voicing* an alternative optimises the weaker variable.

**Limits.** Human face-to-face groups, political-caucus simulation, deliberately constructed hidden
profiles. This is a **different mechanism from anchoring** — it is about which information gets voiced,
not about adjustment from an anchor. Do not blur it into C1.

**LLM-specific counterpart, and it is stronger than the transfer argument.** HiddenBench
(arXiv:2505.11556, Li, Naito & Shirado, 65 tasks): multi-agent LLMs score **30.1%** when information
is distributed across agents versus **80.7%** for a single agent handed everything — a **50.6pp**
collective-reasoning deficit. Per-model post-discussion: Gemini-2.5-Pro 67.1%, Gemini-2.5-Flash 55.0%,
GPT-4.1 23.3%. The deficit **worsens with scale**: 3 agents +34.8% pre-to-post-discussion, 4 agents
+25.0%, 7 agents **+0.6%**; extended communication peaked at 15 rounds then declined.

### C3. A seed example contaminates downstream generation even when its features are known-bad and explicitly prohibited

**Source.** Jansson & Smith 1991 (10.1016/0142-694X(91)90003-F), full text read. Four experiments;
control gets the problem, "fixation" group gets the same problem plus one example solution.

- Exp 1 (bicycle rack, n=25): designs per subject 4.5 vs 4.3 — **no difference in output volume** —
  but suction cups **6% vs 54%**, tyre railings 15% vs 48%.
- Exp 2 (measuring cup for the blind, n=31): designs lacking an overflow device 17% vs 54%; designs
  "highly similar to example" **7% vs 50%**.
- Exp 3 (spill-proof cup, n=35): instructions **forbade straws and mouthpieces**; the example used
  both. Fixation group still produced straws **1% vs 17%** and mouthpieces **10% vs 39%**.
- Exp 4 (13 **professional design engineers**): cords 36% vs 78%, control box 73% vs 100%, front
  opening **9% vs 56%**. Fixation replicated in professionals.

**Two structural facts.** Output *quantity* was unaffected in every experiment — fixation is invisible
in throughput and shows only in the distribution of features. And a stated prohibition was
**empirically insufficient** to keep the prohibited feature out.

**Methodological caveat that must travel with any quotation.** The paper reports **only group means
and percentages**. No p-values, no confidence intervals, no significance tests, no effect sizes
anywhere; cell sizes are 6–18 people. Any round-02 text calling these results "statistically
significant" or quoting a *d* would be **inventing it**. The largest raw gaps are the defensible ones:
+48pp suction cups (Exp 1), +43pp "highly similar to example" (Exp 2), +47pp cords (Exp 4).

**Limits.** Seeing an **artefact** (an example solution), not a numeric estimate. This supports
"example contamination", not classic numeric anchoring.

**LLM counterpart.** Deng, Brucks & Toubia (arXiv:2602.20408) name individual-level fixation in LLMs
directly: *"early outputs constrain subsequent ideation."* Their second named mechanism is
collective — LLMs compress into a single unified distribution rather than reproducing the knowledge
**partitioning** of a human population.

### C4. Removing the incumbent from the board recovers performance that instructions cannot

**Source.** Bilalić, McLeod & Gobet 2008 (10.1016/j.cogpsych.2007.02.001), full text read. A chess
position containing both a familiar 5-move smothered mate and a shorter optimal 3-move mate.

- With the familiar solution available: Grand Masters 100%, International Masters **50%**, Masters
  **18%**, Candidate Masters **0%** found the optimal move.
- With the familiar solution **disabled** by moving one black bishop — same optimal move, same
  players' skill classes: International Masters **100%**, Masters **100%**, Candidate Masters **100%**.
- Quantified: *"The presence of a familiar solution reduced the problem solving abilities of the
  experts to about that of players 3 SDs lower in skill."*

**And a second, independent finding: self-report of search carries no information.** Verbal protocols
show the smothered mate was the first thing **all** players noticed, and all then reported looking for
something shorter. Those who failed spent **8–9s** looking; those who succeeded spent **7–10s**.
Metacognition failed too: 12 experts shown both solutions predicted 86%/74%/59% success for
IMs/Masters/CMs; actual was **50%/18%/0%**.

**Corroborated in the design domain.** Linsey et al. 2010 (10.1115/1.4001110): engineering design
**faculty** — people who research and teach design methods — *"show statistically significant evidence
of design fixation, but only partially perceive its effects."* Their conclusion: *"designers, even
those that study and teach design on a regular basis, do not know when they are being influenced or
fixated by misleading or poor information."*

**Limits.** Chess and engineering design, humans. No LLM replication exists.

### C5. Fixation is dangerous precisely when the incumbent is *good*, not when it is obviously bad

**Source.** Sheridan & Reingold 2013 (10.1371/journal.pone.0075796), 34 players, eye-tracked. When the
familiar move was **advantageous but suboptimal**, 47% of experts and 47% of novices took it —
expertise gave no protection. When the familiar move was an **outright blunder**, all experts and most
novices avoided it and disengaged attention from it.

**Consequence for gating.** This inverts the natural intuition that scrutiny should be spent where
objections exist. A proposal that draws **no objections** is, on this evidence, the higher-risk case.
Round 02 will have to decide whether that argues for an unconditional obligation or for triggering on
**consensus** rather than dissent — the evidence establishes the risk profile, not the trigger.

**Limits.** Chess, humans, eye-tracking as the mechanism (attention biased toward features associated
with the familiar solution).

### C6. *Structured, imposed* re-analysis works; an *opportunity* to reconsider does not

This pair is the cleanest natural experiment in the whole sweep, run on overlapping populations by
overlapping research groups.

- **Works.** Mamede et al. 2010 (10.1001/jama.2010.1276), 36 internal medicine residents. Availability
  bias demonstrated (1.55 on similar cases vs 2.19 on other cases, p=0.03). Then **experimenter-imposed
  structured reflection** produced a significant main effect of reasoning mode, improving accuracy on
  the biased cases for second-years (2.03) and first-years (2.31), **p=0.006**.
- **Does not work.** Monteiro et al. 2015 (10.1007/s11606-015-3369-4), 47 residents diagnosing 16
  cases, then **explicitly given the opportunity** to reflect and revise. Residents *did* engage — the
  door was used. Only **8%** of diagnoses were revised, and accuracy moved **1.20/2 → 1.22/2**
  (t=2.15, p=0.03): statistically detectable, practically negligible. Authors: *"this strategy
  provided minimal benefits."*

**The differentiating variable is imposed-and-structured versus offered-and-self-directed.** Not
reflection per se.

**Aggregate.** Lambe et al. 2016 (10.1136/bmjqs-2015-004417), systematic review of 28 studies in five
categories: *"guided reflection interventions emerged as the most consistently successful across five
studies."* Their own hedge should travel with it: *"further research with refined methodology and more
diverse samples is required before firm recommendations may be made"*, and the evidence base is
largely early-career doctors.

**Limits.** Medical diagnosis, humans, narrow task. No LLM analogue was found.

### C7. The failure being chased is *stopping*, not *not-knowing*

**Source.** Graber, Franklin & Gordon 2005 (10.1001/archinte.165.13.1493), 100 cases of diagnostic
error; 548 contributing factors across 93 non-no-fault cases (5.9 per case); cognitive factors in 74%.
Verbatim: *"Premature closure, ie, the failure to continue considering reasonable alternatives after an
initial diagnosis was reached, was the single most common cause."* And: *"Faulty or inadequate
knowledge was uncommon."*

**Why this bounds the design space.** If the failure is stopping rather than ignorance, adding
capability — a stronger model, a bigger roster, more context — does not by itself address it. Whether
the right response is a changed **exit criterion**, a changed **artifact**, or something else is
exactly what round 02 has to decide; this finding rules out one family, not in another.

**Limits.** Human clinicians, retrospective case review with known selection bias toward discovered
errors.

### C8. Independent generation before interaction beats interacting generation, on both quantity and quality

**Source.** Mullen, Johnson & Salas 1991 (10.1207/s15324834basp1201_1), full PDF read. 18 articles / 20
studies; **34 hypothesis tests for quantity (2,577 individuals in 844 groups)** and 9 for quality.
Combined effect for **quantity r = .572** (Z = 15.324, p = 6.16E-35, fail-safe N = 2827); for
**quality r = .558** (Z = 10.592, p = 1.25E-22, fail-safe N = 253). The "quality compensates for lost
quantity" hypothesis is **refuted**, not merely unsupported. Productivity loss **grows with group
size** (r = .606 quantity; r = .715 quality) and is largest when the comparison is against individuals
working truly **Alone** rather than merely Together.

**Delphi's ingredient set points the same way with a boundary.** Rowe & Wright 1999/2001: Delphi beats
traditional interacting groups **5:1 with one tie** and statistical groups **12:2 with two ties** —
but there is *"no consistent evidence that the technique outperforms other structured group
procedures."* So the measured gain is attributable to **structure** (anonymity, iteration, controlled
feedback, and — the fourth feature the round-01 citation omitted — **statistical aggregation of group
response**), not to Delphi's particular ceremony. Van de Ven & Delbecq 1974 (10.5465/255641) is the
founding NGT comparison but its **full text was not obtained**; the quantitative weight for this claim
rests on Mullen et al., which was.

**Limits.** Human brainstorming and judgment tasks. Two moderators do map suggestively onto a deck —
loss worsens with group size, and the "Together" condition (co-present but not conversing) forfeits
most of the benefit, which is structurally what a shared context window is — but **that mapping is an
inference, not a finding**.

**LLM-side echoes, and they are mixed.** Zhu et al. (arXiv:2601.19921) show diversity-aware
initialization raises Pass@5 from 0.7921→0.9097 (Qwen) and 0.7423→0.9026 (Llama), with +7.1pp on
MMLU-FL over majority voting. But Self-MoA (arXiv:2502.00674) runs the other way on **roster**
diversity: aggregating multiple samples from the single best model beat mixing different LLMs by
**+6.6%** on AlpacaEval 2.0 and **+3.8%** on average, because *"mixing different LLMs often lowers the
average quality of the models."* The two are reconcilable — diversity of *prompt/direction* is
supported, diversity of *participant quality* is not — but round 02 should not cite "diversity helps"
as an undifferentiated claim.

### C9. A protocol can produce artifacts reliably even when its quality effect is unproven

**Source.** DCI (arXiv:2603.11781, Prakash), evaluated on 45 tasks across seven domains. Artifact
production: **decision packets 100% of the time vs ≤16% for baselines; minority reports 98% vs 0%**.
The packet carries the selected option, residual objections, a minority report preserving dissent with
confidence, next actions, and explicit reopen triggers. It also defines **bounded openness**: a new
option may enter late only if materially distinct, plausibly superior, evidence-linked, and before a
cutoff round.

**And the same paper is the strongest caution against heavyweight protocol.** Across all 45 tasks a
**single agent beat it** (8.84 vs 8.24). On routine tasks DCI scored **5.39, −3.19 versus debate**. It
cost **~62× single-agent tokens** (quality-per-token 0.035 vs 2.320). **All three ablations scored
equal or higher than full DCI** (no archetypes +0.37, no typed grammar +0.08, no DCI-CF +0.09) — the
author concedes component contributions are not separable at that sample size. Its gains are confined
to non-routine tasks (+0.95 vs debate, 95% CI [+0.41, +1.54]) and hidden-profile integration (9.56,
its best domain score).

**The honest reading:** the artifact discipline is well-evidenced; the deliberation machinery is not.

### C10. Order of presentation, authority labels, and peer agreement measurably move LLM judgments

- **Order.** Wang et al. ACL 2024 (10.18653/v1/2024.acl-long.511), verbatim: *"the quality ranking of
  candidate responses can be easily hacked by simply altering their order of appearance in the
  context... Vicuna-13B could beat ChatGPT on 66 over 80 tested queries with ChatGPT as an
  evaluator."* Their own remedy, Balanced Position Calibration, evaluates candidates in swapped
  positions.
- **Authority and peer agreement.** Qu, Fu & Hu (arXiv:2606.01637), four open-weight LLMs, seven QA
  datasets, verbatim: *"peer agreement makes it much easier to mislead initially correct models than to
  correct initially wrong ones. Authority labels make models more likely to choose the endorsed answer,
  regardless of whether it is correct."* Note: the **abstract contains no numbers** — any effect size
  attributed to this paper must come from the PDF. And the paper's own prescription is narrower than
  blinding: *"multi-agent LLM systems should verify peer answers rather than simply aggregate them."*
  It never tests a blinding intervention, and its chain-of-thought/reflection result shows generic
  reasoning interventions did **not** reliably help.
- **Sycophancy toward a stated view.** Sharma et al. (arXiv:2310.13548): five state-of-the-art
  assistants *"consistently exhibit sycophancy across four varied free-form text-generation tasks"*,
  and *"when a response matches a user's views, it is more likely to be preferred."* **Scope limit:**
  measured toward a **human user's** stated views in single-model interaction. It does **not** study
  agent-to-agent anchoring. Both kimi-1 and zcode-1 flagged their mapping as inference; round 02 must
  not upgrade it.
- **Self-preference.** Panickssery, Bowman & Feng (arXiv:2404.13076): LLMs distinguish their own
  outputs at non-trivial accuracy, and there is *"a linear correlation between self-recognition
  capability and the strength of self-preference bias."* Scope: single-model evaluation of candidate
  texts, not multi-agent protocol adjudication.

### C11. Convergence is not evidence of correctness, and majorities discard correct minorities at a measurable rate

- **arXiv:2605.00914** (Bertalanič & Fortuna), 10-agent homogeneous teams over 3 rounds: debate **never
  beat** isolated self-correction in any configuration. Qwen2.5-7B/MMLU-Hard 60.7% vs 66.7%
  (p<0.001); Ministral-3-8B/GSM-Hard **20.7% vs 48.3%** (p<0.001). Debate consumed **2.1–3.4×** more
  tokens (up to 28,631 per problem). Three named pathways: sycophantic conformity (modal adoption up to
  **85.5%**), contextual fragility (vulnerability up to **70.0%**), consensus collapse (oracle gap up
  to **32.3pp**). Teams reached **90.1% consensus while accuracy stagnated or fell**. Most damning: an
  **irrelevant-rationale noise control sometimes outperformed real debate**.
  **Scope limit, and it is strong:** N=10 agents, R=3 rounds, only Qwen2.5-7B / Llama-3.1-8B /
  Ministral-3-8B, only GSM-Hard and MMLU-Hard, and the authors confine the conclusion to *"within the
  7–8B parameter class"* with *"homogeneous teams without structured roles"*. **This paper cannot
  license any claim about heterogeneous frontier-model debate.**
- **arXiv:2606.29270** (He et al., "Minority Sentinel"), three **heterogeneous** agents: among divergent
  2:1 cases, **the minority holds the correct answer 25.5% of the time**, a ~10pp theoretical recovery
  margin. Diagnosis: *"Because contemporary LLMs share similar pretraining corpora, their errors are
  strongly correlated, causing the majority to systematically suppress correct minority opinions."*
  Adding agents from the same pretraining distribution does not buy independence.
  Their classifier achieves **81.2% Flip Precision** (not "81.2% accurate" — a distinct quantity) with
  positive Net Gain across all six benchmarks and all 20 seeds (+1.71% over a 74.3% baseline).
- **These two bound the question rather than settling it.** 2605.00914 is homogeneous same-model teams;
  2606.29270 is heterogeneous. Whether a six-lab roster inherits the homogeneous-debate pathology is
  **not answered by either alone**.

### C12. Defaults change outcomes at scale, and the counterfactual is measured

**Source.** Madrian & Shea 2001 (10.1162/003355301753265543). Two distinct effects that must not be
blurred: (a) a default changes **participation** — 401(k) participation at 3–15 months tenure was
**37% without auto-enrollment vs 86% with**; (b) a default changes the **content** of the choice —
**61%** of the auto-enrolled cohort sat at the full default tuple (participate + 3% + 100% money
market) versus **1%** of every other cohort. The 1%-vs-61% comparison is the strong form: almost none
of the default-conforming behaviour is behaviour people would have chosen anyway. The effect is
strongly income-graded: **74%** of employees earning under $20,000/yr stayed at the default vs **30%**
of those earning $70,000–79,999.

**Limits and provenance.** Single large U.S. corporation, cohort comparison around a policy switch —
a natural experiment, not randomized. Numbers above were read from the authors' **NBER WP 7682**
full text; the QJE version is paywalled. The commonly circulated two-decimal figures (37.4% → 85.9%)
were **not** verified. The authors name *"anchoring around the default and a bias for the status quo"*
plus employees reading the default as an **implicit company recommendation** — and that second
mechanism has no analogue in a purely mechanical software default.

### C13. Formal argument that heterogeneous solvers can beat individually stronger similar solvers

**Source.** Hong & Page 2004 (10.1073/pnas.0403723101), verbatim: *"a team of randomly selected agents
outperforms a team comprised of the best-performing agents... as the initial pool of problem solvers
becomes large, the best-performing agents necessarily become similar in the space of problem solvers.
Their relatively greater ability is more than offset by their lack of problem-solving diversity."*

**Two caveats that must travel with it.** (1) **MODELED** — a proof plus simulations, agents defined
as (representation, heuristic) pairs on a value landscape, with stated conditions (large initial pool,
a difficulty condition, agents that always find a local optimum). **No humans, no data, no LLMs.**
Never cite as measured. (2) **Contested** — Thompson, "Does Diversity Trump Ability?", *Notices of the
AMS* 61(9):1024–1030 (2014) argues the popular reading does not follow and that randomness rather than
diversity drives the result; see also Kuehn (*Critical Review* 29(1), 2017), Singer ("Diversity, Not
Randomness, Trumps Ability", *Philosophy of Science*), and arXiv:2307.04709.

### C14. Independently, two lines of evidence say "declared criteria applied consistently" predicts accuracy

The only positive result to emerge from an otherwise negative ACH trial (Dhami, Belton & Mandel 2019,
10.1002/acp.3550, 50 practising UK intelligence analysts): analysts who applied their evidence-
assessment rule **consistently across evidence items** were far more accurate — **80% (4/5) vs 31%
(9/29), p=.021, odds ratio 8.89**. The predictor of accuracy was **consistency of rule application**,
not use of the technique.

Independently, NASA NPR 7123.1D §3.2.18.2 orders its mandated Decision Analysis process as
*"identification of decision criteria, identification of alternatives, analysis of alternatives, and
alternative selection"* — **criteria before candidates**.

**Limits.** The 8.89 odds ratio rests on very small cells (5 and 29). NASA's ordering is **DOCTRINE**
with no evaluation attached. These converge, but neither is strong on its own.

---

## Interventions with evidence AGAINST them

This section is where the sweeps paid for themselves. Each item below is widely recommended and was
**measured not to work**, or measured to work in a direction other than the one usually claimed.

**A1. "Be different from the examples" / generic divergence instructions — measured null.**
Smith, Ward & Schumacher 1993 (10.3758/bf03202751), Exp 3: *"explicitly instructing subjects to create
ideas that were very different from the examples did not decrease conformity to the examples, and
instructing them to conform to the examples significantly increased conformity."* The asymmetry is the
point — the instruction was read and obeyed in the conform direction, so the null in the diverge
direction is **not** a manipulation failure. Corroborated by Jansson & Smith Exp 3, where an explicit
prohibition on straws and mouthpieces still yielded 17% and 39% (see C3).

**A2. Anti-anchoring against named examples, in LLMs — measured worst of its family.**
Ibrahim, Azad & Baten (arXiv:2605.30150) compared six diversification schemes. The **repr** condition —
all calls instructed to avoid 3 shared representative anchors — was the **worst anchorless method**
(entropy AUC 0.49–0.52 vs 0.61 for semantic-direction stratification), and adding a divergence
instruction **hurt it** (entropy gain 0.090 with diverge vs 0.120 without). Peer anchoring's apparent
edge disappears under full-pipeline token accounting: *"anchored regeneration can be strong in
final-pool diversity, but its advantage shrinks under full-pipeline token accounting."* Note this
directly qualifies claude-1.md:99–100's framing of this paper as prior art confirming a
peer-anchored divergence mechanism — the paper's headline finding runs the other way. The paper's
abstract contains **no numeric results**; any percentage attributed to it without opening the PDF is
invented.

**A3. Delay / incubation as a defixation device — measured null for fixation.**
Smith et al. 1993 Exp 2: a 23-minute interpolated task did **not** significantly reduce conformity.
Sio & Ormerod 2009 (10.1037/a0014212) do find a positive incubation effect overall, but moderated —
divergent-thinking tasks benefit most, and filling the period with **high-cognitive-demand** work
shrinks the effect. Their own conclusion: *"the conditions under which incubation can be used as a
practical technique for enhancing problem solving must be designed with care."* An extra deliberation
round is a high-demand filler. The supported form is a **longer preparation phase before** the first
proposal, not more rounds after it.

**A4. Expertise and self-awareness as protection — measured absent.**
Jansson & Smith Exp 4 (professional engineers) and Linsey et al. 2010 (design faculty) both fixate
significantly; Bilalić 2008 quantifies the cost at 3 SDs of skill. Metacognition fails in both
domains: faculty *"only partially perceive"* fixation; 12 chess experts predicted 86%/74%/59% and the
truth was 50%/18%/0%. Self-reported search carries **no information** — failers and finders spent
identical time looking (8–9s vs 7–10s).

**A5. Teaching debiasing as a skill — measured null at n=191.**
Sherbino, Kulasegaram, Howey & Norman 2014 (10.2310/8000.2013.130860), controlled trial, 191 medical
students, cognitive forcing strategy training vs control. Every measure null: initiated a search for a
second diagnosis **52% vs 48% (p=0.91)**; correctly identified it **54% vs 48% (p=0.13)**;
false-positive over-identification 64% vs 77% (p=0.12); uncommon correct diagnosis in availability
cases **45% vs 45% (p=0.98)**. Their 2011 pilot had already found *"application and retention is
poor."* Note the trial also included a **false-positive arm** to detect over-searching — any
alternative-forcing rule should carry the equivalent check that it does not manufacture spurious
alternatives when the incumbent was in fact correct.
*Tension to note honestly:* Lambe et al. 2016's review says cognitive forcing strategies "improved
accuracy and confidence judgements", and Sherbino 2014 is within that review's scope. Read the two
together, not the summary alone.

**A6. Optional reconsideration — measured at ~2% of a scale point.** See C6. Monteiro et al. 2015:
8% of diagnoses revised, 1.20→1.22.

**A7. ACH (analysis of competing hypotheses) — the one randomized trial with practitioners was
null-to-negative.**
Dhami, Belton & Mandel 2019, 50 practising UK intelligence analysts randomly assigned. Binary
accuracy **36% (n=9) vs 33% (n=8)**, χ²(1,N=49)=.04, **p=.845**. Rank-order accuracy 4% vs 4.9%, not
significant. ACH actively degraded discrimination: **80% of the ACH group produced tied ranks vs 19%
of controls**, χ²(1,N=41)=14.86, **p<.001, φ=.60**. ACH-trained analysts **did not follow ACH's own
steps**. Abstract: *"There was mixed evidence for ACH's ability to reduce confirmation bias, and we
observed that ACH may increase judgement inconsistency and error."*
Follow-up (Dhami et al. 2024, 10.1186/s41235-024-00560-y): the **matrix orientation is wrong** —
Study 1 (N=161) found hypotheses-in-**rows** significantly less prone to confirmation bias than the
ACH-style hypotheses-in-columns or prose (**p=.003**), and *"the ACH-style matrix did not confer any
benefits over the other two ways of structuring task information."* Study 2 (N=62 Dutch military
intelligence analysts) found most analysts **already** integrated evidence for and against each
hypothesis and were already sensitive to evidence credibility — challenging the confirmation-bias
premise ACH was built on.

**A8. Premortem — real effect, on the wrong dependent variable.**
Keysor, Wojtyna & Veinott 2020 (SJDM poster, opened): Exp 1 (N=53) premortem reduced **confidence**
more than promortem or control, F(2,51)=4.02, **p=.024, d=.81**; understanding only marginal
(p=.077). Exp 2 (N=43): group premortem **no better** than individual, F(1,42)=0.85, p=.772. **Every
controlled outcome in this literature is confidence, not decision change.** The premortem's own
principal evaluators state: *"Few controlled, randomized experiments have been conducted to evaluate
plan evaluation techniques and more are needed."* The frequently quoted "30% more reasons" figure
traces to Mitchell, Russo & Pennington 1989 (10.1002/bdm.3960020103), which **could not be opened** —
and even if true, its outcome variable is the **count of reasons generated**, not decision quality.

**A9. Temperature as a diversity lever — measured ineffective and mildly harmful.**
Zhu et al. (arXiv:2601.19921) Appendix C: t=1.2 produced **1.42 unique answers vs 1.45 at t=1.0** —
slightly fewer — while degrading instruction-following.

**A10. Over-generation as a coverage mechanism — measured to saturate hard.**
Si, Yang & Hashimoto (arXiv:2409.04109): **4,000** generated ideas per topic deduplicated to **~200**
(95% duplication), with the non-duplicate rate falling monotonically. Their own conclusion: *"This
sets a bottleneck on our inference-time scaling since increasing the number of generated ideas simply
leads to repeating duplicate ideas."* Corroborated by Hayati et al. (arXiv:2311.09799), which finds
saturation at ~7–8 distinct perspective clusters on subjective tasks, and only ~4–5 on hate-speech
labeling — even when asked for 20.

**A11. More rounds and more agents, in the hidden-profile regime — measured to *worsen* things.**
HiddenBench (arXiv:2505.11556): extended communication peaked at **15 rounds** then declined; **7
agents produced +0.6%** pre-to-post-discussion versus +34.8% for 3 agents. **No** prompting strategy
fixed the deficit: cooperative 20.0–24.2%, **conflictual 0–1.7%**, zero-shot CoT 22.2%, and even an
explicit "Share All Information" prompt reached only 46.7%.

**A12. Adversarial / conflictual framing — the worst-scoring condition in the one benchmark that
tested it.** HiddenBench's *conflictual* prompting: **0–1.7%**. Convergent from the human side:
Nemeth EJSP 2001 finds the authentic minority superior to **all three** devil's-advocate variants
(C-audit item 1). And Schweiger, Sandberg & Ragan 1986 (10.5465/255641 → note: 10.5465/255859)
reportedly bought decision *quality* at the price of **acceptance** — but that paper **could not be
opened** (AOM 403) and is **UNVERIFIED**; an open review that quotes it (Lunenburg 2012) gives the
**wrong page range** and misspells an author, which is itself a reason to distrust downstream
summaries. Schwenk's 1990 meta-analysis (10.1016/0749-5978(90)90051-A) also **could not be opened**;
its reported boundary condition — that dialectical inquiry's advantage does not hold on **ill-structured
tasks**, which is what protocol design is — would matter a great deal if true and should be retrieved
before anyone relies on it in either direction.

**A13. LLM-as-Judge overturning a majority — measured negative.**
Minority Sentinel (arXiv:2606.29270): the LLM-as-Judge baseline produced **negative Net Gain
(−1.37%)**, i.e. worse than not intervening, despite higher recall. Logistic regression managed only
+0.68%. The working version used cheap **behavioural signals from the debate log**, not a judge's
opinion. The authors: *"flip safety, not recovery volume, determines intervention value."* Their own
cost is reported honestly: 9 wrong flips overturning correct majorities, 7 of them on
knowledge-intensive MMLU-STEM.

**A14. "Fewer agents changed position" as a success metric — directly indicted.**
Hao et al. (arXiv:2606.00820) decompose the answer-flip rate into three mechanisms: spontaneous
instability, stance-induced conformity, and reasoning-induced persuasion. **37%** of agent-question
observations change under **self-reflection alone**, with no peer pressure at all, in the primary
MMLU-Pro setting. Strict conformity is 29%, predominantly harmful (57–77% correct-to-wrong across
model replications), and even **vacuous** reasoning drives 20–39% error adoption among otherwise
resistant agents. Decisively: *"without correctness labels or self-reflection controls, reducing peer
adoption does not improve accuracy, because harmful and beneficial influence cannot be distinguished."*
Any round-02 proposal whose success metric is "fewer agents changed position" or "less convergence"
has to answer this finding.

**A15. Examples in the prompt — measured to reduce diversity.**
Hayati et al. (arXiv:2311.09799): **five-shot prompting reduced diversity relative to one-shot**
("over-adherence to examples"), and zero-shot criteria prompting **underperformed** zero-shot
free-form (0.3176 vs 0.2885 on Social-Chem-101) — structure without a demonstration made things worse.
Relatedly, Yun et al. (arXiv:2505.18949) find structured chat templates suppress diversity and the
effect **persists under high-temperature sampling**, so it cannot be sampled away (Llama-3-8B
Distinct-2 0.1556→0.2107, Distinct-4 0.4699→0.5971 as structure is removed) — **but** removing
structure costs quality: *"while simple prompting improves generation diversity, it weakens the
model's ability to produce high-quality outputs."* This is a genuine two-sided tradeoff, not a
one-directional finding, and it bears on any protocol that mandates rigid artifact formats.

**A16. Assuming a defixation device transfers across participant classes — measured not to.**
Viswanathan & Linsey 2013 (10.1115/1.4024123) replicated Linsey 2010 with novices: *"both the novice
designers and design faculty fixate to the same extent, whereas the defixation materials have
differential effect on the two groups."* Their warning: *"The effectiveness of such measures may vary
with the level of expertise of the designer."* The **size of fixation** was identical across levels;
the **effectiveness of the cure** was not.

### What did have positive evidence, stated neutrally

Listed for symmetry, not as endorsement. Each is a *finding*, and each competes with others above.

- **Enumerated blacklist, not abstract instruction.** Chrysikou & Weisberg 2005 (n=89, three
  conditions): defixating instructions that **named** the problematic elements diminished fixation,
  where merely **describing** the flaws did not. George & Wiley 2020 (10.3758/s13421-019-01005-4)
  isolates which half does the work: a verbal **list** of common ideas *plus* a warning to avoid those
  specific ideas enhanced originality; the same examples **without** the avoid-instruction produced no
  benefit; and **visually depicted** examples produced fixation that "avoid" only partly repaired.
- **Incumbent removal.** Bilalić 2008's 1-solution condition, 0–18% → 100% (C4).
- **Forced re-representation.** McCaffrey 2012 (10.1177/0956797611429580): the generic-parts technique
  (decompose objects into parts with function-free names) yielded **67% more problems solved** than
  control. *Caveat: the 67% is the abstract's claim; the method detail sits behind a paywall and was
  not retrieved.*
- **Mandated structured re-analysis.** Mamede 2010, p=0.006 (C6).
- **A two-stage exchange-then-decide protocol, in LLMs.** HiddenBench's protocol — each agent must
  contribute 1–2 decision-relevant facts **and** one reason the currently favoured option might be
  wrong, in a separate stage from deciding — took GPT-4.1 from **3.7% → 80.0%** (+76.3pp),
  Gemini-2.5-Flash 17.3%→72.7%, Flash-Lite 4.3%→74.3%. This is the **largest measured effect in the
  entire sweep** and it is a structural rule, not a prompt. Two load-bearing details: the exchange was
  **separated** from the decision, and agents were **never told** information asymmetry existed.
- **Semantic-direction stratification.** arXiv:2605.30150: one planning call names ~5 broad semantic
  directions, generation is allocated across them. Best diversity-per-token in its comparison
  (0.295/0.379/0.207 gain per 100k tokens across GPT-5.4/Claude/Gemini, roughly double self-anchoring)
  at **1.6× token cost** versus 3.0–3.7× for self/peer/representative anchoring, and it **improves**
  quality rather than trading it away. *Scope: creative ideation quality proxies, not task accuracy.*
- **Verbalized Sampling.** arXiv:2510.01171: prompting for a **distribution** over responses with
  probabilities rather than a single answer yields **1.6–2.1×** diversity over direct prompting and
  recovers **66.8%** of the base model's pre-alignment diversity vs **23.8%** for direct prompting.
  Root cause identified as **typicality bias in preference data** — annotators prefer familiar text.
  Larger models benefit **1.5–2× more** than smaller ones.
- **Prompt perturbation.** DivSampling (arXiv:2502.11027): Random Idea Injection gave relative EM@10
  gains of 13.5% (MMLU-Pro), 15.5% (GSM-Hard), 15.4% (HumanEval). *Caveat: the paper asserts a
  diversity–fidelity tradeoff (Assumption 4.2) but presents **no** table where a perturbation
  underperformed baseline — the harmful regime is asserted, not demonstrated.*
- **Escalation with a named destination.** UFMCS Red Team Handbook v5 (2011), verbatim: *"if the staff
  dismisses an observation critical to mission accomplishment, the Red Team needs to inform the staff
  member that resolution is required with the Commander."* **DOCTRINE** — the handbook contains no
  empirical evaluation and says so: *"There are no formulas or simple checklists for Red Teaming...
  There is no simple formula or checklist that guarantees the insights promised by the red teaming
  concept."* It carries its own limits too: principle 7, *"Red Team recommendations must be within the
  ability of the command to implement"*, and the Red Team *"is not a shadow staff"* and *"must
  carefully weigh which items require elevation."*

---

## Engineering precedent

Eleven real processes were checked by fetching the templates, and where a corpus existed, by counting
it. The pattern: **every major process has an alternatives section, almost none make it mandatory, and
exactly one has written the enforcement code — then disabled it.**

| Process | Section | Mandatory? | Enforced by |
|---|---|---|---|
| MADR 4.0.0 | `## Considered Options` | Structurally required in the template (no optional marker) | **Nothing.** Repo CI is `check-links`, `lint`, `pages` + markdownlint only |
| Kubernetes KEP | `## Alternatives` | **Advisory** — not among the eight `(R)` required checklist items | Human SIG approvers; a TOC checker exists but is disabled |
| Rust RFC | `## Rationale and alternatives` + separate `## Prior art` | Advisory | Human, via FCP team sign-off |
| Python PEP | `Rejected Ideas` (PEP 1 item 10) | **Advisory** ("should") — contrast item 6, backwards compatibility, which "**must**" exist or the PEP "may be rejected outright" | Human PEP editors / Steering Council |
| Java JEP | `Alternatives` | **Explicitly optional** — the template marks only `Summary` and `Description` as `// REQUIRED` | Human, OpenJDK Lead |
| IETF RFC (7322) | **none** | **Does not exist** | n/a |
| IETF RFC 7282 | disposition rule, not a section | Normative practice (WG chair judgment) | Human chair; escalable via IESG appeal |
| NASA NPR 7123.1D §3.2.18 | Decision Analysis process | **"shall"**, requirement id **[SE-23]** | Human, at Engineering Technical Authority approval |
| NASA NPR 7123.1D Table G-3 | MCR success criterion 9 | Recommended best practice; deviation "should be justified to the ETA" | Human review board, checked as an **entrance** criterion |
| ISO/IEC/IEEE 42010:2022 | architecture rationale | **UNVERIFIED** — preview truncates before Clause 6 | unknown |

**Quoted text, where it is load-bearing.**

- **MADR**, full template: `## Considered Options` is a flat list of option titles, then
  `Chosen option: "{title of option 1}", because {justification}`. Optional sections carry an explicit
  inline marker (`<!-- This is an optional element. Feel free to remove. -->`); **Considered Options
  carries no such marker**, and it is one of only four headings that survive into the *minimal*
  template — "Pros and Cons of the Options" does not.
- **Kubernetes KEP** template: *"What other approaches did you consider, and why did you rule them
  out? These do not need to be as detailed as the proposal, but should include enough information to
  express the idea and why it was not acceptable."* Note it asks for **two** things per alternative:
  the idea **and** the reason it was rejected.
- **Rust RFC** template, four fixed questions under Rationale and alternatives, including
  *"What is the impact of not doing this?"* and *"If this is a language proposal, could this be done in
  a library or macro instead?"* Prior art is a **separate** section carrying an explicit
  anti-cargo-cult clause: *"precedent set by other languages [...] does not on its own motivate an
  RFC"*, and *"If there is no prior art, that is fine."*
- **JEP** template: *"Did you consider any alternative approaches or technologies? If so then please
  describe them here."* This is the **weakest wording of any process surveyed** — a yes/no question
  about the author's past mental state, answerable with "no" by an author who considered nothing.
  Compare Rust's and KEP's phrasings, which presuppose alternatives existed.
- **PEP 1** item 10, on why Rejected Ideas exists: *"This both helps record the thought process behind
  the final version of the PEP as well as preventing people from bringing up the same rejected idea
  again in subsequent discussions."*
- **RFC 7282 §3**, title verbatim: **"Rough consensus is achieved when all issues are addressed, but
  not necessarily accommodated."** Body: *"The group must truly consider and weigh an issue before the
  objection can be dismissed as being 'in the rough'. [...] the chair [...] is going to have to decide
  that not only has the working group taken the objection seriously, but that it has fully examined the
  ramifications of not making a change to accommodate it."* §2 supplies the elicitation form: *"a chair
  who asks, 'Can anyone not live with choice A?' is more likely to only hear from folks who think that
  choice A is impossible to engineer"*, followed by *"What are the reasons you object to choice A?"*
  §6 is titled "One hundred people for and five people against might not be rough consensus"; §7,
  "Five people for and one hundred people against might still be rough consensus."
- **RFC 7322 §4** — the authoritative required-section list for every RFC — contains **no**
  alternatives, rationale, design-considerations or rejected-ideas section. What the IETF *did* mandate
  is §4.8.5: *"All RFCs must contain a section that discusses the security considerations relevant to
  the specification."* One cross-cutting concern was judged worth a universal mandatory section;
  alternatives analysis was not. **This bounds how strong any "everyone does this" claim can be.**
- **NASA NPR 7123.1D §3.2.18.1**: *"Program/Project Managers **shall** identify and implement an
  ETA-approved Decision Analysis process [...] **[SE-23]**"*, with §3.2.18.2 naming the four ordered
  stages (criteria → alternatives → analysis → selection).
- **NASA Table G-3, MCR success criterion 9**: *"Alternative concepts have adequately considered the
  use of existing assets or products that could satisfy the mission or parts of the mission."* The
  matching **entrance** criteria require *"Alternative concepts that have been analyzed and are ready
  to be reviewed"* and *"Preliminary mission descope options"* as products that must exist **before the
  review can convene**. Tailoring is expected but *"The decision not to tailor and customize
  life-cycle review criteria should be justified to the ETA."*

### The measured compliance data — the most decision-relevant part of this section

**Kubernetes KEPs (n=657 non-template READMEs, counted).**

- Have an `Alternatives` heading at all: **482/657 (73%)** — **27% dropped the section entirely**.
- Heading present but body **empty** after stripping the template comment: **195/482 (40%)**.
- Word counts of the 287 genuinely filled sections: p10=6, p25=18, **median=51**, p75=118, p90=244.
- ≤20 words: 78/287 (27%). ≤50 words: 143/287 (50%).
- **Empty-or-≤20-words as a share of all 657 KEPs: 42%.**

**And the enforcement that would have caught every one of those, written and switched off.**
`kubernetes/enhancements/hack/verify-toc-vs-template.sh` diffs each KEP's TOC against the template's,
taking only deletion lines (`grep -E '^-[^-]'`) so a KEP may add headings but not drop them, and
ignoring headings marked `(Optional)` via `diff -U0 -I '(Optional)'`. **`## Alternatives` carries no
`(Optional)` marker, so a missing one would be flagged.** The script's last lines:

```bash
echo >&2 "Result: ${result}"
# TODO(soltysh): for now this should not fail, but print problems
exit 0
# exit "${result}"
```

The 42% figure is what "print problems, exit 0" buys.

**Python PEPs (n=737, counted).** Overall with a "Rejected Ideas" section: **181/737 (25%)**. But the
rate climbs monotonically with PEP number:

| PEP range | n | with Rejected Ideas |
|---|---|---|
| 0–299 | 120 | **0%** |
| 300–499 | 198 | 2% |
| 500–599 | 100 | 20% |
| 600–699 | 100 | 53% |
| 700–799 | 96 | 71% |
| 800–899 | 41 | **80%** |

A section that was never mandatory, never tooled, and only ever described as "should" went from unused
to near-universal — over roughly **500 PEPs and fifteen years** of consistent human editing.

**These two datasets support opposite conclusions and round 02 will have to weigh them.** KEPs say an
ungated template slot is ignored ~42% of the time **today**, by skilled engineers under public review.
PEPs say an ungated slot can reach ~80% compliance **eventually**, through review culture alone. Both
numbers are real and both were counted here.

**What was not verified.** ISO/IEC/IEEE 42010:2022's normative clause on architecture rationale.
The official free preview (15pp) was opened and confirms the standard's verbal forms and its Clause 4
conformance model — including the notable *"This document is designed such that 'tailoring' is neither
required nor permitted for its use when claims of conformance are made"* — but the preview **ends
mid-Clause 5.2.3, before Clause 6 where the AD requirements live**. The widely repeated formulation
(that an AD shall "provide evidence of the consideration of alternatives and the rationale for the
choices made") was found **only in secondary paraphrase**. Attempts at primary text failed:
iso-architecture.org refused connections, a mirrored 2011 PDF returned an empty body, dokumen.pub 403,
studylib 429. **Do not cite 42010 for an alternatives requirement until someone with ISO access reads
Clause 6.** This is exactly the §15 "untagged claim = RECALL" case: a paywalled clause that sounds
decisive and cannot be checked.

---

## Canonical deck measurements

**Script:** `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/protocol-generation-bias/reference/measure.sh`
Tooling: `/usr/bin/grep` (BSD grep) + `find` — **not** the shell's `grep`, which is ugrep and honours
`.gitignore`, so every zero below is safe against that trap. Two consecutive runs are byte-identical
apart from the timestamp. Generated 2026-08-28T11:54:37Z.

| Disputed quantity | Canonical value |
|---|---|
| Denominator | **89** idea directories; **88** have a `00-prompt.md` |
| `require_model_diversity` set as a frontmatter key | **1/88 including this idea; 0/87 excluding it**; plus 1 prose-only mention |
| `## Adversarial alternative` section | **3 ideas / 4 files**, of 89 ideas (3.4%) |
| `responding-to:` in later-round review artifacts | **63/349 (18.1%)** |
| `### @<other agent>` heading in later-round review artifacts | **25/349 (7.2%)** |

**Every definitional choice that made the earlier numbers disagree.**

1. **"Idea" = immediate subdirectory of `ideas/`** → 89. There are zero loose files. The entire 88/89
   gap is **one directory**: `launch-orphan-hardening`, which has `FINAL.md`, `IMPLEMENTATION.md` and
   `review/` but **no `00-prompt.md`**. Ruling: **89 for any directory question, 88 for any
   frontmatter-key question**, since a key can only exist in a file that exists. Both circulating
   numbers were right about different populations; neither said which.
2. **"Sets a key" = `^key:` inside frontmatter**, not "the string appears anywhere". Frontmatter is the
   lines between a `---` on line 1 and the next bare `---`; all 88 prompts have one. Column-0 anchoring
   excludes both indented keys nested under `roles:` and backticked prose references. **This single
   choice is the whole 0-vs-2 dispute.** `verification-honesty:27` reads "via opt-in
   \`require_model_diversity\`, escalate" — it *describes* the flag and gates nothing. SET and
   PROSE-ONLY are reported in separate columns and never added.
   - Reconciliation: **"0 of 88"** was correct for the deck *before* this idea existed, but used a
     denominator that already counted it. **"2 of 88"** = key-set (1) + prose-mention (1), counting a
     backticked prose reference as adoption.
3. **Self-inclusion reported both ways.** `protocol-generation-bias` (created 2026-08-28) is the only
   setter. **1/88 including it, 0/87 excluding it** — 0/87 and not 0/88, because removing the idea
   removes it from the denominator too.
4. **"Carries the section" = a literal `^## Adversarial alternative` heading; "mentions it" = the
   phrase anywhere.** 4 files / 3 ideas versus **15 files / 7 ideas**. The mention set is dominated by
   files *quoting the protocol rule* (`meta-protocol-change-verification-integrity/consensus.md` alone
   has five) and by `.codex.log` transcripts. **"15 slugs" was a file count relabelled as ideas.**
   Heading-level variants (`###`, bold) were checked and none exist, so the `##` anchor is not
   undercounting. The three ideas: `addon-bidding-scope`, `mas-research-mining` (twice),
   `zcode-adapter`.
5. **Review population = participant `*.md` artifacts only** in `review/round-NN`, NN≥2. Excludes
   `.codex.log` / `.hermes.log` / `.agy.log` transcripts and `DIFF-fixups.txt` (process residue), and
   excludes the one `VOID.md`, which self-declares "facilitator note, not a participant artifact".
   Denominator **349**.
6. **Guards were run, not assumed.** `responding-to:` counted in frontmatter and anywhere-in-file both
   give **63**, so the frontmatter restriction costs nothing here. `### @<other>` equals `### @<any>`
   at **25** — no agent ever addresses only itself.

**Two findings beyond what was asked, both bearing on round 02.**

- The two halves of the later-round reviewer rule **diverge sharply**: **63/349 (18.1%)** carry
  `responding-to:` but only **25/349 (7.2%)** carry a `### @agent` heading. They are not one rule in
  practice and should not be quoted as a single compliance figure.
- `track:` is set on **33** ideas but splits **16 deliberation / 16 standard / 1 fast**, which bears on
  any claim about conditional-rigor adoption.

**Adoption of the other opt-in gates, for context** (SET / prose-only, of 88):
`track:` 33/8 (37.5%) · `checks:` 4/10 (4.5%) · `strict_gate:` 2/6 (2.3%) · `auto_implement:` 1/3
(1.1%) · `require_model_diversity:` 1/1 (1.1%).

---

## What nobody has evidence for

Named open questions. Round 02 will have to decide these **without** empirical support, and should say
so in the artifact rather than reaching for a citation that does not cover the case.

1. **Whether human group-decision findings transfer to LLM agents at all.** Every result in C2, C3,
   C4, C5, C6, C7, C8 and most of the "evidence against" section is human. The transfer is plausible
   and nothing here refutes it — but **no study establishes it**, and Viswanathan & Linsey 2013 (A16)
   is a documented case of a debiasing device failing to transfer even between two classes of *humans*.
2. **Whether a heterogeneous frontier-model roster inherits the homogeneous-debate pathology.**
   arXiv:2605.00914 measures 10 identical 7–8B models and confines its own conclusion to that class;
   arXiv:2606.29270 uses three heterogeneous agents and finds a different picture. Neither settles the
   six-lab case, and the two together only **bound** it.
3. **Whether LLM agents systematically fail to propose the documented off-the-shelf option.** This is
   the deliberation's own motivating case and the literature does **not** answer it. The closest proxy
   found (arXiv:2512.11589, agent-authored PRs: libraries imported in 29.5% of PRs, a **new dependency
   added in only 1.3%**) is **UNVERIFIED** — the paper was not opened — and its own authors flag as
   future work *"whether agents sensibly reuse existing libraries or quietly reimplement
   functionality."* The 1.3% figure is equally consistent with agents correctly working within existing
   dependencies. **This question appears to be empirically open.**
4. **Whether an author-agent self-prefers when judging its own proposal in a deck round.**
   arXiv:2404.13076 measures single-model evaluation of candidate texts and establishes the
   self-recognition→self-preference link. Extending it to multi-agent protocol adjudication is a
   reasonable inference, not a measured result.
5. **Whether the sycophancy result transfers from user-directed to peer-directed.** arXiv:2310.13548
   measures sycophancy toward a **human user's** stated views. Both citing agents flagged the mapping
   as their own inference. Nothing measures LLM-to-LLM proposal anchoring directly.
6. **Whether debate participants re-converge after round 1.** kimi-1's second use of arXiv:2305.19118
   cites Degeneration-of-Thought as PRIMARY evidence for this. Liang et al. establish DoT for
   **single-agent self-reflection** and offer multi-agent debate as the **remedy**; they do not measure
   DoT inside a running debate. The claim may well be true; **this citation does not support it** and
   the inference should be retagged as reasoning.
7. **Whether structured conflict helps on ill-structured tasks.** Schwenk 1990's reported boundary
   condition — that dialectical inquiry's advantage does not hold on ill-structured problems, which is
   what protocol design is — is **UNVERIFIED** (paywalled, null abstract on both Crossref and Semantic
   Scholar). It is the only meta-analysis in the DA/DI literature and would be decision-relevant in
   either direction. One retrieval would settle it.
8. **Whether the Schweiger et al. quality-vs-acceptance tradeoff is real.** **UNVERIFIED** (AOM 403,
   no abstract exposed). If real, it points the *wrong* way for adding a challenge role. A published
   review citing it gets its page range and an author's name wrong, so downstream summaries are not
   trustworthy here.
9. **Whether the Team A/Team B cautionary history is as reported.** Both intelligence.senate.gov and
   cia.gov returned **403 to every route tried**. The reported pattern — a competing team selected for
   disposition produced a confidently wrong answer, with composition blamed by a 1978 Senate
   subcommittee and a 1989 internal CIA review concluding the threat was "substantially overestimated"
   — is **entirely UNVERIFIED**. It is the canonical example cited *for* institutionalised competing
   analysis, and the reported record cuts against it. Do not put it in a ratified artifact unverified.
10. **Whether prospective hindsight produces the "~30%" it is famous for.** Mitchell, Russo &
    Pennington 1989 **could not be opened** (abstract elided by the publisher on both Crossref and
    Semantic Scholar). An enormous amount of premortem advocacy rests on this single 1989 study, and
    the figure is almost always quoted at second hand.
11. **Whether any of this composes.** Every measured intervention in this file was tested **alone**.
    Nothing here measures stratification *plus* an exchange stage *plus* an artifact requirement, and
    DCI's ablations (C9) are the one attempt to decompose a composite protocol — and it found **no
    component earning its keep** at n=45. Round 02 has no evidence about interaction effects among the
    mechanisms it is choosing between.
12. **Whether structure at generation time costs what it buys.** Yun et al. (A15) measure structured
    templates suppressing diversity *and* removing them costing quality. The tradeoff is real in both
    directions and the crossover point for a deliberation protocol is unmeasured. This is a genuine
    design fork the evidence does not resolve.
