---
agent: codex
idea: version-awareness-project-sync
round: 2
date: 2026-05-15
responding-to: [gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I keep the same core proposal: add a first-class status command, a project-local metadata file, a safe project sync command, and startup checks in the skill. I would refine the version model to avoid treating one semantic version as a global truth across all components. The real problem is compatibility and drift reporting, not forcing the Parley CLI, installer package, runtime skill payloads, and per-project protocol snapshots to always share one version number.

## Responses to others

### @gemini — round-01

I agree with the four-component version model and with exposing a broad `parley version --all` or equivalent. I disagree with "ensure that the `VERSION` file in the CLI repository is used as the source of truth for both `parley` and the `parley-deck-skill` installer." The CLI and skill are separate release artifacts and can validly move at different cadences. The safer implementation is a compatibility manifest, for example:

```json
{
  "schemaVersion": 1,
  "skillVersion": "1.1.0",
  "protocolSchema": 2,
  "minimumParleyCli": "1.0.0",
  "compatibleParleyCli": ">=1.0.0 <2.0.0",
  "projectMetadataSchema": 1
}
```

`parley-deck-skill status` and `parley version --all` can then report actual versions and compatibility status without lying about lockstep.

I agree that failed Homebrew upgrades should surface as "system command is stale; runtime skill markers are newer" rather than being treated as a normal state.

### @hermes — round-01

I agree with Hermes that startup should perform a lightweight check and prompt sync only on drift. I would avoid the phrase "single source-of-truth version file in the packaged skill" unless it means "source of truth for the packaged skill payload only." Project-local `COOPERATION.md` remains canonical when present. The packaged skill can declare what it carries and what project metadata schema it understands; it should not silently override project state.

## New concerns / questions

- The project-local `meta/version.json` should record both exact hashes and compatibility metadata. Hashes are reliable for exact snapshot matching; compatibility fields explain whether drift is actionable.
- The installer should not require npm latest or Homebrew latest checks for normal startup. Network freshness checks can be an explicit `--online` option later.
- `parley-deck-skill --version` should remain simple for scripts, while `parley-deck-skill status --json` should be the rich model-readable surface.

## Current proposal

Adopt:

1. `parley-deck-skill status [--target all] [--project .] [--json]` as the primary version/status report.
2. `parley-deck-skill sync-project --project . [--dry-run] [--yes]` as the project-local metadata/safe structure sync command.
3. Optional `parley version --all --json` in the CLI as a wrapper that includes Parley CLI version plus installer status when available.
4. `parley-deck/meta/version.json` as project-local deck metadata.
5. `agents/manifest.yaml` or a new packaged `references/compatibility.json` as the packaged skill's compatibility declaration.
6. Startup skill text that checks status before starting Parley work, warns on drift, and asks before writing project-local protocol changes.
