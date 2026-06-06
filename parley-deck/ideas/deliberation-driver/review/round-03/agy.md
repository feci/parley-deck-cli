---
agent: agy
idea: deliberation-driver
review-round: 3
date: 2026-06-05
reviewed-commit: a83efa8
---

## Summary

Slice 2 successfully implements the consensus gate logic outlined in D6/D7 and the live adapter wiring. The import-direction guarantee (the driver core never importing the app layer) is preserved using dependency injection (`ConsensusOps` interface). However, several issues concerning platform-specific file renaming, locking behavior, and protocol-correct agent selection need to be addressed before slice 2 can be finalized.

## Findings

### MAJOR: Failed Final Draft Leaves Workspace Stranded in Pseudo-Finalized State
* **What's wrong:** In `internal/app/driver_consensus.go`, `DraftFinal` calls `consensus.Finalize`, which creates the scaffold `FINAL.md` and updates `00-prompt.md` to `status: final` before running the drafter agent. If the drafter agent fails (due to timeout, API error, or failing the D7 non-scaffold checks), the driver escalates and exits. However, on the next tick, `Rebuild` sees that `FINAL.md` exists on disk and sets the phase to `PhaseFinal`. In this phase, `Advance` immediately returns `ActionSurfaceOnly`, bypassing drafting and validation completely.
* **Why it matters:** This violates the crash-recovery and idempotency contracts. A transient agent failure permanently blocks the driver from retrying the draft, leaving the workspace stranded with `status: final` but with an empty/scaffold `FINAL.md`.
* **Concrete Fix:** 
  1. In `advanceConsensus`, check the non-scaffold status even if the file exists:
     ```go
     if !fileExists(finalPath) || finalScaffoldReason(finalPath) != "" { ... }
     ```
  2. In `driverConsensusOps.DraftFinal`, if `runDrafter` fails, delete `FINAL.md` and revert `00-prompt.md` to `"consensus"`.

### MAJOR: Facilitator-Drafter Selection Ignores Idea Participants
* **What's wrong:** In `internal/app/driver_consensus.go`, `firstHeadlessAgent` selects the first discovered headless agent on the system that is not `"gemini"` to draft the consensus and the final plan. It does not verify whether the selected agent is actually in the idea's `participants` list.
* **Why it matters:** Under the Parley Deck protocol (§4, §6), the facilitator/drafter must be a participant of the deliberation. If an external agent is selected, it violates protocol boundaries (an outsider synthesizes the consensus) and will fail at runtime if that agent lacks the appropriate API credentials or configuration for the workspace, whereas active participants are guaranteed to be functional.
* **Concrete Fix:** Pass the active `Participants` roster from `app.go` to the `newDriverConsensusOps` constructor and filter `discovered` agents to prioritize or restrict selection to the idea's actual participants.

### MAJOR: Windows process liveness check always fails, breaking driver locking
* **What's wrong:** In `internal/driver/loop.go`, `processAlive` checks PID liveness by sending signal `0` via `proc.Signal(syscall.Signal(0))`. On Windows, `os.FindProcess` always returns a `Process` struct without checking if the process actually exists, and calling `Signal` with any value always returns an error because Unix signals are not supported on Windows. This makes `processAlive` always return `false` on Windows.
* **Why it matters:** Since `processAlive` always returns `false`, `acquireLock` will always treat any existing lock file as stale and reclaim it. This defeats the concurrency lock completely on Windows, allowing multiple racing instances of `parley run --auto` to run concurrently and corrupt the workspace files.
* **Concrete Fix:** Implement a cross-platform liveness check helper (e.g. using conditional compilation for Unix vs Windows, or calling platform-specific system APIs to check if the process is alive).

### MAJOR: Stale consensus file invalidation (`invalidateStale`) fails on Windows when backup exists
* **What's wrong:** In `internal/driver/consensus.go`, `invalidateStale` renames `consensus.md` and `FINAL.md` to `.bak` using `os.Rename`. It silently discards the error. On Windows, `os.Rename` fails with an error if the destination file already exists.
* **Why it matters:** If an idea goes through multiple BLOCK/reopen cycles, the `.bak` files will already exist. On Windows, the rename will fail silently, leaving the stale `consensus.md` and `FINAL.md` files intact. The next `Rebuild` will see these files and incorrectly classify the phase as `PhaseConsensus` or `PhaseFinal`, bypassing the reopened rounds.
* **Concrete Fix:** Explicitly delete the `.bak` file if it exists before renaming:
  ```go
  bakPath := path + ".bak"
  if fileExists(bakPath) {
      _ = os.Remove(bakPath)
  }
  _ = os.Rename(path, bakPath)
  ```

### MINOR: Stale Invalidation before Round Execution is Not Crash-Safe
* **What's wrong:** Under the `TriageBlocked` case in `advanceConsensus`, `invalidateStale` is called to rename the files *before* `d.runner.RunRound` is executed. If `RunRound` fails, the driver escalates and halts, but the consensus file has already been renamed and no new round has been successfully created. On the next tick, `Rebuild` sees no `consensus.md` and only the old round directories, so it classifies the phase as `PhaseRound` for the *previous* round, sees that it is complete, and re-drafts the same blocked consensus.
* **Why it matters:** A temporary failure in initializing the new round causes the driver to forget the BLOCK decision, re-draft the old consensus, and enter an infinite loop.
* **Concrete Fix:** Reorder the operations under `TriageBlocked`: only run `invalidateStale()`, update `00-prompt.md` status, and update the cursor *after* the new round has been successfully initialized by `d.runner.RunRound(ctx, next)`.

### MINOR: Neglecting `syscall.EPERM` in Unix Process Liveness Check
* **What's wrong:** `processAlive` assumes any error from `proc.Signal(syscall.Signal(0))` means the process is dead. On Unix, if the process exists but is owned by a different user, the call returns `syscall.EPERM` (Permission Denied).
* **Why it matters:** If another user owns a process with that PID, `processAlive` will return `false`, causing the driver to reclaim the lock of a live process.
* **Concrete Fix:** Treat `syscall.EPERM` as a sign that the process is alive:
  ```go
  err := proc.Signal(syscall.Signal(0))
  return err == nil || err == syscall.EPERM
  ```

## Open questions

- **Should the facilitator be configurable in 00-prompt.md?** Currently, there is no way to explicitly designate a facilitator agent (it is hardcoded to pick the first headless agent). If an idea specifies a facilitator role, should the driver respect it?
- **How should we handle human-in-the-loop signoffs on local-dir?** If a human participant is in the list, the driver currently escalates under `TriagePartial` since humans cannot be invoked as headless agents. Is this the intended fallback?
