---
idea: skills-cli-install-path
review-round: 16
agent: kimi-1
date: 2026-07-30
reviewed-commit: c45601f
---

## Summary

Reviewed branch `readme-skill-catalogue` at `c45601f` in my own worktree. The author's
numbers all reproduce: full suite 253/0, and the supplied harness's 22 probes classify
exactly as claimed (18 REFUSED, 3 GREEN, 1 RAN-AND-FAILED). The cycle-16 fix is a real
correction of the round-15 defect: content whitespace is now preserved through splicing,
and all of rounds 13–15's exact probes are refused.

The class is still not closed. Cycle 16 taught the guard to remove a blockquote **marker**
from every physical line — but whether a leading `>` is a marker or is *content* depends on
the container the line sits in, and this guard tracks no containers by design. A `>` at the
start of a fenced or indented code line is literal text the reader copies verbatim; the
guard deletes it, executes the mutated text, and certifies the file green while the reader
gets exit 127 (plus a file named `node` created in their cwd by the redirection). A second,
independent hole: a continuation broken at the one point with no whitespace on either side
(`node\` + `--test x`) splices to `node--test x`, which the detector `node\s+--test` cannot
see — so that continuation shape is not refused, contradicting the round-brief's claim that
every continuation shape is. Two MAJOR findings, both measured. `Status: ❌ BLOCK`.

## What I verified, and how

All guard-side numbers from `zsh probe-kimi.sh` against my worktree (each probe writes one
temp file under `skills/`, runs `node --test test/design-addons.test.js`, deletes it; the
baseline file has 12 tests). Reader-side numbers from pasting the published text into
`/bin/sh` and `/bin/zsh` in a scratch directory.

**Author's claims — all reproduced:**

- `git log -1` → `c45601f`; `git status --short` empty before and after all probes;
  `git diff --name-only 931e682..c45601f` → only `test/design-addons.test.js`.
- `npm test` → `ℹ tests 253 / ℹ pass 253 / ℹ fail 0` (3.2s). Claim confirmed.
- Harness: baseline 12/0; P1–P3, P5–P6, P9–P18, P20–P22 REFUSED (11/1); P4, P7, P19 GREEN
  (12/0); P8 RAN-AND-FAILED (11/1). Every one of the 22 claimed classifications matches,
  including the important P8 control showing the guard is not a universal refuser.
- Only two distinct commands are published under `skills/` (grep): the fenced
  `node --test "skills/parley-design-check/test/*.test.js"` and the inline
  `` `node --test "skills/parley-tracker/bin/*.test.js"` `` (twice in `subtask.md`).
  Neither uses a container prefix, a prompt, or a continuation — the two strip sites my
  findings attack are exercised only by adversarial shapes, not by real content.

**My seven added probes — real output:**

```text
K1 fenced-gt-invalid      pass=11 fail=1  RAN-AND-FAILED   (``` fence, line "> node --test no/such/dir")
K2 fenced-gt-valid        pass=12 fail=0  GREEN            (same fence, valid target)
K3 indented-gt-invalid    pass=11 fail=1  RAN-AND-FAILED   (4 spaces + "> node --test no/such/dir")
K4 indented-gt-valid      pass=12 fail=0  GREEN            (4 spaces + ">", valid target)
K5 zero-width-missing     pass=12 fail=0  GREEN            ("node\" + "--test …/missing.test.js")
K6 zero-width-valid       pass=12 fail=0  GREEN            ("node\" + "--test valid target")
K7 tab-gt-invalid         pass=11 fail=1  RAN-AND-FAILED   (TAB + "> node --test no/such/dir")
```

**Reader-side, same published texts:**

```text
$ /bin/sh -c '> node --test "skills/parley-tracker/bin/*.test.js"'
/bin/sh: --test: command not found
exit=127                        (and an empty file named "node" created in cwd)

$ /bin/sh cmd.sh                # cmd.sh = node\ ⏎ --test …/missing.test.js
cmd.sh: line 2: node--test: command not found
exit=127
$ /bin/zsh cmd.sh
cmd.sh:1: command not found: node--test
exit=127
```

No tracked file was modified; the worktree probe file is auto-removed by the harness and my
reader-side scratch directory was deleted.

## Findings

### [MAJOR] The guard strips a `>` it cannot prove is a marker, executes the mutated text, and certifies it green

Cycle 16 removes `^ {0,3}> ?` from **every physical line** (`logicalLines`,
`test/design-addons.test.js:253-259`) on the premise — stated in the comment at lines
249-251 — that "Removing exactly the marker leaves what the reader actually copies out of
the rendered page." That premise is false wherever `>` is **content**: inside a fenced code
block, or on a line indented four spaces or a tab, CommonMark treats the text literally and
the rendered page shows the `>` to be copied. A `>` is a blockquote marker only when the
*container* says so, and this guard consults no containers — that is its founding doctrine
(comment at line 222-226: "containers are never consulted").

Measured (K1+K2): a shipped file publishing

    ```bash
    > node --test "skills/parley-tracker/bin/*.test.js"
    ```

leaves the suite **green at 12/0**, while the same line copied from the rendered page exits
**127** — `> node` is parsed by the shell as a redirection, which *creates/truncates a file
named `node` in the reader's cwd*, and the remaining word `--test` is not a command. K1
(same shape, target `no/such/dir` → RAN-AND-FAILED) proves K2's green is the guard
**executing the stripped text**, not skipping it: the mutation runs, and when the mutation
happens to be valid the broken published line is certified. K3/K4/K7 show the second strip
site reaching the same result: for 4-space or tab lead-in the cycle-16 stripper correctly
declines, but the older `text.replace(/^[\s>]*/, "")` at line 275 then deletes the content
`>` anyway — a strip site the cycle-16 narrative does not mention.

The sharpest part: had the `>` survived to `SUPPORTED_COMMAND` (line 309, anchored at
`^node`), the fail-closed grammar would have **refused** these lines by name. The two
pre-grammar strips remove the evidence the grammar needs. Note the asymmetry the fix must
respect: in P19 the fence sits *inside* a blockquote, the `>` is genuinely a marker, the
reader copies the bare command, and stripping is correct. The same physical `>` is marker
in one container and content in the other; a per-line heuristic cannot tell them apart.

**Direction for the fix:** never strip `>` (or anything else) before execution. If a raw
line is not the command, refuse it — the grammar already does — rather than repairing it
into something executable.

### [MAJOR] The zero-width continuation boundary escapes detection entirely — "every continuation shape is refused" is false

Round 13 broke the command with whitespace on both sides of the boundary (`node \` +
`--test x`); round 15 with indentation on the continuation (`node\` + `  --test x`). Cycle
16 handles both because the shell preserves that whitespace into the splice, so the
detector `node\s+--test` still fires and the spliced unit is refused. The remaining
boundary is the one with **nothing on either side**:

    node\
    --test skills/parley-worktrees/kimi-zero-width-missing.test.js

The shell removes backslash-newline with no substitution → `node--test …` → exit 127,
`node--test: command not found` (measured in both `/bin/sh` and `/bin/zsh`). The guard's
splice is faithful — it produces the same `node--test …` — but then the detector at line
273 requires `\s+` between the tokens, finds none, and **skips the unit entirely**. Every
*detected* spliced unit is emitted with a restored backslash and refused (lines 279-282),
so K5/K6's 12/0 green can only mean "never seen". The round-brief's claim that "every
continuation shape — split after `--test`, inside `node`, inside `--test`, between `node`
and `--test`, with and without a space before the backslash — is REFUSED" is refuted by
this measurement: the without-a-space, without-an-indent shape is not refused, it is
invisible. Detection, not reconstruction, is now the weak link: any detector narrower than
"what the shell will try to execute" can be stepped around by choosing the split point that
erases the detector's separator.

### [MINOR] Sixteen cycles of hand-written markdown+shell reconstruction — narrow the publication form instead

Both MAJORs are the same lesson rounds 03–09 taught about shell, now taught about markdown:
a hand-written approximation diverges from the thing that defines the semantics, in both
directions at once. Cycle 16 fixed reconstruction *fidelity* and immediately exposed that
fidelity was never the whole problem — K1/K2 are reconstructed perfectly *and* wrong,
because correctness requires container context the guard refuses to track, and K5/K6 show
even the detector can't be made container-proof. The rule that closes the class by
construction: **the guard executes a command only when it is published standalone** — a
physical line whose entire raw text matches `SUPPORTED_COMMAND`, or an inline code span
whose entire content matches — and *refuses* everything else that smells like a command,
using a suspicion detector deliberately broader than the execution one (e.g. a line or
backslash-joined group containing both `node` and `--test` as substrings: false positives
become named refusals, which is the fail-closed direction the guard already chose).
Execution bytes then always equal publication bytes; there is no reconstruction left to be
wrong. Cost: P19's blockquoted-fence form and the fixture's `$`-prompt form become
refusals — no shipped file uses either; the two real published commands are already
standalone. Sixteen cycles of edge cases is itself the evidence for this narrowing.

### [NIT] Comments assert the premises this round refutes

Line 226 ("if a shipped file prints a `node --test` command anywhere, in any context, that
command should work"), lines 249-251 ("leaves what the reader actually copies … matches the
shell rather than approximating it in either direction"), and line 274 ("A command never
begins with `>` or `$ `") are each measured false by K1–K7. When the code is fixed, sweep
the comments with it — the cycle-15 comment made the same kind of overclaim, and it cost a
round.

### Trade-off judgment (asked for by the brief)

The named trade-off — a *legitimate* command split across lines is refused rather than run —
I **accept**. No shipped file uses the form, the refusal is loud and named, and the
alternative is the sixteenth revision of a hand-written shell+markdown evaluator. What I do
not accept is the current *unstated* companion: a split command the detector cannot see is
not refused but silently skipped (K5/K6). Refusal must be the floor for anything suspicious,
which is what the MINOR's broader suspicion detector buys.

### Signoff: kimi-1 — 2026-07-30
Status: ❌ BLOCK
