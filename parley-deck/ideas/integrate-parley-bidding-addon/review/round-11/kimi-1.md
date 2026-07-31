---
idea: integrate-parley-bidding-addon
review-round: 11
agent: kimi-1
date: 2026-07-31
reviewed-commit: 12f9071
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

All four round-10 findings are **closed**, each re-measured by me at `12f9071` as uid 501
(scratch copy of the tree at `/tmp/kimi-r11`, library-driven synthetic HOMEs, harness
`r11-experiments.js` / `r11-experiments2.js` left in place there):

- **Install replace, disposal obstacle in the last target** (codex-1/kimi-1 MAJOR) — one
  mode-000 subdirectory buried in an owned destination, and separately `chmod -R a-w` on the
  whole owned destination: both now end preflight with **0 units written**, `ok: false`
  (my R1). The suite's two regressions cover the same arms.
- **Uninstall, one target, one frozen owned add-on** (kimi-1 MAJOR) — **0 units removed, core
  intact** (my R2). The round-4-forbidden end state is no longer reachable through mode bits.
- **Uninstall fleet, foreign marker in the last of fourteen targets** (hermes-1 MAJOR) — **0
  units removed** (my R3). The fleet-wide preflight `installCommand` got in cycle 10 now exists
  in `uninstallCommand` (`preflightUninstallUnit`, installer.js:725), and I confirmed by reading
  that it runs over every unit of every target before the first deletion, ownership unless
  `--force`, removability always.
- **Marker: delete only `markerSchema`, keep `manifest`, tamper payload** (codex-1 MAJOR) — now
  `status: malformed` with `install marker records a manifest but declares no markerSchema`
  (my R4). The genuine 2.0.0 shape (neither field) stays `valid` via the suite's seventh
  regression.
- **kimi-1's round-10 NIT** (the 3.9.6/3.14 wording) — fixed; both IMPLEMENTATION.md entries
  (cycles 11 and 13) now state the Python leg refuses 3.9.6 by design. I verified the refusal
  directly: `PATH` with `/usr/bin` first prints `python3 is 3.9, but the add-on declares >=3.10`
  and stops.

The "0/7 at `9ed2081`, 7/7 at `12f9071`" discrimination claim I validated independently: the
cycle-14 test file run against a `git archive` of `9ed2081` fails exactly the seven new tests
(77 run, 70 pass, 7 fail), and the full gate at `12f9071` passes 332/332 node, 54/54 Python on
3.14.6, manifest check ok (47 files, aggregate unchanged since `714712f`).

Deferred items are unchanged and correctly recorded: only `parley-bidding` ships a
`parley-addon.json` (B3.11), and my round-9 `dirExists` discovery-guard NIT stands deferred per
the cycle-13 entry.

## New findings

**One — a residual in the disposal class, measured, which I judge a recorded follow-up rather
than a release gate.**

The sixth door exists, and it is narrower than the five before it. `firstRemovalObstacle`
decides disposability with `accessSync` mode checks. On macOS there are disposal obstacles
`access(2)` cannot see, and through them the forbidden state of rounds 4 and 10 is reachable
again at `12f9071`, with no flag:

| arm (all at `12f9071`, uid 501) | scenario | measured |
|---|---|---|
| uninstall, one target | `chflags uchg` on **one file** inside an owned add-on (Finder "Locked" on a single file) | preflight passes; **core and all four sibling add-ons removed**, then `parley-bidding: failed` (my E3) |
| uninstall, one target | `chflags uappnd` on the add-on destination directory | same shape: 5 removed, unit failed (my E5) |
| uninstall, one target | `chmod +a '<user> deny delete'` on the add-on destination directory | same shape, EACCES (my E6) |
| uninstall, one target | one `uchg` file inside the **core** destination | core **gutted to the locked file**, then `failed`; remaining units removed (my E8) |
| install, replace (same locked file) | any target | **absorbed**: unit `replaced`, `ok: true`, warning + `.bak` debris — the cycle-14 post-commit guard working as designed (my E4/E9) |

Why I do not gate 2.1.0 on this, stated plainly:

1. **Every realistic freeze idiom is closed, measured.** `chflags -R uchg` on a destination —
   the macOS mirror of round 10's `chmod -R a-w` — is **caught**: clean fleet-wide block, 0
   writes and 0 removals (my E1/E2). So is a Finder lock on the *folder*. The reason is
   accidental but reliable: macOS `access(2)` returns EPERM on a flag-locked **directory**
   (primitive probe), so directory-carrying flags land inside the walk even though the code
   never reads flags. What escapes is selective per-*file* flagging, `uappnd`, and delete-denying
   ACLs — states with no common producer; making one requires deliberate per-entry surgery
   inside a machine-managed directory, the same category of act as deleting files mid-install,
   which no preflight can ever bound.
2. **The path 2.1.0 users actually run is already absorbed.** Install with the same locked file
   commits correctly and reports a warning (E4/E9). After an E3-style partial uninstall,
   `doctor` fails honestly, and re-install without `--force` blocks on ownership with the
   accurate "re-run with --force" message (the partial `rm` eats the marker); the forced
   re-install then commits with a warning (measured).
3. **A complete stdlib fix does not exist.** Node exposes no `st_flags`. The single-locked-file
   arm *is* detectable — `accessSync(file, W_OK)` returns **EPERM** for a `uchg` file vs
   **EACCES** for an ordinary 0444 file (measured, my P1) — so a ~3-line mitigation in
   `firstRemovalObstacle` (check W_OK on non-directory entries, treat EPERM only as an
   obstacle) would close the Finder-reachable arm. But `uappnd` directories and delete-denying
   ACLs pass `access(2)` entirely; closing those needs platform spawns (`find -flags`, `ls
   -lO`), and Linux `chattr +i` / Windows ACLs are the same question again. That is a
   proportionality decision the project should take with this evidence on record — as a
   follow-up with a known-limit note in IMPLEMENTATION.md — not a "looked complete and was not"
   gate, because cycle 14's mechanism (an `access(2)` walk) is stated in the code and the
   record, and its boundary is exactly this.

If the consensus bar is instead "any reachable forbidden state gates the release", converting
this to a cycle-15 finding is mechanical: the harness reproduces E3/E5/E6/E8 on demand.

Two boundary observations, recorded for completeness, neither a finding:

- The walk is **over-strict in the fail-closed direction**: an *empty* 0555 directory is
  blocked ("cannot be emptied") although `rmSync` would succeed (rmdir needs only parent
  write). Reachable only under `--force`; the message is accurate about the remedy (my P2).
- Within a single unit, a failed `rmSync` is inherently partially-destructive (E8's gutted
  core). No preflight can make `rm` transactional; preflight accuracy is the only guard, which
  is what makes the residual above worth recording rather than waving at.

## Release judgement

Releasable as 2.1.0. The create/touch/dispose triangle is closed for everything `access(2)`
can see — mode-bit obstacles at any depth, in any subtree, for both commands, fleet-wide — and
the two deferred items stand recorded. The one thing I would ask before the tag is not code:
record the flag/ACL residual above as a known limit in IMPLEMENTATION.md (with the EPERM-vs-
EACCES mitigation noted as an option), so 2.1.0 ships with the boundary documented rather than
discovered.

## What I verified

- Read end-to-end at `12f9071`: `lib/installer.js` (2187 lines), `lib/addon-manifest.js`,
  `bin/parley-deck-skill.js`; the diffs `9ed2081..12f9071` (cycle 14) and `49fc3ec..12f9071`
  (full scope, incl. the `design-addons.test.js` grammar-sentinel assertions — test-only,
  consistent with round-1 codex-1 MINOR).
- Working tree untouched throughout (`git status` clean, HEAD `12f9071` before and after); all
  experiments ran in `git archive` scratch copies under `/tmp`.
- Full gate in the scratch copy: **332 node tests, 0 fail**; **Python 54/54 on 3.14.6**;
  `/usr/bin/python3` 3.9.6-first PATH → graceful refusal (`>=3.10` declared), confirming only
  the node leg is measured on both interpreters; manifest check ok, 47 files, aggregate
  `sha256:7854adf1…b95a6d`, unchanged.
- Independent discrimination run: cycle-14 tests against `9ed2081`'s installer → exactly the
  7 new regressions fail (0/7), against `12f9071` → 7/7.
- Round-10 arms re-measured closed (R1–R4 above); round-4 forbidden state unreachable through
  mode bits on either command.
- Disposal-class boundary mapped by primitive probe and nine end-to-end scenarios: flagged
  directories caught (access EPERM), flagged files/`uappnd`/deny-ACLs invisible; install arm
  absorbed by the post-commit guard; uninstall arm partial (the finding above).
- Claims in IMPLEMENTATION.md cycles 10–14 and `review/round-09/VOID.md` read and spot-checked
  against my own measurements; the two corrected NIT entries verified in place. I did not read
  the other round-11 reviews before writing this file.
