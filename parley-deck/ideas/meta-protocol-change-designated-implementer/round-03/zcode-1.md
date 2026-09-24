---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
round: 3
date: 2026-09-24
responding-to: [claude-1/round-02, kimi-1/round-02]
---

## Protocol context attestation

Obtained with `parley protocol packet --dir . --phase 2 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`.

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "fallback_reason": "absent (key not present in the rendered JSON)",
  "body_path": ".parley-runtime/protocol-packets/full-phase2-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"
}
```

`shasum -a 256` of `body_path` = `8ce83cde…9db7`, equal to `source_sha256`, `packet_sha256`,
the kickoff authority SHA and the live `parley-deck/COOPERATION.md` (109,928 B / 1,386 lines);
Phase 2/4/5, §9.0, §12.5 and §15 sections re-read from the body this round, remaining blocks
covered by the round-01/02 full readings of the same hash. The envelope's `shadow` block
(`31903915…`, 68,912 B, 36/69 blocks) describes the optimized packet I did not use.

HEAD at write time `73ee923` ("preserve cross-review and owner release direction"). Code and
protocol text remain untouched by the whole run: `git diff 59b0458 HEAD -- internal/
parley-deck/COOPERATION.md` is empty (round-01 baseline) and the working tree adds no code
changes, so every code locator re-verified below is PRIMARY at HEAD. Round-03 was dispatched
via the owner-authorized documented manual fallback (driver still halts on the unavailable
historical worktree `…/scratchpad/f2repo`; `organizer-notes.md` records it). I altered no
accounting, worktree history, budget state or code to route around it.

## Position changes since round 2

Four, and the first two are the substantive resolutions this round exists for.

1. **I withdraw Phase-0 materialization of the global default and adopt claude-1's C-2/S-4
   live layered read.** In round 2 I adopted materialization *citing claude-1's round-1 D3*;
   claude-1 has now withdrawn D3 with code evidence I re-verified this round, and I add a
   finding of my own that undercuts the precedent kimi-1 and I were relying on (see @kimi-1
   below): **§12.5's driver-authored `00-prompt.md` seeding contract is pipeline-scoped** —
   it governs `pipelines/<slug>/` block N+1 kickoff built from a manifest
   (`COOPERATION.md:1153-1155`), not ordinary `ideas/<slug>/` kickoff; for ordinary ideas no
   tool writes `00-prompt.md` today and §9.0 assigns readiness *recording* there to the
   facilitator, not to resolver machinery. With the precedent gone, the balance tips on
   claude-1's three arguments, all now verified: (a) `LoadDefaults` merges `[defaults]` across
   `configLayers` — central `~/.parley/agents.toml` → deck `parley-deck/agents.toml` →
   `agents.local.toml` → `$PARLEY_AGENT_CONFIG`, "later non-empty values win"
   (`internal/config/runtime.go:384-401`, `:415-437`) — and every shipped `[defaults]` key is
   resolved live at use; a materialized `default_implementer` would be the one key with
   bespoke write-then-freeze mechanics; (b) the freeze that matters (no horse-changing
   mid-implementation) is already supplied by the `IMPLEMENTATION.md` pin, which both
   resolvers read first; (c) materialization has asymmetric stale-default semantics: an
   in-flight idea never sees a later-set or later-removed global default (materialization is
   "Phase-0-only, idempotent"), so the owner edits a standing preference and only some ideas
   obey it. Live read gives the semantics the owner's words actually describe — "my current
   standing preference applies to any idea without its own designation" — and it deletes an
   entire mechanics surface (no writer, no idempotency rule, no dry-run carve-out, no
   never-overwrite-an-owner-field rule; kimi-1's open question 1 evaporates).
   **Resumability, answered explicitly:** resolution is recomputed at each driver
   construction (F-1, frozen per run); once Phase 5 writes `IMPLEMENTATION.md` the pin (tier
   1) governs every resume; before that, a changed global default behaves exactly like the
   pre-work metadata edit all three of us already permit (zcode-1 O2, kimi-1 item on mid-idea
   changes). One sentence in the amended text should say this — resolved per run, pinned by
   execution.
2. **The deck-level layer ships — as a corollary of (1), not a new debate.** My round-2 "no
   deck layer, named future extension" (and kimi-1's "ship two levels only") was only
   available *because* materialization could choose to read machine-only. Under the live
   read, `default_implementer` rides `configLayers` and per-idea > deck > machine arrives by
   construction; suppressing the deck layer would make this the one `[defaults]` key that
   ignores `parley-deck/agents.toml` (whose init header literally says a deck's agents.toml
   overrides the central file). I adopt claude-1's S-7 with his both-files test pin. Both
   files ship the key absent.
3. **I withdraw the write-end FINAL repair (my round-2 item 6) from this idea's delta.**
   Claude-1's C-3 withdrawal of D10(i) is mechanically correct and I checked the edge myself:
   today, on an undesignated deck, the drafter is `firstEligibleHeadlessAgent`'s first
   eligible headless participant (`driver_consensus.go:109-123`) while the fallback
   implementer is `eligible[0]` (`driver_impl.go:131-133`); these coincide when the first
   participant is headless and Found, and diverge exactly when they don't — in that edge, a
   new FINAL emitting `implementer: <drafter>` changes who implements on a deck that set
   nothing. That is a changed selection on the unset path, and my own round-2 rule — an
   unset-path behavior change requires explicit owner-visible deviation authority, and
   participant agreement is not that authority — binds me symmetrically when the group is
   2-vs-1 against the change. The repair is real value (the documented chain feeds the
   resolver nothing today, 51/85 historical FINAL.md files unreadable) and it joins the
   read-set widening in kimi-1's deferred follow-up slug
   (`meta-protocol-change-implementer-inheritance-repair`): one follow-up idea owns both
   halves of the legacy inheritance repair, decided visibly, not smuggled.
4. **Observability (`agent.implementer_resolved` + stdout source): I switch from
   designation-only to always-on**, adopting claude-1's S-8 primary proposal and matching
   kimi-1's table ("additive … no gating on either"). What moved me: (a) the precedent is
   exact — `agent.model_diversity` is an always-on durable event that gates nothing
   (`driver_impl.go:198-206`); (b) the stdout carrier already exists on every run —
   `driver: implementing via %s …` (`driver_impl.go:217`) prints today, designations or not,
   so annotating it with the source is a minimal delta to output that already fires; (c) the
   thing the event makes visible is the pre-existing defect this idea exists to fix
   (positional list-order selection, F-1/F-7) — shipping the fix with its visibility dark on
   the default path preserves the invisibility that let the accident persist. It must land in
   `## Agreed trade-offs` as the single declared unset-path addition (it observes; it selects
   nothing).

Two smaller withdrawals, recorded under the responses below: my Q5 actor assignment for pin
corrections (self-dealing hazard — see @claude-1 R/S-6) and my round-2 "pre-dispatch FINAL
validation" designation-only addition (redundant — claude-1's K-1 evidence, re-verified).

## Response to claude-1 (round-02)

**Accepted and credited:** C-1/X-1/X-2 (kimi-1's line numbers; I was one of the two wrong
parties and the lesson — agreement is not evidence, re-measurement is — is applied in this
file by re-verifying every locator I rely on); C-2 (adopted, position change 1); C-3 both
halves (adopted, position change 3 — including your D10(i) reversal, which my round-2 self
did not expect and your edge argument earns); X-3 (my `consensus.go:805` error — I retracted
it myself in round 2; your X-3 runs in my favor and I accept the strengthened reading);
X-5 confirmed-as-mine (below); X-6/X-7 (your inventory verified: `phasedigest.go:319` and
`protosnap.go:388` both read `implementer` today); S-6 (with one withdrawal of my own);
S-9/S-10/S-11/S-12/S-13 as written; R-D (gate scoped to the idea being run — I made your R-3
normative in round 2 and hold it); B-1 (see @kimi-1 — I read kimi-1's round-2 body as already
complying; the residual is a restatement slip, not a live proposal).

**On S-5, one refinement — split the unavailability rule by layer, and I ask you to confirm
or hold.** We converged in round 2 on fail-closed-with-a-pre-built-`Confirm` for the per-idea
designation (tier 2), and I keep that. But your S-5 text gates tier 3 (the global default)
the same way, and I now think your own C-2 cuts against that: once the global default is a
*live standing config* rather than a materialized per-idea act, an unavailable-but-rostered
global designee halts **every idea's kickoff** until someone records a per-idea waiver line —
a standing tax of one gate-service per idea per outage. Your LE-7/LE-11/L699 precedent
(`COOPERATION.md:275`, `:699`) is about a *per-run assigned* role (the goal-done checker);
the global default is standing preference, which is the case kimi-1 and I both assigned to
loud fall-through in round 2. Proposal: **tier 2 unavailable at the §9.0 ping → gate with
your `Confirm` (edit | waive | remove); tier 3 unavailable or not-a-participant → no gate,
loud record (resolved-source event + readiness-table/inbox notice) and fall through to
today's chain; invalid ids (non-participant, declared non-participating facilitator,
present-empty) stay hard gates on both tiers.** The fail-closed property you defended is
preserved where the owner's *specific intent about this idea* exists, and the standing layer
degrades loudly without training operators to wave gates through. If you hold S-5-uniform, I
will sign it — your `Confirm` makes the cost seconds, and the disagreement is a usability
judgment, not a safety one — but the split is my ask, and it is also the only reading under
which kimi-1's round-2 layering and yours compose without a veto.

**On S-3, a merge: keep your present-empty gate, add an explicit `none` opt-out.** Your S-3
gate on present-empty is right (verified: `ReadFrontmatter` does `strings.Cut(line, ":")`
then `meta[key] = value`, so `implementer:` yields a present `""` distinguishable from
absent, `workspace.go:390-394`; and `driver_impl.go:47-51`'s "NEVER silently falls back" is
the codebase's own value). But with tier 3 live, authors need a per-idea escape from the
owner's standing default — my round-2 three-state point survives the materialization
withdrawal in this form: **`implementer: none` (or `none` in any casing, trimmed) = explicit
non-designation — the idea opts out of tier 3 and resolves on today's chain.** Present-empty
stays a gate whose `Confirm` offers exactly the two legal spellings ("name a participant, or
write `none` to opt this idea out of the global default"). Parsing table, final: key absent →
tier 3 may fire; present-empty → fail-closed gate; `none` → opt-out, tier 3 suppressed;
`<id>` → designation, validated. All three of our positions are honored: your gate, kimi-1's
facilitator-symmetry in parsing (present-empty is a state, not absent), my opt-out. Your R-B
(the sibling-field asymmetry) is answered by the code comment you asked for, plus one more
sentence there: `facilitator:` is exclusionary (empty excludes nobody), `implementer:` is
appointive (empty appoints nobody — and now must say which of "nobody deliberately" vs
"nobody by mistake" is meant, which is exactly what `none` vs empty splits).

**On S-6, accepted — and I withdraw my Q5 actor.** My round-2 answer let "the incoming
implementer" log the correction edit. That is self-dealing: an incoming implementer
rewriting the pin that names them is the quiet self-preference this mechanism exists to
remove. Correct rule, merged from your S-6 and kimi-1's recorded-decision framing: the
machine escalates on any tier-1/tier-2(3) disagreement (never silently honors either); the
correction is an owner/author act recorded in frontmatter (deviation-style note naming the
abandoned attempt, or re-aligning the pin in the same recorded edit); the partial tree is
evidence and reassignment over it is a human decision.

**On S-9/X-5, your fix shape is right, with one addition.** Verified again at HEAD:
`newDriverImplOps` filters `eligible` via `role.IneligibleForRoles` and calls
`resolveImplementer(ideaDir, eligible)` (`driver_impl.go:52-64`) while consensus calls its
copy with the raw list it was handed (`consensus.go:730-748`). Unifying the *reader* alone
would freeze the divergence in a shared function. The shared package must export (or share
internally) **both** the reading chain *and* the eligibility filter, and the consensus call
site must pass the same facilitator-filtered list the driver does — otherwise we mint a third
divergence (two filters instead of two readers). Both terminal tails stay byte-pinned by
their existing tests (`app_test.go:1557/1576`, `roundgate_test.go:102`, both re-read).

**On S-8, adopted as always-on** (position change 4), with your `## Agreed trade-offs`
requirement and the conservative designation-only variant dropped — the group now reads 3-0
always-on (your proposal, kimi-1's table, my switch).

## Response to kimi-1 (round-02)

**Accepted and credited:** your corrected line numbers (now 2-of-2 rounds your locators have
been the accurate ones — noted again under the assignment question); position 1 (naming
concession with conditions — the present-empty condition is resolved by the S-3/`none` merge
above); position 3 (ValidateFinal redundancy — verified: `RequiredFinalSections` includes
`## Observable acceptance criteria` (`finalsections.go:18-26`) and the driver enforces
`finalScaffoldReason` before leaving Phase 4 (`internal/driver/consensus.go:55-63`); I also
drop my own round-2 "pre-dispatch FINAL validation" addition as doubly redundant); position
4's layering (adopted into the S-5 refinement above); the D10 deferral-with-slug (position
change 3 adopts your slug for both halves); #3 (`parley init` commented-out emission —
verified the existing `[defaults]` emission is active-keys-with-inline-comments at
`runtime.go:631-645`, so the commented line is a deliberate deviation FINAL must describe in
one sentence, exactly as my round-2 C-b put it); your §15.3 items 1 and 3 (no residual
conflict).

**Opposed: position 2 (materialization) — please switch to the live read.** Three grounds:
(a) your §12.5 precedent does not reach ordinary ideas — the seeding contract is
pipeline-scoped (`COOPERATION.md:1153-1155`, verified: "the driver authors block N+1's
`00-prompt.md`… built from the manifest… the prior block's finalized typed artifact" —
§12 is "Pipeline blocks & action stages"); for `ideas/<slug>/` no tool writes `00-prompt.md`
today and §9.0's readiness recording is the facilitator's act; (b) claude-1's C-2 mechanics
evidence (re-verified, cited above); (c) your own open question 1 (materialization actor and
timing, idempotency, never-overwrite, never-in-flight) is the smell — a mechanism whose
safest description needs four guard clauses is weaker than the one that needs none. Under the
live read your O1 stays settled (global > claim by tier order, not by write order), your
inapplicability fall-through survives unchanged (tier 3 id not in `participants:` → notice,
fall through), and every mechanics question dissolves. Your correction-path framing for
re-entry (recorded decision naming the abandoned attempt) is adopted into the merged S-6.

**One inconsistency to close (this is claude-1's B-1, and I read it as a slip, not a
proposal):** your round-2 body withdraws machine claim handling three separate ways ("no
machine honoring at all", "social act a human/organizer may materialize", the boundary table's
"no parsing, no gate, no event" / "subordination sentence only") — but your "Current proposal
(concise restatement)" still lists "**perfected uncontested claim (only when no field)**" as a
precedence tier between the per-idea field and FINAL inheritance. A tier in the resolution
chain is a machine ranking; it cannot coexist with "no parsing". Please confirm in your
round-03 (or in your consensus signoff) that the body governs and the restatement tier is a
drafting leftover. The FINAL consequence either way: **no claim tier in the machine chain;
the claim stays a social act; on the designation path a contradicting claim is surfaced as
advisory information (never honored, never gating); on the unset path nothing changes at
all** (the `impl-claim` file stays inert — grep re-run this round, still zero Go readers).
With that confirmation, claude-1's B-1 is satisfied without a fourth exchange; if you
actually intend a machine claim tier, that is a genuine block and we should spend our last
cross-review on it — but nothing else in your file supports that reading.

**On the trace-of-prior-fixes claim (§15.3 item 2): I accept your split verdict as the
admissible resolution.** The counts (1 CRITICAL + 8 MAJOR) are CONFIRMED by all of us with
PRIMARY evidence; "all fixed in later cycles" has closure evidence (three fix-up cycles,
`status: complete`, a zero-fix closing consensus) but no per-finding trace; the admissible
verdict for the clause is **closure observed, per-finding trace not performed** — my round-1
UNVERIFIED stands for the trace half, claude-1's quoted finding-headings do not entail it,
and your decomposition is the one FINAL and the post-release recommendation may cite. Nothing
in any acceptance criterion depends on the clause.

**On your readiness condition 2:** satisfied — I accept fail-closed-with-recorded-waiver for
the per-idea designee (I adopted it in round 2, position 1, and keep it), with the global
layer's loud fall-through now carried inside the merged S-5 refinement rather than as my
round-2's automatic tier-4 descent.

## Resolved decisions (the seven live items)

1. **Global default mechanics: live layered read, no materialization.** Claude-1 S-4
   adopted by all three (my withdrawal is position change 1; kimi-1 asked to confirm).
   `[defaults].default_implementer` resolved through `LoadDefaults` at each run start,
   UNSET on ship, emitted commented-out by `parley init` (deliberate deviation from the
   active-key emission precedent, one sentence in FINAL). Deck layer inherited by
   construction: per-idea > deck > machine, pinned by the both-files test. Non-rostered
   global designee → notice + fall through (unanimous from round 2). Unavailable global
   designee → the S-5 layer split above (pending claude-1's confirmation; not a block).
2. **Per-idea empty/none semantics: four states.** absent → tier 3 may fire;
   present-empty → fail-closed gate (Confirm names both legal spellings); `none` → explicit
   non-designation, tier 3 suppressed, today's chain; `<id>` → designation, validated
   fail-closed. Implementable as verified (`ReadFrontmatter` preserves presence;
   `FacilitatorRoleFromMeta` at `facilitator.go:43-53` is the parsing mirror for the
   empty-is-a-state part).
3. **Write-end FINAL repair: OUT of this delta**, deferred with the read-set widening to
   the recorded follow-up slug. Claude-1 C-3 final and kimi-1 deferral now 3-0; my round-2
   lone hold is withdrawn (position change 3). This idea's resolver FINAL.md read set stays
   exactly `{implementer, drafted-by}`.
4. **Claims: no machine tier, ever, in this idea.** kimi-1's restatement tier flagged as a
   slip (body governs — confirmation requested, not blocking); claude-1's B-1 satisfied by
   the withdrawal; advisory surfacing on the designation path only; unset path untouched.
5. **Observability: always-on** `agent.implementer_resolved {idea, implementer, source}` +
   stdout source beside the existing `implementing via` line; `source ∈
   {implementation-md, designation, none, global-default, final-md, positional}`; recorded
   as the single declared unset-path addition in `## Agreed trade-offs`. 3-0.
6. **Re-entry pin conflicts:** machine escalation on disagreement (never silent pin-win);
   correction = owner/author recorded frontmatter act naming the abandoned attempt; the
   implementer never edits a pin that would name themselves. Merged S-6; unanimous material.
7. **This-run assignment (today's protocol, not the owner's future default): I accept
   claude-1 drafter / kimi-1 implementer, and withdraw my zcode-1-implements proposal.**
   Claude-1's two evidence-backed reasons hold and my own round-2 evidence does not outrank
   them: (a) the post-release recommendation needs implementation data points, and
   kimi-1-implementing produces the roster's second (zcode-1 has the first, n=1) instead of
   deepening the same one; (b) across two rounds of dense locator work, kimi-1's line
   numbers are the only set never wrong (X-1/X-2 round 2; this round again) — small sample,
   stated as such, but it is the only directly relevant evidence and this delta spans six
   `implementer` readers across five packages plus three protocol copies; (c) rotation:
   lean-organizer ran kimi-drafts/zcode-implements/claude-reviews; this runs
   claude-drafts/kimi-implements/zcode-reviews, moving every seat once, which is the
   standing-concentration mitigation all three of us wrote down. I hold my standing
   willingness for either role if kimi-1 becomes unavailable or prefers otherwise —
   willingness is not a claim. Perfection happens through today's channels only: `Drafter:
   yes` in claude-1's consensus signoff block; kimi-1's
   `inbox/kimi-1-to-all_<slug>_impl-claim.md` before Phase-5 work begins. Reviewers = the
   two non-implementers (claude-1 and me). This is a this-run assignment under the protocol
   in force; nothing here names or nudges the owner's future global default, which ships
   UNSET and is decided post-release on evidence.

## The mechanism FINAL should specify (deterministic restatement)

Resolution, in order, recomputed at each run start, frozen for that run, executed by one
shared reader + one shared eligibility filter in `internal/protocol`:

1. `IMPLEMENTATION.md` `implementer:` (re-entry pin) — must be a participant; disagreement
   with tiers 2/3 escalates (never silent).
2. `00-prompt.md` `implementer:`: `<id>` (validated: participant, not the declared
   non-participating facilitator — else blocking gate) or `none` (opt-out; skip tier 3);
   present-empty = blocking gate.
3. `[defaults].default_implementer` (deck > machine, live): if the id is a participant →
   designation; if not rostered on this idea → notice + fall through; if unavailable at the
   §9.0 ping → notice + fall through (layer split, pending claude-1 confirm).
4. `FINAL.md` `{implementer, drafted-by}` — read set unchanged, participant-validated.
5. Tails, byte-pinned by existing tests: driver `eligible[0]`; consensus `""` → expect
   everyone.

Tier 2 unavailability at the §9.0 ping = blocking gate with pre-built `Confirm` (re-designate
| `implementer_waived: <id> — <reason> — confirmed <date>` | remove), the §9.0 record shape
(`COOPERATION.md:872-874` form). Mid-Phase-5 death = escalate; partial tree is evidence.
Designation-only additions (all inert without a designation, per the task's
classification requirement): drafter-preference in `firstEligibleHeadlessAgent` (prefer a
non-designee, fail-open to the designee when none exists, degenerate case pinned); Phase-0
relocation of the model-diversity check and the zero-non-implementer-reviewer floor
(`driver_impl.go:282` already hard-stops at Phase 6 — only timing moves); two-participant
kickoff warning (not gate; `COOPERATION.md:746` keeps two-participant runs satisfiable); one
L433 sentence adding designated runs to the self-containment trigger; §15.5-style one-line
concentration record when drafter == implementer. New gates read live idea frontmatter
(never `meta/version.json` — F-8) and are scoped to the idea being launched. All three
COOPERATION.md copies change identically + `meta/protocol-changelog.md`; SKILL.md restates
no Phase-5 rule (verified round 2) with an implementation-time grep check of the skill tree.
Excluded to the follow-up register: read-set widening, write-end FINAL emission, claim
parsing/auto-honoring, trajectory cross-gate (claude-1 K-7 adopted — my round-2 scoped
version was hygiene, not function, and couples this delta to another subsystem),
`--participants` order documentation, `facilitatorConflictGates` all-ideas iteration repair.

## Owner boundary ledger (updated)

*"Ship the mechanism with the global default UNSET, so behaviour without it is exactly
today's."* **Preserved by construction on a deck that sets nothing:** resolver read set
unchanged; both tails and their test pins byte-preserved through the shared-reader refactor;
no new gates, events-with-teeth, or outputs on undesignated ideas; the `impl-claim` file
inert; roles advisory; quorum/roster/signoff untouched; every attended gate untouched
(`parley protocol publish` TTY, core publication owner-attended, channels at the organizer's
release step — S-12's "designation names who executes, never which gate applies" carried
verbatim into FINAL). **Designation-only:** everything in the two paragraphs above.
**The one declared unset-path addition:** the always-on `agent.implementer_resolved` event +
stdout source (observes, selects nothing) — `## Agreed trade-offs` entry mandatory. **Nothing
else touches the unset path** — with the write-end repair withdrawn, my round-2 ledger's only
trade-off line is gone and the boundary is cleaner than any of our round-1/2 proposals.

## User direction (release order and Windows)

Copied verbatim to `source-context/release-order.md` (Slovak original preserved there);
English direction, quoted:

> "Both — release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, label
> Windows as experimental, and hold winget for the CLI. At the same time open a separate
> reviewed Windows idea that fixes the defects and removes the label."

And the ordering: release-1.49.1 writes its done file first, then this idea, then
`windows-portability`, each gated on its predecessor's done file. Status at my write time:
`codex-1-to-user_release-1.49.1_done.md` **exists** in the lean-organizer worktree deck
inbox, so this idea's release is next in the owner's queue once reviewed implementation
lands — the wait condition is already satisfied, nothing for participants to decide there.
Consequences carried into FINAL/delivery: CLI release notes label Windows
experimental/unvalidated, Windows assets retained, **no CLI winget PR**, skill channels
unaffected, no Windows platform architecture work in this idea (it belongs to
`windows-portability`; both claude-1 and kimi-1's Windows assessments already said this delta
adds no Windows surface — no build tags near the touched files — and the owner decision
confirms rather than widens scope). Versions chosen above actual releases at staging time.
This direction is already authorized; I raise no owner question against it.

## Existing alternatives

Per §15.6(a), for elements whose shape changed in this round. Unchanged elements keep my
round-2 table (8 entries, all locators re-verified still valid at HEAD).

| Element changed in round 3 | What already ships | Locator | Why the round-3 shape is right |
|---|---|---|---|
| Global default resolution (live read, was materialization) | `LoadDefaults` merging `[defaults]` across `configLayers` live at use for every existing key; deck-overrides-central init header | `internal/config/runtime.go:384-401`, `:415-437`, `:631-645` | Materialization had a pipeline-only precedent (§12.5, `COOPERATION.md:1153-1155`), bespoke mechanics, and stale-default asymmetry; the live read is the shipped pattern for all sibling keys and needs zero write machinery. |
| Deck layer (shipped, was deferred) | Deck `parley-deck/agents.toml` overriding central `~/.parley/agents.toml` is an existing shipped layer | `runtime.go:384-401` | Under the live read the third level arrives by construction; suppressing it would be the one `[defaults]` key ignoring the deck — the surprising special case, not the safe default. |
| Write-end FINAL emission (withdrawn to follow-up) | Resolver reads `implementer`/`drafted-by` from FINAL.md; no tool template emits them (51/85 historical files unreadable) | `driver_impl.go:117-118`; `consensus.go:764-765`; `COOPERATION.md:413-419` template prescribes `author:` | Emission changes selection in the not-headless-first edge on decks that set nothing — outside "exactly today's"; the defect is real and recorded, owned by the follow-up slug with the read-set widening. |
| Present-empty/`none` parsing (four states) | `FacilitatorRoleFromMeta` empty-equals-undeclared; `ReadFrontmatter` preserves key presence; `gate{Kind,Detail,Confirm}` with one-line-edit confirms | `facilitator.go:43-53`; `workspace.go:390-394`; `preflight.go:337-341` | The gate catches the empty-key mistake fail-closed; `none` gives authors the per-idea opt-out a live global tier requires; both spellings named in the gate's own `Confirm`. |
| Unavailability handling (layer split) | §9.0 exclusion records (`<id> — reason — confirmed <date>`); LE-7/LE-11 halt for per-run roles; readiness ping table | `COOPERATION.md:868-877`, `:275`, `:699` | Per-idea designation is an act about this idea → gate; standing config is not → loud fall-through; the precedent roles that halt are per-run assignments, not standing preferences. |
| Always-on resolved-source event | `agent.model_diversity` always-on non-gating event + `implementing via %s` stdout that already fires every run | `driver_impl.go:198-206`, `:217` | Same shape as the shipped observability precedent; makes the positional-selection accident visible on the exact path where it happens today. |

## Concerns / open questions

1. **S-5 layer split needs claude-1's explicit confirmation** — the only design item where
   I am asking a peer to move rather than recording a merge. My fallback if he holds
   uniform gating: I sign it (stated above); the item is a usability judgment, and FINAL
   should then carry his S-5 verbatim.
2. **kimi-1's restatement claim tier** — confirmation requested that the body governs.
   Treated as a slip; if it is not, that is the one genuine block left and the last
   cross-review (3 of 3) should take it.
3. **`--no-implement` visibility at preflight** (claude-1's item 3, still open): the
   availability gate should not fire for runs that will not reach Phase 5. One
   implementation-time trace; if not visible, gate + waiver line is the exit. Agree with
   claude-1's framing; no design change needed for it.
4. **"Implementation is the heaviest token work of a run" stays UNVERIFIED** (§15.2);
   the 30-40%-of-implementing-organizer-input figure may appear with its caveat, the
   ranking may not. Unchanged, three-way settled.
5. **New, small: the resolved-source enum should record `none` and `global-default`
   distinctly** so an opt-out and an unset deck are never confused in the event stream;
   one constant, worth naming in FINAL so the event schema is stable from day one.

## Readiness for consensus

**Ready.** No substantive block from my side. With this file's four position changes, the
three proposals are compatible on every shipped element: S-1/S-2/S-4 (live read)/S-6/S-7/
S-8 (always-on)/S-9 (+shared filter)/S-10/S-11/S-12/S-13, the four-state parsing, the layer
split (pending one confirmation with a stated fallback), no claim tier, both FINAL repairs
deferred to the named slug, and the assignment accepted. The required recorded decisions for
FINAL are exactly two `## Agreed trade-offs` entries — the always-on event, and the
commented-out init emission as a deliberate deviation — plus the follow-up register. I see
no need for cross-review 3 of 3 unless kimi-1's restatement tier was not a slip or
claude-1 holds S-5-uniform *and* kimi-1 holds the opposite — a conjunction I assess as
unlikely; otherwise the organizer may open consensus after round-03 responses land.

## Process notes

**Evidence ownership (§15).** New verdicts this round, all PRIMARY from re-measurement at
HEAD `73ee923` in this worktree: claude-1's C-2/K-5 (`runtime.go:384-401`, `:415-437`,
re-read in full), K-1 (`finalsections.go:18-26`, `:92`; `internal/driver/consensus.go:55-63`),
S-3 implementability (`workspace.go:390-394`; `facilitator.go:43-53`), X-7 reader inventory
(`phasedigest.go:319`, `protosnap.go:388`, trajectory flag re-read at `trajectory.go:65`),
S-8 precedent (`driver_impl.go:198-206`, `:217`), R-D/preflight shape
(`preflight.go:314-340`, gate-with-`Confirm` verified, all-ideas iteration verified),
test pins (`app_test.go:1557/1576`, `roundgate_test.go:102-107`), protocol lines L234/L275/
L538/L699/L746 re-read, `impl-claim` grep re-run (still zero Go readers). My new
first-canonical claim this round: **§12.5's seeding contract is pipeline-scoped and is not
precedent for materialization on ordinary ideas** (`COOPERATION.md:1150-1155`, read in the
packet body) — claude-1 and kimi-1 own the verdict on it. The §12.5 scoping also slightly
corrects my own round-2 new-concern 1, which cited §9.0-readiness-recording as kickoff-write
precedent without checking §12.5's scope — recorded here as a correction to a claim I own.
X-1/X-2: I was wrong alongside claude-1 and kimi-1 was right; no verdict issued on my own
error beyond the correction. The trace-of-prior-fixes split (kimi-1 §15.3.2) is accepted as
the group resolution; I own no verdict on my own round-1 UNVERIFIED beyond adopting the
decomposition.

**Independence and scope.** I read `00-prompt.md`, `source-context/owner-brief.md`,
`source-context/release-order.md`, the full packet body (hash-verified against the live
authority), all three round-02 artifacts, and `organizer-notes.md` before writing. I wrote
only this file. No implementation, no commit, no publication, no channel action, no edit to
any peer file or to closed records; roster identities not spawned; no accounting, worktree
or budget state altered. The assignment acceptance in item 7 is a proposal acceptance under
today's protocol, perfected only through the prescribed acts (`Drafter: yes` signoff;
impl-claim inbox file), and is not a recommendation for the owner's future global default —
that ships UNSET and is decided post-release on evidence.

<!-- Original language of the owner's quoted direction: Slovak, preserved verbatim in source-context/release-order.md; English translation quoted in the User direction section above. -->

Round-03 is cross-review 2 of 3 under the deliberation cap.
