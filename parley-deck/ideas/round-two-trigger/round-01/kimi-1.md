---
agent: kimi-1
idea: round-two-trigger
round: 1
date: 2026-09-01
---

## Summary

My lens is measurement: what makes this trigger falsifiable rather than ceremonial. Three
findings frame everything else.

**1. The premise's locator is stale — the clause it cites no longer exists (PRIMARY).** The
00-prompt says §15.6(b)'s *existing* language is "round 1 closes with no substantive
disagreement". It is not existing. That sentence was the preamble of the **old** §15.6
("Correlated agreement"), introduced by `bfca39e` (2026-08-04) and **replaced** by `59eb663`
(2026-08-29, protocol-generation-bias IMPLEMENTATION — commit message: "15.6 replaced across all
three COOPERATION.md copies: 1372 B -> 896 B"). Current §15.6 (`parley-deck/COOPERATION.md:1346-1361`,
verified by reading at HEAD = freeze commit `2d17478`) carries (a) existing-alternatives,
(b) the shared-prior record in `consensus.md`, (c) alternatives disposition. **No close-condition,
no round-02 trigger, nothing to machine-detect.** The only stop-condition prose that survives
anywhere is `COOPERATION.md:359` — "Continue until nobody has new substantive objections" — which
nothing reads. So the honest framing is not "make an existing clause checkable" but "the clause
was deliberately removed three days ago by the idea that ratified the carrier thesis, and this
idea must decide whether it is designing something *new*." Whether anything acts on the old
language: nothing can — the text is gone. That is a PRIMARY null result, and it changes what any
protocol text for 2.12.0 would even say.

**2. The checkable fragment is disk state; the rest must be a recorded declaration (PRIMARY).**
Everything a trigger may legitimately evaluate is already disk-observable: track, participant
set, artifact presence/validity, marker scan (`❌ BLOCK`, `DISPUTED`, `Counter-proposal`,
`ALT-`), and the `cross_review_rounds` budget. What is *not* disk-observable is the semantic
residue — whether agreement is substantive, and the protocol-mutation-diversity instrument limit
applies unchanged: keyword screens detect movement, not correctness, and four participants
converging together on a wrong answer is invisible to every check we can build. A trigger that
pretends to evaluate the semantic part will be ceremonial by construction. A trigger that
evaluates the disk part, **records the evaluation whether it fires or holds**, and leaves the
semantic residue as a written facilitator declaration, is falsifiable.

**3. My independent recount reproduces the freeze (PRIMARY).** At HEAD=`2d17478` I count 81 idea
dirs with ≥1 `round-NN` dir and 29 single-round — *including* `round-two-trigger` itself.
Excluding it: **80 and 28**, exactly the frozen numbers, and the 29 names include all four
`deliberation` closes named in the freeze plus `tui-editor-composer` (the one `fast`). The M5
defect recurred live: the corpus moved under measurement again, this time by one idea.

**Position:** build the deterministic evaluator plus its mandatory record; do **not** build a
semantic gate. The conflict of interest is not removable by instrumentation — only *priced*: a
facilitator closing against a recorded "would-fire" must write a disagreement with a
deterministic instrument into the audit trail, where today closing costs nothing and leaves no
trace. §15.5 (`COOPERATION.md:1331-1336`) already makes procedural calls provisional until
signoff — the record gives signers something to ratify against.

## Proposed approach

Advisory, from the measurement lens. Four components, in dependency order.

**P1. Define the checkable close-condition, disk-only.** Inputs: `track:` and `participants:`
from `00-prompt.md`; round-01 artifact presence and validity; a marker scan of round-01 bodies
for `❌ BLOCK`, `DISPUTED`, `Counter-proposal`, `ALT-`; the configured `cross_review_rounds`.
Output exactly one of `would_open_round_02` / `would_hold` / `incomplete`. Anything not in that
input set is not the trigger's business. Track-gating is mandatory: `fast` already skips Phase 2
(`internal/track/track.go:182`; `COOPERATION.md:232`), so the evaluator must be inert there —
the floor-cost constraint makes this constraint-forced, not a design choice.

**P2. Record every evaluation, hold included, or the instrument does not exist.** Today's defect
is that a considered-and-declined round-02 leaves no trace, so "closed correctly" and "closed
early" are byte-identical afterwards. The record must be idea-local, append-only, and its
*absence* must itself be detectable — otherwise we re-ship the same silent-close defect with a
new filename. Determinism and auditability are stated constraints; this is where they bind.

**P3. Falsify by pre-registered replay against the frozen 80 before any protocol text is
drafted.** State the rule first, replay once, publish the raw 80-bit fire/hold vector and its
disagreements with history (52 opened / 28 held). Discrimination targets fixed in advance: the
rule must fire on the two protocol-change deliberation closes
(`meta-protocol-change-devx-speed`, `protocol-restructure-appendices`) and must hold on the
bulk of the 19 pre-track small ideas; everything else is reported, not tuned. If no statable
rule discriminates, that null result closes this idea without a core change — the D9-shaped
tie-break. With n=4 deliberation closes there are 6 pairwise distances; per-idea criteria are
unreadable (D7), so every success metric here is corpus-level. Overfitting guard: one replay,
rule frozen before it, no iterate-to-fit on n=4.

**P4. Measure over-triggering from the other side.** Of 141 round-02+ artifacts carrying
`## Position changes since prior round`, 23 explicitly report no change — a keyword lower bound
(SECONDARY, frozen measurement, instrument limits stated there). A trigger aggressive enough to
catch all 28 closes converts some of the 52 correct opens into no-movement rounds; the honest
cost model reports both directions, in **call counts and wall time only** — provider input-token
telemetry does not exist (`internal/driver/loop.go:174-175`), so no token or currency claim is
admissible.

**Behavior on the 40-of-41 decks without §15.6 (installed core 2.10.0; 2.11.0 staged).** The
fleet half is SECONDARY (from the 00-prompt); the local half is PRIMARY — this deck's
`parley-deck/meta/version.json` reads `deckVersion 2.11.0`, `updatedAt 2026-08-29`, and the
2026-08-07 changelog entry measured §15 present in only 5 of 36 decks, directionally consistent.
On a 2.10.0 deck the evaluator must run **advisory-only and never fail closed**: it reads only
idea-local artifacts and writes only its own record; a deck lacking the prose simply lacks the
duty, not the tooling. Any validator wired to fail on a missing section would break 40 decks —
that failure mode is named here so the mechanism-design role can design against it. Protocol
text (2.12.0) would add the *duty to produce the record*; it cannot retroactively add it
elsewhere, and must say so.

**Carrier recommendation, from coverage evidence.** The auto-driver only exists on the
`parley run` path; openviking finding 5 (PRIMARY there) counted exactly **three** `_index.md`
files in the whole deck because manual facilitation runs neither writer. The frozen 80 are
overwhelmingly manually facilitated. A CLI-only evaluator the driver runs reproduces that 3-of-80
coverage failure, and a CLI-only change that quietly alters deliberation semantics is worse than
a protocol change (A6). So: the *duty* belongs in the round-prompt/template layer where the
carrier thesis says rules actually reach compliance; the deterministic evaluator is the
instrument that keeps the duty honest; `internal/protocol/roundartifact.go:23`
(`ValidateRoundOneArtifact`, §15.6a) is the existence proof that a template duty plus a validator
is shippable. What does **not** change: who decides. The facilitator still decides; the record
makes the decision visible and referenceable.

## Existing alternatives

Mechanisms my proposal builds by hand, each with the closest thing that already ships. Every
locator below was opened and read today at `2d17478`; each is PRIMARY unless tagged otherwise.

**M1 — a deterministic close-condition evaluator over disk state.**
Closest shipping: the auto-driver's round boundary, `internal/driver/driver.go:289-344`
(`advanceRound`) gated by `roundComplete` (`driver.go:352-408`) — a two-signal disk check
(artifact presence + `runner.ValidateRoundArtifact`, plus `responding-to` and per-agent
`### @<other>` headings for round ≥2, plus a terminal `round.completed` event). Its promotion
decision is **purely budgetary**: `CurrentRound >= 1+CrossReviewRounds` (`driver.go:301,310`),
budget from `cross_review_rounds` in `00-prompt.md` (default 1, `internal/driver/transport.go:32-36`),
clamped per track (`internal/track/track.go:182` fast→0; `:217` standard cap 2; `:196-197`
deliberation cap 3). **It never reads artifact content for agreement** — a repo-wide grep for
`converg` over Go sources returns no convergence detection anywhere (hits are comments, prompt
strings, render idempotence). So the evaluator's *shape* ships; the *signal* does not.
**Constraint-forced** (determinism/zero-new-dependencies force a stdlib-only disk evaluator).

**M2 — an append-only, idea-local evaluation record, written on hold as well as fire.**
Closest shipping: `consensus.Reopen`'s audit append (`internal/consensus/consensus.go:397-405`
— renames the blocked draft, appends `## Reopen reason`); the §4 inbox escalation file
convention (`COOPERATION.md:705-718`); and the driver's cursor/event log
(`driver.go` `commitCursor`, `store` events). All three are records of a *decision taken*; none
records a *decision evaluated and declined*. **Constraint-forced** (the auditability constraint:
"if it does not fire, the record must show it was evaluated").

**M3 — a fire action cheaper than a full round.**
Closest shipping: `parley consult` (dispatch `internal/app/app.go:86-87`; prompt builder
`internal/runner/consult.go` — read-only, single agent, stdout-only) and, stronger, the LE-7
goal-done gate: `BuildGoalCheckPrompt` (`internal/runner/consult.go:72-89`) — a fresh agent
checks observable criteria **read-only, once, before close**, ending in a parseable
`GOAL-CHECK: PASS|FAIL` line the driver parses (`internal/app/driver_impl.go:401` builds,
`:421-439` parses). That is the shipped precedent for "a cheap single-agent check at a close
boundary with a machine-readable verdict", which is exactly the fire-action shape.
**Inherited** (the requirement for *a* cheap action is constraint-forced by the cost-conflict
finding; the consult/goal-check carrier is convention I inherit, not something the constraint
dictates).

**M4 — a corpus replay harness for falsification (P3).**
Closest shipping: `parley retro scan` (dispatch `internal/app/retro.go:54`; `retro.Scan` at
`internal/retro/retro.go:57`), which already walks every idea and emits per-idea signals. But
its scoring **prices extra rounds as friction** (`retro.go:132-134`: `Rounds > 1` adds
"extra design rounds" weight 1.0/round, review cycles 2.0/round) and **filters low-friction
ideas out of the coreset** (`retro.go:182,199`) — so an idea that closed too early presents as
*healthy* to retro. Retro scan structurally cannot surface premature closes; the replay must be
a separate read-only pass, not a retro patch. **Constraint-forced** (falsifiability requires a
frozen-corpus replay; D8 requires the freeze).

**M5 — a human/quorum gate on the boundary.**
Closest shipping: §4.0 track table, "deliberation: human gate at each transition"
(`COOPERATION.md:238`), and `parley loop tick`'s §14 brake (`COOPERATION.md:1216-1223`: an
automated loop MUST NOT promote, run, finalize, or reopen; `internal/app/loop_cmd.go:18-102`
implements tick as candidate-only, disabled by default, `--enable` still candidate-only).
The gate already exists on deliberation — **and the four deliberation single-round closes
happened anyway**: a gate without an instrument did not catch them. Evidence, from the corpus
itself, that the missing piece is the instrument, not the authority. **Inherited.**

**Required spot checks, reported:**

- **§15.6(b) "no substantive disagreement" language** — REMOVED 2026-08-29 by `59eb663` (diff
  read in full); nothing acts on it because nothing remains. Current §15.6 preamble
  (`COOPERATION.md:1348-1350`) states only (a) is machine-validated and says so explicitly —
  verified; the one validator is §15.6(a)'s `## Existing alternatives` non-empty check
  (`internal/protocol/roundartifact.go:10,23-28`, enforced on the run path per
  `internal/runner/roundonegate_test.go`).
- **`parley consensus status|draft|signoff`** — ship and work as advertised: dispatch
  `internal/app/app.go:521-537`; `consensus.Status` `internal/consensus/consensus.go:96`,
  `Draft` `:159`; triage states `:24-28` computed from signoff artifacts `:580-588`. None of
  them looks backward at whether round-02 should have happened. **Inherited, not a trigger.**
- **`parley consensus reopen --reason`** — ships (`internal/app/app.go:661-684`;
  `consensus.Reopen` `internal/consensus/consensus.go:370-412`) but fires **only when
  triage=blocked** (`:382-384`): it is a post-consensus back-edge for a BLOCK, not a
  premature-close detector. An idea that closed after round 1 without ever reaching consensus
  has nothing for reopen to act on. **Inherited; wrong object for this idea.**
- **Auto-driver at a round boundary** — covered in M1: budget-only, content-blind.
- **`parley retro scan`** — covered in M4: scores the extra round as cost, filters the
  prematurely-closed out as low-friction.
- **Inbox escalation path** — ships and is the only existing channel by which a non-facilitator
  could challenge a premature close (`COOPERATION.md:695-724`, `blocking: yes|no`, answer quoted
  into the next round file). But it requires someone to *notice* the close; the close itself is
  unrecorded, so there is nothing to notice. **Inherited; downstream of the missing record.**
- **`parley loop tick`** — outer-loop, pre-idea discovery only (M5); §14.2 forbids it from
  reopening or promoting anything. It cannot fire a mid-deliberation trigger; wrong layer.

**Null-result note (required).** "The current judgment call is correct and should stay" remains
a live outcome: 52 of 80 ideas did open round-02, and when cross-review runs it moves positions
(141/141 substantive, 23 no-change, with stated instrument limits). What the corpus cannot tell
us is the per-idea truth of the 28 holds — and my proposal does not pretend to fix that. It
fixes the *measurability* of the decision. If P3's replay shows the facilitator's historical
52/28 split already matches any statable disk rule, the correct verdict is: keep the judgment
call, add only the record.

## Concerns / open questions

1. **The stale locator is a round-01 correction, not a detail.** The idea as briefed ("make
   §15.6(b)'s existing language checkable") targets text deleted three days ago by a ratified
   implementation whose idea *also* ratified the reason it was deleted: an enforceable leg (a)
   shipped, the unenforceable leg (b) was withheld precisely because "printing a rule nothing
   enforces would be the fourth instance of this deck's own defect class" (`59eb663` commit
   message). Any 2.12.0 text is therefore a **new** duty, and must pass the test the removed
   clause failed. Round 2 should re-scope the protocol-text question accordingly; claude-1
   (protocol fit) owns where, if anywhere, it lands.
2. **De jure, opening a round is not the facilitator's act at all.** `COOPERATION.md:337`: "any
   agent may open round N+1". The facilitator-decides framing is de facto (someone must create
   the dir; the driver auto-advances; manual facilitation centralizes it). The conflict-of-interest
   analysis should target the *record*, not the *authority* — the authority is already diffuse
   and unused.
3. **Marker absence ≠ agreement.** A round-01 artifact may carry `Counter-proposal` or `ALT-`
   aimed at the prompt itself; and disagreement can be written without any marker. The disk
   fragment detects "no *recorded* dispute", and the record must say exactly that, in those
   words — anything stronger is unfalsifiable.
4. **The uncheckable residue has a name.** Four related models converging on a wrong answer
   inside one shared frame is §15.6(b)'s actual concern and is invisible to every disk check
   (protocol-mutation-diversity instrument limit, carried). If that residue is to be addressed
   at all, it is a written facilitator declaration at close — and §15.6's preamble is the
   template for shipping that distinction honestly rather than implying a gate.
5. **Coverage is the binding constraint on any carrier.** Only `parley run` runs the driver and
   its writers (three `_index.md` in the whole deck, openviking finding 5). Whatever we build
   that only runs on that path covers a small minority of the corpus. I do not know the exact
   manual-vs-run split of the frozen 80; if hermes-1's mechanism role can measure it, that
   number decides between template-duty-first and CLI-first.
6. **Replay gaming.** With n=4 deliberation closes, any of us can fit a rule to the corpus by
   inspection. Pre-registration (rule published before the one replay) is the only guard I see;
   a second replay on held-out ideas only if the first passes, mirroring D4's shape.

## Risks

- **Re-shipping the removed defect class.** An unvalidatable "MUST NOT close when X" decays
  exactly the way the old §15.6 preamble did — it survived 25 days (2026-08-04 → 2026-08-29)
  before being replaced. §15.6's own preamble (`COOPERATION.md:1349-1350`) is the warning label.
  Mitigation: P1 restricts the rule to disk state; the semantic residue is declared, not ruled.
- **Coverage failure.** A CLI-only evaluator leaves manual facilitation — the majority of the
  corpus — exactly as silent as today, reproducing the 3-of-80 `_index.md` outcome. A6 applies:
  a quiet CLI semantics change is worse than a protocol change, not better.
- **Floor-cost regression on `fast`.** If the trigger fires on a track that skips Phase 2 by
  design (`track.go:182`), it raises the floor cost the constraints forbid raising. Track-gating
  in P1 is the mitigation; it must be tested, not assumed.
- **Version-skew breakage.** A validator that fails closed on decks without the new prose breaks
  the 40-of-41 fleet on 2.10.0. Advisory-only degradation on old decks is a hard requirement,
  stated in the approach so it cannot be traded away later.
- **Always-fire / never-fire collapse.** A trigger that fires on every unanimous round-01 adds
  up to 28 rounds of cost against a lower-bound 16% no-movement rate (23/141) in opened rounds;
  one that never fires is a dead letter. The only ground truth for tuning is the conflicted
  party's own history (52/28) — usable as a prior, never as truth, and the replay must publish
  disagreements with it raw.
- **Measurement overreach.** Cost claims in tokens or currency are impossible today
  (`internal/driver/loop.go:174-175`); any success metric smuggling them in is unfalsifiable and
  should be struck on review — including from my own round-02.
