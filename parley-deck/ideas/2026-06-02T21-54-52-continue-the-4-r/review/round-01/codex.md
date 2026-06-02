---
agent: codex
idea: 2026-06-02T21-54-52-continue-the-4-r
review-round: 1
date: 2026-06-02
---

## Summary

The decider policy and DAG gate creation keep production gates non-auto-approved, and the additive cursor fields look backward compatible. I found two fail-closed gaps: implementation completion can ignore a blocked review consensus, and stdout fallback can leave an invalid protocol artifact behind.

## Findings

### [MAJOR] Block completion ignores blocked review consensus

What: `blockCompleteFunc` treats an implementation block as complete when `ReviewAgreedFixes` returns `found && count == 0`, but it discards the `blocked` return value at `internal/app/pipeline_cmd.go:1238`. A consensus with `outstanding_agreed_fixes: 0` and `blocked: true` will therefore be considered complete by `Advance`/auto completion checks even though `Phase8Decision` would stop on `blocked` first.

Why: Phase 8 is supposed to fail closed. A reviewer BLOCK should prevent advancement until human resolution, not be bypassed by a zero fix count in the same machine contract.

Fix: Include the blocked value in the completion predicate, e.g. `count, blocked, found, err := pipeline.ReviewAgreedFixes(rc); err == nil && found && !blocked && count == 0`, and add a regression test covering `blocked: true` with zero outstanding fixes.

### [MAJOR] Stdout fallback persists invalid artifacts before validation

What: The stdout fallback writes captured stdout to `outputPath` whenever the trimmed stream merely has prefix `---` (`internal/runner/runner.go:319` to `internal/runner/runner.go:322`), then validates afterward (`internal/runner/runner.go:334` to `internal/runner/runner.go:339`). If stdout starts with `---` but is malformed frontmatter or the wrong artifact type, validation fails but the invalid file remains at the protocol path.

Why: The next run sees the artifact already exists and skips unless overwrite is enabled, so one malformed print-only response can poison retries and leave a bad protocol file in place. That weakens the intended strict stdout fallback behavior.

Fix: Validate stdout content in a temporary file before moving it into `outputPath`, require the first line to be exactly `---` without trimming past leading narration, and delete/reject the temp file on validation failure. Add a regression test where stdout begins with `--- not yaml` or wrong frontmatter and verify no artifact file remains.

## Open questions

None.
