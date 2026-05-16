---
idea: version-awareness-project-sync
status: final
author: codex
consensus-date: 2026-05-15
participants: [codex, gemini, hermes]
---

## Final plan / specification

Implement Parley Deck version awareness as a small compatibility/status layer across three surfaces:

1. `parley-deck-skill` installer/status CLI.
2. Optional `parley` CLI wrapper/status integration.
3. Project-local `parley-deck/meta/version.json` plus skill startup checks.

Do not force all components to share one lockstep version. Report actual component versions and a compatibility status.

### Version model

Report these as separate facts:

- `installer.version`: the version of the `parley-deck-skill` command being executed.
- `installer.path`: best-effort resolved executable/package path.
- `runtimeInstalls[]`: per-target installed skill status from existing installer markers, including `target`, `dest`, `marker.version`, `marker.source`, `installedAt`, and whether it matches the installer version.
- `parleyCli.version`: best-effort output from `parley version`, when available.
- `project.deckRoot`: detected project deck root.
- `project.protocolPath`: live `parley-deck/COOPERATION.md`, when present.
- `project.protocolSha256`: hash of the live project protocol.
- `project.metadataPath`: `parley-deck/meta/version.json`, when present.
- `project.metadata`: parsed metadata plus schema validation result.
- `compatibility`: `ok`, `warning`, or `blocked`, with concrete reasons.
- `actions[]`: recommended commands or manual steps.

### Packaged compatibility manifest

Add packaged metadata such as `references/compatibility.json`:

```json
{
  "schemaVersion": 1,
  "skillVersion": "1.1.0",
  "protocolSchema": 1,
  "projectMetadataSchema": 1,
  "minimumParleyCli": "1.0.0",
  "compatibleParleyCli": ">=1.0.0 <2.0.0"
}
```

This manifest describes the packaged skill payload and compatibility expectations. It is not the project source of truth.

### Project metadata

Add `parley-deck/meta/version.json`:

```json
{
  "schemaVersion": 1,
  "deckVersion": "1.1.0",
  "protocolSchema": 1,
  "projectMetadataSchema": 1,
  "source": "parley-deck-skill@1.1.0",
  "protocolSha256": "<sha256>",
  "skillSha256": "<sha256-or-null>",
  "compatibilityManifestSha256": "<sha256-or-null>",
  "updatedAt": "2026-05-15T00:00:00.000Z",
  "updatedBy": "parley-deck-skill sync-project"
}
```

This file lets a model answer project-local version questions without guessing from surrounding files. The live `COOPERATION.md` remains canonical; the metadata describes it.

### Installer commands

Keep:

```text
parley-deck-skill --version
```

as simple script output for the system installer command only.

Add:

```text
parley-deck-skill status [--target auto|all|...] [--project .] [--json]
```

as the primary answer to "what Parley Deck version is installed?" In plain text, group output as:

- System installer
- Runtime skill installs
- Parley CLI
- Project deck
- Compatibility
- Recommended actions

In JSON, return a stable object:

```json
{
  "ok": true,
  "command": "status",
  "installer": {},
  "runtimeInstalls": [],
  "parleyCli": {},
  "project": {},
  "compatibility": {},
  "actions": []
}
```

Add:

```text
parley-deck-skill sync-project --project . [--dry-run] [--yes]
```

Rules:

- Dry-run by default unless `--yes` is passed.
- Create or update `parley-deck/meta/version.json` when safe.
- Create missing metadata directories when safe.
- Do not overwrite `COOPERATION.md` by default.
- If packaged protocol differs from live project protocol, report hashes and a short diff summary, then recommend a protocol-change idea or explicit future migration flag.
- If only metadata is stale and protocol hash matches, update metadata with `--yes`.

### Parley CLI integration

Add later or in the same implementation if low-risk:

```text
parley version --all [--json]
```

It should report `parley version` plus delegate to `parley-deck-skill status --project . --json` when the installer command is available. If the installer is unavailable, it should still report CLI version and project metadata by reading local files.

### Skill startup change

Update `SKILL.md` startup flow:

1. Read live `parley-deck/COOPERATION.md`.
2. Check version/status before starting a new Parley workflow:
   - Prefer `parley-deck-skill status --target all --project . --json` when available.
   - If unavailable, read `parley-deck/meta/version.json`, `.parley-deck-skill-install.json` markers in known runtime locations, and hash `COOPERATION.md`.
3. Summarize meaningful drift to the user before invoking participants.
4. Continue automatically for warnings that do not affect protocol correctness.
5. Ask before running `sync-project --yes` or changing project-local protocol files.
6. Record drift that affects the workflow in `parley-deck/inbox/<facilitator>-to-all_<slug>_version-drift.md`.

### Implementation sequence

1. Installer status command:
   - Extend `parseArgs` command set with `status`.
   - Reuse `resolveTargets()` and `targetStatus()`.
   - Add project discovery and hashing helpers.
   - Add `parley` CLI version probing with timeout/failure capture.
   - Add tests for stale installer vs newer runtime marker reporting.

2. Compatibility/project metadata:
   - Add packaged `references/compatibility.json`.
   - Add metadata schema validation.
   - Add `sync-project --dry-run` and `sync-project --yes`.
   - Add tests for metadata creation, stale metadata, protocol hash mismatch, and no-overwrite behavior.

3. Skill startup guidance:
   - Update source `SKILL.md`.
   - Sync installed skill copies and bundled fallback protocol/metadata as required.
   - Add quality gate: before starting Parley Deck, report installer/runtime/project version drift when detected.

4. Optional Parley CLI wrapper:
   - Add `parley version --all --json`.
   - Prefer invoking `parley-deck-skill status`.
   - Fall back to local metadata reads if the installer is unavailable.

5. Release:
   - Bump `parley-deck-skill`.
   - Publish GitHub tag/release and npm when credentials permit.
   - Update Homebrew tap.
   - Verify `parley-deck-skill --version`, `parley-deck-skill status --json`, installer doctor, local runtime markers, and project metadata.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
