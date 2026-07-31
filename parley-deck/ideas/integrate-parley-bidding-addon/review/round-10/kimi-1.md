---
idea: integrate-parley-bidding-addon
review-round: 10
agent: kimi-1
date: 2026-07-31
reviewed-commit: 9ed2081
---

## Verdict
BLOCK

One MAJOR. The fifth door is real, and it is the one the last four cycles never looked at:
every fix preflighted whether a destination can be **created** (ancestors: existence, type,
write-and-search) and whether it may be **replaced** (ownership). Nothing preflights whether an
**existing destination tree can be removed** — and both mutation paths depend on exactly that:
the install replace path renames the old tree to a backup and `rmSync`s it
(`lib/installer.js:1356-1362`), and uninstall `rmSync`s the tree outright
(`lib/installer.js:1298`). Any directory inside the destination that this process cannot
write — a `chmod -R a-w` tree, or a single mode-000 subdirectory anywhere in it — passes every
preflight and fails mid-fleet / mid-set. All arms below measured at `9ed2081`, uid 501 (not
root), node v26.5.0, macOS, in throwaway homes under `$TMPDIR`; the working tree was never
touched (verified clean before and after).

## Outstanding findings — closed or not

All four cycles' findings are verified **closed**, by re-measurement, not by reading:

- **Cycle 12 + 13 regressions discriminate exactly as claimed.** I extracted
  `lib/installer.js` from `dcd200e`, `3330a6e`, `9ed2081` into a `/tmp` copy of the
  `9ed2081` tree (via `git archive`; repo untouched) and ran the five new tests against
  each. `dcd200e`: the four refusal tests **fail**, and the force=true file-ancestor failure
  message is literally `wrote 78 units before failing` — the round-9 measurement reproduced.
  `3330a6e`: file-ancestor arm **passes**, permission arm **fails**. `9ed2081`: 5/5 pass,
  including the symlinked-home guard. The middle column is real: cycle 12's own regression
  cannot see cycle 13's arm.
- **Cycle 10 fleet preflight**: holds, including under `--force` — probe D below produced
  zero writes with an impossible destination in the last of 14 targets.
- **Cycle 11 dangling-destination doctrine**: holds on the doors I probed. A dangling
  symlink as an *intermediate ancestor* (`~/.aionrs/skills` → nonexistent) under `--force`
  is `blocked` with an accurate message and **zero** writes. A dangling symlink at `dest`
  itself: unforced install `blocked` with the link untouched; `--force` replaces the link
  only. The `stat`-resolves-`lstat`-locates split in `destinationAncestorObstacle` is
  correct on every arm I ran.
- **No write-through or delete-through symlinks** (the two safety properties I cared most
  about verifying myself): with `dest` a symlink to a foreign directory containing a marker
  file, forced install replaced the *link* and left the foreign directory's contents intact;
  forced uninstall removed the *link* only (`rmSync` does not traverse). Both measured.

Recorded follow-ups stay follow-ups:

- **The round-9 kimi-1 NIT (discovery guard `dirExists`) — I agree with the deferral; do
  not promote it to cycle 14.** Re-measured: after a `--no-addons` install, a dangling
  symlink at the unselected `parley-bidding` path is invisible to unflagged `doctor`
  (`ok:true`, only `parley-deck:valid` listed) while a real directory there is reported
  `malformed` (`ok:false`). The asymmetry is exactly as described, and the mutation paths
  are coherent (`install --only parley-bidding` is `blocked` with the accurate ownership
  message; `--force` remediates). Nothing runnable occupies a dangling link, so the green
  is not a false claim about anything installed — unlike the MAJOR below, nothing here
  writes or deletes. It changes discovery semantics rounds 4 and 6 ratified; the follow-up
  list is the right venue.
- **Only `parley-bidding` ships a manifest** (one `valid-unmanaged`, five `malformed` on a
  universal install): unchanged, B3.11 holds, correctly the first follow-up.

## New findings

### MAJOR — destination-tree removability is never preflighted; a predictable removal failure writes/removes the rest of the fleet first

B5 as ratified: "Preflight every unit and destination **before the first write**; a
predictable failure … must produce **zero** writes" (FINAL.md:61-63), and round 4's ratified
uninstall twin: "a refusal on the last unit must not leave earlier ones already deleted."
Four measured arms, one root cause:

- **Arm 1 — install, no `--force`, owned tree made read-only in the last target.**
  Installed `aionrs` normally (all units owned, healthy), then `chmod -R a-w` on
  `~/.aionrs/skills/parley-deck` (dirs 0555, files 0444 — marker still readable, so
  ownership preflight passes without any flag). Ran
  `install --target all --include-undetected`. Measured: **83 units written** across the
  fleet, then `aionrs/parley-deck` reports `action:"failed"` — `ENOTEMPTY` from
  `rmSync(backup)` (node's recursive rm empties bottom-up; it cannot unlink inside 0555
  directories). The rename-based replace had already **succeeded**: the unit's marker
  timestamps differ before/after, and `doctor` run immediately afterwards reports all six
  `aionrs` units `valid`, `ok:true`. So: install says failed, disk says installed, doctor
  says healthy, and a `.parley-deck.<pid>.<ts>.bak` debris tree (the entire old payload) is
  left in the skills directory, invisible to health. Same "78+ writes before a predictable
  failure" shape as cycles 12/13 — one door further along.
- **Arm 2 — same, but only ONE mode-000 subdirectory** (`references/`) inside an otherwise
  fully writable owned tree. Identical outcome: 83 writes, `failed`, `.bak` debris. This
  kills any "check the destination root's mode" fix — the preflight must walk the whole
  destination tree.
- **Arm 3 — uninstall, no `--force`, one owned add-on tree read-only.** Installed `aionrs`,
  `chmod -R a-w` on `parley-bidding` only, then plain `uninstall --target aionrs`.
  Ownership preflight passes (everything owned). Measured: **core and four add-ons
  removed**, then `parley-bidding` `failed` (same `ENOTEMPTY`) and survives **byte-valid**.
  That is precisely the end state round 4's MAJOR measured and forbad — "removed the core
  and then refused the add-on" — reproduced at `9ed2081` with no flag and no exotic setup
  beyond one chmodded directory. Afterwards `doctor` reports the core and four add-ons
  `missing` with `parley-bidding` still `valid`: the honest epitaph for a partial removal
  the contract said could not happen.
- **Arm 4 — uninstall `--force` fleet, a later target's skills dir 0555.** Codex tree fully
  removed, all six `aionrs` units `failed`, all six trees still on disk — a cross-target
  partial uninstall. Same root cause: `uninstallTarget` preflights ownership only, and
  nothing at all under `--force`; nobody checks that the removal can succeed.

Also measured: arm 3 reproduces identically against `49fc3ec`'s `lib/installer.js`
(`uninstallOk:false`, core removed, add-on `failed`) — this door **predates** cycles 10–13.
It is not a regression those fixes introduced; it is the door they never considered. The
fixes asked "can `dest` be created?" (`destinationAncestorObstacle`, ancestors only) and
"may we touch it?" (ownership). The replace and remove paths additionally need "can the
existing tree at `dest` be disposed of?" — and that question is not asked anywhere.

The failure is 100% predictable at preflight by a read-only walk: every directory in an
existing destination tree needs `R_OK|W_OK|X_OK` from this process (the mirror of
`firstCopyObstacle`, over the destination instead of the source), plus `W_OK|X_OK` on
`path.dirname(dest)` for the removal paths. Realism is the same class the last two cycles
were judged MAJOR on (a chmodded directory), and softer than it sounds: 0555 trees are what
`cp -Rp` from a read-only store or a user's `chmod -R a-w` "freeze" produce, and arms 1–3
need **no** `--force`. Regressions for each arm must discriminate against the previous
commit — this idea has already shipped one test that passed before and after.

### NIT (non-blocking) — the "Python leg 54/54 under 3.9.6" wording

IMPLEMENTATION.md (and the round brief) say the suite was measured "under Homebrew python3
3.14.6 and again under `/usr/bin/python3` 3.9.6. Python leg **54/54**." The python leg cannot
run green under 3.9.6: `scripts/run-python-tests.js` hard-fails below the declared `>=3.10`
floor **by design** — measured: `python tests: python3 is 3.9, but the add-on declares
>=3.10`. What genuinely passes under a 3.9.6-first PATH is the **node** leg (325/325,
re-measured). The facts underneath are all sound; the sentence just reads as if both legs
ran green under both interpreters. One clarifying clause, given this idea's history of
claims that claimed more than they showed (D-2, cycle 10's "every destination path").

## Release judgement

Not releasable as 2.1.0. The one thing that must change: preflight **removability of the
existing destination tree** on both mutation paths — install's replace (backup-cleanup)
path and uninstall — so that a predictable removal failure produces zero writes/removals,
with regressions for arms 1–4 above, each confirmed to fail against this commit's
`lib/installer.js`. Everything else I measured at `9ed2081` held.

## What I verified

- **Tree discipline**: `git status --porcelain` clean before and after all work; HEAD stayed
  `9ed2081ff241f36a0b5b96e930be116327ac6fdc`. All older-commit comparisons ran in a
  `git archive` copy under `/tmp`; all probes ran in `$TMPDIR` homes. No edits under
  `skills/parley-bidding/`, no tree mutation of any kind.
- **Suites at `9ed2081`**: `npm test` — 325 node tests, 0 fail; python leg 54/54 across the
  seven files under python 3.14.6; manifest check ok (47 files,
  `sha256:7854adf1…b95a6d`, matching the claimed aggregate). Node leg re-run under a
  3.9.6-first PATH: 325/325. Python leg under 3.9.6 refuses by design (see NIT).
- **Regression matrix** (cycles 12/13): fail/fail at `dcd200e` (with the round-9
  `wrote 78 units` measurement reproduced verbatim), pass/fail at `3330a6e`, pass/pass at
  `9ed2081` — each arm against the previous commit's actual `lib/installer.js`.
- **Code read in full**: `lib/installer.js` (2043 lines), `lib/addon-manifest.js`, the
  cycles 8–13 diffs (`49fc3ec..9ed2081`), the new tests, IMPLEMENTATION.md cycles 10–13,
  `review/round-09/VOID.md`, FINAL.md's B5 text. The `destinationAncestorObstacle` walk
  (`lib/installer.js:1820-1845`) is correct on every arm I could construct: regular-file
  ancestor, mode-000 ancestor, dangling-symlink intermediate, `ENOTDIR` propagation, and
  the symlinked-home happy path.
- **Safety properties probed and held**: dangling intermediate under `--force` → blocked,
  zero writes; symlinked `dest` → unforced blocked with link untouched, forced replaces the
  link only, foreign directory's contents intact; forced uninstall on a symlinked `dest`
  removes the link only. Control probes (fully writable trees) green, no debris.
- **The MAJOR above**: all four arms measured at `9ed2081` (numbers quoted are from my
  probe runs, not from the record), arm 3 additionally at `49fc3ec`; the post-failure
  `doctor` states (all-valid after a "failed" install; missing-core after a partial
  uninstall) measured with a real PATH so the runtime probe resolved python3.
- **The round-9 deferred NIT**: re-measured as described; deferral endorsed.
- **Notes-for-reviewers items**: `test/design-addons.test.js` since round 7 is 22
  insertions / 0 deletions — additive only, no assertion weakened; the self-consistent-swap
  boundary is covered by tests (`test/bidding-addon.test.js:284`, `:844` — `verifyPayload`
  passes by construction, the marker check catches); the legacy `markerSchema` path reads
  as ratified.
