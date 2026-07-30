---
idea: skills-cli-install-path
review-round: 17
agent: hermes-1
date: 2026-07-30
reviewed-commit: 15ed1ad
---

## Summary

The rewrite from a hand-written line scanner to a real CommonMark parser is the right
architectural move and fixes all five round-16 findings. The contract is narrow and
checkable, the provenance rule is sound, and the guard still executes and catches real
failures (P8: RAN-AND-FAILED). All 31 harness probes behave as claimed, and the suite is
253/253.

I found one real gap: a soft line break in prose splits a command across two `text` nodes,
and the guard's `flushRendered` pushes `"\n"` for the `softbreak` node — but GitHub renders
a soft break as a space. A command that spans a soft break in prose renders as a single
copyable line for the reader and as two lines for the guard, so it is not detected and not
refused. This is the same "detection model differs from reader model" class that round 16
identified, now in the parser layer rather than the scanner layer. It does not produce a
false green by running a mutated command — the command is simply invisible. No shipped file
currently uses this pattern. I rate it MAJOR because it is a concrete gap in the "refused
by name" claim, with a one-line fix.

Everything else I attacked held. The GFM extensions I tested (strikethrough, tables,
`<code>`, `<pre>`, `<kbd>`, autolinks, footnotes) are all handled correctly — either
detected as prose and refused by provenance, or detected as code and checked by form. The
"class is closed by construction" claim is substantiated for code nodes and substantially
substantiated for prose, with the soft-break exception noted.

## Dispositions of round-16 findings (verified, not taken on trust)

I verified each fix against the 31-probe harness and with standalone Node scripts that
replicate the guard's `publishedTestCommands` and run the actual `commonmark` parser.

### agy-1: node flag between tokens → command skipped

FIXED. Detection is now `mentionsATestCommand = (s) => /\bnode\b/.test(s) && /--test\b/.test(s)`
— deliberately broader than the acceptance grammar. P23 (`node --no-warnings --test …`) and
P24 (`node -r ./setup.js --test …`) are both REFUSED (detected, then rejected by
`SUPPORTED_COMMAND`). P25 (flags with a valid target) is also REFUSED — the grammar
rejects anything between `node` and `--test`. The broader-than-acceptance design is correct:
the false greens in prior rounds came from detection and acceptance sharing the same pattern.

### codex-1: markdown rendering synthesizes commands (escape, emphasis, entity)

FIXED. The parser resolves escapes, emphasis and entities before the walker sees `text`
nodes. P27 (`node \--test …` → escape renders as `--test`), P28 (`node --**test** …` →
emphasis renders as `--test`), and P29 (`node&#32;--test …` → entity renders as space) are
all detected as `prose` provenance and REFUSED-PROSE. I confirmed the provenance
assertion fires: the guard's error message says "must be published as its own code span."
The parser, not a regex, is what resolves the rendering — so the synthesis is visible to
the guard for the first time.

### kimi-1 and hermes-1: `>` inside a fence stripped and the mutated text executed

FIXED. The parser knows the container. A `>` inside a fenced code block is literal content
of a `code_block` node; a `>` starting a line outside a fence is a blockquote marker that
the parser consumes. P30 (`> node --test …` inside a fence) is REFUSED — the literal
includes the `>`, and `SUPPORTED_COMMAND` rejects it. P19 (`> node --test …` inside a
blockquote-wrapped fence) is GREEN — the parser strips the blockquote marker, the code
block literal is the bare command, and it runs correctly. The distinction the comment
claims — "a fence inside a blockquote yields the bare command; a `>` inside a fence stays
in the literal" — is exactly what the parser provides.

### kimi-1: zero-width continuation boundary

FIXED. P31 (`node\` + `--test …` with no space before the backslash) is REFUSED. The
`codeNodeUnits` function splices the continuation and emits it with the backslash restored
(`"node--test cont/zero-width.test.js \\"`), which `SUPPORTED_COMMAND` rejects.

### hermes-1: `$ ` prompt stripped inside fenced blocks

FIXED. The guard no longer strips `$ `. The `$` is content of the code node, so
`publishedTestCommands` returns `"$ node --test b/dir"` as a `code` provenance command,
and `SUPPORTED_COMMAND` rejects it (the `$` is not in the allowed character class). The
fixture assertion at line 444 confirms `"$ node --test b/dir"` is captured whole, and
line 498 confirms the grammar refuses it.

## What I verified and how

All measurements were run in my isolated worktree at
`/private/tmp/claude-501/…/scratchpad/wt-hermes` (commit `15ed1ad`), not the shared
checkout.

Baseline:

    $ node --test test/design-addons.test.js
    ℹ tests 12  ℹ pass 12  ℹ fail 0

    $ npm test
    ℹ tests 253  ℹ pass 253  ℹ fail 0

Harness (31 probes):

    $ zsh …/scratchpad/probe-hermes.sh
    baseline                     pass=12 fail=0
    P1–P3                        REFUSED
    P4, P7, P19, P26             GREEN
    P8                           RAN-AND-FAILED
    P5–P6, P9–P18, P20–P25       REFUSED
    P27–P29                      REFUSED-PROSE
    P30–P31                      REFUSED

Every classification matches the author's claim. P8 (genuinely-broken-path) confirms the
guard still executes commands and catches real failures — it does not refuse everything.

GFM/CommonMark attack surface — standalone Node scripts using the actual `commonmark`
parser (version 0.31.2, locked in `package-lock.json`):

    GFM strikethrough         prose:FORM-REFUSED  (~~ renders as literal in CommonMark)
    HTML <code> inline        prose:FORM-REFUSED  (raw html_inline, not a code node)
    HTML <pre> block          prose:FORM-OK       (raw html_block → refused by provenance)
    GFM table code span       code:FORM-OK        (CommonMark doesn't parse tables; pipes
                                                   are literal but the backtick span still
                                                   works as a code node)
    GFM table text            prose:FORM-REFUSED  (pipes in text, refused by provenance)
    Case: Node --test         NOT-DETECTED        (detection is case-sensitive — see NIT)
    Case: NODE --test         NOT-DETECTED
    Numeric entity            prose:FORM-OK       (refused by provenance)
    Link destination          prose:FORM-REFUSED
    Ref link def              prose:FORM-REFUSED
    Image alt text            prose:FORM-OK       (refused by provenance)
    HTML <kbd>                prose:FORM-REFUSED
    Setext heading            prose:FORM-OK       (refused by provenance)
    Blockquote prose          prose:FORM-OK       (refused by provenance)
    Link with code            code:FORM-OK        (code span in link text — correct)
    Footnote ref              code:FORM-OK        ([^1]: `…` — code span, correct)
    Split across spans        NOT-DETECTED        (neither span has both tokens — see below)
    Code span entity          code:FORM-REFUSED   (entity is literal in code spans)
    Code span backslash       code:FORM-REFUSED   (backslash is literal in code spans)
    HTML comment              prose:FORM-REFUSED
    Details summary           code:FORM-OK        (code span inside <details>)

The `<pre>` HTML block case deserves note: the guard classifies it as `prose:FORM-OK`
because the raw HTML reaches the `html_block` handler and is pushed to `rendered`. The
provenance is `prose`, so the guard refuses it. This is correct — `<pre>` is not a
CommonMark code node, and a reader's interaction with `<pre>` content depends on the
renderer, not the parser. Refusing it is fail-closed.

The "split across spans" case (`Use \`node\` with \`--test x\``) is NOT-DETECTED — neither
span contains both tokens. This is correct: the reader sees two separate code spans, not
one command. There is nothing to copy as a single command.

Inline code span with newline: the parser normalizes newlines to spaces in code span
literals (CommonMark §6.1), so `node\n--test x` inside a span produces literal
`"node --test x"` — the guard sees the same text the reader copies. No gap here.

Git status after all probes: clean. No tracked files modified. All probe scripts were
untracked and have been removed.

## Findings

### [MAJOR] Soft break in prose splits a command across two text nodes — invisible to the guard

`flushRendered` pushes `"\n"` for `softbreak` nodes, then splits the joined rendered text
by `"\n"` to check each line. But GitHub renders a soft break as a space (CommonMark §6.7:
"rendered as a space in HTML"). A command that spans a soft break in prose is one copyable
line for the reader and two lines for the guard.

Reproduction (standalone script using the actual parser):

    const { Parser } = require("commonmark");
    const md = "node\n--test missing/definitely-missing.test.js";
    // Parser AST: text("node"), softbreak, text("--test missing/definitely-missing.test.js")
    // Guard flushRendered: rendered = ["node", "\n", "--test missing/..."]
    // Joined: "node\n--test missing/..."
    // Split by \n: ["node", "--test missing/..."]
    // Neither line has both \bnode\b AND --test\b → 0 commands detected

Through the actual guard test harness:

    # probe file: skills/__probe_softbreak__.md
    # Content:
    #   To verify, run node
    #   --test "skills/parley-tracker/bin/nonexistent-softbreak.test.js" from root.
    #   Then also run `node --test "skills/parley-tracker/bin/*.test.js"` for real.

    $ node --test test/design-addons.test.js
    ℹ pass 12  ℹ fail 0                          ← GREEN

The soft-break command is not detected. The guard does not run it, does not refuse it,
does not know it exists. A reader who copies from the rendered GitHub page gets:

    node --test "skills/parley-tracker/bin/nonexistent-softbreak.test.js"

…which is a runnable command that exits with "Could not find" — a command that verifies
nothing while looking like it does.

This is the same class as agy-1's round-16 finding (detection narrower than the set of
commands that render), now in the parser's prose-collection layer rather than the
scanner's detection pattern. The contract says "anything else that renders as such a
command is refused by name" — but this command is neither refused nor detected.

No shipped file currently uses this pattern (verified by searching all `.md` files in
`skills/` for a line ending in `node` followed by a line starting with `--test`).

Fix: in `flushRendered`, join `softbreak` nodes with a space (`" "`) instead of `"\n"`,
matching the HTML rendering semantics that the reader sees. `linebreak` (hard break,
trailing spaces or backslash) should keep `"\n"` — the reader sees two lines there.

### [NIT] Detection is case-sensitive — `Node --test` and `NODE --test` are not detected

`mentionsATestCommand` uses `\bnode\b`, which is case-sensitive. `Node --test x` and
`NODE --test x` in a code block are NOT-DETECTED. On case-insensitive filesystems
(macOS default), `node` and `Node` resolve to the same binary, so a reader who copies
`Node --test x` gets a working command. The guard does not see it.

This is a NIT because the project's published commands use lowercase `node` exclusively,
and the convention is universal in Node.js documentation. But it is the same skip-as-green
class: the detection pattern is narrower than the set of commands that work.

### [NIT] The `commonmark` dependency is pinned by `^0.31.2` with a lockfile

The guard depends on `commonmark`'s AST shape and normalization behavior (newline-to-space
in code spans, softbreak vs linebreak distinction, HTML block/inline classification).
These are CommonMark spec behaviors, not implementation quirks, so a spec-conformant
upgrade should not break the guard. The lockfile pins 0.31.2 exactly. A conformance
assertion (e.g., verifying that the parser normalizes newlines in code spans) would
make the dependency's critical properties explicit, but the current arrangement is
tolerable — the behaviors the guard relies on are spec-mandated, not version-specific.

## The "class is closed by construction" claim

The claim is: a published verification command must be a whole code node in canonical
form, and anything else that renders as such a command is refused by name.

For CODE NODES: substantially closed. The parser produces the exact literal a reader
copies from a code span (newlines normalized to spaces per spec), and the exact literal
a reader copies from a code block (newlines preserved). The `codeNodeUnits` function
handles backslash continuations by emitting them with the backslash restored, so the
grammar refuses them. The only gap in code nodes is case sensitivity (NIT above).

For PROSE: the soft-break gap (MAJOR above) means the guard does not see every command
that renders in prose. Commands on a single line are detected and refused by provenance.
Commands split by a soft break are invisible. The fix is one line.

For GFM EXTENSIONS: CommonMark does not parse GFM tables, strikethrough, footnotes, or
autolinks, but the cases I tested are all handled correctly — either the relevant text
reaches the guard as prose (and is refused by provenance) or as a code node (and is
checked by form). GFM tables are the most interesting case: CommonMark treats `|` as
literal text, so a command in a table cell without backticks is prose (refused), and a
command in a table cell with backticks is a code node (checked). Both are correct.

## Whether the guard still verifies anything

Confirmed. P8 (`node --test no/such/dir`) runs through the guard, executes via `/bin/sh`,
fails with "could not find," and the guard reports RAN-AND-FAILED (pass=11 fail=1). The
guard is not a universal refuser — it executes valid commands and catches real failures.

### Signoff: hermes-1 — 2026-07-30
Status: 🟡 ACCEPT WITH RESERVATIONS
