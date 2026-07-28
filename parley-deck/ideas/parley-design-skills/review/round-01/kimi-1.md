---
agent: kimi-1
idea: parley-design-skills
review-round: 01
date: 2026-07-28
reviewed-commit: 726c024
---

## Summary

The implementation is substantively conformant. Under my lens: the C1 mapping table
appears in PDS.md §1 row-for-row; G1's unconditional two-axis test (A1) is stated with
its duplicate-Signature and banned-slop-signature conditions; A3's unattended-stop is
carried almost verbatim in §4 rule 7; the eight artifact kinds each keep the identical
four-part shape (purpose line, rationale paragraph, required-fields table, minimal
example — I checked all eight); tiers are ordinal and spelled number-plus-word with
engine names confined to WEB-ANNEX.md; rule classes, waivers, L1–L4, extension policy,
append-only versioning and the deprecation window all match what consensus ratified;
§5 rule 7 adopts the Phase-6 reviewer as DESIGN-SYSTEM author, closing the carried-open
question as FINAL.md required. The registry-digest in PDS.md's frontmatter genuinely
matches sha256(RULES.md), and a test guards it.

Two MAJOR findings, both in the protocol layer rather than the checker: the G1 remedy
silently drops two of C7's ratified conjuncts, and the U1 deterministic assignment
ships as an unverifiable MUST — no artifact carries `run_id`, and no code or stated
procedure checks the rotation. Neither defeats the design; both are exactly the class
of defect A1 was written about (a protection that reads as present and cannot fire).

Position on the declared deviations. **D-1 (byte rebalance): accept.** C3's adopted
text makes "64 KiB total" the decision; three different itemised per-file splits were
proposed in round-02 (hermes-1's 60 KB, codex-1's 64 KiB, my own 64 KB) and none was
individually ratified, so the per-file numbers in FINAL.md are illustrative and the
total is binding. Holding 64 723 ≤ 65 536 bytes while preserving the mandated
four-part artifact shape is the right trade against the alternatives D-1 names. Two
caveats: the rebalanced budgets sum to 65 KiB (8+22+24+11), which is only safe because
the total is enforced by its own test; and the rebalance left RULES.md with 30 bytes
of headroom (MINOR, below). **D-2 (check-rules without detectors): accept**, with a
counting note. The behaviour is as FINAL.md requires — verified in the report output,
which names every undetectable rule UNJUDGEABLE with its reason. But "five rules"
counts only those whose stated reason is "no detector"; `web:contrast-ratio`,
`web:target-size` and `web:reflow-narrow` are also `enforced-by: check`/`both` with no
detector, masked by their T2 tier. Eight rules, not five; the deviation note should
say so.

## What I verified (commands run, and their result)

- `npm test` from the skill repo root: **158/158 pass** (registry grammar, detectors'
  fixture pairs, refusal path, exit codes, waivers, conformance, byte budgets,
  no-placeholder scan, no-bundled-registry, digest guard, installer discovery/`--only`/
  `--no-addons`/markers). IMPLEMENTATION.md's central claim holds.
- `shasum -a 256 addons/parley-design/references/RULES.md` → `1fbe071e1222…`, matching
  PDS.md's declared `registry-digest`. The C4 signature mechanism is real.
- Rule-id cross-reference: extracted every `core:`/`web:` id referenced anywhere in the
  four doctrine files and diffed against every `id:` the registry declares. Three
  referenced ids do not exist (MINOR, below).
- `node addons/parley-design-check/bin/check.js check --level L2` on
  `test/fixtures/conformance/sound-run` → "claimed L2, verified L2", exit 0, with 16
  rules listed UNJUDGEABLE (not silently skipped) and the 11 web rules out of scope.
- Same command on `collapsed-run` → G1 recomputed from declared positions ("differ on
  1 declared axis; 2 are required"), missing G1/G2 records flagged, malformed outcome
  and 4-graft verdict flagged, "not verified", exit 1. The two-axis gate fires.
- Refusal: copied the checker alone to `/tmp/pds-probe` and ran it on a plain markdown
  file → stderr refusal naming the absent registry, exit **3**; `--registry` pointing at
  a nonexistent path → exit 3 as well (a named registry is never silently replaced).
- `check` on `literal-outside-tokens/fail` → three findings in exact
  `rule-id — violation — remedy` form with stable paths and lines.
- Waivers, all four fixtures by hand: valid suppresses and is recorded under WAIVED
  with expiry and counter-signatory; expired is rejected "treated as absent";
  system-blind widening is rejected ("scopes to the ratified system itself"); wildcard
  (`core:*`) is rejected. Each rejection keeps the finding. Matches C13/§8 exactly.
- `grep -rn "sha256\|run_id\|rotate" addons/parley-design-check/{lib,bin}` → the only
  hit is the registry digest hash. No assignment verification exists (MAJOR, below).
- Byte budgets: SKILL 6 656 / PDS 22 389 / RULES 24 546 / WEB-ANNEX 11 132 = 64 723 ≤
  65 536; each within its rebalanced budget; RULES.md has **30 bytes** of headroom.
- Capability generation: the report's "18 detectors over 18 rule ids, generated from
  lib/detectors" matches the directory contents; `engine.js` builds it from
  `readdirSync`, never from a hand-written list.
- NOTICE.md: credits hallmark (MIT) and impeccable (Apache-2.0) as prior art studied,
  states nothing was copied. I did not diff rule text against those projects; the
  registry prose reads as independent authorship on its face.
- Not re-verified: every detector's internal logic (I read engine/registry/artifacts/
  CLI in full and skimmed detector output via fixtures; the fixture-pair suite passed).

## Findings

### [CRITICAL] None raised.

### [MAJOR] G1's persistent-convergence remedy drops two of C7's ratified conjuncts

**What.** C7 (as amended by A1) is explicit: persistent convergence "never auto-passes:
it proceeds only past the ban list and the category-plus-avoidance test **and**
on-record human ratification with a brief-specific reason, or it returns `ABSTAIN`."
PDS.md §3 rule 1 ships: "persistent convergence proceeds only on recorded human
ratification with a brief-specific reason, or returns `ABSTAIN`." The ban-list test
and the category-plus-avoidance test are gone, and the G1 error string in §3 omits
them too.
**Why it matters.** Consensus put two agent-computed tests in front of the human
ratification escape so that a collapsed set cannot be waved through on goodwill alone.
A facilitator working from PDS.md alone — the fixed read set for this transition is
PDS §3 §4 — will ratify a converged set without either test, and the run will read as
protocol-conformant. This is a normative divergence from a ratified decision in the
document my lens exists to check.
**Fix.** Restore both conjuncts in §3 rule 1 and in the G1 message: persistent
convergence proceeds only past the ban list and the category-plus-avoidance test and
recorded human ratification with a brief-specific reason, else `ABSTAIN`. Say where
the ban list lives (the registry's `slop` rules are the natural home — see Open
questions).

### [MAJOR] The U1 assignment is an unverifiable MUST: no `run_id` home, no verifier

**What.** §4 rule 2 mandates the deterministic assignment
`rotate(sorted(primary_positions), uint32(sha256("PDS/1" || run_id)[0:8]))` and the
brief's minimum-position count. Two holes. (a) The formula needs `run_id`, but §2's
DESIGN-BRIEF required-fields table defines no `run-id` field, so no artifact carries
the input and the rotation is not reproducible even by hand. (b) Nothing verifies it:
the checker's L2 recomputes G1 from declared positions but contains no rotation code
(verified by grep — the only sha256 in the checker is the registry digest), and
`assigned`/`declined` fields and the position-count requirement are never checked.
U1 was adopted as "cheap, offline, reproducible and **checker-verifiable**", and my
own signoff recorded the audit goal as "the checker-verified mapping plus the recorded
decline reasons". As shipped, a run can ignore the assignment entirely and still
verify L2.
**Why it matters.** A1 exists because a gate that cannot fail "is worse than no gate,
because the run reads as protected". A MUST that no party can check is the same shape
of failure, one section over.
**Fix.** Add `run-id` to DESIGN-BRIEF's required fields in §2 (and the checker's
`REQUIRED_FIELDS`), and add an L2 conformance check that recomputes the rotation from
the brief's primary positions and `run-id`, compares each DIRECTION's `assigned`,
records declines, and fails the brief when distinct primary positions < proposers. If
that is judged out of v1 scope, say so explicitly in PDS §12 and IMPLEMENTATION.md —
do not leave the MUST reading as enforced.

### [MINOR] PDS.md's own examples cite three rule ids the registry does not define

**What.** §2 examples use `core:states-incomplete` (CRITIQUE and VERDICT examples),
`core:raw-value-outside-tokens` (AUDIT example) and `core:effect-budget` (WAIVERS
example). The registry ids are `core:interaction-states-incomplete`,
`core:literal-outside-token-layer` and `core:effect-budget-exceeded`. Verified by
diffing every referenced id against every declared id.
**Why it matters.** Under §10 rule 3 a conformant consumer reads each of these as an
unknown id and reports UNJUDGEABLE — the spec's minimal examples teach citations that
its own extension policy launders into non-findings, in the document whose entire
point is citable rule ids.
**Fix.** Rename the four occurrences to the real ids. One line each; no normative
content changes.

### [MINOR] §3's "The strings below are the checker's" is false, and §9's L2 row overstates what the checker records

**What.** Three instances of gate text/code drift. (a) §3 shows
`G1 DISTINCTNESS — directions '<a>' and '<b>' …` with the remedy "Re-diverge once with
the seeded assignment (§4 rule 2), or record human ratification." The checker emits
`pds-check:l2-gate-g1 — G1 DISTINCTNESS: directions 'ledger' and 'atrium' …` with
different wording (verified on the collapsed-run fixture). (b) The G2 string §3 shows
("graft '<n>' from '<handle>' modifies the winner's token file") covers a condition
the checker never detects; it implements outcome-shape, graft-count and graft-fields
only. (c) §9's L2 row requires "every gate of §3 recorded", but the checker requires
only G1 and G2 records — the sound-run fixture verifies L2 with a CONTRACT present
and no G3 record anywhere.
**Why it matters.** The doctrine's founding defect list bans two hand-maintained
representations precisely because they drift; here the spec text and the tool have
already drifted at v1.0.0, and a reviewer citing §3's strings against real checker
output will conclude the run is non-conformant when it is the text that is wrong.
**Fix.** Pick one direction per instance: reword §3's preamble to "canonical message
shapes" (or make the checker emit the §3 strings verbatim where it detects the
condition), and scope §9's L2 row to "every gate whose transition the run has crossed"
— or require a G3 record once a CONTRACT exists. The checker's own SKILL.md already
says "recomputed where the artifacts allow it"; PDS.md should stop claiming more.

### [MINOR] RULES.md has 30 bytes of headroom, against the reason 64 KiB was chosen

**What.** RULES.md is 24 546 of its 24 576-byte budget. C3's corrected rationale says
64 KiB was adopted over 60 KB specifically because the registry "is the file most
likely to need" growth room — and rule ids are append-only, so the file only grows.
The first appended rule forces an immediate rebalance or a prose cut.
**Why it matters.** Not a conformance break — the budget test passes — but the shipped
state negates the stated rationale for the adopted ceiling on day one, and the next
registry edit will arrive bundled with an unrelated budget negotiation.
**Fix.** Buy headroom now within the held 64 KiB total (shave the rebalanced PDS.md
budget, or tighten registry prose — the counterexamples are the compressible part), or
record in IMPLEMENTATION.md/§12 that the next rule addition ships with a deliberate
rebalance commit.

### [NIT] RULES.md's Keys table documents tier values its own data doesn't use

**What.** The Keys table lists `tier` values as `T0 ARTIFACT`, `T1 SOURCE`, `T2 RENDERED`,
`T3 PIXEL`; every `pds-rule` block in RULES.md and WEB-ANNEX.md uses the bare ordinal
(`tier: T0`), matching FINAL.md's grammar. The parser accepts both spellings, so this
is doc-vs-data only.
**Fix.** Amend the Keys table to say the value is the ordinal and the number-plus-word
form is for prose. One line.

## Open questions

1. **Where does G1's "banned-slop signature" ban list live?** §3 G1 and C7 both invoke
   "the ban list", but no artifact defines it or points at one. The registry's `slop`
   rules are the natural home; if so, §3 should say so, which would also make the
   MAJOR-1 fix and any future checker support mechanical.
2. **Run-level verdict semantics.** The sound-run reports `verdict PASS` with 16 of 25
   in-scope rules UNJUDGEABLE. Per-rule, §6 rule 4 is honoured — they are listed under
   UNJUDGEABLE, never counted as passed — but a roll-up PASS on a mostly-unjudged run
   can read as an all-clear. Should the roll-up distinguish "PASS, everything judged"
   from "PASS with unjudged rules" (or report the judged/total count on the verdict
   line)?
