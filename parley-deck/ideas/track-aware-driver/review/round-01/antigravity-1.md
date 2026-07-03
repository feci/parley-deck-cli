---
agent: antigravity-1
idea: track-aware-driver
review-round: 1
date: 2026-07-03
reviewed-commit: a9a0ff4
---

## Summary

The implementation of §4.0 tracks is functionally correct for the `fast` and `standard` paths, and the backward-compatibility for absent tracks is preserved. The `parley classify` tool correctly implements the classification logic with appropriate fail-safes. However, I have found a **[MAJOR]** inconsistency where the explicit `track: deliberation` bypasses the non-solo invariant, and several **[MINOR]** issues regarding the timing of reviewer truncation and boundary conditions.

The suite is green and `go vet` is clean.

## Refutation attempts

- **Absent/Deliberation behaviour:** Confirmed that absent track yields `ApplyOverrides: false`, preserving legacy values for `CrossReviewRounds` and `MaxFixupCycles`. `MinReviewers` defaults to 2, correctly mirroring the legacy LE-11 guard.
- **Classify §4.0 verbatim:** Tested boundary values. `Files=5` is fast, `Files=6` is standard, `Files=15` is standard, `Files=16` is deliberation. `LOC=300` is fast, `LOC=301` is standard, `LOC=1000` is standard, `LOC=1001` is deliberation. Triggers like `security` or `protocol-change` correctly force `deliberation` even for small changes.
- **Hard-rejects:** 
    - `fast + auto_implement` and `fast + strict_gate` correctly escalate via `PolicyFor` error.
    - `non-solo` (0 available reviewers) correctly escalates for `fast` and `standard`.
    - **FAILURE:** `track: deliberation` (explicit) does NOT escalate on solo config. It returns `ApplyOverrides: false` and proceeds, allowing a solo run which violates the stated non-solo invariant for §4.0 tracks.
- **Reviewer truncation:** Confirmed `newDriverImplOps` truncates based on `participants` order. It never drops to 0 because `PolicyFor` (for `fast`/`standard`) ensures at least 1 reviewer is available. 
    - **OBSERVATION:** `round-01` (the idea phase) is NOT truncated. It launches all participants. Truncation only happens once the driver takes over.
- **Refutation non-optional:** Confirmed no changes to validators. `track: fast` still requires a consensus/sign-off, which includes the `## Refutation attempts` check.

## Findings

### [MAJOR] Explicit `track: deliberation` bypasses non-solo hard-reject
The `track.go` comment states: *"No track may drop the non-solo (≥1 independent reviewer) invariant"*. However, `track.PolicyFor` only enforces `availableReviewers < 1` for the `Fast` and `Standard` cases. If a user explicitly declares `track: deliberation` but only one participant is present, the driver will proceed with a solo run (falling back to the implementer as drafter).
**Fix:** In [internal/track/track.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/track/track.go#L532), check for `availableReviewers < 1` before returning the policy for `Deliberation`.

### [MINOR] `round-01` is not track-aware
The reviewer truncation logic lives in `newDriverImplOps`, which is used by the driver. However, `parley run` executes `round-01` using the full `participants` set before the driver is even instantiated. This means a `track: fast` idea with 5 participants will still pay for 5 agents in `round-01`, only truncating to 1 reviewer starting from the consensus/review phases.
**Fix:** Consider deriving the track and truncating the `participants` set in `internal/app/app.go` before calling `runcontrol.Create`.

### [MINOR] Side-effect in `driver.New`
`driver.New` now performs I/O by reading `cfg.IdeaDir`. While necessary for the "single construction chokepoint" design (D2), it introduces a hidden side effect to a constructor that previously only set fields. Existing tests that don't provide a valid `IdeaDir` will silently fall back to legacy behaviour (which is safe), but the dependency on a live filesystem makes the driver harder to test in isolation.

### [NIT] `LOC` can be negative in `classify`
The `track.Classify` logic checks `in.LOC <= 300` for `fast` but does not floor it at 0. `parley classify --loc -100` will be accepted as `fast`.

## Open questions

1. **Why is `CrossReviewRounds` left at `-1` for `Standard`?** `driver.New` defaults it to 1 if `< 0`. If `track: standard` is meant to be explicit, shouldn't it explicitly set it to 1 to ensure it doesn't inherit a 0 from somewhere else?
2. **Backward compatibility for solo runs:** Does the "non-solo invariant" apply to absent tracks too? If yes, `PolicyFor` should check `availableReviewers < 1` even when `present` is false. If no (to preserve legacy), then only the explicit `deliberation` case is missing the check.
