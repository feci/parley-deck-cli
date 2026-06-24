---
agent: antigravity-1
idea: automation-outer-loop
review-round: 2
date: 2026-06-24
---

## Summary

In Round 2, I performed a refutation-style security review of the outer-loop automation changes implemented in response to the agreed fixes (AF1–AF5) in commit `7ff7985`. The primary goal was to verify the robustness of the §14 human brake security boundary under adversarial inputs.

I attempted to break the fixes using several fresh injection vectors and edge cases. While the critical frontmatter injection vector is closed for active exploitation against the current custom parser, I identified one regression (MAJOR) regarding multi-line detail formatting, and one reliability issue (MINOR) regarding poisoned empty directories on failed writes. 

All other fixes (AF2, AF3, AF5) are verified as complete and correct.

## Refutation attempts

I assumed the implemented fixes were incomplete or wrong and tried to bypass them using the following methods:

1. **Adversarial control characters (`\r`, `\x00`, `\t`) in signal fields**: 
   - *Attempt*: Injected carriage returns (`\r`), null bytes (`\x00`), and horizontal tabs (`\t`) in the signal ID.
   - *Result*: `cleanField` successfully mapped all of these to spaces (since they are `< 0x20` or matched explicitly), preventing line-break injection.
2. **Unicode Line/Paragraph Separators (`U+2028` / `U+2029`)**:
   - *Attempt*: Injected `U+2028` and `U+2029` inside the signal ID to see if they could act as line breaks and inject keys downstream.
   - *Result*: `cleanField` does not flatten them because they are $> 0x20$. However, both custom parsers `ReadFrontmatter` (in [workspace.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/workspace.go#L296)) and `readFrontmatterFieldErr` (in [cursor.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/cursor.go#L290)) split lines strictly on `\n`. Thus, the Unicode separators do not cause line splits, and the injected keys are parsed as part of the value of the `source_id` key, resulting in no frontmatter injection.
3. **YAML flow style and key-splitting without newlines**:
   - *Attempt*: Injected `{status: round-01}` and colons (e.g., `source_id: my-id : status: round-01`) without newlines into the ID to split the key.
   - *Result*: Because `ReadFrontmatter` splits on the first colon on each line, flow-style constructs and trailing colons are parsed entirely as the value of the `source_id` key, failing to inject new keys.
4. **Bypassing `validSources`**:
   - *Attempt*: Tested source strings with case variants (`Commit`, `CI`), whitespace, or trailing junk (`commit\njunk`).
   - *Result*: Go map lookup is strict and case-sensitive. Surrounding whitespace is trimmed, but any case changes or trailing junk fail the strict map lookup, causing the tick to reject the signal.
5. **Digest Collision on slugs**:
   - *Attempt*: Analyzed whether boundary shifts (e.g., `a/b` vs `a:b` or separator-shifted signals) could collide under the new slug identity.
   - *Result*: By using `strconv.Quote` on both fields, the concatenated string contains exactly one pair of adjacent double quotes `""` representing the boundary between fields. Since internal quotes are escaped as `\"`, the boundary is mathematically unique and injective. Collision is impossible.

## Findings

### MAJOR — F1: Aggressive sanitization of `Detail` regresses multi-line formatting (AF1 regression)

- **What is wrong**: The `cleanField` helper is applied to `c.Detail` (in [loop.go:210](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L210)). This flattens all newlines, carriage returns, and tabs in the detail field to spaces.
- **Why it matters**: `Detail` is formatted in the markdown *body* of the candidate prompt (`- detail: %s` at [loop.go:232](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L232)), not inside the YAML frontmatter. Since the frontmatter parsers stop reading at the second `---` delimiter, newlines in `Detail` cannot inject frontmatter keys. Flattening `Detail` ruins the readability of legitimate multi-line text (such as stack traces, error logs, or paragraphs) for no security benefit.
- **Concrete fix**: Do not pass `c.Detail` (and optionally `c.Title`, which is also only printed in the body) through `cleanField`. If sanitization is desired for markdown syntax (e.g. preventing raw HTML), it should preserve newlines.

### MINOR — F2: Poisoned empty directory on failed writes skips candidate creation (AF4 completeness)

- **What is wrong**: The slug claim is made via `os.Mkdir` at [loop.go:197](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L197). If this succeeds, but writing the prompt file `00-prompt.md` fails (due to write/permission errors, disk full, or process crash/kill), the directory is left empty. In subsequent ticks, `os.Mkdir` returns `fs.ErrExist` and the function returns `false, nil` (skip), permanently skipping the signal.
- **Why it matters**: A transient write failure or process interruption leaves a "poisoned" empty directory that prevents the loop from ever drafting that candidate.
- **Concrete fix**: Modify `writeCandidate` so that if `os.Mkdir` returns `fs.ErrExist`, it checks if `00-prompt.md` is missing. If it is missing, it should attempt to write the prompt using `os.O_EXCL` so that concurrent ticks serialize safely:
  ```go
  if err := os.Mkdir(dir, 0o755); err != nil {
  	if errors.Is(err, fs.ErrExist) {
  		promptPath := filepath.Join(dir, "00-prompt.md")
  		if _, statErr := os.Stat(promptPath); !errors.Is(statErr, fs.ErrNotExist) {
  			return false, nil // prompt exists, skip
  		}
  	} else {
  		return false, err
  	}
  }
  ```
  Additionally, add a defer block to clean up the directory if the write fails synchronously:
  ```go
  var success bool
  defer func() {
  	if !success {
  		os.RemoveAll(dir)
  	}
  }()
  ```

### MINOR — F3: Digest change breaks deduplication for pre-existing drafts (AF2 regression)

- **What is wrong**: The digest calculation changed from a simple `c.Source + ":" + c.ID` hash to a `strconv.Quote` concatenated hash, and explicit fingerprints are now hashed instead of kept readable.
- **Why it matters**: During migration, any existing candidates drafted under the old slug scheme will not match the new slugs, causing the loop to draft duplicate candidates for the same signals. Additionally, explicit-fingerprint slugs are no longer human-readable.
- **Concrete fix**: If migration path compatibility is important, legacy slug check fallbacks can be introduced in `Tick`, or the transition can be documented as requiring a manual cleanup of old candidate folders.

### NIT — F4: Dead code in `SlugFor`

- **What is wrong**: `SlugFor` checks if `sanitize(c.Source)` is empty and falls back to `"signal"`.
- **Why it matters**: Under the new strict `validSources` validation, empty sources are rejected, making the fallback dead code.
- **Concrete fix**: Remove the fallback check or replace it with a panic/assertion since the source must be validated.

## Open questions

1. **Markdown Escaping**: Should the prompt body variables (like `Detail`) undergo basic markdown escaping to prevent markdown structure breakages (e.g. a signal containing `## Constraints` in its text)?
2. **Connector Schema Evolution**: When live API connectors (GitHub/GitLab) are implemented, will they write to the same `signals.json` file or go through a different boundary validation?
