---
agent: antigravity-1
idea: automation-outer-loop
review-round: 4
date: 2026-06-24
---

## Summary

In Round 4, I re-reviewed the Parley Deck outer-loop automation candidate discovery implementation (`internal/loop/loop.go` and `internal/app/loop_cmd.go`) in refutation mode. 

All Round 3 fixes (**AF10–AF13** and **F3**) are correct, robust, and complete. The original Round 1 critical vulnerability (frontmatter injection) remains fully closed, and the COOPERATION.md §14 human-brake invariant holds securely (no actions beyond drafting `status: candidate` prompts are reachable). 

I found **no new issues** (0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT). The review has converged.

## Refutation attempts

I systematically attempted to break the security boundaries and the Round 3 fixes:

1. **Symlink and File Escape Refutation (AF10 / AF7)**
   - **Planted symlink at `ideas/<slug>`**: `os.Mkdir` fails with `ErrExist`, then `os.Lstat` detects the symlink and rejects it before opening any files. Verified via scratch test script.
   - **Symlink at prompt file `00-prompt.md`**: Checked if we could write outside the directory by planting a symlink (dangling or active) at `00-prompt.md`. `os.OpenFile` with `O_CREATE|O_EXCL` fails with `ErrExist` (verified via scratch test script) and does not write to the symlink target.
   - **Symlink higher up (`ideas/` itself a symlink)**: This is standard resolved path resolution. Because the slug is sanitized to `[a-z0-9-]`, it cannot contain directory traversal elements (`..` or `/`), making directory traversal impossible.
   - **TOCTOU swap between `os.Mkdir` and `os.OpenFile`**: Winning a millisecond race to replace the empty directory with a symlink requires concurrent workspace write permissions, which is outside the untrusted signal injection threat model.
   - **Non-directory file at `ideas/<slug>`**: `os.Mkdir` fails with `ErrExist`, `Lstat` is called, and `!fi.IsDir()` triggers rejection.
   - **Reject path leaks**: Confirmed that the reject path returns early without creating partial files or leaking descriptors.

2. **Detail Column-0 Heading / Frontmatter Injection Refutation (AF11 / AF12 / AF13)**
   - **Unicode line endings**: Verified that all YAML 1.1 line breaks (`\n`, `\r`, NEL `U+0085`, LS `U+2028`, PS `U+2029`) are normalized to `\n` in `indentDetail` and then split, guaranteeing each line gets a four-space prefix.
   - **Alternative whitespace/breaks**: Checked vertical tab (`U+000B`) and form feed (`U+000C`). They are not line breaks in YAML or CommonMark, meaning they cannot start a new line or bypass the four-space prefix.
   - **TrimSpace**: Confirmed `strings.TrimSpace` successfully eliminates leading and trailing spaces/newlines.

3. **Collision Resistance (AF9 / AF2)**
   - Verified that the 128-bit SHA-256 digest (`[:32]` of hex-encoded sum) over a canonical quoted key (`strconv.Quote`) prevents separator-shift collisions and is highly collision-resistant.

4. **Frontmatter Safety and §14 Brake Invariant**
   - Re-confirmed that all signal fields written into the frontmatter go through `cleanField` (which flattens newlines to spaces), preventing YAML key injection.
   - Re-confirmed `internal/loop` has no imports or execution paths to push, merge, run, finalize, or manage quorums.

## Findings

No findings (0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT).

## Open questions

No open questions.
