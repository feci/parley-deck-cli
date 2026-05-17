---
agent: hermes
idea: roadmap-implementation-plan
round: 2
date: 2026-05-17
---

## Summary
Accept with reservations. The convergence on a deterministic round-index + artifact-pruning first slice aligns with operations/resilience priorities: it produces a reusable, local artifact that later error classification, hooks, and memory layers can consume without model calls or external deps.

## Position on proposed roadmap
Accept the revised order (round-index-artifact-pruning before repo-map). It correctly prioritizes context hygiene that benefits every subsequent resilience and tooling slice. Keep the opaque-CLI fence and no direct SDK in the first sequence.

## Required changes before consensus
- Clarify that `_index.md` and sanitized context are generated artifacts, not canonical participant files; they must be git-tracked but clearly labeled "derived".
- Add explicit test that index generation failure never aborts a successful round (warning only).
- Confirm the exact hidden-reasoning fence patterns the sanitizer will support in slice 1.

## Risks
- Index becomes stale if round artifacts are manually edited post-run. Mitigation: treat index as runner-owned and regenerate on demand.
- Over-trimming in sanitizer removes useful operational signals (e.g., rate-limit metadata). Keep the supported block list small and explicit.