---
agent: hermes
idea: tui-live-steering
round: 1
date: 2026-06-06
---

## Summary
Keymap analysis of live.go:updateMain shows Tab/shift+tab, arrows, and N/enter/esc already committed to tab nav, picker, compose, and submit. Adding slash autocomplete on Tab requires conditional dispatch. Kill needs a new dedicated key (ctrl+k) plus confirm sub-mode that does not overlap picker, N, scroll, or steer-enter. Concurrency rules must treat steers as independent one-shot segments with per-agent locks to prevent races on kill vs finish, or dual steers.

## Proposed approach
Resolve Tab collision by making tab navigation conditional: Tab and shift+tab only switch tabs when inputText does NOT start with "/"; when it does, Tab triggers autocomplete complete (longest common prefix or single match) and arrows/enter operate inside the suggest menu (reuse pickerState with new kind pickerSuggest). Move all tab nav fully onto left/right arrows as primary; keep shift+tab/tab as legacy only for non-/ input.

New kill key: ctrl+k when an agent tab is active and agent state==running. Opens a confirm sub-mode (simple yes/no overlay, not full picker) using y/n or enter/esc. No collision with existing keys.

Full key table (normal mode / picker active / autocomplete visible / confirm-kill / answer-compose):

normal (no submodes, input empty or non-/):
- ctrl+c: quit/cancel run
- esc: clear input or quit
- N: enter compose (new idea) if input==""
- tab / shift+tab / left / right / up / down: switch tabs
- shift+up/down, pgup/pgdown, ctrl+u/d, home/end: scroll transcript
- enter: submitSteer or answer
- backspace/ctrl+h, printable: edit inputText
- / : start slash input (shows command hint)

picker active (open/answer):
- arrows/enter/esc/filter runes: picker navigation (unchanged)
- all other keys: no-op

autocomplete visible (input starts with /, suggest menu shown):
- tab: complete prefix / accept
- arrows/enter: select from suggestions
- esc / backspace past / : close suggest, return to normal
- printable (after /): filter suggestions live
- other nav/scroll: blocked while suggest open

confirm-kill (ctrl+k on running agent tab):
- y/enter: confirm kill (emit agent.killed, cancel only that agent ctx)
- n/esc: cancel confirm, return to normal
- all other: no-op

answer-compose (after /answer pick or open question):
- enter: submit answer (color-flipped row)
- esc: cancel answer
- printable/backspace: edit
- tab/nav: blocked

Concurrency rules:
- A steer always spawns a fresh single-agent segment (new segment ID) via CommandFor + runAgent; never reuses an in-flight process.
- If target agent is currently running its round: queue at depth-1 (one pending steer per agent). The running attempt finishes first; queued steer executes immediately after in its own segment.
- Per-agent mutex or state flag (in runner registry) prevents two steers from starting concurrently on same agent.
- Kill (ctrl+k) cancels only the current runAgent ctx for that agent; if a steer is queued it remains queued and can still execute later. A finished agent (normal exit) simply clears the registry entry; a late steer after end-of-run is allowed as a post-run side invocation (recorded with steer.replied).
- State never ambiguous: every attempt writes its own <agentDir>/steer-N/stdout.log and a steer.replied event; transcript tail appends from the attempt log.

Answers to 6 focus questions (keymap/concurrency lens):
1. Spawn immediate fresh single-agent invocation (parallel ok because one-shot); prompt = steer text + last artifact path + transcript tail (last 50 lines). Reply surfaced via new per-attempt log tailed into the same agent transcript buffer.
2. Registry of per-agent context.CancelFunc (created inside runAgent before cmd.Start). Killed agent lands in stateFailed with agent.killed event; round continues. Race-safe: kill checks registry presence before cancel; normal completion removes entry under same lock.
3. Reuse pickerState (new kind pickerSuggest) for autocomplete menu. Tab semantics: complete to longest common prefix; if exactly one match accept full command. No collision because / prefix gates the behavior; N remains special-cased before rune append.
4. Agent's steer reply renders appended into existing transcript buffer from the attempt-specific log file. While running: status line shows "agy replying..." + spinner.
5. Steering running agent: queue depth-1. Two quick steers: second waits. Steer after run ended: allowed as independent invocation. Per-agent lock + segment ID ensures no clobber of artifact dir.
6. New seam on LiveOptions: SteerFunc and KillAgentFunc (injected like LaunchFunc/Cancel). Headless tests: model-driven TUI tests for key dispatch + fake CommandFor that records invocations without exec.

## Concerns / open questions
- Confirm sub-mode needs its own minimal state flag (confirmKillAgentID) to avoid polluting picker.
- Should /autocomplete also offer args for /open and /answer (reusing existing picker) or keep simple command list only?
- Legacy shift+tab users: do we document the arrow-only recommendation?

## Risks
- Conditional Tab risks muscle-memory breakage for users who type / then expect tab-nav; mitigate with clear hint line change.
- Depth-1 queue + per-agent lock prevents corruption but adds visible latency if a long-running agent is steered; owner may perceive "nothing happened".
- Kill confirm must be strictly modal (block all other keys) to avoid accidental steer-Enter during confirm.