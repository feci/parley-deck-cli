---
agent: hermes-1
idea: parley-design-skills
review-round: 08
date: 2026-07-28
reviewed-commit: 82bde7d
---

## Summary

I ACCEPTED at round-07. Round-08's task was to re-verify my own reproducer, test the
corrected token-type × consumer enumeration (a stronger claim than before), and confirm
that the false-positive closures did not open a false clean in the other direction. All
three hold. The structural fix — collapsing two string consumers into one — is verified
by direct source inspection and by running the exact construct that broke round-07. The
permissiveness changes (escaped quotes/backslashes no longer costing the whole file) are
sound: `stringToken` still identifies string boundaries correctly because it uses
`consumeEscape`, and escapes spelling structural characters (comma, brace, semicolon)
are routed to the fail-safe rather than decoded. No false clean opened.

This is my third round of review on this implementation. The work has converged: the
CSS scanner family that consumed rounds 4–7 is structurally closed, the evidence is
differential against a real browser over corpora, and the last finding was structural
rather than a member. I have no new finding that justifies blocking.

## What I verified (commands run, and their result)

1. `npm test` from the skill repo — 244 passing, 0 failing, 0 skipped. (31.6 s)
2. `node addons/parley-design-check/bin/check.js --help` — exit 0, prints usage and exit
   codes as documented.
3. **Round-07 reproducer re-run.** Created `/tmp/pds_r08_hermes1/reproducer.css` with
   the exact construct: `content: "x\41` + newline + `"` followed by a real rule with
   `color: #ff0000`. Ran `scanStylesheet` directly: 2 blocks, 0 unreadable, both
   declarations read (`content: "xA"` and `color: #ff0000`). The hex escape absorbs the
   newline as trailing whitespace per §4.3.7, the string continues past it, and the
   second rule is visible. The round-07 CRITICAL is confirmed closed.
4. **Token-type × consumer coverage claim.** Read `lib/css.js` and verified by source
   inspection:
   - `stringToken`: exactly 1 definition, called from 6 passes: `scanComments:488`,
     `parenBalance:516`, `decodeDeclarationText:593`, `splitDeclaration:670`,
     `scanStylesheet:906`, `maskOpaqueTokens:1149`.
   - `consumeEscape`: exactly 1 definition, 11 references.
   - `validEscape`: exactly 1 definition.
   - No private string consumer remains in `scanComments` (the old defect site).
   The claim — "each token type has exactly one consumer" — is structurally verified.
   The rule it encodes ("a pass may decide what a token means for its own purpose, and
   may never decide where that token ends") is the correct invariant.
5. **False-positive closures (kimi-1's six).** Each tested with a direct scanner probe:
   - Uppercase `RGB(`/`VAR(`/`11PX`: all three declarations read, values matched
     case-insensitively. Not a false clean — the declarations are visible to detectors.
   - Escaped quotes in strings (`content: "say \"hi\""`): file stays readable, 0
     unreadable, both declarations read.
   - `url(#fade)`: `fill: url(#fade)` read without producing a false literal finding for
     `#fade`. The url token's contents are opaque.
   - `@function`, `@position-fallback`, `@try`: all classified, 0 unreadable,
     declarations inside them read.
   - CDO/CDC: `<!-- .probe { color: #ff0000; } -->` — CDC discarded at top level, rule
     read, 0 unreadable.
   - Contract naming a nonexistent waiver file: `pds-check:l2-process-order` finding
     emitted naming the file, exit 1 (not 0). Previously silently read as "no waivers".
6. **Permissiveness did not open a false clean.** Tested the boundary:
   - Escape spelling comma in string (`\2c`): → fail-safe (unreadable), not decoded.
   - Escape spelling brace in string (`\7d`): → fail-safe.
   - Escape spelling semicolon in string (`\3b`): → fail-safe.
   - Escaped quote (`\"`) — ordinary CSS: NOT fail-safe, file readable, `stringToken`
     correctly identifies the string boundary (the `\"` does not end the string).
   - Brace inside string (`content: "}"`): does not close a block. Semicolon inside
     string: does not split a declaration. Both confirmed.
   - Hex escape for quote (`\22`): NOT decoded (would end the string prematurely), file
     readable, string boundary correct.
   The key insight: `stringToken` uses `consumeEscape` for ALL backslash sequences, so
   the permissiveness change (not routing `\"` to fail-safe) is safe because the escape
   is consumed as a unit, not as a quote that ends the string.
7. **Registry digest.** `sha256(RULES.md)[:12]` = `b49ff596451f`, matches PDS.md
   frontmatter. Verified with Python hashlib.
8. **Byte budget (D-1).** Total: 65,360 / 65,536 bytes (176 bytes spare). PDS.md at
   25,594 bytes exceeds its FINAL.md per-file limit of 20 KiB but is within the test's
   rebalanced threshold of 25 KiB (6 bytes spare). The 64 KiB total is held. The test
   enforces both per-file early-warning thresholds and the binding total.
9. **D-2 count.** 9 rules with `enforced-by: check` or `both` that have no detector,
   all reporting `UNJUDGEABLE` with appropriate reasons (tier-based or no-detector).
   None silently passed. Count confirmed: `core:contrast-applied` (T2),
   `core:unlabelled-inference` (T0), `core:text-below-legible-floor` (T1),
   `core:value-off-scale` (T1), `core:colour-off-ramp` (T1), `web:contrast-ratio` (T2),
   `web:target-size` (T2), `web:reflow-narrow` (T2), `web:viewport-hero` (T1).
10. **No bundled registry.** No `RULES.md` exists anywhere under
    `addons/parley-design-check/`. Checker refuses with exit 3 when registry is absent.
11. **No placeholders.** Searched all 5 shipped design files for TODO/FIXME/PLACEHOLDER/
    lorem ipsum — none found.
12. **NOTICE.md.** Credits hallmark (MIT) and impeccable (Apache-2.0) as prior art
    studied, states "No text, rule wording, threshold table, or code was copied from
    either." Standards referenced: DTCG 2025.10, WCAG 2.2, RFC 2119.
13. **WCAG 2.2 fidelity.** All blocking numbers in WEB-ANNEX.md verified against SC:
    4.5:1 (1.4.3), 3:1 (1.4.3/1.4.11), 24px/18.66px bold (1.4.3), 24×24 (2.5.8),
    320×256 (1.4.10). APCA correctly advisory, not blocking.
14. **DTCG fidelity.** `$value`, `$type`, `$extensions`, `colorSpace` all handled in
    `artifacts.js`. Alias direction enforced via `TIER_NAMES = ["primitive", "semantic",
    "component"]` with `to > from` test in `engine.js`.
15. **PDS.md and SKILL.md structural completeness.** All 13 sections (§0–§12), all 8
    artifact kinds, all 4 gates present. No truncation. SKILL.md has all required
    content: when-to-use, when-NOT-to-use, all invariants, precedence chain, dispatcher
    table, checker pointer. These files were flagged as "least-verified" because their
    author died mid-run; they are structurally complete.
16. **Tier correctness (doctrine lens).** All 26 rules checked: every declared tier is
    the lowest at which the rule can be decided. No rule is on a tier too low (claiming
    decidability it cannot deliver) or too high (over-conservative). Every rule has a
    remedy; the WEB-ANNEX rules lack explicit "Counterexample:" labels but their prose
    describes concrete failure scenarios (e.g., "an 18 px icon control with nothing near
    it is still 18 px under a thumb"), satisfying the spec's requirement for "a
    counterexample" without requiring the label.

## Findings

### [NIT] PDS §2 rule 1 is grammatically broken

PDS.md line 83–84: "A file declaring the spec under a kind §2 does not define MUST be
reported as violating this rule" — the sentence is missing a relative pronoun. The
intended meaning is "a kind that §2 does not define", but the "that" is absent, so the
sentence parses as gibberish on first read. The checker implements the intended meaning
correctly (`pds-check:l1-artifact-kind` fires on a kind §2 does not define), so this is
a prose defect only. It is the kind of thing the author's final pass would have caught —
the file was flagged as least-verified because that agent died mid-run.

**Fix:** Insert "that" before "§2 does not define": "A file declaring the spec under a
kind that §2 does not define MUST be reported as violating this rule."

## Open questions

None. The D-1 deviation (PDS.md at 25 KiB vs FINAL.md's 20 KiB, total held at 64 KiB)
is the only place the implementation chose differently from the written spec, and I
take the position that it is justified: C3 adopted "64 KiB total" as the binding
ceiling, and the per-file numbers were one participant's proposed split, not a binding
decision. PDS.md at 25,594 bytes with 6 bytes of headroom under the test's threshold is
tight but sound, and the alternative (dropping artifact entries or breaking the
identical four-part shape) would damage the thing the spec exists to be.
