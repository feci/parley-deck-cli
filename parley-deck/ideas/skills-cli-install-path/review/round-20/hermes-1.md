---
agent: hermes-1
idea: skills-cli-install-path
review-round: 20
date: 2026-07-30
reviewed-commit: 4fdc7c8
responding-to:
  - hermes-1/review/round-19
  - codex-1/review/round-19
  - kimi-1/review/round-19
---

## Summary

Reviewed `4fdc7c8` in the supplied `wt-hermes` worktree. The full suite is 253/253
green. The unmodified supplied harness reproduces the corrected claim exactly: 68
probes, 57 refused, 10 green (3 valid published forms, 1 prose mention, 5 invisible
forms, 1 hard break), and 1 that runs and fails (P8). Cycles 24 and 25 close every
round-19 disposition. The two shipped documents that write `Verify: ` before a span
and the checklist item with the span in parentheses both pass — their commands are
wholly inside one inline code node with `code` provenance and `SUPPORTED_COMMAND`
form.

I found two remaining false greens. codex-1 independently found the first (brace
expansion) and is blocking on it; I have verified their finding and agree it blocks
release. I also found a second, distinct gap: HTML character references that resolve
to `\n` (`&#10;`, `&NewLine;`) produce a text node containing a literal newline,
which `flushVisible` splits on — making the command invisible in the prose pass.
This is a fourth layer (entity resolution producing line-splitting characters) that
the softbreak handling does not cover, because commonmark emits a `text` node with
`\n` for these entities, not a `softbreak` node.

## Round-19 dispositions

### hermes-1 [NIT] — the probe breakdown was arithmetically wrong

**CLOSED.** The supplied script contains 68 `run_probe` calls (P1–P59 plus N1–N9;
the 69th `run_probe` match is the function definition). I counted the actual
results from my run:

```text
REFUSED / REFUSED-PROSE   57
GREEN                     10
RAN-AND-FAILED             1
total                     68
baseline          pass=12 fail=0
```

The ten green probes are exactly: three valid published forms (P4, P7, P19), one
non-command prose mention (P26), five invisible forms (P41, P42, P43, P56, P57),
and one hard break (P45). The implementation record and round-20 brief now state
this breakdown. Corrected.

### codex-1 measured finding — substitution can build the binary name

**CLOSED by cycle 24.** The diff from `f5200f7` to `4ac913e` changes
`mentionsATestCommand` so that a unit carrying `$` or a backtick is a candidate
when either the binary or the flag is recognisable — it no longer requires `node`
to be spelled out. The grammar already excludes `$` and backticks, so the newly
reached candidates are refused rather than executed. Reproduced:

```text
P58 subst-builds-binary      pass=11 fail=1  REFUSED
P59 expansion-builds-binary  pass=11 fail=1  REFUSED
```

### kimi-1 [MAJOR] — the prose pass still read raw characters

**CLOSED by cycle 25.** The diff from `4ac913e` to `4fdc7c8` rewrites the prose
arm of `flushVisible` to build the same quote/backslash-free word view as the code
pass, with an index map (`toRaw`) back into the raw line for provenance. Both
tokens are matched case-insensitively in the word view. All prose arms are now
refused:

```text
N1 prose-escaped-flag        pass=11 fail=1  REFUSED-PROSE
N2 prose-escaped-binary      pass=11 fail=1  REFUSED-PROSE
N3 prose-quote-splice        pass=11 fail=1  REFUSED-PROSE
N4 prose-uppercase-flag      pass=11 fail=1  REFUSED-PROSE
N5 prose-escaped-valid       pass=11 fail=1  REFUSED-PROSE
N8 prose-subst-builds-node   pass=11 fail=1  REFUSED-PROSE
```

The stated partial disposition is also correct. Cycle 24 changed the shared
`mentionsATestCommand` function used by the code-node pass; cycle 25 does not
touch that pass. The three span arms retain their cycle-24 behavior:

```text
N6 span-subst-builds-node    pass=11 fail=1  REFUSED
N7 span-subst-inside-node    pass=11 fail=1  REFUSED
N9 span-subst-valid-target   pass=11 fail=1  REFUSED
```

Only the prose arm remained open after cycle 24, and only that arm changed in
cycle 25. Verified.

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

`node_modules` is the supplied pre-existing link. I removed three empty untracked
files (`guard`, `one`, `softbreak`) created as side effects of probe runs that
resolved to file targets; no tracked file was modified. The probe fixture was
absent after every run.

### Supplied harness

The unmodified `zsh probe-hermes.sh` run reported the full 68-probe battery. The
results match the claim exactly: 57 refused, 10 green, 1 ran-and-failed. P8
proves the guard still executes a supported-looking command and observes its
failure; it is not merely a text linter.

### Supported forms and false-positive surface

The supported-form controls remain green:

```text
P4  plain-inline-valid       pass=12 fail=0  GREEN
P7  double-backtick-span     pass=12 fail=0  GREEN
P19 blockquote-valid-single  pass=12 fail=0  GREEN
P26 prose-mentions-both      pass=12 fail=0  GREEN
```

The shipped occurrences are:

```text
skills/parley-design-check/SKILL.md:372
  node --test "skills/parley-design-check/test/*.test.js"

skills/parley-tracker/templates/subtask.md:68
  Verify: `node --test "skills/parley-tracker/bin/*.test.js"`

skills/parley-tracker/templates/subtask.md:74
  - [ ] AC-3 (Verify: `node --test "skills/parley-tracker/bin/*.test.js"`) — COMMIT-SHA
```

I verified the latter two through the full `publishedTestCommands` function. Both
produce a single occurrence with `origin: "code"` and `SUPPORTED_COMMAND: true`.
The `Verify: ` prefix and the checklist parentheses remain prose-owned; the
command is wholly inside one inline code node, so the occurrence-level provenance
rule admits them. The first is one code-block line, also `code` origin and
supported form. No project-owned command is newly refused.

### Parser-level verification of the `&#10;` entity gap

I used small `node -e` one-liners that only parse strings (per the method
guidance) to inspect what commonmark 0.31.2 produces for `&#10;` in prose:

```text
$ node -e 'const {Parser}=require("commonmark"); let p=new Parser().parse("node&#10;--test no/such/dir"); let w=p.walker(); let ev; while((ev=w.next())) if(ev.entering) console.log(ev.node.type, JSON.stringify(ev.node.literal||""));'

text "node"
text "\n"
text "--test no/such/dir"
```

For comparison, a real source newline:

```text
$ node -e '... parse("node\n--test no/such/dir") ...'

text "node"
softbreak ""
text "--test no/such/dir"
```

The guard handles `softbreak` by emitting a space (`emit(" ", -1)`) — one line,
both tokens, detected as prose. But `&#10;` resolves to a `text` node with a
literal `\n`, which `flushVisible` splits on — two lines, neither with both
tokens, invisible. `&NewLine;` produces the same `text "\n"` structure. Both
render as whitespace to the reader (HTML collapses `\n` to a space), so the
reader copies one runnable command; the guard sees two empty halves.

I confirmed the invisibility through the full `publishedTestCommands` function:

```text
publishedTestCommands("node&#10;--test no/such/dir")          => []
publishedTestCommands("Run node&#10;--test no/such/dir now.") => []
publishedTestCommands("node&NewLine;--test no/such/dir")      => []
publishedTestCommands("node\n--test no/such/dir")             => [{command:"node --test no/such/dir", origin:"prose"}]
```

And through the probe harness with the guard's full test suite:

```text
ENTITY-LF-BROKEN:   pass=12 fail=0  GREEN (invisible)
ENTITY-PROSE-VALID: pass=12 fail=0  GREEN (false green)
```

For contrast, `&#13;` (CR) and `&#9;` (tab) entities do NOT split the line — `\r`
and `\t` are not `\n`, so `flushVisible` keeps them on one line, and both are
detected and refused:

```text
ENTITY-CR-PROSE:    pass=11 fail=1  REFUSED-PROSE
ENTITY-TAB-PROSE:   pass=11 fail=1  REFUSED-PROSE
```

In a fenced code block, entities are NOT resolved — the literal stays
`node&#10;--test no/such/dir`, which `mentionsATestCommand` sees (because
`\bnode\b` matches before `&` and `--test\b` matches after `;`), and
`SUPPORTED_COMMAND` refuses (because `&` is in the exclusion set). So the code
pass is closed; only the prose pass is open:

```text
ENTITY-CODEBLOCK:   pass=11 fail=1  REFUSED
```

## Findings

### [MAJOR] Brace expansion constructs either shell word without `$`, backtick, `\`, or quotes

This is codex-1's finding, which I independently verified. The shared detector
`shellWordView` removes only `\`, `'`, and `"`, and the substitution fallback
triggers only on `$` or backticks:

```js
const shellWordView = (s) => s.replace(/[\\'\"]/g, "");
return /[`$]/.test(s) && (named || flagged);
```

Brace expansion needs none of those characters. A single-element sequence
constructs the exact missing word:

```text
n{o..o}de --test no/such/dir    (builds "node")
node --{t..t}est no/such/dir    (builds "--test")
```

I verified through the probe harness that both are false greens:

```text
BRACE-SINGLE-CHAR-BIN:  pass=12 fail=0  GREEN (FALSE GREEN)
BRACE-SINGLE-CHAR-FLAG: pass=12 fail=0  GREEN (FALSE GREEN)
BRACE-COMMA-BIN:        pass=12 fail=0  GREEN (FALSE GREEN)
```

The shell behavior is measured, not inferred. Single-character ranges expand;
multi-character ranges like `--{test..test}` do NOT (bash/zsh do not expand
multi-char alphabetic ranges to a single value — they pass the literal through):

```text
$ /bin/sh -c 'printf "<%s>\n" n{o..o}de'
<node>
$ /bin/sh -c 'printf "<%s>\n" --{t..t}est'
<--test>
$ /bin/sh -c 'printf "<%s>\n" --{test..test}'
<--{test..test}>
$ /bin/sh -c 'n{o..o}de --test no/such/dir'
Could not find 'no/such/dir'
$ /bin/sh -c 'node --{t..t}est no/such/dir'
Could not find 'no/such/dir'
```

So codex-1's specific reproductions with `--{test..test}` do not actually expand
in the shell — but the single-char forms `n{o..o}de` and `--{t..t}est` DO, and
those are the ones that matter. The class is real: `n{o..o}de` expands to `node`
and runs, while the guard reports success. This blocks release under the same
rule as cycles 23 and 24: a documented verification command can fail for the
reader while the guard reports green.

**Reproduction (fenced):**

```text
$ printf 'Run:\n\n```bash\nn{o..o}de --test no/such/dir\n```\n' > skills/__probe_hermes__.md
$ node --test test/design-addons.test.js
pass=12 fail=0  GREEN
$ /bin/sh -c 'n{o..o}de --test no/such/dir'
Could not find 'no/such/dir'
```

**Blocks release.** Same class as cycle 24, same fix direction: extend the
fail-closed word-building predicate to brace-expansion forms when either half is
recognisable. The symmetry across both passes is sound (both miss it equally); the
shared approximation is incomplete. Fix must preserve canonical target brace
patterns such as `node --test "test/{unit,integration}/*.test.js"` where both
command tokens are already literal.

### [MAJOR] HTML entity `&#10;` / `&NewLine;` produces a line-splitting `\n` in the prose pass

This is a distinct fourth layer. The three layers the brief names — markdown
structure, markdown rendering, shell word construction — are closed. But HTML
character references that resolve to U+000A (line feed) create a `text` node with
a literal `\n` in commonmark's AST, and `flushVisible` splits on `\n` before
matching tokens. The two halves land on separate lines; neither has both `node`
and `--test`; the command is invisible.

The softbreak handling (cycle 21) closed source newlines: a real `\n` in the
markdown source becomes a `softbreak` node, which the guard emits as a space —
one line, both tokens, detected. But `&#10;` and `&NewLine;` are NOT softbreaks;
they are `text` nodes containing `\n`. The guard treats them differently, and the
asymmetry is invisible to the reader: both render as whitespace (HTML collapses
`\n`), so the reader copies one runnable command.

**Reproduction (prose, broken target):**

```text
$ printf 'node&#10;--test no/such/dir\n' > skills/__probe_hermes__.md
$ node --test test/design-addons.test.js
pass=12 fail=0  GREEN (invisible)
```

**Reproduction (prose, valid target — false green):**

```text
$ printf 'Run node&#10;--test "skills/parley-tracker/bin/*.test.js" now.\n' > skills/__probe_hermes__.md
$ node --test test/design-addons.test.js
pass=12 fail=0  GREEN (false green)
```

**Parser evidence** that `&#10;` produces `text "\n"`, not `softbreak`:

```text
node&#10;--test no/such/dir
  => text "node", text "\n", text "--test no/such/dir"

node\n--test no/such/dir  (real newline)
  => text "node", softbreak "", text "--test no/such/dir"
```

**Full extractor confirmation:**

```text
publishedTestCommands("node&#10;--test no/such/dir")          => []
publishedTestCommands("node&NewLine;--test no/such/dir")      => []
publishedTestCommands("node\n--test no/such/dir")             => [{origin:"prose", ...}]
```

Contrast: `&#13;` (CR) and `&#9;` (tab) do NOT split the line and are correctly
detected and refused:

```text
ENTITY-CR-PROSE:  pass=11 fail=1  REFUSED-PROSE
ENTITY-TAB-PROSE: pass=11 fail=1  REFUSED-PROSE
```

In a code block, entities are not resolved (the literal stays
`node&#10;--test no/such/dir`), so the code pass detects and refuses it:

```text
ENTITY-CODEBLOCK: pass=11 fail=1  REFUSED
```

Only the prose pass is open. This is the same pass asymmetry pattern as round 19
(kimi-1): the code pass handles a construct and the prose pass does not — except
here the construct is entity-produced `\n`, not shell quoting.

**Blocks release.** A reader who copies from the rendered page gets
`node --test <targets>` as one command; if the target is broken, the command
fails while the guard reports green. Fix direction: in the `text` case of the
walker, replace `\n` (and `\r`) in the resolved literal with a space before
emitting — the same treatment as `softbreak`. This cannot hide a command, only
reveal one. Alternatively, normalize `text` node literals the same way
`softbreak` is handled. The fix is small and localised to the prose pass's
`text` handler.

### [NIT] codex-1's specific `--{test..test}` reproduction does not expand

codex-1's review uses `node --{test..test} no/such/dir` as the flag-arm
reproduction. I verified that multi-character alphabetic brace ranges do NOT
expand in `/bin/sh` or `zsh` — the shell passes `--{test..test}` through
literally, and node rejects it as a bad option:

```text
$ /bin/sh -c 'printf "<%s>\n" --{test..test}'
<--{test..test}>
$ /bin/sh -c 'node --{test..test} no/such/dir'
node: bad option: --{test..test}
```

The finding's CLASS is correct (brace expansion is invisible to the guard), but
the specific multi-char form does not produce a runnable `--test`. The
single-char forms `n{o..o}de` and `--{t..t}est` DO expand and DO run — those are
the reproductions that matter. This does not change the severity: one verified
false green is enough to block.

## Open questions

The `&#10;` entity gap and the brace expansion gap are both in the shared
detector / prose-pass layer, not in the grammar. Both are fixable without
touching `SUPPORTED_COMMAND` or the provenance rule. The brace expansion fix is
harder (must not over-refuse canonical target patterns like
`"test/{unit,integration}/*.test.js"`), while the entity fix is straightforward
(normalize `\n`/`\r` in text nodes the same way as softbreak). If the implementer
prefers to ship the entity fix first and defer brace expansion, that is a
defensible sequencing — but both must be closed before release.

### Signoff: hermes-1 — 2026-07-30
Status: ❌ BLOCK
