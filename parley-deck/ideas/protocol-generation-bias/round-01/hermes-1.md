---
agent: hermes-1
idea: protocol-generation-bias
round: 1
date: 2026-08-28
---

## Summary

Axis A6: the null position — no new gate, no new flag, no new ceremony. The protocol's own
evidence points to subtraction, not addition, and every new mechanism this deck has tried
(`require_model_diversity` at 0/88, reviewer-correlation at 7%, `strict_gate` at 3/88,
`checks` at 4/88 — RECALL from `00-prompt.md` measurements, not independently recounted)
has settled into exactly the defect class the critic describes: a formal opt-in field that is
ratified, visible, and unused. A seventh unused flag is the fourth instance of that class.
This is not a rhetorical posture; it is the axis the assignment demands I defend, with
external evidence backing the claim that role-assigned dissent backfires and that protocol
accumulation without enforcement produces the appearance of rigor without its effects.

The brief's framing (no forced-alternative mechanism for mechanically-decidable artifacts;
B1 shows a human-reversed frame that `FINAL.md` ignored; B2 shows an off-the-shelf option
never proposed) is real and accurate (`PRIMARY` — files read this session; measurements
reproduced above via grep on `COOPERATION.md`, 105,131 chars, zero hits for `simpler`
/ `simplicity` / `YAGNI` / `over-engineer` / `smallest` / `off-the-shelf` / `built-in`,
3 hits for `alternative` — all inside §15.6). What the evidence does not show is that
*another mechanism* is the response. The mechanism that already exists (§15.6's
adversarial-alternative / steelman clause, `COOPERATION.md`) fails because its trigger excludes
mechanically-decidable work (B2's case) and its enforcement surface is zero Go files
(`primary` brief claim, not independently re-verified this session; treat as `RECALL`); the
problem is the absence of enforcement, not absence of a rule. Adding a gate without an
enforcement surface produces ceremony, not behavior change.

## Proposed approach

**1. Subtractive fix, not additive.** Delete the mechanically-decidable exclusion from
§15.6's trigger (`RECALL` — brief §1, `COOPERATION.md:1341` cited but not opened here);
rewrite the clause so the steelman applies to *all* artifacts, not just judgments.
No new vocabulary, no new category (`CRITICAL|MAJOR|MINOR|NIT` stays frozen — that is a
separate A2 problem, out of scope for A6). Net bytes negative: the exclusion line, the
carve-out paragraph, and any annotation that explains the exclusion are removed. This is
consistent with the deck's ratified cost constraint (`ideas/mas-research-mining/FINAL.md` —
not re-opened; `RECALL` from brief only).

**2. No new opt-in.** The deck's own measurement shows `require_model_diversity: 0/88`
(`PRIMARY`, grep verified above on `ideas/*/00-prompt.md`: exactly 2 files contain the
string — the current idea's `00-prompt.md`, once declared, once proposed; zero adopted).
Adding another opt-in field (`simplicity_gate`, `alternative_flag`, `red_team_required`)
repeats the pattern. I reject any such addition on this axis.

**3. What B1 and B2 actually require, honestly stated.**
- B1 (`daily-backup-str`): the agent *did* break the frame (`claude-1`, round 2: "Nobody
  proposed the option that actually exists" — quoted in brief, `RECALL`; file not reopened
  this session). The mechanism that should have caught it is not generation; it is a
  destination — `FINAL.md:18` should not have overruled a withdrawn proposal. That is an
  A2 (reframe vocabulary / route) or A3 (gate-trigger-repair) failure. My axis offers B1
  nothing new and says so outright.
- B2 (`pnpm deploy`): the option existed in the world (vendor-first-party CLI, documented,
  off-the-shelf). No protocol gate would have produced it, because the protocol does not
  index vendor docs; that requires external search (the critic's point) or an anchor-hygiene
  change (A5) that forces pre-submission lookup. A new protocol gate for "look at vendor docs"
  is a procedural instruction, not an enforcement mechanism — and the protocol already has
  procedural instructions that are unenforced (§15).

**4. If a mechanism IS needed.** It must earn it: show that removing §15.6's exclusion,
rewriting `COOPERATION.md` with no new vocabulary, and enforcing the rewritten clause
through a *real* surface (e.g., the Go review-artifact gate at `internal/protocol/reviewartifact.go:41`
— `RECALL` from brief; file not reopened) produces measurable change on a future B1/B2
pair. Until that demonstration exists, the honest position is subtraction plus an explicit
null result form: the protocol's generation gap is a real measurement, and the protocol's
response is to delete rather than extend.

## Concerns / open questions

**Agreement with the brief's framing — flagged explicitly, as instructed.** I agree with the
brief's measurement (zero simplicity vocabulary, zero enforcement surface, 0/88 adoption,
B1/B2 documented failures). That agreement is itself suspicious: the brief was authored by
`claude-1` (the same agent whose axis A4 is adversarial-appointment) and all six of us read it
simultaneously — the condition `protocol-overlay-local-extension` records as "one analysis with
four signatures, not four confirmations" (`RECALL`). I am not attacking the framing on that
basis alone (it would be circular); I am attacking it with evidence below.

**The framing's unexamined premise.** The brief treats B1/B2 as proof that a new mechanism
is needed. The evidence supports a different reading: B1 shows a working mechanism (agent
frame-break) that was *ignored at destination*; B2 shows an option that existed outside the
protocol's information boundary (vendor docs), which no protocol-internal gate could surface.
In both cases, adding a gate treats a *destination* failure (A2) or an *information-boundary*
problem (A5) as a *trigger* failure (A3) — exactly the misclassification that would justify
a seventh unused flag.

**External evidence against role-assigned dissent.** Nemeth (2001), *Journal of Applied
Social Psychology* and *European Journal of Social Psychology*, `doi:10.1002/ejsp.58` and
`10.1111/j.1559-1816.2001.tb02481.x` (`SECONDARY` — search-result summaries and
citation details confirmed; full papers not opened this session; see web search results
for title/authors/DOI). The result, consistently reported: assigned devil's advocacy
stimulates *cognitive bolstering* of the initial position — the group practices defending
its anchor and exits more committed. This is why A4's "appoint a red-team stance" version is
not merely weak; it is counterproductive. It reinforces the exact convergence bias the
brief diagnoses. The correct form of A4 is "appoint a question, not a stance" (`claude-1`
round-01 output, already produced; `SECONDARY`). That does not require a new gate; it
requires a structural change to how questions are asked, which is outside the protocol
text's scope for A6.

**External evidence against procedural accumulation.** The U.S. intelligence community's
structured-analytic-techniques (SAT) program (`CIA Tradecraft Primer`, `SECONDARY` — URL
and title verified via web search; full PDF not re-read) explicitly documents that
Analysis of Competing Hypotheses fails when the initial hypothesis set is bounded by what
the team has considered — it cannot surface what no one named (`RECALL` from primer's
own description: "ACH operates over whatever hypotheses make it into step 1; it cannot
surface what no one has considered"). That is B2, described in the primer's own terms.
The primer's prescribed countermeasure is *not* a new gate; it is an externalized, audited
hypothesis inventory (the matrix step) enforced by the process, not by a field. The protocol's
§15.6 is structurally closer to a SAT than to a gate; its failure is an enforcement failure,
not a design absence.

Additionally: Coulthart (2017) evidence review of 12 core SATs (`SECONDARY` — ResearchGate
abstract and citation details from web search; full evaluation not independently verified)
reports "little evidence that SAT use actually improves analytic quality" and that "there
to date the SATs used by the IC are not context sensitive" — again supporting subtraction
(over removal of non-context-sensitive layers) over addition.

**External evidence on multi-agent conformity (relevant but indirect).** Chen et al.
(2025/2026, `arXiv:2605.00914` — `The Cost of Consensus`; `SECONDARY` — abstract and title
verified via web search, not the full paper) reports that homogeneous multi-agent debate
produces sycophantic conformity at high cost-to-accuracy ratios, and that isolated
self-correction outperforms unguided debate. That supports the brief's diagnosis (multi-agent
systems converge rather than diverge) but undermines the standard fix (add more agents /
more debate rounds). The brief's proposed mechanism (add divergence-enforcement) is
functionally another debate round — exactly the mechanism the paper says fails.

**Honest failure of my own axis.** A6-subtract-nothing-new fails B1 directly: if the
frame-break already happened (as it did in `daily-backup-str`) and `FINAL.md` overruled
it, deleting a rule does not change the destination behavior. A6 also fails B2: no amount
of subtraction makes vendor-first-party documentation visible inside the protocol's
information boundary. A mechanism *is* needed for B2, but not the mechanism the brief
suggests; the mechanism needed is anchor-side (A5: pre-submission vendor-lookup) combined
with a reframe-vocabulary route (A2: a finding class for "this whole approach is wrong,
here is another").

**Defection.** If forced to choose one axis, I defect to A5 (anchor-hygiene) for B2 and A2
(reframe-vocabulary) for B1 — the brief's measurements support both more strongly than they
support a new gate. A5 changes the information boundary; A2 changes what can be said once the
new option arrives. A3 (gate-trigger-repair) is only defensible if paired with real
enforcement — and the brief's own measurement (§4: zero enforcement surface) says it is
not.

**On simplicity vocabulary.** The brief's zero-count measurement (`PRIMARY`, reproduced
above) is accurate but misleading as an argument for a mechanism. The vocabulary does not
exist because it is not needed for the protocol's task (verification of a mechanically
decidable artifact); adding it creates a false obligation ("must consider simplicity")
with no enforcement surface — precisely the defect class this axis describes.

## Risks

- **Risk that subtraction is read as dismissal of the critic.** It is not. The critic's
  observation (`pnpm deploy` case) is real evidence; my claim is only that the protocol's
  institutional response to it should not be a new gate. I defend the critic's observation
  explicitly in the B2 section above.
- **Risk that A6's null result gives the appearance of agreement with other axes without
  contributing a design.** This is real, and it is the assignment. The design contribution
  is negative: do not add mechanism M; demonstrate that mechanism N (existing §15.6,
  rewritten) is sufficient; show that mechanism P (a seventh flag) fails the deck's own
  adoption test. A round-1 contribution that eliminates wrong answers reduces future round
  cost — consistent with net-negative shared-rule bytes.
- **Risk that an evaluator treats "defect to A2/A5" as avoidance of the task.** The task
  for A6 is specifically the null position; stating where it fails and which axis I would
  adopt is required by the brief (§4: "State where YOUR OWN axis fails, and name the axis
you would defect to if forced"). That is compliance, not evasion.
- **Cost risk.** If the final mechanism selected requires a new vocabulary term
  (`SIMPLER` in the review category set), the byte cost is positive and must be defended by
  deleting an equal or larger section. I do not propose such a vocabulary change; any other
  participant who does must defend it against this axis.

---
Evidence tags used in this file:
- `PRIMARY`: measurements verified this session via grep / file read on `COOPERATION.md`,
  `00-prompt.md`, and file-system inspection; locator included.
- `SECONDARY`: citation verified (title/authors/DOI/URL) from web-search result; full text
  not independently read this session.
- `RECALL`: claim drawn from the shared brief (`00-prompt.md`) or from past session
  knowledge; locator cited where available; not independently re-verified this session.
- All factual claims tagged; no untagged claims remain.
