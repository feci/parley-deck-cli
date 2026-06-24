---
idea: verification-honesty
author: claude-1
created: 2026-06-24
participants: [claude-1, codex-1, hermes-1, antigravity-1]
implementer: claude-1
status: final
spawned_from: loop-engineering-research
checks: go build ./... && go test ./internal/driver/... ./internal/runner/... ./internal/app/... ./internal/protocol/...
---

## Problem / idea

Tier 1 of the `loop-engineering-research` backlog: **verification honesty** — close
the three "false-green" seams the loop-engineering audit exposed, where gates we
believe we have are inert. Bundled because all four items touch the same files
(`internal/runner/phase58.go`, `internal/app/driver_impl.go`, `internal/driver/*`,
`COOPERATION.md` §Phase 6/8).

This idea is **implemented by claude-1** (Phase 5); reviewed by codex-1, hermes-1,
antigravity-1 (Phase 6, refutation-mode). The design is `FINAL.md` (below), drafted
from `loop-engineering-research/FINAL.md` + the four round-01 lenses.

## Scope (LE-1, LE-3, LE-2, LE-4)
- **LE-1** Refutation-default review posture + a machine-validated `## Refutation
  attempts` section so an empty-findings review must show its work.
- **LE-3** Model-diversity guard: warn (and, via opt-in `require_model_diversity`,
  escalate) when every reviewer shares the implementer's model.
- **LE-2** Implement `strict_gate` enforcement (`ReadStrictGate` + `strict_gate_clean`
  / `closing_review_round` close fields + a deterministic finding-scan veto).
- **LE-4** Generalize `RunChecks` beyond Go via a `checks:` command, fail-closed for
  code-writing auto ideas with no checks — which transitively ties the "artifact-wins"
  fix-up override to a real check (the post-fix-up RunChecks gate, hermes #8).

The full per-file spec is in `FINAL.md`.
