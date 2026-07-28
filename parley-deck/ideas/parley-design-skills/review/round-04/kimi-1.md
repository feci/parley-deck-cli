---
agent: kimi-1
idea: parley-design-skills
review-round: 04
date: 2026-07-28
reviewed-commit: e3ca916
---

## Summary

❌ BLOCK. Round-03's four MAJORs are **genuinely closed** — I re-ran every one of my own
reproducers against a pristine `git archive` extraction of `e3ca916`, plus codex-1's three,
and each now fails (or reads) as claimed: the frontmatter lexer refuses the unquoted `#`
loudly, alias direction is strict descent with the doctrine aligned ("points strictly
down"), and the per-agent token sidecar is enforced in all three shapes with the sound-run
fixture rewritten to the conformant one. One finding needs an honest correction: my round-03
MAJOR-1 premise was wrong about browser semantics — per css-syntax-3 §2.2/§5.4 a declaration
containing a bad-string is dropped *wholesale* (recovery seeks the next `;`), so the browser
never applied the `color: #ff0000` in my original probe; the real exposure was rules *after*
the malformed one being swallowed by a runaway string, and that is now caught (probe below).

The block is on the **new fail-safe itself**. It works exactly as designed for every shape it
can see — unterminated string/comment/url, unmatched brace, unbalanced declaration parens,
discarded in-rule text all produce UNJUDGEABLE naming file and reason, exit 4, no roll-up to
PASS. But it is evaded: the scanner treats `{`/`}` as structure regardless of `()`/`[]`
nesting, while css-syntax-3 makes those matched simple blocks whose braces are content. A
`}` inside any parenthesised construct — an ordinary function, `url (` with a space, or an
escaped `\75 rl(` — closes a block in the scanner and not in the browser, and with one
re-balancing `{` every residue check stays quiet: `unreadable: []`, **verified L3, PASS,
exit 0**, while the browser applies a raw `#ff0000`. Three variants confirmed with commands
below. The family the fail-safe was shipped to close is not closed. One CRITICAL.

D-1 position: **ACCEPT, maintained** for the fourth round. 65,360/65,536 (176 spare);
`RULES.md` byte-identical to round-03 (digest `b49ff596451f` recomputed, matches); per-file
thresholds unchanged at 7/25/24/11 KiB.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive e3ca916 | tar -x -C /tmp/pds-r4`;
HEAD `e3ca916e4253`, branch `parley-design-skills`, worktree clean). Probe fixtures were
mutated copies of the shipped `sound-run` fixture under `/tmp/pds-probes-r4`. Nothing in the
skill repo was modified.

- `npm test` (pristine): **227 passing, 0 failing** (212 at round-03).
- Baseline `check --level L3` on the rewritten `sound-run`: `claimed L3, verified L3`, PASS,
  exit 0, digest `b49ff596451f`, 13 rules UNJUDGEABLE, `recusal-not-anchored` for hermes-1.
  Both DIRECTIONs now name adjacent `round-01/<agent>.tokens.json` sidecars; the CONTRACT
  names the winner's. The honest baseline is intact.
- Budgets/digests: `wc -c` = 6,519 + 25,594 + 23,225 + 10,022 = **65,360 ≤ 65,536**;
  `shasum -a 256 RULES.md` = `b49ff596451f` = PDS frontmatter digest; `RULES.md` absent from
  the `17f6619..e3ca916` diff (registry untouched); thresholds in `design-addons.test.js`
  unchanged.
- All eight PDS `yaml` examples parse under the new strict lexer (extracted and run through
  `parseYamlSubset` myself; also pinned by `checker.test.js:1588`). Eight
  `| Field | Requirement |` tables, eight `yaml` fences — the four-part shape holds.
- `x-note: …` extension key on the CONTRACT: accepted, no error — §2 rule 2 honoured.

**Round-03 reproducers, re-run by me:**

- MAJOR-1 (unterminated string): original probe (`.c::before { content: "unterminated` ␍
  `color: #ff0000; }`) still passes — and per css-syntax-3 §2.2 that is now the **correct**
  verdict, because the browser drops the whole malformed declaration (nothing is applied;
  my round-03 claim that it applies `color: #ff0000` was wrong; §2.2: "throw away whatever
  declaration it's currently building, and seek forward until it finds a semicolon (or the
  end of the block)"). The discriminating probe — `.a { content: "unterminated` ␍ `}` ␍
  `.b { color: #ff0000; }`, where the browser drops `.a`'s declaration at the `}` and
  **applies** `.b` — now yields `core:literal-outside-token-layer`, G3 refuted,
  `not verified`, VIOLATION, **exit 1**. The real hole is closed. A governed-prop variant
  (`color: "unterminated` ␍ `#ff0000; }`) over-reports a literal the browser drops — the
  safe direction, and the file is malformed anyway.
- MAJOR-2(a) (`waivers: WAIVERS.md # the waiver file`): fails
  `pds-check:l1-frontmatter-parses` — "an unquoted `"#"` … quote a scalar carrying
  `, [ ] { } or #`, and write a comment as a whole line, never a trailing one" —
  `not verified`, **exit 1**. Quoted control (`waivers: "WAIVERS.md"`) passes, exit 0.
- MAJOR-2(b) (typo'd `waivers: WAVERIS.md`): **still `verified L3`, PASS, exit 0**, waiver
  file silently unread. A directory named `WAIVERS.md` behaves the same (no crash — the
  EISDIR family stays guarded — but silent). Carried as MINOR-1.
- MAJOR-3 (`semantic.muted → {semantic.text}` on the loser's sidecar):
  `pds-check:l3-alias-direction` — "one semantic reaching another descends nothing" — plus
  the G3 refutation, `not verified`, **exit 1**. Shipped `alias-same-tier` fixture covers
  semantic→semantic and component→component. Doctrine and checker now agree on strict
  descent. Closed.
- MAJOR-4 (both DIRECTIONs on shared `../tokens.json`): three `l2-process-order` violations
  (each misnamed sidecar, then the two-to-one resolution), `not verified`, **exit 1**.
  Cross-named variant (claude-1 naming `codex-1.tokens.json`): same. Negative tests at
  `checker.test.js:961–1030` cover misnamed, other-proposer's, and absent sidecars. §1
  rule 3 makes the requirement normative. Closed.
- MINOR (Decider-instructed second-round CRITIQUE at `round-03/codex-1.md`): **still open**
  — `l2-process-order` "filed outside §1's mapping". Nothing in the `17f6619..e3ca916` diff
  touches §4 rule 4 or §1's CRITIQUE row.
- NITs: `TIER_GROUPS` still maps plurals while G3 names singulars (open); the "will not
  take on trust" table still lists five conditions (open — now also missing the sidecar
  check); PDS.md headroom is 6 bytes (25,594/25,600; total has 176 spare).

**Fail-safe attack (three evasion variants, all confirmed; the unbalanced control is caught):**

- ev1: `.a { background: x) \75 rl(}y); color: #ff0000; dummy: z) \75 rl({w: (1); }` —
  `\75 rl` is a hex escape decoding to `url`, so css-syntax-3 §4.3.4/§4.3.6 tokenises
  `\75 rl(}y)` as a url token (content), while the scanner sees `}` and closes `.a` early.
  Result: **verified L3, PASS, exit 0**, `unreadable: []`. Inspecting `scanStylesheet`
  directly: `.a` carries only `background: x) \75 rl(`; `color: #ff0000` vanished; a phantom
  block carries selector `dummy: z) \75 rl(`. The browser applies the literal; no detector
  can see it.
- ev2 (no escapes needed): `.a { background: x) fn(}y); color: #ff0000; dummy: z) fn({w:
  (1); }` — **verified L3, PASS, exit 0**. Any function with a brace in its arguments works;
  this is ordinary, tokeniser-valid CSS.
- ev3: `url (` with a space — **verified L3, PASS, exit 0**.
- ev4 (control, no re-balancing `{`): the fail-safe **fires** — "the closing brace at line
  14 closes no block this scanner had open", all 8 style-reading rules UNJUDGEABLE naming
  file and reason, verdict UNJUDGEABLE, **exit 4**, also without a level claim. The
  mechanism works for the shapes it can see; ev1–ev3 leave it nothing to see.
- Wiring consistency: every detector using `ctx.styles` declares `inputs: ["styles"]`
  (checked all 18); no detector reads stylesheets off-disk itself; `scanStylesheet` is the
  only scan entry point in the check path; markup detectors never extract inline CSS.

**New-obligation defect hunt:** the strict lexer rejects `key:  v` (two spaces), trailing
edge spaces, tabs, and `\` anywhere — each backed by §2 rule 5's grammar line ("key … then
`": `", "no escapes"), each failing loudly, none over-firing on the eight published examples
or the suite fixtures. Sidecar check: non-string `tokens` values degrade to a loud misname
failure; `tokens: ""` on a DIRECTION now fails process-order (stricter than round-03's
honest UNJUDGEABLE — consistent with §1 rule 3). Frontmatter `x-` keys parse. No defect
beyond MINOR-1/2 below.

**Doctrine lens (full `17f6619..e3ca916` PDS.md diff read line by line):** the normative
deltas are exactly four — §1 rule 3 (sidecar, new), §2 DIRECTION's `tokens` row and example
(sidecar naming), §3 G3 ("strictly down"), §12 changelog (records the sidecar) — and each
matches code I verified. Everything else is meaning-preserving compression; I checked every
hunk, including §5 rule 3 dropping "at any point" (the sentence remains universal: "No
numeric aesthetic score is produced") and §4 rule 7's "named in the brief" → "in the brief"
(A3's delegation intact). The mapping table, G1–G4 conditions, ritual, roles/recusal, and
L1–L4 definitions are otherwise untouched, and all eight artifact kinds keep the four-part
shape.

## Findings

### [CRITICAL] The fail-safe is blind to `()`/`[]` nesting: a brace inside any function hides an applied declaration with `unreadable: []` and exit 0

**What.** `lib/css.js`'s main loop treats `{` and `}` as rule structure wherever they appear.
css-syntax-3 §5.4.4 makes `()`, `[]` and `{}` matched simple blocks: a brace inside a
function or bracket is a component value, never rule structure. So
`.a { background: x) fn(}y); color: #ff0000; dummy: z) fn({w: (1); }` tokenises in a
browser as one rule applying `color: #ff0000` (the two odd declarations are invalid and
dropped, but they are *content*), while the scanner closes `.a` at the first `}`, silently
drops `color: #ff0000` as top-level text, opens a phantom block at the second `{`, and ends
balanced: every individually flushed fragment has balanced parentheses, no string/url/
comment is open, every brace matches, nothing is discarded inside a rule — so the fail-safe
records nothing (`unreadable: []`) and the run rolls up to **verified L3, PASS, exit 0**
(ev1/ev2/ev3 above). The escaped form `\75 rl(}y)` needs no malformed input at all — it is
valid CSS that any minifier or hand author can produce. The per-declaration `parenBalance`
check is the fail-safe's only paren guard, and it is evadable fragment-by-fragment: each
flush is balanced while the structure between flushes is wrong.
**Why it matters.** This is the exact failure mode the cycle shipped the fail-safe to
eliminate — a clean L3 certificate beside a raw literal the ratified system really contains,
reachable through the tokenisation family the scanner author correctly diagnosed as
unbounded. The fail-safe was the backstop that was supposed to make the *next* hole cost
rules their verdict instead of passing silently; it does not fire, because the hole leaves
no residue it measures. Round-02 rated the `content: "}"` instance of this shape CRITICAL;
this one is broader (any function, not one quoted property) and needs no malformed input.
**Fix.** Give the main loop css-syntax's block model, not another per-construct patch: a
matched-bracket stack — push on `(`, `[`, `{`; pop on the matching closer; `{` and `}` open
and close rules only when the bracket stack is empty (stray `)`/`]` at depth zero are
component values, as the browser treats them); bracket depth > 0 at EOF, or a closer that
mismatches the open bracket, is a tokenisation failure reported through the existing
`unreadable` channel. This one change kills all three published variants at once — fixing
only the escaped-`url` ident (decoding `\75 rl`) leaves the plain-function variant standing.
Fixtures: ev1, ev2 and ev3 above (each must fail or go UNJUDGEABLE), plus passing-side
controls (`var(--x, {a: b})`-style braces inside functions, and the shipped suite).

### [MINOR-1] A contract whose named waiver file resolves to no readable file is still silently treated as "no waivers"

**What.** Round-03 MAJOR-2 part (b), unchanged: `waivers: WAVERIS.md` (typo) on the CONTRACT
→ `verified L3`, PASS, exit 0, the named file never read; `waivers:` naming a *directory*
behaves identically (no crash, no word). Part (a) — the forbidden-syntax half — is closed.
**Why it matters.** §8's single named waiver file is a control; a run that cannot read it
should say so rather than certify over its absence. Narrowed to MINOR from the round-03
MAJOR because the dangerous direction is impossible: waivers only suppress findings, so a
missing file can only leave findings open (fail-visible the moment anyone relies on one) —
it cannot forge a certificate. But the run's inputs are still misstated in silence.
**Fix.** In the waiver loader, a non-empty `waivers` value that resolves to no readable file
is reported — `NEEDS_REVIEW` (or a `l2-process-order` finding) naming the contract and the
path — not read as null. Fixtures: the typo (reported) and the existing quoted control
(passes).

### [MINOR-2] IMPLEMENTATION.md records nothing past fix-up cycle 4

**What.** The frontmatter says `status: fix-up-cycle-6`, `head-commit: e3ca916`, but the
body ends at cycle 4: no cycle-5/6 entry, no record of the round-03 fixes, no deviation
notes (the strict-descent doctrine ruling; why MAJOR-2 part (b) was left), and the
verification section's count is stale (`212` vs the actual **227**). The claims for this
round reached reviewers through the deck-side prompt and the commit message, not the file
whose job is to carry them.
**Why it matters.** Three rounds running, the review process has operated on "a report is
not evidence — a command and its output is". That rule exists because the fix-up record was
repeatedly optimistic; the record going silent is the same defect in another form, and it
forces every reviewer to reconstruct the claim set from code archaeology.
**Fix.** Add the cycle-5 and cycle-6 sections: fixes applied, deviations from the literal
remedies (with reasons), and the re-run `npm test` count. Five lines per cycle.

### [MINOR-3] (carried from round-03) §4 rule 4's second critique round still has no §1 home

Unchanged this cycle, re-probed: a Decider-instructed second-round CRITIQUE filed at
`round-03/codex-1.md` fails `pds-check:l2-process-order`, so the one exception §4 rule 4
grants is unusable in any checked run. **Fix** (unchanged): one sentence — either §1 names
the second round's home (and `PROCESS_HOMES` accepts it when the instruction is recorded) or
§4 rule 4 states the second round also files under `round-02`. Fixture for the permitted
shape.

### [NIT-1] (carried) Alias-direction group names: checker accepts plurals, doctrine names singulars

`TIER_GROUPS` (`engine.js:1367`) still maps `primitives`/`semantics`/`components`; PDS §3 G3
still names only `primitive`, `semantic`, `component`, so per the doctrine's own text the
conjunct is vacuous for a plural-named document, yet the checker enforces it (now strictly,
which makes the gap slightly wider). One line either side: G3 gains "(singular or plural)"
or `TIER_GROUPS` drops the plurals.

### [NIT-2] (carried, extended) The checker's "will not take on trust" table predates cycles 4–6

`parley-design-check/SKILL.md:169` still lists five recomputed conditions. It omits the
deduplicated rotation, G1's brief-axis count and enumeration check, §1's process-order
locations, the G3/G4 refutations, waiver-ownership exclusion, and now the sidecar check and
the unreadable roll-up (the latter two are documented elsewhere in the file — the trust
table reads as the complete list and is not). Add the rows or retitle the table.

### [NIT-3] (carried, updated) PDS.md has 6 bytes of headroom against its early-warning threshold

25,594 of 25,600 (the binding total has 176 spare: 65,360/65,536). The cycle-6 compression
bought back 4 bytes net of §1 rule 3. Still: the next normative sentence forces compression
or a third rebalance. Recorded so the next cycle plans for it.

## Open questions

1. **(carried) The level/verdict axis split.** An open, honestly-recorded `slop` violation
   against the winner still yields `verified L3` beside verdict VIOLATION, exit 1. My
   reading: consistent with §9, since G1's ban-list conjunct is a sharing test and R-2
   forces the finding into the recorded signature. The owner never confirmed; if that
   reading is wrong, `gateRefutations` needs a G1 arm.
2. **(carried) The second-round critique home** — doctrine question (MINOR-3): where does
   §4 rule 4's exception live?
3. **The scanner strategy, after four holes in one family.** The matched-bracket fix closes
   the *known* family within the zero-dependency constraint. If the owner wants the next
   probe to stop setting the agenda, the doctrine could instead cap what `T1 SOURCE` claims
   to tokenise (e.g. declare the scanner's block model normatively in the checker SKILL.md,
   so "what the checker reads" is itself a specified, reviewable artifact rather than an
   approximation discovered by probes). Not a condition of unblocking — the CRITICAL fix is.
