---
agent: antigravity-1
idea: automation-outer-loop
review-round: 5
date: 2026-06-24
---

## Summary

In Round 5, I re-reviewed the Parley Deck outer-loop automation candidate discovery implementation ([loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go) and [loop_cmd.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/loop_cmd.go)) in refutation mode.

I focused on testing and attempting to break the cycle-4 security fixes: the ancestor symlink-escape defenses (**AF14**), and the broad line separator normalization in Detail blocks (**AF15**).

All verification steps and tests pass. I have confirmed that the cycle-4 fixes are complete, robust, and correct. The round-01 frontmatter injection vulnerability remains fully closed, the 128-bit digest is cryptographically collision-resistant, and §14's human-brake code boundary is strictly enforced.

I found **no new issues** (0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT). The review has converged to signoff.

## Refutation attempts

I systematically attempted to break the security boundaries and the Round 4 fixes:

1. **Ancestor Symlink Escape (AF14)**
   - **Symlink at Grandparent / Deck Root / Ancestors**: I analyzed the [assertInsideDeck](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L218-L232) logic. Even if `deck` itself (i.e. `parley-deck/`) is a symlink or any of its ancestors are symlinks, `filepath.EvalSymlinks` resolves both `deck` and `dir` to their absolute canonical paths (`realDeck` and `realDir`). `filepath.Rel` then correctly evaluates the containment of the resolved directory relative to the resolved deck root, successfully blocking any path escape.
   - **Path Traversal Trick**: I checked if a malicious signal could inject traversing components (such as `..` or `/`) into the `slug` or the prompt file path. Since [SlugFor](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L158-L164) enforces a strict sanitization character set `[a-z0-9-]` on the source and appends a hex digest, the slug is guaranteed to be a single directory level, preventing any path traversal.
   - **Dangling Symlinks**: I verified that a dangling symlink at the leaf `00-prompt.md` is rejected. On macOS/Unix, `os.OpenFile` with `O_CREATE|O_EXCL` fails with `EEXIST` when it encounters a symbolic link (dangling or not) and does not traverse/write to the target.
   - **TOCTOU Race Condition**: I considered if a concurrent local process could swap `ideas/<slug>` with a symlink between `assertInsideDeck` and `os.OpenFile`. Since `parley` is a local CLI command running with the user's own permissions (no privilege escalation boundary), a concurrent attacker with write access to the user's workspace already possesses the capability to write files as the user, meaning this does not represent a security/privilege boundary violation.

2. **Detail Column-0 Heading / Frontmatter Injection (AF15)**
   - **Broad Line Splitters**: I tested the broad set of C0 control characters and Unicode separators. I verified that vertical tab (`\v` or `0x0B`), form feed (`\f` or `0x0C`), the C0 info separators (`U+001C` through `U+001E`), `\r`, `\n`, NEL (`U+0085`), LS (`U+2028`), and PS (`U+2029`) are all correctly caught by the switch in [indentDetail](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L346-L377) and normalized to `\n`, which are then split and indented by four spaces.
   - **Other C0 Controls**: I verified that other C0 controls (like `\x1f` or `\x01`) are safely mapped to a space character `' '` (since `r < 0x20`), preventing any hidden control characters from acting as line splitters.
   - **Consecutive and Mixed Separators**: I confirmed that mixed separators (e.g. `\r\n` or `\v\f`) are correctly normalized and split, with all resulting empty or blank lines properly indented (four-space prefixed), preventing any Detail token from reaching column 0.

3. **Digest Collision Resistance (AF9 / AF2)**
   - I re-verified that the 128-bit SHA-256 digest (`[:32]` of hex sum) over the quoted canonical key prevents any boundary-shifting, separator-shifting, or lossy character mapping collisions, and is cryptographically resistant to birthday or second-preimage attacks.

4. **Frontmatter Safety & §14 Brake**
   - Re-confirmed [cleanField](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L145-L153) flattens all newline characters (including Unicode breaks) to spaces for `Source`, `ID`, and `Title`, ensuring no frontmatter key injection is possible.
   - Re-confirmed `internal/loop` has no imports or execution paths that can perform run, push, merge, or finalize operations.

## Findings

No findings (0 CRITICAL, 0 MAJOR, 0 MINOR, 0 NIT).

## Open questions

No open questions.
