---
idea: mas-research-mining
author: claude-1
created: 2026-08-15
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## Problem / idea

The owner asks: **is there anything in current multi-agent-collaboration research — new ideas or
approaches — that would improve the Parley Deck protocol?**

A deep-research sweep has already run. Six independent lenses (deliberation/debate, failure modes,
orchestration frameworks, verification/judging, context/memory, governance/process) each produced
findings; each lens was then attacked by an independent adversarial verifier that fetched the primary
sources and checked the deck-side claims against this repository. ~48 findings went in; ~11 survived,
~15 were graded overstated-but-real, ~10 were dropped.

**Read the result first — it is the shared context for this idea:**

```
parley-deck/ideas/mas-research-mining/reference/research-brief.md
```

347 lines. It is deliberately **non-recommending**: it separates MECHANISM from CLAIM from EVIDENCE
STRENGTH and stops there. Deciding what, if anything, should change is this idea's work, not the
brief's.

**The unit of output is a ranked, small set of proposals.** For each: the mechanism, the *measured*
Parley problem it touches, what Parley already has, the concrete successor-idea shape, the cost of
being wrong, and — non-optional — **what evidence would show it did not work.**

**Returning few or zero proposals is an acceptable and possibly correct outcome.** The brief's §6
lists fifteen places where this protocol already appears to be ahead of the surveyed field, several
established by verifiers *trying and failing* to find a gap. "Confirmation, not opportunity" is a
valid finding. Do not manufacture parallels; this deck has been damaged by plausible-sounding
structure before.

## The five tensions this idea exists to resolve

Framed as questions, not as positions. Attack them in any order; disagree with the framing if it is
wrong.

**T1 — Measurement before mechanism?** Brief §9 says Parley has never measured: acted-on fraction of
findings per review round; DC-vs-DM classification of fix-up cycles; run-to-run reproducibility of
the same idea; outcome quality normalized by token spend; a compute-matched single-agent control; or
what the 1.6→5.1 review-round series decomposes into (agent error vs harness-forced re-run). One
source's sweep further implies **every historical Parley claim about whether more rounds helped is
compute-confounded — including the 1.6→5.1 series itself.** Does that make instrumentation the only
defensible first move, or is that an excuse to build nothing?

**T2 — Asymmetric context: designers vs reviewers.** Four unrelated lines converge on *give reviewers
less context, not designers* (Anthropic's verification-needs-minimal-context claim; Refute-or-Promote's
cold-start reviewers; ICML's teams-hold-experts-back; and Parley's own measured split — design rounds
flat at 1.6 where round-1 independence IS enforced, review rounds 1.6→5.1 where it is NOT). The
verifier flagged that no single line is strong and that the agreement is the signal. **But this sits
exactly where the deleted context selector sat.** Does the asymmetry survive the standard, or does it
die by it?

**T3 — The admission bar.** §15.4 (`COOPERATION.md:1321`) deliberately does not gate what a reviewer
may report; LE-1 requires a failing case **or** a check — a disjunction. **A CRITICAL asserted from
RECALL can open a full fix-up cycle today.** The brief notes this is the only mechanism in the entire
survey aimed at review *round count* rather than review *quality*. Is the ungated finding channel a
feature protecting genuine dissent, or the cost driver?

**T4 — Stopping on a signal that can lie.** VRR-Stop proves acceptance can rise monotonically while
true validity collapses (0.87 at round 2 → 0.12 at round 6; 55% of instances had a correct plan
repaired into an incorrect one). Two findings wrongly assumed Parley uses fixed-K; it does not
(`:646-666`, trajectory-based). But the accurate framing is sharper: **the paper refutes the heuristic
Parley actually runs** — `:657-659` defines converging as "total findings dropping sharply each pass",
which is precisely the acceptance signal shown to rise while validity falls. Separately, the
detection/miscorrection decomposition reports P(DM|D=1) of 53–94%: **a fix re-BLOCKed on the same
claim is an observed miscorrection that Parley never counts.**

**T5 — Rule count, collisions, and whether tightening the text even works.** Compliance degrades with
instruction *count* by **silent omission** (depth-20 follow rates 60.4% / 43.3% / 20.1%), and — the
negative half — rewriting the instruction set tighter recovers compliance only for weak models
(+11.0pp / +3.3pp / **−1.2pp**; ρ = −0.85, p = 0.004). Separately, WIRE finds within-policy rule
collisions get resolved **silently**, 64.6% non-jointly, with no signal that a conflict occurred; and
the largest deployed agent constitution handles this with **one document-wide priority order**, which
Parley lacks across its 18 sections (§15 has already had to declare pairwise subordination twice).
Note the trap: both results are extrapolation beyond their measured range, and one of them says our
instinctive fix does not work on strong models.

## What the protocol has ALREADY absorbed — do not re-propose these

All PRIMARY, from `COOPERATION.md` ratification backlinks:

| Section | Ratified by idea | Date |
| --- | --- | --- |
| §12 Pipeline blocks & action stages | `meta-protocol-change-end-to-end-pipeline` | 2026-06-02 |
| §13 Retrospective optimization | `meta-protocol-change-rho-retrospective-optimization` | 2026-06-16 |
| §13 amendment — confident-error signal | `meta-protocol-change-fusion-execplans` | 2026-06-18 |
| §13.5 Playbooks | `parley-learn-playbooks` | 2026-07-04 |
| §14 Automated outer loop (loop engineering) — the human brake | `automation-outer-loop` / loop-engineering ideas | 2026-07 |
| §15 Verification integrity | `meta-protocol-change-verification-integrity` | 2026-08-05 |
| §4.0 Conditional-rigor tracks | `meta-protocol-change-devx-speed` | 2026-07 |

Two of these came from prior *research* ideas of exactly this shape:

- **`loop-engineering-research`** already mined the June-2026 "loop engineering" paradigm (Cherny,
  Steinberger, Osmani) and it became §14. **Do not re-propose loop engineering.**
- **`meta-protocol-change-rho-retrospective-optimization`** already mined retrospective /
  self-improving process design and it became §13. **Do not re-propose retro passes.**
- **`cognee-mechanism-mining`** (closed 2026-08-14) evaluated knowledge-graph memory engines and
  adopted **nothing**. Its standing conclusions are binding here: a mechanism that decides what an
  agent sees must prove it never drops an objection; our generated-view pattern beats post-hoc
  dedup; the objection ledger in `protocol-read-cost-regression/FINAL.md` is ratified-but-unbuilt
  and is stronger than the external art we found.

If the research surfaces something that maps onto one of these, the correct output is
**"confirmation, already have it"** — which is a valid and useful finding — not a proposal.

## Open, ratified-but-unbuilt work this research must be judged against

Proposing something that duplicates unbuilt-but-ratified work is worse than proposing nothing,
because it splits authority. These are open:

- **`protocol-phase-scoped-packet`** — send only the protocol sections a phase needs. Ranked #1
  against the measured read cost. Must be built in the *instruction* layer: the Go runner never
  reads `COOPERATION.md` at all.
- **The objection ledger** (`protocol-read-cost-regression/FINAL.md:53-100`) — owner-namespaced IDs,
  exact scoped propositions, SHA-256 provenance, `OPEN|RESOLVED|DEFERRED|SUPERSEDED`, tombstones,
  forced `DISPUTED`. Implemented, then **deleted in v1.43.1**; the failure analysis is unwritten.
- **`meta-protocol-change-subtractive-maintenance`** — net-byte accounting ×3 lockstep copies plus a
  mechanical before/after rule inventory. Agreed 2026-08-14, not opened.
- **Per-idea protocol-version pinning** and the **deck overlay** — ratified, not implemented.

## Constraints

- **No new dependency, no new service, no tool in the normative path.** Output is protocol or
  tooling design, not an install.
- Anything that would change `COOPERATION.md` is a **§7 meta-protocol-change** proposed as a
  *successor idea*. Do not edit the protocol in this idea.
- **Provenance discipline (§15.2).** Every claim about external research is `SECONDARY` at best
  unless you personally read the source. Author self-description is weak evidence — a paper claiming
  a 30% improvement is a *claim*, not a *fact*, until you have seen the benchmark and the ablation.
  Every claim about Parley Deck needs a `PRIMARY` locator (file:line or a command) or a `RECALL` tag.
  **Do not assert a negative about this repository without running the check.**
- **The protocol is the measured cost centre** (105,382 bytes today, +112% in ten weeks, 3.3× read
  cost per call). Any proposal adding net bytes must name what it removes or why it pays for itself,
  and bytes land in **three lockstep copies** — count ×3.
- The measured explosion is in **review/fix-up rounds** (1.6 → 5.1, max 24), not in design rounds
  (1.4 → 1.6, flat). A proposal that adds rigor to the design phase is attacking the wrong term
  unless it argues otherwise explicitly.

## Non-goals

- A literature survey for its own sake. Findings that touch no measured Parley problem are noise.
- Adopting a framework, library, or platform.
- Editing `COOPERATION.md`.
- Re-proposing loop engineering, retrospective optimization, verification integrity, conditional
  rigor, or anything already ratified above.
