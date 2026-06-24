---
agent: antigravity-1
idea: automation-outer-loop
review-round: 3
date: 2026-06-24
---

## Summary

In Round 3, I conducted a refutation-mode review of the outer-loop automation changes for the `automation-outer-loop` idea (Tier 4) in the `parley-deck-cli` repository. I inspected the changes in commit `14f8295` and verified the current state of [internal/loop/loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go) and [internal/app/loop_cmd.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/loop_cmd.go).

Under the refutation model, all fixes (AF6–AF9) and the safety boundary of COOPERATION.md §14 (LE-8 human brake output invariant) were assumed to be incorrect or bypassable until verified otherwise. 

Following extensive testing and code analysis, **no findings (CRITICAL, MAJOR, MINOR, or NIT) were identified.** The implemented fixes are complete, correct, and robust against frontmatter injection, TOCTOU/race conditions, Unicode-based YAML line-break exploits, and deduplication digest collisions.

## Refutation attempts

I attempted to break or bypass each fix and hunt for new regressions introduced in Round 3:

1. **Detail Block Structural Breakout (AF6 Verification)**:
   - *Attempt*: Injected hostile payloads in `Detail` containing YAML frontmatter fences (`\n---\nstatus: round-01\nparticipants: [evil]\n---\n`), markdown headings (`\n## Promotion\n`), and status lines (`\nstatus: round-01`).
   - *Result*: [indentDetail](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L285) correctly splits on newlines and prepends `"    "` (four spaces) to every line. Downstream parsers [ReadFrontmatter](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/workspace.go#L296) and [readFrontmatterFieldErr](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/cursor.go#L290) read line-by-line and return/break immediately upon encountering the second `---` (which closes the actual frontmatter block, statically written at line 9). Since `Detail` is written in the body section (far below the frontmatter block), the frontmatter parsers never scan those lines. 
   - Furthermore, the regular expression parser `reStatus` in [internal/retro/retro.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/retro/retro.go#L50) uses `(?m)^status:\s*(.+)$` to match status at the absolute start of a line. Because every line in the indented detail block starts with four spaces, none of them match the regex, and the parsed status remains `candidate`.

2. **Atomic File Claim Races and Defer Cleanup (AF7 Verification)**:
   - *Attempt*: Analyzed whether concurrent ticks could double-create or clobber the same candidate, or if the `defer` block could clean up a valid prompt.
   - *Result*: The atomic lock is successfully established via `os.OpenFile(promptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)`. At the OS level, `O_EXCL` guarantees that only one process can create and open the file; the other processes receive `fs.ErrExist` and return `false, nil` (skipped), preventing duplicate creation or clobbering. 
   - The `wrote` boolean is only flipped to `true` after the prompt is successfully written. If the write fails, `wrote` remains `false`, and the defer block safely removes the partial file so the next tick can retry. If the write succeeds, the prompt is preserved.

3. **Unicode Line Separators and YAML 1.1 Breaks (AF8 Verification)**:
   - *Attempt*: Injected YAML 1.1 line breaks like Line Separator (`\u2028`), Paragraph Separator (`\u2029`), and Next Line (`\u0085`) into frontmatter fields.
   - *Result*: [cleanField](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L145) explicitly checks all of these characters and maps them to spaces. Other line/control breaks like horizontal tab (`\t`), carriage return (`\r`), Line Feed (`\n`), form feed (`\f`), and vertical tab (`\v`) are also mapped to spaces because they are `< 0x20` or explicitly filtered, preventing any injection.

4. **128-bit Digest Collision and Slug Length (AF9 Verification)**:
   - *Attempt*: Analyzed collision vulnerability and slug length under the 128-bit digest.
   - *Result*: Bumping the digest to 32 hex characters (128 bits of entropy) reduces the collision probability to $2^{-64}$ under birthday constraints, making second-preimage attacks mathematically impossible. The generated slug (e.g. `loop-commit-32hexdigest`) is around 44-50 characters, which is completely acceptable for directory and file names.
   - Within the same run, all signals are processed using the deterministic 32-hex digest, ensuring that duplicate signals are correctly skipped. Since `parley loop` is a brand-new command, no legacy candidate directories exist in decks, meaning there is no risk of losing deduplication state with pre-existing candidates.

5. **Round-01 Critical & §14 Human Brake Invariants**:
   - The round-01 frontmatter injection vulnerability remains fully closed because frontmatter fields (`Source`, `ID`, `Title`) are sanitized via `cleanField` and `Source` is validated against the closed set `validSources`.
   - The loop command and packages do not contain any paths for running deliberations, staffing quorums, implementing, pushing, merging, or finalizing candidates, meaning the COOPERATION.md §14 boundary holds perfectly in code.

## Findings

No findings.

## Open questions

1. **Markdown Indentation and Code Blocks**: Currently, multi-line details are indented with four spaces to render as a preformatted code block in standard markdown. Is there any scenario in downstream TUI displays or parsers where we would prefer a code block fenced with triple backticks (` ``` `) instead of an indented block, or is the 4-space indent behavior preferred for maximum simplicity?
2. **CWD vs Workspace relative signals path**: If a user runs `parley loop tick --signals custom.json`, the path is evaluated relative to the current working directory. This matches standard CLI behavior, but should it be validated or documented if users attempt to pass signals files from outside the deck directory?
