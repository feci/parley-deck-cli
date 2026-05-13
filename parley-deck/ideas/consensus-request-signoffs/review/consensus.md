---
idea: consensus-request-signoffs
cycle: 1
drafted-by: codex
date: 2026-05-13
implementation-pr: https://github.com/feci/parley-deck-cli/pull/12
reviewed-commit: d9337399c238a1e7b72b2d75cf4bbc24590c22ea
---

## Agreed fixes

- Enforce that each invocation adds exactly one new signoff for the expected participant and no signoff for any other participant.
- Detect edits to existing consensus content outside the append-only suffix so the "do not edit any existing line" rule is enforced.
- Print a partial-progress summary before returning a mid-loop failure after one or more successful signoffs.
- Include the `Counter-proposal` field in the prompt's canonical BLOCK example.
- Treat missing/unconfigured/non-installed runner preflight failures consistently as usage/safety failures.
- Adjust the prompt wording to avoid quoted agent IDs in normal prose.
- Add regression tests for forged extra signoff and edited existing consensus content.

## Deferred follow-ups

- Durable child stdout/stderr logs for signoff requests.
- Automatic Git commits after each successful signoff.
- Exporting runner timeout helper behavior instead of mirroring the default timeout logic in app code.
- Cross-process locking for concurrent operators.

## Dismissed findings

- Hermes reported no findings.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-13
Status: ✅ ACCEPT
Notes: The agreed fixes match the review findings and should be addressed in fix-up cycle 1 before re-review.

### Signoff: claude — 2026-05-13
Status: ✅ ACCEPT
Notes: Agreed fixes cover my MAJOR (forged extra signoff) and all MINOR/NIT findings, including the non-signoff edit detection raised in my open questions.

### Signoff: gemini — 2026-05-13
Status: ✅ ACCEPT
Notes: Agreed fixes address my MAJOR finding on file integrity and all MINOR/NIT feedback.
