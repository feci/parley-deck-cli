---
agent: hermes
idea: consensus-workflow-cli
round: 1
date: 2026-05-12
---

## Summary
First slice adds CLI commands to advance from round-01 artifacts to a consensus draft and sequential signoffs. Source of truth stays in parley-deck/ files. CLI uses templates for deterministic structure, GitHub PR transport remains a mirror.

## Proposed approach
- Add `parley consensus draft` that reads round files, applies a fixed template for design-consensus.md, and writes it under the idea.
- Add `parley signoff add` that appends the current agent's signed block (timestamp + hash of prior content) to the consensus file.
- Validation: require non-empty prior round files, enforce one signoff per agent, reject unsigned changes.
- Trade-off: no content understanding or auto-generation of consensus text; agents must supply prose. Keeps scope small and verifiable.
- Reuse existing store and protocol packages; no new TUI for first slice.

## Concerns / open questions
- How to handle merge conflicts on concurrent signoff appends without full Git automation in slice 1.
- Whether review-consensus primitives can reuse the same signoff block format without extra fields.

## Risks
- Agents may bypass CLI and edit files directly, breaking the "each writes its own artifact" rule. Mitigation: hash checks on signoff.
- Overly narrow slice may leave users hand-editing PR descriptions; acceptable for v1.