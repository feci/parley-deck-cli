---
agent: codex-1
idea: parley-learn-playbooks
review-round: 2
date: 2026-07-04
reviewed-commit: 8825dd3
responding-to: [codex-1/review/round-01, hermes-1/review/round-01]
---

## Position changes since prior review round

I withdraw my round-01 MAJOR finding. `internal/app/learn.go` now checks the `parley-deck/playbooks` parent with `os.Lstat`, refuses a symlinked parent before constructing the target write, creates the parent only when absent, and creates the final playbook with `os.OpenFile(..., os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)`. That closes the concrete symlinked-parent write escape I reproduced in round 01 and removes the final-target check/write window from the previous `Lstat` + `WriteFile` path.

The new `TestLearnRejectsSymlinkedPlaybooksDir` covers the regression shape I asked for: a symlinked `playbooks/` parent fails closed and does not create the playbook in the symlink target. I also concur with the fix-up response to Hermes' wording concern: §13.5 now says `parley learn` scaffolds a deterministic skeleton that the author refines, which matches the v1 implementation rather than overselling automatic extraction.

Validation run: `go build ./...` passed; `go test ./internal/app ./internal/protocol` passed; the targeted symlink regression test passed. Literal root `go vet` exits with `no Go files in .../wt-learn-playbooks` because this repo has no package at the module root; package-scoped `go vet ./...` passed.

## Updated findings

zero remaining

## Open questions

None.
