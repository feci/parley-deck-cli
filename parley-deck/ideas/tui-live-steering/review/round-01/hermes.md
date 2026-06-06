---
agent: hermes
idea: tui-live-steering
phase: review
round: 1
date: 2026-06-06
---

## Summary
Keymap/concurrency review against FINAL §A/B/C/D/E: implementation matches ratified spec with no CRITICAL/MAJOR defects; minor test gaps and one edge-case guard noted.

## Findings
MINOR internal/runner/steer.go:197 (runSteerAgent after waiter): steerBusy cleared via defer in happy path and mkdir-err early return, but the depth-1 waiter goroutine path (lines 170-190) lacks explicit clear on ctx cancellation before exec; add defer clearSteerBusy immediately after setting busy to guarantee even on queued-wait abort.
MINOR internal/tui/live.go:620 (updateMain confirm-kill block): modal correctly intercepts all keys first (after ctrl+c) and highest priority; however suggest-mode Esc path (line 580) precedes it when both active — impossible in practice but add explicit "if confirmKillAgentID != \"\" { ... } else if suggest { ... }" guard for future-proofing.
NIT internal/tui/live_test.go:312 (TestLiveSteerSubmit): asserts submitSteer + reply buffer attach and conditional-Tab, but omits explicit test that second steer while busy produces "already replying" status in the input row; add one case exercising the reject path from §E.
NIT internal/runner/steer_test.go:110 (TestRunSteerAttemptRejectsSecond): covers busy-reject and kill-vs-finish races under -race, but the queued-wait poll loop (agentRunning polling) is exercised only indirectly; add a subtest that launches a round attempt then immediately steers to assert depth-1 serialization and no goroutine leak on run end.
NIT internal/app/app.go:1751 and 2039 (liveSteerKillSeams wiring): both TUI launch sites correctly inject SubmitSteer+KillAgent and activateRun copies both seams; no import of runner from tui — seams are plain funcs as required by §D.

## Verdict
ACCEPT (no CRITICAL or MAJOR items that would block merge; all concurrency rules, key-collision sweep across normal/picker/suggest/confirm-kill/answer modes, and §E contract hold).