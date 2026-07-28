---
idea: parley-design-skills
status: fix-up-cycle-4
implementer: claude-1
started: 2026-07-28
completed: 2026-07-28
branch: parley-deck-skill#parley-design-skills
head-commit: 17f6619
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
does. Eight rules are in that state, not five: `core:text-below-legible-floor`,
`core:unlabelled-inference`, `core:value-off-scale`, `core:colour-off-ramp` and
`web:viewport-hero` report "no detector implements this rule", and `web:contrast-ratio`,
`web:target-size` and `web:reflow-narrow` are equally undetected here, reported instead
under the tier that is above this checker (kimi-1, round-01). For the two threshold rules the
reason is structural — their numbers live in annex prose, and copying a calibration into a
tool creates the second representation this project rejected in C2. All eight are visible in
the generated capability output rather than hidden, and the two `system`-class ones are named
on any L3 result as `system-rules-not-decided`, so a verified level is never read without
them.

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

## Fix-up cycle 1 — the agreed fixes AF-1..AF-9
status: superseded by cycle 2
completed: 2026-07-28
head-commit: 17f6619

### Fixes applied
Conformance levels modelled as explicit obligation sets, so a level cannot be certified on
evidence never obtained; waiver counter-signature independence; exit 0 reserved for PASS;
the artifact frontmatter subset ratified and all eight published examples rewritten to it
(4 of 8 failed the shipped parser before); unknown rule ids surfaced as UNJUDGEABLE;
reduced-motion coverage judged by what the block declares rather than by selector presence;
`:focus-visible` replacements found across the stylesheet; G1's banned-slop signature defined
normatively as a derivation over the registry (`slop` class at `T0 ARTIFACT`) rather than
named in three documents and defined in none; G1's two dropped C7 conjuncts restored; the U1
rotation given a `run-id` to hash and a checker that recomputes it.

### Deviations from agreed fixes
The per-file byte thresholds were rebalanced a second time (PDS.md to 25 KiB) because AF-7,
AF-8, AF-9 and AF-4 add irreducible normative text. The binding 64 KiB total held. The
per-file numbers are now documented as early-warning thresholds that deliberately sum above
the total, not as a partition — which also closes a reviewer NIT that they never summed to it.

### Independent verification
7 FIXED, 2 PARTIAL. The verifier ran every original reproducer against both a pre-fix
`git archive` tree and the post-fix tree, so FIXED means the defect reproduced before and does
not reproduce after — not that a test passes.

## Fix-up cycle 2 — the four residual self-attestation doors

An adversarial re-probe of cycle 1 found seven findings fixed, two partial and two new. All
four were reproduced before being fixed, and each carries a fixture that fails against the
cycle-1 code and passes against this one (verified by reverse-patching the cycle-2 hunks into
a copy of the tree and running the suite against it: 9 fail there, 0 here).

- **R-1 — G2's first conjunct was never enforced.** Only "cannot be re-expressed" was checked,
  so a graft that *added* its token to the winner's file re-expressed perfectly and verified
  L3 with `unmet []`. PDS §2 VERDICT now requires `tokens-digest` — the first twelve hex
  characters of sha256 over the winner's token file as ratified, the same digest form §11
  rule 3 already uses — and §3 G2 says the file is re-read and compared against it. An absent
  digest leaves the conjunct `unverified`, never met. Fixtures: *a graft that adds its token to
  the winner's file fails G2, re-expressible or not*; *a verdict recording no tokens-digest
  leaves G2's first conjunct unverified*.
- **R-2 — G1 trusted the run's own `g1-signatures` ledger.** Ban-list findings the run's own
  detectors raised against the direction artifacts are now collected *before* waivers, the
  sharing test is recomputed from them, and a recorded signature that omits an id the same run
  watched fire fails the gate. Detectors name the declared value that evidenced the finding, so
  the observed sets are compared on RULES.md's own terms. Fixtures: *a G1 ledger the run's own
  findings refute cannot verify L2*; *a waiver cannot launder G1*.
- **R-3 — independence was reported where none was established.** With no participant-bearing
  artifact the roster check was skipped entirely and two distinct strings suppressed a finding.
  Independence now requires both ids in the roster the run's artifacts name; with no roster the
  waiver does not suppress, and the rejection is printed as a `waiver rejected:` line. Ordered
  last among the waiver checks so a waiver wrong for a nameable reason is still told that
  reason. Fixtures: *a valid waiver suppresses its finding when the run records both signers*
  (roster present); *with no roster to check the signers against…* and *two ids differing is
  not independence…* (roster absent, the `nobody-1` / `ghost-2` probe).
- **R-4 — two remaining self-attested doors.** (a) Recusal is decided from the critique's
  artifact path, not its `agent` field, and a declared agent that disagrees with its own file
  name fails on its own. (b) `disposition: waived` must resolve to a valid, unexpired,
  independently counter-signed entry in the waiver file. Fixtures: *recusal is decided from the
  artifact path…*; *a waived disposition that resolves to no valid waiver entry fails G2*; *…
  backed by a valid, independent waiver entry passes G2*; *… whose waiver is self-signed fails
  G2 with it*.

`parley-design-check/SKILL.md` gains a "What it will not take on trust, and what it still
cannot bind" section: a table of the five conditions now recomputed rather than believed, and
three bindings this version explicitly does not claim — file naming is an anchor, not a
signature; independence is roster membership, not disinterest; and `tokens-digest` detects
drift from what the VERDICT ratified, not a re-ratification.

Doctrine cost of R-1: 469 bytes. Total is **65,286 of the binding 65,536** (D-1 / consensus
C3), 250 bytes spare. `RULES.md` is untouched, so the declared `registry-digest` still holds.

## Fix-up cycle 3 — the two residuals cycle 2 left open
status: complete
completed: 2026-07-28
head-commit: 17f6619

### Fixes applied
- A `disposition: waived` answer now requires a waiver entry whose scope actually covers the
  winner's work. It previously resolved by rule id alone, so an entry scoped at a path that
  need not even exist answered a violation against the winner. The path test reuses
  `waiverMatches` rather than adding a second implementation.
- A critique filed under an id minted from a proposer's own id (that id plus a suffix across a
  non-alphanumeric boundary) now fails recusal where it critiques that proposer's direction.
  The previous anchor caught impersonation of an existing id but not the strictly easier route
  of minting a new one.

### Deviations from agreed fixes
The verifier's literal remedy — flag any critique author no other artifact records — was not
implemented, because it also flags a legitimate critic who proposed nothing, which `PDS/1.0`
permits and which no artifact names since the spec defines no roster field. Implementing it
literally would fail an honest run; making it pass would mean adding a roster to the doctrine,
which is scope growth. Instead the condition is surfaced per run as `recusal-not-anchored` and
is never a silent pass.

### Residual, stated rather than implied
A proposer filing under an unrelated fresh id reads exactly like a genuine critic who proposed
nothing. Nothing in the artifacts separates them without a facilitator-held roster. This is
recorded as a non-claim in the checker's `SKILL.md` and surfaced per run.

### Verification
`npm test`: 197 passing, 0 failing (158 at review, 182 after cycle 1, 192 after cycle 2).

## Fix-up cycle 4 — the round-02 re-review blocks

Round-02 produced 2 CRITICAL + 1 MAJOR from codex-1 and 1 MAJOR + 1 MINOR from hermes-1.
Every reproducer codex-1 published was re-run through `bin/check.js` before and after.

### Fixes applied

- **C-1, six forgeable-certificate holes (AF-1).**
  - G1 now counts differences over the brief's declared axes and checks every declared
    position against the brief's enumeration, so an invented axis is a violation rather than
    a second difference (`checkG1`).
  - The §4 rule 2 rotation runs over the deduplicated primary positions, and a brief that
    lists a position twice fails the assignment (`checkAssignment`).
  - `pds-check:l2-process-order` reads each artifact against §1's mapping — DIRECTION in
    `round-01`, CRITIQUE in `round-02`, `DESIGN-BRIEF.md`, `consensus.md`, `FINAL.md` — so
    swapped rounds no longer verify (`PROCESS_HOMES`).
  - A gate recorded `pass` is recomputed for G3 and G4 from the run's own outcomes: open
    `system` rules plus the token-graph conjuncts for G3, open `quality` rules for G4
    (`gateRefutations`). PDS §1 rule 2 now states this normatively.
  - The stylesheet scanner is quote- and escape-aware, so `content: "}"` no longer closes a
    block and hides the declarations after it (`lib/css.js`).
  - A file declaring the spec under a kind §2 does not define is a candidate artifact and
    fails the new `pds-check:l1-artifact-kind` obligation. A file declaring the spec and no
    kind stays uninspected: that is how `PDS.md` and `RULES.md` declare the spec they define,
    and PDS §2 rule 1 now says so, so the checker never manufactures findings against the
    doctrine itself.
- **C-2, the waived work's author counter-signing (AF-2).** A scope is resolved to its
  protocol owner through §1's naming (`round-0N/<agent>.*` belongs to `<agent>`), and a
  counter-signer matching that owner is rejected. A counter-signer the run records only
  through an artifact that signer wrote itself is also rejected, reusing the corroboration
  set recusal already reports. PDS §8 rule 2 defines both, and states that work §1 gives no
  owner has no author to exclude — otherwise a waiver on ordinary source could never be
  counter-signed at all.
- **C-3, the single-agent fast path.** Resolved as *outside* Parley. PDS §4 rule 8 and
  `SKILL.md` now say a fast-path run is not a Parley Deck workflow (`COOPERATION.md` §1 makes
  multi-agent execution mandatory, §0.3 rule 2 gives it the process), MUST NOT claim Parley
  verification and MUST NOT claim a level above L1. The dispatcher row is labelled
  accordingly, so nothing implies a solo Parley run.
- **H-1, `waivers: ""` crashed the checker (EISDIR, exit 2).** An empty path field names no
  file and the run continues. The same guard covers `tokens: ""` on the CONTRACT and on a
  DIRECTION, which hermes-1's open question asked about, and every resolved path is confirmed
  to be a file before it is read.
- **H-2, alias direction.** Implemented rather than dropped. PDS §3 G3 defines it by group
  name — `primitive`, `semantic`, `component`; a reference points down that order and a
  primitive holds a value, never another primitive — and states that a document naming none
  of those groups declares no direction, so the conjunct is vacuous rather than failed. The
  checker enforces it as `pds-check:l3-alias-direction`. FINAL.md and the shipped spec now
  agree.

### Deviations from the review's literal remedies
codex-1 asked that an unknown **or missing** kind be an L1 violation. A missing kind is not
treated as one: `PDS.md` and `RULES.md` both declare `spec: PDS/1.0` with no `kind`, and
flagging that would make the doctrine fail its own checker. The line is drawn at a kind that
is present and undefined, which is the reproducer, and PDS §2 rule 1 states the distinction.

### Byte budget
The doctrine grew by the new normative text and was compressed to pay for it: 65,374 → 65,364
of 65,536 bytes (SKILL 6,519 · PDS 25,598 · RULES 23,225 · WEB-ANNEX 10,022). No per-file
threshold moved. `RULES.md` changed, so `PDS.md`'s `registry-digest` was recomputed to
`b49ff596451f` in the same commit (§11 rule 3).

### Verification
`npm test`: 212 passing, 0 failing (197 before this cycle). Fifteen new tests, one per hole
plus the passing-side controls; each was proved to fail against the pre-fix code by reverting
that fix alone and re-running only its covering tests. All ten of codex-1's and hermes-1's
reproducers were re-run through `bin/check.js`: every one that returned PASS/exit 0 before now
returns a violation and an unverified level, and the `waivers: ""` crash now exits 0 clean.
