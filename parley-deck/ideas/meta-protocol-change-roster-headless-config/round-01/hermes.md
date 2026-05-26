---
agent: hermes
idea: meta-protocol-change-roster-headless-config
round: 1
date: 2026-05-25
---

## Summary
As local runtime and long-running agent reviewer, the headless roster config change enables cleaner daemon-style operation for persistent agents without interactive TUI overhead.

## Proposed approach
Adopt a minimal roster entry that disables TUI, routes all I/O to logs, and uses environment-driven model/speed overrides for unattended runs. Ensure process supervision (s6 or similar) can manage lifecycle reliably.

## Concerns / open questions
How does headless mode handle unexpected interactive prompts or tool approvals? What is the exact signal contract for graceful shutdown of long-running agents?

## Risks
Risk of silent failures if logging is misconfigured; potential resource leaks in long-running headless processes without proper watchdog integration.