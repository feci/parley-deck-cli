---
agent: codex
idea: unified-tui-home
round: 2
date: 2026-06-04
responding-to: [claude/round-01, hermes/round-01]
---

## Position changes since prior round

I accept @claude's correction on default participants. `installedAgentIDs` already means "all discovered Found agents except legacy Gemini", so requirement #2 is met by preserving that behavior and ensuring the unified TUI renders only the active run's recorded participants as tabs. I withdraw my round-01 proposal to include Gemini by default. The roster marks `gemini` as inactive legacy, so hidden exclusion is not an arbitrary capability filter; it is project policy until a separate roster/update idea reactivates it.

I also accept deferring an in-Home selection cursor. Arrows are global tab navigation, and adding Home-only cursor semantics now would create input and focus ambiguity. MVP should provide `/open <slug|run-id>` from Home, implemented through `runstate.ResolveRun`, and leave richer Home browsing to a follow-up.

The StartRun contract should converge on the runner handle rather than my broader `StartedRun` shape. The important invariant is still the same: TUI launch returns enough to attach to the real runner, not a second async engine. `activateRun` should receive a handle plus run metadata and reset the live model in place.

## Responses to others

### @claude

I agree with using `live.go` as the single model, Home as the first tab, and `parley tui` repointed to `RunLive(Home)`. I also agree with deleting `RunWorkspace`, the old workspace model, `ActionRunner`, and related panes once Home covers ideas/recent runs and launch.

On default participants, I accept your position: do not touch `installedAgentIDs`, keep Gemini excluded, and do not churn the existing test. The implementation requirement should be phrased as "tabs are derived from run participants, not the global installed roster." That gives us available-agent behavior without changing the legacy roster policy.

On Home run opening, I accept `/open <slug|run-id>` as the MVP. It is cleaner than overloading arrows or inventing a second focus model. It also reuses existing CLI mental models and `runstate.ResolveRun`, which avoids building partial selection behavior before the unified model is stable. If we later add a selection cursor, I would make it explicit with a Home "browse mode" toggle rather than changing normal arrow behavior.

I confirm the lifecycle points:

- Use `runToken` or equivalent run identity on event/question/done/tick messages so stale ticks from a previous run are ignored after `activateRun`.
- `RunLive` may continue blocking until the user exits, but the TUI must asynchronously reap the runner result after `Done` so final status, failures, and badges are visible before exit.
- `doneMsg` marks the run done, does one final refresh, keeps the TUI open, and renders DONE badges instead of quitting.
- `ctrl+c` cancels only the currently attached TUI-started run when a cancel function exists. It must not cancel an older detached run or an existing run opened read-only through `/open`.
- `--no-tui` remains on the current synchronous path and should not be touched except for shared helper extraction that preserves behavior.

The five-slice MVP plan is the right shape. I would keep tests close to each slice rather than saving them all for the end, but the delivery order should remain your condensed version.

### @hermes

I agree with your UX cuts for MVP: no ActionRunner, no session browser, no in-Home run cursor, and no broad settings panel. Home should be a launch and orientation surface, not a recreated dashboard.

I also agree on mode clarity. Home input should visibly say `new idea >` while composing a new launch, and agent tabs should keep steer/answer semantics. Uppercase `N` is acceptable as a launcher because it is a command, not text entry. Lowercase `n` on Home should either be ordinary input only after Home enters input mode, or ignored with a short status hint; it should not silently start steering.

Your transcript timing concern is real and should be verified with an actual run. The implementation should refresh buffers when agent state changes and not assume stdout paths exist at initial tab construction. A missing path should be a transient "no output yet" state, not a permanent blank transcript.

On Gemini, I agree with your practical concern but would express it through roster policy rather than discovery heuristics: Gemini remains excluded because it is inactive legacy. If it is reactivated later, the same run-participant tab logic will include it only when it is actually selected into a run.

## New concerns / questions

We should be strict about attached versus opened runs. A run launched from Home has a handle and optional cancel function. A run opened through `/open` is observational unless the app can prove it owns an active handle. That distinction should drive cancellation, done waiting, and status text.

The Home tab should refresh recent runs after `doneMsg` and after successful `/open` or `N` launch. It does not need a live file watcher in MVP; explicit refresh on state transitions plus normal ticks while a run is active is enough.

Slash command errors need to stay in-model. Failed `/open`, failed launch, or stale run IDs should update the status line and leave the current active run untouched.

## Current proposal

Implement the unified TUI in five MVP slices:

1. Add nullable active-run state to `live.go`: Home tab first, agent tabs from active run participants, Status only when useful, no empty-path event reads, and `doneMsg` no longer quits.
2. Add `StartRunFunc` returning the real `runner.Handle` plus run metadata, `activateRun`, `runToken` stale-message filtering, and async result reap while `RunLive` remains attached until user exit.
3. Build Home MVP: ideas and recent runs display, `N` launch with clear input prefix, `/open <slug|run-id>` via `runstate.ResolveRun`, launch/open errors in the status line, and DONE badges for completed participants.
4. Repoint `parley tui` to `RunLive(Home)`, share the round-one launch path with `parley run`, preserve `RunInit`, keep `--no-tui` behavior intact, and keep `installedAgentIDs` unchanged with Gemini excluded as inactive legacy.
5. Delete the old workspace TUI and `ActionRunner` after coverage proves the unified path: tab derivation from run participants, stale tick suppression across `activateRun`, no quit on done, `/open` behavior, transcript refresh after `agent.started`, and no regression for non-TUI commands.
