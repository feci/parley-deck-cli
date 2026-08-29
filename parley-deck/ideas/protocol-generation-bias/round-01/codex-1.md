---
agent: codex-1
idea: protocol-generation-bias
round: 1
date: 2026-08-28
---

## Summary

I propose a default-on **Anchor Firewall**: constrain what may enter the round-1 brief, disclose context in stages, seal each participant in a read-isolated workspace until submission, and blind both authorship and presentation order during the first appraisal of candidates. A later frame-breaking candidate must re-enter that blind lane before consensus, so the mechanism protects an alternative that appears after round 1 rather than merely generating a cleaner initial slate. **PRIMARY — local check.** I ran `sed -n '181,360p' parley-deck/COOPERATION.md` and `sed -n '901,1080p' parley-deck/COOPERATION.md`; the relevant output says round 1 is written “without reading other agents' round-1 files,” while §11.B calls the default “the same social one” and makes stronger sub-branches optional. Anchor hygiene should therefore be a runner-enforced disclosure property, not another instruction to resist influence.

## Proposed approach

### External evidence and the design consequence

- **PRIMARY — verified this session.** Tversky and Kahneman, “Judgment under Uncertainty: Heuristics and Biases,” *Science* 185(4157), 1974, DOI [10.1126/science.185.4157.1124](https://doi.org/10.1126/science.185.4157.1124), accessible [paper PDF](https://fbaum.unc.edu/teaching/articles/Science_1973_JudgmentUnderUncertainty.pdf), write: “The initial value, or starting point, may be suggested by the formulation of the problem” and “different starting points yield different estimates.” The proposal treats the kickoff formulation itself as an experimental input that must be controlled.

- **PRIMARY — verified this session.** Stasser and Titus, “Pooling of Unshared Information in Group Decision Making,” *Journal of Personality and Social Psychology* 48(6), 1985, DOI [10.1037/0022-3514.48.6.1467](https://doi.org/10.1037/0022-3514.48.6.1467), accessible [paper PDF](https://www.uni-muenster.de/imperia/md/content/psyifp/aeechterhoff/vorlesungkommunikation/stasser_titus_unsharedinfogroupdisc_jpsp1985.pdf), report that discussion “tended to perpetuate, not to correct, members' distorted pictures of the candidates.” The protocol consequence is to preserve a uniquely held option through an independent appraisal before ordinary consensus discussion can make the shared frame more salient.

- **PRIMARY — verified this session.** Rowe and Wright, “The Delphi technique as a forecasting tool: issues and analysis,” *International Journal of Forecasting* 15(4), 1999, DOI [10.1016/S0169-2070(99)00018-7](https://doi.org/10.1016/S0169-2070(99)00018-7), accessible [paper PDF](https://cpb-us-e1.wpmucdn.com/sites.psu.edu/dist/7/144284/files/2021/12/The_Delphi_technique_as_a_forecasting_to-2.pdf), identify “anonymity, iteration, controlled feedback, and the statistical aggregation of group response” as defining features; they explain that anonymity is intended to avoid pressure from “dominant or dogmatic individuals, or from a majority.” I borrow anonymity and controlled feedback, not Delphi's statistical aggregation: evidence, not vote count, still decides.

- **PRIMARY — verified this session.** Wang et al., “Large Language Models are not Fair Evaluators,” ACL 2024, DOI [10.18653/v1/2024.acl-long.511](https://doi.org/10.18653/v1/2024.acl-long.511), [ACL Anthology paper](https://aclanthology.org/2024.acl-long.511/), find that candidate ranking “can be easily hacked by simply altering their order of appearance in the context.” Their balanced-position calibration evaluates candidates in swapped positions. The protocol consequence is pointwise appraisal first and counterbalanced order whenever a decision depends on comparative ranking.

- **PRIMARY — verified this session; preprint, not claimed as peer reviewed.** Qu, Fu, and Hu, “Easier to Mislead Than to Correct: Harmful and Beneficial Revision in LLM Conformity,” arXiv:[2606.01637](https://arxiv.org/abs/2606.01637), report a controlled study across four open-weight models and seven QA datasets; their abstract says, “Authority labels make models more likely to choose the endorsed answer, regardless of whether it is correct.” They also report an asymmetry in which wrong peer agreement caused more harmful revision than correct peer agreement caused repair. The protocol consequence is to hide author, model, role, status, and apparent-majority cues until each participant has recorded an evidence-bearing appraisal.

**PRIMARY-based inference.** These studies do not prove that the exact workflow below will recover the best software design. They do establish relevant failure channels in human discussion, Delphi feedback, LLM evaluation, and simulated LLM peer influence; the design therefore controls formulation, identity, majority cues, feedback timing, and order separately instead of asking a model to “be independent.”

### 1. Make `00-prompt.md` an admissibility-limited container

Keep two explicitly marked participant-content sections in the one existing kickoff artifact:

1. `## Core brief` may contain only the requested outcome or question, hard constraints, observable acceptance criteria, non-goals, material unknowns, and the minimum raw inputs required to understand the task.
2. `## Evidence bundle — disclose after frame receipt` may contain source text, logs, code observations, benchmark records, and necessary history. Every factual item carries its provenance and locator. The bundle may state what happened; it must not tell participants what the event proves.

Coordination frontmatter may retain `idea`, `author`, `track`, and `participants` for the canonical audit trail, but the runner strips author and participant identity from the round-1 model view. The following material is inadmissible before every participant has sealed its own frame receipt:

- facilitator, critic, prior-human, or prior-agent diagnoses and recommendations;
- candidate solutions, mechanism names, or worked solutions unless the user explicitly mandates one as a hard constraint;
- rankings, vote counts, “obvious/best/simplest” labels, predicted consensus, or statements about what other agents believe;
- round-1 `roles:` or lens assignments, model/provider names, status labels, and submission order;
- selected corpus counts or examples offered to establish a causal frame rather than as raw input; and
- any partial participant output.

The kickoff author puts its diagnosis and candidate in its own round-1 file, on the same footing as every other participant. A role or lens may be assigned only after all round-1 submissions are sealed, and it is launcher metadata for cross-review rather than a retroactive addition to `00-prompt.md`.

**RECALL — supplied brief, not independently verified as a repository artifact.** Under this rule, the present brief's eight interpreted findings, six named solution axes, PDS mechanism description, and critic-supplied `pnpm deploy` contrast would not enter the core brief. Raw records needed to test B1 and B2 could enter the evidence bundle, without the labels “missing option,” “obvious,” or “frame reversal.”

### 2. Use one staged state machine, without a new canonical artifact

```text
CORE -> FRAME RECEIPT -> EVIDENCE -> ROUND-1 SEAL
     -> BLIND APPRAISAL -> AUTHORSHIP REVEAL -> OPEN DISCUSSION
```

- **CORE.** The runner supplies only a rendered `Core brief`, with coordination identities removed.
- **FRAME RECEIPT.** Each participant privately writes a short problem decomposition, unknowns, and documentation/search plan. The runner stores a hash and holds the text as scratch state.
- **EVIDENCE.** Only after that participant's receipt is fixed does the runner disclose the same evidence bundle. Evidence items are placed in a seeded cyclic order so no item is systematically first across the roster. A participant may revise its position, but its final round-1 file must reproduce the original receipt verbatim under `## Pre-evidence frame` before its proposal.
- **ROUND-1 SEAL.** Each complete file is held in its isolated stage and hashed. No file, progress count, timestamp, or partial text becomes visible until every expected receipt exists or the existing deadline rule is invoked. The runner then publishes the complete batch.
- **BLIND APPRAISAL.** Before ordinary round 2, the runner removes frontmatter and self-identifying metadata, assigns opaque candidate IDs, and gives each reader a balanced cyclic ordering. Each candidate is first assessed alone against acceptance criteria: strongest case, decisive failure test, missing evidence, and disposition if that evidence holds. If pairwise ranking matters, the reader sees both `A/B` and `B/A`; an order-dependent flip is recorded as unresolved rather than silently averaged away.
- **AUTHORSHIP REVEAL.** After all blind appraisals are sealed, the runner reveals the ID mapping and canonical files. Each participant's normal round-2 artifact includes its unchanged blind appraisal, followed by the currently required named responses. No extra round, signoff, or protocol file is created.

Questions raised during the core stage receive one neutral answer broadcast to every isolated participant at the same disclosure boundary. This prevents a private clarification from becoming an accidental private anchor.

### 3. Enforce isolation at the read boundary

The runner must launch each round-1 participant with a read allowlist containing only the currently disclosed rendered brief, an immutable task/code snapshot, allowed external-research tools, and that participant's scratch state. It must not mount the deck root, `.git`, other staging directories, pull-request state, or other participants' outputs. The only writable target is private scratch; the runner, not the participant process, relays the sealed result to the participant's canonical path.

Do not call a branch, worktree, or instruction “blinding” unless the process read policy actually prevents access to the hidden material. If the runtime cannot enforce the read boundary, it records `anchor-isolation: cooperative` and the round cannot be represented as blinded. Ratification should be coupled to runner enforcement; adding protocol prose first would recreate the deck's unenforced-rule problem.

### 4. Reset the lane for a late frame-breaking candidate

Any participant may mark a concrete later-round candidate `FRAME-BREAK`. A valid marker includes the candidate, the mechanism-family difference, a source or reproducible witness, and which acceptance criterion could change; a slogan or objection without a candidate is malformed. The runner deduplicates equivalent markers and performs a **late-candidate reset**:

1. Extract a neutral candidate card containing the proposal, evidence, constraints, and claimed observable effect.
2. Present it pointwise to every participant with only the core brief. Withhold its author, model, round, the incumbent drafter's position, apparent support count, and prior disposition.
3. Seal each independent appraisal before revealing the candidate's provenance and comparing it with the incumbent in both presentation orders.
4. Block consensus until the next existing round contains those appraisals and `consensus.md` records an evidence-based disposition: adopt, reject with the decisive reason, or run a named test. This is not a vote and does not force adoption.

This reset is the part that makes A5 more than a round-1 generation aid. A candidate that arrives late receives the same clean first hearing as an initial candidate; the accumulated narrative cannot be included in that first hearing as evidence of quality.

### Benchmark B1 — an option exists and the old frame still wins

**RECALL — benchmark supplied by the brief, not independently verified.** In B1, a participant reportedly introduced a native S3-capable option in round 2 and withdrew its own round-1 position, yet the final artifact retained the earlier `vzdump` plus `rclone` design.

The native option would qualify as `FRAME-BREAK`. Every participant would first appraise its evidence and fit against the original outcome without seeing that it was late, who proposed it, that the incumbent already had plurality/commitment, or how the drafter had framed it. Consensus could not close until every appraisal existed and the disposition explained why the native option was adopted, rejected, or sent to a concrete test. A5 still would not guarantee the native option wins; it would make “the existing frame won by inertia” an inadmissible close path.

### Benchmark B2 — the cheap option never appears

**RECALL — benchmark supplied by the brief, not independently verified.** In B2, `pnpm deploy` reportedly remained unproposed until a human named it, while agents elaborated hand-written dependency-copying scripts.

The core brief would state the deployment outcome, repository/package constraints, and acceptance tests, but no hand-written script frame and no `pnpm deploy` answer. Before receiving history, each isolated participant would seal a search plan; the evidence stage could disclose the package manager, manifests, and authoritative documentation locators without supplying a preferred mechanism. A discovered built-in route would therefore arrive without first being cast as an amendment to a script.

**RECALL — bounded hypothesis.** This improves the option's opportunity to be found but cannot ensure that any participant retrieves or recognizes it. If every isolated participant lacks the relevant knowledge or searches the wrong documentation, A5 reproduces B2 with cleaner provenance. That is my axis's principal failure.

### Shared-rule-text and operational cost

The protocol patch for A5 must be rejected unless authoritative shared-rule-text bytes are lower after the change. Fund the compact Anchor Firewall by deleting, not supplementing:

- the duplicated social-independence prose in the Quickstart, §4.0 invariants, Phase 1, §11.A, and §11.B;
- the round-1 `roles:` schema plus repeated role/lens explanations, because round-1 roles are forbidden here; and
- §15.6's current unanimity-and-judgment-triggered body, replacing it with a short cross-reference to blind appraisal, late-candidate reset, and the required `no independent alternative surfaced` record.

Measure the authoritative core before and after with `wc -c`; generated deck views and packaged copies mirror the same rule and do not create a second allowance. If another axis is combined with A5, it spends from the same deletion budget rather than treating A5's savings as free bytes.

Operationally, the design adds one short pre-evidence completion and one blind-appraisal completion per participant, plus runner-side rendering, sealing, and order balancing. It adds no canonical artifact, no quorum member, no signoff, and no mandatory extra round; the two scratch responses are copied into already-required round files. The cost applies only when work has already crossed the threshold for Parley; trivial reversible work remains outside this protocol.

## Concerns / open questions

- **PRIMARY — session record.** This brief anchored me. It preselected A5 and supplied the vocabulary “blinding,” “presentation order,” “staged disclosure,” and “isolated staging,” so this file cannot serve as evidence that I independently discovered the Anchor Firewall. At most, it develops the assigned axis and attacks the brief's own disclosure design.

- **RECALL — alternative framing, intentionally unresolved.** The critic's cases may be caused by retrieval/tool failure, missing domain expertise, reward for elaboration, or a close/adjudication defect rather than anchoring. A5 controls exposure; it does not diagnose which cause dominated either benchmark.

- **PRIMARY-based limitation.** Tversky–Kahneman and Stasser–Titus studied people, Wang et al. studied LLM comparison rather than proposal generation, and Qu et al. used simulated peer answers on QA tasks. Their results motivate controls but do not establish transfer effect sizes for software-design deliberation.

- **RECALL — implementation risk.** Authorship blinding can leak through style, unique citations, specialist knowledge, or self-reference. The protocol should promise removal of explicit cues, not perfect anonymity, and reveal authorship after the first appraisal so ownership remains auditable.

- The evidence allowlist needs a rule for tasks whose starting point is an existing implementation. My answer is to disclose raw code, logs, and required history in the evidence stage while excluding the facilitator's diagnosis; whether that is sufficient should be tested on real ideas before global adoption.

- A user may deliberately mandate a mechanism. That mechanism is then a hard constraint, not an anchor to hide. The rendered brief must distinguish “use X” from “X is an example,” and ask the user when the language is genuinely ambiguous.

- A frivolous `FRAME-BREAK` could stall closure. Require a concrete alternative, differentiating mechanism, witness, and affected acceptance criterion; deduplicate identical candidates, but do not give the facilitator discretion to suppress a well-formed marker.

- **RECALL — explicit defection.** If forced to abandon A5, I would defect to **A4-adversarial-appointment**, specifically a named documentation-first simplicity scout. B2 shows the missing duty: hygiene can give a built-in option a clean hearing, but only an accountable search role makes someone look for it.

## Risks

- **False neutrality.** An allowlist can hide interpretation but cannot remove the kickoff author's selection of facts. The pre-evidence frame and symmetric clarification broadcast expose some of that influence; they do not eliminate it.

- **Prompt starvation.** A core brief that is too thin wastes participants' time or produces irrelevant frames. The evidence stage must release all decision-material facts after the receipt, not use “blinding” as a reason to conceal necessary context.

- **Security theater.** Shared filesystem credentials or unrestricted repository reads defeat staged disclosure. Until the runner has an enforceable read allowlist, the mechanism is cooperative and must be labeled that way.

- **Cost and latency.** Two short hidden-stage completions per participant are material on a six-agent deliberation. The absence of new rounds/artifacts and deletion of the current §15.6 ceremony are the compensation, but measured run data should decide whether the staged calls remain universal or become mechanically triggered by the presence of an evidence bundle.

- **Late-candidate denial of service.** Repeated cosmetic reframes could keep resetting appraisal. Structural eligibility and content-hash deduplication reduce this risk; a hard numeric cap would recreate a route for suppressing the one important late option, so exhaustion should escalate rather than auto-close.

- **No discovery guarantee.** A5 can prevent a supplied or discovered alternative from being born underneath another agent's frame, and can prevent a late alternative from being dismissed under accumulated social cues. It cannot create missing knowledge. Any final design claiming to solve B2 needs A4's accountable search duty or another axis with an equally explicit discovery mechanism.
