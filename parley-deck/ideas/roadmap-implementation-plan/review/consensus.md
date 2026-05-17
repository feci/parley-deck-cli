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
