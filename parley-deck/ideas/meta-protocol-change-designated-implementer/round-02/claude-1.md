---
agent: claude-1
idea: meta-protocol-change-designated-implementer
round: 2
date: 2026-09-24
responding-to: [kimi-1, zcode-1]
---

## Protocol context attestation

Obtained with `parley protocol packet --dir . --phase 2 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`.

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "fallback_reason": null,
  "body_path": ".parley-runtime/protocol-packets/full-phase2-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"
}
```

`fallback_reason` is **absent** from the JSON (not null-valued); rendered `null` above. The renderer
did not refuse. `shasum -a 256` of `body_path` returns `8ce83cde…`, equal to `source_sha256`,
`packet_sha256`, the kickoff authority SHA, and `shasum -a 256 parley-deck/COOPERATION.md` — the
packet is the live authority verbatim, 109,928 bytes / 1,386 lines, read in full. The envelope again
carries a `shadow` block for the packet that was **not** used (`31903915…`, 68,912 B, 36 of 69
blocks); the shadow differs from the phase-1 shadow (`7ec29a81…`, 68,811 B) because the phase flag
changed, which is expected and is not a drift signal about the source.

HEAD is now `ace284d6b31b72bb1223443f9397a1547690d1b5` ("preserve independent analyses and open
cross-review"), one commit above the `59b0458` that all three of us verified against in round 1.
`git diff --stat 59b0458 HEAD -- internal/ parley-deck/COOPERATION.md` is **empty**: the code and the
protocol text are unchanged, so every round-1 code locator remains addressable at HEAD and all three
round-1 verification records stay valid as to content. Only deck artifacts moved.

**Driver gap.** I did not launch myself; the organizer's `continue --auto` halted before round 2 on
the unavailable historical worktree `…/scratchpad/f2repo`, recorded in `organizer-notes.md`. I write
this artifact under that documented manual fallback and have altered no accounting, worktree history,
budget state or code to route around it.

---

## Corrections to my own round-01

Three, and the first is the most useful thing in this file.

**C-1 — my resolver line numbers were wrong, and zcode-1's were wrong the same way.** I wrote the
positional fallback as `driver_impl.go:129-131` and the consensus `return ""` as `consensus.go:775`.
zcode-1 wrote `129-131` and `771-773`. kimi-1 wrote `131-133` and `778`. **kimi-1 is correct.**
Measured at HEAD with `grep -n` and `awk` line numbering on the unmodified files:

```
internal/app/driver_impl.go
  104:  func resolveImplementer(ideaDir string, participants []string) string {
  131:  	if len(participants) > 0 {
  132:  		return participants[0]
  133:  	}
  134:  	return ""
  135:  }

internal/consensus/consensus.go
  751:  func resolveImplementer(ideaDir string, participants []string) string {
  778:  	return ""
  779:  }
```

I flag this beyond the pedantry because it is a live demonstration of the §15 rule this cross-review
is being asked to apply: **zcode-1 and I agreed, and we were both wrong.** Two independent artifacts
converging is not evidence; re-measurement is. The only reason this surfaced is that the organizer
asked for contradictions to be resolved rather than averaged.

**C-2 — I withdraw D3's Phase-0 materialization of the global default into `00-prompt.md`.** My
round-1 argument was that materializing the default at kickoff collapses the precedence question. It
does, but the evidence says the cost is wrong. `[defaults]` is already a **live, layered** read:
`LoadDefaults` (`internal/config/runtime.go:419-437`) walks `configLayers`
(`runtime.go:384-401`) — central `~/.parley/agents.toml`, then deck `parley-deck/agents.toml`, then
`agents.local.toml`, then `$PARLEY_AGENT_CONFIG` — with "later non-empty values win"
(`runtime.go:415-418`). Every existing `[defaults]` key (`ping_tier`, `preferred_transport`,
`roster_change_policy`, `speed`, `timeouts`, `loop`) is resolved that way, live, at use. A
`default_implementer` that alone got written into a canonical artifact at Phase 0 would be the single
key in that block with different mechanics, and it would make the driver a writer of `00-prompt.md`
— an artifact the protocol assigns to the idea's author, which no tool writes today. The freeze that
actually matters (nobody changes horses mid-implementation) is already supplied by the
`IMPLEMENTATION.md` pin, which both resolvers read first. **Revised position: live layered read,
with the resolved id and its source recorded, not materialized.** kimi-1's O1 then has to be answered
on its merits rather than dissolved; I answer it in R-2 below.

**C-3 — I withdraw D10(ii), the FINAL.md key-set widening, and also D10(i), from this idea's
scope.** Round 1 argued for widening the resolver's FINAL.md key set to
`implementer, drafted-by, author, drafter, finalized-by` and for making the driver's FINAL-draft
prompt emit `implementer:`. I named (ii) as a real behaviour change and argued in-scope anyway. Held
against the owner's exact words — *"Ship the mechanism with the global default UNSET, so behaviour
without it is exactly today's"* — both halves fail, and (i) fails for a reason I did not state in
round 1: a deck with no designation that drafts a new FINAL would start emitting `implementer:
<drafter>`, after which the resolver would select the drafter where it previously selected
`participants[0]`. That is a changed selection on the unset path. Scope creep that changes who
executes is exactly the thing the boundary exists to stop, and my round-1 self-justification
("shipping the override on top of a rule that does not run is shipping half") is an argument for a
*second idea*, not for widening this one. See the follow-up register in the boundary ledger.

---

## Contradictory factual verdicts, resolved

Participant agreement is not evidence (§15). Each row was re-measured at HEAD `ace284d6` in this
worktree; every verdict below is `PRIMARY`.

| # | The disagreement | Resolution at HEAD |
|---|---|---|
| X-1 | Positional fallback location: claude-1 `129-131`, zcode-1 `129-131`, kimi-1 `131-133` | **kimi-1 correct.** `driver_impl.go:131-133`, `return ""` at `:134`. See C-1. |
| X-2 | Consensus `return ""`: claude-1 `:775`, zcode-1 `:771-773`, kimi-1 `:778` | **kimi-1 correct.** `consensus.go:778`. |
| X-3 | zcode-1: *"The two copies agree in practice only because the tool-written FINAL template emits `drafted-by` (`consensus.go:805`)."* vs claude-1 F-2/F-3 and kimi-1 A2: no tool writes a readable key into FINAL.md | **zcode-1's locator and conclusion are both WRONG.** `consensus.go:805` is inside `draftTemplate` (func at `:781`), which emits `review-cycle:`, `outstanding_agreed_fixes:`, `blocked:`, `## Agreed fixes`, `## Deferred follow-ups`, `## Dismissed findings` — that is the **Phase-7 review-consensus scaffold**, not FINAL.md. `consensus.go:833` is `designDraftTemplate` (func at `:831`) — the **Phase-3 `consensus.md`** scaffold. `phase58.go:397` is the review-consensus *prompt*. `grep '"FINAL.md"'` across non-test `internal/` returns readers only; **no code writes a FINAL.md template at all** — `buildFinalDraftPrompt` (`driver_consensus.go:174-199`) asks an agent to write it and requires only `idea:` and `status:`. So the copies do not agree "because the template emits `drafted-by`"; they agree on the happy path only when a *human-or-agent-authored* FINAL happens to carry one of the two keys, which 51 of 85 FINAL.md files on this deck do not. This strengthens kimi-1's A2 and my F-2/F-3 rather than qualifying them. |
| X-4 | Brief hypothesis "with `author: user` the implementer is decided by round-01 filing speed": claude-1 `WRONG` (for the driven path), zcode-1 `CONFIRMED` (summary) / `CONFIRMED with qualifications` (table), kimi-1 `PARTIALLY CONFIRMED` | **All three of us measured the same thing and graded it differently; the facts do not conflict.** The prose chain (L407→L409→L443) does say first-filer-drafts and drafter-implements. No code computes a first filer (`grep` for any filing-time comparison returns nothing in `internal/`). The driven path is list order in two different functions. The honest single verdict: **the hypothesis is CONFIRMED of the protocol text and WRONG of the shipped tooling**, and the gap between those two is the actual finding. zcode-1's *summary* sentence ("I verified … is indeed decided by round-01 filing order") overstates zcode-1's own table, which already carries the qualification; I read that as a summary/table mismatch in one artifact, not a third position. |
| X-5 | zcode-1's second divergence axis: driver validates against facilitator-filtered `eligible`, consensus against raw `participants` | **CONFIRMED, and it is zcode-1's find, not mine — I missed it.** `newDriverImplOps` builds `eligible` by dropping `role.IneligibleForRoles(p)` and calls `resolveImplementer(ideaDir, eligible)` (`driver_impl.go:52-64`); `expectedRoundParticipants` (`consensus.go:730`) calls its copy with the raw `participants` it was handed (`consensus.go:734`). In a declared-facilitator run whose FINAL.md names the facilitator, the two resolve differently from identical inputs. This run is a declared-facilitator run, so the axis is live here. |
| X-6 | kimi-1 V14: the third COOPERATION.md copy is UNVERIFIED | **Supplied.** I measured all three in round 1 (F-9): `parley-deck/COOPERATION.md` `8ce83cde…` 109,928 B; `internal/protocol/defaults/COOPERATION.md` `cce5d7d9…`, 16 diff lines all project zones; `…-skill/skills/parley-deck/references/COOPERATION.md` `fc907e59…` 109,772 B, 18 diff lines of the same zones, byte-identical to the published core 2.13.0 named in `source-context/prior-release-handoff.md`. Phase 4 rule at L409/L402/L402, Phase 5 rule at L443/L436/L436. This is my claim, first asserted in round 1; kimi-1 and zcode-1 own the verdict on it, not me. |
| X-7 | kimi-1 A3: `parley trajectory configure --implementer` is a fourth `implementer` surface | **CONFIRMED** — `internal/app/trajectory.go:65` (`implementer := f.String("implementer", …, "frozen patch author")`), requires `ReadChecksContract` to return a non-empty named-checks list and an attended `--yes` (`trajectory.go:105-125`), feeds `trajectory.NewPolicy`. Never consulted by `newDriverImplOps`. kimi-1 is right that it is not a selector. I missed it in round 1; it belongs in the surface inventory. Two further readers neither of us listed: `internal/driver/phasedigest.go:319` (`meta["implementer"]` off IMPLEMENTATION.md) and `internal/tui/protosnap.go:388` (`{"agent","implementer","by"}`). The `implementer` key therefore has **six** readers across five packages today. |

---

## Response to zcode-1

**Agreed, adopted, and credited:** the `facilitator.go` instantiation shape; role scope binding
Phases 5–8 without touching drafting, quorum or signoff weight; `[defaults].default_implementer`
shipped UNSET; the re-entry pin ranked first; preflight failure for a designation naming a
non-participant or naming the declared non-participating facilitator; the O2 rule that pre-work
designation edits are ordinary metadata edits while post-`IMPLEMENTATION.md` edits need a recorded
decision; **X-5**, which is a real divergence I did not find; and the O5-adjacent warning that
SKILL.md prose restatements can drift independently of the three guarded COOPERATION.md copies.

**R-1 — Your "loud fallback" for an unavailable designee: I argue fail-closed-with-a-confirm, and I
think the disagreement is narrower than it reads.** Your design already fails closed on the two
*invalid* cases; the disputed case is only the third — a valid participant who fails the §9.0
liveness ping — where you fall through to tier 4 and I halt. Your reasoning is that facilitator
purity is a *safety* property while implementer designation is a *preference*, so escalation is
disproportionate. Two pieces of evidence cut against that framing:

1. The protocol already escalates on exactly this shape for a role that is also "just" an execution
   assignment. LE-7/LE-11 (`COOPERATION.md:275`): *"a missing, self or **unavailable** checker, a
   failed run, or an inconclusive or reserved verdict leaves completion unverified and escalates."*
   And L699: *"unavailable checker now halts a review-clean close until a human restores an
   independent [checker]."* The goal-done checker is not a safety-critical identity either; the
   protocol halts anyway, because a silent substitution destroys the property the role was created
   for.
2. Falling through to tier 4 does not land on "today's protocol" in the driven path. It lands on
   `participants[0]` (`driver_impl.go:131-133`), because tier 4's FINAL.md keys are usually
   unreadable (X-3). So the loud fallback's destination is **precisely the accidental list-order
   selection the owner is asking to replace**. "Loud" mitigates the surprise but not the outcome.

**Counterproposal that I think gets you your speed argument at full strength.** Do not halt with no
exit — halt with the exit already built. `gate` is `{Kind, Detail, Confirm}` (`preflight.go:337-341`),
and `facilitatorConflictGates` already ships a `Confirm` string that tells the operator the exact
one-line edit. So: an unavailable designee produces
`gate{Kind: "implementer-designation", Detail: "<id> designated but unavailable at the §9.0 ping",
Confirm: "edit ideas/<slug>/00-prompt.md — designate an available participant, or record
implementer_waived: <id> — <reason> — confirmed <date> to fall back to the default chain"}`. That is
one edit, seconds of owner time, and it produces the §9.0-shaped record
(`COOPERATION.md:872-874` uses exactly that `<id> — reason — confirmed <date>` form for exclusions).
The owner never waits on an agent, the fallback is available on demand, and the fallback is
**written down** instead of inferred. If you still prefer automatic fall-through I will not block
consensus over it — but then the fallback must write the same waiver line itself, because an
un-recorded reassignment is the failure mode both of our designs exist to close.

**R-2 — Your global-default > claim ordering: agreed, and I now think it is required, not
preferred.** In round 1 I sidestepped this with materialization; C-2 withdrew that, so it needs a
real answer, and yours is right. If a standing owner designation ranked below a participant's inbox
claim, any participant could defeat a standing owner instruction by filing a file first — restoring
exactly the whoever-moves-first selection the owner is removing. kimi-1's O1 alternative loses on
that ground. kimi-1's underlying worry is still legitimate and is answered without inverting the
order: under a live designation a claim is **recorded and surfaced as dissent**, and the response is
your own O2 one-line per-idea edit. The claim is not silenced; it just does not self-execute.

**R-3 — Your `consensus.go:805` claim is wrong (X-3), and the correction runs in your favour.** The
tool-written-template premise was load-bearing for "the two copies agree in practice"; with it
removed, the copies do not reliably agree at Phase-5 start at all, which strengthens your §8 call for
unification rather than weakening it.

**R-4 — Field naming: I hold `implementer:`; your lean and mine match, and I have an argument neither
of us made.** A distinct `designated_implementer:` has one concrete cost: an owner who types the
obvious key — `implementer:` — into `00-prompt.md` gets **silence**. The designation does nothing,
no gate fires, and the run proceeds on list order. That is a new silent-failure mode introduced by
the very change meant to remove one. With `implementer:` the obvious key is the real key. kimi-1's
greppability concern is real and is answered by the recorded resolution source (below), which
separates intent from outcome by *field*, not by key name.

**R-5 — One place your §7 wording needs tightening.** You write that the driver should resolve the
designation "at the top of the existing chain". It must be **second**, below the `IMPLEMENTATION.md`
pin, or a mid-run `00-prompt.md` edit reassigns an implementation already in flight. Your §2 ladder
has this right; the §7 sentence contradicts it. Flagging so the FINAL does not inherit the looser
phrasing.

---

## Response to kimi-1

**Agreed, adopted, and credited:** the corrected line numbers (C-1, X-1, X-2); A3 and the resolver
inventory (X-7); A1/A2 (independently re-measured, and X-3 strengthens A2); A4 (preflight already has
the gate machinery); A5 (the MIN-7/F15 drift precedent — which is mine and which you are right to
cite against a two-copy design); V12's operational corroboration from the lean-organizer
`organizer-notes.md` that the driver "printed 'drafting FINAL via claude-1' despite Kimi's published
FINAL and claim". **That last one closes my own R-6**: I declared in round 1 that my list-order
finding was a static reading with no observed launch. Your V12 is the observation. I am the owner of
F-1 so I issue no verdict on it, but I record that the evidence I said was missing now exists and was
supplied by a non-owner.

**K-1 — Your dimension-4 pre-dispatch FINAL validation is already shipped; proposing it again is
either redundant or a boundary breach.** You propose *"the driver should validate FINAL's required
sections and observable acceptance criteria before Phase-5 dispatch and fail closed on gaps"*. At
HEAD: `protocol.ValidateFinal` (`internal/protocol/finalsections.go:92`) is documented as *"the ONE
gate for a FINAL.md, used by manual finalization and by the driver"*; it checks `status: final`, the
idea slug, then `FinalIsScaffold` → `MissingFinalSections` over `RequiredFinalSections`
(`finalsections.go:18-26`), which **already contains `## Observable acceptance criteria`**, plus a
250-byte floor and the `FinalScaffoldPlaceholders` token scan. It is wired into the driver's
consensus advance at `internal/driver/consensus.go:57` and `:61` via `finalScaffoldReason`
(`:175-190`), i.e. before the idea can leave Phase 4. So the gate you want exists, fires pre-dispatch,
and fails closed. The only delta left would be judging the *substance* of those sections — which a
validator cannot do and which, if attempted, would change the gate for every deck including decks
with no designation. That is a default-path behaviour change and it is out of bounds. Substance is a
review duty (the drafter's signoff and the Phase-6 reviewers), not a validator's. **Recommend: cite
the shipped gate in FINAL, add no new validation.**

**K-2 — Your 7(b)/(6) automatic claim parsing is the single largest boundary risk in any of the three
round-1 artifacts, and I ask you to withdraw it.** You propose the resolver gain a claim tier and the
driver "honor a perfected claim by resolving it into the same durable channel". Consider the deck
that sets nothing: today an `impl-claim.md` file is inert (F-5/A1 — `grep -rn "impl-claim"
--include="*.go" .` is empty). With claim parsing, that same deck resolves to a different agent than
it does today. That is *behaviour without the mechanism set*, changed. It fails the owner's sentence
directly. Your own O2 already senses the problem from the other side (contested or ambiguous claims
must fail closed rather than guess) — but the cheapest way to never guess is to never parse. **My
answer to O2: the claim stays social; preflight and the driver may *surface* it, never honour it.**
And to keep even the surfacing inside the boundary, surface it only when a designation is present and
the claim contradicts it (a designation-path gate). Claim visibility on the *unset* path is real value
and a real defect (F-5), and it belongs in the same follow-up register as C-3, not in this delta.

**K-3 — O1 (claim above global default): I argue against, per R-2 above.** The ordering you propose
lets a participant file a file and defeat a standing owner instruction, which reinstates the
first-mover selection the idea exists to remove. Your concern survives and is answered by
surfacing + the one-line per-idea edit.

**K-4 — O3 field naming: I argue for `implementer:`, per R-4.** Your intent-versus-record separation
is a genuine benefit and I would take it if it were free. It is not: it creates a key an owner will
plausibly type and the tool will plausibly ignore. Buy the separation with the recorded resolution
source instead.

**K-5 — O5 deck-level default: you cannot ship "the two required levels only"; the deck level comes
for free and must be acknowledged.** `[defaults]` is read by `LoadDefaults` across `configLayers`
(central → deck `parley-deck/agents.toml` → `agents.local.toml` → env), "later non-empty values win"
(`runtime.go:384-401`, `:415-437`), and `parley init` writes the header *"Project-wide policy
defaults; a deck's parley-deck/agents.toml overrides them"* (`runtime.go:634`). So placing
`default_implementer` in `[defaults]` yields **per-idea > deck > machine** by construction. Reading
machine-only would require a deliberate special case and would make this the one `[defaults]` key
that ignores the deck — the surprising outcome, not the safe one. **Recommend: accept the inherited
layering, state it in the protocol text, and pin it with a test that sets the key in both files and
asserts the deck wins.** Both ship absent.

**K-6 — O4 two-participant ideas: agreed, and the floor you want is already code.**
`driver_impl.go:282` returns `no non-implementer reviewers available`, and LE-11's `< 2` reviewer
guard is live. So a designation leaving zero reviewers already hard-stops — at Phase 6. Your
"surface at kickoff rather than at close" instinct is right and is the same relocation I proposed in
D8; see the diversity decision below for the designation-only form that keeps it inside the boundary.

**K-7 — Your 7(e) (trajectory `--implementer` mismatch becomes a gate): I argue out of scope.**
`parley trajectory configure` is an attended command requiring `--yes` and a non-empty named-checks
contract (`trajectory.go:105-125`); it freezes an evidence-policy patch author, not a launch target.
Making the two agree is defensible hygiene, but it couples a designation feature to the trajectory
evidence subsystem, it can fire on decks that use trajectory and no designation, and nothing in the
owner's request asks for it. Record it as a follow-up; do not build it here.

---

## Resolved design decisions

Each item states the decision, the evidence, and which of the three round-1 positions it adopts.

**S-1 — Answer.** Yes: an explicit per-idea designation plus an optional global default, instantiating
`facilitator.go`'s shape. Unanimous across all three round-1 artifacts; nothing in cross-review moved
it.

**S-2 — Field: `implementer: <agent-id>` in `00-prompt.md`, optional.** Adopts zcode-1's and my
lean over kimi-1's `designated_implementer:`. Decisive argument (R-4): a distinct name makes the
obvious key silently inert. Intent-versus-record separation is bought back by S-8.

**S-3 — Present-empty semantics: fail closed, deliberately diverging from `facilitator:`.**
`FacilitatorRoleFromMeta` (`facilitator.go:43-53`) treats `facilitator:` with an empty value as
**undeclared** — `if id == "" { return FacilitatorRole{} }`. Mirroring that for `implementer:` would
mean an owner who types the key and leaves it blank silently gets list order. The divergence is
principled, not arbitrary: `facilitator:` is an *exclusionary* declaration, so an empty one excludes
nobody and is harmless; `implementer:` is a *positive appointment*, so an empty one appoints nobody
and the silent destination is `participants[0]`. It is also cheaply implementable —
`ReadFrontmatter` (`workspace.go:390-394`) does `key, value, ok := strings.Cut(line, ":")` then
`meta[key] = value`, so a line `implementer:` **does** create a map entry with `""`, and
`raw, present := meta[ImplementerKey]` distinguishes present-empty from absent in one line where
`meta[key]` alone cannot. Present-empty → `gate{Kind:"implementer-designation"}` naming the field.
Absent → untouched. This is a designation-path-only behaviour; a deck with no such line is byte-identical.

**S-4 — Precedence (four tiers, no claim tier, no materialization).**

1. `IMPLEMENTATION.md` `implementer:` — the re-entry pin. *(existing, unchanged)*
2. `00-prompt.md` `implementer:` — the designation. *(new; inert when absent)*
3. `[defaults].default_implementer`, layered deck > machine. *(new; UNSET on ship; inert when absent)*
4. Today's chain verbatim: `FINAL.md` `implementer:` / `drafted-by:`, then `participants[0]` in the
   driver / `""` in consensus. *(existing, unchanged)*

With 2 and 3 absent, resolution runs 1 → 4, which is today's chain exactly. Adopts zcode-1's ladder
with my C-2 correction (live read, not materialized) and against kimi-1's O1 (no claim tier).

**S-5 — Unavailable designee: fail closed with the confirm pre-built.** Invalid designation
(non-participant, declared non-participating facilitator, present-empty) → blocking gate. Valid but
unavailable at the §9.0 ping → blocking gate whose `Confirm` offers both exits (re-designate, or
record `implementer_waived: <id> — <reason> — confirmed <date>` and fall to the default chain).
Mid-Phase-5 loss → escalate via `inbox/` and halt; a partial tree is evidence and reassignment over
it is a human decision. Adopts my D4 and kimi-1's item 3 over zcode-1's automatic fall-through, on
the LE-7/LE-11 and L699 precedent and on the fact that the fallback's destination is list order
(R-1). zcode-1's speed objection is answered by the pre-built `Confirm`.

**S-6 — Re-entry pin conflicts: pin wins, but never silently.** When tier 1 and tier 2 both resolve
and **disagree**, the driver escalates rather than quietly honouring the pin. Without this, an owner
correcting a wrong designation after a failed Phase-5 attempt gets no signal that the correction was
ignored (my round-1 Q5). Requires tier 2 to be present, so it is designation-path only. Reconciles
my Q5, zcode-1's O2 and kimi-1's item 3 into one rule.

**S-7 — Machine/deck precedence: per-idea > deck > machine, inherited from `[defaults]`, not
invented.** Per K-5. Both files ship the key absent. Pin with a test that sets it in both and asserts
the deck wins.

**S-8 — Recorded resolution source.** `agent.implementer_resolved {idea, implementer, source}` on
the existing durable store, modelled on `agent.model_diversity` (`driver_impl.go:198-206`), plus the
existing one-line stdout notice beside `driver: implementing via %s` (`driver_impl.go:217`). `source
∈ {implementation-md, designation, global-default, final-md, positional}`. This is what makes
intent-versus-record greppable by field rather than by key name (answering K-4), and what would have
made F-1 and F-7 visible years ago. **Boundary note: this is the one element of the delta that adds
a record on the unset path.** It changes no selection. I propose shipping it always-on and recording
it as an explicit `## Agreed trade-offs` entry; if any participant or the owner reads the boundary
strictly, the conservative variant is to emit it only when tiers 2 or 3 fired, and I will sign that
variant without argument.

**S-9 — Resolver consistency: share the reading, keep both tails.** One resolver in
`internal/protocol` returning `(id string, src Source, ok bool)`; `driver_impl.go` keeps
`if !ok { id, src = eligible[0], SourcePositional }`; `consensus.go` keeps `if !ok { return
participants }`. Both tails stay byte-equivalent and stay pinned by their existing tests —
`TestResolveImplementerFromRoleMetadata` (`internal/app/app_test.go:1557`) asserts the positional
fallback, `TestUnresolvableImplementerExpectsEveryone` (`internal/consensus/roundgate_test.go:102`)
asserts the `""`→everyone path documented at `consensus.go:723-729` as failing *"closed toward asking
for more"*. `parley wait` inherits the fix through `ExpectedRoundParticipants`
(`consensus.go:1085-1089` → `internal/driver/phasedigest.go:364`, its only caller). **X-5 adds a
requirement all three of us missed**: the shared resolver must take the eligibility-filtered list as
an explicit parameter, and the consensus call site must pass the same facilitator-filtered list the
driver does, or unification will preserve the divergence it was built to remove. My round-1 MIN-7
against a zero-caller `consensus.ResolveImplementer` export (lean-organizer, dropped as F15) does not
re-litigate this: that finding was about an export with no callers and a doc describing behaviour it
did not have. This one has two callers plus a third consumer.

**S-10 — Self-contained FINAL: strengthen the text, add no gate.** `COOPERATION.md:433` already
requires `FINAL.md` + `IMPLEMENTATION.md` to be self-contained *"enough that a fresh agent or the
auto-drive driver can implement or resume from them alone, without session transcripts"*, and the
machine floor under it is already `ValidateFinal` (K-1). Designation makes drafter ≠ implementer the
common case, which makes that requirement load-bearing rather than aspirational. Protocol delta: one
sentence in Phase 5 saying the self-containment rule is what makes a designation executable. Plus my
D6 ordering-only tweak — `firstEligibleHeadlessAgent` (`driver_consensus.go:109-123`) prefers a
drafter that is not the designated implementer, **falling back to the implementer when no other
eligible drafter exists**, so a two-participant deck stays satisfiable. It uses the filter slot the
function already has, fires only when a designation exists, and I do **not** propose forbidding
drafter == implementer. Where they coincide, `FINAL.md` records the concentration in one line — §15.5's
existing device applied to a second role, not a new gate.

**S-11 — Reviewer diversity: relocate the existing gate for the designation case only.**
`checkModelDiversity` (`driver_impl.go:179-209`) already emits the event, warns by default, escalates
under `require_model_diversity`, and `driver_impl.go:282` already hard-stops on zero non-implementer
reviewers. Delta: when a designation is present, run the same check at Phase 0, so an unsatisfiable
roster costs an error message rather than a wasted design cycle (my D8, kimi-1's K-6/O4). No new gate
class, no change when no designation exists. For this roster the question is moot —
`parley roster show --scope machine` reports three distinct model vendors — but a two-participant
deck is where it bites. The *standing-concentration* risk (one model accumulates all implementation
context) is not a gate problem; all three of us named it and all three of us agree it belongs in the
post-release evidence report.

**S-12 — Role scope: code and non-code, and it never moves a gate.** Phase 5 already covers non-code
work (`COOPERATION.md:508`), and the owner's words are *"writes the code if needed, or does whatever
work is needed"*. The implementer is *the participant that executes FINAL*, whatever the artifact.
Exactly one at a time; no co-implementers, no split scopes. **Explicitly: designation names who
executes, never which gate applies.** `parley protocol publish` stays TTY-attended, core publication
stays owner-attended, npm/winget/Homebrew publication stays the organizer's attended release step,
and a designated implementer inherits no authority to perform any of them. A non-code designation is
not a route around an attended action. This needs to be a protocol sentence, not an assumption,
because "does whatever work is needed" is exactly the phrasing a future reader could stretch.

**S-13 — Design rounds and signoffs: unchanged; the implementer participates fully.** The owner's
"the others evaluate **it**" is about the work — Phase 6 — which already excludes the implementer
(L234, `consensus.go:730-746`). Removing the implementer from design would shrink quorum (§5) and
contradict Phase 7's explicit *"Each active participant (implementer included) APPENDS their signoff
block"* (L578). Designation changes **who executes Phase 5 and Phase 8, and nothing else** — worth one
protocol sentence because the owner asked the question explicitly. Unanimous across all three.

---

## The owner's boundary: default-path ledger

*"Ship the mechanism with the global default UNSET, so behaviour without it is exactly today's."*

**Preserved byte-identical on a deck that sets neither `implementer:` nor `default_implementer`:**
the resolution chain (S-4 tiers 1→4 collapse to today's two-file walk plus each caller's existing
tail); both terminal fallbacks and both pinned tests (S-9); `checkModelDiversity` timing and
severity (S-11); `firstEligibleHeadlessAgent` drafter ordering (S-10 tweak requires a designation);
`ValidateFinal` / `finalScaffoldReason` (K-1, no change); the `impl-claim` file's inertness (K-2 —
not parsed, not honoured, not gated); `roles:` advisory semantics; quorum, roster, signoff weight;
every attended gate (S-12).

**Designation-path only — fires only when tier 2 or tier 3 is set:** present-empty and invalid-target
gates (S-3); unavailable-designee gate (S-5); pin-conflict escalation (S-6); deck>machine resolution
of the new key (S-7); Phase-0 diversity relocation (S-11); drafter-separation preference (S-10).

**The one unset-path addition, declared rather than hidden:** the `agent.implementer_resolved` event
and its stdout line (S-8). It changes no selection. Explicit `## Agreed trade-offs` entry required;
conservative designation-only variant available and acceptable to me.

**Deliberately excluded and registered as follow-ups, each with its evidence already gathered:**
(a) the FINAL.md key-set widening and the FINAL-draft-prompt key emission (C-3 — would change
selection on 51 of 85 historical FINAL.md files and on every future unset deck); (b) automatic
claim parsing and unset-path claim gating (K-2 — would change selection on any deck holding a claim
file); (c) trajectory `--implementer` cross-checking (K-7); (d) `parley run --participants` order
documentation — I now judge even this a hardening of an accident into a contract (my round-1 R-5) and
would rather the designation be documented as the supported control and the order dependency simply
be *observable* via S-8; (e) `facilitatorConflictGates` iterating every idea in the deck
(`preflight.go:328-339`), which is why an unrelated historical idea blocked this very run's
preflight. **(e) is a constraint on this delta even though its repair is excluded**: the new gate
must be scoped to the idea being run, or this idea ships the same trap it tripped over.

---

## Inherited Windows limitations

Assessed against this idea's delivery criteria, per the release-boundary constraint. What evidence
permits me to resolve, I resolve; what it does not, I leave as the owner's open question rather than
inheriting an authorization.

**Resolved with evidence — this delta adds no Windows surface and cannot resolve the inherited
one.** The change touches `internal/protocol` (new file + `workspace.go`), `internal/app`
(`preflight.go`, `driver_impl.go`, `driver_consensus.go`), `internal/consensus`, and
`internal/config/runtime.go`. None of those specific files is build-tagged; the only `go:build` files
anywhere near them are `internal/app/{budget_attended_windows,budget_attended_other,termios_*}.go`
and `internal/driver/proclive_{unix,windows}.go`, none of which the delta touches. The new config key
rides `configLayers` (`runtime.go:384-401`), which resolves the central path through
`CentralAgentsPath()` and everything else through `filepath.Join` — already portable and already
exercised by every shipped `[defaults]` key. No new file locking, no ACL handling, no path parsing.
**Verdict: this idea neither worsens nor touches the Windows CI failure recorded at run
36009912946 (Ubuntu PASS, macOS PASS, Windows FAIL).**

**Not resolvable here, and must not be silently inherited.** The handoff records two open items:
CLI winget **not submitted**, held on the unresolved Windows distribution decision (line 30), and the
owner decision itself — *"choose explicit deferral with Windows labelled experimental/unvalidated, or
authorize a separate reviewed Windows-portability track"* (line 46). That is the **prior** release's
held item. This idea's release must not treat its own authorization ("release it and deploy it
through all channels directly") as a resolution of it: the owner authorized releasing *this
mechanism*, having anticipated exactly one owner question, and the Windows scope decision is a second,
pre-existing one that predates this idea. **Recommendation to the organizer:** carry the prior
release's Windows posture forward unchanged and unclaimed, state in the done-report that the CLI
winget hold is inherited rather than newly incurred, and make no catalog-availability claim. I
inherit no authorization for Windows platform architecture and propose none — consistent with the
non-goal in `00-prompt.md`.

**One delivery-criteria item that is a release-time check, not a design input.** The next CLI minor
must be chosen above whatever is actually published at staging time; with a 1.49.1 candidate staged
(`481fb657f0e219535edc83ca286bd33891cf83dc`) but explicitly *"NOT tagged, published or installed"*,
that number is not knowable from here. The protocol hunks stage on top of core 2.13.0, whose bytes I
verified in round 1 as identical to the skill worktree's bundled copy (F-9 / X-6), so the base is
measured rather than assumed.

---

## Existing alternatives

Per §15.6(a), for the elements whose shape **changed** in this round. My round-1 table covers the
unchanged elements and is not restated.

| Element changed in round 2 | What already ships | Locator | Why the round-2 shape is the right one |
|---|---|---|---|
| Global default resolution (C-2: live layered read, was Phase-0 materialization) | `LoadDefaults` merging `[defaults]` across `configLayers`, "later non-empty values win"; `parley init` emits the block with a deck-overrides header | `internal/config/runtime.go:384-401`, `:415-437`, `:634-639` | Materialization would make `default_implementer` the only `[defaults]` key with bespoke mechanics and would make the driver a writer of an author-owned artifact. Nothing writes `00-prompt.md` today. |
| Deck-vs-machine layer (S-7) | Already a shipped layer, not a proposal: deck `parley-deck/agents.toml` overrides central `~/.parley/agents.toml` | `runtime.go:384-401`; `COOPERATION.md:59` | kimi-1's O5 "ship two levels only" is not available — the third layer arrives with the block. Accepting it costs a test; refusing it costs a special case. |
| Pre-dispatch FINAL validation (K-1: cite, don't build) | `protocol.ValidateFinal` — documented as "the ONE gate for a FINAL.md, used by manual finalization and by the driver"; `RequiredFinalSections` already includes `## Observable acceptance criteria`; wired pre-Phase-5 | `internal/protocol/finalsections.go:18-26`, `:92`; `internal/driver/consensus.go:57`, `:61`, `:175-190` | The gate kimi-1 proposes exists and fires. The only remaining delta is semantic judgement, which no validator performs and which would change the gate for every unset deck. |
| Zero-reviewer floor (K-6/S-11) | `no non-implementer reviewers available`; LE-11 `< 2` guard; `checkModelDiversity` warn/escalate | `internal/app/driver_impl.go:282`, `:179-209`; `COOPERATION.md:275`, `:538` | The floor exists; only its *timing* moves, and only when a designation is present. |
| Unavailable-role handling (S-5) | LE-7/LE-11 halt-on-unavailable for the goal-done checker; §9.0's `<id> — reason — confirmed <date>` record form; `gate{Kind,Detail,Confirm}` with a pre-built confirm string | `COOPERATION.md:275`, `:699`, `:872-874`; `internal/app/preflight.go:337-341` | Every piece of the fail-closed-with-an-exit design is shipped machinery; the delta is which field it reads. |
| Claim handling (K-2: surface, never parse) | `impl-claim` prose with no machine reader | `COOPERATION.md:443`; `grep -rn "impl-claim" --include="*.go" .` → empty | Parsing changes selection on unset decks. Surfacing on the designation path only keeps the boundary intact and still closes the contradiction case. |

Sources consulted this round: the full phase-2 packet body; `internal/protocol/`, `internal/app/`,
`internal/config/`, `internal/consensus/`, `internal/driver/`, `internal/runner/`, `internal/track/`,
`internal/tui/`; `parley roster show --help`; `git diff 59b0458 HEAD`; all three round-01 artifacts;
`source-context/prior-release-handoff.md`; `organizer-notes.md`.

---

## Concerns / open questions

1. **S-8 is the only unset-path addition and needs an explicit group decision, not silence.** I
   propose always-on with a recorded trade-off; the conservative variant is designation-only. Either
   is fine; deciding it by not mentioning it is not.
2. **X-5's eligibility axis changes the unification design and arrived after all three round-1
   artifacts.** Whoever drafts FINAL must carry it, or S-9 unifies the reading and preserves the
   divergence. I own no verdict on zcode-1's find beyond confirming it; zcode-1 should confirm my
   reading of the fix shape.
3. **`--no-implement` and design-only ideas** (my round-1 Q4, still open). A designation on an idea
   that never reaches Phase 5 names an agent that never acts. Harmless, but the availability gate
   should probably not fire for such a run, and I still have not traced whether `--no-implement` is
   visible at preflight time. Whoever implements should trace it; if it is not visible, the gate
   should fire and the waiver line is the exit.
4. **"Implementation is the heaviest token work of a run" remains UNVERIFIED.** Unchanged from my
   round-1 item 14; zcode-1 and kimi-1 both confirmed the study's 30-40% figure *as quoted* with its
   caveats, which is a different claim. The study's own numbers put re-orientation at ~33-39% and
   compaction at 42% of uncached input. If this sentence appears in FINAL it must appear as
   testimony, never as established (§15.2).
5. **My item-13 evidence asymmetry persists.** I read only the in-deck copy of the token study;
   zcode-1 read the external original. Anyone citing the original's numbers should cite zcode-1's
   reading, not mine.
6. **`source-context/preflight.json` is stale** (my F-8): its `protocolSha256` is `12e4b31c…`, the
   v1.48.0-era text, not this tree's `8ce83cde…`. Harmless for this `source`-role deck because
   `preflight.go` short-circuits to `source-advisory` before any drift comparison, but it remains a
   hard constraint on the new gate: **read live idea frontmatter, never `meta/version.json`.**

---

## Risks

- **R-A — the follow-up register becomes a graveyard.** C-3, K-2, K-7 and boundary items (d)/(e) are
  all real defects this delta deliberately does not fix. Naming them is not fixing them. The honest
  mitigation is that S-8 makes the largest one (list-order selection) *visible* for the first time,
  so the next reader does not have to re-derive it from three files as we did.
- **R-B — S-3's divergence from `facilitator:` is a real inconsistency.** Two sibling role fields
  will treat present-empty differently. I argued the asymmetry is principled (exclusion vs
  appointment), but a future reader will find it surprising, so the reason must be in the code
  comment, not only in this artifact.
- **R-C — three copies, one hunk.** Unchanged from round 1 (R-4) and reinforced by zcode-1's O5:
  the drift guard covers the three COOPERATION.md copies, not SKILL.md prose restatements. A
  byte-comparison of the changed region across all three, plus a grep of the skill tree for Phase-5
  restatements, belongs in the acceptance criteria.
- **R-D — a new gate class is a new way for `parley run` to stop**, and this deck already proved the
  failure mode: an unrelated historical idea's facilitator conflict blocked this run's preflight
  because `facilitatorConflictGates` iterates every idea (`preflight.go:328-339`). Scope the new
  gate to the idea being run. Non-negotiable for me.
- **R-E — designation concentrates implementation in one model.** Unchanged from round 1 (R-2); all
  three of us named it independently. The mitigation is the sequencing the owner already chose:
  ship UNSET, recommend on evidence afterwards.
- **R-F — cross-review agreement is not convergence.** Three artifacts agreed on two wrong line
  numbers between them and one wrong locator (X-1, X-2, X-3). The FINAL drafter should re-measure
  every locator it copies rather than transcribing from any round artifact, mine included.

---

## Consensus readiness

**Ready for consensus, with one substantive block and one required recorded decision.**

**Block B-1 — K-2 must be withdrawn or overruled before FINAL.** Automatic claim parsing changes
implementer selection on decks that set nothing, which contradicts the owner's constraint directly.
I will not sign a FINAL that carries it. This is a narrow, mechanical objection and kimi-1's own O2
already leans the same way; I expect it resolved in one exchange rather than a third round.

**Required recorded decision D-1 — S-8's always-on event.** Not a block: I will sign either variant.
But it is the one place the delta touches the unset path and it must land in `## Agreed trade-offs`
with its reasoning, not be inherited by silence.

Everything else I consider settled: S-1 through S-13 have either unanimous round-1 support or an
evidence-backed resolution above. I do **not** consider C-3, K-2's unset-path half, K-7 or boundary
items (d)/(e) open questions — they are decided *out*, with their evidence preserved for a successor
idea. If either peer wants any of them in, that is a new argument against the owner's sentence and
needs to be made as one.

Cross-review 1 of 3 is, in my reading, sufficient: the remaining delta is one withdrawal and one
recorded trade-off. I see no need for cross-review 3 and would support the organizer moving to
consensus after round 2 responses land, provided B-1 is resolved.

---

## Process notes

**Drafter and implementer for this run (today's protocol, proper acts to follow).** All three of us
recorded willingness in round 1, so per the organizer's instruction we should converge on one drafter
and one implementer before Phase 4/5 rather than race. My proposal, with reasons that do not favour
me:

- **Implementer: kimi-1.** Two evidence-backed reasons. (a) All three of us independently said the
  post-release recommendation needs more than the deck's single implementation data point; kimi-1 has
  drafted here (lean-organizer FINAL) but never implemented, so designating kimi-1 *produces the
  missing data point* rather than deepening an n=1. (b) This delta's dominant risk is locator
  fidelity across three protocol copies and six `implementer` readers (X-6, X-7, R-C, R-F) — and
  kimi-1 was the only one of the three of us whose round-1 code locators were all correct (X-1, X-2).
  That is a small sample and I say so, but it is the only directly relevant evidence on the table.
- **Drafter: not the implementer** (S-10). zcode-1 or me; I volunteer but have no preference and
  will step aside for zcode-1 without discussion. kimi-1 drafted last run, so rotation argues
  against kimi-1 drafting as well as implementing.
- **Reviewers: the two non-implementers.** Note the argument *against* my implementing: the deck's
  only implementation-review evidence is my lean-organizer review (1 CRITICAL, 8 MAJOR, all fixed
  across 3 fix-up cycles). Spending that on implementation removes the deck's one demonstrated
  high-yield reviewer from the reviewer set. I would rather review this than implement it.

I have **not** perfected a claim here, and this is not one. Under the protocol in force a drafter
volunteer is declared in the consensus signoff block (`Drafter: yes`, L409) and an implementer claim
in `inbox/<from>-to-all_<slug>_impl-claim.md` before Phase-5 work begins (L443); neither is a round-2
artifact and I will not pre-empt those phases. If kimi-1 or zcode-1 prefers a different assignment I
will sign it.

**Owner-decision boundary.** Nothing above names the owner's future global default. The role proposal
is for *this run under today's rules*; the global default ships UNSET and the recommendation is a
post-release act on evidence this deck does not yet have. Conflating the two would be the quiet
self-selection this idea exists to remove.

**Evidence ownership (§15).** Verdicts in X-1…X-5 and X-7 are mine as a **non-owner** of kimi-1's and
zcode-1's first-canonical claims. X-6 restates a claim I own (F-9), so I issue no verdict on it and
supply evidence for a non-owner to verdict. C-1 through C-3 are corrections to claims I own, which an
owner may correct. F-1…F-9 remain mine and unverdicted by me; kimi-1's V12 supplies the operational
corroboration my R-6 said was missing for F-1, and a non-owner still owns that verdict. Where I cite
the copied token study or the copied `preflight.json` I am transcribing testimony and have marked it.

**Independence and scope.** I read `00-prompt.md`, `source-context/owner-brief.md`, the phase-2
packet body, all three round-01 artifacts and `organizer-notes.md` before writing, as instructed for
a cross-review round. I wrote only this file. I implemented nothing, released nothing, committed
nothing, edited no other participant's artifact, spawned no roster identity, and altered no
accounting, worktree or budget state.

<!-- Original language of the owner's quoted words: Slovak; translations as supplied verbatim in 00-prompt.md and source-context/owner-brief.md. -->
