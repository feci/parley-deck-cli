---
idea: roadmap-implementation-plan
finalized-by: codex
date: 2026-05-17
consensus: consensus.md
status: final
---

## Decision

Use `parley-deck-ideas.md` as a roadmap source, but deliver the work as small Parley Deck implementation slices.

The current `parley-deck-cli` 1.1.x baseline already includes the bootstrap orchestration surface: init, runtime discovery/config, headless and interactive launch modes, isolated Gemini/Hermes homes, live TUI, HITL questions, status/resume, consensus signoffs, semantic version reporting, and parallel round-one headless execution.

Provider prompt caching and semantic LLM cache are deferred until a direct SDK/API backend exists. The current opaque CLI invocation path cannot reliably set provider cache metadata or inspect provider cache keys.

## Implementation Roadmap

1. `round-index-artifact-pruning`: deterministic round index plus context-only artifact sanitizer.
2. `repo-map-mvp`: deterministic repo map with file tree plus Go symbol extraction using Go standard-library parsing first.
3. `context-pack-wiring`: use round indices, sanitized artifacts, repo maps, and strict budgets in later round/review/signoff prompt builders.
4. `error-classifier-redaction`: classify agent failures and add explicit redaction helpers for generated context.
5. `hooks-mvp`: add conservative pre/post run hook points after typed errors and context artifacts exist.
6. `local-memory-insights`: build opt-in local insights over round indices and finalized ideas.
7. `direct-sdk-backend-pilot`: conditional future slice for provider prompt caching, semantic cache, and direct provider features.

## First Slice: round-index-artifact-pruning

Implement the first slice next.

Required behavior:

- Generate a derived, git-tracked `round-NN/_index.md` after a round completes.
- Include one row or section per participant, including skipped or failed participants.
- Extract compact context from existing artifacts without mutating those artifacts.
- Add a context-only sanitizer for a closed list of hidden-reasoning fences:
  - `<think>...</think>`
  - `<thought>...</thought>`
  - `<thinking>...</thinking>`
- Clearly document that the sanitizer is not secret redaction.
- Keep `_index.md` runner-owned and derived, not participant-owned.
- Keep `_index.md` byte-deterministic for fixed inputs. Avoid wall-clock timestamps in generated content.
- Limit section extraction to H2 headings (`##`) for the first version.
- Include a deterministic approximate token estimate in index metadata using a simple documented heuristic.
- If index generation fails after participant artifacts are written, the round should still be successful and should surface an explicit warning/result.
- If an artifact lacks recognized sections, the index should degrade gracefully rather than failing.

## Tests Expected

- Unit tests for sanitizer behavior on supported fences, malformed blocks, and content that should remain unchanged.
- Unit tests or golden tests proving deterministic `_index.md` output for fixed artifact inputs.
- Tests covering skipped and failed participants in the index.
- Tests proving index generation failure does not abort an otherwise successful round.
- Regression test that source artifact files are not modified by sanitizer/index generation.

## Notes For Implementation PR

- Present the index writer and sanitizer as two separately reviewable behaviors in the PR description.
- Keep repo-map work out of the first slice. It is the next roadmap item.
- Keep secret redaction out of the first slice. It belongs to `error-classifier-redaction`.
