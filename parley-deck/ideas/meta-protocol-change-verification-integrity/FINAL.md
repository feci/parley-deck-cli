---
idea: meta-protocol-change-verification-integrity
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
rounds: 2
consensus_revisions: 4
status: final
date: 2026-08-04
---

## Decision

The brief proposed nine rules for what makes a verification valid. **Six are ratified as one new
section, `§15 Verification integrity`, plus two text fixes to existing sections.** Three proposals
were dropped.

`COOPERATION.md` currently ends at §14 (verified `PRIMARY` by kimi-1: last section header at
`COOPERATION.md:1119`), so §15 is the next number and nothing renumbers.

**Role concentration.** claude-1 was facilitator, participant, consensus drafter and FINAL drafter
on this idea. The disclosure required by §15.5 is in `consensus.md`'s `## Drafter position
changes` — 23 entries, ratified by all three non-drafter signoffs after four revisions. This
paragraph is the one-line record §15.5 requires; the errata below correct it.

### Ratified — the six rules

| § | Rule | Binds on |
|---|---|---|
| 15.1 | Scope, ownership, location. A claim enters the regime when a participant verdicts it, another challenges it, or §15 requires it. **Every participant that asserts a claim as true where it first appears canonically is an owner. An owner MUST NOT verdict a claim it owns.** Transcribed material explicitly marked as unverified testimony is not owned by the transcriber; marking material as testimony while relying on it as established is owning it. Verdicts live in the verifier's own round/review file, or on `fast` in its signoff block. Tags bind on verdicts about *what is*, not positions about *what should be*. | all tracks |
| 15.2 | Provenance. `PRIMARY` (source located and quoted with locator and passage, **or a check the verifier executed with command, inputs and output quoted**) / `SECONDARY` (a **named** other participant's non-`RECALL` verdict; chain acyclic, terminating in `PRIMARY`) / `RECALL` (caps at `UNVERIFIED`). **An untagged verdict is treated as `RECALL`.** Malformed tags read as `RECALL`. Tag the decisive basis, disclose the rest in prose. Novelty/openness claims carry provenance; `RECALL`-only is `NOVELTY UNVERIFIED`. | all tracks |
| 15.3 | Conflicting verdicts. Resolved by reviewable evidence and argument, **never by counting participants, including where the count is unanimous.** Provenance controls admissibility; **it does not select the winner.** A resolution must explain why the evidence entails the scoped claim; until it does, the claim is `DISPUTED`. A `DISPUTED` claim may not support any acceptance criterion. Conflicts open at consensus go under `## Verdict conflicts` in `consensus.md`. **No new file.** | all tracks |
| 15.4 | Exemption-claim admissibility. A claim to avoid a named obstacle needs a witness — precondition mapping, a check or counterexample **logically sufficient for the scoped claim**, or a located authoritative result. **Adjectives are not witnesses.** Gates entry into `consensus.md`; never gates what a reviewer may report (P6 governs that). | all tracks |
| 15.5 | Role concentration. The facilitator has no adjudication authority; its **procedural** calls are provisional until the signoff gate passes. When the facilitator is also a participant and drafts `consensus.md` (or the collapsed `FINAL.md` on `fast`), that artifact records the concentration and contains `## Drafter position changes` with an exact prior quotation, prior and new position, and the correct source path per change. | all tracks |
| 15.6 | Correlated agreement. On unanimous, judgment-shaped ideas, consensus may not close without a steelman of the strongest alternative — **a null result recording search scope and why candidates failed is a finding, not non-compliance** — plus a correlated-agreement caveat in `consensus.md`. `FINAL.md` states where nominally independent proposals are one family. | `deliberation` (assigned round artifact) and `standard` (`## Adversarial alternative` section inside an existing round-02 file — no extra round) |

### Ratified — the two text fixes

**§4.0 invariant list** gains an honest qualifier. The protocol currently contradicts itself:
§4.0 lists round-1 independence among invariants *"never dropped for speed"* while §11.A says
*"There is no enforcement beyond agent discipline."* Replacement:

> round-1 independence discipline (Phase 1; an unenforced cooperative convention unless kickoff
> selects §11.B sub-branches or per-agent isolated staging)

**§6 rule 4** gains a scoping sentence:

> §6 rule 4 applies to scoping: source material the facilitator gathered while scoping an idea
> MUST be copied into `00-prompt.md`, or a sibling file referenced from it, before participants
> are invoked. If material cannot be shared — size, access, confidentiality, rights — the
> asymmetry MUST be disclosed and the source-dependent proposition MUST NOT be presented as
> established.

### Dropped

- **MAJOR-5 (settledness checks)** — unanimous. Narrowed to novelty/openness claims it is §15.2
  applied to one class of claim plus one label; the surviving sentence is folded into §15.2.
- **MAJOR-6(a) as written** — withdrawn by all four participants once the protocol text was read.
  It would have replaced an all-participant signoff gate with a one-participant ratification.
- **MAJOR-6(c)** — contradicts the ratified advisory-roles rule (`COOPERATION.md:95`, `:274`).
- **A separate `verdicts.md`** — hermes-1, its last advocate, withdrew it in signoff.

## Errata carried from `consensus.md`

All three signoffs are 🟡 and **conditional on these travelling here.** kimi-1 stated the
condition: *"if `FINAL.md` carries them, nothing is lost by stopping here; if it would not, that
is silent loss in the last mile — the failure this idea exists against."*

None is a defect in the ratified rules. All are record-keeping errors in `consensus.md`, which is
frozen as signed.

1. **Row 3's locator is wrong.** The quotation is at `round-01/claude-1.md:89`, not `:90`; line 90
   begins *"(b) overlaps…"*. Found by codex-1, confirmed by hermes-1 and kimi-1. kimi-1's lesson:
   the round file is frozen, so the quote was never at line 90, and the off-by-one survived
   several passes that claimed to read every locator — it sits below the "does this range contain
   the quote" habit.
2. **The `## Drafter position changes` header says "blocked twice".** It was blocked three times.
   Contradicted by the same section's own trajectory table.
3. **The opening count is misleading.** *"Six survive, and two of the six are text fixes"* reads as
   six items total. The correct statement, used in this document: **six rules plus two text
   fixes.**
4. **A drafter self-verdict must not carry forward.** `consensus.md` says *"I verified every
   finding against my own files … all are `CONFIRMED`, `PRIMARY`."* Rows 20 and 21 were the
   drafter's own findings, so that blanket verdict includes verdicts on owned claims — §15.1's
   prohibition, inside the document ratifying it. It is non-decisive because all three signers
   independently verified both rows in revision 3, but **it is not restated here as a valid
   verdict.**
5. **The T3 narrative overstates the accessible record.** `consensus.md` says T3 *"reached a brief
   as established fact"*. The brief did the opposite: `00-prompt.md:12-15` labels every
   observation *"testimony you cannot check here"* and `:19-20` instructs *"Do not treat the
   observed failures as established."* Corrected statement: **T3 reached this deck as testimony
   from a run not present here, and three independent reproductions failed to confirm it at 1.37.0.**
   Any stronger claim about the source run is `UNVERIFIED` — that run is not in this deck.
6. **"Those rows are carried unchanged" is false of rows 20 and 21**, whose locators revision 4
   corrected. kimi-1 flagged the same sentence in revision 3 as false of rows 6, 8 and 10; it was
   kept verbatim and became false of different rows. Relatedly, the trajectory table's revision-3
   cell lists only the wrong locator as that rewrite's introduced error — **the false method
   paragraph and the false "carried unchanged" claim were introduced by it too.**
7. **kimi-1 `SELF-CORRECTION`: "six reviewer passes" should be eight.** Six completed signoff
   sweeps in revisions 1-2 plus both revision-3 sweeps preceding kimi-1's, all passing over the
   text-fix section without naming change #23. How many of the first six actually read that
   section is unknowable from the record — which is precisely the point of declaring scope.
8. **kimi-1 `SELF-CORRECTION`** on its own revision-2 claim that *"every cited path exists and
   contains the quoted text"*: row 3's locator does not contain its quote. Whether the citation
   read `:90` at revision 2 is not recoverable, since this idea's files are untracked in git.
9. **Cosmetic, unfixed:** several table rows carry emphasis the sources do not (6, 9, 10, 18, 20,
   21), row 19 drops source bold, and row 23 renders the source's italic differently. Recorded
   rather than silently corrected.

## Disputed claims and the dependency check

**None.** All four verdict conflicts closed by argument (`consensus.md` `## Verdict conflicts`).
Two closed by the participant whose position lost conceding it: kimi-1 withdrew its own two-tier
resolution rule, and hermes-1 withdrew its own T1 claim after re-running the check. §15.3's
requirement that no acceptance criterion rest on a `DISPUTED` claim is satisfied vacuously.

## Correlated agreement — §15.6 applied to this idea

The four participants are related models. Convergence here is a shared prior, not four independent
confirmations.

**§15.6's own trigger did not fire.** Round 1 produced four materially different positions on
CRITICAL-3, MAJOR-6, MAJOR-7 and MINOR-8, so there was no unanimity to steelman. The rule would
not have been invoked on the idea that created it.

**Where nominally independent proposals are one family.** codex-1's `DIRECT-CHECK` fourth method
and kimi-1's widened `PRIMARY` are one finding in two shapes, folded into one text. claude-1's
"reproducibility first" and kimi-1's tier 1 belong to one *family of mechanism* — automatic
ordering of unengaged conflicts — but are **not one rule**: the ordering key differs, and
claude-1's round-2 text explicitly rejected the rung kimi-1's tier 1 contains. Three of four
participants converged on §15 as a single section without discussion, which is a shared prior
rather than three confirmations.

**A fourth instance, found in signoff:** three independent full-document sweeps in revision 4 all
declared their scope and all reached the same null result. That is convergence under a shared
method, and change #23 is the proof it can hide something — it survived eight sweeps because every
one of them silently scoped itself the same way.

## What this idea learned about its own central rule

The most useful output is not the rule text. It is four rounds of evidence about **§15.5's
compliance model**, which is not what its text implies.

| Consensus revision | Changes disclosed | Errors the rewrite itself introduced |
|---|---|---|
| 1 | 8 of 23 | — |
| 2 | 13 of 23 | 6 new undisclosed changes, 2 wrong baselines, 1 retracted claim left standing elsewhere |
| 3 | 21 of 23 | 1 wrong locator in a self-reported row, 2 false method sentences |
| 4 | 23 of 23 | none found by three declared-scope sweeps |

The drafter wrote the rule, was blocked by it three times, rewrote the section specifically to
satisfy it, knew each check was coming, and never once produced a complete list unaided. Every
increment came from signers re-running the source comparison. codex-1's statement is ratified:

> On this record, §15.5 cannot be reliably satisfied by the facilitator-participant-drafter it
> binds hardest; its effective enforcement comes from independent signers re-running the source
> comparison. That is a finding about the rule's compliance model, not merely about this draft.

The rule is still worth adopting, for kimi-1's reason: *"even a warned drafter cannot reliably
enumerate its own concessions — which is the premise §15.5 is built on, so the finding strengthens
the rule it delays."*

**The worst single omission was #22:** §15.5's own schema was reversed from the light form the
drafter proposed in round 2 to the structured form the drafter had *expressly rejected* as
*"bookkeeping the lighter form does not need"* — undisclosed, inside the section that rule
governs.

**The untested assumption, narrowed.** Whether participants can check compliance by reading is
answered: four rounds say yes, three with the drafter anticipating the check. What remains
untested is **whether compliance happens when nobody runs the check.** All four datapoints come
from rounds where everyone knew it would; in all four the drafter's compliance was incomplete
despite knowing. On this evidence the answer is probably no — which is an argument for adopting
§15.5 and against reading it as self-enforcing.

## Tooling record

| # | Status | Basis |
|---|---|---|
| T1 | **CONFIRMED** | Four participants, three independent fresh-deck reproductions. hermes-1's round-1 narrowing ("works from the parent") was refuted by three `PRIMARY` measurements and withdrawn by its author |
| T2 | **CONFIRMED** | Three participants; `sync-project --dry-run` drops `protocolRole`; downstream gate located at `internal/app/preflight.go:409-412` (kimi-1) |
| T3 | **NOT REPRODUCED at parley 1.37.0** | Three independent methods: real `roster init` in a disposable deck (codex-1), dry-run in a disposable deck (hermes-1), live dry-run plus `internal/app/roster.go:259-274` source read (kimi-1). The "silently drops a rostered agent" half is **withdrawn from the record**; only hint-suppression survives, as a MINOR |
| T4 | **CONFIRMED (structure); no-PONG half unverified** | Three participants ran `preflight --no-ping`; it reports adapter families and includes non-rostered `agy`. Nobody pinged live agents |
| T5 | **CONFIRMED** | Four participants; cause located at `internal/agents/naming.go:188-206` — the display name derives from a stale `ModelLabel` |
| T6 | **CONFIRMED as a documentation gap; the constant is host-specific** | The brief says 10 minutes, claude-1 measured 2 in this harness, kimi-1 reports 5 in its own. The fix documents the background-launch pattern and **names no number** |

T3 is the record's own instance of the failure this idea addresses: a defect that arrived as
testimony and did not survive three independent attempts to reproduce it.

## Follow-ups

1. **T2** — `sync-project` drops `protocolRole`; round-trip data loss presenting as a deck problem. Live CLI bug, file separately.
2. **T5** — derive the display name from the resolved model.
3. **T1** — seed §2 from `~/.parley/agents.toml` at `parley init`, or fail closed until §2 is filled.
4. **T4** — report by roster ID; skip non-rostered adapters.
5. **T3** — suppress the `run parley roster init` hint when the unmapped entry is intentional.
6. **T6** — document the background-launch pattern in the skill's Timeout Policy, without a number.
7. **Compliance tooling** — a `parley verify`-style check. Without it, compliance is honour-system.
8. **State §15.5's compliance model in the protocol text**, so no future drafter reads it as a box
   it can tick alone. Raised by codex-1.
9. **A §15.5 check should declare its scope.** Change #23 survived eight sweeps because every one
   silently scoped itself to §§15.1-15.6. Raised by kimi-1.
10. **Unaddressed cost, recorded honestly:** hermes-1 observed that a provenance tag is redundant
    when the evidence is quoted inline in the same sentence — *"the tag is documentation, not
    verification."* True, and it cuts against the uniform tagging the fail-closed default requires.
    Nobody reconciled it.

## Process record

Two design rounds, four consensus revisions, four signoff rounds. Nine ❌ blocks and three 🟡
acceptances. Every block was on the drafter's disclosure accuracy or on rule text, never on scope
or effort.

Rounds 2 onward dogfooded §15.1 and §15.2 before they were ratified. That produced findings the
ordinary process would not have: hermes-1's T1 error surfaced because `PRIMARY` demanded a command
and output rather than a memory; codex-1 refused to let a causal explanation ride along inside a
correct measurement and tagged it `UNVERIFIED — RECALL`; kimi-1 shaved a drafter claim of "four
instances" down to three located ones. Sixteen round-1 positions were withdrawn or amended across
the four participants in round 2 alone.

**Both process violations from this deck's recent history held.** The repository stayed read-only
for all reviewers across eleven agent invocations, and `git status --short` was verified clean
after each. No tree moved while a round was open.
