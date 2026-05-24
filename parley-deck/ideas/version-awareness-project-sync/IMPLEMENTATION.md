---
idea: version-awareness-project-sync
status: complete
implementer: codex
started: 2026-05-15
completed: 2026-05-15
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#main-working-tree
head-commit: uncommitted-working-tree
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Implemented the accepted version-awareness design across the skill installer, runtime skill instructions, Parley CLI wrapper, and project-local metadata.

## Implementation plan / checklist

- [x] Add packaged compatibility metadata to `parley-deck-skill`.
- [x] Add `parley-deck-skill status [--project .] [--json]`.
- [x] Add `parley-deck-skill sync-project --project . [--dry-run] [--yes]`.
- [x] Add project metadata generation at `parley-deck/meta/version.json`.
- [x] Add tests for status reporting, project sync, and no-overwrite protocol behavior.
- [x] Add `parley version --all [--json]` wrapper with status and legacy-version fallback.
- [x] Update source and installed skill startup guidance.
- [x] Run package and CLI tests.
- [x] Run Parley review cycle.

## Deviations from FINAL.md

- `parley version --all` now falls back to `parley-deck-skill --version` when the system command is older and does not support `status`; this reports legacy system installer version instead of failing the whole version report.

## Notes for reviewers

- The implementation intentionally keeps `parley-deck-skill --version` simple.
- `COOPERATION.md` remains canonical; sync writes metadata by default, not protocol content.
- Detected runtime skill copies for Codex, Claude, Gemini, and Hermes were updated from the local 1.1.0 source.
- The Homebrew `parley-deck-skill` command on PATH is still a legacy 1.0.8 install until the tap/release is updated; the new CLI reports that explicitly through the legacy fallback.
- 2026-05-24 review round 01 completed with Gemini and Hermes.
- Fix-up cycle 1 addresses Gemini's non-blocking review findings:
  - added `--dir DIR` support to `parley version --all` so project discovery does not require the current working directory;
  - switched version JSON output to the shared indented `printJSON` helper;
  - avoided duplicate missing-command text when both `parley-deck-skill status` and `parley-deck-skill --version` fail.
- Fix-up verification: `env GOCACHE=/tmp/parley-go-cache GOMODCACHE=/tmp/parley-go-mod /tmp/go1.26.3-linux-arm64/go/bin/go test ./internal/app`.
- 2026-05-24 review round 02 completed with no remaining agreed fixes.
- Review consensus is ready with ACCEPT signoffs from codex, gemini, and hermes. Implementation is complete.
