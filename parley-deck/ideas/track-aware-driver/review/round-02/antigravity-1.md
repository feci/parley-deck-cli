---
agent: antigravity-1
idea: track-aware-driver
review-round: 2
date: 2026-07-03
reviewed-commit: 70bdc5c
---

## Summary

The fix-up cycle 1 successfully addresses all findings from round-01. The §4.0 rigor tracks are now robustly enforced, with the non-solo invariant applied to explicit deliberation runs and the classifier fail-safes correctly handling unknown or negative sizes. Contradiction gates and model-diversity requirements are properly wired at the idea-declaration level, ensuring that runtime flags like `--no-implement` cannot be used to bypass safety invariants.

The full suite is green, and `go vet` is clean.

## Verification of round-01 findings

- **[RESOLVED] explicit track: deliberation now hits the non-solo hard-reject:**
  Verified that `track.PolicyFor` now enforces the `availableReviewers < 1` check for all explicit tracks by moving the guard above the track switch. Added `TestExplicitDeliberationNonSoloEscalates` in `internal/driver/track_test.go` confirms that a solo roster on an explicit deliberation track correctly triggers a `trackErr`.
- **[RESOLVED] classifier: unknown/negative size no longer fast:**
  Verified that `track.Classify` now requires `FilesKnown` and `LOCKnown` to be true, and both counts to be non-negative. Smoke tests confirm that `parley classify` with missing size flags defaults to `standard`, and negative inputs result in a clean exit 2 with a descriptive error message.
- **[RESOLVED] fast + auto_implement + --no-implement now escalates:**
  Verified that `driver.New` now reads the idea-level `auto_implement` and `strict_gate` flags directly from the idea directory using `ReadAutoImplement` and `ReadStrictGate`. This ensures the §4.0 contradiction check in `PolicyFor` uses the declared intent rather than the runtime-masked `cfg.AutoImplement` value.
- **[RESOLVED] fast now forces a model-diverse reviewer:**
  Verified that `internal/app/driver_impl.go` now hard-forces `required = true` for model diversity checks when `track: fast` is present, effectively making model diversity a non-optional gate for solo-reviewer tracks.
- **[RESOLVED] standard caps cross-review at 2:**
  Verified that `track.Policy` now includes `CapCrossReviewRounds`, which is set to 2 for the standard track. `driver.New` correctly applies this cap to `cfg.CrossReviewRounds`.
- **[RESOLVED] full-suite green:**
  Ran `go test ./...` and `go test -count=1 ./internal/runner`. All tests passed. The known sandbox limitation in `internal/runner` did not manifest as a failure in this environment (reported as `ok`).

## New findings (if any)

None. The implementation is clean and follows the consensus design.

## Signoff (Status: ✅ ACCEPT)
