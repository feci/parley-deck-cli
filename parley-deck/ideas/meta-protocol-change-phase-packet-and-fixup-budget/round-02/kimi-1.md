---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 2
date: 2026-08-11
responding-to: [codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Evidence base for this round

Everything below rests on commands I ran this round. Quoted outputs are verbatim.

**E1 — the fix-up-cycle distribution exists and nobody used it.** Command:

```text
$ for f in parley-deck/ideas/*/IMPLEMENTATION.md; do grep -c '^## Fix-up cycle' "$f"; done | sort -rn | awk ...
n=69 sum=139 mean=2.01449 median=1 p90=4 max=25 min=0
distribution: 0×17, 1×34, 2×7, 3×2, 4×3, 5×2, 9×1, 14×1, 15×1, 25×1
```

(PRIMARY — command and output quoted; n=69 = idea dirs containing an `IMPLEMENTATION.md`.
Caveat: a 0 may mean "no fix-up needed" or "predates the cycle-section convention"; counts are
as-recorded, a lower bound if anything.)

**E2 — the recent regime is heavier.** Pairing each idea's `created:`/`date:` frontmatter with its
cycle count and taking the latest 13 (2026-07-18 → 2026-08-10), same command family (PRIMARY, raw
output quoted):

```text
2026-07-18 5 composite-agent-naming-and-roster-reinit
2026-07-28 14 parley-design-skills
2026-07-29 15 skills-cli-install-path
2026-07-29 25 integrate-parley-bidding-addon
2026-07-29 4 readme-skill-catalogue
2026-08-01 4 addon-manifest-coverage
2026-08-04 0 meta-protocol-change-verification-integrity
2026-08-06 0 kimi-opencode-full-adapters
2026-08-06 0 skill-sync-cli-1-39
2026-08-06 5 roster-operations-standard
2026-08-07 0 protocol-overlay-local-extension
2026-08-07 9 meta-protocol-change-global-core-protocol
2026-08-10 3 protocol-read-cost-regression
```

Sorted: 0,0,0,0,3,4,4,5,5,9,14,15,25 → median 4, mean 6.08 (my arithmetic on PRIMARY data).

**E3 — the tail idea's anatomy.** `integrate-parley-bidding-addon` (deliberation, per its
00-prompt): 3 cross-review rounds, **24 review rounds**, **25 fix-up cycles** (PRIMARY —
per-idea directory counts, `ls -d .../round-*` and `.../review/round-*`; the `track: deliberation`
line read from its 00-prompt). The regression FINAL's "max 24" and "fresh MAJORs at rounds 19–24"
are this idea's **review rounds**. And the two counters are one loop (PRIMARY —
`parley-deck/ideas/driver-impl-phase/IMPLEMENTATION.md:75`): "The fix-up cycle number == the
current review round; bounded by MaxFixupCycles."

**E4 — a printed cap was already ignored.** `skills-cli-install-path` is `track: standard`
(00-prompt.md:5, PRIMARY) — printed fix-up cap 2 — and ran **15** fix-up cycles. A
`grep -in 'cap\|escalat\|extension\|authoriz'` of its IMPLEMENTATION.md returns only "capability"
matches; no cap escalation or authorized extension is recorded there (PRIMARY — grep quoted; I
checked only that file, so the claim is scoped to it).

**E5 — the facilitator's PRIMARY reproduces, and extends to cross-review.** `track.go:150-153`
deliberation → `ApplyOverrides: false, CrossReviewRounds: -1`; `driver.go:100-105` defaults
`CrossReviewRounds < 0 → 1` and `MaxFixupCycles <= 0 → 3`; `track_test.go:47-51` asserts
deliberation preserves exactly that (Cross=1, Fixup=3) (PRIMARY — all three read this round). So
the driver silently defaults deliberation **cross-review to 1** as well, while honoring an explicit
`cross_review_rounds` in 00-prompt with no ceiling (`internal/driver/transport.go:32-34`,
`ReadCrossReviewRounds` default 1; deliberation sets no `CapCrossReviewRounds`). §4.0's table says
"unbounded" for both cells (COOPERATION.md:224-229, read this round).

**E6 — the runner never reads the file, and its prompts carry a paraphrase.** `grep -c
'COOPERATION' internal/runner/runner.go internal/runner/phase58.go internal/app/driver_consensus.go`
→ `0, 0, 0` (PRIMARY, quoted). The prompt builders exist and compose strings: `runner.go:821
BuildRoundOnePrompt`, `:971 BuildRoundPrompt`, `phase58.go:187/202/235/339` (PRIMARY — `grep -n
'func Build'` quoted). `BuildRoundOnePrompt` embeds a hand-written "Rules:" block
(`runner.go:828-833`: "Do not edit any other agent's file.", "Do not read or reference other
agents' round-01 answers."; similar at `:984-987`) — protocol rules restated as Go string literals.
The runner already imports `internal/protocol`, which holds the embedded default
(`internal/protocol/workspace.go:21-22`: `//go:embed defaults/COOPERATION.md` — PRIMARY, read this
round).

**E7 — §15 pointers and binding table.** The "Verification verdicts ... follow §15" header appears
at COOPERATION.md:308 (P1), :328 (P2), :356 (P3), :503 (P6) — and **nowhere in Phase 5, 7, or 8**
(PRIMARY — `grep -n '§15'` quoted; full-file scan). §15.7's table (:1362-1372) binds 15.1/15.2 (and
15.3–15.5) on **every track**. §15.1 (:1259): "**An owner MUST NOT issue a verification verdict on
a claim it owns.**" Phase 8 (:588-600) requires `status: complete` plus `## Validation evidence`;
:660-664: "`MaxFixupCycles` ... are escalation thresholds, not close criteria. Hitting the budget
never marks an implementation complete; it requires human review of the trajectory ..." (PRIMARY —
all read this round).

**E8 — the A/B numbers behind 3.3×.** PRIMARY — stable locator
`parley-deck/ideas/protocol-read-cost-regression/round-02/codex-1.md:124`, quoted from grep:
"full 98.7 s (27.3–105.3) versus excerpt 29.9 s (21.1–39.2), `n=3` each."

One correction to the brief itself: D3 cites "this idea's own cycle-4 record." I cannot locate it —
`ls` of this idea's directory shows only `00-prompt.md`, `round-01/`, `round-02/` (PRIMARY, quoted
output). I do not count it as evidence; D3 below stands on E7 instead.

## D1 — the deliberation fix-up cap: I hold at 5, on new anchors

My round-1 anchor was the 5.1 mean with a noted discomfort that it is tail-inflated. E1/E2 replace
that discomfort with the right dataset. Three conclusions, all arithmetic on E1/E2 (my arithmetic;
the data is PRIMARY):

1. **5, 6, and 8 are empirically identical.** Every observed fix-up count above 5 is in
   {9, 14, 15, 25}. A cap of 5, 6, 7, or 8 escalates exactly the same four ideas on the deck's
   entire history; nothing has ever closed in the 6–8 band. The choice among our three numbers is
   pure out-of-sample hedging, so §15.3's "never by counting" is satisfied trivially: the numbers
   do not differ on any evidence we possess, and must be chosen by error-asymmetry, not by
   averaging anchors.
2. **The asymmetry favors the lower number.** A too-low cap costs one recorded escalation (the
   mechanism all three of us specified, and which :660-664 already defines). A too-high cap gives
   the exact pathology this change exists to interrupt several more cycles of rope. Recent-regime
   median is 4 (E2): cap 5 sits one cycle above the current normal, and the four ideas it would
   have escalated (E2's 9, 14, 15, 25) are precisely the ideas a human should have been looking
   at. ~4/13 recent ideas escalating once is the feature, not the cost — provided each escalation
   is one recorded user decision, which is codex-1's own extension mechanism.
3. **Against the other two anchors specifically.** @codex-1 wrote: "Six is a safety threshold ...
   deliberately just above the measured 5.1-round average" (quoted literally). The 5.1 is the
   tail-inflated **review-round** mean of the recent cohort; anchoring "just above" a tail-inflated
   mean writes the tail into the cap. On the right unit (fix-up cycles, E1/E3 confirm they are the
   same loop), the all-time mean is 2.0 and the recent mean 6.08 — so "just above the mean" argues
   variously for 3, 6, or 7 depending on which mean, which shows the mean is the wrong statistic
   for a right-skewed distribution. The recent **median** is 4; one above it is 5. @hermes-1 wrote:
   "8 is high enough that a converging idea closes before hitting it (most ideas converge in 2-4
   cycles)" (quoted literally). True all-time (median 1) but stale: the recent regime's median is 4
   and its tail starts at 9. In-sample, 8 escalates the same four ideas as 5; out-of-sample, it
   buys three extra cycles only for the churn pattern hermes-1's own clause (c) names the cap a
   backstop against. A backstop belongs just outside the observed normal range, not at 2× it.

**Surviving "the driver already enforces 3 and nobody noticed."** E4/E5 are what settle this for
me: the silent 3 bound only default-config driver runs. Every tail idea (9–25 cycles) ran past it,
so those were hand-driven or explicitly configured runs — and E4 shows a `standard` idea running 15
cycles against a printed cap of 2 with no recorded escalation. **Printed and default caps bind only
where enforcement lives.** So the number is the secondary decision; the primary one is that 5 must
land in `track.go` (`ApplyOverrides: true, MaxFixupCycles: 5` for deliberation) and in the §4.0
cell in the same change, with the escalation payload all three of us specified. A text-only change
re-creates today's divergence under a new number; my round-1 "divergence debt" risk said this, E4
now demonstrates it.

## D2 — the prompt-builder path: it can carry a packet, but it is not where rank 1 lives

@hermes-1's summary puts the change in "the runner's prompt builders" while its own PRIMARY
(confirmed, E6) is that the runner never reads `COOPERATION.md`. Resolving the contradiction is
hermes-1's to do; what I can settle is the question posed to @codex-1 and me.

**Mechanically, the path can carry a packet.** The builders compose prompt strings (E6), and the
runner already imports the package holding the embedded protocol default (`workspace.go:21-22`).
Embedding a generated block is a small, reviewable diff.

**But it carries zero cost relief.** Headless runner-built prompts today contain no reference to
`COOPERATION.md` at all (E6: grep 0,0,0) — driver-managed participants never read the protocol, so
the 3.3× read cost (E8) does not exist on that path. The measured cost arises where instructions
tell an agent to read the file: the skill's standing line, §9's checklist, and hand-written
facilitator prompts — this very idea is hand-facilitated, and my round prompts listed files to
read. **Rank 1's cost lever lives in the three text paths plus `parley protocol packet`, exactly as
locked.** The builder path is not rank 1.

**The builder path is still worth taking — as a correctness fix, not a cost fix.** E6's
"Rules:" block (`runner.go:828-833`, `:984-987`) is a hand-maintained paraphrase of §6 rule 3 and
the Phase-1 independence rule embedded as Go string literals: an unguarded fourth copy of protocol
text, the same stale-copy shape this deck has been bitten by twice. Replacing those literals with
the generated packet removes a drift surface and gives headless agents the real text. Counter-
proposal to hermes-1's framing: keep the builder change, but describe it honestly as
paraphrase-elimination, and do not count it toward the speedup.

**One hard condition.** hermes-1's Q2 has the runner generate from the *embedded default*. The
round-1 lock says packets render from **the single live `COOPERATION.md`**. The embedded default is
a copy; the drift guard is a test, not a runtime gate — hermes-1's own risk 5 concedes the
divergence case. So any runner-side generation must hash-compare embedded against the live deck
file at generation time and fall open to the full resolved protocol on mismatch. With that check, I
withdraw all objection to the builder path; without it, the path reintroduces the stale-copy
failure through the back door.

## D3 — §15 in Phases 5 and 8: codex-1's risk flag beats its own table cell

@codex-1 classified §15 in P5/P8 as "fetched before assigning a material verification verdict" and
then flagged: "If agents routinely make completion verdicts there, §15 should simply become
load-bearing in those packets too" (both quoted literally). Decide for the flag:

- **§15.1 + §15.2 are load-bearing in P5/P8.** Not because agents verdict there, but because the
  two things §15 forbids and requires are exactly what Phase 8's completion contract involves: the
  implementer publishes `status: complete` with `## Validation evidence` (E7, :588-600) — a factual
  claim the implementer **owns**, and §15.1 bars self-verdicts on owned claims (:1259, quoted in
  E7), while §15.2's tags bind on every track (§15.7, :1362-1372). An implementer who does not have
  these two subsections in front of it will write "verified: works" as a verdict instead of
  evidence with provenance. This is also already the convergent position: all three round-1 files
  put 15.1/15.2 in every packet (mine via always-core, hermes-1 via its never-cut list, codex-1 via
  its kernel for P1/2/3/6/7 plus trigger elsewhere).
- **§15.3–15.6 stay on-demand, but the trigger becomes mechanical, not judgment.** The P8 packet
  carries one added line: *your `status: complete` is a claim, not a verdict; it is verified by
  re-review and the checks gate, never by you; any challenge re-opens review under Phase 6/7
  rules.* The omission-index digest quotes §15.4's witness requirement **verbatim** (deviation-log
  entries of the form "obstacle X does not apply because …" are exemption claims), and names §15.3
  with its trigger (a contradicting verdict exists). The digest carries the norm; the full text is
  fetched to act on it.
- The absent §15 pointer in the P5/P8 headers (E7) explains how codex-1's table cell happened; it
  does not change the binding, because §15.7's table is track-scoped, not phase-scoped.

## D4 — the expected saving, and the pre-ship measurement

@claude-1 asked for the number (per the round-2 brief: "the one the owner actually asked for"), and
the honest answer is half "here is a hypothesis" and half "it cannot be estimated in advance — but
it can be measured cheaply before ship."

**Hypothesis (mine, arithmetic on PRIMARY inputs; treat as unverified).** A phase packet under my
round-1 mapping sums to ~308 lines (§4.0 :200-257, §3 :164-195, §6 :732-744, escalation+§8
:686-718+:788-812, §15.1-15.2 :1239-1290, plus one phase's sections — ranges PRIMARY, sum mine) plus
an omission index of at most one digest line per heading (69 headings; `grep -c '^#'
COOPERATION.md` → 69, PRIMARY) → **~350-380 lines ≈ 26% of the 1,372-line file**. Linear
interpolation between the E8 points: 29.9 s + 0.26 × (98.7 − 29.9) ≈ **48 s vs 98.7 s — roughly
half the protocol-attributable time, a per-call ratio of ~0.5**. I will not defend this number:
n=3 per arm, arm-A range 27.3–105.3, two-point calibration, linearity assumed. It says only that
the expected effect is large enough to measure, not small enough to argue about.

**The smallest experiment that produces the real number before ship:**

1. Three arms on one fixed task: (a) no protocol, (b) generated packet for the task's phase,
   (c) full file. n=5 per arm — n=3 already proved too small for a stable effect size (codex-1's
   own caveat at `protocol-read-cost-regression/round-02/codex-1.md:14-16`). Randomized arm order,
   model/effort/output cap held constant, wall clock and prompt/cached tokens recorded — the
   controls already specified at `round-02/codex-1.md:237` of that idea (PRIMARY, locators; the
   control list is codex-1's, I am adopting not re-deriving it).
2. **One canary task** whose correct execution requires a rule the packet omits (e.g. an
   `auto_implement`-flagged idea needing §14). Arm (b) passes only if the agent loads the omitted
   section on demand or names it under Concerns. This is the operational form of claude-1's
   constraint that a packet omitting a needed rule is worse than a slow round — speed measured on
   tasks that never test the index would prove nothing about the failure mode that matters.
3. Pre-registered ship gate: packet-arm median ≤ ~60% of full-arm median, **and** canary passed.
   Fifteen invocations plus three canary runs; a day of wall clock, no code beyond the generator
   sketch. It must run before ship because after ship the full-read control arm no longer exists.

**Whole-idea honesty:** this change cuts only the protocol-read term. The other cost term — re-read
of prior rounds via `gatherPriorRounds`/`gatherReviewContext` — is the regression FINAL's rank 2
and is untouched here (PRIMARY — `protocol-read-cost-regression/FINAL.md`, rank-2 table). The
whole-idea saving will be smaller than the per-call saving, and this idea's FINAL must not claim
otherwise. My round-1 said the same; E8's context makes it worth repeating as a scope guard.

## D5 — Phase 8 is the right lever; cover both cells, asymmetrically

Yes, Phase 8 is the right lever — E3 settles it: the measured 7.2× byte growth lives in `review/`,
the worst idea ran 24 review rounds and 25 fix-up cycles, and in the driver shape those are **one
loop** ("The fix-up cycle number == the current review round; bounded by MaxFixupCycles" — E3
locator). A fix-up cap bounds exactly the loop that produced the measured tail. The premise that
Phase 8 bounds "one loop while another stays open" is true only of Phase 2.

On the other cell — §4.0 "Cross-review rounds (Phase 2): unbounded" for deliberation: observed
maximum across the deck is **3** (E-family command: per-idea `round-*` counts; the largest were 3,
on `integrate-parley-bidding-addon`, `session-resume-cache-plan`, `protocol-overlay-local-extension`),
design rounds measured flat (1.4 → 1.6, regression FINAL via 00-prompt), and the driver already
defaults deliberation cross-review to 1 while honoring an explicit 00-prompt value with no ceiling
(E5). That "unbounded" has never bitten anyone. **Counter-proposal to "bound fix-up alone":** cap
both cells in the same patch — Phase 8 at 5 because that is where the measured cost lives, Phase 2
at "escalate past 3" because leaving "unbounded" printed next to a capped cell preserves, inside
one table row, the exact text/tool disagreement this idea exists to close. But rank 3's cost
justification rests on Phase 8 alone; the Phase 2 cell is hygiene, and should be argued as such.

## Per-participant notes

**@codex-1** — I adopt your escalation payload and renewable-ceiling design unchanged; the
disagreement is only the number, answered in D1. Your D3 risk flag was right; D3 above decides it
your way. Your "the first cap must not wait for perfect data" is now stronger than you wrote it:
the data existed in the deck all along (E1/E2), it just had not been counted.

**@hermes-1** — D2 gives your builder path a precise scope and a hard precondition; please answer
the hash-check condition in your own words. On the cap, D1. One round-1 mapping disagreement I will
not let default through silence: your Phase 1 table lists §4.0 as reference-only; I hold it belongs
in every packet, because the table "OVERRIDES the full-lifecycle defaults" (COOPERATION.md:233-237,
PRIMARY — the bolded override paragraph) and an agent without it applies the wrong ceremony without
knowing. Your own Phase 8 row loads §4.0 for exactly this reason ("this is where the budget
lives"); the override logic is phase-independent. Counter-proposal: §4.0 in always-core, 58 lines.

**@claude-1** — your two standing constraints shaped D3 and D4: "design for the omitted-needed-rule
failure first" became the canary arm, and "what reports what a packet omits" stays answered by the
generated, never-curated index. On D4: the advance estimate you asked for does not exist as more
than an interpolation; the day-long experiment above is the deliverable that turns it into a number
before anything ships.

## Concerns / open questions

- **Counterfactual caution on E1/E2 (self-flagged):** the observed distribution was produced
  without a binding cap; under a cap, trajectories change (the 25-cycle idea might have replanned
  at cycle 5 and converged at 8 — or not). The caps-are-identical conclusion holds for the observed
  maxima; the escalation-rate prediction (~4/13 recent ideas) is a hypothesis the first months of
  telemetry should check.
- **Zero-cycle ambiguity:** 17 of 69 ideas recorded no `## Fix-up cycle` section; some predate the
  convention. E1's median is robust to that (the mass is at 0–1 either way), but the all-time mean
  is soft. The recent-cohort stats (E2) are the ones D1 leans on.
- **The Phase 2 cross-review ceiling needs a driver home too**: today `CapCrossReviewRounds` exists
  only for `standard` (E5). If D5's hygiene cap is adopted, deliberation needs the same clamp
  field set — otherwise we print "escalate past 3" and repeat E4's ignored-cap lesson in the same
  patch that fixes it.
