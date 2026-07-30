---
agent: hermes-1
idea: integrate-parley-bidding-addon
round: 3
date: 2026-07-30
responding-to: [user-ruling/manifest-anchor]
---

## Verdict on marker-anchoring

It closes both B3.11 and B3.13 genuinely. Not by moving the hole — by moving the
trust boundary from "what is on disk now" to "what the package shipped at install
time," which is a different and stronger anchor than either name-keyed or
presence-keyed provides.

I verified the mechanism against the actual code. `writeMarker`
(installer.js:1056-1075) already writes a per-skill marker
(`.parley-deck-skill-install.json`) inside each installed directory, recording
`name`, `skill`, `addon`, `version`, `source`, `target`, `scope`, `installedAt`,
and for core units the `addons` selection array. The marker is written to the
temp staging dir at :1035 *before* `validateInstalledPayload` runs at :1036, then
atomically renamed into place with the payload. Adding `manifest: true|false`
plus the aggregate digest to that same marker is a natural extension of a field
that already exists and already carries install-time observations.

**B3.11 — generic and optional.** The four legacy add-ons ship no
`parley-addon.json`. At install time their markers record `manifest: false`.
`doctor`/`status` see `manifest: false` → SKILL.md-only validation → `valid`,
unchanged. No add-on name appears in installer code. The `discoverAddons`
no-registry property (installer.js:778-793) is preserved. Closed.

**B3.13 — required for parley-bidding.** `parley-bidding` ships
`parley-addon.json`. Its marker records `manifest: true` and the aggregate
digest observed in the source payload. `doctor`/`status` see `manifest: true` →
require the manifest, validate per-file SHA-256, and compare the manifest's
aggregate digest to the marker's stored digest. Delete `parley-addon.json` after
install → manifest missing but marker says `manifest: true` → `malformed`. The
B3.13 acceptance test is implementable and passes. Closed.

The reason this is not just moving the hole: presence-keyed asks "is a manifest
on disk now?" — so deleting it changes the answer to "no" and silently downgrades
to SKILL.md-only. Marker-anchored asks "did the package ship a manifest when we
installed this?" — so deleting it does not change the answer, because the answer
is pinned to install-time observation, not current disk state. The question
changed, and the new question is the one B3 actually needs.

## What defeats it

Three things:

1. **Delete the manifest and the marker together.** `readMarker` (installer.js:1150-1156) returns `null` on a missing or unparseable marker. The policy for a missing marker must be specified (see Required additions). If it falls back to `manifest: false` → SKILL.md-only → `valid`, the double deletion defeats the check.

2. **Edit the marker to set `manifest: false`** (or remove the `manifest` field). The marker is a plaintext JSON file in the installed tree, beside the payload. Anyone with write access to the skill directory can edit it.

3. **Swap the manifest for one matching a gutted tree.** Without the aggregate digest stored in the marker, a fresh `parley-addon.json` listing fewer files with matching hashes would validate cleanly against the gutted tree.

Is this acceptable for a defect-detection mechanism? **Yes, and the comparison
makes it clear why.**

B3's stated purpose is catching *defects*: accidental deletion, disk corruption,
partial copies, bad upgrades, botched manual edits. It is not a security boundary
against an adversary with write access to the installed tree. An adversary who
controls the tree defeats all three approaches — they can rewrite the manifest to
match their gutted tree under name-keyed too, because the manifest is also in the
installed tree. No anchor that lives entirely on the same filesystem as the
payload can resist a determined writer. That is not what B3 was ratifiesd to do.

For the defect case, marker-anchored is *strictly stronger* than presence-keyed:

| Scenario | Presence-keyed | Marker-anchored |
|---|---|---|
| Delete manifest only | `valid` (defeated) | `malformed` (caught) |
| Delete manifest + marker | `valid` (defeated) | `valid` if no-marker=fallback (defeated) |
| Swap manifest (no digest in marker) | `valid` (defeated) | `valid` (defeated) |
| Swap manifest (digest in marker) | n/a | `malformed` (caught) |
| Bad upgrade drops manifest | `valid` (defeated) | `malformed` (caught) |
| Disk corruption deletes manifest | `valid` (defeated) | `malformed` (caught) |

Marker-anchored requires a *double fault* (manifest + marker) where
presence-keyed requires only a *single fault* (manifest). That is the
defect-detection margin, and it is real.

The marker is already the trust/ownership token throughout the installer.
`installSkillUnit` (installer.js:914-921) blocks on destinations that exist
without a marker. `uninstallSkillUnit` (installer.js:974-983) refuses to act on
unmarked destinations. Extending the marker to carry manifest presence is
consistent with the existing trust model — the marker already answers "was this
installed by us, and how?" Adding "did it ship a manifest?" is the same kind of
install-time observation as the `addons` selection array already recorded there.

## Required additions (if any)

Two, both already present in the user's proposal but one needs to be promoted
from "also" to "must":

1. **The aggregate digest stored in the marker is a hard requirement, not
   optional.** Without it, defect #3 (manifest swap) is uncaught: a fresh
   `parley-addon.json` matching a gutted tree validates against that tree and
   `doctor` says `valid`. The marker says `manifest: true`, the manifest is
   present, the per-file hashes match the gutted files — everything checks out
   except the one thing that would catch it: the aggregate digest the marker
   recorded at install time no longer matches the manifest's aggregate. This is
   the difference between catching *deletion* only and catching *substitution*
   too. The user wrote "(and the aggregate digest)" with parentheses; I am
   reading that as a required part of the design, and I want that on the record
   explicitly.

   Implementation: `writeMarker` reads `parley-addon.json` from the source
   payload, extracts the aggregate digest, and stores it as
   `aggregateDigest` alongside `manifest: true`. At `doctor`/`status` time,
   after confirming the manifest is present, compare the manifest's aggregate
   to the marker's stored value. Mismatch → `malformed`. This catches any
   manifest swap that changes the file set (different aggregate) — and a swap
   that keeps the same aggregate must keep the same file hashes, meaning the
   same payload, so the per-file validation still holds.

2. **Explicit policy for missing or legacy markers.** A marker without the
   `manifest` field (written by a release that predates this change) or no marker
   at all (foreign install via `npx skills add`, or marker deleted) must be
   treated as `manifest: false` → SKILL.md-only validation. This mirrors the
   existing `addons` field migration: test at installer.test.js:626-641
   ("a legacy core-only marker without an addons field stays healthy") already
   establishes that a missing field means the pre-feature default. The same
   pattern applies: `manifest` absent → `false` → SKILL.md-only → existing
   installs of the four legacy add-ons stay `valid` with no reinstall. This is
   backward-compatible and should be tested the same way.

   The cost is explicit: a foreign install of `parley-bidding` (no marker) would
   pass `doctor` with SKILL.md only. This is acceptable because the installer
   already treats unmarked destinations as "not ours" (blocks install, refuses
   uninstall), and a user who installed via a foreign tool chose that path. The
   manifest requirement binds installer-owned installations, which is where
   `doctor`'s authority comes from.

## Position

Marker-anchoring is the right resolution. It is the only option that closes both
B3.11 and B3.13 without reintroducing a hardcoded add-on registry, and it is
strictly stronger than presence-keyed for the defect-detection purpose B3 was
ratified for. The two additions above are clarifications that pin down what the
user's proposal already implies — the aggregate digest is load-bearing (not
decorative), and the legacy-marker policy must be explicit and tested. With those
two nailed down, the design is implementable in the existing marker/validation
code path with no structural changes.

Amendment: ACCEPT WITH CONDITIONS
