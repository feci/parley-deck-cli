---
idea: integrate-parley-bidding-addon
review-round: 4
agent: hermes-1
date: 2026-07-30
reviewed-commit: b180127
---

## Verdict

ACCEPT

## Round-3 findings — closed or not

### Round-3 ruling (b) — CLOSED

My round-3 ruling asked for a third verdict when the marker is absent but the
manifest verifies, with a distinct status string and a `managed` boolean. The
implementer landed `valid-unmanaged` with `managed: false` and `marker: null`,
exactly as I asked. I re-measured every case from the ruling table:

| case | expected | measured |
|---|---|---|
| marker absent + source ships a manifest + tree verifies | `valid-unmanaged`, `managed:false`, `marker:null` | confirmed (TEST 1) |
| marker present but unreadable | `malformed` | confirmed (TEST: "unreadable or foreign marker") |
| marker naming another installer | `malformed` | confirmed (same test, foreign case) |
| marker absent + manifest absent where source ships one | `malformed` | confirmed (TEST 7: gutted to SKILL.md, no manifest) |
| legacy `markerSchema`-less marker | unchanged (`valid`) | confirmed by reading manifestProblems (line 1408-1411): `schema === undefined` → `return []` — no manifest check, marker name match still passes |

`doctor`'s top-level `ok` (line 381-386) now accepts both `valid` and
`valid-unmanaged` as payload-healthy, and still requires `(!skill.runtime ||
skill.runtime.ok)`. The `writeResult` stderr message (line 1569) correctly
excludes `valid-unmanaged` from the "missing or malformed" bucket. Measured:
a `valid-unmanaged` + unavailable tree produces stderr "operationally
unavailable" (not "missing or malformed"), exit 1.

### Round-3 MINOR — status always exits 0 — CLOSED

I flagged that `status` hardcodes `ok: true` even when it prints
`unavailable:` or `integrity:` lines. The implementer kept the always-0 exit
(documented as the command's original contract) and documented it: CHANGELOG.md
now says "`status` remains informational and always exits 0 … `doctor` is the
health gate." I accepted option (a) in round 3 and did not block. The
documentation closes it.

### codex-1 MINOR — broken shim satisfied the floor — CLOSED

`probePython3` (line 1339-1361) now parses with `/^(\d+)\.(\d+)$/` and requires
both capture groups to be integers. I verified the regex rejects `4.not-a-version`
(no match), `3.10.1` (no match — third part), `3` (no match — no dot), and
accepts `3.10` and `  3.10  ` (trimmed). The test "a broken interpreter shim
does not satisfy the floor" confirms the behavioural fix. The fail-open
direction is closed.

### kimi-1 NIT — cache key narrower than spawn — CLOSED

The probe cache key (line 1343-1345) now joins `PATH`, `PYTHONHOME`,
`PYTHONPATH`, `PYTHONEXECUTABLE`, and `VIRTUAL_ENV` with `\u0001` separators,
using `\u0000` for undefined values. I measured: the same PATH (`/usr/bin`)
with `PYTHONHOME=/nonexistent` produces `probe.found: false` (python3 3.9.6
at `/usr/bin/python3` can't find its stdlib), while the same PATH without
PYTHONHOME produces `probe.found: true` with version 3.9.6. A PATH-only key
would have returned the cached 3.9.6 answer for both; the expanded key
correctly distinguishes them.

### kimi-1 NIT — status explained add-ons but not the core — CLOSED

The `explain` function (line 1592-1599) is now called for both `units[0]` (the
core, at indent `"  "`) and the add-ons (at indent `"    "`). I verified: a
core with a deleted `references/COOPERATION.md` prints `codex: malformed` in
the status text, but the core's missing files are NOT printed (see new finding
below). The core's runtime (when it has one — it doesn't in this package) would
be printed. The explain function covers `problems[]` and `runtime`, which is
what the diff intended.

## New findings

### [MINOR] status text renderer does not print `missing` files for the core

**Where:** lib/installer.js:1592-1599 (`explain` in the `status` text path)

**What:** The `status` text renderer's `explain` function prints
`skill.runtime` (as `unavailable:`) and `skill.problems` (as `integrity:`),
but not `skill.missing`. The `doctor` text renderer (line 1554-1556) does
print `missing:`. So `status` on a tree with a missing core file shows
`codex: malformed` but not which file is missing, while `doctor` on the same
tree shows `  missing: references/COOPERATION.md`.

**Why it matters:** This is a pre-existing gap, not introduced by this diff —
the `status` renderer never printed missing files. The diff's `explain`
function made it more visible by adding per-skill detail that covers
`problems[]` and `runtime` but not `missing[]`. A user who runs `status`
instead of `doctor` on a broken core sees the verdict but not the cause.
Since `status` is explicitly informational (always exits 0) and `doctor` is
the health gate, this is a cosmetic asymmetry, not a correctness defect. The
JSON output carries `missing` in both commands.

**Evidence:** Installed into an isolated HOME, deleted
`references/COOPERATION.md` from the core, ran `status --target codex` (text):
output shows `codex: malformed` with no detail line. Same tree via `doctor`:
`  missing: references/COOPERATION.md`.

**Fix:** Optional: add `if (skill.missing && skill.missing.length > 0)` to the
`explain` function, mirroring the `doctor` renderer. Not blocking — `status`
is informational and the JSON output is complete.

### [NIT] test count discrepancy: 296 claimed, 299 actual

**Where:** IMPLEMENTATION.md line 426 ("all 296/0") vs. actual `node --test`
output

**What:** The implementer records 296 node tests at b180127. The working tree
is at 48d4fe4 (a follow-up commit adding 3 tests), which shows 299. 296 + 3 =
299, so the counts are consistent across commits. Not a defect — just noting
that the reviewed commit (b180127) and the working-tree HEAD (48d4fe4) differ
by 3 test-only commits. All 299 pass.

**Fix:** None needed. If precision matters, note that 48d4fe4 is in the tree
but after the reviewed commit.

## What I verified and found correct

1. **299 node + 54 Python tests, 0 fail.** `npm test` on python3 3.14.6 (homebrew)
   produces 299/299 node tests pass, 54/54 Python tests pass across 7 files, and
   the manifest build check passes (`parley-bidding: ok (47 files)`). Exit 0.

2. **The `valid-unmanaged` path is reachable only for add-ons whose SOURCE ships
   a manifest.** `unmanagedButVerified` (line 1294-1300) checks
   `addonManifest.hasManifest(unit.addon.root)` — the SOURCE root, not the
   installed tree. I measured: a foreign-copied `parley-design` (no source
   manifest) with a FORGED manifest placed in the installed tree still reports
   `malformed` (TEST 25). The source check is the gate, not the dest check. A
   non-manifest add-on cannot launder itself into `valid-unmanaged` by
   manufacturing a manifest in the installed tree.

3. **The core unit can never become `valid-unmanaged`.** The core unit has
   `addon: null` (line 890), and `unmanagedButVerified` returns `false` when
   `unit.addon` is null (line 1295). I measured: a core install with the marker
   deleted reports `malformed` (TEST 12). The follow-up test at
   test/bidding-addon.test.js:768 confirms this directly.

4. **`valid-unmanaged` does not grant ownership.** `install` without `--force`
   on a `valid-unmanaged` tree is blocked by `preflightSkillUnit` (line 913:
   `fs.existsSync(dest) && !fs.existsSync(markerPath(dest))` → blocked). I
   measured: install without `--force` → `action: "blocked"`, no marker
   synthesized (TEST 10). `uninstall` without `--force` → `action: "blocked"`,
   tree preserved (TEST 9a). `install --force` reclaims it: `action: "replaced"`,
   status becomes `valid` with `managed: true` (TEST 8). The follow-up test at
   test/bidding-addon.test.js:735 confirms all three.

5. **The laundering attack is not new.** A self-consistent forged manifest +
   tampered payload with no marker passes `verifyPayload` and gets
   `valid-unmanaged` (TEST 5, TEST 6). But the same attack WITH a forged marker
   gets plain `valid` — the marker is unsigned JSON, equally forgeable. The
   follow-up test at test/bidding-addon.test.js:779 ("ruling (b) adds no
   laundering weakness that (a) did not already have") measures both paths and
   proves the security delta is zero. The manifest system explicitly disclaims
   tamper resistance (addon-manifest.js:11-13: "this is defect detection, not
   tamper resistance"). The marker-present path's binding check
   (manifestProblems line 1452-1457) catches accidental manifest drift, not
   forgery — and the unmanaged path has no marker to drift against, so there is
   nothing to bind. This is the correct design boundary.

6. **`valid-unmanaged` trees are still runtime-probed.** The `runtime` field is
   set when `ok && options && options.probeRuntime` (line 1281), regardless of
   whether `unmanaged` is true. I measured: a `valid-unmanaged` tree with
   `PATH: ""` reports `runtime.ok: false` and `doctor` exits 1 (TEST 13). With
   `PATH: "/usr/bin"` (python3 3.9.6, below the >=3.10 floor), it reports
   `runtime.ok: false` with "python3 is 3.9, but this skill requires >=3.10"
   (TEST 26). Availability reporting is unchanged for unmanaged units.

7. **The `managed` boolean is carried alongside the status.** `doctor --json`
   and `status --json` both include `managed: true` for tool-installed trees
   and `managed: false` for unmanaged or malformed ones. Automation can require
   tool-managed installs by checking `managed === true` without parsing the
   status string. This is what I asked for in round 3.

8. **`--only` correctly surfaces foreign-installed add-ons.** When the core
   marker records `addons: false` (from a `--no-addons` install), the
   foreign-copied `parley-bidding` is invisible to `doctor`/`status` without
   `--only`. With `--only parley-bidding`, both the core and the named add-on
   are inspected: core is `valid` (tool-installed), bidding is
   `valid-unmanaged` (foreign, manifest verifies). `expectedAddonNames` (line
   868-872) makes `--only` the explicit override, which is correct.

9. **The `doctor` stderr message correctly distinguishes broken from
   unavailable.** Line 1569: `broken = skills.some(skill => skill.status !==
   "valid" && skill.status !== "valid-unmanaged")`. A `valid-unmanaged` tree
   with an unavailable runtime produces "operationally unavailable" only, not
   "missing or malformed". Measured in the FINAL test.

10. **Non-manifest add-ons with no marker stay `malformed`.** I copied
    `parley-design` (no source manifest) verbatim with no marker: `malformed`
    (TEST 19). This is the deferred residual — only `parley-bidding` ships a
    manifest, so the other five add-ons have nothing to verify against. This is
    stated in CHANGELOG.md and referenced from README.md. I agree with
    deferring it: closing it means shipping manifests for the remaining units,
    which the ratified design (FINAL.md B3.11) holds unaffected by this change.
    The payloads are usable regardless; only the verdict differs.

11. **The probe cache key expansion is correct.** Five variables are joined
    with `\u0001`, undefined values sentinel to `\u0000`. Two environments with
    the same PATH but different PYTHONHOME get distinct keys and distinct
    probe results (TEST 26 vs TEST 27). The Map grows by one entry per distinct
    environment tuple per process — bounded at 1 for the real CLI.

12. **The anchored version parse is fail-closed.** `/^(\d+)\.(\d+)$/` rejects
    `4.not-a-version`, `3.10.1`, `3`, and empty string. A shim printing garbage
    no longer satisfies the floor. The probe's own command
    (`import sys; print('%d.%d' % sys.version_info[:2])`) always produces `X.Y`
    for real CPython interpreters, so the anchored regex does not reject
    legitimate interpreters.

## Open questions for the implementer

1. The `status` text renderer does not print `skill.missing` (only
   `skill.problems` and `skill.runtime`). The `doctor` renderer does. This is
   pre-existing and not blocking, but the `explain` function added in this diff
   made the asymmetry more visible. Should `explain` also print `missing:` for
   consistency with `doctor`? Optional, cosmetic.

2. The 48d4fe4 follow-up commit (3 tests: ownership, core-never-unmanaged,
   laundering delta) is in the working tree but after the reviewed commit
   b180127. Should it be folded into b180127, or is it intended as a separate
   commit? The tests directly address the adversarial questions I was asked to
   investigate, and they pass. No issue either way — just noting the
   relationship.
