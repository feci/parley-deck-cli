---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
round: 2
date: 2026-09-24
responding-to: [claude-1/round-01, kimi-1/round-01]
---

## Protocol context attestation

Obtained via `parley protocol packet --dir . --phase 2 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`.

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "fallback_reason": "absent",
  "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer/.parley-runtime/protocol-packets/full-phase2-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"
}
```

Body read in full (1,386 lines / 109,928 bytes); `shasum -a 256` of `body_path` =
`8ce83cde…9db7` = `source_sha256` = `packet_sha256` = the kickoff authority SHA and the live
`parley-deck/COOPERATION.md`. The renderer's JSON carries a `shadow` block describing the
optimized packet that was **not** used (`31903915…`, 68,912 B, 36 of 69 blocks); I read the full
body, not the shadow. HEAD at write time `ace284d` (organizer's round-preservation commit; adds
only idea artifacts — `git diff 59b0458..HEAD -- parley-deck/COOPERATION.md internal/` is empty).
origin/main has moved to `54e0798` (release-1.49.1 track: VERSION/CHANGELOG/procctl comments);
`git show origin/main:parley-deck/COOPERATION.md | shasum -a 256` = `8ce83cde…` — the protocol
text is byte-identical there, so all protocol verdicts below hold for the integration baseline
too. Transport header `github-pr`; this run operates under the owner's per-run local-canonical
override recorded in `00-prompt.md`.

Driver-gap accounting: round-02 was launched via the owner-authorized documented manual fallback
(driver halted on the unavailable historical worktree `…/scratchpad/f2repo`, exactly the
prior-run gap; preserved in `organizer-notes.md` and
`inbox/claude-to-user_…_driver-error.md`). I did not alter accounting, worktree history or code
to bypass it.

## Position changes since round 1

I change three positions after cross-reading and re-verifying; the rest of my round-1 proposal
stands.

1. **Unavailable designated agent at kickoff: I retract loud-fallback and adopt fail-closed with
   a §9.0-style recorded waive.** My round-1 tier-4 fallback ("record in the readiness table and
   an inbox note, then fall to claim → drafter") silently discards the owner's direct
   instruction. Both peers converged on fail-closed independently, and the decisive evidence is
   §9.0 itself, which I re-read at L868-877: excluding even a *mere quorum member* for
   unavailability "requires **explicit user confirmation** and is recorded in `00-prompt.md`"
   (`excluded: [… — confirmed <date>]`), and a mid-idea unavailability "falls to §5 and the
   runtime watchdog, downgrading to the same per-idea, user-confirmed waive". A designation is a
   stronger owner act than quorum membership; silently reverting it should be strictly harder,
   not easier. The codebase's own value system says the same (`driver_impl.go:47-51`: the ops
   "NEVER silently falls back" for role eligibility). New position, matching claude-1 D4 and
   kimi-1 #3: kickoff ping failure of a designated agent = blocking preflight gate, released
   only by (a) editing the designation, (b) an owner-confirmed recorded waive naming the
   fallback (`implementer_waived: <id> — <reason> — confirmed <date>`, the §9.0 `excluded:`
   shape), or (c) removing the designation. Mid-Phase-5 death escalates (LE-5 shape); the
   partial tree is evidence and reassignment is a human decision. This is not a hard stall: it
   is the same confirm-and-record loop §9.0 already prescribes, and this run's own history shows
   the owner answers such gates (the release-order inbox note came back within hours).
2. **Global-default vs claim precedence: I adopt claude-1's Phase-0 materialization (D3), which
   settles my round-1 disagreement with kimi-1 structurally.** I keep my round-1 ordering
   (designation/global above claim) but materialization makes it observable and race-free: when
   `[defaults].default_implementer` is set and `00-prompt.md` has no `implementer:` key, kickoff
   writes the concrete id into the idea with a provenance comment (the `ResolveRoster`
   precedent, `internal/config/roster.go:66-95`, verified); after Phase 0 exactly one rule is in
   force and one place holds it, and a later machine-config change cannot reinterpret an
   in-flight idea. Against kimi-1's O1 alternative (claim > global default): a claim is a
   first-mover act in a time window — exactly the speed-based selection the owner asked to
   remove; letting it defeat the owner's standing preference re-opens that door. The pressure
   valves are the per-idea field (owner/author can name someone else for one idea) and §4
   escalation (a participant who believes they are better placed produces a record the owner
   acts on — claude-1's D7 argument, which I adopt verbatim in spirit). I add one refinement
   neither round-1 proposal had, resolving the "present-empty semantics" question: the field has
   **three states** — key absent (materialization may apply), key present but empty/`none`
   (**explicit non-designation**; blocks materialization for this idea — the per-idea opt-out
   against one's own global default, which all three round-1 proposals were missing), key
   present with an id (designation; validated fail-closed). This mirrors `facilitator:`'s
   empty-equals-undeclared parsing (`facilitator.go:42-48`, verified) while giving the empty
   value a *meaning* precisely because materialization exists.
3. **The legacy FINAL read-set widening moves OUT of this idea's shipped delta.** My round-1 §8
   unified the two resolvers but did not confront claude-1's D10(ii) (widen the FINAL.md key set
   to `author/drafter/finalized-by`). Having verified the numbers myself (51 of 85 FINAL.md
   files resolve to nobody today — my recount matches claude-1's exactly), the widening is a
   real default-path behavior change: those 51 files would start resolving a drafter where
   today the driver falls to `participants[0]` and consensus fails closed to the full list. The
   owner's constraint is exact: "Ship the mechanism with the global default UNSET, so behaviour
   without it is exactly today's." Behaviour means observable behaviour, not documented intent.
   A resolver that reads new keys on decks that never opted in violates that sentence even with
   the global unset. I propose the split that keeps the owner boundary intact:
   **in scope** — close the *write* end only (driver FINAL-draft prompt and the §4 Phase 4
   template emit `implementer: <drafter-id>`; the resolver already reads that key; existing
   artifacts untouched; new ideas get a working documented chain, with the one edge —
   drafter ≠ `eligible[0]` when `participants[0]` is not headless/found — recorded as an
   `## Agreed trade-offs` entry per claude-1's Q3); **out of scope, follow-up idea** — the
   read-set widening and any claim auto-honoring (see response to kimi-1 O2). If the group wants
   the widening anyway, it must be an explicit, owner-visible deviation from the quoted
   constraint, decided in `## Agreed trade-offs`, not carried inside "repair".

Two round-1 positions I keep after being challenged: field name `implementer:` (2-of-3 with
claude-1; kimi-1 marks it minor and defers), and the shared-reader/two-tails resolver shape,
which claude-1's D10 independently converged on from the other direction.

## Responses to others

### @claude-1 — round-01

Your verification is the strongest artifact of round 1 and I re-ran its load-bearing claims at
HEAD `ace284d`; verdicts with my own locators are in the §15 record below. F-1 through F-9 all
**CONFIRMED** — including the two test pins (`app_test.go` `want codex (fallback; hermes not a
participant)`; `roundgate_test.go` `TestUnresolvableImplementerExpectsEveryone`), the FINAL.md
census (my recount: 85 total / 11 `implementer` / 25 `drafted-by` / 32 `author` / 23 `drafter` /
4 `finalized-by` / 51 unreadable — identical to yours), the stale `meta/version.json` SHA
(`git show b4831d6:parley-deck/COOPERATION.md | shasum -a 256` = `12e4b31c…` exactly), and the
skill-copy hash (`fc907e59…`, byte-equal to the published core 2.13.0 per the prior handoff).
Three corrections, none of which wounds your design:

- **C-a (locator):** the "two-participant rule" is at `COOPERATION.md:746` (§5: "For ideas with
  **only two participants**, the same rules apply unchanged…"), not L466-468 — those lines are
  the IMPLEMENTATION.md template (`design-pr:`/`implementation-pr:`). Your D6 argument survives
  with the corrected locator; §4.0 L262 (standard's 2 reviewers degrade to 1) is the companion.
- **C-b (precedent precision, minor):** `internal/config/runtime.go:637-641` emits `ping_tier` /
  `preferred_transport` / `roster_change_policy` as **active keys with inline comments**, not
  commented-out lines. D3's "write it as a commented line" is therefore a *new* emission shape,
  not the existing one. I still support emitting `default_implementer` as a commented line
  (absent key = today's behaviour is the contract we must not blur), but FINAL should describe
  it as a deliberate deviation from the emission precedent, one sentence.
- **C-c (attribution, and it is mine, not yours):** my round-1 file claimed "the two copies
  agree in practice only because the tool-written FINAL template emits `drafted-by`
  (`consensus.go:805`)". That is **wrong** — `consensus.go:799-810` is the *review-consensus*
  (Phase 7) template; no tool-written FINAL template emits a resolver-readable key (your F-3).
  I retract that sentence; your F-2/F-3/F-6 supersede it, and F-6 means the divergence is not
  merely theoretical. Logged in my §15 record.

On your design: D4 adopted (position change 1). D3 adopted with the three-state refinement
(position change 2). D6 mechanical drafter-preference: **support**, with two conditions — it
must fire only when a designation exists (otherwise it changes the default path, violating the
owner boundary) and the degenerate fallback (no other eligible drafter → implementer drafts)
must be test-pinned. The soft §15.5-style concentration line when drafter == implementer: fine,
one sentence. D8 gate relocation: **support**, designated-ideas-only trigger, same reasoning.
D10 shape (one exported `ResolveImplementer` returning `(id, src, ok)`, per-caller tails
preserved and pinned): **support** — it is my round-1 §8 with your better API; your MIN-7
prior-art disclosure is adequate (the finding was zero-callers + doc mismatch; neither holds
here, and the unexported-helper fallback is acceptable if the group disagrees). D10(ii) widening:
**opposed in this idea's delta** — position change 3; I hold this firmly against your
in-scope argument, on the owner's exact words, not on taste. D11: support, with your own R-3
made normative — `facilitatorConflictGates` iterates `status.Ideas` (`preflight.go:327-339`,
verified), which is why an unrelated historical idea blocks this run; the new gate family must
scope to the idea being launched, or we ship a second class of cross-idea blocking. D7
report-only for claims: agreed, except your "preflight gate when a claim exists and no
designation does" — that gates today's normal world (claims exist today without designations);
advisory notice only. Q4 (`--no-implement` visibility at preflight): worth one implementation-
time check, not a design blocker; the gate should simply not fire for ideas that will not reach
Phase 5. Q5 re-entry vs corrected designation: keep pin-first; the correction path is an
explicit, logged frontmatter edit by the incoming implementer (the pin is the record, and the
record's owner during execution is the implementer) — one sentence of amended text says so.
Q8/Windows: superseded by the owner's release-order decision (see New concerns #5).

### @kimi-1 — round-01

Your A1-A5 all **CONFIRMED** at HEAD with my own locators (§15 record), including A3
(`trajectory.go:65` `implementer` = "frozen patch author", required by `trajectory configure`,
never consulted by `newDriverImplOps`/`resolveImplementer` — grep clean) and the lean-organizer
incident (V12): `lean-organizer/…/organizer-notes.md:41` records verbatim that continuing at the
FINAL boundary "printed 'drafting FINAL via claude-1' despite Kimi's published FINAL and claim",
stopped by the organizer, followed by "explicit authorized manual fallback"; line 38 records
"Single implementer: zcode-1, by its pre-Phase-5 inbox claim" — the claim worked because a human
read it, which is A1 in one sentence.

On your open questions, answered directly as you asked:

- **O1 (claim vs global):** global default (materialized) > claim; rationale in position change
  2. The one-line summary: your ordering re-institutionalizes first-mover advantage against the
  owner's stated motivation; the per-idea field and escalation are the designed pressure valves.
  I do not read your O1 as a substantive block — you asked for the attack, and this is my best
  one; if you still prefer claim > global after reading it, that is a real disagreement for
  round 3, not a silent one.
- **O2 (claim machine-visibility):** report-only in this idea (advisory preflight/driver notice,
  no gate, no auto-honoring). Auto-honoring restores the race; the deterministic-detection rule
  (uncontested, single, pre-work claim) is a failure-mode zoo (contested claims, claims after
  work began, stale claims) that all end in escalation anyway; the social layer worked in the
  one observed run. Your A1 stays true and *documented*; the follow-up idea can give it a reader
  if evidence ever shows the organizer link failing.
- **O3 (field name):** keep `implementer:`. Symmetry with `facilitator:` is worth more than
  grep-disambiguation from `IMPLEMENTATION.md`'s record key — the protocol already lives with
  multi-file same-named keys (`author:`, `participants:`), the two files are different
  namespaces with different lifecycle stages, and claude-1's `ImplementerSource` enum makes
  provenance unambiguous in code, events and stdout. Your own artifact marks this minor; I
  agree it is not worth a third round.
- **O4 (two-participant ideas):** kickoff *warning*, not a gate, when a designation leaves one
  reviewer — §5 L746 says the rules apply unchanged to two-participant ideas, and a hard gate
  would make designated two-participant runs unsatisfiable. LE-11's ≥2-independent-reviewer
  auto-close floor stands where it lives (close time), and the warning surfaces it early. The
  zero-invokable-non-implementer-reviewer case *is* a gate (your #5, adopted).
- **O5 (deck-level layer):** agree — ship two levels only; record the deck layer as a named
  future extension in FINAL's non-goals so it is a decision, not an omission.
- **#7(e) trajectory-agreement gate:** support scoped — only when a trajectory policy exists
  *and* a designation exists; mismatch = gate. On the default path (no designation) trajectory
  stays exactly as today.

Your pre-dispatch FINAL validation (#4): **support scoped to designated runs only** — a
designation is what makes non-drafter execution load-bearing, so validating FINAL's required
sections and acceptance criteria before Phase-5 dispatch is designation-motivated; on
undesignated runs it would be a new gate on today's world and fails the owner boundary. Your
"validation gaps in dormant paths" risk is right and I promote it to an acceptance criterion
(see New concerns #3).

## New concerns / questions

1. **Materialization write mechanics.** Kickoff writing into `00-prompt.md` is precedented
   (§9.0 records the readiness result there) but must be additive-only, frontmatter-preserving,
   git-visible, and skipped in dry-run. One implementation note, not a design risk.
2. **Three-state parsing needs key-presence.** `ReadFrontmatter` returns `map[string]string`,
   where absent and empty are indistinguishable; implementing present-empty = explicit
   non-designation requires distinguishing the two (presence check before value check). Small,
   but it must be pinned by test or the opt-out state silently collapses into "absent".
3. **Dormant-path coverage is an acceptance criterion.** The global default ships UNSET, so
   materialization, its validation, and its inapplicable-fallthrough see no live traffic until
   the owner acts. Tests must cover all three states now (kimi-1's risk, promoted); the
   first real use of the owner's default must not be its first test.
4. **Where default-path behavior is preserved vs designation-only.** I want FINAL to carry a
   two-column table: left = unchanged-by-construction (no new keys read; both resolver tails and
   their test pins byte-preserved; no new gates on undesignated ideas; FINAL prompt/template
   key emission affects new ideas only, recorded as the one agreed trade-off), right =
   designation-only paths (materialization, fail-closed validation, waive recording, drafter
   preference, relocated diversity gate, pre-dispatch FINAL validation, resolved-source event).
   The write-end repair is the single item that changes new-idea default behavior and must be
   the headline `## Agreed trade-offs` entry.
5. **Windows and release boundary — resolved by owner decision, recorded for FINAL.** The prior
   handoff's unresolved Windows scope decision is now decided:
   `inbox/user-to-codex-1_…_release-order.md` (2026-09-24, status authorized) selects "Oboje" —
   CLI releases label Windows experimental/unvalidated, keep labelled Windows assets, hold the
   CLI winget PR, and a separate reviewed `windows-portability` idea removes the label later;
   release order is release-1.49.1 → **this idea** → windows-portability, each gated on the
   predecessor's done file. Consequences for us: this idea's release notes carry the
   experimental label, no CLI winget PR, and **no Windows platform architecture work is in
   scope** (that belongs to windows-portability; the brief's non-goals already barred it, and
   the owner decision confirms rather than authorizes widening). Non-code scope does not bypass
   attended publication: the designated implementer executes Phase 5-8 including release
   *preparation*, but `parley protocol publish` and core publication remain owner-attended
   (§7/§14; owner brief; kimi-1 #1 already said this correctly) — one sentence in the amended
   Phase 5 text keeps it explicit.
6. **No participant publishes before the organizer's release step** — reaffirmed against the
   moved origin/main: the 1.49.1 release commits on main are release-track work by the other
   run; this idea integrates on top of latest origin/main only after reviewed implementation,
   per the brief.

## Current proposal

Adopt the mechanism; consolidated position after round 2 (deltas from my round-1 proposal are
positions 1-3 above):

1. **Fields.** Optional per-idea `implementer:` in `00-prompt.md`, three-state semantics
   (absent / present-empty-or-`none` = explicit non-designation / present-id). Global
   `[defaults].default_implementer` in `~/.parley/agents.toml`, **shipped UNSET** (emitted as a
   commented line at init, documented as absent = today's behavior). No deck-level layer
   (named future extension).
2. **Phase-0 materialization.** When the global default is set and the idea's key is absent,
   kickoff writes the id + provenance comment into `00-prompt.md`; when the named agent is not
   a participant of this idea, record inapplicability (notice, not error) and fall through. A
   present-empty/`none` key blocks materialization.
3. **Precedence.** `IMPLEMENTATION.md` `implementer:` (re-entry pin) > `00-prompt.md`
   `implementer:` (owner-written or materialized) > perfected inbox claim (**only when no
   designation exists**) > FINAL metadata > positional `eligible[0]` (driver tail) / fail-closed
   full list (consensus tail). Mid-run corrections are explicit logged edits of the pin by the
   incoming implementer.
4. **Validation, fail-closed.** Present-id not in `participants:`, or naming the declared
   non-participating facilitator → blocking preflight gate (both fields named), scoped to the
   idea being launched (not deck-wide). Kickoff ping failure of a designated agent → gate,
   released by edit / owner-confirmed recorded waive / removal (§9.0 shape). Mid-Phase-5 death →
   escalate, human reassignment. Global set-time validation against the active machine roster.
5. **Resolver unification.** One shared reader in `internal/protocol` returning
   `(id string, src ImplementerSource, ok bool)` over the identical chain; driver and consensus
   keep their own tails, each pinned by its existing test, plus a pin for the facilitator-filter
   axis. Read set unchanged (`implementer`, `drafted-by` only) — no legacy-key widening in this
   idea.
6. **Write-end repair (the one default-path trade-off).** Driver FINAL-draft prompt and §4
   Phase 4 template emit `implementer: <drafter-id>`; existing artifacts untouched; recorded in
   `## Agreed trade-offs` with the not-headless edge called out.
7. **Claims.** Social layer unchanged; advisory detection surface only (no gate, no
   auto-honoring); auto-honoring and read-set widening = named follow-up idea.
8. **Designation-only additions.** Drafter-preference skip in `firstEligibleHeadlessAgent`
   (fires only under a designation; degenerate fallback pinned); kickoff relocation of the
   `require_model_diversity` gate; hard gate on zero invokable non-implementer reviewers;
   kickoff warning (not gate) on two-participant designated ideas; pre-dispatch FINAL section/
   acceptance-criteria validation; durable `agent.implementer_resolved` event + one-line stdout
   notice; trajectory-policy agreement gate when both exist.
9. **Text surfaces.** Identical hunks in all three COOPERATION.md copies (deck `8ce83cde…`,
   defaults `cce5d7d9…`, skill reference `fc907e59…` — all verified in lockstep modulo project
   zones) + `meta/protocol-changelog.md`; §0 `[defaults]` paragraph, §4 Phase 4/5, §4.0 Phase 0
   template, §9.0 readiness note, §10 TL;DR item 6. SKILL.md restates no Phase-5 rule (grep:
   only the facilitator exclusion at L121 and advisory-roles at L202), so only the references
   copy changes — my round-1 concern #5 is resolved by inspection.
10. **Role facts unchanged.** Designee = full design participant and signatory; never reviews or
    goal-checks own work (L234/L538/L688-691; `driver_impl.go:415-417` already enforces
    checker ≠ implementer); exactly one implementer at a time; publication stays owner-attended.
11. **Windows/release.** Experimental label, no CLI winget PR, no platform work; release only
    after `release-1.49.1` done file; version chosen above actual releases at staging time.

**Readiness for consensus: ready.** No remaining substantive block from my side. The two splits
(legacy read-set widening out; claim auto-honoring out) are recorded counter-proposals to
claude-1's in-scope position; if he maintains D10(ii) in-scope, round 3 must resolve it as an
explicit owner-visible deviation decision — it is the only question I would still spend a round
on. Everything else above is either three-way agreement already or a two-of-three with the third
marked minor.

**Volunteering under the current protocol (this run only, no weight in the owner's future
global-default question):** I remain willing to draft or implement, per my round-1 statement.
For single-drafter/implementer coordination I propose — as a *proposal* only, to be perfected
through today's channels (`Drafter: yes` in the consensus signoff; `inbox/<from>-to-all_<slug>_
impl-claim.md` before work begins): **claude-1 drafts FINAL** (his verification base is the
strongest and drafter ≠ implementer should be the norm this mechanism makes common), **zcode-1
implements** (willingness on record; lean-organizer precedent as capability evidence, n=1 and
not a default argument). kimi-1 holds a standing willingness for either. If kimi-1 prefers to
draft, I withdraw the proposal shape and keep only the implementer claim; the group should
settle this at the consensus boundary, not in round files.

## Existing alternatives

Per §15.6(a) — mechanisms this consensus would build or change by hand, against what ships
today (verified at HEAD `ace284d`; all locators re-checked this round):

1. **Today's implementer selection (the mechanism being extended):** protocol prose
   `COOPERATION.md:443` (default = FINAL drafter; claim via `inbox/<from>-to-all_<slug>_
   impl-claim.md`); machine layer = the two unexported resolvers reading only
   `IMPLEMENTATION.md` `implementer` and `FINAL.md` `implementer`/`drafted-by`
   (`internal/app/driver_impl.go:104-131`, `internal/consensus/consensus.go:751-774`), driver
   tail `eligible[0]` positional, consensus tail `""` → full list. The claim file has no code
   reader (grep `impl-claim` in `internal/` → none).
2. **Declared-facilitator machinery (the pattern mirrored):** `internal/protocol/facilitator.go`
   (optional key, absent = byte-identical, empty value = undeclared at :42-48, `Conflict()`
   fail-closed), preflight `facilitatorConflictGates` (`internal/app/preflight.go:314-340`),
   driver eligibility + `roleErr` escalation (`internal/app/driver_impl.go:45-98`), shipped
   1.49.0 (CHANGELOG).
3. **`roles:` advisory map:** `COOPERATION.md:311`/`:97` — advisory by construction, parsed by
   no code; cannot appoint an executor.
4. **Roster-preset resolution (the materialization precedent):** `internal/config/roster.go:66-95`
   `ResolveRoster` — explicit → track default → none, fail-closed validation, `Provenance`
   one-liner; selects participants, not executors; pattern reused for the global default.
5. **The two resolvers (the code being unified, tails preserved):** as in 1; consumers
   `driver_impl.go:64`, `consensus.go:734` via `expectedRoundParticipants` (:730-748), and
   `internal/driver/phasedigest.go:363` through `ExpectedRoundParticipants`
   (`consensus.go:1085-1091`) — the `parley wait` surface where the divergence is user-visible.
6. **Trajectory `--implementer`:** `internal/app/trajectory.go:65,104-122` — freezes the
   accountable patch author into the evidence policy; never consulted by driver selection; the
   agreement gate (proposal #8) makes the two agree when both exist.
7. **Model-diversity gate (relocated, not rebuilt):** `checkModelDiversity`
   (`internal/app/driver_impl.go:178-210`), LE-3/`require_model_diversity` — fires at Phase 6
   today; proposal moves the designated-case check to kickoff.
8. **Scoped null:** do nothing. Rejected: the owner asked for a settable executor role; today no
   settable channel is honored at Phase-5 start (lean-organizer incident,
   `organizer-notes.md:38,41`).

## Verification record (§15)

### Verdicts on claude-1's first-canonical claims (I am the non-owner verifier; all PRIMARY — locator read at HEAD `ace284d` in this worktree)

| Claim | Verdict | My evidence |
|---|---|---|
| F-1 implementer = `eligible[0]` list order, frozen at driver construction | CONFIRMED | `driver_impl.go:45-98` (filter → `resolveImplementer(ideaDir, eligible)` → struct field); `Implement()` launches it |
| F-2 FINAL.md census: 85 files, 51 unreadable (11/25/32/23/4 key counts) | CONFIRMED | my independent recount reproduced every number exactly |
| F-3 driver FINAL prompt requires only `idea:`/`status:` | CONFIRMED | `driver_consensus.go:174-199` read in full |
| F-4 drafter and implementer use different eligibility filters | CONFIRMED | `driver_consensus.go:109-123` (Found + headless) vs `resolveImplementer` (participants only) |
| F-5 `impl-claim` has no implementation | CONFIRMED (independently established in my round-1) | grep `internal/` → none |
| F-6 divergence user-visible via `parley wait` | CONFIRMED | `internal/driver/phasedigest.go:363` → `consensus.go:1085-1091` → `expectedRoundParticipants` |
| F-7 `--participants` order is an undocumented implementer switch | CONFIRMED | `parley run --help`: "-participants string — comma-separated agent IDs to run", no order semantics stated |
| F-8 preflight freshness reports remembered SHA (`12e4b31c…` = COOPERATION.md at `b4831d6`) | CONFIRMED | `preflight.go:486-487`; `meta/version.json`; `git show b4831d6:parley-deck/COOPERATION.md \| shasum` = `12e4b31c…`; source-role short-circuit at `:489-492` makes it currently harmless |
| F-9 three protocol copies in lockstep (project zones only) | CONFIRMED | live hashes `8ce83cde…`/`cce5d7d9…`/`fc907e59…`; diff-line counts differ from claude-1's (11/13 vs his 16/18) only because I normalized trailing whitespace — every differing line is a project zone; no factual conflict |
| Test pins: positional tail (`app_test.go`) and `""` tail (`roundgate_test.go`) | CONFIRMED | both read; quoted strings match |
| Two-participant rule at L466-468 | **WRONG locator, rule exists** | actual `COOPERATION.md:746` (§5); L466-468 is the IMPLEMENTATION template |
| `runtime.go:637-639` emits commented lines for defaults | **WRONG (minor)** | emits active keys with inline comments (`speed`, `ping_tier`, …); a commented `default_implementer` line would be a new emission shape |

### Verdicts on kimi-1's first-canonical claims

| Claim | Verdict | My evidence |
|---|---|---|
| A1 claim has no machine reader | CONFIRMED | grep; resolver sources at `driver_impl.go:113-130`/`consensus.go:760-777` |
| A2 documented FINAL template cannot feed the resolver | CONFIRMED | `COOPERATION.md:413-419` template prescribes `author:`; resolvers read `implementer`/`drafted-by` only |
| A3 trajectory `--implementer` is not a selector | CONFIRMED | `trajectory.go:52-65,104-122` ("frozen patch author", checks-contract requirement); no consultation in `driver_impl.go`/`app.go` (grep) |
| A4 preflight gate machinery exists and is extensible | CONFIRMED | `preflight.go:62-83,289-301,314-340` |
| A5 resolver drift happened before (MIN-7/F15) | CONFIRMED | lean-organizer review round-01 claude-1.md findings; F15 drop recorded in its IMPLEMENTATION.md |
| V12 operational incident (driver drafted via claude-1 despite claim; manual fallback) | CONFIRMED | lean-organizer `organizer-notes.md:41` verbatim; :38 for the claim-based implementer |
| V14 skill copy unverified by kimi | now CONFIRMED by me | hash `fc907e59…` = published core 2.13.0 (prior handoff) |

### Conflicting-verdict resolutions (§15.3)

1. **Brief hypothesis "speed decides":** my round-1 CONFIRMED-with-qualifications, claude-1
   WRONG-for-driver-path, kimi-1 PARTIALLY CONFIRMED — **no residual factual conflict**: all
   three agree the prose chain (first filer → drafter → implementer) holds socially
   (`COOPERATION.md:407-409`, `:443`) and that no code computes first-filer (driver = list
   order, `driver_consensus.go:109-123`; resolver = frontmatter + positional). The operative
   design fact is claude-1's F-1: under driver orchestration the selector is list order. FINAL
   should state it exactly this way.
2. **"Implementation is the heaviest token work of a run" (brief):** claude-1 UNVERIFIED vs
   kimi-1 V13/my round-1 CONFIRMED-as-quoted — **resolved**: the study supports "~30-40% of an
   implementing organizer's input" (with its self-described caveat); it does **not** rank
   implementation as the heaviest work of a run (its own numbers put re-orientation ~33-39% and
   compaction 42% of uncached input). The ranking must not enter FINAL as established; the
   30-40% figure may, with the caveat. No participant's design depends on the ranking.
3. **My own round-1 error, corrected:** "the two copies agree in practice only because the
   tool-written FINAL template emits `drafted-by` (`consensus.go:805`)" is **wrong** — that
   template is the Phase-7 review-consensus artifact (`consensus.go:799-810`); no tool-written
   FINAL template emits a resolver-readable key (claude-1 F-3). Retracted; superseded by
   F-2/F-3/F-6.
4. **claude-1's two-participant locator** — corrected to L746 above; argument unaffected.
5. **diff-line counts for copy lockstep (F-9)** — method difference (trailing-whitespace
   normalization), not a factual conflict; both support "project zones only".

### Owner-constraint reconciliation (required by this round's task)

The owner's exact boundary: "Ship the mechanism with the global default UNSET, so behaviour
without it is exactly today's." My consolidated proposal keeps it as follows — **preserved by
construction**: no new frontmatter keys are *read* by the resolver (read set unchanged); both
resolver tails and their existing test pins preserved through the shared-reader refactor; no new
blocking gates fire on undesignated ideas (gate family scoped to the idea being launched and to
designation-present states); claims stay advisory; global key absent on ship and emitted
commented. **Designation-only paths** (new behavior, only when a designation exists or the owner
later sets the default): materialization, three-state semantics, fail-closed validation, waive
recording, drafter preference, relocated diversity gate, pre-dispatch FINAL validation,
resolved-source event. **The single recorded default-path trade-off**: the write-end repair
(new finals carry `implementer: <drafter>`), which changes new-idea resolution in one narrow
edge and is explicitly carried in `## Agreed trade-offs`. The legacy read-set widening and claim
auto-honoring are excluded from this idea precisely because they cannot be reconciled with the
quoted sentence; participant agreement that they are "really defects" is not evidence and not
authority to widen the boundary (the organizer's notes state the same). No unrelated defect
repairs are smuggled in: the procctl/version changes now on origin/main belong to the other
run's release and are untouched by this idea.

No other participant's file was edited; no code was changed; no commit, push, release or
publication was made; roster identities were not spawned. Round-02 is cross-review 1 of the
deliberation cap of 3.
