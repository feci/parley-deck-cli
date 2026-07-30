# Round 20 — kimi-1

## Summary

All three round-19 dispositions verify. Cycle 24 closed codex-1's substitution-builds-the-
binary finding (P58/P59 refused), cycle 25 closed my prose-arm block exactly where it was open
(N1–N5, N8 refused; the span forms were confirmed already closed at cycle 24 before the fix was
written), and hermes-1's arithmetic NIT is corrected — the new breakdown sums. The guard still
verifies real commands: P8 RUNS and fails, the three valid published forms stay green, the two
shipped `Verify:` forms pass. `commonmark` is pinned at 0.31.2 in `package.json`, the lockfile,
and the installed tree.

One new finding, in the arm cycle 25 added: the prose substitution rule is occurrence-level in
intent but **line-level in implementation**. A single properly published span anywhere on a
visible line sets `publishedWhole` and silences the substitution arm for the entire line — so
the exact N8 shape this cycle was written to refuse goes invisible when a valid `Verify:`-style
span shares its line, including via an ordinary soft wrap. Reproduced end-to-end below.
Everything else I probed held, including the boundary I stated in round 19.

## Round-19 dispositions

### hermes-1 ✅ ACCEPT, [NIT] probe-breakdown arithmetic — VERIFIED CORRECTED

Round 19's claim summed to 56 (46 refused + 4 valid + 4 invisible + 1 hard break + 1
runs-and-fails; P57 uncounted). The round-20 claim is 57 refused, 10 green (3 valid published
forms, 1 prose mention, **5** invisible forms, 1 hard break), 1 that RUNS and fails. My count of
the harness run, by verdict: 35 REFUSED + 22 REFUSED-PROSE = 57 refused; GREEN = P4, P7, P19
(valid), P26 (prose mention), P41, P42, P43, P56, **P57** (invisible), P45 (hard break) = 10;
RAN-AND-FAILED = P8 = 1. 57 + 10 + 1 = 68. The breakdown is now internally consistent and names
P57 among the invisible forms. Closed.

### codex-1 — no signoff (tool outage); measured finding — VERIFIED CLOSED by cycle 24

The finding as logged (`n$(printf '')ode --test no/such/dir`, `n${PATH#"$PATH"}ode --test
no/such/dir` — guard green, reader exit 1) is real and is now refused:

```
P58 subst-builds-binary      pass=11 fail=1  REFUSED
P59 expansion-builds-binary  pass=11 fail=1  REFUSED
```

The cycle-24 change to `mentionsATestCommand` (one recognisable half plus a word-building
construct is enough) is present, and the grammar assertions cover both forms. I also confirmed
the cycle-24 commit's stated scope limit: detection was deliberately NOT widened to "any unit
carrying `$`/backtick", so legitimate shell examples (`$HOME` et al.) keep passing. That scope
line matters for the boundary note below.

### kimi-1 ❌ BLOCK — prose pass matched raw characters — VERIFIED CLOSED by cycle 25

```
N1 prose-escaped-flag        pass=11 fail=1  REFUSED-PROSE
N2 prose-escaped-binary      pass=11 fail=1  REFUSED-PROSE
N3 prose-quote-splice        pass=11 fail=1  REFUSED-PROSE
N4 prose-uppercase-flag      pass=11 fail=1  REFUSED-PROSE
N5 prose-escaped-valid       pass=11 fail=1  REFUSED-PROSE
N6 span-subst-builds-node    pass=11 fail=1  REFUSED
N7 span-subst-inside-node    pass=11 fail=1  REFUSED
N8 prose-subst-builds-node   pass=11 fail=1  REFUSED-PROSE
N9 span-subst-valid-target   pass=11 fail=1  REFUSED
```

The briefing's claim that only the prose arm was open is itself verified. I extracted
`publishedTestCommands` from commit `4ac913e` (cycle 24) and ran the round-19 inputs through
that exact code:

```
N6 span-subst-builds-node => [{"command":"$(echo n)ode --test no/such/dir","origin":"code"}]
N7 span-subst-inside-node => [{"command":"n$(echo o)de --test no/such/dir","origin":"code"}]
N9 span-subst-valid-target => [{"command":"$(echo n)ode --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]
N1 prose-escaped-flag => []
N8 prose-subst-builds-node => []
```

Span forms: code occurrences at cycle 24, refused by the grammar — already closed. Prose forms:
invisible at cycle 24 — the only open arm, and the only thing the cycle-25 diff changes
(word view + `toRaw` provenance map in the prose pass, case-insensitive flag, prose
substitution arm). The fix matches the finding's shape precisely. Closed.

## The guard still verifies something

- Baseline (no probe file): `pass=12 fail=0`.
- Full suite: `npm test` → `pass 253 / fail 0` at `4fdc7c8`.
- P8 (`node --test no/such/dir` in a span) still goes through `/bin/sh` and fails the build with
  `published command failed` — the guard executes and checks real output, not just form.
- Valid published forms stay green: P4 (inline span), P7 (double-backtick span), P19
  (blockquoted fence), all `pass=12 fail=0`.

## False-positive surface

- Both shipped forms pass (covered by the green baseline):
  `skills/parley-tracker/templates/subtask.md:68` (`Verify: ` before a span) and `:74`
  (checklist item, span in parentheses). The occurrence-level rule keeps their code provenance.
- Control probe K4 — `Verify: \`node --test "skills/parley-tracker/bin/*.test.js"\` (about $2 of
  compute).` — is GREEN: a stray `$` in prose next to a properly published command is not
  refused. This is the case the `publishedWhole` gate legitimately exists for, and it works.

## Fourth-layer hunt

Probed and dismissed, consistent with the boundary I stated in round 19 and the author's
cycle-24 scope line:

- **Inline span literals cannot carry `\n`.** Verified against the pinned parser: a soft-wrapped
  span's literal is `"node --test foo/bar.test.js"` (newline → space), so no multi-line literal
  can desync the per-line owner map. `commonmark` 0.31.2 confirmed in `package.json`,
  `package-lock.json`, and `node_modules`.
- **Both-halves-substituted (K6: `Run $(echo n)ode --t$(echo e)st no/such/dir now.`)** — GREEN,
  and I judge it OUTSIDE the contract, not a finding: neither half is recognisable, and catching
  it would mean refusing any line carrying a substitution, which cycle 24 explicitly rejected to
  protect legitimate shell examples. Detection has to anchor on something; one recognisable half
  is the ratified anchor. Stated so the edge is on record.
- `nodejs --test x`, homoglyph/zero-width flags, glob flags (`--tes[t]`): unchanged from my
  round-19 boundary — they do not run as the canonical command on the reader's machine as
  published, or are a different binary name the contract does not claim.

## Findings

### [MAJOR] The prose substitution arm is suppressed line-wide: a valid span on the same line makes the N8 shape invisible

Cycle 25 added the prose substitution arm with a gate:

```js
const publishedWhole = [...segments.values()].some(mentionsATestCommand);
...
if (!publishedWhole && /[`$]/.test(line) && mentionsATestCommand(line)) {
  record(line.trim(), "prose");
}
```

`publishedWhole` is computed over the WHOLE visible line. One properly published span anywhere
on that line sets it, and the arm then stays silent for the entire line — including a second,
substitution-built command sitting in the prose of the same line. The word loop cannot see that
second command: its binary name is not spelled out (that is the point of the N8 shape). Nothing
else looks. The command is skipped — and skipping reads as success.

Reproductions (added to the supplied harness as K-probes, run against my worktree at `4fdc7c8`):

```
K1 shared-line-subst-binary  pass=12 fail=0  GREEN      ← guard silent
K2 shared-line-subst-flag    pass=12 fail=0  GREEN      ← guard silent
K3 shared-line-spelled-out   pass=11 fail=1  REFUSED-PROSE   (control)
K4 dollar-in-prose-control   pass=12 fail=0  GREEN      (control)
K5 shared-line-softwrapped   pass=12 fail=0  GREEN      ← guard silent
```

K1's document:

```
Verify: `node --test "skills/parley-tracker/bin/*.test.js"` — or run n$(printf '')ode --test no/such/dir instead.
```

The reader who copies the second command gets a real run of the real command — demonstrated,
not assumed:

```
$ sh -c 'set -- n$(printf "")ode --test no/such/dir; printf "<%s>\n" "$1" "$2" "$3"'
<node> <--test> <no/such/dir>
$ sh -c 'set -- node --t$(echo e)st x; printf "<%s>\n" "$2"'
<--test>
$ sh -c 'n$(printf "")ode --version'
v26.5.0
```

`node --test no/such/dir` exits 1 for that reader while the guard reports 12/12 green. K5 shows
the author does not even need one source line: a soft-wrapped paragraph merges into one visible
line (softbreak = space), so ordinary prose wrapping triggers the same suppression.

Mechanism, isolated with the current extractor (parse-only): on K1 the guard extracts exactly
one occurrence — the valid span. Remove ONLY the `!publishedWhole &&` conjunct from a copy and
the same input yields the prose occurrence too:

```
K1 current guard  => [{"command":"node --test \"skills/parley-tracker/bin/*.test.js\"","origin":"code"}]
K1 unsuppressed   => [{...same code...},{"command":"Verify: node --test ... — or run n$(printf '')ode --test no/such/dir instead.","origin":"prose"}]
```

So detection reaches the command; the line-wide gate alone hides it.

The controls bound the finding precisely. K3 (same line, second command spelled out) IS refused
— the word loop is occurrence-level and does its provenance check per match. K4 (stray `$`, no
second command) stays green — the gate has a real purpose and I am not asking for its removal.
The asymmetry is specific to the substitution arm: the one detector that cannot see spelled-out
tokens is the one gated line-wide.

Why this blocks rather than follows up. The decision rule this guard ratified in round 17 is
occurrence-level — "the question is not which bucket a fragment landed in", and P36 established
that a valid code occurrence must not mask an invalid prose one. This gate is a per-LINE boolean:
a bucket test, reintroduced in the arm written one cycle ago, and it masks exactly the prose
occurrence class (substitution-built, no spelled half) that cycles 24–25 were written to refuse.
N8 alone is refused; N8 next to a correct `Verify:` span passes. A correct occurrence masking a
broken one is the P36 failure in a new guise, and the failure mode is the one this guard exists
to prevent: a published command fails for the reader while the guard reports success. Every
previous instance of a demonstrated skip in a ratified class has been fixed before signoff; the
contrivance threshold was set by N8/P58 themselves, and K1/K2/K5 sit inside that ratified line,
not beyond it (contrast K6, which I explicitly do not raise).

Fix direction (direction, not prescription): make the gate occurrence-grained rather than
line-grained — evaluate the substitution arm against the line's text with published-whole code
segments excluded (owner-wise subtraction), so K4's `$2` stays green while K1's prose command is
recorded and refused by provenance. The data to do this is already in `segments`/`owners`.

### Boundary on record (not a finding)

K6 (`$(echo n)ode --t$(echo e)st no/such/dir`, neither half recognisable anywhere) is GREEN and,
per the cycle-24 scope line and my round-19 boundary, correctly so: closing it means refusing
any line that carries a substitution, which refuses legitimate shell documentation. If the team
ever wants that horizon moved, it is a contract change, not a bug fix.

## What I did NOT find

No asymmetry remains in the word-view pass itself (K3); no rendering, container, splice,
continuation, case, quote, comment, CDATA, script/style, alt-text, soft/hard-break, or
provenance gap beyond the one finding above; no false refusal of anything the project ships or
legitimately needs (K4, baseline, P26).

No tracked files were modified. I extended only my own untracked probe script
(`scratchpad/probe-kimi.sh`, six appended `run_probe` lines) and wrote scratch files under
`/tmp`; `git status` in the worktree shows only the pre-existing untracked `node_modules` link.

### Signoff: kimi-1 — 2026-07-30
Status: ❌ BLOCK
