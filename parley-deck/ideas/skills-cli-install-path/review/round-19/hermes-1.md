---
idea: skills-cli-install-path
review-round: 19
agent: hermes-1
date: 2026-07-30
reviewed-commit: f5200f7
---

## Summary

Reviewed `f5200f7` in the isolated `wt-hermes` worktree. Suite is 253/253. All 57
probes in the supplied harness behave as intended: 46 refused (30 REFUSED +
16 REFUSED-PROSE), 4 valid forms green, 5 invisible forms ignored (see NIT below
on the claim's arithmetic), a hard break left alone, and the genuinely broken
path still RAN-AND-FAILED. Both round-18 MAJORs are confirmed closed, both NITs
are confirmed closed, and the MINOR is correctly recorded as a follow-up. The
two shipped documents that drove the occurrence-level rule still pass. No
remaining false-green surface found.

I hunted for a fourth layer beyond markdown structure, markdown rendering, and
shell word construction. I tested HTML entities (numeric and named), zero-width
characters, CDATA sections, processing instructions, Unicode homoglyphs,
reference links, setext headings, autolinks, and tab-separated tokens. Every
case fell into one of three correct outcomes: detected and refused, detected and
valid, or not detected because the form is genuinely invisible or non-runnable.
No fourth layer exists.

The approximation-versus-lexing choice is sound. I explain why below.

## Per-finding dispositions of round 18

### codex-1 [MAJOR] — a quoted `>` defeats the tag stripper

**CLOSED.** Cycle 22 replaced the `/<[^>]*>/g` tag stripper with two mechanisms:

1. `html_inline` is dropped whole (`emit("", -1)`). An inline node is one
   complete piece of markup, so its literal is never page text. This eliminates
   the quoted-`>` problem entirely for inline HTML.

2. `html_block` is scanned by `visibleTextOfHtml`, a tokenizer that understands
   quoted attributes, comments, CDATA, and processing instructions. It tracks
   quote state when scanning for `>`, so `<span title="1 > 0">` ends at the
   second `>`, not the first.

Measured with the actual parser against all three forms from the finding:

```
no<span title="1 > 0"></span>de --test no/such/dir
  → occ=[{"command":"node --test no/such/dir","origin":"prose"}]  REFUSED-PROSE

node --t<!-- a > b -->est no/such/dir
  → occ=[{"command":"node --test no/such/dir","origin":"prose"}]  REFUSED-PROSE

<div title="a > b">no<span></span>de --test no/such/dir</div>
  → occ=[{"command":"node --test no/such/dir","origin":"prose"}]  REFUSED-PROSE
```

Probes P47, P48, P49 all produce `pass=11 fail=1  REFUSED-PROSE`.

### kimi-1 [MAJOR] — the shell builds the words, not the page

**CLOSED.** Cycle 23 added a `shellWordView` that strips `\`, `'`, `"` from the
text before checking for `node` and `--test`, and a `$` / backtick heuristic
that treats any unit naming `node` while carrying those characters as a
candidate on that ground alone. Both tokens are matched case-insensitively.
`SUPPORTED_COMMAND` remains exact, so every candidate is detected and then
refused by the grammar.

Measured with the actual parser against all five forms from the finding, plus
the valid-target arm:

```
`node --\test no/such/dir`
  → occ=[{"c":"node --\\test no/such/dir","orig":"code","supp":false}]  REFUSED

`nod\e --test no/such/dir`
  → occ=[{"c":"nod\\e --test no/such/dir","orig":"code","supp":false}]  REFUSED

`node --te''st no/such/dir`
  → occ=[{"c":"node --te''st no/such/dir","orig":"code","supp":false}]  REFUSED

`node --TEST no/such/dir`
  → occ=[{"c":"node --TEST no/such/dir","orig":"code","supp":false}]  REFUSED

node --t`echo e`st no/such/dir   (in a fenced block)
  → occ=[{"c":"node --t`echo e`st no/such/dir","orig":"code","supp":false}]  REFUSED

`node --\test "skills/parley-tracker/bin/*.test.js"`  (the inverted-purpose arm)
  → occ=[{"c":"node --\\test \"skills/parley-tracker/bin/*.test.js\"","orig":"code","supp":false}]  REFUSED
```

Probes P50–P55 all produce `pass=11 fail=1  REFUSED`.

### kimi-1 [NIT] — a script body nested inside a `<div>` was policed as page text

**CLOSED.** `visibleTextOfHtml` now checks for `<script>`, `<style>`, or
`<template>` opening tags at any position in the HTML, not only at the start of
the block. The scanner finds the matching close tag and skips the entire range.

```
<div>
<script> var c = "node --test no/such/dir"; </script>
</div>
  → occ=[]  (correctly invisible)
```

Probe P56 produces `pass=12 fail=0  GREEN`.

### kimi-1 [NIT] — a comment containing `>` leaked its tail into visible text

**CLOSED.** `visibleTextOfHtml` handles `<!--` by scanning to `-->`, so a
comment containing `>` no longer leaks.

```
<!-- note a > b: node --test no/such/dir -->
  → occ=[]  (correctly invisible)
```

Probe P57 produces `pass=12 fail=0  GREEN`.

### kimi-1 [MINOR] — the occurrence rule is order-dependent

**RECORDED AS FOLLOW-UP.** The prose pass finds `\bnode\b` first, then looks
for `--test\b` in the rest of the line. `` Use `node` with `--test x` `` has
`node` before `--test`, so the rest of the line includes `--test` and the
occurrence is recorded as prose and refused. `` Pass the `--test` flag to
`node` `` has `--test` before `node`; the rest of the line after `node` is
` when you want the runner.`, which has no `--test`, so it is not recorded.

```
Use `node` with `--test no/such/dir` to verify.
  → occ=[{"command":"node with --test no/such/dir to verify.","origin":"prose"}]  REFUSED

Pass the `--test` flag to `node` when you want the runner.
  → occ=[]  GREEN
```

The refused line is not a runnable command (`with` sits between the tokens), so
refusing it exceeds the contract's stated scope ("anything else that renders as
such a command"). I agree with kimi-1's own reasoning: this is a follow-up, not
a blocker. No shipped file trips it (the suite is green at 253/253).

## What I verified and how

### Isolation

```
$ git rev-parse HEAD
f5200f7831fe5f8ad00032d7643800692cc3fbfd

$ git status --short
?? node_modules
```

`node_modules` is the supplied pre-existing link. No tracked file was modified.
All temporary verification scripts were untracked and have been removed.

### Suite

```
$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
```

### Harness (57 probes)

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
P47–P49                      REFUSED-PROSE
P50–P55                      REFUSED
P56–P57                      GREEN
```

Counted by verdict: 30 REFUSED + 16 REFUSED-PROSE = 46 refused, 10 GREEN, 1
RAN-AND-FAILED. Total 57. Every classification matches the author's claim.

### The guard still verifies real commands

P8 (`node --test no/such/dir` in an inline code span) runs through `/bin/sh`,
fails with "Could not find 'no/such/dir'", and the guard reports `pass=11
fail=1 RAN-AND-FAILED`. The four valid-form probes (P4 inline, P7
double-backtick, P19 blockquote, P26 prose-mention) are all `pass=12 fail=0
GREEN`. I also ran the two commands the shipped documents actually publish:

```
$ node --test "skills/parley-tracker/bin/*.test.js"
ℹ tests 35   ℹ pass 35   ℹ fail 0

$ node --test "skills/parley-design-check/test/*.test.js"
ℹ tests 159  ℹ pass 159  ℹ fail 0
```

### False-positive surface — the two shipped documents

I parsed both shipped documents with the actual `commonmark` parser and the
guard's `publishedTestCommands`:

```
skills/parley-tracker/templates/subtask.md:
  → [{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]

skills/parley-design-check/SKILL.md:
  → [{"command":"node --test \"skills/parley-design-check/test/*.test.js\"","origin":"code"}]
```

Both yield only `origin:"code"` occurrences. The `Verify: ` prefix and the
parenthesized checklist item are prose (owner -1); the command is wholly inside
one code span. The occurrence-level rule passes both for exactly the reason the
comment says. The suite passing at 253/253 confirms both are accepted.

### Hunting for a fourth layer

The task asks: three layers have been closed — markdown structure, markdown
rendering, shell word construction. If a fourth exists, name it.

I tested the following candidates by parsing each with the actual `commonmark`
parser, rendering with `HtmlRenderer`, and calling `publishedTestCommands`:

**HTML entities (numeric and named).** CommonMark decodes entities in `text`
nodes. `&#110;ode --test no/such/dir` produces two text nodes (`"n"` and
`"ode --test no/such/dir"`) that join in the visible buffer as `node --test
no/such/dir` and are detected as prose and refused. Entities for the space
(`&#32;`), hyphens (`&#45;&#45;`), and all-character entities
(`&#110;&#111;&#100;&#101;&#32;&#45;&#45;&#116;&#101;&#115;&#116;`) are all
detected and refused. In code spans, entities are NOT decoded (the literal is
`&#110;ode --test x`), so the reader copies the entity text, which does not run
as `node --test`. Correct in both directions.

**Zero-width characters (ZWNJ U+200C, ZWSP U+200B, ZWJ U+200D).** Inserted
between `no` and `de`: `no\u200cde --test no/such/dir`. The parser keeps them in
the text literal, so `\bnode\b` does not match `no\u200cde`. Not detected. I
verified the shell agrees: `/bin/sh -c 'no\u200cde --test no/such/dir'` reports
`command not found` (exit 127). The character is part of the word, so the shell
does not see `node`. Not a false green.

**CDATA and processing instructions.** `<![CDATA[node --test no/such/dir]]>`
and `<?pi node --test no/such/dir?>` are each a single `html_block` node.
`visibleTextOfHtml` skips CDATA content (scans to `]]>`) and PI content (scans
to `?>`). Not detected. A reader does not see or copy this content. Correct.

**Other forms.** Reference links with code text, setext headings, autolinks,
tab-separated tokens, multiple `--test` flags, `$()` command substitution —
all either detected and refused, or detected and valid, or not a runnable
command. No false green.

No fourth layer found. The three layers that have been closed cover the full
path from source text through markdown rendering through shell word
construction to the executed command.

### The approximation-versus-lexing choice

kimi-1 offered two directions: lex each candidate through `/bin/sh`, or use a
cheaper fail-closed approximation. The approximation was chosen because lexing
a line through a shell risks evaluating the very substitutions it is meant to
inspect.

I judged this choice and agree with it, for three reasons:

1. **The approximation is fail-closed in both directions.** `shellWordView`
   removes `\`, `'`, `"` — quoting characters. Removing them can only JOIN
   characters into tokens, never split them. So the word view can only reveal a
   command that the character view missed, never conceal one. The `$` and
   backtick heuristic adds any unit carrying those characters to the candidate
   set. Both can only increase detection, never decrease it.

2. **The grammar already excludes these characters.** `SUPPORTED_COMMAND` is
   `/^node\s+--test\s+[^`;|&<>$\\]+$/`. Every character that `shellWordView`
   strips or that the heuristic flags is already refused by the grammar. So
   detection is the only gap; closing it is sufficient. The refused commands
   are not guesses — they are grammar-excluded on sight.

3. **Lexing through a shell is unsafe and unnecessary.** To get the shell's
   word list, you must run the shell's lexer, which performs quote removal and
   word splitting. But word splitting for `node --t`echo e`st x` requires
   executing `echo e` — the very substitution the guard exists to inspect. You
   cannot get the word list without executing the substitutions. `sh -n`
   (noexec) checks syntax but does not perform word splitting, so it cannot
   tell you what words the shell would build. The approximation avoids this
   entirely: it detects that a command MIGHT be assembled, and refuses it. It
   does not need to know what the command IS.

## Findings

### [NIT] The claim's probe breakdown sums to 56, not 57

The task prompt says: "57 probes behave as intended — 46 refused, 4 valid
forms green, 4 invisible forms ignored, a hard break left alone, and a
genuinely broken path still RUNS and fails." That breakdown sums to 56
(46 + 4 + 4 + 1 + 1). The 57th probe is P57 (comment-tail-leak), which is
GREEN — the comment with a `>` inside is now correctly skipped by
`visibleTextOfHtml`. It is an invisible form, but it is not counted in the
"4 invisible forms" list, which names comment, script body, nested script
body, and image alt (P41, P42, P56, P43). P57 is the regression test for
kimi-1's round-18 NIT and is correctly GREEN. All 57 probes behave as
intended; the breakdown's arithmetic is off by one. No functional impact.

### Signoff: hermes-1 — 2026-07-30
Status: ✅ ACCEPT
