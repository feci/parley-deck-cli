---
agent: gemini
idea: roadmap-implementation-plan
review_round: 2
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
verdict: ACCEPT
---

## Summary

The fix-up cycle 1 implementation successfully addresses all issues raised in round 1. The code is much safer, more robust, and properly handles edge cases without dropping critical artifact data.

## Findings

None.

## Fix-up verification

I verified that all 9 agreed fixes were implemented correctly:
1. **Malformed opening-fence:** `removeTaggedBlocks` now explicitly checks if the closing tag exists, and if not, gracefully appends the remaining buffer including the unclosed opening tag (`b.WriteString(input[start:])`).
2. **Unicode safe fold:** `indexASCIIEqualFold` and `asciiEqualFold` were implemented cleanly. Iterating byte-by-byte and using `asciiLower` guarantees that multi-byte UTF-8 sequences won't be corrupted or improperly matched against the ASCII needles.
3. **Removed `### ` guard:** The H2 extraction now properly looks for `## ` and skips other headers cleanly using `strings.HasPrefix(trimmed, "#")`.
4. **`os.MkdirAll`:** `writeRoundIndex` now appropriately creates the directory with `0o755` permissions before attempting to write `_index.md`.
5. **Frontmatter token heuristic:** Added `token-heuristic: bytes_div_4` to the output index.
6. **Warning Sentinel:** The fake agent ID for warnings is now cleanly implemented as `runner/index`.
7. **`len(value)`:** `approxTokens` now avoids the intermediate byte slice allocation and uses `len(value)` directly.
8. **Removed trimming:** `SanitizeForContext` is now solely focused on fence removal, correctly preserving whitespace (verified by the new `TestSanitizeForContextPreservesInputWithoutFences` test).
9. **Tests:** All requested test cases for edge behaviors (identity, case-mixing, malformed fences) and the updated frontmatter/sentinel were added and verify the agreed behaviors.

## Risks / open questions

No remaining risks for this slice. The code is ready to be merged.
