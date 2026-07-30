# Round 21 — kimi-1

## Summary

All three round-20 findings verify closed, and hermes-1's NIT correction re-measures true in my
own shells. The claimed probe arithmetic reproduces exactly: **69 refused, 12 green, 1 that RUNS
and fails** out of 82 probes (P1–P59, N1–N9, R1–R5, K1–K5, H1–H4), baseline `pass=12 fail=0`,
full suite `npm test` 253/253 at `4cc45de`. Both corrections to the author's own work hold: the
word-building predicate lives in exactly one place used by both passes, and I ran both halves of
the cycle-26 reproduction rather than trusting it.

I then aimed fourteen supplementary probes (T-cases) at the three new code paths — the
per-code-node exclusion, the entity-newline join, and the brace predicate — plus their
combinations. Thirteen behave as intended. The fourteenth documents a real detection boundary
(glob-built command words) that is filesystem-gated, absent from this repo and its tarball, and
loud rather than silent where the precondition fails. It is a follow-up, not a block.

Nothing I found stops this from being published today.

## Round-20 dispositions

### codex-1 — brace expansion builds a word out of plain letters (cycle 26) — VERIFIED CLOSED

```
R1 brace-builds-binary       pass=11 fail=1  REFUSED
R2 brace-builds-flag         pass=11 fail=1  REFUSED
R3 prose-brace-binary        pass=11 fail=1  REFUSED-PROSE
R4 prose-brace-flag          pass=11 fail=1  REFUSED-PROSE
R5 canonical-brace-target    pass=12 fail=0  GREEN
```

Re-measured the expansion claim myself rather than citing it:

```
$ sh  -c 'printf "<%s>\n" n{o..o}de --test x'        →  <node> <--test> <x>
$ zsh -c 'printf "<%s>\n" n{o..o}de --test x'        →  <node> <--test> <x>
$ sh  -c 'printf "<%s>\n" node --{test..test} x'     →  <node> <--{test..test}> <x>
$ zsh -c 'printf "<%s>\n" node --{test..test} x'     →  <node> <--{test..test}> <x>
$ sh  -c 'printf "<%s>\n" "a{claim,validate}b"'      →  <a{claim,validate}b>
```

The binary arm really expands in both shells; the flag arm and the quoted form do not (this is
also hermes-1's NIT, confirmed below). R5's GREEN is correct and I verified its mechanism: the
quoted braces reach the shell literally, and **Node 26's own runner** expands the brace glob —
`sh -c 'node --test "skills/parley-tracker/bin/{claim,validate}.test.js"'` runs 35 tests, all
pass (both files exist). The diff comment "the executor decides whether the target works" is
accurate.

The shared-predicate correction holds: `buildsWords` is defined once
(`test/design-addons.test.js:278`) and is the only word-building test in the file — called by
`mentionsATestCommand` (`:299`), which serves the code pass via `addCode` (`:416`), and called
directly by the prose synthesis rule (`:485`). There is no second copy to drift.

### kimi-1 — a published span silenced the synthesis rule line-wide (cycle 27) — VERIFIED CLOSED

```
K1 shared-line-subst-binary  pass=11 fail=1  REFUSED-PROSE
K2 shared-line-subst-flag    pass=11 fail=1  REFUSED-PROSE
K3 shared-line-spelled-out   pass=11 fail=1  REFUSED-PROSE
K4 dollar-in-prose-control   pass=12 fail=0  GREEN
K5 shared-line-softwrapped   pass=11 fail=1  REFUSED-PROSE
```

The fix is the one I asked for and no more: publishing code nodes are excluded owner-wise from
the residue (`:464-470`), and the synthesis rule judges what remains (`:485`) — occurrence-level,
matching the round-17 contract. The fixture pins the shared-line case and the `$FOO` control
(`:888-898`). My own round-20 K-probes now all refuse where they must, and the K4 control that
legitimises the gate stays green.

I extended the same shape at the new constructs (see T-cases below): the exclusion does not
silence a brace-built command, a substitution sitting between TWO publishing spans, an
entity-spliced command, or a spelled-out command bracketed by two spans.

### hermes-1 — an entity newline splits the prose line (cycle 28) — VERIFIED CLOSED

```
H1 entity-newline            pass=11 fail=1  REFUSED-PROSE
H2 named-entity-newline      pass=11 fail=1  REFUSED-PROSE
H3 entity-in-html-block      pass=11 fail=1  REFUSED-PROSE
H4 literal-flag-no-expansion pass=11 fail=1  REFUSED
```

Parse-only confirmation of the premise (pinned `commonmark` 0.31.2 — exact pin in
`package.json`, the lockfile, and the installed tree):

```
"Run node&#10;--test x now."     → text "Run node", text "\n", text "--test x now."
"Run node&NewLine;--test x now." → text "Run node", text "\n", text "--test x now."
"Run node&#13;--test x now."     → text "Run node", text "\r", text "--test x now."
"Run node&#9;--test x now."      → text "Run node", text "\t", text "--test x now."
```

The entity forms produce a literal `\n` inside a text node — exactly what cycle 28 joins with a
space (`:535`), matching what HTML does with that newline (whitespace). The `\r` and `\t`
siblings never split `visible.split("\n")` either; T5/T6 below show both are still caught. I
also verified the premise's converse: a real line break is a `softbreak`/`linebreak` node,
never a `\n` inside a text literal, and a multi-line inline code span's literal has its newline
normalised to a space by the parser (`"node --test \\ x"`), so no text-node newline can be
anything but an entity.

### hermes-1 (NIT) — the flag-arm reproduction does not expand — CONFIRMED, DISPOSITION STANDS

Measured above: `--{test..test}` passes through literally in sh and zsh. The kept case is
refused for the reason that is true, and it is worth being precise about what that reason is:
the grammar anchors `^node\s+--test\s` literally (`:584`), so `--{test..test}` fails the anchor
on sight — R2/R4/H4 are REFUSED before any executor sees them, and a reader who pastes it gets
node's exit 9. The fixture comment now labels it as such (`:864-869`). The correction is in,
the case is honestly labelled, and the second half of the cycle-26 reproduction has now been
run by two reviewers independently.

## The guard still verifies something

- Baseline (no probe file): `pass=12 fail=0`. Full suite: `npm test` → `pass 253 / fail 0`.
- P8 (`node --test no/such/dir` in a proper span) still goes through `/bin/sh` and fails the
  build with `published command failed` — the one RAN-AND-FAILED. The guard executes real
  commands and reads real summaries, not just form.
- Valid published forms stay green: P4 (inline span), P7 (double-backtick span), P19
  (blockquoted fence), R5 (brace target the executor resolves).
- Probe arithmetic, counted from the run: 35 REFUSED + 34 REFUSED-PROSE = **69 refused**;
  GREEN = P4, P7, P19, P26, P41, P42, P43, P45, P56, P57, R5, K4 = **12**; RAN-AND-FAILED = P8
  = **1**. 69 + 12 + 1 = 82. The claim to be measured is measured.

## Supplementary probes (round-21 code paths)

Same mechanism as the supplied harness (one temp file under `skills/`, run the guard, report,
remove; my own untracked copy `scratchpad/probe-kimi-r21.sh`):

```
T1  span+brace-binary            pass=11 fail=1  REFUSED-PROSE
T2  span+brace-flag              pass=11 fail=1  REFUSED-PROSE
T3  two-spans+subst              pass=11 fail=1  REFUSED-PROSE
T4  span+entity-newline          pass=11 fail=1  REFUSED-PROSE
T15 spans-bracket-spelled        pass=11 fail=1  REFUSED-PROSE
T5  entity-CR                    pass=11 fail=1  REFUSED-PROSE
T6  entity-tab                   pass=11 fail=1  REFUSED-PROSE
T7  named-entity-html-block      pass=11 fail=1  REFUSED-PROSE
T8  entity-splits-flag-target    pass=11 fail=1  REFUSED-PROSE
T10 brace-comma-binary           pass=11 fail=1  REFUSED-PROSE
T11 brace-no-halves-control      pass=12 fail=0  GREEN
T12 span+brace-control           pass=12 fail=0  GREEN
T13 multiline-span-backslash     pass=11 fail=1  REFUSED
T14 glob-builds-binary           pass=12 fail=0  GREEN      ← see findings
T16 residue-overjoin-control     pass=12 fail=0  GREEN      (inline, not in the script)
```

What each group establishes:

- **T1–T4, T15** — the cycle-27 exclusion does not silence ANY construct class on a shared
  line: braces (T1/T2), a substitution between two publishing spans (T3), an entity-spliced
  command beside a valid span (T4), and a spelled-out command bracketed by two valid spans
  (T15). The K-cases covered substitution; these close the same shape for the two constructs
  added since.
- **T5–T8** — entity whitespace beyond `&#10;`: `\r` (T5) and `\t` (T6) never split the guard's
  line and are still caught; the named entity inside an HTML block (T7) and an entity placed
  between flag and target (T8 — the exact split point that would have made the command
  invisible pre-cycle-28) are refused.
- **T10–T12** — the comma form of the brace predicate (`n{od,od}e`) refuses; braces with
  neither half (T11) and braces in prose beside a valid span (T12) correctly do not fire — the
  predicate does not false-positive on ordinary prose.
- **T13** — a backslash before a soft-wrapped newline INSIDE an inline code span: the parser
  normalises the newline to a space, the literal keeps the `\`, the grammar refuses it.
- **T16** — the exclusion's owner-wise subtraction cannot be tricked into over-joining:
  `no` + valid span + `de --test no/such/dir` leaves a residue that spells `node --test
  no/such/dir`, but no reader can copy that (the rendered line reads `nonode --test …de
  --test …`), and the guard correctly stays silent rather than manufacturing a refusal. The
  join direction that matters — refuse — was already covered by T1–T4.

## False-positive surface

- Both shipped `Verify:` forms (`skills/parley-tracker/templates/subtask.md:68`, `:74`) pass
  with code provenance, as does the third published command
  (`skills/parley-design-check/SKILL.md:372`). The green baseline IS this check.
- Shipped-file sweep for the exotic classes: two hits, neither near a command —
  `skills/parley-design-check/SKILL.md:365` (a `{fail,pass}` brace in a prose span, no `node` /
  `--test` on the line) and `skills/parley-worktrees/SKILL.md:420` (a `$(comm …)` substitution
  inside a fenced shell example, no command words). Twenty-eight cycles in, no shipped file has
  ever used a refused construction, and that is still true.
- Controls K4, T11, T12, T16 all green: stray `$`, stray braces, and exclusion artefacts do not
  refuse legitimate documentation.

## Findings

### [NIT] A glob can build a command word, and detection has no construct for it — **follow-up**

T14: `Run [n]ode --test no/such/dir now.` is GREEN — invisible to both passes. `[n]ode` carries
no backtick, `$`, or brace (so `buildsWords` is false) and spells neither half literally (the
word view removes only `\'"`). The same holds for `node --tes[t] x`.

Why this does not block. A glob is not a word-building construct in the way the ratified three
are: its result depends on the reader's filesystem, not on the characters alone. `[n]ode`
becomes `node` only if a file literally named `node` exists in the directory where the command
is pasted (likewise `--tes[t]` needs a file named `--test`). I checked the precondition: no
such file at the repo root, and `npm pack --dry-run` ships none. Where the precondition fails —
everywhere this package installs — the reader gets a loud `command not found`, and where it
holds they get the command the page plainly shows. The guard itself can never execute a glob
form, because it never sees one: no occurrence is recorded, so nothing false is certified
green. That is the same shape as the K6 boundary I put on record in round 20 (both halves
substituted): detection must anchor on something, and extending the anchor to globs would
police the reader's filesystem, not the page. It belongs with the tracked follow-ups, next to
K6, as a contract question rather than a bug.

## What I did NOT find

No residual line-silencing in any construct class (T1–T4, T15); no entity-whitespace split that
survives (`\n`, `\r`, `\t`, named, HTML block, mid-command); no brace form that slips the
predicate while naming a half; no over-join or over-exclusion in the new residue logic (T16,
K4); no false refusal of anything the project ships; no drift between the passes (one
predicate); no change to the product under review, which has been verified since round 01.

No tracked files were modified. I wrote one untracked probe script
(`scratchpad/probe-kimi-r21.sh`, a copy of the supplied harness with my own `run_probe` cases
and the same temp-file discipline) and ran one inline probe; `git status` in the worktree shows
only the pre-existing untracked `node_modules` link. I did not read codex-1's or hermes-1's
round-21 reviews before filing this one.

### Signoff: kimi-1 — 2026-07-30
Status: ✅ ACCEPT
