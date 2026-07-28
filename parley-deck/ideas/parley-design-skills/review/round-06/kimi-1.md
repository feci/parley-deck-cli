---
agent: kimi-1
idea: parley-design-skills
review-round: 06
date: 2026-07-28
reviewed-commit: aa6b9b3
---

## Summary

❌ BLOCK. The three closures this cycle claims are **genuine, and I verified each by command
and browser**: codex-1's nested-curly probe now reads the colour *and* fires the fail-safe
(exit 1), his `var()` false positive is gone (exit 0, `--ghost` unreported), hermes-1's
`@scope` bare declaration is read and flagged (exit 1), and my own round-05 `@page` MAJOR is
closed with it (`@page sets margin to the spacing literal 16px`, exit 1). The false-positive
side of the cycle is as good as claimed: 0 unreadable across Bootstrap (280 KB, 2,556 blocks,
5,543 declarations), Google Fonts css2, Tailwind v4 `theme.css` + `preflight.css`, open-props,
modern-normalize, all 37 fixtures, and my 20-block ordinary-CSS stress file. The block model
survived every new construct I threw at it, including a nested-rule/custom-property confusion
sweep and a bad-string recovery probe.

But **my own round-05 CRITICAL is not closed.** Re-run byte-for-byte at `aa6b9b3`:
`#url((y)}`-style constructs still produce verified L3, PASS, exit 0, `unreadable: []` while
headless Chromium computes `rgb(255, 0, 0)` — in all three round-05 spellings plus a fourth I
derived this round (`@\75 rl(`). The block model cannot close it because the desync starts one
layer below the block model: `identLikeToken`'s preceding-context guard (`lib/css.js:270`) is
unchanged from round-05, and the payload still travels the silent top-level `;`-flush discard
channel, also unchanged. The "86 constructs, 0 silent holes" differential either did not
include the prefix class or measured it wrong; the harness itself is committed nowhere in the
tree, so its coverage cannot be reviewed — my round-05 open question 1, unanswered. This is
the sixth time a declared closure did not survive a probe.

Cycle 9 is real progress and the right architecture; it is also one token-class guard short of
what it claims. The fix is the one I prescribed in round-05, still one line plus fixtures.

D-1 position: **ACCEPT, maintained for the sixth round.** 65,360/65,536, `registry-digest`
`b49ff596451f` matches PDS frontmatter, doctrine diff `f1c123d..aa6b9b3` empty, so my round-04
line-by-line lens verification of PDS.md carries in full; spot checks re-confirmed §0–§12, the
C1 mapping table, G1–G4 with checker error strings, and A3 at §4 rule 7.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive aa6b9b3 | tar -x -C /tmp/pds-r6`;
HEAD `aa6b9b3e6482`, branch `parley-design-skills`, worktree clean). Probe runs are copies of
the shipped `sound-run` fixture under `/tmp/pds-probes-r6`; browser checks ran headless Google
Chrome (`--headless=new --dump-dom`, `getComputedStyle` + CSSOM written into the DOM) on pages
under `/tmp/pds-browser`. Corpus CSS fetched into `/tmp/pds-corpus`. Nothing in the skill repo
was modified.

- `npm test` (pristine): **233 passing, 0 failing** (230 at round-05).
- Baseline `check --level L3` on unmodified `sound-run`: PASS, **exit 0**, digest
  `b49ff596451f`, 13 UNJUDGEABLE.
- Doctrine: `git diff f1c123d..aa6b9b3 -- addons/parley-design/` **empty**; `wc -c` = 6,519 +
  25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256 RULES.md` = `b49ff596451f` =
  PDS `registry-digest`. Only `lib/css.js`, three detectors and the test file changed.
- Registry refusal: `--registry /tmp/nope` → "rule checks were refused: no parley-design
  registry was found and this checker carries no copy of one", **exit 3**, structural checks
  still ran. No `RULES.md` under the checker (0 found). 18 detector modules.

**My round-05 CRITICAL, re-run (N1/N3/N4, full `check --level L3`, plus browser):**

- N1 `.a { background: x) #url((y)} z); color: #ff0000; dummy: w) #url((a) {b: (1); }`:
  **verified L3, PASS, exit 0**, `unreadable: []`, zero findings. Chromium:
  `COLOR=rgb(255, 0, 0)`, CSSOM `.a { color: rgb(255, 0, 0); }`. **NOT closed.**
  `scanStylesheet` inspection: `.a` carries only `background: x) #url((y)`; the phantom rule
  `dummy: w) #url((a)` opens; `color: #ff0000` discarded at top level with no note — the exact
  round-05 trace.
- N3 (`@url(`) and N4 (`#\75 rl(`): identical — exit 0 both; Chromium `rgb(255, 0, 0)` both.
  **NOT closed.**
- Sanity control: plain `.a { color: #ff0000; }` as `probe.css` in the same run → literal
  finding, VIOLATION. The probe file is genuinely being scanned.

**codex-1 / hermes-1 round-05 findings, re-run:**

- codex-1 nested-curly, byte for byte: literal finding names `.probe sets color to the colour
  literal #ff0000`, fail-safe fires (3 blocks never closed, paren imbalance, text to EOF),
  **exit 1**. **Closed.**
- codex-1 `var()` false positive (`.supports-\[var\(--ghost\)\]`): **PASS, exit 0**, `--ghost`
  nowhere in the report. `var()` uses now come from `declarations[].value` only. **Closed.**
- hermes-1 `@scope (.probe) { color: #ff0000; }`: literal finding names
  `@scope (.probe) sets color…`, **exit 1**; the bare declaration is attributed to a real
  block. **Closed.**

**My round-05 MAJOR and MINORs, re-run:**

- `@page { margin: 16px; }` → `core:literal-outside-token-layer` names
  `@page sets margin to the spacing literal 16px`, **exit 1**. **Closed** — declaration-at-rules
  are declaration blocks now.
- Case probes (`.u1 { color: VAR(--primitive-ink) }`, `.u2 { color: RGB(255, 0, 0) }`):
  **PASS, exit 0**. Still open (MINOR).
- `content: "say \"hi\""`: whole file unreadable, every style rule UNJUDGEABLE, **exit 4**.
  Still open (MINOR).
- `fill: url(#fade)`: false VIOLATION "sets fill to the colour literal #fade", **exit 1**.
  Still open (MINOR).
- `waivers: WAVERIS.md` typo on the CONTRACT: **PASS, exit 0**, the word "waiver" nowhere in
  the output. Still open (MINOR).
- Second critique round at `round-03/hermes-1.md`: still fails `pds-check:l2-process-order`,
  exit 1. Still open (MINOR).

**Block-model attack probes (all `scanStylesheet` direct):**

- codex-1-class siblings: bracket-nested `[{]`, two-deep `fn({{)}}`, and the custom-property
  value block `.a { --x: { color: #ff0000 }; }` — all read correctly or fail-safe; the shipped
  tests pin them. No hole found.
- Escaped-hash spelling `\#url((y)}…` (my attempt at a fifth prefix variant): **safe** — the
  escape starts an ident sequence at the backslash, `\#url` spells `#url` ≠ `url`, so no url
  token opens; the file goes unreadable and the literal is found (exit 1). The open class is
  exactly a *literal* `#` or `@` before a literal-or-escaped `url(`.
- New variant N6 `@\75 rl((y)}…`: **verified L3, PASS, exit 0** — the escaped-`url` spelling of
  N3. Same hole, fourth spelling (in CRITICAL).
- Bad-string recovery (`.a { content: "oops\ncolor: #ff0000; }`): the scanner merges the
  recovered text into one `content` declaration; the literal stays visible inside the value, so
  value-scanning detectors still see it. Pre-existing, accidental safety, no hole found for the
  literal/token rules.
- Nested-rule vs custom-property confusion: `&:hover`, `.b`, `--x: {…}`, escaped `--x`
  (`\2d\2d x`) all classified per Chromium's reading (the shipped suite pins the same).
- At-rule classification: `@media` bare declaration → both discard, no note (agree);
  `@position-fallback { @try { … } }` and `@function --f() { result: … }` → **false
  unreadable** (new MINOR); `@scope` with bare declaration + nested rules → all read.

**False-positive sweep (the risk this cycle took on):**

- Bootstrap 5 (280 KB): 2,556 blocks, 5,543 declarations, **0 unreadable**. Google Fonts css2
  (20 `@font-face`): 0. Tailwind v4 `theme.css` (`@theme default`, 427 declarations) and
  `preflight.css`: 0 — the round-05-cycle @theme false alarm is genuinely gone. open-props,
  modern-normalize: 0. All **37 fixture stylesheets**: 0.
- My ordinary-CSS stress file (braces in strings and attribute selectors, `url(#fragment)`,
  data-URI svg, Tailwind escaped classes incl. `.supports-\[var\(--ghost\)\]`, `calc`/`minmax`
  nesting, keyframes, `@media`/`@container`/`@supports selector()`/`@starting-style`/`@scope`,
  `@property`, CSS nesting, custom-property value block, `@layer` statement, `@import …
  supports()`, `@namespace`): 20 blocks all read with correct at-contexts; **one** note — the
  known escaped-quote one (carried MINOR).

**Carried NITs:** `TIER_GROUPS` (`engine.js:1367`) now maps singulars **and** plurals —
**closed**. The "will not take on trust" table still lists five rows (open). PDS.md
25,594/25,600 — 6 bytes headroom (open). IMPLEMENTATION.md now records cycles 7, 8 and 9 with
their own commits and counts — the round-04/05 record-keeping MINOR is **closed**.

## Findings

### [CRITICAL] (round-05, re-confirmed at aa6b9b3) `#url(` / `@url(` — the url-ident guard still ignores preceding-token context; four spellings certify L3 over a colour Chromium applies

**What.** `identLikeToken` (`lib/css.js:270`) still rejects a url reading only when the
previous code point is an ident code point — the exact line my round-05 finding named,
unchanged. css-syntax §4.3.3/§4.3.4 make a literal `#` or `@` bind the following ident
sequence into a hash token or an at-keyword token, so after either, **no ident-like token
exists and no url token can start** — the `(` opens an ordinary `()`-block in which `}` is
content. The scanner reads `url(a…)` as one url token, ends it at the first `)`, and meets the
`}` at rule depth: the rule closes early, `color: #ff0000` is discarded at top level
(silently — the top-level `;`-flush records nothing when `enclosingAt()` is null), a phantom
rule opens at the re-balancing `{`, and the file ends balanced. Every residue check stays
quiet. Re-confirmed this round in **four** spellings, each beside `sound-run` at
`--level L3`, each with **verified L3, PASS, exit 0, `unreadable: []`** while headless
Chromium computes `rgb(255, 0, 0)`: N1 `#url(`, N3 `@url(`, N4 `#\75 rl(`, and the new N6
`@\75 rl(`. The escaped-preceding forms (`\#url(`) are safe — the escape starts the ident
sequence at the backslash and the spelled value `#url` is not `url` — so the class is pinned:
a literal `#` or `@` immediately before a literal-or-escaped `url(`.
**Why it matters.** This is the failure signature the fail-safe exists to eliminate, reachable
in any stylesheet from one typo'd hash-colour or at-keyword, needing no malformed input beyond
what a minifier emits. It refutes the load-bearing claim of this cycle — "86 constructs …
before 5 silent holes, after 0" — on three (likely five) of the very constructs that made up
the "before" count, re-run at the reviewed commit. And the differential harness behind that
claim is committed nowhere in the tree (searched: no file mentions it), so its construct
matrix cannot be reviewed; the only reproducible coverage is the dozen constructs pinned in
`checker.test.js`, none of which is this class. A verified-L3 certificate beside a raw colour
the ratified system really contains is the worst output this checker can produce.
**Fix.** Unchanged from round-05, and it covers all four spellings:

1. In `identLikeToken`, return null when `text[index - 1]` is `#` or `@` — those two code
   points merge the ident sequence into a hash or at-keyword token, so the `(` is an ordinary
   `()`-block opener, which the block model already nests correctly.
2. Defence in depth, because this desync's payload travels the *silent discard channel*, not
   the residue the fail-safe measures: a `;`-flush (or EOF flush) at top level whose buffer is
   non-empty and does not start with `@` is discarded text — report it through the unreadable
   channel, exactly as in-rule discards already are. The `@import …;` / `@charset …;` /
   `@layer a, b;` statement forms keep their exemption via the `@`.

Fixtures: N1, N3, N4 and N6 above, each of which must fail or go UNJUDGEABLE, plus
passing-side controls — `background: url(#fragment)`, `color: #abc`, `@media` nesting, and
`\#url(` (already safe; pin it). Add the constructs to `checker.test.js` beside
`NESTED_CURLY_BLOCK`, and if the 86-construct harness is the evidence this family ships on,
commit it or its matrix so the next round can review its dimensions instead of its conclusion.

### [MINOR] (new) `@position-fallback`/`@try` and `@function` hit the unclassifiable at-rule guard — false unreadable on shipped CSS

**What.** `@position-fallback --fb { @try { bottom: 0; left: 0; } }` and
`@function --negate(--v) { result: calc(var(--v) * -1); }` each produce
"the declaration … sits directly inside @try / @function, whose body this scanner can read
neither as rules nor as declarations" — the whole file unreadable, every style rule
UNJUDGEABLE, exit 4. Anchor positioning (`@position-try` is already modelled, its
`@position-fallback` wrapper and `@try` blocks are not) is shipped Chromium CSS, and CSS
`@function` is a declaration-holding at-rule entering the same engines; both bodies are
declarations the browser applies.
**Why it matters.** Safe direction — the guard is doing its job on a construct it was never
taught — but it is the same false-alarm species as the Tailwind `@theme` one this cycle
caught before shipping, on standard rather than framework CSS. A gate that fires on ordinary
shipped constructs teaches readers to waive reflexively.
**Fix.** Add `try` to `DECLARATION_AT_RULES` and `position-fallback` to `RULE_AT_RULES`; add
`function` to `DECLARATION_AT_RULES`. Fixtures: both probes above (no unreadable, declarations
visible), and an honestly unclassifiable at-rule (e.g. `@future-thing { color: #ff0000 }`)
still guarded, exit 4.

### [MINOR] (carried) Case is a spelling the decoded layer does not cover: `VAR(--x)`, `RGB(255,0,0)` pass clean while the browser substitutes and applies

Re-probed at `aa6b9b3`: PASS, exit 0; Chromium substitutes `VAR()` and applies `RGB()`.
`varUses`, `COLOUR_LITERAL` and `LENGTH_LITERAL` match case-sensitively; CSS matches all three
case-insensitively. **Fix** (unchanged): the `i` flag on the function-name and unit
alternatives, and a case-insensitive `var(` match on the decoded value. Fixtures: the two
probes (flagged), the lowercase forms (unchanged).

### [MINOR] (carried) An escaped quote or backslash inside a string still makes the whole file unreadable — `content: "say \"hi\""` is ordinary CSS and gets exit 4

Re-probed, and it is the **only** note my 20-block ordinary-CSS stress file produced. The
string token itself parses correctly (`stringToken` isolates it, escape-aware); the refusal
protects value consumers that split on those code points, but those consumers already see the
raw `\"` via the verbatim form. **Fix** (unchanged): keep the token verbatim but do not route
quote/backslash escapes through the unreadable channel; reserve the note for escapes spelling
code points that change value structure (comma, braces, parens). Fixtures: the probe (no
unreadable, no false literal), an escaped-comma string (still noted).

### [MINOR] (carried) The literal matchers still scan `url()` and string contents — `fill: url(#fade)` is a false VIOLATION, exit 1

Re-probed at `aa6b9b3`. `#fade` inside `url()` is an SVG paint-server reference, not a colour.
**Fix** (unchanged): mask `url(…)` spans and string tokens before literal matching — the
scanner already knows both boundaries. Fixtures: the probe (clean), `fill: #fade` outside a
url (still flagged).

### [MINOR] (carried) A contract whose named waiver file resolves to no readable file is still silently treated as "no waivers"

Re-probed: `waivers: WAVERIS.md` → verified L3, PASS, exit 0, the file never mentioned.
Direction remains safe (waivers only suppress), but the run's inputs are misstated in silence.
**Fix** (unchanged): report it — NEEDS_REVIEW or an `l2-process-order` finding naming the
contract and path.

### [MINOR] (carried) §4 rule 4's second critique round still has no §1 home

Re-probed: a Decider-instructed second-round CRITIQUE at `round-03/hermes-1.md` fails
`pds-check:l2-process-order`, exit 1. **Fix** (unchanged): §1 names the second round's home
(and `PROCESS_HOMES` accepts it when the instruction is recorded), or §4 rule 4 states the
second round also files under `round-02`. Fixture for the permitted shape.

### [NIT] (new) The checker SKILL.md's scanner section still says "seven constructs … each read as text now"

`parley-design-check/SKILL.md:88` enumerates seven constructs and states each is "read as text
now rather than as structure". The list predates round-05's nested curly (an eighth) and the
prefix class above (a ninth, open). When the CRITICAL lands, add both — the section is the
scanner's public contract and its construct list is the one place a reader can check the claim
against the probes.

### [NIT] (carried) The "will not take on trust" table still lists five conditions

`parley-design-check/SKILL.md:182`. Cycle 4–6's recomputed conditions (deduplicated rotation,
G1 axis enumeration, process-order locations, G3/G4 refutations, waiver-ownership exclusion,
the sidecar check) are still not rows. Add the rows or retitle the table.

### [NIT] (carried) PDS.md has 6 bytes of headroom against its early-warning threshold

25,594 of 25,600 (total 65,360/65,536, 176 spare). Unchanged this cycle; recorded so the next
normative sentence is planned for.

## Open questions

1. **(carried, sharpened) The differential harness.** Round-05 asked what generated the
   62-construct matrix; cycle 9 cites an 86-construct run, and the harness is still committed
   nowhere I can find. Three of its claimed "after: 0" constructs reproduce as holes at the
   reviewed commit. Will the owner commit the harness or its matrix? "0 silent holes" is now
   the ship criterion; the evidence for it has to be as reviewable as the code.
2. **(answered) Round-05 question 1's substance** — the prefix class was not covered by the
   block-model replacement. Recorded via the CRITICAL, not answered by the owner.
3. **(carried) The level/verdict axis split.** An open, honestly-recorded `slop` violation
   against the winner still yields `verified L3` beside VIOLATION, exit 1. Never confirmed by
   the owner; if my reading is wrong, `gateRefutations` needs a G1 arm.
4. **(carried) The second-round critique home** — MINOR above: where does §4 rule 4's
   exception live?
