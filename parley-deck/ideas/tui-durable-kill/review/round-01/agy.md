---
agent: agy
idea: tui-durable-kill
phase: review
round: 1
date: 2026-06-07
---

## Summary
REQUEST-CHANGES: The process-group spawn and durable TUI key/badge seams are mostly correct, but the attribution safety gate is bypassed if critical recorded fields are missing, and ACP process control integration was entirely omitted from the implementation.

## Findings

1. **CRITICAL**, [procctl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl.go#L98-L127) `Attributed`: The safety gate treats missing persisted fields (such as `BootID` or `ProcStart`) as optional and skips validation if they are empty in the `Spawned` record. If an older run log lacks this metadata or if capture failed to query the OS, the gate only verifies that the PID is alive. Under PID wrap-around/reuse, this will erroneously attribute and SIGKILL an unrelated process.
   - *Fix*: Fail closed in `Attributed` if `supportsDurableKill` is true and critical fields (`BootID`, `ProcStart`, `PGID`, `Command`) are empty or missing in the recorded process state.

2. **CRITICAL**, [acp.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/acp.go#L27-L96) `runACPAgent` & [spawn.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/acp/spawn.go#L106-L124) `Stop`: ACP process group and durable-kill integration was completely omitted. The `agent.started` event written on the ACP path does not contain process identity fields (`pid`, `pgid`, `boot_id`, etc.), meaning these processes are completely invisible to `KillAgentDurable`. Furthermore, the termination path in `Stop` hard-kills the direct child process via `Process.Kill()`, leaving any grandchildren orphaned on cancel or timeout.
   - *Fix*: Capture the process identity metadata using `procctl.Capture` right after spawn, persist it in the ACP `agent.started` event, and route process teardown through `procctl.KillGroup` to ensure the entire tree is reaped.

3. **MAJOR**, [procctl_darwin.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl_darwin.go#L48) `procStart`: macOS start-time check uses `ps -o lstart`, which has only a 1-second wall-clock resolution. If a PID wraps around and is reused by a process with the same command within the same second, it can satisfy the attribution checks and be targeted for termination.
   - *Fix*: Retrieve a high-resolution start timestamp using sysctl (e.g. `kern.proc.pid.<pid>` to extract boot time + process start ticks) or libproc APIs, and store that raw numeric value instead of the formatted `lstart` string.

4. **MAJOR**, [live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L870-L898) `renderHome`: The home banner warning `⚠️ N stale agent process(es) — ctrl+k on their tabs to clear` is not implemented in the TUI overview. Users have no visibility into stale, un-killed agent processes left behind across restarts unless they open a specific run and switch to the stale agent's tab.
   - *Fix*: Scan active runs for stale processes in the Home view, and render the warning banner when any are detected.

5. **MAJOR**, [procctl_unix.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl_unix.go#L25-L46) `KillGroup`: If the process group target `PGID` resolves to parley's own process group (which can happen if `Setsid` fails or is skipped), calling `syscall.Kill(-target, syscall.SIGTERM)` will terminate parley itself and all parent/sibling processes in the current session.
   - *Fix*: Add a self-protection check using `syscall.Getpgid(0)` and return an error if `target == selfPGID`.

6. **MINOR**, [procctl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/procctl/procctl.go#L129-L137) `commandMatches`: The command check matches symmetrically using `strings.HasPrefix(recorded, live)`. If a recorded command was long and the live command is truncated, this is acceptable; however, a shorter recorded command matching a longer, unrelated live command is dangerous.
   - *Fix*: Restrict the matching logic to only check if `live == recorded` or if `strings.HasPrefix(live, recorded)`.

7. **MINOR**, [app.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/app.go#L2087-L2097) `reattachSeams`: Errors from `DurableKillAt` are returned as success outcome strings rather than Go errors. Consequently, the TUI prints messages like `"kill failed: ..."` in green success style (`okStyle`) instead of red warning style (`warnStyle`).
   - *Fix*: Change `KillAgentFunc` or `reattachSeams` to return a proper error when the kill fails, so the TUI displays the failure in the red warning row.

8. **NIT**, [runner.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L309-L383) `runAgent`: The timestamp for `agent.started` is captured at the very beginning of the run step. Slow prompt generation or directories setup might create a discrepancy where the event time noticeably precedes the actual process startup.
   - *Fix*: Capture `time.Now().UTC()` in the `onStarted` callback and use it for the event timestamp.

9. **NIT**, [steer.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/steer.go#L216) `runSteerAgent`: Steer processes do not record process metadata or write `agent.started` events. If parley is restarted during a steer, the process will become an orphan and cannot be durably terminated.
   - *Fix*: This is an acceptable limitation given that steers are transient follow-up replies, but it should be documented.

### Correct Points
- Spawning headless processes with `Setsid` and executing group-kills on context cancellation effectively reaps grandchildren on timeouts.
- The `Wait` goroutine and watcher logic in `execAgentProcess` are clean, with no double-Wait or resource leak regression.
- Liveness checks and durable kill seams are correctly wired to all TUI entry points (resume, open, run, Home), keeping the TUI decoupled from runner.
- The TUI `STALE` badge rendering, key bindings, and confirmation dialogue logic successfully solve the 2-day stale badge issue.
- Windows cross-restart fallback behaves correctly by refusing attribution and clearing the stale status.

## Verdict
REQUEST-CHANGES

### Blocking Items
- **Attributed Safety Gate**: `Attributed` must fail closed if critical verification fields are empty.
- **ACP Integration**: Implement metadata capture and process-group kill for ACP agents to match the spec.
