---
agent: hermes-1
idea: parley-design-skills
review-round: 09
date: 2026-07-28
reviewed-commit: 9121ec2
---

## Summary

I ACCEPTED at round-07 and round-08 with no new findings. Round-09's task was to
re-verify my own round-08 reproducer and its controls, and to check that cycle 12's
@function parameter-list fix did not swing back toward a false alarm — the failure
mode the last two cycles have been trading against.

All three hold. The round-08 CRITICAL (codex-1 and kimi-1, independently) is closed:
`@function --pick(--x: var(--ghost)) { result: var(--ghost); }` now flags `--ghost`
twice (once from the default, once from the body), returns VIOLATION, and exits 1.
The fix did not swing toward a false alarm: the clean controls (true formal with
declared default, true formal with no default) still PASS at exit 0, the Sass
`@function` false-alarm direction is correctly suppressed at the parameter-list
level, and a non-custom-property parameter name goes to the unreadable channel
rather than being guessed in either direction.

This is my fourth round of review. I have no new finding that justifies blocking.
codex-1 filed a MAJOR (quoted `var()` text in `@function` defaults counted as a
real token use); I independently reproduced it and confirmed it is real, but it is
pre-existing — the same unmasked `VAR_REFERENCE` pattern exists in
`declarationVarUses` for ordinary declarations since at least cycle 9 — and it is
the false-alarm direction, not a false clean. I describe it below as a MAJOR with
that context. Two carried NITs and one carried MINOR also remain, all from prior
rounds, none introduced by cycle 12, none blocking.

## What I verified (commands run, and their result)

1. `npm test` from the skill repo — 245 passing, 0 failing, 0 skipped. (19.9 s)
2. `node addons/parley-design-check/bin/check.js --help` — exit 0, usage and exit
   codes as documented.
3. **Git state.** Branch `parley-design-skills`, HEAD `9121ec2`, worktree clean.
   `git diff 82bde7d..9121ec2 --stat` shows only `lib/css.js` (151 insertions, 16
   deletions) and `test/checker.test.js` (163 insertions). Doctrine
   (`addons/parley-design/`) is byte-identical between the two commits.
4. **Round-08 CRITICAL reproducer (codex-1/kimi-1).** Copied the sound-run fixture
   to `/tmp/pds_r09_hermes1/sound-run`, appended the exact probe:
   `:root { --ghost: rgb(255,0,0); } @function --pick(--x: var(--ghost))
   { result: var(--ghost); } .probe { color: --pick(); }`.
   `check --level L3` → `core:token-used-undeclared` flags `--ghost` at line 17
   (twice — once from the default expression, once from the body), verdict
   VIOLATION, exit 1. The round-08 CRITICAL is confirmed closed.
5. **Clean control: true formal + declared default.** Appended
   `@function --pick(--x: var(--color-text-body)) { result: var(--x); }
   .probe { color: --pick(); }` → verdict PASS, verified L3, exit 0. No false
   alarm from the fix.
6. **Clean control: true formal, no default.** Appended
   `@function --double(--x) { result: calc(var(--x) * 2); }
   .probe { color: --double(); }` → verdict PASS, verified L3, exit 0. The body's
   `var(--x)` is correctly suppressed as a bound formal.
7. **Default-only half (kimi-1's second half).** Appended
   `@function --pick(--x: var(--ghost)) { result: 1px; }
   .probe { color: --pick(); }` → VIOLATION, exit 1. The default's `var(--ghost)`
   is collected as a real use even when the body references nothing. This was the
   half that was invisible before the fix.
8. **Escaped `var()` in default.** Appended
   `@function --pick(--x: \76 ar(--ghost)) { result: 1px; }
   .probe { color: --pick(); }` → VIOLATION, exit 1. The escaped `\76 ar` is
   decoded to `var` by `decodeDeclarationText` and the reference is collected.
9. **Body references undeclared token not in prelude.** Appended
   `@function --pick(--x) { result: var(--ghost); }
   .probe { color: --pick(); }` → VIOLATION, exit 1. `--ghost` is not a formal, so
   the body reference is not suppressed.
10. **Non-custom-property parameter name.** Appended
    `@function --f(x: 1) { result: var(--x); }` → unreadable channel:
    "the @function parameter list at line 21 carries "x: 1", which this scanner
    cannot read as a parameter". Exit 1. Not guessed either way — the segment goes
    to unreadable rather than being read as a formal or a use.
11. **Sass @function false-alarm direction.** Ran a direct scanner probe:
    `scanStylesheet("@function px-to-rem($n) { @return $n / 16 * 1rem; }
    .probe { color: var(--color-text-body); }")` → unreadable is about the
    `@return` body (line 2), NOT about the parameter list.
    `has @function parameter list note: false`. `varUses` returns only
    `--color-text-body@4`. The parameter list is correctly ignored because
    `px-to-rem` is a plain ident, not a `--custom-property` name, so
    `functionBinding` returns `nothing` without noting anything.
12. **Round-07/08 reproducer (hex-escape newline in string).** Appended
    `content: "x\41` + newline + `"` followed by `color: #ff0000` → the literal is
    found at line 23, the fail-safe fires on the unclosed block, verdict VIOLATION,
    exit 1. The structural fix from cycle 11 is intact.
13. **Round-08 controls (all five).** Each tested with a direct scanner probe:
    - Escaped quote (`content: "say \"hi\""`): file readable, literal found,
      VIOLATION, exit 1. The permissive change is still safe.
    - Comma escape (`content: "x\2c y"`): fail-safe fires (escape spelling comma
      routed to unreadable), literal still found, VIOLATION, exit 1.
    - Brace in string (`content: "}"`): file readable, brace doesn't close block,
      literal found, VIOLATION, exit 1.
    - Hex escape for quote (`content: "x\22 y"`): file readable, `\22` not decoded,
      string boundary correct, literal found, VIOLATION, exit 1.
    - URL fragment (`fill: url(#fade)`): url contents not scanned for literals,
      outside literal found, VIOLATION, exit 1.
14. **Registry digest.** `sha256(RULES.md)[:12]` = `b49ff596451f`, matches PDS.md
    frontmatter. Verified with Python hashlib.
15. **Byte budget (D-1).** Total: 65,360 / 65,536 bytes (176 bytes spare). PDS.md
    at 25,594 bytes exceeds its FINAL.md per-file limit of 20 KiB but is within the
    test's rebalanced threshold. The 64 KiB total is held. Doctrine byte-identical
    to round-08, so no change.
16. **D-2 count.** 9 rules with `enforced-by: check` or `both` that have no
    detector, all reporting `UNJUDGEABLE`. Confirmed at the CLI in the probe
    output. None silently passed. Count unchanged from round-08.
17. **No bundled registry.** No `RULES.md` exists under
    `addons/parley-design-check/`. Checker refuses with exit 3 when registry is
    absent (`--registry /nonexistent/RULES.md`).
18. **No placeholders.** Searched all 4 doctrine files for TODO/FIXME/PLACEHOLDER/
    lorem ipsum — none found.
19. **@function test.** `node --test --test-name-pattern="@function"` → 1 test
    passing, 0 failing. The test covers: the reproducer (must flag `--ghost`
    twice), the bound-parameter clean control, the declared-default clean control,
    comma splitting with strings/nested functions/brackets, the `returns` clause,
    escaped parameter names, escaped defaults, the unreadable channel for
    non-custom-property names, and the Sass false-alarm suppression.
20. **WCAG 2.2 fidelity.** All blocking numbers in WEB-ANNEX.md verified: 4.5:1
    (1.4.3), 3:1 (1.4.3/1.4.11), 24px/18.66px bold (1.4.3), 24×24 (2.5.8),
    320×256 (1.4.10). APCA correctly advisory.
21. **Doctrine lens (tier correctness).** All 30 rules (19 core + 11 web) checked:
    every declared tier is the lowest at which the rule can be decided. No rule is
    unfalsifiable, mis-classified, or on the wrong evidence tier. Every rule has a
    remedy and a counterexample (the WEB-ANNEX rules describe concrete failure
    scenarios in prose rather than with a labelled "Counterexample:" — the same
    observation I noted at round-08, which is acceptable).

## Findings

### [MAJOR] Opaque `var()` text in string contents is counted as a real token use (codex-1, independently reproduced)

**What is wrong.** `declarationVarUses` and `functionBinding` both run
`VAR_REFERENCE` directly over declaration values and decoded parameter defaults
without first masking opaque string and URL tokens. As a result, literal text
such as `"var(--ghost)"` inside a string is reported as a real custom-property
reference. I confirmed this with direct scanner probes:

- `@function --quoted(--x: "var(--ghost)") { result: var(--x); }` → `varUses`
  returns `--ghost@1` (false positive from the default).
- `.a { content: "var(--nope)"; color: var(--x); }` → `varUses` returns
  `--nope@1, --x@1` (false positive from an ordinary declaration).
- `.a { font-family: "var(--nope)"; color: var(--x); }` → same.

At the CLI against a copied sound L3 fixture, the `@function` quoted-default
case produces `core:token-used-undeclared`, VIOLATION, exit 1 — a false gate
result on valid CSS.

**Why it matters.** This is the false-alarm direction: valid CSS produces a
VIOLATION and exit 1. `maskOpaqueTokens` already exists and is used by the
`literal-outside-tokens` detector (which correctly masks string contents before
matching colour literals), but neither `declarationVarUses` nor `functionBinding`
applies it before running `VAR_REFERENCE`.

**Scope.** This is pre-existing, not introduced by cycle 12. I verified that
`declarationVarUses` in `82bde7d` (round-08 commit) also ran `VAR_REFERENCE`
directly on `declaration.value` without masking — the ordinary-declaration case
has been a false positive since at least the cycle-9 block-model rewrite. Cycle
12's new `functionBinding` inherited the same pattern for the @function default
path. codex-1 rated this MAJOR; I agree with the severity but note that only the
@function default path is new — the ordinary-declaration path is pre-existing.

**Concrete fix.** Apply `maskOpaqueTokens` to the declaration value (and to the
decoded parameter default suffix) before running `VAR_REFERENCE`, the same way
`literal-outside-tokens` already does. Add a negative fixture: a declaration
with `"var(--ghost)"` in a string value must not produce a `token-used-undeclared`
finding; and a clean control where a real `var(--color-text-body)` beside the
string is still found.

### [NIT] PDS §2 rule 1 is grammatically broken (carried from round-08)

PDS.md line 83–84: "A file declaring the spec under a kind §2 does not define MUST
be reported as violating this rule" — the sentence is missing a relative pronoun.
The intended meaning is "a kind that §2 does not define", but the "that" is absent,
so the sentence parses as gibberish on first read. The checker implements the
intended meaning correctly (`pds-check:l1-artifact-kind` fires on a kind §2 does
not define), so this is a prose defect only. The doctrine was not changed in cycle
12, so this carries forward unchanged from round-08.

**Fix:** Insert "that" before "§2 does not define": "A file declaring the spec
under a kind that §2 does not define MUST be reported as violating this rule."

### [NIT] Counts in prose disagree (carried from kimi-1 round-08)

`css.js:11` says "Seven constructs of one family have now been found by probe"
while `css.js:33` and `SKILL.md:88` say "Eight constructs broke it one at a time".
The eighth is the string-consumer desync described three comment blocks below the
stale seven. "Never write a count in prose" is one of the nine defects FINAL.md
says not to reproduce — counts in prose drift, and these two now have. Not
introduced by cycle 12.

**Fix:** Align the number, or drop it from both and point at the enumeration.

### [MINOR] §4 rule 4's second critique round has no §1 home (carried from kimi-1 round-08)

PDS §4 rule 4 permits a second critique round on an explicit Decider instruction,
but §1's mapping names only `round-02` for CRITIQUE, and the checker's
`PROCESS_HOMES` accepts only that home — so a run that exercises the permission
the spec grants fails the L2 process-order check the same spec defines. Doctrine
byte-identical since `f1c123d`; the engine's home list unchanged. Not introduced
by cycle 12.

**Fix:** §1 names the second round's home (and `PROCESS_HOMES` accepts it when the
instruction is recorded), or §4 rule 4 states the second round also files under
`round-02`; one fixture for the permitted shape.

## Open questions

None. The D-1 deviation (PDS.md at 25 KiB vs FINAL.md's 20 KiB, total held at
64 KiB) is the only place the implementation chose differently from the written
spec, and I take the position that it is justified: C3 adopted "64 KiB total" as
the binding ceiling, and the per-file numbers were one participant's proposed
split, not a binding decision. PDS.md at 25,594 bytes with 6 bytes of headroom
under the test's threshold is tight but sound, and the alternative (dropping
artifact entries or breaking the identical four-part shape) would damage the thing
the spec exists to be.

The convergence pattern is real: findings per review round have gone 10, 3, 3, 2,
2, 1, 1, and the last two findings were regressions introduced by the fix before
them, in opposite directions. Cycle 12's fix did not introduce a new regression in
the @function parameter-list path — the parameter list is parsed at top-level
commas through the shared string and escape consumers, only a segment's leading
custom-property identifier is a formal, `var()` references in defaults are
collected as real uses, and a segment whose leading token is not a custom-property
name goes to the unreadable channel rather than being guessed in either direction.

The MAJOR I confirm (codex-1's quoted-`var()` finding) is pre-existing in the
ordinary-declaration path and inherited by the new @function default path. It is
the false-alarm direction — valid CSS produces a VIOLATION and exit 1 — and it
has a concrete fix (apply `maskOpaqueTokens` before `VAR_REFERENCE`). It should be
fixed before a v1.0.0 tag, but it is not a regression from cycle 12 and it is not
a false clean. The round-08 CRITICAL (false clean) is genuinely closed, and the
fix did not swing back toward a false alarm in the @function parameter path
itself.
