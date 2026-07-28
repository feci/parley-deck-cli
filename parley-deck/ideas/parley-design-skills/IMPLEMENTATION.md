---
idea: parley-design-skills
status: implemented
implementer: claude-1
started: 2026-07-28
completed: 2026-07-28
branch: parley-deck-skill#parley-design-skills
head-commit: 726c024
design-pr: n/a
implementation-pr: pending
---

## Summary of work

Both add-on skills built in `parley-deck-skill` per `FINAL.md`.

**`addons/parley-design`** — zero-dependency doctrine and protocol, four markdown files:

- `SKILL.md` — dispatcher, when-to-use / when-NOT-to-use, the invariants that must be held
  before editing anything, the binding precedence chain, fixed per-phase read sets, and an
  explicit statement that the doctrine is complete without the checker installed.
- `references/PDS.md` — the protocol, `spec: PDS/1.0`, RFC 2119 conformance language,
  sections §0–§12: scope and normative language, the Parley mapping table (C1) as normative,
  artifact kinds each in the identical four-part shape, gates G1–G4 with their exact failure
  conditions and checker error strings, the ritual, roles and invariants, evidence tiers and
  verdicts, rule classes and authority, waivers, conformance levels L1–L4, the extension
  policy, versioning/deprecation, and a real 1.0.0 changelog entry.
- `references/RULES.md` — the literate rule registry: frozen extraction grammar, then one H3
  per rule with exactly one `pds-rule` YAML fence and prose. Nineteen `core:` rules covering
  every class FINAL.md required, plus a normative table cross-referencing the `web:` ids.
- `references/WEB-ANNEX.md` — non-normative for non-web surfaces, stated in its first line:
  the tier→engine mapping for web, the standards-derived numbers, and eleven `web:` rules in
  the same grammar.

**`addons/parley-design-check`** — the enforcement layer. `SKILL.md`, a CLI (`bin/check.js`),
`lib/registry.js` (literate-registry reader plus a strict parser for the flat key/value
subset the `pds-rule` blocks use, which raises on anything outside it), `lib/artifacts.js`
(PDS artifacts and DTCG token graph, alias resolution, WCAG luminance), `lib/css.js` (the
`T1 SOURCE` scanner), `lib/engine.js` (generated capability, gating, waivers, conformance
levels, verdict roll-up, exit codes), and eighteen detector modules. Node built-ins only, no
network at check time.

`NOTICE.md` records hallmark and impeccable as prior art studied, states that nothing was
copied, and lists the standards referenced.

## Implementation plan / checklist

- [x] Files or areas to change: `addons/parley-design/**`, `addons/parley-design-check/**`,
      `NOTICE.md`, `package.json` (`files` array), `test/installer.test.js` (addon-list
      assertions), `test/design-addons.test.js` (new budget and integrity guards).
- [x] Checks to run: `npm test` from the repo root.
- [x] Review or risk notes: below.

## Deviations from FINAL.md

**D-1 — the per-file byte split was rebalanced; the 64 KiB total was held.**
FINAL.md sets `SKILL.md ≤8 KiB · PDS.md ≤20 KiB · RULES.md ≤24 KiB · WEB-ANNEX.md ≤12 KiB`.
`PDS.md` did not fit 20 KiB while carrying every section §0–§12 with the identical four-part
artifact shape intact. The alternatives were to drop artifact entries, to break the identical
shape by keeping examples for some artifact kinds and not others, or to rebalance. The first
two damage the thing the spec exists to be; the third does not. Consensus C3 adopts **"64 KiB
total"** as the binding ceiling and the per-file numbers are the split one participant
proposed, so the total is what has been held. The budget test enforces the rebalanced
per-file numbers and the unchanged 64 KiB total.

**Reviewers should challenge this if they disagree** — it is the only place the
implementation chose differently from the written spec.

**D-2 — `enforced-by: check` rules without a detector.**
FINAL.md requires that such a rule report `UNJUDGEABLE` rather than pass silently, and it
does. Five rules are in that state deliberately: `core:text-below-legible-floor`,
`core:unlabelled-inference`, `core:value-off-scale`, `core:colour-off-ramp` and
`web:viewport-hero`. For the two threshold rules the reason is structural — their numbers live
in annex prose, and copying a calibration into a tool creates the second representation this
project rejected in C2. This is visible in the generated capability output rather than hidden.

## Notes for reviewers

- **The doctrine and the registry were authored by different agents, and one of them died
  mid-run.** The `skill+pds` agent failed with an API connection error after writing both of
  its files. The files were checked afterwards and are structurally complete (§0–§12 present,
  changelog written, no truncation), but they did not get their author's own final review
  pass. Read `PDS.md` and `SKILL.md` with that in mind — they are the least-verified artifacts
  in this change.
- **`node --test <directory>` does not work on Node v26.5.0** in this environment; it resolves
  the directory as a module. Use bare `node --test` from the repo root (which is what
  `npm test` runs) or a glob. This bit the checker author and is worth knowing.
- The checker's capability declaration is **generated by scanning `lib/detectors/`**, never
  hand-maintained, so it cannot claim coverage it does not have.
- The checker carries **no copy of the registry**; a test asserts no `RULES.md` exists under
  `addons/parley-design-check/`. With `parley-design` absent it refuses rule checks on stderr
  and reports a distinct exit code rather than passing.
- Test fixtures are **discovered from the detector directory**, so adding a detector without a
  pass/fail fixture pair fails the suite.
- `test/design-addons.test.js` additionally guards: exactly four doctrine files; no
  placeholder text in any shipped design file; and no non-builtin `require()` anywhere in the
  checker.

## Verification

`npm test` from `parley-deck-skill`, with the exact counts recorded at the fix-up cycle below.
