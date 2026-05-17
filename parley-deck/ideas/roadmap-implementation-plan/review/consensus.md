---
idea: roadmap-implementation-plan
drafted-by: codex
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
status: fix-up-needed
---

## Review Summary

Review round 1 found no CRITICAL or MAJOR issues. The implementation delivers the agreed `round-index-artifact-pruning` slice and the test coverage broadly matches `FINAL.md`.

The reviewers identified several low-risk improvements that should be handled before merge because they are small, deterministic, and reduce edge-case ambiguity in the new sanitizer/index code.

## Agreed Fixes

1. Change malformed opening-fence behavior so an unclosed `<think>`, `<thought>`, or `<thinking>` tag does not drop the entire remainder of the artifact.
2. Make case-insensitive fence matching safe for non-ASCII surrounding content; do not index into the original string using offsets from a Unicode-lowercased copy.
3. Remove the dead `### ` guard from H2 extraction or replace it with clearer H2-only logic.
4. Add `os.MkdirAll(roundDir, 0o755)` before writing `_index.md`.
5. Add a stable `token-heuristic: bytes_div_4` field to `_index.md` frontmatter.
6. Replace the warning pseudo-agent ID `index` with a clearly runner-owned sentinel such as `runner/index`.
7. Replace `len([]byte(value))` with `len(value)` in the byte-count token heuristic.
8. Keep `SanitizeForContext` focused on fence removal; avoid undocumented leading/trailing whitespace trimming.
9. Add or adjust tests for no-fence identity, case-mixed fences, malformed opening fences, and the updated warning sentinel/frontmatter behavior.

## Deferred / Non-blocking

- Streaming extraction for very large artifacts is deferred. Current artifact sizes are small enough for whole-file processing.
- Broader Markdown parsing, code-fence-aware sanitizer behavior, and secret redaction remain out of scope for this slice.
- A full golden snapshot test is useful but not required if the existing deterministic assertions are strengthened around the fixed behavior.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The agreed fixes are all local to the new sanitizer/index code and should be handled in fix-up cycle 1 before merge.

### Signoff: claude — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. Agreed fixes 1–9 cover all of my MINOR findings from round-01 (malformed open fence, Unicode-safe fold, dead `### ` guard, MkdirAll, `token-heuristic` frontmatter, `runner/index` sentinel, `len(value)`, trim scoping, and the missing sanitizer/index tests). The two remaining items I raised — `escapeTable("")` returning `"none"` and the `#`/`---` summary-line skip — are NITs I explicitly flagged as non-blocking, so deferring them is fine. Forward-compatibility of `Result.Warning` (a dedicated `Result.Kind` or runner-event type) is also reasonable to defer until a later slice extends the struct.

### Signoff: gemini — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The consensus correctly incorporates the "degrade gracefully" requirement for malformed fences and the Unicode safety improvements. Deferring streaming extraction for now is appropriate given the current scale of the CLI. The agreed fixes 1-9 cover all the minor technical risks identified during my review.

### Signoff: hermes — 2026-05-17
Status: ✅ ACCEPT
Notes: Accept from operations/resilience/tooling lens. Fixes 1-9 (esp. graceful malformed-fence handling, MkdirAll, runner/index sentinel, and expanded sanitizer tests) strengthen runtime robustness and reduce silent-failure surface. No blocking operational concerns remain for this slice.

## Review Cycle 2 Final Consensus

reviewed-commit: 8d57b63

### Review summary

Fix-up cycle 1 applied all nine agreed fixes from review cycle 1. Claude, Gemini, and Hermes completed review round 2. No CRITICAL, MAJOR, MINOR, or blocking findings remain.

Claude noted four NITs that are explicitly non-blocking:

- table rendering of empty cells as `none`;
- summary extraction skipping content lines beginning with `#` or `---`;
- cosmetic CLI alignment for the longer `runner/index` sentinel;
- forward-compatibility if warning-only results are later attached to participant rows.

These are deferred because they do not affect the agreed first-slice behavior or tests.

### Agreed fixes

None for review cycle 2.

### Merge decision

Implementation PR #20 can be merged with the current fix-up commit.
