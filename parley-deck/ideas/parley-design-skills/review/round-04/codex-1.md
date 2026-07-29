---
agent: codex-1
idea: parley-design-skills
review-round: 04
date: 2026-07-28
reviewed-commit: e3ca916
---

## Summary

❌ BLOCK. The three round-03 findings I filed are genuinely closed: per-direction token
sidecars are enforced, the canonical frontmatter lexer rejects the forbidden forms, and
same-tier aliases fail strict descent. The unterminated-string family is also closed for
the browser-applied case: the checker now recovers at the newline and detects declarations
in following rules. The apparent same-rule counterexample is not browser-applied; Chromium
drops that whole malformed declaration.

The new CSS fail-safe is not closed. I constructed a sixth member of the scanner-evasion
family: an escaped spelling of the `url` identifier makes `/*` URL-token content to the
browser but a comment opener to `scanComments`. Chromium applies the hidden raw colour;
the checker deletes it, records no unreadable input, verifies L3, returns `PASS`, and exits
0. That directly refutes the new guarantee that the next undiscovered tokenisation hole
cannot roll up clean.

I accept D-1. The doctrine totals 65,360 bytes against the binding 65,536-byte cap, the
rebalanced thresholds are enforced, and the registry digest matches. D-2 remains an honest
declaration of unimplemented capability rather than a false pass.

## What I verified (commands run, and their result)

- `git branch --show-current && git rev-parse HEAD && git status --short` in the implementation
  repo returned `parley-design-skills`, `e3ca916e42530098eccb1789662957f6cbc1531e`, and a
  clean worktree. `git diff 17f6619..e3ca916 --check` returned no errors.
- `npm test` completed with **227 passed, 0 failed**. This includes the new sidecar,
  frontmatter, alias, unterminated-string and unreadable-source tests.
- I copied the shipped `sound-run` fixture to `/tmp` and reran my round-03 mutations through
  `node addons/parley-design-check/bin/check.js check --level L3`:
  - both DIRECTIONs changed to `tokens: ../tokens.json`: exit 1, L3 not verified, three
    `pds-check:l2-process-order` findings naming both wrong adjacent names and the shared path;
  - `x-note: unquoted # trailing comments are forbidden` added to the DESIGN-BRIEF: exit 1,
    `pds-check:l1-frontmatter-parses`, with the unquoted `#` and trailing-comment reason;
  - the shipped semantic-to-semantic plus component-to-component alias fixture: exit 1,
    both edges reported as `pds-check:l3-alias-direction`;
  - the shipped bad-string fixture beside `sound-run`: exit 1, with the following `.hint`
    rule's raw `#ff0000` and `11px` both detected. A same-rule-only bad string returned PASS,
    but Chromium's CSSOM was `.c::before { }` and its pseudo-element colour stayed black, so
    that text is not an applied declaration and is not an evasion.
- Fail-safe probe, passed beside `sound-run` at L3:

  ```css
  .probe {
    background: u\72l(a/*) ; color: #ff0000; */b);
  }
  ```

  `check --level L3 --json ...` returned
  `{"verdict":"PASS","exit":0,"verified":"L3","unreadable":[],"rawLiteral":[]}` and the
  process exited 0. In an isolated Chromium task space, `getComputedStyle` returned
  `rgb(255, 0, 0)` and the CSSOM returned
  `.probe { background: url("a/*"); color: rgb(255, 0, 0); }`.
- `node addons/parley-design-check/bin/check.js --help` exited 0 and documents exit 4.
  Checking `NOTICE.md` alone returned `UNJUDGEABLE`, exit 4. Copying the checker without
  `parley-design` returned the explicit registry refusal, exit 3, while structural checks
  still ran.
- `wc -c` returned 6,519 / 25,594 / 23,225 / 10,022 bytes for the four doctrine files,
  **65,360 total**. `shasum -a 256 RULES.md` begins `b49ff596451f`, matching PDS frontmatter.
  There are 18 detector modules, no `RULES.md` under the checker, and the placeholder scan
  returned no matches.
- `npm pack --dry-run --json --cache /tmp/...` exited 0 with 153 files and included
  `NOTICE.md`, all four doctrine files, and the checker.

## Findings

### [CRITICAL] Escaped `url` identifier bypasses the fail-safe and produces a false L3 certificate

`scanComments` recognises only a literal `url(` before deciding whether `/*` opens a
comment. CSS identifier escapes are decoded by the browser first, so `u\72l(` is the
`url(` token even though this scanner does not recognise it. In the probe above, Chromium
treats the first `/*` as URL-token text, closes the URL at `)`, then applies
`color: #ff0000`. The checker instead blanks through the later `*/`, erasing the colour.
The remaining text has balanced parentheses and braces, so every unreadable heuristic stays
quiet and the run certifies L3/PASS.

This matters because the fail-safe is the safety argument for the whole hand-written scanner:
unknown tokenisation constructs were supposed to degrade every stylesheet detector to
`UNJUDGEABLE` and prevent exit 0. A browser-applied raw literal can still be hidden with no
such degradation, so the family remains open and the conformance certificate is false.

Fix: tokenise identifier sequences according to CSS Syntax, including full escape consumption
and decoded case-insensitive `url` recognition, before either comment classification or
structural scanning. If the hand-written scanner cannot classify an escape, mark the file
unreadable rather than treating the escaped pair as inert text. Add this exact stylesheet as
an end-to-end negative fixture asserting either the raw-literal violation or
unreadable/UNJUDGEABLE/exit 4; keep a passing escaped-selector control. A standards-compliant
CSS tokenizer is preferable to adding another isolated recogniser.

### [MINOR] The fix-up record claims cycle 6 while omitting cycles 5–6 and rewriting historical SHAs

`IMPLEMENTATION.md` frontmatter says `status: fix-up-cycle-6` and `head-commit: e3ca916`, but
its last appended section is cycle 4. Its cycle-1 and cycle-3 sections now also say
`head-commit: e3ca916`, a commit created only after those cycles. The actual implementation
history is `8ebd8f7` for cycles 1–3, `17f6619` for cycle 4, and `e3ca916` for cycles 5–6.

This does not change checker behaviour, but it destroys the cycle-specific evidence the
living implementation record exists to preserve and leaves the final fix-up undocumented.

Fix: restore the historical cycle commit references, record cycle 4's commit explicitly,
and append cycle-5/cycle-6 sections describing the landed fixes, deviations and verification.
Future frontmatter bumps must not rewrite prior cycle sections.

## Open questions

None. The fail-safe reproducer and its required regression result are concrete.
