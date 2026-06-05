---
agent: codex
idea: deliberation-driver
review-round: 3
date: 2026-06-05
reviewed-commit: a83efa8
---

## Summary

Slice 2 mostly matches the D6/D7 consensus gate: `Advance` is split into round and consensus paths, `ConsensusOps` keeps the driver core out of `internal/app`, partial signoffs are requested synchronously and rechecked, BLOCK reopens are bounded by `MaxRounds`, and the live adapter is wired from `runTask`. The D9 deviation is acceptable for this slice because the import-direction guarantee is preserved by dependency injection: `internal/app` imports `internal/driver`, while `internal/driver` does not import `internal/app`.

Verification passed with workspace caches set explicitly: `GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go build ./... && go vet ./... && go test ./...`.

## Findings

MAJOR: `DraftFinal` can leave the idea durably marked `status: final` before a usable `FINAL.md` exists. In `internal/app/driver_consensus.go`, `DraftFinal` calls `consensus.Finalize`, which writes the scaffold `FINAL.md` and updates `00-prompt.md` to `status: final`, before invoking the drafter agent. If the drafter exits non-zero, times out, is cancelled, or writes content that fails `finalScaffoldReason`, `advanceConsensus` escalates, but the disk state is already `FINAL.md` present plus idea `status: final`. On the next driver tick, `Rebuild` sees `FINAL.md` and returns `PhaseFinal`, so the consensus gate no longer retries or revalidates the final artifact. This violates the D7 "verify non-scaffold before advance" contract and the crash-recovery/idempotent re-entry contract because a failed final drafter can strand the idea in an apparently final phase with only scaffold content. Concrete fix: split finalization into two phases so the drafter writes a candidate first and only then commits final status. For example, draft into `FINAL.md.tmp` or `FINAL.draft.md`, validate with `finalScaffoldReason`, then atomically rename to `FINAL.md` and set `00-prompt.md` to `status: final`; alternatively add a consensus package API that writes the final scaffold without changing idea status and a separate commit step after validation. Also make `advanceConsensus` validate an existing `FINAL.md` before treating Ready/Reserved as final, so crash/retry with a stale scaffold cannot bypass D7.

MINOR: `invalidateStale` ignores rename failures and can silently leave stale `consensus.md` or `FINAL.md` in place after a BLOCK reopen. `consensus.Reopen` already renames `consensus.md` to an aborted round artifact, but `invalidateStale` still needs to protect the driver from any remaining stale files, especially `FINAL.md`. Today `os.Rename(path, path+".bak")` errors are discarded, so a permissions problem, existing `.bak`, or filesystem error can leave `FINAL.md` present; the next `Rebuild` will classify the idea as `PhaseFinal` even though the driver just opened a new round. Concrete fix: make `invalidateStale` return an error, use a collision-safe backup path, and have `advanceConsensus` escalate before running the next round if stale invalidation fails.

## Open questions

- Is "first available headless agent" the intended facilitator-drafter selection rule long term? It is protocol-compatible enough for this slice's D6 facilitator action, but if the drafter is also a participant, the system may want an explicit configured facilitator or deterministic preference to avoid surprising authorship.
- Should `ActionSignoffsRequested` immediately continue into finalization when the post-request status is Ready/Reserved? The current one-action-per-tick behavior is consistent with the driver design and terminates on the next loop iteration, so I do not consider this a defect.
