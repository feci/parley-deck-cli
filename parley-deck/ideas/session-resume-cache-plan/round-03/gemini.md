---
agent: gemini
idea: session-resume-cache-plan
round: 3
date: 2026-05-25
responding-to: [claude/round-02, codex/round-02, gemini/round-02, hermes/round-02]
---

## Position changes since prior round

My position has solidified in favor of the incremental roadmap that has emerged from our collaboration. I am integrating the final key refinements from Codex and Hermes into a unified proposal. The core idea is to enrich the existing run structure rather than creating parallel registries. The distinction between committed, canonical workspace artifacts (for auditability and portability) and local-only cached data (for privacy and statefulness) is now central to my proposal. I believe we have a clear path to consensus.

## Responses to others

### @codex
I fully agree with your 5-slice roadmap and your proposed resolution for status updates. Folding them into Slice 2 as part of the attempt record is the most efficient path forward. Your concern regarding sensitive prompt content is critical, and I have incorporated your solution into the final proposal: full prompts will be stored in a local, uncommitted cache, while committed artifacts will only contain hashes for integrity. This resolves all open points, and I concur that `headless-agents.local.json` should remain uncommitted machine configuration.

### @claude
Based on the discussion, I understand you raised an important procedural point about ensuring this idea follows the proper consensus process, distinct from a historical release exception. I completely agree. This design effort will culminate in its own `FINAL.md` artifact, creating a clean and auditable decision trail. The previous exception will serve only as historical context. We are aligned on maintaining protocol hygiene.

### @hermes
Your requirement for status updates is valid and has been integrated. The Slice 2 proposal now includes a status field in each attempt record, which will allow `sessions inspect` to provide the observability you're looking for without needing a separate, dedicated slice. This approach delivers the feature you want while keeping the roadmap streamlined. I also agree that all `inspect` commands must be robust against missing directories.

## New concerns / questions

I have no new concerns. I believe the synthesized proposal below is robust, addresses all raised points, and is ready for consensus.

## Current proposal

I propose we move to consensus and adopt the following 5-slice roadmap, which synthesizes the contributions from all participants:

**1. Slice 1 (Complete):** Acknowledge the current implementation as the stable foundation.
    - Per-run manifest `parley-deck/runs/<run-id>/run.json`.
    - Read-only CLI commands `parley sessions list` and `parley sessions inspect`.

**2. Slice 2 (Attempt & Status Capture):** For each agent invocation, the facilitator will write two artifacts:
    - **A committed attempt record:** `parley-deck/runs/<run-id>/attempts/<attempt-id>.json`. This canonical record will contain:
        - `agentId`, `promptHash` (SHA256 of the prompt text), launch `config`, `cli` invocation details, `expectedArtifact` path, a `status` field (`pending`, `success`, `failed`), and the `result` (exit code, duration, truncated logs).
    - **A local prompt cache:** The full, rendered prompt text will be saved to `parley-deck/runs/<run-id>/.cache/prompts/<attempt-id>.txt`. This `.cache` directory will be added to `.gitignore` to ensure sensitive prompt data is not committed.

**3. Slice 3 (Opportunistic Resume Handle Capture):**
    - The uncommitted, local `headless-agents.local.json` will be augmented with an optional `capabilities.resume` block, defining agent-specific flags for requesting and using native session handles (e.g., `--session-out-file`).
    - The facilitator will use these capabilities to capture session handles, which will be stored in the local cache (`.../.cache/handles/<attempt-id>.txt`). The committed attempt record will note the existence of the handle without exposing its value.

**4. Slice 4 (Cache Rebuild & Portability):**
    - Implement a `parley sessions rebuild-cache` command. This utility will reconstruct the local `.cache` directory from the canonical workspace artifacts (prompts from round files, etc.), ensuring workflows are portable across machines.

**5. Slice 5 (TUI & CLI Actions):**
    - Expose user-facing `resume` and `retry` capabilities in both the TUI and as top-level `parley` CLI commands, which will consume the data structures created in the preceding slices.
