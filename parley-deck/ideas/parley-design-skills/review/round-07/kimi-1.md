---
agent: kimi-1
idea: parley-design-skills
review-round: 07
date: 2026-07-28
reviewed-commit: 076ded5
---

## Summary

❌ BLOCK. Everything cycle 10 claims to have fixed is genuinely fixed, and I verified each
closure with my own commands and headless Chromium before going back to the class question:
my own round-06 CRITICAL, re-run byte-for-byte in all four spellings (VIOLATION, verified
null, exit 1, the literal named **and** the fail-safe firing), codex-1's `<style>` MAJOR
(PASS, exit 0, `--ghost` unreported, true positives intact), and the three enumeration-found
fixes (the bad-string poison, the hex-escape-swallowed newline, the dimension-unit guard),
each matching Chromium's CSSOM exactly. The enumeration is complete against the §4 token
list: I re-derived all 24 output token types and every per-type verdict in the `css.js`
head record holds — **for the main consumer**. The two named OPENs are over-report exactly
as stated (browser-confirmed). The reverted markup remedy was the right call. npm test
237/237; doctrine byte-identical; D-1 accepted a seventh round.

Then codex-1's round-07 probe reproduced under my hands, in two independent spellings. The
enumeration checked token types against **one** consumer of the string token; the scanner
has two. `scanComments`' string branch skips one code point after a backslash where §4.3.7
consumes up to six hex digits plus one optional trailing whitespace — and when that
whitespace is the newline, the comment pre-pass ends the string one line before the browser
and `scanStylesheet` do. A comment delimiter in the region the two readers disagree on is
then misclassified in either direction: a real comment hides inside the pre-pass's longer
string and is read as structure (codex-1's probe — manufactured declaration, phantom rule,
the colour silently discarded), or browser-live declarations are blanked as "comment" (my
second spelling). Both spellings: **PASS, verified L3, exit 0, `unreadable: []`, while
Chromium computes `rgb(255, 0, 0)`**. That is the seventh time a declared closure — this
time "136 constructs, 0 silent holes" and the token class closed by enumeration — does not
survive a probe.

Honest accounting: I read `scanComments` this round (its ident/hash/at-keyword ordering) and
did not check its string branch against §4.3.7; my enumeration verification named token
types against the main consumer only. The class-level check this round exists for has to
name consumers, not just types — codex-1's probe is the demonstration.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive 076ded5 | tar -x -C /tmp/pds-r7`;
HEAD `076ded5b44c3`, branch `parley-design-skills`, worktree clean). Probe runs are copies
of the shipped `sound-run` fixture under `/tmp/pds-probes-r7`; browser checks ran headless
Google Chrome (`--headless=new --dump-dom`, `getComputedStyle` + CSSOM written into the DOM)
on pages under `/tmp/pds-browser-r7`; the round-06 corpus was reused at `/tmp/pds-corpus`.
Nothing in either repo was modified.

- `npm test` (skill repo at 076ded5): **237 passing, 0 failing** — matches the claim
  (233 → 237).
- Baseline `check --level L3` on the unmodified sound-run: PASS, exit 0.
- Doctrine: `git diff f1c123d..076ded5 -- addons/parley-design/` **empty**; `wc -c` = 6,519 +
  25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256 RULES.md` 12-hex =
  `b49ff596451f` = PDS `registry-digest`. Since `aa6b9b3` only `lib/css.js` and
  `checker.test.js` changed (+600/−39).

**My round-06 CRITICAL, re-run byte-for-byte (N1/N3/N4/N6, full `check --level L3`, plus browser):**

- N1 `#url(`, N3 `@url(`, N4 `#\75 rl(`, N6 `@\75 rl(` — each:
  `core:literal-outside-token-layer` names `.a sets color to the colour literal #ff0000`,
  the fail-safe fires (paren imbalance, 2 blocks never closed, text to EOF), G3 refuted,
  **VIOLATION, verified null, exit 1**. **Closed, all four spellings.**
- Passing-side control (`color: var(--color-text-body)`): PASS, verified L3, exit 0.
  Literal control (`color: #ff0000`): exit 1 — the probe file is genuinely scanned.
- Browser premise re-confirmed for N1: Chromium computes `rgb(255, 0, 0)`, CSSOM
  `.a { color: rgb(255, 0, 0); }`.

**codex-1's round-06 MAJOR, re-run:** the `<style>` selector probe → `markupVarUses` returns
only `--color-text-body@3`; end-to-end **PASS, verified L3, exit 0, no findings**.
**Closed.** True positives intact: `style="color: var(--nope)"` → exit 1;
`stopColor="hsl(var(--nope))"` → exit 1 — the latter being one of the shapes the reverted
remedy would have lost.

**The enumeration, checked rather than member-hunted:** I re-derived the CSS Syntax L3 §4
output-token list (24 types) against the `css.js` head record: ident/function →
`identLikeToken`; hash and at-keyword → `hashOrAtToken` (new this cycle); string/bad-string,
url/bad-url, escapes → consumed above the stack; the six bracket tokens → the §5.4.4 model;
number, percentage, dimension, delim, comma, colon, semicolon, whitespace, CDO, CDC → "no
code point this scanner would read as structure", which is correct per spec (the structural
code points are never delims, and a digit sequence never carries a brace). Every verdict I
could probe held **for the main consumer**. My own member probes, scanner and Chromium: hash
with an escape spelling a brace (`#a\7b b`), at-keyword with an escape spelling a paren
(`@a\28 b`), comment-split idents (`u/**/rl(`, `url/**/(`), digit-prefixed hash (`1#url(`),
hash inside url contents (`url(#a}b)`), `50%\31 px`, an unknown at-rule with a nested rule,
CDO/CDC wrapping — **0 silent holes**, over-reports only, each consistent with the two named
OPENs.

**The two named OPENs, browser-checked:** `width: 1\65 5` → scanner reads `1e5`, Chromium
drops the declaration (computed width unchanged) — over-read. `color: #ff0000 bogus` →
scanner reads it, Chromium drops it (`.c { }`, computed black) — over-read; the scanner can
manufacture a finding here, never hide one. Both exactly as stated, neither worse.

**The three enumeration-found fixes, scanner vs Chromium:** the bad-string phantom-rule probe
→ scanner: `.a2::before` empty, no phantom rule, `unreadable: []`; Chromium CSSOM
`.a2::before { }` — agree. The hex escape swallowing a newline → scanner: `content: "xA"`
and the colour read, `unreadable: []`; Chromium `.probe3 { content: "xA"; color: rgb(255,
0, 0); }` — agree. Bad-string `;`-recovery keeps the following declaration in both. The
dimension guard: `1\31 px` left as written plus fail-safe note, the colour still found;
`1\3b 5px` the same.

**False-positive sweep:** all 37 fixture stylesheets (64 blocks, 107 declarations): 0
unreadable. Round-06 corpus: Bootstrap 2,556 blocks / 5,543 declarations / 0 unreadable
(counts identical to round-06); Google Fonts css2, open-props, modern-normalize, Tailwind
theme/preflight/utilities: 0; my 20-block stress file: the one known escaped-quote note
(carried MINOR).

**codex-1's round-07 CRITICAL, re-run — REPRODUCED, in two spellings:**

- His probe, byte-for-byte (`.probe {` / `content: "x\41`⏎`"; /* x: 1; } */` / `color:
  #ff0000;` / `dummy: y { z: 1;` / `}`): scanner → `.probe` carrying `content="xA"` and a
  manufactured `/* x=1`, plus a phantom `dummy: y` block; `color: #ff0000` discarded;
  `unreadable: []`; `stripComments` output equals the input verbatim — the real comment was
  never blanked. End-to-end: **PASS, verified L3, exit 0, no findings**. Chromium:
  **computes `rgb(255, 0, 0)`**, CSSOM `.probe { content: "xA"; color: rgb(255, 0, 0); }`.
- My second spelling V5 (`.probe { content: "x\41`⏎`/*"; color: #ff0000; dummy: "*/"; }`):
  the pre-pass blanks the browser-live `/*"; color: #ff0000; dummy: "*/` as a comment; the
  scanner reads one content declaration, `unreadable: []`; end-to-end **PASS, verified L3,
  exit 0**; Chromium **computes `rgb(255, 0, 0)`** (`.probe { content: "xA/*"; color:
  rgb(255, 0, 0); }`).
- Control V3 (hex escape swallowing a *space* beside a real comment): read correctly,
  comment stripped, colour found, `unreadable: []` — the desync needs the escape's optional
  whitespace to be the newline.
- The other two backslash-skipping consumers, `parenBalance` and `splitDeclaration`, share
  the one-code-point skip but their string models err *longer* (no bad-string rule) and feed
  only notes and colon-splitting — over-report direction, not a silent-hole vector.

**Carried MINORs re-probed at 076ded5 (all unchanged):** `fill: url(#fade)` still a false
VIOLATION (exit 1); `VAR(--x)` and `RGB(255, 0, 0)` still PASS clean (exit 0) while Chromium
substitutes and applies; `content: "say \"hi\""` still whole-file unreadable (the single
note in my stress file); `@position-fallback`/`@try` and `@function` still false unreadable
(exit 4; the at-rule lists are unchanged); the `WAVERIS.md` typo and the second-critique
home are engine-side and the engine is unchanged since `aa6b9b3` — carried.

**Lens (PDS.md conformance):** doctrine byte-identical to `f1c123d`, so my round-04
line-by-line verification carries; spot re-checks this round: §0–§12 all present, the C1
mapping table verbatim, §3 G1 with the A1 two-axis test plus ban list, category-plus-
avoidance and ratification conjuncts and its canonical error string, §4 rule 7 (A3's
unattended ABSTAIN), all eight artifact kinds as H3 entries in the identical four-part
shape (DIRECTION re-read: purpose line → rationale → required-fields table → example).
RULES.md: 19 rule H3s; digest matches. hermes-1's D-2 undercount verified:
`core:contrast-applied` (RULES.md:95) has no detector and is not among D-2's named eight —
reported UNJUDGEABLE at runtime (safe direction), missing from the disclosure text (NIT).

## Findings

### [CRITICAL] (codex-1, round-07 — independently reproduced, second spelling added) `scanComments` ends a string at the newline a hex escape just swallowed; the two string readers desync over a comment delimiter and the run certifies clean L3 over a colour Chromium applies

**What.** `scanComments`' string branch (`lib/css.js:395-408` in the archived tree) handles a
backslash by skipping exactly one code point (`if (text[cursor] === "\\") cursor += 1;`).
§4.3.7 consumes up to six hex digits and one optional trailing whitespace — which may be the
newline that would otherwise end the string. The main loop's quote branch and `stringToken`
both do this right (`consumeEscape`); the comment pre-pass does not. So when a hex escape's
trailing whitespace **is** the newline, `scanComments` ends the string one line before the
browser and `scanStylesheet` do, and the two readers disagree about what is string content
from there on. A comment delimiter in the disagreed region is misclassified, in either
direction, and both directions certify clean:

- **Spelling A** (codex-1's probe): the real comment hides inside the pre-pass's *longer*
  string and is never blanked; the stack machine reads the comment's contents as structure —
  a manufactured `/* x: 1` declaration, the `}` closing `.probe` early, `color: #ff0000`
  discarded at top level with no note (the round-05/06 silent discard channel, still open),
  a phantom `dummy: y` rule — balanced, so `unreadable: []`, PASS, verified L3, exit 0,
  while Chromium computes `rgb(255, 0, 0)`.
- **Spelling B** (mine, V5): a `/*` inside the browser's string is a comment opener to the
  pre-pass, which blanks the browser-live `"; color: #ff0000; dummy: "` through the next
  `*/`. The colour is deleted before the stack machine runs; the file ends balanced;
  `unreadable: []`, PASS, verified L3, exit 0, while Chromium computes `rgb(255, 0, 0)`.

**Why it matters.** The exact failure signature the fail-safe exists to eliminate, reachable
from ordinary-if-unusual CSS — one escaped character inside a string, a newline, a comment.
It refutes the cycle-10 load-bearing claims at the reviewed commit: "136 constructs … 0
silent holes" and "the token class, closed by enumeration". The enumeration's string verdict
("each is consumed above, before the stack") is true of `stringToken` and the main loop and
false of the second consumer. It is also the second consecutive round in which the
differential matrix did not contain the construct class that broke the claim, and the matrix
is still committed nowhere. And spelling A's payload travels the silent top-level `;`-flush
channel named in my round-05 and round-06 counter-proposals: defence-in-depth item (2) was
never implemented, and it would have converted spelling A into a fail-safe firing.

**Fix.**

1. Make `scanComments` consume strings and escapes through the same
   `stringToken`/`consumeEscape` path as `scanStylesheet` — factor the token walk so the two
   passes share one definition of a string and cannot diverge again (codex-1's remedy;
   correct and minimal).
2. Defence in depth, carried from round-05/06: report a top-level `;`-flush (or EOF flush)
   whose buffer is non-empty and does not start with `@` through the unreadable channel, as
   in-rule discards already are. This catches spelling A's shape even against the next
   consumer desync.
3. Fixtures: codex-1's probe and V5, each of which must fail or go UNJUDGEABLE, plus
   passing-side controls — a hex escape swallowing a space beside a real comment (the V3
   shape), an ordinary comment after a string, a bad-string beside a real comment, and the
   `\`⏎ line continuation.
4. Extend the enumeration record from per-token-type to per-token-type × per-consumer
   (`scanComments`, the main loop, `decodeDeclarationText`), and commit the differential
   harness or its construct matrix — this is the second "0 silent holes" claim the
   uncommitted matrix did not cover.

### [MINOR] (carried, sharpened; also filed by codex-1) The verdict-table pointer dangles, and the enumeration must name consumers, not just token types

**What.** IMPLEMENTATION.md's cycle-10 record says the per-token-type verdict table "lives in
the cycle-10 agent report"; the commit message says it is "in the fix-up record". Neither
location contains it — no agent report exists in either repo (both searched), and the
fix-up record holds summary prose. The committed enumeration is the `css.js` head comment,
which is complete per token type but silent on **which consumer** each verdict covers.
**Why it matters.** This round existed to check the enumeration. The artifact named for that
check is not where two different pointers say it is, and the artifact that does exist frames
coverage per token type — the framing that let the string verdict pass although only one of
the two string consumers was checked. The 136-construct harness is likewise still
uncommitted, so "0 silent holes" remains irreproducible by a reviewer.
**Fix.** Commit the verdict table with rows per §4 token type × per consuming pass, with the
fixture and disposition for each cell, and commit the differential harness or its construct
matrix beside it; correct the pointer in IMPLEMENTATION.md.

### [MINOR] (carried) Case is a spelling the decoded layer does not cover: `VAR(--x)`, `RGB(255,0,0)` pass clean while the browser substitutes and applies

Re-probed at 076ded5: PASS, exit 0 both. **Fix** (unchanged): the `i` flag on the
function-name and unit alternatives, and a case-insensitive `var(` match on the decoded
value. Fixtures: the two probes (flagged), the lowercase forms (unchanged).

### [MINOR] (carried) An escaped quote or backslash inside a string still makes the whole file unreadable

Re-probed: `content: "say \"hi\""` is still the single note in my 20-block stress file
(exit 4). Cycle 10 made the newline boundary precise but kept the quote/backslash routing.
**Fix** (unchanged): keep the token verbatim but reserve the note for escapes spelling
value-structure code points (comma, braces, parens).

### [MINOR] (carried) The literal matchers still scan `url()` contents — `fill: url(#fade)` is a false VIOLATION

Re-probed at 076ded5: exit 1, `#fade` named as a colour literal. **Fix** (unchanged): mask
`url(…)` spans and string tokens before literal matching — the scanner already knows both
boundaries. Fixtures: the probe (clean), `fill: #fade` outside a url (still flagged).

### [MINOR] (carried) `@position-fallback`/`@try` and `@function` still hit the unclassifiable at-rule guard — false unreadable on shipped CSS

Re-probed at 076ded5: exit 4 both; the at-rule lists are unchanged. **Fix** (unchanged): add
`try` and `function` to `DECLARATION_AT_RULES` and `position-fallback` to `RULE_AT_RULES`;
fixtures both probes (no unreadable, declarations visible) and an honestly unclassifiable
at-rule (still guarded, exit 4).

### [MINOR] (carried) A contract whose named waiver file resolves to no readable file is still silently treated as "no waivers"

Engine unchanged since `aa6b9b3`; carried. **Fix** (unchanged): report it — NEEDS_REVIEW or
an `l2-process-order` finding naming the contract and path.

### [MINOR] (carried) §4 rule 4's second critique round still has no §1 home

Engine and doctrine unchanged; carried. **Fix** (unchanged): §1 names the second round's
home (and `PROCESS_HOMES` accepts it when the instruction is recorded), or §4 rule 4 states
the second round also files under `round-02`. Fixture for the permitted shape.

### [NIT] (new) CDO/CDC at stylesheet top level leaves "text runs to end of file" — false unreadable on legacy comment-wrapped CSS

`.a { color: #ff0000; } -->` and `<!-- .a { color: #ff0000; } -->` read the colour correctly
(over-report at worst: the selector comes out `<!-- .a`) but end with an unreadable note and
exit 4, while Chromium drops CDO/CDC at the top level (browser-confirmed: rule present,
colour applied). Legacy idiom, and the 521-file sweep saw none. **Fix:** consume `<!--` and
`-->` as whitespace at rule depth zero, matching the browser's top-level handling.

### [NIT] (carried, sharpened) The checker SKILL.md's scanner section still says "seven constructs … each found by a probe"

`parley-design-check/SKILL.md:88`. The list now predates the hash/at-keyword token (an
eighth), the bad-string phantom, the hex-escape newline and the dimension guard — and "each
found by a probe" now undersells the enumeration, which is the actual assurance argument and
is not mentioned. When the CRITICAL lands, rewrite the section around the enumeration (per
type × per consumer) rather than lengthening the list a ninth time.

### [NIT] (carried) The "will not take on trust" table still lists five conditions

`parley-design-check/SKILL.md:182`. Cycle 4–6's recomputed conditions are still not rows.
Add the rows or retitle the table.

### [NIT] (carried; hermes-1's observation, verified) D-2's disclosure names eight undetected rules; `core:contrast-applied` is a ninth

Verified: RULES.md:95, no detector module, reported UNJUDGEABLE at runtime (safe direction),
absent from D-2's named list in IMPLEMENTATION.md. One-line correction to the disclosure.

### [NIT] (carried) PDS.md has 6 bytes of headroom against its early-warning threshold

25,594 of 25,600 (total 65,360/65,536, 176 spare). Unchanged; recorded so the next normative
sentence is planned for.

## Open questions

1. **(carried, sharpened — now load-bearing twice) The differential harness.** Round-06
   asked for it; cycle 10 answered with a pointer to a "cycle-10 agent report" that exists
   in neither repo, and the 136-construct matrix is still uncommitted. The construct class
   that refuted this round's claim was absent from the matrix, exactly as the prefix class
   was absent from the 86. Will the owner commit the harness or its matrix, per token type
   × per consumer?
2. **(carried) The level/verdict axis split.** An open, honestly-recorded `slop` violation
   against the winner still yields `verified L3` beside VIOLATION, exit 1. Engine unchanged;
   never confirmed by the owner.
3. **(answered) Round-06's class question.** The token-class enumeration is genuine and
   committed at the head of `css.js`, and it is complete per token type for the main
   consumer; what it needed was the consumer dimension. Recorded via the CRITICAL.
