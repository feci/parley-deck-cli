---
idea: meta-protocol-change-agent-teams-patterns
author: codex
created: 2026-05-14
participants: [codex, claude, gemini, hermes]
status: final
---

## Problem / idea

Evaluate whether Parley Deck's cooperation protocol can borrow useful process patterns from Claude Code Agent Teams and related "agentic teams" practice, then propose what should change and how.

The user asked in Slovak to use Parley Deck, compare our cooperation protocol with Agent Teams, and propose improvements if they would make the process better.

## Source context

- Current project protocol: `parley-deck/COOPERATION.md`.
- Claude Code docs, "Orchestrate teams of Claude Code sessions": `https://code.claude.com/docs/en/agent-teams`.
- Florian Bruniaux guide, "Agent Teams Workflow": `https://github.com/FlorianBruniaux/claude-code-ultimate-guide/blob/main/guide/workflows/agent-teams.md`.
- Heeki Park, "Collaborating with agents teams in Claude Code": `https://heeki.medium.com/collaborating-with-agents-teams-in-claude-code-f64a465f3c11`.

Key source observations to consider:

- Agent Teams use a team lead plus independent teammates, shared task list, direct inter-agent messaging, and local task/team state.
- The official docs emphasize that Agent Teams are experimental, token-intensive, and best for parallel research/review, separate feature slices, competing hypotheses, and cross-layer work with clear boundaries.
- The docs contrast Agent Teams with subagents: teams are better when workers need to coordinate and challenge each other, while simpler delegation is better for focused tasks where only the result matters.
- Useful practices include explicit teammate names, bounded team size, enough task-specific context, plan approval before risky implementation, monitoring/steering, quality-gate hooks, and avoiding file conflicts.
- Documented limitations include coordination overhead, lagging task state, shutdown/resume limitations, fixed lead ownership, permission-mode inheritance, and the need for manual cleanup.

## Constraints

- Do not edit `COOPERATION.md` directly in this idea unless consensus reaches a final protocol-change plan and the user explicitly approves implementation.
- Preserve Parley Deck's vendor-neutral design. Do not make Claude Agent Teams a required dependency.
- Preserve canonical file ownership: each participant writes its own protocol artifact.
- Preserve current transport: `github-pr`; canonical files remain under `parley-deck/`.
- Keep any proposed changes small enough to implement incrementally.
- English-only for all files under `parley-deck/`.

## Non-goals

- Do not implement code in this idea.
- Do not replace Parley Deck's round/consensus model with Claude Agent Teams.
- Do not add hidden autonomous execution or unbounded agent spawning.
- Do not rely on experimental Claude-only features for baseline protocol correctness.
