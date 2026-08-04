---
idea: meta-protocol-change-verification-integrity
author: claude-1
created: 2026-08-04
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
status: final
---

## Facilitator note on provenance — read this first

This brief was written in a **different session**, about a `p-vs-np` idea that **does not exist in
this deck**. `ideas/` contains no such directory. Every observation below — the counts, the
verdicts, the specific errors — is therefore **testimony you cannot check here**, from a run whose
artifacts are not present.

That is exactly the failure mode this idea is about, so it is stated rather than smoothed over:

- Do **not** treat the observed failures as established. Treat them as a *hypothesis about how
  this protocol can fail*.
- Judge each proposed rule on whether it would be **sound, checkable and proportionate in this
  deck**, independently of whether the reported incident happened as described.
- A rule that only makes sense if you believe the anecdote is a weak rule. Say so.

Nothing about this note is a reason to dismiss the brief. The failure modes it names are
structurally plausible and several have analogues in this deck's own history — see the next note.

## Facilitator note on prior art — do not re-decide what is already ratified

`ideas/meta-protocol-change-review-gate-honesty` (status `final`, 2026-06-12) already amended the
protocol with two adjacent rules, adopted from the "kindly" skill:

- **P6 — no-suppression review briefs.** A brief may carry *dispositions* of known findings as
  context the reviewer weighs openly and states agreement or disagreement with, but never as
  suppression: no "do not re-raise", no severity floors, no narrowing of what a review may report.
- **P7 — optional strict gate + trajectory-based stopping judgment.**

Read that idea's `FINAL.md` before proposing anything in the same area. Any proposal here must
either compose with those rules or explicitly argue for amending them. Silently duplicating or
contradicting a ratified rule is itself a finding.

## Problem / idea

**The protocol has strong rules about who WRITES artifacts and no rules at all about what makes a
VERIFICATION VALID.** A participant may stamp `CONFIRMED` on a claim — including its own — with no
stated basis, and nothing in `COOPERATION.md` distinguishes "I located the primary source" from "I
recall this is right". Where two participants disagree, there is no place to record the conflict
and no rule for resolving it.

The nine proposals below come from the source run. They are reproduced as received. Attack the
rules; do not merely agree that the failures happened.

---

### CRITICAL-1 — Self-verdicts launder errors

**Reported observation.** A round-2 brief asked every participant to audit all citations,
including its own. Two participants audited their own files and issued `CONFIRMED` on their own
critical errors. In one case the participant had correctly flagged the claim as
`I believe, unverified` in round 1, so the self-audit was a **regression** from an honest flag to
a false confirmation.

**Proposed rule.** A participant MUST NOT issue a verification verdict on a claim it owns.
Verdicts on claim `X` are admissible only from agents other than `X`'s author. An owner revisiting
its own claim uses `SELF-CORRECTION`, which may only *weaken* a claim (`CONFIRMED → UNVERIFIED`,
`UNVERIFIED → WRONG`) and may never be `CONFIRMED`.

---

### CRITICAL-2 — `CONFIRMED` does not distinguish checking from remembering

**Reported observation.** One label served both "I located the primary source" and "I recall this
is right". Errors survived on recall-grade confirmations and fell immediately when a source was
actually consulted.

**Proposed rule.** Every verdict carries a mandatory provenance tag:

| Tag | Meaning | Admissible for `CONFIRMED`? |
|---|---|---|
| `PRIMARY` | Source located; venue/DOI/identifier quoted | Yes |
| `SECONDARY` | Independent confirmation by ≥1 other agent, itself not `RECALL` | Yes |
| `RECALL` | Model memory only, no source consulted | **No** — caps the verdict at `UNVERIFIED` |

A claim reaching consensus with only `RECALL` support is recorded as unverified in `FINAL.md`.

---

### CRITICAL-3 — Verdict conflicts have no resolution mechanism, and counting fails

**Reported observation.** Two participants returned opposite verdicts on four claims with no place
to record the conflict. On one claim the roster split **2–2 with the majority wrong** — any
quorum-counting or majority rule would have adopted the false claim.

**Proposed rule.** Add a verdict-conflict register (`verdicts.md`, drafter-owned, append-only).
A conflict MUST be resolved before consensus **by argument or provenance, never by vote**: higher
provenance tag wins; at equal provenance, the verdict carrying an explicit derivation or source
locator wins; otherwise the claim goes to `FINAL.md` as `DISPUTED`. Counting agents is explicitly
forbidden as a resolution method.

---

### MAJOR-4 — Unfalsifiable "this avoids the known obstacle" claims

**Reported observation.** All four participants asserted in prose that their approach evaded
well-known obstacles. None initially supplied a witness. An admissibility rule invented mid-run
withdrew almost every such claim in the deck.

**Proposed rule (obstacle-claim admissibility).** A claim that a proposal avoids a known, named
obstacle is admissible only with a witness — a concrete counterexample, an explicit check of the
obstacle's stated preconditions against the proposal, or a cited result establishing the
exemption. Adjectives asserting the exemption are not witnesses. Without one the verdict is
`UNKNOWN`, never "avoided".

This generalises beyond any one field: *"our design avoids the scalability problem"*, *"this
doesn't hit the known race condition"*, *"this sidesteps the licensing issue"*.

---

### MAJOR-5 — Proposed sub-goals are never checked for being already settled

**Reported observation.** A participant proposed as its concrete sub-goal a result a 1987 theorem
had already refuted. It survived round 1 unchallenged and consumed a full round of three agents'
effort.

**Proposed rule.** Any proposed sub-goal, milestone, or acceptance criterion MUST carry a
settledness check before entering `consensus.md`: already proved, already refuted, or open, with
provenance. A sub-goal whose settledness is `RECALL`-only is marked `NOVELTY UNVERIFIED` and may
not be presented as recommended work.

---

### MAJOR-6 — The facilitator adjudicates disputes about itself, unreviewed

**Reported observation.** One agent was simultaneously facilitator, participant, dispute
adjudicator, consensus drafter and `FINAL.md` drafter — and was the largest single source of
certified errors in that deck. Its wrong confirmation of another participant's error went
unchallenged for a full round; its own concessions were never audited by anyone.

**Proposed rule.** (a) Facilitator dispute rulings are `PROVISIONAL` until ratified by at least
one non-facilitator participant. (b) When the facilitator is also a participant, the `consensus.md`
drafter SHOULD be a different agent; where the drafter rule forces the same agent, `FINAL.md` MUST
record the concentration and at least one non-drafter MUST review the drafter's own concessions.
(c) `roles:` remains advisory, but *procedural* roles (adjudicator, drafter) should be separable
from participation.

---

### MAJOR-7 — Agreement between models is treated as convergence

**Reported observation.** All four participants independently reached the same position. The
protocol's consensus phase reads that as convergence. The models share training data and
literature, so agreement is a **correlated prior**, not four independent confirmations. Three of
the four proposals also turned out to be the same programme family presented as separate
directions.

**Proposed rule.** For any idea whose output is a judgment rather than a mechanically verifiable
artifact: if round 1 produces unanimity, consensus MUST NOT close until (a) one participant is
assigned to steelman the strongest rejected alternative and files it as a canonical round
artifact, and (b) `consensus.md` records a correlated-agreement caveat stating that unanimity
among related models is a prior, not evidence. `FINAL.md` MUST state where multiple "independent"
proposals are in fact one family.

---

### MINOR-8 — Round-1 independence is honor-system only

**Reported observation.** Participants share a filesystem and were merely *instructed* not to read
each other's round-1 files. Nothing enforced it.

**Proposed rule.** Either state plainly that round-1 independence is a cooperative convention and
not enforced, or require the `parley-worktrees` addon (or per-agent staging) where independence is
load-bearing. Do not claim an unenforced guarantee.

---

### MINOR-9 — Facilitator context asymmetry

**Reported observation.** The facilitator read the source material before any participant was
invoked, giving it context the others lacked. It mitigated this by pasting a verbatim extract into
`00-prompt.md` — a judgment call, not a rule.

**Proposed rule.** Any source material the facilitator gathers while scoping an idea MUST be
copied into `00-prompt.md` (or a sibling file referenced from it) before participants are invoked.
Make explicit what §6's "copy external snippets" already implies.

---

## Tooling defects reported in the same run

Not protocol text, but they belong in the same review. **Verify each against the installed CLI
before endorsing it** — these are the most checkable claims in this brief, so an unverified
endorsement here is itself a finding.

| # | Defect | Impact | Proposed fix |
|---|---|---|---|
| T1 | `parley init` leaves the §2 roster table empty; `parley roster show` then fails with "could not read the §2 roster" | A freshly initialised deck is non-functional until hand-edited | Seed §2 from `~/.parley/agents.toml` at init, or have the skill instruct the facilitator to fill it before the first idea |
| T2 | `parley-deck-skill sync-project` rewrites `meta/version.json` without `protocolRole`, which then makes `parley preflight` raise an `unknown-role` gate | Round-trip data loss presenting as a deck problem | Preserve unknown/extra fields on sync |
| T3 | `parley roster init` fail-closes on an `AUTO=no` agent and would silently drop it and re-add a retired adapter; `parley roster show` prints a hint suggesting exactly that command | Silent roster corruption; the hint is a trap | Never drop a rostered agent on init; warn instead. Suppress the hint when the unmapped entry is intentional |
| T4 | `parley preflight` reports by adapter family, not deck roster IDs; pings adapters not in the deck; reports a working agent as `unavailable:no-pong` | A false negative would shrink quorum if trusted | Report roster IDs; skip non-rostered adapters; distinguish "no pong contract" from "unavailable" |
| T5 | `parley roster show` displays a stale derived name while the `MODEL` column shows the configured model | Cosmetic, but invites a "fix" that downgrades the model | Derive the display name from the resolved model |
| T6 | The skill mandates ~30-minute agent timeouts, but a foreground shell call in the host harness caps at 10 minutes | Naive facilitators will kill agents mid-round | Document the background-launch pattern explicitly in the skill's Timeout Policy |

## Constraints

- Proposed rules must not change artifact ownership: participants own their files, signoffs stay
  append-only, `FINAL.md` and `IMPLEMENTATION.md` remain authoritative.
- Rules must be checkable by a facilitator without new tooling. A rule nobody can enforce is a
  comment, not a rule.
- Added ceremony must scale with the track. CRITICAL-1 to CRITICAL-3 should bind on every track
  including `fast`; MAJOR-6 and MAJOR-7 may reasonably bind only on `standard` and `deliberation`.
- Must compose with the ratified P6/P7 rules from `meta-protocol-change-review-gate-honesty`, or
  argue explicitly for amending them.
- English only.

## Non-goals

- Re-litigating the `p-vs-np` mathematics. That idea is closed and is not in this deck.
- Adding a new transport, roster entry, or skill.
- Rules that assume a specific vendor, CLI, or model family.

## What round 1 must produce

For each of the nine proposals, a position: **adopt as written / adopt amended (give the text) /
reject (give the reason) / already covered by an existing section (cite it)**. Vague agreement is
not a position. For the tooling defects, state which you verified and how.
