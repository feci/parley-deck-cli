---
idea: consensus-workflow-cli
author: codex
created: 2026-05-12
participants: [codex, claude, gemini, hermes]
status: consensus
---

## Problem / idea

Design the next `parley-deck-cli` roadmap slice: CLI support for consensus and signoff workflows.

The current product can initialize a workspace, discover/verify agent runtimes, run round-01 across selected agents, show a live TUI, handle HITL questions/answers, and rebuild status/resume views from durable run state. The remaining roadmap item 1 is:

- Consensus/signoff workflow in the CLI: draft and validate design consensus, collect sequential signoffs, and support review consensus automation.

This idea should define the product and implementation scope for the first shippable slice that lets users move from round files to consensus/finalization without manually hand-crafting every protocol step. It should respect the Parley Deck canonical filesystem protocol and the active `github-pr` coordination transport used by this repository.

## Constraints

- Keep the source of truth in `parley-deck/` files. PR comments, labels, or reviews are mirrors only.
- Use Go and the existing architecture: `internal/protocol`, `internal/runner`, `internal/store`, `internal/tui`, `internal/app`, and the existing agent runtime config.
- Preserve the non-solo requirement: the CLI may facilitate agent signoff, but each participant must still write or append its own canonical artifact.
- Keep the first slice small enough to implement and review safely after design consensus.
- Prefer deterministic validation and templates before model-generated prose. The CLI should not pretend to understand consensus content more deeply than it can prove.
- Support design consensus first, and include review consensus only if it can share the same primitives without broadening the slice too far.
- Maintain compatibility with the current GitHub PR transport, but do not require full GitHub API automation unless the agents agree it is the right first slice.
- English-only for all protocol artifacts and PR text.

## Non-goals

- Full autonomous auto mode.
- GitLab MR automation.
- Release packaging.
- Token/cost telemetry.
- Retrying, supervising, or reattaching live agent processes.
- Replacing human or agent judgment about whether the consensus text is substantively correct.
