---
agent: codex-1
idea: track-aware-driver
review-round: 2
date: 2026-07-03
reviewed-commit: 70bdc5c
---

## Summary

Fix-up cycle 1 resolves the track-aware-driver issues I reported in round 01: explicit `deliberation` now participates in the non-solo hard reject, the classifier no longer treats unknown or negative size as fast-eligible, idea-level `auto_implement` is used for the fast contradiction check, fast review now hard-gates same-model reviewer sets, and explicit standard now caps cross-review rounds at 2.

I cannot sign off because the required `go test ./...` gate is still red for more than the documented sandbox-only runner failure. `internal/driver TestAcquireLockIsExclusive` fails consistently at commit `70bdc5c`, including when isolated with `-count=20`.

## Verification of round-01 findings

- **RESOLVED — explicit track: deliberation now hits the non-solo hard-reject.** `internal/track.PolicyFor` checks `availableReviewers < 1` before the switch for every explicit track. `TestPolicyForDeliberationNonSolo` and `TestExplicitDeliberationNonSoloEscalates` cover the policy and `driver.New` path. Absent-track legacy remains exempt as documented.
- **RESOLVED — classifier unknown/negative size no longer fast.** `track.Inputs` now has `FilesKnown` and `LOCKnown`, `Classify` requires both known and non-negative for fast, and `runClassify` uses `flag.Visit` plus rejects negative counts with exit 2. Smoke checks: `classify --reversible --mechanically-verifiable` returned `standard`; `classify --files 1 --loc -1 --reversible --mechanically-verifiable` exited `2`; `classify --declared fast --reversible --mechanically-verifiable` exited `4`.
- **RESOLVED — fast + idea-level auto_implement + --no-implement now escalates.** `driver.New` calls `track.PolicyFor(..., ReadAutoImplement(cfg.IdeaDir), ReadStrictGate(cfg.IdeaDir))`, so the contradiction check no longer depends on runtime-masked `cfg.AutoImplement`.
- **RESOLVED — fast now forces a model-diverse reviewer.** `checkModelDiversity` sets `required = true` for explicit `track: fast`, regardless of `require_model_diversity` frontmatter. Existing model-diversity tests pass; I did not find a bypass in `OpenReviewRound`.
- **RESOLVED — standard caps cross-review at 2.** `Policy.CapCrossReviewRounds` is set to `2` for explicit standard and `driver.New` clamps configured `CrossReviewRounds` above that value. `TestExplicitStandardCapsCrossReview` and `TestPolicyForStandardCapsCrossReview` pass.
- **OPEN — full-suite verification has an unexpected failure.** `go vet ./internal/{track,driver,app}/...` passed. Track-focused tests passed. However `go test ./...` failed in both `internal/runner TestDurableKillEndToEndRealProcess` and `internal/driver TestAcquireLockIsExclusive`. The runner failure matches the known codex-sandbox "no recorded boot id" limitation; the driver failure does not.

## New findings (if any)

### [MAJOR] Driver lock can grant multiple concurrent holders in the same process

`go test ./internal/driver -run TestAcquireLockIsExclusive -count=20` fails consistently:

```text
driver_test.go:399: acquireLock granted 2 concurrent holders, want exactly 1
```

The code path is `internal/driver/loop.go` `acquireLock`. It creates the lock file with `O_CREATE|O_EXCL`, then writes the PID token afterward. A racing acquirer can observe the newly-created file before the token is written, fail to parse the empty content as a PID, remove the file as if stale, and acquire its own lock. That violates the single-driver contract and explains the persistent test failure.

Suggested fix: treat empty/unparseable lock contents as held rather than stale, or change lock acquisition to publish a fully-written unique ownership token atomically before other contenders can classify the file as reclaimable. Keep `releaseLock` token-checked.

## Signoff

Status: ❌ BLOCK

Notes: The round-01 track findings are resolved, but the required suite has an unexpected persistent `internal/driver` failure. Only `internal/runner TestDurableKillEndToEndRealProcess` was expected to fail in this sandbox.
