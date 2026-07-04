---
idea: named-roster-presets
status: final
drafter: claude-1
track: standard
date: 2026-07-04
participants: [claude-1, codex-1, hermes-1]
---

## Decision

Add named roster presets as pure, expand-at-creation config sugar. `participants:` in
`00-prompt.md` remains the single canonical quorum; presets never change quorum,
signoff, or track semantics. No protocol-text change.

## Design (as ratified in consensus.md)

1. **Config** — `[rosters.<slug>]` with `participants = ["claude-1", ...]` (§2 canonical
   IDs) parsed in `internal/config` via the existing layered chain (built-ins →
   `~/.parley/agents.toml` → deck `parley-deck/agents.toml`), merged per-preset-name
   (replace, not append). Optional `[defaults.track_rosters]` mapping
   `fast`/`standard`/`deliberation` → preset name, merged per-field.

2. **Resolution** — pure `ResolveRoster(name, track, cfg) ([]string, provenance, error)`:
   explicit `--preset` wins; else explicit `participants:`; else
   `track_rosters[track]` when the track is known at creation; else today's behavior.
   Returns the expanded §2 IDs + a provenance string.

3. **Validation (fail-closed, before idea dir creation)** — new
   `protocol.ReadRosterIDs(root)` reads the §2 active roster; expansion errors on:
   unknown preset, empty/duplicate members, member not in §2, member marked inactive,
   collapse-to-facilitator-only (§1 non-solo), selected track-default naming an unknown
   preset. If the §2 table cannot be parsed, expansion fails (no silent fallback).

4. **CLI** — `--preset <slug>` on the idea-creation path; `--preset` + `--participants`
   = hard error. Expansion writes the concrete `participants:` list into `00-prompt.md`
   plus an HTML-comment provenance line `<!-- roster-preset: <slug> (source: <layer>) -->`.
   A `parley preset list` prints presets with participant counts + source layer + a
   flag for presets referencing inactive/missing agents.

5. **Warnings (non-blocking)** — preset size trips a track reviewer minimum;
   model-diversity risk (LE-3); deck override of a central preset.

## Verification (done criteria)

- `go test ./internal/config/... ./internal/protocol/... ./internal/app/...` green,
  including: central preset inherited when deck has none; deck preset replaces central;
  track defaults merge per-track; unknown preset fails before creation; member not in §2
  fails; member in §2 but not discovered fails without shrinking; `--participants` +
  `--preset` fails; generated `00-prompt.md` has only expanded `participants:` + the
  provenance comment.
- `go build ./...`, `go vet`, `gofmt -l` clean.

## Non-goals

No per-preset model/reasoning/transport overrides; no dynamic rosters; no deck-bootstrap
flow change.

## Signoffs

<!-- each participant appends its own block -->
