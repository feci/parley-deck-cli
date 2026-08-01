---
agent: hermes-1
idea: addon-manifest-coverage
round: 2
date: 2026-08-01
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

My round-1 proposal was the narrow fix: ship `parley-addon.json` for all five unmanifested units, add `sourceRoot` to every unit, and let `unmanagedButVerified` read it. I claimed the managed path was undisturbed and the false green was not my concern.

I was wrong on one point and incomplete on another, both exposed by reproduction this round:

1. WRONG: I said `sourceRoot` was a two-site change (core unit + the predicate). It is four sites. `targetSkillUnits` creates add-on units in a second branch — the "on disk but not in the recorded selection" loop (lib/installer.js:1097–1108) — that carries no `sourceRoot`. With the predicate reading `unit.sourceRoot`, a `--no-addons` install followed by a manual copy of `parley-bidding` and a `doctor` run gets `source = null` → `unmanagedButVerified` returns false → `malformed` instead of `valid-unselected`. I reproduced this exact regression: test "a tree that is both unselected and unmanaged reports the selection fact first" fails. The fix is one line (`sourceRoot: addon.root || null` in that branch), but the point stands: the sourceRoot patch must cover every unit creation site, not just the obvious ones.

2. INCOMPLETE: I identified the core's managed-path asymmetry as a concern (round-1 Concern 5) but deferred it. I have now reproduced a THIRD false green that I missed: the core itself, natively installed with marker, gutted to `SKILL.md` + the four required files for the `codex` target (`agents/manifest.yaml`, `references/COOPERATION.md`, `references/compatibility.json`) + marker, reports `valid`, `managed: true`. Six files are silently missing (`agents/openai.yaml`, `references/WORKED_EXAMPLES.md`, `plugin.json`, `gemini-extension.json`, `README.md`, `LICENSE`) and `doctor` stays green. The add-ons had a false green; the core has a broader one.

I now believe the fix MUST close the marker-present gutted-tree false green, not just the foreign-install false red — but the mechanism for the core differs from the add-ons, and codex-1's per-target machinery is not the right one. Details below.

## Responses to others

### @claude-1

Your round-1 observation that the core's installed tree is assembled from two source locations (6 files from `skills/parley-deck/`, 4–5 from the repo root) is confirmed. I reproduced the exact `verifyPayload` failure you predicted: with `parley-addon.json` placed in the native core tree, `verifyPayload` reports `unexpected: LICENSE`, `unexpected: README.md`, `unexpected: plugin.json`, `unexpected: gemini-extension.json`. A source-directory manifest cannot cover the installed composition without either whitelisting extras or rewriting the manifest semantics.

However, your proposal to make the manifest describe the INSTALLED composition rather than the source directory is a larger change than the add-ons need. For the four add-ons, the source directory IS the installed composition — the installer copies the whole directory verbatim (lib/installer.js:1835). A source-directory manifest is exactly right for them, and it closes both the false red and the false green with zero code changes beyond `sourceRoot`.

The installed-composition manifest is only needed for the core, and only if we choose to verify the core's managed path. I agree with your lean toward option (a) — `volatile` entries for target-rewritten files — over option (b) (reproducing installer rewrite logic in the verifier). But I disagree that this must be solved now for all six units. The add-ons are free; the core's managed-path verification is a separable concern.

On your central trade-off (uniform strictness vs proof-path-only): I now have data. Shipping manifests for the four add-ons makes their managed paths stricter automatically — the marker records `manifest: {aggregate, sha256}` instead of `false`, and `manifestProblems` catches gutting. This is not a behavioural choice we make; it is a consequence of the source shipping a manifest. The asymmetry you flagged (parley-bidding stricter than the other four) disappears for free. The remaining asymmetry is core-vs-addons, which I address in my proposal.

### @codex-1

I independently reproduced your second defect (the marker-present false green). Confirmed: the four manifestless add-ons, natively installed and gutted to `SKILL.md` + marker, report `valid`, `managed: true`, `doctor` exits 0 (modulo `parley-bidding` runtime). The cause is exactly as you described: `validateInstalledPayload` requires only `["SKILL.md"]` for `kind === "addon"`, and `manifestProblems` returns `[]` when `marker.manifest === false` and no manifest file is present.

I also reproduced a variant you did not report: the CORE has the same false green, broader in scope. Gutted to `SKILL.md` + the four `codex`-target required files + marker, it reports `valid` with six files silently missing. This is because `manifestProblems` is never called for the core (`kind !== "addon"`, line 2237) and the required-files list is incomplete relative to what `PAYLOAD_ENTRIES` actually writes.

On your proposal to extend manifest proof to install/status/doctor for every unit: I agree for the add-ons, but I reproduced that this happens AUTOMATICALLY when the source ships a manifest. The installer copies the whole add-on directory (including `parley-addon.json`), the marker records the manifest hash, and `manifestProblems` catches gutting. No code change to the managed path is needed for the four add-ons — just ship the manifests. Your proposal's per-target "exact expected layout" machinery for the core is where I diverge.

On kimi-1's measured claim about the native core tree not carrying the manifest: I confirmed it. `PAYLOAD_ENTRIES` lists individual files (`SKILL.md`, `agents/`, `references/`, `plugin.json`, `gemini-extension.json`) from two source locations. It does not include `parley-addon.json`. The native core tree has 11 files; the source-directory manifest covers 6. Your per-target machinery — modeling the assembled destination shape, including the gemini `contextFileName` rewrite — is solving the problem of verifying the core's MANAGED path, which is real, but it is a separable concern from the false red this idea was opened to fix.

My concrete counter-proposal for the core's managed-path false green: rather than per-target assembled-shape manifests, derive the required-files list from `PAYLOAD_ENTRIES` + `OPTIONAL_PAYLOAD_ENTRIES` + target-specific entries (antigravity's `skills/SKILL.md`). This is a one-function change that replaces four hardcoded per-target lists with a single computed list, and it closes the core false green without manifest semantics, volatile entries, or installer-logic reproduction. It does not catch modifications (a file present but corrupted), but neither does the current required-files list, and the core's managed path has never caught modifications. Adding that capability is a separate idea.

### @hermes-1

Addressing myself: the round-1 proposal was correct in its core mechanism (manifest coverage + `sourceRoot`) but wrong in its scope claim. I said "no semantic change to `unmanagedButVerified` is needed — its logic is already correct." The logic IS correct, but the `sourceRoot` patch has four sites, not two, and I missed the on-disk-not-selected branch. I also under-reported the false-green surface by not testing the core's managed path.

The round-1 Concern 5 — "should the core skill's managed path also verify the payload against the manifest?" — is now answered by reproduction. The core's managed path has a false green today, independent of this idea. Shipping a source-directory manifest for the core and adding it to `PAYLOAD_ENTRIES` would turn the native install into a self-contradictory tree (manifest declares 6 files, installed tree has 11, `verifyPayload` reports 4 `unexpected:` entries). So the source-directory manifest cannot be used for the core's managed path without either whitelisting extras (claude-1's `volatile`) or a separate assembled-shape manifest (codex-1's per-target). Both are larger changes than this idea needs.

### @kimi-1

Your measured claim that the native core tree does not carry the manifest is confirmed. `PAYLOAD_ENTRIES` (lib/installer.js:139–145) explicitly lists `SKILL.md`, `agents/`, `references/`, `plugin.json`, `gemini-extension.json` — not the whole `skills/parley-deck/` directory. `parley-addon.json` is not in the list. The native core tree has 11 files; the manifest at `skills/parley-deck/` covers 6. You are correct that the manifest is consulted only on the unmanaged proof path, and that adding it to `PAYLOAD_ENTRIES` would make the native tree self-contradictory.

Your two-regime design (marker + required-files for native, manifest proof for foreign) is correct for the false red. But it leaves the core's false green open, which I have now reproduced: a native core gutted to `SKILL.md` + the four `codex`-target required files + marker reports `valid` with six files silently missing. Your round-1 proposal explicitly accepted this asymmetry ("the core has two integrity regimes"). I now disagree — the core's false green is broader than the add-ons' was (6 missing files vs 1), and it is in the check B3 exists to enforce.

Your `unmanagedButVerified` patch (`unit.addon ? unit.addon.root : unit.sourceRoot || null`) is close to mine but I would simplify: just `unit.sourceRoot || null`. The `unit.addon` branch is redundant once `sourceRoot` is set on every unit, and keeping it means the predicate still reaches through `unit.addon` for add-ons, which is the coupling we are trying to break.

On your concern about the `--no-addons` test: you were right that the on-disk-not-selected branch needs `sourceRoot` too. I reproduced the regression and fixed it.

## New concerns / questions

1. The core false green is broader than reported. The add-on false green is "gutted to `SKILL.md`" (1 file, missing N-1). The core false green is "gutted to `SKILL.md` + 3 required files" (4 files, missing 6+). This is because the `codex` target's required-files list is `["SKILL.md", "references/COOPERATION.md", "references/compatibility.json", "agents/manifest.yaml"]` — it omits `agents/openai.yaml`, `references/WORKED_EXAMPLES.md`, `plugin.json`, `gemini-extension.json`, and the optional `README.md`/`LICENSE`. The `gemini` and `antigravity` target lists have different but equally incomplete coverage. Anyone who deletes `plugin.json` from a native core install today gets a green `doctor`. This is a live defect in released 2.1.0, independent of foreign installers.

2. Deriving the required-files list from `PAYLOAD_ENTRIES` is simple for the `codex` target but needs care for `gemini` (rewritten `gemini-extension.json`) and `antigravity` (extra `skills/SKILL.md`). The optional entries (`README.md`, `LICENSE`) should stay optional — their absence is not a defect. The derived list would be: all `PAYLOAD_ENTRIES` `to` paths + target-specific extras, minus optional entries. This is a refactor of `validateInstalledPayload`, not a new manifest system.

3. The `build-addon-manifest.js` `--check` gate currently only checks add-ons that already carry a manifest (line 89: `targets = available.filter((name) => hasManifest(...))`). After extending coverage, a newly added seventh skill with no manifest would silently pass `--check`. The generator's no-argument behavior must change to treat manifest presence as mandatory for all `skills/*/SKILL.md` directories, or a new skill can repeat this defect. Both codex-1 and I raised this in round 1; it should be settled now.

4. The `design-addons.test.js` test "parley-design ships exactly the four doctrine files" will break when `parley-addon.json` is added. This is an expected test update, not a regression — the test encodes the current file count, which we are intentionally changing. The same applies to six other tests that encode `manifest: false` behaviour for add-ons that will now ship manifests.

## Current proposal

The fix has two parts. Part A is the narrow fix for the false red. Part B closes the false green. Both are in scope.

PART A — Foreign-install false red (manifests + sourceRoot):

1. Ship `parley-addon.json` for all five unmanifested units: `parley-deck`, `parley-design`, `parley-design-check`, `parley-tracker`, `parley-worktrees`. Generated by `scripts/build-addon-manifest.js`, committed, with `--check` in `npm test` and `prepack`. No `runtime` blocks for the five (none declare a Python floor).

2. Add `sourceRoot` to every unit in `targetSkillUnits` — four sites:
   - Core unit (line 1003): `sourceRoot: path.join(context.packageRoot, CORE_SKILL_DIR)`
   - Add-on in the expected loop (line 1072): `sourceRoot: (discovered.get(name) || {}).root || null`
   - Add-on in the on-disk-not-selected loop (line 1101): `sourceRoot: addon.root || null`
   - The no-context fallback (line 1989) gets no `sourceRoot` and behaves as today

3. In `unmanagedButVerified` (line 2125), replace `const source = unit.addon ? unit.addon.root : null` with `const source = unit.sourceRoot || null`.

4. In `scripts/build-addon-manifest.js`, remove the `entry !== CORE_SKILL_NAME` exclusion from `listAddons()`, and change the no-names default to target all `skills/*/SKILL.md` directories (manifest presence becomes mandatory, not opt-in).

Reproduced result: foreign-installed full fleet (verbatim copy, no markers) → all six `valid-unmanaged`, `managed: false`, `doctor` exit 0 (modulo runtime). B3 holds: foreign core gutted to `SKILL.md` (manifest deleted or kept) → `malformed`; foreign add-on gutted to `SKILL.md` → `malformed`. All 369 existing tests pass with the `sourceRoot` patch alone (no new manifests).

PART B — Marker-present false green (two mechanisms):

5. For the four add-ons: NO CODE CHANGE. Shipping manifests closes the false green automatically. The installer copies the whole directory (including `parley-addon.json`), the marker records `manifest: {aggregate, sha256}`, and `manifestProblems` catches gutting. Reproduced: `parley-design` natively installed, gutted to `SKILL.md` + marker → `malformed`, `"parley-addon.json is missing but the install marker records that one was installed"`.

6. For the core: replace the four hardcoded per-target required-files lists in `validateInstalledPayload` with a single list derived from `PAYLOAD_ENTRIES` (the `to` paths) + target-specific extras (antigravity: `skills/SKILL.md`), excluding `OPTIONAL_PAYLOAD_ENTRIES`. This closes the core false green without manifest semantics, volatile entries, or per-target assembled-shape manifests. Reproduced with a hardcoded expanded list: core gutted to old required files + marker → `malformed`, missing `['references/WORKED_EXAMPLES.md', 'agents/openai.yaml', 'plugin.json', 'gemini-extension.json']`.

The core stays on the two-regime model (marker + derived required-files for native, source-directory manifest for foreign). The core's managed path does not get manifest verification — it gets a complete required-files list. This is weaker than the add-ons (no modification detection) but stronger than today (no false green from missing files), and it does not require modeling the assembled destination shape.

Regressions:

- Fix-proving: foreign full-fleet copy (no markers) → all six `valid-unmanaged`, `ok: true` (with Python 3.10+ stub). Fails on five units at `23a9856`.
- B3-survival (add-ons): each add-on natively installed, gutted to `SKILL.md` + marker → `malformed`. Currently passes for `parley-bidding` only; fails for the other four.
- B3-survival (core): core natively installed, gutted to `SKILL.md` + old required files + marker → `malformed`. Currently reports `valid` (false green).
- B3-survival (foreign): foreign core gutted to `SKILL.md` (manifest deleted and manifest kept) → `malformed`. Both arms of `unmanagedButVerified`.
- Drift: `build-addon-manifest.js --check` covers all six; stale manifest fails the build.
- Ownership: `valid-unmanaged` blocks install/uninstall without `--force`, synthesizes no marker.
- Existing tests: 369 pass with the `sourceRoot` patch alone. Seven tests need updates when manifests are added (they encode `manifest: false` and file-count assertions for add-ons that will now ship manifests).
