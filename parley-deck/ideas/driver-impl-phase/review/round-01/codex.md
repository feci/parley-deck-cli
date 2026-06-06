---
agent: codex
idea: driver-impl-phase
review-round: 1
date: 2026-06-06
reviewed-commit: f624c05
---

## Summary

The implementation is mostly aligned with FINAL.md D1-D10: the ImplOps seam is clean, Rebuild is most-terminal-first, Implement and Fixup both require `auto_implement` plus a clean git tree, the no-land boundary is preserved, blocked review consensus escalates, and re-entry is disk-derived. Verification passed with explicit caches:

`GOCACHE=/Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-cli/.gocache`
`GOMODCACHE=/Users/tomasfecko/go/pkg/mod`
`go build ./... && go vet ./... && go test ./... && GOOS=windows GOARCH=amd64 go build ./...`

I found one protocol/safety issue that should be fixed before accepting the phase.

## Findings

### MAJOR: Fix-up opens the next review round without the required RunChecks gate

In `internal/driver/impl.go`, `advanceReview` runs `Fixup(ctx, cycle)` and then immediately archives `review/consensus.md` and calls `OpenReviewRound(ctx, round+1)` (lines 181-190). There is no `RunChecks` call after the fix-up and before the next review round.

This violates FINAL.md D3/D4 and consensus D4, which require `RunChecks` after implementation and after each fix-up, before `OpenReviewRound`; failing checks must escalate immediately with result text. As written, a fix-up can introduce compile/test failures and still spend reviewer agents on the next round. Because Rebuild sees `review/round-NN` first, a crash after opening that round also strands the system in PhaseReview without ever applying the missed checks gate.

Concrete fix: after `Fixup` succeeds, call `d.cfg.Impl.RunChecks(ctx)` before archiving/opening the next review round. If it fails, return `ActionEscalated` with the check output and do not archive the old consensus or open `round+1`. Add a regression test that scripts outstanding fixes, a successful `Fixup`, failing `RunChecks`, and asserts no `OpenReviewRound(round+1)` call.

### MINOR: Implementer selection ignores the FINAL drafter / IMPLEMENTATION frontmatter role contract

`internal/app/driver_impl.go` always selects `participants[0]` as the implementer and review-consensus drafter (lines 35-49). Consensus D10 says the implementer is the FINAL drafter, or IMPLEMENTATION.md frontmatter on re-entry, restricted to participants; reviewers are the non-implementers.

This matters when the FINAL drafter is not the first participant, or after re-entry from an existing IMPLEMENTATION.md authored by a different participant. The current adapter can run code-writing under the wrong agent and include the actual implementer in the reviewer set, weakening implementer/reviewer separation.

Concrete fix: resolve the implementer from `IMPLEMENTATION.md` frontmatter when present, otherwise from `FINAL.md` frontmatter (`implementer` or `drafted-by`), validate it is in `participants`, and only fall back to `participants[0]` when no durable role metadata exists. Build reviewers after that resolved role.

### MINOR: Non-ready IMPLEMENTATION.md statuses escalate instead of awaiting

`advanceImpl` treats any status outside `implemented`, `ready-for-review`, or `fix-up-cycle*` as `ActionEscalated` (internal/driver/impl.go lines 78-84). FINAL.md D6 says malformed/empty should escalate, but not-ready should await.

This makes a normal in-progress implementation artifact behave like a hard failure. A crash/re-entry while the implementer has created `IMPLEMENTATION.md` but not yet marked it review-ready could halt the whole auto-run instead of waiting for the implementation attempt deadline.

Concrete fix: distinguish malformed/empty/unknown statuses from known not-ready statuses such as `in-progress` or `fixing`; return `ActionAwait` for known not-ready states and keep escalation for malformed or unrecognized values. Add a driver test for an in-progress status returning `ActionAwait`.

## Open questions

None.
