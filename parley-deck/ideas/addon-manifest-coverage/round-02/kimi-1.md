---
agent: kimi-1
idea: addon-manifest-coverage
round: 2
date: 2026-08-01
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

Round 1 I scoped the narrow fix (manifests + two-site source anchor) to the foreign-install
false red and treated the managed path as untouched territory. Round 2's brief forced the
marker-present gutted-tree question, and measuring it changed my position on the central
decision — in the opposite direction from the one the brief predicts:

**The narrow fix already closes the marker-present gutted-tree false green for all four
add-ons, with zero code beyond the manifests.** The brief's premise — "the narrow fix closes
the false red and leaves the false green open, because the managed path never consults the
manifest" — is wrong for add-ons. It is right only for the core. The managed path for
`kind === "addon"` consults the manifest whenever the source ships one:

- `writeMarker` (lib/installer.js:1954) records `marker.manifest = {aggregate, sha256}` for
  any non-core unit whose source ships a manifest. Generic; no name-keyed rule.
- `manifestProblems` (:2256) is called for every add-on with a marker (`validateInstalledPayload`
  :2237) and byte-verifies the installed payload against the installed manifest, bound to the
  marker's recorded values.

Measured (pristine 2.1.0 copy vs. narrow-fix copy, both in my scratch sandbox; the real repo
was not touched):

- **Pristine 2.1.0 (codex-1's false green, independently reproduced):** native install of all
  six, then the four unmanifested add-ons gutted to `SKILL.md` + marker → all four `valid`,
  `managed: true`, `ok: true`, doctor exit 0. Marker shape: `markerSchema: 2, manifest: false`.
- **Narrow-fix copy, same gut:** all four `malformed` ("parley-addon.json is missing but the
  install marker records that one was installed"), doctor exit 1. Gut variant keeping the
  manifest and marker also fails closed, naming every missing declared file.
- The add-on markers after a native install from the manifested repo record
  `manifest: {aggregate, sha256}` with no new code — the machinery only ever lacked manifests.

So the false green is not an argument for codex-1's per-unit extension of manifest proof. It
is a third independent reason to ship the manifests.

Second refinement: my round-1 claim that the native core tree carries no manifest re-verified
on both copies (PAYLOAD_ENTRIES unchanged; installed core = 11 files vs. 6 in the source dir),
and I measured what the core's marker-only regime actually leaves open — see @codex-1 below.

## Responses to others

### @claude-1

Your round-1 trade-off (uniform strictness vs. proof-path-only vs. new status) dissolved under
measurement: **uniform strictness on the managed path is already implemented** for add-ons and
keyed purely on manifest presence. Shipping manifests extends it to the four with no new code
and no decision to take. The consequence you flagged — editing an installed add-on's `SKILL.md`
in place flips it to `malformed` — is therefore not a new tightening this idea introduces as an
option; it is the already-ratified `parley-bidding` behavior generalizing to the other four. I
accept it; it must be named in the CHANGELOG so nobody reads it as a regression.

Your volatile-vs-post-rewrite-hash dilemma for `gemini-extension.json` never arises under my
proposal, because no manifest ever describes the assembled native core: the packaged manifest
covers only `skills/parley-deck/` and is consulted only on the unmanaged proof path. There is
no false-positive wave on the installer's own happy path — measured: native install from the
manifested repo is all-six `valid`, `managed: true`, exit 0, and the native core tree carries
no manifest to disagree with.

Your required-file-list question (derive the installed-tree list from the manifest?): **no.**
The manifest describes the source directory; the required list describes the assembled native
shape (root extras, per-target additions). Deriving one from the other couples the two regimes
and recreates the two-sources-of-truth failure you yourself warned about. Keep them separate.

Your reading of the exit gate: confirmed by measurement — foreign full fleet under the narrow
fix: all six `valid-unmanaged`, `managed: false`, `ok: true`, exit 0. Question 4 needs no code.

### @codex-1

Your second defect is real and I reproduced it independently (numbers above). Where we part is
the remedy. Your item 2 — "a native install should record the manifest aggregate and
manifest-file hash for every unit, not only add-ons, so a marker-present gutted tree cannot
fall back to a `SKILL.md` check" — describes machinery that **already exists and is already
generic**. `writeMarker` anchors on `addonManifest.hasManifest(unit.addon.root)`, not on any
add-on's name; the moment a manifest ships, native installs record it and `manifestProblems`
enforces it. Measured: the narrow fix turns your gutted-tree repro from exit 0 into exit 1 with
per-unit `malformed`, both gut variants. No extension of manifest proof to install/status/doctor
is needed for add-ons; that part of your proposal is larger than the problem.

Your per-target "exact expected layout" machinery for the core is the piece I say is
**unnecessary as proposed**, and I reproduced before saying so:

- The native core tree carries no `parley-addon.json` (re-verified), and its managed path never
  consults one (`validateInstalledPayload` calls `manifestProblems` only for `kind === "addon"`;
  the core's kind is the target kind). The packaged core manifest is consulted only on the
  unmanaged proof path, where the foreign copy is verbatim and matches the source-dir manifest
  exactly (measured: foreign fleet all `valid-unmanaged`, exit 0).
- The core full-gut **is already closed today**: a native core gutted to `SKILL.md` + marker is
  `malformed`, exit 1, because the core's required list is four files, not one (measured:
  missing `references/COOPERATION.md`, `references/compatibility.json`,
  `agents/manifest.yaml`). The stated B3 invariant — "gutted down to `SKILL.md`, with or
  without a marker" — holds for the core without any new machinery.
- What your layouts would uniquely close is the **managed-core partial gut**, which I measured
  and confirm is open: deleting six of the eleven installed files (`agents/openai.yaml`,
  `references/WORKED_EXAMPLES.md`, `plugin.json`, `gemini-extension.json`, `README.md`,
  `LICENSE`) from a marker-present native core still reports `valid`, `managed: true`, exit 0.

So the machinery is not solving a problem that does not arise — it is solving a problem that is
real but smaller than its cost. That cost is claude-1's dilemma made structural: the gemini
target rewrites `gemini-extension.json` at stage time, so one file's installed bytes are
destination-dependent; an exact per-target manifest either marks entries volatile (an integrity
hole declared in the file) or re-implements the installer in the verifier (a second source of
truth that can drift from the first). Both are the failure shape this area keeps producing.

Concrete counter-proposal for the partial gut, if the round wants it closed: **extend the
core's per-target required-files lists** (`validateInstalledPayload` :2230-2234) to the full
installed file set per target kind — add `agents/openai.yaml` and
`references/WORKED_EXAMPLES.md` everywhere, and the target-appropriate extras
(`plugin.json`, `gemini-extension.json`, `README.md`, `LICENSE`; `skills/SKILL.md` already
listed for antigravity). That closes deletion of any installed file with a one-line-per-kind
data change, no new concepts, no target-dependent bytes. It does not byte-verify contents — a
modified-but-present file still passes — and I argue that is the right stopping point: the
marker binds provenance on a managed tree, and content-modification detection on the assembled
core is precisely the assembled-shape manifest complexity I am rejecting.

Two pieces of your proposal I do adopt: (a) the generator's no-argument `--check` must treat
manifest presence as mandatory for every `skills/*/SKILL.md` directory, or a seventh skill
repeats this defect silently — measured that after the fix the default refresh covers all six,
but the opt-out shape survives for any future skill; (b) your ownership regressions — the
`valid-unmanaged` fleet must keep blocking install/uninstall without `--force` for all six
units, core included.

Version skew: we agree — anchored proof means a faithful foreign install of repo vY read by
installer vX (X ≠ Y) reports `malformed`. Accepted in round 1, still accepted.

### @hermes-1

Your round-1 analysis holds up almost entirely; two of your concerns resolved by measurement,
one in your favor and one against:

- **Concern 1 (marker-deleted native core → "unexpected: plugin.json"): does not happen.**
  Measured on the narrow-fix copy: delete only the core marker from a native install →
  `malformed` with the *unchanged* message "no parley-deck-skill install marker…". The native
  core tree carries no manifest, so `manifestFileHash(dest)` returns null and
  `unmanagedButVerified` returns false before `verifyPayload(dest)` is ever consulted — the
  `unexpected:` arm is unreachable for the native core. Your residual worry and the
  documentation question around it can both be dropped.
- **Concern 5 (leave the managed core on required-files only): I agree, and I measured exactly
  what that leaves open** — the partial-gut hole described in @codex-1 above (6 of 11 files
  deleted, still `valid`). Your asymmetry note was correct; I now quantify it. If the round
  wants it closed, my counter-proposal is the required-files extension, not codex-1's layouts.
  Note the full gut is already closed by the existing four-file list, so the stated B3
  invariant is not what is at stake.

Your Q2 drift analysis: confirmed in passing — `--check` over all six passes on the narrow-fix
copy, and the stale-source-manifest failure direction you described (fail closed, same as
today) follows from `verifyPayload(source)` being the first gate in `unmanagedButVerified`.

### @kimi-1

Correcting my own round-1 residuals list. I wrote "native core stays marker-anchored… that is
the documented B3 posture" without having measured the marker-present gut cases; the
measurements now complete the picture and change one recommendation. Round-1 me would have
called the migration nudge below scope creep; the cross-version measurement (new concern 1) is
why I now propose it. Everything else from round 1 stands and re-verified: foreign fleet
fix, B3 foreign arms (manifest kept and deleted), ownership fail-closed, the
do-not-add-the-manifest-to-PAYLOAD_ENTRIES guard, committed-not-pack-time manifests.

## New concerns / questions

1. **The false green persists for the installed base until reinstall.** Measured: take the
   pristine-2.1.0 install (four add-ons gutted, `markerSchema: 2, manifest: false` markers) and
   run doctor from the *narrow-fix* package against it → still all `valid`, `managed: true`,
   exit 0. `manifestProblems`'s `declared === false` branch returns healthy because the tree
   genuinely carries no manifest. Shipping manifests protects only installs made after the fix.
   Proposal: extend the existing, ratified precedent — the schema-undefined branch at
   installer.js:2280 already returns "re-install to validate it" when `sourceHasManifest` — to
   the `declared === false && sourceHasManifest` shape. After this fix, no current install can
   legitimately record `manifest: false`, so that marker state can only mean "installed by a
   pre-manifest package", exactly parallel to the 2.0.0 case. Cost, stated honestly: every
   faithful 2.1.0 install of the four add-ons reports `malformed: re-install to validate` after
   the package upgrades, until reinstalled — a deliberate red wave with a one-command repair.
   This is the only code I propose beyond the narrow fix, and it is a decision the round should
   take explicitly, not absorb.
2. **The seventh-skill hole in `--check`.** Endorsing codex-1's point with a measurement
   attached: after the fix, no-names `--check` covers all six only because all six carry
   manifests; a future skill added without one is silently skipped by both refresh and check.
   The generator change must make manifest presence mandatory for every `skills/*/SKILL.md`
   directory in check mode, not merely drop the core exclusion.
3. **For the record in round 3:** the brief's premise that "the managed path never consults
   the manifest" is false for add-ons and true only for the core. If round 3 argues from the
   brief's framing it will over-build. The accurate split: add-ons — manifest consulted on
   both paths once shipped; core — consulted on neither managed path (by kind) nor present in
   the native tree (by copy plan), only on the unmanaged proof path.
4. Question for claude-1: your round-1 file counts 6 source files for the core; the installed
   core carries 11. My required-files counter-proposal needs the per-target installed sets
   enumerated exactly (especially whether antigravity's `skills/SKILL.md` and the optional
   README/LICENSE pair should both be required on every target). Do you see a target where
   requiring the full installed set breaks a legitimate shape? I measured codex only.

## Current proposal

The fix must close **both** defects — and the measured fact is that the narrow fix already
does so for add-ons on both paths. Final shape:

1. **Ship `parley-addon.json` for all six units**, committed, generated by
   `scripts/build-addon-manifest.js` with the core exclusion dropped; no-names `--check`
   mandatory over every `skills/*/SKILL.md` directory; `--check` stays wired in `npm test` and
   `prepack`. No `runtime` blocks for the five. This alone closes: the foreign-install false
   red for the four add-ons, the marker-present gutted-tree false green for the four add-ons
   (measured, zero extra code), and the marker-deleted add-on case becomes provably
   `valid-unmanaged` (bidding's ratified behavior, generalized).
2. **Two-site core anchor** (`sourceRoot` on the core unit; `unmanagedButVerified` reads
   `unit.addon ? unit.addon.root : unit.sourceRoot || null`). Closes the foreign-install false
   red for the core. Measured: foreign fleet all six `valid-unmanaged`, exit 0; all foreign
   B3 gut variants still `malformed`.
3. **Migration nudge (one branch, explicit decision):** `declared === false && sourceHasManifest`
   in `manifestProblems` → "re-install to validate", mirroring the ratified :2280 precedent.
   Closes the false green for the pre-fix installed base. Accept the upgrade red wave
   knowingly, with the CHANGELOG entry saying why and how to repair.
4. **Do not adopt per-target exact-layout manifests for the core.** The native core carries no
   manifest, its managed path consults none, and its full gut is already closed by the
   required-files list (all measured). If the round wants the measured managed-core partial
   gut closed, extend the core's per-target required-files lists to the full installed file
   set — data, not machinery. Do not add the manifest to PAYLOAD_ENTRIES; do not mark entries
   volatile; do not re-implement the copy plan in a verifier.
5. **Regressions** (additions to the round-1 set, all measured this round):
   - codex-1's marker-present gutted-tree test, table-driven over all six: native install,
     gut to `SKILL.md` + marker, assert `malformed` + exit 1. Fails on four units at 2.1.0 —
     the true-positive regression for the false green. Include the manifest-kept gut arm.
   - Migration: a `markerSchema: 2, manifest: false` marker on an add-on whose source ships a
     manifest is flagged re-install-to-validate (item 3).
   - Core: native core gutted to `SKILL.md` + marker stays `malformed` (passes before and
     after); if item 4's required-files extension is taken, native core minus
     `references/WORKED_EXAMPLES.md` + marker becomes `malformed` (fails at 2.1.0).
   - Guard: native-installed core does not carry `parley-addon.json`; marker-deleted native
     core stays `malformed` with the unchanged no-marker message (measured).
   - Ownership: `valid-unmanaged` blocks install/uninstall without `--force` for all six.
6. **Docs and release:** README's "one skill that ships a parley-addon.json" passage rewritten;
   CHANGELOG known limit resolved, plus the two behavior notes (managed add-ons now
   byte-verified on edit; pre-manifest installs flagged for reinstall); version bump; five
   channels.
