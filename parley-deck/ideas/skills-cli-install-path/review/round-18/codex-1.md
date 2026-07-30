---
idea: skills-cli-install-path
review-round: 18
agent: codex-1
date: 2026-07-30
reviewed-commit: 265eb56
responding-to:
  - codex-1/review/round-17
  - hermes-1/review/round-17
  - kimi-1/review/round-17
---

## Summary

Reviewed `265eb56` in the supplied isolated `wt-codex` worktree. The full suite is
253/253. The unchanged 46-probe harness matches the claim exactly: 37 named refusals, four
green controls, three invisible forms ignored, a hard break left alone, and the genuinely
broken path **run and failed**. Every exact round-17 reproduction is fixed. The two shipped
`Verify:` contexts remain code-origin, and both shipped commands execute real tests
(35/35 and 159/159).

The publication class is nevertheless not closed. Cycle 20 strips raw HTML with
`/<[^>]*>/g`; that expression mistakes a `>` inside a quoted attribute (or comment) for the
end of the tag. Valid invisible markup can therefore leave invisible attribute bytes in the
guard's supposed visible-text view, split `node`, and make a rendered broken command
undetectable. The actual guard stayed green at 12/12 while the page rendered
`node --test no/such/dir` and that copied command exited 1. This is one release-blocking
`[MAJOR]`.

## Round-17 finding dispositions

The cycle log groups the two Hermes NITs together; I retain that six-disposition grouping
here.

| Round-17 finding | Measured disposition at `265eb56` |
| --- | --- |
| `codex-1`: visible command assembled from mixed code/prose or stripped tags belonged to no provenance bucket | **The exact reproductions are fixed.** P32–P35 are all `pass=11 fail=1 REFUSED-PROSE`; the fixture records mixed-span, inline-tag and HTML-block forms as prose occurrences. **The HTML part is not closed as a class:** quoted `>` leaves a bypass, reproduced below. |
| `codex-1`: a `Map` let a valid code occurrence mask identical prose | **Fixed.** P36 is `pass=11 fail=1 REFUSED-PROSE`. The extractor now preserves the `code` and `prose` origins separately, and the fixture asserts `code+prose`. |
| `kimi-1`: comments, formatting tags and GFM strikethrough were invisible splice points | **The exact reproductions are fixed.** P37–P40 are all `pass=11 fail=1 REFUSED-PROSE`. Empty comments/tags are removed and `~~` runs expose the rendered command. **The tag/comment removal remains approximate:** the same finding survives when the invisible construct contains `>`, as below. |
| `kimi-1`: invisible comments/scripts were policed, and image-alt code was executed | **Fixed for all reported forms.** P41–P43 are each `pass=12 fail=0 GREEN`; the extractor fixture also asserts that all three commands are absent. |
| `hermes-1`: a soft break was modeled as a newline instead of a rendered space | **Fixed.** P44 is `pass=11 fail=1 REFUSED-PROSE`; P45, the hard-break control, remains `pass=12 fail=0 GREEN`. |
| `hermes-1`: case-sensitive detection and range-pinned parser | **Fixed.** P46 (`Node --test …`) is detected and refused (`pass=11 fail=1`), while `SUPPORTED_COMMAND` remains lowercase-only. `package.json`, the lockfile, and `npm ls commonmark --depth=0` all resolve exactly `commonmark@0.31.2`. |

## What I verified and how

### Isolation and supplied checks

```text
$ git rev-parse HEAD
265eb56b0bfe9a9634b750605853893f23a705c8

$ git status --short
?? node_modules
```

`node_modules` is the supplied pre-existing link. I wrote no probe file and changed no
tracked implementation file.

```text
$ zsh …/scratchpad/probe-codex.sh
baseline                     pass=12 fail=0
P32 mixed-code-span          pass=11 fail=1  REFUSED-PROSE
P33 html-splits-node         pass=11 fail=1  REFUSED-PROSE
P34 html-splits-flag         pass=11 fail=1  REFUSED-PROSE
P35 html-block               pass=11 fail=1  REFUSED-PROSE
P36 code-masks-prose         pass=11 fail=1  REFUSED-PROSE
P37 html-comment-splice      pass=11 fail=1  REFUSED-PROSE
P38 ins-tag-splice           pass=11 fail=1  REFUSED-PROSE
P39 gfm-strikethrough        pass=11 fail=1  REFUSED-PROSE
P40 gfm-strike-flag          pass=11 fail=1  REFUSED-PROSE
P41 invisible-comment        pass=12 fail=0  GREEN
P42 script-content           pass=12 fail=0  GREEN
P43 image-alt-code           pass=12 fail=0  GREEN
P44 softbreak-splits         pass=11 fail=1  REFUSED-PROSE
P45 softbreak-hardbreak      pass=12 fail=0  GREEN
P46 case-Node                pass=11 fail=1  REFUSED
```

Across all 46 probes: 37 refused, P4/P7/P19/P26 green, P41–P43 ignored, P45 green,
and P8 was `pass=11 fail=1 RAN-AND-FAILED`. The harness's concurrency check was unchanged.

```text
$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
```

### The guard still verifies real commands

The negative execution control reached Node rather than either refusal:

```text
P8 genuinely-broken-path     pass=11 fail=1  RAN-AND-FAILED

$ node --test no/such/dir
Could not find 'no/such/dir'
exit 1
```

The valid inline, double-backtick, and blockquote-container forms are green in P4, P7 and
P19. I also ran the two commands the shipped documents actually publish:

```text
$ node --test "skills/parley-tracker/bin/*.test.js"
ℹ tests 35
ℹ pass 35
ℹ fail 0

$ node --test "skills/parley-design-check/test/*.test.js"
ℹ tests 159
ℹ pass 159
ℹ fail 0
```

### New false-positive surface and the shipped contexts

I passed each context independently through the exact current extractor:

```text
Verify: `node --test "skills/parley-tracker/bin/*.test.js"`
[{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]

- [ ] AC-3 (Verify: `node --test "skills/parley-tracker/bin/*.test.js"`) — COMMIT-SHA
[{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]
```

The `Verify: ` prefix and the parenthesized checklist item therefore do not become prose
occurrences. The fenced command in `skills/parley-design-check/SKILL.md` is also extracted as
`code`. A search over shipped Markdown found no current `~~`, split-tag, script/style, or
uppercase-`Node` form that the broader detector newly refuses. I found no false positive
against anything this project currently needs.

## Findings

### [MAJOR] A quoted `>` defeats the HTML tag stripper and leaves a rendered command invisible

Cycle 20 models visible raw-HTML text with:

```js
node.literal.replace(/<[^>]*>/g, "")
```

`[^>]*` cannot distinguish a tag terminator from `>` inside a quoted attribute. This is valid
CommonMark inline HTML:

```markdown
no<span title="1 > 0"></span>de --test no/such/dir
```

Measured parser and renderer output:

```text
AST:
text:"no"
html_inline:"<span title=\"1 > 0\">"
html_inline:"</span>"
text:"de --test no/such/dir"

HTML:
<p>no<span title="1 > 0"></span>de --test no/such/dir</p>

publishedTestCommands:
[]
```

The `span` has no visible content. Whether a downstream sanitizer preserves or removes it,
the page text is `node --test no/such/dir`. The regex stops at the quoted `>` and leaves
attribute debris between `no` and `de`, so `mentionsATestCommand` sees no `node`.

I injected that Markdown as a virtual `skills/*.md` entry by monkey-patching `fs` in memory
before loading the actual test module; no fixture was written. Actual guard output:

```text
✔ every `node --test` command a shipped file publishes runs tests and passes
ℹ tests 12
ℹ pass 12
ℹ fail 0
```

The copied visible command then failed:

```text
$ node --test no/such/dir
Could not find 'no/such/dir'
exit 1
```

This directly violates “anything else that renders as such a command is refused by name” and
is the same silent false-green shape as the prior blocking findings. It blocks release.

Suggested direction: do not parse complete `html_inline` nodes with a tag regex. CommonMark
already gives the opening and closing markup as separate complete nodes, so their literals can
be omitted while intervening `text` nodes remain. Handle `html_block` with a tokenizer that
respects quoted attributes/comments (or fail closed on raw blocks), and add this exact
quoted-`>` case as a regression.

### Signoff: codex-1 — 2026-07-30
Status: ❌ BLOCK
