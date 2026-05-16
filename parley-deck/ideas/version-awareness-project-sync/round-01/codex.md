---
agent: codex
idea: version-awareness-project-sync
round: 1
date: 2026-05-15
---

## Summary

Parley Deck needs an explicit version/status contract rather than relying on each agent to infer state from `package.json`, installer markers, and protocol hashes. The smallest useful design is to add a `status` command to the system installer, persist project-local version metadata under `parley-deck/meta/`, and make the skill startup flow compare those facts before starting a workflow. Updates should be advisory by default and write project-local files only through an explicit `sync` command or user-approved fix.

## Proposed approach

Add a single machine-readable status surface to `parley-deck-skill`:

```text
parley-deck-skill status --target all --project . --json
```

It should report:

- `installer.version`: the package/runtime version of the command being executed.
- `installer.path`: resolved command/package path when available.
- `runtimeInstalls[]`: the existing doctor targets plus marker version/source/installedAt and whether marker version equals installer version.
- `project`: detected deck root, `COOPERATION.md` path, optional `parley-deck/meta/version.json`, protocol hash, bundled fallback hash, and compatibility status.
- `parleyCli.version`: best-effort `parley version` result when available.
- `actions[]`: human-readable recommended fixes, such as "run `parley-deck-skill install --target all --force`" or "run `parley-deck-skill sync-project --project .`".

Keep `doctor` as structural validation and make `status` answer the user's version question. `doctor --json` can include the same version fields later, but `status` should be the one command the model can call when asked "what version do I have?"

Add a project metadata file:

```json
{
  "schemaVersion": 1,
  "deckVersion": "1.0.9",
  "protocolVersion": "2026-05-15",
  "source": "parley-deck-skill@1.0.9",
  "protocolSha256": "<sha>",
  "skillSha256": "<sha>",
  "updatedAt": "2026-05-15T...",
  "updatedBy": "parley-deck-skill"
}
```

Suggested path: `parley-deck/meta/version.json`. It is project-local, canonical for project deck state, and easy for every model to read.

Add a safe sync command:

```text
parley-deck-skill sync-project --project . [--dry-run] [--yes]
```

Default behavior:

- Detect an existing `parley-deck/`.
- Compare `meta/version.json`, live `COOPERATION.md`, and packaged `references/COOPERATION.md`.
- If only metadata is missing or stale and protocol hashes match, update `meta/version.json` automatically with `--yes`; otherwise print a plan.
- If protocol content differs, never overwrite by default. Produce a diff summary and recommend a protocol-change idea or explicit `--write-protocol` flag.
- If structure changes are needed, create missing `meta/` files and directories only when safe.

Update the skill startup flow:

1. Read live `COOPERATION.md`.
2. Run or emulate `parley-deck-skill status --project . --target all --json` when the installer command is available.
3. If system installer version, runtime marker versions, and project deck metadata diverge, summarize the drift before starting.
4. Continue only for non-breaking drift; ask before project sync or protocol overwrite.
5. Record material drift or sync decisions in `parley-deck/inbox/<facilitator>-to-all_<slug>_version-drift.md` when it affects the workflow.

## Concerns / open questions

- The system installer may be stale relative to the source repo, as currently seen with Homebrew `1.0.8` while installed runtime markers are `1.0.9`. `status` must show this clearly instead of treating one source as authoritative.
- `deckVersion` and `protocolVersion` should not pretend semantic compatibility unless we define compatibility rules. Use `schemaVersion` plus exact hashes first.
- Auto-updating `COOPERATION.md` at startup is risky because live project protocol is canonical and may have project-local intentional changes.
- If npm `latest` is blocked or Homebrew is stale, the system still needs a reliable local answer from installed markers and project metadata.

## Risks

- Too many version labels can confuse users. The CLI output should separate "system command", "runtime skill installs", "project deck", and "Parley CLI" explicitly.
- Requiring `parley-deck-skill` to exist for every startup would reduce portability. The skill should use it when available and fall back to reading markers and metadata files directly.
- A sync command can accidentally become a silent protocol migration tool. Keep it dry-run first and require explicit approval for protocol content changes.
