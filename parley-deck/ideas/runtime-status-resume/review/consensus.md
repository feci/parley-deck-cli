---
idea: runtime-status-resume
review-cycle: 1
drafted-by: codex
date: 2026-05-12
reviewed-commit: 4e367ae6354fce5e55e50184f95fe4b0ff7328c1
---

## Agreed fixes

- From claude/review/round-01 [MAJOR] Resume TUI header renders `status=running` for non-terminal runs: change resume-mode header wording so it never prints unqualified `running` without proof of a live process, and add a regression test.
- From claude/review/round-01 and gemini/review/round-01 [MINOR] Outcome and liveness coverage gaps: add deterministic tests for incomplete, failed, and idle run projections.
- From claude/review/round-01 [MINOR] Missing CLI coverage for status/resume edge paths: add focused tests for `status --idea`, workspace `status --json`, and nonexistent resume target behavior.
- From claude/review/round-01 [MINOR] `parley status --run` detail shows `duration=-` for a started-only agent: show elapsed time for currently running agent snapshots.
- From claude/review/round-01 [MINOR] `ResolveRun` cannot distinguish a known idea with no runs from an unknown target: return a clearer error when the idea exists but has no run history.
- From claude/review/round-01 [NIT] Resume footer omits `ctrl+c`: advertise the existing `ctrl+c` exit path in resume mode.
- From claude/review/round-01 [NIT] `ideaForRun` synthesizes a path under `parley-deck/ideas/unknown`: leave the fallback path empty when the run has no known idea slug.

## Deferred follow-ups

- From gemini/review/round-01 [NIT] Redundant workspace reads in `inferParticipants`: defer until the run list grows enough to justify passing preloaded workspace status through the projection layer.
- From gemini/review/round-01 [NIT] `ResolveRun` exact-run precedence over idea slug: keep the current precedence for timestamp-like run IDs; document or revisit only if user-defined run IDs become possible.
- From hermes/review/round-01 [NIT] Minor duplication in state constants: defer until broader state taxonomy work exists.
- From hermes/review/round-01 [MINOR] Additional graceful-degradation test coverage for missing idea slug inference: defer beyond this fix-up unless touched by the agreed tests above.
- From hermes/review/round-01 [NIT] JSON `RunDir` omission is undocumented: defer until the unstable JSON surface is documented for consumers.

## Dismissed findings

- None.

## Signoffs

### Signoff: codex — 2026-05-12
Status: ✅ ACCEPT
Notes: Implementer accepts the review consensus and will apply the agreed fixes in fix-up cycle 1.

### Signoff: claude — 2026-05-12
Status: ✅ ACCEPT
Notes: The agreed fixes cover my MAJOR resume-header finding and related minor/nit items, so codex may proceed to fix-up cycle 1.

### Signoff: gemini — 2026-05-12
Status: ✅ ACCEPT
Notes: The agreed fixes address my coverage gap finding for outcomes and liveness, and I agree with the deferment of my performance and shadowing nits; codex may proceed to fix-up cycle 1.

### Signoff: hermes — 2026-05-12
Status: ✅ ACCEPT
Notes: The consensus handles my non-blocking findings as deferred follow-ups and includes the blocking fix from review round 1, so codex may proceed to fix-up cycle 1.
