---
agent: kimi-1
idea: integrate-parley-bidding-addon
round: 3
date: 2026-07-30
responding-to: [user-ruling/manifest-anchor]
---

## Verdict on marker-anchoring

It closes both, and it is not just the hole relocated.

- **B3.11 holds.** The requirement keys off a per-unit fact observed at install time, not off
  any add-on name. `discoverAddons` (lib/installer.js:778-793) stays purely structural; no code
  path enumerates "add-ons that must have manifests." The other four install with
  `manifest: false` observed in their source payload, and `validateInstalledPayload` keeps
  requiring `["SKILL.md"]` for them (lib/installer.js:1131-1132). Untouched, as ratified.
- **B3.13 holds for new installs.** `parley-bidding` ships `parley-addon.json`, so its marker
  records `manifest: true`; `doctor`/`status` then treat the manifest as required for that
  unit, and deleting it flips the unit to `malformed` (the existing status vocabulary,
  lib/installer.js:1123). The acceptance test is implementable as stated.
- **The hole is genuinely closed, not moved.** Presence-keyed failed because the installed tree
  had no memory of what "correct" looked like. The marker *is* that memory. Crucially, this is
  not a new mechanism smuggled in for the occasion: the installer already trusts the marker to
  distinguish "absent by choice" from "broken" for the add-on set itself (`marker.addons`,
  `markerAddonNames`, lib/installer.js:823-838). Marker-anchored manifest validation is the
  same idiom applied one level down, which is exactly why it satisfies the "generic" clause
  instead of reintroducing the registry the 2.0.0 layout removed.

One seam must be named, though it does not break B3: units installed by a pre-amendment
installer carry markers without the field. The legacy rule must be "missing field = manifest
not required," which is precisely the precedent at lib/installer.js:826-827 (legacy marker
with no `addons` field is reported healthy, not broken). Consequence: a `parley-bidding`
installed by an old installer is not manifest-checked until its next install/upgrade rewrites
the marker. That is a stated, bounded cost, not a contradiction.

## What defeats it

Co-editing. Anyone — or anything — that removes the manifest can also flip `manifest: true` to
`false`, delete the marker file, or rewrite the digest. A botched merge or a sync tool that
restores an older marker over a newer tree defeats it accidentally; a deliberate actor defeats
it trivially.

That is acceptable for the stated purpose, and calling it theatre would be a category error:

- `doctor`/`status` are defect detectors over the operator's own tree, not an integrity
  boundary. An actor with write access to the installed tree can already replace `SKILL.md`
  or the payload itself; local tamper-resistance was never this layer's job, and nothing in B3
  asks for it. Install-time source verification is where integrity belongs.
- The codebase already treats the marker as honest bookkeeping rather than a security root —
  `marker.addons` is equally editable and is trusted for the same absent-vs-broken judgment.
  The amendment inherits that threat model; it does not invent a weaker one.
- For the failure class B3 actually exists to catch — a gutted tree from a partial copy, a
  moved file, a bad merge, an over-eager cleanup — the check is deterministic. It converts a
  *silent* downgrade into a `malformed` report, and it raises the dishonest case from "do
  nothing" to "make a second, consistent edit in a file you have to know exists."

It becomes theatre only if marketed as tamper-evidence. One hard requirement follows from the
co-editing risk, and it is already satisfied by the design: the marker lives inside the unit's
payload dir and must be written by the same atomic stage-then-rename path as the payload
(lib/installer.js:1041), so a crash cannot leave payload-with-manifest next to
marker-without-flag.

## Required additions (if any)

Two are required; one is a constraint, not a feature.

1. **Store the aggregate digest in the marker and verify it.** Without it, a manifest *swap* —
   syntactically valid JSON with different content: a changed add-on id, reduced declared
   requirements — passes both presence and parse checks, and that is the same gutted-tree
   defect class B3 exists to catch, arriving through the front door. The cost is one sha256 at
   install and one at check; `doctor` already reads the file. Pin the algorithm in the stored
   value (`sha256:<hex>`) so it can migrate. Scope the digest to the manifest only; whole-
   payload integrity is out of scope for B3. Missing manifest and digest-mismatch may both
   report `malformed`, but the reason string should say which.
2. **Version the semantics, not just the data.** Old markers lack the field: treat as "not
   required," backfill on the next install/upgrade, and never carry or hand-edit the field
   across installs — `writeMarker` always re-derives it from the observed source payload
   (lib/installer.js:1056-1074 already rewrites the marker wholesale on every install, so this
   is the natural behavior, but it must be stated).
3. **Document the threat model in the design text.** One sentence: this is defect detection
   for accidental tree damage, not tamper-evidence. That sentence is what keeps point 2 of
   this review from being a future objection.

## Position

Marker-anchoring is the correct resolution of the B3.11/B3.13 contradiction: it preserves the
generic discovery the 2.0.0 layout was built for, gives the installed tree the memory
presence-keyed lacked, and reuses an idiom the installer already depends on. With the digest
stored and verified, the legacy-marker rule stated, and the threat model written down, it
satisfies B3 in full.

Amendment: ACCEPT WITH CONDITIONS
