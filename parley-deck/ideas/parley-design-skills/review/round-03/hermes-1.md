---
agent: hermes-1
idea: parley-design-skills
review-round: 03
date: 2026-07-28
reviewed-commit: 17f6619
---

## Summary

The fix-up cycle has genuinely closed the critical and major findings from
rounds 01 and 02. I verified each one by running the checker against
constructed fixtures, not by reading IMPLEMENTATION.md. Six forgeable-
certificate holes (self-signed waiver, ghost signers, G4 refutation,
G3 token-digest comparison, recorded-gate refutation, ban-list derivation)
are now enforced with real code. The empty-path crash (H-1) is gone. The
fast-path contradiction is resolved in both SKILL.md and PDS §4 rule 8.
The L3 alias-direction conjunct now fires on upward and primitive-to-
primitive aliases. The banned-slop signature — my own round-01 MAJOR, the
most serious doctrine gap — is now defined in RULES.md and implemented in
the checker with a derivation, a sharing test, and a self-report
contradiction check. The colorSpace gap (round-01 MINOR) is closed: a bare
hex $value now fails L3. The WCAG SC 2.2.2 source citation (round-01
MINOR) is fixed. SC 1.4.4 and 1.4.12 are now marked informative (round-01
MINOR). The U1 assignment is now checker-verified with a real rotation
function (kimi-1 round-01 MAJOR). The G1 conjuncts (kimi-1 round-01
MAJOR) are restored.

Two issues remain, both already raised by codex-1 in this round and
confirmed by my own probes: same-tier alias edges (semantic→semantic) are
accepted, and the frontmatter parser does not reject trailing `#`
comments. I independently reproduced both. I also found that the shipped
sound-run fixture uses a shared `tokens.json` for both DIRECTIONs despite
PDS §1 mapping DIVERGE to per-agent `<agent>.tokens.json` sidecars — the
checker does not enforce per-agent token identity. The gate error-string
separator drift (em-dash in PDS vs colon in checker) from my round-01
MINOR is still present. D-1 (byte rebalance) remains acceptable.

I accept with reservations: the implementation is materially sound and
the round-01/02 critical findings are genuinely closed, but the same-tier
alias gap and the frontmatter `#` gap should be fixed before a v1.0.0
tag, and the token-sidecar identity question needs an explicit answer.

## What I verified (commands run, and their result)

1. `npm test` — 212/212 pass, 0 fail, exit 0 (1447 ms).
2. `git rev-parse --short HEAD` — `17f6619`, branch `parley-design-skills`.
3. `node addons/parley-design-check/bin/check.js --help` — exit 0, prints
   usage.
4. `node .../check.js check --json` on the sound-run fixture — exit 0,
   verdict PASS, verified L3.
5. `node .../check.js check --json NOTICE.md` — exit 4, verdict
   UNJUDGEABLE (exit 0-for-clean is no longer the case; round-01 codex-1
   MAJOR "unjudgeable exits clean" is closed).
6. `node .../check.js check --json --registry /nonexistent ...` — exit 3,
   registry status absent (C4 refusal holds).
7. Byte budget: `os.path.getsize` on all four doctrine files —
   SKILL.md 6519, PDS.md 25598, RULES.md 23225, WEB-ANNEX.md 10022,
   total 65364 ≤ 65536. All within test thresholds (7168/25600/24576/
   11264). PDS.md is 5118 bytes over FINAL.md's 20480 per-file limit,
   but within the D-1 rebalanced 25600 threshold.
8. Registry digest: `sha256(RULES.md)[:12]` = `b49ff596451f`, matches
   PDS.md frontmatter `registry-digest: b49ff596451f`.
9. No bundled registry: `find addons/parley-design-check -name RULES.md`
   — 0 results.
10. No placeholder text: regex scan for TODO/TBD/FIXME/XXX/PLACEHOLDER/
    lorem ipsum across 48 .md files — 0 matches.
11. WCAG contrast math: `contrastRatio(toSrgb(...), toSrgb(...))` —
    black/white 21.00, red/white 4.00, #777/white 4.48, #767676/white
    4.54, #1b1b1b/white 17.22. All match published WCAG reference values.
12. No non-builtin `require()`: scanned 26 .js files under
    `addons/parley-design-check` — all requires are `node:` builtins or
    relative paths.
13. Rule registry: 30 rules (19 core + 11 web), 0 duplicate ids, 18
    detectors over 18 rule ids, generated from `lib/detectors/`.
14. Ban list derivation: `slop` class at `T0` tier = 4 rules
    (core:category-guessable, core:decoration-unmotivated,
    core:structural-sameness, core:signature-absent-or-mood). Matches
    `banList()` in engine.js.
15. H-1 (empty waivers crash): constructed a fixture with an empty
    WAIVERS.md entries list — exit 4, verdict UNJUDGEABLE, no crash.
    CLOSED.
16. AF-2 (self-signed waiver): constructed a fixture with
    `granted-by: claude-1, counter-signed-by: claude-1` — exit 1,
    verdict VIOLATION, waiver error: "claude-1 counter-signed its own
    waiver, which is not a counter-signature". CLOSED.
17. Ghost signers (no roster): constructed a fixture with no design
    artifacts (hence no participant roster) and a waiver signed by
    `nobody-1` and `ghost-2` — exit 1, waiver error: "the run records no
    participants, so the independence of nobody-1 and ghost-2 cannot be
    established". CLOSED.
18. G4 refutation: constructed a fixture with an AUDIT artifact, G4
    recorded as `pass`, and an unwaived `core:fabricated-evidence`
    quality VIOLATION in a DIRECTION — exit 1, gate-recorded failure:
    "G4 is recorded pass and this run's own findings refute its
    conditions: core:fabricated-evidence". CLOSED.
19. Alias direction — upward (primitive→semantic): constructed
    `primitive.slate` aliasing `semantic.text` — exit 1,
    `l3-alias-direction` finding: "primitive.slate is a primitive and
    aliases semantic.text, which is a semantic". CLOSED.
20. Alias direction — primitive-to-primitive: constructed
    `primitive.ink` aliasing `primitive.slate` — exit 1,
    `l3-alias-direction` finding: "primitive.ink is a primitive and
    aliases the primitive primitive.slate". CLOSED.
21. Alias direction — same-tier (semantic→semantic): constructed
    `semantic.muted` aliasing `semantic.text` — exit 0, verdict PASS,
    verified L3, alias-direction obligation `met`. NOT CLOSED — codex-1
    round-03 MAJOR confirmed.
22. Bare-hex colorSpace: constructed `"$value": "#1b1b1b"` without
    `colorSpace` — exit 1, `l3-colour-space` failed. My round-01 MINOR
    CLOSED.
23. Fast-path contradiction: SKILL.md says "The fast path is not a
    Parley Deck run" and "claims no Parley verification and no level
    above L1 (PDS §4 rule 8)". PDS §4 rule 8 says "The fast path runs
    outside Parley... MUST NOT claim Parley verification, nor a level
    above L1." Both files now agree. CLOSED.
24. Banned-slop signature: PDS §3 G1 now says "defined, with the ban
    list and the sharing test, in RULES.md class slop". RULES.md defines
    it: ban list = slop rules at T0; a signature is the set of ban-list
    ids that fire against a direction; sharing = intersection of 2+ ids
    or 1 id with same evidence. Checker `checkG1()` (L850–1049)
    implements: `banList()` derivation, `sharedSignature()` test,
    `signatureLedger()` reading, and a self-report contradiction check
    against the run's own findings. My round-01 MAJOR CLOSED.
25. U1 assignment: `rotateAssignment()` (L406–411) computes
    `sha256("PDS/1" + runId)` and rotates sorted positions. L815–847
    verifies each DIRECTION's `assigned` field against the recomputed
    rotation and checks decline recording. `run-id` is in
    `REQUIRED_FIELDS["DESIGN-BRIEF"]`. kimi-1 round-01 MAJOR CLOSED.
26. G1 conjuncts: PDS §3 G1 now reads "Persistent convergence never
    auto-passes: it proceeds only past the ban list, past
    core:category-guessable's category-plus-avoidance test, and on
    recorded human ratification with a brief-specific reason." All
    three conjuncts present. kimi-1 round-01 MAJOR CLOSED.
27. WCAG SC 2.2.2: `core:motion-without-reduced-path` now lists
    `sources: [WCAG-2.2-SC-2.2.2, WCAG-2.2-SC-2.3.3]`. My round-01
    MINOR CLOSED.
28. SC 1.4.4/1.4.12 informative: WEB-ANNEX now marks them "Recorded,
    not enforced (informative)". My round-01 MINOR CLOSED.
29. Frontmatter trailing `#` comment: constructed a fixture with
    `x-test: value # trailing comment` in DESIGN-BRIEF frontmatter —
    exit 0, verdict PASS, frontmatter-parses obligation `met`. The
    parser treats `#` as part of the scalar value, not as a comment
    start. codex-1 round-03 MAJOR confirmed.
30. Token-sidecar identity: sound-run fixture has both DIRECTIONs
    pointing to `../tokens.json` (shared), while PDS §1 maps DIVERGE to
    `round-01/<agent>.md + <agent>.tokens.json` (per-agent). Checker
    passes L3. See Finding MAJOR-3.
31. Gate error-string separator: PDS §3 uses em-dash (`G1 DISTINCTNESS
    —`), checker uses colon (`G1 DISTINCTNESS:`). Still divergent. My
    round-01 MINOR still open.
32. `tokens: ""` on DIRECTION (H-1 variant): constructed fixture with
    empty tokens path on a DIRECTION — exit 4, verdict UNJUDGEABLE,
    G2 unverified, no crash. CLOSED.
33. kimi-1 round-01 MINOR (wrong rule ids in PDS examples): checked for
    `core:states-incomplete`, `core:raw-value-outside-tokens`,
    `core:effect-budget` (without -exceeded) — none found. All
    examples now use correct ids. CLOSED.

## Findings

### [MAJOR] Same-tier alias edges are accepted (codex-1 round-03, confirmed)

PDS §3 G3 says "a reference points down that order" where the order is
primitive < semantic < component. The checker's `tokenProblems()` (engine.js
L1326–1340) tests two conditions: `to > from` (upward alias, blocked) and
`from === 0 && to === 0` (primitive-to-primitive, blocked). A same-tier
alias where `to === from` and `from !== 0` (e.g. semantic→semantic,
component→component) falls through both guards and is accepted.

I constructed a fixture with `semantic.muted` aliasing `semantic.text`:
exit 0, verified L3, alias-direction obligation `met`.

The spec says "points down." "Down" means strictly downward in an
ordering. If same-tier were permitted, the spec would say "points down or
within" or "does not point up." The word "down" alone denotes strict
descent. A semantic token aliasing another semantic token is not pointing
down — it is pointing across, and the checker certifies it as L3-clean.

This is the same gap codex-1 identified. I confirm it independently.

Fix: change the guard at L1326 from `to > from` to `to >= from`, or add
an explicit same-tier case: `else if (from !== null && to !== null && from
=== to && from !== 0)`. Ship a fixture with a semantic→semantic alias and
assert it fails. The PDS text needs no change — "points down" already
means strictly down.

### [MAJOR] Frontmatter parser accepts trailing `#` comments (codex-1 round-03, confirmed)

PDS §2 rule 5 (the canonical frontmatter subset) says: "comment a whole
line whose first non-space character is `#`; never trailing." The parser
in `registry.js` `parseScalar()` (L76–88) treats `#` as an ordinary
character in an unquoted scalar. An unquoted value like `x-test: value #
trailing comment` is parsed as the string `value # trailing comment`,
not rejected.

I constructed a fixture with this exact input: exit 0, verdict PASS,
frontmatter-parses obligation `met`.

A parser that accepts more than the spec publishes is how the published
examples stop being the contract (registry.js L15–17 says this itself).
The checker's own comment identifies this as a failure mode to prevent,
and then prevents the line-initial case (L157: `/^\s*#/.test(line)`) but
not the trailing case.

Fix: in `parseScalar()`, after stripping quotes, check for an unquoted
`#` in the raw text and raise a `CheckError`. Alternatively, truncate the
value at the first unquoted `#` (matching YAML semantics) and warn. The
former is stricter and matches the spec's "never trailing" language.

### [MAJOR] Sound-run fixture uses shared token sidecar; checker does not enforce per-agent identity

PDS §1 maps DIVERGE to `round-01/<agent>.md + <agent>.tokens.json` — one
token sidecar per agent. PDS §2 DIRECTION says the `tokens` field is
"This direction's token file." The PDS §2 DIRECTION example shows
`tokens: round-01/codex-1.tokens.json`.

The shipped sound-run fixture has both DIRECTIONs pointing to
`../tokens.json` — a single shared file at the run root. The checker
passes L3 on this fixture. Neither the checker nor any test enforces that
each DIRECTION carries its own token sidecar, or that the winner's
DIRECTION tokens file is distinct from the loser's.

This matters for G2's token-immutability conjunct: if both directions
share one file, a graft that "adds a token to the winner's file" has
modified the file the loser also declared, and the digest comparison
cannot distinguish the two. The immutability protection G2's digest check
provides is weakened when the inputs are not separated.

The counterargument is that PDS §1's `+ <agent>.tokens.json` is a
mapping convention, not a normative requirement — the DIRECTION's
`tokens` field says "this direction's token file," and if two directions
happen to point at the same path, that is the author's choice. But the
fixture is the shipped example that teaches users what a conformant run
looks like, and it teaches a shape PDS §1 does not show.

Fix: either (a) split the sound-run fixture so each DIRECTION has its own
`<agent>.tokens.json` and the CONTRACT points at the winner's, matching
PDS §1 and §2; or (b) add a note to PDS §2 clarifying that the `tokens`
field may point at a shared file when the directions share a token layer,
and explain why G2's digest check is still meaningful in that case. Option
(a) is simpler and matches the spec's own example.

### [MINOR] Gate error-string separator drift (round-01, still open)

PDS §3's gate error strings use an em-dash separator:
`G1 DISTINCTNESS — directions '<a>' and '<b>' differ on...`
The checker's `ledger.fail()` calls use a colon:
`G1 DISTINCTNESS: directions '${a}' and '${b}' differ on...`

This is the same finding from my round-01 MINOR. PDS §3 says these are
"the checker's" error strings, implying they should match what the
checker emits. A reviewer comparing spec text against a checker report
will see different separators and not know which is authoritative. The
drift is across all four gates (G1–G4).

Fix: pick one separator and align both. Either change PDS §3's error
string blocks to use the colon, or change the checker's `ledger.fail()`
calls to use the em-dash. The finding format (`rule-id — violation —
remedy`) already uses em-dash, so changing the checker to match PDS would
also make the checker internally consistent.

### [MINOR] D-1: PDS.md exceeds FINAL.md's per-file limit by 5118 bytes

FINAL.md sets PDS.md's per-file limit at 20 KiB (20480 bytes). The
shipped PDS.md is 25598 bytes — 5118 bytes over. The test enforces the
D-1 rebalanced threshold of 25600, so the green test proves only the
rebalanced criterion, not the binding one.

I accept D-1 in principle: C3 makes the 64 KiB total binding, and the
per-file split was one participant's proposal that was never individually
ratified. The total is held (65364 ≤ 65536). But PDS.md now sits 2 bytes
from its rebalanced ceiling (25598/25600), which is even tighter than
RULES.md's 351-byte headroom (23225/24576) that kimi-1 flagged in
round-01. The next PDS edit will require a rebalance.

This is not a blocker — the total budget is the binding constraint and it
holds. But the shipped state leaves no room for PDS.md to grow, and the
rebalance was supposed to give the registry room, not PDS.md.

### [NIT] DTCG "2025.10" date-pin notation (round-01, now documented)

NOTICE.md now states: "DTCG 2025.10, as the doctrine pins it, names that
editor's draft as of October 2025; the DTCG has published no formal
versioned release, so the string is a date pin, not a W3C Recommendation."
My round-01 NIT is closed.

## Cross-reviewer round-01 findings: closure status

### codex-1 round-01

- **[CRITICAL] L2 and L3 can be falsely certified**: CLOSED. The checker
  now models each level as an explicit obligation set. L2 validates the
  U1 assignment (`rotateAssignment` + position comparison), G1/G2
  recomputation, G3/G4 gate-recorded refutation, and G2 token-digest
  comparison. L3 requires a DTCG token document, alias resolution, alias
  direction, colorSpace, and system rules. A wrong assignment, missing
  G3/G4, or open quality violation now fails the level. Verified by
  probes 18–20, 25.

- **[CRITICAL] Self-signed waiver**: CLOSED. `waiverProblem()` rejects
  equal grantor and counter-signer ids and validates both against the
  run's participant roster. Verified by probe 16 (self-signed rejected)
  and probe 17 (ghost signers rejected without roster).

- **[MAJOR] Unjudgeable run exits clean**: CLOSED. Exit 0 is now reserved
  for PASS. UNJUDGEABLE exits 4, VIOLATION exits 1. Verified by probe 5
  (NOTICE.md → exit 4) and probe 4 (sound-run → exit 0).

- **[MAJOR] Artifact ingestion rejects valid PDS YAML**: PARTIALLY
  CLOSED. The parser was extended to handle flow maps and block lists.
  However, the trailing-`#` gap (codex-1 round-03 MAJOR, my MAJOR-2
  above) shows the parser still accepts constructs the spec forbids.
  The unknown-rule-id traversal is now handled: the checker reports
  UNJUDGEABLE for registry rules it cannot evaluate, and rule-id-bearing
  fields in CRITIQUE/VERDICT/AUDIT/WAIVERS are traversed.

- **[MAJOR] Reduced-motion block passes without reducing motion**: not
  re-verified in this round (outside my doctrine lens). The fixture pair
  for `motion-reduced-path.js` passes in the test suite. I note the
  detector exists and has a fail/pass pair.

- **[MAJOR] Focus-visible replacement reported as absent**: not
  re-verified in this round (outside my doctrine lens). The fixture pair
  for `focus-indication.js` passes in the test suite.

- **[MAJOR] D-1 weakens binding acceptance criterion**: I maintain my
  accept position (see MINOR above). C3's total is binding; the per-file
  split was never individually ratified. codex-1 and I disagree here,
  which is a legitimate consensus split.

- **[MINOR] Heroicons counted as two icon sources**: not re-verified
  (outside my lens). Test suite passes.

- **[MINOR] PDS normative requirements not all numbered**: not
  re-verified. I note the artifact tables still use RFC 2119 keywords
  in field descriptions without numbered rule references, but this is a
  citation-grammar issue, not a conformance break.

### kimi-1 round-01

- **[MAJOR] G1 drops ban list + category test conjuncts**: CLOSED. PDS
  §3 G1 now reads: "Persistent convergence never auto-passes: it
  proceeds only past the ban list, past core:category-guessable's
  category-plus-avoidance test, and on recorded human ratification with
  a brief-specific reason." Verified by probe 26.

- **[MAJOR] U1 assignment unverifiable**: CLOSED. `run-id` is in
  DESIGN-BRIEF's required fields. `rotateAssignment()` computes the
  rotation from `sha256("PDS/1" + runId)`. The checker verifies each
  DIRECTION's `assigned` field against the recomputed rotation. Verified
  by probe 25.

- **[MINOR] PDS examples cite nonexistent rule ids**: CLOSED. Checked
  for `core:states-incomplete`, `core:raw-value-outside-tokens`, and
  `core:effect-budget` (without -exceeded) — none found. All examples
  now use correct ids. Verified by probe 33.

- **[MINOR] §3 strings don't match checker / §9 L2 row overstates**: The
  separator drift is still present (my MINOR above). The §9 L2 row now
  says "a recorded gate for every §3 transition the run crossed" (was
  "every gate of §3 recorded"), which correctly scopes to crossed
  transitions. The string drift remains.

- **[MINOR] RULES.md 30 bytes of headroom**: RULES.md is now 23225/24576
  — 351 bytes of headroom (was 30 bytes at 24546/24576). Improved but
  still tight. PDS.md is now the tightest file at 25598/25600 (2 bytes).

- **[NIT] Tier values in YAML use short form while Keys table says full
  form**: not re-verified. The parser accepts both, so this is
  doc-vs-data only.

## Open questions

1. The same-tier alias gap (MAJOR-1) turns on whether "points down" in
   PDS §3 G3 means strictly down or "not up." The spec's language
   supports strict descent, but the checker implements "not up." Which
   reading is intended? If same-tier is permitted, PDS §3 should say so
   explicitly; if not, the checker guard needs `>=` not `>`.

2. The sound-run fixture uses a shared token sidecar (MAJOR-3). Is this
   a deliberate simplification (the spec's `<agent>.tokens.json` is a
   convention, not a requirement), or is it a fixture that teaches the
   wrong shape? If the former, PDS §2 should clarify that the `tokens`
   field may point at a shared file. If the latter, the fixture should
   be split.

3. The frontmatter `#` gap (MAJOR-2) is in the same parser that reads
   the registry. A trailing `#` in a `pds-rule` block's value would
   silently include the comment in the scalar. Is there a test that
   exercises this against the registry parser specifically, not just
   the artifact frontmatter parser?
