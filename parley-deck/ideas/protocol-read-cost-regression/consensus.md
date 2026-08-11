---
idea: protocol-read-cost-regression
drafted-by: claude-1
date: 2026-08-10
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: consensus-reached
---

## What the owner asked

Why has Parley Deck got slower over the last few versions, and can it be optimized.

## The answer, in one paragraph

It is not the CLI — every command runs in under a second. It is the cost of a round, multiplied by
the number of rounds, and **both terms grew**. @codex-1's reconciliation is the finding that unifies
the round-1 disagreement:

> For the owner's unit of experience — kickoff to completed idea — **review churn dominates**. For a
> single response, **full-protocol loading dominates**. The statements are compatible: *the per-call
> tax is paid again inside every extra cycle.*

## The measurements

**Per call (claude-1, n=3 per arm, same agent, same question, same output length):**

```
arm A, reads COOPERATION.md in full   : median 98.7s  (27.3–105.3)
arm B, given only the relevant excerpt: median 29.9s  (21.1–39.2)     ratio 3.3x
```

Reported honestly: after the **first** replicate only, arm B looked slower and I briefly reported the
hypothesis refuted. Replicates reversed it. n=3 is small and arm A's variance is large.

**Per idea (claude-1, reproduced independently by @codex-1 and @kimi-1), 76 ideas split at
2026-07-01:**

```
design rounds  1.4 → 1.6   (flat)
review rounds  1.6 → 5.1   (max 24)
review bytes   20,237 → 146,290   (7.2x)
```

**Protocol growth:** 720 → 1,359 lines, ~12,300 → ~26,100 tokens in ten weeks; monotonic. @kimi-1
measured what grew inside it: **`MUST` 15 → 37, `MUST NOT` 6 → 15** — obligation density outgrew
prose.

## Where each cost actually lives — the round-1 finding that redirected the work

The two sources have **different levers**, and the diagnosis in `00-prompt.md` named neither
correctly.

- **The protocol is never embedded by the runtime.** `BuildRoundOnePrompt` reads `00-prompt.md` only
  (`internal/runner/runner.go:822`). It enters context through *instructions*: the skill's "Always
  read `COOPERATION.md` first" and **the facilitator's own hand-written prompts**. (@hermes-1;
  verified.)
- **Prior round files ARE embedded by the runtime.** `gatherPriorRounds` concatenates every
  participant artifact from rounds 1..N-1 (`runner.go:936-965`) and the prompt orders "READ every
  prior-round artifact below" (`runner.go:989`); consensus drafting orders another full-history read
  (`internal/app/driver_consensus.go:110,128`). (@codex-1; verified.)
- **The protocol does not require that.** Phase 2 requires addressing every active peer; it never
  requires re-reading every historical version of every position. **The CLI is stricter than the
  protocol it implements**, which makes this an implementation choice changeable **without a §7
  protocol change**.

## What drives the review-round explosion

@kimi-1 corrected the drafter's hypothesis here and the correction is adopted:

- The growth is a **step change dated 2026-07-28**, not a gradual protocol-growth effect — it tracks
  the roster growing, i.e. more reviewers producing more findings.
- The unbounded `deliberation` fix-up cap is therefore the **enabler of the tail, not the driver**.
- Decisively: the two worst ideas kept finding **fresh MAJORs at rounds 19–24**, so a severity floor
  would not have stopped them. A termination rule justified as "late findings are trivial" would be
  built on a false premise.

*(This is @kimi-1's evidence from its own commands. The drafter did not independently re-run the
2026-07-28 attribution and does not present it as verified.)*

## Q4 — prompt compression (the owner's direct question)

**Against, for normative text. Unanimous, reached by three different routes.**

- @hermes-1: *"must not edit another agent's file"* and *"don't edit other agent files"* are not
  equivalent in a system with append-only signoffs and refutation-default review — the difference
  between *must not* and *don't* is the difference between a protocol violation and a style
  suggestion. The 40–58% claim is measured on natural language, not on prose whose modal verbs,
  qualifiers and section locators carry the obligation.
- @kimi-1: the compressor's targets *are* this document's content — and measured the obligation
  density growth above to show it.
- @codex-1: ranked it **last of five** interventions, "never apply to normative text", and only as an
  experiment on already-resolved analysis.
- All three note it would be a **fourth copy** of a document that already has a drift guard because
  its copies diverge — and the skill's copy already differs from the embedded default.

**Narrow exception, agreed:** peer round files are arguments, not obligations. There, use
**structural** compression — frontmatter, verdict states, locators kept verbatim, prose dropped — not
**semantic** compression, which rewrites prose and can change meaning.

**The test that would decide it** (@hermes-1, adopted): take 10 normative rules from §15 and §7,
compress them, and have three agents independently apply both versions to the same review task. If
any agent's application diverges on a CRITICAL-or-MAJOR finding, compression fails. @hermes-1
predicts failure on at least 3 of 10. It is a falsifiable prediction with a pre-declared criterion,
which is why it is recorded rather than the savings estimate.

## Recommended interventions, ranked by saving ÷ risk

| # | Intervention | Needs §7? |
| --- | --- | --- |
| 1 | Phase-scoped protocol packet, with the **instruction** changed to point at it — extractive, dependency-checked, failing back to the full authority | Yes, if it changes normative text or skill behaviour |
| 2 | Replace full-history re-read with **previous round in full + a participant-owned carry-forward ledger**; the consensus drafter keeps the full-history read, because that is where §15.6's duty binds | **No** — implementation choice |
| 3 | A finite fix-up budget with severity- and scope-aware disposition **after** the budget, never auto-pass | Yes |
| 4 | Size alarms on round files — **alarm, never truncate**; a hard cap can suppress the only dissenting argument | No |
| 5 | Compression, confined to already-resolved analysis, as an experiment only | No |

**Rank 2 is the recommendation to act on first:** largest structural saving, quadratic term removed,
and it is the only one that needs no protocol change.

## Never cut

The authoritative protocol; applicable modals, negations, conditions and exceptions; round-1
independent proposals; every live objection or finding; provenance and verdict-conflict data; the
§15.6 audit; `FINAL.md`, the current diff, acceptance criteria and check results; explicit user
rulings; the ownership, no-secrets and §14 human-brake rules; and access to the raw historical
artifacts on demand.

## Drafter position changes

Per §15.5.

1. **The diagnosis in `00-prompt.md`.** I wrote *"The slowdown is not execution time; it is read
   cost."* Withdrawn twice over: @codex-1 showed the separation is false — context processing **is**
   execution — and my own 76-idea measurement showed review volume grew 7.2× against the protocol's
   2.1×. I had measured the design side of a single idea that happened to have **zero** review
   rounds, and generalized from the half of the lifecycle that did not grow.
2. **"The unbounded fix-up cap is the driver."** Withdrawn on @kimi-1's evidence that the growth is
   a step change tracking roster size, and that rounds 19–24 produced fresh MAJORs.
3. **The A/B result.** Reported as refuting the hypothesis after one replicate; reversed by two more.
   Recorded because the error is the same shape as the one in item 1.

## Process note

This idea was run on `track: standard` with prompts that deliberately did **not** instruct agents to
read `COOPERATION.md` in full. All three participants produced full analyses; @codex-1 and @kimi-1
each stated what the restraint cost them — residual uncertainty about unmarked cross-section
dependencies, and one case where a verdict-producing task would have required opening §15.1–15.2.
That is one data point in favour of scoping and is not a result.

If any recommendation here changes normative text or global skill behaviour, it escalates to a
meta-protocol-change idea under §7. **This idea diagnoses and recommends; it does not change the
protocol.**

## Signoffs

Each participant appends its own block. Do not edit another agent's block.

    ### <agent-id>
    verdict: OK | RESERVED | BLOCK
    date: YYYY-MM-DD

    <reasoning; if RESERVED, the reservation; if BLOCK, the counter-proposal>

### codex-1
verdict: RESERVED
date: 2026-08-10

I accept the diagnosis, ranking, and the owner's choice to implement rank 2, but reserve on the
carry-forward contract. This is not agreement by headcount: I independently checked Phase 2 and
`gatherPriorRounds`/`BuildRoundPrompt`. Phase 2 requires an explicit response to every other active
participant and a counter-proposal for disagreement; it does not require every historical artifact
to be reread. Removing that extra runtime instruction therefore needs no §7 change **provided** the
ledger is an implementation-scoped context optimization, not a new artifact-validity or consensus
rule. Missing, invalid, ambiguous, or challenged ledger state must fall back to full history. If the
implementation instead makes a new ledger duty or close condition normative, it does require §7.

The implementation item is: round 2 receives round 1 in full; round N>=3 receives every active
participant's round N-1 artifact in full and the participant-authored ledger state. Each ledger item
must carry an immutable owner-namespaced ID; kind (position, objection/counter-proposal, material
claim, verification verdict, user ruling, or exemption witness); exact scoped proposition rather
than a generated paraphrase; author/claim-owner/verifier identities as applicable; introduced and
current source path plus stable locator and SHA-256; lifecycle state
`OPEN|RESOLVED|DEFERRED|SUPERSEDED`; and an append-only transition history with the actor, round,
reason, and resolution/superseding-item locator. An objection remains live until its owner explicitly
disposes it. Terminal items retain tombstones and never silently disappear. Position changes carry
old/new item IDs and locators.

For every material claim, the ledger must preserve materiality, all claim owners, every verdict's
exact `CONFIRMED|WRONG|UNVERIFIED` state, verifier, `PRIMARY|SECONDARY|RECALL` tag, and decisive
evidence. `PRIMARY` retains the stable source passage or command, inputs, and relevant output;
`SECONDARY` retains the named dependency and an acyclic path to `PRIMARY`; `RECALL` stays explicit.
Contradictory verdicts on the same claim force `DISPUTED` and carry both verdict records plus any
resolution and dependency-check locator; neither the compiler nor a different participant may
invent a resolution. User rulings are carried verbatim with `owner: user`, without transferring
claim ownership.

Validation must reject duplicate or mutated IDs, silent deletion, unauthorized transitions,
dangling/cyclic `supersedes` or `SECONDARY` links, missing active-participant ledgers, unresolved
hashes/locators, and a verdict conflict not marked `DISPUTED`; rejection selects full history and is
visible in the prompt. Any citation, provenance challenge, supersession dispute, or conflict against
an older item also expands the referenced source verbatim. Raw artifacts remain directly available.

The round-3+ failure mode I would target in review is **orphaned dissent at the first compaction
boundary**: a round-1 minority objection is not repeated in round 2, its ledger entry is omitted or
closed by a non-owner, and round 3 receives neither the objection nor a fallback, creating false
convergence. A second required fixture should reword one material claim under a new ID and issue
opposing `PRIMARY` verdicts in different rounds; the system must join them as `DISPUTED` or fall back,
not lose the conflict through claim-ID forking.

This reservation also corrects two overstatements. My round-2 immediate proposal retained full round
1 plus the previous round and treated the validated ledger as a later optimization; the draft does
not record that disagreement. I accept the owner's alternative only under the contract above. And
the change removes repeated historical **prose**, not every quadratic worst case: a cumulative
ledger still grows with newly introduced items, so the implementation should measure packet growth
and avoid claiming a formal complexity reduction it has not demonstrated.

I found no additional borrowed result presented as the drafter's own verification. The A/B timing is
attributed to claude-1, the era counts include my independent reproduction, the runner and Phase 2
claims are independently checked above, and kimi-1's 2026-07-28 attribution and late-MAJOR evidence
are explicitly attributed and expressly not claimed as independently rerun by the drafter.

### hermes-1
verdict: RESERVED
date: 2026-08-10

I accept the diagnosis, the ranking, and the owner's decision to act on rank 2 first. This is not
headcount agreement: I independently verified the two code claims that the rank-2 recommendation
rests on, and they hold. Phase 2 (COOPERATION.md:338-343) requires "address every other active
agent explicitly" and "disagreement requires a counter-proposal" — it says nothing about re-reading
every historical artifact. gatherPriorRounds (runner.go:936-965) loops r=1..N-1 concatenating every
participant file, and BuildRoundPrompt orders "READ every prior-round artifact below"
(runner.go:989). The CLI is stricter than the protocol demands. Removing that instruction is an
implementation choice, not a §7 change. I also confirmed the second path: buildConsensusDraftPrompt
(driver_consensus.go:110,113) independently orders "Read EVERY round artifact under …/round-*/".
The consensus text's decision to keep the full-history read for the consensus drafter is the correct
response, because that is where §15.6's correlated-agreement duty binds (COOPERATION.md:347-349,
line 1327).

I reserve on two concrete implementation items.

Reservation 1 — the carry-forward ledger contract is underspecified in the consensus text. The text
says "previous round in full + a participant-owned carry-forward ledger" but does not enumerate the
ledger's required fields. codex-1's reservation block has the contract; I endorse it and add one
item it does not name: the ledger must carry, for any position that changed since the previous
round, the old and new claim IDs, the source locators for both, and the round in which the change
occurred. This is what §15.5's "Drafter position changes" (COOPERATION.md:1319-1325) requires of
the consensus drafter, and if round N+1 participants cannot reconstruct a drafter's position-change
chain from the ledger, the consensus drafter's §15.5 audit degrades when it later reads only the
ledger plus the previous round. The codex-1 contract already requires "position changes carry old/
new item IDs and locators," so this may be subsumed — but the consensus text itself should state
that the ledger must satisfy §15.5's position-change reconstruction, not just codex-1's reservation.

Reservation 2 — the rank-2 row addresses only the design-side embedder (gatherPriorRounds) and the
consensus drafter (driver_consensus.go). It does not address the third quadratic path that kimi-1
identified in round 2: gatherReviewContext (phase58.go:276-306) concatenates FINAL.md +
IMPLEMENTATION.md + every prior review round into each reviewer call. kimi-1 measured this at
427,970 bytes per reviewer call by round 21 of skills-cli-install-path — the largest of the three
embedders. The consensus text's "Where each cost actually lives" section names only the first two.
If rank 2 ships without touching gatherReviewContext, the quadratic cost persists on the review
side, which is where 7.2× of the byte growth lives. The consensus text's rank-2 row should name
gatherReviewContext as a third code path in the same change, or the recommendation will be read as
complete when it is not. I am not blocking on this because the owner's instruction was to act on
rank 2 as specified and the design-side fix is correct on its own — but the review-side embedder is
the next item in the same commit, not a separate idea.

On the five questions I was asked to check hardest:

1. Rank-2 is correctly specified and genuinely free of a §7 protocol change. Verified above. The
   one-line caveat is that if the ledger becomes a normative artifact-validity or consensus rule
   (codex-1's reservation states this), it does require §7. The consensus text does not state this
   boundary — codex-1's block does, and I concur with it.

2. The ledger must contain (endorsing codex-1's contract and adding the §15.5 position-change
   reconstruction requirement above): immutable owner-namespaced IDs; kind; exact scoped
   proposition; author/claim-owner/verifier identities; source path + locator + SHA-256;
   lifecycle state OPEN|RESOLVED|DEFERRED|SUPERSEDED with append-only transitions; material claims
   with all verdict states, provenance tags, and decisive evidence; DISPUTED flagging for
   contradictory verdicts; user rulings verbatim with owner: user. An objection stays live until
   its owner disposes it. Validation rejects duplicate/mutated IDs, silent deletion, unauthorized
   transitions, dangling/cyclic links, missing ledgers, unresolved hashes, and unflagged verdict
   conflicts; rejection falls back to full history.

3. The round-3+ failure mode I would target in review is the same one codex-1 names: orphaned
   dissent at the first compaction boundary — a round-1 objection is not repeated in round 2, its
   ledger entry is omitted or closed by a non-owner, and round 3 receives neither the objection nor
   a fallback. A second fixture: reword a material claim under a new ID and issue opposing PRIMARY
   verdicts in different rounds; the system must join them as DISPUTED or fall back, not lose the
   conflict through claim-ID forking. A third, specific to this idea's review-side gap: on a
   review round, a finding from round N-2 is closed in the ledger but the code region it flagged
   is re-touched in fix-up; the reviewer in round N sees neither the original finding nor its
   disposition because gatherReviewContext was not changed and the ledger is design-side only.

4. The text represents my positions accurately where it cites them. The "must not / don't"
   distinction and the §15/§7 compression test are quoted correctly. One inaccuracy: the text's
   ranking places gatherPriorRounds as rank 2 and the fix-up budget as rank 3, and says "Rank 2 is
   the recommendation to act on first." In my round-2 Q5 I ranked the fix-up cap #1 and
   gatherPriorRounds #2 — the exponent before the coefficient. The owner has decided to act on
   rank 2 first, which is the owner's call, but the consensus text does not record that my
   original ranking differed. This is a minor omission, not a misrepresentation; the drafter
   position-changes section records the drafter's own changes but not where participants disagreed
   on ranking. I note it for the record.

5. I found no borrowed claim presented as the drafter's own verification. The A/B timing is
   attributed to claude-1. The era counts are attributed with independent reproductions by codex-1
   and kimi-1. The runner and Phase 2 claims are marked "verified." kimi-1's 2026-07-28 attribution
   and late-MAJOR evidence are explicitly marked as not independently rerun by the drafter
   (lines 78-79). The protocol-not-embedded finding is attributed to me as "verified" and I did
   verify it (runner.go:821-822 reads only 00-prompt.md; COOPERATION.md is never embedded by the
   runtime). No issue.

Per §15.6: I am not agreeing mainly because others do. I independently checked the code paths and
the protocol text, and the rank-2 recommendation is sound on its own merits. My reservations are
about implementation completeness (the ledger contract and the review-side embedder), not about the
correctness of the change the owner has decided to make.

### kimi-1
verdict: RESERVED
date: 2026-08-10

I accept the diagnosis, the ranking, and the owner's decision to act on rank 2. Per §15.6, first:
this is not headcount agreement. For this signoff I independently re-ran every claim the rank-2
recommendation rests on rather than inheriting the other signoffs' checks: `BuildRoundOnePrompt`
reads only `00-prompt.md` (runner.go:822 — the protocol is never embedded); `gatherPriorRounds`
concatenates every participant artifact from rounds 1..N-1 (runner.go:936-965) and
`BuildRoundPrompt` orders "READ every prior-round artifact below" (runner.go:989);
`buildConsensusDraftPrompt` orders the drafter's full-history read (driver_consensus.go:113);
`gatherReviewContext` embeds FINAL.md + IMPLEMENTATION.md + every prior review round per reviewer
call (phase58.go:276-306). On the protocol side I read Phase 2 (COOPERATION.md:347-352), §15.5
(:1328-1334) and §15.6 (:1336-1354) in their own text.

1. **Rank 2 is correctly specified and genuinely §7-free — verified, not inherited.** Phase 2
   requires: address every other active agent explicitly, counter-proposal for disagreement,
   continue until nobody has new substantive objections. Nothing in it requires re-reading
   historical artifacts; the round-2+ template's "Position changes since round 1" is the author's
   own position, needing no peer history. The CLI is strictly stricter than the protocol it
   implements, so replacing the embedded history with previous-round-in-full + ledger, and changing
   the READ order to match, is an implementation change — **provided** the ledger remains an
   implementation-scoped context optimization. If ledger validity becomes a normative
   artifact-acceptance or consensus-close rule, it is a §7 change (codex-1's boundary; I endorse it
   after my own read). The drafter keeps the full-history read correctly: §15.6's close conditions
   are drafter-facing, and §15.5 requires exact prior quotations and source round paths that only
   full history (or a complete ledger) supplies.

2. **The ledger contract.** codex-1's reservation block is the design-side spec; I checked it
   against my round-2 Q2 requirements (self-authored, fixed-format, stable objection IDs, DISPUTED
   verbatim, provenance locators, supersedes links, fail-closed fallback to full text) and it covers
   and extends every one — critically, its append-only transition history with actor, round, reason
   and resolution locator is exactly what re-litigation detection needs. I adopt it and add two
   items it does not name: (a) **§15.5 position-change reconstruction** — old/new claim IDs, both
   source locators, and the round of change (hermes-1's reservation 1; without it the drafter's
   §15.5 audit degrades); (b) **code-region tags on disposed findings** — a finding on code whose
   disposition is recorded must rebut that disposition to reopen a cycle (my round-2 Q3 trigger 2);
   without region tags the ledger cannot feed the Stopping judgment (COOPERATION.md:637-647) or
   Phase 2's own "new substantive objections" test.

3. **The round-3+ failure mode I would target in review is re-litigation blindness** — the
   opposite direction of codex-1's orphaned dissent. His is false convergence; mine is false
   *non*-convergence, which is the expensive direction in this idea's own data. Concretely: a
   round-2 objection is disposed with a rebuttal locator; rounds 3+ no longer carry round 2; in
   round 4 a participant re-raises the same objection reworded; nobody can see it was answered
   without expanding the locator, so it counts as a *new* substantive objection and the loop
   re-opens settled ground — the exact pathology the Stopping judgment exists to catch. Fixture:
   dispose an objection at round 2 with a resolution locator, re-raise it reworded at round 4;
   the system must surface the prior disposition verbatim or treat the item as DISPUTED — never
   as new. A second fixture from my write-side concern: a round-3 file with a malformed or
   under-maintained ledger must trigger visible full-history fallback, not silent acceptance —
   otherwise under-maintenance degrades quietly toward orphaned dissent, and repeated fallback
   storms pay the quadratic anyway plus ledger-authoring output tokens, a cost shift the saving
   claim has not measured.

4. **Representation check.** Attributed to me accurately: MUST 15→37 / MUST NOT 6→15; the
   2026-07-28 step change tracking the roster; enabler-not-driver; fresh MAJORs at rounds 19–24;
   the compressor's targets being this document's content. Two gaps for the record. (a) The
   ranking table does not record that my Q5 ranking differed: I ranked the binding termination
   rule #1 (it attacks the dominant count term) and scoping the three embedders #3 as "CLI-only,
   no §7 needed" — so acting on rank 2 first is consistent with my own ranking and I do not
   object; but the record should show that my evidence says the count term dominates, and that
   rank 3's termination rule must be the trajectory-and-re-litigation rule from my round-2 Q3
   (no-convergence stop, re-litigation guard, severity floor + hard backstop that escalates and
   never completes), not only a budget. (b) The "narrow exception, agreed" on structural
   compression of peer round files is accurate only with the participant-owned qualifier: my Q4
   rejected *semantic* compression for peer files too ("it loses there too"); what I proposed was
   the self-authored ledger, which the rank-2 row correctly names participant-owned.

5. **Verified-claims audit.** No borrowed claim is presented as the drafter's own verification.
   The A/B timing is claude-1's and honestly reported with its first-replicate reversal; the era
   counts carry my and codex-1's independent reproductions (mine reproduces his table exactly);
   the runner and Phase 2 claims marked "verified" I re-verified myself above; my 2026-07-28
   attribution and late-MAJOR evidence are explicitly flagged as not independently re-run by the
   drafter (lines 78-79), which is the correct handling. One citation nit, not a misattribution:
   the drafter's full-history read order is at driver_consensus.go:113, inside the function the
   text cites as :110,128.

**Reservations (implementation items).** Reservation 1 — **rank 2 must name
`gatherReviewContext` (phase58.go:276-306) in the same change.** The consensus text's "Where each
cost actually lives" names only the design-side embedder and the consensus drafter; it omits the
third and largest embedder, which I measured at 427,970 B of embedded context per reviewer call
by round 21 of `skills-cli-install-path` (round-2 PRIMARY) — on the side that owns 7.2× of the
byte growth. If rank 2 ships design-side only, the review-side quadratic persists untouched.
Implementation item: the same change adds the review-side packet — FINAL.md + current
IMPLEMENTATION.md in full (the object under review, non-negotiable), the previous review round
in full for fix verification, a findings ledger (finding ID, severity, round introduced, status
open/closed/follow-up, disposition, code-region tag, locator to the round file), older rounds on
demand, fail-closed. hermes-1's reservation 2 is the same item; I record it independently as the
originator of the measurement. Reservation 2 — **the ledger contract is codex-1's block plus
items (a) and (b) above**, and the boundary condition of question 1 (non-normative or §7) is part
of the spec. Reservation 3 — **the saving claim is unmeasured against ledger-authoring cost**:
the ledger shifts read tokens to write tokens every round, and my round-1 write-side hypothesis
is still unrefuted by telemetry (the log is dormant since 2026-06-02). Before rank 2 is declared
a win, one instrumented idea (my round-2 Q1: per-call phase, agent, wall duration, assembled
bytes, output bytes) should confirm the felt-time mechanism — the optimization this idea
recommends is itself subject to the rule my round-2 Q5 stated: never ship an optimization before
the measurement that justifies it.

### claude-1
verdict: OK
date: 2026-08-10

I drafted this, so my OK is not independent verification of it — §15.1 applies to every claim in it
that is mine, and three of them I have already withdrawn in "Drafter position changes".

All three participants signed RESERVED and all three independently re-verified the code claims the
rank-2 recommendation rests on, rather than inheriting each other's checks. That is the outcome I
wanted from a signoff and it is worth more than three OKs would have been.

I accept every reservation as an open implementation item rather than arguing any of them:

1. **@codex-1's boundary is adopted as binding on the implementation.** The ledger is an
   implementation-scoped context optimization. It must NOT become an artifact-validity or consensus
   rule — if it does, the change needs §7 after all. Missing, invalid, ambiguous or challenged
   ledger state falls back to full history.
2. **The ledger contract belongs to @codex-1's and @hermes-1's specifications, not to my one-line
   phrase in this document.** Both reserved on exactly that under-specification and both are right;
   "participant-owned carry-forward ledger" is a label, not a contract.
3. **@codex-1's two failure fixtures are mandatory**: orphaned dissent at the first compaction
   boundary, and claim-ID forking with opposing PRIMARY verdicts.

**Scope correction from @kimi-1's signoff, adopted.** It found a fourth path none of the rounds
named: `gatherReviewContext` (`internal/runner/phase58.go:278-299`) embeds `FINAL.md`,
`IMPLEMENTATION.md` and every prior review round on each reviewer call — the same quadratic shape as
`gatherPriorRounds`, on the review side. Since the measured 7.2× growth is in review, **this is the
larger target of the two**, and the implementation must cover both.

One thing I want on the record because it changes how careful the implementation has to be. Phase 2
rule 1 reads "**Silence = implicit agreement**" (verified, `COOPERATION.md`). The protocol therefore
converts an omission into consent. If a round-1 objection drops out of the ledger and its author does
not restate it, the protocol does not record a lost datum — it records agreement that was never
given. @codex-1's rule that only an objection's owner may dispose of it is the single thing standing
between compaction and a false consensus, and it is not optional.
