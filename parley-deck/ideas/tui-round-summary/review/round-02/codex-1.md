---
agent: codex-1
idea: tui-round-summary
review-round: 2
date: 2026-07-04
reviewed-commit: 2b25f08
responding-to: [codex-1/review/round-01, hermes-1/review/round-01]
---

## Position changes since prior review round

My round-01 [MAJOR] is resolved on this branch. `internal/tui/live.go` now computes
`usedAbove`, reserves the rows below the digest with `reserveBelow`, and only renders
the digest when the remaining budget is at least three rows. The regression test
`TestRenderHomeReservesRunsBelowDigest` is present.

My round-01 [MINOR] is resolved. `internal/driver/digest.go` now caps the extracted
position with `[]rune` before appending the ellipsis, so the driver-side digest text no
longer byte-slices a UTF-8 rune.

Hermes' round-01 MINOR-1 is resolved. `internal/tui/roundsummary.go` renders a leading
`~` marker before truncating the fallback position, so the degraded-extraction signal is
not truncated away.

Hermes' round-01 MINOR-2 is resolved. `internal/driver/driver.go` now reports
`consensus-ready (manual draft)` when `Consensus == nil`, matching the driver's actual
halt at `ActionConsensus`; the wired path still reports `drafting consensus`.

Verification run on commit `2b25f08`:

- `go build ./...` passed.
- `go vet` was run exactly and failed with `no Go files in .../wt-round-summary`; this is
  the repo-root command-shape issue from round 01, not a package vet finding.
- Supplemental `go vet ./...` passed.
- `go test ./internal/driver ./internal/tui` passed.

## Updated findings

zero remaining

The four agreed fixes from round 01 are present in the current tree, and I found no new
implementation finding in the fix-up diff.

## Open questions

- Should future check lists spell the vet command as `go vet ./...` so reviewers do not
  repeatedly hit the root package `no Go files` failure?
