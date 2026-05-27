---
idea: antigravity-agent-migration
status: complete
implemented-by: codex
date: 2026-05-27
review-cycle: 1
---

## Implementation summary

- Added Antigravity CLI `agy` as a first-class built-in headless agent using
  the locally verified `agy --print` interface.
- Kept Gemini discoverable as deprecated legacy compatibility and excluded it
  from default participant selection.
- Updated CLI runtime defaults, docs, examples, protocol rosters, project
  shared `agents.toml`, embedded workspace defaults, tests, changelog, and
  version metadata.
- Added the skill installer `agy` target at
  `~/.gemini/config/plugins/parley-deck`, shipped `plugin.json`, and created
  an Antigravity-compatible installed `skills/SKILL.md` copy.
- Kept the Gemini extension payload and target for legacy users.
- Bumped CLI to `1.5.3` and skill package metadata to `1.2.0`.

## Verification

- `go test ./...` in `parley-deck-cli`: passed.
- `npm test` in `parley-deck-skill`: passed.
- Temporary Antigravity install validation:
  `HOME=$(mktemp -d) node bin/parley-deck-skill.js install --target agy --force --json`
  followed by `agy plugin validate <temp-home>/.gemini/config/plugins/parley-deck`:
  passed with one skill and two agents processed.

## Fix-up cycle 1

- Re-scoped CLI default participant selection to exclude only legacy `gemini`.
- Removed Hermes `--yolo` from built-in and shared defaults.
- Documented the `agy` isolated-home limitation and the `agy --print`
  argument-order constraint.
- Annotated the retained legacy Gemini host mapping.
- Review round 02 verified the fixes with no remaining findings.

## Known release work remaining

- Run review cycle 1 and fix agreed findings before publishing.
- Build/release/publish CLI and skill artifacts.
- Update Homebrew formulas after package tarballs or GitHub assets exist.
- Update WinGet manifests only after Windows release asset hashes are known.
