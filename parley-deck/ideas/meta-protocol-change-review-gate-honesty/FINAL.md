---
idea: meta-protocol-change-review-gate-honesty
status: final
drafted-by: claude
date: 2026-06-12
---

## Summary

Amend the Parley Deck protocol with review-gate honesty rules adopted from the
"kindly" skill: a no-suppression rule for review briefs (dispositions are
weighed openly, never silenced), an opt-in `strict_gate` close standard
(fresh full-scope zero-findings pass), trajectory-based stopping judgment, and
a standing note for consult artifacts. Ships with release 1.24.0 alongside the
sibling CLI idea `runner-hardening-kindly`. Consensus: all four participants
✅ ACCEPT.

## Final plan / specification

The authoritative decision record is consensus.md D1-D6; the editing contract:

1. **Phase 6 — new subsection "Review briefs and dispositions"** (codex base
   text + agy's Reviewer-Prompt line in the disposition shape): briefs must
   not suppress findings; dispositions travel openly; the reviewer states
   concurrence per disposition; disputed findings close only via reviewer
   withdrawal, normal review consensus, or a verbatim-quoted operator ruling.
2. **Phase 8 — new subsection "Strict review gate (optional)"**: the
   `strict_gate: true` frontmatter flag; default close rule unchanged without
   it; with it, closing requires a fresh full-scope Phase 6 pass (the complete
   implementation diff at pass time) reporting zero findings of any severity,
   recorded by the following Phase 7 consensus; fix-verification passes
   converge but never close; NITs and deferrals keep the gate open; a finding
   must be an objective, code-grounded issue — subjective stylistic preference
   is never a finding; mutability only via consensus or recorded operator
   direction.
3. **Phase 8 — new subsection "Stopping judgment"**: trajectory over pass
   counters (converging/churning/blocked with agy's concrete triggers as
   examples); MaxFixupCycles and driver budgets are escalation thresholds,
   never close criteria.
4. **§8 inbox/communication — new "Consults" note**: consult artifacts under
   parley-deck/consults/ are advisory, non-canonical, never quorum evidence;
   promotion requires normal protocol artifacts.
5. **Frontmatter documentation** gains the `strict_gate` line (exact "true"
   enables; ReadAutoImplement parsing precedent for the future ReadStrictGate;
   no driver enforcement this release).
6. **Copies in lockstep**: live parley-deck/COOPERATION.md,
   internal/protocol/defaults/COOPERATION.md, and a dated
   parley-deck/meta/protocol-changelog.md entry; external skill snapshot
   flagged via inbox in the release step.

## Implementer

claude (drafter; edits ship via the sibling idea's implementation branch with
this idea's slug in the protocol-edit commits). Reviewers: codex, agy, hermes.

## Deferred follow-ups

- Driver ReadStrictGate + machine-readable strict-close fields in review
  consensus (no prose scraping).
- The embedded default protocol's broader §12 drift (separate follow-up idea).
