---
agent: kimi-1
idea: parley-design-skills
review-round: 05
date: 2026-07-28
reviewed-commit: f1c123d
---

## Summary

❌ BLOCK. Round-04's two CRITICALs are **genuinely closed**, and I say so on commands, not
on the report: my own `()`/`[]` nesting reproducers (ev1–ev3) now surface the hidden literal
with the fail-safe recording residue where residue exists, and the escaped-`url` family
(codex-1's probe plus hermes-1's first three variants) is caught the same way. The
raw-vs-decoded layer is real too — `col\6fr`, `#\66 f0000`, `\72 ed`, `11p\78` are all caught
with the written spelling named beside the finding, and the block model now claimed in
`SKILL.md` is exactly the fix I prescribed (one matched-bracket stack; `{`/`}`/`;` structural
only at rule depth; mismatched closers preserved per §5.4.7). The doctrine files are
byte-identical to round-04, so my round-04 line-by-line lens verification of PDS.md carries.

But the load-bearing claim of this round — **0 silent holes** — is refuted by a command and
its output. There is an eighth construct in the tokenisation family, and it is not an escape
exotic: a `url` ident immediately preceded by `#` or `@` is not an ident token at all in a
browser (it is the value of a hash or at-keyword token, §4.3.3/§4.3.5), so the `(` opens a
nesting `()`-block — while `identLikeToken`'s guard only excludes preceding ident code
points, so the scanner reads a url token, ends it at the first `)`, and the structural
desync begins. With one re-balancing `{` every residue check stays quiet: **verified L3,
PASS, exit 0, `unreadable: []`**, while headless Chromium computes `rgb(255, 0, 0)`. Three
variants confirmed in the browser. The fail-safe does not fire because the desync leaves no
residue it measures — the same failure signature as round-04, reached through a door the
escape fix did not touch: preceding-token context decides token class.

One correction owed to the record: hermes-1's fourth round-04 variant, `\75\72\6cl`, spells
`urll` (each `\XX` is a separate escape, then a literal `l`), so in a browser the `/*` opens
a real comment and the colour is never applied — Chromium confirms `rgb(0, 0, 0)`. The
checker's current no-finding on that variant is the *correct* verdict; that variant of the
round-04 claim was wrong about browser semantics.

D-1 position: **ACCEPT, maintained** for the fifth round. 65,360/65,536 (176 spare),
`RULES.md` byte-identical (digest `b49ff596451f` matches PDS frontmatter), per-file
thresholds unchanged at 7/25/24/11 KiB, doctrine diff `e3ca916..f1c123d` empty.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive f1c123d | tar -x -C /tmp/pds-r5`;
HEAD `f1c123d41b3d`, branch `parley-design-skills`, worktree clean). Probe fixtures were
copies of the shipped `sound-run` fixture under `/tmp/pds-probes-r5`; browser checks ran
headless Google Chrome (`--headless=new --dump-dom`, `getComputedStyle` written into the
DOM) on pages under `/tmp/pds-browser`. Nothing in the skill repo was modified.

- `npm test` (pristine): **230 passing, 0 failing** (227 at round-04).
- Baseline `check --level L3` on unmodified `sound-run`: `claimed L3, verified L3`, PASS,
  **exit 0**, digest `b49ff596451f`, 13 UNJUDGEABLE, `recusal-not-anchored` for hermes-1.
- Budgets/digests: `wc -c` = 6,519 + 25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**;
  `shasum -a 256 RULES.md` = `b49ff596451f`; `git diff e3ca916..f1c123d -- addons/parley-design/`
  is empty (doctrine untouched); thresholds in `design-addons.test.js` unchanged; the
  detector diffs are `asWritten(...)` annotations plus matching on the decoded fields, as
  the css.js header promises.

**My round-04 CRITICAL, re-run (all four variants, full `check --level L3`):**

- ev1 (`background: x) \75 rl(}y); color: #ff0000; dummy: z) \75 rl({w: (1); }`):
  `core:literal-outside-token-layer` names `.a sets color to the colour literal #ff0000`,
  G3 refuted, `not verified`, VIOLATION, **exit 1**, unreadable note on the unbalanced
  fragment. **Closed.**
- ev2 (plain `fn(}y)`, no escapes) and ev3 (`url (` with a space): identical outcome —
  literal detected, exit 1. **Closed.**
- ev4 (control, no re-balancing brace): fail-safe fires ("does not balance its
  parentheses…", "1 block opened and never closed"), exit 1. The mechanism still works for
  the shapes it can see. **Closed.**

**codex-1 / hermes-1 round-04 CRITICAL, re-run:**

- codex-1's probe (`u\72l(a/*) ; color: #ff0000; */b)`): literal detected, unreadable notes
  on the `*/b)` fragment, G3 refuted, `not verified`, **exit 1**. **Closed.**
- hermes-1's variants `u\72l`, `\75 rl`, `u\000072l`: same. **Closed.**
- hermes-1's variant `\75\72\6cl`: no literal finding — and that is now demonstrably right:
  the ident spells `urll`, so the `(` is a function and `/*…*/` is a comment in the browser
  too. Chromium: `color=rgb(0, 0, 0)` (not applied), `padding=8px` (rule alive). The
  scanner's decoded-spelling logic is more accurate here than round-04's probe suite was.

**The new attack (CRITICAL below), three variants, checker and browser:**

- N1 `.a { background: x) #url((y)} z); color: #ff0000; dummy: w) #url((a) {b: (1); }` —
  checker: **verified L3, PASS, exit 0, `unreadable: []`**, zero findings. Chromium:
  `rgb(255, 0, 0)`. `scanStylesheet` inspection: `.a` carries only
  `background: x) #url((y)`; `color: #ff0000` discarded as top-level text with no note; a
  phantom block carries selector `dummy: w) #url((a)`.
- N3 (same with `@url(`): identical — exit 0, browser applies.
- N4 (same with `#\75 rl(`, escaped spelling after the hash): identical — exit 0, browser
  applies.

**Secondary probes (all beside `sound-run` at `--level L3` unless noted):**

- `@page { margin: 16px; }`: **verified L3, PASS, exit 0** — the declaration is dropped at
  flush because an at-block is not a "rule" (`currentRule()`), with no unreadable note.
  Chromium's CSSOM carries `CSSPageRule:@page { margin: 16px; }`. A governed spacing literal
  the browser applies, unread and unflagged. (MAJOR.)
- `@font-face { font-family: …; src: … }`: dropped the same way, silently; compensated in
  practice because `face-allowlist` judges the *usage* in ordinary rules.
- `@property --xp { … initial-value: #ff0000 … }` + `.usexp { color: var(--xp) }`: Chromium
  applies `rgb(255, 0, 0)` via the registered initial; the checker flags the run anyway —
  `core:token-used-undeclared` on `--xp`, exit 1. Compensated, for the adjacent reason.
- `@keyframes k { from { color: #ff0000; } }`: read correctly (block `from`, at-context
  `keyframes`); the literal is visible to detectors.
- Claimed catches hold end-to-end: `\76 ar(--nope, #ff0000)` → `token-used-undeclared`,
  exit 1; the four decoded spellings → four literal findings with `(written …)` annotations.
- Case probes: `.u1 { color: VAR(--primitive-ink) }` and `.u2 { color: RGB(255, 0, 0) }` —
  **verified L3, PASS, exit 0**; Chromium substitutes `VAR()` (`rgb(0, 255, 0)`) and applies
  `RGB(` (`rgb(255, 0, 0)`). `varUses` and `COLOUR_LITERAL`/`LENGTH_LITERAL` are
  case-sensitive; CSS function names, keywords and units are not. (MINOR.)
- `content: "say \"hi\""` (ordinary CSS): the whole file goes unreadable, every style rule
  UNJUDGEABLE, **exit 4**. Safe direction, false unreadable on ordinary input. (MINOR.)
- `fill: url(#fade)`: false VIOLATION "sets fill to the colour literal #fade" — the literal
  matcher scans `url()`/string contents; `#fade` there is an SVG paint-server reference, not
  a colour. (MINOR.)
- Fixture corpus: all **37 fixture stylesheets** scanned with `scanStylesheet` directly —
  **0 unreadable**. My hand-built ordinary-CSS stress file (braces in strings and attribute
  selectors, `url(#fragment)`, data URIs, Tailwind escaped classes, `calc`/`minmax` nesting,
  keyframes, media queries, escaped `content` glyph): parsed to 18 correct blocks with one
  note — the escaped-quote one above. The "0 new findings / 0 lost / 0 newly unreadable"
  claim is consistent with what I see on the corpus; the two false-positive-class items
  above live outside it.

**Carried items, re-probed:**

- MINOR-1 (round-04): `waivers: WAVERIS.md` (typo) on the CONTRACT → **verified L3, PASS,
  exit 0**, the word "waiver" nowhere in the output; the named file silently unread. **Still
  open.**
- MINOR-3 (round-04): a Decider-instructed second-round CRITIQUE at `round-03/codex-1.md`
  still fails `pds-check:l2-process-order` ("filed outside §1's mapping"). **Still open.**
- NIT-1: `TIER_GROUPS` (`engine.js:1367`) still maps plurals; PDS §3 G3 still names
  singulars. **Open.**
- NIT-2: the "will not take on trust" table (`parley-design-check/SKILL.md:187`) still lists
  five rows — cycle 7-8's recomputed conditions (the unreadable roll-up is documented in the
  new scanner section, but the table that reads as the complete list is not it). **Open.**
- NIT-3: PDS.md 25,594/25,600 — 6 bytes of headroom, unchanged. **Open.**
- MINOR-2 (round-04), extended: `IMPLEMENTATION.md` frontmatter says `fix-up-cycle-8`,
  `head-commit: f1c123d`; the body still ends at cycle 6. Cycles 7-8 — the cycles under
  review — have no record: no fixes-applied, no deviation notes, no verification counts.

**Doctrine lens:** the `e3ca916..f1c123d` doctrine diff is empty, so the mapping table,
G1–G4, the ritual, roles/recusal, L1–L4 and the eight four-part artifact entries stand as
verified in round-04. `parley-design-check/SKILL.md`'s new scanner section states the block
model and the spelled-vs-written contract plainly and answers my round-04 open question 3
exactly as proposed — the scanner's block model is now a specified, reviewable artifact.
That section will need a tenth construct added to its list, or a ninth and the prefix class.

## Findings

### [CRITICAL] `#url(` / `@url(`: the url-ident guard ignores preceding-token context — a brace after a hash or at-keyword closes the rule in the scanner and not in the browser, certifying clean over an applied literal

**What.** `identLikeToken` (`lib/css.js:254`) rejects a url reading only when the previous
code point is an ident code point. css-syntax §4.3.3/§4.3.5 make `#` and `@` bind the
following ident sequence into a hash token or an at-keyword token, so after either, **no
ident-like token exists and no url token can start** — the `(` opens an ordinary `()`-block,
and inside that block a `}` is content (§5.4). The scanner reads `url(a…)` as one url token,
ends it at the first `)`, and then meets the `}` at rule depth: it closes the rule early,
drops `color: #ff0000` as top-level text (silently — the top-level `;`-flush records
nothing), opens a phantom rule at the re-balancing `{`, and ends balanced. Every residue
check stays quiet: `unreadable: []`, **verified L3, PASS, exit 0** (N1/N3/N4 above), while
Chromium computes `rgb(255, 0, 0)` on the probe element in all three variants — plain
`#url(`, `@url(`, and the escaped `#\75 rl(`. No malformed input is needed beyond what any
minifier or hash-colour typo can produce, and the escaped form needs nothing at all.
**Why it matters.** This is the exact failure signature the fail-safe was shipped to
eliminate, and the claim this round rests on — "62 tricky constructs … 0 silent holes" — is
the claim it refutes. Four times this family has been declared closed one probe too early;
the differential harness that ran this round did not include preceding-token context as a
dimension, which is the layer the escape fix left untouched. A verified-L3 certificate
beside a raw colour the ratified system really contains is the worst output this checker
can produce, and it is reachable in any stylesheet on any element.
**Fix.** One line plus its tests, then a second line of defence:

1. In `identLikeToken`, return null when `text[index - 1]` is `#` or `@` — those two code
   points merge the ident sequence into a hash or at-keyword token, so the `(` is an
   ordinary `()`-block opener, which the block model already nests correctly. That is the
   whole primary fix: the guard's job is "no ident-like token starts here", and ident code
   points are not the only context that decides it.
2. Defence in depth, because this desync's payload travelled the *silent discard channels*,
   not the residue the fail-safe measures: (a) a `;`-flush at top level whose buffer is
   non-empty and does not start with `@` is discarded text — report it through the
   unreadable channel, exactly as in-rule discards already are (`@import …;` and
   `@charset …;` keep their exemption via the `@`); (b) declarations directly inside
   declaration-at-rules (`@page`, `@font-face`, `@property`, `@counter-style`) are currently
   dropped with no block and no note — record them or report the discard (see the MAJOR).

Fixtures: N1, N3 and N4 above, each of which must fail or go UNJUDGEABLE, plus passing-side
controls — `background: url(#fragment)`, hash-colour values (`color: #abc`), `@media`
nesting, and the shipped suite. Add the constructs to the differential harness along with a
prefix dimension (`#`, `@`, digit, `-`, `.`, `#`-plus-escape) so the class, not the
instance, is pinned.

### [MAJOR] Declarations directly inside declaration-at-rules are silently discarded — `@page { margin: 16px }` certifies clean

**What.** Inside an at-block, `currentRule()` finds no `rule` entry, so `flushDeclaration`
drops the buffer with no block and no unreadable note. `@page { margin: 16px; }` beside
`sound-run`: **verified L3, PASS, exit 0**, while Chromium's CSSOM carries
`CSSPageRule { margin: 16px; }` — a governed spacing literal the literal rule would have
flagged had the declaration been recorded anywhere. `@font-face` and `@property` blocks drop
the same way; those two happen to be compensated (face-allowlist judges usage in rules;
`token-used-undeclared` fires on the `var(--xp)` reference, exit 1, confirmed), but `@page`
margins have no compensating control, and the scanner's own contract — "it reports what it
could not read as well as what it read" — is broken either way: this text it read fine and
threw away.
**Why it matters.** Narrower than the CRITICAL (print/paged-media surface, one governed
property family), but the same species: a browser-applied declaration the checker neither
reads nor marks unreadable, one round after "0 silent holes" became the ship criterion.
**Fix.** Record at-block direct declarations on the at-rule itself (a synthetic block whose
selector is the at-prelude, so detectors see them in context), or route the discard through
the unreadable channel. Fixtures: the `@page` probe (must be flagged or UNJUDGEABLE), a
`@font-face` block (no crash, no false literal), and `@keyframes` (must keep reading its
inner rules exactly as it does today).

### [MINOR] Case is a spelling the decoded layer does not cover: `VAR()`, `RGB()`, `11PX` pass clean while the browser substitutes and applies

**What.** `varUses` matches `var\(` case-sensitively; `COLOUR_LITERAL` and `LENGTH_LITERAL`
match function names and units case-sensitively. CSS matches all three case-insensitively.
Confirmed in Chromium: `VAR(--decl)` substitutes (`rgb(0, 255, 0)`), `RGB(255, 0, 0)`
applies (`rgb(255, 0, 0)`); confirmed at the checker: `.u1 { color: VAR(--primitive-ink) }`
and `.u2 { color: RGB(255, 0, 0) }` → **verified L3, PASS, exit 0** — the undeclared-token
edge and the literal both invisible. The raw-vs-decoded work this cycle fixed *what* is
matched (spelled vs written) but not the case fold, which is the same class of
spelling-vs-meaning gap one level down.
**Why it matters.** Silent misses in the two rules the L3 certificate leans on hardest;
pre-existing, but squarely inside the family this round claims closed.
**Fix.** Add the `i` flag to `COLOUR_LITERAL`'s function alternative and `LENGTH_LITERAL`'s
unit alternative, and match `var(` case-insensitively in `varUses` (the decoded line, not
the raw). Fixtures: the three probes above (flagged), the lowercase forms (unchanged).

### [MINOR] An escaped quote or backslash inside a string makes the whole file unreadable — `content: "say \"hi\""` is ordinary CSS and gets exit 4

**What.** `decodeStringToken` refuses escapes that spell anything but an ident code point or
a space, so a string containing `\"` or `\\` is kept verbatim **and** routes the file to the
fail-safe: the probe above produced UNJUDGEABLE for every stylesheet rule, `not verified`,
**exit 4**. The string token itself was parsed correctly — `stringToken` isolates it, escape
aware — so the file's tokenisation is not actually in doubt; the note exists to protect
value consumers that split on those code points, but those consumers already see the raw
`\"` today via the verbatim form, which is exactly what would confuse a naive splitter.
**Why it matters.** The fail-safe's cost model is "the construct nobody found costs rules
their verdict"; here the construct is ordinary authored CSS, so the cost lands on honest
files — a stylesheet with an escaped quote in a `content` string can never receive a
verdict. Safe direction, wrong file.
**Fix.** Keep the token verbatim (consumers cope with raw CSS text already) but do not route
quote/backslash escapes through the unreadable channel; reserve the note for escapes
spelling code points that change value structure for known consumers (comma, braces,
parens). Fixtures: the probe (no unreadable, no false literal), and an escaped-comma string
(still noted).

### [MINOR] The literal matchers scan `url()` and string contents — `fill: url(#fade)` is a false VIOLATION

**What.** `COLOUR_LITERAL` and `LENGTH_LITERAL` run over the whole decoded value, including
url token contents (kept verbatim by design) and string contents. `fill: url(#fade)` — an
SVG paint-server reference — produced "sets fill to the colour literal #fade", G3 refuted,
**exit 1**. Same shape: `background: url("16px.png")` would trip `LENGTH_LITERAL`.
**Why it matters.** Round-05's false-positive question is the mirror of the hole question:
a checker that cries literal on references teaches readers to waive reflexively, and every
reflexive waiver is a door for a real one. Pre-existing, but it is in the code under review.
**Fix.** Mask `url(…)` spans and string tokens before literal matching (the scanner already
knows both boundaries), then match. Fixtures: the probe (clean), `fill: #fade` outside a url
(still flagged).

### [MINOR] (carried, extended) IMPLEMENTATION.md records nothing past fix-up cycle 6; cycles 7-8 — the cycles under review — are unrecorded

Frontmatter says `fix-up-cycle-8`, `head-commit: f1c123d`; the body ends at cycle 6. The
round-05 claims (62-construct differential, false-positive diff) reached reviewers through
the deck-side prompt, not the record whose job is to carry them — the same defect as
round-04, one cycle worse. **Fix:** append cycle-7 and cycle-8 sections: fixes applied,
deviations from literal remedies with reasons, verification counts, and the per-cycle
commits. Five lines per cycle. When the CRITICAL above lands, that cycle belongs in the same
record before the next review round, not after it.

### [MINOR] (carried) A contract whose named waiver file resolves to no readable file is still silently treated as "no waivers"

Re-probed: `waivers: WAVERIS.md` → verified L3, PASS, exit 0, the file never mentioned.
Direction remains safe (waivers only suppress), but the run's inputs are misstated in
silence. **Fix** (unchanged): a non-empty `waivers` value resolving to no readable file is
reported — NEEDS_REVIEW or an `l2-process-order` finding naming the contract and path.

### [MINOR] (carried) §4 rule 4's second critique round still has no §1 home

Re-probed: `round-03/codex-1.md` CRITIQUE fails `l2-process-order`. **Fix** (unchanged):
§1 names the second round's home (and `PROCESS_HOMES` accepts it when the instruction is
recorded), or §4 rule 4 states the second round also files under `round-02`. Fixture for the
permitted shape.

### [NIT] (carried) Alias-direction group names: checker accepts plurals, doctrine names singulars

`TIER_GROUPS` (`engine.js:1367`) maps `primitives`/`semantics`/`components`; PDS §3 G3 names
only singulars. One line either side.

### [NIT] (carried) The "will not take on trust" table still lists five conditions

`parley-design-check/SKILL.md:187`. The new scanner section documents the block model and
the unreadable roll-up, but the table that reads as the complete list of recomputed
conditions still omits cycle 4-6's additions (deduplicated rotation, G1 axis enumeration,
process-order locations, G3/G4 refutations, waiver-ownership exclusion, the sidecar check).
Add the rows or retitle the table. Its construct count will also need updating when the
CRITICAL lands.

### [NIT] (carried) PDS.md has 6 bytes of headroom against its early-warning threshold

25,594 of 25,600 (total 65,360/65,536, 176 spare). Unchanged this cycle; recorded so the
next normative sentence is planned for.

## Open questions

1. **(new) The differential harness's construct matrix.** The 62-construct run missed the
   prefix dimension entirely (`#`/`@` before `url(`, digit/`-`/`.` prefixes, hash-plus-escape).
   What generated the matrix, and will the owner publish it? A harness whose dimensions are
   visible is reviewable; one that is not is self-report with better branding — which is how
   "0 silent holes" survived a round it should not have.
2. **(answered this cycle) Round-04 question 3, the scanner strategy.** The owner did
   exactly what was proposed: `SKILL.md` now declares the block model normatively ("The
   block model the scanner claims (§5.4)"). Recorded as answered, with thanks — the model
   declared is the right one; the gap is the token-class guard one layer below it.
3. **(carried) The level/verdict axis split.** An open, honestly-recorded `slop` violation
   against the winner still yields `verified L3` beside VIOLATION, exit 1. My reading
   (consistent with §9, G1's ban-list being a sharing test) was never confirmed by the
   owner; if wrong, `gateRefutations` needs a G1 arm.
4. **(carried) The second-round critique home** — MINOR above: where does §4 rule 4's
   exception live?
