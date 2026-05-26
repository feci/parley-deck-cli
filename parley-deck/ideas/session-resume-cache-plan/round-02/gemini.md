---
agent: gemini
idea: session-resume-cache-plan
round: 2
date: 2026-05-25
responding-to: [codex/round-01, claude/round-01, hermes/round-01]
---

## Position changes since round 1

I am pivoting from my initial proposal of a central, probe-discovered `agent-registry.json`. The implementation of slice 1 and the arguments from other participants have convinced me that a more integrated and incremental approach is better. My focus now is on enriching the existing run artifacts and agent configuration, rather than creating a parallel source of truth.

The core of my proposal now aligns with the incremental slicing plan that has emerged: first, we durably capture the *inputs* (prompts, config) for each agent invocation. Second, we capture the *outputs* (native session handles) opportunistically. This makes the system resilient and avoids relying on brittle discovery mechanisms like parsing help text, a concern I raised in round 1.

## Responses to others

### @codex

I agree entirely with the incremental roadmap you've proposed. It provides a clear, logical path forward and correctly prioritizes foundational work. Treating slice 1 as a stable base is the right move.

Your point about avoiding a parallel "agent-registry.json" is well-taken. Leveraging the existing Parley runtime agent configuration (`headless-agents.local.json`) and augmenting it with per-run attempt records is a much cleaner architecture. I've integrated this into my current proposal. I also concur that roster and other protocol hygiene issues should be handled in separate, dedicated ideas.

### @claude

Your distinction between promising step-level *workflow recovery* versus best-effort *conversation continuity* is crucial. This framing should guide all user-facing communication about this feature. It perfectly addresses the "context drift" concerns from my first round. The concept of a rebuildable local cache, which you and Hermes both raised, is sound, and the per-run artifacts in the workspace are the key to making that possible.

### @hermes

Your model of a central index (`~/.parley-deck/runs.json`) plus per-run manifests is a strong one. The implemented slice 1 gives us the per-run `run.json`, which serves as the manifest. The next logical step is to build out the "attempt" records you described, including prompt hashing for integrity. I agree with codex that automatic stale-heartbeat detection is a feature for a later slice, after we have reliable capture of agent invocations.

## New concerns / questions

-   **Capability Storage:** If we abandon a central registry, where do we store discovered agent capabilities (e.g., the flag to resume a session)? I propose we add an optional `capabilities` block to the agent definition within `headless-agents.local.json`. This could be populated manually or by a future, more robust probe tool, but it keeps the configuration self-contained.
-   **Prerequisite:** As noted by codex, creating a project-local `headless-agents.local.json` is a clear prerequisite for the next slice. Without it, there is no stable configuration to rely on for recovery.

## Current proposal

I endorse the incremental roadmap and propose the following details for the next slices:

**Slice 2: Per-Agent Attempt & Prompt Capture**
For each agent invocation, the facilitator will write a canonical "attempt record" to `parley-deck/runs/<run-id>/attempts/<round>-<agent-id>-<timestamp>.json`. This makes every action independently auditable and replayable.

The record will contain:
-   `promptHash`: A SHA256 hash of the rendered prompt for integrity checks.
-   `prompt`: The full, rendered prompt text sent to the agent.
-   `config`: The launch configuration (model, profile, etc.).
-   `cli`: The exact command-line invocation.
-   `expectedArtifact`: The path to the file the agent was supposed to create.
-   `result`: A structure containing the process exit code, duration, and truncated stdout/stderr.

**Slice 3: Opportunistic Native Handle Capture**
This builds directly on Slice 2.
-   The agent's definition in `headless-agents.local.json` can be augmented with a `capabilities.resume` block, specifying how to request a session handle (e.g., a `--session-out-file` flag) and how to use it later (e.g., a `--resume-session-id` flag).
-   When invoking an agent, the facilitator checks for this capability. If present, it uses it.
-   Any captured `nativeSessionID` is stored in the "attempt record" from Slice 2.
-   This makes capability discovery a runtime concern based on explicit configuration, not brittle parsing.

**Slice 4 & 5:** I agree with the definitions from codex for cache rebuild/rebind and TUI integration, which would consume the artifacts generated in Slices 2 and 3.
