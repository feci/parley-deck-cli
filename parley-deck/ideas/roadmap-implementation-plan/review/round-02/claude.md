---
agent: claude
idea: roadmap-implementation-plan
review_round: 2
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
verdict: COMMENT
---

## Summary

Fix-up cycle 1 cleanly addresses all nine agreed fixes from `review/consensus.md`. The implementation now degrades gracefully on malformed reasoning fences, uses ASCII-only case folding without Unicode offset hazards, creates the round directory before writing `_index.md`, records the `token-heuristic: bytes_div_4` frontmatter field, uses the `runner/index` runner-owned sentinel, drops the dead `### ` guard, swaps to `len(value)` for the byte heuristic, and stops trimming whitespace in `SanitizeForContext`. Tests cover no-fence identity, case-mixed fences, malformed opening fences, the `runner/index` warning path, and the new frontmatter. No CRITICAL, MAJOR, or new MINOR bugs found.

## Findings

No CRITICAL, MAJOR, MINOR, or new-regression findings.

The following are NITs surfaced by re-reading the fixed code; none block merge and several were already flagged as non-blocking in round 1.

- NIT — `internal/runner/round_index.go:268-273`: `escapeTable("")` (and any whitespace-only value) returns the literal string `none`, which renders as the unquoted word `none` in the markdown table. For empty H2 section lists this is intentional via `SectionNames()`, but for an empty artifact name (already remapped to `(none)`) or future empty status values this can be mildly confusing. Pre-existing from round 1 and explicitly accepted as a NIT in consensus.
- NIT — `internal/runner/round_index.go:228-235`: `extractH2Sections` skips summary lines starting with `#` or `---`. The `---` guard exists to skip YAML frontmatter delimiters and horizontal rules, but it also skips legitimate summary paragraphs that happen to start with `---` (uncommon). Pre-existing from round 1, non-blocking.
- NIT — `internal/app/app.go:1318`: `%-8s` width is narrower than the new `runner/index` sentinel (12 chars), so the warning line will be slightly misaligned vs. agent rows. Cosmetic only.
- NIT — `internal/runner/round_index.go:139-141`: `buildRoundIndexEntry` filters `entry.ArtifactName` against `"."` and `string(filepath.Separator)`, but `filepath.Base("")` returns `"."`, so an empty `OutputPath` (e.g., a future warning-only result not from a participant) would render as `(none)`. The current code short-circuits warning rows because `writeRoundIndex` is called before the `runner/index` warning is appended, so this is fine today; flagging it as a forward-compatibility nit if `Result.Warning` is ever set on a participant result.

## Fix-up verification

Mapping each agreed fix to the diff:

1. ✅ Malformed open fence preserved — `round_index.go:34-38` writes `input[start:]` (including the opening fence) when no close is found. Verified by `TestSanitizeForContextPreservesMalformedOpenFence`.
2. ✅ Unicode-safe fold — `indexASCIIEqualFold` / `asciiEqualFold` / `asciiLower` (`round_index.go:43-72`) operate byte-by-byte on the original string and only fold ASCII A–Z. No `strings.ToLower` of the value is taken, so byte offsets remain valid for any UTF-8 input. Verified by `TestSanitizeForContextRemovesCaseMixedFence`.
3. ✅ Dead `### ` guard removed — `extractH2Sections` (`round_index.go:217-241`) checks only `## ` and the summary-skip uses `#`/`---` prefixes. Verified by `TestBuildRoundIndexIsDeterministicAndExtractsH2Only` (the `### Detail` heading must not appear in the index).
4. ✅ `os.MkdirAll(roundDir, 0o755)` — `writeRoundIndex` (`round_index.go:78-82`) creates the directory before `os.WriteFile`.
5. ✅ `token-heuristic: bytes_div_4` frontmatter — emitted in `BuildRoundIndex` (`round_index.go:99`). Verified by `TestBuildRoundIndexIsDeterministicAndExtractsH2Only` assertion `"token-heuristic: bytes_div_4"`.
6. ✅ `runner/index` sentinel — `runner.go:136-145` uses `AgentID: "runner/index"`. Verified by `TestRunRoundOneIndexWriteFailureIsWarning`.
7. ✅ `len(value)` byte heuristic — `approxTokens` (`round_index.go:243-248`) uses `len(value)` directly, no `[]byte(value)` allocation.
8. ✅ No whitespace trim in `SanitizeForContext` — verified by `TestSanitizeForContextPreservesInputWithoutFences`, which asserts byte-equal output for input with leading/trailing spaces.
9. ✅ Tests cover no-fence identity, case-mixed fences, malformed opening fences, updated frontmatter, updated sentinel, and the warning-only event path.

Spot checks for new bugs:

- The warning result is appended after `writeRoundIndex` is called, so the `runner/index` sentinel never appears in the participant table — avoids self-referential rows.
- `MkdirAll` runs before the `WriteFile`, so the warning-path test that pre-creates `_index.md` as a directory still exercises the failure branch (write fails on directory target, not mkdir).
- `removeTaggedBlocks` preserves source AgentID/case in the original input; only the search comparator is folded. Source artifacts remain unmodified (also asserted by `TestBuildRoundIndexIsDeterministicAndExtractsH2Only`).

## Risks / open questions

- None blocking. Deferred items (streaming for very large artifacts, code-fence-aware sanitizer, secret redaction, dedicated `Result.Kind`) remain explicitly out of scope for this slice and are documented in the implementation notes.
