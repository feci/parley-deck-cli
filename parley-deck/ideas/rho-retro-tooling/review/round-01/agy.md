---
agent: agy
idea: rho-retro-tooling
review-round: 1
date: 2026-06-16
reviewed-commit: 984c757
---

## Summary

The `parley retro` tooling implementation is structurally clean, read-only by default, and implements a deterministic coreset selection process. However, there are a few major and minor issues in the failure classification logic, blocker detection, and command-line validation that should be fixed before acceptance.

## Findings

### [MAJOR] Inconsistent design-churn classification excludes 2-round ideas
- **What is wrong**: In `internal/retro/retro.go`, the `classify` function (line 134) checks `s.Rounds > 2` to label an idea as `design-churn`. However, the `score` function (line 111) correctly adds a friction weight of 1.0 when `s.Rounds > 1`. Because an idea with `Rounds == 2` has a score of 1.0 but is classified as `low-friction`, it is skipped entirely during coreset selection in `Select` (lines 156 and 173).
- **Why it matters**: Ideas that went to a second design round (representing friction) are completely discarded from the coreset, while review friction of the same scale (2 review rounds) is correctly bucketed as `review-churn` and included. This leads to inconsistent filtering and data loss.
- **Concrete fix**: In `internal/retro/retro.go`, change `s.Rounds > 2` to `s.Rounds > 1` inside the `classify` switch statement.

### [MAJOR] Blocker detection fails to detect Verdict: BLOCK in review files
- **What is wrong**: The `reBlocker` regular expression is defined as `regexp.MustCompile("Status:\\s*❌")`, which matches signoffs in `consensus.md` files. However, individual reviewer files (`review/round-NN/<agent-id>.md`) follow the standard review template where a block is indicated by `Verdict: BLOCK` or `Verdict\nBLOCK` instead of a `Status:` line.
- **Why it matters**: If an implementation is blocked during the active review round, the scanner will fail to detect it as `Blocked` until a final consensus draft is updated, causing the idea to be misclassified in the meantime.
- **Concrete fix**: Update `reBlocker` in `internal/retro/retro.go` to also search for `Verdict:\s*(?:❌|BLOCK)`.

### [MINOR] Missing directory check in retroPropose
- **What is wrong**: In `internal/app/retro.go`, `retroPropose` (line 98) checks if the file `ideas/<slug>/00-prompt.md` exists, but does not check if the directory `ideas/<slug>` itself already exists.
- **Why it matters**: If the directory exists but does not contain `00-prompt.md`, `propose` will proceed and write to it, potentially polluting or overwriting an existing workspace for another idea.
- **Concrete fix**: Check if the idea directory `ideaDir` exists using `os.Stat(ideaDir)` and fail closed if the directory itself exists.

### [NIT] Test helper itoa can overflow single digits
- **What is wrong**: The `itoa` helper function in `internal/retro/retro_test.go` (line 36) is implemented as `func itoa(i int) string { return string(rune('0' + i)) }`.
- **Why it matters**: This returns incorrect characters (like `:`) if the number of rounds is 10 or greater. While current tests use small values, it is brittle.
- **Concrete fix**: Replace it with the standard `strconv.Itoa(i)`.

### [NIT] Unconstrained --slug input validation
- **What is wrong**: `retroPropose` checks `slug != filepath.Base(slug)` and `strings.HasPrefix(slug, ".")` to prevent path traversal, but does not check for spaces or special characters (e.g., quotes, brackets).
- **Why it matters**: Allowing directories with spaces or shell-sensitive characters can break shell scripts and general command-line usage.
- **Concrete fix**: Add validation to ensure the slug strictly contains lowercase alphanumeric characters and hyphens (matching standard kebab-case conventions).

## Open questions

None.

## Verdict

ACCEPT-WITH-FIXES
