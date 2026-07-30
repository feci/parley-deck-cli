---
agent: codex-1
idea: integrate-parley-bidding-addon
round: 3
date: 2026-07-30
responding-to: [user-ruling/manifest-anchor]
---

## Verdict on marker-anchoring

Marker-anchoring closes B3.11 and B3.13 for installer-managed trees. The rule is keyed to an
observed property of each installed unit, not its name: an add-on that never had a manifest
records `manifest: false` and retains `SKILL.md`-only compatibility; this release's validated
`parley-bidding` source records `manifest: true`, so removing its installed manifest cannot
silently downgrade it.

There is one boundary: the marker cannot detect a manifest omitted before the first install
observes the source. B3 is therefore satisfied only if the existing repository/package
inventory gate establishes that `parley-bidding` ships `parley-addon.json` in the source,
npm tarball, portable payload, and native install. That release-time assertion is not an
installer registry. Once the trusted install event has occurred, the marker closes the
installed-tree hole rather than moving it.

## What defeats it

Anyone able to rewrite both the payload and its marker can set `manifest: false`, replace the
manifest and payload together, or roll back the whole unit to another internally consistent
state. A missing marker also defeats the design if `doctor` treats it as legacy or falls back
to `SKILL.md`-only validation.

That limitation is acceptable for the stated purpose: detection of accidental deletion,
truncation, drift, and botched payload-only edits. It is not security theatre; those are real
and common defects, and the anchor detects them. It is also not tamper resistance. Defending
against a writer deliberately changing both the evidence and its anchor would require a
signed or otherwise trusted digest outside the installed directory.

## Required additions (if any)

1. Version the marker schema. A marker written by the new installer with an absent, malformed,
   or type-invalid manifest field must make the unit `malformed`; it must not be interpreted as
   legacy. Only genuinely older marker schemas may retain the compatibility fallback. An
   expected installed unit with a missing or unreadable marker must also be unhealthy.
2. Store more than the Boolean. For `manifest: true`, record the manifest's declared aggregate
   digest and the raw SHA-256 of `parley-addon.json`. The aggregate must commit to the sorted
   path-plus-file-hash inventory. `doctor` and `status` must recompute the installed inventory,
   compare it with the manifest and marker aggregate, and compare the manifest bytes with the
   marker hash. The aggregate catches a self-consistent manifest/payload swap; the raw hash also
   binds manifest metadata.
3. Validate the source manifest before any destination write, copy into staging, revalidate
   the staged bytes against the computed anchor, write the marker there, validate the complete
   staged unit, and only then perform the atomic replacement. A present manifest that is
   malformed must fail installation; `manifest: false` plus an installed manifest is also an
   inconsistency, not a downgrade path.
4. Add negative tests for manifest deletion, marker deletion/corruption, removal of the new
   field from a current-schema marker, and an internally consistent manifest-plus-payload
   replacement. Keep positive coverage showing that legacy and newly installed manifest-free
   add-ons remain healthy.

## Position

Marker-anchoring is the correct generic design for B3. It preserves the no-registry property
and gives `parley-bidding` durable installed-tree integrity, provided the release inventory
establishes the first manifest and current markers fail closed. Its trust boundary must be
documented as defect detection after a validated install, not protection from a malicious
local writer.

Amendment: ACCEPT WITH CONDITIONS
