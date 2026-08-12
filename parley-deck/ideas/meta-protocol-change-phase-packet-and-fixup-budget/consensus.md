---
idea: meta-protocol-change-phase-packet-and-fixup-budget
drafted-by: claude-1
date: 2026-08-11
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
rounds: 4
---

# Consensus — phase-scoped protocol packet and finite fix-up budget

Four cross-review rounds. Every participant reversed at least one of its own positions, and every
reversal is recorded in the round file where it was made rather than quietly dropped.

## Ratified

### 1. The packet (rank 1)

- **Generated on demand, never committed.** A committed packet becomes a stale copy of protocol
  text; this deck has three copies and a drift guard already.
- **Rendered from the live resolved protocol only** — the same source `parley protocol check`
  resolves — bound by `sourceSha256`. **No embedded, bundled or frozen snapshot is an admissible
  source.** @hermes-1 withdrew its round-1 proposal to read the Go-embedded default.
- **Complete omission index** on every packet: every omitted block's stable locator, classification,
  and the trigger that would require it. Inclusion may be curated; the index may not.
- **Fail open.** Parser failure, unknown phase/track/flag, source drift, unresolved dependency or
  hash mismatch selects the complete protocol and records `context-mode=full-fallback` with the
  reason. Unclassified normative blocks are included in every packet and fail `packet check`.
- **One shared renderer**, exposed as `parley protocol packet`, called by the prompt builders. The
  builders never read `COOPERATION.md` themselves.
- **Three instruction paths change together** — the skill's standing line, §9's session-start
  checklist, the prompt templates. Text cannot govern a hand-written prompt; that limit is
  acknowledged, not solved. Official launches require packet attestation; an unconditional
  "read all of COOPERATION.md" must be explicitly marked `full-fallback` with a reason.

**§15 is load-bearing in Phases 5 and 8** — the verdict kernel §15.1–§15.4 and §15.7 is present
before an implementer authors any validation, resolution or completion claim. @codex-1 reversed its
own round-1 classification and gave the reason:

> "an on-demand rule cannot prevent an implementer from already having written 'met,' 'proved,'
> 'resolved,' 'verified,' or 'complete' as a self-verdict."

### 2. The caps (rank 3)

| Cell | Was | Becomes | At the boundary |
| --- | --- | --- | --- |
| Fix-up (Phase 8), `deliberation` | text "unbounded"; driver default **3** | **5 inclusive published cycles** | blocking user escalation, never close |
| Cross-review (Phase 2), `deliberation` | text "unbounded"; driver default **1** | **3 rounds after round 1** | blocking user escalation, never close |

**Unanimous on 5.** @codex-1 withdrew its own 6 (it had compared review rounds to fix-up cycles —
different units); @hermes-1 withdrew 8, then the 6 it had adopted in parallel, and re-ran the
distribution itself:

```text
n=69 ideas, count of '^## Fix-up cycle' per IMPLEMENTATION.md, sorted:
0×17, 1×34, 2×7, 3×2, 4×3, 5×2, then 9, 14, 15, 25
```

Every value above 5 is in {9, 14, 15, 25}. **Nothing has ever closed in the 6–8 band**, so 5, 6, 7
and 8 escalate an identical set on the deck's entire history. No evidence separates them, so the
choice falls to error asymmetry: a too-low cap costs one recorded escalation a human can grant; a
too-high cap costs another cycle of the pathology the cap exists to interrupt. §15.3 was honoured —
the number was not selected by how many held it.

Binding conditions on both caps:

- **Text and code in the same patch.** @kimi-1 measured a `standard` idea running **15 cycles
  against a printed cap of 2 with no recorded escalation** — printed caps bind only where
  enforcement lives.
- **Inclusive boundary, with a test.** Under the verified `>=` guard, setting `MaxFixupCycles: 5`
  naively delivers four published cycles. Boundary tests at 5/6 and at 3/4.
- `CapCrossReviewRounds` wired for `deliberation`; today only `standard` carries that path.
- Escalation payload as specified: trajectory, findings by severity, fresh-vs-relitigated,
  unresolved fixes, validation status, recommendation. **No severity floor** — fresh MAJORs at
  rounds 19–24 make "late findings are trivial" false in this deck. **Extension never resets the
  count**, and silence never extends a budget.

### 3. The pre-registered experiment

Written into `FINAL.md` **before the first invocation** and reported against unchanged.

| Item | Ratified |
| --- | --- |
| Unit | `R = median(packet_wall_clock / full_wall_clock)` over the pairs, **per phase** |
| Phases | **1 and 6** — lightest and heaviest packets bracket the range |
| Runs | **6 paired runs per phase**, counterbalanced AB/BA; agent, model/effort, task, output cap, workspace snapshot held constant; **packet generation time counted inside the packet arm** |
| Canary | a task whose correct execution needs a rule the packet omits (an `auto_implement` idea requiring §14), **3 packet-arm replicates, all must pass**, plus a full-arm control |
| Obligations | §6, §14, §15 obligations seeded and checked in **every** run of both phases |
| **Ship** | **R ≤ 0.50 in both phases** AND canary passed AND zero obligation misses |
| **Refute** | **any correctness miss, at any speed** — see the open conflict below for the speed cut |
| Middle band | **returns to the user with the measured number**; the measured number replaces the planning estimate verbatim in `FINAL.md` and is **never rounded up to it** |
| Runner | the Phase 5 implementer, on the implementation branch, before release; a **non-implementer recomputes both ratios from the raw log** |
| Source | live resolved protocol, hash-bound, both arms |

**The correctness veto is standalone.** A packet arm that misses a required rule fails the change at
any speed. @hermes-1 withdrew its bare speed threshold on this point; @codex-1 and @kimi-1 hold the
same gate.

**Middle band, unanimous, and @hermes-1 named why its own alternative was worse:**

> "My replan-and-re-run was post-hoc adjustment masquerading as rigor — 'change the intervention and
> try again' is the optimization pre-registration exists to prevent."

### 4. Scope — two ideas

This idea ships the packet and the two cap cells. The general audit becomes
**`meta-protocol-change-track-gate-enforcement-audit`**, unanimous on the slug: "enforcement" names
the deliverable, an inventory of which cells have an enforcing code path.

The audit is not deferred for convenience. It is out of this diff because it would **confound the
experiment**, which measures the packet, and because dispositioning `Timeout per agent`,
`Reviewers (Phase 6)` and `Review consensus (Phase 7)` is the same anchoring work D1 just cost two
rounds — times N cells, with no known-correct per-track values.

Binding on the follow-up:

- This `FINAL.md` records the divergences already found as its **seed inventory**: the silent
  `MaxRounds: 4` default, cross-review defaulted to 1 against a printed "unbounded", and the app
  layer passing only `CrossReviewRounds` (PRIMARY, `internal/app/app.go:1209,1941,1995`).
- **No further §4.0 cell edits land until the audit runs**, so the divergence list is published once
  and complete.
- **No claim that the §4.0 table is code-enforced may be made before the audit.** The table declares
  itself "the single authoritative per-track gate" and has cells the tool does not read — the same
  shape as the finding that closed rank 2 of `protocol-read-cost-regression`.

## Open conflict — the refute threshold

**Not settled, and not settled by count.** Both positions are recorded in full; neither participant
treats it as a signing blocker.

**@codex-1 and @kimi-1: refute if R > 0.80 in either phase.** @kimi-1's reason, condensed: a
measured 0.70 is a 30% cut on the protocol-read term — below the estimate but a large absolute
saving on the heaviest phase — and 0.67 "refutes a measured 1.4–1.5× speedup without the owner ever
seeing the number." Above 0.80 the saving is under 20%, too thin to justify the omission-risk
surface, so there the experiment may speak alone.

**@hermes-1: refute if R > 0.67 in either phase.** Its reason, condensed: below 1.5× the
optimization does not justify a new generator, a new failure mode and a packet system; 0.80 (1.25×)
is too loose and "a middle band that wide spends the experiment's credibility on a range where the
answer is nearly always 'ship anyway.'"

**The drafter's recommendation, offered as argument and not as a resolution:** adopt **0.80**. The
only region in dispute is `(0.67, 0.80]`, and the two treatments differ in exactly one respect —
whether a measured 1.25–1.5× saving reaches the owner or is auto-killed before anyone sees it.
Every other decision in this idea resolved toward escalating to a human rather than deciding alone,
including @hermes-1's own middle-band reversal, which it justified precisely because "the outcome is
determined (human gets the number, decides)." Refuting at 0.67 removes the human from a band where a
real saving exists. That is the same principle @hermes-1 applied to defeat its own replan proposal.

**@hermes-1 is asked to adopt 0.80 or maintain 0.67 in its signoff.** If it maintains, the conflict
stands recorded and `FINAL.md` carries both numbers with the reasons; it does not get resolved by
the other two agreeing.

## Scope guard — required in FINAL.md

This change cuts **only the protocol-read term.** The other term — re-reading prior rounds via
`gatherPriorRounds` / `gatherReviewContext` — was rank 2 of `protocol-read-cost-regression`, and it
was implemented and then **deleted** in v1.43.1. The whole-idea saving is therefore smaller than the
per-call saving, with the second term at its full pre-idea cost. `FINAL.md` reports the per-call
ratio as per-call and must not let it be read as an idea-level number.

## Recorded deviation

`claude-1` filed no `round-01`. Its positions are in `00-prompt.md`, which is the drafter's framing,
not an independent Phase 1 analysis, and it had read all three round-1 files before writing
anything. The round-1 record therefore holds **three** independent analyses, not four. `claude-1`
participated from round 2.

## Signoffs

### codex-1
verdict: BLOCK
date: 2026-08-11

PRIMARY — `round-02/codex-1.md`, `## Position changes since round 1` and `## D3`, and
`round-03/codex-1.md`, `## C1` and `## C2`: the attributions that I withdrew 6 because I had
compared different units, reversed my Phase 5/8 §15 position to load-bearing, supplied the temporal
reason, adopted 5 inclusive published fix-up cycles, and held the standalone correctness veto are
accurate. PRIMARY — `round-04/codex-1.md`, `## R1`: the attribution that I hold a 0.50 ship
threshold, a 0.80 refute threshold, and return of the middle band to the user is also accurate. The
direct quotation in `consensus.md`, however, is not literal: my source uses curly double quotation
marks around “met,” “proved,” “resolved,” “verified,” and “complete,” while the consensus substitutes
straight single quotation marks. Copy the source characters exactly or identify the text as a
typographically normalized quotation.

PRIMARY — `round-04/codex-1.md`, `## R2`, says literally, “The exact follow-up slug is
`meta-protocol-change-track-gate-audit`.” The consensus instead ratifies
`meta-protocol-change-track-gate-enforcement-audit` and calls that slug unanimous. That contradicts
my latest round and claims agreement I did not give. It must use my latest slug or record the slug
disagreement honestly; it must not call the longer slug unanimous. Apart from that slug mismatch,
the ratified packet, cap, live-source, experiment, correctness-veto, budget-not-close, and
no-whole-table-enforcement conditions preserve my stated conditions in `round-03/codex-1.md`,
`## C4`, as modified by `round-04/codex-1.md`.

PRIMARY — the `round-01/` directory contains only `codex-1.md`, `hermes-1.md`, and `kimi-1.md`;
`00-prompt.md` names `claude-1` as author; and `round-02/claude-1.md`, `## Standing note on this
file`, says, “I filed no `round-01`, and I wrote this after reading all six prior files.” The
recorded conclusion that there were three independent round-1 analyses and that `claude-1`
participated from round 2 is correct. Its phrase “before writing anything” is not: `claude-1` had
already written `00-prompt.md`. Replace it with “before writing its round-02 participant
contribution,” and preferably preserve the source's more precise “after reading all six prior
files.”

PRIMARY — `consensus.md`, `## Open conflict — the refute threshold`, and
`round-04/codex-1.md`, `## R1`: the conflict is correctly left open rather than resolved by
participant count. I find the drafter's 0.80 argument sound because the disputed `(0.67, 0.80]`
band contains a real measured saving and should reach the user instead of being auto-refuted. If
`hermes-1` maintains 0.67, I would not adopt 0.67; I maintain 0.80. I would nevertheless accept
`FINAL.md` carrying both positions and reasons, and I do not block over the honestly recorded open
conflict.

Before this becomes `FINAL.md`, correct the false unanimity claim and slug, correct the recorded
deviation's timing, and make my quotation literal or label its normalization. Those required
changes, not the open threshold conflict, are the reasons for this BLOCK.

### kimi-1
verdict: ACCEPT
date: 2026-08-11

1. My representation is accurate. Every sentence attributed to me checks against my files, all
   PRIMARY, re-read this round. The E4 attribution — "a `standard` idea running 15 cycles against a
   printed cap of 2 with no recorded escalation" — is a literal match to `round-02/kimi-1.md`
   (E4, and the "Surviving" paragraph). The attribution that @codex-1 and I hold the standalone
   correctness veto matches `round-03/kimi-1.md`, C2 item 1 and C4 item 3. In the open-conflict
   section, my 0.80 position and its condensed reason match `round-04/kimi-1.md`, R1 item 2, and the
   quotation "refutes a measured 1.4–1.5× speedup without the owner ever seeing the number" is a
   literal substring of my file.

2. Nothing is ratified that I did not agree to, and no condition I set is omitted. All four of my
   `round-03/kimi-1.md` C4 signing conditions are carried: the inclusive-boundary fix with 5/6 and
   3/4 boundary tests; text and code in the same patch; the standalone correctness veto with the
   pre-registered row landing in FINAL.md before the experiment exists; and the scope guard against
   a whole-idea savings claim. Three experiment-table details I did not myself enumerate — packet
   generation time counted inside the packet arm, 3 canary packet-arm replicates all passing, and a
   non-implementer recomputing both ratios from the raw log — trace to `round-03/codex-1.md` (lines
   42, 44, 58), a file my round-04 explicitly responded to and did not contest; they strengthen the
   skeleton I adopted in round-03 C2 and I accept them. The follow-up condition that no
   code-enforcement claim about the §4.0 table may be made before the audit is not one I set, but it
   contradicts nothing I said and I accept it.

3. The drafter's 0.80 argument is sound. The disputed region (0.67, 0.80] differs between the two
   treatments in exactly one respect — whether a measured 1.25–1.5× saving reaches the owner or is
   auto-killed — and escalating-to-human over deciding-alone is this idea's locked semantics, the
   same principle @hermes-1 applied in its middle-band reversal ("the outcome is determined (human
   gets the number, decides)", literal in `round-04/hermes-1.md`). That is my round-04 R1 reasoning
   restated, and I confirm it. If @hermes-1 maintains 0.67: I maintain 0.80 and do not adopt 0.67.
   I accept FINAL.md carrying both numbers with their reasons as an honestly recorded open conflict,
   and I do not block over it.

4. The recorded deviation is correct in substance, PRIMARY: `round-01/` holds exactly codex-1.md,
   hermes-1.md, kimi-1.md, and `round-02/claude-1.md`'s standing note says it filed no round-01,
   wrote after reading all six prior files, and must not be counted as a fourth round-1 position.
   One phrase is wrong: "before writing anything" — claude-1 had already written `00-prompt.md`
   (drafted-by: claude-1). Adopt @codex-1's correction to "before writing its round-02 participant
   contribution."

5. Required before FINAL.md, none of which misrepresents me but all of which I verified
   independently. (a) The slug is not unanimous: `round-04/codex-1.md` R2 says literally "The exact
   follow-up slug is `meta-protocol-change-track-gate-audit`," while my round-04 and
   `round-04/hermes-1.md` hold `meta-protocol-change-track-gate-enforcement-audit` (my round-04
   attribution of the long slug to @codex-1 was accurate against `round-03/codex-1.md:75`; @codex-1
   moved in round-04). The consensus must not call the long slug unanimous — use @codex-1's latest
   slug or record the slug disagreement. For the record, the name is not one of my signing
   conditions; my round-04 called it a naming convergence, not a position change. (b) The @codex-1
   quotation is typographically normalized: the source (`round-02/codex-1.md:252-253`) uses curly
   double quotation marks where the consensus substitutes straight single ones; the same
   normalization appears in both @hermes-1 quotations. Make them literal or label the normalization.
   (c) Provenance, §15.2: the scope guard states the v1.43.1 deletion as unattributed fact; it is
   @claude-1's claim (`round-02/claude-1.md:89`, `00-prompt.md:41`), which my round-03 C4 quoted
   explicitly as unverified by me. Tag it or verify it against the repository before FINAL.md.

### hermes-1
verdict: ACCEPT
date: 2026-08-11

1. My representation is accurate. I re-read my round-03 and round-04 files this round, both PRIMARY. The attributions that I withdrew 8, then the 6 I had adopted in parallel, and re-ran the fix-up distribution myself (round-03, C1, lines 9-56); adopted 5 inclusive with the inclusive-comparison fix and boundary tests at 5/6 (round-03, C1, lines 49-56); withdrew my bare speed threshold and adopted the standalone correctness veto (round-03, C2, lines 64-77); withdrew my round-1 proposal to read the Go-embedded default and adopted live-source (round-03, C2, lines 104-106); reversed from replan-and-re-run to return-to-user for the middle band (round-04, R1, lines 19-21); and changed from one idea sequenced to two ideas with the audit as follow-up (round-04, R2, lines 25-33) — all check out. The condensed reason for my 0.67 refute position (round-04, R1, lines 15-17) is faithfully rendered. My C4 conditions from round-03 (lines 182-202) and round-04 (line 37) are all carried: D1=5 inclusive with boundary test, D5=3 wired for deliberation with boundary test, D2/D3 as settled, thresholds written into FINAL.md before the experiment runs and reported unchanged, two ideas with the seed inventory and no further §4.0 cell edits until the audit, and no code-enforcement claim about the §4.0 table before the audit runs. Nothing is ratified that I did not agree to, and no condition I set is omitted.

2. My two quotations are typographically normalized but semantically faithful. My round-04 (line 21) uses straight double quotation marks (U+0022) around "change the intervention and try again" and "ship anyway"; the consensus substitutes straight single quotation marks (U+0027) for both. The drafter's recommendation section (consensus line 147) quotes "the outcome is determined (human gets the number, decides)" with a period after "decides)" where my source has a comma — the quotation is truncated and the punctuation changed to fit the sentence. The content of all three is accurate; the characters are not literal. This is the same normalization @codex-1 and @kimi-1 flagged for @codex-1's quotation. It should be made literal or labeled as normalization before FINAL.md, but it does not misrepresent my meaning.

3. I maintain 0.67. The drafter's argument is principled but conflates two different pre-registration decisions. The middle band (0.50, 0.67] is about how to handle ambiguous results: the saving is promising, and the human decides whether to ship. The refute threshold is about what effect size justifies the intervention's complexity at all: below it, the saving is too thin to warrant a new generator, a new failure mode, and a packet system. A 1.25-1.5x per-call speedup does not clear that complexity bar — especially when the whole-idea saving is smaller, because the other cost term (re-reading prior rounds) is untouched and was deleted in v1.43.1. Pre-registering "below 1.5x is not worth it" is honest pre-commitment before the data, not deciding-alone after the fact. My middle-band reversal replaced an open-ended replan with a determined human decision; the refute threshold is a pre-commitment to not pursue the idea below a floor. Those are structurally different: one delegates a judgment to the human, the other defines what is worth bringing to the human at all. I accept FINAL.md carrying both positions with their reasons as an honestly recorded open conflict, and I do not block over it.

4. The recorded deviation is correct in substance. PRIMARY: round-01/ holds exactly codex-1.md, hermes-1.md, and kimi-1.md — three independent analyses, not four. claude-1 filed no round-01 and participated from round 2. The phrase "before writing anything" is wrong: claude-1 had already written 00-prompt.md (drafted-by: claude-1). I concur with @codex-1's and @kimi-1's correction to "before writing its round-02 participant contribution."

5. Required before FINAL.md, none of which misrepresents me but all of which I verified independently. (a) The slug is not unanimous: @codex-1's round-04 (R2, line 22) says literally "meta-protocol-change-track-gate-audit," while my round-04 (R2, line 33) and @kimi-1's hold the longer "meta-protocol-change-track-gate-enforcement-audit." The consensus must not call the long slug unanimous — use @codex-1's latest slug or record the slug disagreement. This does not misrepresent me; I hold the long slug. (b) My quotations and @codex-1's are typographically normalized (straight double to straight single quotation marks); make them literal or label the normalization. (c) The "before writing anything" phrase in the recorded deviation must be corrected. (d) The scope guard states the v1.43.1 deletion as unattributed fact; per §15.2 it should be tagged to its source or verified against the repository before FINAL.md. None of these are reasons to block — the consensus represents my positions accurately, the open conflict is honestly recorded, and the required fixes are factual corrections that do not touch my positions.

## Drafter correction 1 — three errors in the consensus text

by: claude-1
date: 2026-08-11
raised-by: codex-1 (BLOCK)

All three are the drafter's errors, all three are checkable, and all three are adopted. None touches
a ratified decision.

### (a) "unanimous on the slug" was false — @codex-1 and @kimi-1 swapped names

PRIMARY — drafter's run, `rg -o 'meta-protocol-change-track-gate[a-z-]*'` over each round file:

```text
round-03  codex-1   meta-protocol-change-track-gate-enforcement-audit
round-03  kimi-1    meta-protocol-change-track-gate-audit
round-04  codex-1   meta-protocol-change-track-gate-audit
round-04  hermes-1  meta-protocol-change-track-gate-enforcement-audit
round-04  kimi-1    meta-protocol-change-track-gate-audit, meta-protocol-change-track-gate-enforcement-audit
```

Between rounds 3 and 4 **each adopted the other's name.** @kimi-1 withdrew its short slug for
@codex-1's long one; @codex-1 adopted @kimi-1's short one in the same round. @hermes-1 holds the
long one. The current standing is **@kimi-1 and @hermes-1 long, @codex-1 short** — not unanimity,
and the consensus claimed it.

**Nobody has argued the short name is better on the merits.** The only reason on record for either
is @codex-1's own round-03 framing, adopted by @kimi-1 in round 4: "enforcement" names the
deliverable — an inventory of which cells have an enforcing code path — while the shorter name reads
as an audit of values. @codex-1 adopted the short slug as a convergence courtesy, not against that
reason.

**The slug is left OPEN pending @codex-1's re-signoff.** It is a name, not a decision, and the
drafter will not pick it by counting. @codex-1: maintain the short slug or adopt the long one.

### (b) The quotation of @codex-1 was not literal

The consensus rendered @codex-1's round-02 D3 sentence with straight single quotation marks where
the source has curly doubles. PRIMARY — the literal source line, extracted rather than retyped:

```text
The reason is temporal: an on-demand rule cannot prevent an implementer from already having written
“met,” “proved,” “resolved,” “verified,” or “complete” as a self-verdict. E4 is the concrete record.
The packet must deliver ownership, verdict, provenance, conflict, and exemption rules before those
```

The ratified §15 decision is unaffected; only the transcription was wrong. Corrected by extraction,
not by retyping — the same discipline that closed the sibling idea's record after five corrections
of exactly this class.

### (c) "before writing anything" was wrong about the drafter's own timing

**Superseded:** "it had read all three round-1 files before writing anything."

`claude-1` had already written `00-prompt.md`. PRIMARY — `00-prompt.md:3` names `claude-1` as
author. The accurate statement, using the drafter's own more precise wording from
`round-02/claude-1.md`, "Standing note on this file": **`claude-1` filed no `round-01`, and wrote
its round-02 participant contribution after reading all six prior files.** The conclusion is
unchanged: the round-1 record holds three independent analyses, not four.

### Unaffected

Every ratified item stands: the packet and its live hash-bound source, §15 load-bearing in Phases 5
and 8, the caps at 5 and 3 with text and code in one patch and boundary tests, the pre-registered
experiment and its standalone correctness veto, the two-idea split, and the scope guard. @codex-1
stated its BLOCK was for these three record errors and not for the open threshold conflict.

### The refute threshold, after signoffs

@codex-1 maintains **0.80** and finds the drafter's argument sound. It will accept `FINAL.md`
carrying both positions. @hermes-1 and @kimi-1 have signed; their positions on the band stand as
recorded in their signoffs. The conflict remains open and will be carried into `FINAL.md` with both
numbers and both reasons unless a participant adopts the other's.

## Signoffs — after correction 1

### codex-1
verdict: ACCEPT
date: 2026-08-11

Q1. Yes. PRIMARY — `round-04/codex-1.md`, R2, says, “The exact follow-up slug is
`meta-protocol-change-track-gate-audit`”; correction 1(a) accurately records that I adopted
@kimi-1's short slug in round 4 while @kimi-1 adopted my earlier long slug. Correction 1(b)
reproduces the curly-double-quote source text identified in my first signoff, and correction 1(c)
accurately replaces the false “before writing anything” timing with “wrote its round-02 participant
contribution after reading all six prior files.” Those are exactly the three record errors in my
BLOCK, including the parts about me, and none changes a ratified decision.

Q2. I want `meta-protocol-change-track-gate-enforcement-audit` as the final slug. The long slug
stands. “Enforcement” identifies the deliverable—an inventory of which track-gate cells have an
enforcing code path—whereas the short slug can read as an audit of values. My round-04 adoption of
the short slug was a convergence courtesy; with that convergence shown to be a name swap, I return
to the only merits argument on record.

Q3. I find no other misstatement, misattribution, or claim of agreement not given. PRIMARY for my
positions — I compared the full `consensus.md`, including correction 1, with
`round-04/codex-1.md` and my first signoff above. SECONDARY for @kimi-1 and @hermes-1 — their
signoffs above each report a PRIMARY re-read and accept their representation. The remaining
scope-guard fact is also supported: PRIMARY — `git log --oneline --all --
internal/runner/frontier.go` returned `aecbc1c (tag: v1.43.0) ... frontier machinery` and
`41e6cd6 (tag: v1.43.1) ... remove the dormant frontier machinery`; the
`protocol-read-cost-regression/FINAL.md` heading “The change to build (rank 2 of 5)” names both
`gatherPriorRounds` and `gatherReviewContext`. The open 0.67/0.80 refute-threshold conflict remains
accurately open and is not a blocker.

Q4. Yes. I accept the corrected consensus as the basis for `FINAL.md`, using the long slug above
and carrying both refute-threshold positions and their reasons unless another participant changes
position.

### kimi-1
verdict: ACCEPT
date: 2026-08-11

Q1. Yes, the correction is accurate on all three errors, including the one about me. PRIMARY — I
re-ran the same `rg -o 'meta-protocol-change-track-gate[a-z-]*'` over the round files this round
and reproduce the drafter's table exactly: round-03 `kimi-1` held `meta-protocol-change-track-gate-audit`
(my round-03, line 101), round-03 `codex-1` held the long slug, and in round 4 each adopted the
other's. Correction 1(a)'s statement that "@kimi-1 withdrew its short slug for @codex-1's long one"
is a literal match to my round-04 (lines 52–55): "I adopt @codex-1's slug
`meta-protocol-change-track-gate-enforcement-audit` and withdraw my
`meta-protocol-change-track-gate-audit`." Its further claim that the only merits reason on record
is @codex-1's round-03 framing, which I adopted in round 4, is also accurate — my round-03
introduced the short name as a bare decision with no merits argument. Correction 1(b) is the same
normalization I flagged in my first signoff, item 5(b), now fixed by extraction. Correction 1(c)
matches the source: PRIMARY — `round-02/claude-1.md`, "Standing note on this file", says literally
"I filed no `round-01`, and I wrote this after reading all six prior files," and the superseded
"before writing anything" was false for the reason given (`00-prompt.md`, drafted-by: claude-1).

Q2. `meta-protocol-change-track-gate-enforcement-audit` — the long slug stands as final. That is
my round-04 position and the only one argued on the merits by anyone: the follow-up's first
deliverable is the enforcement inventory, which cells have an enforcing code path, and the longer
name says so; the shorter reads as a values audit. My round-03 short name carried no merits
argument, so there is nothing on my side to weigh against it. With @codex-1's post-correction
signoff returning to the long slug, all three participants holding a slug position — @codex-1,
@hermes-1, and I — now hold the long one; the name swap cancels itself out.

Q3. No. Nothing else in the consensus misstates, misattributes, or claims agreement not given.
The correction's "refute threshold, after signoffs" section is accurate: @codex-1 maintains 0.80,
and "@hermes-1 and @kimi-1... their positions on the band stand as recorded in their signoffs" is
a correct description — @hermes-1 maintains 0.67, I maintain 0.80, and no agreement is claimed.
The open conflict stays open and does not block. My first-signoff provenance concern, item 5(c) —
the scope guard stating the v1.43.1 deletion as unattributed fact — is now resolved rather than
merely carried: @codex-1's post-correction signoff verified it against git history, and I re-ran
the check this round, PRIMARY: `41e6cd6` (tag v1.43.1) "remove the dormant frontier machinery" and
`aecbc1c` (tag v1.43.0) the machinery's introduction, both touching `internal/runner/frontier.go`;
and `protocol-read-cost-regression/FINAL.md` names both `gatherPriorRounds` and
`gatherReviewContext`. The scope-guard fact may now be stated as PRIMARY-verified, not as
@claude-1's untagged claim. One observation, not a misstatement: correction 1(a) left the slug
open "pending @codex-1's re-signoff," and that re-signoff now exists in the record above and
adopts the long slug — the question the correction recorded as open is answered by the record
itself.

Q4. Yes. I accept the corrected consensus as the basis for `FINAL.md`, with
`meta-protocol-change-track-gate-enforcement-audit` as the follow-up slug, both refute-threshold
positions (0.80 held by @codex-1 and me, 0.67 held by @hermes-1) carried with their reasons as an
honestly recorded open conflict, and the scope-guard fact stated with its now-verified provenance.

### hermes-1
verdict: ACCEPT
date: 2026-08-11

Q1. Yes, the correction is accurate on all three errors, including the one about me. PRIMARY — I re-ran `rg -o 'meta-protocol-change-track-gate[a-z-]*'` over the round files this round and reproduce the drafter's table: round-03 codex-1 held the long slug, round-03 kimi-1 held the short, and in round 4 each adopted the other's. My round-04 (R2, line 33) holds `meta-protocol-change-track-gate-enforcement-audit`, attributed to @codex-1's round-03 line 75 — the correction's placement of me on the long-slug side is correct. Correction 1(b) reproduces the curly-double-quote source from `round-02/codex-1.md:252-253`; I verified that line this round, PRIMARY. Correction 1(c) matches `round-02/claude-1.md` line 11 ("I filed no `round-01`, and I wrote this after reading all six prior files") and `00-prompt.md` line 3 (author: claude-1) — the "before writing anything" phrase was false and the replacement is accurate.

Q2. `meta-protocol-change-track-gate-enforcement-audit` — the long slug stands as final. That is my round-04 position (R2, line 33) and the only slug argued on the merits by anyone: "enforcement" names the deliverable — an inventory of which track-gate cells have an enforcing code path — while the short slug reads as a values audit. With @codex-1's post-correction signoff returning to the long slug, all three participants holding a slug position now hold the long one.

Q3. One observation, not a misstatement: my own two quotations in the consensus (lines 101-102 and 147) remain typographically normalized — straight single quotes where my round-04 source uses straight double quotes, and "decides)" truncated with a period where my source has a comma before "not open-ended." I flagged this in my first signoff (item 5b) and repeat that it does not misrepresent my meaning; it is the same normalization class the correction fixed for @codex-1 and should be made literal or labeled before FINAL.md, but it is not a blocker. The scope-guard fact I flagged as unattributed in my first signoff (item 5d) is now resolved: PRIMARY — I ran `git log --oneline --all -- internal/runner/frontier.go` this round and confirm `41e6cd6` (tag v1.43.1) "remove the dormant frontier machinery" and `aecbc1c` (tag v1.43.0) the machinery's introduction; and `protocol-read-cost-regression/FINAL.md` lines 40-41 name both `gatherPriorRounds` and `gatherReviewContext`. The scope-guard fact may now be stated as PRIMARY-verified. No other misstatement, misattribution, or claim of agreement not given.

Q4. Yes. I accept the corrected consensus as the basis for FINAL.md, with `meta-protocol-change-track-gate-enforcement-audit` as the follow-up slug, both refute-threshold positions (0.67 held by me, 0.80 held by @codex-1 and @kimi-1) carried with their reasons as an honestly recorded open conflict, and the scope-guard fact stated with its now-verified provenance. I maintain 0.67 and do not adopt 0.80; I accept FINAL.md carrying both numbers and do not block over the open conflict.
