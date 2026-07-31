---
idea: integrate-parley-bidding-addon
review-round: 17
agent: kimi-1
date: 2026-07-31
reviewed-commit: fa3da41
---

## Verdict
BLOCK

## Outstanding findings — closed or not

All round-16 findings are closed as far as their fixes reach — one of them does not reach far
enough, and that is the MAJOR below.

- **codex-1 MAJOR (a commit can destroy a resolution chain) — closed for symlink spellings,
  open via mount spellings.** The cycle-20 touchpoint check works: I independently re-measured
  the round-16 scenario at fa3da41 (managed kimi core, `redirect` symlink planted inside it,
  `CODEX_HOME` pointing at the link) and the plan is refused with "resolving one passes
  through the other", nothing written through the link — install side, uninstall side (no
  quarantine, no debris), and identically under `--dry-run`. But the same defect survives when
  the dependent home is spelled through a mount boundary; see MAJOR-1. The fix is correct
  within one spelling and incomplete across two.
- **codex-1 MINOR (unreadable python version failed open) — closed.** `run-python-tests.js`
  now anchors `/^(\d+)\.(\d+)$/` and fails on anything else. Measured: 54/54 under python3
  3.14.6; exit 1 with `python3 is 3.9, but the add-on declares >=3.10` under
  `/usr/bin/python3` 3.9.6 (PATH-restricted).
- **kimi-1 MINOR (malformed `runtime.python` → null floor, comparison skipped) — closed.**
  `declaredPythonFloor` returning null now fails the leg ("refusing to test against an
  unbounded interpreter") instead of skipping the comparison. `runtimeAvailability` was
  already fail-closed on the same spec; the two consumers now agree.
- **hermes-1 NIT-1 (`uninstall --dry-run` missing the per-action flag) — closed.** The spread
  is at lib/installer.js:678 and the regression ("uninstall --dry-run carries the same
  per-action flag install does") is in the suite.
- **kimi-1 NIT (dead `preflightUninstallUnit` + stale comments) — half-closed.** The dead
  function is gone. The stale comment is not: see NIT-1.

Full suite at fa3da41, measured here on a `git archive` copy: **357 node tests, 0 fail**;
python leg **54/54** on 3.14.6; manifest check ok, 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d` — unchanged since
`714712f`, as claimed.

## New findings

**MAJOR-1 — the plan gate compares realpath strings, and realpath is not a canonical form
across mount boundaries: spelling one home through `/System/Volumes/Data` (or any
firmlink/bind-mount double spelling) defeats BOTH cycle-19 containment AND cycle-20
touchpoints, reproducing the round-15 and round-16 MAJOR signatures at fa3da41.**

The question this round put to me was whether "only symlinks are touchpoints" is complete.
Within one path spelling, it is — I probed every escape class I could construct and each is
caught (list at the bottom so round 18 does not redo it). The narrowing's real gap is one
level down: both plan-gate checks reduce to `overlaps(a, b)`, a string prefix test over
`realpath`-derived paths, and `realpath` preserves the caller's mount spelling verbatim.
Measured on this machine:

- `realpathSync("/private/var/...")` → `/private/var/...`;
  `realpathSync("/System/Volumes/Data/private/var/...")` → unchanged — the *same directory*,
  two disjoint realpath strings.
- `/Users` is a firmlink: `lstat` reports a plain directory on **both** spellings (no
  symlink for the touchpoint walk to record), `realpath` keeps each spelling verbatim, and
  `statSync` returns **identical dev:ino**. So the respelling surface is not exotic — every
  home directory on macOS has two spellings, and only env-var convention keeps a plan on one
  of them. (Linux bind mounts, including unprivileged user-namespace ones, are the same
  mechanism; not measured here, stated as expectation.)

`physicalKey` already crosses mounts (`statSync` dev:ino), which is why *equality* survives:
two homes spelled across the mount boundary pointing at the **same** directory are still
refused as "Destination is shared" (measured). What string space cannot see is *nesting* and
*passing-through*. Two arms, both measured at fa3da41, single process, no external actor:

- **Arm A — containment side (round-15 signature).** `KIMI_CODE_HOME=<outer>`,
  `CODEX_HOME=/System/Volumes/Data/private/<outer>/skills/parley-deck`, so codex's
  destination is physically inside kimi's, with the shared ancestor spelled across the mount.
  `physicalKey`: same anchor dev:ino, different tails → not equal → passes. `resolvedDestination`
  strings: disjoint → no containment. Touchpoints: no symlinks anywhere → empty.
  `install --target all` returned **`ok: true`**, codex `installed`, kimi `replaced` — and
  codex's destination was **absent afterwards under both spellings**, no warning, no debris,
  surviving marker `target: kimi`. Commit physics: codex commits first (TARGETS order), kimi's
  staging had materialized its own dest as a plain dir, kimi's commit renamed that dir — with
  codex's fresh tree inside — to backup, and `discardBackup` rmSync'd it silently. (The
  reverse nesting — kimi inside codex — fails *safe*: ENOENT at kimi's commit, fleet rolled
  back, `ok:false`; only an empty materialized shell left at codex's dest. Order-dependent,
  which is worse, not better.)
- **Arm B — touchpoint side (round-16 signature, byte-for-byte).** The exact cycle-20
  regression setup — managed kimi core, `redirect` symlink inside it pointing at `<away>` —
  but `CODEX_HOME` spelled `/System/Volumes/Data/private/<kimiCore>/redirect`. The recorded
  touchpoint comes out mount-spelled, kimi's resolved path is natural-spelled, `overlaps` is
  false, plan passes. Result: **`ok: true`**, codex `installed`, kimi `replaced`, codex's
  destination no longer resolves, **orphan tree at `<away>/skills/parley-deck`**. Precisely
  the state codex-1 measured at a49d68f, reached through the surviving representation gap.

Why this blocks rather than follows up: the invariant cycles 19-20 wrote into the code is
categorical — "No planned destination may equal, contain, or be contained by another",
extended by "resolving one passes through the other" — and it is false as shipped, with the
same single-process false success this review rated MAJOR three rounds running. The
configuration is deliberate, but so was round 15's nested-homes setup, which was ratified and
fixed on the defect, not on configuration likelihood; and the double spelling itself ships on
every macOS machine — no exotic setup beyond typing the alternate prefix once.

Fix direction, one move closes both arms: evaluate containment **and** touchpoint-overlap in
`physicalKey` space (dev:ino of nearest existing ancestor + lowercased tail), which crosses
mounts and symlinks alike. Containment: equal anchor dev:ino and one tail a prefix of the
other ⇒ dependent plan. Touchpoint: the recorded symlink exists by construction, so key its
parent by `lstat` dev:ino and containment-compare that key against every unit's key. I
checked this construction against all four measured shapes (P7a–d): it refuses exactly the
four dependent plans and still passes plain siblings and a plain symlinked home.

**NIT-1 — stale comment carried from round 16.** The `pathEntryExists` doc comment
(lib/installer.js:2299-2301) still names `uninstallSkillUnit` as a caller; that function was
deleted in cycle 19, and `preflightUninstallUnit` — the other half of my round-16 NIT — in
cycle 20. Comment rot left behind by the two cleanups themselves. No behavioral impact.

Probed and dismissed, stated so round 18 does not have to redo it:

- *Completeness of "only symlinks" within one spelling.* (a) Plain-directory chain through
  another destination → final-path containment catches it (measured: home set to a real
  subdirectory of a managed dest — blocked). (b) Symlink stored inside a planned destination,
  pointing out → touchpoints (the cycle-20 regression; re-measured refused, plus its
  `--dry-run` parity). (c) Symlink at a neutral location pointing *at* a planned destination →
  the dependent's realpath passes through the target, so the final paths nest → containment
  (measured refused; the kimi tree byte-untouched, no `skills/` materialized inside it).
  (d) Dangling symlink anywhere in a dependent's chain → `destinationAncestorObstacle`'s
  `statSync` fails ("broken link or unreadable"), the unit is blocked, and the fleet gate
  holds — measured: `ok:false`, the other target NOT installed, nothing materialized through
  the dangling link. (e) *A symlink whose target is later replaced*: if the link resolves at
  plan time, the dependent's final resolved path contains the target → containment; if it is
  dangling → (d). Both roads end at an existing check. (f) *A component becoming a symlink
  between planning and commit*: impossible single-process — the installer creates no symlinks
  (`mkdirSync`, `copyRecursive` which refuses them, `writeFileSync`, renames only).
  (g) *Hardlinked directories*: not user-creatable on macOS or Linux. (h) *A mount point as a
  chain component* (as opposed to a spelling): a commit cannot destroy one — rename of a
  mountpoint fails EBUSY, and rmSync recursion into it fails at the final rmdir, leaving a
  debris warning with the chain intact. (i) *A firmlink as a component*: kernel-managed; the
  installer cannot rename or unlink it. Only the spelling arm escaped; that is MAJOR-1.
- *Uninstall coverage.* `removeFleetAtomically` runs the same `aliasedDestinations` gate
  before any quarantine rename (lib/installer.js:1626) — install and uninstall share the
  check. Measured: the round-16 dependent-home setup refused on uninstall with both
  destinations intact and zero `.removing` debris. With the gate in place, quarantine order
  no longer matters; without it (mount spelling), Arm A's commit-order dependence is the live
  arm.
- *False-positive arms of the narrowing.* Plain sibling homes install fine (no shared-parent
  flagging — the first-version bug stays fixed); a plain symlinked home installs fine; an
  aliased shared home is still refused; equal-destination-via-mount-spelling is still refused
  via `physicalKey`.
- *Commit/revert/housekeeping beyond the gate.* Re-read in full; the safe arm of Arm A
  exercised the revert path live — committed unit reverted, failing unit reported ENOENT
  honestly, `ok:false`. Rollback and revert are rename-only within one parent; no new hazard.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: the plan gate's containment and
touchpoint comparisons must run on physical identity (dev:ino anchor + tail) instead of
realpath strings — closing the mount/firmlink respelling arm of the very invariant cycles 19
and 20 ratified. Scoping the invariant down to one spelling and documenting the limit instead
would contradict three consecutive ratified MAJORs on this defect class; I do not recommend
it. NIT-1 is a one-line comment fix and can ride along.

## What I verified

- Repo untouched: `git status --porcelain` clean before and after, HEAD `fa3da41` throughout.
  All execution in a `git archive` copy under /tmp plus fabricated homes under TMPDIR; every
  temp path I created (archive copy, probe harness, probe homes) was removed afterwards —
  including a check for suite leftovers. (Older `/tmp/parley-r*` debris from previous rounds
  belongs to other reviewers; I left it alone.)
- Read in full: lib/installer.js (2558 lines), lib/addon-manifest.js, the cycle-20 diff
  `a49d68f..fa3da41`, and the round-16 review record.
- Node suite at fa3da41 in the archive copy: **357 tests, 0 fail** (repo `node_modules`
  symlinked read-only into the copy for `commonmark`; removed after).
- Manifest check: `parley-bidding: ok (47 files, sha256:7854adf1…a6d)`.
- Python leg: **54/54** under python3 3.14.6; exit 1 with the floor message under
  `/usr/bin/python3` 3.9.6.
- Adversarial probes, all in-process against the archive copy's installer with fabricated
  env, each in a fresh temp home: dependent-plan refused on install / uninstall / dry-run;
  dangling-link chain refused with fleet gate holding; link-at-managed-dest refused with tree
  untouched; sibling homes pass; symlinked home passes; aliased homes refused;
  equal-via-mount refused; **nested-via-mount false success (Arm A) and
  touchpoint-via-mount false success (Arm B) reproduced** — the two arms of MAJOR-1.
- Measured the filesystem facts the MAJOR rests on: realpath preserves
  `/System/Volumes/Data/private/...` verbatim; `/Users` firmlink is a plain directory to
  `lstat` on both spellings with identical dev:ino.
