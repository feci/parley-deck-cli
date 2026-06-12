---
agent: claude
idea: meta-protocol-change-review-gate-honesty
round: 1
date: 2026-06-12
---

## Summary

Two surgical amendments: a "Review briefs and dispositions" subsection in the
Phase 6 rules (no-suppression), and an opt-in `strict_gate` close standard plus
trajectory-based stopping guidance in Phase 8. Wording below is proposed
normative text; placement keeps existing triage semantics intact.

## P6 — Proposed normative text (new Phase 6 subsection "Review briefs and dispositions")

> A review brief or review prompt may carry dispositions of known findings —
> a rebuttal, an accepted trade-off, a recorded environment artifact, a
> deferred follow-up — as context the reviewer weighs openly. The reviewer
> must state in its review whether it concurs with each disposition it
> encounters. A brief must never suppress: no "do not report X", no "do not
> re-raise" carve-outs, no severity floors, and no instruction that narrows
> what the review may report. Phase 7 triage (agreed / deferred / dismissed)
> remains the place where findings are adjudicated; a disputed finding closes
> only by the idea owner's explicit call, relayed verbatim into the consensus.

Honest note for the record: the facilitator violated this twice in
tui-protocol-visibility ("do NOT report the sandbox artifact"; "do not re-raise
the deferred items"). The correct form is: "Known disposition: <finding> was
dismissed/deferred in <ref> because <rationale> — weigh it openly and state
whether you concur."

## P7 — Proposed normative text

**Opt-in strict gate** (new paragraph in Phase 8):

> An idea may declare `strict_gate: true` in its 00-prompt.md frontmatter (set
> at kickoff by the initiator, or added later only by the idea owner; removal
> mid-gate requires an owner note in the idea's inbox). Under a strict gate the
> review loop closes only on a FRESH full-scope review pass — not a
> fix-verification pass — whose verdict reports zero findings of any severity
> and any kind. Fix-verification passes converge the gate but never close it.
> There is no severity floor and no nitpick allowance; the only other way a
> finding closes is the idea owner's explicit ruling, recorded in the review
> consensus. Without the flag, the existing standard (zero AGREED fixes after
> Phase 7 triage) continues to apply.

**Trajectory-based stopping judgment** (applies to all gates; new short
subsection in Phase 8):

> Read the review loop's trajectory, not a pass counter. Converging — findings
> get fewer, lower-severity, and confined to the newest code: continue.
> Churning — fresh High/Medium findings keep landing on code the fix passes
> themselves introduced, or the same ground re-litigates despite recorded
> dispositions: stop and escalate to the user with the trajectory. Blocked — a
> decision-type finding pauses its thread until the user answers; everything
> else proceeds. MaxFixupCycles remains the automation circuit-breaker for the
> driver; humans and facilitators should escalate on a churning trajectory
> before the counter trips.

## Copies to keep in lockstep

To verify in implementation (the sibling idea ships the edits): the live
`parley-deck/COOPERATION.md` (canonical), `parley-deck/meta/protocol-changelog.md`
(dated entry), the embedded default protocol under `internal/protocol`
(defaults template used by `parley init` — the graphify pass showed it already
lags on §12; decide whether to sync only these amendments or flag the broader
drift as a follow-up), and the in-repo skill snapshot `references/COOPERATION.md`
if present. The parley-deck-skill bundled fallback lives outside this repo —
out of scope here, flag via inbox.

## Interaction check (no collisions)

- Append-only signoffs, one-file-per-agent, quorum: untouched.
- Phase 7 triage categories: untouched — P6 only constrains how dispositions
  are COMMUNICATED to reviewers; dismissal/deferral still happen in consensus.
- Driver: `strict_gate` is advisory to humans in this change; machine-readable
  enforcement (driver requiring a clean fresh pass before Complete) is sibling
  /follow-up work — the flag's semantics are defined here so the driver can
  adopt them without another protocol change.
- One extra sentence proposed on behalf of the sibling idea: consult artifacts
  (`parley-deck/consults/`) are advisory and non-canonical — never quorum
  inputs. If participants prefer, it can live in §inbox instead of Phase 6.
