---
idea: integrate-parley-bidding-addon
review-round: 5
agent: codex-1
date: 2026-07-30
reviewed-commit: 3634cc8
---

## Verdict

BLOCK

## Round-4 findings — closed or not

1. **`valid-unmanaged` trusted the installed manifest — CLOSED.** At `3634cc8`,
   `unmanagedButVerified` first verifies the packaged source, requires the installed
   manifest hash to equal the source manifest hash, and then verifies the installed
   payload. The runtime-field-removal and self-consistent replacement probes now return
   `malformed`.
2. **Mutation ownership used marker-path existence — NOT FULLY CLOSED.** Foreign-name and
   unreadable markers are refused without `--force`, and uninstall now preflights the
   selected set. However, health does not apply the new matching-skill predicate. A marker
   with our `name` but another unit's `skill` is healthy to `doctor` and unowned to install
   and uninstall. See the first MAJOR.
3. **An excluded add-on could remain installed but disappear from health — CLOSED for the
   original unflagged transition.** After a full copy followed by
   `install --force --no-addons`, an ordinary `doctor` surfaces the residual add-on with
   `selected:false` and fails health. The same works for generic `--dest` and project scope.
   The traversal has a separate explicit-selector false-negative, described in the second
   MAJOR.
4. **A non-regular marker counted as absent — CLOSED.** Directory and dangling-symlink
   markers are present-but-unreadable and `malformed`; only `ENOENT` is absent.
5. **The probe cache omitted effective interpreter selectors — NOT FULLY CLOSED.** JSON over
   all environment pairs closes the `PYENV_VERSION` and separator-collision cases, but the
   child process's effective working directory is still absent from the key. See the MINOR
   finding.
6. **Follow-up output-shape defects — CLOSED.** `managed` is present on missing units and
   the `status` text renderer names missing files.

## New findings

### [MAJOR] Health omits the matching-skill half of the shared ownership predicate

**Where:** `lib/installer.js:1323-1357`, `lib/installer.js:1591-1602`

**What:** `installerOwnsDestination` requires a present, readable marker with our package
name and `marker.skill === unit.skill`; install and uninstall call it. `skillUnitStatus`
does not. Its readable-marker branch checks only `marker.name`, so a marker for another
skill reports `valid` and `managed:true`.

**Why it matters:** The same destination is simultaneously healthy and managed according to
`doctor`, but unowned according to both mutation commands. A normal upgrade or uninstall can
therefore refuse a tree immediately after the advertised health gate approved it. This is
the ownership inconsistency round 4 required one predicate to remove.

**Evidence:** In an archive of `3634cc8`, I performed a normal full install and changed only
the installed bidding marker's `skill` from `parley-bidding` to `parley-design`. Payload,
manifest anchor, and every other marker field were unchanged. `doctor.ok` was `true`; the
bidding unit was `status:"valid"`, `managed:true`, with no problems. Both
`install --only parley-bidding` and `uninstall --only parley-bidding` then returned
`action:"blocked"` for bidding and `action:"skipped"` for the core. A sentinel added after
the health probe survived both refusals.

**Fix:** Make health use the same parsed ownership result. At minimum, a readable marker
whose package name is ours but whose `skill` does not equal the unit must add an explicit
identity problem, report `malformed`, and make `doctor` fail. Add core and add-on
skill-mismatch regressions that assert the same ownership answer across doctor, install, and
uninstall.

### [MAJOR] Explicit read selectors falsely relabel recorded add-ons as unselected

**Where:** `lib/installer.js:858-929`, `lib/installer.js:1336-1351`

**What:** `expectedAddonNames` lets an explicit `--only` or `--no-addons` override the core
marker's recorded selection. `targetSkillUnits` then uses that command-local set as the
`seen` set and labels every other packaged add-on directory `selected:false`. The reported
field and problem claim those units are absent from the *recorded* selection even when the
core marker records them.

**Why it matters:** A healthy full install becomes unhealthy merely because a user asks a
read-only command to focus on a subset. The output is also factually wrong and recommends
removing correctly selected, managed skills. This makes `doctor --only ...` unsafe as a
selective health probe and makes the analogous `paths` selector expand rather than narrow
the result.

**Evidence:** A default install recorded all five add-on names in the core marker. On that
unchanged tree, `doctor --only parley-bidding` returned `ok:false` and labeled the other four
recorded add-ons `selected:false`, `malformed`, with “not part of the recorded selection.”
`doctor --no-addons` labeled all five recorded add-ons that way. By contrast, the intended
unflagged residual detection worked under user scope, project scope, and
`--target generic --dest`; sibling directories named `totally-unrelated-skill` and
`parley-bidding-archive` were correctly ignored.

**Fix:** Keep the recorded selection distinct from a read command's requested filter. For
unflagged read commands, continue adding known on-disk directories absent from the core
marker as `selected:false`. For explicit read selectors, inspect only the requested units,
or surface other units without claiming they are absent from the recorded selection and
without failing health. Add default-full-install regressions for `doctor`, `status`, and
`paths` with both explicit selectors.

### [MINOR] The probe cache still ignores the working directory used by a relative `PATH`

**Where:** `lib/installer.js:1423-1447`

**What:** The new cache key serializes the whole environment but not the effective working
directory. `spawnSync` receives no `cwd`, so relative `PATH` entries resolve against the
process working directory.

**Why it matters:** A library process can change working directory while retaining the same
environment. The second health check then reuses an interpreter verdict from a different
executable. This is the remaining narrow form of the round-4 cache false green.

**Evidence:** In one process I used the identical environment `PATH=bin`. With the process
working directory containing `bin/python3` that printed `3.12`, doctor reported
`python3 3.12`. I changed to a directory whose `bin/python3` printed `3.9` and repeated the
call with the same environment. The second result was the cached `3.12`, and doctor stayed
green against the `>=3.10` floor.

**Fix:** Pass the intended effective `cwd` into the runtime probe and `spawnSync`, and include
its resolved value in the cache key. Add the same-relative-`PATH`, different-working-directory
regression in both call orders.

### [MINOR] The validation record incorrectly says full `npm test` passes on Python 3.9

**Where:** `IMPLEMENTATION.md:162`, `IMPLEMENTATION.md:501-505`

**What:** The record says `npm test` produced 305 Node plus 54 Python passes on Python 3.9.6
and 3.14. The committed runner intentionally rejects any interpreter below the manifest's
`>=3.10` floor before running the Python files.

**Why it matters:** The code is correct to enforce its declared floor, but the canonical
validation evidence is not reproducible as written.

**Evidence:** With `/usr/bin/python3` 3.9.6 first on `PATH`, the Node leg passed 305/305 and
then `scripts/run-python-tests.js` exited 1 with “python3 is 3.9, but the add-on declares
>=3.10”; it ran zero Python tests through `npm test`. Run individually with
`PYTHONDONTWRITEBYTECODE=1 python3 -B`, all 54 tests did pass on 3.9.6. The supported-floor
measurements also passed 54/54 on Python 3.10.20, 3.11.15, and 3.14.6.

**Fix:** Correct the record to distinguish the successful 3.9 file-by-file compatibility
probe from the intentionally failing full package gate. State that full `npm test` passed on
supported Python versions only.

## What I verified and found correct

- I reviewed an isolated `git archive` of
  `3634cc81e571046ab639e9fa1caa60ebcfe310dc`; the source worktree had later local changes in
  the same two files and was left untouched. `git diff --check b180127..3634cc8` was clean.
- With the checkout's dependency directory supplied read-only through `NODE_PATH`,
  `PYTHONDONTWRITEBYTECODE=1 npm test` passed on Python 3.14.6: **305/305 Node** and
  **54/54 Python**, plus the 47-file manifest check at
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Source-anchored unmanaged verification rejects a runtime-field deletion and a laundered
  payload/manifest pair. The intact packaged copy still earns `valid-unmanaged` and remains
  runtime-probed.
- Foreign-name and unreadable markers are refused by install without `--force`; uninstall
  preflight preserves the full selected set when one unit is unowned.
- The original `selected:false` transition works at the normal user destination, a generic
  custom destination, and project scope. Unrelated sibling directory names do not enter the
  traversal.
- Actual installs made by tags `v1.0.0`, `v1.4.0`, and `v2.0.0` all upgraded at `3634cc8`
  without `--force`. I found no legitimate released marker shape newly refused by the
  stricter ownership predicate.
- No `__pycache__`, `.pyc`, or `.pyo` appeared in the archive. I did not reopen the
  deliberately deferred manifest work for the other add-ons.

## Open questions for the implementer

1. For read-only commands, are `--only` and `--no-addons` filters over what to inspect, or
   desired-state assertions? Existing output calls the result the *recorded* selection, but
   the implementation currently substitutes the command-local selection.
