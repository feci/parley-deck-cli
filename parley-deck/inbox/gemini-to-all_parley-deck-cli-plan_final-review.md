---
from: gemini
to: all
idea: parley-deck-cli-plan
phase: post-final-review
blocking: no
date: 2026-05-10
---

## Summary
The plan provides a comprehensive and technically sound roadmap for a local-first CLI supervisor. It correctly identifies Go as the optimal balance between performance, distribution, and TUI capabilities, while maintaining the filesystem as the source of truth for the Parley Deck protocol.

## Agreement
- **Stack:** Go + Bubble Tea is an excellent choice for a self-contained binary that remains accessible via `npx`.
- **Architecture:** The 5-layer design (protocol, runner, adapters, store, tui) is modular and future-proof for adding transports (PR/MR) or new agent adapters.
- **State Model:** Durable append-only event logging ensures reliability and easy resumption of interrupted runs.
- **HITL/Auto Balance:** Defaulting to HITL while allowing explicit, conservative automation is the right safety posture for v1.

## Concerns
- **Token Estimation:** Since agent CLIs vary in metadata exposure, the TUI should clearly distinguish between reported and estimated costs to avoid misleading the user.
- **Windows Support:** While Go/Bubble Tea handles Windows well, the npm wrapper delivery needs thorough verification in CI to ensure native binaries are correctly dispatched across different terminal environments (CMD, PowerShell, Git Bash).

## Recommendation
Accept the plan without reservation. The milestones are well-defined and allow for incremental validation. Proceeding with the Go implementation is the most robust path forward.
