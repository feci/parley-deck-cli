---
agent: codex
idea: unified-tui-home
review-cycle: 2
date: 2026-06-04
reviewed-commit: 91e0e64
---

## Summary

AF2-AF6 are correctly applied. Build and tests are green:
`GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go build ./... && GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go test ./...`

AF1 is only partially applied. The explicit cancel-tracking defer was removed, but `N`-launched runs still inherit the top-level CLI context, which is canceled when `parley tui` returns after `/quit` or empty-input `esc`.

## Findings

### [MAJOR] AF1 still cancels or abandons `N`-launched runs on TUI detach

`runTUIViewWithDiscovery` wires `Start: newLaunchFunc(ctx, root, results)` and then returns immediately after `tui.RunLive` exits (`internal/app/app.go:1884-1893`). `newLaunchFunc` creates each launched run with `context.WithCancel(ctx)` (`internal/app/app.go:1916`) and starts the runner on that derived context (`internal/app/app.go:1917`). The parent `ctx` is the top-level `signal.NotifyContext` from `Run`, whose cancel is deferred at `internal/app/app.go:44-45`.

So `/quit` or empty-input `esc` exits the Bubble Tea program, `runTUIViewWithDiscovery` returns, `Run` returns, and the deferred parent cancel fires. That cancellation propagates into every `N`-launched run even though the user did not press `ctrl+c`. The background `handle.Wait()` goroutine is also not a durable detach mechanism once the CLI process is returning.

Concrete fix: do not derive detached `parley tui` launches from the top-level command context that is canceled on normal TUI exit. Give TUI-launched runs ownership that survives TUI detach, while keeping the returned `Cancel` as the attached-run `ctrl+c` path; then make the background reap path durable enough to record the session.

## Agreed Fix Verification

- AF1: CHANGES REQUESTED. The explicit cancel-all list is gone and a background reap goroutine exists, but normal `parley tui` detach still cancels or abandons launched runs through the parent CLI context/process lifetime.
- AF2: VERIFIED. `runTask` keeps `Root` and no longer passes `Start` into `RunLive`; `N` reports unavailable there.
- AF3: VERIFIED. `TestTranscriptPopulatesFromOnDiskRun` writes `events.jsonl` plus a stdout log, drives an `eventsMsg`, asserts a non-empty `codex` buffer, and checks rendered output.
- AF4: VERIFIED. Done state adds `[done]` to the status line and `/quit or esc to exit` to the input hint.
- AF5: VERIFIED. The old workspace model/API and app-side action helpers are removed, init quits back to the caller, and no live references remain to the deleted TUI/app symbols.
- AF6: VERIFIED. Help says `ctrl+c cancel the attached run, else quit`.
