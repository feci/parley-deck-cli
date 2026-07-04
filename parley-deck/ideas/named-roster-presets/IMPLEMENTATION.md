---
idea: named-roster-presets
status: implemented
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#roster-presets-design
head-commit: 2474c6f
design-pr: https://github.com/feci/parley-deck-cli/pull/68
implementation-pr: same
---

## Summary of work

Named roster presets as expand-at-creation config sugar. `participants:` in
`00-prompt.md` stays the canonical quorum. New config layer + pure resolver + §2
validator + `parley preset list` + `--preset`/`--track` on `parley run`.

## Implementation plan / checklist

- [x] `internal/config/runtime.go`: `fileConfig.Rosters map[string]rosterOverride` +
      `globalDefaults.TrackRosters map[string]string`.
- [x] `internal/config/roster.go` (new): `LoadRosterPresets` (layered per-preset-name
      merge, central→deck→env) + pure `ResolveRoster(cfg, preset, track, rosterIDs,
      inactive)` returning expanded §2 IDs + provenance; fail-closed on unknown /
      empty / duplicate / not-in-§2 / inactive; track default applied only when no
      explicit preset and track is known.
- [x] `internal/protocol/roster.go` (new): `ReadRosterIDs(root)` parses the §2 roster
      table (first table only, skips the host-handle table), returns active + inactive
      sets + ok=false when unparseable (caller fails closed on membership).
- [x] `internal/protocol/workspace.go`: `CreateIdeaFull(...track, provenance)` writes an
      optional `track:` line + a `<!-- roster-preset: … -->` provenance comment;
      `CreateIdeaWithExclusions` delegates (base output byte-identical when both empty).
- [x] `internal/runcontrol/runcontrol.go`: `CreateOptions.Track` + `.Provenance`
      threaded into `CreateIdeaFull`.
- [x] `internal/app/app.go` (`parley run`): `--preset` + `--track` flags; expansion
      before selection/preflight; `--preset` + `--participants` = hard error; prints the
      resolved roster + provenance; passes Track/Provenance to Create.
- [x] `internal/app/preset.go` (new) + dispatch: `parley preset list` — presets, source
      layer, track hint, and a ⚠ flag for members missing/inactive in §2.
- [x] Tests: `internal/config/roster_test.go` (layering, explicit preset, track default,
      all five fail-closed cases), `internal/protocol/roster_test.go` (table parse,
      inactive detection, host-table isolation, missing file).
- [x] Checks: `go build ./...`, `go vet`, `gofmt -l` clean; `go test ./internal/config
      ./internal/protocol ./internal/runcontrol ./internal/app` green; manual
      `parley preset list` e2e (shown below).

## E2E (manual)

```
$ parley preset list --dir <deck>
Roster presets:
  council        4 agents  (parley-deck/agents.toml)  [track: deliberation]  ⚠ gemini-1 (inactive)
                 claude-1, codex-1, hermes-1, gemini-1
  pair           2 agents  (parley-deck/agents.toml)  [track: fast]
                 claude-1, codex-1
```

## Deviations from FINAL.md

- Added a `--track` flag to `parley run` (FINAL implied "when the creation entry point
  knows the track"). It also seeds the `track:` frontmatter, which the track-aware
  driver already reads — a small, in-scope complement, not a new concept.
- Provenance is written by `CreateIdeaFull` (new sibling of `CreateIdeaWithExclusions`)
  rather than mutating the existing signature — keeps all other callers untouched.

## Notes for reviewers

- `ReadRosterIDs` returns `ok=false` when the §2 table can't be parsed; `parley run`
  then still blocks empty/duplicate presets but skips membership checks (documented
  fail-open ONLY for the membership dimension, never for non-solo/dup/empty). Confirm
  this matches the consensus "fail closed" intent or should hard-stop instead.
- `parley preset list` strips a leading `list` token before flag parsing so
  `--dir` after `list` works (Go flag stops at the first positional).
