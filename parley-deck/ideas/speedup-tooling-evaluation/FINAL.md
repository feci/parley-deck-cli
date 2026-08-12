---
idea: speedup-tooling-evaluation
status: final
drafted-by: claude-1
date: 2026-08-11
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
---

# FINAL — none of the three tools; do the two protocol changes

## Decision

**Adopt neither cognee, graphify-as-context-selector, nor omniroute.** Implement the two already
ranked protocol changes, with **no new dependency**:

- **rank 1 `protocol-phase-scoped-packet`** — send only the sections a phase needs. Cuts the cost of
  a round (3.3x median wall clock per call), and is paid every round, so it compounds.
- **rank 3 `protocol-fixup-budget`** — a finite fix-up budget that **escalates rather than closes**.
  `deliberation` is currently unbounded, which is why 24-round reviews exist.

Both need §7.

## Why the tools do not apply

Reached independently by @codex-1, @hermes-1 and @kimi-1 in round 1; all three arrived at the null
answer by different routes.

| Tool | What it is | Measured term it touches | Verdict |
| --- | --- | --- | --- |
| graphify | deterministic local code/doc graph, installed, indexes this repo | none directly — bodies still must be read, and round count is protocol-driven | different bottleneck: navigation, not volume or rounds |
| cognee | LLM-mediated memory: ingest, recall, injected context | none | *"solves a problem we do not have"* — this deck **over**-remembers (quadratic history re-send); it does not under-remember |
| omniroute | multi-provider gateway with optional lossy prompt compression | per-call token volume | *"right term, inadmissible mechanism"* — only via lossy compression or a silent model swap |

**The standard applied**, and the reason it is not negotiable: any tool that decides what an agent
sees occupies the exact position of the frontier machinery deleted in 1.43.1, which was removed
because it could not prove it never drops a participant objection. Under Phase 2 rule 1 — *"Silence
= implicit agreement"* — a dropped objection is not a lost datum; it is recorded consent that was
never given.

**graphify is kept, outside the trusted path** (@codex-1): useful for exploration and for *auditing*
a phase packet, never for *selecting* one. Its output stays non-authoritative.

## The structural finding that constrains rank 1

@hermes-1, PRIMARY, verified twice in two separate ideas: the Go runner never reads
`COOPERATION.md` at all — zero references in `internal/runner/runner.go`,
`internal/runner/phase58.go`, `internal/app/driver_consensus.go`.

The 3.3x does not arise in the runtime. It arises in the **instructions**: the skill's "Always read
`COOPERATION.md` first" and the prompts a facilitator writes by hand. **Rank 1 must therefore be
built in the instruction layer.** Built into the Go prompt builder it would touch nothing.

## What must never be delegated to a tool

Selection of normative context. Any mechanism that chooses which rules an agent sees must be able to
prove completeness, and none of the three can. A tool that moves neither measured term is a cost,
not a win.
