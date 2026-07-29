---
agent: kimi-1
idea: parley-design-skills
review-round: 11
date: 2026-07-28
reviewed-commit: 1804985
---

## Summary

✅ ACCEPT. The round-10 MAJOR all three reviewers converged on is genuinely closed, and I
verified the closure differentially rather than by trusting the new test: my own round-10
reproducer — `<div style="content: 'var(--ghost)'; color: var(--color-text-body)">` in a
fresh copy of the shipped sound-run fixture — still produces the false
`core:token-used-undeclared`, VIOLATION, `verified=null`, exit 1 at a pristine extraction of
`1675b6f`, and produces PASS, verified L3, exit 0 at `1804985`. The swing-back controls all
point the right way: a real undeclared inline-style reference still costs the run its
certificate, a host-language JS string still reports, a `<style>`-body string stays silent,
the single-quoted attribute spelling is handled, and a CRLF file blanks the ghost on line 1
while still reporting the real reference on line 3 — the new span-offset math holds
end-to-end, not just in the code.

The boundary that refused the naive fix is real, and I reproduced both sides of it on my own
inputs before judging it. The corpus shape — `style="…background:${s==='All'?'var(--c1)':'rgba(…)'}…"`
inside a JS template literal — still reports `--c1` (exit 1): blanking every `style=` match
would have lost it, which is the false clean, so the exclusion being of the span the CSS
path *read* rather than of everything the pattern matches is load-bearing, not rhetoric. The
counter-shape — a Jinja conditional inside `style=` — passes with no unreadable entry:
routing those spans to the unreadable channel would strip certification from ordinary
shipped template code, the false alarm the doctrine itself equates in damage with a false
clean.

**The disclosed residual is a shipping state, not a defect, and disclosing it is the right
call.** I reproduced the residual itself (a template-literal `style=` span carrying
`content:'var(--ghost)'` reports `--ghost`), so the disclosure is accurate, not aspirational.
It is bounded (one precise parse test — `blocks.length === 1 && unreadable.length === 0` —
decides membership, and both boundary shapes are pinned in the committed test), measured
(72 of 8,307 spans in the implementer's corpus; my independent replication over my own
200-file corpus found 44 spans, 38 withheld, 6 left — all six synthetic template shapes,
zero in 159 real accumulated markup files), one-directional (over-report only; the
false-clean direction is excluded by construction, since anything left in the sweep is still
read), and recoverable (a false finding is exactly what the waiver mechanism exists for,
with scope, expiry and counter-signature). The false finding needs a triple conjunction — a
residual-class span, a `var()` spelling inside a quoted context, and that token undeclared
in the ratified document. The alternative (unreadable channel) converts that narrow,
waiver-able over-report into whole-file UNJUDGEABLE on valid code. On the evidence I
produced, disclosure beats both refused options.

npm test 247/247 at a pristine extraction. Doctrine byte-identical a fourteenth cycle:
65,360/65,536, digest `b49ff596451f` — D-1 accepted an eleventh round. Direct r10-vs-r11
scanner differential over 245 files (200 markup + 45 stylesheets): 11 files changed, every
one the fixed class (a `var()` spelling inside an opaque token in a `style=` attribute no
longer raw-reported; every real reference kept at the same line); 0 new, 0 real lost, 0
newly unreadable.

## What I verified (commands run, and their result)

All probes ran against pristine extractions (`git archive 1804985 | tar -x -C /tmp/pds-r11`;
round-10's extraction of `1675b6f` survived at `/tmp/pds-r10`; reviewed commit is HEAD of
`parley-design-skills`, worktree clean, `git show -s` confirms
`1804985f91975f075198f087e708d2b4b766dce4`). Probe runs are fresh copies of the shipped
`conformance/sound-run` fixture — verified `diff -r`-identical between the two builds —
under `/tmp/pds-probes-r11`. Nothing in either repo was modified.

- `npm test` (pristine extraction, Node v26.5.0): **247 passing, 0 failing** — matches the
  claimed 246 → 247.
- `git diff 1675b6f..1804985 --stat`: only `lib/css.js` (+91/−12 region) and
  `test/checker.test.js` (+150). `git diff` over `addons/parley-design/` empty; `wc -c` =
  6,519 + 25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**; `shasum -a 256 RULES.md` 12-hex =
  `b49ff596451f` = PDS `registry-digest`.
- **The code, read before the probes.** `styleAttributes` now records the exact `[start,
  end)` span of each attribute value (line-split on `\n`, trailing `\r` stripped, offsets
  accumulated over the unstripped line so they stay exact against the original text — the
  math checked by hand and then probed via CRLF below). Each value is scanned as
  `x{<value>}`; the span is blanked from the raw sweep by `blankOffsets` (newline-preserving)
  only where that one scan returned exactly one block and nothing unreadable — the span kept
  out of the sweep is the span the CSS path read, by construction, one scan answering both
  questions. Blanking order is correct: `<style>` bodies and comments are blanked by their
  patterns first, so a `-->` inside an attribute value cannot end a comment early against
  text already blanked out from under it.
- **My round-10 reproducer, differentially (fresh sound-run copies, `--level L3 --json`,
  `/tmp/pds-probes-r11/run-html.sh <build> <probe> <label>`):**

  | probe | `1675b6f` | `1804985` |
  |---|---|---|
  | `style="content: 'var(--ghost)'; color: var(--color-text-body)"` | VIOLATION, `--ghost` named, exit 1 | **PASS, verified L3, exit 0** |

  The fix is load-bearing, not asserted. **Controls at `1804985`:** `style="color:
  var(--really-undeclared)"` → VIOLATION, exit 1 (a real inline-style reference still
  fires); `const cl = "color:var(--ghost)"` in `<script>` → VIOLATION, exit 1 (the
  deliberate host-language non-masking intact); `<style>.probe{content:"var(--ghost)"…}</style>`
  → PASS (the `<style>` half untouched); single-quoted `style='content: "var(--ghost)"…'` →
  PASS; CRLF file with the ghost attribute on line 1 and `style="color:
  var(--really-undeclared)"` on line 3 → VIOLATION naming line 3, the ghost blanked (offset
  and line-number integrity through `blankOffsets`).
- **The corpus's refusal of the naive fix, reproduced:** `el.innerHTML = \`<div
  style="padding:8px;background:${s==='All'?'var(--c1)':'rgba(0,0,0,.5)'};color:#fff">\`` →
  VIOLATION naming `--c1` at probe.html:2, exit 1 — the span fails the declaration-list
  parse (scanner verdict: `blocks=2`), stays in the sweep, and the real reference survives.
  **The Jinja counter-shape:** `style="width: {% if wide %}600px{% else %}300px{% endif %};
  color: var(--color-text-body)"` → PASS, verified L3, no unreadable — demonstrating the
  false alarm the unreadable channel would have caused.
- **The disclosed residual, reproduced:** `el.innerHTML = \`<div
  style="content:'var(--ghost)';width:${w}px">\`` → VIOLATION naming `--ghost`, exit 1. The
  quotes here are literal template text, so post-substitution the browser sees a CSS string
  and resolves nothing — the true false-positive shape (hermes-1's narrowing confirmed:
  inside `${…}` the quotes are JavaScript's and the reference is real; only a var() inside
  the template's literal text over-reports). Disclosure verified as accurate.
- **Independent residual measurement.** `styleAttributes` is not exported, so I replicated
  it from the shipped diff verbatim (same regex, same line split) and applied the shipped
  qualification rule via the exported `scanStylesheet`, over my own 200-file corpus (41
  generated `style=` matrix files + 159 accumulated real probe/browser/fixture markup files
  from rounds 4–10): **44 spans, 38 withheld, 6 left in the sweep** — and all six are my
  synthetic `${…}`/`{%…%}` template shapes; zero real files have a span in the residual
  class. The mechanism classifies exactly as disclosed.
- **Full scanner differential, r10 vs r11** (round-10's `sweep.js`: `scanStylesheet`
  structure + `varUses` + `markupVarUses` per file, over 200 markup + 45 stylesheet files —
  my 8-file CSS corpus, all 37 fixture stylesheets, the matrix, and every accumulated real
  markup file): **11 files changed, all in `muses`, every one the fixed class.** Each
  inspected individually: my round-10 reproducer itself (`--ghost` dropped, `--color-text-body`
  kept), the six opaque-token matrix files (`--ghost-*` dropped), and the mixed files (real
  references kept at the same lines, ghosts dropped — including the CRLF and duplicate-attribute
  cases). **0 new references, 0 real references lost, 0 newly unreadable**, and identical
  output on all 45 stylesheets and all 159 real markup files.
- **Mixed-shape boundary probes** (template/Jinja/unterminated-string spans that also carry
  a real reference): all four report the real reference exactly once. `markupVarUses` dedups
  by `name@line`, so a dirty span found by both paths cannot double-report; the committed
  test's exactly-once assertion on the *clean* path is what keeps that dedup honest.
- **The committed test** ("a style attribute is read as CSS once, and never raw-scanned a
  second time") pins codex-1's reproducer in both quote spellings plus the url form, the
  exactly-once control (including an escaped `\76 ar(--escaped)` that only the CSS path can
  find — proving collection route, not just membership), the four host-language shapes, the
  `<style>` body, **both boundary shapes** (template literal and Jinja) with their real
  references, the end-to-end L3 certification, and three swing-back guards. It is the filed
  remedy with passing-side controls, and it passes.
- **Registry refusal, capability, budgets (mechanical re-checks at this commit):** absent
  registry → explicit refusal on stderr, structural checks still run, exit 3; capability
  line reports 18 detectors generated from `lib/detectors`; the 64 KiB total is enforced by
  `test/design-addons.test.js` (`TOTAL_BUDGET = 64 * 1024`).
- **Lens (PDS conformance):** doctrine byte-identical to what I verified line-by-line in
  round 04 and re-spot-checked in rounds 09/10, so that verification carries. Mechanical
  re-checks at this commit: §0–§12 all present (13 H2 sections); all eight artifact kinds
  present as H3 entries; the four-part shape confirmed structurally per entry (DESIGN-BRIEF
  inspected directly: one-line purpose → rationale sentence → required-fields table → yaml
  minimal example); changelog present; `registry-digest` equals the registry's actual
  sha256/12.

## Findings

### [CRITICAL] None.

### [MAJOR] None.

### [MINOR] (carried, eighth round — still open, still non-blocking) §4 rule 4's second critique round still has no §1 home

**What.** Unchanged since round 03: PDS §4 rule 4 permits a second critique round on an
explicit Decider instruction, §1's mapping names only `round-02` for CRITIQUE, and the
checker's `PROCESS_HOMES` still accepts only that home — a run exercising the permission
fails the L2 process-order check the same spec defines. Doctrine byte-identical since
`f1c123d`; fourteen fix-up cycles were all scoped to the checker, so nothing was ever going
to move this. **Why it matters.** Same as filed: a permitted shape the conformance level
rejects is a trap for the first real run that uses it. It has never blocked and does not
now. **Fix.** Unchanged: §1 names the second round's home (and `PROCESS_HOMES` accepts it
when the instruction is recorded), or §4 rule 4 states the second round also files under
`round-02`; one fixture for the permitted shape.

### [NIT] A multi-line `style=` attribute value never enters the span accounting at all

**What.** `styleAttributes` splits the text into lines before matching, so an attribute
value spanning lines — legal HTML, and produced by template formatters — is never extracted:
no CSS-path read, no span, no blanking. Verified: `<div\nstyle="content: 'var(--ghost-sq)';\ncolor:
var(--undeclared-one)"\n>` reports both `--ghost-sq@2` (a var spelling inside what a browser
treats as an opaque string — the false-finding class) and `--undeclared-one@3` (real),
identically at `1675b6f` and `1804985`. This is the disclosed residual in effect, but it
sits outside the disclosure's wording: such a value never becomes one of the 8,307 "spans",
so it is not among the 72 either, and the 72 are described as "host-language text that
happens to be spelled like an attribute" — which a multi-line value is not; it is genuine
CSS the line-based extractor cannot see. **Why it matters.** Very little, which is why it is
a NIT: pre-existing (identical at the round-10 build), over-report direction, recoverable by
waiver, and absent from every real file I hold (0 occurrences in 159 accumulated real markup
files, and the corpus differential cannot see it either — both builds sweep it identically,
so "0 files changed" does not bound it). It matters only to the accuracy of the disclosure,
which should name this member of the family alongside the 72. **Fix.** Run `STYLE_ATTRIBUTE`
over the whole text instead of per line (its `[^"]*` already spans newlines), derive the
line from the match offset, and let the same `blocks === 1 && unreadable === 0` test decide
withholding — `blankOffsets` already preserves newlines, so the span math needs no change.
Do it with the same differential sweep this cycle ran, because whole-text matching re-opens
quote-pairing questions (an opener and its closer on different lines) that the per-line form
never had to answer. A documentation-only alternative: extend the disclosure sentence to
cover values the extractor cannot match.

## Open questions

1. **(carried, narrowed) The differential matrix.** The committed test now pins the `style=`
   boundary shapes in both directions, and I verified the wider claim by direct r10-vs-r11
   scanner diff over my own 245-file set instead of trusting the sweep numbers (11,675
   references, 0 files changed — self-reported; my replication found 0 real files in the
   residual class). The offer stands: commit the construct matrix as a fixture file.
2. **(carried) The level/verdict axis split.** An open, honestly-recorded violation against
   the winner still yields `verified L3` beside VIOLATION, exit 1. Unchanged; never
   confirmed by the owner as intended. This is the last planned round, so it is now a
   question for the first real run rather than for a review cycle.
