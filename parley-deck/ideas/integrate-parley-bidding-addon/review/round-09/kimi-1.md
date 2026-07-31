---
idea: integrate-parley-bidding-addon
review-round: 9
agent: kimi-1
date: 2026-07-31
reviewed-commit: dcd200e
---

## Verdict
BLOCK

## Outstanding findings — closed or not

codex-1's round-8 MAJOR (F26) is **closed** at `dcd200e`, both halves, measured rather than
credited:

- **Fleet-wide preflight.** `installCommand` builds the full target × unit plan and preflights
  all of it before the first write (`lib/installer.js:618-654`). The regression — an unmarked
  `parley-bidding` destination in the fourteenth target under
  `--target all --include-undetected` — returns every unit `blocked`/`skipped` with zero
  targets written. I re-ran it against an exact `dcd200e` copy: pass. The same four
  regressions fail 4/4 against `ebe269e`'s installer, split 2-pass/2-fail against `3553f47`'s
  (the cycle-10 pair passes, the cycle-11 pair fails on `actual: 'missing'` and on the
  surviving link), and pass 4/4 at `dcd200e` — so the tests do discriminate, and the cycle-11
  record's fail-before/pass-after claim is accurate.
- **`pathEntryExists` (lstat, only `ENOENT` is absence).** A dangling destination symlink is
  now `blocked` in install preflight with zero writes, `malformed` (not `missing`) in health,
  `blocked` in an unforced uninstall, and actually removed by a forced one. I additionally
  walked the remediation journey the blocked message recommends, which no test covers:
  unforced install over a dangling link at an add-on dest → `blocked` with the `--force`
  message; `--force` → `replaced`, a real directory on disk, no temp/backup leftovers, doctor
  `valid`. Same at the core dest — round 8's `ENOTDIR` case is now a clean forced replace.

The ordinary success paths survive the preflight rewrite: fresh install → `installed`,
re-install → `replaced`, doctor all-`valid` across the six skills, uninstall removes the full
set, home left clean. CLI-level smoke of a fleet-blocked install exits 1 with per-unit
dispositions. The `selected`/recorded-selection logic, the manifest/marker anchor, and the
`valid-unmanaged` boundary are untouched by cycles 10–11 and remain covered by the suite I
re-ran in full.

The deferred follow-up is unchanged: only `parley-bidding` ships a manifest, so a universal
foreign copy reports one `valid-unmanaged` and five `malformed`. `FINAL.md` B3.11 holds the
other add-ons unaffected; still a recorded follow-up, not a 2.1.0 blocker.

## New findings

### [MAJOR] `--force` suppresses the only destination check, so an impossible destination still writes a partial fleet (B5)

**Where:** `lib/installer.js:996-1036` at `dcd200e` (`preflightSkillUnit`), reached through
`:618-654` (fleet preflight) and `:1122-1145` (per-target preflight).

**What:** Cycle 10's fleet-wide preflight checks each destination exactly once — the ownership
check, `pathEntryExists(dest) && !installerOwnsDestination(dest, skill) &&
!context.options.force`. That check is gated on `!force`, and there is no other
destination-side check. So under `--force`, whether the destination can exist at all is never
examined: not in the fleet preflight, not in the per-target preflight, not before earlier
targets are written.

**Measured at `dcd200e`** (exact copy of the reviewed tree, installer and tests from
`git show dcd200e:...`):

1. `~/.aionrs/skills` a **regular file**, `install --target all --include-undetected
   --force` → `ok:false`, `aionrs/parley-deck` `failed` (`EEXIST … mkdir
   '.aionrs/skills'`) — **after 13 targets were written on disk**. Without `--force` the same
   layout is `blocked` at preflight with zero writes: the ownership check happens to look at
   the path, and `pathEntryExists` treats the `ENOTDIR` as presence. `--force` removes exactly
   the check that was catching it.
2. `~/.aionrs/skills` a directory with **mode 000**, same command → `ok:false`, `aionrs`
   `failed` (`EACCES` on the staging temp dir) — again **13 targets written**. Same root
   cause, same force-specificity: unforced it is `blocked` with zero writes.

**Why it matters:** B5 requires every unit *and destination* to be preflighted before the
first write and zero writes for a predictable failure. An uncreatable destination is
statically discoverable — a walk to the nearest existing ancestor answers it without touching
anything. Round 8 closed this for ownership blockers; the `--force` path, which is precisely
the path an operator takes when destinations are unusual, kept the hole. A routine
`install --force` upgrade — the documented repair command — can therefore leave a
mixed-version fleet: thirteen runtimes replaced, the fourteenth failed mid-flight.

**Context for the fix-up:** cycle 12 (`3330a6e`) landed while this review was running and is
**not** the commit I reviewed. Its `destinationAncestorObstacle` closes measured case 1 —
I re-ran it there: `blocked` with a clear "is not a directory" message, zero writes, force and
unforced. Measured case 2 (unreadable ancestor) still writes 13 targets before failing at
`3330a6e`: `statSync` on a mode-000 directory succeeds, so the walk sees a directory and stops.
The fix the next round reviews should cover the unwritable/unreadable-ancestor arm of the same
class, or explicitly record why it is out of scope. Either way `3330a6e` needs its own fresh
full-scope strict-gate round.

### [NIT — non-blocking] A dangling symlink at an *unselected* add-on path is invisible to unflagged `doctor`

The round-4 discovery guard in `targetSkillUnits` (`lib/installer.js:978`) uses `dirExists`,
which follows symlinks, so a dangling link named `parley-bidding` outside the recorded
selection is not reported, while a real directory there is. Measured: after a `--no-addons`
install, a dangling link at the bidding dest → unflagged `doctor` reports only `parley-deck`,
`ok:true`. This is not a health lie — nothing usable is installed at that path, which is the
fact the opt-out verification relies on — and the mutation path is coherent: `install --only
parley-bidding` is `blocked` with the accurate "destination exists" message, and `--force`
remediates cleanly (verified above). Cycle 11's doctrine ("only `ENOENT` is absence") applied
to the discovery guard would make the two leftover kinds symmetric. Implementer's discretion;
does not block.

## Release judgement

**No.** `dcd200e` is not releasable as 2.1.0. The one thing that must change: destination
*feasibility* — not only ownership — must be preflighted across the whole plan before the
first write, independently of `--force`, covering a non-directory ancestor and an
unwritable/unreadable one. Cycle 12 (`3330a6e`) is the in-flight attempt; it closes the
regular-file arm as measured above, and the unreadable-ancestor arm is still open there. The
commit that lands the complete fix needs the fresh full-scope round the strict gate requires.

## What I verified

- The working tree was clean at `dcd200e` when this review began and was never mutated by me
  (no resets, no checkouts, no edits; all old-tree comparison done in `/tmp` copies built with
  `rsync` plus `git show <rev>:<path>`). HEAD moved to `3330a6e` (cycle 12) mid-review; every
  `dcd200e` measurement below was taken either before the move — the 320-test count pins the
  suite runs to `dcd200e` code, since cycle 12 adds two tests — or against an exact `dcd200e`
  copy afterwards.
- Read in full: `lib/installer.js` (1993 lines), `lib/addon-manifest.js`, all three test
  files, the diffs `ebe269e..dcd200e` and `49fc3ec..dcd200e`, FINAL.md B3/B5/B3.11, the
  round-8 reviews, consensus.md, and IMPLEMENTATION.md cycles 10–11. Front matter
  `head-commit: dcd200e` matched `git rev-parse HEAD` at review time.
- Full node suite: **320 tests, 0 fail** — with Homebrew python3 3.14.6 on PATH, with
  `/usr/bin/python3` 3.9.6 on PATH, and again inside the exact `dcd200e` copy.
- Python leg: **54/54 across seven files** on 3.14; on 3.9.6 it refuses by design — exit 1,
  "python3 is 3.9, but the add-on declares >=3.10".
- Manifest check: ok, **47 files**, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d` — matching the
  record. Payload hygiene: 0 symlinks, 0 `__pycache__`/`.pyc`/`.pyo`/nested `.gitignore`,
  48 files on disk.
- Regression discrimination (isolated copies, installers swapped via `git show`): the four
  cycle-10/11 regressions fail 4/4 on `ebe269e`, split 2/2 on `3553f47` (cycle-10 pair passes,
  cycle-11 pair fails on `actual: 'missing'` and the surviving link), pass 4/4 on `dcd200e`.
- Success/refusal probes beyond the suite (fresh temp homes): install → reinstall → doctor →
  uninstall round trip clean; unforced install over a dangling link `blocked`, forced
  `replaced` with no temp/backup leftovers and doctor `valid` after, at both an add-on and the
  core destination; dangling selected add-on → `malformed`, `ok:false`; fleet-blocked CLI run
  exits 1 with per-unit dispositions.
- The MAJOR above was reproduced independently at `dcd200e` before I credited it: regular-file
  and mode-000 destination parents, force and unforced, written targets counted on disk.
  Cycle-12 behavior at `3330a6e` was measured the same way for the fix-up's benefit.
