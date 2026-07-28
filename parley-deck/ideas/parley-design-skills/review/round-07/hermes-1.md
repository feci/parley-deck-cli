---
agent: hermes-1
idea: parley-design-skills
review-round: 07
date: 2026-07-28
reviewed-commit: 076ded5
---

## Summary

✅ ACCEPT. All three round-06 findings are closed: my round-05 `@scope` reproducer
(now a round-06 closure), the `#url(`/`@url(` bypass (codex-1 and kimi-1,
independently), and the `<style>` selector `var()` false positive (codex-1). I
re-ran every reproducer through both the scanner and the end-to-end CLI; each
now either reads the hidden declaration or marks the file unreadable, and the
false positive is gone.

The class-level check this round exists for passes. CSS Syntax Level 3 §4
defines twenty token types; I enumerated every one against the scanner and found
no type that can carry a `{`, `}`, `;` or `:` the scanner reads as structure.
The two remaining OPENs are both in the over-report direction, exactly as
stated. The reverted remedy (narrowing markup `var()` to `style` attributes and
utility brackets) was the right call: it would have lost 1,799 real references
the browser resolves.

D-1 (per-file budget rebalance holding the 64 KiB total): ACCEPT, maintained
for the seventh round. D-2 (enforced-by:check rules without detectors): the
disclosure undercounts by one (`core:contrast-applied` is also undetected) but
in the safe direction — it is T2, reported UNJUDGEABLE regardless.

## What I verified (commands run, and their result)

All probes ran against the skill repo at `parley-deck-skill` (branch
`parley-design-skills`, commit `076ded5`, worktree clean). No implementation
files were modified.

1. `npm test` from the repo root: **237 passed, 0 failed**.
2. `git branch --show-current` / `git rev-parse HEAD` / `git status --short`:
   `parley-design-skills`, `076ded5b44c38a`, clean.

**Round-06 finding 1 — `@scope` bare declarations (my round-05 MAJOR):**

Copied the shipped `sound-run` fixture to `/tmp/pds-r7-hermes/sound-run`,
swapped `panel.css` for my round-05 reproducer (`@scope (.probe) { color:
#ff0000; }` + a clean `.probe` rule), and ran `node
addons/parley-design-check/bin/check.js /tmp/pds-r7-hermes/sound-run --level L3
--json`. Result: `verdict: VIOLATION`, `verified:` (null), exit 1. The finding
reads `core:literal-outside-token-layer — panel.css:1 @scope (.probe) sets
color to the colour literal #ff0000`. Scanner output: 2 blocks, 0 unreadable,
the `@scope` block carries `color:#ff0000`. **CLOSED.**

**Round-06 finding 2 — `#url(`/`@url(` bypass (codex-1, kimi-1):**

Ran all four spellings through the scanner and the end-to-end CLI:

- `#url(`: scanner reads `color:#ff0000` in the `.probe` block, 3 unreadable
  reasons (paren imbalance, 2 blocks never closed, text to EOF). CLI: VIOLATION,
  exit 1, `core:literal-outside-token-layer` naming the raw colour.
- `@url(`: identical — VIOLATION, exit 1, raw colour found.
- `#\75 rl(` (escaped hash): identical — VIOLATION, exit 1.
- `@\75 rl(` (escaped at-keyword): identical — VIOLATION, exit 1.

**All four CLOSED.** The `hashOrAtToken` function (css.js:313) consumes `#` and
`@` tokens before `identLikeToken` can see the `url` ident inside them, so the
`(` opens a matched block the stack tracks, not a url token.

**Round-06 finding 3 — `<style>` selector `var()` false positive (codex-1):**

Ran codex-1's exact reproducer (`<style>[data-value="var(--ghost)"] { color:
var(--color-text-body); }</style>`) through `markupVarUses`. Result: only
`--color-text-body@3` is found. `--ghost` is not reported. **CLOSED.** The
`<style>` contents are now parsed as CSS (declarations only via
`declarationVarUses`), and the raw markup search runs after `<style>` bodies are
blanked.

**Class-level check — CSS Syntax Level 3 §4 token type enumeration:**

Ran 22 targeted probes through `scanStylesheet`, one per §4 token type plus the
bad-token recovery paths and the two declared OPENs:

| §4 type | Scanner handling | Verdict |
|---|---|---|
| ident (§4.3.1) | `identLikeToken`, decoded; `url` opens url token | sound |
| function (§4.3.2) | ident + `(`, consumed before stack | sound |
| at-keyword (§4.3.3) | `hashOrAtToken` before `identLikeToken`, decoded | sound |
| hash (§4.3.4) | `hashOrAtToken` before `identLikeToken` | sound (the fix) |
| string (§4.3.5) | consumed in `scanComments` + `scanStylesheet`; `{`/`}`/`;` inside are text | sound |
| url (§4.3.6) | `unquotedUrl`, consumed before stack; `{`/`}`/`;` inside are text | sound |
| escape (§4.3.7) | `consumeEscape`, decoded in idents and declaration text | sound |
| bad-escape (§4.3.8) | reported unreadable (backslash + newline/EOF) | sound |
| number (§4.3.9) | no structural code points | sound |
| dimension (§4.3.10) | number + unit; escaped unit decoded; non-ident-start after digit guarded | sound |
| percentage (§4.3.10) | number + `%` | sound |
| colon (§4.3.12) | structural (splits declarations) | sound |
| semicolon (§4.3.12) | structural (ends declarations) | sound |
| comma (§4.3.12) | content | sound |
| delim (§4.3.12) | content | sound |
| CDO/CDC (§4.3.13) | read as individual chars; no structural effect | sound |
| unicode-range (§4.3.15) | no structural code points | sound |
| bad-string (§4.3.5) | `badString` flag: `{` opens value block, not rule; recovery at newline | sound |
| bad-url (§4.3.14) | reported unreadable (url token still open) | sound |
| bad-comment | reported unreadable (comment never closed) | sound |

No token type can carry a structural code point the scanner misreads. The
enumeration is complete against §4.

**Two remaining OPENs — verified both over-report:**

- OPEN A: escaped exponent in dimension (`1\65 5` → scanner reads `1e5`,
  Chromium discards as number 1 + invalid unit `e5`). Over-report: a detector
  might flag a value the browser doesn't apply. At the T1 SOURCE tier boundary,
  not a scanner defect. Confirmed.
- OPEN B: declarations Chromium discards as invalid (`rgb(300,300,300)` →
  scanner reads it, Chromium clamps). Over-report: the T1 tier cannot know what
  T2 would compute. Declared, not fixable at T1. Confirmed.

**Reverted remedy verification:**

The current `markupVarUses` (css.js:1096) finds `var()` references in three
places: `<style>` declarations (parsed as CSS), `style` attribute values
(parsed as declaration lists), and raw markup text (class attributes, JSX
expressions, etc.). I verified all three paths find references correctly. The
narrowed remedy codex-1 proposed would have restricted the raw markup search to
`style` attributes and supported utility brackets only; IMPLEMENTATION.md
records it lost 1,799 references across 203 files — real references a browser
resolves (e.g. `h-[var(--radix-select-trigger-height)]`). The revert was the
right call: a false clean is the failure mode this rule exists to prevent, and
losing 1,799 real references to close a false positive in `<style>` selector
text would trade a smaller error for a larger one.

**Other verifications:**

- Registry digest: `sha256(RULES.md)[0:12]` = `b49ff596451f`, matches PDS.md
  frontmatter.
- Byte budget: 6,519 + 25,594 + 23,225 + 10,022 = 65,360 ≤ 65,536. All per-file
  within thresholds.
- Registry refusal: `--registry /tmp/nonexistent` → exit 3, "no parley-design
  registry was found and this checker carries no copy of one", structural
  checks still ran.
- No `RULES.md` under `addons/parley-design-check/` (0 found).
- 18 detector files in `lib/detectors/`, matching "18 detectors over 18 rule
  ids" in the capability output.
- No placeholder text in any shipped design file.
- No non-builtin `require()` anywhere in the checker.
- `atRuleName` decodes escaped at-keywords: `@\73 cope` → `name: "scope"`,
  treated as declaration-at-rule, declarations read.
- WCAG 2.2 thresholds in WEB-ANNEX.md: 4.5:1 (SC 1.4.3), 3:1 large text (SC
  1.4.3), 3:1 non-text (SC 1.4.11), 24×24 CSS px (SC 2.5.8), 320 CSS px reflow
  (SC 1.4.10). All correct. APCA advisory, not blocking. Correct.
- DTCG `2025.10` referenced correctly; `colorSpace` requirement present in both
  PDS.md and RULES.md.
- NOTICE.md credits hallmark (MIT) and impeccable (Apache-2.0) as prior art
  studied, states independent authorship.

## Findings

### [MINOR] D-2 undercounts undetected rules by one

**What is wrong.** D-2 names eight `enforced-by: check`/`both` rules without a
detector. There are nine: `core:contrast-applied` (T2, `both`) is also
undetected and is not named. I verified by parsing the registry and comparing
detector `rule` fields against all `enforced-by: check`/`both` rules.

**Why it matters.** The disclosure is in the safe direction —
`core:contrast-applied` is T2, above the v1 tier, so it is reported UNJUDGEABLE
regardless of detector presence. But the count should match the record: a
reader who recomputes it gets 9, not 8, and the mismatch is the kind of thing
that erodes trust in a disclosure even when it changes nothing.

**Concrete fix.** Add `core:contrast-applied` to D-2's list, or restate the
count as "nine rules" and note the ninth is T2-tier UNJUDGEABLE in any case.

### [NIT] PDS.md has 6 bytes of headroom against its early-warning threshold

**What is wrong.** PDS.md is 25,594 of 25,600 bytes. The next normative
sentence that lands in PDS.md will trip the early-warning threshold, which is a
test failure, not a soft warning.

**Why it matters.** It does not affect this review — the budget holds — but
the next amendment that adds normative text to PDS.md will need to either
rebalance again or compress existing prose.

**Concrete fix.** Plan for it when the next normative change lands, or widen
the PDS.md threshold by a small margin now while there is room in the total.

## Open questions

None. The work has converged. All round-06 findings are closed, the class-level
enumeration is complete, and the two remaining OPENs are honestly declared in
the safe direction.
