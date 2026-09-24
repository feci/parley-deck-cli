---
idea: meta-protocol-change-designated-implementer
status: final
author: claude-1
consensus-date: 2026-09-25
participants: [claude-1, kimi-1, zcode-1]
---

# FINAL — designated implementer

Drafted by claude-1, the agreed FINAL drafter (drafter-volunteer note in `inbox/`, `Drafter: yes`
in its consensus signoff, `COOPERATION.md:409`). The declared facilitator `codex-1` is a pure
organizer: not a participant, not a signer, not an implementer, not a code verifier. This artifact
is **static and frozen on publication**; the living companion is `IMPLEMENTATION.md` (Phase 5).

Every one of the five items the consensus left open for explicit peer evaluation (OPEN-1 … OPEN-5)
was accepted by all three participants in their signoff blocks. **This FINAL states the accepted
disposition as the single operative rule in each case. No mutually exclusive alternative remains
operative anywhere below.**

**Locator discipline.** Every code and protocol locator in this document was re-measured at HEAD
`07ceb0a78b9f4ba89c171cdd9946fe24bb05930d`, not transcribed. `git diff --stat 73ee923 HEAD --
internal/ cmd/ parley-deck/COOPERATION.md` is empty: no code or protocol text moved during the run.
Three locators the consensus carried are corrected below (see `## Changes since the signed
consensus`). Per the VC-4 resolution the implementer MUST re-measure again rather than transcribe
from here — three rounds produced locator corrections from all three participants.

## Final plan / specification

The protocol gains an **owner-set designation of the participant who executes FINAL**, expressed as
a per-idea field and an optional machine-level default, resolved through one chain, validated
fail-closed, and **completely dormant until somebody sets it**. Reviewers stay the non-implementers;
no gate, quorum, roster or attended boundary moves.

The delta is stated as **57 numbered rules, R1 – R57**. In R1 – R49 the mechanism rules carry a
tag: **[D]** fires only when a designation is present (tier 2 or tier 3); **[A]** applies always. A
deck that sets neither field gets byte-identical behaviour to today. R50 – R57 are the release and
staging constraints and carry no tag — they bind the organizer's release step, not the code.

### A. The per-idea field and its four states

**R1 [A].** The per-idea field is `implementer:` in `00-prompt.md` frontmatter — not
`designated_implementer:` (ALT-2 rejected; kimi-1 withdrew the distinct name in round 02). One
sentence of protocol text states the asymmetry: the `00-prompt.md` field is a **designation** (an
instruction about who should execute), while `IMPLEMENTATION.md`'s `implementer:` is an **outcome
record** (who did).

**R2 [A].** Parsing distinguishes **four** states from one expression. `ReadFrontmatter`
(`internal/protocol/workspace.go:368-396`) splits on the first colon (`:390`) and stores
`meta[TrimSpace(key)] = TrimSpace(value)` (`:394`), so key presence survives an empty value and
`raw, present := meta[implementerKey]` separates all four:

| Frontmatter | Meaning | Machine behaviour |
|---|---|---|
| key absent | no per-idea designation | tier 3 may fire; otherwise today's chain |
| present, value empty or whitespace | **incomplete designation** | blocking gate (R16) |
| present, value `none` (trimmed, any casing) | **explicit per-idea non-designation** | tier 3 suppressed; today's chain; **no gate** |
| present, value an agent id | designation | validated fail-closed (R16, R18) |

**R3 [A].** The `none` state MUST be test-pinned (T-2). Without that pin the opt-out silently
collapses into "absent" and tier 3 fires on an idea that deliberately declined a designation. This
is zcode-1's round-02 three-state insight, which all three round-01 proposals missed; collapsing
present-empty and `none` into one meaning (ALT-10) is **rejected** — it makes a typo
indistinguishable from an intention, which is the failure mode this idea exists to remove.

**R4 [A].** Malformed values **fail closed and are never repaired**. `ReadFrontmatter` performs no
comment stripping, so `implementer: kimi-1  # from the global default` yields the literal value
`kimi-1  # from the global default`. It MUST fail the eligibility check and gate; it MUST NOT be
trimmed to `kimi-1`. The same hazard is already recorded in code: `CreateIdeaFull` puts provenance
in an HTML comment **below** the frontmatter fence precisely because inside it "`ReadFrontmatter`'s
`key: value` split would ingest it as a junk key" (`internal/protocol/workspace.go:197-198`).

**R5 [A].** Quoted forms `"kimi-1"` and `'kimi-1'` are accepted, because both existing resolvers
already apply `strings.Trim(strings.TrimSpace(meta[k]), "\"'")` — `internal/app/driver_impl.go:125`
and `internal/consensus/consensus.go:772`. The shared reader keeps that trim unchanged.

**R6 [A].** The divergence from the sibling field goes in a **code comment** at
`FacilitatorRoleFromMeta` (`internal/protocol/facilitator.go:44-48`), not only in prose: that
function deliberately collapses present-empty into absent because `facilitator:` is **exclusionary**
(an empty one excludes nobody), whereas `implementer:` is **appointive** (an empty one appoints
nobody, and the silent destination is `eligible[0]`). `none` versus empty is what separates "nobody
deliberately" from "nobody by mistake".

### B. The machine-level default and the config chain

**R7 [A].** The global default is `default_implementer` inside the existing `[defaults]` block. It
is added to the struct pair — `globalDefaults` (`internal/config/runtime.go:283-291`, with the
`toml:"default_implementer"` tag) and `CentralDefaults` (`:315-329`).

**R8 [A].** It is merged in `mergeDefaults` (`:533-576`) following the existing **non-empty string**
pattern used by `Speed`, `PingTier`, `PreferredTransport` and `RosterChangePolicy` (`:534-545`):
`if s := strings.TrimSpace(x); s != "" { out.X = s }`. It MUST NOT be a presence-aware pointer
field. The `[defaults.loop]` pointer fields (`:563-575`, where "a deliberate `= 0` overrides a lower
layer's seed to unlimited") are the precedent for *clearing* and are deliberately **not** used here
(ALT-19 rejected).

**R9 [A].** Consequence of R8, stated because it is not intuitive: **an empty value at a higher
layer does not clear a lower layer's value.** The documented deck-wide suppressor is therefore
`default_implementer = "none"` — non-empty, so it wins the merge, and resolved as "no global
designation".

**R10 [A].** Resolution is a **live layered read**, never materialized into `00-prompt.md`
(ALT-1 rejected; all three participants withdrew Phase-0 materialization in round 03). **No new
`00-prompt.md` writer is added anywhere.**

**R11 [A].** The chain is **four layers deep, inherited not invented**. `configLayers`
(`internal/config/runtime.go:384-401`), low to high:

1. `~/.parley/agents.toml` — central, machine, optional (`:387-389`)
2. `parley-deck/agents.toml` — deck, optional (`:391`)
3. `parley-deck/agents.local.toml` — deck-local, optional (`:392`)
4. `$PARLEY_HEADLESS_AGENT_CONFIG` — highest, **non-optional when set** (`:394-399`)

`LoadDefaults` (`:419-438`) merges `[defaults]` across every layer: "Missing files are skipped;
later non-empty values win" (`:416-417`). Any key in that block inherits all four by construction.
FINAL states the **full chain**; "per-idea > deck > machine" as written in the round files is
shorthand for it. The environment variable's real name is `EnvAgentConfig =
"PARLEY_HEADLESS_AGENT_CONFIG"` (`:17`) — **not** `PARLEY_AGENT_CONFIG`. Reading machine-only
(ALT-3) is rejected: it would make this the one `[defaults]` key that ignores the deck.

**R12 [A].** **It ships UNSET.** `centralDefaultTemplate` (`:627` onward) emits
`default_implementer` **commented out**, a deliberate deviation from that template's shipped
active-key-with-inline-comment shape (`speed = "fast" # …`, `ping_tier = "hosted-pong" # …`,
`:636-639`). One sentence of protocol text records the deviation. There is no deck `agents.toml`
generator — `EnsureCentralDefault` (`:605-622`) writes only the central file — so the deck file
ships the key absent by construction.

### C. One resolution chain, four ranks, no claim rank

**R13 [A].** Dispatch selection resolves in exactly this order, and stops at the first hit:

1. `IMPLEMENTATION.md` `implementer:` — the **re-entry pin**
2. per-idea `00-prompt.md` `implementer:` — the **designation** (R1–R6) **[D]**
3. `[defaults].default_implementer` — the **live layered read** (R7–R12) **[D]**
4. **today's chain, verbatim:** `FINAL.md` `{implementer, drafted-by}` → driver tail `eligible[0]`
   (`internal/app/driver_impl.go:131-133`) / consensus tail `""` (`internal/consensus/consensus.go:778`)
   → expect everyone

Ranks 1 and 4 are exactly the two sources both resolvers read today, in today's order. Ranks 2 and 3
are the new, dormant middle.

**R14 [A].** **There is no claim rank, anywhere.** `inbox/{from}-to-all_{idea-slug}_impl-claim.md`
stays a social act with **no machine reader**: `grep -rn "impl-claim" --include="*.go" .` returns
nothing at HEAD (re-run for this FINAL: 0 hits). ALT-4 — a machine-honoured claim, whether as a
precedence rank or a parser — is **rejected**; it would reintroduce the first-mover race and change
unset-path selection. **No automatic claim reader ships in this idea.** A machine reader is deferred
(register item F1-iv).

**R15 [D].** Under a **live designation**, one subordination sentence is added to the Phase-5 text:
a claim **does not override the designation**. It is recorded, surfaced as advisory information on
the designation path only, and answered by a one-line per-idea edit or a §4 escalation. On the unset
path nothing changes at all: a claim with no designation is today's normal world and raises nothing.

### D. Validity gates, availability, and where the authority lives

**R16 [A].** **Validity gates are hard and fire on any run**, including a design-only run: an
`implementer:` value that is not in this idea's eligible set, or the present-empty state of R2. A
defective frontmatter line is worth stopping on even where Phase 5 is unreachable.

**R17 [D].** **Availability checks attach only to runs that can reach Phase 5.** `--no-implement`
lives in the same flag scope as preflight — `internal/app/app.go:1807` (`parley run`) and `:1159`
(`parley continue`), with preflight at `:1922` and the continue-path consumption near `:1252` — so
the scoping is decidable in-process. One implementation-time trace MUST confirm the flag is visible
where the gate runs; if it is not, the recorded waiver line of R18 is the exit and **no stall is
created** (register item F2).

**R18 [D].** **Tier-2 designee valid but unavailable at the §9.0 ping → blocking gate with the
`Confirm` pre-built**, offering exactly three recorded exits:

1. re-designate — edit `implementer:` to another eligible participant;
2. record `implementer_waived: {agent-id} — {reason} — confirmed {date}`;
3. write `implementer: none`.

The record shape is §9.0's own: "**Excluding** an unavailable agent from this idea's quorum requires
**explicit user confirmation** and is recorded in `00-prompt.md`" (`COOPERATION.md:872-874`), the
shape `CreateIdeaFull` already writes for `excluded:` (`internal/protocol/workspace.go:186-192`).

**R19 [D]. OPEN-1, as signed — the operative rule.** A **tier-3** designee that is a participant but
**fails the §9.0 ping** gets a **one-line notice and falls through to today's chain. There is no
gate at tier 3.** R18's tier-2 gate is unchanged. ALT-14 (a tier-3 blocking gate) is **rejected**.
The decisive reason: §9.0's confirmation requirement protects **quorum exclusion** — the excluded
agent stops deliberating and signing — while a tier-3 fall-through excludes nobody from quorum; it
declines to **appoint**, and its destination is exactly what an unset deck gets today, so **no idea
is worse off than today**, whereas a tier-3 gate makes **every** idea worse off for the duration of
one agent's outage. kimi-1 moved to this disposition in its signoff ("I am persuaded by the
argument, not by the 2-of-3 count"); zcode-1 and claude-1 hold it. Three concessions make the
fall-through non-silent and are **required**, not optional:

1. the notice names the designee, the reason (§9.0 ping failure), and the three exits that produce
   gate-like behaviour in one line — per-idea `implementer:` naming another agent, per-idea
   `implementer: none`, or a recorded `implementer_waived:` line;
2. the resolved-source record distinguishes a tier-3 fall-through caused by unavailability from an
   unset deck and from an opt-out (R32);
3. **invalid** tier-3 values keep their hard failure — only *unavailability* falls through.

**R20 [D]. OPEN-5, as signed — the operative rule.** At **tier 3 only**, an id that is the idea's
declared **non-participating facilitator** is **inapplicable**, exactly like a non-participant:
one-line notice and fall through. At **tier 2** it stays a **hard gate**. ALT-22 (a uniform hard
gate at both tiers, which is how AD-7 was first drafted) is **rejected**.

This is not academic. The owner's chosen machine-global default implementer is `codex-1`, and on
this deck **seven** ideas declare `facilitator: codex-1` — re-measured for this FINAL,
`grep -rln "^facilitator: codex-1" parley-deck/ideas/*/00-prompt.md` returns exactly:
`facilitator-packet-per-phase-bounds`, `meta-protocol-change-consensus-duty-gates`,
`meta-protocol-change-designated-implementer`,
`meta-protocol-change-facilitator-integrity-phase-coverage`,
`meta-protocol-change-lean-organizer`, `release-binary-reproducibility`,
`wait-boundary-vs-published-fixup`. Under a uniform hard gate, the moment the owner sets the default
they chose, **any relaunch of any of those seven would hard-gate on the owner's own standing
preference** — clearable only by deleting the preference or editing seven idea files. kimi-1 and
zcode-1 each re-ran this grep independently before signing.

**R21 [D].** Mechanically, R19 and R20 collapse into one condition: **is the resolved tier-3 id in
`eligible`?** `IneligibleForRoles` (`internal/protocol/facilitator.go:76-78`) is exactly
`Declared && !Participates && id == Facilitator`, and it is already the predicate the driver applies
when building `eligible` (`internal/app/driver_impl.go:54-60`). With
`facilitator_participates: true` set, `IneligibleForRoles` returns false and the designation is
simply valid — so the interaction arises only for pure-organizer declarations.

**R22 [A].** A **tier-3 designee absent from this idea's `participants:`** is **inapplicable, not
invalid** → one-line notice, fall through. Unanimous since round 02. ALT-20 (gate it) is rejected: a
standing preference legitimately predates an idea's roster, and gating would make the global
unusable on any mixed roster.

**R23 [D].** **Mid-Phase-5 loss of the implementer → escalate and halt.** The partial tree is
evidence; reassignment over it is a human decision.

**R24 [D]. OPEN-4, as signed — the operative rule.** **One sentence of protocol text, no code.** The
designation is validated against `participants:` **only**. **Nothing reads `excluded:` today** —
verified for this FINAL: the only non-test hits for the key in Go are the writer
(`internal/protocol/workspace.go:191`) and an unrelated JSON field on a preflight struct
(`internal/app/preflight.go:124`). A §9.0 exclusion of the designee must therefore **also remove the
id from `participants:`**, or the designation stays live. A second wrinkle, recorded: because
`ReadFrontmatter` keys a flat map, multiple `excluded:` lines collapse to the last one. A code-level
cross-check is **deferred** (ALT-21; register item F3). The failure mode is narrow: it needs an
owner to exclude the very agent they designated without editing `participants:`, and R18's tier-2
gate catches the common route to that state, because the exclusion follows a failed ping.

**R25 [A]. Gate siting is non-negotiable, and the authority is not preflight.** Every new gate reads
the **live frontmatter of the idea being launched** — never `parley-deck/meta/version.json` — and is
**scoped to that idea**. The authoritative fail-closed check belongs on the `newDriverImplOps`
`roleErr` path (`internal/app/driver_impl.go:45-99`; `roleErr` set at `:61-63`): per-idea by
construction, already the site whose comment reads "the driver NEVER silently falls back" (`:51`),
already computing `eligible` (`:54-60`) and the implementer (`:64`), and already escalated by four
role actions — `Implement` (`:214`), `OpenReviewRound` (`:278`), the goal-done checker (`:407-408`)
and `Complete` (`:504`). **ALT-7 (preflight as the authority) is rejected** for two measured
reasons:

- In `parley run`, `runTaskPreflight` executes at `internal/app/app.go:1922` while
  `runcontrol.Create` — which reaches `CreateIdeaFull` — is at `:1939`. **Preflight runs before the
  idea exists.**
- For existing ideas it would have to scan the deck, which is the `facilitatorConflictGates` shape
  (`internal/app/preflight.go:323-341`, `for _, idea := range status.Ideas`) that **halted this very
  run's preflight on an unrelated historical idea**. Shipping that shape again would ship the trap
  this run tripped over.

**R26 [D].** Preflight MAY carry an early **idea-scoped** copy of the **validity** gate for
ergonomics only, using the existing `gate{Kind, Detail, Confirm}` shape
(`internal/app/preflight.go:334-338`). It **MUST NOT** iterate all ideas. Whether it ships at all is
an implementer decision (register item F4); the authority is R25 either way.

### E. Re-entry identity

**R27 [A].** Rank 1 wins whenever tiers 1 and 2 agree, or tier 2 is absent — today's behaviour.

**R28 [D].** When tier 1 and tier 2 are both present and **disagree**, the driver **escalates**; it
never silently honours either. Correction is an **owner/author-recorded act**, never the incoming
implementer's own edit (ALT-15 rejected: self-appointment with a paper trail is not authorization;
zcode-1 withdrew the opposite position in round 03 as "self-dealing"). The `Confirm` offers:

- (a) restore the designation to match the pin; or
- (b) record `implementer_reassigned: {old-agent-id} to {new-agent-id} — {reason} — confirmed
  {date}`, which retires the abandoned attempt and re-pins.

A participant **may draft** the line; it is **not valid until confirmed**, and only then does the
incoming implementer log the pin update as the executor of that recorded decision. This is
designation-path only — it requires tier 2 present — so an undesignated idea behaves exactly as
today.

**R29 [D].** **Re-entry identity inside a run.** Resolution is recomputed at each
`newDriverImplOps` construction (`internal/app/driver_impl.go:45`) and **frozen for that dispatch**.
Before a pin exists, the recorded dispatch event is the durable record of what was dispatched, and a
live tier-2/tier-3 change against a recorded dispatch **escalates** rather than reassigning
silently. That comparison **fires only when the recorded source was `designation` or
`global-default`**, so on the unset path nothing compares and a changed `--participants` order
re-resolves exactly as today. `store.Store.Load()` (`internal/store/events.go:68`) supplies the
replay the comparison needs.

**R30 [A].** Resumability under a live read: nothing in the chain depends on a config file's value
at a past moment. Once Phase 5 writes `IMPLEMENTATION.md`, the pin governs every resume; before
that, a changed global behaves exactly like the pre-work metadata edit all three participants
already permit. A resumed run cannot be surprised by an edit it cannot see.

### F. Resolver consistency — and the boundary that keeps review honest

**R31 [A].** The two duplicate unexported `resolveImplementer` copies —
`internal/app/driver_impl.go:104-135` and `internal/consensus/consensus.go:751-779`, logic-identical
today apart from their tails — are replaced by **one implementation in `internal/protocol`**. Both
callers already import that package, so there is no import cycle and no third package is needed. The
shared function takes (a) an **ordered source list** and (b) the **eligibility list**, and returns
`(id, source, ok)`.

**R32 [A]. OPEN-3, as signed — the operative rule. Dispatch selection and review exclusion are two
different reads, and only one of them sees the new tiers.**

- **Dispatch selection** — the driver, `newDriverImplOps` — reads **all four ranks** of R13.
- **Actual-implementer review exclusion** — the consensus side, `expectedRoundParticipants`
  (`internal/consensus/consensus.go:730-748`) — keeps **exactly today's two artifact sources**:
  `IMPLEMENTATION.md{implementer}` → `FINAL.md{implementer, drafted-by}`. **Tiers 2 and 3 never
  enter it.** ALT-12 is **rejected**.

The reason is categorical, not quantitative: **review exclusion must name who actually implemented —
a recorded fact — never who was instructed to. A preference cannot retroactively define who did the
work.** The exposure this closes is concrete: `expectedRoundParticipants` feeds three consumers,
including `reviewConsensusVoters` (`:136-144`, calling at `:143`), the round-completeness check at
`:176`, and the exported `ExpectedRoundParticipants` (`:1089-1091`) that **`parley wait` reuses**.
Had tiers 2–3 entered that read, a live designation could shrink a re-opened review round's expected
artifact set on an idea that never opted into anything. With tiers 2–3 confined to dispatch, that
exposure is **zero** and `parley wait`'s expected-round computation is untouched by this feature in
every state. Both peers confirmed the refutation case does not exist: review rounds open at Phase 6,
after `IMPLEMENTATION.md` is published, so the pin always exists where expectations are computed.

**R33 [A].** **Both terminal tails stay byte-identical** and stay pinned by their existing tests,
which MUST pass **unmodified**: `internal/app/app_test.go:1557`
(`TestResolveImplementerFromRoleMetadata` — pins the driver's `participants[0]` tail and the rank
order), `internal/consensus/roundgate_test.go:102` (`TestUnresolvableImplementerExpectsEveryone`)
and `:111` (`TestFinalDrafterIsTheFallbackImplementer`).

**R34 [A].** **The X-5 eligibility divergence is preserved, parameterised and pinned — not
repaired.** The driver passes the facilitator-filtered `eligible`
(`internal/app/driver_impl.go:64`); the consensus side passes the raw list it was handed
(`internal/consensus/consensus.go:734`). Making the eligibility list an **explicit parameter at both
call sites** satisfies the requirement that unification must not freeze the divergence inside a
shared function: it becomes visible, parameterised and testable instead of accidental. **Changing
the consensus call site to the filtered list stays OUT** (ALT-13 rejected, deferred): the filter is
`Declared && !Participates && id == Facilitator` (`internal/protocol/facilitator.go:76-78`), and
exactly that state is a blocking preflight gate (`FacilitatorRole.Conflict`, `:59-71`) — so the
difference is reachable only when preflight is skipped or from read-only surfaces that never ran it.
**Narrow is not unreachable**, which is why it is pinned rather than assumed away.

**R35 [A].** `expectedRoundParticipants` keeps its **observable output in every state**. So do
`reviewConsensusVoters`' deliberation-track short-circuit (`:140-142`) and
`ExpectedRoundParticipants`.

### G. Designation-only behaviour preferences — inert without a designation

**R36 [D]. Drafter separation.** `firstEligibleHeadlessAgent`
(`internal/app/driver_consensus.go:109-123`) prefers an eligible drafter that is **not** the
designated implementer, falling back to the designee when no other eligible drafter exists, using
the filter slot the function already has (it already skips `role.IneligibleForRoles(p)` at
`:113-115`). The degenerate fallback is test-pinned. `drafter == implementer` is **never
forbidden** — two-participant decks need it — and where they coincide, `IMPLEMENTATION.md` records
the concentration in one §15.5-style line.

**R37 [D]. Kickoff model-diversity check.** `checkModelDiversity`
(`internal/app/driver_impl.go:179-211`) runs **unchanged** where it runs today, from
`OpenReviewRound` (`:286`). Under a designation it **also** runs at kickoff, so an unsatisfiable
roster costs an error message rather than a design cycle. **No new gate class and no severity
change**: the existing `required` computation (frontmatter `require_model_diversity`, forced true on
the `fast` track, `:184-190`) and its always-on `agent.model_diversity` event (`:195-205`) are
untouched.

**R38 [D].** A **two-participant designated idea gets a kickoff warning, not a block**.
`COOPERATION.md:746` keeps the two-participant rules unchanged.

**R39 [A].** **Zero invokable non-implementer reviewers stays the hard stop it already is** —
`internal/app/driver_impl.go:281-283` and `COOPERATION.md:538` — surfaced earlier under a
designation but never relaxed.

**R40 [A]. Text only, zero behaviour.** One sentence at `COOPERATION.md:433` adds designated runs to
the FINAL self-containment trigger list (today: "complex, `auto_implement`, driver-managed, or
pipeline ideas"), because a designation makes drafter ≠ implementer the common case.

**R41 [A]. No new pre-dispatch FINAL validation.** ALT-8 is **rejected as already shipped**:
`protocol.ValidateFinal` (`internal/protocol/finalsections.go:92-106`), built on
`RequiredFinalSections` including `## Observable acceptance criteria` (`:18-26`), is already
enforced before leaving Phase 4 — the driver calls `finalScaffoldReason`
(`internal/driver/consensus.go:57` and `:61`, implemented at `:175-189`, delegating to
`ValidateFinal` at `:185`), and the manual path calls it at `internal/consensus/consensus.go:343`.
Both peers withdrew their new-validation proposals on this evidence.

### H. What a designation never does

**R42 [A]. A designation names who executes, never which gate applies.** One protocol sentence says
so, because the owner's phrase "spravi tu pracu co treba" / "does whatever work is needed" is
stretchable. Unchanged and **not inheritable by a designee**:

- `parley protocol publish` stays TTY-attended and is never worked around;
- **core publication stays owner-attended**;
- channel publication stays the organizer's attended release step.

**R43 [A].** Also unchanged: **quorum, roster, signoff weight**, model-diversity rules in kind, and
`roles:` semantics (`COOPERATION.md:311` — advisory lenses only; they do not change quorum, signoff
weight, artifact ownership, drafter eligibility or roster membership).

**R44 [A]. Independent review is preserved verbatim.** Reviewers are all non-implementers on
`deliberation` (`COOPERATION.md:234`); no merge or completion without an invokable non-implementer
reviewer (`:538`); the goal-done check is by a **fresh non-implementer** (`:689`), and under
LE-7/LE-11 it can only **withhold** a close, never establish one. The designee remains a full design
participant and signatory, and **never reviews or goal-checks its own work**. These requirements are
**independent of every gate this idea adds**: the new gates are additive, and none of them is a
precondition for, or a substitute for, code review or the goal-done check.

### I. Observability

**R45 [D]. OPEN-2, as signed — the operative rule. Designation-only emission.**
`agent.implementer_resolved {idea, implementer, source}` and the extra stdout line fire **whenever a
designation is present at tier 2 or tier 3** — **including** when a present designation was
inapplicable, opted out, or unavailable and the run therefore fell through, because that
fall-through is exactly what must be recorded (R19 concession 2). **On a deck where neither field is
set, nothing new is emitted** and the existing `driver: implementing via %s ...` line
(`internal/app/driver_impl.go:217`) stays **byte-identical**. ALT-11 (always-on emission) is
**rejected**.

The boundary sentence is the owner's, not the participants': "Ship the mechanism with the global
default UNSET, so behaviour without it is exactly today's" (`source-context/owner-brief.md`).
Designation-only is the only variant under which the unset path is byte-identical. An
`## Agreed trade-offs` entry records a decision; **it cannot supply authority an exact owner
sentence withholds.** zcode-1 owns the correction that settled this: its round-03 switch to
always-on rested on "the group now reads 3-0 always-on", which read kimi-1's round-**02** table
while kimi-1's round-03, written in parallel, had moved to designation-only — **no 3-0 ever
existed**, and zcode-1's own rule ("an unset-path behaviour change needs explicit owner authority;
participant agreement is not that authority") binds its switch symmetrically.

**R46 [D].** The event is guarded exactly as `checkModelDiversity` guards its own —
`if o.base.Store != (store.Store{})` (`internal/app/driver_impl.go:195`) — and is emitted in
`Implement` immediately before `runner.RunImplementation` (`:218`).

**R47 [A].** **The source enum MUST keep these distinct values**, so no consumer can confuse them:
`pin` (rank 1), `designation` (rank 2), `global-default` (rank 3), `none` (an explicit opt-out at
either layer), and a **fall-through** value carrying its cause — unavailability versus
inapplicability — distinguishable from both an opt-out and an unset deck. The exact constant names
are an implementer decision (register item F5); **the distinctions are not.**

**R48 [A]. The cost of R45 is recorded, not hidden.** The positional-selection defect of R49 stays
as invisible on the unset path as it is today: this deferral ships with **no in-band mitigation**.
The always-on variant becomes the **first item** of the follow-up slug, where the owner can relax
the boundary deliberately. This link is recorded here on purpose rather than inherited by silence.

### J. What is deliberately not repaired here

**R49 [A]. The legacy inheritance defect is NOT repaired in this idea.** On undesignated decks the
documented drafter-implements chain still does not run and dispatch remains list order. This is the
defect that motivated the idea, and it is left unrepaired **deliberately**, because repairing it
changes unset-path selection, which the owner's sentence forbids. Specifically **out**:

- **ALT-5 — write-end emission** of `implementer:` into new FINALs and the driver draft prompt.
  **REJECTED (deferred).** It would change dispatch on ordinary undesignated ideas through the
  protocol's own volunteer-drafter route (`COOPERATION.md:409`, `Drafter: yes`) — which fired in the
  immediately preceding run — making the first-mover chain mechanically authoritative. zcode-1
  withdrew the last hold on this in round 03, binding its own owner-boundary rule symmetrically.
- **ALT-6 — widening the resolver's `FINAL.md` read set** to `{author, drafter, finalized-by}`.
  **REJECTED (deferred).** Same class: it changes who implements on existing decks that set nothing.
  **No FINAL-inheritance repair ships in this idea.**
- **ALT-16 — `Complete` (`internal/app/driver_impl.go:518` onward) recording the `implementer:` it
  already knows. DEFERRED.** It would close the missing-pin gap measured below, but it changes
  unset-path artifacts.

All of these, plus the claim reader of R14 and the always-on event of R48, belong to **one**
follow-up slug: `meta-protocol-change-implementer-inheritance-repair`. Splitting them across
separate follow-ups invites half a repair.

## Purpose / user-visible outcome

**The owner's question, answered.** "Who does it now? Is it prescribed in the protocol?" —
**prescribed in prose, not implemented in the tooling.** The answer must always be given scoped,
never bare:

- **Prose:** `author: user` routes to "the first agent to have submitted a round-01 file"
  (`COOPERATION.md:409`), and then "**The default implementer is the FINAL drafter**" (`:443`).
- **Shipped tooling:** neither hop is implemented. Both resolver copies read
  `IMPLEMENTATION.md{implementer}` then `FINAL.md{implementer, drafted-by}` and nothing computes a
  first filer; the driver's tail is `eligible[0]` (`internal/app/driver_impl.go:131-133`) and the
  consensus copy's tail is `""` (`internal/consensus/consensus.go:778`). **Dispatch is list order.**

The unscoped sentence "speed decides the implementer" is **false of the code** and must not be
repeated bare.

**User-visible outcome after this ships.**

1. An owner who wants a specific agent to execute a specific idea writes one line —
   `implementer: kimi-1` in that idea's `00-prompt.md` — and the driver dispatches that agent, or
   stops with a pre-built `Confirm` if it cannot.
2. An owner who wants a standing preference sets `default_implementer` once in
   `~/.parley/agents.toml` (or any deck layer) and every idea that does not override it follows,
   while any idea can decline with `implementer: none`.
3. An owner who sets nothing sees **no change whatsoever** — same dispatch, same stdout, same
   events, same expected round sets.
4. In every case the **others still evaluate the work**: reviewers remain the non-implementers, the
   goal-done check remains a fresh non-implementer, and the designee never reviews itself.

**The global product default ships UNSET.** No agent is selected as anybody's persistent default by
this artifact or by the code it specifies.

**Separately, and after shipment:** the owner has **personally chosen `codex-1`** as the
machine-global default implementer (direction of 2026-09-25, `source-context/owner-default-2026-09-25.md`,
mirrored in `inbox/user-to-codex-1_meta-protocol-change-designated-implementer_default-implementer.md`).
Setting it is **post-release owner configuration performed through the shipped mechanism** — not a
shipped default, not an organizer choice, not a participant recommendation. The brief's only
anticipated owner question is therefore **closed by the owner**, and **no participant owes a
default-implementer recommendation.** Every unset-path guarantee above is unchanged by it.

**This run is expressly preserved and unchanged.** The owner's 2026-09-25 direction states that it
"does not rewrite existing deck membership, historical idea artifacts, signatures, or runs already
in flight". For this idea: `codex-1` remains the pure organizer; participants remain `claude-1`,
`kimi-1`, `zcode-1`; **claude-1 drafts FINAL, kimi-1 implements, claude-1 and zcode-1 review** —
Anthropic plus Zhipu reviewing a Moonshot implementation, model-diverse by construction. The owner's
new defaults for **new** runs (organizer `claude-1`, implementer `codex-1`, quorum `codex-1`,
`kimi-1`, `zcode-1`) are roster and orchestration policy **outside this delta**; R43 stands and this
idea changes none of it.

**Seats are a this-run assignment under the protocol in force, not evidence for anyone's default.**
kimi-1's seat here must not be offered as a recommendation, and the owner's choice of `codex-1` must
not be read back into this run. Conflating the two would be the quiet self-selection this idea
exists to remove, in either direction.

## Context & orientation

**Track and why.** `deliberation`, forced twice over: §4.0 lists protocol change (§7) and
`auto_implement` as independent triggers, and this idea is both. Speed came from driver
orchestration and bounded rounds, not from a lighter track.

**Where things live.**

- CLI worktree (this one): `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer`,
  branch `designated-implementer`. Its deck is `parley-deck/`.
- Skill worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill`,
  same branch.
- `protocolRole: source` (`parley-deck/meta/version.json`). The deck's own `COOPERATION.md` is the
  live authority; the global core is downstream. The stale `2.12.0` / `protocolSha256 12e4b31c…`
  metadata in that file is **advisory drift, already disclosed**, and is never current authority
  evidence. No protocol replacement or metadata refresh was performed in this run.

**The three protocol copies, measured.** They differ **only** in project zones:

| Copy | SHA-256 | Bytes | Lines |
|---|---|---:|---:|
| `parley-deck/COOPERATION.md` (deck, authority) | `8ce83cde…9db7` | 109,928 | 1,386 |
| `internal/protocol/defaults/COOPERATION.md` (CLI guarded) | `cce5d7d9…9b8d` | 109,682 | 1,379 |
| skill `skills/parley-deck/references/COOPERATION.md` | `fc907e59…2c9f` | 109,772 | 1,379 |

`diff` between any pair yields hunks at **only** three places — the header lines 3 and 5–7, and the
§2 host-handle rows (deck 154–159). **From deck line 160 / the other copies' line 153 to end of
file, all three are byte-identical** (`tail -n +160` and `tail -n +153` all hash to
`f67af415a13c7fd6dad84fd8cae97e388b98d6c954d693696806a746f7872291`). So one identical hunk applies
to all three at a constant offset: **−1** for lines 8–153, **−7** for lines 160 onward.

**Target sections and their deck line numbers** (subtract 1 or 7 per the rule above for the other
two copies): §0 `[defaults]` sentence at **59**; §4.0 Phase-0 kickoff template at **278**; §4 Phase
4 at **405**; Phase 5 at **441**; §10 TL;DR at **829**; §9.0 readiness at **847**.

**Skill companion — a surface the consensus enumerated incompletely.** See
`## Changes since the signed consensus`, finding (1): `skills/parley-deck/SKILL.md` restates no
Phase-5 implementer rule (true, and verified again here), but the companion reference
`skills/parley-deck/references/ROSTER_AND_PROTOCOL.md:64` **does**: "Phase 5 implementation: default
implementer is the FINAL drafter unless another participant claims it; implementation must follow
`FINAL.md`; deviations go into `IMPLEMENTATION.md`." That one line is **in scope** (R57).

**Prior-run inheritance, assessed rather than inherited.** `source-context/prior-release-handoff.md`
reports incomplete channels and an unresolved Windows scope decision from the preceding release.
Neither is treated as success and neither grants authority to fix unrelated platform architecture.
The owner has since settled the Windows question (see `## Core staging and release constraints`), so
nothing is left for participants to decide there.

**Driver and dispatch gaps, recorded and not routed around.** Three are live and this artifact was
produced under the owner-authorized documented manual fallback because of them:

1. The driver halts cross-review on an unavailable historical worktree under
   `…/scratchpad/f2repo`, and its historical-worktree accounting gate refuses cross-review
   (`organizer-notes.md`, Phase 1→2 and Phase 2;
   `inbox/claude-to-user_meta-protocol-change-designated-implementer_driver-error.md`).
2. `continue --json` on the seeded run still plans round-02 although round-03 is complete on disk;
   its prior cross-review accounting failure is unresolved. **Do not attempt an auto-continuation
   that could replace the drafter or overwrite an existing canonical artifact.**
3. The consensus-draft CLI scaffolds markdown only — it does not launch a drafter.

No accounting, worktree history, budget state, run event or code was altered to route around any of
them, and the unrelated historical idea whose facilitator declaration blocks preflight was not
touched.

## Observable acceptance criteria

Checkable by a reviewer or the driver without judgement. **No criterion below cites an unresolved
verdict** (see `## Verdict status`), and none depends on any unverdicted drafter claim.

**AC-1 — Three-copy fidelity.** A byte comparison of the changed region across all three
`COOPERATION.md` copies shows **identical hunks**, and the pre-existing project-zone differences are
unchanged: after the change, `tail -n +160` of the deck copy and `tail -n +153` of the other two
still hash to one common value, and the three baseline diff hunks (header lines 3, 5–7; deck
154–159) are the **only** remaining differences.

**AC-2 — Changelog.** `parley-deck/meta/protocol-changelog.md` gains an entry in the §7 shape
(date, `Idea:`, `Drafted by:`, `Summary:`).

**AC-3 — Skill companion consistency.** `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md`'s
Phase-5 line states the amended rule, and a grep of the whole skill tree shows **no remaining
Phase-5 default-implementer restatement that contradicts the amended protocol text** outside the
guarded `references/COOPERATION.md`. The reproducible check is
`grep -rn "default implementer" skills/ lib/ bin/` in the skill worktree; every hit must be either
the guarded copy or consistent with it.

**AC-4 — Unset-path invariance.** On a deck setting neither `implementer:` nor
`default_implementer`, all of the following hold:

- `internal/app/app_test.go:1557`, `internal/consensus/roundgate_test.go:102` and `:111` pass
  **unmodified**;
- the `driver: implementing via %s ...` stdout line is **byte-identical**;
- `ExpectedRoundParticipants` returns the **same set as before the change** for design rounds,
  review rounds, the unresolvable case and the FINAL-drafter case;
- the event stream contains **no `agent.implementer_resolved` record**.

**AC-5 — Four-state parse.** Tests demonstrate all four states of R2, including the malformed
inline-comment value of R4 failing closed (never trimmed to a bare id) and the quoted forms of R5
being accepted.

**AC-6 — The opt-out does not collapse.** A test demonstrates that `implementer: none` **suppresses
tier 3** and is distinguishable from an absent key.

**AC-7 — Config precedence.** Tests demonstrate: the same key set in two layers resolves to the
**higher non-empty** layer; an empty value at a higher layer **does not clear** a lower layer's
value; `"none"` at a higher layer **suppresses**.

**AC-8 — Tier-2 gate and its exits.** A test demonstrates the tier-2 unavailability gate blocking,
and each of R18's three exits clearing it.

**AC-9 — Dormant designation paths.** Tests demonstrate that **each** tier-3 state falls through to
today's chain **without gating**: key absent; value naming a **non-participant** (R22); value
`none`; value naming the idea's **own declared non-participating facilitator** (R20); and value
naming an eligible participant that **fails the §9.0 ping** (R19). For the present-designation cases
the designation-only event is emitted with a source value that distinguishes the cause (R47); for
the absent case nothing is emitted.

**AC-10 — Invalid tier-3 still fails.** A test demonstrates that a **malformed** tier-3 value keeps
its hard failure and does not fall through (R19 concession 3).

**AC-11 — Pin conflict escalates.** A test demonstrates that a tier-1/tier-2 disagreement
escalates — never silently honouring either — and that the `implementer_reassigned:` path re-pins
only after the confirmed record exists.

**AC-12 — Degenerate drafter fallback.** A test demonstrates that when the designee is the **only**
eligible drafter, drafting falls back to the designee rather than failing.

**AC-13 — Gate scoping.** A test demonstrates that a designation defect in idea A **does not gate**
a run of idea B, and that no new gate iterates all ideas.

**AC-14 — Two-participant designated kickoff is a warning.** A test demonstrates it warns and
proceeds; it does not block.

**AC-15 — Review exclusion reads the record, not the instruction.** A test demonstrates that with a
tier-2 designation naming agent X and an `IMPLEMENTATION.md` pin naming agent Y, review-round
expectations exclude **Y** — and that `ExpectedRoundParticipants` is unaffected by the designation
in every state (R32).

**AC-16 — Inertness.** `grep -rn "impl-claim" --include="*.go" .` still returns **nothing**, and no
new `00-prompt.md` writer exists — a grep for writes to that filename finds only `CreateIdeaFull`,
`updateIdeaStatus` (`internal/consensus/consensus.go:898-930`) and `pipeline.SeedBlockPrompt`.

**AC-17 — Whole-tree health.** `go build ./...`, `go vet ./...` and `go test ./...` are green, and
`gofmt -l` is clean on every changed file.

**AC-18 — Designation-path demonstration, end to end.** One fixture run with `implementer:` set
dispatches the designee; one fixture run with `default_implementer` set at **two** config layers
dispatches the **higher** layer's id and records the source.

**AC-19 — Ships UNSET.** A freshly generated central config contains `default_implementer`
**commented out** (R12), and a deck created by the tooling contains the key **not at all**.

**AC-20 — Attended boundaries intact.** `parley protocol publish` still refuses without a
controlling terminal, and no code path in the delta allocates a pty, drives publication, or performs
a channel action.

**AC-21 — Independent review actually happened.** Review round 1 carries a review file from **each**
non-implementer (`claude-1` and `zcode-1` for this run) with a populated `## Refutation attempts`
section, and the goal-done check is performed by a **fresh non-implementer** (`COOPERATION.md:689`).
Neither is satisfiable by the implementer, by the organizer, or by any gate this idea adds.

## Test requirements

Tests are **participant-owned and meaningful**: each one must fail if the corresponding rule is
removed. A test that passes against an unimplemented rule does not count.

**T-1.** Four-state parse: absent / present-empty / `none` / id — plus the malformed
inline-comment value and both quoted forms (AC-5, AC-10).

**T-2.** `none` opt-out suppresses tier 3 and does not read as absent (AC-6). This is the pin without
which the opt-out silently disappears.

**T-3.** Config layer precedence, three cases: higher non-empty wins; `""` does not clear; `"none"`
suppresses (AC-7). At least one case must span **two real layers**, not one file.

**T-4.** Tier-2 unavailability gate plus each of its three exits (AC-8).

**T-5.** Pin-versus-designation escalation, and the `implementer_reassigned:` re-pin path (AC-11).

**T-6.** Degenerate drafter fallback: designee is the only eligible drafter (AC-12).

**T-7. Dormant designation paths — all five tier-3 states of AC-9**, each asserting (a) no gate,
(b) fall-through to today's resolved id, and (c) the emitted source value where a designation was
present.

**T-8. Legacy compatibility.** The three existing tests of R33 pass **unmodified**; plus a new test
asserting `ExpectedRoundParticipants` parity across design rounds, review rounds, the unresolvable
case and the FINAL-drafter case, and a test asserting **no** `agent.implementer_resolved` event on
an unset deck (AC-4).

**T-9.** Gate scoping across two ideas in one deck (AC-13).

**T-10.** Two-participant designated kickoff warns rather than blocks (AC-14).

**T-11.** Review exclusion reads the pin, not the designation (AC-15).

## Exact in/out boundaries

### In scope — files that change

1. **`internal/protocol/` (new shared code).** The key constant for `implementer`; the four-state
   designation parser over `00-prompt.md` frontmatter (R2–R6); the shared resolution chain returning
   `(id, source, ok)` and taking an ordered source list plus the eligibility list (R31). Home is
   `internal/protocol` because both callers already import it.
2. **`internal/config/runtime.go`.** `default_implementer` added to the `[defaults]` struct pair
   (`:283-291`, `:315-329`); merged in `mergeDefaults` (`:533-576`) on the **non-empty string**
   pattern, not a pointer; emitted **commented out** in `centralDefaultTemplate` (`:627` onward).
3. **`internal/app/driver_impl.go`.** Resolve through the shared chain in `newDriverImplOps` after
   `eligible` is computed (`:54-64`); carry designation gate failures on the existing `roleErr` path
   so all four role actions escalate (`:214`, `:278`, `:407-408`, `:504`); in `Implement`
   (`:213-223`) emit the designation-only resolved-source event immediately before
   `runner.RunImplementation` (`:218`), guarded as at `:195`, plus the re-entry comparison of R29
   using `store.Store.Load()`; add the kickoff diversity check of R37 **without changing**
   `checkModelDiversity` itself.
4. **`internal/app/driver_consensus.go`.** `firstEligibleHeadlessAgent` (`:109-123`) gains the
   designation-aware drafter preference of R36.
5. **`internal/consensus/consensus.go`.** Replace the duplicate `resolveImplementer` (`:751-779`)
   with a call into the shared chain, passing `[IMPLEMENTATION.md, FINAL.md]` and **the raw
   participants list it was handed** (`:734`). `expectedRoundParticipants` (`:730-748`) keeps its
   behaviour exactly.
6. **`internal/app/preflight.go`.** Optional early **idea-scoped** copy of the **validity** gate
   only, using the existing `gate{Kind, Detail, Confirm}` shape (`:334-338`). **Must not** iterate
   all ideas.
7. **Three `COOPERATION.md` copies** (R40 and the sections listed under `## Context & orientation`),
   **plus `parley-deck/meta/protocol-changelog.md`**.
8. **`skills/parley-deck/references/ROSTER_AND_PROTOCOL.md`** — the single Phase-5 line at `:64`
   (R57).
9. **New tests** per `## Test requirements`.

### Out of scope — what must not change

`checkModelDiversity`'s logic and severity; the `driver: implementing via %s ...` stdout string;
`ValidateFinal` / `RequiredFinalSections` / `finalScaffoldReason`; the `impl-claim` file's
inertness; `roles:` semantics; quorum, roster and signoff weight; `expectedRoundParticipants`'
observable output in any state; both resolver tails; every attended gate; `updateIdeaStatus`
(`internal/consensus/consensus.go:898-930`) and its callers; and **no `00-prompt.md` writer is added
anywhere** — the live read replaces materialization precisely so that no tool writes a designation
into an author-owned artifact.

Also out: the write-end FINAL emission (ALT-5), the FINAL read-set widening (ALT-6), the `Complete`
pin record (ALT-16), any machine reader for `impl-claim` (ALT-4), always-on event emission
(ALT-11), extending tiers 2–3 into the review-exclusion read (ALT-12), unifying the eligibility
filter by changing the consensus call site (ALT-13), a tier-3 unavailability gate (ALT-14), a
uniform tier-2/tier-3 facilitator hard gate (ALT-22), a code-level `excluded:` cross-check (ALT-21),
the `facilitatorConflictGates` all-ideas repair (ALT-17), the trajectory-policy coupling (ALT-9),
and any Windows platform architecture work.

**Non-goals, restated.** No quorum or roster change. No model or effort downgrade. No unrelated code
changes. No change to the 1.49.0 organizer rules beyond what this role needs. **No persistent
default implementer selection** in the shipped product.

## Core staging and release constraints

**R50.** The release order is the owner's, set in
`inbox/user-to-codex-1_meta-protocol-change-designated-implementer_release-order.md` (copied to
`source-context/release-order.md`, Slovak original preserved there as source evidence):
`release-1.49.1` → **this idea** → `windows-portability`, each gated on its predecessor's done file.

**R51.** The predecessor gate is **satisfied**: `codex-1-to-user_release-1.49.1_done.md` exists in
the lean-organizer worktree deck inbox (verified for this FINAL), and `origin/main` is at `c49b464`
"complete independently audited release 1.49.1" with tag `v1.49.1` present. **This is the release
*order* gate only.** This idea still owes, before any release step: a reviewed implementation
(Phases 6–8 to zero agreed fixes), a fresh non-implementer goal-done check, and **independent
channel verification by a participant** of every channel actually used.

**R52. Windows, as the owner decided.** Release notes label Windows **experimental/unvalidated**;
labelled Windows assets are **retained**; **no CLI winget PR** is opened — the owner holds it until
`windows-portability` ships. **Skill channels are unaffected**, so the skill's own winget path is
normal. **No Windows platform architecture work happens in this idea**: the delta adds no
build-tagged file and no new platform surface, and `internal/app/driver_impl.go:543` already refuses
independent-evidence completion on Windows.

**R53. Core staging base, verified.** Global core **2.13.0 is published live**:
`~/.parley/protocol/core/2.13.0/COOPERATION.md`, SHA-256
`fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, 109,772 bytes, mode `0444`
(write-once). It is **byte-identical** to the staged `~/.parley/staging/COOPERATION-2.13.0.md` and to
the skill worktree's bundled references copy, so the base is unambiguous. Stage the next core as
**that published base plus exactly this idea's reviewed hunks — nothing else.**

**R54. Publication stays attended.** `parley protocol publish --version V --from FILE` requires a
controlling terminal. The exact command is **escalated to the owner**; a TTY is **never** allocated
to work around the gate, and no participant invokes core publication.

**R55. Version selection is an organizer act at staging time**: the next minor **above actual
releases**, re-verified then rather than fixed here. Observed at FINAL time: CLI released `v1.49.1`;
skill `2.13.0`; core `2.13.0`.

**R56. Integration.** Integrate the latest `origin/main` **only after** the reviewed
implementation. This run uses **local canonical files and direct main integration, without
development PRs**, as an explicit per-run override of the github-pr transport mechanics; the
**global transport header is not changed**. Merging to main, channel publication and any global
configuration change are the **organizer's attended release step** — no participant performs them,
and none of them is authorized by this FINAL.

**R57. The skill companion edit** (`skills/parley-deck/references/ROSTER_AND_PROTOCOL.md:64`) is a
**one-line consistency edit in the skill worktree**, on the same branch. It changes no skill
behaviour, no installer, no manifest and no channel.

## Idempotence & recovery

**What state matters.** (1) `00-prompt.md` `implementer:` — author-owned, never machine-written.
(2) `[defaults].default_implementer` in up to four config files — owner-owned. (3)
`IMPLEMENTATION.md` `implementer:` — the pin, written once at Phase 5. (4) The durable event stream.
(5) The working trees of the two worktrees.

**Safe to rerun, unchanged, any number of times.** Every read in the chain (R13); every validity and
availability check; `go build` / `go vet` / `go test`; the three-copy hunk comparison; every grep in
the acceptance criteria. Resolution is recomputed at each `newDriverImplOps` construction and frozen
for that dispatch (R29), so re-running a dispatch after a config edit re-resolves rather than
resuming a stale decision — and a **recorded** dispatch under a designation that no longer matches
**escalates** instead of silently switching.

**Idempotent by construction.** The text edits are hunk-shaped: applying them twice is detectable by
AC-1 and by `git diff`. The config field is additive; re-emitting `centralDefaultTemplate` on an
existing file is already a no-op (`EnsureCentralDefault` returns early when the file exists,
`internal/config/runtime.go:610-611`).

**Recovery.**

- **Implementation abandoned mid-Phase-5:** the partial tree is evidence. Reassignment over it is a
  human decision (R23). Do not delete or rewrite it.
- **A gate fires wrongly:** the exits are pre-built (R18) and recorded. Prefer the recorded exit to
  removing the gate.
- **A designation conflicts with a pin:** R28 — escalate, then either restore the designation or
  record `implementer_reassigned:`. Never let the incoming implementer self-appoint.
- **Text drift between the three copies:** AC-1 detects it; re-apply the identical hunk rather than
  hand-editing one copy.
- **Core staging wrong base:** recompute against the published `fc907e59…` (R53) and re-stage. A
  published core release is write-once and is never overwritten.

**Human gates — never automated, never inherited by a designee.** Core publication; channel
publication; `parley protocol publish`; merging to main; any global configuration change (including
setting the owner's `default_implementer`); `implementer_waived:`, `excluded:` and
`implementer_reassigned:` confirmations.

## Known risks / de-risking

**K-1 — The mechanism ships without the legacy inheritance repair (R49), and with no in-band
mitigation (R48).** On undesignated decks the documented drafter-implements chain still does not run,
dispatch remains list order, and — under designation-only emission — nothing makes that visible on
the path where it happens. *De-risking:* the whole repair is one named follow-up slug with the
always-on event as its first item, where the owner can relax the "exactly today's" boundary
deliberately. Recorded rather than inherited by silence.

**K-2 — `parley init` emits `default_implementer` commented out**, deviating from
`centralDefaultTemplate`'s active-keys-with-inline-comments shape. *De-risking:* deliberate — an
active key with a value would ship the global default **SET**, against the owner's sentence. One
documented sentence of protocol text; AC-19 pins it.

**K-3 — A designation concentrates the heaviest work of a run in one agent.** All three participants
raised this independently. *De-risking, all already shipped or inside the delta:* drafter separation
(R36); the kickoff diversity check (R37); reviewers = all non-implementers on `deliberation`
(`COOPERATION.md:234`); the invokable-non-implementer-reviewer floor (`:538`, R39); the fresh
non-implementer goal-done check (`:689`); and seat rotation — this run moves every seat relative to
the previous one. **The residual risk is real and is the owner's to accept when setting a default.**
Note for honesty: the cost framing behind this risk rests on an **UNVERIFIED** ranking claim (see
`## Verdict status`, VC-2); the risk stands on the three participants' independent judgment, not on
that ranking.

**K-4 — The X-5 eligibility divergence is preserved rather than repaired (R34).** *De-risking:* it
becomes an explicit parameter and is test-pinned, so it is visible rather than accidental; unifying
it would change unset-path behaviour in a narrow but **reachable** state.

**K-5 — Four config layers, not two.** No participant asked for `agents.local.toml` and
`$PARLEY_HEADLESS_AGENT_CONFIG` to carry this key; both peers wrote "two levels only" in round 02
and withdrew it. *De-risking:* the key inherits all four by riding `[defaults]`; reading machine-only
would make this the one `[defaults]` key that ignores the deck. Accepted as inherited, stated in the
protocol text, pinned by T-3.

**K-6 — Unexamined surfaces (blind spots, named not claimed).** Nobody traced: (a) **pipeline-block
dispatch** — a designation written into a block prompt seeded by `pipeline.SeedBlockPrompt`
(`internal/pipeline/executor.go:302`), and whether tier 3 applies to block dispatch at all
(`newDriverImplOps` has exactly three call sites, all in `parley run` / `parley continue`:
`internal/app/app.go:1252`, `:1984`, `:2038`); (b) **the TUI** — `implementerOf`
(`internal/tui/protosnap.go:383-394`) reads `implementer` among `{"agent","implementer","by"}` from
`IMPLEMENTATION.md`, and nobody checked whether a designation or a stale implementer displays
correctly; (c) **`participants:` membership change against a live designation** — the rounds covered
`--participants` *order* but not membership, so an idea that ran fine can gate on its next launch if
the designee is removed. *De-risking:* each is on the deferred register as unexamined surface, not as
a claimed defect. R25's gate siting rests on the measured claim that `roleErr` is reached on the
launch routes that matter; **if it is not reached on some route (resume, TUI, pipeline), the gate is
unenforced exactly where it matters, and (a) is the place to look.** The implementer must re-measure
this rather than accept it.

**K-7 — The recommendation-adequacy gap, recorded.** Across three rounds **no participant proposed a
criterion** for when the evidence base behind a default-implementer recommendation would be
*sufficient* rather than merely assembled (the deck's role evidence is one prior implementation,
labelled n=1 by its own author, one prior review, and locator accuracy over three artifacts). The
owner answered the question directly on 2026-09-25, which closes the deliverable — but the gap is
worth naming, and R20 shows that the owner's choice immediately surfaced an interaction none of the
three had modelled.

**K-8 — This FINAL carries two drafter findings no peer has reviewed** (the skill-companion surface
and three locator corrections). *De-risking:* both are flagged as such in
`## Changes since the signed consensus`, both are PRIMARY-evidenced and independently re-checkable
in one command, and both are review-round business. A reviewer who disagrees should file a finding.

## Verdict status

Per §15.2/§15.3. **No acceptance criterion in this FINAL is supported by an unresolved verdict, and
no rule depends on one.**

**VC-1 — RESOLVED, no longer `DISPUTED`.** The claim is owner-brief testimony that claude-1's first
code review of the preceding run "found 1 CRITICAL and 8 MAJOR findings in the implementation, all
fixed in later cycles". The **counts** were CONFIRMED by all three with PRIMARY evidence and were
never in conflict. The conflict was the **"all fixed"** clause, on which zcode-1 held a round-01
`UNVERIFIED` ("per-finding trace not performed"). **zcode-1 withdrew that verdict in its signoff**,
having performed the trace itself with PRIMARY evidence, and recorded that its premise "was true of
the round-02 record and false of the round-03 record" — it was written in parallel with two
independent traces it could not have read. kimi-1 filed an independent PRIMARY trace in round 03;
claude-1 owns Verdict A and its own round-03 trace and issues no verdict on either peer's.

**The clause is admissible only scoped**, and two precisions MUST travel with any citation of it:
**nine findings produced eight fixes** (MAJ-1 and MAJ-2 merged into F2), and **not every disposition
was a repair** (F1 *removed* the disputed change). **"Fixed" means dispositioned-and-closed through
the protocol's own gate.** The **code-content efficacy of each fix was never re-audited, and no
party claimed it** — this FINAL makes **no code-efficacy claim** about that prior run, and none of
its rules or criteria rests on the clause. It feeds only the post-release recommendation evidence
base, which the owner's direct choice has closed.

**VC-2 — RESOLVED as `UNVERIFIED`, and it stays that way.** The claim "implementation is the
heaviest token work of a run" is a **ranking** that nothing in evidence supports. The copied
`source-context/organizer-token-study.md` supports **~30–40% of an *implementing organizer's*
input**, with its own caveats. Under §15.2 the ranking **may not enter this FINAL as established**,
and it does not: the percentage appears only with its caveat, K-3 explicitly does not rest on it,
and no acceptance criterion cites it. Evidence asymmetry recorded: claude-1 read only the in-deck
copy; zcode-1 read the external original — cite zcode-1's reading for the original's numbers.

**VC-3 — RESOLVED, and load-bearing.** "Speed decides the implementer" is **CONFIRMED of the
protocol prose** (`COOPERATION.md:407` → `:409` → `:443`) and **WRONG of the shipped tooling**
(dispatch is list order; no code computes a first filer). `## Purpose / user-visible outcome` states
it scoped, never bare. R49 and the purpose statement both depend on this resolution.

**VC-4 — RESOLVED, with a standing rule.** Locator corrections across the rounds were found by
measurement, never by argument. **No locator in this FINAL was transcribed**, and no locator in the
implementation may be transcribed from any round artifact, the consensus, or this FINAL without
re-measurement.

**Unverdicted drafter claims, carried with their provenance.** Under §15.1 claude-1 owns these and
may not verdict them; each is recorded so a peer can verdict it, and **no rule or criterion depends
on any of them**:

- **The pin census (claude-1's N-5).** Re-measured for this FINAL at HEAD, same result: **75**
  `IMPLEMENTATION.md` files under `parley-deck/ideas/`, **63** carrying a non-empty `implementer:`
  pin, **12** not, and of the 75 exactly **9** with neither a readable pin nor a readable `FINAL.md`
  key — `antigravity-agent-migration`, `consensus-request-signoffs`, `consensus-workflow-cli`,
  `continuous-run-tui`, `repo-map-mvp`, `roadmap-implementation-plan`,
  `roster-operations-standard`, `tui-action-execution`, `tui-workspace-sessions`. **No peer has
  verdicted this.** R32's disposition does **not** depend on it — the reason for R32 is categorical
  and holds at any census value; the census only quantifies what R49's deferral leaves open. If a
  peer re-measures differently, the number changes and the rule does not.
- **The skill-companion finding** at `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md:64`.
  Tagged PRIMARY (grep run, line quoted). New in this FINAL; **no peer has seen it.**
- **claude-1's N-2/N-3/N-4** (the only ordinary-idea kickoff writer is `CreateIdeaFull`; preflight
  runs before the idea exists; materialized provenance cannot live inside the frontmatter fence) are
  first-canonical and unverdicted. R25's gate siting rests on the preflight-ordering measurement and
  **must be re-measured by the implementer.**
- **zcode-1's §12.5 pipeline-scoping claim** carries no peer verdict. R10 does **not** depend on it:
  materialization is rejected on the `LoadDefaults` mechanics, the absence of a hook for
  hand-authored ideas, and the provenance constraint — each independently sufficient.
- **kimi-1's locator-accuracy record** across rounds 01–03 is the only directly relevant measured
  evidence that bore on the implementer choice, and all three files label it a **small sample**. It
  is not cited in support of anything here.

## Alternatives and correlated priors

Per §15.6. **This FINAL contradicts no adoption recorded in the consensus `## Alternatives
disposition` table.** ALT-1 … ALT-22 are all honoured in the text above; the five that were
*proposed* rather than settled at drafting time — ALT-11 (OPEN-2), ALT-12 (OPEN-3), ALT-14 (OPEN-1),
ALT-21/ALT-22 (OPEN-4, OPEN-5) — were each ratified by **all three** signoffs in the direction R19,
R20, R24, R32 and R45 state, and each appears above with its rejected counterpart named.

**§15.6(b) — where nominally independent proposals are one family.** The agreement behind this FINAL
is **not** three independent confirmations:

- **Shared inputs.** All three participants read the same byte-identical protocol packet
  (`8ce83cde…`), the same owner brief, the same code at the same SHA, and each other's files in
  full. Convergence under identical inputs is **weak** evidence of correctness.
- **One author's design, verified by two peers — not three designs converging.** R10/R11 (the live
  layered read), R13 (the single chain), R31/R34 (the unification shape), R36/R37 (drafter
  separation, kickoff diversity) and R41 all **originate in claude-1's round-01/02 proposals** and
  were adopted by the peers with their own verification. Genuine verification — but of one author's
  design. The clearest cases of a position surviving because a **different** participant built it
  are R2–R3's `none` state (zcode-1) and R18's tier-2 fail-closed rule (kimi-1's and zcode-1's
  independent §9.0 reading).
- **A demonstrated adoption cascade, in this run's own record.** Phase-0 materialization was
  proposed in claude-1's round-01, **both** peers adopted it in round 02 explicitly citing that
  proposal, and **both** withdrew it in round 03 when claude-1 withdrew it. The round-02 unanimity
  was an artifact of one participant's proposal travelling — **and it was wrong**, by all three
  participants' later judgment. Any unanimity here that traces to a single originator deserves the
  same suspicion.
- **A demonstrated stale count.** zcode-1's round-03 concluded "the group now reads 3-0 always-on"
  while kimi-1's round-03, written in parallel, had moved the other way. Anyone reading zcode-1's
  file alone would have recorded a unanimity that never existed. zcode-1 owns and corrected this in
  its signoff (R45).
- **Shared blind spots.** All three missed the `excluded:` interaction, the pipeline path, the TUI,
  the recommendation-adequacy question, and — most consequentially — the tier-3/declared-facilitator
  collision of R20, which three rounds of agreement did not surface and **one owner sentence did**.
  **Agreement says nothing about what nobody looked at** (K-6, K-7).

**What would make the agreed position wrong** — each of these refutes a load-bearing rule:

- **R10/R13:** a hook that reliably covers hand-authored ideas at Phase 0. Materialization's only
  fatal defect is that no such hook exists; show one and the freeze argument returns.
- **R8/R9:** a `mergeDefaults` path where an empty higher-layer value **does** clear a lower one —
  the `"none"` suppressor would be unnecessary and the layering story changes.
- **R2:** a `ReadFrontmatter` behaviour that does not preserve key presence — the four-state parse
  collapses to three and the typo case becomes silent.
- **R32:** a consumer of `expectedRoundParticipants` that genuinely needs the **instruction** rather
  than the **record** — for example a pre-implementation review round with no `IMPLEMENTATION.md`
  yet. Both peers checked and found none; if one appears, the round-02 unification shape returns
  with the exposure recorded as an explicit trade-off.
- **R25:** evidence that the `roleErr` path is not reached on some launch route (resume, TUI,
  pipeline) — the gate would be unenforced exactly where it matters. K-6(a) is where to look.
- **R49/K-1:** an owner instruction relaxing "exactly today's", which moves the whole inheritance
  repair back into scope.
- **The whole mechanism:** if the owner's actual intent were that the **organizer** assign the
  implementer per run rather than the **owner** designate it, the per-idea field is the wrong surface
  and the design should be re-opened. Nothing in the brief says that, and the brief's words are the
  owner's own — but **no participant tested that alternative reading.**

## Deferred follow-up register

Controlled and explicit. **Nothing here is implemented by this idea, no quorum is convened or
staffed for any of it, and listing an item is not authorization to start it.** Each entry names its
carrier.

| # | Deferred item | Carrier |
|---|---|---|
| F1 | The whole legacy inheritance repair, as one unit: (i) FINAL read-set widening to `{author, drafter, finalized-by}` (ALT-6); (ii) write-end emission of `implementer:` in the FINAL template and driver draft prompt (ALT-5); (iii) `Complete` recording the pin it already knows (ALT-16), which would close the 12-of-75 gap; (iv) a machine reader for `impl-claim` (ALT-4); (v) **first item:** the always-on `agent.implementer_resolved` variant (ALT-11) | slug `meta-protocol-change-implementer-inheritance-repair` — **not opened by this idea** |
| F2 | `--no-implement` visibility at the chosen gate site: one implementation-time trace (R17) | this idea's implementation, recorded in `IMPLEMENTATION.md`; the waiver line is the exit if it is not visible |
| F3 | A code-level cross-check of a live designation against `excluded:` (ALT-21) | a later idea; R24's protocol sentence ships now |
| F4 | Whether the early preflight copy of the validity gate ships at all (R26) — ergonomics only | this idea's implementer decision |
| F5 | Symbol and file names for the shared reader, the parser and the source enum; and the enum's exact constant names (R47 fixes the required **distinctions**, not the names) | this idea's implementer decision |
| F6 | X-5 eligibility unification (ALT-13) — parameterised and pinned here, repaired later | a later idea |
| F7 | `facilitatorConflictGates`' all-ideas iteration (ALT-17) — unrelated pre-existing defect; the binding requirement here is only that the new gate not copy it | a later idea |
| F8 | Documenting `--participants` ordering (ALT-18) | a later idea |
| F9 | Trajectory-policy / designation agreement gate (ALT-9) | a later idea |
| F10 | Unexamined surfaces of K-6: pipeline-block dispatch; TUI display of a designation or a stale implementer; `participants:` membership change against a live designation | each needs its own scoped idea; **no defect is claimed** |
| F11 | Version selection at staging time (R55) | organizer act under the owner's release order |

## Changes since the signed consensus

§15.5's role-concentration trigger does **not** fire for this artifact: the declared facilitator
`codex-1` is a pure organizer, not a participant, and this is `deliberation`, where consensus and
FINAL are separate. This section is therefore disclosure, not that duty — and it exists because a
fresh implementer and both reviewers must see exactly where this FINAL goes beyond the text all
three signed.

**(1) NEW — a skill-companion text surface the consensus enumerated incompletely.** The consensus
recorded that `skills/parley-deck/SKILL.md` restates no Phase-5 rule, and concluded that only the
guarded references copy changes in the skill tree, with a skill-tree grep as an acceptance check
"against future drift". The SKILL.md statement is **true** and re-verified. But running that grep
**now** finds **present** drift: `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md:64` restates
the very rule being amended — "Phase 5 implementation: default implementer is the FINAL drafter
unless another participant claims it…". Tag PRIMARY (grep run in the skill worktree, line quoted
verbatim). Left alone, the companion would ship contradicting the amended core, and AC-1's sibling
check would be unsatisfiable as worded. This FINAL therefore puts that **one line** in scope (R57,
AC-3). It **extends** the consensus enumeration for consistency; it **contradicts no adoption** —
AD-14 required identical hunks across the protocol copies and mandated the grep, and this is the
grep doing its job. **No peer has reviewed this finding**; a reviewer who disagrees should file it as
a finding rather than treat it as settled.

**(2) Three locator corrections, re-measured at HEAD `07ceb0a`.** Per VC-4 these were found by
measurement, not argument:

- the quote-trim line is `internal/app/driver_impl.go:125` and `internal/consensus/consensus.go:772`
  (the consensus cited `driver_impl.go:121`, which is `if err != nil`);
- the Phase-4 FINAL-validation enforcement site is `internal/driver/consensus.go:57` and `:61`
  calling `finalScaffoldReason` (`:175-189`), which delegates to `protocol.ValidateFinal` at `:185`
  — the consensus cited `internal/driver/consensus.go:55-63` as the `ValidateFinal` site directly;
- `reviewConsensusVoters` begins at `internal/consensus/consensus.go:136` (the consensus cited
  `:143`, which is the call **inside** it), and `ExpectedRoundParticipants` begins at `:1089` (cited
  as `:1090`).

**(3) Measurements added, not changed.** The three-copy byte-identity proof and the constant hunk
offsets under `## Context & orientation`; the published-core identity in R53; the re-run inertness
greps of R14 and R24; the re-measured seven-idea facilitator census in R20 and pin census in
`## Verdict status`. All returned the same results the signed text asserts.

**(4) VC-1's status advanced by zcode-1's own signoff**, not by this drafter: it withdrew Verdict B,
so the mandatory `DISPUTED` heading the consensus anticipated does not apply. The scoped clause and
both precisions are carried anyway, and no code-efficacy claim is made.

Nothing else in the signed text is narrowed, widened or reinterpreted here.

## Protocol attestation

`parley protocol packet --dir . --phase 4 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change --audience
participant --json`:

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "body_path": ".parley-runtime/protocol-packets/full-phase4-deliberation-8ce83cde….md",
  "source": { "role": "source", "transport": "github-pr", "bytes": 109928 },
  "shadow": { "packet_sha256": "26abffe9…", "packet_bytes": 73105, "included_blocks": 37, "omitted_blocks": 32 }
}
```

- **`fallback_reason` is ABSENT**, verified by enumerating the response's top-level keys
  (`body_path`, `context_mode`, `index`, `packet_sha256`, `request`, `shadow`, `source`,
  `source_sha256`) — not by reading a rendered null.
- `source_sha256` = `packet_sha256` = the kickoff authority SHA =
  `shasum -a 256 parley-deck/COOPERATION.md` = `8ce83cde…9db7`; source 109,928 bytes.
- The rendered body was **read end to end**: 1,386 lines / 109,928 bytes, and it hashes to
  `8ce83cde…9db7`, i.e. **byte-identical to the live authority**.
- The `shadow` block describes the **optimized packet that was not used**; the `index`'s
  include/omit flags belong to that shadow audit record, and the phase-4 view classifies §12 as "not
  applicable to this launch" and §13 as "read on demand" while §7 enters on `protocol_change` and
  §15 on `phase 4`. **The full body was used.**
- Transport recorded as `github-pr`. This run's local-canonical-files / direct-main override is the
  organizer's release mechanic (R56), not a transport-header change.
- HEAD at write time: `07ceb0a78b9f4ba89c171cdd9946fe24bb05930d`. `git diff --stat 73ee923 HEAD --
  internal/ cmd/ parley-deck/COOPERATION.md` is empty — no code or protocol text moved during the
  entire run, so every locator verified in rounds 1–3 remains valid as to content and every locator
  here is addressable at HEAD.

**Canonical-artifact and signature check performed before writing this FINAL** (Phase 4's
requirement that every active non-facilitator participant has created the expected canonical
artifacts):

- `round-01/`, `round-02/`, `round-03/` each contain a file from **all three** participants —
  `claude-1.md`, `kimi-1.md`, `zcode-1.md` — nine round artifacts in total.
- `consensus.md` (1,568 lines) exists with `## Signoffs` and **three participants' ACCEPT blocks**.
  `parley consensus status --json meta-protocol-change-designated-implementer` reports
  `"triage": "ready"` with `claude-1` ✅ (line 1208), `kimi-1` ✅ (1277) and `zcode-1` ✅ (1473);
  kimi-1's reaffirmation block after the signoff-sequencing incident is at 1374.
- `inbox/claude-1-to-all_meta-protocol-change-designated-implementer_drafter-volunteer.md` exists and
  was filed before the consensus signoff completed; the claude-1 signoff carries `Drafter: yes`.
- No solo exception is claimed or needed: three participants filed independently.

**What the drafter did and did not do.** Wrote this one file, then published it and set
`00-prompt.md` `status: final` through the normal command. **No** code or skill edit, **no**
`IMPLEMENTATION.md`, **no** implementation, **no** release or channel action, **no** merge or main
integration, **no** global configuration change, **no** roster change, **no** edit to any peer
artifact, round file, signature or closed record, **no** accounting / worktree / budget state
touched, **no** browser driven, **no** synthetic driver transition or forged event, and **no**
auto-continuation attempted on the stale run. No code-verification verdict is issued here: the
implementation is kimi-1's and its review is claude-1's and zcode-1's.

**FINAL is frozen on publication.** If it is later invalidated, open
`meta-protocol-change-designated-implementer-v2`; do not edit this file.

## References

- Consensus (all dispositions, signoffs, alternatives table): `./consensus.md`
- Rounds: `./round-01/`, `./round-02/`, `./round-03/`
- Owner brief (controlling, quoted verbatim with translation): `./source-context/owner-brief.md`
- Owner release order and the Windows decision: `./source-context/release-order.md`
- Owner global-defaults direction of 2026-09-25: `./source-context/owner-default-2026-09-25.md`
- Prior release handoff (inherited limitations, not successes): `./source-context/prior-release-handoff.md`
- Organizer token study (cite with its caveats; see VC-2): `./source-context/organizer-token-study.md`
- Readiness evidence: `./source-context/preflight.json`
- Organizer execution notes and driver gaps: `./organizer-notes.md`
- Organizer usage ledger: `./organizer-usage.md`, `./usage-ledger.jsonl`
- Drafter volunteer note: `../../inbox/claude-1-to-all_meta-protocol-change-designated-implementer_drafter-volunteer.md`
- Owner-default relay and organizer request to reflect it:
  `../../inbox/user-to-codex-1_meta-protocol-change-designated-implementer_default-implementer.md`,
  `../../inbox/codex-1-to-all_meta-protocol-change-designated-implementer_owner-default-update.md`
- Driver escalation preserved: `../../inbox/claude-to-user_meta-protocol-change-designated-implementer_driver-error.md`
- Protocol authority: `../../COOPERATION.md` (`8ce83cde…9db7`); changelog `../../meta/protocol-changelog.md`
- Predecessor release gate: `codex-1-to-user_release-1.49.1_done.md` in the lean-organizer worktree deck inbox
- Prior designated-implementer precedent (1.49.0 `facilitator:` / `facilitator_participates:`):
  `../meta-protocol-change-lean-organizer/`

<!-- The owner's quoted words and release direction are originally Slovak; the verbatim originals are
     preserved in 00-prompt.md and source-context/{owner-brief,release-order,owner-default-2026-09-25}.md
     alongside their English translations, per §6 rule 6. -->
