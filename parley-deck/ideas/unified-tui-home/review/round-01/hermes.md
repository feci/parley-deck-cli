---
agent: hermes
idea: unified-tui-home
review-round: 1
date: 2026-06-04
---

## Summary
Adversarial UX + correctness review of Phase 6 unified TUI (Home + nullable run) per FINAL.md/IMPLEMENTATION.md. Inspected git diff, live.go/live_test.go, app.go:runTUIViewWithDiscovery+newLaunchFunc. Ran go build/vet/test (green). Focus: owner flows, transcript visibility, no-quit-done, N-launch, regressions.

## Findings
### [MAJOR] N (new idea) discoverability and replace semantics unclear after finished run
N is uppercase-only; Home shows "N=new" hint but no visual affordance or lower-case fallback. After DONE, N replaces the finished run (deferred snapshot footer per spec) with no confirmation or prior-run preservation. User may lose context on re-launch. Fix: make 'n' alias, add replace confirmation or stash prior run summary.

### [MAJOR] Agent transcripts may still be invisible for real multi-agent runs
Per D10, refreshBuffers lazily reads stdout post agent.started. live.go uses per-agent agentBuffer + loadFocusTail but no explicit test with actual runner output paths in the provided live_test.go coverage. Owner complaint #3 not demonstrably closed without a live gemini-excluded multi-agent execution trace. Fix: add integration test that asserts non-empty transcript buffers after RunRoundOneAsync completion.

### [MINOR] Potential TUI "stuck" perception on finished runs despite /quit+esc+ctrl+c
doneMsg sets shortState=FIN + badges; no Quit. /quit and esc-at-empty-Home detach per spec, ctrl+c cancels attached run. However, no explicit empty-Home esc path test in live_test.go and no visual "detached" state. User may perceive stuckness after FIN. Fix: ensure status line always shows exit options on done.

### [MINOR] "All available agents" tabs correctly limited via installedAgentIDs
newLaunchFunc calls installedAgentIDs(discovered) which excludes gemini (unchanged policy). Tabs derived only from run participants (FINAL D3). No regression here; only real agents appear. Acceptable.

### [NIT] Dead workspace model left in place
Slice 5 notes deletion deferred to fix-up. Per CLAUDE.md surgical rule this is fine (no breakage now); no current impact on unified flow. No action required in this round.

## Open questions
- Is the timed previous-run snapshot footer still planned post-deferral, or dropped?
- Confirm real multi-agent stdout paths are populated in practice (beyond unit mocks).