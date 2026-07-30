---
idea: skills-cli-install-path
review-round: 20
agent: codex-1
date: 2026-07-30
reviewed-commit: 4fdc7c8
responding-to:
  - hermes-1/review/round-19
  - kimi-1/review/round-19
---

## Summary

Reviewed `4fdc7c8` in the supplied isolated `wt-codex` worktree. The full suite is
253/253. The unmodified supplied harness has 68 probes plus its baseline invocation and
reproduces the corrected claim exactly: 57 refused, 10 green, and one genuinely broken
path that runs and fails. Cycles 24 and 25 close every round-19 disposition without
regressing the three supported forms, the non-command prose mention, the five invisible
forms, the hard break, or the two shipped documents.

I found one remaining false green. Shell-word construction is not closed: brace expansion
can construct either `node` or `--test` without any backslash, quote, `$`, or backtick.
The detector therefore sees neither token pair and skips the command in both the code and
prose passes. Four additions to the supplied `run_probe` list all report GREEN, although
the current `/bin/sh` and `zsh` expand `n{o..o}de` to exactly `node`. This is the same
release-blocking class as cycle 24, not a fourth layer and not a renewed pass asymmetry.

## Round-19 dispositions

### hermes-1 [NIT] — the probe breakdown was arithmetically wrong

**CLOSED.** The supplied script contains 68 `run_probe` calls. I counted the actual
results:

```text
REFUSED / REFUSED-PROSE   57
GREEN                     10
RAN-AND-FAILED             1
total                     68
baseline          pass=12 fail=0
```

The ten green probes are exactly three supported published forms (P4, P7, P19), one
non-command prose mention (P26), five invisible forms (P41, P42, P43, P56, P57), and one
hard break (P45). The implementation record and round-20 brief now state that breakdown.

### codex-1 measured finding — substitution can build the binary name

**CLOSED by cycle 24.** The change from `f5200f7` to `4ac913e` makes a unit containing
`$` or a backtick a candidate when either the binary or flag is recognisable; it no
longer requires `node` to be spelled out. The supplied probes reproduce the closure:

```text
P58 subst-builds-binary      pass=11 fail=1  REFUSED
P59 expansion-builds-binary  pass=11 fail=1  REFUSED
```

The grammar already excludes `$` and backticks, so these newly reached candidates are
refused rather than executed.

### kimi-1 [MAJOR] — the prose pass still read raw characters

**CLOSED by cycle 25.** The rendered-page pass now builds the same quote/backslash-free
word view as the code pass, maps its indices back to the raw line for provenance, and
matches both tokens case-insensitively. All prose arms are now refused:

```text
N1 prose-escaped-flag        pass=11 fail=1  REFUSED-PROSE
N2 prose-escaped-binary      pass=11 fail=1  REFUSED-PROSE
N3 prose-quote-splice        pass=11 fail=1  REFUSED-PROSE
N4 prose-uppercase-flag      pass=11 fail=1  REFUSED-PROSE
N5 prose-escaped-valid       pass=11 fail=1  REFUSED-PROSE
N8 prose-subst-builds-node   pass=11 fail=1  REFUSED-PROSE
```

The stated partial disposition is also correct. Cycle 24 had already changed the shared
`mentionsATestCommand` function used by the code-node pass; cycle 25 does not touch that
pass. At `4fdc7c8`, the three span arms are:

```text
N6 span-subst-builds-node    pass=11 fail=1  REFUSED
N7 span-subst-inside-node    pass=11 fail=1  REFUSED
N9 span-subst-valid-target   pass=11 fail=1  REFUSED
```

Thus only the prose arm remained open after cycle 24, and only that arm changed in cycle
25.

## What I verified and how

### Worktree and suite

```text
$ git rev-parse HEAD
4fdc7c8f6ac876a8663b483b9f140f9ffb2c8c2b

$ git status --short
?? node_modules

$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
```

`node_modules` is the supplied pre-existing link. The probe fixture was absent after every
run. I did not modify a tracked implementation file.

### Supplied harness

The unmodified `zsh probe-codex.sh` run reported:

```text
baseline                     pass=12 fail=0
P8 genuinely-broken-path     pass=11 fail=1  RAN-AND-FAILED
P58 subst-builds-binary      pass=11 fail=1  REFUSED
P59 expansion-builds-binary  pass=11 fail=1  REFUSED
N1–N5                        pass=11 fail=1  REFUSED-PROSE
N6–N7                        pass=11 fail=1  REFUSED
N8                           pass=11 fail=1  REFUSED-PROSE
N9                           pass=11 fail=1  REFUSED
```

The other 57 round-14 through round-18 probes retain their intended classifications. In
particular, P8 proves the guard still executes a supported-looking command and observes
its failure; it is not merely a text linter.

### Supported forms and false-positive surface

The supported-form controls remain green:

```text
P4  plain-inline-valid       pass=12 fail=0  GREEN
P7  double-backtick-span     pass=12 fail=0  GREEN
P19 blockquote-valid-single  pass=12 fail=0  GREEN
P26 prose-mentions-both      pass=12 fail=0  GREEN
```

The shipped occurrences are still:

```text
skills/parley-design-check/SKILL.md:372
  node --test "skills/parley-design-check/test/*.test.js"

skills/parley-tracker/templates/subtask.md:68
  Verify: `node --test "skills/parley-tracker/bin/*.test.js"`

skills/parley-tracker/templates/subtask.md:74
  - [ ] AC-3 (Verify: `node --test "skills/parley-tracker/bin/*.test.js"`) — COMMIT-SHA
```

The first is one code-block line. In the latter two, both tokens are wholly owned by the
same inline code node; `Verify: ` and the checklist parentheses remain prose-owned. The
passing shipped-command test executes the two distinct commands, requires exit zero and
requires at least one passing test. No project-owned command is newly refused, and the
known order-dependent prose over-refusal remains a follow-up rather than a release issue.

## Findings

### [MAJOR] Brace expansion constructs either shell word and is invisible in both passes

The shared detector removes only `\`, `'`, and `"` and treats only `$` and backticks as
generic word-building constructs:

```js
const shellWordView = (s) => s.replace(/[\\'"]/g, "");
// ...
return /[`$]/.test(s) && (named || flagged);
```

Brace expansion needs none of those characters. A one-element sequence constructs the
exact missing word:

```text
n{o..o}de --test no/such/dir
node --{test..test} no/such/dir
```

For the first line, `flagged` is true but `named` and the `$`/backtick fallback are false.
For the second, `named` is true but `flagged` and the fallback are false. Neither line
becomes a candidate.

I appended these two fenced forms and their prose equivalents to a copy of the supplied
harness, changing no harness logic:

```text
R20 brace-builds-binary      pass=12 fail=0  GREEN
R20 brace-builds-flag        pass=12 fail=0  GREEN
R20 prose-brace-binary       pass=12 fail=0  GREEN
R20 prose-brace-flag         pass=12 fail=0  GREEN
```

The shell behavior is measured rather than inferred:

```text
$ /bin/sh -c 'printf "<%s>\n" n{o..o}de'
<node>
$ zsh -c 'printf "<%s>\n" n{o..o}de'
<node>
$ /bin/sh -c 'n{o..o}de --version'
v26.5.0
```

Therefore the first reproduction expands to the same command and bad target that P8 proves
runs and fails, while the guard reports success. The flag arm is the symmetric construction
of the other recognised word. Prose and code both miss it, so cycle 25's symmetry is sound;
the shared approximation is incomplete.

This **blocks release** under the rule used for cycles 23 and 24: a documented verification
command can fail for the reader while the test named for published commands reports green.
It also directly narrows cycle 24's universal rationale that a word-building construct can
produce either word.

Fix direction: extend the fail-closed word-building predicate to actual brace-expansion
forms when either half is recognisable, and cover binary/flag construction in both code and
prose. Preserve ordinary canonical target brace patterns such as
`node --test "test/{unit,integration}/*.test.js"`: both command tokens are already literal,
and the executor should continue deciding whether those targets work.

### Signoff: codex-1 — 2026-07-30
Status: ❌ BLOCK
