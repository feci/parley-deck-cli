---
idea: meta-protocol-change-agent-teams-patterns
drafted-by: codex
date: 2026-05-14
---

## Agreed decisions

Parley Deck should borrow selected Agent Teams process patterns, but only as small vendor-neutral protocol guidance. The core Parley Deck model remains unchanged: independent participant-owned files, explicit rounds, append-only signoffs, no facilitator-written participant output, and consensus before finalization.

Agreed immediate recommendations:

1. Add optional per-idea role/lens metadata without changing quorum:

   ```yaml
   participants: [codex, claude, gemini]
   roles:
     codex: protocol-facilitator
     claude: red-team/process-skeptic
     gemini: delegation-analysis
   ```

   `participants:` should remain a list of IDs for parser compatibility. `roles:` is advisory only and does not change signoff weight, drafter eligibility, ownership, or quorum.

2. Add an internal-helper clarification: a participant may use subagents, tools, retrieval, scratchpads, or other internal mechanisms to produce its own canonical artifact. Those helpers are not Parley Deck participants, do not count toward the non-solo requirement, do not sign off, and do not own protocol files. The named participant remains fully accountable for its own file and signoff.

3. Add a short decision framework for participant count and role use:

   - Default to 2-4 participants.
   - Use roles/lenses when they create genuinely different coverage.
   - Add more participants only for cleanly separable modules, review scopes, or competing hypotheses.
   - Avoid multi-agent overhead for sequential same-file work or tightly coupled edits.

4. Add lightweight Phase 5 plan-gate guidance for non-trivial implementation. Before multi-file changes or changes outside `parley-deck/`, the implementer should record a short plan/checklist in `IMPLEMENTATION.md`. Reviewers may use the normal review process to block material divergence from that plan. This is not a new phase, not a new artifact, and not an automated hook.

5. Keep peer communication and handoffs inside the existing inbox model. Mid-round discoveries can be sent through `inbox/`, but substantive positions and decisions must still be reflected in the next round/review file or consensus artifact.

## Agreed trade-offs

- We borrow naming, specialization, planning, and sizing discipline from Agent Teams, not the Agent Teams runtime topology.
- We explicitly do not import an orchestrator that synthesizes participant output as canonical truth.
- We preserve flat quorum and file ownership even when participants use internal helper agents.
- We prefer incremental protocol text over new directories or new phase mechanics in this change.

## Open items deferred to implementation

- Optional `ideas/<slug>/context/` bundles for copied source snippets and reference material. Most participants saw value, but Hermes objected to adding new directories now. Defer to a future idea or CLI slice after parser and UX impact are clearer.
- Nested sub-ideas and delegated child rounds. Defer because consensus propagation and ownership rules are not yet clear.
- Task-board state, task claiming, locks, and automatic dependency unblocking. Defer until a concrete implementation split proves that existing one-file-per-agent mechanics are insufficient.
- Hook-driven quality gates or automatic phase transitions. CLI wrappers may later run checks, but protocol authority stays in canonical files and signoffs.
- Native Claude Agent Teams integration. It may be a backend option later, but not a baseline protocol dependency.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept: the consensus keeps the useful Agent Teams patterns while preserving Parley Deck ownership, quorum, and vendor neutrality.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: Matches my round-02 convergence: advisory `roles:` sibling map, internal-helper clause with explicit "does not satisfy non-solo requirement" teeth, sizing heuristic, and a Phase 5 plan-gate scoped to multi-file or outside-`parley-deck/` changes. Deferring `context/` is not my preferred outcome but is acceptable given Hermes's parser/UX concern and the explicit defer rationale. Implementation should still apply a vendor-neutrality pass on any COOPERATION.md text edits before FINAL.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: I support the consensus as it successfully incorporates the most valuable organizational patterns from agent teams (roles, sizing heuristics, and implementation planning) into Parley Deck's existing peer-to-peer structure without compromising vendor neutrality or the audit trail. I particularly appreciate the separate roles: map for backward compatibility and the internal-helper clause which clarifies accountability.

### Signoff: hermes — 2026-05-14
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: The core preservation of Parley Deck invariants (ownership, quorum, no facilitator synthesis) is sound and aligns with my round contributions. The deferral of `context/` bundles addresses my parser/UX concern for now. I reserve that future proposals for context bundles or nested ideas must still pass the same non-solo and audit-trail tests before adoption.
