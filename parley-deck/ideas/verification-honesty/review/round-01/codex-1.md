---
agent: codex-1
idea: verification-honesty
review-round: 1
date: 2026-06-24
---

## Summary

I found two issues: one acceptance miss in LE-3 and one weaker validation gap in LE-1. I did not find a strict_gate close-path bypass for the specified heading-based finding scan, and the RunChecks fail-closed path matches the requested resolution order.

## Refutation attempts

- Read `FINAL.md` and `IMPLEMENTATION.md`, then inspected `git diff 1962cb5..HEAD` for the changed Go files and both `COOPERATION.md` copies.
- LE-2 strict_gate close path: traced `advanceReview` in `internal/driver/impl.go:191-220`. I tried to break completion with an uncertified zero-fix consensus, a certified clean consensus with a real severity heading on disk, and the MaxFixupCycles ceiling. The code opens another round when uncertified, escalates on a real heading despite `strict_gate_clean: true`, and escalates at the configured bound.
- LE-2 placeholder scan: traced `reviewRoundHasFindings` / `scanHasRealFinding` in `internal/driver/impl.go:305-355`. The literal `<title>` placeholder is ignored and concrete `### [CRITICAL|MAJOR|MINOR|NIT] ...` headings are counted.
- LE-4 RunChecks: traced `internal/app/driver_impl.go:152-175` and ran the focused app tests. Explicit `checks:` wins, `go.mod` falls back to `go test ./...`, `auto_implement: true` without `checks:` and without `go.mod` fails closed, and design-only non-Go passes.
- LE-3 model diversity: traced `internal/app/driver_impl.go:184-188` against the acceptance requirement for stdout plus an `agent.model_diversity` event. I found stdout/error behavior but no event append.
- Verification run: `git diff --check 1962cb5..HEAD` passed. `go test ./internal/driver -run 'StrictGate|PhaseReview|PhaseImpl'` passed. `go test ./internal/app -run 'RunChecks|ReviewersShareImplementerModel'` passed. `go test ./internal/driver ./internal/app ./internal/runner` passed for driver/app but failed in runner at `TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`, which appears environment-related and not specific to this diff.

## Findings

### [MAJOR] Model-diversity warning is not recorded as the required event

`FINAL.md` LE-3 requires same-model review opening to emit both a stdout warning and an `agent.model_diversity` event. The implementation in `internal/app/driver_impl.go:184-188` only writes a warning to `o.out`, or returns an error when `require_model_diversity: true`; there is no append to the run event store. This matters because unattended/live auto-drive can pass `io.Discard` for output, so the warning can disappear from the durable audit trail and TUI/state consumers cannot observe the guard.

Concrete fix: append a best-effort `store.Event{Type: "agent.model_diversity"}` through `o.base.Store` before either warning or escalating. Include at least idea slug, implementer, reviewers, shared model, whether the gate was required, and action `warn` or `escalate`. Add tests that call `OpenReviewRound` for both warn and `require_model_diversity: true` paths and then load `events.jsonl` to assert the event exists.

### [MINOR] Review validation accepts a refutation-attempts substring instead of a real section

`ValidateReviewArtifact` in `internal/runner/phase58.go:431-438` uses `strings.Contains(data, "## Refutation attempts")`. A review that has no `## Refutation attempts` heading but mentions that text in a sentence or code block still passes validation, even though LE-1 requires the section so empty-findings reviews show their work.

Concrete fix: parse the file line-by-line and require a real level-2 heading whose trimmed line is exactly `## Refutation attempts`, ideally with the same treatment for `## Findings`. Add a negative test where the phrase appears only in body text and validation must fail.

## Open questions

- Should `MaxFixupCycles` bound absolute review round number, as implemented, or only strict closing attempts after the last fix-up? The current absolute-round bound can escalate at the ceiling instead of opening another clean-closing round.
- Should `scanHasRealFinding` only scan inside the `## Findings` section? The current implementation counts severity headings anywhere in the review file; I did not classify that as a bug because `FINAL.md` defines the scan in terms of headings, but section scoping would avoid false vetoes from quoted examples.
