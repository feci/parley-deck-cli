---
agent: kimi-1
idea: parley-design-skills
review-round: 02
date: 2026-07-28
reviewed-commit: 8fc3a18
---

## Summary

❌ BLOCK. The reviewed deck commit `8fc3a18` names implementation HEAD `8ebd8f7`
(branch `parley-design-skills` in `parley-deck-skill`). I pinned every adjudication to a
pristine `git archive` extraction of `8ebd8f7` after discovering the worktree was being
mutated mid-review (uncommitted in-flight fixes to `lib/css.js` +28, `lib/engine.js`
+421/−60, fixture renames `brief.md` → `DESIGN-BRIEF.md` — evidently the start of a fix
for codex-1's round-02 findings; HEAD is still `8ebd8f7`, so none of it is reviewable
yet). On the pristine reviewed tree, codex-1's two CRITICALs reproduce in full with my
own commands, and hermes-1's new MAJOR reproduces. The doctrine layer under my lens is in
good shape — PDS.md conforms to the ratified protocol everywhere I checked — but the
checker still issues clean conformance certificates on evidence it never obtained, and
that is what a conformance checker exists not to do.

Agreed-fix disposition (each verified below, none believed):

- **AF-1: not closed.** Six independent paths to false L1/L2/L3 certificates reproduce on
  `8ebd8f7`, four of them with `PASS` + exit 0. (CRITICAL-1.)
- **AF-2: not closed.** The author of the waived work can still counter-sign; reproduced
  as `PASS`, exit 0, `waived 1`. (CRITICAL-2.)
- **AF-3: closed.** `check NOTICE.md` → verdict `UNJUDGEABLE`, exit **4**; help and the
  checker's SKILL.md document the code.
- **AF-4: closed.** Canonical frontmatter subset is ratified (§2 rule 5) and obeyed by
  all eight examples; an unknown cited rule id passes through as `UNJUDGEABLE` with a
  clear message and fails the `l2-cited-rules` obligation rather than vanishing.
- **AF-5: closed.** A reduced-motion block setting only `color: red` no longer counts as
  coverage; `core:motion-without-reduced-path` fires.
- **AF-6: closed.** Split `outline: none` / `button:focus-visible { outline: … }` raises
  no `core:focus-indication` finding.
- **AF-7 (hermes-1): genuinely closed.** RULES.md:282–291 defines the ban list (derived:
  `slop` class at `T0 ARTIFACT`), the signature, and the sharing test; the checker's
  `banList()` (engine.js:389–392) implements exactly that derivation; §3 G1 cites it.
- **AF-8 (mine): genuinely closed.** §3 rule 1 and the G1 message shape carry all three
  C7 conjuncts — ban list, `core:category-guessable`'s category-plus-avoidance test,
  recorded human ratification with a brief-specific reason — else `ABSTAIN`.
- **AF-9 (mine): partially closed.** The doctrine side landed (`run-id` in §2's
  DESIGN-BRIEF table and example; §4 rule 2 fully pins the hash semantics) and the
  rotation is genuinely recomputed (a wrong `assigned` fails L2 with "the rotation is
  recomputed from the brief, never taken on trust"; a verbatim-recorded `assigned` plus a
  one-line `declined` is the only legal decline, correctly). Two residual holes —
  invented non-brief axes and duplicate primary values — are codex-1's round-02 findings
  and reproduce under CRITICAL-1.
- **AF-10: closed by accepted reasoning.** I maintain my round-01 accept of D-1 (below);
  codex-1 withdrew its block.

R-1..R-4 spot-checks on pristine: token drift after ratification fails G2 on the
`tokens-digest` comparison (R-1 ✓); ghost signers `nobody-1`/`ghost-2` are rejected
("not a participant this run records"), finding unsuppressed (R-3 ✓); recusal is
path-anchored and `recusal-not-anchored` is printed, never a silent pass (cycle-3
non-claim ✓, visible in every run I made).

## What I verified (commands run, and their result)

All checker probes run twice where it mattered: against the live worktree and against a
pristine extraction (`git archive 8ebd8f7 | tar -x -C /tmp/pds-pristine`). Results below
are pristine unless marked worktree.

- `git rev-parse HEAD` / `git status` in the skill repo: HEAD `8ebd8f7`, branch
  `parley-design-skills`, worktree **dirty** (uncommitted in-flight fixes — see Summary).
  Deck repo HEAD `8fc3a18`, matching the reviewed-commit.
- `npm test` in the pristine extraction: **197 passing, 0 failing** — IMPLEMENTATION.md's
  cycle-3 count holds at the reviewed commit.
- Baseline `check --level L2` on `sound-run`: `claimed L2, verified L2`, `verdict PASS`,
  exit 0, 13 rules listed `UNJUDGEABLE` (never skipped), `recusal-not-anchored` printed
  for hermes-1.
- codex-1's AF-1 probes, rebuilt by me on pristine fixture copies under /tmp:
  - rounds swapped (DIRECTIONs→round-02, CRITIQUE→round-01): `verified L2`, PASS, **exit 0**;
  - invented axis `texture` + only one differing brief axis: `verified L2`, PASS, **exit 0**;
  - primary values `[flat, flat, layered]`, both proposers assigned `flat`: `verified L2`,
    PASS, **exit 0** (§4 rule 2's "materially distinct positions" count not enforced);
  - raw `#ff0000` literal with G3 recorded `pass`: violation found, verdict VIOLATION,
    exit 1 — but the certificate line still reads `verified L2` with empty `unmet`;
  - AUDIT recording G4 `pass` beside the open violation: same, `verified L2`;
  - `.trap::before { content: "}"; color: #ff0000; }` appended to panel.css at L3: the
    literal is lost by the parser — `verified L3`, PASS, **exit 0**;
  - `kind: TYPO-BRIEF` file at L1: never inspected, `verified L1`, PASS, **exit 0**.
  I also re-ran codex-1's own surviving probe dirs (/tmp/pds-codex1.6XgLyH) — same
  results on the pristine checker.
- codex-1's AF-2 probe, rebuilt: waiver scoped at `round-01/claude-1.md`,
  `granted-by: codex-1`, `counter-signed-by: claude-1` (the author of the scoped file) →
  `WAIVED`, verdict PASS, **exit 0**.
- hermes-1's MAJOR, rebuilt: CONTRACT with `waivers: ""` → `Error: EISDIR` from
  `loadWaivers` (engine.js:1451: `path.resolve(dir, "")` → the directory itself) — crash,
  no report. Extension I probed myself: `tokens: ""` crashes identically at
  `readTokens` (artifacts.js:72 via engine.js:1235). Same defect class, second site.
- AF-3: `check NOTICE.md` → exit 4, `UNJUDGEABLE`. AF-4: critique citing
  `project:unknown-thing` → `UNJUDGEABLE` pass-through with remedy. AF-5/AF-6 probes as
  above. Wrong `assigned` → `l2-assignment` fails L2. R-1/R-3 probes as above.
- Refusal (C4): checker copied alone to /tmp, run on a plain file → stderr refusal naming
  the absent registry, exit 3, structural checks still ran.
- Doctrine (my lens, full read of PDS.md): §1 mapping table row-for-row C1; all eight
  artifact kinds keep the identical four-part shape (purpose line, rationale paragraph,
  required-fields table, minimal example — checked each); §3 G1–G4 with the A1 two-axis
  test, restored C7 conjuncts, `tokens-digest` G2, message shapes labelled "canonical
  message shapes, not literal output"; §4 ritual with the pinned assignment formula,
  decline valve, A3's unattended stop near-verbatim, C12 fast path; §5 roles/recusal/
  no-self-scoring/800-word caps/declared degradation/Phase-6-reviewer authorship (the
  carried-open question, closed as FINAL.md required); §6 ordinal number-plus-word tiers,
  engine names banished; §7 class authorities; §8 waiver fields incl. the
  neither-grantor-nor-author rule; §9 cumulative L1–L4 with the L2 row scoped to crossed
  transitions (my round-01 MINOR); §10–§12 extension/versioning/changelog all present and
  maintained.
- Registry and budgets: `shasum -a 256 RULES.md` → `f0c38eed1b8d` = PDS frontmatter
  digest. Sizes 6681/24909/23572/10124 = **65286 ≤ 65536**; RULES.md headroom now 1004
  bytes (my round-01 MINOR addressed). Keys-table tier row now documents the bare-ordinal
  form (my round-01 NIT fixed). Stale example ids from my round-01 MINOR: zero matches.
  No `RULES.md` under `parley-design-check/`; zero placeholder strings in shipped files;
  capability line "18 detectors over 18 rule ids, generated from lib/detectors".
- Worktree-only (informational, not part of the reviewed commit): the uncommitted fixes
  catch all six AF-1 probes, reject the author-counter-signed waiver with the right
  message, and survive `waivers: ""`. Directionally correct; uncommitted and unreviewed.

## Findings

### [CRITICAL] AF-1 remains open: clean L1/L2/L3 certificates on evidence never obtained

**What.** Six reproduced paths on `8ebd8f7` (commands above): (a) DIRECTIONs filed in
round-02 and the CRITIQUE in round-01 verify L2 — process order tallies artifact kinds
and never checks §1's normative homes; (b) an axis the brief never declared supplies
G1's second difference — the two-axis count runs over the union of DIRECTION `positions`
keys, not the brief's declared axes §3 names; (c) primary values `[flat, flat, layered]`
assign both proposers `flat` and verify — §4 rule 2's "at least as many materially
distinct primary positions as there are proposers" is not enforced as *distinct*, so the
assignment that exists to guarantee divergence guarantees none; (d) a G3-recorded `pass`
and (e) a G4-recorded `pass` survive the checker's own contradictory findings — the
level certifies on the recorded word while the report lists the violation (exit 1 saves
CI here, but the certificate lies); (f) `content: "}"` hides a following raw colour from
the T1 scanner — verified L3, PASS, exit 0; (g) a `spec: PDS/1.0` file with an unknown
`kind` is silently uninspected — verified L1, PASS, exit 0.
**Why it matters.** This is the exact failure AF-1 was accepted to close, and four of
the seven paths exit 0: a gate that reads as protected and is not. IMPLEMENTATION.md's
cycle-1 record claims levels are modelled "so a level cannot be certified on evidence
never obtained" — that claim is not true at the reviewed commit.
**Fix.** (a) Validate artifact paths against §1's mapping (brief at Phase 0 as
`DESIGN-BRIEF.md`, directions in round-01, critique in round-02). (b) Count G1
differences only over brief-declared axes, and validate each declared position against
its axis's enumeration. (c) Enforce the materially-distinct count (deduplicate before
rotating, or reject duplicate primary positions in the brief). (d)/(e) When the run
crossed G3/G4, feed relevant rule outcomes into the level obligation, not just the
recorded outcome word. (f) Make the CSS scanner quote/escape-aware before treating
`{};:` as structure. (g) Any file declaring `spec: PDS/1.0` with an unknown `kind` is an
L1 violation, never `not-inspected`. Negative fixture per path. The in-flight worktree
already implements most of this; land it as a commit and re-review.

### [CRITICAL] AF-2 remains open: the author of the waived work can counter-sign

**What.** Reproduced: a waiver scoped at `round-01/claude-1.md`, granted by codex-1,
counter-signed by **claude-1** — the author of the scoped work — suppressed the finding;
verdict PASS, exit 0. `waiverProblem()` establishes that grantor and signer are distinct
roster members, but never compares the signer with ownership of the scoped work, which
§8 rule 2 requires ("neither the grantor nor an author of the waived work").
**Why it matters.** Self-approval through a colleague's grant is the same unilateral
suppression C13 exists to prevent, one indirection removed. A `quality` finding — the
class that may block a build — can be retired by the person who caused it.
**Fix.** Resolve an artifact scope to its path-author — the checker already trusts file
paths over declared identity for recusal (engine.js:628–633); reuse that mechanism — and
reject a counter-signer matching it. For scopes whose authorship cannot be established
from the run's artifacts, leave the finding unsuppressed, as §8 rule 2 already requires.
Fixtures: signer-is-scope-author; signer known only via a self-authored artifact.

### [MAJOR] `waivers: ""` (and `tokens: ""`) crash the checker instead of reading as empty

**What.** hermes-1's MAJOR, reproduced: a CONTRACT with `waivers: ""` crashes with
`EISDIR` — engine.js:1451 tests `typeof === "string"` (true for `""`), resolves the
empty string to the contract's own directory, `existsSync` passes, and `loadWaivers`
reads a directory. I extended the probe: `tokens: ""` crashes identically in
`readTokens` (artifacts.js:72 via engine.js:1235). The crash is an uncaught exception —
no report, and not even the documented exit 2.
**Why it matters.** §2 rule 5's canonical subset lists `""` as a valid empty value and
§2 rule 3 says "empty is not absent": a run with no waiver file and no token-layer
override is a legitimate state the spec permits and the checker kills. A CI gate sees a
checker failure on a conformant input.
**Fix.** Guard the empty string at both path-resolution sites (`… === "string" &&
.trim() !== ""`), or check `statSync(p).isFile()` before reading; treat the guarded
empty as absent. Fixtures for `waivers: ""` and `tokens: ""` running to a normal report.

### [MINOR] The fast path's relationship to Parley's non-solo rule is never stated

**What.** codex-1's round-02 MAJOR; I downgrade to MINOR. §4 rule 8 and SKILL.md
prescribe a one-agent fast path while §0.3 rule 2 gives `COOPERATION.md` process
supremacy, and FINAL.md defers any carve-out to a separate meta-protocol idea. The
conflict technically resolves for COOPERATION.md (so the fast path cannot be a Parley
process), but neither file says that, so a user can read the add-on as licensing a solo
run that claims Parley verification.
**Why it matters.** The audit trail cannot distinguish a valid out-of-Parley fast path
from an invalid solo Parley run; the ambiguity is in the document my lens checks.
**Fix.** One sentence in SKILL.md and §4 rule 8: the fast path runs outside a Parley
Deck workflow and its result MUST NOT be recorded as Parley-verified; a solo carve-out
inside Parley goes through the deferred meta-protocol change.

### [MINOR] "Alias direction" is required by FINAL.md, absent from G3, unimplemented

**What.** FINAL.md (line 111) names alias direction part of L3 token integrity; §3 G3
lists alias resolution and cycles but not direction; the checker verifies resolution and
cycles only. Raised by hermes-1 in round-01, accepted into the fix-up pass, never
landed; both round-02 reviewers concur it remains.
**Why it matters.** A binding acceptance criterion the doctrine dropped and the checker
cannot check — the unverifiable-MUST shape again, one size smaller.
**Fix.** Define the permitted direction in G3 (e.g. an alias chain MUST NOT traverse
from `semantic` groups to `primitive` ones), implement it with a negative fixture; or
record an erratum striking it from FINAL.md. Either closes the gap; the current tri-state
does not.

### [NIT] None raised.

(hermes-1's two round-02 NITs — the D-2 count omitting `core:contrast-applied` — are
accurate and cosmetic; no action needed beyond the record.)

## Cross-reviewer responses (round-01 findings)

### codex-1

- **CRITICAL AF-1 (L2/L3 falsely certified): NOT genuinely closed** — six reproduced
  paths above, including two I initially mis-adjudicated: my first css-string and
  unknown-kind probes ran against the mid-review-patched worktree and did not reproduce;
  re-run against pristine `8ebd8f7` they reproduce exactly as codex-1 reported, as do
  codex-1's own surviving probe inputs. Apologies for the detour; the finding stands in
  full.
- **CRITICAL AF-2 (self-counter-signature): NOT genuinely closed** — the equal-identity
  case is fixed (hermes-1's probe and the rejection message confirm), but the
  author-of-the-waived-work case reproduces as a clean PASS.
- **MAJOR AF-3: closed** (exit 4 verified). **MAJOR AF-4: closed** (subset ratified;
  unknown ids pass through as `UNJUDGEABLE`). **MAJOR AF-5, AF-6: closed** (own probes).
- **MAJOR AF-10 (D-1): closed by reasoning**; codex-1's withdrawal is consistent with my
  round-01 position.
- MINORs: heroicons double-count — hermes-1 verified and the `one-package` fixture
  exists; I did not re-run it. Unnumbered normative requirements — closed: §0.5 rule 3
  plus §2 rule 4 give the tables normative force through a numbered rule, verified in my
  full read.

### hermes-1

- **MAJOR AF-7 (banned-slop signature undefined): genuinely closed** — doctrine side
  (RULES.md:282–291, cited from §3 G1) and checker side (`banList()` derives `slop`@`T0`)
  both verified.
- MINORs: SC 2.3.3/2.2.2 sourcing, SC 1.4.4/1.4.12 moved to informative, gate-string
  preamble, Keys-table tier form, DTCG date pin — the four I spot-checked (preamble,
  Keys table, ban-list definition, colour-space) are as hermes-1 reported; bare-hex
  colour without `colorSpace` now fails L3 (own probe). Gate strings now honestly
  labelled shapes rather than claimed literals — my round-01 MINOR is closed the right
  way. Alias direction: see my MINOR above — the one residual from hermes-1's round-01
  set.
- **hermes-1's round-02 claim that "all cross-reviewer findings … are genuinely closed"
  is wrong** on AF-1 and AF-2 (and its own signoff's "AF-9 fixed" needs the two residual
  holes above) — its own new MAJOR, however, reproduces, and its alias-direction MINOR is
  correctly carried.

## Open questions

1. **Roll-up semantics (carried from my round-01).** The sound-run reports
   `verdict PASS` with 13 of the in-scope rules `UNJUDGEABLE`. Per-rule honesty is
   intact, but a top-line PASS on a mostly-unjudged run can read as an all-clear. Should
   the verdict line carry judged/total, or distinguish "PASS, all judged" from "PASS
   with unjudged rules"?
2. **Review hygiene.** The worktree was patched mid-review, and half my early probes
   silently changed meaning until I pinned a `git archive` of the reviewed commit. Should
   the review handoff require a clean worktree (or a tag) so "reviewed-commit" names an
   unambiguous tree? The in-flight fixes also mean a cycle-4 re-review should be able to
   close CRITICAL-1/2 quickly once committed.
3. **Certificate-vs-verdict split.** Even where exit codes save CI (paths d/e), a
   `verified L2` line beside an open gate-relevant violation is a certificate a human
   will quote without the exit code. Should any open violation relevant to a crossed
   gate cap the level line at `not verified`, independent of the verdict roll-up?
