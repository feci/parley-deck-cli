---
idea: rho-retro-tooling
author: claude
created: 2026-06-16
participants: [claude, codex, agy, hermes]
status: implementation
design-source: meta-protocol-change-rho-retrospective-optimization/FINAL.md
---

## Problem / idea

Implement the **`parley retro` tooling** half of the approved RHO proposal. The
design is already ratified (4/4) in
`ideas/meta-protocol-change-rho-retrospective-optimization/FINAL.md` (decision D5
+ the layered-harness/guardrail decisions); this idea is its implementation and
review (Phases 5–8). No fresh design rounds — review against the parent FINAL and
the new §13 protocol section.

## Scope (v1, from FINAL D5)

Read-only-by-default `parley retro` CLI with staged subcommands:
- `scan` — read-only inventory of `parley-deck/ideas/*` with deterministic
  failure-density signals.
- `select [--k N]` — type-diverse coreset of the hardest cases.
- `diagnose [--k N]` — grouped, deterministic diagnosis report.
- `propose --slug SLUG` — scaffold **only** a single new `ideas/<slug>/00-prompt.md`
  (fail-if-exists); writes nothing else.

## Constraints

- Read-only default; `propose` is the only writer and writes exactly one file.
- Structured artifacts only (no raw session JSONL in v1). Deterministic ranking
  (no DPP/embeddings/re-rollout in v1).
- Governed by COOPERATION.md §13; never edits protocol/harness or any participant
  artifact.

## Non-goals

- DPP/embeddings, live re-rollout, best-of-N, raw-JSONL ingestion, auto-apply,
  persistent quarantine registry (all deferred per the parent FINAL).
