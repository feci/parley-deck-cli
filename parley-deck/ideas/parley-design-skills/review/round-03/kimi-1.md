---
agent: kimi-1
idea: parley-design-skills
review-round: 03
date: 2026-07-28
reviewed-commit: 17f6619
---

## Summary

🟡 ACCEPT-WITH-RESERVATIONS. Every round-02 finding I blocked on is now **genuinely
closed, verified by re-running my own reproducers against a pristine `git archive`
extraction of `17f6619`** — all seven forgeable-certificate paths, the waiver-ownership
hole, both empty-path crashes, the fast-path contradiction, and the alias-direction
conjunct. Each now fails (or runs clean) exactly as the fix-up record claims, with the
right messages, and the passing-side controls still pass. The doctrine layer remains
conformant under my lens.

The reservations are in the **new surface the cycle-4 checks added**. My own hunt found
one MAJOR (the new quote-aware CSS scanner still loses a literal inside an unterminated
string) and one MINOR; codex-1's round-03 block names three more, and I independently
reproduced all three before crediting them — two are MAJOR (the frontmatter parser
silently accepts the subset-forbidden unquoted `#`, with a clean exit-0 L3 in the
waivers-path case; §1's per-agent `<agent>.tokens.json` sidecar naming is unenforced and
the shipped sound-run fixture itself violates it, hollowing out G2's graft check), and I
concur with codex-1/hermes-1 that the same-tier alias question is MAJOR doctrine
ambiguity ("points down" vs the checker's `to > from`). Four MAJORs, none systemic like
AF-1, each with a small concrete fix — that is a reservations list, not a block, but
they should land (with their fixtures) before any v1.0.0 tag.

D-1 position: **ACCEPT, maintained** for the third round. The binding 64 KiB total held
at 65,364/65,536; all 30 rule YAML blocks are byte-identical to round-02 (I diffed the
fences, not the prose), and the rebalance remains the honest resolution of C3.

## What I verified (commands run, and their result)

All probes ran against a pristine extraction (`git archive 17f6619 | tar -x -C /tmp/pds-r3`;
HEAD `17f6619a2885`, branch `parley-design-skills`, worktree clean). Probe fixtures were
built under `/tmp/pds-probes` as mutated copies of the shipped `sound-run` fixture.
Nothing in the skill repo was modified.

- `npm test` (pristine): **212 passing, 0 failing** — IMPLEMENTATION.md's cycle-4 count
  holds at the reviewed commit.
- Baseline `check --level L3` on `sound-run`: `claimed L3, verified L3`, PASS, exit 0,
  registry digest `b49ff596451f`, 13 rules UNJUDGEABLE, `recusal-not-anchored` printed
  for hermes-1. The honest baseline is intact.

**Round-02 reproducers, re-run by me (each closed):**

- (a) DIRECTIONs↔CRITIQUE rounds swapped: `not verified`, three
  `pds-check:l2-process-order` violations naming §1's homes, VIOLATION, exit 1.
  (Round-02: verified L2, PASS, exit 0.)
- (b) invented axis `texture` + agreement on one brief axis: `not verified` — G1 counts
  declared axes only ("differ on 1 declared axis; 2 are required") and each invented
  axis is its own violation. Exit 1. (Round-02: verified L2, PASS, exit 0.)
- (c) brief primary `[flat, flat, layered]`, both proposers recording `flat`:
  `not verified` — `l2-assignment` fails on the duplicate enumeration *and* on the
  recorded position mismatching the deduplicated rotation. Exit 1.
- (d) raw `#ff0000` + G3 recorded `pass`: certificate now reads `not verified`;
  `l2-gate-recorded` fires — "G3 is recorded pass and this run's own findings refute its
  conditions: core:literal-outside-token-layer" — plus `l3-system-rules`. Exit 1.
  (Round-02: `verified L2` beside the open violation.)
- (e) AUDIT recording G4 `pass` + open `core:focus-indication`: `not verified`, same
  refutation shape naming both open quality rules. Exit 1.
- (f) `.trap::before { content: "}"; color: #ff0000; }` at L3: the literal is now seen
  (quote-aware scanner), G3 refuted, `not verified`, exit 1. (Round-02: verified L3,
  PASS, exit 0.)
- (g) `kind: TYPO-BRIEF` file: `pds-check:l1-artifact-kind` violation, `not verified`,
  exit 1. Companion probe: a copy of `PDS.md` itself (spec declared, no kind) stays
  uninspected — verified L1, PASS, exit 0, so the doctrine never fails its own checker,
  exactly as §2 rule 1's new sentence states. The deviation from codex-1's literal
  "missing kind is a violation" remedy is the right call.
- AF-2: waiver scoped at `round-01/claude-1.md`, granted by codex-1, counter-signed by
  **claude-1** (the scoped file's author): `waiver rejected: the counter-signer claude-1
  authored the waived work … a self-approval rather than a counter-signature`; finding
  stays open; `not verified`, exit 1. Controls: the same waiver counter-signed by
  **codex-1** (corroborated, not the author) suppresses cleanly — verified L3,
  `waived 1`, exit 0; counter-signed by **hermes-1** (recorded only by its own critique)
  is rejected as uncorroborated. Closed without over-firing.
- `waivers: ""` on the CONTRACT: clean run, verified L3, PASS, exit 0 (round-02: EISDIR
  crash). `tokens: ""` on the CONTRACT: clean. `tokens: ""` on the winning DIRECTION: no
  crash — G2 reads `unverified` ("names no tokens file"), UNJUDGEABLE, exit 4, the
  honest treatment. Closed at all three sites.
- Alias direction: the shipped `alias-direction` fixture fires `l3-alias-direction` on
  both a primitive→primitive alias and a semantic→component alias; the legal
  component→semantic and semantic→primitive aliases pass. The vacuous case is real:
  `sound-run`'s own `color`/`space` groups declare no direction and verify L3. Closed.
- Fast path: PDS §4 rule 8 now says the fast path "runs outside Parley", is "not a
  Parley workflow", "MUST NOT claim Parley verification, nor a level above L1";
  SKILL.md says the same and the dispatcher row reads "Fast path (outside Parley, §4
  rule 8)". Closed — the sentence my round-02 MINOR asked for.
- Regression spot-checks: R-1 token drift (shipped `token-drift` fixture) still fails G2
  on the digest comparison; R-2's observed-vs-recorded signature check fired in my own
  AF-2 probe when the recorded `fires: []` omitted a finding the run watched fire.
- AF-3: `check NOTICE.md` → exit 4. C4 refusal: checker copied alone to /tmp → stderr
  refusal naming the absent registry, exit 3, structural checks still ran. L4 claim:
  `not verified`, UNJUDGEABLE, exit 4.
- Budgets and digests: `wc -c` = 6,519 / 25,598 / 23,225 / 10,022 = **65,364 ≤ 65,536**
  (172 spare; per-file thresholds 7/25/24/11 KiB confirmed in `design-addons.test.js`);
  `shasum -a 256 RULES.md` begins `b49ff596451f` = PDS frontmatter's `registry-digest`.
  All 19 core + 11 web `pds-rule` fences **byte-identical** between `8ebd8f7` and
  `17f6619` (diffed the extracted fences) — the cycle changed rule prose only, never a
  rule. Ban-list definition intact (RULES.md:278–287). Zero placeholder strings; no
  `RULES.md` under the checker; 18 detectors, capability line generated.
- Doctrine lens (full diff of PDS.md `8ebd8f7..17f6619` read line by line): §1 mapping
  table row-for-row C1; §1 rule 2 now makes gate-recomputation normative ("the outcome
  is the gate's conditions, never the word"); §2 rule 1 states the unknown-kind/no-kind
  distinction; §3 G3 defines alias direction by group name including the vacuous case;
  §8 rule 2 carries both new independence conjuncts; §4 rule 2 matches the checker's
  deduplicated `rotateAssignment` (sorted, sha256 `"PDS/1"||run_id`, first 8 hex
  big-endian, modulo count — code read and output compared); all eight artifact kinds
  keep the four-part shape (purpose line, rationale, `| Field | Requirement |` table,
  `yaml` example — checked all eight); §9 L2/L3 definitions unchanged; §10–§12 intact;
  the 1.0.0 changelog entry explicitly records alias direction and §8's authorship test.
  `states:` examples shrinking from 8 entries to 5 is cosmetic — no doctrine text
  enumerates a required state set.

**codex-1's round-03 findings, independently reproduced by me before crediting:**

- `waivers: WAIVERS.md # the waiver file` on the CONTRACT (unquoted `#`, forbidden by §2
  rule 5's scalar grammar): **verified L3, PASS, exit 0**, waiver file silently unread.
  On `kind:` the same trailing text fails loudly as an unknown kind — the silent case is
  specific to fields whose corrupted value the run can ignore. See MAJOR-2.
- `semantic.muted → {semantic.text}` (same-tier edge): no `l3-alias-direction` finding;
  the checker tests `to > from`, not `to >= from`. See MAJOR-3.
- The shipped `sound-run` fixture: both DIRECTIONs name `tokens: ../tokens.json` — one
  shared file at the run root — while §1's DIVERGE row names `round-01/<agent>.md +
  <agent>.tokens.json`. It verifies L3. See MAJOR-4.

## Findings

### [CRITICAL] None.

The two round-02 CRITICALs are closed with commands above; no new critical path found.

### [MAJOR-1] Unterminated CSS string still hides declarations a browser applies

**What.** The new quote-aware scanner (`lib/css.js` `stripComments` and
`parseStylesheet`) treats a newline inside a quoted string as whitespace and keeps the
string open to EOF. Probe: appending

```css
.c::before { content: "unterminated
 color: #ff0000; }
```

to `sound-run`'s `panel.css` yields `claimed L3, verified L3`, verdict PASS, **exit 0**.
But CSS Syntax 3 §4.3.5 ends a string at an unescaped newline (bad-string token, parse
error, recovery continues): a browser drops the `content` declaration and **applies
`color: #ff0000`** in that block. So the checker issues a clean exit-0 L3 certificate
beside a raw literal the ratified system really contains — the `content: "}"` hole's
sibling, reachable only through malformed CSS.
**Why it matters.** The one scanner family this cycle hardened still has an exit-0
evasion, and unterminated strings occur in real stylesheets precisely because browsers
tolerate them — the page "works", the checker passes, the token layer is bypassed.
Narrower than the round-02 hole (invalid input required), hence MAJOR not CRITICAL.
**Fix.** Apply the bad-string rule in both quote loops: an unescaped newline ends the
string and scanning resumes at the next character, so following text parses as
declarations. Fixture pair: the probe above (fail) and a properly terminated block
(pass). Optionally report a string reaching EOF unterminated as NEEDS_REVIEW.

### [MAJOR-2] The frontmatter parser never enforces the subset's scalar-quoting rules; the corruption is silent in the waivers path

**What.** §2 rule 5's scalar grammar requires quoting around `, [ ] { } #` and edge
spaces, and bans trailing comments ("never trailing"); §2 rule 6 says a file whose
frontmatter leaves rule 5 MUST be reported as violating it. The parser
(`lib/registry.js` `parseScalar`) implements none of that: an unquoted ` # ` becomes
literal value text. Loud case: `kind: DESIGN-BRIEF # x` fails as an unknown kind.
Silent case, probed: `waivers: WAIVERS.md # the waiver file` → the corrupted path
resolves to nothing, `waivers` stays null, the run continues — **verified L3, PASS,
exit 0** on a file §2 rule 6 says MUST be reported. A typo'd `waivers: WAVERIS.md`
degrades the same way: a contract naming a waiver file this run cannot read is treated
as a run with no waiver file, silently, whenever no `waived` answer forces the read.
**Why it matters.** `l1-frontmatter-parses` claims "every file declaring the spec
parses as the canonical frontmatter subset (§2 rule 5)"; it in fact only catches files
that fail *this* parser, and this parser accepts syntax the subset forbids — the
certificate asserts a conformity the checker never tested. The silent branch retires a
§8 control (the single named waiver file) without a word.
**Fix.** Two parts. (a) In `parseScalar` (or a rule-5 lint pass over raw lines): reject
an unquoted scalar containing ` #`, leading/trailing whitespace, or a trailing comment
shape — report as `pds-check:l1-frontmatter-parses` with the line. Quoted forms stay
legal. (b) A contract whose named waiver file resolves to no readable file is reported
(NEEDS_REVIEW at minimum), not read as "no waivers". Fixtures: the `waivers:` probe
above (fail), a quoted `waivers: "WAIVERS.md"` (pass), a missing waiver file (reported).

### [MAJOR-3] Same-tier alias edges: the doctrine says "points down", the checker tests "not up"

**What.** §3 G3: "a reference points down that order and a primitive holds a value,
never another primitive." The checker (`engine.js:1326`) raises only when `to > from`,
plus the primitive→primitive case — so `semantic.muted → {semantic.text}` (probed) and
`component.x → {component.y}` pass silently. Read strictly, "points down" forbids
same-tier edges, and the primitive clause is exactly the floor statement of that
reading (a primitive has nothing below it, so it may reference nothing); the checker's
reading ("not up, except primitive→primitive") leaves the primitive clause unmotivated —
why ban peer references at the floor but nowhere else?
**Why it matters.** If strict descent is meant, every L3 certificate over a token file
with same-tier aliases is issued against a G3 violation the checker declined to see —
the same "certificate on evidence never obtained" shape, one conjunct small. If loose
descent is meant, the doctrine's word "down" misleads every reader, and two of three
reviewers independently read it as strict. Either way, doctrine and checker must be
brought to one meaning before 1.0.0 ships.
**Fix.** The owner rules on the meaning, then: strict — change the comparison to
`to >= from` (the primitive case folds into it), add same-tier fail fixtures at both
lower tiers; loose — rewrite G3's clause as "a reference never points up that order;
same-tier references are permitted except between primitives". One word of doctrine or
one character of code, plus a fixture.

### [MAJOR-4] §1's per-agent token sidecars are unenforced, and the sound-run fixture violates them — hollowing G2's graft check

**What.** §1's DIVERGE row names the artifacts `round-01/<agent>.md +
<agent>.tokens.json`, and §2 DIRECTION's `tokens` is "this direction's token file". The
shipped `sound-run` fixture has both DIRECTIONs naming `tokens: ../tokens.json` — one
shared file at the run root — and verifies L3. The new `PROCESS_HOMES` check pins the
markdown half of the row and never the token half. With a shared file, G2's graft
re-expression test ("re-expressed in a token the winner already declares") is vacuous:
the loser's token world and the winner's are the same index, so every graft target
exists by construction. The conformant example the suite ships teaches the shape that
defeats the gate.
**Why it matters.** This is round-02's process-order hole one row over: the level
certifies §1's mapping while §1's mapping is not fully checked, and here the shipped
fixture itself is the violating run, so every consumer learns the wrong shape. The
digest conjunct (R-1) still works; the re-expression conjunct is the casualty.
**Fix.** (a) In the process-order check, resolve each DIRECTION's `tokens` and require
`round-01/<agent>.tokens.json` with the same `<agent>` as the direction's file (a run
that legitimately declares no per-direction tokens should fail or go unverified, not
pass). (b) Rewrite the sound-run fixture with per-agent sidecars
(`round-01/claude-1.tokens.json`, `round-01/codex-1.tokens.json`), re-pinning the
VERDICT's `tokens-digest`. Negative fixture: the shared-`../tokens.json` shape.

### [MINOR] §4 rule 4's second critique round has no §1 home, so a lawful run cannot verify L2

**What.** §4 rule 4: "Exactly one critique round. A second round needs an explicit
Decider instruction and a recorded reason." §1 maps CRITIQUE to `round-02/<agent>.md`,
"exactly one round", and the new `PROCESS_HOMES` check pins every CRITIQUE artifact to
`round-02`. Probe: a Decider-instructed second-round CRITIQUE filed at
`round-03/codex-1.md` fails `pds-check:l2-process-order` ("filed outside §1's
mapping") — the run can never verify L2, so the permitted exception is unusable in any
checked run.
**Why it matters.** A permission the doctrine grants and the checker forbids is the
unverifiable-MUST shape inverted: users who follow §4 rule 4 get an unexplained
process-order failure. Small, but it sits in the mapping table my lens owns.
**Fix.** One sentence: either §1 names the second round's home ("a Decider-instructed
second round files its CRITIQUEs in `round-03`", and `PROCESS_HOMES` accepts it when
such an instruction is recorded) or §4 rule 4 states the second round also files under
`round-02`. Fixture for the permitted shape.

### [NIT] Alias-direction group names: checker accepts plurals, doctrine names singulars

`TIER_GROUPS` maps `primitives`/`semantics`/`components` as well as the singulars, but
PDS §3 G3 says "where a document names a group `primitive`, `semantic` or `component`"
— a document using `primitives` names none of those, so the conjunct is vacuous per the
doctrine's own text, yet the checker enforces it (probed: `primitives.ink` →
`{primitives.slate}` fails `l3-alias-direction`). The checker is stricter than the
doctrine it implements. Fix either side: G3 gains "(singular or plural)" or
`TIER_GROUPS` drops the plurals. One line.

### [NIT] The checker's "will not take on trust" table predates cycle 4

`parley-design-check/SKILL.md:137` lists five recomputed conditions; cycle 4 added at
least five more (the deduplicated rotation, G1's brief-axis count, §1's locations, the
G3/G4 refutations, waiver ownership). They *are* documented — the level-obligations
table two sections later was updated and is accurate — but the trust table reads as the
complete list of what is recomputed rather than believed, and it no longer is. Add the
rows or retitle the table.

### [NIT] PDS.md has 2 bytes of headroom against its early-warning threshold

25,598 of 25,600 bytes (the binding total has 172 spare: 65,364/65,536). The thresholds
are documented as early-warning, so nothing is wrong yet — but the next normative
sentence added to PDS.md forces compression or a third rebalance. Recorded so the next
fix-up cycle plans for it rather than discovering it in CI.

## Cross-reviewer responses

### codex-1 (round-01 findings)

- **CRITICAL AF-1 (forgeable certificates): GENUINELY CLOSED.** In round-02 I reproduced
  all six paths on the pristine tree; this round I re-ran every one against `17f6619`
  ((a)–(g) above) and each now fails with a named obligation and exit 1. The controls
  (sound-run baseline, legitimate waiver) still pass.
- **CRITICAL AF-2 (author of waived work counter-signs): GENUINELY CLOSED.** Rejected
  with the right message; the positive control (independent corroborated signer) still
  suppresses, so the ownership derivation does not over-fire; the uncorroborated-signer
  conjunct behaves.
- **MAJOR AF-3: closed** (exit 4 re-verified this round). **MAJOR AF-4: closed** (and
  extended: `l1-artifact-kind` now catches the typo'd-kind case codex-1 raised under
  AF-1). **MAJOR AF-5, AF-6: closed** (round-02 probes; fixtures remain in the suite).
- **MAJOR AF-10 (D-1): closed by reasoning; I maintain my accept.** The binding total
  held at 65,364/65,536 through a fourth cycle, and the registry content is provably
  unchanged (identical YAML fences), so the rebalance bought normative text, not
  weakened rules.
- Round-02 MAJOR (fast path): **closed** — doctrine and dispatcher now state the fast
  path is outside Parley and caps at L1. Round-02 MINOR (alias direction): **closed as
  shipped** — defined, enforced, fixtured — but its resolution is what exposed MAJOR-3
  (same-tier) and the plural NIT above.
- Round-03 findings: all three **independently reproduced by me** and carried as
  MAJOR-2, MAJOR-3 and MAJOR-4. I join the substance of codex-1's round-03 review; on
  severity I block only on CRITICALs, so these four MAJORs are my reservations rather
  than a block.
- MINORs: heroicons double-count — closed per hermes-1's round-02 probe and the
  `one-package` fixture; not re-run by me this round. Unnumbered normative requirements
  — closed (§0.5 rule 3 + §2 rule 4, re-confirmed in this round's full diff read).

### hermes-1 (round-01 findings)

- **MAJOR AF-7 (banned-slop signature undefined): GENUINELY CLOSED.** Definition intact
  at RULES.md:278–287 after the compression (I diffed the registry fences and read the
  prose), `banList()` still derives `slop`@`T0`, and I watched R-2's observed-vs-recorded
  comparison fire on my own probe this round.
- Round-02 MAJOR (`waivers: ""` crash): **closed** at all three path sites, with
  `tokens: ""` on the winner degrading to an honest `unverified` + exit 4 — its open
  question 1 (other path fields) is answered: `waivers`, contract `tokens` and DIRECTION
  `tokens` are the only path-bearing fields and all three are guarded (`trim()` +
  `isFile()`).
- Round-02 MINOR (alias direction): **closed** as above. Its round-02 open question 3
  (should G4 be recomputed) is answered by `gateRefutations` — probed (e).
- Round-03 MAJORs: the alias same-tier and frontmatter-`#` findings reproduce under my
  own commands and are carried as MAJOR-3 and MAJOR-2; the sound-run-sidecar finding
  reproduces and is MAJOR-4. Its round-03 MINOR (gate-string separator drift "still
  open from round-01") I do **not** carry: §3's preamble labels the strings "canonical
  message shapes, not literal output" and licenses tool variation, which hermes-1's own
  round-02 review accepted as the correct resolution; a colon-vs-em-dash difference
  inside that licence is not a defect.
- Round-02 NITs (D-2 count, `core:contrast-applied` omitted from D-2's list): accurate,
  cosmetic, unchanged; no action needed.
- Round-01 MINORs (SC sourcing, SC 1.4.4/1.4.12 informative, gate-string preamble,
  Keys-table tier form, DTCG date pin): verified closed in round-02; the diff this cycle
  touches none of them.

## Open questions

1. **The level/verdict axis split, one last corner.** An open, honestly-recorded
   `slop` violation against the winner (critique never cited it, verdict never answered
   it, waiver rejected) yields `verified L3` beside verdict VIOLATION, exit 1. I now
   read this as *consistent* with §9: G1's ban-list conjunct is a sharing test, not a
   cleanliness test, and no level obligation is "no open slop findings" — the R-2
   mechanism already forces the finding into the recorded signature, so it cannot be
   hidden. Confirm that reading is intended; if not, `gateRefutations` needs a G1 arm.
   (Supersedes my round-02 question 3: the G3/G4 arms landed and work.)
2. **Alias direction semantics (MAJOR-3).** Strict descent or not-up? The owner should
   rule before the checker and doctrine are aligned; two reviewers read "down" as
   strict, the implementer read it as loose.
3. **Second-round critique home (MINOR).** A doctrine question as much as a checker one
   — the owner should pick where §4 rule 4's exception lives before the checker can
   accept it.
