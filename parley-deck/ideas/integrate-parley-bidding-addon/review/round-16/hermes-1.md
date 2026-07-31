---
idea: integrate-parley-bidding-addon
review-round: 16
agent: hermes-1
date: 2026-07-31
reviewed-commit: a49d68f
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from rounds 14 and 15 are addressed at a49d68f. I verified each closure by reading the code, running the tests, and reproducing the edge cases.

- **Nested destinations (codex-1 round-15 MAJOR): CLOSED.** `aliasedDestinations` (installer.js:1427-1457) now computes a resolved physical path for every planned unit via `resolvedDestination` (installer.js:1374-1396) and rejects any pair where one contains the other, in both nesting directions. I reproduced the exact scenario from codex-1's round-15 report — `CODEX_HOME` set to kimi's planned core destination — and the plan is now refused with `ok: false`, zero writes, and an "overlaps another in this plan" message. I also tested the reverse nesting direction (same result), and verified that uninstall rejects the same configuration with nothing deleted. The regression at test:1655 fails at `26478e9` and passes at `a49d68f`.

- **Uninstall dry/real per-unit agreement (codex-1 round-15 MINOR, kimi-1 round-15 MINOR, hermes-1 round-15 MINOR-1): CLOSED.** The separate real-only preflight block that cycle 18 deleted from install was deleted from uninstall in cycle 19 (the 38-line removal at the top of the cycle-19 diff). `uninstallCommand` (installer.js:649-685) now goes directly to `removeFleetAtomically` for both dry and real, exactly as `installCommand` goes to `installFleetAtomically`. I reproduced my round-15 MINOR-1 scenario — an aliased pair (codex/hermes symlinked) plus a corrupted marker on qwen — and verified that dry and real produce identical per-unit `ok`, `action`, and `message` for every unit. The regression at test:1684 fails at `26478e9` and passes at `a49d68f`.

- **Case-only spellings in one plan (codex-1 round-15 MAJOR, test quality): CLOSED.** The test at test:1628 was rewritten to put two targets (`CODEX_HOME` and `KIMI_CODE_HOME`) with case-only spellings of one directory into a single `install --target all --include-undetected` plan. I confirmed it fails at `d7ab1c3` (the commit it was written to discriminate) and passes at `a49d68f`.

- **Fourth manifest reader in run-python-tests.js (codex-1 round-15 MINOR): CLOSED.** `declaredPythonFloor` (run-python-tests.js:53-61) now calls `addonManifest.readManifest(ADDON_ROOT)` instead of `readFileSync` + `JSON.parse`. A non-readable or symlinked manifest now fails the runner rather than silently becoming "no declared floor." I verified the runner reports the manifest error and exits non-zero when given a symlinked manifest.

- **CHANGELOG Known limits (round-14 ruling): CLOSED.** The "Known limits" heading (CHANGELOG.md) carries the single-writer warning verbatim. The concurrency follow-up does not gate 2.1.0.

## New findings

### NIT-1: Uninstall per-action JSON does not carry the `dryRun` flag that install's does

`installCommand` (installer.js:634-636) propagates the per-unit `dryRun` flag to the per-action wrapper:

```js
...(core.dryRun ? { dryRun: true } : {}),
```

The comment at line 634 explicitly states this flag is "part of the JSON contract." `uninstallCommand` (installer.js:670-677) does not — the per-action wrapper omits the field entirely. The per-unit result inside `removeFleetAtomically` does carry `dryRun: true` (installer.js:1641), but the per-action wrapper that a JSON consumer reads does not propagate it.

Measured: a dry-run uninstall of a healthy single-target fleet produces `actions[0].dryRun: undefined` where a dry-run install of the same shape produces `actions[0].dryRun: true`. The top-level `result.dryRun` is correct for both (`true` for dry, `false` for real).

This is pre-existing (the uninstall path never carried the flag, even before the cycle-19 unification) and cosmetic: no behavior is affected, no test expects it, and the top-level field is correct. A JSON consumer that checks `action.dryRun` rather than `result.dryRun` would see a difference between install and uninstall dry-runs, but both commands' top-level `dryRun` agree. Filing as NIT because the install code's own comment calls the flag part of the JSON contract, and the asymmetry is unannotated.

Does not gate 2.1.0. It can ship as a follow-up: add `...(core.dryRun ? { dryRun: true } : {})` to the uninstall action builder at installer.js:676, matching install.

No other new findings.

## Release judgement

**Releasable as 2.1.0.** Nothing I found must change before release.

The three cycle-19 fixes are complete and their regressions discriminate. The overlap check is computed at planning time on the resolved physical path, and I verified that the relationship between two destinations cannot change between planning and commit in a single process: staging creates only real directories (never symlinks), creating intermediate directories does not change the resolved path (because `realpath` of a newly-created directory is deterministic from its parent), and dangling symlinks are blocked by `destinationAncestorObstacle` before any staging creates their targets. The single result path for uninstall now agrees with install's on every predictable disposition. No fifth manifest reader or validator bypass exists.

NIT-1 is a cosmetic JSON field asymmetry with no behavioral impact.

## What I verified

1. **Test suite.** `npm test` at `a49d68f`: 355/355 node tests pass, 0 fail. Python leg 54/54 on Python 3.14.6; refuses 3.9.6 by design (floor >=3.10). Manifest check: 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since `714712f`. Node tests also pass with `/usr/bin/python3` 3.9.6 first on PATH (355/355). `npm pack --dry-run --json`: 202 files, 48 under `skills/parley-bidding/`, no `__pycache__`/`.pyc`/`.pyo`.

2. **Cycle-19 diff.** Read `git diff 26478e9..a49d68f` — three files: lib/installer.js (overlap check + `resolvedDestination`/`overlaps` functions, uninstall preflight deletion), scripts/run-python-tests.js (shared manifest parser), test/bidding-addon.test.js (three rewritten/new regressions). `git diff --check` clean for both `26478e9..a49d68f` and `49fc3ec..a49d68f`.

3. **Regression discrimination.** Copied `26478e9` and `d7ab1c3` to temporary archives via `git archive`, dropped the cycle-19 test file in, and ran the three new regressions:
   - Nested destinations: **fails at `26478e9`**, passes at `a49d68f`.
   - Uninstall dry/real per-unit agreement: **fails at `26478e9`**, passes at `a49d68f`.
   - Case-only spellings in one plan: **fails at `d7ab1c3`**, passes at `a49d68f`.

4. **Question 1 — can the relationship between two destinations change between planning and commit?** Traced `resolvedDestination` (installer.js:1374-1396): it walks to the nearest existing ancestor, takes its `realpath`, and appends the not-yet-created tail. The final resolved path is invariant to which ancestor is "nearest" because `realpath` of a newly-created directory is `realpath(parent) + "/" + name`. Staging (`copyPayloadAtomically`, installer.js:1713-1771) creates only real directories via `mkdirSync({recursive:true})` — never symlinks. Commit (`commitStagedUnit`, installer.js:1775-1789) uses `renameSync`, which moves real directories. Quarantine (`removeFleetAtomically`, installer.js:1646-1691) also uses `renameSync`. No phase of install or uninstall creates a symlink. I tested the most dangerous scenario I could construct — a symlink to a non-existent target whose target would be materialized by staging (`/home/link -> /home/real`, `CODEX_HOME=/home/real`, `KIMI_CODE_HOME=/home/link`) — and `destinationAncestorObstacle` (installer.js:2289-2314) blocked the dangling-symlink unit in preflight before any staging ran, producing `ok: false` with zero writes. When the symlink target exists at planning time, `resolvedDestination` resolves through it and the paths are equal, caught by the overlap check. **The relationship cannot change in a single process.**

5. **Question 2 — does the single result path for uninstall agree with install's on every disposition?** Traced `removeFleetAtomically` (installer.js:1585-1711) and `installFleetAtomically` (installer.js:1470-1579) unit by unit. Both use `aliasedDestinations` for the overlap check, both apply preflight conditions, both skip remaining units when any unit is blocked, both use a dry-run branch that records future-tense actions without touching the filesystem. The dispositions: `blocked` (alias/preflight), `missing` (absent, uninstall-only, `ok: true`), `skipped` (fleet gate), `remove`/`removed` (dry/real success), `install`/`replace`/`installed`/`replaced` (install dry/real). The dry/real naming difference (`remove` vs `removed`) is by design and parallels install's (`install` vs `installed`). I verified the alias-plus-separate-blocker scenario from my round-15 MINOR-1: dry and real now produce identical per-unit `ok`, `action`, and `message` for all 14 targets. I also verified the no-blocker all-present case: dry says `remove`, real says `removed`, both `ok: true` — matching install's convention. Quarantine/staging/commit failures are inherently unpredictable in dry-run (dry-run doesn't perform the operation), which is by design in both paths. The per-action `dryRun` flag is the only asymmetry (NIT-1).

6. **Question 3 — fifth manifest reader or validator bypass?** Traced every manifest read in the repository. All production reads go through `readManifest` or `verifyPayload`: installer.js `writeMarker` (line 1844), `runtimeAvailability` (line 2041), `manifestProblems` (line 2205 via `verifyPayload`); run-python-tests.js `declaredPythonFloor` (line 54, fixed in cycle 19); build-addon-manifest.js (line 103). The raw `readFileSync` at build-addon-manifest.js:119 is a byte-for-byte string comparison in `--check` mode, guarded by `hasManifest` (which uses `lstatSync` and rejects symlinks), not a parse — its purpose is detecting stale manifests, not consuming manifest policy. `rewriteStagedGeminiManifest` (installer.js:860-871) reads `gemini-extension.json`, not `parley-addon.json`. The Python `manifest.py` in `skills/parley-bidding/scripts/` is the add-on's own tool for release archives, not an installer manifest reader. **No fifth reader or bypass exists.**

7. **Uninstall overlap check.** Verified `removeFleetAtomically` (installer.js:1589) calls `aliasedDestinations`, so the containment check applies to uninstall as well as install. Reproduced nested destinations in uninstall: both units blocked, zero deletions, both dests intact on disk.

8. **Add-on source stability.** `git diff 714712f..a49d68f --stat -- skills/parley-bidding/` is empty — the add-on payload is unchanged since `714712f`.
