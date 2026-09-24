---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
round: 1
date: 2026-09-24
---

## Summary

Yes — the protocol should let the owner designate the implementer, and the right shape is the
pattern the codebase already shipped for the declared facilitator in 1.49.0: an optional
per-idea `implementer:` field in `00-prompt.md` plus an optional global `default_implementer`
in the `[defaults]` block of `~/.parley/agents.toml`, with the global default shipped UNSET so
that an absent field means byte-identical behavior to today. Today there is no designation
mechanism anywhere: the protocol's default is "implementer = FINAL drafter" with a purely
social inbox claim as the only override (COOPERATION.md L443), and the only machine-read
sources are the `implementer`/`drafted-by` frontmatter keys written *after* the fact. I verified
at HEAD that the two unexported `resolveImplementer` copies exist and **disagree** in their
fallback tail (`internal/app/driver_impl.go:104` falls back to `participants[0]`;
`internal/consensus/consensus.go:751` returns `""`), which this idea should unify. I also
verified the brief's hypothesis: with `author: user` the default drafter — and therefore the
default implementer — is indeed decided by round-01 filing order (speed), not suitability,
unless someone volunteers/claims first.

## Proposed approach

### Answer to the design question

Adopt the mechanism. The owner's words ask for exactly the symmetric counterpart of the
organizer role: one agent executes FINAL (code or other work), the others evaluate. The
declared-facilitator feature proves the pattern works end-to-end in this codebase: optional
frontmatter field, absent-field byte-identity, fail-closed preflight validation, driver-side
role eligibility, escalation instead of silent fallback
(`internal/protocol/facilitator.go:15-22`, `internal/app/preflight.go:314-330`,
`internal/app/driver_impl.go:46-62`, CHANGELOG 1.49.0). We should not invent a second pattern;
we should instantiate the existing one for the implementer role.

### 1. Role scope (code and non-code work)

The designated implementer owns **execution of FINAL in Phases 5–8**: creating the
implementation branch, writing/updating `IMPLEMENTATION.md` (plan, progress, decision log,
validation evidence), performing whatever work FINAL prescribes — code, docs, protocol-hunk
application, release preparation, other non-code artifacts — and applying fix-up cycles.

The role explicitly does **not** confer or remove anything else:

- No drafting rights/duties change (Phase 4 stays: `author:` → `author: user` first-filer
  default → `Drafter: yes` volunteer, L405-409).
- No design participation change: the designated implementer remains a full participant in
  rounds 1..N and owes their round files like everyone else.
- No signoff weight change: they still sign consensus.md and review/consensus.md (Phase 7
  already includes the implementer's signoff).
- No review rights over their own work: §6's prohibition stands (reviewers are non-implementers,
  L234; no merge without an invokable non-implementer reviewer, L538; goal-done check stays a
  fresh non-implementer, L688-691).

This mirrors the owner's framing ("does the work, the others evaluate it") without creating a
second privileged role.

### 2. Fields, precedence, and the UNSET global default

Two optional declaration surfaces:

- **Per-idea:** `implementer: <agent-id>` in `00-prompt.md` frontmatter (symmetric with
  `facilitator:`). Absent = today's behavior, mirroring the facilitator contract that a deck
  not declaring the field stays untouched (`facilitator.go:16-18`).
- **Global:** `default_implementer` in the `[defaults]` block of `~/.parley/agents.toml` — the
  established home for owner policy defaults (`ping_tier`, `preferred_transport`,
  `roster_change_policy`, `speed`/`timeouts`; COOPERATION.md §0, the `[defaults]` paragraph
  around L57-59). **Shipped UNSET:** the release writes no such key anywhere and documents that
  absence = today's behavior. Choosing the owner's persistent default is expressly not this
  idea's act; after release the organizer reports a recommendation with evidence and the owner
  sets it.

Resolution precedence for "who executes Phase 5" (each tier only if the ones above yield
nothing):

1. `IMPLEMENTATION.md` `implementer` — the re-entry pin. Once work has started, the durable
   record wins; a mid-run edit of `00-prompt.md` cannot silently switch implementers. (Both
   existing resolvers already rank this first: `driver_impl.go:116-119`, `consensus.go:763-766`.)
2. Per-idea `implementer:` designation (validated against the idea's participants).
3. Global `default_implementer` — applies only if it names a participant of *this* idea.
4. Today's rules unchanged: the FINAL drafter is the default implementer, preempted by an
   `inbox/<from>-to-all_<slug>_impl-claim.md` claim posted before work begins (L443).

Rationale: specific intent beats standing preference beats social default; the re-entry pin
beats everything so a run never changes horses mid-flight. If a global default names X while
Y has claimed via inbox for a specific idea, X implements (standing owner preference beats an
ad-hoc claim); the per-idea field remains the escape hatch. I flag this ordering for
cross-review because it is the one place designation visibly overrides an existing social
mechanism.

### 3. Designated agent unavailable

Distinguish a *safety* property from a *preference*. The facilitator machinery escalates rather
than falls back because organizer purity is safety. Implementer designation is an owner
**preference** about who does the work; when it cannot be satisfied, the correct behavior is a
**loud fall-back to today's mechanism**, not a silent switch and not a hard stall:

- `implementer:` names an id not in `participants:` → **preflight failure**, fail-closed,
  naming both fields (a typo'd or stale designation must never silently resolve to
  `participants[0]` — that is the current driver tail and precisely the failure mode to close).
- `implementer:` names the declared non-participating facilitator → **preflight failure**
  (same shape; the driver's `IneligibleForRoles` would already exclude it from `eligible`, so
  the gate turns a silent fall-through into a kickoff error).
- The designated agent is a valid participant but unreachable at the §9.0 readiness ping →
  record it in the readiness table and an inbox note, then fall to tier 4 (claim → drafter).
  Quorum is NOT shrunk: the agent remains a participant who owes round/review files under the
  §5 async rules; only execution reverts. An owner who prefers the idea to wait instead can
  fix the designation and reopen — the mechanism documents the choice rather than encoding a
  wait.
- Global default names an agent not on this idea's participants → **not an error** (rosters
  differ across decks); it is simply inapplicable: record and fall through. The global value is
  validated at set time against the machine roster (known id, `active = true`).

### 4. Drafter separation and self-contained FINAL

Designation makes implementer ≠ drafter the *expected* case, which promotes the existing
"FINAL + IMPLEMENTATION.md must be self-contained enough for a fresh agent to implement or
resume from them alone" rule (L421-423) from a nice property to the load-bearing contract. No
textual weakening is acceptable; if anything, this idea should add one sentence to Phase 5
making explicit that under a designation the FINAL self-containment requirement is what makes
the designation executable. The protocol need not *force* drafter/implementer separation — the
owner may designate the drafter — but the owner's stated motivation (others evaluate) argues
the recommendation phase should present evidence for a default implementer distinct from the
typical drafter.

### 5. Reviewer model diversity

Reviewer selection is untouched: all non-implementers review on `deliberation` (L234), the
non-implementer-merge gate stands (L538), the goal-done checker stays a fresh non-implementer
(L688-691), and LE-3/`require_model_diversity` continues to gate all-shared-model reviewer sets
(L536). With a three-agent roster and one implementer, the reviewers are the other two agents,
so designation does not by itself create homogeneity. Two standing-default risks remain and
should be measured in the post-release recommendation: (a) if two roster agents share a model,
a designated third could leave an all-one-model reviewer pair — caught only when
`require_model_diversity: true` (this idea sets it; ordinary ideas may not); (b) standing
concentration — one agent accumulates all implementation context and never reviews others'
code. Mitigations: the cheap per-idea override, keeping the implementer a full design
participant, and evidence in the recommendation. The lean-organizer run is relevant precedent:
with three participants and zcode-1 implementing, the reviewers were claude-1 and kimi-1 —
model-diverse by construction.

### 6. The claim mechanism stays as fallback

Yes — tiers 3/4 above. The inbox claim remains the only override when nothing is designated
(the shipped default world), and the escape hatch when a designated agent declines (voluntary
inbox note → tier 4 opens). Under a live designation a non-designated participant's claim does
not override the designation; it may record willingness for the fallback.

### 7. Driver, preflight, and fail-closed validation

- **Preflight** gains a validation family mirroring `facilitatorConflictGates`
  (`preflight.go:314-330`): designation-not-a-participant → gate; designation = declared
  non-participating facilitator → gate; (set-time) global default must name an active machine
  roster id. All fail closed with both fields named.
- **Driver**: `newDriverImplOps` resolves the designation at the top of the existing chain,
  validated against the facilitator-filtered `eligible` list; a designation that names the
  ineligible facilitator surfaces through the existing `roleErr` escalation path rather than a
  silent `eligible[0]`.
- **Skill/CLI surfaces**: the Phase 5 prose in all three COOPERATION.md copies (deck view,
  `internal/protocol/defaults/COOPERATION.md`, the skill's bundled
  `skills/parley-deck/references/COOPERATION.md`), the Phase 0 frontmatter sketch, the §0
  `[defaults]` paragraph, the TL;DR item 6 sentence ("The FINAL drafter is the default
  implementer" gains "; an explicit `implementer:` designation or the owner's global default,
  when set, takes precedence"), a §9.0 preflight note, and the protocol-changelog entry.

### 8. Consistency of the two resolveImplementer copies

Verified at HEAD — they exist and **do not agree**:

- `internal/app/driver_impl.go:104-131`: IMPLEMENTATION.md `implementer` → FINAL.md
  `implementer`/`drafted-by` (validated) → **`participants[0]`** → `""`.
- `internal/consensus/consensus.go:751-774`: the identical chain but **no `participants[0]`
  fallback** — returns `""`, which `expectedRoundParticipants` (L726-748) turns into "require
  the full participant list", i.e. it demands a review file from the implementer too — an
  artifact §6 forbids the implementer from writing (the code comment at L738-741 admits the
  over-fail).
- A second divergence axis: the driver copy validates against the facilitator-filtered
  `eligible` list (`driver_impl.go:46-62`) while the consensus copy validates against raw
  `participants`, so in a declared-facilitator run the two can resolve different implementers
  from the same FINAL.md.

The two copies agree in practice only because the tool-written FINAL template emits
`drafted-by` (`consensus.go:805`). This idea must **unify them into one shared, exported
resolver** (natural home: `internal/protocol` beside `facilitator.go`, or a small shared
package both import) with an explicit two-result API — `(id string, ok bool)` — over one
identical chain (re-entry pin → per-idea designation → global default → FINAL metadata), so
each caller keeps its own *tail* policy visibly at one site: the driver picks `participants[0]`
(it must launch someone), `expectedRoundParticipants` keeps its fail-closed full-list guard for
the genuinely unresolvable case. Tests must pin both tails and the facilitator-filtering axis
so unification changes no Phase 6 gating silently — the lean-organizer audit's codex-1/F2
finding was exactly this class of bookkeeping drift.

### 9. What ships UNSET, and what this idea does not do

The global default ships UNSET (no `[defaults]` key written by the release; absence documented
as today's behavior). No agent is selected as the owner's persistent default here; the
post-release recommendation with evidence is the organizer's done-report duty. Non-goals
reaffirmed: no quorum/roster change, no model/effort downgrade, no organizer-rule change beyond
this role, no change to drafter-selection rules (the driver's separate
`firstEligibleHeadlessAgent` participants-order drafter pick at `driver_consensus.go:85` is the
same order-over-suitability critique but is out of scope), no transport-header change (this
run's local-canonical/direct-main override is per-run, not global).

### Current-protocol voluntary claim

Under the current (unamended) protocol only — and explicitly **not** as a recommendation for
the owner's future global default — I am willing to serve as drafter or implementer for this
idea if the round reaches that point under today's rules (claim/first-filer/drafter-default as
they apply). This claim has no weight in the global-default question.

## Concerns / open questions

1. **Field naming.** `implementer:` in `00-prompt.md` (designation) collides in name with
   `IMPLEMENTATION.md`'s `implementer:` (record). Same word, two files, two lifecycle stages; I
   lean to keeping the symmetric `implementer:` (matching `facilitator:`), but
   `designated_implementer:` is the unambiguous alternative. Cross-review should settle it
   before Phase 5.
2. **Mid-idea designation change.** I propose: edits before work starts are ordinary Phase 0/4
   metadata edits; after `IMPLEMENTATION.md` exists, the re-entry pin outranks and a change
   requires a recorded decision (deviation-style log entry). Is post-FINAL/pre-work change
   without a new consensus acceptable? I believe yes (operational metadata, FINAL frozen), but
   it deserves an explicit sentence in the amended text.
3. **Global-default-vs-claim precedence** (§2 above): standing preference beating an ad-hoc
   claim is my recommendation; others may prefer per-idea claims to always win.
4. **Unavailable-designated fallback vs wait.** I chose loud-fallback over hard-block on speed
   grounds (the owner asked for speed; the fallback is exactly today's protocol). If the owner
   would rather designated ideas wait for the designated agent, that is a one-line variant —
   worth one owner question at most, or default to fallback and document.
5. **Skill-bundle surface.** Beyond the three COOPERATION.md copies, SKILL.md may restate the
   Phase 5 default drafter-implementer rule; the implementation must grep the skill tree for
   restatements so no copy drifts (the drift guard covers COOPERATION.md copies, not prose
   restatements).
6. **Windows scope.** The prior release left an unresolved Windows scope decision; if its
   resolution lands first, this idea integrates on top of it — sequencing belongs to the
   organizer.

## Risks

1. **Two sources of truth** (00-prompt designation vs IMPLEMENTATION.md record) — mitigated by
   the explicit precedence ladder and the re-entry pin; the resolver must be the single point
   that reads both.
2. **Unification regressions in review bookkeeping.** Touching `expectedRoundParticipants`
   can change Phase 6 completeness checks; requires tests pinning both tails and the
   facilitator-filter axis before the merge (precedent: codex-1/F2 audit finding).
3. **Standing concentration under a global default** — one model accumulates implementation
   context, never reviews; mitigated by per-idea override, full design participation, and the
   evidence-backed post-release recommendation, but it is a real long-term diversity cost the
   owner should weigh when setting the default.
4. **Scope creep into organizer rules** — the delta must stay additive and absent-field
   byte-identical (the facilitator contract); any hunk that changes behavior without the new
   fields present is a defect.
5. **Three-copy + release-channel consistency** — protocol hunks must land identically in all
   three COOPERATION.md copies; release sequencing (prior 1.49.x run's incomplete channels,
   core 2.13.0 publish state) is the organizer's boundary, and a stale base would fork the
   copies.

## Existing alternatives

Per §15.6(a) — what the proposal would build by hand, and what the toolchain already ships:

1. **Implementer selection today (the mechanism being extended):** protocol text at
   `parley-deck/COOPERATION.md:443` — default implementer is the FINAL drafter, preempted by an
   `inbox/<from>-to-all_<slug>_impl-claim.md` claim before work begins. Toolchain support: the
   impl-claim file has **no** code support (grep `impl-claim` across `internal/` → only the
   COOPERATION.md prose); the machine-resolvable sources are only the after-the-fact
   frontmatter keys `implementer` (IMPLEMENTATION.md) and `implementer`/`drafted-by`
   (FINAL.md) read at `internal/app/driver_impl.go:116-119` and
   `internal/consensus/consensus.go:763-766`.
2. **`roles:` advisory map:** `00-prompt.md` sketch and semantics at
   `parley-deck/COOPERATION.md:288-290` and `:310-311` — an owner could write
   `roles: {<agent>: implementer}` today, but it is advisory only, is parsed by no code (grep
   `"roles"` across `internal/` → no hits), and the protocol states roles do not change quorum,
   signoff weight, artifact ownership, drafter eligibility, or roster membership. This is the
   no-code-change alternative; it fails the owner's requirement because it selects nothing.
3. **Declared-facilitator machinery (the pattern being mirrored):**
   `internal/protocol/facilitator.go:15-60` (optional field; absent = byte-identical v1.48
   behavior; `Conflict()` fail-closed message), preflight gating at
   `internal/app/preflight.go:314-330` (`facilitatorConflictGates`), driver eligibility
   filtering and escalation at `internal/app/driver_impl.go:46-62`, shipped in 1.49.0
   (CHANGELOG.md "Declared-facilitator runs keep the organizer pure…").
4. **The two existing resolvers (the code being unified):**
   `internal/app/driver_impl.go:104-131` (tail `participants[0]`) and
   `internal/consensus/consensus.go:751-774` (tail `""`), consumers at `driver_impl.go:64` and
   `consensus.go:734/726-748`.
5. **Global policy-defaults home:** `~/.parley/agents.toml` `[defaults]` block — `ping_tier`,
   `preferred_transport`, `roster_change_policy`, `speed`/`timeouts` — described at
   `parley-deck/COOPERATION.md` §0 (the `[defaults]` paragraph, ~L57-59); a global
   `default_implementer` is a new key in an existing, tooled home rather than a new mechanism.

Scoped null note: beyond the above, I consulted `internal/protocol/` (file listing),
`internal/app/driver_consensus.go`, `parley-deck/meta/protocol-changelog.md`, and `CHANGELOG.md`
for adjacent mechanisms; I found no other shipped implementer-selection or role-designation
machinery (no `Drafter: yes` parsing, no first-filer computation anywhere in `internal/`).

## Verification record (§15)

### Protocol packet attestation

Obtained via `parley protocol packet --dir . --phase 1 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`; body read in full
(1,386 lines). `context_mode=full`, `source_sha256=8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7`,
`packet_sha256=8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7` (equal in full
mode; source = live `parley-deck/COOPERATION.md`, protocolRole `source`, 109,928 bytes),
`fallback_reason` absent — matches the kickoff attestation recorded in `00-prompt.md`. Working
tree: HEAD `59b0458` (organizer's idea-open commit) on origin/main `9134c7a`; `git diff
origin/main HEAD -- parley-deck/COOPERATION.md internal/` is empty, so the live protocol file
and code equal origin/main. Transport: `github-pr` (per-run local-canonical override noted in
`00-prompt.md`).

### Brief-locator verdicts (all PRIMARY — file read at the cited lines, this worktree, HEAD)

| Brief claim | Verdict | Evidence (locator) |
|---|---|---|
| Phase 4 (~L407-409): drafter = `author:`; `author: user` → first round-01 filer; volunteer via `Drafter: yes` | CONFIRMED | `parley-deck/COOPERATION.md:407-409` |
| Phase 5 (~L443): default implementer = FINAL drafter; claim via `inbox/<from>-to-all_<slug>_impl-claim.md`; declared facilitator never implements unless `facilitator_participates: true` | CONFIRMED | `parley-deck/COOPERATION.md:443` |
| `facilitator_participates` also ~L883 | CONFIRMED | `parley-deck/COOPERATION.md:883` |
| `roles:` (~L311) advisory lenses; no quorum/signoff/ownership/drafter-eligibility/roster change | CONFIRMED | `parley-deck/COOPERATION.md:310-311`; no code parses `roles` (grep `internal/`) |
| Reviewers all non-implementers (~L234, deliberation) | CONFIRMED | `parley-deck/COOPERATION.md:234` |
| No merge without an invokable non-implementer reviewer (~L538) | CONFIRMED | `parley-deck/COOPERATION.md:538` |
| Goal-done check by a fresh non-implementer (~L689) | CONFIRMED | `parley-deck/COOPERATION.md:688-691` |
| Two separate unexported `resolveImplementer` at `driver_impl.go:104` and `consensus.go:751` | CONFIRMED | `internal/app/driver_impl.go:104`, `internal/consensus/consensus.go:751` |
| "Check whether they agree" | **WRONG (they do not agree)** — same frontmatter chain, but divergent tails (`participants[0]` vs `""`) and divergent validation sets (`eligible` vs raw `participants`) | `driver_impl.go:129-131` vs `consensus.go:771-773`; `driver_impl.go:46-62,64` vs `consensus.go:734` |
| `facilitator:`/`facilitator_participates:` shipped in 1.49.0 | CONFIRMED (version content; channel completeness not re-verified here — prior run's handoff reports incomplete channels) | `VERSION`=1.49.0; `CHANGELOG.md` 1.49.0 entry; `internal/protocol/facilitator.go`; `parley-deck/meta/protocol-changelog.md` 2026-09-23 entry |
| Lean-organizer run: kimi-1 drafted FINAL, zcode-1 implemented | CONFIRMED (as outcomes; the *claim acts* themselves unverified) | `parley-deck/meta/protocol-changelog.md` ("Drafted by: kimi-1 (FINAL); protocol hunks applied by zcode-1 (Phase 5)"); lean-organizer worktree `IMPLEMENTATION.md` frontmatter `implementer: zcode-1` |
| claude-1's first code review: 1 CRITICAL + 8 MAJOR | CONFIRMED | lean-organizer worktree `review/round-01/claude-1.md` — grep counts: 1× `### [CRITICAL]`, 8× `### [MAJOR]` |
| "all fixed in later cycles" | UNVERIFIED (supported, not traced: 3 fix-up cycles + review consensus cycles exist and the run closed `status: complete`; I did not trace each of the 9 findings to its fix) | same `IMPLEMENTATION.md` (`## Fix-up cycle 1/2/3`, `status: complete`) |
| Cost: implementation ≈ 30-40% of organizer input when the organizer implemented | CONFIRMED as quoted, with the source's own caveat: single session, 98.1% pre-v1.48.0, organizer was simultaneously author/drafter/implementer/release manager | organizer-token-study `2026-09-23/README.md` §2 item 1 (snippet copied below) |
| Hypothesis: with `author: user`, default drafter (hence default implementer) decided by round-01 filing speed, not suitability | CONFIRMED as the social default chain, with two qualifications: (a) a volunteer/claim preempts speed in both phases, so speed decides only when nobody steps up; (b) mechanically the driver never computes "first filer" — it reads FINAL.md `drafted-by` when present, else falls to `participants[0]` order, and picks the Phase 3/4 drafter as first *eligible participant* in list order, so the driver's own defaults are order-based too | `COOPERATION.md:407-409` + `:443`; `driver_impl.go:116-131`; `driver_consensus.go:85` (`firstEligibleHeadlessAgent`); no first-filer code anywhere (grep) |

### Snippets copied from out-of-deck sources (§6 rule 4)

From `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/README.md`
(outside this deck):

> Caveat: 98.1% of this session predates the v1.48.0 commit. It measures an organizer that was
> also author, FINAL drafter, implementer and release manager of evidence-first-efficiency, not
> a pure facilitator.
>
> 1. **Self-implementation and self-verification: ~30-40% of input** (49.8% is the upper bound
> that counts all code-level work).

From the lean-organizer worktree deck (different worktree, read-only):

> `parley-deck/meta/protocol-changelog.md`: "Drafted by: kimi-1 (FINAL); protocol hunks applied
> by zcode-1 (Phase 5)"
>
> `IMPLEMENTATION.md` frontmatter: `implementer: zcode-1`, `status: complete`,
> `fix-up-cycle: 3`

Owner Slovak quotations and translations are preserved verbatim in `00-prompt.md` and
`source-context/owner-brief.md`; I rely on those copies and add nothing to them.

No other participant's round-01 file was read before this artifact was saved (round-01/ was
empty at write time). No source code was changed; no commit was made.
