---
agent: gemini
idea: roadmap-implementation-plan
round: 1
date: 2026-05-17
---

## Summary
The current `parley-deck-cli` provides a robust foundation for multi-agent cooperation, but as rounds progress, context bloat and "exploratory noise" will degrade performance and increase costs. As the **roadmap prioritization and context-efficiency lens**, I propose a sequence that prioritizes high-density context and structural awareness. My primary goal is to ensure that agents receive the most relevant information in the most compact form possible.

## Proposed approach

### 1. Assessment of parley-deck-ideas.md
- **Complete**: Parallel headless execution, `parley resume` (checkpointing), isolated agent homes.
- **Partial**: Skills integration (basic runtime exists), Local memory (isolated homes provide basic file persistence).
- **Missing**: Repo maps, context compression (scrubbing/summarization), diff fences, semantic cache, credential guards, sub-agent delegation.
- **Not Feasible (Current Architecture)**: Direct provider prompt caching (the current use of external CLIs prevents passing the required cache control metadata; this would require migrating to direct SDKs or specialized wrappers).

### 2. Ordered Implementation Roadmap
1.  **Context Efficiency I: Artifact Scrubbing & Pruning** (Immediate token savings by removing internal reasoning and redundant metadata from history).
2.  **Structural Context: Repo Maps** (Reduces exploratory file-reading turns by providing a deterministic map of the workspace).
3.  **Multi-agent Orchestration: Sub-agent Delegation Pattern** (Scaling complexity by allowing agents to delegate specialized sub-tasks via the runner).
4.  **Context Efficiency II: Trajectory Compression** (Recursive summarization of historical rounds to maintain a long-term context foundation).
5.  **Operational Resilience: Rate & Credential Guards** (Protection against runaway costs and leaked secrets).
6.  **Memory: Structured Hierarchy & Insights Provider** (Persisting cross-run facts beyond isolated home directories).

### 3. First Implementation Slice: Artifact Scrubbing & Pruning
The first slice implements a post-processing filter in the `runner` that scrubs `<thought>` blocks and redundant frontmatter from previous round artifacts before they are injected into the prompt for subsequent rounds.

- **Objective**: Maximize the effective context window for problem-solving by stripping "internal monologue" that peers do not need for technical consensus.
- **Implementation**: Add a `PruneArtifact(content string) string` utility to `internal/runner` and apply it when building multi-round prompts.

### 4. Slice Details
- **Non-goals**: Modifying the original artifacts on disk; implementing model-side summarization (reserved for Slice 4).
- **Tests**: Unit tests for regex-based scrubbing; integration test verifying that a Round 2 prompt excludes the `<thought>` section of a Round 1 artifact.
- **Checks**: Ensure scrubbing does not break Markdown structure or valid frontmatter required for protocol tracking.

## Concerns / open questions
- **Opaque CLI Outputs**: If an agent CLI does not wrap its reasoning in standard tags (like `<thought>`), scrubbing will be unreliable. We may need to update agent instructions to mandate these tags.
- **Skill Versioning**: We should verify if the pending `parley-deck-skill` 1.1.1 update introduces any output format changes that would conflict with this scrubbing logic.

## Risks
- **Context Loss**: Information intended for peers might be lost if it is accidentally included inside a scrubbed block.
- **Formatting Fragility**: Regex-based scrubbing can be brittle if agents produce malformed Markdown or unexpected nested blocks.
