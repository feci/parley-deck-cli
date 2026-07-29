---
agent: hermes-1
idea: parley-design-skills
review-round: 05
date: 2026-07-28
reviewed-commit: f1c123d
---

## Summary

🟡 ACCEPT-WITH-RESERVATIONS. Both round-04 CRITICALs are genuinely closed — I re-ran my own escaped-url reproducer and kimi-1's ()/[] nesting probe, and both now fail-safe correctly (exit 4 UNJUDGEABLE, or read correctly). The raw-vs-decoded change introduced no false positives: 17 pass-fixture stylesheets, ordinary CSS, and Tailwind-style escaped selectors all read clean. The doctrine is sound: WCAG 2.2 thresholds are correct, DTCG format is correct, no unfalsifiable rules, no placeholders. D-1 and D-2 are accepted.

The reservation is one new MAJOR. The claim of "zero silent holes" does not hold. I constructed a stylesheet using `@scope` where Chrome applies a declaration the checker neither reads nor marks unreadable, producing a false PASS (exit 0). The root cause is localized: `flushDeclaration()` silently drops buffered text when `currentRule()` returns null, even if an `at` context is on the stack. The fix is small. I would ship with a note.

## What I verified (commands run, and their result)

All probes ran against the skill repo at `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill` (branch `parley-design-skills`, commit `f1c123d`). Browser verification used headless Google Chrome (`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome --headless`). No implementation files were modified.

**Round-04 CRITICAL (escaped-url bypass) — CLOSED.**

My round-04 reproducer: `.probe { background: u\72l(a/*) ; color: #ff0000; */b); }`

- Scanner: `scanStylesheet` returns 0 blocks, `unreadable: ["an unquoted url() token opened at line 1 is still open at end of file", ...]` — exit 4, UNJUDGEABLE.
- Checker end-to-end: `node addons/parley-design-check/bin/check.js /tmp/pds-r5/scope_hole.css` — exit 4, file marked unreadable.
- Browser: Chrome `--headless --dump-dom` on the same CSS wrapped in HTML — `COLOR=rgb(255, 0, 0)`, confirming the browser does apply the color (so the fail-safe correctly refuses rather than passing).

Four escape variants (`u\72l(`, `\75 rl(`, `\000075rl(`, `\75\72\6cl`) were tested. The first three trigger UNJUDGEABLE. The fourth (`\75\72\6cl`) decodes to `urll` not `url`, so the browser also treats `/*` as a comment — not a bypass. All correct.

**Round-04 kimi-1 CRITICAL (()/[] nesting blindness) — CLOSED.**

kimi-1's probe: `.probe { background: url(} ) ; color: #ff0000; } .extra { color: #00ff00; }`

- Scanner: 2 blocks, `.probe` with `background=url(} )` and `color=#ff0000`, `.extra` with `color=#00ff00`. 0 unreadable. The `}` inside `url()` is consumed as content, not structure.
- Browser: `PROBE=rgb(255, 0, 0) EXTRA=rgb(0, 255, 0)` — both applied, matching the scanner's reading.
- Escaped variant `\75rl(})`: scanner reads 2 blocks correctly, 0 unreadable. Closed.

**Attack on "zero silent holes" — FOUND ONE (MAJOR, below).**

I tested 30+ CSS constructs against both the scanner and headless Chrome, including:

- At-rules with bare declarations: `@media`, `@supports`, `@container`, `@layer`, `@scope`, `@starting-style`, `@font-face`, `@property`, `@viewport`, `@counter-style`, `@font-feature-values`, `@page`
- Backslash escapes: `col\6fr`, `\63olor`, `#\66 f0000`, `\2f\2a`, `\2022`, escaped semicolons, escaped comment openers/closers
- Comment interactions: comment inside ident, comment inside value, nested-looking comments, escaped-comment-opener
- String handling: semicolon/colon/brace inside strings, semicolon in `url(";...")`
- Backslash-newline (line continuation) in property, value, and selector
- CSS nesting (with and without `&`), pseudo-elements, `@starting-style` nesting
- `@import`/`@charset`/`@namespace` without semicolons
- Null bytes, form feeds, tabs, `!important` variants
- `var()` with fallback, custom property declarations

The only construct where the browser applies a declaration the checker neither reads nor marks unreadable is `@scope` with bare declarations. Full details in the MAJOR finding below.

**False-positive check (raw-vs-decoded change) — CLEAN.**

- 17 pass-fixture CSS files: all read with 0 unreadable, correct block/declaration counts.
- Ordinary CSS (11 rules, 24 declarations, transitions, keyframes, media queries, custom properties): 0 unreadable, 11 blocks.
- Ordinary CSS with non-adversarial escapes (unicode content `\2022`, escaped class names, attribute selectors): 0 unreadable, 3 blocks.
- Tailwind-style escaped selectors (`.w-1\/2`, `.hover\:bg-red-500:hover`, `.md\:w-1\/3`, `@media` with `.md\:block`): 0 unreadable, 6 blocks, all selectors correctly read.

**Other verifications:**

- `npm test`: 230 pass, 0 fail.
- Byte budgets: SKILL.md 6,519 / 7,168; PDS.md 25,594 / 25,600; RULES.md 23,225 / 24,576; WEB-ANNEX.md 10,022 / 11,264. Total 65,360 / 65,536. All within thresholds.
- Registry digest: computed `b49ff596451f`, matches PDS.md frontmatter `b49ff596451f`.
- Checker without registry: refuses with "no parley-design registry was found and this checker carries no copy of one" — exit 3.
- Checker with nonexistent `--registry`: same refusal.
- Capability: "18 detectors over 18 rule ids, generated from lib/detectors" — confirmed 18 detector files in `lib/detectors/`.
- Placeholder check: `grep -rn "TODO|FIXME|PLACEHOLDER|TBD|XXX|lorem ipsum" addons/parley-design/` — no matches.
- D-2 verification: 8 rules without detectors all report UNJUDGEABLE with clear reasons (not silent PASS).

## Findings

### [MAJOR] @scope bare declarations are a silent hole — browser applies, checker passes

**What is wrong.**

Chrome applies bare declarations directly inside `@scope (selector) { ... }` as if scoped to the scope root. The scanner treats `@scope` as an at-rule context and silently drops any bare declaration inside it — producing 0 blocks and 0 unreadable, which yields a false PASS.

Reproducer:

```css
@scope (.probe) { color: #ff0000; }
.probe { color: var(--color-text-body); padding: var(--space-3); border-radius: var(--radius-panel); }
```

Scanner output (`scanStylesheet`):
- blocks: 1 (the `.probe` rule only)
- unreadable: `[]`
- The `color: #ff0000` inside `@scope` is not in any block.

Browser output (headless Chrome `--dump-dom` with `getComputedStyle`):
- `COLOR=rgb(255, 0, 0)` — Chrome applies the `@scope` bare declaration.

Checker end-to-end (with contract + tokens):
- `verdict PASS — violations 0, needs-review 0, unjudgeable 18`
- `EXIT: 0`

The same gap affects spacing literals: `@scope (.probe) { padding: 20px; }` — Chrome applies `PADDING=20px`, scanner reads 0 blocks for the `@scope` declaration, 0 unreadable.

**Why it matters.**

This breaks the load-bearing claim of the checker: "0 silent holes." The whole point of the fail-safe is that a browser-applied declaration is either read or marked unreadable. Here it is neither. The `@scope` construct is syntactically valid CSS that Chrome (105+, May 2022) accepts and applies, and the checker certifies it as PASS.

**Root cause.**

In `css.js` `flushDeclaration()` (line 531): when `currentRule()` returns null (no `rule` type entry on the stack), the buffer is silently reset — no `log.note`, no unreadable reason. The `else` branch at line 548 only fires when `rule` is truthy but the text is not a valid declaration. When `rule` is null and an `at` context is on the stack, the declaration vanishes without a trace.

**Concrete fix.**

In `flushDeclaration()`, add an `else if` branch: when `rule` is null but `text_` is non-empty and the stack contains an `at` context, emit `log.note(\`the declaration at line ${at} sits inside an at-rule block without a rule, so this scanner cannot tell whether a browser applies it\`)`. This marks the file unreadable, which is the correct conservative behavior — the scanner cannot determine whether the at-rule applies bare declarations.

**Severity rationale.**

MAJOR, not CRITICAL. The construct is non-standard: per the CSS @scope spec, `@scope` contains qualified rules, not bare declarations, so Chrome's behavior is an implementation detail. Real CSS would use `:scope { ... }` or a nested selector, both of which the scanner handles correctly. The gap is localized and the fix is ~3 lines. But it is a genuine silent hole that produces a false PASS, so it cannot be dismissed.

**What I could not test.**

I verified only against Chrome. Firefox and Safari do not yet ship `@scope` as of mid-2024, so I could not cross-check whether the bare-declaration behavior is Chrome-specific or will be interoperable. The checker's claim is about Chromium, and Chromium applies it.

### [MINOR] PDS.md and SKILL.md never received their author's final pass

IMPLEMENTATION.md flags that the `skill+pds` agent died on an API error after writing both files. I read both fully. They are structurally complete (§0–§12 present, changelog written, frontmatter correct, registry digest matches). No placeholders, no truncation, no obvious factual errors. The doctrine is internally consistent and the WCAG 2.2 / DTCG references are correct. I note this as a MINOR because the files are the least-verified artifacts and a future change should give them a full editorial pass, but nothing I found rises to a finding.

## Open questions

1. Is Chrome's application of bare declarations inside `@scope` a spec-compliant behavior or a Chrome bug? If it's a bug that will be fixed, the MAJOR finding may be transient. If it's intended, the fix should ship.

2. The `flushDeclaration` gap (dropping text when `currentRule()` returns null) is the same mechanism that would affect any future at-rule that Chrome decides to apply bare declarations in. Should the fix be the targeted `else if` I propose, or should the scanner more broadly refuse to certify any at-rule block that contains bare declarations it cannot model?

3. D-1: I accept the rebalanced per-file budgets (7/25/24/11 KiB as early-warning thresholds summing above the 64 KiB binding total). PDS.md at 25,594/25,600 has 6 bytes of headroom — tight but valid. Is this sustainable as future amendments add normative text?
