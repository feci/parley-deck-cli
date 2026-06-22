---
agent: claude-1
idea: addon-skills-research
round: 1
date: 2026-06-22
---

## Summary

Both tracks share one principle: **the parley-deck core stays thin and agnostic; the
skills carry the mechanics.** Worktrees and ticketing are both "execution surfaces"
that the protocol already has a slot for — worktrees extend §6 conflict-avoidance /
Phase 5 isolation; trackers extend the existing "PR/MR comments are mirrors, not
canonical" rule. Neither should leak vendor or git plumbing into COOPERATION.md.

## Proposed approach

### Track A — worktree-sessions (integration seam)
- **Mapping:** Parley's `github-pr` transport already implies **branch-per-idea**. A
  worktree is the *local* equivalent: `git worktree add ../wt/<idea-slug> -b <branch>`
  → each idea/session (or each Phase-5 implementer) gets an isolated dir over the same
  `.git`. The deck artifacts (`parley-deck/ideas/<slug>/`) live on that branch.
- **Skill owns:** add/list/remove/prune, branch-per-idea naming, per-worktree
  `.env.local`, sparse-checkout to scope files, stale-worktree cleanup, staging-branch
  merge flow, "don't check out the same branch twice" guard.
- **Core owns (thin, separate meta-protocol-change):** one line in §6/Phase 5 —
  "parallel implementers/sessions MAY run in isolated worktrees (branch-per-idea)."
- **Runtime overlap:** Claude Code / Agent tool already do worktree isolation, so the
  skill should *defer to runtime worktree support when present* and only provide the
  manual git flow as fallback. Don't reinvent what the runtime gives.
- **Coordination:** the "shared task doc agents claim/mark" pattern from the wild IS
  Parley's idea + quorum + one-file-per-agent — call that out so users see worktrees
  as the *physical* layer under Parley's *logical* coordination.

### Track B — ticket-tracker (three-audience template)
- **One artifact, three layers** (not three docs): a story has a **Why** (business,
  1 line: "As <role> I want <goal> so that <benefit>"), a **What/How** (technical:
  scope, constraints, NFRs, affected files/paths), and **Acceptance criteria**
  (Given/When/Then + measurable bullets) that double as the AI agent's spec and the
  tester's checklist. Same text serves all three by *layering*, not duplicating.
- **AI-readability rules:** self-contained (no "see Slack"); explicit edge/error/offline
  states + NFRs (AI won't assume them); file/path hints + links to OpenAPI/architecture;
  a machine-parseable metadata block (epic, status, priority, labels, DoD).
- **Hierarchy:** Epic (outcome/value, spans stories) → Story (INVEST, one increment) →
  Technical subtask (implementation unit, can map 1:1 to a Parley idea or a worktree).
- **Generic schema** over vendor fields: define our own neutral field set, then map to
  Jira/Linear/GitHub Issues at the edge. Tracker = mirror; `FINAL.md`/idea = canonical.
- **Parley seam:** `FINAL.md` → epic; its work-items → stories; Phase-5 checklist →
  subtasks. Sync is one-way-canonical (deck → tracker) by default.

## Concerns / open questions

- Worktrees + the deck dir: if `parley-deck/` is committed, per-idea worktrees branch
  it too — need a rule so concurrent ideas don't fight over `parley-deck/ideas/`
  (answer: each idea on its own branch/worktree; merge via staging). Codex should nail
  the exact collision model.
- How prescriptive should the ticket template be vs. configurable? Lean: opinionated
  default template + a small config for field names per tracker.
- Do we ship a `ticket-tracker` that actually *calls* a tracker API/MCP, or one that
  *authors the text* and hands the API call to the user/runtime? I lean author-text +
  optional MCP, to stay vendor-neutral and avoid auth in the skill.

## Risks

- Scope creep: a worktree skill can balloon into a full git-workflow manual. Keep it to
  "isolate parallel Parley work + clean up." 
- Ticket template that's too heavy → nobody fills it. Must stay one screen, with
  optional depth.
- Vendor leak: hardcoding Jira field names would repeat the genericity bug we just fixed.
