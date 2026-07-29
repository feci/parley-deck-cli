---
agent: codex-1
idea: parley-design-skills
review-round: 08
date: 2026-07-28
reviewed-commit: 82bde7d
---

## Summary

❌ BLOCK. The round-07 CRITICAL is genuinely closed: my exact reproducer now exposes the live raw colour, the CLI exits 1, and the six CSS passes share one string/escape implementation. I also tested the stronger token-type × consumer claim with strings containing hex escapes, comment-shaped text, declaration delimiters, parentheses, and colour-shaped text; the consumers remained aligned.

The requested reverse-direction review did, however, find one new false clean in the permissive `@function` closure. The checker mistakes custom-property names in parameter default expressions for formal parameters, then suppresses real `var()` uses of those names in the function body. Chromium resolves the omitted token while the checker reports PASS, verifies L3, records no unreadable input, and exits 0. That is a release-blocking capability error, but it is localized and has a concrete repair.

Apart from that finding, the previously reported case-insensitive function/unit handling, URL fragment exclusion, escaped-string handling, modern at-rule classification, CDO/CDC handling, registry refusal, waiver behavior, generated capability, packaging, and exit-code paths held under direct checks. D-1 is acceptable: the four doctrine files total 65,360 bytes against the 65,536-byte aggregate cap, preserve the required four-part artifact, and remain within the binding early-warning thresholds.

## What I verified (commands run, and their result)

- `git status --short`, `git rev-parse --abbrev-ref HEAD`, and `git rev-parse HEAD` in the implementation repository: clean `parley-design-skills` checkout at `82bde7d78000fbef9848636fec0d37b65d885014`.
- `npm test`: 244 tests passed, 0 failed.
- `node addons/parley-design-check/bin/check.js --help`: documented command surface and exit-code contract present.
- CLI against a directory without `parley-design-tokens.md`: `status: absent`, `verdict: UNJUDGEABLE`, exit 3, with an actionable refusal message.
- CLI against an unchanged sound L3 fixture: PASS, verified L3, exit 0; the generated capability contained all 18 rule/detector records and no unreadable files.
- `npm_config_cache=<temporary cache> npm pack --dry-run`: succeeded with 153 files; the four doctrine files, checker implementation, tests, and `NOTICE.md` are packaged, with no bundled substitute registry or placeholder artifact.
- Byte counts and digest checks: `SKILL.md` 6,519 bytes, `PDS.md` 25,594, `RULES.md` 23,225, and `WEB-ANNEX.md` 10,022, totaling 65,360/65,536; the registry digest prefix is `b49ff596451f` and matches the doctrine.
- My round-07 reproducer, unchanged:

  ```css
  .probe {
    content: "x\41
  "; /* x: 1; } */
    color: #ff0000;
    dummy: y { z: 1;
  }
  ```

  Direct scanning now reads `content: "xA"` and `color: #ff0000`; the CLI reports `core:literal-outside-token-layer`, returns VIOLATION, and exits 1. The malformed trailing block is also reported unreadable. Kimi's V5 spelling likewise preserves `"xA/*"`, reads the colour, reports no unreadable input, and exits 1 on the literal.
- Definition/call-site inspection of `lib/css.js`: exactly one `consumeEscape` and one `stringToken` definition exist. `scanComments`, `scanStylesheet`, `decodeDeclarationText`/`decodeStringToken`, `maskOpaqueTokens`, `splitDeclaration`, and `parenBalance` route string and escape handling through them.
- A direct consumer matrix put the same hex-escaped newline string, containing `/*;:{}(*/`, through comment stripping, stylesheet scanning, declaration decoding, colon splitting, opaque-token masking, and parenthesis balancing. The real comment after the string was removed, the comment-shaped text inside it was preserved, the following declaration was read, the colon inside it did not split the declaration, in-string `#fade` was masked, and outside `#beef` remained visible. Backslash-newline, escaped quote/backslash, ordinary comment, and hex-escape-plus-space controls also held.
- Reverse-direction CLI fixtures covered uppercase `RGB()`/`VAR()`, uppercase `PX`, escaped strings followed by raw colours, `url(#fade)` followed by a raw colour, bound and unbound `@function` variables, `@try`, `@position-fallback`, and CDO/CDC. The violating fixture reported every intended finding without reporting URL/string contents; the declared-token and bound-parameter control passed L3 with no unreadable input.
- Removing the required waiver from the registry produced `pds-check:l2-process-order`, VIOLATION, exit 1. An explicitly empty waiver field with no waiver file passed, as required.
- Chromium 150, through an isolated ego-browser task space, parsed the following as a real `CSSFunctionRule`, returned true from `CSS.supports("color", "--pick()")`, and computed the probe colour as `rgb(255, 0, 0)`:

  ```css
  :root { --ghost: rgb(255,0,0); }
  @function --pick(--x: var(--ghost)) { result: var(--ghost); }
  .probe { color: --pick(); }
  ```

  The same CSS added to a copied sound L3 fixture produced no `varUses`, no finding, no unreadable input, PASS, verified L3, and exit 0 from the checker.

## Findings

### [MAJOR] `@function` defaults turn real token uses into a false clean

**What is wrong:** `functionParameters()` gathers every `--name` anywhere after the opening parenthesis of an `@function` prelude. For `@function --pick(--x: var(--ghost))`, it therefore returns both `--x` and `--ghost`, although only `--x` is a formal parameter. `declarationVarUses()` then suppresses `var(--ghost)` in the body as if it were a bound parameter. It also does not collect the `var(--ghost)` reference in the default expression itself.

**Why it matters:** `--ghost` is a browser-resolved token reference. In the reproduced fixture it is absent from the ratified token document, yet the checker certifies PASS at L3 with exit 0. This is precisely the false-clean direction that the round requested reviewers to test after making the checker more permissive, and it overstates the delivered `core:token-used-undeclared` capability.

**Concrete fix:** parse the parameter list at top-level commas while respecting strings and nested functions; take only each segment's leading custom-property identifier as a bound formal; collect `var()` references from each default expression as real uses; and suppress only true formal parameters while scanning the function body. Add the browser-backed reproduction above as a failing checker fixture, alongside a clean control in which `result: var(--x)` uses only the bound parameter and a clean default that references a token declared in the registry.

## Open questions

None. The blocker is reproduced, browser-confirmed, and has a bounded repair.
