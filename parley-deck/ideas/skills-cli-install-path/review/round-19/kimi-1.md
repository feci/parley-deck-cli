# Round 19 review — kimi-1

## Summary

Reviewed `f5200f7` in the supplied `wt-kimi` worktree (clean except the `node_modules`
link; no tracked file touched, no probe file left behind). Suite reproduces at 253/253.
The supplied 58-run harness matches its claim in behavior on every line: baseline green,
46 refused, 4 valid forms green, **5** invisible forms ignored (the brief says 4 — measured
P41, P42, P43, P56, P57; arithmetic only, every probe behaves as intended), the hard break
left alone, and P8 still RUNS and fails. The two shipped documents still yield exactly one
`origin:"code"` occurrence each, and across **all** shipped `skills/*.md` the extractor
finds exactly those two distinct commands — cycle 23's wider net swallowed nothing shipped.

Round-18 dispositions, measured per finding below: codex-1's MAJOR is genuinely fixed
(cycle 22 is the right shape — drop the inline node, tokenize the block); my NIT is fixed;
my MINOR stands unchanged as the recorded follow-up. My own MAJOR is **not** fully fixed,
and that is this round's finding. Cycle 23 put the shell-word view in `mentionsATestCommand`
— the code-node pass — and left the visible-text pass matching raw characters, so the prose
arm that my round-18 review explicitly measured and listed (`Run node --\test no/such/dir
now.` → `extractor: []`) is still invisible: same characters, same shell, same reader
failure, guard GREEN. And a second arm shows the gap is not only between the passes: the
substitution rule keys on a literal `node`, but a substitution can *produce* `node`, so
`` `$(echo n)ode --test no/such/dir` `` — in a single code span, the canonical container —
defeats even the rewritten pass. One `[MAJOR]`. `Status: ❌ BLOCK`.

## Dispositions of the round-18 findings — verified, not trusted

| finding | verdict | evidence |
|---|---|---|
| `codex-1` [MAJOR]: a quoted `>` (and a `>` in a comment) defeated the tag stripper | **FIXED** | P47/P48/P49 all `pass=11 fail=1 REFUSED-PROSE` at guard level. Extractor on the real functions records each as a prose occurrence: `no<span title="1 > 0"></span>de --test no/such/dir` → `[{"command":"node --test no/such/dir","origin":"prose"}]`, likewise the `<div title="a > b">` block and the `<!-- a > b -->` splice. The design is right: `html_inline` is one complete node and is dropped whole — nothing to parse; `html_block` goes through a quote/comment/CDATA/PI-aware scanner, and an unterminated `<` is kept as text (can only add a refusal). |
| `kimi-1` [MAJOR]: the shell builds the words | **CODE ARMS FIXED — PROSE ARM OPEN** | P50–P55 all `pass=11 fail=1 REFUSED`; extractor records each as a `code` occurrence and `SUPPORTED_COMMAND` rejects every one (measured individually below). **But the prose arm my review measured is still `extractor: []` and guard GREEN** (N1–N5 under Findings). The disposition of this finding is therefore not "fixed"; see [MAJOR]. |
| `kimi-1` [NIT]: script nested in a `<div>` policed; comment with `>` leaked its tail | **FIXED** | P56 and P57 both `pass=12 fail=0 GREEN`; extractor returns `[]` for both — the script body is now skipped wherever it sits, and the comment is scanned to `-->`, so invisible text is no longer policed. The cycle-22 comment claims the common forms; the edges I measured are closed. |
| `kimi-1` [MINOR]: occurrence rule is order-dependent | **UNCHANGED — still a valid follow-up, as recorded** | Extractor still records `Use \`node\` with \`--test no/such/dir\` to verify.` as a prose occurrence (`"node with --test no/such/dir to verify."` → REFUSED-PROSE) while the flag-first phrasing (P26) is GREEN. No shipped file trips it: over all shipped `skills/*.md` the only occurrences are the two code-owned commands. Behavior identical to round 18; it stays a follow-up for the reasons given then. |

No fix regressed another: the full 57-probe battery, the 253-test suite, and my own 11
round-19 probes are consistent end to end.

## What I verified, and how

Worktree state and suite:

```text
$ git rev-parse HEAD && git status --porcelain
f5200f7831fe5f8ad00032d7643800692cc3fbfd
?? node_modules

$ npm test
ℹ tests 253  ℹ pass 253  ℹ fail 0
```

Supplied harness (`zsh …/scratchpad/probe-kimi.sh`, concurrency check intact), counted
from the actual output:

```text
baseline pass=12 fail=0
REFUSED        P1–P3, P5–P6, P9–P18, P20–P25, P30–P31, P46, P50–P55   (30)
REFUSED-PROSE  P27–P29, P32–P40, P44, P47–P49                          (16)
GREEN valid    P4, P7, P19, P26                                        (4)
GREEN ignored  P41, P42, P43, P56, P57                                 (5)
GREEN hard brk P45                                                     (1)
RAN-AND-FAILED P8                                                      (1)
               total 57 — every probe behaves as intended
```

Method note: beyond the harness I loaded the guard's **real** functions out of
`test/design-addons.test.js` (compile the file's source in a `new Function`, neutralizing
only the `test(...)` registration so nothing runs; every measured function is
byte-identical to what the guard executes) and measured `publishedTestCommands`,
`mentionsATestCommand`, `SUPPORTED_COMMAND` and the actual commonmark renders directly.
Detector/grammar truth table from that run:

```text
"node --test \"skills/parley-tracker/bin/*.test.js\""  mentions:true  grammar:true
"node --\\test no/such/dir"        mentions:true  grammar:false   (P50 shape — seen, refused)
"nod\\e --test no/such/dir"        mentions:true  grammar:false   (P52)
"node --te''st no/such/dir"        mentions:true  grammar:false   (P53)
"node --TEST no/such/dir"          mentions:true  grammar:false   (P55)
"node --t`echo e`st no/such/dir"   mentions:true  grammar:false   (P54)
"$(echo n)ode --test no/such/dir"  mentions:FALSE grammar:false   (N6 — never seen)
"n$(echo o)de --test no/such/dir"  mentions:FALSE grammar:false   (N7 — never seen)
```

**The guard still verifies.** P8 and my control C1 (plain `node --test no/such/dir` in a
span) are `pass=11 fail=1 RAN-AND-FAILED` — the command reaches `/bin/sh -c` and fails
honestly. P4/P7/P19/P26 are green through the real executor. The two shipped commands,
run by hand in the worktree:

```text
$ /bin/sh -c 'node --test "skills/parley-tracker/bin/*.test.js"'      → ℹ pass 35  ℹ fail 0
$ /bin/sh -c 'node --test "skills/parley-design-check/test/*.test.js"' → ℹ pass 159 ℹ fail 0
```

**False-positive surface.** Extractor over both shipped documents — the `Verify: ` line
(`subtask.md:68`), the parenthesised checklist span (`subtask.md:74`), and the fenced
block in `SKILL.md`:

```text
skills/parley-tracker/templates/subtask.md -> [{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]
skills/parley-design-check/SKILL.md        -> [{"command":"node --test \"skills/parley-design-check/test/*.test.js\"","origin":"code"}]

all shipped skills/*.md, distinct occurrences: 2 (both origin "code")
```

So the occurrence-level rule still admits exactly what the project legitimately publishes,
and cycle 23's broader detection refuses nothing the project ships. The control matters:
`` `echo $HOME` `` in a span is `extractor: []` — the `$` arm does not fire without a
token — so the net is wider only where a command token is present.

## Findings

### [MAJOR] Layer three was closed for one token in one pass — the shell builds both words, in both contexts

Cycle 23's own header states the principle universally: *"Detection reads the shell's
words, not the page's characters"* and *"A substitution can produce ANY word."* The
implementation applies it to `mentionsATestCommand` only — the code-node pass. The
visible-text pass (`flushVisible`) still scans raw characters:

```js
for (const m of line.matchAll(/\bnode\b/gi)) {   // case-insensitive node — new in cycle 23
  const rest = line.slice(m.index);
  const flag = rest.match(/--test\b/);            // raw characters, case-SENSITIVE — unchanged
```

**Arm A — the prose arm of my own round-18 MAJOR is still invisible.** My round-18
reproduction table listed, and measured, this line: `Run node --\test no/such/dir now.
(prose) → guard: pass=12 fail=0 GREEN, extractor: []`. Cycle 23's claim is *"Every one of
these is now detected and then refused by the grammar."* Measured now, guard level (own
battery `…/scratchpad/probe-kimi-r19.sh`, same one-file discipline as the supplied
harness):

```text
N1 prose-escaped-flag              pass=12 fail=0  GREEN     Run node --\test no/such/dir now.
N1b bare-escaped-flag              pass=12 fail=0  GREEN     node --\test no/such/dir            (whole paragraph)
N2 prose-escaped-binary            pass=12 fail=0  GREEN     Run nod\e --test no/such/dir now.
N3 prose-quote-splice              pass=12 fail=0  GREEN     Run node --te''st no/such/dir now.
N4 prose-uppercase-flag            pass=12 fail=0  GREEN     Run node --TEST no/such/dir now.
N5 prose-escaped-flag-VALID-target pass=12 fail=0  GREEN     Run node --\test "skills/parley-tracker/bin/*.test.js" now.
```

Extractor on each of N1–N5: `[]`. Rendered page (real parser): `<p>Run node --\test
no/such/dir now.</p>` — the backslash survives to the copy, exactly as it does inside the
`<code>` elements of P50–P55, which cycle 23 refuses. What the same shell the guard
delegates to does with the copied text:

```text
$ /bin/sh -c 'node --\test no/such/dir'        exit=1  Could not find 'no/such/dir'
$ /bin/sh -c 'nod\e --test no/such/dir'        exit=1  Could not find 'no/such/dir'
$ /bin/sh -c "node --te''st no/such/dir"       exit=1  Could not find 'no/such/dir'
$ /bin/sh -c 'node --TEST no/such/dir'         exit=9  node: bad option: --TEST
$ /bin/sh -c 'node --\test "skills/parley-tracker/bin/*.test.js"'   → ℹ pass 35 ℹ fail 0
```

N4 is likewise the prose arm of hermes-1's round-17 case finding: cycle 23 made the prose
pass case-insensitive about `node` but not about the flag, in the same line it touched.
N5 is the raison-d'être inversion again: 35 real tests run for the reader, never run by
the test named "every `node --test` command a shipped file publishes runs tests and
passes".

**Arm B — the substitution gate keys on a literal `node` that a substitution need never
spell.** The new rule treats *"a unit that names `node` while carrying a backtick or `$`"*
as a candidate — but the gate runs first:

```js
if (!/\bnode\b/i.test(s) && !/\bnode\b/i.test(words)) return false;
```

A substitution can produce ANY word — including `node`. In a single inline code span, the
contract's own canonical container:

```text
N6 span-subst-builds-node          pass=12 fail=0  GREEN     Run: `$(echo n)ode --test no/such/dir`
N7 span-subst-inside-node          pass=12 fail=0  GREEN     Run: `n$(echo o)de --test no/such/dir`
N8 prose-subst-builds-node         pass=12 fail=0  GREEN     Run $(echo n)ode --test no/such/dir now.
N9 span-subst-node-VALID-target    pass=12 fail=0  GREEN     Run: `$(echo n)ode --test "skills/parley-tracker/bin/*.test.js"`

extractor on each: []          mentionsATestCommand("$(echo n)ode --test no/such/dir") = false
$ /bin/sh -c '$(echo n)ode --test no/such/dir'                              exit=1  Could not find 'no/such/dir'
$ /bin/sh -c '$(echo n)ode --test "skills/parley-tracker/bin/*.test.js"'    → ℹ pass 35 ℹ fail 0
```

This one defeats the pass cycle 23 rewrote, in the container the contract blesses, and it
contradicts the new comment's own sentence. Note that even a perfect no-op lex through
`/bin/sh` would **not** catch it — lexing does not evaluate `$(echo n)` — so this arm is
not an argument for the lexing alternative; the fail-closed rule simply has to trigger on
a substitution that mentions *either* token, not `node` alone.

**Why this blocks rather than follows up.** The standard this deliberation has applied
from round 13 onward is the measured shape: a command the documentation publishes fails
for the reader while the guard reports success — detection narrower than execution,
skipping reads as success. Both arms reproduce it live, one of them in the canonical
container, one with a valid target running 35 real tests unexamined. Arm A is not a new
finding at all: it is the unclosed half of the round-18 MAJOR whose closure the cycle-23
message asserts ("every one of these is now detected") — the finding's own reproduction
list included the prose line. A release-blocking finding whose fix covers the code
containers and misses the prose container is not fixed, and hermes-1's round-19-anticipating
"no remaining false-green surface" is, measured, not the case.

Direction, not prescription: **one detector, two provenance decisions.** Route the visible
line through the same shell-word view the code pass uses — when the word view sees a
candidate the character scan missed, record the occurrence as prose unless both tokens are
wholly owned by one code node (that owner check already exists and is what keeps the two
shipped documents green). And trigger the substitution arm on a unit carrying `` ` `` or
`$` that mentions `node` **or** `--test`: refuse-on-sight needs no evaluation, so nothing
new is executed and the grammar (`[^...$\\...]` already excludes both characters) already
refuses every form this reaches.

**Boundary probed and dismissed, so the contract's edge is stated.** A zero-width or
homoglyph flag (`node --te\u200Bst x`) renders as `--test` to the eye but does not
*execute* as `--test` — node receives a different word and rejects it; that is a typo,
invisible or not, and outside "runs as `node --test`". Glob forms (`--tes[t]`) depend on
the reader's cwd. Neither is a finding; the contract covers what the shell runs as the
command, and every arm above does exactly that.

## The approximation-versus-lexing choice

I named both options in round 18; the approximation was chosen, and I confirm the choice —
for a reason stronger than the one recorded. The recorded reason is that lexing a line
through a live shell risks evaluating the very substitutions under inspection; any
practical "print the words" trick (`set --`, `printf` over unquoted expansion) is
evaluation by another name. Arm B adds the structural reason: even a hypothetical safe,
no-op lexer only *tokenizes* — it does not expand — so `$(echo n)ode` comes back as a word
containing a substitution, and the only sound decision at that point is "refuse on sight
of an unevaluated substitution", which is precisely the fail-closed approximation, one
step later and with a shell in the loop. Lexing buys tokenization fidelity for the
quote-removal cases; `shellWordView` gets those same cases by deletion, and deletion can
only reveal a token, never conceal one. The approximation's failure this cycle is not
that it approximates — it is that it was applied to one of the two detection passes and
gated on one of the two tokens. Fix the application, keep the approximation.

### Signoff: kimi-1 — 2026-07-30
Status: ❌ BLOCK
