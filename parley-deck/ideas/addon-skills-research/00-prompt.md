---
idea: addon-skills-research
author: user
created: 2026-06-22
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: synthesis & parley-deck integration seams (skill vs core boundary)
  codex-1: worktrees — mechanics, pitfalls, exact commands for our case
  hermes-1: ticketing — three-audience writing (business / technical / AI)
  antigravity-1: standards & prior-art for both (INVEST/Gherkin/DoD; worktree patterns)
status: final
---

## Problem / idea

Research to enable building **two vendor-neutral companion skills** for the Parley
Deck ecosystem (core + opt-in plugins, install-all-by-default with `--no-addons`):

1. **`worktree-sessions`** — run multiple Parley sessions / parallel implementers
   over ONE git repo via git worktrees, without collisions.
2. **`ticket-tracker`** — author epics / user stories / technical subtasks for a
   GENERIC ticket/kanban tracker (Jira, Linear, GitHub Issues, Trello, …), written
   so they're readable by **business people, technical people, AND the AI agents
   that will implement them**.

Each participant addresses **both** tracks (lead with your lens). The output is a
`FINAL.md` research/design doc concrete enough to author both `SKILL.md`s from —
templates, commands, decision rules, and the parley-deck integration seam.

## Track A — worktrees (research questions)

- Exact mechanics for our case: `git worktree add/list/remove/prune`, shared `.git`,
  branch-per-idea/session, isolated dir per agent. How does a Parley session or a
  Phase-5 implementer get its own worktree, and how is it cleaned up?
- Decomposition & coordination: split by domain/feature boundary; a shared task doc
  agents claim/mark (this maps onto Parley ideas/quorum — make the seam explicit).
- Merge strategy: integration/staging branch → test → main.
- Pitfalls: can't check out the same branch in two worktrees; per-worktree env
  (`.env.local` in `.gitignore`); stale worktrees → `prune`; submodules; sparse-checkout
  to constrain files; IDE confusion.
- Runtime overlap: native worktree support now exists in Claude Code / Cursor /
  Copilot CLI (late-2025→2026) and the Agent/Workflow tooling has worktree isolation —
  so decide what the SKILL adds vs. what the runtime already provides.
- Seam: what belongs in `worktree-sessions` (the git mechanics/patterns) vs a thin
  note in the parley-deck protocol (§6 conflict-avoidance / Phase 5 "implementers MAY
  be isolated in worktrees").

## Track B — ticketing (research questions)

- Templates for **epic / user story / technical subtask** that satisfy three readers
  at once. Business reader: value/why. Technical reader: how/constraints/NFRs. AI
  agent: explicit, self-contained, no-assumption, machine-parseable.
- Acceptance criteria format: Given/When/Then (Gherkin) and/or measurable bullets;
  specific, testable. AI agents do NOT fill gaps with assumptions, so edge cases,
  error/offline states, and non-functional (perf/security) must be explicit.
- Machine-readable metadata: theme/epic link, status, priority, labels, file/path
  hints, links to API/OpenAPI/architecture, Definition of Done.
- Generic across trackers (don't bake in Jira). How an idea/`FINAL.md` maps to an
  epic + stories + subtasks, and the "tracker is a MIRROR, files stay canonical" rule.

## Reference grounding (from a June-2026 web pass — extend, don't just repeat)

- Worktrees-for-AI: file isolation converts invisible runtime corruption into normal
  git conflicts; decompose by feature boundary; shared task doc for claim/done;
  staging branch before main; `.env.local` per worktree; sparse-checkout to scope
  files; native worktree support shipped in Claude Code/Cursor/Copilot CLI + JetBrains
  2026.1 / VS Code Jul-2025. (mindstudio.ai, augmentcode.com, upsun, parallelcode.app)
- Tickets-for-AI: precision controls output quality; AI doesn't assume gaps → spell
  out edge/NFR; Given-When-Then; metadata-rich stories (theme/epic/release/status/
  priority); 2026 RAG approach indexes OpenAPI/architecture so AI references real
  endpoints; human oversight stays essential; INVEST. (mountaingoatsoftware,
  handsonai.info, infoworld, vegavid, clickup)

## Constraints

- English only. Both skills must be **vendor/tracker/model-agnostic** (we just
  audited the core for genericity — keep that bar).
- Truthful, concrete: templates and commands must be usable as-is.
- Respect the core+plugins design: thin seam in parley-deck, mechanics in the skills.

## Non-goals

- Not writing the skills' `SKILL.md` yet (this idea produces the design that lets us).
- No protocol edits here (any parley-deck carve-out is a separate meta-protocol-change).
- No vendor lock-in (no "Jira-only" assumptions).
