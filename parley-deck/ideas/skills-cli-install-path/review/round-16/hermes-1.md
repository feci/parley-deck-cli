---
idea: skills-cli-install-path
review-round: 16
agent: hermes-1
date: 2026-07-30
reviewed-commit: c45601f
---

## Summary

The cycle-16 fix correctly separates container-marker removal from content-whitespace
preservation, and all 22 supplied probes reproduce exactly as claimed: 15 REFUSED, 3 GREEN,
1 RAN-AND-FAILED, baseline 12/0. The full suite is 253/0 on a clean tree. The guard did not
degrade into a universal refuser — the broken-path probe (P8) still runs and fails.

But cycle-16 reintroduces a class of bug that cycles 10–12 fought and cycle 12 claimed to have
closed: applying a markdown rule to shell text because fence state is not tracked.
`logicalLines` strips `>` from every physical line — including lines inside fenced code blocks,
where `>` is a shell redirect operator, not a blockquote marker. A published command that
starts with `>` inside a fence is stripped to a bare `node --test …`, passes the guard green,
and exits 127 for a reader who copies it.

## What I verified and how

All measurements were run in my isolated worktree at
`/private/tmp/claude-501/…/scratchpad/wt-hermes` (commit `c45601f`), not the shared checkout.

Baseline:

    $ node --test test/design-addons.test.js
    ℹ tests 12  ℹ pass 12  ℹ fail 0

    $ npm test
    ℹ tests 253  ℹ pass 253  ℹ fail 0

Harness (22 probes):

    $ zsh …/scratchpad/probe-hermes.sh
    baseline                     pass=12 fail=0
    P1–P3, P5–P6, P9–P18         pass=11 fail=1  REFUSED
    P4, P7, P19                  pass=12 fail=0  GREEN
    P8                           pass=11 fail=1  RAN-AND-FAILED
    P20–P22                      pass=11 fail=1  REFUSED

Every classification matches the author's claim. The RAN-AND-FAILED control (P8) confirms the
guard still executes and catches real failures rather than refusing everything.

Primary finding — guard output vs. reader output:

    # probe file: skills/__probe_hermes__.md
    ```bash
    > node --test "skills/parley-tracker/bin/*.test.js"
    ```

    $ node --test test/design-addons.test.js
    ℹ pass 12  ℹ fail 0                          ← GREEN

    $ /bin/sh -c '> node --test "skills/parley-tracker/bin/*.test.js"' </dev/null 2>&1
    /bin/sh: --test: command not found
    EXIT=127                                      ← reader gets exit 127

Git status after all probes: clean. No tracked files modified. The `> ` redirect created a
stray file named `node` during one probe; removed it and verified `git status --short` is empty.

## Findings

### [MAJOR] `logicalLines` strips `>` inside fenced code blocks, creating a false green

`logicalLines` (line 256) runs `line.replace(/^ {0,3}> ?/, "")` on every physical line
unconditionally — it does not track whether the line is inside a fenced code block. Outside a
fence, `>` at the start of a line is a markdown blockquote marker and stripping it is correct.
Inside a fence, `>` is literal shell text — a redirect operator — and stripping it changes the
command's meaning.

Reproduction (written to a temporary file under `skills/`, then removed):

    ```bash
    > node --test "skills/parley-tracker/bin/*.test.js"
    ```

The guard strips `> `, sees `node --test "skills/parley-tracker/bin/*.test.js"`, matches
`SUPPORTED_COMMAND`, runs it through `/bin/sh`, reads 35 pass / 0 fail, and reports 12 pass /
0 fail — GREEN.

A reader who copies the line from the rendered fenced block gets the raw `>`:
`/bin/sh` interprets `> node` as "redirect stdout to a file named `node`" and then tries to
execute `--test` as a command. Result: `--test: command not found`, exit 127.

The same false green appears in every fenced form:

| form | guard | reader |
|---|---|---|
| ```` ```bash \n > node --test "…" \n ``` ```` | GREEN (12/0) | exit 127 |
| `~~~bash \n > node --test "…" \n ~~~` | GREEN (12/0) | exit 127 |
| `    > node --test "…"` (indented block) | GREEN (12/0) | exit 127 |
| `> ```bash \n > > node --test "…" \n > ``` `` | GREEN (12/0) | exit 127 |

The indented-block case (4 leading spaces) is doubly affected: `logicalLines`' regex
`/^ {0,3}> ?/` does not match (it allows only 0–3 spaces), but `publishedTestCommands` (line 275)
runs `text.replace(/^[\s>]*/, "")` which strips all leading whitespace and `>` regardless of
count. So the `>` is removed by the second strip even when the first one misses it.

No shipped file currently publishes a `node --test` command with `>` at the start of a line
inside a fence. But the guard is preventive — its own comment says "if a shipped file prints
a `node --test` command anywhere, in any context, that command should work." A future file
that includes a `>` redirect in a fenced block would pass the guard and fail for the reader.

This is the same class as cycles 10–12: the guard applied markdown rules (span parsing, fence
detection) to shell text because it did not track fence state. Cycle 12 removed fence tracking
entirely, replacing it with the discriminator "does a backtick span contain the whole
command?" Cycle 16's blockquote-marker stripping brings back a rule that needs fence state to
be applied correctly. The tension is structural: `>` must be stripped for blockquote
continuations (round 14) and must not be stripped inside fences (this finding), and
distinguishing the two requires knowing whether the line is inside a fence.

Suggested fix: either (a) track fence state so `>` is stripped only outside fences, or
(b) stop stripping `>` in `logicalLines` and instead refuse any command line that still
starts with `>` after whitespace trimming — treating it as unsupported shell syntax rather
than guessing. Option (b) would turn the P19 blockquote-valid-single probe from GREEN to
REFUSED, but no shipped file uses that form, and the trade-off is the same one already
accepted for continuations.

### [MINOR] `publishedTestCommands` strips `$` inside fenced code blocks

Line 275: `text.replace(/^\$\s+/, "")` strips a leading `$` prompt. Inside a fenced block,
`$` is literal text. A reader who copies `$ node --test "…"` gets `$: command not found`
(exit 127) while the guard runs the stripped `node --test "…"` and reports GREEN.

This is an intentional design choice — the comment says "A command never begins with > or $"
— and `$` as a prompt indicator is a ubiquitous documentation convention that most readers
know to omit. I rate it MINOR rather than MAJOR for that reason. But it is the same structural
issue: a noise-stripping rule that assumes `$` is never literal shell text, which is false
inside a fenced block.

### [NIT] Design: 16 cycles of hand-written reconstruction is itself a signal

The guard has been revised 16 times, each revision fixing one edge case of a hand-written
markdown-plus-shell reconstruction. The pattern is consistent: every revision is correct for
the cases its author imagined and wrong for a case a reviewer measured. The root cause is that
reconstructing what a reader copies requires knowing the markdown context (fence state,
blockquote nesting, span boundaries), but the guard deliberately does not track that context
because cycle 12 found that doing so meant "reimplementing CommonMark inside a test."

The tension is unresolved: cycle 16 needs blockquote-marker stripping (a markdown rule) but
cannot apply it correctly without fence state (a markdown context). A narrower rule —
requiring a published verification command to stand alone on its own line with no prefix —
would not close this class by construction, because the `>` is stripped before the grammar
sees it; the grammar never gets the chance to refuse. The class can be closed either by
restoring minimal fence tracking (just enough to know whether `>` is a marker or literal) or
by refusing any line that starts with `>` after whitespace trimming, accepting that the
legitimate blockquoted-fence form (P19) becomes a refusal rather than a green.

## Trade-off judgment

A legitimate command split across lines by a backslash continuation is refused rather than
run. No shipped file uses the form. This is an acceptable fail-closed policy — the same
judgment codex-1 recorded in round 15 and the author restated in cycle 16. I concur.

## Prior finding dispositions (spot-checked)

- Rounds 01–12: the layout move, installer, packaging, Gemini reconciliation, and README
  panel are unchanged since round 13. The addons/ guard fires on reintroduction. Both
  published commands (159/159 and 35/35) still pass. No regression.
- Round 13 MAJOR (continuation before `--test`): FIXED — P10, P11, P12 all REFUSED.
- Round 14 (blockquote continuation): FIXED — P13, P16, P17, P18 all REFUSED.
- Round 15 MAJOR (whitespace deletion merges `node` and `--test`): FIXED — P20, P21, P22
  all REFUSED. The cycle-16 separation of marker removal from whitespace preservation is
  correct for the continuation class.
- Round 15's named divergence ("the guard reconstructs more aggressively than a shell"):
  WITHDRAWN by the author in cycle 16, correctly. The reconstruction now matches what the
  reader copies — except inside fences, where this finding applies.

### Signoff: hermes-1 — 2026-07-30
Status: ❌ BLOCK
