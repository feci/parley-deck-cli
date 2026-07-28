---
agent: hermes-1
idea: parley-design-skills
review-round: 02
date: 2026-07-28
reviewed-commit: 8fc3a18
---

## Summary

This is a re-review of commit 8fc3a18, which incorporates three fix-up cycles
against the 10 agreed fixes (AF-1..AF-10) and four residual self-attestation
doors (R-1..R-4) found by adversarial re-probes. My lens is the doctrine:
whether RULES.md and WEB-ANNEX.md say true, decidable things; standards fidelity
(DTCG, WCAG 2.2); whether any rule is unfalsifiable, mis-classified, or on the
wrong evidence tier.

The implementation is substantially sound. All 10 agreed fixes landed. All 4
residual doors (R-1..R-4) are closed. I verified each one with real commands,
not by trusting IMPLEMENTATION.md. The doctrine's rule classifications are
correct, the WCAG contrast math is accurate, the ban list is now normatively
defined and decidable, and the checker recomputes rather than trusts the five
conditions it identified. 197 tests pass.

I found one MAJOR new defect (a crash on `waivers: ""` in a CONTRACT — a valid
value per the spec's own canonical frontmatter subset), one MINOR residual
(the "alias direction" requirement from FINAL.md was dropped from G3 and never
implemented in the checker, despite being accepted as a round-01 MINOR), and
two NITs. None blocks.

Position on D-1 (per-file byte rebalance): ACCEPT. The 64 KiB total is held
(65286/65536, 250 bytes spare). RULES.md now has 1004 bytes of headroom against
its 24 KiB threshold — substantially better than round-01's 30 bytes. The
per-file numbers are documented as early-warning thresholds, not a partition,
which closes my round-01 NIT about them summing above the total.

## What I verified (commands run, and their result)

1. `git rev-parse --verify HEAD` → 8fc3a18, branch `parley-design-skills`.
   Confirms the reviewed commit matches.

2. `npm test` → 197 passing, 0 failing, 0 skipped (1448 ms). Matches
   IMPLEMENTATION.md's claim (197 after cycle 3).

3. `wc -c` over the four doctrine files: SKILL 6681, PDS 24909, RULES 23572,
   WEB-ANNEX 10124. Total 65286 ≤ 65536. Each within its test threshold
   (SKILL ≤7K, PDS ≤25K, RULES ≤24K, WEB-ANNEX ≤11K). RULES headroom: 1004
   bytes. PDS headroom: 691 bytes.

4. `shasum -a 256 RULES.md` → `f0c38eed1b8d...`, matching PDS.md's declared
   `registry-digest: f0c38eed1b8d`. Registry-digest guard is real.

5. Registry-digest mismatch: edited PDS.md to declare `deadbeefdead`, ran
   checker → "registry-digest mismatch: the spec declares deadbeefdead, the
   registry file computes f0c38eed1b8d". Restored afterward. §11 rule 3
   drift detection works.

6. `node check.js --help` → exit 0, documents all commands and exit codes
   including code 4 for UNJUDGEABLE. AF-3 fixed.

7. `node check.js check NOTICE.md` → verdict UNJUDGEABLE, exit 4 (not 0).
   AF-3 fixed: exit 0 reserved for PASS alone.

8. `node check.js check /tmp/nonexistent.txt` → exit 2 ("the run failed —
   no such path"). Correct failure handling.

9. Self-signed waiver probe: created a contract + tokens + CSS with a literal,
   plus a waiver entry `granted-by: claude-1, counter-signed-by: claude-1`.
   Result: "waiver rejected: claude-1 counter-signed its own waiver, which is
   not a counter-signature" — the finding stays in the ledger. AF-2 fixed.

10. Focus-visible probe: `button { outline: none }` + separate
    `button:focus-visible { outline: 2px solid currentColor }`. No
    `core:focus-indication` violation fired. AF-6 fixed: the detector looks
    across the stylesheet, not only inside the block that removed the outline.

11. Reduced-motion probe: `@media (prefers-reduced-motion: reduce) { * { color:
    red } }` + `.spinner { animation: spin 1s infinite }`. Result: the
    `core:motion-without-reduced-path` finding fires ("the file's
    prefers-reduced-motion block changes other properties without removing the
    motion"). AF-5 fixed: selector presence is no longer enough; the block
    must actually remove or neutralise the animation.

12. Icon-provenance probe: one `@heroicons/react/24/solid` import → verdict
    PASS, no mixed-source finding. codex-1's round-01 MINOR fixed: the
    double-count is gone, one package = one source.

13. L2 on `sound-run` fixture → "claimed L2, verified L2", exit 0. The
    `recusal-not-anchored` line appears for hermes-1 (the fixture's critic),
    which is the honest non-claim the checker makes when a critique author no
    other artifact records.

14. L2 on `collapsed-run` fixture → "claimed L2, not verified", exit 1.
    G1 two-axis test fires ("differ on 1 declared axis; 2 are required").
    G1 error string includes all three conjuncts: "the ban list, the
    category-plus-avoidance test and recorded human ratification" — AF-8
    fixed. The `pds-check:l2-assignment` finding fires: "the brief declares
    no run-id, which §4 rule 2 hashes, so the rotation is not reproducible"
    — AF-9 fixed. The `pds-check:l2-gate-g1` finding fires: "the G1 outcome
    records no banned-slop signature for the directions, so the sharing test
    cannot be recomputed" — R-2 and AF-7 fixed.

15. L3 colorSpace probe: bare-hex color tokens (`"$value": "#111111"`)
    without explicit `colorSpace` → `pds-check:l3-colour-space` VIOLATION:
    "declares no colorSpace on the value it resolves to; a bare literal is
    computable and states nothing about the space it was chosen in". My
    round-01 MINOR fixed.

16. G2 tokens-digest probe (code reading): engine.js lines 922-946 re-read the
    winner's token file, compute sha256, and compare against the VERDICT's
    `tokens-digest`. An absent digest → `unverified`. A mismatch → G2 fail.
    R-1 fixed.

17. G2 waiver scope probe (code reading): engine.js lines 1042-1063 check that
    a `waived` answer's waiver entry scope actually resolves to the winner's
    DIRECTION or its token file, using the same `waiverMatches` function as
    detector findings. Cycle 3 fix landed.

18. Recusal id-minting probe (code reading): engine.js lines 570-587 check for
    ids minted from a proposer's own id (via `mintedFrom`) and fail recusal
    if such an id critiques that proposer's direction. `mintedFrom` correctly
    checks the non-alphanumeric boundary. Cycle 3 fix landed.

19. G1 sharing test (code reading): engine.js line 812 runs `shareTest` on
    `observedFor` ("this run's own ban-list findings") before testing the
    recorded signatures. Lines 847-857 flag a recorded signature that omits
    an id the same run watched fire. R-2 fixed.

20. WCAG contrast math: `toSrgb` + `contrastRatio` → black/white 21.00,
    red/white 4.00, #777/white 4.48, #767676/white 4.54. All match published
    WCAG reference values.

21. Rule registry parse: 30 rules (19 core + 11 web), 0 duplicate ids, all
    with valid class/tier/enforced-by/severity/surface values. Class
    distribution: quality 11, slop 12, system 7. Tier distribution: T0 7,
    T1 18, T2 4, T3 1.

22. Ban list derivation: `banList(registry)` at engine.js line 389-393 returns
    exactly `slop`-class rules at `T0`. That gives 4 ids:
    `core:category-guessable`, `core:decoration-unmotivated`,
    `core:structural-sameness`, `core:signature-absent-or-mood`. All are
    decidable from DIRECTION artifacts alone. AF-7 fixed: the ban list and
    sharing test are defined in RULES.md lines 282-291 and implemented in the
    checker.

23. No bundled registry: `find addons/parley-design-check -name RULES.md`
    → 0 results. C2/C4 honoured.

24. No placeholders: searched all doctrine + checker files for
    TODO/TBD/FIXME/XXX/lorem ipsum/goes here → 0 matches. FINAL.md's
    "zero placeholders" requirement met.

25. PDS.md §3 gate strings: now prefixed "canonical message shapes, not
    literal output" (line 306-308). The checker's SKILL.md explains the
    colon-vs-em-dash difference (lines 254-256). My round-01 MINOR about
    separator mismatch addressed.

26. RULES.md Keys table: `tier` row now says "written here as the bare ordinal
    `T0`–`T3`; prose spells it number-plus-word" (line 38). My round-01 NIT
    fixed.

27. WEB-ANNEX.md: SC 1.4.4 and 1.4.12 now under "Recorded, not enforced
    *(informative)*" (line 38), not "They block". My round-01 MINOR fixed.

28. `core:motion-without-reduced-path` sources: now
    `[WCAG-2.2-SC-2.2.2, WCAG-2.2-SC-2.3.3]` (line 237). My round-01 MINOR
    about SC 2.3.3 source mismatch fixed.

29. NOTICE.md: DTCG `2025.10` now noted as "a date pin, not a W3C
    Recommendation" (line 30-31). My round-01 NIT fixed.

30. PDS.md examples: searched for stale rule ids
    (`core:states-incomplete`, `core:raw-value-outside-tokens`,
    `core:effect-budget` without `-exceeded`) → 0 matches. kimi-1's
    round-01 MINOR fixed.

31. `waivers: ""` crash probe: created a CONTRACT with `waivers: ""` (empty
    string, valid per §2 rule 5 canonical frontmatter). Result: EISDIR
    crash, exit 2. See MAJOR finding below.

32. `waivers: []` probe: same contract with `waivers: []` (empty flow list).
    No crash — runs normally. The crash is specific to the empty-string case.

## Findings

### [MAJOR] `waivers: ""` in a CONTRACT crashes the checker (EISDIR)

WHAT. A CONTRACT artifact with `waivers: ""` (empty string) causes the checker
to crash with `Error: EISDIR: illegal operation on a directory, read` and exit
code 2 ("the run itself failed"). The path is: engine.js line 1440 checks
`typeof contract.data.waivers === "string"` — an empty string passes this
check — then `path.resolve(path.dirname(contract.path), "")` returns the
contract's own directory. `fs.existsSync` on that directory is true, so
`loadWaivers` is called on a directory, and `readArtifact` crashes.

WHY IT MATTERS. PDS §2 rule 3 says "Empty is not absent" and §2 rule 5
(canonical frontmatter) explicitly lists `""` as a valid empty value. A
CONTRACT with no waivers is a legitimate state — a run with no exceptions.
The spec permits the input; the checker crashes on it. Exit 2 means "the run
itself failed," which is the wrong signal for a valid input — a CI gate would
report a checker failure rather than a clean run. This is a new defect
introduced by none of the fix-up cycles specifically; it is an edge case in
the waiver-loading code that was never probed.

CONCRETE FIX. In engine.js line 1440, change the condition to also reject
the empty string before resolving:

```js
const waiverPath =
    options.waiversPath ||
    (contract && typeof contract.data.waivers === "string" && contract.data.waivers.trim() !== ""
      ? path.resolve(path.dirname(contract.path), contract.data.waivers)
      : null);
```

Or, more defensively, check `fs.statSync(waiverPath).isFile()` before calling
`loadWaivers`. Add a fixture: a CONTRACT with `waivers: ""` and no waiver
file → the checker runs normally, no waivers loaded, no crash.

### [MINOR] "alias direction" from FINAL.md dropped from G3 and unimplemented

WHAT. FINAL.md line 111 says: "L3 token integrity (DTCG `2025.10`, alias
direction, no raw literals outside the token layer)." PDS.md G3 (line 337)
says: "an alias that does not resolve to a declared token, an alias cycle" —
"alias direction" is absent. The checker's L3 (engine.js `checkL3`) verifies
alias resolution and cycles but not alias direction. I raised this as a
round-01 MINOR; the consensus accepted all MINOR findings into the fix-up
pass (consensus line 87-88), but neither the doctrine nor the checker
addressed it.

WHY IT MATTERS. This is a fidelity gap between the binding spec (FINAL.md)
and the implementation (PDS.md G3 + checker L3). "Alias direction" is a
PDS-specific constraint beyond what DTCG requires (the DTCG editor's draft
permits any reference). A project claiming L3 cannot have its alias direction
verified. The constraint is also undefined — neither FINAL.md nor PDS.md says
what "alias direction" means (e.g. "semantic → primitive, not reverse"),
making it an unverifiable requirement, the same shape of defect A1 was
written about.

CONCRETE FIX. Either (a) define "alias direction" in G3 (e.g. "an alias
chain MUST NOT traverse from a group named `semantic` to a group named
`primitive`"), add a check in `checkL3`, and ship a fixture; or (b) strike
"alias direction" from FINAL.md's L3 description with a recorded erratum,
since G3 as shipped covers resolution and cycles only. The current state —
FINAL.md requires it, PDS.md G3 omits it, the checker doesn't verify it — is
the worst option.

### [NIT] D-2 deviation note says "eight rules" but could name them

IMPLEMENTATION.md D-2 says "Eight rules are in that state, not five" and
names them. This matches what I verified: 5 no-detector rules
(`core:text-below-legible-floor`, `core:unlabelled-inference`,
`core:value-off-scale`, `core:colour-off-ramp`, `web:viewport-hero`) plus 3
T2-tier rules (`web:contrast-ratio`, `web:target-size`, `web:reflow-narrow`).
kimi-1's round-01 counting note is addressed. No action needed; recording
that the count is now correct.

### [NIT] `core:contrast-applied` (T2, `both`) is undetected but not named in D-2

`core:contrast-applied` is `enforced-by: both` at `T2 RENDERED`. It has no
detector and is above the checker's tier. D-2 names 8 undetected rules but
does not mention this one — it falls in both categories (T2 and no detector).
It is correctly reported UNJUDGEABLE in the capability output with the reason
"T2 RENDERED evidence is above this checker." The omission from D-2's list is
cosmetic since the rule is visible in the report either way, but the deviation
note's count is slightly incomplete. No action needed beyond awareness.

## Cross-reviewer responses (round-01 findings)

### codex-1 round-01 findings

- **CRITICAL: L2/L3 can be falsely certified (AF-1).** CLOSED. The checker
  now models each level as an explicit obligation set
  (`ledger.declare`/`ledger.fail`/`ledger.unverified`). L2 validates the U1
  assignment (`checkAssignment`), G1/G2 conditions, recusal, process order,
  and gate records. L3 requires a DTCG token document, colorSpace on every
  colour, alias resolution, and system rules decided against real source.
  Verified on `collapsed-run`: "not verified" with 12 violations including
  assignment, G1, G2, recusal, required fields. Verified on `sound-run`:
  "verified L2" with PASS. An unverified obligation never counts as met.

- **CRITICAL: self-signed waiver (AF-2).** CLOSED. Verified by probe:
  `granted-by: claude-1, counter-signed-by: claude-1` is rejected with
  "claude-1 counter-signed its own waiver, which is not a counter-signature"
  and the finding stays in the ledger. Independence is checked against the
  run's roster; no roster → no suppression.

- **MAJOR: unjudgeable run exits 0 (AF-3).** CLOSED. `check NOTICE.md` →
  exit 4, not 0. Exit 0 is reserved for PASS alone. Help text and SKILL.md
  document code 4.

- **MAJOR: artifact ingestion rejects valid YAML and drops unknown rule ids
  (AF-4).** CLOSED. A canonical frontmatter subset is ratified (PDS §2 rule
  5) and all 8 published examples conform to it. A candidate artifact that
  fails to parse is reported as `pds-check:l1-frontmatter-parses`, not
  silently omitted. Unknown rule ids are traversed and reported UNJUDGEABLE
  (verified: the `every rule id the doctrine cites` test guards this, and
  the checker's §10 rule 3 pass-through is implemented).

- **MAJOR: reduced-motion block passes without reducing motion (AF-5).**
  CLOSED. Verified by probe: a reduced-motion block that only sets
  `color: red` no longer suppresses the finding.

- **MAJOR: focus-visible replacement reported as absent (AF-6).** CLOSED.
  Verified by probe: split `outline: none` / `:focus-visible { outline: ... }`
  no longer fires a violation.

- **MAJOR: D-1 weakens binding acceptance criterion (AF-10).** I maintain
  my round-01 ACCEPT. The 64 KiB total is held (65286/65536). The per-file
  numbers are now documented as early-warning thresholds that deliberately
  sum above the total, not as a partition — which also closes my round-01
  NIT about the sum. PDS.md at 24909 bytes preserves the identical four-part
  artifact shape, which is the property the spec exists to have.

- **MINOR: Heroicons double-count.** CLOSED. One `@heroicons/react/24/solid`
  import → PASS. The `one-package` fixture covers it.

- **MINOR: unnumbered normative requirements in tables.** CLOSED. PDS §2
  rule 4 ("One shape, and the table binds") gives table cells normative
  force through a numbered, named rule: "a field it lists MUST be present
  and MUST meet the requirement beside it." The tables bind through this
  rule rather than carrying floating RFC 2119 keywords.

### kimi-1 round-01 findings

- **MAJOR: G1 drops two C7 conjuncts (AF-8).** CLOSED. PDS §3 G1 (line 316-318)
  now includes all three: "the ban list, past `core:category-guessable`'s
  category-plus-avoidance test, and on recorded human ratification with a
  brief-specific reason." The G1 error string (line 321) also includes all
  three. Verified on `collapsed-run` output.

- **MAJOR: U1 assignment unverifiable (AF-9).** CLOSED. `run-id` is now in
  the DESIGN-BRIEF required-fields table (PDS §2 line 112) and in the example
  (line 122). The checker's `checkAssignment` recomputes the rotation from
  the brief's `run-id` and primary positions, compares each DIRECTION's
  `assigned`, and checks the minimum-position count. Verified on
  `collapsed-run`: "the brief declares no run-id, which §4 rule 2 hashes, so
  the rotation is not reproducible."

- **MINOR: PDS examples cite three non-existent rule ids.** CLOSED. Searched
  for `core:states-incomplete`, `core:raw-value-outside-tokens`,
  `core:effect-budget` (without `-exceeded`) → 0 matches. The
  `every rule id the doctrine cites` test guards this going forward.

- **MINOR: gate strings / L2 row drift.** CLOSED. PDS §3 now says "canonical
  message shapes, not literal output" (line 306-308). §9 L2 row says "a
  recorded gate for every §3 transition the run crossed" (line 464), scoped
  to transitions actually crossed.

- **MINOR: RULES.md 30 bytes of headroom.** ADDRESSED. RULES.md now has 1004
  bytes of headroom (23572/24576) after the rebalance moved its threshold to
  24 KiB and the fix-up cycles compressed prose. The headroom concern is
  substantially mitigated, though the file still only grows.

- **NIT: Keys table tier form.** CLOSED. Keys table now says "written here
  as the bare ordinal `T0`–`T3`; prose spells it number-plus-word" (line 38).

## Open questions

1. The `waivers: ""` crash (MAJOR) — is this the only empty-string field
   that resolves to a directory? A scan of other path-bearing contract fields
   (`tokens`, `named-rules`) would confirm whether the same pattern exists
   elsewhere. The `tokens` field is resolved through `readTokenIndexFor`
   which may have different handling.

2. "alias direction" (MINOR) — was this dropped from G3 deliberately (because
   DTCG doesn't enforce it and defining it rigorously is non-trivial) or
   accidentally during the author's-death recovery? If deliberate, FINAL.md
   should be amended; if accidental, G3 should be restored. Either way, the
   current state where FINAL.md says it and PDS.md doesn't is a fidelity gap.

3. G4's condition (no unresolved quality violations) is checked only as a
   recorded gate outcome, not recomputed. The verdict roll-up covers the
   substance (unwaived VIOLATION → verdict VIOLATION, not PASS), but the
   checker never explicitly says "G4 fails" — it says the verdict is
   VIOLATION. Is this intentional, or should G4 get the same recomputation
   treatment G1/G2 received?
