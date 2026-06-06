---
agent: hermes
idea: deliberation-driver
round: 1
date: 2026-06-05
---

## Summary
The proposed ~300 LOC port correctly harvests the essential pattern (dispatch_once idempotency via artifact existence, recompute_ready as pure function of durable state, crash-safe re-entry via disk-derived cursor) while safely dropping SQLite-specific machinery. The linear graph simplification works for the degenerate case but the consensus BLOCK→reopen back-edge requires explicit handling in the ready-phase switch to avoid hidden non-linearity.

## Proposed approach
Keep exactly these kanban invariants as load-bearing:
- recompute_ready is a pure function of on-disk artifacts only (no in-memory task graph)
- dispatch_once is enforced by Overwrite=false + os.Stat guard (runner.go:292)
- crash-safe re-entry: Rebuild derives everything from disk files; cursor is cache only
- ready-set computation reduces to "all participant round-N/*.md exist"

Safe to drop: claim_lock/heartbeat/zombie-reaping/blackboard/topological-sort (linear graph has no contention or cycles beyond the explicit reopen edge).

The advisory PID lock + disk cursor is sufficient; the SQLite atomicity gap for "mark in-flight" is closed by the existing atomic Save (pipeline/run.go:82) plus RunRound's idempotent skip.

## Concerns / open questions
The ordered switch for readyPhase is faithful only if the reopen path (consensus.go triage Blocked → Reopen + CurrentRound++) is treated as an explicit transition that forces Rebuild to re-evaluate the highest round-N dir. Without that, a partial round after reopen could be misclassified as "complete".

Crash re-entry must explicitly handle: (1) round-N dir exists but not all agents wrote (stuck gate), (2) cursor deleted mid-RunRound (Rebuild must still see partial artifacts and skip already-written agents), (3) duplicate dispatch prevented solely by the artifact check.

## Risks
Dropping the kernel creates a correctness gap only if "mark in-flight" is needed to prevent concurrent RunRound invocations on the same round; the advisory lock mitigates but a crash between Rebuild and first write could theoretically allow a second process to start the same round. The linear simplification hides the back-edge; the port must surface the BLOCK reopen as a distinct case in the switch or risk duplicate round-N seeding.