---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency-v2
round: 1
date: 2026-09-16
---

## Launch record and read scope

`context_mode=full`, `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`,
packet equal to source; a shadow packet audit was attached as an unapplied diagnostic. I also read
the live `parley-deck/COOPERATION.md` directly, so every `COOPERATION.md:NNN` locator below is to the
live file as read in this session. I did not compute its hash (no shell in this launch), so the
attestation's hash equality is not something I verify here; the content I rely on is read first-hand.

Read: this idea's `00-prompt.md`; `source-context/pilot-preparation.md`; `source-context/provenance.json`;
`../meta-protocol-change-evidence-first-efficiency/FINAL.md`; `../meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md`;
live `COOPERATION.md` in full. Not read: the other participants' `round-01` files (Phase 1 independence,
`COOPERATION.md:331`), and the sibling evaluation repository (blocked — see V9). Nothing here accepts,
verifies or closes the original implementation or PR #73; this is a design position only.

## Summary

The amendment I support is deliberately small: supersede only the prospective *arm membership* words in
D8/AC-X1 (full-six → an explicit full-four ID list), pin the two user-given constants (15 minutes shared
per task-arm; USD 15 total prospective experiment spend), and change nothing else — not the 12 tasks, not
the 36 cells, not the 28 packet runs, not the thresholds, bands, blinding or denominators. Beyond that I
propose a concrete deterministic allocation for the 36 pilot cells, and a spend design whose enforcement
quantity is a *reservation debited before each call*, never the provider's reported cost. The hard part is
honest: with ~340 planned model invocations and a USD 15 ceiling, I can state the arithmetic of the
division but I cannot state feasibility, because no admissible price basis is readable from this launch —
so the amendment must make "zero treatment calls, feasibility unproven, denominators intact" a first-class
reportable outcome rather than a failure to be engineered around.

## Source claims — verification verdicts (§15.1–§15.2)

These are verdicts about **what is** in the named documents. Throughout I keep two things apart: that a
document *contains* a sentence (which I can verify by reading it) and that the sentence is *true of the
world* (which usually I cannot). None of these are claims I own; my proposals below carry no verdicts.

| ID | Claim | Tag | Verdict | Locator |
|---|---|---|---|---|
| V1 | The packet experiment is 12 matched pairs = 24 paired runs, plus 3 canaries and 1 full control = 28 runs | PRIMARY | CONFIRMED | packet FINAL:109–111, :116; original FINAL:104–106, :232 |
| V2 | The live §2 generated roster view contains a header and no data rows | PRIMARY | CONFIRMED | COOPERATION.md:138–141 |
| V3 | The live protocol's cost ceiling is telemetry-gated: "Cost enforcement is telemetry-gated — it applies only once the runner emits `agent.usage` events" | PRIMARY | CONFIRMED | COOPERATION.md:679–680 |
| V4 | Unknown spend must stop or carry a documented conservative reservation, "never zero"; spend-boundary telemetry writes fail closed | PRIMARY | CONFIRMED | original FINAL:50–54; 00-prompt.md:44 |
| V5 | Reported cost is an estimate, not an invoice; missing cost is missing, not zero | PRIMARY | CONFIRMED | pilot-preparation.md:70–72; original FINAL:52 |
| V6 | The (0.67, 0.80] refute band is unresolved and "not resolved by count"; hermes-1 holds one of its two recorded positions; the implementer "must not resolve this by picking one" | PRIMARY | CONFIRMED | packet FINAL:134, :140, :147–149; original FINAL:186–189 |
| V7 | The 12-task pilot uses full context in **every** arm; packet A/B is a separate experiment | PRIMARY | CONFIRMED | original FINAL:97–98; pilot-preparation.md:30–31 |
| V8 | The copied preparation is a draft, explicitly not preregistration, with no treatment calls made | PRIMARY | CONFIRMED | pilot-preparation.md:3–4; provenance.json:3–5 |
| V9 | This participant cannot read the sibling evaluation repository in this launch | PRIMARY (check I ran) | CONFIRMED | `Glob scripts/pilot_*.py` at `…/parley-deck-evaluation` → "is outside …/worktrees/evidence-first-amendment; --restricted confines the file tools to the working directory" |
| V10 | The copied preparation refers to "four executable, four review and four design tasks **below**" but contains no task enumeration in its 78 lines | PRIMARY | CONFIRMED | pilot-preparation.md:18–19, whole file |
| V11 | §4.0's deliberation per-agent timeout is ~30 min, while the experiment ceiling is 15 min shared by all participants and phases of one task-arm | PRIMARY | CONFIRMED (both texts) | COOPERATION.md:237; 00-prompt.md:42–44 |

**V1 arithmetic, shown so it can be checked:** 6 pairs/phase × 2 phases = 12 pairs; 12 pairs × 2 runs =
24; + 3 canaries + 1 control = **28**. One ambiguity is worth killing before freeze: packet FINAL:110 reads
"6 paired runs per phase," which alone could be misread as 6 runs (3 pairs). AC-P2's "the exact 12 matched
packet pairs" (original FINAL:232) and D5's "six matched AB/BA pairs EACH" (:104–105) resolve it to 6
*pairs* per phase. A planner who used the looser reading would silently run a half-sized experiment.

## Assumptions I could not verify (recorded as unverified, not relied on as established)

- **A1.** Existence, content and fitness of `scripts/pilot_{acceptance,grading,analysis}.py` and
  `delivery/2026-09-05/pilot/` (asserted at 00-prompt.md:66–67). UNVERIFIED for me by V9. §6 rule 4
  (COOPERATION.md:751–755) requires that unshareable material's asymmetry be disclosed and that the
  source-dependent proposition not be presented as established — so decision 5's "reuse existing
  preparation scripts and tests" is not assessable by participants under current access.
- **A2.** The 12 task fixtures themselves. By V10 they are not in the participant-visible corpus. If the
  copy is byte-exact as `provenance.json` records, the upstream draft does not enumerate them either.
- **A3.** Any provider pricing or per-call cost. None is readable here and I will not supply one.
- **A4.** That PR #73 is incomplete and that three implementation calls were in flight (00-prompt.md:79–81,
  :97–99). I verify only that the document says so.
- **A5.** My per-cell call-count model (used below for arithmetic only, never for a dollar claim).

## Proposed approach

### P1 — Exact 12-task allocation (decision 1)

Arm roster, frozen as an explicit ID list rather than a `roster show` snapshot (V2 shows the deck's own
generated view carries no rows, so "all six active roster IDs frozen from `parley roster show`"
(pilot-preparation.md:27) is not reconstructible from the deck): **A=claude-1, C=codex-1, K=kimi-1,
Z=zcode-1**. Tasks T01–T04 executable, T05–T08 review, T09–T12 design (pilot-preparation.md:18–19). The
table is deterministic — no RNG, no seed — so it is reproducible from the rule and checkable by anyone.

| Pos | Task | Type | Solo | Duo drafter→critic | Full drafter / critic-lead | Arm order |
|---|---|---|---|---|---|---|
| 1 | T01 | exec | A | A→C | A / C | S,D,F |
| 2 | T05 | review | C | C→K | C / K | S,F,D |
| 3 | T09 | design | K | K→Z | K / Z | D,S,F |
| 4 | T02 | exec | C | Z→A | Z / A | D,F,S |
| 5 | T06 | review | K | A→K | A / C | F,S,D |
| 6 | T10 | design | Z | C→Z | C / K | F,D,S |
| 7 | T03 | exec | K | K→A | K / Z | S,D,F |
| 8 | T07 | review | Z | Z→C | Z / A | S,F,D |
| 9 | T11 | design | A | C→A | A / C | D,S,F |
| 10 | T04 | exec | Z | K→C | C / K | D,F,S |
| 11 | T08 | review | A | Z→K | K / Z | F,S,D |
| 12 | T12 | design | C | A→Z | Z / A | F,D,S |

Balance properties, all checkable against the table: each identity solos exactly 3 tasks, one of each
type; all 6 unordered pairs appear exactly twice, once in each role order; each identity drafts 3 duo
cells and critiques 3; each identity drafts 3 full cells and leads critique on 3, never both on the same
task; each arm occupies each ordinal position exactly 4 times. Two properties are deliberate concessions
to the budget: the **run order interleaves types** (pos 1–3 are exec/review/design), and every consecutive
triple of positions has three distinct solo identities — so a budget stop mid-run leaves a type- and
identity-spread prefix rather than four executable tasks and nothing else. Positions 1–6 cover all six
duo pairs; positions 1–8 keep full-arm drafting perfectly balanced.

Not engineered, therefore recorded: at some positions the solo identity equals the full-arm drafter, and
in the full arm every identity sees every task. Cross-arm carryover is controlled by *fresh session and
fresh workspace per cell* (pilot-preparation.md:31–32), not by identity separation; residual provider-side
caching is unknown and gets reported, not assumed away.

### P2 — The shared clock (decision 1)

One elapsed clock per task-arm cell: started before the first participant call, stopped at the cell's last
artifact, covering every participant, every phase and the orchestration between them, never reset per call
(00-prompt.md:42–44). Overruns terminate the cell and are recorded as `timeout` outcomes retained in the
denominator — they are never extended, and a timed-out member is never dropped (pilot-preparation.md:33).
By V11 this ceiling is *tighter* than §4.0's per-agent default and is per-arm rather than per-agent: in
the full-four arm, four participants plus orchestration share 15 minutes. Expect the full arm to time out
more often than solo **by construction**. That must be pre-registered as a design property; reporting it
afterwards as "larger teams are slower" would be a measurement artifact, and pilot-preparation.md:36–38
already forbids claiming equal compute from equal wall time.

### P3 — Full context in every pilot arm, and the one place that rule does not apply

Every pilot cell launches with `context_mode=full` (or a visible `full-fallback` with its reason) and
never `packet` (V7). A pilot cell whose attestation records `packet` is recorded as `protocol block`, not
silently kept. The complementary guard matters just as much: the **packet experiment is the experiment
that varies context mode**, so the full-context rule must never be back-propagated onto its 28 runs. Two
rules that look contradictory read against different experiments; the amendment should say so in one line
so nobody reconciles them by weakening one.

### P4 — The 28 packet runs, unchanged (preserve thresholds and bands)

Unchanged and restated only so the amendment cannot be read as touching them: 6 AB/BA pairs in each of
phases 1 and 6; packet generation time inside the packet arm; 3 canary replicates all of which must pass,
plus the full control; ship at `R = median(packet/full) ≤ 0.50` in **both** phases with zero obligation
misses; the standalone correctness veto at any speed; the middle band returning the measured number
unrounded; a non-implementer recomputing both ratios (packet FINAL:108–117).

**The band must survive the roster change.** By V6 the (0.67, 0.80] dispute is explicitly unresolved and
explicitly not resolvable by count, and one of its two positions is hermes-1's. hermes-1 leaving the
prospective quorum does not retract it, and zcode-1 does not inherit it. §15.3 is directly on point:
conflicts are "resolved by reviewable evidence and argument, **never by counting participants**"
(COOPERATION.md:1307–1308). If the amendment is silent here, a membership change quietly resolves a
pre-registered decision rule in the surviving side's favour — which is the post-hoc adjustment the packet
FINAL's own withdrawal record (packet FINAL:126–130) exists to prevent. I propose the amendment state, in
one sentence, that both recorded positions and the band survive verbatim and still route to the user.

### P5 — Dollar enforcement versus reported cost (decision 2)

The enforcement quantity is **not** the provider's reported cost. By V5 reported cost is an estimate that
arrives *after* the money is spent; using it as the gate would enforce nothing. The design, reusing what
already binds rather than inventing a regime:

1. **Reserve before work.** Each call is preceded by a durable reservation debit with an idempotent action
   identity — D6's existing read/reserve boundary (original FINAL:118–123), not a second ledger.
2. **The reservation constant comes from the operator.** `r` per call class must rest on an
   operator-supplied price basis or invoice evidence. No participant may supply it: by A3 any number I
   produced would be `RECALL`, and V4 forbids the zero default. If no admissible basis exists at freeze,
   `r` is undefined, the gate cannot be evaluated, and **zero treatment calls are made** (see P8).
3. **Reconcile after, never refund into the cap.** Reported cost is stored beside the reservation for
   reporting and coverage. If reported cost exceeds the reservation, the excess is debited and the cap
   re-checked; if cost is unknown, the reservation stands as the charge. Under-reservation is the only
   unsafe direction, so reservations are set at a documented upper bound and overflow reconciliation makes
   an underestimate self-correcting rather than silent.
4. **Retries are separate spent attempts** (original FINAL:48–49): each needs its own reservation, and a
   failed call's reservation is not returned.
5. **Stop rule.** Before each unit: if `remaining < reservation(unit)`, stop. The stop is a recorded
   outcome, not an error to retry around.
6. **Ledger shape.** Adopt the deck's existing effects-ledger schema — append-only per-record file,
   `planned → executing → succeeded|failed → reconciled`, idempotency key, reconcile-before-retry
   (COOPERATION.md:1150–1154) — as a *schema*, without adopting §12's opt-in pipeline machinery.

By V3, the protocol's own cost ceiling only bites once `agent.usage` is emitted, and that emission is
itself part of the unfinished D2 work (original FINAL:45–49). So **the existence of a working dollar gate
is a precondition to be observed before treatment, not assumed at freeze** — and observing it is exactly
the kind of claim that needs a non-owner's PRIMARY verdict rather than the implementer's own.

### P6 — Funding order, whole units, and grading reserved with treatment

- **Whole-unit funding.** Reserve an entire unit before starting it. A packet *pair* is atomic (a half
  pair yields no ratio at all). A pilot *task block* — 3 arm cells plus their grading — is atomic, because
  the paired analysis is per task (pilot-preparation.md:69–70); scattered orphan cells cost money and
  produce no paired delta.
- **Grading is reserved at the same moment as its treatment.** Otherwise the cap produces candidate
  outputs nobody can score: money spent, measurement zero. Adjudication carries its own reserve.
- **Order.** Packet pairs and canaries first (their 28 runs are exactly specified and carry the ship
  gate), then pilot task blocks in the frozen run order of P1, then grading, then adjudication. This is a
  policy proposal for the group to ratify, not a measurement, and it is the *pre-registered* order —
  which is what stops the not-run set from being chosen after the fact by which results looked bad
  (pilot-preparation.md:45–46 forbids selectively omitting expensive failed tasks).
- **Sub-caps derived from call counts, not invented percentages.** Required spend is
  `28·r_packet + Σ(cells · r_class) + 72·r_grade + adjudication reserve`, with class-specific `r` because a
  full-four cell is roughly four times a solo cell in invocations. If the sum exceeds USD 15, the answer is
  *not* to shrink scope (00-prompt.md:54–55): run the prefix that fits and report the shortfall.
- **Scale, stated as arithmetic and nothing more.** Under A5 — solo ≈ 2 calls, duo ≈ 6, full-four ≈ 12 per
  task — the pilot is ≈ 240 invocations, plus 28 packet runs, plus ≈ 72 gradings: **≈ 340 model
  invocations**, i.e. **USD 15 / 340 ≈ USD 0.044 per invocation** on average. That division is arithmetic
  on two numbers, and I stop there: whether USD 0.044/invocation is achievable is UNVERIFIED by A3, and I
  decline to imply feasibility in either direction. What follows is a freeze requirement, not a guess:
  **the exact planned call count per cell class must be a frozen field**, so the denominator of that
  division is fixed before treatment rather than discovered during it.
- **One question I will not settle alone:** whether AC-P2's non-implementer recomputation draws from the
  USD 15. My position is that it is verification, like review, and 00-prompt.md:23 separates review cost
  from experiment budget — but it is prospective spend, it materially changes feasibility, and it belongs
  in the user's hands rather than in a participant's reading. Recorded as an open question, not decided.

### P7 — Two blind grading-only configurations (decision 3)

- **Disjointness.** Grader identity ∉ {task fixture authors} ∪ {candidate authors}. All four arm members
  are candidate authors in some cell, so **both graders must come from outside the arm set**. D8 already
  permits grader-only configurations as observers "never extra quorum votes" (original FINAL:159–161);
  using a non-arm configuration as a grader therefore does not re-admit anyone to quorum, does not revive
  a prior audit role, and carries no signoff weight. By A2/V10 the fixture authorship is not yet recorded
  anywhere I can read, so disjointness is currently unevaluable — recording fixture authorship is a
  precondition for choosing graders, not a detail.
- **Blinding, enforced mechanically.** Grader inputs are built by allowlist into a manifest — anonymized
  candidate output, task bytes, rubric, deterministic test results — with arm, model, agent, cost and
  order metadata excluded (pilot-preparation.md:51–55). The manifest is hashed; any field outside the
  allowlist fails closed rather than being filtered at read time. A pre-grading scan flags identity-shaped
  tokens (agent IDs, model names, arm labels); leaks that cannot be removed without changing substance are
  **reported, not silently edited** (pilot-preparation.md:55).
- **Family overlap is structural — report it, don't pretend.** With two grader configurations and four
  candidate families, some grading is family-matched. Record the match flag per graded cell and report
  matched-versus-unmatched score deltas as a diagnostic. §13.3 is explicit that a self-preference score is
  diagnostic and never an acceptance criterion; LE-3 (COOPERATION.md:272) is the same shape one layer over.
- **Disagreement.** Adopt pilot-preparation.md:64–67 unchanged as the trigger — any critical-invariant
  conflict, any pass/fail disagreement, or a weighted rubric difference above 15 points — resolved by an
  actual third blind review from a third disjoint configuration if its reserve exists, and otherwise by
  recording the cell **`unresolved`**. A contradiction is never averaged into a score.
- **Budget note that is also a design argument.** For executable tasks the primary signal is the hidden
  deterministic acceptance tests (pilot-preparation.md:58–59), which consume no model calls. The cheapest
  evidence channel is also the least gameable one; graders stay secondary there, which is both better
  measurement and less spend.

### P8 — Honest not-run reporting (decision 4)

Use the source's own state vocabulary verbatim — `not-run`, `provider failure`, `timeout`, `protocol
block`, `valid final output` (pilot-preparation.md:47) — and add only a `not_run_reason ∈ {budget-stop,
beyond-frozen-order, feasibility-unproven}` sub-field, not a new state. Denominators are fixed at freeze
and never re-baselined to what completed: **12 tasks, 36 pilot cells, 28 packet runs, 72 planned
gradings**. Every reported rate shows `k/36`, `k/28`. Failures and timeouts stay in intention-to-treat
counts (original FINAL:157). And the outcome the amendment must name as valid in advance: if no admissible
price basis exists at freeze, the run is **zero treatment calls, 64 cells not-run with reason
`feasibility-unproven`, denominators intact** — reported as an honest unfinished experiment. Per D5 an
unrun trial is unfinished work, not a result (original FINAL:112–113); what is *not* acceptable is
inventing a cap, shrinking the task set, or letting an unbounded run discover the limit by spending.

## Existing alternatives (§15.6(a))

Mechanisms this proposal would otherwise hand-build, against what already ships, with locators:

| Mechanism | Already shipped / specified | Locator | Disposition |
|---|---|---|---|
| Spend ceiling | LE-5 loop budgets: max steps, wall-clock, best-effort cost; `~/.parley [defaults.loop]`; `parley run --max-driver-steps/--max-wall-clock` | COOPERATION.md:673–680 | Extend; do not rebuild. Gap: cost is telemetry-gated (V3), so the dollar gate is a precondition to observe |
| Charge-before-work, idempotent spend identity | D6 shared read/reserve boundary, durable idempotent action identities, fail-closed on corrupt state | original FINAL:118–123 | Reuse as the experiment ledger boundary; no second ledger |
| Effects ledger with reconcile | §12.6/§12.7 append-only per-effect records, idempotency key, reconcile before retry | COOPERATION.md:1150–1154 | Adopt the schema only; §12's pipeline machinery is opt-in and not required here |
| Unknown-spend policy | "visible stop or a documented conservative reservation, never zero" | original FINAL:50–54 | Adopt verbatim; it already decides the hardest case |
| Elapsed ceiling | §4.0 per-agent timeout; D7 liveness/watchdog observation | COOPERATION.md:237; original FINAL:133–147 | Gap: no shared per-arm clock across participants and phases — new, and disclosed as new (P2) |
| Freeze / preregistration format | Packet FINAL's pre-registered table; "Ceilings to Freeze"; write-once freeze, amendments as new versions | packet FINAL:106–117; pilot-preparation.md:3–5, :36–43 | Adopt; one write-once freeze file |
| Cell state vocabulary | not-run / provider failure / timeout / protocol block / valid final output | pilot-preparation.md:47 | Adopt verbatim; add only a reason sub-field |
| Grader model-diversity check | LE-3 `require_model_diversity` (reviewer-vs-implementer, not grader-vs-candidate) | COOPERATION.md:272 | Borrow the shape as a diagnostic flag; not a gate |
| Blind grading harness | Nothing in the protocol. `parley consult` exists but is explicitly advisory and non-canonical — never quorum evidence | COOPERATION.md:820–825 | Build, unless A1 resolves. Sibling `scripts/pilot_*.py` claimed at 00-prompt.md:66–67 is unreadable here (V9) |
| Counterbalanced allocation | Nothing shipped. The fixed seed in the source governs bootstrap analysis, not allocation | pilot-preparation.md:69 | Build as a deterministic table (P1) — no RNG, so nothing to seed or drift |

**Scoped null — sources consulted:** live `COOPERATION.md` (full), both binding FINALs, this idea's
`00-prompt.md`, `source-context/pilot-preparation.md`, `source-context/provenance.json`, and one executed
access check (V9). **Not consulted:** provider pricing (none in the workspace; no web in this launch), the
CLI source tree, the sibling evaluation repository (access-blocked), and peer round-01 files (Phase 1
independence). Where a mechanism is marked "nothing shipped," that is scoped to these sources.

## Proposed narrow amendment

Supersede **only** the following, leaving the original FINAL, its kickoff, rounds and signatures
untouched (00-prompt.md:27–33):

1. **Arm membership.** In D8 and AC-X1, "full six active roster IDs" / "full-six" → **full-four: exactly
   `{claude-1, codex-1, kimi-1, zcode-1}`**, frozen as an explicit ID list (V2: the generated view cannot
   supply it). Every other word of D8/AC-X1 stands.
2. **Constants.** Shared per-task-arm elapsed ceiling = **15 minutes**, covering all participants, all
   phases and orchestration, not reset per call. Total prospective experiment spend ≤ **USD 15**, spanning
   pilot, packet experiment, grading, adjudication and retries.
3. **Construct change, recorded not buried.** The "full" arm now measures a four-participant roster. The
   original six-participant question is **superseded, not answered**; the coordination-overhead contrast is
   smaller than the one D8 specified, and any report must say so rather than presenting full-four results
   against a full-six hypothesis.
4. **The packet band survives the roster change.** Both recorded refute positions, including hermes-1's,
   and the (0.67, 0.80] routing to the user, stand verbatim (V6, §15.3).
5. **No retroactive effect.** No prospective roster change re-attributes, revives or promotes any
   historical artifact, observer note or signature.

Explicitly **not** amended: 12 tasks; 36 pilot cells; 6 pairs per phase and the 28 packet/canary/control
runs; `R ≤ 0.50` in both phases; the correctness veto; all other decision bands; full context in every
pilot arm; ITT denominators and retained failures; two blind non-author graders; the no-vendor-ranking and
no-historical-causal-claim limits. And not amended in the other direction: this is design work, not
implementation acceptance — PR #73's gates remain open.

## Observable gates

| ID | Required observable result |
|---|---|
| AC-V2-1 | A single write-once freeze file records: 12 task IDs with byte hashes, the full 36-cell allocation of P1, per-cell planned call counts by class, grader configurations, rubric and hidden-test hashes, model/effort per ID, the 15-minute ceiling, the USD 15 cap, each reservation constant **with its price basis**, and the funding order — hashed before any treatment call. No call executes whose cell is absent from the freeze. |
| AC-V2-2 | Every pilot cell's participant set equals its frozen set; the full arm is exactly the four IDs; no substitution, no silent drop; a failed or timed-out member remains in the ITT denominator. |
| AC-V2-3 | Each cell records one shared elapsed clock spanning all participants, phases and orchestration; overruns appear as `timeout` outcomes and are never extended. |
| AC-V2-4 | Every call is preceded by a durable idempotent reservation debit; the run stops when `remaining < reservation(unit)`; reported cost is recorded separately, reconciled after, and never used as the gate; unknown cost reconciles to the reservation, never to zero; total charges ≤ USD 15 is verifiable from the ledger alone. |
| AC-V2-5 | Exactly 24 paired runs (6 AB/BA pairs × 2 phases) plus 3 canaries and 1 control = 28; no pair is executed partially; thresholds, veto and bands are unchanged; both recorded band positions are carried forward verbatim. |
| AC-V2-6 | Two grader configurations, each disjoint from fixture authors and candidate authors; each grading input matches an allowlist manifest that provably excludes arm/model/agent/cost/order metadata; family-match flags recorded; disagreement triggers a third blind review or an explicit `unresolved` cell; no contradiction is averaged. |
| AC-V2-7 | Every pilot cell lands in exactly one of `valid final output`, `provider failure`, `timeout`, `protocol block`, `not-run`(+reason); the report shows k/36, k/28, k/72 against frozen denominators; not-run cells are the tail of the pre-registered order. |
| AC-V2-8 | Every pilot cell's launch attestation shows `context_mode` `full` or `full-fallback` with a reason; a `packet` attestation in a pilot cell marks that cell `protocol block`. |
| AC-V2-9 | If no admissible price basis exists at freeze, zero treatment calls occur and the result is reported as feasibility-unproven with all denominators intact — not as a smaller experiment. |
| AC-V2-10 | Before consensus, the task fixtures and any reused preparation scripts are copied into `source-context/` (or listed with hashes and interfaces) so every participant can evaluate them; until then no participant claims reuse or fixture adequacy. |

## Concerns / open questions

1. **Does AC-P2's non-implementer recomputation draw from the USD 15?** Materially changes feasibility.
   My position is no (verification, like review), but this is the user's call, not a participant's reading.
2. **Pair semantics before freeze.** Confirm 24 paired runs, not 12 (V1). Cheap to confirm, expensive to
   discover late.
3. **Which two grader configurations are actually invokable and author-disjoint?** Requires a bounded
   liveness observation per D7 — an exact-PONG-class result, never inferred from silence — and requires
   fixture authorship to be recorded first (A2).
4. **Is zcode-1 observed ready?** 00-prompt.md:79–82 says three implementation calls were still running
   and that readiness must come from completed calls or a new bounded check. Same rule: observation, not
   silence.
5. **Does the 15-minute clock include workspace setup/teardown and grading?** I propose: treatment
   orchestration yes, grading no. It needs to be written down before it is measured.
6. **Where do the 12 fixtures live and who authored them?** By V10 they are not in the participant corpus.

## Risks

- **Feasibility risk, stated plainly.** ≈340 planned invocations against USD 15 (P6). If the price basis
  makes that impossible, the correct outcome is a truthful partial or zero run — the risk is that pressure
  to "complete the experiment" turns into a quiet scope shrink, which is precisely what 00-prompt.md:54–55
  forbids.
- **Gate-does-not-exist risk.** By V3 the dollar ceiling only binds once `agent.usage` flows, and that
  emission is part of unfinished work. Running treatment on an unverified gate is unbounded spending with
  a cap written on paper — the exact "printed caps bind only where enforcement lives" failure recorded at
  packet FINAL:94–96.
- **Over-reservation.** Conservative reservations may exhaust the cap while actual spend is far below it.
  This is the safe direction and is accepted; reserved-versus-reconciled totals must both be reported so
  the gap is visible rather than read as real spend.
- **Ceiling asymmetry.** The shared 15-minute clock systematically disadvantages the full arm (V11/P2).
  Pre-register it; never report it as a quality finding.
- **Truncation bias.** A budget stop makes the completed set a prefix of the frozen order. Mitigated by
  interleaving types and identities across positions (P1), but the completed subset is still not a random
  sample, and the report must say which cells never ran and why.
- **Blinding leakage.** Protocol wording and style leak identity even after metadata stripping
  (pilot-preparation.md:55); with four candidate families and two graders, family overlap is structural.
  Both are reportable diagnostics, not things to assert away.
- **Construct drift.** Six→four changes what the "full" arm means. Comparing the result to D8's original
  hypothesis without saying so would be a silent redefinition of the measured quantity.
- **Dispute erosion by membership change.** Without amendment item 4, hermes-1's departure resolves a
  pre-registered open band by attrition (V6, §15.3).
- **Access asymmetry.** By V9/A1, material the facilitator can read and participants cannot is already
  load-bearing in decision 5. §6 rule 4 requires disclosure and forbids treating the source-dependent
  proposition as established — AC-V2-10 is the concrete remedy.

No self-verification verdict and no signoff is issued here: my proposals above are positions about what
should be done, which carry no provenance tag (§15.1), and the claims I do verdict are other authors'
canonical statements, none of which I own.
