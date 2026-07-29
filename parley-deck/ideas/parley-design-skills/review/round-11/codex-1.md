---
agent: codex-1
idea: parley-design-skills
review-round: 11
date: 2026-07-28
reviewed-commit: 1804985
---

## Summary

✅ ACCEPT. Cycle 14 closes my round-10 finding. The exact inline-style reproducer now
contributes only the real reference from the CSS declaration path; opaque string and URL
contents stay hidden, while real inline-style, escaped-inline-style, `<style>`-body and
host-language controls all retain the intended result.

The 72 of 8,307 `style=` spans left in the raw sweep are a **disclosed, bounded residual,
not a defect**. They are precisely the values the CSS scanner could not establish as one
readable declaration list, including host-language/template fragments that later build real
inline CSS. Sending all 72 to the unreadable channel would turn every such file into
`UNJUDGEABLE`/exit 4 even when the raw sweep preserves a real reference. Raw scanning can
over-report the narrower case where one of those fragments later becomes an opaque CSS token;
my control reproduced that known edge. It cannot create the false-clean direction, is under
one percent of observed spans, and is disclosed with the measured counterfactual. That is a
shipping state.

D-1 is accepted. The written per-file split was deviated from, but the ratified 64 KiB total
is held at 65,360/65,536 bytes, every current early-warning threshold passes, and the
identical four-part artifact shape was preserved. D-2 is also honest: capability is generated
from 18 detector modules, and all nine check/both rules without a detector remain visible
rather than silently passing.

## What I verified (commands run, and their result)

- `git status --short --branch`, `git rev-parse HEAD`, and `git show -s 1804985` in the
  implementation repository: clean `parley-design-skills` checkout at
  `1804985f91975f075198f087e708d2b4b766dce4`.
- Full `git show 1804985` and `git diff --check main...HEAD`: cycle 14 changes only
  `lib/css.js` and `test/checker.test.js`; the span withheld from the raw pass is the exact
  value span returned by the declaration-list consumer. No whitespace errors.
- `npm test`: **247 passed, 0 failed**. This covers registry refusal and parsing, generated
  capability, every detector fixture pair, gate obligations, waiver independence and scope,
  exit codes, doctrine budgets, packaging, installer paths, and the cycle-14 regression.
- Focused
  `node --test --test-name-pattern='a style attribute is read as CSS once' .../checker.test.js`:
  **1 passed, 0 failed**. It includes scanner and end-to-end L3 assertions.
- My round-10 scanner reproducer, rerun directly through `markupVarUses()`:
  `<div style="font-family: 'var(--ghost)'; color: var(--real)">` now returns only
  `--real@1`. A real inline reference returns once; `\76 ar(--escaped)` returns
  `--escaped`; a host-language string still returns its reference; and a `<style>` string
  stays opaque while the adjacent real declaration is found.
- Residual controls through the same exported function: the JS-template attribute retains
  `--c1`; the Jinja conditional retains `--brand` and `--neutral`; and a deliberately
  constructed Jinja branch containing `font-family: 'var(--opaque)'` demonstrates the
  disclosed over-reporting edge. The boundary therefore behaves in both directions exactly
  as disclosed.
- `node .../bin/check.js --help`: command surface and exit-code table match the skill.
  An explicit nonexistent `--registry` refused rule checks on stderr, still emitted the
  generated 18-detector capability, reported `UNJUDGEABLE`, and exited **3**.
- A direct registry/capability comparison found 18 detectors and the documented nine
  check/both rules without one. The suite confirms each is reported `UNJUDGEABLE` under its
  applicable reason.
- `wc -c` and `shasum -a 256`: doctrine total **65,360 bytes**; `RULES.md` digest starts
  `b49ff596451f`, matching `PDS.md`. The placeholder scan was empty.
- I reread the doctrine, protocol, registry, web annex, checker skill, registry parser,
  artifact/token reader, detector modules, engine gating and waiver paths, CLI, tests,
  installer changes and `NOTICE.md`. I found no unsupported capability claim or new
  correctness issue.

## Findings

No findings. The disclosed 72-span residual is bounded and correctly left on the
host-language raw path; routing it wholesale to the unreadable channel would be a broader
false alarm.

## Open questions

None.
