---
agent: hermes
idea: roadmap-implementation-plan
round: 1
date: 2026-05-17
---

## Summary
Focus on operations, resilience, tooling. Current CLI already delivers workspace init, parallel headless runs, TUI, resume, status, HITL, and version surfaces. Missing: structured error classification, credential/rate guards, redaction, checkpointing hooks, local memory providers, and deterministic skill execution paths. Prompt caching and direct SDK calls are non-goals for first slices; keep Go + Bubble Tea.

## Proposed approach
1. Add error classifier + redaction hook (deterministic, local).
2. Credential and rate-limit guard layer with local config.
3. Resume checkpointing for long runs (file-based, no external deps).
4. Local memory provider scaffold (JSONL + simple index).
5. Tooling: parley doctor and hook registry commands.
Trade-off: deterministic local features first; model-dependent automation deferred until after these slices. Small slices ensure one-PR review cycles.

## Concerns / open questions
How to surface rate-guard decisions in TUI without noise? Will hooks be opt-in only or default-on for resilience?

## Risks
Over-engineering guards before real usage data; scope creep into provider SDKs (explicit non-goal). Test: run 10 parallel agents with injected rate errors and verify clean resume.