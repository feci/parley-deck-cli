---
agent: codex-1
idea: parley-design-skills
review-round: 09
date: 2026-07-28
reviewed-commit: 9121ec2
---

## Summary

❌ BLOCK. Cycle 12 genuinely closes my round-08 false clean: the exact reproducer now reports
both live `--ghost` references, returns `VIOLATION`, does not verify L3, and exits 1. The bound
parameter and declared-default controls both remain clean at verified L3 with exit 0.

The repair did swing toward the opposite failure in the new default-value path.
`functionBinding()` runs `VAR_REFERENCE` over the decoded parameter suffix without first masking
opaque string and URL tokens. As a result, literal text such as `"var(--ghost)"` is reported as a
real custom-property reference. Chromium 150 parses the function as a `CSSFunctionRule`, supports
the call, and computes the literal content `"var(--ghost)"`; the checker instead emits
`core:token-used-undeclared`, rejects L3, and exits 1. This is one localized MAJOR with a shared,
concrete fix, but it is a false gate result and therefore blocks release.

D-1 remains acceptable. The doctrine is unchanged at 65,360/65,536 bytes, all early-warning
thresholds hold, and the registry digest remains correct.

## What I verified (commands run, and their result)

- `git status --short --branch`, `git branch --show-current`, and `git rev-parse HEAD` in the
  implementation repository: clean `parley-design-skills` checkout at
  `9121ec20a05988efae0519fbfd04955cefb283ca`.
- `git show --stat 9121ec2` and the full cycle-12 diff: only `lib/css.js` and
  `test/checker.test.js` changed; the repair adds `functionBinding()` and its regression tests.
- `npm test`: 245 passed, 0 failed. The suite directly exercised registry refusal, generated
  capability, detector fixtures, conformance gating, waiver rejection, exit codes, packaging
  invariants, and installer paths.
- `node addons/parley-design-check/bin/check.js --help`: the documented command and exit-code
  surface is present.
- The cycle-12 focused test by name: 1 passed, 0 failed.
- My round-08 reproducer, unchanged, through both `varUses()` and a copied sound L3 fixture:
  two `--ghost` uses at line 2, no unreadable input, `VIOLATION`, no verified L3, process exit 1.
- The two controls named in round 08 through the same end-to-end CLI path: a body using only its
  bound formal, and a default referencing declared `--color-text-body`, each produced no
  findings, no unreadable input, `PASS`, verified L3, process exit 0.
- Direct scanner controls for top-level commas, nested functions, bracket blocks, strings,
  escaped identifiers, the `returns` clause, no-argument functions, and Sass
  `@function px-to-rem($n)`: the list split correctly, ordinary CSS stayed readable, and the
  Sass parameter lists did not enter the unreadable channel.
- A reverse-direction scanner matrix for quoted and URL-shaped defaults:
  `@function --quoted(--x: "var(--ghost)") ...` and
  `@function --pick(--x: url("var(--ghost)")) ...` each produced a false `--ghost` use with no
  unreadable input. `maskOpaqueTokens()` already blanks both contents correctly, but neither
  `functionBinding()` nor `declarationVarUses()` applies it before `VAR_REFERENCE`.
- The quoted-default case through a copied sound L3 fixture: the checker emitted
  `core:token-used-undeclared`, returned `VIOLATION`, left L3 unverified, and exited 1.
- Chromium 150 in an isolated ego-browser task space: the exact round-08 reproducer computed
  `rgb(255, 0, 0)`; the bound and declared-default controls computed their expected green and
  blue; the quoted-default function was a real `CSSFunctionRule`,
  `CSS.supports("content", "--quoted()")` returned true, and computed pseudo-element content was
  the literal `"var(--ghost)"`. The task space was closed after the probe.
- A sound L3 fixture at current HEAD: `PASS`, verified L3, exit 0, 18 detectors generated over
  18 rule ids, no findings and no unreadable input.
- `wc -c` over the doctrine and a SHA-256 digest check: 6,519 + 25,594 + 23,225 + 10,022 =
  65,360 bytes; `RULES.md` computes `b49ff596451f`, matching PDS frontmatter. There is still
  exactly one `consumeEscape` and one `stringToken` definition.

## Findings

### [MAJOR] Opaque `var()` text in a function default becomes a real token use

**What is wrong:** `functionBinding()` decodes everything after a parameter's formal name and
then applies `VAR_REFERENCE` directly to that whole string. Quoted strings and URL contents are
opaque CSS tokens, but their text remains visible to the regular expression:

```css
@function --quoted(--x: "var(--ghost)") returns <string> {
  result: var(--x);
}
.probe::before { content: --quoted(); }
```

The only real body reference is the bound formal `--x`; `"var(--ghost)"` is literal string
content. Chromium preserves it as that literal, while the checker reports `--ghost` as an
undeclared token. The same unmasked collector pattern is present in `declarationVarUses()`, so
ordinary declaration strings and URL contents also manufacture uses.

**Why it matters:** valid CSS receives a false `core:token-used-undeclared` violation, loses its
L3 certificate, and exits 1. This is the exact reverse-direction regression the round asked to
exclude: avoiding a false clean by treating token-shaped text inside opaque tokens as runtime
references makes the checker an unreliable gate in the other direction.

**Concrete fix:** introduce one value-reference collector that runs over the already decoded
declaration value after `maskOpaqueTokens()` blanks string and URL contents, and use it for both
parameter defaults and ordinary declarations. Keep the present decoding before that step so
escaped real references such as `\76 ar(--ghost)` remain visible. Add end-to-end fixtures for a
quoted default and a URL default that must stay clean, alongside the existing round-08
undeclared-default reproducer, bound-formal control, declared-default control, and escaped-real-
reference control.

## Open questions

None.
