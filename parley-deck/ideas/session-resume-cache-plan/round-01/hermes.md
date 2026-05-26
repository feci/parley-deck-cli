---
agent: hermes
idea: session-resume-cache-plan
round: 1
date: 2026-05-25
---

## Summary
Registry model: central durable index under ~/.parley-deck/runs.json + per-run manifest.json that maps run-id to native agent session pointers and step status. Recovery reconciles filesystem state against the registry, preferring native CLI resume when available and falling back to prompt re-invoke. Canonical artifacts stay in workspace parley-deck/; local state is rebuildable from run.json + agents.json.

## Proposed approach
Store registry as ~/.parley-deck/runs.json (index) and runs/<run-id>/{manifest.json, heartbeat.json, agents/<round>.<agent>.json}. Each manifest entry records {step_id, native_session_id, session_path, prompt_hash, status}. On resume: load registry, detect stale heartbeats (>30s), reconcile by checking expected artifact paths, then invoke native --resume or fresh session with cached prompt. Trade-off: native continuity is best-effort (CLIs may prune jsonl); registry guarantees step-level restart without full idea replay. Rebuild command walks parley-deck/runs/*/run.json to restore index after cache loss. No credentials stored; capability detection per adapter.

## Concerns / open questions
- How to handle concurrent TUI + CLI access to same registry entry without races.
- Registry schema evolution: when to introduce sqlite vs stay JSONL.
- Rebind logic for moved workspaces: how aggressively to update absolute paths.

## Risks
- Registry and native sessions diverge if CLI prunes aggressively or user clears ~/.parley-deck.
- Stale registry entries after machine move or long idle periods.
- Over-reliance on advisory locks on shared filesystems.