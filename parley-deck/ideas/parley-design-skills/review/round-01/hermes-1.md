---
agent: hermes-1
idea: parley-design-skills
review-round: 01
date: 2026-07-28
reviewed-commit: 726c024
---

## Summary

The implementation is structurally sound: both add-ons exist in the right shape,
the checker refuses rule checks without a registry (exit 3, verified), the
capability declaration is generated from detector modules (verified, 18 detectors
over 18 rule ids), the WCAG 2.2 contrast math is correct (verified against known
values: 21:1 black-on-white, 4.00 red-on-white, 4.48 #777-on-white), the
registry digest matches (1fbe071e1222, verified), all 158 tests pass, and no
placeholder text ships in any design file. The doctrine files are within the 64
KiB total (64723 bytes, verified). DTCG token handling — $value leaves, group-
level $type inheritance, {path.to.token} alias syntax, $extensions — is faithful
to the DTCG editor's draft.

My lens is the doctrine: whether RULES.md and WEB-ANNEX.md say true, decidable
things; standards fidelity (DTCG, WCAG 2.2); whether any rule is unfalsifiable,
mis-classified, or on the wrong evidence tier. I found one MAJOR issue (an
unfalsifiable G1 condition that the spec mandates but never defines), several
MINOR issues (a WCAG SC source mismatch, blocking numbers without rules, a
colorSpace gap in L3, alias direction not verified), and a few NITs.

D-1 (per-file budget rebalance): I accept it. C3 makes 64 KiB binding, and the
per-file split was one participant's proposal. The total is held (64723 ≤ 65536).
However, RULES.md sits 30 bytes from its 24 KiB ceiling — the registry cannot
grow without blowing the per-file budget, which contradicts C3's stated rationale
("leaves the rule registry room to grow past thirty entries"). The test enforces
24 KiB for RULES.md, so adding even one rule will require either compressing
existing prose or rebalancing again. This is a live constraint, not a defect.

## What I verified (commands run, and their result)

1. `npm test` from the skill repo root — 158 tests, 0 failures, 0 skipped
   (duration 1362 ms).
2. `node addons/parley-design-check/bin/check.js --help` — prints usage, exit 0.
3. `--registry /nonexistent/RULES.md /tmp/test-tokens.json` — refuses rule
   checks on stderr, exit 3, structural checks still ran. Confirms C4 refusal.
4. Capability JSON (`--json --registry /nonexistent`) — 18 detectors, 18 rule
   ids, generated from lib/detectors. Confirms capability is generated, not
   hand-maintained.
5. Registry load via `loadRegistry()` — 30 rules total (19 core + 11 web), 0
   duplicate ids, 0 missing records for declared ids, 0 warnings. Confirms
   registry parses clean.
6. `loadDetectors()` — 18 detector modules, all with required fields, all at
   T0 or T1. Cross-referenced with registry: every detector rule exists in the
   registry; 9 rules enforced by check/both have no detector (5 at T0/T1
   reported "no detector", 4 at T2 reported "tier above this checker").
7. Contrast ratio computation: `toSrgb` + `contrastRatio` — 21:00 for
   #000/#fff, 4.00 for #ff0000/#fff, 4.48 for #777777/#fff, 4.54 for
   #767676/#fff. All match published WCAG reference values.
8. Registry-digest mismatch: edited PDS.md to declare `deadbeefdead`, ran
   checker — report notes contain "registry-digest mismatch: the spec declares
   deadbeefdead, the registry file computes 1fbe071e1222". Restored PDS.md
   afterward. Confirms §11 rule 3 drift detection.
9. File sizes: SKILL.md 6656, PDS.md 22389, RULES.md 24546, WEB-ANNEX.md 11132.
   Total 64723. All within test budgets. RULES.md is 30 bytes from its 24576
   ceiling.
10. `toSrgb('#ff0000')` returns channels without error — a bare hex $value
    passes L3 color checks without colorSpace being declared.
11. Every detector has a fail/pass fixture pair (18/18, verified by directory
    listing).
12. `grep` for "banned" across all doctrine files: appears only in PDS.md G1
    and in the consensus/FINAL specs — never defined in the doctrine itself.
13. Checker G1 implementation: tests 2-axis difference and duplicate
    Signatures; does not test banned-slop signature (not found in engine.js).
14. Checker L3: tests alias resolution, alias cycles, color computability; does
    not test alias direction.
15. WCAG 2.2 SC references in WEB-ANNEX: 1.4.3, 1.4.10, 1.4.11, 2.5.8 — all
    verified against the published standard. SC 1.4.4 and 1.4.12 appear in the
    blocking numbers table but have no corresponding rule.

## Findings

### [MAJOR] G1 "banned-slop signature" is an unfalsifiable MUST

PDS.md §3 G1 says the gate "MUST fail if ... two directions share a banned-slop
signature." Consensus C7 and A1 both say "the banned-slop-signature and
duplicate-Signature checks are retained." FINAL.md item 5 says "plus
banned-slop-signature and duplicate-Signature checks."

But "banned-slop signature" is never defined anywhere in the doctrine. It is
not a rule in RULES.md. It is not a list of patterns. It is not a matching
criterion. The word "banned" appears exactly once in the entire doctrine
(PDS.md:298), and only to name the check, not to define it. Consensus C5 gives
three examples ("a banned font, a purple-to-cyan gradient, the icon-tile feature
card") but these are illustrative, not normative, and none appear as registry
rules or annex entries.

The checker's G1 implementation (engine.js:310-343) tests two of G1's three
conditions: the 2-axis difference and the duplicate-Signature check. The third
condition is silently absent. This is not a missing detector — it is a missing
definition. A facilitator computing G1 by hand cannot know what constitutes a
"banned-slop signature" any more than the checker can.

This is the most serious doctrine issue because G1 is a gate: a failed set
MUST NOT proceed to critique. A MUST condition that nobody can evaluate is
either dead code (it never fires) or a trap (it fires on judgement nobody can
audit). The consensus explicitly retained this check twice (C7, A1), which
makes its absence from the doctrine a fidelity gap, not a deferral.

Concrete fix: either define "banned-slop signature" as a named list in
WEB-ANNEX.md (or a new annex), with the matching criterion (e.g. "a Signature
that names a pattern in the banned list"), and add a detector or an L2 check
for it; or strike the condition from G1 and record the removal as a spec
erratum. The current state — a MUST that cannot be satisfied — is the worst
option.

### [MINOR] WCAG SC 2.3.3 source does not match the rule's scope

`core:motion-without-reduced-path` cites `sources: [WCAG-2.2-SC-2.3.3]` and its
prose says "Any motion the interface starts on its own MUST have a declared path
for a user who has asked for reduced motion."

SC 2.3.3 ("Animation from Interactions", Level AAA) covers motion *triggered by
user interaction* — it says such motion can be disabled unless essential. The
rule covers motion "the interface starts on its own" — i.e., auto-started motion.
The more directly relevant SC is 2.2.2 ("Pause, Stop, Hide", Level AA), which
covers moving/blinking/scrolling information that starts automatically.

The rule's intent (reduced-motion path for auto-started motion) is sound and the
overlap with 2.3.3's vestibular-disorder concern is real. But the source citation
is imprecise: a reader checking the citation against the SC will find the scopes
do not align. Since the `sources` key is how the doctrine distinguishes
standards-derived thresholds from its own calibrations, a misattributed source
weakens the auditability the key exists to provide.

Concrete fix: add `WCAG-2.2-SC-2.2.2` to the sources list (the auto-started
motion SC), keeping 2.3.3 if the interaction-triggered case is also intended.
Alternatively, narrow the rule text to cover interaction-triggered motion if
2.3.3 is the sole intended source.

### [MINOR] WEB-ANNEX blocking numbers list SC 1.4.4 and 1.4.12 with no enforcing rule

The "Blocking numbers" table in WEB-ANNEX.md lists eight WCAG 2.2 success
criteria with the header "These are WCAG 2.2's published thresholds. They
block." Six of the eight have corresponding rules (1.4.3, 1.4.11, 2.5.8, 1.4.10,
2.4.7, 2.4.11). Two do not: SC 1.4.4 (text resize 200%) and SC 1.4.12 (text
spacing).

The section is normative for web (no `(informative)` label), and "They block"
is an imperative claim. But without a rule, there is nothing to block with — no
detector, no agent-judgement criterion, no error string. A web consumer reading
this annex will see two SCs declared as blocking and find no rule to cite, no
finding to emit, and no waiver to file.

This is not a defect if the intent is "these are reference numbers for context"
— but the section says "They block," not "These are reference numbers." The gap
between the claim and the enforcement is the problem.

Concrete fix: either add `web:text-resize` and `web:text-spacing` rules (both
at T2 RENDERED, which would make them UNJUDGEABLE in v1, which is honest), or
move SC 1.4.4 and 1.4.12 out of the "Blocking numbers" table into a separate
"Reference" or "(informative)" subsection. The current state implies enforcement
that does not exist.

### [MINOR] L3 does not verify colorSpace declaration for bare-hex colour tokens

PDS.md G3 says the gate "MUST fail on ... a colour token without colorSpace."
C11 (consensus) says "every colour token MUST declare a colorSpace and MUST be
computable to a displayable value." PDS.md §9 L3 says "colour tokens declare a
space and compute to a displayable value."

The checker's L3 conformance check (engine.js:443-458) calls `toSrgb()` on the
resolved value. `toSrgb()` treats a bare hex string (e.g. `"$value": "#ff0000"`)
as implicitly sRGB and returns channels — it never reports an error. So a colour
token whose $value is a bare hex string passes L3 without ever having declared a
colorSpace. The `colorSpace` check only fires when the value is an object that
lacks both `.hex` and `.colorSpace`.

This is a gap between G3's normative requirement ("MUST fail on a colour token
without colorSpace") and the checker's L3 verification. A bare hex is
computable, but it does not declare a colorSpace — and the spec says the
declaration is mandatory, not just the computability.

The counterargument is that a bare hex is unambiguously sRGB and therefore
implicitly declares its space. That is a reasonable reading of DTCG, but it is
not what PDS.md says, and it is not what C11 says. The doctrine's own words make
the declaration mandatory.

Concrete fix: in `checkL3`, after confirming a token is type `color`, check
whether the token (or its resolved target) carries an explicit `colorSpace` key
or a `$type` that implies one. If the resolved value is a bare string (hex),
emit a VIOLATION for "colour token declares no colorSpace" unless the token's
`$type` or a group-level convention supplies one. Alternatively, amend PDS.md
G3 to say "a colour token without colorSpace that is not computable to a
displayable value" — making computability the sole criterion, which is what the
checker already enforces.

### [MINOR] L3 does not verify alias direction

PDS.md §9 L3 says token integrity includes "alias direction." FINAL.md says the
same. The checker's L3 (engine.js:412-471) verifies that aliases resolve and
that no cycle exists, but it does not verify that aliases flow in one direction
(e.g. semantic → primitive, not reverse). A token graph where a primitive
aliases a semantic token would pass L3.

The DTCG editor's draft does not enforce alias direction — it permits any
reference — so "alias direction" is a PDS-specific constraint beyond what DTCG
requires. Since it is stated as an L3 requirement, its absence from the checker
is a conformance gap: a project claiming L3 cannot have its alias direction
verified.

Concrete fix: either define what "alias direction" means (e.g. "an alias chain
must not traverse from a group named `primitive` to a group named `semantic`"),
add a check for it in L3, and ship a fixture; or strike "alias direction" from
the L3 requirement in PDS.md if the constraint is not yet ready to enforce. As
with the banned-slop finding, a requirement that cannot be checked is worse than
no requirement.

### [MINOR] Gate error strings in PDS.md do not match the checker's

FINAL.md says PDS.md should carry "Gates G1–G4 with exact failure conditions and
the checker's error strings." PDS.md does include error string blocks for all
four gates, but they use an em-dash separator ("G1 DISTINCTNESS — ...") while
the checker's L2 conformance findings use a colon ("G1 DISTINCTNESS: ...").
The checker also uses `pds-check:l2-gate-g1` as the finding's rule id, while
PDS.md's strings use no id prefix.

This matters because PDS.md §3 says these are "the checker's" error strings,
implying they should match what the checker emits. A reviewer comparing the spec
against a checker report will see different strings and not know which is
authoritative.

Concrete fix: align the separators. Either change PDS.md's error string blocks
to use the colon form the checker emits, or change the checker to use the
em-dash form. The finding format (`rule-id — violation — remedy`) already uses
em-dash, so the checker's gate strings using colon is internally inconsistent
too.

### [NIT] Tier values in YAML use short form while the Keys table says full form

The RULES.md Keys table describes the `tier` key as accepting `T0 ARTIFACT`,
`T1 SOURCE`, `T2 RENDERED`, `T3 PIXEL` — the number-plus-word form that C6
mandates for prose. But every `pds-rule` YAML block in the shipped registry uses
the short form (`tier: T0`, `tier: T1`, etc.). The parser accepts both
(`normaliseTier` handles both), so this is not a functional issue, but the
registry contradicts its own Keys table. FINAL.md's YAML key spec shows the
short form, so there is a three-way tension: FINAL.md says short, C6 says
number-plus-word, RULES.md Keys table says full, actual blocks use short.

Concrete fix: either update the Keys table to show `T0 | T1 | T2 | T3` (matching
FINAL.md and the actual blocks), or update the blocks to use the full form. The
former is simpler and the parser already accepts it.

### [NIT] DTCG "2025.10" is not a published DTCG version

The doctrine references "DTCG `2025.10`" as if it were a versioned spec release.
The DTCG has not published formal versioned specifications with date-based
version strings; the URL in NOTICE.md (tr.designtokens.org/format/) points to
an editor's draft. "2025.10" appears to be a project-internal pin meaning "the
draft as of October 2025." This is a reasonable convention, but it implies a
formal DTCG release that does not exist, which could mislead a reader checking
the citation.

Concrete fix: add a parenthetical in NOTICE.md or PDS.md noting that "DTCG
2025.10" refers to the editor's draft as of that date, not a formal W3C
publication.

### [NIT] D-1 per-file budgets sum above 64 KiB

The test's per-file budgets (8 + 22 + 24 + 11 = 65 KiB) sum to 1 KiB above the
64 KiB total. This is not a bug — the total test catches the aggregate — but it
means the per-file budgets are not a tight partition and a file could grow into
the gap without tripping its individual test while still being caught by the
total. This is slightly unusual for a budget enforced "by a test, not by a
comment."

Concrete fix: either tighten the per-file budgets to sum to exactly 64 KiB
(e.g. 8 + 22 + 23 + 11), or add a comment in the test explaining that the total
is the binding constraint and the per-file numbers are early-warning thresholds.

## Open questions

1. The "banned-slop signature" check was retained twice by consensus (C7, A1)
   and appears in FINAL.md. Was its definition intended to live in WEB-ANNEX.md
   or a separate annex, and was it dropped during the author's-death recovery,
   or was it never written? If the latter, the consensus record says it was
   "retained" — but retained from what? There is no prior definition to retain.

2. G3 and G4 have error strings in PDS.md but the checker only recomputes G1
   and G2 (in L2). G3 is partially covered by L3 (alias/cycle/color) and by
   registry detectors (literal, token-declared-unused, token-used-undeclared).
   G4 (no open violations) is not checked at all by the checker — is this
   intentional (G4 is a review-time gate, not a checker gate), or should the
   checker at least verify that no VIOLATION findings remain unwaived at AUDIT?

3. RULES.md is 30 bytes from its 24 KiB ceiling. The consensus chose 64 KiB
   over 60 KiB specifically to "leave the rule registry room to grow past
   thirty entries." With 19 core rules + 11 web rules = 30, the registry is at
   the threshold the budget was designed for, and it cannot grow without
   rebalancing. Is the intent to rebalance the per-file split at the next rule
   addition, or to compress existing prose?
