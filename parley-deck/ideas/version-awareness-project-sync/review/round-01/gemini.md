---
agent: gemini
idea: version-awareness-project-sync
review-round: 1
date: 2026-05-24
reviewed-commit: f69ab0a
---

## Summary
The implementation successfully delivers the core components of the version-awareness design, including the `parley version --all` integration, the project metadata schema, and the legacy fallback mechanism for older skill installers. The code follows existing patterns in the codebase and includes necessary test coverage for the new functionality.

## Findings

### [MINOR] Missing `--dir` support for `parley version --all`
The `parley version --all` command delegates to `parley-deck-skill status` using the current working directory as the project root. While this matches the "delegate to ... --project ." instruction in `FINAL.md`, it is inconsistent with most other `parley` commands (like `status`, `consensus`, `init`) which allow overriding the workspace root via the `--dir` flag. Users running `parley` from outside a project root will receive a "project not found" error from the skill installer even if they intend to target a specific directory.

### [NIT] Inconsistent JSON formatting in `version` command
The `runVersion` command implements its own `writeVersionJSON` helper instead of utilizing the shared `printJSON` helper found in `internal/app/app.go`. This results in the `parley version --all --json` output being unindented (compact), whereas `parley status --json`, `parley sessions list --json`, and others use 2-space indentation. For consistency, `runVersion` should be refactored to use `printJSON`.

### [NIT] Redundant error messaging in legacy fallback
In `internal/app/version_status.go`, the `legacyParleyDeckSkillStatus` function appends the `statusError` (from the failed `status` command) to the error from the `version` probe. If `parley-deck-skill` is missing entirely from the PATH, both probes will fail with the exact same "executable file not found" error, leading to a redundant error string: `unavailable (exec: "parley-deck-skill": ...; version probe failed: exec: "parley-deck-skill": ...)`. A check to see if the error is identical before appending would improve output clarity.
