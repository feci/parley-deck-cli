---
agent: codex-1
idea: parley-design-skills
review-round: 02
date: 2026-07-28
reviewed-commit: 8fc3a18
---

## Summary

❌ BLOCK. The fix-up is substantial and most original defects are genuinely closed, but
AF-1 and AF-2 are not. I reproduced new paths to clean L1/L2/L3 certificates without the
evidence those levels claim, and a waiver can still be counter-signed by the author of the
waived work. These retain the same CRITICAL impact as round 01.

Agreed-fix disposition:

- **AF-1: not closed.** The explicit obligation ledger is real, but several obligations are
  still marked met without validating their binding inputs or conditions.
- **AF-2: not closed.** Equal identities, missing grantors and unknown signers are rejected;
  the required “counter-signer is not an author of the waived work” check is absent.
- **AF-3: closed.** An uninspected run now reports `UNJUDGEABLE` and exits 4; help and the
  checker skill document the code.
- **AF-4: closed for the original defect.** Every published example parses and unknown cited
  ids become `UNJUDGEABLE`; a separate candidate-kind omission remains under AF-1 below.
- **AF-5: closed.** A reduced-motion block changing only an unrelated property no longer
  counts as coverage.
- **AF-6: closed.** The split `outline: none` / `:focus-visible` pattern no longer raises
  `core:focus-indication` (the fixture still raises an unrelated interaction-state finding).
- **AF-7 (hermes-1 MAJOR): genuinely closed.** `RULES.md` now derives the ban list from
  `slop` rules at `T0 ARTIFACT`, defines the signature and sharing test, and the checker uses
  that definition.
- **AF-8 (kimi-1 MAJOR): genuinely closed.** PDS §3 restores the ban list,
  category-plus-avoidance test and recorded human ratification as conjuncts.
- **AF-9 (kimi-1 MAJOR): not genuinely closed.** `run-id` and rotation recomputation landed,
  but duplicate primary values can assign the same position to two proposers, and invented
  non-brief axes can satisfy G1. Both probes verified L2.
- **AF-10: closed by accepted reasoning.** I withdraw my round-01 D-1 block. The review
  consensus resolves the 64 KiB aggregate as binding, the early-warning thresholds are now
  explicit, and the actual total is 65,286 of 65,536 bytes. D-2 is also acceptable: every
  missing detector remains visible as `UNJUDGEABLE`.

## What I verified (commands run, and their result)

- `git` branch/HEAD/status checks in the implementation repo: branch
  `parley-design-skills`, implementation HEAD `8ebd8f7`, clean worktree. Deck commit
  `8fc3a18` contains the cycle-3 review handoff.
- `npm test`: **197 passed, 0 failed**. This independently confirms the declared suite
  count, registry refusal, generated capability tests, original AF regressions, detector
  fixture pairs, installer coverage, byte budgets, digest guard and placeholder scan.
- `npm pack --dry-run --json --cache <temporary-cache>`: succeeded; 141 package entries,
  including `NOTICE.md`, all four doctrine files, and the complete checker tree.
- `node .../check.js --help`: exit 0 and documents exit 4. `--json NOTICE.md`: verdict
  `UNJUDGEABLE`, process/report exit **4**.
- Capability probe over `lib/detectors`: 18 detector files, 18 unique rule ids, no duplicate
  capability claim. Source inspection confirms `loadDetectors()` derives it from the
  directory.
- `wc -c` and `shasum -a 256`: doctrine sizes 6,681 / 24,909 / 23,572 / 10,124 bytes,
  total 65,286; `RULES.md` begins `f0c38eed1b8d`, matching PDS frontmatter.
- Direct reduced-motion fixture run: `core:motion-without-reduced-path` VIOLATION, exit 1.
  Direct split-focus fixture run: no `core:focus-indication` finding.
- Mutated copies of the sound run under `/tmp/pds-codex1.6XgLyH`:
  - swapped DIRECTION files into `round-02` and CRITIQUE into `round-01` → **L2 verified,
    PASS, exit 0**;
  - added an invented direction axis so the pair differed on only one brief axis → **L2
    verified, PASS, exit 0**;
  - used primary values `[flat, flat, layered]` and recorded `flat` for both proposers,
    while keeping two other differing axes → **L2 verified, PASS, exit 0**;
  - introduced a raw colour while the recorded G3 outcome stayed `pass` → the report found
    the violation but still certified **verified L2** with an empty unmet set;
  - added an AUDIT recording G4 `pass` plus open quality findings → still **verified L2**
    with an empty unmet set;
  - appended `.trap::before { content: "}"; color: #ff0000; }` → the source parser lost the
    colour after the quoted brace and returned **verified L3, PASS, exit 0**;
  - added a parseable `spec: PDS/1.0` file with `kind: TYPO-BRIEF` → it was listed as
    `not-inspected`, while **L1 verified, PASS, exit 0**;
  - waived a violation on `round-01/claude-1.md` with `granted-by: codex-1` and
    `counter-signed-by: claude-1` → one waiver applied, the finding disappeared, **PASS,
    exit 0**, even though `claude-1` authored the scoped work.

I did not independently perform a textual similarity comparison against the two prior-art
repositories. `NOTICE.md` is present and accurately states the declared attribution policy.

## Findings

### [CRITICAL] AF-1 remains open: L1/L2/L3 certificates are still forgeable

**What is wrong.** The obligation ledger records only what its checks happen to test:

- G1 counts the union of keys supplied by DIRECTION files, not the brief's declared axes,
  and it never validates each position against the brief's enumeration. An invented axis
  therefore supplies the second difference required for a clean L2 result.
- Assignment checks the count of distinct primary values but rotates the original list.
  Duplicates can assign the same value to multiple proposers while the obligation reports
  `met`.
- “Process order” checks artifact presence, not the normative §1 locations; swapping
  `round-01` and `round-02` artifacts still verifies.
- When a run crossed G3 or G4, L2 validates only the recorded word `pass`. Relevant raw-token
  and open-quality violations do not sink the L2 obligation, contrary to AF-1's agreed
  requirement to validate every applicable gate condition.
- The CSS stack machine treats braces and semicolons inside quoted strings as syntax. A
  common `content: "}"` hides following declarations, including a raw colour, and yielded a
  clean L3 certificate.
- A parseable file declaring `spec: PDS/1.0` but an invalid `kind` is demoted to
  `not-inspected`, so L1 can verify beside an invalid candidate artifact.

**Why it matters.** These are not report-polish defects. The tool emits `verified L2` or
`verified L3`, `PASS`, and exit 0 while the binding process or token-integrity conditions
are false. That is the exact false-certificate failure AF-1 was meant to close.

**Concrete fix.** Validate DIRECTION position keys exactly against the brief axes and each
value against its enumeration; reject duplicate primary positions (or rotate a deduplicated
list) and count G1 differences only on brief axes. Validate artifact paths against §1's
mapping. When G3/G4 transitions were crossed, make their relevant rule outcomes sink the
level obligation. Make the CSS scanner quote/escape-aware before treating `{`, `}`, `;` or
`:` as structure. Treat every file declaring `spec: PDS/1.0` as a candidate and report an
unknown/missing kind as an L1 violation. Add each reproducer above as a negative fixture.

### [CRITICAL] AF-2 remains open: the waived work's author can counter-sign

**What is wrong.** `waiverProblem()` proves that grantor and signer are distinct strings in
the run's inferred participant set, but it never compares the signer with ownership of the
scoped work. My waiver for `round-01/claude-1.md`, granted by `codex-1` and counter-signed
by `claude-1`, suppressed the violation and exited 0.

**Why it matters.** PDS §8 rule 2 requires a counter-signer who is neither the grantor nor an
author of the waived work. The author can currently approve their own exception merely by
making another participant the grantor, so the independence guarantee remains bypassable.

**Concrete fix.** Resolve an artifact scope to its protocol owner and reject any
counter-signer matching that owner. For source whose ownership cannot be established from
the supplied artifacts, require a trusted roster/ownership input or leave the finding
unsuppressed, as §8 already requires. Add fixtures for signer-as-scoped-artifact-author and
for a signer present only through a self-authored identity artifact.

### [MAJOR] The single-agent fast path contradicts the protocol it says wins

**What is wrong.** `SKILL.md` and PDS §4 rule 8 prescribe a one-agent fast path while PDS
§0 rule 2 says `COOPERATION.md` wins on process, and that protocol makes non-solo execution
mandatory for any claimed Parley workflow. FINAL.md explicitly defers a core-protocol
carve-out.

**Why it matters.** A user following the add-on literally is instructed both to load it
alongside Parley Deck and to violate Parley's quorum rule. The audit trail cannot tell
whether the run is a valid fast path outside Parley or an invalid solo Parley run.

**Concrete fix.** For v1, state explicitly in both SKILL.md and PDS that the one-agent fast
path is outside a Parley Deck workflow and MUST NOT claim Parley verification. Any actual
solo carve-out must instead go through the deferred meta-protocol-change idea.

### [MINOR] The binding L3 “alias direction” capability is still absent

FINAL.md requires alias direction as part of L3 token integrity. The checker verifies only
that aliases resolve without cycles, and PDS §9 now omits “alias direction” rather than
defining it. This silently weakens the binding acceptance criterion. Define the permitted
direction (for example semantic-to-primitive), enforce it with a negative fixture, or
ratify a spec amendment before removing the requirement.

### [NIT] None.

## Open questions

None. The blockers have reproducible fixes and need no owner decision.
