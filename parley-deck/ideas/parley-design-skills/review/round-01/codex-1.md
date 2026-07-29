---
agent: codex-1
idea: parley-design-skills
review-round: 01
date: 2026-07-28
reviewed-commit: 726c024
---

## Summary

The implementation is broad, offline, well organized, and backed by a fast green test suite,
but it is not ready to merge. The checker can certify L2/L3 conformance without checking
binding obligations, accepts a waiver signed twice by the same participant, and exits zero
after inspecting nothing. I also reproduced an accessibility false negative and common CSS
and icon-source false positives.

The current registry itself parses cleanly, its digest matches `PDS.md`, capability is
derived from the detector directory, all detector ids resolve, registry absence is refused,
and declared deviation D-2 is reported honestly. I reject deviation D-1: `FINAL.md`, not the
implementation rationale, is binding, and it explicitly sets both the aggregate and
per-file limits.

## What I verified (commands run, and their result)

- `git rev-parse --verify HEAD`, `git branch --show-current`, and
  `git diff-tree --no-commit-id --name-status -r 726c024` in the implementation repo:
  HEAD was `726c0246c41293bdca35c471428657c9ca1aee24`, on
  `parley-design-skills`, with the expected 110-file change.
- `npm test`: 158 tests passed, 0 failed, in about 1.4 seconds.
- `node addons/parley-design-check/bin/check.js --help`: exited 0 and documented the
  advertised commands and codes.
- `git diff 726c024^ 726c024 --check`: clean.
- `wc -c` over the doctrine files: `SKILL.md` 6,656; `PDS.md` 22,389; `RULES.md`
  24,546; `WEB-ANNEX.md` 11,132; aggregate 64,723 bytes.
- `shasum -a 256 addons/parley-design/references/RULES.md`: digest starts
  `1fbe071e1222`, matching `PDS.md`.
- `rg` for placeholders and `find` for a checker-local `RULES.md`: no shipped placeholder
  was found and the checker contains no registry fallback.
- A registry/capability probe using `loadRegistry()` and `loadDetectors()` found 30 rules,
  18 detector rule ids, and exactly the five detector gaps declared in D-2:
  `core:unlabelled-inference`, `core:text-below-legible-floor`,
  `core:value-off-scale`, `core:colour-off-ramp`, and `web:viewport-hero`.
- An in-memory `runCheck()` probe over the shipped sound-run artifacts plus a valid token
  document, with no stylesheet, returned `verified: L3`, `verdict: PASS`, and exit 0 while
  `core:literal-outside-token-layer` was `UNJUDGEABLE`.
- An in-memory mutation changing the sound fixture's `assigned: flat` to
  `assigned: layered` still returned `verified: L2` with no conformance finding.
- An in-memory mutation changing a critique's rule id to `project:unknown-thing` still
  returned `verified: L2`; the unknown id appeared in neither findings nor
  `UNJUDGEABLE`.
- An in-memory mutation making `granted-by` and `counter-signed-by` both `claude-1`
  produced no waiver error and suppressed the literal-value finding.
- Direct detector probes reproduced: a reduced-motion media rule that only sets
  `color: red` suppresses the missing-reduced-path finding; a valid split
  `outline: none` / `:focus-visible { outline: ... }` pattern is reported as a violation;
  one Heroicons import is reported as two icon sources.
- Direct `parseYamlSubset()` probes rejected both an ordinary YAML block mapping and the
  multiline flow-map shape printed in `PDS.md`.
- `node .../check.js --json NOTICE.md` inspected no artifact, token, style, or markup,
  returned report verdict `UNJUDGEABLE`, and nevertheless exited 0.
- `npm pack --dry-run --json` could not run because the machine's npm cache contains
  root-owned files (`EPERM`). Package-tarball contents therefore were not independently
  verified; the installer and package-discovery tests did pass.

## Findings

### [CRITICAL] L2 and L3 can be falsely certified

`lib/engine.js` does not implement the conformance contract it claims to verify.
`checkL2()` validates only fragments of G1/G2 and requires recorded outcomes only for
`G1` and `G2` (lines 268–409). It does not verify U1's deterministic primary-axis
assignment, the assigned position against the direction, G1's banned-signature condition,
G2 token-file immutability or answered winner violations, or recorded G3/G4 outcomes.
`checkL3()` checks part of the token graph but does not make source availability or the
no-literals rule a conformance prerequisite when a registry is present. Finally, the
`verified` calculation at lines 731–739 considers only conformance-class violations and
`pds-check:` unjudgeables, ignoring relevant registry rules that are `UNJUDGEABLE`.

This is observable, not theoretical: a wrong assigned position still verified L2, and a run
with no styles verified L3, reported PASS, and exited 0 while the no-literals rule explicitly
said it could not be judged. That makes a machine conformance certificate weaker than the
binding PDS gates and can authorize Phase 5 or a clean audit on evidence the checker never
obtained.

Fix: model each level as an explicit obligation set. L2 must validate the U1 assignment and
all applicable G1–G4 records/conditions; any condition that cannot be recomputed must make
the level unverified. L3 must require a DTCG `2025.10` token document, declared
`colorSpace` on every colour, valid aliases, and actual source coverage for the no-literals
rule. Relevant registry `VIOLATION`, `NEEDS_REVIEW`, or `UNJUDGEABLE` results must feed the
level result. Add negative fixtures for a wrong assignment, missing G3/G4, modified winner
tokens, missing source, plain-string colour values, and unanswered winner findings.

### [CRITICAL] A participant can counter-sign their own waiver

`waiverProblem()` (`lib/engine.js` lines 187–211) checks only that
`counter-signed-by` is non-empty. It neither requires `granted-by`, compares the two
identities, nor establishes that the counter-signer is a different participant who did not
author the waived work. My probe changed the valid fixture to
`granted-by: claude-1, counter-signed-by: claude-1`; the waiver was accepted and the
underlying violation was removed.

This defeats C13's central protection. A unilateral author can suppress a quality or system
finding while the report falsely says it was counter-signed.

Fix: make the granting/author identity machine-readable and required, reject equal grantor
and counter-signer ids, validate both against the run's participant/ownership metadata, and
leave the finding unsuppressed when independence cannot be established. Add fixtures for a
self-signature, missing grantor, unknown signer, and signer who authored the scoped work.

### [MAJOR] An entirely unjudgeable run exits as clean

The roll-up correctly sets `verdict: UNJUDGEABLE` when `judgedNothing()` is true, but
lines 839–841 assign exit 0 to every registry-loaded result other than `VIOLATION` or
`NEEDS_REVIEW`. Running the CLI on `NOTICE.md` inspected no supported input, reported 19
unjudgeable rules, and exited 0. This contradicts both the CLI's “0 clean” contract and the
nearby source comment that an empty, uninspected run “is not a clean result.”

CI and pre-commit users gate on the process code, not a later manual reading of JSON, so this
turns “the checker checked nothing” into success.

Fix: reserve exit 0 exclusively for a `PASS` verdict. Give an overall `UNJUDGEABLE` result a
documented non-zero code (or explicitly broaden code 1 to all non-clean findings), update
the help/SKILL tables, and add a CLI test using an unsupported or non-artifact input.

### [MAJOR] Artifact ingestion rejects valid PDS YAML and drops unknown rule ids

The flat parser in `lib/registry.js` is adequate for the current `pds-rule` records, but it
is also used as the artifact-frontmatter parser. Lines 140–192 accept only top-level values
and one-line inline list items. They reject ordinary YAML block mappings and even multiline
flow maps printed by `PDS.md` itself, including its `DESIGN-SYSTEM`, `AUDIT`, and `WAIVERS`
examples. `collectInputs()` then demotes parse failures to `not-inspected`, allowing other
artifacts to carry a level claim.

Separately, artifacts that do parse are not traversed for rule references. Replacing a
CRITIQUE finding id with an unknown project id produced neither the required
`UNJUDGEABLE` pass-through nor an error, and L2 still verified. This violates §10 rule 3
and the frozen registry consumer contract.

Fix: either parse the YAML language the spec publishes, or ratify a canonical artifact
frontmatter subset and rewrite every example to that subset; do not silently omit a
candidate PDS artifact from conformance. Then centrally traverse rule-id-bearing fields in
CRITIQUE, VERDICT answers, AUDIT, and WAIVERS, emitting `UNJUDGEABLE` for ids absent from
the loaded registry. Add the exact multiline PDS examples and an unknown-id critique as
fixtures.

### [MAJOR] A reduced-motion media block passes without reducing motion

`motion-reduced-path.js` treats any universal selector inside a
`prefers-reduced-motion` query as covering every animation (lines 26–44), without inspecting
what the rule declares. This input returned no finding:

```css
@media (prefers-reduced-motion: reduce) { * { color: red; } }
.spinner { animation: spin 1s infinite; }
```

The detector therefore claims a `core:motion-without-reduced-path` capability while missing
the exact accessibility violation it exists to catch.

Fix: a reduced block must count as coverage only when its declarations actually remove or
materially neutralize animation for the matched base (for example `animation: none`, a
non-moving alternative, or an explicitly bounded equivalent). Selector presence alone is
not evidence. Add universal and named-selector negative fixtures whose reduced blocks change
unrelated properties.

### [MAJOR] A valid focus-visible replacement is reported as absent

`focus-indication.js` looks for an outline replacement only inside the same CSS block that
removes the outline (lines 52–66). It reports the common valid pattern
`button { outline: none }` plus
`button:focus-visible { outline: 2px solid currentColor }` as a severity-4 quality
violation.

Because quality findings may block on one reproducible report, this false positive can stop
a conforming build and will teach users to distrust the checker.

Fix: correlate outline removal and focus-indicator declarations by normalized selector base
across the stylesheet/cascade approximation. Report a source-tier violation only when no
focus/focus-visible rule supplies a replacement; use `NEEDS_REVIEW` where selector or
cascade ambiguity remains. Add split-rule and `:focus`/`:focus-visible` regression fixtures.

### [MAJOR] D-1 weakens a binding acceptance criterion

`FINAL.md` explicitly binds `PDS.md` to 20 KiB as part of the hard four-file budget.
The implementation's `PDS.md` is 22,389 bytes, 1,909 bytes over that ceiling.
`test/design-addons.test.js` lines 13–22 changes the test to 22 KiB, so the green test proves
only the implementation's replacement criterion, not the binding one.

Consensus C3's aggregate rationale does not erase the later, explicit per-file limits in the
binding `FINAL.md`, and `IMPLEMENTATION.md` cannot amend them unilaterally. I therefore do
not accept D-1.

Fix: reduce `PDS.md` to at most 20,480 bytes while retaining the required artifact shape and
sections, and restore the 8/20/24/12 KiB test. If that is genuinely impossible, ratify a new
specification before changing the acceptance test.

### [MINOR] A single Heroicons package is counted as two icon sources

`icon-provenance.js` runs overlapping generic and named-package regexes independently
(lines 13–21, 35–55). One import from `@heroicons/react/24/solid` is recorded as both
`imported icon package:@heroicons/react/24/solid` and `heroicons`, producing a mixed-source
violation from one source.

This is a slop rule rather than a unilateral blocker, but the false finding still pollutes
the concurrence process.

Fix: canonicalize each import to exactly one source id, with named package recognizers taking
precedence over the generic fallback, and add one-source fixtures for Heroicons,
Font Awesome, Material Icons, and other overlapping names.

### [MINOR] PDS normative requirements are not all numbered and named

`PDS.md` §0 rule 3 says every normative statement is numbered and bold-named, matching
`FINAL.md`. The artifact tables nevertheless contain unnumbered RFC 2119 requirements:
for example CONTRACT `winner`/`tokens` (lines 199–200), AUDIT `tiers`/`level`
(252–254), and WAIVERS id/expiry (276–280). Because §0 rule 2 also says all unlabelled
content is normative, these are not merely informative table notes.

This breaks the promised citation grammar and leaves tools/reviewers unable to cite a stable
rule number for several binding fields.

Fix: move each binding table requirement into a numbered, bold-named rule and make the table
reference that rule, or explicitly mark the table informative and keep all normative force
in numbered rules. Add a structural test that flags RFC 2119 keywords outside numbered
normative items.

## Open questions

1. Is the one-agent PDS fast path explicitly outside a Parley Deck workflow? PDS says
   `COOPERATION.md` wins on process, while the core protocol requires non-solo execution and
   `FINAL.md` defers any core-protocol carve-out. The docs should say whether fast-path work
   must not claim Parley, or pursue the deferred protocol change.
2. Should the audit identity digest cover the annexes whose rules the checker actually
   executes? The report exposes per-file annex digests in JSON, but its canonical
   `registry-digest` and text output cover only `RULES.md`, so a web-rule edit can leave the
   headline digest unchanged.
3. Actual npm tarball contents remain unverified because `npm pack --dry-run` was blocked by
   the machine's root-owned npm cache. Is there a clean packaging CI result for commit
   `726c024` that confirms `NOTICE.md` and both complete add-on trees are included?
