---
agent: kimi-1
idea: parley-design-skills
review-round: 10
date: 2026-07-28
reviewed-commit: 1675b6f
---

## Summary

🟡 ACCEPT-WITH-RESERVATIONS. The round-09 MAJOR is genuinely closed in both directions,
and I verified the closure differentially rather than by trusting the new test: all five
opaque-token probes (codex-1's quoted and url-quoted `@function` defaults, hermes-1's
ordinary-declaration string in both quote spellings, and the unquoted url token) produced
the false `core:token-used-undeclared` VIOLATION with exit 1 and `verified=null` at a
pristine extraction of `9121ec2`, and all five produce PASS, verified L3, exit 0 at
`1675b6f`. The swing-back controls point the right way: the round-08 reproducer still
flags `--ghost` twice and costs the run its certificate, a real `var()` beside a string in
the same value and beside a string default in the same parameter list is still collected,
and an escaped `\76 ar(--ghost)` default is decoded before masking and still flagged — the
decode-then-mask ordering confirmed end-to-end, not just in the code.

The deliberate non-masking of the markup raw sweep — the decision this round puts to me —
I **accept**, on evidence I produced myself: routing host-language lines through the CSS
masker blanks every real shape the corpus carries (`const cl = "color:var(--error)"`,
`className="text-[var(--x)]"`, `stopColor="hsl(var(--primary))"`,
`['hsl(var(--chart-2))']`), and a browser resolves every one of them. That is the measured
2,410 references across 306 of 2,236 files: masking would buy a rare false finding back at
the price of a silent false clean over 13.7% of real markup — the worse direction, and the
one this rule exists to prevent. The asymmetry is principled (a quote is a string-token
delimiter only where CSS is actually being parsed) and both halves of it behave end-to-end
exactly as the new assertion locks them.

One new MAJOR, adjacent to — not inside — this cycle's change: the claim that the
opaque-token rule reaches `style=` attributes is false, because the raw sweep re-reads
those attribute values unmasked, so a `var()` inside a CSS string in an inline style costs
a run its L3 certificate over markup a browser resolves. It is pre-existing (verified
identical at `9121ec2`), narrow, over-reports rather than cleans, and has a concrete
loss-free fix — which is why it is a reservation and not a block: the change under review
is verified correct in both directions, and holding it for a seam three cycles older than
itself conflates the commit's verdict with the surface's history. It should land as a
fast-follow before any v1.0.0 tag. (Written before I read the other round-10 signoffs;
codex-1 and hermes-1 found the same seam independently — codex-1 blocks on it, hermes-1
files it at my disposition. The convergence marks it real; the severity split is over
whether a pre-existing, over-reporting seam holds a verified commit, and on that I side
with hermes-1.) The carried §4-rule-4 MINOR stands a seventh round; the doctrine was
untouched again, so nothing was going to move it.

npm test 246/246. Doctrine byte-identical: 65,360/65,536, digest `b49ff596451f` — D-1
accepted a tenth round. Direct r9-vs-r10 scanner diff over my whole corpus and all 37
fixture stylesheets — structure, var-uses and markup-uses: identical, 0 new, 0 lost, 0
newly unreadable.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive 1675b6f | tar -x -C /tmp/pds-r10`;
reviewed commit is HEAD of `parley-design-skills`, worktree clean; cycle 13 is exactly one
commit over `9121ec2`). Probe runs are fresh copies of the shipped `conformance/sound-run`
fixture (verified `diff -r`-identical to round-09's) under `/tmp/pds-probes-r10`; the
round-06 corpus survived at `/tmp/pds-corpus`; the round-09 build survived at `/tmp/pds-r9`
for the differential. Nothing in either repo was modified.

- `npm test` (pristine extraction, Node v26.5.0): **246 passing, 0 failing** — matches the
  claimed 245 → 246.
- `git diff 9121ec2..1675b6f --stat`: only `lib/css.js` and `test/checker.test.js`. `bin/`
  and `lib/engine.js` byte-identical; `git diff` over `addons/parley-design/` empty;
  `wc -c` = 6,519 + 25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256
  RULES.md` 12-hex = `b49ff596451f` = PDS `registry-digest`.
- Baseline `check --level L3` on the unmodified sound-run: **PASS, exit 0**, 13
  unjudgeables — the round-09 baseline, undisturbed.

**The code, read before the probes.** The fix is the one codex-1 filed: one collector,
`valueVarUses` (`css.js:1332`), which runs `VAR_REFERENCE` over the value after
`maskOpaqueTokens` blanks string and url token contents, and both collectors — the
`@function` default reader in `functionBinding` (`css.js:876-885`) and ordinary
declarations in `declarationVarUses` (`css.js:1372`) — go through it. The masking runs
after each caller's decode, so `\76 ar(` is still found. `VAR_REFERENCE` now has exactly
three occurrences in the file: its definition, the masked collector, and the declared
markup raw sweep — no second raw-text collector survives anywhere in the CSS path. One
`stringToken`, one `consumeEscape`, one `maskOpaqueTokens`: the one-consumer discipline
cycle 11 established is intact.

**My round-09 reproducer batch, re-run byte-for-byte (fresh sound-run copies, `--level L3
--json`):** codex-1's probe → exit 1, `--ghost` named **twice** at probe.css:2, G3 refuted,
`verified=null`; my default-only half → exit 1, named at line 1; clean control
(declared-default-plus-true-formal) → **PASS, verified L3, exit 0**; true formal → exit 1
on the planted literal only, `var(--x)` suppressed; undeclared body → exit 1; plain
`var(--ghost)` → exit 1. All six identical to round-09.

**The round-09 finding, differentially (same probes at both builds, end-to-end):**

| probe | `9121ec2` | `1675b6f` |
|---|---|---|
| `@function --quoted(--x: "var(--ghost)") returns <string>` | VIOLATION, exit 1 | **PASS, verified L3, exit 0** |
| `@function --pick(--x: url("var(--ghost)"))` | VIOLATION, exit 1 | **PASS, verified L3, exit 0** |
| `.probe { content: "var(--ghost)"; … }` | VIOLATION, exit 1 | **PASS, verified L3, exit 0** |
| `.probe { background: url(var(--ghost)); … }` | VIOLATION, exit 1 | **PASS, verified L3, exit 0** |
| `.probe { font-family: 'var(--nope)'; … }` | VIOLATION, exit 1 | **PASS, verified L3, exit 0** |

The fix is load-bearing, not asserted. **Swing-back controls at `1675b6f`:** `content:
"var(--nope)" var(--real)` → `--real` flagged, `--nope` silent; `--f(--x: "var(--ghost)",
--y: var(--gap))` → only `--gap` flagged; `--x: \76 ar(--ghost)` → flagged (decode precedes
mask); the round-08 reproducer → still VIOLATION, exit 1. The oscillation is not swinging.

**The markup asymmetry, probed directly (node, both exported functions):** the raw sweep
finds all five host-language shapes listed above; `maskOpaqueTokens` over the same lines
finds **none** — the mechanism behind the 2,410/306 measurement, reproduced on my own
inputs. End-to-end: a string inside a `<style>` body (`content: "var(--ghost)"`) → PASS,
exit 0; a JS string in markup (`const cl = "color:var(--ghost)"`) → VIOLATION, exit 1; a
markup comment carrying `var(--hidden)` → blanked, not reported; `style=` attributes are
parsed as declaration lists (real uses collected once — see the finding for the seam).

**The committed test** ("an opaque token's contents are not a var() reference, in a
declaration or in an @function default") pins all five reproducers with empty expected
uses, five controls with their uses, both halves of the markup asymmetry, the four
end-to-end L3 certifications, and the round-08 reproducer as the explicit swing-back guard
("Masking that suppressed this would be the false clean cycle 12 closed, returning under a
new name"). It is the filed remedy with passing-side controls, and it passes.

**Sweeps:** direct r9-vs-r10 diff of `scanStylesheet` structure **plus** `varUses` and
`markupVarUses` output over all 37 fixture stylesheets and my 8-file corpus (Bootstrap
274 KB, gfonts, normalize, openprops, Tailwind theme/preflight/utilities, my stress file —
45 files): **byte-identical output, 0 new, 0 lost, 0 newly unreadable.** Real shipped CSS
carries no `var(` text inside strings or url tokens, which is exactly the cycle's claim
that the defect fired on reviewer probes rather than the corpus — and why the probes are
the right place to test a gate.

**Browser premise:** attempted two headless runs of a purpose-built page (quoted-default
pseudo-element content, unquoted `url(var(--ghost))`, string `content`) against this host's
Chrome; both hung without output — the host is contended today (54 live Chrome processes),
as it was in round-09. The premise halves this cycle needs are (a) a string default
computing as literal text — codex-1 verified it in Chromium 150 in round-09 and it is
recorded in the committed test, and (b) an unquoted `url(` consuming one opaque url token —
CSS Syntax §4.3.6, core tokenization that does not drift between builds; round-08's Chrome
150 run on this same binary established the `@function` substitution half. The
checker-side change is what this commit contains, and it is verified above; I did not hold
the review for a redundant browser confirmation. If a later re-run disagrees I will amend.

**Lens (PDS conformance):** doctrine byte-identical to what I verified line-by-line in
round 04 and re-spot-checked in round 09, so that verification carries. Mechanical
re-checks at this commit: §0–§12 all present (13 H2 sections); all eight artifact kinds
present as H3 entries, each carrying a required-fields table, a yaml minimal example and
its prose — the invariant four-part shape holds everywhere; the C1 mapping table stands in
§1; `registry-digest: b49ff596451f` equals the registry's actual sha256/12.

## Findings

### [MAJOR] `style=` attribute values are re-read unmasked by the raw sweep — the opaque-token rule does not actually reach them

**What.** Cycle 13's comment in `markupVarUses` claims: "The spans that really are CSS — a
`<style>` body and a `style` attribute — are parsed as CSS above and collected through
`valueVarUses` … so the opaque-token rule reaches them there." For `<style>` this is true:
those bodies are blanked before the raw sweep, so they are read only through the masked CSS
path (verified end-to-end: PASS). For `style=` attributes it is not: `styleAttributes`
extracts and parses them correctly through `declarationVarUses`/`valueVarUses`, but their
bytes are **not** blanked from the raw sweep, which then re-reads the attribute as markup
text. A `var()` inside a CSS string in an inline style is therefore still reported:

```html
<div style="content: 'var(--ghost)'; color: var(--color-text-body)">x</div>
```

— `core:token-used-undeclared` names `--ghost`, VIOLATION, `verified=null`, exit 1,
verified at the CLI against a fresh sound-run copy. The browser substitutes nothing (the
string token is opaque); this is the same false-finding class round-09 blocked on,
relocated to a much narrower surface. Identical at `9121ec2` (verified same command, same
result), so it is **pre-existing** — introduced with the raw sweep in cycle 10, not by this
cycle — and the new assertion locks the `<style>` half and the host-language half but not
this one, which the prose names.

**Why it matters.** A valid inline style costs a run its L3 certificate and exits 1 — the
false-finding direction the doctrine itself equates in damage with a false clean ("it is
how a gate gets switched off"), and the same verdict class round-09 blocked on. That makes
it MAJOR in kind. What keeps it a reservation rather than a block: the direction is
over-report, not a false clean; the trigger is rare (a string or quoted-url token spelling
`var(--…)` inside a `style=` attribute, with the token undeclared — no real file in
thirteen cycles of corpus sweeps shows the shape); it is pre-existing rather than
introduced by the commit under review; and the fix is a surgical blank with no loss. The
cycle's own prose now asserts coverage it does not have, and the double-read is exactly
the two-readings-of-one-token defect this cycle eliminated everywhere else — so it must
be fixed, but as a fast-follow, not by holding a verified commit.

**Concrete fix.** Blank the `style=` attribute values from the raw sweep the same way
`<style>` bodies are blanked: the spans `styleAttributes` already extracts are precisely
the bytes the CSS path has read, so blanking them loses nothing — anything the extractor
cannot match (an unquoted attribute, a value spanning lines) stays visible to the raw
sweep, and a malformed extracted value already fails safe through the unreadable channel.
Two fixtures beside the existing asymmetry assertions: `style="content: 'var(--ghost)'"`
must not report `--ghost`, and `style="color: var(--real)"` must still report `--real`.

### [MINOR] (carried, seventh round — still open, still non-blocking) §4 rule 4's second critique round still has no §1 home

**What.** Unchanged since round 03: PDS §4 rule 4 permits a second critique round on an
explicit Decider instruction, §1's mapping names only `round-02` for CRITIQUE, and the
checker's `PROCESS_HOMES` (`engine.js:446`) still accepts only that home — a run exercising
the permission fails the L2 process-order check the same spec defines. Doctrine
byte-identical since `f1c123d`; thirteen fix-up cycles were all scoped to the checker, so
nothing was ever going to move this. **Why it matters.** Same as filed: a permitted shape
the conformance level rejects is a trap for the first real run that uses it. It has never
blocked and does not now. **Fix.** Unchanged: §1 names the second round's home (and
`PROCESS_HOMES` accepts it when the instruction is recorded), or §4 rule 4 states the
second round also files under `round-02`; one fixture for the permitted shape.

## Open questions

1. **(carried, narrowed) The differential matrix.** The committed fixture set now also
   pins the opaque-token matrix in both directions, including the markup asymmetry halves.
   What remains self-reported are the full sweep numbers (2,933 files); I verified by
   direct r9-vs-r10 scanner diff on my own 45-file sample instead — identical. The offer
   stands: commit the construct matrix as a fixture file.
2. **(carried) The level/verdict axis split.** An open, honestly-recorded violation
   against the winner still yields `verified L3` beside VIOLATION, exit 1. Unchanged;
   never confirmed by the owner as intended.
