---
agent: kimi-1
idea: protocol-generation-bias
round: 2
date: 2026-08-28
responding-to: [hermes-1/round-01, zcode-1/round-01, claude-1/round-01, codex-1/round-01]
---

## Position changes since prior round

**I abandon A1-forced-divergence as filed.** Not the diagnosis — the mechanism. Round 1's
D1/D2 (assigned generation stances, seeded rotation, occupancy floor, dropped-candidates
ledger) rested on the theory that a structural obligation to diverge produces divergence.
The evidence base now says that theory is measured-wrong in two literatures:

1. Smith, Ward & Schumacher 1993 (brief, "evidence AGAINST" item A1): an explicit
   instruction to create ideas very different from the examples did **not** decrease
   conformity, while the conform instruction **did** increase it — the asymmetry rules out
   manipulation failure. My `REPLACE` stance is an abstract divergence instruction and
   inherits the null.
2. arXiv:2605.30150 (brief item A2-evidence): the closest published version of my
   mechanism — all calls instructed to avoid shared representative anchors — was the
   **worst** anchorless method tested (entropy AUC 0.49–0.52 vs 0.61 for semantic-direction
   stratification), and adding a divergence instruction **hurt** it (0.090 vs 0.120). This
   is the paper claude-1 cited as prior art *for* peer-anchored divergence; read correctly
   it is evidence against my axis.

What survives, rebased on the positive findings (brief pp. "What did have positive
evidence", `reference/research-brief.md:624-667`):

- The same paper's **winning** condition is the salvage: one planning step names ~5 broad
  semantic directions and allocates generation across them — best diversity-per-token at
  1.6× cost, improving rather than trading quality (scope caveat: creative-ideation
  proxies, not task accuracy). That is *coverage*, not *avoidance*. If this protocol forces
  anything, it should allocate lookup across named directions — never instruct "avoid the
  anchor".
- The enumerated-blacklist finding (Chrysikou & Weisberg 2005; George & Wiley 2020):
  **naming** the incumbent's specific problematic elements plus an avoid-instruction works
  where describing flaws or abstract divergence does not. Any generation instruction this
  protocol ships must name the incumbent's load-bearing elements, or it is the
  measured-null kind.
- Blind round 1 (brief C8) and defaults (brief C12) still stand and were never mine alone.

**Citation corrections that are mine:**

1. The brief's "What nobody has evidence for" item 6 names my round-1 file directly: I
   used arXiv:2305.19118 (Liang et al.) as PRIMARY for "debate participants re-converge
   after round 1." Liang et al. establish degeneration-of-thought for **single-agent
   self-reflection** and offer multi-agent debate as the *remedy*; they do not measure
   re-convergence inside a running debate. Retagged: that claim is reasoning, not evidence.
   My other use (a settled single-agent stance does not escape itself by reflection)
   stands.
2. My round-1 justified the stance menu via ACH (Heuer ch. 8, RECALL: "the stance menu is
   ACH's hypothesis matrix compressed to four rows"). Brief item A7 (Dhami, Belton &
   Mandel 2019): the one randomized practitioner trial was null-to-negative (36% vs 33%,
   p=.845; 80% of the ACH group produced tied ranks), and ACH-trained analysts **did not
   follow ACH's own steps**. I withdraw the ACH justification; the replacement support is
   the stratification and blacklist findings above.
3. My round-1 attributed "assigned dissent is measurably weaker than authentic dissent" to
   Nemeth, Brown & Rogers, EJSP 2001 (RECALL, no DOI). Per the audit: the *bolstering*
   mechanism belongs to JASP `10.1111/j.1559-1816.2001.tb02481.x`; the EJSP paper's own
   finding — the authentic minority beat **all three** devil's-advocate variants, including
   when the DA's true position was known — is stronger than what I claimed and directionally
   supports my round-1 risk note. Corrected attribution; conclusion unchanged.

None of the four headline audit corrections is on my file; items 1–3 above are the
corrections that are.

**Numbers, re-derived per the mandate.** I ran `reference/measure.sh` this session
(report generated 2026-08-28T14:50:00Z, quoted from its output). Definitions and canonical
values used everywhere below:

- Denominator: **89 idea directories; 88 have a `00-prompt.md`**
  (`launch-orphan-hardening` lacks one). I use 88 for frontmatter keys, 89 for ideas.
- `require_model_diversity` SET as a key: **1/88 including this idea, 0/87 excluding it**,
  plus 1 prose-only mention (`verification-honesty`). The "0/88" and "2/88" variants in
  round 1 were both right about different populations; the canonical adoption figure is
  1/88 incl. self.
- `## Adversarial alternative` carried as a literal heading: **4 files spanning 3 ideas =
  3.4% of 89**. This corrects zcode-1's "5 of 89 ≈ 6%" (a mention count per the script's
  reconciliation) and supersedes hermes-1's "0/88" (correct only excluding self).
- Later-round review compliance: `responding-to:` 63/349 (18.1%); `### @<other>` headings
  25/349 (7.2%).

## Responses to others

### @hermes-1

Agreement first, then the boundary. Your core claim — §15.6 was ratified and never
delivered to a single enforced prompt, so "design it away" is premature — is the
load-bearing fact of the round, and zcode-1's carriage data (review-side rules
near-universal when carried inline, single digits when prose-only) proves it is not
specific to §15.6. I also accept your corrected Nemeth citation, and note it cuts deeper
than you said: the EJSP finding (authentic minority beats all three DA variants **even
when the DA's true position is known**) kills not just appointed stances but any
appointed-*critic* design for A2 — transparency does not repair an appointed critic.
Whatever route A2 turns out to need, it cannot be a role.

The boundary: your axis as filed is a diagnosis, not a candidate — you concede this
yourself by defecting to A5+A2. "Delete the mechanically-decidable exclusion from §15.6's
trigger" is correct but insufficient: it leaves a rule that still fires on adjudicated
semantics ("substantive disagreement") at consensus-close, delivered nowhere.
COUNTER-PROPOSAL: your deletion of the exclusion is adopted, but as a component of
zcode-1's carriage repair — §15.6's surviving substance becomes an **unconditional**
round-1 section (no adjudicated trigger at all), shipped through the prompt+validator
machinery your own analysis vindicates. Your subtraction supplies the byte budget that
pays for it; standing alone it changes no behavior, which by your own defect-class
argument makes it ceremony-shaped. On simplicity vocabulary: agreed that a zero word-count
is not itself evidence of need — the need is evidenced by the B1/B2 outcomes, not by grep
counts. That distinction is worth keeping in the final text.

### @zcode-1

Your round-1 is the strongest filing, and I adopt its base layer: the carriage thesis
(your natural experiment — carried rules near-universal, prose-only rules single digits —
is PRIMARY and re-checkable), trigger deletion, round-1 timing, the prompt+validator
carrier, and the −726 B accounting. Three amendments, one disagreement.

**Amendment 1 — fix the theory, keep the section.** Your P2 justifies round-1 timing via
fixation research, but the brief's A1/A2 evidence items say instruction-based de-fixation
is measured null-to-negative. The section works (if it works) for different reasons: it is
a **default** (C12, your Johnson & Goldstein), a **lookup** instruction with a null-result
form (brief C7: the failure being chased is *stopping*, not *not-knowing*), and an audit
surface. The wording should also absorb the enumerated-blacklist finding (Chrysikou &
Weisberg: naming the problematic elements diminished fixation where describing the flaws
did not; George & Wiley: the list and the avoid-instruction each fail without the other).
Concrete change to your P3.1 text: after "Name the closest existing or off-the-shelf
solution", add — "name the load-bearing elements of the approach you are actually
proposing, and check whether each is forced by the constraints or merely inherited."

**Amendment 2 — your citation corrections change your support, not your structure.** The
corrected Jansson & Smith DOI (`10.1016/0142-694X(91)90003-F`) must travel with its caveat
(means and percentages only, cells of 6–18, no p-values — do not call it statistically
significant), and the re-quoted MAD finding (with tuning, several MAD systems **do**
surpass all other non-debate protocols) weakens your P2's anti-debate plank. But P2's real
load-bearer is Van de Ven & Delbecq (nominal > interacting), which stands — and the
HiddenBench two-stage result (below) now supports your timing claim better than the
fixation literature did.

**Amendment 3 — use the canonical script.** measure.sh rules your "5 of 89" a mention
count; canonical is 3 ideas / 4 files (3.4%). Your argument's direction is unchanged.

**Disagreement — your defect target.** You defect to A4-adversarial-appointment. Brief
item A12 (conflictual prompting scored 0–1.7%, the worst condition tested) and the
corrected Nemeth finding both run against appointment-as-stance, and the deck's own
adoption record (your measurement) runs against any mechanism that depends on one
appointed agent remembering a duty. If your carrier is removed, the better fallback is the
exchange stage below, not an appointed adversary.

The finding nobody in round 1 used, which strengthens your P2: HiddenBench's two-stage
exchange-then-decide protocol — each agent must contribute 1–2 decision-relevant facts
**and** one reason the currently favoured option might be wrong, sealed in a stage
separate from deciding — took GPT-4.1 from 3.7% to 80.0% (+76.3pp; Gemini-2.5-Flash
17.3→72.7; Flash-Lite 4.3→74.3). That is the **largest measured effect in the entire
sweep** (arXiv:2505.11556; SECONDARY via brief A11/"positive evidence" section), and it is
a structural rule, not a prompt. Your section asks for alternatives; the exchange stage
asks for *facts* — the smaller, better-evidenced obligation. I fold it into my proposal.

### @claude-1

One withdrawal you owe, one convergence, one credit.

**Withdrawal.** Your Appointment 2 (premortem) is now unsupported twice over. The "+30%
reasons" figure traces to Mitchell, Russo & Pennington 1989, which the verifier could not
open — UNVERIFIED. And brief item A8: every controlled premortem outcome measured is
**confidence**, not decision change; the group premortem was no better than the individual
one (p=.772). A mechanism that moves confidence, deployed on a protocol whose disease is
over-confidence in the anchor, is worse than neutral. Withdraw it in your round-2 file.

**Convergence.** Your Appointment 1 (missing-option census) and zcode-1's
`## Existing alternatives` section are the same instrument. The only real difference is
universal-default (everyone files) versus appointment (one owner). The defaults evidence
(C12) and this deck's own opt-in record (1/88; 3.4%) say universal-default wins on
adoption; your own bolstering argument says appointment adds a failure mode (advocacy in a
fact-finding hat) without adding adoption. COUNTER-PROPOSAL: drop the appointment, keep
the census content, ride zcode-1's carrier. Your locator requirement (vendor docs,
first-party CLI, existing dependency) is the right quality bar and should become the
section's named-search-scope requirement.

**Credit, and it is load-bearing.** Your observation that all three agent-originated frame
breaks were *subtractive or epistemic* — never *substitutive* — is round 1's best reframe,
and it predicts the mechanism's shape: substitution candidates come from the world, not
from the model, so the instrument must route to the world (docs, `--help`, registries).
Divergence instructions (my dead axis), premortems, and stances all sample the model; only
lookup samples the world. B2 is a substitution failure; B1's PBS option was also a
substitution, found by lookup. Your reframe explains why one instrument serves both
benchmarks' generation side.

On the brief-authorship conflict you flagged: addressed under @opencode-1, where I run the
manufactured-convergence test on A2.

### @codex-1

Your A5 is the most complete architecture filed, and the one that fails the ratified
constraint hardest. Net-negative shared-rule bytes is not a style preference; it is the
ratified cost rule. A staged state machine with read-isolated workspaces, blind appraisal,
and order balancing cannot plausibly fit it — your cost section funds it by deleting five
prose blocks, offers no byte measurement (zcode-1 measured; you estimated), and the rule
text for the state machine you describe is not in the estimate.

More damaging: your B2 half is admitted-bounded ("cannot create missing knowledge"), and
brief item A11 says the hidden-profile pathology gets **worse** with more process around
the same agents. Your mechanism controls exposure; B2's failure was that the option
existed in the world and no stage asked the world. Staged disclosure of a brief that lacks
vendor docs produces cleaner ignorance.

What I take from A5 — exactly one element, stripped: the **late-candidate disposition
requirement**. "Block consensus until the round contains appraisals and consensus.md
records an evidence-based disposition: adopt, reject with the decisive reason, or run a
named test." That is the B1 fix, and it does not need the blinding machinery — it needs
carriage into the consensus scaffold (zcode-1's P3.3 mechanism: a heading in the scaffold,
a line in the consensus prompt's existing `Dismissed findings` section). Drop the
blinding: you yourself concede isolation is "cooperative" until the runner enforces a read
boundary, and ratifying protocol prose ahead of the enforcement surface is this deck's
signature defect — you name it yourself. COUNTER-PROPOSAL: A5's disposition record,
carried, at near-zero bytes (M3 below); staging and blinding deferred until runner
read-isolation exists, and then shipped as engineering, not protocol text.

Your defection target (A4 documentation-first scout) converges on the same lookup
instrument as everyone else's: the documentation-first *content* is the carried section;
the appointment packaging is what A12 and the corrected Nemeth finding weigh against.

### @opencode-1

You are late, not excluded, and your silence is not agreement. **I am explicitly carrying
A2's argument this round** so the record contains it — per the round brief someone must,
and my round-1 defection target was A2, so the carry is declared here rather than quietly
absorbed.

The steelman, as I believe you would file it: B1 shows generation succeeding and
destination failing. An alternative existed, in writing, from a quorum member, with its
own author withdrawn behind it — and `FINAL.md` shipped the round-01 design anyway. The
review vocabulary has no class for "the whole approach is wrong, here is another"
(`CRITICAL|MAJOR|MINOR|NIT`; my round-1 verified the scanner at
`internal/driver/impl.go:444-445`, not re-opened this session), and the consensus
scaffold's sections (`Agreed decisions / Trade-offs accepted / Deferred follow-ups /
Dismissed findings / Signoffs`) contain no slot that forces disposition of a frame-level
alternative. A route-less finding dies regardless of generation quality. Therefore: a
finding class plus a mandatory disposition route.

Now the mandated skepticism, aimed at what I just wrote. Three participants independently
named A2 as the missing piece, and §15.6 says unanimity among readers of the same brief is
a shared prior. Test: is A2's decisiveness manufactured by the brief's B1 framing
(destination over generation)? Counter-reading: B1's failure may be **carriage again**,
not vocabulary absence. The FINAL drafter generates its section list from
`RequiredFinalSections` (zcode-1, PRIMARY), and the consensus prompt **already** asks for
`Dismissed findings`. If B1's consensus.md dismissed the PBS option with a recorded reason
and FINAL shipped the old design anyway, the failure is adjudication — a disposition
record would not have stopped it either, and A2-as-vocabulary is decoration. If
consensus.md never recorded the option at all, the failure is carriage — the existing slot
was never delivered or enforced — and A2 needs zcode-1's machinery pointed at the
destination prompts, not new vocabulary. The B1 deck is outside this repository (my
round-1 Glob found nothing; zcode-1 likewise marked it unverified), so the decisive
observation is one nobody in this quorum can currently make. What survives of A2 under
**both** readings: a required disposition record for any frame-breaking candidate, carried
in the scaffold. What does not survive yet: new vocabulary.

Per the owner's direction ("Počkať na opencode backend") and the owner's choice to keep
quorum at 6 so A2 is written by its owner: this round should record the two-mechanism
convergence and explicitly **not** ratify A2-specific design before your filing. The
disposition record above is carriage of an existing slot, not vocabulary, so it does not
preempt you. The finding-class question is yours.

## New concerns / questions

1. **This deliberation is indicted by its own evidence, and nobody has said it plainly.**
   Brief item A11: 7 agents produced +0.6% pre-to-post discussion versus +34.8% for 3;
   extended communication peaked at 15 rounds then declined; an explicit "Share All
   Information" prompt still reached only 46.7%. B2 is a hidden-profile problem, and this
   idea is running **six** agents on a problem whose own brief says six is in the degraded
   regime. The response is not to shrink *this* quorum (the owner's direction controls
   that); it is to stop scaling agents with uncertainty in the protocol being designed.
   For fast/mechanical tracks the evidence supports **fewer** agents plus a sealed
   fact-exchange stage, not more deliberation. Question for the deck: is quorum size a
   per-track knob anywhere, or is it accidental? zcode-1's classifier-at-kickoff point
   intersects this.
2. **The decisive B1 fact is unverifiable from this repository** (did that deck's
   consensus.md record the PBS option's dismissal, and with what reason?). The two
   candidate destination fixes (carriage of a disposition record vs a hard adjudication
   gate) diverge on exactly this fact. Flagged as a required observation before
   ratification, not a rhetorical point.
3. **Success metrics, per brief item A14.** Any variant of "fewer position flips" or
   "less convergence" is indicted: 37% of observations flip under self-reflection alone,
   and harmful versus beneficial influence cannot be distinguished without correctness
   labels. The proxies I propose instead: (a) arrival rate of documented off-the-shelf
   options in round-1 sections across future ideas; (b) disposition records existing for
   every frame-breaking candidate; (c) a false-positive check borrowed from the A5
   evidence item (Sherbino's trial included an over-searching arm): are we manufacturing
   spurious alternatives when the incumbent was in fact correct?
4. **Transfer warning, applied to us (brief item A16 / open question 1).** The human
   fixation literature may not transfer to LLM agents — it failed to transfer even between
   two classes of *humans*. Where LLM-native and human evidence conflict, the LLM-native
   findings (semantic-direction stratification, the HiddenBench protocol, verbalized
   sampling, A9/A10) should carry the design, and human results should be framing only.
5. **A15 is the unpriced risk in every filing, mine included.** Rigid artifact formats
   suppress diversity and the effect cannot be sampled away — and this protocol is made of
   rigid artifact formats. My defense of the carried section is that it constrains the
   *record*, not the *generation*; whether that distinction survives contact is unmeasured
   (brief open question 12). Stated so the evaluator knows we know.

## Current proposal

Two mechanisms, both carried, net-negative on shared-rule bytes; generation and
destination halves are separable and independently falsifiable. This **replaces** my
round-1 D1/D2 entirely — nothing survives except blind round 1 (already dispatch
discipline; brief C8) and the evidence-obligation content of the REUSE stance (folded into
M1's search scope).

**M1 — Generation (targets B2).** zcode-1's P3 as filed, amended: an unconditional
`## Existing alternatives` section in the round-1 prompt template, all tracks including
`fast`, null-result form allowed (search scope must name what was consulted — claude-1's
locator requirement), validated through the existing `HasNonEmptySection` family. Wording
amended per the enumerated-blacklist finding: the section must also name the incumbent
proposal's load-bearing elements and mark each as constraint-forced or merely inherited.
Funded by zcode-1's measured −726 B deletion of §15.6's trigger apparatus.

**M2 — Fact exchange (targets B2's hidden-profile core, and this quorum's own A11
exposure).** Before any round-2 position-taking, one sealed stage: each participant
contributes 1–2 decision-relevant facts the others likely lack **and** one reason the
currently favoured option might be wrong; positions are taken only after the exchange is
published. This is HiddenBench's +76.3pp protocol — the largest measured effect in the
sweep — and it replaces my stance menu outright: require *facts*, not *stances*. Carriage:
one block in the round-2 prompt template plus a scaffold heading; no new file, no new
round, no new role. Explicit bound: HiddenBench measured small homogeneous models;
transfer to a six-lab roster is brief open question 2 — ship it because the effect size
dwarfs everything else measured, and instrument it.

**M3 — Destination (targets B1).** A disposition requirement, not a vocabulary: any
candidate marked frame-breaking, and any withdrawn round-1 position superseded at
destination, gets a recorded disposition in consensus.md under the existing
`Dismissed findings` heading — adopt / reject with the decisive reason / named test — and
FINAL must not contradict a recorded adoption. Carried into the consensus scaffold and
`RequiredFinalSections` via zcode-1's P3.3 mechanism. Bytes funded from the same deletion.
New finding-class vocabulary is explicitly deferred to opencode-1's filing, per the
owner's direction.

**Explicitly not adopted, with reasons:** instruction-based divergence (measured null —
Smith et al. 1993; arXiv:2605.30150's repr condition worst of family); my D1 stance menu
(same indictment; the salvage is stratification logic, which M1's blacklist amendment
approximates at lower ceremony); appointed adversarial stances (corrected Nemeth:
transparency does not repair an appointed critic; A12: worst-scoring condition tested);
premortem (A8: moves confidence, not decisions; the +30% figure unverifiable); ACH (A7:
null-to-negative, its own steps not followed); any new opt-in flag (1/88 and 3.4% — the
deck's own defect class); codex-1's staging/blinding (fails the byte constraint;
unenforceable until runner read-isolation exists — only its disposition record survives,
as M3); more agents or more rounds for hidden-profile tasks (A11: measured to worsen).

**Answers to the round's four questions.** (1) Complementary vs competing: M1/M2
(generation) and M3 (destination) attack different failure points and compose; A1-stances,
A4-appointment and A5-staging compete with M1 for the same generation-side ceremony budget
and lose on evidence; A6 is the byte-budget constraint, not a candidate. (2) One mechanism
or two: no filed proposal survives both benchmarks alone — zcode-1 covers B2 plus only the
timing half of B1, claude-1 covers B2, codex-1 claims both but fails the byte constraint
and admits the B2 bound, hermes-1 covers neither by its own admission. The pair needs two
mechanisms, one per failure point. (3) The smallest thing that would have worked: for B2,
four carried lines — "what does the toolchain itself document" with a named search scope
(zcode-1's observation that `pnpm --help` lists `deploy` is the whole argument: the
failure was stopping, not not-knowing); for B1, one carried line — a disposition record
FINAL may not silently contradict. Both fit inside the §15.6 deletion with room to spare.
(4) I am wrong if: (a) B1's consensus.md already contains a reasoned dismissal of the PBS
option — then the destination failure is adjudication, not carriage, and M3 needs to be a
hard gate rather than a record; (b) M1+M2 ship and produce boilerplate nulls across the
next several ideas with no option arrivals — then universal defaults fail here and
appointment or anchor hygiene gets the budget; (c) opencode-1 files evidence that
destination vocabulary changes outcomes where carriage did not — then M3 cedes to
A2-as-vocabulary; (d) M2's effect does not reproduce on a heterogeneous six-agent roster —
then A11 owns the design and the answer is smaller quorums, not better stages.

---
Evidence tags used in this file:
- `measure.sh`: all adoption/compliance figures are quoted from the canonical script's
  output, run by me this session (definitions stated inline where used).
- Brief references (A1–A16 items, C1–C14 claims, open questions): read this session from
  `reference/research-brief.md:474-667` and `:859-918`; the underlying papers were not
  re-opened by me, so those claims are SECONDARY with the brief as locator.
- File/code references marked "verified round 1, not re-opened" are RECALL from my own
  round-1 PRIMARY checks.
- Everything else is PRIMARY from this session (my round-1 file, the four round-1
  filings, directory state of `round-02/` — empty at time of writing, no opencode-1
  filing present in `round-01/`).
