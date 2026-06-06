---
agent: hermes
idea: tui-durable-kill
phase: review
round: 1
date: 2026-06-07
---

## Summary
ACCEPT: attribution gate + process-group spawn + event-log reattach fully solves the 2-day stale running badge; no PID-reuse kill path exists; concurrency and keymap safe; tests prove the contract on darwin.

## Findings
MINOR internal/procctl/procctl.go:118 Attributed: the 5-check gate (boot+alive+exact ProcStart+pgid+command prefix) is sufficient; live in-memory paths (execAgentProcess watcher, Handle.KillAgent active map) call KillGroup only on freshly-captured Spawned from the same Start so never hit Attributed — safe. Darwin lstart capture==verify byte-for-byte via identical ps call. Refuse-across-reboot via bootID. KillGroup -pgid + ESRCH=harmless correct. No other kill paths.

MINOR internal/runner/durablekill.go:32 KillAgentDurable + latestAgentProc: json Number unmarshal to float64 for pid/pgid handled before int cast in Spawned reconstruction; terminal-after-started guard prevents two-parley double-kill.

MINOR internal/runner/runner.go:384 execAgentProcess + onStarted: Start failure never emits agent.started (acceptable, matches spec); single Wait goroutine + ctx-watcher no double-Wait or leak; watcher KillGroup only on the sp it spawned.

NIT internal/tui/live.go: activateRun + openRun copy seams: KillAgentFunc(string,error) + Liveness correctly wired for resume/open; no runner/app import.

NIT internal/app/app.go:1970 reattachSeams + liveSteerKillSeams: all RunLive sites now receive durable kill/liveness; two-parley terminal guard + clearStale on stale both exercised.

NIT internal/procctl/procctl_*_test.go + durablekill_test.go: real macOS processes, grandchild reap, attribution refusal, end-to-end durable kill via event log alone; linux /proc parse and watcher-timeout paths are the only untested but low-risk.

No CRITICAL/MAJOR; concurrency (watcher vs durable vs clearStale) and ctrl+k keymap reuse (kill vs clear) have no collisions.

## Verdict
ACCEPT
No blocking items.