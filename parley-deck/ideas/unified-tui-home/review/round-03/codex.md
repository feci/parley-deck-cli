---
agent: codex
idea: unified-tui-home
review-cycle: 3
date: 2026-06-04
reviewed-commit: 245b9eb
---

## Summary

AF1 is VERIFIED. `runTUIViewWithDiscovery` now owns a `launchReaper`, passes it
into `newLaunchFunc`, and waits with `reaper.waitForActive(stdout)` after
`tui.RunLive` returns and before the command returns. That keeps the parent
signal context live during detach waiting, so normal `/quit` no longer cancels or
abandons in-flight `N`-launched runs, while real SIGINT can still abort and the
attached run's in-TUI `Cancel` path remains unchanged.

`TestLaunchReaperWaitsForInFlightRuns` covers the key detach behavior by proving
`waitForActive` blocks until the tracked in-flight run completes. The requested
build/test command is green:
`GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go build ./... && GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go test ./...`

AF2-AF6 remain intact from round-02. No CRITICAL, MAJOR, MINOR, or NIT findings
remain.
