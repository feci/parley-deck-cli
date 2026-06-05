---
agent: codex
idea: deliberation-driver
review-round: 1
date: 2026-06-05
reviewed-commit: f8c880d
---

## Summary

The slice is well scoped and the round-completion gate does validate artifacts before reconciliation; it does not promote on mere file presence. `GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go build ./... && go vet ./... && go test ./...` passes.

I found correctness issues in the concurrency contract and escalation path, plus one protocol-validity regression from the D4 deviation.

## Findings

MAJOR: `driver.lock` acquisition is not atomic, so two `parley run --auto` processes can both become drivers. `internal/driver/loop.go:103` reads the PID file, then `internal/driver/loop.go:112` overwrites it with `os.WriteFile`; two starters can both observe no live lock and both enter `Advance`. In that race, `internal/driver/driver.go:102` can see `round-02` incomplete in both processes and both can call `RunRound(2)`. `Overwrite=false` does not close this hole because `runner.runAgent` has its own pre-launch `os.Stat` TOCTOU, which D10 explicitly leaves to the driver lock. This violates the documented single-driver contract and can double-launch agents, interleave writes to the same artifacts/logs, and emit duplicate round events. Fix by acquiring the lock with an atomic primitive (`os.OpenFile(path, O_CREATE|O_EXCL|O_WRONLY, 0644)` after stale-lock removal, or `flock`/`F_SETLK`) and make release remove the file only if it still contains this process' ownership token. Add a concurrent acquisition test that starts two goroutines at a barrier and proves only one `RunRound` can occur.

MAJOR: The round>=2 gate is weaker than FINAL D4 and accepts artifacts that do not actually respond to every participant. `internal/driver/driver.go:141` only checks that `responding-to` frontmatter exists, while `runner.ValidateRoundArtifact` for cross-review rounds only requires matching frontmatter and any `## ` heading. That means an artifact with `responding-to: [prior round artifacts]` and a generic `## Summary` can satisfy the gate and allow consensus readiness without per-agent engagement. The implementation rationale says this is because `BuildRoundPrompt` emits `## Responses to other participants`, but `internal/runner/runner.go:624` already tells agents to respond by name and the accepted design made the stricter `### @<other>` headings the machine-checkable evidence. Fix by changing the runner prompt's required file shape to include one `### @<other participant>` under `## Responses to other participants`, then restore the D4 validation for every other active participant. Keep `responding-to` as an additional metadata check, not the only cross-review evidence.

MAJOR: Malformed event logs halt but do not create the durable escalation promised by the design. `internal/driver/driver.go:147` returns an error when `terminalRoundEvent` cannot load `events.jsonl`, and `internal/driver/loop.go:41` only prints `driver: halting` and returns that error. Unlike deadline handling in `internal/driver/loop.go:64`, this path writes no blocking inbox note, so an unattended auto run can stop with no durable artifact explaining the corrupted event log or how to recover. This is especially relevant because D4 says malformed event log goes to escalate, not just stderr. Fix by routing `roundComplete` errors through an escalation writer that records the idea, round, event-log path, parse error, and recovery instructions before returning.

## Open questions

Should reconstructed `round.completed` events also include and validate the `idea` field when scanning terminal events? D4 says the terminal event is for idea+round; `terminalRoundEvent` currently matches only `round`, which is probably harmless for one-idea runs but is looser than the spec.

Should the cross-review acceptance test use the real prompt contract and assert per-participant `### @...` evidence? The current fake artifacts mirror the weakened gate, so they would not catch the D4 regression.
