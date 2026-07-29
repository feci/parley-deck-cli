---
agent: hermes-1
idea: parley-design-skills
review-round: 10
date: 2026-07-28
reviewed-commit: 1675b6f
---

## Summary

ACCEPT-WITH-RESERVATIONS. Cycle 13 closes the round-09 MAJOR (codex-1 and
hermes-1, independently): opaque token contents — a `var()` inside a string or
a url token — are no longer counted as a real custom-property reference. Both
collectors (`declarationVarUses` and `functionBinding`) now route through one
shared `valueVarUses`, which runs `VAR_REFERENCE` over the output of
`maskOpaqueTokens` — the same masker `literal-outside-tokens` already uses. The
fix did not swing toward a false clean: the round-08 reproducer (real `var()`
in a default) still flags `--ghost` and exits 1, and escaped `var()`
(`\76 ar(--ghost)`) is still found because decoding runs before masking.

The deliberate non-masking of the markup raw-line sweep is measured rather than
argued, and I accept the measurement for genuine host-language regions: outside
`<style>` and `style=`, a quote belongs to the host language, not to a CSS
string token. Routing that path through `valueVarUses` would lose 2,410
references across 306 markup files.

One new MAJOR: the raw sweep does not blank `style=` attribute spans before
running `VAR_REFERENCE`, so a `var()` inside a CSS string in a `style=`
attribute is found as a false reference — the false-alarm direction. The parsed
path (step 2) correctly masks it, but the raw sweep (step 3) re-finds it. This
is pre-existing (since cycle 10, not a cycle-13 regression) and fixable with a
one-line change. codex-1 filed the same finding independently.

Three carried NITs/MINOR from prior rounds remain, none introduced by cycle 13,
none blocking.

## What I verified (commands run, and their result)

1. `npm test` — 246 passing, 0 failing, 0 skipped.
2. `node addons/parley-design-check/bin/check.js --help` — exit 0, usage and
   exit codes as documented.
3. **Git state.** Branch `parley-design-skills`, HEAD `1675b6f`, worktree clean.
   `git diff 9121ec2..1675b6f --stat`: only `lib/css.js` (58+, 8-) and
   `test/checker.test.js` (114+). Doctrine byte-identical between the commits.
4. **Round-09 reproducer (scanner-level, 9 probes).** Ran `varUses()` and
   `markupVarUses()` directly:
   - Quoted `@function` default (`"var(--ghost)"`): `(none)`. OK.
   - URL-quoted `@function` default (`url("var(--ghost)")`): `(none)`. OK.
   - String in ordinary declaration (`content: "var(--nope)"`): `--x` only. OK.
   - Font-family string (`'var(--nope)'`): `--x` only. OK.
   - URL token locator (`url(var(--ghost))`): `--x` only. OK.
   - Real default beside string (`"var(--ghost)", --y: var(--gap)`): `--gap`. OK.
   - Escaped var in default (`\76 ar(--ghost)`): `--ghost` — decoded before
     masking. OK.
   - Real var beside string in declaration (`"var(--nope)" var(--real)`):
     `--real`. OK.
   - URL with real var after (`url("/img/x.png") no-repeat, var(--fade)`):
     `--fade`. OK.
   All 9 pass. The fix masks opaque tokens without losing real references.
5. **Round-09 reproducer (end-to-end CLI, `--level L3`).** 6 probes against
   sound-run fixture:
   - Quoted `@function` default: PASS, verified L3, exit 0. OK.
   - URL-quoted `@function` default: PASS, verified L3, exit 0. OK.
   - String in ordinary declaration: PASS, verified L3, exit 0. OK.
   - URL locator text: PASS, verified L3, exit 0. OK.
   - Round-08 real `var()` in default: VIOLATION, exit 1, `--ghost` flagged. OK.
   - Escaped `var()` in default: VIOLATION, exit 1, `--ghost` flagged. OK.
   All 6 pass. The fix did not swing toward a false clean.
6. **Markup asymmetry assertion.** Two scanner probes:
   - `<style>` body with `"var(--nope)"` + `var(--real)`: `--real` only — opaque
     tokens masked through the declaration path. OK.
   - `const cl = "color:var(--error)"`: `--error` — outside `<style>` and
     comments, a quote belongs to the host language. OK.
   Both halves hold, locked by the test suite.
7. **Cycle-13 test by name.** `node --test --test-name-pattern="opaque token"`:
   1 test, 0 failing. Covers 5 reproducer cases, 5 controls, 2 markup asymmetry
   assertions, 4 end-to-end L3 probes, and the round-08 guard.
8. **Cycle-12 test (regression check).** `node --test --test-name-pattern="@function"`:
   2 tests, 0 failing. @function parameter-list fix intact.
9. **Registry digest.** `sha256(RULES.md)[:12]` = `b49ff596451f`, matches
   PDS.md frontmatter. Doctrine byte-identical to round-09.
10. **Byte budget (D-1).** Total: 65,360 / 65,536 (176 bytes spare). PDS.md at
    25,594 bytes exceeds its FINAL.md per-file limit of 20 KiB but is within the
    test's rebalanced threshold. 64 KiB total held. Doctrine byte-identical to
    round-09.
11. **D-2 count.** 9 rules with `enforced-by: check`/`both` and no detector, all
    reporting `UNJUDGEABLE`. None silently passed. Unchanged from round-09.
12. **No bundled registry.** No `RULES.md` under `addons/parley-design-check/`.
    Checker refuses with exit 3 when registry absent.
13. **No placeholders.** Searched all 4 doctrine files for TODO/FIXME/PLACEHOLDER/
    lorem ipsum — none found.
14. **WCAG 2.2 fidelity.** All blocking numbers verified: 4.5:1 (1.4.3), 3:1
    (1.4.3/1.4.11), 24px/18.66px bold (1.4.3), 24x24 (2.5.8), 320x256 (1.4.10).
    APCA correctly advisory, no APCA in any rule YAML.
15. **Doctrine lens (tier correctness).** All 30 rules (19 core + 11 web)
    checked: every rule has required YAML keys, no duplicate ids, severity 0-4
    only. T2/T3 rules (5 total) correctly UNJUDGEABLE. No rule unfalsifiable,
    mis-classified, or on the wrong evidence tier.
16. **NOTICE.md.** Prior-art attribution present: hallmark (MIT), impeccable
    (Apache-2.0), independent authorship stated, standards referenced (DTCG
    2025.10, WCAG 2.2, RFC 2119/8174).
17. **codex-1's round-10 `style=` double-scan finding (independently
    reproduced).** Scanner-level: `markupVarUses('<div style="font-family:
    \'var(--ghost)\'">hello</div>')` returns `--ghost@1` — a false reference.
    End-to-end CLI on sound-run fixture with `probe.html` containing the same
    markup: VIOLATION, exit 1, `core:token-used-undeclared` naming `--ghost`.
    The parsed path (lines 1444-1448) correctly masks the string (returns
    nothing); the raw sweep (line 1450) does not blank `style=` attribute spans,
    so `VAR_REFERENCE` re-finds `--ghost` in the attribute text. Verified
    pre-existing: same result at cycle 12 (`9121ec2`) and cycle 10 (`076ded5`),
    not a cycle-13 regression. The fix (blank `style=` spans before the raw
    sweep) is safe: the parsed path already retains true positives
    (`color: var(--real)` → `--real`), so blanking only removes the false
    positive from the raw sweep.

## Findings

### [MAJOR] `style=` attribute double-scanned: var() inside a CSS string in a style attribute is a false reference

**What.** `markupVarUses` parses `style=` attributes as CSS (lines 1444-1448),
correctly masking opaque tokens. But the raw-line sweep (line 1450) blanks only
`STYLE_ELEMENT` and `MARKUP_COMMENT` spans — it does not blank `style=`
attribute spans. So `VAR_REFERENCE` runs over the `style=` attribute text and
finds `var()` inside CSS strings that the parsed path correctly ignored.

Scanner probe:
```
markupVarUses('<div style="font-family: \'var(--ghost)\'">hello</div>')
  → ["--ghost@1"]   // false: --ghost is inside a CSS string
```

End-to-end: sound-run fixture with `probe.html` containing that markup yields
VIOLATION, exit 1, `core:token-used-undeclared` naming `--ghost`.

**Why it matters.** A `style=` attribute is CSS, not host-language markup. The
parsed path treats it as CSS and masks the string; the raw sweep treats the same
text as host language and does not. The commit message says the non-masking is
for text "outside `<style>` and `style=`" — but the code does not actually blank
`style=` spans in the raw sweep. This is the false-alarm direction (valid CSS
becomes a VIOLATION), not a false clean, and it is pre-existing (since cycle 10),
not a cycle-13 regression. But it is a gap in the asymmetry the commit claims to
have locked, and the same class as the round-09 finding the cycle existed to
close.

**Concrete fix.** Add a `STYLE_ATTRIBUTE` regex constant
(`/\bstyle\s*=\s*("[^"]*"|'[^']*')/gi`, the same pattern `styleAttributes` already
uses) and change line 1450 from:
```
blankSpans(blankSpans(text, STYLE_ELEMENT), MARKUP_COMMENT)
```
to:
```
blankSpans(blankSpans(blankSpans(text, STYLE_ELEMENT), MARKUP_COMMENT), STYLE_ATTRIBUTE)
```
The parsed path (lines 1444-1448) already retains true positives
(`color: var(--real)` → `--real`), so blanking only removes the false positive.
Add a fixture: `style="font-family: 'var(--ghost)'"` must PASS at L3 (exit 0),
and a control `style="color: var(--real)"` must still find `--real`.

### [NIT] PDS §2 rule 1 is grammatically broken (carried from round-08)

PDS.md line 83-84: "A file declaring the spec under a kind §2 does not define
MUST be reported as violating this rule" — the sentence is missing a relative
pronoun ("that"). The checker implements the intended meaning correctly
(`pds-check:l1-artifact-kind`), so this is a prose defect only. Doctrine
byte-identical since at least round-07.

Fix: Insert "that" before "§2 does not define".

### [NIT] Counts in prose disagree (carried from round-08)

`css.js:11` says "Seven constructs of one family have now been found by probe"
while `css.js:33` says "Eight fix-up cycles closed members of one family one at
a time". Both are prose counts in a file whose own doctrine says "never write a
count in prose" — and both are now stale as further cycles have run. Not
introduced by cycle 13.

Fix: Drop both counts and point at the token-type enumeration table.

### [MINOR] §4 rule 4's second critique round has no §1 home (carried from kimi-1 round-08)

PDS §4 rule 4 permits a second critique round on an explicit Decider
instruction, but §1's mapping names only `round-02` for CRITIQUE, and the
checker's `PROCESS_HOMES` accepts only that home. A run that exercises the
permission the spec grants fails the L2 process-order check the same spec
defines. Doctrine byte-identical since `f1c123d`. Not introduced by cycle 13.

Fix: §1 names the second round's home (or §4 rule 4 states it also files under
`round-02`), and `PROCESS_HOMES` accepts it when the instruction is recorded.

## Open questions

None. The D-1 deviation (PDS.md at 25 KiB vs FINAL.md's 20 KiB, total held at
64 KiB) remains the only place the implementation chose differently from the
written spec, and I take the same position as round-09: it is justified. C3
adopted "64 KiB total" as the binding ceiling; the per-file numbers were one
participant's proposed split. PDS.md at 25,594 bytes with 176 bytes of headroom
under the 64 KiB total is tight but sound, and the alternative (dropping
artifact entries or breaking the identical four-part shape) would damage the
thing the spec exists to be.

The deliberate non-masking of the markup raw-line sweep is the right call for
genuine host-language regions. The measurement is in the commit and in the
test: routing that path through `valueVarUses` would lose 2,410 references
across 306 of 2,236 markup files, because outside `<style>` and `style=` a
quote belongs to the host language and not to a CSS string token. The markup
spans that are CSS (`<style>` bodies and `style` attributes) are parsed as CSS
and collected through `valueVarUses` with every other declaration, so the
opaque-token rule reaches them there.

The MAJOR finding above is about the gap between that claim and the code: the
`style=` attribute IS parsed as CSS in step 2, but is NOT blanked in the raw
sweep in step 3, so the raw sweep re-finds `var()` inside CSS strings that the
parsed path correctly masked. The fix is one line and a regex constant, and it
is safe because the parsed path already retains true positives. This should
land before a v1.0.0 tag but does not block merge — it is the false-alarm
direction, pre-existing, and narrow.
