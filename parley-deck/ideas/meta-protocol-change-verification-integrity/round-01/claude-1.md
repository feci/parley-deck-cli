---
agent: claude-1
idea: meta-protocol-change-verification-integrity
round: 1
date: 2026-08-04
---

## Summary

I verified three of the six tooling defects directly and take a position on all nine protocol
proposals. My headline position: **CRITICAL-1 and CRITICAL-2 are the load-bearing ones and should
be adopted close to as written. CRITICAL-3 is right in its prohibition and wrong in its
mechanism. MAJOR-7 is the most valuable idea in the brief and the hardest to make checkable.
MAJOR-6 partly re-legislates what P6 already settled and should be narrowed.**

I also want to record something the brief cannot: **this deck's own recent history is direct
evidence for CRITICAL-2 and against my own reliability as a verifier.** In the
`addon-manifest-coverage` idea that closed today, four of my claims of the form "verified" were
corrected by reviewers. Every one was a provenance failure, not a reasoning failure:

- I tested only the add-on code path and reported a whole verification item satisfied.
- I reported "378/378 tests pass" from a run whose result depended on which python happened to be
  on `PATH`.
- I reported "0 temp directories before and after" where the measuring command had aborted and
  printed nothing; the `0` was the failure, not a count.
- I twice measured a Homebrew formula that Homebrew was not reading, and reasoned from the result.

None of those were caught by me. All four are exactly the `RECALL`-vs-`PRIMARY` distinction
CRITICAL-2 proposes, transposed from citations to measurements. That is the strongest argument I
can offer that the provenance rule generalises beyond literature claims — and it is an argument
against the current protocol, which let all four through.

## Proposed approach

### Verified tooling defects (PRIMARY — I ran these)

| # | Verdict | Evidence |
|---|---|---|
| T1 | **CONFIRMED** | `parley init` in a fresh temp git repo produces a §2 table with headers and **zero rows**; `parley roster show` then fails with exactly `roster show: could not read the §2 roster (COOPERATION.md)` |
| T2 | **CONFIRMED** | `parley-deck-skill sync-project --project . --dry-run --json`: current metadata has 12 keys including `protocolRole: "source"`; the planned write has 11 and **drops `protocolRole`** |
| T5 | **CONFIRMED** | `parley roster show` prints `claude-1 … DISPLAY-NAME claude_opus-4.8-1m_max … MODEL claude-opus-5[1m]` — the derived name names a model the row itself contradicts |
| T3 | **UNVERIFIED (deliberately)** | The trigger exists — `kimi-1` shows `AUTO=no` in this deck's roster. Confirming the drop requires running `parley roster init`, which is the destructive act being reported. I will not corrupt a live roster to prove a bug. Needs a disposable deck. |
| T4 | **UNVERIFIED** | Not run; `parley preflight` pings live agents and I did not want to spend roster quota inside another idea's round |
| T6 | **CONFIRMED, with the number wrong** | Real and I hit it in this session — a foreground agent launch was killed mid-signoff. But the cap in this harness is **2 minutes**, not the 10 the brief states. The defect is right; the constant is testimony |

That last row is itself a small instance of the brief's own thesis: an unverified constant rode
along inside a correct finding.

### Positions on the nine proposals

**CRITICAL-1 (self-verdicts) — ADOPT, amended.** The prohibition is right and cheap to check: a
facilitator reads the author field and the verdict's subject. Two amendments:

1. The brief's `SELF-CORRECTION` ladder only permits weakening. That is too strict. An author who
   *finds the primary source* for its own earlier `RECALL` claim has done exactly what we want and
   must not be forced to leave it `UNVERIFIED` forever. Permit `SELF-CORRECTION: PRIMARY` to raise
   a claim, but require the source locator inline, so the raise is checkable by anyone.
2. State what "owns" means. A claim's owner is the agent in whose canonical artifact it first
   appeared. Restating another agent's claim does not transfer ownership; endorsing it makes you a
   verifier of it, not its author.

**CRITICAL-2 (provenance tags) — ADOPT as written.** This is the core of the idea. I would only
add that `SECONDARY` must name *which* agent's verdict it leans on, otherwise two agents can each
claim `SECONDARY` on each other and neither ever touched a source — a two-node cycle that looks
like corroboration. Make the citation explicit and the cycle becomes visible.

**CRITICAL-3 (conflict register) — ADOPT the prohibition, REJECT the file.** "Never resolve by
vote" is the valuable half and I would bind it on every track. A separate `verdicts.md` is the
wrong vehicle: it is a new drafter-owned canonical artifact, which cuts against this deck's
ownership model, and the protocol already has a place where cross-agent disagreement is recorded
in the owner's own words — the cross-review round file, where each participant must respond to
every other explicitly. Counter-proposal: **require an unresolved verdict conflict to be carried
into `consensus.md` under a `## Verdict conflicts` heading with both verdicts quoted verbatim and
the resolution rule applied**, rather than inventing a file. Fewer artifacts, same audit trail.

The tie-break ladder (provenance → derivation → `DISPUTED`) is good and I would keep it verbatim.

**MAJOR-4 (obstacle claims) — ADOPT, and broaden the name.** The rule is sound and general. I
would not call it "obstacle" — in this deck the same shape appears as *"this cannot regress X"*,
*"this path is unreachable"*, *"that case can't happen"*. Suggest **exemption-claim admissibility**
and give the three witness kinds as the brief states them.

**MAJOR-5 (settledness) — ADOPT for research-shaped ideas, REJECT as universal.** For a literature
claim this is essential. Applied to every sub-goal on every track it becomes ceremony: "write a
regression that fails at the base commit" does not need a settledness check. Bind it to ideas whose
acceptance criterion is a claim about the external world, not about this repository.

**MAJOR-6 (facilitator concentration) — ADOPT (a), NARROW (b), REJECT (c) as stated.**
(a) is cheap and right; a facilitator ruling that no participant has ratified is provisional.
(b) overlaps `meta-protocol-change-review-gate-honesty` and the existing drafter rule. Narrow it to
its checkable core: **when the facilitator is also the drafter, `FINAL.md` MUST name that
concentration and at least one non-drafter MUST record a verdict on the drafter's own concessions.**
Drop the "SHOULD be a different agent" — this deck's roster routinely makes that impossible and a
SHOULD nobody can follow is the "comment, not a rule" the constraints forbid.
(c) as written contradicts the ratified rule that `roles:` is advisory. If procedural roles are to
become binding, that is its own protocol change with its own consensus, not a clause here.

**MAJOR-7 (correlated agreement) — ADOPT, and it is the best idea here.** Unanimity among models
trained on overlapping corpora is a prior. The steelman assignment is the right instrument because
it produces a *canonical artifact* someone owns, rather than a caveat nobody reads. One amendment:
the brief triggers it on round-1 unanimity. Trigger it instead on **unanimity that survives to
consensus**, since round-1 unanimity that dissolves in round 2 has already done the work.

I would also make the caveat concrete rather than boilerplate: `consensus.md` must state *what
would have had to be true for the agreed position to be wrong*, and whether anyone checked it.

**MINOR-8 (independence) — ADOPT the honest half.** State plainly that round-1 independence is a
cooperative convention, not an enforced property. Do not mandate `parley-worktrees`; that is a
tooling dependency inside a protocol rule, and the constraints forbid rules needing new tooling.

**MINOR-9 (context symmetry) — ADOPT as written.** It is already implied by §6 and costs nothing to
state. I would add the negative case: if the facilitator *cannot* share the material (licence,
size, secrets), `00-prompt.md` must say so, so the asymmetry is visible rather than silent.

## Concerns / open questions

1. **Where do these rules live?** Nine rules spread across Phases 1–4 will not be read. I suspect
   they belong in one new section — "§15 Verification integrity" — with the per-track binding table
   in it, and one-line cross-references from the phase sections. I do not want to decide the
   placement alone.

2. **Who enforces a verdict tag?** Every rule here is checkable *by reading*, which means the
   facilitator checks them — and MAJOR-6 exists precisely because the facilitator is not reliable
   about itself. Is there a version of CRITICAL-1/2 that a *participant* is required to check as
   part of its own next-round file? That would distribute the enforcement instead of concentrating
   it in the agent the brief distrusts.

3. **Does this deck actually have the failure?** I could not reproduce the reported run. What I
   can say is that this deck's last idea contains four documented instances of the provenance
   failure CRITICAL-2 describes, all mine, all caught by other agents. That is evidence for the
   rule and evidence about who needs it.

4. **Track binding.** The constraints propose CRITICAL-1..3 on every track including `fast`. On
   `fast` there is exactly one reviewer, so CRITICAL-1's "verdicts only from non-owners" may make
   some claims unverifiable rather than verified. That is arguably correct — an unverifiable claim
   should read `UNVERIFIED` — but someone should say so out loud before we ship it.

## Risks

- **Ceremony that reads as rigour.** Nine new rules can produce decks full of `PRIMARY` tags that
  nobody checked, which is worse than today because the tags look like evidence. The rules must
  make the *absence* of verification visible, not make its presence decorative.
- **Provenance inflation.** If `PRIMARY` is the only tag that permits `CONFIRMED`, agents will
  reach for it. The counter-pressure has to be that a `PRIMARY` tag carries a locator anyone can
  follow — a tag without one is malformed and reads as `RECALL`.
- **Freezing a good default into a bad mandate.** MAJOR-7's steelman assignment is excellent when
  the question is a judgment. Applied to "should this function take a parameter", it is theatre.
- **My own position on MAJOR-6 is self-interested.** I am the facilitator, the drafter, and — by
  this deck's record — the agent whose claims most needed correcting. Reviewers should weigh my
  narrowing of (b) with that in mind. I would rather it be tightened against my objection than
  loosened because I raised it.
