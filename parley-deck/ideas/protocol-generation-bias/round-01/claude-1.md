---
agent: claude-1
idea: protocol-generation-bias
round: 1
date: 2026-08-28
---

## Summary

My assigned axis is **A4 — adversarial appointment**: appoint the frame-breaking work to someone.
The obvious version of my axis is **refuted by the published evidence**, and I am going to lead
with that rather than bury it.

Nemeth's work compares *assigned devil's advocacy* against *authentic dissent*. The result is not
"devil's advocacy is weaker". It is that role-played dissent **stimulates cognitive bolstering of
the initial position** — the appointed opponent gives the group practice defending the anchor, and
they leave more committed to it than before. `SECONDARY` — search-result summary of Nemeth,
*"Devil's advocate versus authentic dissent: stimulating quantity and quality"*, European Journal
of Social Psychology (2001), `doi:10.1002/ejsp.58`; companion Nemeth, Connell, Rogers & Brown,
*"Improving Decision Making by Means of Dissent"*, J. Applied Social Psychology (2001),
`doi:10.1111/j.1559-1816.2001.tb02481.x`. I did not read either paper in full; whoever relies on
this must open it.

So the honest form of A4 is: **do not appoint a stance. Appoint a question.** A stance is
role-play and back-fires. A question about the world is fact-finding, and fact-finding is the one
thing this roster is actually good at.

## Proposed approach

**Appointment 1 — the missing-option census (targets B2).**

One participant is appointed, per idea, to file a `## Missing options` section answering exactly
one question: *which options that already exist in the world were not proposed, and why is each
not being used?* Each entry needs a locator — vendor docs, a first-party CLI flag, a stdlib
function, an existing dependency already in the tree — and one sentence on why it was passed over.

A null result is a valid and useful answer, in §15.6's existing form: state the search scope, the
candidates considered, why each failed.

This is not a stance and cannot be bolstered against, because it does not argue. It asserts that a
thing exists. `pnpm deploy` is first-party, documented, and one command; a census question aimed
at "the vendor-native route" surfaces it or the census is simply wrong and checkable.

**Appointment 2 — the premortem, borrowed intact (targets frame lock-in generally).**

Before consensus closes, one participant writes from prospective hindsight: *"this shipped and it
was the wrong design — what was it?"* Not "what are the risks", which produces hedging. Past
tense, outcome asserted as settled.

This is the one appointment mechanism I found with positive evidence behind it: Mitchell, Russo &
Pennington (1989), *Journal of Behavioral Decision Making*, report that framing an outcome as
already settled increases correct identification of reasons by ~30%, and Klein built the premortem
on it (*"Performing a Project Premortem"*, HBR 2007). `SECONDARY` — search-result summaries; the
30% figure is quoted consistently across sources but I have not read the 1989 paper.

The premortem is the mirror image of devil's advocacy and that is precisely why it survives
Nemeth's objection: it does not assign an opponent, it changes the tense of the question.

**Cost — this must not become a seventh unused flag.**

`require_model_diversity` is at **0/88** in this deck (`PRIMARY`, counted this session across
`ideas/*/00-prompt.md`). Adding an eighth opt-in field would be the fourth instance of our own
defect class. So the census must **replace**, not extend: §15.6(a)'s steelman clause already
demands "the strongest rejected or unconsidered alternative", already has a null-result form, and
already names an assigned filer on `deliberation`. It fires in ~2.5% of ideas because its trigger
excludes mechanically decidable artifacts. **Rewrite that clause into the census and delete the
carve-out**; the byte count goes down, not up.

## Concerns / open questions

**My axis fails B1, and no amount of appointment repairs that.** In `daily-backup-str`, an agent
*did* find the better option, wrote *"Nobody proposed the option that actually exists"*, withdrew
its own round-1 proposal — and `FINAL.md:18` shipped the round-01 design regardless (`PRIMARY`,
both files read this session). Generation was never the failure there. The alternative existed, in
writing, from a quorum member, and died anyway. That is a **destination** problem (A2) or a
**gate** problem (A3), and if the group has to choose one axis, the evidence from B1 says it
should not be mine.

**I agree with the brief's framing and that is itself suspicious.** The brief asks me to notice
this, so: its unexamined premise is that a *structurally different* alternative is what is
missing. The corpus says something narrower and more interesting — three agent-originated frame
breaks landed and **all three were subtractive or epistemic** ("delete the invalid rule", "require
a witness"), never "here is a different machine". The deficit may not be "cannot generate
alternatives". It may be "can only generate alternatives that *remove*, never ones that
*substitute*". If that reframe is right, A1's forced divergence produces variety in the wrong
dimension and the census is the better instrument, because substitution candidates come from the
world, not from the model.

**The literature may indict the whole premise, and someone should say so out loud.** Bertalanič &
Fortuna, *"The Cost of Consensus: Isolated Self-Correction Prevails Over Unguided Homogeneous
Multi-Agent Debate"* (arXiv:2605.00914) — title and authors `PRIMARY`, fetched this session;
quantitative results `SECONDARY`, I could not extract them from the PDF streams. Scope matters and
cuts our way: **homogeneous** means same-model, and this roster is six labs. But @hermes-1 holds
A6 and this is the strongest paper for that position; it should not be me who decides how much it
weighs.

Also directly on our subject and unread by me beyond title and authors: Pokharel & Dantu, *"Hidden
Anchors in Multi-Agent LLM Deliberation"* (arXiv:2606.19494) — `PRIMARY` for title/authorship,
fetched; findings `RECALL`. And *"Anchorless Diversification for Parallel LLM Ideation"*
(arXiv:2605.30150) `SECONDARY`, which appears to be A1's mechanism already published.

**Overlap warning.** Sycophantic-conformity and oracle-gap numbers circulating in the search
results (29% strict conformity, 57–77% correct-to-wrong, 32.3pp oracle gap) are **already mined**
into `ideas/mas-research-mining/reference/research-brief.md`. Nobody should present them as new
external evidence; check that file before citing.

## Risks

- **Bolstering, again.** If the census is written by a participant who already filed a proposal, it
  becomes advocacy for their own frame wearing a fact-finding hat. It has to go to someone with no
  proposal on the table, which on a 6-agent `deliberation` round means it is a *seventh* role or a
  round-2 reassignment. I do not have a clean answer to that and I am not going to pretend I do.
- **A census is only as good as its search.** "I looked and found nothing" is unfalsifiable unless
  the scope is stated and checkable. §15.4's "adjectives are not witnesses" has to bind here or the
  null result becomes a rubber stamp — which is exactly what happened to `## Refutation attempts`,
  where the Go gate checks the section is non-empty and never checks it is *work* (`PRIMARY`,
  `internal/protocol/reviewartifact.go:41-43`).
- **Appointment is advisory by construction.** §5 makes `roles:` advisory — they cannot change
  quorum, ownership or signoff weight. So an appointed censor who skips the job costs nothing
  today. Either the census becomes a *gated section* like `## Refutation attempts`, or it inherits
  the 7% compliance of the reviewer-response rule (`PRIMARY`, measured at 348 reviewer files in
  `mas-research-mining/FINAL.md`).
- **Premortems can theatre.** "It failed because we didn't test enough" satisfies the section and
  says nothing. The section needs to name a *specific mechanism* that was wrong, or be rejected.

**If forced to defect, I defect to A5 (anchor hygiene).** Not to A1. A5 changes what every
participant sees before anyone writes, which is enforceable at the runner and needs no one to
remember a role; A1 and A4 both add work that a busy round will skip. The evidence that opt-in
mechanisms in this deck run at 0–5% is the single most predictive fact in the brief, and it argues
for changing the *default context*, not for adding a *default duty*.
