---
idea: integrate-parley-bidding-addon
review-round: 10
agent: hermes-1
date: 2026-07-31
reviewed-commit: 9ed2081
---

## Verdict
BLOCK

## Outstanding findings — closed or not

All previously raised findings are closed or correctly deferred, with one exception
that I believe is new this round (below).

- **Cycle 10 (fleet-wide install preflight):** CLOSED and verified. `installCommand`
  at `lib/installer.js:618-654` builds the full target × unit plan and preflights every
  unit before the first write. I confirmed a blocker in the last target writes nothing
  in any earlier target by running the scenario directly (see What I verified).

- **Cycle 11 (dangling destinations in health and uninstall):** CLOSED and verified.
  `skillUnitStatus` and `uninstallSkillUnit` both use `pathEntryExists` (lstat, only
  ENOENT is absence). A dangling symlink is `malformed` in doctor, `blocked` in unforced
  uninstall, and removed by `--force` uninstall — all three confirmed by running the
  existing regressions and by independent probes.

- **Cycle 12 (destination ancestor obstacle, file arm):** CLOSED and verified.
  `destinationAncestorObstacle` walks to the nearest existing entry and blocks when it
  is not a directory, regardless of `--force`. I confirmed a regular file at
  `~/.aionrs/skills` produces zero writes under both `--force=false` and `--force=true`.

- **Cycle 13 (destination ancestor obstacle, permission arm):** CLOSED and verified.
  The walk also requires `W_OK | X_OK` on the nearest existing ancestor. I confirmed a
  mode-000 directory at `~/.aionrs/skills` produces zero writes under both force
  settings. Confirmed not running as root (uid 501), so the permission check is
  genuinely exercised.

- **`firstCopyObstacle` source preflight (round 1 / round 8):** CLOSED and verified.
  The mirror faithfully covers what `copyRecursive` refuses or fails on: symlinks,
  non-regular files, unreadable entries. I compared the two functions line by line.
  `copySourcesFor` covers both add-on roots and core package entries. The core-source
  regression and the manifest-free-add-on regression both pass.

- **kimi-1's NIT (dirExists in targetSkillUnits discovery guard):** Correctly deferred.
  I confirmed the behavior: a dangling symlink at an unselected add-on path is invisible
  to unflagged `doctor` while a real directory there is reported. Nothing usable is
  installed at that path, and the mutation paths are coherent (`install --only` blocks
  with an accurate message, `--force` remediates). Non-blocking. I agree with leaving it
  as a follow-up rather than absorbing it, since it changes discovery semantics that
  rounds 4 and 6 ratified.

- **B3.11 manifest follow-up (five remaining skills):** Correctly deferred. Only
  `parley-bidding` ships a manifest. A universal install of all six skills reports one
  `valid-unmanaged` and five `malformed`. This is the first follow-up in the consensus
  and is stated in `CHANGELOG.md`. Not a 2.1.0 blocker.

## New findings

### [MAJOR] Uninstall is not fleet-wide — a blocker in the last target deletes the earlier targets first (B5)

**Where:** `lib/installer.js:665-683` (`uninstallCommand`), `:1127-1155` (`installTarget`
fleet preflight, for contrast), `:1221-1274` (`uninstallTarget`).

**What:** Cycle 10 added a fleet-wide preflight to `installCommand`: the complete
target × unit plan is checked before any write, and a single blocker anywhere returns
every unit as `blocked`/`skipped` with zero writes. `uninstallCommand` never received
the same treatment. It calls `targets.map((target) => uninstallTarget(...))` directly —
no pre-pass, no fleet-wide gate. Each target is preflighted only inside `uninstallTarget`,
immediately before that target is deleted. A predictable blocker in a later target is
discovered after earlier targets have already been removed.

This is the exact mirror of the defect cycle 10 fixed for install. The per-target
blocker check in `uninstallTarget` (lines 1228-1261) provides within-target atomicity —
a refusal on the last unit of one target does not leave earlier units of that target
deleted — but across targets the deletion is sequential and non-atomic.

**Why it matters:** B5 says "The installer is atomic per skill directory, not per
selected set. Preflight every unit and destination before the first write; a
predictable failure must produce zero writes." The requirement names "the installer,"
not "the install command." A deletion is a write (a mutation of the destination tree),
and a partial fleet of deletions is the same unacceptable state as a partial fleet of
installs: the user is left with some runtimes carrying the skill and some without, with
no indication of which is which, when the failure was statically discoverable before any
deletion.

**Evidence:** I reproduced this through `installer.uninstallCommand` at `9ed2081` in an
isolated temporary HOME:

1. Installed `--target all --include-undetected` (all 14 targets, ok:true).
2. Corrupted the last target's core marker (`aionrs`) to name a foreign installer —
   a statically discoverable ownership mismatch, exactly the kind of predictable failure
   B5 names.
3. Ran `uninstall --target all --include-undetected` (unforced).

Result: `ok:false`. The first 13 targets (codex through opencode) were **removed**.
The 14th (aionrs) was **blocked**. Six of the seven checked target directories no
longer had their core; aionrs still did. A partial fleet.

No existing test covers `uninstall --target all` with a blocker in a later target. Every
uninstall regression uses a single target (`codex`). The gap is both untested and
unaddressed.

**Fix:** Add a fleet-wide preflight to `uninstallCommand` structurally identical to the
one `installCommand` already has: build the full target × unit plan, check every unit's
ownership/destination before the first deletion, and if any blocker exists, return every
unit as `blocked`/`skipped` with zero deletions. The ownership check already exists
per-target in `uninstallTarget`; lift it to the fleet level the same way cycle 10 lifted
`preflightSkillUnit`. Add a regression: a foreign marker in the last target asserts no
earlier target was deleted.

**Severity:** MAJOR. This blocks the 2.1.0 release. It is the same B5 class that cycles
10-13 fixed for install, on the uninstall side, reachable through a documented command
(`README.md:230`: `uninstall --target all`).

## Release judgement

Not releasable as 2.1.0. The uninstall fleet-wide preflight must be added — the same
gate cycle 10 gave `installCommand`, applied to `uninstallCommand` — with a regression
that asserts zero deletions when a blocker exists in any target. This is cycle 14.

## What I verified

- **Suite:** `node --test` — 325 tests, 0 fail. Python leg — 54/54 across 7 files under
  Homebrew python3 3.14.6. Manifest check — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, matching
  the stated value unchanged since `714712f`. All three confirmed by running them.

- **Install fleet-wide preflight (cycle 10):** Ran `installCommand` with
  `--target all --include-undetected` and a foreign `parley-bidding` destination in the
  last target. `ok:false`, zero writes in any earlier target. The existing regression
  ("a blocker in the LAST target writes nothing in any earlier target") passes.

- **Destination ancestor obstacle (cycles 12-13):** Ran both arms — regular file at
  `~/.aionrs/skills` and mode-000 directory at `~/.aionrs/skills` — under
  `--force=false` and `--force=true`, each with `--target all --include-undetected`.
  All four: `ok:false`, zero installed/replaced units, `aionrs` blocked with the
  expected message. The symlink-to-real-directory positive case installs correctly.

- **Dangling destination symlink (cycle 11):** Confirmed `pathEntryExists` returns true
  for a dangling link (lstat succeeds), `skillUnitStatus` reports `malformed` (not
  `missing`), unforced uninstall reports `blocked`, and `--force` uninstall removes it.
  All via the existing regressions and independent probes.

- **Source preflight mirror (rounds 1/8):** Compared `firstCopyObstacle`
  (`lib/installer.js:1068-1124`) against `copyRecursive` (`:1412-1426`) line by line.
  The mirror catches every case the copy would fail on: symlinks (both refuse),
  non-regular files (mirror checks `isFile`, copy would throw on `readFileSync`),
  unreadable entries (mirror checks `accessSync(R_OK)`, copy would throw on
  `readFileSync`). `copySourcesFor` covers add-on roots and all core payload entries
  including optional ones. The two existing regressions (symlink in manifest-free
  add-on, symlink in core source) pass with zero writes.

- **Uninstall fleet-wide gap (new finding):** Reproduced through `uninstallCommand` in
  an isolated HOME. Installed all 14 targets, corrupted the last target's core marker
  to a foreign name, ran unforced `uninstall --target all`. 13 targets removed, 1
  blocked — partial fleet. Confirmed no existing test covers `uninstall --target all`
  with a later-target blocker (all uninstall regressions use single-target `codex`).

- **Deferred NIT (dirExists in discovery):** Confirmed a dangling symlink at an
  unselected add-on path is invisible to unflagged `doctor` while a real directory is
  reported. Non-blocking; correctly deferred.

- **Dry-run path:** The fleet-wide preflight is skipped for `--target all` dry-run
  (`if (!context.options.dryRun)` at line 622), but `installSkillUnit` still checks
  ownership at the per-unit level even in dry-run. A dry-run with an impossible
  destination reports `ok:false` with a `blocked` action — the message cites ownership
  rather than the ancestor obstacle, but this is cosmetic: dry-run is non-mutating, so
  no partial fleet can occur.

- **Per-target preflight redundancy:** `installTarget` (line 1132) re-runs
  `preflightSkillUnit` for every unit after the fleet-wide gate already passed. This is
  redundant but harmless — the fleet gate is a strict superset. Not a finding.

- **`destinationAncestorObstacle` termination:** The walk uses `path.dirname` until
  `dir === previous`. On Unix, `path.dirname("/")` returns `/`, so the loop terminates
  at root. For any realistic destination, the home directory (or a temp dir) is the
  nearest existing ancestor and the walk stops there. No infinite-loop risk.

- **Not running as root:** `uid=501(tomasfecko)`. The mode-000 permission arm is
  genuinely exercised.
