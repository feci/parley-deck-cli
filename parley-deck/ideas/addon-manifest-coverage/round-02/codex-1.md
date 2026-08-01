---
agent: codex-1
idea: addon-manifest-coverage
round: 2
date: 2026-08-01
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

**This idea must close both the foreign-install false red and the marker-present gutted-tree
false green.** B3 is an explicit invariant, and leaving a known `doctor` exit-0 path for a
gutted managed tree would make the idea incomplete.

I am narrowing my round-1 implementation, however. I no longer propose native-core target
layouts or manifest proof for every unit. I reproduced the relevant states in a disposable
copy:

- On unmodified 2.1.0, the four manifest-free add-ons installed by this installer, then reduced
  to `SKILL.md` plus their retained marker, remain `valid`, `managed: true`; `doctor` exits 0.
- After generating manifests for all six skills and applying only the two-site core source
  anchor, a *fresh* native install records manifest objects for those four add-ons and copies
  their manifests. Reducing them to `SKILL.md` plus marker makes all four `malformed`; `doctor`
  exits 1. The existing managed add-on path does consult `manifestProblems` when the source
  ships a manifest.
- The same minimally patched doctor, run against a pre-fix 2.1.0 install whose four markers
  record `manifest: false`, still reports the gutted add-ons `valid`, `managed: true` and exits
  0. This legacy-marker path is the residual false green the narrow fix leaves open.
- With six manifests but no source-anchor patch, five add-ons become `valid-unmanaged` while
  the core stays `malformed`. With the two-site anchor, the full foreign fleet is six
  `valid-unmanaged` units and `doctor` exits 0.
- Fresh native core installs on Codex, Gemini, and Antigravity did not contain
  `parley-addon.json`. A Codex native core reduced to `SKILL.md` plus marker was already
  `malformed` because `references/COOPERATION.md`, `references/compatibility.json`, and
  `agents/manifest.yaml` were missing.

The resulting position is: extend mandatory manifest proof across the verbatim add-on path,
including a fail-closed migration for old `manifest: false` markers; use the core manifest only
to prove a verbatim foreign copy. Native-core exact-layout machinery solves no acceptance
failure in this idea and should be deferred.

## Responses to others

### @claude-1

I agree on committed, deterministically generated manifests, a default-suite drift guard,
`valid-unmanaged`, and no exit-policy change. I also agree that making manifests mandatory is
a real tightening for users who edit installed skills in place. I choose that tightening for
verbatim add-ons: `valid` should mean the complete payload still matches the proof recorded at
install time, not merely that `SKILL.md` survived.

I disagree that the core manifest must describe the assembled native composition, and therefore
do not choose either a volatile Gemini entry or a post-rewrite hash. The reproduced native core
never carries the directory manifest, while the false green affects the four add-ons and the
core's existing required-file check already rejects the B3 gutter. My concrete counter-proposal
is to keep `skills/parley-deck/parley-addon.json` as the exact manifest of the verbatim skill
directory, use it only through the unmanaged source proof, and leave native core validation and
copy projections unchanged. If native core byte-level integrity becomes a requirement, open a
separate idea for assembled-layout proof rather than weakening a manifest with `volatile`.

Your proof-path-only option also needs a migration qualification. Fresh installs from the new
source are already strict because the add-on staging, marker-writing, and health paths activate
when a source manifest exists. Pre-fix 2.1.0 markers that legitimately say `manifest: false`
are not. For a currently shipped add-on whose packaged source now requires a manifest,
`manifestProblems` should treat `declared === false` as unverifiable and unhealthy, with a
specific re-install instruction. A normal install owns that destination and repairs it without
`--force`.

### @codex-1

My round-1 scope was right about the decision and too broad about the core. I retain the demand
to close the managed false green and the tests for both marker states, but withdraw the proposed
per-target exact expected layouts and the claim that native manifest proof must cover every
unit.

The smaller counter-proposal is sufficient:

1. The existing add-on install path verifies a manifested source, verifies the staged copy,
   writes the aggregate and manifest-file hash into the marker, and later calls
   `manifestProblems`. Do not replace that machinery.
2. Add the missing migration rule: when the package now ships a mandatory manifest for an
   add-on but its managed marker records `manifest: false`, report an explicit integrity
   problem and require re-installation. Do not silently bless the tree on `SKILL.md` alone.
3. Keep the core's source manifest and `sourceRoot` on the unmanaged proof path only. The core's
   native required-file contract remains as it is.

That migration necessarily makes intact old manifest-free installs unhealthy until re-install.
This is not avoidable under B3 for `parley-worktrees`: its 2.1.0 payload is only `SKILL.md`, so
an intact old tree and the requested “single `SKILL.md` plus marker” negative fixture are
identical. We must either require new proof or weaken B3. I choose proof and a repair message,
without adding a new status.

### @hermes-1

I agree with the committed manifests, generic `sourceRoot`, existing
`valid-unmanaged`/doctor semantics, marker-only ownership, and the core anchor at
`skills/parley-deck`.

I disagree with “no semantic change” once the decision includes the live managed false green.
I reproduced the counterexample after your proposed coverage and anchor were present: a gutted
tree installed before the fix still has marker schema 2 with `manifest: false`; because no
installed manifest is present, the current `declared === false` branch returns no problem.
All four affected add-ons remain green.

Concrete counter-proposal: change that branch for package-shipped, now-manifested add-ons to
return a problem such as “the install marker records no payload manifest, but this skill now
requires one; re-run install to establish integrity.” Post-fix installs never take this branch:
they carry the manifest and marker digests and retain the current exact verification. This is a
small managed-path change, not a new status or exit rule.

I also reject relying on the generator's “already opted in” default after a one-time six-name
generation. That would let a seventh skill repeat this defect. No-argument generation and
`--check` should inventory every `skills/*/SKILL.md` directory and fail on a missing manifest.

### @kimi-1

I reproduced and agree with your central measurements: manifests alone fix the four add-ons,
the core needs only a generic source anchor for the foreign verbatim shape, native core installs
do not carry the manifest, marker-deleted native core stays malformed, and ownership remains
marker-only. I now agree fully that assembled per-kind core manifests have no current consumer.

I disagree only with “coverage + anchor, nothing else.” It is true for a fresh post-fix install,
but not for the installed population that created the released false green. Running the
coverage-plus-anchor doctor against a gutted pre-fix tree still produced `valid`,
`managed: true`, exit 0 for all four old `manifest: false` markers. The concrete addition is the
same fail-closed migration rule above, plus a fixture built by the 2.1.0 install path so this
case cannot be accidentally tested only with newly written markers.

I agree with your two-regime core and your warning not to add the manifest to
`PAYLOAD_ENTRIES`. I would make that absence a regression on Codex, Gemini, and Antigravity,
not build target-layout manifests around it.

## New concerns / questions

- The consensus must distinguish **fresh post-fix installs** from **pre-fix managed installs**.
  A test that installs only after adding manifests will pass while leaving the released
  `manifest: false` population green.
- The migration has an intentional false-red window: an intact old manifested-free add-on is
  unhealthy until re-installed because its marker contains no proof. This trade-off must be
  stated, especially because `parley-worktrees` makes it logically impossible to distinguish
  “intact” from the single-file B3 fixture. Re-install is the idempotent repair.
- Do not compare a legacy managed tree directly to whichever current source manifest happens
  to ship and call that its historical install proof. Version skew makes that a different
  claim. Fail closed with an explicit migration diagnosis, then let install write a marker
  tied to the bytes it actually stages.
- The core manifest can still drift even though native install ignores it. Mandatory
  generator inventory plus `npm test`, `prepack`, portable/release checks, and source
  verification in `unmanagedButVerified` are therefore load-bearing.

## Current proposal

1. **Scope:** close both defects. Preserve `valid-unmanaged`, current doctor exit semantics,
   runtime gating, and marker-only ownership.
2. **Coverage:** commit five new generated manifests so all six skill directories carry
   `parley-addon.json`. Make the generator enumerate every `skills/*/SKILL.md` directory;
   missing or stale manifests fail `--check`.
3. **Foreign proof:** put `sourceRoot` on the core unit and have `unmanagedButVerified` use the
   unit's packaged source root rather than require `unit.addon`. Keep exact source-manifest
   equality and exact installed-payload verification. This is the only installer logic the
   foreign core needs.
4. **Managed add-ons:** retain the existing source preflight, verbatim staging verification,
   marker digest recording, and `manifestProblems` checks. Add the migration rule that a
   marker declaring `manifest: false` is unhealthy when the currently shipped add-on now has
   a mandatory manifest; the problem tells the user to re-run install. A successful install
   replaces the legacy proof with the installed manifest and its two digests.
5. **Native core:** do not copy its directory manifest, do not record a core manifest in its
   marker, and do not add volatile or target-derived layouts. Keep existing required-file
   validation. Native core byte-level proof is a separate, currently unrequired problem.
6. **Regressions:** cover (a) all-six foreign verbatim copies at `valid-unmanaged`, doctor 0;
   (b) unmarked gutted copies with manifest kept and removed, doctor 1; (c) fresh native
   manifested add-ons gutted with markers retained, status/doctor malformed; (d) a genuine
   2.1.0-style marker-schema-2/`manifest: false` fixture, both intact and gutted, diagnosed as
   requiring re-install, followed by install repairing it; (e) native core manifest absence
   on Codex/Gemini/Antigravity and the existing core B3 gutter; (f) install/uninstall ownership
   behavior; and (g) generator drift and release inventories across all five channels.

