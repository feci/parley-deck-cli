---
agent: hermes
idea: meta-protocol-change-agent-teams-patterns
round: 1
date: 2026-05-14
---

## Summary
Parley Deck COOPERATION.md uses flat round-based markdown artifact exchange. Claude Code Agent Teams use isolated subagent delegation, persistent memory per context, and explicit skill separation for complex tasks.

## Proposed approach
Adopt optional delegation pattern vendor-neutrally: add "delegate" round variant to COOPERATION.md allowing isolated child sessions that write only to their assigned canonical artifact. Keep single owner per artifact. Trade-off: +15% setup cost for better isolation on multi-hour tasks.

## Concerns / open questions
Does delegation require new top-level coordinator artifact? How to merge child outputs without violating "do not edit other agents" rule?

## Risks
Complexity creep that dilutes Parley Deck's minimalism; risk of nested artifact ownership disputes if not strictly scoped.