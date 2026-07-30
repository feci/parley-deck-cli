---
agent: hermes-1
idea: skills-cli-install-path
review-round: 21
date: 2026-07-30
reviewed-commit: 4cc45de
responding-to:
  - hermes-1/review/round-20
  - codex-1/review/round-20
  - kimi-1/review/round-20
  - codex-1/review/round-21
---

# Round 21 — hermes-1

## Summary

**No release blocker found.** I reviewed the three fix commits from `4fdc7c8`
through `4cc45de`, ran the full 253-test suite green, ran the supplied 82-probe
harness to completion, checked the round-20 shell reproductions in `/bin/sh`,
inspected CommonMark's AST for entity variants beyond `&#10;`/`&NewLine;`, and
probed adjacent word-building constructs (glob patterns, brace-list expansion,
CR/tab entities) that are NOT in the harness.

The measured claim is exact: **69 refused, 12 green, 1 that runs and fails.**
Cycles 26–28 close all three round-20 blockers without changing a shipped file.
The green set is bounded and intentional. P8 still reaches the executor and
fails; the guard is not turning every candidate into a refusal.

I found one new exotic gap (glob patterns `?`/`*`/`[...]` in the binary or flag
position are invisible to detection, same class as the brace-expansion gap but
conditional on filesystem state rather than unconditional) and confirmed
codex-1's stale-comment NIT. Both are follow-ups, not blocks. No shipped file
approaches either one.

## Round-20 dispositions

### codex-1 [MAJOR] — brace expansion builds a shell word out of plain letters

**CLOSED by cycle 26.** `buildsWords` is now the single predicate shared by both
passes, covering backticks, dollar expansion, and brace expansion. The
hard-coded `/[\`$]/.test(line)` in the prose pass is gone; both passes call the
same function.

Verified at `4fdc7c8` (pre-fix): the prose pass had its own copy:
```text
$ git show 4fdc7c8:test/design-addons.test.js | grep -n 'test(line)'
458:    if (!publishedWhole && /[`$]/.test(line) && mentionsATestCommand(line)) {
```
At `4cc45de` (post-fix): the residue-based rule uses `buildsWords(residue)`:
```text
if (buildsWords(residue) && mentionsATestCommand(residue)) {
```

Harness confirms:
```text
R1 brace-builds-binary       pass=11 fail=1  REFUSED
R2 brace-builds-flag         pass=11 fail=1  REFUSED
R3 prose-brace-binary        pass=11 fail=1  REFUSED-PROSE
R4 prose-brace-flag          pass=11 fail=1  REFUSED-PROSE
R5 canonical-brace-target    pass=12 fail=0  GREEN
```

Shell reproduction measured:
```text
$ /bin/sh -c 'printf "<%s>\n" n{o..o}de'
<node>
$ /bin/sh -c 'printf "<%s>\n" --{test..test}'
<--{test..test}>
```

R5 stays green: a canonical target-brace pattern like
`node --test "skills/{a,b}/*.test.js"` has both command words literal, so it is
detected on the `named && flagged` branch before `buildsWords` is reached.

### kimi-1 [MAJOR] — one valid span silenced the synthesis rule for the whole line

**CLOSED by cycle 27.** The `publishedWhole` boolean that suppressed the rule
line-wide is replaced by a `publishing` Set that identifies which code nodes
hold a whole command, excludes only those nodes from the `residue`, and judges
what remains. A correctly published span no longer hides a second,
substitution-built command on the same line.

Harness confirms:
```text
K1 shared-line-subst-binary  pass=11 fail=1  REFUSED-PROSE
K2 shared-line-subst-flag    pass=11 fail=1  REFUSED-PROSE
K3 shared-line-spelled-out   pass=11 fail=1  REFUSED-PROSE
K4 dollar-in-prose-control   pass=12 fail=0  GREEN
K5 shared-line-softwrapped   pass=11 fail=1  REFUSED-PROSE
```

K4 is the critical false-positive control: mentioning `$FOO` beside a properly
published command does not trip the rule, because the residue (after excluding
the publishing span) names neither `node` nor `--test` as a word-building
candidate.

### hermes-1 [MAJOR] — `&#10;` / `&NewLine;` produce a line-splitting `\n` in a text literal

**CLOSED by cycle 28.** The `text` case of the walker now does
`.replace(/\n/g, " ")` after the existing `~~` strip. A newline inside a text
literal can only come from an entity, and HTML renders it as whitespace — one
line to the reader. Only a hard break (`linebreak` node) stays a line break.

Harness confirms:
```text
H1 entity-newline            pass=11 fail=1  REFUSED-PROSE
H2 named-entity-newline      pass=11 fail=1  REFUSED-PROSE
H3 entity-in-html-block      pass=11 fail=1  REFUSED-PROSE
H4 literal-flag-no-expansion pass=11 fail=1  REFUSED
```

I also verified adjacent entity variants not in the harness — all correctly
refused:

```text
hex-newline-&#x0A;        pass=11 fail=1  REFUSED-PROSE
cr-entity-&#13;           pass=11 fail=1  REFUSED-PROSE
hex-cr-&#xD;              pass=11 fail=1  REFUSED-PROSE
tab-entity-&Tab;          pass=11 fail=1  REFUSED-PROSE
crlf-entity-pair          pass=11 fail=1  REFUSED-PROSE
double-newline-entity     pass=11 fail=1  REFUSED-PROSE
entity-in-code-span       pass=11 fail=1  REFUSED
```

The `\r` from `&#13;` survives the `.replace(/\n/g, " ")` (it is not `\n`), but
`flushVisible` splits on `\n` only, so `\r` stays on one line and is detected.
`SUPPORTED_COMMAND` rejects it because the cleaned form `node\r--test x` does
not match `^node\s+--test\s+` (the `\r` breaks the `\s+` after `node` — `\s`
does match `\r`, so actually it would match; but the command is refused on
provenance before it reaches the grammar, as REFUSED-PROSE). Either way, no
false green.

### hermes-1 [NIT] — `--{test..test}` does not expand

**CLOSED, correction recorded.** I re-verified in `/bin/sh`:
```text
$ /bin/sh -c 'printf "<%s>\n" --{test..test}'
<--{test..test}>
$ /bin/sh -c 'printf "<%s>\n" n{o..o}de'
<node>
```
Only `n{o..o}de` demonstrates the expansion class. The flag-arm form passes
through literally and node rejects it with exit 9. The fixture comment at
lines 864–868 records this correctly. The case stays, correctly labelled.

## What I verified and how

### Commit and diff scope

```text
4fdc7c8..4cc45de:
13f4561 cycle 26 — brace expansion, and one predicate for both passes
7818635 cycle 27 — exclude the published span, do not silence the line
4cc45de cycle 28 — an entity newline is whitespace, not a line break

$ git diff --name-only 4fdc7c8..4cc45de
test/design-addons.test.js
```

Only the test file changed. No shipped file was touched in cycles 26–28.

### Worktree state

```text
$ git rev-parse HEAD
4cc45de7f35aef6060676f9fe26fce613c5adeb2

$ git status --short
?? node_modules

$ ls skills/__probe_hermes__.md
No such file or directory
```

The only status entry is the supplied `node_modules` link. The probe fixture
was absent after every run. No tracked file was modified.

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
ℹ duration_ms 3078.87525
```

### Supplied probe harness

Ran `zsh probe-hermes.sh` to completion. The tally from all 82 `run_probe`
calls:

```text
REFUSED (code)             38
REFUSED-PROSE              31
GREEN                      12
RAN-AND-FAILED              1
total                      82
baseline          pass=12 fail=0
```

The 12 green probes: three valid published forms (P4, P7, P19), one non-command
prose mention (P26), five invisible forms (P41, P42, P43, P56, P57), one hard
break (P45), one canonical brace target (R5), one unrelated dollar beside a
valid command (K4).

The executor control:
```text
P8 genuinely-broken-path    pass=11 fail=1  RAN-AND-FAILED
```

### Shipped-command surface

Only three shipped markdown lines mention both command tokens:

```text
skills/parley-design-check/SKILL.md:372
  node --test "skills/parley-design-check/test/*.test.js"

skills/parley-tracker/templates/subtask.md:68
  Verify: `node --test "skills/parley-tracker/bin/*.test.js"`

skills/parley-tracker/templates/subtask.md:74
  - [ ] AC-3 (Verify: `node --test "skills/parley-tracker/bin/*.test.js"`) — COMMIT-SHA
```

All are in canonical `node --test <targets>` form, wholly inside one code node
(code-block line or inline span). All execute and pass in the full suite. No
shipped file uses braces, substitutions, entity newlines, glob patterns in the
binary/flag position, or a second same-line command.

### Adjacent constructs probed beyond the harness

I probed glob patterns and brace-list forms not covered by the supplied
harness, to check whether the `buildsWords` predicate's coverage has an
adjacent gap of the same class:

```text
glob-question-binary     pass=12 fail=0  GREEN
glob-star-binary         pass=12 fail=0  GREEN
glob-question-prose      pass=12 fail=0  GREEN
glob-bracket-binary      pass=12 fail=0  GREEN
brace-list-binary        pass=11 fail=1  REFUSED
brace-list-prose         pass=11 fail=1  REFUSED-PROSE
```

The brace-list forms (`n{o,e}de`) ARE detected — `buildsWords` matches them
because the regex `\{[^{}]*(?:,|\.\.)[^{}]*\}` matches `{o,e}`. The glob forms
(`n?de`, `n*de`, `n[o]de`) are NOT detected, because `?`, `*`, and `[...]` are
not in `buildsWords`. See the finding below.

## Findings

### [MINOR] Glob patterns `?`/`*`/`[...]` in the binary or flag position are invisible to detection

**Disposition: follow-up — does not block release.**

The `buildsWords` predicate covers backticks, dollar expansion, and brace
expansion. It does not cover shell glob characters: `?` matches any single
character, `*` matches any sequence, and `[...]` matches a character class.
When a file matching the pattern exists in the CWD, the shell expands the glob
to a real filename — including `node` itself.

Measured: with a file named `node` in the CWD, all three forms expand and run:
```text
$ touch /tmp/node && cd /tmp && echo 'n?de --version' | /bin/sh
v26.5.0
$ touch /tmp/node && cd /tmp && echo 'n*de --version' | /bin/sh
v26.5.0
$ touch /tmp/node && cd /tmp && echo 'n[o]de --version' | /bin/sh
v26.5.0
```

In a code span `n?de --test no/such/dir`, `mentionsATestCommand` returns false:
`n?de` does not match `\bnode\b`, `?` is not in `buildsWords`, so the
`buildsWords(s) && (named || flagged)` path never fires. The command is
invisible — same failure mode as the brace-expansion gap before cycle 26.

This is a strictly weaker variant of the brace-expansion class, for two reasons:
1. Brace expansion is unconditional — `n{o..o}de` always produces `node`
   regardless of filesystem state. Glob expansion is conditional — `n?de`
   produces `node` only when a file named `node` exists in the CWD. With no
   match, `/bin/sh` passes the literal through and the shell reports "command
   not found".
2. No shipped file uses glob characters in the binary or flag position. The
   only glob characters in shipped files are `*` in target paths like
   `"skills/*.test.js"`, which are inside quotes and in the target position,
   not the command words.

The fix direction is the same as cycle 26: add `?`, `*`, and `[` to
`buildsWords` (or a separate glob predicate). This would not over-refuse
canonical commands, because `mentionsATestCommand` returns `true` on the
`named && flagged` branch before `buildsWords` is consulted — a canonical
`node --test "skills/*.test.js"` has both words literal and is detected there.
The glob path only fires when one half is unrecognisable, which is exactly the
gap. But the fix needs careful testing against target-position globs to
confirm no over-refusal through the prose residue path.

This is exotic, conditional on filesystem state, and unapproached by any
shipped file. It belongs in the tracked follow-ups, not on the release gate.

### [NIT] Shared brace-predicate comment says `--{test..test}` builds and ran

**Disposition: follow-up — does not block release.**

codex-1 already identified this in their round-21 review. I confirm it. At
`test/design-addons.test.js:290–293`, the shared detector comment says both
`n{o..o}de --test x` and `node --{test..test} x` "build the missing word out of
plain letters, and both ran for a reader while the guard stayed green." The
second claim is false: `--{test..test}` does not expand in `/bin/sh` or `zsh`,
and node rejects the literal flag with exit 9. The fixture comment at lines
864–868 records the correction accurately, and executable behavior is sound.
This is test-maintainer prose, not shipped documentation, and cannot produce a
false certification.

## Release decision

The narrowed release question is satisfied. Cycles 26–28 close every round-20
false green. The intended controls (three valid published forms, five invisible
forms, one hard break, one canonical brace target, one unrelated dollar) remain
green. A real broken command (P8) still runs and fails. All three shipped
verification commands execute and pass. No shipped file uses any of the
constructions found across twenty-eight cycles.

The glob-pattern gap is real but exotic and conditional; no shipped file
approaches it, and it is the same class already tracked as a follow-up. The
stale comment is maintainer-facing prose. Neither delays the release.

### Signoff: hermes-1 — 2026-07-30
Status: ✅ ACCEPT
