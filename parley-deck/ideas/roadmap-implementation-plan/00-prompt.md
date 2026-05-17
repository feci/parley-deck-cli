---
idea: roadmap-implementation-plan
author: user
created: 2026-05-17
participants: [codex, claude, gemini, hermes]
status: round-02
roles:
  codex: facilitator and implementation-sequencing lens
  claude: architecture, UX, and risk lens
  gemini: roadmap prioritization and context-efficiency lens
  hermes: operations, resilience, and tooling lens
---

## Problem / idea

Turn `/Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-ideas.md` into an implementation plan for `parley-deck-cli`, then select the next small slice to implement first.

The source catalogue highlights ideas from Hermes Agent and related tools:

- Token and context efficiency: Anthropic prompt caching, repo maps, trajectory/context compression, think-block scrubbing, diff fences, semantic cache.
- Memory: local memory providers, insights, richer memory hierarchy, user/project modeling.
- Multi-agent orchestration: parallel execution, batch execution, sub-agent delegation, roles, skills, microagents.
- Operations and resilience: credential/rate guards, error classifiers, redaction, hooks, checkpointing/resume.
- Extensibility: skills, MCP or server hosting, optional advanced cognition.

Current project state to consider:

- `parley-deck-cli` is on `main`, clean, and reports `parley 1.1.0`.
- Released tags currently include `v1.0.0` and `v1.1.0`.
- Project metadata and protocol status are valid.
- The local installed `parley-deck-skill` runtime is still `1.1.0`; a newer `1.1.1` skill package exists but this is not a blocker for CLI planning.
- Already implemented: workspace init, agent runtime config, headless and interactive agent invocation modes, isolated Gemini/Hermes homes, live run TUI, HITL questions, `parley status`, `parley resume`, consensus signoff requests, semantic version reporting, and basic release/install surfaces.
- Round-one headless agent execution is already parallelized in the runner.

## Requested output

Each participant should propose:

1. What from `parley-deck-ideas.md` is already complete, partially complete, missing, or not feasible in the current CLI-backed architecture.
2. An ordered implementation roadmap of small reviewable slices.
3. The single first slice that should be implemented next.
4. Risks, non-goals, and tests/checks for that first slice.

The consensus should produce a concrete sequence that can be delivered through Parley Deck implementation PRs one after another.

## Constraints

- Keep slices small enough for one focused implementation and review cycle.
- Prefer deterministic local functionality before model-dependent automation.
- Do not claim provider prompt caching support unless the current invocation path can actually control provider cache metadata. Opaque external CLIs may make this infeasible until a direct API/SDK backend exists.
- Do not add speculative architecture for future tools unless it unblocks the first one or two slices.
- Preserve the current Go CLI and Bubble Tea TUI direction.
- Use English for all canonical Parley Deck artifacts.
- This design idea does not implement code directly; implementation starts after consensus/finalization selects the first slice.

## Non-goals

- No provider billing bypass.
- No direct Anthropic/OpenAI/Gemini SDK migration in the first slice unless participants agree it is the smallest safe step.
- No sandboxed execution, hosted MCP server, voice interface, federation, or knowledge graph in the first sequence.
