---
idea: meta-protocol-change-agent-teams-patterns
review-cycle: 1
drafted-by: codex
date: 2026-05-14
reviewed-commit: 181fd08
---

## Agreed fixes

- From `claude/review/round-01` [MAJOR] `IMPLEMENTATION.md head-commit and transport metadata do not match the working tree`: update `IMPLEMENTATION.md` so it truthfully records that this implementation is currently a local uncommitted working-tree update on top of `181fd08`, with no PR created in this turn, and record that as a transport deviation prompted by the user's direct request.
- From `claude/review/round-01` [MINOR] `Phase 5 plan-gate wording creates a timing ambiguity in IMPLEMENTATION.md creation`: change Phase 5 wording from "creates" to "opens or updates" / "finalizes" and clarify that `IMPLEMENTATION.md` may begin with the plan/checklist before code changes.
- From `claude/review/round-01` [NIT] `Deviations from FINAL.md: None understates the SKILL.md / installed-skill scope expansion`: record the skill and installed-skill propagation as a user-directed scope extension, not as an unmentioned no-op.
- From `gemini/review/round-01` [MINOR] `§8 Inbox mirroring rule missing IMPLEMENTATION.md`: include `IMPLEMENTATION.md` as a canonical mirror target for phase-affecting build decisions.
- From `gemini/review/round-01` [NIT] `Phase 5 Plan-gate Normal Review Process ambiguity`: add light wording that risky implementation plans may be shared for a brief feedback window through the active transport or inbox, without creating a new phase.
- From `hermes/review/round-01` [MINOR] `Phase 5 guidance lacks explicit checklist example`: add a short implementation plan/checklist template.
- From `hermes/review/round-01` [NIT] `Changelog entry omits explicit link to idea directory`: add an explicit idea path reference to the changelog summary.

## Deferred follow-ups

- Native GitHub PR branch mechanics for this already-in-progress local protocol update are not created in this turn. If the user wants the strict GitHub transport mirror, open a follow-up branch/PR that contains the accumulated protocol and skill changes.
- Whether bundled fallback `references/COOPERATION.md` should have a transport-neutral header instead of inheriting `Transport: github-pr` is deferred to a separate protocol-change idea.

## Dismissed findings

- `gemini/review/round-01` [NIT] `Grammatical consistency in §4 Phase 4`: Gemini self-corrected this as a non-issue.

## Signoffs

<!-- Each active participant (implementer included) APPENDS their signoff block. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept review consensus and will apply the agreed documentation/protocol fixes in fix-up cycle 1.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: All three of my round-01 findings (MAJOR audit-trail mismatch, MINOR Phase 5 timing ambiguity, NIT understated SKILL.md scope) are captured in Agreed fixes, and my two open questions are addressed in Deferred follow-ups. Gemini's and Hermes's findings are folded in coherently. No blockers.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: My findings regarding §8 (mirroring in `IMPLEMENTATION.md`) and §4 Phase 5 (plan-gate feedback window) are correctly captured and resolved in the agreed fixes. I also support the resolutions of Claude's and Hermes's findings.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: My two round-01 findings (MINOR checklist example and NIT changelog link) are fully addressed in the agreed fixes. All other resolutions look coherent and preserve protocol invariants. No blockers.
