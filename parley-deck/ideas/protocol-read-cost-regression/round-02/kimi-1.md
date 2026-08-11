---
agent: kimi-1
idea: protocol-read-cost-regression
round: 2
date: 2026-08-10
---

## Summary

The felt slowdown is dominated by the review-cycle **count** term, not the read term — and the
count term is a step change dated 2026-07-28, the day the roster grew, not a gradual
protocol-growth effect. Read cost is real (claude-1's arm A/B: 3.3× per mandated-read call) but
it is the coefficient; the exponent is review rounds. The unbounded `deliberation` fix-up cap is
the enabler of the tail, not the driver: the two worst ideas kept finding **fresh MAJORs at
rounds 19–24**, which no severity floor would have stopped. On Q2: round N should send the
previous round in full plus a participant-owned carry-forward ledger, with the consensus drafter
— where §15.6's duty actually binds — keeping the full-history read. On Q4: against caveman
compression anywhere an agent acts on the text; it is a fourth drifting copy applied to a
document whose modal verbs are its content, and scoping dominates it on both savings and
fidelity.

Nothing from round 1 is restated without re-running it. New PRIMARY evidence this round, all my
own commands:

- **claude-1's era table reproduces exactly.** My re-run over all 76 ideas (`created:` from
  frontmatter; review rounds = `review/round-*` dirs; bytes = `cat` of review `.md` files):
  older 55 ideas — design 1.42, review 1.62, 20,237 B; newer 21 ideas — design 1.62, review
  5.14, 146,290 B. His numbers are solid.
- **But his 07-01 split hides a step at 07-28.** Within his "newer" bucket:
  07-01..07-27 = 9 ideas, mean review **2.33**, 30,853 B; 07-28..08-10 = 12 ideas, mean review
  **7.25**, 232,868 B. Early-July deliberation ideas ran 2–3 review rounds — under the same
  "unbounded" rule. The rule text did not change on 07-28; the roster did (PRIMARY:
  `git log -- parley-deck/agents.toml`: `d70c3f0` 07-28 "activate kimi-1, pin models",
  `6bb64cc` 07-30 "reactivate antigravity-1 as a fifth participant"). The §4.0 cap table itself
  entered 07-03 (PRIMARY: `git log -S unbounded` → `a224621` 2026-07-03), 25 days before the
  explosion it is blamed for.
- **The late findings are real and severe.** `skills-cli-install-path` (`standard`, 21 review
  rounds): round-19 file reports 2 fresh MAJORs + 1 MINOR; round-21 file reports 3 fresh MAJORs
  (`grep '### \['`). `integrate-parley-bidding-addon` (`deliberation`, `strict_gate`, 24 review
  rounds, 3 reviewers every round): rounds 22–24 produced 1 fresh MAJOR + 1 NIT. claude-1's
  mechanism-2 premise — "a MINOR found in round 19 reopens the cycle" — is not what the worst
  cases show. They show fresh MAJORs at rounds 19–24.
- **The quadratic is worse on the review side than the design side.** `gatherReviewContext`
  (internal/runner/phase58.go:278, comment at :276) embeds FINAL.md + IMPLEMENTATION.md + every
  prior review round into each reviewer call — on top of `gatherPriorRounds`
  (runner.go:936-989, verified) on the design side. Measured for the round-21 call of
  `skills-cli-install-path`: FINAL+IMPLEMENTATION 83,260 B + prior review rounds 1–20 at
  344,710 B = **427,970 B of embedded context per reviewer call**. IMPLEMENTATION.md alone is
  76,497 B carrying 15 `## Fix-up cycle` sections. codex-1's round-1 answer named
  `gatherPriorRounds` and `driver_consensus.go` but missed this third embedder; it is the
  largest of the three.
- **The CLI cap barely exists, and isn't what ran these loops.** track.go:149 caps `fast` at 1
  fix-up cycle, :159 caps `standard` at 2; the `Deliberation` case sets no `MaxFixupCycles`, so
  the driver default of 3 applies (driver.go:103-104). The protocol's "unbounded"
  (COOPERATION.md:220) binds only human-facilitated runs — and nothing since 2026-06-02 left
  runner telemetry (newest `parley-deck/runs/` dir: `20260602T195452`; re-verified), consistent
  with the entire slowdown period being skill-facilitated, where neither the driver caps nor the
  protocol caps were enforced by anyone.
- **Telemetry cannot currently settle Q1.** `events.jsonl` carries `duration_ms` but no token
  fields (PRIMARY: field listing of the newest run's log), and the log is dormant since 06-02.

## Q1 — what dominates, and how we would know

**I defend a dominant cause: the review-cycle count term, dated to the 07-28 roster seam.**
Arithmetic first (bytes are the only thing fully measured): 7.2× review-byte growth ≈ 3.2×
round count × ~2.2× per-round size. The protocol grew 2.12× over the same ten weeks — the
coefficient, exactly as claude-1's corrected framing says. Where I go further than claude-1:
his "newer era" is not a trend, it is a step. Everything before 07-28 looks like May–June
(2.33 rounds); everything after looks like a different system (7.25). Three candidate causes
moved inside that window: the roster (07-28/07-30), idea ambition (strict_gate, cross-repo
sources), and nothing in the protocol text — the cap rule was 25 days old and early-July ideas
under it were flat. A rule that permits long loops cannot explain a step that happened 25 days
after the rule arrived, on the exact day a fifth participant joined. Mechanism: more reviewers
→ more findings per round → more cycles to reach zero Agreed fixes; and each round completes at
the slowest of n (codex-1's `agents.toml` observation, "~20 min/round" pinned participant, is
the per-round version of the same effect).

Fermi check on felt time (ESTIMATE, not PRIMARY): claude-1's arm A/B delta is ~69 s/call for a
mandated full-protocol read; at ~40 agent calls per current-era idea that is ≈ 45 min/idea.
Each extra review cycle costs an implementer fix-up call + 2–4 reviewer calls carrying
monotonically growing embedded context (428 KB by round 21, above) + human gaps — tens of
minutes to hours each; +4.9 cycles vs early July is hours per idea. The count term plausibly
outweighs the read term several-fold, and it *multiplies* the read term: every extra round
re-pays the protocol read and carries a larger prior-round embedding.

**The ONE instrumented measurement that settles it:** per-invocation telemetry across one full
`standard` idea run end-to-end on the current roster — for every agent call record phase,
agent, model, wall duration, assembled-prompt bytes, bytes actually read via tool calls (or
provider input/cached/output tokens where the CLI exposes them), and the round wall window
(`max(agent durations)`, per codex-1). Then regress call duration on input vs output bytes with
agent fixed effects. Decision rule: if the input coefficient dominates, read-scoping is the
primary lever and I concede co-dominance; if output bytes + fixed effects dominate, the
write-side/count-side hypothesis holds and we fix counts first; if fixed effects swamp both,
codex-1's critical-path drift wins and neither scoping nor compression is the lever. The
plumbing half-exists (`duration_ms` in `events.jsonl`) but token capture is absent and the log
has been dormant since 06-02 — so step zero is running one idea under measurement at all.
Against codex-1's replay benchmark as the deciding instrument: frozen outputs measure calls the
system no longer makes; the owner's variable is felt wall-clock on live ideas. Replay is the
right second stage (isolating caching effects) after the in-vivo regression names the term.

## Q2 — what round N should send

Design side (`gatherPriorRounds`, runner.go:936-989):

- **Round 2: full round-1 files.** Converging with codex-1 — first non-anchored contact.
- **Round N≥3: the previous round in full + a per-participant carry-forward ledger +
  older rounds on demand.** The ledger is self-authored, fixed-format, machine-checkable for
  presence: position, open objections with stable IDs, `DISPUTED` flags verbatim, provenance
  locators, `supersedes` links. Missing/invalid ledger or a challenged locator fails back to
  full text, never to a silent summary.
- **The consensus drafter and signers keep the full-history read** (`buildConsensusDraftPrompt`,
  driver_consensus.go:112-113 — "Read EVERY round artifact under %s/round-*/", codex-1's
  citation confirmed in substance). This is my main addition to codex-1's design: §15.6's
  correlated-agreement duty is drafter-facing (COOPERATION.md:348, :1327), so preserving the
  full-text read exactly where the duty binds costs one call instead of n, and the detection
  surface is unchanged. What weakens: a round participant's own ability to spot convergence
  patterns across full history. I accept that trade; the participant's check is redundant with
  the drafter's, and any `DISPUTED` forces full-text reads anyway.

Review side (`gatherReviewContext`, phase58.go:278) — where the measured 428 KB lives:

- **Round N: FINAL.md + current IMPLEMENTATION.md in full** (the object under review —
  non-negotiable), **a findings ledger** (finding ID, severity, round introduced, status
  open/closed/follow-up, disposition, locator to the round file), **the previous review round
  in full** for fix verification, older rounds on demand. Review ledgers are nearly mechanical —
  reviewers already maintain them by hand (round-19 hermes-1: "Both round-18 MAJORs are
  confirmed closed, both NITs are confirmed closed, and the MINOR is correctly recorded as a
  follow-up").

What breaks, honestly: (a) quote-level rebuttal of older rounds — mitigated by one-command
expansion; (b) detection of hedged or self-contradictory claims across history — genuinely
weakened for round participants; (c) **re-litigation detection** — the Stopping judgment
(COOPERATION.md:637-647) needs to see "the same ground re-litigated despite open rebuttals", so
the ledger must carry dispositions and code-region tags or review termination loses its input;
(d) digest laundering — a ledger can upgrade a contested claim to settled, hence verbatim
verdict states + fail-closed fallback. Source-side bound, unchanged from round 1: cap round and
review file growth (cite line numbers, don't re-quote blocks); the quadratic is in bytes, and
bounded bytes/round bounds the sum even where full text is sent. claude-1's dichotomy (digests
fine for review, never for design) lands in the right place for the wrong reason: review
tolerates ledgers because findings are discrete objects with severities and dispositions, not
because review text matters less — and review is where his own data puts 7.2× of the bytes.

## Q3 — the unbounded `deliberation` cap

**It is the enabler, not the driver.** Evidence: (1) timing — the rule is 25 days older than the
explosion (above); (2) enforcement — under auto-drive, deliberation is capped at 3 by default
(driver.go:103-104), and "unbounded" (COOPERATION.md:220) is protocol text that binds only the
human-facilitated path these marathons actually ran on; (3) content — the worst cases kept
finding fresh MAJORs at rounds 19–24 (above), so a severity floor alone would not have ended
them; when round 21 of a shell-installer review finds "brace expansion builds a shell word",
the loop is working, and the defect is upstream (implementation density, idea risk). A fourth
datapoint: `integrate-parley-bidding-addon`'s 00-prompt.md today reads `status: round-01,
blocked-by: skills-cli-install-path` — 699,565 review bytes spent and the idea is parked behind
a file-set conflict. That is a sequencing failure, and no read-cost fix touches it.

Also worth noting for claude-1's table: `skills-cli-install-path`'s 21 "review rounds"
decompose into ~8 multi-reviewer rounds and ~13 single-reviewer codex-1 passes (rounds 3–15,
plus an empty round-14; `agy-1` joins at round 16 — the roster grew *mid-idea*). Directory
counts overstate full reviews; the conclusion survives, but future tables should count
reviewer-files, not directories.

**Termination rule proposal** — the protocol already has the right doctrine and it is advisory
and unused: "Stopping judgment" (COOPERATION.md:637-647) judges by trajectory, and budgets are
already "escalation thresholds, not close criteria" (:652-654). Make it binding with three
mechanical triggers:

1. **No-convergence stop:** two consecutive full review rounds whose fresh CRITICAL/MAJOR count
   is not falling → halt the loop and escalate with a trajectory summary. The human chooses:
   new fix-up plan, recorded operator ruling, or abandon. (This is the protocol's own
   "churning" example, promoted from illustrative to normative.)
2. **Re-litigation guard:** a finding on code whose disposition is already recorded must rebut
   the recorded disposition, else it does not reopen a cycle. (Gives the Q2 ledger its second
   job; direct response to "the same ground is re-litigated despite open rebuttals".)
3. **Severity floor + backstop:** after the track's fix-up budget (1/2/3), MINOR/NIT findings on
   non-`strict_gate` ideas are recorded follow-ups, not cycle-openers — formalizing what the
   deck already does informally. CRITICAL/MAJOR always reopen; `strict_gate` keeps its current
   semantics (NITs remain blocking, :622). Replace §4.0's "unbounded" for deliberation with a
   hard backstop (suggest 5 cycles) that escalates and never completes — aligning protocol text
   with the driver semantics that already exist.

Refutation-default review is untouched: reviewers must still try to break everything, every
round, and record their refutation attempts. The rule changes only what the system *does with*
findings — route, escalate — never whether they may be found or recorded. claude-1's
counter-case (severity-aware termination as the fix) answers a premise his own worst examples
refute; my counter-proposal terminates on *trajectory and re-litigation*, which is what those
marathons actually exhibited, and escalates rather than suppresses when the trajectory is
genuinely finding MAJORs — because then the right question is "why is fix-up cycle 15 still
bug-dense", not "how do we stop reviewing".

## Q4 — caveman-style prompt compression: against, honestly argued

For this system, specifically:

1. **It strips the load-bearing tokens.** The compressor's targets — auxiliaries, articles,
   connectives — are this document's normative content. Obligation density outgrew prose
   (round 1, PRIMARY: `MUST` 15→37, `MUST NOT` 6→15); RFC-2119 modality and article scope
   ("the facilitator" vs "a facilitator" in a multi-agent protocol) are exactly what a
   lossy stripper destroys. The strict-gate section depends on such precision: findings must be
   "objective, code-grounded", a "subjective stylistic preference is never a finding at any
   severity" (:617-620) — compress that sentence and review loses its definition.
2. **It manufactures a fourth copy of a document that already drifts in three.** hermes-1
   measured the bundled skill copy at 104,570 B vs the embedded 104,480 B (md5 mismatch);
   §9.0's provenance machinery exists because copies diverge. A compressed runtime rendering
   diverges invisibly — nobody diff-checks a lossy transform — straight into the two drift
   incidents claude-1 cited.
3. **It is dominated on its own metric.** Scoping removes 55–60 KB of 104.5 KB entirely (my
   round-1 phase analysis) with zero fidelity loss on what remains, and claude-1's arm B
   (excerpt) already showed 3.3×. Compression claims 40–58% token reduction on text it *keeps* —
   applied after scoping it would attack the small remainder at nonzero risk to every remaining
   rule. Verification asymmetry: excerpting keeps surviving text byte-identical and checkable;
   compression alters every line and is only statistically verifiable.
4. **Peer round files are the strongest case for it, and it loses there too.** They are
   analyses, not law — but they are evidence under §15: compression shifts line numbers,
   rewrites quotable text, and launders hedges into assertions, which is precisely the
   correlated-agreement failure §15.6 polices. The self-authored ledger (Q2) preserves
   `DISPUTED` verbatim *and* cuts more bytes.

Safe zone: human-facing status summaries to the owner, where nobody acts on modal precision.
Everywhere an agent acts on the text: against. **Deciding test:** a fidelity probe — ~20
normative edge-case questions (e.g., "may a participant silently relax a strict gate?" —
:625-626) answered by agents given compressed vs original text, requiring zero additional
errors — plus a head-to-head against excerpting at equal measured recall. Compression ships
only if it beats scoping on tokens at equal fidelity. Prediction: it cannot, because
scoped-full-fidelity text needs no compression to be smaller than compressed-everything.

## Q5 — ranking by (expected saving)/(risk to rules that catch real errors)

1. **Binding termination rule (Q3).** Attacks the dominant count term. Risk: moderate — a bad
   stop ships worse code — mitigated because CRITICAL/MAJOR always reopen, budgets escalate
   rather than complete, and `strict_gate` is untouched. Needs one §7 idea.
2. **Instruction-path fix.** SKILL.md:12/:24, the §9 checklist line (COOPERATION.md:857), and
   the facilitator's hand-written prompt habit must flip together to "load the core view; the
   facilitator names the phase-relevant reference sections" — hermes-1's lockstep point,
   extended: claude-1's inbox note shows the facilitator path is an instruction surface no
   generated artifact constrains, and this very round-2 prompt scoped our reads (4 files + 1
   note) and still got primary measurements — arm B working, again. Evidence: 3.3×/call. Risk:
   discovery failure; mitigated by the always-loaded safety core + Quickstart map.
3. **Scope the three embedders (Q2).** CLI-only, no §7 needed; attacks the measured 428
   KB/call. Risk: re-litigation blindness, digest laundering — mitigations as in Q2.
4. **Source-side length bounds** on round files and IMPLEMENTATION.md fix-up append growth.
   §7 template change; bounds the quadratic in place; risk is gaming a style cap — frame as
   protocol rule.
5. **Phase-scoped compiler / `render --scope`** (codex-1's registry, hermes-1's flag). Largest
   per-call saving, largest build, depends on the unbuilt core store. Sequence after 1–4 bleed
   off urgency; do not let it block them.
6. **Caveman compression: do not ship** (Q4).

**Must never be cut:** §15 in full (provenance, conflicting verdicts, role concentration,
correlated agreement — the checks that catch real errors; claude-1 credits it with three of his
own this month), §7 change control, §6 rule 3 (never edit another agent's file), §14 human
brake, §1 non-solo, refutation-default review with recorded refutation attempts, no-suppression
and `strict_gate` semantics (:611-626), "budgets escalate, never complete" (:652-654),
append-only signoffs, the drafter's full-history read (Q2), and the human escalation path as
terminator of last resort. Plus one process rule: never ship an optimization before the
measurement that justifies it (Q1) — that rule is what this idea is for.

## Cost of this round's restraint (required disclosure)

This round I read of COOPERATION.md only: the §4.0 table rows (:216-221), the strict-gate and
Stopping judgment text (:605-660), and grep-hit lines (:348, :622, :857, :1327) — roughly 6 KB
of 104.5 KB, plus code and artifacts as cited above. What it cost me: my §6-rule-3 and §14
endorsements remain inherited from round-1 characterizations, not re-verified against their
text; I read §15.6's location and drafter-facing pointer (:348, :1327) but not its full body
this round, so the Q2 claim that correlated-agreement duties bind the drafter rests on :348's
one-line characterization; and I did not read the LE-5 continuation past :660. Nothing in the
five answers depends on text I skipped beyond those three flags. Two restraint rounds, two full
analyses with primary measurements — the scoped-read pattern the owner is paying us to test is
the one this idea's own prompts keep demonstrating.
