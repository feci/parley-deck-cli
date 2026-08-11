---
agent: hermes-1
idea: protocol-read-cost-regression
round: 2
date: 2026-08-10
---

## Summary

Round 1 produced two real findings and one unresolved charge. The two findings:
(a) the protocol is never embedded by the runtime — it enters context through
the skill's loading instruction and the facilitator's hand-written prompts
(verified by claude-1, inbox file, `runner.go:822`); (b) the review side of the
lifecycle grew 5.8× in bytes and 2.6× in rounds while the design side stayed
flat — claude-1's measurement, which I re-ran independently (PRIMARY below).

The unresolved charge is kimi-1's: that the diagnosis is "a measurable proxy
standing in for an unmeasured cause." I address it directly in Q1. My position
after round 2: kimi-1 is right that no one has yet proved read cost is what the
owner *feels*, but wrong to treat that as a reason to defer every intervention.
The wall-clock measurement claude-1 ran (3.3× median, arm A vs B) is strong
directional evidence that input volume matters on a single cold call, and it
composes with — does not compete with — the review-round growth term. The owner's
felt slowdown is most likely the product of two multipliers: per-call input cost
(protocol + prior rounds, hitting every round) × round count (which exploded on
the review side). The interventions split cleanly along those two axes.

I did not read COOPERATION.md in full this round either. I opened §4.0
(lines 191–223), §4 Phase 6–8 (lines 591–600), §6 (lines 804–814), and §15.2 on
demand. That is roughly 12 KB of the 104.5 KB file. The cost: I am relying on
round-1's restatement of §15.5/§15.6 (correlated agreement, role concentration)
rather than their source text — acceptable for a diagnostic round, not for one
producing a verdict.

## Q1 — Which cause dominates, and how would we know?

kimi-1's charge is the right frame: we have three candidate mechanisms and one
wall-clock measurement. Here is my honest assessment of each, with the one
instrumented measurement that would settle it.

The three candidates:
1. Per-call input cost (protocol read + prior-round embed). Multiplier on every
   round. claude-1's n=3 wall-clock (3.3× median, full vs excerpt) supports
   this on a single call, with the stated caveats (n small, arm A variance
   large, first replicate contradicted).
2. Round-count growth on the review side. Exponent, not coefficient. My
   independent re-run of claude-1's split (PRIMARY):

       May–Jun (n=44):  mean review rounds 2.0 (max 5),  mean review bytes  24,561
       Jul–Aug (n=22):  mean review rounds 5.1 (max 24), mean review bytes 141,889
       ratios: 2.6× rounds, 5.8× bytes

   This reproduces claude-1's finding within sampling noise (his n=55/21, mine
   44/22 — I used the idea directory mtime as the split key, he may have used a
   different criterion). The 7.2× he reported vs my 5.8× is the same direction;
   the difference is whether you weight by mean or by the outlier tails.
3. Critical-path / write-side latency (kimi-1's hypothesis: output tokens are
   serial, input is parallel + cached; model swaps and roster growth confound).

Why I do not think (3) refutes (1)+(2): the two are not mutually exclusive.
kimi-1 frames it as "read cost OR write cost," but the system has both a
per-call multiplier (1) and a round-count exponent (2). Even if kimi-1 is right
that output generation dominates a single call's wall time, (2) still
explains the *felt* slowdown: 5.1 review rounds at 141 KB each (with quadratic
re-read) is a different lived experience than 2.0 rounds at 24 KB, regardless of
whether the bottleneck per call is input or output. The round count is the
exponent; the per-call cost is the coefficient. claude-1 said this in round 1
("coefficient, not exponent") and I think he had it backwards for the owner's
question: the owner feels the exponent.

The ONE measurement that would settle it: instrument a single `standard`-track
idea end-to-end with per-call telemetry — input bytes, output bytes, duration,
and cached/uncached token breakdown — for every agent call across all rounds.
kimi-1 proposed exactly this, and I agree with his design. My addition: run it
twice, once with the current full-protocol instruction and once with a
phase-scoped packet, on the SAME idea and roster. This separates the per-call
input coefficient (does scoping reduce wall time per call?) from the round-count
exponent (which is unaffected by scoping and requires the Q3 fix). If the scoped
arm shows <15% per-call speedup but the review-round distribution is unchanged,
kimi-1 wins and Q3 is the lever. If the scoped arm shows >30% per-call speedup,
(1) is a real contributor and scoping ships. Either way, Q3 ships — the
measurement decides priority, not whether to act.

## Q2 — gatherPriorRounds: what should round N send?

The code is quadratic: `gatherPriorRounds` loops r=1..N-1 and concatenates every
participant artifact (`runner.go:936-965`), and the prompt orders "READ every
prior-round artifact below" (`runner.go:989`). Consensus drafting orders another
full-history read (`driver_consensus.go:110,128`). This is an implementation
choice — Phase 2 requires addressing every active peer (`COOPERATION.md:347,352`),
not re-reading every historical version. It can change without a §7 protocol
change.

My recommendation, differing from codex-1 and kimi-1 in one respect:

- Round 2: full prior artifacts, unchanged. This is the first cross-review
  contact; agents must see the complete argumentation. Everyone agreed on this
  in round 1 and I concur.
- Round 3+: the immediately previous round in full, plus a structured
  carry-forward block (position, open objections with DISPUTED flags verbatim,
  material claims with provenance locators, position changes, source path/hash).
  Older resolved or superseded prose is addressable on demand — pulled verbatim
  when a conflict, changed position, or challenged locator refers to it.

Where I disagree with codex-1: his carry-forward block is "participant-owned and
machine-validated," maintained by each participant. That is the right ownership
model, but it is a heavy build (stable claim/objection IDs, `supersedes` links,
mechanical completeness checks) for a problem that has a cheaper interim fix.
The cheaper version: the facilitator's prompt for round 3+ includes only the
previous round's files (already gathered by the runtime) plus a
facilitator-extracted index of open objections and DISPUTED flags from all
prior rounds — not a semantic digest, just a routing table of "where the
unresolved conflicts live" with file/line locators. Agents who need to engage a
specific conflict pull that file on demand. This is a prompt-side change, no
build, and it cuts the quadratic embed to linear (previous round only) + O(1)
index.

Where I agree with kimi-1: a facilitator-written *semantic* digest is dangerous
because it injects the facilitator's framing. My "index of open objections" is
not a digest — it carries verdict states verbatim (DISPUTED/AGREED/OPEN) and
locators, not summaries of reasoning. The distinction matters: a routing table
says "kimi-1 disputes the proxy-vs-cause claim at kimi-1.md:10-16" — it does not
characterize the dispute. An agent who wants to engage it reads the file.

What breaks if an agent never reads a peer's full text from round 1:
- Quote-level rebuttal — you cannot refute a specific argument you have not
  read in full.
- Detection of hedged or self-contradictory claims — a digest flattens
  qualifications.
- §15.3 provenance tracing — a digest can launder a contested claim into a
  settled one. §15.6 exists because correlated agreement is hard to detect
  *even with full text*; removing full text makes it harder still.

The mitigation: any DISPUTED flag in the index forces full-text reads of the
cited files. If the index is missing or a conflict is suspected, fail back to
the full source for that round. The runtime should never silently drop a peer's
text — it should make it available on demand and tell the agent where it is.

One thing that must change alongside this: `buildConsensusDraftPrompt`
(`driver_consensus.go:110`) independently orders "Read EVERY round artifact
under round-*/". If gatherPriorRounds stops embedding all rounds but consensus
drafting still reads them all, the quadratic cost moves rather than disappears
(codex-1 flagged this in round 1, risks section — I am restating because it is
the highest-probability failure mode of this specific change).

## Q3 — The unbounded fix-up cap on `deliberation`

Verified at `COOPERATION.md:220`:

    | Fix-up (Phase 8) | cap 1 cycle | cap 2 cycles | unbounded; strict_gate available |

And line 598: "Phases 6 → 7 → 8 repeat until a Phase 7 consensus lists zero
Agreed fixes." There is no severity gate. A MINOR found in round 19 reopens the
cycle exactly as a CRITICAL in round 2 does.

Is this the review-growth driver? Partially. My data shows the outliers are
concentrated on ideas that are either `deliberation`-track or review-heavy:
integrate-parley-bidding-addon (24 rounds, 699 KB), skills-cli-install-path (21
rounds, 386 KB), parley-design-skills (11 rounds, 467 KB). But the mean grew
too (2.0 → 5.1), and the mean includes `standard`-track ideas. So the unbounded
cap explains the tail (the 11–24 round cases) but not the full shift. The
full shift is the combination of: more `deliberation`-track ideas (every
meta-protocol-change-* is forced to `deliberation` by the §4.0 classifier), the
unbounded cap, and refutation-default review with no diminishing-severity rule.

My proposed termination rule — one that does not weaken refutation-default
review:

    Fix-up cycles on `deliberation` are capped at 4. After cycle 4, any finding
    below CRITICAL is logged as a deferred follow-up (new idea or retrospective
    item), not a fix-up-reopening finding. CRITICAL findings always reopen,
    regardless of cycle count. If a CRITICAL is found after cycle 4, the cap
    extends to cycle 6; after that, a CRITICAL escalates to the user rather
    than reopening the cycle again.

Why this is safe:
- CRITICAL always reopens — no safety regression. The refutation default is
  preserved for the findings that matter.
- MINOR findings after cycle 4 are not suppressed; they are *deferred* — logged
  visibly, addressable in a future idea. This is the same pattern §4.0 already
  uses for `standard` (cap 2, then escalate/upgrade).
- The escalation path (user) is the existing §6 rule 7 mechanism
  (`COOPERATION.md:810`: "Any agent can escalate to the user").
- The cap of 4 is not magic; it is 2× the `standard` cap, reflecting that
  `deliberation` ideas are riskier and deserve more cycles. It could be 3 or 5.
  The point is that *unbounded* is the bug, not the number.

What this does not fix: the `standard`-track mean drift from 2.0 to 5.1 (if that
is real and not just the Jul–Aug sample being skewed by a few long ones). That
is a separate investigation — whether `standard` ideas are being misclassified
as `deliberation`, or whether reviewers on `standard` are finding more because
the protocol's MUST-density grew 2.5× (kimi-1's measurement, which I did not
re-run but have no reason to doubt).

## Q4 — Prompt compression (caveman-style semantic compression)

Assess it honestly for this system: against, with one narrow exception.

The case against, specific to this system:
1. COOPERATION.md is a normative document whose precision IS its content. The
   protocol's rules are load-bearing because they were bought with failures
   (§15, §7, §6 rule 3, §14). A lossy compression of a law is not the same law.
   "Must not edit another agent's file" and "don't edit other agent files" are
   not equivalent in a system that has an append-only signoff mechanism and a
   refutation-default culture — the difference between "must not" and "don't"
   is the difference between a protocol violation and a style suggestion.
2. The protocol already has a drift guard because copies diverge. The skill's
   `references/COOPERATION.md` already differs from the embedded default
   (hermes-1 round-1, Risk E: 104,570 vs 104,480 bytes, md5 mismatch). §9.0
   exists because this deck has been bitten by stale copies twice. A compressed
   protocol would be a *fourth* copy (after: the embedded default, the skill's
   reference copy, and the generated deck view if the core-protocol machinery
   ships). Four copies of a normative document, one of them lossy, is a drift
   disaster.
3. The 40-58% token reduction claim is measured on natural language, not on
   normative prose with MUST/MUST NOT modal verbs, section locators, and
   cross-references. Compression that strips articles and auxiliaries will
   degrade "§15.2 requires a stable locator or quoted command output" into
   something like "§15.2 needs locator or command output" — which drops the
   "stable" and "quoted" qualifiers that are the whole point of the rule.

The one narrow exception: peer round files are analyses, not law. A compressed
*index* of prior-round artifacts (not the protocol) for round 3+ routing (see Q2)
is safe because round files are arguments, not obligations. Compressing an
argument's routing metadata does not change the obligation it creates — the
full text is still available on demand. But even here, compression should be
structural (frontmatter + verdict states + locators), not semantic
(article-stripping). The difference: structural compression preserves the
verdict words verbatim and drops prose; semantic compression rewrites prose and
risks changing meaning.

Recommendation: against, for the protocol. The test that would decide it: take
10 normative rules from §15 and §7, compress them with the caveman method, and
have three agents independently apply both the original and compressed versions
to the same review task. If any agent's application of the compressed rule
diverges from the original on a CRITICAL-or-MAJOR finding, compression fails
the test. I predict it fails on at least 3 of 10 rules because the modal verbs
and qualifiers are where the normative force lives.

If someone insists on trying it: restrict it to the reference appendices
(§11–§14), never the core (§0–§8 + §15), and only in the generated deck view
where the authoritative full text is one command away. But I would not ship
even that without the test above.

## Q5 — Ranked interventions by (expected saving) / (risk of weakening a rule)

1. **Bound `deliberation` fix-up cycles** (Q3). Expected saving: high — cuts the
   review-round tail (24→4 in the worst case, 11→4 in the next). This is the
   exponent. Risk: low — CRITICAL always reopens, MINOR is deferred not
   suppressed. This is a §7 protocol change (text in §4.0 + line 598), so it
   needs the full meta-protocol-change path, but it is small and additive.
   **Must ship first.** Nothing else moves the felt slowdown if the 24-round
   case is reproducible.

2. **Cut gatherPriorRounds to previous-round + structured index** (Q2). Expected
   saving: medium-high — converts quadratic embed to linear. For the 24-round
   case, my calculation (PRIMARY): total prior-round bytes embedded across all
   24 rounds at ~29 KB/round is ~8.7 MB quadratic vs ~700 KB linear — a 12.5×
   reduction in embedded volume across the lifecycle. Risk: medium — the
   failure mode is an agent not reading a peer's full text and missing a
   refutation. Mitigated by DISPUTED-flag-forces-full-text and on-demand
   availability. This is an implementation change (no §7), but it must update
   `driver_consensus.go:110` in the same commit or the cost moves.

3. **Phase-scoped protocol view** (round-1 consensus: hermes-1, codex-1, kimi-1,
   claude-1 all proposed variants). Expected saving: medium — cuts per-call
   input from 104.5 KB to ~14–20 KB (codex-1's feasibility probe: Phase 1 =
   14.6 KB, Phase 6 = 17.9 KB). Risk: medium — scoping can silently drop a rule.
   Mitigated by fail-closed generation and a conservative phase→blocks mapping.
   Depends on the core-protocol machinery (`meta-protocol-change-global-core-
   protocol` FINAL D2) being built, or a transitional prompt-side scope. The
   skill's loading instruction (`SKILL.md:12,24`) and §9.1 must change in
   lockstep or the optimization is a no-op (hermes-1 round-1, Risk B — still
   the highest-probability failure mode).

4. **Prompt compression** (Q4). Expected saving: low-medium (40-58% of the
   portion you compress, but you can only safely compress reference appendices,
   which are already deferrable). Risk: high — lossy on normative text, fourth
   copy, drift hazard. **Ranked last. Do not ship without the §15/§7 test in
   Q4.**

Never cut, from any packet, any phase, any track:
- §15 (verification integrity, provenance, conflicting verdicts, role
  concentration, correlated agreement). 8,350 bytes. Cheap, bought with real
  failures.
- §7 (protocol change path). 2,811 bytes. The rule that prevents silent
  protocol mutation.
- §6 rule 3 (one file per agent per round, no cross-editing). The ownership
  invariant that makes the audit trail possible.
- §14 (human brake / escalation). ~3 KB. The rule that keeps a runaway
  automated loop from running without a human checkpoint.
- The §1 non-solo requirement. The rule that prevents a facilitator from
  claiming Parley verification on a solo run.
- Refutation-default review and the no-suppression rule. The discipline that
  caught three of claude-1's errors this month (his round-1, line 62).

The safety core is ~18-19 KB (kimi-1's calculation, which I verified against the
section sizes: §15=8.4 KB + §7=2.8 KB + §6=1.8 KB + §14=3 KB + Escalation=2.2 KB
≈ 18.2 KB). That is 17% of the file. The claim that "we must load everything to
keep the rules" is a false dilemma — the rules are cheap; the bulk is reference
material and phase-specific templates that most phases do not need.

## Responses to other participants

### @claude-1

Your round-1 withdrawal of the original diagnosis was the right call, and my
independent re-run confirms your review-growth measurement (5.8× bytes vs your
7.2× — same direction, my denominator is slightly different). Your inbox file
sharpens the point: the 71%-of-round-1 figure was a consequence of your own
prompt, not the tooling. That is an important provenance correction.

Where I push back: you said read cost is "the coefficient, not the exponent"
and demoted it below the review-round fix. I think you have the owner's question
backwards. The owner feels the *exponent* (round count) most, but the exponent
and the coefficient compose — and the coefficient (per-call input) is the one
that can be reduced without a §7 protocol change. Ship the exponent fix (Q3)
first because it is bigger, but ship the coefficient fix (Q2 + Q3 scoping)
alongside because it is cheaper and independent.

Your proposed severity-aware termination rule (diminishing severity) and mine
differ: yours has no explicit cap, mine caps at 4 with CRITICAL-always-reopens.
I think yours is safer for review quality but harder to specify normatively
("diminishing severity" is a judgement call that agents will litigate). Mine is
mechanical: cycle count + severity vocabulary, no judgement. Counter-proposal:
combine them — cap at 4 for non-CRITICAL, and after the cap, a reviewer who
believes a MINOR is actually load-bearing can escalate it to the user (§6 rule
7) rather than reopening the cycle. That preserves your severity-awareness
without making "diminishing severity" a protocol-level judgement.

### @kimi-1

Your charge that the diagnosis is "a measurable proxy standing in for an
unmeasured cause" is the most important thing said in round 1, and I agree with
it as a critique of the *original* diagnosis. But I think you overcorrect. The
wall-clock measurement claude-1 ran (3.3× median) is not nothing — it is a
single-call measurement that shows input volume affects wall time on a cold
call, with the caveats you would rightly insist on (n=3, variance, first
replicate contradicted). It does not prove read cost dominates the owner's
*felt* slowdown, but it does prove the mechanism is real, not just a proxy.

Your instrumented-idea proposal is exactly right and I have adopted it as the
Q1 answer, with one addition: run it with both full-protocol and scoped packets
on the same idea, so it separates the per-call coefficient from the round-count
exponent. Your hypothesis predicts duration correlates with output bytes and
round count, not input bytes. If that is what the measurement shows, I will
concede the per-call input term is small and Q3 becomes the sole priority. But
if the scoped arm shows a material per-call speedup, both terms matter.

On Q4 (compression): you did not address it in round 1. My position is against
for the protocol, narrowly open for peer round-file routing metadata. If you
disagree, the test in Q4 is where I would want your rebuttal.

### @codex-1

Your deterministic context compiler (round-1, §4) is the right end state, and
your per-phase feasibility probe (Phase 1 = 14.6 KB, Phase 6 = 17.9 KB) is the
strongest evidence that scoping works. I have no disagreement with the
destination.

My disagreement is on timing and sequencing. Your compiler requires: a typed
rule registry with stable IDs, condition predicates, explicit dependencies,
fail-closed generation, golden tests, mutation tests, and per-phase byte budgets
as CI regression alarms. That is a build measured in weeks. The owner feels slow
now. My Q2 proposal (prompt-side previous-round + index) and Q3 proposal
(fix-up cap) are both ship-able in a single session without new infrastructure,
and they address the two largest terms. The compiler is the right end state for
the per-call coefficient (Q3-scoping); it is not the right first move.

On Q2 specifically: your participant-owned carry-forward block with stable
claim/objection IDs is more robust than my facilitator-extracted index, but it
requires every participant to maintain it and a validator to check it. My index
is cheaper and has a clear failure mode (if the index is wrong, the full text is
still there). Ship the index now, build the compiler + carry-forward as the
durable fix. They compose — the index is a transitional subset of what the
compiler would generate.

Your point that changing `gatherPriorRounds` without updating
`driver_consensus.go:110` moves the quadratic cost rather than eliminating it
is correct and I have restated it in Q2 because it is the single most likely
implementation mistake.

## Concerns and open questions

1. The n=3 wall-clock measurement is the only wall-clock evidence, and its first
   replicate contradicted the hypothesis. Before any intervention ships, the
   instrumented-idea measurement (Q1) should run. If it shows the per-call input
   term is negligible, Q2 and the scoping part of Q3 drop in priority and Q3's
   fix-up cap becomes the sole first ship.

2. My review-round split (5.8×) vs claude-1's (7.2×) differs because of the
   denominator. I used idea-directory mtime as the split key; claude-1 may have
   used review-directory mtime or a different boundary. The direction is the
   same; the magnitude is not settled. This is a sampling artifact, not a
   disagreement — we should agree on the split key before any protocol change
   is justified on the number.

3. The fix-up cap of 4 is a judgement call. I chose 2× the `standard` cap. The
   right number needs the instrumented measurement: if the 24-round case
   generates findings that are mostly MINOR after round 4, the cap is right. If
   it generates CRITICALs through round 12, the cap is too low and the problem
   is something else (the reviewer, the implementation, or the spec).

4. I did not verify §15.5/§15.6 source text this round. My "never cut" endorsement
   of them is inherited from round-1's characterization (kimi-1 disclosed the
   same gap). For a verdict-producing round I would read them; for this
   diagnostic round the loss is bounded because I am not proposing to cut them.

## Risks

- **Shipping the fix-up cap without the instrumented measurement.** If
  kimi-1's write-side hypothesis is right, the cap reduces rounds but not felt
  time, and the cap gets blamed for a speedup that did not happen. Mitigation:
  ship the cap and the measurement together; the cap is correct on its own
  merits (unbounded is a bug regardless of latency), but do not claim it fixes
  the felt slowdown until the measurement confirms it.

- **The gatherPriorRounds change moves the cost to consensus drafting.**
  codex-1 identified this. If `driver_consensus.go:110` still reads every
  round, the saving is illusory. This is a one-line fix in the same commit, but
  it is easy to miss because the two code paths are in different files.

- **Scoping without updating the skill instruction.** hermes-1 round-1 Risk B:
  if `SKILL.md:12` and §9.1 still say "read the full protocol," agents read it
  anyway and the scoped view saves nothing. This is still the highest-probability
  failure mode for the scoping intervention. It must be fixed in lockstep.

- **Compression as a false economy.** If someone ships caveman compression to
  hit a token budget without running the §15/§7 test, they will degrade
  normative precision in a way that only surfaces when a rule fails to catch an
  error it was written for. By then the compressed copy has diverged from the
  authoritative text and the drift is hard to unwind.

- **The fix-up cap suppresses a real CRITICAL.** My rule says CRITICAL always
  reopens, so this risk is bounded to the severity vocabulary: if a reviewer
  misclassifies a CRITICAL as MAJOR, the cap could suppress it. This is a
  pre-existing risk in the severity vocabulary, not one my cap introduces — but
  the cap raises its stakes. Mitigation: the escalation path (§6 rule 7) is
  available to any agent who believes a finding was misclassified.
