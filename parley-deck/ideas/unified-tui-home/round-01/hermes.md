---
agent: hermes
idea: unified-tui-home
round: 1
date: 2026-06-04
---

## Summary
Adversarial UX review of the unified `parley tui` (Home + agent tabs). The Home → N-launch → agent-tabs flow has footgun risk on input meaning; "don't quit on done" leaves stale finished tabs that confuse re-use; per-agent transcripts likely broken by timing (tabs built before agent.started events); 4-5 agent tabs are usable but legacy gemini will pollute; retiring workspace model loses ActionRunner/next-actions which mattered for quick follow-ups. Minimum viable: cut ActionRunner, session browser, run selection on Home; keep only ideas list + N launcher + live tabs + done state.

## Proposed approach
- Home tab (00-prompt:27) renders Status.Ideas + recent runs (from runstate); `N` switches input to `new idea › ` then calls StartRunFunc (app.go:82) which does agents.Discover().Found (discover.go:120) to get all available, launches via runner path, returns RunHandle → liveModel swaps to active run and sets activeTab to first agent (live.go:94).
- Input row meaning: on Home = launcher only (N mode or direct task); on agent tab = steer/answer. No overlap because N is uppercase-only and never conflicts with lowercase steer. Footgun avoided by showing mode prefix ("new idea ›" vs "steer ›").
- Don't quit on done (00-prompt:36): on doneMsg (live.go:91) set m.done=true, render "[done]" status, keep ticking and allow N (replaces run) or /quit. Finished tabs stay visible with last transcript; add per-tab "DONE" badge.
- Per-agent transcript (#3): hypothesize buffer path empty because LiveOptions passed before any agent.started event (app.go runTask vs RunLive), tabs created in init (live.go:97) with empty path, ensureActiveBuffer never re-polled on event. Verify: check agent.started writes stdout path (runmanifest), then refreshBuffers on every event for active run only.
- All agents tabs: usable for 4-5 (claude-cli-layout windowing); drop legacy gemini (discover.go agents list) if Found=false or mark inactive.
- Retire old workspace (app.go:81 runTUI, tui/app.go entire model): lose ActionRunner (tui/app.go:50) and next-actions pane; keep only for v1 if cheap else cut (owner said discard rest).
- MV vs gold: cut run selection on Home, ActionRunner, session browsing, Status tab when no-run; ship only Home + N + agent tabs + done persistence.

## Concerns / open questions
- Input mode collision: if user types lowercase 'n' on Home does it become steer or error? (live.go input routing)
- Done run persistence: after N new run, old finished tabs are replaced cleanly or leak state in m.run?
- Transcript fix: does loadFocusTail get called after agent.started arrives for the exact tab? (live.go:64 buffer logic)
- Legacy agents: gemini in discover.go list but never Found=true in real envs — still appears in "all available"?

## Risks
- Tab overflow with 5 agents breaks lipgloss strip (live.go tab render).
- StartRun inside TUI bypasses runTask path → inconsistent Done channel wiring (app.go:71).
- Users lose quick "next action" execution from old dashboard; must re-learn N flow.