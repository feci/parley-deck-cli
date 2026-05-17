---
agent: codex
idea: roadmap-implementation-plan
round: 1
date: 2026-05-17
---

## Summary

`parley-deck-cli` is past the bootstrap stage. The current code already handles workspace setup, runtime discovery/configuration, headless and interactive invocation, isolated Gemini/Hermes homes, live TUI monitoring, HITL questions, status/resume, consensus signoff requests, and version/status reporting. Round-one headless participants already run in parallel.

The next roadmap should therefore avoid rebuilding orchestration basics. The most useful near-term direction is token/context efficiency that works with today's CLI-backed agents. That means repo maps, deterministic context packs, artifact trimming/scrubbing, and better failure classification before memory layers or provider-specific caching.

## Proposed approach

Status against the idea catalogue:

- Complete or mostly complete: agent runtime config, interactive/headless modes, isolated homes for Gemini/Hermes, live run TUI, HITL Q&A, status/resume, consensus signoff workflow, semantic version/status awareness, and parallel round-one execution.
- Partial: checkpointing/resume exists for run inspection and pending signoffs, but not as a full workflow checkpoint engine. Roles exist as protocol guidance, but not as rich agent personas. TUI agent controls show runtime details, but not full run launch/edit orchestration.
- Missing and feasible now: repo map, deterministic context pack/budgeting, trajectory/artifact compaction, diff-fence conventions for future implementation handoffs, hooks, error classification, redaction helpers, local memory/insights.
- Not feasible as a first slice in the current architecture: Anthropic prompt caching and semantic LLM caching for hosted agents invoked through opaque external CLIs. Those require either direct API/SDK backends or explicit provider CLI support for cache controls and cache keys.

Recommended sequence:

1. `repo-map-mvp`: add a deterministic repository map command/library. Start with file tree plus Go package/type/function extraction via the Go standard library, with text fallback for non-Go files. No tree-sitter dependency yet.
2. `context-pack-mvp`: add a small reusable context packer that can combine the idea prompt, repo map, selected artifacts, and strict byte/token-ish budgets. Include deterministic head/tail trimming and think-block scrubbing.
3. `round-review-context`: wire context packs into later round/review/signoff prompt builders so agents receive compact prior context instead of raw growing transcripts.
4. `agent-error-classifier`: classify agent failures into auth, permission/sandbox, timeout, rate limit, malformed artifact, and unknown. Surface suggested recovery in `status`, TUI, and run results.
5. `hooks-mvp`: add explicit pre/post run hook points for local commands, with conservative defaults and clear failure behavior.
6. `memory-insights-mvp`: add opt-in local SQLite insights after the prompt/context pipeline exists, so memory has a concrete place to feed into.

The first implementation should be `repo-map-mvp`. It is self-contained, testable, low-risk, and becomes useful both as a CLI command and as an input to later context packing. It also avoids the ambiguity around provider prompt caching.

Concrete first-slice shape:

- Add an internal repo-map package that walks a workspace with ignore rules for `.git`, `node_modules`, build outputs, Parley run logs, and common cache directories.
- Extract Go symbols using `go/parser` and `go/ast`; report packages, types, funcs, methods, and imports at a compact level.
- Add `parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]`.
- Keep output deterministic and stable for tests.
- Do not auto-inject the map into agent prompts in this slice; that should be the next slice after the command and data shape are reviewed.

## Concerns / open questions

- The name `context repo-map` should be checked against future `context pack` UX. I prefer the grouped command because it leaves room for `parley context pack` later.
- Ignore rules need to be conservative. Over-aggressive ignores can hide important source files; under-aggressive walks can make the command noisy.
- JSON shape should be explicitly marked as an early developer surface unless we want to commit to long-term compatibility now.

## Risks

- Scope creep into a full tree-sitter/aider clone. The MVP should use Go stdlib parsing and deterministic file summaries only.
- Large repositories can still produce huge maps. The first slice needs file count and byte guardrails.
- If the repo map output is too verbose, it will not help token efficiency. Tests should assert stable compact output for fixtures.
