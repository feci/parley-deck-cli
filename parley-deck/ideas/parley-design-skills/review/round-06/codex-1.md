---
agent: codex-1
idea: parley-design-skills
review-round: 06
date: 2026-07-28
reviewed-commit: aa6b9b3
---

## Summary

❌ BLOCK. Both of my round-05 findings are closed on their exact reproducers, and
hermes-1's `@scope` reproducer is closed too. The nested-curly probe now exposes the applied
colour, returns `VIOLATION`, leaves L3 unverified and exits 1; the Tailwind-style escaped
selector now verifies L3 and exits 0 without manufacturing `--ghost`; and the scoped raw
colour is now found and exits 1.

The replacement block model is nevertheless bypassable. `identLikeToken` still reads `url(`
after `#` or `@` as a URL token, although CSS tokenisation makes the preceding sequence one
hash-token or at-keyword and leaves `(` to open a matched block. Because URL consumption runs
before the stack, the stack never sees those parentheses. The published `#url(` probe still
ends with `unreadable: []`, verifies L3, returns `PASS` and exits 0 while Chromium applies the
hidden raw colour. The stack is correct only over the wrong token stream, so the family is
open.

The `var()` false-positive fix is also incomplete in markup. `markupVarUses` first searches
every raw markup line and only then parses `<style>` contents as CSS. A `var(--ghost)` string
in a `<style>` selector is therefore still reported as an undeclared token reference even
though the only declaration references a declared token.

The new declaration-at-rule modelling itself held under my real-world differential: 338
stylesheets (2.69 MB) produced 0 new findings, 0 verdict changes and 0 unreadability changes
against `aa6b9b3^`. The 11 findings removed by the new declaration-only `var()` search were
all old matches inside comments. `@scope`, `@page`, `@property`, `@counter-style`,
`@font-face`, `@theme` and nested-rule controls behaved in the intended direction. Six
real Tailwind files with `@apply` still fail closed at exit 4, but that limitation is
pre-existing rather than introduced by cycle 9.

I continue to accept D-1: the binding total is 65,360 of 65,536 bytes, and the rebalanced
per-file thresholds remain enforced. D-2 is still disclosed rather than falsely passed.

## What I verified (commands run, and their result)

- In the implementation repo, `git branch --show-current`, `git rev-parse HEAD`, and
  `git status --short` returned `parley-design-skills`,
  `aa6b9b3e6482f2491ef96e4c535eee65320c93a4`, and a clean worktree.
- `npm test` completed with **233 passed, 0 failed**.
- I copied the shipped `sound-run` fixture under `/tmp`, added my round-05 nested-curly
  probe byte for byte, and ran
  `node bin/check.js --registry RULES.md --level L3 --json <run>`. The process exited 1;
  the report returned `VIOLATION`, `verified: null`, named the unbalanced declaration and
  three open blocks, and reported
  `.probe sets color to the colour literal #ff0000`. In isolated Chromium the same source
  still computed `rgb(255, 0, 0)` and CSSOM contained only
  `.probe { color: rgb(255, 0, 0); }`. My CRITICAL reproducer is closed.
- I ran the same direct CLI probe with my exact
  `.supports-\[var\(--ghost\)\]` stylesheet. It returned `PASS`, verified L3, had no
  finding and exited 0. Chromium matched the escaped selector and computed the declaration.
  My round-05 stylesheet false positive is closed.
- I ran hermes-1's `@scope (.probe) { color: #ff0000; }` through the same sound run.
  It returned `VIOLATION`, left L3 unverified, reported the scoped raw colour and exited 1.
  That MAJOR is closed.
- I attacked the token-stream boundary with:

  ```css
  .probe {
    background: x) #url((y)} z);
    color: #ff0000;
    dummy: w) #url((a) {b: (1);
  }
  ```

  Direct `scanStylesheet` output had `unreadable: []`, one `.probe` block containing only
  `background`, and one phantom `dummy: w) #url((a)` block. The end-to-end CLI returned
  `PASS`, `verified: "L3"`, `exit: 0`, no findings, and the process exited 0. Isolated
  Chromium computed `rgb(255, 0, 0)` and CSSOM contained only
  `.probe { color: rgb(255, 0, 0); }`. The same scanner result reproduced for `@url(` and
  escaped `#\75 rl(`.
- For the markup boundary I ran:

  ```html
  <style>
    [data-value="var(--ghost)"] {
      color: var(--color-text-body);
    }
  </style>
  ```

  `markupVarUses` returned both `--ghost@3` and the real
  `--color-text-body@4`. The end-to-end CLI reported
  `core:token-used-undeclared`, left L3 unverified and exited 1. Chromium CSSOM preserved
  `var(--ghost)` only in `selectorText`; the sole declaration value was
  `var(--color-text-body)`.
- I compared current and parent `runCheck` output over all **338 CSS files / 2,688,337
  bytes** under the workspace, each beside the same contract and token graph: **0 new
  findings, 0 verdict/exit changes, 0 unreadability changes**. Eleven old undeclared-token
  findings disappeared; inspection showed every one came from `var()` text in comments.
  A separate scanner count found seven already-unreadable files: six Tailwind `@apply`
  files and one string-escape file, all unchanged from the parent.
- `--help` documented exits 0/1/2/3/4. An explicit nonexistent `--registry` printed the
  refusal, returned overall `UNJUDGEABLE`, and the process exited 3 while structural checks
  still ran.
- `wc -c` returned 6,519 / 25,594 / 23,225 / 10,022 bytes for the four doctrine files,
  **65,360 total**. `shasum -a 256 RULES.md` begins `b49ff596451f`, matching PDS
  frontmatter.
- With a writable temporary npm cache, `npm pack --dry-run --json` exited 0 with 153 files
  and included `NOTICE.md`, all four doctrine files and 126 checker files.

## Findings

### [CRITICAL] Hash- and at-keyword `url` spellings still bypass the matched-block stack and forge L3

**What is wrong.** `identLikeToken` rejects a candidate only when its immediately preceding
character is an ident character. At `#url(` and `@url(` it therefore starts again at the
`u`, calls `unquotedUrl`, and consumes through the first `)`. CSS Syntax instead consumes
`#url` as a hash-token and `@url` as an at-keyword; their following `(` opens a matched
simple block. URL consumption happens before the new stack loop, so the stack never sees
the parentheses that decide what the later braces mean. The probe above desynchronises the
scanner, drops the applied colour as top-level text, opens a phantom rule, finishes balanced,
and leaves the fail-safe silent.

**Why it matters.** This is the exact certificate-forgery family cycle 9 claims to close
with the block model. Chromium applies an unratified raw colour while the checker reads no
colour, records no unreadable input, verifies L3, reports `PASS` and exits 0. A correct stack
over a mis-tokenised stream is not a correct CSS scanner.

**Concrete fix.** Tokenise hash-tokens and at-keywords before ident-like tokens, so an ident
inside either can never independently open a URL token; the narrow equivalent is to refuse
the URL-token path when the ident is immediately bound by `#` or `@`. Add the exact `#url(`,
`@url(` and escaped `#\75 rl(` probes as browser-backed end-to-end regressions requiring the
raw-literal finding or unreadable/UNJUDGEABLE/exit 4, plus clean controls for
`url(#fragment)`, ordinary hash colours and genuine `url(...)`.

### [MAJOR] `<style>` selector text still manufactures undeclared `var()` references

**What is wrong.** `markupVarUses` scans every raw markup line with `VAR_REFERENCE`, including
the contents of `<style>` elements, and afterwards also parses those elements through
`declarationVarUses`. The second pass is correct; the first has already counted selector,
comment and prelude text. In the probe above, `--ghost` exists only inside an attribute
selector string, but it becomes a token use and fails L3.

**Why it matters.** This leaves the round-05 false-positive class open on a normal place CSS
is embedded. CI rejects valid HTML with a token reference the browser never resolves, even
though the real declaration uses a ratified token. The implementation comment explicitly
claims `<style>` contents go through the stylesheet path, but they also go through the
unscoped markup path.

**Concrete fix.** Remove `<style>` contents before the raw-markup search and restrict markup
references to declaration-bearing contexts: parse `<style>` bodies as stylesheets, parse
`style` attributes as declaration lists, and inspect class/className values only for the
supported arbitrary-utility form. Do not search text nodes or unrelated attributes. Add the
exact selector fixture above, plus controls proving inline style and arbitrary-utility
references remain visible.

### [MINOR] The fix-up audit still assigns later commits to earlier cycles

**What is wrong.** `IMPLEMENTATION.md` now includes cycles 7–9, but its claim that every
cycle carries its own commit is still false. Cycles 1 and 3 record `f1c123d` although cycles
1–3 landed at `8ebd8f7`; cycles 5 and 6 also record `f1c123d` although they landed at
`e3ca916`; cycle 2 and cycle 4 carry no cycle-local `head-commit` at all. `git log` shows
`8ebd8f7` for cycles 1–3, `17f6619` for cycle 4, `e3ca916` for cycles 5–6, `f1c123d` for
cycles 7–8, and `aa6b9b3` for cycle 9.

**Why it matters.** It does not affect execution, but it continues to make the audit trail
say earlier fixes were verified at a commit that contains later work, defeating
cycle-by-cycle reproduction.

**Concrete fix.** Record the actual cycle-local commit beside every cycle and leave those
historical fields immutable; keep only the frontmatter `head-commit` moving with current
head.

## Open questions

None. Both blocking reproducers have direct scanner, CLI and Chromium evidence. The
declaration-at-rule risk did not produce a new false positive in my real-world differential.
