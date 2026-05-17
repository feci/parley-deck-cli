---
agent: claude
idea: roadmap-implementation-plan
round: 1
date: 2026-05-17
---

## Summary

Through the architecture/UX/risk lens: the catalogue's highest-leverage themes
that are also safely achievable with the current opaque-CLI invocation model
are (a) deterministic context/artifact compaction, (b) safety hooks at the
protocol boundary (redaction, error classification), and (c) lightweight local
memory built from canonical artifacts. Provider prompt caching, semantic
cache, sub-agent delegation inside opaque CLIs, MCP hosting, and advanced
cognition should stay deferred until a direct SDK backend exists — the
constraints document is right to fence them off. The next implementable slice
should be a single deterministic round-summary index, because every later
context/memory/insights idea wants this artifact as input.

## Proposed approach

### Triage of `parley-deck-ideas.md` against current state

- **Already implemented (don't redo):** workspace init, roles, headless and
  interactive invocation, isolated agent homes, live TUI, HITL questions,
  `status` / `resume`, consensus signoff, version reporting, round-one
  headless parallelism.
- **Partially implemented:**
  - Memory / local memory provider — the canonical `ideas/<idea>/round-NN/`
    layout is already a filesystem memory, but it is unindexed and uncompacted.
  - Error classification — `internal/runner/validation.go` exists; not yet a
    typed taxonomy with exit codes consumable by hooks.
  - Batch execution — round-one parallelism exists; cross-round and
    cross-idea batching do not.
- **Missing but feasible now (deterministic, local):** round-summary index,
  artifact redaction, diff fences in artifact display, repo-map generator,
  hook points (`pre-write`, `post-round`), structured error taxonomy,
  credential/rate guard wrappers around `exec.Cmd`.
- **Missing and currently infeasible / out of scope:** Anthropic prompt
  caching (no control over cache metadata through opaque CLIs), semantic
  cache (needs embedding store + provider), in-process sub-agent delegation
  (CLI subprocess boundary), MCP/server hosting, microagents, advanced
  cognition. Defer until/unless a direct SDK backend lands.

### Ordered slice roadmap (each = one focused PR)

1. **Round-summary index (`round-NN/_index.md`)** — deterministic, generated
   after `runner` finishes a round; per-agent abstract built from heading +
   first paragraph scan; canonical, lossy by design.
2. **Artifact redaction at the protocol write boundary** — single chokepoint
   in `internal/protocol` that strips a configured set of secret patterns
   (env-style tokens, `Authorization:` headers, JWT-shape strings) before any
   agent artifact or HITL question file is written.
3. **Structured error taxonomy + exit codes** — promote ad-hoc errors in
   `runner/validation.go` to a typed set; surface stable exit codes; enables
   later hooks and reliable `resume`.
4. **Hook points (`pre-round`, `post-round`, `pre-write`)** — minimal
   user-script invocation, documented contract, deterministic env. Consumes
   (3).
5. **Repo map command (`parley repo-map`)** — deterministic file/symbol
   outline writable into an idea's round directory so agents can reference
   it via filesystem rather than re-listing.
6. **Local insights memory** — derive per-idea `insights.md` from prior
   round indices (1) using deterministic rules first; LLM-assisted variant
   only after a direct SDK backend exists.
7. **Credential / rate guard wrappers** — uniform `exec.Cmd` wrapper that
   enforces per-provider concurrency and surfaces rate-limit errors through
   (3).
8. **Cross-round / cross-idea batch runner** — extends current round-one
   parallelism; depends on (3) and (7) to be safe.
9. **(Conditional) Direct SDK backend pilot** — only if and when the team
   accepts the migration cost; this unlocks prompt caching, semantic cache,
   in-process sub-agent delegation, and think-block control. Until then,
   keep the opaque-CLI path as the supported invocation surface.

### First slice to implement next

**Slice 1: `round-NN/_index.md` deterministic round-summary writer.**

- Trigger: `internal/runner` after every participant artifact in a round has
  been written and validated.
- Input: the round directory's `<agent>.md` files (existing canonical
  artifacts).
- Output: `round-NN/_index.md` with YAML frontmatter
  (`idea`, `round`, `generated_at`, `participants`) and one section per
  agent containing: the agent's own frontmatter values, the section headings
  it produced, and the first paragraph under each heading (truncated to a
  fixed line cap).
- No model call. No network. Idempotent. Rewrites in place when re-run.
- New CLI surface: `parley round summarize <idea> <round>` for manual /
  catch-up use; runner calls the same library function.
- UX: index file is plain Markdown and human-readable; future agents and
  consensus tooling can reference it instead of re-reading every artifact.

**Why this slice first (trade-off):** It is the smallest deterministic step
that produces a reusable artifact every later slice (insights memory,
context compaction, batch runner summaries) wants as input. It does not
require any provider behavior, does not change existing canonical artifact
semantics, and is reviewable in one focused PR. The trade-off is that it
adds a generated artifact to the canonical tree, which slightly enlarges
review surface for each round — accepted because the per-round savings
downstream will outweigh it.

### Tests / checks for slice 1

- Unit: given a fixture `round-01/` with three agent artifacts, generated
  `_index.md` matches a golden file (exact bytes).
- Unit: rerunning the writer over an existing `_index.md` produces identical
  bytes (idempotency).
- Unit: an agent artifact with malformed YAML frontmatter is skipped with a
  recorded warning rather than aborting the whole index.
- Unit: heading scan respects ATX (`##`) levels and truncates paragraph
  capture at the configured line cap.
- Integration: `parley round summarize <idea> 1` on a real workspace writes
  `_index.md` and exits 0; second invocation is a no-op diff.
- Manual UX: `_index.md` renders cleanly in the existing TUI artifact view
  and in plain `cat`.

## Concerns / open questions

- **Index placement and naming.** `_index.md` keeps it lexicographically
  first in directory listings and visually distinct from agent files; happy
  to defer to a different name if another participant has a strong reason.
- **Whether the index is canonical or derived.** I propose canonical
  (committed alongside agent artifacts) so consensus and later rounds can
  rely on it without regenerating. The alternative (derived, gitignored) is
  cheaper but breaks deterministic replay.
- **Heading-scan heuristic vs. requiring a fixed schema in agent files.**
  Current schema (`## Summary` / `## Proposed approach` / …) is already
  effectively required by the round prompt; I propose the writer assume it
  and degrade gracefully when absent, rather than enforcing a hard schema
  this slice.
- **Confirming the opaque-CLI fence.** I am assuming the project agrees
  prompt caching, semantic cache, and in-process sub-agent delegation stay
  out of the first sequence (constraints already say so). If any participant
  argues for a direct SDK backend as slice 1, that reorders this entire
  roadmap and should be a consensus decision, not a unilateral one.

## Risks

- **Lossy summarization risk:** the index is intentionally lossy. Mitigation:
  canonical agent artifacts remain the source of truth; the index is
  explicitly labeled as a generated abstract; tooling that needs full
  fidelity must read the agent files.
- **Schema drift risk:** if agent prompt templates change their heading
  structure, the index degrades. Mitigation: graceful skip + warning, plus
  fixture-based golden tests that fail loudly when the writer's behavior
  changes.
- **Scope-creep risk on slice 1:** tempting to bundle redaction, hooks, or
  an LLM-summarized variant. Resist — those are slices 2, 4, and 6
  respectively. Each must be reviewable in isolation.
- **False sense of progress on context efficiency:** a deterministic index
  is a real but bounded win. It does not replace provider-level prompt
  caching, and the roadmap should not claim it does. The unlock for the
  bigger context-efficiency wins is the conditional SDK backend (slice 9),
  not this one.
- **Replay / determinism risk:** if `generated_at` is included in the index,
  byte-identical re-runs fail. Mitigation: either omit `generated_at` from
  the body, or pin it to the latest agent-artifact mtime so deterministic
  inputs yield deterministic output.
