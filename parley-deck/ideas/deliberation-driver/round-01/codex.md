---
agent: codex
idea: deliberation-driver
round: 1
date: 2026-06-05
---

## Summary

The proposed `internal/driver` package is the right Go shape if it stays small, boring, and disk-derived. The current failure mode is not missing agent execution machinery; it is missing an owner for the next tick after `round.completed`. I confirmed the core root cause: `runaction` has no open-next-round action, `runplan` moves a complete round directly to consensus and treats BLOCK as inspect-only, `runner.RunRound` already supports `N >= 2`, and consensus already owns Draft, Status, AppendSignoff, Finalize, and Reopen.

I support a linear ordered switch rather than a DAG or scheduler. This driver is not allocating work across workers; it is advancing one idea through one protocol. The idiomatic Go version is a small package with a `Cursor` cache, pure disk inspection, and a `Driver.Advance(ctx, cursor)` method whose dependencies are injected as function fields or a small interface. The cursor must never be authoritative. If it is missing, corrupt, or stale, rebuild from disk and continue.

## Proposed approach

I would structure the package as:

- `cursor.go`: `Cursor`, atomic `Save`, `Load`, and `Rebuild`. Keep schema fields minimal: phase, current round, idea status, rounds run, max rounds, updated time. Do not encode claims or long-running ownership in the cursor.
- `driver.go`: pure readiness functions and one-step `Advance`. This should be a readable ordered switch: final exists, consensus exists, consensus triage, latest round complete, latest round incomplete. No DAG abstraction is needed for the deliberation path.
- `loop.go`: optional polling loop plus `driver.lock` advisory PID file. The lock protects normal CLI double starts; correctness still comes from idempotent gates.

The important test seam is not a fake filesystem; it is a fake executor boundary. Define a narrow dependency surface, for example:

```go
type RoundRunner interface {
	RunRound(ctx context.Context, round int) error
}

type ConsensusOps interface {
	Status() (consensus.Summary, error)
	Draft(round int) error
	RequestSignoffs(ctx context.Context) error
	Finalize() error
	Reopen(reason string) error
}
```

The production adapter should call the existing launch paths. For cross-review rounds, call `runner.RunRound(ctx, opts)` with `Round: next`, `RoundLabel: round-%02d`, and `Overwrite: false`. Do not create a second agent launcher in `internal/driver`. `RunRound` already delegates to `RunRoundOne`, which launches configured participants concurrently, writes `round-NN/<agent>.md`, appends segment events, writes the round index, and emits `round.completed` or `round.incomplete`.

For signoffs, the driver should reuse the existing request-signoffs path rather than building its own process loop. `internal/app/consensus_request_signoffs.go` already resolves missing participants, discovers configured agents, validates launch modes, invokes headless agents through `runner.CommandFor`, writes manual or interactive handoff packets, and validates that each agent appended exactly its own signoff. The driver-facing seam should be extracted so app and driver share it, likely by moving the reusable non-CLI logic into a small internal package or exported app helper. The driver must never synthesize ACCEPT signoffs.

The first slice should prove only round promotion:

- Given a completed `round-01` on disk, `Advance` computes next action `open round-02`.
- `Advance` calls the injected round runner with `2`.
- On success, disk state becomes `round-02`, cursor saves atomically, and subsequent `Advance` does not call `RunRound(2)` again if `round-02` is already complete or in progress.
- The real wiring is gated behind `--auto` and `local-dir` only.

Table tests should cover:

- `round-01 complete -> RunRound(2)` using a fake runner that records calls.
- existing `round-02` artifacts plus `round.completed -> no duplicate RunRound(2)`.
- missing participant artifact or `round.incomplete` event -> no promotion, return await/incomplete.
- missing or corrupt cursor -> `Rebuild` derives the same phase from round directories, prompt status, consensus status, and `FINAL.md`.
- `TriageBlocked -> Reopen -> RunRound(current+1)` with fakes verifying order and MaxRounds enforcement.

## Concerns / open questions

The largest idempotency hole is not `Overwrite=false`; it is deciding when a round is complete. Artifact presence alone is too weak because a concurrent or crashed agent may leave a partial-looking file. The gate should require the expected participant artifacts and a matching `round.completed` event for that round, preferably also accepting the round index if that is the local convention. If the event says `round.incomplete`, do not promote.

`Overwrite=false` makes re-entry mostly safe for already-written artifacts, but it is not a complete concurrency strategy. Two driver processes can both observe round-01 complete and both call `RunRound(2)`. The advisory PID lock is sufficient for the intended CLI model only if the gate remains idempotent and `RunRound(2)` can tolerate this race. Today each agent skips when the output file already exists, but two processes can pass `os.Stat` before either writes. That can lead to duplicate launches for the same agent. I would not add claim locks in the first slice, but I would make the driver lock mandatory for the loop and treat failure to acquire it as a clean stop. Tests should document that correctness is single-driver plus re-entry, not multi-writer execution.

The cursor save should copy the pipeline atomic temp-and-rename pattern, ideally with temp files in the same directory. A corrupt cursor should be ignored only after recording a warning or returning a typed rebuild result; silent rebuilds make debugging hard.

`consensus.Reopen` sets idea status back to the latest existing round. The driver must then run the next round after that, not confuse "reopened to round-01" with "round-01 is the target again." The pure readiness function should compute `nextRound = latestRound + 1` after a BLOCK reopen, bounded by `MaxRounds`.

`runplan` and `runaction` should get `KindOpenNextRound` for visibility, but I would keep execution in `internal/driver`. The plan package should describe the next action; it should not become the driver.

## Risks

The first implementation can accidentally broaden scope into a second orchestrator. Keep phase one to cursor, ready phase, and `round-01 -> round-02` promotion.

The signoff path currently lives under `internal/app`, so reusing it from `internal/driver` may create an import direction problem. Extract the reusable request-signoffs service before driver integration rather than importing app from driver.

The advisory PID lock can leave stale files. That is acceptable if stale detection is conservative: read PID, check whether the process is alive when supported, and otherwise allow a clear manual recovery path. Do not build heartbeat or zombie reaping.

The tests need to avoid live agents. The minimal useful seam is an injected `RunRound(ctx, opts)` function or interface plus injected consensus operations. Full integration tests can use temporary workspaces and fake agent commands later, but unit tests for `Advance` should be deterministic and fast.
