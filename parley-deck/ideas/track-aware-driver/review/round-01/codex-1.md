---
agent: codex-1
idea: track-aware-driver
review-round: 1
date: 2026-07-03
reviewed-commit: a9a0ff4
---

## Summary

I found the track machinery mostly wired through the expected driver/app seams, and the touched-package tests and requested `go vet` pass. I do not consider the implementation ready: the classifier can under-tier unknown/invalid size input to `fast`, `fast` + idea-level `auto_implement` can be bypassed when the run uses `--no-implement`, explicit `track: deliberation` does not participate in the non-solo hard reject, and the required full `go test ./...` gate is red in this environment.

## Refutation attempts

- Acceptance 1, classifier: I traced `internal/app/classify.go` and `internal/track/track.go`, then smoked boundary values with the built CLI. `files=5 loc=300` returns `fast`; `files=6`, `loc=301`, and `files=15 loc=1000` return `standard`; `files=16` and `loc=1001` return `deliberation`. Deliberation flags (`--security`, `--strict-gate`, `--auto-implement`) win over fast-looking size. I broke fail-safe with omitted size and negative LOC: `parley classify --reversible --mechanically-verifiable` returns `fast`, and `--files 1 --loc -1 --reversible --mechanically-verifiable` also returns `fast`.
- Acceptance 2, per-track config and compatibility: `PolicyFor` and `driver.New` set fast to MaxReviewers 1 / MinReviewers 1 / CrossReviewRounds 0 / MaxFixupCycles 1, and explicit standard to reviewer/fix-up caps. Absent track and explicit deliberation apply no overrides and keep the old default CrossReviewRounds/MaxFixupCycles/MinReviewers behavior. The new `ReadTrack(cfg.IdeaDir)` side effect is covered by focused driver unit tests, but not by an app-level regression for all three `driver.New` construction paths.
- Acceptance 3, hard rejects: `fast` + `strict_gate` is read from the idea and should escalate. `fast` + `auto_implement` escalates only when `Config.AutoImplement` is true; app call sites clear that field under `--no-implement`, so the idea-level contradiction can be bypassed. Non-solo is rejected for fast and explicit standard, but not for explicit deliberation.
- Acceptance 4, refutation non-optional: `runner.BuildReviewPrompt` still includes refutation-default instructions, and `runner.ValidateReviewArtifact` requires a non-empty `## Refutation attempts` section. `newDriverImplOps.ReviewRoundComplete` validates every selected reviewer file before completion. I did not find a track flag that disables the validator or sets the reviewer count to zero after a successful fast/standard policy derivation.
- Acceptance 5, backward compatibility and suite: absent and deliberation preserve today's driver knobs. `go test ./internal/track ./internal/driver ./internal/app -count=1` passed, and `go vet ./internal/{track,driver,app}/...` passed. The required `go test ./...` failed persistently in `internal/runner`, so the full-suite acceptance criterion is not green in this review run.
- Reviewer truncation: the app dedupes non-implementer reviewers before truncation, then slices only when `pol.MaxReviewers > 0` and `len(reviewers) > pol.MaxReviewers`, so a successful fast policy does not drop to zero reviewers. Ordering is participant order, which matches the documented implementation deviation but not the stronger model-diversity-preserving ordering discussed in consensus.
- Deviations: driver-internal derivation is acceptable in principle because `driver.New` is the common construction chokepoint, but it is currently coupled to the wrong `AutoImplement` meaning. Deferred per-track timeout and the absent default template are acceptable deviations against FINAL.md as written; they do not break an observable acceptance criterion by themselves.

## Findings

### [MAJOR] Classifier treats unknown or invalid size as fast-eligible

`runClassify` defaults `--files` and `--loc` to `0` (`internal/app/classify.go:21-22`), and `Classify` accepts fast when `Files >= 0`, `Files <= 5`, and `LOC <= 300` without requiring `LOC >= 0` or knowing that the size flags were actually supplied (`internal/track/track.go:95`). As a result, `parley classify --reversible --mechanically-verifiable` returns `fast`, and `parley classify --files 1 --loc -1 --reversible --mechanically-verifiable` also returns `fast`. That violates the fail-safe rule: unknown or nonsensical size is not proof of "<=5 files / <=300 LOC".

Concrete fix: make size knowledge explicit. Use sentinel defaults such as `-1` or add `FilesKnown` / `LOCKnown` to `track.Inputs`, reject negative counts, and only allow fast when both counts are present and non-negative. Add table tests for omitted counts and negative `loc`, plus CLI tests for `--declared fast` on those cases returning exit 4.

### [MAJOR] `fast` + idea-level `auto_implement` can bypass the hard reject with `--no-implement`

The app passes `AutoImplement: driver.ReadAutoImplement(ideaDir) && !noImplement` into `driver.New` (`internal/app/app.go:1167`, `internal/app/app.go:1840`, `internal/app/app.go:1894`). `driver.New` then uses `cfg.AutoImplement` for the track contradiction check (`internal/driver/driver.go:117`). That conflates two different facts: the idea-level `auto_implement: true` is a section 4.0 deliberation trigger, while `--no-implement` is only a runtime brake on code-writing. With `track: fast`, `auto_implement: true`, and `--no-implement`, the driver sees `AutoImplement=false`, applies fast's reduced cross-review policy, and does not escalate the invalid track declaration.

Concrete fix: separate idea metadata from runtime permission. Track validation should use `driver.ReadAutoImplement(ideaDir)` or a distinct `AutoImplementDeclared` field that is not masked by `--no-implement`; keep the masked value only for deciding whether Phase 5/8 may write code. Add app-level regression coverage for `track: fast` + `auto_implement: true` + `--no-implement` escalating before fast policy can finalize the idea.

### [MAJOR] Explicit `track: deliberation` bypasses the non-solo hard reject

The acceptance criterion says the driver hard-rejects a non-solo config with 0 available reviewers, and the consensus says any derived config with 0 available independent reviewers errors. `PolicyFor` only checks `availableReviewers < 1` inside the `Fast` and explicit `Standard` branches (`internal/track/track.go:128-150`). The `Deliberation` branch returns legacy policy without checking reviewer availability (`internal/track/track.go:139-141`), so `driver.New` records no `trackErr` for `track: deliberation` with a one-participant roster (`internal/driver/driver.go:115-120`). Existing tests cover fast non-solo only (`internal/driver/track_test.go:77-84`).

Concrete fix: enforce the non-solo floor before the per-track switch for every explicit track, or at minimum for explicit `track: deliberation` while preserving the documented absent-track legacy path if that compatibility carve-out is intentional. Add `driver.New` tests for explicit deliberation and explicit standard with one participant, and an app-level test that the runtime cannot enter a review path with zero reviewers.

### [MAJOR] Required full-suite verification is not green

The required `go test ./...` command failed twice in this review environment. The persistent failure is `internal/runner TestDurableKillEndToEndRealProcess`: `process verification failed (no recorded boot id); not killed`. The touched-package subset passes, so this may be an existing environment-sensitive runner test, but FINAL.md acceptance criterion 5 and the review prompt require the whole suite to be green.

Concrete fix: make the durable-kill test reliably pass in this supported environment, or explicitly skip it when the platform cannot record a boot id, then rerun `go test ./...` and record the green result in `IMPLEMENTATION.md`.

## Open questions

- Is absent-track legacy allowed to keep bypassing the new non-solo hard reject, or should the non-solo driver gate apply to all ideas regardless of backward compatibility? The current implementation preserves absent legacy behavior, but explicit `track: deliberation` is not clearly exempt from acceptance criterion 3.
- Should `parley classify` require explicit `--files` and `--loc` for any fast result, or should another input source populate those counts before CI uses `--declared`?
