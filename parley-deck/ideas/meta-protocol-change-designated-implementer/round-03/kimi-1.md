---
agent: kimi-1
idea: meta-protocol-change-designated-implementer
round: 3
date: 2026-09-24
responding-to: [claude-1/round-02, zcode-1/round-02]
---

## Protocol context attestation

Obtained via `parley protocol packet --dir . --phase 2 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json` at HEAD
`73ee923151dbdcdc95ed34a614fee7f76692f49f` (branch `designated-implementer`; the commit adds only this
idea's round-02 artifacts and `source-context/release-order.md` over round-02 HEAD `ace284d` —
`git diff ace284d..HEAD -- internal/ parley-deck/COOPERATION.md` is empty, so every verified code and
protocol locator from rounds 1-2 remains addressable):

- `context_mode`: `full`; `source_sha256` = `packet_sha256` =
  `8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7`; `fallback_reason`: absent.
- `body_path`: `.parley-runtime/protocol-packets/full-phase2-deliberation-8ce83cde…md`. The packet is
  byte-identical to round-02's (same SHA trio). I read that body in full (1,386 lines) in round 2 and
  re-verified this session that `shasum -a 256` of the body equals `shasum -a 256
  parley-deck/COOPERATION.md` = `8ce83cde…` — the packet is the live authority verbatim. Operative
  sections re-read in the live file this session: Phase 2/4/5 (L333-360, L405-440, L441-446), §9.0
  readiness/exclusion (L868-877), §12.5 seeding contract (L1153-1154). The shadow block
  (`31903915…`) was again noted and not used.
- This round-03 file is written via the documented manual fallback recorded in `organizer-notes.md`
  (driver halted on the unavailable historical worktree `f2repo` — the unchanged historical-worktree
  accounting gap). I alter no accounting, worktree history, or code around that gap; no code gates
  were bypassed; this is the only file I write.

## User direction (owner release direction — authorized, not an open question)

The controlling inbox direction, copied verbatim to `source-context/release-order.md` (original
Slovak preserved there as source evidence: *"Oboje" — "CLI 1.49.1 hneď vydať pre macOS a Linux cez
GitHub a Homebrew, Windows označiť ako experimentálne a winget pre CLI podržať. Zároveň otvoriť
samostatnú Windows ideu s review, ktorá chyby opraví a označenie odstráni."*), states in English:

> "Both — release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, label Windows as
> experimental, and hold winget for the CLI. At the same time open a separate reviewed Windows idea
> that fixes the defects and removes the label."

Release order (owner-set, to keep main linear): (1) `release-1.49.1` writes
`codex-1-to-user_release-1.49.1_done.md` in the lean-organizer worktree deck inbox; (2) **this idea**
releases only after that done file exists; (3) `windows-portability` releases after this idea's done
file, removes the experimental label and submits the held CLI winget PR once hosted windows-latest CI
is green. Until (3) ships: every CLI release labels Windows experimental/unvalidated, keeps the
Windows assets (labelled), opens no CLI winget PR, and the skill is unaffected.

Consequences for this run, stated as settled boundary rather than assessment: this idea **waits for
the release-1.49.1 done file**, ships with the Windows experimental label and assets retained, holds
the CLI winget PR, leaves skill channels unaffected, and does **no** Windows platform architecture
work — that belongs to the separate `windows-portability` run. claude-1's round-02 framing of the
Windows scope as an open owner question the organizer must carry "unclaimed" is superseded: the owner
has decided, the direction is authorized, and no participant needs to inherit or seek any
authorization here. My round-02 assessment of the mechanism's Windows-neutrality (pure Go, no
build-tagged files in the delta — corroborated by claude-1's file-level check) stands and is now the
only Windows-relevant fact this idea needs.

## Position changes since round 2

1. **I withdraw Phase-0 materialization and adopt claude-1's live layered read (C-2/S-4/S-7).** My
   round-02 change 2 adopted materialization on the strength of claude-1's own round-01 D3; claude-1
   has since withdrawn D3 on evidence I have now re-verified at HEAD, and the verification moves me:
   (a) `LoadDefaults` (`internal/config/runtime.go:419-438`) merges `[defaults]` live across
   `configLayers` (`:384-401`: central `~/.parley/agents.toml` → deck `parley-deck/agents.toml` →
   `agents.local.toml` → `$PARLEY_AGENT_CONFIG`), "later non-empty values win" — every existing
   `[defaults]` key is resolved live at use, and `LoadDefaults` iterates all layers with no
   membership filter, so a deck-level `default_implementer` is merged by construction (claude-1's
   K-5: **CONFIRMED, PRIMARY**); (b) the init template header itself states "Project-wide policy
   defaults; a deck's parley-deck/agents.toml overrides them" (`runtime.go:634`) — the deck layer
   arrives with the block, which retires my round-02 "ship two levels only" framing (my new-concern
   #2): the levels are per-idea > deck > machine, inherited, not invented, both config files shipping
   the key absent; (c) the freeze that matters post-work is already supplied by the
   `IMPLEMENTATION.md` pin both resolvers read first; pre-work, a live-read change is exactly an
   owner metadata edit, which all three of us already allow. Fairness note on my own round-02
   citation: §12.5 (L1153-1154) does let the driver author a **new** block's `00-prompt.md` as
   initiator-owned kickoff material — but materialization would **edit an existing, author-owned**
   idea file's frontmatter, which no tool does today; my precedent was weaker than I credited.
   Resumability under live read (dispatch question): `parley run`/`continue` re-resolves at each
   launch — deterministic from durable inputs (frontmatter + layered config); a mid-idea global
   change pre-pin resolves differently, which is owner intent surfaced by the recorded source, not a
   bug; post-pin the pin wins. Materialization's one residual benefit — insulating an in-flight idea
   from a later machine-config change — is thinner than its cost: a bespoke frontmatter-write path
   (idempotency, dry-run, provenance-comment machinery) plus per-idea repair when the owner changes
   the global, versus one central edit under live read.
2. **Unavailable designee at the global-default layer: I move from loud fall-through to
   fail-closed-with-pre-built-exits (claude-1 S-5, zcode-1 position change 1).** zcode-1's §9.0
   citation is decisive and I verified it verbatim (L868-877): excluding even a *mere quorum member*
   for unavailability "requires **explicit user confirmation** and is recorded in `00-prompt.md`
   (`excluded: [<id> — reason — confirmed <date>]`)", and a mid-idea unavailability "downgrades to the
   same per-idea, user-confirmed waive". A designation is a stronger owner act than quorum
   membership; reverting it silently (even loudly) must be strictly harder, not easier. My round-02
   objection — a hard stall on every idea while one agent is down — is answered under live read: the
   owner unblocks all future kickoffs with **one central edit**, and any single idea with the
   pre-built per-idea exits (re-designate, recorded waiver, or `none`). The layered rule that
   survives: **non-participant** global designee = inapplicable (notice + fall through, never a gate —
   the global predates the idea's roster); **participant-but-unavailable** designee at either layer =
   blocking gate with the §9.0-shaped exits. Invalid designations (non-participant per-idea field,
   pure facilitator, blank) stay hard failures.
3. **Empty/`none` semantics: synthesis.** Adopt claude-1's S-3 for the blank case — present-empty →
   blocking gate naming the field (I verified the mechanics: `ReadFrontmatter`,
   `internal/protocol/workspace.go:390-394`, does `strings.Cut(line, ":")` then `meta[key] = value`,
   so `implementer:` creates a `""` entry and `raw, present := meta[...]` distinguishes blank from
   absent in one line). Adopt zcode-1's three-state insight where it carries semantic load: the
   reserved literal `implementer: none` = **explicit per-idea non-designation** (skips the
   global-default tier for this idea), validated as a reserved token before any id check so no roster
   id can collide. Blank is the typo case and gets the gate; `none` is the opt-out case and gets a
   meaning. Both are designation-path only: a deck that never writes the key is untouched.
4. **Resolved-source event: I move from always-on to designation-only emission.** My round-02
   default-path table carried the `agent.implementer_resolved` event + stdout notice as always-on
   additive. Held against the exact owner sentence — "behaviour without it is exactly today's", with
   zcode-1's correct gloss that behaviour means *observable* behaviour — a new event on the unset
   path is a change, even at zero selection effect. Emitting only when tier 2 or tier 3 fired keeps
   this idea's unset path byte-identical with **no** trade-off entry needed, which all three of us
   can sign (claude-1 explicitly signs the conservative variant; zcode-1 already lists the event
   under designation-only additions). The always-on variant's genuine value — making list-order
   selection visible on today's decks (F-1/F-7) — belongs to the follow-up repair slug's evidence
   discussion, where the boundary can be relaxed deliberately. claude-1's required-recorded-decision
   D-1 is thereby dissolved by decision, not silence.
5. **Precedence restatement correction (dispatch item 4).** My round-02 "Current proposal" ladder
   listed "perfected uncontested claim (only when no field)" as a machine rank above FINAL-drafter
   inheritance — contradicting my own same-file position that the claim gets **no machine honoring at
   all**. zcode-1's round-02 has the identical internal inconsistency (proposal #3 lists a claim rank
   while #7 says social-only, no auto-honoring). Correcting mine and asking zcode-1 to correct
   theirs: the claim is **not a machine rank anywhere**; the sole machine chain is claude-1's S-4
   (adopted verbatim below). The claim stays what it is today — a prose/social act a human or
   organizer may materialize into the per-idea field — plus the one subordination sentence: under a
   present designation a claim never overrides it; the claimant escalates per §4. claude-1's block
   B-1 is dissolved: the parsing proposal it targeted was already withdrawn in my round-02, and the
   residual ladder wording is now corrected.

Positions kept unchanged: field name `implementer:` (present-empty/non-`none` id validated
fail-closed); the deferral slug `meta-protocol-change-implementer-inheritance-repair` for the legacy
read-set widening **and** the write-end emission (see @zcode-1); fail-closed per-idea unavailable
designee with the three recorded exits; pin-conflict escalation; claim social-only; designation-gated
drafter preference and Phase-0 diversity relocation; the L433 self-containment sentence; pre-dispatch
FINAL validation cited as shipped (`protocol.ValidateFinal`), not rebuilt; unification with both
tails pinned; identical hunks in all three protocol copies + `meta/protocol-changelog.md`;
`default_implementer` emitted commented-out by `parley init` as a deliberate, one-sentence-documented
deviation from the active-key emission precedent (`runtime.go:635-639`, zcode-1's C-b, verified).

## Responses to others

### @claude-1 — round-02

**Verdicts on your new round-02 claims (I am a non-owner; all PRIMARY, re-measured at HEAD `73ee923`):**

- **C-2's LoadDefaults evidence: CONFIRMED** (locators and semantics quoted in position change 1),
  including K-5's deck-layer inheritance — `LoadDefaults` applies no membership filter, so the deck
  layer merges by construction.
- **C-3's mechanism for D10(i): CONFIRMED** — a new FINAL emitting `implementer: <drafter>` on an
  unset deck would change later re-entry resolution from `participants[0]` to the drafter; a changed
  selection on the unset path.
- **S-3's `ReadFrontmatter` mechanics: CONFIRMED** (`workspace.go:390-394`).
- **§9.0 as cited by zcode-1 and embedded in your S-5: CONFIRMED** (L868-877, quoted above).

**Substantive responses:**

- **C-2/S-4/S-7: adopted** (position change 1). Your consistency argument is the one that survives
  contact with the code: materialization would have made `default_implementer` the only `[defaults]`
  key with bespoke mechanics and the only tool-edit of an existing author-owned idea file.
- **S-3: adopted, with one addition requested** — the reserved `none` literal as the explicit
  opt-out (position change 3). This preserves zcode-1's genuine three-state contribution (the
  per-idea escape from one's own global default, which all three round-01 proposals missed) without
  giving the blank-typo case silent semantics. Small; asking your assent as drafter.
- **S-5: adopted at both layers** (position change 2), keeping the non-participant/inapplicable
  distinction you did not explicitly split: a global default naming a non-participant of *this* idea
  is not an invalid designation, it is an inapplicable one (notice + fall through), or every mixed
  roster makes the global unusable. Your R-1 point 2 is decisive for the available-but-ping-failed
  case: falling through lands on `participants[0]` — the accidental selection the owner is replacing.
- **S-6: adopted, with authorization pinned.** Pin wins by default; pin-vs-designation conflict
  escalates rather than silently honoring; the correction is an **owner-recorded decision** (the
  `<id> — reason — confirmed <date>` shape, naming the abandoned attempt — the partial tree is
  evidence), and only then does the incoming implementer update the pin/field as the logged executor
  of that decision. zcode-1's "incoming implementer edits the pin" is the right mechanics with the
  wrong authorizer unless this sentence is attached: no implementer self-authorizes a reassignment.
- **S-8: decided designation-only** (position change 4). D-1 dissolved.
- **B-1: dissolved** — the machine-honoring proposal was withdrawn in my round-02; the residual
  ladder inconsistency you did not catch (mine, and zcode-1's #3) is corrected in position change 5.
- **K-1 (ValidateFinal already shipped): agreed and kept** — my round-02 had already withdrawn the
  rebuild; FINAL cites `finalsections.go:92` + the L433 sentence, nothing more.
- **X-5 (eligibility axis): carried.** Your requirement stands in the FINAL spec: the shared
  resolver takes the facilitator-filtered list as an explicit parameter, and the consensus call site
  passes the same list the driver does, or unification preserves the divergence it exists to remove.
- **Windows: your open-question framing is superseded** by the authorized owner direction (User
  direction section above) — recorded as settled boundary, not as a carried recommendation.
- **Drafter offer: accepted** (role assignment below).

### @zcode-1 — round-02

**Verdict on your load-bearing new claim (I am a non-owner):** your §9.0 reading (L868-877) is
**CONFIRMED, PRIMARY** — quoted in position change 2; it moved my global-layer position.

**Substantive responses:**

- **Materialization: I ask you to withdraw it, as I now have** (position change 1). The three
  arguments you made for it are each answered without it: (a) *race-freedom/one-rule-in-force* — the
  stated precedence chain (S-4) is one rule, resolved from durable inputs, with the resolved source
  recorded when a designation fires; (b) *in-flight insulation from a later machine-config change* —
  the dangerous window (post-work) is already closed by the `IMPLEMENTATION.md` pin, which we all
  rank first; the pre-work window is owner prerogative under any design; (c) *owner-act-outranks-
  claim* — preserved directly by the chain: field > global > today's chain, with the claim social and
  subordinate, so a participant filing first still cannot defeat a standing owner instruction. What
  materialization buys beyond that is a bespoke tool-write into an existing author-owned
  `00-prompt.md` (my round-02 §12.5 precedent covers **authoring a new block file**, not editing a
  live idea — I over-cited it), idempotency/dry-run/provenance machinery, and per-idea repair when
  the owner later changes the global — versus one central edit under live read. Your three-state
  semantics survive the withdrawal intact via the `none` literal (position change 3); nothing you
  designed is lost except the write path. Your own C-b diligence (the commented-out emission is a new
  shape) is exactly the kind of bespoke-mechanics cost that live read avoids wholesale.
- **Non-rostered/unavailable global designee (dispatch item 1, second half):** under live read the
  rule is deterministic: global designee not in this idea's `participants:` → inapplicable, notice +
  fall through to today's chain (you and I agreed here already; it is not an error for a standing
  preference to predate a roster). Global designee who is a participant but fails the §9.0 ping →
  your own position-change-1 gate, extended one layer up: blocking gate with the pre-built Confirm
  (re-designate in the idea field, record `implementer_waived: <id> — <reason> — confirmed <date>`,
  or `implementer: none`), the owner's single central edit unblocking everything else. Resumability:
  re-resolution at each launch is deterministic; the waiver/`none`/field-edit records make any
  correction auditable; post-`IMPLEMENTATION.md` the pin wins and a conflict escalates (S-6).
- **Write-end FINAL repair: opposed, and I hold this as firmly as you hold the read-set widening
  out.** Your position change 3 splits write-end (in, as the "headline trade-off") from read-set
  (out). The split fails your own test. The write-end emission matters *only because* it changes what
  the resolver later reads on unset decks: a new FINAL carrying `implementer: <drafter>` re-resolves
  future re-entries from `participants[0]` to the drafter (claude-1's C-3 mechanism, which I
  re-verified). That is the same observable selection change as the widening, one step removed —
  "behaviour without it" changed. And the boundary is not ours to trade: the owner's sentence is
  exact ("behaviour without it is exactly today's"), and an `## Agreed trade-offs` entry **records** a
  decision; it cannot authorize what an exact owner sentence forbids — your own condition ("explicit,
  owner-visible deviation") names the owner, not the participants, as the deviating authority. Both
  halves therefore go to the recorded slug `meta-protocol-change-implementer-inheritance-repair`,
  where the owner can relax the boundary deliberately with the evidence already gathered (51/85
  unreadable FINALs; the positional-tail destination). Deferral-with-slug is the explicit decision
  claude-1's Q3 demanded; it is not silence.
- **Your proposal #3 vs #7:** corrected per position change 5 — the claim holds no machine rank;
  please align #3 to your own #7.
- **X-5: credited and carried** (see @claude-1).
- **Skill-tree restatements:** your inspection result (SKILL.md restates no Phase-5 rule) resolves
  my round-01 concern; the implementation-time grep stays in the acceptance criteria as the cheap
  guard against future drift outside the three guarded copies.
- **Role assignment: countered, with reasons, below.**

## §15.3 conflict resolution — trace-of-prior-fixes (new PRIMARY evidence this round)

The contested claim (owner brief, testimony): claude-1's lean-organizer first code review found 1
CRITICAL + 8 MAJOR, "all fixed in later cycles." Round-1 verdicts: claude-1 `CONFIRMED`, zcode-1
`UNVERIFIED`; my round-02 split admitted the counts (both confirmed, PRIMARY) and ruled the "all
fixed" clause supported-not-traced. This round I performed the missing per-finding record trace in
the lean-organizer worktree deck (read-only; that idea's history untouched):

1. **Finding inventory:** `review/round-01/claude-1.md` carries CRIT-1 (L235) and MAJ-1…MAJ-8
   (L300-518); its own severity line reads "CRITICAL 1 · MAJOR 8 · MINOR 7 · NIT 3 (19 findings)".
2. **Ledger entry:** every one of the nine ids appears in `review/consensus-cycle-01.md`'s agreed-
   fixes ledger (CRIT-1 ×6 mentions, MAJ-1 ×6, MAJ-2 ×5, MAJ-3/4/5 ×4, MAJ-6/7 ×3, MAJ-8 ×7).
3. **Fix mapping:** `IMPLEMENTATION.md`'s per-cycle `### Fixes applied` sections map each id to a
   numbered fix entry (e.g., F1 = CRIT-1 "primary branch — both reviewers signed for removal";
   F2 = MAJ-1 + MAJ-2; G1/G2 carry later-cycle dispositions with reviewer signoffs recorded).
4. **Monotone close:** `outstanding_agreed_fixes` runs 21 (cycle 1) → 10 (cycle 2) → 7 (cycle 3) →
   **0** in `review/consensus.md` ("None — `outstanding_agreed_fixes: 0`") with three unconditional
   ✅ ACCEPT signoffs; `IMPLEMENTATION.md` frontmatter `status: complete`.

**Resolved verdict (I am a non-owner of the brief's claim):** the "all fixed" clause upgrades from
`UNVERIFIED` to **CONFIRMED as record-traced** — each of the nine findings has a ledger entry, a
mapped fix record, and a signed zero-outstanding close under the protocol's own Phase-8 default close
rule. Named residual, honestly scoped: the **code-content efficacy** of each fix was not re-audited
(that would mean re-reviewing each hunk against the committed code), and no party ever claimed it;
"fixed" means dispositioned-fixed through the protocol's own gate. The claim remains immaterial to
this idea's acceptance criteria — it feeds only the post-release default-implementer recommendation's
evidence base, where it now enters as "closure record-traced; code efficacy not re-audited."

Other §15.3 items standing from round-02, unchanged and unopposed in either peer's round-02: the
brief hypothesis is true of the prose/social path and false of the driver path (never state it
unscoped); "implementation is the heaviest token work of a run" stays `UNVERIFIED` — no new evidence
arrived, the study supports only the 30-40%-of-organizer-input share with its caveats, and the
ranking must not enter `FINAL.md` as established.

## The single deterministic chain (FINAL spec, no contradictory ladders)

Machine precedence, resolved per launch from durable inputs, exactly one chain:

1. `IMPLEMENTATION.md` `implementer:` — re-entry pin. If tier 2 also resolves and **disagrees**:
   escalate; correction requires an owner-recorded decision (`<id> — reason — confirmed <date>`,
   naming the abandoned attempt), then the incoming implementer logs the pin/field update.
2. `00-prompt.md` `implementer:` — present-blank → blocking gate naming the field; `none` → explicit
   non-designation (skip tier 3); id → fail-closed validation (participant; never the declared
   non-participating facilitator) at a preflight gate **scoped to the idea being run**, live
   frontmatter only (never `meta/version.json`).
3. `[defaults] default_implementer` via live `LoadDefaults` (per-idea > deck > machine inherited;
   both files ship the key absent; `parley init` emits it commented-out as a documented one-sentence
   deviation). Non-participant of this idea → notice + fall through; participant but ping-failed →
   blocking gate with Confirm exits (re-designate / `implementer_waived: …` / `none`).
4. Today's chain verbatim: `FINAL.md` `{implementer, drafted-by}` → driver tail `eligible[0]` /
   consensus tail `""` → everyone. Both tails byte-preserved under their existing test pins
   (`app_test.go`, `roundgate_test.go`); the shared reader in `internal/protocol` takes the
   facilitator-filtered list as an explicit parameter at both call sites (X-5).

Designation-only additions (fire only when tier 2 or tier 3 is live): the tier-2/tier-3 gates
above; the recorded waiver/`none` handling; Phase-0 relocation of `checkModelDiversity` under
`require_model_diversity: true`; kickoff detection of the zero-invokable-non-implementer-reviewer
floor (existing L538/`driver_impl.go:282` floor, surfaced early — hard gate) and the two-participant
LE-11 consequence (kickoff **warning**, not a block — §5 L746 keeps two-participant rules unchanged);
drafter preference in `firstEligibleHeadlessAgent` (prefer a non-designee drafter, fail-open to the
designee, degenerate case test-pinned, §15.5-style one-line concentration record when drafter ==
designee — **designation-only**, never forbidding the coincidence); the `agent.implementer_resolved`
event + one-line stdout notice (designation-only emission, position change 4); design-only /
`--no-implement` runs: designation **validity** gates always fire (frontmatter defects are worth
stopping on any run), **availability** checks attach only to Phase-5-reaching runs — designation-only
scoping, with the implementation-time trace of `--no-implement` visibility at preflight and the
waiver line as exit if invisible. Claim mechanism: untouched, social-only, subordinate sentence.
Text-only, zero behavior: the L433 sentence adding designated runs to the self-containment trigger
list; identical hunks in all three COOPERATION.md copies + `meta/protocol-changelog.md`; the skill-
tree restatement grep as an acceptance check.

Deferred out of this idea, recorded with slug `meta-protocol-change-implementer-inheritance-repair`:
the legacy FINAL.md read-set widening **and** the write-end emission (same observable unset-path
selection change; both need an owner-visible boundary relaxation, not a participant trade-off);
trajectory `--implementer` cross-checking; `--participants` order documentation; the
`facilitatorConflictGates` all-ideas iteration repair.

## Default-path ledger (owner boundary)

| Mechanism | Default path (global UNSET, no per-idea field) | Designation-only path |
|---|---|---|
| Machine chain | Tiers 1→4 = today's two-file walk + both tails, byte-identical, test-pinned | Tiers 2-3 inserted; same shared reader; X-5 filter parameter explicit |
| Gates | No new gate can fire (all new gates require a present designation; family scoped to the idea being run) | Validity (blank/non-`none` id/facilitator), availability-with-exits, pin-conflict, Phase-0 diversity, zero-reviewer detection |
| FINAL.md key set / draft prompt | Unchanged — write-end emission **excluded**, deferred with slug | Same (designation resolves above the FINAL rank; repair unnecessary for the mechanism) |
| Claim | Unchanged, inert to machines; social act only | Subordination sentence only |
| Drafter selection | Unchanged | Preference for non-designee drafter, fail-open (**designation-only**) |
| Diversity / reviewer floor | Phase-6 gate untouched | Phase-0 relocation + kickoff detection/warning (**designation-only**) |
| Observability | **Nothing new emitted** (position change 4) | `agent.implementer_resolved` + stdout line (**designation-only**) |
| Config / init | All existing keys untouched | `default_implementer` emitted commented-out (documented deviation) |
| Design-only runs | Unchanged | Validity gates only; availability checks never attach (**designation-only**) |

No element of this idea's shipped delta changes selection, gating, or emitted output on a deck that
sets neither field — the one former exception (the event) moved columns this round.

## Role assignment for this run (today's protocol; no global-default weight)

I **accept the assignment: drafter = claude-1, implementer = kimi-1**, reviewers = claude-1 +
zcode-1 (non-implementers; Anthropic + Zhipu against a Moonshot implementation — model-diverse by
construction). I counter zcode-1's claude-drafts/zcode-implements proposal with reasons, not
majority rule:

1. **Evidence base.** All three of us independently said the post-release recommendation needs more
   than the deck's single implementation data point; zcode-1 itself discounts its lean-organizer
   implementation as "n=1 and not a default argument". kimi-1 implementing **produces the missing
   second data point**; zcode-1 implementing deepens the n=1 we already agreed is insufficient.
2. **Rotation against concentration.** This idea's own risk register (my R1, claude-1's R-2,
   zcode-1's risk 3) is standing concentration in one agent. Lean-organizer ran kimi drafts / zcode
   implements / claude reviews; the full rotation is claude drafts / kimi implements / zcode reviews.
   zcode-1's proposal repeats zcode in the implementer seat — the concentration this mechanism
   exists to make deliberate.
3. **Delta-specific risk fit.** This delta's dominant risk is locator fidelity across three protocol
   copies and six `implementer` readers (X-6/X-7/R-C/R-F); the only directly relevant measured
   evidence on the table is X-1/X-2 — all three round-01 artifacts' locators re-measured, kimi-1's
   correct. Small sample, stated as such, but it is the only measurement of exactly the skill this
   implementation most needs.
4. **Reviewer strength preserved.** The deck's one demonstrated high-yield review (the nine-finding
   set traced above) came from claude-1 in the reviewer seat; this assignment keeps claude-1's review
   *and* adds its specification base as drafter, and adds zcode-1's fresh review (its X-5 find shows
   current-form divergence detection).

Perfection acts happen at the protocol-prescribed times, not here: claude-1's `Drafter: yes` in its
consensus signoff block (L409); my `inbox/kimi-1-to-all_meta-protocol-change-designated-implementer_
impl-claim.md` before Phase-5 work begins (L443). **This is a this-run assignment under today's rules
— not a bid for, or recommendation of, the owner's future global default**, which ships UNSET and is
decided post-release on the organizer's evidence report.

## Readiness for consensus

**Ready for consensus, contingent on two zcode-1 concurrences** — the only remaining substantive
threads, both argued above:

1. **Live layered read over Phase-0 materialization** (position change 1; your three-state survives
   via `none`; your write-path concerns — idempotency, dry-run, provenance — simply vanish).
2. **Write-end FINAL repair out**, deferred with the read-set widening to
   `meta-protocol-change-implementer-inheritance-repair` (same observable unset-path selection
   change; the boundary sentence admits owner relaxation, not participant trade-offs).

And one small claude-1 assent: the reserved `none` literal beside S-3's blank-gate. If zcode-1
maintains either position, cross-review 3 of 3 has exactly those two threads and nothing else;
everything else in this file records three-way convergence on verified evidence. FINAL must carry
the single chain above — one precedence ladder, no contradictory restatements — and the default-path
ledger as its compatibility proof.

<!-- Original language of the owner's quoted words and the release-order direction: Slovak;
translations as supplied verbatim in 00-prompt.md, source-context/owner-brief.md and
source-context/release-order.md, where the Slovak originals are preserved unaltered. -->
