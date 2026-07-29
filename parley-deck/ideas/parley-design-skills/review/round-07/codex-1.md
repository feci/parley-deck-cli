---
agent: codex-1
idea: parley-design-skills
review-round: 07
date: 2026-07-28
reviewed-commit: 076ded5
---

## Summary

❌ BLOCK. The round-06 `#url(`/`@url(` bypasses and the `<style>` selector `var()` false positive are genuinely closed. The two declared remaining OPENs are also accurately described as over-reads, and reverting the proposed markup narrowing was the right call.

The class-level conclusion is nevertheless false. The §4 family list is accounted for in the prose, but its string-token verdict is wrong: `scanComments()` and `scanStylesheet()` do not use the same string/escape consumer. A hexadecimal escape inside a string can consume the newline that follows it in Chromium and the main scanner, while the comment pre-pass treats that newline as the end of the string. A later real comment is consequently left in the source and reinterpreted as declarations and blocks. The exact probe below makes Chromium apply `#ff0000` while the checker returns `PASS`, verifies L3, records `unreadable: []`, and exits 0. This refutes both the claimed zero silent holes and the claimed closure of the token class.

I accept deviation D-1. The binding total remains 65,360/65,536 bytes, all four early-warning thresholds hold, and the identical artifact shape in PDS.md is more valuable than preserving the superseded initial per-file split. I found no independent defect in registry parsing, waiver validation, capability generation, packaging, or the doctrine files.

## What I verified (commands run, and their result)

- Confirmed the implementation checkout before testing:

  ```text
  $ git branch --show-current
  parley-design-skills
  $ git rev-parse HEAD
  076ded5b44c38a428fd7e6b7616f94b5bf745b25
  $ git status --short
  [no output]
  ```

- Ran the complete test suite:

  ```text
  $ npm test
  tests 237; pass 237; fail 0
  ```

- Re-ran my round-06 `#url(` and `@url(` fixtures beside the sound L3 run. For each fixture, Chromium still computed `rgb(255, 0, 0)`, while the checker now found `core:literal-outside-token-layer`, recorded the intentionally malformed remainder as unreadable, refused L3 verification, and exited 1:

  ```text
  #url(: verdict VIOLATION; verified null; process exit 1
  @url(: verdict VIOLATION; verified null; process exit 1
  ```

  I also re-ran the exact `<style>` selector case. It returned `PASS`, verified L3, had no unreadable input, reported no `--ghost` use, and exited 0. That is the correct result because Chromium keeps `var(--ghost)` in selector text and resolves no reference from it.

- Checked the [CSS Syntax Level 3 §4 output-token list](https://www.w3.org/TR/css-syntax-3/#tokenization) against the scanner's cycle-10 grouping. The prose accounts for the token families, including the eight named token classes that can contain significant code points and the bracket tokens through the matched-simple-block model. I then checked the consumers, rather than trying another hash/URL spelling. This exposed the inconsistent string handling at `lib/css.js:395-405`: `scanComments()` skips only one code point after `\`, while the shared `stringToken()` used later calls `consumeEscape()` and consumes up to six hexadecimal digits plus the following whitespace.

- Ran this scanner probe:

  ```css
  .probe {
    content: "x\41
  "; /* x: 1; } */
    color: #ff0000;
    dummy: y { z: 1;
  }
  ```

  Chromium's CSSOM contained only:

  ```text
  .probe { content: "xA"; color: rgb(255, 0, 0); }
  computed color: rgb(255, 0, 0)
  ```

  The direct scanner instead left the real comment unstripped, omitted `color`, manufactured a declaration `/* x: 1` and a phantom `dummy: y` block, and returned:

  ```json
  {
    "blocks": [
      {"selector": ".probe", "declarations": [
        {"prop": "content", "value": "\"xA\""},
        {"prop": "/* x", "value": "1"}
      ]},
      {"selector": "dummy: y", "declarations": [{"prop": "z", "value": "1"}]}
    ],
    "unreadable": []
  }
  ```

  I copied the sound-run fixture under `/tmp`, added that stylesheet, and ran the real CLI:

  ```text
  $ node addons/parley-design-check/bin/check.js \
      --registry addons/parley-design/references/RULES.md \
      --level L3 --json /tmp/.../string-run
  verdict PASS
  report exit 0
  verified L3
  unreadable []
  probe findings []
  process exit 0
  ```

- Reproduced both declared residual OPENs in Chromium and the scanner:

  - `border-radius: 1\65 5` is decoded by the scanner to `1e5`; Chromium discards it and computes `0px`.
  - `color: #ff0000 garbage` is retained by the scanner; Chromium drops that declaration and computes black while retaining the following valid declaration.

  Both are over-reads at the declared T1 source boundary. Neither hides a browser-applied declaration, and neither is worse than stated.

- Read the markup remedy and its retained controls (`cn("h-[var(--radix-select-trigger-height)]")`, JSX/SVG `stopColor="hsl(var(--primary))"`, and JavaScript token strings). Restricting discovery to style attributes and the one supported utility form would erase legitimate reference classes. The measured loss of 1,799 references across 203 files is therefore a sound reason to revert that exact remedy; the targeted `<style>` excision is the appropriate fix.

- Exercised the checker contract outside the scanner:

  ```text
  $ node addons/parley-design-check/bin/check.js --help
  [documents exits 0/1/2/3/4 and L1-L4]

  $ node addons/parley-design-check/bin/check.js \
      --registry /tmp/definitely-missing-parley-r07.md \
      --json addons/parley-design-check/test/fixtures/literal-outside-tokens/fail
  registry status absent; verdict UNJUDGEABLE; findings []; exit 3
  stderr: rule checks were refused ... this checker carries no copy

  $ node addons/parley-design-check/bin/check.js \
      --registry addons/parley-design/references/RULES.md \
      --level L3 --json addons/parley-design-check/test/fixtures/conformance/sound-run
  PASS; verified L3; exit 0; 18 generated detectors; 18 rules with detectors
  registry core-rules/1.0.0, digest b49ff596451f; unreadable []
  ```

  I also read `registry.js`, the gating and waiver paths in `engine.js`, all detector modules, the checker tests, the installer tests, and the four doctrine files. The suite's refusal, malformed-registry, self/ghost-signed waiver, widening, expiry, conformance, recusal, and exit-code cases all passed.

- Verified distribution and doctrine constraints:

  ```text
  $ wc -c addons/parley-design/SKILL.md \
      addons/parley-design/references/{PDS.md,RULES.md,WEB-ANNEX.md}
   6519 SKILL.md
  25594 PDS.md
  23225 RULES.md
  10022 WEB-ANNEX.md
  65360 total

  $ npm pack --dry-run --json
  153 files; all four parley-design doctrine files; 126 parley-design-check files
  ```

  The registry digest recomputed to `b49ff596451f`, matching PDS.md and the generated report. The no-placeholder, built-ins-only, single-registry, known-rule-ID, and installer-discovery/package tests passed.

## Findings

### [CRITICAL] The comment pre-pass mis-tokenises hexadecimal escapes inside strings and can forge a clean L3 certificate

**What is wrong:** `scanComments()` claims to read strings with tokenizer ordering, but its loop at `lib/css.js:395-405` advances only one code point after a backslash. CSS Syntax §4.3.7 consumes up to six hexadecimal digits and then one optional whitespace code point. When that whitespace is a newline, `scanComments()` ends the string before the browser and `scanStylesheet()` do. It then decides whether `/*` is a comment from the wrong token stream. The published probe demonstrates the resulting block/declaration desynchronisation.

**Why it matters:** this is a silent under-read of browser-applied CSS, not one of the accepted T1 over-reads. It suppresses the real raw-colour declaration, records no unreadable input, verifies L3, returns `PASS`, and exits 0. A CI gate can therefore issue the exact false-clean certificate the scanner's fail-safe is intended to prevent. It also proves that the cycle-10 “0 silent holes” differential and the string-token verdict did not cover the implementation's separate comment pass.

**Concrete fix:** make `scanComments()` call the same `stringToken()`/`consumeEscape()` path as `scanStylesheet()`; preferably factor the token walk so both passes cannot diverge again. Add the exact string→comment→raw-colour probe above as direct-scanner, CLI L3, and Chromium-backed regressions, with ordinary strings, bad strings, escaped newlines, and hexadecimal escapes whose optional whitespace is and is not a newline as controls. The negative case must either find the applied literal or mark the file unreadable and exit 4. Then rerun and retain the complete §4 differential in reviewable form before claiming the class closed.

### [MINOR] The promised per-token verdict table is absent from the review record

**What is wrong:** IMPLEMENTATION.md says the cycle-10 per-token-type table “lives in the cycle-10 agent report,” but no such report or table is present in the idea tree. The only retained enumeration is grouped prose in `css.js` and IMPLEMENTATION.md.

**Why it matters:** round 7 explicitly depends on reviewing per-type verdicts. Without the table, a reviewer cannot inspect the tested construct, scanner path, browser outcome, and open/closed decision for each §4 token type. The incorrect string verdict shows that prose-level coverage is not equivalent to checking every consuming pass.

**Concrete fix:** retain the actual verdict table with the cycle record, mapping every §4 output token to all relevant scanner passes, browser/control fixtures, measured result, and OPEN/CLOSED disposition. Store the differential construct matrix or harness beside it so “0 silent holes” is reproducible.

## Open questions

None. The block is reproducible and the repair boundary is concrete.
