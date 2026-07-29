---
agent: claude-1
idea: parley-design-skills
round: 2
date: 2026-07-28
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

**I withdraw my six-phase `PDP/1` ladder.** codex-1 and kimi-1 independently made the same
argument and it is correct: a design protocol with its own phase state machine is a second
state machine, which is anti-goal #3 — an anti-goal I endorsed in round-01 and then
violated two paragraphs later by proposing P0–P5. That is exactly the kind of drift the
adversarial round exists to catch, and I would rather be caught now than in review.

What I got wrong specifically: I treated "protocol-shaped" as "has its own phases". It does
not. AG-UI is protocol-shaped because of **typed artifacts, normative language, a version,
an extension policy, and a conformance definition** — none of which requires owning a phase
ladder. `parley-design` can have all five while riding Parley's existing rounds.

Second change: **I withdraw the vendored-fallback rule registry.** I proposed the checker
carry its own copy if `parley-design` is not installed. That is a second hand-maintained
copy, which is precisely the AG-UI `events.proto` rot (16 types vs 33) that the research
names. Fail closed instead.

Third change: **I no longer support any numeric aesthetic score**, including the anchored
weighted rubric in the research strawman. See @kimi-1 below.

## Responses to others

### @codex-1

**Adopting — the profile position, wholesale.** "Reuse Parley Deck's files and phases rather
than create a second state machine" is the single most important sentence in round-01. Your
mapping (round-01 = isolated directions, round-02 = one adversarial critique, consensus =
the Decider's recorded choice, `FINAL.md` = the binding contract) is the right skeleton and
I have adopted it below.

**Adopting — the descriptive-system limit.** *"This descriptive system cannot retroactively
legalize a contract violation"* closes the exact hole I flagged as my own open question #1
in round-01. I asked what happens when the post-build `DESIGN-SYSTEM.md` contradicts the
ratified contract; your answer is better than my "the build wins". The build wins as a
*description*; it does not win as an *authorisation*. A drift between contract and built
system is a finding, not a retroactive amendment.

**Rejecting — five files and 96 KiB.** The number is defended by "phase read sets are
fixed", but a fixed read set is a discipline, not a size argument. Four agents × 96 KiB on
every full run is real cost for content that is mostly a web annex. I hold at four files
and propose a byte ceiling below.

**Rejecting — mandatory deterministic rank-assignment on every full route.** Your SHA-256
roll is checkable, which is its virtue, but it solves a problem we do not have by default.
impeccable needed dice because it had *one* model with a deterministic resonance ranking.
We have four model families. Spending the mechanism unconditionally means it must be
verified unconditionally. Make it the **remedy for a failed distinctness gate**, not a
precondition — same determinism, paid only when needed. Counter-proposal in §E.

### @hermes-1

**Adopting — `enforced-by:` on every rule.** `enforced-by: check:<id> | agent-judgement` is
the best single idea in round-01 for the split line. One annotation simultaneously stops the
doctrine claiming mechanical enforcement it does not have and stops the checker grading
taste it cannot measure. It is now the centre of my §F answer.

**Adopting — DTCG `2025.10` verbatim, and WCAG blocking / APCA advisory.** Concrete,
standards-anchored, and it means our token layer is checkable by tooling we did not write.

**Rejecting — the D0–D7 phase spine.** Same objection as to my own round-01 proposal. Your
phases D3 Decide / D4 Graft / D5 Tokenize are not phases, they are *steps inside* Parley's
consensus and finalization. Keep the **names** as a vocabulary — they are good names and
agents can cite them — but a name is not a state machine. Concretely: your D4 and D5 cannot
have independent exit conditions without inventing artifacts that duplicate `consensus.md`
and `FINAL.md`.

**Rejecting — `rules.json` plus a generated markdown catalog plus a drift test.** You cite
the right cautionary tale and then reintroduce two representations. A drift test only fires
*after* someone forgets to regenerate, and only if the test is run. The literate-file
alternative has no second copy to drift. See §B — this is my main disagreement with you and
I would rather resolve it now than ship two artifacts.

**Gap.** Your file does not label positions F1–F8 explicitly. Your Summary claims them and
the substance is largely present inside Q1–Q7, but consensus has to be assembled by
comparing positions across four agents, and unlabelled positions cost a re-read each. Please
label them in round-02.

### @kimi-1

**Adopting — no holistic scorecards; a findings ledger instead.** This is the sharpest
correction anyone made, including to the research strawman. A weighted 0–10 aesthetic rubric
is contradicted by our own evidence section: LLM judges carry position bias worst at small
quality gaps, self-preference from −38 % to +90 %, verbosity correlation r ≈ .87, and human
inter-rater agreement on aesthetics of 38 %. A number produced under those conditions is
false precision that agents will then optimise. Brief-1's own anti-goal list says "no numeric
design score"; brief-2's D4 scorecard contradicts it, and you are right that brief-1 wins.
**The Decider receives typed findings, not scores.**

**Adopting — anonymity demoted to SHOULD, recusal as the real mechanism.** Four model
families have recognisable prose registers; a reader who has seen three files can often name
the fourth. I flagged this as a limit in round-01 but still leaned on anonymisation. Claiming
blind review we do not have is worse than not claiming it. Recusal is enforceable; anonymity
is aspirational. Say so in the spec.

**Adopting — grafts mechanically forbidden from touching the winner's token file.** Better
than my prose rule ("never a system layer") because it is checkable: the graft may not modify
`*.tokens.json`. A rule a tool can enforce beats a rule an agent must remember.

**Rejecting, partially — "one seeded re-diverge, then convergence is accepted as data".**
I agree the loop must terminate and I agree convergence-as-data beats an infinite rework
loop. But "accepted" is too passive for a signal this expensive to obtain. Convergence that
survives a re-diverge is the most interesting thing the roster can tell us, and it must
**escalate to the Decider with the axis table attached**, not quietly pass. Small change,
material difference: the Decider is told "four independent models converged on this after a
forced re-diverge", which is either strong evidence or a training artefact, and only a human
can say which for a given brief.

## Resolved disagreements A-G

**A. Second state machine, or reuse Parley's? — REUSE.** `parley-design` is a **profile over
the existing deliberation track**, not a phase ladder. Mapping:

| Design step | Parley home | Artifact |
|---|---|---|
| BRIEF | `00-prompt.md` + one human question set | `DESIGN-BRIEF.md` beside it |
| DIVERGE | round-01 (isolation is already structural) | `round-01/<agent>.md` + `<agent>.tokens.json` |
| DISTINCTNESS gate | facilitator, between rounds | recorded in `round-02/` preamble or an inbox note |
| CRITIQUE | round-02, one round | `round-02/<agent>.md` |
| DECIDE + GRAFT | `consensus.md` | Decider's recorded act + ≤3 typed grafts |
| CONTRACT | `FINAL.md` | the binding direction contract |
| APPLY | Phase 5 (ordinary implementer, D8) | code |
| AUDIT | Phase 6–7 review | review artifacts + checker output |
| SYSTEM | after Phase 8 | `DESIGN-SYSTEM.md`, written from shipped code |

What is lost by *not* having our own ladder: nothing that matters. What would be lost by
having one: two competing state machines, two definitions of "done", and a driver that
cannot tell which gate it is at. The phase *names* survive as citable vocabulary.

**B. Registry source of truth — ONE LITERATE FILE.** `references/RULES.md` is simultaneously
the machine source and the human source: each rule is a prose entry whose metadata lives in
a fenced block with a declared, boring grammar. No generated view, no drift test, no second
copy, because there is no second artifact to be out of date. Cost: the checker parses fenced
blocks out of markdown instead of reading JSON — roughly twenty lines, and it fails loudly on
a malformed block. Against `rules.json` + generated catalog: a drift test is a *detector* of
a failure mode that the literate file does not have. Against a pure-prose registry the
checker regex-scrapes: too fragile, and the grammar must be declared anyway.

**C. Size and file count — FOUR FILES, ≤80 KB total.** `SKILL.md` ≤200 lines (dispatcher,
when-to-use, when-NOT-to-use, and the invariants that must be in memory before editing
anything); `references/PDS.md` (the protocol); `references/RULES.md` (the literate registry);
`references/WEB-ANNEX.md` (target-specific hard numbers, explicitly non-normative for other
targets). Defence: N agents each pay the load, so per-file lazy loading — the economy both
prior-art projects rely on — inverts under a roster. Four files with fixed read sets is
cheaper *and* guarantees all four agents read the same bytes.

**D. Evidence-tier vocabulary — `T0`–`T3` plus `UNJUDGEABLE`.** I withdraw my
`stated|source|rendered|measured` names. Two of four participants already proposed T0–T3 and
naming the same thing three ways is the drift we keep warning about. One demand: every tier
carries a fixed one-word gloss in the spec so the numbers are never opaque — **T0 artifact
text · T1 source parse · T2 rendered DOM · T3 pixels** — and `UNJUDGEABLE` is a *compliant*
verdict, never silence.

**E. The dice — heterogeneity first, determinism as the remedy.** No roll on the happy path.
If gate G1 (distinctness) fails, the facilitator assigns each agent a distinct position on
the primary divergence axis, derived as `sha256(idea-slug + round + agent-id)` truncated into
the axis-position list. That is reproducible, auditable, offline, and checkable after the
fact — codex-1's determinism — but paid only when the roster demonstrably failed to diverge —
kimi-1's economy. If G1 fails a second time, the run does not loop: convergence is recorded
with the axis table and **escalated to the Decider**.

**F. The split line and the contract.** `parley-design` owns the doctrine and the registry.
Every rule entry carries `enforced-by: check:<id>` or `enforced-by: agent-judgement`, and
`tier: T0|T1|T2|T3`. `parley-design-check` declares `implements: PDS/1.0` and, in every
report it emits, the `registry-digest` (12 hex of sha256 over `RULES.md`) it ran against, so
a stale signature is detectable. **The checker reads the registry from the installed
`parley-design`; if that skill is absent it MUST refuse rule checks and say so, rather than
fall back to a bundled copy.** Structural and token checks that do not depend on the registry
may still run. Findings are emitted as `rule-id — violation — remedy`, always all three.

**G. Actively harmful.** Three things, named plainly:

1. **The weighted 0–10 aesthetic scorecard** (research strawman D4, and any round-01 echo of
   it). Under our own cited bias measurements it manufactures false precision, and a number
   in an artifact is a number agents will optimise. Replace with a typed findings ledger.
2. **Anonymisation as a MUST.** It cannot be delivered across four recognisable model
   registers, and a claimed-but-ineffective blind is worse than an acknowledged open one.
3. **Any bundled duplicate of the rule registry.** Fail closed, never fall back.

## Final positions on F1-F8

- **F1 — Authority.** Split by class: a `quality` finding with reproducible evidence lets one
  agent BLOCK; a `slop` finding never blocks unilaterally and needs two independent
  concurrences; a `system` finding binds only after a contract is ratified. Re-classifying a
  rule requires a major spec version.
- **F2 — Convergence.** Alarm, not pass. G1 fails on identical positions across all declared
  axes or a shared slop signature; remedy is one forced re-diverge; a second convergence
  escalates to the Decider with the axis table, never auto-passes.
- **F3 — Rendered evidence.** Optional with declared capability. `T0`–`T3` + `UNJUDGEABLE`;
  nobody may sign a layout verdict from evidence below T2; every artifact produced below full
  capability leads with a degradation declaration.
- **F4 — Dice.** Not on the happy path. Deterministic `sha256`-derived axis assignment is the
  remedy for a failed G1 only, and the assignment is recorded so it can be re-derived.
- **F5 — Size.** Four files, ≤80 KB total, `SKILL.md` ≤200 lines, fixed per-phase read sets.
- **F6 — Waivers.** One central waiver file; every waiver names rule id, narrow scope, reason,
  and expiry; a second participant must counter-sign; wildcard and unilateral suppression are
  forbidden; objective legibility and accessibility rules are design-system-blind and cannot
  be waived by widening the system.
- **F7 — Fast path.** Full ritual only for a new visual world or a change to a ratified rule.
  Everything else runs invariants + checker, single agent, no deliberation. Greenfield work
  can never use fast.
- **F8 — Decider.** The human by default; agent input is typed findings and is advisory and
  labelled as such. Unattended runs may let the facilitator decide, but must record
  `decider: agent (unattended)` and the decision is provisional until a human ratifies.
  Critics hold no decision authority; no agent scores or votes on its own direction.

## New concerns / questions

1. **Where does `DESIGN-BRIEF.md` live?** Under the profile model there is no D0 phase to own
   it. I propose it sits beside `00-prompt.md` in the idea directory and is ratified by the
   human in the same message that answers the one question set. If we cannot name its home
   cleanly, that is evidence the profile model has a genuine seam and I want to hear it.
2. **`DESIGN-SYSTEM.md` authorship is still unresolved.** Nobody in round-01 gave a
   satisfying answer. The winner's author knows the intent and is the worst party to describe
   what was actually built; a non-author lacks the build context. My current preference is the
   Phase-6 design reviewer, because reviewing the built code is already its job.
3. **Anti-goal drift is a real risk in the registry.** Every threshold we import was tuned on
   someone else's corpus. We need the contest path (§13 retro, advisory, quorum over
   self-preference) written into the spec on day one, not added later.

## Risks

- **Profile-over-track is right but under-specified.** Riding Parley's rounds is correct; it
  also means the design gates (G1 distinctness, G2 coherence) have no phase of their own to
  live in. They must attach to existing transitions explicitly, or they will be skipped.
- **Four files is a promise that erodes.** File count only ever grows. The ceiling needs to be
  a test, not a sentence.
- **The registry is the product, and it will be wrong at first.** If we cannot contest a rule
  cheaply and on the record, agents will either obey a bad rule or quietly ignore it, and both
  outcomes destroy the registry's authority.
- **We are still four agents who cannot see.** Nothing in the roster renders. Every claim in
  this design about layout, occlusion, or rhythm is `UNJUDGEABLE` at T2 until someone attaches
  rendered evidence, and the spec must make that embarrassing rather than invisible.
