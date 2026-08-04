---
agent: kimi-1
idea: meta-protocol-change-verification-integrity
round: 1
date: 2026-08-04
---

## Summary

The brief names a real structural gap: this protocol regulates who writes artifacts and
never defines what makes a verification valid. I judge that gap genuine independently of
the unverifiable source run — the protocol contains no verdict vocabulary at all, so any
`CONFIRMED` stamp anywhere is currently free-form prose. Of the nine proposals, I would
adopt four with amendments (CRITICAL-1, CRITICAL-2, CRITICAL-3, MAJOR-5), adopt two
essentially as written (MAJOR-4, MINOR-9 as a clarification), scope-limit two
(MAJOR-6 partially already covered, MAJOR-7 deliberation-only), and adopt MINOR-8 as an
honesty fix to an overclaim the protocol actually makes today. Of the six tooling defects
I verified four by execution (T1, T2, T4-structure, T5), one by source inspection plus a
non-destructive dry-run that fails to reproduce it on CLI 1.37.0 (T3), and one by reading
the skill text (T6). I deliberately did not run a live-ping `parley preflight` and did not
run a writing `parley roster init` against the live deck.

The brief's observations are testimony from a run not in this deck. Where a rule only makes
sense if you believe the anecdote, I say so below (MAJOR-5 as written, MAJOR-6(a) as written).

## Proposed approach

Positions per proposal, in brief order. "Verdict" below means the vocabulary these rules
would introduce (`CONFIRMED` / `UNVERIFIED` / `WRONG` / `UNKNOWN` / `DISPUTED`); none of it
exists in `COOPERATION.md` today, which is itself the finding.

### CRITICAL-1 (self-verdicts) — ADOPT AMENDED

The asymmetric core is sound, anecdote-independent, and checkable by reading author lines:
an owner's confirmation of its own claim adds zero evidence (same correlated prior that
made the claim), while an owner's retraction is high-value. But the proposal needs a
two-participant degenerate rule and must not forbid the owner from producing evidence —
only from *verdicting* it. Replacement text:

> A verification verdict on a claim is admissible only from a participant other than the
> claim's author. The author may restate, supply evidence for (commands run, sources
> quoted), or weaken its own claim at any time; an author entry that changes its own
> claim's status is a `SELF-CORRECTION`, valid only when it weakens
> (`CONFIRMED → UNVERIFIED`, `UNVERIFIED → WRONG`, or withdrawal) — never a strengthening
> and never a `CONFIRMED`. A claim counts as verified for consensus only with at least
> one non-author verdict at `SECONDARY` provenance or better (CRITICAL-2). With two
> participants this reduces to: the verifying verdict must come from the other
> participant. On `fast` this is one line in the reviewer's file, not an extra round.

Checkable: any `CONFIRMED` whose author equals the claim owner is invalid on its face;
any strengthening `SELF-CORRECTION` is invalid on its face. No new tooling.

### CRITICAL-2 (provenance tags) — ADOPT AMENDED

Sound and the cheapest of the three criticals. Two defects as written: (a) `PRIMARY` is
literature-centric ("venue/DOI/identifier") — in a software deck the primary source is
usually an executed check, and the rule must admit it or every code claim caps at
`SECONDARY`; (b) no default for an untagged verdict, so the scheme fails open. Replacement
text:

> Every verdict carries exactly one provenance tag:
>
> | Tag | Meaning | Admissible for `CONFIRMED`? |
> |---|---|---|
> | `PRIMARY` | Source located and quoted (venue/DOI/identifier), **or a check actually executed, with the command and its output quoted** | Yes |
> | `SECONDARY` | Independent confirmation by ≥1 other participant, itself not `RECALL` | Yes |
> | `RECALL` | Model memory only, no source consulted | No — caps the verdict at `UNVERIFIED` |
>
> **A verdict with no tag is treated as `RECALL`.** A claim reaching consensus with only
> `RECALL` support is recorded as unverified in `FINAL.md`.

The tag is self-reported and can be lied about — but `PRIMARY` as amended is falsifiable
(the quoted locator or command output can itself be checked), which is what makes it a
rule rather than a comment. Binds on every track; it is one word per verdict.

### CRITICAL-3 (verdict-conflict register) — ADOPT AMENDED

The gap is real and current: Phase 3's "agent silent past deadline is treated as ✅" can
launder an unresolved 2–2 split into adoption, and no `DISPUTED` terminal state exists.
Forbidding vote-counting is correct and consistent with MAJOR-7's correlated-prior
argument. But a mandatory `verdicts.md` per idea is ceremony without signal — most ideas
have no verdict conflicts. Replacement text:

> On the first verdict conflict in an idea, the consensus drafter creates
> `ideas/<slug>/verdicts.md` (drafter-owned, append-only — same ownership shape as
> `consensus.md`; absent when no conflict exists). Each entry: the claim, the conflicting
> verdicts with authors and provenance tags, and the resolution. A conflict MUST be
> resolved before consensus, by argument or provenance, never by counting agents: the
> higher provenance tag wins; at equal provenance, the verdict carrying an explicit
> derivation or source locator wins; otherwise the claim enters `FINAL.md` as `DISPUTED`
> under a mandatory heading, and a `DISPUTED` claim may not be cited as support for any
> acceptance criterion. Counting agents is forbidden as a resolution method. In the
> review phase, disputed findings follow the ratified P6 close rule (reviewer withdrawal,
> review consensus, or verbatim-quoted operator ruling) instead of this register.

That last sentence is the composition with P6 the facilitator note demands: P6 governs
review-phase disputes, this governs design-phase verdicts; no overlap, no contradiction.

### MAJOR-4 (obstacle-claim admissibility) — ADOPT AS WRITTEN

The witness taxonomy (counterexample / explicit precondition check / cited result) is
narrowly scoped to claims about a *known, named* obstacle, which is exactly where
unfalsifiable prose thrives, and "adjectives are not witnesses" is enforceable by any
reviewer asking "where is the witness?" It is strictly a special case of CRITICAL-2
applied to exemption claims, but the special case earns its text because it names the
failure. Sound without the anecdote: unfalsifiable exemption claims need no p-vs-np run
to be a hazard. Binds on all tracks — the cost when no obstacle is named is zero.

### MAJOR-5 (settledness checks) — ADOPT AMENDED (narrowed)

**As written, this rule only makes sense if you believe the anecdote.** "Any proposed
sub-goal, milestone, or acceptance criterion MUST carry a settledness check" taxes every
engineering milestone ("add a retry policy") with a question that has no referent —
acceptance criteria are things we intend to *make* true, not claims about the world.
That is ceremony without signal, and on this deck's actual idea mix it would degrade to
boilerplate `open` stamps. The sound, checkable core is novelty/openness claims.
Replacement text:

> A claim that a problem is open, a result novel, or an approach previously untried MUST
> carry provenance per CRITICAL-2. A novelty or openness claim with `RECALL`-only
> support is recorded `NOVELTY UNVERIFIED` in `consensus.md` and may not be presented as
> recommended work in `FINAL.md`. Milestones and acceptance criteria — states the idea
> intends to bring about — require no settledness check.

### MAJOR-6 (facilitator conflict of interest) — (a) ALREADY COVERED + one-line amendment; (b) ADOPT AMENDED, standard+deliberation; (c) fold into (b)

(a) **Already covered — as written it legislates for a power that does not exist here.**
`COOPERATION.md` gives the facilitator no adjudication power: Phase 3 consensus is
all-participant signoff, any ❌ forces a new round, and unresolvable disputes escalate
to the user (§4 "Escalation to user"). "Facilitator dispute rulings are PROVISIONAL until
ratified" is a rule that only makes sense in the source run, where the facilitator was
"dispute adjudicator" — here there is nothing to ratify because there is nothing to rule.
The one genuine gap is *procedural* calls (declaring a round closed, judging convergence),
so adopt one sentence:

> The facilitator's procedural calls — declaring discussion converged, opening consensus,
> closing a round — are provisional until the corresponding signoff gate passes; the
> signoffs, not the facilitator's judgment, are the close.

Checkable, cheap, and removes the de-facto ruling power without inventing a formal one.

(b) The concentration risk is structural, not anecdote-dependent: the same agent
summarizing others' positions and then drafting the authoritative artifact can smooth
conflicts. Existing partial mitigations: append-only signoffs, and Phase 3's rule that
any participant may block an inaccurate `## Comparison & blind spots`. Adopt amended:

> On `standard` and `deliberation`, when the facilitator is also a participant, the
> `consensus.md` drafter SHOULD be a different agent. Where the drafter rule forces the
> same agent, `FINAL.md` MUST record the role concentration in one line, and the
> drafter MUST list its own position changes since its last round file in
> `consensus.md` — traceable by any participant against the raw round files, which are
> never hidden (Phase 3).

This replaces the uncheckable "a non-drafter MUST review the drafter's concessions"
(review *for what*? a concession is a position change, so make it traceable instead) with
something a participant can actually verify by diffing files.

(c) Advisory role-separability is a comment, not a rule — fold the sentiment into (b)'s
SHOULD and drop the standalone clause.

### MAJOR-7 (agreement treated as convergence) — ADOPT AMENDED (deliberation only)

The correlated-prior argument is true of LLM panels independent of any anecdote — shared
training data means unanimity is not four independent samples. Assigned steelmanning is
the standard groupthink countermeasure and its artifact is trivially checkable (the file
exists or it doesn't). But binding it wherever round 1 is unanimous would force dissent
theater on routine engineering ideas, where unanimity is common and cheap convergence is
correct; `standard` ideas also have Phases 6–8 to catch error mechanically. The
constraint permits binding on `standard` and `deliberation`; I would bind only
`deliberation`, with a mechanical trigger:

> On the `deliberation` track: if round 1 closes with no substantive disagreement between
> participants, consensus MUST NOT close until (a) the facilitator assigns one
> participant to steelman the strongest rejected alternative, filed as a canonical
> round-02 artifact — a signoff MAY block a pro-forma steelman as an inaccurate
> `## Comparison & blind spots` under the existing Phase 3 rule — and (b) `consensus.md`
> records a correlated-agreement caveat: unanimity among related models is a shared
> prior, not independent evidence. `FINAL.md` MUST state where multiple "independent"
> proposals are in fact one family.

Track is already mechanical (§4.0), so the trigger needs no judgment call. Composes with
P7: P7 judges review *trajectory*; this judges round-1 *convergence* — different gates,
no conflict.

### MINOR-8 (round-1 independence) — ADOPT AMENDED (fix a real overclaim)

The demanded honesty is half-present: §11.A already says "There is no enforcement beyond
agent discipline." But §4.0 lists "round-1 independence (Phase 1)" among **invariants
"never dropped for speed"** — that *is* claiming an unenforced guarantee, so the proposal
has a live target in current text. Adopt amended:

> Phase 1 and the §4.0 invariant state plainly: round-1 independence is a cooperative
> convention, ex-post auditable through commit order and timestamps (transports A and B),
> not an enforced property. Where independence is load-bearing (blind review,
> benchmarking), the idea MUST use §11.B sub-branches or per-agent staging
> (`parley-worktrees`), and `00-prompt.md` MUST say so at kickoff.

Cite as partial existing coverage: §11.A Phase 1 paragraph, §11.B "Independence in
Round 1" sub-branch protocol.

### MINOR-9 (facilitator context asymmetry) — ADOPT AMENDED (clarification of §6 rule 4)

Substantially covered already: §6 rule 4 requires copying external snippets into the
deck. What is missing is timing and scope — §6.4 reads as applying to references made
mid-discussion, not to material gathered while scoping. Adopt as a one-line amendment to
§6 rule 4 (or Phase 0):

> §6 rule 4 applies to scoping: any source material the facilitator gathered while
> scoping an idea MUST be copied into `00-prompt.md` (or a sibling file referenced from
> it) before participants are invoked.

Checkable by reading `00-prompt.md`; cost is trivial; cites §6 rule 4 as existing coverage.

### Tooling defects — what I ran, what I did not

- **T1 — VERIFIED by reproduction.** In a disposable deck (`mktemp -d`, `git init`,
  `parley init -dir <tmp>`; removed afterwards): the generated `COOPERATION.md` §2 roster
  table is empty, and `parley roster show` fails with exactly
  `roster show: could not read the §2 roster (COOPERATION.md)`. Endorse the fix: seed §2
  from `~/.parley/agents.toml` at init (init already reads that file for defaults per §0),
  or have init emit a fillable roster scaffold. A freshly initialized deck being
  non-functional is a real defect, not testimony.
- **T2 — VERIFIED by dry-run and source.** `parley-deck-skill sync-project --project .
  --dry-run --json` on the live deck shows the rewritten `meta/version.json` **omits
  `protocolRole: source`** (present in the on-disk metadata). `internal/app/preflight.go:409-412`
  raises the `unknown-role` gate when `protocolRole` is absent. The round-trip data loss
  and its presentation as a deck problem are both confirmed. Endorse "preserve
  unknown/extra fields on sync".
- **T3 — DID NOT REPRODUCE on CLI 1.37.0; not verified as described.** I ran only
  `parley roster init --scope session --dry-run` against the live deck (no writing run,
  per the brief). It proposes *appending* `[roster.<id>]` blocks for all four IDs —
  including `kimi-1`, whose `AUTO=no` — mapping it to `kimi`; nothing is dropped and no
  retired adapter is re-added. Source (`internal/app/roster.go:259-274`) only ever appends
  missing mappings and fail-closes (exit 1) on genuinely unresolvable families; the
  `unmapped — run parley roster init` hint (roster.go:126) prints only when a family
  cannot be resolved, and did not print on the live deck. The fail-close itself is
  arguably correct behavior. The residual valid core: where an unmapped entry is
  *intentional*, the hint is still a standing invitation to a command the operator may not
  want — endorse "suppress the hint when the unmapped entry is intentional" as a MINOR
  tooling fix. The "silently drop" half appears already fixed (or version-specific
  testimony); I will not endorse it as a live defect.
- **T4 — VERIFIED STRUCTURALLY; live-ping half deliberately skipped.** I ran
  `parley preflight -no-ping -json` (presence-only; I chose not to ping live agents —
  skipping is my stated choice, and the `-no-ping` flag exists for exactly this). The
  report keys agents by adapter family (`codex`, `claude`, `agy`, `hermes`, `kimi`), not
  by deck roster IDs (`claude-1`, …), and includes `agy`, which is **not** in this deck's
  §2 roster. "Reports by adapter family; touches adapters not in the deck" is confirmed.
  The `unavailable:no-pong` false-negative half is unverified — that required the live
  ping I skipped. Endorse reporting by roster ID and skipping non-rostered adapters.
- **T5 — VERIFIED by observation and source.** Live `parley roster show`: DISPLAY-NAME
  `claude_opus-4.8-1m_max` while MODEL is `claude-opus-5[1m]`.
  `internal/agents/naming.go:188-206` derives the display name from `ModelLabel` when set,
  so a stale label produces exactly this mismatch, and a "fix the model to match the
  display name" reaction would downgrade the configured model. Endorse deriving the
  display name from the resolved model.
- **T6 — VERIFIED by reading the skill.** `SKILL.md:288` and `:387` mandate 30-minute
  per-agent process timeouts; the word "background" appears nowhere in `SKILL.md`, so the
  background-launch pattern is undocumented. (My own host harness caps foreground shell
  calls at 5 minutes — tighter than the 10 in T6 — so the trap is real for any host with
  a foreground cap.) Endorse documenting the background-launch pattern explicitly in the
  Timeout Policy.

## Concerns / open questions

- These rules introduce a verdict vocabulary the protocol lacks; the amendment must say
  *where* verdicts live (round files, `consensus.md`, the conditional `verdicts.md`) or
  CRITICAL-1..3 are unenforceable conventions. I have assumed verdicts are recorded in
  the authoring agent's own round/review files — ownership unchanged.
- Provenance tags are self-reported. `PRIMARY`-as-executed-check is falsifiable, but
  `SECONDARY` chains and `RECALL` honesty rest on agent discipline like the rest of the
  protocol. The tags buy auditability, not proof; the amendment should not oversell them.
- CRITICAL-3's "higher provenance wins" can still adopt a false claim when both sides
  quote genuine-looking locators (a fabricated `PRIMARY` beats an honest `SECONDARY`).
  The `DISPUTED` escape and the P6 operator-ruling path mitigate; they do not eliminate.
- MAJOR-7's steelman can degenerate into compliance theater. The only check I see is the
  existing "block an inaccurate comparison" signoff rule; if participants won't use it,
  the rule adds a round and buys nothing.
- Scope of "claim": does a verdict regime apply to every assertion in a round file, or
  only to claims a verdict is actually issued on? I assume the latter — anything broader
  is uncheckable.

## Risks

- **Ceremony drift.** Nine obligations at once is how protocols acquire dead text. My
  amendments deliberately make three of them conditional (register only on conflict,
  steelman only on unanimous deliberation, settledness only on novelty claims); if the
  group prefers the unconditional forms, expect the rules to be quietly ignored within a
  few ideas — an unenforced rule is worse than none, per the brief's own constraint.
- **False confidence.** Tags and registers can become ritual compliance — agents stamping
  `PRIMARY` on recall. The mitigation (falsifiable locators/outputs) only works if
  reviewers actually spot-check them; LE-1 refutation-default already demands this muscle.
- **Track mismatch.** Binding MAJOR-7 on `standard` (as the constraint permits but I
  advise against) would tax the most common track for a failure mode that Phases 6–8
  already catch; binding it on `deliberation` only keeps the cost where the risk lives.
- **Version drift in tooling testimony.** T3 did not reproduce on 1.37.0. If the group
  endorses all six defects wholesale on the strength of the brief, that is itself the
  unverified-endorsement failure this idea exists to prevent; the record should show T3
  as "not reproduced at 1.37.0" rather than "confirmed".
