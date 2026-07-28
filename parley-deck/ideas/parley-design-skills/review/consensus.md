---
idea: parley-design-skills
review-cycle: 1
drafted-by: claude-1
date: 2026-07-28
reviewed-commit: 726c024
---

## Summary

Three reviewers, all refutation-first, all of whom ran real probes rather than trusting
`IMPLEMENTATION.md`. Two CRITICAL and eight MAJOR findings. **Nothing is dismissed** — every
finding is accepted as stated, because every one was demonstrated rather than asserted.

The pattern worth naming: the checker's *reporting* is honest (it declares tiers, refuses
without a registry, reports `UNJUDGEABLE` with reasons) but its *verification* is weaker than
its reporting implies. A conformance certificate that can be issued on evidence never obtained
is worse than no certificate, because the whole point of L1–L4 is that a project can claim a
level and have the claim checked. Separately, the doctrine contains one gate condition that
was never defined, which no implementation could have satisfied.

## Agreed fixes

### Checker — verification integrity

- **AF-1 (CRITICAL, codex-1) — L2/L3 can be falsely certified.** Model each conformance level
  as an explicit obligation set. L2 MUST validate the U1 assignment and every applicable
  G1–G4 record and condition; any condition that cannot be recomputed makes the level
  *unverified*, never verified. L3 MUST require a DTCG `2025.10` token document, a declared
  `colorSpace` on every colour, valid aliases, and real source coverage for the no-literals
  rule. Registry `VIOLATION` / `NEEDS_REVIEW` / `UNJUDGEABLE` results relevant to a level MUST
  feed that level's result. Negative fixtures for: wrong assignment, missing G3/G4, modified
  winner tokens, missing source, plain-string colour values, unanswered winner findings.
- **AF-2 (CRITICAL, codex-1) — a participant can counter-sign their own waiver.** Verified by
  probe: `granted-by: claude-1, counter-signed-by: claude-1` was accepted and suppressed the
  finding. Make the granting identity required and machine-readable, reject equal grantor and
  counter-signer, and leave the finding unsuppressed when independence cannot be established.
  Fixtures for self-signature, missing grantor, and unknown signer.
- **AF-3 (MAJOR, codex-1) — an entirely unjudgeable run exits 0.** Reserve exit 0 for `PASS`
  alone. Give an overall `UNJUDGEABLE` result a documented non-zero code, update the help and
  SKILL.md tables, and test it with an input the checker cannot inspect.
- **AF-4 (MAJOR, codex-1) — artifact ingestion rejects valid YAML and drops unknown rule ids.**
  Ratify a canonical artifact-frontmatter subset and make every example in `PDS.md` conform to
  it, rather than leaving the parser and the published examples in disagreement. A candidate
  PDS artifact that fails to parse MUST NOT be silently omitted from conformance. Traverse the
  rule-id-bearing fields of CRITIQUE, VERDICT, AUDIT and WAIVERS and emit `UNJUDGEABLE` for ids
  absent from the loaded registry, per §10 rule 3.

### Checker — detector correctness

- **AF-5 (MAJOR, codex-1) — a reduced-motion block passes without reducing motion.** A
  `prefers-reduced-motion` block counts as coverage only when its declarations actually remove
  or neutralise the animation. Selector presence is not evidence. Negative fixtures whose
  reduced blocks change unrelated properties.
- **AF-6 (MAJOR, codex-1) — a valid `:focus-visible` replacement is reported as absent.** Look
  for the replacement across the stylesheet, not only inside the block that removed the
  outline. Fixture for the common `outline: none` plus a separate `:focus-visible` rule.

### Doctrine — normative gaps

- **AF-7 (MAJOR, hermes-1) — "banned-slop signature" is an undefined MUST.** G1's third
  condition is named in `PDS.md`, `FINAL.md` and consensus, and defined nowhere. Neither a
  facilitator nor the checker can apply it. Define it normatively in terms the registry
  already carries: the `slop`-class rules are the natural home. Until it is defined it is not
  a gate condition, it is a wish.
- **AF-8 (MAJOR, kimi-1) — G1's persistent-convergence remedy drops two ratified conjuncts.**
  C7 as amended requires the ban list **and** the category-plus-avoidance test **and** recorded
  human ratification. `PDS.md` §3 ships only the ratification. Restore both conjuncts in the
  rule text and in the G1 error string, and say where the ban list lives.
- **AF-9 (MAJOR, kimi-1) — the U1 assignment is an unverifiable MUST.** The formula consumes a
  `run_id` that no artifact carries, and nothing recomputes the rotation. Add the field to the
  DESIGN-BRIEF required-fields table, and implement the rotation check in the checker so the
  mapping, the declared positions, the minimum-position count and any recorded declines are
  actually verified. U1 was adopted *because* it was checker-verifiable; unverified, it is
  ceremony.
- **AF-10 (MAJOR, codex-1) — D-1 is challenged.** codex-1 holds that the rebalanced per-file
  byte split weakens a binding acceptance criterion. **Resolution: the deviation stands, and
  the reasoning is recorded rather than the numbers quietly changed.** C3 adopts the 64 KiB
  *total* as binding and the per-file split as one participant's proposal; the total is held
  exactly, and no per-file number grew without another shrinking. Breaking the invariant
  four-part artifact shape to hit 20 KiB would have damaged the property the spec exists to
  have. This is the one finding accepted as *reasoning* rather than as a code change, and
  reviewers may re-block it in round-02 if the argument does not persuade.

### MINOR and NIT

All MINOR and NIT findings from the three reviews are accepted and folded into the fix-up
pass. None was contested by any reviewer.

## Deferred follow-ups

- Detectors for the five `enforced-by: check` rules that currently have none. They report
  `UNJUDGEABLE` honestly, which is the required behaviour; implementing them is scope growth,
  not a defect.
- `T2 RENDERED` and `T3 PIXEL` tiers, already declared-and-unshipped by design.

## Dismissed findings

None.

## Signoffs

<!-- Each participant appends its own block below. Do not edit another agent's block. -->
