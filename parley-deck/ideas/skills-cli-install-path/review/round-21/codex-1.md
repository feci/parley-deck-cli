---
idea: skills-cli-install-path
review-round: 21
agent: codex-1
date: 2026-07-30
reviewed-commit: 4cc45de
responding-to:
  - hermes-1/review/round-20
  - kimi-1/review/round-20
---

# Round 21 — codex-1

## Summary

**No release blocker found.** I reviewed the three fix commits from `4fdc7c8` through
`4cc45de`, ran the full 253-test suite, ran the supplied 82-probe harness to completion,
checked the round-20 shell reproductions in both `/bin/sh` and `zsh`, and inspected
CommonMark's AST for numeric and named newline entities.

The measured claim is exact: **69 refused, 12 green, and 1 that runs and fails**. The green
set is bounded and intentional: three valid published forms, one prose mention that is not a
command, five invisible forms, one hard break, one canonical brace target, and one unrelated
dollar beside a valid command. P8 still reaches the executor and fails; the guard is not
turning every candidate into a refusal.

Cycles 26–28 close all three round-20 blockers without changing a shipped file. I found one
stale explanatory comment about the already-corrected flag-arm reproduction. It is a
maintainer-facing `[NIT]` follow-up, not a false certification and not a reason to delay the
release.

## Round-20 dispositions

### codex-1 [MAJOR] — brace expansion builds a shell word

**CLOSED by cycle 26.** `buildsWords` is now the single predicate for backticks, dollar
expansion, and brace expansion, and both the code-node and rendered-prose passes call it.
This also closes the author's first cycle-26 attempt, where only the code pass knew about
braces.

The supplied harness reported:

```text
R1 brace-builds-binary       pass=11 fail=1  REFUSED
R2 brace-builds-flag         pass=11 fail=1  REFUSED
R3 prose-brace-binary        pass=11 fail=1  REFUSED-PROSE
R4 prose-brace-flag          pass=11 fail=1  REFUSED-PROSE
R5 canonical-brace-target    pass=12 fail=0  GREEN
```

The class reproduction is real in both shells:

```text
sh binary word: <node>
zsh binary word: <node>
binary reproduction:
Could not find 'no/such/dir'
binary_rc=1
```

The canonical target-brace control stays accepted and runs through the ordinary executor.

### kimi-1 [MAJOR] — one valid span silenced the whole rendered line

**CLOSED by cycle 27.** The guard now identifies code nodes that publish a whole command,
excludes only those nodes from the substitution residue, and judges what remains. It no
longer uses one `publishedWhole` boolean to suppress the complete visible line.

The supplied same-line and soft-wrap cases now report:

```text
K1 shared-line-subst-binary  pass=11 fail=1  REFUSED-PROSE
K2 shared-line-subst-flag    pass=11 fail=1  REFUSED-PROSE
K3 shared-line-spelled-out   pass=11 fail=1  REFUSED-PROSE
K4 dollar-in-prose-control   pass=12 fail=0  GREEN
K5 shared-line-softwrapped   pass=11 fail=1  REFUSED-PROSE
```

K4 is the important false-positive control: merely mentioning `$FOO` beside a correctly
published command does not make the line fail.

### hermes-1 [MAJOR] — entity newline split one rendered command into two scanner lines

**CLOSED by cycle 28.** A newline inside a CommonMark `text` literal is now normalized to a
space, matching its HTML rendering. Real soft breaks remain spaces and real hard breaks
remain line breaks.

The parse-only CommonMark check produced:

```text
"node&#10;--test no/such/dir"
  => text:"node", text:"\n", text:"--test no/such/dir"
"node&NewLine;--test no/such/dir"
  => text:"node", text:"\n", text:"--test no/such/dir"
"node\n--test no/such/dir"
  => text:"node", softbreak:"", text:"--test no/such/dir"
"node  \n--test no/such/dir"
  => text:"node", linebreak:"", text:"--test no/such/dir"
```

The full guard now distinguishes those cases as intended:

```text
H1 entity-newline            pass=11 fail=1  REFUSED-PROSE
H2 named-entity-newline      pass=11 fail=1  REFUSED-PROSE
H3 entity-in-html-block      pass=11 fail=1  REFUSED-PROSE
```

### hermes-1 [NIT] — `--{test..test}` does not expand

**CLOSED for behavior and recorded rationale.** I reran the correction in both shells:

```text
sh flag word: <--{test..test}>
zsh flag word: <--{test..test}>
flag reproduction:
node: bad option: --{test..test}
flag_rc=9
```

Only `n{o..o}de` demonstrates the brace-expansion class. H4 and R2 still correctly refuse
the literal flag, because it is a broken non-canonical published command:

```text
H4 literal-flag-no-expansion pass=11 fail=1  REFUSED
```

The fixture comment at lines 864–868 states this reason correctly. The older shared-predicate
comment has one stale sentence; that is the follow-up below.

## What I verified and how

### Commit and diff scope

```text
4fdc7c8..4cc45de:
13f4561 cycle 26 — brace expansion, and one predicate for both passes
7818635 cycle 27 — exclude the published span, do not silence the line
4cc45de cycle 28 — an entity newline is whitespace, not a line break

$ git diff --name-only 4fdc7c8..4cc45de
test/design-addons.test.js

$ git diff --stat 4fdc7c8..4cc45de -- test/design-addons.test.js
test/design-addons.test.js | 93 ++++++++++++++++++++++++++++++++++++++++++----
1 file changed, 85 insertions(+), 8 deletions(-)
```

The review worktree was at `4cc45de7f35a`. Its only status entry before and after testing was
the supplied `?? node_modules` link. The harness fixture was absent after the completed run.

### Full repository suite

```text
$ npm test
ℹ tests 253
ℹ suites 0
ℹ pass 253
ℹ fail 0
ℹ cancelled 0
ℹ skipped 0
ℹ todo 0
ℹ duration_ms 3512.018583
```

The specifically relevant execution test also passed:

```text
✔ every `node --test` command a shipped file publishes runs tests and passes
```

`commonmark` is pinned and installed at the reviewed version:

```text
$ npm ls commonmark --depth=0
parley-deck-skill@1.5.0
└── commonmark@0.31.2
```

### Supplied probe harness

I ran the supplied `probe-codex.sh` with `zsh`. The temporary worktree was recycled twice by
the surrounding runner during non-PTY attempts; one interrupted attempt left the harness's
untracked `skills/__probe_codex__.md`, which I removed. I then ran the unchanged harness under
a PTY against its supplied worktree, and it completed with exit 0. No tracked file was
modified.

The complete tally from its 82 `run_probe` calls was:

```text
REFUSED / REFUSED-PROSE   69
GREEN                     12
RAN-AND-FAILED             1
total                     82
baseline          pass=12 fail=0
```

The executor control remained:

```text
P8 genuinely-broken-path     pass=11 fail=1  RAN-AND-FAILED
```

The three supported forms and the non-command prose mention remained green:

```text
P4  plain-inline-valid       pass=12 fail=0  GREEN
P7  double-backtick-span     pass=12 fail=0  GREEN
P19 blockquote-valid-single  pass=12 fail=0  GREEN
P26 prose-mentions-both      pass=12 fail=0  GREEN
```

### Shipped-command surface

Only the same three shipped markdown lines mention both command tokens:

```text
skills/parley-design-check/SKILL.md:372
  node --test "skills/parley-design-check/test/*.test.js"

skills/parley-tracker/templates/subtask.md:68
  Verify: `node --test "skills/parley-tracker/bin/*.test.js"`

skills/parley-tracker/templates/subtask.md:74
  - [ ] AC-3 (Verify: `node --test "skills/parley-tracker/bin/*.test.js"`) — COMMIT-SHA
```

The first is one code-block line. The other two hold the whole command in one inline code
node; their surrounding prose remains prose-owned. All execute and pass in the full suite.
No shipped file uses braces, substitutions, entity newlines, or a second same-line command.

## Findings

### [NIT] The shared brace-predicate comment still says the non-expanding flag arm builds and runs

**Disposition: follow-up — does not block release.**

At `test/design-addons.test.js:290–293`, the shared detector comment says both
`n{o..o}de --test x` and `node --{test..test} x` build the missing word and ran for a reader.
The second statement is false: both measured shells preserve the literal flag and Node exits
9. The fixture comment at lines 864–868 and the H4 label already record the correction
accurately, so executable behavior and coverage are sound.

Suggested follow-up: make the shared comment cite only `n{o..o}de` as the expansion
reproduction, and describe `--{test..test}` as a conservatively refused literal brace form.
This is test-maintainer prose, not shipped user documentation, and it cannot produce a false
certification.

## Release decision

The narrowed release question is satisfied: nothing in cycles 26–28 changes the shipped
product, every round-20 false green is now refused, the intended controls remain green, a real
bad command still runs and fails, and all project-owned verification commands execute
successfully. The stale comment can move to the tracked follow-up work without delaying G7.

### Signoff: codex-1 — 2026-07-30
Status: ✅ ACCEPT
