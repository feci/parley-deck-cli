---
idea: close-integrity
author: claude-1
created: 2026-06-24
participants: [claude-1, codex-1, hermes-1, antigravity-1]
implementer: claude-1
status: final
spawned_from: loop-engineering-research
checks: go build ./... && go test ./internal/driver/... ./internal/runner/... ./internal/app/... ./internal/protocol/...
---

## Problem / idea

Tier 3 of the `loop-engineering-research` backlog: **close-decision integrity** — make the
moment the auto-driver marks an idea complete trustworthy. Two items:

- **LE-7 goal-done gate.** Before `Complete()` under `auto_implement`/`strict_gate`, a
  fresh agent checks the FINAL.md observable acceptance criteria (loop engineering's
  `/goal`: "a separate model checks whether you're done"). Run once, before close, reusing
  the `consult.go` execution machinery as the separate checker.
- **LE-11 HITL-fatigue guardrails.** Under `auto_implement`, do **not** silently
  auto-complete on a soft `ACCEPT-WITH-RESERVATIONS` triage, and refuse to auto-complete
  with fewer than 2 independent reviewers — both re-create a quiet false-green that a
  fatigued human waves through.

Implemented by claude-1 (Phase 5); reviewed by codex-1, hermes-1, antigravity-1
(Phase 6, refutation). Design = `FINAL.md`.
