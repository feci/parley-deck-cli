---
idea: integrate-parley-bidding-addon
review-round: 12
agent: kimi-1
date: 2026-07-31
reviewed-commit: 5100f34
---

## Verdict
BLOCK

One MAJOR: the disposal class cycle 15 removed by transaction is closed only on the uninstall
side. Install still predicts-then-commits per unit, and the deleted walk was — accidentally —
the last fleet-wide guard on that path. Measured: the round-8 forbidden end state (N targets
changed, failure discovered at target N+1) is reachable again at `5100f34`, and for the
`uchg`-destination arm this is a regression introduced by this cycle's own predicate deletion.

## Outstanding findings — closed or not

All four round-11 findings are **closed**, each re-measured by me at `5100f34` as uid 501
(rsync copies at `/tmp/pd-r12/head` and `/tmp/pd-r12/prev`, library-driven synthetic HOMEs,
probes in `/tmp/pd-r12/probes`; the working tree was never touched — `git status` clean
before and after):

- **CRITICAL — marker `addons` path traversal.** Closed at two independent layers, both
  probed. Layer 1: every recorded name must be a string, not `.`/`..`, separator-free,
  matching `^[A-Za-z0-9][A-Za-z0-9._-]*$`, unduplicated — else the selection is unusable, no
  unit is constructed, the core reports `malformed`, both mutation preflights refuse.
  Layer 2: whatever a name's origin (marker, flag, discovery), its destination must resolve
  to an exact direct child of the skills directory or no unit is built. My end-to-end arm
  through `run()`: marker `addons: ["../../outside-sentinel"]`, `uninstall --force --json` →
  exit 1, `ok:false`, sentinel intact, core intact; the suite additionally asserts no unit's
  `dest` even mentions the sentinel. I probed the edges the regex invites: `../escape`,
  `sub/dir`, `/absolute`, `.hidden`, `foo `, ` foo`, `..foo`, `parley-deck `, en-dash,
  duplicates, and `["parley-deck"]` (a duplicate *destination* in the plan — phase A rolls
  back, 0 removed, tree intact). All fail closed. I also enumerated every other stored-data→
  path route: marker fields `skill`/`version`/`manifest`/`source`/`target`/`scope` are only
  compared or displayed; `--only` names must match discovery; discovery names come from
  `readdir`. The `addons` array was the only route, and it is bricked.
- **MAJOR — the disposal predicate.** Deleted, not patched, and the transaction that replaces
  it is sound in every state I could put it in: phase-A rollback restores the fleet
  **byte-identically** with **zero `.removing` debris** (my hash-tree comparison, B3 — the
  suite only counts `removed:0`); the rollback-*failure* branch keeps the quarantined tree
  and names it in the unit's message, 0 removed (injected double fault, C1); the headline
  `uchg` arm measures exactly as claimed — **84 removed, 0 failed, 1 warning naming the
  debris** (C2); `uappnd` on a skills directory → full rollback, 0 removed (B2); six
  concurrent two-process runs → exactly one winner (84 removed), the loser fails clean, no
  half-deleted trees, no debris (C3); `--dry-run` changes nothing on disk (A6); a live
  symlinked destination is unlinked, never traversed — the target tree keeps its bytes (A5);
  install over a live symlink replaces the link with a real directory and preserves the
  target (A14). The "rename succeeds on exactly the trees whose recursive removal fails"
  premise re-measured at primitive level: `rename(2)` on a `uchg` directory fails EPERM on
  the root itself, succeeds over a tree containing flagged files.
- **MAJOR — legacy exemption scope.** Now keyed to whether the packaged source ships a
  manifest. The rewritten test honestly uses `parley-worktrees` (which shipped in 2.0.0), and
  the `parley-bidding` arm (both fields deleted + payload tampered) is `malformed`. Verified
  in the suite and by reading: a schema-2 marker still validates fully even for a unit whose
  source is absent from the package; only the schema-less exemption consults the source.
- **MINOR — `hashFile` outside `verifyPayload`'s try.** Now `unreadable (EACCES): <file>` in
  the problems list; suite regression present; diff matches the claim.

Discrimination claim, re-measured independently: the cycle-15 test file against `12f9071`'s
`lib/` fails **9 of the 11** diff-touched regressions; the 2 that pass at both are preserved
properties (the genuine legacy shape, and kept-manifest/lost-schema). The record says **3/11**
— see the NIT below. At `5100f34` all 11 pass; full gate: **338 node tests, 0 fail** (both
under python 3.14.6 and with `python3` shimmed to `/usr/bin/python3` 3.9.6), Python leg
**54/54 on 3.14**, manifest check ok (47 files, `sha256:7854adf1…b95a6d`).

Deferred items unchanged and correctly recorded: only `parley-bidding` ships a manifest
(B3.11); my round-9 `dirExists` discovery-guard NIT stands deferred. The recorded known limit
(phase-B debris invisible to `doctor`) is accurate — I left such debris on disk in C2 and no
later command surfaces it.

## New findings

### MAJOR — install still leaves a partial fleet when a commit-time operation fails mid-fleet; the `uchg`-destination arm is a cycle-15 regression

Cycle 15's answer to "removability cannot be predicted" is: transact. It was applied to the
only command that deletes — not to the command that replaces. Install remains
fleet-preflight + sequential per-unit commit, and `copyPayloadAtomically` contains **two
fallible pre-commit operations preflight cannot see**: the staging `mkdir` and the
`rename(dest → .bak)`. Measured at `5100f34`, fleet install `--force` against the last of
fourteen targets:

| arm | mechanism | measured at `5100f34` | at `12f9071` |
|---|---|---|---|
| `chflags uchg` on the destination **directory** | `rename(dest→.bak)` → EPERM | **83 units replaced, then the 84th fails** (B1) | **0 replaced — clean fleet-wide block** |
| `chflags uappnd` on the skills directory | staging `mkdir(.tmp)` → EPERM | **78 replaced, then 6 units fail** (B2) | same door open (pre-existing) |

The first row is the one that must be read carefully. macOS `access(2)` returns **EPERM** on
a flag-locked directory (primitive probe, re-measured), so cycle 14's
`firstRemovalObstacle` walk — wrong as a disposal oracle in both directions, correctly
deleted — was *incidentally* install's fleet guard against the most realistic flag idiom
there is: a locked destination folder (Finder "Locked" / `chflags uchg`). Cycle 15 removed
it from install preflight on the reasoning that "replacement commits by rename, and a failed
backup cleanup is only a warning." That reasoning covers obstacles **inside** the old tree
(locked file, frozen subtree — the commit renames touch only the root, cleanup warns, my
round-11 E4/E9 arm, still absorbed). It does not cover obstacles **on** the root itself —
the backup rename is not housekeeping, it is pre-commit — nor on the parent. IMPLEMENTATION.md's
"Install needed nothing further: it already commits by rename" is, measured, wrong for
exactly the arm the deleted walk used to catch.

The end state is the one this idea has gated on three times: round 8 (install, predictable
blocker found after 13 targets written — MAJOR), round 10 (uninstall, same shape — MAJOR),
round 11 (disposal misprediction — MAJOR). Harm is milder than uninstall's was — nothing is
lost, the failed unit keeps its intact old tree, the error names the operation, re-running
after `nouchg` converges — and the suite's own fleet-preflight tests still pass because
their blockers are predictable. But the deck's revealed posture across three rounds is that
a partial fleet on a mutation command gates the release when a transaction design exists,
and cycle 15 demonstrates one exists in this codebase. On win32 the same door is an open
handle or AV lock on any skill directory — not an exotic state — though I could not measure
that arm here (macOS).

What must change: install's replacement path gets the same fleet transaction shape uninstall
got — stage every unit fleet-wide, set every old destination aside with rollback, commit,
then clean up — or the deck explicitly ratifies this residual as a recorded limit and
amends B5's "no partial fleet" reading and the "install needed nothing further" claim.
A two-phase install has a real crash-window question (dest/.tmp/.bak generations) that
uninstall's did not; that is a design conversation, but it is the conversation cycle 15
started and only finished for one command.

### MINOR — a non-array, non-`false` `addons` fails open as "core-only"

`markerAddonNames` validates the array's *elements* but not its *type*: any non-array value
that isn't an array — `"parley-bidding"`, `42`, `{}`, `true` — takes the same branch as the
legitimate `false`/absent shapes and is read as a recorded core-only selection (measured,
A4/A4b). Consequences: `uninstall --force` then removes **only the core**, leaves all four
add-on trees on disk, and returns `ok:true`; health does not stay green (the on-disk add-ons
report `valid-unselected`, and after the core is gone, `missing`), and a subsequent install
heals the marker. The CRITICAL class stays closed — no path is ever constructed from a
non-array — but this is a residual fail-open branch inside the very function cycle 15 made
fail-closed, reachable by exactly the kind of hand-edit slip (`"parley-bidding"` for
`["parley-bidding"]`) the validation exists for. One branch fixes it without touching the
ratified legacy rule: `undefined`/`false` → core-only; array → validate; anything else →
`markerProblem`. Does not block the release on its own; fix with the MAJOR.

### NIT — the discrimination table's count is off by one

IMPLEMENTATION.md records "3/11 pass at `12f9071`." Running the cycle-15 test file against
`12f9071`'s `lib/`, I measure **2 of the 11** diff-touched tests passing at both commits
(the two preserved properties); the other 9 fail. The third passing test in the record is
presumably a quarantine-adjacent preserved-property test outside the diff set (e.g. the
committed-replacement/backup-cleanup regression). The substance — every test asserting new
behavior discriminates — holds; the count should say 2 of the 11 diff-touched, or name the
eleventh.

### Observations, not findings

- **Crash mid-phase-A** is a state the old design did not have: some destinations renamed
  aside, none deleted, no automatic rollback. It is strictly safer than the old crash
  window (half-deleted trees at the destination): every quarantined tree is intact, named,
  and the next uninstall completes the job (`missing` is the no-op). Recorded here because
  the prompt asked; no change requested.
- **Coverage gap, behavior correct:** the rollback-failure branch (`Rolled back, but the
  directory could not be restored… it is at <aside>`) and rollback byte-fidelity have no
  regression test — I exercised both by injection (C1) and hash comparison (B3) and both
  behave as designed. Worth one regression so the branch cannot rot silently.
- Concurrency, name collisions, `--dry-run`, symlinked destinations, and install-backup
  interaction (the prompt's attack list): all probed, all benign — see What I verified.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: **install must not leave a partial
fleet when a commit-time operation fails mid-fleet.** Either transact the replacement path
the way uninstall now is, or the deck explicitly ratifies the residual as a recorded limit —
in which case IMPLEMENTATION.md's "Install needed nothing further: it already commits by
rename" must be corrected, because B1/B2 falsify it. The `uchg`-destination arm is not
inherited risk; it is a door this cycle opened (`12f9071`: 0 replaced, clean block;
`5100f34`: 83 replaced, then EPERM). Everything round 11 gated on is verifiably closed, and
the MINOR and NIT above ride along cheaply with the MAJOR's cycle.

## What I verified

- Read end-to-end at `5100f34`: `lib/installer.js` (2295 lines), `lib/addon-manifest.js`,
  `bin/parley-deck-skill.js`, `package.json`; full diffs `12f9071..5100f34` (all three
  files) and the `49fc3ec..5100f34` scope (stat + the `design-addons.test.js` portion —
  comment corrections and two grammar-sentinel assertions, test-only).
- Working tree never mutated: no resets/checkouts/edits; no scratch files in the repo. All
  runs in rsync copies `/tmp/pd-r12/{head,prev}` (prev checked out at `12f9071`); probes in
  `/tmp/pd-r12/probes`; synthetic HOMEs under `os.tmpdir()`; all `chflags`/`chmod` arms
  reverted after use.
- Full gate at `5100f34` (copy): **338 node tests, 0 fail**; again with `python3` shimmed to
  `/usr/bin/python3` 3.9.6: **338/338**. Python leg **54/54 on 3.14.6**. Manifest check ok,
  47 files, `sha256:7854adf1…b95a6d`.
- Discrimination: cycle-15 `test/bidding-addon.test.js` against `12f9071`'s `lib/` → 9 fail,
  2 pass (preserved properties); against `5100f34` → all pass.
- Round-11 findings re-measured closed (above), including the original CRITICAL arm
  end-to-end through `run()` with `--json`: exit 1, sentinel intact, core intact.
- Quarantine transaction probed in every state I could construct: phase-A rollback fidelity
  (byte-identical, zero debris), injected rollback double fault (aside kept and named),
  `uchg` file-in-tree arm (84/84 removed, 1 warning naming debris), `uappnd` parent (full
  rollback), six concurrent two-process runs (one winner, loser clean, no half-deleted
  trees), dry-run purity, live-symlink destinations (link removed, target intact), install
  over a live symlink (replaced, target preserved).
- Install-side class measured: B1 (`uchg` dest root → 83 replaced then EPERM at the backup
  rename), B2 (`uappnd` skills dir → 78 replaced then EPERM in staging), and the
  `12f9071` comparison (0 replaced, clean block) plus the primitive `access(2)`-on-`uchg`-dir
  → EPERM that explains the regression. Uninstall under the same flags: fully rolled back.
- Stored-data→path routes enumerated by reading every `path.join`/`path.resolve` consumer:
  the marker `addons` array was the only route from stored data to a constructed path; it is
  validated and confined. All other marker fields are compared or displayed only.
- Record-keeping: IMPLEMENTATION.md cycles 10–15 read; the two claims I could not reproduce
  verbatim are the "3/11" count (NIT above) and "Install needed nothing further" (the
  MAJOR). I did not read the other round-12 reviews before writing this file.
