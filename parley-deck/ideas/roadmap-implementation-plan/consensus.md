---
idea: roadmap-implementation-plan
drafted-by: codex
date: 2026-05-17
---

## Agreed decisions

- Use the new `parley-deck-ideas.md` catalogue as a roadmap source, but deliver it as small Parley Deck implementation slices.
- Treat current `parley-deck-cli` 1.1.x as already having the bootstrap orchestration surface: init, runtime discovery/config, headless and interactive launch modes, isolated Gemini/Hermes homes, live TUI, HITL questions, status/resume, consensus signoffs, semantic version reporting, and parallel round-one headless execution.
- Do not prioritize provider prompt caching or semantic LLM cache until a direct SDK/API backend exists. The current opaque CLI invocation path cannot reliably set provider cache metadata or inspect provider cache keys.
- Choose `round-index-artifact-pruning` as the first implementation slice.
- Keep the first slice deterministic and local:
  - Generate a derived, git-tracked `round-NN/_index.md` after a round completes.
  - Include one row/section per participant, including skipped or failed participants.
  - Extract compact context from existing artifacts without mutating those artifacts.
  - Add a context-only sanitizer for a closed list of hidden-reasoning fences.
- The sanitizer is not secret redaction. Secret redaction belongs to a later `error-classifier-redaction` slice.

## Ordered implementation roadmap

1. `round-index-artifact-pruning`: deterministic round index plus context-only artifact sanitizer.
2. `repo-map-mvp`: deterministic repo map with file tree plus Go symbol extraction using Go standard-library parsing first.
3. `context-pack-wiring`: use round indices, sanitized artifacts, repo maps, and strict budgets in later round/review/signoff prompt builders.
4. `error-classifier-redaction`: classify agent failures and add explicit redaction helpers for generated context.
5. `hooks-mvp`: add conservative pre/post run hook points after typed errors and context artifacts exist.
6. `local-memory-insights`: build opt-in local insights over round indices and finalized ideas.
7. `direct-sdk-backend-pilot`: conditional future slice for provider prompt caching, semantic cache, and direct provider features.

## First slice guardrails

- Supported sanitizer fence patterns are a closed initial set: `<think>...</think>`, `<thought>...</thought>`, and `<thinking>...</thinking>`. New patterns require a later change.
- Sanitization is context-only and must never rewrite source artifacts on disk.
- `_index.md` is a runner-owned derived artifact, not a participant-owned file.
- `_index.md` should be byte-deterministic for fixed inputs. Avoid wall-clock timestamps in the generated body.
- Limit section extraction to H2 headings (`##`) for the first version.
- Include a deterministic approximate token estimate in index metadata using a simple documented heuristic.
- If index generation fails after participant artifacts are written, the round should still be considered successful and should surface an explicit warning/result.
- If an artifact lacks recognized sections, the index should degrade gracefully rather than failing.

## Agreed trade-offs

- Start with round indices and sanitizer before repo maps because every later context, review, and memory slice can consume compact round artifacts.
- Keep repo maps as the second slice; they remain important but are more useful after a context substrate exists.
- Defer hooks and memory until generated context artifacts and typed failure surfaces exist.
- Keep direct provider SDK work conditional because it has a higher migration cost and changes the billing/control model.

## Tests expected for slice 1

- Unit tests for sanitizer behavior on supported fences, malformed blocks, and content that should remain unchanged.
- Unit tests or golden tests proving deterministic `_index.md` output for fixed artifact inputs.
- Tests covering skipped/failed participants in the index.
- Tests proving index generation failure does not abort an otherwise successful round.
- Regression test that source artifact files are not modified by sanitizer/index generation.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The plan picks a small deterministic first slice, keeps provider-specific caching fenced until a direct backend exists, and preserves repo-map as the next delivery after the context substrate.
