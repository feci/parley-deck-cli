---
agent: hermes-1
idea: parley-design-skills
review-round: 04
date: 2026-07-28
reviewed-commit: e3ca916
---

## Summary

BLOCK. The three round-03 findings I filed are genuinely closed: per-direction
token sidecars are enforced, the frontmatter lexer rejects trailing `#` comments,
and same-tier alias edges fail strict descent. The unterminated-string hole is
closed for the browser-applied case. I verified each with a reproducer.

The new CSS fail-safe is evadable. I independently reproduced the escaped-`url`
bypass codex-1 found and extended it to four escape variants, all of which hide
a raw `color: #ff0000` from the literal detector, keep the fail-safe quiet, and
roll the run to a false PASS at exit 0. This is a sixth member of the
scanner-evasion family the fail-safe was built to close. The family is not closed.

D-1 (per-file byte budget rebalance) is acceptable — the 64 KiB total is held.
D-2 (undetected `enforced-by: check` rules) is an honest declaration, not a
silent pass. The doctrine (RULES.md, WEB-ANNEX.md) is sound: WCAG 2.2 thresholds
are correct, rule classifications are right, no rule is unfalsifiable or
mis-classified, and evidence tiers are correctly assigned.

## What I verified (commands run, and their result)

- `npm test` in the skill repo: 227 passed, 0 failed.
- `node addons/parley-design-check/bin/check.js --help`: exit 0, documents exit 3 (no registry) and exit 4 (UNJUDGEABLE).
- Round-03 MAJOR-1 (same-tier alias edges): ran the shipped `alias-same-tier` fixture with `--level L3`. Exit 1, `pds-check:l3-alias-direction` fires for both the semantic-to-semantic and component-to-component edges. CLOSED.
- Round-03 MAJOR-2 (frontmatter trailing `#`): wrote a DIRECTION with `assigned: flat # comment`. Exit 1, `pds-check:l1-frontmatter-parses` rejects it: "an unquoted `#` in `flat # this should be rejected`". CLOSED.
- Round-03 MAJOR-3 (per-direction token sidecar): wrote two DIRECTIONs naming the same `shared.tokens.json`. Exit 1, `pds-check:l2-process-order` fires: "the DIRECTIONs filed as 'claude-1' and 'codex-1' resolve their token files to one path". Also verified the `sound-run` fixture now ships `codex-1.tokens.json` and `claude-1.tokens.json` separately. CLOSED.
- Round-03 MINOR (gate error-string separator drift): ran the `collapsed-run` fixture. G1 finding reads `G1 DISTINCTNESS:` (colon); PDS §3 template reads `G1 DISTINCTNESS —` (em-dash). Still present. NOT CLOSED (MINOR, cosmetic).
- Fail-safe attack suite (55 probes under `/tmp/pds-*`): content-quoted braces, unterminated strings, braces in unquoted `url()`, comment delimiters in `url()`, escaped braces in selectors, CDO/CDC, at-rules (`@layer`, `@container`, `@scope`, `@starting-style`, `@font-face`, `@import`), CSS nesting, CRLF, BOM, null bytes, all control characters 0x00–0x1F + 0x7F, backslash-newline, unbalanced parens. Results: all either correctly detected, correctly UNJUDGEABLE, or correctly not browser-applied (verified by tracing CSS Syntax §4.3.3, §4.3.5, §4.3.6, §4.3.7, §4.3.12).
- Fail-safe evasion (CRITICAL): `u\72l(a/*) ; color: #ff0000; */b)` — exit 0, verdict PASS, 0 literal violations, 0 unreadable. Verified with `scanStylesheet` directly: `stripComments` blanks from `/*` to `*/`, erasing the `color` declaration. The scanner sees only `background: u\72l(a ... b)`. Extended to four escape forms (`u\72l`, `\75 rl`, `u\000072l`, `\75\72\6cl`): all produce false PASS.
- Doctrine file sizes: SKILL.md 6,519 / PDS.md 25,594 / RULES.md 23,225 / WEB-ANNEX.md 10,022 = 65,360 total. Under the 65,536-byte ceiling. Per-file thresholds in `design-addons.test.js`: 7/25/24/11 KiB.
- Checker capability output: 18 rules with detectors, 12 without. All 9 `enforced-by: check`/`both` rules without a detector report `UNJUDGEABLE` (5 as "no detector implements this rule", 4 as "T2/T3 evidence is above this checker"). No silent passes.
- Registry refusal: `--registry /nonexistent` → exit 3. Structural checks still run.
- Detector spot-checks: `web:gradient-text`, `web:edge-stripe`, `core:focus-indication`, `web:motion-defaults`, `core:face-outside-allowlist` — all fire on their target constructs.
- WCAG 2.2 thresholds in WEB-ANNEX.md: SC 1.4.3 (4.5:1/3:1), SC 1.4.11 (3:1), SC 2.5.8 (24×24 CSS px), SC 1.4.10 (320×256 CSS px), SC 2.4.7/2.4.11, SC 2.2.2/2.3.3 — all correct. Large-text definition (24 CSS px, or 18.66 at bold) matches WCAG 2.2.
- Rule classification audit: 19 core + 11 web = 30 rules. No rule is unfalsifiable (each has a concrete counterexample and remedy). No rule is mis-classified (quality/slop/system assignments match their authority model). Evidence tiers are correct (T0 for token-decidable rules, T1 for source-decidable, T2 for computed, T3 for raster).

## Findings

### [CRITICAL] Escaped `url` identifier bypasses the fail-safe and produces a false PASS

**What.** `scanComments` recognises `url(` only via the literal regex `/^url\($/i` (css.js:92). CSS identifier escapes (§4.3.7) are decoded by the browser before token classification, so `u\72l(` is the `url(` token to a browser but not to this scanner. When `/*` appears inside such a url token, `scanComments` treats it as a comment opener and blanks through the next `*/`, erasing every declaration between them. The stack machine then sees balanced braces and parens, so the fail-safe stays quiet.

**Reproducer.**

```css
.probe {
  background: u\72l(a/*) ; color: #ff0000; */b);
  padding: var(--primitive-ink);
}
```

`check --json --registry RULES.md --contract CONTRACT.md <dir>` returns `{"verdict":"PASS"}`, exit 0, 0 literal violations, 0 unreadable entries. In Chromium, `getComputedStyle` returns `rgb(255, 0, 0)` — the browser applies the colour.

I extended the test to four escape variants (`u\72l`, `\75 rl`, `u\000072l`, `\75\72\6cl`); all produce the same false PASS.

**Why it matters.** The fail-safe is the safety argument for the entire hand-written scanner: the comment at css.js:11–22 promises that any tokenisation construct the scanner cannot confidently read will degrade every stylesheet detector to UNJUDGEABLE and prevent exit 0. This construct evades that guarantee — a browser-applied raw literal ships while the checker certifies clean. This is the same family as the original `content: "}"` hole (round-01), reached via a sixth vector. The family is not closed.

**Concrete fix.** In `scanComments`, before checking for `/*`, decode CSS identifier escapes (§4.3.7: `\` followed by 1–6 hex digits, optionally followed by a single whitespace) and check whether the decoded ident is `url`. If an escape sequence appears at a position where a url check would run and the decoded form is `url`, treat the following `(` as opening a url token. Alternatively, fail closed: if `scanComments` encounters a backslash escape at a position where `url` could begin, emit an unreadable note and let the fail-safe handle it. Add the four reproducers above as end-to-end regression tests requiring either a literal finding or an unreadable/UNJUDGEABLE/exit-4 result.

### [MINOR] Gate error-string separator drift (carried from round-01/03)

**What.** The checker emits `G1 DISTINCTNESS:` (colon) but PDS §3 writes `G1 DISTINCTNESS —` (em-dash). Same for G2. This has persisted across three rounds.

**Why it matters.** A downstream consumer that pattern-matches the PDS §3 template to parse gate findings will miss them. The separator is not load-bearing for conformance — the gate fires correctly — but the strings diverge from the spec.

**Concrete fix.** Replace the colon with the em-dash (`—`) in the G1/G2 finding strings to match PDS §3, or update PDS §3 to use the colon. Either way, pick one and align.

## Open questions

- The `scanComments` url check is the only place in the scanner where an escaped ident could cause a structural misread. Are there other positions where the scanner's literal-string checks (`startsWith("/*")`, `indexOf("*/")`, quote matching) could be fooled by decoded-escape forms? A systematic audit of every literal-string comparison in `css.js` against the escape-decoding rules in §4.3.7 would close the question.
- The control-character probes (0x00–0x1F except tab/LF/CR/FF) showed that non-whitespace control characters prepend to the property name, causing the detector to miss the declaration. This is not a fail-safe evasion because the browser also treats the declaration as invalid (the control char is not a valid ident-start). But should the scanner report these as "text inside a rule that is not a declaration" (fail-safe condition f) to be more conservative? Currently it silently parses them as declarations with garbled prop names.
