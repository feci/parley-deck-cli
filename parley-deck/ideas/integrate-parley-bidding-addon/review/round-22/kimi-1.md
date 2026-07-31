---
idea: integrate-parley-bidding-addon
review-round: 22
agent: kimi-1
date: 2026-07-31
reviewed-commit: b1f43e4
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

**Closed, measured on this tree.** Round 21 I blocked on one arm: `resolutionTouchpoints`
accumulated its path by string join and discarded the physical landing `walkRawTarget` returns,
so a link reached through an earlier link had its relative target walked from the spelling, and
every `..` climbed a different tree than the kernel climbs. Cycle 25 implements exactly the
remedy I stated (`lib/installer.js:1495-1496` propagates `landed` into `logical`) plus codex-1's
companion point (`lib/installer.js:1420-1425` anchors a relative target at
`fs.realpathSync` of the link's parent). Reproduced both sides myself:

- At `381e639` (git archive to /tmp, `dist/` excluded, `node_modules` symlinked, the `b1f43e4`
  test file overlaid): `node --test --test-name-pattern="physically sits"` **fails** with
  `true !== false` at the `result.ok` assertion — the install escaped and reported success.
- At `b1f43e4`: the same test **passes**; the backslash pin passes at both commits, exactly as
  its body declares (verified at `381e639` too — the labelling claim is true).
- The archive staging dir was removed after the run.

No other finding of mine from any round stands open. The single-writer limit stays settled per
the round-14 ruling, verbatim in `CHANGELOG.md` under "Known limits".

## Position on the gate

**1 (correct).** I move from round 21's position 3 because the arm I named is closed by the
remedy I stated, and I measured the closure at both commits. Round 21's split closes on the
merits of the fix: codex-1 and I examined the anchor and found it wrong; the anchor is now
physical end to end. `hermes-1` arrived at 1 by examining only the expansion — the right answer
for an incomplete reason; the record should not read round 21's split as having validated that
reading.

On position 2: I was the reviewer who originally leaned toward narrowing. I accept codex-1's
recorded argument, and I verified its premise on this tree — `CHANGELOG.md:49-54` promises
fleet-wide atomicity ("a predictable failure anywhere produces zero writes and zero deletions")
with no exclusion for symlinked runtime homes, and symlinked homes were deliberately supported
by earlier cycles. Narrowing now would withdraw a documented promise, not trim an untouched
edge. If the project ever narrows, the promise must be rewritten in the same change.

My own attempt to find a third position-3 arm, three probes against the live gate (script run
from /tmp, homes deleted afterwards):

- **Dangling intermediate link, both destinations fresh** (`C/skills -> ../K/skills`, `K`
  absent): blocked in preflight with the destination parent named, every other unit skipped,
  **0 writes**, tree afterwards contained only the planted link.
- **Link landing exactly on the other unit's existing destination** (`C/skills ->
  ../K/skills/parley-deck` after seeding kimi): refused — "one contains the other", 0 writes.
  The `statSync`-follows-links anchor in `identityChain` puts kimi's inode inside codex's chain.
- **Valid link onto the other unit's skills dir, both destinations fresh**: refused — "they
  resolve to the same directory", 0 writes.

The three existence regimes line up with the three layers of the gate: what resolves is caught
by identity equality or chain nesting, what crosses is caught by the touchpoint walk, and what
cannot resolve cannot be written through and dies in preflight. I found no input that reaches
phase 2 with two units aliasing one physical directory.

## New findings

None.

## Release judgement

Releasable as 2.1.0. The payload is unchanged since `714712f` (empty diff under `skills/` and
`bin/` against this tree), the manifest hash is unchanged, and the gate — the only thing that
has moved for sixteen cycles — now measures correct at both ends of the round-21 arm and
refuses every regime I could construct. In round 21 I said I found nothing else standing
between this tree and release; with the named arm measured closed, nothing is.

## What I verified

- `git rev-parse HEAD` = `b1f43e4d8fec…`, working tree clean.
- Diff scope: `git diff 49fc3ec..b1f43e4 --stat` touches only `CHANGELOG.md`,
  `lib/addon-manifest.js`, `lib/installer.js`, `scripts/run-python-tests.js`, and the two test
  files — the gate and its records, as claimed. `git diff 714712f..b1f43e4 -- skills/ bin/` is
  empty.
- Read the full cycle-25 diff and the gate itself
  (`lib/installer.js:1330-1557`: `splitAtRoot`, `identityChain`, `entryChain`, `walkRawTarget`,
  `resolutionTouchpoints`, `aliasedDestinations`).
- `node --test` at `b1f43e4`: **368 pass, 0 fail** (~50s).
- Python leg: **54/54 across 7 files on python3 3.14.6**; with `PATH="/usr/bin:$PATH"` the
  runner prints `python3 is 3.9, but the add-on declares >=3.10` and exits **1** — fail-closed
  as designed.
- Manifest: `node scripts/build-addon-manifest.js --check` → **ok, 47 files,
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`**.
- Arm reproduction at `381e639` and confirmation at `b1f43e4`, per the table above; the pin
  test's "passes at 381e639 too" label independently confirmed.
- The three probes under "Position on the gate" (this round's own arm hunt).
- `CHANGELOG.md:106-112` carries the single-writer known-limit verbatim; the atomicity promise
  at `:49-54` contains no symlink exclusion.
- `IMPLEMENTATION.md` cycle-25 record matches every measurement I made independently, including
  the discrimination table and the pin's labelling. `review/round-09/VOID.md` exists.
- `test/design-addons.test.js` in the full-scope diff: comment corrections plus three added
  assertions around the backslash-sentinel arm; no assertion removed or weakened.
- "All seven accumulated arms refused" is covered by the full-suite pass; I did not re-enumerate
  them individually.
- Disk hygiene: the `381e639` archive copy, the probe script, and all probe homes are deleted.
  Older `/tmp/parley-*` and `/tmp/kimi1-*` leftovers from earlier rounds were left in place —
  round 22 is still open for the other two reviewers and I cannot attribute those safely. The
  volume sits at 63% used, 705 GiB free.
