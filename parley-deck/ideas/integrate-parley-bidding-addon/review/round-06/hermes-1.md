---
idea: integrate-parley-bidding-addon
review-round: 6
agent: hermes-1
date: 2026-07-30
reviewed-commit: c673111
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

### Round-4 MINOR — status text renderer does not print `missing` files for the core — CLOSED

My round-4 finding was that the `status` text renderer's `explain` function printed
`skill.problems` and `skill.runtime` but not `skill.missing`, so `status` on a tree
with a missing core file showed the verdict but not the cause, while `doctor` showed
both. The `explain` function at lib/installer.js:1755-1767 now prints `missing:` with
a comment citing "review round 4, hermes-1 MINOR."

Measured: installed into an isolated HOME, deleted
`references/COOPERATION.md` from the core, ran `status --target codex` via the CLI.
Output includes `  missing: references/COOPERATION.md` for the core, matching
`doctor`'s output. Both commands now name the missing file.

### Round-3 conditions (valid-unmanaged with managed boolean) — CLOSED (since round 4)

Re-verified at c673111: `valid-unmanaged` carries `managed: false`, `marker: null`;
the core can never be `valid-unmanaged` (`unmanagedButVerified` returns false when
`unit.addon` is null); a non-manifest add-on cannot launder into `valid-unmanaged`
by manufacturing a manifest in the installed tree (source-anchored check).

## New findings

### [NIT] valid-unselected masks valid-unmanaged for a foreign copy outside the recorded selection

**Where:** lib/installer.js:1366-1372 (status precedence in `skillUnitStatus`)

**What:** The status precedence chain is:

```
malformed > valid-unselected > valid-unmanaged > valid
```

When a foreign-installed copy of `parley-bidding` (no marker, manifest verifies) is
on disk but outside the recorded selection (the core marker records a different
add-on set, or `addons: false`), the `selected === false` check fires before the
`unmanaged` check, so the status is `valid-unselected` rather than `valid-unmanaged`.

This means the status string implies "this tool installed it and it was selected,
but now it is not" — when in fact the tool never installed it at all. The `managed`
boolean is `false` in both cases, which correctly disambiguates.

**Why it matters:** It does not affect behavior. I measured every mutation path:

- `install --force` reclaims it (action: `replaced`, becomes `valid` + `managed: true`)
- `install` without `--force` blocks it (no marker → `installerOwnsDestination` false)
- `uninstall` without `--force` blocks it (same reason, tree preserved)
- `doctor.ok` is `false` — the same as `valid-unselected`, and stricter than
  `valid-unmanaged` (which could pass health with an available runtime)

The problem message ("installed but not part of the recorded selection: remove the
directory, or re-run install including it") is reasonable advice regardless of who
installed the tree. The `managed: false` boolean is the provenance fact; the status
string is the health fact. A consumer that checks `managed === false` to distinguish
"foreign" from "tool-installed-but-unselected" gets the right answer. A consumer
that parses the status string alone would get the wrong provenance — but the code
comments at line 1378-1380 explicitly direct automation to use `managed`, not the
status string, for exactly this distinction.

**Evidence:** In an isolated copy of c673111, installed core + `parley-design` only
into HOME A, then foreign-copied `parley-bidding` (payload only, no marker) from a
full install in HOME B into A's skills directory. `doctor` reported
`parley-bidding: valid-unselected, managed: false, selected: false`. The same
foreign copy placed inside the recorded selection (core marker records
`parley-bidding`) reported `valid-unmanaged, managed: false, selected: true`.

**Fix:** Optional. If status-string precision matters, check `unmanaged` before
`selected === false` so a foreign unmanaged copy outside the selection gets
`valid-unmanaged` (with `selected: false` still available as a separate field).
Or leave as-is: the `managed` boolean is the authoritative provenance signal, and
the current precedence is the stricter health outcome. Not blocking.

## What I verified and found correct

1. **309/309 node tests pass on python3 3.9.6 and 3.14; 54/54 Python tests pass on
   3.14.** On 3.9.6 the Python leg refuses by design ("python3 is 3.9, but the
   add-on declares >=3.10") and zero Python tests execute through `npm test`. The
   corrected validation record (IMPLEMENTATION.md:559-560) is accurate.

2. **Ownership is one answer.** A marker with our package name but a different
   `skill` identity is `malformed` to `doctor`, `blocked` to `install`, and
   `blocked` to `uninstall`. Measured by changing `marker.skill` from
   `parley-bidding` to `parley-design` on an installed tree: all three commands
   agree. The identity problem message ("identifies this directory as
   \"parley-design\", not \"parley-bidding\"") is present. (codex-1 F14, closed.)

3. **A read command's `--only`/`--no-addons` is a filter, not a claim about the
   recorded selection.** On a healthy full install, `doctor --only parley-bidding`
   inspects only the core and `parley-bidding`, both `valid` with empty problems
   and `selected: true`. No recorded add-on is falsely labeled `valid-unselected`
   or accused of being outside the selection. `doctor --no-addons` inspects only
   the core. `paths --only` behaves identically. (codex-1 F15, closed.)

4. **`valid-unselected` is a distinct status, not `malformed`.** An add-on
   installed by an earlier `--only` run, then excluded by a later `--only` run
   that names a different set, reports `valid-unselected` with `managed: true`,
   `selected: false`, `missing: []`, and a problem naming the remedy. Health
   fails (`doctor.ok: false`). Re-including the add-on in a subsequent install
   fixes it: status becomes `valid`, `doctor.ok: true`.

5. **`valid-unselected` via `--no-addons` works.** A full install followed by
   `install --force --no-addons` leaves the bidding add-on on disk as
   `valid-unselected, managed: true, selected: false`. Health fails. The add-on
   did not vanish from `doctor`'s traversal.

6. **The `doctor.ok` gate correctly excludes `valid-unselected`.** Only `valid`
   and `valid-unmanaged` pass (lib/installer.js:384). `valid-unselected` fails
   health regardless of runtime availability. The `runtime` field is `null` for
   `valid-unselected` (since `ok` is false, the probe is skipped), so the stderr
   message reports "installed but outside the recorded selection" without also
   reporting "operationally unavailable" — the selection problem is primary.

7. **The stderr message correctly distinguishes three failure reasons.**
   `writeResult` (lib/installer.js:1728-1737) checks `broken` (missing/malformed),
   `unselected` (valid-unselected), and `unavailable` (runtime not ok) as separate
   conditions, joining them with ", or ". A tree with both a `valid-unselected`
   add-on and a `malformed` add-on produces "missing or malformed, or installed
   but outside the recorded selection."

8. **`status` reports `valid-unselected` but always exits 0.** This is the
   documented contract: `status` is informational, `doctor` is the health gate.
   The `valid-unselected` status is visible in `status` output but does not
   change its exit code.

9. **Legacy markers without a `skill` field are backward-compatible.** A marker
   with `name` matching ours but no `skill` field (or `skill: undefined`) does
   not trigger the identity problem — the guard is
   `state.marker.skill !== undefined && state.marker.skill !== unit.skill`.
   Measured: both `delete marker.skill` and `marker.skill = undefined` produce
   `valid, managed: true` with no identity problem. This preserves compatibility
   with 2.0.0-era markers.

10. **Unrelated sibling directories are ignored.** Directories named
    `totally-unrelated-skill` and `parley-bidding-archive` in the skills directory
    do not appear in `doctor` output and do not affect any skill's problems.

11. **The probe cache key includes the working directory.** `probePython3`
    (lib/installer.js:1460-1476) resolves `cwd` via `path.resolve(cwd)`, includes
    it as the first element of the JSON cache key, and passes it to `spawnSync` as
    `cwd: workingDir`. Two calls sharing an environment but not a directory get
    distinct cache entries. (codex-1 F16, closed.)

12. **`doctor` and `status` pass `cwd` through to the runtime probe.** Both
    commands now include `cwd: context.cwd` in the options passed to
    `targetStatus` (lib/installer.js:374, 395), which flows through to
    `runtimeAvailability` and `probePython3`.

13. **Install and uninstall of `valid-unselected` add-ons behave correctly.**
    An owned unselected add-on (has our marker) can be reclaimed by `install`
    without `--force` (marker proves ownership) or removed by `uninstall`
    without `--force`. A foreign unselected add-on (no marker) is blocked by
    both without `--force`, and the tree is preserved.

14. **The `managed` boolean is uniform across all unit shapes.** `managed` is
    present on every status, including `missing` (where it is `false`) and
    `valid-unselected` (where it is `true` for tool-installed, `false` for
    foreign). JSON automation can rely on its presence.

## Open questions for the implementer

1. The `valid-unselected` masking of `valid-unmanaged` (NIT above) is the one
   edge case where the status string's provenance implication does not match
   reality. The `managed` boolean is correct in all cases. Is the status-string
   precision worth a precedence swap, or is `managed` the intended authority and
   the current stricter health outcome preferred? Not blocking either way.
