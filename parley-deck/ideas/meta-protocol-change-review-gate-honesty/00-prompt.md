---
idea: meta-protocol-change-review-gate-honesty
author: user
created: 2026-06-12
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + drafter of the protocol amendments
  codex: precision of normative language; interaction with driver gates (Phases 6-8)
  agy: reviewer-experience lens — what the rules feel like to follow; failure modes of over-strictness
  hermes: protocol-consistency lens — collisions with existing sections, quorum, signoff rules
transport: local-dir
cross_review_rounds: 1
status: final
---

## Task (owner's directive)

Amend the Parley Deck protocol (COOPERATION.md) with two review-gate practices
adopted from the MIT-licensed "kindly" skill (reference copies live in
`../runner-hardening-kindly/reference/`, esp. kindly-audit.md "Gate Standard"
and "Stopping Judgment"). Ships with release 1.24.0 alongside the sibling CLI
idea `runner-hardening-kindly`. Both points are IN SCOPE by owner decision —
deliberation refines wording and placement, not the scope.

## The two points

**P6 — No-suppression review briefs.**
Today facilitators sometimes instruct reviewers "do NOT report X" / "do not
re-raise the deferred items" (it happened in tui-protocol-visibility for the
sandbox artifact and the deferred follow-ups). Adopt kindly's rule: briefs and
review prompts may carry DISPOSITIONS of known findings (a rebuttal, an accepted
tradeoff, a recorded sandbox artifact) as context the reviewer weighs OPENLY —
the reviewer states in its review whether it concurs — but never as suppression:
no "do not re-raise" carve-outs, no severity floors, no instruction that narrows
what a review may report. A disputed finding closes only by the operator's
(owner's) explicit call, relayed verbatim. Decide: exact normative text and
where it lands (Phase 6 section; review-prompt construction guidance for the
runner's BuildReviewPrompt is the sibling idea's concern but the rule lives here).

**P7 — Optional strict gate + trajectory-based stopping judgment.**
Default Phase 8 today: repeat review until zero AGREED fixes (dismissals and
deferrals negotiated in review consensus); MaxFixupCycles bounds the loop.
Add an OPT-IN stricter close standard for high-stakes ideas (frontmatter flag,
e.g. `strict_gate: true` in 00-prompt.md): the gate closes only on a FRESH
full-scope review pass whose verdict is zero findings of any severity — fix-
verification passes converge the gate but never close it; no severity floor, no
nitpick allowance; operator rulings are the only other way to close a finding.
Also add (for all gates) kindly's trajectory-based stopping guidance: read the
trajectory, not a counter — converging (fewer, lower-severity, confined to new
code → continue), churning (fresh High/Medium findings landing on fix-pass code,
or re-litigated ground → stop and escalate with the trajectory), blocked (a
decision finding pauses its thread until the operator answers). Decide: exact
wording, placement (Phase 8 + a new subsection), interaction with
MaxFixupCycles and the driver's review gate (machine-readable bits are the
sibling idea's concern; flag semantics defined here).

## Constraints

- This is a protocol-change idea: amendments land in the live
  `parley-deck/COOPERATION.md`, the protocol changelog
  (`parley-deck/meta/protocol-changelog.md`), and the embedded default protocol
  + skill snapshot copies if those exist in-repo (verify: internal/protocol
  defaults, references/COOPERATION.md) — participants should confirm the full
  set of copies that must stay in sync.
- Amendments must not contradict existing invariants: append-only signoffs,
  one-file-per-agent ownership, quorum rules, non-solo execution, Phase 7
  review-consensus triage semantics (agreed fixes / deferred / dismissed —
  P6 constrains HOW dismissals are communicated to reviewers, not whether the
  triage categories exist).
- English-only under parley-deck/.

## Questions for round-01

1. P6: exact normative paragraph(s) + where (Phase 6? a new "Review briefs"
   subsection?); how dispositions should be phrased in briefs (template?).
2. P7: exact `strict_gate` semantics (who may set it, when it can be added or
   removed, what "fresh full-scope pass" means operationally for participants);
   trajectory guidance wording; MaxFixupCycles interaction.
3. Which protocol copies in this repo must be updated in lockstep, and does the
   protocol version/changelog entry need anything beyond a dated entry?
