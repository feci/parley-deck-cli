---
agent: hermes
idea: embedded-default-protocol-resync
round: 1
date: 2026-06-13
---

## Summary
The embedded default is the bootstrap template delivered by `parley init`. It must produce a usable, non-misleading starting COOPERATION.md for any new project. Currently it ships parley-deck's own roster, workspace name, and lacks §12. The single allowed mutation (transport swap) must remain the only runtime change.

## Proposed approach
1. §12 propagation: carry the entire section verbatim, including the closing "ratified by idea …" provenance line. It is already-ratified protocol content.
2. Genericization: make only the header lines that are obviously per-project (`Workspace:`, `Created:`, `Protocol synced:`) into placeholders or omit the synced line; replace the §2 roster table and host-handle table with a single-row placeholder example plus the existing prose note that "The roster is project-specific." Everything else (all phases, invariants, §11 mechanics) stays byte-identical.
3. Anti-drift: add a Go test in the workspace package that asserts the embedded default equals the live deck modulo an explicit allowlist of the genericized header/roster lines. This ships with the idea and is the enforceable guard.

## Concerns / open questions
- How verbose should the placeholder roster example be? One illustrative row or an empty table?
- Does the `Protocol synced:` line belong in the template at all, or only in the live deck after the first skill sync?

## Risks
- If the test allowlist grows, drift detection weakens. Keep the allowlist tiny (header + §2 tables only).
- Over-genericizing the template could make the first `parley init` experience feel sparse; balance with the existing prose guidance already present in §2.