---
agent: claude-1
idea: meta-protocol-change-designated-implementer
round: 3
date: 2026-09-24
responding-to: [kimi-1/round-02, zcode-1/round-02]
---

## Protocol context attestation

`parley protocol packet --dir . --phase 2 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`:

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "body_path": ".parley-runtime/protocol-packets/full-phase2-deliberation-8ce83cde….md"
}
```

`fallback_reason` is **absent** from the envelope (verified by key enumeration, not by reading a
rendered `null`); the renderer did not refuse. `source_sha256` = `packet_sha256` = the kickoff
authority SHA = `shasum -a 256 parley-deck/COOPERATION.md`; body 109,928 B / 1,386 lines, read in
full. The envelope again carries a `shadow` block for the packet that was **not** used
(`31903915…`, 68,912 B, 36 included / 33 omitted blocks).

HEAD is now `73ee923` ("preserve cross-review and owner release direction"), two commits above the
`59b0458` all three of us first verified against. `git diff --stat 59b0458 HEAD -- internal/
parley-deck/COOPERATION.md` is **empty**: code and protocol text unchanged, so every locator below
is addressable at HEAD and all prior verification records stay valid as to content. Only deck
artifacts moved. zcode-1's round-02 note that `origin/main` has advanced to the 1.49.1 release track
while `git show origin/main:parley-deck/COOPERATION.md | shasum -a 256` still equals `8ce83cde…` is
consistent with what I measure here; the integration baseline's protocol text is identical.

**Driver gap.** Round 3 was launched by the owner-authorized documented manual fallback; the driver
still halts on the unavailable historical worktree `…/scratchpad/f2repo` (organizer-notes Phase 1→2
and Phase 2, `inbox/claude-to-user_…_driver-error.md`). I altered no accounting, worktree history,
budget state or code to route around it, and I did not touch the unrelated historical idea whose
facilitator declaration blocks preflight.

---

## User direction

The owner's release direction is now controlling and is **not** an open question. Its English
direction, as recorded in `inbox/user-to-codex-1_meta-protocol-change-designated-implementer_release-order.md`
and copied to `source-context/release-order.md`:

> "Both — release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, label Windows as
> experimental, and hold winget for the CLI. At the same time open a separate reviewed Windows idea
> that fixes the defects and removes the label."

Release order: `release-1.49.1` → **this idea** → `windows-portability`, each gated on its
predecessor's done file; this idea waits for `codex-1-to-user_release-1.49.1_done.md` in the
lean-organizer worktree deck inbox. Every CLI release until (3) labels Windows
experimental/unvalidated, **keeps** the labelled Windows assets, opens **no** CLI winget PR; the
skill channels are unaffected. The original-language (Slovak) source is preserved verbatim in
`source-context/release-order.md` alongside its translation; I quote only the English direction here,
per the round instruction.

Consequences I treat as settled, not as design inputs: my round-02 "not resolvable here" item on the
inherited Windows scope decision is **withdrawn as resolved by the owner**. No Windows platform
architecture work is in this idea's scope; the separate run owns the fixes. This idea's release notes
carry the label and no winget PR, and the done-report makes no catalog-availability claim.

---

## Corrections to my own round-02

**C-4 — my C-2 argument rested on a false statement, and I retract that statement.** I wrote that
materialization "would make the driver a writer of `00-prompt.md` — an artifact the protocol assigns
to the idea's author, **which no tool writes today**." That is wrong. `updateIdeaStatus`
(`internal/consensus/consensus.go:898-928`) reads an existing idea's `00-prompt.md`, rewrites or
inserts its `status:` frontmatter line, and writes it back atomically; it is called at three phase
transitions (`consensus.go:194` → `consensus`, `:360` → `final`, `:407` → latest round). Tooling
already maintains a machine-owned frontmatter key inside the author's artifact. My conclusion against
materialization survives, but on different and stronger evidence (D-1 below); the premise I used in
round 2 does not, and zcode-1 and kimi-1 were both entitled to rely on it. This is the third round in
a row where a locator or a "no code does X" claim of mine needed correcting by measurement — the
pattern itself is the finding, and it is why R-F stands.

**C-5 — my round-02 `$PARLEY_AGENT_CONFIG` was a misnamed constant.** It is
`PARLEY_HEADLESS_AGENT_CONFIG` (`internal/config/runtime.go:17`). Immaterial to the argument, but it
is a fourth config layer and it matters to D-1; see N-7.

---

## New measurements this round

All PRIMARY, at HEAD `73ee923` in this worktree. These are the facts the seven decisions turn on.

**N-1 — tooling writes an existing idea's `00-prompt.md`.** `updateIdeaStatus`
(`consensus.go:898-928`), three callers. See C-4.

**N-2 — the only ordinary-idea kickoff writer is `CreateIdeaFull`
(`internal/protocol/workspace.go:178-232`), and it runs only for `parley run TASK`.** It builds the
whole `00-prompt.md` (frontmatter `idea/author/created/participants`, optional `track:`, §9.0
`excluded:` lines, `status: round-01`) and writes it once. Outside pipelines (§12.5,
`internal/pipeline/executor.go:297-362`) and `retro propose` (`internal/app/retro.go:123`), **no code
creates an idea**. Hand-authored ideas — §4.0's documented Phase-0 path, the path the owner brief
prescribed for this very idea, and the path every meta-protocol idea on this deck has taken — never
reach `CreateIdeaFull`.

**N-3 — preflight runs BEFORE the idea exists.** In `parley run`, `runTaskPreflight` is called at
`internal/app/app.go:1921-1922`; `runcontrol.Create` (which reaches `CreateIdeaFull`) is at `:1939`.
A per-idea designation gate located in preflight therefore **cannot see** a designation for a freshly
created idea, and for existing ideas it would have to scan the deck — the exact `facilitatorConflictGates`
shape that blocked this run (my R-D). The reliable choke point is `newDriverImplOps`
(`internal/app/driver_impl.go:45-98`), which is per-idea by construction, already computes `eligible`
and the implementer, and already carries the "escalated, not fallen back" refusal string `roleErr`
that four role actions check (`:214`, `:278`, `:407`, `:504`).

**N-4 — a materialized designation cannot carry machine-readable provenance in frontmatter.**
`CreateIdeaFull` writes roster-preset provenance as an HTML comment **below** the fence, with the
reason in code: *"Provenance is an HTML comment BELOW the frontmatter fence (review fix: inside the
fence, `ReadFrontmatter`'s `key: value` split would ingest it as a junk key)"* (`workspace.go:197-198`).
`ReadFrontmatter` (`workspace.go:368-395`) does `key, value, ok := strings.Cut(line, ":")` then
`meta[key] = strings.TrimSpace(value)` with **no comment stripping**, so an inline
`implementer: kimi-1  # from default_implementer` yields the literal value
`kimi-1  # from default_implementer`, which then fails participant validation — a legitimate
designation turned into a fail-closed gate. A prior review already caught this class once.

**N-5 — the review-exclusion consumer only ever runs after the work exists, and the census bounds the
exposure.** `expectedRoundParticipants` (`consensus.go:730-733`) returns `participants` unchanged
unless `review` is true; it calls the resolver **only** for review rounds. Census on this deck:
**75** `IMPLEMENTATION.md` files, **63** carry a non-empty `implementer:` pin, **12** do not; of the
75, exactly **9** have neither a readable pin nor a readable `FINAL.md` key
(`antigravity-agent-migration`, `consensus-request-signoffs`, `consensus-workflow-cli`,
`continuous-run-tui`, `repo-map-mvp`, `roadmap-implementation-plan`, `roster-operations-standard`,
`tui-action-execution`, `tui-workspace-sessions`). Those 9 are the entire set on which a *live* global
default could alter a review round's expected-artifact set. This measurement drives D-1(c) and it is
new; no round-02 artifact has it.

**N-6 — `--no-implement` is visible in the same flag scope as preflight.** `parley run`:
`no-implement` at `app.go:1807`, `no-preflight` at `:1810`, preflight at `:1921`. `parley continue`:
`:1159`, consumed at `:1253`. My round-02 open question 3 is answered: an availability gate can be
skipped for runs that will not reach Phase 5.

**N-7 — the `[defaults]` block has four layers, not two.** `configLayers` (`runtime.go:384-400`):
central `~/.parley/agents.toml` → deck `parley-deck/agents.toml` → `parley-deck/agents.local.toml` →
`$PARLEY_HEADLESS_AGENT_CONFIG` (non-optional when set, highest). `LoadDefaults` (`:419-436`) merges
across all of them, later non-empty wins. Any key added to that block inherits all four.

**N-8 — the protocol's two templates sit on opposite sides of the resolver.** The IMPLEMENTATION.md
template already prescribes `implementer: <agent-id>` (`COOPERATION.md:455`), which is why the
rank-1 pin is populated in 63 of 75 real artifacts; the FINAL.md template prescribes `author:
<agent-id>` (`:409-415`), which the resolver does not read. This is the mechanical root of kimi-1's
A2 and my F-2/F-3, stated in one line.

**N-9 — `centralDefaultTemplate` emits active keys with inline comments** (`runtime.go:627-648`:
`speed = "fast"  # …`, `ping_tier = "hosted-pong"  # …`). zcode-1's C-b correction to my round-01
claim is confirmed; a commented-out `default_implementer` line is a new emission shape.

**N-10 — nothing asserts the dispatch stdout line.** `grep -rn 'implementing via' --include="*.go"`
returns only the production site (`driver_impl.go:217`). Event consumers switch on known types with
no exhaustive-error (`internal/tui/protocolui.go:220-256`, `internal/runstate/runstate.go:352-451`),
and `agent.model_diversity` (`driver_impl.go:196-206`) is already an event no consumer switches on.
A new event type is tolerated by construction.

**N-11 — `driver_impl.go` already refuses on Windows.** `completeWithWriter` returns *"independent
evidence completion requires a POSIX execution host; Windows runtime is not supported"*
(`driver_impl.go:543`). Pre-existing, in a file this delta touches, and consistent with the owner's
experimental labelling. This delta adds no build-tagged file and no new platform surface; my round-02
Windows assessment otherwise stands unchanged.

---

## The seven decisions

### D-1 — Live layered read; **no** Phase-0 materialization

**Decision: reject materialization. Tier 3 is a live layered read whose resolved id and source are
recorded in the durable event stream, not written into `00-prompt.md`.** Against kimi-1's position
change 2 and zcode-1's position change 2 (both adopted my withdrawn round-01 D3, which C-2 retracted
and which I now reject on new evidence rather than on the premise C-4 corrects).

(a) **Materialization has no hook that covers the ideas this deck actually runs (N-2/N-3).** The only
ordinary-idea writer is `CreateIdeaFull`, reached only by `parley run TASK`. A hand-authored idea —
including this one — never passes through it, and preflight, the other candidate site, executes
before the idea exists. So a materializing implementation would apply the owner's standing
instruction to quick `parley run` ideas and **silently skip** hand-authored ones. Both peers wrote
"kickoff writes the id"; for hand-authored ideas there is no kickoff code to write it. The remaining
option — materialize on first driver invocation against an existing idea — is no longer "Phase 0": it
mutates an author-owned artifact at an arbitrary later time and must then special-case ideas already
past Phase 0, dry-runs, and `--no-implement`. A live read applies uniformly regardless of how the
idea was created.

(b) **Materialization and kimi-1's two-tier availability policy are mutually exclusive (N-4).** kimi-1
wants the per-idea field to fail closed and the global layer to fall through loudly. After
materialization there is no global layer at resolution time: the materialized line is a per-idea field,
byte-indistinguishable from an owner-written one, because provenance must live below the fence or be
ingested as a junk key. To keep the two-tier policy, materialization needs a second machine-read key
(`implementer_source:`) — more surface, and a key whose forgery silently downgrades an owner's
designation to a fall-through. The live read keeps the two layers distinguishable for free, which is
the only way the policy kimi-1 and I both want is implementable at all.

(c) **The freeze materialization buys is already bought, and where it is not, an event buys it more
cheaply and also covers the unset path (N-5/N-8).** Rank 1 is the `IMPLEMENTATION.md` pin, prescribed
by the protocol's own template and present in 63 of 75 real artifacts. The only unpinned window is
Phase-5 dispatch → the implementer's first `IMPLEMENTATION.md` write. For that window: the driver
records `{implementer, source}` at dispatch, and on re-entry compares the live resolution to the
recorded one, escalating on a mismatch. **The comparison fires only when the recorded source was
`designation` or `global-default`** — on the unset path the event is written and nothing compares it,
so a changed `--participants` order still silently re-resolves exactly as today.

**Non-rostered global designee:** tier 3 that names an agent absent from `participants:` is
**inapplicable, not invalid** — a one-line notice and fall-through to tier 4. This is zcode-1's
original instinct and both peers' "record inapplicability" rule; the existing `isParticipant`
validation inside both resolvers (`driver_impl.go:105-112`, `consensus.go:752-759`) already supplies
the mechanism.

**Unavailable global designee:** notice and fall through (kimi-1's argument, accepted: a standing
preference that hard-stops every idea whenever one agent is down is unusable). **Unavailable per-idea
designee: blocking gate with the pre-built `Confirm`** — my S-5, unchanged, on the LE-7/LE-11 and
`COOPERATION.md:699` precedent and on the fact that the fall-through destination is the list-order
selection the owner is replacing. Both peers adopted fail-closed here in round 2; the split is now
by layer, and per (b) that split only exists without materialization.

**Resumability, stated as a rule:** the pin governs re-entry (rank 1); before a pin exists, the
recorded dispatch event governs; a live tier-2/tier-3 change against a recorded dispatch escalates
rather than reassigning silently. Nothing in the resolution chain depends on a config file's value at
a past moment, so a resumed run cannot be surprised by a config edit it cannot see.

### D-2 — Per-idea empty/`none`: four states, and `none` is the opt-out

**Decision: absent → no designation (today's chain). `implementer: none` → explicit per-idea
opt-out, suppresses tier 3, no gate. `implementer:` bare-empty → incomplete designation, blocking
gate naming both exits. Non-empty non-participant (or the declared non-participating facilitator) →
blocking gate.**

This adopts zcode-1's three-state insight, which is better than my S-3 and than kimi-1's
`facilitator:`-symmetry, and it fixes the one thing zcode-1's version cannot do. zcode-1 justified
the opt-out state as "the per-idea opt-out against one's own global default" — correct, and it
survives without materialization, which is why I adopt it despite rejecting the mechanism it was
proposed for. But collapsing bare-empty and `none` into one meaning makes a typo indistinguishable
from an intention: an owner who typed `implementer:` meaning to name someone gets silence, which is
the failure mode this idea exists to remove. Requiring the explicit token separates them
deterministically, at the cost of one documented word.

Mechanically: `raw, present := meta[ImplementerKey]`; `!present` → absent; `strings.TrimSpace(raw) ==
""` → gate; `== "none"` → opt-out; else validate against `participants:`. `ReadFrontmatter`
(`workspace.go:390-394`) creates the map entry for a bare `implementer:` line, so presence and
emptiness are distinguishable in one expression — re-verified this round, and it must be test-pinned
or the opt-out silently collapses into "absent" (zcode-1's concern #2, adopted).

Divergence from `facilitator:` (`facilitator.go:44-48`, empty → undeclared) is principled and must be
in the code comment, not only here: `facilitator:` is exclusionary, so an empty one excludes nobody;
`implementer:` is appointive, so an empty one appoints nobody and the silent destination is
`eligible[0]`.

### D-3 — The write-end FINAL repair stays OUT

**Decision: out of this delta. Deferred with the read-end widening into one named follow-up.** I keep
my C-3 withdrawal and side with kimi-1 against zcode-1's proposal #6, with an argument neither of us
made in round 2.

zcode-1 calls the cost "the one edge — drafter ≠ `eligible[0]` when `participants[0]` is not
headless/found". That understates it. The FINAL drafter is also chosen by the protocol's **documented
volunteer mechanism** (`COOPERATION.md:409`, `Drafter: yes`), and in the immediately preceding run the
volunteer was not `participants[0]` — kimi-1 drafted while the driver announced "drafting FINAL via
claude-1" (lean-organizer `organizer-notes.md:41`, kimi-1's V12, zcode-1's confirmation). So the
write-end repair would change Phase-5 dispatch on every future undesignated idea where a
non-first participant drafts. That is not an edge; it is the normal volunteer path.

And the direction of the change is the decisive point: the drafter-inherits rule **is** the
first-mover chain. The prose route to it is "first agent to submit a round-01 file drafts" (`:407`)
and "the default implementer is the FINAL drafter" (`:443`). Writing `implementer: <drafter>` into
new FINALs would take the selection the owner asked to stop relying on and make it *mechanically
authoritative* on decks that never opted into the new mechanism. Whatever else is true, that cannot be
what "behaviour without it is exactly today's" was protecting.

**One deferred follow-up, not three.** The read-end widening (my C-3 / zcode-1's position change 3),
the write-end emission (zcode-1 #6), a machine reader for claims, and the cheapest real repair I
found this round — having the driver's `Complete` (`driver_impl.go:518-560`, which today only
transitions status) also record the `implementer:` it already knows, closing the 12-of-75 missing
pins in N-5 — all change unset-path selection and all belong in one idea the owner can see whole.
Suggested slug: kimi-1's `meta-protocol-change-implementer-inheritance-repair`. Splitting them across
separate follow-ups invites half a repair.

### D-4 — Claims: one prose rule, one machine ladder, never one blended ladder

**Decision: the claim is a social act with no machine reader, and FINAL must state the two ladders
separately and label them.** All three of us withdrew claim parsing — and all three round-02
restatements then printed a single precedence ladder containing a claim rank. kimi-1: "`IMPLEMENTATION.md`
pin > per-idea field > **perfected uncontested claim (only when no field)** > FINAL-drafter
inheritance > positional tail". zcode-1 #3: "… > **perfected inbox claim (only when no designation
exists)** > FINAL metadata > positional". Mine (S-4) omits the claim but says so only in a
parenthetical. A ladder that mixes a human act with machine ranks is either a specification for a
parser none of us wants, or an ambiguity a FINAL implementer must guess at. My B-1 block is satisfied
as to parsing; this is the residue, and it is a drafting requirement, not a disagreement about intent.

FINAL states:

- **Machine resolution (dispatch), four ranks, no claim:** `IMPLEMENTATION.md implementer:` → the
  per-idea designation → `[defaults].default_implementer` → today's tail (`FINAL.md`
  `implementer`/`drafted-by`, then `eligible[0]` in the driver).
- **Machine resolution (review exclusion), unchanged:** `IMPLEMENTATION.md implementer:` →
  `FINAL.md implementer`/`drafted-by` → `""` → full list. See D-5.
- **Protocol prose:** the default implementer is the FINAL drafter; a participant may claim by
  `inbox/<from>-to-all_<slug>_impl-claim.md` before work begins; **under a live designation a claim
  does not override it** — it is recorded, surfaced as dissent, and answered by a one-line per-idea
  edit or a §4 escalation. FINAL must say explicitly that no code reads the claim file today
  (`grep -rn "impl-claim" --include="*.go"` → empty) and name the follow-up that would give it one.

Advisory surfacing: I accept zcode-1's narrowing of my round-01 position — no gate when a claim exists
without a designation, because that is today's normal world. Surface a contradiction only when a
designation is present.

### D-5 — Designation is a dispatch input only; the review-exclusion read set does not change

**Decision (new this round, and I ask both peers to confirm it): tiers 2 and 3 enter only the
driver's dispatch selection. The consensus-side read that computes review-round expectations keeps
exactly today's two artifact sources.**

Rationale: review exclusion must name **who actually implemented** — a recorded fact — never **who was
instructed to**. A preference cannot retroactively define who did the work. N-5 makes the cost of
getting this wrong concrete: on the 9 ideas with neither a pin nor a readable FINAL key, a live tier-3
read inside `expectedRoundParticipants` would silently shrink a re-opened review round's expected set
for an idea that never opted into anything. With tiers 2–3 confined to dispatch, that exposure is
**zero** and `parley wait`'s expected-round computation is untouched by this feature in every state.

Three consequences worth stating plainly:

1. It is a stronger boundary guarantee than any round-02 proposal offered, including mine.
2. It changes the S-9 unification shape: the shared reader is the **artifact** reader (the two files,
   exactly as today) taking an ordered source list; the driver passes
   `[Implementation, Designation, GlobalDefault, Final]`, consensus passes `[Implementation, Final]`.
   The divergence between the two callers becomes explicit, parameterised and testable instead of
   accidental — which is what X-5 was really asking for.
3. It takes **X-5 off this delta's critical path**. zcode-1's find stands (the driver validates
   against facilitator-filtered `eligible`, consensus against raw `participants`), and it remains a
   real pre-existing divergence, but unification no longer has to resolve it to be correct. It joins
   the D-3 follow-up register. zcode-1 owns that claim; I confirmed it in round 2 and I am not
   re-verdicting it, only re-scoping what this idea must fix.

### D-6 — Observability: event always-on, stdout byte-identical on the unset path

**Decision: `agent.implementer_resolved {idea, implementer, source}` is emitted on every dispatch,
including the unset path. The existing `driver: implementing via %s ...` line
(`driver_impl.go:217`) stays byte-identical; a second stdout line prints only when tier 2 or tier 3
fired.** `source ∈ {implementation-md, designation, global-default, final-md, positional}`. Modelled
on `agent.model_diversity` (`driver_impl.go:196-206`) including its best-effort store guard.

Why always-on, as an argument rather than a preference: D-3 defers a real defect — the documented
drafter-implements chain does not run, and dispatch falls to list order. The honest mitigation for
deferring it is that the run now **records** `source: positional` where today it records nothing.
Make the event designation-only and the deferral loses its mitigation: the unset-path divergence stays
exactly as invisible as it has been, and the next reader re-derives it from three files as we did.
The event is also what D-1(c) compares against on re-entry, so designation-only would leave the unset
path without the record it needs anyway.

Why this is inside the boundary: no selection, no gate, no exit code changes; N-10 shows no test
asserts the dispatch line and that unknown event types are tolerated by every consumer. It is still
the one element that adds a record where none existed, so it lands in `## Agreed trade-offs` with this
reasoning. **kimi-1 and zcode-1 both listed the resolved-source event as additive and fine in round 2
without distinguishing always-on from designation-only; I am asking for the always-on variant
explicitly so the decision is recorded rather than inherited.** If either peer wants the conservative
variant I will sign it, and then D-3's deferral should be re-argued, because the mitigation is gone.

### D-7 — Re-entry pin conflicts: pin wins, never silently, and never by self-appointment

**Decision: rank 1 wins whenever tiers 1 and 2 agree or tier 2 is absent (today's behaviour). When
both are present and disagree, escalate. The pin may be rewritten only by a recorded,
owner-confirmed line — not by the incoming implementer's own hand.**

I accept zcode-1's and kimi-1's "pin wins" for execution continuity. I reject one half of zcode-1's
correction path: *"an explicit, logged frontmatter edit by the incoming implementer (the pin is the
record, and the record's owner during execution is the implementer)"*. An agent editing
`IMPLEMENTATION.md`'s `implementer:` to name itself is self-appointment with a paper trail, which is
the act this mechanism exists to remove; that it is logged does not make it authorized. kimi-1 is
closer — "never by editing the field and pretending" — but "a deviation-style log entry" does not say
whose confirmation makes the reassignment valid.

The authorization shape already exists and is code, not prose: §9.0 records per-idea exclusions as
`excluded: [<roster-id> — reason — confirmed <date>]` in `00-prompt.md` (`COOPERATION.md:872-877`),
and `CreateIdeaFull` writes exactly those lines (`workspace.go:186-192`). So: a changed designation
over a live pin escalates with a `Confirm` offering (a) restore the designation to match the pin, or
(b) record `implementer_reassigned: <old> → <new> — <reason> — confirmed <date>`, which retires the
abandoned attempt and re-pins. The `confirmed <date>` element is the owner's, exactly as in §9.0. A
participant may draft the line; it is not valid until confirmed. Designation-path only: it requires
tier 2 to be present, so an undesignated idea behaves as today.

---

## Designation-only classification of the design-preference changes

Explicitly, per the round instruction, for the two items that are preferences rather than gates —
both **retained and designation-only**, inert with no designation present:

- **Drafter-separation preference (S-10).** `firstEligibleHeadlessAgent`
  (`driver_consensus.go:109-123`) prefers an eligible drafter that is not the designated implementer,
  **falling back to the designee when no other eligible drafter exists**, using the filter slot the
  function already has. Fires only under a designation; the degenerate fallback is test-pinned
  (zcode-1's condition, adopted); drafter == implementer is never forbidden (kimi-1's O4 /
  two-participant case), and where they coincide `FINAL.md` records the concentration in one line.
- **Phase-0 model-diversity relocation (S-11).** The existing `checkModelDiversity`
  (`driver_impl.go:179-210`) runs unchanged at Phase 6; when a designation is present it *also* runs
  at kickoff so an unsatisfiable roster costs an error message rather than a design cycle. No new
  gate class, no severity change. Two-participant designated ideas get a kickoff **warning**, not a
  block (zcode-1's O4 reading, on `COOPERATION.md:746` — his corrected locator, which supersedes my
  round-01 L466-468 — and §4.0's reviewer degradation); zero invokable non-implementer reviewers
  stays the hard stop it already is (`driver_impl.go:282`).
- **Design-only / `--no-implement` runs.** Validity gates (bare-empty, non-participant, declared
  facilitator) fire on any run — a defective frontmatter line is worth stopping on. The *availability*
  gate attaches only to runs that can reach Phase 5, which N-6 shows is decidable in-process.
  My round-02 open question 3 is closed.

---

## Response to kimi-1

**Adopted and credited:** the withdrawal of `designated_implementer:` (D-2 keeps `implementer:`); the
withdrawal of dimension-4 pre-dispatch validation as redundant, with the L433 self-containment
sentence as the surviving delta; your deferral of the D10 read-end widening to a named slug, which
D-3 extends to the write end; the global-layer loud fall-through, which D-1 adopts *because* it
rejects materialization; `parley init` emitting the key commented-out as the physical form of the
boundary — confirmed as a deliberate deviation from the emission shape at `runtime.go:627-648`
(N-9, zcode-1's C-b); your closure of the dormant-path test requirement; and your §15 adjudications,
which I take up below.

**K-8 — your position change 2 (materialization) is the one thing I now block on, and the reason is
your own two-tier availability policy.** You adopted materialization to dissolve O1 *and* proposed
per-idea fail-closed with global loud fall-through. N-4 shows those cannot both ship: after
materialization the two layers are the same line, and the provenance that would distinguish them
cannot live in frontmatter without a second machine-read key. N-2/N-3 add that materialization would
not fire at all for hand-authored ideas — the class this deck runs. Your ordering conclusion is right
and survives intact under a live read: an owner's standing instruction outranks a participant's
claim because the claim has no machine reader and the designation does. You do not need
materialization to get that; you need D-4's two labelled ladders.

**K-9 — your "ship the two required levels only" (concern 2) is still not available, and it now has a
third layer to account for.** `LoadDefaults` merges `[defaults]` across `configLayers`, later
non-empty wins (`runtime.go:384-400`, `:419-436`). Placing `default_implementer` there yields
per-idea > env > deck-local > deck > machine **by construction** (N-7). Reading machine-only requires
a deliberate special case that would make this the one `[defaults]` key ignoring the deck. Accept the
inherited layering, state it in the protocol text, and pin it with a test that sets the key in two
layers and asserts the higher one wins. You and zcode-1 both wrote "two levels only"; the code does
not offer that choice.

**K-10 — your materialization mechanics (idempotent, Phase-0-only, never overwrite an owner field,
never touch an in-flight idea) are exactly right, and they are the specification of a component D-1
deletes.** I record them because if the group overrules D-1, those four properties plus N-4's
below-the-fence constraint are the minimum contract, and the `implementer_source:` key becomes
mandatory rather than optional.

---

## Response to zcode-1

**Adopted and credited:** your three-state per-idea semantics, which D-2 refines rather than replaces
— the opt-out state is a real requirement all three round-01 proposals missed, and it is yours; your
C-a locator correction (`COOPERATION.md:746`, superseding my L466-468) and C-b emission correction
(N-9), both of which I accept without reservation; your position change 1 (fail-closed unavailable
designee), which converges with mine; your retraction of the `consensus.go:805` premise, which
strengthens rather than weakens your own §8; the idea-scoped gate family as a design requirement; the
skill-tree grep, which you completed by inspection (`SKILL.md` restates no Phase-5 rule — only the
facilitator exclusion and advisory roles), closing my R-C's second half; and your confirmation of
F-1…F-9 with your own locators.

**Z-1 — proposal #6, the write-end repair: I hold my withdrawal and I think your cost estimate is
wrong by a category.** D-3 gives the argument: the volunteer-drafter path (`COOPERATION.md:409`,
`Drafter: yes`) is not an edge, it is the protocol's own mechanism, and it fired in the run
immediately before this one. So the repair changes dispatch on ordinary undesignated ideas, in the
direction of making the first-mover chain mechanically authoritative. You framed this as "the single
recorded default-path trade-off"; I think it is the one item in your consolidated proposal that the
owner's sentence actually forbids, and I would rather ship the mechanism complete — which all three
of us agree it is without the repair — than spend the boundary on it. **This is my only substantive
block, and it is narrow: delete proposal #6, keep everything else in your #1–#11.**

**Z-2 — your "preflight gate" location cannot carry the authority you give it (N-3).** Preflight runs
before the idea exists in `parley run`, and for existing ideas it is the deck-wide scan whose
all-ideas iteration you correctly made normative to avoid. Put the authoritative fail-closed check in
`newDriverImplOps`'s `roleErr` path — per-idea by construction, already the "escalated, not fallen
back" site you quoted (`driver_impl.go:47-51`), and checked by four role actions — and let preflight
carry an early idea-scoped copy for ergonomics only. Your requirement survives; its address changes.

**Z-3 — your correction path for a changed designation authorizes the wrong party.** D-7: the
incoming implementer editing the pin to name itself is self-appointment, logged. The §9.0
`confirmed <date>` shape (`COOPERATION.md:872-877`, written by `workspace.go:186-192`) is the
authorization mechanism that already exists; use it.

**Z-4 — your proposal #8 pre-dispatch FINAL validation, scoped to designated runs: agreed on scope,
but note it is already shipped and needs no code.** `protocol.ValidateFinal`
(`internal/protocol/finalsections.go:92`) with `RequiredFinalSections` including
`## Observable acceptance criteria`, wired pre-Phase-5 at `internal/driver/consensus.go:57`/`:61` —
my K-1 to kimi-1, which kimi-1 accepted and withdrew for. Scoping a *new* validation to designated
runs would add a second gate over the same artifact. Cite the shipped one; add nothing.

**Z-5 — your proposal #8 trajectory-agreement gate: I still argue out of scope, and now more
narrowly.** `parley trajectory configure` is attended, requires `--yes` and a non-empty named-checks
contract (`internal/app/trajectory.go:65`, `:105-125`), and freezes an evidence-policy patch author,
not a launch target. Your scoping (fires only when both a policy and a designation exist) removes the
default-path objection, so this is no longer a boundary question — it is a scope question, and the
owner's request does not reach it. Follow-up register, with your scoping recorded so it need not be
re-derived. I will not block consensus if you keep it, but it buys a coupling to the trajectory
subsystem for a mismatch no one has observed.

**Z-6 — I ask you to confirm or refute D-5**, since it re-scopes your X-5 find and changes the
unification shape you and I converged on. I own no verdict on X-5; D-5 does not contest it.

---

## §15 verdict-conflict resolutions

**V-A — "1 CRITICAL and 8 MAJOR findings in the implementation, all fixed in later cycles" (owner
brief testimony). My round-01 CONFIRMED vs zcode-1's UNVERIFIED for the "all fixed" clause, which
kimi-1 adjudicated in zcode-1's favour as "supported-not-traced".** kimi-1's adjudication was correct
**on the evidence then filed**: my round-01 evidence was finding headings, which do not entail fixes.
The organizer asked for the trace, so I performed it rather than re-arguing it. **Resolution:
CONFIRMED, PRIMARY, with the per-finding trace now supplied; zcode-1's UNVERIFIED is superseded by
measurement, not overruled by argument.** Evidence, all in the lean-organizer worktree deck:

- Severity counts: `review/round-01/claude-1.md:670` — "CRITICAL 1 · MAJOR 8 · MINOR 7 · NIT 3 (19
  findings)": CRIT-1, MAJ-1…MAJ-8.
- Finding → fix mapping, cycle 1: `review/consensus-cycle-01.md:53-60` maps CRIT-1→F1,
  MAJ-1+MAJ-2→F2 (merged; ≡ kimi-1 K1-F4), MAJ-3→F3, MAJ-4→F4, MAJ-5→F5, MAJ-6→F6, MAJ-7→F7,
  MAJ-8→F8, with F1/F2/F8 marked contingent on the then-open verdict conflicts VC-1/VC-2/VC-3.
- Application, per fix, in the implementer's own record: `IMPLEMENTATION.md` "Fix-up cycle 1"
  (`:676`) — F1 at `:770` (removal of the unratified gate), F2 at `:781` (`internal/app/wait.go`),
  F3 `:796`, F4 `:801`, F5 `:805`, F6 `:811`, F7 `:818`, F8 `:825`, and on through F9–F21 for the
  MINOR/NIT set.
- Closure: `review/consensus.md` is review-cycle 4 with `outstanding_agreed_fixes: 0`;
  `IMPLEMENTATION.md` frontmatter `status: complete`, `fix-up-cycle: 3`.

Two precisions the unqualified phrase hides, and which must travel with it if it is cited: **nine
findings produced eight fixes** (MAJ-1 and MAJ-2 merged into F2), and **not every disposition was a
repair** — F1 *removed* the disputed change rather than fixing it, and F15 dropped an export. Also,
"applied" is recorded first-person by the implementer; what makes it admissible is that cycles 2–4
are reviewer artifacts over that fix-up and close at zero outstanding. Materiality is unchanged: this
feeds only the post-release recommendation evidence base, never an acceptance criterion — and it is
material to the role assignment below, which is why I traced it rather than leaving it open.

**V-B — "implementation is the heaviest token work of a run".** Unchanged and settled identically by
all three: the copied study supports ~30-40% of an *implementing organizer's* input with its own
caveats; nothing supports the ranking. Stays UNVERIFIED and must not enter FINAL as established
(§15.2). My round-02 items 4 and 5 (the evidence asymmetry — I read only the in-deck copy, zcode-1
the external original) stand; cite zcode-1's reading for the original's numbers.

**V-C — the brief's "speed decides the implementer" hypothesis.** Settled identically by all three in
round 2 via scope decomposition, no residual factual conflict: CONFIRMED of the protocol prose
(`:407`→`:409`→`:443`), WRONG of the shipped tooling (dispatch is list order; no code computes a
first filer). FINAL states it scoped, never bare. D-3 now depends on this resolution, so it is no
longer merely descriptive.

**V-D — locator corrections, closed.** kimi-1's line numbers were right and mine and zcode-1's were
wrong (`driver_impl.go:131-133`/`:134`; `consensus.go:778`) — X-1/X-2, uncontested since. zcode-1's
C-a and C-b correct two more of mine (N-9, `COOPERATION.md:746`). My C-4 and C-5 above correct two
more. No conflict remains; the pattern is the finding, and R-F is the mitigation.

---

## Default-path ledger (updated)

On a deck that sets neither `implementer:` nor `default_implementer`:

**Byte-identical:** the dispatch chain collapses to today's two-file walk plus `eligible[0]`;
**the review-exclusion read set is untouched in every state** (D-5) — `parley wait`'s expected-round
computation cannot be perturbed by this feature at all; both terminal tails and their pinned tests
(`internal/app/app_test.go:1557`, `internal/consensus/roundgate_test.go:102`);
`checkModelDiversity` timing and severity; `firstEligibleHeadlessAgent` ordering; `ValidateFinal` and
`finalScaffoldReason`; the `impl-claim` file's inertness; `roles:` semantics; quorum, roster, signoff
weight; the `driver: implementing via %s` stdout line (D-6); every attended gate — `parley protocol
publish` stays TTY-attended, core publication owner-attended, channel publication the organizer's
attended step, and a designated implementer inherits none of them (S-12, and one protocol sentence
saying so, because "does whatever work is needed" is stretchable).

**Designation-path only (tier 2 or 3 set):** bare-empty and invalid-target gates; the unavailable
per-idea designee gate; `none` opt-out handling; pin-conflict escalation and the
`implementer_reassigned:` confirm; the drafter-separation preference; the kickoff diversity check and
two-participant warning; the re-entry dispatch comparison; the second stdout line.

**The one unset-path addition, declared:** the `agent.implementer_resolved` event (D-6). No selection,
no gate. `## Agreed trade-offs` entry required.

**Deferred, evidence preserved, one slug:** FINAL read-set widening; FINAL write-end emission;
`Complete` recording the pin it already knows (N-5's 12 missing pins); a machine reader for
`impl-claim`; `facilitatorConflictGates`' all-ideas iteration; X-5's eligibility divergence (D-5);
the trajectory agreement check (Z-5). **Constraint on this delta even though its repair is deferred:**
the new gate must be scoped to the idea being run, or this idea ships the trap it tripped over.

---

## Existing alternatives

Per §15.6(a), only for elements whose shape **changed** this round. Round-01 and round-02 tables cover
the rest and are not restated.

| Element changed in round 3 | What already ships | Locator | Why this shape |
|---|---|---|---|
| Global default resolution (D-1: live read, materialization rejected) | `updateIdeaStatus` is the only tool that mutates an existing `00-prompt.md`, and only `status:`; `CreateIdeaFull` is the only ordinary-idea creator and runs only for `parley run TASK`; provenance must sit below the fence | `consensus.go:898-928`; `workspace.go:178-232`, `:197-198` | Materialization has no hook for hand-authored ideas (N-2/N-3) and cannot carry distinguishable provenance (N-4), which makes the two-tier availability policy unimplementable. |
| Designation scope (D-5: dispatch only) | `expectedRoundParticipants` calls the resolver only for review rounds and fails closed to the full list | `consensus.go:730-733`, `:751-779` | Review exclusion must read the record of who implemented, not the instruction. Zero exposure on the 9 ideas measured in N-5. |
| Validation address (Z-2) | `newDriverImplOps` computes `eligible` and the implementer and carries `roleErr`, which four role actions escalate on; "the driver NEVER silently falls back" | `driver_impl.go:45-98`, `:214`, `:278`, `:407`, `:504` | Per-idea by construction, covers run/continue/resume/TUI, and cannot repeat the cross-idea blocking that halted this run's preflight. |
| Empty-value semantics (D-2) | `ReadFrontmatter` distinguishes present-empty from absent; `FacilitatorRoleFromMeta` deliberately collapses them for an exclusionary field | `workspace.go:390-394`; `facilitator.go:44-48` | Four states separate a typo from an intention deterministically; the asymmetry with `facilitator:` is principled and belongs in the code comment. |
| Reassignment authorization (D-7) | §9.0 `excluded: [<id> — reason — confirmed <date>]`, written by code | `COOPERATION.md:872-877`; `workspace.go:186-192` | The owner-confirmation shape exists; an incoming implementer rewriting its own pin does not become authorization by being logged. |

Sources consulted this round: the full phase-2 packet body; `internal/protocol/`, `internal/app/`,
`internal/config/`, `internal/consensus/`, `internal/driver/`, `internal/pipeline/`, `internal/tui/`,
`internal/runstate/`; both peers' round-02 artifacts in full; `00-prompt.md`,
`source-context/owner-brief.md`, `source-context/release-order.md`, `organizer-notes.md`, the
release-order inbox note; the lean-organizer deck's `review/round-01/claude-1.md`,
`review/consensus-cycle-01…03.md`, `review/consensus.md`, `IMPLEMENTATION.md`; `parley --help`.

---

## Role assignment for this run

**I accept kimi-1's and my own proposal: claude-1 drafts FINAL, kimi-1 implements, claude-1 + zcode-1
review. I decline zcode-1's counter (zcode-1 implements), with reasons, and I would sign it without
further argument if kimi-1 declines or is unavailable.**

Not a vote count. The reasons are about what this run produces for the owner's pending decision:

1. **Evidence yield.** The owner must choose a default implementer after release on evidence this deck
   does not yet have. With kimi-1 implementing, the run produces **two** new agent-role data points —
   kimi-1's first implementation and zcode-1's first code review on this deck. With zcode-1
   implementing, it produces **none**: zcode-1's second implementation and kimi-1's second review.
   At identical cost, one assignment informs the pending owner decision and the other does not.
2. **Anti-concentration.** zcode-1 implemented the immediately preceding protocol-core change
   (lean-organizer, `IMPLEMENTATION.md` `implementer: zcode-1`). Designating it again puts two
   consecutive protocol-core implementations in one model — the standing-concentration risk all three
   of us named independently (my R-E, kimi-1's risk 1-2, zcode-1's own risk 3). zcode-1 offered its
   lean-organizer precedent as capability evidence and correctly labelled it n=1 and not a default
   argument; I take it exactly that way, and V-A's trace confirms the work closed at zero outstanding
   fixes across three cycles. The objection is to repetition, not to competence.
3. **Locator fidelity is this delta's dominant risk**, across three protocol copies, six `implementer`
   readers, and — after three rounds of corrections from all three of us — a demonstrated failure mode.
   kimi-1 was the only one of us whose round-01 code locators were all correct (X-1, X-2). Small
   sample, and I say so; it is the only directly relevant evidence on the table.
4. **Why not me.** The deck's only implementation-review data point is my lean-organizer review (19
   findings, 1 CRITICAL / 8 MAJOR, traced in V-A). Spending that on implementation removes the deck's
   one demonstrated high-yield reviewer from the reviewer set. I would rather review this than
   implement it — and I am the drafter under both proposals, which D-3's S-10 separation makes the
   right place for me.

**Risk I record against my own recommendation:** kimi-1's round-01 and round-02 artifacts are the
shortest filed (108 and 118 lines against my 593/619 and zcode-1's 356/435). Length is not quality and
kimi-1's measurements have been the most accurate of the three. But this delta is wide — a new
`internal/protocol` reader, `driver_impl.go`, `driver_consensus.go`, `runtime.go`, three
`COOPERATION.md` copies, the changelog — and a terse implementation of a wide surgical change is where
drift enters. The mitigation is mechanical, not social: FINAL must be self-contained for a
non-drafter (`COOPERATION.md:433`, the rule designation makes load-bearing), and the acceptance
criteria must be checkable without judgement — byte-comparison of the changed region across all three
protocol copies, the two existing tail tests unchanged, and named tests for the four-state parsing,
the layer precedence, the degenerate drafter fallback, and each dormant tier-3 state.

Both claims are perfected at the protocol-prescribed times, not here: my `Drafter: yes` in the
consensus signoff block (`COOPERATION.md:409`), kimi-1's
`inbox/kimi-1-to-all_meta-protocol-change-designated-implementer_impl-claim.md` before Phase-5 work
begins (`:443`). **Nothing above names the owner's future global default.** This is a this-run
assignment under today's rules; the global default ships UNSET and the recommendation is a
post-release act on evidence that does not yet exist. Conflating the two would be the quiet
self-selection this idea exists to remove.

---

## Risks

- **R-G (new) — D-5 arrived in round 3 and changes the unification shape both peers signed.** If
  FINAL is drafted from the round-02 restatements it will unify the reading *including* the new tiers
  and re-introduce the N-5 exposure. Whoever drafts must carry D-5 explicitly or the delta is wrong
  in a way that only shows up on a re-opened review round.
- **R-H (new) — the deferred follow-up is now four defects in one slug.** D-3 bundles the read end,
  the write end, `Complete`'s missing pin, and the claim reader. Bundling is right for owner
  visibility and wrong for the risk that a big follow-up never opens. D-6's always-on event is the
  only thing that keeps the largest of them visible in the meantime; that is the whole argument for
  always-on, and it should be stated in `## Agreed trade-offs` in those terms.
- **R-B/R-D/R-E/R-F unchanged from round 2**: the `facilitator:` present-empty asymmetry needs its
  reason in the code comment; a new gate class is a new way for `parley run` to stop and must be
  idea-scoped (non-negotiable for me); designation concentrates implementation in one model, mitigated
  by the sequencing the owner already chose; and no locator should be transcribed from any round
  artifact, mine included, without re-measurement.
- **R-C half-closed.** zcode-1's skill-tree grep closed the prose-restatement half by inspection. The
  three-copy byte comparison of the changed region remains an acceptance criterion.

---

## Consensus readiness

**Ready for consensus with one substantive block and one confirmation requested.**

**Block B-2 — zcode-1's proposal #6 (write-end FINAL `implementer:` emission) must come out.** D-3 /
Z-1: it changes dispatch on ordinary undesignated ideas via the protocol's own volunteer-drafter path,
in the direction of making the first-mover chain mechanically authoritative. I will not sign a FINAL
that carries it. Deferring it costs nothing — all three of us agree the mechanism is complete without
it — and D-6's always-on record keeps the underlying defect visible. This is narrow and mechanical;
one exchange should settle it.

**Block B-3 — materialization must come out (kimi-1 position change 2, zcode-1 position change 2).**
D-1 / K-8: no hook covers hand-authored ideas, and materialization makes the two-tier availability
policy both peers want unimplementable without a second machine-read key. I hold this as a block
rather than a preference because a FINAL specifying "kickoff writes the id" is not executable as
written — there is no kickoff code on the path this deck uses — and the round instruction requires an
executable deterministic mechanism.

**B-1 (round 2) is satisfied** as to claim parsing: both peers withdrew it. Its residue is D-4's
drafting requirement, not a block.

**Confirmation requested, not a block — D-5.** It is new, it strengthens the boundary beyond anything
proposed in round 2, and it re-scopes zcode-1's X-5. I ask both peers to confirm or refute it; if
either refutes it I will sign the round-02 unification shape with N-5's exposure recorded as an
explicit trade-off.

**Recorded decision D-6** (always-on event) is not a block; I will sign either variant, but the
conservative variant weakens D-3's deferral and that link must be recorded rather than inherited by
silence.

Everything else I consider settled: D-2 and D-7 synthesize all three round-02 positions and I expect
no argument; the release direction is the owner's and closed; the deferred register is decided *out*
with its evidence preserved. **With B-2 and B-3 resolved, cross-review 3 is unnecessary and I would
support the organizer opening consensus on the round-03 responses.**

---

## Process notes

**Evidence ownership (§15).** N-1 through N-11 are first-canonical claims of mine and carry no verdict
from me; kimi-1 and zcode-1 own those verdicts, and N-5's census and D-5's consequence are the two I
most want independently re-measured. C-4 and C-5 correct claims I own. V-A is a verdict on the owner
brief's testimony, supported by a trace I performed; zcode-1 owns the UNVERIFIED it supersedes and may
contest the trace. V-D records corrections others made to claims I own; I accept all of them. Where I
cite the copied token study or `source-context/preflight.json` I am transcribing testimony and have
marked it. My F-1…F-9 remain mine and unverdicted by me; both peers have now confirmed all nine
independently.

**Independence and scope.** I read the current `00-prompt.md`, `source-context/owner-brief.md`,
`source-context/release-order.md`, all three round-02 artifacts in full, `organizer-notes.md`, and the
cited prior passages before writing. I wrote only this file. I implemented nothing, published nothing,
committed nothing, edited no other participant's artifact, spawned no roster identity, altered no
accounting, worktree or budget state, and did not touch the unrelated historical idea whose
facilitator declaration blocks preflight. I drove no browser. The `[codex-1] <slug>: …` commit prefix
and all integration remain the organizer's.

<!-- Original language of the owner's quoted words and release decision: Slovak; verbatim source preserved in 00-prompt.md and source-context/{owner-brief,release-order}.md, English direction quoted above per the round instruction. -->
