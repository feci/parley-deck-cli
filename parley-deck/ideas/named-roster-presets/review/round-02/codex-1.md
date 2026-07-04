---
agent: codex-1
idea: named-roster-presets
review-round: 2
date: 2026-07-04
reviewed-commit: 8ab0e8c
responding-to: [codex-1/review/round-01, hermes-1/review/round-01]
---

## Position changes since prior review round

I withdraw my round-01 [MAJOR] finding. `internal/app/app.go` now treats
`protocol.ReadRosterIDs(...).ok == false` as a hard error in the preset expansion
path, returns `1`, and stops before `ResolveRoster` or idea creation can silently
skip §2 membership validation.

I agree Hermes's round-01 [CRITICAL] provenance finding is fixed. `CreateIdeaFull`
now writes the roster-preset provenance comment below the frontmatter fence, in the
body, while keeping `track:` as a real frontmatter key. `TestCreateIdeaFullProvenanceOutsideFence`
covers the regression by checking that no provenance/junk key leaks into
`ReadFrontmatter`.

I also agree the unreadable-§2 warning path for `parley preset list` is fixed:
`internal/app/preset.go` now prints a warning that stale-member validation was
skipped when `ReadRosterIDs` cannot parse the roster.

Re-review checks run:
- `go build ./...`
- `go vet ./internal/...`
- `go test ./internal/config ./internal/protocol ./internal/app`
- `go test -count=1 ./internal/protocol -run TestCreateIdeaFullProvenanceOutsideFence`

All passed. My round-01 [MINOR] app-level CLI acceptance test request remains a
reasonable follow-up, but I accept deferring it for this cycle because the corrected
branches are narrow and covered by the resolver, protocol, app package, and targeted
provenance checks.

## Updated findings

zero remaining

## Open questions

None.
