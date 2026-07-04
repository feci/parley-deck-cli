---
agent: codex-1
idea: named-roster-presets
review-round: 1
date: 2026-07-04
reviewed-commit: 03a7252
---

## Summary

The core preset expansion path is mostly wired in the right place: `runTask` expands `--preset` / track defaults before participant selection, before preflight, and before idea creation; `--preset` plus `--participants` hard-errors. Layered preset merge is per preset name and replaces the winning participant list. `ReadRosterIDs` is scoped to the first roster table matching the §2 roster shape and detects inactive rows, and `CreateIdeaFull` appears to preserve the legacy prompt bytes when track/provenance are empty.

Verification run:
- `go build ./...` passed.
- `go vet ./...` passed. Exact `go vet` at module root exits with `no Go files in .../wt-roster-presets`, because the root has no Go package files.
- `go test ./internal/config ./internal/protocol ./internal/app ./internal/runcontrol` passed.

## Refutation attempts

- Checked expansion ordering in `internal/app/app.go`: config load and `ResolveRoster` happen before `selectedParticipantIDs`, `runTaskPreflight`, and `runcontrol.Create`; the hard error for `--preset` plus `--participants` is before any expansion side effect.
- Checked merge semantics in `internal/config/roster.go`: each higher layer assigns `out.Presets[name]` with a fresh participant slice, so a deck preset replaces a central preset instead of appending to it.
- Checked `ReadRosterIDs` in `internal/protocol/roster.go` and its tests: it starts only on the roster table header shape, stops when that table ends, and does not parse the later host-handle table; inactive rows are surfaced in the inactive set.
- Checked `CreateIdeaFull` in `internal/protocol/workspace.go`: `CreateIdeaWithExclusions` delegates with empty `track` and `provenance`, and those empty values produce empty interpolations in the existing prompt template.
- Tried to falsify the membership gate using the IMPLEMENTATION.md note about unparseable §2; this produced the MAJOR finding below.

## Findings

### [MAJOR] Preset expansion fails open when the §2 roster cannot be parsed

`FINAL.md` says expansion must fail if the §2 roster table cannot be parsed: "If the §2 table cannot be parsed, expansion fails (no silent fallback)." The implementation does the opposite in the `parley run` preset path. When `protocol.ReadRosterIDs` returns `ok=false`, `runTask` sets `rosterIDs` to nil and continues (`internal/app/app.go:1788`), and `ResolveRoster` only enforces membership when `len(rosterIDs) > 0` (`internal/config/roster.go:110`). That means a malformed or renamed roster table lets a preset containing an installed but non-rostered agent become the canonical `participants:` list, bypassing the §2 quorum authority.

Fix: in the expansion path, treat `!ok` from `ReadRosterIDs` as a hard error before calling `ResolveRoster` or creating any idea. Keep `parley preset list` non-blocking if desired, but print an "unable to validate §2 roster" warning there. Add an app-level test where a preset is configured, `COOPERATION.md` has an unparseable/missing roster table, and `parley run --preset ...` exits nonzero without creating an idea.

### [MINOR] CLI-level preset acceptance criteria are not covered by app tests

The resolver and parser have focused tests, but I could not find `internal/app` tests for the CLI behaviors that are part of the done criteria: `--preset` plus `--participants` hard-error, track-default expansion before selection/preflight, provenance in generated `00-prompt.md`, and the fail-closed behavior when §2 cannot be parsed. The current code appears correct for the first three, but they are exactly the integration points most likely to regress because they depend on call order in `runTask`.

Fix: add `internal/app` tests with fake installed agents and a temp deck config that assert no idea is created for `--preset` plus `--participants`, a track default writes expanded participants plus provenance, and malformed §2 blocks preset expansion after the MAJOR fix.

## Open questions

None.
