---
agent: hermes
idea: runner-hardening-kindly
round: 1
phase: review
date: 2026-06-12
---

## Summary

Reviewed runner-hardening-kindly Phase 6 implementation (D1-D12) per consensus scope. Focus areas: D1 kill ordering and durable-kill races, D3 retry semantics, D9 snapshot filesystem correctness, D10 AppendLine atomicity. Diff inspected: supervision.go, runner.go, reviewsnapshot.go, fsutil.go, acp.go, COOPERATION.md copies, and supporting changes. Tests not executed (environment constraint); no runtime verification performed.

## Findings

### MINOR waitSupervised tick granularity
**File:** internal/runner/supervision.go (waitSupervised loop)
**Why:** 1-second tick means sub-second watchdog windows effectively round up to ~1s. The waitErr vs. watchdog race is mitigated by the "append watchdog event FIRST, then KillGroup, then drain" ordering, but the granularity remains coarse.
**Suggested fix:** Document explicitly that the 1s tick is the accepted granularity; production windows (120s/30m) make this acceptable.

### NIT AppendLine stuck-claim degradation
**File:** internal/fsutil/fsutil.go (AppendLine)
**Why:** After ~5s stuck claim, degrades to unlocked append rather than failing. Matches D10 trade-off.
**Suggested fix:** None required; the rationale is already recorded in consensus.

### NIT Marker vs. durable-kill attribution clarity
**File:** internal/runner/runner.go (attempt loop, move-aside, marker)
**Why:** attempt_id threads correctly through events and procctl marker; D3 retry rule ("retry ONCE only for no_first_output") is implemented. No collision detected.
**Suggested fix:** None.

### NIT Snapshot commit and move-back races
**File:** internal/runner/reviewsnapshot.go
**Why:** Temp-index snapshot commit (clone + worktree on virtio-fs), MoveArtifactBack (copy+fsync+rename within target dir), marker/sweep/step-aside all follow D9 mechanics. No cross-device rename attempted.
**Suggested fix:** None.

## Dispositions

- **Finding/disposition: waitSupervised uses a 1-second tick, so sub-second watchdog windows effectively round up to ~1s. Prior disposition: accepted trade-off (kindly uses the same granularity; production windows are 120s/30m).**  
  **Concurrence:** I concur. The production timeouts render the 1s granularity a non-issue; the ordering guarantee (watchdog event before Kill) is the critical correctness property.

- **Finding/disposition: AppendLine degrades to an UNLOCKED append after a ~5s stuck-claim wait instead of failing. Prior disposition: accepted trade-off (consensus D10: losing the ledger record is worse than a rare interleaved line).**  
  **Concurrence:** I concur. The design correctly prefers durability over perfect atomicity on virtio-fs.

## Verdict

ACCEPT

No blocking issues found. All reviewed areas align with consensus D1-D12. The two trade-offs are acceptable under the stated conditions.