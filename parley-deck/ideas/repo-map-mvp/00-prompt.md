---
idea: repo-map-mvp
author: user
created: 2026-05-17
participants: [codex, claude, gemini, hermes]
roles:
  codex: implementation and CLI UX lens
  claude: architecture and reviewability lens
  gemini: context-efficiency and output-shape lens
  hermes: operations and resilience lens
status: consensus
---

## Problem / idea

Implement the next roadmap slice from `roadmap-implementation-plan/FINAL.md`: `repo-map-mvp`.

Parley Deck needs a deterministic local repository map that agents and humans can inspect before later slices wire compact repository context into prompts. The MVP should reduce exploratory file-reading cost without taking on a tree-sitter/aider clone.

## Constraints

- Keep the first repo-map slice small and reviewable.
- Use Go standard-library parsing for Go symbols first. Do not add tree-sitter or external parser dependencies in this slice.
- Provide a CLI command that is useful immediately, likely under a `context` command group.
- Output must be deterministic for fixed inputs.
- Support at least Markdown output for humans and JSON output for future tooling.
- Include guardrails for large repositories, such as ignore rules and a file limit.
- Keep prompt wiring, context-pack budgets, memory, and secret redaction out of this slice.
- Preserve current Go CLI style and tests.

## Non-goals

- No LLM summarization.
- No semantic search or embeddings.
- No direct SDK/provider work.
- No generated `_index.md` changes beyond consuming existing code patterns where useful.
- No Homebrew/release work inside this design idea; release happens after implementation and review are complete.

## Requested output

Each participant should propose:

1. The exact MVP command shape and output formats.
2. The minimal data model for files and Go symbols.
3. Default ignore rules and file-limit behavior.
4. Tests needed for confidence.
5. Any risks or deferred follow-up slices.
