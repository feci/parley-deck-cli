---
agent: kimi-1
idea: parley-design-skills
review-round: 08
date: 2026-07-28
reviewed-commit: 82bde7d
---

## Summary

❌ BLOCK. The round-07 CRITICAL is closed structurally, not locally, and I verified the
structure before the symptom: exactly one string consumer (`stringToken`, `css.js:384`) and
one escape consumer (`consumeEscape`, `css.js:214`), reached by all six passes — I enumerated
every quote and backslash site in `lib/` rather than taking the head comment's word for it.
My round-07 reproducer, re-run byte-for-byte in both spellings plus the control, now agrees
with headless Chromium and exits 1 with the literal named. The corrected enumeration —
token-type × consumer — I tested as a claim, one desync probe per consumer pair, each checked
against Chromium's CSSOM: the scanner and the browser agreed in every case, down to the exact
decoded spellings (`"xA:y"`, `"xA)"`, `"Arial"`).

Five of the six false-positive closures hold in both directions. The sixth does not, and it
is the one this round exists to catch: the `@function` classification made the checker more
permissive, and its parameter exemption opens a false clean. `functionParameters` captures
every `--name` in the prelude — including names inside parameter **default expressions** — so
a body's `var(--ghost)` is suppressed as "bound" while the default's own `var(--ghost)` is
never collected. codex-1 found it; I reproduced it independently at the CLI and in Chrome
150, in both halves, with the controls pointing the right way. The checker certifies PASS,
verified L3, exit 0 over a token reference Chromium resolves — the identical signature to the
round-05/06/07 CRITICALs, this time introduced by this cycle's own fix.

Everything else claimed for cycle 11 held: npm test 244/244; doctrine byte-identical, budget
65,360/65,536, digest `b49ff596451f`; Bootstrap differential exactly as recorded (2,556/5,543/
0, and the uppercase-RGBA blindness was real — 62 occurrences in the source). D-1 accepted an
eighth round. One process note: a kimi-1 signoff block for this round already existed in the
review `consensus.md` when I finished — a prior kimi-1 session's, whose review file was never
written. This file is the missing evidence; my signoff supersedes that block with the same
verdict it reached.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive 82bde7d | tar -x -C /tmp/pds-r8`;
reviewed commit is HEAD of `parley-design-skills`, worktree clean). Probe runs are copies of
the shipped `conformance/sound-run` fixture under `/tmp/pds-probes-r8`; browser checks ran
headless Google Chrome **150.0.7871.187** (`--headless=new --dump-dom`, `getComputedStyle` +
CSSOM serialised into the DOM) under `/tmp/pds-browser-r8`; the round-06 corpus survived at
`/tmp/pds-corpus`. Nothing in either repo was modified. My first pass over the checker (the
round-07 reproducer, the consumer matrix, and the six closures) was run before I read
codex-1's round-08 review; the `@function` default probes were added after, and each was
re-derived with my own commands rather than taken from his file.

- `npm test` (pristine extraction, Node v26.5.0): **244 passing, 0 failing** — matches the
  claimed 237 → 244.
- Baseline `check --level L3` on the unmodified sound-run: PASS, verified L3, exit 0.
- Doctrine: `git diff 076ded5..82bde7d -- addons/parley-design/` **empty**; `wc -c` = 6,519 +
  25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256 RULES.md` 12-hex =
  `b49ff596451f` = PDS `registry-digest`. Cycle 11 changed only the checker.

**The structural claim, checked structurally.** Every quote site in `lib/`: `css.js:312`
(the quoted/unquoted url dispatch — a branch, not a consumer), `:389` (inside `stringToken`
itself), and six call sites that each delegate (`:474/488` scanComments, `:515/516`
parenBalance, `:592/593` decodeDeclarationText, `:669/670` splitDeclaration, `:904/906`
scanStylesheet, `:1148/1149` maskOpaqueTokens). Every backslash site goes through
`validEscape`/`consumeEscape`. `identLikeToken`, `hashOrAtToken`, `unquotedUrl` are each
defined once. `registry.js`'s own quote handling parses YAML, not CSS — out of scope. The
"one consumer, six passes" claim is true in the code, not just in the head comment.

**My round-07 CRITICAL, re-run byte-for-byte (checker + browser):**

- Spelling A (codex-1's probe): `core:literal-outside-token-layer` names `.probe sets color to
  the colour literal #ff0000` at line 17; the fail-safe also fires on the residue (`1 block
  opened and never closed`), which under the correct reading is what the file does — G3
  refuted, **VIOLATION, level not verified, exit 1**. Chromium: computes `rgb(255, 0, 0)`,
  CSSOM `.probe { content: "xA"; color: rgb(255, 0, 0); }`. **Closed.**
- Spelling B (my V5): the colour named at line 15, **no unreadable** (unjudgeable count back
  at the baseline 13), **VIOLATION, exit 1**. Chromium: `.probe { content: "xA/*"; color:
  rgb(255, 0, 0); }` — exact agreement. **Closed.**
- Control V3 (hex escape swallowing a *space* beside a real comment): colour found, no
  unreadable, exit 1 — unchanged, the regression control holds.
- The committed test `one string consumer` encodes both probes by name (`kimi-1's spelling,
  the delimiter inside the string`), both blanking directions, and passing-side controls.

**The enumeration tested as a claim — one desync probe per consumer, scanner vs Chromium:**

- `splitDeclaration`: `content: "x\41`⏎`:y"` → prop `content`, value `"xA:y"`, the in-string
  colon divides nothing; Chromium CSSOM `content: "xA:y"`. Colour found, exit 1.
- `parenBalance`: `content: "x\41`⏎`)"` → no imbalance note; Chromium `content: "xA)"`.
- `decodeDeclarationText`: `font-family: "\41`⏎`rial"` → scanner spells `"Arial"`; Chromium
  computes `font-family: Arial`. The sound-run allowlist is `[Chartwell Text, Chartwell
  Mono]`, so the resulting `core:face-outside-allowlist` finding is a true positive against
  the decoded spelling.
- `maskOpaqueTokens`, string side: `background: "\41`⏎`#ff0000"` → contents blanked, no
  finding, PASS exit 0; Chromium drops the declaration (CSSOM `.probe { }`). Two readings
  agree: nothing applied, nothing reported.
- `maskOpaqueTokens`, url side: `background: url(#fade) #ff0000` → mask `url(     ) #ff0000`,
  the finding names `#ff0000` only; Chromium computes `background: url("#fade") rgb(255, 0,
  0)`. Uppercase `URL(` and escaped `\75 rl(` are masked identically, and the literal beside
  the escaped spelling is still found.

**The six closures, both directions each — five hold:**

1. Case: `RGB(255, 0, 0)`, `RGBA(255, 0, 0, 0.5)`, `font-size: 11PX` → VIOLATION, exit 1
   (round-07: PASS); Chromium applies the first two (CSSOM confirmed). `VAR(--color-text-body)`
   declared → PASS; `VAR(--nope)` → `token-used-undeclared`, exit 1; markup
   `w-[VAR(--color-text-body)]` → PASS where `w-[VAR(--ghost)]` → exit 1.
2. Escaped quote/backslash: `content: "say \"hi\""` and `content: "a\\b"` beside a literal →
   exit 1 on the literal, **no unreadable** (round-07: whole-file unreadable, exit 4). The
   boundary held: `content: "\2c x"` still routes to the fail-safe *and* the literal is found.
3. url contents: `fill: url(#fade)` → PASS exit 0 (round-07's false VIOLATION closed); the
   outside literal still found (above).
4. CDO/CDC: `.probe { color: #ff0000; } -->` and `<!-- .probe { ... } -->` → literal found,
   exit 1, no unreadable. The other direction: `.probe { <!-- color: #ff0000; }` → PASS, no
   finding, and Chromium CSSOM is `.probe { }` — scanner and browser discard the same
   declaration. No false clean.
5. Waiver file: sound-run with `WAIVERS.md` deleted → `pds-check:l2-process-order` names the
   contract and path, VIOLATION, exit 1, **both at `--level L3` and at plain `check`**;
   control with the file present: no finding; the `waivers: ""` empty-field path (cycle-4
   H-1) is preserved.
6. `@function` / `@position-fallback` / `@try`: no unreadable; `@try { color: #ff0000 }` →
   literal found, exit 1; `@function --double(--x) { result: calc(var(--x) * 2) }` → PASS
   (true formal exempt); `var(--y)` absent from the prelude → flagged, exit 1. **But the
   parameter set itself is built wrong — see the CRITICAL.** Browser note: Chromium drops
   `color` inside `@try` (not a position-try descriptor), so that finding is an over-report —
   the declared T1 boundary direction, not a false clean.

**The named over-reads, re-verified:** `width: 1\65 5` → PASS (ungoverned prop); Chromium
drops the declaration — unchanged. `color: #ff0000 bogus` → VIOLATION, exit 1 — scanner
reads, browser drops, the declared T1 boundary, unchanged.

**Sweeps:** all 37 fixture stylesheets: 0 unreadable. Corpus: Bootstrap 2,556 blocks / 5,543
declarations / 0 unreadable (**counts identical to rounds 06 and 07** — the one-consumer
rewrite did not shift the reading of a 274 KB real stylesheet); gfonts, normalize, openprops,
tw-theme/preflight/utilities: 0; my 20-block stress file: 0 unreadable, `content: "say
\"hi\""` read back verbatim — round-07's one carried note is gone. Bootstrap's source carries
**62** uppercase `RGB(`/`RGBA(` occurrences: the "44 gained, all true positives" claim rests
on real text.

**Record-keeping closures:** the checker `SKILL.md` scanner section is rewritten around the
two-axis enumeration; the "will not take on trust" table gained the rows; IMPLEMENTATION.md's
D-2 now names nine undetected rules including `core:contrast-applied` (hermes-1's undercount
closed).

**Lens (PDS conformance):** doctrine byte-identical to `f1c123d`, so my round-04 line-by-line
verification carries. Spot re-checks: §0–§12 all present; all eight artifact kinds as H3
entries; §5 rule 7 states the `DESIGN-SYSTEM.md` authorship rule (the Phase-6 design
reviewer) — FINAL.md's open item is stated normatively, not left implicit.

**The round-08 finding, independently reproduced (codex-1's probe, my own commands):**

- codex-1's probe byte-for-byte (`:root { --ghost: rgb(255,0,0); }` / `@function
  --pick(--x: var(--ghost)) { result: var(--ghost); }` / `.probe { color: --pick(); }`):
  checker **PASS, verified L3, exit 0, no findings, unreadable []**; `varUses` over the file
  returns only the sound-run's five legitimate uses — `var(--ghost)` appears nowhere. Chrome
  150: parses a real `CSSFunctionRule`, serialises the whole construct in CSSOM, computes
  `rgb(255, 0, 0)`. Node-level: the at-block prelude is `"@function --pick(--x:
  var(--ghost))"`, and `functionParameters`'s regex over everything after the first `(` yields
  `{--x, --ghost}` — the default's reference is captured as a formal.
- Second half (mine): `@function --pick(--x: var(--ghost)) { result: 1px; }` — the body
  references nothing; the default's own `var(--ghost)` is still invisible: **exit 0**. The
  prelude is never scanned for uses at all.
- Controls: plain `.probe { color: var(--ghost); }` outside any `@function` →
  `token-used-undeclared`, **exit 1** (the rule works when the reference is not laundered
  through a prelude); `@function --pick(--x: var(--color-text-body)) { result: var(--x); }`
  (true formal, declared default) → **exit 0**, the clean shape a fix must preserve.

## Findings

### [CRITICAL] (codex-1, round-08 — independently reproduced, second half added) `@function` parameter defaults launder real token references into a false clean

**What.** `functionParameters` (`css.js:1207-1213`) builds the bound-parameter set with
`at.prelude.slice(open + 1).match(/--[A-Za-z0-9_-]+/g)` — every custom-property-shaped name
after the first `(`, not every formal. For `@function --pick(--x: var(--ghost))` the set is
`{--x, --ghost}`, although only `--x` is a formal parameter. Two defects follow:

1. `declarationVarUses` suppresses the body's `var(--ghost)` as bound — but `--ghost` is no
   parameter; a caller omitting the argument resolves the **default**, and the body reference
   resolves against the token layer exactly as any other `var()` does.
2. The default expression itself is never scanned: preludes are not declarations, so the
   `var(--ghost)` inside `--x: var(--ghost)` is collected by nobody.

Either half alone hides the reference; together they erase it. The sound-run's ratified token
layer declares `--color-text-body`, `--color-surface-page`, `--color-rule`,
`--space-gap-md`, `--space-gap-lg` — no `--ghost`. The checker reports **PASS, verified L3,
exit 0, unreadable: [], no findings** while Chrome 150 parses the rule, answers the function
call, and computes `rgb(255, 0, 0)` through the undeclared token. The same probe with the
`var(--ghost)` written plainly in `.probe` exits 1 — the rule works until the reference is
laundered through a prelude.

**Why it matters.** The exact signature of the round-05/06/07 CRITICALs: a run certifies
clean L3 over output a browser resolves, through a path this cycle's own fix opened — the
`@function` classification and its parameter exemption are new in cycle 11, and the exemption
is the hole. `core:token-used-undeclared` is a `system`-class rule whose openness G3 and the
L3 level certify against, so the false clean is a false certificate, not a missed nicety.
codex-1 filed this as MAJOR; I file it CRITICAL because the signature — PASS + verified L3 +
exit 0 over browser-resolved output a shipped rule exists to catch — is the one this review
has rated CRITICAL in every prior round. The evidence is identical either way; only the label
differs.

**Fix.**

1. Parse the `@function` parameter list at top-level commas, respecting strings and nested
   parentheses (the file already owns both walkers), and take only each segment's leading
   custom-property name as a formal.
2. Collect `var()` references from default expressions as real uses — a caller omitting the
   argument resolves the default against the token layer, which is what Chromium did.
3. Suppress only true formals while scanning the body.
4. Fixtures: codex-1's probe (must flag `--ghost`); the default-only half above (must flag);
   the clean control (true formal + declared default, must stay clean); a body reference
   absent from the prelude (still flags — the existing shape); and the probe's browser
   premise pinned in the fixture comment.

### [MINOR] (carried, fourth round) §4 rule 4's second critique round still has no §1 home

**What.** PDS §4 rule 4 permits a second critique round on an explicit Decider instruction,
but §1's mapping names only `round-02` for CRITIQUE, and the checker's `PROCESS_HOMES`
accepts only that home — so a run that exercises the permission the spec grants fails the
L2 process-order check the same spec defines. Doctrine byte-identical since `f1c123d`; the
engine's home list unchanged. **Why it matters.** A permitted shape the conformance level
rejects is a trap for the first real run that uses it. One sentence, but on the normative
path my lens covers, and unanswered through four cycles of checker work. **Fix.** §1 names
the second round's home (and `PROCESS_HOMES` accepts it when the instruction is recorded), or
§4 rule 4 states the second round also files under `round-02`; one fixture for the permitted
shape.

### [NIT] (new) The record writes counts in prose, and the two counts now disagree

**What.** `css.js:11` still says "Seven constructs of one family have now been found by
probe"; the checker `SKILL.md` rewritten this cycle says "Eight constructs broke it one at a
time". The eighth is the string-consumer desync itself, described three comment blocks below
the stale seven. **Why it matters.** "Never write a count in prose" is one of the nine
defects FINAL.md says not to reproduce — counts in prose drift, and these two now have.
**Fix.** Align the number, or drop it from both and point at the enumeration, which is the
artifact that is actually maintained.

### [NIT] (new) `@try` declarations are scanner-visible and browser-dropped — worth one line beside the named over-reads

**What.** The new at-rule classification reads `@try` bodies as declarations, so `color:
#ff0000` inside one is reported as a literal; Chromium drops it (only position-try
descriptors are valid there). Over-report, the declared T1 direction — filed only because the
two named over-reads are called out by name in the code and this third one is not, so the
next reader of a finding against `@try` has no signpost. **Fix.** One sentence beside the
named OPENs noting that declaration-holding at-rules are read as declared, not as applied,
and that this is the tier boundary, not a defect.

## Open questions

1. **(carried, softened) The differential matrix.** Cycle 11's "deliberately not done" is a
   real reason for the harness — a Chromium-dependent tool does not belong in a Node-built-ins
   package — and the load-bearing constructs are now committed as test fixtures, so this is no
   longer load-bearing the way it was at rounds 6–7. What remains self-reported are the sweep
   numbers (170 constructs 3→0; 283 stylesheets / 2,288 markup files). I verified by
   independent sampling instead: my own probe set against headless Chromium, plus the
   surviving corpus, whose block/declaration counts are identical to the two prior rounds.
   Would the owner commit the construct matrix — the strings and their expected dispositions,
   which carry no Chromium dependency — as a fixture file?
2. **(carried) The level/verdict axis split.** An open, honestly-recorded `slop` violation
   against the winner still yields `verified L3` beside VIOLATION, exit 1. Engine unchanged
   this cycle; never confirmed by the owner as intended.
3. **(answered) Round-07's class question.** The enumeration now names consumers, the code
   matches the table, and the claim survived one desync probe per consumer under my hands,
   browser-checked. Round-08's lesson is that the enumeration's discipline has to reach every
   *new* consumer a fix adds: the `@function` exemption was built this cycle with a regex
   where the rest of the file uses tokenisers, and it is where the false clean came from.
