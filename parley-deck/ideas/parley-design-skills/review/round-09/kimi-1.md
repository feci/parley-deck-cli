---
agent: kimi-1
idea: parley-design-skills
review-round: 09
date: 2026-07-28
reviewed-commit: 9121ec2
---

## Summary

✅ ACCEPT. The round-08 CRITICAL is genuinely closed, and the closure did not swing back
toward a false alarm — the two failure directions the last two cycles have been trading.
I re-ran my own round-08 reproducer byte-for-byte against a pristine extraction of the
reviewed commit: codex-1's probe now flags `--ghost` twice (the default's own reference
and the body's, both named at the `@function` line), my default-only second half flags the
default's reference, and both exit 1 with G3 refuted and `verified=null`. All four controls
point the right way: the declared-default-plus-true-formal clean control stays PASS,
verified L3, exit 0; the true-formal probe flags only its planted literal; the body
reference absent from the prelude still flags; the plain `var(--ghost)` control still
flags. The fix is the one both reviewers filed: the parameter list is parsed at top-level
commas through the shared `stringToken`/`consumeEscape` consumers, only a segment's leading
custom-property identifier binds, defaults are decoded and their `var()` references
collected as real uses, and an undividable list goes to the unreadable channel binding
nothing — verified in the code, not just the commit message.

The false-alarm direction is where I spent the new probes, because that is what this fix
could have broken: Sass `@function px-to-rem($n)` and `spacing($base, $mult: 2)` behave
byte-identically to the round-08 build (the `@return` unreadable is pre-existing cycle-11
behaviour, not this cycle's channel), a Sass function with a declaration-shaped body is
fully clean, a comma inside a string default divides nothing, a comma inside `min(1em,
2em)` divides nothing, an escaped `\76 ar(--ghost)` in a default is decoded and flagged,
an escaped `\(` inside the function name is not read as the list opener, `type(<length>)`
annotations cost no use and no binding, and the empty list `--now()` is not a defect. The
unreadable channel itself fails safe in all three shapes I probed (non-`--` leading token,
bad later segment discarding the whole binding, unclosed list) — exit 4 UNJUDGEABLE or
exit 1 VIOLATION, never clean.

npm test 245/245. Baseline sound-run: PASS, exit 0, the same 13 baseline unjudgeables.
Cycle 12 touches only `lib/css.js` and `checker.test.js`; doctrine byte-identical,
65,360/65,536, digest `b49ff596451f` — D-1 accepted a ninth round. Direct r8-vs-r9 scanner
diff over my whole corpus and all 37 fixture stylesheets: identical output, 0 new, 0 lost,
0 newly unreadable — the claimed sweep, verified on my own sample. No new findings.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive 9121ec2 | tar -x -C /tmp/pds-r9`;
reviewed commit is HEAD of `parley-design-skills`, worktree clean; cycle 12 is exactly one
commit over round-08's `82bde7d`). Probe runs are fresh copies of the shipped
`conformance/sound-run` fixture (verified `diff -r`-identical to round-08's) under
`/tmp/pds-probes-r9`; the round-06 corpus survived at `/tmp/pds-corpus`; round-08's build
survived at `/tmp/pds-r8` for direct comparison. Nothing in either repo was modified.

- `npm test` (pristine extraction, Node v26.5.0): **245 passing, 0 failing** — matches the
  claimed 244 → 245.
- `git diff 82bde7d..9121ec2 --stat`: only `addons/parley-design-check/lib/css.js` (+151)
  and `test/checker.test.js` (+163). `git diff ... -- addons/parley-design/` empty;
  `wc -c` = 6,519 + 25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256
  RULES.md` 12-hex = `b49ff596451f` = PDS `registry-digest`.
- Baseline `check --level L3` on the unmodified sound-run: **PASS, exit 0**, 13
  unjudgeables — the round-08 baseline, undisturbed.

**The code, read before the probes.** `functionBinding` (`css.js:772-885`): reads the
function name with `consumeIdentSequence` and binds only when it starts with `--` (a
plain-ident prelude — Sass's decade of `@function px-to-rem($n)` — is not this at-rule and
says nothing); requires `(` immediately after the name, so an escaped `\(` inside the name
is not the opener; walks to the matching `)` with a block stack, strings and escapes going
through the shared `stringToken`/`validEscape`/`consumeEscape`; divides only at commas with
stack depth 1; an unclosed list or a segment whose leading token is not a custom-property
name goes to the unreadable log and binds nothing (the whole binding is discarded, so no
guessed suppression survives); everything after each segment's leading name is decoded via
`decodeDeclarationText` and its `var()` references collected as uses at the `@function`'s
line. `declarationVarUses` consults `block.atBlock.binding`, suppressing only true formals.
The deleted `functionParameters` regex is gone entirely — no second implementation survives.

**My round-08 CRITICAL, re-run byte-for-byte (fresh sound-run copies, `--level L3 --json`):**

- codex-1's probe (`:root { --ghost: rgb(255,0,0); }` / `@function --pick(--x:
  var(--ghost)) { result: var(--ghost); }` / `.probe { color: --pick(); }`): **VIOLATION,
  verified=null, exit 1** — `core:token-used-undeclared` names `--ghost` at probe.css:2
  **twice**: the default's own reference and the body's, exactly the committed test's
  `["--ghost@2", "--ghost@2"]`. G3 refuted by the run's own findings. Round-08: PASS,
  verified L3, exit 0. **Closed.**
- My default-only half (`@function --pick(--x: var(--ghost)) { result: 1px; }`): **exit 1**,
  `--ghost` named at line 1 — the prelude default is now scanned for uses. **Closed.**
- Clean control (`--x: var(--color-text-body)`, body `var(--x)`): **PASS, verified L3,
  exit 0** — the shape the fix had to preserve, preserved.
- True formal (`--double(--x)` body `calc(var(--x) * 2)` beside a planted `#ff0000`): exit 1
  on the literal only; `var(--x)` suppressed. The exemption still works.
- Undeclared body (`var(--nope)` absent from the prelude): exit 1, named. Held.
- Plain control (`.probe { color: var(--ghost); }`): exit 1, named. Held.

**The false-alarm direction (new surface, my own probes):**

- Sass, three shapes: `px-to-rem($n)` / `spacing($base, $mult: 2)` → exit 4 UNJUDGEABLE —
  but the reason is the pre-existing "text inside a rule that is not a declaration" channel
  on `@return $n;`, and both probes produce the **identical** result at the round-08 build
  (`82bde7d`: exit 4 both) — cycle 12 added nothing here. A Sass function with a
  declaration-shaped body (`px-to-rem($n) { result: $n; }`): **PASS, exit 0** — the new
  parameter channel stays silent on plain-ident names, the mid-flight false alarm the
  commit message describes is genuinely closed.
- Undividable lists, all three shapes: `@function --pick(1px)` → exit 4 UNJUDGEABLE with the
  new reason verbatim ("carries \"1px\", which this scanner cannot read as a parameter");
  `@function --f(--x, 1px)` → the whole binding discarded, body's `var(--x)` flagged **and**
  the unreadable note fires, exit 1 — no guessed suppression in either direction; unclosed
  list → the pre-existing block fail-safe, exit 4. Never clean.
- Structure inside the list: `--x: "a,b", --y: var(--ghost)` → exactly one finding,
  `--ghost` (the string comma divided nothing, the formal `--y` suppressed, the default's
  reference collected); `--x: min(1em, 2em), --y: var(--color-text-body)` → PASS, exit 0;
  `--x: \76 ar(--ghost)` → exit 1, `--ghost` named (defaults are decoded before scanning);
  `--pi\(ck(--x)` → PASS, exit 0 (escaped opener inside the name is not the list opener,
  `--x` bound and suppressed); `--x type(<length>): 1px` → PASS, exit 0; two declared
  defaults → PASS, exit 0.
- Exotic nesting (`@media` inside the `@function` body): the scanner attaches the bare
  declaration to the `@function` block — **byte-identical at `82bde7d`**, so pre-existing —
  the binding reaches it (`var(--x)` suppressed, exit 0), and a nested `var(--ghost)` is
  still collected as a use (checked at the scanner: `[{"name":"--ghost","line":1}]`). No
  hole, nothing new.

**The committed fixture** (`checker.test.js`, "an @function binds its parameters and uses
its defaults, and a default is never a binding"): encodes codex-1's probe byte-for-byte
(`FUNCTION_DEFAULT_REPRODUCER` is identical to my round-08 batch file), both clean controls,
the comma-division matrix (nested function, bracket block, string, escaped comma, `returns`
clause, css-type), the unreadable-binds-nothing direction, the Sass false-alarm side, and
the end-to-end assertion naming the round-08 signature: "What this fixture must not produce
is the clean pass it produced at `82bde7d`", with `verdict === "VIOLATION"`,
`verified !== "L3"`, `exit !== 0` pinned. My filed remedy, verbatim, with passing-side
controls.

**Sweeps:** direct r8-vs-r9 diff of `scanStylesheet` output over my whole corpus (Bootstrap
274 KB: 2,556 blocks / 5,543 declarations / 0 unreadable — identical to rounds 06–08 —
gfonts, normalize, openprops, tw-theme/preflight/utilities, my 20-block stress file) and
over all 37 fixture stylesheets: **identical output, 0 new, 0 lost, 0 newly unreadable**.
The markup path is untouched by construction (the diff never leaves `functionBinding` /
`declarationVarUses`; `markupVarUses` is a separate function), so the claimed 2,032-file
markup parity holds by inspection of the diff.

**Browser premise:** the premise half of the round-08 finding does not depend on the commit
under review — it is a property of the browser. Round-08 established it with
`/tmp/pds-browser-r8/fn.html` (codex-1's probe verbatim) on this same Chrome 150.0.7871.187
binary: a real `CSSFunctionRule`, `CSS.supports("color", "--pick()")` true, computed
`rgb(255, 0, 0)`. I attempted the identical headless re-run three times this round and
Chrome never returned (the host is contended today — dozens of live Chrome processes — where
round 08's fresh-profile runs answered); the checker-side evidence above is what this cycle
changed, so I did not hold the review for it. If a later re-run disagrees with round-08's
result I will amend.

**Lens (PDS conformance):** doctrine byte-identical to what I verified line-by-line in
round 04, so that verification carries. Spot re-checks at this commit: §0–§12 all present
(13 H2 sections); all eight artifact kinds present as H3 entries (DESIGN-BRIEF, DIRECTION,
CRITIQUE, VERDICT, CONTRACT, DESIGN-SYSTEM, AUDIT, WAIVERS) — the invariant four-part shape
holds; §5 rule 7 still states the `DESIGN-SYSTEM.md` Phase-6-authorship rule normatively;
`registry-digest: b49ff596451f` still equals the registry's actual sha256/12.

## Findings

None new. The cycle-12 fix closes the round-08 CRITICAL in both directions, the committed
fixtures pin the exact reproducer and both clean controls, and nothing I probed swings back
toward a false alarm.

### [MINOR] (carried, sixth round — still open, still non-blocking) §4 rule 4's second critique round still has no §1 home

**What.** Unchanged since round 03: PDS §4 rule 4 permits a second critique round on an
explicit Decider instruction, §1's mapping names only `round-02` for CRITIQUE, and the
checker's `PROCESS_HOMES` accepts only that home — a run exercising the permission fails
the L2 process-order check the same spec defines. Doctrine byte-identical since `f1c123d`;
twelve fix-up cycles were all scoped to the checker regression chain, so nothing was ever
going to move this. **Why it matters.** Same as filed: a permitted shape the conformance
level rejects is a trap for the first real run that uses it. It has never blocked and does
not now — it is on the record so the accounting stays complete, not to justify another
round. **Fix.** Unchanged: §1 names the second round's home (and `PROCESS_HOMES` accepts it
when the instruction is recorded), or §4 rule 4 states the second round also files under
`round-02`; one fixture for the permitted shape.

## Open questions

1. **(carried) The differential matrix.** The committed fixture set now pins every
   load-bearing construct of this regression chain, cycle 12's included. What remains
   self-reported are the full sweep numbers (639 stylesheets / 2,032 markup files); I
   verified by direct r8-vs-r9 scanner diff on my own corpus instead — identical. The offer
   stands: commit the construct matrix (strings and expected dispositions, no Chromium
   dependency) as a fixture file.
2. **(carried) The level/verdict axis split.** An open, honestly-recorded violation against
   the winner still yields `verified L3` beside VIOLATION, exit 1. Unchanged; never
   confirmed by the owner as intended.
