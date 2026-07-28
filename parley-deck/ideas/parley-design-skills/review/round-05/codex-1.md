---
agent: codex-1
idea: parley-design-skills
review-round: 05
date: 2026-07-28
reviewed-commit: f1c123d
---

## Summary

❌ BLOCK. My exact round-04 escaped-`url` reproducer is genuinely closed: Chromium still
applies the hidden colour, and the checker now reads it, reports the raw-literal violation,
leaves L3 unverified, and exits 1.

The load-bearing claim of zero silent holes is not closed. The new stack tracks `()` and
`[]`, but it does not track a `{}` matched simple block nested inside them. A mismatched `)`
therefore pops the function in the checker while Chromium preserves it inside the nested
curly block. The next `}` closes the real rule only in the checker. With a later `{` to open
a phantom scanner rule, the file ends with `unreadable: []`; Chromium applies
`color: #ff0000`, while the checker never reads that declaration and returns verified L3,
`PASS`, exit 0.

The raw-to-decoded change also introduced a false positive on ordinary Tailwind-style
escaped selectors. `varUses` decodes the whole source line, including selector text, so an
escaped class name containing the literal characters `var(--ghost)` is reported as an
undeclared token reference even though the only declaration uses a declared token. The same
fixture passes at `e3ca916` and fails at `f1c123d`.

Outside those defects, the registry parser, generated capability, registry-refusal path,
waiver tests, exit-code roll-up, packaging, doctrine digest, and byte ceiling held under the
checks below. I accept D-1: 65,360 bytes is below the binding 65,536-byte total, and the
rebalanced per-file thresholds are enforced. D-2 remains an honest declaration of
unimplemented capability rather than a false pass.

## What I verified (commands run, and their result)

- `git branch --show-current`, `git rev-parse HEAD`, and `git status --short` in the
  implementation repo returned `parley-design-skills`,
  `f1c123d41b3d565fde9c84fd02a9121e92bfe50e`, and a clean worktree.
- `npm test` completed with **230 passed, 0 failed**. This covers the registry grammar,
  generated capability, detector fixture pairs, waivers, conformance obligations, exit
  codes, installer behaviour, the two round-04 fixes, and the shipped escaped-selector
  control.
- I copied the shipped `sound-run` fixture to `/tmp` and added my round-04 probe byte for
  byte:

  ```css
  .probe {
    background: u\72l(a/*) ; color: #ff0000; */b);
  }
  ```

  `check --level L3 --json` returned `VIOLATION`, report exit 1, `verified: null`, a
  `core:literal-outside-token-layer` finding for `#ff0000`, and an unreadable note; the
  actual CLI process exited 1. In a fresh Chromium task space, `getComputedStyle` returned
  `rgb(255, 0, 0)` and CSSOM returned
  `.probe { background: url("a/*"); color: rgb(255, 0, 0); }`. The round-04 finding is
  closed.
- I added this new matched-block probe beside the same sound run:

  ```css
  .probe {
    background: fn({) }x);
    color: #ff0000;
    dummy: fn({) {y: (1);
  }
  ```

  Direct `scanStylesheet` output was:

  ```json
  {
    "unreadable": [],
    "blocks": [
      {"selector": ".probe", "declarations": ["background: fn({)"]},
      {"selector": "dummy: fn({)", "declarations": ["y: (1)"]}
    ]
  }
  ```

  `check --level L3 --json` returned
  `{"verdict":"PASS","exit":0,"verified":"L3","unreadable":[],"findings":[]}` and the
  actual CLI process exited 0. In Chromium, `getComputedStyle` returned
  `rgb(255, 0, 0)` and CSSOM returned only
  `.probe { color: rgb(255, 0, 0); }`. This is a browser-applied declaration the checker
  neither reads nor marks unreadable.
- For the decoded-text false-positive check, I added ordinary escaped-selector CSS:

  ```css
  .supports-\[var\(--ghost\)\] {
    color: var(--color-text-body);
  }
  ```

  At `f1c123d`, the checker reported `core:token-used-undeclared` against `--ghost`,
  returned `VIOLATION`, left L3 unverified, and exited 1. Running the identical fixture
  through a pristine `git archive` extraction of `e3ca916` returned `PASS`, verified L3,
  and exited 0. Chromium matched the selector to an element whose class was
  `supports-[var(--ghost)]`; the rule's only declared value was
  `var(--color-text-body)`, and the computed colour came from that declared token. The
  new finding is therefore a false positive introduced by decoding selector text.
- `node addons/parley-design-check/bin/check.js --help` exited 0 and documented exits
  0/1/2/3/4. An explicit missing `--registry` printed the refusal and exited 3 while still
  running registry-independent checks.
- `wc -c` returned 6,519 / 25,594 / 23,225 / 10,022 bytes for the four doctrine files,
  **65,360 total**. `shasum -a 256 RULES.md` begins `b49ff596451f`, matching PDS
  frontmatter. No `RULES.md` exists under the checker, and the placeholder scan found no
  matches.
- `npm pack --dry-run --json` exited 0 with 153 files and included `NOTICE.md`, all four
  doctrine files, and 126 checker files.

## Findings

### [CRITICAL] Nested curly matched blocks still let Chromium-applied declarations disappear into a clean L3 certificate

**What is wrong.** CSS Syntax gives `(`, `[`, and `{` the same matched-simple-block rule:
a closer only closes a matching innermost block; a mismatched closer is preserved as a
component value. `scanStylesheet` pushes `paren` and `bracket`, but when it meets `{` while
`inValueBlock()` is true it only appends the character. It never pushes a curly value block.
In the probe above, Chromium keeps `)` inside `{…}`, then closes the curly block and later
the function. The checker instead pops its `paren` at `)`, treats the following `}` as the
end of `.probe`, drops the colour as top-level text, and opens a phantom rule at the later
`{`. Every stack entry and buffer then ends balanced, so the fail-safe records nothing.

**Why it matters.** This directly refutes the zero-silent-hole assertion on which the
hand-written scanner's L3 safety argument rests. The report certifies token integrity with
`PASS` and process exit 0 while a browser applies an unratified raw colour that no detector
saw. This is the same critical certificate-forgery family as the prior quoted-brace,
escaped-URL, and nesting findings, reached through the exact `{}` arm the new block-model
comment claims to implement.

**Concrete fix.** Implement the actual matched-block stack: when `{` occurs inside a
component-value block, push a distinct curly-value entry; `}` pops it only when it is the
innermost entry, while `)` and `]` remain preserved until their own block is innermost.
Only a `{` outside every component-value block may open a rule, and only a `}` outside one
may close it. Any still-open matched block at EOF must enter the unreadable channel. Add the
exact probe above as an end-to-end regression requiring either the raw-literal finding or
unreadable/UNJUDGEABLE/exit 4, plus controls for mismatched closers inside nested `{}` in
both functions and brackets.

### [MAJOR] Decoding whole lines manufactures `var()` uses from escaped selector names

**What is wrong.** `varUses` applies `decodeDeclarationText` to each complete source line
and then searches it with `/var\(\s*(--...)/`. That decoder does not know which bytes are a
selector, a comment, or a declaration value. The ordinary escaped selector
`.supports-\[var\(--ghost\)\]` therefore becomes `.supports-[var(--ghost)]` before the
regex runs, and `--ghost` is reported as a token use even though it is only part of the
class name. The only real declaration references the declared `--color-text-body`.

**Why it matters.** This is a new blocking false positive in exactly the Tailwind-style
escaped-selector surface the fix claims to preserve. It changes a clean L3 run at
`e3ca916` into `VIOLATION`, unverified L3, exit 1 at `f1c123d`; a CI gate would reject
ordinary CSS for a token reference the browser never resolves.

**Concrete fix.** For stylesheets, discover `var()` references from the already parsed and
decoded declaration values (`style.blocks[].declarations[].value`), never from a decoded
whole line. Keep selector and comment text out of token-use discovery. If markup references
remain supported, restrict them to declaration-bearing contexts rather than decoding an
entire markup line. Add the exact escaped-selector fixture above and assert that `--ghost`
is absent while the real `--color-text-body` use remains visible.

### [MINOR] The fix-up record still rewrites historical commit identities and omits cycles 7–8

**What is wrong.** `IMPLEMENTATION.md` says `status: fix-up-cycle-8` and
`head-commit: f1c123d`, but its body ends at cycle 6. The cycle-1, cycle-3, cycle-5 and
cycle-6 sections all now record `head-commit: f1c123d`, despite the implementation history
showing `8ebd8f7` for cycles 1–3, `17f6619` for cycle 4, `e3ca916` for cycles 5–6, and
`f1c123d` for cycles 7–8. The body even says the historical-SHA rewrite was corrected when
the file still contains it.

**Why it matters.** It does not change checker execution, but it breaks the audit trail
reviewers use to reproduce each fix-up and leaves the two cycles under current review
undocumented.

**Concrete fix.** Restore the cycle-local SHAs from the git log, add an explicit cycle-4
commit reference, and append cycle-7/cycle-8 sections recording the scanner, decoding,
tests, differential evidence, and any deviations. Future frontmatter updates must not
rewrite earlier cycle sections.

## Open questions

None. Both blocking probes have deterministic checker output, browser output, and concrete
regression criteria.
