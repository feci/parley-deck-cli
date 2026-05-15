---
agent: codex
idea: meta-protocol-change-agent-teams-patterns
round: 1
date: 2026-05-14
---

## Summary

Parley Deck and Claude Code Agent Teams share the same core premise: useful multi-agent work needs independent contexts, clear task boundaries, durable coordination state, and a lead/facilitator that synthesizes the result. Parley Deck is stronger on auditability, non-solo guarantees, file ownership, signoff, and transport independence. Agent Teams are stronger on live coordination ergonomics: task lists, direct teammate messaging, named roles, progress visibility, plan approval before risky work, and quality gates.

The best move is not to replace the protocol. We should borrow the process patterns that improve coordination without weakening canonical artifacts.

## Proposed approach

Add a small "team coordination layer" to the protocol and CLI on top of existing phases:

1. Add optional `roles:` and `focus:` metadata to kickoff prompts or participant launch summaries. Example roles: `protocol-analyst`, `implementation-skeptic`, `ux-reviewer`, `test-strategy`, `security-reviewer`. This mirrors Agent Teams' named teammates, while keeping stable Parley agent IDs for ownership.

2. Add an optional task board artifact under each idea:

   ```text
   parley-deck/ideas/<slug>/tasks.md
   ```

   It should be facilitator-owned or append-only, with simple states: `pending`, `claimed`, `blocked`, `done`. This borrows shared task-list visibility without allowing agents to overwrite one another's protocol files.

3. Clarify when to use Parley Deck multi-agent depth:

   - Use 2-4 participants by default.
   - Add more only when work splits cleanly by subsystem, review lens, or competing hypothesis.
   - Prefer a single agent or a lightweight helper for sequential same-file work.

4. Add a "plan gate" before implementation for risky slices. The implementer writes a short implementation plan or checklist before code changes, and reviewers can block the plan before expensive work starts. This adapts Agent Teams' plan approval pattern to Parley Deck's durable review model.

5. Expand inbox semantics into explicit peer messages for mid-round coordination, but keep them lightweight and non-authoritative:

   ```text
   parley-deck/inbox/<from>-to-<agent-or-all>_<slug>_<topic>.md
   ```

   Messages can share discoveries or request context, but substantive positions still belong in the next round/review file.

6. Add quality gates as protocol checklists and later CLI hooks: tests run, artifact shape valid, no participant file missing, no same-file collision if implementation is split, no unresolved `to-user` escalation.

7. Add a cleanup/recovery section for agent runs: detect stale runs, pending handoffs, orphaned interactive sessions, missing signoffs, and unresolved task-board claims.

## Concerns / open questions

- Task boards can become another source of truth. They must remain coordination aids; `FINAL.md`, `IMPLEMENTATION.md`, reviews, and consensus remain authoritative.
- Peer messaging can reduce latency, but it can also bypass cross-review. The protocol should require that any decision-making message be summarized in the next canonical round file.
- Automatic task claiming and locking are attractive, but Parley Deck already avoids file collisions by assigning one artifact per agent. We should not add locks until implementation work is actually split across disjoint code paths.
- Claude Agent Teams are experimental and provider-specific. We can support them as one backend later, but not encode their assumptions into the core protocol.

## Risks

- Over-automation could make Parley Deck less auditable.
- More artifacts could slow small tasks.
- Direct inter-agent messaging could create hidden consensus if not mirrored into round files.
- Role specialization can help reviews, but hard-coded roles would reduce the current vendor-neutral flexibility.
