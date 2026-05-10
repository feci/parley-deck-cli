---
agent: hermes
idea: hermes-integration-smoke-test-verify-that-hermes-can-participate-in-a-parley-deck-round-by-writing-its-own-round-01-artifact-only-2
round: 1
date: 2026-05-10
---

## Summary
Smoke test: Hermes writes only its round-01 artifact using local-directory transport to verify Parley Deck participation.

## Proposed approach
Use write_file on exact target path with required YAML + sections. No other files touched. Human-in-loop default.

## Concerns / open questions
- Does directory pre-exist or must Hermes create parents?
- How to handle auto vs human mode without extra artifacts.

## Risks
- Accidental overwrite blocked per rule (trade-off: safety over speed).
- Scope creep if non-goals ignored (mitigated by strict "only this file").