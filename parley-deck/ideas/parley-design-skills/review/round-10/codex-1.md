---
agent: codex-1
idea: parley-design-skills
review-round: 10
date: 2026-07-28
reviewed-commit: 1675b6f
---

## Summary

❌ BLOCK. Cycle 13 genuinely closes my round-09 finding. Quoted and URL-shaped `var()` text
inside an `@function` default no longer becomes a token use, ordinary declaration strings
and URL tokens follow the same collector, and the real-reference controls from round 08
remain visible. The focused tests and my independent scanner matrix both pass in those
directions.

I accept the deliberate non-masking of genuine host-language regions in markup. The
measurement is persuasive: outside `<style>` and `style=`, quotes belong to the host
language, and masking them as CSS strings would lose references in class utilities and
source strings that build applied CSS. My direct controls preserved those uses.

The implementation does not fully implement that stated boundary, however. A `style=`
attribute is parsed once as CSS and then scanned again by the unmasked raw-markup sweep.
That second pass resurrects `var()` text inside an opaque CSS string as an undeclared token.
This is one localized MAJOR in the exact asymmetry this round asked me to judge.

D-1 is accepted. The doctrine remains 65,360/65,536 bytes, every current early-warning
threshold passes, the registry digest matches, and the required content shape is intact.

## What I verified (commands run, and their result)

- `git status --short --branch` and `git show -s 1675b6f` in the implementation repository:
  clean `parley-design-skills` checkout at
  `1675b6f771479f944eb967e50c93eb7e5cf7b3ce`.
- `git show --stat 1675b6f` and the full `9121ec2..1675b6f` diff: cycle 13 changes only
  `lib/css.js` and `test/checker.test.js`; both CSS collectors now use `valueVarUses()`.
- `npm test`: 246 passed, 0 failed. This includes registry refusal, generated capability,
  detector fixtures, conformance gating, waivers, exit codes, packaging, installer paths,
  doctrine budgets, the cycle-13 regression, and its controls.
- Focused `node --test` for the `@function` and opaque-token cases: 2 passed, 0 failed.
- My round-09 reproducer and controls through direct `varUses()` calls: quoted and URL
  defaults produced no `--ghost` use; ordinary strings and URL tokens exposed only the real
  adjacent reference; a real default, an escaped real default, and a mixed default list
  retained `--ghost` or `--gap`; every input remained readable.
- Markup controls through `markupVarUses()`: a `<style>` string exposed only `--real`; a
  host-language string exposed `--error`; and a utility class exposed `--height`. These
  support the deliberate non-masking decision.
- The missing boundary control:
  `<div style="font-family: 'var(--ghost)'; color: var(--real)">` produced both
  `--ghost` and `--real`. The CSS parser itself produces only `--real`; the raw markup pass
  adds `--ghost` back.
- The same inline-style probe inserted into a copied sound L3 fixture and run through
  `runCheck()`: `core:token-used-undeclared`, `VIOLATION`, L3 unverified, no unreadable
  input, exit 1. A real inline-style `var(--ghost)` and the deliberately unmasked
  host-language-string control also exited 1, confirming the probe distinguishes the
  intended boundary.
- `node addons/parley-design-check/bin/check.js --help`: the documented command and
  exit-code surface is present. An explicit nonexistent registry refused rule checks on
  stderr and exited 3 while reporting the generated 18-detector capability.
- `wc -c` and `shasum -a 256`: doctrine files total 65,360 bytes; `RULES.md` hashes to
  `b49ff596451f...`, matching `PDS.md`.

## Findings

### [MAJOR] Inline style attributes are raw-scanned after CSS parsing

**What is wrong:** `markupVarUses()` correctly parses each `style=` value as a declaration
list, but the later raw-line sweep blanks only `<style>` elements and markup comments. It
does not blank the already-consumed style attribute. Therefore:

```html
<div style="font-family: 'var(--ghost)'; color: var(--color-text-body)">x</div>
```

is read once as CSS, where `'var(--ghost)'` is one opaque string token, and again as raw
markup, where the regular expression reports `--ghost`. This contradicts the commit's own
stated boundary: `style=` is CSS, not a host-language region.

**Why it matters:** valid inline CSS receives a false
`core:token-used-undeclared` finding, loses its L3 certificate, and exits 1. The cycle-13
fix therefore still fails on one of the two CSS spans its measurement explicitly excludes
from the deliberate raw-markup trade-off.

**Concrete fix:** have the existing style-attribute consumer return each match's line span,
and blank those exact spans before the raw `VAR_REFERENCE` sweep, while retaining the
references obtained from `declarationVarUses()`. Add the probe above as an end-to-end L3
fixture that must stay clean, plus controls proving a real inline-style reference still
fires and quoted host-language references outside `style=` remain visible. Re-run the
markup corpus measurement to confirm only the double-counted attribute contents change.

## Open questions

None.
