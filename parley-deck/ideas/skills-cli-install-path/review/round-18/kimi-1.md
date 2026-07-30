# Round 18 review — kimi-1

## Summary

Reviewed `265eb56` in the supplied `wt-kimi` worktree (clean except the `node_modules`
link; diff against `15ed1ad` touches only `test/design-addons.test.js` and the two
`commonmark` pins). Suite reproduces at 253/253. All 46 supplied probes behave exactly as
claimed: 37 refused, P4/P7/P19/P26 green, P41/P42/P43 ignored, P45 left alone, P8
RAN-AND-FAILED. All six round-17 dispositions are genuinely fixed — I re-measured each on
my own fixtures, including the arms the supplied harness does not carry (the `<del>` +
canonical-command arm of my own MAJOR, and the valid-target arm of the image-alt MINOR).
The occurrence-level rule is real: both shipped documents yield only `origin:"code"`
occurrences, so the `Verify: ` lines and the parenthesised checklist span pass for the
reason the comment says they do. The guard still verifies — P8 runs and fails through the
real executor, and the four green forms execute real tests with `pass > 0` asserted from
the child's own summary.

The class is nevertheless still not closed, and the open door is the same one as always:
**detection is narrower than the set of commands a reader can run.** Cycles 17–21 modelled
what markdown renders. What remains unmodelled is what the *shell* does with the rendered
text. Detection matches the literal characters `node` and `--test`, but a shell constructs
words by escape removal, quote removal, and command substitution — so
`` `node --\test no/such/dir` ``, published in a single inline code span (the contract's
own canonical container), is invisible to the guard in **both** passes (the extractor
returns `[]`), while the reader who copies it gets `node --test no/such/dir` and an exit-1
failure. The valid-target arm is worse in principle: the same span with the tracker glob
runs 35 real tests the guard never executed, under a test named "every `node --test`
command a shipped file publishes runs tests and passes". One MAJOR, one MINOR, one NIT.
`Status: ❌ BLOCK`.

## Dispositions of the round-17 findings — verified, not trusted

| finding | verdict | evidence |
|---|---|---|
| `codex-1` [MAJOR]: visible text belongs to no provenance bucket | **FIXED** | P32–P35 (mixed span, `no<span></span>de`, `--te<span></span>st`, `<div>` block) all `pass=11 fail=1 REFUSED-PROSE`. Tags are stripped as markup, so the split tokens join in the visible buffer; the per-character owner map then shows the tokens share no code node. The fixture asserts `originOf("node --test html/inline-node.test.js") === "prose"`. |
| `codex-1` [MINOR]: Map masked a prose occurrence | **FIXED** | `publishedTestCommands` returns an array of `{command, origin}` occurrences. P36 and my D2 (prose dup + code dup, identical text) → `REFUSED-PROSE`; the fixture asserts `originOf("node --test dup/same.test.js") === "code+prose"`. The prose occurrence is refused even though the code one is valid. |
| `kimi-1` [MAJOR]: comments / formatting tags / GFM strikethrough as splice points | **FIXED** | P37–P40 all `REFUSED-PROSE`. My D1 re-runs the arm that mattered most — `node --t<del>e</del>st "skills/parley-tracker/bin/*.test.js"`, the **working** command synthesized out of prose, GREEN in round 17 — now `REFUSED-PROSE`. D3 (comment splice + canonical command) likewise. `~~` runs are dropped from the visible buffer; code-span literals are left alone, correctly, since renderers leave them alone too. |
| `kimi-1` [MINOR]: invisible text policed; image alt **executed** | **FIXED** | P41/P42/P43 and my D4 (multiline `<script>` block) all `pass=12 fail=0` — ignored, not policed. P43 is the discriminating probe for execution: its target is `no/such/dir`, so if the alt-text command were still executed the run would print `RAN-AND-FAILED`; GREEN proves it is neither run nor refused. The fixture asserts none of the three invisible forms is recorded at all. |
| `hermes-1` [MAJOR]: soft break renders as a space | **FIXED** | `softbreak` now emits `" "`, `linebreak` keeps `"\n"`. P44 and my D6 (soft-broken **canonical** command) → `REFUSED-PROSE`. P45 (hard break) → GREEN, and the fixture asserts the spliced hardbreak command is *not* recorded — two lines stay two lines. |
| `hermes-1` [NIT] ×2: case-sensitivity; unpinned parser | **FIXED** | `mentionsATestCommand` is `/\bnode\b/i`; P46 (`Node --test …missing…`) is detected and then `REFUSED` by the still-exact grammar. `package.json` and `package-lock.json` both pin `"commonmark": "0.31.2"` (diff: `^0.31.2` → `0.31.2`). |

No round-17 finding remains open in its original form, and no fix regressed another: the
three cycles landed in sequence and the full 46-probe battery plus my own 16 probes are
consistent end to end.

## What I verified, and how

Worktree state:

```text
$ git rev-parse HEAD && git status --short
265eb56b0bfe9a9634b750605853893f23a705c8
?? node_modules
$ git diff 15ed1ad..265eb56 --stat
 package-lock.json          |   2 +-
 package.json               |   2 +-
 test/design-addons.test.js | 253 ++++++++++++++++++++++++++++++++++++++-------
```

Suite and supplied battery:

```text
$ npm test
ℹ tests 253  ℹ pass 253  ℹ fail 0

$ zsh …/scratchpad/probe-kimi.sh
baseline                     pass=12 fail=0
P1–P3, P5–P6, P9–P18, P20–P25, P30–P31          REFUSED      (26)
P27–P29, P32–P40, P44                           REFUSED-PROSE (12, i.e. 37 refused total)
P4, P7, P19, P26                                GREEN        (valid forms)
P41, P42, P43                                   GREEN        (invisible forms ignored)
P45                                             GREEN        (hard break left alone)
P46                                             REFUSED      (case)
P8                                              RAN-AND-FAILED
```

My own round-18 battery (own fixture `skills/__probe_kimi_r18__.md`, removed after each
probe; full script left at `…/scratchpad/probe-kimi-r18.sh`):

```text
D1 del-tag-canonical-cmd       pass=11 fail=1  REFUSED-PROSE
D2 dup-prose-and-code          pass=11 fail=1  REFUSED-PROSE
D3 comment-splice-canonical    pass=11 fail=1  REFUSED-PROSE
D4 script-block-multiline      pass=12 fail=0  GREEN
D5 image-alt-valid-cmd         pass=12 fail=0  GREEN
D6 softbreak-canonical         pass=11 fail=1  REFUSED-PROSE
F4 prose-flag-order-normal     pass=12 fail=0  GREEN
```

**The guard still verifies, and still refuses only what it should.** P8's broken path runs
through `/bin/sh -c` and fails honestly. P4/P7/P19/P26 execute the two real shipped
commands, with `fail == 0` and `pass > 0` read from the child's summary. The false-positive
check the brief asked for resolves cleanly — every occurrence in both shipped documents is
code-owned, measured directly:

```text
skills/parley-tracker/templates/subtask.md
  → [{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]
skills/parley-design-check/SKILL.md
  → [{"command":"node --test \"skills/parley-design-check/test/*.test.js\"","origin":"code"}]
```

That is the `Verify: ` line (`subtask.md:68`), the parenthesised checklist span
(`subtask.md:74`), and the fenced block (`SKILL.md:371-373`) all passing for exactly the
stated reason: each command's two tokens live in one code node, so the visible-text pass
skips them and the code pass runs them. A whole-visible-line rule would have refused all
three; the occurrence-level rule does not. The wider net has not swallowed anything the
project ships.

## Findings

### [MAJOR] Detection matches the page's characters; the shell constructs the words — escape, quote and substitution synthesis are invisible in both passes

Cycle 17's own header comment states the design principle: *"Detection stays deliberately
BROADER than acceptance. Every false green in seventeen rounds came from the two being the
same pattern: whatever the grammar would refuse, detection also failed to see, so the
command was skipped — and skipping reads as success."* That is precisely what still
happens, one layer below markdown. The grammar already knows `\`, `` ` `` and `$` are
shell-significant — it refuses them on sight — and execution was delegated to `/bin/sh -c`
in cycle 9 precisely because the shell is the authority on what a copied line does. But
detection still asks whether the literal characters `--test` appear, and a shell builds
that word from characters that do not contain it.

Reproduction — each line is the entire body of a one-file probe dropped into `skills/`,
run through the guard exactly as `probe-kimi.sh` does; the extractor was also called
directly (`[]` = the guard records no occurrence at all):

```text
`node --\test no/such/dir`                              guard: pass=12 fail=0 GREEN   extractor: []
`node --\test "skills/parley-tracker/bin/*.test.js"`    guard: pass=12 fail=0 GREEN   extractor: []
`nod\e --test no/such/dir`                              guard: pass=12 fail=0 GREEN   extractor: []
`node --te''st no/such/dir`                             guard: pass=12 fail=0 GREEN   extractor: []
Run node --\test no/such/dir now.        (prose)        guard: pass=12 fail=0 GREEN   extractor: []
```fence: node --t`echo e`st no/such/dir                guard: pass=12 fail=0 GREEN   extractor: []
```

What the reader gets, measured against the same shell the guard delegates to:

```text
$ /bin/sh -c 'node --\test no/such/dir'                        → exit=1  Could not find 'no/such/dir'
$ /bin/sh -c 'node --\test "skills/parley-tracker/bin/*.test.js"' → exit=0  ℹ pass 35  ℹ fail 0
$ /bin/sh -c 'nod\e --test no/such/dir'                        → exit=1  Could not find 'no/such/dir'
$ /bin/sh -c "node --te''st no/such/dir"                       → exit=1  Could not find 'no/such/dir'
$ /bin/sh -c 'node --t`echo e`st no/such/dir'                  → exit=1  Could not find 'no/such/dir'
```

And what the reader sees — the backslash survives rendering into the copyable element, so
the copy really does carry it to the shell:

```text
commonmark render of the first probe:  <p><code>node --\test no/such/dir</code></p>
```

The first probe is the round-1 sin: a command the documentation publishes fails for the
reader while the guard reports success. The second is the guard's raison d'être inverted:
a published command runs 35 real tests and the test named "every `node --test` command a
shipped file publishes runs tests and passes" never ran it. The last is the shape cycle 9
thought it had killed for good — a substitution refused *when seen* is here simply never
seen, because the substitution produces the flag itself.

A smaller instance of the same mismatch: the round-17 case fix covered the binary name but
not the flag. `` `node --TEST no/such/dir` `` is invisible (`extractor: []`) and fails for
the reader (`node: bad option: --TEST`, exit 9). Detection is case-insensitive about
`node` and case-sensitive about the one thing that must match how *node*, not the
filesystem, parses.

What this refutes in writing: the contract ("anything else that renders as such a command
is refused by name" — these *execute* as such a command and are not even detected), the
header comment quoted above (skipping reads as success — verbatim what happens), and the
test's own name for the N2 arm.

One honest asymmetry against rounds 16–17: there, rendering *hid* the synthesis from the
reader's eye; here the reader sees the backslash and only the shell removes it. I do not
think that changes the verdict. The guard's entire premise since cycle 9 is that
copy-paste execution must be trustworthy without the reader re-auditing what they pasted,
and the measured outcome — published command fails for the reader, suite green — is
identical to every MAJOR this deliberation has ratified from round 13 onward. It blocks a
release by the standard every prior round applied; calling it a follow-up would be the
first time that shape was waved through.

Direction, not prescription: detection needs the shell's word view, not the page's
character view. The consistent move is the one cycle 9 made for execution — delegate: lex
each candidate unit through the same `/bin/sh` (a no-op lex, not an execution) before
deciding whether it mentions a command. A cheaper fail-closed approximation is to treat
any unit that mentions `node` and carries `\`, `` ` ``, `$` or adjacent-quote construction
near the flag as detected, and let the strict grammar refuse it — three of those four
characters are already grammar-excluded, so the refusal path is in place; only the
detection gap needs closing.

### [MINOR] The occurrence rule polices split-span prose only when `node` comes first

```text
Use `node` with `--test no/such/dir` to verify.    → pass=11 fail=1  REFUSED-PROSE
   recorded: "node with --test no/such/dir to verify." (origin: prose)
Pass the `--test` flag to `node` when you want the runner.  → pass=12 fail=0  GREEN (P26/F4)
```

The visible line in the first case is *not* a runnable command — `with` sits between the
tokens — so refusing it exceeds the contract's stated scope ("anything else that renders
as such a command"); it polices a mention, not a publication. The author's own fixture
blesses the flag-first phrasing as GREEN, so the rule is order-dependent in a way nothing
documents. No shipped file trips it (the suite is green), it fails closed with a named
message and the file, and the rewording cost is trivial — but a future author who writes
"run `node` with `--test <glob>`" in prose will hit a refusal the contract does not
predict. Follow-up, not a blocker; if the MAJOR's detection work happens, the same pass is
the place to make the rule symmetric or to document the asymmetry.

### [NIT] The invisible-text model is still syntactic at the edges

Two residuals of my round-17 MINOR, both measured:

```text
<!-- note a > b: node --test no/such/dir -->       → REFUSED-PROSE  ("node --test no/such/dir -->")
<div>
<script> var c = "node --test no/such/dir"; </script>
</div>                                             → REFUSED-PROSE  ("node --test no/such/dir\";")
```

The tag stripper `<[^>]*>` stops at the first `>`, so a comment *containing* `>` leaks its
tail into the visible buffer; and the script-body exception only fires when the html_block
*starts* with `<script`, so a script nested in a `<div>` block is still treated as page
text. Both police text no reader can reach — the exact complaint of round 17, narrowed
(for the common forms) but not closed. No shipped file trips either; the suite's greenness
proves absence. Follow-up.

### Signoff: kimi-1 — 2026-07-30
Status: ❌ BLOCK
