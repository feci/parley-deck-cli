---
agent: gemini
idea: roadmap-implementation-plan
round: 2
date: 2026-05-17
---

## Summary
I accept the facilitator's convergence proposal. Combining deterministic round indexing with artifact pruning addresses the most immediate threat to multi-round agent performance: context bloat and exploratory noise. This foundation is essential for the later, more complex slices in the roadmap.

## Position on proposed roadmap
Accept. The reordered roadmap correctly identifies that context hygiene is the primary bottleneck. By delivering `round-index-artifact-pruning` first, we provide a clean, high-density substrate for both humans and future agents to reason about the project trajectory.

## Required changes before consensus
- **Deterministic Token Estimator**: Slice 1 should include a very basic, deterministic character-to-token heuristic (e.g., 4 chars/token) in the index metadata. While not perfectly accurate, it provides a stable signal for "cost awareness" in later context-packing slices.
- **Section Extraction depth**: Explicitly limit section extraction in the index to H2 (`##`) headings to prevent the index itself from becoming bloated in highly structured artifacts.

## Risks
- **Over-Sanitization**: The pruning of reasoning blocks (`<thought>`, etc.) must be strictly limited to context-injection time. We must ensure the original artifacts remain untouched for debugging and auditability.
- **Heuristic Fragility**: As agents evolve, their internal reasoning formats may change. The sanitizer should be designed as a registry of patterns that can be easily extended without re-architecting the runner.
