---
agent: hermes
idea: repo-map-mvp
round: 1
date: 2026-05-17
---

## Summary
Operations and resilience lens: prioritize deterministic, bounded, fault-tolerant repo scanning with clear guardrails for large repos and graceful degradation.

## Proposed approach
Command: `parley context map --format markdown|json --max-files N --ignore .gitignore`

- Output formats: Markdown (human-readable tree + symbol list) and JSON (machine: {files: [{path, symbols: [{name, kind, line}]}]})
- Data model: minimal struct {Path string; Symbols []Symbol}; Symbol {Name, Kind (func/type/var), Line int}
- Defaults: respect .gitignore + built-in ignores (node_modules, .git, vendor, dist); hard cap 5000 files or 100MB total; fail fast on permission errors with partial results + warning
- Resilience: stream output, no full in-memory load; timeout per-file parse (Go stdlib only); always produce valid output even on partial scan

Trade-off: smaller MVP scope means no incremental indexing yet; revisit in follow-up slice.

## Concerns / open questions
- How to surface partial results vs hard failure when hitting limits?
- Should JSON include scan metadata (ignored count, elapsed time) for ops observability?

## Risks
- Stdlib parser misses edge cases in complex Go (trade-off accepted for MVP size).
- Large monorepos may still hit memory if symbols explode; deferred to later resilience hardening.