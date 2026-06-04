---
agent: hermes
idea: unified-tui-home
round: 2
date: 2026-06-04
responding-to: [claude/round-01, codex/round-01]
---

## Position changes since prior round
Adopt convergence on live.go + Home tab + StartRunFunc/activateRun + doneMsg-no-quit + DONE badges + retire workspace + runToken. From adversarial UX lens: the /open path + N-replaces-run is the minimum that avoids tab-cursor footguns, but DONE badges alone may not make "stale finished run" obvious enough when N overwrites. Transcript verification gate is accepted; gemini exclusion stays as legacy.

## Responses to others

### @claude
Agree on /open <slug|run> as the v1 escape hatch for existing runs (no in-Home cursor because ↑/↓ = tabs). But this is brittle UX: users on Home with no run will type N expecting new, or hunt for /open without discovery. Minimum that feels right is a one-line "recent runs" list on Home that shows slug + last status (DONE or active), with hint text "N=new  /open <slug>". No selection cursor yet. On done persistence: DONE badge + [done] status is clear for re-read/steer, but N-replace risks silent loss of the finished transcript context. Counter-proposal: on N during done state, keep the finished run's last tab snapshot in a small "previous" footer line for 10s so user sees what was replaced. This costs almost nothing and prevents "where did my finished run go?" confusion.

### @codex
Agree on StartRunFunc + activateRun + runToken + shared startRoundOne. On gemini: agree keep excluded as roster-legacy (owner "all available" means discoverable Found agents; gemini is explicitly retired and never Found=true in practice). On #3 transcript: agree with claude — refreshBuffers already lazy-resolves via agentByID each tick, so the bug is unverified rather than a known missing line. Phase-5 real multi-agent run gate is the right call; no preemptive code change. On MVP scope: cut exactly as proposed (ActionRunner, session-browser, in-Home-cursor); nothing else missing for v1.

## New concerns / questions
- Home → N flow feels right only if the input prefix ("new idea ›") is shown immediately on uppercase N and never collides with lowercase steer input on agent tabs.
- When N replaces a finished run, does the model fully reset buffers/offset or could a stale transcript flash for one tick?

## Current proposal
Unify on liveModel with Home + optional active run. MVP locked: Home ideas/runs list, N new-idea launch, /open existing, agent tabs only for run participants, done-no-quit + per-tab DONE + [done] status, runToken, retire RunWorkspace/ActionRunner. Verify transcripts in Phase 5; keep gemini excluded. Slices unchanged from round-1 convergence.