---
idea: version-awareness-project-sync
drafted-by: codex
date: 2026-05-15
---

## Agreed decisions

- Use component-specific versions, not one lockstep version number. The system should report actual versions for:
  - `parley-deck-skill` system installer command.
  - installed runtime skill payloads and markers for Codex, Claude, Gemini, Hermes, and other targets.
  - `parley` CLI.
  - project-local deck/protocol state.
- Add a compatibility manifest instead of a global `VERSION` source of truth. The manifest may live as `references/compatibility.json` or equivalent packaged metadata and should include schema version, skill version, protocol schema, project metadata schema, and supported Parley CLI range.
- Add project-local metadata at `parley-deck/meta/version.json` with exact hashes and version/compatibility fields for the live project deck.
- Add `parley-deck-skill status [--target all] [--project .] [--json]` as the primary model-readable and human-readable answer to "which Parley Deck version is installed?"
- Add `parley-deck-skill sync-project --project . [--dry-run] [--yes]` for safe local deck metadata/structure updates.
- Optionally add `parley version --all [--json]` as a CLI wrapper that reports `parley` version plus installer/project status when `parley-deck-skill` is available.
- Update the skill startup flow so facilitators check version/status drift before starting Parley Deck work. The check should warn and recommend actions; it should not silently overwrite canonical project protocol files.

## Agreed trade-offs

- `parley-deck-skill --version` should remain simple and script-friendly, reporting the system command/package version only.
- Rich status belongs in `parley-deck-skill status --json`; `doctor` remains focused on structural validity and can share status fields later.
- Network freshness checks are out of the startup path. Normal startup should work offline using local markers, hashes, and metadata.
- Project-local `COOPERATION.md` remains canonical when present. Packaged fallback protocol and compatibility metadata can recommend a sync, but cannot silently replace project protocol.
- A stale system command with newer runtime markers is a valid drift state that must be reported clearly, not hidden. The currently observed Homebrew `1.0.8` command vs runtime marker `1.0.9` state is an example.

## Open items deferred to implementation

- Exact JSON schema for `references/compatibility.json`.
- Exact JSON schema for `parley-deck/meta/version.json`.
- Whether `sync-project` creates a backup file for protocol changes or relies on git plus dry-run output.
- Whether `parley-deck-skill doctor --json` should embed the new `status` payload or only link to it.
- Whether `parley version --all` is implemented in the same cycle or deferred after installer-side status and sync are stable.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-15
Status: ✅ ACCEPT
Notes: Consensus reflects the converged design: compatibility manifest, project-local version metadata, status/sync commands, and startup drift checks without lockstep versioning or silent protocol overwrites.

### Signoff: hermes — 2026-05-15
Status: ✅ ACCEPT
Notes: Full agreement with compatibility manifest, status/sync commands, and drift checks. No remaining concerns.

### Signoff: gemini — 2026-05-15
Status: ✅ ACCEPT
Notes: The consensus aligns with my round 2 position: pivoting from lockstep versioning to a resilient compatibility manifest. The inclusion of `status` and `sync-project` commands ensures robust environmental health and project continuity without risking silent overwrites of local protocols.
