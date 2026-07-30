---
idea: skills-cli-install-path
review-round: 18
agent: hermes-1
date: 2026-07-30
reviewed-commit: 265eb56
---

## Summary

Reviewed `265eb56` in the isolated `wt-hermes` worktree. The suite is 253/253.
All 46 probes in the supplied harness behave as intended: 37 refused, 4 valid
forms green, 3 invisible forms ignored, a hard break left alone, and a
genuinely broken path still RUNS and fails. All six round-17 findings are
confirmed closed. The two shipped documents that drove the occurrence-level
rule still pass. No remaining false-green surface found.

The guard has reached the point where the publication contract is enforced by
construction for the classes that have been identified across eighteen rounds.
I could not produce a command the documentation publishes that would fail for
a reader while the guard reports success. This is a sign-off.

## Per-finding dispositions of round 17

Each finding was verified independently with a standalone Node script that
replicates the guard's `publishedTestCommands` using the actual `commonmark`
parser (0.31.2, locked), plus the harness probes where applicable.

### codex-1 [MAJOR] — visible text belongs to no provenance bucket

**CLOSED.** The guard now builds a complete visible-text view per block,
tracking per-character ownership. Code span literals are emitted to the
visible buffer with a unique `codeNodeId` as owner; prose text gets owner
`-1`. The decision rule is: do both tokens of the occurrence come from one
and the same code node? If yes, the code pass already has it. If no, it is
recorded as prose and refused.

Measured with the actual parser:

```
node `--test` no/such/dir          → prose:REFUSED-PROSE  (node from prose, --test from code span)
no<span></span>de --test no/such/dir  → prose:REFUSED-PROSE  (tag stripped, "node" assembled from prose)
node --te<span></span>st no/such/dir  → prose:REFUSED-PROSE  (tag stripped, "--test" assembled from prose)
<div>no<span></span>de --test no/such/dir</div> → prose:REFUSED-PROSE  (html_block, tags stripped)
```

Probes P32–P35 all produce `pass=11 fail=1  REFUSED-PROSE`.

### codex-1 [MINOR] — a Map keyed by command text masked a prose occurrence

**CLOSED.** `publishedTestCommands` now returns an **array of occurrences**,
not a Map. Each `{command, origin}` pair is recorded independently (deduped
only on the `origin + " " + command` key, so the same text with different
provenance produces two entries). A valid code occurrence can no longer mask
an invalid prose one with identical text.

Measured:

```
node \--test "skills/parley-tracker/bin/*.test.js"     → prose occurrence
`node --test "skills/parley-tracker/bin/*.test.js"`    → code occurrence
Both recorded: [{origin:"prose",...}, {origin:"code",...}]
```

Probe P36 produces `pass=11 fail=1  REFUSED-PROSE`. The fixture assertion at
line 635 confirms `originOf("node --test dup/same.test.js")` is `"code+prose"`.

### kimi-1 [MAJOR] — comments, formatting tags, GFM strikethrough as splice points

**CLOSED.** Three mechanisms:

- **HTML comments:** `html_inline`/`html_block` literals are now processed with
  `node.literal.replace(/<[^>]*>/g, "")` — tags are stripped, not appended.
  `<!-- -->` becomes empty, so `node --t<!-- -->est` joins to `node --test` in
  the visible buffer, is detected as prose, and refused.
- **Formatting tags** (`<ins>`, `<del>`, `<s>`): same stripping. `<ins>e</ins>`
  yields `e` in the visible buffer; `node --t<ins>e</ins>st` becomes `node --test`.
- **GFM strikethrough** (`~~`): `text` nodes are processed with
  `node.literal.replace(/~~/g, "")`. The parser does not implement GFM, so
  `node --t~~e~~st` is one text literal; stripping `~~` yields `node --test`.

Measured:

```
node --t<!-- -->est no/such/dir    → prose:REFUSED-PROSE
node --t<ins>e</ins>st no/such/dir → prose:REFUSED-PROSE
node --t~~e~~st no/such/dir        → prose:REFUSED-PROSE
node ~~--~~test no/such/dir        → prose:REFUSED-PROSE
```

Probes P37–P40 all produce `pass=11 fail=1  REFUSED-PROSE`.

### kimi-1 [MINOR] — invisible text policed, and image alt text executed

**CLOSED.** Two changes:

- **`<script>`/`<style>`/`<template>` bodies:** the guard checks
  `/^\s*<\s*(script|style|template)\b/i` and emits `""` (nothing) for these,
  so their content is neither policed nor executed.
- **Image alt text:** an `insideImage` counter is maintained. When `insideImage > 0`,
  all nodes are skipped. A `code` node inside an `image` node is never processed
  by either pass.

Measured:

```
<!-- maintenance note: node --test ... -->  → NOT-DETECTED (comment stripped to empty, no command visible)
<script> var cmd = "node --test ..."; </script>  → NOT-DETECTED (script body emitted as "")
![`node --test no/such/dir`](img.png)  → NOT-DETECTED (insideImage, skipped)
```

Probes P41–P43 all produce `pass=12 fail=0  GREEN` — the guard does not fail
the build over invisible text, and does not execute a command from alt text.

### hermes-1 [MAJOR] — a soft break renders as a space

**CLOSED.** `softbreak` nodes now emit `" "` instead of `"\n"`. A command
spanning a soft break in prose is one copyable line for the reader and one
line for the guard. It is detected as prose and refused.

Measured:

```
node\n--test no/such/dir  → prose:REFUSED-PROSE  (command: "node --test no/such/dir")
```

Probe P44 produces `pass=11 fail=1  REFUSED-PROSE`. `linebreak` (hard break)
still emits `"\n"` — the reader sees two lines there, and the guard correctly
does not splice them. Probe P45 (hard break) produces `pass=12 fail=0  GREEN`.

### hermes-1 [NIT] ×2 — case-sensitivity and unpinned parser

**CLOSED.**

- **Case:** `mentionsATestCommand` is now `/\bnode\b/i` — case-insensitive
  detection. `Node --test` and `NODE --test` are detected. `SUPPORTED_COMMAND`
  remains case-sensitive (`/^node\s+--test\s+...`), so these are detected and
  then REFUSED by the grammar, not skipped. Probe P46 produces
  `pass=11 fail=1  REFUSED`.
- **Parser pin:** `commonmark` is now pinned exactly at `0.31.2` in both
  `package.json` (devDependencies) and `package-lock.json`. CI runs `npm ci`.

```
package.json:  "commonmark": "0.31.2"
lockfile:      0.31.2
CI:            npm ci
```

## What I verified and how

### Isolation

```
$ git rev-parse HEAD
265eb56b0bfe9a9634b750605853893f23a705c8

$ git status --short
?? node_modules
```

`node_modules` is the supplied pre-existing link. No tracked file was modified.
All temporary scripts (`__verify_round17.js`, `__hunt*.js`) were untracked and
have been removed. Final status is still only `?? node_modules`.

### Suite

```
$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
```

### Harness (46 probes)

```
$ zsh …/scratchpad/probe-hermes.sh
baseline                     pass=12 fail=0
P1–P3                        REFUSED
P4, P7                       GREEN
P5–P6, P9–P18, P20–P25       REFUSED
P8                           RAN-AND-FAILED
P19, P26                     GREEN
P27–P29                      REFUSED-PROSE
P30–P31                      REFUSED
P32–P36                      REFUSED-PROSE
P37–P40                      REFUSED-PROSE
P41–P43                      GREEN
P44                          REFUSED-PROSE
P45                          GREEN
P46                          REFUSED
```

Every classification matches the author's claim: 37 refused, 4 valid forms
green (P4, P7, P19, P26), 3 invisible forms ignored (P41, P42, P43), a hard
break left alone (P45), and a genuinely broken path RAN-AND-FAILED (P8).

### The guard still verifies (not a universal refuser)

P8 (`node --test no/such/dir` in an inline code span) runs through `/bin/sh`,
fails with "could not find," and the guard reports `pass=11 fail=1
RAN-AND-FAILED`. P4, P7, P19, P26 run real tests and assert `pass > 0` from
the child's own summary. The guard executes valid commands and catches real
failures.

### False-positive surface — the two shipped documents

The occurrence-level rule was chosen over a whole-visible-line rule
specifically to avoid breaking two shipped documents. I confirmed both still
pass by parsing them with the actual `commonmark` parser and checking that
the command in each is a single code node whose literal is the whole command:

1. `skills/parley-tracker/templates/subtask.md:68`:
   `  Verify: \`node --test "skills/parley-tracker/bin/*.test.js"\``
   → Code node literal: `node --test "skills/parley-tracker/bin/*.test.js"`
   → Origin: `code`, SUPPORTED_COMMAND: true → WOULD-RUN

2. `skills/parley-tracker/templates/subtask.md:74`:
   `- [ ] AC-3 (Verify: \`node --test "skills/parley-tracker/bin/*.test.js"\`) — COMMIT-SHA`
   → Code node literal: `node --test "skills/parley-tracker/bin/*.test.js"`
   → Origin: `code`, SUPPORTED_COMMAND: true → WOULD-RUN

In both cases the "Verify: " prefix and the parentheses are prose (owner -1);
the command is wholly inside one code span (owner 0). The `nodeOwner` and
`flagOwner` are both 0, so the prose pass skips it. The code pass records it
as `code`. The suite passing at 253/253 confirms both are accepted.

### Additional attack surface explored

I ran 28 standalone probes against the actual parser to look for remaining
false-green surfaces — cases where a reader would copy a runnable broken
command while the guard reports success. Every case fell into one of three
correct outcomes:

- **Detected as prose, REFUSED-PROSE:** entities (`&#32;`, `&nbsp;`, `&#x9;`,
  `&#110;`), `<wbr>` tags, `<br>` tags, `<b>` content tags, `<ins>`/`<del>`
  formatting tags, GFM strikethrough, HTML comments, blockquote/list soft
  breaks, setext headings, link reference definitions, reference link text,
  `<pre><code>` HTML blocks, `<div>` HTML blocks (single and multiline),
  emphasis spanning soft breaks, empty emphasis.
- **Detected as code, WOULD-RUN or REFUSED:** code span with normalized
  newline (CommonMark §6.1), indented code block, fenced block (all forms),
  link with code span text, footnote definition with code span, blockquote
  with code span, code block with `$` prompt.
- **NOT-DETECTED (correctly invisible):** `<script>`/`<style>`/`<template>`
  bodies, image alt text, `<details><summary>` HTML block (entire block is
  one `html_block` node — backtick syntax inside is not parsed by CommonMark;
  tags stripped, content becomes prose → REFUSED-PROSE).

The `<details>` case is worth noting: CommonMark treats the entire
`<details><summary>Run \`node --test no/such/dir\`</summary></details>` as a
single `html_block` node. The backticks are literal content, not a code span
delimiter. The guard strips the tags, gets `Run \`node --test
no/such/dir\``, detects `node --test` in the visible text, records it as
prose, and refuses it. A reader who copies from a GitHub-rendered `<details>`
block gets the text inside `<summary>`, which includes the backticks — so the
command is not cleanly copyable from this form anyway. Refusing it is
fail-closed and correct.

### Owner tracking across block boundaries

I verified that the `flushVisible` call at `html_block` exit and at
paragraph/heading exit correctly resets the `visible` buffer and `owners`
array, so a code span in a paragraph following an HTML block is not
contaminated by the HTML block's prose, and vice versa.

### `--test` word boundary in detection

The detection regex `/--test\b/` matches `--test` followed by a word
boundary. The hyphen in `--test-reporter` is not a word character, so
`\b` exists between `t` and `-`, and `--test-reporter` IS detected.
`SUPPORTED_COMMAND` then refuses it (it requires `\s+` after `--test`, not
`-`). `--testing` is NOT detected (no word boundary between `t` and `i`).
This is correct: `--testing` is not `--test` and would not run the test
runner.

## Findings

No findings. I could not produce a command the documentation publishes that
would fail for a reader while the guard reports success.

The guard has been hardened across eighteen rounds of review. The
publication contract is narrow and checkable: a verification command must be
the whole text of a single code node in canonical `node --test <targets>`
form. The two-pass design (code pass + visible-text pass with per-character
provenance tracking) correctly enforces this contract. The occurrence-level
rule avoids breaking legitimately formatted shipped documents. Detection is
deliberately broader than acceptance, case-insensitive, and fail-closed.

### Signoff: hermes-1 — 2026-07-30
Status: ✅ ACCEPT
