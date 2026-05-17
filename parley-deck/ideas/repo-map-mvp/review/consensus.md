---
idea: repo-map-mvp
drafted-by: codex
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/22
status: fix-up-needed
---

## Review Summary

Review round 1 found no CRITICAL or MAJOR issues. Gemini and Hermes reported no findings. Claude found three MINOR issues and one NIT that are small and should be fixed before merge.

## Agreed Fixes

1. Ensure `parse_error` never leaks the absolute developer-machine root when a Go file read fails between walk and parse.
2. Add walker coverage for symlink skipping and practical non-regular-file skipping.
3. Add CLI usage coverage for bare `parley context` and an unknown context subcommand.
4. Remove the undocumented `--format md` alias so accepted formats exactly match `markdown|json`.

## Deferred / Dismissed

- `--max-files 0` behavior remains current default fallback behavior for this slice. A stricter numeric validation policy can be designed later if real usage needs it.
- JSON `omitempty` behavior remains accepted for optional fields; the schema lists supported fields, not mandatory emission of empty values.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The agreed fixes are narrow, local, and should be applied in fix-up cycle 1 before merge/release.
