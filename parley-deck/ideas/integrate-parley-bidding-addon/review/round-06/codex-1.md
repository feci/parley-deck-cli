---
idea: integrate-parley-bidding-addon
review-round: 6
agent: codex-1
date: 2026-07-30
reviewed-commit: c673111
---

## Verdict
BLOCK

## Outstanding findings — closed or not

1. **Ownership consistency — NOT FULLY CLOSED.** A marker naming a different skill is now
   `malformed`, and the added regression covers that value. A readable marker that omits
   `skill` entirely still reports `valid` and `managed:true` while install and uninstall
   refuse ownership. See the MAJOR finding.
2. **Read selectors replacing the recorded selection — CLOSED.** On a healthy full install,
   `--only parley-bidding` and `--no-addons` narrow `doctor`, `status`, and `paths` to exactly
   the requested units; no inspected unit becomes `valid-unselected`. Unflagged reads still
   surface genuine residual directories.
3. **Interpreter probe cache omitted effective cwd — CLOSED BY MEASUREMENT.** With the same
   environment and relative `PATH=bin`, working directories whose shims reported 3.12 and
   3.9 produced the correct pass/fail answers in both call orders. `spawnSync` receives the
   resolved cwd and the cache key distinguishes it.
4. **Validation record overstated Python 3.9 coverage — CLOSED.** The record now separates
   the 309-test Node leg from the Python package gate. I reproduced 309/309 Node tests with
   `/usr/bin/python3` 3.9.6 first on `PATH`; the Python runner then exited 1 at the declared
   floor and ran zero tests. On Python 3.14.6, full `npm test` passed 309/309 Node and 54/54
   Python tests.

## New findings

### [MAJOR] An omitted marker skill still makes health claim ownership that mutations deny

**Where:** `lib/installer.js:1342`, `lib/installer.js:1632-1639`,
`test/bidding-addon.test.js:936-960`

**What:** `skillUnitStatus` reports an identity problem only when `marker.skill` is defined
and unequal:

```js
state.marker.skill !== undefined && state.marker.skill !== unit.skill
```

`installerOwnsDestination`, used by install and uninstall, requires exact equality. Therefore
a readable marker with the correct package name but no `skill` property is healthy and
managed to `doctor`, yet unowned to both mutation commands. The new regression changes the
field to another string; it does not test the missing-identity branch.

**Why it matters:** The cycle claims one ownership answer, but a one-field deletion still
creates the same false-green/upgrade-refusal contradiction as round 5. This affects both the
core and add-on units. It is not needed for released-marker compatibility: `v1.0.0`,
`v1.4.0`, and `v2.0.0` all wrote a `skill` identity.

**Evidence:** In an isolated checkout of `c673111`, I performed a normal full install and
deleted only `skill` from the installed core and bidding markers. Payloads, manifests, and
all other marker fields remained unchanged. `doctor.ok` was `true`; both units were
`status:"valid"`, `managed:true`, with no problems, and bidding remained runtime-available.
Both `install --only parley-bidding` and `uninstall --only parley-bidding` then returned
`action:"blocked"` for the core and bidding units.

**Fix:** Require `state.marker.skill === unit.skill` for every readable installer-owned
marker; do not exempt `undefined`. Prefer one shared parsed-ownership predicate for health
and mutations. Add core and add-on regressions that delete `skill` and assert `malformed`,
`managed:false`, failed health, and matching install/uninstall refusal.

## What I verified and found correct

- I reviewed a clean disposable clone at
  `c673111579aa43cf86872647efd3bd07d71c6043`; the source worktree's pre-existing
  `test/bidding-addon.test.js` modification was untouched.
  `git diff --check 3634cc8..c673111` passed.
- With the source checkout's dependency directory supplied read-only through `NODE_PATH`,
  `PYTHONDONTWRITEBYTECODE=1 npm test` passed on Python 3.14.6: **309/309 Node** and
  **54/54 Python**, followed by the 47-file manifest check at
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- `valid-unmanaged` behaves as designed for an intact bidding tree with its marker removed:
  `managed:false`, no integrity problem, runtime still probed, and healthy `doctor`; install
  and uninstall remain fail-closed without `--force`.
- `valid-unselected` preserves the payload verdict: a still-intact, installer-owned tracker
  left outside a later recorded selection was `valid-unselected`, `managed:true`,
  `selected:false`, with no missing files, and made `doctor` fail. Re-including it restored
  `valid`.
- Unflagged `doctor`, `status`, and `paths` all surfaced the same residual
  `valid-unselected` unit. Explicit read filters narrowed all three commands without accusing
  other recorded add-ons.
- A marker whose `skill` is present but names another unit is now `malformed` and unowned, as
  intended. The remaining defect is specifically the omitted property.
- A focused code-only structural graph of `installer.js` and `addon-manifest.js` confirmed
  that health and mutation ownership are separate consumers of `readMarkerState`; direct
  source inspection and runtime probes established the finding.
- The deliberately deferred manifests for the other add-ons were not reopened.

## Open questions for the implementer

None. The missing identity must receive the same answer as a wrong identity.
